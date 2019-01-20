package models

import (
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/utils"
)

// Tx contains fields necessary for an Ethereum transaction with
// an additional field for the TxAttempt.
type Tx struct {
	ID       uint64         `gorm:"primary_key;auto_increment"`
	From     common.Address `gorm:"index;not null"`
	To       common.Address `gorm:"not null"`
	Data     []byte
	Nonce    uint64 `gorm:"index"`
	Value    *Big   `gorm:"type:varchar(255)"`
	GasLimit uint64
	// TxAttempt fields manually included; can't embed another primary_key
	Hash      common.Hash
	GasPrice  *Big `gorm:"type:varchar(255)"`
	Confirmed bool
	Hex       string `gorm:"type:text"`
	SentAt    uint64
}

// EthTx creates a new Ethereum transaction with a given gasPrice in wei
// that is ready to be signed.
func (tx *Tx) EthTx(gasPriceWei *big.Int) *types.Transaction {
	return types.NewTransaction(
		tx.Nonce,
		tx.To,
		tx.Value.ToInt(),
		tx.GasLimit,
		gasPriceWei,
		tx.Data,
	)
}

// AssignTxAttempt assigns the values of the attempt to the top level Tx.
func (tx *Tx) AssignTxAttempt(txat *TxAttempt) {
	tx.Hash = txat.Hash
	tx.GasPrice = txat.GasPrice
	tx.Confirmed = txat.Confirmed
	tx.Hex = txat.Hex
	tx.SentAt = txat.SentAt
}

// TxAttempt is used for keeping track of transactions that
// have been written to the Ethereum blockchain. This makes
// it so that if the network is busy, a transaction can be
// resubmitted with a higher GasPrice.
type TxAttempt struct {
	Hash      common.Hash `gorm:"primary_key;not null"`
	TxID      uint64      `gorm:"index"`
	GasPrice  *Big        `gorm:"type:varchar(255)"`
	Confirmed bool
	Hex       string `gorm:"type:text"`
	SentAt    uint64
	CreatedAt time.Time `gorm:"index"`
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

// UnmarshalJSON parses the raw FunctionSelector and sets the FunctionSelector
// type to the given input.
func (f *FunctionSelector) UnmarshalJSON(input []byte) error {
	var s string
	err := json.Unmarshal(input, &s)
	if err != nil {
		return err
	}

	if utils.HasHexPrefix(s) {
		bytes := common.FromHex(s)
		if len(bytes) != FunctionSelectorLength {
			return errors.New("Function ID must be 4 bytes in length")
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

// BlockHeader represents a block header in the Ethereum blockchain.
// Deliberately does not have required fields because some fields aren't
// present depending on the Ethereum node.
// i.e. Parity does not always send mixHash
type BlockHeader struct {
	ParentHash  common.Hash      `json:"parentHash"`
	UncleHash   common.Hash      `json:"sha3Uncles"`
	Coinbase    common.Address   `json:"miner"`
	Root        common.Hash      `json:"stateRoot"`
	TxHash      common.Hash      `json:"transactionsRoot"`
	ReceiptHash common.Hash      `json:"receiptsRoot"`
	Bloom       types.Bloom      `json:"logsBloom"`
	Difficulty  hexutil.Big      `json:"difficulty"`
	Number      hexutil.Big      `json:"number"`
	GasLimit    hexutil.Uint64   `json:"gasLimit"`
	GasUsed     hexutil.Uint64   `json:"gasUsed"`
	Time        hexutil.Big      `json:"timestamp"`
	Extra       hexutil.Bytes    `json:"extraData"`
	Nonce       types.BlockNonce `json:"nonce"`
	GethHash    common.Hash      `json:"mixHash"`
	ParityHash  common.Hash      `json:"hash"`
}

var emptyHash = common.Hash{}

// Hash will return GethHash if it exists otherwise it returns the ParityHash
func (h BlockHeader) Hash() common.Hash {
	if h.GethHash != emptyHash {
		return h.GethHash
	}
	return h.ParityHash
}

// ToIndexableBlockNumber converts a given BlockHeader to an IndexableBlockNumber
func (h BlockHeader) ToIndexableBlockNumber() *IndexableBlockNumber {
	return NewIndexableBlockNumber(h.Number.ToInt(), h.Hash())
}

// IndexableBlockNumber represents a BlockNumber, BlockHash and the number of Digits in the BlockNumber
type IndexableBlockNumber struct {
	Number Big         `json:"number" gorm:"index;type:varchar(255);not null"`
	Digits int         `json:"digits" gorm:"index"`
	Hash   common.Hash `json:"hash"`
}

// NewIndexableBlockNumber creates an IndexableBlockNumber given a BlockNumber and BlockHash
func NewIndexableBlockNumber(bigint *big.Int, hash common.Hash) *IndexableBlockNumber {
	if bigint == nil {
		return nil
	}
	number := hexutil.Big(*bigint)
	return &IndexableBlockNumber{
		Number: Big(number),
		Digits: len(number.String()) - 2,
		Hash:   hash,
	}
}

// String returns a string representation of this number.
func (l *IndexableBlockNumber) String() string {
	return l.ToInt().String()
}

// ToInt Coerces the value into *big.Int. Also handles nil *IndexableBlockNumber values to
// nil *big.Int.
func (l *IndexableBlockNumber) ToInt() *big.Int {
	if l == nil {
		return nil
	}
	return l.Number.ToInt()
}

// GreaterThan compares BlockNumbers and returns true if the reciever BlockNumber is greater than
// the supplied BlockNumber
func (l *IndexableBlockNumber) GreaterThan(r *IndexableBlockNumber) bool {
	if l == nil {
		return false
	}
	if l != nil && r == nil {
		return true
	}
	return l.ToInt().Cmp(r.ToInt()) > 0
}

// NextInt returns the next BlockNumber as big.int, or nil if nil to represent latest.
func (l *IndexableBlockNumber) NextInt() *big.Int {
	if l == nil {
		return nil
	}
	return new(big.Int).Add(l.ToInt(), big.NewInt(1))
}

// EthSubscription should implement Err() <-chan error and Unsubscribe()
type EthSubscription interface {
	Err() <-chan error
	Unsubscribe()
}
