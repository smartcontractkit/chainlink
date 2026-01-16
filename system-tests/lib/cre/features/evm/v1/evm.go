package v1

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"strings"
	"text/template"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	evmworkflow "github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	chainlinkbig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	tronchangeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/tron"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	corevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	evmblockchain "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/evm"
)

const flag = cre.WriteEVMCapability

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

	tronForwardersSelectors, exist := chainsWithForwarders[blockchain.FamilyTron]
	if exist {
		deployErr := deployTronForwarders(testLogger, creEnv.CldfEnvironment, tronForwardersSelectors, creEnv.ContractVersions)
		if deployErr != nil {
			return nil, errors.Wrap(deployErr, "failed to deploy Tron Keystone forwarder")
		}
	}

	// update node configs to include write-evm (evm v1) configuration
	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return nil, errors.Wrap(wErr, "failed to find worker nodes")
	}

	for _, workerNode := range workerNodes {
		currentConfig := don.NodeSets().NodeSpecs[workerNode.Index].Node.TestConfigOverrides
		updatedConfig, updErr := updateNodeConfig(workerNode, currentConfig, don.NodeSets().ChainCapabilities, creEnv)
		if updErr != nil {
			return nil, errors.Wrapf(updErr, "failed to update node config for node index %d", workerNode.Index)
		}

		don.NodeSets().NodeSpecs[workerNode.Index].Node.TestConfigOverrides = *updatedConfig
	}

	capabilities := []keystone_changeset.DONCapabilityWithConfig{}
	for _, chainID := range don.NodeSets().ChainCapabilities[flag].EnabledChains {
		fullName := corevm.GenerateWriteTargetName(chainID)
		splitName := strings.Split(fullName, "@")

		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   splitName[0],
				Version:        splitName[1],
				CapabilityType: 3, // TARGET
				ResponseType:   1, // OBSERVATION_IDENTICAL
			},
			Config: &capabilitiespb.CapabilityConfig{
				LocalOnly: don.HasOnlyLocalCapabilities(),
			},
		})
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

func (o *EVM) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	consensusDons := dons.DonsWithFlags(cre.ConsensusCapability, cre.ConsensusCapabilityV2)
	chainsWithForwarders := evm.ChainsWithForwarders(creEnv.Blockchains, dons.AsNodeSetWithChainCapabilities())

	// for now we end up configuring forwarders twice, if the same chain has both evm v1 and v2 capabilities enabled
	// it doesn't create any issues, but ideally we wouldn't do that
	evmForwardersSelectors, exist := chainsWithForwarders[blockchain.FamilyEVM]
	if exist {
		for _, don := range consensusDons {
			config, confErr := evm.ConfigureEVMForwarders(testLogger, creEnv.CldfEnvironment, evmForwardersSelectors, don)
			if confErr != nil {
				return errors.Wrap(confErr, "failed to configure EVM forwarders")
			}
			testLogger.Info().Msgf("Configured EVM forwarders: %+v", config)
		}
	}

	_, exist = chainsWithForwarders[blockchain.FamilyTron]
	if exist {
		for _, don := range consensusDons {
			tErr := configureTronForwarder(testLogger, creEnv, don)
			if tErr != nil {
				return errors.Wrap(tErr, "failed to configure Tron forwarders")
			}
		}
	}

	return nil
}

