package solana

import (
	"context"
	"math"

	"fmt"
	"strconv"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

// use these three changesets to add a remote chain to solana
var _ deployment.ChangeSet[AddRemoteChainToRouterConfig] = AddRemoteChainToRouter
var _ deployment.ChangeSet[AddRemoteChainToOffRampConfig] = AddRemoteChainToOffRamp
var _ deployment.ChangeSet[AddRemoteChainToFeeQuoterConfig] = AddRemoteChainToFeeQuoter

type AddRemoteChainToRouterConfig struct {
	ChainSelector uint64
	// UpdatesByChain is a mapping of SVM chain selector -> remote chain selector -> remote chain config update
	UpdatesByChain map[uint64]RouterConfig
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be achieved by calling this function multiple times)
	MCMSSolana *MCMSConfigSolana
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

func (cfg AddRemoteChainToRouterConfig) Validate(e deployment.Environment) error {
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.SolChains[cfg.ChainSelector]
	chain, ok := e.SolChains[cfg.ChainSelector]
	if !ok {
		return fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}
	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}

	if err := ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, solana.PublicKey{}); err != nil {
		return err
	}
	routerProgramAddress, routerConfigPDA, _ := chainState.GetRouterInfo()
	var routerConfigAccount solRouter.Config
	// already validated that router config exists
	_ = chain.GetAccountDataBorshInto(context.Background(), routerConfigPDA, &routerConfigAccount)

	supportedChains := state.SupportedChains()
	for remote := range cfg.UpdatesByChain {
		if _, ok := supportedChains[remote]; !ok {
			return fmt.Errorf("remote chain %d is not supported", remote)
		}
		if remote == routerConfigAccount.SvmChainSelector {
			return fmt.Errorf("cannot add remote chain %d with same chain selector as current chain %d", remote, cfg.ChainSelector)
		}
		if err := state.ValidateRamp(remote, ccipChangeset.OnRamp); err != nil {
			return err
		}
		routerDestChainPDA, err := solState.FindDestChainStatePDA(remote, routerProgramAddress)
		if err != nil {
			return fmt.Errorf("failed to find dest chain state pda for remote chain %d: %w", remote, err)
		}
		if !cfg.UpdatesByChain[remote].IsUpdate {
			var destChainStateAccount solRouter.DestChain
			err = chain.GetAccountDataBorshInto(context.Background(), routerDestChainPDA, &destChainStateAccount)
			if err == nil {
				return fmt.Errorf("remote %d is already configured on solana chain router %d", remote, cfg.ChainSelector)
			}
		}
	}
	return nil
}

// Adds new remote chain configurations
func AddRemoteChainToRouter(e deployment.Environment, cfg AddRemoteChainToRouterConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	s, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	ab := deployment.NewMemoryAddressBook()
	txns, err := doAddRemoteChainToRouter(e, s, cfg, ab)
	if err != nil {
		return deployment.ChangesetOutput{AddressBook: ab}, err
	}

	// create proposals for ixns
	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to add remote chains to Solana", cfg.MCMSSolana.MCMS.MinDelay, txns)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
			AddressBook:           ab,
		}, nil
	}
	return deployment.ChangesetOutput{AddressBook: ab}, nil
}

