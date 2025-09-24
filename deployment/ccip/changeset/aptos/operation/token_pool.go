package operation

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/managed_token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/regulated_token_pool"
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
	case shared.AptosRegulatedTokenPoolType:
		payload, err = regulated_token_pool.Compile(
			in.TokenPoolObjAddress,
			aptosState.CCIPAddress,
			aptosState.MCMSAddress,
			in.TokenPoolObjAddress,
			in.TokenCodeObjAddress,
			deps.AptosChain.DeployerSigner.AccountAddress(), // Unused parameter, since the admin is set on the token not the pool
			true,
		)
	default:
		return nil, fmt.Errorf("unsupported token pool type for DeployTokenPoolModuleOp: %s", in.PoolType.String())
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
	TokenPoolAddress             aptos.AccountAddress
	TokenPoolType                cldf.ContractType
	RemoteChainSelectorsToRemove []uint64
	RemoteChainSelectorsToAdd    []uint64
	RemotePoolAddresses          [][][]byte
	RemoteTokenAddresses         [][]byte
}

// ApplyChainUpdatesOp ...
var ApplyChainUpdatesOp = operations.NewOperation(
	"apply-chain-updates-op",
	Version1_0_0,
	"Apply chain updates to an Aptos token pool",
	applyChainUpdates,
)

func applyChainUpdates(b operations.Bundle, deps AptosDeps, in ApplyChainUpdatesInput) (types.Transaction, error) {
	var (
		moduleInfo bind.ModuleInformation
		function   string
		args       [][]byte
		err        error
	)

	switch in.TokenPoolType {
	case shared.AptosManagedTokenPoolType:
		poolBind := managed_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		moduleInfo, function, _, args, err = poolBind.ManagedTokenPool().Encoder().ApplyChainUpdates(
			in.RemoteChainSelectorsToRemove,
			in.RemoteChainSelectorsToAdd,
			in.RemotePoolAddresses,
			in.RemoteTokenAddresses,
		)
	case shared.AptosRegulatedTokenPoolType:
		poolBind := regulated_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		moduleInfo, function, _, args, err = poolBind.RegulatedTokenPool().Encoder().ApplyChainUpdates(
			in.RemoteChainSelectorsToRemove,
			in.RemoteChainSelectorsToAdd,
			in.RemotePoolAddresses,
			in.RemoteTokenAddresses,
		)
	case shared.BurnMintTokenPool:
		poolBind := burn_mint_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		moduleInfo, function, _, args, err = poolBind.BurnMintTokenPool().Encoder().ApplyChainUpdates(
			in.RemoteChainSelectorsToRemove,
			in.RemoteChainSelectorsToAdd,
			in.RemotePoolAddresses,
			in.RemoteTokenAddresses,
		)
	case shared.LockReleaseTokenPool:
		poolBind := lock_release_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		moduleInfo, function, _, args, err = poolBind.LockReleaseTokenPool().Encoder().ApplyChainUpdates(
			in.RemoteChainSelectorsToRemove,
			in.RemoteChainSelectorsToAdd,
			in.RemotePoolAddresses,
			in.RemoteTokenAddresses,
		)
	default:
		return types.Transaction{}, fmt.Errorf("unsupported token pool type for ApplyChainUpdatesOp: %v", in.TokenPoolType.String())
	}

	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode ApplyChainUpdates for chains: %w", err)
	}

	return utils.GenerateMCMSTx(in.TokenPoolAddress, moduleInfo, function, args)
}

type SetChainRLConfigsInput struct {
	TokenPoolAddress     aptos.AccountAddress
	TokenPoolType        cldf.ContractType
	RemoteChainSelectors []uint64
	OutboundIsEnableds   []bool
	OutboundCapacities   []uint64
	OutboundRates        []uint64
	InboundIsEnableds    []bool
	InboundCapacities    []uint64
	InboundRates         []uint64
}

var SetChainRateLimiterConfigsOp = operations.NewOperation(
	"set-chain-rate-limiter-configs-op",
	Version1_0_0,
	"Set chain rate limiter configs for an Aptos token pool",
	setChainRateLimiterConfigs,
)

