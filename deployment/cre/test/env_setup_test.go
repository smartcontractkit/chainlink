package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
)

func TestSetupEnvV2(t *testing.T) {
	envV2 := SetupEnvV2(t, false)
	env := envV2.Env

	ds := env.DataStore

	// all contracts on registry chain
	registryChainAddrs := ds.Addresses().Filter(datastore.AddressRefByChainSelector(envV2.RegistrySelector))
	require.Len(t, registryChainAddrs, 1) // registry
	// only forwarder on non-home chain
	for sel := range env.BlockChains.EVMChains() {
		chainAddrs := ds.Addresses().Filter(datastore.AddressRefByChainSelector(sel))
		if sel != envV2.RegistrySelector {
			require.Len(t, chainAddrs, 0)
		} else {
			require.Len(t, chainAddrs, 1)
		}
	}
}
