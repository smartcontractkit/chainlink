package wasmtest

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

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

	pkgRelPath := "core/synthetic/test/cmd"
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
	t.Cleanup(func() {
		fingerprintCacheMu.Lock()
		delete(fingerprintCache, pkgRelPath)
		fingerprintCacheMu.Unlock()
	})

	gotCompressed, err := getOrBuildBinary(repoRoot, pkgRelPath, true)
	require.NoError(t, err)
	assert.Equal(t, wantCompressed, gotCompressed)

	gotDecompressed, err := getOrBuildBinary(repoRoot, pkgRelPath, false)
	require.NoError(t, err)
	assert.Equal(t, want, gotDecompressed)
}

// fingerprintFiles hashes a package's GoFiles+EmbedFiles the same way
// computeBuildFingerprint does, without invoking the go toolchain.
func fingerprintFiles(t *testing.T, repoRoot string, pkg listPackage) string {
	t.Helper()
	var digests []fileDigest
	require.NoError(t, hashPackageFiles(repoRoot, pkg, &digests))
	sort.Slice(digests, func(i, j int) bool { return digests[i].path < digests[j].path })
	h := sha256.New()
	require.NoError(t, writeFingerprint(h, "go1.26.4", sha256.Sum256([]byte("gosum")), digests))
	return hex.EncodeToString(h.Sum(nil))
}

// TestFingerprintInvalidation mutates shared files in sequence (base -> go edit ->
// embed edit), so the steps are intentionally ordered and not parallel subtests.
func TestFingerprintInvalidation(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	pkgDir := filepath.Join(repoRoot, "pkg")
	assetsDir := filepath.Join(pkgDir, "assets")
	require.NoError(t, os.MkdirAll(assetsDir, 0o755))

	goFile := filepath.Join(pkgDir, "main.go")
	embedFile := filepath.Join(assetsDir, "data.txt")
	require.NoError(t, os.WriteFile(goFile, []byte("package main\n"), 0o600))
	require.NoError(t, os.WriteFile(embedFile, []byte("v1"), 0o600))

	pkg := listPackage{
		Dir:        pkgDir,
		GoFiles:    []string{"main.go"},
		EmbedFiles: []string{"assets/data.txt"},
	}

	base := fingerprintFiles(t, repoRoot, pkg)

	require.NoError(t, os.WriteFile(goFile, []byte("package main\n\nvar X = 1\n"), 0o600))
	afterGoEdit := fingerprintFiles(t, repoRoot, pkg)
	assert.NotEqual(t, base, afterGoEdit, "go source edit must bust the fingerprint")

	require.NoError(t, os.WriteFile(embedFile, []byte("v2-changed"), 0o600))
	afterEmbedEdit := fingerprintFiles(t, repoRoot, pkg)
	assert.NotEqual(t, afterGoEdit, afterEmbedEdit, "embedded asset edit must bust the fingerprint")
}

func TestFingerprintStableWhenUnchanged(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	pkgDir := filepath.Join(repoRoot, "pkg")
	require.NoError(t, os.MkdirAll(pkgDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(pkgDir, "main.go"), []byte("package main\n"), 0o600))

	pkg := listPackage{Dir: pkgDir, GoFiles: []string{"main.go"}}

	first := fingerprintFiles(t, repoRoot, pkg)
	second := fingerprintFiles(t, repoRoot, pkg)
	assert.Equal(t, first, second)
}

func TestBuildLockPerKey(t *testing.T) {
	t.Parallel()

	lockA := buildLock("pkgA")
	assert.Same(t, lockA, buildLock("pkgA"), "same key must return the same lock")
	assert.NotSame(t, lockA, buildLock("pkgB"), "different keys must return different locks")
}

func TestBuildLockSerializesSameKey(t *testing.T) {
	t.Parallel()

	const goroutines = 25
	var active, maxActive int32
	var wg sync.WaitGroup

	for range goroutines {
		wg.Go(func() {
			mu := buildLock("serialize-test")
			mu.Lock()
			defer mu.Unlock()

			cur := atomic.AddInt32(&active, 1)
			for {
				prev := atomic.LoadInt32(&maxActive)
				if cur <= prev || atomic.CompareAndSwapInt32(&maxActive, prev, cur) {
					break
				}
			}
			time.Sleep(time.Millisecond)
			atomic.AddInt32(&active, -1)
		})
	}
	wg.Wait()

	assert.Equal(t, int32(1), maxActive, "same-key builds must not run concurrently")
}
