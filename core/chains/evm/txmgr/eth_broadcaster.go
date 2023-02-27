package txmgr

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgconn"
	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/core/chains/evm/label"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const (
	// InFlightTransactionRecheckInterval controls how often the EthBroadcaster
	// will poll the unconfirmed queue to see if it is allowed to send another
	// transaction
	InFlightTransactionRecheckInterval = 1 * time.Second

	// TransmitCheckTimeout controls the maximum amount of time that will be
	// spent on the transmit check.
	TransmitCheckTimeout = 2 * time.Second
)

var (
	promTimeUntilBroadcast = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "tx_manager_time_until_tx_broadcast",
		Help: "The amount of time elapsed from when a transaction is enqueued to until it is broadcast.",
		Buckets: []float64{
			float64(500 * time.Millisecond),
			float64(time.Second),
			float64(5 * time.Second),
			float64(15 * time.Second),
			float64(30 * time.Second),
			float64(time.Minute),
			float64(2 * time.Minute),
		},
	}, []string{"evmChainID"})
)

var errEthTxRemoved = errors.New("eth_tx removed")

// TransmitCheckerFactory creates a transmit checker based on a spec.
type TransmitCheckerFactory interface {
	// BuildChecker builds a new TransmitChecker based on the given spec.
	BuildChecker(spec TransmitCheckerSpec) (TransmitChecker, error)
}

// TransmitChecker determines whether a transaction should be submitted on-chain.
type TransmitChecker interface {

	// Check the given transaction. If the transaction should not be sent, an error indicating why
	// is returned. Errors should only be returned if the checker can confirm that a transaction
	// should not be sent, other errors (for example connection or other unexpected errors) should
	// be logged and swallowed.
	Check(ctx context.Context, l logger.Logger, tx EthTx, a EthTxAttempt) error
}

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
	orm       ORM
	db        *sqlx.DB
	q         pg.Q
	ethClient evmclient.Client
	ChainKeyStore
	estimator      gas.Estimator
	resumeCallback ResumeCallback

	// autoSyncNonce, if set, will cause EthBroadcaster to fast-forward the nonce
	// when Start is called
	autoSyncNonce bool

	ethTxInsertListener pg.Subscription
	eventBroadcaster    pg.EventBroadcaster

	keyStates []ethkey.State

	checkerFactory TransmitCheckerFactory

	// triggers allow other goroutines to force EthBroadcaster to rescan the
	// database early (before the next poll interval)
	// Each key has its own trigger
	triggers map[gethCommon.Address]chan struct{}

	chStop chan struct{}
	wg     sync.WaitGroup

	utils.StartStopOnce
}

// NewEthBroadcaster returns a new concrete EthBroadcaster
func NewEthBroadcaster(db *sqlx.DB, ethClient evmclient.Client, config Config, keystore KeyStore,
	eventBroadcaster pg.EventBroadcaster,
	keyStates []ethkey.State, estimator gas.Estimator, resumeCallback ResumeCallback,
	logger logger.Logger, checkerFactory TransmitCheckerFactory, autoSyncNonce bool) *EthBroadcaster {

	triggers := make(map[gethCommon.Address]chan struct{})
	logger = logger.Named("EthBroadcaster")
	return &EthBroadcaster{
		logger:    logger,
		orm:       NewORM(db, logger, config),
		db:        db,
		q:         pg.NewQ(db, logger, config),
		ethClient: ethClient,
		ChainKeyStore: ChainKeyStore{
			chainID:  *ethClient.ChainID(),
			config:   config,
			keystore: keystore,
		},
		estimator:        estimator,
		resumeCallback:   resumeCallback,
		eventBroadcaster: eventBroadcaster,
		keyStates:        keyStates,
		checkerFactory:   checkerFactory,
		triggers:         triggers,
		chStop:           make(chan struct{}),
		wg:               sync.WaitGroup{},
		autoSyncNonce:    autoSyncNonce,
	}
}

