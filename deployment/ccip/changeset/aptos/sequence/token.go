package sequence

import (
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/mcms/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

type DeployTokenSeqOutput struct {
	TokenAddress    aptos.AccountAddress
	TokenObjAddress aptos.AccountAddress
	MCMSOperations  []mcmstypes.BatchOperation
}

var DeployAptosTokenSequence = operations.NewSequence(
	"deploy-aptos-token",
	operation.Version1_0_0,
	"Deploys token and configures",
	deployAptosTokenSequence,
)

func deployAptosTokenSequence(b operations.Bundle, deps operation.AptosDeps, in config.TokenParams) (DeployTokenSeqOutput, error) {
	mcmsOperations := []mcmstypes.BatchOperation{}

	// Cleanup staging area
	cleanupReport, err := operations.ExecuteOperation(b, operation.CleanupStagingAreaOp, deps, deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].MCMSAddress)
	if err != nil {
		return DeployTokenSeqOutput{}, err
	}
	if len(cleanupReport.Output.Transactions) > 0 {
		mcmsOperations = append(mcmsOperations, cleanupReport.Output)
	}

	// Deploy token
	deployTInput := operation.DeployTokenInput{
		Name:   in.Name,
		Symbol: in.Symbol,
	}
	deployTReport, err := operations.ExecuteOperation(b, operation.DeployTokenOp, deps, deployTInput)
	if err != nil {
		return DeployTokenSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployTReport.Output.MCMSOps)...)

	// Deploy token MCMS Registrar
	deployRegReport, err := operations.ExecuteOperation(b, operation.DeployTokenMCMSRegistrarOp, deps, deployTReport.Output.TokenObjAddress)
	if err != nil {
		return DeployTokenSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployRegReport.Output)...)

	// Initialize token
	initTokenInput := operation.InitializeTokenInput{
		TokenObjAddress: deployTReport.Output.TokenObjAddress,
		MaxSupply:       in.MaxSupply,
		Name:            in.Name,
		Symbol:          in.Symbol,
		Decimals:        in.Decimals,
		Icon:            in.Icon,
		Project:         in.Project,
	}
	initTokenReport, err := operations.ExecuteOperation(b, operation.InitializeTokenOp, deps, initTokenInput)
	if err != nil {
		return DeployTokenSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  []types.Transaction{initTokenReport.Output},
	})

	return DeployTokenSeqOutput{
		TokenAddress:    deployTReport.Output.TokenAddress,
		TokenObjAddress: deployTReport.Output.TokenObjAddress,
		MCMSOperations:  mcmsOperations,
	}, nil
}
