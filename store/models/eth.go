package models

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

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

type TxAttempt struct {
	Hash      common.Hash `storm:"id,index,unique"`
	TxID      uint64      `storm:"index"`
	GasPrice  *big.Int
	Confirmed bool
	Hex       string
	SentAt    uint64
}

type FunctionID [FunctionIDLength]byte

const FunctionIDLength = 4

func BytesToFunctionID(b []byte) FunctionID {
	var f FunctionID
	f.SetBytes(b)
	return f
}

func HexToFunctionID(s string) FunctionID { return BytesToFunctionID(common.FromHex(s)) }

func (f FunctionID) String() string        { return hexutil.Encode(f[:]) }
func (f FunctionID) WithoutPrefix() string { return f.String()[2:] }
func (f *FunctionID) SetBytes(b []byte)    { copy(f[:], b[:FunctionIDLength]) }

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
