package solana

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/gagliardetto/solana-go"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// use these three changesets to add a remote chain to solana
var _ cldf.ChangeSet[AddRemoteChainToRouterConfig] = AddRemoteChainToRouter
var _ cldf.ChangeSet[AddRemoteChainToOffRampConfig] = AddRemoteChainToOffRamp
var _ cldf.ChangeSet[AddRemoteChainToFeeQuoterConfig] = AddRemoteChainToFeeQuoter

type AddRemoteChainToRouterConfig struct {
	ChainSelector uint64
	// UpdatesByChain is a mapping of SVM chain selector -> remote chain selector -> remote chain config update
	UpdatesByChain map[uint64]*RouterConfig
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be achieved by calling this function multiple times)
	MCMS *proposalutils.TimelockConfig
}

type RouterConfig struct {
	// if enabling AllowedSender -> it needs to be a complete list
	// onchain just clones what we pass in
	// and tooling does not handle upserts
	// so you have to clone what is in state, edit the list, and then pass into this changeset
	RouterDestinationConfig solRouter.DestChainConfig
	// We have different instructions for add vs update, so we need to know which one to use
	IsUpdate bool
}

func (cfg *AddRemoteChainToRouterConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	chainState := state.SolChains[cfg.ChainSelector]
	chain, ok := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if !ok {
		return fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}
	if err := chainState.ValidateRouterConfig(chain); err != nil {
		return err
	}

	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.Router: true}); err != nil {
		return err
	}
	routerProgramAddress, routerConfigPDA, _ := chainState.GetRouterInfo()
	var routerConfigAccount solRouter.Config
	// already validated that router config exists
	_ = chain.GetAccountDataBorshInto(context.Background(), routerConfigPDA, &routerConfigAccount)

	supportedChains := state.SupportedChains()
	for remote, remoteConfig := range cfg.UpdatesByChain {
		if _, ok := supportedChains[remote]; !ok {
			return fmt.Errorf("remote chain %d is not supported", remote)
		}
		if remote == routerConfigAccount.SvmChainSelector {
			return fmt.Errorf("cannot add remote chain %d with same chain selector as current chain %d", remote, cfg.ChainSelector)
		}
		if err := state.ValidateRamp(remote, shared.OnRamp); err != nil {
			return err
		}
		routerDestChainPDA, err := solState.FindDestChainStatePDA(remote, routerProgramAddress)
		if err != nil {
			return fmt.Errorf("failed to find dest chain state pda for remote chain %d: %w", remote, err)
		}
		var destChainStateAccount solRouter.DestChain
		err = chain.GetAccountDataBorshInto(context.Background(), routerDestChainPDA, &destChainStateAccount)
		if err == nil {
			e.Logger.Infow("remote chain already configured. setting as update", "remoteChainSel", remote)
			remoteConfig.IsUpdate = true
		}
	}
	return nil
}

// Adds new remote chain configurations
func AddRemoteChainToRouter(e cldf.Environment, cfg AddRemoteChainToRouterConfig) (cldf.ChangesetOutput, error) {
	s, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	if err := cfg.Validate(e, s); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	ab := cldf.NewMemoryAddressBook()
	txns, err := doAddRemoteChainToRouter(e, s, cfg, ab)
	if err != nil {
		return cldf.ChangesetOutput{AddressBook: ab}, err
	}

	// create proposals for ixns
	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to add remote chains to Solana", cfg.MCMS.MinDelay, txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
			AddressBook:           ab,
		}, nil
	}
	return cldf.ChangesetOutput{AddressBook: ab}, nil
}

