package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Tx struct {
	ID       uint64 `storm:"id,increment,index"`
	From     string
	To       string
	Data     string
	Nonce    uint64
	Value    *big.Int
	GasLimit *big.Int
	TxAttempt
}

func (self *Tx) EthTx(gasPrice *big.Int) *types.Transaction {
	return types.NewTransaction(
		self.Nonce,
		common.HexToAddress(self.To),
		self.Value,
		self.GasLimit,
		gasPrice,
		common.FromHex(self.Data),
	)
}

type TxAttempt struct {
	Hash      string `storm:"id,index,unique"`
	TxID      uint64 `storm:"index"`
	GasPrice  *big.Int
	Confirmed bool
	Hex       string
	SentAt    uint64
}
