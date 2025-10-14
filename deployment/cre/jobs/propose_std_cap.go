package jobs

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	operations2 "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

var _ cldf.ChangeSetV2[ProposeStandardCapabilityJobInput] = ProposeStandardCapabilityJob{}

type ProposeStandardCapabilityJobInput struct {
	Domain  string `json:"domain" yaml:"domain"`
	DONName string `json:"donName" yaml:"donName"`
	JobName string `json:"jobName" yaml:"jobName"`
	Command string `json:"command" yaml:"command"`
	Config  string `json:"config" yaml:"config"`

	ExternalJobID         string            `json:"externalJobID" yaml:"externalJobID"`                 // Optional
	OracleFactory         pkg.OracleFactory `json:"oracleFactory" yaml:"oracleFactory"`                 // Optional
	GenerateOracleFactory bool              `json:"generateOracleFactory" yaml:"generateOracleFactory"` // Optional

	DONFilters  []offchain.TargetDONFilter `json:"donFilters" yaml:"donFilters"`
	ExtraLabels map[string]string          `json:"extraLabels,omitempty" yaml:"extraLabels,omitempty"`
}

type ProposeStandardCapabilityJob struct{}

func (u ProposeStandardCapabilityJob) VerifyPreconditions(_ cldf.Environment, config ProposeStandardCapabilityJobInput) error {
	if config.JobName == "" {
		return errors.New("jobName is required")
	}
	if config.Command == "" {
		return errors.New("command is required")
	}
	if config.DONName == "" {
		return errors.New("don_name is required")
	}
	if len(config.DONFilters) == 0 {
		return errors.New("DONFilters is required")
	}
	return nil
}

func (u ProposeStandardCapabilityJob) Apply(e cldf.Environment, input ProposeStandardCapabilityJobInput) (cldf.ChangesetOutput, error) {
	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		operations2.ProposeStandardCapabilityJob,
		operations2.ProposeStandardCapabilityJobDeps{Env: e},
		operations2.ProposeStandardCapabilityJobInput{
			Domain:  input.Domain,
			DONName: input.DONName,
			Job: pkg.StandardCapabilityJob{
				JobName:               input.JobName,
				Command:               input.Command,
				Config:                input.Config,
				ExternalJobID:         input.ExternalJobID,
				OracleFactory:         &input.OracleFactory,
				GenerateOracleFactory: input.GenerateOracleFactory,
			},
			DONFilters:  input.DONFilters,
			ExtraLabels: input.ExtraLabels,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to propose standard capability job: %w", err)
	}

	return cldf.ChangesetOutput{
		Reports: []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}
