package memory

import (
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
	cldf_ton_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton/provider"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func getTestTonChainSelectors() []uint64 {
	return []uint64{chainsel.TON_LOCALNET.Selector}
}

func generateChainsTon(t *testing.T, numChains int) []cldf_chain.BlockChain {
	testTonChainSelectors := getTestTonChainSelectors()
	if numChains > 1 {
		t.Fatalf("only one ton chain is supported for now, got %d", numChains)
	}
	if len(testTonChainSelectors) < numChains {
		t.Fatalf("not enough test ton chain selectors available")
	}

	chains := make([]cldf_chain.BlockChain, 0, numChains)
	for i := 0; i < numChains; i++ {
		selector := testTonChainSelectors[i]

		c, err := cldf_ton_provider.NewCTFChainProvider(t, selector,
			cldf_ton_provider.CTFChainProviderConfig{
				Once: once,
			},
		).Initialize(t.Context())
		require.NoError(t, err)

		chains = append(chains, c)
	}

	return chains
}

func createTonChainConfig(chainID string, chain cldf_ton.Chain) chainlink.RawConfig {
	chainConfig := chainlink.RawConfig{}

	chainConfig["Enabled"] = true
	chainConfig["ChainID"] = chainID
	chainConfig["NetworkName"] = "ton-local"
	chainConfig["NetworkNameFull"] = "ton-local"
	chainConfig["Nodes"] = []any{
		map[string]any{
			"Name": "primary",
			"URL":  chain.URL,
		},
	}

	return chainConfig
}
