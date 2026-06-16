package changeset

import (
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	proposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

// ethBalMonMCMSContractTypeForAction selects the MCMS contract type in the datastore for EthBalMon
// timelock proposals, matching operations.generateMCMSProposals.
func ethBalMonMCMSContractTypeForAction(action mcmstypes.TimelockAction) cldf.ContractType {
	if action == mcmstypes.TimelockActionBypass {
		return commontypes.BypasserManyChainMultisig
	}
	return commontypes.ProposerManyChainMultisig
}

// ethBalMonMCMSContractTypeForProposal resolves MCMS datastore type from optional timelock config.
// When cfg is nil or MCMSAction is empty, the action is treated as schedule (non-bypass), so the proposer MCM is used.
func ethBalMonMCMSContractTypeForProposal(cfg *proposalutils.TimelockConfig) cldf.ContractType {
	action := mcmstypes.TimelockActionSchedule
	if cfg != nil && cfg.MCMSAction != "" {
		action = cfg.MCMSAction
	}
	return ethBalMonMCMSContractTypeForAction(action)
}

func ethBalMonProposalTimelockConfig(cfg *proposalutils.TimelockConfig) proposalutils.TimelockConfig {
	if cfg == nil {
		return proposalutils.TimelockConfig{MinDelay: 0}
	}
	return *cfg
}

// deployEthBalMonAcceptOwnershipMCMSAction returns the MCMS timelock action for the post-deploy accept-ownership proposal.
// When cfg is nil or MCMSAction is unset, the deploy flow defaults to bypass.
func deployEthBalMonAcceptOwnershipMCMSAction(cfg *proposalutils.TimelockConfig) mcmstypes.TimelockAction {
	if cfg == nil || cfg.MCMSAction == "" {
		return mcmstypes.TimelockActionBypass
	}
	return cfg.MCMSAction
}

func deployEthBalMonAcceptOwnershipTimelockConfig(cfg *proposalutils.TimelockConfig) proposalutils.TimelockConfig {
	out := proposalutils.TimelockConfig{MinDelay: 0}
	if cfg != nil {
		out = *cfg
	}
	if out.MCMSAction == "" {
		out.MCMSAction = mcmstypes.TimelockActionBypass
	}
	return out
}
