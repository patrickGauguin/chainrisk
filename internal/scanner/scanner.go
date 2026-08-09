package scanner

import (
	"fmt"
	"sync"
	"time"

	"github.com/patrickGauguin/chainrisk/internal/depsdev"
	"github.com/patrickGauguin/chainrisk/internal/github"
	"github.com/patrickGauguin/chainrisk/internal/osv"
	"github.com/patrickGauguin/chainrisk/internal/parser"
	"github.com/patrickGauguin/chainrisk/internal/scorer"
	"github.com/patrickGauguin/chainrisk/internal/types"
)

type Scanner struct {
	Client *github.Client
}

func New(token string) *Scanner {
	return &Scanner{
		Client: github.New(token),
	}
}

func (s *Scanner) Scan(repoURL string) (*types.ScanResult, error) {
	owner, repo, err := github.ParseOwnerRepo(repoURL)
	if err != nil {
		return nil, err
	}

	repoInfo, err := s.Client.GetRepo(owner, repo)
	if err != nil {
		return nil, err
	}

	manifests := []struct {
		path  string
		parse func(string) (types.PackageMeta, []types.Dependency, error)
	}{
		{"package.json", parser.ParsePackageJSON},
		{"go.mod", parser.ParseGoMod},
	}

	var pkgMeta types.PackageMeta
	var deps []types.Dependency
	for _, manifest := range manifests {
		content, err := s.Client.GetFileContent(owner, repo, manifest.path)
		if err != nil {
			continue
		}

		pkg, dependencies, err := manifest.parse(content)
		if err != nil {
			continue
		}

		pkgMeta = pkg
		deps = dependencies
		break
	}

	if deps == nil {
		return nil, fmt.Errorf("no supported manifest file found (package.json, go.mod)")
	}

	for _, dep := range deps {
		depsdev.GetPackageVersion(dep.Ecosystem, dep.Name, dep.Version)
	}

	vulnMap, err := osv.LookupVulnerabilities(deps)
	if err != nil {
		return nil, err
	}

	packagesRisk := []types.PackageRisk{}

	type pkgResult struct {
		pkgRisk types.PackageRisk
		err     error
	}

	ch := make(chan pkgResult, len(deps))
	var wg sync.WaitGroup

	for _, dep := range deps {
		wg.Add(1)
		go func(dep types.Dependency) {
			defer wg.Done()
			pkgVersion, err := depsdev.GetPackageVersion(dep.Ecosystem, dep.Name, dep.Version)
			if err != nil {
				return
			}

			published, parseErr := time.Parse(time.RFC3339, pkgVersion.PublishedAt)
			if parseErr != nil {
				published = time.Time{}
			}

			days := int(time.Since(published).Hours() / 24)

			pkgInfo := types.PackageInfo{PublishedAt: published, IsDefault: pkgVersion.IsDefault, IsDeprecated: pkgVersion.IsDeprecated, DeprecatedReason: pkgVersion.DeprecatedReason, DaysSincePublish: days}

			vulns := vulnMap[dep.Name]
			score := scorer.ScorePackage(vulns, pkgInfo)
			risk := scorer.RiskLevel(score)

			pkgRisk := types.PackageRisk{Dependency: dep, Vulns: vulns, Info: pkgInfo, Score: score, RiskLevel: risk}

			ch <- pkgResult{pkgRisk, nil}
		}(dep)
	}

	wg.Wait()
	close(ch)

	for r := range ch {
		if r.err != nil {
			return nil, r.err
		}

		packagesRisk = append(packagesRisk, r.pkgRisk)
	}

	scanResult := types.ScanResult{Repo: repoInfo, Packages: packagesRisk, PackageMeta: pkgMeta}

	return &scanResult, err
}
