package txmgr

import (
	"context"
	"fmt"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"time"
)

// TODO: Track abandoned transactions
// 1: Get abandoned transaction addresses from confirmer
// 2: Add all abandoned addresses to AbandonedTracker
// 3: Loop to confirm / invalidate transactions
// TODO: Thread safety?

const (
	// TTL is the default time to live for abandoned transactions (6hrs)
	TTL = 6 * time.Hour
)

// AbandonedTx is a transaction who's fromAddress was removed from the Confirmer's enabledAddresses list
type AbandonedTx[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	// invalidTime represents the time at which this transaction is to be marked fatal
	invalidTime time.Time
}

// AbandonedTracker tracks all abandoned transactions
type AbandonedTracker[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	txStore *txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	lggr    logger.Logger
	txes    []AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
}

func NewAbandonedTracker[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
](
	txStore *txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
	lggr logger.Logger,
) AbandonedTracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE] {
	return AbandonedTracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		txStore: txStore,
		lggr:    lggr.Named("Abandoned Tracker"),
		txes:    make([]AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], 0),
	}
}

func (tracker *AbandonedTracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) InsertAbandonedTx(
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
	tracker.txes = append(tracker.txes, AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{
		tx:          tx,
		invalidTime: time.Now().Add(TTL),
	})
}

// MarkFatal sets an abandoned transaction's state to fatal_error
// TODO: Add context
func (tracker *AbandonedTracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) MarkFatal(
	atx AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
	atx.tx.Error.SetValid(fmt.Sprintf(
		"abandoned transaction exceeded time to live of %d hours", int(TTL.Hours())))
	err := (*tracker.txStore).UpdateTxFatalError(context.Background(), atx.tx)
	if err != nil {
		// TODO: Handle error
	}
}

// RunLoop TODO rename this function
func (tracker *AbandonedTracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) RunLoop() {
	temp := make([]AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], 0)

	for _, atx := range tracker.txes {
		switch atx.tx.State {
		case TxConfirmed, TxConfirmedMissingReceipt, TxFatalError:
			// Stop tracking tx when finalized state is obtained
			continue
		case TxInProgress:
			if time.Now().After(atx.invalidTime) {
				tracker.MarkFatal(atx)
			} else {
				temp = append(temp, atx)
			}
		case TxUnstarted, TxUnconfirmed:
			// TODO Handle TxUnstarted, TxUnconfirmed
			//(*tracker.txStore).Abandon() ??
		default:
			tracker.lggr.Panicw(fmt.Sprintf("unhandled transaction state: %v", atx.tx.State))
		}
	}

	tracker.txes = temp
}
