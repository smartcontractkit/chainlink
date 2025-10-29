package v2

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"dario.cat/mergo"
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/durationpb"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	cre_jobs_pkg "github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	credon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/evm"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

const (
	flag                = cre.EVMCapability
	configTemplate      = `{"chainId":{{.ChainID}}, "network":"{{.NetworkFamily}}", "logTriggerPollInterval":{{.LogTriggerPollInterval}}, "creForwarderAddress":"{{.CreForwarderAddress}}", "receiverGasMinimum":{{.ReceiverGasMinimum}}, "nodeAddress":"{{.NodeAddress}}"{{with .LogTriggerSendChannelBufferSize}},"logTriggerSendChannelBufferSize":{{.}}{{end}}{{with .LogTriggerLimitQueryLogSize}},"logTriggerLimitQueryLogSize":{{.}}{{end}}}`
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
			// filter out EVM forwarder selectors that might have been already deployed by evm_v2 capability
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
		currentConfig := don.NodeSets().NodeSpecs[workerNode.Index].Node.TestConfigOverrides
		updatedConfig, updErr := updateNodeConfig(workerNode, don.NodeSets(), currentConfig)
		if updErr != nil {
			return nil, errors.Wrapf(updErr, "failed to update node config for node index %d", workerNode.Index)
		}

		don.NodeSets().NodeSpecs[workerNode.Index].Node.TestConfigOverrides = *updatedConfig
	}

	capabilities := []keystone_changeset.DONCapabilityWithConfig{}
	for _, chainID := range don.NodeSets().ChainCapabilities[flag].EnabledChains {
		selector, selectorErr := chainselectors.SelectorFromChainId(chainID)
		if selectorErr != nil {
			return nil, errors.Wrapf(selectorErr, "failed to get selector from chainID: %d", chainID)
		}

		evmMethodConfigs, err := getEvmMethodConfigs(don.NodeSets())
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
			},
		})
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

func updateNodeConfig(workerNode *cre.NodeMetadata, nodeSet *cre.NodeSet, currentConfig string) (*string, error) {
	chainsFromAddress, err := findNodeAddressPerChain(nodeSet, workerNode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get chains with from address")
	}

	var typedConfig corechainlink.Config
	unmarshallErr := toml.Unmarshal([]byte(currentConfig), &typedConfig)
	if unmarshallErr != nil {
		return nil, errors.Wrapf(unmarshallErr, "failed to unmarshal config for node index %d", workerNode.Index)
	}

	if len(typedConfig.EVM) < len(chainsFromAddress) {
		return nil, fmt.Errorf("not enough EVM chains configured in node index %d to add evm config. Expected at least %d chains, but found %d", workerNode.Index, len(chainsFromAddress), len(typedConfig.EVM))
	}

	for idx, evmChain := range typedConfig.EVM {
		chainID := libc.MustSafeUint64(evmChain.ChainID.Int64())
		addr, ok := chainsFromAddress[chainID]
		if ok {
			// if present means we need fromAddress for this chain
			address, err := types.NewEIP55Address(addr.Hex())
			if err != nil {
				return nil, errors.Wrapf(err, "failed to convert fromAddress to EIP55Address for chain %d", chainID)
			}
			typedConfig.EVM[idx].Workflow.FromAddress = &address
		}
	}

	stringifiedConfig, mErr := toml.Marshal(typedConfig)
	if mErr != nil {
		return nil, errors.Wrapf(mErr, "failed to marshal config for node index %d", workerNode.Index)
	}

	return ptr.Ptr(string(stringifiedConfig)), nil
}

