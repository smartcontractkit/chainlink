package v1_6

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ccipseqs "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// UpdateBidirectionalLanesChangeset enables or disables multiple bidirectional lanes on CCIP.
// It batches all lane updates into a single MCMS proposal.
var UpdateBidirectionalLanesChangeset = cldf.CreateChangeSet(updateBidirectionalLanesLogic, updateBidirectionalLanesPrecondition)

// BidirectionalLaneDefinition indicates two chains that we want to connect.
type BidirectionalLaneDefinition struct {
	// IsDisabled indicates if the lane should be disabled.
	// We use IsDisabled instead of IsEnabled because enabling a lane should be the default action.
	IsDisabled bool
	Chains     [2]ChainDefinition
}

// laneDefinition defines a lane between source and destination.
type laneDefinition struct {
	// Source defines the source chain.
	Source ChainDefinition
	// Dest defines the destination chain.
	Dest ChainDefinition
}

// UpdateBidirectionalLanesConfig is a configuration struct for UpdateBidirectionalLanesChangeset.
type UpdateBidirectionalLanesConfig struct {
	// MCMSConfig defines the MCMS configuration for the changeset.
	MCMSConfig *proposalutils.TimelockConfig
	// Lanes describes the lanes that we want to create.
	Lanes []BidirectionalLaneDefinition
	// TestRouter indicates if we want to enable these lanes on the test router.
	TestRouter bool
}

type UpdateBidirectionalLanesChangesetConfigs struct {
	UpdateFeeQuoterDestsConfig  UpdateFeeQuoterDestsConfig
	UpdateFeeQuoterPricesConfig UpdateFeeQuoterPricesConfig
	UpdateOnRampDestsConfig     UpdateOnRampDestsConfig
	UpdateOffRampSourcesConfig  UpdateOffRampSourcesConfig
	UpdateRouterRampsConfig     UpdateRouterRampsConfig
}

func (c UpdateBidirectionalLanesConfig) BuildConfigs() UpdateBidirectionalLanesChangesetConfigs {
	onRampUpdatesByChain := make(map[uint64]map[uint64]OnRampDestinationUpdate)
	offRampUpdatesByChain := make(map[uint64]map[uint64]OffRampSourceUpdate)
	routerUpdatesByChain := make(map[uint64]RouterUpdates)
	feeQuoterDestUpdatesByChain := make(map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig)
	feeQuoterPriceUpdatesByChain := make(map[uint64]FeeQuoterPriceUpdatePerSource)

	for _, lane := range c.Lanes {
		isEnabled := !lane.IsDisabled
		chainA := lane.Chains[0]
		chainB := lane.Chains[1]

		laneAToB := laneDefinition{
			Source: chainA,
			Dest:   chainB,
		}
		laneBToA := laneDefinition{
			Source: chainB,
			Dest:   chainA,
		}

		for _, laneDef := range []laneDefinition{laneAToB, laneBToA} {
			// Setting the destination on the on ramp
			if onRampUpdatesByChain[laneDef.Source.Selector] == nil {
				onRampUpdatesByChain[laneDef.Source.Selector] = make(map[uint64]OnRampDestinationUpdate)
			}
			onRampUpdatesByChain[laneDef.Source.Selector][laneDef.Dest.Selector] = OnRampDestinationUpdate{
				IsEnabled:        isEnabled,
				TestRouter:       c.TestRouter,
				AllowListEnabled: laneDef.Dest.AllowListEnabled,
			}

			// Setting the source on the off ramp
			if offRampUpdatesByChain[laneDef.Dest.Selector] == nil {
				offRampUpdatesByChain[laneDef.Dest.Selector] = make(map[uint64]OffRampSourceUpdate)
			}
			offRampUpdatesByChain[laneDef.Dest.Selector][laneDef.Source.Selector] = OffRampSourceUpdate{
				IsEnabled:                 isEnabled,
				TestRouter:                c.TestRouter,
				IsRMNVerificationDisabled: laneDef.Source.RMNVerificationDisabled,
			}

			// Setting the on ramp on the source router
			routerUpdatesOnSource := routerUpdatesByChain[laneDef.Source.Selector]
			if routerUpdatesByChain[laneDef.Source.Selector].OnRampUpdates == nil {
				routerUpdatesOnSource.OnRampUpdates = make(map[uint64]bool)
			}
			routerUpdatesOnSource.OnRampUpdates[laneDef.Dest.Selector] = isEnabled
			routerUpdatesByChain[laneDef.Source.Selector] = routerUpdatesOnSource

			// Setting the off ramp on the dest router
			routerUpdatesOnDest := routerUpdatesByChain[laneDef.Dest.Selector]
			if routerUpdatesByChain[laneDef.Dest.Selector].OffRampUpdates == nil {
				routerUpdatesOnDest.OffRampUpdates = make(map[uint64]bool)
			}
			routerUpdatesOnDest.OffRampUpdates[laneDef.Source.Selector] = isEnabled
			routerUpdatesByChain[laneDef.Dest.Selector] = routerUpdatesOnDest

			// Setting the fee quoter destination on the source chain
			if feeQuoterDestUpdatesByChain[laneDef.Source.Selector] == nil {
				feeQuoterDestUpdatesByChain[laneDef.Source.Selector] = make(map[uint64]fee_quoter.FeeQuoterDestChainConfig)
			}
			feeQuoterDestUpdatesByChain[laneDef.Source.Selector][laneDef.Dest.Selector] = laneDef.Dest.FeeQuoterDestChainConfig

			// Setting the destination gas prices on the source chain
			feeQuoterPriceUpdatesOnSource := feeQuoterPriceUpdatesByChain[laneDef.Source.Selector]
			if feeQuoterPriceUpdatesOnSource.GasPrices == nil {
				feeQuoterPriceUpdatesOnSource.GasPrices = make(map[uint64]*big.Int)
			}
			feeQuoterPriceUpdatesOnSource.GasPrices[laneDef.Dest.Selector] = laneDef.Dest.GasPrice
			feeQuoterPriceUpdatesByChain[laneDef.Source.Selector] = feeQuoterPriceUpdatesOnSource
		}
	}

	routerMCMSConfig := c.MCMSConfig
	if c.TestRouter {
		routerMCMSConfig = nil // Test router is never owned by MCMS
	}

	return UpdateBidirectionalLanesChangesetConfigs{
		UpdateFeeQuoterDestsConfig: UpdateFeeQuoterDestsConfig{
			MCMS:           c.MCMSConfig,
			UpdatesByChain: feeQuoterDestUpdatesByChain,
		},
		UpdateFeeQuoterPricesConfig: UpdateFeeQuoterPricesConfig{
			MCMS:          c.MCMSConfig,
			PricesByChain: feeQuoterPriceUpdatesByChain,
		},
		UpdateOnRampDestsConfig: UpdateOnRampDestsConfig{
			MCMS:           c.MCMSConfig,
			UpdatesByChain: onRampUpdatesByChain,
		},
		UpdateOffRampSourcesConfig: UpdateOffRampSourcesConfig{
			MCMS:           c.MCMSConfig,
			UpdatesByChain: offRampUpdatesByChain,
		},
		UpdateRouterRampsConfig: UpdateRouterRampsConfig{
			TestRouter:     c.TestRouter,
			MCMS:           routerMCMSConfig,
			UpdatesByChain: routerUpdatesByChain,
		},
	}
}

