package ring

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
)

func TestStateTransitionDeterminism(t *testing.T) {
	now := time.Unix(0, 0)
	timeToSync := 5 * time.Minute

	current := &ringpb.RoutingState{
		Id:    1,
		State: &ringpb.RoutingState_RoutableShards{RoutableShards: 2},
	}

	// Same inputs should produce identical outputs
	result1, err := NextState(current, 4, now, timeToSync)
	require.NoError(t, err)

	result2, err := NextState(current, 4, now, timeToSync)
	require.NoError(t, err)

	require.Equal(t, result1.Id, result2.Id)
	require.Equal(t, result1.GetTransition().WantShards, result2.GetTransition().WantShards)
	require.Equal(t, result1.GetTransition().LastStableCount, result2.GetTransition().LastStableCount)
	require.Equal(t, result1.GetTransition().ChangesSafeAfter.AsTime(), result2.GetTransition().ChangesSafeAfter.AsTime())
}

// ∀ state, inputs: NextState(state, inputs).Id >= state.Id
func TestFV_StateIDMonotonicity(t *testing.T) {
	timeToSync := 5 * time.Minute
	baseTime := time.Unix(0, 0)

	testCases := []struct {
		name  string
		state *ringpb.RoutingState
		now   time.Time
	}{
		// Steady state cases
		{"steady_same_shards", steadyState(10, 3), baseTime},
		{"steady_more_shards", steadyState(10, 3), baseTime},
		{"steady_fewer_shards", steadyState(10, 3), baseTime},
		// Transition state cases
		{"transition_before_safe", transitionState(10, 3, 5, baseTime.Add(1*time.Hour)), baseTime},
		{"transition_at_safe", transitionState(10, 3, 5, baseTime), baseTime},
		{"transition_after_safe", transitionState(10, 3, 5, baseTime.Add(-1*time.Second)), baseTime},
	}

	shardCounts := []uint32{1, 2, 3, 5, 10}

	for _, tc := range testCases {
		for _, wantShards := range shardCounts {
			t.Run(tc.name, func(t *testing.T) {
				result, err := NextState(tc.state, wantShards, tc.now, timeToSync)
				require.NoError(t, err)

				// INVARIANT: ID never decreases
				require.GreaterOrEqual(t, result.Id, tc.state.Id,
					"state ID must be monotonically non-decreasing")
			})
		}
	}
}

// The state machine only produces valid transitions:
//   - Steady → Steady (when shards unchanged)
//   - Steady → Transition (when shards change)
//   - Transition → Transition (before safety period)
//   - Transition → Steady (after safety period)
func TestFV_ValidStateTransitions(t *testing.T) {
	timeToSync := 5 * time.Minute
	baseTime := time.Unix(0, 0)

	t.Run("steady_to_steady_when_unchanged", func(t *testing.T) {
		for _, shards := range []uint32{1, 2, 3, 5, 10} {
			state := steadyState(1, shards)
			result, err := NextState(state, shards, baseTime, timeToSync)
			require.NoError(t, err)

			// Must remain steady with same shard count
			require.True(t, IsInSteadyState(result))
			require.Equal(t, shards, result.GetRoutableShards())
			require.Equal(t, state.Id, result.Id, "ID unchanged when no transition")
		}
	})

	t.Run("steady_to_transition_when_changed", func(t *testing.T) {
		transitions := [][2]uint32{{1, 2}, {2, 1}, {3, 5}, {5, 3}, {1, 10}}
		for _, tr := range transitions {
			current, want := tr[0], tr[1]
			state := steadyState(1, current)
			result, err := NextState(state, want, baseTime, timeToSync)
			require.NoError(t, err)

			// Must enter transition
			require.False(t, IsInSteadyState(result))
			require.NotNil(t, result.GetTransition())
			require.Equal(t, want, result.GetTransition().WantShards)
			require.Equal(t, current, result.GetTransition().LastStableCount)
			require.Equal(t, state.Id+1, result.Id)
		}
	})

	t.Run("transition_stays_before_safe_time", func(t *testing.T) {
		safeAfter := baseTime.Add(1 * time.Hour)
		for _, wantShards := range []uint32{1, 2, 5} {
			state := transitionState(5, 2, wantShards, safeAfter)
			result, err := NextState(state, wantShards, baseTime, timeToSync)
			require.NoError(t, err)

			// Must remain in transition
			require.False(t, IsInSteadyState(result))
			require.Equal(t, state.Id, result.Id, "ID unchanged while waiting")
		}
	})

	t.Run("transition_completes_after_safe_time", func(t *testing.T) {
		safeAfter := baseTime.Add(-1 * time.Second)
		for _, wantShards := range []uint32{1, 2, 5} {
			state := transitionState(5, 2, wantShards, safeAfter)
			result, err := NextState(state, wantShards, baseTime, timeToSync)
			require.NoError(t, err)

			// Must complete to steady
			require.True(t, IsInSteadyState(result))
			require.Equal(t, wantShards, result.GetRoutableShards())
			require.Equal(t, state.Id+1, result.Id)
		}
	})
}

