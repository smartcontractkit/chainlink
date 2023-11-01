package txmgr

import (
	"context"
	"fmt"
	"golang.org/x/exp/slices"
	"time"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// TTL is the default time to live for abandoned transactions (6hrs)
const TTL = 6 * time.Hour

// AbandonedErrorMsg occurs when an abandoned tx exceeds its time to live
var AbandonedErrorMsg = fmt.Sprintf(
	"abandoned transaction exceeded time to live of %d hours", int(TTL.Hours()))

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
	// fatalTime represents the time at which this transaction is to be marked fatal
	fatalTime time.Time
}

// AbandonedTracker tracks and handles abandoned transactions
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
	ks      *txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ]
	client  *txmgrtypes.TxmClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	chainID CHAIN_ID
	lggr    logger.Logger
	txes    []AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
}

// NewAbandonedTracker creates a new AbandonedTracker
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
	keystore *txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ],
	client *txmgrtypes.TxmClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
	chainID CHAIN_ID,
	lggr logger.Logger,
) AbandonedTracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE] {
	return AbandonedTracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		txStore: txStore,
		ks:      keystore,
		chainID: chainID,
		client:  client,
		lggr:    lggr.Named("AbandonedTracker"),
		txes:    make([]AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], 0),
	}
}

func (tracker *AbandonedTracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) containsTx(
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
	for _, atx := range tracker.txes {
		if atx.tx.ID == tx.ID {
			return true
		}
	}
	return false
}

// insertAbandonedTx inserts a transaction into the tracker
func (tracker *AbandonedTracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) insertAbandonedTx(
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
	if tracker.containsTx(tx) {
		return
	}

	tracker.lggr.Debugw(fmt.Sprintf("inserting tx %v", tx.ID))
	tracker.txes = append(tracker.txes, AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{
		tx:        tx,
		fatalTime: time.Now().Add(TTL),
	})
}

// markFatal sets a transaction's state to fatal_error
func (tracker *AbandonedTracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) markFatal(
	ctx context.Context,
	atx AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
	tracker.lggr.Infow(fmt.Sprintf("tx %v marked as fatal for exceeding ttl", atx.tx.ID))

	atx.tx.Error.SetValid(AbandonedErrorMsg)

	err := (*tracker.txStore).UpdateTxFatalError(ctx, atx.tx)
	if err != nil {
		tracker.lggr.Errorw(fmt.Sprintf("failed to mark tx %v as fatal", atx.tx.ID))
		// TODO: Handle error
	}
}

// getAbandonedAddresses retrieves fromAddress’s in evm.txes that are not present in the Confirmer's enabledAddresses list
func (tracker *AbandonedTracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) getAbandonedAddresses(enabledAddrs []ADDR) ([]ADDR, error) {
	fromAddresses, err := (*tracker.ks).EnabledAddressesForChain(tracker.chainID)
	if err != nil {
		return nil, err
	}

	var abandoned []ADDR
	for _, addr := range fromAddresses {
		if !slices.Contains(enabledAddrs, addr) {
			abandoned = append(abandoned, addr)
		}
	}

	return abandoned, nil
}

// HandleAbandonedTransactions is called by the Confirmer to track and handle all abandoned transactions
func (tracker *AbandonedTracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) HandleAbandonedTransactions(ctx context.Context, enabledAddrs []ADDR) {
	// Find abandoned addresses
	abandonedAddrs, err := tracker.getAbandonedAddresses(enabledAddrs)
	if err != nil {
		// TODO handle error
	}

	// Get abandoned txs from addresses and insert into the tracker
	for _, addr := range abandonedAddrs {
		seq, err := (*tracker.client).SequenceAt(ctx, addr, nil)
		if err != nil {
			// TODO handle error
		}

		tx, err := (*tracker.txStore).FindTxWithSequence(ctx, addr, seq)
		if err != nil {
			// TODO handle error
		}

		tracker.insertAbandonedTx(tx)
	}

	tracker.handleTransactionStates(ctx)
}

// handleTransactionStates handles all abandoned transactions based on their current state.
// Transactions with finalized states are no longer tracked, while transactions which
// exceed their ttl are marked as fatal.
func (tracker *AbandonedTracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) handleTransactionStates(ctx context.Context) {
	temp := make([]AbandonedTx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], 0)

	for _, atx := range tracker.txes {
		switch atx.tx.State {
		case TxConfirmed, TxConfirmedMissingReceipt, TxFatalError:
			// Stop tracking tx when finalized state is obtained
			continue
		case TxInProgress:
			if time.Now().After(atx.fatalTime) {
				tracker.markFatal(ctx, atx)
				continue
			}
			temp = append(temp, atx)
		case TxUnstarted, TxUnconfirmed:
			if time.Now().After(atx.fatalTime) {
				// TODO: Handle cancelling TxUnstarted, TxUnconfirmed
				continue
			}
			temp = append(temp, atx)
		default:
			// This should never happen unless a new transaction state is added
			tracker.lggr.Panicw(fmt.Sprintf("unhandled transaction state: %v", atx.tx.State))
		}
	}

	tracker.txes = temp
}
