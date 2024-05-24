package txmgr

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/google/uuid"
	nullv4 "gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/headtracker"
	iutils "github.com/smartcontractkit/chainlink/v2/common/internal/utils"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// For more information about the Txm architecture, see the design doc:
// https://www.notion.so/chainlink/Txm-Architecture-Overview-9dc62450cd7a443ba9e7dceffa1a8d6b

// ResumeCallback is assumed to be idempotent
type ResumeCallback func(ctx context.Context, id uuid.UUID, result interface{}, err error) error

// TxManager is the main component of the transaction manager.
// It is also the interface to external callers.
//
//go:generate mockery --quiet --recursive --name TxManager --output ./mocks/ --case=underscore --structname TxManager --filename tx_manager.go
type TxManager[
	CHAIN_ID types.ID,
	HEAD types.Head[BLOCK_HASH],
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	SEQ types.Sequence,
	FEE feetypes.Fee,
] interface {
	headtracker.HeadTrackable[HEAD, BLOCK_HASH]
	services.Service
	Trigger(addr ADDR)
	CreateTransaction(ctx context.Context, txRequest txmgrtypes.TxRequest[ADDR, TX_HASH]) (etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	GetForwarderForEOA(ctx context.Context, eoa ADDR) (forwarder ADDR, err error)
	GetForwarderForEOAOCR2Feeds(ctx context.Context, eoa, ocr2AggregatorID ADDR) (forwarder ADDR, err error)
	RegisterResumeCallback(fn ResumeCallback)
	SendNativeToken(ctx context.Context, chainID CHAIN_ID, from, to ADDR, value big.Int, gasLimit uint64) (etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	Reset(addr ADDR, abandon bool) error
	// Find transactions by a field in the TxMeta blob and transaction states
	FindTxesByMetaFieldAndStates(ctx context.Context, metaField string, metaValue string, states []txmgrtypes.TxState, chainID *big.Int) (txes []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	// Find transactions with a non-null TxMeta field that was provided by transaction states
	FindTxesWithMetaFieldByStates(ctx context.Context, metaField string, states []txmgrtypes.TxState, chainID *big.Int) (txes []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	// Find transactions with a non-null TxMeta field that was provided and a receipt block number greater than or equal to the one provided
	FindTxesWithMetaFieldByReceiptBlockNum(ctx context.Context, metaField string, blockNum int64, chainID *big.Int) (txes []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	// Find transactions loaded with transaction attempts and receipts by transaction IDs and states
	FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx context.Context, ids []int64, states []txmgrtypes.TxState, chainID *big.Int) (txes []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	FindEarliestUnconfirmedBroadcastTime(ctx context.Context) (nullv4.Time, error)
	FindEarliestUnconfirmedTxAttemptBlock(ctx context.Context) (nullv4.Int, error)
	CountTransactionsByState(ctx context.Context, state txmgrtypes.TxState) (count uint32, err error)
}

type reset struct {
	// f is the function to execute between stopping/starting the
	// Broadcaster and Confirmer
	f func()
	// done is either closed after running f, or returns error if f could not
	// be run for some reason
	done chan error
}

type Txm[
	CHAIN_ID types.ID,
	HEAD types.Head[BLOCK_HASH],
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	services.StateMachine
	logger                  logger.SugaredLogger
	txStore                 txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	config                  txmgrtypes.TransactionManagerChainConfig
	txConfig                txmgrtypes.TransactionManagerTransactionsConfig
	keyStore                txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ]
	chainID                 CHAIN_ID
	checkerFactory          TransmitCheckerFactory[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	pruneQueueAndCreateLock sync.Mutex

	chHeads        chan HEAD
	trigger        chan ADDR
	reset          chan reset
	resumeCallback ResumeCallback

	chStop   services.StopChan
	chSubbed chan struct{}
	wg       sync.WaitGroup

	reaper           *Reaper[CHAIN_ID]
	resender         *Resender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	broadcaster      *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	confirmer        *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	tracker          *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	fwdMgr           txmgrtypes.ForwarderManager[ADDR]
	txAttemptBuilder txmgrtypes.TxAttemptBuilder[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) RegisterResumeCallback(fn ResumeCallback) {
	b.resumeCallback = fn
	b.broadcaster.SetResumeCallback(fn)
	b.confirmer.SetResumeCallback(fn)
}

// NewTxm creates a new Txm with the given configuration.
func NewTxm[
	CHAIN_ID types.ID,
	HEAD types.Head[BLOCK_HASH],
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
](
	chainId CHAIN_ID,
	cfg txmgrtypes.TransactionManagerChainConfig,
	txCfg txmgrtypes.TransactionManagerTransactionsConfig,
	keyStore txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ],
	lggr logger.Logger,
	checkerFactory TransmitCheckerFactory[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	fwdMgr txmgrtypes.ForwarderManager[ADDR],
	txAttemptBuilder txmgrtypes.TxAttemptBuilder[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	txStore txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
	broadcaster *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	confirmer *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
	resender *Resender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
	tracker *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
) *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE] {
	b := Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		logger:           logger.Sugared(lggr),
		txStore:          txStore,
		config:           cfg,
		txConfig:         txCfg,
		keyStore:         keyStore,
		chainID:          chainId,
		checkerFactory:   checkerFactory,
		chHeads:          make(chan HEAD),
		trigger:          make(chan ADDR),
		chStop:           make(chan struct{}),
		chSubbed:         make(chan struct{}),
		reset:            make(chan reset),
		fwdMgr:           fwdMgr,
		txAttemptBuilder: txAttemptBuilder,
		broadcaster:      broadcaster,
		confirmer:        confirmer,
		resender:         resender,
		tracker:          tracker,
	}

	if txCfg.ResendAfterThreshold() <= 0 {
		b.logger.Info("Resender: Disabled")
	}
	if txCfg.ReaperThreshold() > 0 && txCfg.ReaperInterval() > 0 {
		b.reaper = NewReaper[CHAIN_ID](lggr, b.txStore, cfg, txCfg, chainId)
	} else {
		b.logger.Info("TxReaper: Disabled")
	}

	return &b
}

// Start starts Txm service.
// The provided context can be used to terminate Start sequence.
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Start(ctx context.Context) (merr error) {
	return b.StartOnce("Txm", func() error {
		var ms services.MultiStart
		if err := ms.Start(ctx, b.broadcaster); err != nil {
			return fmt.Errorf("Txm: Broadcaster failed to start: %w", err)
		}
		if err := ms.Start(ctx, b.confirmer); err != nil {
			return fmt.Errorf("Txm: Confirmer failed to start: %w", err)
		}

		if err := ms.Start(ctx, b.txAttemptBuilder); err != nil {
			return fmt.Errorf("Txm: Estimator failed to start: %w", err)
		}

		if err := ms.Start(ctx, b.tracker); err != nil {
			return fmt.Errorf("Txm: Tracker failed to start: %w", err)
		}

		b.logger.Info("Txm starting runLoop")
		b.wg.Add(1)
		go b.runLoop()
		<-b.chSubbed

		if b.reaper != nil {
			b.reaper.Start()
		}

		if b.resender != nil {
			b.resender.Start(ctx)
		}

		if b.fwdMgr != nil {
			if err := ms.Start(ctx, b.fwdMgr); err != nil {
				return fmt.Errorf("Txm: ForwarderManager failed to start: %w", err)
			}
		}

		return nil
	})
}

// Reset stops Broadcaster/Confirmer, executes callback, then starts them again
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Reset(addr ADDR, abandon bool) (err error) {
	ok := b.IfStarted(func() {
		done := make(chan error)
		f := func() {
			if abandon {
				err = b.abandon(addr)
			}
		}

		b.reset <- reset{f, done}
		err = <-done
	})
	if !ok {
		return errors.New("not started")
	}
	return err
}

// abandon, scoped to the key of this txm:
// - marks all pending and inflight transactions fatally errored (note: at this point all transactions are either confirmed or fatally errored)
// this must not be run while Broadcaster or Confirmer are running
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) abandon(addr ADDR) (err error) {
	ctx, cancel := b.chStop.NewCtx()
	defer cancel()
	if err = b.txStore.Abandon(ctx, b.chainID, addr); err != nil {
		return fmt.Errorf("abandon failed to update txes for key %s: %w", addr.String(), err)
	}
	return nil
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Close() (merr error) {
	return b.StopOnce("Txm", func() error {
		close(b.chStop)

		b.txStore.Close()

		if b.reaper != nil {
			b.reaper.Stop()
		}
		if b.resender != nil {
			b.resender.Stop()
		}
		if b.fwdMgr != nil {
			if err := b.fwdMgr.Close(); err != nil {
				merr = errors.Join(merr, fmt.Errorf("Txm: failed to stop ForwarderManager: %w", err))
			}
		}

		b.wg.Wait()

		if err := b.txAttemptBuilder.Close(); err != nil {
			merr = errors.Join(merr, fmt.Errorf("Txm: failed to close TxAttemptBuilder: %w", err))
		}

		return nil
	})
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Name() string {
	return b.logger.Name()
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) HealthReport() map[string]error {
	report := map[string]error{b.Name(): b.Healthy()}

	// only query if txm started properly
	b.IfStarted(func() {
		services.CopyHealth(report, b.broadcaster.HealthReport())
		services.CopyHealth(report, b.confirmer.HealthReport())
		services.CopyHealth(report, b.txAttemptBuilder.HealthReport())
	})

	if b.txConfig.ForwardersEnabled() {
		services.CopyHealth(report, b.fwdMgr.HealthReport())
	}
	return report
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) runLoop() {
	ctx, cancel := b.chStop.NewCtx()
	defer cancel()

	// eb, ec and keyStates can all be modified by the runloop.
	// This is concurrent-safe because the runloop ensures serial access.
	defer b.wg.Done()
	keysChanged, unsub := b.keyStore.SubscribeToKeyChanges(ctx)
	defer unsub()

	close(b.chSubbed)

	var stopped bool
	var stopOnce sync.Once

	// execReset is defined as an inline function here because it closes over
	// eb, ec and stopped
	execReset := func(ctx context.Context, r *reset) {
		// These should always close successfully, since it should be logically
		// impossible to enter this code path with ec/eb in a state other than
		// "Started"
		if err := b.broadcaster.closeInternal(); err != nil {
			b.logger.Panicw(fmt.Sprintf("Failed to Close Broadcaster: %v", err), "err", err)
		}
		if err := b.tracker.closeInternal(); err != nil {
			b.logger.Panicw(fmt.Sprintf("Failed to Close Tracker: %v", err), "err", err)
		}
		if err := b.confirmer.closeInternal(); err != nil {
			b.logger.Panicw(fmt.Sprintf("Failed to Close Confirmer: %v", err), "err", err)
		}
		if r != nil {
			r.f()
			close(r.done)
		}
		var wg sync.WaitGroup
		// three goroutines to handle independent backoff retries starting:
		// - Broadcaster
		// - Confirmer
		// - Tracker
		// If chStop is closed, we mark stopped=true so that the main runloop
		// can check and exit early if necessary
		//
		// execReset will not return until either:
		// 1. Broadcaster, Confirmer, and Tracker all started successfully
		// 2. chStop was closed (txmgr exit)
		wg.Add(3)
		go func() {
			defer wg.Done()
			// Retry indefinitely on failure
			backoff := iutils.NewRedialBackoff()
			for {
				select {
				case <-time.After(backoff.Duration()):
					if err := b.broadcaster.startInternal(ctx); err != nil {
						b.logger.Criticalw("Failed to start Broadcaster", "err", err)
						b.SvcErrBuffer.Append(err)
						continue
					}
					return
				case <-b.chStop:
					stopOnce.Do(func() { stopped = true })
					return
				}
			}
		}()
		go func() {
			defer wg.Done()
			// Retry indefinitely on failure
			backoff := iutils.NewRedialBackoff()
			for {
				select {
				case <-time.After(backoff.Duration()):
					if err := b.tracker.startInternal(ctx); err != nil {
						b.logger.Criticalw("Failed to start Tracker", "err", err)
						b.SvcErrBuffer.Append(err)
						continue
					}
					return
				case <-b.chStop:
					stopOnce.Do(func() { stopped = true })
					return
				}
			}
		}()
		go func() {
			defer wg.Done()
			// Retry indefinitely on failure
			backoff := iutils.NewRedialBackoff()
			for {
				select {
				case <-time.After(backoff.Duration()):
					if err := b.confirmer.startInternal(ctx); err != nil {
						b.logger.Criticalw("Failed to start Confirmer", "err", err)
						b.SvcErrBuffer.Append(err)
						continue
					}
					return
				case <-b.chStop:
					stopOnce.Do(func() { stopped = true })
					return
				}
			}
		}()

		wg.Wait()
	}

	for {
		select {
		case address := <-b.trigger:
			b.broadcaster.Trigger(address)
		case head := <-b.chHeads:
			b.confirmer.mb.Deliver(head)
			b.tracker.mb.Deliver(head.BlockNumber())
		case reset := <-b.reset:
			// This check prevents the weird edge-case where you can select
			// into this block after chStop has already been closed and the
			// previous reset exited early.
			// In this case we do not want to reset again, we would rather go
			// around and hit the stop case.
			if stopped {
				reset.done <- errors.New("Txm was stopped")
				continue
			}
			execReset(ctx, &reset)
		case <-b.chStop:
			// close and exit
			//
			// Note that in some cases Broadcaster and/or Confirmer may
			// be in an Unstarted state here, if execReset exited early.
			//
			// In this case, we don't care about stopping them since they are
			// already "stopped".
			err := b.broadcaster.Close()
			if err != nil && (!errors.Is(err, services.ErrAlreadyStopped) || !errors.Is(err, services.ErrCannotStopUnstarted)) {
				b.logger.Errorw(fmt.Sprintf("Failed to Close Broadcaster: %v", err), "err", err)
			}
			err = b.confirmer.Close()
			if err != nil && (!errors.Is(err, services.ErrAlreadyStopped) || !errors.Is(err, services.ErrCannotStopUnstarted)) {
				b.logger.Errorw(fmt.Sprintf("Failed to Close Confirmer: %v", err), "err", err)
			}
			err = b.tracker.Close()
			if err != nil && (!errors.Is(err, services.ErrAlreadyStopped) || !errors.Is(err, services.ErrCannotStopUnstarted)) {
				b.logger.Errorw(fmt.Sprintf("Failed to Close Tracker: %v", err), "err", err)
			}
			return
		case <-keysChanged:
			// This check prevents the weird edge-case where you can select
			// into this block after chStop has already been closed and the
			// previous reset exited early.
			// In this case we do not want to reset again, we would rather go
			// around and hit the stop case.
			if stopped {
				continue
			}
			enabledAddresses, err := b.keyStore.EnabledAddressesForChain(ctx, b.chainID)
			if err != nil {
				b.logger.Critical("Failed to reload key states after key change")
				b.SvcErrBuffer.Append(err)
				continue
			}
			b.logger.Debugw("Keys changed, reloading", "enabledAddresses", enabledAddresses)

			execReset(ctx, nil)
		}
	}
}

// OnNewLongestChain conforms to HeadTrackable
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) OnNewLongestChain(ctx context.Context, head HEAD) {
	ok := b.IfStarted(func() {
		if b.reaper != nil {
			b.reaper.SetLatestBlockNum(head.BlockNumber())
		}
		b.txAttemptBuilder.OnNewLongestChain(ctx, head)
		select {
		case b.chHeads <- head:
		case <-ctx.Done():
			b.logger.Errorw("Timed out handling head", "blockNum", head.BlockNumber(), "ctxErr", ctx.Err())
		}
	})
	if !ok {
		b.logger.Debugw("Not started; ignoring head", "head", head, "state", b.State())
	}
}

// Trigger forces the Broadcaster to check early for the given address
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Trigger(addr ADDR) {
	select {
	case b.trigger <- addr:
	default:
	}
}

// CreateTransaction inserts a new transaction
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CreateTransaction(ctx context.Context, txRequest txmgrtypes.TxRequest[ADDR, TX_HASH]) (tx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	// Check for existing Tx with IdempotencyKey. If found, return the Tx and do nothing
	// Skipping CreateTransaction to avoid double send
	if txRequest.IdempotencyKey != nil {
		var existingTx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
		existingTx, err = b.txStore.FindTxWithIdempotencyKey(ctx, *txRequest.IdempotencyKey, b.chainID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return tx, fmt.Errorf("Failed to search for transaction with IdempotencyKey: %w", err)
		}
		if existingTx != nil {
			b.logger.Infow("Found a Tx with IdempotencyKey. Returning existing Tx without creating a new one.", "IdempotencyKey", *txRequest.IdempotencyKey)
			return *existingTx, nil
		}
	}

	if err = b.checkEnabled(ctx, txRequest.FromAddress); err != nil {
		return tx, err
	}

	if b.txConfig.ForwardersEnabled() && (!utils.IsZero(txRequest.ForwarderAddress)) {
		fwdPayload, fwdErr := b.fwdMgr.ConvertPayload(txRequest.ToAddress, txRequest.EncodedPayload)
		if fwdErr == nil {
			// Handling meta not set at caller.
			if txRequest.Meta != nil {
				txRequest.Meta.FwdrDestAddress = &txRequest.ToAddress
			} else {
				txRequest.Meta = &txmgrtypes.TxMeta[ADDR, TX_HASH]{
					FwdrDestAddress: &txRequest.ToAddress,
				}
			}
			txRequest.ToAddress = txRequest.ForwarderAddress
			txRequest.EncodedPayload = fwdPayload
		} else {
			b.logger.Errorf("Failed to use forwarder set upstream: %w", fwdErr.Error())
		}
	}

	err = b.txStore.CheckTxQueueCapacity(ctx, txRequest.FromAddress, b.txConfig.MaxQueued(), b.chainID)
	if err != nil {
		return tx, fmt.Errorf("Txm#CreateTransaction: %w", err)
	}

	tx, err = b.pruneQueueAndCreateTxn(ctx, txRequest, b.chainID)
	if err != nil {
		return tx, err
	}

	// Trigger the Broadcaster to check for new transaction
	b.broadcaster.Trigger(txRequest.FromAddress)

	return tx, nil
}

// Calls forwarderMgr to get a proper forwarder for a given EOA.
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) GetForwarderForEOA(ctx context.Context, eoa ADDR) (forwarder ADDR, err error) {
	if !b.txConfig.ForwardersEnabled() {
		return forwarder, fmt.Errorf("forwarding is not enabled, to enable set Transactions.ForwardersEnabled =true")
	}
	forwarder, err = b.fwdMgr.ForwarderFor(ctx, eoa)
	return
}

// GetForwarderForEOAOCR2Feeds calls forwarderMgr to get a proper forwarder for a given EOA and checks if its set as a transmitter on the OCR2Aggregator contract.
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) GetForwarderForEOAOCR2Feeds(ctx context.Context, eoa, ocr2Aggregator ADDR) (forwarder ADDR, err error) {
	if !b.txConfig.ForwardersEnabled() {
		return forwarder, fmt.Errorf("forwarding is not enabled, to enable set Transactions.ForwardersEnabled =true")
	}
	forwarder, err = b.fwdMgr.ForwarderForOCR2Feeds(ctx, eoa, ocr2Aggregator)
	return
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) checkEnabled(ctx context.Context, addr ADDR) error {
	if err := b.keyStore.CheckEnabled(ctx, addr, b.chainID); err != nil {
		return fmt.Errorf("cannot send transaction from %s on chain ID %s: %w", addr, b.chainID.String(), err)
	}
	return nil
}

// SendNativeToken creates a transaction that transfers the given value of native tokens
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SendNativeToken(ctx context.Context, chainID CHAIN_ID, from, to ADDR, value big.Int, gasLimit uint64) (etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	if utils.IsZero(to) {
		return etx, errors.New("cannot send native token to zero address")
	}
	txRequest := txmgrtypes.TxRequest[ADDR, TX_HASH]{
		FromAddress:    from,
		ToAddress:      to,
		EncodedPayload: []byte{},
		Value:          value,
		FeeLimit:       gasLimit,
		Strategy:       NewSendEveryStrategy(),
	}
	etx, err = b.pruneQueueAndCreateTxn(ctx, txRequest, chainID)
	if err != nil {
		return etx, fmt.Errorf("SendNativeToken failed to insert tx: %w", err)
	}

	// Trigger the Broadcaster to check for new transaction
	b.broadcaster.Trigger(from)
	return etx, nil
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxesByMetaFieldAndStates(ctx context.Context, metaField string, metaValue string, states []txmgrtypes.TxState, chainID *big.Int) (txes []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	txes, err = b.txStore.FindTxesByMetaFieldAndStates(ctx, metaField, metaValue, states, chainID)
	return
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxesWithMetaFieldByStates(ctx context.Context, metaField string, states []txmgrtypes.TxState, chainID *big.Int) (txes []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	txes, err = b.txStore.FindTxesWithMetaFieldByStates(ctx, metaField, states, chainID)
	return
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxesWithMetaFieldByReceiptBlockNum(ctx context.Context, metaField string, blockNum int64, chainID *big.Int) (txes []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	txes, err = b.txStore.FindTxesWithMetaFieldByReceiptBlockNum(ctx, metaField, blockNum, chainID)
	return
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx context.Context, ids []int64, states []txmgrtypes.TxState, chainID *big.Int) (txes []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	txes, err = b.txStore.FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx, ids, states, chainID)
	return
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindEarliestUnconfirmedBroadcastTime(ctx context.Context) (nullv4.Time, error) {
	return b.txStore.FindEarliestUnconfirmedBroadcastTime(ctx, b.chainID)
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindEarliestUnconfirmedTxAttemptBlock(ctx context.Context) (nullv4.Int, error) {
	return b.txStore.FindEarliestUnconfirmedTxAttemptBlock(ctx, b.chainID)
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CountTransactionsByState(ctx context.Context, state txmgrtypes.TxState) (count uint32, err error) {
	return b.txStore.CountTransactionsByState(ctx, state, b.chainID)
}

type NullTxManager[
	CHAIN_ID types.ID,
	HEAD types.Head[BLOCK_HASH],
	ADDR types.Hashable,
	TX_HASH, BLOCK_HASH types.Hashable,
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	ErrMsg string
}

func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) OnNewLongestChain(context.Context, HEAD) {
}

// Start does noop for NullTxManager.
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) Start(context.Context) error {
	return nil
}

// Close does noop for NullTxManager.
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) Close() error {
	return nil
}

// Trigger does noop for NullTxManager.
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) Trigger(ADDR) {
	panic(n.ErrMsg)
}
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) CreateTransaction(ctx context.Context, txRequest txmgrtypes.TxRequest[ADDR, TX_HASH]) (etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	return etx, errors.New(n.ErrMsg)
}
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) GetForwarderForEOA(ctx context.Context, addr ADDR) (fwdr ADDR, err error) {
	return fwdr, err
}
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) GetForwarderForEOAOCR2Feeds(ctx context.Context, _, _ ADDR) (fwdr ADDR, err error) {
	return fwdr, err
}

func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) Reset(addr ADDR, abandon bool) error {
	return nil
}

