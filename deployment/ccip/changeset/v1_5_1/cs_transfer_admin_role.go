package v1_5_1

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

var _ deployment.ChangeSet[changeset.TokenAdminRegistryChangesetConfig] = TransferAdminRoleChangeset

func validateTransferAdminRole(
	config token_admin_registry.TokenAdminRegistryTokenConfig,
	sender common.Address,
	externalAdmin common.Address,
	symbol changeset.TokenSymbol,
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
func TransferAdminRoleChangeset(env deployment.Environment, c changeset.TokenAdminRegistryChangesetConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(env, false, validateTransferAdminRole); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid TokenAdminRegistryChangesetConfig: %w", err)
	}
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := changeset.NewDeployerGroup(env, state, c.MCMS).WithDeploymentContext("transfer admin role for tokens on token admin registries")

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
