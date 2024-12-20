package changeset

import (
	"testing"
)

type curseAssertion struct {
	chainSelector uint64
	subject       uint64
}

type CurseTestCase struct {
	useMCMS             bool
	name                string
	curseActionsBuilder []CurseAction
	curseAssertions     []curseAssertion
}

func TestRMNCurse(t *testing.T) {
	t.Parallel()
	testCases := []CurseTestCase{
		// {
		// 	useMCMS: true,
		// 	name:    "with MCMS",
		// },
		{
			useMCMS:             false,
			name:                "without MCMS",
			curseActionsBuilder: []CurseAction{CurseLane(0, 1)},
			curseAssertions: []curseAssertion{
				{chainSelector: 0, subject: 1},
				{chainSelector: 1, subject: 0},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testRmnCurse(t, tc)
		})
	}
}

func testRmnCurse(t *testing.T, tc CurseTestCase) {
	e := NewMemoryEnvironment(t, WithChains(2))

	config := RMNCurseConfig{
		HomeChainSelector: e.HomeChainSel,
		CurseActions:      tc.curseActionsBuilder,
		CurseReason:       "test curse",
	}

	if tc.useMCMS {
		config.MCMS = &MCMSConfig{
			MinDelay: 0,
		}
	}

	NewRMNCurseChangeset(e.Env, config)
}
