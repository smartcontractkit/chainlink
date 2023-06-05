package txmgr

import (
	"context"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
)

// TEST ONLY FUNCTIONS
// these need to be exported for the txmgr tests to continue to work

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) XXXTestSetClient(client txmgrtypes.TxmClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) {
	ec.client = client
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) XXXTestStartInternal() error {
	return eb.startInternal()
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) XXXTestCloseInternal() error {
	return eb.closeInternal()
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) XXXTestDisableUnstartedTxAutoProcessing() {
	eb.processUnstartedTxsImpl = func(ctx context.Context, fromAddress ADDR) (retryable bool, err error) { return false, nil }
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) XXXTestStartInternal() error {
	return ec.startInternal()
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) XXXTestCloseInternal() error {
	return ec.closeInternal()
}

func (er *Resender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, R, FEE, ADD]) XXXTestResendUnconfirmed() error {
	return er.resendUnconfirmed()
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) XXXTestAbandon(addr ADDR) (err error) {
	return b.abandon(addr)
}
