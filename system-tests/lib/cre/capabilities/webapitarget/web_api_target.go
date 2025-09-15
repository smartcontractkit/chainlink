package webapitarget

import (
	"github.com/pkg/errors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability/donlevel"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

const flag = cre.WebAPITargetCapability
const configTemplate = `"""
[rateLimiter]
GlobalRPS = {{.GlobalRPS}}
GlobalBurst = {{.GlobalBurst}}
PerSenderRPS = {{.PerSenderRPS}}
PerSenderBurst = {{.PerSenderBurst}}
"""`

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
			factory.NoOpExtractor, // No runtime values extraction needed
			func(_ *cre.JobSpecInput, _ cre.CapabilityConfig) (string, error) {
				return "__builtin_web-api-target", nil
			},
		)),
		capabilities.WithCapabilityRegistryV1ConfigFn(registerWithV1),
	)
}

func registerWithV1(donFlags []string, _ *cre.CapabilitiesAwareNodeSet) ([]keystone_changeset.DONCapabilityWithConfig, error) {
	var capabilities []keystone_changeset.DONCapabilityWithConfig

	if flags.HasFlag(donFlags, flag) {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "web-api-target",
				Version:        "1.0.0",
				CapabilityType: 3, // TARGET
				ResponseType:   1, // OBSERVATION_IDENTICAL
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return capabilities, nil
}
