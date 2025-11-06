package solana

import (
	"context"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	lockrelease "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_0/lockrelease_token_pool"
	solBurnMintTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/burnmint_token_pool"
	solCommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_common"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_router"
	solLockReleaseTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/lockrelease_token_pool"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	cldfsolana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var _ cldf.ChangeSet[OnboardTokenPoolsForSelfServeConfig] = OnboardTokenPoolsForSelfServe

type OnboardTokenPoolConfig struct {
	TokenMint        solana.PublicKey
	TokenProgramName cldf.ContractType
	ProposedOwner    solana.PublicKey
	PoolType         cldf.ContractType
	Metadata         string
	Override         bool
}

type OnboardTokenPoolsForSelfServeConfig struct {
	ChainSelector        uint64
	RegisterTokenConfigs []OnboardTokenPoolConfig
	MCMS                 *proposalutils.TimelockConfig
}

func (cfg OnboardTokenPoolsForSelfServeConfig) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateRouterConfig(chain); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.Router: true}); err != nil {
		return err
	}
	routerProgramAddress, _, _ := chainState.GetRouterInfo()
	// Duplicate mint detection
	seen := make(map[string]int, len(cfg.RegisterTokenConfigs))
	for i, registerTokenConfig := range cfg.RegisterTokenConfigs {
		if registerTokenConfig.Metadata == "" {
			return fmt.Errorf("RegisterTokenConfigs[%d].Metadata is required for token mint %s", i, registerTokenConfig.TokenMint.String())
		}
		if registerTokenConfig.PoolType != shared.BurnMintTokenPool && registerTokenConfig.PoolType != shared.LockReleaseTokenPool {
			return fmt.Errorf("PoolType not supported: %v", registerTokenConfig.PoolType)
		}
		tokenMint := registerTokenConfig.TokenMint
		mintStr := tokenMint.String()
		if mintStr == "" {
			return fmt.Errorf("TokenMint cannot be empty: %v", registerTokenConfig.TokenMint)
		}
		if firstIdx, dup := seen[mintStr]; dup {
			return fmt.Errorf("duplicate token mint %s found at indexes %d and %d", mintStr, firstIdx, i)
		}
		seen[mintStr] = i
		_, err := GetTokenProgramID(registerTokenConfig.TokenProgramName)
		if err != nil {
			return fmt.Errorf("TokenProgramName not found in registerTokenConfig: %v", registerTokenConfig.TokenProgramName)
		}
		if registerTokenConfig.ProposedOwner.IsZero() {
			return errors.New("token admin registry admin is required")
		}
		tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenMint, routerProgramAddress)
		if err != nil {
			return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w",
				mintStr, routerProgramAddress.String(), err)
		}
		var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
		if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err == nil {
			if !registerTokenConfig.Override {
				return fmt.Errorf("token admin registry already exists for (mint: %s, router: %s)", mintStr, routerProgramAddress.String())
			}
		}
		tokenPoolProgramID := chainState.GetActiveTokenPool(registerTokenConfig.PoolType, shared.CLLMetadata) // This changeset is to register the token pool in the CLL Token Pool Program
		if (tokenPoolProgramID == solana.PublicKey{}) {
			return fmt.Errorf("token pool program ID not found for pool type: %s", registerTokenConfig.PoolType)
		}
		tokenPoolPDA, err := solTokenUtil.TokenPoolConfigAddress(tokenMint, tokenPoolProgramID)
		if err != nil {
			return err
		}
		var tokenPoolAccount lockrelease.State
		if err := chain.GetAccountDataBorshInto(context.Background(), tokenPoolPDA, &tokenPoolAccount); err == nil {
			if !registerTokenConfig.Override {
				return fmt.Errorf("token pool already initialized for (mint: %s, program: %s)", mintStr, tokenPoolProgramID.String())
			}
		}
	}
	return nil
}

