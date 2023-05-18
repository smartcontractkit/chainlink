package config

import (
	"math/big"
	"net/url"
)

type Ethereum interface {
	EthereumHTTPURL() *url.URL
	EthereumSecondaryURLs() []url.URL
	EthereumURL() string
	DefaultChainID() *big.Int
}
