package bulletprooftxmanager

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EthBroadcaster monitors eth_txes for transactions that need to
// be broadcast, assigns nonces and ensures that at least one eth node
// somewhere has received the transaction successfully.
//
// This does not guarantee delivery! A whole host of other things can
// subsequently go wrong such as transactions being evicted from the mempool,
// eth nodes going offline etc. Responsibility for ensuring eventual inclusion
// into the chain falls on the shoulders of the ethConfirmer.
//
// What ethBroadcaster does guarantee is:
// - a monotonic series of increasing nonces for eth_txes that can all eventually be confirmed if you retry enough times
// - transition of eth_txes out of unstarted into either fatal_error or unconfirmed
// - existence of a saved eth_tx_attempt
type EthBroadcaster interface {
	Start() error
	Close() error

	Trigger()

	ProcessUnstartedEthTxs(models.Key) error
}

type ethBroadcaster struct {
	store     *store.Store
	ethClient eth.Client
	config    orm.ConfigReader

	ethTxInsertListener postgres.Subscription
	eventBroadcaster    postgres.EventBroadcaster

	// trigger allows other goroutines to force ethBroadcaster to rescan the
	// database early (before the next poll interval)
	trigger chan struct{}

	ctx       context.Context
	ctxCancel context.CancelFunc
	wg        sync.WaitGroup

	utils.StartStopOnce
}

// NewEthBroadcaster returns a new concrete ethBroadcaster
func NewEthBroadcaster(store *store.Store, config orm.ConfigReader, eventBroadcaster postgres.EventBroadcaster) EthBroadcaster {
	ctx, cancel := context.WithCancel(context.Background())
	return &ethBroadcaster{
		store:            store,
		config:           config,
		ethClient:        store.EthClient,
		trigger:          make(chan struct{}, 1),
		ctx:              ctx,
		ctxCancel:        cancel,
		wg:               sync.WaitGroup{},
		eventBroadcaster: eventBroadcaster,
	}
}

func (eb *ethBroadcaster) Start() error {
	if !eb.OkayToStart() {
		return errors.New("EthBroadcaster is already started")
	}

	var err error
	eb.ethTxInsertListener, err = eb.eventBroadcaster.Subscribe(postgres.ChannelInsertOnEthTx, "")
	if err != nil {
		return errors.Wrap(err, "EthBroadcaster could not start")
	}

	syncer := NewNonceSyncer(eb.store, eb.config, eb.ethClient)
	if err := syncer.SyncAll(eb.ctx); err != nil {
		return errors.Wrap(err, "EthBroadcaster failed to sync with on-chain nonce")
	}

	eb.wg.Add(1)
	go eb.monitorEthTxs()

	eb.wg.Add(1)
	go eb.ethTxInsertTriggerer()

	return nil
}

func (eb *ethBroadcaster) Close() error {
	if !eb.OkayToStop() {
		return errors.New("EthBroadcaster is already stopped")
	}

	if eb.ethTxInsertListener != nil {
		eb.ethTxInsertListener.Close()
	}

	eb.ctxCancel()
	eb.wg.Wait()

	return nil
}

func (eb *ethBroadcaster) Trigger() {
	select {
	case eb.trigger <- struct{}{}:
	default:
	}
}

func (eb *ethBroadcaster) ethTxInsertTriggerer() {
	defer eb.wg.Done()
	for {
		select {
		case <-eb.ethTxInsertListener.Events():
			eb.Trigger()
		case <-eb.ctx.Done():
			return
		}
	}
}