// OnboardTokenPoolsForSelfServe registers a token admin registry for a given token and initializes the token pool in CLL Token Pool Program.
// This changeset is used when the owner of the token pool doesn't have the mint authority over the token, but they want to self serve.
// So, this changeset includes the minimum configuration that CCIP Admin needs to do in the Token Admin Registry and in the Token Pool Program
func OnboardTokenPoolsForSelfServe(e cldf.Environment, cfg OnboardTokenPoolsForSelfServeConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("OnboardTokenPoolsForSelfServe", "cfg", cfg)
	solChainState, routerState, err := loadRouterSolanaState(e, cfg)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	mcmsTxs := []mcmsTypes.Transaction{}
	instructions := [][]solana.Instruction{}
	for _, registerTokenConfig := range cfg.RegisterTokenConfigs {
		// Propose Admin in Token Admin Registry
		proposeTokenAdminRegistryAdminIx, err := generateProposeTokenAdminRegistryAdministratorIx(registerTokenConfig, routerState)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		currentTokenPoolSolanaState, err := loadTokenPoolSolanaState(registerTokenConfig, solChainState)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		tokenInstructions := []solana.Instruction{proposeTokenAdminRegistryAdminIx}
		var initializeTokenPoolIx solana.Instruction
		if !registerTokenConfig.Override {
			// Initialize Token Pool in CLL Program
			initializeTokenPoolIx, err = generateInitializeCLLTokenPoolIx(registerTokenConfig, currentTokenPoolSolanaState)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			tokenInstructions = append(tokenInstructions, initializeTokenPoolIx)
		}
		// Propose new owner of the token pool
		transferTokenPoolOwnershipIx, err := generateTransferTokenPoolOwnershipIx(registerTokenConfig, currentTokenPoolSolanaState)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		tokenInstructions = append(tokenInstructions, transferTokenPoolOwnershipIx)
		e.Logger.Infow("Onboarding Token in ", "TokenProgramID", currentTokenPoolSolanaState.tokenPoolProgramID.String())
		// if the ccip admin is timelock, build mcms transaction
		if cfg.MCMS != nil {
			inputs := []MCMSTxParams{{
				Ix:           proposeTokenAdminRegistryAdminIx,
				ProgramID:    routerState.routerProgramID.String(),
				ContractType: shared.Router}}
			if !registerTokenConfig.Override {
				inputs = append(inputs,
					MCMSTxParams{
						Ix:           initializeTokenPoolIx,
						ProgramID:    currentTokenPoolSolanaState.tokenPoolProgramID.String(),
						ContractType: registerTokenConfig.PoolType})
			}
			inputs = append(inputs,
				MCMSTxParams{
					Ix:           transferTokenPoolOwnershipIx,
					ProgramID:    currentTokenPoolSolanaState.tokenPoolProgramID.String(),
					ContractType: registerTokenConfig.PoolType})
			moreTx, err := BuildManyMCMSTxsFrom(inputs)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			for _, tx := range moreTx {
				mcmsTxs = append(mcmsTxs, *tx)
			}
		} else {
			// the ccip admin will always be deployer key if done without mcms
			instructions = append(instructions, tokenInstructions)
		}
		if !registerTokenConfig.Override {
			// Store in Address Book only first time running this
			newAddresses := cldf.NewMemoryAddressBook()
			tv := cldf.NewTypeAndVersion(registerTokenConfig.TokenProgramName, deployment.Version1_0_0)
			tv.AddLabel(registerTokenConfig.Metadata)                            // Customer Identifier
			tv.AddLabel(registerTokenConfig.PoolType.String())                   // Pool Type
			tv.AddLabel(currentTokenPoolSolanaState.tokenPoolProgramID.String()) // Token Pool Program ID
			err = newAddresses.Save(cfg.ChainSelector, registerTokenConfig.TokenMint.String(), tv)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
		}
	}
	return ExecuteInstructionsAndBuildProposals(e, ExecuteConfig{ChainSelector: cfg.ChainSelector, MCMS: cfg.MCMS, Chain: solChainState.chain}, instructions, mcmsTxs)
}

func generateProposeTokenAdminRegistryAdministratorIx(registerTokenConfig OnboardTokenPoolConfig, routerState routerSolanaState) (solana.Instruction, error) {
	tokenPubKey := registerTokenConfig.TokenMint
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerState.routerProgramID)
	tokenAdminRegistryAdmin := registerTokenConfig.ProposedOwner
	var instruction solana.Instruction
	// the ccip admin signs and makes tokenAdminRegistryAdmin the pending authority of the tokenAdminRegistry PDA, then they need to accept the role
	if !registerTokenConfig.Override {
		tempIx, err := solRouter.NewCcipAdminProposeAdministratorInstruction(
			tokenAdminRegistryAdmin, // customer's admin of the tokenAdminRegistry PDA in the Router
			routerState.routerConfigPDA,
			tokenAdminRegistryPDA, // If invoking the first time, this PDA is created
			tokenPubKey,
			routerState.ccipAdmin,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("failed to generate instruction to propose administrator: %w", err)
		}
		ixData, err := tempIx.Data()
		if err != nil {
			return nil, fmt.Errorf("failed to extract data payload from ccip admin propose admin instruction: %w", err)
		}
		instruction = solana.NewInstruction(routerState.routerProgramID, tempIx.Accounts(), ixData)
	} else {
		// Use this if the proposed token admin registry admin set was incorrect
		overridePendingAdministratorIx, err := solRouter.NewCcipAdminOverridePendingAdministratorInstruction(
			tokenAdminRegistryAdmin, // customer's admin of the tokenAdminRegistry PDA in the Router
			routerState.routerConfigPDA,
			tokenAdminRegistryPDA,
			tokenPubKey,
			routerState.ccipAdmin,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("failed to generate instruction to override pending administrator: %w", err)
		}
		ixData, err := overridePendingAdministratorIx.Data()
		if err != nil {
			return nil, fmt.Errorf("failed to extract data payload from ccip admin override pending admin instruction: %w", err)
		}
		instruction = solana.NewInstruction(routerState.routerProgramID, overridePendingAdministratorIx.Accounts(), ixData)
	}
	return instruction, nil
}

