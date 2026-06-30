package wasmtest

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadFixture(t *testing.T) {
	t.Parallel()

	t.Run("reads and decompresses fixture", func(t *testing.T) {
		t.Parallel()
		pkgDir := t.TempDir()
		testdataDir := filepath.Join(pkgDir, "testdata")
		require.NoError(t, os.MkdirAll(testdataDir, 0o755))
		want := []byte("wasm-bytes")

		var b bytes.Buffer
		bwr := brotli.NewWriter(&b)
		_, err := bwr.Write(want)
		require.NoError(t, err)
		require.NoError(t, bwr.Close())

		require.NoError(t, os.WriteFile(filepath.Join(testdataDir, "output.wasm.br"), b.Bytes(), 0o600))

		got, err := readFixture(pkgDir, false)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("reads compressed fixture as is", func(t *testing.T) {
		t.Parallel()
		pkgDir := t.TempDir()
		testdataDir := filepath.Join(pkgDir, "testdata")
		require.NoError(t, os.MkdirAll(testdataDir, 0o755))
		want := []byte("compressed-bytes")
		require.NoError(t, os.WriteFile(filepath.Join(testdataDir, "output.wasm.br"), want, 0o600))

		got, err := readFixture(pkgDir, true)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("missing fixture errors with regenerate guidance", func(t *testing.T) {
		t.Parallel()
		_, err := readFixture(t.TempDir(), false)
		require.Error(t, err)
		assert.ErrorContains(t, err, "go generate")
	})
}
