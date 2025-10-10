package solana

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	solToken "github.com/gagliardetto/solana-go/programs/token"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solBaseTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/base_token_pool"
	solBurnMintTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/burnmint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/cctp_token_pool"
	solLockReleaseTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/lockrelease_token_pool"
	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/test_token_pool"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset_v1_5_1 "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// use this changeset to add a token pool and lookup table
var _ cldf.ChangeSet[AddTokenPoolAndLookupTableConfig] = AddTokenPoolAndLookupTable

// use this changeset to setup a token pool for a remote chain
var _ cldf.ChangeSet[SetupTokenPoolForRemoteChainConfig] = SetupTokenPoolForRemoteChain

// lock / release ops on LnR token pool
var _ cldf.ChangeSet[LockReleaseLiquidityOpsConfig] = LockReleaseLiquidityOps

// configure token pool allow list
var _ cldf.ChangeSet[ConfigureTokenPoolAllowListConfig] = ConfigureTokenPoolAllowList

// remove from token pool allow list
var _ cldf.ChangeSet[RemoveFromAllowListConfig] = RemoveFromTokenPoolAllowList

// token pool ops
var _ cldf.ChangeSet[TokenPoolOpsCfg] = TokenPoolOps

// sync supported CCTP domains
var _ cldf.ChangeSet[SyncDomainConfig] = SyncDomain

// extend token pool lookup table
var _ cldf.ChangeSet[ExtendTokenPoolLookupTableConfig] = ExtendTokenPoolLookupTable

// append mcms txns generated from solanainstructions
func appendTxs(instructions []solana.Instruction, tokenPool solana.PublicKey, poolType cldf.ContractType, txns *[]mcmsTypes.Transaction) error {
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
	// a pool pda is uniquely identified by (solTokenPubKey, poolType, metadata)
	PoolType                 cldf.ContractType
	TokenPubKey              solana.PublicKey
	Metadata                 string           // tag to identify which client/cll token pool executable to use
	CCTPTokenMessengerMinter solana.PublicKey // required if PoolType is CCTPTokenPool
	CCTPMessageTransmitter   solana.PublicKey // required if PoolType is CCTPTokenPool
}

type AddTokenPoolAndLookupTableConfig struct {
	ChainSelector    uint64
	TokenPoolConfigs []TokenPoolConfig
	MCMS             *proposalutils.TimelockConfig
}

type TokenPoolConfigWithMCM struct {
	ChainSelector    uint64
	TokenPoolConfigs []TokenPoolConfig
	MCMS             *proposalutils.TimelockConfig
}

func (cfg TokenPoolConfigWithMCM) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	for _, tokenPoolCfg := range cfg.TokenPoolConfigs {
		if err := chainState.CommonValidation(e, cfg.ChainSelector, tokenPoolCfg.TokenPubKey); err != nil {
			return err
		}

		if err := chainState.ValidatePoolDeployment(&e, tokenPoolCfg.PoolType, cfg.ChainSelector, tokenPoolCfg.TokenPubKey, false, tokenPoolCfg.Metadata); err != nil {
			return err
		}
	}
	return nil
}

func (cfg TokenPoolConfigWithMCM) ValidateForGlobalInit() error {
	if cfg.ChainSelector == 0 {
		return errors.New("chain selector must be defined")
	}
	for _, tokenPoolCfg := range cfg.TokenPoolConfigs {
		if tokenPoolCfg.PoolType == "" {
			return errors.New("pool type must be defined")
		}
		if tokenPoolCfg.Metadata == "" {
			return errors.New("metadata must be defined")
		}
	}

	return nil
}

type NewMintTokenPoolConfig struct {
	ChainSelector    uint64
	PoolType         cldf.ContractType
	TokenPubKey      solana.PublicKey
	Metadata         string
	MCMS             *proposalutils.TimelockConfig
	NewMintAuthority solana.PublicKey // new mint authority to set for the token pool
	OldMintAuthority solana.PublicKey // Only require when the current mint authority is a multisig
}

func (cfg NewMintTokenPoolConfig) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	if err := chainState.CommonValidation(e, cfg.ChainSelector, cfg.TokenPubKey); err != nil {
		return err
	}

	return chainState.ValidatePoolDeployment(&e, cfg.PoolType, cfg.ChainSelector, cfg.TokenPubKey, false, cfg.Metadata)
}

func (cfg AddTokenPoolAndLookupTableConfig) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	for _, tokenPoolCfg := range cfg.TokenPoolConfigs {
		if err := chainState.CommonValidation(e, cfg.ChainSelector, tokenPoolCfg.TokenPubKey); err != nil {
			return err
		}
		if tokenPoolCfg.PoolType == "" {
			return errors.New("pool type must be defined")
		}
		if tokenPoolCfg.PoolType == shared.CCTPTokenPool {
			if tokenPoolCfg.CCTPMessageTransmitter.IsZero() {
				return errors.New("cctp message transmitter is empty")
			}
			if tokenPoolCfg.CCTPTokenMessengerMinter.IsZero() {
				return errors.New("cctp messenger minter is empty")
			}
		}
		if err := chainState.ValidatePoolDeployment(&e, tokenPoolCfg.PoolType, cfg.ChainSelector, tokenPoolCfg.TokenPubKey, false, tokenPoolCfg.Metadata); err != nil {
			return err
		}
	}
	return nil
}

func AddTokenPoolAndLookupTable(e cldf.Environment, cfg AddTokenPoolAndLookupTableConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chainState := state.SolChains[cfg.ChainSelector]
	if err := cfg.Validate(e, chainState); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	addressBook := cldf.NewMemoryAddressBook()
	routerProgramAddress, _, _ := chainState.GetRouterInfo()
	rmnRemoteAddress := chainState.RMNRemote

	timelockSignerPDA, err := FetchTimelockSigner(e, cfg.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
	}

	var txns []mcmsTypes.Transaction
	for _, tokenPoolCfg := range cfg.TokenPoolConfigs {
		e.Logger.Infow("Adding token pool", "cfg", tokenPoolCfg)
		tokenPubKey := tokenPoolCfg.TokenPubKey
		tokenPool := chainState.GetActiveTokenPool(tokenPoolCfg.PoolType, tokenPoolCfg.Metadata)

		progDataAddr, err := deployment.GetProgramDataAddress(chain.Client, tokenPool)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get program data address for program %s: %w", tokenPool.String(), err)
		}
		authority, _, err := deployment.GetUpgradeAuthority(chain.Client, progDataAddr)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get upgrade authority for program data %s: %w", progDataAddr.String(), err)
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
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create associated token account for tokenpool (mint: %s, pool: %s): %w", tokenPubKey.String(), tokenPool.String(), err)
		}
		instructions := []solana.Instruction{createI}

		var configPDA solana.PublicKey
		// Global Configuration
		configPDA, err = TokenPoolGlobalConfigPDA(tokenPool)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get solana token pool global config PDA: %w", err)
		}

		// initialize token pool config pda
		var poolInitI solana.Instruction
		programData, err := getSolProgramData(e, chain, tokenPool)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get solana token pool program data: %w", err)
		}
		switch tokenPoolCfg.PoolType {
		case shared.BurnMintTokenPool:
			solBurnMintTokenPool.SetProgramID(tokenPool)
			// initialize token pool for token
			poolInitI, err = solBurnMintTokenPool.NewInitializeInstruction(
				poolConfigPDA,
				tokenPubKey,
				authority,
				solana.SystemProgramID,
				tokenPool,
				programData.Address,
				configPDA,
			).ValidateAndBuild()
		case shared.LockReleaseTokenPool:
			solLockReleaseTokenPool.SetProgramID(tokenPool)
			// initialize token pool for token
			poolInitI, err = solLockReleaseTokenPool.NewInitializeInstruction(
				poolConfigPDA,
				tokenPubKey,
				authority,
				solana.SystemProgramID,
				tokenPool,
				programData.Address,
				configPDA,
			).ValidateAndBuild()
		case shared.CCTPTokenPool:
			cctp_token_pool.SetProgramID(tokenPool)
			// initialize token pool for token
			poolInitI, err = cctp_token_pool.NewInitializeInstruction(
				routerProgramAddress,
				rmnRemoteAddress,
				poolConfigPDA,
				tokenPubKey,
				authority,
				solana.SystemProgramID,
				tokenPool,
				programData.Address,
				configPDA,
			).ValidateAndBuild()
		default:
			return cldf.ChangesetOutput{}, fmt.Errorf("invalid pool type: %s", tokenPoolCfg.PoolType)
		}
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}

		if authority.Equals(timelockSignerPDA) {
			tx, err := BuildMCMSTxn(poolInitI, tokenPool.String(), tokenPoolCfg.PoolType)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
			}
			txns = append(txns, *tx)
		} else {
			instructions = append(instructions, poolInitI)
		}

		// fetch current token mint authority
		var tokenMint solToken.Mint
		var mintAuthority string
		err = chain.GetAccountDataBorshInto(context.Background(), tokenPubKey, &tokenMint)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get token mint account data: %w", err)
		}
		if tokenMint.MintAuthority != nil {
			mintAuthority = tokenMint.MintAuthority.String()
		}

		// make pool mint_authority for token
		if tokenPoolCfg.PoolType == shared.BurnMintTokenPool && tokenPubKey != solana.SolMint {
			if mintAuthority == chain.DeployerKey.PublicKey().String() {
				authI, err := solTokenUtil.SetTokenMintAuthority(
					tokenprogramID,
					poolSigner,
					tokenPubKey,
					chain.DeployerKey.PublicKey(),
				)
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
				}
				instructions = append(instructions, authI)
				e.Logger.Infow("Setting mint authority", "poolSigner", poolSigner.String())
			} else {
				e.Logger.Warnw("Token's mint authority is not with deployer key, skipping setting poolSigner as mint authority",
					"poolType", tokenPoolCfg.PoolType, "mintAuthority", mintAuthority,
					"deployer", chain.DeployerKey.PublicKey().String(), "poolSigner", poolSigner.String())
			}
		} else {
			e.Logger.Warnw("PoolType is not a BurnMintTokenPool, skipping setting poolSigner as mint authority",
				"poolType", tokenPoolCfg.PoolType, "mintAuthority", mintAuthority,
				"deployer", chain.DeployerKey.PublicKey().String(), "poolSigner", poolSigner.String())
		}

		// confirm instructions
		if err := chain.Confirm(instructions); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
		}
		e.Logger.Infow("Created new token pool config", "token_pool_ata", tokenPoolATA.String(), "pool_config", poolConfigPDA.String(), "pool_signer", poolSigner.String())

		// add token pool lookup table
		csOutput, err := AddTokenPoolLookupTable(e, TokenPoolLookupTableConfig{
			ChainSelector:            cfg.ChainSelector,
			TokenPubKey:              tokenPoolCfg.TokenPubKey,
			PoolType:                 tokenPoolCfg.PoolType,
			Metadata:                 tokenPoolCfg.Metadata,
			CCTPTokenMessengerMinter: tokenPoolCfg.CCTPTokenMessengerMinter,
			CCTPMessageTransmitter:   tokenPoolCfg.CCTPMessageTransmitter,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to add token pool lookup table: %w", err)
		}
		err = addressBook.Merge(csOutput.AddressBook) //nolint:staticcheck // SA1019: AddressBook is deprecated, migration to DataStore pending
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book: %w", err)
		}
	}

	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to initialize token pools", cfg.MCMS.MinDelay, txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			AddressBook:           addressBook,
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{
		AddressBook: addressBook,
	}, nil
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
	TokenSymbol       shared.TokenSymbol
	PoolType          cldf.ContractType
	PoolVersion       semver.Version
	RateLimiterConfig RateLimiterConfig
	OverrideConfig    bool // set to true if you want to overwrite the remote pool config of a configured remote token pool
}

