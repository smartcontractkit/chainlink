package sequence

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
)

var UpdateCCIPSequence = operations.NewSequence(
	"update-aptos-ccip-sequence",
	operation.Version1_0_0,
	"Update Aptos CCIP contracts",
	updateCCIPSequence,
)

func updateCCIPSequence(b operations.Bundle, deps operation.AptosDeps, in config.UpdateAptosChainConfig) ([]mcmstypes.BatchOperation, error) {
	var mcmsOperations []mcmstypes.BatchOperation

	// Cleanup MCMS staging area if not clear
	cleanupReport, err := operations.ExecuteOperation(b, operation.CleanupStagingAreaOp, deps, deps.OnChainState.MCMSAddress)
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
		deployCCIPReport, err := operations.ExecuteOperation(b, operation.DeployCCIPOp, deps, deployCCIPInput)
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
		deployOnRampReport, err := operations.ExecuteOperation(b, operation.DeployOnRampOp, deps, deployModulesInput)
		if err != nil {
			return nil, err
		}
		mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployOnRampReport.Output)...)
	}

	if in.UpdateOffRamp {
		deployOffRampReport, err := operations.ExecuteOperation(b, operation.DeployOffRampOp, deps, deployModulesInput)
		if err != nil {
			return nil, err
		}
		mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployOffRampReport.Output)...)
	}

	if in.UpdateRouter {
		deployRouterReport, err := operations.ExecuteOperation(b, operation.DeployRouterOp, deps, deployModulesInput)
		if err != nil {
			return nil, err
		}
		mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployRouterReport.Output)...)
	}

	return mcmsOperations, nil
}
