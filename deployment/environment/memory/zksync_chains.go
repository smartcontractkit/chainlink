package memory

import (
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm/provider"
)

// zkSyncTestChainSelectors returns the selectors for the test zkSync chains. We arbitrarily
// start this from the EVM test selector TEST_90000051 and limit the number of chains you can load
// to 10. This avoid conflicts with other selectors.
var zkSyncTestChainSelectors = []uint64{
	chain_selectors.TEST_90000051.Selector,
	chain_selectors.TEST_90000052.Selector,
	chain_selectors.TEST_90000053.Selector,
	chain_selectors.TEST_90000054.Selector,
	chain_selectors.TEST_90000055.Selector,
	chain_selectors.TEST_90000056.Selector,
	chain_selectors.TEST_90000057.Selector,
	chain_selectors.TEST_90000058.Selector,
	chain_selectors.TEST_90000059.Selector,
	chain_selectors.TEST_90000060.Selector,
}

func GenerateChainsZk(t *testing.T, numChains int) []cldf_chain.BlockChain {
	if numChains > len(zkSyncTestChainSelectors) {
		t.Fatalf("not enough test zkSync chain selectors available, max is %d", len(zkSyncTestChainSelectors))
	}

	chains := make([]cldf_chain.BlockChain, 0, numChains)
	for i := range numChains {
		selector := zkSyncTestChainSelectors[i]

		c, err := cldf_evm_provider.NewZkSyncCTFChainProvider(t, selector,
			cldf_evm_provider.ZkSyncCTFChainProviderConfig{
				Once: once,
			},
		).Initialize(t.Context())
		require.NoError(t, err)

		chains = append(chains, c)
		require.NoError(t, err)
	}

	return chains
}
