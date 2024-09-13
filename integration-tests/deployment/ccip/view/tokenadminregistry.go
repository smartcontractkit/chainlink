package view

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

const (
	GetTokensPaginationSize = 20
)

type TokenAdminRegistry struct {
	Contract
	Tokens []common.Address `json:"tokens"`
}

func (ta TokenAdminRegistry) Address() common.Address {
	return common.HexToAddress(ta.Contract.Address)
}

func TokenAdminRegistrySnapshot(taContract TokenAdminRegistryReader) (TokenAdminRegistry, error) {
	tokens, err := getAllConfiguredTokensPaginated(taContract)
	if err != nil {
		return TokenAdminRegistry{}, err
	}
	tv, err := taContract.TypeAndVersion(nil)
	if err != nil {
		return TokenAdminRegistry{}, err
	}
	return TokenAdminRegistry{
		Contract: Contract{
			TypeAndVersion: tv,
			Address:        taContract.Address().Hex(),
		},
		Tokens: tokens,
	}, nil
}

type TokenAdminRegistryReader interface {
	GetAllConfiguredTokens(opts *bind.CallOpts, startIndex uint64, maxCount uint64) ([]common.Address, error)
	ContractState
}

// getAllConfiguredTokensPaginated fetches all configured tokens from the TokenAdminRegistry contract in paginated
// manner to avoid RPC timeouts since the list of configured tokens can grow to be very large over time.
func getAllConfiguredTokensPaginated(taContract TokenAdminRegistryReader) ([]common.Address, error) {
	startIndex := uint64(0)
	allTokens := make([]common.Address, 0)
	fetchedTokens := make([]common.Address, 0)
	for len(fetchedTokens) < GetTokensPaginationSize {
		fetchedTokens, err := taContract.GetAllConfiguredTokens(nil, startIndex, GetTokensPaginationSize)
		if err != nil {
			return nil, err
		}
		allTokens = append(allTokens, fetchedTokens...)
		startIndex += GetTokensPaginationSize
	}
	return allTokens, nil
}
