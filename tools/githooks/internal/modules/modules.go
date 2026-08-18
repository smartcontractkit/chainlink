package modules

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

var relevantExtensions = map[string]bool{
	".go":  true,
	".mod": true,
	".sum": true,
}

// isExcludedFromUnitTests returns true for long-running E2E suites or example scripts.
func isExcludedFromUnitTests(relPath string) bool {
	clean := filepath.ToSlash(filepath.Clean(relPath))
	clean = strings.TrimPrefix(clean, "./")

	if strings.HasPrefix(clean, "system-tests") ||
		strings.HasPrefix(clean, "integration-tests") ||
		strings.HasPrefix(clean, "core/scripts/cre/environment/examples/workflows") {
		return true
	}
	return false
}

// ModulePackages associates a Go module with its changed package paths.
type ModulePackages struct {
	Module   string   // module root directory relative to repo root (e.g. ".", "deployment")
	Packages []string // package paths relative to module root (e.g. ["./core/logger", "./..."])
}

// FindAffectedModules maps a list of changed file paths to unique Go modules and their changed packages.
// Paths can be relative to repoRoot or absolute.
func FindAffectedModules(repoRoot string, files []string) ([]ModulePackages, error) {
	absRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute repo root: %w", err)
	}

	modulePkgMap := make(map[string]map[string]struct{})

	for _, file := range files {
		file = strings.TrimSpace(file)
		if file == "" {
			continue
		}

		ext := filepath.Ext(file)
		if !relevantExtensions[ext] {
			continue
		}

		absFile := file
		if !filepath.IsAbs(absFile) {
			absFile = filepath.Join(absRoot, file)
		}

		fileDir := filepath.Dir(absFile)
		baseName := filepath.Base(absFile)

		// Find nearest enclosing go.mod walking upwards
		dir := fileDir
		var modDir string

		for {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				modDir = dir
				break
			}

			if dir == absRoot || dir == filepath.Dir(dir) {
				break
			}
			dir = filepath.Dir(dir)
		}

		if modDir == "" {
			if _, err := os.Stat(filepath.Join(absRoot, "go.mod")); err == nil {
				modDir = absRoot
			}
		}

		if modDir != "" {
			relMod, err := filepath.Rel(absRoot, modDir)
			if err != nil {
				return nil, fmt.Errorf("failed to get relative module path: %w", err)
			}

			if _, ok := modulePkgMap[relMod]; !ok {
				modulePkgMap[relMod] = make(map[string]struct{})
			}

			// If go.mod or go.sum changed, the whole module needs analysis
			if baseName == "go.mod" || baseName == "go.sum" {
				modulePkgMap[relMod]["./..."] = struct{}{}
				continue
			}

			relPkg, err := filepath.Rel(modDir, fileDir)
			if err != nil {
				return nil, fmt.Errorf("failed to get relative package path: %w", err)
			}

			pkgPattern := "."
			if relPkg != "." {
				pkgPattern = "./" + filepath.ToSlash(relPkg)
			}

			modulePkgMap[relMod][pkgPattern] = struct{}{}
		}
	}

	result := make([]ModulePackages, 0, len(modulePkgMap))
	for mod, pkgSet := range modulePkgMap {
		var pkgs []string
		if _, hasAll := pkgSet["./..."]; hasAll {
			pkgs = []string{"./..."}
		} else {
			pkgs = make([]string, 0, len(pkgSet))
			for pkg := range pkgSet {
				pkgs = append(pkgs, pkg)
			}
			sort.Strings(pkgs)
		}

		result = append(result, ModulePackages{
			Module:   mod,
			Packages: pkgs,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Module < result[j].Module
	})

	return result, nil
}

// FindTestPackages maps changed files to unique Go test package patterns relative to repo root,
// skipping full E2E test suites (system-tests, integration-tests, workflow examples).
func FindTestPackages(repoRoot string, files []string) ([]string, error) {
	absRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute repo root: %w", err)
	}

	pkgSet := make(map[string]struct{})

	for _, file := range files {
		file = strings.TrimSpace(file)
		if file == "" {
			continue
		}

		ext := filepath.Ext(file)
		if !relevantExtensions[ext] {
			continue
		}

		absFile := file
		if !filepath.IsAbs(absFile) {
			absFile = filepath.Join(absRoot, file)
		}

		relFromRoot, err := filepath.Rel(absRoot, absFile)
		if err != nil {
			return nil, fmt.Errorf("failed to get relative path from repo root: %w", err)
		}

		if isExcludedFromUnitTests(relFromRoot) {
			continue
		}

		baseName := filepath.Base(absFile)
		fileDir := filepath.Dir(absFile)

		if baseName == "go.mod" || baseName == "go.sum" {
			relMod, relModErr := filepath.Rel(absRoot, fileDir)
			if relModErr != nil {
				return nil, fmt.Errorf("failed to get relative module path: %w", relModErr)
			}
			if relMod == "." {
				pkgSet["./core/..."] = struct{}{}
				pkgSet["./tools/..."] = struct{}{}
			} else {
				pkgSet["./"+filepath.ToSlash(relMod)+"/..."] = struct{}{}
			}
			continue
		}

		relPkg, relPkgErr := filepath.Rel(absRoot, fileDir)
		if relPkgErr != nil {
			return nil, fmt.Errorf("failed to get relative package path: %w", relPkgErr)
		}

		pkgPattern := "."
		if relPkg != "." {
			pkgPattern = "./" + filepath.ToSlash(relPkg)
		}
		pkgSet[pkgPattern] = struct{}{}
	}

	result := make([]string, 0, len(pkgSet))
	for pkg := range pkgSet {
		result = append(result, pkg)
	}
	sort.Strings(result)
	return result, nil
}

// GetStagedFiles fetches staged files from git.
func GetStagedFiles(ctx context.Context, repoRoot string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "git", "diff", "--cached", "--name-only", "--diff-filter=d")
	cmd.Dir = repoRoot

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git diff failed: %w (output: %s)", err, out.String())
	}

	lines := strings.Split(out.String(), "\n")
	var files []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			files = append(files, trimmed)
		}
	}

	return files, nil
}

// GetChangedFiles fetches changed files comparing against HEAD or push remote.
func GetChangedFiles(ctx context.Context, repoRoot string) ([]string, error) {
	// First check staged files
	staged, err := GetStagedFiles(ctx, repoRoot)
	if err == nil && len(staged) > 0 {
		return staged, nil
	}

	// Fallback to git diff HEAD~1..HEAD or working tree
	cmd := exec.CommandContext(ctx, "git", "diff", "HEAD~1..HEAD", "--name-only", "--diff-filter=d")
	cmd.Dir = repoRoot

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		// Fallback to uncommitted working tree changes
		cmd = exec.CommandContext(ctx, "git", "diff", "--name-only", "--diff-filter=d")
		cmd.Dir = repoRoot
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("git diff failed: %w (output: %s)", err, out.String())
		}
	}

	lines := strings.Split(out.String(), "\n")
	var files []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			files = append(files, trimmed)
		}
	}

	return files, nil
}
