package v1_5_1

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
)

// AddTokensE2E is a changeset that deploys and configures token pools for multiple tokens across multiple chains in a single changeset.
// AddTokensE2E does the following:
//
//  1. Deploys token pool contracts for each token specified in the config
//
//  2. Configures pools -
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
// 3. Proposes admin rights for the token on the token admin registry
// If the token admin is not an external address -
// 4. Accepts admin rights for the token on the token admin registry
// 5. Sets the pool for the token on the token admin registry
var AddTokensE2E = deployment.CreateChangeSet(addTokenE2ELogic, addTokenE2EPreconditionValidation)

type AddTokensE2EConfig struct {
	Tokens map[changeset.TokenSymbol]AddTokenE2EConfig
}

type AddTokenE2EConfig struct {
	DeploymentConfig       DeployTokenPoolContractsConfig
	ConfigurePools         ConfigureTokenPoolContractsConfig
	ConfigureTokenAdminReg changeset.TokenAdminRegistryChangesetConfig
}

func addTokenE2EPreconditionValidation(e deployment.Environment, config AddTokensE2EConfig) error {
	if len(config.Tokens) == 0 {
		return nil
	}
	for token, cfg := range config.Tokens {
		if token != cfg.DeploymentConfig.TokenSymbol {
			return fmt.Errorf("token symbol %s in pool config does not match token %s", cfg.DeploymentConfig.TokenSymbol, token)
		}
		if err := cfg.DeploymentConfig.Validate(e); err != nil {
			return fmt.Errorf("failed to validate token pool for token %s: %w", token, err)
		}
		if err := cfg.ConfigurePools.Validate(e); err != nil {
			return fmt.Errorf("failed to validate token pool configuration for token %s: %w", token, err)
		}
		// skip the registry ownership check by passing nil for the registryConfig check function
		// since we would like to perform propose admin, accept admin, and set pool in the same changeset
		if err := cfg.ConfigureTokenAdminReg.Validate(e, true, nil); err != nil {
			return fmt.Errorf("failed to validate token admin registry for token %s: %w", token, err)
		}
	}
	return nil
}

func addTokenE2ELogic(e deployment.Environment, config AddTokensE2EConfig) (deployment.ChangesetOutput, error) {
	if len(config.Tokens) == 0 {
		return deployment.ChangesetOutput{}, nil
	}

	finalCSOut := &deployment.ChangesetOutput{
		AddressBook: deployment.NewMemoryAddressBook(),
	}
	for _, cfg := range config.Tokens {
		output, err := DeployTokenPoolContractsChangeset(e, cfg.DeploymentConfig)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy token pool for token %s: %w", cfg.DeploymentConfig.TokenSymbol, err)
		}
		if err := globals.MergeChangesetOutput(finalCSOut, output); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge address book for token %s: %w", cfg.DeploymentConfig.TokenSymbol, err)
		}
		e.Logger.Infow("deployed token pool", "token", cfg.DeploymentConfig.TokenSymbol, "addresses", output.AddressBook)
		output, err = ConfigureTokenPoolContractsChangeset(e, cfg.ConfigurePools)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to configure token pool for token %s: %w", cfg.DeploymentConfig.TokenSymbol, err)
		}
		if err := globals.MergeChangesetOutput(finalCSOut, output); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge address book for token %s: %w", cfg.DeploymentConfig.TokenSymbol, err)
		}
		e.Logger.Infow("configured token pool", "token", cfg.DeploymentConfig.TokenSymbol)

		output, err = ProposeAdminRoleChangeset(e, cfg.ConfigureTokenAdminReg)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to propose admin role for token %s: %w", cfg.DeploymentConfig.TokenSymbol, err)
		}
		if err := globals.MergeChangesetOutput(finalCSOut, output); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge address book for token %s: %w", cfg.DeploymentConfig.TokenSymbol, err)
		}
		e.Logger.Infow("proposed admin role", "token", cfg.DeploymentConfig.TokenSymbol, "config", cfg.ConfigureTokenAdminReg)
		// find all tokens for which there is no external admin
		updatedConfigureTokenAdminReg := changeset.TokenAdminRegistryChangesetConfig{
			MCMS:  cfg.ConfigureTokenAdminReg.MCMS,
			Pools: make(map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo),
		}
		for chain, poolInfo := range cfg.ConfigureTokenAdminReg.Pools {
			for symbol, info := range poolInfo {
				if info.ExternalAdmin == utils.ZeroAddress {
					if updatedConfigureTokenAdminReg.Pools[chain] == nil {
						updatedConfigureTokenAdminReg.Pools[chain] = make(map[changeset.TokenSymbol]changeset.TokenPoolInfo)
					}
					updatedConfigureTokenAdminReg.Pools[chain][symbol] = info
				}
			}
		}
		output, err = AcceptAdminRoleChangeset(e, updatedConfigureTokenAdminReg)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to accept admin role for token %s: %w", cfg.DeploymentConfig.TokenSymbol, err)
		}
		if err := globals.MergeChangesetOutput(finalCSOut, output); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge address book for token %s: %w", cfg.DeploymentConfig.TokenSymbol, err)
		}
		e.Logger.Infow("accepted admin role", "token", cfg.DeploymentConfig.TokenSymbol, "config", updatedConfigureTokenAdminReg)
		output, err = SetPoolChangeset(e, updatedConfigureTokenAdminReg)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to set pool for token %s: %w", cfg.DeploymentConfig.TokenSymbol, err)
		}
		if err := globals.MergeChangesetOutput(finalCSOut, output); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge address book for token %s: %w", cfg.DeploymentConfig.TokenSymbol, err)
		}
		e.Logger.Infow("set pool", "token", cfg.DeploymentConfig.TokenSymbol, "config", updatedConfigureTokenAdminReg)
	}

	return *finalCSOut, nil
}
