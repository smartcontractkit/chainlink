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

var _ cldf.ChangeSet[TokenAdminRegistryChangesetConfig] = TransferAdminRoleChangeset

var (
	// TransferAdminRoleChangesetV2 is a changeset that transfers administrator roles for tokens
	// on the token admin registry. It takes array of TokenAdminInfo by chain.
	// The TokenAdminInfo expects to have TokenAddress and AdminAddress. It uses the transferAdminRoleLogic function to execute the
	// transfer and transferAdminRolePrecondition to validate preconditions before execution.
	// This changeset is intended for use when transferring admin rights to other addresses, ensuring that
	// the caller is the current administrator and that the new admin is valid.
	TransferAdminRoleChangesetV2 = cldf.CreateChangeSet(transferAdminRoleLogic, transferAdminRolePrecondition)
)

type TransferAdminRoleConfig struct {
	TransferAdminByChain map[uint64][]TokenAdminInfo `json:"transferAdminByChain"`
	// MCMS defines the delay to use for Timelock (if absent, the changeset will attempt to use the deployer key).
	MCMS *proposalutils.TimelockConfig
}

func validateTransferAdminRole(
	config token_admin_registry.TokenAdminRegistryTokenConfig,
	sender common.Address,
	externalAdmin common.Address,
	symbol shared.TokenSymbol,
	chain cldf_evm.Chain,
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

func transferAdminRolePrecondition(e cldf.Environment, cfg TransferAdminRoleConfig) error {
	if len(cfg.TransferAdminByChain) == 0 {
		return errors.New("at least one chain with token admin info must be specified in TransferAdminRoleConfig")
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, tokenAdminInfo := range cfg.TransferAdminByChain {
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

		// Determine the current admin based on MCMS configuration
		currentAdmin := chain.DeployerKey.From
		if cfg.MCMS != nil {
			currentAdmin = chainState.Timelock.Address()
		}

		if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, currentAdmin, chainState.Timelock.Address(), chainState.TokenAdminRegistry); err != nil {
			return fmt.Errorf("token admin registry failed ownership validation on %s: %w", chain, err)
		}

		// Track seen token addresses to prevent duplicates within the same chain
		seenTokens := make(map[common.Address]bool)

		for _, info := range tokenAdminInfo {
			if info.TokenAddress == utils.ZeroAddress {
				return fmt.Errorf("token address cannot be zero for transfer admin role on chain %d", chainSelector)
			}
			if info.AdminAddress == utils.ZeroAddress {
				return fmt.Errorf("admin address cannot be zero for transfer admin role on chain %d for token %s", chainSelector, info.TokenAddress.Hex())
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

			// Check if we are the current administrator (required to transfer)
			if tokenConfig.Administrator != currentAdmin {
				return fmt.Errorf("cannot transfer admin role for token %s on %s: current administrator is %s, but we are %s", info.TokenAddress.Hex(), chain, tokenConfig.Administrator.Hex(), currentAdmin.Hex())
			}

			// Validate that we're not transferring to the same admin that's already active
			if tokenConfig.Administrator == info.AdminAddress {
				return fmt.Errorf("admin address %s is already the administrator for token %s on %s", info.AdminAddress.Hex(), info.TokenAddress.Hex(), chain)
			}

			// Check if there's a pending administrator that would be overridden
			if tokenConfig.PendingAdministrator != utils.ZeroAddress {
				e.Logger.Warnf("Token %s on chain %s has a pending administrator %s that will be overridden by this transfer to %s",
					info.TokenAddress.Hex(), chain, tokenConfig.PendingAdministrator.Hex(), info.AdminAddress.Hex())
			}
		}
	}

	return nil
}

func transferAdminRoleLogic(e cldf.Environment, cfg TransferAdminRoleConfig) (cldf.ChangesetOutput, error) {
	if len(cfg.TransferAdminByChain) == 0 {
		return cldf.ChangesetOutput{}, nil
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	deployerGroup := deployergroup.NewDeployerGroup(e, state, cfg.MCMS).WithDeploymentContext("transfer admin role for tokens on token admin registries")

	// Count total operations for logging
	totalOperations := 0
	for _, tokenAdminInfo := range cfg.TransferAdminByChain {
		totalOperations += len(tokenAdminInfo)
	}
	e.Logger.Infof("Transferring admin roles for %d tokens across %d chains", totalOperations, len(cfg.TransferAdminByChain))

	for chainSelector, tokenAdminInfo := range cfg.TransferAdminByChain {
		chain := e.BlockChains.EVMChains()[chainSelector]
		chainState := state.Chains[chainSelector]
		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for chain %s (selector: %d): %w", chain, chainSelector, err)
		}

		e.Logger.Infof("Transferring admin roles for %d tokens on chain %s (selector: %d)", len(tokenAdminInfo), chain, chainSelector)

		for i, info := range tokenAdminInfo {
			e.Logger.Debugf("Transferring admin from current to %s for token %s on chain %s (%d/%d)",
				info.AdminAddress.Hex(), info.TokenAddress.Hex(), chain, i+1, len(tokenAdminInfo))

			_, err = chainState.TokenAdminRegistry.TransferAdminRole(opts, info.TokenAddress, info.AdminAddress)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transferAdminRole transaction for token %s on chain %s (selector: %d): %w",
					info.TokenAddress.Hex(), chain, chainSelector, err)
			}
		}
	}

	return deployerGroup.Enact()
}

// TransferAdminRoleChangeset transfers the admin role for tokens on the token admin registry to 3rd parties.
func TransferAdminRoleChangeset(env cldf.Environment, c TokenAdminRegistryChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(env, false, validateTransferAdminRole); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid TokenAdminRegistryChangesetConfig: %w", err)
	}
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := deployergroup.NewDeployerGroup(env, state, c.MCMS).WithDeploymentContext("transfer admin role for tokens on token admin registries")

	for chainSelector, tokenSymbolToPoolInfo := range c.Pools {
		chain := env.BlockChains.EVMChains()[chainSelector]
		chainState := state.Chains[chainSelector]
		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for %s", chain)
		}
		for symbol, poolInfo := range tokenSymbolToPoolInfo {
			_, tokenAddress, err := poolInfo.GetPoolAndTokenAddress(env.GetContext(), symbol, chain, chainState)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get state of %s token on chain %s: %w", symbol, chain, err)
			}
			_, err = chainState.TokenAdminRegistry.TransferAdminRole(opts, tokenAddress, poolInfo.ExternalAdmin)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transferAdminRole transaction for %s on %s registry: %w", symbol, chain, err)
			}
		}
	}

	return deployerGroup.Enact()
}
