package v2

import (
	"bytes"
	"context"
	"fmt"
	"maps"
	"strconv"
	"strings"
	"text/template"
	"time"

	"dario.cat/mergo"
	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/durationpb"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	cre_jobs_pkg "github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	credon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/evm"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/jobhelpers"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

const (
	flag                = cre.EVMCapability
	configTemplate      = `{"chainId":{{printf "%d" .ChainID}}, "network":"{{.NetworkFamily}}", "logTriggerPollInterval":{{printf "%d" .LogTriggerPollInterval}}, "creForwarderAddress":"{{.CreForwarderAddress}}", "receiverGasMinimum":{{.ReceiverGasMinimum}}, "nodeAddress":"{{.NodeAddress}}", "deltaStage":{{printf "%d" .DeltaStage}}{{with .LogTriggerSendChannelBufferSize}},"logTriggerSendChannelBufferSize":{{printf "%d" .}}{{end}}{{with .LogTriggerLimitQueryLogSize}},"logTriggerLimitQueryLogSize":{{printf "%d" .}}{{end}}}`
	registrationRefresh = 20 * time.Second
	registrationExpiry  = 60 * time.Second
	deltaStage          = 500*time.Millisecond + 1*time.Second // block time + 1 second delta
	requestTimeout      = 30 * time.Second
)

type EVM struct{}

func (o *EVM) Flag() cre.CapabilityFlag {
	return flag
}

func (o *EVM) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	chainsWithForwarders := evm.ChainsWithForwarders(creEnv.Blockchains, cre.ConvertToNodeSetWithChainCapabilities(topology.NodeSets()))
	evmForwardersSelectors, exist := chainsWithForwarders[blockchain.FamilyEVM]

	if exist {
		selectorsToDeploy := make([]uint64, 0)
		for _, selector := range evmForwardersSelectors {
			// filter out EVM forwarder selectors that might have been already deployed by evm_v1 capability
			forwarderAddr := contracts.MightGetAddressFromDataStore(creEnv.CldfEnvironment.DataStore, selector, keystone_changeset.KeystoneForwarder.String(), creEnv.ContractVersions[keystone_changeset.KeystoneForwarder.String()], "")
			if forwarderAddr == nil {
				selectorsToDeploy = append(selectorsToDeploy, selector)
			}
		}

		if len(selectorsToDeploy) > 0 {
			deployErr := evm.DeployEVMForwarders(testLogger, creEnv.CldfEnvironment, selectorsToDeploy, creEnv.ContractVersions)
			if deployErr != nil {
				return nil, errors.Wrap(deployErr, "failed to deploy EVM Keystone forwarder")
			}
		}
	}

	// update node configs to include evm v2 configuration
	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return nil, errors.Wrap(wErr, "failed to find worker nodes")
	}
	for _, workerNode := range workerNodes {
		currentConfig := don.MustNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides
		currentConfigPtr := ptr.Ptr(currentConfig)
		don.MustNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides = *currentConfigPtr
	}

	capabilities := []keystone_changeset.DONCapabilityWithConfig{}

	enabledChainIDs, err := don.MustNodeSet().GetEnabledChainIDsForCapability(flag)
	if err != nil {
		return nil, fmt.Errorf("could not find enabled chainIDs for '%s' in don '%s': %w", flag, don.Name, err)
	}

	for _, chainID := range enabledChainIDs {
		selector, selectorErr := chainselectors.SelectorFromChainId(chainID)
		if selectorErr != nil {
			return nil, errors.Wrapf(selectorErr, "failed to get selector from chainID: %d", chainID)
		}

		evmMethodConfigs, err := getEvmMethodConfigs(don.MustNodeSet())
		if err != nil {
			return nil, errors.Wrap(err, "there was an error getting EVM method configs")
		}

		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName: "evm" + ":ChainSelector:" + strconv.FormatUint(selector, 10),
				Version:      "1.0.0",
			},
			Config: &capabilitiespb.CapabilityConfig{
				MethodConfigs: evmMethodConfigs,
				LocalOnly:     don.HasOnlyLocalCapabilities(),
			},
			UseCapRegOCRConfig: true,
		})
	}

	capabilityToOCR3Config := make(map[string]*ocr3.OracleConfig, len(capabilities))
	for _, cap := range capabilities {
		capabilityToOCR3Config[cap.Capability.LabelledName] = contracts.DefaultChainCapabilityOCR3Config()
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
		CapabilityToOCR3Config:  capabilityToOCR3Config,
	}, nil
}

