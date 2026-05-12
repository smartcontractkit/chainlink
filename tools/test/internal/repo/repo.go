package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const rootModulePath = "github.com/smartcontractkit/chainlink/v2"

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
			if mod, ok := modulePathFromGoMod(string(data)); ok && mod == rootModulePath {
				// Exact root module only (not tools/test or other nested modules).
				return dir, nil
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("chainlink repo root not found (no go.mod with module %s) starting from %q", rootModulePath, dir)
		}
		dir = parent
	}
}

// modulePathFromGoMod returns the module path from the first `module` directive,
// skipping leading comments and blank lines (go.mod may legally start with either).
func modulePathFromGoMod(data string) (path string, ok bool) {
	for raw := range strings.SplitSeq(data, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if i := strings.Index(line, "//"); i >= 0 {
			line = strings.TrimSpace(line[:i])
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "module" {
			return fields[1], true
		}
	}
	return "", false
}
