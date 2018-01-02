package store

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink-go/store/models"
	"github.com/smartcontractkit/chainlink-go/utils"
)

type Eth struct {
	*EthClient
	KeyStore *KeyStore
	Config   Config
	ORM      *models.ORM
}

func (self *Eth) CreateTx(to, data string) (*models.EthTx, error) {
	account := self.KeyStore.GetAccount()
	nonce, err := self.GetNonce(account)
	if err != nil {
		return nil, err
	}
	txr, err := self.ORM.CreateEthTx(
		account.Address.String(),
		nonce,
		to,
		data,
		big.NewInt(0),
		big.NewInt(500000),
		big.NewInt(20000000000),
	)
	if err != nil {
		return txr, err
	}

	if err = self.createAttempt(txr); err != nil {
		return txr, err
	}

	return txr, nil
}

func (self *Eth) createAttempt(txr *models.EthTx) error {
	tx := txr.Signable()
	tx, err := self.KeyStore.SignTx(tx, self.Config.ChainID)
	if err != nil {
		return err
	}
	if _, err = txr.NewAttempt(tx); err != nil {
		return err
	}
	if err = self.sendTransaction(tx); err != nil {
		return err
	}
	return self.ORM.Save(txr)
}

func (self *Eth) sendTransaction(tx *types.Transaction) error {
	hex, err := utils.EncodeTxToHex(tx)
	if err != nil {
		return err
	}
	if _, err = self.SendRawTx(hex); err != nil {
		return err
	}
	return nil
}

func (self *Eth) TxConfirmed(txid string) (bool, error) {
	receipt, err := self.GetTxReceipt(txid)
	if err != nil {
		return false, err
	} else if receipt.Unconfirmed() {
		return false, nil
	}

	min := receipt.BlockNumber + self.Config.EthConfMin
	current, err := self.BlockNumber()
	if err != nil {
		return false, err
	}
	return (min <= current), nil
}
