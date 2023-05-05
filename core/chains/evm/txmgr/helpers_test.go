package txmgr

import (
	"context"
)

func (ec *EthConfirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) SetClient(client TxmClient[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) {
	ec.client = client
}

func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) StartInternal() error {
	return eb.startInternal()
}

func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) CloseInternal() error {
	return eb.closeInternal()
}

func (eb *EthBroadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) DisableUnstartedEthTxAutoProcessing() {
	eb.processUnstartedEthTxsImpl = func(ctx context.Context, fromAddress ADDR) (retryable bool, err error) { return false, nil }
}

func (ec *EthConfirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) StartInternal() error {
	return ec.startInternal()
}

func (ec *EthConfirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) CloseInternal() error {
	return ec.closeInternal()
}

func (er *EthResender[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, SEQ, R, FEE, ADD]) ResendUnconfirmed() error {
	return er.resendUnconfirmed()
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) Abandon(addr ADDR) (err error) {
	return b.abandon(addr)
}
