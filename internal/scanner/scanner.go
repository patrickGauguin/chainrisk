package scanner

import (
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

	content, err := s.Client.GetFileContent(owner, repo, "package.json")
	if err != nil {
		return nil, err
	}

	deps, err := parser.ParsePackageJSON(content)
	if err != nil {
		return nil, err
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

	scanResult := types.ScanResult{Repo: repoInfo, Packages: packagesRisk}

	return &scanResult, err
}
