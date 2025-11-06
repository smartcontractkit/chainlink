package changeset

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/sequences"
)

var _ cldf.ChangeSetV2[DeployVaultInput] = DeployVault{}

type DeployVaultInput struct {
	ChainSelector uint64 `json:"chainSelector" yaml:"chainSelector"`
	Qualifier     string `json:"qualifier" yaml:"qualifier"`
}

type DeployVaultDeps struct {
	Env *cldf.Environment
}

type DeployVault struct{}

func (l DeployVault) VerifyPreconditions(e cldf.Environment, config DeployVaultInput) error {
	return nil
}

func (l DeployVault) Apply(e cldf.Environment, config DeployVaultInput) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()

	deploymentReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		sequences.DeployVault,
		sequences.DeployVaultDeps{Env: &e},
		sequences.DeployVaultInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.Qualifier,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy vault contracts: %w", err)
	}

	reports := make([]operations.Report[any, any], 0)
	reports = append(reports, deploymentReport.ToGenericReport())

	// Parse the version string back to semver.Version
	version, err := semver.NewVersion(deploymentReport.Output.Version)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// Create labels from the operation output
	labels := datastore.NewLabelSet()
	for _, label := range deploymentReport.Output.Labels {
		labels.Add(label)
	}

	addressRef := datastore.AddressRef{
		ChainSelector: deploymentReport.Output.ChainSelector,
		Address:       deploymentReport.Output.PluginAddress,
		Qualifier:     deploymentReport.Output.PluginQualifier,
		Type:          datastore.ContractType(deploymentReport.Output.Type),
		Version:       version,
		Labels:        labels,
	}

	if err := ds.Addresses().Add(addressRef); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	addressRef = datastore.AddressRef{
		ChainSelector: deploymentReport.Output.ChainSelector,
		Address:       deploymentReport.Output.DKGAddress,
		Qualifier:     deploymentReport.Output.DKGQualifier,
		Type:          datastore.ContractType(deploymentReport.Output.Type),
		Version:       version,
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
