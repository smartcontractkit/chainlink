package offchain

import (
	"context"
	"fmt"
	"strconv"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodeapiv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	jdtypesv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/pointer"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

type NodesFilter struct {
	DONID        uint64
	EnvLabel     string
	ProductLabel string
	Size         int
	IsBootstrap  bool
}

func (f *NodesFilter) filter() *nodeapiv1.ListNodesRequest_Filter {
	selectors := []*jdtypesv1.Selector{
		{
			Key:   "don_id",
			Op:    jdtypesv1.SelectorOp_EQ,
			Value: pointer.To(strconv.FormatUint(f.DONID, 10)),
		},
		{
			Key:   "environment",
			Op:    jdtypesv1.SelectorOp_EQ,
			Value: &f.EnvLabel,
		},
		{
			Key:   "product",
			Op:    jdtypesv1.SelectorOp_EQ,
			Value: &f.ProductLabel,
		},
	}

	if f.IsBootstrap {
		selectors = append(selectors, &jdtypesv1.Selector{
			Key:   devenv.LabelNodeTypeKey,
			Op:    jdtypesv1.SelectorOp_EQ,
			Value: pointer.To(devenv.LabelNodeTypeValueBootstrap),
		})
	} else {
		selectors = append(selectors, &jdtypesv1.Selector{
			Key:   devenv.LabelNodeTypeKey,
			Op:    jdtypesv1.SelectorOp_EQ,
			Value: pointer.To(devenv.LabelNodeTypeValuePlugin),
		})
	}

	return &nodeapiv1.ListNodesRequest_Filter{
		Selectors: selectors,
	}
}

func fetchNodesFromJD(ctx context.Context, env deployment.Environment, nodeFilters *NodesFilter) (nodes []*nodeapiv1.Node, err error) {
	filter := nodeFilters.filter()

	resp, err := env.Offchain.ListNodes(ctx, &nodeapiv1.ListNodesRequest{Filter: filter})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes for DON %w", err)
	}
	if len(resp.Nodes) != nodeFilters.Size {
		return nil, fmt.Errorf("expected %d nodes for DON(%d), got %d", nodeFilters.Size, nodeFilters.DONID, len(resp.Nodes))
	}

	return resp.Nodes, nil
}

func ProposeJobs(ctx context.Context, env deployment.Environment, workflowJobSpec string, workflowName *string, nodeFilters *NodesFilter) (deployment.ChangesetOutput, error) {
	out := deployment.ChangesetOutput{
		Jobs: []deployment.ProposedJob{},
	}
	// Fetch nodes
	nodes, err := fetchNodesFromJD(ctx, env, nodeFilters)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get nodes: %w", err)
	}
	// Propose jobs
	jobLabels := []*ptypes.Label{
		&ptypes.Label{
			Key:   "don_id",
			Value: pointer.To(strconv.FormatUint(nodeFilters.DONID, 10)),
		},
	}
	if workflowName != nil {
		jobLabels = append(jobLabels, &ptypes.Label{
			Key:   "workflow_name",
			Value: workflowName,
		})
	}

	for _, node := range nodes {
		env.Logger.Debugf("Proposing jof for node %s", node.Name)
		resp, err := env.Offchain.ProposeJob(ctx, &jobv1.ProposeJobRequest{
			NodeId: node.Id,
			Spec:   workflowJobSpec,
			Labels: jobLabels,
		})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to propose job: %w", err)
		}
		env.Logger.Debugf("Job proposed %s", resp.Proposal.JobId)
		out.Jobs = append(out.Jobs, deployment.ProposedJob{
			JobID: resp.Proposal.JobId,
			Node:  node.Id,
			Spec:  resp.Proposal.Spec,
		})
	}
	return out, nil
}

func DeleteJobs(ctx context.Context, env deployment.Environment, jobIDs []string, workflowName string) {
	if len(jobIDs) == 0 {
		env.Logger.Debugf("jobIDs not present. Listing jobs to delete via workflow name")
		jobSelectors := []*jdtypesv1.Selector{
			{
				Key:   "workflow_name",
				Op:    jdtypesv1.SelectorOp_EQ,
				Value: &workflowName,
			},
		}

		listJobResponse, err := env.Offchain.ListJobs(ctx, &jobv1.ListJobsRequest{
			Filter: &jobv1.ListJobsRequest_Filter{
				Selectors: jobSelectors,
			},
		})
		if err != nil {
			env.Logger.Errorf("Failed to list jobs before deleting: %v", err)
			return
		}
		for _, job := range listJobResponse.Jobs {
			jobIDs = append(jobIDs, job.Id)
		}
	}

	for _, jobID := range jobIDs {
		env.Logger.Debugf("Deleting job %s", jobID)
		_, err := env.Offchain.DeleteJob(ctx, &jobv1.DeleteJobRequest{
			IdOneof: &jobv1.DeleteJobRequest_Id{Id: jobID},
		})
		if err != nil {
			env.Logger.Errorf("Failed to delete job %s: %v", jobID, err)
		} else {
			env.Logger.Debugf("Job %s deleted)", jobID)
		}
	}
}
