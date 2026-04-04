package scorer

import (
	"github.com/patrickGauguin/chainrisk/internal/types"
)

func ScorePackage(vulns []types.Vulnerability, daysSinceLastPush int) int {

	securityScore := 0
	for _, vuln := range vulns {
		switch vuln.Severity {
		case "CRITICAL":
			securityScore += 40
		case "HIGH":
			securityScore += 20
		case "MEDIUM":
			securityScore += 8
		case "LOW":
			securityScore += 2
		case "UNKNOWN":
			securityScore += 5
		default:
			securityScore += 0
		}
	}

	securityScore = int(float64(securityScore) * 0.6)

	maintainerScore := 0
	switch true {
	case daysSinceLastPush > 730:
		maintainerScore = 40
	case daysSinceLastPush > 365:
		maintainerScore = 25
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
