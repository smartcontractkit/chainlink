package bigint

import "math/big"

func UBigInt(i uint64) *big.Int {
	return new(big.Int).SetUint64(i)
}

func E18Mult(amount uint64) *big.Int {
	return new(big.Int).Mul(UBigInt(amount), UBigInt(1e18))
}

// EDecMult scales amount by the number of decimals
func EDecMult(amount uint64, decimals int64) *big.Int {
	return new(big.Int).Mul(
		UBigInt(amount),
		new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil),
	)
}
