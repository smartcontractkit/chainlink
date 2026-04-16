package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const rootModuleLine = "module github.com/smartcontractkit/chainlink/v2"

// RootFromWd walks parents from the working directory to find the chainlink v2 repo root.
func RootFromWd() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return RootFrom(wd)
}

// RootFrom walks parents from dir until go.mod declares the chainlink v2 module.
func RootFrom(dir string) (string, error) {
	dir = filepath.Clean(dir)
	for {
		modPath := filepath.Join(dir, "go.mod")
		data, err := os.ReadFile(modPath)
		if err == nil {
			firstLine := strings.TrimSpace(strings.Split(string(data), "\n")[0])
			// Exact root module only (not tools/test or other nested modules).
			if firstLine == rootModuleLine {
				return dir, nil
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("chainlink repo root not found (no go.mod with %s) starting from %q", rootModuleLine, dir)
		}
		dir = parent
	}
}
