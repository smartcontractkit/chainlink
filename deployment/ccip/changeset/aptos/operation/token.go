package operation

import (
	"fmt"
	"math/big"

	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/managed_token_faucet"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"

	"github.com/smartcontractkit/chainlink-aptos/bindings/managed_token"
	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
)

const managedTokenStateSeed = "managed_token::managed_token::token_state"

type DeployTokenInput struct {
	Name        string
	Symbol      string
	MCMSAddress aptos.AccountAddress
}

type DeployTokenOutput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	TokenAddress           aptos.AccountAddress
	TokenOwnerAddress      aptos.AccountAddress
	MCMSOps                []types.Operation
}

// DeployTokenOp generates proposal to deploy a token
var DeployTokenOp = operations.NewOperation(
	"deploy-token-op",
	Version1_0_0,
	"Deploy a managed token instance",
	deployToken,
)

func deployToken(b operations.Bundle, deps AptosDeps, in DeployTokenInput) (DeployTokenOutput, error) {
	mcmsContract := mcmsbind.Bind(in.MCMSAddress, deps.AptosChain.Client)

	// Calculate token address
	managedTokenSeed := fmt.Sprintf("%s::%s", in.Name, in.Symbol) // Use name and symbol as seed for uniqueness
	managedTokenObjectAddress, err := mcmsContract.MCMSRegistry().GetNewCodeObjectAddress(nil, []byte(managedTokenSeed))
	if err != nil {
		return DeployTokenOutput{}, fmt.Errorf("failed to GetNewCodeObjectAddress: %w", err)
	}
	managedTokenOwnerAddress, err := mcmsContract.MCMSRegistry().GetNewCodeObjectOwnerAddress(nil, []byte(managedTokenSeed))
	if err != nil {
		return DeployTokenOutput{}, fmt.Errorf("failed to GetNewCodeObjectOwnerAddress: %w", err)
	}

	// Calculate token Metadata Address
	managedTokenStateAddress := managedTokenObjectAddress.NamedObjectAddress([]byte(managedTokenStateSeed))
	managedTokenMetadataAddress := managedTokenStateAddress.NamedObjectAddress([]byte(in.Symbol))

	// Compile and create deploy operation for the token
	managedTokenPayload, err := managed_token.Compile(managedTokenObjectAddress)
	if err != nil {
		return DeployTokenOutput{}, fmt.Errorf("failed to compile managed_token package: %w", err)
	}
	ops, err := utils.CreateChunksAndStage(managedTokenPayload, mcmsContract, deps.AptosChain.Selector, managedTokenSeed, nil)
	if err != nil {
		return DeployTokenOutput{}, fmt.Errorf("failed to create chunks for manage_token deployment: %w", err)
	}

	return DeployTokenOutput{
		TokenCodeObjectAddress: managedTokenObjectAddress,
		TokenAddress:           managedTokenMetadataAddress,
		TokenOwnerAddress:      managedTokenOwnerAddress,
		MCMSOps:                ops,
	}, nil
}

type DeployTokenRegistrarInput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	MCMSAddress            aptos.AccountAddress
}

// DeployTokenMCMSRegistrarOp generates proposal to deploy a MCMS registrar on a token package
var DeployTokenMCMSRegistrarOp = operations.NewOperation(
	"deploy-token-mcms-registrar-op",
	Version1_0_0,
	"Deploy token MCMS registrar onto managed token code object",
	deployTokenMCMSRegistrar,
)

func deployTokenMCMSRegistrar(b operations.Bundle, deps AptosDeps, in DeployTokenRegistrarInput) ([]types.Operation, error) {
	mcmsContract := mcmsbind.Bind(in.MCMSAddress, deps.AptosChain.Client)

	// Deploy MCMS Registrar
	mcmsRegistrarPayload, err := managed_token.CompileMCMSRegistrar(in.TokenCodeObjectAddress, in.MCMSAddress, true)
	if err != nil {
		return nil, fmt.Errorf("failed to compile MCMS registrar: %w", err)
	}
	ops, err := utils.CreateChunksAndStage(mcmsRegistrarPayload, mcmsContract, deps.AptosChain.Selector, "", &in.TokenCodeObjectAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create chunks for token pool: %w", err)
	}

	return ops, nil
}

