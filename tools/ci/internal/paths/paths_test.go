package paths

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // t.Chdir mutates process-wide state
func TestResolveFromRepoRoot_FallsBackToRepoRoot(t *testing.T) {
	repoRoot := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(repoRoot, "tools", "ci"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(repoRoot, "package.json"), []byte("{}"), 0o600))

	t.Chdir(filepath.Join(repoRoot, "tools", "ci"))

	got := ResolveFromRepoRoot("package.json")
	require.Equal(t, filepath.Join(repoRoot, "package.json"), got)
}

//nolint:paralleltest // t.Chdir mutates process-wide state
func TestResolveFromRepoRoot_PrefersLocalFile(t *testing.T) {
	repoRoot := t.TempDir()
	ciDir := filepath.Join(repoRoot, "tools", "ci")
	require.NoError(t, os.MkdirAll(ciDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(ciDir, "CHANGELOG.md"), []byte("local"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(repoRoot, "CHANGELOG.md"), []byte("root"), 0o600))

	t.Chdir(ciDir)

	got := ResolveFromRepoRoot("CHANGELOG.md")
	require.Equal(t, filepath.Join(ciDir, "CHANGELOG.md"), got)
}

//nolint:paralleltest // t.Chdir mutates process-wide state
func TestResolveFromRepoRoot_ReturnsInputWhenNothingMatches(t *testing.T) {
	repoRoot := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(repoRoot, "tools", "ci"), 0o755))
	t.Chdir(filepath.Join(repoRoot, "tools", "ci"))

	require.Equal(t, "does-not-exist.toml", ResolveFromRepoRoot("does-not-exist.toml"))
}

//nolint:paralleltest // t.Chdir mutates process-wide state
func TestResolveFromRepoRoot_PassesThroughEmptyAndAbsolute(t *testing.T) {
	repoRoot := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(repoRoot, "tools", "ci"), 0o755))
	t.Chdir(filepath.Join(repoRoot, "tools", "ci"))

	require.Empty(t, ResolveFromRepoRoot(""))

	abs := filepath.Join(repoRoot, "package.json")
	require.Equal(t, abs, ResolveFromRepoRoot(abs))
}
