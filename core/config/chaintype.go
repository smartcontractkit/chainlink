package config

// ChainType denotes the chain or network to work with
type ChainType string

//nolint
const (
	ChainArbitrum ChainType = "arbitrum"
	ChainMetis    ChainType = "metis"
	ChainOptimism ChainType = "optimism"
	ChainXDai     ChainType = "xdai"
)

// IsValid returns true if the ChainType value is known or empty.
func (c ChainType) IsValid() bool {
	switch c {
	case "", ChainArbitrum, ChainMetis, ChainOptimism, ChainXDai:
		return true
	}
	return false
}

// IsOptimismClone returns true if chain is Optimism or a clone
func (c ChainType) IsOptimismClone() bool {
	switch c {
	case ChainMetis, ChainOptimism:
		return true
	default:
		return false
	}
}

// IsL2 returns true if this chain is a Layer 2 chain. Notably the block numbers
// used for log searching are different from calling block.number
func (c ChainType) IsL2() bool {
	switch c {
	case ChainArbitrum, ChainMetis, ChainOptimism:
		return true

	case ChainXDai:
		fallthrough
	default:
		return false
	}
}
