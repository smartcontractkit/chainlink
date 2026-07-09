package wasmtest

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/require"
)

const (
	modulePath   = "github.com/smartcontractkit/chainlink/v2"
	cacheDirName = ".wasm-cache"
)

var buildFlags = []string{"-trimpath", "-buildvcs=false"}

var getRepoRoot = sync.OnceValues(func() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("could not get caller info")
	}

	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("could not find repo root (no go.mod found)")
		}
		dir = parent
	}
})

var (
	binaryCache   = make(map[string][]byte)
	binaryCacheMu sync.RWMutex

	fingerprintCache   = make(map[string]string)
	fingerprintCacheMu sync.RWMutex

	buildLocksMu sync.Mutex
	buildLocks   = make(map[string]*sync.Mutex)
)

// GetTestBinary returns a WASM binary for the given package path relative to the repo root
// (e.g. core/capabilities/compute/test/simple/cmd). The binary is built lazily on first use
// and cached in .wasm-cache/ using a content-addressed filename.
func GetTestBinary(tb testing.TB, outputPath string, compress bool) []byte {
	tb.Helper()

	cacheKey := fmt.Sprintf("%s:%t", outputPath, compress)
	if cached, ok := loadBinaryCache(cacheKey); ok {
		return cached
	}

	// Serialize per package (not per cacheKey) so parallel tests requesting the
	// same binary build it once, while different packages still build in parallel.
	mu := buildLock(outputPath)
	mu.Lock()
	defer mu.Unlock()

	if cached, ok := loadBinaryCache(cacheKey); ok {
		return cached
	}

	repoRoot, err := getRepoRoot()
	require.NoError(tb, err, "failed to get repo root: %s", err)

	binary, err := getOrBuildBinary(repoRoot, outputPath, compress)
	require.NoError(tb, err)

	storeBinaryCache(cacheKey, binary)

	return binary
}

func loadBinaryCache(cacheKey string) ([]byte, bool) {
	binaryCacheMu.RLock()
	defer binaryCacheMu.RUnlock()
	cached, ok := binaryCache[cacheKey]
	if !ok {
		return nil, false
	}
	res := make([]byte, len(cached))
	copy(res, cached)
	return res, true
}

func storeBinaryCache(cacheKey string, binary []byte) {
	cp := make([]byte, len(binary))
	copy(cp, binary)
	binaryCacheMu.Lock()
	binaryCache[cacheKey] = cp
	binaryCacheMu.Unlock()
}

func buildLock(key string) *sync.Mutex {
	buildLocksMu.Lock()
	defer buildLocksMu.Unlock()
	mu, ok := buildLocks[key]
	if !ok {
		mu = &sync.Mutex{}
		buildLocks[key] = mu
	}
	return mu
}

func getOrBuildBinary(repoRoot, pkgRelPath string, compress bool) ([]byte, error) {
	if _, err := pkgSlug(pkgRelPath); err != nil {
		return nil, err
	}

	fingerprint, err := buildFingerprint(pkgRelPath)
	if err != nil {
		return nil, fmt.Errorf("compute build fingerprint for %s: %w", pkgRelPath, err)
	}

	cachePath, err := cacheFilePath(repoRoot, pkgRelPath, fingerprint)
	if err != nil {
		return nil, err
	}

	compressed, err := readCacheFile(cachePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("read WASM cache %s: %w", cachePath, err)
		}

		compressed, err = buildAndCacheBinary(repoRoot, pkgRelPath, cachePath)
		if err != nil {
			return nil, err
		}
	}

	if compress {
		res := make([]byte, len(compressed))
		copy(res, compressed)
		return res, nil
	}

	return decompressBinary(compressed)
}

func cacheFilePath(repoRoot, pkgRelPath, fingerprint string) (string, error) {
	slug, err := pkgSlug(pkgRelPath)
	if err != nil {
		return "", err
	}
	return filepath.Join(repoRoot, cacheDirName, fmt.Sprintf("%s-%s.wasm.br", slug, fingerprint)), nil
}

