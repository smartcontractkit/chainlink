package tontestutils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/mod/modfile"
)

func GetTONSha() (version string, err error) {
	modFilePath, err := getModFilePath()
	if err != nil {
		return "", err
	}
	depVersion, err := getTONCcipDependencyVersion(modFilePath)
	if err != nil {
		return "", err
	}
	tokens := strings.Split(depVersion, "-")
	if len(tokens) == 3 {
		version := tokens[len(tokens)-1]

		return version, nil
	}

	return "", fmt.Errorf("invalid go.mod version: %s", depVersion)
}

func getTONCcipDependencyVersion(gomodPath string) (string, error) {
	const dependency = "github.com/smartcontractkit/chainlink-ton"

	gomod, err := os.ReadFile(gomodPath)
	if err != nil {
		return "", err
	}

	modFile, err := modfile.ParseLax("go.mod", gomod, nil)
	if err != nil {
		return "", err
	}

	for _, dep := range modFile.Require {
		if dep.Mod.Path == dependency {
			return dep.Mod.Version, nil
		}
	}

	return "", fmt.Errorf("dependency %s not found", dependency)
}

func getModFilePath() (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("runtime.Caller failed to retrieve current file")
	}

	// Get the root directory by walking up from current file until we find go.mod
	rootDir := filepath.Dir(currentFile)
	for {
		if _, err := os.Stat(filepath.Join(rootDir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(rootDir)
		if parent == rootDir {
			return "", errors.New("could not find project root directory containing go.mod")
		}
		rootDir = parent
	}
	return filepath.Join(rootDir, "go.mod"), nil
}
