package txmgr

import (
	"context"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
)

func (ec *EthConfirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) SetClient(client txmgrtypes.TxmClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) {
	ec.client = client
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, UNIT]) StartInternal() error {
	return eb.startInternal()
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, UNIT]) CloseInternal() error {
	return eb.closeInternal()
}

func (eb *Broadcaster[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, UNIT]) DisableUnstartedTxAutoProcessing() {
	eb.processUnstartedTxsImpl = func(ctx context.Context, fromAddress ADDR) (retryable bool, err error) { return false, nil }
}

func (ec *EthConfirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) StartInternal() error {
	return ec.startInternal()
}

func (ec *EthConfirmer[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]) CloseInternal() error {
	return ec.closeInternal()
}

func (er *Resender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, R, FEE, ADD]) ResendUnconfirmed() error {
	return er.resendUnconfirmed()
}

func (b *Txm[CHAIN_ID, HEAD, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD, UNIT]) Abandon(addr ADDR) (err error) {
	return b.abandon(addr)
}
