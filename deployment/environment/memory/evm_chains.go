package memory

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm/provider"
)

type EVMChain struct {
	Backend     *simulated.Backend
	DeployerKey *bind.TransactOpts
	Users       []*bind.TransactOpts
}

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

// funcAddress funds to an EVM address using a given transaction options.
func fundAddress(t *testing.T, from *bind.TransactOpts, to common.Address, amount *big.Int, backend *simulated.Backend) {
	ctx := t.Context()
	nonce, err := backend.Client().PendingNonceAt(ctx, from.From)
	require.NoError(t, err)
	gp, err := backend.Client().SuggestGasPrice(ctx)
	require.NoError(t, err)
	rawTx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gp,
		Gas:      21000,
		To:       &to,
		Value:    amount,
	})
	signedTx, err := from.Signer(from.From, rawTx)
	require.NoError(t, err)
	err = backend.Client().SendTransaction(ctx, signedTx)
	require.NoError(t, err)
	backend.Commit()
}
