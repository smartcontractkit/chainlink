package aptos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	router "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router"
	"github.com/smartcontractkit/chainlink-aptos/bindings/compile"
	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	module_mcms "github.com/smartcontractkit/chainlink-aptos/bindings/mcms/mcms"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/mcms"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	"github.com/smartcontractkit/mcms/types"
)

// CsDeployAptosChain deploys CCIP Package for Aptos chains
var CsDeployAptosChain deployment.ChangeSetV2[DeployAptosChainConfig] = CsDeployAptosChainImp{}

type CsDeployAptosChainImp struct {
	onChainState map[uint64]AptosCCIPChainState
	env          deployment.Environment
	ab           *deployment.AddressBookMap
	proposals    []mcms.Proposal
	opCount      uint64
}

func (cs CsDeployAptosChainImp) VerifyPreconditions(env deployment.Environment, config DeployAptosChainConfig) error {
	// Validate configs
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid DeployAptosChainConfig: %w", err)
	}

	// Validate env
	state, err := LoadOnchainStateAptos(env)
	if err != nil {
		return fmt.Errorf("failed to load existing onchain state: %w", err)
	}
	var errs []error
	for chainSel := range config.ContractParamsPerChain {
		if _, ok := env.AptosChains[chainSel]; !ok {
			errs = append(errs, fmt.Errorf("env not found for chain: %d", chainSel))
		}
		chainState := state[chainSel]
		if chainState.AptosMCMSObjAddr != (aptos.AccountAddress{}) {
			// If MCMS is already deployed for chain, skip this validation
			continue
		}
		mcmsDeployCfg, ok := config.MCMSConfigPerChain[chainSel]
		if !ok {
			errs = append(errs, fmt.Errorf("either MCMS deploy configs or MCMS address should be provided, chain: %d", chainSel))
		}
		err := mcmsDeployCfg.Validate()
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid MCMS config for chain %d: %w", chainSel, err))
		}
	}

	return errors.Join(errs...)
}

