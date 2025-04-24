package shared

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mr-tron/base58"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"
)

const (
	GetTokensPaginationSize = 20
)

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
