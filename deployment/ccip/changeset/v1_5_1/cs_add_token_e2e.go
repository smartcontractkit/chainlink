package v1_5_1

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/mcms"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/erc20"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/erc677"
)

// AddTokensE2E is a changeset that deploys and configures token pools for multiple tokens across multiple chains in a single changeset.
// AddTokensE2E does the following:
//
//  1. Deploys tokens ( specifically TestTokens) optionally if DeployTokenConfig is provided and
//     populates the pool deployment configuration for each token.
//
//  2. Deploys token pool contracts for each token specified in the config.
//     If the token deployment config is provided, pool deployment configuration DeployTokenPoolContractsConfig is not required.
//     It will use the token deployment config to deploy the token and
//     populate DeployTokenPoolContractsConfig.
//     If the token deployment config is not provided, pool deployment configuration DeployTokenPoolContractsConfig is mandatory.
//
//  3. Configures pools -
//     If the chain is already supported -
//
//     i. it updates the rate limits for the chain
//     ii. it adds a new remote pool if the token pool on the remote chain is being updated
//
//     If the chain is not supported -
//
//     i. it adds chain support with the desired rate limits
//     i. it adds the desired remote pool addresses to the token pool on the chain
//     iii. if there used to be an existing token pool on tokenadmin_registry, it adds the remote pool addresses of that token pool to ensure 0 downtime
//
// 4. Proposes admin rights for the token on the token admin registry
//
// If the token admin is not an external address -
// 5. Accepts admin rights for the token on the token admin registry
// 6. Sets the pool for the token on the token admin registry
var AddTokensE2E = deployment.CreateChangeSet(addTokenE2ELogic, addTokenE2EPreconditionValidation)

type E2ETokenAndPoolConfig struct {
	TokenDeploymentConfig *DeployTokenConfig    // TokenDeploymentConfig is optional. If provided, it will be used to deploy the token and populate the pool deployment configuration.
	DeployPoolConfig      *DeployTokenPoolInput // Deployment configuration for pools is not needed if tokenDeploymentConfig is provided. This will be populated from the tokenDeploymentConfig if it is provided.
	PoolVersion           semver.Version
	ExternalAdmin         common.Address // ExternalAdmin is the external administrator of the token pool on the registry.
	RateLimiterConfig     RateLimiterPerChain
	// OverrideTokenSymbol is the token symbol to use to override against main symbol (ex: override to clCCIP-LnM when the main token symbol is CCIP-LnM)
	// WARNING: This should only be used in exceptional cases where the token symbol on a particular chain differs from the main tokenSymbol
	OverrideTokenSymbol changeset.TokenSymbol
}

type AddTokenE2EConfig struct {
	PoolConfig   map[uint64]E2ETokenAndPoolConfig
	IsTestRouter bool

	// internal fields - To be populated from the PoolConfig.
	// User do not need to populate these fields.
	deployPool             DeployTokenPoolContractsConfig
	configurePools         ConfigureTokenPoolContractsConfig
	configureTokenAdminReg changeset.TokenAdminRegistryChangesetConfig
}

// newConfigurePoolAndTokenAdminRegConfig populated internal fields in AddTokenE2EConfig.
// It creates the configuration for deploying and configuring token pools and token admin registry.
// It then validates the configuration.
func (c *AddTokenE2EConfig) newConfigurePoolAndTokenAdminRegConfig(e deployment.Environment, symbol changeset.TokenSymbol, timelockCfg *proposalutils.TimelockConfig) error {
	c.deployPool = DeployTokenPoolContractsConfig{
		TokenSymbol:  symbol,
		NewPools:     make(map[uint64]DeployTokenPoolInput),
		IsTestRouter: c.IsTestRouter,
	}
	c.configurePools = ConfigureTokenPoolContractsConfig{
		TokenSymbol: symbol,
		MCMS:        nil, // as token pools are deployed as part of the changeset, the pools will still be owned by the deployer key
		PoolUpdates: make(map[uint64]TokenPoolConfig),
	}
	c.configureTokenAdminReg = changeset.TokenAdminRegistryChangesetConfig{
		MCMS:  timelockCfg,
		Pools: make(map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo),
	}
	for chain, poolCfg := range c.PoolConfig {
		c.deployPool.NewPools[chain] = *poolCfg.DeployPoolConfig
		c.configurePools.PoolUpdates[chain] = TokenPoolConfig{
			ChainUpdates:        poolCfg.RateLimiterConfig,
			Type:                poolCfg.DeployPoolConfig.Type,
			Version:             poolCfg.PoolVersion,
			OverrideTokenSymbol: poolCfg.OverrideTokenSymbol,
		}

		// Populate the TokenAdminRegistryChangesetConfig for each chain.
		if _, ok := c.configureTokenAdminReg.Pools[chain]; !ok {
			c.configureTokenAdminReg.Pools[chain] = make(map[changeset.TokenSymbol]changeset.TokenPoolInfo)
		}
		c.configureTokenAdminReg.Pools[chain][symbol] = changeset.TokenPoolInfo{
			Version:       poolCfg.PoolVersion,
			ExternalAdmin: poolCfg.ExternalAdmin,
			Type:          poolCfg.DeployPoolConfig.Type,
		}
	}
	if err := c.deployPool.Validate(e); err != nil {
		return fmt.Errorf("failed to validate deploy pool config: %w", err)
	}
	// rest of the validation should be done after token pools are deployed
	return nil
}

