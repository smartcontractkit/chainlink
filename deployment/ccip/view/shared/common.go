package shared

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mr-tron/base58"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	registryModuleOwnerCustomv15 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/registry_module_owner_custom"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"
	registryModuleOwnerCustomv16 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/registry_module_owner_custom"
)

const (
	GetTokensPaginationSize = 20
)

type RegistryModulesView struct {
	TypeAndVersion     string `json:"typeAndVersion,omitempty"`
	TokenAdminRegistry string `json:"tokenAdminRegistry,omitempty"`
}

func GetSupportedTokens(taContract *token_admin_registry.TokenAdminRegistry) ([]common.Address, error) {
	startIndex := uint64(0)
	allTokens := make([]common.Address, 0)
	for {
		fetchedTokens, err := taContract.GetAllConfiguredTokens(nil, startIndex, GetTokensPaginationSize)
		if err != nil {
			return nil, err
		}
		allTokens = append(allTokens, fetchedTokens...)
		startIndex += GetTokensPaginationSize
		if len(fetchedTokens) < GetTokensPaginationSize {
			break
		}
	}
	return allTokens, nil
}

func GetAddressFromBytes(chainSelector uint64, address []byte) string {
	family, err := chain_selectors.GetSelectorFamily(chainSelector)
	if err != nil {
		return "invalid chain selector"
	}

	switch family {
	case chain_selectors.FamilyEVM:
		return strings.ToLower(common.BytesToAddress(address).Hex())
	case chain_selectors.FamilySolana:
		return base58.Encode(address)
	default:
		return "unsupported chain family"
	}
}

func GetRegistryModuleView(registryModule any, tokenAdminRegistry common.Address) (RegistryModulesView, error) {
	switch module := registryModule.(type) {
	case *registryModuleOwnerCustomv15.RegistryModuleOwnerCustom:
		tv, err := module.TypeAndVersion(nil)
		if err != nil {
			return RegistryModulesView{}, err
		}
		return RegistryModulesView{
			TypeAndVersion:     tv,
			TokenAdminRegistry: tokenAdminRegistry.Hex(),
		}, nil

	case *registryModuleOwnerCustomv16.RegistryModuleOwnerCustom:
		tv, err := module.TypeAndVersion(nil)
		if err != nil {
			return RegistryModulesView{}, err
		}
		return RegistryModulesView{
			TypeAndVersion:     tv,
			TokenAdminRegistry: tokenAdminRegistry.Hex(),
		}, nil
	default:
		return RegistryModulesView{}, fmt.Errorf("unsupported registry module type: %T", module)
	}
}
