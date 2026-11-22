package queue

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
)

func Test_extractEvent(t *testing.T) {
	payload, err := anypb.New(timestamppb.New(time.Now()))
	require.NoError(t, err)
	event := &pb.CoordinatedEvent{
		WorkflowID:   "want",
		TriggerCapID: "capID",
		TriggerIndex: 1,
		ObservedAt:   &timestamppb.Timestamp{Seconds: 10, Nanos: 42},
		Deadline:     timestamppb.New(time.Now()),
		Event: &pb.TriggerResponse{
			Event: &pb.TriggerEvent{
				TriggerType: "type",
				Id:          "id",
				Outputs:     nil,
				Payload:     payload,
			},
			Error: "asdf",
		},
	}
	otherPayload, err := anypb.New(&pb.CapabilityResponse{})
	require.NoError(t, err)
	others := &pb.CoordinatedEvent{
		WorkflowID:   "foo",
		TriggerCapID: "bar",
		TriggerIndex: 99,
		ObservedAt:   &timestamppb.Timestamp{Nanos: 99},
		Deadline:     &timestamppb.Timestamp{Seconds: 1000},
		Event: &pb.TriggerResponse{
			Event: &pb.TriggerEvent{
				TriggerType: "baz",
				Id:          "trigger-id",
				Outputs:     nil,
				Payload:     otherPayload,
			},
			Error: "none",
		},
	}
	for _, tt := range []struct {
		name     string
		len, idx int
	}{
		{"first", 4, 0},
		{"last", 4, 3},
		{"middle", 4, 2},
		{"bounds", 4, 4},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var events []*pb.CoordinatedEvent
			for i := range tt.len {
				if i == tt.idx {
					events = append(events, event)
				} else {
					events = append(events, others)
				}
			}
			b, err := proto.MarshalOptions{Deterministic: true}.Marshal(&pb.CoordinatedEventList{Events: events})
			require.NoError(t, err)
			got, _, err := extractEvent(b, tt.idx)
			if tt.idx < tt.len {
				require.True(t, proto.Equal(event, got))
			} else {
				require.Error(t, err)
			}
		})
	}
}

func Test_coordinatedEventsFieldNumber(t *testing.T) {
	fd := (&pb.CoordinatedEventList{}).ProtoReflect().Descriptor().Fields().ByName("events")
	require.NotNil(t, fd)
	require.Equal(t, protowire.Number(coordinatedEventsFieldNumber), fd.Number())
}

var errBoom = errors.New("boom")

// handleBytes returns the wire encoding of a minimal valid blob handle; distinct b gives distinct handles.
func handleBytes(b byte) []byte {
	root := make([]byte, 32)
	root[0] = b
	return append([]byte{1, 0x0a, 0x20}, root...)
}

func newBlobCache(f ocr3_1types.BlobFetcher) *blobCache {
	return &blobCache{fetcher: f, cache: make(map[string]func() ([]byte, error))}
}

// countingFetcher counts FetchBlob calls and returns a fixed payload/error.
type countingFetcher struct {
	mu      sync.Mutex
	calls   int
	payload []byte
	err     error
	gate    chan struct{} // if set, FetchBlob blocks until closed
}

func (f *countingFetcher) FetchBlob(context.Context, ocr3_1types.BlobHandle) ([]byte, error) {
	f.mu.Lock()
	f.calls++
	f.mu.Unlock()
	if f.gate != nil {
		<-f.gate
	}
	return f.payload, f.err
}

func (f *countingFetcher) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.calls
}

func marshalEvents(t testing.TB, events ...*pb.CoordinatedEvent) []byte {
	t.Helper()
	b, err := proto.MarshalOptions{Deterministic: true}.Marshal(&pb.CoordinatedEventList{Events: events})
	require.NoError(t, err)
	return b
}

