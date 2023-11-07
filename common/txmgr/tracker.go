package txmgr

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"slices"
	"time"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	// TTL is the default time to live for abandoned transactions (6hrs)
	TTL = 6 * time.Hour
)

// AbandonedErrorMsg occurs when an abandoned tx exceeds its time to live
var AbandonedErrorMsg = fmt.Sprintf(
	"abandoned transaction exceeded time to live of %d hours", int(TTL.Hours()))

// AbandonedTx is a transaction who's 'FromAddress' was removed from Confirmer's enabled addresses list
type AbandonedTx[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
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
	// txCache stores abandoned transactions by ID
	txCache map[int64]AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
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
	}
}

func (tracker *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SetEnabledAddresses(enabledAddrs []ADDR) {
	for _, addr := range enabledAddrs {
		tracker.enabledAddrs[addr] = true
	}
}

// TrackAbandonedTxes called once to find and inserts all abandoned txes into the tracker
func (tracker *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) TrackAbandonedTxes(ctx context.Context) {
	if tracker.enabledAddrs == nil {
		tracker.lggr.Panicw("enabledAddresses not set to track abandoned txes")
		return
	}

	nonFinalizedTxes, err := tracker.txStore.GetNonFinalizedTransactions(ctx)
	if err != nil {
		tracker.lggr.Errorw("failed to get non finalized txes from txStore")
		return
	}

	for _, tx := range nonFinalizedTxes {
		// Check if tx is abandoned
		if tracker.enabledAddrs[tx.FromAddress] {
			continue
		}

		// Ignore duplicate txes
		if _, contains := tracker.txCache[tx.ID]; contains {
			continue
		}

		tracker.insertTx(tx)
	}
}

// HandleAbandonedTxes is called by the Confirmer to update abandoned transactions
func (tracker *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) HandleAbandonedTxes(ctx context.Context) {
	for id, atx := range tracker.txCache {
		if finalized := tracker.finalizeTx(ctx, atx); finalized {
			delete(tracker.txCache, id)
		}
	}
}

// insertTx inserts a transaction into the tracker as an AbandonedTx
func (tracker *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) insertTx(
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
	if _, contains := tracker.txCache[tx.ID]; contains {
		return
	}

	tracker.txCache[tx.ID] = AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{
		tx:        tx,
		fatalTime: time.Now().Add(TTL),
	}
	tracker.lggr.Debugw(fmt.Sprintf("inserted tx %v", tx.ID))
}

// GetAbandonedAddresses returns list of abandoned addresses being tracked
func (tracker *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) GetAbandonedAddresses() []ADDR {
	var addrs []ADDR
	for _, atx := range tracker.txCache {
		if atx.isValid() && !slices.Contains(addrs, atx.tx.FromAddress) {
			addrs = append(addrs, atx.tx.FromAddress)
		}
	}
	return addrs
}

// finalizeTx tries to finalize a transaction based on its current state.
// Returns true if the transaction was finalized.
func (tracker *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) finalizeTx(
	ctx context.Context, atx AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
	// TODO: Query db again to see if tx with this ID was updated?
	switch atx.tx.State {
	case TxConfirmed, TxConfirmedMissingReceipt, TxFatalError:
		return true
	case TxInProgress:
		if atx.isValid() {
			break
		}

		// TODO: Confirm tx status on chain in case it was confirmed
		if err := tracker.finalizeFatal(ctx, atx.tx); err != nil {
			tracker.lggr.Errorw(err.Error())
			break
		}
		return true
	case TxUnstarted, TxUnconfirmed:
		// TODO: Handle TxUnstarted, TxUnconfirmed
	default:
		tracker.lggr.Panicw(fmt.Sprintf("unhandled transaction state: %v", atx.tx.State))
	}

	return false
}

// finalizeFatal sets a transaction's state to fatal_error
func (tracker *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) finalizeFatal(
	ctx context.Context,
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	tx.Error.SetValid(AbandonedErrorMsg)

	err := tracker.txStore.UpdateTxFatalError(ctx, tx)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to mark tx %v as fatal", tx.ID))
	}

	tracker.lggr.Infow(fmt.Sprintf("tx %v marked fatal for exceeding ttl", tx.ID))
	return nil
}
