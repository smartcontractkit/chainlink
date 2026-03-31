package readcontract

import (
	"bytes"
	"context"
	"fmt"
	"maps"
	"text/template"

	"dario.cat/mergo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	credon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	creblockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/jobhelpers"
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
	enabledChainIDs, err := don.MustNodeSet().GetEnabledChainIDsForCapability(flag)
	if err != nil {
		return nil, fmt.Errorf("could not find enabled chainIDs for '%s' in don '%s': %w", flag, don.Name, err)
	}

	for _, chainID := range enabledChainIDs {
		labelledName, labelErr := capabilityLabelForChain(creEnv, chainID)
		if labelErr != nil {
			return nil, labelErr
		}

		capConfig := &capabilitiespb.CapabilityConfig{
			LocalOnly: don.HasOnlyLocalCapabilities(),
		}

		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   labelledName,
				Version:        "1.0.0",
				CapabilityType: 1,
			},
			Config: capConfig,
		})
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

func capabilityLabelForChain(creEnv *cre.Environment, chainID uint64) (string, error) {
	for _, bc := range creEnv.Blockchains {
		if bc.ChainID() != chainID {
			continue
		}

		switch {
		case bc.IsFamily(blockchain.FamilyEVM), bc.IsFamily(blockchain.FamilyTron):
			return fmt.Sprintf("read-contract-evm-%d", chainID), nil
		default:
			return "", fmt.Errorf("read-contract is not supported for chain family %s on chainID %d", bc.ChainFamily(), chainID)
		}
	}

	return "", fmt.Errorf("could not find blockchain for read-contract chainID %d", chainID)
}

const configTemplate = `{"chainId":{{printf "%d" .ChainID}},"network":"{{.NetworkFamily}}"}`

func (o *ReadContract) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
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

	enabledChainIDs, err := nodeSet.GetEnabledChainIDsForCapability(flag)
	if err != nil {
		return fmt.Errorf("could not find enabled chainIDs for '%s' in don '%s': %w", flag, don.Name, err)
	}

	results := make([]map[string][]string, len(enabledChainIDs))
	group, groupCtx := errgroup.WithContext(ctx)
	group.SetLimit(jobhelpers.Parallelism(len(enabledChainIDs)))

	for i, chainID := range enabledChainIDs {
		group.Go(func() error {
			if _, findErr := findBlockchainByChainID(creEnv, chainID); findErr != nil {
				return findErr
			}

			capabilityConfig, resolveErr := cre.ResolveCapabilityConfig(nodeSet, flag, cre.ChainCapabilityScope(chainID))
			if resolveErr != nil {
				return fmt.Errorf("could not resolve capability config for '%s' on chain %d: %w", flag, chainID, resolveErr)
			}

			command, cErr := standardcapability.GetCommand(capabilityConfig.BinaryName)
			if cErr != nil {
				return errors.Wrap(cErr, "failed to get command for Read Contract capability")
			}

			templateData := maps.Clone(capabilityConfig.Values)
			templateData["ChainID"] = chainID

			tmpl, tmplErr := template.New(flag + "-config").Parse(configTemplate)
			if tmplErr != nil {
				return errors.Wrapf(tmplErr, "failed to parse %s config template", flag)
			}

			var configBuffer bytes.Buffer
			if executeErr := tmpl.Execute(&configBuffer, templateData); executeErr != nil {
				return errors.Wrapf(executeErr, "failed to execute %s config template", flag)
			}
			configStr := configBuffer.String()

			if validateErr := credon.ValidateTemplateSubstitution(configStr, flag); validateErr != nil {
				return fmt.Errorf("%s template validation failed: %w\nRendered template: %s", flag, validateErr, configStr)
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

			specs := make(map[string][]string)
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

			select {
			case <-groupCtx.Done():
				return groupCtx.Err()
			default:
			}

			results[i] = specs
			return nil
		})
	}

	if wErr := group.Wait(); wErr != nil {
		return wErr
	}

	specs, mErr := jobhelpers.MergeSpecsByIndex(results)
	if mErr != nil {
		return mErr
	}
	if len(specs) == 0 {
		return nil
	}

	approveErr := jobs.Approve(ctx, creEnv.CldfEnvironment.Offchain, dons, specs)
	if approveErr != nil {
		return fmt.Errorf("failed to approve Read Contract jobs: %w", approveErr)
	}

	return nil
}

func findBlockchainByChainID(creEnv *cre.Environment, chainID uint64) (creblockchains.Blockchain, error) {
	for _, bc := range creEnv.Blockchains {
		if bc.ChainID() == chainID {
			return bc, nil
		}
	}

	return nil, fmt.Errorf("could not find blockchain for read-contract chainID %d", chainID)
}
