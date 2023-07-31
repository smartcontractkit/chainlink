package fee

import (
	"math/big"

	bigmath "github.com/smartcontractkit/chainlink/v2/core/utils/big_math"

	"github.com/shopspring/decimal"
)

func FeePriceLimit(userSpecifiedMax, maxFee *big.Int) *big.Int {
	return bigmath.Min(userSpecifiedMax, maxFee)
}

func CapFeePrice(calculatedFeePrice, userSpecifiedMax, maxFeePriceWei *big.Int, feeLimit uint32, multiplier float32) (maxFeePrice *big.Int, chainSpecificFeeLimit uint32) {
	chainSpecificFeeLimit = ApplyMultiplier(feeLimit, multiplier)
	maxFeePrice = FeePriceLimit(userSpecifiedMax, maxFeePriceWei)
	return bigmath.Min(calculatedFeePrice, maxFeePrice), chainSpecificFeeLimit
}

func ApplyMultiplier(feeLimit uint32, multiplier float32) uint32 {
	return uint32(decimal.NewFromBigInt(big.NewInt(0).SetUint64(uint64(feeLimit)), 0).Mul(decimal.NewFromFloat32(multiplier)).IntPart())
}

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
