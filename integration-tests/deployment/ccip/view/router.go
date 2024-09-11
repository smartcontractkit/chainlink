package view

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type Router struct {
	Contract
	WrappedNative string            `json:"wrappedNative,omitempty"`
	ARMProxy      string            `json:"armProxy,omitempty"`
	OnRamps       map[uint64]string `json:"onRamps,omitempty"`  // Map of DestinationChainSelectors to OnRamp Addresses
	OffRamps      map[uint64]string `json:"offRamps,omitempty"` // Map of SourceChainSelectors to a list of OffRamp Addresses
}

func (r Router) Address() common.Address {
	return common.HexToAddress(r.Contract.Address)
}

func (r Router) DestinationChainSelectors() []uint64 {
	selectors := make([]uint64, 0, len(r.OffRamps))
	for selector := range r.OffRamps {
		selectors = append(selectors, selector)
	}
	return selectors
}

func (r Router) SourceChainSelectors() []uint64 {
	selectors := make([]uint64, 0, len(r.OnRamps))
	for selector := range r.OnRamps {
		selectors = append(selectors, selector)
	}
	return selectors
}

func RouterSnapshot(r RouterReader) (Router, error) {
	tv, err := r.TypeAndVersion(nil)
	if err != nil {
		return Router{}, err
	}
	wrappedNative, err := r.GetWrappedNative(nil)
	if err != nil {
		return Router{}, err
	}
	armProxy, err := r.GetArmProxy(nil)
	if err != nil {
		return Router{}, err
	}
	onRamps := make(map[uint64]string)
	offRamps := make(map[uint64]string)
	offRampList, err := r.GetOffRamps(nil)
	if err != nil {
		return Router{}, err
	}
	for _, offRamp := range offRampList {
		offRamps[offRamp.SourceChainSelector] = offRamp.OffRamp
	}
	for selector := range offRamps {
		onRamp, err := r.GetOnRamp(nil, selector)
		if err != nil {
			return Router{}, err
		}
		onRamps[selector] = onRamp.Hex()
	}
	return Router{
		Contract: Contract{
			Address:        r.Address().Hex(),
			TypeAndVersion: tv,
		},
		WrappedNative: wrappedNative.Hex(),
		ARMProxy:      armProxy.Hex(),
		OnRamps:       onRamps,
		OffRamps:      offRamps,
	}, nil
}

type RouterOffRamp struct {
	SourceChainSelector uint64 `json:"sourceChainSelector"`
	OffRamp             string `json:"offRamp"`
}

type RouterReader interface {
	ContractState
	GetOffRamps(opts *bind.CallOpts) ([]RouterOffRamp, error)
	GetOnRamp(opts *bind.CallOpts, destChainSelector uint64) (common.Address, error)
	GetWrappedNative(opts *bind.CallOpts) (common.Address, error)
	GetArmProxy(opts *bind.CallOpts) (common.Address, error)
}
