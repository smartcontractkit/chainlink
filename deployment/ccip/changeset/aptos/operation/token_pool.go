package operation

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/managed_token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/compile"
	"github.com/smartcontractkit/chainlink-aptos/bindings/managed_token"
	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

type DeployTokenPoolPackageOutput struct {
	TokenPoolObjectAddress aptos.AccountAddress
	MCMSOps                []types.Operation
}

// DeployTokenPoolPackageOp deploys token pool package to Token Object Address
var DeployTokenPoolPackageOp = operations.NewOperation(
	"deploy-token-pool-package-op",
	Version1_0_0,
	"Deploy Aptos token pool package",
	deployTokenPoolPackage,
)

func deployTokenPoolPackage(b operations.Bundle, deps AptosDeps, poolSeed string) (DeployTokenPoolPackageOutput, error) {
	aptosState := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector]
	mcmsContract := mcmsbind.Bind(aptosState.MCMSAddress, deps.AptosChain.Client)

	// Calculate pool address
	tokenPoolObjectAddress, err := mcmsContract.MCMSRegistry().GetNewCodeObjectAddress(nil, []byte(poolSeed))
	if err != nil {
		return DeployTokenPoolPackageOutput{}, fmt.Errorf("failed to GetNewCodeObjectAddress for pool seed %s: %w", poolSeed, err)
	}

	payload, err := token_pool.Compile(tokenPoolObjectAddress, aptosState.CCIPAddress, aptosState.MCMSAddress)
	if err != nil {
		return DeployTokenPoolPackageOutput{}, fmt.Errorf("failed to compile token pool: %w", err)
	}
	ops, err := utils.CreateChunksAndStage(payload, mcmsContract, deps.AptosChain.Selector, poolSeed, nil)
	if err != nil {
		return DeployTokenPoolPackageOutput{}, fmt.Errorf("failed to create chunks for token pool: %w", err)
	}

	return DeployTokenPoolPackageOutput{
		TokenPoolObjectAddress: tokenPoolObjectAddress,
		MCMSOps:                ops,
	}, nil
}

type DeployTokenPoolModuleInput struct {
	PoolType             cldf.ContractType
	TokenObjAddress      aptos.AccountAddress // TODO change this to metadata address, and determine object if needed
	TokenPoolObjAddress  aptos.AccountAddress
	InitialAdministrator aptos.AccountAddress
}

// DeployTokenPoolModuleOp deploys token pool module to Token Object Address
var DeployTokenPoolModuleOp = operations.NewOperation(
	"deploy-token-pool-module-op",
	Version1_0_0,
	"Deploy Aptos token pool module",
	deployTokenPoolModule,
)

func deployTokenPoolModule(b operations.Bundle, deps AptosDeps, in DeployTokenPoolModuleInput) ([]types.Operation, error) {
	aptosState := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector]
	mcmsContract := mcmsbind.Bind(aptosState.MCMSAddress, deps.AptosChain.Client)

	var ops []types.Operation

	var (
		payload compile.CompiledPackage
		err     error
	)
	switch in.PoolType {
	case shared.AptosManagedTokenPoolType:
		payload, err = managed_token_pool.Compile(
			in.TokenPoolObjAddress,
			aptosState.CCIPAddress,
			aptosState.MCMSAddress,
			in.TokenPoolObjAddress,
			in.TokenObjAddress,
			in.InitialAdministrator,
			true,
		)
	case shared.BurnMintTokenPool:
		payload, err = burn_mint_token_pool.Compile(
			in.TokenPoolObjAddress,
			aptosState.CCIPAddress,
			aptosState.MCMSAddress,
			in.TokenPoolObjAddress,
			in.TokenObjAddress, // TODO this should be the metadata address
			in.InitialAdministrator,
			true,
		)
	case shared.LockReleaseTokenPool:
		payload, err = lock_release_token_pool.Compile(
			in.TokenPoolObjAddress,
			aptosState.CCIPAddress,
			aptosState.MCMSAddress,
			in.TokenPoolObjAddress,
			in.TokenObjAddress, // TODO this should be the metadata address
			in.InitialAdministrator,
			true,
		)
	default:
		return nil, fmt.Errorf("invalid token pool type: %s", in.PoolType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to compile token pool: %w", err)
	}
	ops, err = utils.CreateChunksAndStage(payload, mcmsContract, deps.AptosChain.Selector, "", &in.TokenPoolObjAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create chunks for token pool: %w", err)
	}

	return ops, nil
}

