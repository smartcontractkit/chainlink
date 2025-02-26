package solana

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/gagliardetto/solana-go"

	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	"github.com/smartcontractkit/mcms"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

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
	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.ChainSelector, cfg.MCMSSolana); err != nil {
		return err
	}
	routerUsingMcms := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock
	if err := ccipChangeset.ValidateOwnershipSolana(&e, chain, routerUsingMcms, chainState.Router, ccipChangeset.Router); err != nil {
		return fmt.Errorf("failed to validate ownership: %w", err)
	}
	tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)
	if err != nil {
		return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), chainState.Router.String(), err)
	}
	var tokenAdminRegistryAccount solRouter.TokenAdminRegistry
	if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err == nil {
		return fmt.Errorf("token admin registry already exists for (mint: %s, router: %s)", tokenPubKey.String(), chainState.Router.String())
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

	// verified
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)
	tokenAdminRegistryAdmin := solana.MustPublicKeyFromBase58(cfg.TokenAdminRegistryAdmin)

	var instruction *solRouter.Instruction
	var err error
	routerUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock
	// var authority solana.PublicKey
	// if routerUsingMCMS {
	// 	authority, err = FetchTimelockSigner(e, cfg.ChainSelector)
	// 	if err != nil {
	// 		return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
	// 	}
	// } else {
	// 	authority = chain.DeployerKey.PublicKey()
	// }
	switch cfg.RegisterType {
	// the ccip admin signs and makes tokenAdminRegistryAdmin the authority of the tokenAdminRegistry PDA
	case ViaGetCcipAdminInstruction:
		instruction, err = solRouter.NewCcipAdminProposeAdministratorInstruction(
			tokenAdminRegistryAdmin, // admin of the tokenAdminRegistry PDA
			chainState.RouterConfigPDA,
			tokenAdminRegistryPDA, // this gets created
			tokenPubKey,
			chain.DeployerKey.PublicKey(),
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
	case ViaOwnerInstruction:
		// the token mint authority signs and makes itself the authority of the tokenAdminRegistry PDA
		instruction, err = solRouter.NewOwnerProposeAdministratorInstruction(
			tokenAdminRegistryAdmin, // admin of the tokenAdminRegistry PDA
			chainState.RouterConfigPDA,
			tokenAdminRegistryPDA, // this gets created
			tokenPubKey,
			chain.DeployerKey.PublicKey(), // (token mint authority) becomes the authority of the tokenAdminRegistry PDA
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
	}
	if routerUsingMCMS {
		ixn := instruction
		programID := chainState.Router.String()
		contractType := ccipChangeset.Router
		// tx, err := BuildMCMSTxn(instruction, chainState.Router.String(), ccipChangeset.Router)
		// if err != nil {
		// 	return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		// }
		data, err := ixn.Data()
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to extract data: %w", err)
		}
		tx, err := mcmsSolana.NewTransaction(
			programID,
			data,
			big.NewInt(0),        // e.g. value
			ixn.Accounts(),       // pass along needed accounts
			string(contractType), // some string identifying the target
			[]string{},           // any relevant metadata
		)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to RegisterTokenAdminRegistry in Solana", cfg.MCMSSolana.MCMS.MinDelay, []mcmsTypes.Transaction{tx})
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
	// CurrentRegistryAdminPrivateKey string
	MCMSSolana *MCMSConfigSolana
}

func (cfg TransferAdminRoleTokenAdminRegistryConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}

	timelockSigner, err := FetchTimelockSigner(e, cfg.ChainSelector)
	if err != nil {
		return fmt.Errorf("failed to fetch timelock signer: %w", err)
	}
	newRegistryAdminPubKey := solana.MustPublicKeyFromBase58(cfg.NewRegistryAdminPublicKey)

	if timelockSigner.Equals(newRegistryAdminPubKey) {
		return fmt.Errorf("new registry admin public key (%s) cannot be the same as current registry admin public key (%s) for token %s",
			newRegistryAdminPubKey.String(),
			timelockSigner.String(),
			tokenPubKey.String(),
		)
	}

	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.ChainSelector, cfg.MCMSSolana); err != nil {
		return err
	}
	routerUsingMcms := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock
	if err := ccipChangeset.ValidateOwnershipSolana(&e, chain, routerUsingMcms, chainState.Router, ccipChangeset.Router); err != nil {
		return fmt.Errorf("failed to validate ownership: %w", err)
	}
	tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)
	if err != nil {
		return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), chainState.Router.String(), err)
	}
	var tokenAdminRegistryAccount solRouter.TokenAdminRegistry
	if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err != nil {
		return fmt.Errorf("token admin registry not found for (mint: %s, router: %s), cannot transfer admin role", tokenPubKey.String(), chainState.Router.String())
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

	// verified
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)

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
		chainState.RouterConfigPDA,
		tokenAdminRegistryPDA,
		tokenPubKey,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}
	if routerUsingMCMS {
		tx, err := BuildMCMSTxn(ix1, chainState.Router.String(), ccipChangeset.Router)
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
	// instructions := []solana.Instruction{ix1}
	// // the existing authority will have to sign the transfer
	// if err := chain.Confirm(instructions, solCommonUtil.AddSigners(currentRegistryAdminPrivateKey)); err != nil {
	// 	return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	// }
	return deployment.ChangesetOutput{}, nil
}

