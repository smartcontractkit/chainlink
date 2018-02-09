package models

import (
	"encoding/json"
	"errors"
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
	GasLimit *big.Int
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
	Hash      common.Hash `storm:"id,index,unique"`
	TxID      uint64      `storm:"index"`
	GasPrice  *big.Int
	Confirmed bool
	Hex       string
	SentAt    uint64
}

// FunctionID is the first four bytes of the call data for a
// function call and specifies the function to be called.
type FunctionID [FunctionIDLength]byte

// FunctionIDLength should always be a length of 4 as a byte.
const FunctionIDLength = 4

// BytesToFunctionID converts the given bytes to a FunctionID.
func BytesToFunctionID(b []byte) FunctionID {
	var f FunctionID
	f.SetBytes(b)
	return f
}

// HexToFunctionID converts the given string to a FunctionID.
func HexToFunctionID(s string) FunctionID { return BytesToFunctionID(common.FromHex(s)) }

// String returns the FunctionID as a string type.
func (f FunctionID) String() string { return hexutil.Encode(f[:]) }

// WithoutPrefix returns the FunctionID as a string without the '0x' prefix.
func (f FunctionID) WithoutPrefix() string { return f.String()[2:] }

// SetBytes sets the FunctionID to that of the given bytes (will trim).
func (f *FunctionID) SetBytes(b []byte) { copy(f[:], b[:FunctionIDLength]) }

// UnmarshalJSON parses the raw FunctionID and sets the FunctionID
// type to the given input.
func (f *FunctionID) UnmarshalJSON(input []byte) error {
	var s string
	err := json.Unmarshal(input, &s)
	if err != nil {
		return err
	}

	bytes := common.FromHex(s)
	if len(bytes) != FunctionIDLength {
		return errors.New("Function ID must be 4 bytes in length")
	}

	f.SetBytes(bytes)
	return nil
}
