package v1_6

import (
	"fmt"
	"math/big"
	"slices"

	"github.com/Masterminds/semver/v3"

	"github.com/ethereum/go-ethereum/common"

	fqv2ops "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/operations/fee_quoter"
	fqv2seq "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/sequences"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ccipseqs "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"

	mcmslib "github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
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

	ds, err := shared.PopulateDataStore(e.ExistingAddresses)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to populate datastore from existing addresses: %w", err)
	}
	// Merge the environment's DataStore as feequoter v2 is only available in DS
	if e.DataStore != nil {
		if err := ds.Merge(e.DataStore); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge environment datastore: %w", err)
		}
	}

	feeQuoterDestsInput := configs.UpdateFeeQuoterDestsConfig.ToSequenceInput(state)
	feeQuoterPricesInput := configs.UpdateFeeQuoterPricesConfig.ToSequenceInput(state)
	feeQuoterVersionsByChain, err := resolveFeeQuoterTargets(ds, &feeQuoterDestsInput, &feeQuoterPricesInput)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	v2FeeQuoterChains := make(map[uint64]struct{})
	v1FeeQuoterDestsUpdates := make(map[uint64]opsutil.EVMCallInput[[]fee_quoter.FeeQuoterDestChainConfigArgs])
	v1FeeQuoterPriceUpdates := make(map[uint64]opsutil.EVMCallInput[fee_quoter.InternalPriceUpdates])

	for chainSel, update := range feeQuoterDestsInput.UpdatesByChain {
		version, ok := feeQuoterVersionsByChain[chainSel]
		if ok && version.Major() >= 2 {
			v2FeeQuoterChains[chainSel] = struct{}{}
			continue
		}
		v1FeeQuoterDestsUpdates[chainSel] = update
	}
	for chainSel, update := range feeQuoterPricesInput.UpdatesByChain {
		version, ok := feeQuoterVersionsByChain[chainSel]
		if ok && version.Major() >= 2 {
			v2FeeQuoterChains[chainSel] = struct{}{}
			continue
		}
		v1FeeQuoterPriceUpdates[chainSel] = update
	}

	report, err := operations.ExecuteSequence(e.OperationsBundle, ccipseqs.UpdateLanesSequence, e.BlockChains.EVMChains(), ccipseqs.UpdateLanesSequenceInput{
		FeeQuoterApplyDestChainConfigUpdatesSequenceInput: ccipseqs.FeeQuoterApplyDestChainConfigUpdatesSequenceInput{
			UpdatesByChain: v1FeeQuoterDestsUpdates,
		},
		FeeQuoterUpdatePricesSequenceInput: ccipseqs.FeeQuoterUpdatePricesSequenceInput{
			UpdatesByChain: v1FeeQuoterPriceUpdates,
		},
		OffRampApplySourceChainConfigUpdatesSequenceInput: configs.UpdateOffRampSourcesConfig.ToSequenceInput(state),
		OnRampApplyDestChainConfigUpdatesSequenceInput:    configs.UpdateOnRampDestsConfig.ToSequenceInput(state),
		RouterApplyRampUpdatesSequenceInput:               configs.UpdateRouterRampsConfig.ToSequenceInput(state),
	})
	output, err := opsutil.AddEVMCallSequenceToCSOutput(
		e,
		cldf.ChangesetOutput{},
		report,
		err,
		state.EVMMCMSStateByChain(),
		mcmsConfig,
		"Update lanes on CCIP",
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	if len(v2FeeQuoterChains) == 0 {
		return output, nil
	}

	// Execute v2 FeeQuoter update sequences and append their batch
	v2ChainSels := make([]uint64, 0, len(v2FeeQuoterChains))
	for chainSel := range v2FeeQuoterChains {
		v2ChainSels = append(v2ChainSels, chainSel)
	}
	slices.Sort(v2ChainSels)

	var v2BatchOps []mcmstypes.BatchOperation
	for _, chainSel := range v2ChainSels {
		fqUpdate := fqv2seq.FeeQuoterUpdate{
			ChainSelector:     chainSel,
			ExistingAddresses: ds.Addresses().Filter(datastore.AddressRefByChainSelector(chainSel)),
		}
		if dests, ok := feeQuoterDestsInput.UpdatesByChain[chainSel]; ok {
			destCfgs, err := ConvertV16FeeQuoterDestUpdatesToV2(dests.CallInput)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert v1.6 fee quoter destination updates for chain %d: %w", chainSel, err)
			}
			fqUpdate.DestChainConfigs = destCfgs
		}
		if prices, ok := feeQuoterPricesInput.UpdatesByChain[chainSel]; ok {
			fqUpdate.PriceUpdates = ConvertV16FeeQuoterPriceUpdatesToV2(prices.CallInput)
		}

		v2Report, err := operations.ExecuteSequence(
			e.OperationsBundle,
			fqv2seq.SequenceFeeQuoterUpdate,
			e.BlockChains,
			fqUpdate,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute v2 FeeQuoter update sequence on chain %d: %w", chainSel, err)
		}
		output.Reports = append(output.Reports, v2Report.ExecutionReports...)
		v2BatchOps = append(v2BatchOps, v2Report.Output.BatchOps...)
	}

	if mcmsConfig == nil || len(v2BatchOps) == 0 {
		return output, nil
	}

	output.MCMSTimelockProposals = append(output.MCMSTimelockProposals, mcmslib.TimelockProposal{
		Operations: v2BatchOps,
	})
	aggProposal, err := proposalutils.AggregateProposalsV2(
		e,
		proposalutils.MCMSStates{MCMSEVMState: state.EVMMCMSStateByChain()},
		output.MCMSTimelockProposals,
		"Update lanes on CCIP",
		mcmsConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to aggregate MCMS proposals: %w", err)
	}
	output.MCMSTimelockProposals = []mcmslib.TimelockProposal{*aggProposal}

	return output, nil
}

