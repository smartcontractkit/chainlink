package cmd_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/cmd"
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

// initRepoWithOrigin builds a repo whose fake origin default branch is pinned
// at the initial commit, then returns the repo dir.
func initRepoWithOrigin(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	git(t, dir, "init")
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "core/config"), 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "core/config/foo.go"), []byte("package config\n"), 0o600))
	git(t, dir, "add", ".")
	git(t, dir, "commit", "-m", "base")
	baseSHA := git(t, dir, "rev-parse", "HEAD")
	git(t, dir, "update-ref", "refs/remotes/origin/develop", baseSHA)
	git(t, dir, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/develop")
	return dir
}

//nolint:paralleltest // t.Chdir is process-global; subtests cannot run in parallel
func TestGenerateCmdDiffBase(t *testing.T) {
	//nolint:paralleltest // t.Chdir is process-global
	t.Run("discovers files from earlier branch commits, not just staged", func(t *testing.T) {
		dir := initRepoWithOrigin(t)

		// First branch commit under core/config; nothing staged. CI would
		// regenerate config docs for this diff, so the hook must attempt it too
		// (and fail here only because the temp repo lacks the real generator
		// package).
		require.NoError(t, os.WriteFile(filepath.Join(dir, "core/config/foo.go"), []byte("package config // changed\n"), 0o600))
		git(t, dir, "commit", "-am", "branch work")

		// Second branch commit unrelated to core/config. A HEAD~1-based diff
		// now misses the config change entirely; the merge-base diff still sees it.
		require.NoError(t, os.WriteFile(filepath.Join(dir, "other.go"), []byte("package main\n"), 0o600))
		git(t, dir, "add", "other.go")
		git(t, dir, "commit", "-m", "unrelated work")

		t.Chdir(dir)
		root := cmd.NewRootCmd()
		root.SetArgs([]string{"generate"})

		err := root.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "generate config docs")
	})

	//nolint:paralleltest // t.Chdir is process-global
	t.Run("no branch changes means no generators run", func(t *testing.T) {
		dir := initRepoWithOrigin(t)

		t.Chdir(dir)
		root := cmd.NewRootCmd()
		root.SetArgs([]string{"generate"})

		require.NoError(t, root.Execute())
	})
}
