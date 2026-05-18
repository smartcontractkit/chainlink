package txmgr

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
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
	// batchSize is the number of txes to fetch from the txStore at once
	batchSize = 1000
)

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
	txStore  txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	keyStore txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ]
	chainID  CHAIN_ID
	lggr     logger.Logger

	lock         sync.Mutex
	enabledAddrs map[ADDR]bool
	txCache      map[int64]ADDR // cache tx fromAddress by txID

	ttl time.Duration
	mb  *mailbox.Mailbox[int64]

	initSync  sync.Mutex
	wg        sync.WaitGroup
	chStop    services.StopChan
	isStarted bool
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
		txCache:      map[int64]ADDR{},
		ttl:          defaultTTL,
		mb:           mailbox.NewSingle[int64](),
		lock:         sync.Mutex{},
		wg:           sync.WaitGroup{},
	}
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Start(ctx context.Context) (err error) {
	tr.lggr.Info("Abandoned transaction tracking enabled")
	return tr.StartOnce("Tracker", func() error {
		return tr.startInternal(ctx)
	})
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) startInternal(ctx context.Context) (err error) {
	tr.initSync.Lock()
	defer tr.initSync.Unlock()

	if err := tr.setEnabledAddresses(ctx); err != nil {
		return fmt.Errorf("failed to set enabled addresses: %w", err)
	}
	tr.lggr.Infof("enabled addresses set for chainID %v", tr.chainID)

	tr.chStop = make(chan struct{})
	tr.wg.Add(1)
	go tr.runLoop(tr.chStop.NewCtx())
	tr.isStarted = true
	return nil
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Close() error {
	return tr.StopOnce("Tracker", func() error {
		return tr.closeInternal()
	})
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) closeInternal() error {
	tr.initSync.Lock()
	defer tr.initSync.Unlock()

	tr.lggr.Info("stopping tracker")
	if !tr.isStarted {
		return fmt.Errorf("tracker is not started: %w", services.ErrAlreadyStopped)
	}

	close(tr.chStop)
	tr.wg.Wait()
	tr.isStarted = false
	return nil
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) runLoop(ctx context.Context, cancel context.CancelFunc) {
	defer tr.wg.Done()
	defer cancel()

	if err := tr.trackAbandonedTxes(ctx); err != nil {
		tr.lggr.Errorf("failed to track abandoned txes: %v", err)
		return
	}
	if err := tr.handleTxesByState(ctx); err != nil {
		tr.lggr.Errorf("failed to handle txes by state: %v", err)
		return
	}
	if tr.AbandonedTxCount() == 0 {
		tr.lggr.Info("no abandoned txes found, skipping runLoop")
		return
	}
	tr.lggr.Infof("%d abandoned txes found, starting runLoop", tr.AbandonedTxCount())

	ttlExceeded := time.NewTicker(tr.ttl)
	defer ttlExceeded.Stop()
	for {
		select {
		case <-tr.mb.Notify():
			for {
				blockHeight := tr.mb.RetrieveLatestAndClear()
				if blockHeight == 0 {
					break
				}
				if err := tr.handleTxesByState(ctx); err != nil {
					tr.lggr.Errorf("failed to handle txes by state: %v", err)
					return
				}
				if tr.AbandonedTxCount() == 0 {
					tr.lggr.Info("all abandoned txes handled, stopping runLoop")
					return
				}
			}
		case <-ttlExceeded.C:
			tr.lggr.Info("ttl exceeded")
			tr.markAllTxesFatal(ctx)
			return
		case <-ctx.Done():
			return
		}
	}
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) GetAbandonedAddresses() []ADDR {
	tr.lock.Lock()
	defer tr.lock.Unlock()

	abandonedAddrs := make([]ADDR, len(tr.txCache))
	for _, fromAddress := range tr.txCache {
		abandonedAddrs = append(abandonedAddrs, fromAddress)
	}
	return abandonedAddrs
}

// AbandonedTxCount returns the number of abandoned txes currently being tracked
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) AbandonedTxCount() int {
	tr.lock.Lock()
	defer tr.lock.Unlock()
	return len(tr.txCache)
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) IsStarted() bool {
	tr.initSync.Lock()
	defer tr.initSync.Unlock()
	return tr.isStarted
}

