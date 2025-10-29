package v1_5_1

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/mcms"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc20_with_drip"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc20"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/erc20"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/erc677"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// AddTokensE2E is a changeset that deploys and configures token pools for multiple tokens across multiple chains in a single changeset.
// AddTokensE2E does the following:
//
//  1. Deploys tokens (specifically TestTokens) optionally if DeployTokenConfig is provided and
//     populates the pool deployment configuration for each token.
//
//  2. Deploys token pool contracts for each token specified in the config.
//     If the token deployment config is provided, pool deployment configuration DeployTokenPoolContractsConfig is not required.
//     It will use the token deployment config to deploy the token and
//     populate DeployTokenPoolContractsConfig.
//     If the token deployment config is not provided, pool deployment configuration DeployTokenPoolContractsConfig is mandatory.
//
// 3. Proposes admin rights for the token on the token admin registry
//
//		If the token admin is not an external address -
//		3a. Accepts admin rights for the token on the token admin registry
//		3b. Sets the pool for the token on the token admin registry
//
//	 4. Configures pools -
//		   The configuration is set by via the input field ConfigurePools. This is so that we can have flexibility during configuration.
//		   For ex: deploy chainB token, pool and configure with already deployed chainA tokenPool.
//		   If config is not provided in the input, it skips this step.
var AddTokensE2E = cldf.CreateChangeSet(addTokenE2ELogic, addTokenE2EPreconditionValidation)

type E2ETokenAndPoolConfig struct {
	// TokenDeploymentConfig is optional. If provided, it will be used to deploy the token
	// and populate the pool deployment configuration.
	TokenDeploymentConfig *DeployTokenConfig `json:"tokenDeploymentConfig,omitempty"`

	// Deployment configuration for pools is not needed if tokenDeploymentConfig is provided.
	// This will be populated from the tokenDeploymentConfig if it is provided.
	DeployPoolConfig *DeployTokenPoolInput `json:"deployPoolConfig,omitempty"`

	// Version of the pool being deployed.
	PoolVersion semver.Version `json:"poolVersion"`

	// ExternalAdmin is the external administrator of the token pool on the registry.
	ExternalAdmin common.Address `json:"externalAdmin"`

	// OverrideTokenSymbol is the token symbol to use to override against main symbol
	// (ex: override to clCCIP-LnM when the main token symbol is CCIP-LnM)
	// WARNING: This should only be used in exceptional cases where the token symbol on
	// a particular chain differs from the main tokenSymbol
	OverrideTokenSymbol shared.TokenSymbol `json:"overrideTokenSymbol,omitempty"`
}

type AddTokenE2EConfig struct {
	// Map of chain ID to E2ETokenAndPoolConfig.
	PoolConfig map[uint64]E2ETokenAndPoolConfig `json:"poolConfig"`

	// Whether this is a test router configuration.
	IsTestRouter bool `json:"isTestRouter"`

	// Configures the pools, if empty, deployed pools aren't configured.
	ConfigurePools ConfigureTokenPoolContractsConfig `json:"configurePools"`

	// Internal fields - To be populated from the PoolConfig.
	// User does not need to populate these fields.
	deployPool             *DeployTokenPoolContractsConfig
	configureTokenAdminReg TokenAdminRegistryChangesetConfig
}

