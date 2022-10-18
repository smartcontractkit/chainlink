package config

import (
	"fmt"
	"strings"
)

// ChainType denotes the chain or network to work with
type ChainType string

// nolint
const (
	ChainArbitrum        ChainType = "arbitrum"
	ChainMetis           ChainType = "metis"
	ChainOptimism        ChainType = "optimism"
	ChainOptimismBedrock ChainType = "optimismBedrock"
	ChainXDai            ChainType = "xdai"
)

var ErrInvalidChainType = fmt.Errorf("must be one of %s or omitted", strings.Join([]string{
	string(ChainArbitrum), string(ChainMetis), string(ChainOptimism), string(ChainXDai), string(ChainOptimismBedrock),
}, ", "))

// IsValid returns true if the ChainType value is known or empty.
func (c ChainType) IsValid() bool {
	switch c {
	case "", ChainArbitrum, ChainMetis, ChainOptimism, ChainOptimismBedrock, ChainXDai:
		return true
	}
	return false
}

// IsL2 returns true if this chain is a Layer 2 chain. Notably:
//   - the block numbers used for log searching are different from calling block.number
//   - gas bumping is not supported, since there is no tx mempool
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
