package ring

import (
	"math"
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

func healthyShardSet(aos []types.AttributedObservation) map[uint32]bool {
	votes := map[uint32]int{}
	for _, ao := range aos {
		obs := &ringpb.Observation{}
		if proto.Unmarshal(ao.Observation, obs) != nil {
			continue
		}
		for id, status := range obs.ShardStatus {
			if status != nil && status.IsHealthy {
				votes[id]++
			}
		}
	}
	healthy := map[uint32]bool{}
	for id, n := range votes {
		if n > fuzzF {
			healthy[id] = true
		}
	}
	return healthy
}

func FuzzPlugin_OutcomeInvariants(f *testing.F) {
	obsA := seedObservation(f, 2, []string{"a", "b"})
	obsB := seedObservation(f, 3, []string{"b", "c"})
	obsC := seedObservation(f, 2, nil)
	f.Add(uint64(1), uint32(2), false, int64(0), obsA, obsB, obsC, []byte{})
	f.Add(uint64(7), uint32(3), true, int64(1_700_000_000+3600), obsA, obsB, obsC, obsA)
	f.Add(uint64(7), uint32(3), true, int64(1_600_000_000), obsA, obsB, obsC, obsA)
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

		prev := priorOutcome(t, seqNr, priorShards, transition, safeAfterUnix)
		outctx := ocr3types.OutcomeContext{SeqNr: seqNr, PreviousOutcome: prev}

		first, err1 := newFuzzPlugin(t).Outcome(t.Context(), outctx, nil, aos)
		second, err2 := newFuzzPlugin(t).Outcome(t.Context(), outctx, nil, aos)
		require.Equal(t, err1 == nil, err2 == nil)
		require.Equal(t, first, second, "same inputs must yield identical outcome")
		if err1 != nil {
			return
		}

		outcome := &ringpb.Outcome{}
		require.NoError(t, proto.Unmarshal(first, outcome))
		require.NotNil(t, outcome.State)

		healthy := healthyShardSet(aos)
		for wfID, route := range outcome.Routes {
			require.True(t, route.DonId == 0 || healthy[route.DonId], "workflow %q routed to unhealthy shard %d", wfID, route.DonId)
		}

		if prev != nil && seqNr != math.MaxUint64 {
			require.GreaterOrEqual(t, outcome.State.Id, seqNr, "state id must not go backwards")
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

	f.Fuzz(func(t *testing.T, raw []byte) {
		p := newFuzzPlugin(t)
		ao := types.AttributedObservation{Observation: raw}
		if err := p.ValidateObservation(t.Context(), ocr3types.OutcomeContext{}, nil, ao); err != nil {
			return
		}
		_, err := p.Outcome(t.Context(), ocr3types.OutcomeContext{SeqNr: 1}, nil, []types.AttributedObservation{ao})
		require.NoError(t, err, "a validated observation must produce an outcome")
	})
}