func (cfg EVMRemoteConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState, evmChainSelector uint64) error {
	// add evm family check
	if cfg.TokenSymbol == "" {
		return errors.New("token symbol must be defined")
	}
	err := cldf.IsValidChainSelector(evmChainSelector)
	if err != nil {
		return fmt.Errorf("failed to validate chain selector %d: %w", evmChainSelector, err)
	}
	chain, ok := e.BlockChains.EVMChains()[evmChainSelector]
	if !ok {
		return fmt.Errorf("chain with selector %d does not exist in environment", evmChainSelector)
	}
	chainState, ok := state.EVMChainState(evmChainSelector)
	if !ok {
		return fmt.Errorf("%s does not exist in state", chain.String())
	}
	// Ensure that the inputted type is known
	if _, typeOk := shared.TokenPoolTypes[cfg.PoolType]; !typeOk {
		return fmt.Errorf("%s is not a known token pool type", cfg.PoolType)
	}
	// Ensure that the inputted version is known
	if _, versionOk := shared.TokenPoolVersions[cfg.PoolVersion]; !versionOk {
		return fmt.Errorf("%s is not a known token pool version", cfg.PoolVersion)
	}
	// Ensure that a pool with given symbol, type and version is known to the environment
	_, getPoolOk := ccipChangeset_v1_5_1.GetTokenPoolAddressFromSymbolTypeAndVersion(chainState, chain, cfg.TokenSymbol, cfg.PoolType, cfg.PoolVersion)
	if !getPoolOk {
		return fmt.Errorf("token pool does not exist on %s with symbol %s, type %s, and version %s", chain.String(), cfg.TokenSymbol, cfg.PoolType, cfg.PoolVersion)
	}
	err = cfg.RateLimiterConfig.Validate()
	return err
}

type RemoteChainTokenPoolConfig struct {
	// a pool pda is uniquely identified by (solTokenPubKey, poolType, metadata)
	SolTokenPubKey   solana.PublicKey
	SolPoolType      cldf.ContractType
	Metadata         string // tag to identify which client/cll token pool executable to use
	EVMRemoteConfigs map[uint64]EVMRemoteConfig
}

func (cfg RemoteChainTokenPoolConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState, solChainSelector uint64, mcms *proposalutils.TimelockConfig) error {
	chainState := state.SolChains[solChainSelector]
	if err := chainState.CommonValidation(e, solChainSelector, cfg.SolTokenPubKey); err != nil {
		return err
	}
	chain := e.BlockChains.SolanaChains()[solChainSelector]
	if cfg.SolPoolType == "" {
		return errors.New("pool type must be defined")
	}

	if err := chainState.ValidatePoolDeployment(&e, cfg.SolPoolType, solChainSelector, cfg.SolTokenPubKey, true, cfg.Metadata); err != nil {
		return err
	}

	if err := ValidateMCMSConfigSolana(e, mcms, chain, chainState, cfg.SolTokenPubKey, cfg.Metadata, map[cldf.ContractType]bool{}); err != nil {
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

type SetupTokenPoolForRemoteChainConfig struct {
	SolChainSelector       uint64
	RemoteTokenPoolConfigs []RemoteChainTokenPoolConfig
	MCMS                   *proposalutils.TimelockConfig
}

func (cfg SetupTokenPoolForRemoteChainConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	for _, remoteTokenPoolConfig := range cfg.RemoteTokenPoolConfigs {
		if err := remoteTokenPoolConfig.Validate(e, state, cfg.SolChainSelector, cfg.MCMS); err != nil {
			return err
		}
	}
	return nil
}

func getOnChainEVMPoolConfig(e cldf.Environment, state stateview.CCIPOnChainState, evmChainSelector uint64, evmRemoteConfig EVMRemoteConfig) (solBaseTokenPool.RemoteConfig, error) {
	evmChain := e.BlockChains.EVMChains()[evmChainSelector]
	evmChainState := state.MustGetEVMChainState(evmChainSelector)
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

func InitGlobalConfigTokenPoolProgram(e cldf.Environment, cfg TokenPoolConfigWithMCM) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Setting up token pool global config", "cfg", cfg)

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	solChainState := state.SolChains[cfg.ChainSelector]
	if err := cfg.ValidateForGlobalInit(); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	chainState := state.SolChains[cfg.ChainSelector]
	routerProgramAddress, _, _ := chainState.GetRouterInfo()
	rmnRemoteAddress := chainState.RMNRemote
	timelockSignerPDA, err := FetchTimelockSigner(e, cfg.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
	}

	var txns []mcmsTypes.Transaction
	for _, tokenPoolCfg := range cfg.TokenPoolConfigs {
		tokenPool := solChainState.GetActiveTokenPool(tokenPoolCfg.PoolType, tokenPoolCfg.Metadata)

		useMcms := solanastateview.IsSolanaProgramOwnedByTimelock(
			&e,
			chain,
			solChainState,
			tokenPoolCfg.PoolType,
			tokenPoolCfg.TokenPubKey,
			tokenPoolCfg.Metadata,
		)
		var authority solana.PublicKey
		if useMcms {
			// If MCMS is used, the authority is the timelock signer PDA
			authority = timelockSignerPDA
		} else {
			// If MCMS is not used, the authority is the deployer key
			authority = chain.DeployerKey.PublicKey()
		}

		var configPDA solana.PublicKey
		// Global Configuration
		configPDA, err = TokenPoolGlobalConfigPDA(tokenPool)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get solana token pool global config PDA: %w", err)
		}

		// If configPDA already exists, we assume the global config is already initialized
		if _, err := chain.Client.GetAccountInfo(context.Background(), configPDA); err == nil {
			e.Logger.Infow("Global config already initialized", "configPDA", configPDA.String())
			return cldf.ChangesetOutput{}, nil
		}

		programData, err := getSolProgramData(e, chain, tokenPool)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get solana token pool program data: %w", err)
		}

		var initGlobalConfigIx solana.Instruction
		switch tokenPoolCfg.PoolType {
		case shared.BurnMintTokenPool:
			solBurnMintTokenPool.SetProgramID(tokenPool)
			initGlobalConfigIx, err = solBurnMintTokenPool.NewInitGlobalConfigInstruction(routerProgramAddress, rmnRemoteAddress, configPDA, authority, solana.SystemProgramID, tokenPool, programData.Address).ValidateAndBuild()
		case shared.LockReleaseTokenPool:
			solLockReleaseTokenPool.SetProgramID(tokenPool)
			initGlobalConfigIx, err = solLockReleaseTokenPool.NewInitGlobalConfigInstruction(routerProgramAddress, rmnRemoteAddress, configPDA, authority, solana.SystemProgramID, tokenPool, programData.Address).ValidateAndBuild()
		case shared.CCTPTokenPool:
			// CCTP token pool should not need separate global config initialization in normal use cases. Global config is initialized in the deployment changeset.
			cctp_token_pool.SetProgramID(tokenPool)
			initGlobalConfigIx, err = cctp_token_pool.NewInitGlobalConfigInstruction(configPDA, authority, solana.SystemProgramID, tokenPool, programData.Address).ValidateAndBuild()
		default:
			return cldf.ChangesetOutput{}, fmt.Errorf("invalid token pool type: %w", err)
		}
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build ix to init global config: %w", err)
		}

		instructions := []solana.Instruction{initGlobalConfigIx}

		if useMcms {
			err := appendTxs(instructions, tokenPool, tokenPoolCfg.PoolType, &txns)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate mcms txn: %w", err)
			}
		} else {
			if err := chain.Confirm(instructions); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}
	}

	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to init global config in Solana Token Pool", cfg.MCMS.MinDelay, txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}

func EnableSelfServedInTokenPoolProgram(e cldf.Environment, cfg TokenPoolConfigWithMCM) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Enable self served token pool", "cfg", cfg)

	return modifySelfServedConfig(e, cfg, true)
}

func DisableSelfServedInTokenPoolProgram(e cldf.Environment, cfg TokenPoolConfigWithMCM) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Disable self served token pool", "cfg", cfg)

	return modifySelfServedConfig(e, cfg, false)
}

