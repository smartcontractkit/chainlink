package store

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/tevino/abool"
	"go.uber.org/multierr"
)

// DefaultGasLimit sets the default gas limit for outgoing transactions.
// if updating DefaultGasLimit, be sure it matches with the
// DefaultGasLimit specified in evm/test/Oracle_test.js
const DefaultGasLimit uint64 = 500000
const nonceReloadLimit uint = 1

// ErrPendingConnection is the error returned if TxManager is not connected.
var ErrPendingConnection = errors.New("Cannot talk to chain, pending connection")

// TxManager represents an interface for interacting with the blockchain
type TxManager interface {
	HeadTrackable
	Connected() bool
	Register(accounts []accounts.Account)
	CreateTx(to common.Address, data []byte) (*models.Tx, error)
	CreateTxWithGas(to common.Address, data []byte, gasPriceWei *big.Int, gasLimit uint64) (*models.Tx, error)
	CreateTxWithEth(from, to common.Address, value *assets.Eth) (*models.Tx, error)
	BumpGasUntilSafe(hash common.Hash) (*models.TxReceipt, error)
	ContractLINKBalance(wr models.WithdrawalRequest) (assets.Link, error)
	WithdrawLINK(wr models.WithdrawalRequest) (common.Hash, error)
	GetLINKBalance(address common.Address) (*assets.Link, error)
	NextActiveAccount() *ManagedAccount

	GetEthBalance(address common.Address) (*assets.Eth, error)
	SubscribeToNewHeads(channel chan<- models.BlockHeader) (models.EthSubscription, error)
	GetBlockByNumber(hex string) (models.BlockHeader, error)
	SubscribeToLogs(channel chan<- models.Log, q ethereum.FilterQuery) (models.EthSubscription, error)
	GetLogs(q ethereum.FilterQuery) ([]models.Log, error)
}

//go:generate mockgen -package=mocks -destination=../internal/mocks/tx_manager_mocks.go github.com/smartcontractkit/chainlink/core/store TxManager

// EthTxManager contains fields for the Ethereum client, the KeyStore,
// the local Config for the application, and the database.
type EthTxManager struct {
	*EthClient
	keyStore            *KeyStore
	config              Config
	orm                 *orm.ORM
	registeredAccounts  []accounts.Account
	availableAccounts   []*ManagedAccount
	availableAccountIdx int
	accountsMutex       *sync.Mutex
	connected           *abool.AtomicBool
}

// NewEthTxManager constructs an EthTxManager using the passed variables and
// initializing internal variables.
func NewEthTxManager(ethClient *EthClient, config Config, keyStore *KeyStore, orm *orm.ORM) *EthTxManager {
	return &EthTxManager{
		EthClient:     ethClient,
		config:        config,
		keyStore:      keyStore,
		orm:           orm,
		accountsMutex: &sync.Mutex{},
		connected:     abool.New(),
	}
}

// Register activates accounts for outgoing transactions and client side
// nonce management.
func (txm *EthTxManager) Register(accts []accounts.Account) {
	txm.accountsMutex.Lock()
	defer txm.accountsMutex.Unlock()

	cp := make([]accounts.Account, len(accts))
	copy(cp, accts)
	txm.registeredAccounts = cp
}

// Connected returns a bool indicating whether or not it is connected.
func (txm *EthTxManager) Connected() bool {
	return txm.connected.IsSet()
}

// Connect iterates over the available accounts to retrieve their nonce
// for client side management.
func (txm *EthTxManager) Connect(bn *models.Head) error {
	txm.accountsMutex.Lock()
	defer txm.accountsMutex.Unlock()

	txm.availableAccounts = []*ManagedAccount{}
	var merr error
	for _, a := range txm.registeredAccounts {
		ma, err := txm.activateAccount(a)
		merr = multierr.Append(merr, err)
		if err == nil {
			txm.availableAccounts = append(txm.availableAccounts, ma)
		}
	}

	txm.connected.Set()
	return merr
}

// Disconnect marks this instance as disconnected.
func (txm *EthTxManager) Disconnect() {
	txm.connected.UnSet()
}

