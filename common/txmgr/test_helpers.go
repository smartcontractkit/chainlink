package txmgr

import (
	"context"
	"fmt"
	"time"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
)

// TEST ONLY FUNCTIONS
// these need to be exported for the txmgr tests to continue to work

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) XXXTestSetClient(client txmgrtypes.TxmClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) {
	ec.client = client
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) XXXTestSetTTL(ttl time.Duration) {
	tr.ttl = ttl
}

func (tr *Tracker[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) XXXDeliverBlock(blockHeight int64) {
	tr.mb.Deliver(blockHeight)
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) XXXTestStartInternal(ctx context.Context) error {
	return eb.startInternal(ctx)
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) XXXTestCloseInternal() error {
	return eb.closeInternal()
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) XXXTestDisableUnstartedTxAutoProcessing() {
	eb.processUnstartedTxsImpl = func(ctx context.Context, fromAddress ADDR) (retryable bool, err error) { return false, nil }
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) XXXTestStartInternal() error {
	ctx, cancel := ec.stopCh.NewCtx()
	defer cancel()
	return ec.startInternal(ctx)
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) XXXTestCloseInternal() error {
	return ec.closeInternal()
}

func (er *Resender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) XXXTestResendUnconfirmed() error {
	ctx, cancel := er.stopCh.NewCtx()
	defer cancel()
	return er.resendUnconfirmed(ctx)
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) XXXTestAbandon(addr ADDR) (err error) {
	return b.abandon(addr)
}

func (b *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) XXXTestInsertTx(fromAddr ADDR, tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	as, ok := b.addressStates[fromAddr]
	if !ok {
		return fmt.Errorf("address not found: %s", fromAddr)
	}

	as.allTxs[tx.ID] = tx
	if tx.IdempotencyKey != nil && *tx.IdempotencyKey != "" {
		as.idempotencyKeyToTx[*tx.IdempotencyKey] = tx
	}
	for i := 0; i < len(tx.TxAttempts); i++ {
		txAttempt := tx.TxAttempts[i]
		as.attemptHashToTxAttempt[txAttempt.Hash] = &txAttempt
	}

	switch tx.State {
	case TxUnstarted:
		as.unstartedTxs.AddTx(tx)
	case TxInProgress:
		as.inprogressTx = tx
	case TxUnconfirmed:
		as.unconfirmedTxs[tx.ID] = tx
	case TxConfirmed:
		as.confirmedTxs[tx.ID] = tx
	case TxConfirmedMissingReceipt:
		as.confirmedMissingReceiptTxs[tx.ID] = tx
	case TxFatalError:
		as.fatalErroredTxs[tx.ID] = tx
	}

	return nil
}

func (b *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) XXXTestFindTxs(
	txStates []txmgrtypes.TxState,
	filter func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool,
	txIDs ...int64,
) []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	txs := []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	for _, as := range b.addressStates {
		txs = append(txs, as.findTxs(txStates, filter, txIDs...)...)
	}

	return txs
}
