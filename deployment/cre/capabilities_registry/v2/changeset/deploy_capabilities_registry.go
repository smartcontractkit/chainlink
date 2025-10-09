package changeset

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
)

var _ cldf.ChangeSetV2[DeployCapabilitiesRegistryInput] = DeployCapabilitiesRegistry{}

// DeployCapabilitiesRegistryInput must be JSON Serializable with no private fields
type DeployCapabilitiesRegistryInput struct {
	ChainSelector uint64 `json:"chainSelector"`
	Qualifier     string `json:"qualifier,omitempty"`
}

type DeployCapabilitiesRegistryDeps struct {
	Env *cldf.Environment
}

type DeployCapabilitiesRegistry struct{}

func (l DeployCapabilitiesRegistry) VerifyPreconditions(e cldf.Environment, config DeployCapabilitiesRegistryInput) error {
	if config.ChainSelector == 0 {
		return errors.New("chain selector must be provided")
	}
	_, err := chain_selectors.GetChainIDFromSelector(config.ChainSelector) // validate chain selector
	if err != nil {
		return fmt.Errorf("could not resolve chain selector %d: %w", config.ChainSelector, err)
	}
	return nil
}

func (l DeployCapabilitiesRegistry) Apply(e cldf.Environment, config DeployCapabilitiesRegistryInput) (cldf.ChangesetOutput, error) {
	// build your custom dependencies needed in the operation
	deps := contracts.DeployCapabilitiesRegistryDeps{
		Env: &e,
	}

	reports := make([]operations.Report[any, any], 0)
	capabilitiesRegistryDeploymentReport, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.DeployCapabilitiesRegistry, deps, contracts.DeployCapabilitiesRegistryInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.Qualifier,
		},
		operations.WithRetry[contracts.DeployCapabilitiesRegistryInput, contracts.DeployCapabilitiesRegistryDeps](),
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	reports = append(reports, capabilitiesRegistryDeploymentReport.ToGenericReport())

	// Create datastore and populate it with the deployed contract information
	ds := datastore.NewMemoryDataStore()

	// Parse the version string back to semver.Version
	version, err := semver.NewVersion(capabilitiesRegistryDeploymentReport.Output.Version)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// Create labels from the operation output
	labels := datastore.NewLabelSet()
	for _, label := range capabilitiesRegistryDeploymentReport.Output.Labels {
		labels.Add(label)
	}

	addressRef := datastore.AddressRef{
		ChainSelector: capabilitiesRegistryDeploymentReport.Output.ChainSelector,
		Address:       capabilitiesRegistryDeploymentReport.Output.Address,
		Type:          datastore.ContractType(capabilitiesRegistryDeploymentReport.Output.Type),
		Version:       version,
		Qualifier:     capabilitiesRegistryDeploymentReport.Output.Qualifier,
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
