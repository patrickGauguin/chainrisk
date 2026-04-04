package osv

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/patrickGauguin/chainrisk/internal/types"
)

type queryRequest struct {
	Queries []osvQuery `json:"queries"`
}

type osvQuery struct {
	Package osvPackage `json:"package"`
	Version string     `json:"version"`
}

type osvPackage struct {
	Name      string `json:"name"`
	Ecosystem string `json:"ecosystem"`
}

type batchResponse struct {
	Results []queryResult `json:"results"`
}

type queryResult struct {
	Vulns []osvVuln `json:"vulns"`
}

type osvVuln struct {
	ID       string        `json:"id"`
	Summary  string        `json:"summary"`
	Severity []osvSeverity `json:"severity"`
}

type osvSeverity struct {
	Type  string `json:"type"`
	Score string `json:"score"`
}

func LookupVulnerabilities(deps []types.Dependency) (map[string][]types.Vulnerability, error) {
	queryRequest := queryRequest{}
	queryList := []osvQuery{}

	for _, dep := range deps {
		osvPackage := osvPackage{Name: dep.Name, Ecosystem: dep.Ecosystem}

		osvQuery := osvQuery{Package: osvPackage, Version: dep.Version}
		queryList = append(queryList, osvQuery)
	}

	queryRequest.Queries = queryList

	req, err := json.Marshal(queryRequest)
	if err != nil {
		return nil, err
	}

	url := "https://api.osv.dev/v1/querybatch"
	resp, err := http.Post(url, "application/json", bytes.NewReader(req))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, err
	}

	respJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	batchResponse := batchResponse{}
	json.Unmarshal(respJSON, &batchResponse)

	vulnMap := map[string][]types.Vulnerability{}

	for i, result := range batchResponse.Results {
		name := deps[i].Name
		vulns := []types.Vulnerability{}

		for _, vuln := range result.Vulns {
			severity := "UNKNOWN"
			if len(vuln.Severity) > 0 {
				severity = vuln.Severity[0].Score
			}

			vulns = append(vulns, types.Vulnerability{ID: vuln.ID, Severity: severity, Summary: vuln.Summary})
		}

		vulnMap[name] = vulns
	}

	return vulnMap, err
}