func setChainRateLimiterConfigs(b operations.Bundle, deps AptosDeps, in SetChainRLConfigsInput) (types.Transaction, error) {
	var (
		moduleInfo bind.ModuleInformation
		function   string
		args       [][]byte
		err        error
	)

	switch in.TokenPoolType {
	case shared.AptosManagedTokenPoolType:
		poolBind := managed_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		moduleInfo, function, _, args, err = poolBind.ManagedTokenPool().Encoder().SetChainRateLimiterConfigs(
			in.RemoteChainSelectors,
			in.OutboundIsEnableds,
			in.OutboundCapacities,
			in.OutboundRates,
			in.InboundIsEnableds,
			in.InboundCapacities,
			in.InboundRates,
		)
	case shared.AptosRegulatedTokenPoolType:
		poolBind := regulated_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		moduleInfo, function, _, args, err = poolBind.RegulatedTokenPool().Encoder().SetChainRateLimiterConfigs(
			in.RemoteChainSelectors,
			in.OutboundIsEnableds,
			in.OutboundCapacities,
			in.OutboundRates,
			in.InboundIsEnableds,
			in.InboundCapacities,
			in.InboundRates,
		)
	case shared.BurnMintTokenPool:
		poolBind := burn_mint_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		moduleInfo, function, _, args, err = poolBind.BurnMintTokenPool().Encoder().SetChainRateLimiterConfigs(
			in.RemoteChainSelectors,
			in.OutboundIsEnableds,
			in.OutboundCapacities,
			in.OutboundRates,
			in.InboundIsEnableds,
			in.InboundCapacities,
			in.InboundRates,
		)
	case shared.LockReleaseTokenPool:
		poolBind := lock_release_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		moduleInfo, function, _, args, err = poolBind.LockReleaseTokenPool().Encoder().SetChainRateLimiterConfigs(
			in.RemoteChainSelectors,
			in.OutboundIsEnableds,
			in.OutboundCapacities,
			in.OutboundRates,
			in.InboundIsEnableds,
			in.InboundCapacities,
			in.InboundRates,
		)
	default:
		return types.Transaction{}, fmt.Errorf("unsupported token pool type for SetChainRateLimiterConfigsOp: %v", in.TokenPoolType.String())
	}
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode SetChainRateLimiterConfigs for chains: %w", err)
	}

	return utils.GenerateMCMSTx(in.TokenPoolAddress, moduleInfo, function, args)
}

type AddRemotePoolsInput struct {
	TokenPoolAddress     aptos.AccountAddress
	TokenPoolType        cldf.ContractType
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

	switch in.TokenPoolType {
	case shared.AptosManagedTokenPoolType:
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
	case shared.AptosRegulatedTokenPoolType:
		poolBind := regulated_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		for i, selector := range in.RemoteChainSelectors {
			moduleInfo, function, _, args, err := poolBind.RegulatedTokenPool().Encoder().AddRemotePool(selector, in.RemotePoolAddresses[i])
			if err != nil {
				return nil, fmt.Errorf("failed to encode AddRemotePools for remote selector %d: %w", selector, err)
			}
			tx, err := utils.GenerateMCMSTx(in.TokenPoolAddress, moduleInfo, function, args)
			if err != nil {
				return nil, err
			}
			txs[i] = tx
		}
	case shared.BurnFromMintTokenPool:
		poolBind := burn_mint_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		for i, selector := range in.RemoteChainSelectors {
			moduleInfo, function, _, args, err := poolBind.BurnMintTokenPool().Encoder().AddRemotePool(selector, in.RemotePoolAddresses[i])
			if err != nil {
				return nil, fmt.Errorf("failed to encode AddRemotePools for remote selector %d: %w", selector, err)
			}
			tx, err := utils.GenerateMCMSTx(in.TokenPoolAddress, moduleInfo, function, args)
			if err != nil {
				return nil, err
			}
			txs[i] = tx
		}
	case shared.LockReleaseTokenPool:
		poolBind := lock_release_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		for i, selector := range in.RemoteChainSelectors {
			moduleInfo, function, _, args, err := poolBind.LockReleaseTokenPool().Encoder().AddRemotePool(selector, in.RemotePoolAddresses[i])
			if err != nil {
				return nil, fmt.Errorf("failed to encode AddRemotePools for remote selector %d: %w", selector, err)
			}
			tx, err := utils.GenerateMCMSTx(in.TokenPoolAddress, moduleInfo, function, args)
			if err != nil {
				return nil, err
			}
			txs[i] = tx
		}
	default:
		return nil, fmt.Errorf("unsupported token pool type for AddRemotePoolsOp: %v", in.TokenPoolType.String())
	}

	return txs, nil
}

// ########################
// # Token Pool Ownership #
// ########################

type TransferTokenPoolOwnershipInput struct {
	TokenPoolAddress aptos.AccountAddress
	To               aptos.AccountAddress
	TokenPoolType    cldf.ContractType
}

var TransferTokenPoolOwnershipOp = operations.NewOperation(
	"transfer-token-pool-ownerhip-op",
	Version1_0_0,
	"Initiated the ownership transfer of a managed/BnM/LnR token pool to a given address",
	transferTokenPoolOwnership,
)

