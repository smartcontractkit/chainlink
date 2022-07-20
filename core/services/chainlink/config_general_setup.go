package chainlink

import (
	"math/big"
	"net/url"
)

func (g *generalConfig) DefaultChainID() *big.Int         { panic("unimplemented") }
func (g *generalConfig) EthereumHTTPURL() *url.URL        { panic("unimplemented") }
func (g *generalConfig) EthereumNodes() string            { panic("unimplemented") }
func (g *generalConfig) EthereumSecondaryURLs() []url.URL { panic("unimplemented") }
func (g *generalConfig) EthereumURL() string              { panic("unimplemented") }

func (g *generalConfig) SolanaNodes() string   { panic("unimplemented") }
func (g *generalConfig) TerraNodes() string    { panic("unimplemented") }
func (g *generalConfig) StarkNetNodes() string { panic("unimplemented") }