// newConfigurePoolAndTokenAdminRegConfig populated internal fields in AddTokenE2EConfig.
// It creates the configuration for deploying and configuring token pools and token admin registry.
// It then validates the configuration.
func (c *AddTokenE2EConfig) newConfigurePoolAndTokenAdminRegConfig(e cldf.Environment, symbol shared.TokenSymbol, timelockCfg *proposalutils.TimelockConfig) error {
	c.deployPool = &DeployTokenPoolContractsConfig{
		TokenSymbol:  symbol,
		NewPools:     make(map[uint64]DeployTokenPoolInput),
		IsTestRouter: c.IsTestRouter,
	}

	c.configureTokenAdminReg = TokenAdminRegistryChangesetConfig{
		MCMS:  timelockCfg,
		Pools: make(map[uint64]map[shared.TokenSymbol]TokenPoolInfo),
	}
	for chain, poolCfg := range c.PoolConfig {
		tpCfg := TokenPoolConfig{}
		poolInfo := TokenPoolInfo{
			Version:       poolCfg.PoolVersion,
			ExternalAdmin: poolCfg.ExternalAdmin,
		}

		c.deployPool.NewPools[chain] = *poolCfg.DeployPoolConfig
		tpCfg.Type = poolCfg.DeployPoolConfig.Type
		poolInfo.Type = poolCfg.DeployPoolConfig.Type

		// Populate the TokenAdminRegistryChangesetConfig for each chain.
		if _, ok := c.configureTokenAdminReg.Pools[chain]; !ok {
			c.configureTokenAdminReg.Pools[chain] = make(map[shared.TokenSymbol]TokenPoolInfo)
		}
		c.configureTokenAdminReg.Pools[chain][symbol] = poolInfo
	}

	if err := c.deployPool.Validate(e); err != nil {
		return fmt.Errorf("failed to validate deploy pool config: %w", err)
	}

	// The rest of the validation should be done after token pools are deployed
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
			TokenType:          p.TokenDeploymentConfig.Type,
			AllowList:          p.TokenDeploymentConfig.PoolAllowList,
			AcceptLiquidity:    p.TokenDeploymentConfig.AcceptLiquidity,
			CCIPAdmin:          p.TokenDeploymentConfig.CCIPAdmin,
		}
		deployTokenCfg[chain] = tp // Add the pool configuration for the chain to the deployment config.
		p.DeployPoolConfig = &tp
		c.PoolConfig[chain] = p
	}
	return nil
}

type DeployTokenConfig struct {
	// TokenName is the full name of the token.
	TokenName string `json:"tokenName"`

	// TokenSymbol is the symbol for the token (e.g., LINK, USDC).
	TokenSymbol shared.TokenSymbol `json:"tokenSymbol"`

	// TokenDecimals specifies how many decimals the token uses.
	// Needed for BurnMintToken only.
	TokenDecimals uint8 `json:"tokenDecimals,omitempty"`

	// MaxSupply defines the maximum supply for the token.
	// Needed for BurnMintToken only.
	MaxSupply *big.Int `json:"maxSupply,omitempty"`

	// Type is the contract type of the token (e.g., BurnMintToken, StandardERC20).
	Type cldf.ContractType `json:"type"`

	// PoolType is the type of the token pool that will be deployed for this token.
	PoolType cldf.ContractType `json:"poolType"`

	// PoolAllowList is a list of addresses allowed to interact with the pool.
	PoolAllowList []common.Address `json:"poolAllowList,omitempty"`

	// AcceptLiquidity defines whether the pool should accept liquidity.
	AcceptLiquidity *bool `json:"acceptLiquidity,omitempty"`

	// MintTokenForRecipients is a map of recipient address to amount to be transferred or minted
	// and provided the minting role after token deployment.
	MintTokenForRecipients map[common.Address]*big.Int `json:"mintTokenForRecipients,omitempty"`

	// PreMint is the amount of tokens to pre-mint for the token.
	PreMint *big.Int `json:"preMint,omitempty"`

	// CCIPAdmin is the address of the CCIP admin for the token and will have default admin role. This is specifically
	// for BurnMintERC20 token.
	CCIPAdmin common.Address `json:"ccipAdmin,omitempty"`
}

func (c *DeployTokenConfig) Validate() error {
	if c.TokenName == "" {
		return errors.New("token name must be defined")
	}
	if c.TokenDecimals == 0 && c.Type == shared.BurnMintToken {
		return errors.New("token decimals must be defined for BurnMintToken type")
	}
	if c.MaxSupply == nil && c.Type == shared.BurnMintToken {
		return errors.New("max supply must be defined for BurnMintToken type")
	}
	if _, ok := shared.TokenPoolTypes[c.PoolType]; !ok {
		return fmt.Errorf("token pool type not supported %s", c.PoolType)
	}
	if _, ok := shared.TokenTypes[c.Type]; !ok {
		return fmt.Errorf("token type not supported %s", c.Type)
	}

	if c.Type == shared.BurnMintERC20Token && c.MaxSupply != nil && c.PreMint != nil {
		if c.PreMint.Cmp(c.MaxSupply) > 0 {
			return fmt.Errorf("preMint amount %s cannot be greater than max supply %s for BurnMintERC20 token type",
				c.PreMint.String(), c.MaxSupply.String())
		}
	}
	return nil
}

