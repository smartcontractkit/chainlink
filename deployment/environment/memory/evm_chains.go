package memory

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm/provider"
)

// evmTestChainSelectors returns the selectors for the test EVM chains. We arbitrarily
// start this from the EVM test selector TEST_90000001 and limit the number of chains you can load
// to 10. This avoid conflicts with other selectors.
var evmTestChainSelectors = []uint64{
	chain_selectors.TEST_90000001.Selector,
	chain_selectors.TEST_90000002.Selector,
	chain_selectors.TEST_90000003.Selector,
	chain_selectors.TEST_90000004.Selector,
	chain_selectors.TEST_90000005.Selector,
	chain_selectors.TEST_90000006.Selector,
	chain_selectors.TEST_90000007.Selector,
	chain_selectors.TEST_90000008.Selector,
	chain_selectors.TEST_90000009.Selector,
	chain_selectors.TEST_90000010.Selector,
}

// GenerateChainsEVM generates a number of simulated EVM chains for testing purposes.
func generateChainsEVM(t *testing.T, numChains int, numUsers int) []cldf_chain.BlockChain {
	if numChains > len(evmTestChainSelectors) {
		require.Failf(t, "not enough test EVM chain selectors available", "max is %d",
			len(evmTestChainSelectors),
		)
	}

	chains := make([]cldf_chain.BlockChain, 0, numChains)
	for i := range numChains {
		selector := evmTestChainSelectors[i]

		c, err := cldf_evm_provider.NewSimChainProvider(t, selector,
			cldf_evm_provider.SimChainProviderConfig{
				NumAdditionalAccounts: uint(numUsers), //nolint:gosec // G115: This is for testing purposes only and should not overflow.
			},
		).Initialize(t.Context())
		require.NoError(t, err)

		chains = append(chains, c)
	}

	return chains
}

func generateChainsEVMWithIDs(t *testing.T, chainIDs []uint64, numUsers int) []cldf_chain.BlockChain {
	chains := make([]cldf_chain.BlockChain, 0, len(chainIDs))
	for _, cid := range chainIDs {
		// Determine the selector for the chain ID
		details, err := chain_selectors.GetChainDetailsByChainIDAndFamily(
			strconv.FormatUint(cid, 10), chain_selectors.FamilyEVM,
		)
		require.NoError(t, err, "selector is not found for chain id: %d", cid)

		c, err := cldf_evm_provider.NewSimChainProvider(t, details.ChainSelector,
			cldf_evm_provider.SimChainProviderConfig{
				NumAdditionalAccounts: uint(numUsers), //nolint:gosec // G115: This is for testing purposes only and should not overflow.
			},
		).Initialize(t.Context())
		require.NoError(t, err)

		chains = append(chains, c)
	}

	return chains
}
