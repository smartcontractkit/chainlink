package paths

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindRepoRootFromSubdir(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module test\n\ngo 1.26.4\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".tool-versions"), []byte("golang 1.26.4\n"), 0o600))

	sub := filepath.Join(root, "integration-tests")
	require.NoError(t, os.Mkdir(sub, 0o755))
	require.NoError(t, os.Chdir(sub))
	t.Cleanup(func() { _ = os.Chdir(root) })

	got, err := findRepoRoot()
	require.NoError(t, err)

	wantRoot, err := filepath.EvalSymlinks(root)
	require.NoError(t, err)
	gotRoot, err := filepath.EvalSymlinks(got)
	require.NoError(t, err)
	require.Equal(t, wantRoot, gotRoot)
}
