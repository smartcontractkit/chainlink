package operation

import (
	"encoding/json"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/operations"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	"github.com/smartcontractkit/mcms/types"
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

func updateOnRampDests(b operations.Bundle, deps AptosDeps, in UpdateOnRampDestsInput) ([]types.Operation, error) {
	// Bind CCIP Package
	ccipAddress := deps.OnChainState.CCIPAddress
	ccipBind := ccip.Bind(ccipAddress, deps.AptosChain.Client)

	// Transform the updates into the format expected by the Aptos contract
	var destChainSelectors []uint64
	var destChainEnabled []bool
	var destChainAllowlistEnabled []bool

	// Process each destination chain config update
	for destChainSelector, update := range in.Updates {
		destChainSelectors = append(destChainSelectors, destChainSelector)
		destChainEnabled = append(destChainEnabled, update.IsEnabled)
		destChainAllowlistEnabled = append(destChainAllowlistEnabled, update.AllowListEnabled)
	}

	if len(destChainSelectors) == 0 {
		b.Logger.Infow("No OnRamp destination updates to apply")
		return []types.Operation{}, nil
	}

	// Encode the update operation
	moduleInfo, function, _, args, err := ccipBind.Onramp().Encoder().ApplyDestChainConfigUpdates(
		destChainSelectors,
		destChainEnabled,
		destChainAllowlistEnabled,
	)
	if err != nil {
		return []types.Operation{}, fmt.Errorf("failed to encode ApplyDestChainConfigUpdates for OnRamp: %w", err)
	}

	// Create MCMS operation
	additionalFields := aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}
	afBytes, err := json.Marshal(additionalFields)
	if err != nil {
		return []types.Operation{}, fmt.Errorf("failed to marshal additional fields: %w", err)
	}

	operation := types.Operation{
		ChainSelector: types.ChainSelector(deps.AptosChain.Selector),
		Transaction: types.Transaction{
			To:               ccipAddress.StringLong(),
			Data:             aptosmcms.ArgsToData(args),
			AdditionalFields: afBytes,
		},
	}

	b.Logger.Infow("Adding OnRamp destination config update operation",
		"chainCount", len(destChainSelectors))

	return []types.Operation{operation}, nil
}
