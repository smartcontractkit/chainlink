package onchain

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
)

func TestDeployerAddress_DefaultAnvilKey(t *testing.T) {
	t.Parallel()

	addr, err := deployerAddress(blockchain.DefaultAnvilPrivateKey)
	require.NoError(t, err)
	require.Equal(t, common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"), addr)
}

func TestDeployerAddress_ExplicitKey(t *testing.T) {
	t.Parallel()

	addr, err := deployerAddress("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	require.NoError(t, err)
	require.Equal(t, common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"), addr)
}

func TestSyncAddressBook_CopiesAllAddressesFromDataStore(t *testing.T) {
	t.Parallel()

	memDs := datastore.NewMemoryDataStore()
	version := semver.MustParse("1.0.0")
	require.NoError(t, memDs.Addresses().Add(datastore.AddressRef{
		ChainSelector: 1337,
		Address:       "0xCapReg",
		Type:          "CapabilitiesRegistry",
		Version:       version,
	}))
	require.NoError(t, memDs.Addresses().Add(datastore.AddressRef{
		ChainSelector: 1337,
		Address:       "0xForwarder",
		Type:          "KeystoneForwarder",
		Version:       version,
		Qualifier:     "zone-c-workflow",
	}))

	d := &Deployer{}
	state := &domain.StateFile{}
	env := &cldf.Environment{DataStore: memDs.Seal()}

	require.NoError(t, d.syncAddressBook(env, state))

	require.Equal(t, "0xCapReg", state.GetAddress("CapabilitiesRegistry"))
	require.True(t, state.HasAddress("KeystoneForwarder"))
	require.Len(t, state.Addresses, 2)

	forwarder := state.Addresses[1]
	require.Equal(t, "0xForwarder", forwarder.Address)
	require.Equal(t, "zone-c-workflow", forwarder.Qualifier)
	require.Equal(t, domain.ChainSelector(1337), forwarder.ChainSelector)
}

func TestSyncAddressBook_ReSyncUpdatesExistingRef(t *testing.T) {
	t.Parallel()

	memDs := datastore.NewMemoryDataStore()
	version := semver.MustParse("1.0.0")
	require.NoError(t, memDs.Addresses().Add(datastore.AddressRef{
		ChainSelector: 1337,
		Address:       "0xCapReg",
		Type:          "CapabilitiesRegistry",
		Version:       version,
	}))

	d := &Deployer{}
	state := &domain.StateFile{}
	env := &cldf.Environment{DataStore: memDs.Seal()}
	require.NoError(t, d.syncAddressBook(env, state))
	require.NoError(t, d.syncAddressBook(env, state))

	require.Len(t, state.Addresses, 1, "syncing the same datastore twice must not duplicate entries")
}
