package operation

import (
	"encoding/json"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	router "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router"
	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	aptoscfg "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/mcms"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	"github.com/smartcontractkit/mcms/types"
)

type CCIPDeploymentOperations struct {
	Env          deployment.Environment
	Ab           *deployment.AddressBookMap
	AptosChain   deployment.AptosChain
	OnChainState changeset.AptosCCIPChainState
	CCIPConfig   aptoscfg.ChainContractParams
	Proposals    *[]mcms.Proposal
	MCMSOpCount  uint64
}

// GenerateCleanupStagingProposal generates cleanup MCMS operations for the staging area if it's not clean
func (op *CCIPDeploymentOperations) GenerateCleanupStagingProposal() error {
	// Check resources first to see if staging is clean
	IsMCMSStagingAreaClean, err := utils.IsMCMSStagingAreaClean(op.AptosChain.Client, op.OnChainState.MCMSAddress)
	if err != nil {
		return fmt.Errorf("failed to check if MCMS staging area is clean: %w", err)
	}
	if IsMCMSStagingAreaClean {
		op.Env.Logger.Infow("MCMS Staging Area already clean", "addr", op.OnChainState.MCMSAddress.String())
		return nil
	}

	// Bind MCMS contract
	mcmsContract := mcmsbind.Bind(op.OnChainState.MCMSAddress, op.AptosChain.Client)
	mcmsAddress := mcmsContract.Address()

	// Get cleanup staging operations
	var operations []types.Operation
	moduleInfo, function, _, args, err := mcmsContract.MCMSDeployer().Encoder().CleanupStagingArea()
	if err != nil {
		return fmt.Errorf("failed to EncodeCleanupStagingArea: %w", err)
	}
	additionalFields := aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}
	afBytes, err := json.Marshal(additionalFields)
	if err != nil {
		return fmt.Errorf("failed to marshal additional fields: %w", err)
	}
	operations = append(operations, types.Operation{
		ChainSelector: types.ChainSelector(op.AptosChain.Selector),
		Transaction: types.Transaction{
			To:               mcmsAddress.StringLong(),
			Data:             aptosmcms.ArgsToData(args),
			AdditionalFields: afBytes,
		},
	})

	// Generate cleanup proposal
	proposal, nextOpCount, err := utils.GenerateProposal(
		op.AptosChain.Client,
		mcmsContract.Address(),
		op.AptosChain.Selector,
		operations,
		"Cleanup Staging Area",
		op.MCMSOpCount,
	)
	op.MCMSOpCount = nextOpCount
	if err != nil {
		return fmt.Errorf("failed to create deploy proposal: %w", err)
	}
	*op.Proposals = append(*op.Proposals, *proposal)

	return nil
}

// GenerateDeployCCIPProposal generates deployment MCMS operations for the CCIP package
func (op *CCIPDeploymentOperations) GenerateDeployCCIPProposal() (*aptos.AccountAddress, error) {
	// Validate there's no package deployed
	if op.OnChainState.CCIPAddress != aptos.AccountZero {
		op.Env.Logger.Infow("CCIP Package already deployed", "addr", op.OnChainState.CCIPAddress.String())
		return &op.OnChainState.CCIPAddress, nil
	}

	// Compile, chunk and get CCIP deploy operations
	mcmsContract := mcmsbind.Bind(op.OnChainState.MCMSAddress, op.AptosChain.Client)
	ccipObjectAddress, operations, err := op.getCCIPDeployOperations(mcmsContract, op.AptosChain.Selector)
	if err != nil {
		return nil, fmt.Errorf("failed to compile and create deploy operations: %w", err)
	}

	// Save the address of the CCIP object
	typeAndVersion := deployment.NewTypeAndVersion(changeset.AptosCCIPType, deployment.Version1_6_0)
	op.Ab.Save(op.AptosChain.Selector, ccipObjectAddress.String(), typeAndVersion)
	op.OnChainState.CCIPAddress = ccipObjectAddress

	// Generate deploy proposal
	proposal, nextOpCount, err := utils.GenerateProposal(op.AptosChain.Client, mcmsContract.Address(), op.AptosChain.Selector, operations, "Deploy CCIP Package", op.MCMSOpCount)
	op.MCMSOpCount = nextOpCount
	if err != nil {
		return nil, fmt.Errorf("failed to create deploy proposal: %w", err)
	}
	*op.Proposals = append(*op.Proposals, *proposal)

	return &ccipObjectAddress, nil
}

