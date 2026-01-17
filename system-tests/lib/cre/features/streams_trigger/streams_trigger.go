package streams_trigger

import (
	"context"

	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

const flag = cre.StreamsTriggerCapability

type StreamsTrigger struct{}

func (s *StreamsTrigger) Flag() cre.CapabilityFlag {
	return flag
}

// PreEnvStartup registers the streams-trigger capability in the on-chain Capabilities Registry.
// This capability is created by the LLO plugin's CRE transmitter when an LLO job is deployed.
func (s *StreamsTrigger) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	// Register streams-trigger@2.0.0 in the Capabilities Registry
	// This is a DON-level trigger capability that emits LLO OCR reports
	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   "streams-trigger",
			Version:        "2.0.0",
			CapabilityType: 0, // TRIGGER
		},
		Config: &capabilitiespb.CapabilityConfig{
			// LocalOnly = false means this capability can be consumed by remote DONs
			LocalOnly: false,
		},
	}}

	testLogger.Info().
		Str("capability", "streams-trigger@2.0.0").
		Str("don", don.Name).
		Msg("Registering streams-trigger capability for DON")

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

// PostEnvStartup is called after the environment is started.
// For streams-trigger, there's no additional job deployment needed here
// because the capability is created by the LLO plugin when its job is deployed.
//
// NOTE: An LLO job must be deployed separately to this DON for the capability
// to actually produce events. The LLO job should use the CRE transmitter:
//
//	[pluginConfig]
//	transmitterType = "cre"
//	[pluginConfig.transmitterConfig]
//	triggerCapabilityName = "streams-trigger@2.0.0"
func (s *StreamsTrigger) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	testLogger.Info().
		Str("capability", "streams-trigger@2.0.0").
		Str("don", don.Name).
		Msg("streams-trigger capability registered. Deploy an LLO job to this DON to enable the trigger.")

	// No additional job deployment for streams-trigger
	// The capability is instantiated by the LLO plugin's CRE transmitter
	return nil
}
