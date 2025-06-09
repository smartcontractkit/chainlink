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

// use this changeset to disable a remote chain on solana
var _ cldf.ChangeSet[DisableRemoteChainConfig] = DisableRemoteChain

type DisableRemoteChainConfig struct {
	ChainSelector uint64
	RemoteChains  []uint64
	MCMS          *proposalutils.TimelockConfig
}

func (cfg DisableRemoteChainConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateRouterConfig(chain); err != nil {
		return err
	}
	if err := chainState.ValidateFeeQuoterConfig(chain); err != nil {
		return err
	}
	if err := chainState.ValidateOffRampConfig(chain); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.FeeQuoter: true, shared.OffRamp: true}); err != nil {
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
		if err := state.ValidateRamp(remote, shared.OnRamp); err != nil {
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

func DisableRemoteChain(e cldf.Environment, cfg DisableRemoteChainConfig) (cldf.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	s, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	txns, err := doDisableRemoteChain(e, s, cfg)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// create proposals for ixns
	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to disable remote chains in Solana", cfg.MCMS.MinDelay, txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}
	return cldf.ChangesetOutput{}, nil
}

func doDisableRemoteChain(
	e cldf.Environment,
	s stateview.CCIPOnChainState,
	cfg DisableRemoteChainConfig) ([]mcmsTypes.Transaction, error) {
	txns := make([]mcmsTypes.Transaction, 0)
	ixns := make([]solana.Instruction, 0)
	chainSel := cfg.ChainSelector
	chain := e.BlockChains.SolanaChains()[chainSel]
	chainState := s.SolChains[chainSel]
	feeQuoterID := s.SolChains[chainSel].FeeQuoter
	offRampID := s.SolChains[chainSel].OffRamp
	feeQuoterUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.FeeQuoter,
		solana.PublicKey{},
		"")
	offRampUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.OffRamp,
		solana.PublicKey{},
		"")

	for _, remoteChainSel := range cfg.RemoteChains {
		// verified while loading state
		fqDestChainPDA, _, _ := solState.FindFqDestChainPDA(remoteChainSel, feeQuoterID)
		offRampSourceChainPDA, _, _ := solState.FindOfframpSourceChainPDA(remoteChainSel, s.SolChains[chainSel].OffRamp)

		solFeeQuoter.SetProgramID(feeQuoterID)
		authority := GetAuthorityForIxn(
			&e,
			chain,
			chainState,
			cfg.MCMS,
			shared.FeeQuoter,
			solana.PublicKey{},
			"",
		)
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
			tx, err := BuildMCMSTxn(feeQuoterIx, feeQuoterID.String(), shared.FeeQuoter)
			if err != nil {
				return txns, fmt.Errorf("failed to create transaction: %w", err)
			}
			txns = append(txns, *tx)
		} else {
			ixns = append(ixns, feeQuoterIx)
		}

		solOffRamp.SetProgramID(offRampID)
		authority = GetAuthorityForIxn(
			&e,
			chain,
			chainState,
			cfg.MCMS,
			shared.OffRamp,
			solana.PublicKey{},
			"")
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
			tx, err := BuildMCMSTxn(offRampIx, offRampID.String(), shared.OffRamp)
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
