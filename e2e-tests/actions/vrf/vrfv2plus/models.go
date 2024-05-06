package vrfv2plus

import (
	"github.com/smartcontractkit/chainlink/e2e-tests/contracts"
)

type VRFV2PlusWrapperContracts struct {
	VRFV2PlusWrapper  contracts.VRFV2PlusWrapper
	LoadTestConsumers []contracts.VRFv2PlusWrapperLoadTestConsumer
}
