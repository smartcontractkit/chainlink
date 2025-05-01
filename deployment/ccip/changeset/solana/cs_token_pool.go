package solana

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	ccipChangeset_v1_5_1 "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
)

// use this changeset to add a token pool and lookup table
var _ deployment.ChangeSet[TokenPoolConfig] = AddTokenPoolAndLookupTable

// use this changeset to setup a token pool for a remote chain
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

// append mcms txns generated from solanainstructions
func appendTxs(instructions []solana.Instruction, tokenPool solana.PublicKey, poolType deployment.ContractType, txns *[]mcmsTypes.Transaction) error {
	for _, ixn := range instructions {
		tx, err := BuildMCMSTxn(ixn, tokenPool.String(), poolType)
		if err != nil {
			return fmt.Errorf("failed to generate mcms txn: %w", err)
		}
		if tx == nil {
			return errors.New("mcms txn unexpectedly nil")
		}
		*txns = append(*txns, *tx)
	}
	return nil
}

// get diff of pool addresses
func poolDiff(existingPoolAddresses []solBaseTokenPool.RemoteAddress, newPoolAddresses []solBaseTokenPool.RemoteAddress) []solBaseTokenPool.RemoteAddress {
	var result []solBaseTokenPool.RemoteAddress
	// for every new address, check if it exists in the existing pool addresses
	for _, newAddr := range newPoolAddresses {
		exists := false
		for _, existingAddr := range existingPoolAddresses {
			if bytes.Equal(existingAddr.Address, newAddr.Address) {
				exists = true
				break
			}
		}
		if !exists {
			result = append(result, newAddr)
		}
	}
	return result
}

// get pool pdas
func getPoolPDAs(
	solTokenPubKey solana.PublicKey, poolAddress solana.PublicKey, remoteChainSelector uint64,
) (poolConfigPDA solana.PublicKey, remoteChainConfigPDA solana.PublicKey) {
	poolConfigPDA, _ = solTokenUtil.TokenPoolConfigAddress(solTokenPubKey, poolAddress)
	remoteChainConfigPDA, _, _ = solTokenUtil.TokenPoolChainConfigPDA(remoteChainSelector, solTokenPubKey, poolAddress)
	return poolConfigPDA, remoteChainConfigPDA
}

type TokenPoolConfig struct {
	ChainSelector uint64
	PoolType      solTestTokenPool.PoolType
	TokenPubKey   solana.PublicKey
}

func (cfg TokenPoolConfig) Validate(e deployment.Environment) error {
	if err := commonValidation(e, cfg.ChainSelector, cfg.TokenPubKey); err != nil {
		return err
	}

	return validatePoolDeployment(&e, cfg.PoolType, cfg.ChainSelector, cfg.TokenPubKey, false)
}

func AddTokenPoolAndLookupTable(e deployment.Environment, cfg TokenPoolConfig) (deployment.ChangesetOutput, error) {
	e.Logger.Infow("Adding token pool", "token_pubkey", cfg.TokenPubKey)
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := cfg.TokenPubKey
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
	programData, err := getSolProgramData(e, chain, tokenPool)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get solana token pool program data: %w", err)
	}
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
			tokenPool,
			programData.Address,
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
			tokenPool,
			programData.Address,
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
		e.Logger.Infow("Setting mint authority", "poolSigner", poolSigner.String())
	}

	// add signer here if authority is different from deployer key
	if err := chain.Confirm(instructions); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	e.Logger.Infow("Created new token pool config", "token_pool_ata", tokenPoolATA.String(), "pool_config", poolConfigPDA.String(), "pool_signer", poolSigner.String())

	csOutput, err := AddTokenPoolLookupTable(e, TokenPoolLookupTableConfig{
		ChainSelector: cfg.ChainSelector,
		TokenPubKey:   cfg.TokenPubKey,
		PoolType:      cfg.PoolType,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to add token pool lookup table: %w", err)
	}

	return csOutput, nil
}

// SETUP REMOTE CHAIN TOKEN POOL FOR A GIVEN TOKEN
type RateLimiterConfig struct {
	// Inbound is the rate limiter config for inbound transfers from a remote chain.
	Inbound solBaseTokenPool.RateLimitConfig
	// Outbound is the rate limiter config for outbound transfers to a remote chain.
	Outbound solBaseTokenPool.RateLimitConfig
}

