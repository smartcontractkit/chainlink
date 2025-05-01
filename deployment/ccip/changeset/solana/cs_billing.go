package solana

import (
	"context"
	"fmt"

	solBinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

// use this changeset to add a billing token to solana
var _ deployment.ChangeSet[BillingTokenConfig] = AddBillingTokenChangeset

// use this changeset to add a token transfer fee for a remote chain to solana
var _ deployment.ChangeSet[TokenTransferFeeForRemoteChainConfig] = AddTokenTransferFeeForRemoteChain

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
	if err := ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, tokenPubKey); err != nil {
		return err
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
	feeQuoterAddress solana.PublicKey,
	routerAddress solana.PublicKey,
) ([]mcmsTypes.Transaction, error) {
	txns := make([]mcmsTypes.Transaction, 0)
	tokenPubKey := solana.MustPublicKeyFromBase58(billingTokenConfig.Mint.String())
	tokenBillingPDA, _, _ := solState.FindFqBillingTokenConfigPDA(tokenPubKey, feeQuoterAddress)
	// we dont need to handle test router here because we explicitly create this and token2022Receiver for test router
	billingSignerPDA, _, _ := solState.FindFeeBillingSignerPDA(routerAddress)
	tokenProgramID, _ := chainState.TokenToTokenProgram(tokenPubKey)
	token2022Receiver, _, _ := solTokenUtil.FindAssociatedTokenAddress(tokenProgramID, tokenPubKey, billingSignerPDA)
	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(feeQuoterAddress)
	feeQuoterUsingMCMS := mcms != nil && mcms.FeeQuoterOwnedByTimelock

	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		mcms,
		ccipChangeset.FeeQuoter,
		solana.PublicKey{})
	if err != nil {
		return txns, fmt.Errorf("failed to get authority for ixn: %w", err)
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

	txns, err := AddBillingToken(e, chain, chainState, cfg.Config, cfg.MCMSSolana, cfg.IsUpdate, chainState.FeeQuoter, chainState.Router)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	if !cfg.IsUpdate {
		tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
		tokenBillingPDA, _, _ := solState.FindFqBillingTokenConfigPDA(tokenPubKey, chainState.FeeQuoter)
		if err := extendLookupTable(e, chain, chainState.OffRamp, []solana.PublicKey{tokenBillingPDA}); err != nil {
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
type TokenTransferFeeForRemoteChainConfig struct {
	ChainSelector       uint64
	RemoteChainSelector uint64
	Config              solFeeQuoter.TokenTransferFeeConfig
	TokenPubKey         string
	MCMSSolana          *MCMSConfigSolana
}

func (cfg TokenTransferFeeForRemoteChainConfig) Validate(e deployment.Environment) error {
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

	return nil
}

// TODO: rename this, i dont think this is for billing, this is more for token transfer config/fees
func AddTokenTransferFeeForRemoteChain(e deployment.Environment, cfg TokenTransferFeeForRemoteChainConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	remoteBillingPDA, _, _ := solState.FindFqPerChainPerTokenConfigPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.FeeQuoter)
	feeQuoterUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.FeeQuoterOwnedByTimelock

	if err := ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, tokenPubKey); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.FeeQuoter,
		solana.PublicKey{})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
	solFeeQuoter.SetProgramID(chainState.FeeQuoter)
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
	if err := extendLookupTable(e, chain, chainState.OffRamp, []solana.PublicKey{remoteBillingPDA}); err != nil {
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

// Price Update Changesets are in case of emergency as normally offramp will call this as part of normal operations
type UpdatePricesConfig struct {
	ChainSelector     uint64
	TokenPriceUpdates []solFeeQuoter.TokenPriceUpdate
	GasPriceUpdates   []solFeeQuoter.GasPriceUpdate
	PriceUpdater      solana.PublicKey
	MCMSSolana        *MCMSConfigSolana
}

func (cfg UpdatePricesConfig) Validate(e deployment.Environment) error {
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateFeeQuoterConfig(chain, chainState); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, solana.PublicKey{}); err != nil {
		return err
	}
	if cfg.PriceUpdater.IsZero() {
		return fmt.Errorf("price updater is zero for chain %d", cfg.ChainSelector)
	}
	for _, update := range cfg.TokenPriceUpdates {
		billingConfigPDA, _, _ := solState.FindFqBillingTokenConfigPDA(update.SourceToken, chainState.FeeQuoter)
		var token0ConfigAccount solFeeQuoter.BillingTokenConfigWrapper
		err = chain.GetAccountDataBorshInto(e.GetContext(), billingConfigPDA, &token0ConfigAccount)
		if err != nil {
			return fmt.Errorf("failed to find billing token config for (mint: %s, feeQuoter: %s): %w", update.SourceToken.String(), chainState.FeeQuoter.String(), err)
		}
	}
	for _, update := range cfg.GasPriceUpdates {
		fqDestPDA, _, _ := solState.FindFqDestChainPDA(update.DestChainSelector, chainState.FeeQuoter)
		var destChainConfig solFeeQuoter.DestChainConfig
		err = chain.GetAccountDataBorshInto(e.GetContext(), fqDestPDA, &destChainConfig)
		if err != nil {
			return fmt.Errorf("failed to find dest chain config for (destSelector: %d, feeQuoter: %s): %w", update.DestChainSelector, chainState.FeeQuoter.String(), err)
		}
	}

	return nil
}

func UpdatePrices(e deployment.Environment, cfg UpdatePricesConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	s, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chainSel := cfg.ChainSelector
	chain := e.SolChains[chainSel]
	feeQuoterID := s.SolChains[chainSel].FeeQuoter
	feeQuoterUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.FeeQuoterOwnedByTimelock

	// verified while loading state
	fqAllowedPriceUpdaterPDA, _, _ := solState.FindFqAllowedPriceUpdaterPDA(cfg.PriceUpdater, feeQuoterID)
	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(feeQuoterID)

	solFeeQuoter.SetProgramID(feeQuoterID)
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.FeeQuoter,
		solana.PublicKey{})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
	raw := solFeeQuoter.NewUpdatePricesInstruction(
		cfg.TokenPriceUpdates,
		cfg.GasPriceUpdates,
		authority,
		fqAllowedPriceUpdaterPDA,
		feeQuoterConfigPDA,
	)
	for _, update := range cfg.TokenPriceUpdates {
		billingTokenConfigPDA, _, _ := solState.FindFqBillingTokenConfigPDA(update.SourceToken, feeQuoterID)
		raw.AccountMetaSlice.Append(solana.Meta(billingTokenConfigPDA).WRITE())
	}
	for _, update := range cfg.GasPriceUpdates {
		fqDestPDA, _, _ := solState.FindFqDestChainPDA(update.DestChainSelector, feeQuoterID)
		raw.AccountMetaSlice.Append(solana.Meta(fqDestPDA).WRITE())
	}
	ix, err := raw.ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}

	if feeQuoterUsingMCMS {
		tx, err := BuildMCMSTxn(ix, s.SolChains[chainSel].FeeQuoter.String(), ccipChangeset.FeeQuoter)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to NewUpdatePricesInstruction in Solana", cfg.MCMSSolana.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	if err = chain.Confirm([]solana.Instruction{ix}); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}

type ModifyPriceUpdaterConfig struct {
	ChainSelector      uint64
	PriceUpdater       solana.PublicKey
	PriceUpdaterAction PriceUpdaterAction
	MCMSSolana         *MCMSConfigSolana
}

type PriceUpdaterAction int

const (
	AddUpdater PriceUpdaterAction = iota
	RemoveUpdater
)

func (cfg ModifyPriceUpdaterConfig) Validate(e deployment.Environment) error {
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateFeeQuoterConfig(chain, chainState); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, solana.PublicKey{}); err != nil {
		return err
	}
	if cfg.PriceUpdater.IsZero() {
		return fmt.Errorf("price updater is zero for chain %d", cfg.ChainSelector)
	}
	return nil
}

