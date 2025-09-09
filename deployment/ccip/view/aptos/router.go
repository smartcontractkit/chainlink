package aptos

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	aptosCommon "github.com/smartcontractkit/chainlink/deployment/common/view/aptos"
)

type RouterView struct {
	aptosCommon.ContractMetaData

	IsTestRouter bool              `json:"isTestRouter"`
	OnRamps      map[uint64]string `json:"onRamps"`  // Map of DestinationChainSelector to OnRampAddress
	OffRamps     map[uint64]string `json:"offRamps"` // Map of DestinationChainSelector to OffRampAddress
}

func GenerateRouterView(chain cldf_aptos.Chain, routerAddress aptos.AccountAddress, offRampAddresses []aptos.AccountAddress, isTestRouter bool) (RouterView, error) {
	boundRouter := ccip_router.Bind(routerAddress, chain.Client)

	typeAndVersion, err := boundRouter.Router().TypeAndVersion(nil)
	if err != nil {
		return RouterView{}, fmt.Errorf("failed to get typeAndVersion of router %s: %w", routerAddress.StringLong(), err)
	}
	owner, err := boundRouter.Router().Owner(nil)
	if err != nil {
		return RouterView{}, fmt.Errorf("failed to get owner of router %s: %w", routerAddress.StringLong(), err)
	}

	// OnRamps
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

	// OffRamps
	// Since on Aptos, we're not tracking the offramps in the router, we're instead iterating over all known offramps
	// and for each are checking if it has a source chain config set for the current router.
	offRamps := make(map[uint64]string)
	for _, offRampAddress := range offRampAddresses {
		boundOffRamp := ccip_offramp.Bind(offRampAddress, chain.Client)
		sourceChainSelectors, sourceChainConfigs, err := boundOffRamp.Offramp().GetAllSourceChainConfigs(nil)
		if err != nil {
			return RouterView{}, fmt.Errorf("failed to get sourceChainConfigs of offRamp %s: %w", offRampAddress.StringLong(), err)
		}
		for i, sourceChainConfig := range sourceChainConfigs {
			if sourceChainConfig.Router == routerAddress {
				sourceChainSelector := sourceChainSelectors[i]
				offRamps[sourceChainSelector] = offRampAddress.StringLong()
			}
		}
	}

	return RouterView{
		ContractMetaData: aptosCommon.ContractMetaData{
			Address:        routerAddress.StringLong(),
			Owner:          owner.StringLong(),
			TypeAndVersion: typeAndVersion,
		},
		IsTestRouter: isTestRouter,
		OnRamps:      onRamps,
		OffRamps:     offRamps,
	}, nil
}
