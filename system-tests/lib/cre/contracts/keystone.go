package contracts

import (
	"context"
	"fmt"
	"slices"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	"github.com/smartcontractkit/chainlink/deployment"
	creforwarder "github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	ks_sol "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"
	ks_sol_seq "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana/sequence"
	ks_sol_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana/sequence/operation"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

const (
	OCR3ContractQualifier         = "capability_ocr3"
	DONTimeContractQualifier      = "capability_dontime"
	VaultOCR3ContractQualifier    = "capability_vault"
	VaultDKGOCR3ContractQualifier = "capability_vault_dkg"
	ConsensusV2ContractQualifier  = "capability_consensus"
)

type DeployKeystoneContractsInput struct {
	CldfBlockchains           map[uint64]cldf_chain.BlockChain
	CtfBlockchains            []*cre.WrappedBlockchainOutput
	ContractVersions          map[string]string
	WithV2Registries          bool
	CapabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet
}

type DeployKeystoneContractsOutput struct {
	Env             *cldf.Environment
	MemoryDataStore *datastore.MemoryDataStore
}

func DeployKeystoneContracts(
	ctx context.Context,
	testLogger zerolog.Logger,
	singleFileLogger logger.Logger,
	input DeployKeystoneContractsInput,
) (*DeployKeystoneContractsOutput, error) {
	memoryDatastore := datastore.NewMemoryDataStore()
	allChainsCLDEnvironment := &cldf.Environment{
		Logger:            singleFileLogger,
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		DataStore:         memoryDatastore.Seal(),
		GetContext: func() context.Context {
			return ctx
		},
		BlockChains: cldf_chain.NewBlockChains(input.CldfBlockchains),
	}
	allChainsCLDEnvironment.OperationsBundle = operations.NewBundle(allChainsCLDEnvironment.GetContext, singleFileLogger, operations.NewMemoryReporter())

	evmForwardersSelectors := make([]uint64, 0)
	solForwardersSelectors := make([]uint64, 0)
	for _, bcOut := range input.CtfBlockchains {
		for _, donMetadata := range input.CapabilitiesAwareNodeSets {
			if slices.Contains(evmForwardersSelectors, bcOut.ChainSelector) {
				continue
			}
			// consider we have just 1 solana chain
			if bcOut.SolChain != nil {
				solForwardersSelectors = append(solForwardersSelectors, bcOut.SolChain.ChainSelector)
				continue
			}
			if flags.RequiresForwarderContract(donMetadata.ComputedCapabilities, bcOut.ChainID) {
				evmForwardersSelectors = append(evmForwardersSelectors, bcOut.ChainSelector)
			}
		}
	}

	var allNodeFlags []string
	for i := range input.CapabilitiesAwareNodeSets {
		nodeFlags, err := flags.NodeSetFlags(input.CapabilitiesAwareNodeSets[i])
		if err != nil {
			continue
		}
		allNodeFlags = append(allNodeFlags, nodeFlags...)
	}
	vaultOCR3AddrFlag := flags.HasFlag(allNodeFlags, cre.VaultCapability)
	evmOCR3AddrFlag := flags.HasFlagForAnyChain(allNodeFlags, cre.EVMCapability)
	consensusV2AddrFlag := flags.HasFlag(allNodeFlags, cre.ConsensusCapabilityV2)

	chainsWithEVMCapability := ChainsWithEVMCapability(input.CtfBlockchains, input.CapabilitiesAwareNodeSets)
	homeChainOutput := input.CtfBlockchains[0]

	// use CLD to deploy the registry contracts, which are required before constructing the node TOML configs
	homeChainSelector := homeChainOutput.ChainSelector
	deployRegistrySeq := ks_contracts_op.DeployRegistryContractsSequence
	if input.WithV2Registries {
		deployRegistrySeq = ks_contracts_op.DeployV2RegistryContractsSequence
	}

	registryContractsReport, seqErr := operations.ExecuteSequence(
		allChainsCLDEnvironment.OperationsBundle,
		deployRegistrySeq,
		ks_contracts_op.DeployContractsSequenceDeps{
			Env: allChainsCLDEnvironment,
		},
		ks_contracts_op.DeployRegistryContractsSequenceInput{
			RegistryChainSelector: homeChainSelector,
		},
	)
	if seqErr != nil {
		return nil, errors.Wrap(seqErr, "failed to deploy Keystone contracts")
	}

	if err := allChainsCLDEnvironment.ExistingAddresses.Merge(registryContractsReport.Output.AddressBook); err != nil { //nolint:staticcheck // won't migrate now
		return nil, errors.Wrap(err, "failed to merge address book with Keystone contracts addresses")
	}

	if err := memoryDatastore.Merge(registryContractsReport.Output.Datastore); err != nil {
		return nil, errors.Wrap(err, "failed to merge datastore with Keystone contracts addresses")
	}
	if len(evmForwardersSelectors) > 0 {
		// deploy evm forwarders
		evmForwardersReport, seqErr2 := operations.ExecuteSequence(
			allChainsCLDEnvironment.OperationsBundle,
			creforwarder.DeploySequence,
			creforwarder.DeploySequenceDeps{
				Env: allChainsCLDEnvironment,
			},
			creforwarder.DeploySequenceInput{
				Targets: evmForwardersSelectors,
			},
		)
		if seqErr2 != nil {
			return nil, errors.Wrap(seqErr2, "failed to deploy evm forwarder")
		}

		if seqErr2 = allChainsCLDEnvironment.ExistingAddresses.Merge(evmForwardersReport.Output.AddressBook); seqErr2 != nil { //nolint:staticcheck // won't migrate now
			return nil, errors.Wrap(seqErr2, "failed to merge address book with Keystone contracts addresses")
		}

		if seqErr2 = memoryDatastore.Merge(evmForwardersReport.Output.Datastore); seqErr2 != nil {
			return nil, errors.Wrap(seqErr2, "failed to merge datastore with Keystone contracts addresses")
		}

		for _, forwarderSelector := range evmForwardersSelectors {
			forwarderAddr := MustGetAddressFromMemoryDataStore(memoryDatastore, forwarderSelector, keystone_changeset.KeystoneForwarder.String(), input.ContractVersions[keystone_changeset.KeystoneForwarder.String()], "")
			testLogger.Info().Msgf("Deployed Forwarder %s contract on chain %d at %s", input.ContractVersions[keystone_changeset.KeystoneForwarder.String()], forwarderSelector, forwarderAddr)
		}
	}

	// deploy solana forwarders
	for _, sel := range solForwardersSelectors {
		populateContracts := map[string]datastore.ContractType{
			deployment.KeystoneForwarderProgramName: ks_sol.ForwarderContract,
		}
		version := semver.MustParse(input.ContractVersions[ks_sol.ForwarderContract.String()])

		// Forwarder for solana is predeployed on chain spin-up. We jus need to add it to memory datastore here
		errp := memory.PopulateDatastore(memoryDatastore.AddressRefStore, populateContracts,
			version, ks_sol.DefaultForwarderQualifier, sel)
		if errp != nil {
			return nil, errors.Wrap(errp, "failed to populate datastore with predeployed contracts")
		}
		out, err := operations.ExecuteSequence(
			allChainsCLDEnvironment.OperationsBundle,
			ks_sol_seq.DeployForwarderSeq,
			ks_sol_op.Deps{
				Env:       *allChainsCLDEnvironment,
				Chain:     allChainsCLDEnvironment.BlockChains.SolanaChains()[sel],
				Datastore: memoryDatastore.Seal(),
			},
			ks_sol_seq.DeployForwarderSeqInput{
				ChainSel:     sel,
				ProgramName:  deployment.KeystoneForwarderProgramName,
				Qualifier:    ks_sol.DefaultForwarderQualifier,
				ContractType: ks_sol.ForwarderContract,
				Version:      version,
			},
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to deploy sol forwarder")
		}

		err = memoryDatastore.AddressRefStore.Add(datastore.AddressRef{
			Address:       out.Output.State.String(),
			ChainSelector: sel,
			Version:       semver.MustParse(input.ContractVersions[ks_sol.ForwarderState.String()]),
			Qualifier:     ks_sol.DefaultForwarderQualifier,
			Type:          ks_sol.ForwarderState,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to add address to the datastore for Solana Forwarder state")
		}

		testLogger.Info().Msgf("Deployed Forwarder %s contract on Solana chain chain %d programID: %s state: %s", input.ContractVersions[ks_sol.ForwarderContract.String()], sel, out.Output.ProgramID.String(), out.Output.State.String())
	}

	wfRegAddr := MustGetAddressFromMemoryDataStore(memoryDatastore, homeChainSelector, keystone_changeset.WorkflowRegistry.String(), input.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")
	testLogger.Info().Msgf("Deployed Workflow Registry %s contract on chain %d at %s", input.ContractVersions[keystone_changeset.WorkflowRegistry.String()], homeChainSelector, wfRegAddr)

	capRegAddr := MustGetAddressFromMemoryDataStore(memoryDatastore, homeChainSelector, keystone_changeset.CapabilitiesRegistry.String(), input.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()], "")
	testLogger.Info().Msgf("Deployed Capabilities Registry %s contract on chain %d at %s", input.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()], homeChainSelector, capRegAddr)

	// deploy the various ocr contracts
	// TODO move this deeper into the stack when we have all the p2p ids and can deploy and configure in one sequence
	// deploy OCR3 contract
	// we deploy OCR3 contract with a qualifier, so that we can distinguish it from other OCR3 contracts (Vault, EVM, ConsensusV2)
	_, seqErr = deployOCR3Contract(OCR3ContractQualifier, homeChainSelector, allChainsCLDEnvironment, memoryDatastore)
	if seqErr != nil {
		return nil, fmt.Errorf("failed to deploy OCR3 contract %w", seqErr)
	}
	ocr3Addr := MustGetAddressFromMemoryDataStore(memoryDatastore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], OCR3ContractQualifier)
	testLogger.Info().Msgf("Deployed OCR3 %s contract on chain %d at %s", input.ContractVersions[keystone_changeset.OCR3Capability.String()], homeChainSelector, ocr3Addr)

	// deploy DONTime contract
	_, seqErr = deployOCR3Contract(DONTimeContractQualifier, homeChainSelector, allChainsCLDEnvironment, memoryDatastore)
	if seqErr != nil {
		return nil, fmt.Errorf("failed to deploy DONTime contract %w", seqErr)
	}
	donTimeAddr := MustGetAddressFromMemoryDataStore(memoryDatastore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], DONTimeContractQualifier)
	testLogger.Info().Msgf("Deployed OCR3 %s (DON Time) contract on chain %d at %s", input.ContractVersions[keystone_changeset.OCR3Capability.String()], homeChainSelector, donTimeAddr)

	// deploy Vault OCR3 contract
	if vaultOCR3AddrFlag {
		_, seqErr = deployOCR3Contract(VaultOCR3ContractQualifier, homeChainSelector, allChainsCLDEnvironment, memoryDatastore)
		if seqErr != nil {
			return nil, fmt.Errorf("failed to deploy Vault OCR3 contract %w", seqErr)
		}
		_, seqErr = deployOCR3Contract(VaultDKGOCR3ContractQualifier, homeChainSelector, allChainsCLDEnvironment, memoryDatastore)
		if seqErr != nil {
			return nil, fmt.Errorf("failed to deploy DKG contract %w", seqErr)
		}

		vaultOCR3Addr := MustGetAddressFromMemoryDataStore(memoryDatastore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], VaultOCR3ContractQualifier)
		testLogger.Info().Msgf("Deployed OCR3 %s (Vault) contract on chain %d at %s", input.ContractVersions[keystone_changeset.OCR3Capability.String()], homeChainSelector, vaultOCR3Addr)

		vaultDKGOCR3Addr := MustGetAddressFromMemoryDataStore(memoryDatastore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], VaultDKGOCR3ContractQualifier)
		testLogger.Info().Msgf("Deployed OCR3 %s (DKG) contract on chain %d at %s", input.ContractVersions[keystone_changeset.OCR3Capability.String()], homeChainSelector, vaultDKGOCR3Addr)
	}

	// deploy EVM OCR3 contracts
	if evmOCR3AddrFlag {
		for chainID, selector := range chainsWithEVMCapability {
			qualifier := ks_contracts_op.CapabilityContractIdentifier(uint64(chainID))
			_, seqErr = deployOCR3Contract(qualifier, homeChainSelector, allChainsCLDEnvironment, memoryDatastore)
			if seqErr != nil {
				return nil, fmt.Errorf("failed to deploy EVM OCR3 contract for chainID %d, selector %d: %w", chainID, selector, seqErr)
			}

			evmOCR3Addr := MustGetAddressFromMemoryDataStore(memoryDatastore, homeChainSelector, keystone_changeset.OCR3Capability.String(), "1.0.0", qualifier)
			testLogger.Info().Msgf("Deployed EVM OCR3 contract (chainID %d) on chainID: %d, selector: %d, at: %s", chainID, homeChainOutput.ChainID, homeChainSelector, evmOCR3Addr)
		}
	}

	// deploy Consensus V2 OCR3 contract
	if consensusV2AddrFlag {
		_, seqErr = deployOCR3Contract(ConsensusV2ContractQualifier, homeChainSelector, allChainsCLDEnvironment, memoryDatastore)
		if seqErr != nil {
			return nil, fmt.Errorf("failed to deploy Consensus V2 OCR3 contract %w", seqErr)
		}
		consensusV2OCR3Addr := MustGetAddressFromMemoryDataStore(memoryDatastore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], ConsensusV2ContractQualifier)
		testLogger.Info().Msgf("Deployed Consensus V2 OCR3 %s contract on chain %d at %s", input.ContractVersions[keystone_changeset.OCR3Capability.String()], homeChainSelector, consensusV2OCR3Addr)
	}
	allChainsCLDEnvironment.DataStore = memoryDatastore.Seal()

	return &DeployKeystoneContractsOutput{
		Env:             allChainsCLDEnvironment,
		MemoryDataStore: memoryDatastore,
	}, nil
}