func deployTronForwarders(testLogger zerolog.Logger, cldfEnv *cldf.Environment, chainSelectors []uint64, contractVersions map[cre.ContractType]*semver.Version) error {
	memoryDatastore, mErr := contracts.NewDataStoreFromExisting(cldfEnv.DataStore)
	if mErr != nil {
		return fmt.Errorf("failed to create memory datastore: %w", mErr)
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

func configureTronForwarder(testLogger zerolog.Logger, creEnv *cre.Environment, don *cre.Don) error {
	triggerOptions := cldf_tron.DefaultTriggerOptions()
	triggerOptions.FeeLimit = 1_000_000_000

	wfNodeIDs := []string{}
	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return fmt.Errorf("failed to find worker nodes for Tron configuration: %w", wErr)
	}

	for _, node := range workerNodes {
		wfNodeIDs = append(wfNodeIDs, node.Keys.P2PKey.PeerID.String())
	}

	registryChain, rErr := creEnv.RegistryChain()
	if rErr != nil {
		return fmt.Errorf("failed to get registry chain: %w", rErr)
	}

	asEVM, ok := registryChain.(*evmblockchain.Blockchain)
	if !ok {
		return fmt.Errorf("registry chain is not *evmblockchain.Blockchain, but %T", registryChain)
	}

	capabilitiesRegistryAddress := contracts.MustGetAddressFromDataStore(creEnv.CldfEnvironment.DataStore, creEnv.RegistryChainSelector, keystone_changeset.CapabilitiesRegistry.String(), creEnv.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()], "")
	capReg, capErr := kcr.NewCapabilitiesRegistry(common.HexToAddress(capabilitiesRegistryAddress), asEVM.SethClient.Client)
	if capErr != nil {
		return fmt.Errorf("failed to create capabilities registry instance: %w", capErr)
	}

	configChangeset := commonchangeset.Configure(tronchangeset.ConfigureForwarder{}, &tronchangeset.ConfigureForwarderRequest{
		WFDonName:        don.Name,
		WFNodeIDs:        wfNodeIDs,
		RegistryChainSel: creEnv.RegistryChainSelector,
		Chains:           make(map[uint64]struct{}),
		TriggerOptions:   triggerOptions,
		Registry:         capReg,
	})

	_, err := commonchangeset.Apply(nil, *creEnv.CldfEnvironment, configChangeset)
	if err != nil {
		return fmt.Errorf("failed to configure Tron forwarders using changesets: %w", err)
	}

	testLogger.Info().Msgf("Configured TRON forwarder for v1 consensus on chain: %d", creEnv.RegistryChainSelector)

	return nil
}

func updateNodeConfig(workerNode *cre.NodeMetadata, currentConfig string, chainCapabilityConfigs map[string]*cre.ChainCapabilityConfig, creEnv *cre.Environment) (*string, error) {
	writeEvmConfigs := []writeEVMData{}

	// for each worker node find all supported chains and node's public address for each chain
	for _, chainID := range chainCapabilityConfigs[flag].EnabledChains {
		chain, exists := chain_selectors.ChainByEvmChainID(chainID)
		if !exists {
			return nil, errors.Errorf("failed to find selector for chain ID %d", chainID)
		}

		evmData := writeEVMData{
			ChainID:       chainID,
			ChainSelector: chain.Selector,
		}

		evmData.ForwarderAddress = contracts.MustGetAddressFromDataStore(creEnv.CldfEnvironment.DataStore, chain.Selector, keystone_changeset.KeystoneForwarder.String(), creEnv.ContractVersions[keystone_changeset.KeystoneForwarder.String()], "")
		evmKey, ok := workerNode.Keys.EVM[chainID]
		if !ok {
			return nil, fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
		}
		evmData.FromAddress = evmKey.PublicAddress

		var mergeErr error
		evmData, mergeErr = mergeDefaultAndRuntimeConfigValues(evmData, creEnv.CapabilityConfigs, chainCapabilityConfigs, chainID)
		if mergeErr != nil {
			return nil, errors.Wrap(mergeErr, "failed to merge default and runtime write-evm config values")
		}

		writeEvmConfigs = append(writeEvmConfigs, evmData)
	}

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

	return ptr.Ptr(string(stringifiedConfig)), nil
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
	GasLimitDefault = {{printf "%v" .GasLimitDefault}}
	TxAcceptanceState = {{printf "%v" .TxAcceptanceState}}
	PollPeriod = '{{.PollPeriod}}'
	AcceptanceTimeout = '{{.AcceptanceTimeout}}'
`
