package changeset

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/sequences"
)

var _ cldf.ChangeSetV2[DeployOCR3Input] = DeployOCR3{}

type DeployOCR3Input struct {
	ChainSelector uint64                      `json:"chainSelector" yaml:"chainSelector"`
	Qualifier     string                      `json:"qualifier" yaml:"qualifier"`
	Dons          []contracts.ConfigureCREDON `json:"dons" yaml:"dons"`
	OracleConfig  *ocr3.OracleConfig          `json:"oracleConfig" yaml:"oracleConfig"`
	DryRun        bool                        `json:"dryRun" yaml:"dryRun"`
	MCMSConfig    *ocr3.MCMSConfig            `json:"mcmsConfig" yaml:"mcmsConfig"`
}

type DeployOCR3Deps struct {
	Env *cldf.Environment
}

type DeployOCR3 struct{}

func (l DeployOCR3) VerifyPreconditions(e cldf.Environment, config DeployOCR3Input) error {
	return nil
}

func (l DeployOCR3) Apply(e cldf.Environment, config DeployOCR3Input) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()

	ocr3DeploymentReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		sequences.DeployOCR3,
		sequences.DeployOCR3Deps{Env: &e},
		sequences.DeployOCR3Input{
			RegistryChainSel: config.ChainSelector,
			Qualifier:        config.Qualifier,

			DONs:         config.Dons,
			OracleConfig: config.OracleConfig,
			DryRun:       config.DryRun,
			MCMSConfig:   config.MCMSConfig,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy ocr3 contract: %w", err)
	}

	reports := make([]operations.Report[any, any], 0)
	reports = append(reports, ocr3DeploymentReport.ToGenericReport())

	// Parse the version string back to semver.Version
	version, err := semver.NewVersion(ocr3DeploymentReport.Output.Version)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// Create labels from the operation output
	labels := datastore.NewLabelSet()
	for _, label := range ocr3DeploymentReport.Output.Labels {
		labels.Add(label)
	}

	addressRef := datastore.AddressRef{
		ChainSelector: ocr3DeploymentReport.Output.ChainSelector,
		Address:       ocr3DeploymentReport.Output.Address,
		Type:          datastore.ContractType(ocr3DeploymentReport.Output.Type),
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
