package store

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"regexp"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/tevino/abool"
	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v3"
)

const (
	// Linear backoff is used so worst-case transaction time increases quadratically with this number
	nonceReloadLimit int = 3

	// The base time for the backoff
	nonceReloadBackoffBaseTime = 3 * time.Second
)

var (
	// ErrPendingConnection is the error returned if TxManager is not connected.
	ErrPendingConnection = errors.New("Cannot talk to chain, pending connection")

	promNumGasBumps = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tx_manager_num_gas_bumps",
		Help: "Number of gas bumps",
	})

	promGasBumpExceedsLimit = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tx_manager_gas_bump_exceeds_limit",
		Help: "Number of times gas bumping failed from exceeding the configured limit. Any counts of this type indicate a serious problem.",
	})

	promGasBumpUnderpricedReplacement = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tx_manager_gas_bump_underpriced_replacement",
		Help: "Number of underpriced replacement errors received while trying to bump gas. Counts of this type most likely indicate some kind of misconfiguration or problem.",
	})

	promTxAttemptFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tx_manager_tx_attempt_failed",
		Help: "Number of tx attempts that failed. Tx attempts should not fail in normal operation.",
	})
)

//go:generate mockery --name TxManager --output ../internal/mocks/ --case=underscore

// TxManager represents an interface for interacting with the blockchain
type TxManager interface {
	HeadTrackable
	Connected() bool
	Register(accounts []accounts.Account)

	CreateTx(to common.Address, data []byte) (*models.Tx, error)
	CreateTxWithGas(surrogateID null.String, to common.Address, data []byte, gasPriceWei *big.Int, gasLimit uint64) (*models.Tx, error)
	CreateTxWithEth(from, to common.Address, value *assets.Eth) (*models.Tx, error)
	CheckAttempt(txAttempt *models.TxAttempt, blockHeight uint64) (*types.Receipt, AttemptState, error)

	BumpGasUntilSafe(hash common.Hash) (*types.Receipt, AttemptState, error)

	ContractLINKBalance(wr models.WithdrawalRequest) (assets.Link, error)
	GetLINKBalance(address common.Address) (*assets.Link, error)
	NextActiveAccount() *ManagedAccount

	SignedRawTxWithBumpedGas(originalTx models.Tx, gasLimit uint64, gasPrice big.Int) ([]byte, error)

	eth.Client
}

