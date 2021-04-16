package models

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type EthTxState string
type EthTxAttemptState string

const (
	EthTxUnstarted               = EthTxState("unstarted")
	EthTxInProgress              = EthTxState("in_progress")
	EthTxFatalError              = EthTxState("fatal_error")
	EthTxUnconfirmed             = EthTxState("unconfirmed")
	EthTxConfirmed               = EthTxState("confirmed")
	EthTxConfirmedMissingReceipt = EthTxState("confirmed_missing_receipt")

	EthTxAttemptInProgress      = EthTxAttemptState("in_progress")
	EthTxAttemptInsufficientEth = EthTxAttemptState("insufficient_eth")
	EthTxAttemptBroadcast       = EthTxAttemptState("broadcast")
)

type EthTaskRunTx struct {
	TaskRunID uuid.UUID
	EthTxID   int64
	EthTx     EthTx
}

type EthTx struct {
	ID             int64
	Nonce          *int64
	FromAddress    common.Address
	ToAddress      common.Address
	EncodedPayload []byte
	Value          assets.Eth
	GasLimit       uint64
	Error          *string
	// BroadcastAt is updated every time an attempt for this eth_tx is re-sent
	// In almost all cases it will be within a second or so of the actual send time.
	BroadcastAt   *time.Time
	CreatedAt     time.Time
	State         EthTxState
	EthTxAttempts []EthTxAttempt `gorm:"->"`
}

func (e EthTx) GetError() error {
	if e.Error == nil {
		return nil
	}
	return errors.New(*e.Error)
}

// GetID allows EthTx to be used as jsonapi.MarshalIdentifier
func (e EthTx) GetID() string {
	return fmt.Sprintf("%d", e.ID)
}

type EthTxAttempt struct {
	ID                      int64
	EthTxID                 int64
	EthTx                   EthTx
	GasPrice                utils.Big
	SignedRawTx             []byte
	Hash                    common.Hash
	CreatedAt               time.Time
	BroadcastBeforeBlockNum *int64
	State                   EthTxAttemptState
	EthReceipts             []EthReceipt `gorm:"foreignKey:TxHash;references:Hash;association_foreignkey:Hash;->"`
}

type EthReceipt struct {
	ID               int64
	TxHash           common.Hash
	BlockHash        common.Hash
	BlockNumber      int64
	TransactionIndex uint
	Receipt          []byte
	CreatedAt        time.Time
}

// GetSignedTx decodes the SignedRawTx into a types.Transaction struct
func (a EthTxAttempt) GetSignedTx() (*types.Transaction, error) {
	s := rlp.NewStream(bytes.NewReader(a.SignedRawTx), 0)
	signedTx := new(types.Transaction)
	if err := signedTx.DecodeRLP(s); err != nil {
		logger.Error("could not decode RLP")
		return nil, err
	}
	return signedTx, nil
}

// Head represents a BlockNumber, BlockHash.
type Head struct {
	ID         uint64
	Hash       common.Hash
	Number     int64
	ParentHash common.Hash
	Parent     *Head `gorm:"-"`
	Timestamp  time.Time
	CreatedAt  time.Time
}

// NewHead returns a Head instance.
func NewHead(number *big.Int, blockHash common.Hash, parentHash common.Hash, timestamp uint64) Head {
	return Head{
		Number:     number.Int64(),
		Hash:       blockHash,
		ParentHash: parentHash,
		Timestamp:  time.Unix(int64(timestamp), 0),
	}
}

// EarliestInChain recurses through parents until it finds the earliest one
func (h Head) EarliestInChain() Head {
	for {
		if h.Parent != nil {
			h = *h.Parent
		} else {
			break
		}
	}
	return h
}

// IsInChain returns true if the given hash matches the hash of a head in the chain
func (h Head) IsInChain(blockHash common.Hash) bool {
	for {
		if h.Hash == blockHash {
			return true
		}
		if h.Parent != nil {
			h = *h.Parent
		} else {
			break
		}
	}
	return false
}

