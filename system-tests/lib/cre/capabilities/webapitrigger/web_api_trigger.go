package webapitrigger

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	webapiregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/webapi"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	donlevel "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability/donlevel"
)

const webAPITriggerConfigTemplate = `""`

func New() (*capabilities.Capability, error) {
	perDonJobSpecFactory := factory.NewCapabilityJobSpecFactory(
		donlevel.IsEnabled,
		donlevel.EnabledChains,
		donlevel.ConfigResolver,
		donlevel.JobName,
	)

	return capabilities.New(
		cre.WebAPITriggerCapability,
		capabilities.WithJobSpecFn(perDonJobSpecFactory.BuildJobSpecFn(
			cre.WebAPITriggerCapability,
			webAPITriggerConfigTemplate,
			factory.NoOpExtractor, // No runtime values extraction needed
			func(_ *cre.JobSpecInput, _ cre.CapabilityConfig) (string, error) {
				return "__builtin_web-api-trigger", nil
			},
		)),
		capabilities.WithCapabilityRegistryV1ConfigFn(webapiregistry.TriggerCapabilityRegistryConfigFn),
	)
}
