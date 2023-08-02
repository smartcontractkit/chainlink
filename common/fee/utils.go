package fee

import (
	"math/big"

	"github.com/shopspring/decimal"
)

func ApplyMultiplier(feeLimit uint32, multiplier float32) uint32 {
	return uint32(decimal.NewFromBigInt(big.NewInt(0).SetUint64(uint64(feeLimit)), 0).Mul(decimal.NewFromFloat32(multiplier)).IntPart())
}

// TODO: Remove function when we move utils.Big into `chainlink-relay`
// Max is a duplicate in big.Utils package.
// Max returns the maximum of the two given big.Int values.
func max(x, y *big.Int) *big.Int {
	if x.Cmp(y) == 1 {
		return x
	}
	return y
}

// Returns the input value increased by the given percentage.
func addPercentage(value *big.Int, percentage uint16) *big.Int {
	bumped := new(big.Int)
	bumped.Mul(value, big.NewInt(int64(100+percentage)))
	bumped.Div(bumped, big.NewInt(100))
	return bumped
}

// Returns the fee in its chain specific unit.
type feeUnitToChainUnit func(fee *big.Int) string
