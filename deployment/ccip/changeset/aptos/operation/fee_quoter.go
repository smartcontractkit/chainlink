package operation

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	aptos_fee_quoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	"github.com/smartcontractkit/chainlink/deployment/operations"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	"github.com/smartcontractkit/mcms/types"
)

// UpdateFeeQuoterDestsInput contains configuration for updating FeeQuoter destination configs
type UpdateFeeQuoterDestsInput struct {
	MCMSAddress aptos.AccountAddress
	Updates     map[uint64]aptos_fee_quoter.DestChainConfig
}

// UpdateFeeQuoterDestsOp operation to update FeeQuoter destination configurations
var UpdateFeeQuoterDestsOp = operations.NewOperation(
	"update-fee-quoter-dests-op",
	Version1_0_0,
	"Updates FeeQuoter destination chain configurations",
	updateFeeQuoterDests,
)

func updateFeeQuoterDests(b operations.Bundle, deps AptosDeps, in UpdateFeeQuoterDestsInput) ([]types.Transaction, error) {
	// Bind CCIP Package
	ccipAddress := deps.OnChainState.CCIPAddress
	ccipBind := ccip.Bind(ccipAddress, deps.AptosChain.Client)

	// Process each destination chain config update
	var txs []types.Transaction

	for destChainSelector, destConfig := range in.Updates {
		// Encode the update operation
		moduleInfo, function, _, args, err := ccipBind.FeeQuoter().Encoder().ApplyDestChainConfigUpdates(
			destChainSelector,
			destConfig.IsEnabled,
			destConfig.MaxNumberOfTokensPerMsg,
			destConfig.MaxDataBytes,
			destConfig.MaxPerMsgGasLimit,
			destConfig.DestGasOverhead,
			destConfig.DestGasPerPayloadByteBase,
			destConfig.DestGasPerPayloadByteHigh,
			destConfig.DestGasPerPayloadByteThreshold,
			destConfig.DestDataAvailabilityOverheadGas,
			destConfig.DestGasPerDataAvailabilityByte,
			destConfig.DestDataAvailabilityMultiplierBps,
			destConfig.ChainFamilySelector,
			destConfig.EnforceOutOfOrder,
			destConfig.DefaultTokenFeeUsdCents,
			destConfig.DefaultTokenDestGasOverhead,
			destConfig.DefaultTxGasLimit,
			destConfig.GasMultiplierWeiPerEth,
			destConfig.GasPriceStalenessThreshold,
			destConfig.NetworkFeeUsdCents,
		)
		if err != nil {
			return []types.Transaction{}, fmt.Errorf("failed to encode ApplyDestChainConfigUpdates for chain %d: %w", destChainSelector, err)
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

		b.Logger.Infow("Adding FeeQuoter destination config update operation",
			"destChainSelector", destChainSelector,
			"isEnabled", destConfig.IsEnabled)
	}

	return txs, nil
}

// UpdateFeeQuoterPricesInput contains configuration for updating FeeQuoter price configs
type UpdateFeeQuoterPricesInput struct {
	MCMSAddress aptos.AccountAddress
	Prices      FeeQuoterPriceUpdatePerSource
}

type FeeQuoterPriceUpdatePerSource struct {
	TokenPrices map[string]*big.Int // token address (string) -> price
	GasPrices   map[uint64]*big.Int // dest chain -> gas price
}

// UpdateFeeQuoterPricesOp operation to update FeeQuoter prices
var UpdateFeeQuoterPricesOp = operations.NewOperation(
	"update-fee-quoter-prices-op",
	Version1_0_0,
	"Updates FeeQuoter token and gas prices",
	updateFeeQuoterPrices,
)

func updateFeeQuoterPrices(b operations.Bundle, deps AptosDeps, in UpdateFeeQuoterPricesInput) ([]types.Transaction, error) {
	var txs []types.Transaction

	// Bind CCIP Package
	ccipAddress := deps.OnChainState.CCIPAddress
	ccipBind := ccip.Bind(ccipAddress, deps.AptosChain.Client)

	// Bind MCMS Package
	mcmsAddress := deps.OnChainState.MCMSAddress
	mcmsBind := mcms.Bind(mcmsAddress, deps.AptosChain.Client)

	// Add CCIP Owner address to update token prices allow list
	ccipOwnerAddress, err := mcmsBind.MCMSRegistry().GetRegisteredOwnerAddress(nil, ccipAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get CCIP owner address: %w", err)
	}
	moduleInfo, function, _, args, err := ccipBind.Auth().Encoder().ApplyAllowedOfframpUpdates(nil, []aptos.AccountAddress{ccipOwnerAddress})
	if err != nil {
		return nil, fmt.Errorf("failed to encode ApplyAllowedOfframpUpdates: %w", err)
	}
	additionalFields := aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}
	afBytes, err := json.Marshal(additionalFields)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal additional fields: %w", err)
	}

	txs = append(txs, types.Transaction{
		To:               ccipAddress.StringLong(),
		Data:             aptosmcms.ArgsToData(args),
		AdditionalFields: afBytes,
	})

	// Convert token prices and gas prices to format expected by Aptos contract
	var sourceTokens []aptos.AccountAddress
	var sourceUsdPerToken []*big.Int
	var gasDestChainSelectors []uint64
	var gasUsdPerUnitGas []*big.Int

	// Process token prices if any
	for tokenAddr, price := range in.Prices.TokenPrices {
		address := aptos.AccountAddress{}
		err := address.ParseStringRelaxed(tokenAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Aptos token address %s: %w", tokenAddr, err)
		}
		sourceTokens = append(sourceTokens, address)
		sourceUsdPerToken = append(sourceUsdPerToken, price)
	}

	// Process gas prices if any
	for destChainSel, gasPrice := range in.Prices.GasPrices {
		gasDestChainSelectors = append(gasDestChainSelectors, destChainSel)
		gasUsdPerUnitGas = append(gasUsdPerUnitGas, gasPrice)
	}

	// Generate MCMS tx to update prices
	if len(sourceTokens) == 0 && len(gasDestChainSelectors) == 0 {
		b.Logger.Infow("No price updates to apply")
		return txs, nil
	}

	// Encode the update tx
	moduleInfo, function, _, args, err = ccipBind.FeeQuoter().Encoder().UpdatePrices(
		sourceTokens,
		sourceUsdPerToken,
		gasDestChainSelectors,
		gasUsdPerUnitGas,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode UpdatePrices: %w", err)
	}

	additionalFields = aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}
	afBytes, err = json.Marshal(additionalFields)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal additional fields: %w", err)
	}

	txs = append(txs, types.Transaction{
		To:               ccipAddress.StringLong(),
		Data:             aptosmcms.ArgsToData(args),
		AdditionalFields: afBytes,
	})

	b.Logger.Infow("Adding FeeQuoter price update operation",
		"tokenPriceCount", len(sourceTokens),
		"gasPriceCount", len(gasDestChainSelectors),
	)

	return txs, nil
}
