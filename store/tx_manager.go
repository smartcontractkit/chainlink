package store

import (
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"sync"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
	"github.com/smartcontractkit/chainlink/utils"
	"go.uber.org/multierr"
)

const defaultGasLimit uint64 = 500000
const nonceReloadLimit uint = 1

// TxManager contains fields for the Ethereum client, the KeyStore,
// the local Config for the application, and the database.
type TxManager struct {
	*EthClient
	keyStore      *KeyStore
	config        Config
	orm           *orm.ORM
	activeAccount *ActiveAccount
}

// CreateTx signs and sends a transaction to the Ethereum blockchain.
func (txm *TxManager) CreateTx(to common.Address, data []byte) (*models.Tx, error) {
	return txm.createTxWithNonceReload(to, data, 0)
}

func (txm *TxManager) createTxWithNonceReload(to common.Address, data []byte, nrc uint) (*models.Tx, error) {
	if txm.activeAccount == nil {
		return nil, errors.New("Must activate an account before creating a transaction")
	}

	blkNum, err := txm.GetBlockNumber()
	if err != nil {
		return nil, err
	}

	var tx *models.Tx
	err = txm.activeAccount.GetAndIncrementNonce(func(nonce uint64) error {
		tx, err = txm.orm.CreateTx(
			txm.activeAccount.Address,
			nonce,
			to,
			data,
			big.NewInt(0),
			defaultGasLimit,
		)
		if err != nil {
			return err
		}

		logger.Infow(fmt.Sprintf("Created ETH transaction, attempt #: %v", nrc), []interface{}{"from", txm.activeAccount.Address.String(), "to", to.String()}...)
		gasPrice := txm.config.EthGasPriceDefault
		var txa *models.TxAttempt
		txa, err = txm.createAttempt(tx, &gasPrice, blkNum)
		if err != nil {
			txm.orm.DeleteStruct(tx)
			txm.orm.DeleteStruct(txa)

			return fmt.Errorf("TxManager CreateTX %v", err)
		}

		return nil
	})

	if err != nil {
		nonceErr, _ := regexp.MatchString("nonce .*too low", err.Error())
		if nonceErr {
			if nrc >= nonceReloadLimit {
				err = fmt.Errorf(
					"Transaction reattempt limit reached for 'nonce is too low' error. Limit: %v, Reattempt: %v",
					nonceReloadLimit,
					nrc,
				)
				return tx, err
			}

			logger.Warnw("Transaction nonce is too low. Reloading the nonce from the network and reattempting the transaction.")
			err = txm.ReloadNonce()
			if err != nil {
				return tx, fmt.Errorf("TxManager CreateTX ReloadNonce %v", err)
			}

			return txm.createTxWithNonceReload(to, data, nrc+1)
		}
	}

	return tx, err
}

// MeetsMinConfirmations returns true if the given transaction hash has been
// confirmed on the blockchain.
func (txm *TxManager) MeetsMinConfirmations(hash common.Hash) (bool, error) {
	blkNum, err := txm.GetBlockNumber()
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
	if err := txm.orm.One("ID", attempts[0].TxID, &tx); err != nil {
		return false, err
	}

	var merr error
	for _, txat := range attempts {
		success, err := txm.checkAttempt(&tx, &txat, blkNum)
		merr = multierr.Combine(merr, err)
		if success {
			return success, merr
		}
	}
	return false, merr
}

// WithdrawLink withdraws the given amount of LINK from the contract to the configured withdrawal address
func (txm *TxManager) WithdrawLink(wr models.WithdrawalRequest) (common.Hash, error) {
	functionSelector := models.HexToFunctionSelector("f3fef3a3") // withdraw(address _recipient, uint256 _amount)

	amount := (*big.Int)(wr.Amount)
	data, err := utils.HexToBytes(
		functionSelector.String(),
		common.ToHex(common.LeftPadBytes(wr.Address.Bytes(), utils.EVMWordByteLen)),
		utils.EVMHexNumber(amount),
	)
	tx, err := txm.CreateTx(*txm.config.OracleContractAddress, data)
	if err != nil {
		return common.Hash{}, err
	}

	return tx.Hash, nil
}

