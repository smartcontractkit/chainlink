package txmgr

import (
	"context"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
)

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) SetClient(client txmgrtypes.TxmClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) {
	ec.client = client
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) StartInternal() error {
	return eb.startInternal()
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) CloseInternal() error {
	return eb.closeInternal()
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) DisableUnstartedTxAutoProcessing() {
	eb.processUnstartedTxsImpl = func(ctx context.Context, fromAddress ADDR) (retryable bool, err error) { return false, nil }
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) StartInternal() error {
	return ec.startInternal()
}

func (ec *Confirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) CloseInternal() error {
	return ec.closeInternal()
}

func (er *Resender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, R, FEE, ADD]) ResendUnconfirmed() error {
	return er.resendUnconfirmed()
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, FEE_UNIT]) Abandon(addr ADDR) (err error) {
	return b.abandon(addr)
}
