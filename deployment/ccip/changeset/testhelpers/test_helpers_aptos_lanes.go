package testhelpers

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	chainsel "github.com/smartcontractkit/chain-selectors"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	_ "github.com/smartcontractkit/chainlink-aptos/deployment/ccip/adapters"              // register Aptos deploy/lane/curse/mcms adapters
	_ "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v1_6_0/sequences" // register EVM deploy/lane adapters
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	cs_ccip "github.com/smartcontractkit/chainlink-ccip/deployment/utils/changesets"
	ccipmcms "github.com/smartcontractkit/chainlink-ccip/deployment/utils/mcms"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

// addAptosMixedLane wires an Aptos↔peer lane via ConnectChains, then reconciles any
// EVM router mapping the mixed-family adapters may have skipped.
func addAptosMixedLane(
	t *testing.T,
	e *DeployedEnv,
	state stateview.CCIPOnChainState,
	from, to uint64,
	fromFamily, toFamily string,
	isTestRouter bool,
	gasPrices map[uint64]*big.Int,
	tokenPrices map[string]*big.Int,
	fqCfg fee_quoter.FeeQuoterDestChainConfig,
) error {
	changesets := addLaneAptosChangesets(t, from, to, isTestRouter, gasPrices, tokenPrices, fqCfg)
	var err error
	e.Env, _, err = commoncs.ApplyChangesets(t, e.Env, changesets)
	if err != nil {
		return err
	}
	return ensureEVMRouterLaneConfig(t, e, state, from, to, fromFamily, toFamily, isTestRouter)
}

// addLaneAptosChangesets connects an Aptos chain with its (EVM or Aptos) peer via the
// generic lanes.ConnectChains changeset, which configures both legs of the lane using the
// adapters registered by the chainlink-aptos and chainlink-ccip EVM adapter packages.
// Gas prices are keyed by chain selector (a missing entry falls back to the adapter default).
// Token prices are attached to the source chain definition; the Aptos source leg uses them
// when seeding FeeQuoter prices.
func addLaneAptosChangesets(
	t *testing.T,
	srcChainSelector, destChainSelector uint64,
	isTestRouter bool,
	gasPrices map[uint64]*big.Int,
	tokenPrices map[string]*big.Int,
	fqCfg fee_quoter.FeeQuoterDestChainConfig,
) []commoncs.ConfiguredChangeSet {
	srcFamily, err := chainsel.GetSelectorFamily(srcChainSelector)
	require.NoError(t, err)
	destFamily, err := chainsel.GetSelectorFamily(destChainSelector)
	require.NoError(t, err)

	if srcFamily != chainsel.FamilyAptos && destFamily != chainsel.FamilyAptos {
		t.Fatalf("At least one of the provided source/destination chains has to be Aptos. srcFamily: %v destFamily: %v", srcFamily, destFamily)
	}

	destFQOverride := feeQuoterDestChainConfigOverride(fqCfg)

	makeDefinition := func(selector uint64) lanes.ChainDefinition {
		definition := lanes.ChainDefinition{
			Selector:               selector,
			GasPrice:               gasPrices[selector],
			RMNVerificationEnabled: false,
			AllowListEnabled:       false,
		}
		if selector == srcChainSelector {
			definition.TokenPrices = tokenPrices
		}
		if selector == destChainSelector {
			definition.FeeQuoterDestChainConfigOverrides = &destFQOverride
		}
		return definition
	}

	validUntil, err := mcmsValidUntil(time.Now().Add(24 * time.Hour))
	require.NoError(t, err)

	return []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			lanes.ConnectChains(lanes.GetLaneAdapterRegistry(), cs_ccip.GetRegistry()),
			lanes.ConnectChainsConfig{
				MCMS: ccipmcms.Input{
					ValidUntil:     validUntil,
					TimelockDelay:  mcmstypes.NewDuration(time.Second),
					TimelockAction: mcmstypes.TimelockActionSchedule,
				},
				Lanes: []lanes.LaneConfig{
					{
						Version:    semver.MustParse("1.6.0"),
						ChainA:     makeDefinition(srcChainSelector),
						ChainB:     makeDefinition(destChainSelector),
						TestRouter: isTestRouter,
					},
				},
			},
		),
	}
}

