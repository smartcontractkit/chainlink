package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"

	solBaseTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/base_token_pool"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solBurnMintTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/example_burnmint_token_pool"
	solLockReleaseTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/example_lockrelease_token_pool"
	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

var _ deployment.ChangeSet[TokenPoolConfig] = AddTokenPool
var _ deployment.ChangeSet[RemoteChainTokenPoolConfig] = SetupTokenPoolForRemoteChain

func validatePoolDeployment(s ccipChangeset.SolCCIPChainState, poolType solTestTokenPool.PoolType, selector uint64) error {
	switch poolType {
	case solTestTokenPool.BurnAndMint_PoolType:
		if s.BurnMintTokenPool.IsZero() {
			return fmt.Errorf("token pool of type BurnAndMint not found in existing state, deploy the token pool first for chain %d", selector)
		}
	case solTestTokenPool.LockAndRelease_PoolType:
		if s.LockReleaseTokenPool.IsZero() {
			return fmt.Errorf("token pool of type LockAndRelease not found in existing state, deploy the token pool first for chain %d", selector)
		}
	default:
		return fmt.Errorf("invalid pool type: %s", poolType)
	}
	return nil
}

type TokenPoolConfig struct {
	ChainSelector uint64
	PoolType      solTestTokenPool.PoolType
	Authority     string
	TokenPubKey   string
}

func (cfg TokenPoolConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]

	if _, err := chainState.TokenToTokenProgram(tokenPubKey); err != nil {
		return fmt.Errorf("failed to get token program for token address %s: %w", tokenPubKey.String(), err)
	}

	if err := validatePoolDeployment(chainState, cfg.PoolType, cfg.ChainSelector); err != nil {
		return err
	}

	var tokenPool solana.PublicKey
	var poolConfigAccount interface{}

	switch cfg.PoolType {
	case solTestTokenPool.BurnAndMint_PoolType:
		tokenPool = chainState.BurnMintTokenPool
		poolConfigAccount = solBurnMintTokenPool.State{}
	case solTestTokenPool.LockAndRelease_PoolType:
		tokenPool = chainState.LockReleaseTokenPool
		poolConfigAccount = solLockReleaseTokenPool.State{}
	default:
		return fmt.Errorf("invalid pool type: %s", cfg.PoolType)
	}

	poolConfigPDA, err := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
	if err != nil {
		return fmt.Errorf("failed to get token pool config address (mint: %s, pool: %s): %w", tokenPubKey.String(), tokenPool.String(), err)
	}
	chain := e.SolChains[cfg.ChainSelector]
	if err := chain.GetAccountDataBorshInto(context.Background(), poolConfigPDA, &poolConfigAccount); err == nil {
		return fmt.Errorf("token pool config already exists for (mint: %s, pool: %s, type: %s)", tokenPubKey.String(), tokenPool.String(), cfg.PoolType)
	}
	return nil
}

func AddTokenPool(e deployment.Environment, cfg TokenPoolConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	authorityPubKey := solana.MustPublicKeyFromBase58(cfg.Authority)
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
			chainState.Router,
			poolConfigPDA,
			tokenPubKey,
			authorityPubKey, // this is assumed to be chain.DeployerKey for now (owner of token pool)
			solana.SystemProgramID,
		).ValidateAndBuild()
	case solTestTokenPool.LockAndRelease_PoolType:
		// initialize token pool for token
		poolInitI, err = solLockReleaseTokenPool.NewInitializeInstruction(
			chainState.Router,
			poolConfigPDA,
			tokenPubKey,
			authorityPubKey, // this is assumed to be chain.DeployerKey for now (owner of token pool)
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
			authorityPubKey,
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
}

func (cfg RemoteChainTokenPoolConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)
	if err := commonValidation(e, cfg.SolChainSelector, tokenPubKey); err != nil {
		return err
	}
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.SolChainSelector]
	chain := e.SolChains[cfg.SolChainSelector]

	if err := validatePoolDeployment(chainState, cfg.PoolType, cfg.SolChainSelector); err != nil {
		return err
	}

	var tokenPool solana.PublicKey
	var poolConfigAccount interface{}
	var remoteChainConfigAccount interface{}

	switch cfg.PoolType {
	case solTestTokenPool.BurnAndMint_PoolType:
		tokenPool = chainState.BurnMintTokenPool
		poolConfigAccount = solBurnMintTokenPool.State{}
		remoteChainConfigAccount = solBurnMintTokenPool.ChainConfig{}
	case solTestTokenPool.LockAndRelease_PoolType:
		tokenPool = chainState.LockReleaseTokenPool
		poolConfigAccount = solLockReleaseTokenPool.State{}
		remoteChainConfigAccount = solLockReleaseTokenPool.ChainConfig{}
	default:
		return fmt.Errorf("invalid pool type: %s", cfg.PoolType)
	}

	// check if pool config exists (cannot do remote setup without it)
	poolConfigPDA, err := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
	if err != nil {
		return fmt.Errorf("failed to get token pool config address (mint: %s, pool: %s): %w", tokenPubKey.String(), tokenPool.String(), err)
	}
	if err := chain.GetAccountDataBorshInto(context.Background(), poolConfigPDA, &poolConfigAccount); err != nil {
		return fmt.Errorf("token pool config not found (mint: %s, pool: %s, type: %s): %w", tokenPubKey.String(), tokenPool.String(), cfg.PoolType, err)
	}

	// check if this remote chain is already configured for this token
	remoteChainConfigPDA, _, err := solTokenUtil.TokenPoolChainConfigPDA(cfg.RemoteChainSelector, tokenPubKey, tokenPool)
	if err != nil {
		return fmt.Errorf("failed to get token pool remote chain config pda (remoteSelector: %d, mint: %s, pool: %s): %w", cfg.RemoteChainSelector, tokenPubKey.String(), tokenPool.String(), err)
	}
	if err := chain.GetAccountDataBorshInto(context.Background(), remoteChainConfigPDA, &remoteChainConfigAccount); err == nil {
		return fmt.Errorf("remote chain config already exists for (remoteSelector: %d, mint: %s, pool: %s, type: %s)", cfg.RemoteChainSelector, tokenPubKey.String(), tokenPool.String(), cfg.PoolType)
	}
	return nil
}

func SetupTokenPoolForRemoteChain(e deployment.Environment, cfg RemoteChainTokenPoolConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.SolChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.SolChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)

	var instructions []solana.Instruction
	var err error
	switch cfg.PoolType {
	case solTestTokenPool.BurnAndMint_PoolType:
		instructions, err = getInstructionsForBurnMint(chain, chainState, cfg)
	case solTestTokenPool.LockAndRelease_PoolType:
		instructions, err = getInstructionsForLockRelease(chain, chainState, cfg)
	default:
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid pool type: %s", cfg.PoolType)
	}
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}

	err = chain.Confirm(instructions)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	e.Logger.Infow("Configured token pool for remote chain", "remote_chain_selector", cfg.RemoteChainSelector, "token_pubkey", tokenPubKey.String())
	return deployment.ChangesetOutput{}, nil
}

func getInstructionsForBurnMint(
	chain deployment.SolChain,
	chainState ccipChangeset.SolCCIPChainState,
	cfg RemoteChainTokenPoolConfig,
) ([]solana.Instruction, error) {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)
	poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, chainState.BurnMintTokenPool)
	remoteChainConfigPDA, _, _ := solTokenUtil.TokenPoolChainConfigPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.BurnMintTokenPool)
	solBurnMintTokenPool.SetProgramID(chainState.BurnMintTokenPool)
	ixConfigure, err := solBurnMintTokenPool.NewInitChainRemoteConfigInstruction(
		cfg.RemoteChainSelector,
		tokenPubKey,
		cfg.RemoteConfig,
		poolConfigPDA,
		remoteChainConfigPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("failed to generate instructions: %w", err)
	}
	ixRates, err := solBurnMintTokenPool.NewSetChainRateLimitInstruction(
		cfg.RemoteChainSelector,
		tokenPubKey,
		cfg.InboundRateLimit,
		cfg.OutboundRateLimit,
		poolConfigPDA,
		remoteChainConfigPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("failed to generate instructions: %w", err)
	}
	ixAppend, err := solBurnMintTokenPool.NewAppendRemotePoolAddressesInstruction(
		cfg.RemoteChainSelector,
		tokenPubKey,
		cfg.RemoteConfig.PoolAddresses, // i dont know why this is a list (is it for different types of pool of the same token?)
		poolConfigPDA,
		remoteChainConfigPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("failed to generate instructions: %w", err)
	}
	return []solana.Instruction{ixConfigure, ixRates, ixAppend}, nil
}

