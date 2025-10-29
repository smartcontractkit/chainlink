package soltestutils

import (
	"testing"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// PreloadMCMS provides a convenience function to preload the MCMS program artifacts and address
// book for a given selector.
func PreloadMCMS(t *testing.T, selector uint64) (string, map[string]string, *cldf.AddressBookMap) {
	t.Helper()

	programsPath, programIDs := ProgramsForMCMS(t)
	ab := PreloadAddressBookWithMCMSPrograms(t, selector)

	return programsPath, programIDs, ab
}
