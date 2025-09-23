package operations

import (
	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

const (
	BootstrapNodeTypeKey = "bootstrap"
	PluginNodeType       = "plugin"
)

type ProposeJobSpecDeps struct {
	Env cldf.Environment
}

type ProposeJobSpecInput struct {
	Domain  string
	DONName string

	Spec string

	DONFilters []offchain.TargetDONFilter
	JobLabels  map[string]string

	IsBootstrap bool
}

type ProposeJobSpecOutput struct {
	Specs map[string][]string
}

var ProposeJobSpec = operations.NewOperation[ProposeJobSpecInput, ProposeJobSpecOutput, ProposeJobSpecDeps](
	"propose-job-spec-op",
	semver.MustParse("1.0.0"),
	"Propose Job Spec",
	func(b operations.Bundle, deps ProposeJobSpecDeps, input ProposeJobSpecInput) (ProposeJobSpecOutput, error) {
		b.Logger.Debugw("Proposing job", "DON", input.DONName, "domain", input.Domain, "environment", deps.Env.Name)
		req := pkg.ProposeJobRequest{
			Spec:      input.Spec,
			DONName:   input.DONName,
			Env:       deps.Env.Name,
			JobLabels: input.JobLabels,
		}

		nodeType := PluginNodeType
		if input.IsBootstrap {
			nodeType = BootstrapNodeTypeKey
		}
		filter := &node.ListNodesRequest_Filter{
			Selectors: []*ptypes.Selector{
				{
					Key:   "type",
					Op:    ptypes.SelectorOp_EQ,
					Value: &nodeType,
				},
			},
		}
		for _, f := range input.DONFilters {
			filter = f.AddToFilter(filter)
		}

		req.DONFilter = filter

		specs, err := pkg.ProposeJob(b.GetContext(), deps.Env, req)
		if err != nil {
			return ProposeJobSpecOutput{}, err
		}

		return ProposeJobSpecOutput{
			Specs: specs,
		}, nil
	},
)
