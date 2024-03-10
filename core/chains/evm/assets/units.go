package assets

import (
	"math/big"

	"golang.org/x/exp/constraints"

	"github.com/ethereum/go-ethereum/params"
)

func GWei[T constraints.Signed](n T) *Wei {
	w := big.NewInt(int64(n))
	w.Mul(w, big.NewInt(params.GWei))
	return NewWei(w)
}

// UEther converts units of micro-ether (terawei) into wei
func UEther[T constraints.Signed](n T) *Wei {
	w := big.NewInt(int64(n))
	w.Mul(w, big.NewInt(params.GWei*1000))
	return NewWei(w)
}

func Ether[T constraints.Signed](n T) *Wei {
	w := big.NewInt(int64(n))
	w.Mul(w, big.NewInt(params.Ether))
	return NewWei(w)
}
