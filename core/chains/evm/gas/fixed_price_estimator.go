package gas

import (
	"context"

	"github.com/pkg/errors"

	commonfee "github.com/smartcontractkit/chainlink/v2/common/fee"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var _ EvmEstimator = (*fixedPriceEstimator)(nil)

type fixedPriceEstimator struct {
	config   fixedPriceEstimatorConfig
	bhConfig fixedPriceEstimatorBlockHistoryConfig
	lggr     logger.SugaredLogger
}

type fixedPriceEstimatorConfig interface {
	BumpThreshold() uint64
	FeeCapDefault() *assets.Wei
	LimitMultiplier() float32
	PriceDefault() *assets.Wei
	TipCapDefault() *assets.Wei
	PriceMax() *assets.Wei
	Mode() string
	bumpConfig
}

type fixedPriceEstimatorBlockHistoryConfig interface {
	EIP1559FeeCapBufferBlocks() uint16
}

// NewFixedPriceEstimator returns a new "FixedPrice" estimator which will
// always use the config default values for gas prices and limits
func NewFixedPriceEstimator(cfg fixedPriceEstimatorConfig, bhCfg fixedPriceEstimatorBlockHistoryConfig, lggr logger.Logger) EvmEstimator {
	return &fixedPriceEstimator{cfg, bhCfg, logger.Sugared(lggr.Named("FixedPriceEstimator"))}
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

func (f *fixedPriceEstimator) GetLegacyGas(_ context.Context, _ []byte, gasLimit uint32, maxGasPriceWei *assets.Wei, _ ...feetypes.Opt) (gasPrice *assets.Wei, chainSpecificGasLimit uint32, err error) {
	gasPrice = f.config.PriceDefault()
	gasPrice, chainSpecificGasLimit = capGasPrice(gasPrice, maxGasPriceWei, f.config.PriceMax(), gasLimit, f.config.LimitMultiplier())
	return
}

func (f *fixedPriceEstimator) BumpLegacyGas(_ context.Context, originalGasPrice *assets.Wei, originalGasLimit uint32, maxGasPriceWei *assets.Wei, _ []EvmPriorAttempt) (gasPrice *assets.Wei, gasLimit uint32, err error) {
	return BumpLegacyGasPriceOnly(f.config, f.lggr, f.config.PriceDefault(), originalGasPrice, originalGasLimit, maxGasPriceWei)
}

func (f *fixedPriceEstimator) GetDynamicFee(_ context.Context, originalGasLimit uint32, maxGasPriceWei *assets.Wei) (d DynamicFee, chainSpecificGasLimit uint32, err error) {
	gasTipCap := f.config.TipCapDefault()

	if gasTipCap == nil {
		return d, 0, errors.New("cannot calculate dynamic fee: EthGasTipCapDefault was not set")
	}
	chainSpecificGasLimit = commonfee.ApplyMultiplier(originalGasLimit, f.config.LimitMultiplier())

	var feeCap *assets.Wei
	if f.config.BumpThreshold() == 0 {
		// Gas bumping is disabled, just use the max fee cap
		feeCap = getMaxGasPrice(maxGasPriceWei, f.config.PriceMax())
	} else {
		// Need to leave headroom for bumping so we fallback to the default value here
		feeCap = f.config.FeeCapDefault()
	}

	return DynamicFee{
		FeeCap: feeCap,
		TipCap: gasTipCap,
	}, chainSpecificGasLimit, nil
}

func (f *fixedPriceEstimator) BumpDynamicFee(_ context.Context, originalFee DynamicFee, originalGasLimit uint32, maxGasPriceWei *assets.Wei, _ []EvmPriorAttempt) (bumped DynamicFee, chainSpecificGasLimit uint32, err error) {
	return BumpDynamicFeeOnly(f.config, f.bhConfig.EIP1559FeeCapBufferBlocks(), f.lggr, f.config.TipCapDefault(), nil, originalFee, originalGasLimit, maxGasPriceWei)
}
