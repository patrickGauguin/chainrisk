package osv

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	ID               string           `json:"id"`
	Summary          string           `json:"summary"`
	Severity         []osvSeverity    `json:"severity"`
	DataBaseSpecific dataBaseSpecific `json:"database_specific"`
}

type osvSeverity struct {
	Type  string `json:"type"`
	Score string `json:"score"`
}

type dataBaseSpecific struct {
	Severity string `json:"severity"`
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
	fmt.Println(string(respJSON[:500])) // first 500 chars

	batchResponse := batchResponse{}
	json.Unmarshal(respJSON, &batchResponse)

	vulnMap := map[string][]types.Vulnerability{}

	for i, result := range batchResponse.Results {
		name := deps[i].Name
		vulns := []types.Vulnerability{}

		for _, vuln := range result.Vulns {
			newVuln, err := fetchVulnDetails(vuln.ID)
			if err != nil {
				return nil, err
			}

			vulns = append(vulns, newVuln)
		}

		vulnMap[name] = vulns
	}

	return vulnMap, err
}

func fetchVulnDetails(id string) (types.Vulnerability, error) {
	url := "https://api.osv.dev/v1/vulns/" + id

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return types.Vulnerability{}, err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return types.Vulnerability{}, err
	}

	defer resp.Body.Close()

	respJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.Vulnerability{}, err
	}

	type smallOsvVuln struct {
		ID               string           `json:"id"`
		Summary          string           `json:"summary"`
		DataBaseSpecific dataBaseSpecific `json:"database_specific"`
	}

	smallVuln := smallOsvVuln{}
	json.Unmarshal(respJSON, &smallVuln)

	severity := smallVuln.DataBaseSpecific.Severity
	if severity == "MODERATE" {
		severity = "MEDIUM"
	}
	if severity == "" {
		severity = "UNKNOWN"
	}

	vuln := types.Vulnerability{ID: smallVuln.ID, Summary: smallVuln.Summary, Severity: severity}

	return vuln, err
}
