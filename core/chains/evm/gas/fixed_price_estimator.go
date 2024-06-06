package gas

import (
	"context"

	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonfee "github.com/smartcontractkit/chainlink/v2/common/fee"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

var _ EvmEstimator = (*fixedPriceEstimator)(nil)

type fixedPriceEstimator struct {
	config   fixedPriceEstimatorConfig
	bhConfig fixedPriceEstimatorBlockHistoryConfig
	lggr     logger.SugaredLogger
	l1Oracle rollups.L1Oracle
}
type bumpConfig interface {
	LimitMultiplier() float32
	PriceMax() *assets.Wei
	BumpPercent() uint16
	BumpMin() *assets.Wei
	TipCapDefault() *assets.Wei
}

type fixedPriceEstimatorConfig interface {
	BumpThreshold() uint64
	FeeCapDefault() *assets.Wei
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
func NewFixedPriceEstimator(cfg fixedPriceEstimatorConfig, ethClient feeEstimatorClient, bhCfg fixedPriceEstimatorBlockHistoryConfig, lggr logger.Logger, l1Oracle rollups.L1Oracle) EvmEstimator {
	return &fixedPriceEstimator{cfg, bhCfg, logger.Sugared(logger.Named(lggr, "FixedPriceEstimator")), l1Oracle}
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

func (f *fixedPriceEstimator) GetLegacyGas(_ context.Context, _ []byte, gasLimit uint64, maxGasPriceWei *assets.Wei, _ ...feetypes.Opt) (*assets.Wei, uint64, error) {
	gasPrice := commonfee.CalculateFee(f.config.PriceDefault().ToInt(), maxGasPriceWei.ToInt(), f.config.PriceMax().ToInt())
	chainSpecificGasLimit := gasLimit
	return assets.NewWei(gasPrice), chainSpecificGasLimit, nil
}

func (f *fixedPriceEstimator) BumpLegacyGas(
	_ context.Context,
	originalGasPrice *assets.Wei,
	originalGasLimit uint64,
	maxGasPriceWei *assets.Wei,
	_ []EvmPriorAttempt,
) (*assets.Wei, uint64, error) {
	gasPrice, err := commonfee.CalculateBumpedFee(
		f.lggr,
		f.config.PriceDefault().ToInt(),
		originalGasPrice.ToInt(),
		maxGasPriceWei.ToInt(),
		f.config.PriceMax().ToInt(),
		f.config.BumpMin().ToInt(),
		f.config.BumpPercent(),
		assets.FormatWei,
	)
	if err != nil {
		return nil, 0, err
	}

	chainSpecificGasLimit := originalGasLimit
	return assets.NewWei(gasPrice), chainSpecificGasLimit, err
}

func (f *fixedPriceEstimator) GetDynamicFee(_ context.Context, maxGasPriceWei *assets.Wei) (d DynamicFee, err error) {
	gasTipCap := f.config.TipCapDefault()

	if gasTipCap == nil {
		return d, pkgerrors.New("cannot calculate dynamic fee: EthGasTipCapDefault was not set")
	}

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
	}, nil
}

func (f *fixedPriceEstimator) BumpDynamicFee(
	_ context.Context,
	originalFee DynamicFee,
	maxGasPriceWei *assets.Wei,
	_ []EvmPriorAttempt,
) (bumped DynamicFee, err error) {
	return BumpDynamicFeeOnly(
		f.config,
		f.bhConfig.EIP1559FeeCapBufferBlocks(),
		f.lggr,
		f.config.TipCapDefault(),
		nil,
		originalFee,
		maxGasPriceWei,
	)
}

func (f *fixedPriceEstimator) L1Oracle() rollups.L1Oracle {
	return f.l1Oracle
}

func (f *fixedPriceEstimator) Name() string                                          { return f.lggr.Name() }
func (f *fixedPriceEstimator) Ready() error                                          { return nil }
func (f *fixedPriceEstimator) HealthReport() map[string]error                        { return map[string]error{} }
func (f *fixedPriceEstimator) Close() error                                          { return nil }
func (f *fixedPriceEstimator) OnNewLongestChain(_ context.Context, _ *evmtypes.Head) {}
