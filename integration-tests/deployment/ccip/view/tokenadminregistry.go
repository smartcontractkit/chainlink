package view

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type TokenAdminRegistry struct {
	Contract
	Tokens []common.Address `json:"tokens"`
}

func (ta TokenAdminRegistry) Address() common.Address {
	return common.HexToAddress(ta.Contract.Address)
}

func TokenAdminRegistrySnapshot(taContract TokenAdminRegistryGetter) (TokenAdminRegistry, error) {
	tokens, err := taContract.GetAllConfiguredTokens(nil, 0, 10)
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

type TokenAdminRegistryGetter interface {
	GetAllConfiguredTokens(opts *bind.CallOpts, startIndex uint64, maxCount uint64) ([]common.Address, error)
	TypeAndVersion(opts *bind.CallOpts) (string, error)
	Address() common.Address
}
