package assets

import (
	"math/big"

	"golang.org/x/exp/constraints"

	"github.com/ethereum/go-ethereum/params"
)

func ItoWei[T constraints.Integer](n T) *Wei {
	w := big.NewInt(int64(n))
	w.Mul(w, big.NewInt(params.Wei))
	return NewWei(w)
}

func ItoGWei[T constraints.Integer](n T) *Wei {
	w := big.NewInt(int64(n))
	w.Mul(w, big.NewInt(params.GWei))
	return NewWei(w)
}

// ItoUEther converts units of micro-ether (terawei) into wei
func ItoUEther[T constraints.Integer](n T) *Wei {
	w := big.NewInt(int64(n))
	w.Mul(w, big.NewInt(params.GWei*1000))
	return NewWei(w)
}

func ItoEther[T constraints.Integer](n T) *Wei {
	w := big.NewInt(int64(n))
	w.Mul(w, big.NewInt(params.Ether))
	return NewWei(w)
}
