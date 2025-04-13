package operation

import (
	"encoding/json"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/operations"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	"github.com/smartcontractkit/mcms/types"
)

// UpdateOffRampSourcesInput contains configuration for updating OffRamp sources
type UpdateOffRampSourcesInput struct {
	MCMSAddress aptos.AccountAddress
	Updates     map[uint64]v1_6.OffRampSourceUpdate
}

// UpdateOffRampSourcesOp operation to update OffRamp source configurations
var UpdateOffRampSourcesOp = operations.NewOperation(
	"update-offramp-sources-op",
	Version1_0_0,
	"Updates OffRamp source chain configurations",
	updateOffRampSources,
)

func updateOffRampSources(b operations.Bundle, deps AptosDeps, in UpdateOffRampSourcesInput) ([]types.Transaction, error) {
	// Bind CCIP Package
	ccipAddress := deps.OnChainState.CCIPAddress
	offrampBind := ccip_offramp.Bind(ccipAddress, deps.AptosChain.Client)

	// Transform the updates into the format expected by the Aptos contract
	var sourceChainSelectors []uint64
	var sourceChainEnabled []bool
	var sourceChainRMNVerificationDisabled []bool
	var sourceChainOnRamp [][]byte

	for sourceChainSelector, update := range in.Updates {
		sourceChainSelectors = append(sourceChainSelectors, sourceChainSelector)
		sourceChainEnabled = append(sourceChainEnabled, update.IsEnabled)
		sourceChainRMNVerificationDisabled = append(sourceChainRMNVerificationDisabled, update.IsRMNVerificationDisabled)

		onRampBytes, err := deps.CCIPOnChainState.GetOnRampAddressBytes(sourceChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to get onRamp address for source chain %d: %w", sourceChainSelector, err)
		}
		sourceChainOnRamp = append(sourceChainOnRamp, onRampBytes)
	}

	if len(sourceChainSelectors) == 0 {
		b.Logger.Infow("No OffRamp source updates to apply")
		return nil, nil
	}

	// Encode the update operation
	moduleInfo, function, _, args, err := offrampBind.Offramp().Encoder().ApplySourceChainConfigUpdates(
		sourceChainSelectors,
		sourceChainEnabled,
		sourceChainRMNVerificationDisabled,
		sourceChainOnRamp,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode ApplySourceChainConfigUpdates for OffRamp: %w", err)
	}

	// Create MCMS operation
	additionalFields := aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}
	afBytes, err := json.Marshal(additionalFields)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal additional fields: %w", err)
	}

	b.Logger.Infow("Adding OffRamp source config update operation",
		"chainCount", len(sourceChainSelectors))

	return []types.Transaction{{
		To:               ccipAddress.StringLong(),
		Data:             aptosmcms.ArgsToData(args),
		AdditionalFields: afBytes,
	}}, nil
}
