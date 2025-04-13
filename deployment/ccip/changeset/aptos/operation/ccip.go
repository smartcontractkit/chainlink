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
	"github.com/smartcontractkit/chainlink/deployment/operations"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	"github.com/smartcontractkit/mcms/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

// CleanupStagingArea Operation
type CleanupStagingAreaInput struct {
	MCMSAddress aptos.AccountAddress
}

var CleanupStagingAreaOp = operations.NewOperation(
	"cleanup-staging-area-op",
	Version1_0_0,
	"Cleans up MCMS staging area if it's not already clean",
	cleanupStagingArea,
)

func cleanupStagingArea(b operations.Bundle, deps AptosDeps, in CleanupStagingAreaInput) ([]mcmstypes.Operation, error) {
	// Check resources first to see if staging is clean
	IsMCMSStagingAreaClean, err := utils.IsMCMSStagingAreaClean(deps.AptosChain.Client, in.MCMSAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to check if MCMS staging area is clean: %w", err)
	}
	if IsMCMSStagingAreaClean {
		b.Logger.Infow("MCMS Staging Area already clean", "addr", in.MCMSAddress.String())
		return nil, nil
	}

	// Bind MCMS contract
	mcmsContract := mcmsbind.Bind(in.MCMSAddress, deps.AptosChain.Client)
	mcmsAddress := mcmsContract.Address()

	// Get cleanup staging operations
	var operations []types.Operation
	moduleInfo, function, _, args, err := mcmsContract.MCMSDeployer().Encoder().CleanupStagingArea()
	if err != nil {
		return nil, fmt.Errorf("failed to EncodeCleanupStagingArea: %w", err)
	}
	additionalFields := aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}
	afBytes, err := json.Marshal(additionalFields)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal additional fields: %w", err)
	}
	operations = append(operations, types.Operation{
		ChainSelector: types.ChainSelector(deps.AptosChain.Selector),
		Transaction: types.Transaction{
			To:               mcmsAddress.StringLong(),
			Data:             aptosmcms.ArgsToData(args),
			AdditionalFields: afBytes,
		},
	})

	return operations, nil
}

// GenerateDeployCCIPProposal Operation generates deployment MCMS operations for the CCIP package
type DeployCCIPInput struct {
	MCMSAddress aptos.AccountAddress
}

type DeployCCIPOutput struct {
	CCIPAddress    aptos.AccountAddress
	MCMSOperations []mcmstypes.Operation
}

var GenerateDeployCCIPProposalOp = operations.NewOperation(
	"deploy-ccip-op",
	Version1_0_0,
	"Deploys CCIP Package for Aptos Chain",
	generateDeployCCIPProposal,
)

func generateDeployCCIPProposal(b operations.Bundle, deps AptosDeps, in DeployCCIPInput) (DeployCCIPOutput, error) {
	// Validate there's no package deployed
	if deps.OnChainState.CCIPAddress != aptos.AccountZero {
		b.Logger.Infow("CCIP Package already deployed", "addr", deps.OnChainState.CCIPAddress.String())
		return DeployCCIPOutput{CCIPAddress: deps.OnChainState.CCIPAddress}, nil
	}

	// Compile, chunk and get CCIP deploy operations
	mcmsContract := mcmsbind.Bind(in.MCMSAddress, deps.AptosChain.Client)
	ccipObjectAddress, operations, err := getCCIPDeployMCMSOps(mcmsContract, deps.AptosChain.Selector)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to compile and create deploy operations: %w", err)
	}

	// Save the address of the CCIP object
	typeAndVersion := deployment.NewTypeAndVersion(changeset.AptosCCIPType, deployment.Version1_6_0)
	deps.AB.Save(deps.AptosChain.Selector, ccipObjectAddress.String(), typeAndVersion)
	deps.OnChainState.CCIPAddress = ccipObjectAddress

	return DeployCCIPOutput{
		CCIPAddress:    ccipObjectAddress,
		MCMSOperations: operations,
	}, nil
}

