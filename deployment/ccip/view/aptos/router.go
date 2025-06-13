package aptos

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	aptosCommon "github.com/smartcontractkit/chainlink/deployment/common/view/aptos"
)

type RouterView struct {
	aptosCommon.ContractMetaData

	OnRamps map[uint64]string `json:"onRamps"` // Map of DestinationChainSelector to OnRampAddress
}

func GenerateRouterView(chain cldf_aptos.Chain, routerAddress aptos.AccountAddress) (RouterView, error) {
	boundRouter := ccip_router.Bind(routerAddress, chain.Client)

	typeAndVersion, err := boundRouter.Router().TypeAndVersion(nil)
	if err != nil {
		return RouterView{}, fmt.Errorf("failed to get typeAndVersion of router %s: %w", routerAddress.StringLong(), err)
	}
	owner, err := boundRouter.Router().Owner(nil)
	if err != nil {
		return RouterView{}, fmt.Errorf("failed to get owner of router %s: %w", routerAddress.StringLong(), err)
	}

	destinationChainSelectors, err := boundRouter.Router().GetDestChains(nil)
	if err != nil {
		return RouterView{}, fmt.Errorf("failed to get destChainSelectors of router %s: %w", routerAddress.StringLong(), err)
	}
	onrampVersions, err := boundRouter.Router().GetOnRampVersions(nil, destinationChainSelectors)
	if err != nil {
		return RouterView{}, fmt.Errorf("failed to get onRamp versions of router %s: %w", routerAddress.StringLong(), err)
	}

	onRamps := make(map[uint64]string, len(onrampVersions))
	for i, destChainSelector := range destinationChainSelectors {
		onRampAddress, err := boundRouter.Router().GetOnRampForVersion(nil, onrampVersions[i])
		if err != nil {
			return RouterView{}, fmt.Errorf("failed to get onRamp for version %d of router %s: %w", onrampVersions[i], routerAddress.StringLong(), err)
		}
		onRamps[destChainSelector] = onRampAddress.StringLong()
	}

	return RouterView{
		ContractMetaData: aptosCommon.ContractMetaData{
			Address:        routerAddress.StringLong(),
			Owner:          owner.StringLong(),
			TypeAndVersion: typeAndVersion,
		},
		OnRamps: onRamps,
	}, nil
}
