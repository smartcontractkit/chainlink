package ring

import (
	"maps"
	"math"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
)

const (
	fuzzN = 4
	fuzzF = 1
	// An honest arbiter reports a small shard set; consistent hashing panics beyond its 997 partitions, which a
	// healthy-shard quorum (needing an honest vote) can never reach.
	maxFuzzShardIDs = 64
)

func newFuzzPlugin(t *testing.T) *Plugin {
	t.Helper()
	cfg := ocr3types.ReportingPluginConfig{N: fuzzN, F: fuzzF}
	p, err := NewPlugin(NewStore(), &mockArbiter{}, cfg, logger.Nop(), nil)
	require.NoError(t, err)
	return p
}

func seedObservation(t testing.TB, wantShards uint32, workflows []string) []byte {
	b, err := proto.Marshal(&ringpb.Observation{
		ShardStatus: toShardStatus(map[uint32]bool{0: true, 1: true, 2: false}),
		WorkflowIds: workflows,
		Now:         timestamppb.New(time.Unix(1_700_000_000, 0)),
		WantShards:  wantShards,
	})
	require.NoError(t, err)
	return b
}

func seedObservationWith(t testing.TB, obs *ringpb.Observation) []byte {
	b, err := proto.Marshal(obs)
	require.NoError(t, err)
	return b
}

func priorOutcome(t *testing.T, seqNr uint64, shards uint32, transition bool, safeAfterUnix int64) ocr3types.Outcome {
	t.Helper()
	if shards == 0 {
		return nil
	}
	state := &ringpb.RoutingState{Id: seqNr, State: &ringpb.RoutingState_RoutableShards{RoutableShards: shards}}
	if transition {
		state.State = &ringpb.RoutingState_Transition{Transition: &ringpb.Transition{
			WantShards:       shards,
			LastStableCount:  shards + 1,
			ChangesSafeAfter: timestamppb.New(time.Unix(safeAfterUnix, 0)),
		}}
	}
	b, err := proto.MarshalOptions{Deterministic: true}.Marshal(&ringpb.Outcome{State: state})
	require.NoError(t, err)
	return b
}

func parseObservations(aos []types.AttributedObservation) []*ringpb.Observation {
	out := make([]*ringpb.Observation, 0, len(aos))
	for _, ao := range aos {
		obs := &ringpb.Observation{}
		if proto.Unmarshal(ao.Observation, obs) == nil {
			out = append(out, obs)
		}
	}
	return out
}

func distinctShardIDs(observations []*ringpb.Observation) int {
	ids := map[uint32]struct{}{}
	for _, obs := range observations {
		for id := range obs.ShardStatus {
			ids[id] = struct{}{}
		}
	}
	return len(ids)
}

// healthyShards is the set of shards reported healthy by more than F oracles.
func healthyShards(observations []*ringpb.Observation) []uint32 {
	votes := map[uint32]int{}
	for _, obs := range observations {
		for id, status := range obs.ShardStatus {
			if status != nil && status.IsHealthy {
				votes[id]++
			}
		}
	}
	var healthy []uint32
	for id, n := range votes {
		if n > fuzzF {
			healthy = append(healthy, id)
		}
	}
	slices.Sort(healthy)
	return healthy
}

func medianWantShards(observations []*ringpb.Observation) uint32 {
	votes := make([]uint32, 0, len(observations))
	for _, obs := range observations {
		votes = append(votes, obs.WantShards)
	}
	slices.Sort(votes)
	return votes[len(votes)/2]
}

func medianNow(observations []*ringpb.Observation) time.Time {
	nows := make([]time.Time, 0, len(observations))
	for _, obs := range observations {
		nows = append(nows, obs.Now.AsTime())
	}
	slices.SortFunc(nows, time.Time.Compare)
	return nows[len(nows)/2]
}