// OnNewHead does nothing; exists to comply with interface.
func (txm *EthTxManager) OnNewHead(*models.Head) {}

// CreateTx signs and sends a transaction to the Ethereum blockchain.
func (txm *EthTxManager) CreateTx(to common.Address, data []byte) (*models.Tx, error) {
	return txm.CreateTxWithGas(to, data, txm.config.EthGasPriceDefault(), DefaultGasLimit)
}

// CreateTxWithGas signs and sends a transaction to the Ethereum blockchain.
func (txm *EthTxManager) CreateTxWithGas(to common.Address, data []byte, gasPriceWei *big.Int, gasLimit uint64) (*models.Tx, error) {
	ma, err := txm.nextAccount()
	if err != nil {
		return nil, err
	}

	gasPriceWei, gasLimit = normalize(gasPriceWei, gasLimit, txm.config)
	return txm.createTxWithNonceReload(ma, to, data, gasPriceWei, gasLimit, 0)
}

// CreateTxWithEth signs and sends a transaction with some ETH to transfer.
func (txm *EthTxManager) CreateTxWithEth(from, to common.Address, value *assets.Eth) (*models.Tx, error) {
	ma := txm.getAccount(from)
	if ma == nil {
		return nil, errors.New("account does not exist")
	}

	return txm.createEthTxWithNonceReload(ma, to, []byte{}, txm.config.EthGasPriceDefault(), DefaultGasLimit, value, 0)
}

func (txm *EthTxManager) nextAccount() (*ManagedAccount, error) {
	if !txm.Connected() {
		return nil, errors.Wrap(ErrPendingConnection, "EthTxManager#nextAccount")
	}

	ma := txm.NextActiveAccount()
	if ma == nil {
		return nil, errors.New("Must connect and activate an account before creating a transaction")
	}

	return ma, nil
}

func normalize(gasPriceWei *big.Int, gasLimit uint64, config Config) (*big.Int, uint64) {
	if !config.Dev() {
		return config.EthGasPriceDefault(), DefaultGasLimit
	}

	if gasPriceWei == nil {
		gasPriceWei = config.EthGasPriceDefault()
	}

	if gasLimit == 0 {
		gasLimit = DefaultGasLimit
	}

	return gasPriceWei, gasLimit
}

func (txm *EthTxManager) createEthTxWithNonceReload(
	ma *ManagedAccount,
	to common.Address,
	data []byte,
	gasPriceWei *big.Int,
	gasLimit uint64,
	value *assets.Eth,
	nrc uint) (*models.Tx, error) {

	blkNum, err := txm.getBlockNumber()
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
			value.ToInt(),
			gasLimit,
		)
		if err != nil {
			return err
		}

		logger.Infow(fmt.Sprintf("Created ETH transaction, attempt #: %v", nrc), "from", ma.Address.String(), "to", to.String())
		_, err = txm.createAttempt(tx, gasPriceWei, blkNum)
		if err != nil {
			merr := multierr.Append(err, txm.orm.DeleteTransaction(tx))
			return fmt.Errorf("TxManager CreateTX %v", merr)
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
			err = ma.ReloadNonce(txm)
			if err != nil {
				return tx, fmt.Errorf("TxManager CreateTX ReloadNonce %v", err)
			}

			return txm.createTxWithNonceReload(ma, to, data, gasPriceWei, gasLimit, nrc+1)
		}
	}

	return tx, err
}

func (txm *EthTxManager) createTxWithNonceReload(
	ma *ManagedAccount,
	to common.Address,
	data []byte,
	gasPriceWei *big.Int,
	gasLimit uint64,
	nrc uint) (*models.Tx, error) {
	return txm.createEthTxWithNonceReload(
		ma,
		to,
		data,
		gasPriceWei,
		gasLimit,
		assets.NewEth(0),
		nrc)
}

// GetLINKBalance returns the balance of LINK at the given address
func (txm *EthTxManager) GetLINKBalance(address common.Address) (*assets.Link, error) {
	contractAddress := common.HexToAddress(txm.config.LinkContractAddress())
	balance, err := txm.GetERC20Balance(address, contractAddress)
	if err != nil {
		return assets.NewLink(0), err
	}
	return (*assets.Link)(balance), nil
}

