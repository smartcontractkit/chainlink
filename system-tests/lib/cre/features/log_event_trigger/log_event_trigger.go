package logeventtrigger

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability/chainlevel"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
)

const flag = cre.LogEventTriggerCapability

type LogEventTrigger struct{}

func (o *LogEventTrigger) Flag() cre.CapabilityFlag {
	return flag
}

func (o *LogEventTrigger) PreEnvStartup(
	testLogger zerolog.Logger,
	registryChainSelector uint64,
	cldfEnv *cldf.Environment,
	provider infra.Provider,
	topology *cre.Topology,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	capabilityConfigs cre.CapabilityConfigs, // move to Topology
	contractVersions map[string]string,
	gatewayJobConfigs map[cre.NodeUUID]*config.GatewayConfig,
) (*cre.PreEnvStartupOutput, error) {
	donsMetadata := topology.DonsMetadataWithFlag(o.Flag())
	if len(donsMetadata) == 0 {
		return nil, nil
	}

	capabilities := make(map[uint64][]keystone_changeset.DONCapabilityWithConfig)
	for _, donMetadata := range donsMetadata {
		if capabilities[donMetadata.ID] == nil {
			capabilities[donMetadata.ID] = []keystone_changeset.DONCapabilityWithConfig{}
		}

		for _, chainID := range donMetadata.CapabilitiesAwareNodeSet().GetChainCapabilityConfigs()[o.Flag()].EnabledChains {
			capabilities[donMetadata.ID] = append(capabilities[donMetadata.ID], keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   fmt.Sprintf("log-event-trigger-evm-%d", chainID),
					Version:        "1.0.0",
					CapabilityType: 0, // TRIGGER
					ResponseType:   0, // REPORT
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfigs: capabilities,
	}, nil
}

const configTemplate = `"""
{
	"chainId": "{{.ChainID}}",
	"network": "{{.NetworkFamily}}",
	"lookbackBlocks": {{.LookbackBlocks}},
	"pollPeriod": {{.PollPeriod}}
}
"""`

func (o *LogEventTrigger) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	creEnv *cre.Environment,
	nodeSetOutput []*cre.WrappedNodeOutput,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	contractVersions map[string]string,
	provider infra.Provider,
	capabilityConfigs map[string]cre.CapabilityConfig,
) error {
	dons := creEnv.DonTopology.DonsWithFlag(flag)
	if len(dons) == 0 {
		return nil
	}

	perDonJobSpecFactory, fErr := factory.NewCapabilityJobSpecFactory(
		chainlevel.CapabilityEnabler,
		chainlevel.EnabledChainsProvider,
		chainlevel.ConfigResolver,
		chainlevel.JobNamer,
	)

	if fErr != nil {
		return errors.Wrap(fErr, "failed to create capability job spec factory")
	}

	bcOuts := make([]*blockchain.Output, len(blockchainOutputs))
	for i, b := range blockchainOutputs {
		bcOuts[i] = b.BlockchainOutput
	}

	donsToJobSpecs, specErr := perDonJobSpecFactory.BuildJobSpec(
		flag,
		configTemplate,
		func(chainID uint64, _ *cre.Node) map[string]any {
			return map[string]any{
				"ChainID":       chainID,
				"NetworkFamily": "evm",
			}
		},
		factory.BinaryPathBuilder,
	)(&cre.JobSpecInput{
		CldEnvironment: creEnv.CldfEnvironment,
		DonTopology:    creEnv.DonTopology,
		InfraInput:     provider,
		/// Capabilities:  // not needed,
		NodeSets:          creEnv.DonTopology.Dons.AsNodeSetWithChainCapabilities(),
		CapabilityConfigs: capabilityConfigs,
	})
	if specErr != nil {
		return fmt.Errorf("failed to build job spec for http action capability: %w", specErr)
	}

	for _, don := range dons {
		jobSpecs, ok := donsToJobSpecs[don.ID]
		if !ok {
			continue
		}
		// pass whole topology, since some jobs might need to be created on multiple DONs
		jobErr := jobs.Create(ctx, creEnv.CldfEnvironment.Offchain, creEnv.DonTopology, jobSpecs)
		if jobErr != nil {
			return fmt.Errorf("failed to create http action jobs for don %s: %w", don.Name, jobErr)
		}
	}

	return nil
}
