package ccip

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms/types"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonSolana "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

const (
	Evm1 = uint64(0)
	Evm2 = uint64(1)
	Sol1 = uint64(2)
)

type curseAssertion struct {
	chainID     uint64
	subject     uint64
	globalCurse bool
	cursed      bool
}

type CurseTestCase struct {
	name                string
	curseActionsBuilder func(mapIDToSelectorFunc) []v1_6.CurseAction
	curseAssertions     []curseAssertion
}

type mapIDToSelectorFunc func(uint64) uint64

var testCases = []CurseTestCase{
	{
		name: "lane",
		curseActionsBuilder: func(mapIDToSelector mapIDToSelectorFunc) []v1_6.CurseAction {
			return []v1_6.CurseAction{v1_6.CurseLaneBidirectionally(mapIDToSelector(Evm1), mapIDToSelector(Evm2))}
		},
		curseAssertions: []curseAssertion{
			{chainID: Evm1, subject: Evm2, cursed: true},
			{chainID: Evm1, subject: Sol1, cursed: false},
			{chainID: Evm2, subject: Evm1, cursed: true},
			{chainID: Evm2, subject: Sol1, cursed: false},
			{chainID: Sol1, subject: Evm1, cursed: false},
			{chainID: Sol1, subject: Evm2, cursed: false},
		},
	},
	{
		name: "solana lane",
		curseActionsBuilder: func(mapIDToSelector mapIDToSelectorFunc) []v1_6.CurseAction {
			return []v1_6.CurseAction{v1_6.CurseLaneBidirectionally(mapIDToSelector(Evm1), mapIDToSelector(Sol1))}
		},
		curseAssertions: []curseAssertion{
			{chainID: Evm1, subject: Evm2, cursed: false},
			{chainID: Evm1, subject: Sol1, cursed: true},
			{chainID: Evm2, subject: Evm1, cursed: false},
			{chainID: Evm2, subject: Sol1, cursed: false},
			{chainID: Sol1, subject: Evm1, cursed: true},
			{chainID: Sol1, subject: Evm2, cursed: false},
		},
	},
	{
		name: "lane duplicate",
		curseActionsBuilder: func(mapIDToSelector mapIDToSelectorFunc) []v1_6.CurseAction {
			return []v1_6.CurseAction{
				v1_6.CurseLaneBidirectionally(mapIDToSelector(Evm1), mapIDToSelector(Sol1)),
				v1_6.CurseLaneBidirectionally(mapIDToSelector(Evm1), mapIDToSelector(Sol1))}
		},
		curseAssertions: []curseAssertion{
			{chainID: Evm1, subject: Evm2, cursed: false},
			{chainID: Evm1, subject: Sol1, cursed: true},
			{chainID: Evm2, subject: Evm1, cursed: false},
			{chainID: Evm2, subject: Sol1, cursed: false},
			{chainID: Sol1, subject: Evm1, cursed: true},
			{chainID: Sol1, subject: Evm2, cursed: false},
		},
	},
	{
		name: "chain",
		curseActionsBuilder: func(mapIDToSelector mapIDToSelectorFunc) []v1_6.CurseAction {
			return []v1_6.CurseAction{v1_6.CurseChain(mapIDToSelector(Evm1))}
		},
		curseAssertions: []curseAssertion{
			{chainID: Evm1, globalCurse: true, cursed: true},
			{chainID: Evm2, subject: Evm1, cursed: true},
			{chainID: Evm2, subject: Sol1, cursed: false},
			{chainID: Sol1, subject: Evm1, cursed: true},
			{chainID: Sol1, subject: Evm2, cursed: false},
		},
	},
	{
		name: "solana chain",
		curseActionsBuilder: func(mapIDToSelector mapIDToSelectorFunc) []v1_6.CurseAction {
			return []v1_6.CurseAction{v1_6.CurseChain(mapIDToSelector(Sol1))}
		},
		curseAssertions: []curseAssertion{
			{chainID: Sol1, globalCurse: true, cursed: true},
			{chainID: Evm1, subject: Evm2, cursed: false},
			{chainID: Evm1, subject: Sol1, cursed: true},
			{chainID: Evm2, subject: Evm1, cursed: false},
			{chainID: Evm2, subject: Sol1, cursed: true},
			{chainID: Sol1, subject: Evm1, cursed: true},
			{chainID: Sol1, subject: Evm2, cursed: true},
		},
	},
	{
		name: "chain duplicate",
		curseActionsBuilder: func(mapIDToSelector mapIDToSelectorFunc) []v1_6.CurseAction {
			return []v1_6.CurseAction{v1_6.CurseChain(mapIDToSelector(Evm1)), v1_6.CurseChain(mapIDToSelector(Evm1))}
		},
		curseAssertions: []curseAssertion{
			{chainID: Evm1, globalCurse: true, cursed: true},
			{chainID: Evm2, subject: Evm1, cursed: true},
			{chainID: Evm2, subject: Sol1, cursed: false},
			{chainID: Sol1, subject: Evm1, cursed: true},
			{chainID: Sol1, subject: Evm2, cursed: false},
		},
	},
	{
		name: "solana chain duplicate",
		curseActionsBuilder: func(mapIDToSelector mapIDToSelectorFunc) []v1_6.CurseAction {
			return []v1_6.CurseAction{v1_6.CurseChain(mapIDToSelector(Sol1)), v1_6.CurseChain(mapIDToSelector(Sol1))}
		},
		curseAssertions: []curseAssertion{
			{chainID: Sol1, globalCurse: true, cursed: true},
			{chainID: Evm1, subject: Evm2, cursed: false},
			{chainID: Evm1, subject: Sol1, cursed: true},
			{chainID: Evm2, subject: Evm1, cursed: false},
			{chainID: Evm2, subject: Sol1, cursed: true},
			{chainID: Sol1, subject: Evm1, cursed: true},
			{chainID: Sol1, subject: Evm2, cursed: true},
		},
	},
	{
		name: "chain and lanes",
		curseActionsBuilder: func(mapIDToSelector mapIDToSelectorFunc) []v1_6.CurseAction {
			return []v1_6.CurseAction{v1_6.CurseChain(mapIDToSelector(Evm1)), v1_6.CurseLaneBidirectionally(mapIDToSelector(Evm2), mapIDToSelector(Sol1))}
		},
		curseAssertions: []curseAssertion{
			{chainID: Evm1, globalCurse: true, cursed: true},
			{chainID: Evm2, subject: Evm1, cursed: true},
			{chainID: Evm2, subject: Sol1, cursed: true},
			{chainID: Sol1, subject: Evm1, cursed: true},
			{chainID: Sol1, subject: Evm2, cursed: true},
		},
	},
	{
		name: "all chain",
		curseActionsBuilder: func(mapIDToSelector mapIDToSelectorFunc) []v1_6.CurseAction {
			return []v1_6.CurseAction{v1_6.CurseGloballyAllChains()}
		},
		curseAssertions: []curseAssertion{
			{chainID: Evm1, globalCurse: true, cursed: true},
			{chainID: Evm2, globalCurse: true, cursed: true},
			{chainID: Sol1, globalCurse: true, cursed: true},
		},
	},
}

