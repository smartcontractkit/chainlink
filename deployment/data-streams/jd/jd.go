package jd

import (
	"context"
	"fmt"

	nodeapiv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	jdtypesv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/pointer"
)

const (
	ProductLabel = "data-streams"
)

// Source for a set of JD filters to apply when fetching a DON.
// Should uniquely identify a set of nodes in JD which belong to a single DON.
type ListFilter struct {
	DONID    uint64
	DONName  string
	EnvLabel string
	Size     int // Expected number of nodes in the result
}

func (f *ListFilter) bootstrappersFilter() *nodeapiv1.ListNodesRequest_Filter {
	return &nodeapiv1.ListNodesRequest_Filter{
		Selectors: []*jdtypesv1.Selector{
			{
				Key: fmt.Sprintf("don-%d-%s", f.DONID, f.DONName),
				Op:  jdtypesv1.SelectorOp_EXIST,
			},
			{
				Key:   "nodeType",
				Op:    jdtypesv1.SelectorOp_EQ,
				Value: pointer.To("bootstrap"),
			},
			{
				Key:   "environment",
				Op:    jdtypesv1.SelectorOp_EQ,
				Value: &f.EnvLabel,
			},
			{
				Key:   "product",
				Op:    jdtypesv1.SelectorOp_EQ,
				Value: pointer.To("data-streams"),
			},
		},
	}
}

func FetchDONBootstrappersFromJD(ctx context.Context, jd deployment.OffchainClient, filter *ListFilter) (nodes []*nodeapiv1.Node, err error) {
	jdFilter := filter.bootstrappersFilter()
	resp, err := jd.ListNodes(ctx, &nodeapiv1.ListNodesRequest{Filter: jdFilter})
	if err != nil {
		return nil, fmt.Errorf("failed to list bootstrap nodes for DON %d - %s: %w", filter.DONID, filter.DONName, err)
	}

	if len(resp.Nodes) != filter.Size {
		return nil, fmt.Errorf("expected %d bootstrap nodes for DON(%d,%s), got %d", filter.Size, filter.DONID, filter.DONName, len(resp.Nodes))
	}

	return resp.Nodes, nil
}