// Start starts EthBroadcaster service.
// The provided context can be used to terminate Start sequence.
func (eb *EthBroadcaster) Start(ctx context.Context) error {
	return eb.StartOnce("EthBroadcaster", func() (err error) {
		eb.ethTxInsertListener, err = eb.eventBroadcaster.Subscribe(pg.ChannelInsertOnEthTx, "")
		if err != nil {
			return errors.Wrap(err, "EthBroadcaster could not start")
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

// Close closes the EthBroadcaster
func (eb *EthBroadcaster) Close() error {
	return eb.StopOnce("EthBroadcaster", func() error {
		if eb.ethTxInsertListener != nil {
			eb.ethTxInsertListener.Close()
		}

		close(eb.chStop)
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
		eb.logger.Debugf("Unstarted; ignoring trigger for %s", addr.Hex())
	}
}

func (eb *EthBroadcaster) ethTxInsertTriggerer() {
	defer eb.wg.Done()
	for {
		select {
		case ev, ok := <-eb.ethTxInsertListener.Events():
			if !ok {
				eb.logger.Debug("ethTxInsertListener channel closed, exiting trigger loop")
				return
			}
			hexAddr := ev.Payload
			address := gethCommon.HexToAddress(hexAddr)
			eb.Trigger(address)
		case <-eb.chStop:
			return
		}
	}
}

func (eb *EthBroadcaster) newNonceSyncBackoff() backoff.Backoff {
	return backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    5 * time.Second,
		Jitter: true,
	}
}

func (eb *EthBroadcaster) newResendBackoff() backoff.Backoff {
	return backoff.Backoff{
		Min:    1 * time.Second,
		Max:    15 * time.Second,
		Jitter: true,
	}
}

func (eb *EthBroadcaster) monitorEthTxs(k ethkey.State, triggerCh chan struct{}) {
	defer eb.wg.Done()

	ctx, cancel := utils.ContextFromChan(eb.chStop)
	defer cancel()

	if eb.autoSyncNonce {
		eb.logger.Debugw("Auto-syncing nonce", "address", k.Address)
		eb.SyncNonce(ctx, k)
		if ctx.Err() != nil {
			return
		}
	} else {
		eb.logger.Debugw("Skipping nonce auto-sync", "address", k.Address)
	}

	// errorRetryCh allows retry on exponential backoff in case of timeout or
	// other unknown error
	var errorRetryCh <-chan time.Time
	bf := eb.newResendBackoff()

	for {
		pollDBTimer := time.NewTimer(utils.WithJitter(eb.config.TriggerFallbackDBPollInterval()))

		err, retryable := eb.ProcessUnstartedEthTxs(ctx, k)
		if err != nil {
			eb.logger.Errorw("Error occurred while handling eth_tx queue in ProcessUnstartedEthTxs", "err", err)
		}
		// On retryable errors we implement exponential backoff retries. This
		// handles intermittent connectivity, remote RPC races, timing issues etc
		if retryable {
			pollDBTimer.Reset(utils.WithJitter(eb.config.TriggerFallbackDBPollInterval()))
			errorRetryCh = time.After(bf.Duration())
		} else {
			bf = eb.newResendBackoff()
			errorRetryCh = nil
		}

		select {
		case <-ctx.Done():
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
		case <-errorRetryCh:
			// Error backoff period reached
			continue
		}
	}
}

// syncNonce tries to sync the key nonce, retrying indefinitely until success
func (eb *EthBroadcaster) SyncNonce(ctx context.Context, k ethkey.State) {
	if k.Disabled {
		eb.logger.Infow("Skipping nonce sync for disabled key", "address", k.Address)
		return
	}
	syncer := NewNonceSyncer(eb.db, eb.logger, eb.config, eb.ethClient, eb.ChainKeyStore.keystore)
	nonceSyncRetryBackoff := eb.newNonceSyncBackoff()
	if err := syncer.Sync(ctx, k); err != nil {
		// Enter retry loop with backoff
		var attempt int
		eb.logger.Errorw("Failed to sync with on-chain nonce", "address", k.Address, "attempt", attempt, "err", err)
		for {
			select {
			case <-eb.chStop:
				return
			case <-time.After(nonceSyncRetryBackoff.Duration()):
				attempt++

				if err := syncer.Sync(ctx, k); err != nil {
					if attempt > 5 {
						eb.logger.Criticalw("Failed to sync with on-chain nonce", "address", k.Address, "attempt", attempt, "err", err)
					} else {
						eb.logger.Warnw("Failed to sync with on-chain nonce", "address", k.Address, "attempt", attempt, "err", err)
					}
					continue
				}
				return
			}
		}
	}
}

// ProcessUnstartedEthTxs picks up and handles all eth_txes in the queue
// revive:disable:error-return
func (eb *EthBroadcaster) ProcessUnstartedEthTxs(ctx context.Context, keyState ethkey.State) (err error, retryable bool) {
	return eb.processUnstartedEthTxs(ctx, keyState.Address.Address())
}

// NOTE: This MUST NOT be run concurrently for the same address or it could
// result in undefined state or deadlocks.
// First handle any in_progress transactions left over from last time.
// Then keep looking up unstarted transactions and processing them until there are none remaining.
func (eb *EthBroadcaster) processUnstartedEthTxs(ctx context.Context, fromAddress gethCommon.Address) (err error, retryable bool) {
	var n uint
	mark := time.Now()
	defer func() {
		if n > 0 {
			eb.logger.Debugw("Finished processUnstartedEthTxs", "address", fromAddress, "time", time.Since(mark), "n", n, "id", "eth_broadcaster")
		}
	}()

	err, retryable = eb.handleAnyInProgressEthTx(ctx, fromAddress)
	if err != nil {
		return errors.Wrap(err, "processUnstartedEthTxs failed on handleAnyInProgressEthTx"), retryable
	}
	for {
		maxInFlightTransactions := eb.config.EvmMaxInFlightTransactions()
		if maxInFlightTransactions > 0 {
			nUnconfirmed, err := CountUnconfirmedTransactions(eb.q, fromAddress, eb.chainID)
			if err != nil {
				return errors.Wrap(err, "CountUnconfirmedTransactions failed"), true
			}
			if nUnconfirmed >= maxInFlightTransactions {
				nUnstarted, err := CountUnstartedTransactions(eb.q, fromAddress, eb.chainID)
				if err != nil {
					return errors.Wrap(err, "CountUnstartedTransactions failed"), true
				}
				eb.logger.Warnw(fmt.Sprintf(`Transaction throttling; %d transactions in-flight and %d unstarted transactions pending (maximum number of in-flight transactions is %d per key). %s`, nUnconfirmed, nUnstarted, maxInFlightTransactions, label.MaxInFlightTransactionsWarning), "maxInFlightTransactions", maxInFlightTransactions, "nUnconfirmed", nUnconfirmed, "nUnstarted", nUnstarted)
				time.Sleep(InFlightTransactionRecheckInterval)
				continue
			}
		}
		etx, err := eb.nextUnstartedTransactionWithNonce(fromAddress)
		if err != nil {
			return errors.Wrap(err, "processUnstartedEthTxs failed on nextUnstartedTransactionWithNonce"), true
		}
		if etx == nil {
			return nil, false
		}
		n++
		var a EthTxAttempt
		keySpecificMaxGasPriceWei := eb.config.KeySpecificMaxGasPriceWei(etx.FromAddress)
		if eb.config.EvmEIP1559DynamicFees() {
			fee, gasLimit, err := eb.estimator.GetDynamicFee(ctx, etx.GasLimit, keySpecificMaxGasPriceWei)
			if err != nil {
				return errors.Wrap(err, "failed to get dynamic gas fee"), true
			}
			a, err = eb.NewDynamicFeeAttempt(*etx, fee, gasLimit)
			if err != nil {
				return errors.Wrap(err, "processUnstartedEthTxs failed on NewDynamicFeeAttempt"), true
			}
		} else {
			gasPrice, gasLimit, err := eb.estimator.GetLegacyGas(ctx, etx.EncodedPayload, etx.GasLimit, keySpecificMaxGasPriceWei)
			if err != nil {
				return errors.Wrap(err, "failed to estimate gas"), true
			}
			a, err = eb.NewLegacyAttempt(*etx, gasPrice, gasLimit)
			if err != nil {
				return errors.Wrap(err, "processUnstartedEthTxs failed on NewLegacyAttempt"), true
			}
		}

		if err := eb.saveInProgressTransaction(etx, &a); errors.Is(err, errEthTxRemoved) {
			eb.logger.Debugw("eth_tx removed", "etxID", etx.ID, "subject", etx.Subject)
			continue
		} else if err != nil {
			return errors.Wrap(err, "processUnstartedEthTxs failed on saveInProgressTransaction"), true
		}

		if err, retryable := eb.handleInProgressEthTx(ctx, *etx, a, time.Now()); err != nil {
			return errors.Wrap(err, "processUnstartedEthTxs failed on handleAnyInProgressEthTx"), retryable
		}
	}
}

// handleInProgressEthTx checks if there is any transaction
// in_progress and if so, finishes the job
func (eb *EthBroadcaster) handleAnyInProgressEthTx(ctx context.Context, fromAddress gethCommon.Address) (err error, retryable bool) {
	etx, err := eb.getInProgressEthTx(fromAddress)
	if err != nil {
		return errors.Wrap(err, "handleAnyInProgressEthTx failed"), true
	}
	if etx != nil {
		if err, retryable := eb.handleInProgressEthTx(ctx, *etx, etx.EthTxAttempts[0], etx.CreatedAt); err != nil {
			return errors.Wrap(err, "handleAnyInProgressEthTx failed"), retryable
		}
	}
	return nil, false
}

// getInProgressEthTx returns either 0 or 1 transaction that was left in
// an unfinished state because something went screwy the last time. Most likely
// the node crashed in the middle of the ProcessUnstartedEthTxs loop.
// It may or may not have been broadcast to an eth node.
func (eb *EthBroadcaster) getInProgressEthTx(fromAddress gethCommon.Address) (etx *EthTx, err error) {
	etx = new(EthTx)
	err = eb.q.Get(etx, `SELECT * FROM eth_txes WHERE from_address = $1 and state = 'in_progress'`, fromAddress.Bytes())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "getInProgressEthTx failed while loading eth tx")
	}
	if err = eb.orm.LoadEthTxAttempts(etx); err != nil {
		return nil, errors.Wrap(err, "getInProgressEthTx failed while loading EthTxAttempts")
	}
	if len(etx.EthTxAttempts) != 1 || etx.EthTxAttempts[0].State != EthTxAttemptInProgress {
		return nil, errors.Errorf("invariant violation: expected in_progress transaction %v to have exactly one unsent attempt. "+
			"Your database is in an inconsistent state and this node will not function correctly until the problem is resolved", etx.ID)
	}
	return etx, errors.Wrap(err, "getInProgressEthTx failed")
}

// There can be at most one in_progress transaction per address.
// Here we complete the job that we didn't finish last time.
func (eb *EthBroadcaster) handleInProgressEthTx(ctx context.Context, etx EthTx, attempt EthTxAttempt, initialBroadcastAt time.Time) (error, bool) {
	if etx.State != EthTxInProgress {
		return errors.Errorf("invariant violation: expected transaction %v to be in_progress, it was %s", etx.ID, etx.State), false
	}

	checkerSpec, err := etx.GetChecker()
	if err != nil {
		return errors.Wrap(err, "parsing transmit checker"), false
	}

	checker, err := eb.checkerFactory.BuildChecker(checkerSpec)
	if err != nil {
		return errors.Wrap(err, "building transmit checker"), false
	}

	lgr := etx.GetLogger(eb.logger.With(
		"gasPrice", attempt.GasPrice,
		"gasTipCap", attempt.GasTipCap,
		"gasFeeCap", attempt.GasFeeCap,
	))

	// If the transmit check does not complete within the timeout, the transaction will be sent
	// anyway.
	checkCtx, cancel := context.WithTimeout(ctx, TransmitCheckTimeout)
	defer cancel()
	err = checker.Check(checkCtx, lgr, etx, attempt)
	if errors.Is(err, context.Canceled) {
		lgr.Warn("Transmission checker timed out, sending anyway")
	} else if err != nil {
		etx.Error = null.StringFrom(err.Error())
		lgr.Warnw("Transmission checker failed, fatally erroring transaction.", "err", err)
		return eb.saveFatallyErroredTransaction(lgr, &etx), true
	}
	cancel()

	sendError := sendTransaction(ctx, eb.ethClient, attempt, etx, lgr)

	if sendError.Fatal() {
		lgr.Criticalw("Fatal error sending transaction", "err", sendError, "etx", etx)
		etx.Error = null.StringFrom(sendError.Error())
		// Attempt is thrown away in this case; we don't need it since it never got accepted by a node
		return eb.saveFatallyErroredTransaction(lgr, &etx), true
	}

	etx.InitialBroadcastAt = &initialBroadcastAt
	etx.BroadcastAt = &initialBroadcastAt

	if sendError.IsNonceTooLowError() || sendError.IsTransactionAlreadyMined() || sendError.IsReplacementUnderpriced() {
		// There are four scenarios that this can happen:
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
		// behavior. It is possible that the eth client could have already
		// sent this exact same transaction even if this is our first time
		// calling SendTransaction().
		//
		// SCENARIO 4 (most likely)
		//
		// A sendonly node got the transaction in first.
		//
		// In all scenarios, the correct thing to do is assume success for now
		// and hand off to the eth confirmer to get the receipt (or mark as
		// failed).
		sendError = nil
	}

	if sendError.IsTerminallyUnderpriced() {
		return eb.tryAgainBumpingGas(ctx, lgr, sendError, etx, attempt, initialBroadcastAt)
	}

	// L2-specific cases
	if sendError.L2FeeTooLow() || sendError.IsL2FeeTooHigh() || sendError.IsL2Full() {
		if eb.ChainKeyStore.config.ChainType().IsL2() {
			return eb.tryAgainWithNewEstimation(ctx, lgr, sendError, etx, attempt, initialBroadcastAt)
		}
		return errors.Wrap(sendError, "this error type only handled for L2s"), false
	}

	if sendError.IsNonceTooHighError() {
		// Nethermind specific error. Nethermind throws a NonceGap error when the tx nonce is
		// greater than current_nonce + tx_count_in_mempool, instead of keeping the tx in mempool.
		// This can happen if previous transactions haven't reached the client yet.
		// The correct thing to do is assume success for now and let the eth_confirmer retry until
		// the nonce gap gets filled by the previous transactions.
		lgr.Warnw("Transaction has a nonce gap.", "err", sendError.Error())
		return sendError, true
	}

	if sendError.IsTemporarilyUnderpriced() {
		// If we can't even get the transaction into the mempool at all, assume
		// success (even though the transaction will never confirm) and hand
		// off to the ethConfirmer to bump gas periodically until we _can_ get
		// it in
		lgr.Infow("Transaction temporarily underpriced", "err", sendError.Error())
		sendError = nil
	}

	if sendError.IsInsufficientEth() {
		lgr.Criticalw(fmt.Sprintf("Tx 0x%x with type 0x%d was rejected due to insufficient eth: %s\n"+
			"ACTION REQUIRED: Chainlink wallet with address 0x%x is OUT OF FUNDS",
			attempt.Hash, attempt.TxType, sendError.Error(), etx.FromAddress,
		), "err", sendError)
		// NOTE: This bails out of the entire cycle and essentially "blocks" on
		// any transaction that gets insufficient_eth. This is OK if a
		// transaction with a large VALUE blocks because this always comes last
		// in the processing list.
		// If it blocks because of a transaction that is expensive due to large
		// gas limit, we could have smaller transactions "above" it that could
		// theoretically be sent, but will instead be blocked.
		return sendError, true
	}

	if sendError == nil {
		// We want to observe the time until the first _successful_ broadcast.
		// Since we can re-enter this method by way of tryAgainBumpingGas,
		// and we pass the same initialBroadcastAt timestamp there, when we re-enter
		// this function we'll be using the same initialBroadcastAt.
		observeTimeUntilBroadcast(eb.chainID, etx.CreatedAt, time.Now())
		return eb.saveAttempt(&etx, attempt, EthTxAttemptBroadcast), true
	}

	// In the case of timeout, we fall back to the backoff retry loop and
	// subsequent tries ought to resend the exact same in-progress transaction
	// attempt and get a definitive answer on what happened
	if sendError.IsTimeout() {
		return errors.Wrapf(sendError, "timeout while sending transaction %s (eth_tx ID %d)", attempt.Hash.Hex(), etx.ID), true
	}

	// Unknown error here. All bets are off in this case, it is possible the
	// transaction could have been accepted. We may be running on an
	// unsupported RPC or chain.
	//
	// The most conservative course of action would be to retry this
	// transaction forever (or until success) however this can lead to nodes
	// getting stuck if we are on an unsupported new chain.
	//
	// We can continue in a kind of gracefully degraded manner if we check the
	// chain for its view on our latest nonce. If it has been incremented, then
	// it accepted the transaction despite the error and we can move forwards
	// assuming success in this case.

	if sendError.IsTxFeeExceedsCap() {
		lgr.Criticalw(fmt.Sprintf("Sending transaction failed; %s", label.RPCTxFeeCapConfiguredIncorrectlyWarning),
			"etx", etx,
			"attempt", attempt,
			"err", sendError,
			"id", "RPCTxFeeCapExceeded",
		)
		// Note that we may have broadcast to multiple nodes and had it
		// accepted by one of them! It is not guaranteed that all nodes share
		// the same tx fee cap. That is why we must treat this as an unknown
		// error that may have been confirmed.
		//
		// If there is only one RPC node, or all RPC nodes have the same
		// configured cap, this transaction will get stuck and keep repeating
		// forever until the issue is resolved.
	} else {
		lgr.Criticalw("Unknown error occurred while handling eth_tx queue in ProcessUnstartedEthTxs. This chain/RPC client may not be supported. Urgent resolution required, Chainlink is currently operating in a degraded state and may miss transactions", "err", sendError, "etx", etx, "attempt", attempt)
	}

	nextNonce, err := eb.ethClient.PendingNonceAt(ctx, etx.FromAddress)
	if err != nil {
		err = multierr.Combine(err, sendError)
		return errors.Wrapf(err, "failed to fetch latest pending nonce after encountering unknown RPC error while sending transaction"), true
	}

	if nextNonce > math.MaxInt64 {
		return errors.Errorf("nonce overflow, got: %v", nextNonce), true
	}

	if int64(nextNonce) > *etx.Nonce {
		// Despite the error, the RPC node considers the previously sent
		// transaction to have been accepted. In this case, the right thing to
		// do is assume success and hand off to EthConfirmer
		return eb.saveAttempt(&etx, attempt, EthTxAttemptBroadcast), true
	}

	// Either the unknown error prevented the transaction from being mined, or
	// it has not yet propagated to the mempool, or there is some race on the
	// remote RPC.
	//
	// In all cases, the best thing we can do is go into a retry loop and keep
	// trying to send the transaction over again.
	return errors.Wrapf(sendError, "retryable error while sending transaction %s (eth_tx ID %d)", attempt.Hash.Hex(), etx.ID), true
}

// Finds next transaction in the queue, assigns a nonce, and moves it to "in_progress" state ready for broadcast.
// Returns nil if no transactions are in queue
func (eb *EthBroadcaster) nextUnstartedTransactionWithNonce(fromAddress gethCommon.Address) (*EthTx, error) {
	etx := &EthTx{}
	if err := findNextUnstartedTransactionFromAddress(eb.db, etx, fromAddress, eb.chainID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Finish. No more transactions left to process. Hoorah!
			return nil, nil
		}
		return nil, errors.Wrap(err, "findNextUnstartedTransactionFromAddress failed")
	}

	nonce, err := eb.getNextNonce(etx.FromAddress)
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
	return eb.q.Transaction(func(tx pg.Queryer) error {
		query, args, e := tx.BindNamed(insertIntoEthTxAttemptsQuery, attempt)
		if e != nil {
			return errors.Wrap(e, "failed to BindNamed")
		}
		err := tx.Get(attempt, query, args...)
		if err != nil {
			switch e := err.(type) {
			case *pgconn.PgError:
				if e.ConstraintName == "eth_tx_attempts_eth_tx_id_fkey" {
					return errEthTxRemoved
				}
			}
			return errors.Wrap(err, "saveInProgressTransaction failed to create eth_tx_attempt")
		}
		err = tx.Get(etx, `UPDATE eth_txes SET nonce=$1, state=$2, broadcast_at=$3, initial_broadcast_at=$4 WHERE id=$5 RETURNING *`, etx.Nonce, etx.State, etx.BroadcastAt, etx.InitialBroadcastAt, etx.ID)
		return errors.Wrap(err, "saveInProgressTransaction failed to save eth_tx")
	})
}

// Finds earliest saved transaction that has yet to be broadcast from the given address
func findNextUnstartedTransactionFromAddress(db *sqlx.DB, etx *EthTx, fromAddress gethCommon.Address, chainID big.Int) error {
	err := db.Get(etx, `SELECT * FROM eth_txes WHERE from_address = $1 AND state = 'unstarted' AND evm_chain_id = $2 ORDER BY value ASC, created_at ASC, id ASC`, fromAddress, chainID.String())
	return errors.Wrap(err, "failed to findNextUnstartedTransactionFromAddress")
}

func (eb *EthBroadcaster) saveAttempt(etx *EthTx, attempt EthTxAttempt, NewAttemptState EthTxAttemptState, callbacks ...func(tx pg.Queryer) error) error {
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
	return eb.q.Transaction(func(tx pg.Queryer) error {
		if err := eb.incrementNextNonce(etx.FromAddress, *etx.Nonce, pg.WithQueryer(tx)); err != nil {
			return errors.Wrap(err, "saveUnconfirmed failed")
		}
		if err := tx.Get(etx, `UPDATE eth_txes SET state=$1, error=$2, broadcast_at=$3, initial_broadcast_at=$4 WHERE id = $5 RETURNING *`, etx.State, etx.Error, etx.BroadcastAt, etx.InitialBroadcastAt, etx.ID); err != nil {
			return errors.Wrap(err, "saveUnconfirmed failed to save eth_tx")
		}
		if err := tx.Get(&attempt, `UPDATE eth_tx_attempts SET state = $1 WHERE id = $2 RETURNING *`, attempt.State, attempt.ID); err != nil {
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

func (eb *EthBroadcaster) tryAgainBumpingGas(ctx context.Context, lgr logger.Logger, sendError *evmclient.SendError, etx EthTx, attempt EthTxAttempt, initialBroadcastAt time.Time) (err error, retryable bool) {
	lgr.With(
		"sendError", sendError,
		"attemptGasFeeCap", attempt.GasFeeCap,
		"attemptGasPrice", attempt.GasPrice,
		"attemptGasTipCap", attempt.GasTipCap,
		"maxGasPriceConfig", eb.config.EvmMaxGasPriceWei(),
	).Errorf("attempt gas price %v was rejected by the eth node for being too low. "+
		"Eth node returned: '%s'. "+
		"Will bump and retry. ACTION REQUIRED: This is a configuration error. "+
		"Consider increasing ETH_GAS_PRICE_DEFAULT (current value: %s)",
		attempt.GasPrice, sendError.Error(), eb.config.EvmGasPriceDefault().String())
	switch attempt.TxType {
	case 0x0:
		return eb.tryAgainBumpingLegacyGas(ctx, lgr, etx, attempt, initialBroadcastAt)
	case 0x2:
		return eb.tryAgainBumpingDynamicFeeGas(ctx, lgr, etx, attempt, initialBroadcastAt)
	default:
		err = errors.Errorf("invariant violation: Attempt %v had unrecognised transaction type %v"+
			"This is a bug! Please report to https://github.com/smartcontractkit/chainlink/issues", attempt.ID, attempt.TxType)
		logger.Sugared(eb.logger).AssumptionViolation(err.Error())
		return err, false
	}
}

func (eb *EthBroadcaster) tryAgainBumpingLegacyGas(ctx context.Context, lgr logger.Logger, etx EthTx, attempt EthTxAttempt, initialBroadcastAt time.Time) (err error, retryable bool) {
	keySpecificMaxGasPriceWei := eb.config.KeySpecificMaxGasPriceWei(etx.FromAddress)
	bumpedGasPrice, bumpedGasLimit, err := eb.estimator.BumpLegacyGas(ctx, attempt.GasPrice, etx.GasLimit, keySpecificMaxGasPriceWei, nil)
	if err != nil {
		return errors.Wrap(err, "tryAgainBumpingLegacyGas failed"), true
	}
	if bumpedGasPrice.Cmp(attempt.GasPrice) == 0 || bumpedGasPrice.Cmp(eb.config.EvmMaxGasPriceWei()) >= 0 {
		return errors.Errorf("hit gas price bump ceiling, will not bump further"), true // TODO: Is this terminal or retryable? Is it possible to send unsaved attempts here?
	}
	return eb.tryAgainWithNewLegacyGas(ctx, lgr, etx, attempt, initialBroadcastAt, bumpedGasPrice, bumpedGasLimit)
}

func (eb *EthBroadcaster) tryAgainBumpingDynamicFeeGas(ctx context.Context, lgr logger.Logger, etx EthTx, attempt EthTxAttempt, initialBroadcastAt time.Time) (err error, retryable bool) {
	keySpecificMaxGasPriceWei := eb.config.KeySpecificMaxGasPriceWei(etx.FromAddress)
	bumpedFee, bumpedGasLimit, err := eb.estimator.BumpDynamicFee(ctx, attempt.DynamicFee(), etx.GasLimit, keySpecificMaxGasPriceWei, nil)
	if err != nil {
		return errors.Wrap(err, "tryAgainBumpingDynamicFeeGas failed"), true
	}
	if bumpedFee.TipCap.Cmp(attempt.GasTipCap) == 0 || bumpedFee.FeeCap.Cmp(attempt.GasFeeCap) == 0 || bumpedFee.TipCap.Cmp(eb.config.EvmMaxGasPriceWei()) >= 0 || bumpedFee.TipCap.Cmp(eb.config.EvmMaxGasPriceWei()) >= 0 {
		return errors.Errorf("hit gas price bump ceiling, will not bump further"), true // TODO: Is this terminal or retryable? Is it possible to send unsaved attempts here?
	}
	return eb.tryAgainWithNewDynamicFeeGas(ctx, lgr, etx, attempt, initialBroadcastAt, bumpedFee, bumpedGasLimit)
}

func (eb *EthBroadcaster) tryAgainWithNewEstimation(ctx context.Context, lgr logger.Logger, sendError *evmclient.SendError, etx EthTx, attempt EthTxAttempt, initialBroadcastAt time.Time) (err error, retryable bool) {
	if attempt.TxType == 0x2 {
		err = errors.Errorf("re-estimation is not supported for EIP-1559 transactions. Eth node returned error: %v. This is a bug", sendError.Error())
		logger.Sugared(eb.logger).AssumptionViolation(err.Error())
		return err, false
	}
	keySpecificMaxGasPriceWei := eb.config.KeySpecificMaxGasPriceWei(etx.FromAddress)
	gasPrice, gasLimit, err := eb.estimator.GetLegacyGas(ctx, etx.EncodedPayload, etx.GasLimit, keySpecificMaxGasPriceWei, gas.OptForceRefetch)
	if err != nil {
		return errors.Wrap(err, "tryAgainWithNewEstimation failed to estimate gas"), true
	}
	lgr.Warnw("L2 rejected transaction due to incorrect fee, re-estimated and will try again",
		"etxID", etx.ID, "err", err, "newGasPrice", gasPrice, "newGasLimit", gasLimit)
	return eb.tryAgainWithNewLegacyGas(ctx, lgr, etx, attempt, initialBroadcastAt, gasPrice, gasLimit)
}

func (eb *EthBroadcaster) tryAgainWithNewLegacyGas(ctx context.Context, lgr logger.Logger, etx EthTx, attempt EthTxAttempt, initialBroadcastAt time.Time, newGasPrice *assets.Wei, newGasLimit uint32) (err error, retyrable bool) {
	replacementAttempt, err := eb.NewLegacyAttempt(etx, newGasPrice, newGasLimit)
	if err != nil {
		return errors.Wrap(err, "tryAgainWithNewLegacyGas failed"), true
	}

	if err = eb.orm.SaveReplacementInProgressAttempt(attempt, &replacementAttempt); err != nil {
		return errors.Wrap(err, "tryAgainWithNewLegacyGas failed"), true
	}
	lgr.Debugw("Bumped legacy gas on initial send", "oldGasPrice", attempt.GasPrice, "newGasPrice", newGasPrice)
	return eb.handleInProgressEthTx(ctx, etx, replacementAttempt, initialBroadcastAt)
}

func (eb *EthBroadcaster) tryAgainWithNewDynamicFeeGas(ctx context.Context, lgr logger.Logger, etx EthTx, attempt EthTxAttempt, initialBroadcastAt time.Time, newDynamicFee gas.DynamicFee, newGasLimit uint32) (err error, retyrable bool) {
	replacementAttempt, err := eb.NewDynamicFeeAttempt(etx, newDynamicFee, newGasLimit)
	if err != nil {
		return errors.Wrap(err, "tryAgainWithNewDynamicFeeGas failed"), true
	}

	if err = eb.orm.SaveReplacementInProgressAttempt(attempt, &replacementAttempt); err != nil {
		return errors.Wrap(err, "tryAgainWithNewDynamicFeeGas failed"), true
	}
	lgr.Debugw("Bumped dynamic fee gas on initial send", "oldFee", attempt.DynamicFee(), "newFee", newDynamicFee)
	return eb.handleInProgressEthTx(ctx, etx, replacementAttempt, initialBroadcastAt)
}

func (eb *EthBroadcaster) saveFatallyErroredTransaction(lgr logger.Logger, etx *EthTx) error {
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
			lgr.Debugw("callback missing or already resumed", "etxID", etx.ID)
		} else if err != nil {
			return errors.Wrap(err, "failed to resume pipeline")
		}
	}
	etx.Nonce = nil
	etx.State = EthTxFatalError
	return eb.q.Transaction(func(tx pg.Queryer) error {
		if _, err := tx.Exec(`DELETE FROM eth_tx_attempts WHERE eth_tx_id = $1`, etx.ID); err != nil {
			return errors.Wrapf(err, "saveFatallyErroredTransaction failed to delete eth_tx_attempt with eth_tx.ID %v", etx.ID)
		}
		return errors.Wrap(tx.Get(etx, `UPDATE eth_txes SET state=$1, error=$2, broadcast_at=NULL, initial_broadcast_at=NULL, nonce=NULL WHERE id=$3 RETURNING *`, etx.State, etx.Error, etx.ID), "saveFatallyErroredTransaction failed to save eth_tx")
	})
}

func (eb *EthBroadcaster) getNextNonce(address gethCommon.Address) (nonce int64, err error) {
	return eb.ChainKeyStore.keystore.GetNextNonce(address, &eb.chainID)
}

func (eb *EthBroadcaster) incrementNextNonce(address gethCommon.Address, currentNonce int64, qopts ...pg.QOpt) error {
	return eb.ChainKeyStore.keystore.IncrementNextNonce(address, &eb.chainID, currentNonce, qopts...)
}

func observeTimeUntilBroadcast(chainID big.Int, createdAt, broadcastAt time.Time) {
	duration := float64(broadcastAt.Sub(createdAt))
	promTimeUntilBroadcast.WithLabelValues(chainID.String()).Observe(duration)
}
