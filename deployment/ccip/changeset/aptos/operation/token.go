package operation

import (
	"fmt"
	"math/big"

	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"

	managed_token "github.com/smartcontractkit/chainlink-aptos/bindings/managed_token"
	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
)

const managedTokenStateSeed = "managed_token::managed_token::token_state"

type DeployTokenInput struct {
	Name        string
	Symbol      string
	MCMSAddress aptos.AccountAddress
}

type DeployTokenOutput struct {
	TokenObjAddress   aptos.AccountAddress
	TokenAddress      aptos.AccountAddress
	TokenOwnerAddress aptos.AccountAddress
	MCMSOps           []types.Operation
}

// DeployTokenOp generates proposal to deploy a token
var DeployTokenOp = operations.NewOperation(
	"deploy-token-op",
	Version1_0_0,
	"deploy token",
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
		return DeployTokenOutput{}, fmt.Errorf("failed to compile managed token: %w", err)
	}
	ops, err := utils.CreateChunksAndStage(managedTokenPayload, mcmsContract, deps.AptosChain.Selector, managedTokenSeed, nil)
	if err != nil {
		return DeployTokenOutput{}, fmt.Errorf("failed to create chunks for token pool: %w", err)
	}

	return DeployTokenOutput{
		TokenObjAddress:   managedTokenObjectAddress,
		TokenAddress:      managedTokenMetadataAddress,
		TokenOwnerAddress: managedTokenOwnerAddress,
		MCMSOps:           ops,
	}, nil
}

type DeployTokenRegistrarInput struct {
	TokenObjAddress aptos.AccountAddress
	MCMSAddress     aptos.AccountAddress
}

// DeployTokenMCMSRegistrarOp generates proposal to deploy a MCMS registrar on a token package
var DeployTokenMCMSRegistrarOp = operations.NewOperation(
	"deploy-token-mcms-registrar-op",
	Version1_0_0,
	"deploy token MCMS registrar on token package",
	deployTokenMCMSRegistrar,
)

func deployTokenMCMSRegistrar(b operations.Bundle, deps AptosDeps, in DeployTokenRegistrarInput) ([]types.Operation, error) {
	mcmsContract := mcmsbind.Bind(in.MCMSAddress, deps.AptosChain.Client)

	// Deploy MCMS Registrar
	mcmsRegistrarPayload, err := managed_token.CompileMCMSRegistrar(in.TokenObjAddress, in.MCMSAddress, true)
	if err != nil {
		return nil, fmt.Errorf("failed to compile MCMS registrar: %w", err)
	}
	ops, err := utils.CreateChunksAndStage(mcmsRegistrarPayload, mcmsContract, deps.AptosChain.Selector, "", &in.TokenObjAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create chunks for token pool: %w", err)
	}

	return ops, nil
}

type InitializeTokenInput struct {
	TokenObjAddress aptos.AccountAddress
	MaxSupply       *big.Int
	Name            string
	Symbol          string
	Decimals        byte
	Icon            string
	Project         string
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
	boundManagedToken := managed_token.Bind(in.TokenObjAddress, deps.AptosChain.Client)
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
	tx, err := utils.GenerateMCMSTx(in.TokenObjAddress, moduleInfo, function, args)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	return tx, nil
}
