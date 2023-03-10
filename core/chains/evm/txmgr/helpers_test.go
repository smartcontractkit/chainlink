package txmgr

import evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"

func SetEthClientOnEthConfirmer(ethClient evmclient.Client, ethConfirmer *EthConfirmer) {
	ethConfirmer.ethClient = ethClient
}

func SetResumeCallbackOnEthBroadcaster(resumeCallback ResumeCallback, ethBroadcaster *EthBroadcaster) {
	ethBroadcaster.resumeCallback = resumeCallback
}

func (er *EthResender) ResendUnconfirmed() error {
	return er.resendUnconfirmed()
}
