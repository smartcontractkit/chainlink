package soltestutils

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"

	cldfsolana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

// RegisterMCMSPrograms registers the MCMS programs in the datastore for a given selector.
func RegisterMCMSPrograms(t *testing.T, selector uint64, ds *datastore.MemoryDataStore) {
	t.Helper()

	err := ds.AddressRefStore.Add(datastore.AddressRef{
		Address:       directory[ProgMCM].ID,
		ChainSelector: selector,
		Type:          datastore.ContractType(commontypes.ManyChainMultisigProgram),
		Version:       semver.MustParse("1.0.0"),
	})
	require.NoError(t, err)

	err = ds.AddressRefStore.Add(datastore.AddressRef{
		Address:       directory[ProgAccessController].ID,
		ChainSelector: selector,
		Type:          datastore.ContractType(commontypes.AccessControllerProgram),
		Version:       semver.MustParse("1.0.0"),
	})
	require.NoError(t, err)

	err = ds.AddressRefStore.Add(datastore.AddressRef{
		Address:       directory[ProgTimelock].ID,
		ChainSelector: selector,
		Type:          datastore.ContractType(commontypes.RBACTimelockProgram),
		Version:       semver.MustParse("1.0.0"),
	})
	require.NoError(t, err)
}

// PreloadAddressBookWithMCMSPrograms creates and returns an address book containing preloaded MCMS
// Solana program addresses for the specified selector.
func PreloadAddressBookWithMCMSPrograms(t *testing.T, selector uint64) *cldf.AddressBookMap {
	t.Helper()

	ab := cldf.NewMemoryAddressBook()

	tv := cldf.NewTypeAndVersion(commontypes.ManyChainMultisigProgram, deployment.Version1_0_0)
	err := ab.Save(selector, directory[ProgMCM].ID, tv)
	require.NoError(t, err)

	tv = cldf.NewTypeAndVersion(commontypes.AccessControllerProgram, deployment.Version1_0_0)
	err = ab.Save(selector, directory[ProgAccessController].ID, tv)
	require.NoError(t, err)

	tv = cldf.NewTypeAndVersion(commontypes.RBACTimelockProgram, deployment.Version1_0_0)
	err = ab.Save(selector, directory[ProgTimelock].ID, tv)
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