type InitializeTokenInput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	MaxSupply              *big.Int
	Name                   string
	Symbol                 string
	Decimals               byte
	Icon                   string
	Project                string
}

// DeployTokenMCMSRegistrarOp generates proposal to deploy a MCMS registrar on a token package
var InitializeTokenOp = operations.NewOperation(
	"initialize-token-op",
	Version1_0_0,
	"initialize token",
	initializeToken,
)

func initializeToken(b operations.Bundle, deps AptosDeps, in InitializeTokenInput) (types.Transaction, error) {
	// Initialize managed token
	var maxSupply **big.Int
	if in.MaxSupply != nil {
		maxSupply = &in.MaxSupply
	}
	boundManagedToken := managed_token.Bind(in.TokenCodeObjectAddress, deps.AptosChain.Client)
	moduleInfo, function, _, args, err := boundManagedToken.ManagedToken().Encoder().Initialize(
		maxSupply,
		in.Name,
		in.Symbol,
		in.Decimals,
		in.Icon,
		in.Project,
	)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode initialize function: %w", err)
	}

	// Create MCMS tx
	tx, err := utils.GenerateMCMSTx(in.TokenCodeObjectAddress, moduleInfo, function, args)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	return tx, nil
}

type MintTokensInput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	To                     aptos.AccountAddress
	Amount                 uint64
}

var MintTokensOp = operations.NewOperation(
	"mint-tokens-op",
	Version1_0_0,
	"Mints tokens to a target account",
	mintTokens,
)

func mintTokens(b operations.Bundle, deps AptosDeps, in MintTokensInput) (types.Transaction, error) {
	boundManagedToken := managed_token.Bind(in.TokenCodeObjectAddress, deps.AptosChain.Client)
	moduleInfo, function, _, args, err := boundManagedToken.ManagedToken().Encoder().Mint(in.To, in.Amount)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode mint function: %w", err)
	}

	// Create MCMS tx
	tx, err := utils.GenerateMCMSTx(in.TokenCodeObjectAddress, moduleInfo, function, args)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	return tx, nil
}

type ApplyAllowedMintersInput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	MintersToAdd           []aptos.AccountAddress
	MintersToRemove        []aptos.AccountAddress
}

// GrantMinterPermissionsOp operation to grant minter permissions
var ApplyAllowedMintersOp = operations.NewOperation(
	"apply-allowed-minters-op",
	Version1_0_0,
	"Applies the given minters remove/add to the managed token",
	applyAllowedMinters,
)

func applyAllowedMinters(b operations.Bundle, deps AptosDeps, in ApplyAllowedMintersInput) (types.Transaction, error) {
	tokenContract := managed_token.Bind(in.TokenCodeObjectAddress, deps.AptosChain.Client)

	moduleInfo, function, _, args, err := tokenContract.ManagedToken().Encoder().ApplyAllowedMinterUpdates(in.MintersToRemove, in.MintersToAdd)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode ApplyAllowedMinterUpdates: %w", err)
	}

	return utils.GenerateMCMSTx(in.TokenCodeObjectAddress, moduleInfo, function, args)
}

type ApplyAllowedBurnersInput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	BurnersToAdd           []aptos.AccountAddress
	BurnersToRemove        []aptos.AccountAddress
}

// GrantBurnerPermissionsOp operation to grant burner permissions
var ApplyAllowedBurnersOp = operations.NewOperation(
	"apply-allowed-burners-op",
	Version1_0_0,
	"Applies the given burners remove/add to the managed token",
	applyAllowedBurners,
)

