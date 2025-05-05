package sequence

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

var UpdateCCIPSequence = operations.NewSequence(
	"update-aptos-ccip-sequence",
	operation.Version1_0_0,
	"Update Aptos CCIP contracts",
	updateCCIPSequence,
)

func updateCCIPSequence(b operations.Bundle, deps operation.AptosDeps, in config.UpdateAptosChainConfig) ([]mcmstypes.BatchOperation, error) {
	var mcmsOperations []mcmstypes.BatchOperation

	// Cleanup staging area
	cleanupInput := operation.CleanupStagingAreaInput{
		MCMSAddress: deps.OnChainState.MCMSAddress,
	}
	cleanupReport, err := operations.ExecuteOperation(b, operation.CleanupStagingAreaOp, deps, cleanupInput)
	if err != nil {
		return nil, err
	}
	if len(cleanupReport.Output.Transactions) > 0 {
		mcmsOperations = append(mcmsOperations, cleanupReport.Output)
	}

	if in.UpdateCCIP {
		deployCCIPInput := operation.DeployCCIPInput{
			MCMSAddress: deps.OnChainState.MCMSAddress,
			IsUpdate:    in.UpdateCCIP,
		}
		deployCCIPReport, err := operations.ExecuteOperation(b, operation.GenerateDeployCCIPProposalOp, deps, deployCCIPInput)
		if err != nil {
			return nil, err
		}
		mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployCCIPReport.Output.MCMSOperations)...)
	}

	deployModulesInput := operation.DeployModulesInput{
		MCMSAddress: deps.OnChainState.MCMSAddress,
		CCIPAddress: deps.OnChainState.CCIPAddress,
	}

	if in.UpdateOnRamp {
		deployOnRampReport, err := operations.ExecuteOperation(b, operation.GenerateDeployOnRampProposalOp, deps, deployModulesInput)
		if err != nil {
			return nil, err
		}
		mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployOnRampReport.Output)...)
	}

	if in.UpdateOffRamp {
		deployOffRampReport, err := operations.ExecuteOperation(b, operation.GenerateDeployOffRampProposalOp, deps, deployModulesInput)
		if err != nil {
			return nil, err
		}
		mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployOffRampReport.Output)...)
	}

	if in.UpdateRouter {
		deployRouterReport, err := operations.ExecuteOperation(b, operation.GenerateDeployRouterProposalOp, deps, deployModulesInput)
		if err != nil {
			return nil, err
		}
		mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployRouterReport.Output)...)
	}

	return mcmsOperations, nil
}
