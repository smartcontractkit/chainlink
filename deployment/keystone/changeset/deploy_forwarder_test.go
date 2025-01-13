package changeset_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
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
		resp, err := changeset.DeployForwarder(env, changeset.DeployForwarderRequest{})
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
				te := test.SetupTestEnv(t, test.TestConfig{
					WFDonConfig:     test.DonConfig{N: 4},
					AssetDonConfig:  test.DonConfig{N: 4},
					WriterDonConfig: test.DonConfig{N: 4},
					NumChains:       nChains,
				})

				var wfNodes []string
				for id := range te.WFNodes {
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
				require.Empty(t, csOut.Proposals)
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
				te := test.SetupTestEnv(t, test.TestConfig{
					WFDonConfig:     test.DonConfig{N: 4},
					AssetDonConfig:  test.DonConfig{N: 4},
					WriterDonConfig: test.DonConfig{N: 4},
					NumChains:       nChains,
					UseMCMS:         true,
				})

				var wfNodes []string
				for id := range te.WFNodes {
					wfNodes = append(wfNodes, id)
				}

				cfg := changeset.ConfigureForwardContractsRequest{
					WFDonName:        "test-wf-don",
					WFNodeIDs:        wfNodes,
					RegistryChainSel: te.RegistrySelector,
					MCMSConfig:       &changeset.MCMSConfig{MinDuration: 0},
				}
				csOut, err := changeset.ConfigureForwardContracts(te.Env, cfg)
				require.NoError(t, err)
				require.Len(t, csOut.Proposals, nChains)
				require.Nil(t, csOut.AddressBook)

				timelockContracts := make(map[uint64]*proposalutils.TimelockExecutionContracts)
				for selector, contractSet := range te.ContractSets() {
					require.NotNil(t, contractSet.Timelock)
					require.NotNil(t, contractSet.CallProxy)
					timelockContracts[selector] = &proposalutils.TimelockExecutionContracts{
						Timelock:  contractSet.Timelock,
						CallProxy: contractSet.CallProxy,
					}
				}
				_, err = commonchangeset.ApplyChangesets(t, te.Env, timelockContracts, []commonchangeset.ChangesetApplication{
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
