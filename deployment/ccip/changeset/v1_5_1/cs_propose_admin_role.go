package v1_5_1

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

// ProposeAdminRoleChangeset is a changeset that proposes admin rights for tokens on the token admin registry.
// To be able to propose admin rights, the caller must own the token admin registry and the token must not already have an administrator.
// If you want to propose admin role for an external address, you can set the ExternalAdmin field in the TokenPoolInfo within TokenAdminRegistryChangesetConfig.
var _ cldf.ChangeSet[TokenAdminRegistryChangesetConfig] = ProposeAdminRoleChangeset

var (
	// ProposeAdminRoleChangesetV2 is a changeset that proposes administrator roles for tokens
	// on the token admin registry. It takes array of tokenAdminInfo by chain.
	// The tokenAdminInfo expect to have tokenAddress and adminAddress It uses the proposeAdminRoleLogic function to execute the
	// proposal and proposeAdminRolePrecondition to validate preconditions before execution.
	// This changeset is intended for use when assigning admin rights to tokens, ensuring that
	// the caller owns the registry and that the token does not already have an administrator.
	ProposeAdminRoleChangesetV2 = cldf.CreateChangeSet(proposeAdminRoleLogic, proposeAdminRolePrecondition)
)

type ProposeAdminRoleConfig struct {
	ProposeAdminByChain map[uint64][]TokenAdminInfo `json:"proposeAdminByChain"`
	// MCMS defines the delay to use for Timelock (if absent, the changeset will attempt to use the deployer key).
	MCMS *proposalutils.TimelockConfig
	// OverridePendingAdmin allows overriding existing pending administrators if set to true.
	// Use with caution as this will replace any existing pending admin proposals.
	OverridePendingAdmin bool `json:"overridePendingAdmin"`
}

type TokenAdminInfo struct {
	TokenAddress common.Address `json:"tokenAddress"`
	AdminAddress common.Address `json:"adminAddress"`
}

func proposeAdminRolePrecondition(e cldf.Environment, cfg ProposeAdminRoleConfig) error {
	if len(cfg.ProposeAdminByChain) == 0 {
		return errors.New("at least one chain with token admin info must be specified in ProposeAdminRoleConfig")
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, tokenAdminInfo := range cfg.ProposeAdminByChain {
		chain := e.BlockChains.EVMChains()[chainSelector]
		chainState, ok := state.EVMChainState(chainSelector)
		if !ok {
			return fmt.Errorf("%s does not exist in state", chain)
		}
		if tokenAdminRegistry := chainState.TokenAdminRegistry; tokenAdminRegistry == nil {
			return fmt.Errorf("missing tokenAdminRegistry on %s", chain)
		}
		if err := stateview.ValidateChain(e, state, chainSelector, cfg.MCMS); err != nil {
			return fmt.Errorf("failed to validate chain %d: %w", chainSelector, err)
		}
		if len(tokenAdminInfo) == 0 {
			return fmt.Errorf("no token admin info provided for chain selector %d", chainSelector)
		}

		if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, chain.DeployerKey.From, chainState.Timelock.Address(), chainState.TokenAdminRegistry); err != nil {
			return fmt.Errorf("token admin registry failed ownership validation on %s: %w", chain, err)
		}

		// Track seen token addresses to prevent duplicates within the same chain
		seenTokens := make(map[common.Address]bool)

		for _, info := range tokenAdminInfo {
			if info.TokenAddress == utils.ZeroAddress {
				return fmt.Errorf("token address cannot be zero for propose admin role on chain %d", chainSelector)
			}
			if info.AdminAddress == utils.ZeroAddress {
				return fmt.Errorf("admin address cannot be zero for propose admin role on chain %d for token %s", chainSelector, info.TokenAddress.Hex())
			}

			// Check for duplicate token addresses within the same chain
			if seenTokens[info.TokenAddress] {
				return fmt.Errorf("duplicate token address %s found for chain %d", info.TokenAddress.Hex(), chainSelector)
			}
			seenTokens[info.TokenAddress] = true

			// Validate that admin address is not the same as token address (would be weird)
			if info.AdminAddress == info.TokenAddress {
				return fmt.Errorf("admin address cannot be the same as token address %s on chain %d", info.TokenAddress.Hex(), chainSelector)
			}

			tokenConfig, err := chainState.TokenAdminRegistry.GetTokenConfig(&bind.CallOpts{Context: e.GetContext()}, info.TokenAddress)
			if err != nil {
				return fmt.Errorf("failed to get config of token with address %s from registry on %s: %w", info.TokenAddress.Hex(), chain, err)
			}

			// Check if token already has an administrator
			if tokenConfig.Administrator != utils.ZeroAddress {
				return fmt.Errorf("cannot propose admin role for token %s on %s: token already has an administrator (%s)", info.TokenAddress.Hex(), chain, tokenConfig.Administrator.Hex())
			}

			// Check if there's already a pending administrator for this token
			if tokenConfig.PendingAdministrator != utils.ZeroAddress {
				if !cfg.OverridePendingAdmin {
					return fmt.Errorf("cannot propose admin role for token %s on %s: token already has a pending administrator (%s). Set OverridePendingAdmin=true to override", info.TokenAddress.Hex(), chain, tokenConfig.PendingAdministrator.Hex())
				}
				// Log warning when overriding pending admin
				e.Logger.Warnf("Overriding existing pending administrator %s for token %s on chain %s with new admin %s",
					tokenConfig.PendingAdministrator.Hex(), info.TokenAddress.Hex(), chain, info.AdminAddress.Hex())
			}

			// Validate that we're not proposing the same admin that's already pending or active
			if tokenConfig.PendingAdministrator == info.AdminAddress {
				return fmt.Errorf("admin address %s is already the pending administrator for token %s on %s", info.AdminAddress.Hex(), info.TokenAddress.Hex(), chain)
			}
			if tokenConfig.Administrator == info.AdminAddress {
				return fmt.Errorf("admin address %s is already the administrator for token %s on %s", info.AdminAddress.Hex(), info.TokenAddress.Hex(), chain)
			}
		}
	}

	return nil
}

