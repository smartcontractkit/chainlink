package chainlink

import (
	"math/big"
	"net/url"

	v2 "github.com/smartcontractkit/chainlink/core/config/v2"
)

func (g *generalConfig) DefaultChainID() *big.Int         { panic(v2.ErrUnsupported) }
func (g *generalConfig) EthereumHTTPURL() *url.URL        { panic(v2.ErrUnsupported) }
func (g *generalConfig) EthereumNodes() string            { panic(v2.ErrUnsupported) }
func (g *generalConfig) EthereumSecondaryURLs() []url.URL { panic(v2.ErrUnsupported) }
func (g *generalConfig) EthereumURL() string              { panic(v2.ErrUnsupported) }

func (g *generalConfig) SolanaNodes() string   { panic(v2.ErrUnsupported) }
func (g *generalConfig) TerraNodes() string    { panic(v2.ErrUnsupported) }
func (g *generalConfig) StarkNetNodes() string { panic(v2.ErrUnsupported) }
