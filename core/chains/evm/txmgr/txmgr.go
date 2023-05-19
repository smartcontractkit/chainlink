package txmgr

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"
	"golang.org/x/exp/maps"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
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
	CHAIN_ID txmgrtypes.ID,
	HEAD commontypes.Head[BLOCK_HASH],
	ADDR commontypes.Hashable,
	TX_HASH commontypes.Hashable,
	BLOCK_HASH commontypes.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ txmgrtypes.Sequence,
	FEE txmgrtypes.Fee,
	ADD any,
] interface {
	txmgrtypes.HeadTrackable[HEAD, BLOCK_HASH]
	services.ServiceCtx
	Trigger(addr ADDR)
	CreateEthTransaction(newTx txmgrtypes.NewTx[ADDR, TX_HASH], qopts ...pg.QOpt) (etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD], err error)
	GetForwarderForEOA(eoa ADDR) (forwarder ADDR, err error)
	RegisterResumeCallback(fn ResumeCallback)
	SendEther(chainID *big.Int, from, to ADDR, value assets.Eth, gasLimit uint32) (etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD], err error)
	Reset(f func(), addr ADDR, abandon bool) error
}

type reset struct {
	// f is the function to execute between stopping/starting the
	// EthBroadcaster and EthConfirmer
	f func()
	// done is either closed after running f, or returns error if f could not
	// be run for some reason
	done chan error
}

type Txm[
	CHAIN_ID txmgrtypes.ID,
	HEAD commontypes.Head[BLOCK_HASH],
	ADDR commontypes.Hashable,
	TX_HASH commontypes.Hashable,
	BLOCK_HASH commontypes.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ txmgrtypes.Sequence,
	FEE txmgrtypes.Fee,
	ADD any,
	FEE_UNIT txmgrtypes.Unit,
] struct {
	utils.StartStopOnce
	logger           logger.Logger
	txStore          txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]
	db               *sqlx.DB
	q                pg.Q
	config           EvmTxmConfig
	keyStore         txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ]
	eventBroadcaster pg.EventBroadcaster
	chainID          CHAIN_ID
	checkerFactory   TransmitCheckerFactory[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]

	chHeads        chan HEAD
	trigger        chan ADDR
	reset          chan reset
	resumeCallback ResumeCallback

	chStop   chan struct{}
	chSubbed chan struct{}
	wg       sync.WaitGroup

	reaper           *Reaper[CHAIN_ID]
	resender         *Resender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE, R, ADD]
	ethBroadcaster   *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]
	ethConfirmer     *EthConfirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]
	fwdMgr           txmgrtypes.ForwarderManager[ADDR]
	txAttemptBuilder txmgrtypes.TxAttemptBuilder[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]
	nonceSyncer      NonceSyncer[ADDR, TX_HASH, BLOCK_HASH]
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) RegisterResumeCallback(fn ResumeCallback) {
	b.resumeCallback = fn
	b.ethBroadcaster.SetResumeCallback(fn)
	b.ethConfirmer.SetResumeCallback(fn)
}

// NewTxm creates a new Txm with the given configuration.
func NewTxm(
	db *sqlx.DB,
	ethClient evmclient.Client,
	cfg EvmTxmConfig,
	keyStore EvmKeyStore,
	eventBroadcaster pg.EventBroadcaster,
	lggr logger.Logger,
	checkerFactory EvmTransmitCheckerFactory,
	fwdMgr EvmFwdMgr,
	txAttemptBuilder EvmTxAttemptBuilder,
	txStore EvmTxStore,
	nonceSyncer EvmNonceSyncer,
	ethBroadcaster *EvmBroadcaster,
	ethConfirmer *EvmConfirmer,
	ethResender *EvmResender,
	q pg.Q,
) *EvmTxm {
	b := EvmTxm{
		StartStopOnce:    utils.StartStopOnce{},
		logger:           lggr,
		txStore:          txStore,
		db:               db,
		q:                q,
		config:           cfg,
		keyStore:         keyStore,
		eventBroadcaster: eventBroadcaster,
		chainID:          ethClient.ConfiguredChainID(),
		checkerFactory:   checkerFactory,
		chHeads:          make(chan *evmtypes.Head),
		trigger:          make(chan common.Address),
		chStop:           make(chan struct{}),
		chSubbed:         make(chan struct{}),
		reset:            make(chan reset),
		fwdMgr:           fwdMgr,
		txAttemptBuilder: txAttemptBuilder,
		nonceSyncer:      nonceSyncer,
		ethBroadcaster:   ethBroadcaster,
		ethConfirmer:     ethConfirmer,
		resender:         ethResender,
	}

	if cfg.TxResendAfterThreshold() <= 0 {
		b.logger.Info("EthResender: Disabled")
	}
	if cfg.TxReaperThreshold() > 0 && cfg.TxReaperInterval() > 0 {
		b.reaper = NewEvmReaper(lggr, b.txStore, cfg, ethClient.ConfiguredChainID())
	} else {
		b.logger.Info("EthTxReaper: Disabled")
	}

	return &b
}

