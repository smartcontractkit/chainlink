package pkg

import (
	"context"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	nodeapiv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

type ProposeJobRequest struct {
	Spec      string // toml
	JobLabels map[string]string
	TargetDON *offchain.DONFilter
}

// ProposeJob sends a single job spec to all the nodes in the DON indicated by `req.DonToUpdate`.
func ProposeJob(ctx context.Context, e cldf.Environment, req ProposeJobRequest) (map[string][]string, error) {
	nodes, err := offchain.FetchNodesFromJD(ctx, e.Offchain, req.TargetDON)
	if err != nil {
		return nil, fmt.Errorf("failed to get don nodes: %w", err)
	}

	jobSpecs := map[string][]string{}
	for _, node := range nodes {
		p2pID := offchain.GetP2pLabel(node.GetLabels())

		e.Logger.Debugw("Proposing job", logLabels(req, node)...)
		req := offchain.ProposeJobRequest{
			Job:            req.Spec,
			Domain:         offchain.ProductLabel,
			Environment:    req.TargetDON.EnvLabel,
			NodeLabels:     map[string]string{offchain.P2pIDLabel: p2pID},
			JobLabels:      req.JobLabels,
			OffchainClient: e.Offchain,
			Lggr:           e.Logger,
		}
		err = offchain.ProposeJob(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to propose job: %w", err)
		}

		jobSpecs[node.Id] = append(jobSpecs[node.Id], req.Job)
	}

	return jobSpecs, nil
}

func logLabels(req ProposeJobRequest, node *nodeapiv1.Node) []any {
	p2pID := offchain.GetP2pLabel(node.GetLabels())

	labels := []any{
		"nodeName",
		node.Name,
		"nodeID",
		node.Id,
		"p2pID",
		p2pID,
		"target DON",
		req.TargetDON.DONName,
	}
	for k, v := range req.JobLabels {
		labels = append(labels, k, v)
	}

	return labels
}
