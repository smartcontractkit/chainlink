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
	cs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

// use this changeset to disable a remote chain on solana
var _ deployment.ChangeSet[DisableRemoteChainConfig] = DisableRemoteChain

type DisableRemoteChainConfig struct {
	ChainSelector uint64
	RemoteChains  []uint64
	MCMSSolana    *MCMSConfigSolana
}

func (cfg DisableRemoteChainConfig) Validate(e deployment.Environment) error {
	state, err := cs.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}
	if err := validateFeeQuoterConfig(chain, chainState); err != nil {
		return err
	}
	if err := validateOffRampConfig(chain, chainState); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, solana.PublicKey{}); err != nil {
		return err
	}
	var routerConfigAccount solRouter.Config
	// already validated that router config exists
	_ = chain.GetAccountDataBorshInto(context.Background(), chainState.RouterConfigPDA, &routerConfigAccount)

	supportedChains := state.SupportedChains()
	for _, remote := range cfg.RemoteChains {
		if _, ok := supportedChains[remote]; !ok {
			return fmt.Errorf("remote chain %d is not supported", remote)
		}
		if remote == routerConfigAccount.SvmChainSelector {
			return fmt.Errorf("cannot disable remote chain %d with same chain selector as current chain %d", remote, cfg.ChainSelector)
		}
		if err := state.ValidateRamp(remote, cs.OnRamp); err != nil {
			return err
		}
		routerDestChainPDA, err := solState.FindDestChainStatePDA(remote, chainState.Router)
		if err != nil {
			return fmt.Errorf("failed to find dest chain state pda for remote chain %d: %w", remote, err)
		}
		var destChainStateAccount solRouter.DestChain
		err = chain.GetAccountDataBorshInto(context.Background(), routerDestChainPDA, &destChainStateAccount)
		if err != nil {
			return fmt.Errorf("remote %d is not configured on solana chain %d", remote, cfg.ChainSelector)
		}
	}
	return nil
}

func DisableRemoteChain(e deployment.Environment, cfg DisableRemoteChainConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	s, err := cs.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	txns, err := doDisableRemoteChain(e, s, cfg)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// create proposals for ixns
	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to disable remote chains in Solana", cfg.MCMSSolana.MCMS.MinDelay, txns)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}
	return deployment.ChangesetOutput{}, nil
}

func doDisableRemoteChain(
	e deployment.Environment,
	s cs.CCIPOnChainState,
	cfg DisableRemoteChainConfig) ([]mcmsTypes.Transaction, error) {
	txns := make([]mcmsTypes.Transaction, 0)
	ixns := make([]solana.Instruction, 0)
	chainSel := cfg.ChainSelector
	chain := e.SolChains[chainSel]
	feeQuoterID := s.SolChains[chainSel].FeeQuoter
	offRampID := s.SolChains[chainSel].OffRamp
	feeQuoterUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.FeeQuoterOwnedByTimelock
	offRampUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.OffRampOwnedByTimelock

	for _, remoteChainSel := range cfg.RemoteChains {
		// verified while loading state
		fqDestChainPDA, _, _ := solState.FindFqDestChainPDA(remoteChainSel, feeQuoterID)
		offRampSourceChainPDA, _, _ := solState.FindOfframpSourceChainPDA(remoteChainSel, s.SolChains[chainSel].OffRamp)

		solFeeQuoter.SetProgramID(feeQuoterID)
		authority, err := GetAuthorityForIxn(
			&e,
			chain,
			cfg.MCMSSolana,
			cs.FeeQuoter,
			solana.PublicKey{})
		if err != nil {
			return txns, fmt.Errorf("failed to get authority for ixn: %w", err)
		}
		feeQuoterIx, err := solFeeQuoter.NewDisableDestChainInstruction(
			remoteChainSel,
			s.SolChains[chainSel].FeeQuoterConfigPDA,
			fqDestChainPDA,
			authority,
		).ValidateAndBuild()
		if err != nil {
			return txns, fmt.Errorf("failed to generate instructions: %w", err)
		}
		if feeQuoterUsingMCMS {
			tx, err := BuildMCMSTxn(feeQuoterIx, feeQuoterID.String(), cs.FeeQuoter)
			if err != nil {
				return txns, fmt.Errorf("failed to create transaction: %w", err)
			}
			txns = append(txns, *tx)
		} else {
			ixns = append(ixns, feeQuoterIx)
		}

		solOffRamp.SetProgramID(offRampID)
		authority, err = GetAuthorityForIxn(
			&e,
			chain,
			cfg.MCMSSolana,
			cs.OffRamp,
			solana.PublicKey{})
		if err != nil {
			return txns, fmt.Errorf("failed to get authority for ixn: %w", err)
		}
		offRampIx, err := solOffRamp.NewDisableSourceChainSelectorInstruction(
			remoteChainSel,
			offRampSourceChainPDA,
			s.SolChains[chainSel].OffRampConfigPDA,
			authority,
		).ValidateAndBuild()
		if err != nil {
			return txns, fmt.Errorf("failed to generate instructions: %w", err)
		}
		if offRampUsingMCMS {
			tx, err := BuildMCMSTxn(offRampIx, offRampID.String(), cs.OffRamp)
			if err != nil {
				return txns, fmt.Errorf("failed to create transaction: %w", err)
			}
			txns = append(txns, *tx)
		} else {
			ixns = append(ixns, offRampIx)
		}
		if len(ixns) > 0 {
			err = chain.Confirm(ixns)
			if err != nil {
				return txns, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}
	}

	return txns, nil
}
