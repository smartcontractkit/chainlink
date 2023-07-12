package gas

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	commonfee "github.com/smartcontractkit/chainlink/v2/common/fee"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/label"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

var (
	ErrBumpGasExceedsLimit = errors.New("gas bump exceeds limit")
	ErrBump                = errors.New("gas bump failed")
	ErrConnectivity        = errors.New("transaction propagation issue: transactions are not being mined")
)

func IsBumpErr(err error) bool {
	return err != nil && (errors.Is(err, ErrBumpGasExceedsLimit) || errors.Is(err, ErrBump) || errors.Is(err, ErrConnectivity))
}

// EvmFeeEstimator provides a unified interface that wraps EvmEstimator and can determine if legacy or dynamic fee estimation should be used
//
//go:generate mockery --quiet --name EvmFeeEstimator --output ./mocks/ --case=underscore
type EvmFeeEstimator interface {
	services.ServiceCtx
	commontypes.HeadTrackable[*evmtypes.Head, common.Hash]

	GetFee(ctx context.Context, calldata []byte, feeLimit uint32, maxFeePrice *assets.Wei, opts ...feetypes.Opt) (fee EvmFee, chainSpecificFeeLimit uint32, err error)
	BumpFee(ctx context.Context, originalFee EvmFee, feeLimit uint32, maxFeePrice *assets.Wei, attempts []EvmPriorAttempt) (bumpedFee EvmFee, chainSpecificFeeLimit uint32, err error)
}

// NewEstimator returns the estimator for a given config
func NewEstimator(lggr logger.Logger, ethClient evmclient.Client, cfg Config, geCfg evmconfig.GasEstimator) EvmFeeEstimator {
	bh := geCfg.BlockHistory()
	s := geCfg.Mode()
	lggr.Infow(fmt.Sprintf("Initializing EVM gas estimator in mode: %s", s),
		"estimatorMode", s,
		"batchSize", bh.BatchSize(),
		"blockDelay", bh.BlockDelay(),
		"blockHistorySize", bh.BlockHistorySize(),
		"eip1559FeeCapBufferBlocks", bh.EIP1559FeeCapBufferBlocks(),
		"transactionPercentile", bh.TransactionPercentile(),
		"eip1559DynamicFees", geCfg.EIP1559DynamicFees(),
		"gasBumpPercent", geCfg.BumpPercent(),
		"gasBumpThreshold", geCfg.BumpThreshold(),
		"bumpMin", geCfg.BumpMin(),
		"feeCapDefault", geCfg.FeeCapDefault(),
		"limitMultiplier", geCfg.LimitMultiplier(),
		"priceDefault", geCfg.PriceDefault(),
		"tipCapDefault", geCfg.TipCapDefault(),
		"tipCapMin", geCfg.TipCapMin(),
		"priceMax", geCfg.PriceMax(),
		"priceMin", geCfg.PriceMin(),
	)
	df := geCfg.EIP1559DynamicFees()
	switch s {
	case "Arbitrum":
		return NewWrappedEvmEstimator(NewArbitrumEstimator(lggr, geCfg, ethClient, ethClient), df)
	case "BlockHistory":
		return NewWrappedEvmEstimator(NewBlockHistoryEstimator(lggr, ethClient, cfg, geCfg, bh, *ethClient.ConfiguredChainID()), df)
	case "FixedPrice":
		return NewWrappedEvmEstimator(NewFixedPriceEstimator(geCfg, bh, lggr), df)
	case "Optimism2", "L2Suggested":
		return NewWrappedEvmEstimator(NewL2SuggestedPriceEstimator(lggr, ethClient), df)
	default:
		lggr.Warnf("GasEstimator: unrecognised mode '%s', falling back to FixedPriceEstimator", s)
		return NewWrappedEvmEstimator(NewFixedPriceEstimator(geCfg, bh, lggr), df)
	}
}

// DynamicFee encompasses both FeeCap and TipCap for EIP1559 transactions
type DynamicFee struct {
	FeeCap *assets.Wei
	TipCap *assets.Wei
}

