package node

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type stringTransformer func(string) string

func NoOpTransformFn(value string) string {
	return value
}

func KeyExtractingTransformFn(value string) string {
	parts := strings.Split(value, "_")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return value
}

func ToP2PID(node devenv.Node, transformFn stringTransformer) (string, error) {
	for _, label := range node.Labels() {
		if label.Key == devenv.NodeLabelP2PIDType {
			if label.Value == nil {
				return "", fmt.Errorf("p2p label value is nil for node %s", node.Name)
			}
			return transformFn(*label.Value), nil
		}
	}

	return "", fmt.Errorf("p2p label not found for node %s", node.Name)
}

const (
	RoleLabelKey = "role"
	HostLabelKey = "host"
	NodeIndexKey = "node_index"
)

// copied from Bala's unmerged PR: https://github.com/smartcontractkit/chainlink/pull/15751
// TODO: remove this once the PR is merged and import his function
// IMPORTANT ADDITION:  prefix to differentiate between the different DONs
func GetNodeInfo(nodeOut *ns.Output, prefix string, bootstrapNodeCount int) ([]devenv.NodeInfo, error) {
	var nodeInfo []devenv.NodeInfo
	for i := 1; i <= len(nodeOut.CLNodes); i++ {
		p2pURL, err := url.Parse(nodeOut.CLNodes[i-1].Node.DockerP2PUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to parse p2p url: %w", err)
		}
		if i <= bootstrapNodeCount {
			nodeInfo = append(nodeInfo, devenv.NodeInfo{
				IsBootstrap: true,
				Name:        fmt.Sprintf("%s_bootstrap-%d", prefix, i),
				P2PPort:     p2pURL.Port(),
				CLConfig: nodeclient.ChainlinkConfig{
					URL:        nodeOut.CLNodes[i-1].Node.HostURL,
					Email:      nodeOut.CLNodes[i-1].Node.APIAuthUser,
					Password:   nodeOut.CLNodes[i-1].Node.APIAuthPassword,
					InternalIP: nodeOut.CLNodes[i-1].Node.InternalIP,
				},
				Labels: map[string]string{
					HostLabelKey: nodeOut.CLNodes[i-1].Node.ContainerName,
					NodeIndexKey: strconv.Itoa(i - 1),
					RoleLabelKey: types.BootstrapNode,
				},
			})
		} else {
			nodeInfo = append(nodeInfo, devenv.NodeInfo{
				IsBootstrap: false,
				Name:        fmt.Sprintf("%s_node-%d", prefix, i),
				P2PPort:     p2pURL.Port(),
				CLConfig: nodeclient.ChainlinkConfig{
					URL:        nodeOut.CLNodes[i-1].Node.HostURL,
					Email:      nodeOut.CLNodes[i-1].Node.APIAuthUser,
					Password:   nodeOut.CLNodes[i-1].Node.APIAuthPassword,
					InternalIP: nodeOut.CLNodes[i-1].Node.InternalIP,
				},
				Labels: map[string]string{
					HostLabelKey: nodeOut.CLNodes[i-1].Node.ContainerName,
					NodeIndexKey: strconv.Itoa(i - 1),
					RoleLabelKey: types.WorkerNode,
				},
			})
		}
	}
	return nodeInfo, nil
}

func FindOneWithLabel(nodes *devenv.DON, wantedLabel *ptypes.Label) (*devenv.Node, error) {
	if wantedLabel == nil {
		return nil, errors.New("label is nil")
	}
	for _, node := range nodes.Nodes {
		for _, label := range node.Labels() {
			if wantedLabel.Key == label.Key && equalLabels(wantedLabel.Value, label.Value) {
				return &node, nil
			}
		}
	}
	return nil, fmt.Errorf("node with label %s=%s not found", wantedLabel.Key, *wantedLabel.Value)
}

func FindManyWithLabel(nodes *devenv.DON, wantedLabel *ptypes.Label) ([]*devenv.Node, error) {
	if wantedLabel == nil {
		return nil, errors.New("label is nil")
	}

	var foundNodes []*devenv.Node

	for _, node := range nodes.Nodes {
		for _, label := range node.Labels() {
			if wantedLabel.Key == label.Key && equalLabels(wantedLabel.Value, label.Value) {
				foundNodes = append(foundNodes, &node)
			}
		}
	}

	if len(foundNodes) == 0 {
		return nil, fmt.Errorf("node with label %s=%s not found", wantedLabel.Key, *wantedLabel.Value)
	}

	return foundNodes, nil
}

func equalLabels(first, second *string) bool {
	if first == nil && second == nil {
		return true
	}
	if first == nil || second == nil {
		return false
	}
	return *first == *second
}
