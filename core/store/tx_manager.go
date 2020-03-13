package store

import (
	"bytes"
	"fmt"
	"math/big"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/pkg/errors"

	"chainlink/core/assets"
	"chainlink/core/eth"
	"chainlink/core/logger"
	"chainlink/core/store/models"
	"chainlink/core/store/orm"
	"chainlink/core/utils"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tevino/abool"
	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v3"
)

// DefaultGasLimit sets the default gas limit for outgoing transactions.
// if updating DefaultGasLimit, be sure it matches with the
// DefaultGasLimit specified in evm/test/Oracle_test.js
const DefaultGasLimit uint64 = 500000

// Linear backoff is used so worst-case transaction time increases quadratically with this number
const nonceReloadLimit int = 3

// The base time for the backoff
const nonceReloadBackoffBaseTime = 3 * time.Second

// ErrPendingConnection is the error returned if TxManager is not connected.
var ErrPendingConnection = errors.New("Cannot talk to chain, pending connection")

//go:generate mockery -name TxManager -output ../internal/mocks/ -case=underscore

// TxManager represents an interface for interacting with the blockchain
type TxManager interface {
	HeadTrackable
	Connected() bool
	Register(accounts []accounts.Account)

	CreateTx(to common.Address, data []byte) (*models.Tx, error)
	CreateTxWithGas(surrogateID null.String, to common.Address, data []byte, gasPriceWei *big.Int, gasLimit uint64) (*models.Tx, error)
	CreateTxWithEth(from, to common.Address, value *assets.Eth) (*models.Tx, error)
	CheckAttempt(txAttempt *models.TxAttempt, blockHeight uint64) (*eth.TxReceipt, AttemptState, error)

	BumpGasUntilSafe(hash common.Hash) (*eth.TxReceipt, AttemptState, error)

	ContractLINKBalance(wr models.WithdrawalRequest) (assets.Link, error)
	WithdrawLINK(wr models.WithdrawalRequest) (common.Hash, error)
	GetLINKBalance(address common.Address) (*assets.Link, error)
	NextActiveAccount() *ManagedAccount

	eth.Client
}

// EthTxManager contains fields for the Ethereum client, the KeyStore,
// the local Config for the application, and the database.
type EthTxManager struct {
	eth.Client
	keyStore            *KeyStore
	config              orm.ConfigReader
	orm                 *orm.ORM
	registeredAccounts  []accounts.Account
	availableAccounts   []*ManagedAccount
	availableAccountIdx int
	accountsMutex       *sync.Mutex
	connected           *abool.AtomicBool
	currentHead         models.Head

	chHeads  chan *models.Head
	chNewTxs chan *models.Tx
	chStop   chan struct{}

	processPendingTxs services.SleeperTask
	nonceTooLowErrors map[uint64]int

	createTxMu sync.Mutex
}

// NewEthTxManager constructs an EthTxManager using the passed variables and
// initializing internal variables.
func NewEthTxManager(client eth.Client, config orm.ConfigReader, keyStore *KeyStore, orm *orm.ORM) *EthTxManager {
	txm := &EthTxManager{
		Client:                client,
		config:                config,
		keyStore:              keyStore,
		orm:                   orm,
		accountsMutex:         &sync.Mutex{},
		connected:             abool.New(),
		chHeads:               make(chan *models.Head),
		chUnsentTxs:           make(chan *models.Tx),
		chStop:                make(chan struct{}),
		nonceTooLowErrors:     make(map[uint64]int),
		processPendingTxsTask: services.NewSleeperTask(),
	}
	return txm
}

