package depsdev

import (
	"encoding/json"
	"io"
	"net/http"
)

type PackageVersion struct {
	PublishedAt      string `json:"publishedAt"`
	IsDefault        bool   `json:"isDefault"`
	IsDeprecated     bool   `json:"isDeprecated"`
	DeprecatedReason string `json:"deprecatedReason"`
}

func GetPackageVersion(ecosystem, name, version string) (PackageVersion, error) {
	url := "https://api.deps.dev/v3/systems/" + ecosystem + "/packages/" + name + "/versions/" + version

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PackageVersion{}, err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return PackageVersion{}, err
	}

	defer resp.Body.Close()

	respJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return PackageVersion{}, err
	}

	packageVersion := PackageVersion{}
	json.Unmarshal(respJSON, &packageVersion)

	return packageVersion, err
}
