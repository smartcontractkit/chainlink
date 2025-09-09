package solana

import (
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	solanaStateView "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var (
	ApproveToken                    = cldf.CreateChangeSet(logicApproveToken, preconditionsApproveToken)
	ApproveTokenForFeeBillingSigner = cldf.CreateChangeSet(logicApproveTokenForFeeBillingSigner, preconditionsApproveTokenForFeeBillingSigner)
)

type ApproveTokenForFeeBillingSignerConfig struct {
	ChainSelector uint64
	TokenProgram  solana.PublicKey
	Amount        uint64
	Decimals      uint8
}

func preconditionsApproveTokenForFeeBillingSigner(e cldf.Environment, cfg ApproveTokenForFeeBillingSignerConfig) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	_, found := state.SolChains[cfg.ChainSelector]
	if !found {
		return fmt.Errorf("failed to get chain state for selector %d", cfg.ChainSelector)
	}
	if cfg.TokenProgram == (solana.PublicKey{}) {
		return errors.New("token program is not set")
	}

	return nil
}

func logicApproveTokenForFeeBillingSigner(e cldf.Environment, cfg ApproveTokenForFeeBillingSignerConfig) (cldf.ChangesetOutput, error) {
	state, _ := stateview.LoadOnchainState(e)

	chainState := state.SolChains[cfg.ChainSelector]

	billingSignerPDA, _, err := solState.FindFeeBillingSignerPDA(chainState.Router)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get billing signer PDA: %w", err)
	}

	err = doApproveTokenTransfer(
		e,
		chainState,
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

type ApproveTokenConfig struct {
	ChainSelector    uint64
	AddressToApprove solana.PublicKey
	TokenProgram     solana.PublicKey
	Amount           uint64
	Decimals         uint8
}

func preconditionsApproveToken(e cldf.Environment, cfg ApproveTokenConfig) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	_, found := state.SolChains[cfg.ChainSelector]
	if !found {
		return fmt.Errorf("failed to get chain state for selector %d", cfg.ChainSelector)
	}

	if cfg.AddressToApprove.Equals(solana.PublicKey{}) {
		return errors.New("address to approve is not set")
	}
	if cfg.TokenProgram.Equals(solana.PublicKey{}) {
		return errors.New("token program is not set")
	}

	return nil
}

func logicApproveToken(e cldf.Environment, cfg ApproveTokenConfig) (cldf.ChangesetOutput, error) {
	state, _ := stateview.LoadOnchainState(e)

	chainState := state.SolChains[cfg.ChainSelector]

	err := doApproveTokenTransfer(
		e,
		chainState,
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
	solChain := e.BlockChains.SolanaChains()[chainSelector]

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

	e.Logger.Infof("Running TokenApprovedChecked (owner ATA = '%s', approved account = '%s', token = '%s', amount = %d, decimals = %d)",
		deployerATA.String(),
		addressToApprove.String(),
		tokenPubKey.String(),
		amount,
		decimals,
	)

	if err = solChain.Confirm([]solana.Instruction{ix}); err != nil {
		e.Logger.Errorf("Failed to confirm instructions TokenApproveChecked for chain %s err %v", solChain.String(), err)
		return err
	}

	if err != nil {
		return fmt.Errorf("failed to create approve instruction: %w", err)
	}

	return nil
}
