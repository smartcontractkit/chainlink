package changeset

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/erc20"
	"github.com/smartcontractkit/chainlink/deployment"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
)

var CurrentTokenPoolVersion semver.Version = deployment.Version1_5_1

var TokenTypes = map[deployment.ContractType]struct{}{
	BurnMintToken: {},
	ERC20Token:    {},
	ERC677Token:   {},
}

var TokenPoolTypes = map[deployment.ContractType]struct{}{
	BurnMintTokenPool:              {},
	BurnWithFromMintTokenPool:      {},
	BurnFromMintTokenPool:          {},
	LockReleaseTokenPool:           {},
	USDCTokenPool:                  {},
	HybridLockReleaseUSDCTokenPool: {},
}

var TokenPoolVersions map[semver.Version]struct{} = map[semver.Version]struct{}{
	deployment.Version1_5_1: struct{}{},
}

// TokenPoolInfo defines the type & version of a token pool, along with an optional external administrator.
type TokenPoolInfo struct {
	// Type is the type of the token pool.
	Type deployment.ContractType
	// Version is the version of the token pool.
	Version semver.Version
	// ExternalAdmin is the external administrator of the token pool on the registry.
	ExternalAdmin common.Address
}

func (t TokenPoolInfo) Validate() error {
	// Ensure that the inputted type is known
	if _, ok := TokenPoolTypes[t.Type]; !ok {
		return fmt.Errorf("%s is not a known token pool type", t.Type)
	}

	// Ensure that the inputted version is known
	if _, ok := TokenPoolVersions[t.Version]; !ok {
		return fmt.Errorf("%s is not a known token pool version", t.Version)
	}

	return nil
}

// GetConfigOnRegistry fetches the token's config on the token admin registry.
func (t TokenPoolInfo) GetConfigOnRegistry(
	ctx context.Context,
	symbol TokenSymbol,
	chain deployment.Chain,
	state CCIPChainState,
) (token_admin_registry.TokenAdminRegistryTokenConfig, error) {
	_, tokenAddress, err := t.GetPoolAndTokenAddress(ctx, symbol, chain, state)
	if err != nil {
		return token_admin_registry.TokenAdminRegistryTokenConfig{}, fmt.Errorf("failed to get token pool and token address for %s token on %s: %w", symbol, chain, err)
	}
	tokenAdminRegistry := state.TokenAdminRegistry
	tokenConfig, err := tokenAdminRegistry.GetTokenConfig(&bind.CallOpts{Context: ctx}, tokenAddress)
	if err != nil {
		return token_admin_registry.TokenAdminRegistryTokenConfig{}, fmt.Errorf("failed to get config of %s token with address %s from registry on %s: %w", symbol, tokenAddress, chain, err)
	}
	return tokenConfig, nil
}

// GetPoolAndTokenAddress returns pool bindings and the token address.
func (t TokenPoolInfo) GetPoolAndTokenAddress(
	ctx context.Context,
	symbol TokenSymbol,
	chain deployment.Chain,
	state CCIPChainState,
) (*token_pool.TokenPool, common.Address, error) {
	tokenPoolAddress, ok := GetTokenPoolAddressFromSymbolTypeAndVersion(state, chain, symbol, t.Type, t.Version)
	if !ok {
		return nil, utils.ZeroAddress, fmt.Errorf("token pool does not exist on %s with symbol %s, type %s, and version %s", chain, symbol, t.Type, t.Version)
	}
	tokenPool, err := token_pool.NewTokenPool(tokenPoolAddress, chain.Client)
	if err != nil {
		return nil, utils.ZeroAddress, fmt.Errorf("failed to connect token pool with address %s on chain %s to token pool bindings: %w", tokenPoolAddress, chain, err)
	}
	tokenAddress, err := tokenPool.GetToken(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, utils.ZeroAddress, fmt.Errorf("failed to get token from pool with address %s on %s: %w", tokenPool.Address(), chain, err)
	}
	return tokenPool, tokenAddress, nil
}

