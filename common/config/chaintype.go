package config

import (
	"fmt"
	"strings"
)

// ChainType denotes the chain or network to work with
type ChainType int

const (
	ChainTypeNone ChainType = iota
	ChainArbitrum
	ChainCelo
	ChainGnosis
	ChainKroma
	ChainMetis
	ChainOptimismBedrock
	ChainScroll
	ChainWeMix
	ChainXLayer
	ChainZkSync
)

func ChainTypeFromSlug(slug string) (ChainType, error) {
	switch slug {
	case "":
		return ChainTypeNone, nil
	case "arbitrum":
		return ChainArbitrum, nil
	case "celo":
		return ChainCelo, nil
	case "gnosis", "xdai":
		return ChainGnosis, nil
	case "kroma":
		return ChainKroma, nil
	case "metis":
		return ChainMetis, nil
	case "optimismBedrock":
		return ChainOptimismBedrock, nil
	case "scroll":
		return ChainScroll, nil
	case "wemix":
		return ChainWeMix, nil
	case "xlayer":
		return ChainXLayer, nil
	case "zksync":
		return ChainZkSync, nil
	default:
		return ChainTypeNone, ErrInvalidChainType
	}
}

func (c ChainType) String() string {
	switch c {
	case ChainArbitrum:
		return "arbitrum"
	case ChainCelo:
		return "celo"
	case ChainGnosis:
		return "gnosis"
	case ChainKroma:
		return "kroma"
	case ChainMetis:
		return "metis"
	case ChainOptimismBedrock:
		return "optimismBedrock"
	case ChainScroll:
		return "scroll"
	case ChainWeMix:
		return "wemix"
	case ChainXLayer:
		return "xlayer"
	case ChainZkSync:
		return "zksync"
	default:
		return ""
	}
}

var ErrInvalidChainType = fmt.Errorf("must be one of %s or omitted", strings.Join([]string{
	ChainArbitrum.String(),
	ChainCelo.String(),
	ChainGnosis.String(),
	ChainKroma.String(),
	ChainMetis.String(),
	ChainOptimismBedrock.String(),
	ChainScroll.String(),
	ChainWeMix.String(),
	ChainXLayer.String(),
	ChainZkSync.String(),
}, ", "))

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
