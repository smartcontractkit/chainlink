package gas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/label"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	bigmath "github.com/smartcontractkit/chainlink/core/utils/big_math"
)

var (
	ErrBumpGasExceedsLimit = errors.New("gas bump exceeds limit")
	ErrBump                = errors.New("gas bump failed")
)

func IsBumpErr(err error) bool {
	return err != nil && (errors.Is(err, ErrBumpGasExceedsLimit) || errors.Is(err, ErrBump))
}

// NewEstimator returns the estimator for a given config
func NewEstimator(lggr logger.Logger, ethClient evmclient.Client, cfg Config) Estimator {
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
	case "BlockHistory":
		return NewBlockHistoryEstimator(lggr, ethClient, cfg, *ethClient.ChainID())
	case "FixedPrice":
		return NewFixedPriceEstimator(cfg, lggr)
	case "Optimism2", "L2Suggested":
		return NewL2SuggestedEstimator(lggr, cfg, ethClient)
	default:
		lggr.Warnf("GasEstimator: unrecognised mode '%s', falling back to FixedPriceEstimator", s)
		return NewFixedPriceEstimator(cfg, lggr)
	}
}

// DynamicFee encompasses both FeeCap and TipCap for EIP1559 transactions
type DynamicFee struct {
	FeeCap *big.Int
	TipCap *big.Int
}

// Estimator provides an interface for estimating gas price and limit
//go:generate mockery --name Estimator --output ./mocks/ --case=underscore
type Estimator interface {
	OnNewLongestChain(context.Context, *evmtypes.Head)
	Start(context.Context) error
	Close() error
	// Calculates initial gas fee for non-EIP1559 transaction
	// maxGasPriceWei parameter is the highest possible gas fee cap that the function will return
	GetLegacyGas(calldata []byte, gasLimit uint64, maxGasPriceWei *big.Int, opts ...Opt) (gasPrice *big.Int, chainSpecificGasLimit uint64, err error)
	// Increases gas price and/or limit for non-EIP1559 transactions
	// if the bumped gas fee is greater than maxGasPriceWei, the method returns an error
	BumpLegacyGas(originalGasPrice *big.Int, gasLimit uint64, maxGasPriceWei *big.Int) (bumpedGasPrice *big.Int, chainSpecificGasLimit uint64, err error)
	// Calculates initial gas fee for gas for EIP1559 transactions
	// maxGasPriceWei parameter is the highest possible gas fee cap that the function will return
	GetDynamicFee(gasLimit uint64, maxGasPriceWei *big.Int) (fee DynamicFee, chainSpecificGasLimit uint64, err error)
	// Increases gas price and/or limit for non-EIP1559 transactions
	// if the bumped gas fee or tip caps are greater than maxGasPriceWei, the method returns an error
	BumpDynamicFee(original DynamicFee, gasLimit uint64, maxGasPriceWei *big.Int) (bumped DynamicFee, chainSpecificGasLimit uint64, err error)
}

// Opt is an option for a gas estimator
type Opt int

const (
	// OptForceRefetch forces the estimator to bust a cache if necessary
	OptForceRefetch Opt = iota
)

func applyMultiplier(gasLimit uint64, multiplier float32) uint64 {
	return uint64(decimal.NewFromBigInt(big.NewInt(0).SetUint64(gasLimit), 0).Mul(decimal.NewFromFloat32(multiplier)).IntPart())
}

