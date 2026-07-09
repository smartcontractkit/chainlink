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
	"io/fs"
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
	buildFlags   = "-trimpath -buildvcs=false"
)

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
)

// GetTestBinary returns a WASM binary for the given package path relative to the repo root
// (e.g. core/capabilities/compute/test/simple/cmd). The binary is built lazily on first use
// and cached in .wasm-cache/ using a content-addressed filename.
func GetTestBinary(tb testing.TB, outputPath string, compress bool) []byte {
	tb.Helper()

	cacheKey := fmt.Sprintf("%s:%t", outputPath, compress)
	binaryCacheMu.RLock()
	cached, ok := binaryCache[cacheKey]
	binaryCacheMu.RUnlock()
	if ok {
		res := make([]byte, len(cached))
		copy(res, cached)
		return res
	}

	binaryCacheMu.Lock()
	defer binaryCacheMu.Unlock()
	if cached, ok = binaryCache[cacheKey]; ok {
		res := make([]byte, len(cached))
		copy(res, cached)
		return res
	}

	repoRoot, err := getRepoRoot()
	require.NoError(tb, err, "failed to get repo root: %s", err)

	binary, err := getOrBuildBinary(repoRoot, outputPath, compress)
	require.NoError(tb, err)

	cachedCopy := make([]byte, len(binary))
	copy(cachedCopy, binary)
	binaryCache[cacheKey] = cachedCopy

	return binary
}

func getOrBuildBinary(repoRoot, pkgRelPath string, compress bool) ([]byte, error) {
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

		compressed, err = buildAndCacheBinary(repoRoot, pkgRelPath, fingerprint, cachePath)
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

func buildAndCacheBinary(repoRoot, pkgRelPath, fingerprint, cachePath string) ([]byte, error) {
	pkgPath := modulePath + "/" + filepath.ToSlash(pkgRelPath)
	compressed, err := buildBinary(pkgPath)
	if err != nil {
		return nil, fmt.Errorf("build WASM for %s: %w", pkgRelPath, err)
	}

	cacheDir := filepath.Join(repoRoot, cacheDirName)
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return nil, fmt.Errorf("create WASM cache dir: %w", err)
	}

	tmpPath := cachePath + ".tmp." + fingerprint
	if err := os.WriteFile(tmpPath, compressed, 0o600); err != nil {
		return nil, fmt.Errorf("write WASM cache temp file: %w", err)
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
	Dir string `json:"Dir"`
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

	dirs, err := listDepDirs(pkgPath)
	if err != nil {
		return "", err
	}

	type fileDigest struct {
		path string
		hash string
	}
	var digests []fileDigest

	for _, dir := range dirs {
		if !strings.HasPrefix(dir, repoRoot+string(filepath.Separator)) && dir != repoRoot {
			continue
		}

		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				name := d.Name()
				if name == "testdata" || strings.HasPrefix(name, ".") {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(path, ".go") {
				return nil
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			sum := sha256.Sum256(data)
			rel, err := filepath.Rel(repoRoot, path)
			if err != nil {
				return err
			}
			digests = append(digests, fileDigest{
				path: filepath.ToSlash(rel),
				hash: hex.EncodeToString(sum[:]),
			})
			return nil
		})
		if err != nil {
			return "", fmt.Errorf("hash sources in %s: %w", dir, err)
		}
	}

	sort.Slice(digests, func(i, j int) bool {
		return digests[i].path < digests[j].path
	})

	h := sha256.New()
	_, _ = io.WriteString(h, goVersion)
	_, _ = io.WriteString(h, "\nwasip1\nwasm\n")
	_, _ = io.WriteString(h, buildFlags)
	_, _ = io.WriteString(h, "\n")
	_, _ = h.Write(goSumDigest[:])
	_, _ = io.WriteString(h, "\n")
	for _, d := range digests {
		_, _ = io.WriteString(h, d.path)
		_, _ = io.WriteString(h, "\n")
		_, _ = io.WriteString(h, d.hash)
		_, _ = io.WriteString(h, "\n")
	}

	return hex.EncodeToString(h.Sum(nil)[:16]), nil
}

func listDepDirs(pkgPath string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "list", "-deps", "-json", pkgPath) // #nosec
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("go list -deps %s: %w", pkgPath, err)
	}

	seen := make(map[string]struct{})
	var dirs []string
	dec := json.NewDecoder(bytes.NewReader(out))
	for {
		var pkg listPackage
		if err := dec.Decode(&pkg); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("decode go list output: %w", err)
		}
		if pkg.Dir == "" {
			continue
		}
		if _, ok := seen[pkg.Dir]; ok {
			continue
		}
		seen[pkg.Dir] = struct{}{}
		dirs = append(dirs, pkg.Dir)
	}

	sort.Strings(dirs)
	return dirs, nil
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
	cmd := exec.CommandContext(cmdCtx, "go", "build", "-trimpath", "-buildvcs=false", "-o", buildPath, pkgPath) // #nosec
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