func ConvertV16FeeQuoterDestUpdatesToV2(in []fee_quoter.FeeQuoterDestChainConfigArgs) ([]fqv2ops.DestChainConfigArgs, error) {
	out := make([]fqv2ops.DestChainConfigArgs, 0, len(in))
	for _, cfg := range in {
		if cfg.DestChainConfig.NetworkFeeUSDCents > uint32(^uint16(0)) {
			return nil, fmt.Errorf(
				"network fee USD cents %d for destination chain %d exceeds uint16 max",
				cfg.DestChainConfig.NetworkFeeUSDCents,
				cfg.DestChainSelector,
			)
		}
		out = append(out, fqv2ops.DestChainConfigArgs{
			DestChainSelector: cfg.DestChainSelector,
			DestChainConfig: fqv2ops.DestChainConfig{
				IsEnabled:                   cfg.DestChainConfig.IsEnabled,
				MaxDataBytes:                cfg.DestChainConfig.MaxDataBytes,
				MaxPerMsgGasLimit:           cfg.DestChainConfig.MaxPerMsgGasLimit,
				DestGasOverhead:             cfg.DestChainConfig.DestGasOverhead,
				DestGasPerPayloadByteBase:   cfg.DestChainConfig.DestGasPerPayloadByteBase,
				ChainFamilySelector:         cfg.DestChainConfig.ChainFamilySelector,
				DefaultTokenFeeUSDCents:     cfg.DestChainConfig.DefaultTokenFeeUSDCents,
				DefaultTokenDestGasOverhead: cfg.DestChainConfig.DefaultTokenDestGasOverhead,
				DefaultTxGasLimit:           cfg.DestChainConfig.DefaultTxGasLimit,
				NetworkFeeUSDCents:          uint16(cfg.DestChainConfig.NetworkFeeUSDCents), //nolint:gosec // value is range-checked above
				LinkFeeMultiplierPercent:    fqv2seq.LinkFeeMultiplierPercent,
			},
		})
	}
	return out, nil
}

