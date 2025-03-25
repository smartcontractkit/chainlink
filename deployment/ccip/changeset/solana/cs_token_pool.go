package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solBaseTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/base_token_pool"
	solBurnMintTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/burnmint_token_pool"
	solCommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_common"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solLockReleaseTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/lockrelease_token_pool"
	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

var _ deployment.ChangeSet[TokenPoolConfig] = AddTokenPool
var _ deployment.ChangeSet[RemoteChainTokenPoolConfig] = SetupTokenPoolForRemoteChain

func validatePoolDeployment(
	e *deployment.Environment,
	poolType solTestTokenPool.PoolType,
	selector uint64,
	tokenPubKey solana.PublicKey,
	validatePoolConfig bool) error {
	state, _ := ccipChangeset.LoadOnchainState(*e)
	chainState := state.SolChains[selector]
	chain := e.SolChains[selector]

	var tokenPool solana.PublicKey
	var poolConfigAccount interface{}

	if _, err := chainState.TokenToTokenProgram(tokenPubKey); err != nil {
		return fmt.Errorf("failed to get token program for token address %s: %w", tokenPubKey.String(), err)
	}
	switch poolType {
	case solTestTokenPool.BurnAndMint_PoolType:
		if chainState.BurnMintTokenPool.IsZero() {
			return fmt.Errorf("token pool of type BurnAndMint not found in existing state, deploy the token pool first for chain %d", selector)
		}
		tokenPool = chainState.BurnMintTokenPool
		poolConfigAccount = solBurnMintTokenPool.State{}
	case solTestTokenPool.LockAndRelease_PoolType:
		if chainState.LockReleaseTokenPool.IsZero() {
			return fmt.Errorf("token pool of type LockAndRelease not found in existing state, deploy the token pool first for chain %d", selector)
		}
		tokenPool = chainState.LockReleaseTokenPool
		poolConfigAccount = solLockReleaseTokenPool.State{}
	default:
		return fmt.Errorf("invalid pool type: %s", poolType)
	}

	if validatePoolConfig {
		poolConfigPDA, err := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
		if err != nil {
			return fmt.Errorf("failed to get token pool config address (mint: %s, pool: %s): %w", tokenPubKey.String(), tokenPool.String(), err)
		}
		if err := chain.GetAccountDataBorshInto(context.Background(), poolConfigPDA, &poolConfigAccount); err != nil {
			return fmt.Errorf("token pool config not found (mint: %s, pool: %s, type: %s): %w", tokenPubKey.String(), tokenPool.String(), poolType, err)
		}
	}
	return nil
}

type TokenPoolConfig struct {
	ChainSelector uint64
	PoolType      solTestTokenPool.PoolType
	TokenPubKey   string
}

func (cfg TokenPoolConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}

	return validatePoolDeployment(&e, cfg.PoolType, cfg.ChainSelector, tokenPubKey, false)
}

func AddTokenPool(e deployment.Environment, cfg TokenPoolConfig) (deployment.ChangesetOutput, error) {
	e.Logger.Infow("Adding token pool", "token_pubkey", cfg.TokenPubKey)
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	tokenPool := solana.PublicKey{}

	if cfg.PoolType == solTestTokenPool.BurnAndMint_PoolType {
		tokenPool = chainState.BurnMintTokenPool
		solBurnMintTokenPool.SetProgramID(tokenPool)
	} else if cfg.PoolType == solTestTokenPool.LockAndRelease_PoolType {
		tokenPool = chainState.LockReleaseTokenPool
		solLockReleaseTokenPool.SetProgramID(tokenPool)
	}

	// verified
	tokenprogramID, _ := chainState.TokenToTokenProgram(tokenPubKey)
	poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
	poolSigner, _ := solTokenUtil.TokenPoolSignerAddress(tokenPubKey, tokenPool)
	routerProgramAddress, _, _ := chainState.GetRouterInfo()
	rmnRemoteAddress := chainState.RMNRemote
	// ata for token pool
	createI, tokenPoolATA, err := solTokenUtil.CreateAssociatedTokenAccount(
		tokenprogramID,
		tokenPubKey,
		poolSigner,
		chain.DeployerKey.PublicKey(),
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create associated token account for tokenpool (mint: %s, pool: %s): %w", tokenPubKey.String(), tokenPool.String(), err)
	}
	instructions := []solana.Instruction{createI}

	var poolInitI solana.Instruction
	switch cfg.PoolType {
	case solTestTokenPool.BurnAndMint_PoolType:
		// initialize token pool for token
		poolInitI, err = solBurnMintTokenPool.NewInitializeInstruction(
			routerProgramAddress,
			rmnRemoteAddress,
			poolConfigPDA,
			tokenPubKey,
			chain.DeployerKey.PublicKey(), // a token pool will only ever be added by the deployer key.
			solana.SystemProgramID,
		).ValidateAndBuild()
	case solTestTokenPool.LockAndRelease_PoolType:
		// initialize token pool for token
		poolInitI, err = solLockReleaseTokenPool.NewInitializeInstruction(
			routerProgramAddress,
			rmnRemoteAddress,
			poolConfigPDA,
			tokenPubKey,
			chain.DeployerKey.PublicKey(), // a token pool will only ever be added by the deployer key.
			solana.SystemProgramID,
		).ValidateAndBuild()
	default:
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid pool type: %s", cfg.PoolType)
	}
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}

	instructions = append(instructions, poolInitI)

	if cfg.PoolType == solTestTokenPool.BurnAndMint_PoolType && tokenPubKey != solana.SolMint {
		// make pool mint_authority for token
		authI, err := solTokenUtil.SetTokenMintAuthority(
			tokenprogramID,
			poolSigner,
			tokenPubKey,
			chain.DeployerKey.PublicKey(),
		)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
		instructions = append(instructions, authI)
	}

	// add signer here if authority is different from deployer key
	if err := chain.Confirm(instructions); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	e.Logger.Infow("Created new token pool config", "token_pool_ata", tokenPoolATA.String(), "pool_config", poolConfigPDA.String(), "pool_signer", poolSigner.String())
	e.Logger.Infow("Set mint authority", "poolSigner", poolSigner.String())

	return deployment.ChangesetOutput{}, nil
}

