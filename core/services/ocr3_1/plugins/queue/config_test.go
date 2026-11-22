package queue

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
)

func validLimits() ocr3_1types.ReportingPluginLimits {
	return ocr3_1types.ReportingPluginLimits{
		MaxObservationBytes:                             ocr3_1types.MaxMaxObservationBytes,
		MaxKeyValueModifiedKeys:                         numKVKeys,
		MaxKeyValueModifiedKeysPlusValuesBytes:          ocr3_1types.MaxMaxKeyValueModifiedKeysPlusValuesBytes,
		MaxBlobPayloadBytes:                             ocr3_1types.MaxMaxBlobPayloadBytes,
		MaxPerOracleUnexpiredBlobCumulativePayloadBytes: 2 * maxObservedBlobHandles * ocr3_1types.MaxMaxBlobPayloadBytes,
		MaxPerOracleUnexpiredBlobCount:                  100,
	}
}

func TestPluginConfig_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mutate  func(l *ocr3_1types.ReportingPluginLimits)
		wantErr string // empty means valid
	}{
		{name: "valid"},
		{"observation bytes unset", func(l *ocr3_1types.ReportingPluginLimits) { l.MaxObservationBytes = 0 }, "MaxObservationBytes"},
		{"observation bytes above libocr max", func(l *ocr3_1types.ReportingPluginLimits) {
			l.MaxObservationBytes = ocr3_1types.MaxMaxObservationBytes + 1
		}, "out of range"},
		{"blob payload unset", func(l *ocr3_1types.ReportingPluginLimits) { l.MaxBlobPayloadBytes = 0 }, "MaxBlobPayloadBytes"},
		{"blob payload above libocr max", func(l *ocr3_1types.ReportingPluginLimits) {
			l.MaxBlobPayloadBytes = ocr3_1types.MaxMaxBlobPayloadBytes + 1
		}, "out of range"},
		{"kv bytes unset", func(l *ocr3_1types.ReportingPluginLimits) { l.MaxKeyValueModifiedKeysPlusValuesBytes = 0 }, "MaxKeyValueModifiedKeysPlusValuesBytes"},
		{"too few modified keys", func(l *ocr3_1types.ReportingPluginLimits) { l.MaxKeyValueModifiedKeys = numKVKeys - 1 }, "MaxKeyValueModifiedKeys"},
		{"blob count unset", func(l *ocr3_1types.ReportingPluginLimits) { l.MaxPerOracleUnexpiredBlobCount = 0 }, "MaxPerOracleUnexpiredBlobCount"},
		{"blob count below one full observation", func(l *ocr3_1types.ReportingPluginLimits) {
			l.MaxPerOracleUnexpiredBlobCount = maxObservedBlobHandles - 1
		}, "MaxPerOracleUnexpiredBlobCount"},
		{"cumulative blob bytes cannot hold a full observation", func(l *ocr3_1types.ReportingPluginLimits) {
			l.MaxPerOracleUnexpiredBlobCumulativePayloadBytes = maxObservedBlobHandles*l.MaxBlobPayloadBytes - 1
		}, "MaxPerOracleUnexpiredBlobCumulativePayloadBytes"},
		{"exactly one full observation of blobs is enough", func(l *ocr3_1types.ReportingPluginLimits) {
			l.MaxPerOracleUnexpiredBlobCumulativePayloadBytes = maxObservedBlobHandles * l.MaxBlobPayloadBytes
			l.MaxPerOracleUnexpiredBlobCount = maxObservedBlobHandles
		}, ""},
		{"negative", func(l *ocr3_1types.ReportingPluginLimits) { l.MaxObservationBytes = -1 }, "MaxObservationBytes"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := &PluginConfig{Limits: validLimits()}
			if tt.mutate != nil {
				tt.mutate(&cfg.Limits)
			}
			err := cfg.Validate()
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tt.wantErr)
		})
	}

	t.Run("reports every problem", func(t *testing.T) {
		t.Parallel()
		err := (&PluginConfig{}).Validate()
		require.Error(t, err)
		for _, field := range []string{"MaxObservationBytes", "MaxBlobPayloadBytes", "MaxKeyValueModifiedKeysPlusValuesBytes", "MaxPerOracleUnexpiredBlobCount", "MaxKeyValueModifiedKeys"} {
			assert.ErrorContains(t, err, field)
		}
	})
}

