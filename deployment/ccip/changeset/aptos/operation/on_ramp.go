package operation

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_onramp"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
)

// UpdateOnRampDestsInput contains configuration for updating OnRamp destinations
type UpdateOnRampDestsInput struct {
	MCMSAddress aptos.AccountAddress
	Updates     map[uint64]v1_6.OnRampDestinationUpdate
}

// UpdateOnRampDestsOp operation to update OnRamp destination configurations
var UpdateOnRampDestsOp = operations.NewOperation(
	"update-onramp-dests-op",
	Version1_0_0,
	"Updates OnRamp destination chain configurations",
	updateOnRampDests,
)

func updateOnRampDests(b operations.Bundle, deps AptosDeps, in UpdateOnRampDestsInput) ([]mcmstypes.Transaction, error) {
	var txs []mcmstypes.Transaction

	aptosState := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector]
	// Bind CCIP Package
	ccipAddress := aptosState.CCIPAddress
	onrampBind := ccip_onramp.Bind(ccipAddress, deps.AptosChain.Client)

	// Transform the updates into the format expected by the Aptos contract
	var destChainSelectors []uint64
	var destChainRouters []aptos.AccountAddress
	var destChainAllowlistEnabled []bool

	// Get routers state addresses
	var testRouterStateAddress aptos.AccountAddress
	var routerStateAddress aptos.AccountAddress
	if aptosState.TestRouterAddress != (aptos.AccountAddress{}) {
		testRouter := ccip_router.Bind(aptosState.TestRouterAddress, deps.AptosChain.Client)
		stateAddress, err := testRouter.Router().GetStateAddress(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get test router state address: %w", err)
		}
		testRouterStateAddress = stateAddress
	}
	router := ccip_router.Bind(ccipAddress, deps.AptosChain.Client)
	routerStateAddress, err := router.Router().GetStateAddress(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get router state address: %w", err)
	}

	// Process each destination chain config update
	for destChainSelector, update := range in.Updates {
		// destChainRouters
		if !update.IsEnabled {
			destChainRouters = append(destChainRouters, aptos.AccountAddress{})
			continue
		}
		if update.TestRouter {
			destChainRouters = append(destChainRouters, testRouterStateAddress)
		} else {
			destChainRouters = append(destChainRouters, routerStateAddress)
		}
		// destChainSelectors
		destChainSelectors = append(destChainSelectors, destChainSelector)
		// destChainAllowlistEnabled
		destChainAllowlistEnabled = append(destChainAllowlistEnabled, update.AllowListEnabled)
	}

	if len(destChainSelectors) == 0 {
		b.Logger.Infow("No OnRamp destination updates to apply")
		return nil, nil
	}

	// Encode the update operation
	moduleInfo, function, _, args, err := onrampBind.Onramp().Encoder().ApplyDestChainConfigUpdates(
		destChainSelectors,
		destChainRouters,
		destChainAllowlistEnabled,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode ApplyDestChainConfigUpdates for OnRamp: %w", err)
	}
	tx, err := utils.GenerateMCMSTx(ccipAddress, moduleInfo, function, args)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	txs = append(txs, tx)

	b.Logger.Infow("Adding OnRamp destination config update operation",
		"chainCount", len(destChainSelectors))

	return txs, nil
}
