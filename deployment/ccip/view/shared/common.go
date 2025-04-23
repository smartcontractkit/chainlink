package shared

import (
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
		// cropped left in case of long bytes sequence
		evmAddress := common.BytesToAddress(address)
		// happy-path: evm address is non-zero
		// happy-path: if raw address is no longer than 20 bytes, there is no bytes left to check
		if evmAddress != (common.Address{}) || len(address) <= 20 {
			return evmAddress.Hex()
		}
		// if raw address longer than 20 bytes and its left-cropped version is 0-address, we should right-crop it
		return common.BytesToAddress(address[:20]).Hex()
	case chain_selectors.FamilySolana:
		return base58.Encode(address)
	default:
		return "unsupported chain family"
	}
}
