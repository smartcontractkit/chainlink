package v2

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"slices"
	"strconv"
	"strings"
	"time"

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
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	keystone_contracts "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/ocr"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/ocr/chainlevel"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/evm"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
)

const (
	flag                = cre.EVMCapability
	configTemplate      = `'{"chainId":{{.ChainID}}, "network":"{{.NetworkFamily}}", "logTriggerPollInterval":{{.LogTriggerPollInterval}}, "creForwarderAddress":"{{.CreForwarderAddress}}", "receiverGasMinimum":{{.ReceiverGasMinimum}}, "nodeAddress":"{{.NodeAddress}}"{{with .LogTriggerSendChannelBufferSize}},"logTriggerSendChannelBufferSize":{{.}}{{end}}{{with .LogTriggerLimitQueryLogSize}},"logTriggerLimitQueryLogSize":{{.}}{{end}}}'`
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
	testLogger zerolog.Logger,
	registryChainSelector uint64,
	cldfEnv *cldf.Environment,
	provider infra.Provider,
	topology *cre.Topology,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	capabilityConfigs cre.CapabilityConfigs,
	contractVersions map[string]string,
	gatewayConfigs map[cre.NodeUUID]*config.GatewayConfig,
) (*cre.PreEnvStartupOutput, error) {
	donsMetadata := topology.DonsMetadataWithFlag(flag)
	if len(donsMetadata) == 0 {
		return nil, nil
	}

	// deploy EVM forwarders if needed
	evmForwardersSelectors := make([]uint64, 0)
	for _, bcOut := range blockchainOutputs {
		for _, donMetadata := range topology.CapabilitiesAwareNodeSets() {
			if slices.Contains(evmForwardersSelectors, bcOut.ChainSelector) {
				continue
			}

			if !strings.EqualFold(bcOut.BlockchainOutput.Family, blockchain.FamilyEVM) && !strings.EqualFold(bcOut.BlockchainOutput.Family, blockchain.FamilyTron) {
				continue
			}

			if flags.RequiresForwarderContract(donMetadata.ComputedCapabilities, bcOut.ChainID) {
				if strings.EqualFold(bcOut.BlockchainOutput.Family, blockchain.FamilyTron) {
					continue
				} else {
					// deploy EVM forwarder only if not deployed yet (evm_v2 capability high have deployed it already)
					forwarderAddr := contracts.MightGetAddressFromDataStore(cldfEnv.DataStore, bcOut.ChainSelector, keystone_changeset.KeystoneForwarder.String(), contractVersions[keystone_changeset.KeystoneForwarder.String()], "")
					if forwarderAddr == nil {
						evmForwardersSelectors = append(evmForwardersSelectors, bcOut.ChainSelector)
					}
				}
			}
		}
	}

	if len(evmForwardersSelectors) > 0 {
		deployErr := evm.DeployEVMForwarders(testLogger, cldfEnv, evmForwardersSelectors, contractVersions)
		if deployErr != nil {
			return nil, errors.Wrap(deployErr, "failed to deploy EVM Keystone forwarder")
		}
	}

	// add node config for EVM.Workflow
	for _, donMetadata := range donsMetadata {
		workerNodes, wErr := donMetadata.Workers()
		if wErr != nil {
			return nil, errors.Wrap(wErr, "failed to find worker nodes")
		}
		for _, workerNode := range workerNodes {
			chainsFromAddress, err := findNodeAddressPerChain(donMetadata.CapabilitiesAwareNodeSet(), workerNode)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get chains with from address")
			}

			currentConfig := donMetadata.CapabilitiesAwareNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides

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

			donMetadata.CapabilitiesAwareNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides = string(stringifiedConfig)
		}
	}

	capabilities := make(map[uint64][]keystone_changeset.DONCapabilityWithConfig)
	for _, donMetadata := range donsMetadata {
		if capabilities[donMetadata.ID] == nil {
			capabilities[donMetadata.ID] = []keystone_changeset.DONCapabilityWithConfig{}
		}

		for _, chainID := range donMetadata.CapabilitiesAwareNodeSet().ChainCapabilities[flag].EnabledChains {
			selector, selectorErr := chainselectors.SelectorFromChainId(chainID)
			if selectorErr != nil {
				return nil, errors.Wrapf(selectorErr, "failed to get selector from chainID: %d", chainID)
			}

			evmMethodConfigs, err := getEvmMethodConfigs(donMetadata.CapabilitiesAwareNodeSet())
			if err != nil {
				return nil, errors.Wrap(err, "there was an error getting EVM method configs")
			}

			capabilities[donMetadata.ID] = append(capabilities[donMetadata.ID], keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName: "evm" + ":ChainSelector:" + strconv.FormatUint(selector, 10),
					Version:      "1.0.0",
				},
				Config: &capabilitiespb.CapabilityConfig{
					MethodConfigs: evmMethodConfigs,
				},
			})
		}
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfigs: capabilities,
	}, nil
}