// ADD TOKEN POOL FOR REMOTE CHAIN
type RemoteChainTokenPoolConfig struct {
	SolChainSelector    uint64
	RemoteChainSelector uint64
	SolTokenPubKey      string
	PoolType            solTestTokenPool.PoolType
	// this is actually derivable from on chain given token symbol
	RemoteConfig      solBaseTokenPool.RemoteConfig
	InboundRateLimit  solBaseTokenPool.RateLimitConfig
	OutboundRateLimit solBaseTokenPool.RateLimitConfig
	MCMSSolana        *MCMSConfigSolana
	IsUpdate          bool
}

func (cfg RemoteChainTokenPoolConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)
	if err := commonValidation(e, cfg.SolChainSelector, tokenPubKey); err != nil {
		return err
	}
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.SolChainSelector]
	chain := e.SolChains[cfg.SolChainSelector]

	if err := validatePoolDeployment(&e, cfg.PoolType, cfg.SolChainSelector, tokenPubKey, true); err != nil {
		return err
	}

	if err := ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, tokenPubKey); err != nil {
		return err
	}

	var tokenPool solana.PublicKey
	var remoteChainConfigAccount interface{}

	switch cfg.PoolType {
	case solTestTokenPool.BurnAndMint_PoolType:
		tokenPool = chainState.BurnMintTokenPool
		remoteChainConfigAccount = solBurnMintTokenPool.ChainConfig{}
	case solTestTokenPool.LockAndRelease_PoolType:
		tokenPool = chainState.LockReleaseTokenPool
		remoteChainConfigAccount = solLockReleaseTokenPool.ChainConfig{}
	default:
		return fmt.Errorf("invalid pool type: %s", cfg.PoolType)
	}

	// check if this remote chain is already configured for this token
	remoteChainConfigPDA, _, err := solTokenUtil.TokenPoolChainConfigPDA(cfg.RemoteChainSelector, tokenPubKey, tokenPool)
	if err != nil {
		return fmt.Errorf("failed to get token pool remote chain config pda (remoteSelector: %d, mint: %s, pool: %s): %w", cfg.RemoteChainSelector, tokenPubKey.String(), tokenPool.String(), err)
	}
	err = chain.GetAccountDataBorshInto(context.Background(), remoteChainConfigPDA, &remoteChainConfigAccount)

	if !cfg.IsUpdate && err == nil {
		return fmt.Errorf("remote chain config already exists for (remoteSelector: %d, mint: %s, pool: %s, type: %s)", cfg.RemoteChainSelector, tokenPubKey.String(), tokenPool.String(), cfg.PoolType)
	} else if cfg.IsUpdate && err != nil {
		return fmt.Errorf("remote chain config not found for (remoteSelector: %d, mint: %s, pool: %s, type: %s): %w", cfg.RemoteChainSelector, tokenPubKey.String(), tokenPool.String(), cfg.PoolType, err)
	}
	return nil
}

func SetupTokenPoolForRemoteChain(e deployment.Environment, cfg RemoteChainTokenPoolConfig) (deployment.ChangesetOutput, error) {
	e.Logger.Infow("Setting up token pool for remote chain", "remote_chain_selector", cfg.RemoteChainSelector, "token_pubkey", cfg.SolTokenPubKey)
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.SolChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.SolChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)

	var instructions []solana.Instruction
	var txns []mcmsTypes.Transaction
	var err error
	switch cfg.PoolType {
	case solTestTokenPool.BurnAndMint_PoolType:
		instructions, err = getInstructionsForBurnMint(e, chain, chainState, cfg)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
		if cfg.MCMSSolana != nil && cfg.MCMSSolana.BurnMintTokenPoolOwnedByTimelock[tokenPubKey] {
			for _, ixn := range instructions {
				tx, err := BuildMCMSTxn(ixn, chainState.BurnMintTokenPool.String(), ccipChangeset.BurnMintToken)
				if err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate mcms txn: %w", err)
				}
				txns = append(txns, *tx)
			}
		}
	case solTestTokenPool.LockAndRelease_PoolType:
		instructions, err = getInstructionsForLockRelease(e, chain, chainState, cfg)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
		if cfg.MCMSSolana != nil && cfg.MCMSSolana.LockReleaseTokenPoolOwnedByTimelock[tokenPubKey] {
			for _, ixn := range instructions {
				tx, err := BuildMCMSTxn(ixn, chainState.LockReleaseTokenPool.String(), ccipChangeset.LockReleaseTokenPool)
				if err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate mcms txn: %w", err)
				}
				txns = append(txns, *tx)
			}
		}
	default:
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid pool type: %s", cfg.PoolType)
	}

	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.SolChainSelector, "proposal to edit token pools in Solana", cfg.MCMSSolana.MCMS.MinDelay, txns)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}
	if err = chain.Confirm(instructions); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	e.Logger.Infow("Configured token pool for remote chain", "remote_chain_selector", cfg.RemoteChainSelector, "token_pubkey", tokenPubKey.String())
	return deployment.ChangesetOutput{}, nil
}

