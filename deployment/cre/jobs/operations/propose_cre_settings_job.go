package operations

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

type ProposeCRESettingsJobsDeps struct {
	Env cldf.Environment
}
type ProposeCRESettingsJobsInput struct {
	Domain      string
	DONName     string
	DONFilters  []offchain.TargetDONFilter
	ExtraLabels map[string]string

	Settings string // toml
}

type ProposeCRESettingsJobsOutput struct {
	Specs map[string][]string
}

var ProposeCRESettingsJobs = operations.NewOperation[
	ProposeCRESettingsJobsInput,
	ProposeCRESettingsJobsOutput,
	ProposeCRESettingsJobsDeps,
](
	"propose-cre-settings-job-op",
	semver.MustParse("1.0.0"),
	"Propose CRESettings Job",
	func(b operations.Bundle, deps ProposeCRESettingsJobsDeps, input ProposeCRESettingsJobsInput) (output ProposeCRESettingsJobsOutput, err error) {
		job := pkg.CRESettingsJob{Settings: input.Settings}
		jobSpec, err := job.ResolveJob()
		if err != nil {
			return ProposeCRESettingsJobsOutput{}, fmt.Errorf("failed to resolve job spec: %w", err)
		}

		report, err := operations.ExecuteOperation(b, ProposeJobSpec, ProposeJobSpecDeps(deps), ProposeJobSpecInput{
			Domain:     input.Domain,
			DONName:    input.DONName,
			JobLabels:  input.ExtraLabels,
			DONFilters: input.DONFilters,
			Spec:       jobSpec,
		})
		if err != nil {
			return ProposeCRESettingsJobsOutput{}, fmt.Errorf("failed to propose cre settings job: %w", err)
		}

		return ProposeCRESettingsJobsOutput{
			Specs: report.Output.Specs,
		}, nil
	},
)
