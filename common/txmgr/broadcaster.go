package txmgr

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v4"

	clienttypes "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/label"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	// InFlightTransactionRecheckInterval controls how often the Broadcaster
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

var ErrEthTxRemoved = errors.New("eth_tx removed")

type ProcessUnstartedTxs[ADDR types.Hashable] func(ctx context.Context, fromAddress ADDR) (retryable bool, err error)

// TransmitCheckerFactory creates a transmit checker based on a spec.
type TransmitCheckerFactory[
	CHAIN_ID txmgrtypes.ID,
	ADDR types.Hashable,
	TX_HASH, BLOCK_HASH types.Hashable,
	SEQ txmgrtypes.Sequence,
	FEE feetypes.Fee,
] interface {
	// BuildChecker builds a new TransmitChecker based on the given spec.
	BuildChecker(spec txmgrtypes.TransmitCheckerSpec[ADDR]) (TransmitChecker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error)
}

// TransmitChecker determines whether a transaction should be submitted on-chain.
type TransmitChecker[
	CHAIN_ID txmgrtypes.ID,
	ADDR types.Hashable,
	TX_HASH, BLOCK_HASH types.Hashable,
	SEQ txmgrtypes.Sequence,
	FEE feetypes.Fee,
] interface {

	// Check the given transaction. If the transaction should not be sent, an error indicating why
	// is returned. Errors should only be returned if the checker can confirm that a transaction
	// should not be sent, other errors (for example connection or other unexpected errors) should
	// be logged and swallowed.
	Check(ctx context.Context, l logger.Logger, tx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], a txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error
}

// Broadcaster monitors eth_txes for transactions that need to
// be broadcast, assigns sequences and ensures that at least one eth node
// somewhere has received the transaction successfully.
//
// This does not guarantee delivery! A whole host of other things can
// subsequently go wrong such as transactions being evicted from the mempool,
// eth nodes going offline etc. Responsibility for ensuring eventual inclusion
// into the chain falls on the shoulders of the ethConfirmer.
//
// What Broadcaster does guarantee is:
// - a monotonic series of increasing sequences for eth_txes that can all eventually be confirmed if you retry enough times
// - transition of eth_txes out of unstarted into either fatal_error or unconfirmed
// - existence of a saved eth_tx_attempt
type Broadcaster[
	CHAIN_ID txmgrtypes.ID,
	HEAD types.Head[BLOCK_HASH],
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	SEQ txmgrtypes.Sequence,
	FEE feetypes.Fee,
] struct {
	logger  logger.Logger
	txStore txmgrtypes.TransactionStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, SEQ, FEE]
	client  txmgrtypes.TransactionClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	txmgrtypes.TxAttemptBuilder[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	sequenceSyncer SequenceSyncer[ADDR, TX_HASH, BLOCK_HASH]
	resumeCallback ResumeCallback
	chainID        CHAIN_ID
	config         txmgrtypes.BroadcasterConfig
	txConfig       txmgrtypes.BroadcasterTransactionsConfig
	listenerConfig txmgrtypes.BroadcasterListenerConfig

	// autoSyncSequence, if set, will cause Broadcaster to fast-forward the sequence
	// when Start is called
	autoSyncSequence bool

	txInsertListener        pg.Subscription
	eventBroadcaster        pg.EventBroadcaster
	processUnstartedTxsImpl ProcessUnstartedTxs[ADDR]

	ks               txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ]
	enabledAddresses []ADDR

	checkerFactory TransmitCheckerFactory[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]

	// triggers allow other goroutines to force Broadcaster to rescan the
	// database early (before the next poll interval)
	// Each key has its own trigger
	triggers map[ADDR]chan struct{}

	chStop utils.StopChan
	wg     sync.WaitGroup

	initSync  sync.Mutex
	isStarted bool
	utils.StartStopOnce

	parseAddr func(string) (ADDR, error)
}

