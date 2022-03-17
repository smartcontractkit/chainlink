package gas

import (
	"context"
	"math/big"

	"github.com/pkg/errors"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
)

var _ Estimator = &fixedPriceEstimator{}

type fixedPriceEstimator struct {
	config Config
	lggr   logger.SugaredLogger
}

// NewFixedPriceEstimator returns a new "FixedPrice" estimator which will
// always use the config default values for gas prices and limits
func NewFixedPriceEstimator(cfg Config, lggr logger.Logger) Estimator {
	return &fixedPriceEstimator{cfg, logger.Sugared(lggr.Named("FixedPriceEstimator"))}
}

func (f *fixedPriceEstimator) Start(context.Context) error                           { return nil }
func (f *fixedPriceEstimator) Close() error                                          { return nil }
func (f *fixedPriceEstimator) OnNewLongestChain(_ context.Context, _ *evmtypes.Head) {}

func (f *fixedPriceEstimator) GetLegacyGas(_ []byte, gasLimit uint64, _ ...Opt) (gasPrice *big.Int, chainSpecificGasLimit uint64, err error) {
	gasPrice = f.config.EvmGasPriceDefault()
	chainSpecificGasLimit = applyMultiplier(gasLimit, f.config.EvmGasLimitMultiplier())
	return
}

func (f *fixedPriceEstimator) BumpLegacyGas(originalGasPrice *big.Int, originalGasLimit uint64) (gasPrice *big.Int, gasLimit uint64, err error) {
	return BumpLegacyGasPriceOnly(f.config, f.lggr, f.config.EvmGasPriceDefault(), originalGasPrice, originalGasLimit)
}

func (f *fixedPriceEstimator) GetDynamicFee(originalGasLimit uint64) (d DynamicFee, chainSpecificGasLimit uint64, err error) {
	gasTipCap := f.config.EvmGasTipCapDefault()
	if gasTipCap == nil {
		return d, 0, errors.New("cannot calculate dynamic fee: EthGasTipCapDefault was not set")
	}
	chainSpecificGasLimit = applyMultiplier(originalGasLimit, f.config.EvmGasLimitMultiplier())

	var feeCap *big.Int
	if f.config.EvmGasBumpThreshold() == 0 {
		// Gas bumping is disabled, just use the max fee cap
		feeCap = f.config.EvmMaxGasPriceWei()
	} else {
		// Need to leave headroom for bumping so we fallback to the default value here
		feeCap = f.config.EvmGasFeeCapDefault()
	}

	return DynamicFee{
		FeeCap: feeCap,
		TipCap: gasTipCap,
	}, chainSpecificGasLimit, nil
}

func (f *fixedPriceEstimator) BumpDynamicFee(originalFee DynamicFee, originalGasLimit uint64) (bumped DynamicFee, chainSpecificGasLimit uint64, err error) {
	return BumpDynamicFeeOnly(f.config, f.lggr, f.config.EvmGasTipCapDefault(), nil, originalFee, originalGasLimit)
}
