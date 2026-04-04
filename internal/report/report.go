package report

import (
	"fmt"

	"github.com/patrickGauguin/chainrisk/internal/types"
)

func PrintTerminal(result *types.ScanResult) error {
	fmt.Printf("NAME: %s | LANGUAGE: %s | STARS: %d | LAST_PUSHED: %s:\n", result.Repo.FullName, result.Repo.Language, result.Repo.Stars, result.Repo.LastPushed.Format("2006-01-12"))

	for _, pkg := range result.Packages {
		if pkg.Score > 0 {
			fmt.Printf("  [name - %s][version - %s][score - %d][risk - %s]:\n", pkg.Dependency.Name, pkg.Dependency.Version, pkg.Score, pkg.RiskLevel)

			if len(pkg.Vulns) > 0 {
				for _, v := range pkg.Vulns {
					fmt.Printf("      [%s] %s - %s\n", v.Severity, v.ID, v.Summary)
				}
			}
		}
	}

	return nil
}
