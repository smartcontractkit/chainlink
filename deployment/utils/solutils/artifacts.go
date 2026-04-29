package solutils

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	retry "github.com/avast/retry-go/v4"
	"golang.org/x/mod/modfile"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// DownloadChainlinkCCIPProgramArtifacts downloads CCIP program artifacts from the
// smartcontractkit/chainlink-ccip GitHub repository.
//
// The function downloads a tar.gz archive containing Solana program binaries and extracts
// them to the specified target directory. If sha is empty, it automatically resolves
// the version by parsing the "github.com/smartcontractkit/chainlink-ccip/chains/solana"
// dependency from the nearest go.mod file.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - targetDir: Directory where extracted artifacts will be stored
//   - sha: Git commit SHA or version identifier. If empty, auto-resolved from go.mod
//   - lggr: Logger for progress and debug information. Can be nil to disable logging
//
// Returns an error if the download fails, extraction fails, or SHA resolution fails.
func DownloadChainlinkCCIPProgramArtifacts(ctx context.Context, targetDir string, sha string, lggr logger.Logger) error {
	const (
		owner = "smartcontractkit"
		repo  = "chainlink-ccip"
		name  = "artifacts.tar.gz"
	)

	if sha == "" {
		version, err := getDependencySHA("github.com/smartcontractkit/chainlink-ccip/chains/solana")
		if err != nil {
			return err
		}
		sha = version
	}
	tag := "solana-artifacts-localtest-" + sha

	if lggr != nil {
		lggr.Infof("Downloading chainlink-ccip program artifacts (tag = %s)", tag)
	}

	return downloadProgramArtifacts(ctx, githubReleaseURL(owner, repo, tag, name), targetDir, lggr)
}

// DownloadChainlinkSolanaProgramArtifacts downloads Solana program artifacts from the
// smartcontractkit/chainlink-solana GitHub repository.
//
// The function downloads a tar.gz archive containing Solana program binaries and extracts
// them to the specified target directory. If sha is empty, a hardcoded default SHA
// "b0f7cd3fbdbb" is used for compatibility.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - targetDir: Directory where extracted artifacts will be stored
//   - sha: Git commit SHA or version identifier. If empty, uses default "b0f7cd3fbdbb"
//   - lggr: Logger for progress and debug information. Can be nil to disable logging
//
// Returns an error if the download fails or extraction fails.
func DownloadChainlinkSolanaProgramArtifacts(ctx context.Context, targetDir string, sha string, lggr logger.Logger) error {
	const (
		owner = "smartcontractkit"
		repo  = "chainlink-solana"
		name  = "artifacts.tar.gz"
	)

	if sha == "" {
		sha = "b0f7cd3fbdbb"
	}

	tag := "solana-artifacts-localtest-" + sha

	if lggr != nil {
		lggr.Infof("Downloading Solana chainlink-solana program artifacts (tag = %s)", tag)
	}

	return downloadProgramArtifacts(ctx, githubReleaseURL(owner, repo, tag, name), targetDir, lggr)
}

// downloadProgramArtifacts downloads and extracts program artifacts from a GitHub release URL.
// It retries up to 5 times with exponential backoff on transient errors (5xx, network failures).
// 4xx responses are returned immediately as unrecoverable. Extraction happens into a sibling
// temp directory and is atomically promoted via os.Rename only on full success.
func downloadProgramArtifacts(ctx context.Context, url string, targetDir string, lggr logger.Logger) error {
	const requestTimeout = 5 * time.Minute

	parentDir := filepath.Dir(targetDir)
	if err := os.MkdirAll(parentDir, 0o755); err != nil {
		return fmt.Errorf("creating parent dir %q: %w", parentDir, err)
	}

	client := &http.Client{Timeout: requestTimeout}

	// tmpDir is reset on each attempt so partial extraction never bleeds into the next.
	var tmpDir string

	err := retry.Do(
		func() error {
			// Clean up any partial extraction from the previous attempt.
			if tmpDir != "" {
				_ = os.RemoveAll(tmpDir)
			}
			var mkErr error
			tmpDir, mkErr = os.MkdirTemp(parentDir, ".artifacts-tmp-*")
			if mkErr != nil {
				return retry.Unrecoverable(fmt.Errorf("creating temp dir: %w", mkErr))
			}

			req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if reqErr != nil {
				return retry.Unrecoverable(reqErr)
			}

			res, doErr := client.Do(req)
			if doErr != nil {
				return doErr
			}
			defer func() {
				_, _ = io.Copy(io.Discard, res.Body)
				_ = res.Body.Close()
			}()

			if res.StatusCode >= 400 && res.StatusCode < 500 {
				// 4xx: artifact doesn't exist or auth failed — retrying won't help.
				return retry.Unrecoverable(fmt.Errorf("download failed with status %d - could not download tar.gz release artifact (url = %q)", res.StatusCode, url))
			}
			if res.StatusCode != http.StatusOK {
				// 5xx or other — retryable.
				return fmt.Errorf("download failed with status %d - could not download tar.gz release artifact (url = %q)", res.StatusCode, url)
			}

			return extractTarGz(res.Body, tmpDir, lggr)
		},
		retry.Attempts(5),
		retry.Delay(500*time.Millisecond),
		retry.DelayType(retry.BackOffDelay),
		retry.MaxDelay(8*time.Second),
		retry.Context(ctx),
		retry.OnRetry(func(n uint, err error) {
			if lggr != nil {
				lggr.Infof("Artifact download attempt %d failed (%v), retrying", n+1, err)
			}
		}),
	)

	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return err
	}

	// Atomically promote tmpDir to targetDir (temp dir is a sibling, same filesystem).
	if err := os.RemoveAll(targetDir); err != nil {
		_ = os.RemoveAll(tmpDir)
		return fmt.Errorf("removing existing target dir %q: %w", targetDir, err)
	}
	if err := os.Rename(tmpDir, targetDir); err != nil {
		_ = os.RemoveAll(tmpDir)
		return fmt.Errorf("promoting artifacts dir %q -> %q: %w", tmpDir, targetDir, err)
	}
	return nil
}

