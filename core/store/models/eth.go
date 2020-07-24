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
	null "gopkg.in/guregu/null.v3"
)

type EthTxState string
type EthTxAttemptState string

const (
	EthTxUnstarted   = EthTxState("unstarted")
	EthTxInProgress  = EthTxState("in_progress")
	EthTxFatalError  = EthTxState("fatal_error")
	EthTxUnconfirmed = EthTxState("unconfirmed")
	EthTxConfirmed   = EthTxState("confirmed")

	EthTxAttemptInProgress = EthTxAttemptState("in_progress")
	EthTxAttemptBroadcast  = EthTxAttemptState("broadcast")
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
	BroadcastAt    *time.Time
	CreatedAt      time.Time
	State          EthTxState
	EthTxAttempts  []EthTxAttempt `gorm:"association_autoupdate:false;association_autocreate:false"`
}

func (e EthTx) GetError() error {
	if e.Error == nil {
		return nil
	}
	return errors.New(*e.Error)
}

// GetID allows EthTx to be used as jsonapi.MarshalIdentifier
func (e EthTx) GetID() string {
	return string(e.ID)
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
	EthReceipts             []EthReceipt `gorm:"foreignkey:TxHash;association_foreignkey:Hash;association_autoupdate:false;association_autocreate:false"`
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

// Tx contains fields necessary for an Ethereum transaction with
// an additional field for the TxAttempt.
type Tx struct {
	ID uint64 `gorm:"primary_key;auto_increment"`

	// SurrogateID is used to look up a transaction using a secondary ID, used to
	// associate jobs with transactions so that we don't double spend in certain
	// failure scenarios
	SurrogateID null.String `gorm:"index;unique"`

	Attempts []*TxAttempt `json:"-"`

	From     common.Address `gorm:"index;not null"`
	To       common.Address `gorm:"not null"`
	Data     []byte         `gorm:"not null"`
	Nonce    uint64         `gorm:"index;not null"`
	Value    *utils.Big     `gorm:"not null"`
	GasLimit uint64         `gorm:"not null"`

	// TxAttempt fields manually included; can't embed another primary_key
	Hash        common.Hash `gorm:"not null"`
	GasPrice    *utils.Big  `gorm:"not null"`
	Confirmed   bool        `gorm:"not null"`
	SentAt      uint64      `gorm:"not null"`
	SignedRawTx []byte      `gorm:"not null"`
	CreatedAt   time.Time   `json:"-"`
	UpdatedAt   time.Time   `json:"-"`
}

// String implements Stringer for Tx
func (tx *Tx) String() string {
	return fmt.Sprintf("Tx(ID: %d, From: %s, To: %s, Hash: %s, SentAt: %d)",
		tx.ID,
		tx.From.String(),
		tx.To.String(),
		tx.Hash.String(),
		tx.SentAt)
}

// EthTx creates a new Ethereum transaction with a given gasPrice in wei
// that is ready to be signed.
func (tx Tx) EthTx(gasPriceWei *big.Int) *types.Transaction {
	return types.NewTransaction(
		tx.Nonce,
		tx.To,
		tx.Value.ToInt(),
		tx.GasLimit,
		gasPriceWei,
		tx.Data,
	)
}

// TxAttempt is used for keeping track of transactions that
// have been written to the Ethereum blockchain. This makes
// it so that if the network is busy, a transaction can be
// resubmitted with a higher GasPrice.
type TxAttempt struct {
	ID uint64 `gorm:"primary_key;auto_increment"`

	TxID uint64 `gorm:"index;type:bigint REFERENCES txes(id) ON DELETE CASCADE"`
	Tx   *Tx    `json:"-" gorm:"PRELOAD:false;foreignkey:TxID"`

	CreatedAt time.Time `gorm:"index;not null"`

	Hash        common.Hash `gorm:"index;not null"`
	GasPrice    *utils.Big  `gorm:"type:varchar(78);not null"`
	Confirmed   bool        `gorm:"not null"`
	SentAt      uint64      `gorm:"not null"`
	SignedRawTx []byte      `gorm:"not null"`
	UpdatedAt   time.Time   `json:"-"`
}

// String implements Stringer for TxAttempt
func (txa *TxAttempt) String() string {
	return fmt.Sprintf("TxAttempt{ID: %d, TxID: %d, Hash: %s, SentAt: %d, Confirmed: %t}",
		txa.ID,
		txa.TxID,
		txa.Hash.String(),
		txa.SentAt,
		txa.Confirmed)
}

// GetID returns the ID of this structure for jsonapi serialization.
func (txa TxAttempt) GetID() string {
	return txa.Hash.Hex()
}

// GetName returns the pluralized "type" of this structure for jsonapi serialization.
func (txa TxAttempt) GetName() string {
	return "txattempts"
}

// SetID is used to set the ID of this structure when deserializing from jsonapi documents.
func (txa *TxAttempt) SetID(value string) error {
	txa.Hash = common.HexToHash(value)
	return nil
}

func HighestPricedTxAttemptPerTx(items []TxAttempt) []TxAttempt {
	highestPricedSet := map[uint64]TxAttempt{}
	for _, item := range items {
		if currentHighest, ok := highestPricedSet[item.TxID]; ok {
			if currentHighest.GasPrice.ToInt().Cmp(item.GasPrice.ToInt()) == -1 {
				highestPricedSet[item.TxID] = item
			}
		} else {
			highestPricedSet[item.TxID] = item
		}
	}
	highestPriced := make([]TxAttempt, len(highestPricedSet))
	i := 0
	for _, attempt := range highestPricedSet {
		highestPriced[i] = attempt
		i++
	}
	return highestPriced
}

// Head represents a BlockNumber, BlockHash.
type Head struct {
	ID         uint64
	Hash       common.Hash
	Number     int64
	ParentHash common.Hash
	Parent     *Head
	Timestamp  time.Time
	CreatedAt  time.Time
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

// NewHead returns a Head instance.
func NewHead(number *big.Int, blockHash common.Hash, parentHash common.Hash, timestamp uint64) Head {
	return Head{
		Number:     number.Int64(),
		Hash:       blockHash,
		ParentHash: parentHash,
		Timestamp:  time.Unix(int64(timestamp), 0),
	}
}

// NewHeadFromBlockHeader returns a new Head from geth's types.Header
func NewHeadFromBlockHeader(h types.Header) Head {
	return NewHead(h.Number, h.Hash(), h.ParentHash, h.Time)
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

// WeiPerEth is amount of Wei currency units in one Eth.
var WeiPerEth = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

// This data can contain anything and is submitted by user on-chain, so we must
// be extra careful how we interact with it
type UntrustedBytes []byte

type Log = types.Log

var emptyHash = common.Hash{}

// TxReceipt holds the block number and the transaction hash of a signed
// transaction that has been written to the blockchain.
type TxReceipt struct {
	BlockNumber *utils.Big   `json:"blockNumber"`
	BlockHash   *common.Hash `json:"blockHash"`
	Hash        common.Hash  `json:"transactionHash"`
	Logs        []Log        `json:"logs"`
}

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

// WithoutPrefix returns the FunctionSelector as a string without the '0x' prefix.
func (f FunctionSelector) WithoutPrefix() string { return f.String()[2:] }

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
func (f FunctionSelector) Scan(value interface{}) error {
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

// SafeByteSlice returns an error on out of bounds access to a byte array, where a
// normal slice would panic instead
func (ary UntrustedBytes) SafeByteSlice(start int, end int) ([]byte, error) {
	if end > len(ary) || start > end || start < 0 || end < 0 {
		var empty []byte
		return empty, errors.New("out of bounds slice access")
	}
	return ary[start:end], nil
}