func doAddRemoteChainToRouter(
	e cldf.Environment,
	s stateview.CCIPOnChainState,
	cfg AddRemoteChainToRouterConfig,
	ab cldf.AddressBook) ([]mcmsTypes.Transaction, error) {
	txns := make([]mcmsTypes.Transaction, 0)
	chainSel := cfg.ChainSelector
	updates := cfg.UpdatesByChain
	chain := e.BlockChains.SolanaChains()[chainSel]
	chainState := s.SolChains[chainSel]
	ccipRouterID, routerConfigPDA, _ := s.SolChains[chainSel].GetRouterInfo()
	offRampID := s.SolChains[chainSel].OffRamp
	routerUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.Router,
		solana.PublicKey{},
		"",
	)
	lookUpTableEntries := make([]solana.PublicKey, 0)
	// router setup
	solRouter.SetProgramID(ccipRouterID)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.Router,
		solana.PublicKey{},
		"",
	)

	for remoteChainSel, update := range updates {
		// verified while loading state
		routerRemoteStatePDA, _ := solState.FindDestChainStatePDA(remoteChainSel, ccipRouterID)
		allowedOffRampRemotePDA, _ := solState.FindAllowedOfframpPDA(remoteChainSel, offRampID, ccipRouterID)

		if update.IsUpdate {
			routerIx, err := solRouter.NewUpdateDestChainConfigInstruction(
				remoteChainSel,
				// TODO: this needs to be merged with what the user is sending in and whats their onchain.
				// right now, the user will have to send the final version of the config.
				update.RouterDestinationConfig,
				routerRemoteStatePDA,
				routerConfigPDA,
				authority,
				solana.SystemProgramID,
			).ValidateAndBuild()
			if err != nil {
				return txns, fmt.Errorf("failed to generate update router config instructions: %w", err)
			}
			e.Logger.Infow("update router config for remote chain", "remoteChainSel", remoteChainSel)
			if routerUsingMCMS {
				tx, err := BuildMCMSTxn(routerIx, ccipRouterID.String(), shared.Router)
				if err != nil {
					return txns, fmt.Errorf("failed to create update router config transaction: %w", err)
				}
				txns = append(txns, *tx)
			} else {
				err = chain.Confirm([]solana.Instruction{routerIx})
				if err != nil {
					return txns, fmt.Errorf("failed to confirm update router config instructions: %w", err)
				}
			}
		} else {
			// new remote chain
			lookUpTableEntries = append(lookUpTableEntries,
				routerRemoteStatePDA,
			)
			// generate instructions
			routerIx, err := solRouter.NewAddChainSelectorInstruction(
				remoteChainSel,
				update.RouterDestinationConfig,
				routerRemoteStatePDA,
				routerConfigPDA,
				authority,
				solana.SystemProgramID,
			).ValidateAndBuild()
			if err != nil {
				return txns, fmt.Errorf("failed to generate add router config instructions: %w", err)
			}
			e.Logger.Infow("add router config for remote chain", "remoteChainSel", remoteChainSel)
			routerOfframpIx, err := solRouter.NewAddOfframpInstruction(
				remoteChainSel,
				offRampID,
				allowedOffRampRemotePDA,
				routerConfigPDA,
				authority,
				solana.SystemProgramID,
			).ValidateAndBuild()
			if err != nil {
				return txns, fmt.Errorf("failed to generate instructions: %w", err)
			}
			e.Logger.Infow("add offramp to router for remote chain", "remoteChainSel", remoteChainSel)
			if routerUsingMCMS {
				// build transactions if mcms
				for _, ix := range []solana.Instruction{routerIx, routerOfframpIx} {
					tx, err := BuildMCMSTxn(ix, ccipRouterID.String(), shared.Router)
					if err != nil {
						return txns, fmt.Errorf("failed to create add router config transaction: %w", err)
					}
					txns = append(txns, *tx)
				}
			} else {
				// confirm ixns if not mcms
				err = chain.Confirm([]solana.Instruction{routerIx, routerOfframpIx})
				if err != nil {
					return txns, fmt.Errorf("failed to confirm add router config instructions: %w", err)
				}
			}
			// add to address book
			tv := cldf.NewTypeAndVersion(shared.RemoteDest, deployment.Version1_0_0)
			remoteChainSelStr := strconv.FormatUint(remoteChainSel, 10)
			tv.AddLabel(remoteChainSelStr)
			err = ab.Save(chainSel, routerRemoteStatePDA.String(), tv)
			if err != nil {
				return txns, fmt.Errorf("failed to save dest chain state to address book: %w", err)
			}
		}
	}

	if len(lookUpTableEntries) > 0 {
		err := extendLookupTable(e, chain, offRampID, lookUpTableEntries)
		if err != nil {
			return txns, fmt.Errorf("failed to extend lookup table: %w", err)
		}
	}

	return txns, nil
}

type AddRemoteChainToFeeQuoterConfig struct {
	ChainSelector uint64
	// UpdatesByChain is a mapping of SVM chain selector -> remote chain selector -> remote chain config update
	UpdatesByChain map[uint64]*FeeQuoterConfig
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be achieved by calling this function multiple times)
	MCMS *proposalutils.TimelockConfig
}

