package changeset

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/workflow_registry/v2/changeset/operations/contracts"
)

var _ cldf.ChangeSetV2[DeployWorkflowRegistryInput] = DeployWorkflowRegistry{}

// DeployWorkflowRegistryInput must be JSON Serializable with no private fields
type DeployWorkflowRegistryInput struct {
	ChainSelector uint64 `json:"chainSelector"`
	Qualifier     string `json:"qualifier,omitempty"`
}

type DeployWorkflowRegistryDeps struct {
	Env *cldf.Environment
}

type DeployWorkflowRegistry struct{}

func (l DeployWorkflowRegistry) VerifyPreconditions(e cldf.Environment, config DeployWorkflowRegistryInput) error {
	return nil
}

func (l DeployWorkflowRegistry) Apply(e cldf.Environment, config DeployWorkflowRegistryInput) (cldf.ChangesetOutput, error) {
	// build your custom dependencies needed in the operation
	deps := contracts.DeployWorkflowRegistryOpDeps{
		Env: &e,
	}

	reports := make([]operations.Report[any, any], 0)
	workflowRegistryDeploymentReport, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.DeployWorkflowRegistryOp, deps, contracts.DeployWorkflowRegistryOpInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.Qualifier,
		},
		operations.WithRetry[contracts.DeployWorkflowRegistryOpInput, contracts.DeployWorkflowRegistryOpDeps](),
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	reports = append(reports, workflowRegistryDeploymentReport.ToGenericReport())

	// Create datastore and populate it with the deployed contract information
	ds := datastore.NewMemoryDataStore()

	// Parse the version string back to semver.Version
	version, err := semver.NewVersion(workflowRegistryDeploymentReport.Output.Version)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// Create labels from the operation output
	labels := datastore.NewLabelSet()
	for _, label := range workflowRegistryDeploymentReport.Output.Labels {
		labels.Add(label)
	}

	addressRef := datastore.AddressRef{
		ChainSelector: workflowRegistryDeploymentReport.Output.ChainSelector,
		Address:       workflowRegistryDeploymentReport.Output.Address,
		Type:          datastore.ContractType(workflowRegistryDeploymentReport.Output.Type),
		Version:       version,
		Qualifier:     workflowRegistryDeploymentReport.Output.Qualifier,
		Labels:        labels,
	}

	if err := ds.Addresses().Add(addressRef); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		DataStore: ds,
		Reports:   reports,
	}, nil
}
