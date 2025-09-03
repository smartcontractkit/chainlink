package offchain

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	nodeapiv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	jdtypesv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
)

// labels used in JD to identify nodes and jobs
const (
	ProductLabel              = "cre"
	P2pIDLabel                = "p2p_id"
	WorkflowOwnerAddressLabel = "workflow_owner"
	WorkflowNameLabel         = "workflow_name"
	GatewayNameLabel          = "gateway_name"
	CapabilityLabel           = "capability_name"
)

// Source for a set of JD filters to apply when fetching a DON.
// Should uniquely identify a set of nodes in JD which belong to a single DON.
type DONFilter struct {
	DONName      string            `json:"don_name" yaml:"don_name"`
	EnvLabel     string            `json:"env_label" yaml:"env_label"`
	ProductLabel string            `json:"product_label" yaml:"product_label"`
	ExtraLabels  map[string]string `json:"extra_labels" yaml:"extra_labels"`
	Size         int               `json:"size" yaml:"size"`
}

func (f *DONFilter) filter() *nodeapiv1.ListNodesRequest_Filter {
	filters := &nodeapiv1.ListNodesRequest_Filter{
		Selectors: []*jdtypesv1.Selector{
			{
				Key: "don-" + f.DONName,
				Op:  jdtypesv1.SelectorOp_EXIST,
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
		},
	}
	for k, v := range f.ExtraLabels {
		filters.Selectors = append(filters.Selectors, &jdtypesv1.Selector{
			Key:   k,
			Op:    jdtypesv1.SelectorOp_EQ,
			Value: &v,
		})
	}

	return filters
}

func FetchNodesFromJD(ctx context.Context, jd cldf_offchain.Client, filter *DONFilter) (nodes []*nodeapiv1.Node, err error) {
	jdFilter := filter.filter()
	resp, err := jd.ListNodes(ctx, &nodeapiv1.ListNodesRequest{Filter: jdFilter})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	if len(resp.Nodes) != filter.Size {
		b, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, fmt.Errorf("expected %d nodes, got %d in %v", filter.Size, len(resp.Nodes), string(b))
	}
	slices.SortFunc(resp.Nodes, func(a, b *nodeapiv1.Node) int {
		return strings.Compare(a.Name, b.Name)
	})
	for _, node := range resp.Nodes {
		if GetP2pLabel(node.GetLabels()) == "" {
			return nil, fmt.Errorf("node %s has no non-empty p2p_id label: %v", node.Name, node)
		}
	}

	return resp.Nodes, nil
}

func GetP2pLabel(labels []*jdtypesv1.Label) string {
	for _, label := range labels {
		if label.GetKey() == P2pIDLabel {
			return label.GetValue()
		}
	}

	return ""
}
