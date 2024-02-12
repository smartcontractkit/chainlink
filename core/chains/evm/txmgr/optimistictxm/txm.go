package txm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	nullv4 "gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// Type aliases for EVM
type (
	Tx        = txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee]
	TxAttempt = txmgrtypes.TxAttempt[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee]
	TxRequest = txmgrtypes.TxRequest[common.Address, common.Hash]
	TxStore   = txmgrtypes.TxStore[common.Address, *big.Int, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee]
)

// ResumeCallback is assumed to be idempotent
type ResumeCallback func(id uuid.UUID, result interface{}, err error) error

//go:generate mockery --quiet --recursive --name TxManager --output ./mocks/ --case=underscore --structname TxManager --filename tx_manager.go
type TxManager interface {
	types.HeadTrackable[*evmtypes.Head, common.Hash]
	services.Service
	Trigger(addr common.Address)
	CreateTransaction(ctx context.Context, txRequest TxRequest) (etx Tx, err error)
	GetForwarderForEOA(eoa common.Address) (forwarder common.Address, err error)
	RegisterResumeCallback(fn ResumeCallback)
	SendNativeToken(ctx context.Context, chainID *big.Int, from, to common.Address, value big.Int, gasLimit uint32) (etx Tx, err error)
	Reset(addr common.Address, abandon bool) error
	// Find transactions by a field in the TxMeta blob and transaction states
	FindTxesByMetaFieldAndStates(ctx context.Context, metaField string, metaValue string, states []txmgrtypes.TxState, chainID *big.Int) (txes []*Tx, err error)
	// Find transactions with a non-null TxMeta field that was provided by transaction states
	FindTxesWithMetaFieldByStates(ctx context.Context, metaField string, states []txmgrtypes.TxState, chainID *big.Int) (txes []*Tx, err error)
	// Find transactions with a non-null TxMeta field that was provided and a receipt block number greater than or equal to the one provided
	FindTxesWithMetaFieldByReceiptBlockNum(ctx context.Context, metaField string, blockNum int64, chainID *big.Int) (txes []*Tx, err error)
	// Find transactions loaded with transaction attempts and receipts by transaction IDs and states
	FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx context.Context, ids []big.Int, states []txmgrtypes.TxState, chainID *big.Int) (txes []*Tx, err error)
	FindEarliestUnconfirmedBroadcastTime(ctx context.Context) (nullv4.Time, error)
	FindEarliestUnconfirmedTxAttemptBlock(ctx context.Context) (nullv4.Int, error)
	CountTransactionsByState(ctx context.Context, state txmgrtypes.TxState) (count uint32, err error)
}

type Txm struct {
	services.StateMachine
	logger                  logger.SugaredLogger
	txStore                 TxStore
	config                  txmgrtypes.TransactionManagerChainConfig
	txConfig                txmgrtypes.TransactionManagerTransactionsConfig
	keyStore                txmgrtypes.KeyStore[common.Address, *big.Int, evmtypes.Nonce]
	chainID                 *big.Int
	pruneQueueAndCreateLock sync.Mutex

	chHeads        chan *evmtypes.Head
	trigger        chan common.Address
	resumeCallback ResumeCallback

	chStop   services.StopChan
	chSubbed chan struct{}
	wg       sync.WaitGroup

	reaper           *Reaper
	resender         *Resender
	broadcaster      *Broadcaster
	fwdMgr           txmgrtypes.ForwarderManager[common.Address]
	txAttemptBuilder txmgrtypes.TxAttemptBuilder[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee]
	sequenceSyncer   SequenceSyncer
}

func (b *Txm) RegisterResumeCallback(fn ResumeCallback) {
	b.resumeCallback = fn
}

func NewTxm(
	chainId *big.Int,
	cfg txmgrtypes.TransactionManagerChainConfig,
	txCfg txmgrtypes.TransactionManagerTransactionsConfig,
	keyStore txmgrtypes.KeyStore[common.Address, *big.Int, evmtypes.Nonce],
	lggr logger.Logger,
	fwdMgr txmgrtypes.ForwarderManager[common.Address],
	txAttemptBuilder txmgrtypes.TxAttemptBuilder[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee],
	txStore txmgrtypes.TxStore[common.Address, *big.Int, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee],
	sequenceSyncer SequenceSyncer,
	broadcaster *Broadcaster,
	resender *Resender,
	reaper *Reaper,
) *Txm {
	b := Txm{
		logger:           logger.Sugared(lggr),
		txStore:          txStore,
		config:           cfg,
		txConfig:         txCfg,
		keyStore:         keyStore,
		chainID:          chainId,
		chHeads:          make(chan *evmtypes.Head),
		trigger:          make(chan common.Address),
		chStop:           make(chan struct{}),
		chSubbed:         make(chan struct{}),
		fwdMgr:           fwdMgr,
		txAttemptBuilder: txAttemptBuilder,
		sequenceSyncer:   sequenceSyncer,
		broadcaster:      broadcaster,
		resender:         resender,
	}

	if txCfg.ResendAfterThreshold() <= 0 {
		b.logger.Info("Resender: Disabled")
	}
	if txCfg.ReaperThreshold() > 0 && txCfg.ReaperInterval() > 0 {
		b.reaper = reaper
	} else {
		b.logger.Info("TxReaper: Disabled")
	}

	return &b
}