func getInstructionsForBurnMint(
	e deployment.Environment,
	chain deployment.SolChain,
	chainState ccipChangeset.SolCCIPChainState,
	cfg RemoteChainTokenPoolConfig,
) ([]solana.Instruction, error) {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)
	poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, chainState.BurnMintTokenPool)
	remoteChainConfigPDA, _, _ := solTokenUtil.TokenPoolChainConfigPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.BurnMintTokenPool)
	solBurnMintTokenPool.SetProgramID(chainState.BurnMintTokenPool)
	ixns := make([]solana.Instruction, 0)
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.BurnMintTokenPool,
		tokenPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
	if !cfg.IsUpdate {
		ixConfigure, err := solBurnMintTokenPool.NewInitChainRemoteConfigInstruction(
			cfg.RemoteChainSelector,
			tokenPubKey,
			cfg.RemoteConfig,
			poolConfigPDA,
			remoteChainConfigPDA,
			authority,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("failed to generate instructions: %w", err)
		}
		ixns = append(ixns, ixConfigure)
	} else {
		ixConfigure, err := solBurnMintTokenPool.NewEditChainRemoteConfigInstruction(
			cfg.RemoteChainSelector,
			tokenPubKey,
			cfg.RemoteConfig,
			poolConfigPDA,
			remoteChainConfigPDA,
			authority,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("failed to generate instructions: %w", err)
		}
		ixns = append(ixns, ixConfigure)
	}
	ixRates, err := solBurnMintTokenPool.NewSetChainRateLimitInstruction(
		cfg.RemoteChainSelector,
		tokenPubKey,
		cfg.InboundRateLimit,
		cfg.OutboundRateLimit,
		poolConfigPDA,
		remoteChainConfigPDA,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("failed to generate instructions: %w", err)
	}
	ixns = append(ixns, ixRates)
	if len(cfg.RemoteConfig.PoolAddresses) > 0 {
		ixAppend, err := solBurnMintTokenPool.NewAppendRemotePoolAddressesInstruction(
			cfg.RemoteChainSelector,
			tokenPubKey,
			cfg.RemoteConfig.PoolAddresses, // i dont know why this is a list (is it for different types of pool of the same token?)
			poolConfigPDA,
			remoteChainConfigPDA,
			authority,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("failed to generate instructions: %w", err)
		}
		ixns = append(ixns, ixAppend)
	}
	return ixns, nil
}

func getInstructionsForLockRelease(
	e deployment.Environment,
	chain deployment.SolChain,
	chainState ccipChangeset.SolCCIPChainState,
	cfg RemoteChainTokenPoolConfig,
) ([]solana.Instruction, error) {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)
	poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, chainState.LockReleaseTokenPool)
	remoteChainConfigPDA, _, _ := solTokenUtil.TokenPoolChainConfigPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.LockReleaseTokenPool)
	solLockReleaseTokenPool.SetProgramID(chainState.LockReleaseTokenPool)
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.LockReleaseTokenPool,
		tokenPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
	ixns := make([]solana.Instruction, 0)
	if !cfg.IsUpdate {
		ixConfigure, err := solLockReleaseTokenPool.NewInitChainRemoteConfigInstruction(
			cfg.RemoteChainSelector,
			tokenPubKey,
			cfg.RemoteConfig,
			poolConfigPDA,
			remoteChainConfigPDA,
			authority,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("failed to generate instructions: %w", err)
		}
		ixns = append(ixns, ixConfigure)
	} else {
		ixConfigure, err := solLockReleaseTokenPool.NewEditChainRemoteConfigInstruction(
			cfg.RemoteChainSelector,
			tokenPubKey,
			cfg.RemoteConfig,
			poolConfigPDA,
			remoteChainConfigPDA,
			authority,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("failed to generate instructions: %w", err)
		}
		ixns = append(ixns, ixConfigure)
	}
	ixRates, err := solLockReleaseTokenPool.NewSetChainRateLimitInstruction(
		cfg.RemoteChainSelector,
		tokenPubKey,
		cfg.InboundRateLimit,
		cfg.OutboundRateLimit,
		poolConfigPDA,
		remoteChainConfigPDA,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("failed to generate instructions: %w", err)
	}
	ixns = append(ixns, ixRates)
	if len(cfg.RemoteConfig.PoolAddresses) > 0 {
		ixAppend, err := solLockReleaseTokenPool.NewAppendRemotePoolAddressesInstruction(
			cfg.RemoteChainSelector,
			tokenPubKey,
			cfg.RemoteConfig.PoolAddresses, // i dont know why this is a list (is it for different types of pool of the same token?)
			poolConfigPDA,
			remoteChainConfigPDA,
			authority,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("failed to generate instructions: %w", err)
		}
		ixns = append(ixns, ixAppend)
	}
	return ixns, nil
}

