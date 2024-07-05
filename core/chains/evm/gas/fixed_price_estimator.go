package gas

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type FixedPriceEstimator struct {
	lggr     logger.SugaredLogger
	config   fixedPriceEstimatorConfig
	l1Oracle rollups.L1Oracle
}

//go:generate mockery --quiet --name fixedPriceEstimatorConfig --output ./mocks/ --case=underscore --structname FixedPriceEstimatorConfig
type fixedPriceEstimatorConfig interface {
	PriceDefault() *assets.Wei
	TipCapDefault() *assets.Wei
	FeeCapDefault() *assets.Wei

	BumpPercent() uint16
}

func NewFixedPriceEstimator(lggr logger.Logger, config fixedPriceEstimatorConfig, l1Oracle rollups.L1Oracle) EvmEstimator {
	return &FixedPriceEstimator{logger.Sugared(logger.Named(lggr, "FixedPriceEstimator")), config, l1Oracle}
}

func (f *FixedPriceEstimator) Start(context.Context) error { return nil }

func (f *FixedPriceEstimator) GetLegacyGas(_ context.Context, _ []byte, gasLimit uint64, maxPrice *assets.Wei, _ ...feetypes.Opt) (*assets.Wei, uint64, error) {
	gasPrice := assets.WeiMin(f.config.PriceDefault(), maxPrice)
	return gasPrice, gasLimit, nil
}

func (f *FixedPriceEstimator) BumpLegacyGas(_ context.Context, originalGasPrice *assets.Wei, gasLimit uint64, maxPrice *assets.Wei, _ []EvmPriorAttempt) (*assets.Wei, uint64, error) {
	// Sanitize original fee input
	if originalGasPrice == nil || originalGasPrice.Cmp(maxPrice) >= 0 {
		return nil, 0, fmt.Errorf("error while retrieving original gas price: originalGasPrice: %s. Maximum price configured: %s", originalGasPrice, maxPrice)
	}

	bumpedGasPrice := originalGasPrice.AddPercentage(f.config.BumpPercent())
	bumpedGasPrice = assets.WeiMin(bumpedGasPrice, maxPrice)
	return bumpedGasPrice, gasLimit, nil
}

func (f *FixedPriceEstimator) GetDynamicFee(_ context.Context, maxPrice *assets.Wei) (d DynamicFee, err error) {
	maxPriorityFeePerGas := assets.WeiMin(f.config.TipCapDefault(), maxPrice)
	maxFeePerGas := assets.WeiMin(f.config.FeeCapDefault(), maxPrice)

	return DynamicFee{FeeCap: maxFeePerGas, TipCap: maxPriorityFeePerGas}, nil
}

func (f *FixedPriceEstimator) BumpDynamicFee(_ context.Context, originalFee DynamicFee, maxPrice *assets.Wei, _ []EvmPriorAttempt) (bumpedFee DynamicFee, err error) {
	// Sanitize original fee input
	if originalFee.FeeCap == nil ||
		originalFee.TipCap == nil ||
		((originalFee.TipCap.Cmp(originalFee.FeeCap)) > 0) ||
		(originalFee.FeeCap.Cmp(maxPrice) >= 0) {
		return bumpedFee, fmt.Errorf("error while retrieving original dynamic fees: (originalFeePerGas: %s - originalPriorityFeePerGas: %s). Maximum price configured: %s",
			originalFee.FeeCap, originalFee.TipCap, maxPrice)
	}

	bumpedMaxPriorityFeePerGas := originalFee.TipCap.AddPercentage(f.config.BumpPercent())
	bumpedMaxFeePerGas := originalFee.FeeCap.AddPercentage(f.config.BumpPercent())

	bumpedMaxFeePerGas = assets.WeiMin(bumpedMaxFeePerGas, maxPrice)
	bumpedMaxPriorityFeePerGas = assets.WeiMin(bumpedMaxPriorityFeePerGas, maxPrice)

	bumpedFee = DynamicFee{FeeCap: bumpedMaxFeePerGas, TipCap: bumpedMaxPriorityFeePerGas}
	return bumpedFee, nil
}

func (f *FixedPriceEstimator) L1Oracle() rollups.L1Oracle {
	return f.l1Oracle
}

func (f *FixedPriceEstimator) Name() string                                          { return f.lggr.Name() }
func (f *FixedPriceEstimator) Ready() error                                          { return nil }
func (f *FixedPriceEstimator) HealthReport() map[string]error                        { return map[string]error{} }
func (f *FixedPriceEstimator) Close() error                                          { return nil }
func (f *FixedPriceEstimator) OnNewLongestChain(_ context.Context, _ *evmtypes.Head) {}