type FeeQuoterConfig struct {
	FeeQuoterDestinationConfig solFeeQuoter.DestChainConfig
	// We have different instructions for add vs update, so we need to know which one to use
	IsUpdate bool
}

func (cfg *AddRemoteChainToFeeQuoterConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	chainState := state.SolChains[cfg.ChainSelector]
	chain, ok := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if !ok {
		return fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}

	if err := chainState.ValidateFeeQuoterConfig(chain); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.FeeQuoter: true}); err != nil {
		return err
	}
	supportedChains := state.SupportedChains()
	for remote, remoteConfig := range cfg.UpdatesByChain {
		if _, ok := supportedChains[remote]; !ok {
			return fmt.Errorf("remote chain %d is not supported", remote)
		}
		if err := state.ValidateRamp(remote, shared.OnRamp); err != nil {
			return err
		}
		fqRemoteChainPDA, _, err := solState.FindFqDestChainPDA(remote, chainState.FeeQuoter)
		if err != nil {
			return fmt.Errorf("failed to find dest chain state pda for remote chain %d: %w", remote, err)
		}
		var destChainStateAccount solFeeQuoter.DestChain
		err = chain.GetAccountDataBorshInto(context.Background(), fqRemoteChainPDA, &destChainStateAccount)
		if err == nil {
			e.Logger.Infow("remote chain already configured. setting as update", "remoteChainSel", remote)
			remoteConfig.IsUpdate = true
		}
	}
	return nil
}

// Adds new remote chain configurations
func AddRemoteChainToFeeQuoter(e cldf.Environment, cfg AddRemoteChainToFeeQuoterConfig) (cldf.ChangesetOutput, error) {
	s, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	if err := cfg.Validate(e, s); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	ab := cldf.NewMemoryAddressBook()
	txns, err := doAddRemoteChainToFeeQuoter(e, s, cfg, ab)
	if err != nil {
		return cldf.ChangesetOutput{AddressBook: ab}, err
	}

	// create proposals for ixns
	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to add remote chains to Solana", cfg.MCMS.MinDelay, txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
			AddressBook:           ab,
		}, nil
	}
	return cldf.ChangesetOutput{AddressBook: ab}, nil
}

func doAddRemoteChainToFeeQuoter(
	e cldf.Environment,
	s stateview.CCIPOnChainState,
	cfg AddRemoteChainToFeeQuoterConfig,
	ab cldf.AddressBook) ([]mcmsTypes.Transaction, error) {
	txns := make([]mcmsTypes.Transaction, 0)
	chainSel := cfg.ChainSelector
	updates := cfg.UpdatesByChain
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
	lookUpTableEntries := make([]solana.PublicKey, 0)
	// fee quoter setup
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
	for remoteChainSel, update := range updates {
		// verified while loading state
		fqRemoteChainPDA, _, _ := solState.FindFqDestChainPDA(remoteChainSel, feeQuoterID)
		var feeQuoterIx solana.Instruction
		var err error
		if update.IsUpdate {
			feeQuoterIx, err = solFeeQuoter.NewUpdateDestChainConfigInstruction(
				remoteChainSel,
				// TODO: this needs to be merged with what the user is sending in and whats their onchain.
				// right now, the user will have to send the final version of the config.
				update.FeeQuoterDestinationConfig,
				s.SolChains[chainSel].FeeQuoterConfigPDA,
				fqRemoteChainPDA,
				authority,
			).ValidateAndBuild()
			e.Logger.Infow("update fee quoter config for remote chain", "remoteChainSel", remoteChainSel)
		} else {
			lookUpTableEntries = append(lookUpTableEntries,
				fqRemoteChainPDA,
			)
			feeQuoterIx, err = solFeeQuoter.NewAddDestChainInstruction(
				remoteChainSel,
				update.FeeQuoterDestinationConfig,
				s.SolChains[chainSel].FeeQuoterConfigPDA,
				fqRemoteChainPDA,
				authority,
				solana.SystemProgramID,
			).ValidateAndBuild()
			e.Logger.Infow("add fee quoter config for remote chain", "remoteChainSel", remoteChainSel)
		}
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
			err = chain.Confirm([]solana.Instruction{feeQuoterIx})
			if err != nil {
				return txns, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}
	}

	if len(lookUpTableEntries) > 0 {
		err := extendLookupTable(e, chain, offRampID, lookUpTableEntries)
		if err != nil {
			return txns, fmt.Errorf("failed to extend lookup table: %w", err)
		}
	}

	return txns, nil
}

