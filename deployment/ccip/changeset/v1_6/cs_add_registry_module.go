package v1_6

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
)

var _ cldf.ChangeSet[AddRegistryModuleConfig] = AddRegistryModuleChangeset

type AddRegistryModuleConfig struct {
	ChainSelectors      []uint64         // which chains to add registry modules on
	RegistryModuleAddrs []common.Address // addresses of the registry modules to add (must match length of ChainSelectors)
}

func (c AddRegistryModuleConfig) Validate(e cldf.Environment) error {
	if len(c.ChainSelectors) == 0 {
		return fmt.Errorf("no chain selectors provided")
	}

	if len(c.RegistryModuleAddrs) == 0 {
		return fmt.Errorf("no registry module addresses provided")
	}

	if len(c.ChainSelectors) != len(c.RegistryModuleAddrs) {
		return fmt.Errorf("chain selectors and registry module addresses must have the same length")
	}

	for _, chainSel := range c.ChainSelectors {
		if err := cldf.IsValidChainSelector(chainSel); err != nil {
			return fmt.Errorf("invalid chain selector %d: %w", chainSel, err)
		}

		if _, exists := e.BlockChains.EVMChains()[chainSel]; !exists {
			return fmt.Errorf("chain %d not found in environment", chainSel)
		}
	}

	for i, addr := range c.RegistryModuleAddrs {
		if addr == (common.Address{}) {
			return fmt.Errorf("registry module address at index %d is zero address", i)
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

	for i, chainSel := range cfg.ChainSelectors {
		chain := e.BlockChains.EVMChains()[chainSel]
		chainState, exists := state.Chains[chainSel]

		if !exists {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain state not found for chain %d", chainSel)
		}

		// Check if TokenAdminRegistry exists
		if chainState.TokenAdminRegistry == nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("TokenAdminRegistry not found on chain %d", chainSel)
		}

		registryModuleAddr := cfg.RegistryModuleAddrs[i]

		// Check if registry module is already added
		isAlreadyModule, err := chainState.TokenAdminRegistry.IsRegistryModule(nil, registryModuleAddr)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to check if registry module is already added on chain %d: %w", chainSel, err)
		}

		if isAlreadyModule {
			// Check if it's a 1.6 registry module
			is16Module := isRegistryModule16(chainState, registryModuleAddr)
			if is16Module {
				e.Logger.Infow("RegistryModule 1.6 already added to TokenAdminRegistry",
					"chain", chainSel,
					"registryModule", registryModuleAddr.Hex())
				continue
			} else {
				// It's not a 1.6 module, remove the old one and add the new one
				e.Logger.Infow("Found non-1.6 RegistryModule, updating to 1.6",
					"chain", chainSel,
					"oldRegistryModule", registryModuleAddr.Hex())

				// Remove the old registry module
				removeTx, err := chainState.TokenAdminRegistry.RemoveRegistryModule(chain.DeployerKey, registryModuleAddr)
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to remove old registry module from TokenAdminRegistry on chain %d: %w", chainSel, err)
				}

				_, err = chain.Confirm(removeTx)
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm registry module removal transaction on chain %d: %w", chainSel, err)
				}

				e.Logger.Infow("Removed old RegistryModule from TokenAdminRegistry",
					"chain", chainSel,
					"registryModule", registryModuleAddr.Hex())
			}
		}

		// Add the RegistryModule to TokenAdminRegistry. Case: no registry module existed OR removed the old one
		e.Logger.Infow("Adding RegistryModule 1.6 to TokenAdminRegistry",
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

		e.Logger.Infow("Successfully added RegistryModule 1.6 to TokenAdminRegistry",
			"chain", chainSel,
			"registryModule", registryModuleAddr.Hex())
	}

	return cldf.ChangesetOutput{}, nil
}

// isRegistryModule16 checks if the given registry module address is a 1.6 version
func isRegistryModule16(chainState evm.CCIPChainState, registryModuleAddr common.Address) bool {
	// Check if the address exists in the 1.6 registry modules map
	for _, module16 := range chainState.RegistryModules1_6 {
		if module16.Address() == registryModuleAddr {
			return true
		}
	}
	return false
}
