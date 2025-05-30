package conversions

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
)

func MustSafeInt(input uint64) int {
	maxInt := uint64(1<<(strconv.IntSize-1) - 1) // Max value for int (platform dependent)
	if input > maxInt {
		panic(fmt.Errorf("uint64 %d exceeds int max value", input))
	}
	return int(input)
}

func MustSafeUint64(input int64) uint64 {
	if input < 0 {
		panic(fmt.Errorf("int64 %d is below uint64 min value", input))
	}
	// No need for max value check since int64's max value is always less than uint64's max value
	return uint64(input)
}

func MustSafeInt64(input uint64) int64 {
	if input > math.MaxInt64 {
		panic(fmt.Errorf("uint64 %d exceeds int64 max value", input))
	}
	return int64(input)
}

func MustSafeUint32(input int) uint32 {
	if input < 0 {
		panic(fmt.Errorf("int %d is below uint32 min value", input))
	}
	if input > math.MaxUint32 {
		panic(fmt.Errorf("int %d exceeds uint32 max value", input))
	}
	return uint32(input)
}

func MustSafeUint8(input int) uint8 {
	if input < 0 {
		panic(fmt.Errorf("int %d is below uint8 min value", input))
	}
	if input > math.MaxUint8 {
		panic(fmt.Errorf("int %d exceeds uint8 max value", input))
	}
	return uint8(input)
}

func Float64ToBigInt(f float64) *big.Int {
	bigFloat := new(big.Float).SetFloat64(f)

	bigInt := new(big.Int)
	bigFloat.Int(bigInt) // Truncate towards zero

	return bigInt
}
