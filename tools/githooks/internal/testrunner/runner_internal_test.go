package testrunner

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNeedsBuild(t *testing.T) {
	t.Parallel()

	t.Run("binary missing requires build", func(t *testing.T) {
		t.Parallel()

		root := t.TempDir()
		binPath := filepath.Join(root, "tools", "test", ".bin", "test")
		srcDir := filepath.Join(root, "tools", "test")
		require.NoError(t, os.MkdirAll(srcDir, 0o750))

		needs, err := needsBuild(binPath, srcDir)
		require.NoError(t, err)
		assert.True(t, needs)
	})

	t.Run("binary newer than all sources skips build", func(t *testing.T) {
		t.Parallel()

		root := t.TempDir()
		binPath := filepath.Join(root, "tools", "test", ".bin", "test")
		srcDir := filepath.Join(root, "tools", "test")
		require.NoError(t, os.MkdirAll(srcDir, 0o750))
		require.NoError(t, os.WriteFile(filepath.Join(srcDir, "main.go"), []byte("package main\n"), 0o600))
		require.NoError(t, os.MkdirAll(filepath.Dir(binPath), 0o750))
		require.NoError(t, os.WriteFile(binPath, []byte("x"), 0o600))

		past := time.Now().Add(-time.Hour)
		require.NoError(t, os.Chtimes(filepath.Join(srcDir, "main.go"), past, past))

		needs, err := needsBuild(binPath, srcDir)
		require.NoError(t, err)
		assert.False(t, needs)
	})

	t.Run("source newer than binary requires build", func(t *testing.T) {
		t.Parallel()

		root := t.TempDir()
		binPath := filepath.Join(root, "tools", "test", ".bin", "test")
		srcDir := filepath.Join(root, "tools", "test")
		require.NoError(t, os.MkdirAll(srcDir, 0o750))
		srcFile := filepath.Join(srcDir, "main.go")
		require.NoError(t, os.WriteFile(srcFile, []byte("package main\n"), 0o600))
		require.NoError(t, os.MkdirAll(filepath.Dir(binPath), 0o750))
		require.NoError(t, os.WriteFile(binPath, []byte("x"), 0o600))

		past := time.Now().Add(-time.Hour)
		require.NoError(t, os.Chtimes(binPath, past, past))
		future := time.Now().Add(time.Hour)
		require.NoError(t, os.Chtimes(srcFile, future, future))

		needs, err := needsBuild(binPath, srcDir)
		require.NoError(t, err)
		assert.True(t, needs)
	})

	t.Run("missing source directory returns error", func(t *testing.T) {
		t.Parallel()

		root := t.TempDir()
		binPath := filepath.Join(root, "tools", "test", ".bin", "test")
		require.NoError(t, os.MkdirAll(filepath.Dir(binPath), 0o750))
		require.NoError(t, os.WriteFile(binPath, []byte("x"), 0o600))

		_, err := needsBuild(binPath, filepath.Join(root, "tools", "missing"))
		require.Error(t, err)
	})
}
