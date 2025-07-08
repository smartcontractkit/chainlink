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
	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"
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

// use this changeset to set pool on token admin registry
var _ cldf.ChangeSet[SetPoolConfig] = SetPool

type RegisterTokenAdminRegistryType int

const (
	ViaGetCcipAdminInstruction RegisterTokenAdminRegistryType = iota
	ViaOwnerInstruction
)

type RegisterTokenConfig struct {
	TokenPubKey             solana.PublicKey
	TokenAdminRegistryAdmin solana.PublicKey
	RegisterType            RegisterTokenAdminRegistryType
	Override                bool
}

type RegisterTokenAdminRegistryConfig struct {
	ChainSelector        uint64
	RegisterTokenConfigs []RegisterTokenConfig
	MCMS                 *proposalutils.TimelockConfig
}

func (cfg RegisterTokenAdminRegistryConfig) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateRouterConfig(chain); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.Router: true}); err != nil {
		return err
	}
	routerProgramAddress, _, _ := chainState.GetRouterInfo()

	for _, registerTokenConfig := range cfg.RegisterTokenConfigs {
		if registerTokenConfig.RegisterType != ViaGetCcipAdminInstruction && registerTokenConfig.RegisterType != ViaOwnerInstruction {
			return fmt.Errorf("invalid register type, valid types are %d and %d", ViaGetCcipAdminInstruction, ViaOwnerInstruction)
		}
		if registerTokenConfig.TokenAdminRegistryAdmin.IsZero() {
			return errors.New("token admin registry admin is required")
		}
		tokenPubKey := registerTokenConfig.TokenPubKey
		if err := chainState.CommonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
			return err
		}
		tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
		if err != nil {
			return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), routerProgramAddress.String(), err)
		}
		var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
		if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err == nil {
			if !registerTokenConfig.Override {
				return fmt.Errorf("token admin registry already exists for (mint: %s, router: %s)", tokenPubKey.String(), routerProgramAddress.String())
			}
		}
	}

	return nil
}

func RegisterTokenAdminRegistry(e cldf.Environment, cfg RegisterTokenAdminRegistryConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("RegisterTokenAdminRegistry", "cfg", cfg)
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chainState, ok := state.SolChains[cfg.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}
	if err := cfg.Validate(e, chainState); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	routerProgramAddress, routerConfigPDA, _ := chainState.GetRouterInfo()
	solRouter.SetProgramID(routerProgramAddress)
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
		shared.Router,
		solana.PublicKey{},
		"")
	mcmsTxs := []mcmsTypes.Transaction{}

	for _, registerTokenConfig := range cfg.RegisterTokenConfigs {
		tokenPubKey := registerTokenConfig.TokenPubKey
		tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
		tokenAdminRegistryAdmin := registerTokenConfig.TokenAdminRegistryAdmin
		var instruction *solRouter.Instruction
		switch registerTokenConfig.RegisterType {
		// the ccip admin signs and makes tokenAdminRegistryAdmin the authority of the tokenAdminRegistry PDA
		case ViaGetCcipAdminInstruction:
			if registerTokenConfig.Override {
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
			if registerTokenConfig.Override {
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

		// if mcms build the transaction
		// else just confirm it
		if routerUsingMCMS {
			tx, err := BuildMCMSTxn(instruction, routerProgramAddress.String(), shared.Router)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
			}
			mcmsTxs = append(mcmsTxs, *tx)

		} else {
			// if we want to have a different authority, we will need to add the corresponding signer here
			// for now we are assuming both token owner and ccip admin will always be deployer key if done without mcms
			instructions := []solana.Instruction{instruction}
			if err := chain.Confirm(instructions); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}
	}

	if len(mcmsTxs) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to RegisterTokenAdminRegistry in Solana", cfg.MCMS.MinDelay, mcmsTxs)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}

type TrasnferTokenAdminConfig struct {
	TokenPubKey               solana.PublicKey
	NewRegistryAdminPublicKey solana.PublicKey
}

// TRANSFER AND ACCEPT TOKEN ADMIN REGISTRY
type TransferAdminRoleTokenAdminRegistryConfig struct {
	ChainSelector             uint64
	TransferTokenAdminConfigs []TrasnferTokenAdminConfig
	MCMS                      *proposalutils.TimelockConfig
}

func (cfg TransferAdminRoleTokenAdminRegistryConfig) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateRouterConfig(chain); err != nil {
		return err
	}
	routerProgramAddress, _, _ := chainState.GetRouterInfo()
	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.Router: true}); err != nil {
		return err
	}

	for _, transferTokenAdminConfig := range cfg.TransferTokenAdminConfigs {
		tokenPubKey := transferTokenAdminConfig.TokenPubKey
		if err := chainState.CommonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
			return err
		}
		newRegistryAdminPubKey := transferTokenAdminConfig.NewRegistryAdminPublicKey
		tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
		if err != nil {
			return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), routerProgramAddress.String(), err)
		}
		var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
		if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err != nil {
			return fmt.Errorf("token admin registry not found for (mint: %s, router: %s), cannot transfer admin role", tokenPubKey.String(), routerProgramAddress.String())
		}
		currentAdmin := tokenAdminRegistryAccount.Administrator
		if currentAdmin.Equals(newRegistryAdminPubKey) {
			return fmt.Errorf("new registry admin public key (%s) cannot be the same as current registry admin public key (%s) for token %s",
				newRegistryAdminPubKey.String(),
				currentAdmin.String(),
				tokenPubKey.String(),
			)
		}
	}

	return nil
}

