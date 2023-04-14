package txmgr

import (
	"context"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
)

func (ec *EthConfirmer[ADDR, TX_HASH, BLOCK_HASH]) SetEthClient(ethClient evmclient.Client) {
	ec.ethClient = ethClient
}

func (eb *EthBroadcaster[ADDR, TX_HASH, BLOCK_HASH]) StartInternal() error {
	return eb.startInternal()
}

func (eb *EthBroadcaster[ADDR, TX_HASH, BLOCK_HASH]) CloseInternal() error {
	return eb.closeInternal()
}

func (eb *EthBroadcaster[ADDR, TX_HASH, BLOCK_HASH]) DisableUnstartedEthTxAutoProcessing() {
	eb.processUnstartedEthTxsImpl = func(ctx context.Context, fromAddress ADDR) (retryable bool, err error) { return false, nil }
}

func (ec *EthConfirmer[ADDR, TX_HASH, BLOCK_HASH]) StartInternal() error {
	return ec.startInternal()
}

func (ec *EthConfirmer[ADDR, TX_HASH, BLOCK_HASH]) CloseInternal() error {
	return ec.closeInternal()
}

func (er *EthResender[ADDR, TX_HASH, BLOCK_HASH]) ResendUnconfirmed() error {
	return er.resendUnconfirmed()
}
