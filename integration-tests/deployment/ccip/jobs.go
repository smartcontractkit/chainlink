package ccipdeployment

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	nodev1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/node/v1"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// In our case, the only address needed is the cap registry which is actually an env var.
// and will pre-exist for our deployment. So the job specs only depend on the environment operators.
func NewCCIPJobSpecs(nodeIds []string, oc deployment.OffchainClient) (map[string][]string, error) {
	// Generate a set of brand new job specs for CCIP for a specific environment
	// (including NOPs) and new addresses.
	// We want to assign one CCIP capability job to each node. And node with
	// an addr we'll list as bootstrapper.
	// Find the bootstrap nodes
	bootstrapMp := make(map[string]struct{})
	for _, node := range nodeIds {
		// TODO: Filter should accept multiple nodes
		nodeChainConfigs, err := oc.ListNodeChainConfigs(context.Background(), &nodev1.ListNodeChainConfigsRequest{Filter: &nodev1.ListNodeChainConfigsRequest_Filter{
			NodeId: node,
		}})
		if err != nil {
			return nil, err
		}
		for _, chainConfig := range nodeChainConfigs.ChainConfigs {
			if chainConfig.Ocr2Config.IsBootstrap {
				bootstrapMp[fmt.Sprintf("%s@%s",
					// p2p_12D3... -> 12D3...
					chainConfig.Ocr2Config.P2PKeyBundle.PeerId[4:], chainConfig.Ocr2Config.Multiaddr)] = struct{}{}
			}
		}
	}
	var bootstraps []string
	for b := range bootstrapMp {
		bootstraps = append(bootstraps, b)
	}
	nodesToJobSpecs := make(map[string][]string)
	for _, node := range nodeIds {
		// TODO: Filter should accept multiple.
		nodeChainConfigs, err := oc.ListNodeChainConfigs(context.Background(), &nodev1.ListNodeChainConfigsRequest{Filter: &nodev1.ListNodeChainConfigsRequest_Filter{
			NodeId: node,
		}})
		if err != nil {
			return nil, err
		}
		spec, err := validate.NewCCIPSpecToml(validate.SpecArgs{
			P2PV2Bootstrappers:     bootstraps,
			CapabilityVersion:      CapabilityVersion,
			CapabilityLabelledName: CapabilityLabelledName,
			OCRKeyBundleIDs: map[string]string{
				// TODO: Validate that that all EVM chains are using the same keybundle.
				relay.NetworkEVM: nodeChainConfigs.ChainConfigs[0].Ocr2Config.OcrKeyBundle.BundleId,
			},
			// TODO: validate that all EVM chains are using the same keybundle
			P2PKeyID:     nodeChainConfigs.ChainConfigs[0].Ocr2Config.P2PKeyBundle.PeerId,
			RelayConfigs: nil,
			PluginConfig: map[string]any{},
		})
		if err != nil {
			return nil, err
		}
		nodesToJobSpecs[node] = append(nodesToJobSpecs[node], spec)
	}
	return nodesToJobSpecs, nil
}