func TestRMNCurse(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name+"_NO_MCMS", func(t *testing.T) {
			runRmnCurseTest(t, tc)
		})
		t.Run(tc.name+"_MCMS", func(t *testing.T) {
			runRmnCurseMCMSTest(t, tc, types.TimelockActionSchedule)
		})
	}
}

func TestRMNCurseBypass(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name+"_MCMS", func(t *testing.T) {
			runRmnCurseMCMSTest(t, tc, types.TimelockActionBypass)
		})
	}
}

func TestRMNCurseIdempotent(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name+"_CURSE_IDEMPOTENT_NO_MCMS", func(t *testing.T) {
			runRmnCurseIdempotentTest(t, tc)
		})
	}
}

func TestRMNUncurseIdempotent(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name+"_UNCURSE_IDEMPOTENT_NO_MCMS", func(t *testing.T) {
			runRmnUncurseIdempotentTest(t, tc)
		})
	}
}

func TestRMNUncurse(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name+"_UNCURSE", func(t *testing.T) {
			runRmnUncurseTest(t, tc)
		})
		t.Run(tc.name+"_UNCURSE_MCMS", func(t *testing.T) {
			runRmnUncurseMCMSTest(t, tc, types.TimelockActionSchedule)
		})
	}
}

func TestRMNUncurseBypass(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name+"_UNCURSE_MCMS", func(t *testing.T) {
			runRmnUncurseMCMSTest(t, tc, types.TimelockActionBypass)
		})
	}
}

