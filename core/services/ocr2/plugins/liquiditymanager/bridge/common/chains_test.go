package common

import (
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func TestSupports(t *testing.T) {
	tests := []struct {
		src, dest models.NetworkSelector
		expected  bool
	}{
		{
			src:      models.NetworkSelector(chainsel.ETHEREUM_MAINNET.Selector),
			dest:     models.NetworkSelector(chainsel.ETHEREUM_MAINNET_ARBITRUM_1.Selector),
			expected: true,
		},
		{
			src:      models.NetworkSelector(chainsel.ETHEREUM_MAINNET.Selector),
			dest:     models.NetworkSelector(chainsel.ETHEREUM_MAINNET_OPTIMISM_1.Selector),
			expected: true,
		},
		{
			src:      models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector),
			dest:     models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector),
			expected: true,
		},
		{
			src:      models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector),
			dest:     models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA_OPTIMISM_1.Selector),
			expected: true,
		},
		{
			src:      models.NetworkSelector(chainsel.ETHEREUM_MAINNET_ARBITRUM_1.Selector),
			dest:     models.NetworkSelector(chainsel.ETHEREUM_MAINNET.Selector),
			expected: true,
		},
		{
			src:      models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector),
			dest:     models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector),
			expected: true,
		},
		{
			src:      models.NetworkSelector(chainsel.ETHEREUM_MAINNET_OPTIMISM_1.Selector),
			dest:     models.NetworkSelector(chainsel.ETHEREUM_MAINNET.Selector),
			expected: true,
		},
		{
			src:      models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA_OPTIMISM_1.Selector),
			dest:     models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector),
			expected: true,
		},
		{
			src:      models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA_OPTIMISM_1.Selector),
			dest:     models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector),
			expected: false,
		},
		{
			src:      models.NetworkSelector(chainsel.ETHEREUM_MAINNET.Selector),
			dest:     models.NetworkSelector(chainsel.ETHEREUM_MAINNET.Selector),
			expected: false,
		},
		{
			src:      models.NetworkSelector(chainsel.ETHEREUM_MAINNET_ARBITRUM_1.Selector),
			dest:     models.NetworkSelector(chainsel.ETHEREUM_MAINNET_OPTIMISM_1.Selector),
			expected: false,
		},
		{
			src:      models.NetworkSelector(chainsel.ETHEREUM_MAINNET.Selector),
			dest:     models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector),
			expected: false,
		},
		{
			src:      models.NetworkSelector(chainsel.ETHEREUM_MAINNET.Selector),
			dest:     models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA_OPTIMISM_1.Selector),
			expected: false,
		},
		{
			src:      models.NetworkSelector(chainsel.AREON_MAINNET.Selector),
			dest:     models.NetworkSelector(chainsel.ETHEREUM_MAINNET_ARBITRUM_1.Selector),
			expected: false,
		},
		{
			src:      models.NetworkSelector(chainsel.AREON_MAINNET.Selector),
			dest:     models.NetworkSelector(chainsel.AVALANCHE_MAINNET.Selector),
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run("Test", func(t *testing.T) {
			result := Supports(tc.src, tc.dest)
			if result != tc.expected {
				t.Errorf("Supports(%v, %v) = %v, want %v", tc.src, tc.dest, result, tc.expected)
			}
		})
	}
}