func TestBlobCache_FetchBlob(t *testing.T) {
	t.Parallel()

	t.Run("concurrent fetches of one handle hit the fetcher once", func(t *testing.T) {
		t.Parallel()
		f := &countingFetcher{payload: []byte("data"), gate: make(chan struct{})}
		c := newBlobCache(f)

		var wg sync.WaitGroup
		results := make([][]byte, 20)
		for i := range results {
			wg.Add(1)
			go func() {
				defer wg.Done()
				b, err := c.fetchBlob(t.Context(), handleBytes(1))
				assert.NoError(t, err)
				results[i] = b
			}()
		}
		close(f.gate)
		wg.Wait()

		assert.Equal(t, 1, f.count())
		for _, r := range results {
			assert.Equal(t, []byte("data"), r)
		}
	})

	t.Run("distinct handles are fetched separately", func(t *testing.T) {
		t.Parallel()
		f := &countingFetcher{payload: []byte("data")}
		c := newBlobCache(f)
		for _, h := range [][]byte{handleBytes(1), handleBytes(2), handleBytes(1), handleBytes(2)} {
			_, err := c.fetchBlob(t.Context(), h)
			require.NoError(t, err)
		}
		assert.Equal(t, 2, f.count())
	})

	t.Run("invalid handle errors without fetching", func(t *testing.T) {
		t.Parallel()
		f := &countingFetcher{}
		_, err := newBlobCache(f).fetchBlob(t.Context(), []byte("not a handle"))
		require.Error(t, err)
		assert.Zero(t, f.count())
	})

	t.Run("fetch error is returned and cached for the round", func(t *testing.T) {
		t.Parallel()
		f := &countingFetcher{err: errBoom}
		c := newBlobCache(f)
		for range 3 {
			_, err := c.fetchBlob(t.Context(), handleBytes(1))
			require.ErrorIs(t, err, errBoom)
		}
		assert.Equal(t, 1, f.count())
	})
}

func TestFetchCoordinatedEventFromBlob(t *testing.T) {
	t.Parallel()

	_, e0 := newCoordinatedEvent(t, "wf", "e0")
	_, e1 := newCoordinatedEvent(t, "wf", "e1")
	e0.ObservedAt, e0.Deadline, e1.ObservedAt, e1.Deadline = nil, nil, nil, nil
	payload := marshalEvents(t, e0, e1)

	t.Run("returns the element and its raw bytes", func(t *testing.T) {
		t.Parallel()
		got, raw, err := fetchCoordinatedEventFromBlob(t.Context(), newBlobCache(&countingFetcher{payload: payload}), handleBytes(1), 1)
		require.NoError(t, err)
		assert.True(t, proto.Equal(e1, got))
		want, err := proto.MarshalOptions{Deterministic: true}.Marshal(e1)
		require.NoError(t, err)
		assert.Equal(t, want, raw, "raw bytes must equal the deterministic marshaling the hash was computed over")
	})

	t.Run("fetch error", func(t *testing.T) {
		t.Parallel()
		_, _, err := fetchCoordinatedEventFromBlob(t.Context(), newBlobCache(&countingFetcher{err: errBoom}), handleBytes(1), 0)
		require.ErrorIs(t, err, errBoom)
	})

	t.Run("index out of range", func(t *testing.T) {
		t.Parallel()
		_, _, err := fetchCoordinatedEventFromBlob(t.Context(), newBlobCache(&countingFetcher{payload: payload}), handleBytes(1), 2)
		require.Error(t, err)
	})
}