// Start starts Txm service.
// The provided context can be used to terminate Start sequence.
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) Start(ctx context.Context) (merr error) {
	return b.StartOnce("Txm", func() error {
		var ms services.MultiStart
		if err := ms.Start(ctx, b.ethBroadcaster); err != nil {
			return errors.Wrap(err, "Txm: EthBroadcaster failed to start")
		}
		if err := ms.Start(ctx, b.ethConfirmer); err != nil {
			return errors.Wrap(err, "Txm: EthConfirmer failed to start")
		}

		if err := ms.Start(ctx, b.txAttemptBuilder); err != nil {
			return errors.Wrap(err, "Txm: Estimator failed to start")
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
				return errors.Wrap(err, "Txm: EVMForwarderManager failed to start")
			}
		}

		return nil
	})
}

// Reset stops EthBroadcaster/EthConfirmer, executes callback, then starts them
// again
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) Reset(callback func(), addr ADDR, abandon bool) (err error) {
	ok := b.IfStarted(func() {
		done := make(chan error)
		f := func() {
			callback()
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
// this must not be run while EthBroadcaster or EthConfirmer are running
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) abandon(addr ADDR) (err error) {
	err = b.txStore.Abandon(b.chainID, addr)
	return errors.Wrapf(err, "abandon failed to update eth_txes for key %s", addr.String())
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) Close() (merr error) {
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
				return errors.Wrap(err, "Txm: failed to stop EVMForwarderManager")
			}
		}

		b.wg.Wait()

		b.txAttemptBuilder.Close()

		return nil
	})
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) Name() string {
	return b.logger.Name()
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) HealthReport() map[string]error {
	report := map[string]error{b.Name(): b.StartStopOnce.Healthy()}

	// only query if txm started properly
	b.IfStarted(func() {
		maps.Copy(report, b.ethBroadcaster.HealthReport())
		maps.Copy(report, b.ethConfirmer.HealthReport())
		maps.Copy(report, b.txAttemptBuilder.HealthReport())
	})

	if b.config.UseForwarders() {
		maps.Copy(report, b.fwdMgr.HealthReport())
	}
	return report
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) runLoop() {
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
		if err := b.ethBroadcaster.closeInternal(); err != nil {
			b.logger.Panicw(fmt.Sprintf("Failed to Close EthBroadcaster: %v", err), "err", err)
		}
		if err := b.ethConfirmer.closeInternal(); err != nil {
			b.logger.Panicw(fmt.Sprintf("Failed to Close EthConfirmer: %v", err), "err", err)
		}
		if r != nil {
			r.f()
			close(r.done)
		}
		var wg sync.WaitGroup
		// two goroutines to handle independent backoff retries starting:
		// - EthBroadcaster
		// - EthConfirmer
		// If chStop is closed, we mark stopped=true so that the main runloop
		// can check and exit early if necessary
		//
		// execReset will not return until either:
		// 1. Both EthBroadcaster and EthConfirmer started successfully
		// 2. chStop was closed (txmgr exit)
		wg.Add(2)
		go func() {
			defer wg.Done()
			// Retry indefinitely on failure
			backoff := utils.NewRedialBackoff()
			for {
				select {
				case <-time.After(backoff.Duration()):
					if err := b.ethBroadcaster.startInternal(); err != nil {
						b.logger.Criticalw("Failed to start EthBroadcaster", "err", err)
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
					if err := b.ethConfirmer.startInternal(); err != nil {
						b.logger.Criticalw("Failed to start EthConfirmer", "err", err)
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
			b.ethBroadcaster.Trigger(address)
		case head := <-b.chHeads:
			b.ethConfirmer.mb.Deliver(head)
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
			// Note that in some cases EthBroadcaster and/or EthConfirmer may
			// be in an Unstarted state here, if execReset exited early.
			//
			// In this case, we don't care about stopping them since they are
			// already "stopped", hence the usage of utils.EnsureClosed.
			if err := utils.EnsureClosed(b.ethBroadcaster); err != nil {
				b.logger.Panicw(fmt.Sprintf("Failed to Close EthBroadcaster: %v", err), "err", err)
			}
			if err := utils.EnsureClosed(b.ethConfirmer); err != nil {
				b.logger.Panicw(fmt.Sprintf("Failed to Close EthConfirmer: %v", err), "err", err)
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
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) OnNewLongestChain(ctx context.Context, head HEAD) {
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

// Trigger forces the EthBroadcaster to check early for the given address
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) Trigger(addr ADDR) {
	select {
	case b.trigger <- addr:
	default:
	}
}

// CreateEthTransaction inserts a new transaction
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) CreateEthTransaction(newTx txmgrtypes.NewTx[ADDR, TX_HASH], qs ...pg.QOpt) (tx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD], err error) {
	if err = b.checkEnabled(newTx.FromAddress); err != nil {
		return tx, err
	}

	if b.config.UseForwarders() && (!utils.IsZero(newTx.ForwarderAddress)) {
		fwdPayload, fwdErr := b.fwdMgr.ConvertPayload(newTx.ToAddress, newTx.EncodedPayload)
		if fwdErr == nil {
			// Handling meta not set at caller.
			if newTx.Meta != nil {
				newTx.Meta.FwdrDestAddress = &newTx.ToAddress
			} else {
				newTx.Meta = &txmgrtypes.TxMeta[ADDR, TX_HASH]{
					FwdrDestAddress: &newTx.ToAddress,
				}
			}
			newTx.ToAddress = newTx.ForwarderAddress
			newTx.EncodedPayload = fwdPayload
		} else {
			b.logger.Errorf("Failed to use forwarder set upstream: %s", fwdErr.Error())
		}
	}

	err = b.txStore.CheckEthTxQueueCapacity(newTx.FromAddress, b.config.MaxQueuedTransactions(), b.chainID, qs...)
	if err != nil {
		return tx, errors.Wrap(err, "Txm#CreateEthTransaction")
	}

	tx, err = b.txStore.CreateEthTransaction(newTx, b.chainID, qs...)
	return
}

// Calls forwarderMgr to get a proper forwarder for a given EOA.
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) GetForwarderForEOA(eoa ADDR) (forwarder ADDR, err error) {
	if !b.config.UseForwarders() {
		return forwarder, errors.Errorf("Forwarding is not enabled, to enable set EVM.Transactions.ForwardersEnabled =true")
	}
	forwarder, err = b.fwdMgr.ForwarderFor(eoa)
	return
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) checkEnabled(addr ADDR) error {
	err := b.keyStore.CheckEnabled(addr, b.chainID)
	return errors.Wrapf(err, "cannot send transaction from %s on chain ID %s", addr, b.chainID.String())
}

// SendEther creates a transaction that transfers the given value of ether
func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) SendEther(chainID CHAIN_ID, from, to ADDR, value assets.Eth, gasLimit uint32) (etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD], err error) {
	// TODO: Remove this hard-coding on evm package
	if utils.IsZero(to) {
		return etx, errors.New("cannot send ether to zero address")
	}
	etx = txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]{
		FromAddress:    from,
		ToAddress:      to,
		EncodedPayload: []byte{},
		Value:          *value.ToInt(),
		FeeLimit:       gasLimit,
		State:          EthTxUnstarted,
		ChainID:        chainID,
	}
	err = b.txStore.InsertEthTx(&etx)
	return etx, errors.Wrap(err, "SendEther failed to insert eth_tx")
}

