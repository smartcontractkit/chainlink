package tools

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Target represents a runnable Go test target.
type Target struct {
	Name        string `json:"name"`
	Dir         string `json:"dir"`
	Packages    string `json:"packages"`
	IsSubmodule bool   `json:"is_submodule"`
}

// DiscoverTargets scans toolsDir relative to repoRoot to find standalone Go submodules
// and root-module packages located under toolsDir.
func DiscoverTargets(repoRoot, toolsDir string) ([]Target, error) {
	var submodules []Target
	var hasRootTests bool

	// Find all submodules with go.mod
	err := filepath.WalkDir(toolsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if d.Name() == "go.mod" {
			dir := filepath.Dir(path)
			relDir, err := filepath.Rel(repoRoot, dir)
			if err != nil {
				return err
			}
			submodules = append(submodules, Target{
				Name:        relDir,
				Dir:         relDir,
				Packages:    "./...",
				IsSubmodule: true,
			})
			return nil
		}

		if strings.HasSuffix(d.Name(), "_test.go") {
			// Check if this test file belongs to a submodule or root module
			// A file belongs to root if no parent directory between toolsDir and file has go.mod
			cur := filepath.Dir(path)
			inSubmodule := false
			for {
				if _, err := os.Stat(filepath.Join(cur, "go.mod")); err == nil {
					inSubmodule = true
					break
				}
				if cur == toolsDir || cur == repoRoot || cur == "." || cur == "/" {
					break
				}
				parent := filepath.Dir(cur)
				if parent == cur {
					break
				}
				cur = parent
			}
			if !inSubmodule {
				hasRootTests = true
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	var targets []Target
	if hasRootTests {
		targets = append(targets, Target{
			Name:        "tools/root",
			Dir:         ".",
			Packages:    "./tools/...",
			IsSubmodule: false,
		})
	}

	sort.Slice(submodules, func(i, j int) bool {
		return submodules[i].Name < submodules[j].Name
	})
	targets = append(targets, submodules...)

	return targets, nil
}

// MatrixOptions contains filtering parameters for matrix generation.
type MatrixOptions struct {
	EventName    string
	ChangedFiles []string
	All          bool
}

// ComputeMatrix filters discovered targets against changed files and event triggers.
func ComputeMatrix(targets []Target, opts MatrixOptions) []Target {
	if opts.All || opts.EventName == "schedule" || opts.EventName == "workflow_dispatch" {
		return targets
	}

	// If any workflow definition or GitHub action changed, run all targets.
	for _, file := range opts.ChangedFiles {
		if strings.HasPrefix(file, ".github/workflows/") || strings.HasPrefix(file, ".github/actions/") {
			return targets
		}
	}

	var matched []Target
	for _, target := range targets {
		if targetIsChanged(target, targets, opts.ChangedFiles) {
			matched = append(matched, target)
		}
	}

	return matched
}

func targetIsChanged(target Target, allTargets []Target, changedFiles []string) bool {
	for _, file := range changedFiles {
		if target.IsSubmodule {
			if strings.HasPrefix(file, target.Dir+"/") || file == target.Dir {
				return true
			}
		} else {
			// Root tools target: matches root go.mod/go.sum or files under tools/ not belonging to a submodule
			if file == "go.mod" || file == "go.sum" {
				return true
			}
			if strings.HasPrefix(file, "tools/") {
				inSubmodule := false
				for _, other := range allTargets {
					if other.IsSubmodule && (strings.HasPrefix(file, other.Dir+"/") || file == other.Dir) {
						inSubmodule = true
						break
					}
				}
				if !inSubmodule {
					return true
				}
			}
		}
	}
	return false
}
