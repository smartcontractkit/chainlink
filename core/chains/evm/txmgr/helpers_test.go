package txmgr

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/common/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
)

func processUnstartedEthTxsNoOp[ADDR types.Hashable[ADDR]](ctx context.Context, fromAddress ADDR) (err error, retryable bool) {
	return nil, false
}

func SetEthClientOnEthConfirmer(ethClient evmclient.Client, ethConfirmer *EvmConfirmer) {
	ethConfirmer.ethClient = ethClient
}

func SetResumeCallbackOnEthBroadcaster(resumeCallback ResumeCallback, ethBroadcaster *EvmBroadcaster) {
	ethBroadcaster.resumeCallback = resumeCallback
}

func (eb *EthBroadcaster[ADDR, TX_HASH, BLOCK_HASH]) StartInternal() error {
	return eb.startInternal()
}

func (eb *EthBroadcaster[ADDR, TX_HASH, BLOCK_HASH]) CloseInternal() error {
	return eb.closeInternal()
}

func (eb *EthBroadcaster[ADDR, TX_HASH, BLOCK_HASH]) DisableUnstartedEthTxAutoProcessing() {
	eb.processUnstartedEthTxsImpl = processUnstartedEthTxsNoOp[ADDR]
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
