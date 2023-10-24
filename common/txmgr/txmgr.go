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
	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-relay/pkg/services"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// For more information about the Txm architecture, see the design doc:
// https://www.notion.so/chainlink/Txm-Architecture-Overview-9dc62450cd7a443ba9e7dceffa1a8d6b

// ResumeCallback is assumed to be idempotent
type ResumeCallback func(id uuid.UUID, result interface{}, err error) error

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
	types.HeadTrackable[HEAD, BLOCK_HASH]
	services.Service
	Trigger(addr ADDR)
	CreateTransaction(ctx context.Context, txRequest txmgrtypes.TxRequest[ADDR, TX_HASH]) (etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	GetForwarderForEOA(eoa ADDR) (forwarder ADDR, err error)
	RegisterResumeCallback(fn ResumeCallback)
	SendNativeToken(ctx context.Context, chainID CHAIN_ID, from, to ADDR, value big.Int, gasLimit uint32) (etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	Reset(addr ADDR, abandon bool) error
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
	utils.StartStopOnce
	logger         logger.Logger
	txStore        txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	config         txmgrtypes.TransactionManagerChainConfig
	txConfig       txmgrtypes.TransactionManagerTransactionsConfig
	keyStore       txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ]
	chainID        CHAIN_ID
	checkerFactory TransmitCheckerFactory[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]

	chHeads        chan HEAD
	trigger        chan ADDR
	reset          chan reset
	resumeCallback ResumeCallback

	chStop   chan struct{}
	chSubbed chan struct{}
	wg       sync.WaitGroup

	reaper           *Reaper[CHAIN_ID]
	resender         *Resender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	broadcaster      *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	confirmer        *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	fwdMgr           txmgrtypes.ForwarderManager[ADDR]
	txAttemptBuilder txmgrtypes.TxAttemptBuilder[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	sequenceSyncer   SequenceSyncer[ADDR, TX_HASH, BLOCK_HASH, SEQ]
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
	sequenceSyncer SequenceSyncer[ADDR, TX_HASH, BLOCK_HASH, SEQ],
	broadcaster *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	confirmer *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
	resender *Resender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE] {
	b := Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		logger:           lggr,
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
		sequenceSyncer:   sequenceSyncer,
		broadcaster:      broadcaster,
		confirmer:        confirmer,
		resender:         resender,
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
			return pkgerrors.Wrap(err, "Txm: Broadcaster failed to start")
		}
		if err := ms.Start(ctx, b.confirmer); err != nil {
			return pkgerrors.Wrap(err, "Txm: Confirmer failed to start")
		}

		if err := ms.Start(ctx, b.txAttemptBuilder); err != nil {
			return pkgerrors.Wrap(err, "Txm: Estimator failed to start")
		}

		b.wg.Add(1)
		go b.runLoop()
		<-b.chSubbed

		if b.reaper != nil {
			b.reaper.Start()
		}

		if b.resender != nil {
			b.resender.Start()
		}

		if b.fwdMgr != nil {
			if err := ms.Start(ctx, b.fwdMgr); err != nil {
				return pkgerrors.Wrap(err, "Txm: ForwarderManager failed to start")
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
	ctx, cancel := utils.StopChan(b.chStop).NewCtx()
	defer cancel()
	err = b.txStore.Abandon(ctx, b.chainID, addr)
	return pkgerrors.Wrapf(err, "abandon failed to update txes for key %s", addr.String())
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
				merr = errors.Join(merr, pkgerrors.Wrap(err, "Txm: failed to stop ForwarderManager"))
			}
		}

		b.wg.Wait()

		if err := b.txAttemptBuilder.Close(); err != nil {
			merr = errors.Join(merr, pkgerrors.Wrap(err, "Txm: failed to close TxAttemptBuilder"))
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
	// eb, ec and keyStates can all be modified by the runloop.
	// This is concurrent-safe because the runloop ensures serial access.
	defer b.wg.Done()
	keysChanged, unsub := b.keyStore.SubscribeToKeyChanges()
	defer unsub()

	close(b.chSubbed)

	var stopped bool
	var stopOnce sync.Once

	// execReset is defined as an inline function here because it closes over
	// eb, ec and stopped
	execReset := func(r *reset) {
		// These should always close successfully, since it should be logically
		// impossible to enter this code path with ec/eb in a state other than
		// "Started"
		if err := b.broadcaster.closeInternal(); err != nil {
			b.logger.Panicw(fmt.Sprintf("Failed to Close Broadcaster: %v", err), "err", err)
		}
		if err := b.confirmer.closeInternal(); err != nil {
			b.logger.Panicw(fmt.Sprintf("Failed to Close Confirmer: %v", err), "err", err)
		}
		if r != nil {
			r.f()
			close(r.done)
		}
		var wg sync.WaitGroup
		// two goroutines to handle independent backoff retries starting:
		// - Broadcaster
		// - Confirmer
		// If chStop is closed, we mark stopped=true so that the main runloop
		// can check and exit early if necessary
		//
		// execReset will not return until either:
		// 1. Both Broadcaster and Confirmer started successfully
		// 2. chStop was closed (txmgr exit)
		wg.Add(2)
		go func() {
			defer wg.Done()
			// Retry indefinitely on failure
			backoff := utils.NewRedialBackoff()
			for {
				select {
				case <-time.After(backoff.Duration()):
					if err := b.broadcaster.startInternal(); err != nil {
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
			backoff := utils.NewRedialBackoff()
			for {
				select {
				case <-time.After(backoff.Duration()):
					if err := b.confirmer.startInternal(); err != nil {
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
			execReset(&reset)
		case <-b.chStop:
			// close and exit
			//
			// Note that in some cases Broadcaster and/or Confirmer may
			// be in an Unstarted state here, if execReset exited early.
			//
			// In this case, we don't care about stopping them since they are
			// already "stopped", hence the usage of utils.EnsureClosed.
			if err := utils.EnsureClosed(b.broadcaster); err != nil {
				b.logger.Panicw(fmt.Sprintf("Failed to Close Broadcaster: %v", err), "err", err)
			}
			if err := utils.EnsureClosed(b.confirmer); err != nil {
				b.logger.Panicw(fmt.Sprintf("Failed to Close Confirmer: %v", err), "err", err)
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
			enabledAddresses, err := b.keyStore.EnabledAddressesForChain(b.chainID)
			if err != nil {
				b.logger.Criticalf("Failed to reload key states after key change")
				b.SvcErrBuffer.Append(err)
				continue
			}
			b.logger.Debugw("Keys changed, reloading", "enabledAddresses", enabledAddresses)

			execReset(nil)
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
			return tx, pkgerrors.Wrap(err, "Failed to search for transaction with IdempotencyKey")
		}
		if existingTx != nil {
			b.logger.Infow("Found a Tx with IdempotencyKey. Returning existing Tx without creating a new one.", "IdempotencyKey", *txRequest.IdempotencyKey)
			return *existingTx, nil
		}
	}

	if err = b.checkEnabled(txRequest.FromAddress); err != nil {
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
			b.logger.Errorf("Failed to use forwarder set upstream: %s", fwdErr.Error())
		}
	}

	err = b.txStore.CheckTxQueueCapacity(ctx, txRequest.FromAddress, b.txConfig.MaxQueued(), b.chainID)
	if err != nil {
		return tx, pkgerrors.Wrap(err, "Txm#CreateTransaction")
	}

	tx, err = b.txStore.CreateTransaction(ctx, txRequest, b.chainID)
	return
}

// Calls forwarderMgr to get a proper forwarder for a given EOA.
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) GetForwarderForEOA(eoa ADDR) (forwarder ADDR, err error) {
	if !b.txConfig.ForwardersEnabled() {
		return forwarder, pkgerrors.Errorf("Forwarding is not enabled, to enable set Transactions.ForwardersEnabled =true")
	}
	forwarder, err = b.fwdMgr.ForwarderFor(eoa)
	return
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) checkEnabled(addr ADDR) error {
	err := b.keyStore.CheckEnabled(addr, b.chainID)
	return pkgerrors.Wrapf(err, "cannot send transaction from %s on chain ID %s", addr, b.chainID.String())
}

// SendNativeToken creates a transaction that transfers the given value of native tokens
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SendNativeToken(ctx context.Context, chainID CHAIN_ID, from, to ADDR, value big.Int, gasLimit uint32) (etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
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
	etx, err = b.txStore.CreateTransaction(ctx, txRequest, chainID)
	return etx, pkgerrors.Wrap(err, "SendNativeToken failed to insert tx")
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
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) GetForwarderForEOA(addr ADDR) (fwdr ADDR, err error) {
	return fwdr, err
}
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) Reset(addr ADDR, abandon bool) error {
	return nil
}

// SendNativeToken does nothing, null functionality
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) SendNativeToken(ctx context.Context, chainID CHAIN_ID, from, to ADDR, value big.Int, gasLimit uint32) (etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
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
