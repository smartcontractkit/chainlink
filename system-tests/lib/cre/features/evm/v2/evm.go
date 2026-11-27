package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"math"
	"slices"
	"sort"
	"strconv"
	"text/template"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/durationpb"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	credon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

const (
	flag                = cre.EVMCapability
	specConfigTemplate  = `{"chainId":{{printf "%d" .ChainID}}, "network":"{{.NetworkFamily}}", "logTriggerPollInterval":{{printf "%d" .LogTriggerPollInterval}}, "creForwarderAddress":"{{.CreForwarderAddress}}", "receiverGasMinimum":{{.ReceiverGasMinimum}}, "deltaStage":{{printf "%d" .DeltaStage}}{{with .LogTriggerSendChannelBufferSize}},"logTriggerSendChannelBufferSize":{{printf "%d" .}}{{end}}{{with .LogTriggerLimitQueryLogSize}},"logTriggerLimitQueryLogSize":{{printf "%d" .}}{{end}}}`
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
	chainsWithForwarders := chainsWithForwarders(creEnv.Blockchains, cre.ConvertToNodeSetWithChainCapabilities(topology.NodeSets()))
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
			deployErr := deployEVMForwarders(testLogger, creEnv.CldfEnvironment, selectorsToDeploy, creEnv.ContractVersions)
			if deployErr != nil {
				return nil, errors.Wrap(deployErr, "failed to deploy EVM Keystone forwarder")
			}
		}
	}

	enabledChainIDs, err := don.MustNodeSet().GetEnabledChainIDsForCapability(flag)
	if err != nil {
		return nil, fmt.Errorf("could not find enabled chainIDs for '%s' in don '%s': %w", flag, don.Name, err)
	}

	capabilities := []keystone_changeset.DONCapabilityWithConfig{}

	for _, chainID := range enabledChainIDs {
		selector, selectorErr := chainselectors.SelectorFromChainId(chainID)
		if selectorErr != nil {
			return nil, errors.Wrapf(selectorErr, "failed to get selector from chainID: %d", chainID)
		}

		evmMethodConfigs, err := getEvmMethodConfigs(don.MustNodeSet())
		if err != nil {
			return nil, errors.Wrap(err, "there was an error getting EVM method configs")
		}

		capabilityConfig, resolveErr := cre.ResolveCapabilityConfig(don.MustNodeSet(), flag, cre.ChainCapabilityScope(chainID))
		if resolveErr != nil {
			return nil, fmt.Errorf("could not resolve capability config for '%s' on chain %d: %w", flag, chainID, resolveErr)
		}

		creForwarderAddress, addrErr := creEnv.CldfEnvironment.DataStore.Addresses().Get(datastore.NewAddressRefKey(
			selector,
			datastore.ContractType(keystone_changeset.KeystoneForwarder.String()),
			semver.MustParse("1.0.0"),
			"",
		))
		if addrErr != nil {
			return nil, errors.Wrapf(addrErr, "failed to get CRE Forwarder address for chain %d", chainID)
		}

		specConfig, specErr := buildSpecConfig(chainID, creForwarderAddress.Address, capabilityConfig.Values)
		if specErr != nil {
			return nil, errors.Wrapf(specErr, "failed to build EVM spec config for chain %d", chainID)
		}

		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName: "evm" + ":ChainSelector:" + strconv.FormatUint(selector, 10),
				Version:      "1.0.0",
			},
			Config: &capabilitiespb.CapabilityConfig{
				MethodConfigs: evmMethodConfigs,
				LocalOnly:     don.HasOnlyLocalCapabilities(),
				SpecConfig:    specConfig,
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
	// The EVM capability is launched by the node from its registry spec_config
	// (see PreEnvStartup), so no worker job spec is proposed here.

	// configure EVM forwarders for DONs that run consensus
	consensusDons := dons.DonsWithFlags(cre.ConsensusCapability)

	// Forwarders may be configured when multiple DONs share chains with EVM capability; duplicate configuration is harmless.
	chainsWithEVMCapability := chainsWithEVMCapability(creEnv.Blockchains, dons.DonsWithFlag(flag))
	if len(chainsWithEVMCapability) > 0 {
		evmChainsWithForwarders := make([]uint64, 0)
		for _, chainSelector := range chainsWithEVMCapability {
			evmChainsWithForwarders = append(evmChainsWithForwarders, uint64(chainSelector))
		}
		for _, don := range consensusDons {
			config, confErr := configureEVMForwarders(testLogger, creEnv.CldfEnvironment, evmChainsWithForwarders, don)
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

// buildSpecConfig renders the EVM capability's spec_config for one chain. Every node of the DON
// shares it, so it carries no per-node values: nodeAddress was only a monitoring label and is
// filled in by the node when absent.
func buildSpecConfig(chainID uint64, creForwarderAddress string, capabilityValues map[string]any) (*valuespb.Map, error) {
	templateData, err := credon.ApplyRuntimeValues(maps.Clone(capabilityValues), buildRuntimeValues(chainID, "evm", creForwarderAddress))
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply runtime values")
	}

	tmpl, err := template.New("evmConfig").Parse(specConfigTemplate)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s spec config template", flag)
	}

	var configBuffer bytes.Buffer
	if execErr := tmpl.Execute(&configBuffer, templateData); execErr != nil {
		return nil, errors.Wrapf(execErr, "failed to execute %s spec config template", flag)
	}
	rendered := configBuffer.String()

	if validateErr := credon.ValidateTemplateSubstitution(rendered, flag); validateErr != nil {
		return nil, fmt.Errorf("%s template validation failed: %w\nRendered template: %s", flag, validateErr, rendered)
	}

	var cfg map[string]any
	if err = json.Unmarshal([]byte(rendered), &cfg); err != nil {
		return nil, errors.Wrapf(err, "failed to decode rendered %s spec config", flag)
	}
	for key, v := range cfg {
		// encoding/json decodes every number as float64; keep integral ones integral.
		if f, ok := v.(float64); ok && f == math.Trunc(f) {
			cfg[key] = int64(f)
		}
	}

	wrapped, err := values.NewMap(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to wrap %s spec config", flag)
	}
	return values.ProtoMap(wrapped), nil
}

// buildRuntimeValues creates runtime-generated  values for any keys not specified in TOML
func buildRuntimeValues(chainID uint64, networkFamily, creForwarderAddress string) map[string]any {
	return map[string]any{
		"ChainID":             chainID,
		"NetworkFamily":       networkFamily,
		"CreForwarderAddress": creForwarderAddress,
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

func deployEVMForwarders(testLogger zerolog.Logger, cldfEnv *cldf.Environment, chainSelectors []uint64, contractVersions map[cre.ContractType]*semver.Version) error {
	memoryDatastore, mErr := contracts.NewDataStoreFromExisting(cldfEnv.DataStore)
	if mErr != nil {
		return fmt.Errorf("failed to create memory datastore: %w", mErr)
	}

	evmForwardersReport, deployErr := operations.ExecuteSequence(
		cldfEnv.OperationsBundle,
		forwarder.DeploySequence,
		forwarder.DeploySequenceDeps{
			Env: cldfEnv,
		},
		forwarder.DeploySequenceInput{
			Targets: chainSelectors,
		},
	)
	if deployErr != nil {
		return errors.Wrap(deployErr, "failed to deploy evm forwarder")
	}

	if err := memoryDatastore.Merge(evmForwardersReport.Output.Datastore); err != nil {
		return errors.Wrap(err, "failed to merge datastore with Keystone contracts addresses")
	}

	for _, selector := range chainSelectors {
		forwarderAddr := contracts.MustGetAddressFromMemoryDataStore(memoryDatastore, selector, keystone_changeset.KeystoneForwarder.String(), contractVersions[keystone_changeset.KeystoneForwarder.String()], "")
		testLogger.Info().Msgf("Deployed EVM Forwarder %s contract on chain %d at %s", contractVersions[keystone_changeset.KeystoneForwarder.String()], selector, forwarderAddr)
	}

	cldfEnv.DataStore = memoryDatastore.Seal()

	return nil
}

func configureEVMForwarders(testLogger zerolog.Logger, cldfEnv *cldf.Environment, chainSelectors []uint64, ocr3DON *cre.Don) (*forwarder.Config, error) {
	forwarderCfg := forwarder.DonConfiguration{
		Name:    ocr3DON.Name,
		ID:      libc.MustSafeUint32FromUint64(ocr3DON.ID),
		F:       ocr3DON.F,
		Version: 1, // TODO this should be dynamic, but we don't have cap reg configured at this point, can we get that version from forwarder contract?
		NodeIDs: ocr3DON.KeystoneDONConfig().NodeIDs,
	}

	if len(chainSelectors) == 0 {
		for _, chain := range cldfEnv.BlockChains.EVMChains() {
			chainSelectors = append(chainSelectors, chain.Selector)
		}
	}

	chainsByQualifier := make(map[string]map[uint64]struct{})
	for _, selector := range chainSelectors {
		refs := cldfEnv.DataStore.Addresses().Filter(
			datastore.AddressRefByChainSelector(selector),
			datastore.AddressRefByType(datastore.ContractType(keystone_changeset.KeystoneForwarder.String())),
		)
		if len(refs) == 0 {
			return nil, fmt.Errorf("failed to resolve deployed forwarder for chain selector %d", selector)
		}

		for _, ref := range refs {
			if chainsByQualifier[ref.Qualifier] == nil {
				chainsByQualifier[ref.Qualifier] = make(map[uint64]struct{})
			}
			chainsByQualifier[ref.Qualifier][selector] = struct{}{}
		}
	}

	qualifiers := make([]string, 0, len(chainsByQualifier))
	for qualifier := range chainsByQualifier {
		qualifiers = append(qualifiers, qualifier)
	}
	sort.Strings(qualifiers)

	var configuredConfig forwarder.Config
	for _, qualifier := range qualifiers {
		fout, err := operations.ExecuteSequence(
			cldfEnv.OperationsBundle,
			forwarder.ConfigureSeq,
			forwarder.ConfigureSeqDeps{
				Env: cldfEnv,
			},
			forwarder.ConfigureSeqInput{
				DON:       forwarderCfg,
				Qualifier: qualifier,
				Chains:    chainsByQualifier[qualifier],
			},
		)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to configure forwarders with qualifier %q", qualifier)
		}
		configuredConfig = fout.Output.Config
	}

	return &configuredConfig, nil
}

func chainsWithForwarders(blockchains []blockchains.Blockchain, nodeSets []cre.NodeSetWithCapabilityConfigs) map[string][]uint64 {
	chainsWithForwarders := make(map[string][]uint64)

	for _, bcOut := range blockchains {
		for _, nodeSet := range nodeSets {
			if chainSelectors, familyExists := chainsWithForwarders[bcOut.ChainFamily()]; familyExists {
				if slices.Contains(chainSelectors, bcOut.ChainSelector()) {
					continue
				}
			}

			if !bcOut.IsFamily(chainselectors.FamilyEVM) && !bcOut.IsFamily(chainselectors.FamilyTron) {
				continue
			}

			if flags.RequiresForwarderContract(nodeSet.GetCapabilityFlags(), bcOut.ChainID()) {
				if _, exists := chainsWithForwarders[bcOut.ChainFamily()]; !exists {
					chainsWithForwarders[bcOut.ChainFamily()] = []uint64{}
				}
				chainsWithForwarders[bcOut.ChainFamily()] = append(chainsWithForwarders[bcOut.ChainFamily()], bcOut.ChainSelector())
			}
		}
	}

	return chainsWithForwarders
}
