package gas

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
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

type EvmFeeEstimator txmgrtypes.FeeEstimator[*evmtypes.Head, EvmFee, *assets.Wei, common.Hash]

// NewEstimator returns the estimator for a given config
func NewEstimator(lggr logger.Logger, ethClient evmclient.Client, cfg Config) EvmFeeEstimator {

	s := cfg.GasEstimatorMode()
	lggr.Infow(fmt.Sprintf("Initializing EVM gas estimator in mode: %s", s),
		"estimatorMode", s,
		"batchSize", cfg.BlockHistoryEstimatorBatchSize(),
		"blockDelay", cfg.BlockHistoryEstimatorBlockDelay(),
		"blockHistorySize", cfg.BlockHistoryEstimatorBlockHistorySize(),
		"eip1559FeeCapBufferBlocks", cfg.BlockHistoryEstimatorEIP1559FeeCapBufferBlocks(),
		"transactionPercentile", cfg.BlockHistoryEstimatorTransactionPercentile(),
		"eip1559DynamicFees", cfg.EvmEIP1559DynamicFees(),
		"gasBumpPercent", cfg.EvmGasBumpPercent(),
		"gasBumpThreshold", cfg.EvmGasBumpThreshold(),
		"gasBumpWei", cfg.EvmGasBumpWei(),
		"feeCapDefault", cfg.EvmGasFeeCapDefault(),
		"gasLimitMultiplier", cfg.EvmGasLimitMultiplier(),
		"gasPriceDefault", cfg.EvmGasPriceDefault(),
		"gasTipCapDefault", cfg.EvmGasTipCapDefault(),
		"gasTipCapMinimum", cfg.EvmGasTipCapMinimum(),
		"maxGasPriceWei", cfg.EvmMaxGasPriceWei(),
		"minGasPriceWei", cfg.EvmMinGasPriceWei(),
	)
	switch s {
	case "Arbitrum":
		return NewWrappedEvmEstimator(NewArbitrumEstimator(lggr, cfg, ethClient, ethClient), cfg)
	case "BlockHistory":
		return NewWrappedEvmEstimator(NewBlockHistoryEstimator(lggr, ethClient, cfg, *ethClient.ConfiguredChainID()), cfg)
	case "FixedPrice":
		return NewWrappedEvmEstimator(NewFixedPriceEstimator(cfg, lggr), cfg)
	case "Optimism2", "L2Suggested":
		return NewWrappedEvmEstimator(NewL2SuggestedPriceEstimator(lggr, ethClient), cfg)
	default:
		lggr.Warnf("GasEstimator: unrecognised mode '%s', falling back to FixedPriceEstimator", s)
		return NewWrappedEvmEstimator(NewFixedPriceEstimator(cfg, lggr), cfg)
	}
}

// DynamicFee encompasses both FeeCap and TipCap for EIP1559 transactions
type DynamicFee struct {
	FeeCap *assets.Wei
	TipCap *assets.Wei
}

type EvmPriorAttempt interface {
	txmgrtypes.PriorAttempt[EvmFee, common.Hash]

	GetGasPrice() *assets.Wei
	DynamicFee() DynamicFee
}

type evmPriorAttempt struct {
	txmgrtypes.PriorAttempt[EvmFee, common.Hash]
}

func (e evmPriorAttempt) GetGasPrice() *assets.Wei {
	return e.Fee().Legacy
}

func (e evmPriorAttempt) DynamicFee() DynamicFee {
	fee := e.Fee().Dynamic
	if fee == nil {
		return DynamicFee{}
	}
	return *fee
}

func MakeEvmPriorAttempts(attempts []txmgrtypes.PriorAttempt[EvmFee, common.Hash]) (out []EvmPriorAttempt) {
	for i := range attempts {
		out = append(out, MakeEvmPriorAttempt(attempts[i]))
	}
	return out
}

func MakeEvmPriorAttempt(a txmgrtypes.PriorAttempt[EvmFee, common.Hash]) EvmPriorAttempt {
	e := evmPriorAttempt{a}
	return &e
}

