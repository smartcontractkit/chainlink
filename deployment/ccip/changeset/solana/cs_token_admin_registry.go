package solana

import (
	"context"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solCommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_common"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// use these changesets to register a token admin registry, transfer the admin role, and accept the admin role
var _ cldf.ChangeSet[RegisterTokenAdminRegistryConfig] = RegisterTokenAdminRegistry
var _ cldf.ChangeSet[TransferAdminRoleTokenAdminRegistryConfig] = TransferAdminRoleTokenAdminRegistry
var _ cldf.ChangeSet[AcceptAdminRoleTokenAdminRegistryConfig] = AcceptAdminRoleTokenAdminRegistry

type RegisterTokenAdminRegistryType int

const (
	ViaGetCcipAdminInstruction RegisterTokenAdminRegistryType = iota
	ViaOwnerInstruction
)

type RegisterTokenAdminRegistryConfig struct {
	ChainSelector           uint64
	TokenPubKey             solana.PublicKey
	TokenAdminRegistryAdmin string
	RegisterType            RegisterTokenAdminRegistryType
	Override                bool
	MCMS                    *proposalutils.TimelockConfig
}

func (cfg RegisterTokenAdminRegistryConfig) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	if cfg.RegisterType != ViaGetCcipAdminInstruction && cfg.RegisterType != ViaOwnerInstruction {
		return fmt.Errorf("invalid register type, valid types are %d and %d", ViaGetCcipAdminInstruction, ViaOwnerInstruction)
	}

	if cfg.TokenAdminRegistryAdmin == "" {
		return errors.New("token admin registry admin is required")
	}

	tokenPubKey := cfg.TokenPubKey
	if err := chainState.CommonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateRouterConfig(chain); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.Router: true}); err != nil {
		return err
	}
	routerProgramAddress, _, _ := chainState.GetRouterInfo()
	tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
	if err != nil {
		return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), routerProgramAddress.String(), err)
	}
	var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
	if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err == nil {
		return fmt.Errorf("token admin registry already exists for (mint: %s, router: %s)", tokenPubKey.String(), routerProgramAddress.String())
	}
	return nil
}

func RegisterTokenAdminRegistry(e cldf.Environment, cfg RegisterTokenAdminRegistryConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("RegisterTokenAdminRegistry", "cfg", cfg)
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chainState := state.SolChains[cfg.ChainSelector]
	if err := cfg.Validate(e, chainState); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	tokenPubKey := cfg.TokenPubKey
	routerProgramAddress, routerConfigPDA, _ := chainState.GetRouterInfo()
	solRouter.SetProgramID(routerProgramAddress)

	// verified
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
	tokenAdminRegistryAdmin := solana.MustPublicKeyFromBase58(cfg.TokenAdminRegistryAdmin)

	var instruction *solRouter.Instruction
	routerUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.Router,
		solana.PublicKey{},
		"")
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.Router,
		solana.PublicKey{},
		"")
	switch cfg.RegisterType {
	// the ccip admin signs and makes tokenAdminRegistryAdmin the authority of the tokenAdminRegistry PDA
	case ViaGetCcipAdminInstruction:
		if cfg.Override {
			instruction, err = solRouter.NewCcipAdminOverridePendingAdministratorInstruction(
				tokenAdminRegistryAdmin, // admin of the tokenAdminRegistry PDA
				routerConfigPDA,
				tokenAdminRegistryPDA, // this gets created
				tokenPubKey,
				authority,
				solana.SystemProgramID,
			).ValidateAndBuild()
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
		} else {
			instruction, err = solRouter.NewCcipAdminProposeAdministratorInstruction(
				tokenAdminRegistryAdmin, // admin of the tokenAdminRegistry PDA
				routerConfigPDA,
				tokenAdminRegistryPDA, // this gets created
				tokenPubKey,
				authority,
				solana.SystemProgramID,
			).ValidateAndBuild()
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
		}
	case ViaOwnerInstruction:
		if cfg.Override {
			instruction, err = solRouter.NewOwnerOverridePendingAdministratorInstruction(
				tokenAdminRegistryAdmin, // admin of the tokenAdminRegistry PDA
				routerConfigPDA,
				tokenAdminRegistryPDA, // this gets created
				tokenPubKey,
				authority,
				solana.SystemProgramID,
			).ValidateAndBuild()
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
		} else {
			// the token mint authority signs and makes itself the authority of the tokenAdminRegistry PDA
			instruction, err = solRouter.NewOwnerProposeAdministratorInstruction(
				tokenAdminRegistryAdmin, // admin of the tokenAdminRegistry PDA
				routerConfigPDA,
				tokenAdminRegistryPDA, // this gets created
				tokenPubKey,
				authority, // (token mint authority) becomes the authority of the tokenAdminRegistry PDA
				solana.SystemProgramID,
			).ValidateAndBuild()
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
		}
	}
	if routerUsingMCMS {
		tx, err := BuildMCMSTxn(instruction, routerProgramAddress.String(), shared.Router)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to RegisterTokenAdminRegistry in Solana", cfg.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}
	// if we want to have a different authority, we will need to add the corresponding signer here
	// for now we are assuming both token owner and ccip admin will always be deployer key
	instructions := []solana.Instruction{instruction}
	if err := chain.Confirm(instructions); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return cldf.ChangesetOutput{}, nil
}

