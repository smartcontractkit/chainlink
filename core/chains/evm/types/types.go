package types

import (
	"database/sql/driver"
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
type NewNode struct {
	Name       string      `json:"name"`
	EVMChainID utils.Big   `json:"evmChainId"`
	WSURL      null.String `json:"wsURL" db:"ws_url"`
	HTTPURL    null.String `json:"httpURL" db:"http_url"`
	SendOnly   bool        `json:"sendOnly"`
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
type ChainConfigORM interface {
	StoreString(chainID utils.Big, key, val string) error
	Clear(chainID utils.Big, key string) error
}

type ORM interface {
	Chain(id utils.Big, qopts ...pg.QOpt) (chain DBChain, err error)
	Chains(offset, limit int, qopts ...pg.QOpt) ([]DBChain, int, error)
	CreateChain(id utils.Big, config *ChainCfg, qopts ...pg.QOpt) (DBChain, error)
	UpdateChain(id utils.Big, enabled bool, config *ChainCfg, qopts ...pg.QOpt) (DBChain, error)
	DeleteChain(id utils.Big, qopts ...pg.QOpt) error
	GetChainsByIDs(ids []utils.Big) (chains []DBChain, err error)
	EnabledChains(...pg.QOpt) ([]DBChain, error)

	CreateNode(data Node, qopts ...pg.QOpt) (Node, error)
	DeleteNode(id int32, qopts ...pg.QOpt) error
	GetNodesByChainIDs(chainIDs []utils.Big, qopts ...pg.QOpt) (nodes []Node, err error)
	NodeNamed(string, ...pg.QOpt) (Node, error)
	Nodes(offset, limit int, qopts ...pg.QOpt) ([]Node, int, error)
	NodesForChain(chainID utils.Big, offset, limit int, qopts ...pg.QOpt) ([]Node, int, error)

	ChainConfigORM

	SetupNodes([]Node, []utils.Big) error
	EnsureChains([]utils.Big, ...pg.QOpt) error
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
type ChainCfg struct {
	BlockHistoryEstimatorBlockDelay                null.Int
	BlockHistoryEstimatorBlockHistorySize          null.Int
	BlockHistoryEstimatorEIP1559FeeCapBufferBlocks null.Int
	ChainType                                      null.String
	EthTxReaperThreshold                           *models.Duration
	EthTxResendAfterThreshold                      *models.Duration
	EvmEIP1559DynamicFees                          null.Bool
	EvmFinalityDepth                               null.Int
	EvmGasBumpPercent                              null.Int
	EvmGasBumpTxDepth                              null.Int
	EvmGasBumpWei                                  *assets.Wei
	EvmGasFeeCapDefault                            *assets.Wei
	EvmGasLimitDefault                             null.Int
	EvmGasLimitMax                                 null.Int
	EvmGasLimitMultiplier                          null.Float
	EvmGasLimitOCRJobType                          null.Int
	EvmGasLimitDRJobType                           null.Int
	EvmGasLimitVRFJobType                          null.Int
	EvmGasLimitFMJobType                           null.Int
	EvmGasLimitKeeperJobType                       null.Int
	EvmGasPriceDefault                             *assets.Wei
	EvmGasTipCapDefault                            *assets.Wei
	EvmGasTipCapMinimum                            *assets.Wei
	EvmHeadTrackerHistoryDepth                     null.Int
	EvmHeadTrackerMaxBufferSize                    null.Int
	EvmHeadTrackerSamplingInterval                 *models.Duration
	EvmLogBackfillBatchSize                        null.Int
	EvmLogPollInterval                             *models.Duration
	EvmLogKeepBlocksDepth                          null.Int
	EvmMaxGasPriceWei                              *assets.Wei
	EvmNonceAutoSync                               null.Bool
	EvmUseForwarders                               null.Bool
	EvmRPCDefaultBatchSize                         null.Int
	FlagsContractAddress                           null.String
	GasEstimatorMode                               null.String
	KeySpecific                                    map[string]ChainCfg
	LinkContractAddress                            null.String
	OperatorFactoryAddress                         null.String
	MinIncomingConfirmations                       null.Int
	MinimumContractPayment                         *assets.Link
	OCRObservationTimeout                          *models.Duration
	NodeNoNewHeadsThreshold                        *models.Duration
}

func (c *ChainCfg) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, c)
}

func (c *ChainCfg) Value() (driver.Value, error) {
	return json.Marshal(c)
}

type DBChain = chains.DBChain[utils.Big, *ChainCfg]

// TODO: https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
type Node struct {
	ID         int32
	Name       string
	EVMChainID utils.Big
	WSURL      null.String `db:"ws_url"`
	HTTPURL    null.String `db:"http_url"`
	SendOnly   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	// State doesn't exist in the DB, it's used to hold an in-memory state for
	// rendering
	State string `db:"-"`
}

// Receipt represents an ethereum receipt.
//
// Copied from go-ethereum: https://github.com/ethereum/go-ethereum/blob/ce9a289fa48e0d2593c4aaa7e207c8a5dd3eaa8a/core/types/receipt.go#L50
//
// We use our own version because Geth's version specifies various
// gencodec:"required" fields which cause unhelpful errors when unmarshalling
// from an empty JSON object which can happen in the batch fetcher.
type Receipt struct {
	PostState         []byte          `json:"root"`
	Status            uint64          `json:"status"`
	CumulativeGasUsed uint64          `json:"cumulativeGasUsed"`
	Bloom             gethTypes.Bloom `json:"logsBloom"`
	Logs              []*Log          `json:"logs"`
	TxHash            common.Hash     `json:"transactionHash"`
	ContractAddress   common.Address  `json:"contractAddress"`
	GasUsed           uint64          `json:"gasUsed"`
	BlockHash         common.Hash     `json:"blockHash,omitempty"`
	BlockNumber       *big.Int        `json:"blockNumber,omitempty"`
	TransactionIndex  uint            `json:"transactionIndex"`
}

// FromGethReceipt converts a gethTypes.Receipt to a Receipt
func FromGethReceipt(gr *gethTypes.Receipt) *Receipt {
	if gr == nil {
		return nil
	}
	logs := make([]*Log, len(gr.Logs))
	for i, glog := range gr.Logs {
		logs[i] = FromGethLog(glog)
	}
	return &Receipt{
		gr.PostState,
		gr.Status,
		gr.CumulativeGasUsed,
		gr.Bloom,
		logs,
		gr.TxHash,
		gr.ContractAddress,
		gr.GasUsed,
		gr.BlockHash,
		gr.BlockNumber,
		gr.TransactionIndex,
	}
}

// IsZero returns true if receipt is the zero receipt
// Batch calls to the RPC will return a pointer to an empty Receipt struct
// Easiest way to check if the receipt was missing is to see if the hash is 0x0
// Real receipts will always have the TxHash set
func (r Receipt) IsZero() bool {
	return r.TxHash == utils.EmptyHash
}

// IsUnmined returns true if the receipt is for a TX that has not been mined yet.
// Supposedly according to the spec this should never happen, but Parity does
// it anyway.
func (r Receipt) IsUnmined() bool {
	return r.BlockHash == utils.EmptyHash
}

// MarshalJSON marshals Receipt as JSON.
// Copied from: https://github.com/ethereum/go-ethereum/blob/ce9a289fa48e0d2593c4aaa7e207c8a5dd3eaa8a/core/types/gen_receipt_json.go
func (r Receipt) MarshalJSON() ([]byte, error) {
	type Receipt struct {
		PostState         hexutil.Bytes   `json:"root"`
		Status            hexutil.Uint64  `json:"status"`
		CumulativeGasUsed hexutil.Uint64  `json:"cumulativeGasUsed"`
		Bloom             gethTypes.Bloom `json:"logsBloom"`
		Logs              []*Log          `json:"logs"`
		TxHash            common.Hash     `json:"transactionHash"`
		ContractAddress   common.Address  `json:"contractAddress"`
		GasUsed           hexutil.Uint64  `json:"gasUsed"`
		BlockHash         common.Hash     `json:"blockHash,omitempty"`
		BlockNumber       *hexutil.Big    `json:"blockNumber,omitempty"`
		TransactionIndex  hexutil.Uint    `json:"transactionIndex"`
	}
	var enc Receipt
	enc.PostState = r.PostState
	enc.Status = hexutil.Uint64(r.Status)
	enc.CumulativeGasUsed = hexutil.Uint64(r.CumulativeGasUsed)
	enc.Bloom = r.Bloom
	enc.Logs = r.Logs
	enc.TxHash = r.TxHash
	enc.ContractAddress = r.ContractAddress
	enc.GasUsed = hexutil.Uint64(r.GasUsed)
	enc.BlockHash = r.BlockHash
	enc.BlockNumber = (*hexutil.Big)(r.BlockNumber)
	enc.TransactionIndex = hexutil.Uint(r.TransactionIndex)
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (r *Receipt) UnmarshalJSON(input []byte) error {
	type Receipt struct {
		PostState         *hexutil.Bytes   `json:"root"`
		Status            *hexutil.Uint64  `json:"status"`
		CumulativeGasUsed *hexutil.Uint64  `json:"cumulativeGasUsed"`
		Bloom             *gethTypes.Bloom `json:"logsBloom"`
		Logs              []*Log           `json:"logs"`
		TxHash            *common.Hash     `json:"transactionHash"`
		ContractAddress   *common.Address  `json:"contractAddress"`
		GasUsed           *hexutil.Uint64  `json:"gasUsed"`
		BlockHash         *common.Hash     `json:"blockHash,omitempty"`
		BlockNumber       *hexutil.Big     `json:"blockNumber,omitempty"`
		TransactionIndex  *hexutil.Uint    `json:"transactionIndex"`
	}
	var dec Receipt
	if err := json.Unmarshal(input, &dec); err != nil {
		return errors.Wrap(err, "could not unmarshal receipt")
	}
	if dec.PostState != nil {
		r.PostState = *dec.PostState
	}
	if dec.Status != nil {
		r.Status = uint64(*dec.Status)
	}
	if dec.CumulativeGasUsed != nil {
		r.CumulativeGasUsed = uint64(*dec.CumulativeGasUsed)
	}
	if dec.Bloom != nil {
		r.Bloom = *dec.Bloom
	}
	r.Logs = dec.Logs
	if dec.TxHash != nil {
		r.TxHash = *dec.TxHash
	}
	if dec.ContractAddress != nil {
		r.ContractAddress = *dec.ContractAddress
	}
	if dec.GasUsed != nil {
		r.GasUsed = uint64(*dec.GasUsed)
	}
	if dec.BlockHash != nil {
		r.BlockHash = *dec.BlockHash
	}
	if dec.BlockNumber != nil {
		r.BlockNumber = (*big.Int)(dec.BlockNumber)
	}
	if dec.TransactionIndex != nil {
		r.TransactionIndex = uint(*dec.TransactionIndex)
	}
	return nil
}

func (r *Receipt) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, r)
}

