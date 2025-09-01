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

var _ cldf.ChangeSetV2[UploadStandardCapabilityJobInput] = UploadStandardCapabilityJob{}

type UploadStandardCapabilityJobInput struct {
	JobName string `json:"jobName" yaml:"jobName"`
	Command string `json:"command" yaml:"command"`
	Config  string `json:"config" yaml:"config"`

	ExternalJobID string            `json:"externalJobID" yaml:"externalJobID"` // Optional
	OracleFactory pkg.OracleFactory `json:"oracleFactory" yaml:"oracleFactory"` // Optional

	TargetDON *offchain.DONFilter `json:"targetDON" yaml:"targetDON"`
}

type UploadStandardCapabilityJob struct{}

func (u UploadStandardCapabilityJob) VerifyPreconditions(e cldf.Environment, config UploadStandardCapabilityJobInput) error {
	if config.JobName == "" {
		return errors.New("jobName is required")
	}
	if config.Command == "" {
		return errors.New("command is required")
	}
	if config.TargetDON == nil {
		return errors.New("targetDON is required")
	}
	return nil
}

func (u UploadStandardCapabilityJob) Apply(e cldf.Environment, input UploadStandardCapabilityJobInput) (cldf.ChangesetOutput, error) {
	report, err := operations.ExecuteOperation(
		e.OperationsBundle,
		operations2.UploadStandardCapabilityJob,
		operations2.UploadStandardCapabilityJobDeps{Env: e},
		operations2.UploadStandardCapabilityJobInput{
			Job: pkg.StandardCapabilityJob{
				JobName:       input.JobName,
				Command:       input.Command,
				Config:        input.Config,
				ExternalJobID: input.ExternalJobID,
				OracleFactory: input.OracleFactory,
			},
			TargetDON: input.TargetDON,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to upload standard capability job: %w", err)
	}

	return cldf.ChangesetOutput{
		Reports: []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}
