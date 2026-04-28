package handler

import (
	"math"
	"math/big"
	"strings"

	ethmath "github.com/ethereum/go-ethereum/common/math"
)

const (
	zeroPrefix   = "0x"
	upkeepPrefix = "UPx"
)

// LeastSignificant32 returns the least significant 32 bits of the input as a uint64.
func LeastSignificant32(num *big.Int) uint64 {
	max32 := big.NewInt(math.MaxUint32)
	return big.NewInt(0).And(num, max32).Uint64()
}

// ParseUpkeepId parses the upkeep id input string to a big int pointer.
func ParseUpkeepId(upkeepIdStr string) (*big.Int, bool) {
	if strings.HasPrefix(upkeepIdStr, upkeepPrefix) {
		upkeepIdStr = zeroPrefix + upkeepIdStr[len(upkeepPrefix):]
	}

	upkeepId, ok := ethmath.ParseBig256(upkeepIdStr)
	if !ok {
		return ethmath.ParseBig256(zeroPrefix + upkeepIdStr)
	}
	return upkeepId, ok
}
