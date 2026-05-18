package keeper

import (
	"math"
	"math/big"
	"strings"

	ethmath "github.com/ethereum/go-ethereum/common/math"
)

const (
	ZeroPrefix   = "0x"
	UpkeepPrefix = "UPx"
)

// LeastSignificant32 returns the least significant 32 bits of the input as a big int
func LeastSignificant32(num *big.Int) uint64 {
	max32 := big.NewInt(math.MaxUint32)
	return big.NewInt(0).And(num, max32).Uint64()
}

// ParseUpkeepId parses the upkeep id input string to a big int pointer. It can handle the following 4 formats:
// 1. decimal format like 123471239047239047243709...
// 2. hex format like AbC13D354eFF...
// 3. 0x-prefixed hex like 0xAbC13D354eFF...
// 4. Upkeep-prefixed hex like UPxAbC13D354eFF...
func ParseUpkeepId(upkeepIdStr string) (*big.Int, bool) {
	if strings.HasPrefix(upkeepIdStr, UpkeepPrefix) {
		upkeepIdStr = ZeroPrefix + upkeepIdStr[len(UpkeepPrefix):]
	}

	// this handles cases 1, 3, 4
	upkeepId, ok := ethmath.ParseBig256(upkeepIdStr)
	if !ok {
		// this handles case 2 or returns (nil, false)
		return ethmath.ParseBig256(ZeroPrefix + upkeepIdStr)
	}
	return upkeepId, ok
}