// BumpGasUntilSafe returns true if the given transaction hash has been
// confirmed on the blockchain.
func (txm *EthTxManager) BumpGasUntilSafe(hash common.Hash) (*models.TxReceipt, error) {
	blkNum, err := txm.getBlockNumber()
	if err != nil {
		return nil, errors.Wrap(err, "BumpGasUntilSafe getBlockNumber")
	}
	tx, attempts, err := txm.getTxAndAttempts(hash)
	if err != nil {
		return nil, errors.Wrap(err, "BumpGasUntilSafe getTxAndAttempts")
	}

	var merr error
	for _, txat := range attempts {
		receipt, state, err := txm.checkAttempt(tx, &txat, blkNum)
		if state == safe || state == confirmed {
			return receipt, err // success, so all other attempt errors can be ignored.
		}
		merr = multierr.Append(merr, err)
	}
	return nil, merr
}

func (txm *EthTxManager) getTxAndAttempts(hash common.Hash) (*models.Tx, []models.TxAttempt, error) {
	tx, err := txm.orm.FindTxByAttempt(hash)
	if err != nil {
		return nil, []models.TxAttempt{}, err
	}
	attempts, err := txm.orm.TxAttemptsFor(tx.ID)
	if err != nil {
		return nil, []models.TxAttempt{}, err
	}
	return tx, attempts, nil
}

// ContractLINKBalance returns the balance for the contract associated with this
// withdrawal request, or any errors
func (txm *EthTxManager) ContractLINKBalance(wr models.WithdrawalRequest) (assets.Link, error) {
	contractAddress := &wr.ContractAddress
	if (*contractAddress == common.Address{}) {
		if txm.config.OracleContractAddress() == nil {
			return assets.Link{}, errors.New(
				"OracleContractAddress not set; cannot check LINK balance")
		}
		contractAddress = txm.config.OracleContractAddress()
	}

	linkBalance, err := txm.GetLINKBalance(*contractAddress)
	if err != nil {
		return assets.Link{}, multierr.Combine(
			fmt.Errorf("Could not check LINK balance for %v",
				contractAddress),
			err)
	}
	return *linkBalance, nil
}

// WithdrawLINK withdraws the given amount of LINK from the contract to the
// configured withdrawal address. If wr.ContractAddress is empty (zero address),
// funds are withdrawn from configured OracleContractAddress.
func (txm *EthTxManager) WithdrawLINK(wr models.WithdrawalRequest) (common.Hash, error) {
	oracle, err := GetContract("Oracle")
	if err != nil {
		return common.Hash{}, err
	}

	data, err := oracle.EncodeMessageCall("withdraw", wr.DestinationAddress, (*big.Int)(wr.Amount))
	if err != nil {
		return common.Hash{}, err
	}

	contractAddress := &wr.ContractAddress
	if (*contractAddress == common.Address{}) {
		if txm.config.OracleContractAddress() == nil {
			return common.Hash{}, errors.New(
				"OracleContractAddress not set; cannot withdraw")
		}
		contractAddress = txm.config.OracleContractAddress()
	}

	tx, err := txm.CreateTx(*contractAddress, data)
	if err != nil {
		return common.Hash{}, err
	}

	return tx.Hash, nil
}

func (txm *EthTxManager) createAttempt(
	tx *models.Tx,
	gasPriceWei *big.Int,
	blkNum uint64,
) (*models.TxAttempt, error) {
	ma := txm.getAccount(tx.From)
	if ma == nil {
		return nil, fmt.Errorf("Unable to locate %v as an available account in EthTxManager. Has TxManager been started or has the address been removed?", tx.From.Hex())
	}
	etx := tx.EthTx(gasPriceWei)
	etx, err := txm.keyStore.SignTx(ma.Account, etx, txm.config.ChainID())
	if err != nil {
		return nil, err
	}

	a, err := txm.orm.AddTxAttempt(tx, etx, blkNum)
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
		return errors.Wrapf(err, "TxManager#sendTransaction with nonce %d", tx.Nonce())
	}
	return nil
}

type attemptState int

