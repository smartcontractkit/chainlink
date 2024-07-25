package common

import (
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

type Chains map[uint64][]uint64

func Supports(src, dest models.NetworkSelector) bool {
	if chains[uint64(src)] == nil {
		return false
	}
	for _, d := range chains[uint64(src)] {
		if d == uint64(dest) {
			return true
		}
	}
	return false
}

// Supported source chain -> destination chains for bridge transfers
var chains = Chains{
	// Source = Ethereum
	chainsel.ETHEREUM_MAINNET.Selector: []uint64{
		chainsel.ETHEREUM_MAINNET_ARBITRUM_1.Selector,
		chainsel.ETHEREUM_MAINNET_OPTIMISM_1.Selector,
	},
	chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector: []uint64{
		chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector,
		chainsel.ETHEREUM_TESTNET_SEPOLIA_OPTIMISM_1.Selector,
	},
	// Source = Arbitrum
	chainsel.ETHEREUM_MAINNET_ARBITRUM_1.Selector: []uint64{
		chainsel.ETHEREUM_MAINNET.Selector,
	},
	chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector: []uint64{
		chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector,
	},
	// Source = Optimism
	chainsel.ETHEREUM_MAINNET_OPTIMISM_1.Selector: []uint64{
		chainsel.ETHEREUM_MAINNET.Selector,
	},
	chainsel.ETHEREUM_TESTNET_SEPOLIA_OPTIMISM_1.Selector: []uint64{
		chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector,
	},
}
