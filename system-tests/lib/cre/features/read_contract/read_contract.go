package readcontract

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	creblockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"

	"dario.cat/mergo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/durationpb"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	credon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
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
	capabilityToOCR3Config := map[string]*ocr3.OracleConfig{}
	enabledChainIDs, err := don.MustNodeSet().GetEnabledChainIDsForCapability(flag)
	if err != nil {
		return nil, fmt.Errorf("could not find enabled chainIDs for '%s' in don '%s': %w", flag, don.Name, err)
	}

	for _, chainID := range enabledChainIDs {
		// If this chain ID is an Aptos chain in the environment, register the Aptos view capability
		// (workflow expects aptos:ChainSelector:<selector>@1.0.0). Otherwise register EVM read-contract.
		var (
			labelledName       string
			useCapRegOCRConfig bool
			methodConfigs      map[string]*capabilitiespb.CapabilityMethodConfig
			localOnly          = don.HasOnlyLocalCapabilities()
		)
		if aptosChain := findAptosChainByChainID(creEnv.Blockchains, chainID); aptosChain != nil {
			// Version is provided separately in CapabilitiesRegistryCapability.Version below.
			// Keep LabelledName versionless to avoid creating IDs like "...@1.0.0@1.0.0".
			labelledName = "aptos:ChainSelector:" + strconv.FormatUint(aptosChain.ChainSelector(), 10)
			useCapRegOCRConfig = true
			methodConfigs = getAptosMethodConfigs()
			// Keep Aptos capability remote-enabled even in single-DON local CRE topologies.
			localOnly = false
			capabilityToOCR3Config[labelledName] = contracts.DefaultChainCapabilityOCR3Config()
		} else {
			labelledName = fmt.Sprintf("read-contract-evm-%d", chainID)
		}
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   labelledName,
				Version:        "1.0.0",
				CapabilityType: 1, // ACTION
			},
			Config: &capabilitiespb.CapabilityConfig{
				MethodConfigs: methodConfigs,
				LocalOnly:     localOnly,
			},
			UseCapRegOCRConfig: useCapRegOCRConfig,
		})
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
		CapabilityToOCR3Config:  capabilityToOCR3Config,
	}, nil
}

// findAptosChainByChainID returns an Aptos blockchain in the slice with the given chain ID, or nil.
func findAptosChainByChainID(blockchains []creblockchains.Blockchain, chainID uint64) creblockchains.Blockchain {
	for _, bc := range blockchains {
		if bc.IsFamily("aptos") && bc.ChainID() == chainID {
			return bc
		}
	}
	return nil
}

const configTemplate = `{"chainId":{{printf "%d" .ChainID}},"network":"{{.NetworkFamily}}"}`
const aptosConfigTemplate = `{"chainId":"{{.ChainID}}","network":"aptos","creForwarderAddress":"{{.CREForwarderAddress}}"}`
const aptosZeroForwarderHex = "0x0000000000000000000000000000000000000000000000000000000000000000"
const aptosRequestTimeout = 30 * time.Second

func getAptosMethodConfigs() map[string]*capabilitiespb.CapabilityMethodConfig {
	return map[string]*capabilitiespb.CapabilityMethodConfig{
		"View": {
			RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
				RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
					TransmissionSchedule:      capabilitiespb.TransmissionSchedule_AllAtOnce,
					RequestTimeout:            durationpb.New(aptosRequestTimeout),
					ServerMaxParallelRequests: 10,
					RequestHasherType:         capabilitiespb.RequestHasherType_Simple,
				},
			},
		},
	}
}