// GenerateDeployRouterProposal generates deployment MCMS operations for the Router module
func (op *CCIPDeploymentOperations) GenerateDeployRouterProposal(ccipObjectAddress *aptos.AccountAddress) error {
	// TODO: is there a way to check if module exists?
	// Compile, chunk and get Router deploy operations
	mcmsContract := mcmsbind.Bind(op.OnChainState.MCMSAddress, op.AptosChain.Client)
	operations, err := op.getRouterDeployOperations(mcmsContract, ccipObjectAddress)
	if err != nil {
		return fmt.Errorf("failed to compile and create deploy operations: %w", err)
	}

	// Generate deploy proposal
	proposal, nextOpCount, err := utils.GenerateProposal(op.AptosChain.Client, mcmsContract.Address(), op.AptosChain.Selector, operations, "Deploy Router Package", op.MCMSOpCount)
	op.MCMSOpCount = nextOpCount
	if err != nil {
		return fmt.Errorf("failed to create deploy proposal: %w", err)
	}
	*op.Proposals = append(*op.Proposals, *proposal)

	return nil
}

func (op *CCIPDeploymentOperations) GenerateInitializeCCIPProposal(ccipObjectAddress *aptos.AccountAddress) error {
	var operations []types.Operation
	ccipBind := ccip.Bind(*ccipObjectAddress, op.AptosChain.Client)

	// Config OnRamp
	moduleInfo, function, _, args, err := ccipBind.Onramp().Encoder().Initialize(
		op.AptosChain.Selector,
		op.CCIPConfig.OnRampParams.AllowlistAdmin,
		op.CCIPConfig.OnRampParams.DestChainSelectors,
		op.CCIPConfig.OnRampParams.DestChainEnabled,
		op.CCIPConfig.OnRampParams.DestChainAllowlistEnabled,
	)
	if err != nil {
		return fmt.Errorf("failed to encode onramp initialize: %w", err)
	}
	mcmsOp, err := generateMCMSOperation(op.AptosChain.Selector, ccipObjectAddress, moduleInfo, function, args)
	if err != nil {
		return fmt.Errorf("failed to generate MCMS operations for OnRamp Initialize: %w", err)
	}
	operations = append(operations, mcmsOp)

	// Config OffRamp
	moduleInfo, function, _, args, err = ccipBind.Offramp().Encoder().Initialize(
		op.AptosChain.Selector,
		op.CCIPConfig.OffRampParams.PermissionlessExecutionThreshold,
		op.CCIPConfig.OffRampParams.SourceChainSelectors,
		op.CCIPConfig.OffRampParams.SourceChainIsEnabled,
		op.CCIPConfig.OffRampParams.IsRMNVerificationDisabled,
		op.CCIPConfig.OffRampParams.SourceChainsOnRamp,
	)
	if err != nil {
		return fmt.Errorf("failed to encode offramp initialize: %w", err)
	}
	mcmsOp, err = generateMCMSOperation(op.AptosChain.Selector, ccipObjectAddress, moduleInfo, function, args)
	if err != nil {
		return fmt.Errorf("failed to generate MCMS operations for OffRamp Initialize: %w", err)
	}
	operations = append(operations, mcmsOp)

	// Config FeeQuoter
	moduleInfo, function, _, args, err = ccipBind.FeeQuoter().Encoder().Initialize(
		op.AptosChain.Selector,
		op.CCIPConfig.FeeQuoterParams.LinkToken,
		op.CCIPConfig.FeeQuoterParams.TokenPriceStalenessThreshold,
		op.CCIPConfig.FeeQuoterParams.FeeTokens,
	)
	if err != nil {
		return fmt.Errorf("failed to encode feequoter initialize: %w", err)
	}
	mcmsOp, err = generateMCMSOperation(op.AptosChain.Selector, ccipObjectAddress, moduleInfo, function, args)
	if err != nil {
		return fmt.Errorf("failed to generate MCMS operations for FeeQuoter Initialize: %w", err)
	}
	operations = append(operations, mcmsOp)

	// Config RMNRemote
	moduleInfo, function, _, args, err = ccipBind.RMNRemote().Encoder().Initialize(op.AptosChain.Selector)
	if err != nil {
		return fmt.Errorf("failed to encode rmnremote initialize: %w", err)
	}
	mcmsOp, err = generateMCMSOperation(op.AptosChain.Selector, ccipObjectAddress, moduleInfo, function, args)
	if err != nil {
		return fmt.Errorf("failed to generate MCMS operations for RMNRemote Initialize: %w", err)
	}
	operations = append(operations, mcmsOp)

	// Generate config proposal
	proposal, nextOpCount, err := utils.GenerateProposal(
		op.AptosChain.Client,
		op.OnChainState.MCMSAddress,
		op.AptosChain.Selector,
		operations,
		"Initialize CCIP Package",
		op.MCMSOpCount,
	)
	op.MCMSOpCount = nextOpCount
	if err != nil {
		return fmt.Errorf("failed to create deploy proposal: %w", err)
	}
	*op.Proposals = append(*op.Proposals, *proposal)

	return nil
}

