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

	nodeChainConfigs, err := oc.ListNodeChainConfigs(context.Background(), &nodev1.ListNodeChainConfigsRequest{Filter: &nodev1.ListNodeChainConfigsRequest_Filter{
		NodeIds: nodeIds,
	}})
	if err != nil {
		return nil, err
	}
	if len(nodeChainConfigs.ChainConfigs) != len(nodeIds) {
		return nil, fmt.Errorf("expected %d chain configs, got %d", len(nodeIds), len(nodeChainConfigs.ChainConfigs))
	}
	var p2pV2Bootstrappers []string
	// Find the bootstrap nodes
	for _, chainConfig := range nodeChainConfigs.ChainConfigs {
		if chainConfig.Ocr2Config.IsBootstrap {
			p2pV2Bootstrappers = append(p2pV2Bootstrappers, fmt.Sprintf("%s@%s",
				// p2p_12D3... -> 12D3...
				chainConfig.Ocr2Config.P2PKeyBundle.PeerId[4:], chainConfig.Ocr2Config.Multiaddr))
		}
	}

	nodesToJobSpecs := make(map[string][]string)

	for i, chainConfig := range nodeChainConfigs.ChainConfigs {
		// only set P2PV2Bootstrappers in the job spec if the node is a plugin node.
		if chainConfig.Ocr2Config.IsBootstrap {
			continue
		}
		spec, err := validate.NewCCIPSpecToml(validate.SpecArgs{
			P2PV2Bootstrappers:     p2pV2Bootstrappers,
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
		nodesToJobSpecs[nodeIds[i]] = append(nodesToJobSpecs[nodeIds[i]], spec)
	}
	return nodesToJobSpecs, nil
}
