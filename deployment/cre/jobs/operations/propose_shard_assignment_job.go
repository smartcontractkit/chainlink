package operations

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

type ProposeShardAssignmentJobDeps struct {
	Env cldf.Environment
}

type ProposeShardAssignmentJobInput struct {
	Domain          string
	Environment     string
	DONName         string
	DONFilters      []offchain.TargetDONFilter
	ExtraLabels     map[string]string
	ShardAssignment string // toml
}

type ProposeShardAssignmentJobOutput struct {
	Specs map[string][]string
}

var ProposeShardAssignmentJob = operations.NewOperation[
	ProposeShardAssignmentJobInput,
	ProposeShardAssignmentJobOutput,
	ProposeShardAssignmentJobDeps,
](
	"propose-shard-assignment-job-op",
	semver.MustParse("1.0.0"),
	"Propose ShardAssignment CRESettings Job",
	func(b operations.Bundle, deps ProposeShardAssignmentJobDeps, input ProposeShardAssignmentJobInput) (output ProposeShardAssignmentJobOutput, err error) {
		job := pkg.ShardAssignmentJob{ShardAssignment: input.ShardAssignment}
		jobSpec, err := job.ResolveJob()
		if err != nil {
			return ProposeShardAssignmentJobOutput{}, fmt.Errorf("failed to resolve shard assignment job spec: %w", err)
		}

		report, err := operations.ExecuteOperation(b, ProposeJobSpec, ProposeJobSpecDeps(deps), ProposeJobSpecInput{
			Domain:      input.Domain,
			Environment: input.Environment,
			DONName:     input.DONName,
			JobLabels:   input.ExtraLabels,
			DONFilters:  input.DONFilters,
			Spec:        jobSpec,
		})
		if err != nil {
			return ProposeShardAssignmentJobOutput{}, fmt.Errorf("failed to propose shard assignment job: %w", err)
		}

		return ProposeShardAssignmentJobOutput{
			Specs: report.Output.Specs,
		}, nil
	},
)
