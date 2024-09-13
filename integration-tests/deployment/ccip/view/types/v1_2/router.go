package v1_2

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types"
	router1_2 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

type Router struct {
	types.ContractMetaData
	WrappedNative common.Address            `json:"wrappedNative,omitempty"`
	ARMProxy      common.Address            `json:"armProxy,omitempty"`
	OnRamps       map[uint64]common.Address `json:"onRamps,omitempty"`  // Map of DestinationChainSelectors to OnRamp Addresses
	OffRamps      map[uint64]common.Address `json:"offRamps,omitempty"` // Map of SourceChainSelectors to a list of OffRamp Addresses
}

func (r *Router) Snapshot(routerContractMeta types.ContractMetaData, _ []types.ContractMetaData, client bind.ContractBackend) error {
	r.ContractMetaData = routerContractMeta
	if err := r.ContractMetaData.Validate(); err != nil {
		return fmt.Errorf("snapshot error for Router: %w", err)
	}
	rContract, err := router1_2.NewRouter(routerContractMeta.Address, client)
	if err != nil {
		return fmt.Errorf("snapshot error for Router: failed to get router contract: %w", err)
	}
	wrappedNative, err := rContract.GetWrappedNative(nil)
	if err != nil {
		return err
	}
	armProxy, err := rContract.GetArmProxy(nil)
	if err != nil {
		return err
	}
	onRamps := make(map[uint64]common.Address)
	offRamps := make(map[uint64]common.Address)
	offRampList, err := rContract.GetOffRamps(nil)
	if err != nil {
		return err
	}
	for _, offRamp := range offRampList {
		offRamps[offRamp.SourceChainSelector] = offRamp.OffRamp
	}
	for selector := range offRamps {
		onRamp, err := rContract.GetOnRamp(nil, selector)
		if err != nil {
			return err
		}
		onRamps[selector] = onRamp
	}
	r.WrappedNative = wrappedNative
	r.ARMProxy = armProxy
	r.OnRamps = onRamps
	r.OffRamps = offRamps
	return nil
}
