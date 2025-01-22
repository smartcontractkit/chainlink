package changeset

import (
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

var _ deployment.ChangeSet[any] = CCIPCapabilityJobspecChangeset

// CCIPCapabilityJobspecChangeset returns the job specs for the CCIP capability.
// The caller needs to propose these job specs to the offchain system.
func CCIPCapabilityJobspecChangeset(env deployment.Environment, _ any) (deployment.ChangesetOutput, error) {
	nodes, err := deployment.NodeInfo(env.NodeIDs, env.Offchain)
	if err != nil {
		return deployment.ChangesetOutput{}, err
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
				CapabilityVersion:      internal.CapabilityVersion,
				CapabilityLabelledName: internal.CapabilityLabelledName,
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
				CapabilityVersion:      internal.CapabilityVersion,
				CapabilityLabelledName: internal.CapabilityLabelledName,
				OCRKeyBundleIDs:        map[string]string{},
				// TODO: validate that all EVM chains are using the same keybundle
				P2PKeyID:     node.PeerID.String(),
				RelayConfigs: nil,
				PluginConfig: map[string]any{},
			})
		}
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		nodesToJobSpecs[node.NodeID] = append(nodesToJobSpecs[node.NodeID], spec)
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: nil,
		JobSpecs:    nodesToJobSpecs,
	}, nil
}
