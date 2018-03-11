package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

// Tx contains fields necessary for an Ethereum transaction with
// an additional field for the TxAttempt.
type Tx struct {
	ID       uint64 `storm:"id,increment,index"`
	From     common.Address
	To       common.Address
	Data     []byte
	Nonce    uint64
	Value    *big.Int
	GasLimit uint64
	TxAttempt
}

// EthTx creates a new Ethereum transaction with a given gasPrice
// that is ready to be signed.
func (tx *Tx) EthTx(gasPrice *big.Int) *types.Transaction {
	return types.NewTransaction(
		tx.Nonce,
		tx.To,
		tx.Value,
		tx.GasLimit,
		gasPrice,
		tx.Data,
	)
}

// TxAttempt is used for keeping track of transactions that
// have been written to the Ethereum blockchain. This makes
// it so that if the network is busy, a transaction can be
// resubmitted with a higher GasPrice.
type TxAttempt struct {
	Hash      common.Hash `storm:"id,unique"`
	TxID      uint64      `storm:"index"`
	GasPrice  *big.Int
	Confirmed bool
	Hex       string
	SentAt    uint64
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

	bytes := common.FromHex(s)
	if len(bytes) != FunctionSelectorLength {
		return errors.New("Function ID must be 4 bytes in length")
	}

	f.SetBytes(bytes)
	return nil
}

// Represents a block header in the Ethereum blockchain.
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

func (h BlockHeader) Hash() common.Hash {
	if !common.EmptyHash(h.GethHash) {
		return h.GethHash
	}
	return h.ParityHash
}

func (h BlockHeader) IndexableBlockNumber() *IndexableBlockNumber {
	return NewIndexableBlockNumber(h.Number.ToInt(), h.Hash())
}

type IndexableBlockNumber struct {
	Number hexutil.Big `json:"number" storm:"id,unique"`
	Digits int         `json:"digits" storm:"index"`
	Hash   common.Hash `json:"hash"`
}

func NewIndexableBlockNumber(bigint *big.Int, hashes ...common.Hash) *IndexableBlockNumber {
	if bigint == nil {
		return nil
	}
	var hash common.Hash
	if len(hashes) > 0 {
		hash = hashes[0]
	}
	number := hexutil.Big(*bigint)
	return &IndexableBlockNumber{
		Number: number,
		Digits: len(number.String()) - 2,
		Hash:   hash,
	}
}

// Coerces the value into *big.Int. Also handles nil *IndexableBlockNumber values to
// nil *big.Int.
func (n *IndexableBlockNumber) ToInt() *big.Int {
	if n == nil {
		return nil
	}
	return n.Number.ToInt()
}

// Return a hex string representation of the block number, or empty string if nil.
func (n *IndexableBlockNumber) String() string {
	if n == nil {
		return ""
	}
	return n.Number.String()
}

func (n *IndexableBlockNumber) FriendlyString() string {
	return fmt.Sprintf("#%v (%v)", n.ToInt(), n.String())
}

func (l *IndexableBlockNumber) GreaterThan(r *IndexableBlockNumber) bool {
	if l == nil {
		return false
	}
	if l != nil && r == nil {
		return true
	}
	return l.ToInt().Cmp(r.ToInt()) > 0
}

func (l *IndexableBlockNumber) NextInt() *big.Int {
	if l == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Add(l.ToInt(), big.NewInt(1))
}

func (l *IndexableBlockNumber) NextNumber() *IndexableBlockNumber {
	if l != nil {
		return NewIndexableBlockNumber(l.NextInt(), l.Hash)
	}
	return NewIndexableBlockNumber(l.NextInt())
}

type EthSubscription interface {
	Err() <-chan error
	Unsubscribe()
}