func (c *AddTokenE2EConfig) newDeployTokenPoolConfigAfterTokenDeployment(tokenAddresses map[uint64]common.Address) error {
	deployTokenCfg := make(map[uint64]DeployTokenPoolInput) // This will hold the deployment configuration for each token.
	for chain, p := range c.PoolConfig {
		tokenAddress, ok := tokenAddresses[chain]
		if !ok {
			// If the token address is not found for the chain, return an error.
			return fmt.Errorf("token address not found for chain %d", chain)
		}
		if p.TokenDeploymentConfig == nil {
			continue
		}
		tp := DeployTokenPoolInput{
			TokenAddress:       tokenAddress,                          // The address of the token deployed on the chain.
			LocalTokenDecimals: p.TokenDeploymentConfig.TokenDecimals, // The decimals of the token deployed on the chain.
			Type:               p.TokenDeploymentConfig.PoolType,      // The type of the token pool (e.g. LockRelease, BurnMint).
			AllowList:          p.TokenDeploymentConfig.PoolAllowList,
			AcceptLiquidity:    p.TokenDeploymentConfig.AcceptLiquidity,
		}
		deployTokenCfg[chain] = tp // Add the pool configuration for the chain to the deployment config.
		p.DeployPoolConfig = &tp
		c.PoolConfig[chain] = p
	}
	return nil
}

type DeployTokenConfig struct {
	TokenName       string
	TokenSymbol     changeset.TokenSymbol
	TokenDecimals   uint8    // needed for BurnMintToken only
	MaxSupply       *big.Int // needed for BurnMintToken only
	Type            deployment.ContractType
	PoolType        deployment.ContractType // This is the type of the token pool that will be deployed for this token.
	PoolAllowList   []common.Address
	AcceptLiquidity *bool
}

func (c *DeployTokenConfig) Validate() error {
	if c.TokenName == "" {
		return errors.New("token name must be defined")
	}
	if c.TokenDecimals == 0 && c.Type == changeset.BurnMintToken {
		return errors.New("token decimals must be defined for BurnMintToken type")
	}
	if c.MaxSupply == nil && c.Type == changeset.BurnMintToken {
		return errors.New("max supply must be defined for BurnMintToken type")
	}
	if _, ok := changeset.TokenPoolTypes[c.PoolType]; !ok {
		return fmt.Errorf("token pool type not supported %s", c.PoolType)
	}
	if _, ok := changeset.TokenTypes[c.Type]; !ok {
		return fmt.Errorf("token type not supported %s", c.Type)
	}
	return nil
}

type AddTokensE2EConfig struct {
	Tokens map[changeset.TokenSymbol]AddTokenE2EConfig
	MCMS   *proposalutils.TimelockConfig
}

func addTokenE2EPreconditionValidation(e deployment.Environment, config AddTokensE2EConfig) error {
	if len(config.Tokens) == 0 {
		return nil
	}
	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for token, cfg := range config.Tokens {
		for chain, poolCfg := range cfg.PoolConfig {
			if err := changeset.ValidateChain(e, state, chain, config.MCMS); err != nil {
				return fmt.Errorf("failed to validate chain %d: %w", chain, err)
			}
			if (poolCfg.DeployPoolConfig != nil) == (poolCfg.TokenDeploymentConfig != nil) {
				return fmt.Errorf("must provide either DeploymentConfig or TokenDeploymentConfig for token %s: cannot provide both or neither", token)
			}
			if poolCfg.TokenDeploymentConfig != nil {
				if poolCfg.TokenDeploymentConfig.TokenSymbol != token {
					return fmt.Errorf("token symbol %s in token deployment config does not match token %s", poolCfg.TokenDeploymentConfig.TokenSymbol, token)
				}
				if err := poolCfg.TokenDeploymentConfig.Validate(); err != nil {
					return fmt.Errorf("failed to validate token deployment config for token %s: %w", token, err)
				}
				// the rest of the internal fields are populated from the PoolConfig and it will be validated once the tokens are deployed
			} else {
				if poolCfg.DeployPoolConfig == nil {
					return fmt.Errorf("must provide pool DeploymentConfig for token %s when TokenDeploymentConfig is not provided", token)
				}
				if err := poolCfg.DeployPoolConfig.Validate(e.GetContext(), e.Chains[chain], state.Chains[chain], token); err != nil {
					return fmt.Errorf("failed to validate token pool config for token %s: %w", token, err)
				}
				// populate the internal fields for deploying and configuring token pools and token admin registry and validate them
				err := cfg.newConfigurePoolAndTokenAdminRegConfig(e, token, config.MCMS)
				if err != nil {
					return err
				}
				config.Tokens[token] = cfg
			}
		}
	}
	return nil
}

