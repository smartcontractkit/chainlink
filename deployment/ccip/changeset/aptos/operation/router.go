package operation

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router"
	aptos_router "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router/router"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
)

// UpdateRouterDestInput contains configuration for updating FeeQuoter destination configs
type UpdateRouterDestInput struct {
	MCMSAddress aptos.AccountAddress
	Updates     []aptos_router.OnRampSet
}

// UpdateRouterOp...
var UpdateRouterOp = operations.NewOperation(
	"update-router-op",
	Version1_0_0,
	"Updates Router destination chain configurations",
	updateRouter,
)

func updateRouter(b operations.Bundle, deps AptosDeps, in UpdateRouterDestInput) (mcmstypes.Transaction, error) {
	// Bind CCIP Package
	ccipAddress := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].CCIPAddress
	routerBind := ccip_router.Bind(ccipAddress, deps.AptosChain.Client)

	// Process each destination chain config update
	var destChainSelectors []uint64
	var onRampVersions [][]byte
	for _, update := range in.Updates {
		destChainSelectors = append(destChainSelectors, update.DestChainSelector)
		onRampVersions = append(onRampVersions, update.OnRampVersion)
	}
	moduleInfo, function, _, args, err := routerBind.Router().Encoder().SetOnRampVersions(destChainSelectors, onRampVersions)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to encode ApplyDestChainConfigUpdates for chains %d: %w", deps.AptosChain.Selector, err)
	}

	tx, err := utils.GenerateMCMSTx(ccipAddress, moduleInfo, function, args)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to create transaction: %w", err)
	}
	b.Logger.Infow("Adding Router destination config update operation")

	return tx, nil
}
