package gas

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

var (
	promNumGasBumps = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tx_manager_num_gas_bumps",
		Help: "Number of gas bumps",
	})

	promGasBumpExceedsLimit = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tx_manager_gas_bump_exceeds_limit",
		Help: "Number of times gas bumping failed from exceeding the configured limit. Any counts of this type indicate a serious problem.",
	})
)

func NewEstimator(ethClient eth.Client, config Config) Estimator {
	s := config.GasEstimatorMode()
	switch s {
	case "BlockHistory":
		return NewBlockHistoryEstimator(ethClient, config)
	case "FixedPrice":
		return NewFixedPriceEstimator(config)
	case "Optimism":
		return NewOptimismEstimator(config, ethClient)
	default:
		logger.Warnf("GasEstimator: unrecognised mode '%s', falling back to FixedPriceEstimator", s)
		return NewFixedPriceEstimator(config)
	}
}

// Estimator provides an interface for estimating gas price and limit
type Estimator interface {
	OnNewLongestChain(context.Context, models.Head)
	Start() error
	Close() error
	EstimateGas(calldata []byte, gasLimit uint64, opts ...Opt) (gasPrice *big.Int, chainSpecificGasLimit uint64, err error)
	BumpGas(originalGasPrice *big.Int, gasLimit uint64) (bumpedGasPrice *big.Int, chainSpecificGasLimit uint64, err error)
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
	ChainID() *big.Int
	EthFinalityDepth() uint
	EthGasBumpPercent() uint16
	EthGasBumpWei() *big.Int
	EthGasLimitMultiplier() float32
	EthGasPriceDefault() *big.Int
	EthMaxGasPriceWei() *big.Int
	EthMinGasPriceWei() *big.Int
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
	Number       int64
	Hash         common.Hash
	ParentHash   common.Hash
	Transactions []Transaction
}

type blockInternal struct {
	Number       string
	Hash         common.Hash
	ParentHash   common.Hash
	Transactions []Transaction
}

// MarshalJSON implements json marshalling for Block
func (b Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(blockInternal{
		Int64ToHex(b.Number),
		b.Hash,
		b.ParentHash,
		b.Transactions,
	})
}

// UnmarshalJSON unmarshals to a Block
func (b *Block) UnmarshalJSON(data []byte) error {
	bi := blockInternal{}
	if err := json.Unmarshal(data, &bi); err != nil {
		return errors.Wrapf(err, "failed to unmarshal to blockInternal, got: '%s'", data)
	}
	n, err := hexutil.DecodeBig(bi.Number)
	if err != nil {
		return errors.Wrapf(err, "failed to decode block number while unmarshalling block, got: '%s'", data)
	}
	*b = Block{
		n.Int64(),
		bi.Hash,
		bi.ParentHash,
		bi.Transactions,
	}
	return nil
}

type transactionInternal struct {
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Gas      *hexutil.Uint64 `json:"gas"`
}

// Transaction represents an ethereum transaction
// Use our own type because geth's type has validation failures on e.g. zero
// gas used, which can occur on other chains.
// This type is only used for the block history estimator, and can be expensive to unmarshal. Don't add unnecessary fields here.
type Transaction struct {
	GasPrice *big.Int
	GasLimit uint64
}

// UnmarshalJSON unmarshals a Transaction
func (t *Transaction) UnmarshalJSON(data []byte) error {
	ti := transactionInternal{}
	if err := json.Unmarshal(data, &ti); err != nil {
		return errors.Wrapf(err, "failed to unmarshal to transactionInternal, got: '%s'", data)
	}
	if ti.Gas == nil {
		return errors.Errorf("expected 'gas' to not be null, got: '%s'", data)
	}
	*t = Transaction{
		(*big.Int)(ti.GasPrice),
		uint64(*ti.Gas),
	}
	return nil
}

// BumpGasPriceOnly will increase the price and apply multiplier to the gas limit
func BumpGasPriceOnly(config Config, originalGasPrice *big.Int, originalGasLimit uint64) (gasPrice *big.Int, chainSpecificGasLimit uint64, err error) {
	gasPrice, err = bumpGasPrice(config, originalGasPrice)
	if err != nil {
		return nil, 0, err
	}
	chainSpecificGasLimit = applyMultiplier(originalGasLimit, config.EthGasLimitMultiplier())
	return
}

// bumpGasPrice computes the next gas price to attempt as the largest of:
// - A configured percentage bump (ETH_GAS_BUMP_PERCENT) on top of the baseline price.
// - A configured fixed amount of Wei (ETH_GAS_PRICE_WEI) on top of the baseline price.
// The baseline price is the maximum of the previous gas price attempt and the node's current gas price.
func bumpGasPrice(config Config, originalGasPrice *big.Int) (*big.Int, error) {
	baselinePrice := max(originalGasPrice, config.EthGasPriceDefault())

	var priceByPercentage = new(big.Int)
	priceByPercentage.Mul(baselinePrice, big.NewInt(int64(100+config.EthGasBumpPercent())))
	priceByPercentage.Div(priceByPercentage, big.NewInt(100))

	var priceByIncrement = new(big.Int)
	priceByIncrement.Add(baselinePrice, config.EthGasBumpWei())

	bumpedGasPrice := max(priceByPercentage, priceByIncrement)
	if bumpedGasPrice.Cmp(config.EthMaxGasPriceWei()) > 0 {
		promGasBumpExceedsLimit.Inc()
		return config.EthMaxGasPriceWei(), errors.Errorf("bumped gas price of %s would exceed configured max gas price of %s (original price was %s). %s",
			bumpedGasPrice.String(), config.EthMaxGasPriceWei(), originalGasPrice.String(), static.EthNodeConnectivityProblemLabel)
	} else if bumpedGasPrice.Cmp(originalGasPrice) == 0 {
		// NOTE: This really shouldn't happen since we enforce minimums for
		// ETH_GAS_BUMP_PERCENT and ETH_GAS_BUMP_WEI in the config validation,
		// but it's here anyway for a "belts and braces" approach
		return bumpedGasPrice, errors.Errorf("bumped gas price of %s is equal to original gas price of %s."+
			" ACTION REQUIRED: This is a configuration error, you must increase either "+
			"ETH_GAS_BUMP_PERCENT or ETH_GAS_BUMP_WEI", bumpedGasPrice.String(), originalGasPrice.String())
	}
	// TODO: Move/fix these
	promNumGasBumps.Inc()
	return bumpedGasPrice, nil
}

func max(a, b *big.Int) *big.Int {
	if a.Cmp(b) >= 0 {
		return a
	}
	return b
}
