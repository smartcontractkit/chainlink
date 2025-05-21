package conversions

import (
	"fmt"
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
	if input > int64(^uint64(0)>>1) {
		panic(fmt.Errorf("int64 %d exceeds uint64 max value", input))
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

func MustSafeUint8(input int) uint8 {
	if input < 0 {
		panic(fmt.Errorf("int %d is below uint8 min value", input))
	}
	maxUint8 := (1 << 8) - 1
	if input > maxUint8 {
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
