package v1_6

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var _ cldf.ChangeSet[AddRegistryModuleConfig] = AddRegistryModuleChangeset

type AddRegistryModuleConfig struct {
	ChainSelectors     []uint64       // which chains to add registry modules on
	RegistryModuleAddr common.Address // address of the registry module to add (if same across chains)
	// Optional: map of chain selector to registry module address if different per chain
	RegistryModuleAddrs map[uint64]common.Address
}

func (c AddRegistryModuleConfig) Validate(e cldf.Environment) error {
	if len(c.ChainSelectors) == 0 {
		return fmt.Errorf("no chain selectors provided")
	}

	// Either use single address or per-chain addresses, not both
	if (c.RegistryModuleAddr != common.Address{}) && len(c.RegistryModuleAddrs) > 0 {
		return fmt.Errorf("cannot specify both single registry module address and per-chain addresses")
	}

	// Must specify at least one way to get the address
	if (c.RegistryModuleAddr == common.Address{}) && len(c.RegistryModuleAddrs) == 0 {
		return fmt.Errorf("must specify registry module address(es)")
	}

	for _, chainSel := range c.ChainSelectors {
		if err := cldf.IsValidChainSelector(chainSel); err != nil {
			return fmt.Errorf("invalid chain selector %d: %w", chainSel, err)
		}

		if _, exists := e.BlockChains.EVMChains()[chainSel]; !exists {
			return fmt.Errorf("chain %d not found in environment", chainSel)
		}

		// If using per-chain addresses, ensure all chains have an address
		if len(c.RegistryModuleAddrs) > 0 {
			if _, exists := c.RegistryModuleAddrs[chainSel]; !exists {
				return fmt.Errorf("no registry module address specified for chain %d", chainSel)
			}
		}
	}

	return nil
}

func AddRegistryModuleChangeset(e cldf.Environment, cfg AddRegistryModuleConfig) (cldf.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid config: %w", err)
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	for _, chainSel := range cfg.ChainSelectors {
		chain := e.BlockChains.EVMChains()[chainSel]
		chainState, exists := state.Chains[chainSel]

		if !exists {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain state not found for chain %d", chainSel)
		}

		// Check if TokenAdminRegistry exists
		if chainState.TokenAdminRegistry == nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("TokenAdminRegistry not found on chain %d", chainSel)
		}

		// Get the registry module address for this chain
		var registryModuleAddr common.Address
		if len(cfg.RegistryModuleAddrs) > 0 {
			registryModuleAddr = cfg.RegistryModuleAddrs[chainSel]
		} else {
			registryModuleAddr = cfg.RegistryModuleAddr
		}

		// Check if registry module is already added
		isAlreadyModule, err := chainState.TokenAdminRegistry.IsRegistryModule(nil, registryModuleAddr)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to check if registry module is already added on chain %d: %w", chainSel, err)
		}

		if isAlreadyModule {
			e.Logger.Infow("RegistryModule already added to TokenAdminRegistry",
				"chain", chainSel,
				"registryModule", registryModuleAddr.Hex())
			continue
		}

		// Add the RegistryModule to TokenAdminRegistry
		e.Logger.Infow("Adding RegistryModule to TokenAdminRegistry",
			"chain", chainSel,
			"registryModule", registryModuleAddr.Hex())

		tx, err := chainState.TokenAdminRegistry.AddRegistryModule(chain.DeployerKey, registryModuleAddr)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to add registry module to TokenAdminRegistry on chain %d: %w", chainSel, err)
		}

		_, err = chain.Confirm(tx)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm registry module registration transaction on chain %d: %w", chainSel, err)
		}

		e.Logger.Infow("Successfully added RegistryModule to TokenAdminRegistry",
			"chain", chainSel,
			"registryModule", registryModuleAddr.Hex())
	}

	return cldf.ChangesetOutput{}, nil
}
