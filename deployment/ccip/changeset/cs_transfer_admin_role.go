package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/token_admin_registry"
)

var _ deployment.ChangeSet[TokenAdminRegistryChangesetConfig] = TransferAdminRoleChangeset

func validateTransferAdminRole(
	config token_admin_registry.TokenAdminRegistryTokenConfig,
	sender common.Address,
	externalAdmin common.Address,
	symbol TokenSymbol,
	chain deployment.Chain,
) error {
	if externalAdmin == utils.ZeroAddress {
		return errors.New("external admin must be defined")
	}
	// We must be the administrator
	if config.Administrator != sender {
		return fmt.Errorf("unable to transfer admin role for %s token on %s: %s is not the administrator (%s)", symbol, chain, sender, config.Administrator)
	}
	return nil
}

// TransferAdminRoleChangeset transfers the admin role for tokens on the token admin registry to 3rd parties.
func TransferAdminRoleChangeset(env deployment.Environment, c TokenAdminRegistryChangesetConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(env, false, validateTransferAdminRole); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid TokenAdminRegistryChangesetConfig: %w", err)
	}
	state, err := LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := NewDeployerGroup(env, state, c.MCMS).WithDeploymentContext("transfer admin role for tokens on token admin registries")

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
			_, err = chainState.TokenAdminRegistry.TransferAdminRole(opts, tokenAddress, poolInfo.ExternalAdmin)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transferAdminRole transaction for %s on %s registry: %w", symbol, chain, err)
			}
		}
	}

	return deployerGroup.Enact()
}
