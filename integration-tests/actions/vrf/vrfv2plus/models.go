package vrfv2plus

import (
	"context"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

type VRFV2PlusWrapperContracts struct {
	VRFV2PlusWrapper  contracts.VRFV2PlusWrapper
	LoadTestConsumers []contracts.VRFv2PlusWrapperLoadTestConsumer
}

type LoadTestConsumer interface {
	GetLoadTestMetrics(ctx context.Context) (*contracts.VRFV2PlusLoadTestMetrics, error)
}
