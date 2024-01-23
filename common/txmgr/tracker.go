package txmgr

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

const (
	// defaultTTL is the default time to live for abandoned transactions
	// After this TTL, the TXM stops tracking abandoned Txs.
	defaultTTL = 6 * time.Hour
	// handleTxesTimeout represents a sanity limit on how long handleTxesByState
	// should take to complete
	handleTxesTimeout = 10 * time.Minute
)

// AbandonedTx is a transaction who's 'FromAddress' was removed from the KeyStore(by the Node Operator).
// Thus, any new attempts for this Tx can't be signed/created. This means no fee bumping can be done.
// However, the Tx may still have live attempts in the chain's mempool, and could get confirmed on the
// chain as-is. Thus, the TXM should not directly discard this Tx.
type AbandonedTx[ADDR types.Hashable] struct {
	id          int64
	fromAddress ADDR
}

// Tracker tracks all transactions which have abandoned fromAddresses.
// The fromAddresses can be deleted by Node Operators from the KeyStore. In such cases,
// existing in-flight transactions for these fromAddresses are considered abandoned too.
// Since such Txs can still have attempts on chain's mempool, these could still be confirmed.
// This tracker just tracks such Txs for some time, in case they get confirmed as-is.
type Tracker[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	services.StateMachine
	txStore      txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	keyStore     txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ]
	chainID      CHAIN_ID
	lggr         logger.Logger
	enabledAddrs map[ADDR]bool
	txCache      map[int64]AbandonedTx[ADDR]
	ttl          time.Duration
	lock         sync.Mutex
	mb           *mailbox.Mailbox[int64]
	wg           sync.WaitGroup
	isStarted    bool
	ctx          context.Context
	ctxCancel    context.CancelFunc
}

func NewTracker[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
](
	txStore txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
	keyStore txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ],
	chainID CHAIN_ID,
	lggr logger.Logger,
) *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE] {
	return &Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		txStore:      txStore,
		keyStore:     keyStore,
		chainID:      chainID,
		lggr:         logger.Named(lggr, "TxMgrTracker"),
		enabledAddrs: map[ADDR]bool{},
		txCache:      map[int64]AbandonedTx[ADDR]{},
		ttl:          defaultTTL,
		mb:           mailbox.NewSingle[int64](),
		lock:         sync.Mutex{},
		wg:           sync.WaitGroup{},
	}
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Start(_ context.Context) (err error) {
	tr.lggr.Info("Abandoned transaction tracking enabled")
	return tr.StartOnce("Tracker", func() error {
		return tr.startInternal()
	})
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) startInternal() (err error) {
	tr.lock.Lock()
	defer tr.lock.Unlock()

	tr.ctx, tr.ctxCancel = context.WithCancel(context.Background())

	if err := tr.setEnabledAddresses(); err != nil {
		return fmt.Errorf("failed to set enabled addresses: %w", err)
	}
	tr.lggr.Info("Enabled addresses set")

	if err := tr.trackAbandonedTxes(tr.ctx); err != nil {
		return fmt.Errorf("failed to track abandoned txes: %w", err)
	}

	tr.isStarted = true
	if len(tr.txCache) == 0 {
		tr.lggr.Info("no abandoned txes found, skipping runLoop")
		return nil
	}

	tr.lggr.Infof("%d abandoned txes found, starting runLoop", len(tr.txCache))
	tr.wg.Add(1)
	go tr.runLoop()
	return nil
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Close() error {
	tr.lock.Lock()
	defer tr.lock.Unlock()
	return tr.StopOnce("Tracker", func() error {
		return tr.closeInternal()
	})
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) closeInternal() error {
	tr.lggr.Info("stopping tracker")
	if !tr.isStarted {
		return fmt.Errorf("tracker not started")
	}
	tr.ctxCancel()
	tr.wg.Wait()
	tr.isStarted = false
	return nil
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) runLoop() {
	defer tr.wg.Done()
	ttlExceeded := time.NewTicker(tr.ttl)
	defer ttlExceeded.Stop()
	for {
		select {
		case <-tr.mb.Notify():
			for {
				if tr.ctx.Err() != nil {
					return
				}
				blockHeight, exists := tr.mb.Retrieve()
				if !exists {
					break
				}
				if err := tr.HandleTxesByState(tr.ctx, blockHeight); err != nil {
					tr.lggr.Errorw(fmt.Errorf("failed to handle txes by state: %w", err).Error())
				}
			}
		case <-ttlExceeded.C:
			tr.lggr.Info("ttl exceeded")
			tr.MarkAllTxesFatal(tr.ctx)
			return
		case <-tr.ctx.Done():
			return
		}
	}
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) GetAbandonedAddresses() []ADDR {
	tr.lock.Lock()
	defer tr.lock.Unlock()

	if !tr.isStarted {
		return []ADDR{}
	}

	abandonedAddrs := make([]ADDR, len(tr.txCache))
	for _, atx := range tr.txCache {
		abandonedAddrs = append(abandonedAddrs, atx.fromAddress)
	}
	return abandonedAddrs
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) IsStarted() bool {
	tr.lock.Lock()
	defer tr.lock.Unlock()
	return tr.isStarted
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) setEnabledAddresses() error {
	enabledAddrs, err := tr.keyStore.EnabledAddressesForChain(tr.chainID)
	if err != nil {
		return fmt.Errorf("failed to get enabled addresses for chain: %w", err)
	}

	if len(enabledAddrs) == 0 {
		tr.lggr.Warnf("enabled address list is empty")
	}

	for _, addr := range enabledAddrs {
		tr.enabledAddrs[addr] = true
	}
	return nil
}

// trackAbandonedTxes called once to find and insert all abandoned txes into the tracker.
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) trackAbandonedTxes(ctx context.Context) (err error) {
	if tr.isStarted {
		return fmt.Errorf("tracker already started")
	}

	tr.lggr.Info("Retrieving non fatal transactions from txStore")
	nonFatalTxes, err := tr.txStore.GetNonFatalTransactions(ctx, tr.chainID)
	if err != nil {
		return fmt.Errorf("failed to get non fatal txes from txStore: %w", err)
	}

	// insert abandoned txes
	for _, tx := range nonFatalTxes {
		if !tr.enabledAddrs[tx.FromAddress] {
			tr.insertTx(tx)
		}
	}

	if err := tr.handleTxesByState(ctx, 0); err != nil {
		return fmt.Errorf("failed to handle txes by state: %w", err)
	}

	return nil
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) HandleTxesByState(ctx context.Context, blockHeight int64) error {
	tr.lock.Lock()
	defer tr.lock.Unlock()
	tr.ctx, tr.ctxCancel = context.WithTimeout(ctx, handleTxesTimeout)
	defer tr.ctxCancel()
	return tr.handleTxesByState(ctx, blockHeight)
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) handleTxesByState(ctx context.Context, blockHeight int64) error {
	tr.lggr.Info("Handling transactions by state")

	for id, atx := range tr.txCache {
		tx, err := tr.txStore.GetTxByID(ctx, atx.id)
		if err != nil {
			return fmt.Errorf("failed to get tx by ID: %w", err)
		}

		switch tx.State {
		case TxConfirmed:
			if err := tr.handleConfirmedTx(tx, blockHeight); err != nil {
				return fmt.Errorf("failed to handle confirmed txes: %w", err)
			}
		case TxConfirmedMissingReceipt, TxUnconfirmed:
			// Keep tracking tx
		case TxInProgress, TxUnstarted:
			// Tx could never be sent on chain even once. That means that we need to sign
			// an attempt to even broadcast this Tx to the chain. Since the fromAddress
			// is deleted, we can't sign it.
			errMsg := "The FromAddress for this Tx was deleted before this Tx could be broadcast to the chain."
			if err := tr.markTxFatal(ctx, tx, errMsg); err != nil {
				return fmt.Errorf("failed to mark tx as fatal: %w", err)
			}
			delete(tr.txCache, id)
		case TxFatalError:
			delete(tr.txCache, id)
		default:
			tr.lggr.Errorw(fmt.Sprintf("unhandled transaction state: %v", tx.State))
		}
	}

	return nil
}

