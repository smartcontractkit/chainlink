package ccipdeployment

import (
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// In our case, the only address needed is the cap registry which is actually an env var.
// and will pre-exist for our deployment. So the job specs only depend on the environment operators.
func NewCCIPJobSpecs(nodeIds []string, oc deployment.OffchainClient) (map[string][]string, error) {
	nodes, err := deployment.NodeInfo(nodeIds, oc)
	if err != nil {
		return nil, err
	}
	// Generate a set of brand new job specs for CCIP for a specific environment
	// (including NOPs) and new addresses.
	// We want to assign one CCIP capability job to each node. And node with
	// an addr we'll list as bootstrapper.
	// Find the bootstrap nodes

	nodesToJobSpecs := make(map[string][]string)
	for _, node := range nodes {
		var spec string
		var err error
		if !node.IsBootstrap {
			spec, err = validate.NewCCIPSpecToml(validate.SpecArgs{
				P2PV2Bootstrappers:     nodes.BootstrapLocators(),
				CapabilityVersion:      CapabilityVersion,
				CapabilityLabelledName: CapabilityLabelledName,
				OCRKeyBundleIDs: map[string]string{
					// TODO: Validate that that all EVM chains are using the same keybundle.
					relay.NetworkEVM: node.FirstOCRKeybundle().KeyBundleID,
				},
				P2PKeyID:     node.PeerID.String(),
				RelayConfigs: nil,
				PluginConfig: map[string]any{},
			})
		} else {
			spec, err = validate.NewCCIPSpecToml(validate.SpecArgs{
				P2PV2Bootstrappers:     []string{}, // Intentionally empty for bootstraps.
				CapabilityVersion:      CapabilityVersion,
				CapabilityLabelledName: CapabilityLabelledName,
				OCRKeyBundleIDs:        map[string]string{},
				// TODO: validate that all EVM chains are using the same keybundle
				P2PKeyID:     node.PeerID.String(),
				RelayConfigs: nil,
				PluginConfig: map[string]any{},
			})
		}
		if err != nil {
			return nil, err
		}
		nodesToJobSpecs[node.NodeID] = append(nodesToJobSpecs[node.NodeID], spec)
	}
	return nodesToJobSpecs, nil
}
