package consensus

import (
	"context"
	"fmt"

	"dario.cat/mergo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
)

const flag = cre.ConsensusCapability
const consensusLabelledName = "consensus"

type Consensus struct{}

func (c *Consensus) Flag() cre.CapabilityFlag {
	return flag
}

func (c *Consensus) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   consensusLabelledName,
			Version:        "1.0.0-alpha",
			CapabilityType: 2, // CONSENSUS
			ResponseType:   0, // REPORT
		},
		Config: &capabilitiespb.CapabilityConfig{
			LocalOnly: don.HasOnlyLocalCapabilities(),
		},
		UseCapRegOCRConfig: true,
	}}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
		CapabilityToOCR3Config: map[string]*ocr3.OracleConfig{
			consensusLabelledName: contracts.DefaultOCR3Config(),
		},
		CapabilityToExtraSignerFamilies: cre.CapabilityToExtraSignerFamilies(
			cre.OCRExtraSignerFamilies(creEnv.Blockchains),
			consensusLabelledName,
		),
	}, nil
}

const ContractQualifier = "capability_consensus"

// configTemplate defines the JSON template for consensus capability configuration.
// This allows overriding limits and other settings from capability_defaults.toml.
// If empty, the capability will use hardcoded defaults.
const configTemplate = `{}`

func (c *Consensus) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	jobsErr := createJobs(
		ctx,
		don,
		dons,
		creEnv,
	)
	if jobsErr != nil {
		return fmt.Errorf("failed to create OCR3 jobs: %w", jobsErr)
	}

	return nil
}

func createJobs(
	ctx context.Context,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	capabilityConfig, ok := don.GetCapabilityConfig(flag)
	if !ok {
		return fmt.Errorf("config for '%s' capability not found for %s DON", flag, don.GetName())
	}

	command, commandErr := standardcapability.GetCommand(capabilityConfig.BinaryName)
	if commandErr != nil {
		return fmt.Errorf("failed to get command for consensus capability: %w", commandErr)
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

	configStr, configErr := buildCapabilityConfig(
		flag,
		configTemplate,
		capabilityConfig,
	)
	if configErr != nil {
		return fmt.Errorf("failed to build consensus capability config: %w", configErr)
	}

	bootstrap, isBootstrap := dons.Bootstrap()
	if !isBootstrap {
		return errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
	}

	specs := make(map[string][]string)
	// Create node job
	if nodeSpecs, err := proposeNodeJob(creEnv, don, command, []string{formatBootstrapPeer(bootstrap)}, configStr); err != nil {
		return err
	} else if err := mergo.Merge(&specs, nodeSpecs); err != nil {
		return fmt.Errorf("failed to merge node job specs: %w", err)
	}

	if err := jobs.Approve(ctx, creEnv.CldfEnvironment.Offchain, dons, specs); err != nil {
		return fmt.Errorf("failed to approve Consensus jobs: %w", err)
	}

	return nil
}
