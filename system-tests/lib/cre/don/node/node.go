package node

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

var (
	NodeTypeKey            = cre.NodeTypeKey
	NodeIDKey              = cre.NodeIDKey
	NodeOCR2KeyBundleIDKey = cre.NodeOCR2KeyBundleIDKey
	NodeOCRFamiliesKey     = cre.NodeOCRFamiliesKey
	DONIDKey               = cre.DONIDKey
	EnvironmentKey         = cre.EnvironmentKey
	ProductKey             = cre.ProductKey
	DONNameKey             = cre.DONNameKey
)

// ocr2 keys depend on report's target chain family
func CreateNodeOCR2KeyBundleIDKey(chainFamily string) string {
	return NodeOCR2KeyBundleIDKey + "_" + chainFamily
}

func CreateNodeOCRFamiliesListValue(families []string) string {
	return strings.Join(families, ",")
}

func ExtractBundleKeysPerFamily(n *cre.NodeMetadata) (map[string]string, error) {
	keyBundlesFamilies, fErr := FindLabelValue(n, cre.NodeOCRFamiliesKey)
	if fErr != nil {
		return nil, fmt.Errorf("failed to get ocr families bundle id from worker node labels: %w", fErr)
	}

	supportedFamilies := strings.Split(keyBundlesFamilies, ",")

	bundlesPerFamily := make(map[string]string)
	for _, family := range supportedFamilies {
		kBundle, kbErr := FindLabelValue(n, CreateNodeOCR2KeyBundleIDKey(family))
		if kbErr != nil {
			return nil, fmt.Errorf("failed to get ocr bundle id from worker node labels for family %s err: %w", family, kbErr)
		}
		bundlesPerFamily[family] = kBundle
	}

	return bundlesPerFamily, nil
}

// copied from Bala's unmerged PR: https://github.com/smartcontractkit/chainlink/pull/15751
// TODO: remove this once the PR is merged and import his function
// IMPORTANT ADDITION: prefix to differentiate between the different DONs
func GetNodeInfo(nodeOut *ns.Output, prefix string, donID uint64, bootstrapNodeCount int) ([]devenv.NodeInfo, error) {
	var nodeInfo []devenv.NodeInfo
	for i := 1; i <= len(nodeOut.CLNodes); i++ {
		p2pURL, err := url.Parse(nodeOut.CLNodes[i-1].Node.InternalP2PUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to parse p2p url: %w", err)
		}

		info := devenv.NodeInfo{
			P2PPort: p2pURL.Port(),
			CLConfig: nodeclient.ChainlinkConfig{
				URL:        nodeOut.CLNodes[i-1].Node.ExternalURL,
				Email:      nodeOut.CLNodes[i-1].Node.APIAuthUser,
				Password:   nodeOut.CLNodes[i-1].Node.APIAuthPassword,
				InternalIP: nodeOut.CLNodes[i-1].Node.InternalIP,
			},
			Labels: map[string]string{
				"don-" + prefix: "true",
				ProductKey:      "keystone",
				EnvironmentKey:  "local",
				DONIDKey:        strconv.FormatUint(donID, 10),
				DONNameKey:      prefix,
			},
		}

		if i <= bootstrapNodeCount {
			info.IsBootstrap = true
			info.Name = fmt.Sprintf("%s_bootstrap-%d", prefix, i)
			info.Labels[NodeTypeKey] = cre.BootstrapNode
		} else {
			info.IsBootstrap = false
			info.Name = fmt.Sprintf("%s_node-%d", prefix, i)
			info.Labels[NodeTypeKey] = cre.WorkerNode
		}

		nodeInfo = append(nodeInfo, info)
	}
	return nodeInfo, nil
}

func FindOneWithLabel(nodes []*cre.NodeMetadata, wantedLabel *cre.Label, labelMatcherFn labelMatcherFn) (*cre.NodeMetadata, error) {
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

func FindManyWithLabel(nodes []*cre.NodeMetadata, wantedLabel *cre.Label, labelMatcherFn labelMatcherFn) ([]*cre.NodeMetadata, error) {
	if wantedLabel == nil {
		return nil, errors.New("label is nil")
	}

	var foundNodes []*cre.NodeMetadata

	for _, node := range nodes {
		for _, label := range node.Labels {
			if wantedLabel.Key == label.Key && labelMatcherFn(wantedLabel.Value, label.Value) {
				foundNodes = append(foundNodes, node)
			}
		}
	}

	return foundNodes, nil
}

func HasLabel(node *cre.NodeMetadata, labelKey string) bool {
	for _, label := range node.Labels {
		if label.Key == labelKey {
			return true
		}
	}
	return false
}

func FindLabelValue(node *cre.NodeMetadata, labelKey string) (string, error) {
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