// tokenPool defines behavior common to all token pools.
type tokenPool interface {
	GetToken(opts *bind.CallOpts) (common.Address, error)
	TypeAndVersion(*bind.CallOpts) (string, error)
}

// tokenPoolMetadata defines the token pool version version and symbol of the corresponding token.
type tokenPoolMetadata struct {
	Version semver.Version
	Symbol  TokenSymbol
}

// NewTokenPoolWithMetadata returns a token pool along with its metadata.
func NewTokenPoolWithMetadata[P tokenPool](
	ctx context.Context,
	newTokenPool func(address common.Address, backend bind.ContractBackend) (P, error),
	poolAddress common.Address,
	chainClient deployment.OnchainClient,
) (P, tokenPoolMetadata, error) {
	pool, err := newTokenPool(poolAddress, chainClient)
	if err != nil {
		return pool, tokenPoolMetadata{}, fmt.Errorf("failed to connect address %s with token pool bindings: %w", poolAddress, err)
	}
	tokenAddress, err := pool.GetToken(&bind.CallOpts{Context: ctx})
	if err != nil {
		return pool, tokenPoolMetadata{}, fmt.Errorf("failed to get token address from pool with address %s: %w", poolAddress, err)
	}
	typeAndVersionStr, err := pool.TypeAndVersion(&bind.CallOpts{Context: ctx})
	if err != nil {
		return pool, tokenPoolMetadata{}, fmt.Errorf("failed to get type and version from pool with address %s: %w", poolAddress, err)
	}
	_, versionStr, err := ccipconfig.ParseTypeAndVersion(typeAndVersionStr)
	if err != nil {
		return pool, tokenPoolMetadata{}, fmt.Errorf("failed to parse type and version of pool with address %s: %w", poolAddress, err)
	}
	version, err := semver.NewVersion(versionStr)
	if err != nil {
		return pool, tokenPoolMetadata{}, fmt.Errorf("failed parsing version %s of pool with address %s: %w", versionStr, poolAddress, err)
	}
	token, err := erc20.NewERC20(tokenAddress, chainClient)
	if err != nil {
		return pool, tokenPoolMetadata{}, fmt.Errorf("failed to connect address %s with ERC20 bindings: %w", tokenAddress, err)
	}
	symbol, err := token.Symbol(&bind.CallOpts{Context: ctx})
	if err != nil {
		return pool, tokenPoolMetadata{}, fmt.Errorf("failed to fetch symbol from token with address %s: %w", tokenAddress, err)
	}
	return pool, tokenPoolMetadata{
		Symbol:  TokenSymbol(symbol),
		Version: *version,
	}, nil
}

// GetTokenPoolAddressFromSymbolTypeAndVersion returns the token pool address in the environment linked to a particular symbol, type, and version
func GetTokenPoolAddressFromSymbolTypeAndVersion(
	chainState CCIPChainState,
	chain deployment.Chain,
	symbol TokenSymbol,
	poolType deployment.ContractType,
	version semver.Version,
) (common.Address, bool) {
	switch poolType {
	case BurnMintTokenPool:
		if tokenPools, ok := chainState.BurnMintTokenPools[symbol]; ok {
			if tokenPool, ok := tokenPools[version]; ok {
				return tokenPool.Address(), true
			}
		}
	case BurnFromMintTokenPool:
		if tokenPools, ok := chainState.BurnFromMintTokenPools[symbol]; ok {
			if tokenPool, ok := tokenPools[version]; ok {
				return tokenPool.Address(), true
			}
		}
	case BurnWithFromMintTokenPool:
		if tokenPools, ok := chainState.BurnWithFromMintTokenPools[symbol]; ok {
			if tokenPool, ok := tokenPools[version]; ok {
				return tokenPool.Address(), true
			}
		}
	case LockReleaseTokenPool:
		if tokenPools, ok := chainState.LockReleaseTokenPools[symbol]; ok {
			if tokenPool, ok := tokenPools[version]; ok {
				return tokenPool.Address(), true
			}
		}
	case USDCTokenPool:
		if tokenPool, ok := chainState.USDCTokenPools[version]; ok {
			return tokenPool.Address(), true
		}
	case HybridLockReleaseUSDCTokenPool:
		if tokenPool, ok := chainState.USDCTokenPools[version]; ok {
			return tokenPool.Address(), true
		}
	}

	return utils.ZeroAddress, false
}

