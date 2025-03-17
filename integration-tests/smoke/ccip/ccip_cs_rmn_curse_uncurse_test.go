package ccip

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
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
			return []v1_6.CurseAction{v1_6.CurseLaneBidirectionally(mapIDToSelector(0), mapIDToSelector(1))}
		},
		curseAssertions: []curseAssertion{
			{chainID: 0, subject: 1, cursed: true},
			{chainID: 0, subject: 2, cursed: false},
			{chainID: 1, subject: 0, cursed: true},
			{chainID: 1, subject: 2, cursed: false},
			{chainID: 2, subject: 0, cursed: false},
			{chainID: 2, subject: 1, cursed: false},
		},
	},
	{
		name: "lane duplicate",
		curseActionsBuilder: func(mapIDToSelector mapIDToSelectorFunc) []v1_6.CurseAction {
			return []v1_6.CurseAction{
				v1_6.CurseLaneBidirectionally(mapIDToSelector(0), mapIDToSelector(1)),
				v1_6.CurseLaneBidirectionally(mapIDToSelector(0), mapIDToSelector(1))}
		},
		curseAssertions: []curseAssertion{
			{chainID: 0, subject: 1, cursed: true},
			{chainID: 0, subject: 2, cursed: false},
			{chainID: 1, subject: 0, cursed: true},
			{chainID: 1, subject: 2, cursed: false},
			{chainID: 2, subject: 0, cursed: false},
			{chainID: 2, subject: 1, cursed: false},
		},
	},
	{
		name: "chain",
		curseActionsBuilder: func(mapIDToSelector mapIDToSelectorFunc) []v1_6.CurseAction {
			return []v1_6.CurseAction{v1_6.CurseChain(mapIDToSelector(0))}
		},
		curseAssertions: []curseAssertion{
			{chainID: 0, globalCurse: true, cursed: true},
			{chainID: 1, subject: 0, cursed: true},
			{chainID: 1, subject: 2, cursed: false},
			{chainID: 2, subject: 0, cursed: true},
			{chainID: 2, subject: 1, cursed: false},
		},
	},
	{
		name: "chain duplicate",
		curseActionsBuilder: func(mapIDToSelector mapIDToSelectorFunc) []v1_6.CurseAction {
			return []v1_6.CurseAction{v1_6.CurseChain(mapIDToSelector(0)), v1_6.CurseChain(mapIDToSelector(0))}
		},
		curseAssertions: []curseAssertion{
			{chainID: 0, globalCurse: true, cursed: true},
			{chainID: 1, subject: 0, cursed: true},
			{chainID: 1, subject: 2, cursed: false},
			{chainID: 2, subject: 0, cursed: true},
			{chainID: 2, subject: 1, cursed: false},
		},
	},
	{
		name: "chain and lanes",
		curseActionsBuilder: func(mapIDToSelector mapIDToSelectorFunc) []v1_6.CurseAction {
			return []v1_6.CurseAction{v1_6.CurseChain(mapIDToSelector(0)), v1_6.CurseLaneBidirectionally(mapIDToSelector(1), mapIDToSelector(2))}
		},
		curseAssertions: []curseAssertion{
			{chainID: 0, globalCurse: true, cursed: true},
			{chainID: 1, subject: 0, cursed: true},
			{chainID: 1, subject: 2, cursed: true},
			{chainID: 2, subject: 0, cursed: true},
			{chainID: 2, subject: 1, cursed: true},
		},
	},
}

func TestRMNCurse(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name+"_NO_MCMS", func(t *testing.T) {
			runRmnCurseTest(t, tc)
		})
		t.Run(tc.name+"_MCMS", func(t *testing.T) {
			runRmnCurseMCMSTest(t, tc)
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
			runRmnUncurseMCMSTest(t, tc)
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

func runRmnUncurseTest(t *testing.T, tc CurseTestCase) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))

	mapIDToSelector := func(id uint64) uint64 {
		return e.Env.AllChainSelectors()[id]
	}

	verifyNoActiveCurseOnAllChains(t, &e)

	config := v1_6.RMNCurseConfig{
		CurseActions: tc.curseActionsBuilder(mapIDToSelector),
		Reason:       "test curse",
	}

	_, err := v1_6.RMNCurseChangeset(e.Env, config)
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, tc, mapIDToSelector)

	_, err = v1_6.RMNUncurseChangeset(e.Env, config)
	require.NoError(t, err)

	verifyNoActiveCurseOnAllChains(t, &e)
}

func transferRMNContractToMCMS(t *testing.T, e *testhelpers.DeployedEnv, state changeset.CCIPOnChainState, timelocksPerChain map[uint64]*proposalutils.TimelockExecutionContracts) {
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

	contractsByChain[e.HomeChainSel] = append(contractsByChain[e.HomeChainSel], state.Chains[e.HomeChainSel].RMNHome.Address())

	// This is required because RMN Contracts is initially owned by the deployer
	_, err := commonchangeset.Apply(t, e.Env, timelocksPerChain,
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelock),
			commonchangeset.TransferToMCMSWithTimelockConfig{
				ContractsByChain: contractsByChain,
				MinDelay:         0,
			},
		),
	)
	require.NoError(t, err)
}

