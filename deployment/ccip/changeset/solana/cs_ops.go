package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// use these changesets to set the default code version
var _ cldf.ChangeSet[SetDefaultCodeVersionConfig] = SetDefaultCodeVersion

// use these changesets to update the SVM chain selector
var _ cldf.ChangeSet[UpdateSvmChainSelectorConfig] = UpdateSvmChainSelector

// use these changesets to update the enable manual execution after
var _ cldf.ChangeSet[UpdateEnableManualExecutionAfterConfig] = UpdateEnableManualExecutionAfter

// use these changesets to configure the CCIP version
var _ cldf.ChangeSet[ConfigureCCIPVersionConfig] = ConfigureCCIPVersion

// use these changesets to remove the offramp
var _ cldf.ChangeSet[RemoveOffRampConfig] = RemoveOffRamp

type SetDefaultCodeVersionConfig struct {
	ChainSelector uint64
	VersionEnum   uint8
	MCMS          *proposalutils.TimelockConfig
}

func (cfg SetDefaultCodeVersionConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateRouterConfig(chain); err != nil {
		return err
	}
	if err := chainState.ValidateOffRampConfig(chain); err != nil {
		return err
	}
	if err := chainState.ValidateFeeQuoterConfig(chain); err != nil {
		return err
	}
	return ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.FeeQuoter: true, shared.OffRamp: true, shared.Router: true})
}

func SetDefaultCodeVersion(e cldf.Environment, cfg SetDefaultCodeVersionConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Setting default code version", "chain_selector", cfg.ChainSelector, "new_code_version", cfg.VersionEnum)
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	chainState := state.SolChains[cfg.ChainSelector]
	routerUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.Router,
		solana.PublicKey{},
		"")
	offRampUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.OffRamp,
		solana.PublicKey{},
		"")
	feeQuoterUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.FeeQuoter,
		solana.PublicKey{},
		"")
	txns := make([]mcmsTypes.Transaction, 0)
	ixns := make([]solana.Instruction, 0)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.Router,
		solana.PublicKey{},
		"")
	solRouter.SetProgramID(chainState.Router)
	ixn, err := solRouter.NewSetDefaultCodeVersionInstruction(
		solRouter.CodeVersion(cfg.VersionEnum),
		chainState.RouterConfigPDA,
		authority,
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if routerUsingMCMS {
		tx, err := BuildMCMSTxn(ixn, chainState.Router.String(), shared.Router)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		txns = append(txns, *tx)
	} else {
		ixns = append(ixns, ixn)
	}

	authority = GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.OffRamp,
		solana.PublicKey{},
		"")
	solOffRamp.SetProgramID(chainState.OffRamp)
	ixn2, err := solOffRamp.NewSetDefaultCodeVersionInstruction(
		solOffRamp.CodeVersion(cfg.VersionEnum),
		chainState.OffRampConfigPDA,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if offRampUsingMCMS {
		tx, err := BuildMCMSTxn(ixn2, chainState.OffRamp.String(), shared.OffRamp)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		txns = append(txns, *tx)
	} else {
		ixns = append(ixns, ixn2)
	}

	authority = GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.FeeQuoter,
		solana.PublicKey{},
		"")
	solFeeQuoter.SetProgramID(chainState.FeeQuoter)
	ixn3, err := solFeeQuoter.NewSetDefaultCodeVersionInstruction(
		solFeeQuoter.CodeVersion(cfg.VersionEnum),
		chainState.FeeQuoterConfigPDA,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if feeQuoterUsingMCMS {
		tx, err := BuildMCMSTxn(ixn3, chainState.FeeQuoter.String(), shared.FeeQuoter)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		txns = append(txns, *tx)
	} else {
		ixns = append(ixns, ixn3)
	}

	if len(ixns) > 0 {
		if err = chain.Confirm(ixns); err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}

	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to SetDefaultCodeVersion in Solana", cfg.MCMS.MinDelay, txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}

type UpdateSvmChainSelectorConfig struct {
	OldChainSelector uint64
	NewChainSelector uint64
	MCMS             *proposalutils.TimelockConfig
}

func (cfg UpdateSvmChainSelectorConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	chainState := state.SolChains[cfg.OldChainSelector]
	chain := e.BlockChains.SolanaChains()[cfg.OldChainSelector]
	if err := chainState.ValidateRouterConfig(chain); err != nil {
		return err
	}
	if err := chainState.ValidateOffRampConfig(chain); err != nil {
		return err
	}
	return ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.Router: true, shared.OffRamp: true})
}

