package sequence

import (
	"fmt"
	"math/big"

	"github.com/aptos-labs/aptos-go-sdk"
	chainsel "github.com/smartcontractkit/chain-selectors"
	aptos_fee_quoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	config "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/mcms/types"
)

type UpdateAptosLanesSeqInput struct {
	UpdateFeeQuoterDestsConfig  operation.UpdateFeeQuoterDestsInput
	UpdateFeeQuoterPricesConfig operation.UpdateFeeQuoterPricesInput
	UpdateOnRampDestsConfig     operation.UpdateOnRampDestsInput
	UpdateOffRampSourcesConfig  operation.UpdateOffRampSourcesInput
}

// UpdateAptosLanesSequence orchestrates operations to update Aptos lanes
var UpdateAptosLanesSequence = operations.NewSequence(
	"update-aptos-lanes-sequence",
	operation.Version1_0_0,
	"Update Aptos CCIP lanes to enable or disable connections",
	updateAptosLanesSequence,
)

func updateAptosLanesSequence(b operations.Bundle, deps operation.AptosDeps, in UpdateAptosLanesSeqInput) (types.BatchOperation, error) {
	var mcmsTxs []types.Transaction

	// 1. Update FeeQuoters with destination configs
	b.Logger.Info("Updating destination configs on FeeQuoters")
	feeQuoterDestReport, err := operations.ExecuteOperation(b, operation.UpdateFeeQuoterDestsOp, deps, in.UpdateFeeQuoterDestsConfig)
	if err != nil {
		return types.BatchOperation{}, fmt.Errorf("failed to update FeeQuoter destinations: %w", err)
	}
	mcmsTxs = append(mcmsTxs, feeQuoterDestReport.Output...)

	// 2. Configure destinations on OnRamps
	b.Logger.Info("Updating destination configs on OnRamps")
	onRampReport, err := operations.ExecuteOperation(b, operation.UpdateOnRampDestsOp, deps, in.UpdateOnRampDestsConfig)
	if err != nil {
		return types.BatchOperation{}, fmt.Errorf("failed to update OnRamp destinations: %w", err)
	}
	mcmsTxs = append(mcmsTxs, onRampReport.Output...)

	// 3. Configure sources on OffRamps
	b.Logger.Info("Updating source configs on OffRamps")
	offRampReport, err := operations.ExecuteOperation(b, operation.UpdateOffRampSourcesOp, deps, in.UpdateOffRampSourcesConfig)
	if err != nil {
		return types.BatchOperation{}, fmt.Errorf("failed to update OffRamp sources: %w", err)
	}
	mcmsTxs = append(mcmsTxs, offRampReport.Output...)

	// TODO: This is not working
	// 4. Update FeeQuoters with gas prices
	b.Logger.Info("Updating gas prices on FeeQuoters")
	feeQuoterPricesReport, err := operations.ExecuteOperation(b, operation.UpdateFeeQuoterPricesOp, deps, in.UpdateFeeQuoterPricesConfig)
	if err != nil {
		return types.BatchOperation{}, fmt.Errorf("failed to update FeeQuoter prices: %w", err)
	}
	mcmsTxs = append(mcmsTxs, feeQuoterPricesReport.Output...)

	return types.BatchOperation{
		ChainSelector: types.ChainSelector(deps.AptosChain.Selector),
		Transactions:  mcmsTxs,
	}, nil
}

