package txmgr

import (
	"context"
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
	return ec.startInternal(ec.ctx)
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) XXXTestCloseInternal() error {
	return ec.closeInternal()
}

func (er *Resender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) XXXTestResendUnconfirmed() error {
	return er.resendUnconfirmed(er.ctx)
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) XXXTestAbandon(addr ADDR) (err error) {
	return b.abandon(addr)
}
