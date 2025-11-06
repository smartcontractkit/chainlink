package changeset

import (
	"errors"
	"fmt"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/operations/contracts"
)

var _ cldf.ChangeSetV2[DeployOCR3Input] = DeployOCR3{}

type DeployOCR3Input struct {
	ChainSelector uint64   `json:"chainSelector" yaml:"chainSelector"`
	Qualifier     string   `json:"qualifier" yaml:"qualifier"`
	Labels        []string `json:"labels" yaml:"labels"`
}

type DeployOCR3Deps struct {
	Env *cldf.Environment
}

type DeployOCR3 struct{}

func (l DeployOCR3) VerifyPreconditions(_ cldf.Environment, input DeployOCR3Input) error {
	if input.ChainSelector == 0 {
		return errors.New("chainSelector is required")
	}
	_, err := chain_selectors.GetChainIDFromSelector(input.ChainSelector) // validate chain selector
	if err != nil {
		return fmt.Errorf("could not resolve chain selector %d: %w", input.ChainSelector, err)
	}
	if input.Qualifier == "" {
		return errors.New("qualifier is required")
	}
	return nil
}

func (l DeployOCR3) Apply(e cldf.Environment, config DeployOCR3Input) (cldf.ChangesetOutput, error) {
	ocr3DeploymentReport, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.DeployOCR3,
		contracts.DeployOCR3Deps{Env: &e},
		contracts.DeployOCR3Input{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.Qualifier,
			Labels:        config.Labels,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy ocr3 contract: %w", err)
	}

	return cldf.ChangesetOutput{
		DataStore: ocr3DeploymentReport.Output.Datastore,
		Reports:   []operations.Report[any, any]{ocr3DeploymentReport.ToGenericReport()},
	}, nil
}
