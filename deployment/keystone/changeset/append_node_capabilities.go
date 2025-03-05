package changeset

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/types"

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
	c, contractSet, err := req.convert(env)
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
		timelocksPerChain := map[uint64]string{
			c.Chain.Selector: contractSet.Timelock.Address().Hex(),
		}
		proposerMCMSes := map[uint64]string{
			c.Chain.Selector: contractSet.ProposerMcm.Address().Hex(),
		}
		inspector, err := proposalutils.McmsInspectorForChain(env, req.RegistryChainSel)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		inspectorPerChain := map[uint64]sdk.Inspector{
			req.RegistryChainSel: inspector,
		}

		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			env,
			timelocksPerChain,
			proposerMCMSes,
			inspectorPerChain,
			[]types.BatchOperation{*r.Ops},
			"proposal to set update node capabilities",
			req.MCMSConfig.MinDuration,
		)
		if err != nil {
			return out, fmt.Errorf("failed to build proposal: %w", err)
		}
		out.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}
	}
	return out, nil
}

func (req *AppendNodeCapabilitiesRequest) convert(e deployment.Environment) (*internal.AppendNodeCapabilitiesRequest, *ContractSet, error) {
	if err := req.Validate(e); err != nil {
		return nil, nil, fmt.Errorf("failed to validate UpdateNodeCapabilitiesRequest: %w", err)
	}
	registryChain := e.Chains[req.RegistryChainSel] // exists because of the validation above
	resp, err := GetContractSets(e.Logger, &GetContractSetsRequest{
		Chains:      map[uint64]deployment.Chain{req.RegistryChainSel: registryChain},
		AddressBook: e.ExistingAddresses,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get contract sets: %w", err)
	}
	contractSet := resp.ContractSets[req.RegistryChainSel]

	return &internal.AppendNodeCapabilitiesRequest{
		Chain:                registryChain,
		CapabilitiesRegistry: contractSet.CapabilitiesRegistry,
		P2pToCapabilities:    req.P2pToCapabilities,
		UseMCMS:              req.UseMCMS(),
	}, &contractSet, nil
}
