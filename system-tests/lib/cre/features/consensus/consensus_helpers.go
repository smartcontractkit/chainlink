package consensus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"dario.cat/mergo"
	"github.com/pkg/errors"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	credon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
)

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

func proposeNodeJob(creEnv *cre.Environment, don *cre.Don, command string, bootstrapPeers []string, configStr string) (map[string][]string, error) {
	capRegVersion, ok := creEnv.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()]
	if !ok {
		return nil, errors.New("CapabilitiesRegistry version not found in contract versions")
	}

	inputs := job_types.JobSpecInput{
		"command":            command,
		"chainSelectorEVM":   creEnv.RegistryChainSelector,
		"bootstrapPeers":     bootstrapPeers,
		"useCapRegOCRConfig": true,
		"capRegVersion":      capRegVersion.String(),
	}

	// Add config if provided (allows overriding limits from capability_defaults.toml)
	if configStr != "" {
		inputs["config"] = configStr
	}

	// Add non-EVM OCR selectors when present so consensus can select the correct
	// offchain key bundle path for report generation.
	for _, blockchain := range creEnv.Blockchains {
		if blockchain.IsFamily(chainselectors.FamilyAptos) {
			inputs["chainSelectorAptos"] = blockchain.ChainSelector()
			continue
		}
		if blockchain.IsFamily(chainselectors.FamilySolana) {
			inputs["chainSelectorSolana"] = blockchain.ChainSelector()
			break
		}
	}

	input := cre_jobs.ProposeJobSpecInput{
		Domain:      offchain.ProductLabel,
		Environment: cre.EnvironmentName,
		DONName:     don.Name,
		JobName:     "consensus-worker",
		ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: don.Name},
		},
		Template: job_types.Consensus,
		Inputs:   inputs,
	}

	proposer := cre_jobs.ProposeJobSpec{}
	if verErr := proposer.VerifyPreconditions(*creEnv.CldfEnvironment, input); verErr != nil {
		return nil, fmt.Errorf("precondition verification failed for Consensus node job: %w", verErr)
	}

	report, applyErr := proposer.Apply(*creEnv.CldfEnvironment, input)
	if applyErr != nil {
		if strings.Contains(applyErr.Error(), "no aptos ocr2 config for node") {
			return nil, fmt.Errorf(
				"failed to propose Consensus node job spec: %w; Aptos workflows require Aptos OCR2 key bundles on all workflow DON nodes",
				applyErr,
			)
		}
		return nil, fmt.Errorf("failed to propose Consensus node job spec: %w", applyErr)
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