type AddRemoteChainToOffRampConfig struct {
	ChainSelector uint64
	// UpdatesByChain is a mapping of SVM chain selector -> remote chain selector -> remote chain config update
	UpdatesByChain map[uint64]*OffRampConfig
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be achieved by calling this function multiple times)
	MCMS *proposalutils.TimelockConfig
}

type OffRampConfig struct {
	// source
	EnabledAsSource bool
	// We have different instructions for add vs update, so we need to know which one to use
	IsUpdate bool
}

func (cfg *AddRemoteChainToOffRampConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	chainState := state.SolChains[cfg.ChainSelector]
	chain, ok := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if !ok {
		return fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}

	if err := chainState.ValidateOffRampConfig(chain); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.OffRamp: true}); err != nil {
		return err
	}

	supportedChains := state.SupportedChains()
	for remote, remoteConfig := range cfg.UpdatesByChain {
		if _, ok := supportedChains[remote]; !ok {
			return fmt.Errorf("remote chain %d is not supported", remote)
		}
		if err := state.ValidateRamp(remote, shared.OnRamp); err != nil {
			return err
		}
		offRampRemoteStatePDA, _, err := solState.FindOfframpSourceChainPDA(remote, chainState.OffRamp)
		if err != nil {
			return fmt.Errorf("failed to find dest chain state pda for remote chain %d: %w", remote, err)
		}
		var destChainStateAccount solOffRamp.SourceChain
		err = chain.GetAccountDataBorshInto(context.Background(), offRampRemoteStatePDA, &destChainStateAccount)
		if err == nil {
			e.Logger.Infow("remote chain already configured. setting as update", "remoteChainSel", remote)
			remoteConfig.IsUpdate = true
		}
	}
	return nil
}

// Adds new remote chain configurations
func AddRemoteChainToOffRamp(e cldf.Environment, cfg AddRemoteChainToOffRampConfig) (cldf.ChangesetOutput, error) {
	s, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	if err := cfg.Validate(e, s); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	ab := cldf.NewMemoryAddressBook()
	txns, err := doAddRemoteChainToOffRamp(e, s, cfg, ab)
	if err != nil {
		return cldf.ChangesetOutput{AddressBook: ab}, err
	}

	// create proposals for ixns
	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to add remote chains to Solana", cfg.MCMS.MinDelay, txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
			AddressBook:           ab,
		}, nil
	}
	return cldf.ChangesetOutput{AddressBook: ab}, nil
}