func modifySelfServedConfig(e cldf.Environment, cfg TokenPoolConfigWithMCM, enabled bool) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	solChainState := state.SolChains[cfg.ChainSelector]
	if err := cfg.ValidateForGlobalInit(); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	var txns []mcmsTypes.Transaction
	for _, tokenPoolConfig := range cfg.TokenPoolConfigs {
		tokenPool := solChainState.GetActiveTokenPool(tokenPoolConfig.PoolType, tokenPoolConfig.Metadata)
		var configPDA solana.PublicKey
		// Global Configuration
		configPDA, err = TokenPoolGlobalConfigPDA(tokenPool)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get solana token pool global config PDA: %w", err)
		}

		useMcms := solanastateview.IsSolanaProgramOwnedByTimelock(
			&e,
			chain,
			solChainState,
			tokenPoolConfig.PoolType,
			tokenPoolConfig.TokenPubKey,
			tokenPoolConfig.Metadata,
		)

		// Checking that configPDA exists, so the update method will not fail
		if !useMcms {
			_, err = chain.Client.GetAccountInfo(context.Background(), configPDA)
			if err != nil {
				e.Logger.Infow("Global config not initialized", "configPDA", configPDA.String())
				return cldf.ChangesetOutput{}, nil
			}
		}

		programData, err := getSolProgramData(e, chain, tokenPool)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get solana token pool program data: %w", err)
		}

		timelockSignerPDA, err := FetchTimelockSigner(e, cfg.ChainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
		}
		var authority solana.PublicKey
		if useMcms {
			// If MCMS is used, the authority is the timelock signer PDA
			authority = timelockSignerPDA
		} else {
			// If MCMS is not used, the authority is the deployer key
			authority = chain.DeployerKey.PublicKey()
		}

		var initGlobalConfigIx solana.Instruction
		switch tokenPoolConfig.PoolType {
		case shared.BurnMintTokenPool:
			solBurnMintTokenPool.SetProgramID(tokenPool)
			initGlobalConfigIx, err = solBurnMintTokenPool.NewUpdateSelfServedAllowedInstruction(enabled, configPDA, authority, tokenPool, programData.Address).ValidateAndBuild()
		case shared.LockReleaseTokenPool:
			solLockReleaseTokenPool.SetProgramID(tokenPool)
			initGlobalConfigIx, err = solLockReleaseTokenPool.NewUpdateSelfServedAllowedInstruction(enabled, configPDA, authority, tokenPool, programData.Address).ValidateAndBuild()
		default:
			return cldf.ChangesetOutput{}, fmt.Errorf("invalid token pool type: %w", err)
		}
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build ix to init global config: %w", err)
		}

		instructions := []solana.Instruction{initGlobalConfigIx}

		if useMcms {
			err := appendTxs(instructions, tokenPool, tokenPoolConfig.PoolType, &txns)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate mcms txn: %w", err)
			}
		} else {
			if err := chain.Confirm(instructions); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}
	}

	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to init global config in Solana Token Pool", cfg.MCMS.MinDelay, txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}

func ModifyMintAuthority(e cldf.Environment, cfg NewMintTokenPoolConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Use multisig as mint authority", "cfg", cfg)

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	solChainState := state.SolChains[cfg.ChainSelector]
	if err := cfg.Validate(e, solChainState); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	tokenPubKey := cfg.TokenPubKey
	tokenProgram, _ := solChainState.TokenToTokenProgram(tokenPubKey)
	tokenPool := solChainState.GetActiveTokenPool(cfg.PoolType, cfg.Metadata)

	switch cfg.PoolType {
	case shared.BurnMintTokenPool:
		solBurnMintTokenPool.SetProgramID(tokenPool)
	case shared.LockReleaseTokenPool:
		return cldf.ChangesetOutput{}, nil
	default:
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid token pool type: %w", err)
	}

	newMintAuthority := cfg.NewMintAuthority
	tokenPoolSigner, _ := solTokenUtil.TokenPoolSignerAddress(tokenPubKey, tokenPool)

	poolConfig, err := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to calculate the pool config: %w", err)
	}
	programData, err := getSolProgramData(e, chain, tokenPool)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get solana token pool program data: %w", err)
	}

	useMcms := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		solChainState,
		shared.BurnMintTokenPool,
		tokenPubKey,
		cfg.Metadata,
	)

	var txns []mcmsTypes.Transaction
	if useMcms {
		timelockSigner, err := FetchTimelockSigner(e, cfg.ChainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
		}

		ix, err := solBurnMintTokenPool.NewTransferMintAuthorityToMultisigInstruction(
			poolConfig,
			tokenPubKey,
			tokenProgram,
			tokenPoolSigner,
			timelockSigner,
			newMintAuthority,
			tokenPool,
			programData.Address).ValidateAndBuild()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build ix to transfer mint authority to multisig: %w", err)
		}

		err = appendTxs([]solana.Instruction{ix}, tokenPool, cfg.PoolType, &txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate mcms txn: %w", err)
		}
	} else {
		builder := solBurnMintTokenPool.NewTransferMintAuthorityToMultisigInstruction(
			poolConfig,
			tokenPubKey,
			tokenProgram,
			tokenPoolSigner,
			chain.DeployerKey.PublicKey(),
			newMintAuthority,
			tokenPool,
			programData.Address)

		// Old mint authority is required only if the current mint authority is a multisig
		if (cfg.OldMintAuthority != solana.PublicKey{}) {
			builder.AccountMetaSlice = append(builder.AccountMetaSlice, solana.Meta(cfg.OldMintAuthority))
		}

		ix, err := builder.ValidateAndBuild()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build ix to transfer mint authority to multisig: %w", err)
		}

		if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
		}
	}

	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to init global config in Solana Token Pool", cfg.MCMS.MinDelay, txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}

func SetupTokenPoolForRemoteChain(e cldf.Environment, cfg SetupTokenPoolForRemoteChainConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Setting up token pool for remote chain", "cfg", cfg)
	envState, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	solChainState := envState.SolChains[cfg.SolChainSelector]
	if err := cfg.Validate(e, envState); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.SolChainSelector]
	var txns []mcmsTypes.Transaction
	for _, remoteTokenPoolConfig := range cfg.RemoteTokenPoolConfigs {
		tokenPubKey := remoteTokenPoolConfig.SolTokenPubKey
		tokenPool := solChainState.GetActiveTokenPool(remoteTokenPoolConfig.SolPoolType, remoteTokenPoolConfig.Metadata)
		useMcms := solanastateview.IsSolanaProgramOwnedByTimelock(
			&e,
			chain,
			solChainState,
			remoteTokenPoolConfig.SolPoolType,
			tokenPubKey,
			remoteTokenPoolConfig.Metadata,
		)
		switch remoteTokenPoolConfig.SolPoolType {
		case shared.BurnMintTokenPool:
			solBurnMintTokenPool.SetProgramID(tokenPool)
			for evmChainSelector, evmRemoteConfig := range remoteTokenPoolConfig.EVMRemoteConfigs {
				e.Logger.Infow("Setting up bnm token pool for remote chain", "remote_chain_selector", evmChainSelector, "token_pubkey", tokenPubKey.String(), "pool_type", remoteTokenPoolConfig.SolPoolType)
				chainIxs, err := getInstructionsForBurnMint(e, chain, envState, solChainState, remoteTokenPoolConfig, evmChainSelector, evmRemoteConfig)
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
				}
				if useMcms {
					err := appendTxs(chainIxs, tokenPool, remoteTokenPoolConfig.SolPoolType, &txns)
					if err != nil {
						return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate mcms txn: %w", err)
					}
				} else {
					if err := chain.Confirm(chainIxs); err != nil {
						return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
					}
				}
			}

		case shared.LockReleaseTokenPool:
			solLockReleaseTokenPool.SetProgramID(tokenPool)
			for evmChainSelector, evmRemoteConfig := range remoteTokenPoolConfig.EVMRemoteConfigs {
				e.Logger.Infow("Setting up lnr token pool for remote chain", "remote_chain_selector", evmChainSelector, "token_pubkey", tokenPubKey.String(), "pool_type", remoteTokenPoolConfig.SolPoolType)
				chainIxs, err := getInstructionsForLockRelease(e, chain, envState, solChainState, remoteTokenPoolConfig, evmChainSelector, evmRemoteConfig)
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
				}
				if useMcms {
					err := appendTxs(chainIxs, tokenPool, remoteTokenPoolConfig.SolPoolType, &txns)
					if err != nil {
						return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate mcms txn: %w", err)
					}
				} else {
					if err := chain.Confirm(chainIxs); err != nil {
						return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
					}
				}
			}
		case shared.CCTPTokenPool:
			cctp_token_pool.SetProgramID(tokenPool)
			for evmChainSelector, evmRemoteConfig := range remoteTokenPoolConfig.EVMRemoteConfigs {
				e.Logger.Infow("Setting up CCTP token pool for remote chain", "remote_chain_selector", evmChainSelector, "token_pubkey", tokenPubKey.String())
				chainIxs, err := getInstructionsForCCTP(e, chain, envState, solChainState, remoteTokenPoolConfig, evmChainSelector, evmRemoteConfig)
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
				}
				if useMcms {
					err := appendTxs(chainIxs, tokenPool, remoteTokenPoolConfig.SolPoolType, &txns)
					if err != nil {
						return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate mcms txn: %w", err)
					}
				} else {
					if err := chain.Confirm(chainIxs); err != nil {
						return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
					}
				}
			}
		default:
			return cldf.ChangesetOutput{}, fmt.Errorf("invalid pool type: %s", remoteTokenPoolConfig.SolPoolType)
		}
	}

	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.SolChainSelector, "proposal to edit token pools in Solana", cfg.MCMS.MinDelay, txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}

// checks if the evmChainSelector is supported for the given token and pool type
func isSupportedChain(chain cldf_solana.Chain, solTokenPubKey solana.PublicKey, solPoolAddress solana.PublicKey, poolType cldf.ContractType, evmChainSelector uint64) (bool, solBaseTokenPool.BaseChain, error) {
	// check if this remote chain is already configured for this token
	remoteChainConfigPDA, _, err := solTokenUtil.TokenPoolChainConfigPDA(evmChainSelector, solTokenPubKey, solPoolAddress)
	if err != nil {
		return false, solBaseTokenPool.BaseChain{}, fmt.Errorf("failed to get token pool remote chain config pda (remoteSelector: %d, mint: %s, pool: %s): %w", evmChainSelector, solTokenPubKey.String(), solPoolAddress.String(), err)
	}
	var base solBaseTokenPool.BaseChain
	switch poolType {
	case shared.BurnMintTokenPool, shared.LockReleaseTokenPool:
		var remoteChainConfigAccount solTestTokenPool.ChainConfig
		err = chain.GetAccountDataBorshInto(context.Background(), remoteChainConfigPDA, &remoteChainConfigAccount)
		if err != nil { // not a supported chain for this combination of token and pool type
			return false, solBaseTokenPool.BaseChain{}, nil
		}
		base = remoteChainConfigAccount.Base
	case shared.CCTPTokenPool:
		var remoteChainConfigAccount cctp_token_pool.ChainConfig
		err = chain.GetAccountDataBorshInto(context.Background(), remoteChainConfigPDA, &remoteChainConfigAccount)
		if err != nil { // not a supported chain for this combination of token and pool type
			return false, solBaseTokenPool.BaseChain{}, nil
		}
		base = remoteChainConfigAccount.Base
	}
	return true, base, nil
}

