package txmgr

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/sqlx"
	"golang.org/x/exp/maps"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/null"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// Config encompasses config used by txmgr package
// Unless otherwise specified, these should support changing at runtime
//
//go:generate mockery --quiet --recursive --name Config --output ./mocks/ --case=underscore --structname Config --filename config.go
type Config interface {
	gas.Config
	pg.QConfig
	EthTxReaperInterval() time.Duration
	EthTxReaperThreshold() time.Duration
	EthTxResendAfterThreshold() time.Duration
	EvmGasBumpThreshold() uint64
	EvmGasBumpTxDepth() uint16
	EvmGasLimitDefault() uint32
	EvmMaxInFlightTransactions() uint32
	EvmMaxQueuedTransactions() uint64
	EvmNonceAutoSync() bool
	EvmUseForwarders() bool
	EvmRPCDefaultBatchSize() uint32
	KeySpecificMaxGasPriceWei(addr common.Address) *assets.Wei
	TriggerFallbackDBPollInterval() time.Duration
}

// For more information about the Txm architecture, see the design doc:
// https://www.notion.so/chainlink/Txm-Architecture-Overview-9dc62450cd7a443ba9e7dceffa1a8d6b

// ResumeCallback is assumed to be idempotent
type ResumeCallback func(id uuid.UUID, result interface{}, err error) error

