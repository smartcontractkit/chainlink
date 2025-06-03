package operation

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/managed_token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/token_pool"
	managed_token "github.com/smartcontractkit/chainlink-aptos/bindings/managed_token"
	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/mcms/types"
)

// DeployTokenPoolPackageOp deploys token pool package to Token Object Address
var DeployTokenPoolPackageOp = operations.NewOperation(
	"deploy-token-pool-package-op",
	Version1_0_0,
	"Deploy Aptos token pool package",
	deployTokenPoolPackage,
)

func deployTokenPoolPackage(b operations.Bundle, deps AptosDeps, tokenObjAddress aptos.AccountAddress) ([]types.Operation, error) {
	aptosState := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector]
	mcmsContract := mcmsbind.Bind(aptosState.MCMSAddress, deps.AptosChain.Client)

	payload, err := token_pool.Compile(tokenObjAddress, aptosState.CCIPAddress, aptosState.MCMSAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to compile token pool: %w", err)
	}
	ops, err := utils.CreateChunksAndStage(payload, mcmsContract, deps.AptosChain.Selector, "", &tokenObjAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create chunks for token pool: %w", err)
	}

	return ops, nil
}

type DeployTokenPoolModuleInput struct {
	PoolType        cldf.ContractType
	TokenObjAddress aptos.AccountAddress
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

	ccipOwnerAddr, err := mcmsContract.MCMSRegistry().GetRegisteredOwnerAddress(nil, aptosState.CCIPAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get preexisting code object owner address: %w", err)
	}

	var ops []types.Operation

	switch in.PoolType {
	case shared.BurnMintTokenPool:
		payload, err := managed_token_pool.Compile(
			in.TokenObjAddress,
			aptosState.CCIPAddress,
			aptosState.MCMSAddress,
			in.TokenObjAddress,
			in.TokenObjAddress,
			ccipOwnerAddr,
			true,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to compile token pool: %w", err)
		}
		ops, err = utils.CreateChunksAndStage(payload, mcmsContract, deps.AptosChain.Selector, "", &in.TokenObjAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to create chunks for token pool: %w", err)
		}
	case shared.LockReleaseTokenPool:
		payload, err := lock_release_token_pool.Compile(
			in.TokenObjAddress,
			aptosState.CCIPAddress,
			aptosState.MCMSAddress,
			in.TokenObjAddress,
			in.TokenObjAddress,
			ccipOwnerAddr,
			true,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to compile token pool: %w", err)
		}
		ops, err = utils.CreateChunksAndStage(payload, mcmsContract, deps.AptosChain.Selector, "", &in.TokenObjAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to create chunks for token pool: %w", err)
		}
	default:
		return nil, fmt.Errorf("invalid token pool type: %s", in.PoolType)
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

func grantMinterPermissions(b operations.Bundle, deps AptosDeps, tokenObjAddress aptos.AccountAddress) (types.Transaction, error) {
	managedTokenPoolStateAddress := tokenObjAddress.ResourceAccount([]byte("CcipManagedTokenPool"))
	tokenContract := managed_token.Bind(tokenObjAddress, deps.AptosChain.Client)

	moduleInfo, function, _, args, err := tokenContract.ManagedToken().Encoder().ApplyAllowedMinterUpdates([]aptos.AccountAddress{}, []aptos.AccountAddress{managedTokenPoolStateAddress})
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode ApplyAllowedMinterUpdates: %w", err)
	}

	return utils.GenerateMCMSTx(tokenObjAddress, moduleInfo, function, args)
}

// GrantBurnerPermissionsOp operation to grant burner permissions
var GrantBurnerPermissionsOp = operations.NewOperation(
	"grant-burner-permissions-op",
	Version1_0_0,
	"Grant Burner permissions to the token pool state address",
	grantBurnerPermissions,
)

func grantBurnerPermissions(b operations.Bundle, deps AptosDeps, tokenObjAddress aptos.AccountAddress) (types.Transaction, error) {
	managedTokenPoolStateAddress := tokenObjAddress.ResourceAccount([]byte("CcipManagedTokenPool"))
	tokenContract := managed_token.Bind(tokenObjAddress, deps.AptosChain.Client)

	moduleInfo, function, _, args, err := tokenContract.ManagedToken().Encoder().ApplyAllowedBurnerUpdates([]aptos.AccountAddress{}, []aptos.AccountAddress{managedTokenPoolStateAddress})
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode ApplyAllowedBurnerUpdates: %w", err)
	}

	return utils.GenerateMCMSTx(tokenObjAddress, moduleInfo, function, args)
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
