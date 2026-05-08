package state_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

// TestAddressesForChain_preservesAddressBookLabelsWhenDatastoreRowHasEmptyLabels guards against datastore
// map entries overwriting labeled ManyChainMultiSig rows from the address book with empty labels (which breaks MCMS role binding).
func TestAddressesForChain_preservesAddressBookLabelsWhenDatastoreRowHasEmptyLabels(t *testing.T) {
	t.Parallel()

	sel := chain_selectors.TEST_90000001.Selector
	addr := common.HexToAddress("0x1111111111111111111111111111111111111111").Hex()

	ab := cldf.NewMemoryAddressBook()
	abTV := cldf.NewTypeAndVersion(commontypes.ManyChainMultisig, deployment.Version1_0_0)
	abTV.Labels.Add(commontypes.CancellerRole.String())
	require.NoError(t, ab.Save(sel, addr, abTV))

	store := datastore.NewMemoryDataStore()
	require.NoError(t, store.Addresses().Add(datastore.AddressRef{
		ChainSelector: sel,
		Address:       addr,
		Type:          datastore.ContractType(commontypes.ManyChainMultisig),
		Version:       &deployment.Version1_0_0,
		Labels:        datastore.NewLabelSet(),
	}))

	env := cldf.Environment{
		ExistingAddresses: ab,
		DataStore:         store.Seal(),
	}

	merged, err := state.AddressesForChain(env, sel, "")
	require.NoError(t, err)
	got := merged[addr]
	require.True(t, got.Labels.Contains(commontypes.CancellerRole.String()), "expected address-book CANCELLER label preserved")
}

// TestAddressesForChain_unionsLabelsWhenBothSourcesHaveLabels merges disjoint role labels on the same address/type/version.
func TestAddressesForChain_unionsLabelsWhenBothSourcesHaveLabels(t *testing.T) {
	t.Parallel()

	sel := chain_selectors.TEST_90000002.Selector
	addr := common.HexToAddress("0x2222222222222222222222222222222222222222").Hex()

	ab := cldf.NewMemoryAddressBook()
	abTV := cldf.NewTypeAndVersion(commontypes.ManyChainMultisig, deployment.Version1_0_0)
	abTV.Labels.Add(commontypes.ProposerRole.String())
	require.NoError(t, ab.Save(sel, addr, abTV))

	store := datastore.NewMemoryDataStore()
	require.NoError(t, store.Addresses().Add(datastore.AddressRef{
		ChainSelector: sel,
		Address:       addr,
		Type:          datastore.ContractType(commontypes.ManyChainMultisig),
		Version:       &deployment.Version1_0_0,
		Labels:        datastore.NewLabelSet(commontypes.BypasserRole.String()),
	}))

	env := cldf.Environment{
		ExistingAddresses: ab,
		DataStore:         store.Seal(),
	}

	merged, err := state.AddressesForChain(env, sel, "")
	require.NoError(t, err)
	got := merged[addr]
	require.True(t, got.Labels.Contains(commontypes.ProposerRole.String()))
	require.True(t, got.Labels.Contains(commontypes.BypasserRole.String()))
}
