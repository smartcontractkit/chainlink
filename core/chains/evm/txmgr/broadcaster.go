package txmgr

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/jpillora/backoff"
	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/common/client"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/label"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
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

var ErrTxRemoved = errors.New("tx removed")

//var (
//	promTimeUntilBroadcast = promauto.NewHistogramVec(prometheus.HistogramOpts{
//		Name: "tx_manager_time_until_tx_broadcast",
//		Help: "The amount of time elapsed from when a transaction is enqueued to until it is broadcast.",
//		Buckets: []float64{
//			float64(500 * time.Millisecond),
//			float64(time.Second),
//			float64(5 * time.Second),
//			float64(15 * time.Second),
//			float64(30 * time.Second),
//			float64(time.Minute),
//			float64(2 * time.Minute),
//		},
//	}, []string{"chainID"})
//)

type BroadcasterTxStore interface {
	CountUnconfirmedTransactions(context.Context, common.Address, *big.Int) (uint32, error)
	CountUnstartedTransactions(context.Context, common.Address, *big.Int) (uint32, error)
	FindNextUnstartedTransactionFromAddress(context.Context, common.Address, *big.Int) (*Tx, error)
	GetTxInProgress(context.Context, common.Address) (*Tx, error)
	SaveReplacementInProgressAttempt(context.Context, TxAttempt, *TxAttempt) error
	UpdateTxAttemptInProgressToBroadcast(context.Context, *Tx, TxAttempt, txmgrtypes.TxAttemptState) error
	UpdateTxCallbackCompleted(context.Context, uuid.UUID, *big.Int) error
	UpdateTxFatalError(context.Context, *Tx) error
	UpdateTxUnstartedToInProgress(context.Context, *Tx, *TxAttempt) error
}

type BroadcasterClient interface {
	ConfiguredChainID() *big.Int
	PendingNonceAt(context.Context, common.Address) (uint64, error)
	SendTransactionReturnCode(context.Context, *types.Transaction, common.Address) (client.SendTxReturnCode, error)
}

type BroadcasterTxAttemptBuilder interface {
	NewBumpAttempt(context.Context, Tx, TxAttempt, []TxAttempt, logger.Logger) (TxAttempt, error)
	NewAttempt(context.Context, Tx, logger.Logger, ...feetypes.Opt) (TxAttempt, error)
}

type SequenceTracker interface {
	LoadNextSequences(context.Context, []common.Address)
	GetNextSequence(context.Context, common.Address) (evmtypes.Nonce, error)
	GenerateNextSequence(common.Address, evmtypes.Nonce)
	SyncSequence(context.Context, common.Address, services.StopChan)
}

type BroadcasterKeyStore interface {
	EnabledAddressesForChain(context.Context, *big.Int) ([]common.Address, error)
}

type BroadcasterConfig struct {
	FallbackPollInterval time.Duration
	MaxInFlight          uint32
	NonceAutoSync        bool
}

// Broadcaster monitors transaction attempts for transactions that need to
// be broadcast, assigns sequences and ensures that at least one node
// somewhere has received the transaction successfully.
//
// This does not guarantee delivery! A whole host of other things can
// subsequently go wrong such as transactions being evicted from the mempool,
// nodes going offline etc. Responsibility for ensuring eventual inclusion
// into the chain falls on the shoulders of the confirmer.
//
// What Broadcaster does guarantee is:
// - a monotonic series of increasing sequences for txes that can all eventually be confirmed if you retry enough times
// - transition of txes out of unstarted into either fatal_error or unconfirmed
// - existence of a saved tx_attempt
type Broadcaster struct {
	services.StateMachine
	txAttemptBuilder BroadcasterTxAttemptBuilder
	lggr             logger.SugaredLogger
	txStore          BroadcasterTxStore
	client           BroadcasterClient
	chainID          *big.Int
	config           BroadcasterConfig

	ks               BroadcasterKeyStore
	sequenceTracker  SequenceTracker
	resumeCallback   txmgr.ResumeCallback
	enabledAddresses []common.Address

	checkerFactory TransmitCheckerFactory

	triggers map[common.Address]chan struct{}

	chStop services.StopChan
	wg     sync.WaitGroup

	initSync  sync.Mutex
	isStarted bool
}

