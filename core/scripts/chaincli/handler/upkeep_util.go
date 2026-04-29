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

// ParseUpkeepID parses the upkeep id input string to a big int pointer.
func ParseUpkeepID(upkeepIDStr string) (*big.Int, bool) {
	if strings.HasPrefix(upkeepIDStr, upkeepPrefix) {
		upkeepIDStr = zeroPrefix + upkeepIDStr[len(upkeepPrefix):]
	}

	upkeepID, ok := ethmath.ParseBig256(upkeepIDStr)
	if !ok {
		return ethmath.ParseBig256(zeroPrefix + upkeepIDStr)
	}
	return upkeepID, ok
}
