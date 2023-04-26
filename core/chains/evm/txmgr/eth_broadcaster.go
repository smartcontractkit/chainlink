package txmgr

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v4"

	clienttypes "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/label"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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

type ProcessUnstartedEthTxs[ADDR types.Hashable] func(ctx context.Context, fromAddress ADDR) (retryable bool, err error)

// TransmitCheckerFactory creates a transmit checker based on a spec.
type TransmitCheckerFactory[ADDR types.Hashable, TX_HASH types.Hashable] interface {
	// BuildChecker builds a new TransmitChecker based on the given spec.
	BuildChecker(spec txmgrtypes.TransmitCheckerSpec[ADDR]) (TransmitChecker[ADDR, TX_HASH], error)
}

// TransmitChecker determines whether a transaction should be submitted on-chain.
type TransmitChecker[ADDR types.Hashable, TX_HASH types.Hashable] interface {

	// Check the given transaction. If the transaction should not be sent, an error indicating why
	// is returned. Errors should only be returned if the checker can confirm that a transaction
	// should not be sent, other errors (for example connection or other unexpected errors) should
	// be logged and swallowed.
	Check(ctx context.Context, l logger.Logger, tx EthTx[ADDR, TX_HASH], a EthTxAttempt[ADDR, TX_HASH]) error
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
type EthBroadcaster[
	CHAIN_ID txmgrtypes.ID,
	HEAD txmgrtypes.Head,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R any,
	SEQ txmgrtypes.Sequence,
	FEE txmgrtypes.Fee,
] struct {
	logger    logger.Logger
	txStore   txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, txmgrtypes.NewTx[ADDR, TX_HASH], *evmtypes.Receipt, EthTx[ADDR, TX_HASH], EthTxAttempt[ADDR, TX_HASH], SEQ]
	ethClient evmclient.Client
	txmgrtypes.TxAttemptBuilder[HEAD, gas.EvmFee, ADDR, TX_HASH, EthTx[ADDR, TX_HASH], EthTxAttempt[ADDR, TX_HASH], SEQ]
	nonceSyncer    NonceSyncer[ADDR, TX_HASH, BLOCK_HASH]
	resumeCallback ResumeCallback
	chainID        CHAIN_ID
	config         EvmBroadcasterConfig

	// autoSyncNonce, if set, will cause EthBroadcaster to fast-forward the nonce
	// when Start is called
	autoSyncNonce bool

	ethTxInsertListener        pg.Subscription
	eventBroadcaster           pg.EventBroadcaster
	processUnstartedEthTxsImpl ProcessUnstartedEthTxs[ADDR]

	// TODO: When EthTx is generalized for Nonce, then change below nonce type
	// to SEQ. https://smartcontract-it.atlassian.net/browse/BCI-1129
	ks               txmgrtypes.KeyStore[ADDR, CHAIN_ID, evmtypes.Nonce]
	enabledAddresses []ADDR

	checkerFactory TransmitCheckerFactory[ADDR, TX_HASH]

	// triggers allow other goroutines to force EthBroadcaster to rescan the
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

// NewEthBroadcaster returns a new concrete EthBroadcaster
func NewEthBroadcaster(
	txStore EvmTxStore,
	ethClient evmclient.Client,
	config EvmBroadcasterConfig,
	keystore EvmKeyStore,
	eventBroadcaster pg.EventBroadcaster,
	txAttemptBuilder EvmTxAttemptBuilder,
	nonceSyncer EvmNonceSyncer,
	logger logger.Logger,
	checkerFactory EvmTransmitCheckerFactory,
	autoSyncNonce bool,
) *EvmBroadcaster {

	logger = logger.Named("EthBroadcaster")
	b := &EvmBroadcaster{
		logger:           logger,
		txStore:          txStore,
		ethClient:        ethClient,
		TxAttemptBuilder: txAttemptBuilder,
		nonceSyncer:      nonceSyncer,
		chainID:          ethClient.ConfiguredChainID(),
		config:           config,
		eventBroadcaster: eventBroadcaster,
		ks:               keystore,
		checkerFactory:   checkerFactory,
		initSync:         sync.Mutex{},
		isStarted:        false,
		autoSyncNonce:    autoSyncNonce,
		parseAddr:        stringToGethAddress, // note: still evm-specific
	}

	b.processUnstartedEthTxsImpl = b.processUnstartedEthTxs
	return b
}

// Start starts EthBroadcaster service.
// The provided context can be used to terminate Start sequence.
func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Start(_ context.Context) error {
	return eb.StartOnce("EthBroadcaster", func() (err error) {
		return eb.startInternal()
	})
}

// startInternal can be called multiple times, in conjunction with closeInternal. The TxMgr uses this functionality to reset broadcaster multiple times in its own lifetime.
func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) startInternal() error {
	eb.initSync.Lock()
	defer eb.initSync.Unlock()
	if eb.isStarted {
		return errors.New("EthBroadcaster is already started")
	}
	var err error
	eb.ethTxInsertListener, err = eb.eventBroadcaster.Subscribe(pg.ChannelInsertOnEthTx, "")
	if err != nil {
		return errors.Wrap(err, "EthBroadcaster could not start")
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

// Close closes the EthBroadcaster
func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Close() error {
	return eb.StopOnce("EthBroadcaster", func() error {
		return eb.closeInternal()
	})
}

func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) closeInternal() error {
	eb.initSync.Lock()
	defer eb.initSync.Unlock()
	if !eb.isStarted {
		return errors.Wrap(utils.ErrAlreadyStopped, "EthBroadcaster is not started")
	}
	if eb.ethTxInsertListener != nil {
		eb.ethTxInsertListener.Close()
	}
	close(eb.chStop)
	eb.wg.Wait()
	eb.isStarted = false
	return nil
}

func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SetResumeCallback(callback ResumeCallback) {
	eb.resumeCallback = callback
}

func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Name() string {
	return eb.logger.Name()
}

func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) HealthReport() map[string]error {
	return map[string]error{eb.Name(): eb.StartStopOnce.Healthy()}
}

