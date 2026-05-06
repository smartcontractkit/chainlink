package solana

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/fee_quoter"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// use this changeset to disable a billing token on solana fee quoter
var _ cldf.ChangeSet[DisableBillingTokenConfig] = DisableBillingTokenChangeset

// DisableBillingTokenConfig is the config for disabling a billing token on the Solana FeeQuoter.
// It sets enabled=false on the existing BillingTokenConfigWrapper PDA for the given mint,
// preserving all other fields in the config.
type DisableBillingTokenConfig struct {
	ChainSelector uint64
	TokenMint     solana.PublicKey
	MCMS          *proposalutils.TimelockConfig
}

func (cfg DisableBillingTokenConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	chainState, chainExists := state.SolChains[cfg.ChainSelector]
	if !chainExists {
		return fmt.Errorf("chain %d not found in existing state", cfg.ChainSelector)
	}
	if cfg.TokenMint.IsZero() {
		return errors.New("token mint must be set")
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateFeeQuoterConfig(chain); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.FeeQuoter: true}); err != nil {
		return err
	}

	// verify that the billing token config PDA exists
	billingConfigPDA, _, err := solState.FindFqBillingTokenConfigPDA(cfg.TokenMint, chainState.FeeQuoter)
	if err != nil {
		return fmt.Errorf("failed to find billing token config pda (mint: %s, feeQuoter: %s): %w",
			cfg.TokenMint.String(), chainState.FeeQuoter.String(), err)
	}
	var configAccount solFeeQuoter.BillingTokenConfigWrapper
	if err := chain.GetAccountDataBorshInto(context.Background(), billingConfigPDA, &configAccount); err != nil {
		return fmt.Errorf("billing token config does not exist for mint %s: %w", cfg.TokenMint.String(), err)
	}
	if !configAccount.Config.Enabled {
		return fmt.Errorf("billing token %s is already disabled", cfg.TokenMint.String())
	}

	return nil
}

// DisableBillingTokenChangeset disables a billing token on the Solana FeeQuoter by
// calling UpdateBillingTokenConfig with enabled=false. The existing config fields
// (mint, usd_per_token, premium_multiplier) are preserved from on-chain state.
func DisableBillingTokenChangeset(e cldf.Environment, cfg DisableBillingTokenConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	chainState := state.SolChains[cfg.ChainSelector]

	runSafely(func() {
		solFeeQuoter.SetProgramID(chainState.FeeQuoter)
	})

	// read existing config from chain
	billingConfigPDA, _, _ := solState.FindFqBillingTokenConfigPDA(cfg.TokenMint, chainState.FeeQuoter)
	var existingConfig solFeeQuoter.BillingTokenConfigWrapper
	if err := chain.GetAccountDataBorshInto(context.Background(), billingConfigPDA, &existingConfig); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to read billing token config: %w", err)
	}

	// build the updated config with enabled=false, preserving other fields
	updatedConfig := existingConfig.Config
	updatedConfig.Enabled = false
	// refresh the timestamp
	updatedConfig.UsdPerToken.Timestamp = time.Now().Unix()

	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(chainState.FeeQuoter)
	feeQuoterUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.FeeQuoter,
		solana.PublicKey{},
		"")

	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		shared.FeeQuoter,
		solana.PublicKey{},
		"",
	)

	ix, err := solFeeQuoter.NewUpdateBillingTokenConfigInstruction(
		updatedConfig,
		feeQuoterConfigPDA,
		billingConfigPDA,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build UpdateBillingTokenConfig instruction: %w", err)
	}

	if feeQuoterUsingMCMS {
		tx, err := BuildMCMSTxn(ix, chainState.FeeQuoter.String(), shared.FeeQuoter)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to disable billing token on Solana FeeQuoter", cfg.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm DisableBillingToken: %w", err)
	}

	e.Logger.Infow("Billing token disabled", "chainSelector", cfg.ChainSelector, "tokenMint", cfg.TokenMint.String())
	return cldf.ChangesetOutput{}, nil
}
