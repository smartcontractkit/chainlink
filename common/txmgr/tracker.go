package txmgr

import (
	"context"
	"fmt"
	"slices"
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
	txStore      txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	lggr         logger.Logger
	enabledAddrs map[ADDR]bool
	txCache      map[int64]AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	ttl          time.Duration
	lock         sync.Mutex
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
		txStore:      txStore,
		lggr:         lggr.Named("Tracker"),
		enabledAddrs: map[ADDR]bool{},
		txCache:      map[int64]AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{},
		ttl:          defaultTTL,
		lock:         sync.Mutex{},
	}
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SetEnabledAddresses(enabledAddrs []ADDR) {
	tr.lock.Lock()
	defer tr.lock.Unlock()
	if len(enabledAddrs) == 0 {
		tr.lggr.Warnf("enabled address list is empty")
	}
	for _, addr := range enabledAddrs {
		tr.enabledAddrs[addr] = true
	}
}

// TrackAbandonedTxes called once to find and inserts all abandoned txes into the tracker
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) TrackAbandonedTxes(ctx context.Context) error {
	tr.lock.Lock()
	defer tr.lock.Unlock()
	if len(tr.enabledAddrs) == 0 {
		return errors.New("enabledAddresses not set")
	}

	nonFinalizedTxes, err := tr.txStore.GetNonFinalizedTransactions(ctx)
	if err != nil {
		tr.lggr.Errorw("failed to get non finalized txes from txStore")
		return nil
	}

	for _, tx := range nonFinalizedTxes {
		// Check if tx is abandoned
		if tr.enabledAddrs[tx.FromAddress] {
			continue
		}

		if _, contains := tr.txCache[tx.ID]; contains {
			continue
		}

		tr.insertTx(tx)
	}

	return nil
}

// HandleAbandonedTxes is called by the Confirmer to update abandoned transactions
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) HandleAbandonedTxes(ctx context.Context) {
	tr.lock.Lock()
	defer tr.lock.Unlock()
	for id, atx := range tr.txCache {
		if finalized := tr.finalizeTx(ctx, atx); finalized {
			delete(tr.txCache, id)
		}
	}
}

// GetAbandonedAddresses returns list of abandoned addresses being tracked
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) GetAbandonedAddresses() []ADDR {
	tr.lock.Lock()
	defer tr.lock.Unlock()
	var addrs []ADDR
	for _, atx := range tr.txCache {
		if atx.isValid() && !slices.Contains(addrs, atx.fromAddress) {
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
// Returns true if the transaction was finalized.
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) finalizeTx(
	ctx context.Context, atx AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
	tx, err := tr.getTx(ctx, atx)
	if err != nil {
		tr.lggr.Errorw(err.Error())
		return false
	}

	switch tx.State {
	case TxConfirmed, TxConfirmedMissingReceipt, TxFatalError:
		return true
	case TxInProgress:
		if atx.isValid() {
			break
		}
		if err := tr.finalizeFatal(ctx, tx); err != nil {
			tr.lggr.Errorw(err.Error())
			break
		}
		return true
	case TxUnstarted:
		if err := tr.rebroadcastTx(ctx, tx); err != nil {
			tr.lggr.Errorw(err.Error())
		}
		if atx.isValid() {
			break
		}

		if err := tr.finalizeFatal(ctx, tx); err != nil {
			tr.lggr.Errorw(err.Error())
			break
		}
		return true
	case TxUnconfirmed:
		// TODO: Handle TxUnconfirmed
	default:
		tr.lggr.Panicw(fmt.Sprintf("unhandled transaction state: %v", tx.State))
	}

	return false
}

// rebroadcastTx sets a transaction's state for rebroadcasting
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) rebroadcastTx(
	ctx context.Context,
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {

	err := tr.txStore.UpdateTxUnstartedToInProgress(ctx, tx, &tx.TxAttempts[0])
	if err != nil {
		return errors.Wrap(err, "failed to rebroadcast transaction")
	}

	tr.lggr.Infow(fmt.Sprintf("tx %v set for rebroadcasting", tx.ID))
	return nil
}

// finalizeFatal sets a transaction's state to fatal_error
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) finalizeFatal(
	ctx context.Context,
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {

	tx.Error.SetValid(fmt.Sprintf(
		"abandoned transaction exceeded time to live of %d hours", int(tr.ttl.Hours())))

	err := tr.txStore.UpdateTxFatalError(ctx, tx)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to mark tx %v as fatal", tx.ID))
	}

	tr.lggr.Infow(fmt.Sprintf("tx %v marked fatal for exceeding ttl", tx.ID))
	return nil
}
