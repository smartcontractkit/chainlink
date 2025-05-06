package v1_6

import (
	"bytes"
	"fmt"

	"github.com/pelletier/go-toml/v2"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip"
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
	existingSpecs, err := acceptedOrPendingAcceptedJobSpecs(env, nodes)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	// We first generate the job specs for the CCIP capability. so that if
	// there are any errors in the job specs, we can fail early.
	// Generate a set of new job specs for CCIP for a specific environment
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
				CapabilityVersion:      ccip.CapabilityVersion,
				CapabilityLabelledName: ccip.CapabilityLabelledName,
				OCRKeyBundleIDs:        keyBundles,
				P2PKeyID:               node.PeerID.String(),
				RelayConfigs:           nil,
				PluginConfig:           map[string]any{},
			})
		} else {
			spec, err = validate.NewCCIPSpecToml(validate.SpecArgs{
				P2PV2Bootstrappers:     []string{}, // Intentionally empty for bootstraps.
				CapabilityVersion:      ccip.CapabilityVersion,
				CapabilityLabelledName: ccip.CapabilityLabelledName,
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
	// Now we propose the job specs to the offchain system.
	var Jobs []deployment.ProposedJob
	for nodeID, jobs := range nodesToJobSpecs {
		nodeID := nodeID
		for _, job := range jobs {
			Jobs = append(Jobs, deployment.ProposedJob{
				Node: nodeID,
				Spec: job,
			})
			if !isJobProposed(env, nodeID, job) {
				// we add a label to the job spec so that we can identify it later
				labels := []*ptypes.Label{
					{
						Key:   "node_id",
						Value: &nodeID,
					},
				}
				res, err := env.Offchain.ProposeJob(env.GetContext(),
					&jobv1.ProposeJobRequest{
						NodeId: nodeID,
						Spec:   job,
						Labels: labels,
					})
				if err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("failed to propose job: %w", err)
				}
				Jobs[len(Jobs)-1].JobID = res.Proposal.JobId
			}
		}
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: nil,
		Jobs:        Jobs,
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

// acceptedOrPendingAcceptedJobSpecs returns a map of nodeID to job specs that are either accepted or pending review
// or proposed
func acceptedOrPendingAcceptedJobSpecs(env deployment.Environment, nodes deployment.Nodes) (map[string][]string, error) {
	existingSpecs := make(map[string][]string)
	for _, node := range nodes {
		jobs, err := env.Offchain.ListJobs(env.GetContext(), &jobv1.ListJobsRequest{
			Filter: &jobv1.ListJobsRequest_Filter{
				NodeIds: []string{node.NodeID},
			},
		})
		if err != nil {
			return make(map[string][]string), fmt.Errorf("failed to list jobs for node %s: %w", node.NodeID, err)
		}
		for _, j := range jobs.Jobs {
			// skip deleted jobs
			if j.DeletedAt != nil {
				continue
			}
			for _, propID := range j.ProposalIds {
				jbProposal, err := env.Offchain.GetProposal(env.GetContext(), &jobv1.GetProposalRequest{
					Id: propID,
				})
				if err != nil {
					return nil, fmt.Errorf("failed to get job proposal %s on node %s: %w", propID, node.NodeID, err)
				}
				if jbProposal.Proposal == nil {
					return nil, fmt.Errorf("job proposal %s on node %s is nil", propID, node.NodeID)
				}
				if jbProposal.Proposal.Status == jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED ||
					jbProposal.Proposal.Status == jobv1.ProposalStatus_PROPOSAL_STATUS_PENDING ||
					jbProposal.Proposal.Status == jobv1.ProposalStatus_PROPOSAL_STATUS_PROPOSED {
					existingSpecs[node.NodeID] = append(existingSpecs[node.NodeID], jbProposal.Proposal.Spec)
				}
			}
		}
	}
	return existingSpecs, nil
}

// isJobProposed returns true if the job spec is already proposed to a node
// currently there is no way to check if a job spec is already proposed to a specific node
// so we return false TODO implement this check when it is available
// there is no downside to proposing the same job spec to a node multiple times
// It will only be accepted once by node
func isJobProposed(env deployment.Environment, nodeID string, spec string) bool {
	// implement this when it is available
	return false
}
