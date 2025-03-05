package solana

import (
	"context"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

type RegisterTokenAdminRegistryType int

const (
	ViaGetCcipAdminInstruction RegisterTokenAdminRegistryType = iota
	ViaOwnerInstruction
)

type RegisterTokenAdminRegistryConfig struct {
	ChainSelector           uint64
	TokenPubKey             string
	TokenAdminRegistryAdmin string
	RegisterType            RegisterTokenAdminRegistryType
	MCMSSolana              *MCMSConfigSolana
	TestRouter              bool
}

func (cfg RegisterTokenAdminRegistryConfig) Validate(e deployment.Environment) error {
	if cfg.RegisterType != ViaGetCcipAdminInstruction && cfg.RegisterType != ViaOwnerInstruction {
		return fmt.Errorf("invalid register type, valid types are %d and %d", ViaGetCcipAdminInstruction, ViaOwnerInstruction)
	}

	if cfg.TokenAdminRegistryAdmin == "" {
		return errors.New("token admin registry admin is required")
	}

	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateRouterConfig(chain, chainState, cfg.TestRouter); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.ChainSelector, cfg.MCMSSolana); err != nil {
		return err
	}
	routerUsingMcms := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock
	if !cfg.TestRouter {
		if err := ccipChangeset.ValidateOwnershipSolana(&e, chain, routerUsingMcms, chainState.Router, ccipChangeset.Router); err != nil {
			return fmt.Errorf("failed to validate ownership: %w", err)
		}
	}
	routerProgramAddress, _, _ := chainState.GetRouterInfo(cfg.TestRouter)
	tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
	if err != nil {
		return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), routerProgramAddress.String(), err)
	}
	var tokenAdminRegistryAccount solRouter.TokenAdminRegistry
	if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err == nil {
		return fmt.Errorf("token admin registry already exists for (mint: %s, router: %s)", tokenPubKey.String(), routerProgramAddress.String())
	}
	return nil
}

func RegisterTokenAdminRegistry(e deployment.Environment, cfg RegisterTokenAdminRegistryConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	routerProgramAddress, routerConfigPDA, _ := chainState.GetRouterInfo(cfg.TestRouter)
	solRouter.SetProgramID(routerProgramAddress)

	// verified
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
	tokenAdminRegistryAdmin := solana.MustPublicKeyFromBase58(cfg.TokenAdminRegistryAdmin)

	var instruction *solRouter.Instruction
	var err error
	routerUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock
	var authority solana.PublicKey
	if routerUsingMCMS {
		authority, err = FetchTimelockSigner(e, cfg.ChainSelector)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
		}
	} else {
		authority = chain.DeployerKey.PublicKey()
	}
	switch cfg.RegisterType {
	// the ccip admin signs and makes tokenAdminRegistryAdmin the authority of the tokenAdminRegistry PDA
	case ViaGetCcipAdminInstruction:
		instruction, err = solRouter.NewCcipAdminProposeAdministratorInstruction(
			tokenAdminRegistryAdmin, // admin of the tokenAdminRegistry PDA
			routerConfigPDA,
			tokenAdminRegistryPDA, // this gets created
			tokenPubKey,
			authority,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
	case ViaOwnerInstruction:
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
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
	}
	if routerUsingMCMS {
		tx, err := BuildMCMSTxn(instruction, routerProgramAddress.String(), ccipChangeset.Router)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to RegisterTokenAdminRegistry in Solana", cfg.MCMSSolana.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}
	// if we want to have a different authority, we will need to add the corresponding signer here
	// for now we are assuming both token owner and ccip admin will always be deployer key
	instructions := []solana.Instruction{instruction}
	if err := chain.Confirm(instructions); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return deployment.ChangesetOutput{}, nil
}

// TRANSFER AND ACCEPT TOKEN ADMIN REGISTRY
type TransferAdminRoleTokenAdminRegistryConfig struct {
	ChainSelector             uint64
	TokenPubKey               string
	NewRegistryAdminPublicKey string
	MCMSSolana                *MCMSConfigSolana
	TestRouter                bool
}

func (cfg TransferAdminRoleTokenAdminRegistryConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateRouterConfig(chain, chainState, cfg.TestRouter); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.ChainSelector, cfg.MCMSSolana); err != nil {
		return err
	}
	currentAdmin := chain.DeployerKey.PublicKey()
	routerUsingMcms := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock
	var err error
	if routerUsingMcms {
		currentAdmin, err = FetchTimelockSigner(e, cfg.ChainSelector)
		if err != nil {
			return fmt.Errorf("failed to fetch timelock signer: %w", err)
		}
	}

	newRegistryAdminPubKey := solana.MustPublicKeyFromBase58(cfg.NewRegistryAdminPublicKey)

	if currentAdmin.Equals(newRegistryAdminPubKey) {
		return fmt.Errorf("new registry admin public key (%s) cannot be the same as current registry admin public key (%s) for token %s",
			newRegistryAdminPubKey.String(),
			currentAdmin.String(),
			tokenPubKey.String(),
		)
	}

	if !cfg.TestRouter {
		if err := ccipChangeset.ValidateOwnershipSolana(&e, chain, routerUsingMcms, chainState.Router, ccipChangeset.Router); err != nil {
			return fmt.Errorf("failed to validate ownership: %w", err)
		}
	}
	routerProgramAddress, _, _ := chainState.GetRouterInfo(cfg.TestRouter)
	tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
	if err != nil {
		return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), routerProgramAddress.String(), err)
	}
	var tokenAdminRegistryAccount solRouter.TokenAdminRegistry
	if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err != nil {
		return fmt.Errorf("token admin registry not found for (mint: %s, router: %s), cannot transfer admin role", tokenPubKey.String(), routerProgramAddress.String())
	}
	return nil
}

