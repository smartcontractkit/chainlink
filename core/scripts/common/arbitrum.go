package common

const (
	ArbitrumGoerliChainID  int64 = 421613
	ArbitrumOneChainID     int64 = 42161
	ArbitrumSepoliaChainID int64 = 421614
	OptimismChainID        int64 = 10
	OptimismSepoliaChainID int64 = 11155420
	BaseChainID            int64 = 8453
	BaseSepoliaChainID     int64 = 84532
)

// IsArbitrumChainID returns true if and only if the given chain ID corresponds
// to an Arbitrum chain (testnet or mainnet).
func IsArbitrumChainID(chainID int64) bool {
	return chainID == ArbitrumGoerliChainID || chainID == ArbitrumOneChainID || chainID == ArbitrumSepoliaChainID
}

// IsArbitrumChainID returns true if and only if the given chain ID corresponds
// to an Optimism or Base chain (testnet or mainnet).
func IsOptimismOrBaseChainID(chainID int64) bool {
	return chainID == OptimismSepoliaChainID || chainID == OptimismChainID || chainID == BaseSepoliaChainID || chainID == BaseChainID
}
