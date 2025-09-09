package httptrigger

import (
	"github.com/pkg/errors"

	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	donlevel "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability/donlevel"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

const flag = cre.HTTPTriggerCapability
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
			httpTriggerConfigTemplate,
			factory.NoOpExtractor,
			factory.BinaryPathBuilder,
		)),
		capabilities.WithGatewayJobHandlerConfigFn(handlerConfig),
		capabilities.WithCapabilityRegistryV1ConfigFn(registerWithV1),
	)
}

func handlerConfig(donMetadata *cre.DonMetadata) (cre.HandlerTypeToConfig, error) {
	// if any of the DONs have http action or http trigger capability, we need to add a http handler to the jobspec for the gateway node
	if !flags.HasFlag(donMetadata.Flags, flag) {
		return nil, nil
	}

	return map[string]string{coregateway.HTTPCapabilityType: `
ServiceName = "workflows"
[gatewayConfig.Dons.Handlers.Config]
maxTriggerRequestDurationMs = 5_000
metadataPullIntervalMs = 1_000
metadataAggregationIntervalMs = 1_000
[gatewayConfig.Dons.Handlers.Config.NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10
[gatewayConfig.Dons.Handlers.Config.UserRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10`}, nil
}

func registerWithV1(donFlags []string, _ *cre.CapabilitiesAwareNodeSet) ([]keystone_changeset.DONCapabilityWithConfig, error) {
	var capabilities []keystone_changeset.DONCapabilityWithConfig

	if flags.HasFlag(donFlags, flag) {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "http-trigger",
				Version:        "1.0.0-alpha",
				CapabilityType: 0, // TRIGGER
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return capabilities, nil
}