func (r *Receipt) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// Log represents a contract log event.
//
// Copied from go-ethereum: https://github.com/ethereum/go-ethereum/blob/ce9a289fa48e0d2593c4aaa7e207c8a5dd3eaa8a/core/types/log.go
//
// We use our own version because Geth's version specifies various
// gencodec:"required" fields which cause unhelpful errors when unmarshalling
// from an empty JSON object which can happen in the batch fetcher.
type Log struct {
	Address     common.Address `json:"address"`
	Topics      []common.Hash  `json:"topics"`
	Data        []byte         `json:"data"`
	BlockNumber uint64         `json:"blockNumber"`
	TxHash      common.Hash    `json:"transactionHash"`
	TxIndex     uint           `json:"transactionIndex"`
	BlockHash   common.Hash    `json:"blockHash"`
	Index       uint           `json:"logIndex"`
	Removed     bool           `json:"removed"`
}

// FromGethLog converts a gethTypes.Log to a Log
func FromGethLog(gl *gethTypes.Log) *Log {
	if gl == nil {
		return nil
	}
	return &Log{
		gl.Address,
		gl.Topics,
		gl.Data,
		gl.BlockNumber,
		gl.TxHash,
		gl.TxIndex,
		gl.BlockHash,
		gl.Index,
		gl.Removed,
	}
}

