package gas

import (
	"context"

	commonfee "github.com/smartcontractkit/chainlink/v2/common/fee"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var _ EvmEstimator = (*fixedPriceEstimator)(nil)

type fixedPriceEstimator struct {
	config     fixedPriceEstimatorConfigWrapper
	bumpConfig wrappedBumpConfig
	bhConfig   fixedPriceEstimatorBlockHistoryConfig
	lggr       logger.SugaredLogger
}

type fixedPriceEstimatorBlockHistoryConfig interface {
	EIP1559FeeCapBufferBlocks() uint16
}

// NewFixedPriceEstimator returns a new "FixedPrice" estimator which will
// always use the config default values for gas prices and limits
func NewFixedPriceEstimator(config fixedPriceEstimatorConfigWrapper, bumpCfg wrappedBumpConfig, bhCfg fixedPriceEstimatorBlockHistoryConfig, lggr logger.Logger) EvmEstimator {
	return &fixedPriceEstimator{config, bumpCfg, bhCfg, logger.Sugared(lggr.Named("FixedPriceEstimator"))}
}

func (f *fixedPriceEstimator) Start(context.Context) error {
	if f.config.BumpThreshold() == 0 && f.config.Mode() == "FixedPrice" {
		// EvmGasFeeCapDefault is ignored if fixed estimator mode is on and gas bumping is disabled
		if f.config.FeeCapDefault().Cmp(f.config.PriceMax()) != 0 {
			f.lggr.Infof("You are using FixedPrice estimator with gas bumping disabled. EVM.GasEstimator.PriceMax (value: %s) will be used as the FeeCap for transactions", f.config.PriceMax())
		}
	}
	return nil
}

func (f *fixedPriceEstimator) Name() string                                          { return f.lggr.Name() }
func (f *fixedPriceEstimator) Ready() error                                          { return nil }
func (f *fixedPriceEstimator) HealthReport() map[string]error                        { return map[string]error{} }
func (f *fixedPriceEstimator) Close() error                                          { return nil }
func (f *fixedPriceEstimator) OnNewLongestChain(_ context.Context, _ *evmtypes.Head) {}

func (f *fixedPriceEstimator) GetLegacyGas(_ context.Context, _ []byte, gasLimit uint32, maxGasPriceWei *assets.Wei, _ ...feetypes.Opt) (*assets.Wei, uint32, error) {
	gasPrice, chainSpecificGasLimit, err := commonfee.GetLegacyGas(f.config, f.bumpConfig, gasLimit, maxGasPriceWei.ToInt())
	return assets.NewWei(gasPrice), chainSpecificGasLimit, err
}

func (f *fixedPriceEstimator) BumpLegacyGas(_ context.Context, originalGasPrice *assets.Wei, originalGasLimit uint32, maxGasPriceWei *assets.Wei, _ []EvmPriorAttempt) (*assets.Wei, uint32, error) {
	gasPrice, chainSpecificGasLimit, err := commonfee.BumpLegacyGasPriceOnly(f.bumpConfig, f.lggr, f.config.PriceDefault(), originalGasPrice.ToInt(), originalGasLimit, maxGasPriceWei.ToInt())
	return assets.NewWei(gasPrice), chainSpecificGasLimit, err
}

func (f *fixedPriceEstimator) GetDynamicFee(_ context.Context, originalGasLimit uint32, maxGasPriceWei *assets.Wei) (d DynamicFee, chainSpecificGasLimit uint32, err error) {
	feeCap, tipCap, chainSpecificGasLimit, err := commonfee.GetDynamicFee(f.config, originalGasLimit, maxGasPriceWei.ToInt())
	if err != nil {
		return d, 0, err
	}

	return DynamicFee{
		FeeCap: assets.NewWei(feeCap),
		TipCap: assets.NewWei(tipCap),
	}, chainSpecificGasLimit, nil
}

func (f *fixedPriceEstimator) BumpDynamicFee(_ context.Context, originalFee DynamicFee, originalGasLimit uint32, maxGasPriceWei *assets.Wei, _ []EvmPriorAttempt) (bumped DynamicFee, chainSpecificGasLimit uint32, err error) {
	return BumpDynamicFeeOnly(f.bumpConfig.config, f.bhConfig.EIP1559FeeCapBufferBlocks(), f.lggr, f.config.config.TipCapDefault(), nil, originalFee, originalGasLimit, maxGasPriceWei)
}
