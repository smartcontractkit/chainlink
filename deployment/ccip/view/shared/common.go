package shared

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/token_admin_registry"
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
