package modules_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
)

// git runs a git command in dir and fails the test on error. Signing and user
// identity are pinned per-command so host gitconfig cannot leak in.
func git(t *testing.T, dir string, args ...string) string {
	t.Helper()

	base := make([]string, 0, 6+len(args))
	base = append(base, "-c", "commit.gpgsign=false", "-c", "user.email=test@example.com", "-c", "user.name=test")
	cmd := exec.CommandContext(t.Context(), "git", append(base, args...)...) //nolint:gosec // test helper: fixed binary, args are test-controlled
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "git %s: %s", strings.Join(args, " "), out)
	return strings.TrimSpace(string(out))
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
}

func TestFindAffectedModules(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	modDirs := []string{
		"", // root module
		"deployment",
		"integration-tests",
		"system-tests/lib",
		"core/scripts/cre/workflows/cron",
	}

	for _, mod := range modDirs {
		modPath := filepath.Join(tmpDir, mod)
		require.NoError(t, os.MkdirAll(modPath, 0o700))
		require.NoError(t, os.WriteFile(filepath.Join(modPath, "go.mod"), []byte("module test\n"), 0o600))
	}

	tests := []struct {
		name     string
		files    []string
		expected []modules.ModulePackages
	}{
		{
			name:  "root module package file",
			files: []string{"core/services/app.go"},
			expected: []modules.ModulePackages{
				{Module: ".", Packages: []string{"./core/services"}},
			},
		},
		{
			name:  "root module root-level file",
			files: []string{"main.go"},
			expected: []modules.ModulePackages{
				{Module: ".", Packages: []string{"."}},
			},
		},
		{
			name:  "submodule package file",
			files: []string{"deployment/environment/env.go"},
			expected: []modules.ModulePackages{
				{Module: "deployment", Packages: []string{"./environment"}},
			},
		},
		{
			name:  "submodule root-level file",
			files: []string{"deployment/main.go"},
			expected: []modules.ModulePackages{
				{Module: "deployment", Packages: []string{"."}},
			},
		},
		{
			name:  "nested submodule package file",
			files: []string{"core/scripts/cre/workflows/cron/pkg/main.go"},
			expected: []modules.ModulePackages{
				{Module: "core/scripts/cre/workflows/cron", Packages: []string{"./pkg"}},
			},
		},
		{
			name: "multiple files deduplicate packages within module",
			files: []string{
				"core/services/app.go",
				"core/services/db.go",
				"core/logger/log.go",
				"deployment/environment/env.go",
				"deployment/environment/node.go",
				"integration-tests/smoke/test.go",
			},
			expected: []modules.ModulePackages{
				{Module: ".", Packages: []string{"./core/logger", "./core/services"}},
				{Module: "deployment", Packages: []string{"./environment"}},
				{Module: "integration-tests", Packages: []string{"./smoke"}},
			},
		},
		{
			name:  "go.mod change triggers all packages in module",
			files: []string{"deployment/go.mod", "deployment/environment/env.go"},
			expected: []modules.ModulePackages{
				{Module: "deployment", Packages: []string{"./..."}},
			},
		},
		{
			name:  "go.sum change triggers all packages in module",
			files: []string{"go.sum"},
			expected: []modules.ModulePackages{
				{Module: ".", Packages: []string{"./..."}},
			},
		},
		{
			name:  "absolute path handling",
			files: []string{filepath.Join(tmpDir, "system-tests/lib/suite/test.go")},
			expected: []modules.ModulePackages{
				{Module: "system-tests/lib", Packages: []string{"./suite"}},
			},
		},
		{
			name:     "empty files slice",
			files:    []string{},
			expected: []modules.ModulePackages{},
		},
		{
			name:     "ignores non-go files",
			files:    []string{"README.md", "docs/architecture.png"},
			expected: []modules.ModulePackages{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mods, err := modules.FindAffectedModules(tmpDir, tc.files)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, mods)
		})
	}
}

func TestGetMergeBase(t *testing.T) {
	t.Parallel()

	t.Run("returns merge-base of HEAD and origin default branch", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		git(t, dir, "init")
		writeFile(t, dir, "a.go", "package a\n")
		git(t, dir, "add", ".")
		git(t, dir, "commit", "-m", "base")
		baseSHA := git(t, dir, "rev-parse", "HEAD")

		// Simulate a clone's remote default branch pinned at the base commit.
		git(t, dir, "update-ref", "refs/remotes/origin/develop", baseSHA)
		git(t, dir, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/develop")

		writeFile(t, dir, "a.go", "package a // changed\n")
		git(t, dir, "commit", "-am", "branch work")

		got := modules.GetMergeBase(t.Context(), dir)
		assert.Equal(t, baseSHA, got)
	})

	t.Run("falls back to HEAD when no origin default branch exists", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		git(t, dir, "init")
		writeFile(t, dir, "a.go", "package a\n")
		git(t, dir, "add", ".")
		git(t, dir, "commit", "-m", "base")

		got := modules.GetMergeBase(t.Context(), dir)
		assert.Equal(t, "HEAD", got)
	})

	t.Run("falls back to HEAD when no common ancestor exists", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		git(t, dir, "init")
		writeFile(t, dir, "a.go", "package a\n")
		git(t, dir, "add", ".")
		git(t, dir, "commit", "-m", "base")
		defaultBranch := git(t, dir, "branch", "--show-current")

		// Orphan history on the fake remote default branch: no shared ancestor.
		git(t, dir, "checkout", "--orphan", "other")
		writeFile(t, dir, "b.go", "package b\n")
		git(t, dir, "add", ".")
		git(t, dir, "commit", "-m", "unrelated")
		git(t, dir, "update-ref", "refs/remotes/origin/develop", "HEAD")
		git(t, dir, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/develop")
		git(t, dir, "checkout", defaultBranch)

		got := modules.GetMergeBase(t.Context(), dir)
		assert.Equal(t, "HEAD", got)
	})
}

func TestGetChangedFilesSince(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	git(t, dir, "init")
	writeFile(t, dir, "a.go", "package a\n")
	writeFile(t, dir, "keep.go", "package keep\n")
	writeFile(t, dir, "del.go", "package del\n")
	writeFile(t, dir, "index_only.go", "package index\n")
	git(t, dir, "add", ".")
	git(t, dir, "commit", "-m", "base")
	baseSHA := git(t, dir, "rev-parse", "HEAD")

	// Committed change and committed deletion since the base.
	writeFile(t, dir, "a.go", "package a // changed\n")
	git(t, dir, "rm", "-q", "del.go")
	git(t, dir, "commit", "-am", "branch work")

	// Staged addition and unstaged modification.
	writeFile(t, dir, "c.go", "package c\n")
	git(t, dir, "add", "c.go")
	writeFile(t, dir, "keep.go", "package keep // unstaged\n")

	// Index-only staged change: modify, stage, then revert working tree to base content.
	writeFile(t, dir, "index_only.go", "package index // staged\n")
	git(t, dir, "add", "index_only.go")
	writeFile(t, dir, "index_only.go", "package index\n")

	files, err := modules.GetChangedFilesSince(t.Context(), dir, baseSHA)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"a.go", "c.go", "index_only.go", "keep.go"}, files)
}
