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

type DeployWorkflowRegistryOpDeps struct {
	Env *cldf.Environment
}

type DeployWorkflowRegistryInput struct {
	ChainSelector uint64
}

type DeployWorkflowRegistryOutput struct {
	Addresses   datastore.AddressRefStore
	AddressBook cldf.AddressBook // Keeping the address book for backward compatibility, as not everything has been migrated to datastore
}

// DeployWorkflowRegistryOp is an operation that deploys the Workflow Registry contract.
var DeployWorkflowRegistryOp = operations.NewOperation[DeployWorkflowRegistryInput, DeployWorkflowRegistryOutput, DeployWorkflowRegistryOpDeps](
	"deploy-workflow-registry-op",
	semver.MustParse("1.0.0"),
	"Deploy WorkflowRegistry Contract",
	func(b operations.Bundle, deps DeployWorkflowRegistryOpDeps, input DeployWorkflowRegistryInput) (DeployWorkflowRegistryOutput, error) {
		workfloRegistryOutput, err := workflow_registry_changeset.DeployV2(*deps.Env, &keystone_changeset.DeployRequestV2{
			ChainSel: input.ChainSelector,
		})
		if err != nil {
			return DeployWorkflowRegistryOutput{}, fmt.Errorf("DeployWorkflowRegistryOp error: failed to deploy Workflow Registry contract: %w", err)
		}
		return DeployWorkflowRegistryOutput{
			Addresses:   workfloRegistryOutput.DataStore.Addresses(),
			AddressBook: workfloRegistryOutput.AddressBook, //nolint:staticcheck // keeping the address book since not everything has been migrated to datastore
		}, nil
	},
)