// SendNativeToken does nothing, null functionality
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) SendNativeToken(ctx context.Context, chainID CHAIN_ID, from, to ADDR, value big.Int, gasLimit uint64) (etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	return etx, errors.New(n.ErrMsg)
}

func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) Ready() error {
	return nil
}
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) Name() string {
	return "NullTxManager"
}
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) HealthReport() map[string]error {
	return map[string]error{}
}
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) RegisterResumeCallback(fn ResumeCallback) {
}
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) FindTxesByMetaFieldAndStates(ctx context.Context, metaField string, metaValue string, states []txmgrtypes.TxState, chainID *big.Int) (txes []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	return txes, errors.New(n.ErrMsg)
}
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) FindTxesWithMetaFieldByStates(ctx context.Context, metaField string, states []txmgrtypes.TxState, chainID *big.Int) (txes []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	return txes, errors.New(n.ErrMsg)
}
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) FindTxesWithMetaFieldByReceiptBlockNum(ctx context.Context, metaField string, blockNum int64, chainID *big.Int) (txes []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	return txes, errors.New(n.ErrMsg)
}
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx context.Context, ids []int64, states []txmgrtypes.TxState, chainID *big.Int) (txes []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	return txes, errors.New(n.ErrMsg)
}

func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) FindEarliestUnconfirmedBroadcastTime(ctx context.Context) (nullv4.Time, error) {
	return nullv4.Time{}, errors.New(n.ErrMsg)
}

