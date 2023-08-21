package deployments

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rhea"
)

const (
	BATCH_GAS_LIMIT                     = 5_000_000
	FEE_UPDATE_HEARTBEAT                = 24 * time.Hour
	FEE_UPDATE_DEVIATION_PPB            = 10e7 // 10%
	FEE_UPDATE_DEVIATION_PPB_FAST_CHAIN = 20e7 // 20%
	// This boosts the fee paid every 10x every 15s, since fees artificially low (0.1 on testnet) and
	// we have source finality artificially low. On fast chains transactions become available for execution
	// within 30s - 1min, and we want to avoid waiting for a full root snooze so we make sure they
	// are boosted back to cost immediately.
	RELATIVE_BOOST_PER_WAIT_HOUR = 2400
	INFLIGHT_CACHE_EXPIRY        = 3 * time.Minute
	ROOT_SNOOZE_TIME             = 5 * time.Minute
)

func getFinalityDepth(chain rhea.Chain) uint32 {
	// NOTE most of these is still way artificially low, but we aim for quick iteration on testnet.
	// Optimism, polygon and arbitrum in particular are known for decent sized reorgs, so we set
	// those higher than others.
	var finalityDepthPerChain = map[rhea.Chain]uint32{
		// Testnets
		rhea.Goerli:         4, // We use finality tag
		rhea.OptimismGoerli: 5, // We use finality tag
		rhea.Quorum:         4,
		rhea.AvaxFuji:       35,  // Implements strange finality tag, use 35 instead
		rhea.PolygonMumbai:  550, // No finality tag available
		rhea.ArbitrumGoerli: 5,   // We use finality tag
		rhea.Sepolia:        4,   // We use finality tag
		rhea.BASEGoerli:     5,   // We use finality tag
		rhea.AvaxAnzSubnet:  1,   // devnet without normal block production
		rhea.BSCTestnet:     2,   //  We use finality tag
		// Mainnets
		// We've made a bit of an effort to make these realistic, but they are *not* production-ready parameters, only use them for testing!!!
		rhea.Ethereum: 96, // 3 epochs to be safe
		rhea.Avax:     2,  // Should be 1 theoretically
		rhea.Polygon:  600,
		rhea.Optimism: 5000, // Should be 1 theoretically
		rhea.Arbitrum: 5000,
	}

	if val, ok := finalityDepthPerChain[chain]; ok {
		return val
	}
	panic(fmt.Sprintf("Finality depth for %s not found", chain))
}

func getOptimisticConfirmations(chain rhea.Chain) uint32 {
	var optimisticConfirmations = map[rhea.Chain]uint32{
		// Testnets
		rhea.Goerli:         4,
		rhea.Sepolia:        4,
		rhea.OptimismGoerli: 4,
		rhea.AvaxFuji:       1,
		rhea.PolygonMumbai:  4,
		rhea.ArbitrumGoerli: 1,
		rhea.Quorum:         1,
		rhea.BASEGoerli:     4,
		rhea.AvaxAnzSubnet:  1, // devnet without normal block production
		rhea.BSCTestnet:     2,
		// Mainnets
		rhea.Ethereum: 2,
		rhea.Avax:     1,
		rhea.Polygon:  20,
		rhea.Optimism: 1,
		rhea.Arbitrum: 1,
	}

	if val, ok := optimisticConfirmations[chain]; ok {
		return val
	}
	panic(fmt.Sprintf("Optimistic confirmations for %s not found", chain))
}

func getMaxGasPrice(chain rhea.Chain) uint64 {
	var maxGasPricePerChain = map[rhea.Chain]uint64{
		// Testnets
		rhea.Goerli:         200e9,
		rhea.Sepolia:        200e9,
		rhea.OptimismGoerli: 200e9,
		rhea.AvaxFuji:       200e9,
		rhea.PolygonMumbai:  200e9,
		rhea.ArbitrumGoerli: 200e9,
		rhea.Quorum:         200e9,
		rhea.AvaxAnzSubnet:  200e9,
		rhea.BASEGoerli:     200e9,
		rhea.BSCTestnet:     200e9,
		// Mainnets
		rhea.Ethereum: 200e9,
		rhea.Avax:     200e9,
		rhea.Polygon:  200e9,
		rhea.Optimism: 200e9,
		rhea.Arbitrum: 200e9,
	}

	if val, ok := maxGasPricePerChain[chain]; ok {
		return val
	}
	panic(fmt.Sprintf("Max gas price for %s not found", chain))
}