// MarshalJSON marshals as JSON.
func (l Log) MarshalJSON() ([]byte, error) {
	type Log struct {
		Address     common.Address `json:"address"`
		Topics      []common.Hash  `json:"topics"`
		Data        hexutil.Bytes  `json:"data"`
		BlockNumber hexutil.Uint64 `json:"blockNumber"`
		TxHash      common.Hash    `json:"transactionHash"`
		TxIndex     hexutil.Uint   `json:"transactionIndex"`
		BlockHash   common.Hash    `json:"blockHash"`
		Index       hexutil.Uint   `json:"logIndex"`
		Removed     bool           `json:"removed"`
	}
	var enc Log
	enc.Address = l.Address
	enc.Topics = l.Topics
	enc.Data = l.Data
	enc.BlockNumber = hexutil.Uint64(l.BlockNumber)
	enc.TxHash = l.TxHash
	enc.TxIndex = hexutil.Uint(l.TxIndex)
	enc.BlockHash = l.BlockHash
	enc.Index = hexutil.Uint(l.Index)
	enc.Removed = l.Removed
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (l *Log) UnmarshalJSON(input []byte) error {
	type Log struct {
		Address     *common.Address `json:"address"`
		Topics      []common.Hash   `json:"topics"`
		Data        *hexutil.Bytes  `json:"data"`
		BlockNumber *hexutil.Uint64 `json:"blockNumber"`
		TxHash      *common.Hash    `json:"transactionHash"`
		TxIndex     *hexutil.Uint   `json:"transactionIndex"`
		BlockHash   *common.Hash    `json:"blockHash"`
		Index       *hexutil.Uint   `json:"logIndex"`
		Removed     *bool           `json:"removed"`
	}
	var dec Log
	if err := json.Unmarshal(input, &dec); err != nil {
		return errors.Wrap(err, "could not unmarshal log")
	}
	if dec.Address != nil {
		l.Address = *dec.Address
	}
	l.Topics = dec.Topics
	if dec.Data != nil {
		l.Data = *dec.Data
	}
	if dec.BlockNumber != nil {
		l.BlockNumber = uint64(*dec.BlockNumber)
	}
	if dec.TxHash != nil {
		l.TxHash = *dec.TxHash
	}
	if dec.TxIndex != nil {
		l.TxIndex = uint(*dec.TxIndex)
	}
	if dec.BlockHash != nil {
		l.BlockHash = *dec.BlockHash
	}
	if dec.Index != nil {
		l.Index = uint(*dec.Index)
	}
	if dec.Removed != nil {
		l.Removed = *dec.Removed
	}
	return nil
}

type AddressArray []common.Address

func (a *AddressArray) Scan(src interface{}) error {
	baArray := pgtype.ByteaArray{}
	err := baArray.Scan(src)
	if err != nil {
		return errors.Wrap(err, "Expected BYTEA[] column for AddressArray")
	}
	if baArray.Status != pgtype.Present || len(baArray.Dimensions) > 1 {
		return errors.Errorf("Expected AddressArray to be 1-dimensional. Dimensions = %v", baArray.Dimensions)
	}

	for i, ba := range baArray.Elements {
		addr := common.Address{}
		if ba.Status != pgtype.Present {
			return errors.Errorf("Expected all addresses in AddressArray to be non-NULL.  Got AddressArray[%d] = NULL", i)
		}
		err = addr.Scan(ba.Bytes)
		if err != nil {
			return err
		}
		*a = append(*a, addr)
	}

	return nil
}

type HashArray []common.Hash

func (h *HashArray) Scan(src interface{}) error {
	baArray := pgtype.ByteaArray{}
	err := baArray.Scan(src)
	if err != nil {
		return errors.Wrap(err, "Expected BYTEA[] column for HashArray")
	}
	if baArray.Status != pgtype.Present || len(baArray.Dimensions) > 1 {
		return errors.Errorf("Expected HashArray to be 1-dimensional. Dimensions = %v", baArray.Dimensions)
	}

	for i, ba := range baArray.Elements {
		hash := common.Hash{}
		if ba.Status != pgtype.Present {
			return errors.Errorf("Expected all addresses in HashArray to be non-NULL.  Got HashArray[%d] = NULL", i)
		}
		err = hash.Scan(ba.Bytes)
		if err != nil {
			return err
		}
		*h = append(*h, hash)
	}
	return err
}
