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

// LegacyGas Price
// TODO: What does legacy gas price mean? define that
// BumpLegacyGasPriceOnly will increase the price and apply multiplier to the gas limit
func BumpLegacyGasPriceOnlyBigInt(cfg feetypes.BumpConfig, lggr logger.SugaredLogger, currentGasPrice, originalGasPrice *big.Int, originalGasLimit uint32, maxGasPriceWei *big.Int) (gasPrice *big.Int, chainSpecificGasLimit uint32, err error) {
	gasPrice, err = bumpGasPriceBigInt(cfg, lggr, currentGasPrice, originalGasPrice, maxGasPriceWei)
	if err != nil {
		return nil, 0, err
	}
	chainSpecificGasLimit = ApplyMultiplier(originalGasLimit, cfg.LimitMultiplier())
	return
}

// bumpGasPrice computes the next gas price to attempt as the largest of:
// - A configured percentage bump (EVM.GasEstimator.BumpPercent) on top of the baseline price.
// - A configured fixed amount of Wei (ETH_GAS_PRICE_WEI) on top of the baseline price.
// The baseline price is the maximum of the previous gas price attempt and the node's current gas price.
func bumpGasPriceBigInt(cfg feetypes.BumpConfig, lggr logger.SugaredLogger, currentGasPrice, originalGasPrice, maxGasPriceWei *big.Int) (*big.Int, error) {
	maxGasPrice := getMaxGasPriceBigInt(maxGasPriceWei, cfg.PriceMax()) // Make a wrapper config
	bumpedGasPrice := bumpFeePriceBigInt(originalGasPrice, cfg.BumpPercent(), cfg.BumpMin())

	// Update bumpedGasPrice if currentGasPrice is higher than bumpedGasPrice and within maxGasPrice
	bumpedGasPrice = maxBumpedFeeBigInt(lggr, currentGasPrice, bumpedGasPrice, maxGasPrice, "gas price")

	if bumpedGasPrice.Cmp(maxGasPrice) > 0 {
		return maxGasPrice, errors.Wrapf(ErrBumpGasExceedsLimit, "bumped gas price of %s would exceed configured max gas price of %s (original price was %s). %s",
			bumpedGasPrice.String(), maxGasPrice, originalGasPrice.String(), label.NodeConnectivityProblemWarning)
	} else if bumpedGasPrice.Cmp(originalGasPrice) == 0 {
		// NOTE: This really shouldn't happen since we enforce minimums for
		// EVM.GasEstimator.BumpPercent and EVM.GasEstimator.BumpMin in the config validation,
		// but it's here anyway for a "belts and braces" approach
		return bumpedGasPrice, errors.Wrapf(ErrBump, "bumped gas price of %s is equal to original gas price of %s."+
			" ACTION REQUIRED: This is a configuration error, you must increase either "+
			"EVM.GasEstimator.BumpPercent or EVM.GasEstimator.BumpMin", bumpedGasPrice.String(), originalGasPrice.String())
	}
	return bumpedGasPrice, nil
}

func getMaxGasPriceBigInt(userSpecifiedMax, maxGasPriceWei *big.Int) *big.Int {
	return GetCeilingFeePrice(userSpecifiedMax, maxGasPriceWei)
}

func bumpFeePriceBigInt(originalFeePrice *big.Int, feeBumpPercent uint16, feeBumpUnits *big.Int) *big.Int {

	linearFeePrice := new(big.Int)
	linearFeePrice.Add(originalFeePrice, feeBumpUnits)
	percentageFeePrice := AddPercentage(originalFeePrice, feeBumpPercent)
	// Find which is higher using max
	return Max(linearFeePrice, percentageFeePrice)
}

func maxBumpedFeeBigInt(lggr logger.SugaredLogger, currentFeePrice, bumpedFeePrice, maxGasPrice *big.Int, feeType string) *big.Int {
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

// TODO: Refactor this into another package
// Max returns the maximum of the two given values.
func Max(x, y *big.Int) *big.Int {
	if x.Cmp(y) == 1 {
		return x
	}
	return y
}

// TODO: Refactor this into another package
func AddPercentage(value *big.Int, percentage uint16) *big.Int {
	bumped := new(big.Int)
	bumped.Mul(value, big.NewInt(int64(100+percentage)))
	bumped.Div(bumped, big.NewInt(100))
	return bumped
}

// GetLegacyGas computes the gas price and chain specific gas limit for a transaction.
func GetLegacyGas(cfg feetypes.FixedPriceEstimatorConfig, bumpCfg feetypes.BumpConfig, gasLimit uint32, maxGasPriceUnit *big.Int) (gasPrice *big.Int, chainSpecificGasLimit uint32, err error) {
	gasPrice = cfg.PriceDefault()
	gasPrice, chainSpecificGasLimit = CapFeePrice(gasPrice, maxGasPriceUnit, bumpCfg.PriceMax(), gasLimit, bumpCfg.LimitMultiplier())
	return
}

// Dynamic Fee

func GetDynamicFee(cfg feetypes.FixedPriceEstimatorConfig, originalGasLimit uint32, maxGasPriceWei *big.Int) (feeCap, tipCap *big.Int, chainSpecificGasLimit uint32, err error) {
	tipCap = cfg.TipCapDefault()

	if tipCap == nil {
		return big.NewInt(0), big.NewInt(0), 0, errors.New("cannot calculate dynamic fee: EthGasTipCapDefault was not set")
	}

	chainSpecificGasLimit = ApplyMultiplier(originalGasLimit, cfg.LimitMultiplier())
	feeCap = GetFeeCap(cfg, originalGasLimit, maxGasPriceWei)

	return feeCap, tipCap, chainSpecificGasLimit, nil

}

func GetFeeCap(cfg feetypes.FixedPriceEstimatorConfig, originalGasLimit uint32, maxGasPriceWei *big.Int) (feeCap *big.Int) {
	if cfg.BumpThreshold() == 0 {
		// Gas bumping is disabled, just use the max fee cap
		feeCap = getMaxGasPriceBigInt(maxGasPriceWei, cfg.PriceMax())
	} else {
		// Need to leave headroom for bumping so we fallback to the default value here
		feeCap = cfg.FeeCapDefault()
	}
	return feeCap
}

// // BumpDynamicFeeOnly bumps the tip cap and max gas price if necessary
// func BumpDynamicFeeOnly(config bumpConfig, feeCapBufferBlocks uint16, lggr logger.SugaredLogger, currentTipCap, currentBaseFee *assets.Wei, originalFee DynamicFee, originalGasLimit uint32, maxGasPriceWei *assets.Wei) (bumped DynamicFee, chainSpecificGasLimit uint32, err error) {
// 	bumped, err = bumpDynamicFee(config, feeCapBufferBlocks, lggr, currentTipCap, currentBaseFee, originalFee, maxGasPriceWei)
// 	if err != nil {
// 		return bumped, 0, err
// 	}
// 	chainSpecificGasLimit = ApplyMultiplier(originalGasLimit, config.LimitMultiplier())
// 	return
// }
