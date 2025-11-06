package solana

import (
	"context"
	"errors"
	"fmt"
	"math"

	solBinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/fee_quoter"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	ccipcommoncs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/helpers/pointer"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// use this changeset to add a token transfer fee for a remote chain to solana (used for very specific cases)
var _ cldf.ChangeSetV2[TokenTransferFeeForRemoteChainConfigV2] = AddTokenTransferFeeForRemoteChainV2

// use this changeset to add a billing token to solana
var _ cldf.ChangeSet[BillingTokenConfig] = AddBillingTokenChangeset

// use this changeset to add a token transfer fee for a remote chain to solana (used for very specific cases)
var _ cldf.ChangeSet[TokenTransferFeeForRemoteChainConfig] = AddTokenTransferFeeForRemoteChain

// use this changeset to withdraw billed funds on solana
var _ cldf.ChangeSet[WithdrawBilledFundsConfig] = WithdrawBilledFunds

// use this changeset to update prices for token and gas price updates on solana (emergency use only)
var _ cldf.ChangeSet[UpdatePricesConfig] = UpdatePrices

// use this changeset to set max fee juels per msg on solana (emergency use only)
var _ cldf.ChangeSet[SetMaxFeeJuelsPerMsgConfig] = SetMaxFeeJuelsPerMsg

// use this changeset to update price updaters on solana (emergency use only)
var _ cldf.ChangeSet[ModifyPriceUpdaterConfig] = ModifyPriceUpdater

// ADD BILLING TOKEN
type BillingTokenConfig struct {
	ChainSelector uint64
	Config        solFeeQuoter.BillingTokenConfig
	MCMS          *proposalutils.TimelockConfig

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
	mcms *proposalutils.TimelockConfig,
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

	tokenPubKey := cfg.Config.Mint
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
	ChainSelector      uint64
	RemoteChainConfigs map[uint64]solFeeQuoter.TokenTransferFeeConfig
	TokenPubKey        solana.PublicKey
	MCMS               *proposalutils.TimelockConfig
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
	for _, config := range cfg.RemoteChainConfigs {
		if config.DestBytesOverhead < 32 {
			e.Logger.Infow("dest bytes overhead is less than minimum. Setting to minimum value",
				"destBytesOverhead", config.DestBytesOverhead,
				"minDestBytesOverhead", MinDestBytesOverhead)
			config.DestBytesOverhead = MinDestBytesOverhead
		}
		if config.MinFeeUsdcents > config.MaxFeeUsdcents {
			return fmt.Errorf("min fee %d cannot be greater than max fee %d", config.MinFeeUsdcents, config.MaxFeeUsdcents)
		}
	}

	return ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.FeeQuoter: true})
}

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
	tokenPubKey := cfg.TokenPubKey
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
		"")
	solFeeQuoter.SetProgramID(chainState.FeeQuoter)
	txns := make([]mcmsTypes.Transaction, 0)
	for remoteChainSelector, config := range cfg.RemoteChainConfigs {
		remoteBillingPDA, _, _ := solState.FindFqPerChainPerTokenConfigPDA(remoteChainSelector, tokenPubKey, chainState.FeeQuoter)

		ix, err := solFeeQuoter.NewSetTokenTransferFeeConfigInstruction(
			remoteChainSelector,
			tokenPubKey,
			config,
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

		e.Logger.Infow("Token billing set for remote chain", "chainSelector ", cfg.ChainSelector, "remoteChainSelector ", remoteChainSelector, "tokenPubKey", tokenPubKey.String())

		if feeQuoterUsingMCMS {
			tx, err := BuildMCMSTxn(ix, chainState.FeeQuoter.String(), shared.FeeQuoter)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
			}
			txns = append(txns, *tx)
		}
	}

	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to set billing token for remote chain to Solana", cfg.MCMS.MinDelay, txns)
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
		raw.Append(solana.Meta(billingTokenConfigPDA).WRITE())
	}
	for _, update := range cfg.GasPriceUpdates {
		fqDestPDA, _, _ := solState.FindFqDestChainPDA(update.DestChainSelector, feeQuoterID)
		raw.Append(solana.Meta(fqDestPDA).WRITE())
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
	PriceUpdater       solana.PublicKey   // price updater to add or remove
	PriceUpdaterAction PriceUpdaterAction // add or remove price updater
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
	TransferAll   bool                          // transfer all or specific amount
	Amount        uint64                        // amount to transfer
	TokenPubKey   solana.PublicKey              // billing token to transfer
	MCMS          *proposalutils.TimelockConfig // timelock config for mcms
}

