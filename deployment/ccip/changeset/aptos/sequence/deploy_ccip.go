package sequence

import (
	"github.com/aptos-labs/aptos-go-sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
)

type DeployCCIPSeqInput struct {
	MCMSAddress      aptos.AccountAddress
	CCIPConfig       config.ChainContractParams
	LinkTokenAddress aptos.AccountAddress
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

	// Cleanup MCMS staging area if not clear
	cleanupReport, err := operations.ExecuteOperation(b, operation.CleanupStagingAreaOp, deps, in.MCMSAddress)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	if len(cleanupReport.Output.Transactions) > 0 {
		mcmsOperations = append(mcmsOperations, cleanupReport.Output)
	}

	// Generate batch operations to deploy CCIP package
	deployCCIPInput := operation.DeployCCIPInput{
		MCMSAddress: in.MCMSAddress,
		IsUpdate:    false,
	}
	deployCCIPReport, err := operations.ExecuteOperation(b, operation.DeployCCIPOp, deps, deployCCIPInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	ccipAddress := deployCCIPReport.Output.CCIPAddress
	// For CCIP deployment the txs cannot be batched - it'd exceed Aptos API limits
	// so they're converted to batch operations with single transactions in each batch
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployCCIPReport.Output.MCMSOperations)...)

	// Generate batch operations to deploy CCIP modules
	deployModulesInput := operation.DeployModulesInput{
		MCMSAddress: in.MCMSAddress,
		CCIPAddress: ccipAddress,
	}
	// OnRamp module
	deployOnRampReport, err := operations.ExecuteOperation(b, operation.DeployOnRampOp, deps, deployModulesInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployOnRampReport.Output)...)
	// OffRamp module
	deployOffRampReport, err := operations.ExecuteOperation(b, operation.DeployOffRampOp, deps, deployModulesInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployOffRampReport.Output)...)
	// Router module
	deployRouterReport, err := operations.ExecuteOperation(b, operation.DeployRouterOp, deps, deployModulesInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployRouterReport.Output)...)

	// Generate batch operations to initialize CCIP
	initCCIPInput := operation.InitializeCCIPInput{
		MCMSAddress:      in.MCMSAddress,
		CCIPAddress:      ccipAddress,
		CCIPConfig:       in.CCIPConfig,
		LinkTokenAddress: in.LinkTokenAddress,
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
