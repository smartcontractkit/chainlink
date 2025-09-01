package operations

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

type UploadStandardCapabilityJobDeps struct {
	Env cldf.Environment
}

type UploadStandardCapabilityJobInput struct {
	Job       pkg.StandardCapabilityJob
	TargetDON *offchain.DONFilter
}

type UploadStandardCapabilityJobOutput struct {
	Specs map[string][]string
}

var UploadStandardCapabilityJob = operations.NewOperation[UploadStandardCapabilityJobInput, UploadStandardCapabilityJobOutput, UploadStandardCapabilityJobDeps](
	"upload-standard-capability-job-op",
	semver.MustParse("1.0.0"),
	"Upload Standard Capability Job",
	func(b operations.Bundle, deps UploadStandardCapabilityJobDeps, input UploadStandardCapabilityJobInput) (UploadStandardCapabilityJobOutput, error) {
		if err := input.Job.Validate(); err != nil {
			return UploadStandardCapabilityJobOutput{}, fmt.Errorf("invalid job: %w", err)
		}

		spec, err := input.Job.Resolve()
		if err != nil {
			return UploadStandardCapabilityJobOutput{}, fmt.Errorf("failed to resolve job: %w", err)
		}

		specs, err := pkg.ProposeJob(b.GetContext(), deps.Env, pkg.ProposeJobRequest{
			Spec: spec,
			JobLabels: map[string]string{
				offchain.CapabilityLabel: input.Job.JobName,
			},
			TargetDON: input.TargetDON,
		})
		if err != nil {
			return UploadStandardCapabilityJobOutput{}, fmt.Errorf("failed to propose job: %w", err)
		}

		return UploadStandardCapabilityJobOutput{
			Specs: specs,
		}, nil
	},
)
