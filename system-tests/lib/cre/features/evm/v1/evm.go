package features

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"slices"
	"sort"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	corevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	evmworkflow "github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	chainlinkbig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	tronchangeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/tron"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const flag = cre.WriteEVMCapability

type EVM struct{}

func (o *EVM) Flag() cre.CapabilityFlag {
	return flag
}

func (o *EVM) PreDONStartup(
	testLogger zerolog.Logger,
	registryChainSelector uint64,
	cldfEnv *cldf.Environment,
	provider infra.Provider,
	topology *cre.Topology,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	capabilityConfigs cre.CapabilityConfigs,
	contractVersions map[string]string,
) error {
	// deploy forwarder contracts, if needed
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
					evmForwardersSelectors = append(evmForwardersSelectors, bcOut.ChainSelector)
				}
			}
		}
	}

	if len(evmForwardersSelectors) == 0 && len(tronForwardersSelectors) == 0 {
		return nil
	}

	memoryDatastore := datastore.NewMemoryDataStore()

	// load all existing addresses into memory datastore
	mergeErr := memoryDatastore.Merge(cldfEnv.DataStore)
	if mergeErr != nil {
		return fmt.Errorf("failed to merge existing datastore into memory datastore: %w", mergeErr)
	}

	// deploy evm forwarders
	if len(evmForwardersSelectors) > 0 {
		evmForwardersReport, seqErr2 := operations.ExecuteSequence(
			cldfEnv.OperationsBundle,
			forwarder.DeploySequence,
			forwarder.DeploySequenceDeps{
				Env: cldfEnv,
			},
			forwarder.DeploySequenceInput{
				Targets: evmForwardersSelectors,
			},
		)
		if seqErr2 != nil {
			return errors.Wrap(seqErr2, "failed to deploy evm forwarder")
		}

		if seqErr2 = cldfEnv.ExistingAddresses.Merge(evmForwardersReport.Output.AddressBook); seqErr2 != nil { //nolint:staticcheck // won't migrate now
			return errors.Wrap(seqErr2, "failed to merge address book with Keystone contracts addresses")
		}

		if seqErr2 = memoryDatastore.Merge(evmForwardersReport.Output.Datastore); seqErr2 != nil {
			return errors.Wrap(seqErr2, "failed to merge datastore with Keystone contracts addresses")
		}

		for _, forwarderSelector := range evmForwardersSelectors {
			forwarderAddr := MustGetAddressFromMemoryDataStore(memoryDatastore, forwarderSelector, keystone_changeset.KeystoneForwarder.String(), contractVersions[keystone_changeset.KeystoneForwarder.String()], "")
			testLogger.Info().Msgf("Deployed Forwarder %s contract on chain %d at %s", contractVersions[keystone_changeset.KeystoneForwarder.String()], forwarderSelector, forwarderAddr)
		}
	}

	// deploy tron forwarders
	if len(tronForwardersSelectors) > 0 {
		tronErr := deployTronForwarders(cldfEnv, tronForwardersSelectors)
		if tronErr != nil {
			return errors.Wrap(tronErr, "failed to deploy Tron Keystone forwarder contracts using changesets")
		}

		err := memoryDatastore.Merge(cldfEnv.DataStore)
		if err != nil {
			return errors.Wrap(err, "failed to merge Tron deployment results into main datastore")
		}
	}

	// update the CRE environment datastore to include the newly deployed contracts
	cldfEnv.DataStore = memoryDatastore.Seal()

	// update node configs to include write-evm config
	for _, donMetadata := range topology.DonsMetadata.List() {
		workerNodes, wErr := donMetadata.WorkerNodes()
		if wErr != nil {
			return errors.Wrap(wErr, "failed to find worker nodes")
		}

		for _, workerNode := range workerNodes {
			// // get all the forwarders and add workflow config (FromAddress + Forwarder) for chains that have write-evm enabled
			data := []writeEVMData{}
			for _, chainID := range donMetadata.CapabilitiesAwareNodeSet().ChainCapabilities[flag].EnabledChains {
				chain, exists := chain_selectors.ChainByEvmChainID(chainID)
				if !exists {
					return errors.Errorf("failed to find selector for chain ID %d", chainID)
				}

				evmData := writeEVMData{
					ChainID:       chainID,
					ChainSelector: chain.Selector,
				}

				forwarderAddress, fErr := findForwarderAddress(chain, cldfEnv.ExistingAddresses)
				if fErr != nil {
					return errors.Errorf("failed to find forwarder address for chain %d", chain.Selector)
				}
				evmData.ForwarderAddress = forwarderAddress.Hex()

				evmKey, ok := workerNode.Keys.EVM[chainID]
				if !ok {
					return fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
				}
				evmData.FromAddress = evmKey.PublicAddress

				var mergeErr error
				evmData, mergeErr = mergeDefaultAndRuntimeConfigValues(evmData, capabilityConfigs, donMetadata.CapabilitiesAwareNodeSet().ChainCapabilities, chainID)
				if mergeErr != nil {
					return errors.Wrap(mergeErr, "failed to merge default and runtime write-evm config values")
				}

				data = append(data, evmData)
			}

			currentConfig := donMetadata.CapabilitiesAwareNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides

			var typedConfig corechainlink.Config
			unmarshallErr := toml.Unmarshal([]byte(currentConfig), &typedConfig)
			if unmarshallErr != nil {
				return errors.Wrapf(unmarshallErr, "failed to unmarshal config for node index %d", workerNode.Index)
			}

			if len(typedConfig.EVM) < len(data) {
				return fmt.Errorf("not enough EVM chains configured in node index %d to add write-evm (evm v1) config. Expected at least %d chains, but found %d", workerNode.Index, len(data), len(typedConfig.EVM))
			}

			for _, writeEVMInput := range data {
				chainFound := false
				for idx, evmChain := range typedConfig.EVM {
					chainIDIsEqual := evmChain.ChainID.Cmp(chainlinkbig.New(big.NewInt(libc.MustSafeInt64(writeEVMInput.ChainID)))) == 0
					if chainIDIsEqual {
						evmWorkflow, evmErr := buildEVMWorkflowConfig(writeEVMInput)
						if evmErr != nil {
							return errors.Wrap(evmErr, "failed to build EVM workflow config")
						}

						typedConfig.EVM[idx].Workflow = *evmWorkflow
						typedConfig.EVM[idx].Transactions.ForwardersEnabled = ptr.Ptr(true)

						chainFound = true
						break
					}
				}

				if !chainFound {
					return fmt.Errorf("failed to find EVM chain with ID %d in the config of node index %d to add write-evm config", writeEVMInput.ChainID, workerNode.Index)
				}
			}

			stringifiedConfig, mErr := toml.Marshal(typedConfig)
			if mErr != nil {
				return errors.Wrapf(mErr, "failed to marshal config for node index %d", workerNode.Index)
			}

			donMetadata.CapabilitiesAwareNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides = string(stringifiedConfig)
		}
	}

	return nil
}