func ConvertV16FeeQuoterPriceUpdatesToV2(in fee_quoter.InternalPriceUpdates) fqv2ops.PriceUpdates {
	out := fqv2ops.PriceUpdates{
		TokenPriceUpdates: make([]fqv2ops.TokenPriceUpdate, 0, len(in.TokenPriceUpdates)),
		GasPriceUpdates:   make([]fqv2ops.GasPriceUpdate, 0, len(in.GasPriceUpdates)),
	}
	for _, tokenPrice := range in.TokenPriceUpdates {
		out.TokenPriceUpdates = append(out.TokenPriceUpdates, fqv2ops.TokenPriceUpdate{
			SourceToken: tokenPrice.SourceToken,
			UsdPerToken: tokenPrice.UsdPerToken,
		})
	}
	for _, gasPrice := range in.GasPriceUpdates {
		out.GasPriceUpdates = append(out.GasPriceUpdates, fqv2ops.GasPriceUpdate{
			DestChainSelector: gasPrice.DestChainSelector,
			UsdPerUnitGas:     gasPrice.UsdPerUnitGas,
		})
	}
	return out
}

func resolveFeeQuoterTargets(
	ds *datastore.MemoryDataStore,
	destsInput *ccipseqs.FeeQuoterApplyDestChainConfigUpdatesSequenceInput,
	pricesInput *ccipseqs.FeeQuoterUpdatePricesSequenceInput,
) (map[uint64]semver.Version, error) {
	versionsByChain := make(map[uint64]semver.Version)

	resolve := func(chainSel uint64) error {
		if _, ok := versionsByChain[chainSel]; ok {
			return nil
		}
		chainAddresses := ds.Addresses().Filter(datastore.AddressRefByChainSelector(chainSel))
		addr, version, err := resolveUpdateLanesFeeQuoterAddressAndVersion(chainAddresses, chainSel)
		if err != nil {
			return fmt.Errorf("failed to resolve FeeQuoter target on chain %d: %w", chainSel, err)
		}
		versionsByChain[chainSel] = version

		if update, ok := destsInput.UpdatesByChain[chainSel]; ok {
			update.Address = addr
			destsInput.UpdatesByChain[chainSel] = update
		}
		if update, ok := pricesInput.UpdatesByChain[chainSel]; ok {
			update.Address = addr
			pricesInput.UpdatesByChain[chainSel] = update
		}
		return nil
	}

	for chainSel := range destsInput.UpdatesByChain {
		if err := resolve(chainSel); err != nil {
			return nil, err
		}
	}
	for chainSel := range pricesInput.UpdatesByChain {
		if err := resolve(chainSel); err != nil {
			return nil, err
		}
	}

	return versionsByChain, nil
}

func resolveUpdateLanesFeeQuoterAddressAndVersion(
	addresses []datastore.AddressRef,
	chainSel uint64,
) (common.Address, semver.Version, error) {
	// Find the FeeQuoter with the highest version for this chain
	var bestRef datastore.AddressRef
	var bestVersion *semver.Version

	for _, ref := range addresses {
		if ref.ChainSelector != chainSel {
			continue
		}
		if ref.Type != datastore.ContractType(fqv2ops.ContractType) {
			continue
		}
		if ref.Version == nil {
			continue
		}
		if bestVersion == nil || ref.Version.GreaterThan(bestVersion) {
			bestVersion = ref.Version
			bestRef = ref
		}
	}

	if bestVersion == nil {
		return common.Address{}, semver.Version{}, fmt.Errorf("no fee quoter address found for chain %d", chainSel)
	}

	if !common.IsHexAddress(bestRef.Address) {
		return common.Address{}, semver.Version{}, fmt.Errorf("invalid fee quoter address %q for chain %d", bestRef.Address, chainSel)
	}

	return common.HexToAddress(bestRef.Address), *bestVersion, nil
}