func (txm *EthTxManager) Start() {
	go func() {
		defer close(txm.chDone)
		for {
			select {
			case <-txm.chStop:
				return

			case tx := <-txm.chUnsentTxs:
				txm.processUnsentTx(tx)

			case head := <-txm.chHeads:
				txm.processPendingTxs.WakeUp(func() {
					txm.processPendingTxs(head.Number)
				})
			}
		}
	}()
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
	var merr error
	txm.accountsMutex.Lock()
	defer txm.accountsMutex.Unlock()

	txm.availableAccounts = []*ManagedAccount{}
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

// OnNewHead kicks off a round of processing unconfirmed txs.
func (txm *EthTxManager) OnNewHead(head *models.Head) {
	txm.currentHead = *head
	select {
	case txm.chHeads <- *head:
	case txm.chStop:
	}
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
	tx, err := txm.createTx(surrogateID, ma, to, data, gasPriceWei, gasLimit, nil)
	if err != nil {
		return nil, err
	}

	select {
	case txm.chNewTxs <- tx:
	case <-txm.chStop:
	}
	return tx, nil
}

// CreateTxWithEth signs and sends a transaction with some ETH to transfer.
func (txm *EthTxManager) CreateTxWithEth(from, to common.Address, value *assets.Eth) (*models.Tx, error) {
	ma := txm.getAccount(from)
	if ma == nil {
		return nil, errors.New("account does not exist")
	}

	tx, err := txm.createTx(null.String{}, ma, to, []byte{}, txm.config.EthGasPriceDefault(), DefaultGasLimit, value)
	if err != nil {
		return nil, err
	}

	select {
	case txm.chNewTxs <- tx:
	case <-txm.chStop:
	}
	return tx, nil
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

	// Save the unsent tx to the DB
	blockHeight := uint64(txm.currentHead.Number)
	tx, err := txm.newTx(
		ma.Account,
		0,
		to,
		value.ToInt(),
		gasLimit,
		gasPriceWei,
		data,
		&ma.Address,
		0,
	)
	if err != nil {
		return nil, errors.Wrap(err, "TxManager#sendInitialTx newTx")
	}

	tx.SurrogateID = surrogateID
	tx, err = txm.orm.CreateTx(tx)
	if err != nil {
		return nil, errors.Wrap(err, "TxManager#sendInitialTx CreateTx")
	}
	return tx, nil
}

func (txm *EthTxManager) processPendingTxs(blockHeight uint64) {
	// Fetch all unconfirmed txs along with their txAttempts sorted descending by gas price.
	txs, err := p.orm.UnconfirmedTxsWithAttemptsByGasPrice()
	if err != nil {
		logger.Error(err)
		return
	}

TxLoop:
	for _, tx := range txs {
		for attemptIndex, txAttempt := range txAttempts {
			_, state, err := txm.processAttempt(tx, attemptIndex, blockHeight)
			if err != nil {
				logger.Warnf("Failed to rebroadcast tx %v: %v", attempt.Hash.Hex(), err)
				continue TxLoop
			}

			if state == Safe || state == Unconfirmable || state == Confirmed {
				continue TxLoop
			}
		}
	}

	// Fetch all unsent txs
	for _, tx := range txs {
		txm.processUnsentTx(tx)
	}
}

func (txm *EthTxManager) processUnsentTx(tx *models.Tx) {
	dbNonce, err := txm.orm.GetLastNonce(ma.Address)
	if err != nil {
		panic("TODO")
	}

	err = ma.ReloadNonce()
	if err != nil {
		panic("TODO")
	}

	// Try to guess the most likely nonce by fetching from the DB as well as from
	// the Ethereum node.  Highest wins.  If that nonce doesn't succeed, we adjust
	// it manually based on the error we've received until it does succeed or we
	// hit the retry limit.  When the retry limit is hit, the tx will be reattempted
	// upon receiving new heads.
	var likelyNonce uint64
	if ma.Nonce() > dbNonce {
		likelyNonce = ma.Nonce()
	} else {
		likelyNonce = dbNonce
	}

	for nrc := 0; nrc < nonceReloadLimit+1; nrc++ {
		tx.Nonce = likelyNonce

		txAttempt, err = txm.createAttempt(tx, tx.GasPrice, blockHeight)
		if isNonceTooLowError(err) || isUnderpricedReplacementError(err) {
			likelyNonce++
			continue

		} else if isNonceTooHighError(err) {
			if likelyNonce == 0 {
				panic("is this even possible?")
			}
			likelyNonce--
			continue

		} else if err != nil {
			logger.Errorw(
				"Tx #0: error sending new transaction",
				"nonce", tx.Nonce, "gasPriceWei", gasPriceWei, "gasLimit", gasLimit, "error", err.Error(),
			)
			return nil, errors.Wrap(err, "TxManager#sendInitialTx SendRawTx")
		}

		logger.Debugw("Added Tx attempt #0", "txID", tx.ID, "txAttemptID", txAttempt.ID)
		return
	}

	logger.Error(fmt.Errorf(
		"Transaction reattempt limit reached for 'nonce is too low' error. Limit: %v",
		nonceReloadLimit,
	))
}

var (
	nonceTooLowRegex            = regexp.MustCompile("nonce .*too low")
	nonceTooHighRegex           = regexp.MustCompile("nonce .*too high")
	underpricedReplacementRegex = regexp.MustCompile("(same hash was already imported|replacement transaction underpriced)")
)

// FIXME: There are probably other types of errors here that are symptomatic of a nonce that is too low
func isNonceTooLowError(err error) bool {
	return err != nil && nonceTooLowRegex.MatchString(err.Error())
}

func isNonceTooHighError(err error) bool {
	return err != nil && nonceTooHighRegex.MatchString(err.Error())
}

func isUnderpricedReplacementError(err error) bool {
	return err != nil && underpricedReplacementRegex.MatchString(err.Error())
}

// newTx returns a newly signed Ethereum Transaction
func (txm *EthTxManager) newTx(
	account accounts.Account,
	nonce uint64,
	to common.Address,
	amount *big.Int,
	gasLimit uint64,
	gasPrice *big.Int,
	data []byte,
	from *common.Address,
	sentAt uint64) (*models.Tx, error) {

	transaction := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)

	transaction, err := txm.keyStore.SignTx(account, transaction, txm.config.ChainID())
	if err != nil {
		return nil, errors.Wrap(err, "TxManager newTx.SignTx")
	}

	rlp := new(bytes.Buffer)
	if err := transaction.EncodeRLP(rlp); err != nil {
		return nil, errors.Wrap(err, "TxManager newTx.EncodeRLP")
	}

	return &models.Tx{
		From:        *from,
		SentAt:      sentAt,
		To:          *transaction.To(),
		Nonce:       transaction.Nonce(),
		Data:        transaction.Data(),
		Value:       utils.NewBig(transaction.Value()),
		GasLimit:    transaction.Gas(),
		GasPrice:    utils.NewBig(transaction.GasPrice()),
		Hash:        transaction.Hash(),
		SignedRawTx: hexutil.Encode(rlp.Bytes()),
	}, nil
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

// GetAvailableAccount retrieves a managed account if it one matches the address given.
func (txm *EthTxManager) GetAvailableAccount(from common.Address) *ManagedAccount {
	for _, a := range txm.availableAccounts {
		if a.Address == from {
			return a
		}
	}
	return nil
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
	oracle, err := eth.GetContract("Oracle")
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
func (txm *EthTxManager) CheckAttempt(txAttempt *models.TxAttempt, blockHeight uint64) (*eth.TxReceipt, AttemptState, error) {
	receipt, err := txm.GetTxReceipt(txAttempt.Hash)
	if err != nil {
		return nil, Unknown, errors.Wrap(err, "CheckAttempt GetTxReceipt failed")
	}

	if receipt.Unconfirmed() {
		return receipt, Unconfirmed, nil
	}

	minimumConfirmations := new(big.Int).SetUint64(txm.config.MinOutgoingConfirmations())
	confirmedAt := new(big.Int).Add(minimumConfirmations, receipt.BlockNumber.ToInt())

	confirmedAt.Sub(confirmedAt, big.NewInt(1)) // confirmed at block counts as 1 conf

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
	// Unsent means we couldn't send the transaction right away when it was
	// created, and it should be reattempted
	Unsent
	// Unconfirmed means that a transaction has had no confirmations at all
	Unconfirmed
	// Confirmed means that a transaftion has had at least one transaction, but
	// not enough to satisfy the minimum number of confirmations configuration
	// option
	Confirmed
	// Safe has the required number of confirmations or more
	Safe
	// Unconfirmable means that we've exceeded the threshold for retrying a
	// given transaction and have given up on it.
	Unconfirmable
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
) (*eth.TxReceipt, AttemptState, error) {
	if tx.Unconfirmable {
		return nil, Unconfirmable, nil
	}

	jobRunID := tx.SurrogateID.ValueOrZero()
	txAttempt := tx.Attempts[attemptIndex]

	receipt, state, err := txm.CheckAttempt(txAttempt, blockHeight)

	switch state {
	case Safe:
		txm.updateLastSafeNonce(tx)
		return receipt, state, txm.handleSafe(tx, attemptIndex)

	case Confirmed:
		logger.Debugw(
			fmt.Sprintf("Tx #%d is %s", attemptIndex, state),
			"txHash", txAttempt.Hash.String(),
			"txID", txAttempt.TxID,
			"receiptBlockNumber", receipt.BlockNumber.ToInt(),
			"currentBlockNumber", blockHeight,
			"receiptHash", receipt.Hash.Hex(),
			"jobRunId", jobRunID,
		)

		return receipt, state, nil

	case Unconfirmed:
		attemptLimit := txm.config.TxAttemptLimit()
		if len(tx.Attempts) >= int(attemptLimit) {
			logger.Warnw(
				fmt.Sprintf("Tx #%d is %s, has met TxAttemptLimit", attemptIndex, state),
				"txAttemptLimit", attemptLimit,
				"txHash", txAttempt.Hash.String(),
				"txID", txAttempt.TxID,
				"jobRunId", jobRunID,
			)
			err := txm.handleUnconfirmable(tx)
			return receipt, Unconfirmable, err
		}

		if !txm.hasTxAttemptMetGasBumpThreshold(tx, attemptIndex, blockHeight) {
			logger.Debugw(
				fmt.Sprintf("Tx #%d is %s", attemptIndex, state),
				"txHash", txAttempt.Hash.String(),
				"txID", txAttempt.TxID,
				"jobRunId", jobRunID,
			)
			return receipt, state, nil
		}
		logger.Debugw(
			fmt.Sprintf("Tx #%d is %s, bumping gas", attemptIndex, state),
			"txHash", txAttempt.Hash.String(),
			"txID", txAttempt.TxID,
			"currentBlockNumber", blockHeight,
			"jobRunId", jobRunID,
		)
		err = txm.bumpGas(tx, attemptIndex, blockHeight)
		if isNonceTooLowError(err) {
			// A tx with this nonce has already been included into a block.
			// This means either that:
			//    1. the tx has confirmed in between when we fetched its receipt and now (this
			//       includes the case where a chain reorganization has caused many of our pending
			//       txs to confirm at once)
			//    2. we're talking to a load balanced Ethereum provider.  The node we fetched the
			//       receipt from doesn't know about the tx, and one we're resubmitting to does.
			//    3. the account was used by something other than the CL node
			//
			// The safest thing to do in cases 1 and 2 is to simply wait for more certainty.
			// Case 3 is unrecoverable, so after receiving this error a certain number of times,
			// we mark the tx as Unconfirmable.
			txm.nonceTooLowErrors[tx.ID]++
			if txm.nonceTooLowErrors[tx.ID] > txm.config.MaxNonceTooLowErrors() {
				err := txm.orm.MarkTxUnconfirmable(tx)
				if err != nil {
					logger.Errorf("can't mark tx unconfirmable: %v", err)
					continue TxLoop
				}
				delete(txm.nonceTooLowErrors, tx.ID)
			}

		} else if isUnderpricedReplacementError(err) {
			// This tx has somehow already made it to the mempool, but with a higher gas price
			// than any of the attempts that we already know about.  Unfortunately, there's no
			// way to track the other attempt because we don't know its hash.
			// This means either that:
			//     1. we've failed to record a previous submission in the DB
			//     2. the account was used by something other than the CL node (possibly another CL node?)
			//
			// We need to investigate scenarios where this case occurs in production.  For now,
			// the only safe thing to do here is to mark it as Unconfirmable right away.
			err := txm.orm.MarkTxUnconfirmable(tx)
			if err != nil {
				logger.Errorf("can't mark tx unconfirmable: %v", err)
				continue TxLoop
			}
		}

		return receipt, state, err

	default:
		logger.Debugw(
			fmt.Sprintf("Tx #%d is %s, error fetching receipt", attemptIndex, state),
			"txHash", txAttempt.Hash.String(),
			"txID", txAttempt.TxID,
			"jobRunId", jobRunID,
			"error", err,
		)
		return nil, Unknown, errors.Wrap(err, "processAttempt CheckAttempt failed")
	}
}

func (txm *EthTxManager) updateLastSafeNonce(tx *models.Tx) {
	for _, a := range txm.availableAccounts {
		if tx.From == a.Address {
			a.updateLastSafeNonce(tx.Nonce)
		}
	}
}

// hasTxAttemptMetGasBumpThreshold returns true if the current block height
// exceeds the configured gas bump threshold, indicating that it is time for a
// new transaction attempt to be created with an increased gas price
func (txm *EthTxManager) hasTxAttemptMetGasBumpThreshold(
	txAttempt *models.TxAttempt,
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
	txAttempt *models.Tx) error {

	if err := txm.orm.MarkTxSafe(tx, txAttempt); err != nil {
		return errors.Wrap(err, "handleSafe MarkTxSafe failed")
	}

	minimumConfirmations := txm.config.MinOutgoingConfirmations()
	ethBalance, linkBalance, balanceErr := txm.GetETHAndLINKBalances(tx.From)

	logger.Infow(
		fmt.Sprintf("Tx #%d is safe", txAttempt.TxID),
		"minimumConfirmations", minimumConfirmations,
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

	newTxAttempt, err := txm.newTx(
		ma.Account,
		tx.Nonce,
		tx.To,
		tx.Value.ToInt(),
		tx.GasLimit,
		gasPriceWei,
		tx.Data,
		&ma.Address,
		blockHeight,
	)
	if err != nil {
		return nil, errors.Wrap(err, "createAttempt#newTx failed")
	}

	if _, err = txm.SendRawTx(newTxAttempt.SignedRawTx); err != nil {
		return nil, errors.Wrap(err, "createAttempt#SendRawTx failed")
	}

	txAttempt, err := txm.orm.AddTxAttempt(tx, newTxAttempt)
	if err != nil {
		return nil, errors.Wrap(err, "createAttempt#AddTxAttempt failed")
	}

	logger.Debugw(fmt.Sprintf("Added Tx attempt #%d", len(tx.Attempts)+1), "txID", tx.ID, "txAttemptID", txAttempt.ID)

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
	nonce         uint64
	lastSafeNonce uint64
	mutex         *sync.Mutex
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
	logger.Debugw("Got new network nonce", "nonce", nonce)
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

func (a *ManagedAccount) updateLastSafeNonce(latest uint64) {
	if latest > a.lastSafeNonce {
		a.lastSafeNonce = latest
	}
}
