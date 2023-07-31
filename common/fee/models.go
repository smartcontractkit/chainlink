package fee

import (
	"math/big"

	"github.com/pkg/errors"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/label" // TODO: Remove import from core
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var (
	ErrBumpGasExceedsLimit = errors.New("gas bump exceeds limit")
	ErrBump                = errors.New("gas bump failed")
	ErrConnectivity        = errors.New("transaction propagation issue: transactions are not being mined")
)

func IsBumpErr(err error) bool {
	return err != nil && (errors.Is(err, ErrBumpGasExceedsLimit) || errors.Is(err, ErrBump) || errors.Is(err, ErrConnectivity))
}

// TODO: What does legacy gas price mean? define that
// BumpLegacyGasPriceOnly will increase the price and apply multiplier to the gas limit
func BumpLegacyGasPriceOnly(cfg feetypes.BumpConfig, lggr logger.SugaredLogger, currentGasPrice, originalGasPrice *big.Int, originalGasLimit uint32, maxGasPrice *big.Int) (gasPrice *big.Int, chainSpecificGasLimit uint32, err error) {
	gasPrice, err = bumpGasPrice(cfg, lggr, currentGasPrice, originalGasPrice, maxGasPrice)
	if err != nil {
		return nil, 0, err
	}
	chainSpecificGasLimit = ApplyMultiplier(originalGasLimit, cfg.LimitMultiplier())
	return
}

// bumpGasPrice computes the next gas price to attempt as the largest of:
// - A configured percentage bump (GasEstimator.BumpPercent) on top of the baseline price.
// - A configured fixed amount of Unit (ETH_GAS_PRICE_Unit) on top of the baseline price.
// The baseline price is the maximum of the previous gas price attempt and the node's current gas price.
func bumpGasPrice(cfg feetypes.BumpConfig, lggr logger.SugaredLogger, currentGasPrice, originalGasPrice, maxGasPriceInput *big.Int) (*big.Int, error) {
	maxGasPrice := getMaxGasPrice(maxGasPriceInput, cfg.PriceMax()) // Make a wrapper config
	bumpedGasPrice := bumpFeePrice(originalGasPrice, cfg.BumpPercent(), cfg.BumpMin())

	// Update bumpedGasPrice if currentGasPrice is higher than bumpedGasPrice and within maxGasPrice
	bumpedGasPrice = maxBumpedFee(lggr, currentGasPrice, bumpedGasPrice, maxGasPrice, "gas price")

	if bumpedGasPrice.Cmp(maxGasPrice) > 0 {
		return maxGasPrice, errors.Wrapf(ErrBumpGasExceedsLimit, "bumped gas price of %s would exceed configured max gas price of %s (original price was %s). %s",
			bumpedGasPrice.String(), maxGasPrice, originalGasPrice.String(), label.NodeConnectivityProblemWarning)
	} else if bumpedGasPrice.Cmp(originalGasPrice) == 0 {
		// NOTE: This really shouldn't happen since we enforce minimums for
		// GasEstimator.BumpPercent and GasEstimator.BumpMin in the config validation,
		// but it's here anyway for a "belts and braces" approach
		return bumpedGasPrice, errors.Wrapf(ErrBump, "bumped gas price of %s is equal to original gas price of %s."+
			" ACTION REQUIRED: This is a configuration error, you must increase either "+
			"GasEstimator.BumpPercent or GasEstimator.BumpMin", bumpedGasPrice.String(), originalGasPrice.String())
	}
	return bumpedGasPrice, nil
}

func getMaxGasPrice(userSpecifiedMax, maxGasPrice *big.Int) *big.Int {
	return FeePriceLimit(userSpecifiedMax, maxGasPrice)
}

func bumpFeePrice(originalFeePrice *big.Int, feeBumpPercent uint16, feeBumpUnits *big.Int) *big.Int {
	linearFeePrice := new(big.Int)
	linearFeePrice.Add(originalFeePrice, feeBumpUnits)
	percentageFeePrice := AddPercentage(originalFeePrice, feeBumpPercent)
	return Max(linearFeePrice, percentageFeePrice)
}

func maxBumpedFee(lggr logger.SugaredLogger, currentFeePrice, bumpedFeePrice, maxGasPrice *big.Int, feeType string) *big.Int {
	if currentFeePrice != nil {
		if currentFeePrice.Cmp(maxGasPrice) > 0 {
			// Shouldn't happen because the estimator should not be allowed to
			// estimate a higher gas than the maximum allowed
			lggr.AssumptionViolationf("Ignoring current %s of %s that would exceed max %s of %s", feeType, currentFeePrice.String(), feeType, maxGasPrice.String())
		} else if bumpedFeePrice.Cmp(currentFeePrice) < 0 {
			// If the current gas price is higher than the old price bumped, use that instead
			bumpedFeePrice = currentFeePrice
		}
	}
	return bumpedFeePrice
}

// GetLegacyGas computes the gas price and chain specific gas limit for a transaction.
func GetLegacyGas(cfg feetypes.FixedPriceEstimatorConfig, bumpCfg feetypes.BumpConfig, gasLimit uint32, maxGasPrice *big.Int) (gasPrice *big.Int, chainSpecificGasLimit uint32, err error) {
	gasPrice = cfg.PriceDefault()
	gasPrice, chainSpecificGasLimit = CapFeePrice(gasPrice, maxGasPrice, bumpCfg.PriceMax(), gasLimit, bumpCfg.LimitMultiplier())
	return
}

func GetDynamicFee(cfg feetypes.FixedPriceEstimatorConfig, originalGasLimit uint32, maxGasPrice *big.Int) (feeCap, tipCap *big.Int, chainSpecificGasLimit uint32, err error) {
	tipCap = cfg.TipCapDefault()

	if tipCap == nil {
		return big.NewInt(0), big.NewInt(0), 0, errors.New("cannot calculate dynamic fee: EthGasTipCapDefault was not set")
	}

	chainSpecificGasLimit = ApplyMultiplier(originalGasLimit, cfg.LimitMultiplier())
	feeCap = GetFeeCap(cfg, originalGasLimit, maxGasPrice)

	return feeCap, tipCap, chainSpecificGasLimit, nil

}

func GetFeeCap(cfg feetypes.FixedPriceEstimatorConfig, originalGasLimit uint32, maxGasPrice *big.Int) (feeCap *big.Int) {
	if cfg.BumpThreshold() == 0 {
		// Gas bumping is disabled, just use the max fee cap
		feeCap = getMaxGasPrice(maxGasPrice, cfg.PriceMax())
	} else {
		// Need to leave headroom for bumping so we fallback to the default value here
		feeCap = cfg.FeeCapDefault()
	}
	return feeCap
}

// BumpDynamicFeeOnly bumps the tip cap and max gas price if necessary
func BumpDynamicFeeOnly(cfg feetypes.BumpConfig, feeCapBufferBlocks uint16, lggr logger.SugaredLogger, currentTipCap, currentBaseFee *big.Int, originalFeeCap, originalTipCap *big.Int, originalGasLimit uint32, maxGasPrice *big.Int) (bumpedFeeCap, bumpedTipCap *big.Int, chainSpecificGasLimit uint32, err error) {
	bumpedFeeCap, bumpedTipCap, err = bumpDynamicFee(cfg, feeCapBufferBlocks, lggr, currentTipCap, currentBaseFee, originalFeeCap, originalTipCap, maxGasPrice)
	if err != nil {
		return bumpedFeeCap, bumpedTipCap, 0, err
	}
	chainSpecificGasLimit = ApplyMultiplier(originalGasLimit, cfg.LimitMultiplier())
	return
}

func bumpDynamicFee(cfg feetypes.BumpConfig, feeCapBufferBlocks uint16, lggr logger.SugaredLogger, currentTipCap, currentBaseFee *big.Int, originalFeeCap, originalTipCap *big.Int, maxGasPriceInput *big.Int) (bumpedFeeCap, bumpedTipCap *big.Int, err error) {
	maxGasPrice := getMaxGasPrice(maxGasPriceInput, cfg.PriceMax()) // TODO: Rename gas to fee
	baselineTipCap := Max(originalTipCap, cfg.TipCapDefault())
	bumpedTipCap = bumpFeePrice(baselineTipCap, cfg.BumpPercent(), cfg.BumpMin())

	// Update bumpedTipCap if currentTipCap is higher than bumpedTipCap and within maxGasPrice
	bumpedTipCap = maxBumpedFee(lggr, currentTipCap, bumpedTipCap, maxGasPrice, "tip cap")

	if bumpedTipCap.Cmp(maxGasPrice) > 0 {
		return bumpedFeeCap, bumpedTipCap, errors.Wrapf(ErrBumpGasExceedsLimit, "bumped tip cap of %s would exceed configured max gas price of %s (original fee: tip cap %s, fee cap %s). %s",
			bumpedTipCap.String(), maxGasPrice, originalTipCap.String(), originalFeeCap.String(), label.NodeConnectivityProblemWarning)
	} else if bumpedTipCap.Cmp(originalTipCap) <= 0 {
		// NOTE: This really shouldn't happen since we enforce minimums for
		// GasEstimator.BumpPercent and GasEstimator.BumpMin in the config validation,
		// but it's here anyway for a "belts and braces" approach
		return bumpedFeeCap, bumpedTipCap, errors.Wrapf(ErrBump, "bumped gas tip cap of %s is less than or equal to original gas tip cap of %s."+
			" ACTION REQUIRED: This is a configuration error, you must increase either "+
			"GasEstimator.BumpPercent or GasEstimator.BumpMin", bumpedTipCap.String(), originalTipCap)
	}

	// Always bump the FeeCap by at least the bump percentage (should be greater than or
	// equal to than geth's configured bump minimum which is 10%)
	// See: https://github.com/ethereum/go-ethereum/blob/bff330335b94af3643ac2fb809793f77de3069d4/core/tx_list.go#L298
	// TODO: is this generalisable, looking at the comments above?
	// TODO: check for correctness
	bumpedFeeCap = Max(
		AddPercentage(originalFeeCap, cfg.BumpPercent()),
		new(big.Int).Add(originalFeeCap, cfg.BumpMin()),
	)

	if currentBaseFee != nil {
		if currentBaseFee.Cmp(maxGasPrice) > 0 {
			lggr.Warnf("Ignoring current base fee of %s which is greater than max gas price of %s", currentBaseFee.String(), maxGasPrice.String())
		} else {
			currentFeeCap := calcFeeCap(currentBaseFee, int(feeCapBufferBlocks), bumpedTipCap, maxGasPrice)
			bumpedFeeCap = Max(bumpedFeeCap, currentFeeCap)
		}
	}

	if bumpedFeeCap.Cmp(maxGasPrice) > 0 {
		return bumpedFeeCap, bumpedTipCap, errors.Wrapf(ErrBumpGasExceedsLimit, "bumped fee cap of %s would exceed configured max gas price of %s (original fee: tip cap %s, fee cap %s). %s",
			bumpedFeeCap.String(), maxGasPrice, originalTipCap.String(), originalFeeCap.String(), label.NodeConnectivityProblemWarning)
	}
	return bumpedFeeCap, bumpedTipCap, nil
}

func calcFeeCap(latestAvailableBaseFeePerGas *big.Int, bufferBlocks int, tipCap *big.Int, maxGasPrice *big.Int) (feeCap *big.Int) {
	const maxBaseFeeIncreasePerBlock float64 = 1.125 // Todo: generalise this?

	baseFee := new(big.Float)
	baseFee.SetInt(latestAvailableBaseFeePerGas)
	// Find out the worst case base fee before we should bump
	multiplier := big.NewFloat(maxBaseFeeIncreasePerBlock)
	for i := 0; i < bufferBlocks; i++ {
		baseFee.Mul(baseFee, multiplier)
	}

	baseFeeInt, _ := baseFee.Int(nil)
	feeCap = baseFeeInt.Add(baseFeeInt, tipCap)

	if feeCap.Cmp(maxGasPrice) > 0 {
		return maxGasPrice
	}
	return feeCap
}
