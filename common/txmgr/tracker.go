package txmgr

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	// defaultTTL is the default time to live for abandoned transactions
	defaultTTL = 6 * time.Hour
)

// AbandonedTx is a transaction who's 'FromAddress' was removed from Confirmer's enabled addresses list
type AbandonedTx[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	id          int64
	fromAddress ADDR
	// fatalTime represents the time at which this transaction is to be marked fatal
	fatalTime time.Time
}

// isValid returns false when it's past fatal time for this AbandonedTx
func (atx *AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) isValid() bool {
	return time.Now().Before(atx.fatalTime)
}

// Tracker tracks and finalizes abandoned transactions
type Tracker[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	txStore         txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	lggr            logger.Logger
	enabledAddrs    map[ADDR]bool
	txCache         map[int64]AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	ttl             time.Duration
	lock            sync.Mutex
	setEnabledAddrs bool
	isTracking      bool
}

// NewTracker creates a new Tracker
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
	lggr logger.Logger,
) *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE] {
	return &Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		txStore:         txStore,
		lggr:            lggr.Named("Tracker"),
		enabledAddrs:    map[ADDR]bool{},
		txCache:         map[int64]AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{},
		ttl:             defaultTTL,
		lock:            sync.Mutex{},
		setEnabledAddrs: false,
		isTracking:      false,
	}
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SetEnabledAddresses(enabledAddrs []ADDR) error {
	tr.lock.Lock()
	defer tr.lock.Unlock()
	if tr.isTracking {
		return errors.New("cannot set enabled addresses while already tracking")
	}

	if len(enabledAddrs) == 0 {
		tr.lggr.Warnf("enabled address list is empty")
	}

	for _, addr := range enabledAddrs {
		tr.enabledAddrs[addr] = true
	}
	tr.setEnabledAddrs = true

	return nil
}

// TrackAbandonedTxes called once to find and inserts all abandoned txes into the tracker
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) TrackAbandonedTxes(ctx context.Context) (err error) {
	tr.lock.Lock()
	defer tr.lock.Unlock()
	if tr.isTracking {
		return errors.New("already tracking abandoned txes")
	}
	defer func() {
		if err == nil {
			tr.isTracking = true
		}
	}()

	if !tr.setEnabledAddrs {
		return errors.New("enabledAddresses not set")
	}

	nonFinalizedTxes, err := tr.txStore.GetNonFinalizedTransactions(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get non finalized txes from txStore")
	}

	for _, tx := range nonFinalizedTxes {
		// Check if tx is abandoned
		if tr.enabledAddrs[tx.FromAddress] {
			continue
		}
		tr.insertTx(tx)
	}

	return nil
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) IsTracking() bool {
	tr.lock.Lock()
	defer tr.lock.Unlock()
	return tr.isTracking
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Reset() {
	tr.lock.Lock()
	defer tr.lock.Unlock()
	tr.isTracking = false
	tr.setEnabledAddrs = false
	tr.txCache = map[int64]AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	tr.enabledAddrs = map[ADDR]bool{}
}

// HandleAbandonedTxes is called by the Confirmer to update abandoned transactions
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) HandleAbandonedTxes(ctx context.Context) error {
	tr.lock.Lock()
	defer tr.lock.Unlock()
	if !tr.isTracking {
		return errors.New("not isTracking abandoned txes")
	}

	for id, atx := range tr.txCache {
		if finalized := tr.finalizeTx(ctx, atx); finalized {
			delete(tr.txCache, id)
		}
	}
	return nil
}

// GetAbandonedAddresses returns list of abandoned addresses being tracked
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) GetAbandonedAddresses() []ADDR {
	tr.lock.Lock()
	defer tr.lock.Unlock()
	if !tr.isTracking {
		tr.lggr.Warn("not tracking abandoned txes")
	}

	addrs := make([]ADDR, len(tr.txCache))
	for _, atx := range tr.txCache {
		if atx.isValid() {
			addrs = append(addrs, atx.fromAddress)
		}
	}
	return addrs
}

// insertTx inserts a transaction into the tracker as an AbandonedTx
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) insertTx(
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
	if _, contains := tr.txCache[tx.ID]; contains {
		return
	}

	tr.txCache[tx.ID] = AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{
		id:          tx.ID,
		fromAddress: tx.FromAddress,
		fatalTime:   time.Now().Add(tr.ttl),
	}
	tr.lggr.Debugw(fmt.Sprintf("inserted tx %v", tx.ID))
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) getTx(
	ctx context.Context,
	atx AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) (
	*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	tx, err := tr.txStore.GetTxByID(ctx, atx.id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx by ID from txStore")
	}
	return tx, nil
}

// finalizeTx tries to finalize a transaction based on its current state.
// Transactions exceeding ttl are marked fatal.
// Returns true if the transaction was finalized.
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) finalizeTx(
	ctx context.Context, atx AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
	tx, err := tr.getTx(ctx, atx)
	if err != nil {
		tr.lggr.Errorw(errors.Wrap(err, "failed to get Tx").Error())
		return false
	}

	finalized := false
	switch tx.State {
	case TxConfirmed, TxConfirmedMissingReceipt, TxFatalError:
		finalized = true
	case TxInProgress, TxUnstarted, TxUnconfirmed:
		if atx.isValid() {
			break
		}
		if err := tr.finalizeFatal(ctx, tx); err != nil {
			tr.lggr.Errorw(err.Error())
			break
		}
		finalized = true
	default:
		tr.lggr.Panicw(fmt.Sprintf("unhandled transaction state: %v", tx.State))
	}

	return finalized
}

// finalizeFatal sets a transaction's state to fatal_error
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) finalizeFatal(
	ctx context.Context,
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	tx.Error.SetValid(fmt.Sprintf(
		"abandoned transaction exceeded time to live of %d hours", int(tr.ttl.Hours())))

	tx.State = TxInProgress
	err := tr.txStore.UpdateTxFatalError(ctx, tx)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to mark tx %v as fatal", tx.ID))
	}

	tr.lggr.Infow(fmt.Sprintf("tx %v marked fatal for exceeding ttl", tx.ID))
	return nil
}
