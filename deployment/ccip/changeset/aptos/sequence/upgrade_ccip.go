package sequence

import (
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/dependency"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
)

var UpgradeCCIPSequence = operations.NewSequence(
	"upgrade-aptos-ccip-sequence",
	operation.Version1_0_0,
	"Upgrade Aptos CCIP contracts",
	upgradeCCIPSequence,
)

func upgradeCCIPSequence(b operations.Bundle, deps dependency.AptosDeps, in config.UpgradeAptosChainConfig) ([]mcmstypes.BatchOperation, error) {
	var mcmsOperations []mcmstypes.BatchOperation
	mcmsAddress := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].MCMSAddress
	// Cleanup staging area
	cleanupReport, err := operations.ExecuteOperation(b, operation.CleanupStagingAreaOp, deps, mcmsAddress)
	if err != nil {
		return nil, err
	}
	if len(cleanupReport.Output.Transactions) > 0 {
		mcmsOperations = append(mcmsOperations, cleanupReport.Output)
	}

	if in.UpgradeCCIP {
		deployCCIPInput := operation.DeployCCIPInput{
			MCMSAddress: mcmsAddress,
			IsUpgrade:   in.UpgradeCCIP,
		}
		deployCCIPReport, err := operations.ExecuteOperation(b, operation.DeployCCIPOp, deps, deployCCIPInput)
		if err != nil {
			return nil, err
		}
		mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployCCIPReport.Output.MCMSOperations)...)
	}

	deployModulesInput := operation.DeployModulesInput{
		MCMSAddress: mcmsAddress,
		CCIPAddress: deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].CCIPAddress,
	}

	if in.UpgradeOnRamp {
		deployOnRampReport, err := operations.ExecuteOperation(b, operation.DeployOnRampOp, deps, deployModulesInput)
		if err != nil {
			return nil, err
		}
		mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployOnRampReport.Output)...)
	}

	if in.UpgradeOffRamp {
		deployOffRampReport, err := operations.ExecuteOperation(b, operation.DeployOffRampOp, deps, deployModulesInput)
		if err != nil {
			return nil, err
		}
		mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployOffRampReport.Output)...)
	}

	if in.UpgradeRouter {
		deployRouterReport, err := operations.ExecuteOperation(b, operation.DeployRouterOp, deps, deployModulesInput)
		if err != nil {
			return nil, err
		}
		mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployRouterReport.Output)...)
	}

	return mcmsOperations, nil
}
