package tron_test

import (
	"fmt"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/tron"
)

func TestConfigureForwarder(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		nChains int
	}{
		{
			nChains: 1,
		},
	}

	t.Run("Should configure forwarder", func(t *testing.T) {
		for _, tcase := range testCases {
			nChains := tcase.nChains
			name := fmt.Sprintf("nChains=%d", nChains)

			t.Run(name, func(t *testing.T) {
				lggr := logger.Test(t)

				env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, memory.MemoryEnvironmentConfig{
					TronChains: nChains,
				})
				tronSel := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyTron))[0]

				// configure don for solana chain
				te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
					WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4, ChainSelectors: []uint64{tronSel}},
					AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
					WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
					NumChains:       nChains,
				})

				tronChain := env.BlockChains.TronChains()[tronSel]
				blockchains := make(map[uint64]cldf_chain.BlockChain)

				blockchains[tronSel] = tronChain

				for _, ch := range te.Env.BlockChains.All() {
					blockchains[ch.ChainSelector()] = ch
				}

				te.Env.BlockChains = cldf_chain.NewBlockChains(blockchains)

				deployOptions := cldf_tron.DefaultDeployOptions()
				deployOptions.FeeLimit = 1_000_000_000

				deployChangeset := commonchangeset.Configure(tron.DeployForwarder{},
					&tron.DeployForwarderRequest{
						ChainSelectors: []uint64{tronSel},
						Qualifier:      "my-test-forwarder",
						DeployOptions:  deployOptions,
					},
				)

				var wfNodes []string
				for _, id := range te.GetP2PIDs("wfDon") {
					wfNodes = append(wfNodes, id.String())
				}

				triggerOptions := cldf_tron.DefaultTriggerOptions()
				triggerOptions.FeeLimit = 1_000_000_000

				configureChangeset := commonchangeset.Configure(tron.ConfigureForwarder{},
					&tron.ConfigureForwarderRequest{
						WFDonName:        "test-wf-don",
						WFNodeIDs:        wfNodes,
						RegistryChainSel: te.RegistrySelector,
						TriggerOptions:   triggerOptions,
					},
				)

				env, _, err := commonchangeset.ApplyChangesets(t, te.Env, []commonchangeset.ConfiguredChangeSet{deployChangeset, configureChangeset})
				require.NoError(t, err)
			})
		}
	})
}