// Trigger forces the monitor for a particular address to recheck for new eth_txes
// Logs error and does nothing if address was not registered on startup
func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Trigger(addr ADDR) {
	if eb.isStarted {
		triggerCh, exists := eb.triggers[addr]
		if !exists {
			// ignoring trigger for address which is not registered with this EthBroadcaster
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

func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) ethTxInsertTriggerer() {
	defer eb.wg.Done()
	for {
		select {
		case ev, ok := <-eb.ethTxInsertListener.Events():
			if !ok {
				eb.logger.Debug("ethTxInsertListener channel closed, exiting trigger loop")
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

func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) newNonceSyncBackoff() backoff.Backoff {
	return backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    5 * time.Second,
		Jitter: true,
	}
}

func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) newResendBackoff() backoff.Backoff {
	return backoff.Backoff{
		Min:    1 * time.Second,
		Max:    15 * time.Second,
		Jitter: true,
	}
}

func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) monitorEthTxs(addr ADDR, triggerCh chan struct{}) {
	defer eb.wg.Done()

	ctx, cancel := eb.chStop.NewCtx()
	defer cancel()

	if eb.autoSyncNonce {
		eb.logger.Debugw("Auto-syncing nonce", "address", addr.String())
		eb.SyncNonce(ctx, addr)
		if ctx.Err() != nil {
			return
		}
	} else {
		eb.logger.Debugw("Skipping nonce auto-sync", "address", addr.String())
	}

	// errorRetryCh allows retry on exponential backoff in case of timeout or
	// other unknown error
	var errorRetryCh <-chan time.Time
	bf := eb.newResendBackoff()

	for {
		pollDBTimer := time.NewTimer(utils.WithJitter(eb.config.TriggerFallbackDBPollInterval()))

		retryable, err := eb.processUnstartedEthTxsImpl(ctx, addr)
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
func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SyncNonce(ctx context.Context, addr ADDR) {
	nonceSyncRetryBackoff := eb.newNonceSyncBackoff()
	if err := eb.nonceSyncer.Sync(ctx, addr); err != nil {
		// Enter retry loop with backoff
		var attempt int
		eb.logger.Errorw("Failed to sync with on-chain nonce", "address", addr.String(), "attempt", attempt, "err", err)
		for {
			select {
			case <-eb.chStop:
				return
			case <-time.After(nonceSyncRetryBackoff.Duration()):
				attempt++

				if err := eb.nonceSyncer.Sync(ctx, addr); err != nil {
					if attempt > 5 {
						eb.logger.Criticalw("Failed to sync with on-chain nonce", "address", addr.String(), "attempt", attempt, "err", err)
						eb.SvcErrBuffer.Append(err)
					} else {
						eb.logger.Warnw("Failed to sync with on-chain nonce", "address", addr.String(), "attempt", attempt, "err", err)
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
func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) ProcessUnstartedEthTxs(ctx context.Context, addr ADDR) (retryable bool, err error) {
	return eb.processUnstartedEthTxs(ctx, addr)
}

// NOTE: This MUST NOT be run concurrently for the same address or it could
// result in undefined state or deadlocks.
// First handle any in_progress transactions left over from last time.
// Then keep looking up unstarted transactions and processing them until there are none remaining.
func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) processUnstartedEthTxs(ctx context.Context, fromAddress ADDR) (retryable bool, err error) {
	var n uint
	mark := time.Now()
	defer func() {
		if n > 0 {
			eb.logger.Debugw("Finished processUnstartedEthTxs", "address", fromAddress, "time", time.Since(mark), "n", n, "id", "eth_broadcaster")
		}
	}()

	err, retryable = eb.handleAnyInProgressEthTx(ctx, fromAddress)
	if err != nil {
		return retryable, errors.Wrap(err, "processUnstartedEthTxs failed on handleAnyInProgressEthTx")
	}
	for {
		maxInFlightTransactions := eb.config.MaxInFlightTransactions()
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
		etx, err := eb.nextUnstartedTransactionWithNonce(fromAddress)
		if err != nil {
			return true, errors.Wrap(err, "processUnstartedEthTxs failed on nextUnstartedTransactionWithNonce")
		}
		if etx == nil {
			return false, nil
		}
		n++
		var a EthTxAttempt[ADDR, TX_HASH]
		var retryable bool
		a, _, _, retryable, err = eb.NewTxAttempt(ctx, *etx, eb.logger)
		if err != nil {
			return retryable, errors.Wrap(err, "processUnstartedEthTxs failed on NewAttempt")
		}

		if err := eb.txStore.UpdateEthTxUnstartedToInProgress(etx, &a); errors.Is(err, errEthTxRemoved) {
			eb.logger.Debugw("eth_tx removed", "etxID", etx.ID, "subject", etx.Subject)
			continue
		} else if err != nil {
			return true, errors.Wrap(err, "processUnstartedEthTxs failed on UpdateEthTxUnstartedToInProgress")
		}

		if err, retryable := eb.handleInProgressEthTx(ctx, *etx, a, time.Now()); err != nil {
			return retryable, errors.Wrap(err, "processUnstartedEthTxs failed on handleAnyInProgressEthTx")
		}
	}
}

// handleInProgressEthTx checks if there is any transaction
// in_progress and if so, finishes the job
func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) handleAnyInProgressEthTx(ctx context.Context, fromAddress ADDR) (err error, retryable bool) {
	etx, err := eb.txStore.GetEthTxInProgress(fromAddress)
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

// This function is used to pass the queryer from the txmgr to the keystore.
// It is inevitable we have to pass the queryer because we need the keystate's next nonce to be incremented
// atomically alongside the transition from `in_progress` to `broadcast` so it is ready for the next transaction
func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) incrementNextNonceAtomic(tx pg.Queryer, etx EthTx[ADDR, TX_HASH]) error {
	if err := eb.incrementNextNonce(etx.FromAddress, evmtypes.Nonce(*etx.Nonce), pg.WithQueryer(tx)); err != nil {
		return errors.Wrap(err, "saveUnconfirmed failed")
	}
	return nil
}

// There can be at most one in_progress transaction per address.
// Here we complete the job that we didn't finish last time.
func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) handleInProgressEthTx(ctx context.Context, etx EthTx[ADDR, TX_HASH], attempt EthTxAttempt[ADDR, TX_HASH], initialBroadcastAt time.Time) (error, bool) {
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

	signedTx, err := attempt.GetSignedTx()
	if err != nil {
		lgr.Criticalw("Fatal error signing transaction", "err", err, "etx", etx)
		etx.Error = null.StringFrom(err.Error())
		return eb.saveFatallyErroredTransaction(lgr, &etx), true
	}

	// TODO: When eth client is generalized, remove this address conversion logic below
	// https://smartcontract-it.atlassian.net/browse/BCI-852
	fromAddress, err := stringToGethAddress(etx.FromAddress.String())
	if err != nil {
		return errors.Wrapf(err, "failed to do address format conversion"), true
	}

	lgr.Debugw("Sending transaction", "ethTxAttemptID", attempt.ID, "txHash", attempt.Hash, "err", err, "meta", etx.Meta, "gasLimit", etx.GasLimit, "attempt", attempt, "etx", etx)
	errType, err := eb.ethClient.SendTransactionReturnCode(ctx, signedTx, fromAddress)

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
		observeTimeUntilBroadcast(eb.chainID, etx.CreatedAt, time.Now())
		return eb.txStore.UpdateEthTxAttemptInProgressToBroadcast(&etx, attempt, txmgrtypes.TxAttemptBroadcast, func(tx pg.Queryer) error {
			return eb.incrementNextNonceAtomic(tx, etx)
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
		lgr.Criticalw(`Unknown error occurred while handling eth_tx queue in ProcessUnstartedEthTxs. This chain/RPC client may not be supported. `+
			`Urgent resolution required, Chainlink is currently operating in a degraded state and may miss transactions`, "err", err, "etx", etx, "attempt", attempt)
		nextNonce, e := eb.ethClient.PendingNonceAt(ctx, fromAddress)
		if e != nil {
			err = multierr.Combine(e, err)
			return errors.Wrapf(err, "failed to fetch latest pending nonce after encountering unknown RPC error while sending transaction"), true
		}
		if nextNonce > math.MaxInt64 {
			return errors.Errorf("nonce overflow, got: %v", nextNonce), true
		}
		if int64(nextNonce) > *etx.Nonce {
			// Despite the error, the RPC node considers the previously sent
			// transaction to have been accepted. In this case, the right thing to
			// do is assume success and hand off to EthConfirmer
			return eb.txStore.UpdateEthTxAttemptInProgressToBroadcast(&etx, attempt, txmgrtypes.TxAttemptBroadcast, func(tx pg.Queryer) error {
				return eb.incrementNextNonceAtomic(tx, etx)
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

// Finds next transaction in the queue, assigns a nonce, and moves it to "in_progress" state ready for broadcast.
// Returns nil if no transactions are in queue
func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) nextUnstartedTransactionWithNonce(fromAddress ADDR) (*EthTx[ADDR, TX_HASH], error) {
	etx := &EthTx[ADDR, TX_HASH]{}
	if err := eb.txStore.FindNextUnstartedTransactionFromAddress(etx, fromAddress, eb.chainID); err != nil {
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
	nonceVal := nonce.Int64()
	etx.Nonce = &nonceVal
	return etx, nil
}

func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) tryAgainBumpingGas(ctx context.Context, lgr logger.Logger, txError error, etx EthTx[ADDR, TX_HASH], attempt EthTxAttempt[ADDR, TX_HASH], initialBroadcastAt time.Time) (err error, retryable bool) {
	lgr.With(
		"sendError", txError,
		"attemptGasFeeCap", attempt.GasFeeCap,
		"attemptGasPrice", attempt.GasPrice,
		"attemptGasTipCap", attempt.GasTipCap,
		"maxGasPriceConfig", eb.config.MaxFeePrice(),
	).Errorf("attempt gas price %v was rejected by the eth node for being too low. "+
		"Eth node returned: '%s'. "+
		"Will bump and retry. ACTION REQUIRED: This is a configuration error. "+
		"Consider increasing EVM.GasEstimator.PriceDefault (current value: %s)",
		attempt.GasPrice, txError.Error(), eb.config.FeePriceDefault().String())

	replacementAttempt, bumpedFee, bumpedFeeLimit, retryable, err := eb.NewBumpTxAttempt(ctx, etx, attempt, nil, lgr)
	if err != nil {
		return errors.Wrap(err, "tryAgainBumpFee failed"), retryable
	}

	return eb.saveTryAgainAttempt(ctx, lgr, etx, attempt, replacementAttempt, initialBroadcastAt, bumpedFee, bumpedFeeLimit)
}

func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) tryAgainWithNewEstimation(ctx context.Context, lgr logger.Logger, txError error, etx EthTx[ADDR, TX_HASH], attempt EthTxAttempt[ADDR, TX_HASH], initialBroadcastAt time.Time) (err error, retryable bool) {
	if attempt.TxType == 0x2 {
		err = errors.Errorf("re-estimation is not supported for EIP-1559 transactions. Eth node returned error: %v. This is a bug", txError.Error())
		logger.Sugared(eb.logger).AssumptionViolation(err.Error())
		return err, false
	}

	replacementAttempt, fee, feeLimit, retryable, err := eb.NewTxAttemptWithType(ctx, etx, lgr, attempt.TxType, txmgrtypes.OptForceRefetch)
	if err != nil {
		return errors.Wrap(err, "tryAgainWithNewEstimation failed to build new attempt"), retryable
	}
	lgr.Warnw("L2 rejected transaction due to incorrect fee, re-estimated and will try again",
		"etxID", etx.ID, "err", err, "newGasPrice", fee, "newGasLimit", feeLimit)

	return eb.saveTryAgainAttempt(ctx, lgr, etx, attempt, replacementAttempt, initialBroadcastAt, fee, feeLimit)
}

func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) saveTryAgainAttempt(ctx context.Context, lgr logger.Logger, etx EthTx[ADDR, TX_HASH], attempt EthTxAttempt[ADDR, TX_HASH], replacementAttempt EthTxAttempt[ADDR, TX_HASH], initialBroadcastAt time.Time, newFee gas.EvmFee, newFeeLimit uint32) (err error, retyrable bool) {
	if err = eb.txStore.SaveReplacementInProgressAttempt(attempt, &replacementAttempt); err != nil {
		return errors.Wrap(err, "tryAgainWithNewFee failed"), true
	}
	lgr.Debugw("Bumped fee on initial send", "oldFee", attempt.Fee().String(), "newFee", newFee.String(), "newFeeLimit", newFeeLimit)
	return eb.handleInProgressEthTx(ctx, etx, replacementAttempt, initialBroadcastAt)
}

func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) saveFatallyErroredTransaction(lgr logger.Logger, etx *EthTx[ADDR, TX_HASH]) error {
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
	return eb.txStore.UpdateEthTxFatalError(etx)
}

func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) getNextNonce(address ADDR) (nonce evmtypes.Nonce, err error) {
	return eb.ks.NextSequence(address, eb.chainID)
}

func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) incrementNextNonce(address ADDR, currentNonce evmtypes.Nonce, qopts ...pg.QOpt) error {
	return eb.ks.IncrementNextSequence(address, eb.chainID, currentNonce, qopts...)
}

func observeTimeUntilBroadcast[CHAIN_ID txmgrtypes.ID](chainID CHAIN_ID, createdAt, broadcastAt time.Time) {
	duration := float64(broadcastAt.Sub(createdAt))
	promTimeUntilBroadcast.WithLabelValues(chainID.String()).Observe(duration)
}
