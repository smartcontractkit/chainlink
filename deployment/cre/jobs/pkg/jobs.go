package pkg

import (
	"context"
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	nodeapiv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

type ProposeJobRequest struct {
	Spec      string // toml
	DONName   string
	Env       string // staging, testnet, mainnet, etc...
	JobLabels map[string]string
	DONFilter *nodeapiv1.ListNodesRequest_Filter
}

// ProposeJob sends a single job spec to all the nodes in the DON indicated by `req.DONFilter`.
func ProposeJob(ctx context.Context, e cldf.Environment, req ProposeJobRequest) (map[string][]string, error) {
	nodes, err := offchain.FetchNodesFromJD(ctx, e.Offchain, req.DONFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get don nodes: %w", err)
	}
	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes found for DON `%s` with filters %+v", req.DONName, req.DONFilter)
	}

	jobSpecs := map[string][]string{}
	for _, node := range nodes {
		e.Logger.Debugw("Proposing job", logLabels(req, node)...)
		offchainReq := offchain.ProposeJobRequest{
			Job:            req.Spec,
			Domain:         offchain.ProductLabel,
			Environment:    req.Env,
			PublicKeys:     []string{node.GetPublicKey()},
			JobLabels:      req.JobLabels,
			OffchainClient: e.Offchain,
			Lggr:           e.Logger,
			ExtraSelectors: req.DONFilter.GetSelectors(),
		}
		err = offchain.ProposeJob(ctx, offchainReq)
		if err != nil {
			return nil, fmt.Errorf("failed to propose job: %w", err)
		}

		jobSpecs[node.Id] = append(jobSpecs[node.Id], offchainReq.Job)
	}
	if len(jobSpecs) == 0 {
		return nil, errors.New("no jobs were proposed")
	}

	return jobSpecs, nil
}

func logLabels(req ProposeJobRequest, node *nodeapiv1.Node) []any {
	labels := []any{
		"nodeName",
		node.Name,
		"nodeID",
		node.Id,
		"publicKey",
		node.PublicKey,
		"target DON",
		req.DONName,
	}
	for k, v := range req.JobLabels {
		labels = append(labels, k, v)
	}

	return labels
}
