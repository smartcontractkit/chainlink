package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

// ADD BILLING TOKEN
type BillingTokenConfig struct {
	ChainSelector uint64
	TokenPubKey   string
	Config        solFeeQuoter.BillingTokenConfig
	// We have different instructions for add vs update, so we need to know which one to use
	IsUpdate   bool
	MCMSSolana *MCMSConfigSolana
}

func (cfg BillingTokenConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}

	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	if err := validateFeeQuoterConfig(chain, chainState); err != nil {
		return err
	}
	if _, err := chainState.TokenToTokenProgram(tokenPubKey); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.ChainSelector, cfg.MCMSSolana); err != nil {
		return err
	}
	feeQuoterUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.FeeQuoterOwnedByTimelock
	if err := ccipChangeset.ValidateOwnershipSolana(&e, chain, feeQuoterUsingMCMS, chainState.FeeQuoter, ccipChangeset.FeeQuoter); err != nil {
		return fmt.Errorf("failed to validate ownership: %w", err)
	}
	// check if already setup
	billingConfigPDA, _, err := solState.FindFqBillingTokenConfigPDA(tokenPubKey, chainState.FeeQuoter)
	if err != nil {
		return fmt.Errorf("failed to find billing token config pda (mint: %s, feeQuoter: %s): %w", tokenPubKey.String(), chainState.FeeQuoter.String(), err)
	}
	if !cfg.IsUpdate {
		var token0ConfigAccount solFeeQuoter.BillingTokenConfigWrapper
		if err := chain.GetAccountDataBorshInto(context.Background(), billingConfigPDA, &token0ConfigAccount); err == nil {
			return fmt.Errorf("billing token config already exists for (mint: %s, feeQuoter: %s)", tokenPubKey.String(), chainState.FeeQuoter.String())
		}
	}
	return nil
}

func AddBillingToken(
	e deployment.Environment,
	chain deployment.SolChain,
	chainState ccipChangeset.SolCCIPChainState,
	billingTokenConfig solFeeQuoter.BillingTokenConfig,
	mcms *MCMSConfigSolana,
	isUpdate bool,
) ([]mcmsTypes.Transaction, error) {
	txns := make([]mcmsTypes.Transaction, 0)
	tokenPubKey := solana.MustPublicKeyFromBase58(billingTokenConfig.Mint.String())
	tokenBillingPDA, _, _ := solState.FindFqBillingTokenConfigPDA(tokenPubKey, chainState.FeeQuoter)
	// we dont need to handle test router here because we explicitly create this and token2022Receiver for test router
	billingSignerPDA, _, _ := solState.FindFeeBillingSignerPDA(chainState.Router)
	tokenProgramID, _ := chainState.TokenToTokenProgram(tokenPubKey)
	token2022Receiver, _, _ := solTokenUtil.FindAssociatedTokenAddress(tokenProgramID, tokenPubKey, billingSignerPDA)
	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(chainState.FeeQuoter)
	feeQuoterUsingMCMS := mcms != nil && mcms.FeeQuoterOwnedByTimelock
	timelockSigner, err := FetchTimelockSigner(e, chain.Selector)
	if err != nil {
		return txns, fmt.Errorf("failed to fetch timelock signer: %w", err)
	}
	var authority solana.PublicKey
	if feeQuoterUsingMCMS {
		authority = timelockSigner
	} else {
		authority = chain.DeployerKey.PublicKey()
	}
	var ixConfig solana.Instruction
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
			token2022Receiver,
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
		tx, err := BuildMCMSTxn(ixConfig, chainState.FeeQuoter.String(), ccipChangeset.FeeQuoter)
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

func AddBillingTokenChangeset(e deployment.Environment, cfg BillingTokenConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]

	solFeeQuoter.SetProgramID(chainState.FeeQuoter)

	txns, err := AddBillingToken(e, chain, chainState, cfg.Config, cfg.MCMSSolana, cfg.IsUpdate)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	if !cfg.IsUpdate {
		tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
		tokenBillingPDA, _, _ := solState.FindFqBillingTokenConfigPDA(tokenPubKey, chainState.FeeQuoter)

		addressLookupTable, err := ccipChangeset.FetchOfframpLookupTable(e.GetContext(), chain, chainState.OffRamp)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get offramp reference addresses: %w", err)
		}

		if err := solCommonUtil.ExtendLookupTable(
			e.GetContext(),
			chain.Client,
			addressLookupTable,
			*chain.DeployerKey,
			[]solana.PublicKey{tokenBillingPDA},
		); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to extend lookup table: %w", err)
		}
		e.Logger.Infow("Billing token added", "chainSelector", cfg.ChainSelector, "tokenPubKey", tokenPubKey.String())
	}

	// create proposals for ixns
	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to add billing token to Solana", cfg.MCMSSolana.MCMS.MinDelay, txns)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return deployment.ChangesetOutput{}, nil
}