func (o *EVM) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	chainsWithEVMCapability := chainsWithEVMCapability(creEnv.Blockchains, dons.DonsWithFlag(flag))
	for chainID, selector := range chainsWithEVMCapability {
		qualifier := ks_contracts_op.CapabilityContractIdentifier(uint64(chainID))
		_, _, seqErr := contracts.DeployOCR3Contract(testLogger, qualifier, creEnv.RegistryChainSelector, creEnv.CldfEnvironment, creEnv.ContractVersions)
		if seqErr != nil {
			return fmt.Errorf("failed to deploy EVM OCR3 contract for chainID %d, selector %d: %w", chainID, selector, seqErr)
		}
	}

	jobsErr := createJobs(
		ctx,
		don,
		dons,
		creEnv,
	)
	if jobsErr != nil {
		return jobsErr
	}

	// TODO should we make sure that log poller is listening before we try to configure contracts?

	// configure OCR3 contracts
	for chainID := range chainsWithEVMCapability {
		qualifier := ks_contracts_op.CapabilityContractIdentifier(uint64(chainID))
		// we have deployed OCR3 contract for each EVM chain on the registry chain to avoid a situation when more than 1 OCR contract (of any type) has the same address
		// because in past that violeted a DB constraint for offchain reporting jobs. Now there is no such limitation, but still it's better to have unique addresses to avoid confusion.
		evmOCR3Addr := contracts.MustGetAddressFromDataStore(creEnv.CldfEnvironment.DataStore, creEnv.RegistryChainSelector, keystone_changeset.OCR3Capability.String(), "1.0.0", qualifier)
		var evmDON *cre.Don
		for _, don := range dons.DonsWithFlag(cre.EVMCapability) {
			if flags.HasFlagForChain(don.Flags, cre.EVMCapability, uint64(chainID)) {
				evmDON = don
				break
			}
		}

		if evmDON == nil {
			return fmt.Errorf("failed to find DON for EVM chainID %d. This should never happen", chainID)
		}

		ocr3Config, ocr3confErr := contracts.DefaultChainCapabilityOCR3Config()
		if ocr3confErr != nil {
			return fmt.Errorf("failed to get default OCR3 config: %w", ocr3confErr)
		}

		_, err := operations.ExecuteOperation(
			creEnv.CldfEnvironment.OperationsBundle,
			ks_contracts_op.ConfigureOCR3Op,
			ks_contracts_op.ConfigureOCR3OpDeps{
				Env: creEnv.CldfEnvironment,
			},
			ks_contracts_op.ConfigureOCR3OpInput{
				ContractAddress: ptr.Ptr(common.HexToAddress(evmOCR3Addr)),
				ChainSelector:   creEnv.RegistryChainSelector,
				DON:             evmDON.KeystoneDONConfig(),
				Config:          evmDON.ResolveORC3Config(ocr3Config),
				DryRun:          false,
			},
		)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to configure EVM OCR3 contract for chainID: %d, address:%s", uint64(chainID), evmOCR3Addr))
		}
	}

	// configure EVM forwarders
	consensusDons := dons.DonsWithFlags(cre.ConsensusCapability, cre.ConsensusCapabilityV2)

	// for now we end up configuring forwarders twice, if the same chain has both evm v1 and v2 capabilities enabled
	// it doesn't create any issues, but ideally we wouldn't do that
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
	specs := make(map[string][]string)

	capabilityConfig, ok := creEnv.CapabilityConfigs[flag]
	if !ok {
		return fmt.Errorf("%s config not found in capabilities config: %v", flag, creEnv.CapabilityConfigs)
	}

	command, cErr := standardcapability.GetCommand(capabilityConfig.BinaryPath, creEnv.Provider)
	if cErr != nil {
		return errors.Wrap(cErr, "failed to get command for cron capability")
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

	bootstrap, isBootstrap := dons.Bootstrap()
	if !isBootstrap {
		return errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
	}

	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return errors.Wrap(wErr, "failed to find worker nodes")
	}

	chainConfig, ok := nodeSet.GetChainCapabilityConfigs()[flag]
	if !ok {
		return fmt.Errorf("could not find capability config for capability %s in node set %s", flag, nodeSet.GetName())
	}

	for _, chainID := range chainConfig.EnabledChains {
		chainSelector, selErr := chainselectors.SelectorFromChainId(chainID)
		if selErr != nil {
			return errors.Wrapf(selErr, "failed to get chain selector from chainID %d", chainID)
		}
		chainIDStr := strconv.FormatUint(chainID, 10)
		qualifier := ks_contracts_op.CapabilityContractIdentifier(chainID)

		ocr3Key := datastore.NewAddressRefKey(
			// we have deployed OCR3 contract for each EVM chain on the registry chain to avoid a situation when more than 1 OCR contract (of any type) has the same address
			// because that violates a DB constraint for offchain reporting jobs
			// this can be removed once https://smartcontract-it.atlassian.net/browse/PRODCRE-804 is done and we can deploy OCR3 contract for each EVM chain on that chain
			creEnv.RegistryChainSelector,
			datastore.ContractType(keystone_changeset.OCR3Capability.String()),
			semver.MustParse("1.0.0"),
			qualifier,
		)
		ocr3ConfigContractAddress, err := creEnv.CldfEnvironment.DataStore.Addresses().Get(ocr3Key)
		if err != nil {
			return errors.Wrapf(err, "failed contract address for key %s and chainID %d", ocr3Key, chainID)
		}

		bootInput := cre_jobs.ProposeJobSpecInput{
			Domain:      offchain.ProductLabel,
			Environment: cre.EnvironmentName,
			DONName:     bootstrap.DON.Name,
			JobName:     fmt.Sprintf("evm-v2-bootstrap-%d", chainID),
			ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: bootstrap.DON.Name},
			},
			Template: job_types.BootstrapOCR3,
			Inputs: job_types.JobSpecInput{
				"chainSelector":     chainSelector,
				"contractQualifier": qualifier,
			},
		}

		bootVerErr := cre_jobs.ProposeJobSpec{}.VerifyPreconditions(*creEnv.CldfEnvironment, bootInput)
		if bootVerErr != nil {
			return fmt.Errorf("precondition verification failed for EVM v2 bootstrap job for chainID %d: %w", chainID, bootVerErr)
		}

		bootReport, bootErr := cre_jobs.ProposeJobSpec{}.Apply(*creEnv.CldfEnvironment, bootInput)
		if bootErr != nil {
			return fmt.Errorf("failed to propose EVM v2 bootstrap job spec for chainID %d: %w", chainID, bootErr)
		}

		for _, r := range bootReport.Reports {
			out, ok := r.Output.(cre_jobs_ops.ProposeOCR3BootstrapJobOutput)
			if !ok {
				return fmt.Errorf("unable to cast to ProposeOCR3BootstrapJobOutput, actual type: %T", r.Output)
			}
			mErr := mergo.Merge(&specs, out.Specs, mergo.WithAppendSlice)
			if mErr != nil {
				return fmt.Errorf("failed to merge bootstrap job specs: %w", mErr)
			}
		}

		_, templateData, rErr := envconfig.ResolveCapabilityForChain(flag, nodeSet.GetChainCapabilityConfigs(), capabilityConfig.Config, chainID)
		if rErr != nil {
			return errors.Wrap(rErr, "failed to resolve capability config for chain")
		}

		for _, workerNode := range workerNodes {
			evmKey, ok := workerNode.Keys.EVM[chainID]
			if !ok {
				return fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
			}
			nodeAddress := evmKey.PublicAddress.Hex()

			creForwarderKey := datastore.NewAddressRefKey(
				chainSelector,
				datastore.ContractType(keystone_changeset.KeystoneForwarder.String()),
				semver.MustParse("1.0.0"),
				"",
			)
			creForwarderAddress, err := creEnv.CldfEnvironment.DataStore.Addresses().Get(creForwarderKey)
			if err != nil {
				return errors.Wrap(err, "failed to get CRE Forwarder address")
			}

			runtimeFallbacks := buildRuntimeValues(chainID, "evm", creForwarderAddress.Address, nodeAddress)
			var aErr error
			templateData, aErr = credon.ApplyRuntimeValues(templateData, runtimeFallbacks)
			if aErr != nil {
				return errors.Wrap(aErr, "failed to apply runtime values")
			}

			tmpl, err := template.New("evmConfig").Parse(configTemplate)
			if err != nil {
				return errors.Wrapf(err, "failed to parse %s config template", flag)
			}

			var configBuffer bytes.Buffer
			if err := tmpl.Execute(&configBuffer, templateData); err != nil {
				return errors.Wrapf(err, "failed to execute %s config template", flag)
			}

			configStr := configBuffer.String()

			if err := credon.ValidateTemplateSubstitution(configStr, flag); err != nil {
				return errors.Wrapf(err, "%s template validation failed", flag)
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

			workerInput := cre_jobs.ProposeJobSpecInput{
				Domain:      offchain.ProductLabel,
				Environment: cre.EnvironmentName,
				DONName:     don.Name,
				JobName:     fmt.Sprintf("evm-v2-worker-%d", chainID),
				ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
				DONFilters: []offchain.TargetDONFilter{
					{Key: offchain.FilterKeyDONName, Value: don.Name},
					{Key: "p2p_id", Value: workerNode.Keys.PeerID()}, // required since each node requires a different config (it contains its own from address)
				},
				Template: job_types.EVM,
				Inputs: job_types.JobSpecInput{
					"command": command,
					"config":  configStr,
					"oracleFactory": cre_jobs_pkg.OracleFactory{
						Enabled:            true,
						ChainID:            chainIDStr,
						BootstrapPeers:     bootstrapPeers,
						OCRContractAddress: ocr3ConfigContractAddress.Address,
						OCRKeyBundleID:     evmKeyBundle,
						TransmitterID:      nodeAddress,
						OnchainSigningStrategy: cre_jobs_pkg.OnchainSigningStrategy{
							StrategyName: strategyName,
							Config:       workerNode.Keys.OCR2BundleIDs,
						},
					},
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
	}

	approveErr := jobs.Approve(ctx, creEnv.CldfEnvironment.Offchain, dons, specs)
	if approveErr != nil {
		return fmt.Errorf("failed to approve EVM v2 jobs: %w", approveErr)
	}

	return nil
}

// func createJobs(
// 	ctx context.Context,
// 	don *cre.Don,
// 	dons *cre.Dons,
// 	creEnv *cre.Environment,
// ) error {
// 	jobSpecs := []*jobv1.ProposeJobRequest{}

// 	capabilityConfig, ok := creEnv.CapabilityConfigs[flag]
// 	if !ok {
// 		return fmt.Errorf("%s config not found in capabilities config: %v", flag, creEnv.CapabilityConfigs)
// 	}

// 	bootstrapNode, isBootstrap := dons.Bootstrap()
// 	if !isBootstrap {
// 		return errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
// 	}

// 	workerNodes, wErr := don.Workers()
// 	if wErr != nil {
// 		return errors.Wrap(wErr, "failed to find worker nodes")
// 	}

// 	var nodeSet cre.NodeSetWithCapabilityConfigs
// 	for _, ns := range dons.AsNodeSetWithChainCapabilities() {
// 		if ns.GetName() == don.Name {
// 			nodeSet = ns
// 			break
// 		}
// 	}
// 	if nodeSet == nil {
// 		return fmt.Errorf("could not find node set for Don named '%s'", don.Name)
// 	}

// 	command, cErr := standardcapability.GetCommand(capabilityConfig.BinaryPath, creEnv.Provider)
// 	if cErr != nil {
// 		return errors.Wrap(cErr, "failed to get command for cron capability")
// 	}

// 	chainConfig, ok := nodeSet.GetChainCapabilityConfigs()[flag]
// 	if !ok {
// 		return fmt.Errorf("could not find capability config for capability %s in node set %s", flag, nodeSet.GetName())
// 	}

// 	for _, chainID := range chainConfig.EnabledChains {
// 		chainIDStr := strconv.FormatUint(chainID, 10)
// 		chain, ok := chainselectors.ChainByEvmChainID(chainID)
// 		if !ok {
// 			return fmt.Errorf("failed to get chain selector for chain ID %d", chainID)
// 		}

// 		_, templateData, rErr := envconfig.ResolveCapabilityForChain(flag, nodeSet.GetChainCapabilityConfigs(), capabilityConfig.Config, chainID)
// 		if rErr != nil {
// 			return errors.Wrap(rErr, "failed to resolve capability config for chain")
// 		}

// 		ocr3Key := datastore.NewAddressRefKey(
// 			// we have deployed OCR3 contract for each EVM chain on the registry chain to avoid a situation when more than 1 OCR contract (of any type) has the same address
// 			// because that violates a DB constraint for offchain reporting jobs
// 			// this can be removed once https://smartcontract-it.atlassian.net/browse/PRODCRE-804 is done and we can deploy OCR3 contract for each EVM chain on that chain
// 			creEnv.RegistryChainSelector,
// 			datastore.ContractType(keystone_changeset.OCR3Capability.String()),
// 			semver.MustParse("1.0.0"),
// 			ks_contracts_op.CapabilityContractIdentifier(chainID),
// 		)
// 		ocr3ConfigContractAddress, err := creEnv.CldfEnvironment.DataStore.Addresses().Get(ocr3Key)
// 		if err != nil {
// 			return errors.Wrapf(err, "failed contract address for key %s and chainID %d", ocr3Key, chainID)
// 		}

// 		jobSpecs = append(jobSpecs, ocr.BootstrapJobSpec(bootstrapNode.JobDistributorDetails.NodeID, flag, ocr3ConfigContractAddress.Address, chainID))

// 		for _, workerNode := range workerNodes {
// 			evmKey, ok := workerNode.Keys.EVM[chainID]
// 			if !ok {
// 				return fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
// 			}
// 			nodeAddress := evmKey.PublicAddress.Hex()

// 			evmKeyBundle, ok := workerNode.Keys.OCR2BundleIDs[chainselectors.FamilyEVM] // we can always expect evm bundle key id present since evm is the registry chain
// 			if !ok {
// 				return errors.New("failed to get key bundle id for evm family")
// 			}

// 			bootstrapPeers := []string{fmt.Sprintf("%s@%s:%d", strings.TrimPrefix(bootstrapNode.Keys.PeerID(), "p2p_"), bootstrapNode.Host, cre.OCRPeeringPort)}

// 			strategyName := "single-chain"
// 			if len(workerNode.Keys.OCR2BundleIDs) > 1 {
// 				strategyName = "multi-chain"
// 			}

// 			oracleFactoryConfigInstance := job.OracleFactoryConfig{
// 				Enabled:            true,
// 				ChainID:            chainIDStr,
// 				BootstrapPeers:     bootstrapPeers,
// 				OCRContractAddress: ocr3ConfigContractAddress.Address,
// 				OCRKeyBundleID:     evmKeyBundle,
// 				TransmitterID:      nodeAddress,
// 				OnchainSigning: job.OnchainSigningStrategy{
// 					StrategyName: strategyName,
// 					Config:       workerNode.Keys.OCR2BundleIDs,
// 				},
// 			}

// 			type OracleFactoryConfigWrapper struct {
// 				OracleFactory job.OracleFactoryConfig `toml:"oracle_factory"`
// 			}
// 			wrapper := OracleFactoryConfigWrapper{OracleFactory: oracleFactoryConfigInstance}

// 			var oracleBuffer bytes.Buffer
// 			if errEncoder := toml.NewEncoder(&oracleBuffer).Encode(wrapper); errEncoder != nil {
// 				return errors.Wrap(errEncoder, "failed to encode oracle factory config to TOML")
// 			}
// 			oracleStr := strings.ReplaceAll(oracleBuffer.String(), "\n", "\n\t")

// 			creForwarderKey := datastore.NewAddressRefKey(
// 				chain.Selector,
// 				datastore.ContractType(keystone_changeset.KeystoneForwarder.String()),
// 				semver.MustParse("1.0.0"),
// 				"",
// 			)
// 			creForwarderAddress, err := creEnv.CldfEnvironment.DataStore.Addresses().Get(creForwarderKey)
// 			if err != nil {
// 				return errors.Wrap(err, "failed to get CRE Forwarder address")
// 			}

// 			runtimeFallbacks := buildRuntimeValues(chainID, "evm", creForwarderAddress.Address, nodeAddress)
// 			var aErr error
// 			templateData, aErr = credon.ApplyRuntimeValues(templateData, runtimeFallbacks)
// 			if aErr != nil {
// 				return errors.Wrap(aErr, "failed to apply runtime values")
// 			}

// 			tmpl, err := template.New("evmConfig").Parse(configTemplate)
// 			if err != nil {
// 				return errors.Wrapf(err, "failed to parse %s config template", flag)
// 			}

// 			var configBuffer bytes.Buffer
// 			if err := tmpl.Execute(&configBuffer, templateData); err != nil {
// 				return errors.Wrapf(err, "failed to execute %s config template", flag)
// 			}

// 			configStr := configBuffer.String()

// 			if err := credon.ValidateTemplateSubstitution(configStr, flag); err != nil {
// 				return errors.Wrapf(err, "%s template validation failed", flag)
// 			}

// 			jobSpec := standardcapability.WorkerJobSpec(workerNode.JobDistributorDetails.NodeID, fmt.Sprintf("%s-%d", flag, chainID), command, configStr, oracleStr)
// 			jobSpec.Labels = []*ptypes.Label{{Key: cre.CapabilityLabelKey, Value: ptr.Ptr(flag)}}
// 			jobSpecs = append(jobSpecs, jobSpec)
// 		}
// 	}

// 	jobErr := jobs.Create(ctx, creEnv.CldfEnvironment.Offchain, dons, jobSpecs)

// 	if jobErr != nil {
// 		return fmt.Errorf("failed to create EVM OCR3 jobs for don %s: %w", don.Name, jobErr)
// 	}

// 	return nil
// }

// buildRuntimeValues creates runtime-generated  values for any keys not specified in TOML
func buildRuntimeValues(chainID uint64, networkFamily, creForwarderAddress, nodeAddress string) map[string]any {
	return map[string]any{
		"ChainID":             chainID,
		"NetworkFamily":       networkFamily,
		"CreForwarderAddress": creForwarderAddress,
		"NodeAddress":         nodeAddress,
	}
}

func findNodeAddressPerChain(nodeSet *cre.NodeSet, workerNode *cre.NodeMetadata) (map[uint64]common.Address, error) {
	// get all the forwarders and add workflow config (FromAddress) for chains that have evm enabled
	data := make(map[uint64]common.Address)
	for _, chainID := range nodeSet.ChainCapabilities[flag].EnabledChains {
		evmKey, ok := workerNode.Keys.EVM[chainID]
		if !ok {
			return nil, fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
		}
		data[chainID] = evmKey.PublicAddress
	}

	return data, nil
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
				TransmissionSchedule:      capabilitiespb.TransmissionSchedule_OneAtATime,
				DeltaStage:                durationpb.New(deltaStage),
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
				TransmissionSchedule:      capabilitiespb.TransmissionSchedule_AllAtOnce,
				RequestTimeout:            durationpb.New(requestTimeout),
				ServerMaxParallelRequests: 10,
				RequestHasherType:         capabilitiespb.RequestHasherType_Simple,
			},
		},
	}
}