type NullTxManager[
	CHAIN_ID txmgrtypes.ID,
	HEAD commontypes.Head[BLOCK_HASH],
	ADDR commontypes.Hashable,
	TX_HASH, BLOCK_HASH commontypes.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ txmgrtypes.Sequence,
	FEE txmgrtypes.Fee,
	ADD any,
] struct {
	ErrMsg string
}

func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) OnNewLongestChain(context.Context, *evmtypes.Head) {
}

// Start does noop for NullTxManager.
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) Start(context.Context) error {
	return nil
}

// Close does noop for NullTxManager.
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) Close() error {
	return nil
}

// Trigger does noop for NullTxManager.
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) Trigger(ADDR) {
	panic(n.ErrMsg)
}
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) CreateEthTransaction(txmgrtypes.NewTx[ADDR, TX_HASH], ...pg.QOpt) (etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD], err error) {
	return etx, errors.New(n.ErrMsg)
}
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) GetForwarderForEOA(addr ADDR) (fwdr ADDR, err error) {
	return fwdr, err
}
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) Reset(f func(), addr ADDR, abandon bool) error {
	return nil
}

// SendEther does nothing, null functionality
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) SendEther(chainID *big.Int, from, to ADDR, value assets.Eth, gasLimit uint32) (etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD], err error) {
	return etx, errors.New(n.ErrMsg)
}

func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) Ready() error {
	return nil
}
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) Name() string {
	return "NullTxManager"
}
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) HealthReport() map[string]error {
	return map[string]error{}
}
func (n *NullTxManager[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) RegisterResumeCallback(fn ResumeCallback) {
}