func updateBidirectionalLanesPrecondition(e cldf.Environment, c UpdateBidirectionalLanesConfig) error {
	configs := c.BuildConfigs()

	return UpdateLanesPrecondition(e, configs)
}

func UpdateLanesPrecondition(e cldf.Environment, configs UpdateBidirectionalLanesChangesetConfigs) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	err = configs.UpdateFeeQuoterDestsConfig.Validate(e)
	if err != nil {
		return fmt.Errorf("failed to validate UpdateFeeQuoterDestsConfig: %w", err)
	}

	err = configs.UpdateFeeQuoterPricesConfig.Validate(e)
	if err != nil {
		return fmt.Errorf("failed to validate UpdateFeeQuoterPricesConfig: %w", err)
	}

	err = configs.UpdateOnRampDestsConfig.Validate(e)
	if err != nil {
		return fmt.Errorf("failed to validate UpdateOnRampDestsConfig: %w", err)
	}

	err = configs.UpdateOffRampSourcesConfig.Validate(e, state)
	if err != nil {
		return fmt.Errorf("failed to validate UpdateOffRampSourcesConfig: %w", err)
	}

	err = configs.UpdateRouterRampsConfig.Validate(e, state)
	if err != nil {
		return fmt.Errorf("failed to validate UpdateRouterRampsConfig: %w", err)
	}

	return nil
}

func updateBidirectionalLanesLogic(e cldf.Environment, c UpdateBidirectionalLanesConfig) (cldf.ChangesetOutput, error) {
	configs := c.BuildConfigs()

	return UpdateLanesLogic(e, c.MCMSConfig, configs)
}

// UpdateLanesLogic is the main logic for updating lanes. Configs provided can be unidirectional
// TODO: UpdateBidirectionalLanesChangesetConfigs name is misleading, it also accepts unidirectional lane updates
func UpdateLanesLogic(e cldf.Environment, mcmsConfig *proposalutils.TimelockConfig, configs UpdateBidirectionalLanesChangesetConfigs) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	report, err := operations.ExecuteSequence(e.OperationsBundle, ccipseqs.UpdateLanesSequence, e.BlockChains.EVMChains(), ccipseqs.UpdateLanesSequenceInput{
		FeeQuoterApplyDestChainConfigUpdatesSequenceInput: configs.UpdateFeeQuoterDestsConfig.ToSequenceInput(state),
		FeeQuoterUpdatePricesSequenceInput:                configs.UpdateFeeQuoterPricesConfig.ToSequenceInput(state),
		OffRampApplySourceChainConfigUpdatesSequenceInput: configs.UpdateOffRampSourcesConfig.ToSequenceInput(state),
		OnRampApplyDestChainConfigUpdatesSequenceInput:    configs.UpdateOnRampDestsConfig.ToSequenceInput(state),
		RouterApplyRampUpdatesSequenceInput:               configs.UpdateRouterRampsConfig.ToSequenceInput(state),
	})
	return opsutil.AddEVMCallSequenceToCSOutput(
		e,
		state,
		cldf.ChangesetOutput{},
		report,
		err,
		mcmsConfig,
		"Update lanes on CCIP 1.6.0",
	)
}