func validateRateLimiterConfig(rateLimiterConfig solBaseTokenPool.RateLimitConfig) error {
	if rateLimiterConfig.Enabled {
		if rateLimiterConfig.Rate >= rateLimiterConfig.Capacity || rateLimiterConfig.Rate == 0 {
			return errors.New("rate must be greater than 0 and less than capacity if enabled")
		}
	} else {
		if rateLimiterConfig.Rate != 0 || rateLimiterConfig.Capacity != 0 {
			return errors.New("rate and capacity must be 0 if not enabled")
		}
	}
	return nil
}

func (cfg RateLimiterConfig) Validate() error {
	if err := validateRateLimiterConfig(cfg.Inbound); err != nil {
		return err
	}
	err := validateRateLimiterConfig(cfg.Outbound)
	return err
}

type EVMRemoteConfig struct {
	TokenSymbol       ccipChangeset.TokenSymbol
	PoolType          deployment.ContractType
	PoolVersion       semver.Version
	RateLimiterConfig RateLimiterConfig
	OverrideConfig    bool
}

func (cfg EVMRemoteConfig) Validate(e deployment.Environment, state ccipChangeset.CCIPOnChainState, evmChainSelector uint64) error {
	// add evm family check
	if cfg.TokenSymbol == "" {
		return errors.New("token symbol must be defined")
	}
	err := deployment.IsValidChainSelector(evmChainSelector)
	if err != nil {
		return fmt.Errorf("failed to validate chain selector %d: %w", evmChainSelector, err)
	}
	chain, ok := e.Chains[evmChainSelector]
	if !ok {
		return fmt.Errorf("chain with selector %d does not exist in environment", evmChainSelector)
	}
	chainState, ok := state.Chains[evmChainSelector]
	if !ok {
		return fmt.Errorf("%s does not exist in state", chain.String())
	}
	// Ensure that the inputted type is known
	if _, typeOk := ccipChangeset.TokenPoolTypes[cfg.PoolType]; !typeOk {
		return fmt.Errorf("%s is not a known token pool type", cfg.PoolType)
	}
	// Ensure that the inputted version is known
	if _, versionOk := ccipChangeset.TokenPoolVersions[cfg.PoolVersion]; !versionOk {
		return fmt.Errorf("%s is not a known token pool version", cfg.PoolVersion)
	}
	// Ensure that a pool with given symbol, type and version is known to the environment
	_, getPoolOk := ccipChangeset.GetTokenPoolAddressFromSymbolTypeAndVersion(chainState, chain, cfg.TokenSymbol, cfg.PoolType, cfg.PoolVersion)
	if !getPoolOk {
		return fmt.Errorf("token pool does not exist on %s with symbol %s, type %s, and version %s", chain.String(), cfg.TokenSymbol, cfg.PoolType, cfg.PoolVersion)
	}
	err = cfg.RateLimiterConfig.Validate()
	return err
}

type RemoteChainTokenPoolConfig struct {
	SolChainSelector uint64
	SolTokenPubKey   solana.PublicKey
	SolPoolType      solTestTokenPool.PoolType
	EVMRemoteConfigs map[uint64]EVMRemoteConfig
	MCMSSolana       *MCMSConfigSolana
}

func (cfg RemoteChainTokenPoolConfig) Validate(e deployment.Environment) error {
	if err := commonValidation(e, cfg.SolChainSelector, cfg.SolTokenPubKey); err != nil {
		return err
	}
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.SolChainSelector]
	chain := e.SolChains[cfg.SolChainSelector]

	if err := validatePoolDeployment(&e, cfg.SolPoolType, cfg.SolChainSelector, cfg.SolTokenPubKey, true); err != nil {
		return err
	}

	if err := ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, cfg.SolTokenPubKey); err != nil {
		return err
	}
	// validate EVMRemoteConfig
	for evmChainSelector, evmRemoteConfig := range cfg.EVMRemoteConfigs {
		if err := evmRemoteConfig.Validate(e, state, evmChainSelector); err != nil {
			return err
		}
	}
	return nil
}