func TestRMNCurseConfigValidate(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name+"_VALIDATE", func(t *testing.T) {
			runRmnCurseConfigValidateTest(t, tc)
		})
	}
}

func TestRMNCurseNoConnectedLanes(t *testing.T) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2), testhelpers.WithSolChains(1))

	mapIDToSelector := func(id uint64) uint64 {
		return v1_6.GetAllCursableChainsSelector(e.Env)[id]
	}

	verifyNoActiveCurseOnAllChains(t, &e)

	config := v1_6.RMNCurseConfig{
		CurseActions: []v1_6.CurseAction{
			v1_6.CurseChain(mapIDToSelector(Evm1)),
		},
		Reason:                   "test curse",
		IncludeNotConnectedLanes: false, // This will filter out non connected lanes
	}

	_, err := v1_6.RMNCurseChangeset(e.Env, config)
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, CurseTestCase{
		curseAssertions: []curseAssertion{
			{chainID: Evm1, globalCurse: true, cursed: true},
			{chainID: Evm1, subject: Evm2, cursed: true}, // 0 is globally cursed return true for everything
			{chainID: Evm1, subject: Sol1, cursed: true}, // 0 is globally cursed return true for everything
			{chainID: Evm2, subject: Evm1, cursed: false},
			{chainID: Evm2, subject: Sol1, cursed: false},
			{chainID: Sol1, subject: Evm1, cursed: false},
			{chainID: Sol1, subject: Evm2, cursed: false},
		},
	}, mapIDToSelector)
}

func TestRMNCurseOneConnectedLanes(t *testing.T) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2), testhelpers.WithSolChains(1))

	mapIDToSelector := func(id uint64) uint64 {
		return v1_6.GetAllCursableChainsSelector(e.Env)[id]
	}

	_, err := v1_6.UpdateOffRampSourcesChangeset(e.Env,
		v1_6.UpdateOffRampSourcesConfig{
			UpdatesByChain: map[uint64]map[uint64]v1_6.OffRampSourceUpdate{
				mapIDToSelector(Evm1): { // to
					mapIDToSelector(Evm2): { // from
						IsEnabled:                 true,
						TestRouter:                false,
						IsRMNVerificationDisabled: true,
					},
				},
			},
		},
	)
	require.NoError(t, err)

	verifyNoActiveCurseOnAllChains(t, &e)

	config := v1_6.RMNCurseConfig{
		CurseActions: []v1_6.CurseAction{
			v1_6.CurseChain(mapIDToSelector(Evm1)),
		},
		Reason:                   "test curse",
		IncludeNotConnectedLanes: false, // This will filter out non connected lanes
	}

	_, err = v1_6.RMNCurseChangeset(e.Env, config)
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, CurseTestCase{
		curseAssertions: []curseAssertion{
			{chainID: Evm1, globalCurse: true, cursed: true},
			{chainID: Evm1, subject: Evm2, cursed: true},
			{chainID: Evm1, subject: Sol1, cursed: true}, // 2 is not connected to 0 but 0 is globally cursed return true for everything
			{chainID: Evm2, subject: Evm1, cursed: true},
			{chainID: Evm2, subject: Sol1, cursed: false},
			{chainID: Sol1, subject: Evm1, cursed: false}, // 2 is not connected to 0
			{chainID: Sol1, subject: Evm2, cursed: false},
		},
	}, mapIDToSelector)
}

