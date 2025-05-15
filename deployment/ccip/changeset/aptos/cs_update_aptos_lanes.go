package aptos

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	seq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	aptosstate "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/aptos"
)

var _ cldf.ChangeSetV2[config.UpdateAptosLanesConfig] = AddAptosLanes{}

// AddAptosLane implements adding a new lane to an existing Aptos CCIP deployment
type AddAptosLanes struct{}

func (cs AddAptosLanes) VerifyPreconditions(env cldf.Environment, cfg config.UpdateAptosLanesConfig) error {
	// TODO: Implement verification logic - check chain selector validity, MCMS configuration, etc.
	// Placeholder implementation to show expected structure

	// This EVM specific changeset will be called from within this Aptos changeset, hence, we're verifying it here
	// TODO: this is an anti-pattern, change this once EVM changesets are refactored as Operations
	evmUpdateCfg := config.ToEVMUpdateLanesConfig(cfg)
	err := v1_6.UpdateLanesPrecondition(env, evmUpdateCfg)
	if err != nil {
		return err
	}
	return nil
}

func (cs AddAptosLanes) Apply(env cldf.Environment, cfg config.UpdateAptosLanesConfig) (cldf.ChangesetOutput, error) {
	timeLockProposals := []mcms.TimelockProposal{}
	mcmsOperations := []types.BatchOperation{}
	seqReports := make([]operations.Report[any, any], 0)

	// Add lane on EVM chains
	// TODO: applying a changeset within another changeset is an anti-pattern. Using it here until EVM is refactored into Operations
	evmUpdatesInput := config.ToEVMUpdateLanesConfig(cfg)
	out, err := v1_6.UpdateLanesLogic(env, cfg.EVMMCMSConfig, evmUpdatesInput)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	timeLockProposals = append(timeLockProposals, out.MCMSTimelockProposals...)

	// Add lane on Aptos chains
	// Execute UpdateAptosLanesSequence for each aptos chain
	state, err := aptosstate.LoadOnchainStateAptos(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}
	ccipState, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load EVM onchain state: %w", err)
	}

	updateInputsByAptosChain := seq.ConvertToUpdateAptosLanesSeqInput(state, cfg)
	for aptosChainSel, sequenceInput := range updateInputsByAptosChain {
		deps := operation.AptosDeps{
			AptosChain:       env.AptosChains[aptosChainSel],
			OnChainState:     state[aptosChainSel],
			CCIPOnChainState: ccipState,
		}
		// Execute the sequence
		updateSeqReport, err := operations.ExecuteSequence(env.OperationsBundle, seq.UpdateAptosLanesSequence, deps, sequenceInput)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		seqReports = append(seqReports, updateSeqReport.ExecutionReports...)
		mcmsOperations = append(mcmsOperations, updateSeqReport.Output)

		// Generate MCMS proposals
		proposal, err := utils.GenerateProposal(
			deps.AptosChain.Client,
			state[aptosChainSel].MCMSAddress,
			deps.AptosChain.Selector,
			mcmsOperations,
			"Update lanes on Aptos chain",
			*cfg.AptosMCMSConfig,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal for Aptos chain %d: %w", aptosChainSel, err)
		}
		timeLockProposals = append(timeLockProposals, *proposal)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: timeLockProposals,
		Reports:               seqReports,
	}, nil
}
