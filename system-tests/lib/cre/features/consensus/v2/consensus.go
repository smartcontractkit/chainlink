package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"dario.cat/mergo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	credon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/consensus"
)

const flag = cre.ConsensusCapabilityV2

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
			LabelledName:   "consensus",
			Version:        "1.0.0-alpha",
			CapabilityType: 2, // CONSENSUS
			ResponseType:   0, // REPORT
		},
		Config: &capabilitiespb.CapabilityConfig{
			LocalOnly: don.HasOnlyLocalCapabilities(),
		},
	}}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
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
	_, ocr3ContractAddr, ocrErr := contracts.DeployOCR3Contract(testLogger, ContractQualifier, creEnv.RegistryChainSelector, creEnv.CldfEnvironment, creEnv.ContractVersions)
	if ocrErr != nil {
		return fmt.Errorf("failed to deploy OCR3 (consensus v2) contract: %w", ocrErr)
	}

	jobsErr := createJobs(
		ctx,
		don,
		dons,
		creEnv,
	)
	if jobsErr != nil {
		return fmt.Errorf("failed to create OCR3 jobs: %w", jobsErr)
	}

	// wait for LP to be started (otherwise it won't pick up contract's configuration events)
	if err := consensus.WaitForLogPollerToBeHealthy(don); err != nil {
		return fmt.Errorf("failed while waiting for Log Poller to become healthy: %w", err)
	}

	ocr3Config, ocr3confErr := contracts.DefaultOCR3Config()
	if ocr3confErr != nil {
		return fmt.Errorf("failed to get default OCR3 config: %w", ocr3confErr)
	}

	chain, ok := creEnv.CldfEnvironment.BlockChains.EVMChains()[creEnv.RegistryChainSelector]
	if !ok {
		return fmt.Errorf("chain with selector %d not found in environment", creEnv.RegistryChainSelector)
	}

	strategy, err := strategies.CreateStrategy(
		chain,
		*creEnv.CldfEnvironment,
		nil,
		nil,
		*ocr3ContractAddr,
		"PostEnvStartup - Configure OCR3 Contract - v2 Consensus",
	)
	if err != nil {
		return fmt.Errorf("failed to create strategy: %w", err)
	}

	_, ocr3Err := operations.ExecuteOperation(
		creEnv.CldfEnvironment.OperationsBundle,
		ks_contracts_op.ConfigureOCR3Op,
		ks_contracts_op.ConfigureOCR3OpDeps{
			Env:      creEnv.CldfEnvironment,
			Strategy: strategy,
		},
		ks_contracts_op.ConfigureOCR3OpInput{
			ContractAddress: ocr3ContractAddr,
			ChainSelector:   creEnv.RegistryChainSelector,
			DON:             don.KeystoneDONConfig(),
			Config:          don.ResolveORC3Config(ocr3Config),
			DryRun:          false,
		},
	)

	if ocr3Err != nil {
		return fmt.Errorf("failed to configure OCR3 contract: %w", ocr3Err)
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

	command, commandErr := standardcapability.GetCommand(capabilityConfig.BinaryPath, creEnv.Provider)
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

	// Create bootstrap job
	if bootSpecs, err := proposeBootstrapJob(creEnv, bootstrap, don); err != nil {
		return err
	} else if err := mergo.Merge(&specs, bootSpecs, mergo.WithAppendSlice); err != nil {
		return fmt.Errorf("failed to merge bootstrap job specs: %w", err)
	}

	// Create node job
	if nodeSpecs, err := proposeNodeJob(creEnv, don, command, []string{formatBootstrapPeer(bootstrap)}, configStr); err != nil {
		return err
	} else if err := mergo.Merge(&specs, nodeSpecs); err != nil {
		return fmt.Errorf("failed to merge node job specs: %w", err)
	}

	if err := jobs.Approve(ctx, creEnv.CldfEnvironment.Offchain, dons, specs); err != nil {
		return fmt.Errorf("failed to approve Consensus v2 jobs: %w", err)
	}

	return nil
}

// buildCapabilityConfig generates a JSON config string from template and config data.
// Returns empty string if no config is provided, allowing the capability to use defaults.
func buildCapabilityConfig(
	capabilityFlag cre.CapabilityFlag,
	configTemplate string,
	capConfig cre.CapabilityConfig,
) (string, error) {
	// Merge global defaults with DON-specific overrides
	templateData := capConfig.Values

	// If no config provided, return empty string (capability will use hardcoded defaults)
	if len(templateData) == 0 {
		return "", nil
	}

	// When template is "{}", marshal config map directly to JSON to pass through arbitrary fields from TOML
	if strings.TrimSpace(configTemplate) == "{}" {
		configBytes, err := json.Marshal(templateData)
		if err != nil {
			return "", errors.Wrapf(err, "failed to marshal %s config to JSON", capabilityFlag)
		}
		return string(configBytes), nil
	}

	// For template-based configs, parse and execute the template
	tmpl, err := template.New(capabilityFlag + "-config").Parse(configTemplate)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse %s config template", capabilityFlag)
	}

	var configBuffer bytes.Buffer
	if err := tmpl.Execute(&configBuffer, templateData); err != nil {
		return "", errors.Wrapf(err, "failed to execute %s config template", capabilityFlag)
	}
	configStr := configBuffer.String()

	if err := credon.ValidateTemplateSubstitution(configStr, capabilityFlag); err != nil {
		return "", errors.Wrapf(err, "%s config template validation failed", capabilityFlag)
	}

	return configStr, nil
}