func getNewSetupInstructionsForBurnMint(
	e cldf.Environment,
	chain cldf_solana.Chain,
	chainState solanastateview.CCIPChainState,
	cfg RemoteChainTokenPoolConfig,
	evmChainSelector uint64,
	rateLimiterConfig RateLimiterConfig,
	onChainEVMPoolConfig solBaseTokenPool.RemoteConfig,
) ([]solana.Instruction, error) {
	e.Logger.Infow("getNewSetupInstructionsForBurnMint", "remote_chain_selector", evmChainSelector, "token_pubkey", cfg.SolTokenPubKey.String(), "pool_type", cfg.SolPoolType)
	tokenPool := chainState.GetActiveTokenPool(cfg.SolPoolType, cfg.Metadata)
	tokenPubKey := cfg.SolTokenPubKey
	poolConfigPDA, remoteChainConfigPDA := getPoolPDAs(tokenPubKey, tokenPool, evmChainSelector)
	ixns := make([]solana.Instruction, 0)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.SolPoolType,
		tokenPubKey,
		cfg.Metadata,
	)
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

	// there is a bug on the token pool contract which does not allow us to set the actual rate limits directly
	// we have to setup dummy limits first and then update it
	// so here we are checking if the user is trying to set actual limits,
	// and if so, we are setting dummy limits first and then updating it
	// if the user is just sending in enabled=false anyway, then no need to set dummy limits first
	if rateLimiterConfig.Inbound.Enabled || rateLimiterConfig.Outbound.Enabled {
		ixDummyRates, err := solBurnMintTokenPool.NewSetChainRateLimitInstruction(
			evmChainSelector,
			tokenPubKey,
			solBaseTokenPool.RateLimitConfig{
				Enabled:  false,
				Capacity: 0,
				Rate:     0,
			},
			solBaseTokenPool.RateLimitConfig{
				Enabled:  false,
				Capacity: 0,
				Rate:     0,
			},
			poolConfigPDA,
			remoteChainConfigPDA,
			authority,
		).ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("failed to generate instructions: %w", err)
		}
		ixns = append(ixns, ixDummyRates)
	}
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
	e cldf.Environment,
	chain cldf_solana.Chain,
	envState stateview.CCIPOnChainState,
	solChainState solanastateview.CCIPChainState,
	cfg RemoteChainTokenPoolConfig,
	evmChainSelector uint64,
	evmRemoteConfig EVMRemoteConfig,
) ([]solana.Instruction, error) {
	e.Logger.Infow("getInstructionsForBurnMint", "remote_chain_selector", evmChainSelector, "token_pubkey", cfg.SolTokenPubKey.String(), "pool_type", cfg.SolPoolType)
	tokenPubKey := cfg.SolTokenPubKey
	tokenPool := solChainState.GetActiveTokenPool(cfg.SolPoolType, cfg.Metadata)
	poolConfigPDA, remoteChainConfigPDA := getPoolPDAs(tokenPubKey, tokenPool, evmChainSelector)
	ixns := make([]solana.Instruction, 0)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		solChainState,
		cfg.SolPoolType,
		tokenPubKey,
		cfg.Metadata,
	)

	onChainEVMPoolConfig, err := getOnChainEVMPoolConfig(e, envState, evmChainSelector, evmRemoteConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get onchain evm config: %w", err)
	}

	isSupportedChain, baseChainConfigAccount, err := isSupportedChain(chain, tokenPubKey, tokenPool, shared.BurnMintTokenPool, evmChainSelector)
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
		if !bytes.Equal(baseChainConfigAccount.Remote.TokenAddress.Address, onChainEVMPoolConfig.TokenAddress.Address) || evmRemoteConfig.OverrideConfig {
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
			poolAddresses := baseChainConfigAccount.Remote.PoolAddresses
			// translate to base
			baseAddresses := make([]solBaseTokenPool.RemoteAddress, len(poolAddresses))
			for i, cfg := range poolAddresses {
				baseAddresses[i] = solBaseTokenPool.RemoteAddress{
					Address: cfg.Address,
				}
			}
			diff := poolDiff(baseAddresses, onChainEVMPoolConfig.PoolAddresses)
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
		ixns, err = getNewSetupInstructionsForBurnMint(e, chain, solChainState, cfg, evmChainSelector, evmRemoteConfig.RateLimiterConfig, onChainEVMPoolConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to generate instructions: %w", err)
		}
	}

	return ixns, nil
}

func getNewSetupInstructionsForLockRelease(
	e cldf.Environment,
	chain cldf_solana.Chain,
	chainState solanastateview.CCIPChainState,
	cfg RemoteChainTokenPoolConfig,
	evmChainSelector uint64,
	rateLimiterConfig RateLimiterConfig,
	onChainEVMPoolConfig solBaseTokenPool.RemoteConfig,
) ([]solana.Instruction, error) {
	e.Logger.Infow("getNewSetupInstructionsForLockRelease", "remote_chain_selector", evmChainSelector, "token_pubkey", cfg.SolTokenPubKey.String(), "pool_type", cfg.SolPoolType)
	tokenPubKey := cfg.SolTokenPubKey
	tokenPool := chainState.GetActiveTokenPool(cfg.SolPoolType, cfg.Metadata)
	poolConfigPDA, remoteChainConfigPDA := getPoolPDAs(tokenPubKey, tokenPool, evmChainSelector)
	ixns := make([]solana.Instruction, 0)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.SolPoolType,
		tokenPubKey,
		cfg.Metadata,
	)
	e.Logger.Infow("getNewSetupInstructionsForLockRelease", "authority", authority.String())
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

	if rateLimiterConfig.Inbound.Enabled || rateLimiterConfig.Outbound.Enabled {
		ixDummyRates, err := solLockReleaseTokenPool.NewSetChainRateLimitInstruction(
			evmChainSelector,
			tokenPubKey,
			solBaseTokenPool.RateLimitConfig{
				Enabled:  false,
				Capacity: 0,
				Rate:     0,
			},
			solBaseTokenPool.RateLimitConfig{
				Enabled:  false,
				Capacity: 0,
				Rate:     0,
			},
			poolConfigPDA,
			remoteChainConfigPDA,
			authority,
		).ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("failed to generate instructions: %w", err)
		}
		ixns = append(ixns, ixDummyRates)
	}
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
	e cldf.Environment,
	chain cldf_solana.Chain,
	envState stateview.CCIPOnChainState,
	solChainState solanastateview.CCIPChainState,
	cfg RemoteChainTokenPoolConfig,
	evmChainSelector uint64,
	evmRemoteConfig EVMRemoteConfig,
) ([]solana.Instruction, error) {
	e.Logger.Infow("getInstructionsForLockRelease", "remote_chain_selector", evmChainSelector, "token_pubkey", cfg.SolTokenPubKey.String(), "pool_type", cfg.SolPoolType)
	tokenPubKey := cfg.SolTokenPubKey
	tokenPool := solChainState.GetActiveTokenPool(cfg.SolPoolType, cfg.Metadata)
	poolConfigPDA, remoteChainConfigPDA := getPoolPDAs(tokenPubKey, tokenPool, evmChainSelector)
	ixns := make([]solana.Instruction, 0)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		solChainState,
		cfg.SolPoolType,
		tokenPubKey,
		cfg.Metadata,
	)

	onChainEVMPoolConfig, err := getOnChainEVMPoolConfig(e, envState, evmChainSelector, evmRemoteConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get on chain evm pool config: %w", err)
	}

	isSupportedChain, baseChainConfigAccount, err := isSupportedChain(chain, tokenPubKey, tokenPool, shared.LockReleaseTokenPool, evmChainSelector)
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
		if !bytes.Equal(baseChainConfigAccount.Remote.TokenAddress.Address, onChainEVMPoolConfig.TokenAddress.Address) || evmRemoteConfig.OverrideConfig {
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
			poolAddresses := baseChainConfigAccount.Remote.PoolAddresses
			// translate to base
			baseAddresses := make([]solBaseTokenPool.RemoteAddress, len(poolAddresses))
			for i, cfg := range poolAddresses {
				baseAddresses[i] = solBaseTokenPool.RemoteAddress{
					Address: cfg.Address,
				}
			}
			// diff between [existing remote pool addresses on solana chain] vs [what was just derived from evm chain]
			diff := poolDiff(baseAddresses, onChainEVMPoolConfig.PoolAddresses)
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
		ixns, err = getNewSetupInstructionsForLockRelease(e, chain, solChainState, cfg, evmChainSelector, evmRemoteConfig.RateLimiterConfig, onChainEVMPoolConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to generate instructions: %w", err)
		}
	}

	return ixns, nil
}

