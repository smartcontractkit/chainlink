package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/token_admin_registry"
)

var _ deployment.ChangeSet[TokenAdminRegistryChangesetConfig] = AcceptAdminRoleChangeset

func validateAcceptAdminRole(
	config token_admin_registry.TokenAdminRegistryTokenConfig,
	sender common.Address,
	externalAdmin common.Address,
	symbol TokenSymbol,
	chain deployment.Chain,
) error {
	// We must be the pending administrator
	if config.PendingAdministrator != sender {
		return fmt.Errorf("unable to accept admin role for %s token on %s: %s is not the pending administrator (%s)", symbol, chain, sender, config.PendingAdministrator)
	}
	return nil
}

// AcceptAdminRoleChangeset accepts admin rights for tokens on the token admin registry.
func AcceptAdminRoleChangeset(env deployment.Environment, c TokenAdminRegistryChangesetConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(env, false, validateAcceptAdminRole); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid TokenAdminRegistryChangesetConfig: %w", err)
	}
	state, err := LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := NewDeployerGroup(env, state, c.MCMS).WithDeploymentContext("accept admin role for tokens on token admin registries")

	for chainSelector, tokenSymbolToPoolInfo := range c.Pools {
		chain := env.Chains[chainSelector]
		chainState := state.Chains[chainSelector]
		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get deployer for %s", chain)
		}
		for symbol, poolInfo := range tokenSymbolToPoolInfo {
			_, tokenAddress, err := poolInfo.GetPoolAndTokenAddress(env.GetContext(), symbol, chain, chainState)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to get state of %s token on chain %s: %w", symbol, chain, err)
			}
			_, err = chainState.TokenAdminRegistry.AcceptAdminRole(opts, tokenAddress)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to create acceptAdminRole transaction for %s on %s registry: %w", symbol, chain, err)
			}
		}
	}

	return deployerGroup.Enact()
}
