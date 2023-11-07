package txmgr

import (
	"context"
	"fmt"
	"slices"
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
	// ttl is the default time to live for abandoned transactions
	ttl time.Duration
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
	}
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SetEnabledAddresses(enabledAddrs []ADDR) {
	for _, addr := range enabledAddrs {
		tr.enabledAddrs[addr] = true
	}
}

// TrackAbandonedTxes called once to find and inserts all abandoned txes into the tracker
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) TrackAbandonedTxes(ctx context.Context) error {
	if len(tr.enabledAddrs) == 0 {
		tr.lggr.Errorw("enabledAddresses not set to track abandoned txes")
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
	for id, atx := range tr.txCache {
		if finalized := tr.finalizeTx(ctx, atx); finalized {
			delete(tr.txCache, id)
		}
	}
}

// insertTx inserts a transaction into the tracker as an AbandonedTx
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) insertTx(
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
	if _, contains := tr.txCache[tx.ID]; contains {
		return
	}

	tr.txCache[tx.ID] = AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{
		tx:        tx,
		fatalTime: time.Now().Add(tr.ttl),
	}
	tr.lggr.Debugw(fmt.Sprintf("inserted tx %v", tx.ID))
}

// GetAbandonedAddresses returns list of abandoned addresses being tracked
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) GetAbandonedAddresses() []ADDR {
	var addrs []ADDR
	for _, atx := range tr.txCache {
		if atx.isValid() && !slices.Contains(addrs, atx.tx.FromAddress) {
			addrs = append(addrs, atx.tx.FromAddress)
		}
	}
	return addrs
}

// finalizeTx tries to finalize a transaction based on its current state.
// Returns true if the transaction was finalized.
func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) finalizeTx(
	ctx context.Context, atx AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
	// TODO: Query db to update state of transaction

	switch atx.tx.State {
	case TxConfirmed, TxConfirmedMissingReceipt, TxFatalError:
		return true
	case TxInProgress:
		if atx.isValid() {
			break
		}

		if err := tr.finalizeFatal(ctx, atx.tx); err != nil {
			tr.lggr.Errorw(err.Error())
			break
		}
		return true
	case TxUnstarted:
		if err := tr.rebroadcastTx(ctx, atx.tx); err != nil {
			tr.lggr.Errorw(err.Error())
		}
	case TxUnconfirmed:
		// TODO: Handle TxUnconfirmed
	default:
		tr.lggr.Panicw(fmt.Sprintf("unhandled transaction state: %v", atx.tx.State))
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
