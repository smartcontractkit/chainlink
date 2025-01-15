package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ deployment.ChangeSet[*AppendNodeCapabilitiesRequest] = AppendNodeCapabilities

// AppendNodeCapabilitiesRequest is a request to add capabilities to the existing capabilities of nodes in the registry
type AppendNodeCapabilitiesRequest = MutateNodeCapabilitiesRequest

// AppendNodeCapabilities adds any new capabilities to the registry, merges the new capabilities with the existing capabilities
// of the node, and updates the nodes in the registry host the union of the new and existing capabilities.
func AppendNodeCapabilities(env deployment.Environment, req *AppendNodeCapabilitiesRequest) (deployment.ChangesetOutput, error) {
	c, err := req.convert(env)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	r, err := internal.AppendNodeCapabilitiesImpl(env.Logger, c)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	out := deployment.ChangesetOutput{}
	if req.UseMCMS() {
		if r.Ops == nil {
			return out, errors.New("expected MCMS operation to be non-nil")
		}
		timelocksPerChain := map[uint64]common.Address{
			c.Chain.Selector: c.ContractSet.Timelock.Address(),
		}
		proposerMCMSes := map[uint64]*gethwrappers.ManyChainMultiSig{
			c.Chain.Selector: c.ContractSet.ProposerMcm,
		}

		proposal, err := proposalutils.BuildProposalFromBatches(
			timelocksPerChain,
			proposerMCMSes,
			[]timelock.BatchChainOperation{*r.Ops},
			"proposal to set update node capabilities",
			req.MCMSConfig.MinDuration,
		)
		if err != nil {
			return out, fmt.Errorf("failed to build proposal: %w", err)
		}
		out.Proposals = []timelock.MCMSWithTimelockProposal{*proposal}
	}
	return out, nil
}

func (req *AppendNodeCapabilitiesRequest) convert(e deployment.Environment) (*internal.AppendNodeCapabilitiesRequest, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate UpdateNodeCapabilitiesRequest: %w", err)
	}
	registryChain, ok := e.Chains[req.RegistryChainSel]
	if !ok {
		return nil, fmt.Errorf("registry chain selector %d does not exist in environment", req.RegistryChainSel)
	}
	resp, err := internal.GetContractSets(e.Logger, &internal.GetContractSetsRequest{
		Chains:      map[uint64]deployment.Chain{req.RegistryChainSel: registryChain},
		AddressBook: e.ExistingAddresses,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get contract sets: %w", err)
	}
	contracts := resp.ContractSets[req.RegistryChainSel]

	return &internal.AppendNodeCapabilitiesRequest{
		Chain:             registryChain,
		ContractSet:       &contracts,
		P2pToCapabilities: req.P2pToCapabilities,
		UseMCMS:           req.UseMCMS(),
	}, nil
}
