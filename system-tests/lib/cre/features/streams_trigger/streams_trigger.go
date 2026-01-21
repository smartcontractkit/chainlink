package streams_trigger

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/durationpb"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	// data-streams-deploy provides battle-tested changesets for LLO contract deployment and job specs.
	// See: https://github.com/smartcontractkit/data-streams-deploy
	//
	// Available changesets:
	//   - configurator_v0_5_0.DeployConfiguratorChangeset - Deploy Configurator contract
	//   - configurator_v0_5_0.SetProductionConfigChangeset - Set OCR production config
	//   - channel_config_store.DeployChannelConfigStoreChangeset - Deploy ChannelConfigStore
	//   - channel_config_store.SetChannelDefinitionChangeset - Set channel definitions
	//
	// Available job specs:
	//   - jobs.LLOJobSpec - LLO job with CRE transmitter support
	//   - jobs.Transmitter - CRE transmitter configuration
)

const (
	// Remote trigger configuration
	registrationRefresh = 20 * time.Second
	registrationExpiry  = 60 * time.Second
)

const flag = cre.StreamsTriggerCapability

type StreamsTrigger struct{}

func (s *StreamsTrigger) Flag() cre.CapabilityFlag {
	return flag
}

// PreEnvStartup registers the streams-trigger capability in the on-chain Capabilities Registry.
// This capability is created by the LLO plugin's CRE transmitter when an LLO job is deployed.
//
// Contract deployment (Configurator, ChannelConfigStore) should be done using changesets from
// github.com/smartcontractkit/data-streams-deploy:
//   - configurator_v0_5_0.DeployConfiguratorChangeset
//   - channel_config_store.DeployChannelConfigStoreChangeset
func (s *StreamsTrigger) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	// Calculate minResponsesToAggregate based on F+1 threshold
	// For a 4-node DON with F=1, this is 2
	numNodes := len(don.NodesMetadata)
	if numNodes <= 0 {
		numNodes = 4 // default to 4 nodes
	}
	faultyNodes := uint32((numNodes - 1) / 3)
	minResponsesToAggregate := faultyNodes + 1

	// Register streams-trigger@2.0.0 in the Capabilities Registry
	// This is a DON-level trigger capability that emits LLO OCR reports
	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   "streams-trigger",
			Version:        "2.0.0",
			CapabilityType: 0, // TRIGGER
		},
		Config: &capabilitiespb.CapabilityConfig{
			LocalOnly: don.HasOnlyLocalCapabilities(),
			// Configure the remote trigger config for cross-DON communication
			// This tells the TriggerSubscriber how many responses to wait for before aggregating
			RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
				RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
					RegistrationRefresh:     durationpb.New(registrationRefresh),
					RegistrationExpiry:      durationpb.New(registrationExpiry),
					MinResponsesToAggregate: minResponsesToAggregate,
					MessageExpiry:           durationpb.New(2 * registrationExpiry),
					MaxBatchSize:            25,
					BatchCollectionPeriod:   durationpb.New(200 * time.Millisecond),
				},
			},
		},
	}}

	testLogger.Info().
		Str("capability", "streams-trigger@2.0.0").
		Str("don", don.Name).
		Bool("localOnly", don.HasOnlyLocalCapabilities()).
		Uint32("minResponsesToAggregate", minResponsesToAggregate).
		Msg("Registering streams-trigger capability for DON with remote trigger config")

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

// PostEnvStartup is called after the environment is started.
//
// For streams-trigger to produce events, an LLO job must be deployed to this DON.
// This is typically done by the test using data-streams-deploy changesets.
//
// Example LLO job deployment using data-streams-deploy:
//
//	import (
//	    configurator_v0_5_0 "github.com/smartcontractkit/data-streams-deploy/changeset/configurator/v0_5_0"
//	    channel_config_store "github.com/smartcontractkit/data-streams-deploy/changeset/channel-config-store"
//	    dsdeploy_jobs "github.com/smartcontractkit/data-streams-deploy/jobs"
//	)
//
//	// 1. Deploy Configurator
//	configurator_v0_5_0.DeployConfiguratorChangeset.Apply(env, configurator_v0_5_0.DeployConfiguratorConfig{
//	    ChainsToDeploy: []uint64{chainSelector},
//	})
//
//	// 2. Deploy ChannelConfigStore
//	channel_config_store.DeployChannelConfigStoreChangeset.Apply(env, channel_config_store.DeployChannelConfigStoreConfig{
//	    ChainsToDeploy: []uint64{chainSelector},
//	})
//
//	// 3. Set OCR config
//	configurator_v0_5_0.SetProductionConfigChangeset.Apply(env, configurator_v0_5_0.SetProductionConfig{...})
//
//	// 4. Set channel definitions
//	channel_config_store.SetChannelDefinitionChangeset.Apply(env, channel_config_store.SetChannelDefinitionsConfig{...})
//
//	// 5. Deploy LLO job with CRE transmitter
//	jobSpec := &dsdeploy_jobs.LLOJobSpec{
//	    PluginConfig: dsdeploy_jobs.PluginConfigLLO{
//	        Transmitters: []dsdeploy_jobs.Transmitter{{
//	            Type: "cre",
//	            Opts: dsdeploy_jobs.TransmitterOpts{
//	                TriggerCapabilityName:    "streams-trigger",
//	                TriggerCapabilityVersion: "2.0.0",
//	            },
//	        }},
//	    },
//	}
//
// See system-tests/tests/smoke/cre/llo_deployment_helpers.go for helper functions.
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
		Msg("streams-trigger capability registered - deploy LLO jobs to enable trigger events")

	// NOTE: Actual LLO contract and job deployment is done by the test, not here.
	// This is because:
	// 1. Contract addresses (Configurator, ChannelConfigStore) may vary per test
	// 2. Stream job configuration is test-specific (stream IDs, report formats)
	// 3. The data-streams-deploy changesets provide flexible deployment options
	//
	// See:
	// - system-tests/tests/smoke/cre/llo_deployment_helpers.go - Helper functions using data-streams-deploy
	// - system-tests/tests/smoke/cre/v2_llo_streams_trigger_test.go - E2E test example

	return nil
}