func getOnChainEVMPoolConfig(e deployment.Environment, state ccipChangeset.CCIPOnChainState, evmChainSelector uint64, evmRemoteConfig EVMRemoteConfig) (solBaseTokenPool.RemoteConfig, error) {
	evmChain := e.Chains[evmChainSelector]
	evmChainState := state.Chains[evmChainSelector]
	evmTokenPool, evmTokenAddress, _, evmErr := ccipChangeset_v1_5_1.GetTokenStateFromPoolEVM(context.Background(), evmRemoteConfig.TokenSymbol, evmRemoteConfig.PoolType, evmRemoteConfig.PoolVersion, evmChain, evmChainState)
	if evmErr != nil {
		return solBaseTokenPool.RemoteConfig{}, fmt.Errorf("failed to get token evm token pool and token address: %w", evmErr)
	}
	evmTokenPoolAddress := evmTokenPool.Address()
	evmTokenDecimals, err := evmTokenPool.GetTokenDecimals(&bind.CallOpts{Context: context.Background()})
	if err != nil {
		return solBaseTokenPool.RemoteConfig{}, fmt.Errorf("failed to get token decimals: %w", err)
	}
	onChainEVMRemoteConfig := solBaseTokenPool.RemoteConfig{
		TokenAddress: solBaseTokenPool.RemoteAddress{
			Address: common.LeftPadBytes(evmTokenAddress.Bytes(), 32),
		},
		PoolAddresses: []solBaseTokenPool.RemoteAddress{
			{
				Address: evmTokenPoolAddress.Bytes(),
			},
		},
		Decimals: evmTokenDecimals,
	}
	return onChainEVMRemoteConfig, nil
}

func SetupTokenPoolForRemoteChain(e deployment.Environment, cfg RemoteChainTokenPoolConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.SolChainSelector]
	envState, _ := ccipChangeset.LoadOnchainState(e)
	solChainState := envState.SolChains[cfg.SolChainSelector]
	tokenPubKey := cfg.SolTokenPubKey

	var instructions []solana.Instruction
	var txns []mcmsTypes.Transaction
	switch cfg.SolPoolType {
	case solTestTokenPool.BurnAndMint_PoolType:
		for evmChainSelector, evmRemoteConfig := range cfg.EVMRemoteConfigs {
			e.Logger.Infow("Setting up token pool for remote chain", "remote_chain_selector", evmChainSelector, "token_pubkey", tokenPubKey.String(), "pool_type", cfg.SolPoolType)
			chainIxs, err := getInstructionsForBurnMint(e, chain, envState, solChainState, cfg, evmChainSelector, evmRemoteConfig)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
			instructions = append(instructions, chainIxs...)
		}
		if cfg.MCMSSolana != nil && cfg.MCMSSolana.BurnMintTokenPoolOwnedByTimelock[tokenPubKey] {
			err := appendTxs(instructions, solChainState.BurnMintTokenPool, ccipChangeset.BurnMintTokenPool, &txns)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate mcms txn: %w", err)
			}
		}
	case solTestTokenPool.LockAndRelease_PoolType:
		for evmChainSelector, evmRemoteConfig := range cfg.EVMRemoteConfigs {
			e.Logger.Infow("Setting up token pool for remote chain", "remote_chain_selector", evmChainSelector, "token_pubkey", tokenPubKey.String(), "pool_type", cfg.SolPoolType)
			chainIxs, err := getInstructionsForLockRelease(e, chain, envState, solChainState, cfg, evmChainSelector, evmRemoteConfig)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
			instructions = append(instructions, chainIxs...)
		}
		if cfg.MCMSSolana != nil && cfg.MCMSSolana.LockReleaseTokenPoolOwnedByTimelock[tokenPubKey] {
			err := appendTxs(instructions, solChainState.LockReleaseTokenPool, ccipChangeset.LockReleaseTokenPool, &txns)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate mcms txn: %w", err)
			}
		}
	default:
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid pool type: %s", cfg.SolPoolType)
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
	if err := chain.Confirm(instructions); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return deployment.ChangesetOutput{}, nil
}