const (
	unconfirmed attemptState = iota
	confirmed
	safe
)

func (txm *EthTxManager) checkAttempt(
	tx *models.Tx,
	txat *models.TxAttempt,
	blkNum uint64,
) (*models.TxReceipt, attemptState, error) {
	receipt, err := txm.GetTxReceipt(txat.Hash)
	if err != nil {
		return nil, unconfirmed, errors.Wrap(err, "checkAttempt GetTxReceipt")
	}

	if receipt.Unconfirmed() {
		return txm.handleUnconfirmed(tx, txat, blkNum)
	}
	return txm.handleConfirmed(tx, txat, receipt, blkNum)
}

// GetETHAndLINKBalances attempts to retrieve the ethereum node's perception of
// the latest ETH and LINK balances for the active account on the txm, or an
// error on failure.
func (txm *EthTxManager) GetETHAndLINKBalances(address common.Address) (*big.Int, *assets.Link, error) {
	linkBalance, linkErr := txm.GetLINKBalance(address)
	ethBalance, ethErr := txm.EthClient.GetWeiBalance(address)
	merr := multierr.Append(linkErr, ethErr)
	return ethBalance, linkBalance, merr
}

// handleConfirmed checks whether a tx is confirmed, and records and reports it
// as such if so. Its bool return value is true if the tx is confirmed and it
// was successfully recorded as confirmed.
func (txm *EthTxManager) handleConfirmed(
	tx *models.Tx,
	txat *models.TxAttempt,
	rcpt *models.TxReceipt,
	blkNum uint64,
) (*models.TxReceipt, attemptState, error) {
	minConfs := big.NewInt(int64(txm.config.MinOutgoingConfirmations()))
	confirmedAt := big.NewInt(0).Add(minConfs, rcpt.BlockNumber.ToInt())
	confirmedAt.Sub(confirmedAt, big.NewInt(1)) // 0 based indexing since rcpt is 1 conf

	logger.Debugw(
		fmt.Sprintf("Confirmed TX: %d attempt %s waiting on %v confirmations", txat.TxID, txat.Hash.Hex(), minConfs),
		"txHash", txat.Hash.String(),
		"txid", txat.TxID,
		"nonce", tx.Nonce,
		"gasPrice", txat.GasPrice.String(),
		"from", tx.From.Hex(),
		"receiptBlockNumber", rcpt.BlockNumber.ToInt(),
		"currentBlockNumber", blkNum,
		"receiptHash", rcpt.Hash.Hex(),
		"confirmedAt", confirmedAt,
	)

	if big.NewInt(int64(blkNum)).Cmp(confirmedAt) == -1 {
		return nil, confirmed, nil
	}

	if err := txm.orm.MarkTxSafe(tx, txat); err != nil {
		return nil, confirmed, err
	}

	ethBalance, linkBalance, balanceErr := txm.GetETHAndLINKBalances(tx.From)
	logger.Infow(
		fmt.Sprintf("Confirmed TX %d: %s", txat.TxID, txat.Hash.String()),
		"txHash", txat.Hash.String(),
		"txid", txat.TxID,
		"nonce", tx.Nonce,
		"ethBalance", ethBalance,
		"linkBalance", linkBalance,
		"receipt", rcpt,
		"txat", txat,
		"err", balanceErr,
	)

	return rcpt, safe, nil
}

func (txm *EthTxManager) handleUnconfirmed(
	tx *models.Tx,
	txAttempt *models.TxAttempt,
	blkNum uint64,
) (*models.TxReceipt, attemptState, error) {
	if !isLatestAttempt(tx, txAttempt) {
		return nil, unconfirmed, nil
	}

	logParams := []interface{}{
		"txHash", txAttempt.Hash.String(),
		"txid", txAttempt.TxID,
		"nonce", tx.Nonce,
		"gasPrice", txAttempt.GasPrice.String(),
		"from", tx.From.Hex(),
		"blkNum", blkNum,
		"sentAt", txAttempt.SentAt,
	}
	if blkNum >= txAttempt.SentAt+txm.config.EthGasBumpThreshold() {
		logger.Debugw(
			fmt.Sprintf("Unconfirmed TX %d attempt %s, bumping gas", txAttempt.TxID, txAttempt.Hash.Hex()),
			logParams...,
		)
		return nil, unconfirmed, txm.bumpGas(txAttempt, blkNum)
	}
	logger.Infow(
		fmt.Sprintf("Unconfirmed TX %d has not met gas bump threshold", txAttempt.TxID),
		logParams...,
	)
	return nil, unconfirmed, nil
}