// setEnabledAddresses is called on startup to set the enabled addresses for the chain
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) setEnabledAddresses(ctx context.Context) error {
	tr.lock.Lock()
	defer tr.lock.Unlock()

	enabledAddrs, err := tr.keyStore.EnabledAddressesForChain(ctx, tr.chainID)
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

// trackAbandonedTxes called on startup to find and insert all abandoned txes into the tracker.
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) trackAbandonedTxes(ctx context.Context) (err error) {
	return sqlutil.Batch(func(offset, limit uint) (count uint, err error) {
		var enabledAddrs []ADDR
		for addr := range tr.enabledAddrs {
			enabledAddrs = append(enabledAddrs, addr)
		}

		nonFatalTxes, err := tr.txStore.GetAbandonedTransactionsByBatch(ctx, tr.chainID, enabledAddrs, offset, limit)
		if err != nil {
			return 0, fmt.Errorf("failed to get non fatal txes from txStore: %w", err)
		}
		// insert abandoned txes
		tr.lock.Lock()
		for _, tx := range nonFatalTxes {
			if !tr.enabledAddrs[tx.FromAddress] {
				tr.txCache[tx.ID] = tx.FromAddress
				tr.lggr.Debugf("inserted tx %v", tx.ID)
			}
		}
		tr.lock.Unlock()
		return uint(len(nonFatalTxes)), nil
	}, batchSize)
}

// handleTxesByState handles all txes in the txCache by their state
// It's called on every new blockHeight and also on startup to handle all txes in the txCache
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) handleTxesByState(ctx context.Context) error {
	tr.lock.Lock()
	defer tr.lock.Unlock()
	ctx, cancel := context.WithTimeout(ctx, handleTxesTimeout)
	defer cancel()

	for id := range tr.txCache {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		tx, err := tr.txStore.GetTxByID(ctx, id)
		if err != nil {
			tr.lggr.Errorf("failed to get tx by ID: %v", err)
			continue
		}
		if tx == nil {
			tr.lggr.Warnf("tx with ID %v no longer exists, removing from tracker", id)
			delete(tr.txCache, id)
			continue
		}

		switch tx.State {
		case TxConfirmed:
			// TODO: Handle finalized state https://smartcontract-it.atlassian.net/browse/BCI-2920
		case TxConfirmedMissingReceipt, TxUnconfirmed:
			// Keep tracking tx
		case TxInProgress, TxUnstarted:
			// Tx could never be sent on chain even once. That means that we need to sign
			// an attempt to even broadcast this Tx to the chain. Since the fromAddress
			// is deleted, we can't sign it.
			errMsg := "The FromAddress for this Tx was deleted before this Tx could be broadcast to the chain."
			if err := tr.markTxFatal(ctx, tx, errMsg); err != nil {
				tr.lggr.Errorf("failed to mark tx as fatal: %v", err)
				continue
			}
			delete(tr.txCache, id)
		case TxFatalError:
			delete(tr.txCache, id)
		default:
			tr.lggr.Errorf("unhandled transaction state: %v", tx.State)
		}
	}

	return nil
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

// markAllTxesFatal tries to mark all txes in the txCache as fatal and removes them from the cache
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) markAllTxesFatal(ctx context.Context) {
	tr.lock.Lock()
	defer tr.lock.Unlock()

	errMsg := fmt.Sprintf(
		"tx abandoned: fromAddress for this tx was deleted and existing attempts didn't finalize onchain within %d hours",
		int(tr.ttl.Hours()))

	for id := range tr.txCache {
		tx, err := tr.txStore.GetTxByID(ctx, id)
		if err != nil {
			tr.lggr.Errorf("failed to get tx by ID: %v", err)
			delete(tr.txCache, id)
			continue
		}

		if err := tr.markTxFatal(ctx, tx, errMsg); err != nil {
			tr.lggr.Errorf("failed to mark tx as abandoned: %v", err)
		}
		delete(tr.txCache, id)
	}
}
