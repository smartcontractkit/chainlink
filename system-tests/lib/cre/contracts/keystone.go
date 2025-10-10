package contracts

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

type DeployKeystoneContractsInput struct {
	CldfEnvironment           *cldf.Environment
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

	solForwardersSelectors := make([]uint64, 0)
	for _, bcOut := range input.CtfBlockchains {
		// consider we have just 1 solana chain
		if bcOut.SolChain != nil {
			solForwardersSelectors = append(solForwardersSelectors, bcOut.SolChain.ChainSelector)
			continue
		}
	}

	homeChainOutput := input.CtfBlockchains[0]

	// use CLD to deploy the registry contracts, which are required before constructing the node TOML configs
	homeChainSelector := homeChainOutput.ChainSelector
	deployRegistrySeq := ks_contracts_op.DeployRegistryContractsSequence
	if input.WithV2Registries {
		deployRegistrySeq = ks_contracts_op.DeployV2RegistryContractsSequence
	}

	registryContractsReport, seqErr := operations.ExecuteSequence(
		input.CldfEnvironment.OperationsBundle,
		deployRegistrySeq,
		ks_contracts_op.DeployContractsSequenceDeps{
			Env: input.CldfEnvironment,
		},
		ks_contracts_op.DeployRegistryContractsSequenceInput{
			RegistryChainSelector: homeChainSelector,
		},
	)
	if seqErr != nil {
		return nil, errors.Wrap(seqErr, "failed to deploy Keystone contracts")
	}

	if err := input.CldfEnvironment.ExistingAddresses.Merge(registryContractsReport.Output.AddressBook); err != nil { //nolint:staticcheck // won't migrate now
		return nil, errors.Wrap(err, "failed to merge address book with Keystone contracts addresses")
	}

	if err := memoryDatastore.Merge(registryContractsReport.Output.Datastore); err != nil {
		return nil, errors.Wrap(err, "failed to merge datastore with Keystone contracts addresses")
	}

	wfRegAddr := MustGetAddressFromMemoryDataStore(memoryDatastore, homeChainSelector, keystone_changeset.WorkflowRegistry.String(), input.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")
	testLogger.Info().Msgf("Deployed Workflow Registry %s contract on chain %d at %s", input.ContractVersions[keystone_changeset.WorkflowRegistry.String()], homeChainSelector, wfRegAddr)

	capRegAddr := MustGetAddressFromMemoryDataStore(memoryDatastore, homeChainSelector, keystone_changeset.CapabilitiesRegistry.String(), input.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()], "")
	testLogger.Info().Msgf("Deployed Capabilities Registry %s contract on chain %d at %s", input.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()], homeChainSelector, capRegAddr)

	input.CldfEnvironment.DataStore = memoryDatastore.Seal()

	return &DeployKeystoneContractsOutput{
		Env:             input.CldfEnvironment,
		MemoryDataStore: memoryDatastore,
	}, nil
}

func DeployOCR3Contract(logger zerolog.Logger, qualifier string, selector uint64, env *cldf.Environment, contractVersions map[string]string) (*ks_contracts_op.DeployOCR3ContractSequenceOutput, *common.Address, error) {
	memoryDatastore := datastore.NewMemoryDataStore()

	// load all existing addresses into memory datastore
	mergeErr := memoryDatastore.Merge(env.DataStore)
	if mergeErr != nil {
		return nil, nil, fmt.Errorf("failed to merge existing datastore into memory datastore: %w", mergeErr)
	}

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
		return nil, nil, fmt.Errorf("failed to deploy OCR3 contract '%s' on chain %d: %w", qualifier, selector, err)
	}
	// TODO: CRE-742 remove address book
	if err = env.ExistingAddresses.Merge(ocr3DeployReport.Output.AddressBook); err != nil { //nolint:staticcheck // won't migrate now
		return nil, nil, fmt.Errorf("failed to merge address book with OCR3 contract address for '%s' on chain %d: %w", qualifier, selector, err)
	}
	if err = memoryDatastore.Merge(ocr3DeployReport.Output.Datastore); err != nil {
		return nil, nil, fmt.Errorf("failed to merge datastore with OCR3 contract address for '%s' on chain %d: %w", qualifier, selector, err)
	}

	address := MustGetAddressFromMemoryDataStore(memoryDatastore, selector, keystone_changeset.OCR3Capability.String(), contractVersions[keystone_changeset.OCR3Capability.String()], qualifier)
	logger.Info().Msgf("Deployed OCR3 %s contract on chain %d at %s [qualifier: %s]", contractVersions[keystone_changeset.OCR3Capability.String()], selector, address, qualifier)

	env.DataStore = memoryDatastore.Seal()

	return &ocr3DeployReport.Output, &address, nil
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

func MightGetAddressFromDataStore(dataStore datastore.DataStore, chainSel uint64, contractType string, version string, qualifier string) *common.Address {
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
