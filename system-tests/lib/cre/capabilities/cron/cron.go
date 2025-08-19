package cron

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	cronregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/cron"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	donlevel "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability/donlevel"
)

const cronConfigTemplate = `""` // Empty config by default

func New() (*capabilities.Capability, error) {
	perDonJobSpecFactory := factory.NewCapabilityJobSpecFactory(
		donlevel.IsEnabled,
		donlevel.EnabledChains,
		donlevel.ConfigResolver,
		donlevel.JobName,
	)

	return capabilities.New(
		cre.CronCapability,
		capabilities.WithJobSpecFn(perDonJobSpecFactory.BuildJobSpecFn(
			cre.CronCapability,
			cronConfigTemplate,
			factory.NoOpExtractor, // No runtime values extraction needed
			factory.BinaryPathBuilder,
		)),
		capabilities.WithCapabilityRegistryV1ConfigFn(cronregistry.CapabilityRegistryConfigFn),
	)
}
