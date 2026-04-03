package parser

import (
	"encoding/json"

	"github.com/patrickGauguin/chainrisk/internal/types"
)

type packageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func ParsePackageJSON(content string) ([]types.Dependency, error) {
	var pkg packageJSON

	err := json.Unmarshal([]byte(content), &pkg)
	if err != nil {
		return nil, err
	}

	dependencies := []types.Dependency{}
	for name, version := range pkg.Dependencies {
		dependency := types.Dependency{Name: name, Version: version, Ecosystem: "npm"}
		dependencies = append(dependencies, dependency)
	}

	for name, version := range pkg.DevDependencies {
		dependency := types.Dependency{Name: name, Version: version, Ecosystem: "npm", IsDev: true}
		dependencies = append(dependencies, dependency)
	}

	return dependencies, err
}
