package operations

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

// ProposeRingJobDeps contains the dependencies for proposing a Ring job.
type ProposeRingJobDeps struct {
	Env cldf.Environment
}

// ProposeRingJobInput contains the input for proposing a Ring job.
type ProposeRingJobInput struct {
	Domain  string
	EnvName string

	DONName string
	JobName string

	ContractAddress  string
	ChainSelectorEVM uint64
	ShardConfigAddr  string
	BootstrapperUrls []string

	DONFilters  []offchain.TargetDONFilter
	ExtraLabels map[string]string
}

// ProposeRingJobOutput contains the output of proposing a Ring job.
type ProposeRingJobOutput struct {
	Specs map[string][]string
}

// ProposeRingJob is the sequence for proposing Ring jobs.
var ProposeRingJob = operations.NewSequence[ProposeRingJobInput, ProposeRingJobOutput, ProposeRingJobDeps](
	"propose-ring-job-seq",
	semver.MustParse("1.0.0"),
	"Propose Ring Job",
	func(b operations.Bundle, deps ProposeRingJobDeps, input ProposeRingJobInput) (ProposeRingJobOutput, error) {
		filters := &nodev1.ListNodesRequest_Filter{}
		for _, f := range input.DONFilters {
			filters = offchain.TargetDONFilter{
				Key:   f.Key,
				Value: f.Value,
			}.AddToFilter(filters)
		}
		// We only want to target plugin nodes for Ring jobs.
		filters = offchain.TargetDONFilter{
			Key:   "type",
			Value: "plugin",
		}.AddToFilter(filters)
		nodes, err := pkg.FetchNodesFromJD(b.GetContext(), deps.Env, pkg.FetchNodesRequest{
			Domain:  input.Domain,
			Filters: filters,
		})
		if err != nil {
			return ProposeRingJobOutput{}, fmt.Errorf("failed to fetch nodes from JD: %w", err)
		}

		nodeToCSAKey := make(map[string]string)
		for _, n := range nodes {
			nodeToCSAKey[n.Id] = n.GetPublicKey()
		}

		specs, err := pkg.BuildRingJobConfigSpecs(
			deps.Env.Offchain, deps.Env.Logger, input.ContractAddress, input.ChainSelectorEVM,
			nodes, input.BootstrapperUrls, input.DONName, input.JobName, input.ShardConfigAddr,
		)
		if err != nil {
			return ProposeRingJobOutput{}, fmt.Errorf("failed to build Ring job config specs: %w", err)
		}

		finalSpecs := make(map[string][]string)

		var errs []error
		for _, spec := range specs {
			// Let's limit the target to the specific node for this spec.
			filters := []offchain.TargetDONFilter{
				{
					Key:   offchain.FilterKeyCSAPublicKey,
					Value: nodeToCSAKey[spec.NodeID],
				},
			}
			filters = append(filters, input.DONFilters...)
			opReport, opErr := operations.ExecuteOperation(b, ProposeJobSpec, ProposeJobSpecDeps(deps), ProposeJobSpecInput{
				Domain:     input.Domain,
				DONName:    input.DONName,
				Spec:       spec.Spec,
				DONFilters: filters,
				JobLabels:  input.ExtraLabels,
			})
			if opErr != nil {
				// Do not fail the sequence if a single proposal fails, make it through all proposals.
				errs = append(errs, fmt.Errorf("error proposing Ring job to node %s: %w", spec.NodeID, opErr))
				continue
			}

			for nodeID, s := range opReport.Output.Specs {
				finalSpecs[nodeID] = append(finalSpecs[nodeID], s...)
			}
		}

		return ProposeRingJobOutput{
			Specs: finalSpecs,
		}, errors.Join(errs...)
	},
)