// TxManager is the main component of the transaction manager.
// It is also the interface to external callers.
//
//go:generate mockery --quiet --recursive --name TxManager --output ./mocks/ --case=underscore --structname TxManager --filename tx_manager.go
type TxManager[ADDR types.Hashable[ADDR], TX_HASH types.Hashable[TX_HASH], BLOCK_HASH types.Hashable[BLOCK_HASH]] interface {
	txmgrtypes.HeadTrackable[*evmtypes.Head]
	services.ServiceCtx
	Trigger(addr ADDR)
	CreateEthTransaction(newTx NewTx[ADDR], qopts ...pg.QOpt) (etx txmgrtypes.Transaction, err error)
	GetForwarderForEOA(eoa ADDR) (forwarder ADDR, err error)
	RegisterResumeCallback(fn ResumeCallback)
	SendEther(chainID *big.Int, from, to ADDR, value assets.Eth, gasLimit uint32) (etx EthTx[ADDR, TX_HASH], err error)
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

type Txm[ADDR types.Hashable[ADDR], TX_HASH types.Hashable[TX_HASH], BLOCK_HASH types.Hashable[BLOCK_HASH]] struct {
	utils.StartStopOnce
	logger           logger.Logger
	txStorageService txmgrtypes.TxStore[ADDR, big.Int, TX_HASH, BLOCK_HASH, NewTx[ADDR], *evmtypes.Receipt, EthTx[ADDR, TX_HASH], EthTxAttempt[ADDR, TX_HASH], int64, int64]
	db               *sqlx.DB
	q                pg.Q
	ethClient        evmclient.Client
	config           Config
	keyStore         txmgrtypes.KeyStore[ADDR, *big.Int, int64]
	eventBroadcaster pg.EventBroadcaster
	chainID          big.Int
	checkerFactory   TransmitCheckerFactory[ADDR, TX_HASH]

	chHeads        chan *evmtypes.Head
	trigger        chan ADDR
	reset          chan reset
	resumeCallback ResumeCallback

	chStop   chan struct{}
	chSubbed chan struct{}
	wg       sync.WaitGroup

	reaper           *Reaper
	ethResender      *EthResender[ADDR, TX_HASH, BLOCK_HASH]
	ethBroadcaster   *EthBroadcaster[ADDR, TX_HASH, BLOCK_HASH]
	ethConfirmer     *EthConfirmer[ADDR, TX_HASH, BLOCK_HASH]
	fwdMgr           txmgrtypes.ForwarderManager[ADDR]
	txAttemptBuilder txmgrtypes.TxAttemptBuilder[*evmtypes.Head, gas.EvmFee, ADDR, TX_HASH, EthTx[ADDR, TX_HASH], EthTxAttempt[ADDR, TX_HASH]]
	nonceSyncer      NonceSyncer[ADDR, TX_HASH, BLOCK_HASH]
}

func (b *Txm[ADDR, TX_HASH, BLOCK_HASH]) RegisterResumeCallback(fn ResumeCallback) {
	b.resumeCallback = fn
	b.ethBroadcaster.SetResumeCallback(fn)
	b.ethConfirmer.SetResumeCallback(fn)
}

// NewTxm creates a new Txm with the given configuration.
func NewTxm(
	db *sqlx.DB,
	ethClient evmclient.Client,
	cfg Config,
	keyStore EvmKeyStore,
	eventBroadcaster pg.EventBroadcaster,
	lggr logger.Logger,
	checkerFactory EvmTransmitCheckerFactory,
	fwdMgr EvmFwdMgr,
	txAttemptBuilder EvmTxAttemptBuilder,
	txStorageService EvmTxStore,
	nonceSyncer EvmNonceSyncer,
	ethBroadcaster *EvmBroadcaster,
	ethConfirmer *EvmConfirmer,
	ethResender *EvmResender,
) *EvmTxm {
	b := EvmTxm{
		StartStopOnce:    utils.StartStopOnce{},
		logger:           lggr,
		txStorageService: txStorageService,
		db:               db,
		q:                pg.NewQ(db, lggr, cfg),
		ethClient:        ethClient,
		config:           cfg,
		keyStore:         keyStore,
		eventBroadcaster: eventBroadcaster,
		chainID:          *ethClient.ChainID(),
		checkerFactory:   checkerFactory,
		chHeads:          make(chan *evmtypes.Head),
		trigger:          make(chan *evmtypes.Address),
		chStop:           make(chan struct{}),
		chSubbed:         make(chan struct{}),
		reset:            make(chan reset),
		fwdMgr:           fwdMgr,
		txAttemptBuilder: txAttemptBuilder,
		nonceSyncer:      nonceSyncer,
		ethBroadcaster:   ethBroadcaster,
		ethConfirmer:     ethConfirmer,
		ethResender:      ethResender,
	}

	if cfg.EthTxResendAfterThreshold() <= 0 {
		b.logger.Info("EthResender: Disabled")
	}
	if cfg.EthTxReaperThreshold() > 0 && cfg.EthTxReaperInterval() > 0 {
		b.reaper = NewReaper(lggr, db, cfg, *ethClient.ChainID())
	} else {
		b.logger.Info("EthTxReaper: Disabled")
	}

	return &b
}

// Start starts Txm service.
// The provided context can be used to terminate Start sequence.
func (b *Txm[ADDR, TX_HASH, BLOCK_HASH]) Start(ctx context.Context) (merr error) {
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

		if b.ethResender != nil {
			b.ethResender.Start()
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
func (b *Txm[ADDR, TX_HASH, BLOCK_HASH]) Reset(callback func(), addr ADDR, abandon bool) (err error) {
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
func (b *Txm[ADDR, TX_HASH, BLOCK_HASH]) abandon(addr ADDR) (err error) {
	gethAddr, err := getGethAddressFromADDR(addr)
	if err != nil {
		return errors.Wrapf(err, "failed to do address format conversion")
	}
	_, err = b.q.Exec(`UPDATE eth_txes SET state='fatal_error', nonce = NULL, error = 'abandoned' WHERE state IN ('unconfirmed', 'in_progress', 'unstarted') AND evm_chain_id = $1 AND from_address = $2`, b.chainID.String(), gethAddr)
	return errors.Wrapf(err, "abandon failed to update eth_txes for key %s", addr.String())
}

func (b *Txm[ADDR, TX_HASH, BLOCK_HASH]) Close() (merr error) {
	return b.StopOnce("Txm", func() error {
		close(b.chStop)

		b.txStorageService.Close()

		if b.reaper != nil {
			b.reaper.Stop()
		}
		if b.ethResender != nil {
			b.ethResender.Stop()
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

func (b *Txm[ADDR, TX_HASH, BLOCK_HASH]) Name() string {
	return b.logger.Name()
}

func (b *Txm[ADDR, TX_HASH, BLOCK_HASH]) HealthReport() map[string]error {
	report := map[string]error{b.Name(): b.StartStopOnce.Healthy()}

	// only query if txm started properly
	b.IfStarted(func() {
		maps.Copy(report, b.ethBroadcaster.HealthReport())
		maps.Copy(report, b.ethConfirmer.HealthReport())
		maps.Copy(report, b.txAttemptBuilder.HealthReport())
	})

	if b.config.EvmUseForwarders() {
		maps.Copy(report, b.fwdMgr.HealthReport())
	}
	return report
}

func (b *Txm[ADDR, TX_HASH, BLOCK_HASH]) runLoop() {
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
			b.ethBroadcaster.Trigger(address.String())
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
			enabledAddresses, err := b.keyStore.EnabledAddressesForChain(&b.chainID)
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
func (b *Txm[ADDR, TX_HASH, BLOCK_HASH]) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
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
func (b *Txm[ADDR, TX_HASH, BLOCK_HASH]) Trigger(addr ADDR) {
	select {
	case b.trigger <- addr:
	default:
	}
}

type NewTx[ADDR types.Hashable[ADDR]] struct {
	FromAddress      ADDR
	ToAddress        ADDR
	EncodedPayload   []byte
	GasLimit         uint32
	Meta             *EthTxMeta
	ForwarderAddress ADDR

	// Pipeline variables - if you aren't calling this from ethtx task within
	// the pipeline, you don't need these variables
	MinConfirmations  null.Uint32
	PipelineTaskRunID *uuid.UUID

	Strategy txmgrtypes.TxStrategy

	// Checker defines the check that should be run before a transaction is submitted on chain.
	Checker TransmitCheckerSpec
}

// CreateEthTransaction inserts a new transaction
func (b *Txm[ADDR, TX_HASH, BLOCK_HASH]) CreateEthTransaction(newTx NewTx[ADDR], qs ...pg.QOpt) (tx txmgrtypes.Transaction, err error) {
	if err = b.checkEnabled(newTx.FromAddress); err != nil {
		return tx, err
	}

	if b.config.EvmUseForwarders() && (!newTx.ForwarderAddress.Empty()) {
		fwdPayload, fwdErr := b.fwdMgr.ConvertPayload(newTx.ToAddress, newTx.EncodedPayload)
		if fwdErr == nil {
			// Handling meta not set at caller.
			var gethToAddr common.Address
			gethToAddr, err = getGethAddressFromADDR(newTx.ToAddress)
			if err != nil {
				return tx, errors.Wrapf(err, "failed to do address format conversion")
			}

			if newTx.Meta != nil {
				newTx.Meta.FwdrDestAddress = &gethToAddr
			} else {
				newTx.Meta = &EthTxMeta{
					FwdrDestAddress: &gethToAddr,
				}
			}
			newTx.ToAddress = newTx.ForwarderAddress
			newTx.EncodedPayload = fwdPayload
		} else {
			b.logger.Errorf("Failed to use forwarder set upstream: %s", fwdErr.Error())
		}
	}

	err = b.txStorageService.CheckEthTxQueueCapacity(newTx.FromAddress, b.config.EvmMaxQueuedTransactions(), b.chainID, qs...)
	if err != nil {
		return tx, errors.Wrap(err, "Txm#CreateEthTransaction")
	}

	tx, err = b.txStorageService.CreateEthTransaction(newTx, b.chainID, qs...)
	return
}

// Calls forwarderMgr to get a proper forwarder for a given EOA.
func (b *Txm[ADDR, TX_HASH, BLOCK_HASH]) GetForwarderForEOA(eoa ADDR) (forwarder ADDR, err error) {
	if !b.config.EvmUseForwarders() {
		return forwarder, errors.Errorf("Forwarding is not enabled, to enable set EVM.Transactions.ForwardersEnabled =true")
	}
	forwarder, err = b.fwdMgr.ForwarderFor(eoa)
	return
}

func (b *Txm[ADDR, TX_HASH, BLOCK_HASH]) checkEnabled(addr ADDR) error {
	err := b.keyStore.CheckEnabled(addr, &b.chainID)
	return errors.Wrapf(err, "cannot send transaction from %s on chain ID %s", addr, b.chainID.String())
}

// SendEther creates a transaction that transfers the given value of ether
func (b *Txm[ADDR, TX_HASH, BLOCK_HASH]) SendEther(chainID *big.Int, from, to ADDR, value assets.Eth, gasLimit uint32) (etx EthTx[ADDR, TX_HASH], err error) {
	// TODO: Remove this hard-coding on evm package
	if to.Empty() {
		return etx, errors.New("cannot send ether to zero address")
	}
	etx = EthTx[ADDR, TX_HASH]{
		FromAddress:    from,
		ToAddress:      to,
		EncodedPayload: []byte{},
		Value:          value,
		GasLimit:       gasLimit,
		State:          EthTxUnstarted,
		EVMChainID:     *utils.NewBig(chainID),
	}
	err = b.txStorageService.InsertEthTx(&etx)
	return etx, errors.Wrap(err, "SendEther failed to insert eth_tx")
}

// send broadcasts the transaction to the ethereum network, writes any relevant
// data onto the attempt and returns an error (or nil) depending on the status
func sendTransaction[ADDR types.Hashable[ADDR], TX_HASH types.Hashable[TX_HASH]](ctx context.Context, ethClient evmclient.Client, a EthTxAttempt[ADDR, TX_HASH], e EthTx[ADDR, TX_HASH], logger logger.Logger) *evmclient.SendError {
	signedTx, err := a.GetSignedTx()
	if err != nil {
		return evmclient.NewFatalSendError(err)
	}

	err = ethClient.SendTransaction(ctx, signedTx)

	a.EthTx = e // for logging
	logger.Debugw("Sent transaction", "ethTxAttemptID", a.ID, "txHash", a.Hash, "err", err, "meta", e.Meta, "gasLimit", e.GasLimit, "attempt", a)
	sendErr := evmclient.NewSendError(err)
	if sendErr.IsTransactionAlreadyInMempool() {
		logger.Debugw("Transaction already in mempool", "txHash", a.Hash, "nodeErr", sendErr.Error())
		return nil
	}
	return sendErr
}

// sendEmptyTransaction sends a transaction with 0 Eth and an empty payload to the burn address
// May be useful for clearing stuck nonces
func sendEmptyTransaction[ADDR types.Hashable[ADDR], TX_HASH types.Hashable[TX_HASH]](
	ctx context.Context,
	ethClient evmclient.Client,
	txAttemptBuilder txmgrtypes.TxAttemptBuilder[*evmtypes.Head, gas.EvmFee, ADDR, TX_HASH, EthTx[ADDR, TX_HASH], EthTxAttempt[ADDR, TX_HASH]],
	nonce uint64,
	gasLimit uint32,
	gasPriceWei int64,
	fromAddress ADDR,
) (_ *gethTypes.Transaction, err error) {
	defer utils.WrapIfError(&err, "sendEmptyTransaction failed")

	attempt, err := txAttemptBuilder.NewEmptyTxAttempt(nonce, gasLimit, gas.EvmFee{Legacy: assets.NewWeiI(gasPriceWei)}, fromAddress)
	if err != nil {
		return nil, err
	}

	signedTx, err := attempt.GetSignedTx()
	if err != nil {
		return nil, err
	}

	err = ethClient.SendTransaction(ctx, signedTx)
	return signedTx, err
}

type NullTxManager[ADDR types.Hashable[ADDR], TX_HASH types.Hashable[TX_HASH], BLOCK_HASH types.Hashable[BLOCK_HASH]] struct {
	ErrMsg string
}

func (n *NullTxManager[ADDR, TX_HASH, BLOCK_HASH]) OnNewLongestChain(context.Context, *evmtypes.Head) {
}

// Start does noop for NullTxManager.
func (n *NullTxManager[ADDR, TX_HASH, BLOCK_HASH]) Start(context.Context) error { return nil }

// Close does noop for NullTxManager.
func (n *NullTxManager[ADDR, TX_HASH, BLOCK_HASH]) Close() error { return nil }

// Trigger does noop for NullTxManager.
func (n *NullTxManager[ADDR, TX_HASH, BLOCK_HASH]) Trigger(ADDR) { panic(n.ErrMsg) }
func (n *NullTxManager[ADDR, TX_HASH, BLOCK_HASH]) CreateEthTransaction(NewTx[ADDR], ...pg.QOpt) (etx txmgrtypes.Transaction, err error) {
	return etx, errors.New(n.ErrMsg)
}
func (n *NullTxManager[ADDR, TX_HASH, BLOCK_HASH]) GetForwarderForEOA(addr ADDR) (fwdr ADDR, err error) {
	return fwdr, err
}
func (n *NullTxManager[ADDR, TX_HASH, BLOCK_HASH]) Reset(f func(), addr ADDR, abandon bool) error {
	return nil
}

// SendEther does nothing, null functionality
func (n *NullTxManager[ADDR, TX_HASH, BLOCK_HASH]) SendEther(chainID *big.Int, from, to ADDR, value assets.Eth, gasLimit uint32) (etx EthTx[ADDR, TX_HASH], err error) {
	return etx, errors.New(n.ErrMsg)
}

func (n *NullTxManager[ADDR, TX_HASH, BLOCK_HASH]) Ready() error { return nil }
func (n *NullTxManager[ADDR, TX_HASH, BLOCK_HASH]) Name() string { return "NullTxManager" }
func (n *NullTxManager[ADDR, TX_HASH, BLOCK_HASH]) HealthReport() map[string]error {
	return map[string]error{}
}
func (n *NullTxManager[ADDR, TX_HASH, BLOCK_HASH]) RegisterResumeCallback(fn ResumeCallback) {}
