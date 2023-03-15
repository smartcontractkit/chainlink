package gas

import (
	"context"

	"github.com/pkg/errors"

	txmgrtypes "github.com/smartcontractkit/chainlink/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/core/assets"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
)

var _ EvmEstimator[*evmtypes.Head] = &fixedPriceEstimator[*evmtypes.Head]{}

type fixedPriceEstimator[H txmgrtypes.Head] struct {
	config Config
	lggr   logger.SugaredLogger
}

// NewFixedPriceEstimator returns a new "FixedPrice" estimator which will
// always use the config default values for gas prices and limits
func NewFixedPriceEstimator(cfg Config, lggr logger.Logger) EvmEstimator[*evmtypes.Head] {
	return &fixedPriceEstimator[*evmtypes.Head]{cfg, logger.Sugared(lggr.Named("FixedPriceEstimator"))}
}

func (f *fixedPriceEstimator[H]) Start(context.Context) error {
	if f.config.EvmGasBumpThreshold() == 0 && f.config.GasEstimatorMode() == "FixedPrice" {
		// EvmGasFeeCapDefault is ignored if fixed estimator mode is on and gas bumping is disabled
		if f.config.EvmGasFeeCapDefault().Cmp(f.config.EvmMaxGasPriceWei()) != 0 {
			f.lggr.Infof("You are using FixedPrice estimator with gas bumping disabled. ETH_MAX_GAS_PRICE_WEI (value: %s) will be used as the FeeCap for transactions", f.config.EvmMaxGasPriceWei())
		}
	}
	return nil
}
func (f *fixedPriceEstimator[H]) Close() error                             { return nil }
func (f *fixedPriceEstimator[H]) OnNewLongestChain(_ context.Context, _ H) {}

func (f *fixedPriceEstimator[H]) GetLegacyGas(_ context.Context, _ []byte, gasLimit uint32, maxGasPriceWei *assets.Wei, _ ...txmgrtypes.Opt) (gasPrice *assets.Wei, chainSpecificGasLimit uint32, err error) {
	gasPrice = f.config.EvmGasPriceDefault()
	chainSpecificGasLimit = applyMultiplier(gasLimit, f.config.EvmGasLimitMultiplier())
	gasPrice = capGasPrice(gasPrice, maxGasPriceWei, f.config)
	return
}

func (f *fixedPriceEstimator[H]) BumpLegacyGas(_ context.Context, originalGasPrice *assets.Wei, originalGasLimit uint32, maxGasPriceWei *assets.Wei, _ []EvmPriorAttempt) (gasPrice *assets.Wei, gasLimit uint32, err error) {
	return BumpLegacyGasPriceOnly(f.config, f.lggr, f.config.EvmGasPriceDefault(), originalGasPrice, originalGasLimit, maxGasPriceWei)
}

func (f *fixedPriceEstimator[H]) GetDynamicFee(_ context.Context, originalGasLimit uint32, maxGasPriceWei *assets.Wei) (d DynamicFee, chainSpecificGasLimit uint32, err error) {
	gasTipCap := f.config.EvmGasTipCapDefault()
	if gasTipCap == nil {
		return d, 0, errors.New("cannot calculate dynamic fee: EthGasTipCapDefault was not set")
	}
	chainSpecificGasLimit = applyMultiplier(originalGasLimit, f.config.EvmGasLimitMultiplier())

	var feeCap *assets.Wei
	if f.config.EvmGasBumpThreshold() == 0 {
		// Gas bumping is disabled, just use the max fee cap
		feeCap = getMaxGasPrice(maxGasPriceWei, f.config)
	} else {
		// Need to leave headroom for bumping so we fallback to the default value here
		feeCap = f.config.EvmGasFeeCapDefault()
	}

	return DynamicFee{
		FeeCap: feeCap,
		TipCap: gasTipCap,
	}, chainSpecificGasLimit, nil
}

func (f *fixedPriceEstimator[H]) BumpDynamicFee(_ context.Context, originalFee DynamicFee, originalGasLimit uint32, maxGasPriceWei *assets.Wei, _ []EvmPriorAttempt) (bumped DynamicFee, chainSpecificGasLimit uint32, err error) {
	return BumpDynamicFeeOnly(f.config, f.lggr, f.config.EvmGasTipCapDefault(), nil, originalFee, originalGasLimit, maxGasPriceWei)
}