type EvmPriorAttempt struct {
	ChainSpecificFeeLimit   uint32
	BroadcastBeforeBlockNum *int64
	TxHash                  common.Hash
	TxType                  int
	GasPrice                *assets.Wei
	DynamicFee              DynamicFee
}

// Estimator provides an interface for estimating gas price and limit
//
//go:generate mockery --quiet --name EvmEstimator --output ./mocks/ --case=underscore
type EvmEstimator interface {
	commontypes.HeadTrackable[*evmtypes.Head, common.Hash]
	services.ServiceCtx

	// GetLegacyGas Calculates initial gas fee for non-EIP1559 transaction
	// maxGasPriceWei parameter is the highest possible gas fee cap that the function will return
	GetLegacyGas(ctx context.Context, calldata []byte, gasLimit uint32, maxGasPriceWei *assets.Wei, opts ...feetypes.Opt) (gasPrice *assets.Wei, chainSpecificGasLimit uint32, err error)
	// BumpLegacyGas Increases gas price and/or limit for non-EIP1559 transactions
	// if the bumped gas fee is greater than maxGasPriceWei, the method returns an error
	// attempts must:
	//   - be sorted in order from highest price to lowest price
	//   - all be of transaction type 0x0 or 0x1
	BumpLegacyGas(ctx context.Context, originalGasPrice *assets.Wei, gasLimit uint32, maxGasPriceWei *assets.Wei, attempts []EvmPriorAttempt) (bumpedGasPrice *assets.Wei, chainSpecificGasLimit uint32, err error)
	// GetDynamicFee Calculates initial gas fee for gas for EIP1559 transactions
	// maxGasPriceWei parameter is the highest possible gas fee cap that the function will return
	GetDynamicFee(ctx context.Context, gasLimit uint32, maxGasPriceWei *assets.Wei) (fee DynamicFee, chainSpecificGasLimit uint32, err error)
	// BumpDynamicFee Increases gas price and/or limit for non-EIP1559 transactions
	// if the bumped gas fee or tip caps are greater than maxGasPriceWei, the method returns an error
	// attempts must:
	//   - be sorted in order from highest price to lowest price
	//   - all be of transaction type 0x2
	BumpDynamicFee(ctx context.Context, original DynamicFee, gasLimit uint32, maxGasPriceWei *assets.Wei, attempts []EvmPriorAttempt) (bumped DynamicFee, chainSpecificGasLimit uint32, err error)
}

var _ feetypes.Fee = (*EvmFee)(nil)

type EvmFee struct {
	// legacy fees
	Legacy *assets.Wei

	// dynamic/EIP1559 fees
	DynamicFeeCap *assets.Wei
	DynamicTipCap *assets.Wei
}

func (fee EvmFee) String() string {
	return fmt.Sprintf("{Legacy: %s, DynamicFeeCap: %s, DynamicTipCap: %s}", fee.Legacy, fee.DynamicFeeCap, fee.DynamicTipCap)
}

func (fee EvmFee) ValidDynamic() bool {
	return fee.DynamicFeeCap != nil && fee.DynamicTipCap != nil
}

// WrappedEvmEstimator provides a struct that wraps the EVM specific dynamic and legacy estimators into one estimator that conforms to the generic FeeEstimator
type WrappedEvmEstimator struct {
	EvmEstimator
	EIP1559Enabled bool
}

var _ EvmFeeEstimator = (*WrappedEvmEstimator)(nil)

func NewWrappedEvmEstimator(e EvmEstimator, eip1559Enabled bool) EvmFeeEstimator {
	return &WrappedEvmEstimator{
		EvmEstimator:   e,
		EIP1559Enabled: eip1559Enabled,
	}
}

