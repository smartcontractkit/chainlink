package stateview

import (
	"testing"

	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestMergeAddressesFromBothSources(t *testing.T) {
	chainSelector := chain_selectors.ETHEREUM_MAINNET.Selector
	
	// Create a mock environment with both AddressBook and DataStore
	addressBook := cldf.NewMemoryAddressBook()
	err := addressBook.Save(chainSelector, "0x1234567890123456789012345678901234567890", 
		cldf.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0))
	require.NoError(t, err)
	
	dataStore := datastore.NewMemoryDataStore()
	err = dataStore.Addresses().Add(datastore.AddressRef{
		ChainSelector: chainSelector,
		Address:       "0xABCDEF1234567890123456789012345678901234",
		Type:          datastore.ContractType(types.RBACTimelock),
		Version:       &deployment.Version1_0_0,
	})
	require.NoError(t, err)
	
	env := cldf.Environment{
		ExistingAddresses: addressBook,
		DataStore:         dataStore.Seal(),
	}
	
	// Test the merge function
	mergedAddresses, err := mergeAddressesFromBothSources(env, chainSelector)
	require.NoError(t, err)
	
	// Should have addresses from both sources
	require.Len(t, mergedAddresses, 2)
	require.Contains(t, mergedAddresses, "0x1234567890123456789012345678901234567890")
	require.Contains(t, mergedAddresses, "0xABCDEF1234567890123456789012345678901234")
}

func TestMergeAddressesOnlyAddressBook(t *testing.T) {
	chainSelector := chain_selectors.ETHEREUM_MAINNET.Selector
	
	// Create environment with only AddressBook
	addressBook := cldf.NewMemoryAddressBook()
	err := addressBook.Save(chainSelector, "0x1234567890123456789012345678901234567890", 
		cldf.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0))
	require.NoError(t, err)
	
	env := cldf.Environment{
		ExistingAddresses: addressBook,
		DataStore:         nil, // No DataStore
	}
	
	// Test the merge function
	mergedAddresses, err := mergeAddressesFromBothSources(env, chainSelector)
	require.NoError(t, err)
	
	// Should have address from AddressBook only
	require.Len(t, mergedAddresses, 1)
	require.Contains(t, mergedAddresses, "0x1234567890123456789012345678901234567890")
}

func TestMergeAddressesOnlyDataStore(t *testing.T) {
	chainSelector := chain_selectors.ETHEREUM_MAINNET.Selector
	
	// Create environment with only DataStore
	dataStore := datastore.NewMemoryDataStore()
	err := dataStore.Addresses().Add(datastore.AddressRef{
		ChainSelector: chainSelector,
		Address:       "0xABCDEF1234567890123456789012345678901234",
		Type:          datastore.ContractType(types.RBACTimelock),
		Version:       &deployment.Version1_0_0,
	})
	require.NoError(t, err)
	
	addressBook := cldf.NewMemoryAddressBook()
	
	env := cldf.Environment{
		ExistingAddresses: addressBook,
		DataStore:         dataStore.Seal(),
	}
	
	// Test the merge function
	mergedAddresses, err := mergeAddressesFromBothSources(env, chainSelector)
	require.NoError(t, err)
	
	// Should have address from DataStore only
	require.Len(t, mergedAddresses, 1)
	require.Contains(t, mergedAddresses, "0xABCDEF1234567890123456789012345678901234")
}