// The tests below calibrate the constants in config.go (and the libocr limits they must fit within) by measuring real
// protobuf sizes. They log a table (run with -v) and assert the invariants that keep the plugin inside libocr's hard limits.
// Inputs that are assumptions, rather than measurements, are marked as such.

// hexID returns a deterministic 64 character hex string, like workflow and execution IDs.
func hexID(prefix string, i int) string {
	h := sha256.Sum256([]byte(fmt.Sprint(prefix, "-", i)))
	return hex.EncodeToString(h[:])
}

// randTS returns a current-era timestamp with random nanos, so varint sizes are realistic (most nanos need 5 bytes).
func randTS(rng *rand.Rand) *timestamppb.Timestamp {
	return &timestamppb.Timestamp{Seconds: 1_800_000_000 + rng.Int63n(100_000_000), Nanos: rng.Int31n(1_000_000_000)}
}

// summarize returns the average and maximum of f(i) over n samples.
func summarize(n int, f func(i int) int) (avg float64, maxVal int) {
	var sum int
	for i := range n {
		v := f(i)
		sum += v
		maxVal = max(maxVal, v)
	}
	return float64(sum) / float64(n), maxVal
}

// observedEventFor is an event as it appears in an observation, with distinct realistic IDs.
func observedEventFor(rng *rand.Rand, i int) *pb.ConsensusQueueEvent {
	hash := sha256.Sum256([]byte(hexID("hash", i)))
	return &pb.ConsensusQueueEvent{
		WorkflowID:   hexID("wf", i),
		ExecutionID:  hexID("exec", i),
		ObservedAt:   randTS(rng),
		Deadline:     randTS(rng),
		TimelessHash: hash[:],
		Blobs:        []*pb.CoordinatedEventBlobHandle{{Blob: uint32(rng.Intn(maxObservedBlobHandles)), Element: uint32(rng.Intn(maxObservedEvents))}},
	}
}

// oracleCounts are DON sizes to calibrate for, as (N, F=(N-1)/3).
var oracleCounts = []int{4, 7, 16, 31, 100}

func TestCalibration_ObservationFitsLimit(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(1))
	avgEvent, maxEvent := summarize(maxObservedEvents, func(i int) int { return proto.Size(observedEventFor(rng, i)) })

	// worst case: every event is for a different workflow, so each also adds a capacity entry
	worstCaseObservation := func(n int) *pb.ConsensusQueueObservation {
		f := (n - 1) / 3
		obs := &pb.ConsensusQueueObservation{
			Expiration:         randTS(rng),
			Capacity:           1<<32 - 1,
			WorkflowCapacities: map[string]uint32{},
		}
		for i := range maxObservedEvents {
			e := observedEventFor(rng, i)
			obs.Events = append(obs.Events, e)
			obs.WorkflowCapacities[e.WorkflowID] = 1<<32 - 1
		}
		for range maxObservedBlobHandles { // the real handles take up no more than the placeholders reserved for them
			obs.BlobHandles = append(obs.BlobHandles, make([]byte, ocr3_1types.BlobHandleMarshalledBytesUpperBound(n, f)))
		}
		return obs
	}

	t.Logf("event in an observation: avg %.0fB, max %dB (IDs are 64 hex chars)", avgEvent, maxEvent)
	for _, n := range append(oracleCounts, 256) {
		size := proto.Size(worstCaseObservation(n))
		headroom := ocr3_1types.MaxMaxObservationBytes - size
		t.Logf("N=%-3d %d events + %d blob handles (%dB each): %d of %d bytes (%.0f%% of libocr's maximum, %dB headroom)",
			n, maxObservedEvents, maxObservedBlobHandles, ocr3_1types.BlobHandleMarshalledBytesUpperBound(n, (n-1)/3),
			size, ocr3_1types.MaxMaxObservationBytes, 100*float64(size)/ocr3_1types.MaxMaxObservationBytes, headroom)
		if n <= 100 { // N=256 is libocr's theoretical maximum, and is only logged
			assert.LessOrEqual(t, size, ocr3_1types.MaxMaxObservationBytes,
				"a full observation (maxObservedEvents, maxObservedBlobHandles) must fit in libocr's observation limit at N=%d", n)
		}
	}
}