// Start starts Txm service.
// The provided context can be used to terminate Start sequence.
func (b *Txm) Start(ctx context.Context) (merr error) {
	return b.StartOnce("Txm", func() error {
		var ms services.MultiStart
		if err := ms.Start(ctx, b.broadcaster); err != nil {
			return fmt.Errorf("Txm: Broadcaster failed to start: %w", err)
		}

		if err := ms.Start(ctx, b.txAttemptBuilder); err != nil {
			return fmt.Errorf("Txm: Estimator failed to start: %w", err)
		}

		b.logger.Info("Txm starting runLoop")
		b.wg.Add(1)
		go b.runLoop()
		<-b.chSubbed

		if b.resender != nil {
			b.resender.Start()
		}

		if b.reaper != nil {
			b.reaper.Start()
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
func (b *Txm) Reset(addr common.Address, abandon bool) (err error) {
	ok := b.IfStarted(func() {
		done := make(chan error)
		if abandon {
			err = b.abandon(addr)
		}

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
func (b *Txm) abandon(addr common.Address) (err error) {
	ctx, cancel := b.chStop.NewCtx()
	defer cancel()
	if err = b.txStore.Abandon(ctx, b.chainID, addr); err != nil {
		return fmt.Errorf("abandon failed to update txes for key %s: %w", addr.String(), err)
	}
	return nil
}

func (b *Txm) Close() (merr error) {
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

func (b *Txm) Name() string {
	return b.logger.Name()
}

func (b *Txm) HealthReport() map[string]error {
	report := map[string]error{b.Name(): b.Healthy()}

	// only query if txm started properly
	b.IfStarted(func() {
		services.CopyHealth(report, b.broadcaster.HealthReport())
		//services.CopyHealth(report, b.confirmer.HealthReport())
		services.CopyHealth(report, b.txAttemptBuilder.HealthReport())
	})

	if b.txConfig.ForwardersEnabled() {
		services.CopyHealth(report, b.fwdMgr.HealthReport())
	}
	return report
}

func (b *Txm) runLoop() {
	// eb, ec and keyStates can all be modified by the runloop.
	// This is concurrent-safe because the runloop ensures serial access.
	defer b.wg.Done()
	keysChanged, unsub := b.keyStore.SubscribeToKeyChanges()
	defer unsub()

	close(b.chSubbed)

	for {
		select {
		case address := <-b.trigger:
			b.broadcaster.Trigger(address)
		case <-b.chHeads:
			//b.confirmer.mb.Deliver(head)
			// Tracker currently disabled for BCI-2638; refactor required
			// b.tracker.mb.Deliver(head.BlockNumber())
		case <-b.chStop:
			// close and exit

			err := b.broadcaster.Close()
			if err != nil && (!errors.Is(err, services.ErrAlreadyStopped) || !errors.Is(err, services.ErrCannotStopUnstarted)) {
				b.logger.Errorw(fmt.Sprintf("Failed to Close Broadcaster: %v", err), "err", err)
			}
			return
		case <-keysChanged:
			// This check prevents the weird edge-case where you can select
			// into this block after chStop has already been closed and the
			// previous reset exited early.
			// In this case we do not want to reset again, we would rather go
			// around and hit the stop case.
			enabledAddresses, err := b.keyStore.EnabledAddressesForChain(b.chainID)
			if err != nil {
				b.logger.Critical("Failed to reload key states after key change")
				b.SvcErrBuffer.Append(err)
				continue
			}
			b.logger.Debugw("Keys changed, reloading", "enabledAddresses", enabledAddresses)

		}
	}
}

// OnNewLongestChain conforms to HeadTrackable
func (b *Txm) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
	ok := b.IfStarted(func() {
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
func (b *Txm) Trigger(addr common.Address) {
	select {
	case b.trigger <- addr:
	default:
	}
}

// CreateTransaction inserts a new transaction
func (b *Txm) CreateTransaction(ctx context.Context, txRequest TxRequest) (tx Tx, err error) {
	// Check for existing Tx with IdempotencyKey. If found, return the Tx and do nothing
	// Skipping CreateTransaction to avoid double send
	if txRequest.IdempotencyKey != nil {
		var existingTx *Tx
		existingTx, err = b.txStore.FindTxWithIdempotencyKey(ctx, *txRequest.IdempotencyKey, b.chainID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return tx, fmt.Errorf("failed to search for transaction with IdempotencyKey: %w", err)
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
				txRequest.Meta = &txmgrtypes.TxMeta[common.Address, common.Hash]{
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
func (b *Txm) GetForwarderForEOA(eoa common.Address) (forwarder common.Address, err error) {
	if !b.txConfig.ForwardersEnabled() {
		return forwarder, fmt.Errorf("forwarding is not enabled, to enable set Transactions.ForwardersEnabled =true")
	}
	forwarder, err = b.fwdMgr.ForwarderFor(eoa)
	return
}

func (b *Txm) checkEnabled(addr common.Address) error {
	if err := b.keyStore.CheckEnabled(addr, b.chainID); err != nil {
		return fmt.Errorf("cannot send transaction from %s on chain ID %s: %w", addr, b.chainID.String(), err)
	}
	return nil
}

// SendNativeToken creates a transaction that transfers the given value of native tokens
func (b *Txm) SendNativeToken(ctx context.Context, chainID *big.Int, from, to common.Address, value big.Int, gasLimit uint32) (etx Tx, err error) {
	if utils.IsZero(to) {
		return etx, errors.New("cannot send native token to zero address")
	}
	txRequest := TxRequest{
		FromAddress:    from,
		ToAddress:      to,
		EncodedPayload: []byte{},
		Value:          value,
		FeeLimit:       gasLimit,
		Strategy:       txmgr.NewSendEveryStrategy(),
	}
	etx, err = b.pruneQueueAndCreateTxn(ctx, txRequest, chainID)
	if err != nil {
		return etx, fmt.Errorf("SendNativeToken failed to insert tx: %w", err)
	}

	// Trigger the Broadcaster to check for new transaction
	b.broadcaster.Trigger(from)
	return etx, nil
}

func (b *Txm) FindTxesByMetaFieldAndStates(ctx context.Context, metaField string, metaValue string, states []txmgrtypes.TxState, chainID *big.Int) (txes []*Tx, err error) {
	txes, err = b.txStore.FindTxesByMetaFieldAndStates(ctx, metaField, metaValue, states, chainID)
	return
}

func (b *Txm) FindTxesWithMetaFieldByStates(ctx context.Context, metaField string, states []txmgrtypes.TxState, chainID *big.Int) (txes []*Tx, err error) {
	txes, err = b.txStore.FindTxesWithMetaFieldByStates(ctx, metaField, states, chainID)
	return
}

func (b *Txm) FindTxesWithMetaFieldByReceiptBlockNum(ctx context.Context, metaField string, blockNum int64, chainID *big.Int) (txes []*Tx, err error) {
	txes, err = b.txStore.FindTxesWithMetaFieldByReceiptBlockNum(ctx, metaField, blockNum, chainID)
	return
}

func (b *Txm) FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx context.Context, ids []big.Int, states []txmgrtypes.TxState, chainID *big.Int) (txes []*Tx, err error) {
	txes, err = b.txStore.FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx, ids, states, chainID)
	return
}

func (b *Txm) FindEarliestUnconfirmedBroadcastTime(ctx context.Context) (nullv4.Time, error) {
	return b.txStore.FindEarliestUnconfirmedBroadcastTime(ctx, b.chainID)
}

func (b *Txm) FindEarliestUnconfirmedTxAttemptBlock(ctx context.Context) (nullv4.Int, error) {
	return b.txStore.FindEarliestUnconfirmedTxAttemptBlock(ctx, b.chainID)
}

func (b *Txm) CountTransactionsByState(ctx context.Context, state txmgrtypes.TxState) (count uint32, err error) {
	return b.txStore.CountTransactionsByState(ctx, state, b.chainID)
}

func (b *Txm) pruneQueueAndCreateTxn(
	ctx context.Context,
	txRequest txmgrtypes.TxRequest[common.Address, common.Hash],
	chainID *big.Int,
) (
	tx Tx,
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
