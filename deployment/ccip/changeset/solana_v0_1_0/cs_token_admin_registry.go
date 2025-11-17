package solana

import (
	"context"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solCommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_0/ccip_common"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_0/ccip_router"
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
	ccipAdmin := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		shared.Router,
		solana.PublicKey{},
		"")
	timelockSignerPDA, err := FetchTimelockSigner(e, cfg.ChainSelector)
	if err != nil {
		return fmt.Errorf("failed to fetch timelock signer: %w", err)
	}

	for _, registerTokenConfig := range cfg.RegisterTokenConfigs {
		if registerTokenConfig.RegisterType != ViaGetCcipAdminInstruction && registerTokenConfig.RegisterType != ViaOwnerInstruction {
			return fmt.Errorf("invalid register type, valid types are %d", ViaGetCcipAdminInstruction)
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
		if registerTokenConfig.RegisterType == ViaGetCcipAdminInstruction && ccipAdmin.Equals(timelockSignerPDA) && cfg.MCMS == nil {
			return errors.New("ccip admin is the timelock signer, but no mcms config is provided, hence this changeset cannot sign for the registration")
		}
	}

	return nil
}

// RegisterTokenAdminRegistry registers a token admin registry for a given token
// you can register using the ccipAdminRole which can be the deployer key or timelock signer
// you can register using the token mint authority which can be the deployer key only
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

	deployerKey := chain.DeployerKey.PublicKey()
	timelockSignerPDA, err := FetchTimelockSigner(e, cfg.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
	}

	ccipAdmin := GetAuthorityForIxn(
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
		case ViaGetCcipAdminInstruction:
			// the ccip admin signs and makes tokenAdminRegistryAdmin the authority of the tokenAdminRegistry PDA
			if registerTokenConfig.Override {
				instruction, err = solRouter.NewCcipAdminOverridePendingAdministratorInstruction(
					tokenAdminRegistryAdmin, // admin of the tokenAdminRegistry PDA
					routerConfigPDA,
					tokenAdminRegistryPDA, // this gets created
					tokenPubKey,
					ccipAdmin,
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
					ccipAdmin,
					solana.SystemProgramID,
				).ValidateAndBuild()
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
				}
			}
		case ViaOwnerInstruction:
			// only works if the token mint authority is the deployer key
			if registerTokenConfig.Override {
				instruction, err = solRouter.NewOwnerOverridePendingAdministratorInstruction(
					tokenAdminRegistryAdmin, // admin of the tokenAdminRegistry PDA
					routerConfigPDA,
					tokenAdminRegistryPDA, // this gets created
					tokenPubKey,
					deployerKey,
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
					deployerKey, // (token mint authority) becomes the authority of the tokenAdminRegistry PDA
					solana.SystemProgramID,
				).ValidateAndBuild()
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
				}
			}
		}

		// as ccip admin is proposing the admin role, it needs to sign the transaction
		// if the ccip admin is timelock, build mcms transaction
		// else just confirm it
		if registerTokenConfig.RegisterType == ViaGetCcipAdminInstruction && ccipAdmin.Equals(timelockSignerPDA) {
			tx, err := BuildMCMSTxn(instruction, routerProgramAddress.String(), shared.Router)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
			}
			mcmsTxs = append(mcmsTxs, *tx)
		} else {
			// if we want to have a different authority, we will need to add the corresponding signer here
			// the ccip admin will always be deployer key if done without mcms
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

	deployerKey := chain.DeployerKey.PublicKey()
	timelockSignerPDA, err := FetchTimelockSigner(e, cfg.ChainSelector)
	if err != nil {
		return fmt.Errorf("failed to fetch timelock signer: %w", err)
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
		// if current admin is not the deployer key or timelock signer, we cannot transfer the admin role
		if !currentAdmin.Equals(deployerKey) && !currentAdmin.Equals(timelockSignerPDA) {
			return fmt.Errorf("current registry admin public key (%s) is not the deployer key (%s) or timelock signer (%s) for token %s, hence this changeset cannot sign for the transfer",
				currentAdmin.String(),
				deployerKey.String(),
				timelockSignerPDA.String(),
				tokenPubKey.String(),
			)
		}
		if currentAdmin.Equals(timelockSignerPDA) && cfg.MCMS == nil {
			return fmt.Errorf("current registry admin public key (%s) is the timelock signer (%s) for token %s, but no mcms config is provided, hence this changeset cannot sign for the transfer",
				currentAdmin.String(),
				timelockSignerPDA.String(),
				tokenPubKey.String(),
			)
		}
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

	timelockSignerPDA, err := FetchTimelockSigner(e, cfg.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
	}

	mcmsTxs := []mcmsTypes.Transaction{}

	for _, transferTokenAdminConfig := range cfg.TransferTokenAdminConfigs {
		tokenPubKey := transferTokenAdminConfig.TokenPubKey
		tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
		newRegistryAdminPubKey := transferTokenAdminConfig.NewRegistryAdminPublicKey

		var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
		if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get token admin registry account: %w", err)
		}
		currentAdmin := tokenAdminRegistryAccount.Administrator

		instruction, err := solRouter.NewTransferAdminRoleTokenAdminRegistryInstruction(
			newRegistryAdminPubKey,
			routerConfigPDA,
			tokenAdminRegistryPDA,
			tokenPubKey,
			currentAdmin,
		).ValidateAndBuild()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
		// when transferring admin role away from timelock, we will need to build mcms transaction
		if currentAdmin.Equals(timelockSignerPDA) {
			tx, err := BuildMCMSTxn(instruction, routerProgramAddress.String(), shared.Router)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
			}
			mcmsTxs = append(mcmsTxs, *tx)
		} else { // already confirmed that admin is either deployer key or timelock signer
			instructions := []solana.Instruction{instruction}
			if err := chain.Confirm(instructions); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}
	}

	// when transferring admin role away from timelock, we will need to build mcms transaction
	if len(mcmsTxs) > 0 {
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
	timelockSignerPDA, err := FetchTimelockSigner(e, cfg.ChainSelector)
	if err != nil {
		return fmt.Errorf("failed to fetch timelock signer: %w", err)
	}
	deployerKey := chain.DeployerKey.PublicKey()
	for _, acceptAdminRoleTokenConfig := range cfg.AcceptAdminRoleTokenConfigs {
		tokenPubKey := acceptAdminRoleTokenConfig.TokenPubKey
		if err := chainState.CommonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
			return err
		}
		if acceptAdminRoleTokenConfig.SkipRegistryCheck {
			continue
		}
		tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
		if err != nil {
			return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), routerProgramAddress.String(), err)
		}
		var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
		if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err != nil {
			return fmt.Errorf("token admin registry not found for (mint: %s, router: %s), cannot accept admin role", tokenPubKey.String(), routerProgramAddress.String())
		}
		// if pending admin is not the deployer key or timelock signer, we cannot accept the admin role
		if !tokenAdminRegistryAccount.PendingAdministrator.Equals(deployerKey) && !tokenAdminRegistryAccount.PendingAdministrator.Equals(timelockSignerPDA) {
			return fmt.Errorf("pending registry admin role is not the deployer key (%s) or timelock signer (%s) for token %s, pending admin is %s",
				deployerKey.String(),
				timelockSignerPDA.String(),
				tokenPubKey.String(),
				tokenAdminRegistryAccount.PendingAdministrator.String(),
			)
		}
		if tokenAdminRegistryAccount.PendingAdministrator.Equals(timelockSignerPDA) && cfg.MCMS == nil {
			return errors.New("pending registry admin role is the timelock signer, but no mcms config is provided, hence this changeset cannot sign for the acceptance")
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

	timelockSignerPDA, err := FetchTimelockSigner(e, cfg.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
	}
	// verified
	routerProgramAddress, routerConfigPDA, _ := chainState.GetRouterInfo()
	solRouter.SetProgramID(routerProgramAddress)
	mcmsTxs := []mcmsTypes.Transaction{}
	for _, acceptAdminRoleTokenConfig := range cfg.AcceptAdminRoleTokenConfigs {
		tokenPubKey := acceptAdminRoleTokenConfig.TokenPubKey
		tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
		var pendingAdmin solana.PublicKey
		// if skip registry check is true, then we are registering and accepting in the same batch, so while generating the instruction, we will use the timelock signer as the pending admin
		if acceptAdminRoleTokenConfig.SkipRegistryCheck {
			pendingAdmin = timelockSignerPDA
		} else {
			var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
			if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("token admin registry not found for (mint: %s, router: %s), cannot accept admin role", tokenPubKey.String(), routerProgramAddress.String())
			}
			pendingAdmin = tokenAdminRegistryAccount.PendingAdministrator
		}

		instruction, err := solRouter.NewAcceptAdminRoleTokenAdminRegistryInstruction(
			routerConfigPDA,
			tokenAdminRegistryPDA,
			tokenPubKey,
			pendingAdmin,
		).ValidateAndBuild()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
		if pendingAdmin.Equals(timelockSignerPDA) {
			tx, err := BuildMCMSTxn(instruction, routerProgramAddress.String(), shared.Router)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
			}
			mcmsTxs = append(mcmsTxs, *tx)
		} else { // already confirmed that pending admin is either deployer key or timelock signer
			instructions := []solana.Instruction{instruction}
			if err := chain.Confirm(instructions); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}
	}

	if len(mcmsTxs) > 0 {
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
	PoolType          cldf.ContractType
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
	deployerKey := chain.DeployerKey.PublicKey()
	timelockSignerPDA, err := FetchTimelockSigner(e, cfg.ChainSelector)
	if err != nil {
		return fmt.Errorf("failed to fetch timelock signer: %w", err)
	}
	for _, tokenConfig := range cfg.SetPoolTokenConfigs {
		tokenPubKey := tokenConfig.TokenPubKey
		if err := chainState.CommonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
			return err
		}
		if tokenConfig.PoolType == "" {
			return errors.New("pool type must be defined")
		}

		if tokenConfig.Metadata == "" {
			return errors.New("metadata must be defined")
		}
		if lut, ok := chainState.TokenPoolLookupTable[tokenPubKey][tokenConfig.PoolType][tokenConfig.Metadata]; !ok || lut.IsZero() {
			return fmt.Errorf("token pool lookup table not found for (mint: %s)", tokenPubKey.String())
		}
		if tokenConfig.SkipRegistryCheck {
			continue
		}
		tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
		if err != nil {
			return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), routerProgramAddress.String(), err)
		}
		var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
		if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err != nil {
			return fmt.Errorf("token admin registry not found for (mint: %s, router: %s), cannot set pool", tokenPubKey.String(), routerProgramAddress.String())
		}
		if !tokenAdminRegistryAccount.Administrator.Equals(deployerKey) && !tokenAdminRegistryAccount.Administrator.Equals(timelockSignerPDA) {
			return fmt.Errorf("token admin registry admin public key (%s) is not the deployer key (%s) or timelock signer (%s) for token %s, cannot set pool",
				tokenAdminRegistryAccount.Administrator.String(),
				deployerKey.String(),
				timelockSignerPDA.String(),
				tokenPubKey.String(),
			)
		}
		if tokenAdminRegistryAccount.Administrator.Equals(timelockSignerPDA) && cfg.MCMS == nil {
			return errors.New("registry admin role is the timelock signer, but no mcms config is provided, hence this changeset cannot sign for the set pool")
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
	timelockSignerPDA, err := FetchTimelockSigner(e, cfg.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
	}

	mcmsTxs := []mcmsTypes.Transaction{}
	for _, tokenConfig := range cfg.SetPoolTokenConfigs {
		tokenPubKey := tokenConfig.TokenPubKey
		tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
		lookupTablePubKey := chainState.TokenPoolLookupTable[tokenPubKey][tokenConfig.PoolType][tokenConfig.Metadata]

		var currentAdmin solana.PublicKey
		// if skip registry check is true, then we are registering and setting pool in the same batch, so while generating the instruction, we will use the timelock signer as the current admin
		if tokenConfig.SkipRegistryCheck {
			currentAdmin = timelockSignerPDA
		} else {
			var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
			if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("token admin registry not found for (mint: %s, router: %s), cannot set pool", tokenPubKey.String(), routerProgramAddress.String())
			}
			currentAdmin = tokenAdminRegistryAccount.Administrator
		}
		base := solRouter.NewSetPoolInstruction(
			cfg.WritableIndexes,
			routerConfigPDA,
			tokenAdminRegistryPDA,
			tokenPubKey,
			lookupTablePubKey,
			currentAdmin,
		)
		base.AccountMetaSlice = append(base.AccountMetaSlice, solana.Meta(lookupTablePubKey))
		instruction, err := base.ValidateAndBuild()
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		if currentAdmin.Equals(timelockSignerPDA) {
			tx, err := BuildMCMSTxn(instruction, routerProgramAddress.String(), shared.Router)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
			}
			mcmsTxs = append(mcmsTxs, *tx)
		} else { // already confirmed that admin is either deployer key or timelock signer
			if err = chain.Confirm([]solana.Instruction{instruction}); err != nil {
				return cldf.ChangesetOutput{}, err
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
