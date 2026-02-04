package v1_6_1

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-contract-examples/chains/evm/gobindings/generated/1_6_1/proxy_admin"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	_ cldf.ChangeSet[ProxyAdminChangesetConfig] = TransferProxyAdminOwnership
)

type ProxyAdmin struct {
	Symbol   string
	NewOwner common.Address
}

type ProxyAdminChangesetConfig struct {
	Tokens map[uint64]map[string]ProxyAdmin
	MCMS   *proposalutils.TimelockConfig
}

func (c ProxyAdminChangesetConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, tokens := range c.Tokens {
		for token, config := range tokens {
			chainState, ok := state.EVMChainState(chainSelector)
			if !ok {
				return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
			}

			chain, ok := e.BlockChains.EVMChains()[chainSelector]
			if !ok {
				return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
			}

			if _, ok := chainState.TransparentUpgradeableProxy[shared.TokenSymbol(config.Symbol)]; !ok {
				return fmt.Errorf("transparent upgradeable proxy does not exist for %s token on %s", token, chain.Name())
			}

			if _, ok := chainState.ProxyAdmin[shared.TokenSymbol(config.Symbol)]; !ok {
				return fmt.Errorf("proxy admin does not exist for %s token on %s", token, chain.Name())
			}
		}
	}
	return nil
}

// TransferProxyAdminOwnership transfers ownership of ProxyAdmin contracts to new owners as specified in the config.
func TransferProxyAdminOwnership(e cldf.Environment, c ProxyAdminChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid ProxyAdminChangesetConfig: %w", err)
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := deployergroup.NewDeployerGroup(e, state, c.MCMS).WithDeploymentContext("transfer ownership of ProxyAdmin")

	for chainSelector, tokens := range c.Tokens {
		chain := e.BlockChains.EVMChains()[chainSelector]

		chainState, ok := state.EVMChainState(chainSelector)
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
		}

		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for %s", chain)
		}

		for token, config := range tokens {
			if config.NewOwner == (common.Address{}) {
				return cldf.ChangesetOutput{}, fmt.Errorf("NewOwner address is required for %s token on chain %s", token, chain.Name())
			}

			transparent, ok := chainState.TransparentUpgradeableProxy[shared.TokenSymbol(config.Symbol)]

			if !ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("TransparentUpgradeableProxy does not exist for %s token on %s", token, chain.Name())
			}

			storageBytes, err := chain.Client.StorageAt(e.GetContext(), transparent.Address(), shared.AdminSlot, nil)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get storage at slot %s for TransparentUpgradeableProxy for %s token on %s: %w", shared.AdminSlot, token, chain, err)
			}

			proxyAdmin := common.BytesToAddress(storageBytes)
			proxy, err := proxy_admin.NewProxyAdmin(proxyAdmin, chain.Client)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to connect to ProxyAdmin at %s for %s token on %s: %w", proxyAdmin, token, chain, err)
			}

			if _, err := proxy.TransferOwnership(opts, config.NewOwner); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transferOwnership transaction to %s for ProxyAdmin at %s for %s token on %s: %w", config.NewOwner, proxyAdmin, token, chain.Name(), err)
			}
		}
	}

	return deployerGroup.Enact()
}
