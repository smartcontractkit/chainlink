package gas

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	bigmath "github.com/smartcontractkit/chainlink-common/pkg/utils/big_math"

	commonfee "github.com/smartcontractkit/chainlink/v2/common/fee"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/label"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// EstimateGasBuffer is a multiplier applied to estimated gas when the EstimateLimit feature is enabled
const EstimateGasBuffer = float32(1.15)

// EvmFeeEstimator provides a unified interface that wraps EvmEstimator and can determine if legacy or dynamic fee estimation should be used
type EvmFeeEstimator interface {
	services.Service
	headtracker.HeadTrackable[*evmtypes.Head, common.Hash]

	// L1Oracle returns the L1 gas price oracle only if the chain has one, e.g. OP stack L2s and Arbitrum.
	L1Oracle() rollups.L1Oracle
	GetFee(ctx context.Context, calldata []byte, feeLimit uint64, maxFeePrice *assets.Wei, fromAddress, toAddress *common.Address, opts ...feetypes.Opt) (fee EvmFee, estimatedFeeLimit uint64, err error)
	BumpFee(ctx context.Context, originalFee EvmFee, feeLimit uint64, maxFeePrice *assets.Wei, attempts []EvmPriorAttempt) (bumpedFee EvmFee, chainSpecificFeeLimit uint64, err error)

	// GetMaxCost returns the total value = max price x fee units + transferred value
	GetMaxCost(ctx context.Context, amount assets.Eth, calldata []byte, feeLimit uint64, maxFeePrice *assets.Wei, fromAddress, toAddress *common.Address, opts ...feetypes.Opt) (*big.Int, error)
}

type feeEstimatorClient interface {
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	ConfiguredChainID() *big.Int
	HeadByNumber(ctx context.Context, n *big.Int) (*evmtypes.Head, error)
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	FeeHistory(ctx context.Context, blockCount uint64, rewardPercentiles []float64) (feeHistory *ethereum.FeeHistory, err error)
}

// NewEstimator returns the estimator for a given config
func NewEstimator(lggr logger.Logger, ethClient feeEstimatorClient, cfg Config, geCfg evmconfig.GasEstimator) (EvmFeeEstimator, error) {
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
		"estimateLimit", geCfg.EstimateLimit(),
	)
	df := geCfg.EIP1559DynamicFees()

	// create l1Oracle only if it is supported for the chain
	var l1Oracle rollups.L1Oracle
	lggr.Infow("Checking if chain type is roll up", "chainType", cfg.ChainType(), "isRollUp", rollups.IsRollupWithL1Support(cfg.ChainType()))
	if rollups.IsRollupWithL1Support(cfg.ChainType()) {
		var err error
		l1Oracle, err = rollups.NewL1GasOracle(lggr, ethClient, cfg.ChainType())
		if err != nil {
			return nil, fmt.Errorf("failed to initialize L1 oracle: %w", err)
		}
	}
	var newEstimator func(logger.Logger) EvmEstimator
	switch s {
	case "Arbitrum":
		arbOracle, err := rollups.NewArbitrumL1GasOracle(lggr, ethClient)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Arbitrum L1 oracle: %w", err)
		}
		newEstimator = func(l logger.Logger) EvmEstimator {
			return NewArbitrumEstimator(lggr, geCfg, ethClient, arbOracle)
		}
	case "BlockHistory":
		newEstimator = func(l logger.Logger) EvmEstimator {
			return NewBlockHistoryEstimator(lggr, ethClient, cfg, geCfg, bh, ethClient.ConfiguredChainID(), l1Oracle)
		}
	case "FixedPrice":
		newEstimator = func(l logger.Logger) EvmEstimator {
			return NewFixedPriceEstimator(geCfg, ethClient, bh, lggr, l1Oracle)
		}
	case "L2Suggested", "SuggestedPrice":
		newEstimator = func(l logger.Logger) EvmEstimator {
			return NewSuggestedPriceEstimator(lggr, ethClient, geCfg, l1Oracle)
		}
	case "FeeHistory":
		newEstimator = func(l logger.Logger) EvmEstimator {
			ccfg := FeeHistoryEstimatorConfig{
				BumpPercent:      geCfg.BumpPercent(),
				CacheTimeout:     geCfg.FeeHistory().CacheTimeout(),
				EIP1559:          geCfg.EIP1559DynamicFees(),
				BlockHistorySize: uint64(geCfg.BlockHistory().BlockHistorySize()),
				RewardPercentile: float64(geCfg.BlockHistory().TransactionPercentile()),
			}
			return NewFeeHistoryEstimator(lggr, ethClient, ccfg, ethClient.ConfiguredChainID(), l1Oracle)
		}

	default:
		lggr.Warnf("GasEstimator: unrecognised mode '%s', falling back to FixedPriceEstimator", s)
		newEstimator = func(l logger.Logger) EvmEstimator {
			return NewFixedPriceEstimator(geCfg, ethClient, bh, lggr, l1Oracle)
		}
	}
	return NewEvmFeeEstimator(lggr, newEstimator, df, geCfg, ethClient), nil
}

