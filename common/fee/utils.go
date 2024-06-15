package fee

import (
	"fmt"
	"math"
	"math/big"

	"github.com/shopspring/decimal"
)

func ApplyMultiplier(feeLimit uint64, multiplier float32) (uint64, error) {
	result := decimal.NewFromBigInt(big.NewInt(0).SetUint64(feeLimit), 0).Mul(decimal.NewFromFloat32(multiplier))

	if result.GreaterThan(decimal.NewFromBigInt(big.NewInt(0).SetUint64(math.MaxUint64), 0)) {
		return 0, fmt.Errorf("integer overflow when applying multiplier of %f to fee limit of %d", multiplier, feeLimit)
	}
	return result.BigInt().Uint64(), nil
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