// checks if the evmChainSelector is supported for the given token and pool type
func isSupportedChain(chain deployment.SolChain, solTokenPubKey solana.PublicKey, solPoolAddress solana.PublicKey, evmChainSelector uint64) (bool, solTestTokenPool.ChainConfig, error) {
	var remoteChainConfigAccount solTestTokenPool.ChainConfig
	// check if this remote chain is already configured for this token
	remoteChainConfigPDA, _, err := solTokenUtil.TokenPoolChainConfigPDA(evmChainSelector, solTokenPubKey, solPoolAddress)
	if err != nil {
		return false, solTestTokenPool.ChainConfig{}, fmt.Errorf("failed to get token pool remote chain config pda (remoteSelector: %d, mint: %s, pool: %s): %w", evmChainSelector, solTokenPubKey.String(), solPoolAddress.String(), err)
	}
	err = chain.GetAccountDataBorshInto(context.Background(), remoteChainConfigPDA, &remoteChainConfigAccount)
	if err != nil { // not a supported chain for this combination of token and pool type
		return false, solTestTokenPool.ChainConfig{}, nil
	}
	return true, remoteChainConfigAccount, nil
}

func getNewSetuptInstructionsForBurnMint(
	e deployment.Environment,
	chain deployment.SolChain,
	chainState ccipChangeset.SolCCIPChainState,
	cfg RemoteChainTokenPoolConfig,
	evmChainSelector uint64,
	rateLimiterConfig RateLimiterConfig,
	onChainEVMPoolConfig solBaseTokenPool.RemoteConfig,
) ([]solana.Instruction, error) {
	tokenPubKey := cfg.SolTokenPubKey
	poolConfigPDA, remoteChainConfigPDA := getPoolPDAs(tokenPubKey, chainState.BurnMintTokenPool, evmChainSelector)
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
	onChainEVMPoolConfigWithoutPoolAddress := solBaseTokenPool.RemoteConfig{
		TokenAddress:  onChainEVMPoolConfig.TokenAddress,
		PoolAddresses: []solBaseTokenPool.RemoteAddress{},
		Decimals:      onChainEVMPoolConfig.Decimals,
	}

	ixConfigure, err := solBurnMintTokenPool.NewInitChainRemoteConfigInstruction(
		evmChainSelector,
		tokenPubKey,
		onChainEVMPoolConfigWithoutPoolAddress,
		poolConfigPDA,
		remoteChainConfigPDA,
		authority,
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("failed to generate instructions: %w", err)
	}
	ixns = append(ixns, ixConfigure)

	ixRates, err := solBurnMintTokenPool.NewSetChainRateLimitInstruction(
		evmChainSelector,
		tokenPubKey,
		rateLimiterConfig.Inbound,
		rateLimiterConfig.Outbound,
		poolConfigPDA,
		remoteChainConfigPDA,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("failed to generate instructions: %w", err)
	}
	ixns = append(ixns, ixRates)

	ixAppend, err := solBurnMintTokenPool.NewAppendRemotePoolAddressesInstruction(
		evmChainSelector,
		tokenPubKey,
		onChainEVMPoolConfig.PoolAddresses, // evm supports multiple remote pools per token
		poolConfigPDA,
		remoteChainConfigPDA,
		authority,
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("failed to generate instructions: %w", err)
	}
	ixns = append(ixns, ixAppend)

	return ixns, nil
}