func TransferAdminRoleTokenAdminRegistry(e cldf.Environment, cfg TransferAdminRoleTokenAdminRegistryConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("TransferAdminRoleTokenAdminRegistry", "cfg", cfg)
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chainState, ok := state.SolChains[cfg.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}
	if err := cfg.Validate(e, chainState); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	routerProgramAddress, routerConfigPDA, _ := chainState.GetRouterInfo()
	solRouter.SetProgramID(routerProgramAddress)

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
		shared.Router,
		solana.PublicKey{},
		"")

	mcmsTxs := []mcmsTypes.Transaction{}

	for _, transferTokenAdminConfig := range cfg.TransferTokenAdminConfigs {
		tokenPubKey := transferTokenAdminConfig.TokenPubKey
		tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
		newRegistryAdminPubKey := transferTokenAdminConfig.NewRegistryAdminPublicKey
		instruction, err := solRouter.NewTransferAdminRoleTokenAdminRegistryInstruction(
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
			tx, err := BuildMCMSTxn(instruction, routerProgramAddress.String(), shared.Router)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
			}
			mcmsTxs = append(mcmsTxs, *tx)
		} else {
			instructions := []solana.Instruction{instruction}
			if err := chain.Confirm(instructions); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}
	}

	if routerUsingMCMS {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to TransferAdminRoleTokenAdminRegistry in Solana", cfg.MCMS.MinDelay, mcmsTxs)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}

// ACCEPT TOKEN ADMIN REGISTRY

type AcceptAdminRoleTokenConfig struct {
	TokenPubKey       solana.PublicKey
	SkipRegistryCheck bool
}

type AcceptAdminRoleTokenAdminRegistryConfig struct {
	ChainSelector               uint64
	AcceptAdminRoleTokenConfigs []AcceptAdminRoleTokenConfig
	MCMS                        *proposalutils.TimelockConfig
}

func (cfg AcceptAdminRoleTokenAdminRegistryConfig) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateRouterConfig(chain); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.Router: true}); err != nil {
		return err
	}
	routerProgramAddress, _, _ := chainState.GetRouterInfo()

	for _, acceptAdminRoleTokenConfig := range cfg.AcceptAdminRoleTokenConfigs {
		tokenPubKey := acceptAdminRoleTokenConfig.TokenPubKey
		if err := chainState.CommonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
			return err
		}
		// can only be deployer key or timelock signer
		newAdmin := chain.DeployerKey.PublicKey()
		if cfg.MCMS != nil {
			timelockSignerPDA, err := FetchTimelockSigner(e, cfg.ChainSelector)
			if err != nil {
				return fmt.Errorf("failed to fetch timelock signer: %w", err)
			}
			newAdmin = timelockSignerPDA
		}
		if !acceptAdminRoleTokenConfig.SkipRegistryCheck {
			tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
			if err != nil {
				return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), routerProgramAddress.String(), err)
			}
			var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
			if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err != nil {
				return fmt.Errorf("token admin registry not found for (mint: %s, router: %s), cannot accept admin role", tokenPubKey.String(), routerProgramAddress.String())
			}
			// this will be hit if
			// you register with timelock but accept without mcms config
			// register with deployer key but accept with mcms config
			if !tokenAdminRegistryAccount.PendingAdministrator.Equals(newAdmin) {
				return fmt.Errorf("new admin public key (%s) does not match pending registry admin role (%s) for token %s",
					newAdmin.String(),
					tokenAdminRegistryAccount.PendingAdministrator.String(),
					tokenPubKey.String(),
				)
			}
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
	chainState, ok := state.SolChains[cfg.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}
	if err := cfg.Validate(e, chainState); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	routerUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.Router,
		solana.PublicKey{},
		"")
	authority := chain.DeployerKey.PublicKey()
	if cfg.MCMS != nil {
		timelockSignerPDA, err := FetchTimelockSigner(e, cfg.ChainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
		}
		authority = timelockSignerPDA
	}
	// verified
	routerProgramAddress, routerConfigPDA, _ := chainState.GetRouterInfo()
	solRouter.SetProgramID(routerProgramAddress)
	mcmsTxs := []mcmsTypes.Transaction{}
	for _, acceptAdminRoleTokenConfig := range cfg.AcceptAdminRoleTokenConfigs {
		tokenPubKey := acceptAdminRoleTokenConfig.TokenPubKey
		tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
		instruction, err := solRouter.NewAcceptAdminRoleTokenAdminRegistryInstruction(
			routerConfigPDA,
			tokenAdminRegistryPDA,
			tokenPubKey,
			authority,
		).ValidateAndBuild()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
		if routerUsingMCMS {
			tx, err := BuildMCMSTxn(instruction, routerProgramAddress.String(), shared.Router)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
			}
			mcmsTxs = append(mcmsTxs, *tx)
		} else {
			// pending admin is deployer key
			instructions := []solana.Instruction{instruction}
			if err := chain.Confirm(instructions); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}
	}

	if routerUsingMCMS {
		// We will only be able to accept the admin role if the pending admin is the timelock signer
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to AcceptAdminRoleTokenAdminRegistry in Solana", cfg.MCMS.MinDelay, mcmsTxs)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}

// SET POOL

type SetPoolTokenConfig struct {
	TokenPubKey       solana.PublicKey
	PoolType          *solTestTokenPool.PoolType
	Metadata          string
	SkipRegistryCheck bool // set to true when you want to register and set pool in the same proposal
}

type SetPoolConfig struct {
	ChainSelector       uint64
	SetPoolTokenConfigs []SetPoolTokenConfig
	WritableIndexes     []uint8
	MCMS                *proposalutils.TimelockConfig
}

func (cfg SetPoolConfig) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateRouterConfig(chain); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.Router: true}); err != nil {
		return err
	}
	routerProgramAddress, _, _ := chainState.GetRouterInfo()
	for _, tokenConfig := range cfg.SetPoolTokenConfigs {
		tokenPubKey := tokenConfig.TokenPubKey
		if err := chainState.CommonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
			return err
		}
		if tokenConfig.PoolType == nil {
			return errors.New("pool type must be defined")
		}

		if tokenConfig.Metadata == "" {
			return errors.New("metadata must be defined")
		}
		if lut, ok := chainState.TokenPoolLookupTable[tokenPubKey][*tokenConfig.PoolType][tokenConfig.Metadata]; !ok || lut.IsZero() {
			return fmt.Errorf("token pool lookup table not found for (mint: %s)", tokenPubKey.String())
		}
		if !tokenConfig.SkipRegistryCheck {
			tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
			if err != nil {
				return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), routerProgramAddress.String(), err)
			}
			var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
			if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err != nil {
				return fmt.Errorf("token admin registry not found for (mint: %s, router: %s), cannot set pool", tokenPubKey.String(), routerProgramAddress.String())
			}
		}
	}

	return nil
}

