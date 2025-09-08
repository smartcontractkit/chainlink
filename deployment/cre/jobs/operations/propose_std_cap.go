package operations

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	"github.com/smartcontractkit/chainlink/deployment/helpers/pointer"
)

const FilterKeyDONName = "don_name"

type ProposeStandardCapabilityJobDeps struct {
	Env cldf.Environment
}

type ProposeStandardCapabilityJobInput struct {
	DONName     string
	Job         pkg.StandardCapabilityJob
	DONFilters  []TargetDONFilter
	ExtraLabels map[string]string
}

type TargetDONFilter struct {
	Key   string
	Value string
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

		filter := &node.ListNodesRequest_Filter{
			Selectors: []*ptypes.Selector{
				{
					Key:   "type",
					Op:    ptypes.SelectorOp_EQ,
					Value: pointer.To("plugin"),
				},
			},
		}
		for _, f := range input.DONFilters {
			// DON name is a key, so we just check for its existence instead of equality
			if f.Key == FilterKeyDONName {
				filter.Selectors = append(filter.Selectors, &ptypes.Selector{
					Op:  ptypes.SelectorOp_EXIST,
					Key: f.Value,
				})
			} else {
				filter.Selectors = append(filter.Selectors, &ptypes.Selector{
					Op:    ptypes.SelectorOp_EQ,
					Key:   f.Key,
					Value: &f.Value,
				})
			}
		}

		specs, err := pkg.ProposeJob(b.GetContext(), deps.Env, pkg.ProposeJobRequest{
			Spec:      spec,
			DONName:   input.DONName,
			Env:       deps.Env.Name,
			JobLabels: jobLabels,
			DONFilter: filter,
		})
		if err != nil {
			return ProposeStandardCapabilityJobOutput{}, fmt.Errorf("failed to propose job: %w", err)
		}

		return ProposeStandardCapabilityJobOutput{
			Specs: specs,
		}, nil
	},
)