func doAddRemoteChainToOffRamp(
	e cldf.Environment,
	s stateview.CCIPOnChainState,
	cfg AddRemoteChainToOffRampConfig,
	ab cldf.AddressBook) ([]mcmsTypes.Transaction, error) {
	txns := make([]mcmsTypes.Transaction, 0)
	chainSel := cfg.ChainSelector
	updates := cfg.UpdatesByChain
	chain := e.BlockChains.SolanaChains()[chainSel]
	chainState := s.SolChains[chainSel]
	offRampID := s.SolChains[chainSel].OffRamp
	offRampUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.OffRamp,
		solana.PublicKey{},
		"")
	lookUpTableEntries := make([]solana.PublicKey, 0)
	solOffRamp.SetProgramID(offRampID)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.OffRamp,
		solana.PublicKey{},
		"",
	)

	for remoteChainSel, update := range updates {
		// verified while loading state
		offRampRemoteStatePDA, _, _ := solState.FindOfframpSourceChainPDA(remoteChainSel, offRampID)
		// offramp setup
		validSourceChainConfig, err := getSourceChainConfig(s, remoteChainSel, update.EnabledAsSource)
		if err != nil {
			return txns, fmt.Errorf("failed to get source chain config: %w", err)
		}

		var offRampIx solana.Instruction
		if update.IsUpdate {
			offRampIx, err = solOffRamp.NewUpdateSourceChainConfigInstruction(
				remoteChainSel,
				validSourceChainConfig,
				offRampRemoteStatePDA,
				s.SolChains[chainSel].OffRampConfigPDA,
				authority,
			).ValidateAndBuild()
			if err != nil {
				return txns, fmt.Errorf("failed to generate instructions: %w", err)
			}
			e.Logger.Infow("update offramp config for remote chain", "remoteChainSel", remoteChainSel)
		} else {
			lookUpTableEntries = append(lookUpTableEntries,
				offRampRemoteStatePDA,
			)
			offRampIx, err = solOffRamp.NewAddSourceChainInstruction(
				remoteChainSel,
				validSourceChainConfig,
				offRampRemoteStatePDA,
				s.SolChains[chainSel].OffRampConfigPDA,
				authority,
				solana.SystemProgramID,
			).ValidateAndBuild()
			if err != nil {
				return txns, fmt.Errorf("failed to generate instructions: %w", err)
			}
			e.Logger.Infow("add offramp config for remote chain", "remoteChainSel", remoteChainSel)
			remoteChainSelStr := strconv.FormatUint(remoteChainSel, 10)
			tv := cldf.NewTypeAndVersion(shared.RemoteSource, deployment.Version1_0_0)
			tv.AddLabel(remoteChainSelStr)
			err = ab.Save(chainSel, offRampRemoteStatePDA.String(), tv)
			if err != nil {
				return txns, fmt.Errorf("failed to save source chain state to address book: %w", err)
			}
		}

		if offRampUsingMCMS {
			tx, err := BuildMCMSTxn(offRampIx, offRampID.String(), shared.OffRamp)
			if err != nil {
				return txns, fmt.Errorf("failed to create transaction: %w", err)
			}
			txns = append(txns, *tx)
		} else {
			err = chain.Confirm([]solana.Instruction{offRampIx})
			if err != nil {
				return txns, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}
	}

	if len(lookUpTableEntries) > 0 {
		err := extendLookupTable(e, chain, offRampID, lookUpTableEntries)
		if err != nil {
			return txns, fmt.Errorf("failed to extend lookup table: %w", err)
		}
	}

	return txns, nil
}

func getSourceChainConfig(s stateview.CCIPOnChainState, remoteChainSel uint64, enabledAsSource bool) (solOffRamp.SourceChainConfig, error) {
	var onRampAddress solOffRamp.OnRampAddress
	// already verified, skipping errcheck
	addressBytes, _ := s.GetOnRampAddressBytes(remoteChainSel)
	copy(onRampAddress.Bytes[:], addressBytes)
	addressBytesLen := len(addressBytes)
	if addressBytesLen < 0 || addressBytesLen > math.MaxUint32 {
		return solOffRamp.SourceChainConfig{}, fmt.Errorf("address bytes length %d is outside valid uint32 range", addressBytesLen)
	}
	onRampAddress.Len = uint32(addressBytesLen)
	validSourceChainConfig := solOffRamp.SourceChainConfig{
		OnRamp:    onRampAddress,
		IsEnabled: enabledAsSource,
	}
	return validSourceChainConfig, nil
}

func extendLookupTable(e cldf.Environment, chain cldf_solana.Chain, offRampID solana.PublicKey, lookUpTableEntries []solana.PublicKey) error {
	addressLookupTable, err := solanastateview.FetchOfframpLookupTable(e.GetContext(), chain, offRampID)
	if err != nil {
		return fmt.Errorf("failed to get offramp reference addresses: %w", err)
	}

	addresses, err := solCommonUtil.GetAddressLookupTable(
		e.GetContext(),
		chain.Client,
		addressLookupTable)
	if err != nil {
		return fmt.Errorf("failed to get address lookup table: %w", err)
	}

	// calculate diff and add new entries
	seen := make(map[solana.PublicKey]bool)
	toAdd := make([]solana.PublicKey, 0)
	for _, entry := range addresses {
		seen[entry] = true
	}
	for _, entry := range lookUpTableEntries {
		if _, ok := seen[entry]; !ok {
			toAdd = append(toAdd, entry)
		}
	}
	if len(toAdd) == 0 {
		e.Logger.Infow("no new entries to add to lookup table")
		return nil
	}

	e.Logger.Debugw("Populating lookup table", "keys", toAdd)
	if err := solCommonUtil.ExtendLookupTable(
		e.GetContext(),
		chain.Client,
		addressLookupTable,
		*chain.DeployerKey,
		toAdd,
	); err != nil {
		return fmt.Errorf("failed to extend lookup table: %w", err)
	}
	return nil
}
