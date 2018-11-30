package store

import (
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"sync"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
	"github.com/smartcontractkit/chainlink/utils"
	"go.uber.org/multierr"
)

const defaultGasLimit uint64 = 500000
const nonceReloadLimit uint = 1

// TxManager represents an interface for interacting with the blockchain
type TxManager interface {
	Start(accounts []accounts.Account) error
	CreateTx(to common.Address, data []byte) (*models.Tx, error)
	MeetsMinConfirmations(hash common.Hash) (bool, error)
	WithdrawLink(wr models.WithdrawalRequest) (common.Hash, error)
	GetLinkBalance(address common.Address) (*assets.Link, error)
	GetActiveAccount() *ManagedAccount

	GetEthBalance(address common.Address) (*assets.Eth, error)
	SubscribeToNewHeads(channel chan<- models.BlockHeader) (models.EthSubscription, error)
	GetBlockByNumber(hex string) (models.BlockHeader, error)
	SubscribeToLogs(channel chan<- Log, q ethereum.FilterQuery) (models.EthSubscription, error)
	GetLogs(q ethereum.FilterQuery) ([]Log, error)
}

// EthTxManager contains fields for the Ethereum client, the KeyStore,
// the local Config for the application, and the database.
type EthTxManager struct {
	*EthClient
	keyStore            *KeyStore
	config              Config
	orm                 *orm.ORM
	availableAccounts   []*ManagedAccount
	availableAccountIdx int
}

// CreateTx signs and sends a transaction to the Ethereum blockchain.
func (txm *EthTxManager) CreateTx(to common.Address, data []byte) (*models.Tx, error) {
	ma := txm.GetActiveAccount()
	if ma == nil {
		return nil, errors.New("Must activate an account before creating a transaction")
	}

	return txm.createTxWithNonceReload(ma, to, data, 0)
}

func (txm *EthTxManager) createTxWithNonceReload(ma *ManagedAccount, to common.Address, data []byte, nrc uint) (*models.Tx, error) {
	blkNum, err := txm.GetBlockNumber()
	if err != nil {
		return nil, err
	}

	var tx *models.Tx
	err = ma.GetAndIncrementNonce(func(nonce uint64) error {
		tx, err = txm.orm.CreateTx(
			ma.Address,
			nonce,
			to,
			data,
			big.NewInt(0),
			defaultGasLimit,
		)
		if err != nil {
			return err
		}

		logger.Infow(fmt.Sprintf("Created ETH transaction, attempt #: %v", nrc), "from", ma.Address.String(), "to", to.String())
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
			err = txm.ReloadNonce(ma)
			if err != nil {
				return tx, fmt.Errorf("TxManager CreateTX ReloadNonce %v", err)
			}

			return txm.createTxWithNonceReload(ma, to, data, nrc+1)
		}
	}

	return tx, err
}

// GetLinkBalance returns the balance of LINK at the given address
func (txm *EthTxManager) GetLinkBalance(address common.Address) (*assets.Link, error) {
	contractAddress := common.HexToAddress(txm.config.LinkContractAddress)
	balance, err := txm.GetERC20Balance(address, contractAddress)
	if err != nil {
		return assets.NewLink(0), err
	}
	return (*assets.Link)(balance), nil
}

// MeetsMinConfirmations returns true if the given transaction hash has been
// confirmed on the blockchain.
func (txm *EthTxManager) MeetsMinConfirmations(hash common.Hash) (bool, error) {
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
		merr = multierr.Append(merr, err)
		if success {
			return success, merr
		}
	}
	return false, merr
}

// WithdrawLink withdraws the given amount of LINK from the contract to the configured withdrawal address
func (txm *EthTxManager) WithdrawLink(wr models.WithdrawalRequest) (common.Hash, error) {
	functionSelector := models.HexToFunctionSelector("f3fef3a3") // withdraw(address _recipient, uint256 _amount)

	amount := (*big.Int)(wr.Amount)
	data, err := utils.ConcatBytes(
		functionSelector.Bytes(),
		common.LeftPadBytes(wr.Address.Bytes(), utils.EVMWordByteLen),
		common.LeftPadBytes(amount.Bytes(), utils.EVMWordByteLen),
	)

	if txm.config.OracleContractAddress == nil {
		return common.Hash{}, errors.New("OracleContractAddress not set can not withdraw")
	}
	tx, err := txm.CreateTx(*txm.config.OracleContractAddress, data)
	if err != nil {
		return common.Hash{}, err
	}

	return tx.Hash, nil
}

func (txm *EthTxManager) createAttempt(
	tx *models.Tx,
	gasPrice *big.Int,
	blkNum uint64,
) (*models.TxAttempt, error) {
	ma := txm.getAccount(tx.From)
	if ma == nil {
		return nil, fmt.Errorf("Unable to locate %v as an available account in EthTxManager. Has TxManager been started or has the address been removed?", tx.From.Hex())
	}
	etx := tx.EthTx(gasPrice)
	etx, err := txm.keyStore.SignTx(ma.Account, etx, txm.config.ChainID)
	if err != nil {
		return nil, err
	}

	a, err := txm.orm.AddAttempt(tx, etx, blkNum)
	if err != nil {
		return nil, err
	}
	return a, txm.sendTransaction(etx)
}

func (txm *EthTxManager) sendTransaction(tx *types.Transaction) error {
	hex, err := utils.EncodeTxToHex(tx)
	if err != nil {
		return err
	}
	if _, err = txm.SendRawTx(hex); err != nil {
		return fmt.Errorf("TxManager sendTransaction: %v", err)
	}
	return nil
}