// DynamicFee encompasses both FeeCap and TipCap for EIP1559 transactions
type DynamicFee struct {
	FeeCap *assets.Wei
	TipCap *assets.Wei
}

type EvmPriorAttempt struct {
	ChainSpecificFeeLimit   uint64
	BroadcastBeforeBlockNum *int64
	TxHash                  common.Hash
	TxType                  int
	GasPrice                *assets.Wei
	DynamicFee              DynamicFee
}

// Estimator provides an interface for estimating gas price and limit
type EvmEstimator interface {
	headtracker.HeadTrackable[*evmtypes.Head, common.Hash]
	services.Service

	// GetLegacyGas Calculates initial gas fee for non-EIP1559 transaction
	// maxGasPriceWei parameter is the highest possible gas fee cap that the function will return
	GetLegacyGas(ctx context.Context, calldata []byte, gasLimit uint64, maxGasPriceWei *assets.Wei, opts ...feetypes.Opt) (gasPrice *assets.Wei, chainSpecificGasLimit uint64, err error)
	// BumpLegacyGas Increases gas price and/or limit for non-EIP1559 transactions
	// if the bumped gas fee is greater than maxGasPriceWei, the method returns an error
	// attempts must:
	//   - be sorted in order from highest price to lowest price
	//   - all be of transaction type 0x0 or 0x1
	BumpLegacyGas(ctx context.Context, originalGasPrice *assets.Wei, gasLimit uint64, maxGasPriceWei *assets.Wei, attempts []EvmPriorAttempt) (bumpedGasPrice *assets.Wei, chainSpecificGasLimit uint64, err error)
	// GetDynamicFee Calculates initial gas fee for gas for EIP1559 transactions
	// maxGasPriceWei parameter is the highest possible gas fee cap that the function will return
	GetDynamicFee(ctx context.Context, maxGasPriceWei *assets.Wei) (fee DynamicFee, err error)
	// BumpDynamicFee Increases gas price and/or limit for non-EIP1559 transactions
	// if the bumped gas fee or tip caps are greater than maxGasPriceWei, the method returns an error
	// attempts must:
	//   - be sorted in order from highest price to lowest price
	//   - all be of transaction type 0x2
	BumpDynamicFee(ctx context.Context, original DynamicFee, maxGasPriceWei *assets.Wei, attempts []EvmPriorAttempt) (bumped DynamicFee, err error)

	L1Oracle() rollups.L1Oracle
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

// evmFeeEstimator provides a struct that wraps the EVM specific dynamic and legacy estimators into one estimator that conforms to the generic FeeEstimator
type evmFeeEstimator struct {
	services.StateMachine
	lggr logger.Logger
	EvmEstimator
	EIP1559Enabled bool
	geCfg          GasEstimatorConfig
	ethClient      feeEstimatorClient
}

var _ EvmFeeEstimator = (*evmFeeEstimator)(nil)

func NewEvmFeeEstimator(lggr logger.Logger, newEstimator func(logger.Logger) EvmEstimator, eip1559Enabled bool, geCfg GasEstimatorConfig, ethClient feeEstimatorClient) EvmFeeEstimator {
	lggr = logger.Named(lggr, "WrappedEvmEstimator")
	return &evmFeeEstimator{
		lggr:           lggr,
		EvmEstimator:   newEstimator(lggr),
		EIP1559Enabled: eip1559Enabled,
		geCfg:          geCfg,
		ethClient:      ethClient,
	}
}

func (e *evmFeeEstimator) Name() string {
	return e.lggr.Name()
}

func (e *evmFeeEstimator) Start(ctx context.Context) error {
	return e.StartOnce(e.Name(), func() error {
		if err := e.EvmEstimator.Start(ctx); err != nil {
			return pkgerrors.Wrap(err, "failed to start EVMEstimator")
		}
		l1Oracle := e.L1Oracle()
		if l1Oracle != nil {
			if err := l1Oracle.Start(ctx); err != nil {
				return pkgerrors.Wrap(err, "failed to start L1Oracle")
			}
		}
		return nil
	})
}
func (e *evmFeeEstimator) Close() error {
	return e.StopOnce(e.Name(), func() error {
		var errEVM, errOracle error

		errEVM = pkgerrors.Wrap(e.EvmEstimator.Close(), "failed to stop EVMEstimator")
		l1Oracle := e.L1Oracle()
		if l1Oracle != nil {
			errOracle = pkgerrors.Wrap(l1Oracle.Close(), "failed to stop L1Oracle")
		}

		if errEVM != nil {
			return errEVM
		}
		return errOracle
	})
}

func (e *evmFeeEstimator) Ready() error {
	var errEVM, errOracle error

	errEVM = e.EvmEstimator.Ready()
	l1Oracle := e.L1Oracle()
	if l1Oracle != nil {
		errOracle = l1Oracle.Ready()
	}

	if errEVM != nil {
		return errEVM
	}
	return errOracle
}

func (e *evmFeeEstimator) HealthReport() map[string]error {
	report := map[string]error{e.Name(): e.Healthy()}
	services.CopyHealth(report, e.EvmEstimator.HealthReport())

	l1Oracle := e.L1Oracle()
	if l1Oracle != nil {
		services.CopyHealth(report, l1Oracle.HealthReport())
	}

	return report
}

func (e *evmFeeEstimator) L1Oracle() rollups.L1Oracle {
	return e.EvmEstimator.L1Oracle()
}

// GetFee returns an initial estimated gas price and gas limit for a transaction
// The gas limit provided by the caller can be adjusted by gas estimation or for 2D fees
func (e *evmFeeEstimator) GetFee(ctx context.Context, calldata []byte, feeLimit uint64, maxFeePrice *assets.Wei, fromAddress, toAddress *common.Address, opts ...feetypes.Opt) (fee EvmFee, estimatedFeeLimit uint64, err error) {
	var chainSpecificFeeLimit uint64
	// get dynamic fee
	if e.EIP1559Enabled {
		var dynamicFee DynamicFee
		dynamicFee, err = e.EvmEstimator.GetDynamicFee(ctx, maxFeePrice)
		if err != nil {
			return
		}
		fee.DynamicFeeCap = dynamicFee.FeeCap
		fee.DynamicTipCap = dynamicFee.TipCap
		chainSpecificFeeLimit = feeLimit
	} else {
		// get legacy fee
		fee.Legacy, chainSpecificFeeLimit, err = e.EvmEstimator.GetLegacyGas(ctx, calldata, feeLimit, maxFeePrice, opts...)
		if err != nil {
			return
		}
	}

	estimatedFeeLimit, err = e.estimateFeeLimit(ctx, chainSpecificFeeLimit, calldata, fromAddress, toAddress)
	return
}

func (e *evmFeeEstimator) GetMaxCost(ctx context.Context, amount assets.Eth, calldata []byte, feeLimit uint64, maxFeePrice *assets.Wei, fromAddress, toAddress *common.Address, opts ...feetypes.Opt) (*big.Int, error) {
	fees, gasLimit, err := e.GetFee(ctx, calldata, feeLimit, maxFeePrice, fromAddress, toAddress, opts...)
	if err != nil {
		return nil, err
	}

	var gasPrice *assets.Wei
	if e.EIP1559Enabled {
		gasPrice = fees.DynamicFeeCap
	} else {
		gasPrice = fees.Legacy
	}

	fee := new(big.Int).Mul(gasPrice.ToInt(), big.NewInt(int64(gasLimit)))
	amountWithFees := new(big.Int).Add(amount.ToInt(), fee)
	return amountWithFees, nil
}

func (e *evmFeeEstimator) BumpFee(ctx context.Context, originalFee EvmFee, feeLimit uint64, maxFeePrice *assets.Wei, attempts []EvmPriorAttempt) (bumpedFee EvmFee, chainSpecificFeeLimit uint64, err error) {
	// validate only 1 fee type is present
	if (!originalFee.ValidDynamic() && originalFee.Legacy == nil) || (originalFee.ValidDynamic() && originalFee.Legacy != nil) {
		err = pkgerrors.New("only one dynamic or legacy fee can be defined")
		return
	}

	// bump fee based on what fee the tx has previously used (not based on config)
	// bump dynamic original
	if originalFee.ValidDynamic() {
		var bumpedDynamic DynamicFee
		bumpedDynamic, err = e.EvmEstimator.BumpDynamicFee(ctx,
			DynamicFee{
				TipCap: originalFee.DynamicTipCap,
				FeeCap: originalFee.DynamicFeeCap,
			}, maxFeePrice, attempts)
		if err != nil {
			return
		}
		chainSpecificFeeLimit, err = commonfee.ApplyMultiplier(feeLimit, e.geCfg.LimitMultiplier())
		bumpedFee.DynamicFeeCap = bumpedDynamic.FeeCap
		bumpedFee.DynamicTipCap = bumpedDynamic.TipCap
		return
	}

	// bump legacy fee
	bumpedFee.Legacy, chainSpecificFeeLimit, err = e.EvmEstimator.BumpLegacyGas(ctx, originalFee.Legacy, feeLimit, maxFeePrice, attempts)
	if err != nil {
		return
	}
	chainSpecificFeeLimit, err = commonfee.ApplyMultiplier(chainSpecificFeeLimit, e.geCfg.LimitMultiplier())
	return
}

func (e *evmFeeEstimator) estimateFeeLimit(ctx context.Context, feeLimit uint64, calldata []byte, fromAddress, toAddress *common.Address) (estimatedFeeLimit uint64, err error) {
	// Use the feeLimit * LimitMultiplier as the provided gas limit since this multiplier is applied on top of the caller specified gas limit
	providedGasLimit, err := commonfee.ApplyMultiplier(feeLimit, e.geCfg.LimitMultiplier())
	if err != nil {
		return estimatedFeeLimit, err
	}
	// Use provided fee limit by default if EstimateLimit is disabled
	if !e.geCfg.EstimateLimit() {
		return providedGasLimit, nil
	}
	// Create call msg for gas limit estimation
	// Skip setting Gas to avoid capping the results of the estimation
	callMsg := ethereum.CallMsg{
		To:   toAddress,
		Data: calldata,
	}
	if fromAddress != nil {
		callMsg.From = *fromAddress
	}
	estimatedGas, estimateErr := e.ethClient.EstimateGas(ctx, callMsg)
	if estimateErr != nil {
		if providedGasLimit > 0 {
			// Do not return error if estimate gas failed, we can still use the provided limit instead since it is an upper limit
			e.lggr.Errorw("failed to estimate gas limit. falling back to the provided gas limit with multiplier", "callMsg", callMsg, "providedGasLimitWithMultiplier", providedGasLimit, "error", estimateErr)
			return providedGasLimit, nil
		}
		return estimatedFeeLimit, fmt.Errorf("gas estimation failed and provided gas limit is 0: %w", estimateErr)
	}
	e.lggr.Debugw("estimated gas", "estimatedGas", estimatedGas, "providedGasLimitWithMultiplier", providedGasLimit)
	// Return error if estimated gas without the buffer exceeds the provided gas limit, if provided
	// Transaction would be destined to run out of gas and fail
	if providedGasLimit > 0 && estimatedGas > providedGasLimit {
		e.lggr.Errorw("estimated gas exceeds provided gas limit with multiplier", "estimatedGas", estimatedGas, "providedGasLimitWithMultiplier", providedGasLimit)
		return estimatedFeeLimit, commonfee.ErrFeeLimitTooLow
	}
	// Apply EstimateGasBuffer to the estimated gas limit
	estimatedFeeLimit, err = commonfee.ApplyMultiplier(estimatedGas, EstimateGasBuffer)
	if err != nil {
		return
	}
	// If provided gas limit is not 0, fallback to it if the buffer causes the estimated gas limit to exceed it
	// The provided gas limit should be used as an upper bound to avoid unexpected behavior for products
	if providedGasLimit > 0 && estimatedFeeLimit > providedGasLimit {
		e.lggr.Debugw("estimated gas limit with buffer exceeds the provided gas limit with multiplier. falling back to the provided gas limit with multiplier", "estimatedGasLimit", estimatedFeeLimit, "providedGasLimitWithMultiplier", providedGasLimit)
		estimatedFeeLimit = providedGasLimit
	}

	return
}

// Config defines an interface for configuration in the gas package
type Config interface {
	ChainType() chaintype.ChainType
	FinalityDepth() uint32
	FinalityTagEnabled() bool
}

type GasEstimatorConfig interface {
	EIP1559DynamicFees() bool
	BumpPercent() uint16
	BumpThreshold() uint64
	BumpMin() *assets.Wei
	FeeCapDefault() *assets.Wei
	LimitMax() uint64
	LimitMultiplier() float32
	PriceDefault() *assets.Wei
	TipCapDefault() *assets.Wei
	TipCapMin() *assets.Wei
	PriceMin() *assets.Wei
	PriceMax() *assets.Wei
	Mode() string
	EstimateLimit() bool
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

// BumpLegacyGasPriceOnly will increase the price
func BumpLegacyGasPriceOnly(cfg bumpConfig, lggr logger.SugaredLogger, currentGasPrice, originalGasPrice *assets.Wei, maxGasPriceWei *assets.Wei) (gasPrice *assets.Wei, err error) {
	gasPrice, err = bumpGasPrice(cfg, lggr, currentGasPrice, originalGasPrice, maxGasPriceWei)
	if err != nil {
		return nil, err
	}
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
		return maxGasPrice, pkgerrors.Wrapf(commonfee.ErrBumpFeeExceedsLimit, "bumped gas price of %s would exceed configured max gas price of %s (original price was %s). %s",
			bumpedGasPrice.String(), maxGasPrice, originalGasPrice.String(), label.NodeConnectivityProblemWarning)
	} else if bumpedGasPrice.Cmp(originalGasPrice) == 0 {
		// NOTE: This really shouldn't happen since we enforce minimums for
		// EVM.GasEstimator.BumpPercent and EVM.GasEstimator.BumpMin in the config validation,
		// but it's here anyway for a "belts and braces" approach
		return bumpedGasPrice, pkgerrors.Wrapf(commonfee.ErrBump, "bumped gas price of %s is equal to original gas price of %s."+
			" ACTION REQUIRED: This is a configuration error, you must increase either "+
			"EVM.GasEstimator.BumpPercent or EVM.GasEstimator.BumpMin", bumpedGasPrice.String(), originalGasPrice.String())
	}
	return bumpedGasPrice, nil
}

// BumpDynamicFeeOnly bumps the tip cap and max gas price if necessary
func BumpDynamicFeeOnly(config bumpConfig, feeCapBufferBlocks uint16, lggr logger.SugaredLogger, currentTipCap, currentBaseFee *assets.Wei, originalFee DynamicFee, maxGasPriceWei *assets.Wei) (bumped DynamicFee, err error) {
	bumped, err = bumpDynamicFee(config, feeCapBufferBlocks, lggr, currentTipCap, currentBaseFee, originalFee, maxGasPriceWei)
	if err != nil {
		return bumped, err
	}
	return
}

// bumpDynamicFee computes the next tip cap to attempt as the largest of:
// - A configured percentage bump (EVM.GasEstimator.BumpPercent) on top of the baseline tip cap.
// - A configured fixed amount of Wei (ETH_GAS_PRICE_WEI) on top of the baseline tip cap.
// The baseline tip cap is the maximum of the previous tip cap attempt and the node's current tip cap.
// It increases the max fee cap by BumpPercent
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
		return bumpedFee, pkgerrors.Wrapf(commonfee.ErrBumpFeeExceedsLimit, "bumped tip cap of %s would exceed configured max gas price of %s (original fee: tip cap %s, fee cap %s). %s",
			bumpedTipCap.String(), maxGasPrice, originalFee.TipCap.String(), originalFee.FeeCap.String(), label.NodeConnectivityProblemWarning)
	} else if bumpedTipCap.Cmp(originalFee.TipCap) <= 0 {
		// NOTE: This really shouldn't happen since we enforce minimums for
		// EVM.GasEstimator.BumpPercent and EVM.GasEstimator.BumpMin in the config validation,
		// but it's here anyway for a "belts and braces" approach
		return bumpedFee, pkgerrors.Wrapf(commonfee.ErrBump, "bumped gas tip cap of %s is less than or equal to original gas tip cap of %s."+
			" ACTION REQUIRED: This is a configuration error, you must increase either "+
			"EVM.GasEstimator.BumpPercent or EVM.GasEstimator.BumpMin", bumpedTipCap.String(), originalFee.TipCap.String())
	}

	// Always bump the FeeCap by at least the bump percentage (should be greater than or
	// equal to than geth's configured bump minimum which is 10%)
	// See: https://github.com/ethereum/go-ethereum/blob/bff330335b94af3643ac2fb809793f77de3069d4/core/tx_list.go#L298
	bumpedFeeCap := bumpFeePrice(originalFee.FeeCap, cfg.BumpPercent(), cfg.BumpMin())

	if currentBaseFee != nil {
		if currentBaseFee.Cmp(maxGasPrice) > 0 {
			lggr.Warnf("Ignoring current base fee of %s which is greater than max gas price of %s", currentBaseFee.String(), maxGasPrice.String())
		} else {
			currentFeeCap := calcFeeCap(currentBaseFee, int(feeCapBufferBlocks), bumpedTipCap, maxGasPrice)
			bumpedFeeCap = assets.WeiMax(bumpedFeeCap, currentFeeCap)
		}
	}

	if bumpedFeeCap.Cmp(maxGasPrice) > 0 {
		return bumpedFee, pkgerrors.Wrapf(commonfee.ErrBumpFeeExceedsLimit, "bumped fee cap of %s would exceed configured max gas price of %s (original fee: tip cap %s, fee cap %s). %s",
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
	return assets.NewWei(bigmath.Min(userSpecifiedMax.ToInt(), maxGasPriceWei.ToInt()))
}

func capGasPrice(calculatedGasPrice, userSpecifiedMax, maxGasPriceWei *assets.Wei) *assets.Wei {
	maxGasPrice := commonfee.CalculateFee(calculatedGasPrice.ToInt(), userSpecifiedMax.ToInt(), maxGasPriceWei.ToInt())
	return assets.NewWei(maxGasPrice)
}
