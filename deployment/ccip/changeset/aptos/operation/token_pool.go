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
	PoolType            cldf.ContractType
	TokenCodeObjAddress aptos.AccountAddress
	TokenAddress        aptos.AccountAddress
	TokenPoolObjAddress aptos.AccountAddress
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
			in.TokenCodeObjAddress,
			true,
		)
	case shared.BurnMintTokenPool:
		payload, err = burn_mint_token_pool.Compile(
			in.TokenPoolObjAddress,
			aptosState.CCIPAddress,
			aptosState.MCMSAddress,
			in.TokenPoolObjAddress,
			in.TokenAddress,
			true,
		)
	case shared.LockReleaseTokenPool:
		payload, err = lock_release_token_pool.Compile(
			in.TokenPoolObjAddress,
			aptosState.CCIPAddress,
			aptosState.MCMSAddress,
			in.TokenPoolObjAddress,
			in.TokenAddress,
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

type AddRemotePoolsInput struct {
	TokenPoolAddress     aptos.AccountAddress
	RemoteChainSelectors []uint64
	RemotePoolAddresses  [][]byte
}

var AddRemotePoolsOp = operations.NewOperation(
	"add-remote-pools-op",
	Version1_0_0,
	"Adds new remote pools to an Aptos token pool",
	addRemotePools,
)

func addRemotePools(b operations.Bundle, deps AptosDeps, in AddRemotePoolsInput) ([]types.Transaction, error) {
	txs := make([]types.Transaction, len(in.RemoteChainSelectors))
	poolBind := managed_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
	for i, selector := range in.RemoteChainSelectors {
		moduleInfo, function, _, args, err := poolBind.ManagedTokenPool().Encoder().AddRemotePool(selector, in.RemotePoolAddresses[i])
		if err != nil {
			return nil, fmt.Errorf("failed to encode AddRemotePools for remote selector %d: %w", selector, err)
		}
		tx, err := utils.GenerateMCMSTx(in.TokenPoolAddress, moduleInfo, function, args)
		if err != nil {
			return nil, err
		}
		txs[i] = tx
	}
	return txs, nil
}