func NewBroadcaster[
	CHAIN_ID txmgrtypes.ID,
	HEAD types.Head[BLOCK_HASH],
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	SEQ txmgrtypes.Sequence,
	FEE feetypes.Fee,
](
	txStore txmgrtypes.TransactionStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, SEQ, FEE],
	client txmgrtypes.TransactionClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	config txmgrtypes.BroadcasterConfig,
	txConfig txmgrtypes.BroadcasterTransactionsConfig,
	listenerConfig txmgrtypes.BroadcasterListenerConfig,
	keystore txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ],
	eventBroadcaster pg.EventBroadcaster,
	txAttemptBuilder txmgrtypes.TxAttemptBuilder[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	sequenceSyncer SequenceSyncer[ADDR, TX_HASH, BLOCK_HASH],
	logger logger.Logger,
	checkerFactory TransmitCheckerFactory[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	autoSyncSequence bool,
	parseAddress func(string) (ADDR, error),
) *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	logger = logger.Named("Broadcaster")
	b := &Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{
		logger:           logger,
		txStore:          txStore,
		client:           client,
		TxAttemptBuilder: txAttemptBuilder,
		sequenceSyncer:   sequenceSyncer,
		chainID:          client.ConfiguredChainID(),
		config:           config,
		txConfig:         txConfig,
		listenerConfig:   listenerConfig,
		eventBroadcaster: eventBroadcaster,
		ks:               keystore,
		checkerFactory:   checkerFactory,
		autoSyncSequence: autoSyncSequence,
		parseAddr:        parseAddress,
	}

	b.processUnstartedTxsImpl = b.processUnstartedTxs
	return b
}

// Start starts Broadcaster service.
// The provided context can be used to terminate Start sequence.
func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) Start(_ context.Context) error {
	return eb.StartOnce("Broadcaster", func() (err error) {
		return eb.startInternal()
	})
}

// startInternal can be called multiple times, in conjunction with closeInternal. The TxMgr uses this functionality to reset broadcaster multiple times in its own lifetime.
func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) startInternal() error {
	eb.initSync.Lock()
	defer eb.initSync.Unlock()
	if eb.isStarted {
		return errors.New("Broadcaster is already started")
	}
	var err error
	eb.txInsertListener, err = eb.eventBroadcaster.Subscribe(pg.ChannelInsertOnEthTx, "")
	if err != nil {
		return errors.Wrap(err, "Broadcaster could not start")
	}
	eb.enabledAddresses, err = eb.ks.EnabledAddressesForChain(eb.chainID)
	if err != nil {
		return errors.Wrap(err, "Broadcaster: failed to load EnabledAddressesForChain")
	}

	if len(eb.enabledAddresses) > 0 {
		eb.logger.Debugw(fmt.Sprintf("Booting with %d keys", len(eb.enabledAddresses)), "keys", eb.enabledAddresses)
	} else {
		eb.logger.Warnf("Chain %s does not have any eth keys, no transactions will be sent on this chain", eb.chainID.String())
	}
	eb.chStop = make(chan struct{})
	eb.wg = sync.WaitGroup{}
	eb.wg.Add(len(eb.enabledAddresses))
	eb.triggers = make(map[ADDR]chan struct{})
	for _, addr := range eb.enabledAddresses {
		triggerCh := make(chan struct{}, 1)
		eb.triggers[addr] = triggerCh
		go eb.monitorEthTxs(addr, triggerCh)
	}

	eb.wg.Add(1)
	go eb.ethTxInsertTriggerer()

	eb.isStarted = true
	return nil
}

// Close closes the Broadcaster
func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) Close() error {
	return eb.StopOnce("Broadcaster", func() error {
		return eb.closeInternal()
	})
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) closeInternal() error {
	eb.initSync.Lock()
	defer eb.initSync.Unlock()
	if !eb.isStarted {
		return errors.Wrap(utils.ErrAlreadyStopped, "Broadcaster is not started")
	}
	if eb.txInsertListener != nil {
		eb.txInsertListener.Close()
	}
	close(eb.chStop)
	eb.wg.Wait()
	eb.isStarted = false
	return nil
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) SetResumeCallback(callback ResumeCallback) {
	eb.resumeCallback = callback
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) Name() string {
	return eb.logger.Name()
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) HealthReport() map[string]error {
	return map[string]error{eb.Name(): eb.StartStopOnce.Healthy()}
}

