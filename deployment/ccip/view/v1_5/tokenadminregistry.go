package v1_5

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/token_admin_registry"
)

const (
	GetTokensPaginationSize = 20
)

type TokenAdminRegistryView struct {
	types.ContractMetaData
	Tokens []common.Address `json:"tokens"`
}

func GenerateTokenAdminRegistryView(taContract *token_admin_registry.TokenAdminRegistry) (TokenAdminRegistryView, error) {
	if taContract == nil {
		return TokenAdminRegistryView{}, errors.New("token admin registry contract is nil")
	}
	tokens, err := getAllConfiguredTokensPaginated(taContract)
	if err != nil {
		return TokenAdminRegistryView{}, fmt.Errorf("view error for token admin registry: %w", err)
	}
	tvMeta, err := types.NewContractMetaData(taContract, taContract.Address())
	if err != nil {
		return TokenAdminRegistryView{}, fmt.Errorf("metadata error for token admin registry: %w", err)
	}
	return TokenAdminRegistryView{
		ContractMetaData: tvMeta,
		Tokens:           tokens,
	}, nil
}

// getAllConfiguredTokensPaginated fetches all configured tokens from the TokenAdminRegistry contract in paginated
// manner to avoid RPC timeouts since the list of configured tokens can grow to be very large over time.
func getAllConfiguredTokensPaginated(taContract *token_admin_registry.TokenAdminRegistry) ([]common.Address, error) {
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