func (eb *ethBroadcaster) monitorEthTxs() {
	defer eb.wg.Done()
	for {
		pollDBTimer := time.NewTimer(utils.WithJitter(eb.config.TriggerFallbackDBPollInterval()))

		keys, err := eb.store.SendKeys()

		if err != nil {
			logger.Error(errors.Wrap(err, "monitorEthTxs failed getting key"))
		} else {
			var wg sync.WaitGroup

			// It is safe to process separate keys concurrently
			// NOTE: This design will block one key if another takes a really long time to execute
			wg.Add(len(keys))
			for _, key := range keys {
				go func(k models.Key) {
					if err := eb.ProcessUnstartedEthTxs(k); err != nil {
						logger.Errorw("Error in ProcessUnstartedEthTxs", "error", err)
					}

					wg.Done()
				}(key)
			}

			wg.Wait()
		}

		select {
		case <-eb.ctx.Done():
			// NOTE: See: https://godoc.org/time#Timer.Stop for an explanation of this pattern
			if !pollDBTimer.Stop() {
				<-pollDBTimer.C
			}
			return
		case <-eb.trigger:
			if !pollDBTimer.Stop() {
				<-pollDBTimer.C
			}
			continue
		case <-pollDBTimer.C:
			continue
		}
	}
}

func (eb *ethBroadcaster) ProcessUnstartedEthTxs(key models.Key) error {
	return eb.store.AdvisoryLocker.WithAdvisoryLock(context.TODO(), postgres.AdvisoryLockClassID_EthBroadcaster, key.ID, func() error {
		return eb.processUnstartedEthTxs(key.Address.Address())
	})
}

// NOTE: This MUST NOT be run concurrently for the same address or it could
// result in undefined state or deadlocks.
// First handle any in_progress transactions left over from last time.
// Then keep looking up unstarted transactions and processing them until there are none remaining.
func (eb *ethBroadcaster) processUnstartedEthTxs(fromAddress gethCommon.Address) error {
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
		etx, err := eb.nextUnstartedTransactionWithNonce(fromAddress)
		if err != nil {
			return errors.Wrap(err, "processUnstartedEthTxs failed")
		}
		if etx == nil {
			return nil
		}
		n++
		a, err := newAttempt(eb.store, *etx, eb.config.EthGasPriceDefault())
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
func (eb *ethBroadcaster) handleAnyInProgressEthTx(fromAddress gethCommon.Address) error {
	etx, err := getInProgressEthTx(eb.store, fromAddress)
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
func getInProgressEthTx(store *store.Store, fromAddress gethCommon.Address) (*models.EthTx, error) {
	etx := &models.EthTx{}
	err := store.DB.Preload("EthTxAttempts").First(etx, "from_address = ? AND state = 'in_progress'", fromAddress.Bytes()).Error
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
func (eb *ethBroadcaster) handleInProgressEthTx(etx models.EthTx, attempt models.EthTxAttempt, initialBroadcastAt time.Time) error {
	if etx.State != models.EthTxInProgress {
		return errors.Errorf("invariant violation: expected transaction %v to be in_progress, it was %s", etx.ID, etx.State)
	}

	ctx, cancel := context.WithTimeout(eb.ctx, maxEthNodeRequestTime)
	defer cancel()
	sendError := sendTransaction(ctx, eb.ethClient, attempt)

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
		return saveFatallyErroredTransaction(eb.store, &etx)
	}

	if sendError.Fatal() {
		logger.Errorw("EthBroadcaster: fatal error sending transaction", "ethTxID", etx.ID, "error", sendError, "gasLimit", etx.GasLimit, "gasPrice", attempt.GasPrice)
		etx.Error = sendError.StrPtr()
		// Attempt is thrown away in this case; we don't need it since it never got accepted by a node
		return saveFatallyErroredTransaction(eb.store, &etx)
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
		return saveAttempt(eb.store, &etx, attempt, models.EthTxAttemptInsufficientEth)
	}

	if sendError == nil {
		return saveAttempt(eb.store, &etx, attempt, models.EthTxAttemptBroadcast)
	}

	// Any other type of error is considered temporary or resolvable by the
	// node operator, but will likely prevent other transactions from working.
	// Safest thing to do is bail out and wait for the next poll.
	return errors.Wrapf(sendError, "error while sending transaction %v", etx.ID)
}

// Finds next transaction in the queue, assigns a nonce, and moves it to "in_progress" state ready for broadcast.
// Returns nil if no transactions are in queue
func (eb *ethBroadcaster) nextUnstartedTransactionWithNonce(fromAddress gethCommon.Address) (*models.EthTx, error) {
	etx := &models.EthTx{}
	if err := findNextUnstartedTransactionFromAddress(eb.store.DB, etx, fromAddress); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Finish. No more transactions left to process. Hoorah!
			return nil, nil
		}
		return nil, errors.Wrap(err, "findNextUnstartedTransactionFromAddress failed")
	}

	nonce, err := GetNextNonce(eb.store.DB, etx.FromAddress)
	if err != nil {
		return nil, err
	}
	etx.Nonce = &nonce
	return etx, nil
}