// ADD BILLING TOKEN FOR REMOTE CHAIN
type BillingTokenForRemoteChainConfig struct {
	ChainSelector       uint64
	RemoteChainSelector uint64
	Config              solFeeQuoter.TokenTransferFeeConfig
	TokenPubKey         string
	MCMSSolana          *MCMSConfigSolana
}

func (cfg BillingTokenForRemoteChainConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateFeeQuoterConfig(chain, chainState); err != nil {
		return fmt.Errorf("fee quoter validation failed: %w", err)
	}
	// check if desired state already exists
	remoteBillingPDA, _, err := solState.FindFqPerChainPerTokenConfigPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.FeeQuoter)
	if err != nil {
		return fmt.Errorf("failed to find remote billing token config pda for (remoteSelector: %d, mint: %s, feeQuoter: %s): %w", cfg.RemoteChainSelector, tokenPubKey.String(), chainState.FeeQuoter.String(), err)
	}
	var remoteBillingAccount solFeeQuoter.PerChainPerTokenConfig
	if err := chain.GetAccountDataBorshInto(context.Background(), remoteBillingPDA, &remoteBillingAccount); err == nil {
		return fmt.Errorf("billing token config already exists for (remoteSelector: %d, mint: %s, feeQuoter: %s)", cfg.RemoteChainSelector, tokenPubKey.String(), chainState.FeeQuoter.String())
	}
	return nil
}

// TODO: rename this, i dont think this is for billing, this is more for token transfer config/fees
func AddBillingTokenForRemoteChain(e deployment.Environment, cfg BillingTokenForRemoteChainConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	remoteBillingPDA, _, _ := solState.FindFqPerChainPerTokenConfigPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.FeeQuoter)

	if err := ValidateMCMSConfigSolana(e, cfg.ChainSelector, cfg.MCMSSolana); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	feeQuoterUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.FeeQuoterOwnedByTimelock
	if err := ccipChangeset.ValidateOwnershipSolana(&e, chain, feeQuoterUsingMCMS, chainState.FeeQuoter, ccipChangeset.FeeQuoter); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to validate ownership: %w", err)
	}
	timelockSigner, err := FetchTimelockSigner(e, chain.Selector)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
	}

	var authority solana.PublicKey
	if feeQuoterUsingMCMS {
		authority = timelockSigner
	} else {
		authority = chain.DeployerKey.PublicKey()
	}
	ix, err := solFeeQuoter.NewSetTokenTransferFeeConfigInstruction(
		cfg.RemoteChainSelector,
		tokenPubKey,
		cfg.Config,
		chainState.FeeQuoterConfigPDA,
		remoteBillingPDA,
		authority,
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}
	if !feeQuoterUsingMCMS {
		if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
		}
	}

	addressLookupTable, err := ccipChangeset.FetchOfframpLookupTable(e.GetContext(), chain, chainState.OffRamp)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get offramp reference addresses: %w", err)
	}

	if err := solCommonUtil.ExtendLookupTable(
		e.GetContext(),
		chain.Client,
		addressLookupTable,
		*chain.DeployerKey,
		[]solana.PublicKey{remoteBillingPDA},
	); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to extend lookup table: %w", err)
	}

	e.Logger.Infow("Token billing set for remote chain", "chainSelector ", cfg.ChainSelector, "remoteChainSelector ", cfg.RemoteChainSelector, "tokenPubKey", tokenPubKey.String())

	if feeQuoterUsingMCMS {
		tx, err := BuildMCMSTxn(ix, chainState.FeeQuoter.String(), ccipChangeset.FeeQuoter)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to set billing token for remote chain to Solana", cfg.MCMSSolana.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return deployment.ChangesetOutput{}, nil
}
