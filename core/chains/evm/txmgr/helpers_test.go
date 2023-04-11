package txmgr

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/common/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
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

func StartInternalOnBroadcaster(eb *EvmBroadcaster) error {
	return eb.startInternal()
}

func CloseInternalOnBroadcaster(eb *EvmBroadcaster) error {
	return eb.closeInternal()
}

func DisableUnstartedEthTxAutoProcessingOnBroadcaster(eb *EvmBroadcaster) {
	eb.processUnstartedEthTxsImpl = processUnstartedEthTxsNoOp[*evmtypes.Address]
}

func StartInternalOnConfirmer(ec *EvmConfirmer) error {
	return ec.startInternal()
}

func CloseInternalOnConfirmer(ec *EvmConfirmer) error {
	return ec.closeInternal()
}

func (er *EthResender[ADDR, TX_HASH, BLOCK_HASH]) ResendUnconfirmed() error {
	return er.resendUnconfirmed()
}
