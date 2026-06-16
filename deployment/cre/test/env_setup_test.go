package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
)

func TestNewTestHarness(t *testing.T) {
	h := NewTestHarness(t)

	// all contracts on registry chain
	registryChainAddrs := h.Runtime.Environment().DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(h.RegistrySelector))
	require.Len(t, registryChainAddrs, 1) // registry
	require.Equal(t, datastore.ContractType("CapabilitiesRegistry"), registryChainAddrs[0].Type)

	for sel := range h.Runtime.Environment().BlockChains.EVMChains() {
		chainAddrs := h.Runtime.Environment().DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(sel))
		if sel != h.RegistrySelector {
			require.Empty(t, chainAddrs)
		} else {
			require.Len(t, chainAddrs, 1) // Only the registry should have addresses
		}
	}
}
