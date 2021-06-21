package bulletprooftxmanager

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// InFlightTransactionRecheckInterval controls how often the EthBroadcaster
// will poll the unconfirmed queue to see if it is allowed to send another
// transaction
const InFlightTransactionRecheckInterval = 1 * time.Second

// EthBroadcaster monitors eth_txes for transactions that need to
// be broadcast, assigns nonces and ensures that at least one eth node
// somewhere has received the transaction successfully.
//
// This does not guarantee delivery! A whole host of other things can
// subsequently go wrong such as transactions being evicted from the mempool,
// eth nodes going offline etc. Responsibility for ensuring eventual inclusion
// into the chain falls on the shoulders of the ethConfirmer.
//
// What EthBroadcaster does guarantee is:
// - a monotonic series of increasing nonces for eth_txes that can all eventually be confirmed if you retry enough times
// - transition of eth_txes out of unstarted into either fatal_error or unconfirmed
// - existence of a saved eth_tx_attempt
type EthBroadcaster struct {
	db             *gorm.DB
	ethClient      eth.Client
	config         Config
	keystore       KeyStore
	advisoryLocker postgres.AdvisoryLocker

	ethTxInsertListener postgres.Subscription
	eventBroadcaster    postgres.EventBroadcaster

	keys []ethkey.Key

	// triggers allow other goroutines to force EthBroadcaster to rescan the
	// database early (before the next poll interval)
	// Each key has its own trigger
	triggers map[gethCommon.Address]chan struct{}

	ctx       context.Context
	ctxCancel context.CancelFunc
	wg        sync.WaitGroup

	utils.StartStopOnce
}

// NewEthBroadcaster returns a new concrete EthBroadcaster
func NewEthBroadcaster(db *gorm.DB, ethClient eth.Client, config Config, keystore KeyStore, advisoryLocker postgres.AdvisoryLocker, eventBroadcaster postgres.EventBroadcaster, allKeys []ethkey.Key) *EthBroadcaster {
	ctx, cancel := context.WithCancel(context.Background())
	triggers := make(map[gethCommon.Address]chan struct{})
	return &EthBroadcaster{
		db:               db,
		ethClient:        ethClient,
		config:           config,
		keystore:         keystore,
		advisoryLocker:   advisoryLocker,
		eventBroadcaster: eventBroadcaster,
		keys:             allKeys,
		triggers:         triggers,
		ctx:              ctx,
		ctxCancel:        cancel,
		wg:               sync.WaitGroup{},
	}
}

func (eb *EthBroadcaster) Start() error {
	return eb.StartOnce("EthBroadcaster", func() (err error) {
		eb.ethTxInsertListener, err = eb.eventBroadcaster.Subscribe(postgres.ChannelInsertOnEthTx, "")
		if err != nil {
			return errors.Wrap(err, "EthBroadcaster could not start")
		}

		if eb.config.EthNonceAutoSync() {
			syncer := NewNonceSyncer(eb.db, eb.ethClient)
			if err := syncer.SyncAll(eb.ctx, eb.keys); err != nil {
				return errors.Wrap(err, "EthBroadcaster failed to sync with on-chain nonce")
			}
		}

		eb.wg.Add(len(eb.keys))
		for _, k := range eb.keys {
			triggerCh := make(chan struct{}, 1)
			eb.triggers[k.Address.Address()] = triggerCh
			go eb.monitorEthTxs(k, triggerCh)
		}

		eb.wg.Add(1)
		go eb.ethTxInsertTriggerer()

		return nil
	})
}

func (eb *EthBroadcaster) Close() error {
	return eb.StopOnce("EthBroadcaster", func() error {
		if eb.ethTxInsertListener != nil {
			eb.ethTxInsertListener.Close()
		}

		eb.ctxCancel()
		eb.wg.Wait()

		return nil
	})
}

// Trigger forces the monitor for a particular address to recheck for new eth_txes
// Logs error and does nothing if address was not registered on startup
func (eb *EthBroadcaster) Trigger(addr gethCommon.Address) {
	ok := eb.IfStarted(func() {
		triggerCh, exists := eb.triggers[addr]
		if !exists {
			var registeredAddrs []gethCommon.Address
			for addr := range eb.triggers {
				registeredAddrs = append(registeredAddrs, addr)
			}
			logger.Errorw(fmt.Sprintf("EthBroadcaster: attempted trigger for address %s which is not registered", addr.Hex()), "registeredAddrs", registeredAddrs)
			return
		}
		select {
		case triggerCh <- struct{}{}:
		default:
		}
	})

	if !ok {
		logger.Debugf("EthBroadcaster: unstarted; ignoring trigger for %s", addr.Hex())
	}
}

