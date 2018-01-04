package store

import (
	"fmt"
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
		return nil, err
	}
	blkNum, err := self.BlockNumber()
	if err != nil {
		return nil, err
	}

	gasPrice := self.Config.EthGasPriceDefault
	_, err = self.createAttempt(txr, gasPrice, blkNum)
	if err != nil {
		return txr, err
	}

	return txr, nil
}

func (self *Eth) EnsureTxConfirmed(hash string) (bool, error) {
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
	txr := models.EthTx{}
	if err := self.ORM.One("ID", attempts[0].EthTxID, &txr); err != nil {
		return false, err
	}

	for _, txat := range attempts {
		success, err := self.checkAttempt(&txr, txat, blkNum)
		if success {
			return success, err
		}
	}
	return false, nil
}

func (self *Eth) createAttempt(
	txr *models.EthTx,
	gasPrice *big.Int,
	blkNum uint64,
) (*models.EthTxAttempt, error) {
	signable := txr.Signable(gasPrice)
	signable, err := self.KeyStore.SignTx(signable, self.Config.ChainID)
	if err != nil {
		return nil, err
	}

	a, err := self.ORM.AddAttempt(txr, signable, blkNum)
	if err != nil {
		return nil, err
	}
	return a, self.sendTransaction(signable)
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

func (self *Eth) getAttempts(hash string) ([]*models.EthTxAttempt, error) {
	attempt := &models.EthTxAttempt{}
	if err := self.ORM.One("Hash", hash, attempt); err != nil {
		return []*models.EthTxAttempt{}, err
	}
	attempts, err := self.ORM.AttemptsFor(attempt.EthTxID)
	if err != nil {
		return []*models.EthTxAttempt{}, err
	}
	return attempts, nil
}

func (self *Eth) checkAttempt(
	txr *models.EthTx,
	txat *models.EthTxAttempt,
	blkNum uint64,
) (bool, error) {
	receipt, err := self.GetTxReceipt(txat.Hash)
	if err != nil {
		return false, err
	}

	if receipt.Unconfirmed() {
		return self.handleUnconfirmed(txr, txat, blkNum)
	}
	return self.handleConfirmed(txr, txat, receipt, blkNum)
}

func (self *Eth) handleConfirmed(
	txr *models.EthTx,
	txat *models.EthTxAttempt,
	rcpt *TxReceipt,
	blkNum uint64,
) (bool, error) {

	safeAt := rcpt.BlockNumber + self.Config.EthMinConfirmations
	if blkNum < safeAt {
		return false, nil
	}

	if err := self.ORM.ConfirmTx(txr, txat); err != nil {
		return false, err
	}
	return true, nil
}

func (self *Eth) handleUnconfirmed(
	txr *models.EthTx,
	txat *models.EthTxAttempt,
	blkNum uint64,
) (bool, error) {
	bumpable := txr.Hash == txat.Hash
	pastThreshold := blkNum >= txat.SentAt+self.Config.EthGasBumpThreshold
	if bumpable && pastThreshold {
		return false, self.bumpGas(txat, blkNum)
	}
	return false, nil
}

func (self *Eth) bumpGas(txat *models.EthTxAttempt, blkNum uint64) error {
	txr := &models.EthTx{}
	if err := self.ORM.One("ID", txat.EthTxID, txr); err != nil {
		return err
	}
	gasPrice := new(big.Int).Add(txat.GasPrice, self.Config.EthGasBumpWei)
	_, err := self.createAttempt(txr, gasPrice, blkNum)
	if err != nil {
		return err
	}
	return self.ORM.Save(txat)
}