func (o *EVM) PostDONStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	creEnv *cre.Environment,
	nodeSetOutput []*cre.WrappedNodeOutput,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	contractVersions map[string]string,
) (*cre.PostDONStartupOutput, error) {
	// configure forwarders contracts
	allAddresses, addrErr := creEnv.CldfEnvironment.ExistingAddresses.Addresses() //nolint:staticcheck // ignore SA1019 as ExistingAddresses is deprecated but still used
	if addrErr != nil {
		return nil, errors.Wrap(addrErr, "failed to get addresses from address book")
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

	dons, donsErr := toDons(creEnv)
	if donsErr != nil {
		return nil, fmt.Errorf("failed to convert to dons: %w", donsErr)
	}

	consensusV1DON, err := dons.shouldBeOneDon(cre.ConsensusCapability)
	if err != nil {
		return nil, fmt.Errorf("failed to get consensus v1 DON: %w", err)
	}

	// configure EVM forwarders
	if len(evmChainsWithForwarders) > 0 {
		forwarderCfg := forwarder.DonConfiguration{
			Name:    consensusV1DON.Name,
			ID:      consensusV1DON.id,
			F:       consensusV1DON.F,
			Version: 1, // TODO this should be dynamic, but we don't have cap reg configured at this point
			NodeIDs: consensusV1DON.Nops[0].Nodes,
		}
		fout, err3 := operations.ExecuteSequence(
			creEnv.CldfEnvironment.OperationsBundle,
			forwarder.ConfigureSeq,
			forwarder.ConfigureSeqDeps{
				Env: creEnv.CldfEnvironment,
			},
			forwarder.ConfigureSeqInput{
				DON:    forwarderCfg,
				Chains: evmChainsWithForwarders,
			},
		)
		if err3 != nil {
			return nil, errors.Wrap(err3, "failed to configure forwarders")
		}

		testLogger.Info().Msgf("Configured forwarders for v1 consensus: %+v", fout.Output.Config)
	}

	// configure TRON forwarders
	if len(tronChainsWithForwarders) > 0 {
		tErr := configureTronForwarders(creEnv.CldfEnvironment, creEnv.DonTopology.HomeChainSelector, creEnv.DonTopology)
		if tErr != nil {
			return nil, errors.Wrap(tErr, "failed to configure TRON forwarders")
		}
	}

	// return capabilities registry configuration data
	capabilities := make(map[int][]keystone_changeset.DONCapabilityWithConfig)

	for donIdx, don := range creEnv.DonTopology.Dons.List() {
		if capabilities[donIdx] == nil {
			capabilities[donIdx] = []keystone_changeset.DONCapabilityWithConfig{}
		}
		if don.HasFlag(flag) {
			for _, chainID := range don.ChainCapabilities()[flag].EnabledChains {
				fullName := corevm.GenerateWriteTargetName(chainID)
				splitName := strings.Split(fullName, "@")

				capabilities[donIdx] = append(capabilities[donIdx], keystone_changeset.DONCapabilityWithConfig{
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
	}

	return &cre.PostDONStartupOutput{
		DONCapabilityWithConfigs: capabilities,
	}, nil
}

/*
	CODE BELOW WAS COPIED FROM VAROIUS PLACES IN THE SYSTEM TESTS AND (SOMETIMES) MODIFIED
	TO SHOWCASE THE CONCEPT OF FEATURES.

	IN THE FUTURE, WE SHOULD REFACTOR THE CODE TO AVOID DUPLICATIONS.
*/

// for now copy from system-tests/lib/cre/contracts/keystone.go
func MustGetAddressFromMemoryDataStore(dataStore *datastore.MemoryDataStore, chainSel uint64, contractType string, version string, qualifier string) common.Address {
	key := datastore.NewAddressRefKey(
		chainSel,
		datastore.ContractType(contractType),
		semver.MustParse(version),
		qualifier,
	)
	addrRef, err := dataStore.Addresses().Get(key)
	if err != nil {
		panic(fmt.Sprintf("Failed to get %s %s (qualifier=%s) address for chain %d: %s", contractType, version, qualifier, chainSel, err.Error()))
	}
	return common.HexToAddress(addrRef.Address)
}

func deployTronForwarders(env *cldf.Environment, chainSelectors []uint64) error {
	deployOptions := cldf_tron.DefaultDeployOptions()
	deployOptions.FeeLimit = 1_000_000_000

	deployChangeset := commonchangeset.Configure(tronchangeset.DeployForwarder{}, &tronchangeset.DeployForwarderRequest{
		ChainSelectors: chainSelectors,
		Qualifier:      "",
		DeployOptions:  deployOptions,
	})

	updatedEnv, err := commonchangeset.Apply(nil, *env, deployChangeset)
	if err != nil {
		return fmt.Errorf("failed to deploy Tron forwarders using changesets: %w", err)
	}

	env.ExistingAddresses = updatedEnv.ExistingAddresses //nolint:staticcheck // won't migrate now

	if updatedEnv.DataStore != nil {
		memoryDS := datastore.NewMemoryDataStore()
		err = memoryDS.Merge(env.DataStore)
		if err != nil {
			return fmt.Errorf("failed to merge existing datastore: %w", err)
		}
		err = memoryDS.Merge(updatedEnv.DataStore)
		if err != nil {
			return fmt.Errorf("failed to merge updated datastore: %w", err)
		}
		env.DataStore = memoryDS.Seal()
	}

	return nil
}

// copied from system-tests/lib/cre/contracts/contracts.go
func configureTronForwarders(env *cldf.Environment, registryChainSelector uint64, donTopology *cre.DonTopology) error {
	triggerOptions := cldf_tron.DefaultTriggerOptions()
	triggerOptions.FeeLimit = 1_000_000_000

	var wfNodeIDs []string
	for _, don := range donTopology.Dons.List() {
		if flags.HasOnlyOneFlag(don.Flags, cre.GatewayDON) {
			continue
		}

		workerNodes, wErr := don.WorkerNodes()
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

	return nil
}

// copied from system-tests/lib/cre/capabilities/writeevm/write_evm.go
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
			cre.WriteEVMCapability,
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

// copied from system-tests/lib/cre/contracts/contracts.go and modified
func toDons(
	creEnv *cre.Environment,
) (*dons, error) {
	dons := &dons{
		c:        make(map[string]donConfig),
		offChain: creEnv.CldfEnvironment.Offchain,
	}

	for _, don := range creEnv.DonTopology.Dons.List() {
		// if it's only a gateway DON, we don't want to register it with the Capabilities Registry
		// since it doesn't have any capabilities
		if flags.HasOnlyOneFlag(don.Flags, cre.GatewayDON) {
			continue
		}

		var capabilities []keystone_changeset.DONCapabilityWithConfig

		// // check what capabilities each DON has and register them with Capabilities Registry contract
		// for _, configFn := range input.CapabilityRegistryConfigFns {
		// 	if configFn == nil {
		// 		continue
		// 	}

		// 	enabledCapabilities, err2 := configFn(don.Flags, input.NodeSets[donIdx])
		// 	if err2 != nil {
		// 		return nil, errors.Wrap(err2, "failed to get capabilities from config function")
		// 	}

		// 	capabilities = append(capabilities, enabledCapabilities...)
		// }

		workerNodes, wErr := don.WorkerNodes()
		if wErr != nil {
			return nil, errors.Wrap(wErr, "failed to find worker nodes")
		}

		donPeerIDs := make([]string, len(workerNodes))
		for i, node := range workerNodes {
			// we need to use p2pID here with the "p2p_" prefix
			donPeerIDs[i] = node.Keys.P2PKey.PeerID.String()
		}

		forwarderF := (len(workerNodes) - 1) / 3
		if forwarderF == 0 {
			if flags.HasFlag(don.Flags, cre.ConsensusCapability) || flags.HasFlag(don.Flags, cre.ConsensusCapabilityV2) {
				return nil, fmt.Errorf("incorrect number of worker nodes: %d. Resulting F must conform to formula: mod((N-1)/3) > 0", len(workerNodes))
			}
			// for other capabilities, we can use 1 as F
			forwarderF = 1
		}

		// we only need to assign P2P IDs to NOPs, since `ConfigureInitialContractsChangeset` method
		// will take care of creating DON to Nodes mapping
		nop := keystone_changeset.NOP{
			Name:  fmt.Sprintf("NOP for %s DON", don.Name),
			Nodes: donPeerIDs,
		}
		donName := don.Name + "-don"
		c := keystone_changeset.DonCapabilities{
			Name:         donName,
			F:            libc.MustSafeUint8(forwarderF),
			Nops:         []keystone_changeset.NOP{nop},
			Capabilities: capabilities,
		}

		dons.c[donName] = donConfig{
			id:              uint32(don.ID), //nolint:gosec // G115
			DonCapabilities: c,
			flags:           don.Flags,
		}
	}

	return dons, nil
}

type dons struct {
	c        map[string]donConfig
	offChain offchain.Client
}

func (d *dons) donsOrderedByID() []donConfig {
	out := make([]donConfig, 0, len(d.c))
	for _, don := range d.c {
		out = append(out, don)
	}

	// Use sort library to sort by ID
	sort.Slice(out, func(i, j int) bool {
		return out[i].id < out[j].id
	})

	return out
}

func (d *dons) ListByFlag(flag cre.CapabilityFlag) ([]donConfig, error) {
	out := make([]donConfig, 0)
	for _, don := range d.donsOrderedByID() {
		if flags.HasFlag(don.flags, flag) {
			out = append(out, don)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("don with flag %s not found", flag)
	}
	return out, nil
}

func (d *dons) shouldBeOneDon(flag cre.CapabilityFlag) (donConfig, error) {
	dons, err := d.ListByFlag(flag)
	if err != nil {
		return donConfig{}, err
	}
	if len(dons) != 1 {
		return donConfig{}, fmt.Errorf("expected exactly one DON with flag %s, found %d", flag, len(dons))
	}
	return dons[0], nil
}

type donConfig struct {
	id uint32 // the DON id as registered in the capabilities registry
	keystone_changeset.DonCapabilities
	flags []cre.CapabilityFlag
}