func (eb *EthBroadcaster) ethTxInsertTriggerer() {
	defer eb.wg.Done()
	for {
		select {
		case ev := <-eb.ethTxInsertListener.Events():
			hexAddr := ev.Payload
			address := gethCommon.HexToAddress(hexAddr)
			eb.Trigger(address)
		case <-eb.ctx.Done():
			return
		}
	}
}

func (eb *EthBroadcaster) monitorEthTxs(k ethkey.Key, triggerCh chan struct{}) {
	defer eb.wg.Done()
	for {
		pollDBTimer := time.NewTimer(utils.WithJitter(eb.config.TriggerFallbackDBPollInterval()))

		if err := eb.ProcessUnstartedEthTxs(k); err != nil {
			logger.Errorw("Error in ProcessUnstartedEthTxs", "error", err)
		}

		select {
		case <-eb.ctx.Done():
			// NOTE: See: https://godoc.org/time#Timer.Stop for an explanation of this pattern
			if !pollDBTimer.Stop() {
				<-pollDBTimer.C
			}
			return
		case <-triggerCh:
			// EthTx was inserted
			if !pollDBTimer.Stop() {
				<-pollDBTimer.C
			}
			continue
		case <-pollDBTimer.C:
			// DB poller timed out
			continue
		}
	}
}

func (eb *EthBroadcaster) ProcessUnstartedEthTxs(key ethkey.Key) error {
	return eb.advisoryLocker.WithAdvisoryLock(context.TODO(), postgres.AdvisoryLockClassID_EthBroadcaster, key.ID, func() error {
		return eb.processUnstartedEthTxs(key.Address.Address())
	})
}

// NOTE: This MUST NOT be run concurrently for the same address or it could
// result in undefined state or deadlocks.
// First handle any in_progress transactions left over from last time.
// Then keep looking up unstarted transactions and processing them until there are none remaining.
func (eb *EthBroadcaster) processUnstartedEthTxs(fromAddress gethCommon.Address) error {
	var n uint = 0
	mark := time.Now()
	defer func() {
		if n > 0 {
			logger.Debugw("EthBroadcaster: finished processUnstartedEthTxs", "address", fromAddress, "time", time.Since(mark), "n", n, "id", "eth_broadcaster")
		}
	}()

	if err := eb.handleAnyInProgressEthTx(fromAddress); err != nil {
		return errors.Wrap(err, "processUnstartedEthTxs failed")
	}
	for {
		maxInFlightTransactions := eb.config.EthMaxInFlightTransactions()
		if maxInFlightTransactions > 0 {
			nUnconfirmed, err := CountUnconfirmedTransactions(eb.db, fromAddress)
			if err != nil {
				return errors.Wrap(err, "CountUnconfirmedTransactions failed")
			}
			if nUnconfirmed >= maxInFlightTransactions {
				logger.Warnw(fmt.Sprintf(`EthBroadcaster: transaction throttling; maximum number of in-flight transactions is %d per key. If this happens a lot, you might need to increase ETH_MAX_IN_FLIGHT_TRANSACTIONS. %s`, maxInFlightTransactions, EthMaxInFlightTransactionsWarningLabel), "nUnconfirmed", nUnconfirmed)
				time.Sleep(InFlightTransactionRecheckInterval)
				continue
			}
		}
		etx, err := eb.nextUnstartedTransactionWithNonce(fromAddress)
		if err != nil {
			return errors.Wrap(err, "processUnstartedEthTxs failed")
		}
		if etx == nil {
			return nil
		}
		n++
		a, err := newAttempt(context.TODO(), eb.ethClient, eb.keystore, eb.config, *etx, nil)
		if err != nil {
			return errors.Wrap(err, "processUnstartedEthTxs failed")
		}

		if err := eb.saveInProgressTransaction(etx, &a); err != nil {
			return errors.Wrap(err, "processUnstartedEthTxs failed")
		}

		if err := eb.handleInProgressEthTx(*etx, a, time.Now()); err != nil {
			return errors.Wrap(err, "processUnstartedEthTxs failed")
		}
	}
}

