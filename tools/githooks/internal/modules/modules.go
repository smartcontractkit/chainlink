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

// isExcludedFromUnitTests returns true for long-running E2E suites (deployment, system-tests,
// integration-tests, example scripts, or non-unit tool directories).
func isExcludedFromUnitTests(relPath string) bool {
	clean := filepath.ToSlash(filepath.Clean(relPath))
	clean = strings.TrimPrefix(clean, "./")

	if strings.HasPrefix(clean, "deployment") ||
		strings.HasPrefix(clean, "system-tests") ||
		strings.HasPrefix(clean, "integration-tests") ||
		strings.HasPrefix(clean, "tools/secrets") ||
		strings.HasPrefix(clean, "tools/docker") ||
		strings.HasPrefix(clean, "tools/benchmark") ||
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

// FindModuleDir walks up from path until it locates a go.mod, bounded by repoRoot.
// Returns the absolute directory containing the nearest go.mod, or "" if none.
func FindModuleDir(repoRoot, path string) (string, error) {
	absRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute repo root: %w", err)
	}

	absPath := path
	if !filepath.IsAbs(absPath) {
		absPath = filepath.Join(absRoot, path)
	}

	dir := filepath.Dir(filepath.Clean(absPath))
	for {
		if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
			return dir, nil
		}

		if dir == absRoot || dir == filepath.Dir(dir) {
			break
		}
		dir = filepath.Dir(dir)
	}

	if _, statErr := os.Stat(filepath.Join(absRoot, "go.mod")); statErr == nil {
		return absRoot, nil
	}

	return "", nil
}

// ReadModulePath returns the module path declared in the go.mod at modDir.
func ReadModulePath(modDir string) (string, error) {
	goModPath := filepath.Join(modDir, "go.mod")
	raw, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", goModPath, err)
	}

	for line := range strings.SplitSeq(string(raw), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module")), nil
		}
	}

	return "", fmt.Errorf("no module path found in %s", goModPath)
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

		modDir, err := FindModuleDir(absRoot, absFile)
		if err != nil {
			return nil, err
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

// FindTestModules maps changed files to unique Go modules and their package paths,
// skipping excluded E2E suites (deployment, system-tests, integration-tests, workflow examples).
func FindTestModules(repoRoot string, files []string) ([]ModulePackages, error) {
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

		relFromRoot, relErr := filepath.Rel(absRoot, absFile)
		if relErr != nil {
			return nil, fmt.Errorf("failed to get relative path from repo root: %w", relErr)
		}

		if isExcludedFromUnitTests(relFromRoot) {
			continue
		}

		fileDir := filepath.Dir(absFile)
		baseName := filepath.Base(absFile)

		modDir, findErr := FindModuleDir(absRoot, absFile)
		if findErr != nil {
			return nil, findErr
		}

		if modDir != "" {
			relMod, relModErr := filepath.Rel(absRoot, modDir)
			if relModErr != nil {
				return nil, fmt.Errorf("failed to get relative module path: %w", relModErr)
			}

			if _, ok := modulePkgMap[relMod]; !ok {
				modulePkgMap[relMod] = make(map[string]struct{})
			}

			// If go.mod or go.sum changed, a dependency bump can affect any
			// package in the module, so run tests on the entire module.
			if baseName == "go.mod" || baseName == "go.sum" {
				modulePkgMap[relMod]["./..."] = struct{}{}
				continue
			}

			relPkg, relPkgErr := filepath.Rel(modDir, fileDir)
			if relPkgErr != nil {
				return nil, fmt.Errorf("failed to get relative package path: %w", relPkgErr)
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

// GetMergeBase returns the merge-base of HEAD with the origin default branch
// (refs/remotes/origin/HEAD). This matches the diff base CI uses for
// only-new-issues, so local lint results align with CI. Falls back to "HEAD"
// when the remote default branch is missing or shares no history.
func GetMergeBase(ctx context.Context, repoRoot string) string {
	refCmd := exec.CommandContext(ctx, "git", "symbolic-ref", "--short", "refs/remotes/origin/HEAD")
	refCmd.Dir = repoRoot
	refOut, err := refCmd.Output()
	if err != nil {
		return "HEAD"
	}
	defaultBranch := strings.TrimSpace(string(refOut))
	if defaultBranch == "" {
		return "HEAD"
	}

	mbCmd := exec.CommandContext(ctx, "git", "merge-base", "HEAD", defaultBranch)
	mbCmd.Dir = repoRoot
	mbOut, err := mbCmd.Output()
	if err != nil {
		return "HEAD"
	}
	if sha := strings.TrimSpace(string(mbOut)); sha != "" {
		return sha
	}
	return "HEAD"
}

// GetChangedFilesSince returns files changed relative to rev, including
// committed branch changes, staged changes, and unstaged working-tree changes
// (git diff <rev> semantics). Deleted files are excluded.
func GetChangedFilesSince(ctx context.Context, repoRoot, rev string) ([]string, error) {
	// Union the working-tree diff (captures committed + unstaged changes) with the
	// index diff (captures committed + staged changes), so index-only staged changes
	// aren't missed.
	diffArgs := [][]string{
		{"diff", "--name-only", "--diff-filter=d", rev, "--"},
		{"diff", "--cached", "--name-only", "--diff-filter=d", rev, "--"},
	}

	fileSet := make(map[string]struct{})
	for _, args := range diffArgs {
		cmd := exec.CommandContext(ctx, "git", args...)
		cmd.Dir = repoRoot
		out, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("git %s failed: %w (output: %s)", strings.Join(args, " "), err, string(out))
		}
		for line := range strings.SplitSeq(string(out), "\n") {
			if trimmed := strings.TrimSpace(line); trimmed != "" {
				fileSet[trimmed] = struct{}{}
			}
		}
	}

	files := make([]string, 0, len(fileSet))
	for f := range fileSet {
		files = append(files, f)
	}
	sort.Strings(files)
	return files, nil
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

	var (
		out    bytes.Buffer
		errOut bytes.Buffer
	)
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	if err := cmd.Run(); err != nil {
		// Fallback to uncommitted working tree changes. Reset buffers so error
		// output never mixes into the parsed file list.
		out.Reset()
		errOut.Reset()

		cmd = exec.CommandContext(ctx, "git", "diff", "--name-only", "--diff-filter=d")
		cmd.Dir = repoRoot
		cmd.Stdout = &out
		cmd.Stderr = &errOut
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("git diff failed: %w (output: %s)", err, errOut.String())
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
