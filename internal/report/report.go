package report

import (
	"fmt"
	"sort"

	"github.com/patrickGauguin/chainrisk/internal/types"
)

func PrintTerminal(result *types.ScanResult) error {
	var prodPackages []types.PackageRisk
	var devPackages []types.PackageRisk

	for _, pkg := range result.Packages {
		if pkg.Dependency.IsDev {
			devPackages = append(devPackages, pkg)
		} else {
			prodPackages = append(prodPackages, pkg)
		}
	}

	fmt.Printf("NAME: %s | LANGUAGE: %s | STARS: %d | LAST_PUSHED: %s:\n", result.Repo.FullName, result.Repo.Language, result.Repo.Stars, result.Repo.LastPushed.Format("2006-01-12"))

	printSummary(prodPackages, devPackages)

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
	severityRank := map[string]int{
		"CRITICAL": 0,
		"HIGH":     1,
		"MEDIUM":   2,
		"LOW":      3,
		"UNKNOWN":  4,
	}

	for _, pkg := range packages {
		fmt.Printf("  [name - %s][version - %s][score - %d][risk - %s]:\n", pkg.Dependency.Name, pkg.Dependency.Version, pkg.Score, pkg.RiskLevel)

		if len(pkg.Vulns) > 0 {
			sort.Slice(pkg.Vulns, func(i, j int) bool {
				return severityRank[pkg.Vulns[i].Severity] < severityRank[pkg.Vulns[j].Severity]
			})

			for _, v := range pkg.Vulns {
				fmt.Printf("      [%s] %s - %s\n", v.Severity, v.ID, v.Summary)
			}
		}
	}
}

func printSummary(prod, dev []types.PackageRisk) {
	riskRank := map[string]int{
		"CRITICAL": 4,
		"HIGH":     3,
		"MEDIUM":   2,
		"LOW":      1,
		"SAFE":     0,
	}

	highestRisk := "SAFE"
	for _, pkg := range prod {
		if riskRank[pkg.RiskLevel] > riskRank[highestRisk] {
			highestRisk = pkg.RiskLevel
		}
	}

	criticalCount := 0
	highCount := 0
	mediumCount := 0
	lowCount := 0
	safeCount := 0

	allPackages := prod

	for _, pkg := range dev {
		allPackages = append(allPackages, pkg)
	}

	for _, pkg := range allPackages {
		switch pkg.RiskLevel {
		case "CRITICAL":
			criticalCount++
		case "HIGH":
			highCount++
		case "MEDIUM":
			mediumCount++
		case "LOW":
			lowCount++
		case "SAFE":
			safeCount++
		}
	}

	fmt.Printf("OVERALL RISK: %s\n", highestRisk)
	fmt.Println("─────────────────────────────")
	fmt.Printf("Total packages		%d\n", len(prod)+len(dev))
	fmt.Printf("CRITICAL		%d\n", criticalCount)
	fmt.Printf("HIGH			%d\n", highCount)
	fmt.Printf("MEDIUM			%d\n", mediumCount)
	fmt.Printf("LOW			%d\n", lowCount)
	fmt.Printf("SAFE			%d\n", safeCount)
}