func UpdateSvmChainSelector(e cldf.Environment, cfg UpdateSvmChainSelectorConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Updating SVM chain selector", "old_chain_selector", cfg.OldChainSelector, "new_chain_selector", cfg.NewChainSelector)
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	chain := e.BlockChains.SolanaChains()[cfg.OldChainSelector]
	chainState := state.SolChains[cfg.OldChainSelector]
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)

	routerUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.Router,
		solana.PublicKey{},
		"")
	offRampUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.OffRamp,
		solana.PublicKey{},
		"")
	txns := make([]mcmsTypes.Transaction, 0)
	ixns := make([]solana.Instruction, 0)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.Router,
		solana.PublicKey{},
		"")
	solRouter.SetProgramID(chainState.Router)
	ixn, err := solRouter.NewUpdateSvmChainSelectorInstruction(
		cfg.NewChainSelector,
		routerConfigPDA,
		authority,
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if routerUsingMCMS {
		tx, err := BuildMCMSTxn(ixn, chainState.Router.String(), shared.Router)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		txns = append(txns, *tx)
	} else {
		ixns = append(ixns, ixn)
	}

	authority = GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.OffRamp,
		solana.PublicKey{},
		"")
	solOffRamp.SetProgramID(chainState.OffRamp)
	ixn2, err := solOffRamp.NewUpdateSvmChainSelectorInstruction(
		cfg.NewChainSelector,
		chainState.OffRampConfigPDA,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if offRampUsingMCMS {
		tx, err := BuildMCMSTxn(ixn2, chainState.OffRamp.String(), shared.OffRamp)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		txns = append(txns, *tx)
	} else {
		ixns = append(ixns, ixn2)
	}

	if len(ixns) > 0 {
		if err = chain.Confirm(ixns); err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}

	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.OldChainSelector, "proposal to UpdateSvmChainSelector in Solana", cfg.MCMS.MinDelay, txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}

type UpdateEnableManualExecutionAfterConfig struct {
	ChainSelector         uint64
	EnableManualExecution int64
	MCMS                  *proposalutils.TimelockConfig
}

func (cfg UpdateEnableManualExecutionAfterConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateOffRampConfig(chain); err != nil {
		return err
	}
	return ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.OffRamp: true})
}

func UpdateEnableManualExecutionAfter(e cldf.Environment, cfg UpdateEnableManualExecutionAfterConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Updating enable manual execution after", "chain_selector", cfg.ChainSelector, "enable_manual_execution_after", cfg.EnableManualExecution)
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	chainState := state.SolChains[cfg.ChainSelector]
	offRampUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.OffRamp,
		solana.PublicKey{},
		"")
	txns := make([]mcmsTypes.Transaction, 0)
	ixns := make([]solana.Instruction, 0)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.OffRamp,
		solana.PublicKey{},
		"")
	solOffRamp.SetProgramID(chainState.OffRamp)
	ixn2, err := solOffRamp.NewUpdateEnableManualExecutionAfterInstruction(
		cfg.EnableManualExecution,
		chainState.OffRampConfigPDA,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if offRampUsingMCMS {
		tx, err := BuildMCMSTxn(ixn2, chainState.OffRamp.String(), shared.OffRamp)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		txns = append(txns, *tx)
	} else {
		ixns = append(ixns, ixn2)
	}

	if len(ixns) > 0 {
		if err = chain.Confirm(ixns); err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}

	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to UpdateEnableManualExecutionAfter in Solana", cfg.MCMS.MinDelay, txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
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
	MCMS              *proposalutils.TimelockConfig
}

func (cfg ConfigureCCIPVersionConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateRouterConfig(chain); err != nil {
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
	return ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.Router: true})
}

func ConfigureCCIPVersion(e cldf.Environment, cfg ConfigureCCIPVersionConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	chainState := state.SolChains[cfg.ChainSelector]
	destChainStatePDA, _ := solState.FindDestChainStatePDA(cfg.DestChainSelector, chainState.Router)

	routerUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.Router,
		solana.PublicKey{},
		"")
	txns := make([]mcmsTypes.Transaction, 0)
	ixns := make([]solana.Instruction, 0)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.Router,
		solana.PublicKey{},
		"")
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
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
		}
	} else if cfg.Operation == Rollback {
		ixn, err = solRouter.NewRollbackCcipVersionForDestChainInstruction(
			cfg.DestChainSelector,
			destChainStatePDA,
			chainState.RouterConfigPDA,
			authority,
		).ValidateAndBuild()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
		}
	}

	if routerUsingMCMS {
		tx, err := BuildMCMSTxn(ixn, chainState.Router.String(), shared.Router)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		txns = append(txns, *tx)
	} else {
		ixns = append(ixns, ixn)
	}

	if len(ixns) > 0 {
		if err = chain.Confirm(ixns); err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}

	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to ConfigureCCIPVersion in Solana", cfg.MCMS.MinDelay, txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}

type RemoveOffRampConfig struct {
	ChainSelector uint64
	OffRamp       solana.PublicKey
	MCMS          *proposalutils.TimelockConfig
}

func (cfg RemoveOffRampConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateRouterConfig(chain); err != nil {
		return err
	}
	return ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.Router: true})
}

func RemoveOffRamp(e cldf.Environment, cfg RemoveOffRampConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	chainState := state.SolChains[cfg.ChainSelector]
	routerUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.Router,
		solana.PublicKey{},
		"")
	txns := make([]mcmsTypes.Transaction, 0)
	ixns := make([]solana.Instruction, 0)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.Router,
		solana.PublicKey{},
		"")
	solRouter.SetProgramID(chainState.Router)
	ixn, err := solRouter.NewRemoveOfframpInstruction(
		cfg.ChainSelector,
		cfg.OffRamp,
		chainState.OffRamp,
		chainState.RouterConfigPDA,
		authority,
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if routerUsingMCMS {
		tx, err := BuildMCMSTxn(ixn, chainState.Router.String(), shared.Router)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		txns = append(txns, *tx)
	} else {
		ixns = append(ixns, ixn)
	}

	if len(ixns) > 0 {
		if err = chain.Confirm(ixns); err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}

	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to RemoveOffRamp in Solana", cfg.MCMS.MinDelay, txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}