// ACCEPT TOKEN ADMIN REGISTRY
type AcceptAdminRoleTokenAdminRegistryConfig struct {
	ChainSelector uint64
	TokenPubKey   string
	// NewRegistryAdminPrivateKey string
	MCMSSolana *MCMSConfigSolana
}

func (cfg AcceptAdminRoleTokenAdminRegistryConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.ChainSelector, cfg.MCMSSolana); err != nil {
		return err
	}
	routerUsingMcms := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock
	if err := ccipChangeset.ValidateOwnershipSolana(&e, chain, routerUsingMcms, chainState.Router, ccipChangeset.Router); err != nil {
		return fmt.Errorf("failed to validate ownership: %w", err)
	}
	tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)
	if err != nil {
		return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), chainState.Router.String(), err)
	}
	var tokenAdminRegistryAccount solRouter.TokenAdminRegistry
	if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err != nil {
		return fmt.Errorf("token admin registry not found for (mint: %s, router: %s), cannot accept admin role", tokenPubKey.String(), chainState.Router.String())
	}
	// check if accepting admin is the pending admin
	newRegistryAdminPublicKey, err := FetchTimelockSigner(e, cfg.ChainSelector)
	if err != nil {
		return fmt.Errorf("failed to fetch timelock signer: %w", err)
	}
	if !tokenAdminRegistryAccount.PendingAdministrator.Equals(newRegistryAdminPublicKey) {
		return fmt.Errorf("new admin public key (%s) does not match pending registry admin role (%s) for token %s",
			newRegistryAdminPublicKey.String(),
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
	// chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)

	// verified
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)
	newRegistryAdminPublicKey, err := FetchTimelockSigner(e, cfg.ChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
	}

	routerUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock
	ix1, err := solRouter.NewAcceptAdminRoleTokenAdminRegistryInstruction(
		chainState.RouterConfigPDA,
		tokenAdminRegistryPDA,
		tokenPubKey,
		newRegistryAdminPublicKey,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}
	if routerUsingMCMS {
		tx, err := BuildMCMSTxn(ix1, chainState.Router.String(), ccipChangeset.Router)
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

	// instructions := []solana.Instruction{ix1}
	// // the new authority will have to sign the acceptance
	// if err := chain.Confirm(instructions, solCommonUtil.AddSigners(newRegistryAdminPrivateKey)); err != nil {
	// 	return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	// }
	return deployment.ChangesetOutput{}, nil
}
