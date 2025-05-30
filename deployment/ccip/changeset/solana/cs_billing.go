package solana

import (
	"context"
	"fmt"

	solBinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// use this changeset to add a billing token to solana
var _ cldf.ChangeSet[BillingTokenConfig] = AddBillingTokenChangeset

// use this changeset to add a token transfer fee for a remote chain to solana
var _ cldf.ChangeSet[TokenTransferFeeForRemoteChainConfig] = AddTokenTransferFeeForRemoteChain

// ADD BILLING TOKEN
type BillingTokenConfig struct {
	ChainSelector uint64
	TokenPubKey   string
	Config        solFeeQuoter.BillingTokenConfig
	// We have different instructions for add vs update, so we need to know which one to use
	IsUpdate bool
	MCMS     *proposalutils.TimelockConfig
}

func (cfg *BillingTokenConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
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
	mcms *proposalutils.TimelockConfig,
	isUpdate bool,
	feeQuoterAddress solana.PublicKey,
	routerAddress solana.PublicKey,
) ([]mcmsTypes.Transaction, error) {
	txns := make([]mcmsTypes.Transaction, 0)
	tokenPubKey := solana.MustPublicKeyFromBase58(billingTokenConfig.Mint.String())
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
		mcms,
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

func AddBillingTokenChangeset(e cldf.Environment, cfg BillingTokenConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	chainState := state.SolChains[cfg.ChainSelector]

	solFeeQuoter.SetProgramID(chainState.FeeQuoter)

	txns, err := AddBillingToken(e, chain, chainState, cfg.Config, cfg.MCMS, cfg.IsUpdate, chainState.FeeQuoter, chainState.Router)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	tokenBillingPDA, _, _ := solState.FindFqBillingTokenConfigPDA(tokenPubKey, chainState.FeeQuoter)
	if err := extendLookupTable(e, chain, chainState.OffRamp, []solana.PublicKey{tokenBillingPDA}); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to extend lookup table: %w", err)
	}
	e.Logger.Infow("Billing token added", "chainSelector", cfg.ChainSelector, "tokenPubKey", tokenPubKey.String())

	// create proposals for ixns
	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to add billing token to Solana", cfg.MCMS.MinDelay, txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}

// ADD BILLING TOKEN FOR REMOTE CHAIN
type TokenTransferFeeForRemoteChainConfig struct {
	ChainSelector       uint64
	RemoteChainSelector uint64
	Config              solFeeQuoter.TokenTransferFeeConfig
	TokenPubKey         string
	MCMS                *proposalutils.TimelockConfig
}

const MinDestBytesOverhead = 32

func (cfg TokenTransferFeeForRemoteChainConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
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

// TODO: rename this, i dont think this is for billing, this is more for token transfer config/fees
func AddTokenTransferFeeForRemoteChain(e cldf.Environment, cfg TokenTransferFeeForRemoteChainConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	remoteBillingPDA, _, _ := solState.FindFqPerChainPerTokenConfigPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.FeeQuoter)
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
		cfg.MCMS,
		shared.FeeQuoter,
		solana.PublicKey{},
		"")
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
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}
	if !feeQuoterUsingMCMS {
		if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
		}
	}
	if err := extendLookupTable(e, chain, chainState.OffRamp, []solana.PublicKey{remoteBillingPDA}); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to extend lookup table: %w", err)
	}

	e.Logger.Infow("Token billing set for remote chain", "chainSelector ", cfg.ChainSelector, "remoteChainSelector ", cfg.RemoteChainSelector, "tokenPubKey", tokenPubKey.String())

	if feeQuoterUsingMCMS {
		tx, err := BuildMCMSTxn(ix, chainState.FeeQuoter.String(), shared.FeeQuoter)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to set billing token for remote chain to Solana", cfg.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}

// Price Update Changesets are in case of emergency as normally offramp will call this as part of normal operations
type UpdatePricesConfig struct {
	ChainSelector     uint64
	TokenPriceUpdates []solFeeQuoter.TokenPriceUpdate
	GasPriceUpdates   []solFeeQuoter.GasPriceUpdate
	PriceUpdater      solana.PublicKey
	MCMS              *proposalutils.TimelockConfig
}

func (cfg UpdatePricesConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateFeeQuoterConfig(chain); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.FeeQuoter: true}); err != nil {
		return err
	}
	if cfg.PriceUpdater.IsZero() {
		return fmt.Errorf("price updater is zero for chain %d", cfg.ChainSelector)
	}
	var err error
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

func UpdatePrices(e cldf.Environment, cfg UpdatePricesConfig) (cldf.ChangesetOutput, error) {
	s, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, s); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	chainSel := cfg.ChainSelector
	chain := e.BlockChains.SolanaChains()[chainSel]
	chainState := s.SolChains[chainSel]
	feeQuoterID := s.SolChains[chainSel].FeeQuoter
	feeQuoterUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.FeeQuoter,
		solana.PublicKey{},
		"")

	// verified while loading state
	fqAllowedPriceUpdaterPDA, _, _ := solState.FindFqAllowedPriceUpdaterPDA(cfg.PriceUpdater, feeQuoterID)
	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(feeQuoterID)

	solFeeQuoter.SetProgramID(feeQuoterID)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.FeeQuoter,
		solana.PublicKey{},
		"")
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
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}

	if feeQuoterUsingMCMS {
		tx, err := BuildMCMSTxn(ix, s.SolChains[chainSel].FeeQuoter.String(), shared.FeeQuoter)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to NewUpdatePricesInstruction in Solana", cfg.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	if err = chain.Confirm([]solana.Instruction{ix}); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	return cldf.ChangesetOutput{}, nil
}