func applyAllowedBurners(b operations.Bundle, deps AptosDeps, in ApplyAllowedBurnersInput) (types.Transaction, error) {
	tokenContract := managed_token.Bind(in.TokenCodeObjectAddress, deps.AptosChain.Client)

	moduleInfo, function, _, args, err := tokenContract.ManagedToken().Encoder().ApplyAllowedBurnerUpdates(in.BurnersToRemove, in.BurnersToAdd)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode ApplyAllowedBurnerUpdates: %w", err)
	}

	return utils.GenerateMCMSTx(in.TokenCodeObjectAddress, moduleInfo, function, args)
}

type DeployTokenFaucetInput struct {
	MCMSAddress            aptos.AccountAddress
	TokenCodeObjectAddress aptos.AccountAddress
}

var DeployTokenFaucetOp = operations.NewOperation(
	"deploy-token-faucet-op",
	Version1_0_0,
	"Deploy the faucet package onto a managed token code object",
	deployTokenFaucet,
)

func deployTokenFaucet(b operations.Bundle, deps AptosDeps, in DeployTokenFaucetInput) ([]types.Operation, error) {
	boundMcmsContract := mcmsbind.Bind(in.MCMSAddress, deps.AptosChain.Client)

	managedTokenFaucetPayload, err := managed_token_faucet.Compile(in.TokenCodeObjectAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to compile managed_token_faucet package: %w", err)
	}
	ops, err := utils.CreateChunksAndStage(managedTokenFaucetPayload, boundMcmsContract, deps.AptosChain.Selector, "", &in.TokenCodeObjectAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create chunks for managed_token_faucet deployment: %w", err)
	}

	return ops, nil
}

type TransferTokenOwnershipInput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	To                     aptos.AccountAddress
}

var TransferTokenOwnershipOp = operations.NewOperation(
	"transfer-token-ownership-op",
	Version1_0_0,
	"Initiates the ownership transfer of a managed token to a given address",
	transferTokenOwnership,
)

func transferTokenOwnership(b operations.Bundle, deps AptosDeps, in TransferTokenOwnershipInput) (types.Transaction, error) {
	tokenContract := managed_token.Bind(in.TokenCodeObjectAddress, deps.AptosChain.Client)

	moduleInfo, function, _, args, err := tokenContract.ManagedToken().Encoder().TransferOwnership(in.To)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode TransferOwnership: %w", err)
	}

	return utils.GenerateMCMSTx(in.TokenCodeObjectAddress, moduleInfo, function, args)
}

type AcceptTokenOwnershipInput struct {
	TokenCodeObjectAddress aptos.AccountAddress
}

var AcceptTokenOwnershipOp = operations.NewOperation(
	"accept-token-ownership-op",
	Version1_0_0,
	"Accepts ownership of a managed token",
	acceptTokenOwnership,
)

func acceptTokenOwnership(b operations.Bundle, deps AptosDeps, in AcceptTokenOwnershipInput) (types.Transaction, error) {
	tokenContract := managed_token.Bind(in.TokenCodeObjectAddress, nil)

	moduleInfo, function, _, args, err := tokenContract.ManagedToken().Encoder().AcceptOwnership()
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode AcceptOwnership: %w", err)
	}

	return utils.GenerateMCMSTx(in.TokenCodeObjectAddress, moduleInfo, function, args)
}

type ExecuteTokenOwnershipTransferInput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	To                     aptos.AccountAddress
}

var ExecuteTokenOwnershipTransferOp = operations.NewOperation(
	"execute-token-ownership-transfer-op",
	Version1_0_0,
	"Executes the ownership transfer of a managed token, after ownership has been accepted by the receiver",
	executeTokenOwnershipTransfer,
)

func executeTokenOwnershipTransfer(b operations.Bundle, deps AptosDeps, in ExecuteTokenOwnershipTransferInput) (types.Transaction, error) {
	tokenContract := managed_token.Bind(in.TokenCodeObjectAddress, nil)

	moduleInfo, function, _, args, err := tokenContract.ManagedToken().Encoder().ExecuteOwnershipTransfer(in.To)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode AcceptOwnership: %w", err)
	}

	return utils.GenerateMCMSTx(in.TokenCodeObjectAddress, moduleInfo, function, args)
}
