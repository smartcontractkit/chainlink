package v1_6

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/registry_module_owner_custom"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
)

var _ cldf.ChangeSet[DeployRegistryModuleConfig] = DeployRegistryModuleChangeset

type DeployRegistryModuleConfig struct {
	ChainSelectors []uint64 // which chains to deploy the registry module on
}

func (c DeployRegistryModuleConfig) Validate(e cldf.Environment) error {
	if len(c.ChainSelectors) == 0 {
		return errors.New("no chain selectors provided")
	}

	// Load state to validate chain states and TokenAdminRegistry
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for _, chainSel := range c.ChainSelectors {
		if err := cldf.IsValidChainSelector(chainSel); err != nil {
			return fmt.Errorf("invalid chain selector %d: %w", chainSel, err)
		}

		// Check if chain state exists
		chainState, exists := state.Chains[chainSel]
		if !exists {
			return fmt.Errorf("chain state not found for chain %d", chainSel)
		}

		// Check if TokenAdminRegistry exists
		if chainState.TokenAdminRegistry == nil {
			return fmt.Errorf("TokenAdminRegistry not found on chain %d", chainSel)
		}
	}

	return nil
}

func DeployRegistryModuleChangeset(e cldf.Environment, cfg DeployRegistryModuleConfig) (cldf.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid config: %w", err)
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	addressBook := cldf.NewMemoryAddressBook()
	ds := datastore.NewMemoryDataStore()

	for _, chainSel := range cfg.ChainSelectors {
		chain := e.BlockChains.EVMChains()[chainSel]
		chainState := state.Chains[chainSel]

		// Check if we need to deploy RegistryModuleOwnerCustom 1.6.0
		needsDeploy := NeedsRegistryModule16Deployment(chainState)

		if !needsDeploy {
			e.Logger.Infow("RegistryModuleOwnerCustom 1.6.0 already deployed", "chain", chainSel)
			continue
		}

		e.Logger.Infow("Deploying RegistryModuleOwnerCustom 1.6.0", "chain", chainSel)

		registryModule, err := cldf.DeployContract(e.Logger, chain, addressBook,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*registry_module_owner_custom.RegistryModuleOwnerCustom] {
				var (
					regModAddr common.Address
					tx         *types.Transaction
					regMod     *registry_module_owner_custom.RegistryModuleOwnerCustom
					err2       error
				)

				if chain.IsZkSyncVM {
					regModAddr, _, regMod, err2 = registry_module_owner_custom.DeployRegistryModuleOwnerCustomZk(
						nil,
						chain.ClientZkSyncVM,
						chain.DeployerKeyZkSyncVM,
						chain.Client,
						chainState.TokenAdminRegistry.Address(),
					)
					// ZkSync deployment doesn't return a transaction, so tx remains nil
				} else {
					regModAddr, tx, regMod, err2 = registry_module_owner_custom.DeployRegistryModuleOwnerCustom(
						chain.DeployerKey,
						chain.Client,
						chainState.TokenAdminRegistry.Address(),
					)
				}

				return cldf.ContractDeploy[*registry_module_owner_custom.RegistryModuleOwnerCustom]{
					Address:  regModAddr,
					Contract: regMod,
					Tx:       tx,
					Tv:       cldf.NewTypeAndVersion(shared.RegistryModule, deployment.Version1_6_0),
					Err:      err2,
				}
			})

		if err != nil {
			return cldf.ChangesetOutput{DataStore: ds}, fmt.Errorf("failed to deploy registry module on chain %d: %w", chainSel, err)
		}

		// Add the address reference to the datastore
		if err = ds.Addresses().Add(datastore.AddressRef{
			ChainSelector: chainSel,
			Address:       registryModule.Address.Hex(),
			Type:          datastore.ContractType(shared.RegistryModule),
			Version:       semver.MustParse("1.6.0"),
			Labels: datastore.NewLabelSet(
				"RegistryModuleOwnerCustom 1.6.0",
			),
		}); err != nil {
			return cldf.ChangesetOutput{DataStore: ds},
				fmt.Errorf("failed to save address ref for chain %d: %w", chainSel, err)
		}

		e.Logger.Infow("Successfully deployed RegistryModuleOwnerCustom 1.6.0",
			"chain", chainSel,
			"address", registryModule.Address.Hex())
	}

	return cldf.ChangesetOutput{
		AddressBook: addressBook,
		DataStore:   ds,
	}, nil
}

// needsRegistryModule16Deployment checks if we need to deploy RegistryModuleOwnerCustom 1.6.0
// Returns true if:
// - No 1.6.0 registry modules exist on the chain
// - Only non-1.6.0 versions exist
func NeedsRegistryModule16Deployment(chainState evm.CCIPChainState) bool {
	// Check if any 1.6.0 registry modules exist
	return len(chainState.RegistryModules1_6) == 0
}
