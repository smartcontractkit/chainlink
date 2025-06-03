package config

import (
	"math/big"
)

type TokenParams struct {
	MaxSupply *big.Int
	Name      string
	Symbol    string
	Decimals  byte
	Icon      string
	Project   string
}
