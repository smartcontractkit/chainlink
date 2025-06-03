package operation

import (
	"fmt"
	"math/big"

	"github.com/aptos-labs/aptos-go-sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	aptos_fee_quoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
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

func updateFeeQuoterDests(b operations.Bundle, deps AptosDeps, in UpdateFeeQuoterDestsInput) ([]mcmstypes.Transaction, error) {
	// Bind CCIP Package
	ccipAddress := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].CCIPAddress
	ccipBind := ccip.Bind(ccipAddress, deps.AptosChain.Client)

	// Process each destination chain config update
	var txs []mcmstypes.Transaction

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
			return nil, fmt.Errorf("failed to encode ApplyDestChainConfigUpdates for chain %d: %w", destChainSelector, err)
		}

		tx, err := utils.GenerateMCMSTx(ccipAddress, moduleInfo, function, args)
		if err != nil {
			return nil, fmt.Errorf("failed to create transaction: %w", err)
		}

		txs = append(txs, tx)

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

func updateFeeQuoterPrices(b operations.Bundle, deps AptosDeps, in UpdateFeeQuoterPricesInput) ([]mcmstypes.Transaction, error) {
	var txs []mcmstypes.Transaction

	// Bind CCIP Package
	ccipAddress := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].CCIPAddress
	ccipBind := ccip.Bind(ccipAddress, deps.AptosChain.Client)

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
		return nil, nil
	}

	// Encode the update tx
	moduleInfo, function, _, args, err := ccipBind.FeeQuoter().Encoder().UpdatePrices(
		sourceTokens,
		sourceUsdPerToken,
		gasDestChainSelectors,
		gasUsdPerUnitGas,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode UpdatePrices: %w", err)
	}
	tx, err := utils.GenerateMCMSTx(ccipAddress, moduleInfo, function, args)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	txs = append(txs, tx)

	b.Logger.Infow("Adding FeeQuoter price update operation",
		"tokenPriceCount", len(sourceTokens),
		"gasPriceCount", len(gasDestChainSelectors),
	)

	return txs, nil
}