// handleConfirmedTx removes a transaction from the tracker if it's been finalized on chain
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) handleConfirmedTx(
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	blockHeight int64,
) error {
	finalized, err := tr.txStore.IsTxFinalized(tr.ctx, blockHeight, tx.ID, tr.chainID)
	if err != nil {
		return fmt.Errorf("failed to check if tx is finalized: %w", err)
	}

	if finalized {
		delete(tr.txCache, tx.ID)
	}

	return nil
}

// insertTx inserts a transaction into the tracker as an AbandonedTx
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) insertTx(
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
	if _, contains := tr.txCache[tx.ID]; contains {
		return
	}

	tr.txCache[tx.ID] = AbandonedTx[ADDR]{
		id:          tx.ID,
		fromAddress: tx.FromAddress,
	}
	tr.lggr.Debugw(fmt.Sprintf("inserted tx %v", tx.ID))
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) markTxFatal(ctx context.Context,
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	errMsg string) error {
	tx.Error.SetValid(errMsg)

	// Set state to TxInProgress so the tracker can attempt to mark it as fatal
	tx.State = TxInProgress
	if err := tr.txStore.UpdateTxFatalError(ctx, tx); err != nil {
		return fmt.Errorf("failed to mark tx %v as abandoned: %w", tx.ID, err)
	}
	return nil
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) MarkAllTxesFatal(ctx context.Context) {
	tr.lock.Lock()
	defer tr.lock.Unlock()
	errMsg := fmt.Sprintf(
		"fromAddress for this Tx was deleted, and existing attempts onchain didn't finalize within %d hours, thus this Tx was abandoned.",
		int(tr.ttl.Hours()))

	for _, atx := range tr.txCache {
		tx, err := tr.txStore.GetTxByID(ctx, atx.id)
		if err != nil {
			tr.lggr.Errorw(fmt.Errorf("failed to get tx by ID: %w", err).Error())
			continue
		}

		if err := tr.markTxFatal(ctx, tx, errMsg); err != nil {
			tr.lggr.Errorw(fmt.Errorf("failed to mark tx as abandoned: %w", err).Error())
		}
	}
}