// Trigger forces the monitor for a particular address to recheck for new eth_txes
// Logs error and does nothing if address was not registered on startup
func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) Trigger(addr ADDR) {
	if eb.isStarted {
		triggerCh, exists := eb.triggers[addr]
		if !exists {
			// ignoring trigger for address which is not registered with this Broadcaster
			return
		}
		select {
		case triggerCh <- struct{}{}:
		default:
		}
	} else {
		eb.logger.Debugf("Unstarted; ignoring trigger for %s", addr)
	}
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) ethTxInsertTriggerer() {
	defer eb.wg.Done()
	for {
		select {
		case ev, ok := <-eb.txInsertListener.Events():
			if !ok {
				eb.logger.Debug("txInsertListener channel closed, exiting trigger loop")
				return
			}
			addr, err := eb.parseAddr(ev.Payload)
			if err != nil {
				eb.logger.Errorw("failed to parse address in trigger", "error", err)
				continue
			}
			eb.Trigger(addr)
		case <-eb.chStop:
			return
		}
	}
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) newSequenceSyncBackoff() backoff.Backoff {
	return backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    5 * time.Second,
		Jitter: true,
	}
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) newResendBackoff() backoff.Backoff {
	return backoff.Backoff{
		Min:    1 * time.Second,
		Max:    15 * time.Second,
		Jitter: true,
	}
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) monitorEthTxs(addr ADDR, triggerCh chan struct{}) {
	defer eb.wg.Done()

	ctx, cancel := eb.chStop.NewCtx()
	defer cancel()

	if eb.autoSyncSequence {
		eb.logger.Debugw("Auto-syncing sequence", "address", addr.String())
		eb.SyncSequence(ctx, addr)
		if ctx.Err() != nil {
			return
		}
	} else {
		eb.logger.Debugw("Skipping sequence auto-sync", "address", addr.String())
	}

	// errorRetryCh allows retry on exponential backoff in case of timeout or
	// other unknown error
	var errorRetryCh <-chan time.Time
	bf := eb.newResendBackoff()

	for {
		pollDBTimer := time.NewTimer(utils.WithJitter(eb.listenerConfig.FallbackPollInterval()))

		retryable, err := eb.processUnstartedTxsImpl(ctx, addr)
		if err != nil {
			eb.logger.Errorw("Error occurred while handling eth_tx queue in ProcessUnstartedTxs", "err", err)
		}
		// On retryable errors we implement exponential backoff retries. This
		// handles intermittent connectivity, remote RPC races, timing issues etc
		if retryable {
			pollDBTimer.Reset(utils.WithJitter(eb.listenerConfig.FallbackPollInterval()))
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

// syncSequence tries to sync the key sequence, retrying indefinitely until success
func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) SyncSequence(ctx context.Context, addr ADDR) {
	sequenceSyncRetryBackoff := eb.newSequenceSyncBackoff()
	if err := eb.sequenceSyncer.Sync(ctx, addr); err != nil {
		// Enter retry loop with backoff
		var attempt int
		eb.logger.Errorw("Failed to sync with on-chain sequence", "address", addr.String(), "attempt", attempt, "err", err)
		for {
			select {
			case <-eb.chStop:
				return
			case <-time.After(sequenceSyncRetryBackoff.Duration()):
				attempt++

				if err := eb.sequenceSyncer.Sync(ctx, addr); err != nil {
					if attempt > 5 {
						eb.logger.Criticalw("Failed to sync with on-chain sequence", "address", addr.String(), "attempt", attempt, "err", err)
						eb.SvcErrBuffer.Append(err)
					} else {
						eb.logger.Warnw("Failed to sync with on-chain sequence", "address", addr.String(), "attempt", attempt, "err", err)
					}
					continue
				}
				return
			}
		}
	}
}

// ProcessUnstartedTxs picks up and handles all eth_txes in the queue
// revive:disable:error-return
func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) ProcessUnstartedTxs(ctx context.Context, addr ADDR) (retryable bool, err error) {
	return eb.processUnstartedTxs(ctx, addr)
}

// NOTE: This MUST NOT be run concurrently for the same address or it could
// result in undefined state or deadlocks.
// First handle any in_progress transactions left over from last time.
// Then keep looking up unstarted transactions and processing them until there are none remaining.
func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) processUnstartedTxs(ctx context.Context, fromAddress ADDR) (retryable bool, err error) {
	var n uint
	mark := time.Now()
	defer func() {
		if n > 0 {
			eb.logger.Debugw("Finished processUnstartedTxs", "address", fromAddress, "time", time.Since(mark), "n", n, "id", "eth_broadcaster")
		}
	}()

	err, retryable = eb.handleAnyInProgressEthTx(ctx, fromAddress)
	if err != nil {
		return retryable, errors.Wrap(err, "processUnstartedTxs failed on handleAnyInProgressEthTx")
	}
	for {
		maxInFlightTransactions := eb.txConfig.MaxInFlight()
		if maxInFlightTransactions > 0 {
			nUnconfirmed, err := eb.txStore.CountUnconfirmedTransactions(fromAddress, eb.chainID)
			if err != nil {
				return true, errors.Wrap(err, "CountUnconfirmedTransactions failed")
			}
			if nUnconfirmed >= maxInFlightTransactions {
				nUnstarted, err := eb.txStore.CountUnstartedTransactions(fromAddress, eb.chainID)
				if err != nil {
					return true, errors.Wrap(err, "CountUnstartedTransactions failed")
				}
				eb.logger.Warnw(fmt.Sprintf(`Transaction throttling; %d transactions in-flight and %d unstarted transactions pending (maximum number of in-flight transactions is %d per key). %s`, nUnconfirmed, nUnstarted, maxInFlightTransactions, label.MaxInFlightTransactionsWarning), "maxInFlightTransactions", maxInFlightTransactions, "nUnconfirmed", nUnconfirmed, "nUnstarted", nUnstarted)
				select {
				case <-time.After(InFlightTransactionRecheckInterval):
				case <-ctx.Done():
					return false, context.Cause(ctx)
				}
				continue
			}
		}
		etx, err := eb.nextUnstartedTransactionWithSequence(fromAddress)
		if err != nil {
			return true, errors.Wrap(err, "processUnstartedTxs failed on nextUnstartedTransactionWithSequence")
		}
		if etx == nil {
			return false, nil
		}
		n++
		var a txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
		var retryable bool
		a, _, _, retryable, err = eb.NewTxAttempt(ctx, *etx, eb.logger)
		if err != nil {
			return retryable, errors.Wrap(err, "processUnstartedTxs failed on NewAttempt")
		}

		if err := eb.txStore.UpdateTxUnstartedToInProgress(etx, &a); errors.Is(err, ErrEthTxRemoved) {
			eb.logger.Debugw("eth_tx removed", "etxID", etx.ID, "subject", etx.Subject)
			continue
		} else if err != nil {
			return true, errors.Wrap(err, "processUnstartedTxs failed on UpdateTxUnstartedToInProgress")
		}

		if err, retryable := eb.handleInProgressEthTx(ctx, *etx, a, time.Now()); err != nil {
			return retryable, errors.Wrap(err, "processUnstartedTxs failed on handleAnyInProgressEthTx")
		}
	}
}

