package operation

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	aptosutils "github.com/smartcontractkit/chainlink-aptos/relayer/utils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
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

func updateOffRampSources(b operations.Bundle, deps AptosDeps, in UpdateOffRampSourcesInput) ([]mcmstypes.Transaction, error) {
	var txs []mcmstypes.Transaction

	// Bind CCIP Package
	ccipAddress := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].CCIPAddress
	offRampBind := ccip_offramp.Bind(ccipAddress, deps.AptosChain.Client)

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
	moduleInfo, function, _, args, err := offRampBind.Offramp().Encoder().ApplySourceChainConfigUpdates(
		sourceChainSelectors,
		sourceChainEnabled,
		sourceChainRMNVerificationDisabled,
		sourceChainOnRamp,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode ApplySourceChainConfigUpdates for OffRamp: %w", err)
	}
	tx, err := utils.GenerateMCMSTx(ccipAddress, moduleInfo, function, args)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	txs = append(txs, tx)

	return txs, nil
}

// SetOcr3ConfigInput contains configuration for setting the OCR3 config on the OffRamp
type SetOcr3ConfigInput struct {
	OcrPluginType types.PluginType
	OCRConfigArgs internal.MultiOCR3BaseOCRConfigArgsAptos
}

// SetOcr3ConfigOp operation sets commit or exec OCR3 configuration for OffRamp
var SetOcr3ConfigOp = operations.NewOperation(
	"set-ocr3-config-op",
	Version1_0_0,
	"Sets OCR3 configuration for OffRamp",
	setOcr3Config,
)

func setOcr3Config(b operations.Bundle, deps AptosDeps, in SetOcr3ConfigInput) (mcmstypes.Transaction, error) {
	// Bind CCIP Package
	ccipAddress := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].CCIPAddress
	offRampBind := ccip_offramp.Bind(ccipAddress, deps.AptosChain.Client)

	var transmitters []aptos.AccountAddress
	for _, transmitter := range in.OCRConfigArgs.Transmitters {
		address, err := aptosutils.PublicKeyBytesToAddress(transmitter)
		if err != nil {
			return mcmstypes.Transaction{}, fmt.Errorf("failed to convert transmitter to address: %w", err)
		}
		transmitters = append(transmitters, address)
	}

	moduleInfo, function, _, args, err := offRampBind.Offramp().Encoder().SetOcr3Config(
		in.OCRConfigArgs.ConfigDigest[:],
		uint8(in.OcrPluginType),
		in.OCRConfigArgs.F,
		in.OCRConfigArgs.IsSignatureVerificationEnabled,
		in.OCRConfigArgs.Signers,
		transmitters,
	)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to encode SetOcr3Config for commit: %w", err)
	}
	tx, err := utils.GenerateMCMSTx(ccipAddress, moduleInfo, function, args)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to generate MCMS operations for OffRamp Initialize: %w", err)
	}

	return tx, nil
}
