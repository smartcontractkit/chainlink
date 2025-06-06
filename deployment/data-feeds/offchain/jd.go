package offchain

import (
	"context"
	"fmt"
	"strconv"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodeapiv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	jdtypesv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/helpers/pointer"
)

type NodesFilter struct {
	DONID        uint64 // Required
	EnvLabel     string
	ProductLabel string
	Size         int
	IsBootstrap  bool
	NodeIDs      []string // Optional, if other filters are provided
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

func fetchNodesFromJD(ctx context.Context, env cldf.Environment, nodeFilters *NodesFilter) (nodes []*nodeapiv1.Node, err error) {
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

func getNodes(ctx context.Context, env cldf.Environment, nodeIDs []string) ([]*nodeapiv1.Node, error) {
	nodes := make([]*nodeapiv1.Node, 0, len(nodeIDs))
	for _, nodeID := range nodeIDs {
		resp, err := env.Offchain.GetNode(ctx, &nodeapiv1.GetNodeRequest{Id: nodeID})
		if err != nil {
			return nil, fmt.Errorf("failed to get node %s: %w", nodeID, err)
		}
		nodes = append(nodes, resp.Node)
	}
	return nodes, nil
}

func ProposeJobs(ctx context.Context, env cldf.Environment, workflowJobSpec string, workflowName *string, nodeFilters *NodesFilter) (cldf.ChangesetOutput, error) {
	out := cldf.ChangesetOutput{
		Jobs: []cldf.ProposedJob{},
	}
	var nodes []*nodeapiv1.Node
	var err error

	// Use node IDs if provided
	if len(nodeFilters.NodeIDs) > 0 {
		env.Logger.Debugf("nodeIDs provided. Fetching nodes for node IDs %s", nodeFilters.NodeIDs)
		nodes, err = getNodes(ctx, env, nodeFilters.NodeIDs)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get nodes for ndoe Ids %s: %w", nodeFilters.NodeIDs, err)
		}
	} else {
		// Fetch nodes based on filter
		nodes, err = fetchNodesFromJD(ctx, env, nodeFilters)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to filter nodes: %w", err)
		}
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
		env.Logger.Debugf("Proposing job for node %s", node.Name)
		resp, err := env.Offchain.ProposeJob(ctx, &jobv1.ProposeJobRequest{
			NodeId: node.Id,
			Spec:   workflowJobSpec,
			Labels: jobLabels,
		})
		if err != nil {
			env.Logger.Errorf("failed to propose job: %s", err)
			continue
		}
		env.Logger.Debugf("Job proposed %s", resp.Proposal.JobId)
		out.Jobs = append(out.Jobs, cldf.ProposedJob{
			JobID: resp.Proposal.JobId,
			Node:  node.Id,
			Spec:  resp.Proposal.Spec,
		})
	}
	return out, nil
}

func DeleteJobs(ctx context.Context, env cldf.Environment, jobIDs []string, workflowName string) {
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