func (e WrappedEvmEstimator) GetFee(ctx context.Context, calldata []byte, feeLimit uint32, maxFeePrice *assets.Wei, opts ...feetypes.Opt) (fee EvmFee, chainSpecificFeeLimit uint32, err error) {
	// get dynamic fee
	if e.EIP1559Enabled {
		var dynamicFee DynamicFee
		dynamicFee, chainSpecificFeeLimit, err = e.EvmEstimator.GetDynamicFee(ctx, feeLimit, maxFeePrice)
		fee.DynamicFeeCap = dynamicFee.FeeCap
		fee.DynamicTipCap = dynamicFee.TipCap
		return
	}

	// get legacy fee
	fee.Legacy, chainSpecificFeeLimit, err = e.EvmEstimator.GetLegacyGas(ctx, calldata, feeLimit, maxFeePrice, opts...)
	return
}

func (e WrappedEvmEstimator) BumpFee(ctx context.Context, originalFee EvmFee, feeLimit uint32, maxFeePrice *assets.Wei, attempts []EvmPriorAttempt) (bumpedFee EvmFee, chainSpecificFeeLimit uint32, err error) {
	// validate only 1 fee type is present
	if (!originalFee.ValidDynamic() && originalFee.Legacy == nil) || (originalFee.ValidDynamic() && originalFee.Legacy != nil) {
		err = errors.New("only one dynamic or legacy fee can be defined")
		return
	}

	// bump fee based on what fee the tx has previously used (not based on config)
	// bump dynamic original
	if originalFee.ValidDynamic() {
		var bumpedDynamic DynamicFee
		bumpedDynamic, chainSpecificFeeLimit, err = e.EvmEstimator.BumpDynamicFee(ctx,
			DynamicFee{
				TipCap: originalFee.DynamicTipCap,
				FeeCap: originalFee.DynamicFeeCap,
			}, feeLimit, maxFeePrice, attempts)
		bumpedFee.DynamicFeeCap = bumpedDynamic.FeeCap
		bumpedFee.DynamicTipCap = bumpedDynamic.TipCap
		return
	}

	// bump legacy fee
	bumpedFee.Legacy, chainSpecificFeeLimit, err = e.EvmEstimator.BumpLegacyGas(ctx, originalFee.Legacy, feeLimit, maxFeePrice, attempts)
	return
}

// Config defines an interface for configuration in the gas package
//
//go:generate mockery --quiet --name Config --output ./mocks/ --case=underscore
type Config interface {
	ChainType() config.ChainType
	FinalityDepth() uint32
}

type GasEstimatorConfig interface {
	EIP1559DynamicFees() bool
	BumpPercent() uint16
	BumpThreshold() uint64
	BumpMin() *assets.Wei
	FeeCapDefault() *assets.Wei
	LimitMax() uint32
	LimitMultiplier() float32
	PriceDefault() *assets.Wei
	TipCapDefault() *assets.Wei
	TipCapMin() *assets.Wei
	PriceMin() *assets.Wei
	PriceMax() *assets.Wei
	Mode() string
}

type BlockHistoryConfig interface {
	evmconfig.BlockHistory
}

// Int64ToHex converts an int64 into go-ethereum's hex representation
func Int64ToHex(n int64) string {
	return hexutil.EncodeBig(big.NewInt(n))
}

// HexToInt64 performs the inverse of Int64ToHex
// Returns 0 on invalid input
func HexToInt64(input interface{}) int64 {
	switch v := input.(type) {
	case string:
		big, err := hexutil.DecodeBig(v)
		if err != nil {
			return 0
		}
		return big.Int64()
	case []byte:
		big, err := hexutil.DecodeBig(string(v))
		if err != nil {
			return 0
		}
		return big.Int64()
	default:
		return 0
	}
}

type bumpConfig interface {
	LimitMultiplier() float32
	PriceMax() *assets.Wei
	BumpPercent() uint16
	BumpMin() *assets.Wei
	TipCapDefault() *assets.Wei
}

