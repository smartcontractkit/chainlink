package v1_6

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/registry_module_owner_custom"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
)

var _ cldf.ChangeSet[DeployRegistryModuleConfig] = DeployRegistryModuleChangeset

type DeployRegistryModuleConfig struct {
	ChainSelectors []uint64 //which chains to deploy the registry module on
}

func (c DeployRegistryModuleConfig) Validate(e cldf.Environment) error {
	if len(c.ChainSelectors) == 0 {
		return fmt.Errorf("no chain selectors provided")
	}

	for _, chainSel := range c.ChainSelectors {
		if err := cldf.IsValidChainSelector(chainSel); err != nil {
			return fmt.Errorf("invalid chain selector %d: %w", chainSel, err)
		}

		if _, exists := e.BlockChains.EVMChains()[chainSel]; !exists {
			return fmt.Errorf("chain %d not found in environment", chainSel)
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

		// Check if we need to deploy RegistryModuleOwnerCustom 1.6.0
		needsDeploy, err := needsRegistryModule16Deployment(chainState)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to check registry module status on chain %d: %w", chainSel, err)
		}

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
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy registry module on chain %d: %w", chainSel, err)
		}

		e.Logger.Infow("Successfully deployed RegistryModuleOwnerCustom 1.6.0",
			"chain", chainSel,
			"address", registryModule.Address.Hex())
	}

	return cldf.ChangesetOutput{
		AddressBook: addressBook,
	}, nil
}

// needsRegistryModule16Deployment checks if we need to deploy RegistryModuleOwnerCustom 1.6.0
// Returns true if:
// - No 1.6.0 registry modules exist on the chain
// - Only non-1.6.0 versions exist
func needsRegistryModule16Deployment(chainState evm.CCIPChainState) (bool, error) {
	// Check if any 1.6.0 registry modules exist
	if len(chainState.RegistryModules1_6) > 0 {
		// 1.6.0 version already exists
		return false, nil
	}

	// No 1.6.0 version exists, we need to deploy it
	return true, nil
}