func (cfg WithdrawBilledFundsConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	tokenPubKey := cfg.TokenPubKey
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
	tokenPubKey := cfg.TokenPubKey
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

// ADD BILLING TOKEN FOR REMOTE CHAIN V2
//
// AddTokenTransferFeeForRemoteChainV2 wraps the original AddTokenTransferFeeForRemoteChain changeset with several devex
// improvements. In the V2 version you can provide a partial (or empty) config and any missing fields will be autofilled
// with sensible defaults. It's also possible to provide multiple tokens as input and all of them will be processed with
// a single MCMS proposal. A few other benefits include more detailed log messages and a new input schema structure that
// aligns more closely with the EVM equivalent of this changeset. This changeset also handles no-ops gracefully and uses
// the CLDF changeset V2 API.
var AddTokenTransferFeeForRemoteChainV2 = cldf.CreateChangeSet(
	addTokenTransferFeeForRemoteChainV2Logic,
	addTokenTransferFeeForRemoteChainV2Precondition,
)

type TokenTransferFeeForRemoteChainConfigV2 struct {
	// Map of source chain selector (solana family) => destination chain selector (any family) => config
	InputsByChain map[uint64]map[uint64]TokenTransferFeeForRemoteChainConfigArgsV2

	// Required MCMS config
	MCMS *proposalutils.TimelockConfig
}

type TokenTransferFeeForRemoteChainConfigArgsV2 struct {
	// Maps each token address to its token fee config
	TokenAddressToFeeConfig map[solana.PublicKey]OptionalFeeQuoterTokenTransferFeeConfig
}

type OptionalFeeQuoterTokenTransferFeeConfig struct {
	MinFeeUsdcents    *uint32
	MaxFeeUsdcents    *uint32
	DeciBps           *uint16
	DestGasOverhead   *uint32
	DestBytesOverhead *uint32
	IsEnabled         *bool
}

func (cfg TokenTransferFeeForRemoteChainConfigV2) buildOrchestrateChangesetsConfig(env cldf.Environment) (ccipcommoncs.OrchestrateChangesetsConfig, error) {
	env.Logger.Info("building orchestrate changesets config")
	if cfg.MCMS == nil {
		return ccipcommoncs.OrchestrateChangesetsConfig{}, errors.New("MCMS config is required")
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return ccipcommoncs.OrchestrateChangesetsConfig{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	changesets := []ccipcommoncs.WithConfig{}
	for srcSelector, dst := range cfg.InputsByChain {
		// get solana chain state and solana chain info
		err := stateview.ValidateChain(env, state, srcSelector, cfg.MCMS)
		if err != nil {
			return ccipcommoncs.OrchestrateChangesetsConfig{}, fmt.Errorf("failed to validate src chain (src = %d): %w", srcSelector, err)
		}
		solChainState, ok := state.SolChains[srcSelector]
		if !ok {
			return ccipcommoncs.OrchestrateChangesetsConfig{}, fmt.Errorf("failed to find solana chain with selector '%d' in state", srcSelector)
		}
		solChain, ok := env.BlockChains.SolanaChains()[srcSelector]
		if !ok {
			return ccipcommoncs.OrchestrateChangesetsConfig{}, fmt.Errorf("failed to find solana chain with selector '%d' in environment", srcSelector)
		}

		// prepare for iteration
		remoteChainConfigsByToken := map[solana.PublicKey]map[uint64]solFeeQuoter.TokenTransferFeeConfig{}
		env.Logger.Infof(
			"successfully found Solana fee quoter in state for selector %d: %s",
			srcSelector, solChainState.FeeQuoter.String(),
		)

		// 1st pass -> populate remote chain configs for each (token, dst) pair
		env.Logger.Infof("building remote chain configs for chain %d", srcSelector)
		for dstSelector, dstConfig := range dst {
			// perform basic validations on the first pass
			if srcSelector == dstSelector {
				return ccipcommoncs.OrchestrateChangesetsConfig{}, fmt.Errorf("destination chain cannot be the same as source chain (src = %d, dst = %d)", srcSelector, dstSelector)
			}
			if err := stateview.ValidateChain(env, state, dstSelector, cfg.MCMS); err != nil {
				return ccipcommoncs.OrchestrateChangesetsConfig{}, fmt.Errorf("failed to validate dst chain (src = %d, dst = %d): %w", srcSelector, dstSelector, err)
			}

			// build the config map
			for tokenAddress, feeConfig := range dstConfig.TokenAddressToFeeConfig {
				// get the remote billing PDA for the given (token, dst) pair
				env.Logger.Infof("processing inputs src = %d, dst = %d, token = %s", srcSelector, dstSelector, tokenAddress.String())
				remoteBillingPDA, _, err := solState.FindFqPerChainPerTokenConfigPDA(dstSelector, tokenAddress, solChainState.FeeQuoter)
				if err != nil {
					return ccipcommoncs.OrchestrateChangesetsConfig{}, fmt.Errorf("failed to find remote billing token config pda (src = %d, dst = %d, token = %s): %w", srcSelector, dstSelector, tokenAddress.String(), err)
				}

				// get the token transfer fee config from the fee quoter - if it doesn't exist, then the zero struct will be returned and `IsEnabled` will be `false`
				env.Logger.Infof("remote billing PDA = %s", remoteBillingPDA.String())
				var curConfig solFeeQuoter.PerChainPerTokenConfig
				err = solChain.GetAccountDataBorshInto(env.GetContext(), remoteBillingPDA, &curConfig)
				if !errors.Is(err, rpc.ErrNotFound) && err != nil {
					return ccipcommoncs.OrchestrateChangesetsConfig{}, fmt.Errorf("failed to deserialize PerChainPerTokenConfig (src = %d, dst = %d, token = %s, pda = %s): %w", srcSelector, dstSelector, tokenAddress.String(), remoteBillingPDA.String(), err)
				}

				// if the config has not been set yet, we auto-fill any missing input fields with sensible defaults
				env.Logger.Infof("current config = %+v", curConfig)
				if !curConfig.TokenTransferConfig.IsEnabled {
					// only use sensible defaults to fill in missing fields - do not overwrite anything that the user provided
					minFeeUsdCents := pointer.Coalesce(feeConfig.MinFeeUsdcents, uint32(175_000))
					maxFeeUsdCents := pointer.Coalesce(feeConfig.MaxFeeUsdcents, math.MaxUint32)
					destGasOverhead := pointer.Coalesce(feeConfig.DestGasOverhead, uint32(90_000))
					destBytesOverhead := pointer.Coalesce(feeConfig.DestBytesOverhead, uint32(32))
					deciBps := pointer.Coalesce(feeConfig.DeciBps, uint16(0))
					isEnabled := pointer.Coalesce(feeConfig.IsEnabled, true)
					env.Logger.Infof("config is not set - populating missing fields in user input with sensible defaults: %+v", feeConfig)

					// fill in the missing values in-place
					feeConfig.MinFeeUsdcents = &minFeeUsdCents
					feeConfig.MaxFeeUsdcents = &maxFeeUsdCents
					feeConfig.DeciBps = &deciBps
					feeConfig.DestGasOverhead = &destGasOverhead
					feeConfig.DestBytesOverhead = &destBytesOverhead
					feeConfig.IsEnabled = &isEnabled
					env.Logger.Infof("missing fields in user input have been auto-filled with sensible defaults: %+v", feeConfig)
				}

				// at this point, we're either using inputs from the user (highest precedence), fallback values from the chain, or pre-defined sensible defaults
				newConfig := solFeeQuoter.TokenTransferFeeConfig{
					MinFeeUsdcents:    pointer.Coalesce(feeConfig.MinFeeUsdcents, curConfig.TokenTransferConfig.MinFeeUsdcents),
					MaxFeeUsdcents:    pointer.Coalesce(feeConfig.MaxFeeUsdcents, curConfig.TokenTransferConfig.MaxFeeUsdcents),
					DeciBps:           pointer.Coalesce(feeConfig.DeciBps, curConfig.TokenTransferConfig.DeciBps),
					DestGasOverhead:   pointer.Coalesce(feeConfig.DestGasOverhead, curConfig.TokenTransferConfig.DestGasOverhead),
					DestBytesOverhead: pointer.Coalesce(feeConfig.DestBytesOverhead, curConfig.TokenTransferConfig.DestBytesOverhead),
					IsEnabled:         pointer.Coalesce(feeConfig.IsEnabled, curConfig.TokenTransferConfig.IsEnabled),
				}

				// make sure that the config is still valid after merge
				if newConfig.MinFeeUsdcents >= newConfig.MaxFeeUsdcents {
					return ccipcommoncs.OrchestrateChangesetsConfig{}, fmt.Errorf("min fee (%d) must be less than max fee (%d)", newConfig.MinFeeUsdcents, newConfig.MaxFeeUsdcents)
				}

				// check if the new config is different from the on-chain config
				isDifferent := newConfig.MinFeeUsdcents != curConfig.TokenTransferConfig.MinFeeUsdcents ||
					newConfig.MaxFeeUsdcents != curConfig.TokenTransferConfig.MaxFeeUsdcents ||
					newConfig.DeciBps != curConfig.TokenTransferConfig.DeciBps ||
					newConfig.DestGasOverhead != curConfig.TokenTransferConfig.DestGasOverhead ||
					newConfig.DestBytesOverhead != curConfig.TokenTransferConfig.DestBytesOverhead ||
					newConfig.IsEnabled != curConfig.TokenTransferConfig.IsEnabled

				// only perform an update if the new config is different from the on-chain config
				env.Logger.Infof("constructed new token transfer fee config: %+v", newConfig)
				if !isDifferent {
					env.Logger.Infof(
						"skipping update since input config is the same as on-chain config (src=%d, dst=%d, token=%s): %+v",
						srcSelector, dstSelector, tokenAddress.String(), curConfig,
					)
					continue
				}

				// update the config map
				if _, ok := remoteChainConfigsByToken[tokenAddress]; !ok {
					remoteChainConfigsByToken[tokenAddress] = map[uint64]solFeeQuoter.TokenTransferFeeConfig{}
				}
				remoteChainConfigsByToken[tokenAddress][dstSelector] = newConfig
			}
		}

		// 2nd pass -> populate token transfer fee configs and changesets
		env.Logger.Infof("detected %d token(s) to configure", len(remoteChainConfigsByToken))
		for tokenPubKey, remoteChainConfigs := range remoteChainConfigsByToken {
			if len(remoteChainConfigs) == 0 {
				env.Logger.Infof("detected no changes for token %s - skipping", tokenPubKey.String())
				continue
			}

			tokenTransferFeeConfig := TokenTransferFeeForRemoteChainConfig{
				RemoteChainConfigs: remoteChainConfigs,
				ChainSelector:      srcSelector,
				TokenPubKey:        tokenPubKey,
				MCMS:               cfg.MCMS,
			}

			err := tokenTransferFeeConfig.Validate(env, state)
			if err != nil {
				return ccipcommoncs.OrchestrateChangesetsConfig{}, fmt.Errorf(
					"validation failed for token transfer fee config (src=%d, token=%s, input=%+v): %w",
					srcSelector, tokenPubKey.String(), tokenTransferFeeConfig, err,
				)
			}

			changesets = append(
				changesets,
				ccipcommoncs.CreateGenericChangeSetWithConfig(
					cldf.CreateLegacyChangeSet(AddTokenTransferFeeForRemoteChain),
					tokenTransferFeeConfig,
				),
			)
		}
	}

	return ccipcommoncs.OrchestrateChangesetsConfig{
		Description: "Apply token transfer fee configs from Solana -> *",
		ChangeSets:  changesets,
		MCMS:        cfg.MCMS,
	}, nil
}

func addTokenTransferFeeForRemoteChainV2Precondition(env cldf.Environment, cfg TokenTransferFeeForRemoteChainConfigV2) error {
	if len(cfg.InputsByChain) == 0 {
		env.Logger.Info("no inputs were provided - exiting precondition stage gracefully")
		return nil
	}

	input, err := cfg.buildOrchestrateChangesetsConfig(env)
	if err != nil {
		return fmt.Errorf("failed to build OrchestrateChangesetsConfig: %w", err)
	}

	if len(input.ChangeSets) == 0 {
		env.Logger.Info("no changesets to orchestrate - exiting precondition stage gracefully")
		return nil
	}

	return ccipcommoncs.OrchestrateChangesets.VerifyPreconditions(env, input)
}

func addTokenTransferFeeForRemoteChainV2Logic(env cldf.Environment, cfg TokenTransferFeeForRemoteChainConfigV2) (cldf.ChangesetOutput, error) {
	if len(cfg.InputsByChain) == 0 {
		env.Logger.Info("no inputs were provided - exiting apply stage gracefully")
		return cldf.ChangesetOutput{}, nil
	}

	input, err := cfg.buildOrchestrateChangesetsConfig(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build OrchestrateChangesetsConfig: %w", err)
	}

	if len(input.ChangeSets) == 0 {
		env.Logger.Info("no changesets to orchestrate - exiting apply stage gracefully")
		return cldf.ChangesetOutput{}, nil
	}

	return ccipcommoncs.OrchestrateChangesets.Apply(env, input)
}
