package operation

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_onramp"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router"
	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	aptoscfg "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
)

// OP: DeployCCIPOp deploys the CCIP package on Aptos chain
type DeployCCIPInput struct {
	MCMSAddress aptos.AccountAddress
	IsUpdate    bool
}

type DeployCCIPOutput struct {
	CCIPAddress    aptos.AccountAddress
	MCMSOperations []types.Operation
}

var DeployCCIPOp = operations.NewOperation(
	"deploy-ccip-op",
	Version1_0_0,
	"Deploys CCIP Package for Aptos Chain",
	deployCCIP,
)

func deployCCIP(b operations.Bundle, deps AptosDeps, in DeployCCIPInput) (DeployCCIPOutput, error) {
	onChainState := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector]
	// Validate there's no package deployed XOR is update
	if (onChainState.CCIPAddress == (aptos.AccountAddress{})) == (in.IsUpdate) {
		if in.IsUpdate {
			b.Logger.Infow("Trying to update a non-deployed package", "addr", onChainState.CCIPAddress.String())
			return DeployCCIPOutput{}, fmt.Errorf("CCIP package not deployed on Aptos chain %d", deps.AptosChain.Selector)
		}
		b.Logger.Infow("CCIP Package already deployed", "addr", onChainState.CCIPAddress.String())
		return DeployCCIPOutput{CCIPAddress: onChainState.CCIPAddress}, nil
	}

	// Compile, chunk and get CCIP deploy operations
	mcmsContract := mcmsbind.Bind(in.MCMSAddress, deps.AptosChain.Client)
	ccipObjectAddress, operations, err := getCCIPDeployMCMSOps(mcmsContract, deps.AptosChain.Selector, onChainState.CCIPAddress)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to compile and create deploy operations: %w", err)
	}
	if in.IsUpdate {
		return DeployCCIPOutput{
			MCMSOperations: operations,
		}, nil
	}

	return DeployCCIPOutput{
		CCIPAddress:    ccipObjectAddress,
		MCMSOperations: operations,
	}, nil
}

func getCCIPDeployMCMSOps(mcmsContract mcmsbind.MCMS, chainSel uint64, ccipAddress aptos.AccountAddress) (aptos.AccountAddress, []types.Operation, error) {
	// Calculate addresses of the owner and the object
	var ccipObjectAddress aptos.AccountAddress
	var err error
	if ccipAddress != (aptos.AccountAddress{}) {
		ccipObjectAddress = ccipAddress
	} else {
		ccipObjectAddress, err = mcmsContract.MCMSRegistry().GetNewCodeObjectAddress(nil, []byte(ccip.DefaultSeed))
		if err != nil {
			return ccipObjectAddress, []types.Operation{}, fmt.Errorf("failed to calculate object address: %w", err)
		}
	}

	// Compile Package
	payload, err := ccip.Compile(ccipObjectAddress, mcmsContract.Address(), ccipAddress == aptos.AccountAddress{})
	if err != nil {
		return ccipObjectAddress, []types.Operation{}, fmt.Errorf("failed to compile: %w", err)
	}

	// Create chunks and stage operations
	var operations []types.Operation
	if ccipAddress == (aptos.AccountAddress{}) {
		operations, err = utils.CreateChunksAndStage(payload, mcmsContract, chainSel, ccip.DefaultSeed, nil)
	} else {
		operations, err = utils.CreateChunksAndStage(payload, mcmsContract, chainSel, "", &ccipObjectAddress)
	}
	if err != nil {
		return ccipObjectAddress, operations, fmt.Errorf("failed to create chunks and stage for %d: %w", chainSel, err)
	}

	return ccipObjectAddress, operations, nil
}

// DeployModulesInput is the input for every module deployment operation
type DeployModulesInput struct {
	MCMSAddress aptos.AccountAddress
	CCIPAddress aptos.AccountAddress
}