func getNewSetupInstructionsForCCTP(
	e cldf.Environment,
	chain cldf_solana.Chain,
	chainState solanastateview.CCIPChainState,
	cfg RemoteChainTokenPoolConfig,
	evmChainSelector uint64,
	rateLimiterConfig RateLimiterConfig,
	onChainEVMPoolConfig cctp_token_pool.RemoteConfig,
) ([]solana.Instruction, error) {
	e.Logger.Infow("getNewSetupInstructionsForCCTP", "remote_chain_selector", evmChainSelector, "token_pubkey", cfg.SolTokenPubKey.String())
	tokenPubKey := cfg.SolTokenPubKey
	tokenPool := chainState.CCTPTokenPool
	contractType := shared.CCTPTokenPool
	poolConfigPDA, remoteChainConfigPDA := getPoolPDAs(tokenPubKey, tokenPool, evmChainSelector)
	ixns := make([]solana.Instruction, 0)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		contractType,
		tokenPubKey,
		cfg.Metadata,
	)
	e.Logger.Infow("getNewSetupInstructionsForCCTP", "authority", authority.String())
	onChainEVMPoolConfigWithoutPoolAddress := cctp_token_pool.RemoteConfig{
		TokenAddress:  onChainEVMPoolConfig.TokenAddress,
		PoolAddresses: []cctp_token_pool.RemoteAddress{},
		Decimals:      onChainEVMPoolConfig.Decimals,
	}

	ixConfigure, err := cctp_token_pool.NewInitChainRemoteConfigInstruction(
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

	if rateLimiterConfig.Inbound.Enabled || rateLimiterConfig.Outbound.Enabled {
		ixDummyRates, err := cctp_token_pool.NewSetChainRateLimitInstruction(
			evmChainSelector,
			tokenPubKey,
			cctp_token_pool.RateLimitConfig{
				Enabled:  false,
				Capacity: 0,
				Rate:     0,
			},
			cctp_token_pool.RateLimitConfig{
				Enabled:  false,
				Capacity: 0,
				Rate:     0,
			},
			poolConfigPDA,
			remoteChainConfigPDA,
			authority,
		).ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("failed to generate instructions: %w", err)
		}
		ixns = append(ixns, ixDummyRates)
	}
	ixRates, err := cctp_token_pool.NewSetChainRateLimitInstruction(
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

	ixAppend, err := cctp_token_pool.NewAppendRemotePoolAddressesInstruction(
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

func getInstructionsForCCTP(
	e cldf.Environment,
	chain cldf_solana.Chain,
	envState stateview.CCIPOnChainState,
	solChainState solanastateview.CCIPChainState,
	cfg RemoteChainTokenPoolConfig,
	evmChainSelector uint64,
	evmRemoteConfig EVMRemoteConfig,
) ([]solana.Instruction, error) {
	e.Logger.Infow("getInstructionsForCCTP", "remote_chain_selector", evmChainSelector, "token_pubkey", cfg.SolTokenPubKey.String())
	tokenPubKey := cfg.SolTokenPubKey
	tokenPool := solChainState.CCTPTokenPool
	contractType := shared.CCTPTokenPool
	poolConfigPDA, remoteChainConfigPDA := getPoolPDAs(tokenPubKey, tokenPool, evmChainSelector)
	ixns := make([]solana.Instruction, 0)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		solChainState,
		contractType,
		tokenPubKey,
		cfg.Metadata,
	)

	onChainEVMPoolConfig, err := getOnChainEVMPoolConfig(e, envState, evmChainSelector, evmRemoteConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get on chain evm pool config: %w", err)
	}

	isSupportedChain, baseChainConfigAccount, err := isSupportedChain(chain, tokenPubKey, tokenPool, shared.CCTPTokenPool, evmChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to check if chain is supported: %w", err)
	}
	if isSupportedChain {
		// override the rate limits if the chain is already supported
		e.Logger.Infof("overriding rate limits for chain %d", evmChainSelector)
		ixRates, err := cctp_token_pool.NewSetChainRateLimitInstruction(
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
		if !bytes.Equal(baseChainConfigAccount.Remote.TokenAddress.Address, onChainEVMPoolConfig.TokenAddress.Address) || evmRemoteConfig.OverrideConfig {
			e.Logger.Infof("overriding remote config for chain %d", evmChainSelector)
			ixConfigure, err := cctp_token_pool.NewEditChainRemoteConfigInstruction(
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
			diff := poolDiff(baseChainConfigAccount.Remote.PoolAddresses, onChainEVMPoolConfig.PoolAddresses)
			if len(diff) > 0 {
				e.Logger.Infof("adding new pool addresses for chain %d", evmChainSelector)
				ixAppend, err := cctp_token_pool.NewAppendRemotePoolAddressesInstruction(
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
		ixns, err = getNewSetupInstructionsForCCTP(e, chain, solChainState, cfg, evmChainSelector, evmRemoteConfig.RateLimiterConfig, onChainEVMPoolConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to generate instructions: %w", err)
		}
	}

	return ixns, nil
}

// ADD TOKEN POOL LOOKUP TABLE
type TokenPoolLookupTableConfig struct {
	ChainSelector            uint64
	TokenPubKey              solana.PublicKey
	PoolType                 cldf.ContractType
	Metadata                 string
	CCTPTokenMessengerMinter solana.PublicKey
	CCTPMessageTransmitter   solana.PublicKey
}

func (cfg TokenPoolLookupTableConfig) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	if err := chainState.CommonValidation(e, cfg.ChainSelector, cfg.TokenPubKey); err != nil {
		return err
	}
	if cfg.PoolType == "" {
		return errors.New("pool type must be defined")
	}
	if cfg.Metadata == "" {
		return errors.New("metadata must be defined")
	}
	return chainState.ValidatePoolDeployment(&e, cfg.PoolType, cfg.ChainSelector, cfg.TokenPubKey, false, cfg.Metadata)
}

// this changeset is called in AddTokenPoolAndLookupTable
// call this indepently only for some very specific reason, otherwise this should not be called and
// AddTokenPoolAndLookupTable should be called instead
func AddTokenPoolLookupTable(e cldf.Environment, cfg TokenPoolLookupTableConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Adding token pool lookup table", "cfg", cfg)
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chainState := state.SolChains[cfg.ChainSelector]
	if err := cfg.Validate(e, chainState); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	ctx := e.GetContext()
	client := chain.Client

	authorityPrivKey := chain.DeployerKey // assuming the authority is the deployer key
	tokenPubKey := cfg.TokenPubKey
	tokenPool := chainState.GetActiveTokenPool(cfg.PoolType, cfg.Metadata)
	routerProgramAddress, _, _ := chainState.GetRouterInfo()
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
	tokenPoolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
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
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create lookup table for token pool (mint: %s): %w", tokenPubKey.String(), err)
	}
	typeVersion := cldf.NewTypeAndVersion(shared.TokenPoolLookupTable, deployment.Version1_0_0)
	typeVersion.Labels.Add(tokenPubKey.String())
	typeVersion.Labels.Add(cfg.Metadata)
	typeVersion.Labels.Add(cfg.PoolType.String())

	var list solana.PublicKeySlice
	switch cfg.PoolType {
	case shared.BurnMintTokenPool, shared.LockReleaseTokenPool:
		list = solana.PublicKeySlice{
			table,                 // 0
			tokenAdminRegistryPDA, // 1
			tokenPool,             // 2
			tokenPoolConfigPDA,    // 3 - writable
			poolTokenAccount,      // 4 - writable
			tokenPoolSigner,       // 5
			tokenProgram,          // 6
			tokenPubKey,           // 7 - writable
			feeTokenConfigPDA,     // 8
			routerPoolSignerPDA,   // 9
		}
	case shared.CCTPTokenPool:
		messageTransmitterPDA, _, err := solana.FindProgramAddress([][]byte{[]byte("message_transmitter")}, cfg.CCTPMessageTransmitter)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to calculate message transmitter: %w", err)
		}
		tokenMessenger, _, err := solana.FindProgramAddress([][]byte{[]byte("token_messenger")}, cfg.CCTPTokenMessengerMinter)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to calculate token messenger: %w", err)
		}
		tokenMinter, _, err := solana.FindProgramAddress([][]byte{[]byte("token_minter")}, cfg.CCTPTokenMessengerMinter)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to calculate token minter: %w", err)
		}
		minterLocalToken, _, err := solana.FindProgramAddress([][]byte{[]byte("local_token"), tokenPubKey.Bytes()}, cfg.CCTPTokenMessengerMinter)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to calculate minter local token: %w", err)
		}
		eventAuthority, _, err := solana.FindProgramAddress([][]byte{[]byte("__event_authority")}, cfg.CCTPTokenMessengerMinter)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to calculate token message minter authority: %w", err)
		}
		list = append(list,
			table,                 // 0
			tokenAdminRegistryPDA, // 1
			tokenPool,             // 2
			tokenPoolConfigPDA,    // 3
			poolTokenAccount,      // 4 - writable
			tokenPoolSigner,       // 5 - writable
			tokenProgram,          // 6
			tokenPubKey,           // 7 - writable
			feeTokenConfigPDA,     // 8
			routerPoolSignerPDA,   // 9
			// -- CCTP custom entries --
			messageTransmitterPDA,        // 10 - writable
			cfg.CCTPTokenMessengerMinter, // 11
			solana.SystemProgramID,       // 12
			cfg.CCTPMessageTransmitter,   // 13
			tokenMessenger,               // 14
			tokenMinter,                  // 15
			minterLocalToken,             // 16 - writable
			eventAuthority,               // 17
		)
	default:
		return cldf.ChangesetOutput{}, errors.New("unsupported pool type")
	}
	if err = solCommonUtil.ExtendLookupTable(ctx, client, table, *authorityPrivKey, list); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to extend lookup table for token pool (mint: %s): %w", tokenPubKey.String(), err)
	}
	if err := solCommonUtil.AwaitSlotChange(ctx, client); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to await slot change while extending lookup table: %w", err)
	}
	newAddressBook := cldf.NewMemoryAddressBook()
	tv := cldf.NewTypeAndVersion(shared.TokenPoolLookupTable, deployment.Version1_0_0)
	tv.Labels.Add(tokenPubKey.String())
	tv.Labels.Add(cfg.PoolType.String())
	tv.Labels.Add(cfg.Metadata)
	if err := newAddressBook.Save(cfg.ChainSelector, table.String(), tv); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save tokenpool address lookup table: %w", err)
	}
	e.Logger.Infow("Added token pool lookup table", "token_pubkey", tokenPubKey.String())

	return cldf.ChangesetOutput{
		AddressBook: newAddressBook,
	}, nil
}

// CONFIGURE TOKEN POOL ALLOW LIST
type ConfigureTokenPoolAllowListConfig struct {
	SolChainSelector uint64
	// a pool pda is uniquely identified by (solTokenPubKey, poolType, metadata)
	SolTokenPubKey string
	PoolType       cldf.ContractType
	Metadata       string // tag to identify which client/cll token pool executable to use
	// input only the ones you want to add, onchain throws error when we pass already configured accounts
	Accounts []solana.PublicKey
	Enabled  bool // enable or disable the allow list
	MCMS     *proposalutils.TimelockConfig
}

func (cfg ConfigureTokenPoolAllowListConfig) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)

	if cfg.PoolType == "" {
		return errors.New("pool type must be defined")
	}
	chain := e.BlockChains.SolanaChains()[cfg.SolChainSelector]
	if err := chainState.CommonValidation(e, cfg.SolChainSelector, tokenPubKey); err != nil {
		return err
	}
	if err := chainState.ValidatePoolDeployment(&e, cfg.PoolType, cfg.SolChainSelector, tokenPubKey, true, cfg.Metadata); err != nil {
		return err
	}
	return ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, tokenPubKey, cfg.Metadata, map[cldf.ContractType]bool{})
}

