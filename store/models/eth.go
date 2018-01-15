package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
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
