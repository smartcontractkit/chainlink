package v1_2

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
)

type RouterView struct {
	types.ContractMetaData
	IsTestRouter  bool                      `json:"isTestRouter"`
	WrappedNative common.Address            `json:"wrappedNative,omitempty"`
	ARMProxy      common.Address            `json:"armProxy,omitempty"`
	OnRamps       map[uint64]common.Address `json:"onRamps,omitempty"`  // Map of DestinationChainSelectors to OnRamp Addresses
	OffRamps      map[uint64]common.Address `json:"offRamps,omitempty"` // Map of SourceChainSelectors to a list of OffRamp Addresses
}

func GenerateRouterView(r *router.Router, isTestRouter bool) (RouterView, error) {
	meta, err := types.NewContractMetaData(r, r.Address())
	if err != nil {
		return RouterView{}, fmt.Errorf("view error to get router metadata: %w", err)
	}
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
		IsTestRouter:     isTestRouter,
		WrappedNative:    wrappedNative,
		ARMProxy:         armProxy,
		OnRamps:          onRamps,
		OffRamps:         offRamps,
	}, nil
}

// From the perspective of the OnRamp, the destination chains are the source chains for the OffRamp.
func GetRemoteChainSelectors(routerContract *router.Router) ([]uint64, error) {
	remoteSelectors := make([]uint64, 0)
	offRamps, err := routerContract.GetOffRamps(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get offRamps from router: %w", err)
	}
	// lanes are bidirectional, so we get the list of source chains to know which chains are supported as destinations as well
	for _, offRamp := range offRamps {
		remoteSelectors = append(remoteSelectors, offRamp.SourceChainSelector)
	}

	return remoteSelectors, nil
}
