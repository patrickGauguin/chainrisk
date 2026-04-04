package report

import (
	"fmt"
	"sort"

	"github.com/patrickGauguin/chainrisk/internal/types"
)

func PrintTerminal(result *types.ScanResult) error {
	fmt.Printf("NAME: %s | LANGUAGE: %s | STARS: %d | LAST_PUSHED: %s:\n", result.Repo.FullName, result.Repo.Language, result.Repo.Stars, result.Repo.LastPushed.Format("2006-01-12"))

	var prodPackages []types.PackageRisk
	var devPackages []types.PackageRisk

	for _, pkg := range result.Packages {
		if pkg.Dependency.IsDev {
			devPackages = append(devPackages, pkg)
		} else {
			prodPackages = append(prodPackages, pkg)
		}
	}

	sort.Slice(prodPackages, func(i, j int) bool {
		return prodPackages[i].Score > prodPackages[j].Score
	})

	sort.Slice(devPackages, func(i, j int) bool {
		return devPackages[i].Score > devPackages[j].Score
	})

	fmt.Printf("PRODUCTION DEPENDENCIES (%d packages)\n", len(prodPackages))
	threshold := pickThreshold(prodPackages)
	printPackages(filterByScore(prodPackages, threshold))

	fmt.Printf("DEV DEPENDENCIES (%d packages)\n", len(devPackages))
	threshold = pickThreshold(devPackages)
	printPackages(filterByScore(devPackages, threshold))

	return nil
}

func filterByScore(packages []types.PackageRisk, minScore int) []types.PackageRisk {
	var result []types.PackageRisk
	for _, pkg := range packages {
		if pkg.Score >= minScore {
			result = append(result, pkg)
		}
	}
	return result
}

func pickThreshold(packages []types.PackageRisk) int {
	for _, pkg := range packages {
		if pkg.Score >= 25 {
			return 25
		}
	}
	for _, pkg := range packages {
		if pkg.Score >= 10 {
			return 10
		}
	}
	return 0
}

func printPackages(packages []types.PackageRisk) {
	for _, pkg := range packages {
		fmt.Printf("  [name - %s][version - %s][score - %d][risk - %s]:\n", pkg.Dependency.Name, pkg.Dependency.Version, pkg.Score, pkg.RiskLevel)

		if len(pkg.Vulns) > 0 {
			for _, v := range pkg.Vulns {
				fmt.Printf("      [%s] %s - %s\n", v.Severity, v.ID, v.Summary)
			}
		}
	}
}