// ∀ transition: completion occurs iff now >= safeAfter
func TestFV_SafetyPeriodEnforcement(t *testing.T) {
	timeToSync := 5 * time.Minute
	baseTime := time.Unix(0, 0)

	// Test various time offsets relative to safeAfter
	offsets := []time.Duration{
		-1 * time.Hour,
		-1 * time.Minute,
		-1 * time.Second,
		-1 * time.Nanosecond,
		0,
		1 * time.Nanosecond,
		1 * time.Second,
		1 * time.Minute,
		1 * time.Hour,
	}

	for _, offset := range offsets {
		safeAfter := baseTime
		now := baseTime.Add(offset)
		state := transitionState(1, 2, 5, safeAfter)

		result, err := NextState(state, 5, now, timeToSync)
		require.NoError(t, err)

		shouldComplete := !now.Before(safeAfter)
		didComplete := IsInSteadyState(result)

		require.Equal(t, shouldComplete, didComplete,
			"offset=%v: safety period enforcement failed", offset)
	}
}

// When entering transition, WantShards equals the requested shard count
// When completing transition, final shard count equals WantShards
func TestFV_TransitionPreservesTarget(t *testing.T) {
	timeToSync := 5 * time.Minute
	baseTime := time.Unix(0, 0)

	for _, currentShards := range []uint32{1, 2, 3, 5} {
		for _, wantShards := range []uint32{1, 2, 3, 5} {
			if currentShards == wantShards {
				continue // No transition occurs
			}

			// Step 1: Enter transition
			state := steadyState(0, currentShards)
			afterEnter, err := NextState(state, wantShards, baseTime, timeToSync)
			require.NoError(t, err)
			require.Equal(t, wantShards, afterEnter.GetTransition().WantShards,
				"transition must preserve target shard count")

			// Step 2: Complete transition (after safety period)
			afterComplete, err := NextState(afterEnter, wantShards, baseTime.Add(timeToSync+time.Second), timeToSync)
			require.NoError(t, err)
			require.Equal(t, wantShards, afterComplete.GetRoutableShards(),
				"completed state must have target shard count")
		}
	}
}

// ∀ transition: ∃ time t where transition completes (no infinite loops)
func TestFV_EventualCompletion(t *testing.T) {
	timeToSync := 5 * time.Minute
	baseTime := time.Unix(0, 0)

	state := steadyState(0, 2)

	// Enter transition
	state, err := NextState(state, 5, baseTime, timeToSync)
	require.NoError(t, err)
	require.False(t, IsInSteadyState(state))

	// Simulate time progression - must complete within safety period
	completionTime := baseTime.Add(timeToSync)
	state, err = NextState(state, 5, completionTime, timeToSync)
	require.NoError(t, err)

	require.True(t, IsInSteadyState(state), "transition must eventually complete")
}

// ∀ state: exactly one of (IsInSteadyState, IsInTransition) is true
func TestFV_StateTypeExclusivity(t *testing.T) {
	states := []*ringpb.RoutingState{
		steadyState(0, 1),
		steadyState(5, 3),
		transitionState(0, 1, 2, time.Now()),
		transitionState(5, 3, 5, time.Now().Add(time.Hour)),
	}

	for i, state := range states {
		isSteady := IsInSteadyState(state)
		_, isTransition := state.State.(*ringpb.RoutingState_Transition)

		require.NotEqual(t, isSteady, isTransition,
			"state %d: exactly one state type must be true", i)
	}
}

// IsInSteadyState(nil) = false (safe handling of nil)
// NextState(nil, ...) returns error (explicit failure)
func TestFV_NilStateSafety(t *testing.T) {
	require.False(t, IsInSteadyState(nil), "nil state must not be steady")

	_, err := NextState(nil, 1, time.Now(), time.Minute)
	require.Error(t, err, "NextState must reject nil input")
}

func steadyState(id uint64, shards uint32) *ringpb.RoutingState {
	return &ringpb.RoutingState{
		Id:    id,
		State: &ringpb.RoutingState_RoutableShards{RoutableShards: shards},
	}
}

func transitionState(id uint64, lastStable, wantShards uint32, safeAfter time.Time) *ringpb.RoutingState {
	return &ringpb.RoutingState{
		Id: id,
		State: &ringpb.RoutingState_Transition{
			Transition: &ringpb.Transition{
				WantShards:       wantShards,
				LastStableCount:  lastStable,
				ChangesSafeAfter: timestamppb.New(safeAfter),
			},
		},
	}
}
