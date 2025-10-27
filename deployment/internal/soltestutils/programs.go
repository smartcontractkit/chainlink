package soltestutils

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/utils/solutils"
)

var (
	// onceCCIP is used to ensure that the program artifacts from the chainlink-ccip repository are only downloaded once.
	onceCCIP = &sync.Once{}
	// onceSolana is used to ensure that the program artifacts from the chainlink-solana repository are only downloaded once.
	onceSolana = &sync.Once{} //nolint:unused // Will be used once all tests are migrated to use this package
)

// MCMSProgramIDs is a map of predeployed MCMS Solana program IDs used in tests.
var MCMSProgramIDs = map[string]string{
	"mcm":               "5vNJx78mz7KVMjhuipyr9jKBKcMrKYGdjGkgE4LUmjKk",
	"timelock":          "DoajfR5tK24xVw51fWcawUZWhAXD8yrBJVacc13neVQA",
	"access_controller": "6KsN58MTnRQ8FfPaXHiFPPFGDRioikj9CdPvPxZJdCjb",
}

// MCMSPrograms downloads the MCMS program artifacts and returns the path to the cached artifacts
// and the map of program IDs to paths.
//
// This can be used to preload the MCMS program artifacts into a test environment as arguments to
// the WithSolanaContainer function.
//
// TODO: Remove the dependency on the memory package by extracting the download logic into a
// separate solutils package.
func ProgramsForMCMS(t *testing.T) (string, map[string]string) {
	t.Helper()

	targetDir := t.TempDir()

	// Download the MCMS program artifacts
	cachePath := downloadChainlinkCCIPProgramArtifacts(t)

	// Copy the specific artifacts to the path provided
	for name := range MCMSProgramIDs {
		src := filepath.Join(cachePath, name+".so")
		dst := filepath.Join(targetDir, name+".so")

		// Copy the cached artifacts to the target directory
		srcFile, err := os.Open(src)
		require.NoError(t, err)

		dstFile, err := os.Create(dst)
		require.NoError(t, err)

		_, err = io.Copy(dstFile, srcFile)
		require.NoError(t, err)

		srcFile.Close()
		dstFile.Close()
	}

	// Return the path to the cached artifacts and the map of program IDs to paths
	return targetDir, MCMSProgramIDs
}

// downloadCLSolanaProgramArtifacts downloads the Chainlink Solana program artifacts.
//
// The artifacts that are downloaded contain both the CCIP and MCMS program artifacts (even though
// this is called "CCIP" program artifacts).
func downloadCLSolanaProgramArtifacts(t *testing.T) string { //nolint:unused // Will be used once all tests are migrated to use this package
	t.Helper()

	cachePath := programsCachePath()

	onceSolana.Do(func() {
		err := solutils.DownloadChainlinkSolanaProgramArtifacts(t.Context(), cachePath, "", nil)
		require.NoError(t, err)
	})

	return cachePath
}

// downloadChainlinkCCIPProgramArtifacts downloads the Chainlink CCIP program artifacts for the
// test environment.
//
// The artifacts that are downloaded contain both the CCIP and MCMS program artifacts (even though
// this is called "CCIP" program artifacts).
func downloadChainlinkCCIPProgramArtifacts(t *testing.T) string {
	t.Helper()

	cachePath := programsCachePath()

	onceCCIP.Do(func() {
		err := solutils.DownloadChainlinkCCIPProgramArtifacts(t.Context(), cachePath, "", nil)
		require.NoError(t, err)
	})

	return cachePath
}

// programsCachePath returns the path to the cache directory for the program artifacts.
//
// This is used to cache the program artifacts so that they do not need to be downloaded every time
// the tests are run.
//
// The cache directory is located in the same directory as the current file.
func programsCachePath() string {
	// Get the directory of the current file
	_, currentFile, _, _ := runtime.Caller(0)

	dir := filepath.Dir(currentFile)

	return filepath.Join(dir, "programs_cache")
}