func Test_extractEvent_edgeCases(t *testing.T) {
	t.Parallel()

	_, e := newCoordinatedEvent(t, "wf", "e")
	e.ObservedAt, e.Deadline = nil, nil
	eb, err := proto.MarshalOptions{Deterministic: true}.Marshal(e)
	require.NoError(t, err)

	t.Run("empty and nil blob", func(t *testing.T) {
		t.Parallel()
		_, _, err := extractEvent(nil, 0)
		require.Error(t, err)
		_, _, err = extractEvent([]byte{}, 0)
		require.Error(t, err)
	})

	t.Run("negative index", func(t *testing.T) {
		t.Parallel()
		_, _, err := extractEvent(marshalEvents(t, e), -1)
		require.Error(t, err)
	})

	t.Run("unknown fields between events are skipped", func(t *testing.T) {
		t.Parallel()
		var b []byte
		b = protowire.AppendTag(b, 2, protowire.VarintType) // unknown field 2
		b = protowire.AppendVarint(b, 12345)
		b = protowire.AppendTag(b, coordinatedEventsFieldNumber, protowire.BytesType)
		b = protowire.AppendBytes(b, eb)
		b = protowire.AppendTag(b, 3, protowire.BytesType) // unknown field 3
		b = protowire.AppendBytes(b, []byte("junk"))
		got, raw, err := extractEvent(b, 0)
		require.NoError(t, err)
		assert.True(t, proto.Equal(e, got))
		assert.Equal(t, eb, raw)
		_, _, err = extractEvent(b, 1)
		require.Error(t, err, "unknown fields are not counted as elements")
	})

	t.Run("wrong wire type for events field", func(t *testing.T) {
		t.Parallel()
		b := protowire.AppendTag(nil, coordinatedEventsFieldNumber, protowire.VarintType)
		b = protowire.AppendVarint(b, 1)
		_, _, err := extractEvent(b, 0)
		require.Error(t, err)
	})

	t.Run("truncated element", func(t *testing.T) {
		t.Parallel()
		b := marshalEvents(t, e)
		_, _, err := extractEvent(b[:len(b)-3], 0)
		require.Error(t, err)
	})

	t.Run("truncated tag", func(t *testing.T) {
		t.Parallel()
		_, _, err := extractEvent([]byte{0xff}, 0)
		require.Error(t, err)
	})

	t.Run("element that is not a valid event returns an error and its bytes", func(t *testing.T) {
		t.Parallel()
		b := protowire.AppendTag(nil, coordinatedEventsFieldNumber, protowire.BytesType)
		b = protowire.AppendBytes(b, []byte{0xff, 0xff, 0xff})
		got, raw, err := extractEvent(b, 0)
		require.Error(t, err)
		assert.Nil(t, got)
		assert.NotNil(t, raw)
	})
}

// hintBlobs records expiration hints and can be told to fail.
type hintBlobs struct {
	memBlobs
	hints []ocr3_1types.BlobExpirationHint
	err   error
}

func (h *hintBlobs) BroadcastBlob(ctx context.Context, payload []byte, hint ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	if h.err != nil {
		return ocr3_1types.BlobHandle{}, h.err
	}
	h.hints = append(h.hints, hint)
	return h.memBlobs.BroadcastBlob(ctx, payload, hint)
}

func eventsWithDeadline(t testing.TB, n int, deadline time.Time) []*pb.CoordinatedEvent {
	t.Helper()
	var events []*pb.CoordinatedEvent
	for i := range n {
		_, e := newCoordinatedEvent(t, "wf", string(rune('a'+i)))
		e.Deadline = timestamppb.New(deadline)
		events = append(events, e)
	}
	return events
}