func NewBroadcaster(
	txAttemptBuilder BroadcasterTxAttemptBuilder,
	lggr logger.Logger,
	txStore BroadcasterTxStore,
	client BroadcasterClient,
	config BroadcasterConfig,
	keystore BroadcasterKeyStore,
	sequenceTracker SequenceTracker,
	checkerFactory TransmitCheckerFactory,
) *Broadcaster {
	lggr = logger.Named(lggr, "Broadcaster")
	return &Broadcaster{
		txAttemptBuilder: txAttemptBuilder,
		lggr:             logger.Sugared(lggr),
		txStore:          txStore,
		client:           client,
		chainID:          client.ConfiguredChainID(),
		config:           config,
		ks:               keystore,
		checkerFactory:   checkerFactory,
		sequenceTracker:  sequenceTracker,
	}
}

// Start starts Broadcaster service.
// The provided context can be used to terminate Start sequence.
func (b *Broadcaster) Start(ctx context.Context) error {
	return b.StartOnce("Broadcaster", func() (err error) {
		return b.startInternal(ctx)
	})
}

// startInternal can be called multiple times, in conjunction with closeInternal. The TxMgr uses this functionality to reset broadcaster multiple times in its own lifetime.
func (b *Broadcaster) startInternal(ctx context.Context) error {
	b.initSync.Lock()
	defer b.initSync.Unlock()
	if b.isStarted {
		return errors.New("Broadcaster is already started")
	}
	var err error
	b.enabledAddresses, err = b.ks.EnabledAddressesForChain(ctx, b.chainID)
	if err != nil {
		return fmt.Errorf("Broadcaster: failed to load EnabledAddressesForChain: %w", err)
	}

	if len(b.enabledAddresses) > 0 {
		b.lggr.Debugw(fmt.Sprintf("Booting with %d keys", len(b.enabledAddresses)), "keys", b.enabledAddresses)
	} else {
		b.lggr.Warnf("Chain %s does not have any keys, no transactions will be sent on this chain", b.chainID.String())
	}
	b.chStop = make(chan struct{})
	b.wg = sync.WaitGroup{}
	b.wg.Add(len(b.enabledAddresses))
	b.triggers = make(map[common.Address]chan struct{})
	b.sequenceTracker.LoadNextSequences(ctx, b.enabledAddresses)
	for _, addr := range b.enabledAddresses {
		triggerCh := make(chan struct{}, 1)
		b.triggers[addr] = triggerCh
		go b.monitorTxs(addr, triggerCh)
	}

	b.isStarted = true
	return nil
}

// Close closes the Broadcaster
func (b *Broadcaster) Close() error {
	return b.StopOnce("Broadcaster", func() error {
		return b.closeInternal()
	})
}

func (b *Broadcaster) closeInternal() error {
	b.initSync.Lock()
	defer b.initSync.Unlock()
	if !b.isStarted {
		return fmt.Errorf("Broadcaster is not started: %w", services.ErrAlreadyStopped)
	}
	close(b.chStop)
	b.wg.Wait()
	b.isStarted = false
	return nil
}

// Trigger forces the monitor for a particular address to recheck for new txes
// Logs error and does nothing if address was not registered on startup
func (b *Broadcaster) Trigger(addr common.Address) {
	if b.isStarted {
		triggerCh, exists := b.triggers[addr]
		if !exists {
			// ignoring trigger for address which is not registered with this Broadcaster
			return
		}
		select {
		case triggerCh <- struct{}{}:
		default:
		}
	} else {
		b.lggr.Debugf("Unstarted; ignoring trigger for %s", addr)
	}
}

func (b *Broadcaster) SetResumeCallback(callback txmgr.ResumeCallback) {
	b.resumeCallback = callback
}

func newResendBackoff() backoff.Backoff {
	return backoff.Backoff{
		Min:    1 * time.Second,
		Max:    15 * time.Second,
		Jitter: true,
	}
}

