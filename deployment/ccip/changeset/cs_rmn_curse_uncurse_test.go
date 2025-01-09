package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/stretchr/testify/require"
)

type curseAssertion struct {
	chainId      uint64
	subject      uint64
	global_curse bool
	cursed       bool
}

type CurseTestCase struct {
	name                string
	curseActionsBuilder func(mapIdToSelectorFunc) []CurseAction
	curseAssertions     []curseAssertion
}

type mapIdToSelectorFunc func(uint64) uint64

var testCases = []CurseTestCase{
	{
		name: "lane",
		curseActionsBuilder: func(mapIdToSelector mapIdToSelectorFunc) []CurseAction {
			return []CurseAction{CurseLane(mapIdToSelector(0), mapIdToSelector(1))}
		},
		curseAssertions: []curseAssertion{
			{chainId: 0, subject: 1, cursed: true},
			{chainId: 0, subject: 2, cursed: false},
			{chainId: 1, subject: 0, cursed: true},
			{chainId: 1, subject: 2, cursed: false},
			{chainId: 2, subject: 0, cursed: false},
			{chainId: 2, subject: 1, cursed: false},
		},
	},
	{
		name: "lane duplicate",
		curseActionsBuilder: func(mapIdToSelector mapIdToSelectorFunc) []CurseAction {
			return []CurseAction{CurseLane(mapIdToSelector(0), mapIdToSelector(1)), CurseLane(mapIdToSelector(0), mapIdToSelector(1))}
		},
		curseAssertions: []curseAssertion{
			{chainId: 0, subject: 1, cursed: true},
			{chainId: 0, subject: 2, cursed: false},
			{chainId: 1, subject: 0, cursed: true},
			{chainId: 1, subject: 2, cursed: false},
			{chainId: 2, subject: 0, cursed: false},
			{chainId: 2, subject: 1, cursed: false},
		},
	},
	{
		name: "chain",
		curseActionsBuilder: func(mapIdToSelector mapIdToSelectorFunc) []CurseAction {
			return []CurseAction{CurseChain(mapIdToSelector(0))}
		},
		curseAssertions: []curseAssertion{
			{chainId: 0, global_curse: true, cursed: true},
			{chainId: 1, subject: 0, cursed: true},
			{chainId: 1, subject: 2, cursed: false},
			{chainId: 2, subject: 0, cursed: true},
			{chainId: 2, subject: 1, cursed: false},
		},
	},
	{
		name: "chain duplicate",
		curseActionsBuilder: func(mapIdToSelector mapIdToSelectorFunc) []CurseAction {
			return []CurseAction{CurseChain(mapIdToSelector(0)), CurseChain(mapIdToSelector(0))}
		},
		curseAssertions: []curseAssertion{
			{chainId: 0, global_curse: true, cursed: true},
			{chainId: 1, subject: 0, cursed: true},
			{chainId: 1, subject: 2, cursed: false},
			{chainId: 2, subject: 0, cursed: true},
			{chainId: 2, subject: 1, cursed: false},
		},
	},
	{
		name: "chain and lanes",
		curseActionsBuilder: func(mapIdToSelector mapIdToSelectorFunc) []CurseAction {
			return []CurseAction{CurseChain(mapIdToSelector(0)), CurseLane(mapIdToSelector(1), mapIdToSelector(2))}
		},
		curseAssertions: []curseAssertion{
			{chainId: 0, global_curse: true, cursed: true},
			{chainId: 1, subject: 0, cursed: true},
			{chainId: 1, subject: 2, cursed: true},
			{chainId: 2, subject: 0, cursed: true},
			{chainId: 2, subject: 1, cursed: true},
		},
	},
}

func TestRMNCurse(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name+"_NO_MCMS", func(t *testing.T) {
			testRmnCurse(t, tc)
		})
		t.Run(tc.name+"_MCMS", func(t *testing.T) {
			testRmnCurseMCMS(t, tc)
		})
	}
}

func TestRMNUncurse(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name+"_UNCURSE", func(t *testing.T) {
			testRmnUncurse(t, tc)
		})
		t.Run(tc.name+"_UNCURSE_MCMS", func(t *testing.T) {
			testRmnUncurseMCMS(t, tc)
		})
	}
}

func TestRMNCurseConfigValidate(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name+"_VALIDATE", func(t *testing.T) {
			testRmnCurseConfigValidate(t, tc)
		})
	}
}

func testRmnUncurse(t *testing.T, tc CurseTestCase) {
	e := NewMemoryEnvironment(t, WithChains(3))

	mapIdToSelector := func(id uint64) uint64 {
		return e.Env.AllChainSelectors()[id]
	}

	verifyNoActiveCurseOnAllChains(t, &e)

	config := RMNCurseConfig{
		CurseActions: tc.curseActionsBuilder(mapIdToSelector),
		Reason:       "test curse",
	}

	_, err := NewRMNCurseChangeset(e.Env, config)
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, tc, mapIdToSelector)

	_, err = NewRMNUncurseChangeset(e.Env, config)
	require.NoError(t, err)

	verifyNoActiveCurseOnAllChains(t, &e)
}

