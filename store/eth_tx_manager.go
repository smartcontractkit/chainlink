package store

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EthTxManager struct {
	KeyStore *KeyStore
	Eth      *Eth
	Config   Config
}

func (self *EthTxManager) CreateTx(to, data string) (*types.Transaction, error) {
	tx, err := self.NewSignedTx(to, data)
	if err != nil {
		return tx, err
	}
	hex, err := encodeTxToHex(tx)
	if err != nil {
		return tx, err
	}
	if _, err = self.Eth.SendRawTx(hex); err != nil {
		return tx, err
	}
	return tx, nil
}

func (self *EthTxManager) NewSignedTx(to, data string) (*types.Transaction, error) {
	account := self.KeyStore.GetAccount()
	nonce, err := self.Eth.GetNonce(account)
	if err != nil {
		return nil, err
	}
	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(to),
		big.NewInt(0),
		big.NewInt(500000),
		big.NewInt(20000000000),
		common.FromHex(data),
	)
	return self.KeyStore.SignTx(tx, self.Config.ChainID)
}

func (self *EthTxManager) TxConfirmed(txid string) (bool, error) {
	receipt, err := self.Eth.GetTxReceipt(txid)
	if err != nil {
		return false, err
	} else if common.EmptyHash(receipt.TxHash) {
		return false, nil
	}
	return true, nil
}

func encodeTxToHex(tx *types.Transaction) (string, error) {
	rlp := new(bytes.Buffer)
	if err := tx.EncodeRLP(rlp); err != nil {
		return "", err
	}
	return common.ToHex(rlp.Bytes()), nil
}
