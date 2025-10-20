package logeventtrigger

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability/chainlevel"
)

const flag = cre.LogEventTriggerCapability

type LogEventTrigger struct{}

func (o *LogEventTrigger) Flag() cre.CapabilityFlag {
	return flag
}

func (o *LogEventTrigger) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	donsMetadata := topology.DonsMetadataWithFlag(flag)
	if len(donsMetadata) == 0 {
		return nil, nil
	}

	capabilities := make(map[uint64][]keystone_changeset.DONCapabilityWithConfig)
	for _, donMetadata := range donsMetadata {
		if capabilities[donMetadata.ID] == nil {
			capabilities[donMetadata.ID] = []keystone_changeset.DONCapabilityWithConfig{}
		}

		for _, chainID := range donMetadata.CapabilitiesAwareNodeSet().GetChainCapabilityConfigs()[flag].EnabledChains {
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
	donTopology *cre.DonTopology,
	creEnv *cre.Environment,
) error {
	dons := donTopology.DonsWithFlag(flag)
	if len(dons) == 0 {
		return nil
	}

	perDonJobSpecFactory, fErr := factory.NewCapabilityJobSpecFactory(
		creEnv.RegistryChainSelector,
		chainlevel.CapabilityEnabler,
		chainlevel.EnabledChainsProvider,
		chainlevel.ConfigResolver,
		chainlevel.JobNamer,
	)

	if fErr != nil {
		return errors.Wrap(fErr, "failed to create capability job spec factory")
	}

	bcOuts := make([]*blockchain.Output, len(creEnv.Blockchains))
	for i, b := range creEnv.Blockchains {
		bcOuts[i] = b.CtfOutput()
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
		CreEnvironment: creEnv,
		DonTopology:    donTopology,
		NodeSets:       donTopology.Dons.AsNodeSetWithChainCapabilities(),
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
		jobErr := jobs.Create(ctx, creEnv.CldfEnvironment.Offchain, donTopology, jobSpecs)
		if jobErr != nil {
			return fmt.Errorf("failed to create http action jobs for don %s: %w", don.Name, jobErr)
		}
	}

	return nil
}