func transferRMNContractToMCMS(t *testing.T, e *DeployedEnv, state CCIPOnChainState, timelocksPerChain map[uint64]*proposalutils.TimelockExecutionContracts) {
	contractsByChain := make(map[uint64][]common.Address)
	rmnRemoteAddressesByChain := buildRMNRemoteAddressPerChain(e.Env, state)
	for chainSelector, rmnRemoteAddress := range rmnRemoteAddressesByChain {
		contractsByChain[chainSelector] = []common.Address{rmnRemoteAddress}
	}

	contractsByChain[e.HomeChainSel] = append(contractsByChain[e.HomeChainSel], state.Chains[e.HomeChainSel].RMNHome.Address())

	// This is required because RMN Contracts is initially owned by the deployer
	_, err := commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.TransferToMCMSWithTimelock),
			Config: commonchangeset.TransferToMCMSWithTimelockConfig{
				ContractsByChain: contractsByChain,
				MinDelay:         0,
			},
		},
	})
	require.NoError(t, err)
}

func testRmnUncurseMCMS(t *testing.T, tc CurseTestCase) {
	e := NewMemoryEnvironment(t, WithChains(3))

	mapIdToSelector := func(id uint64) uint64 {
		return e.Env.AllChainSelectors()[id]
	}

	config := RMNCurseConfig{
		CurseActions: tc.curseActionsBuilder(mapIdToSelector),
		Reason:       "test curse",
		MCMS:         &MCMSConfig{MinDelay: 0},
	}

	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)

	verifyNoActiveCurseOnAllChains(t, &e)

	timelocksPerChain := buildTimelockPerChain(e.Env, state)

	transferRMNContractToMCMS(t, &e, state, timelocksPerChain)

	_, err = commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(NewRMNCurseChangeset),
			Config:    config,
		},
	})
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, tc, mapIdToSelector)

	_, err = commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(NewRMNUncurseChangeset),
			Config:    config,
		},
	})
	require.NoError(t, err)

	verifyNoActiveCurseOnAllChains(t, &e)
}

func testRmnCurseConfigValidate(t *testing.T, tc CurseTestCase) {
	e := NewMemoryEnvironment(t, WithChains(3))

	mapIdToSelector := func(id uint64) uint64 {
		return e.Env.AllChainSelectors()[id]
	}

	config := RMNCurseConfig{
		CurseActions: tc.curseActionsBuilder(mapIdToSelector),
		Reason:       "test curse",
	}

	err := config.Validate(e.Env)
	require.NoError(t, err)
}

func testRmnCurse(t *testing.T, tc CurseTestCase) {
	e := NewMemoryEnvironment(t, WithChains(3))

	mapIdToSelector := func(id uint64) uint64 {
		return e.Env.AllChainSelectors()[id]
	}

	verifyNoActiveCurseOnAllChains(t, &e)

	config := RMNCurseConfig{
		CurseActions: tc.curseActionsBuilder(mapIdToSelector),
		Reason:       "test curse",
	}

	_, err := NewRMNCurseChangeset(e.Env, config)
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, tc, mapIdToSelector)
}

func testRmnCurseMCMS(t *testing.T, tc CurseTestCase) {
	e := NewMemoryEnvironment(t, WithChains(3))

	mapIdToSelector := func(id uint64) uint64 {
		return e.Env.AllChainSelectors()[id]
	}

	config := RMNCurseConfig{
		CurseActions: tc.curseActionsBuilder(mapIdToSelector),
		Reason:       "test curse",
		MCMS:         &MCMSConfig{MinDelay: 0},
	}

	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)

	verifyNoActiveCurseOnAllChains(t, &e)

	timelocksPerChain := buildTimelockPerChain(e.Env, state)

	transferRMNContractToMCMS(t, &e, state, timelocksPerChain)

	_, err = commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(NewRMNCurseChangeset),
			Config:    config,
		},
	})
	require.NoError(t, err)

	verifyTestCaseAssertions(t, &e, tc, mapIdToSelector)
}

func verifyTestCaseAssertions(t *testing.T, e *DeployedEnv, tc CurseTestCase, mapIdToSelector mapIdToSelectorFunc) {
	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)

	for _, assertion := range tc.curseAssertions {
		cursedSubject := SelectorToSubject(mapIdToSelector(assertion.subject))
		if assertion.global_curse {
			cursedSubject = GlobalCurseSubject()
		}

		isCursed, err := state.Chains[mapIdToSelector(assertion.chainId)].RMNRemote.IsCursed(nil, cursedSubject)
		require.NoError(t, err)
		require.Equal(t, assertion.cursed, isCursed, "chain %d subject %d", assertion.chainId, assertion.subject)
	}
}

func verifyNoActiveCurseOnAllChains(t *testing.T, e *DeployedEnv) {
	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)

	for _, chain := range e.Env.Chains {
		isCursed, err := state.Chains[chain.Selector].RMNRemote.IsCursed0(nil)
		require.NoError(t, err)
		require.False(t, isCursed, "chain %d", chain.Selector)
	}
}
