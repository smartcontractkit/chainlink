package v1_5

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/registry_module_owner_custom"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
)

type RegistryModulesView struct {
	types.ContractMetaData
	TokenAdminRegistry common.Address `json:"tokenAdminRegistry,omitempty"`
}

func GenerateRegistryModulesView(
	registryModulesContractv1_5 []*registry_module_owner_custom.RegistryModuleOwnerCustom,
	tokenAdminRegistry *token_admin_registry.TokenAdminRegistry,
) ([]RegistryModulesView, error) {
	registryModules := make([]RegistryModulesView, 0)
	for _, registryModuleContract := range registryModulesContractv1_5 {
		tv, err := registryModuleContract.TypeAndVersion(nil)
		if err != nil {
			return []RegistryModulesView{}, err
		}
		registryModules = append(registryModules, RegistryModulesView{
			ContractMetaData: types.ContractMetaData{
				Address:        registryModuleContract.Address(),
				TypeAndVersion: tv,
			},
			TokenAdminRegistry: tokenAdminRegistry.Address(),
		})
	}

	return registryModules, nil
}
