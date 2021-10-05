package bulletprooftxmanager

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/jackc/pgconn"
	"github.com/lib/pq"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/gas"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gopkg.in/guregu/null.v4"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// InFlightTransactionRecheckInterval controls how often the EthBroadcaster
// will poll the unconfirmed queue to see if it is allowed to send another
// transaction
const InFlightTransactionRecheckInterval = 1 * time.Second

var errEthTxRemoved = errors.New("eth_tx removed")

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
	logger    logger.Logger
	db        *gorm.DB
	ethClient eth.Client
	ChainKeyStore
	estimator      gas.Estimator
	resumeCallback ResumeCallback

	ethTxInsertListener postgres.Subscription
	eventBroadcaster    postgres.EventBroadcaster

	keyStates []ethkey.State

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
func NewEthBroadcaster(db *gorm.DB, ethClient eth.Client, config Config, keystore KeyStore,
	eventBroadcaster postgres.EventBroadcaster,
	keyStates []ethkey.State, estimator gas.Estimator, resumeCallback ResumeCallback,
	logger logger.Logger) *EthBroadcaster {

	ctx, cancel := context.WithCancel(context.Background())
	triggers := make(map[gethCommon.Address]chan struct{})
	return &EthBroadcaster{
		logger:    logger,
		db:        db,
		ethClient: ethClient,
		ChainKeyStore: ChainKeyStore{
			chainID:  *ethClient.ChainID(),
			config:   config,
			keystore: keystore,
		},
		estimator:        estimator,
		eventBroadcaster: eventBroadcaster,
		keyStates:        keyStates,
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

		if eb.config.EvmNonceAutoSync() {
			syncer := NewNonceSyncer(eb.db, eb.ethClient)
			if err := syncer.SyncAll(eb.ctx, eb.keyStates); err != nil {
				return errors.Wrap(err, "EthBroadcaster failed to sync with on-chain nonce")
			}
		}

		eb.wg.Add(len(eb.keyStates))
		for _, k := range eb.keyStates {
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
			// ignoring trigger for address which is not registered with this EthBroadcaster
			return
		}
		select {
		case triggerCh <- struct{}{}:
		default:
		}
	})

	if !ok {
		eb.logger.Debugf("EthBroadcaster: unstarted; ignoring trigger for %s", addr.Hex())
	}
}

func (eb *EthBroadcaster) ethTxInsertTriggerer() {
	defer eb.wg.Done()
	for {
		select {
		case ev, ok := <-eb.ethTxInsertListener.Events():
			if !ok {
				eb.logger.Debug("EthBroadcaster: ethTxInsertListener channel closed, exiting trigger loop")
				return
			}
			hexAddr := ev.Payload
			address := gethCommon.HexToAddress(hexAddr)
			eb.Trigger(address)
		case <-eb.ctx.Done():
			return
		}
	}
}

