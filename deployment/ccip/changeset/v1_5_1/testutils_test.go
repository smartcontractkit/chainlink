package v1_5_1_test

import (
	"bytes"
	"math/big"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/stretchr/testify/require"
)

// validateMemberOfTokenPoolPair performs checks required to validate that a token pool is fully configured for cross-chain transfer.
func validateMemberOfTokenPoolPair(
	t *testing.T,
	state stateview.CCIPOnChainState,
	tokenPool *token_pool.TokenPool,
	expectedRemotePools []common.Address,
	tokens map[uint64]*cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677],
	tokenSymbol shared.TokenSymbol,
	chainSelector uint64,
	rate *big.Int,
	capacity *big.Int,
	expectedOwner common.Address,
) {
	// Verify that the owner is expected
	owner, err := tokenPool.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, expectedOwner, owner)

	// Fetch the supported remote chains
	supportedChains, err := tokenPool.GetSupportedChains(nil)
	require.NoError(t, err)

	// Verify that the rate limits and remote addresses are correct
	for _, supportedChain := range supportedChains {
		inboundConfig, err := tokenPool.GetCurrentInboundRateLimiterState(nil, supportedChain)
		require.NoError(t, err)
		require.True(t, inboundConfig.IsEnabled)
		require.Equal(t, capacity, inboundConfig.Capacity)
		require.Equal(t, rate, inboundConfig.Rate)

		outboundConfig, err := tokenPool.GetCurrentOutboundRateLimiterState(nil, supportedChain)
		require.NoError(t, err)
		require.True(t, outboundConfig.IsEnabled)
		require.Equal(t, capacity, outboundConfig.Capacity)
		require.Equal(t, rate, outboundConfig.Rate)

		remoteTokenAddress, err := tokenPool.GetRemoteToken(nil, supportedChain)
		require.NoError(t, err)
		require.Equal(t, common.LeftPadBytes(tokens[supportedChain].Address.Bytes(), 32), remoteTokenAddress)

		remotePoolAddresses, err := tokenPool.GetRemotePools(nil, supportedChain)
		require.NoError(t, err)

		require.Len(t, remotePoolAddresses, len(expectedRemotePools))
		expectedRemotePoolAddressesBytes := make([][]byte, len(expectedRemotePools))
		for i, remotePool := range expectedRemotePools {
			expectedRemotePoolAddressesBytes[i] = common.LeftPadBytes(remotePool.Bytes(), 32)
		}
		sort.Slice(expectedRemotePoolAddressesBytes, func(i, j int) bool {
			return bytes.Compare(expectedRemotePoolAddressesBytes[i], expectedRemotePoolAddressesBytes[j]) < 0
		})
		sort.Slice(remotePoolAddresses, func(i, j int) bool {
			return bytes.Compare(remotePoolAddresses[i], remotePoolAddresses[j]) < 0
		})
		for i := range expectedRemotePoolAddressesBytes {
			require.Equal(t, expectedRemotePoolAddressesBytes[i], remotePoolAddresses[i])
		}
	}
}
