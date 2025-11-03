package pkg

import (
	"context"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

type NodesByLabelsRequest struct {
	Domain      string
	Environment string
	Filters     []offchain.NodeLabelFilter
}

func (req NodesByLabelsRequest) Filter() *nodev1.ListNodesRequest_Filter {
	filter := &nodev1.ListNodesRequest_Filter{
		Selectors: []*ptypes.Selector{
			{
				Key:   "product",
				Op:    ptypes.SelectorOp_EQ,
				Value: &req.Domain,
			},
			{
				Key:   "environment",
				Op:    ptypes.SelectorOp_EQ,
				Value: &req.Environment,
			},
		},
	}

	for _, f := range req.Filters {
		filter = f.AddToFilter(filter)
	}
	return filter
}

type FetchNodeChainConfigsResponse struct {
	NodeID       string
	ChainConfigs []*nodev1.ChainConfig
}

func FetchNodeChainConfigsFromJD(ctx context.Context, e cldf.Environment, filter *nodev1.ListNodesRequest_Filter) ([]FetchNodeChainConfigsResponse, error) {
	resp, err := e.Offchain.ListNodes(ctx, &nodev1.ListNodesRequest{Filter: filter})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	nodeIDs := []string{}
	for _, n := range resp.Nodes {
		nodeIDs = append(nodeIDs, n.Id)
	}

	chainConfigResp, err := e.Offchain.ListNodeChainConfigs(
		ctx,
		&nodev1.ListNodeChainConfigsRequest{
			Filter: &nodev1.ListNodeChainConfigsRequest_Filter{NodeIds: nodeIDs},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain config for nodes: %w", err)
	}

	m := map[string][]*nodev1.ChainConfig{}
	for _, cc := range chainConfigResp.ChainConfigs {
		if _, ok := m[cc.NodeId]; !ok {
			m[cc.NodeId] = []*nodev1.ChainConfig{}
		}

		m[cc.NodeId] = append(m[cc.NodeId], cc)
	}

	fetchResp := []FetchNodeChainConfigsResponse{}
	for nid, ccfgs := range m {
		fetchResp = append(fetchResp, FetchNodeChainConfigsResponse{
			NodeID:       nid,
			ChainConfigs: ccfgs,
		})
	}

	return fetchResp, err
}
