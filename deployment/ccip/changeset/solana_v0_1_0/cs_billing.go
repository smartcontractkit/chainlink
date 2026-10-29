package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_0/fee_quoter"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
)

// ADD BILLING TOKEN
type BillingTokenConfig struct {
	ChainSelector uint64
	Config        solFeeQuoter.BillingTokenConfig
	MCMS          *cldfproposalutils.TimelockConfig

	// inferred from state
	IsUpdate bool
}

func (cfg *BillingTokenConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	tokenPubKey := cfg.Config.Mint
	chainState := state.SolChains[cfg.ChainSelector]
	if err := chainState.CommonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateFeeQuoterConfig(chain); err != nil {
		return err
	}
	if _, err := chainState.TokenToTokenProgram(tokenPubKey); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.FeeQuoter: true}); err != nil {
		return err
	}
	// check if already setup
	billingConfigPDA, _, err := solState.FindFqBillingTokenConfigPDA(tokenPubKey, chainState.FeeQuoter)
	if err != nil {
		return fmt.Errorf("failed to find billing token config pda (mint: %s, feeQuoter: %s): %w", tokenPubKey.String(), chainState.FeeQuoter.String(), err)
	}
	var token0ConfigAccount solFeeQuoter.BillingTokenConfigWrapper
	if err := chain.GetAccountDataBorshInto(context.Background(), billingConfigPDA, &token0ConfigAccount); err == nil {
		e.Logger.Infow("Billing token already exists. Configuring as update", "chainSelector", cfg.ChainSelector, "tokenPubKey", tokenPubKey.String())
		cfg.IsUpdate = true
	}
	return nil
}

func AddBillingToken(
	e cldf.Environment,
	chain cldf_solana.Chain,
	chainState solanastateview.CCIPChainState,
	billingTokenConfig solFeeQuoter.BillingTokenConfig,
	mcms *cldfproposalutils.TimelockConfig,
	isUpdate bool,
	feeQuoterAddress solana.PublicKey,
	routerAddress solana.PublicKey,
) ([]mcmsTypes.Transaction, error) {
	txns := make([]mcmsTypes.Transaction, 0)
	tokenPubKey := billingTokenConfig.Mint
	tokenBillingPDA, _, _ := solState.FindFqBillingTokenConfigPDA(tokenPubKey, feeQuoterAddress)
	// we dont need to handle test router here because we explicitly create this and token Receiver for test router
	billingSignerPDA, _, _ := solState.FindFeeBillingSignerPDA(routerAddress)
	tokenProgramID, _ := chainState.TokenToTokenProgram(tokenPubKey)
	tokenReceiver, _, _ := solTokenUtil.FindAssociatedTokenAddress(tokenProgramID, tokenPubKey, billingSignerPDA)
	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(feeQuoterAddress)
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
	var ixConfig solana.Instruction
	var err error
	if isUpdate {
		ixConfig, err = solFeeQuoter.NewUpdateBillingTokenConfigInstruction(
			billingTokenConfig,
			feeQuoterConfigPDA,
			tokenBillingPDA,
			authority,
		).ValidateAndBuild()
	} else {
		ixConfig, err = solFeeQuoter.NewAddBillingTokenConfigInstruction(
			billingTokenConfig,
			feeQuoterConfigPDA,
			tokenBillingPDA,
			tokenProgramID,
			tokenPubKey,
			tokenReceiver,
			authority, // ccip admin
			billingSignerPDA,
			ata.ProgramID,
			solana.SystemProgramID,
		).ValidateAndBuild()
	}
	if err != nil {
		return txns, fmt.Errorf("failed to generate instructions: %w", err)
	}
	if feeQuoterUsingMCMS {
		tx, err := BuildMCMSTxn(ixConfig, chainState.FeeQuoter.String(), shared.FeeQuoter)
		if err != nil {
			return txns, fmt.Errorf("failed to create transaction: %w", err)
		}
		txns = append(txns, *tx)
	} else {
		if err := chain.Confirm([]solana.Instruction{ixConfig}); err != nil {
			return txns, fmt.Errorf("failed to confirm instructions: %w", err)
		}
	}

	return txns, nil
}

// ADD BILLING TOKEN FOR REMOTE CHAIN
type TokenTransferFeeForRemoteChainConfig struct {
	ChainSelector       uint64
	RemoteChainSelector uint64
	// need to provide complete config, onchain does not do an upsert, it does a overwrite
	Config      solFeeQuoter.TokenTransferFeeConfig
	TokenPubKey solana.PublicKey
	MCMS        *cldfproposalutils.TimelockConfig
}

const MinDestBytesOverhead = 32

func (cfg TokenTransferFeeForRemoteChainConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	tokenPubKey := cfg.TokenPubKey
	chainState := state.SolChains[cfg.ChainSelector]
	if err := chainState.CommonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateFeeQuoterConfig(chain); err != nil {
		return fmt.Errorf("fee quoter validation failed: %w", err)
	}
	if cfg.Config.DestBytesOverhead < 32 {
		e.Logger.Infow("dest bytes overhead is less than minimum. Setting to minimum value",
			"destBytesOverhead", cfg.Config.DestBytesOverhead,
			"minDestBytesOverhead", MinDestBytesOverhead)
		cfg.Config.DestBytesOverhead = MinDestBytesOverhead
	}
	if cfg.Config.MinFeeUsdcents > cfg.Config.MaxFeeUsdcents {
		return fmt.Errorf("min fee %d cannot be greater than max fee %d", cfg.Config.MinFeeUsdcents, cfg.Config.MaxFeeUsdcents)
	}

	return ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.FeeQuoter: true})
}