func (o *ReadContract) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	specs := make(map[string][]string)

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

	var aptosBootstrapPeers []string
	var capRegVersion string
	aptosOracleFactoryInputsReady := false

	for _, chainID := range enabledChainIDs {
		aptosChain := findAptosChainByChainID(creEnv.Blockchains, chainID)
		if aptosChain != nil {
			if !aptosOracleFactoryInputsReady {
				bootstrapNode, ok := dons.Bootstrap()
				if !ok {
					return errors.New("bootstrap node not found; required for Aptos OCR bootstrap peers")
				}
				aptosBootstrapPeers = []string{
					fmt.Sprintf("%s@%s:%d", strings.TrimPrefix(bootstrapNode.Keys.PeerID(), "p2p_"), bootstrapNode.Host, cre.OCRPeeringPort),
				}

				capRegSemver, ok := creEnv.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()]
				if !ok {
					return errors.New("CapabilitiesRegistry version not found in contract versions")
				}
				capRegVersion = capRegSemver.String()
				aptosOracleFactoryInputsReady = true
			}

			capabilityConfig, resolveErr := cre.ResolveCapabilityConfig(nodeSet, cre.WriteAptosCapability, cre.ChainCapabilityScope(chainID))
			if resolveErr != nil {
				return fmt.Errorf("could not resolve capability config for '%s' on chain %d: %w", cre.WriteAptosCapability, chainID, resolveErr)
			}
			command, cErr := standardcapability.GetCommand(capabilityConfig.BinaryName)
			if cErr != nil {
				return errors.Wrap(cErr, "failed to get command for Aptos capability")
			}

			forwarderAddress := aptosZeroForwarderHex
			if creEnv.AptosForwarderAddresses != nil {
				if a := creEnv.AptosForwarderAddresses[aptosChain.ChainSelector()]; a != "" {
					forwarderAddress = a
				}
			}

			tmpl, tmplErr := template.New("aptos-config").Parse(aptosConfigTemplate)
			if tmplErr != nil {
				return errors.Wrapf(tmplErr, "failed to parse Aptos config template")
			}

			templateData := map[string]string{
				"ChainID":             strconv.FormatUint(chainID, 10),
				"CREForwarderAddress": forwarderAddress,
			}

			var configBuffer bytes.Buffer
			if err := tmpl.Execute(&configBuffer, templateData); err != nil {
				return errors.Wrapf(err, "failed to execute Aptos config template")
			}
			configStr := configBuffer.String()
			if err := credon.ValidateTemplateSubstitution(configStr, string(cre.WriteAptosCapability)); err != nil {
				return fmt.Errorf("Aptos template validation failed: %w\nRendered: %s", err, configStr)
			}

			workerInput := cre_jobs.ProposeJobSpecInput{
				Domain:      offchain.ProductLabel,
				Environment: cre.EnvironmentName,
				DONName:     don.Name,
				JobName:     "aptos-worker-" + strconv.FormatUint(chainID, 10),
				ExtraLabels: map[string]string{cre.CapabilityLabelKey: string(cre.WriteAptosCapability)},
				DONFilters: []offchain.TargetDONFilter{
					{Key: offchain.FilterKeyDONName, Value: don.Name},
				},
				Template: job_types.Aptos,
				Inputs: job_types.JobSpecInput{
					"command":            command,
					"config":             configStr,
					"chainSelectorEVM":   creEnv.RegistryChainSelector,
					"chainSelectorAptos": aptosChain.ChainSelector(),
					"bootstrapPeers":     aptosBootstrapPeers,
					"useCapRegOCRConfig": true,
					"capRegVersion":      capRegVersion,
				},
			}

			workerVerErr := cre_jobs.ProposeJobSpec{}.VerifyPreconditions(*creEnv.CldfEnvironment, workerInput)
			if workerVerErr != nil {
				return fmt.Errorf("precondition verification failed for Aptos worker job: %w", workerVerErr)
			}
			workerReport, workerErr := cre_jobs.ProposeJobSpec{}.Apply(*creEnv.CldfEnvironment, workerInput)
			if workerErr != nil {
				return fmt.Errorf("failed to propose Aptos worker job spec: %w", workerErr)
			}

			for _, r := range workerReport.Reports {
				out, ok := r.Output.(cre_jobs_ops.ProposeStandardCapabilityJobOutput)
				if !ok {
					return fmt.Errorf("unable to cast to ProposeStandardCapabilityJobOutput, actual type: %T", r.Output)
				}
				mErr := mergo.Merge(&specs, out.Specs, mergo.WithAppendSlice)
				if mErr != nil {
					return fmt.Errorf("failed to merge Aptos worker job specs: %w", mErr)
				}
			}
			continue
		}
		capabilityConfig, resolveErr := cre.ResolveCapabilityConfig(nodeSet, flag, cre.ChainCapabilityScope(chainID))
		if resolveErr != nil {
			return fmt.Errorf("could not resolve capability config for '%s' on chain %d: %w", flag, chainID, resolveErr)
		}

		command, cErr := standardcapability.GetCommand(capabilityConfig.BinaryName)
		if cErr != nil {
			return errors.Wrap(cErr, "failed to get command for Read Contract capability")
		}

		templateData := capabilityConfig.Values
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
			return fmt.Errorf("%s template validation failed: %w\nRendered template: %s", flag, err, configStr)
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
