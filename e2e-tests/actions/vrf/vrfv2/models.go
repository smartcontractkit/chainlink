package vrfv2

import (
	"github.com/smartcontractkit/chainlink/e2e-tests/contracts"
)

type VRFV2WrapperContracts struct {
	VRFV2Wrapper      contracts.VRFV2Wrapper
	LoadTestConsumers []contracts.VRFv2WrapperLoadTestConsumer
}
