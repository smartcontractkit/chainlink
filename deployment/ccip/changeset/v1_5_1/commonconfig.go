package v1_5_1

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"
	deployment2 "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	changeset2 "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// TokenAdminRegistryChangesetConfig defines a config for all token admin registry actions.
type TokenAdminRegistryChangesetConfig struct {
	// MCMS defines the delay to use for Timelock (if absent, the changeset will attempt to use the deployer key).
	MCMS *proposalutils.TimelockConfig
	// Pools defines the pools corresponding to the tokens we want to accept admin role for.
	Pools map[uint64]map[shared.TokenSymbol]TokenPoolInfo
	// SkipOwnershipValidation indicates whether or not to skip admin ownership validation of token in the registry.
	// it is skipped when propose admin,set admin and set pool operations are done as part of one changeset.
	SkipOwnershipValidation bool
}

// validateTokenAdminRegistryChangeset validates all token admin registry changesets.
func (c TokenAdminRegistryChangesetConfig) Validate(
	env cldf.Environment,
	mustBeOwner bool,
	registryConfigCheck func(
		config token_admin_registry.TokenAdminRegistryTokenConfig,
		sender common.Address,
		externalAdmin common.Address,
		symbol shared.TokenSymbol,
		chain cldf_evm.Chain,
	) error,
) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	err = state.EnforceMCMSUsageIfProd(env.GetContext(), c.MCMS)
	if err != nil {
		return err
	}
	for chainSelector, symbolToPoolInfo := range c.Pools {
		err := stateview.ValidateChain(env, state, chainSelector, c.MCMS)
		if err != nil {
			return fmt.Errorf("failed to validate chain selector %d: %w", chainSelector, err)
		}
		chain, ok := env.BlockChains.EVMChains()[chainSelector]
		if !ok {
			return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
		}
		chainState, ok := state.EVMChainState(chainSelector)
		if !ok {
			return fmt.Errorf("%s does not exist in state", chain)
		}
		if tokenAdminRegistry := chainState.TokenAdminRegistry; tokenAdminRegistry == nil {
			return fmt.Errorf("missing tokenAdminRegistry on %s", chain)
		}

		// Validate that the token admin registry is owned by the address that will be actioning the transactions (i.e. Timelock or deployer key)
		// However, most token admin registry actions aren't owner-protected. They just require you to be the admin.
		if mustBeOwner {
			if err := changeset2.ValidateOwnership(env.GetContext(), c.MCMS != nil, chain.DeployerKey.From, chainState.Timelock.Address(), chainState.TokenAdminRegistry); err != nil {
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

// TokenPoolInfo defines the type & version of a token pool, along with an optional external administrator.
type TokenPoolInfo struct {
	// Type is the type of the token pool.
	Type deployment2.ContractType
	// Version is the version of the token pool.
	Version semver.Version
	// ExternalAdmin is the external administrator of the token pool on the registry.
	ExternalAdmin common.Address
}

func (t TokenPoolInfo) Validate() error {
	// Ensure that the inputted type is known
	if _, ok := shared.TokenPoolTypes[t.Type]; !ok {
		return fmt.Errorf("%s is not a known token pool type", t.Type)
	}

	// Ensure that the inputted version is known
	if _, ok := shared.TokenPoolVersions[t.Version]; !ok {
		return fmt.Errorf("%s is not a known token pool version", t.Version)
	}

	return nil
}

// GetConfigOnRegistry fetches the token's config on the token admin registry.
func (t TokenPoolInfo) GetConfigOnRegistry(
	ctx context.Context,
	symbol shared.TokenSymbol,
	chain cldf_evm.Chain,
	state evm.CCIPChainState,
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
	symbol shared.TokenSymbol,
	chain cldf_evm.Chain,
	state evm.CCIPChainState,
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

// GetTokenPoolAddressFromSymbolTypeAndVersion returns the token pool address in the environment linked to a particular symbol, type, and version
func GetTokenPoolAddressFromSymbolTypeAndVersion(
	chainState evm.CCIPChainState,
	chain cldf_evm.Chain,
	symbol shared.TokenSymbol,
	poolType deployment2.ContractType,
	version semver.Version,
) (common.Address, bool) {
	switch poolType {
	case shared.BurnMintTokenPool:
		if tokenPools, ok := chainState.BurnMintTokenPools[symbol]; ok {
			if tokenPool, ok := tokenPools[version]; ok {
				return tokenPool.Address(), true
			}
		}
	case shared.BurnFromMintTokenPool:
		if tokenPools, ok := chainState.BurnFromMintTokenPools[symbol]; ok {
			if tokenPool, ok := tokenPools[version]; ok {
				return tokenPool.Address(), true
			}
		}
	case shared.BurnWithFromMintTokenPool:
		if tokenPools, ok := chainState.BurnWithFromMintTokenPools[symbol]; ok {
			if tokenPool, ok := tokenPools[version]; ok {
				return tokenPool.Address(), true
			}
		}
	case shared.LockReleaseTokenPool:
		if tokenPools, ok := chainState.LockReleaseTokenPools[symbol]; ok {
			if tokenPool, ok := tokenPools[version]; ok {
				return tokenPool.Address(), true
			}
		}
	case shared.USDCTokenPool:
		if tokenPool, ok := chainState.USDCTokenPools[version]; ok {
			return tokenPool.Address(), true
		}
	case shared.HybridLockReleaseUSDCTokenPool:
		if tokenPool, ok := chainState.USDCTokenPools[version]; ok {
			return tokenPool.Address(), true
		}
	}

	return utils.ZeroAddress, false
}
