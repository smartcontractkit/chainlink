package config

import (
	"errors"
	"math/big"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

type TokenParams struct {
	MaxSupply *big.Int
	Name      string
	Symbol    shared.TokenSymbol
	Decimals  byte
	Icon      string
	Project   string
}

func (tp TokenParams) Validate() error {
	if tp.MaxSupply != nil && tp.MaxSupply.Sign() < 0 {
		return errors.New("maxSupply must be a positive integer")
	}
	if tp.Name == "" {
		return errors.New("name cannot be empty")
	}
	if tp.Symbol == "" {
		return errors.New("symbol cannot be empty")
	}
	if tp.Decimals < 0 || tp.Decimals > 18 {
		return errors.New("decimals must be between 0 and 18")
	}
	return nil
}