// eventOfPayload is an event as it appears in a blob, with an opaque trigger payload of the given size.
func eventOfPayload(i, payload int) *pb.CoordinatedEvent {
	any, err := anypb.New(wrapperspb.Bytes(make([]byte, payload)))
	if err != nil {
		panic(err)
	}
	return &pb.CoordinatedEvent{
		WorkflowID: hexID("wf", i), TriggerCapID: "evm@1.0.0", TriggerIndex: int64(i % 3),
		Event: &pb.TriggerResponse{Event: &pb.TriggerEvent{TriggerType: "evm@1.0.0", Id: hexID("evt", i) + ":" + hexID("blk", i) + ":7", Payload: any}},
	}
}

// blobsNeeded is how many blobs maxObservedEvents events with this payload take, with each blob full to blobLimit.
func blobsNeeded(payload, blobLimit int) (blobs, perBlob, element int) {
	b, err := proto.MarshalOptions{Deterministic: true}.Marshal(eventOfPayload(0, payload))
	if err != nil {
		panic(err)
	}
	element = protowire.SizeTag(coordinatedEventsFieldNumber) + protowire.SizeBytes(len(b))
	perBlob = blobLimit / element
	if perBlob == 0 {
		return 0, 0, element // does not fit at all
	}
	return (maxObservedEvents + perBlob - 1) / perBlob, perBlob, element
}

func TestCalibration_BlobPacking(t *testing.T) {
	t.Parallel()

	// Trigger payload sizes are ASSUMPTIONS (the plugin does not interpret payloads): adjust to what triggers really emit.
	kinds := []struct {
		name    string
		payload int
		typical bool // typical loads must fit maxObservedEvents in maxObservedBlobHandles blobs
	}{
		{"cron tick", 128, true},
		{"EVM log", 1 << 10, true},
		{"EVM log with large data", 8 << 10, true},
		{"HTTP trigger request", 32 << 10, true},
	}
	for _, k := range kinds {
		blobs, perBlob, element := blobsNeeded(k.payload, ocr3_1types.MaxMaxBlobPayloadBytes)
		t.Logf("%-24s payload %6dB -> %6dB per event in a blob: %5d events per %dMiB blob, %2d blobs for %d events (%d handles available)",
			k.name, k.payload, element, perBlob, ocr3_1types.MaxMaxBlobPayloadBytes>>20, blobs, maxObservedEvents, maxObservedBlobHandles)
		if k.typical {
			assert.LessOrEqual(t, blobs, maxObservedBlobHandles, "%s: maxObservedEvents must be broadcastable in maxObservedBlobHandles blobs", k.name)
		}
	}

	// the largest payload for which a full observation (maxObservedEvents) still fits in maxObservedBlobHandles blobs
	limit := sort.Search(ocr3_1types.MaxMaxBlobPayloadBytes, func(p int) bool {
		blobs, _, _ := blobsNeeded(p, ocr3_1types.MaxMaxBlobPayloadBytes)
		return blobs == 0 || blobs > maxObservedBlobHandles
	}) - 1
	blobs, perBlob, element := blobsNeeded(limit, ocr3_1types.MaxMaxBlobPayloadBytes)
	t.Logf("%d events fit in %d blobs up to a payload of %dB (%dB per event, %d per blob, %d blobs)",
		maxObservedEvents, maxObservedBlobHandles, limit, element, perBlob, blobs)
	assert.Positive(t, limit, "a full observation must be able to carry some payload")
	assert.LessOrEqual(t, blobs, maxObservedBlobHandles)
	// beyond that, Observation drops the events that do not fit, rather than failing
}

