package operation

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	aptos_fee_quoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
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

	// TODO
	multipliers := make([]uint64, len(sourceTokens))
	for i := range multipliers {
		multipliers[i] = uint64(1)
	}
	moduleInfo, function, _, args, err = ccipBind.FeeQuoter().Encoder().ApplyPremiumMultiplierWeiPerEthUpdates(sourceTokens, multipliers)
	if err != nil {
		return nil, fmt.Errorf("failed to encode ApplyPremiumMultiplierWeiPerEth: %w", err)
	}
	tx, err := aptosmcms.NewTransaction(
		moduleInfo.PackageName,
		moduleInfo.ModuleName,
		function,
		ccipAddress,
		aptosmcms.ArgsToData(args),
		"FeeQuoter",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate transaction: %w", err)
	}
	txs = append(txs, tx)

	return txs, nil
}

// UpdateTokenTransferFeeConfigsInput ...
type UpdateTokenTransferFeeConfigsInput struct {
	MCMSAddress        aptos.AccountAddress
	AddTokenConfigs    []aptos_fee_quoter.TokenTransferFeeConfigAdded
	RemoveTokenConfigs []aptos_fee_quoter.TokenTransferFeeConfigRemoved
}

// UpdateTokenTransferCfgOp operation to update FeeQuoter prices
var UpdateTokenTransferCfgOp = operations.NewOperation(
	"update-token-transfer-config-op",
	Version1_0_0,
	"Updates Token Transfer Fee Configs",
	updateTokenTransferCfg,
)

func updateTokenTransferCfg(b operations.Bundle, deps AptosDeps, in UpdateTokenTransferFeeConfigsInput) ([]types.Transaction, error) {
	var txs []types.Transaction

	// Bind CCIP Package
	ccipAddress := deps.OnChainState.CCIPAddress
	ccipBind := ccip.Bind(ccipAddress, deps.AptosChain.Client)

	// Encode the update tx
	argsByDestChain := toApplyTokenTransferFeeConfigUpdatesArgs(in.AddTokenConfigs, in.RemoveTokenConfigs)
	for destChainSelector, args := range argsByDestChain {
		moduleInfo, function, _, args, err := ccipBind.FeeQuoter().Encoder().ApplyTokenTransferFeeConfigUpdates(
			destChainSelector,
			args.addTokens,
			args.addMinFeeUsdCents,
			args.addMaxFeeUsdCents,
			args.addDeciBps,
			args.addDestGasOverhead,
			args.addDestBytesOverhead,
			args.addIsEnabled,
			args.removeTokens,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to encode UpdatePrices: %w", err)
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
		b.Logger.Infow("Adding TokenTransferFeeConfig update transaction", "destChainSelector", destChainSelector)
	}

	return txs, nil
}

type applyTokenTransferFeeConfigUpdatesArgs struct {
	destChainSelector    uint64
	addTokens            []aptos.AccountAddress
	addMinFeeUsdCents    []uint32
	addMaxFeeUsdCents    []uint32
	addDeciBps           []uint16
	addDestGasOverhead   []uint32
	addDestBytesOverhead []uint32
	addIsEnabled         []bool
	removeTokens         []aptos.AccountAddress
}

func toApplyTokenTransferFeeConfigUpdatesArgs(
	AddTokenConfigs []aptos_fee_quoter.TokenTransferFeeConfigAdded,
	RemovTokenConfigs []aptos_fee_quoter.TokenTransferFeeConfigRemoved,
) map[uint64]applyTokenTransferFeeConfigUpdatesArgs {
	argsByDestChain := make(map[uint64]applyTokenTransferFeeConfigUpdatesArgs)

	// Process added token configs
	for _, config := range AddTokenConfigs {
		args := argsByDestChain[config.DestChainSelector]

		args.destChainSelector = config.DestChainSelector
		args.addTokens = append(args.addTokens, config.Token)
		args.addMinFeeUsdCents = append(args.addMinFeeUsdCents, config.TokenTransferFeeConfig.MinFeeUsdCents)
		args.addMaxFeeUsdCents = append(args.addMaxFeeUsdCents, config.TokenTransferFeeConfig.MaxFeeUsdCents)
		args.addDeciBps = append(args.addDeciBps, config.TokenTransferFeeConfig.DeciBps)
		args.addDestGasOverhead = append(args.addDestGasOverhead, config.TokenTransferFeeConfig.DestGasOverhead)
		args.addDestBytesOverhead = append(args.addDestBytesOverhead, config.TokenTransferFeeConfig.DestBytesOverhead)
		args.addIsEnabled = append(args.addIsEnabled, config.TokenTransferFeeConfig.IsEnabled)

		argsByDestChain[config.DestChainSelector] = args
	}

	// Process removed token configs
	for _, config := range RemovTokenConfigs {
		args := argsByDestChain[config.DestChainSelector]

		args.destChainSelector = config.DestChainSelector
		args.removeTokens = append(args.removeTokens, config.Token)

		argsByDestChain[config.DestChainSelector] = args
	}

	return argsByDestChain
}
