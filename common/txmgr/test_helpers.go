package txmgr

import (
	"context"
	"time"
)

// TEST ONLY FUNCTIONS
// these need to be exported for the txmgr tests to continue to work

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

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) XXXTestAbandon(addr ADDR) (err error) {
	return b.abandon(addr)
}
