package changeset

import (
	"errors"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
)

type AddNopsRequest = struct {
	RegistryChainSel uint64
	Nops             []kcr.CapabilitiesRegistryNodeOperator

	MCMSConfig *MCMSConfig // if non-nil, the changes will be proposed using MCMS.
}

var _ deployment.ChangeSet[*AddNopsRequest] = AddNops

func AddNops(env deployment.Environment, req *AddNopsRequest) (deployment.ChangesetOutput, error) {
	for _, nop := range req.Nops {
		env.Logger.Infow("input NOP", "address", nop.Admin, "name", nop.Name)
	}
	registryChain, ok := env.Chains[req.RegistryChainSel]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("registry chain selector %d does not exist in environment", req.RegistryChainSel)
	}
	cs, err := GetContractSets(env.Logger, &GetContractSetsRequest{
		Chains:      map[uint64]deployment.Chain{req.RegistryChainSel: registryChain},
		AddressBook: env.ExistingAddresses,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get contract sets: %w", err)
	}
	contractSet, exists := cs.ContractSets[req.RegistryChainSel]
	if !exists {
		return deployment.ChangesetOutput{}, fmt.Errorf("contract set not found for chain %d", req.RegistryChainSel)
	}

	useMCMS := req.MCMSConfig != nil
	req2 := internal.RegisterNOPSRequest{
		Env:                   &env,
		RegistryChainSelector: req.RegistryChainSel,
		Nops:                  req.Nops,
		UseMCMS:               useMCMS,
	}
	resp, err := internal.RegisterNOPS(env.GetContext(), env.Logger, req2)

	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	env.Logger.Infow("registered NOPs", "nops", resp.Nops)
	out := deployment.ChangesetOutput{}
	if useMCMS {
		if resp.Ops == nil {
			return out, errors.New("expected MCMS operation to be non-nil")
		}
		timelocksPerChain := map[uint64]gethcommon.Address{
			registryChain.Selector: contractSet.Timelock.Address(),
		}
		proposerMCMSes := map[uint64]*gethwrappers.ManyChainMultiSig{
			registryChain.Selector: contractSet.ProposerMcm,
		}

		proposal, err := proposalutils.BuildProposalFromBatches(
			timelocksPerChain,
			proposerMCMSes,
			[]timelock.BatchChainOperation{*resp.Ops},
			"proposal to add NOPs",
			req.MCMSConfig.MinDuration,
		)
		if err != nil {
			return out, fmt.Errorf("failed to build proposal: %w", err)
		}
		out.Proposals = []timelock.MCMSWithTimelockProposal{*proposal}
	}

	return out, nil
}