// isLatestAttempt returns true only if the attempt is the last
// attempt associated with the transaction, alluding to the fact that
// it has the highest gas price after subsequent bumps.
func isLatestAttempt(tx *models.Tx, txat *models.TxAttempt) bool {
	return tx.Hash == txat.Hash
}

func (txm *EthTxManager) bumpGas(txat *models.TxAttempt, blkNum uint64) error {
	tx, err := txm.orm.FindTx(txat.TxID)
	if err != nil {
		return errors.Wrapf(err, "bumpGas for txid %v", txat.TxID)
	}
	gasPrice := new(big.Int).Add(txat.GasPrice.ToInt(), txm.config.EthGasBumpWei())
	bumpedTxAt, err := txm.createAttempt(tx, gasPrice, blkNum)
	if err != nil {
		return errors.Wrapf(err, "bumpGas from tx %s", txat.Hash.Hex())
	}
	logger.Infow(fmt.Sprintf("Bumped gas to %v for TX %v", txat.TxID, gasPrice), "bumpSource", txat.Hash.Hex(), "txat", bumpedTxAt, "nonce", tx.Nonce)
	return nil
}

// NextActiveAccount uses round robin to select a managed account
// from the list of available accounts as defined in Register(...)
func (txm *EthTxManager) NextActiveAccount() *ManagedAccount {
	txm.accountsMutex.Lock()
	defer txm.accountsMutex.Unlock()

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
	txm.accountsMutex.Lock()
	defer txm.accountsMutex.Unlock()

	for _, a := range txm.availableAccounts {
		if a.Address == from {
			return a
		}
	}

	return nil
}

// ActivateAccount retrieves an account's nonce from the blockchain for client
// side management in ManagedAccount.
func (txm *EthTxManager) activateAccount(account accounts.Account) (*ManagedAccount, error) {
	nonce, err := txm.GetNonce(account.Address)
	if err != nil {
		return nil, err
	}

	return NewManagedAccount(account, nonce), nil
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

// Nonce returns the client side managed nonce.
func (a *ManagedAccount) Nonce() uint64 {
	return a.nonce
}

// ReloadNonce fetch and update the current nonce via eth_getTransactionCount
func (a *ManagedAccount) ReloadNonce(txm *EthTxManager) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	nonce, err := txm.GetNonce(a.Address)
	if err != nil {
		return fmt.Errorf("TxManager ReloadNonce: %v", err)
	}
	a.nonce = nonce
	return nil
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

func (txm *EthTxManager) getBlockNumber() (uint64, error) {
	if !txm.Connected() {
		return 0, errors.Wrap(ErrPendingConnection, "EthTxManager#getBlockNumber")
	}

	return txm.GetBlockNumber()
}

// Contract holds the solidity contract's parsed ABI
type Contract struct {
	ABI abi.ABI
}

// GetContract loads the contract JSON file from ../evm/build/contracts
// and parses the ABI JSON contents into an abi.ABI object
func GetContract(name string) (*Contract, error) {
	box := packr.NewBox("../../evm/build/contracts")
	jsonFile, err := box.Find(name + ".json")
	if err != nil {
		return nil, errors.New("unable to read contract JSON")
	}

	abiBytes := gjson.GetBytes(jsonFile, "abi")
	abiParsed, err := abi.JSON(strings.NewReader(abiBytes.Raw))
	if err != nil {
		return nil, err
	}

	return &Contract{abiParsed}, nil
}

// EncodeMessageCall encodes method name and arguments into a byte array
// to conform with the contract's ABI
func (contract *Contract) EncodeMessageCall(method string, args ...interface{}) ([]byte, error) {
	return contract.ABI.Pack(method, args...)
}
