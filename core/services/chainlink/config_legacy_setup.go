package chainlink

import (
	"math/big"
	"net/url"
)

func (l *legacyGeneralConfig) DefaultChainID() *big.Int         { panic("unimplemented") }
func (l *legacyGeneralConfig) EthereumHTTPURL() *url.URL        { panic("unimplemented") }
func (l *legacyGeneralConfig) EthereumNodes() string            { panic("unimplemented") }
func (l *legacyGeneralConfig) EthereumSecondaryURLs() []url.URL { panic("unimplemented") }
func (l *legacyGeneralConfig) EthereumURL() string              { panic("unimplemented") }

func (l *legacyGeneralConfig) SolanaNodes() string   { panic("unimplemented") }
func (l *legacyGeneralConfig) TerraNodes() string    { panic("unimplemented") }
func (l *legacyGeneralConfig) StarkNetNodes() string { panic("unimplemented") }
