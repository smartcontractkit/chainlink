package evm

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"
	"text/template"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"google.golang.org/protobuf/types/known/durationpb"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	evmworkflow "github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	chainlinkbig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/ocr"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/ocr/chainlevel"
)

const (
	flag                = cre.EVMCapability
	configTemplate      = `'{"chainId":{{.ChainID}},"network":"{{.NetworkFamily}}","logTriggerPollInterval":{{.LogTriggerPollInterval}}, "creForwarderAddress":"{{.CreForwarderAddress}}","receiverGasMinimum":{{.ReceiverGasMinimum}},"nodeAddress":"{{.NodeAddress}}"}'`
	registrationRefresh = 20 * time.Second
	registrationExpiry  = 60 * time.Second
	deltaStage          = 500*time.Millisecond + 1*time.Second // block time + 1 second delta
	requestTimeout      = 30 * time.Second
)

func New(registryChainID uint64) (*capabilities.Capability, error) {
	registryChainSelector, err := chainselectors.SelectorFromChainId(registryChainID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get selector from registry chainID: %d", registryChainID)
	}

	return capabilities.New(
		flag,
		capabilities.WithJobSpecFn(jobSpecWithRegistryChainSelector(registryChainSelector)),
		capabilities.WithCapabilityRegistryV1ConfigFn(registerWithV1),
		capabilities.WithNodeConfigTransformerFn(transformNodeConfig),
	)
}

