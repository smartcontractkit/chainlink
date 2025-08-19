package webapitarget

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	webapiregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/webapi"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability/donlevel"
)

const webAPITargetConfigTemplate = `"""
[rateLimiter]
GlobalRPS = {{.GlobalRPS}}
GlobalBurst = {{.GlobalBurst}}
PerSenderRPS = {{.PerSenderRPS}}
PerSenderBurst = {{.PerSenderBurst}}
"""`

func New() (*capabilities.Capability, error) {
	perDonJobSpecFactory := factory.NewCapabilityJobSpecFactory(
		donlevel.IsEnabled,
		donlevel.EnabledChains,
		donlevel.ConfigResolver,
		donlevel.JobName,
	)

	return capabilities.New(
		cre.WebAPITargetCapability,
		capabilities.WithJobSpecFn(perDonJobSpecFactory.BuildJobSpecFn(
			cre.WebAPITargetCapability,
			webAPITargetConfigTemplate,
			factory.NoOpExtractor, // No runtime values extraction needed
			func(_ *cre.JobSpecInput, _ cre.CapabilityConfig) (string, error) {
				return "__builtin_web-api-target", nil
			},
		)),
		capabilities.WithCapabilityRegistryV1ConfigFn(webapiregistry.TargetCapabilityRegistryConfigFn),
	)
}
