package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

var _ deployment.ChangeSet[*UpdateDonRequest] = UpdateDon

// CapabilityConfig is a struct that holds a capability and its configuration
type CapabilityConfig = internal.CapabilityConfig

type UpdateDonRequest = internal.UpdateDonRequest

type UpdateDonRequest2 struct {
	RegistryChainSel  uint64
	P2PIDs            []p2pkey.PeerID    // this is the unique identifier for the don
	CapabilityConfigs []CapabilityConfig // if Config subfield is nil, a default config is used
	UseMCMS           bool

	HackUseMulitProposal bool
}
type UpdateDonResponse struct {
	DonInfo kcr.CapabilitiesRegistryDONInfo
}

// UpdateDon updates the capabilities of a Don
// This a complex action in practice that involves registering missing capabilities, adding the nodes, and updating
// the capabilities of the DON
func UpdateDon(env deployment.Environment, req *UpdateDonRequest) (deployment.ChangesetOutput, error) {
	_, err := internal.UpdateDon(env.Logger, req)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to update don: %w", err)
	}

	return deployment.ChangesetOutput{}, nil
}

var (
	HACK_USE_MULTI_PROPOSAL = true
)

// UpdateDon updates the capabilities of a Don
// This a complex action in practice that involves registering missing capabilities, adding the nodes, and updating
// the capabilities of the DON
func UpdateDon2(env deployment.Environment, req *UpdateDonRequest2) (deployment.ChangesetOutput, error) {
	appendResult, err := AppendNodeCapabilities(env, appendRequest(req))
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to append node capabilities: %w", err)
	}

	ur, err := updateDonRequest(env, req)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create update don request: %w", err)
	}
	updateResult, err := internal.UpdateDon2(env.Logger, ur)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to update don: %w", err)
	}

	cresp, err := kslib.GetContractSets(env.Logger, &kslib.GetContractSetsRequest{
		Chains:      env.Chains,
		AddressBook: env.ExistingAddresses,
	})
	contracts := cresp.ContractSets[req.RegistryChainSel]
	out := deployment.ChangesetOutput{}
	if req.UseMCMS {
		if updateResult.Ops == nil {
			return out, fmt.Errorf("expected MCMS operation to be non-nil")
		}
		if len(appendResult.Proposals) == 0 {
			return out, fmt.Errorf("expected append node capabilities to return proposals")
		}

		out.Proposals = appendResult.Proposals
		timelocksPerChain := map[uint64]common.Address{
			req.RegistryChainSel: contracts.Timelock.Address(),
		}
		proposerMCMSes := map[uint64]*gethwrappers.ManyChainMultiSig{
			req.RegistryChainSel: contracts.ProposerMcm,
		}

		proposal, err := proposalutils.BuildProposalFromBatches(
			timelocksPerChain,
			proposerMCMSes,
			[]timelock.BatchChainOperation{*updateResult.Ops},
			"proposal to set update node capabilities",
			0,
		)
		if err != nil {
			return out, fmt.Errorf("failed to build proposal: %w", err)
		}
		if HACK_USE_MULTI_PROPOSAL {
			out.Proposals = append(out.Proposals, *proposal)
		} else {
			out.Proposals[0].Transactions = append(out.Proposals[0].Transactions, proposal.Transactions...)
		}
		//out.Proposals = append(out.Proposals, *proposal)
	}
	return out, nil

}

func appendRequest(r *UpdateDonRequest2) *AppendNodeCapabilitiesRequest {
	out := &AppendNodeCapabilitiesRequest{
		RegistryChainSel:  r.RegistryChainSel,
		P2pToCapabilities: make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability),
		UseMCMS:           r.UseMCMS,
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

func updateDonRequest(env deployment.Environment, r *UpdateDonRequest2) (*UpdateDonRequest, error) {
	resp, err := kslib.GetContractSets(env.Logger, &kslib.GetContractSetsRequest{
		Chains:      env.Chains,
		AddressBook: env.ExistingAddresses,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get contract sets: %w", err)
	}
	contractSet := resp.ContractSets[r.RegistryChainSel]

	return &UpdateDonRequest{
		Chain:             env.Chains[r.RegistryChainSel],
		ContractSet:       &contractSet,
		P2PIDs:            r.P2PIDs,
		CapabilityConfigs: r.CapabilityConfigs,
		UseMCMS:           r.UseMCMS,
	}, nil
}