func ConfigureTokenPoolAllowList(e cldf.Environment, cfg ConfigureTokenPoolAllowListConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infof("Configuring token pool allowlist for token %s", cfg.SolTokenPubKey)
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chainState := state.SolChains[cfg.SolChainSelector]
	if err := cfg.Validate(e, chainState); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	chain := e.BlockChains.SolanaChains()[cfg.SolChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)

	var ix solana.Instruction
	tokenPool := chainState.GetActiveTokenPool(cfg.PoolType, cfg.Metadata)
	poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
	tokenPoolUsingMcms := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		cfg.PoolType,
		tokenPubKey,
		cfg.Metadata,
	)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.PoolType,
		tokenPubKey,
		cfg.Metadata,
	)
	switch cfg.PoolType {
	case shared.BurnMintTokenPool:
		poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
		solBurnMintTokenPool.SetProgramID(tokenPool)
		ix, err = solBurnMintTokenPool.NewConfigureAllowListInstruction(
			cfg.Accounts,
			cfg.Enabled,
			poolConfigPDA,
			tokenPubKey,
			authority,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
	case shared.LockReleaseTokenPool:
		poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
		solLockReleaseTokenPool.SetProgramID(tokenPool)
		ix, err = solLockReleaseTokenPool.NewConfigureAllowListInstruction(
			cfg.Accounts,
			cfg.Enabled,
			poolConfigPDA,
			tokenPubKey,
			authority,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
	case shared.CCTPTokenPool:
		cctp_token_pool.SetProgramID(tokenPool)
		ix, err = cctp_token_pool.NewConfigureAllowListInstruction(
			cfg.Accounts,
			cfg.Enabled,
			poolConfigPDA,
			tokenPubKey,
			authority,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
	default:
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid pool type: %s", cfg.PoolType)
	}
	if tokenPoolUsingMcms {
		tx, err := BuildMCMSTxn(ix, tokenPool.String(), cfg.PoolType)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.SolChainSelector, "proposal to ConfigureTokenPoolAllowList in Solana", cfg.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	e.Logger.Infow("Configured token pool allowlist", "token_pubkey", tokenPubKey.String())
	return cldf.ChangesetOutput{}, nil
}

// REMOVE FROM TOKEN POOL ALLOW LIST
type RemoveFromAllowListConfig struct {
	SolChainSelector uint64
	// a pool pda is uniquely identified by (solTokenPubKey, poolType, metadata)
	SolTokenPubKey string
	PoolType       cldf.ContractType
	Metadata       string             // tag to identify which client/cll token pool executable to use
	Accounts       []solana.PublicKey // accounts to remove from allow list
	MCMS           *proposalutils.TimelockConfig
}

func (cfg RemoveFromAllowListConfig) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)
	if err := chainState.CommonValidation(e, cfg.SolChainSelector, tokenPubKey); err != nil {
		return err
	}
	if cfg.PoolType == "" {
		return errors.New("pool type must be defined")
	}
	chain := e.BlockChains.SolanaChains()[cfg.SolChainSelector]
	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, tokenPubKey, cfg.Metadata, map[cldf.ContractType]bool{}); err != nil {
		return err
	}
	return chainState.ValidatePoolDeployment(&e, cfg.PoolType, cfg.SolChainSelector, tokenPubKey, true, cfg.Metadata)
}

func RemoveFromTokenPoolAllowList(e cldf.Environment, cfg RemoveFromAllowListConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infof("Removing from token pool allowlist for token %s", cfg.SolTokenPubKey)
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chainState := state.SolChains[cfg.SolChainSelector]
	if err := cfg.Validate(e, chainState); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.SolChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)
	tokenPool := chainState.GetActiveTokenPool(cfg.PoolType, cfg.Metadata)
	poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)

	var ix solana.Instruction
	tokenPoolUsingMcms := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		cfg.PoolType,
		tokenPubKey,
		cfg.Metadata,
	)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.PoolType,
		tokenPubKey,
		cfg.Metadata,
	)
	switch cfg.PoolType {
	case shared.BurnMintTokenPool:
		poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
		solBurnMintTokenPool.SetProgramID(tokenPool)
		ix, err = solBurnMintTokenPool.NewRemoveFromAllowListInstruction(
			cfg.Accounts,
			poolConfigPDA,
			tokenPubKey,
			authority,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
	case shared.LockReleaseTokenPool:
		poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
		solLockReleaseTokenPool.SetProgramID(tokenPool)
		ix, err = solLockReleaseTokenPool.NewRemoveFromAllowListInstruction(
			cfg.Accounts,
			poolConfigPDA,
			tokenPubKey,
			authority,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
	case shared.CCTPTokenPool:
		cctp_token_pool.SetProgramID(tokenPool)
		ix, err = cctp_token_pool.NewRemoveFromAllowListInstruction(
			cfg.Accounts,
			poolConfigPDA,
			tokenPubKey,
			authority,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
	default:
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid pool type: %s", cfg.PoolType)
	}
	if tokenPoolUsingMcms {
		tx, err := BuildMCMSTxn(ix, tokenPool.String(), cfg.PoolType)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.SolChainSelector, "proposal to RemoveFromTokenPoolAllowList in Solana", cfg.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	e.Logger.Infow("Configured token pool allowlist", "token_pubkey", tokenPubKey.String())
	return cldf.ChangesetOutput{}, nil
}

// LOCK/UNLOCK LIQUIDITY
type LockReleaseLiquidityOpsConfig struct {
	SolChainSelector uint64
	// a pool pda is uniquely identified by (solTokenPubKey, poolType, metadata)
	// poolType is only LockAndRelease_PoolType for this migration
	SolTokenPubKey string
	SetCfg         *SetLiquidityConfig
	LiquidityCfg   *LiquidityConfig
	RebalancerCfg  *RebalancerConfig
	MCMS           *proposalutils.TimelockConfig
	Metadata       string
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

func (cfg LockReleaseLiquidityOpsConfig) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)
	if err := chainState.CommonValidation(e, cfg.SolChainSelector, tokenPubKey); err != nil {
		return err
	}
	chain := e.BlockChains.SolanaChains()[cfg.SolChainSelector]
	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, tokenPubKey, cfg.Metadata, map[cldf.ContractType]bool{}); err != nil {
		return err
	}
	return chainState.ValidatePoolDeployment(&e, shared.LockReleaseTokenPool, cfg.SolChainSelector, tokenPubKey, true, cfg.Metadata)
}

func LockReleaseLiquidityOps(e cldf.Environment, cfg LockReleaseLiquidityOpsConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infof("Locking/Unlocking liquidity for token %s", cfg.SolTokenPubKey)
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chainState := state.SolChains[cfg.SolChainSelector]
	if err := cfg.Validate(e, chainState); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.SolChainSelector]
	tokenPool := chainState.GetActiveTokenPool(shared.LockReleaseTokenPool, cfg.Metadata)
	solLockReleaseTokenPool.SetProgramID(tokenPool)
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)
	poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
	tokenPoolUsingMcms := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.LockReleaseTokenPool,
		tokenPubKey,
		cfg.Metadata,
	)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		shared.LockReleaseTokenPool,
		tokenPubKey,
		cfg.Metadata,
	)
	ixns := make([]solana.Instruction, 0)
	if cfg.SetCfg != nil {
		ix, err := solLockReleaseTokenPool.NewSetCanAcceptLiquidityInstruction(
			cfg.SetCfg.Enabled,
			poolConfigPDA,
			tokenPubKey,
			authority,
		).ValidateAndBuild()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
		ixns = append(ixns, ix)
	}
	if cfg.LiquidityCfg != nil {
		tokenProgram, _ := chainState.TokenToTokenProgram(tokenPubKey)
		poolSigner, _ := solTokenUtil.TokenPoolSignerAddress(tokenPubKey, tokenPool)
		poolConfigAccount := solLockReleaseTokenPool.State{}
		_ = chain.GetAccountDataBorshInto(context.Background(), poolConfigPDA, &poolConfigAccount)
		if cfg.LiquidityCfg.Amount <= 0 {
			return cldf.ChangesetOutput{}, fmt.Errorf("invalid amount: %d", cfg.LiquidityCfg.Amount)
		}
		tokenAmount := uint64(cfg.LiquidityCfg.Amount) // #nosec G115 - we check the amount above
		switch cfg.LiquidityCfg.Type {
		case Provide:
			outDec, outVal, err := solTokenUtil.TokenBalance(
				e.GetContext(),
				chain.Client,
				cfg.LiquidityCfg.RemoteTokenAccount,
				cldf_solana.SolDefaultCommitment)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get token balance: %w", err)
			}
			if outVal < cfg.LiquidityCfg.Amount {
				return cldf.ChangesetOutput{}, fmt.Errorf("insufficient token balance: %d < %d", outVal, cfg.LiquidityCfg.Amount)
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
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to TokenApproveChecked: %w", err)
			}
			if err = chain.Confirm([]solana.Instruction{ix1}); err != nil {
				e.Logger.Errorw("Failed to confirm instructions for TokenApproveChecked", "chain", chain.String(), "err", err)
				return cldf.ChangesetOutput{}, err
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
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
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
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
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
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
		ixns = append(ixns, ix)
	}

	if tokenPoolUsingMcms {
		txns := make([]mcmsTypes.Transaction, 0)
		err := appendTxs(ixns, tokenPool, shared.LockReleaseTokenPool, &txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate mcms txn: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.SolChainSelector, "proposal to RemoveFromTokenPoolAllowList in Solana", cfg.MCMS.MinDelay, txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	err = chain.Confirm(ixns)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return cldf.ChangesetOutput{}, nil
}

// TOKEN POOL OPS
type TokenPoolOpsCfg struct {
	SolChainSelector uint64
	// a pool pda is uniquely identified by (solTokenPubKey, poolType, metadata)
	SolTokenPubKey string
	PoolType       cldf.ContractType
	Metadata       string          // tag to identify which client/cll token pool executable to use
	DeleteChainCfg *DeleteChainCfg // remove remote pool config corresponding to the set (solTokenPubKey, poolType, metadata, remoteChainSelector)
	SetRouterCfg   *SetRouterCfg   // set router address on token pool config pda
	MCMS           *proposalutils.TimelockConfig
}

type DeleteChainCfg struct {
	RemoteChainSelector uint64
}

type SetRouterCfg struct {
	Router solana.PublicKey
}

func (cfg TokenPoolOpsCfg) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)
	if cfg.PoolType == "" {
		return errors.New("pool type must be defined")
	}
	chainState := state.SolChains[cfg.SolChainSelector]
	chain := e.BlockChains.SolanaChains()[cfg.SolChainSelector]
	if err := chainState.CommonValidation(e, cfg.SolChainSelector, tokenPubKey); err != nil {
		return err
	}
	if err := chainState.ValidatePoolDeployment(&e, cfg.PoolType, cfg.SolChainSelector, tokenPubKey, true, cfg.Metadata); err != nil {
		return err
	}
	if cfg.DeleteChainCfg != nil {
		var remoteChainConfigAccount any

		tokenPool := chainState.GetActiveTokenPool(cfg.PoolType, cfg.Metadata)
		switch cfg.PoolType {
		case shared.BurnMintTokenPool:
			remoteChainConfigAccount = solBurnMintTokenPool.ChainConfig{}
		case shared.LockReleaseTokenPool:
			remoteChainConfigAccount = solLockReleaseTokenPool.ChainConfig{}
		case shared.CCTPTokenPool:
			remoteChainConfigAccount = cctp_token_pool.ChainConfig{}
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
	return ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, tokenPubKey, cfg.Metadata, map[cldf.ContractType]bool{})
}

// remove remote pool config corresponding to the set (solTokenPubKey, poolType, metadata, remoteChainSelector)
// set router address on token pool config pda
func TokenPoolOps(e cldf.Environment, cfg TokenPoolOpsCfg) (cldf.ChangesetOutput, error) {
	e.Logger.Infof("Setting pool config for token %s", cfg.SolTokenPubKey)
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.SolChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.SolTokenPubKey)
	chainState := state.SolChains[cfg.SolChainSelector]
	var ix solana.Instruction
	ixns := make([]solana.Instruction, 0)
	tokenPool := chainState.GetActiveTokenPool(cfg.PoolType, cfg.Metadata)
	tokenPoolUsingMcms := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		cfg.PoolType,
		tokenPubKey,
		cfg.Metadata,
	)

	programData, err := getSolProgramData(e, chain, tokenPool)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get solana token pool program data: %w", err)
	}
	poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
	remoteChainConfigPDA, _, _ := solTokenUtil.TokenPoolChainConfigPDA(cfg.DeleteChainCfg.RemoteChainSelector, tokenPubKey, tokenPool)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.PoolType,
		tokenPubKey,
		cfg.Metadata,
	)

	switch cfg.PoolType {
	case shared.BurnMintTokenPool:
		solBurnMintTokenPool.SetProgramID(tokenPool)
		if cfg.DeleteChainCfg != nil {
			ix, err = solBurnMintTokenPool.NewDeleteChainConfigInstruction(
				cfg.DeleteChainCfg.RemoteChainSelector,
				tokenPubKey,
				poolConfigPDA,
				remoteChainConfigPDA,
				authority,
			).ValidateAndBuild()
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
			ixns = append(ixns, ix)
		}
		if cfg.SetRouterCfg != nil {
			ix, err = solBurnMintTokenPool.NewSetRouterInstruction(
				cfg.SetRouterCfg.Router,
				poolConfigPDA,
				tokenPubKey,
				authority,
				tokenPool,
				programData.Address,
			).ValidateAndBuild()
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
			ixns = append(ixns, ix)
		}
	case shared.LockReleaseTokenPool:
		solLockReleaseTokenPool.SetProgramID(tokenPool)
		if cfg.DeleteChainCfg != nil {
			ix, err = solLockReleaseTokenPool.NewDeleteChainConfigInstruction(
				cfg.DeleteChainCfg.RemoteChainSelector,
				tokenPubKey,
				poolConfigPDA,
				remoteChainConfigPDA,
				authority,
			).ValidateAndBuild()
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
			ixns = append(ixns, ix)
		}
		if cfg.SetRouterCfg != nil {
			ix, err = solLockReleaseTokenPool.NewSetRouterInstruction(
				cfg.SetRouterCfg.Router,
				poolConfigPDA,
				tokenPubKey,
				authority,
				tokenPool,
				programData.Address,
			).ValidateAndBuild()
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
			ixns = append(ixns, ix)
		}
	case shared.CCTPTokenPool:
		cctp_token_pool.SetProgramID(tokenPool)
		if cfg.DeleteChainCfg != nil {
			ix, err = cctp_token_pool.NewDeleteChainConfigInstruction(
				cfg.DeleteChainCfg.RemoteChainSelector,
				tokenPubKey,
				poolConfigPDA,
				remoteChainConfigPDA,
				authority,
			).ValidateAndBuild()
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
			ixns = append(ixns, ix)
		}
		if cfg.SetRouterCfg != nil {
			ix, err = cctp_token_pool.NewSetRouterInstruction(
				cfg.SetRouterCfg.Router,
				poolConfigPDA,
				tokenPubKey,
				authority,
				tokenPool,
				programData.Address,
			).ValidateAndBuild()
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
			ixns = append(ixns, ix)
		}
	default:
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid pool type: %s", cfg.PoolType)
	}
	if tokenPoolUsingMcms {
		tx, err := BuildMCMSTxn(ix, tokenPool.String(), cfg.PoolType)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.SolChainSelector, "proposal to ConfigureTokenPoolAllowList in Solana", cfg.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	if err := chain.Confirm(ixns); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	e.Logger.Infow("Configured token pool allowlist", "token_pubkey", tokenPubKey.String())
	return cldf.ChangesetOutput{}, nil
}

