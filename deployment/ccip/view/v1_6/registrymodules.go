package v1_6

import (
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/registry_module_owner_custom"
)

type RegistryModulesView struct {
	TypeAndVersion     string `json:"typeAndVersion,omitempty"`
	TokenAdminRegistry string `json:"tokenAdminRegistry,omitempty"`
}

func GenerateRegistryModulesView(
	registryModulesContractv1_6 []*registry_module_owner_custom.RegistryModuleOwnerCustom,
	tokenAdminRegistry *token_admin_registry.TokenAdminRegistry,
) (map[string]RegistryModulesView, error) {
	registryModules := make(map[string]RegistryModulesView, 0)
	for _, registryModuleContract := range registryModulesContractv1_6 {
		tv, err := registryModuleContract.TypeAndVersion(nil)
		if err != nil {
			return map[string]RegistryModulesView{}, err
		}
		registryModules[registryModuleContract.Address().Hex()] = RegistryModulesView{
			TypeAndVersion:     tv,
			TokenAdminRegistry: tokenAdminRegistry.Address().Hex(),
		}
	}

	return registryModules, nil
}
