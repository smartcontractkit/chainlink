package bulletprooftxmanager

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/eth"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	gethAccounts "github.com/ethereum/go-ethereum/accounts"
	gethCommon "github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const (
	// databasePollInterval indicates how long to wait each time before polling
	// the database for new eth_transactions to send
	databasePollInterval = 1 * time.Second

	// maxEthNodeRequestTime is the worst case time we will wait for a response
	// from the eth node before we consider it to be an error
	maxEthNodeRequestTime = 2 * time.Minute

	// EthBroadcaster advisory lock class ID
	ethBroadcasterAdvisoryLockClassID = 0
)

type EthBroadcaster interface {
	Start() error
	Stop() error

	ProcessUnbroadcastEthTransactions(models.Key) error
}

// fatal means this transaction can never be accepted even with a different nonce or higher gas price
type sendError struct {
	fatal bool
	err   error
}

func (f *sendError) Error() string {
	return f.err.Error()
}

func (f *sendError) StrPtr() *string {
	e := f.err.Error()
	return &e
}

func (s *sendError) Fatal() bool {
	return s != nil && s.fatal
}

// Geth/parity returns this error if a transaction with this nonce already
// exists either on-chain or in the mempool.
//
// There are two scenarios in which this can happen:
// 1. The private key has been used to send at least one transaction from another wallet
// 2. The chainlink node crashed before being able to save the broadcastAt timestamp, indicating
//    that we are trying to send the exact same transaction twice (but it was already mined into a block).
//
// We can know which it is, because if we crashed there will be an unfinishedEthTransaction in the database.
// TODO: Probably needs a unit test
func (s *sendError) isNonceAlreadyUsedError() bool {
	// TODO: Add parity error
	return s != nil && s.err != nil && (s.err.Error() == "nonce too low" || s.err.Error() == "replacement transaction underpriced")
}

// Geth/parity returns this error if the transaction is already in the node's mempool
func (s *sendError) isTransactionAlreadyInMempool() bool {
	// TODO: Needs parity errors here
	return s.err != nil && strings.HasPrefix(s.Error(), "known transaction:")
}

// TODO: Write doc
func (s *sendError) isTerminallyUnderpriced() bool {
	// TODO: geth/parity errors
	return s.err != nil && (s.Error() == "transaction underpriced")
}

func NewFatalSendError(s string) *sendError {
	return &sendError{err: errors.New(s), fatal: true}
}

func FatalSendError(e error) *sendError {
	if e == nil {
		return nil
	}
	return &sendError{err: e, fatal: true}
}

func SendError(e error) *sendError {
	if e == nil {
		return nil
	}
	fatal := isFatalSendError(e)
	return &sendError{err: e, fatal: fatal}
}

// ethBroadcaster monitors eth_transactions for transactions that need to
// be broadcast, assigns nonces and ensures that at least one eth node
// somewhere has received the transaction successfully.
//
// This does not guarantee delivery! A whole host of other things can
// subsequently go wrong such as transctions being evicted from the mempool,
// eth nodes going offline etc. Responsibility for ensuring eventual inclusion
// into the chain falls on the shoulders of the ethConfirmer.
//
// What ethBroadcaster does guarantee is:
// - a monotic series of increasing nonces for eth_transactions that can be confirmed if you retry enough times
// - existence of a saved eth_transaction_attempt
type ethBroadcaster struct {
	store             *store.Store
	gethClientWrapper store.GethClientWrapper
	config            orm.ConfigReader

	started    bool
	stateMutex sync.RWMutex

	chStop chan struct{}
	chDone chan struct{}
}

func NewEthBroadcaster(store *store.Store, gethClientWrapper store.GethClientWrapper, config orm.ConfigReader) EthBroadcaster {
	return &ethBroadcaster{
		store:             store,
		gethClientWrapper: gethClientWrapper,
		config:            config,
		chStop:            make(chan struct{}),
		chDone:            make(chan struct{}),
	}
}

func (eb *ethBroadcaster) Start() error {
	if !eb.config.EnableBulletproofTxManager() {
		return nil
	}

	eb.stateMutex.Lock()
	defer eb.stateMutex.Unlock()
	if eb.started {
		return errors.New("already started")
	}
	go eb.monitorEthTransactions()
	eb.started = true

	return nil
}

func (eb *ethBroadcaster) Stop() error {
	eb.stateMutex.Lock()
	defer eb.stateMutex.Unlock()
	if !eb.started {
		return nil
	}
	eb.started = false
	close(eb.chStop)
	<-eb.chDone

	return nil
}

