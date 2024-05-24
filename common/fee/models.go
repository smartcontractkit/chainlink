package fee

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/chains/label"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	bigmath "github.com/smartcontractkit/chainlink-common/pkg/utils/big_math"
)

var (
	ErrBumpFeeExceedsLimit = errors.New("fee bump exceeds limit")
	ErrBump                = errors.New("fee bump failed")
	ErrConnectivity        = errors.New("transaction propagation issue: transactions are not being mined")
)

func IsBumpErr(err error) bool {
	return err != nil && (errors.Is(err, ErrBumpFeeExceedsLimit) || errors.Is(err, ErrBump) || errors.Is(err, ErrConnectivity))
}

// CalculateFee computes the fee price for a transaction.
// The fee price is the minimum of:
// - max fee price specified, default fee price and max fee price for the node.
func CalculateFee(
	maxFeePrice, defaultPrice, maxFeePriceConfigured *big.Int,
) *big.Int {
	maxFeePriceAllowed := bigmath.Min(maxFeePrice, maxFeePriceConfigured)
	return bigmath.Min(defaultPrice, maxFeePriceAllowed)
}

// CalculateBumpedFee computes the next fee price to attempt as the largest of:
// - A configured percentage bump (bumpPercent) on top of the baseline price.
// - A configured fixed amount of Unit (bumpMin) on top of the baseline price.
// The baseline price is the maximum of the previous fee price attempt and the node's current fee price.
func CalculateBumpedFee(
	lggr logger.SugaredLogger,
	currentfeePrice, originalfeePrice, maxFeePriceInput, maxBumpPrice, bumpMin *big.Int,
	bumpPercent uint16,
	toChainUnit feeUnitToChainUnit,
) (*big.Int, error) {
	maxFeePrice := bigmath.Min(maxFeePriceInput, maxBumpPrice)
	bumpedFeePrice := MaxBumpedFee(originalfeePrice, bumpPercent, bumpMin)

	// Update bumpedFeePrice if currentfeePrice is higher than bumpedFeePrice and within maxFeePrice
	bumpedFeePrice = maxFee(lggr, currentfeePrice, bumpedFeePrice, maxFeePrice, "fee price", toChainUnit)

	if bumpedFeePrice.Cmp(maxFeePrice) > 0 {
		return maxFeePrice, fmt.Errorf("bumped fee price of %s would exceed configured max fee price of %s (original price was %s). %s: %w",
			toChainUnit(bumpedFeePrice), toChainUnit(maxFeePrice), toChainUnit(originalfeePrice), label.NodeConnectivityProblemWarning, ErrBumpFeeExceedsLimit)
	} else if bumpedFeePrice.Cmp(originalfeePrice) == 0 {
		// NOTE: This really shouldn't happen since we enforce minimums for
		// FeeEstimator.BumpPercent and FeeEstimator.BumpMin in the config validation,
		// but it's here anyway for a "belts and braces" approach
		return bumpedFeePrice, fmt.Errorf("bumped fee price of %s is equal to original fee price of %s."+
			" ACTION REQUIRED: This is a configuration error, you must increase either "+
			"FeeEstimator.BumpPercent or FeeEstimator.BumpMin: %w", toChainUnit(bumpedFeePrice), toChainUnit(bumpedFeePrice), ErrBump)
	}
	return bumpedFeePrice, nil
}

// Returns highest bumped fee price of originalFeePrice bumped by fixed units or percentage.
func MaxBumpedFee(originalFeePrice *big.Int, feeBumpPercent uint16, feeBumpUnits *big.Int) *big.Int {
	return bigmath.Max(
		addPercentage(originalFeePrice, feeBumpPercent),
		new(big.Int).Add(originalFeePrice, feeBumpUnits),
	)
}

// Returns the max of currentFeePrice, bumpedFeePrice, and maxFeePrice
func maxFee(lggr logger.SugaredLogger, currentFeePrice, bumpedFeePrice, maxFeePrice *big.Int, feeType string, toChainUnit feeUnitToChainUnit) *big.Int {
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