// BumpLegacyGasPriceOnly will increase the price and apply multiplier to the gas limit
func BumpLegacyGasPriceOnly(cfg bumpConfig, lggr logger.SugaredLogger, currentGasPrice, originalGasPrice *assets.Wei, originalGasLimit uint32, maxGasPriceWei *assets.Wei) (gasPrice *assets.Wei, chainSpecificGasLimit uint32, err error) {
	gasPrice, err = bumpGasPrice(cfg, lggr, currentGasPrice, originalGasPrice, maxGasPriceWei)
	if err != nil {
		return nil, 0, err
	}
	chainSpecificGasLimit = commonfee.ApplyMultiplier(originalGasLimit, cfg.LimitMultiplier())
	return
}

// bumpGasPrice computes the next gas price to attempt as the largest of:
// - A configured percentage bump (EVM.GasEstimator.BumpPercent) on top of the baseline price.
// - A configured fixed amount of Wei (ETH_GAS_PRICE_WEI) on top of the baseline price.
// The baseline price is the maximum of the previous gas price attempt and the node's current gas price.
func bumpGasPrice(cfg bumpConfig, lggr logger.SugaredLogger, currentGasPrice, originalGasPrice, maxGasPriceWei *assets.Wei) (*assets.Wei, error) {
	maxGasPrice := getMaxGasPrice(maxGasPriceWei, cfg.PriceMax())
	bumpedGasPrice := bumpFeePrice(originalGasPrice, cfg.BumpPercent(), cfg.BumpMin())

	// Update bumpedGasPrice if currentGasPrice is higher than bumpedGasPrice and within maxGasPrice
	bumpedGasPrice = maxBumpedFee(lggr, currentGasPrice, bumpedGasPrice, maxGasPrice, "gas price")

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

// BumpDynamicFeeOnly bumps the tip cap and max gas price if necessary
func BumpDynamicFeeOnly(config bumpConfig, feeCapBufferBlocks uint16, lggr logger.SugaredLogger, currentTipCap, currentBaseFee *assets.Wei, originalFee DynamicFee, originalGasLimit uint32, maxGasPriceWei *assets.Wei) (bumped DynamicFee, chainSpecificGasLimit uint32, err error) {
	bumped, err = bumpDynamicFee(config, feeCapBufferBlocks, lggr, currentTipCap, currentBaseFee, originalFee, maxGasPriceWei)
	if err != nil {
		return bumped, 0, err
	}
	chainSpecificGasLimit = commonfee.ApplyMultiplier(originalGasLimit, config.LimitMultiplier())
	return
}

// bumpDynamicFee computes the next tip cap to attempt as the largest of:
// - A configured percentage bump (EVM.GasEstimator.BumpPercent) on top of the baseline tip cap.
// - A configured fixed amount of Wei (ETH_GAS_PRICE_WEI) on top of the baseline tip cap.
// The baseline tip cap is the maximum of the previous tip cap attempt and the node's current tip cap.
// It increases the max fee cap by GasBumpPercent
//
// NOTE: We would prefer to have set a large FeeCap and leave it fixed, bumping
// the Tip only. Unfortunately due to a flaw of how EIP-1559 is implemented we
// have to bump FeeCap by at least 10% each time we bump the tip cap.
// See: https://github.com/ethereum/go-ethereum/issues/24284
func bumpDynamicFee(cfg bumpConfig, feeCapBufferBlocks uint16, lggr logger.SugaredLogger, currentTipCap, currentBaseFee *assets.Wei, originalFee DynamicFee, maxGasPriceWei *assets.Wei) (bumpedFee DynamicFee, err error) {
	maxGasPrice := getMaxGasPrice(maxGasPriceWei, cfg.PriceMax())
	baselineTipCap := assets.MaxWei(originalFee.TipCap, cfg.TipCapDefault())
	bumpedTipCap := bumpFeePrice(baselineTipCap, cfg.BumpPercent(), cfg.BumpMin())

	// Update bumpedTipCap if currentTipCap is higher than bumpedTipCap and within maxGasPrice
	bumpedTipCap = maxBumpedFee(lggr, currentTipCap, bumpedTipCap, maxGasPrice, "tip cap")

	if bumpedTipCap.Cmp(maxGasPrice) > 0 {
		return bumpedFee, errors.Wrapf(ErrBumpGasExceedsLimit, "bumped tip cap of %s would exceed configured max gas price of %s (original fee: tip cap %s, fee cap %s). %s",
			bumpedTipCap.String(), maxGasPrice, originalFee.TipCap.String(), originalFee.FeeCap.String(), label.NodeConnectivityProblemWarning)
	} else if bumpedTipCap.Cmp(originalFee.TipCap) <= 0 {
		// NOTE: This really shouldn't happen since we enforce minimums for
		// EVM.GasEstimator.BumpPercent and EVM.GasEstimator.BumpMin in the config validation,
		// but it's here anyway for a "belts and braces" approach
		return bumpedFee, errors.Wrapf(ErrBump, "bumped gas tip cap of %s is less than or equal to original gas tip cap of %s."+
			" ACTION REQUIRED: This is a configuration error, you must increase either "+
			"EVM.GasEstimator.BumpPercent or EVM.GasEstimator.BumpMin", bumpedTipCap.String(), originalFee.TipCap.String())
	}

	// Always bump the FeeCap by at least the bump percentage (should be greater than or
	// equal to than geth's configured bump minimum which is 10%)
	// See: https://github.com/ethereum/go-ethereum/blob/bff330335b94af3643ac2fb809793f77de3069d4/core/tx_list.go#L298
	bumpedFeeCap := assets.MaxWei(
		originalFee.FeeCap.AddPercentage(cfg.BumpPercent()),
		originalFee.FeeCap.Add(cfg.BumpMin()),
	)

	if currentBaseFee != nil {
		if currentBaseFee.Cmp(maxGasPrice) > 0 {
			lggr.Warnf("Ignoring current base fee of %s which is greater than max gas price of %s", currentBaseFee.String(), maxGasPrice.String())
		} else {
			currentFeeCap := calcFeeCap(currentBaseFee, int(feeCapBufferBlocks), bumpedTipCap, maxGasPrice)
			bumpedFeeCap = assets.WeiMax(bumpedFeeCap, currentFeeCap)
		}
	}

	if bumpedFeeCap.Cmp(maxGasPrice) > 0 {
		return bumpedFee, errors.Wrapf(ErrBumpGasExceedsLimit, "bumped fee cap of %s would exceed configured max gas price of %s (original fee: tip cap %s, fee cap %s). %s",
			bumpedFeeCap.String(), maxGasPrice, originalFee.TipCap.String(), originalFee.FeeCap.String(), label.NodeConnectivityProblemWarning)
	}

	return DynamicFee{FeeCap: bumpedFeeCap, TipCap: bumpedTipCap}, nil
}

func bumpFeePrice(originalFeePrice *assets.Wei, feeBumpPercent uint16, feeBumpUnits *assets.Wei) *assets.Wei {
	bumpedFeePrice := assets.MaxWei(
		originalFeePrice.AddPercentage(feeBumpPercent),
		originalFeePrice.Add(feeBumpUnits),
	)
	return bumpedFeePrice
}

func maxBumpedFee(lggr logger.SugaredLogger, currentFeePrice, bumpedFeePrice, maxGasPrice *assets.Wei, feeType string) *assets.Wei {
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

func getMaxGasPrice(userSpecifiedMax, maxGasPriceWei *assets.Wei) *assets.Wei {
	return assets.NewWei(commonfee.GetMaxFeePrice(userSpecifiedMax.ToInt(), maxGasPriceWei.ToInt()))
}

func capGasPrice(calculatedGasPrice, userSpecifiedMax, maxGasPriceWei *assets.Wei, gasLimit uint32, multiplier float32) (*assets.Wei, uint32) {
	maxGasPrice, chainSpecificGasLimit := commonfee.CapFeePrice(calculatedGasPrice.ToInt(), userSpecifiedMax.ToInt(), maxGasPriceWei.ToInt(), gasLimit, multiplier)
	return assets.NewWei(maxGasPrice), chainSpecificGasLimit
}
