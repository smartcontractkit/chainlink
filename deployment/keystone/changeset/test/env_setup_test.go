package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
)

func TestSetupEnv(t *testing.T) {
	t.Parallel()

	t.Run("test env with in memory nodes", func(t *testing.T) {
		for _, useMCMS := range []bool{true, false} {
			te := SetupContractTestEnv(t, EnvWrapperConfig{
				WFDonConfig:     DonConfig{Name: "wfDon", N: 4},
				AssetDonConfig:  DonConfig{Name: "assetDon", N: 4},
				WriterDonConfig: DonConfig{Name: "writerDon", N: 4},
				NumChains:       3,
				UseMCMS:         useMCMS,
			})
			t.Run(fmt.Sprintf("set up test env using MCMS: %t", useMCMS), func(t *testing.T) {
				require.NotNil(t, te.Env.ExistingAddresses)
				require.NotNil(t, te.Env.DataStore)
				addrs := te.Env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(te.RegistrySelector))
				wantAddrCnt := 4 // default setup include capabilities registry, workflow registry, forwarder and ocr3
				if useMCMS {
					wantAddrCnt += 5 // time lock, call proxy, canceller, bypass, proposer
				}
				require.Len(t, addrs, wantAddrCnt)
				require.Len(t, te.Env.BlockChains.EVMChains(), 3)
				require.NotEmpty(t, te.RegistrySelector)
				require.NotNil(t, te.Env.Offchain)
				// one forwarder on each chain
				forwardersByChain := te.OwnedForwarders()
				require.Len(t, forwardersByChain, 3)
				for _, forwarders := range forwardersByChain {
					require.Len(t, forwarders, 1)
					require.NotNil(t, forwarders[0])
					require.NotNil(t, forwarders[0].Contract)
				}
				r, err := te.Env.Offchain.ListNodes(t.Context(), &node.ListNodesRequest{})
				require.NoError(t, err)
				require.Len(t, r.Nodes, 12)
				for _, donNames := range []string{"wfDon", "assetDon", "writerDon"} {
					require.Len(t, te.GetP2PIDs(donNames), 4, "don %s should have 4 p2p ids", donNames)
				}
			})
		}
	})

	t.Run("test env with view only, non functional node stubs", func(t *testing.T) {
		for _, useMCMS := range []bool{true, false} {
			te := SetupContractTestEnv(t, EnvWrapperConfig{
				WFDonConfig:     DonConfig{Name: "wfDon", N: 4},
				AssetDonConfig:  DonConfig{Name: "assetDon", N: 4},
				WriterDonConfig: DonConfig{Name: "writerDon", N: 4},
				NumChains:       3,
				UseMCMS:         useMCMS,
			})
			t.Run(fmt.Sprintf("set up test env using MCMS: %t", useMCMS), func(t *testing.T) {
				require.NotNil(t, te.Env.ExistingAddresses)
				require.Len(t, te.Env.BlockChains.EVMChains(), 3)
				require.NotEmpty(t, te.RegistrySelector)
				require.NotNil(t, te.Env.Offchain)
				r, err := te.Env.Offchain.ListNodes(t.Context(), &node.ListNodesRequest{})
				require.NoError(t, err)
				require.Len(t, r.Nodes, 12)
				for _, donNames := range []string{"wfDon", "assetDon", "writerDon"} {
					require.Len(t, te.GetP2PIDs(donNames), 4, "don %s should have 4 p2p ids", donNames)
				}
			})
		}
	})
}
