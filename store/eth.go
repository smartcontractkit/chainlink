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
	)
	if err != nil {
		return txr, err
	}

	gasPrice := self.Config.EthGasPriceDefault
	if err = self.createAttempt(txr, gasPrice); err != nil {
		return txr, err
	}

	return txr, nil
}

func (self *Eth) createAttempt(txr *models.EthTx, gasPrice *big.Int) error {
	tx := txr.Signable(gasPrice)
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
	return self.ORM.SaveTx(txr)
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

func (self *Eth) EnsureTxConfirmed(txid string) (bool, error) {
	receipt, err := self.GetTxReceipt(txid)
	if err != nil {
		return false, err
	}
	blkNum, err := self.BlockNumber()
	if err != nil {
		return false, err
	}
	txat := &models.EthTxAttempt{}
	if err := self.ORM.One("TxID", txid, txat); err != nil {
		return false, err
	}
	if receipt.Unconfirmed() {
		return self.handleUnconfirmed(receipt, txat, blkNum)
	}
	return self.handleConfirmed(receipt, txat, blkNum)
}

func (self *Eth) handleConfirmed(
	rcpt *TxReceipt,
	txat *models.EthTxAttempt,
	blkNum uint64,
) (bool, error) {
	safeAt := rcpt.BlockNumber + self.Config.EthMinConfirmations
	return blkNum >= safeAt, nil
}

func (self *Eth) handleUnconfirmed(
	rcpt *TxReceipt,
	txat *models.EthTxAttempt,
	blkNum uint64,
) (bool, error) {
	if blkNum >= txat.SentAt+self.Config.EthGasBumpThreshold {
		return false, self.bumpGas(txat)
	}
	return false, nil
}

func (self *Eth) bumpGas(txat *models.EthTxAttempt) error {
	txr := &models.EthTx{}
	if err := self.ORM.One("ID", txat.EthTxID, txr); err != nil {
		return err
	}
	gasPrice := new(big.Int).Add(txr.GasPrice(), self.Config.EthGasBumpWei)
	return self.createAttempt(txr, gasPrice)
}
