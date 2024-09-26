package vrfv2plus

import (
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

type VRFV2PlusWrapperContracts struct {
	VRFV2PlusWrapper contracts.VRFV2PlusWrapper
	WrapperConsumers []contracts.VRFv2PlusWrapperLoadTestConsumer
}