// handleInProgressEthTx checks if there is any transaction
// in_progress and if so, finishes the job
func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) handleAnyInProgressEthTx(ctx context.Context, fromAddress ADDR) (err error, retryable bool) {
	etx, err := eb.txStore.GetTxInProgress(fromAddress)
	if err != nil {
		return errors.Wrap(err, "handleAnyInProgressEthTx failed"), true
	}
	if etx != nil {
		if err, retryable := eb.handleInProgressEthTx(ctx, *etx, etx.TxAttempts[0], etx.CreatedAt); err != nil {
			return errors.Wrap(err, "handleAnyInProgressEthTx failed"), retryable
		}
	}
	return nil, false
}

// This function is used to pass the queryer from the txmgr to the keystore.
// It is inevitable we have to pass the queryer because we need the keystate's next sequence to be incremented
// atomically alongside the transition from `in_progress` to `broadcast` so it is ready for the next transaction
func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) incrementNextSequenceAtomic(tx pg.Queryer, etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	if err := eb.incrementNextSequence(etx.FromAddress, *etx.Sequence, pg.WithQueryer(tx)); err != nil {
		return errors.Wrap(err, "saveUnconfirmed failed")
	}
	return nil
}

// There can be at most one in_progress transaction per address.
// Here we complete the job that we didn't finish last time.
func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) handleInProgressEthTx(ctx context.Context, etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], initialBroadcastAt time.Time) (error, bool) {
	if etx.State != TxInProgress {
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

	lgr := etx.GetLogger(eb.logger.With("fee", attempt.TxFee))

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

	lgr.Infow("Sending transaction", "ethTxAttemptID", attempt.ID, "txHash", attempt.Hash, "err", err, "meta", etx.Meta, "feeLimit", etx.FeeLimit, "attempt", attempt, "etx", etx)
	errType, err := eb.client.SendTransactionReturnCode(ctx, etx, attempt, lgr)

	if errType != clienttypes.Fatal {
		etx.InitialBroadcastAt = &initialBroadcastAt
		etx.BroadcastAt = &initialBroadcastAt
	}

	switch errType {
	case clienttypes.Fatal:
		eb.SvcErrBuffer.Append(err)
		etx.Error = null.StringFrom(err.Error())
		return eb.saveFatallyErroredTransaction(lgr, &etx), true
	case clienttypes.TransactionAlreadyKnown:
		fallthrough
	case clienttypes.Successful:
		// Either the transaction was successful or one of the following four scenarios happened:
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
		// account and sent a transaction on this sequence.
		//
		// In this case, the onus is on the node operator since this is
		// explicitly unsupported.
		//
		// If it turns out to have been an external wallet, we will never get a
		// receipt for this transaction and it will eventually be marked as
		// errored.
		//
		// The end result is that we will NOT SEND a transaction for this
		// sequence.
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
		observeTimeUntilBroadcast(eb.chainID, etx.CreatedAt, time.Now())
		return eb.txStore.UpdateTxAttemptInProgressToBroadcast(&etx, attempt, txmgrtypes.TxAttemptBroadcast, func(tx pg.Queryer) error {
			return eb.incrementNextSequenceAtomic(tx, etx)
		}), true
	case clienttypes.Underpriced:
		return eb.tryAgainBumpingGas(ctx, lgr, err, etx, attempt, initialBroadcastAt)
	case clienttypes.InsufficientFunds:
		// NOTE: This bails out of the entire cycle and essentially "blocks" on
		// any transaction that gets insufficient_eth. This is OK if a
		// transaction with a large VALUE blocks because this always comes last
		// in the processing list.
		// If it blocks because of a transaction that is expensive due to large
		// gas limit, we could have smaller transactions "above" it that could
		// theoretically be sent, but will instead be blocked.
		eb.SvcErrBuffer.Append(err)
		fallthrough
	case clienttypes.Retryable:
		return err, true
	case clienttypes.FeeOutOfValidRange:
		return eb.tryAgainWithNewEstimation(ctx, lgr, err, etx, attempt, initialBroadcastAt)
	case clienttypes.Unsupported:
		return err, false
	case clienttypes.ExceedsMaxFee:
		// Broadcaster: Note that we may have broadcast to multiple nodes and had it
		// accepted by one of them! It is not guaranteed that all nodes share
		// the same tx fee cap. That is why we must treat this as an unknown
		// error that may have been confirmed.
		// If there is only one RPC node, or all RPC nodes have the same
		// configured cap, this transaction will get stuck and keep repeating
		// forever until the issue is resolved.
		lgr.Criticalw(`RPC node rejected this tx as outside Fee Cap`)
		fallthrough
	default:
		// Every error that doesn't fall under one of the above categories will be treated as Unknown.
		fallthrough
	case clienttypes.Unknown:
		eb.SvcErrBuffer.Append(err)
		lgr.Criticalw(`Unknown error occurred while handling eth_tx queue in ProcessUnstartedTxs. This chain/RPC client may not be supported. `+
			`Urgent resolution required, Chainlink is currently operating in a degraded state and may miss transactions`, "err", err, "etx", etx, "attempt", attempt)
		nextSequence, e := eb.client.PendingSequenceAt(ctx, etx.FromAddress)
		if e != nil {
			err = multierr.Combine(e, err)
			return errors.Wrapf(err, "failed to fetch latest pending sequence after encountering unknown RPC error while sending transaction"), true
		}
		if nextSequence.Int64() > (*etx.Sequence).Int64() {
			// Despite the error, the RPC node considers the previously sent
			// transaction to have been accepted. In this case, the right thing to
			// do is assume success and hand off to EthConfirmer
			return eb.txStore.UpdateTxAttemptInProgressToBroadcast(&etx, attempt, txmgrtypes.TxAttemptBroadcast, func(tx pg.Queryer) error {
				return eb.incrementNextSequenceAtomic(tx, etx)
			}), true
		}
		// Either the unknown error prevented the transaction from being mined, or
		// it has not yet propagated to the mempool, or there is some race on the
		// remote RPC.
		//
		// In all cases, the best thing we can do is go into a retry loop and keep
		// trying to send the transaction over again.
		return errors.Wrapf(err, "retryable error while sending transaction %s (eth_tx ID %d)", attempt.Hash.String(), etx.ID), true
	}

}

// Finds next transaction in the queue, assigns a sequence, and moves it to "in_progress" state ready for broadcast.
// Returns nil if no transactions are in queue
func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) nextUnstartedTransactionWithSequence(fromAddress ADDR) (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	etx := &txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	if err := eb.txStore.FindNextUnstartedTransactionFromAddress(etx, fromAddress, eb.chainID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Finish. No more transactions left to process. Hoorah!
			return nil, nil
		}
		return nil, errors.Wrap(err, "findNextUnstartedTransactionFromAddress failed")
	}

	sequence, err := eb.getNextSequence(etx.FromAddress)
	if err != nil {
		return nil, err
	}
	etx.Sequence = &sequence
	return etx, nil
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) tryAgainBumpingGas(ctx context.Context, lgr logger.Logger, txError error, etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], initialBroadcastAt time.Time) (err error, retryable bool) {
	lgr.With(
		"sendError", txError,
		"attemptFee", attempt.TxFee,
		"maxGasPriceConfig", eb.config.MaxFeePrice(),
	).Errorf("attempt fee %v was rejected by the node for being too low. "+
		"Node returned: '%s'. "+
		"Will bump and retry. ACTION REQUIRED: This is a configuration error. "+
		"Consider increasing FeeEstimator.PriceDefault (current value: %s)",
		attempt.TxFee, txError.Error(), eb.config.FeePriceDefault())

	replacementAttempt, bumpedFee, bumpedFeeLimit, retryable, err := eb.NewBumpTxAttempt(ctx, etx, attempt, nil, lgr)
	if err != nil {
		return errors.Wrap(err, "tryAgainBumpFee failed"), retryable
	}

	return eb.saveTryAgainAttempt(ctx, lgr, etx, attempt, replacementAttempt, initialBroadcastAt, bumpedFee, bumpedFeeLimit)
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) tryAgainWithNewEstimation(ctx context.Context, lgr logger.Logger, txError error, etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], initialBroadcastAt time.Time) (err error, retryable bool) {
	if attempt.TxType == 0x2 {
		err = errors.Errorf("re-estimation is not supported for EIP-1559 transactions. Eth node returned error: %v. This is a bug", txError.Error())
		logger.Sugared(eb.logger).AssumptionViolation(err.Error())
		return err, false
	}

	replacementAttempt, fee, feeLimit, retryable, err := eb.NewTxAttemptWithType(ctx, etx, lgr, attempt.TxType, feetypes.OptForceRefetch)
	if err != nil {
		return errors.Wrap(err, "tryAgainWithNewEstimation failed to build new attempt"), retryable
	}
	lgr.Warnw("L2 rejected transaction due to incorrect fee, re-estimated and will try again",
		"etxID", etx.ID, "err", err, "newGasPrice", fee, "newGasLimit", feeLimit)

	return eb.saveTryAgainAttempt(ctx, lgr, etx, attempt, replacementAttempt, initialBroadcastAt, fee, feeLimit)
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) saveTryAgainAttempt(ctx context.Context, lgr logger.Logger, etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], replacementAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], initialBroadcastAt time.Time, newFee FEE, newFeeLimit uint32) (err error, retyrable bool) {
	if err = eb.txStore.SaveReplacementInProgressAttempt(attempt, &replacementAttempt); err != nil {
		return errors.Wrap(err, "tryAgainWithNewFee failed"), true
	}
	lgr.Debugw("Bumped fee on initial send", "oldFee", attempt.TxFee.String(), "newFee", newFee.String(), "newFeeLimit", newFeeLimit)
	return eb.handleInProgressEthTx(ctx, etx, replacementAttempt, initialBroadcastAt)
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) saveFatallyErroredTransaction(lgr logger.Logger, etx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	if etx.State != TxInProgress {
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
	return eb.txStore.UpdateTxFatalError(etx)
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) getNextSequence(address ADDR) (sequence SEQ, err error) {
	return eb.ks.NextSequence(address, eb.chainID)
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) incrementNextSequence(address ADDR, currentSequence SEQ, qopts ...pg.QOpt) error {
	return eb.ks.IncrementNextSequence(address, eb.chainID, currentSequence, qopts...)
}

func observeTimeUntilBroadcast[CHAIN_ID txmgrtypes.ID](chainID CHAIN_ID, createdAt, broadcastAt time.Time) {
	duration := float64(broadcastAt.Sub(createdAt))
	promTimeUntilBroadcast.WithLabelValues(chainID.String()).Observe(duration)
}
