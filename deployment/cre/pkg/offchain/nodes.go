package offchain

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strings"

	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	nodeapiv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/helpers/pointer"
)

// labels used in JD to identify nodes and jobs
const (
	ProductLabel              = "cre"
	WorkflowOwnerAddressLabel = "workflow_owner"
	WorkflowNameLabel         = "workflow_name"
	GatewayNameLabel          = "gateway_name"
	CapabilityLabel           = "capability_name"
)

// Matches the labels at
// https://github.com/smartcontractkit/chainlink/blob/b7f9f23f3aeae5d0cfd003c57bd1d9d19e2ddb80/deployment/environment/devenv/don.go#L38-L38
const (
	labelNodeTypeKey            = "type"
	labelNodeTypeValueBootstrap = "bootstrap"
	labelNodeTypeValuePlugin    = "plugin"
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

// RegisterNode registers a single node with the job distributor. It errors if the node is already registered.
func RegisterNode(
	ctx context.Context,
	jd cldf_offchain.Client,
	name string,
	csaKey string,
	isBootstrap bool,
	domain string, // domain key
	environment string,
	extraLabels map[string]string,
) (string, error) {
	labels := make([]*ptypes.Label, 0)
	labels = append(labels, &ptypes.Label{
		Key:   "product",
		Value: &domain,
	})
	labels = append(labels, &ptypes.Label{
		Key:   "environment",
		Value: &environment,
	})

	// Sort extraLabels keys to ensure deterministic label ordering
	extraLabelKeys := make([]string, 0, len(extraLabels))
	for key := range extraLabels {
		extraLabelKeys = append(extraLabelKeys, key)
	}
	sort.Strings(extraLabelKeys)

	for _, key := range extraLabelKeys {
		labels = append(labels, &ptypes.Label{
			Key:   key,
			Value: pointer.To(extraLabels[key]),
		})
	}
	if isBootstrap {
		labels = append(labels, &ptypes.Label{
			Key:   labelNodeTypeKey,
			Value: pointer.To(labelNodeTypeValueBootstrap),
		})
	} else {
		labels = append(labels, &ptypes.Label{
			Key:   labelNodeTypeKey,
			Value: pointer.To(labelNodeTypeValuePlugin),
		})
	}
	resp, err := jd.RegisterNode(ctx, &nodeapiv1.RegisterNodeRequest{
		Name:      name,
		PublicKey: csaKey,
		Labels:    labels,
	})
	if err != nil {
		return "", fmt.Errorf("failed to register node %s : %w", name, err)
	}
	if resp == nil || resp.Node == nil || resp.Node.Id == "" {
		return "", fmt.Errorf("failed to register node %s, blank response received", name)
	}

	return resp.Node.Id, nil
}