func registerWithV1(_ []string, nodeSetInput *cre.CapabilitiesAwareNodeSet) ([]keystone_changeset.DONCapabilityWithConfig, error) {
	capabilities := make([]keystone_changeset.DONCapabilityWithConfig, 0)

	if nodeSetInput == nil {
		return nil, errors.New("node set input is nil")
	}

	// it's fine if there are no chain capabilities
	if nodeSetInput.ChainCapabilities == nil {
		return nil, nil
	}

	if _, ok := nodeSetInput.ChainCapabilities[flag]; !ok {
		return nil, nil
	}

	for _, chainID := range nodeSetInput.ChainCapabilities[flag].EnabledChains {
		selector, selectorErr := chainselectors.SelectorFromChainId(chainID)
		if selectorErr != nil {
			return nil, errors.Wrapf(selectorErr, "failed to get selector from chainID: %d", chainID)
		}

		evmMethodConfigs, err := getEvmMethodConfigs(nodeSetInput)
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

	return capabilities, nil
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

// buildRuntimeValues creates runtime-generated  values for any keys not specified in TOML
func buildRuntimeValues(chainID uint64, networkFamily, creForwarderAddress, nodeAddress string) map[string]any {
	return map[string]any{
		"ChainID":             chainID,
		"NetworkFamily":       networkFamily,
		"CreForwarderAddress": creForwarderAddress,
		"NodeAddress":         nodeAddress,
	}
}

func jobSpecWithRegistryChainSelector(registryChainSelector uint64) cre.JobSpecFn {
	return func(input *cre.JobSpecInput) (cre.DonsToJobSpecs, error) {
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
			creForwarderAddress, err := input.CldEnvironment.DataStore.Addresses().Get(creForwarderKey)
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

		return ocr.GenerateJobSpecsForStandardCapabilityWithOCR(
			input.DonTopology,
			input.CldEnvironment.DataStore,
			input.CapabilitiesAwareNodeSets,
			input.InfraInput,
			flag,
			contracts.CapabilityContractIdentifier,
			dataStoreOCR3ContractKeyProvider,
			chainlevel.CapabilityEnabler,
			chainlevel.EnabledChainsProvider,
			generateJobSpec,
			chainlevel.ConfigMerger,
			input.CapabilityConfigs,
		)
	}
}

type evmData struct {
	ChainID          uint64
	ChainSelector    uint64
	FromAddress      common.Address
	ForwarderAddress string
	WorkflowConfig   map[string]any // Configuration for EVM.Workflow section
}

// TODO PLEX-1732: refactor this method to not duplicate system-tests/lib/cre/capabilities/writeevm/write_evm.go, or guarantee it only looks for fromAddress to add it to the chain's workflow YAML element.
func transformNodeConfig(input cre.GenerateConfigsInput, existingConfigs cre.NodeIndexToConfigOverride) (cre.NodeIndexToConfigOverride, error) {
	if input.NodeSet == nil {
		return nil, errors.New("node set input is nil")
	}

	if input.NodeSet.ChainCapabilities == nil || input.NodeSet.ChainCapabilities[flag] == nil {
		return existingConfigs, nil
	}

	if input.CapabilityConfigs == nil {
		return nil, errors.New("additional capabilities configs are nil, but are required to configure the evm capability")
	}

	workflowNodeSet, wErr := node.FindManyWithLabel(input.DonMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.WorkerNode}, node.EqualLabels)
	if wErr != nil {
		return nil, errors.Wrap(wErr, "failed to find worker nodes")
	}

	for nodeIdx := range workflowNodeSet {
		var nodeIndex int
		for _, label := range workflowNodeSet[nodeIdx].Labels {
			if label.Key == node.IndexKey {
				var nErr error
				nodeIndex, nErr = strconv.Atoi(label.Value)
				if nErr != nil {
					return nil, errors.Wrap(nErr, "failed to convert node index to int")
				}
			}
		}

		// // get all the forwarders and add workflow config (FromAddress + Forwarder) for chains that have evm enabled
		data := []evmData{}
		for _, chainID := range input.NodeSet.ChainCapabilities[flag].EnabledChains {
			chain, exists := chainselectors.ChainByEvmChainID(chainID)
			if !exists {
				return nil, errors.Errorf("failed to find selector for chain ID %d", chainID)
			}

			evmDataValues := evmData{
				ChainID:       chainID,
				ChainSelector: chain.Selector,
			}

			forwarderAddress, fErr := findForwarderAddress(chain, input.AddressBook)
			if fErr != nil {
				return nil, errors.Errorf("failed to find forwarder address for chain %d", chain.Selector)
			}
			evmDataValues.ForwarderAddress = forwarderAddress.Hex()

			ethAddress, addrErr := findNodeEthAddressAddress(chain.Selector, workflowNodeSet[nodeIdx].Labels)
			if addrErr != nil {
				return nil, errors.Wrapf(addrErr, "failed to get ETH address for chain %d for node at index %d", chain.Selector, nodeIdx)
			}
			evmDataValues.FromAddress = *ethAddress

			var mergeErr error
			evmDataValues, mergeErr = mergeDefaultAndRuntimeConfigValues(evmDataValues, input.CapabilityConfigs, input.NodeSet.ChainCapabilities, chainID)
			if mergeErr != nil {
				return nil, errors.Wrap(mergeErr, "failed to merge default and runtime evm config values")
			}

			data = append(data, evmDataValues)
		}

		if len(existingConfigs) < nodeIndex+1 {
			return nil, errors.Errorf("missing config for node index %d", nodeIndex)
		}

		currentConfig := existingConfigs[nodeIndex]

		var typedConfig corechainlink.Config
		unmarshallErr := toml.Unmarshal([]byte(currentConfig), &typedConfig)
		if unmarshallErr != nil {
			return nil, errors.Wrapf(unmarshallErr, "failed to unmarshal config for node index %d", nodeIndex)
		}

		if len(typedConfig.EVM) < len(data) {
			return nil, fmt.Errorf("not enough EVM chains configured in node index %d to add evm config. Expected at least %d chains, but found %d", nodeIndex, len(data), len(typedConfig.EVM))
		}

		for _, evmInput := range data {
			chainFound := false
		INNER:
			for idx, evmChain := range typedConfig.EVM {
				chainIDIsEqual := evmChain.ChainID.Cmp(chainlinkbig.New(big.NewInt(libc.MustSafeInt64(evmInput.ChainID)))) == 0
				if chainIDIsEqual {
					evmWorkflow, evmErr := buildEVMWorkflowConfig(evmInput)
					if evmErr != nil {
						return nil, errors.Wrap(evmErr, "failed to build EVM workflow config")
					}

					typedConfig.EVM[idx].Workflow = *evmWorkflow
					typedConfig.EVM[idx].Transactions.ForwardersEnabled = ptr.Ptr(true)

					chainFound = true
					break INNER
				}
			}

			if !chainFound {
				return nil, fmt.Errorf("failed to find EVM chain with ID %d in the config of node index %d to add evm config", evmInput.ChainID, nodeIndex)
			}
		}

		stringifiedConfig, mErr := toml.Marshal(typedConfig)
		if mErr != nil {
			return nil, errors.Wrapf(mErr, "failed to marshal config for node index %d", nodeIndex)
		}

		existingConfigs[nodeIndex] = string(stringifiedConfig)
	}

	return existingConfigs, nil
}

func findForwarderAddress(chain chainselectors.Chain, addressBook deployment.AddressBook) (*common.Address, error) {
	addrsForChains, addErr := addressBook.AddressesForChain(chain.Selector)
	if addErr != nil {
		return nil, errors.Wrap(addErr, "failed to get addresses from address book")
	}

	for addr, addrValue := range addrsForChains {
		if addrValue.Type == keystone_changeset.KeystoneForwarder {
			return ptr.Ptr(common.HexToAddress(addr)), nil
		}
	}

	return nil, errors.Errorf("failed to find forwarder address for chain %d", chain.Selector)
}

func findNodeEthAddressAddress(chainSelector uint64, nodeLabels []*cre.Label) (*common.Address, error) {
	expectedAddressKey := node.AddressKeyFromSelector(chainSelector)
	for _, label := range nodeLabels {
		if label.Key == expectedAddressKey {
			if label.Value == "" {
				return nil, errors.Errorf("%s label value is empty", expectedAddressKey)
			}
			return ptr.Ptr(common.HexToAddress(label.Value)), nil
		}
	}

	return nil, errors.Errorf("failed to get from address for chain %d", chainSelector)
}

func mergeDefaultAndRuntimeConfigValues(data evmData, defaultCapabilityConfigs cre.CapabilityConfigs, nodeSetChainCapabilities map[string]*cre.ChainCapabilityConfig, chainID uint64) (evmData, error) {
	if evmConfig, ok := defaultCapabilityConfigs[flag]; ok {
		_, mergedConfig, rErr := envconfig.ResolveCapabilityForChain(
			cre.EVMCapability,
			nodeSetChainCapabilities,
			evmConfig.Config,
			chainID,
		)
		if rErr != nil {
			return data, errors.Wrapf(rErr, "failed to resolve evm config for chain %d", chainID)
		}

		runtimeValues := map[string]any{
			"FromAddress":      data.FromAddress.Hex(),
			"ForwarderAddress": data.ForwarderAddress,
		}

		var mErr error
		data.WorkflowConfig, mErr = don.ApplyRuntimeValues(mergedConfig, runtimeValues)
		if mErr != nil {
			return data, errors.Wrap(mErr, "failed to apply runtime values")
		}
	}

	return data, nil
}

const evmWorkflowConfigTemplate = `
	FromAddress = '{{.FromAddress}}'
	ForwarderAddress = '{{.ForwarderAddress}}'
	GasLimitDefault = {{.GasLimitDefault}}
	TxAcceptanceState = {{.TxAcceptanceState}}
	PollPeriod = '{{.PollPeriod}}'
	AcceptanceTimeout = '{{.AcceptanceTimeout}}'
`

func buildEVMWorkflowConfig(evmInput evmData) (*evmworkflow.Workflow, error) {
	var evmWorkflow evmworkflow.Workflow

	tmpl, tErr := template.New("evmWorkflowConfig").Parse(evmWorkflowConfigTemplate)
	if tErr != nil {
		return nil, errors.Wrap(tErr, "failed to parse evm workflow config template")
	}
	var configBuffer bytes.Buffer
	if executeErr := tmpl.Execute(&configBuffer, evmInput.WorkflowConfig); executeErr != nil {
		return nil, errors.Wrap(executeErr, "failed to execute evm workflow config template")
	}

	configStr := configBuffer.String()
	if err := don.ValidateTemplateSubstitution(configStr, flag); err != nil {
		return nil, errors.Wrapf(err, "%s template validation failed", flag)
	}

	unmarshallErr := toml.Unmarshal([]byte(configStr), &evmWorkflow)
	if unmarshallErr != nil {
		return nil, errors.Wrapf(unmarshallErr, "failed to unmarshal EVM.Workflow config for chain %d, err: %s", evmInput.ChainID, unmarshallErr.Error())
	}

	return &evmWorkflow, nil
}
