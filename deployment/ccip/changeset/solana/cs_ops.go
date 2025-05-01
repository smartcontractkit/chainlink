package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

// use these changesets to set the default code version
var _ deployment.ChangeSet[SetDefaultCodeVersionConfig] = SetDefaultCodeVersion

// use these changesets to update the SVM chain selector
var _ deployment.ChangeSet[UpdateSvmChainSelectorConfig] = UpdateSvmChainSelector

// use these changesets to update the enable manual execution after
var _ deployment.ChangeSet[UpdateEnableManualExecutionAfterConfig] = UpdateEnableManualExecutionAfter

// use these changesets to configure the CCIP version
var _ deployment.ChangeSet[ConfigureCCIPVersionConfig] = ConfigureCCIPVersion

// use these changesets to remove the offramp
var _ deployment.ChangeSet[RemoveOffRampConfig] = RemoveOffRamp

type SetDefaultCodeVersionConfig struct {
	ChainSelector uint64
	VersionEnum   uint8
	MCMSSolana    *MCMSConfigSolana
}

func (cfg SetDefaultCodeVersionConfig) Validate(e deployment.Environment) error {
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}
	if err := validateOffRampConfig(chain, chainState); err != nil {
		return err
	}
	if err := validateFeeQuoterConfig(chain, chainState); err != nil {
		return err
	}
	return ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, solana.PublicKey{})
}

func SetDefaultCodeVersion(e deployment.Environment, cfg SetDefaultCodeVersionConfig) (deployment.ChangesetOutput, error) {
	e.Logger.Infow("Setting default code version", "chain_selector", cfg.ChainSelector, "new_code_version", cfg.VersionEnum)
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]

	routerUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock
	offRampUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.OffRampOwnedByTimelock
	feeQuoterUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.FeeQuoterOwnedByTimelock
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
	ixn, err := solRouter.NewSetDefaultCodeVersionInstruction(
		solRouter.CodeVersion(cfg.VersionEnum),
		chainState.RouterConfigPDA,
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
	ixn2, err := solOffRamp.NewSetDefaultCodeVersionInstruction(
		solOffRamp.CodeVersion(cfg.VersionEnum),
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

	authority, err = GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.FeeQuoter,
		solana.PublicKey{})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
	solFeeQuoter.SetProgramID(chainState.FeeQuoter)
	ixn3, err := solFeeQuoter.NewSetDefaultCodeVersionInstruction(
		solFeeQuoter.CodeVersion(cfg.VersionEnum),
		chainState.FeeQuoterConfigPDA,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if feeQuoterUsingMCMS {
		tx, err := BuildMCMSTxn(ixn3, chainState.FeeQuoter.String(), ccipChangeset.FeeQuoter)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		txns = append(txns, *tx)
	} else {
		ixns = append(ixns, ixn3)
	}

	if len(ixns) > 0 {
		if err = chain.Confirm(ixns); err != nil {
			return deployment.ChangesetOutput{}, err
		}
	}

	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to SetDefaultCodeVersion in Solana", cfg.MCMSSolana.MCMS.MinDelay, txns)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return deployment.ChangesetOutput{}, nil
}

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

type UpdateEnableManualExecutionAfterConfig struct {
	ChainSelector         uint64
	EnableManualExecution int64
	MCMSSolana            *MCMSConfigSolana
}

func (cfg UpdateEnableManualExecutionAfterConfig) Validate(e deployment.Environment) error {
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateOffRampConfig(chain, chainState); err != nil {
		return err
	}
	return ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, solana.PublicKey{})
}

func UpdateEnableManualExecutionAfter(e deployment.Environment, cfg UpdateEnableManualExecutionAfterConfig) (deployment.ChangesetOutput, error) {
	e.Logger.Infow("Updating enable manual execution after", "chain_selector", cfg.ChainSelector, "enable_manual_execution_after", cfg.EnableManualExecution)
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]

	offRampUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.OffRampOwnedByTimelock
	txns := make([]mcmsTypes.Transaction, 0)
	ixns := make([]solana.Instruction, 0)
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.OffRamp,
		solana.PublicKey{})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
	solOffRamp.SetProgramID(chainState.OffRamp)
	ixn2, err := solOffRamp.NewUpdateEnableManualExecutionAfterInstruction(
		cfg.EnableManualExecution,
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
			e, cfg.ChainSelector, "proposal to UpdateEnableManualExecutionAfter in Solana", cfg.MCMSSolana.MCMS.MinDelay, txns)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return deployment.ChangesetOutput{}, nil
}

type CCIPVersionOp int

const (
	Bump CCIPVersionOp = iota
	Rollback
)

type ConfigureCCIPVersionConfig struct {
	ChainSelector     uint64
	DestChainSelector uint64
	Operation         CCIPVersionOp
	MCMSSolana        *MCMSConfigSolana
}

func (cfg ConfigureCCIPVersionConfig) Validate(e deployment.Environment) error {
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}
	routerDestChainPDA, err := solState.FindDestChainStatePDA(cfg.DestChainSelector, chainState.Router)
	if err != nil {
		return fmt.Errorf("failed to find dest chain state pda for remote chain %d: %w", cfg.DestChainSelector, err)
	}
	var destChainStateAccount solRouter.DestChain
	err = chain.GetAccountDataBorshInto(context.Background(), routerDestChainPDA, &destChainStateAccount)
	if err != nil {
		return fmt.Errorf("remote %d is not configured on solana chain %d", cfg.DestChainSelector, cfg.ChainSelector)
	}
	return ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, solana.PublicKey{})
}

func ConfigureCCIPVersion(e deployment.Environment, cfg ConfigureCCIPVersionConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	destChainStatePDA, _ := solState.FindDestChainStatePDA(cfg.DestChainSelector, chainState.Router)

	routerUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock
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
	var ixn solana.Instruction
	if cfg.Operation == Bump {
		ixn, err = solRouter.NewBumpCcipVersionForDestChainInstruction(
			cfg.DestChainSelector,
			destChainStatePDA,
			chainState.RouterConfigPDA,
			authority,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
		}
	} else if cfg.Operation == Rollback {
		ixn, err = solRouter.NewRollbackCcipVersionForDestChainInstruction(
			cfg.DestChainSelector,
			destChainStatePDA,
			chainState.RouterConfigPDA,
			authority,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
		}
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

	if len(ixns) > 0 {
		if err = chain.Confirm(ixns); err != nil {
			return deployment.ChangesetOutput{}, err
		}
	}

	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to ConfigureCCIPVersion in Solana", cfg.MCMSSolana.MCMS.MinDelay, txns)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return deployment.ChangesetOutput{}, nil
}

type RemoveOffRampConfig struct {
	ChainSelector uint64
	OffRamp       solana.PublicKey
	MCMSSolana    *MCMSConfigSolana
}

func (cfg RemoveOffRampConfig) Validate(e deployment.Environment) error {
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}
	return ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, solana.PublicKey{})
}

func RemoveOffRamp(e deployment.Environment, cfg RemoveOffRampConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]

	routerUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock
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
	var ixn solana.Instruction
	ixn, err = solRouter.NewRemoveOfframpInstruction(
		cfg.ChainSelector,
		cfg.OffRamp,
		chainState.OffRamp,
		chainState.RouterConfigPDA,
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

	if len(ixns) > 0 {
		if err = chain.Confirm(ixns); err != nil {
			return deployment.ChangesetOutput{}, err
		}
	}

	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to RemoveOffRamp in Solana", cfg.MCMSSolana.MCMS.MinDelay, txns)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return deployment.ChangesetOutput{}, nil
}