func proposeAdminRoleLogic(e cldf.Environment, cfg ProposeAdminRoleConfig) (cldf.ChangesetOutput, error) {
	if len(cfg.ProposeAdminByChain) == 0 {
		return cldf.ChangesetOutput{}, nil
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	deployerGroup := deployergroup.NewDeployerGroup(e, state, cfg.MCMS).WithDeploymentContext("propose admin role for tokens on token admin registries")

	// Count total operations for logging
	totalOperations := 0
	for _, tokenAdminInfo := range cfg.ProposeAdminByChain {
		totalOperations += len(tokenAdminInfo)
	}
	e.Logger.Infof("Proposing admin roles for %d tokens across %d chains", totalOperations, len(cfg.ProposeAdminByChain))

	for chainSelector, tokenAdminInfo := range cfg.ProposeAdminByChain {
		chain := e.BlockChains.EVMChains()[chainSelector]
		chainState := state.Chains[chainSelector]
		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for chain %s (selector: %d): %w", chain, chainSelector, err)
		}

		e.Logger.Infof("Proposing admin roles for %d tokens on chain %s (selector: %d)", len(tokenAdminInfo), chain, chainSelector)

		for i, info := range tokenAdminInfo {
			e.Logger.Debugf("Proposing admin %s for token %s on chain %s (%d/%d)",
				info.AdminAddress.Hex(), info.TokenAddress.Hex(), chain, i+1, len(tokenAdminInfo))

			_, err = chainState.TokenAdminRegistry.ProposeAdministrator(opts, info.TokenAddress, info.AdminAddress)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create proposeAdministrator transaction for token %s on chain %s (selector: %d): %w",
					info.TokenAddress.Hex(), chain, chainSelector, err)
			}
		}
	}

	return deployerGroup.Enact()
}

func validateProposeAdminRole(
	config token_admin_registry.TokenAdminRegistryTokenConfig,
	sender common.Address,
	externalAdmin common.Address,
	symbol shared.TokenSymbol,
	chain cldf_evm.Chain,
) error {
	// To propose ourselves as admin of the token, two things must be true.
	//   1. We own the token admin registry
	//   2. An admin does not exist exist yet
	// We've already validated that we own the registry during ValidateOwnership, so we only need to check the 2nd condition
	if config.Administrator != utils.ZeroAddress {
		return fmt.Errorf("unable to propose %s as admin of %s token on %s: token already has an administrator (%s)", sender, symbol, chain, config.Administrator)
	}
	return nil
}

// ProposeAdminRoleChangeset proposes admin rights for tokens on the token admin registry.
func ProposeAdminRoleChangeset(env cldf.Environment, c TokenAdminRegistryChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(env, true, validateProposeAdminRole); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid TokenAdminRegistryChangesetConfig: %w", err)
	}
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := deployergroup.NewDeployerGroup(env, state, c.MCMS).WithDeploymentContext("propose admin role for tokens on token admin registries")

	for chainSelector, tokenSymbolToPoolInfo := range c.Pools {
		chain := env.BlockChains.EVMChains()[chainSelector]
		chainState := state.Chains[chainSelector]
		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for %s", chain)
		}
		desiredAdmin := chainState.Timelock.Address()
		if c.MCMS == nil {
			desiredAdmin = chain.DeployerKey.From
		}
		for symbol, poolInfo := range tokenSymbolToPoolInfo {
			if poolInfo.ExternalAdmin != utils.ZeroAddress {
				desiredAdmin = poolInfo.ExternalAdmin
			}
			_, tokenAddress, err := poolInfo.GetPoolAndTokenAddress(env.GetContext(), symbol, chain, chainState)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get state of %s token on chain %s: %w", symbol, chain, err)
			}
			_, err = chainState.TokenAdminRegistry.ProposeAdministrator(opts, tokenAddress, desiredAdmin)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create proposeAdministrator transaction for %s on %s registry: %w", symbol, chain, err)
			}
		}
	}

	return deployerGroup.Enact()
}