// handleInProgressEthTx checks if there is any transaction
// in_progress and if so, finishes the job
func (eb *EthBroadcaster) handleAnyInProgressEthTx(fromAddress gethCommon.Address) error {
	etx, err := getInProgressEthTx(eb.db, fromAddress)
	if err != nil {
		return errors.Wrap(err, "handleAnyInProgressEthTx failed")
	}
	if etx != nil {
		if err := eb.handleInProgressEthTx(*etx, etx.EthTxAttempts[0], etx.CreatedAt); err != nil {
			return errors.Wrap(err, "handleAnyInProgressEthTx failed")
		}
	}
	return nil
}

// getInProgressEthTx returns either 0 or 1 transaction that was left in
// an unfinished state because something went screwy the last time. Most likely
// the node crashed in the middle of the ProcessUnstartedEthTxs loop.
// It may or may not have been broadcast to an eth node.
func getInProgressEthTx(db *gorm.DB, fromAddress gethCommon.Address) (*models.EthTx, error) {
	etx := &models.EthTx{}
	err := db.Preload("EthTxAttempts").First(etx, "from_address = ? AND state = 'in_progress'", fromAddress.Bytes()).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if len(etx.EthTxAttempts) != 1 || etx.EthTxAttempts[0].State != models.EthTxAttemptInProgress {
		return nil, errors.Errorf("invariant violation: expected in_progress transaction %v to have exactly one unsent attempt. "+
			"Your database is in an inconsistent state and this node will not function correctly until the problem is resolved", etx.ID)
	}
	return etx, errors.Wrap(err, "getInProgressEthTx failed")
}

// There can be at most one in_progress transaction per address.
// Here we complete the job that we didn't finish last time.
func (eb *EthBroadcaster) handleInProgressEthTx(etx models.EthTx, attempt models.EthTxAttempt, initialBroadcastAt time.Time) error {
	if etx.State != models.EthTxInProgress {
		return errors.Errorf("invariant violation: expected transaction %v to be in_progress, it was %s", etx.ID, etx.State)
	}

	sendError := sendTransaction(context.TODO(), eb.ethClient, attempt, etx)

	if sendError.IsTooExpensive() {
		logger.Errorw("EthBroadcaster: transaction gas price was rejected by the eth node for being too high. Consider increasing your eth node's RPCTxFeeCap (it is suggested to run geth with no cap i.e. --rpc.gascap=0 --rpc.txfeecap=0)",
			"ethTxID", etx.ID,
			"err", sendError,
			"gasPrice", attempt.GasPrice,
			"gasLimit", etx.GasLimit,
			"id", "RPCTxFeeCapExceeded",
		)
		etx.Error = sendError.StrPtr()
		// Attempt is thrown away in this case; we don't need it since it never got accepted by a node
		return saveFatallyErroredTransaction(eb.db, &etx)
	}

	if sendError.Fatal() {
		logger.Errorw("EthBroadcaster: fatal error sending transaction", "ethTxID", etx.ID, "error", sendError, "gasLimit", etx.GasLimit, "gasPrice", attempt.GasPrice)
		etx.Error = sendError.StrPtr()
		// Attempt is thrown away in this case; we don't need it since it never got accepted by a node
		return saveFatallyErroredTransaction(eb.db, &etx)
	}

	etx.BroadcastAt = &initialBroadcastAt

	if sendError.IsNonceTooLowError() || sendError.IsReplacementUnderpriced() {
		// There are three scenarios that this can happen:
		//
		// SCENARIO 1
		//
		// This is resuming a previous crashed run. In this scenario, it is
		// likely that our previous transaction was the one who was confirmed,
		// in which case we hand it off to the eth confirmer to get the
		// receipt.
		//
		// SCENARIO 2
		//
		// It is also possible that an external wallet can have messed with the
		// account and sent a transaction on this nonce.
		//
		// In this case, the onus is on the node operator since this is
		// explicitly unsupported.
		//
		// If it turns out to have been an external wallet, we will never get a
		// receipt for this transaction and it will eventually be marked as
		// errored.
		//
		// The end result is that we will NOT SEND a transaction for this
		// nonce.
		//
		// SCENARIO 3
		//
		// The network/eth client can be assumed to have at-least-once delivery
		// behaviour. It is possible that the eth client could have already
		// sent this exact same transaction even if this is our first time
		// calling SendTransaction().
		//
		// In all scenarios, the correct thing to do is assume success for now
		// and hand off to the eth confirmer to get the receipt (or mark as
		// failed).
		sendError = nil
	}

	if sendError.IsTerminallyUnderpriced() {
		return eb.tryAgainWithHigherGasPrice(sendError, etx, attempt, initialBroadcastAt)
	}

	if sendError.IsTemporarilyUnderpriced() {
		// If we can't even get the transaction into the mempool at all, assume
		// success (even though the transaction will never confirm) and hand
		// off to the ethConfirmer to bump gas periodically until we _can_ get
		// it in
		logger.Infow("EthBroadcaster: Transaction temporarily underpriced", "ethTxID", etx.ID, "err", sendError.Error(), "gasPriceWei", attempt.GasPrice.String())
		sendError = nil
	}

	if sendError.IsInsufficientEth() {
		logger.Errorw(fmt.Sprintf("EthBroadcaster: EthTxAttempt %v (hash 0x%x) at gas price (%s Wei) was rejected due to insufficient eth. "+
			"The eth node returned %s. "+
			"ACTION REQUIRED: Chainlink wallet with address 0x%x is OUT OF FUNDS",
			attempt.ID, attempt.Hash, attempt.GasPrice.String(), sendError.Error(), etx.FromAddress,
		), "ethTxID", etx.ID, "err", sendError)
		return saveAttempt(eb.db, &etx, attempt, models.EthTxAttemptInsufficientEth)
	}

	if sendError == nil {
		return saveAttempt(eb.db, &etx, attempt, models.EthTxAttemptBroadcast)
	}

	// Any other type of error is considered temporary or resolvable by the
	// node operator, but will likely prevent other transactions from working.
	// Safest thing to do is bail out and wait for the next poll.
	return errors.Wrapf(sendError, "error while sending transaction %v", etx.ID)
}

