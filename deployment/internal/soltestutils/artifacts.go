package soltestutils

import (
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
	onceSolana = &sync.Once{}
)

// downloadFunc is a function type for downloading program artifacts
type downloadFunc func(t *testing.T) string

// downloadChainlinkSolanaProgramArtifacts downloads the Chainlink Solana program artifacts.
//
// The artifacts that are downloaded contain both the CCIP and MCMS program artifacts (even though
// this is called "CCIP" program artifacts).
func downloadChainlinkSolanaProgramArtifacts(t *testing.T) string {
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
