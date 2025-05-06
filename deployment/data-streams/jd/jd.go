package jd

import (
	"context"
	"fmt"
	"slices"

	nodeapiv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	jdtypesv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/pointer"
)

// Source for a set of JD filters to apply when fetching a DON.
// Should uniquely identify a set of nodes in JD which belong to a single DON.
type ListFilter struct {
	DONID             uint64
	DONName           string
	EnvLabel          string
	NumOracleNodes    int // Expected number of oracle nodes in the result
	NumBootstrapNodes int // Expected number of bootstrap nodes in the result
}

func (f *ListFilter) bootstrappersFilter() *nodeapiv1.ListNodesRequest_Filter {
	return &nodeapiv1.ListNodesRequest_Filter{
		Selectors: []*jdtypesv1.Selector{
			{
				Key: utils.DonIdentifier(f.DONID, f.DONName),
				Op:  jdtypesv1.SelectorOp_EXIST,
			},
			{
				Key:   devenv.LabelNodeTypeKey,
				Op:    jdtypesv1.SelectorOp_EQ,
				Value: pointer.To(devenv.LabelNodeTypeValueBootstrap),
			},
			{
				Key:   devenv.LabelEnvironmentKey,
				Op:    jdtypesv1.SelectorOp_EQ,
				Value: &f.EnvLabel,
			},
			{
				Key:   devenv.LabelProductKey,
				Op:    jdtypesv1.SelectorOp_EQ,
				Value: pointer.To(utils.ProductLabel),
			},
		},
	}
}

// oraclesFilter is used to fetch all oracle (non-bootstrap) nodes in a DON.
func (f *ListFilter) oraclesFilter() *nodeapiv1.ListNodesRequest_Filter {
	return &nodeapiv1.ListNodesRequest_Filter{
		Selectors: []*jdtypesv1.Selector{
			{
				Key: utils.DonIdentifier(f.DONID, f.DONName),
				Op:  jdtypesv1.SelectorOp_EXIST,
			},
			{
				Key:   devenv.LabelNodeTypeKey,
				Op:    jdtypesv1.SelectorOp_EQ,
				Value: pointer.To(devenv.LabelNodeTypeValuePlugin),
			},
			{
				Key:   devenv.LabelEnvironmentKey,
				Op:    jdtypesv1.SelectorOp_EQ,
				Value: &f.EnvLabel,
			},
			{
				Key:   devenv.LabelProductKey,
				Op:    jdtypesv1.SelectorOp_EQ,
				Value: pointer.To(utils.ProductLabel),
			},
		},
	}
}

// FetchDONBootstrappersFromJD fetches all bootstrap nodes which match the given filter *and* their name is in the nodeNames list.
func FetchDONBootstrappersFromJD(ctx context.Context, jd deployment.OffchainClient, filter *ListFilter, nodeNames []string) ([]*nodeapiv1.Node, error) {
	jdFilter := filter.bootstrappersFilter()
	resp, err := jd.ListNodes(ctx, &nodeapiv1.ListNodesRequest{Filter: jdFilter})
	if err != nil {
		return nil, fmt.Errorf("failed to list bootstrap nodes for DON %d - %s: %w", filter.DONID, filter.DONName, err)
	}
	nodes := make([]*nodeapiv1.Node, 0, filter.NumBootstrapNodes)
	for _, node := range resp.Nodes {
		idx := slices.IndexFunc(nodeNames, func(name string) bool {
			return node.Name == name
		})
		if idx >= 0 {
			nodes = append(nodes, node)
		}
	}

	if len(nodes) != filter.NumBootstrapNodes {
		return nil, fmt.Errorf("expected %d bootstrap nodes for DON(%d,%s), got %d", filter.NumBootstrapNodes, filter.DONID, filter.DONName, len(nodes))
	}

	return nodes, nil
}

// FetchDONOraclesFromJD fetches all oracle nodes which match the given filter *and* their name is in the nodeNames list.
func FetchDONOraclesFromJD(ctx context.Context, jd deployment.OffchainClient, filter *ListFilter, nodeNames []string) ([]*nodeapiv1.Node, error) {
	jdFilter := filter.oraclesFilter()
	resp, err := jd.ListNodes(ctx, &nodeapiv1.ListNodesRequest{Filter: jdFilter})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes for DON %d - %s: %w", filter.DONID, filter.DONName, err)
	}
	nodes := make([]*nodeapiv1.Node, 0, filter.NumOracleNodes)
	for _, node := range resp.Nodes {
		idx := slices.IndexFunc(nodeNames, func(name string) bool {
			return node.Name == name
		})
		if idx >= 0 {
			nodes = append(nodes, node)
		}
	}

	if len(nodes) != filter.NumOracleNodes {
		return nil, fmt.Errorf("expected %d oracle nodes for DON(%d,%s), got %d", filter.NumOracleNodes, filter.DONID, filter.DONName, len(nodes))
	}

	return nodes, nil
}