// TokenAdminRegistryChangesetConfig defines a config for all token admin registry actions.
type TokenAdminRegistryChangesetConfig struct {
	// MCMS defines the delay to use for Timelock (if absent, the changeset will attempt to use the deployer key).
	MCMS *proposalutils.TimelockConfig
	// Pools defines the pools corresponding to the tokens we want to accept admin role for.
	Pools map[uint64]map[TokenSymbol]TokenPoolInfo
	// SkipOwnershipValidation indicates whether or not to skip admin ownership validation of token in the registry.
	// it is skipped when propose admin,set admin and set pool operations are done as part of one changeset.
	SkipOwnershipValidation bool
}

// validateTokenAdminRegistryChangeset validates all token admin registry changesets.
func (c TokenAdminRegistryChangesetConfig) Validate(
	env deployment.Environment,
	mustBeOwner bool,
	registryConfigCheck func(
		config token_admin_registry.TokenAdminRegistryTokenConfig,
		sender common.Address,
		externalAdmin common.Address,
		symbol TokenSymbol,
		chain deployment.Chain,
	) error,
) error {
	state, err := LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	err = state.EnforceMCMSUsageIfProd(env.GetContext(), c.MCMS)
	if err != nil {
		return err
	}
	for chainSelector, symbolToPoolInfo := range c.Pools {
		err := ValidateChain(env, state, chainSelector, c.MCMS)
		if err != nil {
			return fmt.Errorf("failed to validate chain selector %d: %w", chainSelector, err)
		}
		chain, ok := env.Chains[chainSelector]
		if !ok {
			return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
		}
		chainState, ok := state.Chains[chainSelector]
		if !ok {
			return fmt.Errorf("%s does not exist in state", chain)
		}
		if tokenAdminRegistry := chainState.TokenAdminRegistry; tokenAdminRegistry == nil {
			return fmt.Errorf("missing tokenAdminRegistry on %s", chain)
		}

		// Validate that the token admin registry is owned by the address that will be actioning the transactions (i.e. Timelock or deployer key)
		// However, most token admin registry actions aren't owner-protected. They just require you to be the admin.
		if mustBeOwner {
			if err := commoncs.ValidateOwnership(env.GetContext(), c.MCMS != nil, chain.DeployerKey.From, chainState.Timelock.Address(), chainState.TokenAdminRegistry); err != nil {
				return fmt.Errorf("token admin registry failed ownership validation on %s: %w", chain, err)
			}
		}
		for symbol, poolInfo := range symbolToPoolInfo {
			if err := poolInfo.Validate(); err != nil {
				return fmt.Errorf("failed to validate token pool info for %s token on chain %s: %w", symbol, chain, err)
			}

			tokenConfigOnRegistry, err := poolInfo.GetConfigOnRegistry(env.GetContext(), symbol, chain, chainState)
			if err != nil {
				return fmt.Errorf("failed to get state of %s token on chain %s: %w", symbol, chain, err)
			}

			fromAddress := chain.DeployerKey.From // "We" are either the deployer key or the timelock
			if c.MCMS != nil {
				fromAddress = chainState.Timelock.Address()
			}

			if !c.SkipOwnershipValidation {
				err = registryConfigCheck(tokenConfigOnRegistry, fromAddress, poolInfo.ExternalAdmin, symbol, chain)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