func (eb *ethBroadcaster) monitorEthTransactions() {
	defer close(eb.chDone)
	for {
		pollDatabaseTimer := time.NewTimer(databasePollInterval)

		keys, err := eb.store.Keys()

		if err != nil {
			logger.Error(err)
		} else {
			var wg sync.WaitGroup

			// It is safe to process separate keys concurrently
			// NOTE: This design will block one key if another takes a really long time to execute
			for _, key := range keys {
				if key == nil {
					logger.Error("key was unexpectedly nil. This should never happen")
					continue
				}
				wg.Add(1)
				go func(k models.Key) {
					if err := eb.ProcessUnbroadcastEthTransactions(k); err != nil {
						// NOTE: retries if this function errors are unbounded,
						// since they can be due to things like network errors
						// etc
						logger.Error(err)
					}
					wg.Done()
				}(*key)
			}

			wg.Wait()
		}

		select {
		case <-eb.chStop:
			return
		// TODO: can add <-eb.trigger channel for allowing other goroutines to manually trigger it early
		case <-pollDatabaseTimer.C:
			continue
		}
	}
}

func (eb *ethBroadcaster) ProcessUnbroadcastEthTransactions(key models.Key) error {
	ctx := context.Background()
	conn, err := eb.store.GetRawDB().DB().Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := eb.lock(ctx, conn, key.ID); err != nil {
		return err
	}
	defer eb.unlock(ctx, conn, key.ID)
	return eb.processUnbroadcastEthTransactions(key.Address.Address())
}

// TODO: write this doc
// NOTE: This MUST NOT be run concurrently for the same key or it will break things!
// TODO: Enforce this with a specific advisory lock (and move away from the global one)
func (eb *ethBroadcaster) processUnbroadcastEthTransactions(fromAddress gethCommon.Address) error {
	logger.Debugf("ProcessUnbroadcastEthTransactions start for %s", fromAddress.Hex())

	if err := eb.handleAnyUnfinishedEthTransaction(fromAddress); err != nil {
		return err
	}

	for {
		etx, err := nextUnbroadcastTransactionWithNonce(eb.store, fromAddress)
		if err != nil {
			// Break loop
			return err
		}
		if etx == nil {
			logger.Debugf("ProcessUnbroadcastEthTransactions finish for %s", fromAddress.Hex())
			// Finished
			return nil
		}

		gasPrice := eb.config.EthGasPriceDefault()
		etxAttempt := &models.EthTransactionAttempt{}
		sendError := eb.send(etx, etxAttempt, gasPrice)

		if sendError.Fatal() {
			etx.Error = sendError.StrPtr()
			err := saveFatallyErroredTransaction(eb.store, etx)
			if err != nil {
				return err
			}
			continue
		} else if sendError.isNonceAlreadyUsedError() {
			if err := eb.handleExternalWalletUsedNonce(etx, etxAttempt); err != nil {
				return err
			}
			continue
		} else if sendError != nil {
			return sendError.err
		}

		if err := saveBroadcastTransaction(eb.store, etx, etxAttempt); err != nil {
			return err
		}
	}
}

// TODO: docs
func (eb *ethBroadcaster) handleAnyUnfinishedEthTransaction(fromAddress gethCommon.Address) error {
	unfinishedEthTransaction, err := getUnfinishedEthTransaction(eb.store, fromAddress)
	if err != nil {
		return err
	}
	if unfinishedEthTransaction != nil {
		if err := eb.handleUnfinishedEthTransaction(unfinishedEthTransaction); err != nil {
			return err
		}
	}
	return nil
}

// TODO: Document exactly what the potential implications are from this
func (eb *ethBroadcaster) handleExternalWalletUsedNonce(etx *models.EthTransaction, etxAttempt *models.EthTransactionAttempt) error {
	// At all costs we avoid possible gaps in the nonce sequence. This means we may fail to send transactions, or send them twice and have one revert
	logger.Errorf("nonce of %v was too low for eth_transaction %v. Address %s has been used by another wallet. This is NOT SUPPORTED by chainlink and can lead to lost or reverted transactions.", *etx.Nonce, etx.ID, etx.FromAddress.String())

	clonedEtx := cloneForRebroadcast(etx)

	return eb.store.Transaction(func(db *gorm.DB) error {
		// Handle this case by assuming the particular transaction is broadcast already and handing off to the confirmer
		// We MUST do this to avoid gaps in the nonce sequence

		// We cannot know when the transaction was broadcast so just assume it was at the time of creation
		broadcastAt := etx.CreatedAt
		etx.BroadcastAt = &broadcastAt
		if err := saveBroadcastTransaction(eb.store, etx, etxAttempt); err != nil {
			return err
		}
		return db.Save(&clonedEtx).Error
	})
}