func (b *Broadcaster) monitorTxs(addr common.Address, triggerCh chan struct{}) {
	defer b.wg.Done()

	ctx, cancel := b.chStop.NewCtx()
	defer cancel()

	if b.config.NonceAutoSync {
		b.lggr.Debugw("Auto-syncing sequence", "address", addr.String())
		b.sequenceTracker.SyncSequence(ctx, addr, b.chStop)
		if ctx.Err() != nil {
			return
		}
	} else {
		b.lggr.Debugw("Skipping sequence auto-sync", "address", addr.String())
	}

	var errorRetryCh <-chan time.Time
	bf := newResendBackoff()

	for {
		pollDBTimer := time.NewTimer(utils.WithJitter(b.config.FallbackPollInterval))

		retryable, err := b.ProcessUnstartedTxs(ctx, addr)
		if err != nil {
			b.lggr.Errorw("Error occurred while handling tx queue in ProcessUnstartedTxs", "err", err)
		}
		// On retryable errors we implement exponential backoff retries. This
		// handles intermittent connectivity, remote RPC races, timing issues etc
		if retryable {
			pollDBTimer.Reset(utils.WithJitter(b.config.FallbackPollInterval))
			errorRetryCh = time.After(bf.Duration())
		} else {
			bf = newResendBackoff()
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
			// tx was inserted
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

// ProcessUnstartedTxs picks up and handles all txes in the queue
// NOTE: This MUST NOT be run concurrently for the same address or it could
// result in undefined state or deadlocks.
// First handle any in_progress transactions left over from last time.
// Then keep looking up unstarted transactions and processing them until there are none remaining.
func (b *Broadcaster) ProcessUnstartedTxs(ctx context.Context, fromAddress common.Address) (retryable bool, err error) {
	retryable, err = b.handleAnyInProgressTx(ctx, fromAddress)
	if err != nil {
		return retryable, fmt.Errorf("ProcessUnstartedTxs failed on handleAnyInProgressTx: %w", err)
	}

	for {
		maxInFlightTransactions := b.config.MaxInFlight
		if maxInFlightTransactions > 0 {
			nUnconfirmed, err := b.txStore.CountUnconfirmedTransactions(ctx, fromAddress, b.chainID)
			if err != nil {
				return true, fmt.Errorf("CountUnconfirmedTransactions failed: %w", err)
			}
			if nUnconfirmed >= maxInFlightTransactions {
				nUnstarted, err := b.txStore.CountUnstartedTransactions(ctx, fromAddress, b.chainID)
				if err != nil {
					return true, fmt.Errorf("CountUnstartedTransactions failed: %w", err)
				}
				b.lggr.Warnw(fmt.Sprintf(`Transaction throttling; %d transactions in-flight and %d unstarted transactions pending (maximum number of in-flight transactions is %d per key). %s`, nUnconfirmed, nUnstarted, maxInFlightTransactions, label.MaxInFlightTransactionsWarning), "maxInFlightTransactions", maxInFlightTransactions, "nUnconfirmed", nUnconfirmed, "nUnstarted", nUnstarted)
				select {
				case <-time.After(InFlightTransactionRecheckInterval):
				case <-ctx.Done():
					return false, context.Cause(ctx)
				}
				continue
			}
		}
		etx, err := b.nextUnstartedTransactionWithSequence(fromAddress)
		if err != nil {
			return true, fmt.Errorf("processUnstartedTxs failed on nextUnstartedTransactionWithSequence: %w", err)
		}
		if etx == nil {
			return false, nil
		}

		if retryable, err := b.handleUnstartedTx(ctx, etx); err != nil {
			return retryable, fmt.Errorf("processUnstartedTxs failed on handleUnstartedTx: %w", err)
		}
	}
}

// handleInProgressTx checks if there is any transaction
// in_progress and if so, finishes the job
func (b *Broadcaster) handleAnyInProgressTx(ctx context.Context, fromAddress common.Address) (retryable bool, err error) {
	tx, err := b.txStore.GetTxInProgress(ctx, fromAddress)
	if err != nil {
		return true, fmt.Errorf("handleAnyInProgressTx failed: %w", err)
	}
	if tx != nil {
		if retryable, err := b.handleInProgressTx(ctx, *tx, tx.TxAttempts[0], tx.CreatedAt); err != nil {
			return retryable, fmt.Errorf("handleAnyInProgressTx failed: %w", err)
		}
	}
	return false, nil
}

// Finds next transaction in the queue, assigns a sequence, and moves it to "in_progress" state ready for broadcast.
// Returns nil if no transactions are in queue
func (b *Broadcaster) nextUnstartedTransactionWithSequence(fromAddress common.Address) (*Tx, error) {
	ctx, cancel := b.chStop.NewCtx()
	defer cancel()
	etx, err := b.txStore.FindNextUnstartedTransactionFromAddress(ctx, fromAddress, b.chainID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("findNextUnstartedTransactionFromAddress failed: %w", err)
	}

	sequence, err := b.sequenceTracker.GetNextSequence(ctx, etx.FromAddress)
	if err != nil {
		return nil, err
	}
	etx.Sequence = &sequence
	return etx, nil
}

func (b *Broadcaster) handleUnstartedTx(ctx context.Context, tx *Tx) (bool, error) {
	if tx.State != txmgr.TxUnstarted {
		return false, fmt.Errorf("invariant violation: expected transaction %v to be unstarted, it was %s", tx.ID, tx.State)
	}

	attempt, err := b.txAttemptBuilder.NewAttempt(ctx, *tx, b.lggr)
	if err != nil {
		return false, fmt.Errorf("processUnstartedTxs failed on NewAttempt: %w", err)
	}

	checkerSpec, err := tx.GetChecker()
	if err != nil {
		return false, fmt.Errorf("parsing transmit checker: %w", err)
	}

	checker, err := b.checkerFactory.BuildChecker(checkerSpec)
	if err != nil {
		return false, fmt.Errorf("building transmit checker: %w", err)
	}

	lgr := tx.GetLogger(b.lggr.With("fee", attempt.TxFee))

	// If the transmit check does not complete within the timeout, the transaction will be sent
	// anyway.
	// It's intentional that we only run `Check` for unstarted transactions.
	// Running it on other states might lead to nonce duplication, as we might mark applied transactions as fatally errored.

	checkCtx, cancel := context.WithTimeout(ctx, TransmitCheckTimeout)
	defer cancel()
	err = checker.Check(checkCtx, lgr, *tx, attempt)
	if errors.Is(err, context.Canceled) {
		lgr.Warn("Transmission checker timed out, sending anyway")
	} else if err != nil {
		tx.Error = null.StringFrom(err.Error())
		lgr.Warnw("Transmission checker failed, fatally erroring transaction.", "err", err)
		return true, b.saveFatallyErroredTransaction(tx)
	}
	cancel()

	if err = b.txStore.UpdateTxUnstartedToInProgress(ctx, tx, &attempt); errors.Is(err, ErrTxRemoved) {
		b.lggr.Debugw("tx removed", "txID", tx.ID, "subject", tx.Subject)
		return false, nil
	} else if err != nil {
		return true, fmt.Errorf("processUnstartedTxs failed on UpdateTxUnstartedToInProgress: %w", err)
	}

	return b.handleInProgressTx(ctx, *tx, attempt, time.Now())
}

func (b *Broadcaster) saveFatallyErroredTransaction(tx *Tx) error {
	ctx, cancel := b.chStop.NewCtx()
	defer cancel()
	if tx.State != txmgr.TxInProgress && tx.State != txmgr.TxUnstarted {
		return fmt.Errorf("can only transition to fatal_error from in_progress or unstarted, transaction is currently %s", tx.State)
	}
	if !tx.Error.Valid {
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
	if tx.PipelineTaskRunID.Valid && b.resumeCallback != nil && tx.SignalCallback {
		err := b.resumeCallback(tx.PipelineTaskRunID.UUID, nil, fmt.Errorf("fatal error while sending transaction: %s", tx.Error.String))
		if errors.Is(err, sql.ErrNoRows) {
			b.lggr.Debugw("callback missing or already resumed", "etxID", tx.ID)
		} else if err != nil {
			return fmt.Errorf("failed to resume pipeline: %w", err)
		} else {
			if err := b.txStore.UpdateTxCallbackCompleted(ctx, tx.PipelineTaskRunID.UUID, b.chainID); err != nil {
				return err
			}
		}
	}
	return b.txStore.UpdateTxFatalError(ctx, tx)
}

// There can be at most one in_progress transaction per address.
// Here we complete the job that we didn't finish last time.
func (b *Broadcaster) handleInProgressTx(ctx context.Context, tx Tx, attempt TxAttempt, initialBroadcastAt time.Time) (bool, error) {
	if tx.State != txmgr.TxInProgress {
		return false, fmt.Errorf("invariant violation: expected transaction %v to be in_progress, it was %s", tx.ID, tx.State)
	}

	signedTx, err := GetGethSignedTx(attempt.SignedRawTx)
	if err != nil {
		return false, fmt.Errorf("error while sending transaction %s (tx ID %d): %w", attempt.Hash.String(), tx.ID, err)
	}

	errType, err := b.client.SendTransactionReturnCode(ctx, signedTx, tx.FromAddress)
	if errType != client.Fatal {
		tx.InitialBroadcastAt = &initialBroadcastAt
		tx.BroadcastAt = &initialBroadcastAt
	}

	// Info log is a lightweight representation of tx and its attempt.
	// For less frequently used fields such as Metadata or SignedRawTx use debug.
	b.lggr.Infow("Sent transaction", "tx", tx.PrettyPrint(), "attempt", attempt.PrettyPrint(), "error", err)
	b.lggr.Debug("tx", tx, "attempt", attempt)

	switch errType {
	case client.Fatal:
		b.SvcErrBuffer.Append(err)
		tx.Error = null.StringFrom(err.Error())
		return true, b.saveFatallyErroredTransaction(&tx)
	case client.TransactionAlreadyKnown:
		fallthrough
	case client.Successful:
		// Scenario 1:
		// Transaction was successfully received by the RPC
		//
		// Scenario 2:
		//
		// The network has already received the transaction. This can happen for numerous reasons, for example
		// we mass transmitted transaction to multiple RPCsw or we restarted the node.
		//
		// Scenario 3:
		//
		// An external wallet has messed with the account and sent a transaction on this sequence.
		// This is not supported.
		//
		// If it turns out to have been an external wallet, we will never get a
		// receipt for this transaction and it will eventually be marked as
		// errored by the TXM.
		//observeTimeUntilBroadcast(b.chainID, tx.CreatedAt, time.Now())
		err = b.txStore.UpdateTxAttemptInProgressToBroadcast(ctx, &tx, attempt, txmgrtypes.TxAttemptBroadcast)
		if err != nil {
			return true, err
		}
		// Increment sequence if successfully broadcasted
		b.sequenceTracker.GenerateNextSequence(tx.FromAddress, *tx.Sequence)
		return false, err
	case client.Underpriced:
		return b.tryAgainWithNewEstimation(ctx, tx, attempt, initialBroadcastAt, true)
	case client.FeeOutOfValidRange:
		return b.tryAgainWithNewEstimation(ctx, tx, attempt, initialBroadcastAt, false)
	case client.InsufficientFunds:
		// Treat this error as retryable. The node will stop transmitting transactions until it gets funded since
		// this transactions is holding the nonce.
		b.SvcErrBuffer.Append(err)
		fallthrough
	case client.Retryable:
		return true, err
	case client.Unsupported:
		return false, err
	case client.ExceedsMaxFee:
		// Note that we may have broadcast to multiple nodes and had it
		// accepted by one of them! It is not guaranteed that all nodes share
		// the same tx fee cap. That is why we must treat this as an unknown
		// error that may have been confirmed.
		// If there is only one RPC node, or all RPC nodes have the same
		// configured cap, this transaction will get stuck and keep repeating
		// forever until the issue is resolved.
		b.lggr.Criticalw(`RPC node rejected this tx as outside Fee Cap`, "attempt", attempt)
		fallthrough
	default:
		// Every error that doesn't fall under one of the above categories will be treated as Unknown.
		fallthrough
	case client.Unknown:
		b.SvcErrBuffer.Append(err)
		b.lggr.Criticalw(`Unknown error occurred while handling tx queue in ProcessUnstartedTxs. This chain/RPC client may not be supported. `+
			`Urgent resolution required, Chainlink is currently operating in a degraded state and may miss transactions`, "attempt", attempt)
		nextSequence, e := b.client.PendingNonceAt(ctx, tx.FromAddress)
		if e != nil {
			err = multierr.Combine(e, err)
			return true, fmt.Errorf("failed to fetch latest pending nonce after encountering unknown RPC error while sending transaction: %w", err)
		}
		if nextSequence > (*tx.Sequence).Uint64() {
			// Despite the error, the RPC node considers the previously sent
			// transaction to have been accepted. In this case, the right thing to
			// do is assume success and hand off to Confirmer

			e = b.txStore.UpdateTxAttemptInProgressToBroadcast(ctx, &tx, attempt, txmgrtypes.TxAttemptBroadcast)
			if e != nil {
				err = multierr.Combine(e, err)
				return true, err
			}
			// Increment sequence if successfully broadcasted
			b.sequenceTracker.GenerateNextSequence(tx.FromAddress, *tx.Sequence)
			return false, err
		}
		// Either the unknown error prevented the transaction from being mined, or
		// it has not yet propagated to the mempool, or there is some race on the
		// remote RPC.
		//
		// In all cases, the best thing we can do is go into a retry loop and keep
		// trying to send the transaction over again.
		return true, fmt.Errorf("retryable error while sending transaction %s (tx ID %d): %w", attempt.Hash.String(), tx.ID, err)
	}

}

func (b *Broadcaster) tryAgainWithNewEstimation(ctx context.Context, tx Tx, attempt TxAttempt, initialBroadcastAt time.Time, bump bool) (retryable bool, err error) {
	var replacementAttempt TxAttempt
	if bump {
		replacementAttempt, err = b.txAttemptBuilder.NewBumpAttempt(ctx, tx, attempt, nil, b.lggr)
		if err != nil {
			return false, fmt.Errorf("tryAgainBumpingGas failed: %w", err)
		}
	}
	replacementAttempt, err = b.txAttemptBuilder.NewAttempt(ctx, tx, b.lggr, feetypes.OptForceRefetch)
	if err != nil {
		return false, fmt.Errorf("tryAgainWithNewEstimation failed to build new attempt: %w", err)
	}

	b.lggr.Warnw("RPC rejected transaction due to incorrect fee, re-estimated and will try again",
		"txID", tx.ID, "err", err, "newGasPrice", replacementAttempt.TxFee, "newGasLimit", replacementAttempt.ChainSpecificFeeLimit)

	return b.saveTryAgainAttempt(ctx, tx, attempt, replacementAttempt, initialBroadcastAt, replacementAttempt.TxFee, replacementAttempt.ChainSpecificFeeLimit)
}

func (b *Broadcaster) saveTryAgainAttempt(
	ctx context.Context,
	tx Tx,
	attempt TxAttempt,
	replacementAttempt TxAttempt,
	initialBroadcastAt time.Time,
	newFee gas.EvmFee,
	newFeeLimit uint64,
) (retyrable bool, err error) {
	if err = b.txStore.SaveReplacementInProgressAttempt(ctx, attempt, &replacementAttempt); err != nil {
		return true, fmt.Errorf("tryAgainWithNewFee failed: %w", err)
	}
	b.lggr.Infow("Bumped fee on initial send", "txID", tx.ID, "oldFee", attempt.TxFee.String(), "newFee", newFee.String(), "newFeeLimit", newFeeLimit)
	return b.handleInProgressTx(ctx, tx, replacementAttempt, initialBroadcastAt)
}

//func observeTimeUntilBroadcast(chainID *big.Int, createdAt, broadcastAt time.Time) {
//	duration := float64(broadcastAt.Sub(createdAt))
//	promTimeUntilBroadcast.WithLabelValues(chainID.String()).Observe(duration)
//}
