package readcontract

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"dario.cat/mergo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	credon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

const flag = cre.ReadContractCapability

type ReadContract struct{}

func (o *ReadContract) Flag() cre.CapabilityFlag {
	return flag
}

func (o *ReadContract) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	capabilities := []keystone_changeset.DONCapabilityWithConfig{}
	for _, chainID := range don.NodeSets().GetChainCapabilityConfigs()[flag].EnabledChains {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   fmt.Sprintf("read-contract-evm-%d", chainID),
				Version:        "1.0.0",
				CapabilityType: 1, // ACTION
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

const configTemplate = `{"chainId":{{.ChainID}},"network":"{{.NetworkFamily}}"}`

func (o *ReadContract) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	specs := make(map[string][]string)

	capabilityConfig, ok := creEnv.CapabilityConfigs[flag]
	if !ok {
		return errors.Errorf("%s config not found in capabilities config. Make sure you have set it in the TOML config", flag)
	}

	command, cErr := standardcapability.GetCommand(capabilityConfig.BinaryPath, creEnv.Provider)
	if cErr != nil {
		return errors.Wrap(cErr, "failed to get command for Read Contract capability")
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

	chainCapConfig, ok := nodeSet.GetChainCapabilityConfigs()[flag]
	if !ok || chainCapConfig == nil {
		return fmt.Errorf("could not find chain capability config for '%s' in don '%s'", flag, don.Name)
	}

	for _, chainID := range chainCapConfig.EnabledChains {
		_, templateData, tErr := envconfig.ResolveCapabilityForChain(flag, nodeSet.GetChainCapabilityConfigs(), capabilityConfig.Config, chainID)
		if tErr != nil {
			return errors.Wrapf(tErr, "failed to resolve capability config for chain %d", chainID)
		}
		templateData["ChainID"] = chainID

		tmpl, tmplErr := template.New(flag + "-config").Parse(configTemplate)
		if tmplErr != nil {
			return errors.Wrapf(tmplErr, "failed to parse %s config template", flag)
		}

		var configBuffer bytes.Buffer
		if err := tmpl.Execute(&configBuffer, templateData); err != nil {
			return errors.Wrapf(err, "failed to execute %s config template", flag)
		}
		configStr := configBuffer.String()

		if err := credon.ValidateTemplateSubstitution(configStr, flag); err != nil {
			return errors.Wrapf(err, "%s template validation failed", flag)
		}

		workerInput := cre_jobs.ProposeJobSpecInput{
			Domain:      offchain.ProductLabel,
			Environment: cre.EnvironmentName,
			DONName:     don.Name,
			JobName:     fmt.Sprintf("read-contract-worker-%d", chainID),
			ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: don.Name},
			},
			Template: job_types.ReadContract,
			Inputs: job_types.JobSpecInput{
				"command": command,
				"config":  configStr,
			},
		}

		workerVerErr := cre_jobs.ProposeJobSpec{}.VerifyPreconditions(*creEnv.CldfEnvironment, workerInput)
		if workerVerErr != nil {
			return fmt.Errorf("precondition verification failed for Read Contract worker job: %w", workerVerErr)
		}

		workerReport, workerErr := cre_jobs.ProposeJobSpec{}.Apply(*creEnv.CldfEnvironment, workerInput)
		if workerErr != nil {
			return fmt.Errorf("failed to propose Read Contract worker job spec: %w", workerErr)
		}

		for _, r := range workerReport.Reports {
			out, ok := r.Output.(cre_jobs_ops.ProposeStandardCapabilityJobOutput)
			if !ok {
				return fmt.Errorf("unable to cast to ProposeStandardCapabilityJobOutput, actual type: %T", r.Output)
			}
			mErr := mergo.Merge(&specs, out.Specs, mergo.WithAppendSlice)
			if mErr != nil {
				return fmt.Errorf("failed to merge worker job specs: %w", mErr)
			}
		}
	}

	approveErr := jobs.Approve(ctx, creEnv.CldfEnvironment.Offchain, dons, specs)
	if approveErr != nil {
		return fmt.Errorf("failed to approve Read Contract jobs: %w", approveErr)
	}

	return nil
}
