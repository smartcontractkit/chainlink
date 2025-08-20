package mock

import (
	"github.com/pkg/errors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	donlevel "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability/donlevel"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

const flag = cre.MockCapability
const configTemplate = `"""port={{.Port}}"""`

func New() (*capabilities.Capability, error) {
	perDonJobSpecFactory, fErr := factory.NewCapabilityJobSpecFactory(
		donlevel.CapabilityEnabler,
		donlevel.EnabledChainsProvider,
		donlevel.ConfigResolver,
		donlevel.JobNamer,
	)

	if fErr != nil {
		return nil, errors.Wrap(fErr, "failed to create capability job spec factory")
	}

	return capabilities.New(
		flag,
		capabilities.WithJobSpecFn(perDonJobSpecFactory.BuildJobSpec(
			flag,
			configTemplate,
			factory.NoOpExtractor,
			factory.BinaryPathBuilder,
		)),
		capabilities.WithCapabilityRegistryV1ConfigFn(registerWithV1),
	)
}

func registerWithV1(donFlags []string, _ *cre.CapabilitiesAwareNodeSet) ([]keystone_changeset.DONCapabilityWithConfig, error) {
	var capabilities []keystone_changeset.DONCapabilityWithConfig

	// Configuration details are irrelevant here; this is only included to satisfy validation requirements
	if flags.HasFlag(donFlags, flag) {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "mock",
				Version:        "1.0.0",
				CapabilityType: 0, // TRIGGER
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return capabilities, nil
}