func ModifyPriceUpdater(e deployment.Environment, cfg ModifyPriceUpdaterConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	s, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chainSel := cfg.ChainSelector
	chain := e.SolChains[chainSel]
	feeQuoterID := s.SolChains[chainSel].FeeQuoter
	feeQuoterUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.FeeQuoterOwnedByTimelock

	// verified while loading state
	fqAllowedPriceUpdaterPDA, _, _ := solState.FindFqAllowedPriceUpdaterPDA(cfg.PriceUpdater, feeQuoterID)
	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(feeQuoterID)

	solFeeQuoter.SetProgramID(feeQuoterID)
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.FeeQuoter,
		solana.PublicKey{})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
	var ix solana.Instruction
	if cfg.PriceUpdaterAction == AddUpdater {
		ix, err = solFeeQuoter.NewAddPriceUpdaterInstruction(
			cfg.PriceUpdater,
			fqAllowedPriceUpdaterPDA,
			feeQuoterConfigPDA,
			authority,
			solana.SystemProgramID,
		).ValidateAndBuild()
	} else {
		ix, err = solFeeQuoter.NewRemovePriceUpdaterInstruction(
			cfg.PriceUpdater,
			fqAllowedPriceUpdaterPDA,
			feeQuoterConfigPDA,
			authority,
			solana.SystemProgramID,
		).ValidateAndBuild()
	}
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}

	if feeQuoterUsingMCMS {
		tx, err := BuildMCMSTxn(ix, s.SolChains[chainSel].FeeQuoter.String(), ccipChangeset.FeeQuoter)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to NewUpdatePricesInstruction in Solana", cfg.MCMSSolana.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	if err = chain.Confirm([]solana.Instruction{ix}); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}

type WithdrawBilledFundsConfig struct {
	ChainSelector uint64
	TransferAll   bool
	Amount        uint64
	TokenPubKey   string
	MCMSSolana    *MCMSConfigSolana
}

func (cfg WithdrawBilledFundsConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}
	if err := validateFeeAggregatorConfig(chain, chainState); err != nil {
		return err
	}
	return ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, tokenPubKey)
}

