package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/cmd"
)

//nolint:paralleltest // t.Chdir is process-global; subtests cannot run in parallel
func TestTidyCmdDiffBase(t *testing.T) {
	//nolint:paralleltest // t.Chdir is process-global
	t.Run("tidies modules touched by earlier branch commits, not just staged", func(t *testing.T) {
		dir := initRepoWithOrigin(t)

		// Submodule with a broken go.mod: any tidy attempt against it fails,
		// which makes "did the command even try this module" observable.
		require.NoError(t, os.MkdirAll(filepath.Join(dir, "sub"), 0o700))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "sub/go.mod"), []byte("module\n"), 0o600))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "sub/b.go"), []byte("package sub\n"), 0o600))
		git(t, dir, "add", ".")
		git(t, dir, "commit", "-m", "branch work touching sub module")

		// Second commit unrelated to the submodule. A HEAD~1-based diff now
		// misses the submodule entirely; the merge-base diff still sees it.
		require.NoError(t, os.WriteFile(filepath.Join(dir, "other.go"), []byte("package main\n"), 0o600))
		git(t, dir, "add", "other.go")
		git(t, dir, "commit", "-m", "unrelated work")

		t.Chdir(dir)
		root := cmd.NewRootCmd()
		root.SetArgs([]string{"tidy"})

		err := root.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "module sub")
	})

	//nolint:paralleltest // t.Chdir is process-global
	t.Run("no branch changes means no modules tidied", func(t *testing.T) {
		dir := initRepoWithOrigin(t)

		t.Chdir(dir)
		root := cmd.NewRootCmd()
		root.SetArgs([]string{"tidy"})

		require.NoError(t, root.Execute())
	})
}