// getUnfinishedEthTransaction returns either 0 or 1 transaction that was left in
// an unfinished state because something went screwy the last time. Most likely
// the node crashed in the middle of the ProcessUnbroadcastEthTransactions loop.
// It may or may not have been broadcast to an eth node.
func getUnfinishedEthTransaction(store *store.Store, fromAddress gethCommon.Address) (*models.EthTransaction, error) {
	etx := &models.EthTransaction{}
	err := store.GetRawDB().First(etx, "from_address = ? AND broadcast_at IS NULL AND nonce IS NOT NULL", fromAddress.Bytes()).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return etx, err
}

// TODO: docs
func (eb *ethBroadcaster) handleUnfinishedEthTransaction(ethTransaction *models.EthTransaction) error {
	gasPrice := eb.config.EthGasPriceDefault()
	ethTransactionAttempt := &models.EthTransactionAttempt{}

	sendError := eb.send(ethTransaction, ethTransactionAttempt, gasPrice)
	if sendError.Fatal() {
		errString := sendError.Error()
		ethTransaction.Error = &errString
		return saveFatallyErroredTransaction(eb.store, ethTransaction)
	} else if sendError.isNonceAlreadyUsedError() {
		logger.Warnf("A transaction with nonce %v has already been confirmed. Either the node crashed on a previous run, or address %s has been used by another wallet. Assuming transaction was sent successfully", *ethTransaction.Nonce, ethTransaction.FromAddress.String())
		// Cannot really know BroadcastAt for certain since the node could have crashed an indeterminate time ago
		// CreatedAt is our best guess
		// NOTE: Could add additional column 'started_at' to do better but probably not very important
		broadcastAt := ethTransaction.CreatedAt
		ethTransaction.BroadcastAt = &broadcastAt
		return saveBroadcastTransaction(eb.store, ethTransaction, ethTransactionAttempt)
	} else if sendError != nil {
		return sendError
	}

	return saveBroadcastTransaction(eb.store, ethTransaction, ethTransactionAttempt)
}

// TODO: Write short doc
func nextUnbroadcastTransactionWithNonce(store *store.Store, fromAddress gethCommon.Address) (*models.EthTransaction, error) {
	ethTransaction := &models.EthTransaction{}
	if err := findNextUnbroadcastTransactionFromAddress(store.GetRawDB(), ethTransaction, fromAddress); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// Finish. No more unbroadcasted transactions left to process. Hoorah!
			return nil, nil
		}
		return nil, err
	}

	nonce, err := GetNextNonce(store.GetRawDB(), ethTransaction.FromAddress)
	if err != nil {
		return nil, err
	}
	ethTransaction.Nonce = &nonce
	if err := store.GetRawDB().Save(ethTransaction).Error; err != nil {
		return nil, err
	}
	return ethTransaction, nil
}

func findNextUnbroadcastTransactionFromAddress(tx *gorm.DB, ethTransaction *models.EthTransaction, fromAddress gethCommon.Address) error {
	return tx.
		Where("nonce IS NULL AND error IS NULL AND broadcast_at IS NULL AND from_address = ?", fromAddress).
		Order("created_at ASC, id ASC").
		First(ethTransaction).
		Error
}

func saveBroadcastTransaction(store *store.Store, ethTransaction *models.EthTransaction, attempt *models.EthTransactionAttempt) error {
	if ethTransaction.BroadcastAt == nil {
		return errors.New("broadcastAt must be set")
	}
	if ethTransaction.Nonce == nil {
		return errors.New("nonce must be set")
	}
	// TODO: Convert these to use TransactionWithAdvisoryLock
	return store.Transaction(func(tx *gorm.DB) error {
		if err := IncrementNextNonce(tx, ethTransaction.FromAddress, *ethTransaction.Nonce); err != nil {
			return err
		}
		if err := tx.Save(ethTransaction).Error; err != nil {
			return err
		}
		return tx.Save(attempt).Error
	})
}

func saveTransactionWithoutNonce(store *store.Store, ethTransaction *models.EthTransaction) error {
	if ethTransaction.Nonce != nil {
		return errors.New("nonce must be nil")
	}
	if ethTransaction.BroadcastAt != nil {
		return errors.New("broadcastAt must be nil")
	}
	return store.GetRawDB().Save(ethTransaction).Error
}

