package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink-go/utils"
)

type EthTx struct {
	ID       uint64 `storm:"id,increment"`
	From     string
	To       string
	Data     string
	Nonce    uint64
	Value    *big.Int
	GasLimit *big.Int
	GasPrice *big.Int
	Attempts []*EthTxAttempt `storm:"inline"`
}

func (self *EthTx) NewAttempt(tx *types.Transaction) (*EthTxAttempt, error) {
	hex, err := utils.EncodeTxToHex(tx)
	if err != nil {
		return nil, err
	}
	attempt := &EthTxAttempt{
		TxID:      tx.Hash().String(),
		GasPrice:  tx.GasPrice(),
		Confirmed: false,
		Hex:       hex,
	}

	self.Attempts = append(self.Attempts, attempt)
	return attempt, nil
}

func (self *EthTx) TxID() string {
	return self.Attempts[len(self.Attempts)-1].TxID
}

func (self *EthTx) Signable() *types.Transaction {
	return types.NewTransaction(
		self.Nonce,
		common.HexToAddress(self.To),
		self.Value,
		self.GasLimit,
		self.GasPrice,
		common.FromHex(self.Data),
	)
}

type EthTxAttempt struct {
	TxID      string
	GasPrice  *big.Int
	Confirmed bool
	Hex       string
}