func pkgSlug(pkgRelPath string) (string, error) {
	clean := filepath.Clean(pkgRelPath)
	if clean == "." || strings.HasPrefix(clean, "..") {
		return "", fmt.Errorf("invalid package path: %s", pkgRelPath)
	}
	return strings.ReplaceAll(clean, string(filepath.Separator), "_"), nil
}

func readCacheFile(cachePath string) ([]byte, error) {
	return os.ReadFile(cachePath)
}

func decompressBinary(compressed []byte) ([]byte, error) {
	var b bytes.Buffer
	bwr := brotli.NewReader(bytes.NewReader(compressed))
	if _, err := io.Copy(&b, bwr); err != nil {
		return nil, fmt.Errorf("decompress WASM cache failed: %w", err)
	}
	return b.Bytes(), nil
}

func buildAndCacheBinary(repoRoot, pkgRelPath, cachePath string) ([]byte, error) {
	pkgPath := modulePath + "/" + filepath.ToSlash(pkgRelPath)
	compressed, err := buildBinary(pkgPath)
	if err != nil {
		return nil, fmt.Errorf("build WASM for %s: %w", pkgRelPath, err)
	}

	cacheDir := filepath.Join(repoRoot, cacheDirName)
	if err = os.MkdirAll(cacheDir, 0o755); err != nil {
		return nil, fmt.Errorf("create WASM cache dir: %w", err)
	}

	// Unique temp file per writer: concurrent test processes building the same
	// package must not share a temp path, or interleaved writes corrupt the cache.
	tmp, err := os.CreateTemp(cacheDir, filepath.Base(cachePath)+".tmp.*")
	if err != nil {
		return nil, fmt.Errorf("create WASM cache temp file: %w", err)
	}
	tmpPath := tmp.Name()
	if _, err := tmp.Write(compressed); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
		return nil, fmt.Errorf("write WASM cache temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return nil, fmt.Errorf("close WASM cache temp file: %w", err)
	}
	if err := os.Rename(tmpPath, cachePath); err != nil {
		_ = os.Remove(tmpPath)
		if existing, readErr := os.ReadFile(cachePath); readErr == nil {
			return existing, nil
		}
		return nil, fmt.Errorf("write WASM cache file: %w", err)
	}

	return compressed, nil
}

func buildFingerprint(pkgRelPath string) (string, error) {
	fingerprintCacheMu.RLock()
	if fp, ok := fingerprintCache[pkgRelPath]; ok {
		fingerprintCacheMu.RUnlock()
		return fp, nil
	}
	fingerprintCacheMu.RUnlock()

	pkgPath := modulePath + "/" + filepath.ToSlash(pkgRelPath)
	fp, err := computeBuildFingerprint(pkgPath)
	if err != nil {
		return "", err
	}

	fingerprintCacheMu.Lock()
	fingerprintCache[pkgRelPath] = fp
	fingerprintCacheMu.Unlock()

	return fp, nil
}

type listPackage struct {
	Dir        string   `json:"Dir"`
	GoFiles    []string `json:"GoFiles"`
	EmbedFiles []string `json:"EmbedFiles"`
}

type fileDigest struct {
	path string
	hash string
}

func computeBuildFingerprint(pkgPath string) (string, error) {
	repoRoot, err := getRepoRoot()
	if err != nil {
		return "", err
	}

	goVersion, err := goEnv("GOVERSION")
	if err != nil {
		return "", err
	}

	goSumPath := filepath.Join(repoRoot, "go.sum")
	goSum, err := os.ReadFile(goSumPath)
	if err != nil {
		return "", fmt.Errorf("read go.sum: %w", err)
	}
	goSumDigest := sha256.Sum256(goSum)

	pkgs, err := listDeps(pkgPath)
	if err != nil {
		return "", err
	}

	var digests []fileDigest

	for _, pkg := range pkgs {
		if pkg.Dir == "" {
			continue
		}
		if pkg.Dir != repoRoot && !strings.HasPrefix(pkg.Dir, repoRoot+string(filepath.Separator)) {
			continue
		}

		if err := hashPackageFiles(repoRoot, pkg, &digests); err != nil {
			return "", fmt.Errorf("hash sources in %s: %w", pkg.Dir, err)
		}
	}

	sort.Slice(digests, func(i, j int) bool {
		return digests[i].path < digests[j].path
	})

	h := sha256.New()
	if err := writeFingerprint(h, goVersion, goSumDigest, digests); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)[:16]), nil
}