func TestRMNCurseOneConnectedLanesGlobalOnly(t *testing.T) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2), testhelpers.WithSolChains(1))

	mapIDToSelector := func(id uint64) uint64 {
		return v1_6.GetAllCursableChainsSelector(e.Env)[id]
	}

	_, err := v1_6.UpdateOffRampSourcesChangeset(e.Env,
		v1_6.UpdateOffRampSourcesConfig{
			UpdatesByChain: map[uint64]map[uint64]v1_6.OffRampSourceUpdate{
				mapIDToSelector(Evm1): { // to
					mapIDToSelector(Evm2): { // from
						IsEnabled:                 true,
						TestRouter:                false,
						IsRMNVerificationDisabled: true,
					},
				},
			},
		},
	)
	require.NoError(t, err)

	verifyNoActiveCurseOnAllChains(t, &e)

	config := v1_6.RMNCurseConfig{
		CurseActions: []v1_6.CurseAction{
			v1_6.CurseGloballyOnlyOnChain(mapIDToSelector(Evm1)),
		},
		Reason:                   "test curse",
		IncludeNotConnectedLanes: false, // This will filter out non connected lanes
	}

	_, err = v1_6.RMNCurseChangeset(e.Env, config)
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, CurseTestCase{
		curseAssertions: []curseAssertion{
			{chainID: Evm1, globalCurse: true, cursed: true},
			{chainID: Evm1, subject: Evm2, cursed: true},
			{chainID: Evm1, subject: Sol1, cursed: true}, // 2 is not connected to 0 but 0 is globally cursed return true for everything
			{chainID: Evm2, subject: Evm1, cursed: false},
			{chainID: Evm2, subject: Sol1, cursed: false},
			{chainID: Sol1, subject: Evm1, cursed: false},
			{chainID: Sol1, subject: Evm2, cursed: false},
		},
	}, mapIDToSelector)
}

func TestRMNCurseOneConnectedLanesLaneOnlyOnSource(t *testing.T) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2), testhelpers.WithSolChains(1))

	mapIDToSelector := func(id uint64) uint64 {
		return v1_6.GetAllCursableChainsSelector(e.Env)[id]
	}

	_, err := v1_6.UpdateOffRampSourcesChangeset(e.Env,
		v1_6.UpdateOffRampSourcesConfig{
			UpdatesByChain: map[uint64]map[uint64]v1_6.OffRampSourceUpdate{
				mapIDToSelector(Evm1): { // to
					mapIDToSelector(Evm2): { // from
						IsEnabled:                 true,
						TestRouter:                false,
						IsRMNVerificationDisabled: true,
					},
				},
			},
		},
	)
	require.NoError(t, err)

	verifyNoActiveCurseOnAllChains(t, &e)

	config := v1_6.RMNCurseConfig{
		CurseActions: []v1_6.CurseAction{
			v1_6.CurseLaneOnlyOnSource(mapIDToSelector(Evm1), mapIDToSelector(Evm2)),
		},
		Reason:                   "test curse",
		IncludeNotConnectedLanes: false, // This will filter out non connected lanes
	}

	_, err = v1_6.RMNCurseChangeset(e.Env, config)
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, CurseTestCase{
		curseAssertions: []curseAssertion{
			{chainID: Evm1, subject: Evm2, cursed: true},
			{chainID: Evm1, subject: Sol1, cursed: false},
			{chainID: Evm2, subject: Evm1, cursed: false},
			{chainID: Evm2, subject: Sol1, cursed: false},
			{chainID: Sol1, subject: Evm1, cursed: false},
			{chainID: Sol1, subject: Evm2, cursed: false},
		},
	}, mapIDToSelector)
}

