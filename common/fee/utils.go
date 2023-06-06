package fee

import (
	"math/big"

	bigmath "github.com/smartcontractkit/chainlink/v2/core/utils/big_math"

	"github.com/shopspring/decimal"
)

func GetMaxFeePrice(userSpecifiedMax, maxFee *big.Int) *big.Int {
	return bigmath.Min(userSpecifiedMax, maxFee)
}

func CapFeePrice(calculatedFeePrice, userSpecifiedMax, maxFeePriceWei *big.Int, feeLimit uint32, multiplier float32) (maxFeePrice *big.Int, chainSpecificFeeLimit uint32) {
	chainSpecificFeeLimit = ApplyMultiplier(feeLimit, multiplier)
	maxFeePrice = GetMaxFeePrice(userSpecifiedMax, maxFeePriceWei)
	return bigmath.Min(calculatedFeePrice, maxFeePrice), chainSpecificFeeLimit
}

func ApplyMultiplier(feeLimit uint32, multiplier float32) uint32 {
	return uint32(decimal.NewFromBigInt(big.NewInt(0).SetUint64(uint64(feeLimit)), 0).Mul(decimal.NewFromFloat32(multiplier)).IntPart())
}