// hashPackageFiles hashes exactly the files go compiles for this package under the
// target GOOS/GOARCH (GoFiles) plus its embedded assets (EmbedFiles). Using go list's
// file set (rather than walking the dir) omits build-tag-excluded and _test.go files
// that don't affect the WASM output, and correctly captures //go:embed inputs.
func hashPackageFiles(repoRoot string, pkg listPackage, digests *[]fileDigest) error {
	root, err := os.OpenRoot(pkg.Dir)
	if err != nil {
		return err
	}
	defer root.Close()

	names := make([]string, 0, len(pkg.GoFiles)+len(pkg.EmbedFiles))
	names = append(names, pkg.GoFiles...)
	names = append(names, pkg.EmbedFiles...)

	for _, name := range names {
		data, err := root.ReadFile(name)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(data)
		rel, err := filepath.Rel(repoRoot, filepath.Join(pkg.Dir, filepath.FromSlash(name)))
		if err != nil {
			return err
		}
		*digests = append(*digests, fileDigest{
			path: filepath.ToSlash(rel),
			hash: hex.EncodeToString(sum[:]),
		})
	}
	return nil
}

func writeFingerprint(h io.Writer, goVersion string, goSumDigest [sha256.Size]byte, digests []fileDigest) error {
	if _, err := io.WriteString(h, goVersion); err != nil {
		return fmt.Errorf("hash go version: %w", err)
	}
	if _, err := io.WriteString(h, "\nwasip1\nwasm\n"); err != nil {
		return fmt.Errorf("hash platform: %w", err)
	}
	if _, err := io.WriteString(h, strings.Join(buildFlags, " ")); err != nil {
		return fmt.Errorf("hash build flags: %w", err)
	}
	if _, err := io.WriteString(h, "\n"); err != nil {
		return fmt.Errorf("hash separator: %w", err)
	}
	if _, err := h.Write(goSumDigest[:]); err != nil {
		return fmt.Errorf("hash go.sum digest: %w", err)
	}
	if _, err := io.WriteString(h, "\n"); err != nil {
		return fmt.Errorf("hash separator: %w", err)
	}
	for _, d := range digests {
		if _, err := io.WriteString(h, d.path); err != nil {
			return fmt.Errorf("hash source path: %w", err)
		}
		if _, err := io.WriteString(h, "\n"); err != nil {
			return fmt.Errorf("hash separator: %w", err)
		}
		if _, err := io.WriteString(h, d.hash); err != nil {
			return fmt.Errorf("hash source digest: %w", err)
		}
		if _, err := io.WriteString(h, "\n"); err != nil {
			return fmt.Errorf("hash separator: %w", err)
		}
	}
	return nil
}

func listDeps(pkgPath string) ([]listPackage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "list", "-deps", "-json", pkgPath) // #nosec
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")
	// Capture stderr separately: stdout carries the JSON we decode, so it must stay
	// clean, but stderr holds the compiler/module message we want on failure.
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("go list -deps %s failed: %s: %w", pkgPath, strings.TrimSpace(stderr.String()), err)
	}

	var pkgs []listPackage
	dec := json.NewDecoder(bytes.NewReader(out))
	for {
		var pkg listPackage
		if err := dec.Decode(&pkg); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("decode go list output: %w", err)
		}
		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}

func goEnv(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "env", key) // #nosec
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("go env %s: %w", key, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func buildBinary(pkgPath string) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "wasmtest-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	cmdCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	buildPath := filepath.Join(tmpDir, "output.wasm")
	args := append([]string{"build"}, buildFlags...)
	args = append(args, "-o", buildPath, pkgPath)
	cmd := exec.CommandContext(cmdCtx, "go", args...) // #nosec
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("build failed: %s %w", string(output), err)
	}

	binary, err := os.ReadFile(buildPath)
	if err != nil {
		return nil, fmt.Errorf("read file failed: %w", err)
	}

	var b bytes.Buffer
	bwr := brotli.NewWriter(&b)
	if _, err = bwr.Write(binary); err != nil {
		return nil, err
	}
	if err = bwr.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
