package jobs

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	operations2 "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

var _ cldf.ChangeSetV2[DistributeJobSpecInput] = DistributeJobSpec{}

type JobSpecTemplate string

const (
	Cron JobSpecTemplate = "cron"
)

type DistributeJobSpecInput struct {
	Environment string `json:"environment" yaml:"environment"`
	Domain      string `json:"domain" yaml:"domain"`

	TargetDON *offchain.DONFilter `json:"target_don" yaml:"target_don"`

	JobName  string          `json:"job_name" yaml:"job_name"`
	Template JobSpecTemplate `json:"template" yaml:"template"`

	// Inputs is a map of input variables to be used in the job spec template.
	// These will vary based on the template used, and will be validated differently
	// for each template type.
	Inputs types.JobSpecInput `json:"inputs" yaml:"inputs"`
}

type DistributeJobSpec struct{}

func (u DistributeJobSpec) VerifyPreconditions(e cldf.Environment, config DistributeJobSpecInput) error {
	if config.Environment == "" {
		return errors.New("environment is required")
	}

	if config.Domain == "" {
		return errors.New("domain is required")
	}

	switch config.Template {
	case Cron:
	default:
		return fmt.Errorf("unsupported template: %s", config.Template)
	}

	if config.Inputs == nil {
		return errors.New("inputs are required")
	}

	return nil
}

func (u DistributeJobSpec) Apply(e cldf.Environment, input DistributeJobSpecInput) (cldf.ChangesetOutput, error) {
	var report operations.Report[any, any]
	switch input.Template {
	case Cron: // This will hold all standard capabilities jobs as we add support for them.
		job, err := input.Inputs.ToStandardCapabilityJob(input.JobName)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert inputs to standard capability job: %w", err)
		}

		r, rErr := operations.ExecuteOperation(
			e.OperationsBundle,
			operations2.UploadStandardCapabilityJob,
			operations2.UploadStandardCapabilityJobDeps{Env: e},
			operations2.UploadStandardCapabilityJobInput{
				Job:       job,
				TargetDON: input.TargetDON,
			},
		)
		if rErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to distribute standard capability job: %w", rErr)
		}

		report = r.ToGenericReport()
	default:
		return cldf.ChangesetOutput{}, fmt.Errorf("unsupported template: %s", input.Template)
	}

	return cldf.ChangesetOutput{
		Reports: []operations.Report[any, any]{report},
	}, nil
}
