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
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	capabilities := []keystone_changeset.DONCapabilityWithConfig{}

	for _, chainID := range don.CapabilitiesAwareNodeSet().GetChainCapabilityConfigs()[flag].EnabledChains {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   fmt.Sprintf("log-event-trigger-evm-%d", chainID),
				Version:        "1.0.0",
				CapabilityType: 0, // TRIGGER
				ResponseType:   0, // REPORT
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
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
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
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

	var nodeSet cre.NodeSetWithCapabilityConfigs
	for _, ns := range dons.AsNodeSetWithChainCapabilities() {
		if ns.GetName() == don.Name {
			nodeSet = ns
			break
		}
	}
	if nodeSet == nil {
		return fmt.Errorf("could not find node set for Don named '%s'", don.Name)
	}

	jobSpecs, specErr := perDonJobSpecFactory.BuildJobSpec(
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
		Don:            don,
		NodeSet:        nodeSet,
	})
	if specErr != nil {
		return fmt.Errorf("failed to build job spec for http action capability: %w", specErr)
	}
	if len(jobSpecs) == 0 {
		return fmt.Errorf("no job specs created for '%s' capability, even though it is enabled", flag)
	}

	// pass all dons, since some jobs might need to be created on multiple ones
	jobErr := jobs.Create(ctx, creEnv.CldfEnvironment.Offchain, dons, jobSpecs)
	if jobErr != nil {
		return fmt.Errorf("failed to create http action jobs for don %s: %w", don.Name, jobErr)
	}

	return nil
}
