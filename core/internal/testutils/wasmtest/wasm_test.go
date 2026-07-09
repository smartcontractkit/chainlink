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

func TestDecompressBinary(t *testing.T) {
	t.Parallel()

	want := []byte("wasm-bytes")
	var b bytes.Buffer
	bwr := brotli.NewWriter(&b)
	_, err := bwr.Write(want)
	require.NoError(t, err)
	require.NoError(t, bwr.Close())

	got, err := decompressBinary(b.Bytes())
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestCacheFilePath(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	path, err := cacheFilePath(repoRoot, "core/capabilities/compute/test/simple/cmd", "abc123")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(repoRoot, ".wasm-cache", "core_capabilities_compute_test_simple_cmd-abc123.wasm.br"), path)
}

func TestPkgSlugRejectsTraversal(t *testing.T) {
	t.Parallel()

	_, err := pkgSlug("../escape")
	require.Error(t, err)
}

func TestReadCacheFile(t *testing.T) {
	t.Parallel()

	cacheDir := t.TempDir()
	want := []byte("compressed-bytes")
	cachePath := filepath.Join(cacheDir, "fixture.wasm.br")
	require.NoError(t, os.WriteFile(cachePath, want, 0o600))

	got, err := readCacheFile(cachePath)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestGetOrBuildBinaryUsesCache(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(repoRoot, "go.mod"), []byte("module example.com/test\n"), 0o600))

	pkgRelPath := "core/capabilities/compute/test/simple/cmd"
	fingerprint := "deadbeefcafebabe"
	cacheDir := filepath.Join(repoRoot, cacheDirName)
	require.NoError(t, os.MkdirAll(cacheDir, 0o755))

	slug, err := pkgSlug(pkgRelPath)
	require.NoError(t, err)
	cachePath := filepath.Join(cacheDir, slug+"-"+fingerprint+".wasm.br")
	want := []byte("cached-wasm-bytes")
	var compressedBuf bytes.Buffer
	bwr := brotli.NewWriter(&compressedBuf)
	_, err = bwr.Write(want)
	require.NoError(t, err)
	require.NoError(t, bwr.Close())
	wantCompressed := compressedBuf.Bytes()
	require.NoError(t, os.WriteFile(cachePath, wantCompressed, 0o600))

	fingerprintCacheMu.Lock()
	fingerprintCache[pkgRelPath] = fingerprint
	fingerprintCacheMu.Unlock()

	gotCompressed, err := getOrBuildBinary(repoRoot, pkgRelPath, true)
	require.NoError(t, err)
	assert.Equal(t, wantCompressed, gotCompressed)

	gotDecompressed, err := getOrBuildBinary(repoRoot, pkgRelPath, false)
	require.NoError(t, err)
	assert.Equal(t, want, gotDecompressed)
}