func getInstructionsForBurnMint(
	e deployment.Environment,
	chain deployment.SolChain,
	envState ccipChangeset.CCIPOnChainState,
	solChainState ccipChangeset.SolCCIPChainState,
	cfg RemoteChainTokenPoolConfig,
	evmChainSelector uint64,
	evmRemoteConfig EVMRemoteConfig,
) ([]solana.Instruction, error) {
	tokenPubKey := cfg.SolTokenPubKey
	poolConfigPDA, remoteChainConfigPDA := getPoolPDAs(tokenPubKey, solChainState.BurnMintTokenPool, evmChainSelector)
	solBurnMintTokenPool.SetProgramID(solChainState.BurnMintTokenPool)
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

	onChainEVMPoolConfig, err := getOnChainEVMPoolConfig(e, envState, evmChainSelector, evmRemoteConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get onchain evm config: %w", err)
	}

	isSupportedChain, remoteChainConfigAccount, err := isSupportedChain(chain, tokenPubKey, solChainState.BurnMintTokenPool, evmChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to check if chain is supported: %w", err)
	}

	if isSupportedChain {
		// override the rate limits if the chain is already supported
		e.Logger.Infof("overriding rate limits for chain %d", evmChainSelector)
		ixRates, err := solBurnMintTokenPool.NewSetChainRateLimitInstruction(
			evmChainSelector,
			tokenPubKey,
			evmRemoteConfig.RateLimiterConfig.Inbound,
			evmRemoteConfig.RateLimiterConfig.Outbound,
			poolConfigPDA,
			remoteChainConfigPDA,
			authority,
		).ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("failed to generate instructions: %w", err)
		}
		ixns = append(ixns, ixRates)

		// if the token address has changed or if the override config flag is set, edit the remote config (just overwrite the existing remote config)
		if !bytes.Equal(remoteChainConfigAccount.Base.Remote.TokenAddress.Address, onChainEVMPoolConfig.TokenAddress.Address) || evmRemoteConfig.OverrideConfig {
			e.Logger.Infof("overriding remote config for chain %d", evmChainSelector)
			ixConfigure, err := solBurnMintTokenPool.NewEditChainRemoteConfigInstruction(
				evmChainSelector,
				tokenPubKey,
				onChainEVMPoolConfig,
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
			// diff between [existing remote pool addresses on solana chain] vs [what was just derived from evm chain]
			diff := poolDiff(remoteChainConfigAccount.Base.Remote.PoolAddresses, onChainEVMPoolConfig.PoolAddresses)
			if len(diff) > 0 {
				e.Logger.Infof("adding new pool addresses for chain %d", evmChainSelector)
				ixAppend, err := solBurnMintTokenPool.NewAppendRemotePoolAddressesInstruction(
					evmChainSelector,
					tokenPubKey,
					diff, // evm supports multiple remote pools per token
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
		}
	} else {
		ixns, err = getNewSetuptInstructionsForBurnMint(e, chain, solChainState, cfg, evmChainSelector, evmRemoteConfig.RateLimiterConfig, onChainEVMPoolConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to generate instructions: %w", err)
		}
	}

	return ixns, nil
}

func getNewSetuptInstructionsForLockRelease(
	e deployment.Environment,
	chain deployment.SolChain,
	chainState ccipChangeset.SolCCIPChainState,
	cfg RemoteChainTokenPoolConfig,
	evmChainSelector uint64,
	rateLimiterConfig RateLimiterConfig,
	onChainEVMPoolConfig solBaseTokenPool.RemoteConfig,
) ([]solana.Instruction, error) {
	tokenPubKey := cfg.SolTokenPubKey
	poolConfigPDA, remoteChainConfigPDA := getPoolPDAs(tokenPubKey, chainState.LockReleaseTokenPool, evmChainSelector)
	ixns := make([]solana.Instruction, 0)
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.LockReleaseTokenPool,
		tokenPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
	onChainEVMPoolConfigWithoutPoolAddress := solBaseTokenPool.RemoteConfig{
		TokenAddress:  onChainEVMPoolConfig.TokenAddress,
		PoolAddresses: []solBaseTokenPool.RemoteAddress{},
		Decimals:      onChainEVMPoolConfig.Decimals,
	}

	ixConfigure, err := solLockReleaseTokenPool.NewInitChainRemoteConfigInstruction(
		evmChainSelector,
		tokenPubKey,
		onChainEVMPoolConfigWithoutPoolAddress,
		poolConfigPDA,
		remoteChainConfigPDA,
		authority,
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("failed to generate instructions: %w", err)
	}
	ixns = append(ixns, ixConfigure)

	ixRates, err := solLockReleaseTokenPool.NewSetChainRateLimitInstruction(
		evmChainSelector,
		tokenPubKey,
		rateLimiterConfig.Inbound,
		rateLimiterConfig.Outbound,
		poolConfigPDA,
		remoteChainConfigPDA,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("failed to generate instructions: %w", err)
	}
	ixns = append(ixns, ixRates)

	ixAppend, err := solLockReleaseTokenPool.NewAppendRemotePoolAddressesInstruction(
		evmChainSelector,
		tokenPubKey,
		onChainEVMPoolConfig.PoolAddresses, // evm supports multiple remote pools per token
		poolConfigPDA,
		remoteChainConfigPDA,
		authority,
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("failed to generate instructions: %w", err)
	}
	ixns = append(ixns, ixAppend)

	return ixns, nil
}

func getInstructionsForLockRelease(
	e deployment.Environment,
	chain deployment.SolChain,
	envState ccipChangeset.CCIPOnChainState,
	solChainState ccipChangeset.SolCCIPChainState,
	cfg RemoteChainTokenPoolConfig,
	evmChainSelector uint64,
	evmRemoteConfig EVMRemoteConfig,
) ([]solana.Instruction, error) {
	tokenPubKey := cfg.SolTokenPubKey
	poolConfigPDA, remoteChainConfigPDA := getPoolPDAs(tokenPubKey, solChainState.LockReleaseTokenPool, evmChainSelector)
	solLockReleaseTokenPool.SetProgramID(solChainState.LockReleaseTokenPool)
	ixns := make([]solana.Instruction, 0)
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.LockReleaseTokenPool,
		tokenPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get authority for ixn: %w", err)
	}

	onChainEVMPoolConfig, err := getOnChainEVMPoolConfig(e, envState, evmChainSelector, evmRemoteConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get on chain evm pool config: %w", err)
	}

	isSupportedChain, remoteChainConfigAccount, err := isSupportedChain(chain, tokenPubKey, solChainState.LockReleaseTokenPool, evmChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to check if chain is supported: %w", err)
	}
	if isSupportedChain {
		// override the rate limits if the chain is already supported
		e.Logger.Infof("overriding rate limits for chain %d", evmChainSelector)
		ixRates, err := solLockReleaseTokenPool.NewSetChainRateLimitInstruction(
			evmChainSelector,
			tokenPubKey,
			evmRemoteConfig.RateLimiterConfig.Inbound,
			evmRemoteConfig.RateLimiterConfig.Outbound,
			poolConfigPDA,
			remoteChainConfigPDA,
			authority,
		).ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("failed to generate instructions: %w", err)
		}
		ixns = append(ixns, ixRates)
		if !bytes.Equal(remoteChainConfigAccount.Base.Remote.TokenAddress.Address, onChainEVMPoolConfig.TokenAddress.Address) || evmRemoteConfig.OverrideConfig {
			e.Logger.Infof("overriding remote config for chain %d", evmChainSelector)
			ixConfigure, err := solLockReleaseTokenPool.NewEditChainRemoteConfigInstruction(
				evmChainSelector,
				tokenPubKey,
				onChainEVMPoolConfig,
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
			// diff between [existing remote pool addresses on solana chain] vs [what was just derived from evm chain]
			diff := poolDiff(remoteChainConfigAccount.Base.Remote.PoolAddresses, onChainEVMPoolConfig.PoolAddresses)
			if len(diff) > 0 {
				e.Logger.Infof("adding new pool addresses for chain %d", evmChainSelector)
				ixAppend, err := solLockReleaseTokenPool.NewAppendRemotePoolAddressesInstruction(
					evmChainSelector,
					tokenPubKey,
					diff, // evm supports multiple remote pools per token
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
		}
	} else {
		ixns, err = getNewSetuptInstructionsForLockRelease(e, chain, solChainState, cfg, evmChainSelector, evmRemoteConfig.RateLimiterConfig, onChainEVMPoolConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to generate instructions: %w", err)
		}
	}

	return ixns, nil
}

// ADD TOKEN POOL LOOKUP TABLE
type TokenPoolLookupTableConfig struct {
	ChainSelector uint64
	TokenPubKey   solana.PublicKey
	PoolType      solTestTokenPool.PoolType
}

func (cfg TokenPoolLookupTableConfig) Validate(e deployment.Environment) error {
	if err := commonValidation(e, cfg.ChainSelector, cfg.TokenPubKey); err != nil {
		return err
	}
	return validatePoolDeployment(&e, cfg.PoolType, cfg.ChainSelector, cfg.TokenPubKey, false)
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
	tokenPubKey := cfg.TokenPubKey
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
	routerPoolSignerPDA, _, _ := solState.FindExternalTokenPoolsSignerPDA(tokenPool, routerProgramAddress)

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
		routerPoolSignerPDA,     // 9
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
	TokenPubKey     solana.PublicKey
	WritableIndexes []uint8
	MCMSSolana      *MCMSConfigSolana
}

func (cfg SetPoolConfig) Validate(e deployment.Environment) error {
	tokenPubKey := cfg.TokenPubKey
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
	tokenPubKey := cfg.TokenPubKey
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

// input only the ones you want to add
// onchain throws error when we pass already configured accounts
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
		ix, err = solLockReleaseTokenPool.NewRemoveFromAllowListInstruction(
			cfg.Accounts,
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
				cldf.SolDefaultCommitment)
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
		err = appendTxs(ixns, tokenPool, ccipChangeset.LockReleaseTokenPool, &txns)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate mcms txn: %w", err)
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