func generateInitializeCLLTokenPoolIx(config OnboardTokenPoolConfig, state tokenPoolSolanaState) (solana.Instruction, error) {
	switch config.PoolType {
	case shared.BurnMintTokenPool:
		solBurnMintTokenPool.SetProgramID(state.tokenPoolProgramID)
		return solBurnMintTokenPool.NewInitializeInstruction(
			state.poolConfigPDA,
			config.TokenMint,
			state.upgradeAuthority,
			solana.SystemProgramID,
			state.tokenPoolProgramID,
			state.programDataAddress,
			state.configPDA,
		).ValidateAndBuild()
	case shared.LockReleaseTokenPool:
		solLockReleaseTokenPool.SetProgramID(state.tokenPoolProgramID)
		return solLockReleaseTokenPool.NewInitializeInstruction(
			state.poolConfigPDA,
			config.TokenMint,
			state.upgradeAuthority,
			solana.SystemProgramID,
			state.tokenPoolProgramID,
			state.programDataAddress,
			state.configPDA,
		).ValidateAndBuild()
	default:
		return nil, errors.New("invalid token pool type")
	}
}

func generateTransferTokenPoolOwnershipIx(config OnboardTokenPoolConfig, state tokenPoolSolanaState) (solana.Instruction, error) {
	switch config.PoolType {
	case shared.BurnMintTokenPool:
		solBurnMintTokenPool.SetProgramID(state.tokenPoolProgramID)
		return solBurnMintTokenPool.NewTransferOwnershipInstruction(
			config.ProposedOwner,
			state.poolConfigPDA,
			config.TokenMint,
			state.upgradeAuthority,
		).ValidateAndBuild()
	case shared.LockReleaseTokenPool:
		solLockReleaseTokenPool.SetProgramID(state.tokenPoolProgramID)
		return solLockReleaseTokenPool.NewTransferOwnershipInstruction(
			config.ProposedOwner,
			state.poolConfigPDA,
			config.TokenMint,
			state.upgradeAuthority,
		).ValidateAndBuild()
	default:
		return nil, errors.New("invalid token pool type")
	}
}

type globalState struct {
	chain      cldfsolana.Chain
	chainState solanastateview.CCIPChainState
}

type routerSolanaState struct {
	routerProgramID solana.PublicKey
	routerConfigPDA solana.PublicKey
	ccipAdmin       solana.PublicKey
}

func loadRouterSolanaState(e cldf.Environment, cfg OnboardTokenPoolsForSelfServeConfig) (globalState, routerSolanaState, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return globalState{}, routerSolanaState{}, err
	}
	chainState, ok := state.SolChains[cfg.ChainSelector]
	if !ok {
		return globalState{}, routerSolanaState{}, fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}

	if err := cfg.Validate(e, chainState); err != nil {
		return globalState{}, routerSolanaState{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	routerProgramAddress, routerConfigPDA, _ := chainState.GetRouterInfo()
	ccipAdmin := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		shared.Router,
		solana.PublicKey{},
		"",
	)
	return globalState{
			chain:      chain,
			chainState: chainState,
		}, routerSolanaState{
			routerProgramID: routerProgramAddress,
			routerConfigPDA: routerConfigPDA,
			ccipAdmin:       ccipAdmin,
		}, nil
}

type tokenPoolSolanaState struct {
	tokenPoolProgramID solana.PublicKey
	poolConfigPDA      solana.PublicKey
	configPDA          solana.PublicKey
	programDataAddress solana.PublicKey
	upgradeAuthority   solana.PublicKey
}

func loadTokenPoolSolanaState(cfg OnboardTokenPoolConfig, state globalState) (tokenPoolSolanaState, error) {
	tokenPoolProgramID := state.chainState.GetActiveTokenPool(cfg.PoolType, shared.CLLMetadata) // This changeset is to set up the token pool in the CLL Program
	if (tokenPoolProgramID == solana.PublicKey{}) {
		return tokenPoolSolanaState{}, fmt.Errorf("token pool program ID not found for pool type: %s", cfg.PoolType)
	}
	poolConfigPDA, err := solTokenUtil.TokenPoolConfigAddress(cfg.TokenMint, tokenPoolProgramID)
	if err != nil {
		return tokenPoolSolanaState{}, err
	}
	configPDA, err := TokenPoolGlobalConfigPDA(tokenPoolProgramID)
	if err != nil {
		return tokenPoolSolanaState{}, fmt.Errorf("failed to get solana token pool global config PDA: %w", err)
	}
	progDataAddr, err := deployment.GetProgramDataAddress(state.chain.Client, tokenPoolProgramID)
	if err != nil {
		return tokenPoolSolanaState{}, fmt.Errorf("failed to get program data address for program %s: %w", tokenPoolProgramID.String(), err)
	}
	upgradeAuthority, _, err := deployment.GetUpgradeAuthority(state.chain.Client, progDataAddr)
	if err != nil {
		return tokenPoolSolanaState{}, fmt.Errorf("failed to get upgrade authority for program data %s: %w", progDataAddr.String(), err)
	}
	return tokenPoolSolanaState{
		tokenPoolProgramID: tokenPoolProgramID,
		poolConfigPDA:      poolConfigPDA,
		configPDA:          configPDA,
		programDataAddress: progDataAddr,
		upgradeAuthority:   upgradeAuthority,
	}, nil
}