func TestRMNCurseOneConnectedLanesSolana(t *testing.T) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2), testhelpers.WithSolChains(1))

	mapIDToSelector := func(id uint64) uint64 {
		return v1_6.GetAllCursableChainsSelector(e.Env)[id]
	}

	_, err := ccipChangesetSolana.AddRemoteChainToOffRamp(e.Env, ccipChangesetSolana.AddRemoteChainToOffRampConfig{
		ChainSelector: mapIDToSelector(Sol1),
		UpdatesByChain: map[uint64]*ccipChangesetSolana.OffRampConfig{
			mapIDToSelector(Evm1): {
				EnabledAsSource: true,
				IsUpdate:        false,
			},
		},
	})
	require.NoError(t, err)

	verifyNoActiveCurseOnAllChains(t, &e)

	config := v1_6.RMNCurseConfig{
		CurseActions: []v1_6.CurseAction{
			v1_6.CurseChain(mapIDToSelector(Evm1)),
		},
		Reason:                   "test curse",
		IncludeNotConnectedLanes: false, // This will filter out non connected lanes
	}

	_, err = v1_6.RMNCurseChangeset(e.Env, config)
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, CurseTestCase{
		curseAssertions: []curseAssertion{
			{chainID: Evm1, globalCurse: true, cursed: true},
			{chainID: Evm1, subject: Evm2, cursed: true},
			{chainID: Evm1, subject: Sol1, cursed: true},
			{chainID: Evm2, subject: Evm1, cursed: false},
			{chainID: Evm2, subject: Sol1, cursed: false},
			{chainID: Sol1, subject: Evm1, cursed: true},
			{chainID: Sol1, subject: Evm2, cursed: false},
		},
	}, mapIDToSelector)
}

func runRmnUncurseTest(t *testing.T, tc CurseTestCase) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2), testhelpers.WithSolChains(1))

	mapIDToSelector := func(id uint64) uint64 {
		return v1_6.GetAllCursableChainsSelector(e.Env)[id]
	}

	verifyNoActiveCurseOnAllChains(t, &e)

	config := v1_6.RMNCurseConfig{
		CurseActions:             tc.curseActionsBuilder(mapIDToSelector),
		Reason:                   "test curse",
		IncludeNotConnectedLanes: true,
	}

	_, err := v1_6.RMNCurseChangeset(e.Env, config)
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, tc, mapIDToSelector)

	_, err = v1_6.RMNUncurseChangeset(e.Env, config)
	require.NoError(t, err)

	verifyNoActiveCurseOnAllChains(t, &e)
}

func transferRMNContractToMCMS(t *testing.T, e *testhelpers.DeployedEnv, state stateview.CCIPOnChainState) {
	contractsByChain := make(map[uint64][]common.Address)
	rmnRemotePerChain := v1_6.BuildRMNRemotePerChain(e.Env, state)
	rmnRemoteAddressesByChain := make(map[uint64]common.Address)
	for chain, remote := range rmnRemotePerChain {
		if remote == nil {
			continue
		}
		rmnRemoteAddressesByChain[chain] = remote.Address()
	}
	for chainSelector, rmnRemoteAddress := range rmnRemoteAddressesByChain {
		contractsByChain[chainSelector] = []common.Address{rmnRemoteAddress}
	}

	contractsByChain[e.HomeChainSel] = append(contractsByChain[e.HomeChainSel], state.MustGetEVMChainState(e.HomeChainSel).RMNHome.Address())
	timelocksPerChain := deployergroup.BuildTimelockPerChain(e.Env, state)
	// This is required because RMN Contracts is initially owned by the deployer
	_, err := commonchangeset.Apply(t, e.Env, timelocksPerChain,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2),
			commonchangeset.TransferToMCMSWithTimelockConfig{
				ContractsByChain: contractsByChain,
				MCMSConfig: proposalutils.TimelockConfig{
					MinDelay: 0 * time.Second,
				},
			},
		),
	)
	require.NoError(t, err)

	chainSelectorSolana := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))
	for _, solChain := range chainSelectorSolana {
		_, _ = testhelpers.TransferOwnershipSolana(t, &e.Env, solChain, true,
			ccipChangesetSolana.CCIPContractsToTransfer{
				Router:    true,
				FeeQuoter: true,
				OffRamp:   true,
				RMNRemote: true,
			})
	}

	cfgAmounts := commonSolana.AmountsToTransfer{
		ProposeMCM:   100 * solana.LAMPORTS_PER_SOL,
		CancellerMCM: 350 * solana.LAMPORTS_PER_SOL,
		BypasserMCM:  75 * solana.LAMPORTS_PER_SOL,
		Timelock:     83 * solana.LAMPORTS_PER_SOL,
	}
	amountsPerChain := make(map[uint64]commonSolana.AmountsToTransfer)
	for _, chainSelector := range e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana)) {
		amountsPerChain[chainSelector] = cfgAmounts
	}
	config := commonSolana.FundMCMSignerConfig{
		AmountsPerChain: amountsPerChain,
	}

	changesetInstance := commonSolana.FundMCMSignersChangeset{}

	_, _, err = commonchangeset.ApplyChangesetsV2(t, e.Env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(changesetInstance, config),
	})
	require.NoError(t, err)
}