// EthTxManager contains fields for the Ethereum client, the KeyStore,
// the local Config for the application, and the database.
type EthTxManager struct {
	eth.Client
	keyStore            KeyStoreInterface
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
func NewEthTxManager(client eth.Client, config orm.ConfigReader, keyStore KeyStoreInterface, orm *orm.ORM) *EthTxManager {
	return &EthTxManager{
		Client:        client,
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
	if txm.config.EnableBulletproofTxManager() {
		return nil
	}
	var merr error
	func() {
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
	}()

	// Upon connecting/reconnecting, rebroadcast any transactions that are still unconfirmed
	attempts, err := txm.orm.UnconfirmedTxAttempts()
	if err != nil {
		merr = multierr.Append(merr, err)
		return merr
	}

	attempts = models.HighestPricedTxAttemptPerTx(attempts)

	for _, attempt := range attempts {
		ma := txm.getAccount(attempt.Tx.From)
		if ma == nil {
			logger.Warnf("Trying to rebroadcast tx %v, could not find account %v", attempt.Hash.Hex(), attempt.Tx.From.Hex())
			continue
		} else if ma.Nonce() > attempt.Tx.Nonce {
			// Do not rebroadcast txs with nonces that are lower than our current nonce
			continue
		}

		logger.Infof("Rebroadcasting tx %v", attempt.Hash.Hex())

		_, err = txm.SendRawTx(attempt.SignedRawTx)
		if err != nil && !isNonceTooLowError(err) {
			logger.Warnf("Failed to rebroadcast tx %v: %v", attempt.Hash.Hex(), err)
		}
	}

	return merr
}

// Disconnect marks this instance as disconnected.
func (txm *EthTxManager) Disconnect() {
	if txm.config.EnableBulletproofTxManager() {
		return
	}
	txm.connected.UnSet()
}

// OnNewLongestChain does nothing; exists to comply with interface.
func (txm *EthTxManager) OnNewLongestChain(head models.Head) {
	if txm.config.EnableBulletproofTxManager() {
		return
	}
	txm.currentHead = head
}

// CreateTx signs and sends a transaction to the Ethereum blockchain.
func (txm *EthTxManager) CreateTx(to common.Address, data []byte) (*models.Tx, error) {
	return txm.CreateTxWithGas(null.String{}, to, data, txm.config.EthGasPriceDefault(), txm.config.EthGasLimitDefault())
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

	return txm.createTx(null.String{}, ma, to, []byte{}, txm.config.EthGasPriceDefault(), txm.config.EthGasLimitDefault(), value)
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
		return config.EthGasPriceDefault(), config.EthGasLimitDefault()
	}

	if gasPriceWei == nil {
		gasPriceWei = config.EthGasPriceDefault()
	}

	if gasLimit == 0 {
		gasLimit = config.EthGasLimitDefault()
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

	for nrc := 0; nrc < nonceReloadLimit+1; nrc++ {
		tx, err := txm.sendInitialTx(surrogateID, ma, to, data, gasPriceWei, gasLimit, value)
		if err == nil {
			return tx, nil
		}

		if !isNonceTooLowError(err) {
			return nil, errors.Wrap(err, "TxManager#retryInitialTx sendInitialTx")
		}

		logger.Warnw(
			"Tx #0: another tx with this nonce already exists, will retry with network nonce",
			"nonce", tx.Nonce, "gasPriceWei", gasPriceWei, "gasLimit", gasLimit, "error", err.Error(),
		)

		// Linear backoff
		time.Sleep(time.Duration(nrc+1) * nonceReloadBackoffBaseTime)

		logger.Warnw(
			"Tx #0: another tx with this nonce already exists, retrying with network nonce",
			"nonce", tx.Nonce, "gasPriceWei", gasPriceWei, "gasLimit", gasLimit, "error", err.Error(),
		)

		err = ma.ReloadNonce(txm)
		if err != nil {
			return nil, errors.Wrap(err, "TxManager#retryInitialTx ReloadNonce")
		}
	}

	return nil, fmt.Errorf(
		"transaction reattempt limit reached for 'nonce is too low' error. Limit: %v",
		nonceReloadLimit,
	)
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

	var err error
	var tx *models.Tx

	err = ma.GetAndIncrementNonce(func(nonce uint64) error {
		blockHeight := uint64(txm.currentHead.Number)
		tx, err = txm.newTx(
			ma.Account,
			nonce,
			to,
			value.ToInt(),
			gasLimit,
			gasPriceWei,
			data,
			&ma.Address,
			blockHeight,
		)
		if err != nil {
			return errors.Wrap(err, "TxManager#sendInitialTx newTx")
		}

		tx.SurrogateID = surrogateID
		tx, err = txm.orm.CreateTx(tx)
		if err != nil {
			return errors.Wrap(err, "TxManager#sendInitialTx CreateTx")
		}

		_, err = txm.SendRawTx(tx.SignedRawTx)
		if err != nil {
			return errors.Wrap(err, "TxManager#sendInitialTx SendRawTx")
		}

		txAttempt, e := txm.orm.AddTxAttempt(tx, tx)
		if e != nil {
			return errors.Wrap(e, "TxManager#sendInitialTx AddTxAttempt")
		}

		logger.Debugw("Added Tx attempt #0", "txID", tx.ID, "txAttemptID", txAttempt.ID)

		return nil
	})

	return tx, err
}

var (
	nonceTooLowRegex                       = regexp.MustCompile("(nonce .*too low|same hash was already imported|replacement transaction underpriced)")
	replacementTransactionUnderpricedRegex = regexp.MustCompile("replacement transaction underpriced")
)

// FIXME: There are probably other types of errors here that are symptomatic of a nonce that is too low
func isNonceTooLowError(err error) bool {
	return err != nil && nonceTooLowRegex.MatchString(err.Error())
}

func isUnderPricedReplacementError(err error) bool {
	return err != nil && replacementTransactionUnderpricedRegex.MatchString(err.Error())
}

// SignedRawTxWithBumpedGas takes a transaction and generates a new signed TX from it with the provided params
func (txm *EthTxManager) SignedRawTxWithBumpedGas(originalTx models.Tx, gasLimit uint64, gasPrice big.Int) ([]byte, error) {
	ma := txm.getAccount(originalTx.From)
	rlp := new(bytes.Buffer)
	if ma == nil {
		return nil, fmt.Errorf("unable to locate %v as an available account in EthTxManager. Has TxManager been started or has the address been removed?", originalTx.From.Hex())
	}

	transaction := types.NewTransaction(originalTx.Nonce, originalTx.To, originalTx.Value.ToInt(), gasLimit, &gasPrice, originalTx.Data)

	transaction, err := txm.keyStore.SignTx(ma.Account, transaction, txm.config.ChainID())
	if err != nil {
		return nil, err
	}

	if err := transaction.EncodeRLP(rlp); err != nil {
		return nil, err
	}
	return rlp.Bytes(), nil
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
		SignedRawTx: rlp.Bytes(),
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

// BumpGasUntilSafe process a collection of related TxAttempts, trying to get
// at least one TxAttempt into a safe state, bumping gas if needed
func (txm *EthTxManager) BumpGasUntilSafe(hash common.Hash) (*types.Receipt, AttemptState, error) {
	tx, _, err := txm.orm.FindTxByAttempt(hash)
	if err != nil {
		return nil, Unknown, errors.Wrap(err, "BumpGasUntilSafe FindTxByAttempt")
	}

	receipt, state, err := txm.checkChainForConfirmation(tx)
	if err != nil || state != Unconfirmed {
		return receipt, state, err
	}

	return txm.checkAccountForConfirmation(tx)
}

func (txm *EthTxManager) checkChainForConfirmation(tx *models.Tx) (*types.Receipt, AttemptState, error) {
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

func (txm *EthTxManager) checkAccountForConfirmation(tx *models.Tx) (*types.Receipt, AttemptState, error) {
	ma := txm.GetAvailableAccount(tx.From)

	if ma != nil && ma.lastSafeNonce > tx.Nonce {
		tx.Confirmed = true
		tx.Hash = utils.EmptyHash
		if err := txm.orm.SaveTx(tx); err != nil {
			return nil, Safe, fmt.Errorf("BumpGasUntilSafe error saving Tx confirmation to the database")
		}
		return nil, Safe, fmt.Errorf("BumpGasUntilSafe a version of the Ethereum Transaction from %v with nonce %v was not recorded in the database", tx.From, tx.Nonce)
	}

	return nil, Unconfirmed, nil
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
			fmt.Errorf("could not check LINK balance for %v",
				contractAddress),
			err)
	}
	return *linkBalance, nil
}

// CheckAttempt retrieves a receipt for a TxAttempt, and check if it meets the
// minimum number of confirmations
func (txm *EthTxManager) CheckAttempt(txAttempt *models.TxAttempt, blockHeight uint64) (*types.Receipt, AttemptState, error) {
	receipt, err := txm.TransactionReceipt(context.TODO(), txAttempt.Hash)
	if err == ethereum.NotFound {
		return nil, Unconfirmed, nil
	} else if err != nil {
		return nil, Unknown, errors.Wrap(err, "CheckAttempt TransactionReceipt failed")
	}

	if models.ReceiptIsUnconfirmed(receipt) {
		return receipt, Unconfirmed, nil
	}

	minimumConfirmations := new(big.Int).SetUint64(txm.config.MinRequiredOutgoingConfirmations())
	confirmedAt := new(big.Int).Add(minimumConfirmations, receipt.BlockNumber)

	confirmedAt.Sub(confirmedAt, big.NewInt(1)) // confirmed at block counts as 1 conf

	if big.NewInt(int64(blockHeight)).Cmp(confirmedAt) < 0 {
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
	// Confirmed means that a transaction has had at least one confirmation, but
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
) (*types.Receipt, AttemptState, error) {
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
			"receiptBlockNumber", receipt.BlockNumber,
			"currentBlockNumber", blockHeight,
			"receiptHash", receipt.TxHash.Hex(),
			"jobRunId", jobRunID,
		)

		return receipt, state, nil

	case Unconfirmed:
		attemptLimit := txm.config.TxAttemptLimit()
		if attemptIndex >= int(attemptLimit) {
			logger.Warnw(
				fmt.Sprintf("Tx #%d is %s, has met TxAttemptLimit", attemptIndex, state),
				"txAttemptLimit", attemptLimit,
				"txHash", txAttempt.Hash.String(),
				"txID", txAttempt.TxID,
				"jobRunId", jobRunID,
			)
			return receipt, state, nil
		}

		if isLatestAttempt(tx, attemptIndex) && txm.hasTxAttemptMetGasBumpThreshold(tx, attemptIndex, blockHeight) {
			logger.Debugw(
				fmt.Sprintf("Tx #%d is %s, bumping gas", attemptIndex, state),
				"txHash", txAttempt.Hash.String(),
				"txID", txAttempt.TxID,
				"currentBlockNumber", blockHeight,
				"jobRunId", jobRunID,
			)
			err = txm.bumpGas(tx, attemptIndex, blockHeight)
		} else {
			logger.Debugw(
				fmt.Sprintf("Tx #%d is %s", attemptIndex, state),
				"txHash", txAttempt.Hash.String(),
				"txID", txAttempt.TxID,
				"jobRunId", jobRunID,
			)
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

	minimumConfirmations := txm.config.MinRequiredOutgoingConfirmations()
	logger.Infow(
		fmt.Sprintf("Tx #%d is safe", attemptIndex),
		"minimumConfirmations", minimumConfirmations,
		"txHash", txAttempt.Hash.String(),
		"txID", txAttempt.TxID,
	)

	return nil
}

func (txm *EthTxManager) BumpGasByIncrement(originalGasPrice *big.Int) *big.Int {
	return BumpGas(txm.config, originalGasPrice)
}

// BumpGas computes the next gas price to attempt as the largest of:
// - A configured percentage bump (ETH_GAS_BUMP_PERCENT) on top of the baseline price.
// - A configured fixed amount of Wei (ETH_GAS_PRICE_WEI) on top of the baseline price.
// The baseline price is the maximum of the previous gas price attempt and the node's current gas price.
func BumpGas(config orm.ConfigReader, originalGasPrice *big.Int) *big.Int {
	baselinePrice := max(originalGasPrice, config.EthGasPriceDefault())

	var priceByPercentage = new(big.Int)
	priceByPercentage.Mul(baselinePrice, big.NewInt(int64(100+config.EthGasBumpPercent())))
	priceByPercentage.Div(priceByPercentage, big.NewInt(100))

	var priceByIncrement = new(big.Int)
	priceByIncrement.Add(baselinePrice, config.EthGasBumpWei())

	bumpedGasPrice := max(priceByPercentage, priceByIncrement)
	if bumpedGasPrice.Cmp(config.EthMaxGasPriceWei()) > 0 {
		return config.EthMaxGasPriceWei()
	}
	return bumpedGasPrice
}

func max(a, b *big.Int) *big.Int {
	if a.Cmp(b) >= 0 {
		return a
	}
	return b
}

// bumpGas attempts a new transaction with an increased gas cost
func (txm *EthTxManager) bumpGas(tx *models.Tx, attemptIndex int, blockHeight uint64) error {
	txAttempt := tx.Attempts[attemptIndex]

	originalGasPrice := txAttempt.GasPrice.ToInt()

	bumpedGasPrice := txm.BumpGasByIncrement(originalGasPrice)

	for {
		promNumGasBumps.Inc()
		if bumpedGasPrice.Cmp(txm.config.EthMaxGasPriceWei()) > 0 {
			// NOTE: In the current design, a new tx attempt will be created even if this one returns error.
			// If we do hit this scenario, we will keep creating new attempts that are guaranteed to fail
			// until CHAINLINK_TX_ATTEMPT_LIMIT is reached
			promGasBumpExceedsLimit.Inc()
			err := fmt.Errorf("bumped gas price of %v would exceed maximum configured limit of %v, set by ETH_MAX_GAS_PRICE_WEI", bumpedGasPrice, txm.config.EthMaxGasPriceWei())
			logger.Error(err)
			return err
		}
		bumpedTxAttempt, err := txm.createAttempt(tx, bumpedGasPrice, blockHeight)
		if isUnderPricedReplacementError(err) {
			// This is not expected if we have bumped at least geth's required
			// amount.
			promGasBumpUnderpricedReplacement.Inc()
			logger.Warnw(fmt.Sprintf("Gas bump was rejected by ethereum node as underpriced, bumping again. Your value of ETH_GAS_BUMP_PERCENT (%v) may be set too low", txm.config.EthGasBumpPercent()),
				"originalGasPrice", originalGasPrice, "bumpedGasPrice", bumpedGasPrice,
			)
			bumpedGasPrice = txm.BumpGasByIncrement(bumpedGasPrice)
			continue
		}
		if err != nil {
			promTxAttemptFailed.Inc()
			e := errors.Wrapf(err, "bumpGas from Tx #%s", txAttempt.Hash.Hex())
			logger.Error(e)
			return e
		}

		logger.Infow(
			fmt.Sprintf("Tx #%d created with bumped gas %v", attemptIndex+1, bumpedGasPrice),
			"originalTxHash", txAttempt.Hash,
			"newTxHash", bumpedTxAttempt.Hash)

		return nil
	}
}

// createAttempt adds a new transaction attempt to a transaction record
func (txm *EthTxManager) createAttempt(
	tx *models.Tx,
	gasPriceWei *big.Int,
	blockHeight uint64,
) (*models.TxAttempt, error) {
	ma := txm.getAccount(tx.From)
	if ma == nil {
		return nil, fmt.Errorf("unable to locate %v as an available account in EthTxManager. Has TxManager been started or has the address been removed?", tx.From.Hex())
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
	nonce, err := txm.PendingNonceAt(context.TODO(), account.Address)
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
	nonce, err := txm.PendingNonceAt(context.TODO(), a.Address)
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
