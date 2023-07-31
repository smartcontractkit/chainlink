package fee

import (
	"math/big"

	"github.com/pkg/errors"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/label" // TODO: Remove import from core
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var (
	ErrBumpFeeExceedsLimit = errors.New("fee bump exceeds limit")
	ErrBump                = errors.New("fee bump failed")
	ErrConnectivity        = errors.New("transaction propagation issue: transactions are not being mined")
)

func IsBumpErr(err error) bool {
	return err != nil && (errors.Is(err, ErrBumpFeeExceedsLimit) || errors.Is(err, ErrBump) || errors.Is(err, ErrConnectivity))
}

// CalculateBumpedFee will increase the price and apply multiplier to the fee limit
func CalculateBumpedFee(cfg feetypes.BumpConfig, lggr logger.SugaredLogger, currentfeePrice, originalfeePrice *big.Int, originalFeeLimit uint32, maxFeePrice *big.Int) (*big.Int, uint32, error) {
	feePrice, err := bumpFeePrice(cfg, lggr, currentfeePrice, originalfeePrice, maxFeePrice)
	if err != nil {
		return nil, 0, err
	}
	chainSpecificFeeLimit := ApplyMultiplier(originalFeeLimit, cfg.LimitMultiplier())
	return feePrice, chainSpecificFeeLimit, nil
}

// bumpfeePrice computes the next fee price to attempt as the largest of:
// - A configured percentage bump (FeeEstimator.BumpPercent) on top of the baseline price.
// - A configured fixed amount of Unit (FEE_PRICE_Unit) on top of the baseline price.
// The baseline price is the maximum of the previous fee price attempt and the node's current fee price.
func bumpFeePrice(cfg feetypes.BumpConfig, lggr logger.SugaredLogger, currentfeePrice, originalfeePrice, maxFeePriceInput *big.Int) (*big.Int, error) {
	maxFeePrice := getMaxFeePrice(maxFeePriceInput, cfg.PriceMax()) // Make a wrapper config
	bumpedFeePrice := bumpFeePriceByPercentage(originalfeePrice, cfg.BumpPercent(), cfg.BumpMin())

	// Update bumpedFeePrice if currentfeePrice is higher than bumpedFeePrice and within maxFeePrice
	bumpedFeePrice = maxBumpedFee(lggr, currentfeePrice, bumpedFeePrice, maxFeePrice, "fee price")

	if bumpedFeePrice.Cmp(maxFeePrice) > 0 {
		return maxFeePrice, errors.Wrapf(ErrBumpFeeExceedsLimit, "bumped fee price of %s would exceed configured max fee price of %s (original price was %s). %s",
			bumpedFeePrice.String(), maxFeePrice, originalfeePrice.String(), label.NodeConnectivityProblemWarning)
	} else if bumpedFeePrice.Cmp(originalfeePrice) == 0 {
		// NOTE: This really shouldn't happen since we enforce minimums for
		// FeeEstimator.BumpPercent and FeeEstimator.BumpMin in the config validation,
		// but it's here anyway for a "belts and braces" approach
		return bumpedFeePrice, errors.Wrapf(ErrBump, "bumped fee price of %s is equal to original fee price of %s."+
			" ACTION REQUIRED: This is a configuration error, you must increase either "+
			"FeeEstimator.BumpPercent or FeeEstimator.BumpMin", bumpedFeePrice.String(), originalfeePrice.String())
	}
	return bumpedFeePrice, nil
}

func getMaxFeePrice(userSpecifiedMax, maxFeePrice *big.Int) *big.Int {
	return FeePriceLimit(userSpecifiedMax, maxFeePrice)
}

func bumpFeePriceByPercentage(originalFeePrice *big.Int, feeBumpPercent uint16, feeBumpUnits *big.Int) *big.Int {
	return max(new(big.Int).Add(originalFeePrice, feeBumpUnits), addPercentage(originalFeePrice, feeBumpPercent))
}

func maxBumpedFee(lggr logger.SugaredLogger, currentFeePrice, bumpedFeePrice, maxFeePrice *big.Int, feeType string) *big.Int {
	if currentFeePrice != nil {
		if currentFeePrice.Cmp(maxFeePrice) > 0 {
			// Shouldn't happen because the estimator should not be allowed to
			// estimate a higher fee than the maximum allowed
			lggr.AssumptionViolationf("Ignoring current %s of %s that would exceed max %s of %s", feeType, currentFeePrice.String(), feeType, maxFeePrice.String())
		} else if bumpedFeePrice.Cmp(currentFeePrice) < 0 {
			// If the current fee price is higher than the old price bumped, use that instead
			bumpedFeePrice = currentFeePrice
		}
	}
	return bumpedFeePrice
}

// CalculateFee computes the fee price and chain specific fee limit for a transaction.
func CalculateFee(cfg feetypes.FixedPriceEstimatorConfig, bumpCfg feetypes.BumpConfig, feeLimit uint32, maxFeePrice *big.Int) (feePrice *big.Int, chainSpecificFeeLimit uint32, err error) {
	feePrice = cfg.PriceDefault()
	feePrice, chainSpecificFeeLimit = CapFeePrice(feePrice, maxFeePrice, bumpCfg.PriceMax(), feeLimit, bumpCfg.LimitMultiplier())
	return
}

func GetDynamicFee(cfg feetypes.FixedPriceEstimatorConfig, originalFeeLimit uint32, maxFeePrice *big.Int) (feeCap, tipCap *big.Int, chainSpecificFeeLimit uint32, err error) {
	tipCap = cfg.TipCapDefault()

	if tipCap == nil {
		return big.NewInt(0), big.NewInt(0), 0, errors.New("cannot calculate dynamic fee: FeeTipCapDefault was not set")
	}

	chainSpecificFeeLimit = ApplyMultiplier(originalFeeLimit, cfg.LimitMultiplier())
	feeCap = GetFeeCap(cfg, originalFeeLimit, maxFeePrice)

	return feeCap, tipCap, chainSpecificFeeLimit, nil

}

func GetFeeCap(cfg feetypes.FixedPriceEstimatorConfig, originalFeeLimit uint32, maxFeePrice *big.Int) (feeCap *big.Int) {
	if cfg.BumpThreshold() == 0 {
		// Fee bumping is disabled, just use the max fee cap
		feeCap = getMaxFeePrice(maxFeePrice, cfg.PriceMax())
	} else {
		// Need to leave headroom for bumping so we fallback to the default value here
		feeCap = cfg.FeeCapDefault()
	}
	return feeCap
}

// BumpDynamicFeeOnly bumps the tip cap and max fee price if necessary
func CalculateBumpDynamicFee(cfg feetypes.BumpConfig, feeCapBufferBlocks uint16, lggr logger.SugaredLogger, currentTipCap, currentBaseFee *big.Int, originalFeeCap, originalTipCap *big.Int, originalFeeLimit uint32, maxFeePrice *big.Int) (bumpedFeeCap, bumpedTipCap *big.Int, chainSpecificFeeLimit uint32, err error) {
	bumpedFeeCap, bumpedTipCap, err = bumpDynamicFee(cfg, feeCapBufferBlocks, lggr, currentTipCap, currentBaseFee, originalFeeCap, originalTipCap, maxFeePrice)
	if err != nil {
		return bumpedFeeCap, bumpedTipCap, 0, err
	}
	chainSpecificFeeLimit = ApplyMultiplier(originalFeeLimit, cfg.LimitMultiplier())
	return
}

func bumpDynamicFee(cfg feetypes.BumpConfig, feeCapBufferBlocks uint16, lggr logger.SugaredLogger, currentTipCap, currentBaseFee *big.Int, originalFeeCap, originalTipCap *big.Int, maxFeePriceInput *big.Int) (bumpedFeeCap, bumpedTipCap *big.Int, err error) {
	maxFeePrice := getMaxFeePrice(maxFeePriceInput, cfg.PriceMax())
	baselineTipCap := max(originalTipCap, cfg.TipCapDefault())
	bumpedTipCap = bumpFeePriceByPercentage(baselineTipCap, cfg.BumpPercent(), cfg.BumpMin())

	// Update bumpedTipCap if currentTipCap is higher than bumpedTipCap and within maxFeePrice
	bumpedTipCap = maxBumpedFee(lggr, currentTipCap, bumpedTipCap, maxFeePrice, "tip cap")

	if bumpedTipCap.Cmp(maxFeePrice) > 0 {
		return bumpedFeeCap, bumpedTipCap, errors.Wrapf(ErrBumpFeeExceedsLimit, "bumped tip cap of %s would exceed configured max fee price of %s (original fee: tip cap %s, fee cap %s). %s",
			bumpedTipCap.String(), maxFeePrice, originalTipCap.String(), originalFeeCap.String(), label.NodeConnectivityProblemWarning)
	} else if bumpedTipCap.Cmp(originalTipCap) <= 0 {
		// NOTE: This really shouldn't happen since we enforce minimums for
		// FeeEstimator.BumpPercent and FeeEstimator.BumpMin in the config validation,
		// It's here for extra precaution
		return bumpedFeeCap, bumpedTipCap, errors.Wrapf(ErrBump, "bumped fee tip cap of %s is less than or equal to original fee tip cap of %s."+
			" ACTION REQUIRED: This is a configuration error, you must increase either "+
			"FeeEstimator.BumpPercent or FeeEstimator.BumpMin", bumpedTipCap.String(), originalTipCap)
	}

	// Always bump the FeeCap by at least the bump percentage
	// For geth the configured bump is 10%
	bumpedFeeCap = max(
		addPercentage(originalFeeCap, cfg.BumpPercent()),
		new(big.Int).Add(originalFeeCap, cfg.BumpMin()),
	)

	if currentBaseFee != nil {
		if currentBaseFee.Cmp(maxFeePrice) > 0 {
			lggr.Warnf("Ignoring current base fee of %s which is greater than max fee price of %s", currentBaseFee.String(), maxFeePrice.String())
		} else {
			currentFeeCap := calcFeeCap(currentBaseFee, int(feeCapBufferBlocks), bumpedTipCap, maxFeePrice)
			bumpedFeeCap = max(bumpedFeeCap, currentFeeCap)
		}
	}

	if bumpedFeeCap.Cmp(maxFeePrice) > 0 {
		return bumpedFeeCap, bumpedTipCap, errors.Wrapf(ErrBumpFeeExceedsLimit, "bumped fee cap of %s would exceed configured max fee price of %s (original fee: tip cap %s, fee cap %s). %s",
			bumpedFeeCap.String(), maxFeePrice, originalTipCap.String(), originalFeeCap.String(), label.NodeConnectivityProblemWarning)
	}
	return bumpedFeeCap, bumpedTipCap, nil
}

func calcFeeCap(latestAvailableBaseFeePerUnit *big.Int, bufferBlocks int, tipCap *big.Int, maxFeePrice *big.Int) (feeCap *big.Int) {
	const maxBaseFeeIncreasePerBlock float64 = 1.125

	baseFee := new(big.Float)
	baseFee.SetInt(latestAvailableBaseFeePerUnit)
	// Find out the worst case base fee before we bump
	multiplier := big.NewFloat(maxBaseFeeIncreasePerBlock)
	for i := 0; i < bufferBlocks; i++ {
		baseFee.Mul(baseFee, multiplier)
	}

	baseFeeInt, _ := baseFee.Int(nil)
	feeCap = baseFeeInt.Add(baseFeeInt, tipCap)

	if feeCap.Cmp(maxFeePrice) > 0 {
		return maxFeePrice
	}
	return feeCap
}
