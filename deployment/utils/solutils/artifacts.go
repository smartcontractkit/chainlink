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
//
// This internal function handles the HTTP download of a tar.gz archive and extracts all
// regular files to the target directory. It creates parent directories as needed and
// logs each extracted file if a logger is provided.
//
// The function performs the following steps:
//  1. Downloads the tar.gz archive from the provided URL
//  2. Decompresses the gzip stream
//  3. Extracts each regular file from the tar archive
//  4. Creates necessary parent directories
//  5. Writes files to the target directory using only the base filename
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - url: Full URL to the tar.gz release asset
//   - targetDir: Directory where extracted files will be stored
//   - lggr: Logger for progress information. Can be nil to disable logging
//
// Returns an error if the download fails, decompression fails, or file extraction fails.
func downloadProgramArtifacts(ctx context.Context, url string, targetDir string, lggr logger.Logger) error {
	// Download the artifact
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d - could not download tar.gz release artifact (url = '%s')", res.StatusCode, url)
	}

	// Extract the artifact to the target directory
	gzipReader, err := gzip.NewReader(res.Body)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	// Protection against decompression bombs
	const (
		maxFiles     = 1000              // Maximum number of files to extract
		maxTotalSize = 500 * 1024 * 1024 // Maximum total extraction size (500MB)
	)
	var (
		fileCount int
		totalSize int64
	)

	for {
		header, err := tarReader.Next()
		// End of tar archive
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		// Skip non-regular files
		if header.Typeflag != tar.TypeReg {
			continue
		}

		// Check limits to prevent decompression bombs
		fileCount++
		if fileCount > maxFiles {
			return fmt.Errorf("archive contains too many files (limit: %d)", maxFiles)
		}

		if totalSize+header.Size > maxTotalSize {
			return fmt.Errorf("archive total size exceeds limit (limit: %d bytes)", maxTotalSize)
		}

		// Copy the file to the target directory
		outPath := filepath.Join(targetDir, filepath.Base(header.Name))
		if err := os.MkdirAll(filepath.Dir(outPath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.Create(outPath)
		if err != nil {
			return err
		}

		// Limit individual file size to 100MB to prevent decompression bombs
		const maxFileSize = 100 * 1024 * 1024 // 100MB
		limitedReader := io.LimitReader(tarReader, maxFileSize)
		bytesWritten, err := io.Copy(outFile, limitedReader)
		if err != nil {
			outFile.Close()
			return err
		}

		// Update total size counter
		totalSize += bytesWritten

		if lggr != nil {
			lggr.Infof("Extracted Solana chainlink-solana artifact: %s", outPath)
		}

		outFile.Close()
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
