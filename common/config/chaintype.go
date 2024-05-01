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
	ChainCelo            ChainType = "celo"
	ChainGnosis          ChainType = "gnosis"
	ChainKroma           ChainType = "kroma"
	ChainMetis           ChainType = "metis"
	ChainOptimismBedrock ChainType = "optimismBedrock"
	ChainScroll          ChainType = "scroll"
	ChainWeMix           ChainType = "wemix"
	ChainXDai            ChainType = "xdai" // Deprecated: use ChainGnosis instead
	ChainXLayer          ChainType = "xlayer"
	ChainZkSync          ChainType = "zksync"
)

var ErrInvalidChainType = fmt.Errorf("must be one of %s or omitted", strings.Join([]string{
	string(ChainArbitrum),
	string(ChainCelo),
	string(ChainGnosis),
	string(ChainKroma),
	string(ChainMetis),
	string(ChainOptimismBedrock),
	string(ChainScroll),
	string(ChainWeMix),
	string(ChainXLayer),
	string(ChainZkSync),
}, ", "))

// IsValid returns true if the ChainType value is known or empty.
func (c ChainType) IsValid() bool {
	switch c {
	case "", ChainArbitrum, ChainCelo, ChainGnosis, ChainKroma, ChainMetis, ChainOptimismBedrock, ChainScroll, ChainWeMix, ChainXDai, ChainXLayer, ChainZkSync:
		return true
	}
	return false
}

// IsL2 returns true if this chain is a Layer 2 chain. Notably:
//   - the block numbers used for log searching are different from calling block.number
//   - gas bumping is not supported, since there is no tx mempool
func (c ChainType) IsL2() bool {
	switch c {
	case ChainArbitrum, ChainMetis:
		return true
	default:
		return false
	}
}
