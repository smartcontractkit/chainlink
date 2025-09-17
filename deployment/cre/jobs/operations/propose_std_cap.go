package operations

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

type ProposeStandardCapabilityJobDeps struct {
	Env cldf.Environment
}

type ProposeStandardCapabilityJobInput struct {
	Domain      string
	DONName     string
	Job         pkg.StandardCapabilityJob
	DONFilters  []offchain.TargetDONFilter
	ExtraLabels map[string]string
}

type ProposeStandardCapabilityJobOutput struct {
	Specs map[string][]string
}

var ProposeStandardCapabilityJob = operations.NewOperation[ProposeStandardCapabilityJobInput, ProposeStandardCapabilityJobOutput, ProposeStandardCapabilityJobDeps](
	"propose-standard-capability-job-op",
	semver.MustParse("1.0.0"),
	"Propose Standard Capability Job",
	func(b operations.Bundle, deps ProposeStandardCapabilityJobDeps, input ProposeStandardCapabilityJobInput) (ProposeStandardCapabilityJobOutput, error) {
		if err := input.Job.Validate(); err != nil {
			return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("invalid job: %w", err)
		}

		spec, err := input.Job.Resolve()
		if err != nil {
			return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to resolve job: %w", err)
		}

		jobLabels := map[string]string{
			offchain.CapabilityLabel: input.Job.JobName,
		}
		for k, v := range input.ExtraLabels {
			jobLabels[k] = v
		}

		report, err := operations.ExecuteOperation(b, ProposeJobSpec, ProposeJobSpecDeps(deps), ProposeJobSpecInput{
			Domain:     input.Domain,
			DONName:    input.DONName,
			Spec:       spec,
			JobLabels:  jobLabels,
			DONFilters: input.DONFilters,
		})
		if err != nil {
			return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to propose job: %w", err)
		}

		return ProposeStandardCapabilityJobOutput{
			Specs: report.Output.Specs,
		}, nil
	},
)