func InitializeStateVersion(e cldf.Environment, cfg TokenPoolConfigWithMCM) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Init state version for old tp", "cfg", cfg)

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	solChainState := state.SolChains[cfg.ChainSelector]
	if err := cfg.Validate(e, solChainState); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	var txns []mcmsTypes.Transaction
	for _, tokenPoolConfig := range cfg.TokenPoolConfigs {
		tokenPubKey := tokenPoolConfig.TokenPubKey
		tokenPool := solChainState.GetActiveTokenPool(tokenPoolConfig.PoolType, tokenPoolConfig.Metadata)
		poolConfig, err := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to calculate the pool configg: %w", err)
		}

		// This operation is permisionless, so we don't need to check ownership
		var initializeStateVersionIx solana.Instruction
		switch tokenPoolConfig.PoolType {
		case shared.BurnMintTokenPool:
			solBurnMintTokenPool.SetProgramID(tokenPool)
			initializeStateVersionIx, err = solBurnMintTokenPool.NewInitializeStateVersionInstruction(
				tokenPubKey,
				poolConfig).ValidateAndBuild()
		case shared.LockReleaseTokenPool:
			solLockReleaseTokenPool.SetProgramID(tokenPool)
			initializeStateVersionIx, err = solLockReleaseTokenPool.NewInitializeStateVersionInstruction(
				tokenPubKey,
				poolConfig).ValidateAndBuild()
		default:
			return cldf.ChangesetOutput{}, fmt.Errorf("invalid token pool type: %w", err)
		}
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build ix to init global config: %w", err)
		}

		useMcms := solanastateview.IsSolanaProgramOwnedByTimelock(
			&e,
			chain,
			solChainState,
			tokenPoolConfig.PoolType,
			tokenPubKey,
			tokenPoolConfig.Metadata,
		)

		instructions := []solana.Instruction{initializeStateVersionIx}

		if useMcms {
			err := appendTxs(instructions, tokenPool, tokenPoolConfig.PoolType, &txns)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate mcms txn: %w", err)
			}
		} else {
			if err := chain.Confirm(instructions); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}
	}

	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to init global config in Solana Token Pool", cfg.MCMS.MinDelay, txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}

func TokenPoolGlobalConfigPDA(programID solana.PublicKey) (solana.PublicKey, error) {
	addr, _, err := solana.FindProgramAddress([][]byte{[]byte("config")}, programID)
	return addr, err
}

type SyncDomainConfig struct {
	ChainSelector uint64
	// cctpChainConfigMap maps chain selectors to their associated CctpChainConfig
	CCTPChainConfigMap map[uint64]CctpChainConfig
	MCMS               *proposalutils.TimelockConfig
}

type CctpChainConfig struct {
	Domain            uint32
	DestinationCaller solana.PublicKey
}

func (cfg SyncDomainConfig) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	// Validate map contains configs
	if len(cfg.CCTPChainConfigMap) == 0 {
		return errors.New("CCTP chain config map is empty")
	}
	// Validate USDC pool exists in state
	if chainState.CCTPTokenPool.IsZero() {
		return errors.New("CCTP token pool does not exist in state")
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := solanastateview.ValidateOwnershipSolana(&e, chain, cfg.MCMS != nil, chainState.CCTPTokenPool, shared.CCTPTokenPool, chainState.USDCToken); err != nil {
		return fmt.Errorf("failed to validate ownership for cctp token pool: %w", err)
	}
	// Validate chain configs are initialized for each chain selector
	for chainSel := range cfg.CCTPChainConfigMap {
		supported, _, err := isSupportedChain(chain, chainState.USDCToken, chainState.CCTPTokenPool, shared.CCTPTokenPool, chainSel)
		if err != nil {
			return fmt.Errorf("failed to validate if remote chain %d is supported: %w", chainSel, err)
		}
		if !supported {
			return fmt.Errorf("chain config not initialized for selector %d", chainSel)
		}
	}
	return nil
}