func addTokenE2ELogic(env deployment.Environment, config AddTokensE2EConfig) (deployment.ChangesetOutput, error) {
	if len(config.Tokens) == 0 {
		return deployment.ChangesetOutput{}, nil
	}
	// use a clone of env to avoid modifying the original env
	e := env.Clone()
	finalCSOut := &deployment.ChangesetOutput{
		AddressBook: deployment.NewMemoryAddressBook(),
	}
	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	for token, cfg := range config.Tokens {
		e.Logger.Infow("starting token addition operations for", "token", token, "chains", maps.Keys(cfg.PoolConfig))
		tokenDeployCfg := make(map[uint64]DeployTokenConfig)
		for chain, poolCfg := range cfg.PoolConfig {
			if poolCfg.TokenDeploymentConfig != nil {
				tokenDeployCfg[chain] = *poolCfg.TokenDeploymentConfig
			}
		}
		// deploy token pools if token deployment config is provided and populate pool deployment configuration
		if len(tokenDeployCfg) > 0 {
			deployedTokens, ab, err := deployTokenPools(e, tokenDeployCfg)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			if err := cfg.newDeployTokenPoolConfigAfterTokenDeployment(deployedTokens); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to populate pool deployment configuration: %w", err)
			}
			e.Logger.Infow("deployed token and created pool deployment config", "token", token)
			if err := finalCSOut.AddressBook.Merge(ab); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge address book for token %s: %w", token, err)
			}
			// populate the configuration for deploying and configuring token pools and token admin registry
			if err := cfg.newConfigurePoolAndTokenAdminRegConfig(e, token, config.MCMS); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to populate configuration for "+
					"deploying and configuring token pools and token admin registry: %w", err)
			}
		}
		output, err := DeployTokenPoolContractsChangeset(e, cfg.deployPool)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy token pool for token %s: %w", token, err)
		}
		if err := deployment.MergeChangesetOutput(e, finalCSOut, output); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge address book for token %s: %w", token, err)
		}
		newAddresses, err := output.AddressBook.Addresses()
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get addresses from address book: %w", err)
		}
		e.Logger.Infow("deployed token pool", "token", token, "addresses", newAddresses)
		if err := cfg.configurePools.Validate(e); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to validate configure pool config: %w", err)
		}
		// Validate the configure token admin reg config.
		// As we will perform proposing admin, accepting admin and setting pool on same changeset
		// we are only validating the propose admin role.
		if err := cfg.configureTokenAdminReg.Validate(e, true, validateProposeAdminRole); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to validate configure token admin reg config: %w", err)
		}
		output, err = ConfigureTokenPoolContractsChangeset(e, cfg.configurePools)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to configure token pool for token %s: %w", token, err)
		}
		if err := deployment.MergeChangesetOutput(e, finalCSOut, output); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge changeset output after configuring token pool for token %s: %w", token, err)
		}
		e.Logger.Infow("configured token pool", "token", token)

		output, err = ProposeAdminRoleChangeset(e, cfg.configureTokenAdminReg)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to propose admin role for token %s: %w", token, err)
		}
		if err := deployment.MergeChangesetOutput(e, finalCSOut, output); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to changeset output after configuring token admin reg for token %s: %w",
				token, err)
		}
		e.Logger.Infow("proposed admin role", "token", token, "config", cfg.configureTokenAdminReg)

		// find all tokens for which there is no external admin
		// for those tokens, accept the admin role and set the pool
		updatedConfigureTokenAdminReg := changeset.TokenAdminRegistryChangesetConfig{
			MCMS:  config.MCMS,
			Pools: make(map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo),
			// SkipOwnershipValidation is set to true as we are accepting admin role and setting token pool as part of one changeset
			SkipOwnershipValidation: true,
		}
		for chain, poolInfo := range cfg.configureTokenAdminReg.Pools {
			for symbol, info := range poolInfo {
				if info.ExternalAdmin == utils.ZeroAddress {
					if updatedConfigureTokenAdminReg.Pools[chain] == nil {
						updatedConfigureTokenAdminReg.Pools[chain] = make(map[changeset.TokenSymbol]changeset.TokenPoolInfo)
					}
					updatedConfigureTokenAdminReg.Pools[chain][symbol] = info
				}
			}
		}
		// if there are no tokens for which there is no external admin, continue to next token
		if len(updatedConfigureTokenAdminReg.Pools) == 0 {
			continue
		}
		output, err = AcceptAdminRoleChangeset(e, updatedConfigureTokenAdminReg)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to accept admin role for token %s: %w", token, err)
		}
		if err := deployment.MergeChangesetOutput(e, finalCSOut, output); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge address book for token %s: %w", token, err)
		}
		e.Logger.Infow("accepted admin role", "token", token, "config", updatedConfigureTokenAdminReg)
		output, err = SetPoolChangeset(e, updatedConfigureTokenAdminReg)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to set pool for token %s: %w", token, err)
		}
		if err := deployment.MergeChangesetOutput(e, finalCSOut, output); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge address book for token %s: %w", token, err)
		}
		e.Logger.Infow("set pool", "token", token, "config", updatedConfigureTokenAdminReg)
	}
	// if there are multiple proposals, aggregate them so that we don't have to propose them separately
	if len(finalCSOut.MCMSTimelockProposals) > 1 {
		aggregatedProposals, err := proposalutils.AggregateProposals(
			e, state.EVMMCMSStateByChain(), finalCSOut.MCMSTimelockProposals, nil,
			"Add Tokens E2E", config.MCMS)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to aggregate proposals: %w", err)
		}
		finalCSOut.MCMSTimelockProposals = []mcms.TimelockProposal{*aggregatedProposals}
	}
	return *finalCSOut, nil
}