// Config defines an interface for configuration in the gas package
//go:generate mockery --name Config --output ./mocks/ --case=underscore
type Config interface {
	BlockHistoryEstimatorBatchSize() uint32
	BlockHistoryEstimatorBlockDelay() uint16
	BlockHistoryEstimatorBlockHistorySize() uint16
	BlockHistoryEstimatorTransactionPercentile() uint16
	BlockHistoryEstimatorEIP1559FeeCapBufferBlocks() uint16
	ChainType() config.ChainType
	EvmEIP1559DynamicFees() bool
	EvmFinalityDepth() uint32
	EvmGasBumpPercent() uint16
	EvmGasBumpThreshold() uint64
	EvmGasBumpWei() *big.Int
	EvmGasFeeCapDefault() *big.Int
	EvmGasLimitMultiplier() float32
	EvmGasPriceDefault() *big.Int
	EvmGasTipCapDefault() *big.Int
	EvmGasTipCapMinimum() *big.Int
	EvmMaxGasPriceWei() *big.Int
	EvmMinGasPriceWei() *big.Int
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

// Block represents an ethereum block
// This type is only used for the block history estimator, and can be expensive to unmarshal. Don't add unnecessary fields here.
type Block struct {
	Number        int64
	Hash          common.Hash
	ParentHash    common.Hash
	BaseFeePerGas *big.Int
	Timestamp     time.Time
	Transactions  []Transaction
}

type blockInternal struct {
	Number        string
	Hash          common.Hash
	ParentHash    common.Hash
	BaseFeePerGas *hexutil.Big
	Timestamp     hexutil.Uint64
	Transactions  []Transaction
}

// MarshalJSON implements json marshalling for Block
func (b Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(blockInternal{
		Int64ToHex(b.Number),
		b.Hash,
		b.ParentHash,
		(*hexutil.Big)(b.BaseFeePerGas),
		(hexutil.Uint64)(uint64(b.Timestamp.Unix())),
		b.Transactions,
	})
}

var ErrMissingBlock = errors.New("missing block")

// UnmarshalJSON unmarshals to a Block
func (b *Block) UnmarshalJSON(data []byte) error {
	var bi *blockInternal
	if err := json.Unmarshal(data, &bi); err != nil {
		return errors.Wrapf(err, "failed to unmarshal to blockInternal, got: '%s'", data)
	}
	if bi == nil {
		return errors.WithStack(ErrMissingBlock)
	}
	n, err := hexutil.DecodeBig(bi.Number)
	if err != nil {
		return errors.Wrapf(err, "failed to decode block number while unmarshalling block, got: '%s'", data)
	}
	*b = Block{
		n.Int64(),
		bi.Hash,
		bi.ParentHash,
		(*big.Int)(bi.BaseFeePerGas),
		time.Unix((int64((uint64)(bi.Timestamp))), 0),
		bi.Transactions,
	}
	return nil
}

type TxType uint8

// NOTE: Need to roll our own unmarshaller since geth's hexutil.Uint64 does not
// handle double zeroes e.g. 0x00
func (txt *TxType) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte(`"0x00"`)) {
		data = []byte(`"0x0"`)
	}
	var hx hexutil.Uint64
	if err := (&hx).UnmarshalJSON(data); err != nil {
		return err
	}
	if hx > math.MaxUint8 {
		return errors.Errorf("expected 'type' to fit into a single byte, got: '%s'", data)
	}
	*txt = TxType(hx)
	return nil
}

type transactionInternal struct {
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Type                 *TxType         `json:"type"`
	Hash                 common.Hash     `json:"hash"`
}

// Transaction represents an ethereum transaction
// Use our own type because geth's type has validation failures on e.g. zero
// gas used, which can occur on other chains.
// This type is only used for the block history estimator, and can be expensive to unmarshal. Don't add unnecessary fields here.
type Transaction struct {
	GasPrice             *big.Int
	GasLimit             uint64
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
	Type                 TxType
	Hash                 common.Hash
}

const LegacyTxType = TxType(0x0)

// UnmarshalJSON unmarshals a Transaction
func (t *Transaction) UnmarshalJSON(data []byte) error {
	ti := transactionInternal{}
	if err := json.Unmarshal(data, &ti); err != nil {
		return errors.Wrapf(err, "failed to unmarshal to transactionInternal, got: '%s'", data)
	}
	if ti.Gas == nil {
		return errors.Errorf("expected 'gas' to not be null, got: '%s'", data)
	}
	if ti.Type == nil {
		tpe := LegacyTxType
		ti.Type = &tpe
	}
	*t = Transaction{
		(*big.Int)(ti.GasPrice),
		uint64(*ti.Gas),
		(*big.Int)(ti.MaxFeePerGas),
		(*big.Int)(ti.MaxPriorityFeePerGas),
		*ti.Type,
		ti.Hash,
	}
	return nil
}

