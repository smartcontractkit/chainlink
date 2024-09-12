package view

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

type Router struct {
	types.Contract
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

func RouterSnapshot(r *router.Router) (Router, error) {
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
		offRamps[offRamp.SourceChainSelector] = offRamp.OffRamp.Hex()
	}
	for selector := range offRamps {
		onRamp, err := r.GetOnRamp(nil, selector)
		if err != nil {
			return Router{}, err
		}
		onRamps[selector] = onRamp.Hex()
	}
	return Router{
		Contract: types.Contract{
			Address:        r.Address().Hex(),
			TypeAndVersion: tv,
		},
		WrappedNative: wrappedNative.Hex(),
		ARMProxy:      armProxy.Hex(),
		OnRamps:       onRamps,
		OffRamps:      offRamps,
	}, nil
}