// expectedNextState spells out the routing state machine: a steady state moves to a transition when the agreed shard
// count changes, and a transition settles only once its safe time has passed.
func expectedNextState(prior *ringpb.RoutingState, wantShards uint32, now time.Time) *ringpb.RoutingState {
	switch s := prior.State.(type) {
	case *ringpb.RoutingState_RoutableShards:
		if s.RoutableShards == wantShards {
			return &ringpb.RoutingState{Id: prior.Id, State: &ringpb.RoutingState_RoutableShards{RoutableShards: wantShards}}
		}
		return &ringpb.RoutingState{Id: prior.Id + 1, State: &ringpb.RoutingState_Transition{Transition: &ringpb.Transition{
			WantShards:       wantShards,
			LastStableCount:  s.RoutableShards,
			ChangesSafeAfter: timestamppb.New(now.Add(DefaultTimeToSync)),
		}}}
	case *ringpb.RoutingState_Transition:
		if now.Before(s.Transition.ChangesSafeAfter.AsTime()) {
			return prior
		}
		return &ringpb.RoutingState{Id: prior.Id + 1, State: &ringpb.RoutingState_RoutableShards{RoutableShards: s.Transition.WantShards}}
	}
	return nil
}

func requireOutcomeInvariants(t *testing.T, raw []byte, prev ocr3types.Outcome, seqNr uint64, observations []*ringpb.Observation) {
	t.Helper()
	outcome := &ringpb.Outcome{}
	require.NoError(t, proto.Unmarshal(raw, outcome))
	require.NotNil(t, outcome.State)

	workflows := map[string]struct{}{}
	for _, obs := range observations {
		for _, id := range obs.WorkflowIds {
			workflows[id] = struct{}{}
		}
	}
	require.Equal(t, slices.Sorted(maps.Keys(workflows)), slices.Sorted(maps.Keys(outcome.Routes)), "every observed workflow must be routed, and only those")

	healthy := healthyShards(observations)
	for wfID, route := range outcome.Routes {
		if len(healthy) == 0 {
			require.Zero(t, route.DonId, "workflow %q must fall back to shard 0 when no shard is healthy", wfID)
			continue
		}
		want, err := locateShard(newShardRing(healthy), wfID)
		require.NoError(t, err)
		require.Equal(t, want, route.DonId, "workflow %q must be placed by consistent hashing over the healthy shards", wfID)
	}

	wantShards := medianWantShards(observations)
	prior := &ringpb.Outcome{State: &ringpb.RoutingState{Id: seqNr, State: &ringpb.RoutingState_RoutableShards{RoutableShards: wantShards}}}
	if prev != nil {
		prior = &ringpb.Outcome{}
		require.NoError(t, proto.Unmarshal(prev, prior))
	}
	expected := expectedNextState(prior.State, wantShards, medianNow(observations))
	require.True(t, proto.Equal(expected, outcome.State), "routing state must follow the transition rules: want %v, got %v", expected, outcome.State)
}

