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
	TokenParams config.TokenParams
	MCMSAddress aptos.AccountAddress
}

type DeployTokenSeqOutput struct {
	TokenAddress      aptos.AccountAddress
	TokenObjAddress   aptos.AccountAddress
	TokenOwnerAddress aptos.AccountAddress
	MCMSOperations    []mcmstypes.BatchOperation
}

var DeployAptosTokenSequence = operations.NewSequence(
	"deploy-aptos-token",
	operation.Version1_0_0,
	"Deploys token and configures",
	deployAptosTokenSequence,
)

func deployAptosTokenSequence(b operations.Bundle, deps operation.AptosDeps, in DeployTokenSeqInput) (DeployTokenSeqOutput, error) {
	mcmsOperations := []mcmstypes.BatchOperation{}

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
		TokenObjAddress: deployTReport.Output.TokenObjAddress,
		MCMSAddress:     in.MCMSAddress,
	}
	deployRegReport, err := operations.ExecuteOperation(b, operation.DeployTokenMCMSRegistrarOp, deps, deployTokenRegistrarIn)
	if err != nil {
		return DeployTokenSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployRegReport.Output)...)

	// Initialize token
	initTokenInput := operation.InitializeTokenInput{
		TokenObjAddress: deployTReport.Output.TokenObjAddress,
		MaxSupply:       in.TokenParams.MaxSupply,
		Name:            in.TokenParams.Name,
		Symbol:          string(in.TokenParams.Symbol),
		Decimals:        in.TokenParams.Decimals,
		Icon:            in.TokenParams.Icon,
		Project:         in.TokenParams.Project,
	}
	initTokenReport, err := operations.ExecuteOperation(b, operation.InitializeTokenOp, deps, initTokenInput)
	if err != nil {
		return DeployTokenSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  []mcmstypes.Transaction{initTokenReport.Output},
	})

	return DeployTokenSeqOutput{
		TokenAddress:      deployTReport.Output.TokenAddress,
		TokenObjAddress:   deployTReport.Output.TokenObjAddress,
		TokenOwnerAddress: deployTReport.Output.TokenOwnerAddress,
		MCMSOperations:    mcmsOperations,
	}, nil
}
