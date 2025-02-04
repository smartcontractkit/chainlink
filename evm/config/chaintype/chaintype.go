package chaintype

import (
	"fmt"
	"strings"
)

type ChainType string

const (
	ChainArbitrum        ChainType = "arbitrum"
	ChainAstar           ChainType = "astar"
	ChainCelo            ChainType = "celo"
	ChainGnosis          ChainType = "gnosis"
	ChainHedera          ChainType = "hedera"
	ChainKroma           ChainType = "kroma"
	ChainMantle          ChainType = "mantle"
	ChainMetis           ChainType = "metis"
	ChainOptimismBedrock ChainType = "optimismBedrock"
	ChainSei             ChainType = "sei"
	ChainScroll          ChainType = "scroll"
	ChainWeMix           ChainType = "wemix"
	ChainXLayer          ChainType = "xlayer"
	ChainZkEvm           ChainType = "zkevm"
	ChainZkSync          ChainType = "zksync"
	ChainZircuit         ChainType = "zircuit"
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
	case "", ChainArbitrum, ChainAstar, ChainCelo, ChainGnosis, ChainHedera, ChainKroma, ChainMantle, ChainMetis, ChainOptimismBedrock, ChainSei, ChainScroll, ChainWeMix, ChainXLayer, ChainZkEvm, ChainZkSync, ChainZircuit:
		return true
	}
	return false
}

func FromSlug(slug string) ChainType {
	switch slug {
	case "arbitrum":
		return ChainArbitrum
	case "astar":
		return ChainAstar
	case "celo":
		return ChainCelo
	case "gnosis":
		return ChainGnosis
	case "hedera":
		return ChainHedera
	case "kroma":
		return ChainKroma
	case "mantle":
		return ChainMantle
	case "metis":
		return ChainMetis
	case "optimismBedrock":
		return ChainOptimismBedrock
	case "sei":
		return ChainSei
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
	case "zircuit":
		return ChainZircuit
	default:
		return ChainType(slug)
	}
}

type Config struct {
	value ChainType
	slug  string
}

func NewConfig(slug string) *Config {
	return &Config{
		value: FromSlug(slug),
		slug:  slug,
	}
}

func (c *Config) MarshalText() ([]byte, error) {
	if c == nil {
		return nil, nil
	}
	return []byte(c.slug), nil
}

func (c *Config) UnmarshalText(b []byte) error {
	c.slug = string(b)
	c.value = FromSlug(c.slug)
	return nil
}

func (c *Config) Slug() string {
	if c == nil {
		return ""
	}
	return c.slug
}

func (c *Config) ChainType() ChainType {
	if c == nil {
		return ""
	}
	return c.value
}

func (c *Config) String() string {
	if c == nil {
		return ""
	}
	return string(c.value)
}

var ErrInvalid = fmt.Errorf("must be one of %s or omitted", strings.Join([]string{
	string(ChainArbitrum),
	string(ChainAstar),
	string(ChainCelo),
	string(ChainGnosis),
	string(ChainHedera),
	string(ChainKroma),
	string(ChainMantle),
	string(ChainMetis),
	string(ChainOptimismBedrock),
	string(ChainSei),
	string(ChainScroll),
	string(ChainWeMix),
	string(ChainXLayer),
	string(ChainZkEvm),
	string(ChainZkSync),
	string(ChainZircuit),
}, ", "))
