package soltestutils

import (
	"io"
	"maps"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/utils/solutils"
)

// LoadKeystonePrograms loads the Keystone and MCMS program artifacts into the given directory.
//
// Returns a map of program names to IDs.
func LoadKeystonePrograms(t *testing.T, dir string) map[string]string {
	t.Helper()

	keystoneProgIDs := loadProgramArtifacts(t, solutils.KeystoneProgramNames, downloadChainlinkSolanaProgramArtifacts, dir)
	_, mcmsProgIDs := LoadMCMSPrograms(t, dir)

	progIDs := make(map[string]string, len(keystoneProgIDs)+len(mcmsProgIDs))
	maps.Copy(progIDs, keystoneProgIDs)
	maps.Copy(progIDs, mcmsProgIDs)

	return progIDs
}

// LoadMCMSPrograms loads the MCMS program artifacts into the given directory.
//
// Returns the path to the temporary test directory and a map of program names to IDs.
func LoadMCMSPrograms(t *testing.T, dir string) (string, map[string]string) {
	t.Helper()

	progIDs := loadProgramArtifacts(t,
		solutils.MCMSProgramNames, downloadChainlinkCCIPProgramArtifacts, dir,
	)

	return dir, progIDs
}

// LoadDataFeedsPrograms loads the Data Feeds and MCMS program artifacts into the given directory.
//
// Returns a map of program names to IDs.
func LoadDataFeedsPrograms(t *testing.T, dir string) map[string]string {
	t.Helper()

	dataFeedsProgIDs := loadProgramArtifacts(t, solutils.DataFeedsProgramNames, downloadChainlinkSolanaProgramArtifacts, dir)
	_, mcmsProgIDs := LoadMCMSPrograms(t, dir)

	progIDs := make(map[string]string, len(dataFeedsProgIDs)+len(mcmsProgIDs))
	maps.Copy(progIDs, dataFeedsProgIDs)
	maps.Copy(progIDs, mcmsProgIDs)

	return progIDs
}

// LoadCCIPPrograms loads the CCIP program artifacts into the given directory.
func LoadCCIPPrograms(t *testing.T, dir string) map[string]string {
	t.Helper()

	ccipProgIDs := loadProgramArtifacts(t, solutils.CCIPProgramNames, downloadChainlinkCCIPProgramArtifacts, dir)
	_, mcmsProgIDs := LoadMCMSPrograms(t, dir)

	progIDs := make(map[string]string, len(ccipProgIDs)+len(mcmsProgIDs))
	maps.Copy(progIDs, ccipProgIDs)
	maps.Copy(progIDs, mcmsProgIDs)

	return progIDs
}

// PreloadMCMS provides a convenience function to preload the MCMS program artifacts and address
// book for a given selector.
//
// TODO: Clean up this function to use the new LoadMCMSPrograms function.
func PreloadMCMS(t *testing.T, selector uint64) (string, map[string]string, *cldf.AddressBookMap) {
	t.Helper()

	dir := t.TempDir()

	_, programIDs := LoadMCMSPrograms(t, dir)

	ab := PreloadAddressBookWithMCMSPrograms(t, selector)

	return dir, programIDs, ab
}

// loadProgramArtifacts is a helper function that loads program artifacts into a temporary test directory.
// It downloads artifacts using the provided download function and copies the specified programs.
//
// Returns the map of program names to IDs.
func loadProgramArtifacts(t *testing.T, programNames []string, downloadFn downloadFunc, targetDir string) map[string]string {
	t.Helper()

	// Download the program artifacts using the provided download function
	cachePath := downloadFn(t)

	progIDs := make(map[string]string, len(programNames))

	// Copy the specific artifacts to the target directory and add the program ID to the map
	for _, name := range programNames {
		id := solutils.GetProgramID(name)
		require.NotEmpty(t, id, "program id not found for program name: %s", name)

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

		// Add the program ID to the map
		progIDs[name] = id
	}

	// Return the path to the cached artifacts and the map of program IDs
	return progIDs
}
