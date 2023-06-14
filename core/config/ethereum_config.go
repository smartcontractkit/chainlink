package config

import (
	"math/big"
)

type Ethereum interface {
	DefaultChainID() *big.Int
}