// TRANSFER AND ACCEPT TOKEN ADMIN REGISTRY
type TransferAdminRoleTokenAdminRegistryConfig struct {
	ChainSelector             uint64
	TokenPubKey               string
	NewRegistryAdminPublicKey string
	MCMS                      *proposalutils.TimelockConfig
}

func (cfg TransferAdminRoleTokenAdminRegistryConfig) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	if err := chainState.CommonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateRouterConfig(chain); err != nil {
		return err
	}
	currentAdmin := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.Router,
		solana.PublicKey{},
		"",
	)

	newRegistryAdminPubKey := solana.MustPublicKeyFromBase58(cfg.NewRegistryAdminPublicKey)

	if currentAdmin.Equals(newRegistryAdminPubKey) {
		return fmt.Errorf("new registry admin public key (%s) cannot be the same as current registry admin public key (%s) for token %s",
			newRegistryAdminPubKey.String(),
			currentAdmin.String(),
			tokenPubKey.String(),
		)
	}

	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.Router: true}); err != nil {
		return err
	}
	routerProgramAddress, _, _ := chainState.GetRouterInfo()
	tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
	if err != nil {
		return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), routerProgramAddress.String(), err)
	}
	var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
	if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err != nil {
		return fmt.Errorf("token admin registry not found for (mint: %s, router: %s), cannot transfer admin role", tokenPubKey.String(), routerProgramAddress.String())
	}
	return nil
}

func TransferAdminRoleTokenAdminRegistry(e cldf.Environment, cfg TransferAdminRoleTokenAdminRegistryConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chainState := state.SolChains[cfg.ChainSelector]
	if err := cfg.Validate(e, chainState); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	routerProgramAddress, routerConfigPDA, _ := chainState.GetRouterInfo()
	solRouter.SetProgramID(routerProgramAddress)
	// verified
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
	newRegistryAdminPubKey := solana.MustPublicKeyFromBase58(cfg.NewRegistryAdminPublicKey)

	routerUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.Router,
		solana.PublicKey{},
		"")
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.Router,
		solana.PublicKey{},
		"")
	ix1, err := solRouter.NewTransferAdminRoleTokenAdminRegistryInstruction(
		newRegistryAdminPubKey,
		routerConfigPDA,
		tokenAdminRegistryPDA,
		tokenPubKey,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}
	if routerUsingMCMS {
		tx, err := BuildMCMSTxn(ix1, routerProgramAddress.String(), shared.Router)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to TransferAdminRoleTokenAdminRegistry in Solana", cfg.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}
	// the existing authority will have to sign the transfer
	if err := chain.Confirm([]solana.Instruction{ix1}); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return cldf.ChangesetOutput{}, nil
}

// ACCEPT TOKEN ADMIN REGISTRY
type AcceptAdminRoleTokenAdminRegistryConfig struct {
	ChainSelector     uint64
	TokenPubKey       solana.PublicKey
	MCMS              *proposalutils.TimelockConfig
	SkipRegistryCheck bool // set to true when you want to register and accept in the same proposal
}

func (cfg AcceptAdminRoleTokenAdminRegistryConfig) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	tokenPubKey := cfg.TokenPubKey
	if err := chainState.CommonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateRouterConfig(chain); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.Router: true}); err != nil {
		return err
	}

	newAdmin := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.Router,
		solana.PublicKey{},
		"",
	)

	if !cfg.SkipRegistryCheck {
		routerProgramAddress, _, _ := chainState.GetRouterInfo()
		tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
		if err != nil {
			return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), routerProgramAddress.String(), err)
		}
		var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
		if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err != nil {
			return fmt.Errorf("token admin registry not found for (mint: %s, router: %s), cannot accept admin role", tokenPubKey.String(), routerProgramAddress.String())
		}
		if !tokenAdminRegistryAccount.PendingAdministrator.Equals(newAdmin) {
			return fmt.Errorf("new admin public key (%s) does not match pending registry admin role (%s) for token %s",
				newAdmin.String(),
				tokenAdminRegistryAccount.PendingAdministrator.String(),
				tokenPubKey.String(),
			)
		}
	}
	return nil
}

func AcceptAdminRoleTokenAdminRegistry(e cldf.Environment, cfg AcceptAdminRoleTokenAdminRegistryConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("AcceptAdminRoleTokenAdminRegistry", "cfg", cfg)
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chainState := state.SolChains[cfg.ChainSelector]
	if err := cfg.Validate(e, chainState); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	tokenPubKey := cfg.TokenPubKey
	// verified
	routerProgramAddress, routerConfigPDA, _ := chainState.GetRouterInfo()
	solRouter.SetProgramID(routerProgramAddress)
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)

	routerUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.Router,
		solana.PublicKey{},
		"")
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.Router,
		solana.PublicKey{},
		"")
	ix1, err := solRouter.NewAcceptAdminRoleTokenAdminRegistryInstruction(
		routerConfigPDA,
		tokenAdminRegistryPDA,
		tokenPubKey,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}
	if routerUsingMCMS {
		// We will only be able to accept the admin role if the pending admin is the timelock signer
		tx, err := BuildMCMSTxn(ix1, routerProgramAddress.String(), shared.Router)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to AcceptAdminRoleTokenAdminRegistry in Solana", cfg.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}
	if err := chain.Confirm([]solana.Instruction{ix1}); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return cldf.ChangesetOutput{}, nil
}