func doAddRemoteChainToRouter(
	e deployment.Environment,
	s ccipChangeset.CCIPOnChainState,
	cfg AddRemoteChainToRouterConfig,
	ab deployment.AddressBook) ([]mcmsTypes.Transaction, error) {
	txns := make([]mcmsTypes.Transaction, 0)
	chainSel := cfg.ChainSelector
	updates := cfg.UpdatesByChain
	chain := e.SolChains[chainSel]
	ccipRouterID, routerConfigPDA, _ := s.SolChains[chainSel].GetRouterInfo()
	offRampID := s.SolChains[chainSel].OffRamp
	routerUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock
	lookUpTableEntries := make([]solana.PublicKey, 0)
	// router setup
	solRouter.SetProgramID(ccipRouterID)
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.Router,
		solana.PublicKey{})
	if err != nil {
		return txns, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
	for remoteChainSel, update := range updates {
		ixns := make([]solana.Instruction, 0)
		// verified while loading state
		routerRemoteStatePDA, _ := solState.FindDestChainStatePDA(remoteChainSel, ccipRouterID)
		allowedOffRampRemotePDA, _ := solState.FindAllowedOfframpPDA(remoteChainSel, offRampID, ccipRouterID)

		if !update.IsUpdate {
			lookUpTableEntries = append(lookUpTableEntries,
				routerRemoteStatePDA,
			)
		}

		var routerIx solana.Instruction
		if update.IsUpdate {
			routerIx, err = solRouter.NewUpdateDestChainConfigInstruction(
				remoteChainSel,
				// TODO: this needs to be merged with what the user is sending in and whats their onchain.
				// right now, the user will have to send the final version of the config.
				update.RouterDestinationConfig,
				routerRemoteStatePDA,
				routerConfigPDA,
				authority,
				solana.SystemProgramID,
			).ValidateAndBuild()
			e.Logger.Infow("update router config for remote chain", "remoteChainSel", remoteChainSel)
		} else {
			routerIx, err = solRouter.NewAddChainSelectorInstruction(
				remoteChainSel,
				update.RouterDestinationConfig,
				routerRemoteStatePDA,
				routerConfigPDA,
				authority,
				solana.SystemProgramID,
			).ValidateAndBuild()
			e.Logger.Infow("add router config for remote chain", "remoteChainSel", remoteChainSel)
		}
		if err != nil {
			return txns, fmt.Errorf("failed to generate instructions: %w", err)
		}
		if routerUsingMCMS {
			tx, err := BuildMCMSTxn(routerIx, ccipRouterID.String(), ccipChangeset.Router)
			if err != nil {
				return txns, fmt.Errorf("failed to create transaction: %w", err)
			}
			txns = append(txns, *tx)
		} else {
			ixns = append(ixns, routerIx)
		}

		if !update.IsUpdate {
			routerOfframpIx, err := solRouter.NewAddOfframpInstruction(
				remoteChainSel,
				offRampID,
				allowedOffRampRemotePDA,
				routerConfigPDA,
				authority,
				solana.SystemProgramID,
			).ValidateAndBuild()
			e.Logger.Infow("add offramp to router for remote chain", "remoteChainSel", remoteChainSel)
			if err != nil {
				return txns, fmt.Errorf("failed to generate instructions: %w", err)
			}
			if routerUsingMCMS {
				tx, err := BuildMCMSTxn(routerOfframpIx, ccipRouterID.String(), ccipChangeset.Router)
				if err != nil {
					return txns, fmt.Errorf("failed to create transaction: %w", err)
				}
				txns = append(txns, *tx)
			} else {
				ixns = append(ixns, routerOfframpIx)
			}
		}

		// confirm ixns if any
		if len(ixns) > 0 {
			err = chain.Confirm(ixns)
			if err != nil {
				return txns, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}

		if !update.IsUpdate {
			tv := deployment.NewTypeAndVersion(ccipChangeset.RemoteDest, deployment.Version1_0_0)
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
	UpdatesByChain map[uint64]FeeQuoterConfig
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be achieved by calling this function multiple times)
	MCMSSolana *MCMSConfigSolana
}

type FeeQuoterConfig struct {
	FeeQuoterDestinationConfig solFeeQuoter.DestChainConfig
	// We have different instructions for add vs update, so we need to know which one to use
	IsUpdate bool
}

func (cfg AddRemoteChainToFeeQuoterConfig) Validate(e deployment.Environment) error {
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.SolChains[cfg.ChainSelector]
	chain, ok := e.SolChains[cfg.ChainSelector]
	if !ok {
		return fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}

	if err := validateFeeQuoterConfig(chain, chainState); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, solana.PublicKey{}); err != nil {
		return err
	}
	supportedChains := state.SupportedChains()
	for remote := range cfg.UpdatesByChain {
		if _, ok := supportedChains[remote]; !ok {
			return fmt.Errorf("remote chain %d is not supported", remote)
		}
		if err := state.ValidateRamp(remote, ccipChangeset.OnRamp); err != nil {
			return err
		}
		fqRemoteChainPDA, _, err := solState.FindFqDestChainPDA(remote, chainState.FeeQuoter)
		if err != nil {
			return fmt.Errorf("failed to find dest chain state pda for remote chain %d: %w", remote, err)
		}
		if !cfg.UpdatesByChain[remote].IsUpdate {
			var destChainStateAccount solFeeQuoter.DestChain
			err = chain.GetAccountDataBorshInto(context.Background(), fqRemoteChainPDA, &destChainStateAccount)
			if err == nil {
				return fmt.Errorf("remote %d is already configured on solana chain feequoter %d", remote, cfg.ChainSelector)
			}
		}
	}
	return nil
}

// Adds new remote chain configurations
func AddRemoteChainToFeeQuoter(e deployment.Environment, cfg AddRemoteChainToFeeQuoterConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	s, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	ab := deployment.NewMemoryAddressBook()
	txns, err := doAddRemoteChainToFeeQuoter(e, s, cfg, ab)
	if err != nil {
		return deployment.ChangesetOutput{AddressBook: ab}, err
	}

	// create proposals for ixns
	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to add remote chains to Solana", cfg.MCMSSolana.MCMS.MinDelay, txns)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
			AddressBook:           ab,
		}, nil
	}
	return deployment.ChangesetOutput{AddressBook: ab}, nil
}

