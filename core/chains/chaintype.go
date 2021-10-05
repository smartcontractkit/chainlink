package chains

type ChainType string

const (
	Arbitrum ChainType = "arbitrum"
	ExChain  ChainType = "exchain"
	Optimism ChainType = "optimism"
	XDai     ChainType = "xdai"
)

// IsValid returns true if the ChainType value is known or empty.
func (c ChainType) IsValid() bool {
	switch c {
	case "", Arbitrum, ExChain, Optimism, XDai:
		return true
	}
	return false
}

// IsL2 returns true if this chain is a Layer 2 chain. Notably the block numbers
// used for log searching are different from calling block.number
func (c ChainType) IsL2() bool {
	switch c {
	case Arbitrum, Optimism:
		return true

	case XDai, ExChain:
		fallthrough
	default:
		return false
	}
}