func deployOCR3Contract(qualifier string, selector uint64, env *cldf.Environment, ds datastore.MutableDataStore) (*ks_contracts_op.DeployOCR3ContractSequenceOutput, error) {
	ocr3DeployReport, err := operations.ExecuteSequence(
		env.OperationsBundle,
		ks_contracts_op.DeployOCR3ContractsSequence,
		ks_contracts_op.DeployOCR3ContractSequenceDeps{
			Env: env,
		},
		ks_contracts_op.DeployOCR3ContractSequenceInput{
			ChainSelector: selector,
			Qualifier:     qualifier,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy OCR3 contract '%s' on chain %d: %w", qualifier, selector, err)
	}
	// TODO: CRE-742 remove address book
	if err = env.ExistingAddresses.Merge(ocr3DeployReport.Output.AddressBook); err != nil { //nolint:staticcheck // won't migrate now
		return nil, fmt.Errorf("failed to merge address book with OCR3 contract address for '%s' on chain %d: %w", qualifier, selector, err)
	}
	if err = ds.Merge(ocr3DeployReport.Output.Datastore); err != nil {
		return nil, fmt.Errorf("failed to merge datastore with OCR3 contract address for '%s' on chain %d: %w", qualifier, selector, err)
	}
	return &ocr3DeployReport.Output, nil
}

func ChainsWithEVMCapability(chains []*cre.WrappedBlockchainOutput, nodeSets []*cre.CapabilitiesAwareNodeSet) map[ks_contracts_op.EVMChainID]ks_contracts_op.Selector {
	chainsWithEVMCapability := make(map[ks_contracts_op.EVMChainID]ks_contracts_op.Selector)
	for _, chain := range chains {
		for _, donMetadata := range nodeSets {
			if flags.HasFlagForChain(donMetadata.ComputedCapabilities, cre.EVMCapability, chain.ChainID) {
				if chainsWithEVMCapability[ks_contracts_op.EVMChainID(chain.ChainID)] != 0 {
					continue
				}
				chainsWithEVMCapability[ks_contracts_op.EVMChainID(chain.ChainID)] = ks_contracts_op.Selector(chain.ChainSelector)
			}
		}
	}

	return chainsWithEVMCapability
}

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

func MightGetAddressFromMemoryDataStore(dataStore *datastore.MemoryDataStore, chainSel uint64, contractType string, version string, qualifier string) *common.Address {
	key := datastore.NewAddressRefKey(
		chainSel,
		datastore.ContractType(contractType),
		semver.MustParse(version),
		qualifier,
	)

	addrRef, err := dataStore.Addresses().Get(key)
	if err != nil {
		return nil
	}

	return ptr.Ptr(common.HexToAddress(addrRef.Address))
}

func MustGetAddressFromDataStore(dataStore datastore.DataStore, chainSel uint64, contractType string, version string, qualifier string) string {
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
	return addrRef.Address
}
