package compute

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	computeregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/compute"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	donlevel "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability/donlevel"
)

const customComputeConfigTemplate = `"""
NumWorkers = {{.NumWorkers}}
[rateLimiter]
globalRPS = {{.GlobalRPS}}
globalBurst = {{.GlobalBurst}}
perSenderRPS = {{.PerSenderRPS}}
perSenderBurst = {{.PerSenderBurst}}
"""`

func New() (*capabilities.Capability, error) {
	perDonJobSpecFactory := factory.NewCapabilityJobSpecFactory(
		donlevel.IsEnabled,
		donlevel.EnabledChains,
		donlevel.ConfigResolver,
		donlevel.JobName,
	)

	return capabilities.New(
		cre.CustomComputeCapability,
		capabilities.WithJobSpecFn(perDonJobSpecFactory.BuildJobSpecFn(
			cre.CustomComputeCapability,
			customComputeConfigTemplate,
			factory.NoOpExtractor, // No runtime values extraction needed
			func(_ *cre.JobSpecInput, _ cre.CapabilityConfig) (string, error) {
				return "__builtin_custom-compute-action", nil
			},
		)),
		capabilities.WithCapabilityRegistryV1ConfigFn(computeregistry.CapabilityRegistryConfigFn),
	)
}
