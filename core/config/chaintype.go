package config

// ChainType denotes the chain or network to work with
type ChainType string

//nolint
const (
	ChainArbitrum ChainType = "arbitrum"
	ChainExChain  ChainType = "exchain"
	ChainOptimism ChainType = "optimism"
	ChainXDai     ChainType = "xdai"
)

// IsValid returns true if the ChainType value is known or empty.
func (c ChainType) IsValid() bool {
	switch c {
	case "", ChainArbitrum, ChainExChain, ChainOptimism, ChainXDai:
		return true
	}
	return false
}

// IsL2 returns true if this chain is a Layer 2 chain. Notably the block numbers
// used for log searching are different from calling block.number
func (c ChainType) IsL2() bool {
	switch c {
	case ChainArbitrum, ChainOptimism:
		return true

	case ChainXDai, ChainExChain:
		fallthrough
	default:
		return false
	}
}