// BumpLegacyGasPriceOnly will increase the price and apply multiplier to the gas limit
func BumpLegacyGasPriceOnly(cfg Config, lggr logger.SugaredLogger, currentGasPrice, originalGasPrice *big.Int, originalGasLimit uint64, maxGasPriceWei *big.Int) (gasPrice *big.Int, chainSpecificGasLimit uint64, err error) {
	gasPrice, err = bumpGasPrice(cfg, lggr, currentGasPrice, originalGasPrice, maxGasPriceWei)
	if err != nil {
		return nil, 0, err
	}
	chainSpecificGasLimit = applyMultiplier(originalGasLimit, cfg.EvmGasLimitMultiplier())
	return
}

// bumpGasPrice computes the next gas price to attempt as the largest of:
// - A configured percentage bump (ETH_GAS_BUMP_PERCENT) on top of the baseline price.
// - A configured fixed amount of Wei (ETH_GAS_PRICE_WEI) on top of the baseline price.
// The baseline price is the maximum of the previous gas price attempt and the node's current gas price.
func bumpGasPrice(cfg Config, lggr logger.SugaredLogger, currentGasPrice, originalGasPrice *big.Int, maxGasPriceWei *big.Int) (*big.Int, error) {
	maxGasPrice := getMaxGasPrice(maxGasPriceWei, cfg)
	var priceByPercentage = new(big.Int)
	priceByPercentage.Mul(originalGasPrice, big.NewInt(int64(100+cfg.EvmGasBumpPercent())))
	priceByPercentage.Div(priceByPercentage, big.NewInt(100))

	var priceByIncrement = new(big.Int)
	priceByIncrement.Add(originalGasPrice, cfg.EvmGasBumpWei())

	bumpedGasPrice := bigmath.Max(priceByPercentage, priceByIncrement)
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
		// ETH_GAS_BUMP_PERCENT and ETH_GAS_BUMP_WEI in the config validation,
		// but it's here anyway for a "belts and braces" approach
		return bumpedGasPrice, errors.Wrapf(ErrBump, "bumped gas price of %s is equal to original gas price of %s."+
			" ACTION REQUIRED: This is a configuration error, you must increase either "+
			"ETH_GAS_BUMP_PERCENT or ETH_GAS_BUMP_WEI", bumpedGasPrice.String(), originalGasPrice.String())
	}
	return bumpedGasPrice, nil
}

// BumpDynamicFeeOnly bumps the tip cap and max gas price if necessary
func BumpDynamicFeeOnly(config Config, lggr logger.SugaredLogger, currentTipCap *big.Int, currentBaseFee *big.Int, originalFee DynamicFee, originalGasLimit uint64, maxGasPriceWei *big.Int) (bumped DynamicFee, chainSpecificGasLimit uint64, err error) {
	bumped, err = bumpDynamicFee(config, lggr, currentTipCap, currentBaseFee, originalFee, maxGasPriceWei)
	if err != nil {
		return bumped, 0, err
	}
	chainSpecificGasLimit = applyMultiplier(originalGasLimit, config.EvmGasLimitMultiplier())
	return
}