func getCCIPDeployMCMSOps(mcmsContract mcmsbind.MCMS, chainSel uint64) (aptos.AccountAddress, []types.Operation, error) {
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

// GenerateDeployRouterProposal generates deployment MCMS operations for the Router module
type DeployRouterInput struct {
	MCMSAddress aptos.AccountAddress
	CCIPAddress aptos.AccountAddress
}

var GenerateDeployRouterProposalOp = operations.NewOperation(
	"deploy-router-op",
	Version1_0_0,
	"Deploys Router Package for CCIP",
	generateDeployRouterProposal,
)

func generateDeployRouterProposal(b operations.Bundle, deps AptosDeps, in DeployRouterInput) ([]mcmstypes.Operation, error) {
	// TODO: is there a way to check if module exists?
	// Compile, chunk and get Router deploy operations
	mcmsContract := mcmsbind.Bind(in.MCMSAddress, deps.AptosChain.Client)
	operations, err := getRouterDeployMCMSOps(mcmsContract, in.CCIPAddress, deps.AptosChain.Selector)
	if err != nil {
		return nil, fmt.Errorf("failed to compile and create deploy operations: %w", err)
	}

	return operations, nil
}

func getRouterDeployMCMSOps(
	mcmsContract mcmsbind.MCMS,
	ccipObjectAddress aptos.AccountAddress,
	chainSel uint64,
) ([]types.Operation, error) {
	// Compile Package
	payload, err := router.Compile(ccipObjectAddress, mcmsContract.Address())
	if err != nil {
		return []types.Operation{}, fmt.Errorf("failed to compile: %w", err)
	}

	// Create chunks and stage operations
	operations, err := utils.CreateChunksAndStage(payload, mcmsContract, chainSel, "", &ccipObjectAddress)
	if err != nil {
		return operations, fmt.Errorf("failed to create chunks and stage for %d: %w", chainSel, err)
	}

	return operations, nil
}

// InitializeCCIP Operation
type InitializeCCIPInput struct {
	MCMSAddress aptos.AccountAddress
	CCIPAddress aptos.AccountAddress
	CCIPConfig  aptoscfg.ChainContractParams
}

var InitializeCCIPOp = operations.NewOperation(
	"initialize-ccip-op",
	Version1_0_0,
	"Initializes CCIP components with configuration parameters",
	generateInitializeCCIPProposal,
)

func generateInitializeCCIPProposal(b operations.Bundle, deps AptosDeps, in InitializeCCIPInput) ([]types.Operation, error) {
	var operations []types.Operation
	ccipBind := ccip.Bind(in.CCIPAddress, deps.AptosChain.Client)

	// Config OnRamp
	moduleInfo, function, _, args, err := ccipBind.Onramp().Encoder().Initialize(
		deps.AptosChain.Selector,
		in.CCIPConfig.OnRampParams.AllowlistAdmin,
		in.CCIPConfig.OnRampParams.DestChainSelectors,
		in.CCIPConfig.OnRampParams.DestChainEnabled,
		in.CCIPConfig.OnRampParams.DestChainAllowlistEnabled,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode onramp initialize: %w", err)
	}
	mcmsOp, err := generateMCMSOperation(deps.AptosChain.Selector, in.CCIPAddress, moduleInfo, function, args)
	if err != nil {
		return nil, fmt.Errorf("failed to generate MCMS operations for OnRamp Initialize: %w", err)
	}
	operations = append(operations, mcmsOp)

	// Config OffRamp
	moduleInfo, function, _, args, err = ccipBind.Offramp().Encoder().Initialize(
		deps.AptosChain.Selector,
		in.CCIPConfig.OffRampParams.PermissionlessExecutionThreshold,
		in.CCIPConfig.OffRampParams.SourceChainSelectors,
		in.CCIPConfig.OffRampParams.SourceChainIsEnabled,
		in.CCIPConfig.OffRampParams.IsRMNVerificationDisabled,
		in.CCIPConfig.OffRampParams.SourceChainsOnRamp,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode offramp initialize: %w", err)
	}
	mcmsOp, err = generateMCMSOperation(deps.AptosChain.Selector, in.CCIPAddress, moduleInfo, function, args)
	if err != nil {
		return nil, fmt.Errorf("failed to generate MCMS operations for OffRamp Initialize: %w", err)
	}
	operations = append(operations, mcmsOp)

	// Config FeeQuoter
	moduleInfo, function, _, args, err = ccipBind.FeeQuoter().Encoder().Initialize(
		deps.AptosChain.Selector,
		in.CCIPConfig.FeeQuoterParams.LinkToken,
		in.CCIPConfig.FeeQuoterParams.TokenPriceStalenessThreshold,
		in.CCIPConfig.FeeQuoterParams.FeeTokens,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode feequoter initialize: %w", err)
	}
	mcmsOp, err = generateMCMSOperation(deps.AptosChain.Selector, in.CCIPAddress, moduleInfo, function, args)
	if err != nil {
		return nil, fmt.Errorf("failed to generate MCMS operations for FeeQuoter Initialize: %w", err)
	}
	operations = append(operations, mcmsOp)

	// Config RMNRemote
	moduleInfo, function, _, args, err = ccipBind.RMNRemote().Encoder().Initialize(deps.AptosChain.Selector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode rmnremote initialize: %w", err)
	}
	mcmsOp, err = generateMCMSOperation(deps.AptosChain.Selector, in.CCIPAddress, moduleInfo, function, args)
	if err != nil {
		return nil, fmt.Errorf("failed to generate MCMS operations for RMNRemote Initialize: %w", err)
	}
	operations = append(operations, mcmsOp)

	return operations, nil

}

// generateMCMSOperation is a helper function that generates a MCMS operation for the given parameters
func generateMCMSOperation(chainSel uint64, toAddress aptos.AccountAddress, moduleInfo bind.ModuleInformation, function string, args [][]byte) (types.Operation, error) {
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
