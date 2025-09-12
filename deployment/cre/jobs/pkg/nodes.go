package pkg

import (
	"context"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

type FetchNodesRequest struct {
	Domain  string
	Filters []offchain.TargetDONFilter
}

func FetchNodesFromJD(ctx context.Context, e cldf.Environment, req FetchNodesRequest) ([]*nodev1.Node, error) {
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
				Value: &e.Name,
			},
		},
	}

	for _, f := range req.Filters {
		filter.Selectors = append(filter.Selectors, f.ToJDSelector())
	}

	return offchain.FetchNodesFromJD(ctx, e.Offchain, filter)
}