// ADD TOKEN POOL LOOKUP TABLE
type TokenPoolLookupTableConfig struct {
	ChainSelector uint64
	TokenPubKey   string
	PoolType      solTestTokenPool.PoolType
}

func (cfg TokenPoolLookupTableConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	return validatePoolDeployment(&e, cfg.PoolType, cfg.ChainSelector, tokenPubKey, false)
}

func AddTokenPoolLookupTable(e deployment.Environment, cfg TokenPoolLookupTableConfig) (deployment.ChangesetOutput, error) {
	e.Logger.Infow("Adding token pool lookup table", "token_pubkey", cfg.TokenPubKey)
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chain := e.SolChains[cfg.ChainSelector]
	ctx := e.GetContext()
	client := chain.Client
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	authorityPrivKey := chain.DeployerKey // assuming the authority is the deployer key
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	tokenPool := solana.PublicKey{}
	if cfg.PoolType == solTestTokenPool.BurnAndMint_PoolType {
		tokenPool = chainState.BurnMintTokenPool
	} else if cfg.PoolType == solTestTokenPool.LockAndRelease_PoolType {
		tokenPool = chainState.LockReleaseTokenPool
	}
	routerProgramAddress, _, _ := chainState.GetRouterInfo()
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
	tokenPoolChainConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
	tokenPoolSigner, _ := solTokenUtil.TokenPoolSignerAddress(tokenPubKey, tokenPool)
	tokenProgram, _ := chainState.TokenToTokenProgram(tokenPubKey)
	poolTokenAccount, _, _ := solTokenUtil.FindAssociatedTokenAddress(tokenProgram, tokenPubKey, tokenPoolSigner)
	feeTokenConfigPDA, _, _ := solState.FindFqBillingTokenConfigPDA(tokenPubKey, chainState.FeeQuoter)

	// the 'table' address is not derivable
	// but this will be stored in tokenAdminRegistryPDA as a part of the SetPool changeset
	// and tokenAdminRegistryPDA is derivable using token and router address
	table, err := solCommonUtil.CreateLookupTable(ctx, client, *authorityPrivKey)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create lookup table for token pool (mint: %s): %w", tokenPubKey.String(), err)
	}
	list := solana.PublicKeySlice{
		table,                   // 0
		tokenAdminRegistryPDA,   // 1
		tokenPool,               // 2
		tokenPoolChainConfigPDA, // 3 - writable
		poolTokenAccount,        // 4 - writable
		tokenPoolSigner,         // 5
		tokenProgram,            // 6
		tokenPubKey,             // 7 - writable
		feeTokenConfigPDA,       // 8
	}
	if err = solCommonUtil.ExtendLookupTable(ctx, client, table, *authorityPrivKey, list); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to extend lookup table for token pool (mint: %s): %w", tokenPubKey.String(), err)
	}
	if err := solCommonUtil.AwaitSlotChange(ctx, client); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to await slot change while extending lookup table: %w", err)
	}
	newAddressBook := deployment.NewMemoryAddressBook()
	tv := deployment.NewTypeAndVersion(ccipChangeset.TokenPoolLookupTable, deployment.Version1_0_0)
	tv.Labels.Add(tokenPubKey.String())
	if err := newAddressBook.Save(cfg.ChainSelector, table.String(), tv); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to save tokenpool address lookup table: %w", err)
	}
	e.Logger.Infow("Added token pool lookup table", "token_pubkey", tokenPubKey.String())
	return deployment.ChangesetOutput{
		AddressBook: newAddressBook,
	}, nil
}

type SetPoolConfig struct {
	ChainSelector   uint64
	TokenPubKey     string
	WritableIndexes []uint8
	MCMSSolana      *MCMSConfigSolana
}

func (cfg SetPoolConfig) Validate(e deployment.Environment) error {
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
	if err := ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, tokenPubKey); err != nil {
		return err
	}
	routerProgramAddress, _, _ := chainState.GetRouterInfo()
	tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
	if err != nil {
		return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), routerProgramAddress.String(), err)
	}
	var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
	if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err != nil {
		return fmt.Errorf("token admin registry not found for (mint: %s, router: %s), cannot set pool", tokenPubKey.String(), routerProgramAddress.String())
	}
	if _, ok := chainState.TokenPoolLookupTable[tokenPubKey]; !ok {
		return fmt.Errorf("token pool lookup table not found for (mint: %s)", tokenPubKey.String())
	}
	return nil
}

