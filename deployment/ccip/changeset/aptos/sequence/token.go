package sequence

import (
	"github.com/aptos-labs/aptos-go-sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
)

type DeployTokenSeqInput struct {
	MCMSAddress aptos.AccountAddress
	TokenParams config.TokenParams
	TokenMint   *config.TokenMint
}

type DeployTokenSeqOutput struct {
	TokenAddress        aptos.AccountAddress
	TokenCodeObjAddress aptos.AccountAddress
	TokenOwnerAddress   aptos.AccountAddress
	MCMSOperations      []mcmstypes.BatchOperation
}

var DeployAptosTokenSequence = operations.NewSequence(
	"deploy-aptos-token",
	operation.Version1_0_0,
	"Deploys token and configures",
	deployAptosTokenSequence,
)

func deployAptosTokenSequence(b operations.Bundle, deps operation.AptosDeps, in DeployTokenSeqInput) (DeployTokenSeqOutput, error) {
	var mcmsOperations []mcmstypes.BatchOperation
	var txs []mcmstypes.Transaction

	// Cleanup staging area
	cleanupReport, err := operations.ExecuteOperation(b, operation.CleanupStagingAreaOp, deps, in.MCMSAddress)
	if err != nil {
		return DeployTokenSeqOutput{}, err
	}
	if len(cleanupReport.Output.Transactions) > 0 {
		mcmsOperations = append(mcmsOperations, cleanupReport.Output)
	}

	// Deploy token
	deployTInput := operation.DeployTokenInput{
		Name:        in.TokenParams.Name,
		Symbol:      string(in.TokenParams.Symbol),
		MCMSAddress: in.MCMSAddress,
	}
	deployTReport, err := operations.ExecuteOperation(b, operation.DeployTokenOp, deps, deployTInput)
	if err != nil {
		return DeployTokenSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployTReport.Output.MCMSOps)...)

	// Deploy token MCMS Registrar
	deployTokenRegistrarIn := operation.DeployTokenRegistrarInput{
		TokenCodeObjectAddress: deployTReport.Output.TokenCodeObjectAddress,
		MCMSAddress:            in.MCMSAddress,
	}
	deployRegReport, err := operations.ExecuteOperation(b, operation.DeployTokenMCMSRegistrarOp, deps, deployTokenRegistrarIn)
	if err != nil {
		return DeployTokenSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployRegReport.Output)...)

	// Initialize token
	initTokenInput := operation.InitializeTokenInput{
		TokenCodeObjectAddress: deployTReport.Output.TokenCodeObjectAddress,
		MaxSupply:              in.TokenParams.MaxSupply,
		Name:                   in.TokenParams.Name,
		Symbol:                 string(in.TokenParams.Symbol),
		Decimals:               in.TokenParams.Decimals,
		Icon:                   in.TokenParams.Icon,
		Project:                in.TokenParams.Project,
	}
	initTokenReport, err := operations.ExecuteOperation(b, operation.InitializeTokenOp, deps, initTokenInput)
	if err != nil {
		return DeployTokenSeqOutput{}, err
	}
	txs = append(txs, initTokenReport.Output)

	// Mint test tokens
	if in.TokenMint != nil {
		mintTokenInput := operation.MintTokensInput{
			TokenCodeObjectAddress: deployTReport.Output.TokenCodeObjectAddress,
			To:                     in.TokenMint.To,
			Amount:                 in.TokenMint.Amount,
		}
		mintTokenReport, err := operations.ExecuteOperation(b, operation.MintTokensOp, deps, mintTokenInput)
		if err != nil {
			return DeployTokenSeqOutput{}, err
		}
		txs = append(txs, mintTokenReport.Output)
	}

	mcmsOperations = append(mcmsOperations, mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  txs,
	})

	return DeployTokenSeqOutput{
		TokenAddress:        deployTReport.Output.TokenAddress,
		TokenCodeObjAddress: deployTReport.Output.TokenCodeObjectAddress,
		TokenOwnerAddress:   deployTReport.Output.TokenOwnerAddress,
		MCMSOperations:      mcmsOperations,
	}, nil
}

type DeployTokenFaucetSeqInput struct {
	MCMSAddress         aptos.AccountAddress
	TokenCodeObjAddress aptos.AccountAddress
}

var DeployTokenFaucetSequence = operations.NewSequence(
	"deploy-aptos-token-faucet",
	operation.Version1_0_0,
	"Deploys a token faucet onto an existing manage_token instance",
	deployAptosTokenFaucetSequence,
)

func deployAptosTokenFaucetSequence(b operations.Bundle, deps operation.AptosDeps, in DeployTokenFaucetSeqInput) ([]mcmstypes.BatchOperation, error) {
	var mcmsOperations []mcmstypes.BatchOperation

	// Cleanup staging area
	cleanupReport, err := operations.ExecuteOperation(b, operation.CleanupStagingAreaOp, deps, in.MCMSAddress)
	if err != nil {
		return nil, err
	}
	if len(cleanupReport.Output.Transactions) > 0 {
		mcmsOperations = append(mcmsOperations, cleanupReport.Output)
	}

	// Deploy token faucet module
	deployTokenFaucetInput := operation.DeployTokenFaucetInput{
		MCMSAddress:            in.MCMSAddress,
		TokenCodeObjectAddress: in.TokenCodeObjAddress,
	}
	deployTokenFaucetReport, err := operations.ExecuteOperation(b, operation.DeployTokenFaucetOp, deps, deployTokenFaucetInput)
	if err != nil {
		return nil, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployTokenFaucetReport.Output)...)

	// Grant Mint rights to ManagedTokenFaucet signer
	managedTokenFaucetStateAddress := in.TokenCodeObjAddress.NamedObjectAddress([]byte("ManagedTokenFaucet"))
	applyMintersInput := operation.ApplyAllowedMintersInput{
		TokenCodeObjectAddress: in.TokenCodeObjAddress,
		MintersToAdd:           []aptos.AccountAddress{managedTokenFaucetStateAddress},
	}
	applyAllowedMintersReport, err := operations.ExecuteOperation(b, operation.ApplyAllowedMintersOp, deps, applyMintersInput)
	if err != nil {
		return nil, err
	}
	mcmsOperations = append(mcmsOperations, mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  []mcmstypes.Transaction{applyAllowedMintersReport.Output},
	})

	return mcmsOperations, nil
}