func runRmnUncurseMCMSTest(t *testing.T, tc CurseTestCase, action types.TimelockAction) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2), testhelpers.WithSolChains(1))

	mapIDToSelector := func(id uint64) uint64 {
		return v1_6.GetAllCursableChainsSelector(e.Env)[id]
	}

	config := v1_6.RMNCurseConfig{
		CurseActions: tc.curseActionsBuilder(mapIDToSelector),
		Reason:       "test curse",
		MCMS: &proposalutils.TimelockConfig{
			MinDelay:   1 * time.Second,
			MCMSAction: action,
		},
		IncludeNotConnectedLanes: true,
	}

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	verifyNoActiveCurseOnAllChains(t, &e)

	transferRMNContractToMCMS(t, &e, state)

	_, _, err = commonchangeset.ApplyChangesetsV2(t, e.Env,
		[]commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.RMNCurseChangeset),
				config,
			)},
	)
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, tc, mapIDToSelector)

	_, _, err = commonchangeset.ApplyChangesetsV2(t, e.Env,
		[]commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.RMNUncurseChangeset),
				config,
			)},
	)
	require.NoError(t, err)

	verifyNoActiveCurseOnAllChains(t, &e)
}

func runRmnCurseConfigValidateTest(t *testing.T, tc CurseTestCase) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2), testhelpers.WithSolChains(1))

	mapIDToSelector := func(id uint64) uint64 {
		return v1_6.GetAllCursableChainsSelector(e.Env)[id]
	}

	config := v1_6.RMNCurseConfig{
		CurseActions:             tc.curseActionsBuilder(mapIDToSelector),
		Reason:                   "test curse",
		IncludeNotConnectedLanes: true,
	}

	err := config.Validate(e.Env)
	require.NoError(t, err)
}

func runRmnCurseTest(t *testing.T, tc CurseTestCase) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2), testhelpers.WithSolChains(1))

	mapIDToSelector := func(id uint64) uint64 {
		return v1_6.GetAllCursableChainsSelector(e.Env)[id]
	}

	verifyNoActiveCurseOnAllChains(t, &e)

	config := v1_6.RMNCurseConfig{
		CurseActions:             tc.curseActionsBuilder(mapIDToSelector),
		Reason:                   "test curse",
		IncludeNotConnectedLanes: true,
	}

	_, err := v1_6.RMNCurseChangeset(e.Env, config)
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, tc, mapIDToSelector)
}

func runRmnCurseIdempotentTest(t *testing.T, tc CurseTestCase) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2), testhelpers.WithSolChains(1))

	mapIDToSelector := func(id uint64) uint64 {
		return v1_6.GetAllCursableChainsSelector(e.Env)[id]
	}

	verifyNoActiveCurseOnAllChains(t, &e)

	config := v1_6.RMNCurseConfig{
		CurseActions:             tc.curseActionsBuilder(mapIDToSelector),
		Reason:                   "test curse",
		IncludeNotConnectedLanes: true,
	}

	_, err := v1_6.RMNCurseChangeset(e.Env, config)
	require.NoError(t, err)

	_, err = v1_6.RMNCurseChangeset(e.Env, config)
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, tc, mapIDToSelector)
}

