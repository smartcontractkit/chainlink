package config

import "fmt"

// TODO Replace with https://github.com/smartcontractkit/ccip-chain-selectors after making it public
var EvmChainIdToChainSelector = map[uint64]uint64{
	// Testnets
	97:       13264668187771770619, // BSC Testnet
	420:      2664363617261496610,  // Optimism Goerli
	1000:     11787463284727550157, // Tests
	1337:     3379446385462418246,  // Tests
	2337:     12922642891491394802, // Tests
	43113:    14767482510784806043, // Avalanche Fuji
	76578:    781901677223027175,   // AvaxAnzSubnet
	84531:    5790810961207155433,  // BASE Goerli
	80001:    12532609583862916517, // Polygon Mumbai
	421613:   6101244977088475029,  // Arbitrum Goerli
	11155111: 16015286601757825753, // Sepolia

	// Mainnets
	1:     5009297550715157269,  // Ethereum
	10:    3734403246176062136,  // Optimism
	56:    11344663589394136015, // BSC
	137:   4051577828743386545,  // Polygon
	8453:  15971525489660198786, // BASE
	42161: 4949039107694359620,  // Arbitrum
	43114: 6433500567565415381,  // Avalanche
}

func ChainIdFromSelector(chainSelectorId uint64) (uint64, error) {
	for k, v := range EvmChainIdToChainSelector {
		if v == chainSelectorId {
			return k, nil
		}
	}
	return 0, fmt.Errorf("chain not found for chain selector %d", chainSelectorId)
}

func SelectorFromChainId(chainId uint64) (uint64, error) {
	if chainSelectorId, exist := EvmChainIdToChainSelector[chainId]; exist {
		return chainSelectorId, nil
	}
	return 0, fmt.Errorf("chain selector not found for chain %d", chainId)
}

// CommitPluginConfig contains the plugin specific variables for the ccip.CCIPCommit plugin.
// We use ID here to keep it as general as possible, e.g. abstracting for chains which don't have an address concept.
type CommitPluginConfig struct {
	SourceStartBlock, DestStartBlock int64  // Only for first time job add.
	OffRamp                          string `json:"offRamp"`
	// TokenPricesUSDPipeline should contain a token price pipeline for the following tokens:
	//		The SOURCE chain wrapped native
	// 		The DESTINATION supported tokens (including fee tokens) as defined in destination OffRamp and PriceRegistry.
	TokenPricesUSDPipeline string `json:"tokenPricesUSDPipeline"`
}

// ExecutionPluginConfig contains the plugin specific variables for the ccip.CCIPExecution plugin.
type ExecutionPluginConfig struct {
	SourceStartBlock, DestStartBlock int64 // Only for first time job add.
}