func saveFatallyErroredTransaction(store *store.Store, ethTransaction *models.EthTransaction) error {
	if ethTransaction.Error == nil {
		return errors.New("error must be set")
	}
	if ethTransaction.Nonce == nil {
		return errors.New("expected transaction to have a nonce")
	}
	ethTransaction.Nonce = nil
	return store.GetRawDB().Save(ethTransaction).Error
}

// GetNextNonce returns keys.next_nonce for the given address
func GetNextNonce(db *gorm.DB, address gethCommon.Address) (int64, error) {
	var nonce *int64
	row := db.Raw("SELECT next_nonce FROM keys WHERE address = ?", address).Row()
	if err := row.Scan(&nonce); err != nil {
		logger.Error(err)
		return 0, err
	}
	return *nonce, nil
}

// IncrementNextNonce increments keys.next_nonce by 1
func IncrementNextNonce(db *gorm.DB, address gethCommon.Address, currentNonce int64) error {
	res := db.Exec("UPDATE keys SET next_nonce = next_nonce + 1 WHERE address = ? AND next_nonce = ?", address.Bytes(), currentNonce)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		// TODO: Should probably reload nonce from eth client in this case since some invariant has been violated and it's a complete disaster
		return errors.New("could not increment nonce because no rows matched query. Either the key is missing or the nonce has been modified by an external process. This is an unrecoverable error")
	}
	return nil
}

// TODO: Write this doc
// NOTE: it can modify the EthTransaction and the EthTransactionAttempt in
// memory but will not save them
// Returning error here indicates that it may succeed on retry
func (eb *ethBroadcaster) send(etx *models.EthTransaction, attempt *models.EthTransactionAttempt, initialGasPrice *big.Int) *sendError {
	if etx == nil || attempt == nil {
		return NewFatalSendError("etx and etxAttempt must be non-nil")
	}
	if etx.Nonce == nil {
		return NewFatalSendError("cannot send transaction without nonce")
	}
	account, err := eb.store.KeyStore.GetAccountByAddress(etx.FromAddress)
	if err != nil {
		return FatalSendError(errors.Wrapf(err, "Error getting account %s for transaction %v", etx.FromAddress.String(), etx.ID))
	}

	transaction := gethTypes.NewTransaction(uint64(*etx.Nonce), etx.ToAddress, etx.Value.ToInt(), etx.GasLimit, initialGasPrice, etx.EncodedPayload)
	signedTx, signedTxBytes, err := eb.signTx(account, transaction, eb.config.ChainID())
	if err != nil {
		return FatalSendError(errors.Wrapf(err, "Error using account %s to sign transaction %v", etx.FromAddress.String(), etx.ID))
	}

	attempt.SignedRawTx = signedTxBytes
	attempt.EthTransactionID = etx.ID
	attempt.GasPrice = *utils.NewBig(initialGasPrice)

	sendErr := sendTransaction(eb.gethClientWrapper, signedTx)
	broadcastAt := time.Now()

	if sendErr.Fatal() {
		return sendErr
	}

	etx.BroadcastAt = &broadcastAt

	if sendErr == nil {
		return nil
	}

	// Bump gas if necessary
	if sendErr.isTerminallyUnderpriced() {
		logger.Errorf("transaction %v was underpriced at %v wei. You should increase your configured ETH_GAS_PRICE_DEFAULT (currently set to %v wei)", etx.ID, initialGasPrice, eb.config.EthGasPriceDefault())
		newGasPrice := eb.bumpGas(initialGasPrice)
		logger.Infof("retrying transaction %v with new gas price of %v wei", etx.ID, newGasPrice.Int64())
		return eb.send(etx, attempt, newGasPrice)
	} else if sendErr.isTransactionAlreadyInMempool() {
		logger.Debugf("transaction %v already in mempool", etx.ID)
		return nil
	}
	return sendErr
}

func (eb *ethBroadcaster) signTx(account gethAccounts.Account, tx *gethTypes.Transaction, chainID *big.Int) (*gethTypes.Transaction, []byte, error) {
	signedTx, err := eb.store.KeyStore.SignTx(account, tx, chainID)
	if err != nil {
		return nil, nil, err
	}
	rlp := new(bytes.Buffer)
	if err := signedTx.EncodeRLP(rlp); err != nil {
		return nil, nil, err
	}
	return signedTx, rlp.Bytes(), nil

}

func sendTransaction(gethClientWrapper store.GethClientWrapper, signedTransaction *gethTypes.Transaction) *sendError {
	err := gethClientWrapper.GethClient(func(gethClient eth.GethClient) error {
		ctx, cancel := context.WithTimeout(context.Background(), maxEthNodeRequestTime)
		defer cancel()
		return gethClient.SendTransaction(ctx, signedTransaction)
	})

	return SendError(err)
}

