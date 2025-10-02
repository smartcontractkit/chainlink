package offchain

import (
	"context"
	"fmt"
	"slices"
	"strings"

	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	nodeapiv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
)

// labels used in JD to identify nodes and jobs
const (
	ProductLabel              = "cre"
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

	return resp.Nodes, nil
}
