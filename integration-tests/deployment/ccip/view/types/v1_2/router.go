package v1_2

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

type RouterView struct {
	types.ContractMetaData
	WrappedNative common.Address            `json:"wrappedNative,omitempty"`
	ARMProxy      common.Address            `json:"armProxy,omitempty"`
	OnRamps       map[uint64]common.Address `json:"onRamps,omitempty"`  // Map of DestinationChainSelectors to OnRamp Addresses
	OffRamps      map[uint64]common.Address `json:"offRamps,omitempty"` // Map of SourceChainSelectors to a list of OffRamp Addresses
}

func GenerateRouterView(r *router.Router) (RouterView, error) {
	meta, err := types.NewContractMetaData(r, r.Address())
	wrappedNative, err := r.GetWrappedNative(nil)
	if err != nil {
		return RouterView{}, fmt.Errorf("view error to get router wrapped native: %w", err)
	}
	armProxy, err := r.GetArmProxy(nil)
	if err != nil {
		return RouterView{}, fmt.Errorf("view error to get router arm proxy: %w", err)
	}
	onRamps := make(map[uint64]common.Address)
	offRamps := make(map[uint64]common.Address)
	offRampList, err := r.GetOffRamps(nil)
	if err != nil {
		return RouterView{}, fmt.Errorf("view error to get router offRamps: %w", err)
	}
	for _, offRamp := range offRampList {
		offRamps[offRamp.SourceChainSelector] = offRamp.OffRamp
	}
	for selector := range offRamps {
		onRamp, err := r.GetOnRamp(nil, selector)
		if err != nil {
			return RouterView{}, fmt.Errorf("view error to get router onRamp: %w", err)
		}
		onRamps[selector] = onRamp
	}
	return RouterView{
		ContractMetaData: meta,
		WrappedNative:    wrappedNative,
		ARMProxy:         armProxy,
		OnRamps:          onRamps,
		OffRamps:         offRamps,
	}, nil
}