func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) FindEarliestUnconfirmedTxAttemptBlock(ctx context.Context) (nullv4.Int, error) {
	return nullv4.Int{}, errors.New(n.ErrMsg)
}

func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) CountTransactionsByState(ctx context.Context, state txmgrtypes.TxState) (count uint32, err error) {
	return count, errors.New(n.ErrMsg)
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) pruneQueueAndCreateTxn(
	ctx context.Context,
	txRequest txmgrtypes.TxRequest[ADDR, TX_HASH],
	chainID CHAIN_ID,
) (
	tx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	err error,
) {
	b.pruneQueueAndCreateLock.Lock()
	defer b.pruneQueueAndCreateLock.Unlock()

	pruned, err := txRequest.Strategy.PruneQueue(ctx, b.txStore)
	if err != nil {
		return tx, err
	}
	if len(pruned) > 0 {
		b.logger.Warnw(fmt.Sprintf("Pruned %d old unstarted transactions", len(pruned)),
			"subject", txRequest.Strategy.Subject(),
			"pruned-tx-ids", pruned,
		)
	}

	tx, err = b.txStore.CreateTransaction(ctx, txRequest, chainID)
	if err != nil {
		return tx, err
	}
	b.logger.Debugw("Created transaction",
		"fromAddress", txRequest.FromAddress,
		"toAddress", txRequest.ToAddress,
		"meta", txRequest.Meta,
		"transactionID", tx.ID,
	)

	return tx, nil
}
