package solana

import (
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanaStateView "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
)

type ApproveTokensForFeeBillingSignerConfig struct {
	ChainSelector uint64
	TokenProgram  solana.PublicKey
	Amount        uint64
	Decimals      uint8
}

func (cfg ApproveTokensForFeeBillingSignerConfig) Validate(e cldf.Environment) (solanaStateView.CCIPChainState, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return solanaStateView.CCIPChainState{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainState, found := state.SolChains[cfg.ChainSelector]
	if !found {
		return solanaStateView.CCIPChainState{}, fmt.Errorf("failed to get chain state for selector %d", cfg.ChainSelector)
	}
	if cfg.TokenProgram == (solana.PublicKey{}) {
		return solanaStateView.CCIPChainState{}, errors.New("token program is not set")
	}

	return chainState, nil
}

func ApproveTokensForFeeBillingSigner(e cldf.Environment, cfg ApproveTokensForFeeBillingSignerConfig) (cldf.ChangesetOutput, error) {
	state, err := cfg.Validate(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to validate config: %w", err)
	}

	billingSignerPDA, _, err := solState.FindFeeBillingSignerPDA(state.Router)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get billing signer PDA: %w", err)
	}

	err = doApproveTokenTransfer(
		e,
		state,
		cfg.ChainSelector,
		billingSignerPDA,
		cfg.TokenProgram,
		cfg.Amount,
		cfg.Decimals,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to approve token transfer: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}

type ApproveTokensConfig struct {
	ChainSelector    uint64
	AddressToApprove solana.PublicKey
	TokenProgram     solana.PublicKey
	Amount           uint64
	Decimals         uint8
}

func (cfg ApproveTokensConfig) Validate(e cldf.Environment) (solanaStateView.CCIPChainState, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return solanaStateView.CCIPChainState{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainState, found := state.SolChains[cfg.ChainSelector]
	if !found {
		return solanaStateView.CCIPChainState{}, fmt.Errorf("failed to get chain state for selector %d", cfg.ChainSelector)
	}

	if cfg.AddressToApprove == (solana.PublicKey{}) {
		return solanaStateView.CCIPChainState{}, errors.New("address to approve is not set")
	}
	if cfg.TokenProgram == (solana.PublicKey{}) {
		return solanaStateView.CCIPChainState{}, errors.New("token program is not set")
	}

	return chainState, nil
}

func ApproveTokens(e cldf.Environment, cfg ApproveTokensConfig) (cldf.ChangesetOutput, error) {
	state, err := cfg.Validate(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to validate config: %w", err)
	}

	err = doApproveTokenTransfer(
		e,
		state,
		cfg.ChainSelector,
		cfg.AddressToApprove,
		cfg.TokenProgram,
		cfg.Amount,
		cfg.Decimals,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to approve token transfer: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}

func doApproveTokenTransfer(
	e cldf.Environment,
	state solanaStateView.CCIPChainState,
	chainSelector uint64,
	addressToApprove solana.PublicKey,
	tokenPubKey solana.PublicKey,
	amount uint64,
	decimals uint8,
) error {
	solChain := e.SolChains[chainSelector]

	tokenProgram, err := state.TokenToTokenProgram(tokenPubKey)
	if err != nil {
		return fmt.Errorf("failed to get token program: %w", err)
	}

	deployerATA, _, err := solTokenUtil.FindAssociatedTokenAddress(
		tokenProgram,
		tokenPubKey,
		solChain.DeployerKey.PublicKey(),
	)
	if err != nil {
		return fmt.Errorf("failed to find associated token address: %w", err)
	}

	ix, err := solTokenUtil.TokenApproveChecked(
		amount,
		decimals,
		tokenProgram,
		deployerATA,
		tokenPubKey,
		addressToApprove,
		solChain.DeployerKey.PublicKey(),
		solana.PublicKeySlice{},
	)
	if err != nil {
		return fmt.Errorf("failed to TokenApproveChecked: %w", err)
	}

	e.Logger.Infow("Running TokenApprovedChecked (owner ATA = '%s', approved account = '%s', token = '%s', amount = %d, decimals = %d)",
		deployerATA.String(),
		addressToApprove.String(),
		tokenPubKey.String(),
		amount,
		decimals,
	)

	if err = solChain.Confirm([]solana.Instruction{ix}); err != nil {
		e.Logger.Errorw("Failed to confirm instructions for TokenApproveChecked", "chain", solChain.String(), "err", err)
		return err
	}

	if err != nil {
		return fmt.Errorf("failed to create approve instruction: %w", err)
	}

	return nil
}