func (txm *TxManager) createAttempt(
	tx *models.Tx,
	gasPrice *big.Int,
	blkNum uint64,
) (*models.TxAttempt, error) {
	etx := tx.EthTx(gasPrice)
	etx, err := txm.keyStore.SignTx(etx, txm.config.ChainID)
	if err != nil {
		return nil, err
	}

	a, err := txm.orm.AddAttempt(tx, etx, blkNum)
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
		return fmt.Errorf("TxManager sendTransaction: %v", err)
	}
	return nil
}

func (txm *TxManager) getAttempts(hash common.Hash) ([]models.TxAttempt, error) {
	attempt := &models.TxAttempt{}
	if err := txm.orm.One("Hash", hash, attempt); err != nil {
		return []models.TxAttempt{}, err
	}
	attempts, err := txm.orm.AttemptsFor(attempt.TxID)
	if err != nil {
		return []models.TxAttempt{}, err
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

	minConfs := big.NewInt(int64(txm.config.MinOutgoingConfirmations))
	rcptBlkNum := rcpt.BlockNumber.ToBig()
	safeAt := minConfs.Add(rcptBlkNum, minConfs)
	safeAt.Sub(safeAt, big.NewInt(1)) // 0 based indexing since rcpt is 1 conf
	if big.NewInt(int64(blkNum)).Cmp(safeAt) == -1 {
		return false, nil
	}

	if err := txm.orm.ConfirmTx(tx, txat); err != nil {
		return false, err
	}
	logger.Infow(fmt.Sprintf("Confirmed tx %v", txat.Hash.String()), "txat", txat, "receipt", rcpt)
	return true, nil
}

func (txm *TxManager) handleUnconfirmed(
	tx *models.Tx,
	txat *models.TxAttempt,
	blkNum uint64,
) (bool, error) {
	bumpable := tx.Hash == txat.Hash
	pastThreshold := blkNum >= txat.SentAt+txm.config.EthGasBumpThreshold
	if bumpable && pastThreshold {
		return false, txm.bumpGas(txat, blkNum)
	}
	return false, nil
}

func (txm *TxManager) bumpGas(txat *models.TxAttempt, blkNum uint64) error {
	tx := &models.Tx{}
	if err := txm.orm.One("ID", txat.TxID, tx); err != nil {
		return err
	}
	gasPrice := new(big.Int).Add(txat.GasPrice, &txm.config.EthGasBumpWei)
	txat, err := txm.createAttempt(tx, gasPrice, blkNum)
	logger.Infow(fmt.Sprintf("Bumping gas to %v for transaction %v", gasPrice, txat.Hash.String()), "txat", txat)
	return err
}

// GetActiveAccount returns a copy of the TxManager's active nonce managed
// account.
func (txm *TxManager) GetActiveAccount() *ActiveAccount {
	if txm.activeAccount == nil {
		return nil
	}
	return &ActiveAccount{
		Account: txm.activeAccount.Account,
		nonce:   txm.activeAccount.nonce,
	}
}

// ActivateAccount retrieves an account's nonce from the blockchain for client
// side management in ActiveAccount.
func (txm *TxManager) ActivateAccount(account accounts.Account) error {
	nonce, err := txm.GetNonce(account.Address)
	if err != nil {
		return err
	}

	txm.activeAccount = &ActiveAccount{Account: account, nonce: nonce}
	return nil
}

// ReloadNonce fetch and update the current nonce via eth_getTransactionCount
func (txm *TxManager) ReloadNonce() error {
	nonce, err := txm.GetNonce(txm.activeAccount.Address)
	if err != nil {
		return fmt.Errorf("TxManager ReloadNonce: %v", err)
	}
	txm.activeAccount.nonce = nonce
	return nil
}

// ActiveAccount holds the account information alongside a client managed nonce
// to coordinate outgoing transactions.
type ActiveAccount struct {
	accounts.Account
	nonce uint64
	mutex sync.Mutex
}

// GetNonce returns the client side managed nonce.
func (a *ActiveAccount) GetNonce() uint64 {
	return a.nonce
}

// GetAndIncrementNonce will Yield the current nonce to a callback function and increment it once the
// callback has finished executing
func (a *ActiveAccount) GetAndIncrementNonce(callback func(uint64) error) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	err := callback(a.nonce)
	if err == nil {
		a.nonce = a.nonce + 1
	}

	return err
}