func runRmnUncurseIdempotentTest(t *testing.T, tc CurseTestCase) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2), testhelpers.WithSolChains(1))

	mapIDToSelector := func(id uint64) uint64 {
		return v1_6.GetAllCursableChainsSelector(e.Env)[id]
	}

	verifyNoActiveCurseOnAllChains(t, &e)

	config := v1_6.RMNCurseConfig{
		CurseActions:             tc.curseActionsBuilder(mapIDToSelector),
		Reason:                   "test curse",
		IncludeNotConnectedLanes: true,
	}

	_, err := v1_6.RMNCurseChangeset(e.Env, config)
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, tc, mapIDToSelector)

	_, err = v1_6.RMNUncurseChangeset(e.Env, config)
	require.NoError(t, err)

	_, err = v1_6.RMNUncurseChangeset(e.Env, config)
	require.NoError(t, err)

	verifyNoActiveCurseOnAllChains(t, &e)
}

func runRmnCurseMCMSTest(t *testing.T, tc CurseTestCase, action types.TimelockAction) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2), testhelpers.WithSolChains(1))

	mapIDToSelector := func(id uint64) uint64 {
		return v1_6.GetAllCursableChainsSelector(e.Env)[id]
	}

	config := v1_6.RMNCurseConfig{
		CurseActions: tc.curseActionsBuilder(mapIDToSelector),
		Reason:       "test curse",
		MCMS: &proposalutils.TimelockConfig{
			MinDelay:   1 * time.Second,
			MCMSAction: action,
		},
		IncludeNotConnectedLanes: true,
	}

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	verifyNoActiveCurseOnAllChains(t, &e)

	transferRMNContractToMCMS(t, &e, state)

	_, _, err = commonchangeset.ApplyChangesetsV2(t, e.Env,
		[]commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.RMNCurseChangeset),
				config,
			)},
	)
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, tc, mapIDToSelector)
}

func verifyTestCaseAssertions(t *testing.T, e *testhelpers.DeployedEnv, tc CurseTestCase, mapIDToSelector mapIDToSelectorFunc) {
	cursableChains, err := v1_6.GetCursableChains(e.Env)
	require.NoError(t, err)

	for _, assertion := range tc.curseAssertions {
		family, err := chain_selectors.GetSelectorFamily(mapIDToSelector(assertion.chainID))
		require.NoError(t, err)

		cursedSubject := globals.FamilyAwareSelectorToSubject(mapIDToSelector(assertion.subject), family)
		if assertion.globalCurse {
			cursedSubject = globals.GlobalCurseSubject()
		}

		isCursed, err := cursableChains[mapIDToSelector(assertion.chainID)].IsSubjectCursed(cursedSubject)
		require.NoError(t, err)
		require.Equal(t, assertion.cursed, isCursed, "chain %d subject %d", assertion.chainID, assertion.subject)
	}
}

func verifyNoActiveCurseOnAllChains(t *testing.T, e *testhelpers.DeployedEnv) {
	cursableChains, err := v1_6.GetCursableChains(e.Env)
	require.NoError(t, err)

	for chainSelector, chain := range cursableChains {
		for selector := range cursableChains {
			family, err := chain_selectors.GetSelectorFamily(chainSelector)
			require.NoError(t, err)

			isCursed, err := chain.IsSubjectCursed(globals.FamilyAwareSelectorToSubject(selector, family))
			require.NoError(t, err)
			require.False(t, isCursed, "chain %d subject %d", chainSelector, globals.FamilyAwareSelectorToSubject(selector, family))
		}
	}
}

type ForceOptionTestCase struct {
	name               string
	force              bool
	expectedOperations int
	expectProposal     bool
	applyChangeset     func(
		env *testhelpers.DeployedEnv,
		config v1_6.RMNCurseConfig,
		t *testing.T,
	) deployment.ChangesetOutput
}

