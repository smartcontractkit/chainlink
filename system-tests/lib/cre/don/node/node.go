package node

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

const (
	NodeTypeKey            = "type"
	HostLabelKey           = "host"
	IndexKey               = "node_index"
	ExtraRolesKey          = "extra_roles"
	NodeIDKey              = "node_id"
	NodeOCR2KeyBundleIDKey = "ocr2_key_bundle_id"
	NodeP2PIDKey           = "p2p_id"
)

func AddressKeyFromSelector(chainSelector uint64) string {
	return strconv.FormatUint(chainSelector, 10) + "_public_address"
}

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

func ToP2PID(node *types.NodeMetadata, transformFn stringTransformer) (string, error) {
	for _, label := range node.Labels {
		if label.Key == NodeP2PIDKey {
			if label.Value == "" {
				return "", errors.New("p2p label value is empty for node")
			}
			return transformFn(label.Value), nil
		}
	}

	return "", errors.New("p2p label not found for node")
}

// copied from Bala's unmerged PR: https://github.com/smartcontractkit/chainlink/pull/15751
// TODO: remove this once the PR is merged and import his function
// IMPORTANT ADDITION: prefix to differentiate between the different DONs
func GetNodeInfo(nodeOut *ns.Output, prefix string, bootstrapNodeCount int) ([]devenv.NodeInfo, error) {
	var nodeInfo []devenv.NodeInfo
	for i := 1; i <= len(nodeOut.CLNodes); i++ {
		p2pURL, err := url.Parse(nodeOut.CLNodes[i-1].Node.InternalP2PUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to parse p2p url: %w", err)
		}
		if i <= bootstrapNodeCount {
			nodeInfo = append(nodeInfo, devenv.NodeInfo{
				IsBootstrap: true,
				Name:        fmt.Sprintf("%s_bootstrap-%d", prefix, i),
				P2PPort:     p2pURL.Port(),
				CLConfig: nodeclient.ChainlinkConfig{
					URL:        nodeOut.CLNodes[i-1].Node.ExternalURL,
					Email:      nodeOut.CLNodes[i-1].Node.APIAuthUser,
					Password:   nodeOut.CLNodes[i-1].Node.APIAuthPassword,
					InternalIP: nodeOut.CLNodes[i-1].Node.InternalIP,
				},
				Labels: map[string]string{
					NodeTypeKey: types.BootstrapNode,
				},
			})
		} else {
			nodeInfo = append(nodeInfo, devenv.NodeInfo{
				IsBootstrap: false,
				Name:        fmt.Sprintf("%s_node-%d", prefix, i),
				P2PPort:     p2pURL.Port(),
				CLConfig: nodeclient.ChainlinkConfig{
					URL:        nodeOut.CLNodes[i-1].Node.ExternalURL,
					Email:      nodeOut.CLNodes[i-1].Node.APIAuthUser,
					Password:   nodeOut.CLNodes[i-1].Node.APIAuthPassword,
					InternalIP: nodeOut.CLNodes[i-1].Node.InternalIP,
				},
				Labels: map[string]string{
					NodeTypeKey: types.WorkerNode,
				},
			})
		}
	}
	return nodeInfo, nil
}

func FindOneWithLabel(nodes []*types.NodeMetadata, wantedLabel *types.Label, labelMatcherFn labelMatcherFn) (*types.NodeMetadata, error) {
	if wantedLabel == nil {
		return nil, errors.New("label is nil")
	}
	for _, node := range nodes {
		for _, label := range node.Labels {
			if wantedLabel.Key == label.Key && labelMatcherFn(wantedLabel.Value, label.Value) {
				return node, nil
			}
		}
	}
	return nil, fmt.Errorf("node with label %s=%s not found", wantedLabel.Key, wantedLabel.Value)
}

func FindManyWithLabel(nodes []*types.NodeMetadata, wantedLabel *types.Label, labelMatcherFn labelMatcherFn) ([]*types.NodeMetadata, error) {
	if wantedLabel == nil {
		return nil, errors.New("label is nil")
	}

	var foundNodes []*types.NodeMetadata

	for _, node := range nodes {
		for _, label := range node.Labels {
			if wantedLabel.Key == label.Key && labelMatcherFn(wantedLabel.Value, label.Value) {
				foundNodes = append(foundNodes, node)
			}
		}
	}

	return foundNodes, nil
}

func FindLabelValue(node *types.NodeMetadata, labelKey string) (string, error) {
	for _, label := range node.Labels {
		if label.Key == labelKey {
			if label.Value == "" {
				return "", fmt.Errorf("label %s found, but its value is empty", labelKey)
			}
			return label.Value, nil
		}
	}

	return "", fmt.Errorf("label %s not found", labelKey)
}

type labelMatcherFn func(first, second string) bool

func EqualLabels(first, second string) bool {
	return first == second
}

func LabelContains(first, second string) bool {
	return strings.Contains(first, second)
}