func (o *EVM) PostEnvStartup(
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
		return jobsErr
	}

	// configure EVM forwarders
	consensusDons := dons.DonsWithFlags(cre.ConsensusCapability, cre.ConsensusCapabilityV2)

	// for now we end up configuring forwarders twice, if the same chain has both evm v1 and v2 capabilities enabled
	// it doesn't create any issues, but ideally we wouldn't do that
	chainsWithEVMCapability := chainsWithEVMCapability(creEnv.Blockchains, dons.DonsWithFlag(flag))
	if len(chainsWithEVMCapability) > 0 {
		evmChainsWithForwarders := make([]uint64, 0)
		for chainID := range chainsWithEVMCapability {
			evmChainsWithForwarders = append(evmChainsWithForwarders, uint64(chainID))
		}
		for _, don := range consensusDons {
			config, confErr := evm.ConfigureEVMForwarders(testLogger, creEnv.CldfEnvironment, evmChainsWithForwarders, don)
			if confErr != nil {
				return errors.Wrap(confErr, "failed to configure EVM forwarders")
			}
			testLogger.Info().Msgf("Configured EVM forwarders: %+v", config)
		}
	}

	return nil
}

func chainsWithEVMCapability(chains []blockchains.Blockchain, dons []*cre.Don) map[ks_contracts_op.EVMChainID]ks_contracts_op.Selector {
	chainsWithEVMCapability := make(map[ks_contracts_op.EVMChainID]ks_contracts_op.Selector)
	for _, chain := range chains {
		for _, don := range dons {
			if flags.HasFlagForChain(don.Flags, cre.EVMCapability, chain.ChainID()) {
				if chainsWithEVMCapability[ks_contracts_op.EVMChainID(chain.ChainID())] != 0 {
					continue
				}
				chainsWithEVMCapability[ks_contracts_op.EVMChainID(chain.ChainID())] = ks_contracts_op.Selector(chain.ChainSelector())
			}
		}
	}

	return chainsWithEVMCapability
}