// HashAtHeight returns the hash of the block at the given heigh, if it is in the chain.
// If not in chain, returns the zero hash
func (h Head) HashAtHeight(blockNum int64) common.Hash {
	for {
		if h.Number == blockNum {
			return h.Hash
		}
		if h.Parent != nil {
			h = *h.Parent
		} else {
			break
		}
	}
	return common.Hash{}
}

// ChainLength returns the length of the chain followed by recursively looking up parents
func (h Head) ChainLength() uint32 {
	l := uint32(1)

	for {
		if h.Parent != nil {
			l++
			h = *h.Parent
		} else {
			break
		}
	}
	return l
}

// String returns a string representation of this number.
func (h *Head) String() string {
	return h.ToInt().String()
}

// ToInt return the height as a *big.Int. Also handles nil by returning nil.
func (h *Head) ToInt() *big.Int {
	if h == nil {
		return nil
	}
	return big.NewInt(h.Number)
}

// GreaterThan compares BlockNumbers and returns true if the receiver BlockNumber is greater than
// the supplied BlockNumber
func (h *Head) GreaterThan(r *Head) bool {
	if h == nil {
		return false
	}
	if h != nil && r == nil {
		return true
	}
	return h.Number > r.Number
}

// NextInt returns the next BlockNumber as big.int, or nil if nil to represent latest.
func (h *Head) NextInt() *big.Int {
	if h == nil {
		return nil
	}
	return new(big.Int).Add(h.ToInt(), big.NewInt(1))
}

func (h *Head) UnmarshalJSON(bs []byte) error {
	type head struct {
		Hash       common.Hash    `json:"hash"`
		Number     *hexutil.Big   `json:"number"`
		ParentHash common.Hash    `json:"parentHash"`
		Timestamp  hexutil.Uint64 `json:"timestamp"`
	}

	var jsonHead head
	err := json.Unmarshal(bs, &jsonHead)
	if err != nil {
		return err
	}

	if jsonHead.Number == nil {
		*h = Head{}
		return nil
	}

	h.Hash = jsonHead.Hash
	h.Number = (*big.Int)(jsonHead.Number).Int64()
	h.ParentHash = jsonHead.ParentHash
	h.Timestamp = time.Unix(int64(jsonHead.Timestamp), 0).UTC()
	return nil
}

func (h *Head) MarshalJSON() ([]byte, error) {
	type head struct {
		Hash       *common.Hash    `json:"hash,omitempty"`
		Number     *hexutil.Big    `json:"number,omitempty"`
		ParentHash *common.Hash    `json:"parentHash,omitempty"`
		Timestamp  *hexutil.Uint64 `json:"timestamp,omitempty"`
	}

	var jsonHead head
	if h.Hash != (common.Hash{}) {
		jsonHead.Hash = &h.Hash
	}
	jsonHead.Number = (*hexutil.Big)(big.NewInt(int64(h.Number)))
	if h.ParentHash != (common.Hash{}) {
		jsonHead.ParentHash = &h.ParentHash
	}
	if h.Timestamp != (time.Time{}) {
		t := hexutil.Uint64(h.Timestamp.UTC().Unix())
		jsonHead.Timestamp = &t
	}
	return json.Marshal(jsonHead)
}

// WeiPerEth is amount of Wei currency units in one Eth.
var WeiPerEth = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

type Log = types.Log

var emptyHash = common.Hash{}

// Unconfirmed returns true if the transaction is not confirmed.
func ReceiptIsUnconfirmed(txr *types.Receipt) bool {
	return txr == nil || txr.TxHash == emptyHash || txr.BlockNumber == nil
}

// ChainlinkFulfilledTopic is the signature for the event emitted after calling
// ChainlinkClient.validateChainlinkCallback(requestId). See
// ../../evm-contracts/src/v0.6/ChainlinkClient.sol
var ChainlinkFulfilledTopic = utils.MustHash("ChainlinkFulfilled(bytes32)")

// ReceiptIndicatesRunLogFulfillment returns true if this tx receipt is the result of a
// fulfilled run log.
func ReceiptIndicatesRunLogFulfillment(txr types.Receipt) bool {
	for _, log := range txr.Logs {
		if log.Topics[0] == ChainlinkFulfilledTopic {
			return true
		}
	}
	return false
}

