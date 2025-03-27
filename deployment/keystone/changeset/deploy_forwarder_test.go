package changeset_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	forwarder "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder_1_0_0"
)

func TestDeployForwarder(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-111")
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

		chainSel := env.AllChainSelectors()[1]
		// only forwarder on chain 1
		require.NotEqual(t, registrySel, chainSel)
		oaddrs, err := resp.AddressBook.AddressesForChain(chainSel)
		require.NoError(t, err)
		require.Len(t, oaddrs, 1)
		for _, tv := range oaddrs {
			labelsList := tv.Labels.List()
			require.Len(t, labelsList, 2, "expected exactly 2 labels")
			require.Contains(t, labelsList[0], internal.DeploymentBlockLabel)
			require.Contains(t, labelsList[1], internal.DeploymentHashLabel)
		}
	})
}

func TestConfigureForwarders(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		nChains      int
		ExcludeChain bool // if true, configuration should be applied to all except one chain
	}{
		{
			nChains: 1,
		},
		{
			nChains:      3,
			ExcludeChain: true,
		},
	}

	excludeChainsIfNeeded := func(excludeChains bool, env deployment.Environment) (uint64, map[uint64]struct{}) {
		if !excludeChains {
			return 0, nil
		}

		var chainToExclude uint64
		filteredChains := make(map[uint64]struct{})
		for chainID := range env.Chains {
			// we do not really care which chain to exclude, so pick the first one
			if chainToExclude == 0 {
				chainToExclude = chainID
				continue
			}
			filteredChains[chainID] = struct{}{}
		}

		return chainToExclude, filteredChains
	}

	iteratorToSlice := func(iterator *forwarder.KeystoneForwarderConfigSetIterator) (result []*forwarder.KeystoneForwarderConfigSet, err error) {
		defer func(iterator *forwarder.KeystoneForwarderConfigSetIterator) {
			_ = iterator.Close()
		}(iterator)
		for iterator.Next() {
			result = append(result, iterator.Event)
			if iterator.Error() != nil {
				return nil, iterator.Error()
			}
		}
		return
	}

	requireConfigUpdate := func(t *testing.T, forwarder *forwarder.KeystoneForwarder, skippedConfigSet bool) {
		configsIterator, err := forwarder.FilterConfigSet(&bind.FilterOpts{}, nil, nil)
		require.NoError(t, err)
		configs, err := iteratorToSlice(configsIterator)
		require.NoError(t, err)
		if skippedConfigSet {
			require.Len(t, configs, 1) // once configuration is applied during initial setup
		} else {
			require.Len(t, configs, 2)
		}
	}

	t.Run("no mcms ", func(t *testing.T) {
		for _, testCase := range testCases {
			nChains := testCase.nChains
			name := fmt.Sprintf("nChains=%d", nChains)
			t.Run(name, func(t *testing.T) {
				te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
					WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
					AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
					WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
					NumChains:       nChains,
				})

				var wfNodes []string
				for _, id := range te.GetP2PIDs("wfDon") {
					wfNodes = append(wfNodes, id.String())
				}

				cfg := changeset.ConfigureForwardContractsRequest{
					WFDonName:        "test-wf-don",
					WFNodeIDs:        wfNodes,
					RegistryChainSel: te.RegistrySelector,
				}

				var chainToExclude uint64
				chainToExclude, cfg.Chains = excludeChainsIfNeeded(testCase.ExcludeChain, te.Env)

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
					requireConfigUpdate(t, cs.Forwarder, chainToExclude == selector)
				}
			})
		}
	})

	t.Run("with mcms", func(t *testing.T) {
		for _, testCase := range testCases {
			nChains := testCase.nChains
			name := fmt.Sprintf("nChains=%d", nChains)
			t.Run(name, func(t *testing.T) {
				te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
					WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
					AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
					WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
					NumChains:       nChains,
					UseMCMS:         true,
				})

				var wfNodes []string
				for _, id := range te.GetP2PIDs("wfDon") {
					wfNodes = append(wfNodes, id.String())
				}

				cfg := changeset.ConfigureForwardContractsRequest{
					WFDonName:        "test-wf-don",
					WFNodeIDs:        wfNodes,
					RegistryChainSel: te.RegistrySelector,
					MCMSConfig:       &changeset.MCMSConfig{MinDuration: 0},
				}

				var chainToExclude uint64
				chainToExclude, cfg.Chains = excludeChainsIfNeeded(testCase.ExcludeChain, te.Env)

				csOut, err := changeset.ConfigureForwardContracts(te.Env, cfg)
				require.NoError(t, err)
				expectedProposals := nChains
				if testCase.ExcludeChain {
					expectedProposals--
				}

				//nolint:staticcheck // migration will be done in a separate PR
				require.Len(t, csOut.Proposals, expectedProposals)
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
				_, err = commonchangeset.Apply(t, te.Env, timelockContracts,
					commonchangeset.Configure(
						deployment.CreateLegacyChangeSet(changeset.ConfigureForwardContracts),
						cfg,
					),
				)
				require.NoError(t, err)

				for selector, cs := range te.ContractSets() {
					requireConfigUpdate(t, cs.Forwarder, chainToExclude == selector)
				}
			})
		}
	})
}