func (op *CCIPDeploymentOperations) getCCIPDeployOperations(mcmsContract mcmsbind.MCMS, chainSel uint64) (aptos.AccountAddress, []types.Operation, error) {
	// Calculate addresses of the owner and the object
	ccipObjectAddress, err := mcmsContract.MCMSRegistry().GetNewCodeObjectAddress(nil, []byte(ccip.DefaultSeed))
	if err != nil {
		return ccipObjectAddress, []types.Operation{}, fmt.Errorf("failed to calculate object address: %w", err)
	}

	// Compile Package
	payload, err := ccip.Compile(ccipObjectAddress, mcmsContract.Address(), true)
	if err != nil {
		return ccipObjectAddress, []types.Operation{}, fmt.Errorf("failed to compile: %w", err)
	}

	// Create chunks and stage operations
	operations, err := utils.CreateChunksAndStage(payload, mcmsContract, chainSel, ccip.DefaultSeed, nil)
	if err != nil {
		return ccipObjectAddress, operations, fmt.Errorf("failed to create chunks and stage for %d: %w", chainSel, err)
	}

	return ccipObjectAddress, operations, nil
}

func (op *CCIPDeploymentOperations) getRouterDeployOperations(
	mcmsContract mcmsbind.MCMS,
	ccipObjectAddress *aptos.AccountAddress,
) ([]types.Operation, error) {
	// Compile Package
	payload, err := router.Compile(*ccipObjectAddress, mcmsContract.Address())
	if err != nil {
		return []types.Operation{}, fmt.Errorf("failed to compile: %w", err)
	}

	// Create chunks and stage operations
	operations, err := utils.CreateChunksAndStage(payload, mcmsContract, op.AptosChain.Selector, "", ccipObjectAddress)
	if err != nil {
		return operations, fmt.Errorf("failed to create chunks and stage for %d: %w", op.AptosChain.Selector, err)
	}

	return operations, nil
}

func generateMCMSOperation(chainSel uint64, toAddress *aptos.AccountAddress, moduleInfo bind.ModuleInformation, function string, args [][]byte) (types.Operation, error) {
	additionalFields := aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}
	afBytes, err := json.Marshal(additionalFields)
	if err != nil {
		return types.Operation{}, fmt.Errorf("failed to marshal additional fields: %w", err)
	}
	return types.Operation{
		ChainSelector: types.ChainSelector(chainSel),
		Transaction: types.Transaction{
			To:               toAddress.StringLong(),
			Data:             aptosmcms.ArgsToData(args),
			AdditionalFields: afBytes,
		},
	}, nil
}
