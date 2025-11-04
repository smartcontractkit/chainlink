package soltestutils

import (
	"io"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/utils/solutils"
)

// Program names
const (
	// Keystone Programs
	ProgKeystoneForwarder = "keystone_forwarder"

	// Data Feeds Programs
	ProgDataFeedsCache = "data_feeds_cache"

	// MCMS Programs
	ProgMCM              = "mcm"
	ProgTimelock         = "timelock"
	ProgAccessController = "access_controller"
)

// Program names grouped by their usage.
var (
	MCMSProgramNames      = []string{ProgMCM, ProgTimelock, ProgAccessController}
	KeystoneProgramNames  = []string{ProgKeystoneForwarder}
	DataFeedsProgramNames = []string{ProgDataFeedsCache}
)

var (
	// onceCCIP is used to ensure that the program artifacts from the chainlink-ccip repository are only downloaded once.
	onceCCIP = &sync.Once{}
	// onceSolana is used to ensure that the program artifacts from the chainlink-solana repository are only downloaded once.
	onceSolana = &sync.Once{} //nolint:unused // Will be used once all tests are migrated to use this package
)

// Repositories that contain the program artifacts.
const (
	repoCCIP   = "chainlink-ccip"
	repoSolana = "chainlink-solana"
)

// ProgramEntry contains the information about a program.
type ProgramEntry struct {
	// ID is the program ID of the program.
	ID string

	// Repo is the repository name of where the program is located.
	Repo string
}

// Programs maps the program name to the program information.
//
// This is the source of truth for the program IDs and repositories.
var directory = programDirectory{
	// MCMS Programs
	ProgMCM:              {ID: "5vNJx78mz7KVMjhuipyr9jKBKcMrKYGdjGkgE4LUmjKk", Repo: repoCCIP},
	ProgTimelock:         {ID: "DoajfR5tK24xVw51fWcawUZWhAXD8yrBJVacc13neVQA", Repo: repoCCIP},
	ProgAccessController: {ID: "6KsN58MTnRQ8FfPaXHiFPPFGDRioikj9CdPvPxZJdCjb", Repo: repoCCIP},

	// Keystone Programs
	ProgKeystoneForwarder: {ID: "whV7Q5pi17hPPyaPksToDw1nMx6Lh8qmNWKFaLRQ4wz", Repo: repoSolana},

	// Data Feeds Programs
	ProgDataFeedsCache: {ID: "3kX63udXtYcsdj2737Wi2KGd2PhqiKPgAFAxstrjtRUa", Repo: repoSolana},
}

// downloadFunc is a function type for downloading program artifacts
type downloadFunc func(t *testing.T) string

// programDirectory maps the program name to the program information.
type programDirectory map[string]ProgramEntry

// LoadMCMSArtifacts loads the MCMS program artifacts into the temporary test directory.
//
// Returns the map of program names to IDs.
func (d programDirectory) LoadMCMSArtifacts(t *testing.T, dirPath string) map[string]string {
	return d.loadProgramArtifacts(t, MCMSProgramNames, downloadChainlinkCCIPProgramArtifacts, dirPath)
}

// LoadKeystonePrograms loads the Keystone and MCMS program artifacts into the temporary test directory.
//
// Returns the map of program names to IDs.
func (d programDirectory) LoadKeystonePrograms(t *testing.T, dirPath string) map[string]string {
	keystoneProgIDs := d.loadProgramArtifacts(t, KeystoneProgramNames, downloadChainlinkSolanaProgramArtifacts, dirPath)
	mcmsProgIDs := d.loadProgramArtifacts(t, MCMSProgramNames, downloadChainlinkCCIPProgramArtifacts, dirPath)

	progIDs := make(map[string]string, len(keystoneProgIDs)+len(mcmsProgIDs))
	maps.Copy(progIDs, keystoneProgIDs)
	maps.Copy(progIDs, mcmsProgIDs)

	return progIDs
}

// LoadDataFeedsPrograms loads the Data Feeds and MCMS program artifacts into the temporary test directory.
//
// Returns the map of program names to IDs.
func (d programDirectory) LoadDataFeedsPrograms(t *testing.T, dirPath string) map[string]string {
	dataFeedsProgIDs := d.loadProgramArtifacts(t, DataFeedsProgramNames, downloadChainlinkSolanaProgramArtifacts, dirPath)
	mcmsProgIDs := d.loadProgramArtifacts(t, MCMSProgramNames, downloadChainlinkCCIPProgramArtifacts, dirPath)

	progIDs := make(map[string]string, len(dataFeedsProgIDs)+len(mcmsProgIDs))
	maps.Copy(progIDs, dataFeedsProgIDs)
	maps.Copy(progIDs, mcmsProgIDs)

	return progIDs
}

// loadProgramArtifacts is a helper function that loads program artifacts into a temporary test directory.
// It downloads artifacts using the provided download function and copies the specified programs.
//
// Returns the map of program names to IDs.
func (d programDirectory) loadProgramArtifacts(t *testing.T, programNames []string, downloadFn downloadFunc, targetDir string) map[string]string {
	t.Helper()

	// Download the program artifacts using the provided download function
	cachePath := downloadFn(t)

	// Get the program IDs for the specified programs
	progIDs := d.ProgramIDs(programNames)

	// Copy the specific artifacts to the target directory
	for name := range progIDs {
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

	// Return the path to the cached artifacts and the map of program IDs
	return progIDs
}

// ProgramIDs returns the map of program names to IDs.
func (d programDirectory) ProgramIDs(names []string) map[string]string {
	ids := make(map[string]string, len(names))

	for _, name := range names {
		ids[name] = d[name].ID
	}

	return ids
}

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