func TestBroadcastAsBlob(t *testing.T) {
	t.Parallel()

	t.Run("all events fit", func(t *testing.T) {
		t.Parallel()
		blobs := &hintBlobs{}
		c := testPlugin(t, sliceQueue{}, blobs, ocr3_1types.MaxMaxObservationBytes)
		events := eventsWithDeadline(t, 3, time.Now().Add(10*time.Second))

		handle, n, err := c.broadcastAsBlob(t.Context(), blobs, events, 100)
		require.NoError(t, err)
		assert.Equal(t, 3, n)
		assert.NotEmpty(t, handle)
		require.Len(t, blobs.payloads, 1)
		got, _, err := extractEvent(blobs.payloads[0], 2)
		require.NoError(t, err)
		assert.True(t, proto.Equal(events[2], got))
	})

	t.Run("trims to the payload limit", func(t *testing.T) {
		t.Parallel()
		blobs := &hintBlobs{}
		c := testPlugin(t, sliceQueue{}, blobs, ocr3_1types.MaxMaxObservationBytes)
		events := eventsWithDeadline(t, 4, time.Now().Add(10*time.Second))
		c.limits.MaxBlobPayloadBytes = proto.Size(&pb.CoordinatedEventList{Events: events[:2]})

		_, n, err := c.broadcastAsBlob(t.Context(), blobs, events, 100)
		require.NoError(t, err)
		assert.Equal(t, 2, n)
		assert.LessOrEqual(t, len(blobs.payloads[0]), c.limits.MaxBlobPayloadBytes)
	})

	t.Run("nothing fits", func(t *testing.T) {
		t.Parallel()
		blobs := &hintBlobs{}
		c := testPlugin(t, sliceQueue{}, blobs, ocr3_1types.MaxMaxObservationBytes)
		c.limits.MaxBlobPayloadBytes = 1
		_, n, err := c.broadcastAsBlob(t.Context(), blobs, eventsWithDeadline(t, 2, time.Now().Add(time.Minute)), 100)
		require.Error(t, err)
		assert.Zero(t, n)
		assert.Empty(t, blobs.payloads)
	})

	t.Run("no events", func(t *testing.T) {
		t.Parallel()
		blobs := &hintBlobs{}
		c := testPlugin(t, sliceQueue{}, blobs, ocr3_1types.MaxMaxObservationBytes)
		_, _, err := c.broadcastAsBlob(t.Context(), blobs, nil, 100)
		require.Error(t, err)
	})

	t.Run("broadcast error is wrapped", func(t *testing.T) {
		t.Parallel()
		blobs := &hintBlobs{err: errBoom}
		c := testPlugin(t, sliceQueue{}, blobs, ocr3_1types.MaxMaxObservationBytes)
		_, _, err := c.broadcastAsBlob(t.Context(), blobs, eventsWithDeadline(t, 1, time.Now().Add(time.Minute)), 100)
		require.ErrorIs(t, err, errBoom)
	})

	t.Run("expiration hint follows the latest deadline", func(t *testing.T) {
		t.Parallel()
		blobs := &hintBlobs{}
		c := testPlugin(t, sliceQueue{}, blobs, ocr3_1types.MaxMaxObservationBytes) // 1s rounds
		events := eventsWithDeadline(t, 2, time.Now().Add(5*time.Second))
		events[0].Deadline = timestamppb.New(time.Now().Add(60 * time.Second)) // the latest

		_, _, err := c.broadcastAsBlob(t.Context(), blobs, events, 100)
		require.NoError(t, err)
		hint := blobs.hints[0].(ocr3_1types.BlobExpirationHintSequenceNumber)
		// 60s to the latest deadline + 60s grace + 30s minimum lifetime, at 1s per round
		assert.InDelta(t, 100+150, int(hint.SeqNr), 1)
	})
}

// Known issues in broadcastAsBlob's expiration math (SHOULD FIX): these fail until the hint is clamped.
func TestBroadcastAsBlob_ExpirationHintIsClamped(t *testing.T) {
	t.Parallel()

	t.Run("already expired deadline does not wrap the sequence number", func(t *testing.T) {
		t.Parallel()
		blobs := &hintBlobs{}
		c := testPlugin(t, sliceQueue{}, blobs, ocr3_1types.MaxMaxObservationBytes)
		events := eventsWithDeadline(t, 1, time.Now().Add(-time.Hour))

		_, _, err := c.broadcastAsBlob(t.Context(), blobs, events, 100)
		require.NoError(t, err)
		hint := blobs.hints[0].(ocr3_1types.BlobExpirationHintSequenceNumber)
		assert.GreaterOrEqual(t, hint.SeqNr, uint64(100))
		assert.Less(t, hint.SeqNr, uint64(1_000_000), "negative duration must not convert to a huge uint64")
	})

	t.Run("zero round interval does not panic", func(t *testing.T) {
		t.Parallel()
		blobs := &hintBlobs{}
		c := testPlugin(t, sliceQueue{}, blobs, ocr3_1types.MaxMaxObservationBytes)
		c.onchainCfg.EstimatedRoundInterval = 0
		require.NotPanics(t, func() {
			_, _, _ = c.broadcastAsBlob(t.Context(), blobs, eventsWithDeadline(t, 1, time.Now().Add(time.Minute)), 100)
		})
	})
}