func formatBootstrapPeer(bootstrap *cre.Node) string {
	return fmt.Sprintf("%s@%s:%d",
		strings.TrimPrefix(bootstrap.Keys.PeerID(), "p2p_"),
		bootstrap.Host,
		cre.OCRPeeringPort)
}

func proposeBootstrapJob(creEnv *cre.Environment, bootstrap *cre.Node, don *cre.Don) (map[string][]string, error) {
	input := cre_jobs.ProposeJobSpecInput{
		Domain:      offchain.ProductLabel,
		Environment: cre.EnvironmentName,
		DONName:     bootstrap.DON.Name,
		JobName:     "consensus-v2-bootstrap-" + don.Name,
		ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: bootstrap.DON.Name},
		},
		Template: job_types.BootstrapOCR3,
		Inputs: job_types.JobSpecInput{
			"chainSelector":     creEnv.RegistryChainSelector,
			"contractQualifier": ContractQualifier,
		},
	}

	proposer := cre_jobs.ProposeJobSpec{}
	if verErr := proposer.VerifyPreconditions(*creEnv.CldfEnvironment, input); verErr != nil {
		return nil, fmt.Errorf("precondition verification failed for Consensus v2 bootstrap job: %w", verErr)
	}

	report, applyErr := proposer.Apply(*creEnv.CldfEnvironment, input)
	if applyErr != nil {
		return nil, fmt.Errorf("failed to propose Consensus v2 bootstrap job spec: %w", applyErr)
	}

	specs := make(map[string][]string)
	for _, r := range report.Reports {
		out, ok := r.Output.(cre_jobs_ops.ProposeOCR3BootstrapJobOutput)
		if !ok {
			return nil, fmt.Errorf("unable to cast to ProposeOCR3BootstrapJobOutput, actual type: %T", r.Output)
		}
		if mergeErr := mergo.Merge(&specs, out.Specs, mergo.WithAppendSlice); mergeErr != nil {
			return nil, fmt.Errorf("failed to merge bootstrap job specs: %w", mergeErr)
		}
	}

	return specs, nil
}

func proposeNodeJob(creEnv *cre.Environment, don *cre.Don, command string, bootstrapPeers []string, configStr string) (map[string][]string, error) {
	inputs := job_types.JobSpecInput{
		"command":           command,
		"chainSelectorEVM":  creEnv.RegistryChainSelector,
		"contractQualifier": ContractQualifier,
		"bootstrapPeers":    bootstrapPeers,
	}

	// Add config if provided (allows overriding limits from capability_defaults.toml)
	if configStr != "" {
		inputs["config"] = configStr
	}

	// Add Solana chain selector if present
	for _, blockchain := range creEnv.Blockchains {
		if blockchain.IsFamily(chainselectors.FamilySolana) {
			inputs["chainSelectorSolana"] = blockchain.ChainSelector()
			break
		}
	}

	input := cre_jobs.ProposeJobSpecInput{
		Domain:      offchain.ProductLabel,
		Environment: cre.EnvironmentName,
		DONName:     don.Name,
		JobName:     "consensus-v2-worker",
		ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: don.Name},
		},
		Template: job_types.Consensus,
		Inputs:   inputs,
	}

	proposer := cre_jobs.ProposeJobSpec{}
	if verErr := proposer.VerifyPreconditions(*creEnv.CldfEnvironment, input); verErr != nil {
		return nil, fmt.Errorf("precondition verification failed for Consensus v2 node job: %w", verErr)
	}

	report, applyErr := proposer.Apply(*creEnv.CldfEnvironment, input)
	if applyErr != nil {
		return nil, fmt.Errorf("failed to propose Consensus v2 node job spec: %w", applyErr)
	}

	specs := make(map[string][]string)
	for _, r := range report.Reports {
		out, ok := r.Output.(cre_jobs_ops.ProposeStandardCapabilityJobOutput)
		if !ok {
			return nil, fmt.Errorf("unable to cast to ProposeStandardCapabilityJobOutput, actual type: %T", r.Output)
		}
		if mergeErr := mergo.Merge(&specs, out.Specs); mergeErr != nil {
			return nil, fmt.Errorf("failed to merge node job specs: %w", mergeErr)
		}
	}

	return specs, nil
}