var forceOptionTestCases = []ForceOptionTestCase{
	{
		name: "RMNUncurseForceOptionFalse",
		applyChangeset: func(env *testhelpers.DeployedEnv, config v1_6.RMNCurseConfig, t *testing.T) deployment.ChangesetOutput {
			cs, err := v1_6.RMNUncurseChangeset(env.Env, config)
			require.NoError(t, err)
			return cs
		},
		force:          false,
		expectProposal: false,
	},
	{
		name: "RMNCurseForceOptionFalse",
		applyChangeset: func(env *testhelpers.DeployedEnv, config v1_6.RMNCurseConfig, t *testing.T) deployment.ChangesetOutput {
			_, _, err := commonchangeset.ApplyChangesetsV2(t, env.Env,
				[]commonchangeset.ConfiguredChangeSet{
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(v1_6.RMNCurseChangeset),
						config,
					)},
			)
			require.NoError(t, err)

			cs, err := v1_6.RMNCurseChangeset(env.Env, config)
			require.NoError(t, err)

			return cs
		},
		force:          false,
		expectProposal: false,
	},
	{
		name: "RMNUncurseForceOptionTrue",
		applyChangeset: func(env *testhelpers.DeployedEnv, config v1_6.RMNCurseConfig, t *testing.T) deployment.ChangesetOutput {
			cs, err := v1_6.RMNUncurseChangeset(env.Env, config)
			require.NoError(t, err)
			return cs
		},
		force:              true,
		expectProposal:     true,
		expectedOperations: 3, // 3 operations for 3 chains
	},
	{
		name: "RMNCurseForceOptionTrue",
		applyChangeset: func(env *testhelpers.DeployedEnv, config v1_6.RMNCurseConfig, t *testing.T) deployment.ChangesetOutput {
			_, _, err := commonchangeset.ApplyChangesetsV2(t, env.Env,
				[]commonchangeset.ConfiguredChangeSet{
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(v1_6.RMNCurseChangeset),
						config,
					)},
			)
			require.NoError(t, err)

			cs, err := v1_6.RMNCurseChangeset(env.Env, config)
			require.NoError(t, err)

			return cs
		},
		force:              true,
		expectProposal:     true,
		expectedOperations: 3, // 3 operations for 3 chains
	},
}

func TestGetAllCursableChainsEmptyWhenNoRMNRemote(t *testing.T) {
	e, _ := testhelpers.NewMemoryEnvironment(
		t, testhelpers.WithNumOfChains(2), testhelpers.WithSolChains(1),
		testhelpers.WithPrerequisiteDeploymentOnly(nil),
	)

	cursableChains, err := v1_6.GetCursableChains(e.Env)
	require.NoError(t, err)
	require.NotNil(t, cursableChains)
	require.Empty(t, cursableChains)
}

func TestGetAllCursableChainsWithRMNRemote(t *testing.T) {
	e, _ := testhelpers.NewMemoryEnvironment(
		t, testhelpers.WithNumOfChains(2), testhelpers.WithSolChains(1),
	)

	cursableChains, err := v1_6.GetCursableChains(e.Env)
	require.NoError(t, err)
	require.NotNil(t, cursableChains)
	require.Len(t, cursableChains, 3)
}

func TestRMNUncurseForceOption(t *testing.T) {
	for _, tc := range forceOptionTestCases {
		t.Run(tc.name, func(t *testing.T) {
			e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2), testhelpers.WithSolChains(1))

			mapIDToSelector := func(id uint64) uint64 {
				return v1_6.GetAllCursableChainsSelector(e.Env)[id]
			}

			state, err := stateview.LoadOnchainState(e.Env)
			require.NoError(t, err)

			transferRMNContractToMCMS(t, &e, state)

			// Apply a curse changeset to create an active curse
			config := v1_6.RMNCurseConfig{
				CurseActions:             []v1_6.CurseAction{v1_6.CurseChain(mapIDToSelector(Evm1))},
				Reason:                   "test curse",
				IncludeNotConnectedLanes: true,
				MCMS: &proposalutils.TimelockConfig{
					MinDelay: 1 * time.Second,
				},
				Force: tc.force,
			}

			cs := tc.applyChangeset(&e, config, t)

			if tc.expectProposal {
				require.Len(t, cs.MCMSTimelockProposals, 1)
				proposal := cs.MCMSTimelockProposals[0]
				require.Len(t, proposal.Operations, tc.expectedOperations)
			} else {
				require.Empty(t, cs.MCMSTimelockProposals)
			}
		})
	}
}
