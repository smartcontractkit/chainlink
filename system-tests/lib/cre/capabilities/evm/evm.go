package evm

import (
	"bytes"
	"fmt"
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
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
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
	configTemplate      = `'{"chainId":{{.ChainID}}, "network":"{{.NetworkFamily}}", "logTriggerPollInterval":{{.LogTriggerPollInterval}}, "creForwarderAddress":"{{.CreForwarderAddress}}", "receiverGasMinimum":{{.ReceiverGasMinimum}}, "nodeAddress":"{{.NodeAddress}}"{{with .LogTriggerSendChannelBufferSize}},"logTriggerSendChannelBufferSize":{{.}}{{end}}{{with .LogTriggerLimitQueryLogSize}},"logTriggerLimitQueryLogSize":{{.}}{{end}}}'`
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

// transformNodeConfig modifies the node config to add any required values for the evm capability
// specifically it adds the fromAddress for each chain that has evm enabled which will be used for the WriteReport method
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

	workerNodes, wErr := input.DonMetadata.Workers()
	if wErr != nil {
		return nil, errors.Wrap(wErr, "failed to find worker nodes")
	}

	for _, workerNode := range workerNodes {
		chainsFromAddress, err := findNodeAddressPerChain(input, workerNode)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get chains with from address")
		}

		if len(existingConfigs) < workerNode.Index+1 {
			return nil, errors.Errorf("missing config for node index %d", workerNode.Index)
		}

		currentConfig := existingConfigs[workerNode.Index]

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

		existingConfigs[workerNode.Index] = string(stringifiedConfig)
	}

	return existingConfigs, nil
}

func findNodeAddressPerChain(input cre.GenerateConfigsInput, workerNode *cre.NodeMetadata) (map[uint64]common.Address, error) {
	// get all the forwarders and add workflow config (FromAddress) for chains that have evm enabled
	data := make(map[uint64]common.Address)
	for _, chainID := range input.NodeSet.ChainCapabilities[flag].EnabledChains {
		evmKey, ok := workerNode.Keys.EVM[chainID]
		if !ok {
			return nil, fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
		}
		data[chainID] = evmKey.PublicAddress
	}

	return data, nil
}
