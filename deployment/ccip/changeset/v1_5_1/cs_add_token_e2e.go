package v1_5_1

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/mcms"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/erc20"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/erc677"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
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
var AddTokensE2E = cldf.CreateChangeSet(addTokenE2ELogic, addTokenE2EPreconditionValidation)

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
	TokenName              string
	TokenSymbol            changeset.TokenSymbol
	TokenDecimals          uint8    // needed for BurnMintToken only
	MaxSupply              *big.Int // needed for BurnMintToken only
	Type                   deployment.ContractType
	PoolType               deployment.ContractType // This is the type of the token pool that will be deployed for this token.
	PoolAllowList          []common.Address
	AcceptLiquidity        *bool
	MintTokenForRecipients map[common.Address]*big.Int // MintTokenForRecipients is a map of recipient address to amount to be transferred or minted and provided minting role after token deployment.
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
			deployedTokens, ab, err := deployTokens(e, tokenDeployCfg)
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

func deployTokens(e deployment.Environment, tokenDeployCfg map[uint64]DeployTokenConfig) (map[uint64]common.Address, deployment.AddressBook, error) {
	ab := deployment.NewMemoryAddressBook()
	tokenAddresses := make(map[uint64]common.Address) // This will hold the token addresses for each chain.
	for selector, cfg := range tokenDeployCfg {
		switch cfg.Type {
		case changeset.BurnMintToken:
			token, err := cldf.DeployContract(e.Logger, e.Chains[selector], ab,
				func(chain deployment.Chain) cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
					tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
						e.Chains[selector].DeployerKey,
						e.Chains[selector].Client,
						cfg.TokenName,
						string(cfg.TokenSymbol),
						cfg.TokenDecimals,
						cfg.MaxSupply,
					)
					return cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
						Address:  tokenAddress,
						Contract: token,
						Tv:       deployment.NewTypeAndVersion(changeset.BurnMintToken, deployment.Version1_0_0),
						Tx:       tx,
						Err:      err,
					}
				},
			)
			if err != nil {
				return nil, ab, fmt.Errorf("failed to deploy BurnMintERC677 token "+
					"%s on chain %d: %w", cfg.TokenName, selector, err)
			}
			if err := addMinterAndMintToken(e, selector, token.Contract, e.Chains[selector].DeployerKey.From,
				new(big.Int).Mul(big.NewInt(1_000), big.NewInt(1_000_000_000))); err != nil {
				return nil, ab, fmt.Errorf("failed to add minter and mint token "+
					"%s on chain %d: %w", cfg.TokenName, selector, err)
			}
			if len(cfg.MintTokenForRecipients) > 0 {
				for recipient, amount := range cfg.MintTokenForRecipients {
					if err := addMinterAndMintToken(e, selector, token.Contract, recipient,
						amount); err != nil {
						return nil, ab, fmt.Errorf("failed to add minter and mint "+
							"token %s on chain %d: %w", cfg.TokenName, selector, err)
					}
				}
			}

			tokenAddresses[selector] = token.Address
		case changeset.ERC20Token:
			token, err := cldf.DeployContract(e.Logger, e.Chains[selector], ab,
				func(chain deployment.Chain) cldf.ContractDeploy[*erc20.ERC20] {
					tokenAddress, tx, token, err := erc20.DeployERC20(
						e.Chains[selector].DeployerKey,
						e.Chains[selector].Client,
						cfg.TokenName,
						string(cfg.TokenSymbol),
					)
					return cldf.ContractDeploy[*erc20.ERC20]{
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
			token, err := cldf.DeployContract(e.Logger, e.Chains[selector], ab,
				func(chain deployment.Chain) cldf.ContractDeploy[*erc677.ERC677] {
					tokenAddress, tx, token, err := erc677.DeployERC677(
						e.Chains[selector].DeployerKey,
						e.Chains[selector].Client,
						cfg.TokenName,
						string(cfg.TokenSymbol),
					)
					return cldf.ContractDeploy[*erc677.ERC677]{
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

// grantAccessToPool grants the token pool contract access to mint and burn tokens.
func grantAccessToPool(
	ctx context.Context,
	chain deployment.Chain,
	tpAddress common.Address,
	tokenAddress common.Address,
) error {
	token, err := burn_mint_erc677.NewBurnMintERC677(tokenAddress, chain.Client)
	if err != nil {
		return fmt.Errorf("failed to connect address %s with erc677 bindings: %w", tokenAddress, err)
	}
	owner, err := token.Owner(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get owner of token %s: %w", tokenAddress, err)
	}
	// check if the owner is the deployer key and in that case grant access to the token pool
	if owner == chain.DeployerKey.From {
		tx, err := token.GrantMintAndBurnRoles(chain.DeployerKey, tpAddress)
		if err != nil {
			return fmt.Errorf("failed to grant mint and burn roles to token pool address: %s for token: %s %w", tpAddress, tokenAddress, err)
		}
		if _, err = chain.Confirm(tx); err != nil {
			return fmt.Errorf("failed to wait for transaction %s on chain %d: %w", tx.Hash().Hex(), chain.Selector, err)
		}
	}

	return nil
}

// addMinterAndMintToken adds the minter role to the recipient and mints the specified amount of tokens to the recipient's address.
func addMinterAndMintToken(env deployment.Environment, selector uint64, token *burn_mint_erc677.BurnMintERC677, recipient common.Address, amount *big.Int) error {
	deployerKey := env.Chains[selector].DeployerKey
	ctx := env.GetContext()
	// check if owner is the deployer key
	owner, err := token.Owner(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get owner of token %s on chain %d: %w", token.Address().Hex(), selector, err)
	}
	if owner != deployerKey.From {
		return fmt.Errorf("owner of token %s on chain %d is not the deployer key", token.Address().Hex(), selector)
	}
	// Grant minter role to the given address
	tx, err := token.GrantMintRole(deployerKey, recipient)
	if err != nil {
		return fmt.Errorf("failed to grant mint role to %s on chain %d: %w", recipient.Hex(), selector, err)
	}
	if _, err := env.Chains[selector].Confirm(tx); err != nil {
		return fmt.Errorf("failed to wait for transaction %s on chain %d: %w", tx.Hash().Hex(), selector, err)
	}
	env.Logger.Infow("Transaction granting mint role mined successfully",
		"Hash", tx.Hash().Hex(), "Selector", selector)

	// Mint tokens to the given address and verify the balance
	tx, err = token.Mint(deployerKey, recipient, amount)
	if err != nil {
		return fmt.Errorf("failed to mint %s tokens to %s on chain %d: %w",
			token.Address().Hex(), recipient.Hex(), selector, err)
	}
	if _, err := env.Chains[selector].Confirm(tx); err != nil {
		return fmt.Errorf("failed to wait for transaction %s on chain %d: %w",
			tx.Hash().Hex(), selector, err)
	}
	env.Logger.Infow("Transaction minting token mined successfully", "Hash", tx.Hash().Hex())

	balance, err := token.BalanceOf(&bind.CallOpts{Context: ctx}, recipient)
	if err != nil {
		return fmt.Errorf("failed to get balance of %s on chain %d: %w",
			recipient.Hex(), selector, err)
	}
	if balance.Cmp(amount) != 0 {
		return fmt.Errorf("expected balance of %s, got %s",
			amount.String(), balance.String())
	}
	symbol, err := token.Symbol(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get token symbol for %s on chain %d: %w",
			token.Address().Hex(), selector, err)
	}
	env.Logger.Infow("Recipient added as minter and token minted",
		"Address", recipient.Hex(),
		"Balance", balance.String(), "Token Symbol", symbol,
		"Token address", token.Address().Hex())

	return nil
}