// this sets the writable indexes of the token pool lookup table
func SetPool(e deployment.Environment, cfg SetPoolConfig) (deployment.ChangesetOutput, error) {
	e.Logger.Infof("Setting pool config for token %s", cfg.TokenPubKey)
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	routerProgramAddress, routerConfigPDA, _ := chainState.GetRouterInfo()
	solRouter.SetProgramID(routerProgramAddress)
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
	lookupTablePubKey := chainState.TokenPoolLookupTable[tokenPubKey]

	routerUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.Router,
		tokenPubKey)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
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
		return deployment.ChangesetOutput{}, err
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

	if err = chain.Confirm([]solana.Instruction{instruction}); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	e.Logger.Infow("Set pool config", "token_pubkey", tokenPubKey.String())
	return deployment.ChangesetOutput{}, nil
}

type ConfigureTokenPoolAllowListConfig struct {
	SolChainSelector uint64
	SolTokenPubKey   string
	PoolType         solTestTokenPool.PoolType
	Accounts         []solana.PublicKey
	// whether or not the given accounts are being added to the allow list or removed
	// i.e. true = add, false = remove
	Enabled    bool
	MCMSSolana *MCMSConfigSolana
}

func (cfg ConfigureTokenPoolAllowListConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.SolChains[cfg.SolChainSelector]
	chain := e.SolChains[cfg.SolChainSelector]
	if err := commonValidation(e, cfg.SolChainSelector, tokenPubKey); err != nil {
		return err
	}
	if err := validatePoolDeployment(&e, cfg.PoolType, cfg.SolChainSelector, tokenPubKey, true); err != nil {
		return err
	}
	return ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, tokenPubKey)
}

func ConfigureTokenPoolAllowList(e deployment.Environment, cfg ConfigureTokenPoolAllowListConfig) (deployment.ChangesetOutput, error) {
	e.Logger.Infof("Configuring token pool allowlist for token %s", cfg.SolTokenPubKey)
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.SolChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.SolChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)

	var ix solana.Instruction
	var tokenPoolUsingMcms bool
	var programID solana.PublicKey
	var contractType deployment.ContractType
	switch cfg.PoolType {
	case solTestTokenPool.BurnAndMint_PoolType:
		poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, chainState.BurnMintTokenPool)
		solBurnMintTokenPool.SetProgramID(chainState.BurnMintTokenPool)
		programID = chainState.BurnMintTokenPool
		contractType = ccipChangeset.BurnMintTokenPool
		tokenPoolUsingMcms = cfg.MCMSSolana != nil && cfg.MCMSSolana.BurnMintTokenPoolOwnedByTimelock[tokenPubKey]
		authority, err := GetAuthorityForIxn(
			&e,
			chain,
			cfg.MCMSSolana,
			ccipChangeset.BurnMintTokenPool,
			tokenPubKey)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
		}
		ix, err = solBurnMintTokenPool.NewConfigureAllowListInstruction(
			cfg.Accounts,
			cfg.Enabled,
			poolConfigPDA,
			tokenPubKey,
			authority,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
	case solTestTokenPool.LockAndRelease_PoolType:
		poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, chainState.LockReleaseTokenPool)
		solLockReleaseTokenPool.SetProgramID(chainState.LockReleaseTokenPool)
		programID = chainState.LockReleaseTokenPool
		contractType = ccipChangeset.LockReleaseTokenPool
		tokenPoolUsingMcms = cfg.MCMSSolana != nil && cfg.MCMSSolana.LockReleaseTokenPoolOwnedByTimelock[tokenPubKey]
		authority, err := GetAuthorityForIxn(
			&e,
			chain,
			cfg.MCMSSolana,
			ccipChangeset.LockReleaseTokenPool,
			tokenPubKey)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
		}
		ix, err = solLockReleaseTokenPool.NewConfigureAllowListInstruction(
			cfg.Accounts,
			cfg.Enabled,
			poolConfigPDA,
			tokenPubKey,
			authority,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
	default:
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid pool type: %s", cfg.PoolType)
	}
	if tokenPoolUsingMcms {
		tx, err := BuildMCMSTxn(ix, programID.String(), contractType)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.SolChainSelector, "proposal to ConfigureTokenPoolAllowList in Solana", cfg.MCMSSolana.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
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
	e.Logger.Infow("Configured token pool allowlist", "token_pubkey", tokenPubKey.String())
	return deployment.ChangesetOutput{}, nil
}

type RemoveFromAllowListConfig struct {
	SolChainSelector uint64
	SolTokenPubKey   string
	PoolType         solTestTokenPool.PoolType
	Accounts         []solana.PublicKey
	MCMSSolana       *MCMSConfigSolana
}

func (cfg RemoveFromAllowListConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)
	if err := commonValidation(e, cfg.SolChainSelector, tokenPubKey); err != nil {
		return err
	}
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.SolChains[cfg.SolChainSelector]
	chain := e.SolChains[cfg.SolChainSelector]
	if err := ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, tokenPubKey); err != nil {
		return err
	}
	return validatePoolDeployment(&e, cfg.PoolType, cfg.SolChainSelector, tokenPubKey, true)
}

