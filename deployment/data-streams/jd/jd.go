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
type DONFilter struct {
	DONID    uint64
	DONName  string
	EnvLabel string
}

func (f *DONFilter) bootstrappersFilter() *nodeapiv1.ListNodesRequest_Filter {
	return &nodeapiv1.ListNodesRequest_Filter{
		Selectors: []*jdtypesv1.Selector{
			{
				Key: fmt.Sprintf("don-%s-%s", f.DONID, f.DONName),
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

func FetchDONBootstrappersFromJD(ctx context.Context, jd deployment.OffchainClient, filter *DONFilter) (nodes []*nodeapiv1.Node, err error) {
	jdFilter := filter.bootstrappersFilter()
	resp, err := jd.ListNodes(ctx, &nodeapiv1.ListNodesRequest{Filter: jdFilter})
	if err != nil {
		return nil, fmt.Errorf("failed to list bootstrap nodes for DON %s - %s: %w", filter.DONID, filter.DONName, err)
	}

	return resp.Nodes, nil
}
