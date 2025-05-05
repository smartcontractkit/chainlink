package sequence

import (
	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

// Deploy CCIP Sequence
type DeployCCIPSeqInput struct {
	MCMSAddress aptos.AccountAddress
	CCIPConfig  config.ChainContractParams
}

type DeployCCIPSeqOutput struct {
	CCIPAddress    aptos.AccountAddress
	MCMSOperations []mcmstypes.BatchOperation
}

var DeployCCIPSequence = operations.NewSequence(
	"deploy-aptos-ccip-sequence",
	operation.Version1_0_0,
	"Deploy Aptos CCIP contracts and initialize them",
	deployCCIPSequence,
)

func deployCCIPSequence(b operations.Bundle, deps operation.AptosDeps, in DeployCCIPSeqInput) (DeployCCIPSeqOutput, error) {
	var mcmsOperations []mcmstypes.BatchOperation

	// Cleanup staging area
	cleanupInput := operation.CleanupStagingAreaInput{
		MCMSAddress: in.MCMSAddress,
	}
	cleanupReport, err := operations.ExecuteOperation(b, operation.CleanupStagingAreaOp, deps, cleanupInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	if len(cleanupReport.Output.Transactions) > 0 {
		mcmsOperations = append(mcmsOperations, cleanupReport.Output)
	}

	// Generate proposal to deploy CCIP package
	deployCCIPInput := operation.DeployCCIPInput{
		MCMSAddress: in.MCMSAddress,
	}
	deployCCIPReport, err := operations.ExecuteOperation(b, operation.GenerateDeployCCIPProposalOp, deps, deployCCIPInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	ccipAddress := deployCCIPReport.Output.CCIPAddress
	// For CCIP deployment the txs cannot be batched - it'd exceed Aptos API limits
	// so it's converted to batch operations with single transactions in each
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployCCIPReport.Output.MCMSOperations)...)

	deployModulesInput := operation.DeployModulesInput{
		MCMSAddress: in.MCMSAddress,
		CCIPAddress: ccipAddress,
	}
	// Generate proposal to deploy OnRamp module
	deployOnRampReport, err := operations.ExecuteOperation(b, operation.GenerateDeployOnRampProposalOp, deps, deployModulesInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployOnRampReport.Output)...)

	// Generate proposal to deploy OffRamp module
	deployOffRampReport, err := operations.ExecuteOperation(b, operation.GenerateDeployOffRampProposalOp, deps, deployModulesInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployOffRampReport.Output)...)

	// Generate proposal to deploy Router module
	deployRouterReport, err := operations.ExecuteOperation(b, operation.GenerateDeployRouterProposalOp, deps, deployModulesInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployRouterReport.Output)...)

	// Generate proposal to initialize CCIP
	initCCIPInput := operation.InitializeCCIPInput{
		MCMSAddress: in.MCMSAddress,
		CCIPAddress: ccipAddress,
		CCIPConfig:  in.CCIPConfig,
	}
	initCCIPReport, err := operations.ExecuteOperation(b, operation.InitializeCCIPOp, deps, initCCIPInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, initCCIPReport.Output)

	return DeployCCIPSeqOutput{
		CCIPAddress:    ccipAddress,
		MCMSOperations: mcmsOperations,
	}, nil
}
