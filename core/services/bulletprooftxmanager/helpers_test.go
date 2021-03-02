package bulletprooftxmanager

import "github.com/smartcontractkit/chainlink/core/services/eth"

func SetEthClientOnEthConfirmer(ethClient eth.Client, ethConfirmer *ethConfirmer) {
	ethConfirmer.ethClient = ethClient
}