type AddTokensE2EConfig struct {
	// Tokens is a map from token symbol to its E2E configuration.
	Tokens map[shared.TokenSymbol]AddTokenE2EConfig `json:"tokens"`

	// MCMS is the optional TimelockConfig used for governance proposals.
	MCMS *proposalutils.TimelockConfig `json:"mcms,omitempty"`
}

func addTokenE2EPreconditionValidation(e cldf.Environment, config AddTokensE2EConfig) error {
	if len(config.Tokens) == 0 {
		return nil
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for token, cfg := range config.Tokens {
		for chain, poolCfg := range cfg.PoolConfig {
			if err := stateview.ValidateChain(e, state, chain, config.MCMS); err != nil {
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
				if err := poolCfg.DeployPoolConfig.Validate(e.GetContext(), e.BlockChains.EVMChains()[chain], state.MustGetEVMChainState(chain), token); err != nil {
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

func addTokenE2ELogic(env cldf.Environment, config AddTokensE2EConfig) (cldf.ChangesetOutput, error) {
	if len(config.Tokens) == 0 {
		return cldf.ChangesetOutput{}, nil
	}
	// use a clone of env to avoid modifying the original env
	e := env.Clone()
	finalCSOut := &cldf.ChangesetOutput{
		AddressBook: cldf.NewMemoryAddressBook(),
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	for token, cfg := range config.Tokens {
		e.Logger.Infow("starting token addition operations for", "token", token, "chains", maps.Keys(cfg.PoolConfig))
		tokenDeployCfg := make(map[uint64]DeployTokenConfig)
		for chain, poolCfg := range cfg.PoolConfig {
			if poolCfg.TokenDeploymentConfig != nil {
				tokenDeployCfg[chain] = *poolCfg.TokenDeploymentConfig
			}
		}

		//  1. Deploys tokens ( specifically TestTokens) optionally if DeployTokenConfig is provided and
		//     populates the pool deployment configuration for each token.
		if len(tokenDeployCfg) > 0 {
			deployedTokens, ab, err := deployTokens(e, tokenDeployCfg)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			if err := cfg.newDeployTokenPoolConfigAfterTokenDeployment(deployedTokens); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to populate pool deployment configuration: %w", err)
			}
			e.Logger.Infow("deployed token and created pool deployment config", "token", token)
			if err := finalCSOut.AddressBook.Merge(ab); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book for token %s: %w", token, err)
			}
			// populate the configuration for deploying and configuring token pools and token admin registry
			if err := cfg.newConfigurePoolAndTokenAdminRegConfig(e, token, config.MCMS); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to populate configuration for "+
					"deploying and configuring token pools and token admin registry: %w", err)
			}
		}

		if cfg.deployPool != nil {
			//  2. Deploys token pool contracts for each token specified in the config.
			//     If the token deployment config is provided, pool deployment configuration DeployTokenPoolContractsConfig is not required.
			//     It will use the token deployment config to deploy the token and
			//     populate DeployTokenPoolContractsConfig.
			//     If the token deployment config is not provided, pool deployment configuration DeployTokenPoolContractsConfig is mandatory.
			output, err := DeployTokenPoolContractsChangeset(e, *cfg.deployPool)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy token pool for token %s: %w", token, err)
			}
			if err := cldf.MergeChangesetOutput(e, finalCSOut, output); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book for token %s: %w", token, err)
			}
			newAddresses, err := output.AddressBook.Addresses() //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get addresses from address book: %w", err)
			}
			e.Logger.Infow("deployed token pool", "token", token, "addresses", newAddresses)

			// Validate the configure token admin reg config.
			// As we will perform proposing admin, accepting admin and setting pool on same changeset
			// we are only validating the propose admin role.
			if err := cfg.configureTokenAdminReg.Validate(e, true, validateProposeAdminRole); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to validate configure token admin reg config: %w", err)
			}

			// for each chain in pool config, trigger transfer ownership of the deployed poolAddress
			if config.MCMS != nil {
				for chainID, addrMap := range newAddresses {
					var addresses []common.Address
					for addr := range addrMap {
						addresses = append(addresses, common.HexToAddress(addr))
					}

					transferOwnershipProposalOutput, err := commoncs.TransferToMCMSWithTimelockV2(e, commoncs.TransferToMCMSWithTimelockConfig{
						ContractsByChain: map[uint64][]common.Address{
							chainID: addresses,
						},
						MCMSConfig: *config.MCMS,
					})
					if err != nil {
						return cldf.ChangesetOutput{}, fmt.Errorf("failed to run TransferToMCMSWithTimelock on chain with selector %d: %w", chainID, err)
					}

					if err := cldf.MergeChangesetOutput(e, finalCSOut, transferOwnershipProposalOutput); err != nil {
						return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge changeset output after transferring ownership of pool(s): %w", err)
					}
				}
			}
		}

		// 3. Proposes admin rights for the token on the token admin registry
		//
		//		If the token admin is not an external address -
		//		3a. Accepts admin rights for the token on the token admin registry
		//		3b. Sets the pool for the token on the token admin registry
		output, err := ProposeAdminRoleChangeset(e, cfg.configureTokenAdminReg)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to propose admin role for token %s: %w", token, err)
		}
		if err := cldf.MergeChangesetOutput(e, finalCSOut, output); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to changeset output after configuring token admin reg for token %s: %w",
				token, err)
		}
		e.Logger.Infow("proposed admin role", "token", token, "config", cfg.configureTokenAdminReg)

		// find all tokens for which there is no external admin
		// for those tokens, accept the admin role and set the pool
		updatedConfigureTokenAdminReg := TokenAdminRegistryChangesetConfig{
			MCMS:  config.MCMS,
			Pools: make(map[uint64]map[shared.TokenSymbol]TokenPoolInfo),
			// SkipOwnershipValidation is set to true as we are accepting admin role and setting token pool as part of one changeset
			SkipOwnershipValidation: true,
		}
		for chain, poolInfo := range cfg.configureTokenAdminReg.Pools {
			for symbol, info := range poolInfo {
				if info.ExternalAdmin == utils.ZeroAddress {
					if updatedConfigureTokenAdminReg.Pools[chain] == nil {
						updatedConfigureTokenAdminReg.Pools[chain] = make(map[shared.TokenSymbol]TokenPoolInfo)
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
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to accept admin role for token %s: %w", token, err)
		}
		if err := cldf.MergeChangesetOutput(e, finalCSOut, output); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book for token %s: %w", token, err)
		}
		e.Logger.Infow("accepted admin role", "token", token, "config", updatedConfigureTokenAdminReg)
		output, err = SetPoolChangeset(e, updatedConfigureTokenAdminReg)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to set pool for token %s: %w", token, err)
		}
		if err := cldf.MergeChangesetOutput(e, finalCSOut, output); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book for token %s: %w", token, err)
		}
		e.Logger.Infow("set pool", "token", token, "config", updatedConfigureTokenAdminReg)

		//	 4. Configures pools -
		//		   The configuration is set by via the input field ConfigurePools. This is so that we can have flexibility during configuration.
		//		   For ex: deploy chainB token, pool and configure with already deployed chainA tokenPool.
		//		   If config is not provided in the input, it skips this step.
		if len(cfg.ConfigurePools.PoolUpdates) == 0 {
			e.Logger.Infow("only deployment done, the pools haven't been configured")
			continue
		}

		if err := cfg.ConfigurePools.Validate(e); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to validate configure pool config: %w", err)
		}

		output, err = ConfigureTokenPoolContractsChangeset(e, cfg.ConfigurePools)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure token pool for token %s: %w", token, err)
		}
		if err := cldf.MergeChangesetOutput(e, finalCSOut, output); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge changeset output after configuring token pool for token %s: %w", token, err)
		}
		e.Logger.Infow("configured token pool", "token", token)
	}
	// if there are multiple proposals, aggregate them so that we don't have to propose them separately
	if len(finalCSOut.MCMSTimelockProposals) > 1 {
		aggregatedProposals, err := proposalutils.AggregateProposals(
			e, state.EVMMCMSStateByChain(), nil, finalCSOut.MCMSTimelockProposals,
			"Add Tokens E2E", config.MCMS)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to aggregate proposals: %w", err)
		}
		finalCSOut.MCMSTimelockProposals = []mcms.TimelockProposal{*aggregatedProposals}
	}
	return *finalCSOut, nil
}