func getInstructionsForLockRelease(
	chain deployment.SolChain,
	chainState ccipChangeset.SolCCIPChainState,
	cfg RemoteChainTokenPoolConfig,
) ([]solana.Instruction, error) {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)
	poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, chainState.LockReleaseTokenPool)
	remoteChainConfigPDA, _, _ := solTokenUtil.TokenPoolChainConfigPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.LockReleaseTokenPool)
	solLockReleaseTokenPool.SetProgramID(chainState.LockReleaseTokenPool)
	ixConfigure, err := solLockReleaseTokenPool.NewInitChainRemoteConfigInstruction(
		cfg.RemoteChainSelector,
		tokenPubKey,
		cfg.RemoteConfig,
		poolConfigPDA,
		remoteChainConfigPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("failed to generate instructions: %w", err)
	}
	ixRates, err := solLockReleaseTokenPool.NewSetChainRateLimitInstruction(
		cfg.RemoteChainSelector,
		tokenPubKey,
		cfg.InboundRateLimit,
		cfg.OutboundRateLimit,
		poolConfigPDA,
		remoteChainConfigPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("failed to generate instructions: %w", err)
	}
	ixAppend, err := solLockReleaseTokenPool.NewAppendRemotePoolAddressesInstruction(
		cfg.RemoteChainSelector,
		tokenPubKey,
		cfg.RemoteConfig.PoolAddresses, // i dont know why this is a list (is it for different types of pool of the same token?)
		poolConfigPDA,
		remoteChainConfigPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("failed to generate instructions: %w", err)
	}
	return []solana.Instruction{ixConfigure, ixRates, ixAppend}, nil
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
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	_, err := chainState.TokenToTokenProgram(tokenPubKey)
	if err != nil {
		return fmt.Errorf("failed to get token program for token address %s: %w", tokenPubKey.String(), err)
	}
	return validatePoolDeployment(chainState, cfg.PoolType, cfg.ChainSelector)
}

func AddTokenPoolLookupTable(e deployment.Environment, cfg TokenPoolLookupTableConfig) (deployment.ChangesetOutput, error) {
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
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)
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
	ChainSelector                     uint64
	TokenPubKey                       string
	TokenAdminRegistryAdminPrivateKey string
	WritableIndexes                   []uint8
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
	tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)
	if err != nil {
		return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), chainState.Router.String(), err)
	}
	var tokenAdminRegistryAccount solRouter.TokenAdminRegistry
	if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err != nil {
		return fmt.Errorf("token admin registry not found for (mint: %s, router: %s), cannot set pool", tokenPubKey.String(), chainState.Router.String())
	}
	if _, ok := chainState.TokenPoolLookupTable[tokenPubKey]; !ok {
		return fmt.Errorf("token pool lookup table not found for (mint: %s)", tokenPubKey.String())
	}
	return nil
}

// this sets the writable indexes of the token pool lookup table
func SetPool(e deployment.Environment, cfg SetPoolConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)
	tokenAdminRegistryAdminPrivKey := solana.MustPrivateKeyFromBase58(cfg.TokenAdminRegistryAdminPrivateKey)
	lookupTablePubKey := chainState.TokenPoolLookupTable[tokenPubKey]

	base := solRouter.NewSetPoolInstruction(
		cfg.WritableIndexes,
		routerConfigPDA,
		tokenAdminRegistryPDA,
		tokenPubKey,
		lookupTablePubKey,
		tokenAdminRegistryAdminPrivKey.PublicKey(),
	)

	base.AccountMetaSlice = append(base.AccountMetaSlice, solana.Meta(lookupTablePubKey))
	instruction, err := base.ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	instructions := []solana.Instruction{instruction}
	err = chain.Confirm(instructions, solCommonUtil.AddSigners(tokenAdminRegistryAdminPrivKey))
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	e.Logger.Infow("Set pool config", "token_pubkey", tokenPubKey.String())
	return deployment.ChangesetOutput{}, nil
}