// OP: DeployRouterOp generates deployment MCMS operations for the Router module
var DeployRouterOp = operations.NewOperation(
	"deploy-router-op",
	Version1_0_0,
	"Generates MCMS proposals that deployes Router module on CCIP package",
	deployRouter,
)

func deployRouter(b operations.Bundle, deps AptosDeps, in DeployModulesInput) ([]types.Operation, error) {
	// TODO: is there a way to check if module exists?
	mcmsContract := mcmsbind.Bind(in.MCMSAddress, deps.AptosChain.Client)
	// Compile Package
	payload, err := ccip_router.Compile(in.CCIPAddress, mcmsContract.Address(), true)
	if err != nil {
		return []types.Operation{}, fmt.Errorf("failed to compile: %w", err)
	}
	// Create chunks and stage operations
	operations, err := utils.CreateChunksAndStage(payload, mcmsContract, deps.AptosChain.Selector, "", &in.CCIPAddress)
	if err != nil {
		return operations, fmt.Errorf("failed to create chunks and stage for %d: %w", deps.AptosChain.Selector, err)
	}

	return operations, nil
}

// OP: DeployOffRampOp generates deployment MCMS operations for the OffRamp module
var DeployOffRampOp = operations.NewOperation(
	"deploy-offramp-op",
	Version1_0_0,
	"Generates MCMS proposals that deployes OffRamp module on CCIP package",
	deployOffRamp,
)

func deployOffRamp(b operations.Bundle, deps AptosDeps, in DeployModulesInput) ([]types.Operation, error) {
	mcmsContract := mcmsbind.Bind(in.MCMSAddress, deps.AptosChain.Client)
	// Compile Package
	payload, err := ccip_offramp.Compile(in.CCIPAddress, mcmsContract.Address(), true)
	if err != nil {
		return []types.Operation{}, fmt.Errorf("failed to compile: %w", err)
	}
	// Create chunks and stage operations
	operations, err := utils.CreateChunksAndStage(payload, mcmsContract, deps.AptosChain.Selector, "", &in.CCIPAddress)
	if err != nil {
		return operations, fmt.Errorf("failed to create chunks and stage for %d: %w", deps.AptosChain.Selector, err)
	}
	return operations, nil
}

// OP: DeployOnRampOp generates MCMS proposals for the OnRamp module deployment
var DeployOnRampOp = operations.NewOperation(
	"deploy-onramp-op",
	Version1_0_0,
	"Generates MCMS proposals that deployes OnRamp module on CCIP package",
	deployOnRamp,
)

func deployOnRamp(b operations.Bundle, deps AptosDeps, in DeployModulesInput) ([]types.Operation, error) {
	mcmsContract := mcmsbind.Bind(in.MCMSAddress, deps.AptosChain.Client)
	// Compile Package
	payload, err := ccip_onramp.Compile(in.CCIPAddress, mcmsContract.Address(), true)
	if err != nil {
		return []types.Operation{}, fmt.Errorf("failed to compile: %w", err)
	}
	// Create chunks and stage operations
	operations, err := utils.CreateChunksAndStage(payload, mcmsContract, deps.AptosChain.Selector, "", &in.CCIPAddress)
	if err != nil {
		return operations, fmt.Errorf("failed to create chunks and stage for %d: %w", deps.AptosChain.Selector, err)
	}
	return operations, nil
}

// OP: InitializeCCIP Operation
type InitializeCCIPInput struct {
	MCMSAddress aptos.AccountAddress
	CCIPAddress aptos.AccountAddress
	CCIPConfig  aptoscfg.ChainContractParams
}

var InitializeCCIPOp = operations.NewOperation(
	"initialize-ccip-op",
	Version1_0_0,
	"Initializes CCIP components with configuration parameters",
	initializeCCIP,
)

