package gateway

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	gatewayconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config/gateway"
	gatewayjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway"
)

func New(extraAllowedPorts []int, extraAllowedIPs []string, extraAllowedIPsCIDR []string) (*capabilities.Capability, error) {
	return capabilities.New(
		cre.GatewayDON,
		capabilities.WithNodeConfigFn(gatewayconfig.GenerateConfigFn),
		capabilities.WithJobSpecFn(gatewayjobs.JobSpecFn(extraAllowedPorts, extraAllowedIPs, extraAllowedIPsCIDR)),
	)
}