func (o *EVM) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	creEnv *cre.Environment,
	nodeSetOutput []*cre.WrappedNodeOutput,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	contractVersions map[string]string,
	provider infra.Provider,
	capabilityConfigs map[string]cre.CapabilityConfig,
) error {
	dons := creEnv.DonTopology.DonsWithFlag(flag)
	if len(dons) == 0 {
		return nil
	}

	chainsWithEVMCapability := chainsWithEVMCapability(blockchainOutputs, dons)

	for chainID, selector := range chainsWithEVMCapability {
		qualifier := ks_contracts_op.CapabilityContractIdentifier(uint64(chainID))
		_, _, seqErr := contracts.DeployOCR3Contract(testLogger, qualifier, creEnv.DonTopology.HomeChainSelector, creEnv.CldfEnvironment, contractVersions)
		if seqErr != nil {
			return fmt.Errorf("failed to deploy EVM OCR3 contract for chainID %d, selector %d: %w", chainID, selector, seqErr)
		}
	}

	// create jobs
	jobsErr := createJobs(
		ctx,
		creEnv.CldfEnvironment,
		creEnv.DonTopology.HomeChainSelector,
		creEnv.DonTopology,
		provider,
		capabilityConfigs,
	)
	if jobsErr != nil {
		return jobsErr
	}

	// TODO should we make sure that log poller is listening before we try to configure contracts?

	// configure contracts
	for chainID := range chainsWithEVMCapability {
		qualifier := ks_contracts_op.CapabilityContractIdentifier(uint64(chainID))
		// we have deployed OCR3 contract for each EVM chain on the registry chain to avoid a situation when more than 1 OCR contract (of any type) has the same address
		// because that violates a DB constraint for offchain reporting jobs
		evmOCR3Addr := contracts.MustGetAddressFromDataStore(creEnv.CldfEnvironment.DataStore, creEnv.DonTopology.HomeChainSelector, keystone_changeset.OCR3Capability.String(), "1.0.0", qualifier)
		var evmDON *cre.DON
		for _, don := range creEnv.DonTopology.DonsWithFlag(cre.EVMCapability) {
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
				ChainSelector:   creEnv.DonTopology.HomeChainSelector,
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
	consensusVersion := "v1"
	consensusDON, oneErr := creEnv.DonTopology.OneDonWithFlag(cre.ConsensusCapability)
	if oneErr != nil {
		// if v1 consensus DON is not found, let's try v2. We should have exactly one DON with either v1 or v2 consensus
		consensusDON, oneErr = creEnv.DonTopology.OneDonWithFlag(cre.ConsensusCapabilityV2)
		consensusVersion = "v2"
		if oneErr != nil {
			return errors.New("failed to find DON with consensus v1 or v2 capability")
		}
	}

	// for now we end up configuring forwarders twice, if the same chain has both evm v1 and v2 capabilities enabled
	// it doesn't create any issues, but ideally we wouldn't do that
	if len(chainsWithEVMCapability) > 0 {
		evmChainsWithForwarders := make(map[uint64]struct{})
		for chainID := range chainsWithEVMCapability {
			evmChainsWithForwarders[uint64(chainID)] = struct{}{}
		}
		if evmErr := evm.ConfigureEVMForwarders(testLogger, creEnv.CldfEnvironment, evmChainsWithForwarders, consensusDON, consensusVersion); evmErr != nil {
			return errors.Wrap(evmErr, "failed to configure EVM forwarders")
		}
	}

	return nil
}

func chainsWithEVMCapability(chains []*cre.WrappedBlockchainOutput, dons []*cre.DON) map[ks_contracts_op.EVMChainID]ks_contracts_op.Selector {
	chainsWithEVMCapability := make(map[ks_contracts_op.EVMChainID]ks_contracts_op.Selector)
	for _, chain := range chains {
		for _, don := range dons {
			if flags.HasFlagForChain(don.Flags, cre.EVMCapability, chain.ChainID) {
				if chainsWithEVMCapability[ks_contracts_op.EVMChainID(chain.ChainID)] != 0 {
					continue
				}
				chainsWithEVMCapability[ks_contracts_op.EVMChainID(chain.ChainID)] = ks_contracts_op.Selector(chain.ChainSelector)
			}
		}
	}

	return chainsWithEVMCapability
}

func createJobs(
	ctx context.Context,
	cldfEnv *cldf.Environment,
	registryChainSelector uint64,
	donTopology *cre.DonTopology,
	provider infra.Provider,
	capabilityConfigs map[string]cre.CapabilityConfig,
) error {
	generateJobSpec := func(logger zerolog.Logger, chainID uint64, nodeAddress string, mergedConfig map[string]any) (string, error) {
		cs, ok := chainselectors.EvmChainIdToChainSelector()[chainID]
		if !ok {
			return "", fmt.Errorf("chain selector not found for chainID: %d", chainID)
		}

		creForwarderKey := datastore.NewAddressRefKey(
			cs,
			datastore.ContractType(keystone_changeset.KeystoneForwarder.String()),
			semver.MustParse("1.0.0"),
			"",
		)
		creForwarderAddress, err := cldfEnv.DataStore.Addresses().Get(creForwarderKey)
		if err != nil {
			return "", errors.Wrap(err, "failed to get CRE Forwarder address")
		}

		logger.Debug().Msgf("Found CRE Forwarder contract on chain %d at %s", chainID, creForwarderAddress.Address)

		runtimeFallbacks := buildRuntimeValues(chainID, "evm", creForwarderAddress.Address, nodeAddress)

		templateData, aErr := don.ApplyRuntimeValues(mergedConfig, runtimeFallbacks)
		if aErr != nil {
			return "", errors.Wrap(aErr, "failed to apply runtime values")
		}

		tmpl, err := template.New("evmConfig").Parse(configTemplate)
		if err != nil {
			return "", errors.Wrapf(err, "failed to parse %s config template", flag)
		}

		var configBuffer bytes.Buffer
		if err := tmpl.Execute(&configBuffer, templateData); err != nil {
			return "", errors.Wrapf(err, "failed to execute %s config template", flag)
		}

		configStr := configBuffer.String()

		if err := don.ValidateTemplateSubstitution(configStr, flag); err != nil {
			return "", errors.Wrapf(err, "%s template validation failed", flag)
		}

		return configStr, nil
	}

	dataStoreOCR3ContractKeyProvider := func(contractName string, _ uint64) datastore.AddressRefKey {
		return datastore.NewAddressRefKey(
			// we have deployed OCR3 contract for each EVM chain on the registry chain to avoid a situation when more than 1 OCR contract (of any type) has the same address
			// because that violates a DB constraint for offchain reporting jobs
			// this can be removed once https://smartcontract-it.atlassian.net/browse/PRODCRE-804 is done and we can deploy OCR3 contract for each EVM chain on that chain
			registryChainSelector,
			datastore.ContractType(keystone_changeset.OCR3Capability.String()),
			semver.MustParse("1.0.0"),
			contractName,
		)
	}

	donsToJobSpecs, jErr := ocr.GenerateJobSpecsForStandardCapabilityWithOCR(
		donTopology,
		cldfEnv.DataStore,
		donTopology.Dons.AsNodeSetWithChainCapabilities(),
		provider,
		flag,
		keystone_contracts.CapabilityContractIdentifier,
		dataStoreOCR3ContractKeyProvider,
		chainlevel.CapabilityEnabler,
		chainlevel.EnabledChainsProvider,
		generateJobSpec,
		chainlevel.ConfigMerger,
		capabilityConfigs,
	)
	if jErr != nil {
		return errors.Wrap(jErr, "failed to generate EVM OCR3 job specs")
	}

	for _, don := range donTopology.Dons.List() {
		jobSpecs, ok := donsToJobSpecs[don.ID]
		if !ok {
			continue
		}
		jobErr := jobs.Create(ctx, cldfEnv.Offchain, donTopology, jobSpecs)

		if jobErr != nil {
			return fmt.Errorf("failed to create EVM OCR3 jobs for don %s: %w", don.Name, jobErr)
		}
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
	}
}

func findNodeAddressPerChain(nodeSet *cre.CapabilitiesAwareNodeSet, workerNode *cre.NodeMetadata) (map[uint64]common.Address, error) {
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
func getEvmMethodConfigs(nodeSetInput *cre.CapabilitiesAwareNodeSet) (map[string]*capabilitiespb.CapabilityMethodConfig, error) {
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

	triggerConfig, err := logTriggerConfig(nodeSetInput)
	if err != nil {
		return nil, errors.Wrap(err, "failed get config for LogTrigger")
	}

	evmMethodConfigs["LogTrigger"] = triggerConfig
	evmMethodConfigs["WriteReport"] = writeReportActionConfig()
	return evmMethodConfigs, nil
}

func logTriggerConfig(nodeSetInput *cre.CapabilitiesAwareNodeSet) (*capabilitiespb.CapabilityMethodConfig, error) {
	faultyNodes, faultyErr := nodeSetInput.MaxFaultyNodes()
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
