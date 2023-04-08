package txmgr

import (
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
)

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

func SetIsUnitTestInstanceOnBroadcaster(eb *EvmBroadcaster) {
	eb.isUnitTestInstance = true
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