func transferTokenPoolOwnership(b operations.Bundle, deps AptosDeps, in TransferTokenPoolOwnershipInput) (types.Transaction, error) {
	var (
		moduleInfo bind.ModuleInformation
		function   string
		args       [][]byte
		err        error
	)
	switch in.TokenPoolType {
	case shared.AptosManagedTokenPoolType:
		tokenPoolContract := managed_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		moduleInfo, function, _, args, err = tokenPoolContract.ManagedTokenPool().Encoder().TransferOwnership(in.To)
	case shared.BurnMintTokenPool:
		tokenPoolContract := burn_mint_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		moduleInfo, function, _, args, err = tokenPoolContract.BurnMintTokenPool().Encoder().TransferOwnership(in.To)
	case shared.LockReleaseTokenPool:
		tokenPoolContract := lock_release_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		moduleInfo, function, _, args, err = tokenPoolContract.LockReleaseTokenPool().Encoder().TransferOwnership(in.To)
	default:
		return types.Transaction{}, fmt.Errorf("unsupported token pool type for TransferTokenPoolOwnershipOp: %s", in.TokenPoolType.String())
	}
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode TransferOwnership: %w", err)
	}

	return utils.GenerateMCMSTx(in.TokenPoolAddress, moduleInfo, function, args)
}

type AcceptTokenPoolOwnershipInput struct {
	TokenPoolAddress aptos.AccountAddress
	TokenPoolType    cldf.ContractType
}

var AcceptTokenPoolOwnershipOp = operations.NewOperation(
	"accept-token-pool-ownership-op",
	Version1_0_0,
	"Accepts ownership of a managed/BnM/LnR token pool",
	acceptTokenPoolOwnership,
)

func acceptTokenPoolOwnership(b operations.Bundle, deps AptosDeps, in AcceptTokenPoolOwnershipInput) (types.Transaction, error) {
	var (
		moduleInfo bind.ModuleInformation
		function   string
		args       [][]byte
		err        error
	)
	switch in.TokenPoolType {
	case shared.AptosManagedTokenPoolType:
		tokenPoolContract := managed_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		moduleInfo, function, _, args, err = tokenPoolContract.ManagedTokenPool().Encoder().AcceptOwnership()
	case shared.BurnMintTokenPool:
		tokenPoolContract := burn_mint_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		moduleInfo, function, _, args, err = tokenPoolContract.BurnMintTokenPool().Encoder().AcceptOwnership()
	case shared.LockReleaseTokenPool:
		tokenPoolContract := lock_release_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		moduleInfo, function, _, args, err = tokenPoolContract.LockReleaseTokenPool().Encoder().AcceptOwnership()
	default:
		return types.Transaction{}, fmt.Errorf("unsupported token pool type for AcceptTokenPoolOwnershipOp: %s", in.TokenPoolType.String())
	}
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode AcceptOwnership: %w", err)
	}

	return utils.GenerateMCMSTx(in.TokenPoolAddress, moduleInfo, function, args)
}

type ExecuteTokenPoolOwnershipTransferInput struct {
	TokenPoolAddress aptos.AccountAddress
	To               aptos.AccountAddress
	TokenPoolType    cldf.ContractType
}

var ExecuteTokenPoolOwnershipTransferOp = operations.NewOperation(
	"execute-token-pool-ownership-transfer-op",
	Version1_0_0,
	"Executes the ownership transfer of a managed/BnM/LnR token pool",
	executeTokenPoolOwnershipTransfer,
)

func executeTokenPoolOwnershipTransfer(b operations.Bundle, deps AptosDeps, in ExecuteTokenPoolOwnershipTransferInput) (types.Transaction, error) {
	var (
		moduleInfo bind.ModuleInformation
		function   string
		args       [][]byte
		err        error
	)
	switch in.TokenPoolType {
	case shared.AptosManagedTokenPoolType:
		tokenPoolContract := managed_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		moduleInfo, function, _, args, err = tokenPoolContract.ManagedTokenPool().Encoder().ExecuteOwnershipTransfer(in.To)
	case shared.BurnMintTokenPool:
		tokenPoolContract := burn_mint_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		moduleInfo, function, _, args, err = tokenPoolContract.BurnMintTokenPool().Encoder().ExecuteOwnershipTransfer(in.To)
	case shared.LockReleaseTokenPool:
		tokenPoolContract := lock_release_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		moduleInfo, function, _, args, err = tokenPoolContract.LockReleaseTokenPool().Encoder().ExecuteOwnershipTransfer(in.To)
	default:
		return types.Transaction{}, fmt.Errorf("unsupported token pool type for ExecuteTokenPoolOwnershipTransferInput: %s", in.TokenPoolType.String())
	}
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode ExecuteOwnershipTransfer: %w", err)
	}

	return utils.GenerateMCMSTx(in.TokenPoolAddress, moduleInfo, function, args)
}
