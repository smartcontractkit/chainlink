package changeset_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	creforwarder "github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
)

func TestDeployForwarder(t *testing.T) {
	t.Parallel()

	registrySel := chain_selectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{registrySel}),
	))
	require.NoError(t, err)

	t.Run("should deploy forwarder", func(t *testing.T) {
		err = rt.Exec(
			runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployForwarderV2), &changeset.DeployRequestV2{
				ChainSel:  registrySel,
				Qualifier: "my-test-forwarder",
			}),
		)
		require.NoError(t, err)

		// registry, ocr3, forwarder should be deployed on registry chain
		addrs := rt.State().DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(registrySel))
		require.Len(t, addrs, 1)
		fa := rt.State().DataStore.Addresses().Filter(datastore.AddressRefByQualifier("my-test-forwarder"))
		require.Len(t, fa, 1, "expected to find 'my-test-forwarder' qualifier")
		l := fa[0].Labels.List()
		require.Len(t, l, 2, "expected exactly 2 labels")
		require.Contains(t, l[0], creforwarder.DeploymentBlockLabel)
		require.Contains(t, l[1], creforwarder.DeploymentHashLabel)
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

	excludeChainsIfNeeded := func(excludeChains bool, env cldf.Environment) (uint64, map[uint64]struct{}) {
		if !excludeChains {
			return 0, nil
		}

		var chainToExclude uint64
		filteredChains := make(map[uint64]struct{})
		for chainID := range env.BlockChains.EVMChains() {
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
				require.Empty(t, csOut.MCMSTimelockProposals)
				// check that forwarder
				// TODO set up a listener to check that the forwarder is configured
				forwardersByChain := te.OwnedForwarders()
				for selector := range te.Env.BlockChains.EVMChains() {
					forwarders, ok := forwardersByChain[selector]
					require.True(t, ok)
					require.NotNil(t, forwarders)
					require.Len(t, forwarders, 1)
					f := forwarders[0]
					requireConfigUpdate(t, f.Contract, chainToExclude == selector)
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

				require.Len(t, csOut.MCMSTimelockProposals, expectedProposals)
				require.Nil(t, csOut.AddressBook)

				x := te.OwnedForwarders()
				for _, forwardersByChain := range x {
					require.Len(t, forwardersByChain, 1)
					f := forwardersByChain[0]
					require.NotNil(t, f.McmsContracts.Timelock)
					require.NotNil(t, f.McmsContracts.CallProxy)
				}
				_, err = commonchangeset.Apply(t, te.Env,
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(changeset.ConfigureForwardContracts),
						cfg,
					),
				)
				require.NoError(t, err)

				for selector, forwardersByChain := range te.OwnedForwarders() {
					require.Len(t, forwardersByChain, 1)
					f := forwardersByChain[0]
					requireConfigUpdate(t, f.Contract, chainToExclude == selector)
				}
			})
		}
	})
}
