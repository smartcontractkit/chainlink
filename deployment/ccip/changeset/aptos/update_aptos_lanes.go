package aptos

import (
	"errors"
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	config "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	seq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var _ cldf.ChangeSetV2[config.UpdateAptosLanesConfig] = AddAptosLanes{}

// AddAptosLane implements adding a new lane to an existing Aptos CCIP deployment
type AddAptosLanes struct{}

func (cs AddAptosLanes) VerifyPreconditions(env cldf.Environment, cfg config.UpdateAptosLanesConfig) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load existing Aptos onchain state: %w", err)
	}
	supportedChains := state.SupportedChains()
	if cfg.AptosMCMSConfig == nil {
		return errors.New("config for Aptos MCMS is required for AddAptosLanes changeset")
	}
	// For every configured lane validate Aptos source or destination chain definitions
	for _, laneCfg := range cfg.Lanes {
		// Source cannot be an unknown.
		if _, ok := supportedChains[laneCfg.Source.GetSelector()]; !ok {
			return fmt.Errorf("source chain %d is not a supported", laneCfg.Source.GetSelector())
		}
		// Destination cannot be an unknown.
		if _, ok := supportedChains[laneCfg.Dest.GetSelector()]; !ok {
			return fmt.Errorf("destination chain %d is not a supported", laneCfg.Dest.GetSelector())
		}
		if laneCfg.Source.GetChainFamily() == chainsel.FamilyAptos {
			aptosChain, exists := env.BlockChains.AptosChains()[laneCfg.Source.GetSelector()]
			if !exists {
				return fmt.Errorf("source Aptos chain %d is not in env", laneCfg.Source.GetSelector())
			}
			err := laneCfg.Source.(config.AptosChainDefinition).Validate(
				aptosChain.Client,
				state.AptosChains[laneCfg.Source.GetSelector()],
			)
			if err != nil {
				return fmt.Errorf("failed to validate Aptos source chain %d: %w", laneCfg.Source.GetSelector(), err)
			}
		}
		if laneCfg.Dest.GetChainFamily() == chainsel.FamilyAptos {
			aptosChain, exists := env.BlockChains.AptosChains()[laneCfg.Dest.GetSelector()]
			if !exists {
				return fmt.Errorf("destination Aptos chain %d is not in env", laneCfg.Dest.GetSelector())
			}
			err := laneCfg.Dest.(config.AptosChainDefinition).Validate(
				aptosChain.Client,
				state.AptosChains[laneCfg.Dest.GetSelector()],
			)
			if err != nil {
				return fmt.Errorf("failed to validate Aptos destination chain %d: %w", laneCfg.Dest.GetSelector(), err)
			}
		}
	}

	// This EVM specific changeset will be called from within this Aptos changeset, hence, we're verifying it here
	// TODO: this is an anti-pattern, change this once EVM changesets are refactored as Operations
	evmUpdateCfg := config.ToEVMUpdateLanesConfig(cfg)
	err = v1_6.UpdateLanesPrecondition(env, evmUpdateCfg)
	if err != nil {
		return err
	}
	return nil
}

func (cs AddAptosLanes) Apply(env cldf.Environment, cfg config.UpdateAptosLanesConfig) (cldf.ChangesetOutput, error) {
	var (
		timeLockProposals []mcms.TimelockProposal
		mcmsOperations    []mcmstypes.BatchOperation
	)

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
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}

	updateInputsByAptosChain := seq.ToAptosUpdateLanesConfig(state.AptosChains, cfg)
	for aptosChainSel, sequenceInput := range updateInputsByAptosChain {
		deps := operation.AptosDeps{
			AptosChain:       env.BlockChains.AptosChains()[aptosChainSel],
			CCIPOnChainState: state,
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
			state.AptosChains[aptosChainSel].MCMSAddress,
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
