package solana

import (
	"encoding/binary"
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solRmnRemote "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/rmn_remote"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

var (
	GlobalCurse = solRmnRemote.CurseSubject{
		Value: [16]uint8{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
	}
)

type CurseConfig struct {
	ChainSelector       uint64
	GlobalCurse         bool
	RemoteChainSelector uint64
	MCMSSolana          *MCMSConfigSolana
}

func (cfg CurseConfig) Validate(e deployment.Environment) error {
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	if cfg.GlobalCurse && cfg.RemoteChainSelector != 0 {
		return fmt.Errorf("remote chain selector must be 0 if global curse is true. chain: %d", cfg.ChainSelector)
	}
	if !cfg.GlobalCurse && cfg.RemoteChainSelector == 0 {
		return fmt.Errorf("remote chain selector must be non-zero if global curse is false. chain: %d", cfg.ChainSelector)
	}
	if !cfg.GlobalCurse {
		routerDestChainPDA, err := solState.FindDestChainStatePDA(cfg.RemoteChainSelector, chainState.Router)
		if err != nil {
			return fmt.Errorf("failed to find dest chain state pda for remote chain %d: %w", cfg.RemoteChainSelector, err)
		}
		var destChainStateAccount solRouter.DestChain
		if err = chain.GetAccountDataBorshInto(e.GetContext(), routerDestChainPDA, &destChainStateAccount); err != nil {
			return fmt.Errorf("remote %d does not seem to be supported by the router, hence cursing is not possible", cfg.RemoteChainSelector)
		}
	}
	return ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, solana.PublicKey{})
}

func curseExists(e deployment.Environment, chain deployment.SolChain, curseSubject solRmnRemote.CurseSubject, rmnRemoteCursesPDA, rmnRemoteConfigPDA solana.PublicKey) (bool, error) {
	ix, err := solRmnRemote.NewVerifyNotCursedInstruction(
		curseSubject,
		rmnRemoteCursesPDA,
		rmnRemoteConfigPDA,
	).ValidateAndBuild()
	if err != nil {
		return false, fmt.Errorf("failed to generate instructions: %w", err)
	}
	if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
		e.Logger.Infof("Curse already exists for chain %d and curse subject %v", chain.Selector, curseSubject)
		return true, nil
	}
	return false, nil
}

func ApplyCurse(e deployment.Environment, cfg CurseConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to validate curse config: %w", err)
	}
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	rmnRemoteConfigPDA := chainState.RMNRemoteConfigPDA
	solRmnRemote.SetProgramID(chainState.RMNRemote)
	rmnRemoteCursesPDA := chainState.RMNRemoteCursesPDA

	// validate curse subject
	var curseSubject solRmnRemote.CurseSubject
	if cfg.GlobalCurse {
		curseSubject = GlobalCurse
		e.Logger.Info("Checking if global curse is already in place")
		exists, err := curseExists(e, chain, GlobalCurse, rmnRemoteCursesPDA, rmnRemoteConfigPDA)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to check if global curse exists: %w", err)
		}
		if exists {
			e.Logger.Infof("Global curse is already in place for chain %d", cfg.ChainSelector)
			return deployment.ChangesetOutput{}, fmt.Errorf("global curse already exists for chain %d", cfg.ChainSelector)
		}
		e.Logger.Info("Applying global curse")
	} else {
		binary.LittleEndian.PutUint64(curseSubject.Value[:], cfg.RemoteChainSelector)
		e.Logger.Infow("Checking if remote chain curse is already in place", "remoteChainSelector", cfg.RemoteChainSelector)
		// this will also fail if global curse exists
		exists, err := curseExists(e, chain, curseSubject, rmnRemoteCursesPDA, rmnRemoteConfigPDA)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to check if remote chain curse exists: %w", err)
		}
		if exists {
			e.Logger.Infow("Remote chain curse is already in place", "remoteChainSelector", cfg.RemoteChainSelector)
			return deployment.ChangesetOutput{}, fmt.Errorf("remote chain curse already exists for chain %d", cfg.RemoteChainSelector)
		}
		e.Logger.Infow("Applying remote chain curse for chain", "remoteChainSelector", cfg.RemoteChainSelector)
	}

	rmnRemoteUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.RMNRemoteOwnedByTimelock
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.RMNRemote,
		solana.PublicKey{})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
	ix, err := solRmnRemote.NewCurseInstruction(
		curseSubject,
		rmnRemoteConfigPDA,
		authority,
		rmnRemoteCursesPDA,
		solana.SystemProgramID,
	).ValidateAndBuild()

	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}

	if rmnRemoteUsingMCMS {
		tx, err := BuildMCMSTxn(ix, chainState.RMNRemote.String(), ccipChangeset.RMNRemote)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to ApplyCurse in Solana", cfg.MCMSSolana.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}

	return deployment.ChangesetOutput{}, nil
}

func RemoveCurse(e deployment.Environment, cfg CurseConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to validate curse config: %w", err)
	}
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	rmnRemoteConfigPDA := chainState.RMNRemoteConfigPDA
	solRmnRemote.SetProgramID(chainState.RMNRemote)
	rmnRemoteCursesPDA := chainState.RMNRemoteCursesPDA

	// validate curse subject
	var curseSubject solRmnRemote.CurseSubject
	if cfg.GlobalCurse {
		curseSubject = GlobalCurse
		e.Logger.Info("Checking if global curse is already in place")
		exists, err := curseExists(e, chain, GlobalCurse, rmnRemoteCursesPDA, rmnRemoteConfigPDA)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to check if global curse exists: %w", err)
		}
		if !exists {
			e.Logger.Info("Global curse is not in place")
			return deployment.ChangesetOutput{}, fmt.Errorf("global curse is not in place for chain %d", cfg.ChainSelector)
		}
		e.Logger.Info("Applying global curse")
	} else {
		binary.LittleEndian.PutUint64(curseSubject.Value[:], cfg.RemoteChainSelector)
		e.Logger.Infow("Checking if remote chain curse is already in place", "remoteChainSelector", cfg.RemoteChainSelector)
		// this will also fail if global curse exists
		exists, err := curseExists(e, chain, curseSubject, rmnRemoteCursesPDA, rmnRemoteConfigPDA)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to check if remote chain curse exists: %w", err)
		}
		if !exists {
			e.Logger.Infow("Remote chain curse is not in place", "remoteChainSelector", cfg.RemoteChainSelector)
			return deployment.ChangesetOutput{}, fmt.Errorf("remote chain curse is not in place for chain %d", cfg.RemoteChainSelector)
		}
		e.Logger.Infow("Removing remote chain curse for chain", "remoteChainSelector", cfg.RemoteChainSelector)
	}

	rmnRemoteUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.RMNRemoteOwnedByTimelock
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.RMNRemote,
		solana.PublicKey{})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
	ix, err := solRmnRemote.NewUncurseInstruction(
		curseSubject,
		rmnRemoteConfigPDA,
		authority,
		rmnRemoteCursesPDA,
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}
	if rmnRemoteUsingMCMS {
		tx, err := BuildMCMSTxn(ix, chainState.RMNRemote.String(), ccipChangeset.RMNRemote)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to RemoveCurse in Solana", cfg.MCMSSolana.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}
	if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return deployment.ChangesetOutput{}, nil
}