// TODO: Need to handle 'nonce too high'

// Geth/parity returns these errors if the transaction failed in such a way that:
// 1. It can NEVER be included into a block
// 2. Resending the transaction will never change that outcome
// TODO: This probably should have unit tests
// TODO: Better name? Unconfirmable transaction error?
func isFatalSendError(err error) bool {
	if err == nil {
		return false
	}
	switch err.Error() {
	// Geth errors
	// See: https://github.com/ethereum/go-ethereum/blob/b9df7ecdc3d3685180ceb29665bab59e9f614da5/core/tx_pool.go#L516
	case "exceeds block gas limit", "invalid sender", "negative value", "oversized data", "gas uint64 overflow", "intrinsic gas too low":
		return true
	// TODO: Add parity here, and can we use error codes?
	// See: https://github.com/openethereum/openethereum/blob/master/rpc/src/v1/helpers/errors.rs#L420
	default:
		return false
	}
}

// GetDefaultAddress queries the database for the address of the primary default ethereum key
func GetDefaultAddress(store *store.Store) (gethCommon.Address, error) {
	defaultKey, err := getDefaultKey(store)
	if err != nil {
		return gethCommon.Address{}, err
	}
	return defaultKey.Address.Address(), err
}

// NOTE: We can add more advanced logic here later such as sorting by priority
// etc
func getDefaultKey(store *store.Store) (models.Key, error) {
	availableKeys, err := store.Keys()
	if err != nil {
		return models.Key{}, err
	}
	if len(availableKeys) == 0 {
		return models.Key{}, errors.New("no keys available")
	}
	return *availableKeys[0], nil
}

// TODO: Unit test?
func cloneForRebroadcast(etx *models.EthTransaction) models.EthTransaction {
	return models.EthTransaction{
		Nonce:          nil,
		FromAddress:    etx.FromAddress,
		ToAddress:      etx.ToAddress,
		EncodedPayload: etx.EncodedPayload,
		Value:          etx.Value,
		GasLimit:       etx.GasLimit,
		BroadcastAt:    nil,
	}
}

// TODO: This is copied from tx_manager which is suboptimal. Consider copying unit tests also.
// bumpGas returns a new gas price increased by the larger of:
// - A configured percentage bump (ETH_GAS_BUMP_PERCENT)
// - A configured fixed amount of Wei (ETH_GAS_PRICE_WEI)
func (eb *ethBroadcaster) bumpGas(originalGasPrice *big.Int) *big.Int {
	// Similar logic is used in geth
	// See: https://github.com/ethereum/go-ethereum/blob/8d7aa9078f8a94c2c10b1d11e04242df0ea91e5b/core/tx_list.go#L255
	// And: https://github.com/ethereum/go-ethereum/blob/8d7aa9078f8a94c2c10b1d11e04242df0ea91e5b/core/tx_pool.go#L171
	percentageMultiplier := big.NewInt(100 + int64(eb.config.EthGasBumpPercent()))
	minimumGasBumpByPercentage := new(big.Int).Div(
		new(big.Int).Mul(
			originalGasPrice,
			percentageMultiplier,
		),
		big.NewInt(100),
	)
	minimumGasBumpByIncrement := new(big.Int).Add(originalGasPrice, eb.config.EthGasBumpWei())
	if minimumGasBumpByIncrement.Cmp(minimumGasBumpByPercentage) < 0 {
		return minimumGasBumpByPercentage
	}
	return minimumGasBumpByIncrement
}

func (eb *ethBroadcaster) lock(ctx context.Context, conn *sql.Conn, keyID int32) error {
	gotLock := false
	rows, err := conn.QueryContext(ctx, "SELECT pg_try_advisory_lock($1, $2)", ethBroadcasterAdvisoryLockClassID, keyID)
	defer rows.Close()
	if err != nil {
		return err
	}
	gotRow := rows.Next()
	if !gotRow {
		return errors.New("query unexpectedly returned 0 rows")
	}
	if err := rows.Scan(&gotLock); err != nil {
		return err
	}
	if gotLock {
		return nil
	}
	return fmt.Errorf("could not get advisory lock for key %v", keyID)
}

func (eb *ethBroadcaster) unlock(ctx context.Context, conn *sql.Conn, keyID int32) error {
	_, err := conn.ExecContext(ctx, "SELECT pg_advisory_unlock($1, $2)", ethBroadcasterAdvisoryLockClassID, keyID)
	return err
}
