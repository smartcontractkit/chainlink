package config

import (
	"errors"
	"math/big"

	"github.com/aptos-labs/aptos-go-sdk"

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
	if tp.MaxSupply != nil && tp.MaxSupply.Sign() <= 0 {
		return errors.New("maxSupply must be a positive integer or nil")
	}
	if tp.Name == "" {
		return errors.New("name cannot be empty")
	}
	if tp.Symbol == "" {
		return errors.New("symbol cannot be empty")
	}
	if tp.Decimals < 1 || tp.Decimals > 8 {
		return errors.New("decimals must be between 1 and 8")
	}
	return nil
}

type TokenMint struct {
	Amount uint64
	To     aptos.AccountAddress
}
