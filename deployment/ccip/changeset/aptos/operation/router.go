package operation

import (
	"encoding/json"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router"
	aptos_router "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router/router"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	"github.com/smartcontractkit/mcms/types"
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

func updateRouter(b operations.Bundle, deps AptosDeps, in UpdateRouterDestInput) ([]types.Transaction, error) {
	// Bind CCIP Package
	ccipAddress := deps.OnChainState.CCIPAddress
	routerBind := ccip_router.Bind(ccipAddress, deps.AptosChain.Client)

	// Process each destination chain config update
	var txs []types.Transaction

	var destChainSelectors []uint64
	var onRampVersions [][]byte
	for _, update := range in.Updates {
		destChainSelectors = append(destChainSelectors, update.DestChainSelector)
		onRampVersions = append(onRampVersions, update.OnRampVersion)
	}
	moduleInfo, function, _, args, err := routerBind.Router().Encoder().SetOnRampVersions(destChainSelectors, onRampVersions)
	if err != nil {
		return []types.Transaction{}, fmt.Errorf("failed to encode ApplyDestChainConfigUpdates for chains %d: %w", uint64(14767482510784806043), err)
	}

	additionalFields := aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}
	afBytes, err := json.Marshal(additionalFields)
	if err != nil {
		return []types.Transaction{}, fmt.Errorf("failed to marshal additional fields: %w", err)
	}

	txs = append(txs, types.Transaction{
		To:               ccipAddress.StringLong(),
		Data:             aptosmcms.ArgsToData(args),
		AdditionalFields: afBytes,
	})

	b.Logger.Infow("Adding Router destination config update operation")

	return txs, nil
}