// maxEventFramingOverhead must cover the framing of an event of any size a blob can hold, including at length prefix boundaries.
func TestMaxEventFramingOverhead(t *testing.T) {
	t.Parallel()

	for _, size := range []int{0, 1, 127, 128, 129, 1<<14 - 1, 1 << 14, 1<<21 - 1, 1 << 21, ocr3_1types.MaxMaxBlobPayloadBytes - maxEventFramingOverhead} {
		framed := protowire.SizeTag(coordinatedEventsFieldNumber) + protowire.SizeBytes(size)
		assert.LessOrEqual(t, framed-size, maxEventFramingOverhead, "event of %d bytes", size)
	}
	assert.Equal(t, 5, maxEventFramingOverhead, "1 tag byte plus a 4 byte length for blobs up to libocr's 5MiB maximum")
}

func TestBlobExpirySeqNr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		remaining time.Duration
		roundDur  time.Duration
		want      uint64 // added to seqNr
	}{
		{"remaining time and the minimum lifetime, in rounds", 60 * time.Second, time.Second, 90},
		{"partial rounds round up", 60*time.Second + time.Millisecond, time.Second, 91},
		{"slower rounds need fewer", 60 * time.Second, 10 * time.Second, 9},
		{"already expired still gets the minimum lifetime", -time.Hour, time.Second, 30},
		{"minimum rounds when rounds are long", 0, time.Hour, minBlobLifetimeRounds},
		{"unknown round interval gets the minimum rounds", time.Hour, 0, minBlobLifetimeRounds},
		{"negative round interval gets the minimum rounds", time.Hour, -time.Second, minBlobLifetimeRounds},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, 1000+tt.want, blobExpirySeqNr(1000, tt.remaining, tt.roundDur))
		})
	}

	t.Run("never expires before the events in it do", func(t *testing.T) {
		t.Parallel()
		// rounds are never faster than the estimate, so the sequence number the blob expires at is always reached
		// after the remaining time has passed
		for _, remaining := range []time.Duration{0, time.Second, 90 * time.Second, 11 * time.Minute, time.Hour} {
			for _, roundDur := range []time.Duration{100 * time.Millisecond, time.Second, 3 * time.Second} {
				rounds := blobExpirySeqNr(0, remaining, roundDur)
				assert.GreaterOrEqual(t, time.Duration(rounds)*roundDur, remaining+minBlobLifetime, "%s remaining at %s per round", remaining, roundDur)
			}
		}
	})
}

func TestBroadcastAsBlob_GraceIsOptional(t *testing.T) {
	t.Parallel()

	blobs := &hintBlobs{}
	c := testPlugin(t, sliceQueue{}, blobs, ocr3_1types.MaxMaxObservationBytes) // 1s rounds
	events := eventsWithDeadline(t, 1, time.Now().Add(60*time.Second))

	for _, tt := range []struct {
		name  string
		grace time.Duration
		err   error
		want  int // sequence number offset: 60s + grace + 30s at 1s per round
	}{
		{"no grace", 0, nil, 90},
		{"grace", 2 * time.Minute, nil, 210},
		{"error still uses the returned default", 2 * time.Minute, errBoom, 210},
		{"negative grace is none", -time.Hour, nil, 90},
	} {
		c.expirationDur = func(context.Context) (time.Duration, error) { return tt.grace, tt.err }
		blobs.hints = nil
		_, _, err := c.broadcastAsBlob(t.Context(), blobs, events, 0)
		require.NoError(t, err, tt.name)
		assert.InDelta(t, tt.want, int(blobs.hints[0].(ocr3_1types.BlobExpirationHintSequenceNumber).SeqNr), 1, tt.name)
	}
}

// broadcastAsBlob (used by tests: Observation divides events between blobs up front) broadcasts a blob with as many coordinated events as possible, returning the blob handle and count of events included.
func (c *consensusQueuePlugin) broadcastAsBlob(ctx context.Context, broadcaster ocr3_1types.BlobBroadcaster, events []*pb.CoordinatedEvent,
	seqNr uint64) (blobHandle []byte, n int, err error) {
	n = c.blobPrefix(events)
	if n == 0 {
		return nil, 0, errors.New("failed to package any events in blob")
	}

	blobHandle, err = c.broadcastBlob(ctx, broadcaster, events[:n], seqNr)
	if err != nil {
		return nil, -1, err
	}

	return blobHandle, n, nil
}
