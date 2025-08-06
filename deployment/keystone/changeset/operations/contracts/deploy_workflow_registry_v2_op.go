package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	workflow_registry_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/workflowregistry"
)

type DeployWorkflowRegistryV2OpDeps struct {
	Env *cldf.Environment
}

type DeployWorkflowRegistryV2Input struct {
	ChainSelector uint64
}

type DeployWorkflowRegistryV2Output struct {
	Addresses datastore.AddressRefStore
}

// DeployWorkflowRegistryV2Op is an operation that deploys the Workflow Registry v2 contract.
var DeployWorkflowRegistryV2Op = operations.NewOperation[DeployWorkflowRegistryV2Input, DeployWorkflowRegistryV2Output, DeployWorkflowRegistryV2OpDeps](
	"deploy-workflow-registry-v2-op",
	semver.MustParse("1.0.0"),
	"Deploy WorkflowRegistry v2 Contract",
	func(b operations.Bundle, deps DeployWorkflowRegistryV2OpDeps, input DeployWorkflowRegistryV2Input) (DeployWorkflowRegistryV2Output, error) {
		workflowRegistryOutput, err := workflow_registry_changeset.DeployWorkflowRegistryV2Single(*deps.Env, &keystone_changeset.DeployRequestV2{
			ChainSel: input.ChainSelector,
		})
		if err != nil {
			return DeployWorkflowRegistryV2Output{}, fmt.Errorf("DeployWorkflowRegistryV2Op error: failed to deploy Workflow Registry v2 contract: %w", err)
		}
		return DeployWorkflowRegistryV2Output{
			Addresses: workflowRegistryOutput.DataStore.Addresses(),
		}, nil
	},
)

// DeployWorkflowRegistryV2MultiChainOp is an operation that deploys the Workflow Registry v2 contract to multiple chains.
var DeployWorkflowRegistryV2MultiChainOp = operations.NewOperation[keystone_changeset.DeployWorkflowRegistryV2Request, DeployWorkflowRegistryV2Output, DeployWorkflowRegistryV2OpDeps](
	"deploy-workflow-registry-v2-multichain-op",
	semver.MustParse("1.0.0"),
	"Deploy WorkflowRegistry v2 Contract to Multiple Chains",
	func(b operations.Bundle, deps DeployWorkflowRegistryV2OpDeps, input keystone_changeset.DeployWorkflowRegistryV2Request) (DeployWorkflowRegistryV2Output, error) {
		workflowRegistryOutput, err := workflow_registry_changeset.DeployV2MultiChain(*deps.Env, workflow_registry_changeset.DeployV2MultiChainRequest{
			ChainSelectors: input.ChainSelectors,
		})
		if err != nil {
			return DeployWorkflowRegistryV2Output{}, fmt.Errorf("DeployWorkflowRegistryV2MultiChainOp error: failed to deploy Workflow Registry v2 contract: %w", err)
		}
		return DeployWorkflowRegistryV2Output{
			Addresses: workflowRegistryOutput.DataStore.Addresses(),
		}, nil
	},
)
