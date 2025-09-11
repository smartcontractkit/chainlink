package jobs

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	operations2 "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
)

var _ cldf.ChangeSetV2[ProposeJobSpecInput] = ProposeJobSpec{}

type ProposeJobSpecInput struct {
	Environment string `json:"environment" yaml:"environment"`
	Domain      string `json:"domain" yaml:"domain"`

	DONName    string                        `json:"don_name" yaml:"don_name"`
	DONFilters []operations2.TargetDONFilter `json:"don_filters" yaml:"don_filters"`

	JobName     string                    `json:"job_name" yaml:"job_name"`
	Template    job_types.JobSpecTemplate `json:"template" yaml:"template"`
	ExtraLabels map[string]string         `json:"extra_labels,omitempty" yaml:"extra_labels,omitempty"`

	// Inputs is a map of input variables to be used in the job spec template.
	// These will vary based on the template used, and will be validated differently
	// for each template type.
	Inputs job_types.JobSpecInput `json:"inputs" yaml:"inputs"`
}

type ProposeJobSpec struct{}

func (u ProposeJobSpec) VerifyPreconditions(_ cldf.Environment, config ProposeJobSpecInput) error {
	if config.Environment == "" {
		return errors.New("environment is required")
	}

	if config.Domain == "" {
		return errors.New("domain is required")
	}

	if config.DONName == "" {
		return errors.New("don_name is required")
	}

	if len(config.DONFilters) == 0 {
		return errors.New("don_filters is required")
	}

	if config.JobName == "" {
		return errors.New("job_name is required")
	}

	switch config.Template {
	case job_types.Cron, job_types.BootstrapOCR3:
	default:
		return fmt.Errorf("unsupported template: %s", config.Template)
	}

	if config.Inputs == nil {
		return errors.New("inputs are required")
	}

	return nil
}

func (u ProposeJobSpec) Apply(e cldf.Environment, input ProposeJobSpecInput) (cldf.ChangesetOutput, error) {
	e.Logger.Debugw("environment", "name", e.Name)
	var report operations.Report[any, any]
	switch input.Template {
	case job_types.Cron: // This will hold all standard capabilities jobs as we add support for them.
		job, err := input.Inputs.ToStandardCapabilityJob(input.JobName)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert inputs to standard capability job: %w", err)
		}

		r, rErr := operations.ExecuteOperation(
			e.OperationsBundle,
			operations2.ProposeStandardCapabilityJob,
			operations2.ProposeStandardCapabilityJobDeps{Env: e},
			operations2.ProposeStandardCapabilityJobInput{
				Job:         job,
				DONName:     input.DONName,
				DONFilters:  input.DONFilters,
				ExtraLabels: input.ExtraLabels,
			},
		)
		if rErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to propose standard capability job: %w", rErr)
		}

		report = r.ToGenericReport()
	case job_types.BootstrapOCR3:
		jobInput, err := input.Inputs.ToOCR3BootstrapJobInput()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert inputs to OCR3 bootstrap job input: %w", err)
		}

		addrRefKey := pkg.GetOCR3CapabilityV2AddressRefKey(jobInput.ChainSelector, jobInput.ContractQualifier)
		contractAddrRef, err := e.DataStore.Addresses().Get(addrRefKey)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get OCR3 contract address for chain selector %d and qualifier %s: %w", jobInput.ChainSelector, jobInput.ContractQualifier, err)
		}

		r, rErr := operations.ExecuteOperation(
			e.OperationsBundle,
			operations2.ProposeOCR3BootstrapJob,
			operations2.ProposeOCR3BootstrapJobDeps{Env: e},
			operations2.ProposeOCR3BootstrapJobInput{
				Domain:           input.Domain,
				DONName:          input.DONName,
				ContractID:       contractAddrRef.Address,
				EnvironmentLabel: input.Environment,
				ChainSelectorEVM: jobInput.ChainSelector,
				JobName:          input.JobName,
				DONFilters:       input.DONFilters,
				ExtraLabels:      input.ExtraLabels,
			},
		)
		if rErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to propose OCR3 bootstrap job: %w", rErr)
		}

		report = r.ToGenericReport()
	default:
		return cldf.ChangesetOutput{}, fmt.Errorf("unsupported template: %s", input.Template)
	}

	return cldf.ChangesetOutput{
		Reports: []operations.Report[any, any]{report},
	}, nil
}
