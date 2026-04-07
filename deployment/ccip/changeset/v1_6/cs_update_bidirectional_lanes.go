package v1_6

import (
	"fmt"
	"math/big"
	"slices"

	"github.com/Masterminds/semver/v3"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	fqv2ops "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/operations/fee_quoter"
	fqv2seq "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/sequences"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/nonce_manager"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ccipseqs "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"

	mcmslib "github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
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
	// SkipNonceManagerUpdates skips auto-detection of v1.5 contracts and NonceManager previous ramps updates, default=false
	SkipNonceManagerUpdates bool
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

	return UpdateLanesLogic(e, c.MCMSConfig, c.SkipNonceManagerUpdates, configs)
}

// UpdateLanesLogic configures CCIP lanes by updating OnRamp destinations, OffRamp sources,
// Router ramps, FeeQuoter dest chain configs, and FeeQuoter gas prices across all specified chains.
// On chains where a v2 FeeQuoter is deployed alongside the active v1.6 FeeQuoter, both are updated.
// Already-configured destinations are skipped to ensure idempotency. Configs provided can be unidirectional
// TODO: UpdateBidirectionalLanesChangesetConfigs name is misleading, it also accepts unidirectional lane updates
func UpdateLanesLogic(e cldf.Environment, mcmsConfig *proposalutils.TimelockConfig, skipNonceManagerUpdates bool, configs UpdateBidirectionalLanesChangesetConfigs) (cldf.ChangesetOutput, error) {
	var loadOpts []stateview.LoadOption
	if !skipNonceManagerUpdates {
		loadOpts = append(loadOpts, stateview.WithLoadLegacyContracts(true))
	}
	state, err := stateview.LoadOnchainState(e, loadOpts...)
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
	feeQuoterVersionsByChain, v2FQAddresses, err := resolveFeeQuoterTargets(ds, &feeQuoterDestsInput, &feeQuoterPricesInput)
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
			// Don't skip — still update the active v1 FeeQuoter below.
		}
		filtered, err := FilterOutExistingDestChainConfigs(e, update.Address, chainSel, update.CallInput)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if len(filtered) > 0 {
			update.CallInput = filtered
			v1FeeQuoterDestsUpdates[chainSel] = update
		}
	}
	for chainSel, update := range feeQuoterPricesInput.UpdatesByChain {
		version, ok := feeQuoterVersionsByChain[chainSel]
		if ok && version.Major() >= 2 {
			v2FeeQuoterChains[chainSel] = struct{}{}
			// Don't skip — still update the active v1 FeeQuoter below.
		}
		v1FeeQuoterPriceUpdates[chainSel] = update
	}

	// Build NonceManager updates by auto-detecting v1.5 contracts
	var nonceManagerInput ccipseqs.NonceManagerUpdatesSequenceInput
	if !skipNonceManagerUpdates {
		nonceManagerInput, err = buildNonceManagerUpdatesFromV15Contracts(e, state, configs, mcmsConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build nonce manager updates: %w", err)
		}
		if len(nonceManagerInput.UpdatesByChain) > 0 {
			e.Logger.Infow("Auto-detected v1.5 contracts, adding NonceManager previous ramps updates",
				"chainsWithUpdates", len(nonceManagerInput.UpdatesByChain))
		}
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
		NonceManagerUpdatesSequenceInput:                  nonceManagerInput,
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

	// Execute v2 FeeQuoter update sequences on chains that have a v2 FQ deployed
	v2FQChainSels := make([]uint64, 0, len(v2FeeQuoterChains))
	for chainSel := range v2FeeQuoterChains {
		v2FQChainSels = append(v2FQChainSels, chainSel)
	}
	slices.Sort(v2FQChainSels)

	var v2BatchOps []mcmstypes.BatchOperation
	for _, chainSel := range v2FQChainSels {
		fqUpdate := fqv2seq.FeeQuoterUpdate{
			ChainSelector:     chainSel,
			ExistingAddresses: ds.Addresses().Filter(datastore.AddressRefByChainSelector(chainSel)),
		}
		if dests, ok := feeQuoterDestsInput.UpdatesByChain[chainSel]; ok {
			destCfgs, err := ConvertV16FeeQuoterDestUpdatesToV2(dests.CallInput)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert v1.6 fee quoter destination updates for chain %d: %w", chainSel, err)
			}
			destCfgs, err = FilterOutExistingDestChainConfigs(e, v2FQAddresses[chainSel], chainSel, destCfgs)
			if err != nil {
				return cldf.ChangesetOutput{}, err
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

// destChainConfigType constrains the types accepted by FilterOutExistingDestChainConfigs.
type destChainConfigType interface {
	fee_quoter.FeeQuoterDestChainConfigArgs | fqv2ops.DestChainConfigArgs
}

// FilterOutExistingDestChainConfigs removes destination chain configs where the destination
// is already enabled on-chain. It automatically selects the correct FeeQuoter binding
// based on the concrete config type.
func FilterOutExistingDestChainConfigs[T destChainConfigType](
	e cldf.Environment,
	fqAddr common.Address,
	chainSel uint64,
	destCfgs []T,
) ([]T, error) {
	if len(destCfgs) == 0 {
		return destCfgs, nil
	}

	var isDestEnabled func(T) (uint64, bool, error)

	switch any(destCfgs[0]).(type) {
	case fee_quoter.FeeQuoterDestChainConfigArgs:
		fq, err := fee_quoter.NewFeeQuoter(fqAddr, e.BlockChains.EVMChains()[chainSel].Client)
		if err != nil {
			return nil, fmt.Errorf("failed to bind FeeQuoter on chain %d: %w", chainSel, err)
		}
		isDestEnabled = func(cfg T) (uint64, bool, error) {
			destSel := any(cfg).(fee_quoter.FeeQuoterDestChainConfigArgs).DestChainSelector
			onChain, err := fq.GetDestChainConfig(&bind.CallOpts{Context: e.GetContext()}, destSel)
			if err != nil {
				return destSel, false, err
			}
			return destSel, onChain.IsEnabled, nil
		}
	case fqv2ops.DestChainConfigArgs:
		fq, err := fqv2ops.NewFeeQuoterContract(fqAddr, e.BlockChains.EVMChains()[chainSel].Client)
		if err != nil {
			return nil, fmt.Errorf("failed to bind v2 FeeQuoter on chain %d: %w", chainSel, err)
		}
		isDestEnabled = func(cfg T) (uint64, bool, error) {
			destSel := any(cfg).(fqv2ops.DestChainConfigArgs).DestChainSelector
			onChain, err := fq.GetDestChainConfig(&bind.CallOpts{Context: e.GetContext()}, destSel)
			if err != nil {
				return destSel, false, err
			}
			return destSel, onChain.IsEnabled, nil
		}
	}

	filtered := make([]T, 0, len(destCfgs))
	for _, destCfg := range destCfgs {
		destSel, enabled, err := isDestEnabled(destCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to query existing dest chain config on chain %d for dest %d: %w",
				chainSel, destSel, err)
		}
		if enabled {
			e.Logger.Infow("skipping dest chain config already present on FeeQuoter",
				"sourceChain", chainSel,
				"destChain", destSel,
			)
			continue
		}
		filtered = append(filtered, destCfg)
	}

	return filtered, nil
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
) (map[uint64]semver.Version, map[uint64]common.Address, error) {
	versionsByChain := make(map[uint64]semver.Version)
	v2Addresses := make(map[uint64]common.Address)

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

		if version.Major() >= 2 {
			// Store v2 address separately; keep the active v1 address
			// (from on-chain state) in destsInput/pricesInput so the
			// v1 path still updates the active FeeQuoter.
			v2Addresses[chainSel] = addr
		} else {
			if update, ok := destsInput.UpdatesByChain[chainSel]; ok {
				update.Address = addr
				destsInput.UpdatesByChain[chainSel] = update
			}
			if update, ok := pricesInput.UpdatesByChain[chainSel]; ok {
				update.Address = addr
				pricesInput.UpdatesByChain[chainSel] = update
			}
		}
		return nil
	}

	for chainSel := range destsInput.UpdatesByChain {
		if err := resolve(chainSel); err != nil {
			return nil, nil, err
		}
	}
	for chainSel := range pricesInput.UpdatesByChain {
		if err := resolve(chainSel); err != nil {
			return nil, nil, err
		}
	}

	return versionsByChain, v2Addresses, nil
}

func resolveUpdateLanesFeeQuoterAddressAndVersion(
	addresses []datastore.AddressRef,
	chainSel uint64,
) (common.Address, semver.Version, error) {
	return shared.ResolveFeeQuoterAddressAndVersion(addresses, chainSel)
}

// buildNonceManagerUpdatesFromV15Contracts auto-detects v1.5 OnRamp/OffRamp contracts for the lanes
// being updated and builds NonceManager previous ramps updates to preserve nonce continuity
func buildNonceManagerUpdatesFromV15Contracts(
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	configs UpdateBidirectionalLanesChangesetConfigs,
	mcmsConfig *proposalutils.TimelockConfig,
) (ccipseqs.NonceManagerUpdatesSequenceInput, error) {
	updates := make(map[uint64]ccipseqs.NonceManagerUpdateInput)

	// Collect all unique directed pairs (both directions) from the lane configs.
	pairs := make(map[uint64]map[uint64]struct{})
	addPair := func(src, dest uint64) {
		if pairs[src] == nil {
			pairs[src] = make(map[uint64]struct{})
		}
		pairs[src][dest] = struct{}{}
	}

	for srcChain, dests := range configs.UpdateOnRampDestsConfig.UpdatesByChain {
		for destChain := range dests {
			addPair(srcChain, destChain)
			addPair(destChain, srcChain) // reverse direction
		}
	}
	for destChain, sources := range configs.UpdateOffRampSourcesConfig.UpdatesByChain {
		for srcChain := range sources {
			addPair(srcChain, destChain)
			addPair(destChain, srcChain) // reverse direction
		}
	}

	// Check each directed pair for v1.5 contracts
	for sourceChain, destChains := range pairs {
		for destChain := range destChains {
			if err := maybeAddPreviousRampUpdate(e, state, updates, sourceChain, destChain, mcmsConfig); err != nil {
				return ccipseqs.NonceManagerUpdatesSequenceInput{}, err
			}
		}
	}

	return ccipseqs.NonceManagerUpdatesSequenceInput{UpdatesByChain: updates}, nil
}

func maybeAddPreviousRampUpdate(
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	updates map[uint64]ccipseqs.NonceManagerUpdateInput,
	sourceChain, destChain uint64,
	mcmsConfig *proposalutils.TimelockConfig,
) error {
	chainState := state.Chains[sourceChain]
	if chainState.NonceManager == nil {
		return nil
	}

	// Check if v1.5 OnRamp exists
	var hasV15OnRamp bool
	if chainState.EVM2EVMOnRamp != nil {
		if onRamp := chainState.EVM2EVMOnRamp[destChain]; onRamp != nil && onRamp.Address() != (common.Address{}) {
			hasV15OnRamp = true
		}
	}

	// Check if v1.5 OffRamp exists
	var hasV15OffRamp bool
	if chainState.EVM2EVMOffRamp != nil {
		if offRamp := chainState.EVM2EVMOffRamp[destChain]; offRamp != nil && offRamp.Address() != (common.Address{}) {
			hasV15OffRamp = true
		}
	}

	// Require both v1.5 OnRamp and OffRamp to exist
	if !hasV15OnRamp || !hasV15OffRamp {
		return nil
	}

	wantOnRamp := chainState.EVM2EVMOnRamp[destChain].Address()
	wantOffRamp := chainState.EVM2EVMOffRamp[destChain].Address()

	// Check if already configured
	prevRamps, err := chainState.NonceManager.GetPreviousRamps(&bind.CallOpts{Context: e.GetContext()}, destChain)
	if err != nil {
		return fmt.Errorf("failed to get previous ramps for chain %d -> %d: %w", sourceChain, destChain, err)
	}

	// Skip only if the CORRECT addresses are already configured
	if prevRamps.PrevOnRamp == wantOnRamp && prevRamps.PrevOffRamp == wantOffRamp {
		return nil
	}

	if mcmsConfig != nil {
		if chainState.Timelock == nil {
			return fmt.Errorf("timelock not deployed on chain %d", sourceChain)
		}
		evmChain, ok := e.BlockChains.EVMChains()[sourceChain]
		if !ok {
			return fmt.Errorf("chain %d not found in environment", sourceChain)
		}
		if err := commoncs.ValidateOwnership(
			e.GetContext(),
			true,
			evmChain.DeployerKey.From,
			chainState.Timelock.Address(),
			chainState.NonceManager,
		); err != nil {
			return fmt.Errorf("NonceManager ownership validation failed on chain %d: %w", sourceChain, err)
		}
	}

	prevRampArgs := nonce_manager.NonceManagerPreviousRampsArgs{
		RemoteChainSelector:   destChain,
		OverrideExistingRamps: true, // Always override to support multiple lanes with same remote chain
		PrevRamps: nonce_manager.NonceManagerPreviousRamps{
			PrevOnRamp:  wantOnRamp,
			PrevOffRamp: wantOffRamp,
		},
	}

	e.Logger.Infow("Registering v1.5 ramp addresses in NonceManager",
		"sourceChain", sourceChain,
		"destChain", destChain,
		"v15OnRamp", wantOnRamp.Hex(),
		"v15OffRamp", wantOffRamp.Hex())

	update := updates[sourceChain]
	if update.PreviousRampsArgs == nil {
		update.PreviousRampsArgs = &opsutil.EVMCallInput[[]nonce_manager.NonceManagerPreviousRampsArgs]{
			Address:       chainState.NonceManager.Address(),
			ChainSelector: sourceChain,
			CallInput:     make([]nonce_manager.NonceManagerPreviousRampsArgs, 0),
			NoSend:        mcmsConfig != nil,
		}
	}
	update.PreviousRampsArgs.CallInput = append(update.PreviousRampsArgs.CallInput, prevRampArgs)
	updates[sourceChain] = update

	return nil
}
