package offchain

import (
	"context"
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

func FetchNodesFromJD(ctx context.Context, jd cldf_offchain.Client, filter *nodeapiv1.ListNodesRequest_Filter) (nodes []*nodeapiv1.Node, err error) {
	resp, err := jd.ListNodes(ctx, &nodeapiv1.ListNodesRequest{Filter: filter})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
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