func RemoveFromTokenPoolAllowList(e deployment.Environment, cfg RemoveFromAllowListConfig) (deployment.ChangesetOutput, error) {
	e.Logger.Infof("Removing from token pool allowlist for token %s", cfg.SolTokenPubKey)
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.SolChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.SolChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)

	var ix solana.Instruction
	var tokenPoolUsingMcms bool
	var programID solana.PublicKey
	var contractType deployment.ContractType
	switch cfg.PoolType {
	case solTestTokenPool.BurnAndMint_PoolType:
		poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, chainState.BurnMintTokenPool)
		solBurnMintTokenPool.SetProgramID(chainState.BurnMintTokenPool)
		programID = chainState.BurnMintTokenPool
		contractType = ccipChangeset.BurnMintTokenPool
		tokenPoolUsingMcms = cfg.MCMSSolana != nil && cfg.MCMSSolana.BurnMintTokenPoolOwnedByTimelock[tokenPubKey]
		authority, err := GetAuthorityForIxn(
			&e,
			chain,
			cfg.MCMSSolana,
			ccipChangeset.BurnMintTokenPool,
			tokenPubKey)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
		}
		ix, err = solBurnMintTokenPool.NewRemoveFromAllowListInstruction(
			cfg.Accounts,
			poolConfigPDA,
			tokenPubKey,
			authority,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
	case solTestTokenPool.LockAndRelease_PoolType:
		poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, chainState.LockReleaseTokenPool)
		solLockReleaseTokenPool.SetProgramID(chainState.LockReleaseTokenPool)
		programID = chainState.LockReleaseTokenPool
		contractType = ccipChangeset.LockReleaseTokenPool
		tokenPoolUsingMcms = cfg.MCMSSolana != nil && cfg.MCMSSolana.LockReleaseTokenPoolOwnedByTimelock[tokenPubKey]
		authority, err := GetAuthorityForIxn(
			&e,
			chain,
			cfg.MCMSSolana,
			ccipChangeset.LockReleaseTokenPool,
			tokenPubKey)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
		}
		ix, err = solLockReleaseTokenPool.NewRemoveFromAllowListInstruction(
			cfg.Accounts,
			poolConfigPDA,
			tokenPubKey,
			authority,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
	default:
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid pool type: %s", cfg.PoolType)
	}
	if tokenPoolUsingMcms {
		tx, err := BuildMCMSTxn(ix, programID.String(), contractType)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.SolChainSelector, "proposal to RemoveFromTokenPoolAllowList in Solana", cfg.MCMSSolana.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
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
	e.Logger.Infow("Configured token pool allowlist", "token_pubkey", tokenPubKey.String())
	return deployment.ChangesetOutput{}, nil
}

type LockReleaseLiquidityOpsConfig struct {
	SolChainSelector uint64
	SolTokenPubKey   string
	SetCfg           *SetLiquidityConfig
	LiquidityCfg     *LiquidityConfig
	RebalancerCfg    *RebalancerConfig
	MCMSSolana       *MCMSConfigSolana
}

type SetLiquidityConfig struct {
	Enabled bool
}
type LiquidityOperation int

const (
	Provide LiquidityOperation = iota
	Withdraw
)

type LiquidityConfig struct {
	Amount             int
	RemoteTokenAccount solana.PublicKey
	Type               LiquidityOperation
}

type RebalancerConfig struct {
	Rebalancer solana.PublicKey
}

func (cfg LockReleaseLiquidityOpsConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)
	if err := commonValidation(e, cfg.SolChainSelector, tokenPubKey); err != nil {
		return err
	}
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.SolChains[cfg.SolChainSelector]
	chain := e.SolChains[cfg.SolChainSelector]
	if err := ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, tokenPubKey); err != nil {
		return err
	}
	return validatePoolDeployment(&e, solTestTokenPool.LockAndRelease_PoolType, cfg.SolChainSelector, tokenPubKey, true)
}