// Finds next transaction in the queue, assigns a nonce, and moves it to "in_progress" state ready for broadcast.
// Returns nil if no transactions are in queue
func (eb *EthBroadcaster) nextUnstartedTransactionWithNonce(fromAddress gethCommon.Address) (*models.EthTx, error) {
	etx := &models.EthTx{}
	if err := findNextUnstartedTransactionFromAddress(eb.db, etx, fromAddress); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Finish. No more transactions left to process. Hoorah!
			return nil, nil
		}
		return nil, errors.Wrap(err, "findNextUnstartedTransactionFromAddress failed")
	}

	nonce, err := GetNextNonce(eb.db, etx.FromAddress)
	if err != nil {
		return nil, err
	}
	etx.Nonce = &nonce
	return etx, nil
}

func (eb *EthBroadcaster) saveInProgressTransaction(etx *models.EthTx, attempt *models.EthTxAttempt) error {
	if etx.State != models.EthTxUnstarted {
		return errors.Errorf("can only transition to in_progress from unstarted, transaction is currently %s", etx.State)
	}
	if attempt.State != models.EthTxAttemptInProgress {
		return errors.New("attempt state must be in_progress")
	}
	etx.State = models.EthTxInProgress
	return postgres.GormTransactionWithDefaultContext(eb.db, func(tx *gorm.DB) error {
		if err := tx.Create(attempt).Error; err != nil {
			return errors.Wrap(err, "saveInProgressTransaction failed to create eth_tx_attempt")
		}
		return errors.Wrap(tx.Save(etx).Error, "saveInProgressTransaction failed to save eth_tx")
	})
}

// Finds earliest saved transaction that has yet to be broadcast from the given address
func findNextUnstartedTransactionFromAddress(db *gorm.DB, etx *models.EthTx, fromAddress gethCommon.Address) error {
	return db.
		Where("from_address = ? AND state = 'unstarted'", fromAddress).
		Order("value ASC, created_at ASC, id ASC").
		First(etx).
		Error
}

func saveAttempt(db *gorm.DB, etx *models.EthTx, attempt models.EthTxAttempt, newAttemptState models.EthTxAttemptState, callbacks ...func(tx *gorm.DB) error) error {
	if etx.State != models.EthTxInProgress {
		return errors.Errorf("can only transition to unconfirmed from in_progress, transaction is currently %s", etx.State)
	}
	if attempt.State != models.EthTxAttemptInProgress {
		return errors.New("attempt must be in in_progress state")
	}
	if !(newAttemptState == models.EthTxAttemptBroadcast || newAttemptState == models.EthTxAttemptInsufficientEth) {
		return errors.Errorf("new attempt state must be broadcast or insufficient_eth, got: %s", newAttemptState)
	}
	etx.State = models.EthTxUnconfirmed
	attempt.State = newAttemptState
	return postgres.GormTransactionWithDefaultContext(db, func(tx *gorm.DB) error {
		if err := IncrementNextNonce(tx, etx.FromAddress, *etx.Nonce); err != nil {
			return errors.Wrap(err, "saveUnconfirmed failed")
		}
		if err := tx.Save(etx).Error; err != nil {
			return errors.Wrap(err, "saveUnconfirmed failed to save eth_tx")
		}
		if err := tx.Save(&attempt).Error; err != nil {
			return errors.Wrap(err, "saveUnconfirmed failed to save eth_tx_attempt")
		}
		for _, f := range callbacks {
			if err := f(tx); err != nil {
				return errors.Wrap(err, "saveUnconfirmed failed")
			}
		}
		return nil
	})
}