// Estimator provides an interface for estimating gas price and limit
//
//go:generate mockery --quiet --name EvmEstimator --output ./mocks/ --case=underscore
type EvmEstimator interface {
	txmgrtypes.HeadTrackable[*evmtypes.Head]
	services.ServiceCtx

	// GetLegacyGas Calculates initial gas fee for non-EIP1559 transaction
	// maxGasPriceWei parameter is the highest possible gas fee cap that the function will return
	GetLegacyGas(ctx context.Context, calldata []byte, gasLimit uint32, maxGasPriceWei *assets.Wei, opts ...txmgrtypes.Opt) (gasPrice *assets.Wei, chainSpecificGasLimit uint32, err error)
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

var _ txmgrtypes.Fee = (*EvmFee)(nil)

type EvmFee struct {
	Legacy  *assets.Wei
	Dynamic *DynamicFee
}

func (fee EvmFee) String() string {
	return fmt.Sprintf("{Legacy: %s, Dynamic: %+v}", fee.Legacy, fee.Dynamic)
}

// WrappedEvmEstimator provides a struct that wraps the EVM specific dynamic and legacy estimators into one estimator that conforms to the generic FeeEstimator
type WrappedEvmEstimator struct {
	EvmEstimator
	EIP1559Enabled bool
}

var _ EvmFeeEstimator = (*WrappedEvmEstimator)(nil)

func NewWrappedEvmEstimator(e EvmEstimator, cfg Config) EvmFeeEstimator {
	return &WrappedEvmEstimator{
		EvmEstimator:   e,
		EIP1559Enabled: cfg.EvmEIP1559DynamicFees(),
	}
}

func (e WrappedEvmEstimator) GetFee(ctx context.Context, calldata []byte, feeLimit uint32, maxFeePrice *assets.Wei, opts ...txmgrtypes.Opt) (fee EvmFee, chainSpecificFeeLimit uint32, err error) {
	// get dynamic fee
	if e.EIP1559Enabled {
		var dynamicFee DynamicFee
		dynamicFee, chainSpecificFeeLimit, err = e.EvmEstimator.GetDynamicFee(ctx, feeLimit, maxFeePrice)
		fee.Dynamic = &dynamicFee
		return
	}

	// get legacy fee
	fee.Legacy, chainSpecificFeeLimit, err = e.EvmEstimator.GetLegacyGas(ctx, calldata, feeLimit, maxFeePrice, opts...)
	return
}

func (e WrappedEvmEstimator) BumpFee(ctx context.Context, originalFee EvmFee, feeLimit uint32, maxFeePrice *assets.Wei, attempts []txmgrtypes.PriorAttempt[EvmFee, common.Hash]) (bumpedFee EvmFee, chainSpecificFeeLimit uint32, err error) {
	// validate only 1 fee type is present
	if (originalFee.Dynamic == nil && originalFee.Legacy == nil) || (originalFee.Dynamic != nil && originalFee.Legacy != nil) {
		err = errors.New("only one dynamic or legacy fee can be defined")
		return
	}

	// convert PriorAttempts to EvmPriorAttempts
	evmAttempts := MakeEvmPriorAttempts(attempts)

	// bump fee based on what fee the tx has previously used (not based on config)
	// bump dynamic original
	if originalFee.Dynamic != nil {
		var bumpedDynamic DynamicFee
		bumpedDynamic, chainSpecificFeeLimit, err = e.EvmEstimator.BumpDynamicFee(ctx, *originalFee.Dynamic, feeLimit, maxFeePrice, evmAttempts)
		bumpedFee.Dynamic = &bumpedDynamic
		return
	}

	// bump legacy fee
	bumpedFee.Legacy, chainSpecificFeeLimit, err = e.EvmEstimator.BumpLegacyGas(ctx, originalFee.Legacy, feeLimit, maxFeePrice, evmAttempts)
	return
}

func applyMultiplier(gasLimit uint32, multiplier float32) uint32 {
	return uint32(decimal.NewFromBigInt(big.NewInt(0).SetUint64(uint64(gasLimit)), 0).Mul(decimal.NewFromFloat32(multiplier)).IntPart())
}

// Config defines an interface for configuration in the gas package
//
//go:generate mockery --quiet --name Config --output ./mocks/ --case=underscore
type Config interface {
	BlockHistoryEstimatorBatchSize() uint32
	BlockHistoryEstimatorBlockDelay() uint16
	BlockHistoryEstimatorBlockHistorySize() uint16
	BlockHistoryEstimatorCheckInclusionPercentile() uint16
	BlockHistoryEstimatorCheckInclusionBlocks() uint16
	BlockHistoryEstimatorEIP1559FeeCapBufferBlocks() uint16
	BlockHistoryEstimatorTransactionPercentile() uint16
	ChainType() config.ChainType
	EvmEIP1559DynamicFees() bool
	EvmFinalityDepth() uint32
	EvmGasBumpPercent() uint16
	EvmGasBumpThreshold() uint64
	EvmGasBumpWei() *assets.Wei
	EvmGasFeeCapDefault() *assets.Wei
	EvmGasLimitMax() uint32
	EvmGasLimitMultiplier() float32
	EvmGasPriceDefault() *assets.Wei
	EvmGasTipCapDefault() *assets.Wei
	EvmGasTipCapMinimum() *assets.Wei
	EvmMaxGasPriceWei() *assets.Wei
	EvmMinGasPriceWei() *assets.Wei
	GasEstimatorMode() string
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

// BumpLegacyGasPriceOnly will increase the price and apply multiplier to the gas limit
func BumpLegacyGasPriceOnly(cfg Config, lggr logger.SugaredLogger, currentGasPrice, originalGasPrice *assets.Wei, originalGasLimit uint32, maxGasPriceWei *assets.Wei) (gasPrice *assets.Wei, chainSpecificGasLimit uint32, err error) {
	gasPrice, err = bumpGasPrice(cfg, lggr, currentGasPrice, originalGasPrice, maxGasPriceWei)
	if err != nil {
		return nil, 0, err
	}
	chainSpecificGasLimit = applyMultiplier(originalGasLimit, cfg.EvmGasLimitMultiplier())
	return
}

// bumpGasPrice computes the next gas price to attempt as the largest of:
// - A configured percentage bump (EVM.GasEstimator.BumpPercent) on top of the baseline price.
// - A configured fixed amount of Wei (ETH_GAS_PRICE_WEI) on top of the baseline price.
// The baseline price is the maximum of the previous gas price attempt and the node's current gas price.
func bumpGasPrice(cfg Config, lggr logger.SugaredLogger, currentGasPrice, originalGasPrice *assets.Wei, maxGasPriceWei *assets.Wei) (*assets.Wei, error) {
	maxGasPrice := getMaxGasPrice(maxGasPriceWei, cfg)

	bumpedGasPrice := assets.MaxWei(
		originalGasPrice.AddPercentage(cfg.EvmGasBumpPercent()),
		originalGasPrice.Add(cfg.EvmGasBumpWei()),
	)

	if currentGasPrice != nil {
		if currentGasPrice.Cmp(maxGasPrice) > 0 {
			// Shouldn't happen because the estimator should not be allowed to
			// estimate a higher gas than the maximum allowed
			lggr.AssumptionViolationf("Ignoring current gas price of %s that would exceed max gas price of %s", currentGasPrice.String(), maxGasPrice.String())
		} else if bumpedGasPrice.Cmp(currentGasPrice) < 0 {
			// If the current gas price is higher than the old price bumped, use that instead
			bumpedGasPrice = currentGasPrice
		}
	}
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
func BumpDynamicFeeOnly(config Config, lggr logger.SugaredLogger, currentTipCap, currentBaseFee *assets.Wei, originalFee DynamicFee, originalGasLimit uint32, maxGasPriceWei *assets.Wei) (bumped DynamicFee, chainSpecificGasLimit uint32, err error) {
	bumped, err = bumpDynamicFee(config, lggr, currentTipCap, currentBaseFee, originalFee, maxGasPriceWei)
	if err != nil {
		return bumped, 0, err
	}
	chainSpecificGasLimit = applyMultiplier(originalGasLimit, config.EvmGasLimitMultiplier())
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
func bumpDynamicFee(cfg Config, lggr logger.SugaredLogger, currentTipCap, currentBaseFee *assets.Wei, originalFee DynamicFee, maxGasPriceWei *assets.Wei) (bumpedFee DynamicFee, err error) {
	maxGasPrice := getMaxGasPrice(maxGasPriceWei, cfg)
	baselineTipCap := assets.MaxWei(originalFee.TipCap, cfg.EvmGasTipCapDefault())

	bumpedTipCap := assets.MaxWei(
		baselineTipCap.AddPercentage(cfg.EvmGasBumpPercent()),
		baselineTipCap.Add(cfg.EvmGasBumpWei()),
	)

	if currentTipCap != nil {
		if currentTipCap.Cmp(maxGasPrice) > 0 {
			lggr.AssumptionViolationf("Ignoring current tip cap of %s that would exceed max gas price of %s", currentTipCap.String(), maxGasPrice.String())
		} else if bumpedTipCap.Cmp(currentTipCap) < 0 {
			// If the current gas tip cap is higher than the old tip cap with bump applied, use that instead
			bumpedTipCap = currentTipCap
		}
	}
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
		originalFee.FeeCap.AddPercentage(cfg.EvmGasBumpPercent()),
		originalFee.FeeCap.Add(cfg.EvmGasBumpWei()),
	)

	if currentBaseFee != nil {
		if currentBaseFee.Cmp(maxGasPrice) > 0 {
			lggr.Warnf("Ignoring current base fee of %s which is greater than max gas price of %s", currentBaseFee.String(), maxGasPrice.String())
		} else {
			currentFeeCap := calcFeeCap(currentBaseFee, cfg, bumpedTipCap, maxGasPrice)
			bumpedFeeCap = assets.WeiMax(bumpedFeeCap, currentFeeCap)
		}
	}

	if bumpedFeeCap.Cmp(maxGasPrice) > 0 {
		return bumpedFee, errors.Wrapf(ErrBumpGasExceedsLimit, "bumped fee cap of %s would exceed configured max gas price of %s (original fee: tip cap %s, fee cap %s). %s",
			bumpedFeeCap.String(), maxGasPrice, originalFee.TipCap.String(), originalFee.FeeCap.String(), label.NodeConnectivityProblemWarning)
	}

	return DynamicFee{FeeCap: bumpedFeeCap, TipCap: bumpedTipCap}, nil
}

func getMaxGasPrice(userSpecifiedMax *assets.Wei, config Config) *assets.Wei {
	return assets.WeiMin(config.EvmMaxGasPriceWei(), userSpecifiedMax)
}

func capGasPrice(calculatedGasPrice, userSpecifiedMax *assets.Wei, config Config) *assets.Wei {
	maxGasPrice := getMaxGasPrice(userSpecifiedMax, config)
	return assets.WeiMin(calculatedGasPrice, maxGasPrice)
}
