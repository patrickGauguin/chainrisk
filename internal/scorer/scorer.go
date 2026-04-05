package scorer

import (
	"math"

	"github.com/patrickGauguin/chainrisk/internal/types"
)

func ScorePackage(vulns []types.Vulnerability, info types.PackageInfo) int {

	criticalCount := 0
	highCount := 0
	mediumCount := 0
	lowCount := 0

	for _, vuln := range vulns {
		switch vuln.Severity {
		case "CRITICAL":
			criticalCount += 1
		case "HIGH":
			highCount += 1
		case "MEDIUM":
			mediumCount += 1
		case "LOW":
			lowCount += 1
		case "UNKNOWN":
			mediumCount += 1
		}
	}

	criticalScore := int(90 * math.Log1p(float64(criticalCount)))
	highScore := int(25 * math.Log1p(float64(highCount)))
	mediumScore := int(12 * math.Log1p(float64(mediumCount)))
	lowScore := int(4 * math.Log1p(float64(lowCount)))

	securityScore := criticalScore + highScore + mediumScore + lowScore

	notLatest := !info.IsDefault

	maintainerScore := 0
	if info.IsDeprecated {
		maintainerScore += 15
		if notLatest {
			maintainerScore += 5
		}
	}

	switch true {
	case info.DaysSincePublish > 730:
		maintainerScore += 5
		if notLatest {
			maintainerScore += 2
		}
	case info.DaysSincePublish > 365:
		maintainerScore += 3
	case info.DaysSincePublish > 180:
		maintainerScore += 2
	}

	total := securityScore + maintainerScore

	if total > 100 {
		total = 100
	}

	return total
}

func RiskLevel(score int) string {
	switch {
	case score >= 75:
		return "CRITICAL"
	case score >= 50:
		return "HIGH"
	case score >= 25:
		return "MEDIUM"
	case score >= 10:
		return "LOW"
	default:
		return "SAFE"
	}
}
