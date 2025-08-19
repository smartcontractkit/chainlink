package httptrigger

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	httpregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/httptrigger"
	httpactionhandler "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway/handlers/httpaction"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	donlevel "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability/donlevel"
)

const httpTriggerConfigTemplate = `"""
{
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
		cre.HTTPTriggerCapability,
		capabilities.WithJobSpecFn(perDonJobSpecFactory.BuildJobSpecFn(
			cre.HTTPTriggerCapability,
			httpTriggerConfigTemplate,
			factory.NoOpExtractor, // No runtime values extraction needed
			factory.BinaryPathBuilder,
		)),
		capabilities.WithGatewayJobHandlerConfigFn(httpactionhandler.HandlerConfigFn),
		capabilities.WithCapabilityRegistryV1ConfigFn(httpregistry.CapabilityRegistryConfigFn),
	)
}
