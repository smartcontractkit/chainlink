package changeset

import (
	"bytes"
	"fmt"

	"github.com/pelletier/go-toml/v2"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/validate"
	corejob "github.com/smartcontractkit/chainlink/v2/core/services/job"
)

var _ deployment.ChangeSet[any] = CCIPCapabilityJobspecChangeset

// CCIPCapabilityJobspecChangeset returns the job specs for the CCIP capability.
// The caller needs to propose these job specs to the offchain system.
func CCIPCapabilityJobspecChangeset(env deployment.Environment, _ any) (deployment.ChangesetOutput, error) {
	nodes, err := deployment.NodeInfo(env.NodeIDs, env.Offchain)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	// find existing jobs
	existingSpecs := make(map[string][]string)
	for _, node := range nodes {
		jobs, err := env.Offchain.ListJobs(env.GetContext(), &job.ListJobsRequest{
			Filter: &job.ListJobsRequest_Filter{
				NodeIds: []string{node.NodeID},
			},
		})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to list jobs for node %s: %w", node.NodeID, err)
		}
		for _, j := range jobs.Jobs {
			for _, propID := range j.ProposalIds {
				jbProposal, err := env.Offchain.GetProposal(env.GetContext(), &job.GetProposalRequest{
					Id: propID,
				})
				if err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("failed to get job proposal %s on node %s: %w", propID, node.NodeID, err)
				}
				existingSpecs[node.NodeID] = append(existingSpecs[node.NodeID], jbProposal.Proposal.Spec)
			}
		}
	}
	// Generate a set of brand new job specs for CCIP for a specific environment
	// (including NOPs) and new addresses.
	// We want to assign one CCIP capability job to each node. And node with
	// an addr we'll list as bootstrapper.
	// Find the bootstrap nodes
	nodesToJobSpecs := make(map[string][]string)
	for _, node := range nodes {
		// pick first keybundle of each type (by default nodes will auto create one key of each type for each defined chain family)
		keyBundles := map[string]string{}
		for details, config := range node.SelToOCRConfig {
			family, err := chainsel.GetSelectorFamily(details.ChainSelector)
			if err != nil {
				env.Logger.Warnf("skipping unknown/invalid chain family for selector %v", details.ChainSelector)
				continue
			}
			_, exists := keyBundles[family]
			if exists {
				if keyBundles[family] != config.KeyBundleID {
					return deployment.ChangesetOutput{}, fmt.Errorf("multiple different %v OCR keys found for node %v", family, node.PeerID)
				}
				continue
			}
			keyBundles[family] = config.KeyBundleID
		}

		var spec string
		var err error
		if !node.IsBootstrap {
			spec, err = validate.NewCCIPSpecToml(validate.SpecArgs{
				P2PV2Bootstrappers:     nodes.BootstrapLocators(),
				CapabilityVersion:      internal.CapabilityVersion,
				CapabilityLabelledName: internal.CapabilityLabelledName,
				OCRKeyBundleIDs:        keyBundles,
				P2PKeyID:               node.PeerID.String(),
				RelayConfigs:           nil,
				PluginConfig:           map[string]any{},
			})
		} else {
			spec, err = validate.NewCCIPSpecToml(validate.SpecArgs{
				P2PV2Bootstrappers:     []string{}, // Intentionally empty for bootstraps.
				CapabilityVersion:      internal.CapabilityVersion,
				CapabilityLabelledName: internal.CapabilityLabelledName,
				OCRKeyBundleIDs:        map[string]string{},
				P2PKeyID:               node.PeerID.String(),
				RelayConfigs:           nil,
				PluginConfig:           map[string]any{},
			})
		}
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		// If the spec already exists, don't propose it again
		specExists := false
		if existingSpecs[node.NodeID] != nil {
			for _, existingSpec := range existingSpecs[node.NodeID] {
				specExists, err = areCCIPSpecsEqual(existingSpec, spec)
				if err != nil {
					return deployment.ChangesetOutput{}, err
				}
			}
		}
		if !specExists {
			nodesToJobSpecs[node.NodeID] = append(nodesToJobSpecs[node.NodeID], spec)
		}
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: nil,
		JobSpecs:    nodesToJobSpecs,
		Jobs:        nil,
	}, nil
}

func areCCIPSpecsEqual(existingSpecStr, newSpecStr string) (bool, error) {
	var existingCCIPSpec, newSpec corejob.CCIPSpec
	err := toml.Unmarshal([]byte(existingSpecStr), &existingCCIPSpec)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal existing job spec: %w", err)
	}
	err = toml.Unmarshal([]byte(newSpecStr), &newSpec)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal new job spec: %w", err)
	}
	existingOCRKey, err := existingCCIPSpec.OCRKeyBundleIDs.Value()
	if err != nil {
		return false, fmt.Errorf("failed to get OCRKeyBundleIDs from existing job spec: %w", err)
	}

	newOCRKey, err := newSpec.OCRKeyBundleIDs.Value()
	if err != nil {
		return false, fmt.Errorf("failed to get OCRKeyBundleIDs from new job spec: %w", err)
	}
	p2pBootstrapperValue, err := existingCCIPSpec.P2PV2Bootstrappers.Value()
	if err != nil {
		return false, fmt.Errorf("failed to get P2PV2Bootstrappers from existing job spec: %w", err)
	}
	pluginConfigValue, err := existingCCIPSpec.PluginConfig.Value()
	if err != nil {
		return false, fmt.Errorf("failed to get PluginConfig from existing job spec: %w", err)
	}
	relayConfigValue, err := existingCCIPSpec.RelayConfigs.Value()
	if err != nil {
		return false, fmt.Errorf("failed to get RelayConfigs from existing job spec: %w", err)
	}
	p2pBootstrapperValueNew, err := newSpec.P2PV2Bootstrappers.Value()
	if err != nil {
		return false, fmt.Errorf("failed to get P2PV2Bootstrappers from new job spec: %w", err)
	}
	pluginConfigValueNew, err := newSpec.PluginConfig.Value()
	if err != nil {
		return false, fmt.Errorf("failed to get PluginConfig from new job spec: %w", err)
	}
	relayConfigValueNew, err := newSpec.RelayConfigs.Value()
	if err != nil {
		return false, fmt.Errorf("failed to get RelayConfigs from new job spec: %w", err)
	}

	return existingCCIPSpec.CapabilityLabelledName == newSpec.CapabilityLabelledName &&
		existingCCIPSpec.CapabilityVersion == newSpec.CapabilityVersion &&
		bytes.Equal(existingOCRKey.([]byte), newOCRKey.([]byte)) &&
		existingCCIPSpec.P2PKeyID == newSpec.P2PKeyID &&
		p2pBootstrapperValue == p2pBootstrapperValueNew &&
		bytes.Equal(pluginConfigValue.([]byte), pluginConfigValueNew.([]byte)) &&
		bytes.Equal(relayConfigValue.([]byte), relayConfigValueNew.([]byte)), nil
}