type ModifyPriceUpdaterConfig struct {
	ChainSelector      uint64
	PriceUpdater       solana.PublicKey
	PriceUpdaterAction PriceUpdaterAction
	MCMS               *proposalutils.TimelockConfig
}

type PriceUpdaterAction int

const (
	AddUpdater PriceUpdaterAction = iota
	RemoveUpdater
)

func (cfg ModifyPriceUpdaterConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateFeeQuoterConfig(chain); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.FeeQuoter: true}); err != nil {
		return err
	}
	if cfg.PriceUpdater.IsZero() {
		return fmt.Errorf("price updater is zero for chain %d", cfg.ChainSelector)
	}
	return nil
}

func ModifyPriceUpdater(e cldf.Environment, cfg ModifyPriceUpdaterConfig) (cldf.ChangesetOutput, error) {
	s, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, s); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	chainSel := cfg.ChainSelector
	chain := e.BlockChains.SolanaChains()[chainSel]
	chainState := s.SolChains[chainSel]
	feeQuoterID := s.SolChains[chainSel].FeeQuoter
	feeQuoterUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.FeeQuoter,
		solana.PublicKey{},
		"")

	// verified while loading state
	fqAllowedPriceUpdaterPDA, _, _ := solState.FindFqAllowedPriceUpdaterPDA(cfg.PriceUpdater, feeQuoterID)
	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(feeQuoterID)

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
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}

	if feeQuoterUsingMCMS {
		tx, err := BuildMCMSTxn(ix, s.SolChains[chainSel].FeeQuoter.String(), shared.FeeQuoter)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to NewUpdatePricesInstruction in Solana", cfg.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	if err = chain.Confirm([]solana.Instruction{ix}); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	return cldf.ChangesetOutput{}, nil
}

type WithdrawBilledFundsConfig struct {
	ChainSelector uint64
	TransferAll   bool
	Amount        uint64
	TokenPubKey   string
	MCMS          *proposalutils.TimelockConfig
}

func (cfg WithdrawBilledFundsConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	chainState := state.SolChains[cfg.ChainSelector]
	if err := chainState.CommonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateRouterConfig(chain); err != nil {
		return err
	}
	if err := chainState.ValidateFeeAggregatorConfig(chain); err != nil {
		return err
	}
	return ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.Router: true})
}

func WithdrawBilledFunds(e cldf.Environment, cfg WithdrawBilledFundsConfig) (cldf.ChangesetOutput, error) {
	s, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, s); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	chainSel := cfg.ChainSelector
	chain := e.BlockChains.SolanaChains()[chainSel]
	chainState := s.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	billingSignerPDA, _, _ := solState.FindFeeBillingSignerPDA(chainState.Router)
	tokenProgramID, _ := chainState.TokenToTokenProgram(tokenPubKey)
	tokenReceiverPDA, _, _ := solTokenUtil.FindAssociatedTokenAddress(tokenProgramID, tokenPubKey, billingSignerPDA)
	feeAggregatorATA, _, _ := solTokenUtil.FindAssociatedTokenAddress(tokenProgramID, tokenPubKey, chainState.GetFeeAggregator(chain))
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)
	routerUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.Router,
		solana.PublicKey{},
		"")

	solRouter.SetProgramID(chainState.Router)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.Router,
		solana.PublicKey{},
		"",
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
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
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}

	if routerUsingMCMS {
		tx, err := BuildMCMSTxn(ix, s.SolChains[chainSel].Router.String(), shared.Router)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to WithdrawBilledFunds in Solana", cfg.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	if err = chain.Confirm([]solana.Instruction{ix}); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	return cldf.ChangesetOutput{}, nil
}

type SetMaxFeeJuelsPerMsgConfig struct {
	ChainSelector     uint64
	MaxFeeJuelsPerMsg solBinary.Uint128
	MCMS              *proposalutils.TimelockConfig
}

func (cfg SetMaxFeeJuelsPerMsgConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	chainState, chainExists := state.SolChains[cfg.ChainSelector]
	if !chainExists {
		return fmt.Errorf("chain %d not found in existing state", cfg.ChainSelector)
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]

	if err := chainState.ValidateFeeQuoterConfig(chain); err != nil {
		return err
	}

	return ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.FeeQuoter: true})
}

func SetMaxFeeJuelsPerMsg(e cldf.Environment, cfg SetMaxFeeJuelsPerMsgConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]

	fqConfig, _, _ := solState.FindConfigPDA(chainState.FeeQuoter)
	fqUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.FeeQuoter,
		solana.PublicKey{},
		"")

	solFeeQuoter.SetProgramID(chainState.FeeQuoter)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.FeeQuoter,
		solana.PublicKey{},
		"")
	instruction, err := solFeeQuoter.NewSetMaxFeeJuelsPerMsgInstruction(
		cfg.MaxFeeJuelsPerMsg,
		fqConfig,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if fqUsingMCMS {
		tx, err := BuildMCMSTxn(instruction, chainState.FeeQuoter.String(), shared.FeeQuoter)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to SetMaxFeeJuelsPerMsg in Solana", cfg.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	if err := chain.Confirm([]solana.Instruction{instruction}); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}
