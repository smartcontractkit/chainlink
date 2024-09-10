package view

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type NonceManager struct {
	Contract
	AuthorizedCallers []common.Address `json:"authorizedCallers,omitempty"`
}

func (nm NonceManager) Address() common.Address {
	return common.HexToAddress(nm.Contract.Address)
}

func NonceManagerSnapshot(nm NonceManagerReader) (NonceManager, error) {
	authorizedCallers, err := nm.GetAllAuthorizedCallers(nil)
	if err != nil {
		return NonceManager{}, err
	}
	tv, err := nm.TypeAndVersion(nil)
	if err != nil {
		return NonceManager{}, err
	}
	return NonceManager{
		Contract: Contract{
			TypeAndVersion: tv,
			Address:        nm.Address().Hex(),
		},
		// TODO: these can be resolved using an address book
		AuthorizedCallers: authorizedCallers,
	}, nil
}

type NonceManagerReader interface {
	GetAllAuthorizedCallers(opts *bind.CallOpts) ([]common.Address, error)
	ContractState
}
