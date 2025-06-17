package v1_5_1

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

// ProposeAdminRoleChangeset is a changeset that proposes admin rights for tokens on the token admin registry.
// To be able to propose admin rights, the caller must own the token admin registry and the token must not already have an administrator.
// If you want to propose admin role for an external address, you can set the ExternalAdmin field in the TokenPoolInfo within TokenAdminRegistryChangesetConfig.
var _ cldf.ChangeSet[TokenAdminRegistryChangesetConfig] = ProposeAdminRoleChangeset

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
