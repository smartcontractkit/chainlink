package chaintype

import (
	"fmt"
	"strings"
)

type ChainType string

const (
	ChainArbitrum        ChainType = "arbitrum"
	ChainCelo            ChainType = "celo"
	ChainGnosis          ChainType = "gnosis"
	ChainKroma           ChainType = "kroma"
	ChainMetis           ChainType = "metis"
	ChainOptimismBedrock ChainType = "optimismBedrock"
	ChainScroll          ChainType = "scroll"
	ChainWeMix           ChainType = "wemix"
	ChainXLayer          ChainType = "xlayer"
	ChainZkEvm           ChainType = "zkevm"
	ChainZkSync          ChainType = "zksync"
)

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

func (c ChainType) IsValid() bool {
	switch c {
	case "", ChainArbitrum, ChainCelo, ChainGnosis, ChainKroma, ChainMetis, ChainOptimismBedrock, ChainScroll, ChainWeMix, ChainXLayer, ChainZkEvm, ChainZkSync:
		return true
	}
	return false
}

func ChainTypeFromSlug(slug string) ChainType {
	switch slug {
	case "arbitrum":
		return ChainArbitrum
	case "celo":
		return ChainCelo
	case "gnosis":
		return ChainGnosis
	case "kroma":
		return ChainKroma
	case "metis":
		return ChainMetis
	case "optimismBedrock":
		return ChainOptimismBedrock
	case "scroll":
		return ChainScroll
	case "wemix":
		return ChainWeMix
	case "xlayer":
		return ChainXLayer
	case "zkevm":
		return ChainZkEvm
	case "zksync":
		return ChainZkSync
	default:
		return ChainType(slug)
	}
}

type ChainTypeConfig struct {
	value ChainType
	slug  string
}

func NewChainTypeConfig(slug string) *ChainTypeConfig {
	return &ChainTypeConfig{
		value: ChainTypeFromSlug(slug),
		slug:  slug,
	}
}

func (c *ChainTypeConfig) MarshalText() ([]byte, error) {
	if c == nil {
		return nil, nil
	}
	return []byte(c.slug), nil
}

func (c *ChainTypeConfig) UnmarshalText(b []byte) error {
	c.slug = string(b)
	c.value = ChainTypeFromSlug(c.slug)
	return nil
}

func (c *ChainTypeConfig) Slug() string {
	if c == nil {
		return ""
	}
	return c.slug
}

func (c *ChainTypeConfig) ChainType() ChainType {
	if c == nil {
		return ""
	}
	return c.value
}

func (c *ChainTypeConfig) String() string {
	if c == nil {
		return ""
	}
	return string(c.value)
}

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
	string(ChainZkEvm),
	string(ChainZkSync),
}, ", "))
