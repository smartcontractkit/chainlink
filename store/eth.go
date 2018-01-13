package store

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
)

type Eth struct {
	*EthClient
	KeyStore     *KeyStore
	Config       Config
	ORM          *models.ORM
}

func (self *Eth) CreateTx(to common.Address, data []byte) (*models.Tx, error) {
	account := self.KeyStore.GetAccount()
	nonce, err := self.GetNonce(account)
	if err != nil {
		return nil, err
	}
	tx, err := self.ORM.CreateTx(
		account.Address,
		nonce,
		to,
		data,
		big.NewInt(0),
		big.NewInt(500000),
	)
	if err != nil {
		return nil, err
	}
	blkNum, err := self.BlockNumber()
	if err != nil {
		return nil, err
	}

	gasPrice := self.Config.EthGasPriceDefault
	_, err = self.createAttempt(tx, gasPrice, blkNum)
	if err != nil {
		return tx, err
	}

	return tx, nil
}

func (self *Eth) EnsureTxConfirmed(hash common.Hash) (bool, error) {
	blkNum, err := self.BlockNumber()
	if err != nil {
		return false, err
	}
	attempts, err := self.getAttempts(hash)
	if err != nil {
		return false, err
	}
	if len(attempts) == 0 {
		return false, fmt.Errorf("Can only ensure transactions with attempts")
	}
	tx := models.Tx{}
	if err := self.ORM.One("ID", attempts[0].TxID, &tx); err != nil {
		return false, err
	}

	for _, txat := range attempts {
		success, err := self.checkAttempt(&tx, txat, blkNum)
		if success {
			return success, err
		}
	}
	return false, nil
}

func (self *Eth) createAttempt(
	tx *models.Tx,
	gasPrice *big.Int,
	blkNum uint64,
) (*models.TxAttempt, error) {
	etx := tx.EthTx(gasPrice)
	etx, err := self.KeyStore.SignTx(etx, self.Config.ChainID)
	if err != nil {
		return nil, err
	}

	a, err := self.ORM.AddAttempt(tx, etx, blkNum)
	if err != nil {
		return nil, err
	}
	return a, self.sendTransaction(etx)
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

func (self *Eth) getAttempts(hash common.Hash) ([]*models.TxAttempt, error) {
	attempt := &models.TxAttempt{}
	if err := self.ORM.One("Hash", hash, attempt); err != nil {
		return []*models.TxAttempt{}, err
	}
	attempts, err := self.ORM.AttemptsFor(attempt.TxID)
	if err != nil {
		return []*models.TxAttempt{}, err
	}
	return attempts, nil
}

func (self *Eth) checkAttempt(
	tx *models.Tx,
	txat *models.TxAttempt,
	blkNum uint64,
) (bool, error) {
	receipt, err := self.GetTxReceipt(txat.Hash)
	if err != nil {
		return false, err
	}

	if receipt.Unconfirmed() {
		return self.handleUnconfirmed(tx, txat, blkNum)
	}
	return self.handleConfirmed(tx, txat, receipt, blkNum)
}

func (self *Eth) handleConfirmed(
	tx *models.Tx,
	txat *models.TxAttempt,
	rcpt *TxReceipt,
	blkNum uint64,
) (bool, error) {

	safeAt := rcpt.BlockNumber + self.Config.EthMinConfirmations
	if blkNum < safeAt {
		return false, nil
	}

	if err := self.ORM.ConfirmTx(tx, txat); err != nil {
		return false, err
	}
	return true, nil
}

func (self *Eth) handleUnconfirmed(
	tx *models.Tx,
	txat *models.TxAttempt,
	blkNum uint64,
) (bool, error) {
	bumpable := tx.Hash == txat.Hash
	pastThreshold := blkNum >= txat.SentAt+self.Config.EthGasBumpThreshold
	if bumpable && pastThreshold {
		return false, self.bumpGas(txat, blkNum)
	}
	return false, nil
}

func (self *Eth) bumpGas(txat *models.TxAttempt, blkNum uint64) error {
	tx := &models.Tx{}
	if err := self.ORM.One("ID", txat.TxID, tx); err != nil {
		return err
	}
	gasPrice := new(big.Int).Add(txat.GasPrice, self.Config.EthGasBumpWei)
	_, err := self.createAttempt(tx, gasPrice, blkNum)
	if err != nil {
		return err
	}
	return self.ORM.Save(txat)
}