func TransferAdminRoleTokenAdminRegistry(e deployment.Environment, cfg TransferAdminRoleTokenAdminRegistryConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	routerProgramAddress, routerConfigPDA, _ := chainState.GetRouterInfo(cfg.TestRouter)
	solRouter.SetProgramID(routerProgramAddress)
	// verified
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
	newRegistryAdminPubKey := solana.MustPublicKeyFromBase58(cfg.NewRegistryAdminPublicKey)

	routerUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock
	var authority solana.PublicKey
	var err error
	if routerUsingMCMS {
		authority, err = FetchTimelockSigner(e, cfg.ChainSelector)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
		}
	} else {
		authority = e.SolChains[cfg.ChainSelector].DeployerKey.PublicKey()
	}
	ix1, err := solRouter.NewTransferAdminRoleTokenAdminRegistryInstruction(
		newRegistryAdminPubKey,
		routerConfigPDA,
		tokenAdminRegistryPDA,
		tokenPubKey,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}
	if routerUsingMCMS {
		tx, err := BuildMCMSTxn(ix1, routerProgramAddress.String(), ccipChangeset.Router)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to TransferAdminRoleTokenAdminRegistry in Solana", cfg.MCMSSolana.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}
	// the existing authority will have to sign the transfer
	chain := e.SolChains[cfg.ChainSelector]
	if err := chain.Confirm([]solana.Instruction{ix1}); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return deployment.ChangesetOutput{}, nil
}

// ACCEPT TOKEN ADMIN REGISTRY
type AcceptAdminRoleTokenAdminRegistryConfig struct {
	ChainSelector uint64
	TokenPubKey   string
	MCMSSolana    *MCMSConfigSolana
	TestRouter    bool
}

func (cfg AcceptAdminRoleTokenAdminRegistryConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateRouterConfig(chain, chainState, cfg.TestRouter); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.ChainSelector, cfg.MCMSSolana); err != nil {
		return err
	}
	routerUsingMcms := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock
	newAdmin := chain.DeployerKey.PublicKey()
	var err error
	if routerUsingMcms {
		newAdmin, err = FetchTimelockSigner(e, cfg.ChainSelector)
		if err != nil {
			return fmt.Errorf("failed to fetch timelock signer: %w", err)
		}
	}
	if !cfg.TestRouter {
		if err := ccipChangeset.ValidateOwnershipSolana(&e, chain, routerUsingMcms, chainState.Router, ccipChangeset.Router); err != nil {
			return fmt.Errorf("failed to validate ownership: %w", err)
		}
	}
	routerProgramAddress, _, _ := chainState.GetRouterInfo(cfg.TestRouter)
	tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
	if err != nil {
		return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), routerProgramAddress.String(), err)
	}
	var tokenAdminRegistryAccount solRouter.TokenAdminRegistry
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
	return nil
}

func AcceptAdminRoleTokenAdminRegistry(e deployment.Environment, cfg AcceptAdminRoleTokenAdminRegistryConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)

	// verified
	routerProgramAddress, routerConfigPDA, _ := chainState.GetRouterInfo(cfg.TestRouter)
	solRouter.SetProgramID(routerProgramAddress)
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)

	routerUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock
	var authority solana.PublicKey
	var err error
	if routerUsingMCMS {
		authority, err = FetchTimelockSigner(e, cfg.ChainSelector)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
		}
	} else {
		authority = chain.DeployerKey.PublicKey()
	}
	ix1, err := solRouter.NewAcceptAdminRoleTokenAdminRegistryInstruction(
		routerConfigPDA,
		tokenAdminRegistryPDA,
		tokenPubKey,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}
	if routerUsingMCMS {
		// We will only be able to accept the admin role if the pending admin is the timelock signer
		tx, err := BuildMCMSTxn(ix1, routerProgramAddress.String(), ccipChangeset.Router)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to AcceptAdminRoleTokenAdminRegistry in Solana", cfg.MCMSSolana.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}
	if err := chain.Confirm([]solana.Instruction{ix1}); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return deployment.ChangesetOutput{}, nil
}