func deployTokens(e cldf.Environment, tokenDeployCfg map[uint64]DeployTokenConfig) (map[uint64]common.Address, cldf.AddressBook, error) {
	ab := cldf.NewMemoryAddressBook()
	tokenAddresses := make(map[uint64]common.Address) // This will hold the token addresses for each chain.
	for selector, cfg := range tokenDeployCfg {
		switch cfg.Type {
		case shared.BurnMintToken:
			token, err := cldf.DeployContract(e.Logger, e.BlockChains.EVMChains()[selector], ab,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
					tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
						e.BlockChains.EVMChains()[selector].DeployerKey,
						e.BlockChains.EVMChains()[selector].Client,
						cfg.TokenName,
						string(cfg.TokenSymbol),
						cfg.TokenDecimals,
						cfg.MaxSupply,
					)
					return cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
						Address:  tokenAddress,
						Contract: token,
						Tv:       cldf.NewTypeAndVersion(shared.BurnMintToken, deployment.Version1_0_0),
						Tx:       tx,
						Err:      err,
					}
				},
			)
			if err != nil {
				return nil, ab, fmt.Errorf("failed to deploy BurnMintERC677 token "+
					"%s on chain %d: %w", cfg.TokenName, selector, err)
			}
			if err := addMinterAndMintTokenERC677(e, selector, token.Contract, e.BlockChains.EVMChains()[selector].DeployerKey.From,
				new(big.Int).Mul(big.NewInt(1_000), big.NewInt(1_000_000_000))); err != nil {
				return nil, ab, fmt.Errorf("failed to add minter and mint token "+
					"%s on chain %d: %w", cfg.TokenName, selector, err)
			}
			if len(cfg.MintTokenForRecipients) > 0 {
				for recipient, amount := range cfg.MintTokenForRecipients {
					if err := addMinterAndMintTokenERC677(e, selector, token.Contract, recipient,
						amount); err != nil {
						return nil, ab, fmt.Errorf("failed to add minter and mint "+
							"token %s on chain %d: %w", cfg.TokenName, selector, err)
					}
				}
			}

			tokenAddresses[selector] = token.Address
		case shared.ERC20Token:
			token, err := cldf.DeployContract(e.Logger, e.BlockChains.EVMChains()[selector], ab,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*erc20.ERC20] {
					tokenAddress, tx, token, err := erc20.DeployERC20(
						e.BlockChains.EVMChains()[selector].DeployerKey,
						e.BlockChains.EVMChains()[selector].Client,
						cfg.TokenName,
						string(cfg.TokenSymbol),
					)
					return cldf.ContractDeploy[*erc20.ERC20]{
						Address:  tokenAddress,
						Contract: token,
						Tv:       cldf.NewTypeAndVersion(shared.ERC20Token, deployment.Version1_0_0),
						Tx:       tx,
						Err:      err,
					}
				},
			)

			if err != nil {
				return nil, ab, fmt.Errorf("failed to deploy ERC20 token %s on chain %d: %w", cfg.TokenName, selector, err)
			}
			tokenAddresses[selector] = token.Address
		case shared.ERC677Token:
			token, err := cldf.DeployContract(e.Logger, e.BlockChains.EVMChains()[selector], ab,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*erc677.ERC677] {
					tokenAddress, tx, token, err := erc677.DeployERC677(
						e.BlockChains.EVMChains()[selector].DeployerKey,
						e.BlockChains.EVMChains()[selector].Client,
						cfg.TokenName,
						string(cfg.TokenSymbol),
					)
					return cldf.ContractDeploy[*erc677.ERC677]{
						Address:  tokenAddress,
						Contract: token,
						Tv:       cldf.NewTypeAndVersion(shared.ERC677Token, deployment.Version1_0_0),
						Tx:       tx,
						Err:      err,
					}
				},
			)
			if err != nil {
				return nil, ab, fmt.Errorf("failed to deploy ERC677 token %s on chain %d: %w", cfg.TokenName, selector, err)
			}
			tokenAddresses[selector] = token.Address
		case shared.ERC677TokenHelper:
			token, err := cldf.DeployContract(e.Logger, e.BlockChains.EVMChains()[selector], ab,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*burn_mint_erc20_with_drip.BurnMintERC20] {
					tokenAddress, tx, token, err := burn_mint_erc20_with_drip.DeployBurnMintERC20(
						e.BlockChains.EVMChains()[selector].DeployerKey,
						e.BlockChains.EVMChains()[selector].Client,
						cfg.TokenName,
						string(cfg.TokenSymbol),
						18,
						big.NewInt(0),
						big.NewInt(0),
					)
					return cldf.ContractDeploy[*burn_mint_erc20_with_drip.BurnMintERC20]{
						Address:  tokenAddress,
						Contract: token,
						Tv:       cldf.NewTypeAndVersion(shared.ERC677TokenHelper, deployment.Version1_0_0),
						Tx:       tx,
						Err:      err,
					}
				},
			)
			if err != nil {
				return nil, ab, fmt.Errorf("failed to deploy ERC677 token %s on chain %d: %w", cfg.TokenName, selector, err)
			}

			if err := addMinterAndMintTokenERC677Helper(e, selector, token.Contract, e.BlockChains.EVMChains()[selector].DeployerKey.From,
				new(big.Int).Mul(big.NewInt(1_000), big.NewInt(1_000_000_000))); err != nil {
				return nil, ab, fmt.Errorf("failed to add minter and mint token "+
					"%s on chain %d: %w", cfg.TokenName, selector, err)
			}
			if len(cfg.MintTokenForRecipients) > 0 {
				for recipient, amount := range cfg.MintTokenForRecipients {
					if err := addMinterAndMintTokenERC677Helper(e, selector, token.Contract, recipient,
						amount); err != nil {
						return nil, ab, fmt.Errorf("failed to add minter and mint "+
							"token %s on chain %d: %w", cfg.TokenName, selector, err)
					}
				}
			}
			tokenAddresses[selector] = token.Address
		case shared.BurnMintERC20Token:
			token, err := cldf.DeployContract(e.Logger, e.BlockChains.EVMChains()[selector], ab,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*burn_mint_erc20.BurnMintERC20] {
					if cfg.MaxSupply == nil {
						cfg.MaxSupply = big.NewInt(0)
					}
					if cfg.PreMint == nil {
						cfg.PreMint = big.NewInt(0)
					}
					tokenAddress, tx, token, err := burn_mint_erc20.DeployBurnMintERC20(
						e.BlockChains.EVMChains()[selector].DeployerKey,
						e.BlockChains.EVMChains()[selector].Client,
						cfg.TokenName,
						string(cfg.TokenSymbol),
						cfg.TokenDecimals,
						cfg.MaxSupply,
						cfg.PreMint,
					)
					return cldf.ContractDeploy[*burn_mint_erc20.BurnMintERC20]{
						Address:  tokenAddress,
						Contract: token,
						Tv:       cldf.NewTypeAndVersion(shared.BurnMintERC20Token, deployment.Version1_0_0),
						Tx:       tx,
						Err:      err,
					}
				},
			)
			if err != nil {
				return nil, ab, fmt.Errorf("failed to deploy BurnMintERC20Token "+
					"%s on chain %d: %w", cfg.TokenName, selector, err)
			}

			if err := setCCIPAdminForBurnMintERC20Token(e, selector, token.Contract, cfg.CCIPAdmin); err != nil {
				return nil, ab, fmt.Errorf("failed to set CCIP admin for %s on chain %d: %w",
					cfg.CCIPAdmin, selector, err)
			}

			if err := grantDefaultAdminRoleForBurnMintERC20Token(e, selector, token.Contract, cfg.CCIPAdmin); err != nil {
				return nil, ab, fmt.Errorf("failed to grant default admin role for BurnMintERC20 token %s on chain %d: %w",
					cfg.TokenName, selector, err)
			}

			tokenAddresses[selector] = token.Address
		default:
			return nil, ab, fmt.Errorf("unsupported token %s type %s for deployment on chain %d", cfg.TokenName, cfg.Type, selector)
		}
	}

	return tokenAddresses, ab, nil
}