func runRmnUncurseMCMSTest(t *testing.T, tc CurseTestCase) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))

	mapIDToSelector := func(id uint64) uint64 {
		return e.Env.AllChainSelectors()[id]
	}

	config := v1_6.RMNCurseConfig{
		CurseActions: tc.curseActionsBuilder(mapIDToSelector),
		Reason:       "test curse",
		MCMS:         &changeset.MCMSConfig{MinDelay: 0},
	}

	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	verifyNoActiveCurseOnAllChains(t, &e)

	timelocksPerChain := changeset.BuildTimelockPerChain(e.Env, state)

	transferRMNContractToMCMS(t, &e, state, timelocksPerChain)

	_, err = commonchangeset.Apply(t, e.Env, timelocksPerChain,
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(v1_6.RMNCurseChangeset),
			config,
		),
	)
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, tc, mapIDToSelector)

	_, err = commonchangeset.Apply(t, e.Env, timelocksPerChain,
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(v1_6.RMNUncurseChangeset),
			config,
		),
	)
	require.NoError(t, err)

	verifyNoActiveCurseOnAllChains(t, &e)
}

func runRmnCurseConfigValidateTest(t *testing.T, tc CurseTestCase) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))

	mapIDToSelector := func(id uint64) uint64 {
		return e.Env.AllChainSelectors()[id]
	}

	config := v1_6.RMNCurseConfig{
		CurseActions: tc.curseActionsBuilder(mapIDToSelector),
		Reason:       "test curse",
	}

	err := config.Validate(e.Env)
	require.NoError(t, err)
}

func runRmnCurseTest(t *testing.T, tc CurseTestCase) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))

	mapIDToSelector := func(id uint64) uint64 {
		return e.Env.AllChainSelectors()[id]
	}

	verifyNoActiveCurseOnAllChains(t, &e)

	config := v1_6.RMNCurseConfig{
		CurseActions: tc.curseActionsBuilder(mapIDToSelector),
		Reason:       "test curse",
	}

	_, err := v1_6.RMNCurseChangeset(e.Env, config)
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, tc, mapIDToSelector)
}

func runRmnCurseIdempotentTest(t *testing.T, tc CurseTestCase) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))

	mapIDToSelector := func(id uint64) uint64 {
		return e.Env.AllChainSelectors()[id]
	}

	verifyNoActiveCurseOnAllChains(t, &e)

	config := v1_6.RMNCurseConfig{
		CurseActions: tc.curseActionsBuilder(mapIDToSelector),
		Reason:       "test curse",
	}

	_, err := v1_6.RMNCurseChangeset(e.Env, config)
	require.NoError(t, err)

	_, err = v1_6.RMNCurseChangeset(e.Env, config)
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, tc, mapIDToSelector)
}

func runRmnUncurseIdempotentTest(t *testing.T, tc CurseTestCase) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))

	mapIDToSelector := func(id uint64) uint64 {
		return e.Env.AllChainSelectors()[id]
	}

	verifyNoActiveCurseOnAllChains(t, &e)

	config := v1_6.RMNCurseConfig{
		CurseActions: tc.curseActionsBuilder(mapIDToSelector),
		Reason:       "test curse",
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

func runRmnCurseMCMSTest(t *testing.T, tc CurseTestCase) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))

	mapIDToSelector := func(id uint64) uint64 {
		return e.Env.AllChainSelectors()[id]
	}

	config := v1_6.RMNCurseConfig{
		CurseActions: tc.curseActionsBuilder(mapIDToSelector),
		Reason:       "test curse",
		MCMS:         &changeset.MCMSConfig{MinDelay: 0},
	}

	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	verifyNoActiveCurseOnAllChains(t, &e)

	timelocksPerChain := changeset.BuildTimelockPerChain(e.Env, state)

	transferRMNContractToMCMS(t, &e, state, timelocksPerChain)

	_, err = commonchangeset.Apply(t, e.Env, timelocksPerChain,
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(v1_6.RMNCurseChangeset),
			config,
		),
	)
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, tc, mapIDToSelector)
}

func verifyTestCaseAssertions(t *testing.T, e *testhelpers.DeployedEnv, tc CurseTestCase, mapIDToSelector mapIDToSelectorFunc) {
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	for _, assertion := range tc.curseAssertions {
		cursedSubject := globals.SelectorToSubject(mapIDToSelector(assertion.subject))
		if assertion.globalCurse {
			cursedSubject = globals.GlobalCurseSubject()
		}

		isCursed, err := state.Chains[mapIDToSelector(assertion.chainID)].RMNRemote.IsCursed(nil, cursedSubject)
		require.NoError(t, err)
		require.Equal(t, assertion.cursed, isCursed, "chain %d subject %d", assertion.chainID, assertion.subject)
	}
}

func verifyNoActiveCurseOnAllChains(t *testing.T, e *testhelpers.DeployedEnv) {
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	for _, chain := range e.Env.Chains {
		isCursed, err := state.Chains[chain.Selector].RMNRemote.IsCursed0(nil)
		require.NoError(t, err)
		require.False(t, isCursed, "chain %d", chain.Selector)
	}
}