func initializeCCIP(b operations.Bundle, deps AptosDeps, in InitializeCCIPInput) (types.BatchOperation, error) {
	var txs []types.Transaction

	// Config OnRamp with empty lane configs. We're only able to get router address after deploying the router module
	onrampBind := ccip_onramp.Bind(in.CCIPAddress, deps.AptosChain.Client)
	moduleInfo, function, _, args, err := onrampBind.Onramp().Encoder().Initialize(
		deps.AptosChain.Selector,
		in.CCIPConfig.OnRampParams.FeeAggregator,
		in.CCIPConfig.OnRampParams.AllowlistAdmin,
		[]uint64{},
		[]aptos.AccountAddress{},
		[]bool{},
	)
	if err != nil {
		return types.BatchOperation{}, fmt.Errorf("failed to encode onramp initialize: %w", err)
	}
	mcmsTx, err := utils.GenerateMCMSTx(in.CCIPAddress, moduleInfo, function, args)
	if err != nil {
		return types.BatchOperation{}, fmt.Errorf("failed to generate MCMS operations for OnRamp Initialize: %w", err)
	}
	txs = append(txs, mcmsTx)

	// Config OffRamp
	offrampBind := ccip_offramp.Bind(in.CCIPAddress, deps.AptosChain.Client)
	moduleInfo, function, _, args, err = offrampBind.Offramp().Encoder().Initialize(
		deps.AptosChain.Selector,
		in.CCIPConfig.OffRampParams.PermissionlessExecutionThreshold,
		in.CCIPConfig.OffRampParams.SourceChainSelectors,
		in.CCIPConfig.OffRampParams.SourceChainIsEnabled,
		in.CCIPConfig.OffRampParams.IsRMNVerificationDisabled,
		in.CCIPConfig.OffRampParams.SourceChainsOnRamp,
	)
	if err != nil {
		return types.BatchOperation{}, fmt.Errorf("failed to encode offramp initialize: %w", err)
	}
	mcmsTx, err = utils.GenerateMCMSTx(in.CCIPAddress, moduleInfo, function, args)
	if err != nil {
		return types.BatchOperation{}, fmt.Errorf("failed to generate MCMS operations for OffRamp Initialize: %w", err)
	}
	txs = append(txs, mcmsTx)

	// Config FeeQuoter and RMNRemote
	ccipBind := ccip.Bind(in.CCIPAddress, deps.AptosChain.Client)

	moduleInfo, function, _, args, err = ccipBind.FeeQuoter().Encoder().Initialize(
		deps.AptosChain.Selector,
		in.CCIPConfig.FeeQuoterParams.LinkToken,
		in.CCIPConfig.FeeQuoterParams.TokenPriceStalenessThreshold,
		in.CCIPConfig.FeeQuoterParams.FeeTokens,
	)
	if err != nil {
		return types.BatchOperation{}, fmt.Errorf("failed to encode feequoter initialize: %w", err)
	}
	mcmsTx, err = utils.GenerateMCMSTx(in.CCIPAddress, moduleInfo, function, args)
	if err != nil {
		return types.BatchOperation{}, fmt.Errorf("failed to generate MCMS operations for FeeQuoter Initialize: %w", err)
	}
	txs = append(txs, mcmsTx)

	moduleInfo, function, _, args, err = ccipBind.RMNRemote().Encoder().Initialize(deps.AptosChain.Selector)
	if err != nil {
		return types.BatchOperation{}, fmt.Errorf("failed to encode rmnremote initialize: %w", err)
	}
	mcmsTx, err = utils.GenerateMCMSTx(in.CCIPAddress, moduleInfo, function, args)
	if err != nil {
		return types.BatchOperation{}, fmt.Errorf("failed to generate MCMS operations for RMNRemote Initialize: %w", err)
	}
	txs = append(txs, mcmsTx)

	return types.BatchOperation{
		ChainSelector: types.ChainSelector(deps.AptosChain.Selector),
		Transactions:  txs,
	}, nil
}
