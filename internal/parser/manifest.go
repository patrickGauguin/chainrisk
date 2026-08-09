package parser

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/patrickGauguin/chainrisk/internal/types"
)

type packageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func ParsePackageJSON(content string) (types.PackageMeta, []types.Dependency, error) {
	var pkg packageJSON

	if content == "" {
		return types.PackageMeta{}, nil, fmt.Errorf("no package.json found ParsePackageJSON")
	}

	err := json.Unmarshal([]byte(content), &pkg)
	if err != nil {
		return types.PackageMeta{}, nil, err
	}

	pkgMeta := types.PackageMeta{Name: pkg.Name, Version: pkg.Version, Ecosystem: "npm"}

	dependencies := []types.Dependency{}
	for name, version := range pkg.Dependencies {
		dependency := types.Dependency{Name: name, Version: cleanVersion(version), Ecosystem: "npm"}
		dependencies = append(dependencies, dependency)
	}

	for name, version := range pkg.DevDependencies {
		dependency := types.Dependency{Name: name, Version: cleanVersion(version), Ecosystem: "npm", IsDev: true}
		dependencies = append(dependencies, dependency)
	}

	return pkgMeta, dependencies, err
}

func ParseRequirementsTxt(content string) (types.PackageMeta, []types.Dependency, error) {
	return types.PackageMeta{}, nil, nil
}

func ParseGoMod(content string) (types.PackageMeta, []types.Dependency, error) {
	var pkgData types.PackageMeta
	var dependencies []types.Dependency

	lines := strings.Split(content, "\n")
	inBlock := false
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		} else {
			switch fields[0] {
			case "module":
				pkgData = makePkgMeta(fields[1])
			case "require":
				if fields[1] == "(" {
					inBlock = true
					continue
				} else {
					dependencies = append(dependencies, parseDependency(fields[1:]))
				}
			case ")":
				inBlock = false
				continue
			default:
				if inBlock {
					dependencies = append(dependencies, parseDependency(fields))
				} else {
					continue
				}
			}
		}
	}

	return pkgData, dependencies, nil
}

func makePkgMeta(field string) types.PackageMeta {
	pkgMeta := types.PackageMeta{Name: field, Version: "", Ecosystem: "go"}
	return pkgMeta
}

func parseDependency(fields []string) types.Dependency {
	version := fields[1]
	isDev := len(fields) >= 4 && fields[3] == "indirect"

	dependency := types.Dependency{Name: fields[0], Version: version, Ecosystem: "go", IsDev: isDev}

	return dependency
}

func ParseCargoToml(content string) (types.PackageMeta, []types.Dependency, error) {
	return types.PackageMeta{}, nil, nil
}

func cleanVersion(v string) string {
	v = strings.TrimPrefix(v, "^")
	v = strings.TrimPrefix(v, "~")
	v = strings.TrimPrefix(v, ">=")
	v = strings.TrimPrefix(v, "<=")
	v = strings.TrimPrefix(v, ">")
	return strings.TrimSpace(v)
}