// Convert config.UpdateAptosLanesConfig into a map[uint64]UpdateAptosLanesSeqInput
func ConvertToUpdateAptosLanesSeqInput(aptosChains map[uint64]changeset.AptosCCIPChainState, cfg config.UpdateAptosLanesConfig) map[uint64]UpdateAptosLanesSeqInput {
	updateInputsByAptosChain := make(map[uint64]UpdateAptosLanesSeqInput)

	// Group the operations by Aptos chain
	for _, lane := range cfg.Lanes {
		// Process lanes with Aptos as the source chain
		if lane.Source.GetChainFamily() == chainsel.FamilyAptos {
			source := lane.Source.(config.AptosChainDefinition)
			if _, exists := updateInputsByAptosChain[source.Selector]; !exists {
				updateInputsByAptosChain[source.Selector] = UpdateAptosLanesSeqInput{}
			}
			mcmsAddress := aptosChains[source.Selector].MCMSAddress
			setAptosSourceUpdates(lane, updateInputsByAptosChain, cfg.TestRouter, mcmsAddress)
		}

		// Process lanes with Aptos as the destination chain
		if lane.Dest.GetChainFamily() == chainsel.FamilyAptos {
			dest := lane.Dest.(config.AptosChainDefinition)
			if _, exists := updateInputsByAptosChain[dest.Selector]; !exists {
				updateInputsByAptosChain[dest.Selector] = UpdateAptosLanesSeqInput{}
			}
			mcmsAddress := aptosChains[dest.Selector].MCMSAddress
			setAptosDestinationUpdates(lane, updateInputsByAptosChain, cfg.TestRouter, mcmsAddress)
		}
	}

	return updateInputsByAptosChain
}

func setAptosSourceUpdates(lane config.LaneConfig, updateInputsByAptosChain map[uint64]UpdateAptosLanesSeqInput, isTestRouter bool, mcmsAddress aptos.AccountAddress) {
	source := lane.Source.(config.AptosChainDefinition)
	dest := lane.Dest.(config.EVMChainDefinition)
	isEnabled := !lane.IsDisabled

	// Setting the destination on the on ramp
	input := updateInputsByAptosChain[source.Selector]
	input.UpdateOnRampDestsConfig.MCMSAddress = mcmsAddress
	if input.UpdateOnRampDestsConfig.Updates == nil {
		input.UpdateOnRampDestsConfig.Updates = make(map[uint64]v1_6.OnRampDestinationUpdate)
	}
	input.UpdateOnRampDestsConfig.Updates[dest.Selector] = v1_6.OnRampDestinationUpdate{
		IsEnabled:        isEnabled,
		TestRouter:       isTestRouter,
		AllowListEnabled: dest.AllowListEnabled,
	}

	// Setting gas prices updates
	input.UpdateFeeQuoterPricesConfig.MCMSAddress = mcmsAddress
	if input.UpdateFeeQuoterPricesConfig.Prices.GasPrices == nil {
		input.UpdateFeeQuoterPricesConfig.Prices.GasPrices = make(map[uint64]*big.Int)
	}
	input.UpdateFeeQuoterPricesConfig.Prices.GasPrices[dest.Selector] = dest.GasPrice

	// Setting the fee quoter destination on the source chain
	input.UpdateFeeQuoterDestsConfig.MCMSAddress = mcmsAddress
	if input.UpdateFeeQuoterDestsConfig.Updates == nil {
		input.UpdateFeeQuoterDestsConfig.Updates = make(map[uint64]aptos_fee_quoter.DestChainConfig)
	}
	input.UpdateFeeQuoterDestsConfig.Updates[dest.Selector] = dest.GetConvertedAptosFeeQuoterConfig()

	updateInputsByAptosChain[source.Selector] = input
}

func setAptosDestinationUpdates(lane config.LaneConfig, updateInputsByAptosChain map[uint64]UpdateAptosLanesSeqInput, isTestRouter bool, mcmsAddress aptos.AccountAddress) {
	source := lane.Source.(config.EVMChainDefinition)
	dest := lane.Dest.(config.AptosChainDefinition)
	isEnabled := !lane.IsDisabled

	// Setting off ramp updates
	input := updateInputsByAptosChain[dest.Selector]
	input.UpdateOffRampSourcesConfig.MCMSAddress = mcmsAddress
	if input.UpdateOffRampSourcesConfig.Updates == nil {
		input.UpdateOffRampSourcesConfig.Updates = make(map[uint64]v1_6.OffRampSourceUpdate)
	}
	input.UpdateOffRampSourcesConfig.Updates[source.Selector] = v1_6.OffRampSourceUpdate{
		IsEnabled:                 isEnabled,
		TestRouter:                isTestRouter,
		IsRMNVerificationDisabled: source.RMNVerificationDisabled,
	}

	updateInputsByAptosChain[dest.Selector] = input
}
