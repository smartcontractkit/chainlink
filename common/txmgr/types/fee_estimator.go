package gas

import (
	"context"

	"github.com/pkg/errors"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var _ EvmEstimator = &fixedPriceEstimator{}

type fixedPriceEstimator struct {
	config Config
	lggr   logger.SugaredLogger
}

// NewFixedPriceEstimator returns a new "FixedPrice" estimator which will
// always use the config default values for gas prices and limits
func NewFixedPriceEstimator(cfg Config, lggr logger.Logger) Estimator {
	return &fixedPriceEstimator{cfg, logger.Sugared(lggr.Named("FixedPriceEstimator"))}
}

func (f *fixedPriceEstimator) Start(context.Context) error {
	if f.config.EvmGasBumpThreshold() == 0 && f.config.GasEstimatorMode() == "FixedPrice" {
		// EvmGasFeeCapDefault is ignored if fixed estimator mode is on and gas bumping is disabled
		if f.config.EvmGasFeeCapDefault().Cmp(f.config.EvmMaxGasPriceWei()) != 0 {
			f.lggr.Infof("You are using FixedPrice estimator with gas bumping disabled. EVM.GasEstimator.PriceMax (value: %s) will be used as the FeeCap for transactions", f.config.EvmMaxGasPriceWei())
		}
	}
	return nil
}
func (f *fixedPriceEstimator) Name() string                                          { return f.lggr.Name() }
func (f *fixedPriceEstimator) Ready() error                                          { return nil }
func (f *fixedPriceEstimator) HealthReport() map[string]error                        { return map[string]error{} }
func (f *fixedPriceEstimator) Close() error                                          { return nil }
func (f *fixedPriceEstimator) OnNewLongestChain(_ context.Context, _ *evmtypes.Head) {}

func (f *fixedPriceEstimator) BumpFee(_ context.Context, originalFee txmgrtypes.Fee, originalGasLimit uint32, maxGasPriceWei txmgrtypes.Fee, _ []EvmPriorAttempt, feeType txmgrtypes.FeeType) (fee txmgrtypes.Fee, chainSpecificGasLimit uint32, err error) {
	if feeType == txmgrtypes.EvmLegacy {
		return BumpLegacyGasPriceOnly(f.config, f.lggr, f.config.EvmGasPriceDefault(), originalFee.(*assets.Wei), originalGasLimit, maxGasPriceWei.(*assets.Wei))
	} else if feeType == txmgrtypes.EvmDynamic {
		return BumpDynamicFeeOnly(f.config, f.lggr, f.config.EvmGasTipCapDefault(), nil, originalFee.(DynamicFee), originalGasLimit, maxGasPriceWei.(*assets.Wei))
	}
	return nil, 0, errors.Errorf("unknown fee type %v", feeType)
}

func (f *fixedPriceEstimator) BumpLegacyGas(_ context.Context, originalGasPrice *assets.Wei, originalGasLimit uint32, maxGasPriceWei *assets.Wei, _ []EvmPriorAttempt) (gasPrice *assets.Wei, chainSpecificGasLimit uint32, err error) {
	return BumpLegacyGasPriceOnly(f.config, f.lggr, f.config.EvmGasPriceDefault(), originalGasPrice, originalGasLimit, maxGasPriceWei)
}

func (f *fixedPriceEstimator) BumpDynamicFee(_ context.Context, originalFee DynamicFee, originalGasLimit uint32, maxGasPriceWei *assets.Wei, _ []EvmPriorAttempt) (bumped DynamicFee, chainSpecificGasLimit uint32, err error) {
	return BumpDynamicFeeOnly(f.config, f.lggr, f.config.EvmGasTipCapDefault(), nil, originalFee, originalGasLimit, maxGasPriceWei)
}

// Chain Agnostic Gas Estimator to get fee based on fee type
func (f *fixedPriceEstimator) GetFee(_ context.Context, _ []byte, gasLimit uint32, maxGasPriceWei txmgrtypes.Fee, feeType txmgrtypes.FeeType, _ ...txmgrtypes.Opt) (fee txmgrtypes.Fee, chainSpecificGasLimit uint32, err error) {
	if feeType == txmgrtypes.EvmLegacy {
		return f.GetLegacyGas(context.Background(), nil, gasLimit, maxGasPriceWei.(*assets.Wei))
	} else if feeType == txmgrtypes.EvmDynamic {
		dynamicFee, chainSpecificGasLimit, err := f.GetDynamicFee(context.Background(), gasLimit, maxGasPriceWei.(*assets.Wei))
		return dynamicFee, chainSpecificGasLimit, err
	}
	return
}

func (f *fixedPriceEstimator) GetLegacyGas(_ context.Context, _ []byte, gasLimit uint32, maxGasPriceWei *assets.Wei, _ ...txmgrtypes.Opt) (feeCap *assets.Wei, chainSpecificGasLimit uint32, err error) {
	chainSpecificGasLimit = applyMultiplier(gasLimit, f.config.EvmGasLimitMultiplier())
	feeCap, err = getFeeCap(maxGasPriceWei, f.config, txmgrtypes.EvmLegacy)

	return
}

func (f *fixedPriceEstimator) GetDynamicFee(_ context.Context, originalGasLimit uint32, maxGasPriceWei *assets.Wei) (d DynamicFee, chainSpecificGasLimit uint32, err error) {
	gasTipCap := f.config.EvmGasTipCapDefault()
	if gasTipCap == nil {
		return d, 0, errors.New("cannot calculate dynamic fee: EthGasTipCapDefault was not set")
	}
	chainSpecificGasLimit = applyMultiplier(originalGasLimit, f.config.EvmGasLimitMultiplier())

	feeCap, err := getFeeCap(maxGasPriceWei, f.config, txmgrtypes.EvmDynamic)

	return DynamicFee{
		FeeCap: feeCap,
		TipCap: gasTipCap,
	}, chainSpecificGasLimit, err
}

// Returns fee cap based on fee type
func getFeeCap(maxGasPriceWei *assets.Wei, cfg Config, feeType txmgrtypes.FeeType) (*assets.Wei, error) {
	if feeType == txmgrtypes.EvmLegacy {
		return getLegacyFeeCap(maxGasPriceWei, cfg), nil
	} else if feeType == txmgrtypes.EvmDynamic {
		return getDynamicFeeCap(maxGasPriceWei, cfg), nil
	}

	return nil, errors.Errorf("unknown fee type %v", feeType)
}

func getLegacyFeeCap(maxGasPriceWei *assets.Wei, cfg Config) *assets.Wei {
	return capGasPrice(cfg.EvmGasPriceDefault(), maxGasPriceWei, cfg)
}

func getDynamicFeeCap(maxGasPriceWei *assets.Wei, cfg Config) *assets.Wei {
	if cfg.EvmGasBumpThreshold() == 0 {
		// Gas bumping is disabled, just use the max fee cap
		return getMaxGasPrice(maxGasPriceWei, cfg)
	} else {
		// Need to leave headroom for bumping so we fallback to the default value here
		return cfg.EvmGasFeeCapDefault()
	}
}
