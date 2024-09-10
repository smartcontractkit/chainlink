package view

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type Router struct {
	Contract
	WrappedNative string              `json:"wrappedNative"`
	ARMProxy      string              `json:"armProxy"`
	OnRamps       map[uint64]string   `json:"onRamps"`  // Map of DestinationChainSelectors to OnRamp Addresses
	OffRamps      map[uint64][]string `json:"offRamps"` // Map of SourceChainSelectors to a list of OffRamp Addresses
}

func (r Router) Address() common.Address {
	return common.HexToAddress(r.Contract.Address)
}

func RouterSnapshot(r RouterReader) (Router, error) {
	wrappedNative, err := r.GetWrappedNative(nil)
	if err != nil {
		return Router{}, err
	}
	armProxy, err := r.GetArmProxy(nil)
	if err != nil {
		return Router{}, err
	}
	onRamps := make(map[uint64]string)
	offRamps := make(map[uint64][]string)
	offRampList, err := r.GetOffRamps(nil)
	if err != nil {
		return Router{}, err
	}
	for _, offRamp := range offRampList {
		offRamps[offRamp.SourceChainSelector] = append(offRamps[offRamp.SourceChainSelector], offRamp.OffRamp.Hex())
	}
	return Router{
		Contract:      Contract{Address: r.Address().Hex()},
		WrappedNative: wrappedNative.Hex(),
		ARMProxy:      armProxy.Hex(),
		OnRamps:       onRamps,
		OffRamps:      offRamps,
	}, nil
}

type RouterReader interface {
	GetOffRamps(opts *bind.CallOpts, SourceChainSelector uint64) (common.Address, error)
	GetOnRamp(opts *bind.CallOpts, destChainSelector uint64) (common.Address, error)
	GetWrappedNative(opts *bind.CallOpts) (common.Address, error)
	GetArmProxy(opts *bind.CallOpts) (common.Address, error)
}
