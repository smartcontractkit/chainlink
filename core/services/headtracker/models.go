package headtracker

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

// Block represents an ethereum block
// This type is only used for the block fetcher, and can be expensive to unmarshal. Don't add unnecessary fields here.
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

// Transaction represents an ethereum transaction
// Use our own type because geth's type has validation failures on e.g. zero
// gas used, which can occur on other chains.
// This type is only used for the block fetcher, and can be expensive to unmarshal. Don't add unnecessary fields here.
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

type BlockFetcherFakeConfig struct {
	GasUpdaterTransactionPercentileField uint16
	EthMaxGasPriceWeiField               *big.Int
	EthMinGasPriceWeiField               *big.Int
	EthFinalityDepthField                uint
	BlockBackfillDepthField              uint64
	BlockFetcherBatchSizeField           uint32
	EthHeadTrackerHistoryDepthField      uint
	GasUpdaterBatchSizeField             uint32
	GasUpdaterBlockDelayField            uint16
	GasUpdaterBlockHistorySizeField      uint16
	ChainIDField                         *big.Int
}

func NewBlockFetcherConfigWithDefaults() *BlockFetcherFakeConfig {
	return &BlockFetcherFakeConfig{
		EthFinalityDepthField:           42,
		BlockBackfillDepthField:         50,
		BlockFetcherBatchSizeField:      2,
		EthHeadTrackerHistoryDepthField: 100,
		GasUpdaterBatchSizeField:        0,
		GasUpdaterBlockDelayField:       0,
		GasUpdaterBlockHistorySizeField: 2,
		ChainIDField:                    big.NewInt(0),
	}
}

func (config BlockFetcherFakeConfig) GasUpdaterTransactionPercentile() uint16 {
	return config.GasUpdaterTransactionPercentileField
}

func (config BlockFetcherFakeConfig) EthMaxGasPriceWei() *big.Int {
	return config.EthMaxGasPriceWeiField
}

func (config BlockFetcherFakeConfig) EthMinGasPriceWei() *big.Int {
	return config.EthMinGasPriceWeiField
}

func (config BlockFetcherFakeConfig) ChainID() *big.Int {
	return config.ChainIDField
}

func (config BlockFetcherFakeConfig) SetEthGasPriceDefault(value *big.Int) error {
	*config.EthMaxGasPriceWeiField = *value
	return nil
}

// EthFinalityDepth provides a mock function with given fields:
func (config BlockFetcherFakeConfig) EthFinalityDepth() uint {
	return config.EthFinalityDepthField
}

// BlockBackfillDepth provides a mock function with given fields:
func (config BlockFetcherFakeConfig) BlockBackfillDepth() uint64 {
	return config.BlockBackfillDepthField
}

// BlockFetcherBatchSize provides a mock function with given fields:
func (config BlockFetcherFakeConfig) BlockFetcherBatchSize() uint32 {
	return config.BlockFetcherBatchSizeField
}

// EthHeadTrackerHistoryDepth provides a mock function with given fields:
func (config BlockFetcherFakeConfig) EthHeadTrackerHistoryDepth() uint {
	return config.EthHeadTrackerHistoryDepthField
}

// GasUpdaterBatchSize provides a mock function with given fields:
func (config BlockFetcherFakeConfig) GasUpdaterBatchSize() uint32 {
	return config.GasUpdaterBatchSizeField
}

// GasUpdaterBlockDelay provides a mock function with given fields:
func (config BlockFetcherFakeConfig) GasUpdaterBlockDelay() uint16 {
	return config.GasUpdaterBlockDelayField
}

// GasUpdaterBlockHistorySize provides a mock function with given fields:
func (config BlockFetcherFakeConfig) GasUpdaterBlockHistorySize() uint16 {
	return config.GasUpdaterBlockHistorySizeField
}