func feeQuoterDestChainConfigOverride(cfg fee_quoter.FeeQuoterDestChainConfig) lanes.FeeQuoterDestChainConfigOverride {
	lanesCfg := evmBindingFQToLanesFQ(cfg)
	return func(dst *lanes.FeeQuoterDestChainConfig) {
		*dst = lanesCfg
	}
}

func evmBindingFQToLanesFQ(cfg fee_quoter.FeeQuoterDestChainConfig) lanes.FeeQuoterDestChainConfig {
	return lanes.FeeQuoterDestChainConfig{
		OverrideExistingConfig:      true,
		IsEnabled:                   cfg.IsEnabled,
		MaxDataBytes:                cfg.MaxDataBytes,
		MaxPerMsgGasLimit:           cfg.MaxPerMsgGasLimit,
		DestGasOverhead:             cfg.DestGasOverhead,
		DestGasPerPayloadByteBase:   cfg.DestGasPerPayloadByteBase,
		ChainFamilySelector:         binary.BigEndian.Uint32(cfg.ChainFamilySelector[:]),
		DefaultTokenFeeUSDCents:     cfg.DefaultTokenFeeUSDCents,
		DefaultTokenDestGasOverhead: cfg.DefaultTokenDestGasOverhead,
		DefaultTxGasLimit:           cfg.DefaultTxGasLimit,
		// #nosec G115 - test FeeQuoter configs use small network-fee values (e.g. 10 cents).
		NetworkFeeUSDCents: uint16(cfg.NetworkFeeUSDCents),
		V1Params: &lanes.FeeQuoterV1Params{
			MaxNumberOfTokensPerMsg:           cfg.MaxNumberOfTokensPerMsg,
			DestGasPerPayloadByteHigh:         cfg.DestGasPerPayloadByteHigh,
			DestGasPerPayloadByteThreshold:    cfg.DestGasPerPayloadByteThreshold,
			DestDataAvailabilityOverheadGas:   cfg.DestDataAvailabilityOverheadGas,
			DestGasPerDataAvailabilityByte:    cfg.DestGasPerDataAvailabilityByte,
			DestDataAvailabilityMultiplierBps: cfg.DestDataAvailabilityMultiplierBps,
			EnforceOutOfOrder:                 cfg.EnforceOutOfOrder,
			GasMultiplierWeiPerEth:            cfg.GasMultiplierWeiPerEth,
			GasPriceStalenessThreshold:        cfg.GasPriceStalenessThreshold,
		},
	}
}

// ensureEVMRouterLaneConfig verifies EVM router on-ramp/off-ramp mappings for mixed
// Aptos lanes. ConnectChains should configure these via adapters; this is a narrow
// fallback that applies router-only changesets when the on-chain mapping is missing.
func ensureEVMRouterLaneConfig(
	t *testing.T,
	e *DeployedEnv,
	state stateview.CCIPOnChainState,
	from, to uint64,
	fromFamily, toFamily string,
	isTestRouter bool,
) error {
	if fromFamily == chainsel.FamilyEVM && toFamily == chainsel.FamilyAptos {
		return ensureEVMLaneLeg(
			t, e,
			func() (bool, error) { return evmSourceRouterWired(e, state, from, to, isTestRouter) },
			func() []commoncs.ConfiguredChangeSet { return addEVMSourceRouterRampChangesets(from, to, isTestRouter) },
			fmt.Sprintf("EVM source router on chain %d missing on-ramp for destination %d", from, to),
		)
	}

	if fromFamily == chainsel.FamilyAptos && toFamily == chainsel.FamilyEVM {
		return ensureEVMLaneLeg(
			t, e,
			func() (bool, error) { return evmDestRouterWired(e, state, to, from, isTestRouter) },
			func() []commoncs.ConfiguredChangeSet { return addEVMDestRouterRampChangesets(to, from, isTestRouter) },
			fmt.Sprintf("EVM destination router on chain %d missing off-ramp for source %d", to, from),
		)
	}

	return nil
}