func (txm *EthTxManager) getAttempts(hash common.Hash) ([]models.TxAttempt, error) {
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

func (txm *EthTxManager) checkAttempt(
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

// GetETHAndLINKBalances attempts to retrieve the ethereum node's perception of
// the latest ETH and LINK balances for the active account on the txm, or an
// error on failure.
func (txm *EthTxManager) GetETHAndLINKBalances() (*big.Int, *assets.Link, error) {
	ma := txm.GetActiveAccount()
	if ma == nil {
		return big.NewInt(0), assets.NewLink(0), fmt.Errorf(
			"Could not find activeAccount for which to report new balances")
	}
	address := ma.Account.Address
	linkBalance, linkErr := txm.GetLinkBalance(address)
	ethBalance, ethErr := txm.EthClient.GetWeiBalance(address)
	merr := multierr.Append(linkErr, ethErr)
	return ethBalance, linkBalance, merr
}

// ethAndLINKBalancesReport constructs the log message reporting on the current
// ETH and LINK balances, or the error which occurred while retrieving them from
// the ethereum node (merr)
func ethAndLINKBalancesReport(ethBalance *big.Int, linkBalance *assets.Link, merr error) string {
	if merr == nil {
		return fmt.Sprintf(
			"New ETH balance: %v. New LINK balance: %v",
			ethBalance, linkBalance)
	}
	return fmt.Sprintf(
		"Failed to retrieve LINK or ETH balance following confirmation! %v",
		merr)
}

// handleConfirmed checks whether a tx is confirmed, and records and reports it
// as such if so. Its bool return value is true if the tx is confirmed and it
// was successfully recorded as confirmed.
func (txm *EthTxManager) handleConfirmed(
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

	logMessage := fmt.Sprintf("Confirmed tx %v", txat.Hash.String())
	ethBalance, linkBalance, merr := txm.GetETHAndLINKBalances()
	balanceMessage := ethAndLINKBalancesReport(ethBalance, linkBalance, merr)
	logger.Infow(logMessage+" "+balanceMessage,
		"txat", txat, "receipt", rcpt)

	return true, nil
}

func (txm *EthTxManager) handleUnconfirmed(
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

func (txm *EthTxManager) bumpGas(txat *models.TxAttempt, blkNum uint64) error {
	tx := &models.Tx{}
	if err := txm.orm.One("ID", txat.TxID, tx); err != nil {
		return err
	}
	gasPrice := new(big.Int).Add(txat.GasPrice, &txm.config.EthGasBumpWei)
	txat, err := txm.createAttempt(tx, gasPrice, blkNum)
	if err != nil {
		return err
	}
	logger.Infow(fmt.Sprintf("Bumping gas to %v for transaction %v", gasPrice, txat.Hash.String()), "txat", txat)
	return nil
}

// GetActiveAccount uses round robing to select a managed account
// from the list of available accounts as defined in Start(...)
func (txm *EthTxManager) GetActiveAccount() *ManagedAccount {
	if len(txm.availableAccounts) == 0 {
		return nil
	}

	current := txm.availableAccountIdx
	txm.availableAccountIdx++
	if txm.availableAccountIdx >= len(txm.availableAccounts) {
		txm.availableAccountIdx = 0
	}
	return txm.availableAccounts[current]
}

func (txm *EthTxManager) getAccount(from common.Address) *ManagedAccount {
	for _, a := range txm.availableAccounts {
		if a.Address == from {
			return a
		}
	}

	return nil
}

// Start activates accounts for outgoing transactions and client side
// nonce management.
func (txm *EthTxManager) Start(accounts []accounts.Account) error {
	var merr error
	for _, a := range accounts {
		merr = multierr.Append(merr, txm.activateAccount(a))
	}

	return merr
}

// ActivateAccount retrieves an account's nonce from the blockchain for client
// side management in ManagedAccount.
func (txm *EthTxManager) activateAccount(account accounts.Account) error {
	nonce, err := txm.GetNonce(account.Address)
	if err != nil {
		return err
	}

	txm.availableAccounts = append(txm.availableAccounts, NewManagedAccount(account, nonce))
	return nil
}

// ReloadNonce fetch and update the current nonce via eth_getTransactionCount
func (txm *EthTxManager) ReloadNonce(ma *ManagedAccount) error {
	nonce, err := txm.GetNonce(ma.Address)
	if err != nil {
		return fmt.Errorf("TxManager ReloadNonce: %v", err)
	}
	ma.nonce = nonce
	return nil
}

// ManagedAccount holds the account information alongside a client managed nonce
// to coordinate outgoing transactions.
type ManagedAccount struct {
	accounts.Account
	nonce uint64
	mutex *sync.Mutex
}

// NewManagedAccount creates a managed account that handles nonce increments
// locally.
func NewManagedAccount(a accounts.Account, nonce uint64) *ManagedAccount {
	return &ManagedAccount{Account: a, nonce: nonce, mutex: &sync.Mutex{}}
}

// GetNonce returns the client side managed nonce.
func (a *ManagedAccount) GetNonce() uint64 {
	return a.nonce
}

// GetAndIncrementNonce will Yield the current nonce to a callback function and increment it once the
// callback has finished executing
func (a *ManagedAccount) GetAndIncrementNonce(callback func(uint64) error) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	err := callback(a.nonce)
	if err == nil {
		a.nonce = a.nonce + 1
	}

	return err
}
