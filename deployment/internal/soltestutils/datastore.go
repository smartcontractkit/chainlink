package soltestutils

import (
	"testing"

	"github.com/stretchr/testify/require"

	cldfsolana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

// PreloadAddressBookWithMCMSPrograms creates and returns an address book containing preloaded MCMS
// Solana program addresses for the specified selector.
func PreloadAddressBookWithMCMSPrograms(t *testing.T, selector uint64) *cldf.AddressBookMap {
	t.Helper()

	ab := cldf.NewMemoryAddressBook()

	tv := cldf.NewTypeAndVersion(commontypes.ManyChainMultisigProgram, deployment.Version1_0_0)
	err := ab.Save(selector, MCMSProgramIDs["mcm"], tv)
	require.NoError(t, err)

	tv = cldf.NewTypeAndVersion(commontypes.AccessControllerProgram, deployment.Version1_0_0)
	err = ab.Save(selector, MCMSProgramIDs["access_controller"], tv)
	require.NoError(t, err)

	tv = cldf.NewTypeAndVersion(commontypes.RBACTimelockProgram, deployment.Version1_0_0)
	err = ab.Save(selector, MCMSProgramIDs["timelock"], tv)
	require.NoError(t, err)

	return ab
}

// GetMCMSStateFromAddressBook retrieves the state of the Solana MCMS contracts on the given chain.
func GetMCMSStateFromAddressBook(
	t *testing.T, ab cldf.AddressBook, chain cldfsolana.Chain,
) *state.MCMSWithTimelockStateSolana {
	t.Helper()

	addresses, err := ab.AddressesForChain(chain.Selector)
	require.NoError(t, err)

	mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	require.NoError(t, err)

	return mcmState
}