// SyncDomain adds or removes CCTP domain configs from the Solana CCTP token pool
func SyncDomain(e cldf.Environment, cfg SyncDomainConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Syncing USDC domains", "cfg", cfg)
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chainState := state.SolChains[cfg.ChainSelector]
	if err := cfg.Validate(e, chainState); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]

	cctpTokenPool := chainState.CCTPTokenPool
	usdcToken := chainState.USDCToken

	useMcms := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.CCTPTokenPool,
		usdcToken,
		shared.CLLMetadata,
	)
	timelockSignerPDA, err := FetchTimelockSigner(e, cfg.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
	}
	var authority solana.PublicKey
	if useMcms {
		// If MCMS is used, the authority is the timelock signer PDA
		authority = timelockSignerPDA
	} else {
		// If MCMS is not used, the authority is the deployer key
		authority = chain.DeployerKey.PublicKey()
	}

	statePDA, err := solTokenUtil.TokenPoolConfigAddress(usdcToken, cctpTokenPool)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to calculate token pool config PDA: %w", err)
	}

	var txns []mcmsTypes.Transaction
	cctp_token_pool.SetProgramID(cctpTokenPool)
	for remoteChainSel, cctpConfig := range cfg.CCTPChainConfigMap {
		e.Logger.Infow("Setting up USDC token pool CCTP config for remote chain", "remote_chain_selector", remoteChainSel, "domain", cctpConfig.Domain, "destination_caller", cctpConfig.DestinationCaller.String())

		chainConfigPDA, _, err := solTokenUtil.TokenPoolChainConfigPDA(remoteChainSel, usdcToken, cctpTokenPool)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to calculate token pool config PDA: %w", err)
		}

		ix, err := cctp_token_pool.NewEditChainRemoteConfigCctpInstruction(
			remoteChainSel,
			usdcToken,
			cctp_token_pool.CctpChain{
				DomainId:          cctpConfig.Domain,
				DestinationCaller: cctpConfig.DestinationCaller,
			},
			statePDA,
			chainConfigPDA,
			authority,
		).ValidateAndBuild()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
		if useMcms {
			err := appendTxs([]solana.Instruction{ix}, cctpTokenPool, shared.CCTPTokenPool, &txns)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate mcms txn: %w", err)
			}
		} else {
			if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}
	}

	if len(txns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to edit USDC token pool CCTP config in Solana", cfg.MCMS.MinDelay, txns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}

type ExtendTokenPoolLookupTableConfig struct {
	// NOTE: the default behavior would be to block duplicates, but if you need to run the changeset without this restriction, then you can set this to true
	SkipValidationsForDuplicates bool
	ChainSelector                uint64
	TokenPubKey                  solana.PublicKey
	PoolType                     cldf.ContractType
	Metadata                     string
	Accounts                     solana.PublicKeySlice
}

func (cfg ExtendTokenPoolLookupTableConfig) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	if len(cfg.Accounts) == 0 {
		return nil
	}

	if cfg.TokenPubKey.IsZero() {
		return errors.New("required field 'TokenPubKey' is empty or the zero address")
	}

	if cfg.PoolType == "" {
		return errors.New("required field 'PoolType' is empty or not provided")
	}

	if cfg.Metadata == "" {
		return errors.New("required field 'Metadata' is empty or not provided")
	}

	if !cfg.SkipValidationsForDuplicates {
		acctSet := map[string]bool{}
		for _, acct := range cfg.Accounts {
			key := acct.String()
			if _, exists := acctSet[key]; exists {
				return fmt.Errorf("field 'Accounts' has 1 or more duplicate public keys: %s", key)
			}
			acctSet[key] = true
		}
	}

	solChain, exist := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if !exist {
		return fmt.Errorf("chain selector '%d' was not found in environment", cfg.ChainSelector)
	}

	lutPublKey, err := chainState.GetTokenPoolLookupTableAddress(cfg.TokenPubKey, cfg.PoolType, cfg.Metadata)
	if err != nil {
		return fmt.Errorf("failed to get lookup table address: %w", err)
	}

	lutEntries, err := solCommonUtil.GetAddressLookupTable(e.GetContext(), solChain.Client, lutPublKey)
	if err != nil {
		return fmt.Errorf("failed to get entries of address lookup table at address '%s': %w", lutPublKey.String(), err)
	}

	if !cfg.SkipValidationsForDuplicates {
		alutSet := map[string]bool{}
		for _, entry := range lutEntries {
			alutSet[entry.String()] = true
		}

		duplSet := map[string]bool{}
		for _, acct := range cfg.Accounts {
			key := acct.String()
			if _, exists := alutSet[key]; exists {
				duplSet[key] = true
			}
		}

		if len(duplSet) > 0 {
			return fmt.Errorf(
				"refusing to extend lookup table at address '%s' - one or more input accounts overlap with existing LUT entries: [ %s ]",
				lutPublKey.String(),
				strings.Join(slices.AppendSeq([]string{}, maps.Keys(duplSet)), ", "),
			)
		}
	}

	entries := solana.PublicKeySlice{}
	entries.Append(lutEntries...)
	e.Logger.Infof(
		"Validations complete - lookup table at address '%s' has %d accounts and currently looks like this: [ %s ]",
		lutPublKey.String(),
		len(entries),
		strings.Join(entries.ToBase58(), ", "),
	)

	return nil
}

func ExtendTokenPoolLookupTable(e cldf.Environment, cfg ExtendTokenPoolLookupTableConfig) (cldf.ChangesetOutput, error) {
	if len(cfg.Accounts) == 0 {
		e.Logger.Warn("no accounts were provided - exiting early as there is nothing to do")
		return cldf.ChangesetOutput{}, nil
	}

	chainState, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	solState, exist := chainState.SolChains[cfg.ChainSelector]
	if !exist {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain selector '%d' was not found in state", cfg.ChainSelector)
	}

	solChain, exist := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if !exist {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain selector '%d' was not found in environment", cfg.ChainSelector)
	}

	lutPublKey, err := solState.GetTokenPoolLookupTableAddress(cfg.TokenPubKey, cfg.PoolType, cfg.Metadata)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get lookup table address: %w", err)
	}

	if solChain.DeployerKey == nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("solana deployer key is nil for selector '%d'", cfg.ChainSelector)
	}

	if err := solChain.DeployerKey.Validate(); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("solana deployer key is invalid for selector '%d': %w", cfg.ChainSelector, err)
	}

	e.Logger.Infof(
		"Extending lookup table at address '%s' with %d accounts: [ %s ]",
		lutPublKey.String(),
		cfg.Accounts.Len(),
		strings.Join(cfg.Accounts.ToBase58(), ", "),
	)

	ctx := e.GetContext()
	if err := solCommonUtil.ExtendLookupTable(
		ctx,
		solChain.Client,
		lutPublKey,
		*solChain.DeployerKey,
		cfg.Accounts,
	); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to extend lookup table: %w", err)
	}

	// Displays the new LUT or fallbacks to displaying the newly added accounts if an error occurs
	if lutEntries, err := solCommonUtil.GetAddressLookupTable(ctx, solChain.Client, lutPublKey); err != nil {
		e.Logger.Warnf("could not display full LUT due to error - showing only newly added entries instead: %v", err)
		e.Logger.Infof(
			"Lookup table at address '%s' was successfully extended with %d accounts: [ %s ]",
			lutPublKey.String(),
			cfg.Accounts.Len(),
			strings.Join(cfg.Accounts.ToBase58(), ", "),
		)
	} else {
		accounts := solana.PublicKeySlice{}
		accounts.Append(lutEntries...)
		e.Logger.Infof(
			"Lookup table at address '%s' was successfully extended with %d accounts - the LUT looks like this now: [ %s ]",
			lutPublKey.String(),
			cfg.Accounts.Len(),
			strings.Join(accounts.ToBase58(), ", "),
		)
	}

	return cldf.ChangesetOutput{}, nil
}

type RateLimitAdminConfig struct {
	SolTokenPubKey    string
	PoolType          cldf.ContractType
	Metadata          string
	NewRateLimitAdmin solana.PublicKey
}

type SetRateLimitAdminConfig struct {
	SolChainSelector      uint64
	RateLimitAdminConfigs []RateLimitAdminConfig
	MCMS                  *proposalutils.TimelockConfig
}

func (cfg SetRateLimitAdminConfig) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	for _, rateLimitAdminConfig := range cfg.RateLimitAdminConfigs {
		tokenPubKey := solana.MustPublicKeyFromBase58(rateLimitAdminConfig.SolTokenPubKey)
		tokenPool := chainState.GetActiveTokenPool(rateLimitAdminConfig.PoolType, rateLimitAdminConfig.Metadata)
		if rateLimitAdminConfig.PoolType == "" {
			return errors.New("pool type must be defined")
		}
		chain := e.BlockChains.SolanaChains()[cfg.SolChainSelector]
		if err := chainState.CommonValidation(e, cfg.SolChainSelector, tokenPubKey); err != nil {
			return err
		}
		if err := chainState.ValidatePoolDeployment(&e, rateLimitAdminConfig.PoolType, cfg.SolChainSelector, tokenPubKey, true, rateLimitAdminConfig.Metadata); err != nil {
			return err
		}
		poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
		programData := solTestTokenPool.State{}
		if err := chain.GetAccountDataBorshInto(context.Background(), poolConfigPDA, &programData); err != nil {
			return err
		}
	}
	return nil
}

func SetRateLimitAdmin(e cldf.Environment, cfg SetRateLimitAdminConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chainState := state.SolChains[cfg.SolChainSelector]
	if err := cfg.Validate(e, chainState); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	mcmsTxs := []mcmsTypes.Transaction{}

	for _, rateLimitAdminConfig := range cfg.RateLimitAdminConfigs {
		chain := e.BlockChains.SolanaChains()[cfg.SolChainSelector]
		tokenPubKey := solana.MustPublicKeyFromBase58(rateLimitAdminConfig.SolTokenPubKey)

		var ix solana.Instruction
		tokenPool := chainState.GetActiveTokenPool(rateLimitAdminConfig.PoolType, rateLimitAdminConfig.Metadata)
		poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
		tokenPoolUsingMcms := solanastateview.IsSolanaProgramOwnedByTimelock(
			&e,
			chain,
			chainState,
			rateLimitAdminConfig.PoolType,
			tokenPubKey,
			rateLimitAdminConfig.Metadata,
		)
		authority := GetAuthorityForIxn(
			&e,
			chain,
			chainState,
			rateLimitAdminConfig.PoolType,
			tokenPubKey,
			rateLimitAdminConfig.Metadata,
		)
		switch rateLimitAdminConfig.PoolType {
		case shared.BurnMintTokenPool:
			solBurnMintTokenPool.SetProgramID(tokenPool)
			ix, err = solBurnMintTokenPool.NewSetRateLimitAdminInstruction(
				tokenPubKey,
				rateLimitAdminConfig.NewRateLimitAdmin,
				poolConfigPDA,
				authority,
			).ValidateAndBuild()
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
		case shared.LockReleaseTokenPool:
			solLockReleaseTokenPool.SetProgramID(tokenPool)
			ix, err = solLockReleaseTokenPool.NewSetRateLimitAdminInstruction(
				tokenPubKey,
				rateLimitAdminConfig.NewRateLimitAdmin,
				poolConfigPDA,
				authority,
			).ValidateAndBuild()
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
		case shared.CCTPTokenPool:
			cctp_token_pool.SetProgramID(tokenPool)
			ix, err = cctp_token_pool.NewSetRateLimitAdminInstruction(
				tokenPubKey,
				rateLimitAdminConfig.NewRateLimitAdmin,
				poolConfigPDA,
				authority,
			).ValidateAndBuild()
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
		default:
			return cldf.ChangesetOutput{}, fmt.Errorf("invalid pool type: %s", rateLimitAdminConfig.PoolType)
		}
		if tokenPoolUsingMcms {
			tx, err := BuildMCMSTxn(ix, tokenPool.String(), rateLimitAdminConfig.PoolType)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
			}
			mcmsTxs = append(mcmsTxs, *tx)
		} else {
			if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}
	}

	if len(mcmsTxs) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.SolChainSelector, "proposal to SetRateLimitAdmin in Solana", cfg.MCMS.MinDelay, mcmsTxs)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}
