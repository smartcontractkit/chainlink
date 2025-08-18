package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
)

type DeployWorkflowRegistryOpDeps struct {
	Env *cldf.Environment
}

type DeployWorkflowRegistryOpInput struct {
	ChainSelector uint64
	Qualifier     string
}

type DeployWorkflowRegistryOpOutput struct {
	Address       string
	ChainSelector uint64
	Qualifier     string
	Type          string
	Version       string
	Labels        []string
}

// DeployWorkflowRegistryOp is an operation that deploys the V2 Workflow Registry contract.
// This atomic operation performs the single side effect of deploying and registering the contract.
var DeployWorkflowRegistryOp = operations.NewOperation(
	"deploy-workflow-registry-v2-op",
	semver.MustParse("1.0.0"),
	"Deploy WorkflowRegistry V2 Contract",
	func(b operations.Bundle, deps DeployWorkflowRegistryOpDeps, input DeployWorkflowRegistryOpInput) (DeployWorkflowRegistryOpOutput, error) {
		lggr := deps.Env.Logger

		// Get the target chain
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return DeployWorkflowRegistryOpOutput{}, fmt.Errorf("chain not found for selector %d", input.ChainSelector)
		}

		// Deploy the V2 WorkflowRegistry contract
		workflowRegistryAddr, tx, workflowRegistry, err := workflow_registry_wrapper_v2.DeployWorkflowRegistry(
			chain.DeployerKey,
			chain.Client,
		)
		if err != nil {
			return DeployWorkflowRegistryOpOutput{}, fmt.Errorf("failed to deploy WorkflowRegistry V2: %w", err)
		}

		// Wait for deployment confirmation
		_, err = chain.Confirm(tx)
		if err != nil {
			return DeployWorkflowRegistryOpOutput{}, fmt.Errorf("failed to confirm WorkflowRegistry V2 deployment: %w", err)
		}

		// Get type and version from the deployed contract
		tvStr, err := workflowRegistry.TypeAndVersion(&bind.CallOpts{})
		if err != nil {
			return DeployWorkflowRegistryOpOutput{}, fmt.Errorf("failed to get type and version: %w", err)
		}

		tv, err := cldf.TypeAndVersionFromString(tvStr)
		if err != nil {
			return DeployWorkflowRegistryOpOutput{}, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
		}

		lggr.Infof("Deployed %s on chain selector %d at address %s", tv.String(), chain.Selector, workflowRegistryAddr.String())

		return DeployWorkflowRegistryOpOutput{
			Address:       workflowRegistryAddr.String(),
			ChainSelector: input.ChainSelector,
			Qualifier:     input.Qualifier,
			Type:          string(tv.Type),
			Version:       tv.Version.String(),
			Labels:        tv.Labels.List(),
		}, nil
	},
)
