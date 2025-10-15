package jobs

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	operations2 "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/sequences"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

var _ cldf.ChangeSetV2[ProposeJobSpecInput] = ProposeJobSpec{}

type ProposeJobSpecInput struct {
	Environment string `json:"environment" yaml:"environment"`
	Domain      string `json:"domain" yaml:"domain"`

	DONName    string                     `json:"donName" yaml:"donName"`
	DONFilters []offchain.TargetDONFilter `json:"donFilters" yaml:"donFilters"`

	JobName     string                    `json:"jobName" yaml:"jobName"`
	Template    job_types.JobSpecTemplate `json:"template" yaml:"template"`
	ExtraLabels map[string]string         `json:"extraLabels,omitempty" yaml:"extraLabels,omitempty"`

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
	case job_types.EVM:
		if err := verifyEVMJobSpecInputs(config.Inputs); err != nil {
			return fmt.Errorf("invalid inputs for EVM job spec: %w", err)
		}
	case job_types.Cron, job_types.BootstrapOCR3, job_types.OCR3, job_types.Gateway, job_types.HTTPTrigger, job_types.HTTPAction, job_types.BootstrapVault, job_types.Consensus:
	default:
		return fmt.Errorf("unsupported template: %s", config.Template)
	}

	if config.Inputs == nil {
		return errors.New("inputs are required")
	}

	return nil
}

func (u ProposeJobSpec) Apply(e cldf.Environment, input ProposeJobSpecInput) (cldf.ChangesetOutput, error) {
	var report operations.Report[any, any]
	switch input.Template {
	// This will hold all standard capabilities jobs as we add support for them.
	case job_types.EVM, job_types.Cron, job_types.HTTPTrigger, job_types.HTTPAction, job_types.Consensus:
		// Only consensus generates an oracle factory, for now...
		job, err := input.Inputs.ToStandardCapabilityJob(input.JobName, input.Template == job_types.Consensus)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert inputs to standard capability job: %w", err)
		}

		r, rErr := operations.ExecuteSequence(
			e.OperationsBundle,
			operations2.ProposeStandardCapabilityJob,
			operations2.ProposeStandardCapabilityJobDeps{Env: e},
			operations2.ProposeStandardCapabilityJobInput{
				Job:         job,
				Domain:      input.Domain,
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
		jobInput := pkg.BootstrapJobInput{}
		err := input.Inputs.UnmarshalTo(&jobInput)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert inputs to OCR3 bootstrap job input: %w", err)
		}

		addrRefKey := pkg.GetOCR3CapabilityAddressRefKey(uint64(jobInput.ChainSelector), jobInput.ContractQualifier)
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
				ChainSelectorEVM: uint64(jobInput.ChainSelector),
				JobName:          input.JobName,
				DONFilters:       input.DONFilters,
				ExtraLabels:      input.ExtraLabels,
			},
		)
		if rErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to propose OCR3 bootstrap job: %w", rErr)
		}

		report = r.ToGenericReport()
	case job_types.OCR3:
		jobInput, err := input.Inputs.ToOCR3JobConfigInput()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert inputs to OCR3 job input: %w", err)
		}

		addrRefKey := pkg.GetOCR3CapabilityAddressRefKey(uint64(jobInput.ChainSelectorEVM), jobInput.ContractQualifier)
		contractAddrRef, err := e.DataStore.Addresses().Get(addrRefKey)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get OCR3 contract address for chain selector %d and qualifier %s: %w", jobInput.ChainSelectorEVM, jobInput.ContractQualifier, err)
		}

		dkgContractAddr := ""
		if jobInput.DKGContractQualifier != "" {
			dkgContractRefKey := pkg.GetOCR3CapabilityAddressRefKey(uint64(jobInput.ChainSelectorEVM), jobInput.DKGContractQualifier)
			addr, err := e.DataStore.Addresses().Get(dkgContractRefKey)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get OCR3 contract address for chain selector %d and qualifier %s: %w", jobInput.ChainSelectorEVM, jobInput.ContractQualifier, err)
			}

			dkgContractAddr = addr.Address
		}

		r, rErr := operations.ExecuteSequence(
			e.OperationsBundle,
			operations2.ProposeOCR3Job,
			operations2.ProposeOCR3JobDeps{Env: e},
			operations2.ProposeOCR3JobInput{
				Domain:               input.Domain,
				EnvName:              input.Environment,
				DONName:              input.DONName,
				JobName:              input.JobName,
				TemplateName:         jobInput.TemplateName,
				ContractAddress:      contractAddrRef.Address,
				ChainSelectorEVM:     uint64(jobInput.ChainSelectorEVM),
				ChainSelectorAptos:   uint64(jobInput.ChainSelectorAptos),
				BootstrapperOCR3Urls: jobInput.BootstrapperOCR3Urls,
				DKGContractAddress:   dkgContractAddr,
				DONFilters:           input.DONFilters,
				ExtraLabels:          input.ExtraLabels,
			},
		)

		if rErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to propose OCR3 job: %w", rErr)
		}

		report = r.ToGenericReport()
	case job_types.Gateway:
		typedInputs := operations2.ProposeGatewayJobInput{
			Domain:     input.Domain,
			DONFilters: input.DONFilters,
			JobLabels:  input.ExtraLabels,
		}
		err := input.Inputs.UnmarshalTo(&typedInputs)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to unmarshal inputs to gateway job input: %w", err)
		}

		r, rErr := operations.ExecuteOperation(
			e.OperationsBundle,
			operations2.ProposeGatewayJob,
			operations2.ProposeGatewayJobDeps{Env: e},
			typedInputs,
		)
		if rErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to propose gateway job: %w", rErr)
		}

		report = r.ToGenericReport()
	case job_types.BootstrapVault:
		jobInput := pkg.VaultBootstrapJobsInput{}
		err := input.Inputs.UnmarshalTo(&jobInput)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert inputs to OCR3 bootstrap job input: %w", err)
		}

		r, rErr := operations.ExecuteSequence(
			e.OperationsBundle,
			sequences.ProposeVaultBootstrapJobs,
			sequences.ProposeVaultBootstrapJobsDeps{Env: e},
			sequences.ProposeVaultBootstrapJobsInput{
				Domain:                  input.Domain,
				DONName:                 input.DONName,
				ContractQualifierPrefix: jobInput.ContractQualifierPrefix,
				EnvironmentLabel:        input.Environment,
				ChainSelectorEVM:        uint64(jobInput.ChainSelector),
				JobName:                 input.JobName,
				DONFilters:              input.DONFilters,
				ExtraLabels:             input.ExtraLabels,
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
