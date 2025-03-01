package solana

import (
	"fmt"

	"github.com/gagliardetto/solana-go"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

type UpdateSvmChainSelectorConfig struct {
	OldChainSelector uint64
	NewChainSelector uint64
	MCMSSolana       *MCMSConfigSolana
}

func (cfg UpdateSvmChainSelectorConfig) Validate(e deployment.Environment) error {
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.SolChains[cfg.OldChainSelector]
	chain := e.SolChains[cfg.OldChainSelector]
	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}
	if err := validateOffRampConfig(chain, chainState); err != nil {
		return err
	}
	return ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, solana.PublicKey{})
}

func UpdateSvmChainSelector(e deployment.Environment, cfg UpdateSvmChainSelectorConfig) (deployment.ChangesetOutput, error) {
	e.Logger.Infow("Updating SVM chain selector", "old_chain_selector", cfg.OldChainSelector, "new_chain_selector", cfg.NewChainSelector)
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.OldChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.OldChainSelector]
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)

	routerUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock
	offRampUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.OffRampOwnedByTimelock
	txns := make([]mcmsTypes.Transaction, 0)
	ixns := make([]solana.Instruction, 0)
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.Router,
		solana.PublicKey{})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
	solRouter.SetProgramID(chainState.Router)
	ixn, err := solRouter.NewUpdateSvmChainSelectorInstruction(
		cfg.NewChainSelector,
		routerConfigPDA,
		authority,
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if routerUsingMCMS {
		tx, err := BuildMCMSTxn(ixn, chainState.Router.String(), ccipChangeset.Router)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		txns = append(txns, *tx)
	} else {
		ixns = append(ixns, ixn)
	}

	authority, err = GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.OffRamp,
		solana.PublicKey{})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
	solOffRamp.SetProgramID(chainState.OffRamp)
	ixn2, err := solOffRamp.NewUpdateSvmChainSelectorInstruction(
		cfg.NewChainSelector,
		chainState.OffRampConfigPDA,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if offRampUsingMCMS {
		tx, err := BuildMCMSTxn(ixn2, chainState.OffRamp.String(), ccipChangeset.OffRamp)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		txns = append(txns, *tx)
	} else {
		ixns = append(ixns, ixn2)
	}

	if len(ixns) > 0 {
		if err = chain.Confirm(ixns); err != nil {
			return deployment.ChangesetOutput{}, err
		}
	}

	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.OldChainSelector, "proposal to UpdateSvmChainSelector in Solana", cfg.MCMSSolana.MCMS.MinDelay, txns)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return deployment.ChangesetOutput{}, nil
}