func deployTokenPools(e deployment.Environment, tokenDeployCfg map[uint64]DeployTokenConfig) (map[uint64]common.Address, deployment.AddressBook, error) {
	ab := deployment.NewMemoryAddressBook()
	tokenAddresses := make(map[uint64]common.Address) // This will hold the token addresses for each chain.
	for selector, cfg := range tokenDeployCfg {
		switch cfg.Type {
		case changeset.BurnMintToken:
			token, err := deployment.DeployContract(e.Logger, e.Chains[selector], ab,
				func(chain deployment.Chain) deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
					tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
						e.Chains[selector].DeployerKey,
						e.Chains[selector].Client,
						cfg.TokenName,
						string(cfg.TokenSymbol),
						cfg.TokenDecimals,
						cfg.MaxSupply,
					)
					return deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
						Address:  tokenAddress,
						Contract: token,
						Tv:       deployment.NewTypeAndVersion(changeset.BurnMintToken, deployment.Version1_0_0),
						Tx:       tx,
						Err:      err,
					}
				},
			)
			if err != nil {
				return nil, ab, fmt.Errorf("failed to deploy BurnMintERC677 token %s on chain %d: %w", cfg.TokenName, selector, err)
			}
			tokenAddresses[selector] = token.Address
		case changeset.ERC20Token:
			token, err := deployment.DeployContract(e.Logger, e.Chains[selector], ab,
				func(chain deployment.Chain) deployment.ContractDeploy[*erc20.ERC20] {
					tokenAddress, tx, token, err := erc20.DeployERC20(
						e.Chains[selector].DeployerKey,
						e.Chains[selector].Client,
						cfg.TokenName,
						string(cfg.TokenSymbol),
					)
					return deployment.ContractDeploy[*erc20.ERC20]{
						Address:  tokenAddress,
						Contract: token,
						Tv:       deployment.NewTypeAndVersion(changeset.ERC20Token, deployment.Version1_0_0),
						Tx:       tx,
						Err:      err,
					}
				},
			)
			if err != nil {
				return nil, ab, fmt.Errorf("failed to deploy ERC20 token %s on chain %d: %w", cfg.TokenName, selector, err)
			}
			tokenAddresses[selector] = token.Address
		case changeset.ERC677Token:
			token, err := deployment.DeployContract(e.Logger, e.Chains[selector], ab,
				func(chain deployment.Chain) deployment.ContractDeploy[*erc677.ERC677] {
					tokenAddress, tx, token, err := erc677.DeployERC677(
						e.Chains[selector].DeployerKey,
						e.Chains[selector].Client,
						cfg.TokenName,
						string(cfg.TokenSymbol),
					)
					return deployment.ContractDeploy[*erc677.ERC677]{
						Address:  tokenAddress,
						Contract: token,
						Tv:       deployment.NewTypeAndVersion(changeset.ERC677Token, deployment.Version1_0_0),
						Tx:       tx,
						Err:      err,
					}
				},
			)
			if err != nil {
				return nil, ab, fmt.Errorf("failed to deploy ERC677 token %s on chain %d: %w", cfg.TokenName, selector, err)
			}
			tokenAddresses[selector] = token.Address
		default:
			return nil, ab, fmt.Errorf("unsupported token %s type %s for deployment on chain %d", cfg.TokenName, cfg.Type, selector)
		}
	}
	return tokenAddresses, ab, nil
}
