package conversions

import (
	"fmt"
	"math/big"
)

func MustSafeUint64(input int64) uint64 {
	if input < 0 {
		panic(fmt.Errorf("int64 %d is below uint64 min value", input))
	}
	return uint64(input)
}

func MustSafeUint32(input int) uint32 {
	if input < 0 {
		panic(fmt.Errorf("int %d is below uint32 min value", input))
	}
	maxUint32 := (1 << 32) - 1
	if input > maxUint32 {
		panic(fmt.Errorf("int %d exceeds uint32 max value", input))
	}
	return uint32(input)
}

func Float64ToBigInt(f float64) *big.Int {
	f *= 100

	bigFloat := new(big.Float).SetFloat64(f)

	bigInt := new(big.Int)
	bigFloat.Int(bigInt) // Truncate towards zero

	return bigInt
}
