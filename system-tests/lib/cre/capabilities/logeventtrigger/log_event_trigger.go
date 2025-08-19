package logeventtrigger

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	logeventtriggerregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/logevent"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	chainlevel "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability/chainlevel"
)

const logEventTriggerConfigTemplate = `"""
{
	"chainId": "{{.ChainID}}",
	"network": "{{.NetworkFamily}}",
	"lookbackBlocks": {{.LookbackBlocks}},
	"pollPeriod": {{.PollPeriod}}
}
"""`

func New() (*capabilities.Capability, error) {
	perChainJobSpecFactory := factory.NewCapabilityJobSpecFactory(
		chainlevel.IsEnabled,
		chainlevel.EnabledChains,
		chainlevel.ConfigResolver,
		chainlevel.JobName,
	)

	return capabilities.New(
		cre.LogTriggerCapability,
		capabilities.WithJobSpecFn(perChainJobSpecFactory.BuildJobSpecFn(
			cre.LogTriggerCapability,
			logEventTriggerConfigTemplate,
			func(chainID uint64, _ *cre.NodeMetadata) map[string]any {
				return map[string]any{
					"ChainID":       chainID,
					"NetworkFamily": "evm",
				}
			},
			factory.BinaryPathBuilder,
		)),
		capabilities.WithCapabilityRegistryV1ConfigFn(logeventtriggerregistry.CapabilityRegistryConfigFn),
	)
}
