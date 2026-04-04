package scanner

import (
	"time"

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

	vulnMap, err := osv.LookupVulnerabilities(deps)
	if err != nil {
		return nil, err
	}

	packagesRisk := []types.PackageRisk{}

	for _, dep := range deps {
		vulns := vulnMap[dep.Name]
		days := int(time.Since(repoInfo.LastPushed).Hours() / 24)
		score := scorer.ScorePackage(vulns, days)
		risk := scorer.RiskLevel(score)

		packageRisk := types.PackageRisk{Dependency: dep, Vulns: vulns, Score: score, RiskLevel: risk}

		packagesRisk = append(packagesRisk, packageRisk)
	}

	scanResult := types.ScanResult{Repo: repoInfo, Packages: packagesRisk}

	return &scanResult, err
}