// GrantMinterPermissionsOp operation to grant minter permissions
var GrantMinterPermissionsOp = operations.NewOperation(
	"grant-minter-permissions-op",
	Version1_0_0,
	"Grant Minter permissions to the token pool state address",
	grantMinterPermissions,
)

type GrantRolePermissionsInput struct {
	TokenObjAddress       aptos.AccountAddress
	TokenPoolStateAddress aptos.AccountAddress
}

func grantMinterPermissions(b operations.Bundle, deps AptosDeps, in GrantRolePermissionsInput) (types.Transaction, error) {
	tokenContract := managed_token.Bind(in.TokenObjAddress, deps.AptosChain.Client)

	moduleInfo, function, _, args, err := tokenContract.ManagedToken().Encoder().ApplyAllowedMinterUpdates([]aptos.AccountAddress{}, []aptos.AccountAddress{in.TokenPoolStateAddress})
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode ApplyAllowedMinterUpdates: %w", err)
	}

	return utils.GenerateMCMSTx(in.TokenObjAddress, moduleInfo, function, args)
}

// GrantBurnerPermissionsOp operation to grant burner permissions
var GrantBurnerPermissionsOp = operations.NewOperation(
	"grant-burner-permissions-op",
	Version1_0_0,
	"Grant Burner permissions to the token pool state address",
	grantBurnerPermissions,
)

func grantBurnerPermissions(b operations.Bundle, deps AptosDeps, in GrantRolePermissionsInput) (types.Transaction, error) {
	tokenContract := managed_token.Bind(in.TokenObjAddress, deps.AptosChain.Client)

	moduleInfo, function, _, args, err := tokenContract.ManagedToken().Encoder().ApplyAllowedBurnerUpdates([]aptos.AccountAddress{}, []aptos.AccountAddress{in.TokenPoolStateAddress})
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode ApplyAllowedBurnerUpdates: %w", err)
	}

	return utils.GenerateMCMSTx(in.TokenObjAddress, moduleInfo, function, args)
}

type ApplyChainUpdatesInput struct {
	RemoteChainSelectorsToRemove []uint64
	RemoteChainSelectorsToAdd    []uint64
	RemotePoolAddresses          [][][]byte
	RemoteTokenAddresses         [][]byte
	TokenPoolAddress             aptos.AccountAddress
}

// ApplyChainUpdatesOp ...
var ApplyChainUpdatesOp = operations.NewOperation(
	"apply-chain-updates-op",
	Version1_0_0,
	"Apply chain updates to Aptos token pool",
	applyChainUpdates,
)

func applyChainUpdates(b operations.Bundle, deps AptosDeps, in ApplyChainUpdatesInput) (types.Transaction, error) {
	poolBind := managed_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
	moduleInfo, function, _, args, err := poolBind.ManagedTokenPool().Encoder().ApplyChainUpdates(
		in.RemoteChainSelectorsToRemove,
		in.RemoteChainSelectorsToAdd,
		in.RemotePoolAddresses,
		in.RemoteTokenAddresses,
	)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode ApplyChainUpdates for chains: %w", err)
	}

	return utils.GenerateMCMSTx(in.TokenPoolAddress, moduleInfo, function, args)
}

type SetChainRLConfigsInput struct {
	RemoteChainSelectors []uint64
	OutboundIsEnableds   []bool
	OutboundCapacities   []uint64
	OutboundRates        []uint64
	InboundIsEnableds    []bool
	InboundCapacities    []uint64
	InboundRates         []uint64
	TokenPoolAddress     aptos.AccountAddress
}

var SetChainRateLimiterConfigsOp = operations.NewOperation(
	"set-chain-rate-limiter-configs-op",
	Version1_0_0,
	"Set chain rate limiter configs for Aptos token pool",
	setChainRateLimiterConfigs,
)

func setChainRateLimiterConfigs(b operations.Bundle, deps AptosDeps, in SetChainRLConfigsInput) (types.Transaction, error) {
	poolBind := managed_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
	moduleInfo, function, _, args, err := poolBind.ManagedTokenPool().Encoder().SetChainRateLimiterConfigs(
		in.RemoteChainSelectors,
		in.OutboundIsEnableds,
		in.OutboundCapacities,
		in.OutboundRates,
		in.InboundIsEnableds,
		in.InboundCapacities,
		in.InboundRates,
	)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode SetChainRateLimiterConfigs for chains: %w", err)
	}

	return utils.GenerateMCMSTx(in.TokenPoolAddress, moduleInfo, function, args)
}
