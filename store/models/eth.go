package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EthTx struct {
	ID       uint64 `storm:"id,increment,index"`
	From     string
	To       string
	Data     string
	Nonce    uint64
	Value    *big.Int
	GasLimit *big.Int
	EthTxAttempt
}

func (self *EthTx) Signable(gasPrice *big.Int) *types.Transaction {
	return types.NewTransaction(
		self.Nonce,
		common.HexToAddress(self.To),
		self.Value,
		self.GasLimit,
		gasPrice,
		common.FromHex(self.Data),
	)
}

type EthTxAttempt struct {
	TxID      string `storm:"id,index,unique"`
	EthTxID   uint64 `storm:"index"`
	GasPrice  *big.Int
	Confirmed bool
	Hex       string
	SentAt    uint64
	Bumped    bool
}
