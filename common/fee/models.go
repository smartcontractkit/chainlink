package fee

import (
	"math/big"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/common/chains/label"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	bigmath "github.com/smartcontractkit/chainlink/v2/core/utils/big_math"
)

var (
	ErrBumpFeeExceedsLimit = errors.New("fee bump exceeds limit")
	ErrBump                = errors.New("fee bump failed")
	ErrConnectivity        = errors.New("transaction propagation issue: transactions are not being mined")
)

func IsBumpErr(err error) bool {
	return err != nil && (errors.Is(err, ErrBumpFeeExceedsLimit) || errors.Is(err, ErrBump) || errors.Is(err, ErrConnectivity))
}

// CalculateFee computes the fee price and chain specific fee limit for a transaction.
func CalculateFee(
	maxFeePrice, defaultPrice, maxBumpPrice *big.Int,
) (feePrice *big.Int) {
	feePrice = CapFeePrice(defaultPrice, maxFeePrice, maxBumpPrice)
	return
}

// CalculateBumpedFee will increase the price and apply multiplier to the fee limit.
func CalculateBumpedFee(
	lggr logger.SugaredLogger,
	currentFeePrice, originalFeePrice, maxFeePrice,
	maxBumpPrice, bumpMin *big.Int,
	bumpPercent uint16,
	toChainUnit func(*big.Int) string,
) (*big.Int, error) {
	feePrice, err := bumpFeePrice(lggr, currentFeePrice, originalFeePrice, maxFeePrice, maxBumpPrice, bumpMin, bumpPercent, toChainUnit)
	if err != nil {
		return nil, err
	}
	return feePrice, nil
}

// chainSpecificFeeLimit := ApplyMultiplier(originalFeeLimit, bumpLimitMultiplier)

// bumpfeePrice computes the next fee price to attempt as the largest of:
// - A configured percentage bump (FeeEstimator.BumpPercent) on top of the baseline price.
// - A configured fixed amount of Unit (FEE_PRICE_Unit) on top of the baseline price.
// The baseline price is the maximum of the previous fee price attempt and the node's current fee price.
func bumpFeePrice(
	lggr logger.SugaredLogger,
	currentfeePrice, originalfeePrice, maxFeePriceInput, maxBumpPrice, bumpMin *big.Int,
	bumpPercent uint16,
	toChainUnit FeeUnitToChainUnit,
) (*big.Int, error) {
	maxFeePrice := bigmath.Min(maxFeePriceInput, maxBumpPrice)
	bumpedFeePrice := maxBumpedFee(originalfeePrice, bumpPercent, bumpMin)

	// Update bumpedFeePrice if currentfeePrice is higher than bumpedFeePrice and within maxFeePrice
	bumpedFeePrice = maxFee(lggr, currentfeePrice, bumpedFeePrice, maxFeePrice, "fee price", toChainUnit)

	if bumpedFeePrice.Cmp(maxFeePrice) > 0 {
		return maxFeePrice, errors.Wrapf(ErrBumpFeeExceedsLimit, "bumped fee price of %s would exceed configured max fee price of %s (original price was %s). %s",
			toChainUnit(bumpedFeePrice), maxFeePrice, toChainUnit(originalfeePrice), label.NodeConnectivityProblemWarning)
	} else if bumpedFeePrice.Cmp(originalfeePrice) == 0 {
		// NOTE: This really shouldn't happen since we enforce minimums for
		// FeeEstimator.BumpPercent and FeeEstimator.BumpMin in the config validation,
		// but it's here anyway for a "belts and braces" approach
		return bumpedFeePrice, errors.Wrapf(ErrBump, "bumped fee price of %s is equal to original fee price of %s."+
			" ACTION REQUIRED: This is a configuration error, you must increase either "+
			"FeeEstimator.BumpPercent or FeeEstimator.BumpMin", toChainUnit(bumpedFeePrice), toChainUnit(bumpedFeePrice))
	}
	return bumpedFeePrice, nil
}

// Returns highest bumped fee price of originalFeePrice bumped by fixed units or percentage.
func maxBumpedFee(originalFeePrice *big.Int, feeBumpPercent uint16, feeBumpUnits *big.Int) *big.Int {
	return max(
		addPercentage(originalFeePrice, feeBumpPercent),
		new(big.Int).Add(originalFeePrice, feeBumpUnits),
	)
}

// Returns the max of currentFeePrice, bumpedFeePrice, and maxFeePrice
func maxFee(lggr logger.SugaredLogger, currentFeePrice, bumpedFeePrice, maxFeePrice *big.Int, feeType string, toChainUnit FeeUnitToChainUnit) *big.Int {
	if currentFeePrice == nil {
		return bumpedFeePrice
	}
	if currentFeePrice.Cmp(maxFeePrice) > 0 {
		// Shouldn't happen because the estimator should not be allowed to
		// estimate a higher fee than the maximum allowed
		lggr.AssumptionViolationf("Ignoring current %s of %s that would exceed max %s of %s", feeType, toChainUnit(currentFeePrice), feeType, toChainUnit(maxFeePrice))
	} else if bumpedFeePrice.Cmp(currentFeePrice) < 0 {
		// If the current fee price is higher than the old price bumped, use that instead
		return currentFeePrice
	}
	return bumpedFeePrice
}
