package bulletprooftxmanager

import (
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

func SetEthClientOnEthConfirmer(ethClient eth.Client, ethConfirmer *EthConfirmer) {
	ethConfirmer.ethClient = ethClient
}

func SetResumeCallbackOnEthBroadcaster(resumeCallback ResumeCallback, ethBroadcaster *EthBroadcaster) {
	ethBroadcaster.resumeCallback = resumeCallback
}
