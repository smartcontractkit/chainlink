package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type ContractMetaData struct {
	TypeAndVersion string         `json:"typeAndVersion,omitempty"`
	Address        common.Address `json:"address,omitempty"`
	Owner          common.Address `json:"owner,omitempty"`
}

func NewContractMetaData(tv Meta, addr common.Address) (ContractMetaData, error) {
	tvStr, err := tv.TypeAndVersion(nil)
	if err != nil {
		return ContractMetaData{}, fmt.Errorf("failed to get type and version addr %s: %w", addr.String(), err)
	}
	owner, err := tv.Owner(nil)
	if err != nil {
		return ContractMetaData{}, fmt.Errorf("failed to get owner addr %s: %w", addr.String(), err)
	}
	return ContractMetaData{
		TypeAndVersion: tvStr,
		Address:        addr,
		Owner:          owner,
	}, nil
}

type Meta interface {
	TypeAndVersion(opts *bind.CallOpts) (string, error)
	Owner(opts *bind.CallOpts) (common.Address, error)
}
