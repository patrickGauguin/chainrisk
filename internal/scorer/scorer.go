package scorer

import (
	"math"

	"github.com/patrickGauguin/chainrisk/internal/types"
)

func ScorePackage(vulns []types.Vulnerability, daysSinceLastPush int) int {

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

	maintainerScore := 0
	switch true {
	case daysSinceLastPush > 730:
		maintainerScore = 30
	case daysSinceLastPush > 365:
		maintainerScore = 20
	case daysSinceLastPush > 180:
		maintainerScore = 10
	default:
		maintainerScore = 0
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