func createJobs(
	ctx context.Context,
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

	bootstrap, isBootstrap := dons.Bootstrap()
	if !isBootstrap {
		return errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
	}

	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return errors.Wrap(wErr, "failed to find worker nodes")
	}

	enabledChainIDs, err := nodeSet.GetEnabledChainIDsForCapability(flag)
	if err != nil {
		return fmt.Errorf("could not find enabled chainIDs for '%s' in don '%s': %w", flag, don.Name, err)
	}

	registryChainID, rcErr := chainselectors.ChainIdFromSelector(creEnv.RegistryChainSelector)
	if rcErr != nil {
		return fmt.Errorf("failed to get chain ID from registry chain selector %d: %w", creEnv.RegistryChainSelector, rcErr)
	}

	type proposalWork struct {
		chainID          uint64
		chainIDStr       string
		chainSelector    uint64
		capabilityConfig cre.CapabilityConfig
		command          string
		workerNode       *cre.Node
	}

	workItems := make([]proposalWork, 0, len(enabledChainIDs)*len(workerNodes))
	for _, chainID := range enabledChainIDs {
		chainSelector, selErr := chainselectors.SelectorFromChainId(chainID)
		if selErr != nil {
			return errors.Wrapf(selErr, "failed to get chain selector from chainID %d", chainID)
		}
		chainIDStr := strconv.FormatUint(chainID, 10)

		capabilityConfig, resolveErr := cre.ResolveCapabilityConfig(nodeSet, flag, cre.ChainCapabilityScope(chainID))
		if resolveErr != nil {
			return fmt.Errorf("could not resolve capability config for '%s' on chain %d: %w", flag, chainID, resolveErr)
		}

		command, cErr := standardcapability.GetCommand(capabilityConfig.BinaryName)
		if cErr != nil {
			return errors.Wrap(cErr, "failed to get command for Read Contract capability")
		}

		for _, workerNode := range workerNodes {
			workItems = append(workItems, proposalWork{
				chainID:          chainID,
				chainIDStr:       chainIDStr,
				chainSelector:    chainSelector,
				capabilityConfig: capabilityConfig,
				command:          command,
				workerNode:       workerNode,
			})
		}
	}

	results := make([]map[string][]string, len(workItems))
	group, groupCtx := errgroup.WithContext(ctx)
	group.SetLimit(jobhelpers.Parallelism(len(workItems)))

	for i, workItem := range workItems {
		group.Go(func() error {
			chainID := workItem.chainID
			chainSelector := workItem.chainSelector
			workerNode := workItem.workerNode

			evmKey, ok := workerNode.Keys.EVM[chainID]
			if !ok {
				return fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
			}
			nodeAddress := evmKey.PublicAddress.Hex()

			evmRegistryKey, ok := workerNode.Keys.EVM[registryChainID]
			if !ok {
				return fmt.Errorf("failed to get registry EVM key (chainID %d, node index %d) enabledChainIDs: %v", registryChainID, workerNode.Index, enabledChainIDs)
			}
			nodeRegistryAddress := evmRegistryKey.PublicAddress.Hex()

			creForwarderKey := datastore.NewAddressRefKey(
				chainSelector,
				datastore.ContractType(keystone_changeset.KeystoneForwarder.String()),
				semver.MustParse("1.0.0"),
				"",
			)
			creForwarderAddress, cErr := creEnv.CldfEnvironment.DataStore.Addresses().Get(creForwarderKey)
			if cErr != nil {
				return errors.Wrap(cErr, "failed to get CRE Forwarder address")
			}

			runtimeFallbacks := buildRuntimeValues(chainID, "evm", creForwarderAddress.Address, nodeAddress)
			templateData := maps.Clone(workItem.capabilityConfig.Values)

			var aErr error
			templateData, aErr = credon.ApplyRuntimeValues(templateData, runtimeFallbacks)
			if aErr != nil {
				return errors.Wrap(aErr, "failed to apply runtime values")
			}

			tmpl, tErr := template.New("evmConfig").Parse(configTemplate)
			if tErr != nil {
				return errors.Wrapf(tErr, "failed to parse %s config template", flag)
			}

			var configBuffer bytes.Buffer
			if execErr := tmpl.Execute(&configBuffer, templateData); execErr != nil {
				return errors.Wrapf(execErr, "failed to execute %s config template", flag)
			}

			configStr := configBuffer.String()

			if validateErr := credon.ValidateTemplateSubstitution(configStr, flag); validateErr != nil {
				return fmt.Errorf("%s template validation failed: %w\nRendered template: %s", flag, validateErr, configStr)
			}

			evmKeyBundle, ok := workerNode.Keys.OCR2BundleIDs[chainselectors.FamilyEVM] // we can always expect evm bundle key id present since evm is the registry chain
			if !ok {
				return errors.New("failed to get key bundle id for evm family")
			}

			bootstrapPeers := []string{fmt.Sprintf("%s@%s:%d", strings.TrimPrefix(bootstrap.Keys.PeerID(), "p2p_"), bootstrap.Host, cre.OCRPeeringPort)}

			strategyName := "single-chain"
			if len(workerNode.Keys.OCR2BundleIDs) > 1 {
				strategyName = "multi-chain"
			}

			capRegVersion, ok := creEnv.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()]
			if !ok {
				return errors.New("CapabilitiesRegistry version not found in contract versions")
			}
			registryAddrRefKey := cre_jobs_pkg.GetCapRegAddressRefKey(creEnv.RegistryChainSelector, "", capRegVersion.String())
			registryContractAddrRef, err := creEnv.CldfEnvironment.DataStore.Addresses().Get(registryAddrRefKey)
			if err != nil {
				return fmt.Errorf("failed to get contract address for ref key %s: %w", registryAddrRefKey, err)
			}

			workerInput := cre_jobs.ProposeJobSpecInput{
				Domain:      offchain.ProductLabel,
				Environment: cre.EnvironmentName,
				DONName:     don.Name,
				JobName:     fmt.Sprintf("evm-capabilities-v2-%d", chainID),
				ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
				DONFilters: []offchain.TargetDONFilter{
					{Key: offchain.FilterKeyDONName, Value: don.Name},
					{Key: "p2p_id", Value: workerNode.Keys.PeerID()}, // required since each node requires a different config (it contains its own from address)
				},
				Template: job_types.EVM,
				Inputs: job_types.JobSpecInput{
					"command": workItem.command,
					"config":  configStr,
					"oracleFactory": cre_jobs_pkg.OracleFactory{
						Enabled:            true,
						ChainID:            workItem.chainIDStr,
						BootstrapPeers:     bootstrapPeers,
						OCRContractAddress: registryContractAddrRef.Address,
						OCRKeyBundleID:     evmKeyBundle,
						TransmitterID:      nodeRegistryAddress,
						OnchainSigningStrategy: cre_jobs_pkg.OnchainSigningStrategy{
							StrategyName: strategyName,
							Config:       workerNode.Keys.OCR2BundleIDs,
						},
					},
					"useCapRegOCRConfig": true,
					"capRegVersion":      capRegVersion.String(),
				},
			}

			workerVerErr := cre_jobs.ProposeJobSpec{}.VerifyPreconditions(*creEnv.CldfEnvironment, workerInput)
			if workerVerErr != nil {
				return fmt.Errorf("precondition verification failed for EVM v2 worker job: %w", workerVerErr)
			}

			workerReport, workerErr := cre_jobs.ProposeJobSpec{}.Apply(*creEnv.CldfEnvironment, workerInput)
			if workerErr != nil {
				return fmt.Errorf("failed to propose EVM v2 worker job spec: %w", workerErr)
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

	approveErr := jobs.Approve(ctx, creEnv.CldfEnvironment.Offchain, dons, specs)
	if approveErr != nil {
		return fmt.Errorf("failed to approve EVM v2 jobs: %w", approveErr)
	}

	return nil
}

// buildRuntimeValues creates runtime-generated  values for any keys not specified in TOML
func buildRuntimeValues(chainID uint64, networkFamily, creForwarderAddress, nodeAddress string) map[string]any {
	return map[string]any{
		"ChainID":             chainID,
		"NetworkFamily":       networkFamily,
		"CreForwarderAddress": creForwarderAddress,
		"NodeAddress":         nodeAddress,
		"DeltaStage":          deltaStage,
	}
}

// getEvmMethodConfigs returns the method configs for all EVM methods we want to support, if any method is missing it
// will not be reached by the node when running evm capability in remote don
func getEvmMethodConfigs(nodeSet *cre.NodeSet) (map[string]*capabilitiespb.CapabilityMethodConfig, error) {
	evmMethodConfigs := map[string]*capabilitiespb.CapabilityMethodConfig{}

	// the read actions should be all defined in the proto that are neither a LogTrigger type, not a WriteReport type
	// see the RPC methods to map here: https://github.com/smartcontractkit/chainlink-protos/blob/main/cre/capabilities/blockchain/evm/v1alpha/client.proto
	readActions := []string{
		"CallContract",
		"FilterLogs",
		"BalanceAt",
		"EstimateGas",
		"GetTransactionByHash",
		"GetTransactionReceipt",
		"HeaderByNumber",
	}
	for _, action := range readActions {
		evmMethodConfigs[action] = readActionConfig()
	}

	triggerConfig, err := logTriggerConfig(nodeSet)
	if err != nil {
		return nil, errors.Wrap(err, "failed get config for LogTrigger")
	}

	evmMethodConfigs["LogTrigger"] = triggerConfig
	evmMethodConfigs["WriteReport"] = writeReportActionConfig()
	return evmMethodConfigs, nil
}

func logTriggerConfig(nodeSet *cre.NodeSet) (*capabilitiespb.CapabilityMethodConfig, error) {
	faultyNodes, faultyErr := nodeSet.MaxFaultyNodes()
	if faultyErr != nil {
		return nil, errors.Wrap(faultyErr, "failed to get faulty nodes")
	}

	return &capabilitiespb.CapabilityMethodConfig{
		RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh:     durationpb.New(registrationRefresh),
				RegistrationExpiry:      durationpb.New(registrationExpiry),
				MinResponsesToAggregate: faultyNodes + 1,
				MessageExpiry:           durationpb.New(2 * registrationExpiry),
				MaxBatchSize:            25,
				BatchCollectionPeriod:   durationpb.New(200 * time.Millisecond),
			},
		},
	}, nil
}

func writeReportActionConfig() *capabilitiespb.CapabilityMethodConfig {
	return &capabilitiespb.CapabilityMethodConfig{
		RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
			RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
				RequestTimeout:            durationpb.New(requestTimeout),
				ServerMaxParallelRequests: 10,
				RequestHasherType:         capabilitiespb.RequestHasherType_WriteReportExcludeSignatures,
			},
		},
	}
}

func readActionConfig() *capabilitiespb.CapabilityMethodConfig {
	return &capabilitiespb.CapabilityMethodConfig{
		RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
			RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
				RequestTimeout:            durationpb.New(requestTimeout),
				ServerMaxParallelRequests: 10,
				RequestHasherType:         capabilitiespb.RequestHasherType_Simple,
			},
		},
	}
}
