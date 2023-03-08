package txmgr

import evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"

func SetEthClientOnEthConfirmer[HEAD any](ethClient evmclient.Client, ethConfirmer *EthConfirmer[HEAD]) {
	ethConfirmer.ethClient = ethClient
}

func SetResumeCallbackOnEthBroadcaster[HEAD any](resumeCallback ResumeCallback, ethBroadcaster *EthBroadcaster[HEAD]) {
	ethBroadcaster.resumeCallback = resumeCallback
}

func (er *EthResender) ResendUnconfirmed() error {
	return er.resendUnconfirmed()
}
