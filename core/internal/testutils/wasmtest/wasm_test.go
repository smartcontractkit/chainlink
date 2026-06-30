package wasmtest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixtureFileName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		compress bool
		want     string
	}{
		{name: "uncompressed", compress: false, want: "output.wasm"},
		{name: "compressed", compress: true, want: "output.wasm.br"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, fixtureFileName(tc.compress))
		})
	}
}

func TestReadFixture(t *testing.T) {
	t.Parallel()

	t.Run("reads committed fixture", func(t *testing.T) {
		t.Parallel()
		pkgDir := t.TempDir()
		testdataDir := filepath.Join(pkgDir, "testdata")
		require.NoError(t, os.MkdirAll(testdataDir, 0o755))
		want := []byte("wasm-bytes")
		require.NoError(t, os.WriteFile(filepath.Join(testdataDir, "output.wasm"), want, 0o600))

		got, err := readFixture(pkgDir, false)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("reads compressed fixture", func(t *testing.T) {
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