func (eb *EthBroadcaster) monitorEthTxs(k ethkey.State, triggerCh chan struct{}) {
	defer eb.wg.Done()
	for {
		pollDBTimer := time.NewTimer(utils.WithJitter(eb.config.TriggerFallbackDBPollInterval()))

		if err := eb.ProcessUnstartedEthTxs(k); err != nil {
			eb.logger.Errorw("Error in ProcessUnstartedEthTxs", "error", err)
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

func (eb *EthBroadcaster) ProcessUnstartedEthTxs(keyState ethkey.State) error {
	return eb.processUnstartedEthTxs(keyState.Address.Address())
}

// NOTE: This MUST NOT be run concurrently for the same address or it could
// result in undefined state or deadlocks.
// First handle any in_progress transactions left over from last time.
// Then keep looking up unstarted transactions and processing them until there are none remaining.
func (eb *EthBroadcaster) processUnstartedEthTxs(fromAddress gethCommon.Address) error {
	var n uint
	mark := time.Now()
	defer func() {
		if n > 0 {
			eb.logger.Debugw("EthBroadcaster: finished processUnstartedEthTxs", "address", fromAddress, "time", time.Since(mark), "n", n, "id", "eth_broadcaster")
		}
	}()

	if err := eb.handleAnyInProgressEthTx(fromAddress); err != nil {
		return errors.Wrap(err, "processUnstartedEthTxs failed")
	}
	for {
		maxInFlightTransactions := eb.config.EvmMaxInFlightTransactions()
		if maxInFlightTransactions > 0 {
			nUnconfirmed, err := CountUnconfirmedTransactions(eb.db, fromAddress, eb.chainID)
			if err != nil {
				return errors.Wrap(err, "CountUnconfirmedTransactions failed")
			}
			if nUnconfirmed >= maxInFlightTransactions {
				nUnstarted, err := CountUnstartedTransactions(eb.db, fromAddress, eb.chainID)
				if err != nil {
					return errors.Wrap(err, "CountUnstartedTransactions failed")
				}
				eb.logger.Warnw(fmt.Sprintf(`EthBroadcaster: transaction throttling; %d transactions in-flight and %d unstarted transactions pending (maximum number of in-flight transactions is %d per key). %s`, nUnconfirmed, nUnstarted, maxInFlightTransactions, static.EvmMaxInFlightTransactionsWarningLabel), "maxInFlightTransactions", maxInFlightTransactions, "nUnconfirmed", nUnconfirmed, "nUnstarted", nUnstarted)
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
		var a EthTxAttempt
		if eb.config.EvmEIP1559DynamicFees() {
			fee, gasLimit, err := eb.estimator.GetDynamicFee(etx.GasLimit)
			if err != nil {
				return errors.Wrap(err, "failed to get dynamic gas fee")
			}
			a, err = eb.NewDynamicFeeAttempt(*etx, fee, gasLimit)
			if err != nil {
				return errors.Wrap(err, "processUnstartedEthTxs failed")
			}
		} else {
			gasPrice, gasLimit, err := eb.estimator.GetLegacyGas(etx.EncodedPayload, etx.GasLimit)
			if err != nil {
				return errors.Wrap(err, "failed to estimate gas")
			}
			a, err = eb.NewLegacyAttempt(*etx, gasPrice, gasLimit)
			if err != nil {
				return errors.Wrap(err, "processUnstartedEthTxs failed")
			}
		}

		if err := eb.saveInProgressTransaction(etx, &a); errors.Is(err, errEthTxRemoved) {
			eb.logger.Debugw("EthBroadcaster: eth_tx removed", "etxID", etx.ID, "subject", etx.Subject)
			continue
		} else if err != nil {
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
func getInProgressEthTx(db *gorm.DB, fromAddress gethCommon.Address) (*EthTx, error) {
	etx := &EthTx{}
	err := db.Preload("EthTxAttempts").First(etx, "from_address = ? AND state = 'in_progress'", fromAddress.Bytes()).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if len(etx.EthTxAttempts) != 1 || etx.EthTxAttempts[0].State != EthTxAttemptInProgress {
		return nil, errors.Errorf("invariant violation: expected in_progress transaction %v to have exactly one unsent attempt. "+
			"Your database is in an inconsistent state and this node will not function correctly until the problem is resolved", etx.ID)
	}
	return etx, errors.Wrap(err, "getInProgressEthTx failed")
}

// SimulationTimeout must be short since simulation adds latency to
// broadcasting a tx which can negatively affect response time
const SimulationTimeout = 2 * time.Second

// There can be at most one in_progress transaction per address.
// Here we complete the job that we didn't finish last time.
func (eb *EthBroadcaster) handleInProgressEthTx(etx EthTx, attempt EthTxAttempt, initialBroadcastAt time.Time) error {
	if etx.State != EthTxInProgress {
		return errors.Errorf("invariant violation: expected transaction %v to be in_progress, it was %s", etx.ID, etx.State)
	}
	parentCtx := context.TODO()

	if etx.Simulate {
		simulationCtx, cancel := context.WithTimeout(parentCtx, SimulationTimeout)
		defer cancel()
		if b, err := simulateTransaction(simulationCtx, eb.ethClient, attempt, etx); err != nil {
			if jErr := eth.ExtractRPCError(err); jErr != nil {
				eb.logger.Errorw("BulletproofTxManager: Transaction reverted during simulation", "ethTxAttemptID", attempt.ID, "txHash", attempt.Hash, "err", err, "rpcErr", jErr.String(), "returnValue", b.String())
				etx.Error = null.StringFrom(fmt.Sprintf("transaction reverted during simulation: %s", jErr.String()))
				return eb.saveFatallyErroredTransaction(&etx)
			}
			logger.Warnw("BulletproofTxManager: Transaction simulation failed, will attempt to send anyway", "ethTxAttemptID", attempt.ID, "txHash", attempt.Hash, "err", err, "returnValue", b.String())
		} else {
			logger.Debugw("BulletproofTxManager: Transaction simulation succeeded", "ethTxAttemptID", attempt.ID, "txHash", attempt.Hash, "returnValue", b.String())
		}
	}

	sendError := sendTransaction(parentCtx, eb.ethClient, attempt, etx, eb.logger)

	if sendError.IsTooExpensive() {
		eb.logger.Errorw("EthBroadcaster: transaction gas price was rejected by the eth node for being too high. Consider increasing your eth node's RPCTxFeeCap (it is suggested to run geth with no cap i.e. --rpc.gascap=0 --rpc.txfeecap=0)",
			"ethTxID", etx.ID,
			"err", sendError,
			"gasPrice", attempt.GasPrice,
			"gasLimit", etx.GasLimit,
			"id", "RPCTxFeeCapExceeded",
		)
		etx.Error = null.StringFrom(sendError.Error())
		// Attempt is thrown away in this case; we don't need it since it never got accepted by a node
		return eb.saveFatallyErroredTransaction(&etx)
	}

	if sendError.Fatal() {
		eb.logger.Errorw("EthBroadcaster: fatal error sending transaction", "ethTxID", etx.ID, "error", sendError, "gasLimit", etx.GasLimit, "gasPrice", attempt.GasPrice)
		etx.Error = null.StringFrom(sendError.Error())
		// Attempt is thrown away in this case; we don't need it since it never got accepted by a node
		return eb.saveFatallyErroredTransaction(&etx)
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
		return eb.tryAgainBumpingGas(sendError, etx, attempt, initialBroadcastAt)
	}

	if sendError.IsFeeTooLow() || sendError.IsFeeTooHigh() {
		return eb.tryAgainWithNewEstimation(sendError, etx, attempt, initialBroadcastAt)
	}

	if sendError.IsTemporarilyUnderpriced() {
		// If we can't even get the transaction into the mempool at all, assume
		// success (even though the transaction will never confirm) and hand
		// off to the ethConfirmer to bump gas periodically until we _can_ get
		// it in
		eb.logger.Infow("EthBroadcaster: Transaction temporarily underpriced", "ethTxID", etx.ID, "err", sendError.Error(), "gasPriceWei", attempt.GasPrice.String())
		sendError = nil
	}

	if sendError.IsInsufficientEth() {
		eb.logger.Errorw(fmt.Sprintf("EthBroadcaster: tx 0x%x with type 0x%d was rejected due to insufficient eth. "+
			"The eth node returned %s. "+
			"ACTION REQUIRED: Chainlink wallet with address 0x%x is OUT OF FUNDS",
			attempt.Hash, attempt.TxType, sendError.Error(), etx.FromAddress,
		), "ethTxID", etx.ID, "err", sendError, "gasPrice", attempt.GasPrice,
			"gasTipCap", attempt.GasTipCap, "gasFeeCap", attempt.GasFeeCap)
		// NOTE: This bails out of the entire cycle and essentially "blocks" on
		// any transaction that gets insufficient_eth. This is OK if a
		// transaction with a large VALUE blocks because this always comes last
		// in the processing list.
		// If it blocks because of a transaction that is expensive due to large
		// gas limit, we could have smaller transactions "above" it that could
		// theoretically be sent, but will instead be blocked.
		return sendError
	}

	if sendError == nil {
		return saveAttempt(eb.db, &etx, attempt, EthTxAttemptBroadcast)
	}

	// Any other type of error is considered temporary or resolvable by the
	// node operator, but will likely prevent other transactions from working.
	// Safest thing to do is bail out and wait for the next poll.
	return errors.Wrapf(sendError, "error while sending transaction %v", etx.ID)
}

// Finds next transaction in the queue, assigns a nonce, and moves it to "in_progress" state ready for broadcast.
// Returns nil if no transactions are in queue
func (eb *EthBroadcaster) nextUnstartedTransactionWithNonce(fromAddress gethCommon.Address) (*EthTx, error) {
	etx := &EthTx{}
	if err := findNextUnstartedTransactionFromAddress(eb.db, etx, fromAddress, eb.chainID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Finish. No more transactions left to process. Hoorah!
			return nil, nil
		}
		return nil, errors.Wrap(err, "findNextUnstartedTransactionFromAddress failed")
	}

	nonce, err := GetNextNonce(eb.db, etx.FromAddress, &eb.chainID)
	if err != nil {
		return nil, err
	}
	etx.Nonce = &nonce
	return etx, nil
}

func (eb *EthBroadcaster) saveInProgressTransaction(etx *EthTx, attempt *EthTxAttempt) error {
	if etx.State != EthTxUnstarted {
		return errors.Errorf("can only transition to in_progress from unstarted, transaction is currently %s", etx.State)
	}
	if attempt.State != EthTxAttemptInProgress {
		return errors.New("attempt state must be in_progress")
	}
	etx.State = EthTxInProgress
	return postgres.GormTransactionWithDefaultContext(eb.db, func(tx *gorm.DB) error {
		err := tx.Create(attempt).Error
		if err != nil {
			switch e := err.(type) {
			case *pq.Error:
				if e.Constraint == "eth_tx_attempts_eth_tx_id_fkey" {
					return errEthTxRemoved
				}
			case *pgconn.PgError:
				if e.ConstraintName == "eth_tx_attempts_eth_tx_id_fkey" {
					return errEthTxRemoved
				}
			}
			return errors.Wrap(err, "saveInProgressTransaction failed to create eth_tx_attempt")
		}
		return errors.Wrap(tx.Save(etx).Error, "saveInProgressTransaction failed to save eth_tx")
	})
}

// Finds earliest saved transaction that has yet to be broadcast from the given address
func findNextUnstartedTransactionFromAddress(db *gorm.DB, etx *EthTx, fromAddress gethCommon.Address, chainID big.Int) error {
	return db.
		Where("from_address = ? AND state = 'unstarted' AND evm_chain_id = ?", fromAddress, chainID.String()).
		Order("value ASC, created_at ASC, id ASC").
		First(etx).
		Error
}

func saveAttempt(db *gorm.DB, etx *EthTx, attempt EthTxAttempt, NewAttemptState EthTxAttemptState, callbacks ...func(tx *gorm.DB) error) error {
	if etx.State != EthTxInProgress {
		return errors.Errorf("can only transition to unconfirmed from in_progress, transaction is currently %s", etx.State)
	}
	if attempt.State != EthTxAttemptInProgress {
		return errors.New("attempt must be in in_progress state")
	}
	if !(NewAttemptState == EthTxAttemptBroadcast) {
		return errors.Errorf("new attempt state must be broadcast, got: %s", NewAttemptState)
	}
	etx.State = EthTxUnconfirmed
	attempt.State = NewAttemptState
	return postgres.GormTransactionWithDefaultContext(db, func(tx *gorm.DB) error {
		if err := IncrementNextNonce(tx, etx.FromAddress, etx.EVMChainID.ToInt(), *etx.Nonce); err != nil {
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

func (eb *EthBroadcaster) tryAgainBumpingGas(sendError *eth.SendError, etx EthTx, attempt EthTxAttempt, initialBroadcastAt time.Time) error {
	if attempt.TxType == 0x2 {
		return errors.New("bumping gas on initial send is not supported for EIP-1559 transactions")
	}
	bumpedGasPrice, bumpedGasLimit, err := eb.estimator.BumpLegacyGas(attempt.GasPrice.ToInt(), etx.GasLimit)
	if err != nil {
		return errors.Wrap(err, "tryAgainWithHigherGasPrice failed")
	}
	eb.logger.Errorw(fmt.Sprintf("default gas price %v wei was rejected by the eth node for being too low. "+
		"Eth node returned: '%s'. "+
		"Bumping to %v wei and retrying. ACTION REQUIRED: This is a configuration error. "+
		"Consider increasing ETH_GAS_PRICE_DEFAULT", eb.config.EvmGasPriceDefault(), sendError.Error(), bumpedGasPrice), "err", err)
	if bumpedGasPrice.Cmp(attempt.GasPrice.ToInt()) == 0 && bumpedGasPrice.Cmp(eb.config.EvmMaxGasPriceWei()) == 0 {
		return errors.Errorf("Hit gas price bump ceiling, will not bump further. This is a terminal error")
	}
	return eb.tryAgainWithNewGas(etx, attempt, initialBroadcastAt, bumpedGasPrice, bumpedGasLimit)
}

func (eb *EthBroadcaster) tryAgainWithNewEstimation(sendError *eth.SendError, etx EthTx, attempt EthTxAttempt, initialBroadcastAt time.Time) error {
	gasPrice, gasLimit, err := eb.estimator.GetLegacyGas(etx.EncodedPayload, etx.GasLimit, gas.OptForceRefetch)
	if err != nil {
		return errors.Wrap(err, "tryAgainWithNewEstimation failed to estimate gas")
	}
	eb.logger.Debugw("Optimism rejected transaction due to incorrect fee, re-estimated and will try again",
		"etxID", etx.ID, "err", err, "newGasPrice", gasPrice, "newGasLimit", gasLimit)
	return eb.tryAgainWithNewGas(etx, attempt, initialBroadcastAt, gasPrice, gasLimit)
}

func (eb *EthBroadcaster) tryAgainWithNewGas(etx EthTx, attempt EthTxAttempt, initialBroadcastAt time.Time, newGasPrice *big.Int, newGasLimit uint64) error {
	replacementAttempt, err := eb.NewLegacyAttempt(etx, newGasPrice, newGasLimit)
	if err != nil {
		return errors.Wrap(err, "tryAgainWithHigherGasPrice failed")
	}

	if err = saveReplacementInProgressAttempt(eb.db, attempt, &replacementAttempt); err != nil {
		return errors.Wrap(err, "tryAgainWithHigherGasPrice failed")
	}
	return eb.handleInProgressEthTx(etx, replacementAttempt, initialBroadcastAt)
}

func (eb *EthBroadcaster) saveFatallyErroredTransaction(etx *EthTx) error {
	if etx.State != EthTxInProgress {
		return errors.Errorf("can only transition to fatal_error from in_progress, transaction is currently %s", etx.State)
	}
	if !etx.Error.Valid {
		return errors.New("expected error field to be set")
	}
	// NOTE: It's simpler to not do this transactionally for now (would require
	// refactoring pipeline runner resume to use postgres events)
	//
	// There is a very tiny possibility of the following:
	//
	// 1. We get a fatal error on the tx, resuming the pipeline with error
	// 2. Crash or failure during persist of fatal errored tx
	// 3. On the subsequent run the tx somehow succeeds and we save it as successful
	//
	// Now we have an errored pipeline even though the tx succeeded. This case
	// is relatively benign and probably nobody will ever run into it in
	// practice, but something to be aware of.
	if etx.PipelineTaskRunID.Valid && eb.resumeCallback != nil {
		err := eb.resumeCallback(etx.PipelineTaskRunID.UUID, nil, errors.Errorf("fatal error while sending transaction: %s", etx.Error.String))
		if errors.Is(err, sql.ErrNoRows) {
			eb.logger.Debugw("callback missing or already resumed", "etxID", etx.ID)
		} else if err != nil {
			return errors.Wrap(err, "failed to resume pipeline")
		}
	}
	etx.Nonce = nil
	etx.State = EthTxFatalError
	return postgres.GormTransactionWithDefaultContext(eb.db, func(tx *gorm.DB) error {
		if err := tx.Exec(`DELETE FROM eth_tx_attempts WHERE eth_tx_id = ?`, etx.ID).Error; err != nil {
			return errors.Wrapf(err, "saveFatallyErroredTransaction failed to delete eth_tx_attempt with eth_tx.ID %v", etx.ID)
		}
		return errors.Wrap(tx.Save(etx).Error, "saveFatallyErroredTransaction failed to save eth_tx")
	})
}

// GetNextNonce returns keys.next_nonce for the given address
func GetNextNonce(db *gorm.DB, address gethCommon.Address, chainID *big.Int) (int64, error) {
	var nonce int64
	row := db.Raw("SELECT next_nonce FROM eth_key_states WHERE address = ? AND evm_chain_id = ?", address, chainID.String()).Row()
	if err := row.Scan(&nonce); err != nil {
		return 0, errors.Wrap(err, "GetNextNonce failed scanning row")
	}
	return nonce, nil
}

// IncrementNextNonce increments keys.next_nonce by 1
func IncrementNextNonce(db *gorm.DB, address gethCommon.Address, chainID *big.Int, currentNonce int64) error {
	res := db.Exec("UPDATE eth_key_states SET next_nonce = next_nonce + 1, updated_at = NOW() WHERE address = ? AND next_nonce = ? AND evm_chain_id = ?", address, currentNonce, chainID.String())
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
