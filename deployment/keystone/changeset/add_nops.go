package changeset

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

type AddNopsRequest struct {
	RegistryChainSel uint64
	Nops             []kcr.CapabilitiesRegistryNodeOperator

	MCMSConfig *MCMSConfig // if non-nil, the changes will be proposed using MCMS.

	RegistryRef datastore.AddressRefKey
}

func (req *AddNopsRequest) Validate(env deployment.Environment) error {
	if len(req.Nops) == 0 {
		return errors.New("no NOPs provided")
	}
	if req.RegistryChainSel == 0 {
		return errors.New("registry chain selector is required")
	}
	if err := shouldUseDatastore(env, req.RegistryRef); err != nil {
		return fmt.Errorf("invalid registry reference: %w", err)
	}
	return nil
}

var _ deployment.ChangeSet[*AddNopsRequest] = AddNops

func AddNops(env deployment.Environment, req *AddNopsRequest) (deployment.ChangesetOutput, error) {
	if err := req.Validate(env); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid request: %w", err)
	}
	for _, nop := range req.Nops {
		env.Logger.Infow("input NOP", "address", nop.Admin, "name", nop.Name)
	}
	registryChain, ok := env.Chains[req.RegistryChainSel]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("registry chain selector %d does not exist in environment", req.RegistryChainSel)
	}
	capReg, err := loadCapabilityRegistry(registryChain, env, req.RegistryRef)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load capability registry: %w", err)
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
		timelocksPerChain := map[uint64]string{
			registryChain.Selector: capReg.McmsContracts.Timelock.Address().Hex(),
		}
		proposerMCMSes := map[uint64]string{
			registryChain.Selector: capReg.McmsContracts.ProposerMcm.Address().Hex(),
		}
		inspector, err := proposalutils.McmsInspectorForChain(env, req.RegistryChainSel)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		inspectorPerChain := map[uint64]mcmssdk.Inspector{
			req.RegistryChainSel: inspector,
		}

		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			env,
			timelocksPerChain,
			proposerMCMSes,
			inspectorPerChain,
			[]mcmstypes.BatchOperation{*resp.Ops},
			"proposal to add NOPs",
			proposalutils.TimelockConfig{MinDelay: req.MCMSConfig.MinDuration},
		)
		if err != nil {
			return out, fmt.Errorf("failed to build proposal: %w", err)
		}
		out.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}
	}

	return out, nil
}