func FuzzPlugin_OutcomeInvariants(f *testing.F) {
	obsA := seedObservation(f, 2, []string{"a", "b"})
	obsB := seedObservation(f, 3, []string{"b", "c"})
	obsC := seedObservation(f, 2, nil)
	obsWant3 := seedObservation(f, 3, []string{"a"})
	f.Add(uint64(1), uint32(2), false, int64(0), obsA, obsB, obsC, []byte{})
	f.Add(uint64(7), uint32(3), true, int64(1_700_000_000+3600), obsA, obsB, obsC, obsA)
	f.Add(uint64(7), uint32(3), true, int64(1_600_000_000), obsA, obsB, obsC, obsA)
	f.Add(uint64(9), uint32(2), false, int64(0), obsWant3, obsWant3, obsWant3, obsWant3)
	// shard 5 is healthy according to a single oracle, which is not enough to route to it
	obsShard5 := seedObservationWith(f, &ringpb.Observation{
		ShardStatus: toShardStatus(map[uint32]bool{0: true, 1: true, 5: true}),
		WorkflowIds: []string{"a", "b", "c", "d"},
		Now:         timestamppb.New(time.Unix(1_700_000_000, 0)),
		WantShards:  2,
	})
	f.Add(uint64(3), uint32(2), false, int64(0), obsShard5, obsA, obsA, obsC)
	f.Add(uint64(0), uint32(0), false, int64(0), []byte{}, []byte{0xff, 0xff}, []byte{0x0a}, []byte{0x12, 0x00})

	f.Fuzz(func(t *testing.T, seqNr uint64, priorShards uint32, transition bool, safeAfterUnix int64, o0, o1, o2, o3 []byte) {
		var aos []types.AttributedObservation
		validator := newFuzzPlugin(t)
		for i, raw := range [][]byte{o0, o1, o2, o3} {
			ao := types.AttributedObservation{Observation: raw, Observer: commontypes.OracleID(i)}
			if validator.ValidateObservation(t.Context(), ocr3types.OutcomeContext{}, nil, ao) == nil {
				aos = append(aos, ao)
			}
		}
		if len(aos) == 0 {
			return
		}
		observations := parseObservations(aos)
		if distinctShardIDs(observations) > maxFuzzShardIDs {
			t.Skip("more shards than an arbiter reports")
		}

		prev := priorOutcome(t, seqNr, priorShards, transition, safeAfterUnix)
		outctx := ocr3types.OutcomeContext{SeqNr: seqNr, PreviousOutcome: prev}

		first, err1 := newFuzzPlugin(t).Outcome(t.Context(), outctx, nil, aos)
		second, err2 := newFuzzPlugin(t).Outcome(t.Context(), outctx, nil, aos)
		require.Equal(t, err1 == nil, err2 == nil)
		require.Equal(t, first, second, "same inputs must yield identical outcome")
		if err1 != nil {
			return
		}

		if seqNr != math.MaxUint64 {
			requireOutcomeInvariants(t, first, prev, seqNr, observations)
		}

		reports, err := newFuzzPlugin(t).Reports(t.Context(), seqNr, first)
		require.NoError(t, err)
		require.Len(t, reports, 1)
		require.Equal(t, first, ocr3types.Outcome(reports[0].ReportWithInfo.Report), "report must carry the outcome verbatim")
		again, err := newFuzzPlugin(t).Reports(t.Context(), seqNr, first)
		require.NoError(t, err)
		require.Equal(t, reports, again, "same outcome must yield identical reports")
	})
}

func FuzzPlugin_ValidatedObservationYieldsOutcome(f *testing.F) {
	f.Add(seedObservation(f, 2, []string{"a"}))
	f.Add([]byte{})
	f.Add([]byte{0x10, 0x01})
	f.Add(seedObservationWith(f, &ringpb.Observation{WantShards: 2}))                         // no timestamp
	f.Add(seedObservationWith(f, &ringpb.Observation{Now: timestamppb.New(time.Unix(1, 0))})) // no shard count

	f.Fuzz(func(t *testing.T, raw []byte) {
		p := newFuzzPlugin(t)
		ao := types.AttributedObservation{Observation: raw}
		if err := p.ValidateObservation(t.Context(), ocr3types.OutcomeContext{}, nil, ao); err != nil {
			return
		}

		observation := &ringpb.Observation{}
		require.NoError(t, proto.Unmarshal(raw, observation))
		require.NotNil(t, observation.Now, "a validated observation must carry a timestamp")
		require.NotZero(t, observation.WantShards, "a validated observation must carry a shard count")

		out, err := p.Outcome(t.Context(), ocr3types.OutcomeContext{SeqNr: 1}, nil, []types.AttributedObservation{ao})
		require.NoError(t, err, "a validated observation must produce an outcome")
		requireOutcomeInvariants(t, out, nil, 1, []*ringpb.Observation{observation})
	})
}