func doAddRemoteChainToFeeQuoter(
	e deployment.Environment,
	s ccipChangeset.CCIPOnChainState,
	cfg AddRemoteChainToFeeQuoterConfig,
	ab deployment.AddressBook) ([]mcmsTypes.Transaction, error) {
	txns := make([]mcmsTypes.Transaction, 0)
	chainSel := cfg.ChainSelector
	updates := cfg.UpdatesByChain
	chain := e.SolChains[chainSel]
	feeQuoterID := s.SolChains[chainSel].FeeQuoter
	offRampID := s.SolChains[chainSel].OffRamp
	feeQuoterUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.FeeQuoterOwnedByTimelock
	lookUpTableEntries := make([]solana.PublicKey, 0)
	// fee quoter setup
	solFeeQuoter.SetProgramID(feeQuoterID)
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.FeeQuoter,
		solana.PublicKey{})
	if err != nil {
		return txns, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
	for remoteChainSel, update := range updates {
		ixns := make([]solana.Instruction, 0)
		// verified while loading state
		fqRemoteChainPDA, _, _ := solState.FindFqDestChainPDA(remoteChainSel, feeQuoterID)
		if !update.IsUpdate {
			lookUpTableEntries = append(lookUpTableEntries,
				fqRemoteChainPDA,
			)
		}

		var feeQuoterIx solana.Instruction
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
			tx, err := BuildMCMSTxn(feeQuoterIx, feeQuoterID.String(), ccipChangeset.FeeQuoter)
			if err != nil {
				return txns, fmt.Errorf("failed to create transaction: %w", err)
			}
			txns = append(txns, *tx)
		} else {
			ixns = append(ixns, feeQuoterIx)
		}

		// confirm ixns if any
		if len(ixns) > 0 {
			err = chain.Confirm(ixns)
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
	UpdatesByChain map[uint64]OffRampConfig
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be achieved by calling this function multiple times)
	MCMSSolana *MCMSConfigSolana
}

type OffRampConfig struct {
	// source
	EnabledAsSource bool
	// We have different instructions for add vs update, so we need to know which one to use
	IsUpdate bool
}

func (cfg AddRemoteChainToOffRampConfig) Validate(e deployment.Environment) error {
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.SolChains[cfg.ChainSelector]
	chain, ok := e.SolChains[cfg.ChainSelector]
	if !ok {
		return fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}

	if err := validateOffRampConfig(chain, chainState); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, solana.PublicKey{}); err != nil {
		return err
	}

	supportedChains := state.SupportedChains()
	for remote := range cfg.UpdatesByChain {
		if _, ok := supportedChains[remote]; !ok {
			return fmt.Errorf("remote chain %d is not supported", remote)
		}
		if err := state.ValidateRamp(remote, ccipChangeset.OnRamp); err != nil {
			return err
		}
		offRampRemoteStatePDA, _, err := solState.FindOfframpSourceChainPDA(remote, chainState.OffRamp)
		if err != nil {
			return fmt.Errorf("failed to find dest chain state pda for remote chain %d: %w", remote, err)
		}
		if !cfg.UpdatesByChain[remote].IsUpdate {
			var destChainStateAccount solOffRamp.SourceChain
			err = chain.GetAccountDataBorshInto(context.Background(), offRampRemoteStatePDA, &destChainStateAccount)
			if err == nil {
				return fmt.Errorf("remote %d is already configured on solana chain offramp %d", remote, cfg.ChainSelector)
			}
		}
	}
	return nil
}