func LockReleaseLiquidityOps(e deployment.Environment, cfg LockReleaseLiquidityOpsConfig) (deployment.ChangesetOutput, error) {
	e.Logger.Infof("Locking/Unlocking liquidity for token %s", cfg.SolTokenPubKey)
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.SolChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.SolChainSelector]
	tokenPool := chainState.LockReleaseTokenPool

	solLockReleaseTokenPool.SetProgramID(tokenPool)
	programID := tokenPool
	contractType := ccipChangeset.LockReleaseTokenPool
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)
	poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
	tokenPoolUsingMcms := cfg.MCMSSolana != nil && cfg.MCMSSolana.LockReleaseTokenPoolOwnedByTimelock[tokenPubKey]
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.LockReleaseTokenPool,
		tokenPubKey)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
	ixns := make([]solana.Instruction, 0)
	if cfg.SetCfg != nil {
		ix, err := solLockReleaseTokenPool.NewSetCanAcceptLiquidityInstruction(
			cfg.SetCfg.Enabled,
			poolConfigPDA,
			tokenPubKey,
			authority,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
		ixns = append(ixns, ix)
	}
	if cfg.LiquidityCfg != nil {
		tokenProgram, _ := chainState.TokenToTokenProgram(tokenPubKey)
		poolSigner, _ := solTokenUtil.TokenPoolSignerAddress(tokenPubKey, tokenPool)
		poolConfigAccount := solLockReleaseTokenPool.State{}
		_ = chain.GetAccountDataBorshInto(context.Background(), poolConfigPDA, &poolConfigAccount)
		if cfg.LiquidityCfg.Amount <= 0 {
			return deployment.ChangesetOutput{}, fmt.Errorf("invalid amount: %d", cfg.LiquidityCfg.Amount)
		}
		tokenAmount := uint64(cfg.LiquidityCfg.Amount) // #nosec G115 - we check the amount above
		switch cfg.LiquidityCfg.Type {
		case Provide:
			outDec, outVal, err := tokens.TokenBalance(
				e.GetContext(),
				chain.Client,
				cfg.LiquidityCfg.RemoteTokenAccount,
				deployment.SolDefaultCommitment)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to get token balance: %w", err)
			}
			if outVal < cfg.LiquidityCfg.Amount {
				return deployment.ChangesetOutput{}, fmt.Errorf("insufficient token balance: %d < %d", outVal, cfg.LiquidityCfg.Amount)
			}
			ix1, err := solTokenUtil.TokenApproveChecked(
				tokenAmount,
				outDec,
				tokenProgram,
				cfg.LiquidityCfg.RemoteTokenAccount,
				tokenPubKey,
				poolSigner,
				chain.DeployerKey.PublicKey(),
				solana.PublicKeySlice{},
			)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to TokenApproveChecked: %w", err)
			}
			if err = chain.Confirm([]solana.Instruction{ix1}); err != nil {
				e.Logger.Errorw("Failed to confirm instructions for TokenApproveChecked", "chain", chain.String(), "err", err)
				return deployment.ChangesetOutput{}, err
			}
			ix, err := solLockReleaseTokenPool.NewProvideLiquidityInstruction(
				tokenAmount,
				poolConfigPDA,
				tokenProgram,
				tokenPubKey,
				poolSigner,
				poolConfigAccount.Config.PoolTokenAccount,
				cfg.LiquidityCfg.RemoteTokenAccount,
				authority,
			).ValidateAndBuild()
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
			ixns = append(ixns, ix)
		case Withdraw:
			ix, err := solLockReleaseTokenPool.NewWithdrawLiquidityInstruction(
				tokenAmount,
				poolConfigPDA,
				tokenProgram,
				tokenPubKey,
				poolSigner,
				poolConfigAccount.Config.PoolTokenAccount,
				cfg.LiquidityCfg.RemoteTokenAccount,
				authority,
			).ValidateAndBuild()
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
			ixns = append(ixns, ix)
		}
	}
	if cfg.RebalancerCfg != nil {
		ix, err := solLockReleaseTokenPool.NewSetRebalancerInstruction(
			cfg.RebalancerCfg.Rebalancer,
			poolConfigPDA,
			tokenPubKey,
			authority,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
		ixns = append(ixns, ix)
	}

	if tokenPoolUsingMcms {
		txns := make([]mcmsTypes.Transaction, 0)
		for _, ixn := range ixns {
			tx, err := BuildMCMSTxn(ixn, programID.String(), contractType)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
			}
			txns = append(txns, *tx)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.SolChainSelector, "proposal to RemoveFromTokenPoolAllowList in Solana", cfg.MCMSSolana.MCMS.MinDelay, txns)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	err = chain.Confirm(ixns)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return deployment.ChangesetOutput{}, nil
}

type TokenPoolOpsCfg struct {
	SolChainSelector uint64
	SolTokenPubKey   string
	DeleteChainCfg   *DeleteChainCfg
	SetRouterCfg     *SetRouterCfg
	PoolType         solTestTokenPool.PoolType
	MCMSSolana       *MCMSConfigSolana
}

type DeleteChainCfg struct {
	RemoteChainSelector uint64
}

type SetRouterCfg struct {
	Router solana.PublicKey
}

func (cfg TokenPoolOpsCfg) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.SolChains[cfg.SolChainSelector]
	chain := e.SolChains[cfg.SolChainSelector]
	if err := commonValidation(e, cfg.SolChainSelector, tokenPubKey); err != nil {
		return err
	}
	if err := validatePoolDeployment(&e, cfg.PoolType, cfg.SolChainSelector, tokenPubKey, true); err != nil {
		return err
	}
	if cfg.DeleteChainCfg != nil {
		var tokenPool solana.PublicKey
		var remoteChainConfigAccount interface{}

		switch cfg.PoolType {
		case solTestTokenPool.BurnAndMint_PoolType:
			tokenPool = chainState.BurnMintTokenPool
			remoteChainConfigAccount = solBurnMintTokenPool.ChainConfig{}
		case solTestTokenPool.LockAndRelease_PoolType:
			tokenPool = chainState.LockReleaseTokenPool
			remoteChainConfigAccount = solLockReleaseTokenPool.ChainConfig{}
		default:
			return fmt.Errorf("invalid pool type: %s", cfg.PoolType)
		}
		// check if this remote chain is already configured for this token
		remoteChainConfigPDA, _, err := solTokenUtil.TokenPoolChainConfigPDA(cfg.DeleteChainCfg.RemoteChainSelector, tokenPubKey, tokenPool)
		if err != nil {
			return fmt.Errorf("failed to get token pool remote chain config pda (remoteSelector: %d, mint: %s, pool: %s): %w", cfg.DeleteChainCfg.RemoteChainSelector, tokenPubKey.String(), tokenPool.String(), err)
		}
		err = chain.GetAccountDataBorshInto(context.Background(), remoteChainConfigPDA, &remoteChainConfigAccount)

		if err != nil {
			return fmt.Errorf("remote chain config not found for (remoteSelector: %d, mint: %s, pool: %s, type: %s): %w", cfg.DeleteChainCfg.RemoteChainSelector, tokenPubKey.String(), tokenPool.String(), cfg.PoolType, err)
		}
	}
	if cfg.SetRouterCfg != nil {
		if cfg.SetRouterCfg.Router.IsZero() {
			return fmt.Errorf("invalid router address: %s", cfg.SetRouterCfg.Router.String())
		}
	}
	return ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, tokenPubKey)
}

