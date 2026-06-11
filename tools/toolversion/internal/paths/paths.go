package paths

import (
	"os"
	"path/filepath"
)

// Config holds manifest file locations.
type Config struct {
	Root             string
	ToolVersionsFile string
	GoToolsFile      string
	Makefile         string
	GoMod            string
}

// FromEnv resolves paths from CHAINLINK_ROOT (or repo root discovery) and manifest overrides.
func FromEnv() (Config, error) {
	root := os.Getenv("CHAINLINK_ROOT")
	if root == "" {
		var err error
		root, err = findRepoRoot()
		if err != nil {
			return Config{}, err
		}
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return Config{}, err
	}

	tv := os.Getenv("TOOL_VERSIONS_FILE")
	if tv == "" {
		tv = filepath.Join(root, ".tool-versions")
	}
	gt := os.Getenv("GO_TOOLS_FILE")
	if gt == "" {
		gt = filepath.Join(root, "tools", "go-tools.txt")
	}

	return Config{
		Root:             root,
		ToolVersionsFile: tv,
		GoToolsFile:      gt,
		Makefile:         filepath.Join(root, "GNUmakefile"),
		GoMod:            filepath.Join(root, "go.mod"),
	}, nil
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return findRepoRootFrom(dir)
}

func findRepoRootFrom(dir string) (string, error) {
	for {
		if fileExists(filepath.Join(dir, "go.mod")) && fileExists(filepath.Join(dir, ".tool-versions")) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", os.ErrNotExist
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