func (cs CsDeployAptosChainImp) Apply(env deployment.Environment, config DeployAptosChainConfig) (deployment.ChangesetOutput, error) {
	cs.ab = deployment.NewMemoryAddressBook()
	state, err := LoadOnchainStateAptos(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	cs.onChainState = state
	cs.env = env

	// For each aptos chain in the config generate proposals
	for chainSel := range config.ContractParamsPerChain {
		// Cleanup MCMS staging area
		err := cs.generateCleanupStagingProposal(chainSel)
		if err != nil {
			err := fmt.Errorf("failed to generate cleanup staging proposal for chain %d : %w", chainSel, err)
			env.Logger.Error(err)
			return deployment.ChangesetOutput{AddressBook: cs.ab}, err
		}
		// Generate proposals - Deploy CCIP package
		ccipObjectAddress, err := cs.generateDeployCCIPProposal(chainSel)
		if err != nil {
			err := fmt.Errorf("failed to generate CCIP deploy proposal for chain %d : %w", chainSel, err)
			env.Logger.Error(err)
			return deployment.ChangesetOutput{AddressBook: cs.ab}, err
		}
		// Generate proposals - Deploy Router package
		err = cs.generateDeployRouterProposal(chainSel, ccipObjectAddress)
		if err != nil {
			err := fmt.Errorf("failed to generate Router deploy proposal for chain %d : %w", chainSel, err)
			env.Logger.Error(err)
			return deployment.ChangesetOutput{AddressBook: cs.ab}, err
		}
		// TODO: Generate proposals - Initialize contracts
	}

	return deployment.ChangesetOutput{
		AddressBook:   cs.ab,
		MCMSProposals: cs.proposals,
	}, nil
}

// generateCleanupStagingProposal generates cleanup MCMS operations for the staging area if it's not clean
func (cs *CsDeployAptosChainImp) generateCleanupStagingProposal(chainSel uint64) error {
	chainState := cs.onChainState[chainSel]
	aptosChain := cs.env.AptosChains[chainSel]

	// Check resources first to see if staging is clean
	IsMCMSStagingAreaClean, err := IsMCMSStagingAreaClean(aptosChain.Client, chainState.AptosMCMSObjAddr)
	if err != nil {
		return fmt.Errorf("failed to check if MCMS staging area is clean: %w", err)
	}
	if IsMCMSStagingAreaClean {
		cs.env.Logger.Infow("MCMS Staging Area already clean", "addr", chainState.AptosMCMSObjAddr.String())
		return nil
	}

	// Bind MCMS contract
	mcmsContract := mcmsbind.Bind(chainState.AptosMCMSObjAddr, aptosChain.Client)

	// Get cleanup staging operations
	var operations []types.Operation
	module, function, _, args, err := mcmsContract.MCMSDeployer.EncodeCleanupStagingArea()
	if err != nil {
		return fmt.Errorf("failed to EncodeCleanupStagingArea: %w", err)
	}
	additionalFields := aptosmcms.AdditionalFields{
		ModuleName: module.Name,
		Function:   function,
	}
	afBytes, err := json.Marshal(additionalFields)
	if err != nil {
		return fmt.Errorf("failed to marshal additional fields: %w", err)
	}
	operations = append(operations, types.Operation{
		ChainSelector: types.ChainSelector(chainSel),
		Transaction: types.Transaction{
			To:               mcmsContract.Address.StringLong(),
			Data:             module_mcms.ArgsToData(args),
			AdditionalFields: afBytes,
		},
	})

	// Generate cleanup proposal
	proposal, err := cs.generateDeployProposal(mcmsContract, chainSel, operations, "Cleanup Staging Area")
	if err != nil {
		return fmt.Errorf("failed to create deploy proposal: %w", err)
	}
	cs.proposals = append(cs.proposals, *proposal)

	return nil
}

// generateDeployCCIPProposal generates deployment MCMS operations for the CCIP package
func (cs *CsDeployAptosChainImp) generateDeployCCIPProposal(chainSel uint64) (*aptos.AccountAddress, error) {
	chainState := cs.onChainState[chainSel]
	aptosChain := cs.env.AptosChains[chainSel]

	// Validate there's no package deployed
	if (chainState.AptosCCIPObjAddr != aptos.AccountAddress{}) {
		cs.env.Logger.Infow("CCIP Package already deployed", "addr", chainState.AptosCCIPObjAddr.String())
		return &chainState.AptosCCIPObjAddr, nil
	}

	// Compile, chunk and get CCIP deploy operations
	mcmsContract := mcmsbind.Bind(chainState.AptosMCMSObjAddr, aptosChain.Client)
	ccipObjectAddress, operations, err := cs.getCCIPDeployOperations(mcmsContract, chainSel)
	if err != nil {
		return nil, fmt.Errorf("failed to compile and create deploy operations: %w", err)
	}

	// Save the address of the CCIP object
	typeAndVersion := deployment.NewTypeAndVersion(AptosCCIPType, deployment.Version1_0_0)
	cs.ab.Save(chainSel, ccipObjectAddress.String(), typeAndVersion)

	// Generate deploy proposal
	proposal, err := cs.generateDeployProposal(mcmsContract, chainSel, operations, "Deploy CCIP Package")
	if err != nil {
		return nil, fmt.Errorf("failed to create deploy proposal: %w", err)
	}
	cs.proposals = append(cs.proposals, *proposal)

	return &ccipObjectAddress, nil
}

func (cs *CsDeployAptosChainImp) getCCIPDeployOperations(mcmsContract mcmsbind.MCMS, chainSel uint64) (aptos.AccountAddress, []types.Operation, error) {
	// Calculate addresses of the owner and the object
	ccipObjectAddress, err := mcmsContract.MCMSRegistry.GetNewCodeObjectAddress(nil, ccip.DefaultSeed)
	if err != nil {
		return ccipObjectAddress, []types.Operation{}, fmt.Errorf("failed to calculate object address: %w", err)
	}

	// Compile Package
	payload, err := ccip.Compile(ccipObjectAddress, mcmsContract.Address, true)
	if err != nil {
		return ccipObjectAddress, []types.Operation{}, fmt.Errorf("failed to compile: %w", err)
	}

	// Create chunks and stage operations
	operations, err := createChunksAndStage(payload, mcmsContract, chainSel, ccip.DefaultSeed, ccipObjectAddress)
	if err != nil {
		return ccipObjectAddress, operations, fmt.Errorf("failed to create chunks and stage for %d: %w", chainSel, err)
	}

	return ccipObjectAddress, operations, nil
}

// generateDeployRouterProposal generates deployment MCMS operations for the CCIP package
func (cs *CsDeployAptosChainImp) generateDeployRouterProposal(chainSel uint64, ccipObjectAddress *aptos.AccountAddress) error {
	chainState := cs.onChainState[chainSel]
	aptosChain := cs.env.AptosChains[chainSel]

	// TODO: is there a way to check if module exists?

	// Compile, chunk and get Router deploy operations
	mcmsContract := mcmsbind.Bind(chainState.AptosMCMSObjAddr, aptosChain.Client)
	operations, err := cs.getRouterDeployOperations(mcmsContract, chainSel, ccipObjectAddress)
	if err != nil {
		return fmt.Errorf("failed to compile and create deploy operations: %w", err)
	}

	// Generate deploy proposal
	proposal, err := cs.generateDeployProposal(mcmsContract, chainSel, operations, "Deploy Router Package")
	if err != nil {
		return fmt.Errorf("failed to create deploy proposal: %w", err)
	}
	cs.proposals = append(cs.proposals, *proposal)

	return nil
}

func (cs *CsDeployAptosChainImp) getRouterDeployOperations(
	mcmsContract mcmsbind.MCMS,
	chainSel uint64,
	ccipObjectAddress *aptos.AccountAddress,
) ([]types.Operation, error) {
	// Compile Package
	payload, err := router.Compile(*ccipObjectAddress, mcmsContract.Address, true)
	if err != nil {
		return []types.Operation{}, fmt.Errorf("failed to compile: %w", err)
	}

	// Create chunks and stage operations
	operations, err := createChunksAndStage(payload, mcmsContract, chainSel, "", *ccipObjectAddress)
	if err != nil {
		return operations, fmt.Errorf("failed to create chunks and stage for %d: %w", chainSel, err)
	}

	return operations, nil
}

func createChunksAndStage(
	payload compile.CompiledPackage,
	mcmsContract mcmsbind.MCMS,
	chainSel uint64,
	seed string,
	codeObjectAddress aptos.AccountAddress,
) ([]types.Operation, error) {
	// Validate seed XOR codeObjectAddress, one and only one must be provided
	if (seed != "") == (codeObjectAddress != aptos.AccountAddress{}) {
		return nil, fmt.Errorf("either provide seed to publishToObject or objectAddress to upgradeObjectCode")
	}

	var operations []types.Operation

	// Create chunks
	chunks, err := bind.CreateChunks(payload, bind.ChunkSizeInBytes)
	if err != nil {
		return operations, fmt.Errorf("failed to create chunks: %w", err)
	}

	// Stage chunks with mcms_deployer module and execute with the last one
	for i, chunk := range chunks {
		var (
			module   aptos.ModuleId
			function string
			args     [][]byte
			err      error
		)

		// First chunks get staged, the last one gets published or upgraded
		if i != len(chunks)-1 {
			module, function, _, args, err = mcmsContract.MCMSDeployer.EncodeStageCodeChunk(chunk.Metadata, chunk.CodeIndices, chunk.Chunks)
		} else if seed != "" {
			module, function, _, args, err = mcmsContract.MCMSDeployer.EncodeStageCodeChunkAndPublishToObject(chunk.Metadata, chunk.CodeIndices, chunk.Chunks, seed)
		} else {
			module, function, _, args, err = mcmsContract.MCMSDeployer.EncodeStageCodeChunkAndUpgradeObjectCode(chunk.Metadata, chunk.CodeIndices, chunk.Chunks, codeObjectAddress)
		}
		if err != nil {
			return operations, fmt.Errorf("failed to encode chunk %d: %w", i, err)
		}
		additionalFields := aptosmcms.AdditionalFields{
			ModuleName: module.Name,
			Function:   function,
		}
		afBytes, err := json.Marshal(additionalFields)
		if err != nil {
			return operations, fmt.Errorf("failed to marshal additional fields: %w", err)
		}
		operations = append(operations, types.Operation{
			ChainSelector: types.ChainSelector(chainSel),
			Transaction: types.Transaction{
				To:               mcmsContract.Address.StringLong(),
				Data:             module_mcms.ArgsToData(args),
				AdditionalFields: afBytes,
			},
		})
	}

	return operations, nil
}

func (cs *CsDeployAptosChainImp) generateDeployProposal(mcmsContract mcmsbind.MCMS, chainSel uint64, operations []types.Operation, description string) (*mcms.Proposal, error) {
	if cs.opCount == 0 {
		// Create MCMS inspector
		inspector := aptosmcms.NewInspector(cs.env.AptosChains[chainSel].Client)
		startingOpCount, err := inspector.GetOpCount(context.Background(), mcmsContract.Address.StringLong())
		if err != nil {
			return nil, fmt.Errorf("failed to get starting op count: %w", err)
		}
		cs.opCount = startingOpCount
	}

	// Create proposal builder
	validUntil := time.Now().Add(time.Hour * 24).Unix()
	proposalBuilder := mcms.NewProposalBuilder().
		SetVersion("v1").
		SetValidUntil(uint32(validUntil)).
		SetDescription(description).
		SetOverridePreviousRoot(true).
		AddChainMetadata(
			types.ChainSelector(chainSel),
			types.ChainMetadata{
				StartingOpCount: cs.opCount,
				MCMAddress:      mcmsContract.Address.StringLong(),
			},
		)

	// Add operations and build
	for _, op := range operations {
		proposalBuilder.AddOperation(op)
	}
	proposal, err := proposalBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build proposal: %w", err)
	}

	cs.opCount += uint64(len(operations))
	return proposal, nil
}