func (eb *ethBroadcaster) saveInProgressTransaction(etx *models.EthTx, attempt *models.EthTxAttempt) error {
	if etx.State != models.EthTxUnstarted {
		return errors.Errorf("can only transition to in_progress from unstarted, transaction is currently %s", etx.State)
	}
	if attempt.State != models.EthTxAttemptInProgress {
		return errors.New("attempt state must be in_progress")
	}
	etx.State = models.EthTxInProgress
	return eb.store.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit(clause.Associations).Create(attempt).Error; err != nil {
			return errors.Wrap(err, "saveInProgressTransaction failed to create eth_tx_attempt")
		}
		return errors.Wrap(tx.Omit(clause.Associations).Save(etx).Error, "saveInProgressTransaction failed to save eth_tx")
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

func saveAttempt(store *store.Store, etx *models.EthTx, attempt models.EthTxAttempt, newAttemptState models.EthTxAttemptState, callbacks ...func(tx *gorm.DB) error) error {
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
	return store.Transaction(func(tx *gorm.DB) error {
		if err := IncrementNextNonce(tx, etx.FromAddress, *etx.Nonce); err != nil {
			return errors.Wrap(err, "saveUnconfirmed failed")
		}
		if err := tx.Omit(clause.Associations).Save(etx).Error; err != nil {
			return errors.Wrap(err, "saveUnconfirmed failed to save eth_tx")
		}
		if err := tx.Omit(clause.Associations).Save(&attempt).Error; err != nil {
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

func (eb *ethBroadcaster) tryAgainWithHigherGasPrice(sendError *eth.SendError, etx models.EthTx, attempt models.EthTxAttempt, initialBroadcastAt time.Time) error {
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
	replacementAttempt, err := newAttempt(eb.store, etx, bumpedGasPrice)
	if err != nil {
		return errors.Wrap(err, "tryAgainWithHigherGasPrice failed")
	}

	if err := saveReplacementInProgressAttempt(eb.store, attempt, &replacementAttempt); err != nil {
		return errors.Wrap(err, "tryAgainWithHigherGasPrice failed")
	}
	return eb.handleInProgressEthTx(etx, replacementAttempt, initialBroadcastAt)
}

func saveFatallyErroredTransaction(store *store.Store, etx *models.EthTx) error {
	if etx.State != models.EthTxInProgress {
		return errors.Errorf("can only transition to fatal_error from in_progress, transaction is currently %s", etx.State)
	}
	if etx.Error == nil {
		return errors.New("expected error field to be set")
	}
	etx.Nonce = nil
	etx.State = models.EthTxFatalError
	return store.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`DELETE FROM eth_tx_attempts WHERE eth_tx_id = ?`, etx.ID).Error; err != nil {
			return errors.Wrapf(err, "saveFatallyErroredTransaction failed to delete eth_tx_attempt with eth_tx.ID %v", etx.ID)
		}
		return errors.Wrap(tx.Omit(clause.Associations).Save(etx).Error, "saveFatallyErroredTransaction failed to save eth_tx")
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
		return errors.New("invariant violation: could not increment nonce because no rows matched query. " +
			"Either the key is missing or the nonce has been modified by an external process. This is an unrecoverable error")
	}
	return nil
}
