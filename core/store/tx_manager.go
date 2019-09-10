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
	"gopkg.in/guregu/null.v3"

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
const nonceReloadLimit int = 1

// ErrPendingConnection is the error returned if TxManager is not connected.
var ErrPendingConnection = errors.New("Cannot talk to chain, pending connection")

// TxManager represents an interface for interacting with the blockchain
type TxManager interface {
	HeadTrackable
	Connected() bool
	Register(accounts []accounts.Account)

	CreateTx(to common.Address, data []byte) (*models.Tx, error)
	CreateTxWithGas(surrogateID null.String, to common.Address, data []byte, gasPriceWei *big.Int, gasLimit uint64) (*models.Tx, error)
	CreateTxWithEth(from, to common.Address, value *assets.Eth) (*models.Tx, error)
	CheckAttempt(txAttempt *models.TxAttempt, blockHeight uint64) (*models.TxReceipt, AttemptState, error)

	BumpGasUntilSafe(hash common.Hash) (*models.TxReceipt, AttemptState, error)

	ContractLINKBalance(wr models.WithdrawalRequest) (assets.Link, error)
	WithdrawLINK(wr models.WithdrawalRequest) (common.Hash, error)
	GetLINKBalance(address common.Address) (*assets.Link, error)
	NextActiveAccount() *ManagedAccount

	GetEthBalance(address common.Address) (*assets.Eth, error)
	SubscribeToNewHeads(channel chan<- models.BlockHeader) (models.EthSubscription, error)
	GetBlockByNumber(hex string) (models.BlockHeader, error)
	SubscribeToLogs(channel chan<- models.Log, q ethereum.FilterQuery) (models.EthSubscription, error)
	GetLogs(q ethereum.FilterQuery) ([]models.Log, error)
	GetTxReceipt(common.Hash) (*models.TxReceipt, error)
	GetChainID() (*big.Int, error)
}

//go:generate mockgen -package=mocks -destination=../internal/mocks/tx_manager_mocks.go github.com/smartcontractkit/chainlink/core/store TxManager

// EthTxManager contains fields for the Ethereum client, the KeyStore,
// the local Config for the application, and the database.
type EthTxManager struct {
	EthClient
	keyStore            *KeyStore
	config              orm.ConfigReader
	orm                 *orm.ORM
	registeredAccounts  []accounts.Account
	availableAccounts   []*ManagedAccount
	availableAccountIdx int
	accountsMutex       *sync.Mutex
	connected           *abool.AtomicBool
	currentHead         models.Head
}

