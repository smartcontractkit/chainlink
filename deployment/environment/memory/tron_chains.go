package memory

import (
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_tron_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron/provider"
)

func getTestTronChainSelectors() []uint64 {
	return []uint64{chainsel.TRON_TESTNET_NILE.Selector}
}

func generateChainsTron(t *testing.T, numChains int) []cldf_chain.BlockChain {
	testTronChainSelectors := getTestTronChainSelectors()
	if numChains > 1 {
		t.Fatalf("only one tron chain is supported for now, got %d", numChains)
	}
	if len(testTronChainSelectors) < numChains {
		t.Fatalf("not enough test tron chain selectors available")
	}

	chains := make([]cldf_chain.BlockChain, 0, numChains)
	for i := range numChains {
		selector := testTronChainSelectors[i]

		ctfDefault, err := cldf_tron_provider.SignerGenCTFDefault()
		require.NoError(t, err)
		c, err := cldf_tron_provider.NewCTFChainProvider(t, selector,
			cldf_tron_provider.CTFChainProviderConfig{
				Once:              once,
				DeployerSignerGen: ctfDefault,
			},
		).Initialize(t.Context())
		require.NoError(t, err)

		chains = append(chains, c)
	}

	return chains
}
