package v1_5_1

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type TokenPoolFactory interface {
	Address() common.Address
	TypeAndVersion(*bind.CallOpts) (string, error)
}

type TokenPoolFactoryView struct {
	TypeAndVersion string         `json:"typeAndVersion,omitempty"`
	Address        common.Address `json:"address,omitempty"`
}

func GenerateTokenPoolFactoryView(factory TokenPoolFactory) (TokenPoolFactoryView, error) {
	typeAndVersion, err := factory.TypeAndVersion(nil)
	if err != nil {
		return TokenPoolFactoryView{}, err
	}

	return TokenPoolFactoryView{
		TypeAndVersion: typeAndVersion,
		Address:        factory.Address(),
	}, nil
}