// NewEthTxManager constructs an EthTxManager using the passed variables and
// initializing internal variables.
func NewEthTxManager(client EthClient, config orm.ConfigReader, keyStore *KeyStore, orm *orm.ORM) *EthTxManager {
	return &EthTxManager{
		EthClient:     client,
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

	if bn != nil {
		txm.currentHead = *bn
	}
	txm.connected.Set()
	return merr
}

// Disconnect marks this instance as disconnected.
func (txm *EthTxManager) Disconnect() {
	txm.connected.UnSet()
}

// OnNewHead does nothing; exists to comply with interface.
func (txm *EthTxManager) OnNewHead(head *models.Head) {
	txm.currentHead = *head
}

// CreateTx signs and sends a transaction to the Ethereum blockchain.
func (txm *EthTxManager) CreateTx(to common.Address, data []byte) (*models.Tx, error) {
	return txm.CreateTxWithGas(null.String{}, to, data, txm.config.EthGasPriceDefault(), DefaultGasLimit)
}

// CreateTxWithGas signs and sends a transaction to the Ethereum blockchain.
func (txm *EthTxManager) CreateTxWithGas(surrogateID null.String, to common.Address, data []byte, gasPriceWei *big.Int, gasLimit uint64) (*models.Tx, error) {
	ma, err := txm.nextAccount()
	if err != nil {
		return nil, err
	}

	gasPriceWei, gasLimit = normalizeGasParams(gasPriceWei, gasLimit, txm.config)
	return txm.createTx(surrogateID, ma, to, data, gasPriceWei, gasLimit, nil)
}

// CreateTxWithEth signs and sends a transaction with some ETH to transfer.
func (txm *EthTxManager) CreateTxWithEth(from, to common.Address, value *assets.Eth) (*models.Tx, error) {
	ma := txm.getAccount(from)
	if ma == nil {
		return nil, errors.New("account does not exist")
	}

	return txm.createTx(null.String{}, ma, to, []byte{}, txm.config.EthGasPriceDefault(), DefaultGasLimit, value)
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

func normalizeGasParams(gasPriceWei *big.Int, gasLimit uint64, config orm.ConfigReader) (*big.Int, uint64) {
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

// createTx creates an ethereum transaction, and retries to submit the
// transaction if a nonce too low error is returned
func (txm *EthTxManager) createTx(
	surrogateID null.String,
	ma *ManagedAccount,
	to common.Address,
	data []byte,
	gasPriceWei *big.Int,
	gasLimit uint64,
	value *assets.Eth) (*models.Tx, error) {

	tx, err := txm.sendInitialTx(surrogateID, ma, to, data, gasPriceWei, gasLimit, value)
	for nrc := 0; isNonceTooLowError(err); nrc++ {
		logger.Warnw("Tx #0: nonce too low, retrying with network nonce")

		if nrc >= nonceReloadLimit {
			err = fmt.Errorf(
				"Transaction reattempt limit reached for 'nonce is too low' error. Limit: %v",
				nonceReloadLimit,
			)
			return nil, err
		}

		err = txm.retryInitialTx(tx, ma, gasPriceWei)
	}

	return tx, err
}

// sendInitialTx creates the initial Tx record + attempt for an Ethereum Tx,
// there should only ever be one of those for a "job"
func (txm *EthTxManager) sendInitialTx(
	surrogateID null.String,
	ma *ManagedAccount,
	to common.Address,
	data []byte,
	gasPriceWei *big.Int,
	gasLimit uint64,
	value *assets.Eth) (*models.Tx, error) {

	var tx *models.Tx
	err := ma.GetAndIncrementNonce(func(nonce uint64) error {
		ethTx, err := txm.newEthTx(
			ma.Account,
			nonce,
			to,
			value.ToInt(),
			gasLimit,
			gasPriceWei,
			data)
		if err != nil {
			return errors.Wrap(err, "TxManager#sendInitialTx newEthTx")
		}

		blockHeight := uint64(txm.currentHead.Number)
		tx, err = txm.orm.CreateTx(surrogateID, ethTx, &ma.Address, blockHeight)
		if err != nil {
			return errors.Wrap(err, "TxManager#sendInitialTx CreateTx")
		}

		logger.Debugw(fmt.Sprintf("Adding Tx attempt #%d", 0), "txID", tx.ID)

		_, err = txm.SendRawTx(tx.SignedRawTx)
		if err != nil {
			return errors.Wrap(err, "TxManager#sendInitialTx SendRawTx")
		}

		return nil
	})

	// XXX: Small subtlety here: return the tx as it is initialized and is passed to retryInitialTx
	if err != nil {
		err = errors.Wrap(err, "TxManager#sendInitialTx ma#GetAndIncrementNonce")
	}
	return tx, err
}

// retryInitialTx is used to update the Tx record and attempt for an Ethereum
// Tx when a Tx is reattempted because of a nonce too low error
func (txm *EthTxManager) retryInitialTx(
	tx *models.Tx,
	ma *ManagedAccount,
	gasPriceWei *big.Int) error {

	err := ma.ReloadNonce(txm)
	if err != nil {
		return errors.Wrap(err, "TxManager#retryInitialTx ReloadNonce")
	}

	err = ma.GetAndIncrementNonce(func(nonce uint64) error {
		ethTx, err := txm.newEthTx(
			ma.Account,
			nonce,
			tx.To,
			tx.Value.ToInt(),
			tx.GasLimit,
			gasPriceWei,
			tx.Data)
		if err != nil {
			return errors.Wrap(err, "TxManager#retryInitialTx newEthTx")
		}

		blockHeight := uint64(txm.currentHead.Number)
		err = txm.orm.UpdateTx(tx, ethTx, &ma.Address, blockHeight)
		if err != nil {
			return errors.Wrap(err, "TxManager#retryInitialTx UpdateTx")
		}

		_, err = txm.SendRawTx(tx.SignedRawTx)
		if err != nil {
			return errors.Wrap(err, "TxManager#retryInitialTx SendRawTx")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "TxManager#retryInitialTx ma#GetAndIncrementNonce")
	}
	return nil
}

var (
	nonceTooLowRegex = regexp.MustCompile("(nonce .*too low|same hash was already imported)")
)

func isNonceTooLowError(err error) bool {
	return err != nil && nonceTooLowRegex.MatchString(err.Error())
}

// newEthTx returns a newly signed Ethereum Transaction
func (txm *EthTxManager) newEthTx(
	account accounts.Account,
	nonce uint64,
	to common.Address,
	amount *big.Int,
	gasLimit uint64,
	gasPrice *big.Int,
	data []byte) (*types.Transaction, error) {

	ethTx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)

	ethTx, err := txm.keyStore.SignTx(account, ethTx, txm.config.ChainID())
	if err != nil {
		return nil, errors.Wrap(err, "TxManager keyStore.SignTx")
	}

	return ethTx, nil
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

// BumpGasUntilSafe process a collection of related TxAttempts, trying to get
// at least one TxAttempt into a safe state, bumping gas if needed
func (txm *EthTxManager) BumpGasUntilSafe(hash common.Hash) (*models.TxReceipt, AttemptState, error) {
	tx, _, err := txm.orm.FindTxByAttempt(hash)
	if err != nil {
		return nil, Unknown, errors.Wrap(err, "BumpGasUntilSafe FindTxByAttempt")
	}

	receipt, state, err := txm.checkChainForConfirmation(tx)
	if err != nil || state != Unconfirmed {
		return receipt, state, err
	}

	return txm.checkDBForConfirmation(tx)
}

func (txm *EthTxManager) checkChainForConfirmation(tx *models.Tx) (*models.TxReceipt, AttemptState, error) {
	blockHeight := uint64(txm.currentHead.Number)

	var merr error
	// Process attempts in reverse, since the attempt with the highest gas is
	// likely to be confirmed first
	for attemptIndex := len(tx.Attempts) - 1; attemptIndex >= 0; attemptIndex-- {
		receipt, state, err := txm.processAttempt(tx, attemptIndex, blockHeight)
		if state == Safe || state == Confirmed {
			return receipt, state, err // success, so all other attempt errors can be ignored.
		}
		merr = multierr.Append(merr, err)
	}

	return nil, Unconfirmed, merr
}

func (txm *EthTxManager) checkDBForConfirmation(tx *models.Tx) (*models.TxReceipt, AttemptState, error) {
	later, err := txm.orm.FindLaterConfirmedTx(tx)
	if err != nil {
		return nil, Unknown, errors.Wrap(err, "BumpGasUntilSafe checkDBForConfirmation")
	} else if later == nil {
		return nil, Unconfirmed, nil
	}

	tx.Confirmed = true
	tx.Hash = utils.EmptyHash
	if err = txm.orm.SaveTx(tx); err != nil {
		return nil, Confirmed, fmt.Errorf("BumpGasUntilSafe error saving Tx confirmation to the database")
	}

	err = fmt.Errorf("BumpGasUntilSafe a version of the Ethereum Transaction from %v with nonce %v", tx.From, tx.Nonce)
	return nil, Confirmed, err
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

// GetETHAndLINKBalances attempts to retrieve the ethereum node's perception of
// the latest ETH and LINK balances for the active account on the txm, or an
// error on failure.
func (txm *EthTxManager) GetETHAndLINKBalances(address common.Address) (*assets.Eth, *assets.Link, error) {
	linkBalance, linkErr := txm.GetLINKBalance(address)
	ethBalance, ethErr := txm.GetEthBalance(address)
	merr := multierr.Append(linkErr, ethErr)
	return ethBalance, linkBalance, merr
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

// CheckAttempt retrieves a receipt for a TxAttempt, and check if it meets the
// minimum number of confirmations
func (txm *EthTxManager) CheckAttempt(txAttempt *models.TxAttempt, blockHeight uint64) (*models.TxReceipt, AttemptState, error) {
	receipt, err := txm.GetTxReceipt(txAttempt.Hash)
	if err != nil {
		return nil, Unknown, errors.Wrap(err, "CheckAttempt GetTxReceipt failed")
	}

	if receipt.Unconfirmed() {
		return receipt, Unconfirmed, nil
	}

	minimumConfirmations := new(big.Int).SetUint64(txm.config.MinOutgoingConfirmations())
	confirmedAt := new(big.Int).Add(minimumConfirmations, receipt.BlockNumber.ToInt())

	// 0 based indexing since receipt is 1 conf
	confirmedAt.Sub(confirmedAt, big.NewInt(1))

	if new(big.Int).SetUint64(blockHeight).Cmp(confirmedAt) == -1 {
		return receipt, Confirmed, nil
	}

	return receipt, Safe, nil
}

// AttemptState enumerates the possible states of a transaction attempt as it
// gets accepted and confirmed by the blockchain
type AttemptState int

const (
	// Unknown is returned when the state of a transaction could not be
	// determined because of an error
	Unknown AttemptState = iota
	// Unconfirmed means that a transaction has had no confirmations at all
	Unconfirmed
	// Confirmed means that a transaftion has had at least one transaction, but
	// not enough to satisfy the minimum number of confirmations configuration
	// option
	Confirmed
	// Safe has the required number of confirmations or more
	Safe
)

// String conforms to the Stringer interface for AttemptState
func (a AttemptState) String() string {
	switch a {
	case Unconfirmed:
		return "unconfirmed"
	case Confirmed:
		return "confirmed"
	case Safe:
		return "safe"
	default:
		return "unknown"
	}
}

// processAttempt checks the state of a transaction attempt on the blockchain
// and decides if it is safe, needs bumping or more confirmations are needed to
// decide
func (txm *EthTxManager) processAttempt(
	tx *models.Tx,
	attemptIndex int,
	blockHeight uint64,
) (*models.TxReceipt, AttemptState, error) {
	txAttempt := tx.Attempts[attemptIndex]

	logger.Debugw(
		fmt.Sprintf("Tx #%d checking on-chain state", attemptIndex),
		"txHash", txAttempt.Hash.String(),
		"txID", txAttempt.TxID,
	)

	receipt, state, err := txm.CheckAttempt(txAttempt, blockHeight)
	if err != nil {
		return nil, Unknown, errors.Wrap(err, "processAttempt CheckAttempt failed")
	}

	logger.Debugw(
		fmt.Sprintf("Tx #%d is %s", attemptIndex, state),
		"txHash", txAttempt.Hash.String(),
		"txID", txAttempt.TxID,
		"receiptBlockNumber", receipt.BlockNumber.ToInt(),
		"currentBlockNumber", blockHeight,
		"receiptHash", receipt.Hash.Hex(),
	)

	switch state {
	case Safe:
		return receipt, state, txm.handleSafe(tx, attemptIndex)

	case Confirmed: // nothing to do, need to wait
		return receipt, state, nil

	case Unconfirmed:
		attemptLimit := txm.config.TxAttemptLimit()
		if attemptIndex >= int(attemptLimit) {
			logger.Warnw(
				fmt.Sprintf("Tx #%d has met TxAttemptLimit", attemptIndex),
				"txAttemptLimit", attemptLimit,
				"txHash", txAttempt.Hash.String(),
				"txID", txAttempt.TxID,
			)
			return receipt, state, nil
		}

		if isLatestAttempt(tx, attemptIndex) && txm.hasTxAttemptMetGasBumpThreshold(tx, attemptIndex, blockHeight) {
			logger.Debugw(
				fmt.Sprintf("Tx #%d has met gas bump threshold, bumping gas", attemptIndex),
				"txHash", txAttempt.Hash.String(),
				"txID", txAttempt.TxID,
			)
			err = txm.bumpGas(tx, attemptIndex, blockHeight)
		}

		return receipt, state, err
	}

	panic("invariant violated, 'Unknown' state returned without error")
}

// hasTxAttemptMetGasBumpThreshold returns true if the current block height
// exceeds the configured gas bump threshold, indicating that it is time for a
// new transaction attempt to be created with an increased gas price
func (txm *EthTxManager) hasTxAttemptMetGasBumpThreshold(
	tx *models.Tx,
	attemptIndex int,
	blockHeight uint64) bool {

	gasBumpThreshold := txm.config.EthGasBumpThreshold()
	txAttempt := tx.Attempts[attemptIndex]

	return blockHeight >= txAttempt.SentAt+gasBumpThreshold
}

// isLatestAttempt returns true only if the attempt is the last
// attempt associated with the transaction, alluding to the fact that
// it has the highest gas price after subsequent bumps.
func isLatestAttempt(tx *models.Tx, attemptIndex int) bool {
	return attemptIndex+1 == len(tx.Attempts)
}

// handleSafe marks a transaction as safe, no more work needs to be done
func (txm *EthTxManager) handleSafe(
	tx *models.Tx,
	attemptIndex int) error {
	txAttempt := tx.Attempts[attemptIndex]

	if err := txm.orm.MarkTxSafe(tx, txAttempt); err != nil {
		return errors.Wrap(err, "handleSafe MarkTxSafe failed")
	}

	minimumConfirmations := txm.config.MinOutgoingConfirmations()
	ethBalance, linkBalance, balanceErr := txm.GetETHAndLINKBalances(tx.From)

	logger.Infow(
		fmt.Sprintf("Tx #%d got minimum confirmations (%d)", attemptIndex, minimumConfirmations),
		"txHash", txAttempt.Hash.String(),
		"txID", txAttempt.TxID,
		"ethBalance", ethBalance,
		"linkBalance", linkBalance,
		"err", balanceErr,
	)

	return nil
}

// bumpGas creates a new transaction attempt with an increased gas cost
func (txm *EthTxManager) bumpGas(tx *models.Tx, attemptIndex int, blockHeight uint64) error {
	txAttempt := tx.Attempts[attemptIndex]

	originalGasPrice := txAttempt.GasPrice.ToInt()
	bumpedGasPrice := new(big.Int).Add(originalGasPrice, txm.config.EthGasBumpWei())

	bumpedTxAttempt, err := txm.createAttempt(tx, bumpedGasPrice, blockHeight)
	if err != nil {
		return errors.Wrapf(err, "bumpGas from Tx #%s", txAttempt.Hash.Hex())
	}

	logger.Infow(
		fmt.Sprintf("Tx #%d created with bumped gas %v", attemptIndex+1, bumpedGasPrice),
		"originalTxHash", txAttempt.Hash,
		"newTxHash", bumpedTxAttempt.Hash)
	return nil
}

// createAttempt adds a new transaction attempt to a transaction record
func (txm *EthTxManager) createAttempt(
	tx *models.Tx,
	gasPriceWei *big.Int,
	blockHeight uint64,
) (*models.TxAttempt, error) {
	ma := txm.getAccount(tx.From)
	if ma == nil {
		return nil, fmt.Errorf("Unable to locate %v as an available account in EthTxManager. Has TxManager been started or has the address been removed?", tx.From.Hex())
	}
	etx := tx.EthTx(gasPriceWei)
	etx, err := txm.keyStore.SignTx(ma.Account, etx, txm.config.ChainID())
	if err != nil {
		return nil, errors.Wrap(err, "createAttempt#SignTx failed")
	}

	logger.Debugw(fmt.Sprintf("Adding Tx attempt #%d", len(tx.Attempts)+1), "txID", tx.ID)

	txAttempt, err := txm.orm.AddTxAttempt(tx, etx, blockHeight)
	if err != nil {
		return nil, errors.Wrap(err, "createAttempt#AddTxAttempt failed")
	}

	if _, err = txm.SendRawTx(txAttempt.SignedRawTx); err != nil {
		return nil, errors.Wrap(err, "createAttempt#SendRawTx failed")
	}

	return txAttempt, nil
}

// NextActiveAccount uses round robin to select a managed account
// from the list of available accounts as defined in Register(...)
func (txm *EthTxManager) NextActiveAccount() *ManagedAccount {
	txm.accountsMutex.Lock()
	defer txm.accountsMutex.Unlock()

	if len(txm.availableAccounts) == 0 {
		return nil
	}

	account := txm.availableAccounts[txm.availableAccountIdx]
	txm.availableAccountIdx = (txm.availableAccountIdx + 1) % len(txm.availableAccounts)
	return account
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