func WithdrawBilledFunds(e deployment.Environment, cfg WithdrawBilledFundsConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	s, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chainSel := cfg.ChainSelector
	chain := e.SolChains[chainSel]
	chainState := s.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	billingSignerPDA, _, _ := solState.FindFeeBillingSignerPDA(chainState.Router)
	tokenProgramID, _ := chainState.TokenToTokenProgram(tokenPubKey)
	tokenReceiverPDA, _, _ := solTokenUtil.FindAssociatedTokenAddress(tokenProgramID, tokenPubKey, billingSignerPDA)
	feeAggregatorATA, _, _ := solTokenUtil.FindAssociatedTokenAddress(tokenProgramID, tokenPubKey, chainState.GetFeeAggregator(chain))
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)
	routerUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock

	solRouter.SetProgramID(chainState.Router)
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.Router,
		solana.PublicKey{})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
	ix, err := solRouter.NewWithdrawBilledFundsInstruction(
		cfg.TransferAll,
		cfg.Amount,
		tokenPubKey,
		tokenReceiverPDA,
		feeAggregatorATA,
		tokenProgramID,
		billingSignerPDA,
		routerConfigPDA,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}

	if routerUsingMCMS {
		tx, err := BuildMCMSTxn(ix, s.SolChains[chainSel].Router.String(), ccipChangeset.Router)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to WithdrawBilledFunds in Solana", cfg.MCMSSolana.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	if err = chain.Confirm([]solana.Instruction{ix}); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}

type SetMaxFeeJuelsPerMsgConfig struct {
	ChainSelector     uint64
	MaxFeeJuelsPerMsg solBinary.Uint128
	MCMSSolana        *MCMSConfigSolana
}

func (cfg SetMaxFeeJuelsPerMsgConfig) Validate(e deployment.Environment) error {
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState, chainExists := state.SolChains[cfg.ChainSelector]
	if !chainExists {
		return fmt.Errorf("chain %d not found in existing state", cfg.ChainSelector)
	}
	chain := e.SolChains[cfg.ChainSelector]

	if err := validateFeeQuoterConfig(chain, chainState); err != nil {
		return err
	}

	return ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, solana.PublicKey{})
}

func SetMaxFeeJuelsPerMsg(e deployment.Environment, cfg SetMaxFeeJuelsPerMsgConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]

	fqConfig, _, _ := solState.FindConfigPDA(chainState.FeeQuoter)
	fqUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.FeeQuoterOwnedByTimelock

	solFeeQuoter.SetProgramID(chainState.FeeQuoter)
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.FeeQuoter,
		solana.PublicKey{})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
	instruction, err := solFeeQuoter.NewSetMaxFeeJuelsPerMsgInstruction(
		cfg.MaxFeeJuelsPerMsg,
		fqConfig,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if fqUsingMCMS {
		tx, err := BuildMCMSTxn(instruction, chainState.FeeQuoter.String(), ccipChangeset.FeeQuoter)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to SetMaxFeeJuelsPerMsg in Solana", cfg.MCMSSolana.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	if err := chain.Confirm([]solana.Instruction{instruction}); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}

	return deployment.ChangesetOutput{}, nil
}
