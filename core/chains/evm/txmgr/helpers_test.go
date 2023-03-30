package txmgr

import (
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
)

func SetEthClientOnEthConfirmer(ethClient evmclient.Client, ethConfirmer *EvmEthConfirmer) {
	ethConfirmer.ethClient = ethClient
}

func SetResumeCallbackOnEthBroadcaster(resumeCallback ResumeCallback, ethBroadcaster *EvmEthBroadcaster) {
	ethBroadcaster.resumeCallback = resumeCallback
}

func (er *EthResender[ADDR, TX_HASH, BLOCK_HASH]) ResendUnconfirmed() error {
	return er.resendUnconfirmed()
}