// this sets the writable indexes of the token pool lookup table
func SetPool(e cldf.Environment, cfg SetPoolConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Setting pool config", "cfg", cfg)
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chainState, ok := state.SolChains[cfg.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}
	if err := cfg.Validate(e, chainState); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	routerProgramAddress, routerConfigPDA, _ := chainState.GetRouterInfo()
	solRouter.SetProgramID(routerProgramAddress)
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	routerUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.Router,
		solana.PublicKey{},
		"",
	)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		shared.Router,
		solana.PublicKey{},
		"")
	mcmsTxs := []mcmsTypes.Transaction{}
	for _, tokenConfig := range cfg.SetPoolTokenConfigs {
		tokenPubKey := tokenConfig.TokenPubKey
		tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
		lookupTablePubKey := chainState.TokenPoolLookupTable[tokenPubKey][*tokenConfig.PoolType][tokenConfig.Metadata]
		base := solRouter.NewSetPoolInstruction(
			cfg.WritableIndexes,
			routerConfigPDA,
			tokenAdminRegistryPDA,
			tokenPubKey,
			lookupTablePubKey,
			authority,
		)
		base.AccountMetaSlice = append(base.AccountMetaSlice, solana.Meta(lookupTablePubKey))
		instruction, err := base.ValidateAndBuild()
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		if routerUsingMCMS {
			tx, err := BuildMCMSTxn(instruction, routerProgramAddress.String(), shared.Router)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
			}
			mcmsTxs = append(mcmsTxs, *tx)
		} else {
			if err = chain.Confirm([]solana.Instruction{instruction}); err != nil {
				return cldf.ChangesetOutput{}, err
			}
		}
	}

	if routerUsingMCMS {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to RegisterTokenAdminRegistry in Solana", cfg.MCMS.MinDelay, mcmsTxs)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}
