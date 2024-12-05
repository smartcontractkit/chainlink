package changeset_test

import (
	"fmt"
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

func TestDeployForwarder(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1, // nodes unused but required in config
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	registrySel := env.AllChainSelectors()[0]

	t.Run("should deploy forwarder", func(t *testing.T) {
		ab := deployment.NewMemoryAddressBook()

		// deploy forwarder
		env.ExistingAddresses = ab
		resp, err := changeset.DeployForwarder(env, registrySel)
		require.NoError(t, err)
		require.NotNil(t, resp)
		// registry, ocr3, forwarder should be deployed on registry chain
		addrs, err := resp.AddressBook.AddressesForChain(registrySel)
		require.NoError(t, err)
		require.Len(t, addrs, 1)

		// only forwarder on chain 1
		require.NotEqual(t, registrySel, env.AllChainSelectors()[1])
		oaddrs, err := resp.AddressBook.AddressesForChain(env.AllChainSelectors()[1])
		require.NoError(t, err)
		require.Len(t, oaddrs, 1)
	})
}

func TestConfigureForwarders(t *testing.T) {
	t.Parallel()

	t.Run("no mcms ", func(t *testing.T) {
		for _, nChains := range []int{1, 3} {
			name := fmt.Sprintf("nChains=%d", nChains)
			t.Run(name, func(t *testing.T) {
				te := SetupTestEnv(t, TestConfig{
					WFDonConfig:     DonConfig{N: 4},
					AssetDonConfig:  DonConfig{N: 4},
					WriterDonConfig: DonConfig{N: 4},
					NumChains:       nChains,
				})

				var wfNodes []string
				for id, _ := range te.WFNodes {
					wfNodes = append(wfNodes, id)
				}

				cfg := changeset.ConfigureForwardContractsRequest{
					WFDonName:        "test-wf-don",
					WFNodeIDs:        wfNodes,
					RegistryChainSel: te.RegistrySelector,
				}
				csOut, err := changeset.ConfigureForwardContracts(te.Env, cfg)
				require.NoError(t, err)
				require.Nil(t, csOut.AddressBook)
				require.Len(t, csOut.Proposals, 0)
				// check that forwarder
				// TODO set up a listener to check that the forwarder is configured
				contractSet := te.ContractSets()
				for selector := range te.Env.Chains {
					cs, ok := contractSet[selector]
					require.True(t, ok)
					require.NotNil(t, cs.Forwarder)
				}
			})
		}
	})

	t.Run("with mcms", func(t *testing.T) {
		for _, nChains := range []int{1, 3} {
			name := fmt.Sprintf("nChains=%d", nChains)
			t.Run(name, func(t *testing.T) {
				te := SetupTestEnv(t, TestConfig{
					WFDonConfig:     DonConfig{N: 4},
					AssetDonConfig:  DonConfig{N: 4},
					WriterDonConfig: DonConfig{N: 4},
					NumChains:       nChains,
					UseMCMS:         true,
				})

				var wfNodes []string
				for id, _ := range te.WFNodes {
					wfNodes = append(wfNodes, id)
				}

				cfg := changeset.ConfigureForwardContractsRequest{
					WFDonName:        "test-wf-don",
					WFNodeIDs:        wfNodes,
					RegistryChainSel: te.RegistrySelector,
					UseMCMS:          true,
				}
				csOut, err := changeset.ConfigureForwardContracts(te.Env, cfg)
				require.NoError(t, err)
				require.Len(t, csOut.Proposals, nChains)
				require.Nil(t, csOut.AddressBook)

				timelocks := make(map[uint64]*gethwrappers.RBACTimelock)
				for selector, contractSet := range te.ContractSets() {
					require.NotNil(t, contractSet.Timelock)
					timelocks[selector] = contractSet.Timelock
				}
				_, err = commonchangeset.ApplyChangesets(t, te.Env, timelocks, []commonchangeset.ChangesetApplication{
					{
						Changeset: commonchangeset.WrapChangeSet(changeset.ConfigureForwardContracts),
						Config:    cfg,
					},
				})
				require.NoError(t, err)

			})
		}
	})

}