// Adds new remote chain configurations
func AddRemoteChainToOffRamp(e deployment.Environment, cfg AddRemoteChainToOffRampConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	s, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	ab := deployment.NewMemoryAddressBook()
	txns, err := doAddRemoteChainToSolana(e, s, cfg, ab)
	if err != nil {
		return deployment.ChangesetOutput{AddressBook: ab}, err
	}

	// create proposals for ixns
	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to add remote chains to Solana", cfg.MCMSSolana.MCMS.MinDelay, txns)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
			AddressBook:           ab,
		}, nil
	}
	return deployment.ChangesetOutput{AddressBook: ab}, nil
}

func doAddRemoteChainToSolana(
	e deployment.Environment,
	s ccipChangeset.CCIPOnChainState,
	cfg AddRemoteChainToOffRampConfig,
	ab deployment.AddressBook) ([]mcmsTypes.Transaction, error) {
	txns := make([]mcmsTypes.Transaction, 0)
	chainSel := cfg.ChainSelector
	updates := cfg.UpdatesByChain
	chain := e.SolChains[chainSel]
	offRampID := s.SolChains[chainSel].OffRamp
	offRampUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.OffRampOwnedByTimelock
	lookUpTableEntries := make([]solana.PublicKey, 0)
	solOffRamp.SetProgramID(offRampID)
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.OffRamp,
		solana.PublicKey{})
	if err != nil {
		return txns, fmt.Errorf("failed to get authority for ixn: %w", err)
	}

	for remoteChainSel, update := range updates {
		ixns := make([]solana.Instruction, 0)
		// verified while loading state
		offRampRemoteStatePDA, _, _ := solState.FindOfframpSourceChainPDA(remoteChainSel, offRampID)
		if !update.IsUpdate {
			lookUpTableEntries = append(lookUpTableEntries,
				offRampRemoteStatePDA,
			)
		}

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
			e.Logger.Infow("update offramp config for remote chain", "remoteChainSel", remoteChainSel)
		} else {
			offRampIx, err = solOffRamp.NewAddSourceChainInstruction(
				remoteChainSel,
				validSourceChainConfig,
				offRampRemoteStatePDA,
				s.SolChains[chainSel].OffRampConfigPDA,
				authority,
				solana.SystemProgramID,
			).ValidateAndBuild()
			e.Logger.Infow("add offramp config for remote chain", "remoteChainSel", remoteChainSel)
		}
		if err != nil {
			return txns, fmt.Errorf("failed to generate instructions: %w", err)
		}
		if offRampUsingMCMS {
			tx, err := BuildMCMSTxn(offRampIx, offRampID.String(), ccipChangeset.OffRamp)
			if err != nil {
				return txns, fmt.Errorf("failed to create transaction: %w", err)
			}
			txns = append(txns, *tx)
		} else {
			ixns = append(ixns, offRampIx)
		}

		// confirm ixns if any
		if len(ixns) > 0 {
			err = chain.Confirm(ixns)
			if err != nil {
				return txns, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}

		if !update.IsUpdate {
			remoteChainSelStr := strconv.FormatUint(remoteChainSel, 10)
			tv := deployment.NewTypeAndVersion(ccipChangeset.RemoteSource, deployment.Version1_0_0)
			tv.AddLabel(remoteChainSelStr)
			err = ab.Save(chainSel, offRampRemoteStatePDA.String(), tv)
			if err != nil {
				return txns, fmt.Errorf("failed to save source chain state to address book: %w", err)
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

func getSourceChainConfig(s ccipChangeset.CCIPOnChainState, remoteChainSel uint64, enabledAsSource bool) (solOffRamp.SourceChainConfig, error) {
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

func extendLookupTable(e deployment.Environment, chain deployment.SolChain, offRampID solana.PublicKey, lookUpTableEntries []solana.PublicKey) error {
	addressLookupTable, err := ccipChangeset.FetchOfframpLookupTable(e.GetContext(), chain, offRampID)
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
