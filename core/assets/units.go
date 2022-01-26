package assets

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

func Wei(n int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(n), big.NewInt(params.Wei))
}

func GWei(n int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(n), big.NewInt(params.GWei))
}

// UEther converts units of micro-ether (terawei) into wei
func UEther(n int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(n), big.NewInt(params.GWei*1000))
}

func Ether(n int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(n), big.NewInt(params.Ether))
}