// extractTarGz decompresses a gzipped tar stream and writes regular files into outDir.
func extractTarGz(r io.Reader, outDir string, lggr logger.Logger) error {
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	// Protection against decompression bombs.
	const (
		maxFiles     = 1000
		maxTotalSize = 500 * 1024 * 1024 // 500MB
		maxFileSize  = 100 * 1024 * 1024 // 100MB per file
	)
	var (
		fileCount int
		totalSize int64
	)

	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		fileCount++
		if fileCount > maxFiles {
			return fmt.Errorf("archive contains too many files (limit: %d)", maxFiles)
		}
		if totalSize+header.Size > maxTotalSize {
			return fmt.Errorf("archive total size exceeds limit (limit: %d bytes)", maxTotalSize)
		}

		outPath := filepath.Join(outDir, filepath.Base(header.Name))
		if err := os.MkdirAll(filepath.Dir(outPath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.Create(outPath)
		if err != nil {
			return err
		}

		bytesWritten, err := io.Copy(outFile, io.LimitReader(tarReader, maxFileSize))
		outFile.Close()
		if err != nil {
			return err
		}
		totalSize += bytesWritten

		if lggr != nil {
			lggr.Infof("Extracted Solana artifact: %s", outPath)
		}
	}

	return nil
}

// githubReleaseURL constructs a GitHub release asset download URL.
//
// Builds a URL in the format: https://github.com/{owner}/{repo}/releases/download/{tag}/{name}
//
// Parameters:
//   - owner: GitHub repository owner (e.g., "smartcontractkit")
//   - repo: Repository name (e.g., "chainlink-ccip")
//   - tag: Release tag or version (e.g., "solana-artifacts-localtest-abc123")
//   - name: Asset filename (e.g., "artifacts.tar.gz")
//
// Returns the complete download URL for the GitHub release asset.
func githubReleaseURL(owner string, repo string, tag string, name string) string {
	return fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s", owner, repo, tag, name)
}

// getModFilePath locates the nearest go.mod file by traversing up the directory tree.
//
// Starting from the current source file's directory, this function walks up the
// filesystem hierarchy until it finds a go.mod file. This is useful for locating
// the project root and parsing dependency information.
//
// The search stops when either:
//   - A go.mod file is found (returns the full path)
//   - The filesystem root is reached (returns an error)
//
// Returns the absolute path to the go.mod file, or an error if none is found.
func getModFilePath() (string, error) {
	_, currentFile, _, _ := runtime.Caller(0)

	rootDir := filepath.Dir(currentFile)
	for {
		modPath := filepath.Join(rootDir, "go.mod")
		if _, err := os.Stat(modPath); err == nil {
			return modPath, nil
		}

		// Move up one directory
		parent := filepath.Dir(rootDir)

		// If we've reached the filesystem root, stop
		if parent == rootDir {
			return "", errors.New("go.mod file not found in any parent directory")
		}

		rootDir = parent
	}
}

// getDependencyVersion extracts the version of a specific dependency from a go.mod file.
//
// This function parses the go.mod file at the given path and searches for the specified
// dependency in the require section. It uses the golang.org/x/mod/modfile package for
// robust parsing that handles various go.mod formats.
//
// Parameters:
//   - modFilePath: Absolute path to the go.mod file to parse
//   - depPath: Full module path of the dependency (e.g., "github.com/user/repo")
//
// Returns the version string as specified in the go.mod file (e.g., "v1.2.3" or
// "v0.0.0-20230101000000-abc123def456"), or an error if the dependency is not found
// or the file cannot be parsed.
func getDependencyVersion(modFilePath, depPath string) (string, error) {
	gomod, err := os.ReadFile(modFilePath)
	if err != nil {
		return "", err
	}

	modFile, err := modfile.ParseLax("go.mod", gomod, nil)
	if err != nil {
		return "", err
	}

	for _, dep := range modFile.Require {
		if dep.Mod.Path == depPath {
			return dep.Mod.Version, nil
		}
	}

	return "", fmt.Errorf("dependency %s not found", depPath)
}

// getDependencySHA extracts the commit SHA from a dependency version in go.mod.
//
// This function combines go.mod file discovery and dependency version parsing to extract
// the commit SHA from pseudo-versions. It expects dependency versions in the format
// "v0.0.0-YYYYMMDDHHMMSS-{12-char-sha}" and returns the SHA portion.
//
// Parameters:
//   - depPath: Full module path of the dependency to look up
//
// Returns the 12-character commit SHA, or an error if the go.mod file cannot be found,
// the dependency is not present, or the version format is invalid.
func getDependencySHA(depPath string) (version string, err error) {
	modFilePath, err := getModFilePath()
	if err != nil {
		return "", err
	}

	ver, err := getDependencyVersion(modFilePath, depPath)
	if err != nil {
		return "", err
	}
	tokens := strings.Split(ver, "-")
	if len(tokens) == 3 {
		version := tokens[len(tokens)-1]
		return version, nil
	}

	return "", fmt.Errorf("invalid go.mod version: %s", ver)
}