// FunctionSelector is the first four bytes of the call data for a
// function call and specifies the function to be called.
type FunctionSelector [FunctionSelectorLength]byte

// FunctionSelectorLength should always be a length of 4 as a byte.
const FunctionSelectorLength = 4

// BytesToFunctionSelector converts the given bytes to a FunctionSelector.
func BytesToFunctionSelector(b []byte) FunctionSelector {
	var f FunctionSelector
	f.SetBytes(b)
	return f
}

// HexToFunctionSelector converts the given string to a FunctionSelector.
func HexToFunctionSelector(s string) FunctionSelector {
	return BytesToFunctionSelector(common.FromHex(s))
}

// String returns the FunctionSelector as a string type.
func (f FunctionSelector) String() string { return hexutil.Encode(f[:]) }

// Bytes returns the FunctionSelector as a byte slice
func (f FunctionSelector) Bytes() []byte { return f[:] }

// SetBytes sets the FunctionSelector to that of the given bytes (will trim).
func (f *FunctionSelector) SetBytes(b []byte) { copy(f[:], b[:FunctionSelectorLength]) }

var hexRegexp *regexp.Regexp = regexp.MustCompile("^[0-9a-fA-F]*$")

func unmarshalFromString(s string, f *FunctionSelector) error {
	if utils.HasHexPrefix(s) {
		if !hexRegexp.Match([]byte(s)[2:]) {
			return fmt.Errorf("function selector %s must be 0x-hex encoded", s)
		}
		bytes := common.FromHex(s)
		if len(bytes) != FunctionSelectorLength {
			return errors.New("function ID must be 4 bytes in length")
		}
		f.SetBytes(bytes)
	} else {
		bytes, err := utils.Keccak256([]byte(s))
		if err != nil {
			return err
		}
		f.SetBytes(bytes[0:4])
	}
	return nil
}

// UnmarshalJSON parses the raw FunctionSelector and sets the FunctionSelector
// type to the given input.
func (f *FunctionSelector) UnmarshalJSON(input []byte) error {
	var s string
	err := json.Unmarshal(input, &s)
	if err != nil {
		return err
	}
	return unmarshalFromString(s, f)
}

// MarshalJSON returns the JSON encoding of f
func (f FunctionSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

// Value returns this instance serialized for database storage
func (f FunctionSelector) Value() (driver.Value, error) {
	return f.Bytes(), nil
}

// Scan returns the selector from its serialization in the database
func (f *FunctionSelector) Scan(value interface{}) error {
	temp, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unable to convent %v of type %T to FunctionSelector", value, value)
	}
	if len(temp) != FunctionSelectorLength {
		return fmt.Errorf("function selector %v should have length %d, but has length %d",
			temp, FunctionSelectorLength, len(temp))
	}
	copy(f[:], temp)
	return nil
}

// This data can contain anything and is submitted by user on-chain, so we must
// be extra careful how we interact with it
type UntrustedBytes []byte

// SafeByteSlice returns an error on out of bounds access to a byte array, where a
// normal slice would panic instead
func (ary UntrustedBytes) SafeByteSlice(start int, end int) ([]byte, error) {
	if end > len(ary) || start > end || start < 0 || end < 0 {
		var empty []byte
		return empty, errors.New("out of bounds slice access")
	}
	return ary[start:end], nil
}

type blockInternal struct {
	Number       string
	Hash         common.Hash
	ParentHash   common.Hash
	Transactions []types.Transaction
}

// Int64ToHex converts an int64 into go-ethereum's hex representation
func Int64ToHex(n int64) string {
	return hexutil.EncodeBig(big.NewInt(n))
}

// Block represents an ethereum block
type Block struct {
	Number       int64
	Hash         common.Hash
	ParentHash   common.Hash
	Transactions []types.Transaction
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
		return err
	}
	n, err := hexutil.DecodeBig(bi.Number)
	if err != nil {
		return err
	}
	*b = Block{
		n.Int64(),
		bi.Hash,
		bi.ParentHash,
		bi.Transactions,
	}
	return nil
}
