package mock

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	mockregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/mock"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	donlevel "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability/donlevel"
)

const mockConfigTemplate = `"""port={{.Port}}"""`

func New() (*capabilities.Capability, error) {
	perDonJobSpecFactory := factory.NewCapabilityJobSpecFactory(
		donlevel.IsEnabled,
		donlevel.EnabledChains,
		donlevel.ConfigResolver,
		donlevel.JobName,
	)

	return capabilities.New(
		cre.MockCapability,
		capabilities.WithJobSpecFn(perDonJobSpecFactory.BuildJobSpecFn(
			cre.MockCapability,
			mockConfigTemplate,
			factory.NoOpExtractor, // No runtime values extraction needed
			factory.BinaryPathBuilder,
		)),
		capabilities.WithCapabilityRegistryV1ConfigFn(mockregistry.CapabilityRegistryConfigFn),
	)
}