// bumpDynamicFee computes the next tip cap to attempt as the largest of:
// - A configured percentage bump (ETH_GAS_BUMP_PERCENT) on top of the baseline tip cap.
// - A configured fixed amount of Wei (ETH_GAS_PRICE_WEI) on top of the baseline tip cap.
// The baseline tip cap is the maximum of the previous tip cap attempt and the node's current tip cap.
// It increases the max fee cap by GasBumpPercent
//
// NOTE: We would prefer to have set a large FeeCap and leave it fixed, bumping
// the Tip only. Unfortunately due to a flaw of how EIP-1559 is implemented we
// have to bump FeeCap by at least 10% each time we bump the tip cap.
// See: https://github.com/ethereum/go-ethereum/issues/24284
func bumpDynamicFee(cfg Config, lggr logger.SugaredLogger, currentTipCap, currentBaseFee *big.Int, originalFee DynamicFee, maxGasPriceWei *big.Int) (bumpedFee DynamicFee, err error) {
	maxGasPrice := getMaxGasPrice(maxGasPriceWei, cfg)
	baselineTipCap := bigmath.Max(originalFee.TipCap, cfg.EvmGasTipCapDefault())

	bumpedTipCap := increaseByPercentageOrIncrement(baselineTipCap, cfg.EvmGasBumpPercent(), cfg.EvmGasBumpWei())

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
		// ETH_GAS_BUMP_PERCENT and ETH_GAS_BUMP_WEI in the config validation,
		// but it's here anyway for a "belts and braces" approach
		return bumpedFee, errors.Wrapf(ErrBump, "bumped gas tip cap of %s is less than or equal to original gas tip cap of %s."+
			" ACTION REQUIRED: This is a configuration error, you must increase either "+
			"ETH_GAS_BUMP_PERCENT or ETH_GAS_BUMP_WEI", bumpedTipCap.String(), originalFee.TipCap.String())
	}

	// Always bump the FeeCap by at least the bump percentage (should be greater than or
	// equal to than geth's configured bump minimum which is 10%)
	// See: https://github.com/ethereum/go-ethereum/blob/bff330335b94af3643ac2fb809793f77de3069d4/core/tx_list.go#L298
	bumpedFeeCap := increaseByPercentageOrIncrement(originalFee.FeeCap, cfg.EvmGasBumpPercent(), cfg.EvmGasBumpWei())

	if currentBaseFee != nil {
		if currentBaseFee.Cmp(maxGasPrice) > 0 {
			lggr.Warnf("Ignoring current base fee of %s which is greater than max gas price of %s", currentBaseFee.String(), maxGasPrice.String())
		} else {
			currentFeeCap := calcFeeCap(currentBaseFee, cfg, bumpedTipCap, maxGasPrice)
			bumpedFeeCap = bigmath.Max(bumpedFeeCap, currentFeeCap)
		}
	}

	if bumpedFeeCap.Cmp(maxGasPrice) > 0 {
		return bumpedFee, errors.Wrapf(ErrBumpGasExceedsLimit, "bumped fee cap of %s would exceed configured max gas price of %s (original fee: tip cap %s, fee cap %s). %s",
			bumpedFeeCap.String(), maxGasPrice, originalFee.TipCap.String(), originalFee.FeeCap.String(), label.NodeConnectivityProblemWarning)
	}

	return DynamicFee{FeeCap: bumpedFeeCap, TipCap: bumpedTipCap}, nil
}

// Returns whichever is greater, the percentage bump or the bump by fixed increment
func increaseByPercentageOrIncrement(original *big.Int, percentage uint16, increment *big.Int) (bumped *big.Int) {
	percentageBump := increaseByPercentage(original, percentage)

	incrementBump := new(big.Int).Add(original, increment)

	return bigmath.Max(percentageBump, incrementBump)
}

func increaseByPercentage(original *big.Int, percentage uint16) (bumped *big.Int) {
	bumped = new(big.Int)
	bumped.Set(original)
	bumped.Mul(original, big.NewInt(int64(100+percentage)))
	bumped.Div(bumped, big.NewInt(100))
	return
}

func getMaxGasPrice(userSpecifiedMax *big.Int, config Config) *big.Int {
	return bigmath.Min(config.EvmMaxGasPriceWei(), userSpecifiedMax)
}

func capGasPrice(calculatedGasPrice, userSpecifiedMax *big.Int, config Config) *big.Int {
	maxGasPrice := getMaxGasPrice(userSpecifiedMax, config)
	return bigmath.Min(calculatedGasPrice, maxGasPrice)
}
