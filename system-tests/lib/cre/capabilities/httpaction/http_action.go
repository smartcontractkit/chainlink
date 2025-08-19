package httpaction

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	httpactionregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/httpaction"
	httpactionhandler "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway/handlers/httpaction"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	donlevel "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability/donlevel"
)

const httpActionConfigTemplate = `"""
{
	"proxyMode": "{{.ProxyMode}}",
	"incomingRateLimiter": {
		"globalBurst": {{.IncomingGlobalBurst}},
		"globalRPS": {{.IncomingGlobalRPS}},
		"perSenderBurst": {{.IncomingPerSenderBurst}},
		"perSenderRPS": {{.IncomingPerSenderRPS}}
	},
	"outgoingRateLimiter": {
		"globalBurst": {{.OutgoingGlobalBurst}},
		"globalRPS": {{.OutgoingGlobalRPS}},
		"perSenderBurst": {{.OutgoingPerSenderBurst}},
		"perSenderRPS": {{.OutgoingPerSenderRPS}}
	}
}
"""`

func New() (*capabilities.Capability, error) {
	perDonJobSpecFactory := factory.NewCapabilityJobSpecFactory(
		donlevel.IsEnabled,
		donlevel.EnabledChains,
		donlevel.ConfigResolver,
		donlevel.JobName,
	)

	return capabilities.New(
		cre.HTTPActionCapability,
		capabilities.WithJobSpecFn(perDonJobSpecFactory.BuildJobSpecFn(
			cre.HTTPActionCapability,
			httpActionConfigTemplate,
			factory.NoOpExtractor, // No runtime values extraction needed
			factory.BinaryPathBuilder,
		)),
		capabilities.WithGatewayJobHandlerConfigFn(httpactionhandler.HandlerConfigFn),
		capabilities.WithCapabilityRegistryV1ConfigFn(httpactionregistry.CapabilityRegistryConfigFn),
	)
}