func TokenPoolOps(e deployment.Environment, cfg TokenPoolOpsCfg) (deployment.ChangesetOutput, error) {
	e.Logger.Infof("Setting pool config for token %s", cfg.SolTokenPubKey)
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.SolChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.SolChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)

	var ix solana.Instruction
	var tokenPoolUsingMcms bool
	var programID solana.PublicKey
	var contractType deployment.ContractType
	ixns := make([]solana.Instruction, 0)
	switch cfg.PoolType {
	case solTestTokenPool.BurnAndMint_PoolType:
		poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, chainState.BurnMintTokenPool)
		remoteChainConfigPDA, _, _ := solTokenUtil.TokenPoolChainConfigPDA(cfg.DeleteChainCfg.RemoteChainSelector, tokenPubKey, chainState.BurnMintTokenPool)
		solBurnMintTokenPool.SetProgramID(chainState.BurnMintTokenPool)
		programID = chainState.BurnMintTokenPool
		contractType = ccipChangeset.BurnMintTokenPool
		tokenPoolUsingMcms = cfg.MCMSSolana != nil && cfg.MCMSSolana.BurnMintTokenPoolOwnedByTimelock[tokenPubKey]
		authority, err := GetAuthorityForIxn(
			&e,
			chain,
			cfg.MCMSSolana,
			ccipChangeset.BurnMintTokenPool,
			tokenPubKey)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
		}
		if cfg.DeleteChainCfg != nil {
			ix, err = solBurnMintTokenPool.NewDeleteChainConfigInstruction(
				cfg.DeleteChainCfg.RemoteChainSelector,
				tokenPubKey,
				poolConfigPDA,
				remoteChainConfigPDA,
				authority,
			).ValidateAndBuild()
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
			ixns = append(ixns, ix)
		}
		if cfg.SetRouterCfg != nil {
			ix, err = solBurnMintTokenPool.NewSetRouterInstruction(
				cfg.SetRouterCfg.Router,
				poolConfigPDA,
				tokenPubKey,
				authority,
			).ValidateAndBuild()
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
			ixns = append(ixns, ix)
		}
	case solTestTokenPool.LockAndRelease_PoolType:
		poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, chainState.LockReleaseTokenPool)
		remoteChainConfigPDA, _, _ := solTokenUtil.TokenPoolChainConfigPDA(cfg.DeleteChainCfg.RemoteChainSelector, tokenPubKey, chainState.LockReleaseTokenPool)
		solLockReleaseTokenPool.SetProgramID(chainState.LockReleaseTokenPool)
		programID = chainState.LockReleaseTokenPool
		contractType = ccipChangeset.LockReleaseTokenPool
		tokenPoolUsingMcms = cfg.MCMSSolana != nil && cfg.MCMSSolana.LockReleaseTokenPoolOwnedByTimelock[tokenPubKey]
		authority, err := GetAuthorityForIxn(
			&e,
			chain,
			cfg.MCMSSolana,
			ccipChangeset.LockReleaseTokenPool,
			tokenPubKey)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
		}
		if cfg.DeleteChainCfg != nil {
			ix, err = solLockReleaseTokenPool.NewDeleteChainConfigInstruction(
				cfg.DeleteChainCfg.RemoteChainSelector,
				tokenPubKey,
				poolConfigPDA,
				remoteChainConfigPDA,
				authority,
			).ValidateAndBuild()
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
			ixns = append(ixns, ix)
		}
		if cfg.SetRouterCfg != nil {
			ix, err = solLockReleaseTokenPool.NewSetRouterInstruction(
				cfg.SetRouterCfg.Router,
				poolConfigPDA,
				tokenPubKey,
				authority,
			).ValidateAndBuild()
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
			ixns = append(ixns, ix)
		}
	default:
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid pool type: %s", cfg.PoolType)
	}
	if tokenPoolUsingMcms {
		tx, err := BuildMCMSTxn(ix, programID.String(), contractType)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.SolChainSelector, "proposal to ConfigureTokenPoolAllowList in Solana", cfg.MCMSSolana.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	if err := chain.Confirm(ixns); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	e.Logger.Infow("Configured token pool allowlist", "token_pubkey", tokenPubKey.String())
	return deployment.ChangesetOutput{}, nil
}
