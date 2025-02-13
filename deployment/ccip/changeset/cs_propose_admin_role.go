package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/token_admin_registry"
)

var _ deployment.ChangeSet[TokenAdminRegistryChangesetConfig] = ProposeAdminRoleChangeset

func validateProposeAdminRole(
	config token_admin_registry.TokenAdminRegistryTokenConfig,
	sender common.Address,
	externalAdmin common.Address,
	symbol TokenSymbol,
	chain deployment.Chain,
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
func ProposeAdminRoleChangeset(env deployment.Environment, c TokenAdminRegistryChangesetConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(env, true, validateProposeAdminRole); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid TokenAdminRegistryChangesetConfig: %w", err)
	}
	state, err := LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := NewDeployerGroup(env, state, c.MCMS).WithDeploymentContext("propose admin role for tokens on token admin registries")

	for chainSelector, tokenSymbolToPoolInfo := range c.Pools {
		chain := env.Chains[chainSelector]
		chainState := state.Chains[chainSelector]
		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get deployer for %s", chain)
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
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to get state of %s token on chain %s: %w", symbol, chain, err)
			}
			_, err = chainState.TokenAdminRegistry.ProposeAdministrator(opts, tokenAddress, desiredAdmin)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to create proposeAdministrator transaction for %s on %s registry: %w", symbol, chain, err)
			}
		}
	}

	return deployerGroup.Enact()
}
