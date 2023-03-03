package ethgo

import "math/big"

func convert(val uint64, decimals int64) *big.Int {
	v := big.NewInt(int64(val))
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil)
	return v.Mul(v, exp)
}

// Ether converts a value to the ether unit with 18 decimals
func Ether(i uint64) *big.Int {
	return convert(i, 18)
}

// Gwei converts a value to the gwei unit with 9 decimals
func Gwei(i uint64) *big.Int {
	return convert(i, 9)
}
