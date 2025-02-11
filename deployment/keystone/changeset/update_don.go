package changeset

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

var _ deployment.ChangeSet[*UpdateDonRequest] = UpdateDon

// CapabilityConfig is a struct that holds a capability and its configuration
type CapabilityConfig = internal.CapabilityConfig

type UpdateDonRequest struct {
	RegistryChainSel  uint64
	P2PIDs            []p2pkey.PeerID    // this is the unique identifier for the don
	CapabilityConfigs []CapabilityConfig // if Config subfield is nil, a default config is used

	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	MCMSConfig *MCMSConfig
}

func (r *UpdateDonRequest) Validate() error {
	if len(r.P2PIDs) == 0 {
		return errors.New("p2pIDs is required")
	}
	if len(r.CapabilityConfigs) == 0 {
		return errors.New("capabilityConfigs is required")
	}
	return nil
}

func (r UpdateDonRequest) UseMCMS() bool {
	return r.MCMSConfig != nil
}

type UpdateDonResponse struct {
	DonInfo kcr.CapabilitiesRegistryDONInfo
}

// UpdateDon updates the capabilities of a Don
// This a complex action in practice that involves registering missing capabilities, adding the nodes, and updating
// the capabilities of the DON
func UpdateDon(env deployment.Environment, req *UpdateDonRequest) (deployment.ChangesetOutput, error) {
	appendResult, err := AppendNodeCapabilities(env, appendRequest(req))
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to append node capabilities: %w", err)
	}

	ur, err := updateDonRequest(env, req)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create update don request: %w", err)
	}
	updateResult, err := internal.UpdateDon(env.Logger, ur)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to update don: %w", err)
	}

	out := deployment.ChangesetOutput{}
	if req.UseMCMS() {
		if updateResult.Ops == nil {
			return out, errors.New("expected MCMS operation to be non-nil")
		}
		if len(appendResult.Proposals) == 0 {
			return out, errors.New("expected append node capabilities to return proposals")
		}

		out.Proposals = appendResult.Proposals

		// add the update don to the existing batch
		// this makes the proposal all-or-nothing because all the operations are in the same batch, there is only one tr
		// transaction and only one proposal
		out.Proposals[0].Transactions[0].Batch = append(out.Proposals[0].Transactions[0].Batch, updateResult.Ops.Batch...)
	}
	return out, nil
}

func appendRequest(r *UpdateDonRequest) *AppendNodeCapabilitiesRequest {
	out := &AppendNodeCapabilitiesRequest{
		RegistryChainSel:  r.RegistryChainSel,
		P2pToCapabilities: make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability),
		MCMSConfig:        r.MCMSConfig,
	}
	for _, p2pid := range r.P2PIDs {
		if _, exists := out.P2pToCapabilities[p2pid]; !exists {
			out.P2pToCapabilities[p2pid] = make([]kcr.CapabilitiesRegistryCapability, 0)
		}
		for _, cc := range r.CapabilityConfigs {
			out.P2pToCapabilities[p2pid] = append(out.P2pToCapabilities[p2pid], cc.Capability)
		}
	}
	return out
}

func updateDonRequest(env deployment.Environment, r *UpdateDonRequest) (*internal.UpdateDonRequest, error) {
	resp, err := internal.GetContractSets(env.Logger, &internal.GetContractSetsRequest{
		Chains:      env.Chains,
		AddressBook: env.ExistingAddresses,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get contract sets: %w", err)
	}
	contractSet := resp.ContractSets[r.RegistryChainSel]

	return &internal.UpdateDonRequest{
		Chain:                env.Chains[r.RegistryChainSel],
		CapabilitiesRegistry: contractSet.CapabilitiesRegistry,
		P2PIDs:               r.P2PIDs,
		CapabilityConfigs:    r.CapabilityConfigs,
		UseMCMS:              r.UseMCMS(),
	}, nil
}
