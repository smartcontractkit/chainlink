package v1

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"text/template"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	evmworkflow "github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	chainlinkbig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	tronchangeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/tron"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	corevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/evm"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const flag = cre.WriteEVMCapability

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
	gatewayJobConfigs map[cre.NodeUUID]*config.GatewayConfig,
) (*cre.PreEnvStartupOutput, error) {
	donsMetadata := topology.DonsMetadataWithFlag(flag)
	if len(donsMetadata) == 0 {
		return nil, nil
	}

	evmForwardersSelectors := make([]uint64, 0)
	tronForwardersSelectors := make([]uint64, 0)
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
					testLogger.Info().Msgf("Preparing Tron Keystone Forwarder deployment for chain %d", bcOut.ChainID)
					tronForwardersSelectors = append(tronForwardersSelectors, bcOut.ChainSelector)
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

	if len(tronForwardersSelectors) > 0 {
		deployErr := deployTronForwarders(testLogger, cldfEnv, tronForwardersSelectors, contractVersions)
		if deployErr != nil {
			return nil, errors.Wrap(deployErr, "failed to deploy Tron Keystone forwarder")
		}
	}

	// update node configs to include write-evm (evm v1) configuration
	for _, donMetadata := range donsMetadata {
		workerNodes, wErr := donMetadata.Workers()
		if wErr != nil {
			return nil, errors.Wrap(wErr, "failed to find worker nodes")
		}

		for _, workerNode := range workerNodes {
			writeEvmConfigs := []writeEVMData{}

			// for each worker node find all supported chains and node's public address for each chain
			for _, chainID := range donMetadata.CapabilitiesAwareNodeSet().ChainCapabilities[flag].EnabledChains {
				chain, exists := chain_selectors.ChainByEvmChainID(chainID)
				if !exists {
					return nil, errors.Errorf("failed to find selector for chain ID %d", chainID)
				}

				evmData := writeEVMData{
					ChainID:       chainID,
					ChainSelector: chain.Selector,
				}

				forwarderAddress, fErr := findForwarderAddress(chain, cldfEnv.ExistingAddresses) //nolint:staticcheck // won't migrate now
				if fErr != nil {
					return nil, errors.Errorf("failed to find forwarder address for chain %d", chain.Selector)
				}
				evmData.ForwarderAddress = forwarderAddress.Hex()

				evmKey, ok := workerNode.Keys.EVM[chainID]
				if !ok {
					return nil, fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
				}
				evmData.FromAddress = evmKey.PublicAddress

				var mergeErr error
				evmData, mergeErr = mergeDefaultAndRuntimeConfigValues(evmData, capabilityConfigs, donMetadata.CapabilitiesAwareNodeSet().ChainCapabilities, chainID)
				if mergeErr != nil {
					return nil, errors.Wrap(mergeErr, "failed to merge default and runtime write-evm config values")
				}

				writeEvmConfigs = append(writeEvmConfigs, evmData)
			}

			currentConfig := donMetadata.CapabilitiesAwareNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides

			var typedConfig corechainlink.Config
			unmarshallErr := toml.Unmarshal([]byte(currentConfig), &typedConfig)
			if unmarshallErr != nil {
				return nil, errors.Wrapf(unmarshallErr, "failed to unmarshal config for node index %d", workerNode.Index)
			}

			if len(typedConfig.EVM) < len(writeEvmConfigs) {
				return nil, fmt.Errorf("not enough EVM chains configured in node index %d to add write-evm (evm v1) config. Expected at least %d chains, but found %d", workerNode.Index, len(writeEvmConfigs), len(typedConfig.EVM))
			}

			for _, w := range writeEvmConfigs {
				chainFound := false
				for idx, evmChain := range typedConfig.EVM {
					chainIDIsEqual := evmChain.ChainID.Cmp(chainlinkbig.New(big.NewInt(libc.MustSafeInt64(w.ChainID)))) == 0
					if chainIDIsEqual {
						evmWorkflow, evmErr := buildEVMWorkflowConfig(w)
						if evmErr != nil {
							return nil, errors.Wrap(evmErr, "failed to build EVM workflow config")
						}

						typedConfig.EVM[idx].Workflow = *evmWorkflow
						typedConfig.EVM[idx].Transactions.ForwardersEnabled = ptr.Ptr(true)

						chainFound = true
						break
					}
				}

				if !chainFound {
					return nil, fmt.Errorf("failed to find EVM chain with ID %d in the config of node index %d to add write-evm config", w.ChainID, workerNode.Index)
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
			fullName := corevm.GenerateWriteTargetName(chainID)
			splitName := strings.Split(fullName, "@")

			capabilities[donMetadata.ID] = append(capabilities[donMetadata.ID], keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   splitName[0],
					Version:        splitName[1],
					CapabilityType: 3, // TARGET
					ResponseType:   1, // OBSERVATION_IDENTICAL
				},
				Config: &capabilitiespb.CapabilityConfig{},
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

	allAddresses, addrErr := creEnv.CldfEnvironment.ExistingAddresses.Addresses() //nolint:staticcheck // ignore SA1019 as ExistingAddresses is deprecated but still used
	if addrErr != nil {
		return errors.Wrap(addrErr, "failed to get addresses from address book")
	}

	evmChainsWithForwarders := make(map[uint64]struct{})
	tronChainsWithForwarders := make(map[uint64]struct{})
	for chainSelector, addresses := range allAddresses {
		for _, typeAndVersion := range addresses {
			if typeAndVersion.Type == keystone_changeset.KeystoneForwarder {
				for _, bcOut := range blockchainOutputs {
					if bcOut.ChainSelector == chainSelector {
						if !strings.EqualFold(bcOut.BlockchainOutput.Family, blockchain.FamilyTron) && !strings.EqualFold(bcOut.BlockchainOutput.Family, blockchain.FamilyEVM) {
							continue
						}

						if strings.EqualFold(bcOut.BlockchainOutput.Family, blockchain.FamilyTron) {
							tronChainsWithForwarders[chainSelector] = struct{}{}
						}

						evmChainsWithForwarders[chainSelector] = struct{}{}
						break
					}
				}
			}
		}
	}

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
	if len(evmChainsWithForwarders) > 0 {
		if evmErr := evm.ConfigureEVMForwarders(testLogger, creEnv.CldfEnvironment, evmChainsWithForwarders, consensusDON, consensusVersion); evmErr != nil {
			return errors.Wrap(evmErr, "failed to configure EVM forwarders")
		}
	}

	if len(tronChainsWithForwarders) > 0 {
		tErr := configureTronForwarders(testLogger, creEnv.CldfEnvironment, creEnv.DonTopology.HomeChainSelector, dons)
		if tErr != nil {
			return errors.Wrap(tErr, "failed to configure TRON forwarders")
		}
	}

	return nil
}

func deployTronForwarders(testLogger zerolog.Logger, cldfEnv *cldf.Environment, chainSelectors []uint64, contractVersions map[string]string) error {
	memoryDatastore := datastore.NewMemoryDataStore()

	// load all existing addresses into memory datastore
	mergeErr := memoryDatastore.Merge(cldfEnv.DataStore)
	if mergeErr != nil {
		return fmt.Errorf("failed to merge existing datastore into memory datastore: %w", mergeErr)
	}

	deployOptions := cldf_tron.DefaultDeployOptions()
	deployOptions.FeeLimit = 1_000_000_000

	deployChangeset := commonchangeset.Configure(tronchangeset.DeployForwarder{}, &tronchangeset.DeployForwarderRequest{
		ChainSelectors: chainSelectors,
		Qualifier:      "",
		DeployOptions:  deployOptions,
	})

	updatedEnv, err := commonchangeset.Apply(nil, *cldfEnv, deployChangeset)
	if err != nil {
		return fmt.Errorf("failed to deploy Tron forwarders using changesets: %w", err)
	}

	cldfEnv.ExistingAddresses = updatedEnv.ExistingAddresses //nolint:staticcheck // won't migrate now

	if updatedEnv.DataStore != nil {
		err = memoryDatastore.Merge(updatedEnv.DataStore)
		if err != nil {
			return fmt.Errorf("failed to merge updated datastore: %w", err)
		}
		cldfEnv.DataStore = memoryDatastore.Seal()

		for _, selector := range chainSelectors {
			forwarderAddr := contracts.MustGetAddressFromMemoryDataStore(memoryDatastore, selector, keystone_changeset.KeystoneForwarder.String(), contractVersions[keystone_changeset.KeystoneForwarder.String()], "")
			testLogger.Info().Msgf("Deployed Tron Forwarder %s contract on chain %d at %s", contractVersions[keystone_changeset.KeystoneForwarder.String()], selector, forwarderAddr)
		}
	}

	return nil
}

// func configureEVMForwarders(testLogger zerolog.Logger, cldfEnv *cldf.Environment, chainsWithForwarders map[uint64]struct{}, ocr3DON *cre.DON) error {
// 	forwarderCfg := forwarder.DonConfiguration{
// 		Name:    ocr3DON.Name,
// 		ID:      libc.MustSafeUint32FromUint64(ocr3DON.ID),
// 		F:       ocr3DON.F,
// 		Version: 1, // TODO this should be dynamic, but we don't have cap reg configured at this point
// 		NodeIDs: ocr3DON.KeystoneDONConfig().NodeIDs,
// 	}
// 	fout, err3 := operations.ExecuteSequence(
// 		cldfEnv.OperationsBundle,
// 		forwarder.ConfigureSeq,
// 		forwarder.ConfigureSeqDeps{
// 			Env: cldfEnv,
// 		},
// 		forwarder.ConfigureSeqInput{
// 			DON:    forwarderCfg,
// 			Chains: chainsWithForwarders,
// 		},
// 	)
// 	if err3 != nil {
// 		return errors.Wrap(err3, "failed to configure forwarders")
// 	}

// 	testLogger.Info().Msgf("Configured forwarders for v1 consensus: %+v", fout.Output.Config)

// 	return nil
// }

func configureTronForwarders(testLogger zerolog.Logger, env *cldf.Environment, registryChainSelector uint64, dons []*cre.DON) error {
	triggerOptions := cldf_tron.DefaultTriggerOptions()
	triggerOptions.FeeLimit = 1_000_000_000

	var wfNodeIDs []string
	for _, don := range dons {
		workerNodes, wErr := don.Workers()
		if wErr != nil {
			return fmt.Errorf("failed to find worker nodes for Tron configuration: %w", wErr)
		}

		for _, node := range workerNodes {
			wfNodeIDs = append(wfNodeIDs, node.Keys.P2PKey.PeerID.String())
		}
	}

	configChangeset := commonchangeset.Configure(tronchangeset.ConfigureForwarder{}, &tronchangeset.ConfigureForwarderRequest{
		WFDonName:        "workflow-don",
		WFNodeIDs:        wfNodeIDs,
		RegistryChainSel: registryChainSelector,
		Chains:           make(map[uint64]struct{}),
		TriggerOptions:   triggerOptions,
	})

	_, err := commonchangeset.Apply(nil, *env, configChangeset)
	if err != nil {
		return fmt.Errorf("failed to configure Tron forwarders using changesets: %w", err)
	}

	testLogger.Info().Msgf("Configured TRON forwarder for v1 consensus on chain: %d", registryChainSelector)

	return nil
}

func findForwarderAddress(chain chain_selectors.Chain, addressBook cldf.AddressBook) (*common.Address, error) {
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

func mergeDefaultAndRuntimeConfigValues(data writeEVMData, defaultCapabilityConfigs cre.CapabilityConfigs, nodeSetChainCapabilities map[string]*cre.ChainCapabilityConfig, chainID uint64) (writeEVMData, error) {
	if writeEvmConfig, ok := defaultCapabilityConfigs[flag]; ok {
		_, mergedConfig, rErr := envconfig.ResolveCapabilityForChain(
			flag,
			nodeSetChainCapabilities,
			writeEvmConfig.Config,
			chainID,
		)
		if rErr != nil {
			return data, errors.Wrapf(rErr, "failed to resolve write-evm config for chain %d", chainID)
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

func buildEVMWorkflowConfig(writeEVMInput writeEVMData) (*evmworkflow.Workflow, error) {
	var evmWorkflow evmworkflow.Workflow

	tmpl, tErr := template.New("evmWorkflowConfig").Parse(evmWorkflowConfigTemplate)
	if tErr != nil {
		return nil, errors.Wrap(tErr, "failed to parse evm workflow config template")
	}
	var configBuffer bytes.Buffer
	if executeErr := tmpl.Execute(&configBuffer, writeEVMInput.WorkflowConfig); executeErr != nil {
		return nil, errors.Wrap(executeErr, "failed to execute evm workflow config template")
	}

	configStr := configBuffer.String()
	if err := don.ValidateTemplateSubstitution(configStr, flag); err != nil {
		return nil, errors.Wrapf(err, "%s template validation failed", flag)
	}

	unmarshallErr := toml.Unmarshal([]byte(configStr), &evmWorkflow)
	if unmarshallErr != nil {
		return nil, errors.Wrapf(unmarshallErr, "failed to unmarshal EVM.Workflow config for chain %d", writeEVMInput.ChainID)
	}

	return &evmWorkflow, nil
}

type writeEVMData struct {
	ChainID          uint64
	ChainSelector    uint64
	FromAddress      common.Address
	ForwarderAddress string
	WorkflowConfig   map[string]any // Configuration for EVM.Workflow section
}

const evmWorkflowConfigTemplate = `
	FromAddress = '{{.FromAddress}}'
	ForwarderAddress = '{{.ForwarderAddress}}'
	GasLimitDefault = {{.GasLimitDefault}}
	TxAcceptanceState = {{.TxAcceptanceState}}
	PollPeriod = '{{.PollPeriod}}'
	AcceptanceTimeout = '{{.AcceptanceTimeout}}'
`