func TestCalibration_KVState(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(2))

	// consensus and queued events carry one blob ref each
	queuedAvg, queuedMax := summarize(1000, func(i int) int {
		e := observedEventFor(rng, i)
		e.Blobs = []*pb.CoordinatedEventBlobHandle{{Blob: uint32(rng.Intn(maxObservedBlobHandles * 31)), Element: uint32(rng.Intn(maxObservedEvents))}}
		return protowire.SizeBytes(proto.Size(e)) + protowire.SizeTag(1) // as an element of ConsensusQueueEventList
	})
	t.Logf("queued/consensus event in KV: avg %.0fB, max %dB", queuedAvg, queuedMax)
	queuedCapacity := maxKeyValueValueBytes / queuedMax
	t.Logf("the %dMiB value limit holds %d queued events (maxObservedEvents is %d)", maxKeyValueValueBytes>>20, queuedCapacity, maxObservedEvents)
	assert.GreaterOrEqual(t, queuedCapacity, 2*maxObservedEvents, "the queue must hold at least two full observation rounds' worth of events")

	for _, n := range []int{4, 7, 16, 31} {
		f := (n - 1) / 3

		// an observed event waits at between f+1 and 2f votes, each with a node ID and two timestamps
		observedAvg, observedMax := summarize(1000, func(i int) int {
			e := &pb.ConsensusQueueObservedEvent{WorkflowID: hexID("wf", i), ExecutionID: hexID("exec", i)}
			for v := range 2 * f {
				e.NodeIds = append(e.NodeIds, byte(v))
				e.ObservedAts = append(e.ObservedAts, randTS(rng))
				e.Deadlines = append(e.Deadlines, randTS(rng))
			}
			for range 2 * f { // the blob of each vote
				e.Blobs = append(e.Blobs, &pb.CoordinatedEventBlobHandle{Blob: uint32(rng.Intn(maxObservedBlobHandles * n)), Element: uint32(rng.Intn(maxObservedEvents))})
			}
			entry := &pb.ConsensusQueueObservedEventEntry{TimelessHash: make([]byte, sha256.Size), Event: e}
			return protowire.SizeBytes(proto.Size(entry)) + protowire.SizeTag(1)
		})
		observedCapacity := maxKeyValueValueBytes / observedMax
		t.Logf("N=%-2d observed event at 2f=%d votes in KV: avg %.0fB, max %dB -> %d fit in the value limit", n, 2*f, observedAvg, observedMax, observedCapacity)
		assert.GreaterOrEqual(t, observedCapacity, maxObservedEvents, "N=%d: a full round of events must fit in the observed set", n)

		// each of N nodes can broadcast maxObservedBlobHandles blobs in a round; their handles are kept while events reference them
		handles := maxObservedBlobHandles * n * ocr3_1types.BlobHandleMarshalledBytesUpperBound(n, f)
		t.Logf("N=%-2d one round of blob handles from every node: %dB (%.0f%% of the value limit)", n, handles, 100*float64(handles)/maxKeyValueValueBytes)
		assert.LessOrEqual(t, handles, maxKeyValueValueBytes, "N=%d: a round's blob handles must fit in the blob handles value", n)
	}

	// All four keys can be full at once: observed, queued and consensus events, and blob handles.
	total := numKVKeys*maxKeyValueValueBytes + numKVKeys*len(consensusEventsKey) + len(observedEventsKey)
	t.Logf("all %d KV values full: %dMiB of libocr's %dMiB modified-keys-plus-values limit", numKVKeys, total>>20, ocr3_1types.MaxMaxKeyValueModifiedKeysPlusValuesBytes>>20)
	assert.LessOrEqual(t, total, ocr3_1types.MaxMaxKeyValueModifiedKeysPlusValuesBytes)
	assert.LessOrEqual(t, maxKeyValueValueBytes, ocr3_1types.MaxMaxKeyValueValueBytes)
	assert.LessOrEqual(t, numKVKeys, ocr3_1types.MaxMaxKeyValueModifiedKeys)
}

// Blobs are kept until their events expire (the workflow's timeout plus the grace period), so the blob quotas are set by
// how long that is in rounds. This logs what a sustained load needs, using the real defaults.
func TestCalibration_BlobLifetime(t *testing.T) {
	t.Parallel()

	timeout := cresettings.Default.PerWorkflow.TriggerEventQueueTimeout.DefaultValue
	grace := cresettings.Default.ConsensusQueueExpirationGraceInterval.DefaultValue
	t.Logf("event timeout %s, grace %s", timeout, grace)
	assert.LessOrEqual(t, grace*10, timeout, "the grace period is slack for clock skew and consensus latency, a small fraction of the timeout")

	for _, roundDur := range []time.Duration{250 * time.Millisecond, 500 * time.Millisecond, time.Second, 2 * time.Second} {
		rounds := blobExpirySeqNr(0, timeout+grace, roundDur)
		perRoundBlobs := maxObservedBlobHandles // a full observation every round: the worst case
		t.Logf("DeltaRound %-6s: a blob is kept for %5d rounds (%s); sustained full observations would hold %6d blobs = %5.1f GiB per oracle",
			roundDur, rounds, time.Duration(rounds)*roundDur, perRoundBlobs*int(rounds),
			float64(perRoundBlobs)*float64(rounds)*float64(ocr3_1types.MaxMaxBlobPayloadBytes)/(1<<30))
		assert.GreaterOrEqual(t, time.Duration(rounds)*roundDur, timeout+grace, "blobs must outlive their events")
	}
}
