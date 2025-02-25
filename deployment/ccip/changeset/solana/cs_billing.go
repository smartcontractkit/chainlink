package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"

	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

var _ deployment.ChangeSet[BillingTokenConfig] = AddBillingTokenChangeset
var _ deployment.ChangeSet[BillingTokenForRemoteChainConfig] = AddBillingTokenForRemoteChain

// ADD BILLING TOKEN
type BillingTokenConfig struct {
	ChainSelector uint64
	TokenPubKey   string
	Config        solFeeQuoter.BillingTokenConfig
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
	// check if already setup
	billingConfigPDA, _, err := solState.FindFqBillingTokenConfigPDA(tokenPubKey, chainState.FeeQuoter)
	if err != nil {
		return fmt.Errorf("failed to find billing token config pda (mint: %s, feeQuoter: %s): %w", tokenPubKey.String(), chainState.FeeQuoter.String(), err)
	}
	var token0ConfigAccount solFeeQuoter.BillingTokenConfigWrapper
	if err := chain.GetAccountDataBorshInto(context.Background(), billingConfigPDA, &token0ConfigAccount); err == nil {
		return fmt.Errorf("billing token config already exists for (mint: %s, feeQuoter: %s)", tokenPubKey.String(), chainState.FeeQuoter.String())
	}
	return nil
}

func AddBillingToken(
	e deployment.Environment,
	chain deployment.SolChain,
	chainState ccipChangeset.SolCCIPChainState,
	billingConfig solFeeQuoter.BillingTokenConfig,
) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(billingConfig.Mint.String())
	tokenBillingPDA, _, _ := solState.FindFqBillingTokenConfigPDA(tokenPubKey, chainState.FeeQuoter)
	billingSignerPDA, _, _ := solState.FindFeeBillingSignerPDA(chainState.Router)
	tokenProgramID, _ := chainState.TokenToTokenProgram(tokenPubKey)
	token2022Receiver, _, _ := solTokenUtil.FindAssociatedTokenAddress(tokenProgramID, tokenPubKey, billingSignerPDA)
	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(chainState.FeeQuoter)
	ixConfig, cerr := solFeeQuoter.NewAddBillingTokenConfigInstruction(
		billingConfig,
		feeQuoterConfigPDA,
		tokenBillingPDA,
		tokenProgramID,
		tokenPubKey,
		token2022Receiver,
		chain.DeployerKey.PublicKey(), // ccip admin
		billingSignerPDA,
		ata.ProgramID,
		solana.SystemProgramID,
	).ValidateAndBuild()
	if cerr != nil {
		return fmt.Errorf("failed to generate instructions: %w", cerr)
	}
	instructions := []solana.Instruction{ixConfig}
	if err := chain.Confirm(instructions); err != nil {
		return fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return nil
}

func AddBillingTokenChangeset(e deployment.Environment, cfg BillingTokenConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]

	solFeeQuoter.SetProgramID(chainState.FeeQuoter)

	if err := AddBillingToken(e, chain, chainState, cfg.Config); err != nil {
		return deployment.ChangesetOutput{}, err
	}

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
	return deployment.ChangesetOutput{}, nil
}

// ADD BILLING TOKEN FOR REMOTE CHAIN
type BillingTokenForRemoteChainConfig struct {
	ChainSelector       uint64
	RemoteChainSelector uint64
	Config              solFeeQuoter.TokenTransferFeeConfig
	TokenPubKey         string
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

func AddBillingTokenForRemoteChain(e deployment.Environment, cfg BillingTokenForRemoteChainConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	remoteBillingPDA, _, _ := solState.FindFqPerChainPerTokenConfigPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.FeeQuoter)

	ix, err := solFeeQuoter.NewSetTokenTransferFeeConfigInstruction(
		cfg.RemoteChainSelector,
		tokenPubKey,
		cfg.Config,
		chainState.FeeQuoterConfigPDA,
		remoteBillingPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}
	instructions := []solana.Instruction{ix}
	if err := chain.Confirm(instructions); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
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
	return deployment.ChangesetOutput{}, nil
}