// addMinterAndMintTokenERC677 adds the minter role to the recipient and mints the specified amount of tokens to the recipient's address.
func addMinterAndMintTokenERC677(env cldf.Environment, selector uint64, token *burn_mint_erc677.BurnMintERC677, recipient common.Address, amount *big.Int) error {
	return addMinterAndMintTokenHelper(env, selector, token, recipient, amount)
}

// addMinterAndMintTokenERC677Helper adds the minter role to the recipient and mints the specified amount of tokens to the recipient's address.
func addMinterAndMintTokenERC677Helper(env cldf.Environment, selector uint64, token *burn_mint_erc20_with_drip.BurnMintERC20, recipient common.Address, amount *big.Int) error {
	baseToken, err := burn_mint_erc677.NewBurnMintERC677(token.Address(), env.BlockChains.EVMChains()[selector].Client)
	if err != nil {
		return fmt.Errorf("failed to cast helper to base token: %w", err)
	}
	return addMinterAndMintTokenHelper(env, selector, baseToken, recipient, amount)
}

func addMinterAndMintTokenHelper(env cldf.Environment, selector uint64, token *burn_mint_erc677.BurnMintERC677, recipient common.Address, amount *big.Int) error {
	deployerKey := env.BlockChains.EVMChains()[selector].DeployerKey
	chain := env.BlockChains.EVMChains()[selector]
	ctx := env.GetContext()

	// Grant minter role to the given address
	err := addMinterForERC677Token(env, chain, token.Address(), recipient)
	if err != nil {
		return fmt.Errorf("failed to grant Mint & Burn role to %s on chain %d: %w", recipient.Hex(), selector, err)
	}

	// Mint tokens to the given address and verify the balance
	tx, err := token.Mint(deployerKey, recipient, amount)
	if err != nil {
		return fmt.Errorf("failed to mint %s tokens to %s on chain %d: %w",
			token.Address().Hex(), recipient.Hex(), selector, err)
	}
	if _, err := env.BlockChains.EVMChains()[selector].Confirm(tx); err != nil {
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

// addMinterAndBurnerForBurnMintERC20Token adds the burner and minter role to the specified address.
func addMinterAndBurnerForBurnMintERC20Token(env cldf.Environment, selector uint64, tokenAddress common.Address, poolAddress common.Address) error {
	return addMinterAndBurnerForBurnMintERC20TokenHelper(env, selector, tokenAddress, poolAddress)
}

// addMinterAndBurnerForBurnMintERC20TokenHelper is a helper function that adds the minter and burner role to the specified address for BurnMintERC20 token.
func addMinterAndBurnerForBurnMintERC20TokenHelper(env cldf.Environment, selector uint64, tokenAddress common.Address, poolAddress common.Address) error {
	deployerKey := env.BlockChains.EVMChains()[selector].DeployerKey
	ctx := env.GetContext()

	token, err := burn_mint_erc20.NewBurnMintERC20(tokenAddress, env.BlockChains.EVMChains()[selector].Client)
	if err != nil {
		return fmt.Errorf("failed to connect address %s with burn_mint_erc20 bindings: %w", tokenAddress.Hex(), err)
	}

	mintRole, err := token.MINTERROLE(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get mint role of token %s on chain %d: %w", token.Address().Hex(), selector, err)
	}

	hasMintRole, err := token.HasRole(&bind.CallOpts{Context: ctx}, mintRole, poolAddress)
	if err != nil {
		return fmt.Errorf("failed to check if pool has mint role for token %s on chain %d: %w", token.Address().Hex(), selector, err)
	}

	burnRole, err := token.BURNERROLE(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get burn role of token %s on chain %d: %w", token.Address().Hex(), selector, err)
	}

	hasBurnRole, err := token.HasRole(&bind.CallOpts{Context: ctx}, burnRole, poolAddress)
	if err != nil {
		return fmt.Errorf("failed to check if pool has burn role for token %s on chain %d: %w", token.Address().Hex(), selector, err)
	}

	if hasMintRole && hasBurnRole {
		env.Logger.Infow("Pool already has mint and burn role for token", "Token", token.Address().Hex(), "Selector", selector)
	} else {
		tx, err := token.GrantMintAndBurnRoles(deployerKey, poolAddress)

		if err != nil {
			return fmt.Errorf("failed to grant mint and burn role to %s on chain %d: %w", poolAddress.Hex(), selector, err)
		}

		if _, err := env.BlockChains.EVMChains()[selector].Confirm(tx); err != nil {
			return fmt.Errorf("failed to wait for transaction %s on chain %d: %w", tx.Hash().Hex(), selector, err)
		}

		env.Logger.Infow("Transaction granting mint and burn role mined successfully",
			"Hash", tx.Hash().Hex(), "Selector", selector)
	}

	return nil
}

// setCCIPAdminForBurnMintERC20Token sets the CCIP admin for the BurnMintERC20 token. Expectation is that the CCIP admin is
// timelock or customer (multi-sig) address.
func setCCIPAdminForBurnMintERC20Token(env cldf.Environment, selector uint64, token *burn_mint_erc20.BurnMintERC20, address common.Address) error {
	deployerKey := env.BlockChains.EVMChains()[selector].DeployerKey

	tx, err := token.SetCCIPAdmin(deployerKey, address)
	if err != nil {
		return fmt.Errorf("failed to set CCIP admin for token %s on chain %d: %w", token.Address().Hex(), selector, err)
	}
	if _, err := env.BlockChains.EVMChains()[selector].Confirm(tx); err != nil {
		return fmt.Errorf("failed to wait for transaction %s on chain %d: %w", tx.Hash().Hex(), selector, err)
	}
	env.Logger.Infow("Transaction setting CCIP admin mined successfully",
		"Hash", tx.Hash().Hex(), "Selector", selector)

	return nil
}

// grantDefaultAdminRoleForBurnMintERC20Token grants the default admin role for the BurnMintERC20 token to the specified address
// which is expected to be the pool address.
func grantDefaultAdminRoleForBurnMintERC20Token(env cldf.Environment, selector uint64, token *burn_mint_erc20.BurnMintERC20, address common.Address) error {
	deployerKey := env.BlockChains.EVMChains()[selector].DeployerKey
	ctx := env.GetContext()

	adminRole, err := token.DEFAULTADMINROLE(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get default admin role on token %s on chain %d: %w", token.Address().Hex(), selector, err)
	}

	hasRole, err := token.HasRole(&bind.CallOpts{Context: ctx}, adminRole, address)
	if err != nil {
		return fmt.Errorf("failed to check if address %s has default admin role for token %s on chain %d: %w", address, token.Address().Hex(), selector, err)
	}

	if hasRole {
		env.Logger.Infow("Pool already has default admin role for token", "Token", token.Address().Hex(), "Pool", address, "Selector", selector)
	} else {
		tx, err := token.GrantRole(deployerKey, adminRole, address)

		if err != nil {
			return fmt.Errorf("failed to grant default admin role for token %s on chain %d: %w", token.Address().Hex(), selector, err)
		}
		if _, err := env.BlockChains.EVMChains()[selector].Confirm(tx); err != nil {
			return fmt.Errorf("failed to wait for transaction %s on chain %d: %w", tx.Hash().Hex(), selector, err)
		}
		env.Logger.Infow("Transaction granting default admin role mined successfully",
			"Hash", tx.Hash().Hex(), "Selector", selector)
	}

	return nil
}

func addMinterForERC677Token(env cldf.Environment, chain cldf_evm.Chain, tokenAddress common.Address, poolAddress common.Address) error {
	_, err := operations.ExecuteOperation(env.OperationsBundle, ccipops.GrantMintAndBurnRolesERC677Op, chain, opsutil.EVMCallInput[common.Address]{
		Address:       tokenAddress,
		ChainSelector: chain.Selector,
		CallInput:     poolAddress,
	})

	if err != nil {
		return fmt.Errorf("failed to grant mint and burn roles: %w", err)
	}

	return nil
}
