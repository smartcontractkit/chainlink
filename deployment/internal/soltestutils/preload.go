package soltestutils

import (
	"testing"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// LoadKeystonePrograms loads the Keystone and MCMS program artifacts into the given directory.
//
// Returns a map of program names to IDs.
func LoadKeystonePrograms(t *testing.T, dir string) map[string]string {
	t.Helper()

	return directory.LoadKeystonePrograms(t, dir)
}

// LoadMCMSPrograms loads the MCMS program artifacts into the given directory.
//
// Returns the path to the temporary test directory and a map of program names to IDs.
func LoadMCMSPrograms(t *testing.T, dir string) (string, map[string]string) {
	t.Helper()

	return dir, directory.LoadMCMSArtifacts(t, dir)
}

// LoadDataFeedsPrograms loads the Data Feeds and MCMS program artifacts into the given directory.
//
// Returns a map of program names to IDs.
func LoadDataFeedsPrograms(t *testing.T, dir string) map[string]string {
	t.Helper()

	return directory.LoadDataFeedsPrograms(t, dir)
}

// PreloadMCMS provides a convenience function to preload the MCMS program artifacts and address
// book for a given selector.
//
// TODO: Clean up this function to use the new LoadMCMSPrograms function.
func PreloadMCMS(t *testing.T, selector uint64) (string, map[string]string, *cldf.AddressBookMap) {
	t.Helper()

	dir := t.TempDir()

	programIDs := directory.LoadMCMSArtifacts(t, dir)

	ab := PreloadAddressBookWithMCMSPrograms(t, selector)

	return dir, programIDs, ab
}
