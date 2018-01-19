package store

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
)

type TxManager struct {
	*EthClient
	KeyStore *KeyStore
	Config   Config
	ORM      *models.ORM
}

func (txm *TxManager) CreateTx(to common.Address, data []byte) (*models.Tx, error) {
	account := txm.KeyStore.GetAccount()
	nonce, err := txm.GetNonce(account)
	if err != nil {
		return nil, err
	}
	tx, err := txm.ORM.CreateTx(
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
	blkNum, err := txm.BlockNumber()
	if err != nil {
		return nil, err
	}

	gasPrice := &txm.Config.EthGasPriceDefault
	_, err = txm.createAttempt(tx, gasPrice, blkNum)
	if err != nil {
		return tx, err
	}

	return tx, nil
}

func (txm *TxManager) EnsureTxConfirmed(hash common.Hash) (bool, error) {
	blkNum, err := txm.BlockNumber()
	if err != nil {
		return false, err
	}
	attempts, err := txm.getAttempts(hash)
	if err != nil {
		return false, err
	}
	if len(attempts) == 0 {
		return false, fmt.Errorf("Can only ensure transactions with attempts")
	}
	tx := models.Tx{}
	if err := txm.ORM.One("ID", attempts[0].TxID, &tx); err != nil {
		return false, err
	}

	for _, txat := range attempts {
		success, err := txm.checkAttempt(&tx, txat, blkNum)
		if success {
			return success, err
		}
	}
	return false, nil
}

func (txm *TxManager) createAttempt(
	tx *models.Tx,
	gasPrice *big.Int,
	blkNum uint64,
) (*models.TxAttempt, error) {
	etx := tx.EthTx(gasPrice)
	etx, err := txm.KeyStore.SignTx(etx, txm.Config.ChainID)
	if err != nil {
		return nil, err
	}

	a, err := txm.ORM.AddAttempt(tx, etx, blkNum)
	if err != nil {
		return nil, err
	}
	return a, txm.sendTransaction(etx)
}

func (txm *TxManager) sendTransaction(tx *types.Transaction) error {
	hex, err := utils.EncodeTxToHex(tx)
	if err != nil {
		return err
	}
	if _, err = txm.SendRawTx(hex); err != nil {
		return err
	}
	return nil
}

func (txm *TxManager) getAttempts(hash common.Hash) ([]*models.TxAttempt, error) {
	attempt := &models.TxAttempt{}
	if err := txm.ORM.One("Hash", hash, attempt); err != nil {
		return []*models.TxAttempt{}, err
	}
	attempts, err := txm.ORM.AttemptsFor(attempt.TxID)
	if err != nil {
		return []*models.TxAttempt{}, err
	}
	return attempts, nil
}

func (txm *TxManager) checkAttempt(
	tx *models.Tx,
	txat *models.TxAttempt,
	blkNum uint64,
) (bool, error) {
	receipt, err := txm.GetTxReceipt(txat.Hash)
	if err != nil {
		return false, err
	}

	if receipt.Unconfirmed() {
		return txm.handleUnconfirmed(tx, txat, blkNum)
	}
	return txm.handleConfirmed(tx, txat, receipt, blkNum)
}

func (txm *TxManager) handleConfirmed(
	tx *models.Tx,
	txat *models.TxAttempt,
	rcpt *TxReceipt,
	blkNum uint64,
) (bool, error) {

	safeAt := rcpt.BlockNumber + txm.Config.EthMinConfirmations
	if blkNum < safeAt {
		return false, nil
	}

	if err := txm.ORM.ConfirmTx(tx, txat); err != nil {
		return false, err
	}
	return true, nil
}

func (txm *TxManager) handleUnconfirmed(
	tx *models.Tx,
	txat *models.TxAttempt,
	blkNum uint64,
) (bool, error) {
	bumpable := tx.Hash == txat.Hash
	pastThreshold := blkNum >= txat.SentAt+txm.Config.EthGasBumpThreshold
	if bumpable && pastThreshold {
		return false, txm.bumpGas(txat, blkNum)
	}
	return false, nil
}

func (txm *TxManager) bumpGas(txat *models.TxAttempt, blkNum uint64) error {
	tx := &models.Tx{}
	if err := txm.ORM.One("ID", txat.TxID, tx); err != nil {
		return err
	}
	gasPrice := new(big.Int).Add(txat.GasPrice, &txm.Config.EthGasBumpWei)
	_, err := txm.createAttempt(tx, gasPrice, blkNum)
	if err != nil {
		return err
	}
	return txm.ORM.Save(txat)
}
