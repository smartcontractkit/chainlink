package common

const (
	ArbitrumGoerliChainID  int64 = 421613
	ArbitrumSepoliaChainID int64 = 421614
	ArbitrumOneChainID     int64 = 42161
)

// IsArbitrumChainID returns true if and only if the given chain ID corresponds
// to an Arbitrum chain (testnet or mainnet).
func IsArbitrumChainID(chainID int64) bool {
	return chainID == ArbitrumGoerliChainID || chainID == ArbitrumSepoliaChainID || chainID == ArbitrumOneChainID
}
