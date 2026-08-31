package testhelpers

import (
	"math/big"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"
	_ "github.com/smartcontractkit/chainlink-ccip/chains/solana/deployment/v1_6_0/sequences" // register Solana deploy/lane/mcms adapters
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	cs_ccip "github.com/smartcontractkit/chainlink-ccip/deployment/utils/changesets"
	ccipmcms "github.com/smartcontractkit/chainlink-ccip/deployment/utils/mcms"

	suilanes "github.com/smartcontractkit/chainlink-sui/deployment/lanes" // registers SuiAdapter and provides WithConnectChainsEnvironment

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

// addSuiSolanaMixedLane wires a Sui<->Solana lane via ConnectChains. Unlike the Aptos mixed
// lane, there is no EVM leg to reconcile, so no ensureEVMRouterLaneConfig fallback is needed.
// ConnectChains is bidirectional: it configures both Sui->Solana and Solana->Sui legs using the
// SuiAdapter and SolanaAdapter registered by the chainlink-sui lanes package and the Solana
// v1_6_0 sequences package (blank-imported above).
//
// The SuiAdapter resolves Sui package IDs from the active CLDF environment rather than the
// datastore, so ConnectChains MUST run inside suilanes.WithConnectChainsEnvironment (see
// chainlink-sui/deployment/lanes/env.go). Other chain families are unaffected by that scope.
func addSuiSolanaMixedLane(
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
	changesets := addSuiSolanaLaneChangesets(t, from, to, isTestRouter, gasPrices, tokenPrices, fqCfg)
	// SuiAdapter address getters read chain metadata from the env scope set below; run the
	// ConnectChains changeset application inside it.
	return suilanes.WithConnectChainsEnvironment(e.Env, func() error {
		var err error
		e.Env, _, err = commoncs.ApplyChangesets(t, e.Env, changesets)
		return err
	})
}

// addSuiSolanaLaneChangesets connects a Sui chain with its Solana peer via the generic
// lanes.ConnectChains changeset, which configures both legs of the lane using the adapters
// registered by the chainlink-sui lanes package and the chainlink-ccip Solana v1_6_0 sequences
// package. Gas prices are keyed by chain selector (a missing entry falls back to the adapter
// default). Token prices are attached to the source chain definition; the source leg uses them
// when seeding FeeQuoter prices. The FeeQuoterDestChainConfig override is attached to the
// destination definition, mirroring addLaneAptosChangesets.
func addSuiSolanaLaneChangesets(
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

	if (srcFamily != chainsel.FamilySui && srcFamily != chainsel.FamilySolana) ||
		(destFamily != chainsel.FamilySui && destFamily != chainsel.FamilySolana) {
		t.Fatalf("addSuiSolanaLaneChangesets requires a Sui<->Solana lane. srcFamily: %v destFamily: %v", srcFamily, destFamily)
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