func (eb *EthBroadcaster) tryAgainWithHigherGasPrice(sendError *eth.SendError, etx models.EthTx, attempt models.EthTxAttempt, initialBroadcastAt time.Time) error {
	bumpedGasPrice, err := BumpGas(eb.config, attempt.GasPrice.ToInt())
	if err != nil {
		return errors.Wrap(err, "tryAgainWithHigherGasPrice failed")
	}
	logger.Errorw(fmt.Sprintf("default gas price %v wei was rejected by the eth node for being too low. "+
		"Eth node returned: '%s'. "+
		"Bumping to %v wei and retrying. ACTION REQUIRED: This is a configuration error. "+
		"Consider increasing ETH_GAS_PRICE_DEFAULT", eb.config.EthGasPriceDefault(), sendError.Error(), bumpedGasPrice), "err", err)
	if bumpedGasPrice.Cmp(attempt.GasPrice.ToInt()) == 0 && bumpedGasPrice.Cmp(eb.config.EthMaxGasPriceWei()) == 0 {
		return errors.Errorf("Hit gas price bump ceiling, will not bump further. This is a terminal error")
	}
	ctx, cancel := eth.DefaultQueryCtx()
	defer cancel()

	replacementAttempt, err := newAttempt(ctx, eb.ethClient, eb.keystore, eb.config, etx, bumpedGasPrice)
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "tryAgainWithHigherGasPrice failed, context deadline exceeded")
	} else if err != nil {
		return errors.Wrap(err, "tryAgainWithHigherGasPrice failed")
	}

	if err := saveReplacementInProgressAttempt(eb.db, attempt, &replacementAttempt); err != nil {
		return errors.Wrap(err, "tryAgainWithHigherGasPrice failed")
	}
	return eb.handleInProgressEthTx(etx, replacementAttempt, initialBroadcastAt)
}

func saveFatallyErroredTransaction(db *gorm.DB, etx *models.EthTx) error {
	if etx.State != models.EthTxInProgress {
		return errors.Errorf("can only transition to fatal_error from in_progress, transaction is currently %s", etx.State)
	}
	if etx.Error == nil {
		return errors.New("expected error field to be set")
	}
	etx.Nonce = nil
	etx.State = models.EthTxFatalError
	return postgres.GormTransactionWithDefaultContext(db, func(tx *gorm.DB) error {
		if err := tx.Exec(`DELETE FROM eth_tx_attempts WHERE eth_tx_id = ?`, etx.ID).Error; err != nil {
			return errors.Wrapf(err, "saveFatallyErroredTransaction failed to delete eth_tx_attempt with eth_tx.ID %v", etx.ID)
		}
		return errors.Wrap(tx.Save(etx).Error, "saveFatallyErroredTransaction failed to save eth_tx")
	})
}

// GetNextNonce returns keys.next_nonce for the given address
func GetNextNonce(db *gorm.DB, address gethCommon.Address) (int64, error) {
	var nonce int64
	row := db.Raw("SELECT next_nonce FROM keys WHERE address = ?", address).Row()
	if err := row.Scan(&nonce); err != nil {
		return 0, errors.Wrap(err, "GetNextNonce failed scanning row")
	}
	return nonce, nil
}

// IncrementNextNonce increments keys.next_nonce by 1
func IncrementNextNonce(db *gorm.DB, address gethCommon.Address, currentNonce int64) error {
	res := db.Exec("UPDATE keys SET next_nonce = next_nonce + 1, updated_at = NOW() WHERE address = ? AND next_nonce = ?", address.Bytes(), currentNonce)
	if res.Error != nil {
		return errors.Wrap(res.Error, "IncrementNextNonce failed to update keys")
	}
	if res.RowsAffected == 0 {
		var key ethkey.Key
		db.Where("address = ?", address.Bytes()).First(&key)
		return errors.New("invariant violation: could not increment nonce because no rows matched query. " +
			"Either the key is missing or the nonce has been modified by an external process. This is an unrecoverable error")
	}
	return nil
}