func ensureEVMLaneLeg(
	t *testing.T,
	e *DeployedEnv,
	isWired func() (bool, error),
	apply func() []commoncs.ConfiguredChangeSet,
	verifyErrMsg string,
) error {
	wired, err := isWired()
	if err != nil {
		return err
	}
	if wired {
		return nil
	}

	changesets := apply()
	if len(changesets) > 0 {
		e.Env, _, err = commoncs.ApplyChangesets(t, e.Env, changesets)
		if err != nil {
			return err
		}
	}

	wired, err = isWired()
	if err != nil {
		return err
	}
	if !wired {
		return fmt.Errorf("%s", verifyErrMsg)
	}
	return nil
}

func evmSourceRouterWired(e *DeployedEnv, state stateview.CCIPOnChainState, evmSource, remoteDest uint64, isTestRouter bool) (bool, error) {
	chainState := state.MustGetEVMChainState(evmSource)
	routerContract := chainState.Router
	if isTestRouter {
		routerContract = chainState.TestRouter
	}
	if routerContract == nil {
		return false, fmt.Errorf("router not found for EVM source chain %d", evmSource)
	}
	if chainState.OnRamp == nil {
		return false, fmt.Errorf("on-ramp not found for EVM source chain %d", evmSource)
	}

	onRamp, err := routerContract.GetOnRamp(&bind.CallOpts{Context: e.Env.GetContext()}, remoteDest)
	if err != nil {
		return false, fmt.Errorf("get on-ramp for EVM source chain %d and destination %d: %w", evmSource, remoteDest, err)
	}

	return onRamp == chainState.OnRamp.Address(), nil
}

func evmDestRouterWired(e *DeployedEnv, state stateview.CCIPOnChainState, evmDest, remoteSource uint64, isTestRouter bool) (bool, error) {
	chainState := state.MustGetEVMChainState(evmDest)
	routerContract := chainState.Router
	if isTestRouter {
		routerContract = chainState.TestRouter
	}
	if routerContract == nil {
		return false, fmt.Errorf("router not found for EVM destination chain %d", evmDest)
	}
	if chainState.OffRamp == nil {
		return false, fmt.Errorf("off-ramp not found for EVM destination chain %d", evmDest)
	}

	isOffRamp, err := routerContract.IsOffRamp(&bind.CallOpts{Context: e.Env.GetContext()}, remoteSource, chainState.OffRamp.Address())
	if err != nil {
		return false, fmt.Errorf("check off-ramp for EVM destination chain %d and source %d: %w", evmDest, remoteSource, err)
	}

	return isOffRamp, nil
}

func addEVMSourceRouterRampChangesets(from, to uint64, isTestRouter bool) []commoncs.ConfiguredChangeSet {
	return []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateRouterRampsChangeset),
			v1_6.UpdateRouterRampsConfig{
				TestRouter: isTestRouter,
				UpdatesByChain: map[uint64]v1_6.RouterUpdates{
					from: {
						OnRampUpdates: map[uint64]bool{
							to: true,
						},
					},
				},
			},
		),
	}
}

func addEVMDestRouterRampChangesets(to, from uint64, isTestRouter bool) []commoncs.ConfiguredChangeSet {
	return []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateRouterRampsChangeset),
			v1_6.UpdateRouterRampsConfig{
				TestRouter: isTestRouter,
				UpdatesByChain: map[uint64]v1_6.RouterUpdates{
					to: {
						OffRampUpdates: map[uint64]bool{
							from: true,
						},
					},
				},
			},
		),
	}
}
