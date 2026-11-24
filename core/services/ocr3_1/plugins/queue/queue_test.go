package queue

import (
	"bytes"
	"cmp"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"iter"
	"maps"
	"math"
	"math/rand"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric/noop"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
)

// libocr calls NewReportingPlugin with no recover, so a panic here takes down
// the whole node. This guards the constructor and the plugin lifecycle.
func TestConsensusQueuePluginFactory_NewReportingPlugin(t *testing.T) {
	t.Parallel()

	cfg := &PluginConfig{Limits: ocr3_1types.ReportingPluginLimits{MaxKeyValueModifiedKeys: 1}}
	factory, err := NewConsensusQueuePluginFactory(logger.Test(t), cfg, limits.Factory{}, noopCoordinator{}, noopQueue{}, noop.NewMeterProvider().Meter("test"))
	require.NoError(t, err)

	plugin, info, err := factory.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{N: 4, F: 1}, nil)
	require.NoError(t, err)
	require.NotNil(t, plugin)

	info1, ok := info.(ocr3_1types.ReportingPluginInfo1)
	require.True(t, ok, "unexpected ReportingPluginInfo type %T", info)
	require.Equal(t, "ConsensusQueueReportingPlugin", info1.Name)
	require.Equal(t, cfg.Limits, info1.Limits)

	require.NoError(t, plugin.Close())
}

type noopCoordinator struct{}

func (n noopCoordinator) Start(ctx context.Context) error {
	return nil
}

func (n noopCoordinator) Close() error {
	return nil
}

func (n noopCoordinator) Ready() error {
	return nil
}

func (n noopCoordinator) HealthReport() map[string]error {
	return map[string]error{}
}

func (n noopCoordinator) Name() string {
	return "NoopTriggerCoordinator"
}

func (n noopCoordinator) Ack(ctx context.Context, triggerCapID, triggerRegistrationID, eventID string) error {
	return nil
}

func (n noopCoordinator) RegisterTriggers(ctx context.Context, subscriber triggers.Subscriber, params triggers.RegistrationParams) ([]string, error) {
	return nil, nil
}

func (n noopCoordinator) UnregisterTriggers(ctx context.Context, workflowID string) error {
	return nil
}

func (n noopCoordinator) AvailableCapacity(ctx context.Context) (triggers.Capacity, error) {
	return triggers.Capacity{}, nil
}

func (n noopCoordinator) StartExecution(ctx context.Context, event triggers.CoordinatedEvent) error {
	return nil
}

func (n noopCoordinator) WorkflowOrg(string) (string, bool) { return "", false }

type noopQueue struct{}

func (n noopQueue) Put(ctx context.Context, event triggers.CoordinatedEvent) error {
	return nil
}

func (n noopQueue) Drain(ctx context.Context) ([]triggers.CoordinatedEvent, error) {
	return nil, nil
}

func (n noopQueue) Stats(ctx context.Context) (triggers.Stats, error) {
	return triggers.Stats{}, nil
}

func (n noopQueue) Sorted(context.Context, int, int) iter.Seq2[[]byte, *pb.CoordinatedEvent] {
	return nil
}

func (n noopQueue) Remove(ctx context.Context, execIDs iter.Seq[string]) {}

func newCoordinatedEvent(t testing.TB, wf, eventID string) (hash []byte, event *pb.CoordinatedEvent) {
	t.Helper()
	event = &pb.CoordinatedEvent{
		WorkflowID: wf, TriggerCapID: "trigger@1.0.0",
		Event: &pb.TriggerResponse{Event: &pb.TriggerEvent{Id: eventID}},
	}
	b, err := event.MarshalTimeless() // hash is over the event without timestamps
	require.NoError(t, err)
	sum := sha256.Sum256(b)
	hash = sum[:]
	event = proto.CloneOf(event)
	event.ObservedAt = ts(1000)
	event.Deadline = ts(2000)
	return hash, event
}

// testWF is a valid workflow ID (32 bytes of hex), as events must have.
var testWF = hexID("wf", 0)

type hashedEvent struct {
	hash  []byte
	event *pb.CoordinatedEvent
}

// sliceQueue serves a fixed set of events from Sorted.
type sliceQueue struct {
	noopQueue
	events []hashedEvent
}

func (q sliceQueue) Sorted(_ context.Context, limit, _ int) iter.Seq2[[]byte, *pb.CoordinatedEvent] {
	sorted := slices.Clone(q.events)
	slices.SortFunc(sorted, func(a, b hashedEvent) int { // same order as triggers.CentralQueue
		return cmp.Or(a.event.Deadline.AsTime().Compare(b.event.Deadline.AsTime()), bytes.Compare(a.hash, b.hash))
	})
	return func(yield func([]byte, *pb.CoordinatedEvent) bool) {
		for i, e := range sorted {
			if i >= limit || !yield(e.hash, proto.CloneOf(e.event)) {
				return
			}
		}
	}
}

// testBlobHandle returns a valid handle (libocr offers no constructor) built from its wire encoding. i identifies it.
func testBlobHandle(i int) ocr3_1types.BlobHandle {
	var h ocr3_1types.BlobHandle
	if err := h.UnmarshalBinary(handleBytes(byte(i))); err != nil { //nolint:gosec // G115: a small test value
		panic(err)
	}
	return h
}

func marshalHandle(h ocr3_1types.BlobHandle) string {
	b, err := h.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return string(b)
}

// memBlobs is an in-memory blob store. Each broadcast gets a distinct handle (the first is handleBytes(1), the
// second handleBytes(2), ...), and fetches are served by handle. payloads are in broadcast order.
type memBlobs struct {
	payloads [][]byte
	byHandle map[string][]byte
}

// memBlobsMu guards every memBlobs, which are used concurrently when blobs are broadcast or fetched in parallel.
var memBlobsMu sync.Mutex

func (m *memBlobs) BroadcastBlob(_ context.Context, payload []byte, _ ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	memBlobsMu.Lock()
	defer memBlobsMu.Unlock()
	m.payloads = append(m.payloads, payload)
	h := testBlobHandle(len(m.payloads))
	m.put(h, payload)
	return h, nil
}

func (m *memBlobs) put(h ocr3_1types.BlobHandle, payload []byte) {
	if m.byHandle == nil {
		m.byHandle = make(map[string][]byte)
	}
	m.byHandle[marshalHandle(h)] = payload
}

// store registers a payload under a handle that was written to KV state rather than broadcast.
func (m *memBlobs) store(handle []byte, payload []byte) {
	memBlobsMu.Lock()
	defer memBlobsMu.Unlock()
	var h ocr3_1types.BlobHandle
	if err := h.UnmarshalBinary(handle); err != nil {
		panic(err)
	}
	m.payloads = append(m.payloads, payload)
	m.put(h, payload)
}

func (m *memBlobs) FetchBlob(_ context.Context, h ocr3_1types.BlobHandle) ([]byte, error) {
	memBlobsMu.Lock()
	defer memBlobsMu.Unlock()
	payload, ok := m.byHandle[marshalHandle(h)]
	if !ok {
		return nil, errors.New("unknown blob handle")
	}
	return payload, nil
}

func testPlugin(t testing.TB, q sliceQueue, blobs ocr3_1types.BlobBroadcastFetcher, maxObsBytes int) *consensusQueuePlugin {
	t.Helper()
	limits := ocr3_1types.ReportingPluginLimits{
		MaxObservationBytes:                    maxObsBytes,
		MaxBlobPayloadBytes:                    ocr3_1types.MaxMaxBlobPayloadBytes,
		MaxKeyValueModifiedKeysPlusValuesBytes: ocr3_1types.MaxMaxKeyValueModifiedKeysPlusValuesBytes,
	}
	return newConsensusQueuePlugin(logger.Test(t),
		ocr3types.ReportingPluginConfig{N: 4, F: 1, OracleID: 0, EstimatedRoundInterval: time.Second},
		limits, noopCoordinator{}, q,
		func(context.Context) (time.Duration, error) { return time.Minute, nil }, blobs, testMetrics(t))
}

// testMetrics returns metrics that record nowhere.
func testMetrics(tb testing.TB) *pluginMetrics {
	tb.Helper()
	m, err := newPluginMetrics(noop.NewMeterProvider().Meter("test"), "digest")
	require.NoError(tb, err)
	return m
}

func TestObserveEvent(t *testing.T) {
	t.Parallel()

	hash, event := newCoordinatedEvent(t, "wf1", "evt-1")
	obs := &pb.ConsensusQueueObservation{WorkflowCapacities: map[string]uint32{}}

	newWF, err := observeEvent(obs, hash, event, map[string]int{"wf1": 3})
	require.NoError(t, err)
	assert.True(t, newWF)
	assert.Equal(t, uint32(3), obs.WorkflowCapacities["wf1"])
	require.Len(t, obs.Events, 1)

	got := obs.Events[0]
	wantExec, err := workflows.GenerateExecutionIDWithTriggerIndex("wf1", "evt-1", 0)
	require.NoError(t, err)
	assert.Equal(t, wantExec, got.ExecutionID)
	assert.Equal(t, hash, got.TimelessHash)
	assert.Equal(t, int64(1000), got.ObservedAt.Seconds)
	assert.Equal(t, int64(2000), got.Deadline.Seconds)
	assert.Nil(t, event.ObservedAt, "timestamps are stripped from the event so it hashes consistently")
	assert.Nil(t, event.Deadline)

	// second event in the same workflow doesn't re-add capacity
	hash2, event2 := newCoordinatedEvent(t, "wf1", "evt-2")
	newWF, err = observeEvent(obs, hash2, event2, map[string]int{"wf1": 99})
	require.NoError(t, err)
	assert.False(t, newWF)
	assert.Equal(t, uint32(3), obs.WorkflowCapacities["wf1"])
	assert.Equal(t, uint32(1), obs.Events[1].Blobs[0].Element, "element index follows position")
}

func TestObservationsAggregate(t *testing.T) {
	t.Parallel()
	const f = 1

	mkAO := func(id commontypes.OracleID, obs *pb.ConsensusQueueObservation) ocrtypes.AttributedObservation {
		b, err := proto.Marshal(obs)
		require.NoError(t, err)
		return ocrtypes.AttributedObservation{Observer: id, Observation: b}
	}
	base := func(capacity uint32, wfCaps map[string]uint32) *pb.ConsensusQueueObservation {
		return &pb.ConsensusQueueObservation{Expiration: ts(100), Capacity: capacity, WorkflowCapacities: wfCaps}
	}

	t.Run("capacity and expiration without events", func(t *testing.T) {
		t.Parallel()
		var res observations
		res.aggregate(t.Context(), f, []ocrtypes.AttributedObservation{
			mkAO(0, base(10, map[string]uint32{testWF: 5})),
			mkAO(1, base(1, map[string]uint32{testWF: 1})),
			mkAO(2, base(7, map[string]uint32{testWF: 9})),
		})
		require.NoError(t, res.err)
		assert.Equal(t, uint32(7), res.fPlus1thCap, "F+1th lowest of 1,7,10")
		assert.Equal(t, uint32(5), res.fPlus1thCapPerWF[testWF], "F+1th lowest of 1,5,9")
		assert.Equal(t, time.Unix(100, 0).UTC(), res.medianExpiration.UTC())
	})

	t.Run("a byzantine low or high capacity cannot move the result outside the honest range", func(t *testing.T) {
		t.Parallel()
		var res observations
		res.aggregate(t.Context(), f, []ocrtypes.AttributedObservation{
			mkAO(0, base(0, nil)),     // byzantine low
			mkAO(1, base(5, nil)),     // honest
			mkAO(2, base(6, nil)),     // honest
			mkAO(3, base(1<<31, nil)), // byzantine high
		})
		assert.Equal(t, uint32(5), res.fPlus1thCap)
	})

	t.Run("garbage observation errors", func(t *testing.T) {
		t.Parallel()
		var res observations
		res.aggregate(t.Context(), f, []ocrtypes.AttributedObservation{{Observer: 0, Observation: []byte{0xff, 0xff}}})
		require.Error(t, res.err)
	})

	t.Run("events are merged by hash across observers", func(t *testing.T) {
		t.Parallel()
		hash, event := newCoordinatedEvent(t, testWF, "evt")
		mkObs := func() *pb.ConsensusQueueObservation {
			o := base(5, map[string]uint32{testWF: 5})
			o.BlobHandles = [][]byte{[]byte("handle")}
			_, err := observeEvent(o, hash, proto.CloneOf(event), map[string]int{testWF: 5})
			require.NoError(t, err)
			return o
		}
		var res observations
		require.NotPanics(t, func() {
			res.aggregate(t.Context(), f, []ocrtypes.AttributedObservation{mkAO(0, mkObs()), mkAO(1, mkObs()), mkAO(3, mkObs())})
		})
		require.NoError(t, res.err)
		require.Len(t, res.newObservedEvents, 1)
		got := res.newObservedEvents[string(hash)]
		require.NotNil(t, got)
		assert.Equal(t, []byte{0, 1, 3}, got.NodeIds)
		assert.Len(t, got.Deadlines, 3)
		assert.Equal(t, map[uint8]blobRef{0: {"handle", 0}, 1: {"handle", 0}, 3: {"handle", 0}}, res.observedBlobs[string(hash)], "each voter's blob")
	})
}

func TestBroadcastObservationBlobs(t *testing.T) {
	t.Parallel()

	t.Run("no events", func(t *testing.T) {
		t.Parallel()
		c := testPlugin(t, sliceQueue{}, &memBlobs{}, ocr3_1types.MaxMaxObservationBytes)
		ok, err := c.broadcastObservationBlobs(t.Context(), &pb.ConsensusQueueObservation{}, &memBlobs{}, nil, 1)
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("single blob holds all events", func(t *testing.T) {
		t.Parallel()
		blobs := &memBlobs{}
		c := testPlugin(t, sliceQueue{}, blobs, ocr3_1types.MaxMaxObservationBytes)
		obs := &pb.ConsensusQueueObservation{WorkflowCapacities: map[string]uint32{}}
		events := make([]*pb.CoordinatedEvent, 0, 3)
		for _, id := range []string{"a", "b", "c"} {
			h, e := newCoordinatedEvent(t, testWF, id)
			events = append(events, e)
			_, err := observeEvent(obs, h, e, nil)
			require.NoError(t, err)
		}
		ok, err := c.broadcastObservationBlobs(t.Context(), obs, blobs, events, 1)
		require.NoError(t, err)
		assert.True(t, ok)
		assert.Len(t, blobs.payloads, 1)
		assert.Len(t, obs.BlobHandles, 1)
		assert.Len(t, obs.Events, 3)
	})

	t.Run("event and observation counts must match", func(t *testing.T) {
		t.Parallel()
		blobs := &memBlobs{}
		c := testPlugin(t, sliceQueue{}, blobs, ocr3_1types.MaxMaxObservationBytes)
		_, e := newCoordinatedEvent(t, testWF, "a")
		_, err := c.broadcastObservationBlobs(t.Context(), &pb.ConsensusQueueObservation{}, blobs, []*pb.CoordinatedEvent{e}, 1)
		require.Error(t, err)
	})
}

// Observation output must pass ValidateObservation on a peer.
func TestObservation_PassesValidateObservation(t *testing.T) {
	t.Parallel()

	events := make([]hashedEvent, 0, 2)
	for _, id := range []string{"a", "b"} {
		h, e := newCoordinatedEvent(t, testWF, id)
		events = append(events, hashedEvent{h, e})
	}
	blobs := &memBlobs{}
	c := testPlugin(t, sliceQueue{events: events}, blobs, ocr3_1types.MaxMaxObservationBytes)

	raw, err := c.Observation(t.Context(), 1, ocrtypes.AttributedQuery{}, memKV{}, blobs)
	require.NoError(t, err)

	var obs pb.ConsensusQueueObservation
	require.NoError(t, proto.Unmarshal(raw, &obs))
	assert.Len(t, obs.Events, 2)
	assert.Len(t, obs.BlobHandles, 1, "only real blob handles, no size-estimation placeholders")
	for _, e := range obs.Events {
		assert.Less(t, int(e.Blobs[0].Blob), len(obs.BlobHandles))
	}

	err = c.ValidateObservation(t.Context(), 1, ocrtypes.AttributedQuery{}, ocrtypes.AttributedObservation{Observer: 0, Observation: raw}, memKV{}, blobs)
	require.NoError(t, err)
}

func TestObservation_NoEvents(t *testing.T) {
	t.Parallel()

	c := testPlugin(t, sliceQueue{}, &memBlobs{}, ocr3_1types.MaxMaxObservationBytes)
	raw, err := c.Observation(t.Context(), 1, ocrtypes.AttributedQuery{}, memKV{}, &memBlobs{})
	require.NoError(t, err)

	var obs pb.ConsensusQueueObservation
	require.NoError(t, proto.Unmarshal(raw, &obs))
	assert.Empty(t, obs.Events)
	assert.Empty(t, obs.BlobHandles, "size-estimation placeholders must not be sent")
	require.NoError(t, c.ValidateObservation(t.Context(), 1, ocrtypes.AttributedQuery{}, ocrtypes.AttributedObservation{Observation: raw}, memKV{}, &memBlobs{}))
}

// When the observation is size-limited, Observation must still succeed with the events that fit.
func TestObservation_SizeLimited(t *testing.T) {
	t.Parallel()

	events := make([]hashedEvent, 0, 3)
	for _, id := range []string{"a", "b", "c"} {
		h, e := newCoordinatedEvent(t, testWF, id)
		events = append(events, hashedEvent{h, e})
	}

	// Find the observation size with only the first event and placeholders, then allow just that much (plus a little slack).
	probe := &pb.ConsensusQueueObservation{
		BlobHandles:        make([][]byte, maxObservedBlobHandles),
		Expiration:         timestamppb.New(time.Now()),
		WorkflowCapacities: map[string]uint32{},
	}
	c0 := testPlugin(t, sliceQueue{}, &memBlobs{}, ocr3_1types.MaxMaxObservationBytes)
	for i := range probe.BlobHandles {
		probe.BlobHandles[i] = make([]byte, c0.blobHandleMarshalledBytesUpperBound)
	}
	_, e0 := newCoordinatedEvent(t, testWF, "a")
	_, err := observeEvent(probe, events[0].hash, e0, nil)
	require.NoError(t, err)
	limit := proto.Size(probe) + 8

	blobs := &memBlobs{}
	c := testPlugin(t, sliceQueue{events: events}, blobs, limit)
	raw, err := c.Observation(t.Context(), 1, ocrtypes.AttributedQuery{}, memKV{}, blobs)
	require.NoError(t, err, "hitting the size limit must not fail the observation")

	var obs pb.ConsensusQueueObservation
	require.NoError(t, proto.Unmarshal(raw, &obs))
	assert.Len(t, obs.Events, 1)
}

// blobEvent is a coordinated event as it appears in a blob (no timestamps), plus its observation form.
func blobEvent(t testing.TB, wf, id string) (ce *pb.CoordinatedEvent, qe *pb.ConsensusQueueEvent) {
	t.Helper()
	ce = &pb.CoordinatedEvent{
		WorkflowID: wf, TriggerCapID: "trigger@1.0.0",
		Event: &pb.TriggerResponse{Event: &pb.TriggerEvent{Id: id}},
	}
	b, err := proto.MarshalOptions{Deterministic: true}.Marshal(ce)
	require.NoError(t, err)
	hash := sha256.Sum256(b)
	execID, err := workflows.GenerateExecutionIDWithTriggerIndex(wf, id, 0)
	require.NoError(t, err)
	qe = &pb.ConsensusQueueEvent{
		WorkflowID: wf, ExecutionID: execID, TimelessHash: hash[:],
		ObservedAt: ts(1), Deadline: ts(2),
		Blobs: []*pb.CoordinatedEventBlobHandle{{Blob: 0, Element: 0}},
	}
	return ce, qe
}

func TestValidateObservedEvent(t *testing.T) {
	t.Parallel()

	handles := [][]byte{handleBytes(1)}

	tests := []struct {
		name string
		// mutate may change the observed event, and/or the blob contents (returned payload replaces the default one)
		mutate  func(ce *pb.CoordinatedEvent, qe *pb.ConsensusQueueEvent)
		fetch   error
		wantErr string // empty means valid
	}{
		{name: "valid"},
		{name: "missing workflow id", mutate: func(_ *pb.CoordinatedEvent, qe *pb.ConsensusQueueEvent) { qe.WorkflowID = "" }, wantErr: "workflow id is missing"},
		{name: "workflow id not hex", mutate: func(_ *pb.CoordinatedEvent, qe *pb.ConsensusQueueEvent) { qe.WorkflowID = "not-a-workflow-id" }, wantErr: "invalid event workflow id"},
		{name: "workflow id too short", mutate: func(_ *pb.CoordinatedEvent, qe *pb.ConsensusQueueEvent) { qe.WorkflowID = "abcd" }, wantErr: "invalid event workflow id"},
		{name: "workflow id too long", mutate: func(_ *pb.CoordinatedEvent, qe *pb.ConsensusQueueEvent) { qe.WorkflowID = hexID("a", 0) + "00" }, wantErr: "invalid event workflow id"},
		{name: "missing execution id", mutate: func(_ *pb.CoordinatedEvent, qe *pb.ConsensusQueueEvent) { qe.ExecutionID = "" }, wantErr: "execution id is missing"},
		{name: "missing observed at", mutate: func(_ *pb.CoordinatedEvent, qe *pb.ConsensusQueueEvent) { qe.ObservedAt = nil }, wantErr: "observation ts is missing"},
		{name: "missing deadline", mutate: func(_ *pb.CoordinatedEvent, qe *pb.ConsensusQueueEvent) { qe.Deadline = nil }, wantErr: "deadline is missing"},
		{name: "missing hash", mutate: func(_ *pb.CoordinatedEvent, qe *pb.ConsensusQueueEvent) { qe.TimelessHash = nil }, wantErr: "hash is missing"},
		{name: "no blob ref", mutate: func(_ *pb.CoordinatedEvent, qe *pb.ConsensusQueueEvent) { qe.Blobs = nil }, wantErr: "blob is missing"},
		{name: "two blob refs", mutate: func(_ *pb.CoordinatedEvent, qe *pb.ConsensusQueueEvent) { qe.Blobs = append(qe.Blobs, qe.Blobs[0]) }, wantErr: "more than one blob"},
		{name: "blob index out of range", mutate: func(_ *pb.CoordinatedEvent, qe *pb.ConsensusQueueEvent) { qe.Blobs[0].Blob = uint32(len(handles)) }, wantErr: "out of range"}, //nolint:gosec // G115: a small test value
		{name: "element out of range", mutate: func(_ *pb.CoordinatedEvent, qe *pb.ConsensusQueueEvent) { qe.Blobs[0].Element = 5 }, wantErr: "failed to fetch event"},
		{name: "blob fetch fails", fetch: errBoom, wantErr: "boom"},
		{name: "workflow id differs from blob", mutate: func(_ *pb.CoordinatedEvent, qe *pb.ConsensusQueueEvent) { qe.WorkflowID = hexID("other", 0) }, wantErr: "does not match event workflow ID"},
		{name: "execution id differs from blob", mutate: func(_ *pb.CoordinatedEvent, qe *pb.ConsensusQueueEvent) { qe.ExecutionID = "forged" }, wantErr: "does not match event execution ID"},
		{name: "blob carries observed at", mutate: func(ce *pb.CoordinatedEvent, _ *pb.ConsensusQueueEvent) { ce.ObservedAt = ts(1) }, wantErr: "observation timestamp is not nil"},
		{name: "blob carries deadline", mutate: func(ce *pb.CoordinatedEvent, _ *pb.ConsensusQueueEvent) { ce.Deadline = ts(1) }, wantErr: "deadline timestamp is not nil"},
		{name: "hash differs from blob", mutate: func(_ *pb.CoordinatedEvent, qe *pb.ConsensusQueueEvent) { qe.TimelessHash = make([]byte, sha256.Size) }, wantErr: "does not match observation hash"},
		{name: "wrong length hash", mutate: func(_ *pb.CoordinatedEvent, qe *pb.ConsensusQueueEvent) { qe.TimelessHash = []byte{1} }, wantErr: "must be 32 bytes"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ce, qe := blobEvent(t, testWF, "evt")
			if tt.mutate != nil {
				tt.mutate(ce, qe)
			}
			f := &countingFetcher{payload: marshalEvents(t, ce), err: tt.fetch}

			err := validateObservedEvent(t.Context(), newBlobCache(f), qe, handles)
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tt.wantErr)
		})
	}

	// A peer controls the blob contents, so a malformed event must be rejected, never panic.
	t.Run("blob event without trigger event does not panic", func(t *testing.T) {
		t.Parallel()
		_, qe := blobEvent(t, testWF, "evt")
		for _, ce := range []*pb.CoordinatedEvent{
			{WorkflowID: testWF}, // no Event
			{WorkflowID: testWF, Event: &pb.TriggerResponse{}}, // no inner Event
		} {
			f := &countingFetcher{payload: marshalEvents(t, ce)}
			require.NotPanics(t, func() {
				require.Error(t, validateObservedEvent(t.Context(), newBlobCache(f), qe, handles))
			})
		}
	})
}

// observeN builds a valid observation from a plugin with n queued events with ascending deadlines.
func observeN(t *testing.T, n int) (*consensusQueuePlugin, *memBlobs, *pb.ConsensusQueueObservation) {
	t.Helper()
	events := make([]hashedEvent, 0, n)
	for i := range n {
		h, e := newCoordinatedEvent(t, testWF, fmt.Sprintf("evt-%d", i))
		e.Deadline = ts(int64(2000 + 10*i))
		events = append(events, hashedEvent{h, e})
	}
	blobs := &memBlobs{}
	c := testPlugin(t, sliceQueue{events: events}, blobs, ocr3_1types.MaxMaxObservationBytes)
	raw, err := c.Observation(t.Context(), 1, ocrtypes.AttributedQuery{}, memKV{}, blobs)
	require.NoError(t, err)
	var obs pb.ConsensusQueueObservation
	require.NoError(t, proto.Unmarshal(raw, &obs))
	require.Len(t, obs.Events, n)
	return c, blobs, &obs
}

func validate(t *testing.T, c *consensusQueuePlugin, blobs *memBlobs, obs *pb.ConsensusQueueObservation) error {
	t.Helper()
	raw, err := proto.Marshal(obs)
	require.NoError(t, err)
	return c.ValidateObservation(t.Context(), 1, ocrtypes.AttributedQuery{}, ocrtypes.AttributedObservation{Observation: raw}, memKV{}, blobs)
}

func TestValidateObservation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mutate  func(obs *pb.ConsensusQueueObservation)
		wantErr string
	}{
		{name: "valid"},
		{name: "too many blob handles", mutate: func(o *pb.ConsensusQueueObservation) {
			for len(o.BlobHandles) <= maxObservedBlobHandles {
				o.BlobHandles = append(o.BlobHandles, o.BlobHandles[0])
			}
		}, wantErr: "blob handles are limited"},
		{name: "missing expiration", mutate: func(o *pb.ConsensusQueueObservation) { o.Expiration = nil }, wantErr: "expiration is missing"},
		{name: "too many events", mutate: func(o *pb.ConsensusQueueObservation) {
			o.Events = make([]*pb.ConsensusQueueEvent, maxObservedEvents+1)
		}, wantErr: "events are limited"},
		{name: "invalid event reports its index", mutate: func(o *pb.ConsensusQueueObservation) { o.Events[1].WorkflowID = "" }, wantErr: "event 1 failed validation"},
		{name: "event with nil entry does not panic", mutate: func(o *pb.ConsensusQueueObservation) { o.Events[1] = nil }, wantErr: "event 1 failed validation"},
		{name: "missing workflow capacity", mutate: func(o *pb.ConsensusQueueObservation) { delete(o.WorkflowCapacities, testWF) }, wantErr: "missing capacity"},
		{name: "duplicate event", mutate: func(o *pb.ConsensusQueueObservation) { o.Events[2] = proto.CloneOf(o.Events[1]) }, wantErr: "duplicate workflow execution ID"},
		{name: "unsorted", mutate: func(o *pb.ConsensusQueueObservation) { o.Events[0], o.Events[1] = o.Events[1], o.Events[0] }, wantErr: "must be sorted"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, blobs, obs := observeN(t, 3)
			if tt.mutate != nil {
				tt.mutate(obs)
			}
			err := validate(t, c, blobs, obs)
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tt.wantErr)
		})
	}

	t.Run("garbage", func(t *testing.T) {
		t.Parallel()
		c, blobs, _ := observeN(t, 1)
		err := c.ValidateObservation(t.Context(), 1, ocrtypes.AttributedQuery{}, ocrtypes.AttributedObservation{Observation: []byte{0xff, 0xff}}, memKV{}, blobs)
		require.ErrorContains(t, err, "unmarshal")
	})

	t.Run("empty observation is valid", func(t *testing.T) {
		t.Parallel()
		c, blobs, obs := observeN(t, 0)
		require.NoError(t, validate(t, c, blobs, obs))
	})

	// Completeness (C4): fields that feed StateTransition must be bound to what an honest node reports.
	t.Run("workflow capacities for workflows without events are bounded", func(t *testing.T) {
		t.Parallel()
		c, blobs, obs := observeN(t, 1)
		for i := range 10_000 {
			obs.WorkflowCapacities[fmt.Sprintf("junk-%d", i)] = 1
		}
		assert.Error(t, validate(t, c, blobs, obs), "unbounded junk workflow IDs must be rejected")
	})
}

func TestObservationQuorum(t *testing.T) {
	t.Parallel()

	c := testPlugin(t, sliceQueue{}, &memBlobs{}, ocr3_1types.MaxMaxObservationBytes) // N=4, F=1
	aos := func(n int) (out []ocrtypes.AttributedObservation) {
		for i := range n {
			out = append(out, ocrtypes.AttributedObservation{Observer: commontypes.OracleID(i)})
		}
		return
	}
	for n, want := range map[int]bool{0: false, 1: false, 2: false, 3: true, 4: true} {
		got, err := c.ObservationQuorum(t.Context(), 1, ocrtypes.AttributedQuery{}, aos(n), memKV{}, nil)
		require.NoError(t, err)
		assert.Equal(t, want, got, "%d of N=4,F=1 observations (quorum is 2F+1=3)", n)
	}
}

func TestReportlessPluginMethods(t *testing.T) {
	t.Parallel()

	c := testPlugin(t, sliceQueue{}, &memBlobs{}, ocr3_1types.MaxMaxObservationBytes)

	q, err := c.Query(t.Context(), 1, memKV{}, nil)
	require.NoError(t, err)
	assert.Nil(t, q)

	reports, err := c.Reports(t.Context(), 1, nil)
	require.NoError(t, err)
	assert.Empty(t, reports)

	ok, err := c.ShouldAcceptAttestedReport(t.Context(), 1, ocr3types.ReportWithInfo[[]byte]{})
	require.NoError(t, err)
	assert.False(t, ok)
	ok, err = c.ShouldTransmitAcceptedReport(t.Context(), 1, ocr3types.ReportWithInfo[[]byte]{})
	require.NoError(t, err)
	assert.False(t, ok)
}

// nodeOpts are the non-event fields of a node's observation.
type nodeOpts struct {
	capacity   uint32
	wfCaps     map[string]uint32
	expiration int64 // seconds
	drain      bool  // reports KV drain mode
}

var stdNode = nodeOpts{capacity: 10, wfCaps: map[string]uint32{testWF: 10}, expiration: 100}

// nodeObs builds a node's observation of the events, referencing a single (shared) blob.
func nodeObs(t testing.TB, o nodeOpts, events ...hashedEvent) *pb.ConsensusQueueObservation {
	t.Helper()
	obs := &pb.ConsensusQueueObservation{
		BlobHandles:        [][]byte{handleBytes(1)},
		Expiration:         ts(o.expiration),
		Capacity:           o.capacity,
		WorkflowCapacities: maps.Clone(o.wfCaps),
		KvDrainActive:      o.drain,
	}
	if obs.WorkflowCapacities == nil {
		obs.WorkflowCapacities = map[string]uint32{}
	}
	for _, e := range events {
		_, err := observeEvent(obs, e.hash, proto.CloneOf(e.event), nil)
		require.NoError(t, err)
	}
	return obs
}

func ao(t testing.TB, id int, obs *pb.ConsensusQueueObservation) ocrtypes.AttributedObservation {
	t.Helper()
	b, err := proto.Marshal(obs)
	require.NoError(t, err)
	return ocrtypes.AttributedObservation{Observer: commontypes.OracleID(id), Observation: b} //nolint:gosec // G115: a small test value
}

// mkEvents makes n events for workflow testWF with deadlines base, base+1, ...
func mkEvents(t testing.TB, n int, base int64) []hashedEvent {
	t.Helper()
	out := make([]hashedEvent, 0, n)
	for i := range n {
		h, e := newCoordinatedEvent(t, testWF, fmt.Sprintf("evt-%02d", i))
		e.Deadline = ts(base + int64(i))
		out = append(out, hashedEvent{h, e})
	}
	return out
}

func execIDOf(t testing.TB, id string) string {
	t.Helper()
	v, err := workflows.GenerateExecutionIDWithTriggerIndex(testWF, id, 0)
	require.NoError(t, err)
	return v
}

func stateTransition(t testing.TB, c *consensusQueuePlugin, kv ocr3_1types.KeyValueStateReadWriter, aos ...ocrtypes.AttributedObservation) error {
	t.Helper()
	_, err := c.StateTransition(t.Context(), 1, ocrtypes.AttributedQuery{}, aos, kv, nil)
	return err
}

func loadState(t testing.TB, kv memKV) *kvState {
	t.Helper()
	var s kvState
	require.NoError(t, s.loadFrom(kv))
	return &s
}

// consensusIDs returns the execution IDs of the events committed to execute, and checks every blob ref resolves to the shared blob.
func consensusIDs(t *testing.T, kv memKV) []string {
	t.Helper()
	events, handles, err := readConsensusState(kv)
	require.NoError(t, err)
	ids := make([]string, 0, len(events))
	for _, e := range events {
		require.NotEmpty(t, e.Blobs, "consensus event %s must reference a blob", e.ExecutionID)
		for _, b := range e.Blobs {
			require.Less(t, int(b.Blob), len(handles))
			assert.Equal(t, handleBytes(1), handles[b.Blob])
		}
		ids = append(ids, e.ExecutionID)
	}
	return ids
}

func queuedIDs(t *testing.T, kv memKV) []string {
	t.Helper()
	s := loadState(t, kv)
	ids := make([]string, 0, len(s.queuedEvents))
	for _, e := range s.queuedEvents {
		assert.Contains(t, s.queuedBlobs, string(e.TimelessHash), "queued event %s must keep its blob ref", e.ExecutionID)
		ids = append(ids, e.ExecutionID)
	}
	slices.Sort(ids)
	return ids
}

func kvTotal(kv memKV) (n int) {
	for k, v := range kv {
		n += len(k) + len(v)
	}
	return
}

func newSTPlugin(t testing.TB) *consensusQueuePlugin {
	t.Helper()
	return testPlugin(t, sliceQueue{}, &memBlobs{}, ocr3_1types.MaxMaxObservationBytes) // N=4, F=1
}

func TestStateTransition_Lifecycle(t *testing.T) {
	t.Parallel()

	c := newSTPlugin(t)
	kv := memKV{}
	e := mkEvents(t, 1, 2000)

	// round 1: f+1 nodes observe: tracked, not queued
	require.NoError(t, stateTransition(t, c, kv,
		ao(t, 0, nodeObs(t, stdNode, e...)), ao(t, 1, nodeObs(t, stdNode, e...)), ao(t, 2, nodeObs(t, stdNode))))
	st := loadState(t, kv)
	require.Contains(t, st.observedEvents, string(e[0].hash), "observed events must be persisted")
	assert.Equal(t, []byte{0, 1}, st.observedEvents[string(e[0].hash)].NodeIds)
	assert.Equal(t, map[uint8]blobRef{0: {string(handleBytes(1)), 0}, 1: {string(handleBytes(1)), 0}}, st.observedBlobs[string(e[0].hash)])
	assert.Empty(t, st.queuedEvents)
	assert.Empty(t, consensusIDs(t, kv))

	// round 2: a third node observes: promoted, and executed in the same round
	require.NoError(t, stateTransition(t, c, kv,
		ao(t, 1, nodeObs(t, stdNode)), ao(t, 2, nodeObs(t, stdNode, e...)), ao(t, 3, nodeObs(t, stdNode))))
	st = loadState(t, kv)
	assert.Empty(t, st.observedEvents)
	assert.Empty(t, st.queuedEvents)
	assert.Equal(t, []string{execIDOf(t, "evt-00")}, consensusIDs(t, kv))

	// round 3: consensus events are only for one round
	require.NoError(t, stateTransition(t, c, kv, ao(t, 0, nodeObs(t, stdNode)), ao(t, 1, nodeObs(t, stdNode)), ao(t, 2, nodeObs(t, stdNode))))
	assert.Empty(t, consensusIDs(t, kv))
}

func TestStateTransition_QuorumInOneRound(t *testing.T) {
	t.Parallel()

	c := newSTPlugin(t)
	kv := memKV{}
	e := mkEvents(t, 1, 2000)
	require.NoError(t, stateTransition(t, c, kv,
		ao(t, 0, nodeObs(t, stdNode, e...)), ao(t, 1, nodeObs(t, stdNode, e...)), ao(t, 3, nodeObs(t, stdNode, e...))))

	assert.Equal(t, []string{execIDOf(t, "evt-00")}, consensusIDs(t, kv))
	assert.Empty(t, loadState(t, kv).observedEvents)
}

func TestStateTransition_SingleObserverIsIgnored(t *testing.T) {
	t.Parallel()

	c := newSTPlugin(t)
	kv := memKV{}
	e := mkEvents(t, 1, 2000)
	require.NoError(t, stateTransition(t, c, kv,
		ao(t, 0, nodeObs(t, stdNode, e...)), ao(t, 1, nodeObs(t, stdNode)), ao(t, 2, nodeObs(t, stdNode))))

	assert.Empty(t, loadState(t, kv).observedEvents, "a lone (possibly byzantine) observer cannot introduce an event")
	assert.Empty(t, queuedIDs(t, kv))
	assert.Empty(t, consensusIDs(t, kv))
}

func TestStateTransition_Capacity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		global        [3]uint32 // each node's reported global capacity
		wfCap         uint32
		events        int
		wantConsensus int
	}{
		{"plenty", [3]uint32{9, 9, 9}, 9, 3, 3},
		{"global F+1th lowest is zero", [3]uint32{0, 0, 9}, 9, 2, 0},
		{"global F+1th lowest is one", [3]uint32{9, 1, 1}, 9, 2, 1},
		{"byzantine high capacity cannot raise it", [3]uint32{1, 1, 1 << 31}, 9, 3, 1},
		{"workflow capacity zero", [3]uint32{9, 9, 9}, 0, 2, 0},
		{"workflow capacity one", [3]uint32{9, 9, 9}, 1, 3, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := newSTPlugin(t)
			kv := memKV{}
			e := mkEvents(t, tt.events, 2000)
			var aos []ocrtypes.AttributedObservation
			for i, g := range tt.global {
				aos = append(aos, ao(t, i, nodeObs(t, nodeOpts{capacity: g, wfCaps: map[string]uint32{testWF: tt.wfCap}, expiration: 100}, e...)))
			}
			require.NoError(t, stateTransition(t, c, kv, aos...))

			got := consensusIDs(t, kv)
			require.Len(t, got, tt.wantConsensus)
			for i := range got { // earliest deadlines go first
				assert.Equal(t, execIDOf(t, fmt.Sprintf("evt-%02d", i)), got[i])
			}
			assert.Len(t, queuedIDs(t, kv), tt.events-tt.wantConsensus, "events without capacity stay queued, with their blob refs")
		})
	}

	t.Run("queued events run when capacity returns", func(t *testing.T) {
		t.Parallel()
		c := newSTPlugin(t)
		kv := memKV{}
		e := mkEvents(t, 2, 2000)
		none := nodeOpts{capacity: 0, wfCaps: map[string]uint32{testWF: 0}, expiration: 100}
		require.NoError(t, stateTransition(t, c, kv, ao(t, 0, nodeObs(t, none, e...)), ao(t, 1, nodeObs(t, none, e...)), ao(t, 2, nodeObs(t, none, e...))))
		require.Empty(t, consensusIDs(t, kv))
		require.Len(t, queuedIDs(t, kv), 2)

		require.NoError(t, stateTransition(t, c, kv, ao(t, 0, nodeObs(t, stdNode)), ao(t, 1, nodeObs(t, stdNode)), ao(t, 2, nodeObs(t, stdNode))))
		assert.Equal(t, []string{execIDOf(t, "evt-00"), execIDOf(t, "evt-01")}, consensusIDs(t, kv))
		assert.Empty(t, queuedIDs(t, kv))
	})
}

func TestStateTransition_Expiration(t *testing.T) {
	t.Parallel()

	late := nodeOpts{capacity: 10, wfCaps: map[string]uint32{testWF: 10}, expiration: 3000}

	t.Run("expired events are dropped instead of executed", func(t *testing.T) {
		t.Parallel()
		c := newSTPlugin(t)
		kv := memKV{}
		e := mkEvents(t, 2, 2000) // deadlines before the 3000 expiration
		require.NoError(t, stateTransition(t, c, kv, ao(t, 0, nodeObs(t, late, e...)), ao(t, 1, nodeObs(t, late, e...)), ao(t, 2, nodeObs(t, late, e...))))
		assert.Empty(t, consensusIDs(t, kv))
		assert.Empty(t, queuedIDs(t, kv))
	})

	t.Run("expired observations are pruned", func(t *testing.T) {
		t.Parallel()
		c := newSTPlugin(t)
		kv := memKV{}
		e := mkEvents(t, 1, 2000)
		require.NoError(t, stateTransition(t, c, kv, ao(t, 0, nodeObs(t, stdNode, e...)), ao(t, 1, nodeObs(t, stdNode, e...)), ao(t, 2, nodeObs(t, stdNode))))
		require.NotEmpty(t, loadState(t, kv).observedEvents)

		require.NoError(t, stateTransition(t, c, kv, ao(t, 0, nodeObs(t, late)), ao(t, 1, nodeObs(t, late)), ao(t, 2, nodeObs(t, late))))
		assert.Empty(t, loadState(t, kv).observedEvents)
	})
}

// deterministicScenario runs two rounds over a mixed set of events, optionally shuffling the second round's observations.
func deterministicScenario(t *testing.T, shuffleSeed int64) memKV {
	t.Helper()
	e := mkEvents(t, 7, 2000)
	pick := func(idx ...int) (out []hashedEvent) {
		for _, i := range idx {
			out = append(out, e[i])
		}
		return
	}
	c := newSTPlugin(t)
	kv := memKV{}
	// round 1: e5, e6 observed by two nodes
	require.NoError(t, stateTransition(t, c, kv,
		ao(t, 0, nodeObs(t, stdNode, pick(5, 6)...)), ao(t, 1, nodeObs(t, stdNode, pick(5, 6)...)), ao(t, 2, nodeObs(t, stdNode))))
	// round 2 votes: e0:3 e1:3 e2:2 e3:2 e4:2 (+ e5 gains a vote and is promoted)
	aos := []ocrtypes.AttributedObservation{
		ao(t, 0, nodeObs(t, nodeOpts{capacity: 3, wfCaps: map[string]uint32{testWF: 2}, expiration: 90}, pick(0, 1, 2)...)),
		ao(t, 1, nodeObs(t, nodeOpts{capacity: 5, wfCaps: map[string]uint32{testWF: 3}, expiration: 95}, pick(0, 1, 3)...)),
		ao(t, 2, nodeObs(t, nodeOpts{capacity: 4, wfCaps: map[string]uint32{testWF: 9}, expiration: 99}, pick(0, 2, 3, 4, 5)...)),
		ao(t, 3, nodeObs(t, nodeOpts{capacity: 9, wfCaps: map[string]uint32{testWF: 1}, expiration: 80}, pick(1, 4)...)),
	}
	if shuffleSeed >= 0 {
		rand.New(rand.NewSource(shuffleSeed)).Shuffle(len(aos), func(i, j int) { aos[i], aos[j] = aos[j], aos[i] })
	}
	require.NoError(t, stateTransition(t, c, kv, aos...))
	return kv
}

// Same inputs must always produce identical KV bytes (guards against map iteration order leaking into state).
func TestStateTransition_Deterministic(t *testing.T) {
	t.Parallel()

	want := deterministicScenario(t, -1)
	require.NotEmpty(t, consensusIDs(t, want), "scenario should exercise the consensus path")
	for i := range 40 {
		require.True(t, maps.EqualFunc(want, deterministicScenario(t, -1), slices.Equal[[]byte]), "run %d differs", i)
	}
}

// Hardening: the output also does not depend on the order observations arrive in. libocr hands every oracle the same
// list in the same order, so this is not strictly required, but sorting by observer would make it so cheaply.
func TestStateTransition_ObservationOrderIndependent(t *testing.T) {
	t.Parallel()

	want := deterministicScenario(t, -1)
	for seed := range int64(40) {
		got := deterministicScenario(t, seed)
		require.True(t, maps.EqualFunc(want, got, slices.Equal[[]byte]), "seed %d: observation order changed the KV bytes", seed)
	}
}

func TestStateTransition_Errors(t *testing.T) {
	t.Parallel()

	good := func(t *testing.T) []ocrtypes.AttributedObservation {
		return []ocrtypes.AttributedObservation{ao(t, 0, nodeObs(t, stdNode)), ao(t, 1, nodeObs(t, stdNode)), ao(t, 2, nodeObs(t, stdNode))}
	}

	t.Run("garbage observation", func(t *testing.T) {
		t.Parallel()
		kv := memKV{}
		aos := append(good(t), ocrtypes.AttributedObservation{Observer: 3, Observation: []byte{0xff, 0xff}})
		require.Error(t, stateTransition(t, newSTPlugin(t), kv, aos...))
		assert.Empty(t, kv, "nothing may be written on error")
	})

	t.Run("cancelled context", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		_, err := newSTPlugin(t).StateTransition(ctx, 1, ocrtypes.AttributedQuery{}, good(t), memKV{}, nil)
		require.ErrorIs(t, err, context.Canceled)
	})

	t.Run("read error", func(t *testing.T) {
		t.Parallel()
		require.ErrorIs(t, stateTransition(t, newSTPlugin(t), &errKV{memKV: memKV{}, readErrKey: queuedEventsKey}, good(t)...), errBoom)
	})

	for _, key := range allKeys {
		t.Run("write error "+key, func(t *testing.T) {
			t.Parallel()
			require.ErrorIs(t, stateTransition(t, newSTPlugin(t), &errKV{memKV: memKV{}, writeErrKey: key}, good(t)...), errBoom)
		})
	}
}

// within runs fn, failing the test if it does not return in time (e.g. an unbounded pruning loop).
func within(t *testing.T, d time.Duration, fn func()) {
	t.Helper()
	done := make(chan struct{})
	go func() { defer close(done); fn() }()
	select {
	case <-done:
	case <-time.After(d):
		t.Fatalf("did not complete within %s", d)
	}
}

func TestStateTransition_SizeLimits(t *testing.T) {
	t.Parallel()

	t.Run("observed events are pruned to fit, least endorsed first", func(t *testing.T) {
		t.Parallel()
		c := newSTPlugin(t)
		const limit = 2000
		c.limits.MaxKeyValueModifiedKeysPlusValuesBytes = limit
		kv := memKV{}
		e := mkEvents(t, 30, 2000)

		within(t, 10*time.Second, func() {
			require.NoError(t, stateTransition(t, c, kv,
				ao(t, 0, nodeObs(t, stdNode, e...)), ao(t, 1, nodeObs(t, stdNode, e...)), ao(t, 2, nodeObs(t, stdNode))))
		})

		assert.LessOrEqual(t, kvTotal(kv), limit)
		st := loadState(t, kv)
		require.NotEmpty(t, st.observedEvents)
		require.Less(t, len(st.observedEvents), 30)
		// equal votes: ties break on hash ascending, so the retained are the k smallest hashes
		all := make([]string, 0, len(e))
		for _, ev := range e {
			all = append(all, string(ev.hash))
		}
		slices.Sort(all)
		kept := make([]string, 0, len(st.observedEvents))
		for h := range st.observedEvents {
			kept = append(kept, h)
		}
		assert.ElementsMatch(t, all[:len(kept)], kept)
	})

	t.Run("newly queued events are pruned to fit, latest deadline first", func(t *testing.T) {
		t.Parallel()
		c := newSTPlugin(t)
		const limit = 3000
		c.limits.MaxKeyValueModifiedKeysPlusValuesBytes = limit
		kv := memKV{}
		e := mkEvents(t, 30, 2000)
		none := nodeOpts{capacity: 0, wfCaps: map[string]uint32{testWF: 0}, expiration: 100} // nothing executes

		within(t, 10*time.Second, func() {
			require.NoError(t, stateTransition(t, c, kv,
				ao(t, 0, nodeObs(t, none, e...)), ao(t, 1, nodeObs(t, none, e...)), ao(t, 2, nodeObs(t, none, e...))))
		})

		assert.LessOrEqual(t, kvTotal(kv), limit)
		got := queuedIDs(t, kv)
		require.NotEmpty(t, got)
		require.Less(t, len(got), 30)
		want := make([]string, 0, len(got))
		for i := range got { // the earliest deadlines survive
			want = append(want, execIDOf(t, fmt.Sprintf("evt-%02d", i)))
		}
		assert.ElementsMatch(t, want, got)
	})

	t.Run("events that do not fit in the consensus set go back to the queue, none are lost", func(t *testing.T) {
		t.Parallel()
		const n = 20
		e := mkEvents(t, n, 2000)
		aos := func() []ocrtypes.AttributedObservation {
			big := nodeOpts{capacity: 100, wfCaps: map[string]uint32{testWF: 100}, expiration: 100}
			return []ocrtypes.AttributedObservation{ao(t, 0, nodeObs(t, big, e...)), ao(t, 1, nodeObs(t, big, e...)), ao(t, 2, nodeObs(t, big, e...))}
		}

		// unconstrained: everything executes
		full := memKV{}
		require.NoError(t, stateTransition(t, newSTPlugin(t), full, aos()...))
		require.Len(t, consensusIDs(t, full), n)

		// limit leaves room for the whole state but not for all of it in the consensus set (capped at half the limit)
		c := newSTPlugin(t)
		c.limits.MaxKeyValueModifiedKeysPlusValuesBytes = kvTotal(full) + 100
		kv := memKV{}
		require.NoError(t, stateTransition(t, c, kv, aos()...))

		executing, queued := consensusIDs(t, kv), queuedIDs(t, kv)
		assert.NotEmpty(t, executing)
		assert.Less(t, len(executing), n)
		assert.ElementsMatch(t, slices.Concat(executing, queued), func() (all []string) {
			for i := range n {
				all = append(all, execIDOf(t, fmt.Sprintf("evt-%02d", i)))
			}
			return
		}(), "every event is either executing or still queued, exactly once")
	})
}

// recCoordinator records StartExecution calls.
type recCoordinator struct {
	noopCoordinator
	mu       sync.Mutex
	started  []triggers.CoordinatedEvent
	startErr error
	capacity triggers.Capacity
	capErr   error
}

func (r *recCoordinator) AvailableCapacity(context.Context) (triggers.Capacity, error) {
	return r.capacity, r.capErr
}

func (r *recCoordinator) StartExecution(_ context.Context, e triggers.CoordinatedEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.started = append(r.started, e)
	return r.startErr
}

func (r *recCoordinator) startedEvents() []triggers.CoordinatedEvent {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]triggers.CoordinatedEvent(nil), r.started...)
}

// recQueue records removals from the local queue.
type recQueue struct {
	sliceQueue
	mu                      sync.Mutex
	removed                 []string
	sortedLimit, sortedSize int // arguments of the last Sorted call
}

func (q *recQueue) Sorted(ctx context.Context, limit, maxEventBytes int) iter.Seq2[[]byte, *pb.CoordinatedEvent] {
	q.mu.Lock()
	q.sortedLimit, q.sortedSize = limit, maxEventBytes
	q.mu.Unlock()
	return q.sliceQueue.Sorted(ctx, limit, maxEventBytes)
}

func (q *recQueue) Remove(_ context.Context, ids iter.Seq[string]) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for id := range ids {
		q.removed = append(q.removed, id)
	}
}

func (q *recQueue) removedIDs() []string {
	q.mu.Lock()
	defer q.mu.Unlock()
	return append([]string(nil), q.removed...)
}

// errFetchBlobs fails every fetch.
type errFetchBlobs struct{ memBlobs }

func (errFetchBlobs) FetchBlob(context.Context, ocr3_1types.BlobHandle) ([]byte, error) {
	return nil, errBoom
}

func newCommitPlugin(t testing.TB, coord triggers.Coordinator, q triggers.Queue, blobs ocr3_1types.BlobBroadcastFetcher) *consensusQueuePlugin {
	t.Helper()
	limits := ocr3_1types.ReportingPluginLimits{
		MaxObservationBytes:                    ocr3_1types.MaxMaxObservationBytes,
		MaxBlobPayloadBytes:                    ocr3_1types.MaxMaxBlobPayloadBytes,
		MaxKeyValueModifiedKeysPlusValuesBytes: ocr3_1types.MaxMaxKeyValueModifiedKeysPlusValuesBytes,
	}
	c := newConsensusQueuePlugin(logger.Test(t),
		ocr3types.ReportingPluginConfig{N: 4, F: 1, OracleID: 0, EstimatedRoundInterval: time.Second},
		limits, coord, q, func(context.Context) (time.Duration, error) { return time.Minute, nil }, blobs, testMetrics(t))
	require.NoError(t, c.Start(t.Context()))
	t.Cleanup(func() { require.NoError(t, c.Close()) })
	return c
}

// seedKV writes consensus (executing) and queued events to kv, all backed by one blob. Event i has medianized
// timestamps ObservedAt=4000+i and Deadline=5000+i. It returns the execution IDs in order consensus then queued.
func seedKV(t testing.TB, kv memKV, blobs *memBlobs, consensus, queued []string) (execIDs []string) {
	t.Helper()
	s := newKVState()
	ces := make([]*pb.CoordinatedEvent, 0, len(consensus)+len(queued))
	for i, id := range append(append([]string(nil), consensus...), queued...) {
		ce, qe := blobEvent(t, testWF, id)
		qe.ObservedAt, qe.Deadline = ts(int64(4000+i)), ts(int64(5000+i))
		qe.Blobs = nil
		ces = append(ces, ce)
		execIDs = append(execIDs, qe.ExecutionID)
		s.queuedBlobs.add(string(qe.TimelessHash), string(handleBytes(1)), uint32(i))
		if i < len(consensus) {
			s.consensusEvents = append(s.consensusEvents, qe)
		} else {
			s.queuedEvents.push(qe)
		}
	}
	blobs.store(handleBytes(1), marshalEvents(t, ces...))
	s.sizeOf() // prepares blob refs
	require.NoError(t, s.writeTo(&stateWriter{kv: kv}))
	return execIDs
}

func TestCommitted(t *testing.T) {
	t.Parallel()

	t.Run("starts every consensus event with the agreed timestamps and clears the local queue", func(t *testing.T) {
		t.Parallel()
		coord, q, blobs, kv := &recCoordinator{}, &recQueue{}, &memBlobs{}, memKV{}
		c := newCommitPlugin(t, coord, q, blobs)
		ids := seedKV(t, kv, blobs, []string{"a", "b"}, []string{"c"})

		require.NoError(t, c.Committed(t.Context(), 1, kv))

		started := coord.startedEvents()
		require.Len(t, started, 2)
		byTrigger := map[string]triggers.CoordinatedEvent{}
		for _, e := range started {
			byTrigger[e.Event.Event.ID] = e
		}
		for i, id := range []string{"a", "b"} {
			e, ok := byTrigger[id]
			require.True(t, ok, "event %s not started", id)
			assert.Equal(t, testWF, e.WorkflowID)
			assert.True(t, time.Unix(int64(4000+i), 0).Equal(e.ObservedAt), "agreed observedAt")
			assert.True(t, time.Unix(int64(5000+i), 0).Equal(e.Deadline), "agreed deadline")
		}
		assert.ElementsMatch(t, ids, q.removedIDs(), "both executing and queued events leave the local queue")
	})

	t.Run("nothing to do", func(t *testing.T) {
		t.Parallel()
		coord, q := &recCoordinator{}, &recQueue{}
		c := newCommitPlugin(t, coord, q, &memBlobs{})
		require.NoError(t, c.Committed(t.Context(), 1, memKV{}))
		assert.Empty(t, coord.startedEvents())
		assert.Empty(t, q.removedIDs())
	})

	t.Run("corrupt state is an error and starts nothing", func(t *testing.T) {
		t.Parallel()
		for _, key := range []string{consensusEventsKey, queuedEventsKey, blobHandlesKey} {
			coord := &recCoordinator{}
			c := newCommitPlugin(t, coord, &recQueue{}, &memBlobs{})
			require.Error(t, c.Committed(t.Context(), 1, memKV{key: {0xff, 0xff, 0xff}}), key)
		}
	})

	t.Run("a failing start does not fail the commit or block the others", func(t *testing.T) {
		t.Parallel()
		coord, q, blobs, kv := &recCoordinator{startErr: errBoom}, &recQueue{}, &memBlobs{}, memKV{}
		c := newCommitPlugin(t, coord, q, blobs)
		seedKV(t, kv, blobs, []string{"a", "b", "c"}, nil)

		require.NoError(t, c.Committed(t.Context(), 1, kv))
		assert.Len(t, coord.startedEvents(), 3)
	})

	t.Run("blob fetch failure is logged, not returned", func(t *testing.T) {
		t.Parallel()
		coord, q, kv := &recCoordinator{}, &recQueue{}, memKV{}
		blobs := &errFetchBlobs{}
		c := newCommitPlugin(t, coord, q, blobs)
		seedKV(t, kv, &blobs.memBlobs, []string{"a"}, nil)

		require.NoError(t, c.Committed(t.Context(), 1, kv))
		assert.Empty(t, coord.startedEvents())
	})

	t.Run("out of range blob ref does not panic", func(t *testing.T) {
		t.Parallel()
		coord, blobs, kv := &recCoordinator{}, &memBlobs{}, memKV{}
		c := newCommitPlugin(t, coord, &recQueue{}, blobs)
		_, qe := blobEvent(t, testWF, "a")
		qe.Blobs = []*pb.CoordinatedEventBlobHandle{{Blob: 9, Element: 0}}
		b, err := proto.Marshal(&pb.ConsensusQueueEventList{Events: []*pb.ConsensusQueueEvent{qe}})
		require.NoError(t, err)
		kv[consensusEventsKey] = b

		require.NotPanics(t, func() { require.NoError(t, c.Committed(t.Context(), 1, kv)) })
		assert.Empty(t, coord.startedEvents())
	})

	t.Run("event without a usable blob is skipped but later ones still start", func(t *testing.T) {
		t.Parallel()
		coord, blobs, kv := &recCoordinator{}, &memBlobs{}, memKV{}
		c := newCommitPlugin(t, coord, &recQueue{}, blobs)
		seedKV(t, kv, blobs, []string{"a", "b"}, nil)

		// first event points past the end of the blob
		events, handles, err := readConsensusState(kv)
		require.NoError(t, err)
		events[0].Blobs[0].Element = 99
		l, err := proto.Marshal(&pb.ConsensusQueueEventList{Events: events})
		require.NoError(t, err)
		kv[consensusEventsKey] = l
		_ = handles

		require.NoError(t, c.Committed(t.Context(), 1, kv))
		started := coord.startedEvents()
		require.Len(t, started, 1)
		assert.Equal(t, "b", started[0].Event.Event.ID)
	})
}

func TestObservation_StartsExecutionsFromState(t *testing.T) {
	t.Parallel()

	coord, q, blobs, kv := &recCoordinator{}, &recQueue{}, &memBlobs{}, memKV{}
	c := newCommitPlugin(t, coord, q, blobs)
	seedKV(t, kv, blobs, []string{"a"}, []string{"b"})

	_, err := c.Observation(t.Context(), 1, ocrtypes.AttributedQuery{}, kv, blobs)
	require.NoError(t, err)

	started := coord.startedEvents()
	require.Len(t, started, 1, "Committed may not have been called, so Observation starts executions too")
	assert.Equal(t, "a", started[0].Event.Event.ID)
	assert.Len(t, q.removedIDs(), 2)
}

func observe(t testing.TB, c *consensusQueuePlugin, kv ocr3_1types.KeyValueStateReader, blobs ocr3_1types.BlobBroadcastFetcher) (*pb.ConsensusQueueObservation, error) {
	t.Helper()
	raw, err := c.Observation(t.Context(), 1, ocrtypes.AttributedQuery{}, kv, blobs)
	if err != nil {
		return nil, err
	}
	var obs pb.ConsensusQueueObservation
	require.NoError(t, proto.Unmarshal(raw, &obs))
	return &obs, nil
}

func TestObservation_Fields(t *testing.T) {
	t.Parallel()

	t.Run("capacity errors fail the observation", func(t *testing.T) {
		t.Parallel()
		c := newCommitPlugin(t, &recCoordinator{capErr: errBoom}, &recQueue{}, &memBlobs{})
		_, err := observe(t, c, memKV{}, &memBlobs{})
		require.ErrorIs(t, err, errBoom)
	})

	t.Run("corrupt state fails the observation", func(t *testing.T) {
		t.Parallel()
		for _, key := range []string{observedEventsKey, queuedEventsKey, consensusEventsKey} {
			c := newCommitPlugin(t, &recCoordinator{}, &recQueue{}, &memBlobs{})
			_, err := observe(t, c, memKV{key: {0xff, 0xff, 0xff}}, &memBlobs{})
			require.Error(t, err, key)
		}
	})

	t.Run("expiration grace errors are tolerated", func(t *testing.T) {
		t.Parallel()
		c := newCommitPlugin(t, &recCoordinator{}, &recQueue{}, &memBlobs{})
		c.expirationDur = func(context.Context) (time.Duration, error) { return 0, errBoom }
		before := time.Now()
		obs, err := observe(t, c, memKV{}, &memBlobs{})
		require.NoError(t, err)
		assert.False(t, obs.Expiration.AsTime().Before(before.Add(-time.Second)), "falls back to no grace period")
	})

	t.Run("expiration is observation time minus grace", func(t *testing.T) {
		t.Parallel()
		c := newCommitPlugin(t, &recCoordinator{}, &recQueue{}, &memBlobs{}) // 1 minute grace
		before := time.Now()
		obs, err := observe(t, c, memKV{}, &memBlobs{})
		require.NoError(t, err)
		assert.WithinDuration(t, before.Add(-time.Minute), obs.Expiration.AsTime(), 5*time.Second)
	})

	t.Run("global capacity is clamped to uint32", func(t *testing.T) {
		t.Parallel()
		c := newCommitPlugin(t, &recCoordinator{capacity: triggers.Capacity{Global: math.MaxInt}}, &recQueue{}, &memBlobs{})
		obs, err := observe(t, c, memKV{}, &memBlobs{})
		require.NoError(t, err)
		assert.Equal(t, uint32(math.MaxUint32), obs.Capacity)
	})

	t.Run("workflow capacities cover workflows with events in state or locally, and nothing else", func(t *testing.T) {
		t.Parallel()
		blobs, kv := &memBlobs{}, memKV{}
		h, e := newCoordinatedEvent(t, "local-wf", "local")
		coord := &recCoordinator{capacity: triggers.Capacity{Global: 5, PerWorkflow: map[string]int{"local-wf": 3, "queued-wf": 7, "observed-wf": 8, "unrelated": 9}}}
		c := newCommitPlugin(t, coord, &recQueue{sliceQueue: sliceQueue{events: []hashedEvent{{h, e}}}}, blobs)

		// state: a queued event for queued-wf, an observed event for observed-wf
		s := newKVState()
		s.queuedEvents.push(&pb.ConsensusQueueEvent{WorkflowID: "queued-wf", ExecutionID: "x", TimelessHash: []byte("q"), ObservedAt: ts(1), Deadline: ts(2)})
		s.observedEvents["o"] = obsEvent("observed-wf", "y", 100, 5, 6)
		s.sizeOf()
		require.NoError(t, s.writeTo(&stateWriter{kv: kv}))

		obs, err := observe(t, c, kv, blobs)
		require.NoError(t, err)
		assert.Equal(t, map[string]uint32{"local-wf": 3, "queued-wf": 7, "observed-wf": 8}, obs.WorkflowCapacities)
	})

	t.Run("events this node already observed in state are not re-observed", func(t *testing.T) {
		t.Parallel()
		blobs, kv := &memBlobs{}, memKV{}
		mine, theirs := mkEvents(t, 2, 2000), mkEvents(t, 4, 2000)[2:]
		c := newCommitPlugin(t, &recCoordinator{}, &recQueue{sliceQueue: sliceQueue{events: append(append([]hashedEvent{}, mine...), theirs...)}}, blobs)

		s := newKVState()
		s.observedEvents[string(mine[0].hash)] = obsEvent(testWF, "m0", 100, 0, 1) // includes OracleID 0 (this node)
		s.observedEvents[string(mine[1].hash)] = obsEvent(testWF, "m1", 100, 2, 3) // does not
		s.sizeOf()
		require.NoError(t, s.writeTo(&stateWriter{kv: kv}))

		obs, err := observe(t, c, kv, blobs)
		require.NoError(t, err)
		got := make([][]byte, 0, len(obs.Events))
		for _, e := range obs.Events {
			got = append(got, e.TimelessHash)
		}
		assert.NotContains(t, got, mine[0].hash, "already observed by this node")
		assert.Contains(t, got, mine[1].hash, "observed by others only, so this node adds its vote")
		assert.Contains(t, got, theirs[0].hash)
		assert.Contains(t, got, theirs[1].hash)
	})
}

func TestObservation_BlobBatching(t *testing.T) {
	t.Parallel()

	// checkRefs verifies every observed event's blob ref resolves to that event in a broadcast blob.
	checkRefs := func(t *testing.T, obs *pb.ConsensusQueueObservation, blobs *memBlobs) {
		t.Helper()
		for i, e := range obs.Events {
			require.Less(t, int(e.Blobs[0].Blob), len(obs.BlobHandles), "event %d", i)
			payload, ok := blobs.byHandle[string(obs.BlobHandles[e.Blobs[0].Blob])] // broadcast order is not defined, so go by handle
			require.True(t, ok, "event %d", i)
			_, raw, err := extractEvent(payload, int(e.Blobs[0].Element))
			require.NoError(t, err, "event %d", i)
			sum := sha256.Sum256(raw)
			assert.Equal(t, e.TimelessHash, sum[:], "event %d resolves to its own blob element", i)
		}
	}
	events := func(t *testing.T, n int) []hashedEvent {
		return mkEvents(t, n, time.Now().Add(time.Hour).Unix())
	}

	t.Run("events spread over several blobs keep correct refs", func(t *testing.T) {
		t.Parallel()
		blobs := &memBlobs{}
		ev := events(t, 6)
		c := testPlugin(t, sliceQueue{events: ev}, blobs, ocr3_1types.MaxMaxObservationBytes)
		c.limits.MaxBlobPayloadBytes = proto.Size(&pb.CoordinatedEventList{Events: []*pb.CoordinatedEvent{ev[0].event, ev[1].event}}) // two per blob

		obs, err := observe(t, c, memKV{}, blobs)
		require.NoError(t, err)
		require.Len(t, obs.Events, 6)
		assert.Len(t, obs.BlobHandles, 3)
		assert.Len(t, blobs.payloads, 3)
		checkRefs(t, obs, blobs)
	})

	t.Run("more events than blob handles are truncated", func(t *testing.T) {
		t.Parallel()
		blobs := &memBlobs{}
		ev := events(t, 15)
		c := testPlugin(t, sliceQueue{events: ev}, blobs, ocr3_1types.MaxMaxObservationBytes)
		c.limits.MaxBlobPayloadBytes = proto.Size(&pb.CoordinatedEventList{Events: []*pb.CoordinatedEvent{ev[0].event}}) // one per blob

		obs, err := observe(t, c, memKV{}, blobs)
		require.NoError(t, err)
		require.Len(t, obs.BlobHandles, maxObservedBlobHandles)
		require.Len(t, obs.Events, maxObservedBlobHandles, "events beyond the last blob are dropped")
		checkRefs(t, obs, blobs)
		for i, e := range obs.Events { // earliest first
			assert.Equal(t, ev[i].hash, e.TimelessHash)
		}
		assert.Contains(t, obs.WorkflowCapacities, testWF)
	})

	t.Run("size limited observations keep only the events that fit", func(t *testing.T) {
		t.Parallel()
		blobs := &memBlobs{}
		ev := events(t, 10)
		c := testPlugin(t, sliceQueue{events: ev}, blobs, ocr3_1types.MaxMaxObservationBytes)

		// Observation reserves space for maxObservedBlobHandles handles up front, so size a probe the same way for 5 events.
		probe := &pb.ConsensusQueueObservation{Expiration: ts(time.Now().Unix()), WorkflowCapacities: map[string]uint32{}}
		for range maxObservedBlobHandles {
			probe.BlobHandles = append(probe.BlobHandles, make([]byte, c.blobHandleMarshalledBytesUpperBound))
		}
		for _, e := range ev[:5] {
			_, err := observeEvent(probe, e.hash, proto.CloneOf(e.event), nil)
			require.NoError(t, err)
		}
		c.limits.MaxObservationBytes = proto.Size(probe) + 16 // slack for timestamp nanos; a 6th event needs far more

		obs, err := observe(t, c, memKV{}, blobs)
		require.NoError(t, err)
		require.Len(t, obs.Events, 5)
		assert.LessOrEqual(t, proto.Size(obs), c.limits.MaxObservationBytes)
		checkRefs(t, obs, blobs)
		for i, e := range obs.Events {
			assert.Equal(t, ev[i].hash, e.TimelessHash, "the earliest events are kept")
		}
	})
}

// Known issues (SHOULD FIX): libocr says a BlobBroadcast failure should not make Observation fail, and one event that
// can never be broadcast must not wedge this node's observations forever. These fail until handled.
func TestObservation_BlobFailuresAreNotFatal(t *testing.T) {
	t.Parallel()

	t.Run("broadcast error", func(t *testing.T) {
		t.Parallel()
		blobs := &hintBlobs{err: errBoom}
		c := testPlugin(t, sliceQueue{events: mkEvents(t, 2, time.Now().Add(time.Hour).Unix())}, blobs, ocr3_1types.MaxMaxObservationBytes)
		obs, err := observe(t, c, memKV{}, blobs)
		require.NoError(t, err, "a failed broadcast should degrade to an observation without those events")
		assert.Empty(t, obs.Events)
	})

	t.Run("event larger than a blob", func(t *testing.T) {
		t.Parallel()
		blobs := &memBlobs{}
		c := testPlugin(t, sliceQueue{events: mkEvents(t, 2, time.Now().Add(time.Hour).Unix())}, blobs, ocr3_1types.MaxMaxObservationBytes)
		c.limits.MaxBlobPayloadBytes = 1
		obs, err := observe(t, c, memKV{}, blobs)
		require.NoError(t, err, "an unbroadcastable event should be dropped, not fail every observation")
		assert.Empty(t, obs.Events)
	})
}

// The queue drops events that can never be observed, so the plugin must pass it the largest event that fits in a blob.
func TestObservation_PassesMaxEventBytesToQueue(t *testing.T) {
	t.Parallel()

	q := &recQueue{}
	c := newCommitPlugin(t, &recCoordinator{}, q, &memBlobs{})
	c.limits.MaxBlobPayloadBytes = 1000
	_, err := observe(t, c, memKV{}, &memBlobs{})
	require.NoError(t, err)

	assert.Equal(t, maxObservedEvents, q.sortedLimit)
	assert.Equal(t, 1000-maxEventFramingOverhead, q.sortedSize)

	t.Run("limits too small to hold any event disable filtering rather than dropping everything", func(t *testing.T) {
		c.limits.MaxBlobPayloadBytes = 0
		_, err := observe(t, c, memKV{}, &memBlobs{})
		require.NoError(t, err)
		assert.Zero(t, q.sortedSize)
	})
}

// An event of exactly maxEventBytes must fit alone in a blob of the configured size, for any limit libocr allows.
func TestMaxEventBytes(t *testing.T) {
	t.Parallel()

	eventOfSize := func(k int) *pb.CoordinatedEvent { // event whose encoded size grows with k
		e, _ := blobEvent(t, testWF, "evt")
		e.Event.Event.TriggerType = strings.Repeat("x", k)
		return e
	}
	for _, limit := range []int{6, 127, 128, 129, 1 << 14, 1<<21 - 1, 1 << 21, ocr3_1types.MaxMaxBlobPayloadBytes} {
		t.Run(strconv.Itoa(limit), func(t *testing.T) {
			t.Parallel()
			c := testPlugin(t, sliceQueue{}, &memBlobs{}, ocr3_1types.MaxMaxObservationBytes)
			c.limits.MaxBlobPayloadBytes = limit
			maxSize := c.maxEventBytes()
			require.GreaterOrEqual(t, maxSize, 0)
			require.LessOrEqual(t, maxSize, limit)

			// largest event that passes the filter
			if proto.Size(eventOfSize(0)) > maxSize {
				return // nothing fits at this limit
			}
			k := sort.Search(limit, func(k int) bool { return proto.Size(eventOfSize(k)) > maxSize }) - 1
			e := eventOfSize(k)
			require.LessOrEqual(t, proto.Size(e), maxSize)
			assert.LessOrEqual(t, proto.Size(&pb.CoordinatedEventList{Events: []*pb.CoordinatedEvent{e}}), limit,
				"an event passing the filter must fit alone in a blob")
		})
	}
}

func TestValidateObservation_Fields(t *testing.T) {
	t.Parallel()

	badTS := &timestamppb.Timestamp{Seconds: 1, Nanos: 2_000_000_000}
	tests := []struct {
		name    string
		mutate  func(c *consensusQueuePlugin, o *pb.ConsensusQueueObservation)
		wantErr string
	}{
		{"invalid expiration", func(_ *consensusQueuePlugin, o *pb.ConsensusQueueObservation) { o.Expiration = badTS }, "invalid expiration"},
		{"invalid event observed at", func(_ *consensusQueuePlugin, o *pb.ConsensusQueueObservation) { o.Events[0].ObservedAt = badTS }, "invalid event observation ts"},
		{"invalid event deadline", func(_ *consensusQueuePlugin, o *pb.ConsensusQueueObservation) { o.Events[0].Deadline = badTS }, "invalid event deadline"},
		{"blob handle longer than any real handle", func(c *consensusQueuePlugin, o *pb.ConsensusQueueObservation) {
			o.BlobHandles = append(o.BlobHandles, make([]byte, c.blobHandleMarshalledBytesUpperBound+1))
		}, "more than the maximum"},
		{"empty hash", func(_ *consensusQueuePlugin, o *pb.ConsensusQueueObservation) { o.Events[0].TimelessHash = []byte{} }, "hash is missing"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, blobs, obs := observeN(t, 2)
			tt.mutate(c, obs)
			require.ErrorContains(t, validate(t, c, blobs, obs), tt.wantErr)
		})
	}
}

func TestObservationsAggregate_Omissions(t *testing.T) {
	t.Parallel()
	const f = 1

	mk := func(wfCaps map[string]uint32) *pb.ConsensusQueueObservation {
		return &pb.ConsensusQueueObservation{Expiration: ts(100), Capacity: 5, WorkflowCapacities: wfCaps}
	}
	aoOf := func(t *testing.T, obs ...*pb.ConsensusQueueObservation) (out []ocrtypes.AttributedObservation) {
		for i, o := range obs {
			out = append(out, ao(t, i, o))
		}
		return
	}

	t.Run("a workflow one node omits still gets the capacity the others report", func(t *testing.T) {
		t.Parallel()
		var res observations
		res.aggregate(t.Context(), f, aoOf(t, mk(map[string]uint32{testWF: 9}), mk(map[string]uint32{testWF: 9}), mk(nil)))
		require.NoError(t, res.err)
		assert.Equal(t, uint32(9), res.fPlus1thCapPerWF[testWF], "F omissions count as zero, which sorted[F] tolerates")
	})

	t.Run("a workflow most nodes omit gets zero", func(t *testing.T) {
		t.Parallel()
		var res observations
		res.aggregate(t.Context(), f, aoOf(t, mk(map[string]uint32{testWF: 9}), mk(nil), mk(nil)))
		require.NoError(t, res.err)
		assert.Contains(t, res.fPlus1thCapPerWF, testWF)
		assert.Zero(t, res.fPlus1thCapPerWF[testWF])
	})

	t.Run("too few observations is an error, not a panic", func(t *testing.T) {
		t.Parallel()
		var res observations
		require.NotPanics(t, func() { res.aggregate(t.Context(), f, aoOf(t, mk(nil))) })
		require.Error(t, res.err)
		require.NotPanics(t, func() { res.aggregate(t.Context(), f, nil) })
		require.Error(t, res.err)

		kv := memKV{}
		require.Error(t, stateTransition(t, newSTPlugin(t), kv, aoOf(t, mk(nil))...))
		assert.Empty(t, kv)
	})
}

func TestNewConsensusQueuePluginFactory_Validation(t *testing.T) {
	t.Parallel()

	cfg := &PluginConfig{}
	tests := []struct {
		name  string
		cfg   *PluginConfig
		coord triggers.Coordinator
		queue triggers.Queue
	}{
		{"nil cfg", nil, noopCoordinator{}, noopQueue{}},
		{"nil coordinator", cfg, nil, noopQueue{}},
		{"nil queue", cfg, noopCoordinator{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f, err := NewConsensusQueuePluginFactory(logger.Test(t), tt.cfg, limits.Factory{}, tt.coord, tt.queue, noop.NewMeterProvider().Meter("test"))
			require.Error(t, err)
			assert.Nil(t, f)
		})
	}
}

func TestNewConsensusQueuePluginFactory_NilMeter(t *testing.T) {
	t.Parallel()
	f, err := NewConsensusQueuePluginFactory(logger.Test(t), &PluginConfig{}, limits.Factory{}, noopCoordinator{}, noopQueue{}, nil)
	require.ErrorContains(t, err, "meter")
	assert.Nil(t, f)
}

func TestBroadcastAsBlob_LargestFit(t *testing.T) {
	t.Parallel()

	blobs := &hintBlobs{}
	c := testPlugin(t, sliceQueue{}, blobs, ocr3_1types.MaxMaxObservationBytes)
	events := eventsWithDeadline(t, 50, time.Now().Add(time.Hour))
	sizeOf := func(n int) int { return proto.Size(&pb.CoordinatedEventList{Events: events[:n]}) }

	for _, want := range []int{1, 2, 17, 36, 49, 50} {
		c.limits.MaxBlobPayloadBytes = sizeOf(want) // exactly fits want events
		_, n, err := c.broadcastAsBlob(t.Context(), blobs, events, 1)
		require.NoError(t, err)
		assert.Equal(t, want, n, "limit == size of %d events", want)

		if want > 1 {
			c.limits.MaxBlobPayloadBytes = sizeOf(want) - 1 // one byte short
			_, n, err = c.broadcastAsBlob(t.Context(), blobs, events, 1)
			require.NoError(t, err)
			assert.Equal(t, want-1, n, "one byte short of %d events", want)
		}
	}
}

// multiBlobObservation builds an observation of 6 events spread over 3 blobs (two per blob).
func multiBlobObservation(t *testing.T) (*consensusQueuePlugin, *memBlobs, *pb.ConsensusQueueObservation) {
	t.Helper()
	ev := mkEvents(t, 6, time.Now().Add(time.Hour).Unix())
	blobs := &memBlobs{}
	c := testPlugin(t, sliceQueue{events: ev}, blobs, ocr3_1types.MaxMaxObservationBytes)
	c.limits.MaxBlobPayloadBytes = proto.Size(&pb.CoordinatedEventList{Events: []*pb.CoordinatedEvent{ev[0].event, ev[1].event}})
	obs, err := observe(t, c, memKV{}, blobs)
	require.NoError(t, err)
	require.Len(t, obs.Events, 6)
	require.Len(t, obs.BlobHandles, 3)
	return c, blobs, obs
}

func TestValidateObservation_MultipleBlobs(t *testing.T) {
	t.Parallel()

	t.Run("an observation spanning several blobs is valid", func(t *testing.T) {
		t.Parallel()
		c, blobs, obs := multiBlobObservation(t)
		require.NoError(t, validate(t, c, blobs, obs))
	})

	t.Run("an event pointing at another blob is rejected", func(t *testing.T) {
		t.Parallel()
		c, blobs, obs := multiBlobObservation(t)
		obs.Events[0].Blobs[0].Blob = 2 // blob 2 holds different events
		require.Error(t, validate(t, c, blobs, obs))
	})

	t.Run("an element past the end of its blob is rejected", func(t *testing.T) {
		t.Parallel()
		c, blobs, obs := multiBlobObservation(t)
		obs.Events[5].Blobs[0].Element = 2 // blobs hold two events
		require.Error(t, validate(t, c, blobs, obs))
	})

	t.Run("a blob that cannot be fetched is rejected", func(t *testing.T) {
		t.Parallel()
		c, blobs, obs := multiBlobObservation(t)
		delete(blobs.byHandle, marshalHandle(testBlobHandle(2))) // the second blob is unavailable
		require.ErrorContains(t, validate(t, c, blobs, obs), "unknown blob handle")
	})

	t.Run("the same blob is fetched once however many events it holds", func(t *testing.T) {
		t.Parallel()
		c, blobs, obs := multiBlobObservation(t)
		counting := &fetchCountingBlobs{memBlobs: blobs}
		raw, err := proto.Marshal(obs)
		require.NoError(t, err)
		require.NoError(t, c.ValidateObservation(t.Context(), 1, ocrtypes.AttributedQuery{}, ocrtypes.AttributedObservation{Observation: raw}, memKV{}, counting))
		assert.Equal(t, 3, int(counting.fetches.Load()), "6 events over 3 blobs")
	})
}

type fetchCountingBlobs struct {
	*memBlobs
	fetches atomic.Int32
}

func (f *fetchCountingBlobs) FetchBlob(ctx context.Context, h ocr3_1types.BlobHandle) ([]byte, error) {
	f.fetches.Add(1)
	return f.memBlobs.FetchBlob(ctx, h)
}

// ctxBlobs fails like libocr does when its context is done.
type ctxBlobs struct{}

func (ctxBlobs) BroadcastBlob(ctx context.Context, _ []byte, _ ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	return ocr3_1types.BlobHandle{}, ctx.Err()
}

func (ctxBlobs) FetchBlob(ctx context.Context, _ ocr3_1types.BlobHandle) ([]byte, error) {
	return nil, ctx.Err()
}

func cancelledCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	return ctx
}

// A cancelled context may never turn into a successful but different result: determinism requires that a function
// that returns without an error behaves the same however its context behaves. Each path must fail instead.
func TestContextCancellation(t *testing.T) {
	t.Parallel()

	t.Run("aggregate reports cancellation as an error instead of returning a partial result", func(t *testing.T) {
		t.Parallel()
		e := mkEvents(t, 2, 2000)
		withEvents := []ocrtypes.AttributedObservation{
			ao(t, 0, nodeObs(t, stdNode, e...)), ao(t, 1, nodeObs(t, stdNode, e...)), ao(t, 2, nodeObs(t, stdNode, e...)),
		}
		withoutEvents := []ocrtypes.AttributedObservation{ao(t, 0, nodeObs(t, stdNode)), ao(t, 1, nodeObs(t, stdNode)), ao(t, 2, nodeObs(t, stdNode))}
		for name, aos := range map[string][]ocrtypes.AttributedObservation{"with events": withEvents, "without events": withoutEvents} {
			var res observations
			res.aggregate(cancelledCtx(t), 1, aos)
			require.ErrorIs(t, res.err, context.Canceled, name)
		}
	})

	t.Run("StateTransition fails and writes nothing", func(t *testing.T) {
		t.Parallel()
		e := mkEvents(t, 1, 2000)
		kv := memKV{}
		_, err := newSTPlugin(t).StateTransition(cancelledCtx(t), 1, ocrtypes.AttributedQuery{},
			[]ocrtypes.AttributedObservation{ao(t, 0, nodeObs(t, stdNode, e...)), ao(t, 1, nodeObs(t, stdNode, e...)), ao(t, 2, nodeObs(t, stdNode, e...))}, kv, nil)
		require.ErrorIs(t, err, context.Canceled)
		assert.Empty(t, kv)
	})

	t.Run("ValidateObservation does not accept events it could not fetch", func(t *testing.T) {
		t.Parallel()
		c, _, obs := observeN(t, 2)
		raw, err := proto.Marshal(obs)
		require.NoError(t, err)
		err = c.ValidateObservation(cancelledCtx(t), 1, ocrtypes.AttributedQuery{}, ocrtypes.AttributedObservation{Observation: raw}, memKV{}, ctxBlobs{})
		require.ErrorIs(t, err, context.Canceled)
	})

	t.Run("Observation fails instead of degrading to an observation without events", func(t *testing.T) {
		t.Parallel()
		ev := mkEvents(t, 2, time.Now().Add(time.Hour).Unix())
		c := testPlugin(t, sliceQueue{events: ev}, ctxBlobs{}, ocr3_1types.MaxMaxObservationBytes)
		_, err := c.Observation(cancelledCtx(t), 1, ocrtypes.AttributedQuery{}, memKV{}, ctxBlobs{})
		require.ErrorIs(t, err, context.Canceled)
	})

	t.Run("Observation fails if one of several blob broadcasts is cancelled", func(t *testing.T) {
		t.Parallel()
		ev := mkEvents(t, 4, time.Now().Add(time.Hour).Unix())
		ctx, cancel := context.WithCancel(t.Context())
		blobs := &cancelOnSecondBroadcast{cancel: cancel}
		c := testPlugin(t, sliceQueue{events: ev}, blobs, ocr3_1types.MaxMaxObservationBytes)
		c.limits.MaxBlobPayloadBytes = proto.Size(&pb.CoordinatedEventList{Events: []*pb.CoordinatedEvent{ev[0].event, ev[1].event}})
		_, err := c.Observation(ctx, 1, ocrtypes.AttributedQuery{}, memKV{}, blobs)
		require.ErrorIs(t, err, context.Canceled, "must fail, not return an observation with only the blob that succeeded")
	})

	t.Run("Committed returns promptly and does not panic", func(t *testing.T) {
		t.Parallel()
		coord, blobs, kv := &recCoordinator{}, &memBlobs{}, memKV{}
		c := newCommitPlugin(t, coord, &recQueue{}, blobs)
		seedKV(t, kv, blobs, []string{"a", "b"}, []string{"c"})
		within(t, 5*time.Second, func() {
			require.NoError(t, c.Committed(cancelledCtx(t), 1, kv))
		})
	})
}

// cancelOnSecondBroadcast succeeds once, then cancels the context and fails like libocr would. Calls may be concurrent.
type cancelOnSecondBroadcast struct {
	memBlobs
	cancel context.CancelFunc
	calls  atomic.Int32
}

func (c *cancelOnSecondBroadcast) BroadcastBlob(ctx context.Context, payload []byte, hint ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	if c.calls.Add(1) > 1 {
		c.cancel()
		return ocr3_1types.BlobHandle{}, ctx.Err()
	}
	return c.memBlobs.BroadcastBlob(ctx, payload, hint)
}

// failBlobsWith fails to broadcast the blob that starts with the given event, and broadcasts the others.
type failBlobsWith struct {
	memBlobs
	firstID string
}

func (f *failBlobsWith) BroadcastBlob(ctx context.Context, payload []byte, hint ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	if first, _, err := extractEvent(payload, 0); err == nil && first.Event.Event.Id == f.firstID {
		return ocr3_1types.BlobHandle{}, errBoom
	}
	return f.memBlobs.BroadcastBlob(ctx, payload, hint)
}

func TestObservation_PartialBlobFailure(t *testing.T) {
	t.Parallel()

	ev := mkEvents(t, 6, time.Now().Add(time.Hour).Unix()) // two events per blob
	perBlob := proto.Size(&pb.CoordinatedEventList{Events: []*pb.CoordinatedEvent{ev[0].event, ev[1].event}})

	for _, tt := range []struct {
		name       string
		failBlobAt int // index of the first event of the blob that fails
		wantEvents int
	}{
		{"first", 0, 4},
		{"middle", 2, 4},
		{"last", 4, 4},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			blobs := &failBlobsWith{firstID: ev[tt.failBlobAt].event.Event.Event.Id}
			c := testPlugin(t, sliceQueue{events: ev}, blobs, ocr3_1types.MaxMaxObservationBytes)
			c.limits.MaxBlobPayloadBytes = perBlob

			obs, err := observe(t, c, memKV{}, blobs)
			require.NoError(t, err)
			require.Len(t, obs.Events, tt.wantEvents)
			assert.Len(t, obs.BlobHandles, tt.wantEvents/2)
			var i, oi int
			for _, e := range obs.Events {
				for tt.failBlobAt == oi {
					oi += 2
				}
				assert.Equal(t, ev[oi].hash, e.TimelessHash, "events stay in order")
				assert.Equal(t, uint32(i/2), e.Blobs[0].Blob)
				assert.Equal(t, uint32(i%2), e.Blobs[0].Element)
				i++
				oi++
			}
			if tt.wantEvents > 0 {
				assert.Contains(t, obs.WorkflowCapacities, testWF, "capacities cover the events kept")
				require.NoError(t, validate(t, c, &blobs.memBlobs, obs), "what remains is a valid observation")
			}
		})
	}
}

// ---- Benchmarks of full (packed to the limit) rounds ----
//
// Run with: go test -run '^$' -bench Full -benchtime 3x ./core/services/ocr3_1/plugins/queue/
//
// KV and blob calls sleep, to model what the plugin waits on in a real round. The durations are ASSUMPTIONS, not
// measurements from a running DON: libocr's KV is a local embedded database (sub-millisecond), while a blob broadcast
// waits for availability signatures from peers over p2p (tens of milliseconds, growing with size) and a fetch pulls the
// blob from a peer. Replace the numbers with measured ones to re-calibrate. The "avg" profile makes every call take
// the average, and "max" makes every call take the maximum, so together they bound the round.

// ioProfile is how long each KV and blob operation takes: a fixed part plus a part per MiB moved.
type ioProfile struct {
	name                               string
	kvRead, kvWrite, kvPerMiB          time.Duration
	blobBroadcast, blobBroadcastPerMiB time.Duration
	blobFetch, blobFetchPerMiB         time.Duration
}

var ioProfiles = []ioProfile{
	{name: "noIO"},
	{
		name:   "avgIO",
		kvRead: 200 * time.Microsecond, kvWrite: 300 * time.Microsecond, kvPerMiB: time.Millisecond,
		blobBroadcast: 30 * time.Millisecond, blobBroadcastPerMiB: 10 * time.Millisecond,
		blobFetch: 10 * time.Millisecond, blobFetchPerMiB: 5 * time.Millisecond,
	},
	{
		name:   "maxIO",
		kvRead: 2 * time.Millisecond, kvWrite: 3 * time.Millisecond, kvPerMiB: 3 * time.Millisecond,
		blobBroadcast: 450 * time.Millisecond, blobBroadcastPerMiB: 60 * time.Millisecond, // 750ms for a 5MiB blob: observed worst case is 600-800ms
		blobFetch: 100 * time.Millisecond, blobFetchPerMiB: 15 * time.Millisecond,
	},
}

func perMiB(d time.Duration, n int) time.Duration { return d * time.Duration(n) / (1 << 20) }

// delayKV adds the profile's latency to KV reads and writes.
type delayKV struct {
	ocr3_1types.KeyValueStateReadWriter
	p ioProfile
}

func (d delayKV) Read(k []byte) ([]byte, error) {
	v, err := d.KeyValueStateReadWriter.Read(k)
	time.Sleep(d.p.kvRead + perMiB(d.p.kvPerMiB, len(v)))
	return v, err
}

func (d delayKV) Write(k, v []byte) error {
	time.Sleep(d.p.kvWrite + perMiB(d.p.kvPerMiB, len(v)))
	return d.KeyValueStateReadWriter.Write(k, v)
}

// delayBlobs adds the profile's latency to blob broadcasts and fetches.
type delayBlobs struct {
	*memBlobs
	p ioProfile
}

func (d delayBlobs) BroadcastBlob(ctx context.Context, payload []byte, hint ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	time.Sleep(d.p.blobBroadcast + perMiB(d.p.blobBroadcastPerMiB, len(payload)))
	return d.memBlobs.BroadcastBlob(ctx, payload, hint)
}

func (d delayBlobs) FetchBlob(ctx context.Context, h ocr3_1types.BlobHandle) ([]byte, error) {
	v, err := d.memBlobs.FetchBlob(ctx, h)
	time.Sleep(d.p.blobFetch + perMiB(d.p.blobFetchPerMiB, len(v)))
	return v, err
}

// packedEvents makes n events across 50 workflows, each with a trigger payload of the given size, in deadline order.
func packedEvents(tb testing.TB, n, payload int) []hashedEvent {
	tb.Helper()
	out := make([]hashedEvent, 0, n)
	for i := range n {
		e := eventOfPayload(i, payload)
		e.WorkflowID = hexID(testWF, i%50)
		b, err := e.MarshalTimeless()
		require.NoError(tb, err)
		sum := sha256.Sum256(b)
		e.ObservedAt = ts(1_800_000_000 + int64(i))
		e.Deadline = ts(4_000_000_000 + int64(i))
		out = append(out, hashedEvent{sum[:], e})
	}
	return out
}

// packShape is a full observation: maxObservedEvents events of one payload size.
type packShape struct {
	name    string
	payload int
}

func packShapes() []packShape {
	// the largest payload that still lets maxObservedEvents events fill maxObservedBlobHandles blobs
	maxPayload := sort.Search(ocr3_1types.MaxMaxBlobPayloadBytes, func(p int) bool {
		blobs, _, _ := blobsNeeded(p, ocr3_1types.MaxMaxBlobPayloadBytes)
		return blobs == 0 || blobs > maxObservedBlobHandles
	}) - 1
	return []packShape{
		{"1000xCron128B", 128},
		{"1000xEVM8KB", 8 << 10},
		{fmt.Sprintf("1000xMaxPacked%dKB", maxPayload>>10), maxPayload},
	}
}

func fullLimitsPlugin(b *testing.B, events []hashedEvent, blobs ocr3_1types.BlobBroadcastFetcher) *consensusQueuePlugin {
	b.Helper()
	c := testPlugin(b, sliceQueue{events: events}, blobs, ocr3_1types.MaxMaxObservationBytes)
	c.limits.MaxBlobPayloadBytes = ocr3_1types.MaxMaxBlobPayloadBytes
	c.limits.MaxKeyValueModifiedKeysPlusValuesBytes = ocr3_1types.MaxMaxKeyValueModifiedKeysPlusValuesBytes
	return c
}

// BenchmarkFullObservation measures building and broadcasting an observation packed to maxObservedEvents.
func BenchmarkFullObservation(b *testing.B) {
	for _, shape := range packShapes() {
		for _, io := range ioProfiles {
			b.Run(shape.name+"/"+io.name, func(b *testing.B) {
				events := packedEvents(b, maxObservedEvents, shape.payload)
				store := &memBlobs{}
				blobs := delayBlobs{store, io}
				c := fullLimitsPlugin(b, events, blobs)
				kv := delayKV{memKV{}, io}

				var raw []byte
				b.ResetTimer()
				for i := range b.N {
					store.payloads, store.byHandle = nil, nil
					var err error
					raw, err = c.Observation(b.Context(), uint64(i+1), ocrtypes.AttributedQuery{}, kv, blobs)
					require.NoError(b, err)
				}
				b.StopTimer()

				var obs pb.ConsensusQueueObservation
				require.NoError(b, proto.Unmarshal(raw, &obs))
				b.ReportMetric(float64(len(obs.Events)), "events/op")
				b.ReportMetric(float64(len(obs.BlobHandles)), "blobs/op")
				b.ReportMetric(float64(len(raw))/1024, "obsKiB/op")
			})
		}
	}
}

// BenchmarkFullValidateObservation measures validating a peer's packed observation: fetching its blobs and checking every event.
func BenchmarkFullValidateObservation(b *testing.B) {
	for _, shape := range packShapes() {
		for _, io := range ioProfiles {
			b.Run(shape.name+"/"+io.name, func(b *testing.B) {
				events := packedEvents(b, maxObservedEvents, shape.payload)
				store := &memBlobs{}
				c := fullLimitsPlugin(b, events, store)
				raw, err := c.Observation(b.Context(), 1, ocrtypes.AttributedQuery{}, memKV{}, store) // built without delays
				require.NoError(b, err)

				blobs, kv := delayBlobs{store, io}, delayKV{memKV{}, io}
				b.ResetTimer()
				for range b.N {
					require.NoError(b, c.ValidateObservation(b.Context(), 1, ocrtypes.AttributedQuery{}, ocrtypes.AttributedObservation{Observation: raw}, kv, blobs))
				}
			})
		}
	}
}

// BenchmarkFullStateTransition measures a round where every node observes maxObservedEvents events against a populated KV state.
func BenchmarkFullStateTransition(b *testing.B) {
	for _, n := range []int{4, 31} {
		for _, io := range ioProfiles {
			b.Run(fmt.Sprintf("N=%d/%s", n, io.name), func(b *testing.B) {
				f := (n - 1) / 3
				c := newSTPlugin(b)
				c.onchainCfg.N, c.onchainCfg.F = n, f
				c.limits.MaxKeyValueModifiedKeysPlusValuesBytes = ocr3_1types.MaxMaxKeyValueModifiedKeysPlusValuesBytes
				open := nodeOpts{capacity: 1 << 20, wfCaps: map[string]uint32{}, expiration: 100}
				all := packedEvents(b, 2500, 0)
				for _, e := range all {
					open.wfCaps[e.event.WorkflowID] = 1 << 20
				}
				quorum := func(events []hashedEvent, nodes int) (aos []ocrtypes.AttributedObservation) {
					for node := range n {
						var seen []hashedEvent
						if node < nodes {
							seen = events
						}
						aos = append(aos, ao(b, node, nodeObs(b, open, seen...)))
					}
					return aos
				}

				// prior state: 1000 events queued, 500 waiting at f+1 votes. Capacity 0 so that nothing executed yet.
				none := open
				none.capacity = 0
				seed := memKV{}
				setup := func(events []hashedEvent, nodes int) []ocrtypes.AttributedObservation {
					var aos []ocrtypes.AttributedObservation
					for node := range n {
						var seen []hashedEvent
						if node < nodes {
							seen = events
						}
						aos = append(aos, ao(b, node, nodeObs(b, none, seen...)))
					}
					return aos
				}
				require.NoError(b, stateTransition(b, c, seed, setup(all[:1000], n)...))
				require.NoError(b, stateTransition(b, c, seed, setup(all[1000:1500], f+1)...))

				// this round: 500 new events from everyone, the 500 waiting events gain votes, 1000 more are new
				aos := quorum(all[1000:], n)
				var written int
				b.ResetTimer()
				for range b.N {
					b.StopTimer()
					kv := &countingKV{memKV: maps.Clone(seed)}
					b.StartTimer()
					require.NoError(b, stateTransition(b, c, delayKV{kv, io}, aos...))
					written = kv.written
				}
				b.StopTimer()
				b.ReportMetric(float64(written)/1024, "kvWrittenKiB/op")
				b.ReportMetric(float64(len(aos[0].Observation))/1024, "obsKiB/node")
			})
		}
	}
}

type countingKV struct {
	memKV
	written int
}

func (c *countingKV) Write(k, v []byte) error {
	c.written += len(k) + len(v)
	return c.memKV.Write(k, v)
}

// BenchmarkFullCommit measures starting a full consensus set, which fetches its blobs (concurrently) and starts each event.
func BenchmarkFullCommit(b *testing.B) {
	for _, shape := range packShapes() {
		for _, io := range ioProfiles {
			b.Run(shape.name+"/"+io.name, func(b *testing.B) {
				events := packedEvents(b, maxObservedEvents, shape.payload)
				store := &memBlobs{}
				coord := &recCoordinator{}
				c := newCommitPlugin(b, coord, &recQueue{}, delayBlobs{store, io})
				c.limits.MaxBlobPayloadBytes = ocr3_1types.MaxMaxBlobPayloadBytes

				// the consensus set references blobs the way a StateTransition would have written it
				s := newKVState()
				stripped := make([]*pb.CoordinatedEvent, len(events))
				for i, e := range events {
					stripped[i] = proto.CloneOf(e.event)
					stripped[i].ObservedAt, stripped[i].Deadline = nil, nil
				}
				rest := stripped
				for len(rest) > 0 {
					handle, n, err := c.broadcastAsBlob(b.Context(), store, rest, 1)
					require.NoError(b, err)
					for i := range n {
						s.queuedBlobs.add(string(events[len(stripped)-len(rest)+i].hash), string(handle), uint32(i))
					}
					rest = rest[n:]
				}
				for _, e := range events {
					s.consensusEvents = append(s.consensusEvents, &pb.ConsensusQueueEvent{
						WorkflowID: e.event.WorkflowID, ExecutionID: hexID("exec", len(s.consensusEvents)),
						ObservedAt: e.event.ObservedAt, Deadline: e.event.Deadline, TimelessHash: e.hash,
					})
				}
				s.sizeOf()
				kv := memKV{}
				require.NoError(b, s.writeTo(&stateWriter{kv: kv}))

				b.ResetTimer()
				for range b.N {
					coord.mu.Lock()
					coord.started = nil
					coord.mu.Unlock()
					require.NoError(b, c.Committed(b.Context(), 1, kv))
				}
				b.StopTimer()
				b.ReportMetric(float64(len(coord.startedEvents())), "started/op")
				b.ReportMetric(float64(len(store.payloads)), "blobs/op")
			})
		}
	}
}

// observationSize must equal the real size at every step, and must not be fooled by the cached size of the observation
// itself, which goes stale as events are appended.
func TestObservationSize(t *testing.T) {
	t.Parallel()

	obs := &pb.ConsensusQueueObservation{
		BlobHandles:        [][]byte{make([]byte, 300), make([]byte, 300)},
		Expiration:         randTS(rand.New(rand.NewSource(1))),
		Capacity:           7,
		WorkflowCapacities: map[string]uint32{},
	}
	rng := rand.New(rand.NewSource(2))
	for i := range 200 {
		e := observedEventFor(rng, i)
		obs.Events = append(obs.Events, e)
		if i%3 == 0 { // some events add a workflow, with capacities of all varint sizes
			obs.WorkflowCapacities[e.WorkflowID] = uint32(rng.Int63n(1 << (1 + rng.Intn(31)))) //nolint:gosec // G115: bounded by 2^31
		}
		require.Equal(t, proto.MarshalOptions{Deterministic: true}.Size(obs), observationSize(obs), "after %d events", i+1)
	}

	t.Run("sizing the observation itself with cached sizes would be wrong", func(t *testing.T) {
		before := proto.Size(obs) // stores obs's size
		obs.Events = append(obs.Events, observedEventFor(rng, 1000))
		actual := proto.Size(proto.Clone(obs)) // a fresh message has nothing cached
		stale := proto.MarshalOptions{UseCachedSize: true}.Size(obs)
		assert.Equal(t, before, stale, "the cached size of obs is that of before the event was added")
		assert.NotEqual(t, actual, stale)
		assert.Equal(t, actual, observationSize(obs), "but the shallow copy is exact")
	})
}

// BenchmarkFullStateTransitionOverLimit measures rounds that must prune state to fit: a burst that overflows the
// observed set, and a full queue that is sent 3000 more events.
func BenchmarkFullStateTransitionOverLimit(b *testing.B) {
	const n = 4
	open := nodeOpts{capacity: 1 << 20, wfCaps: map[string]uint32{}, expiration: 100}
	none := open
	none.capacity = 0 // nothing executes, so events pile up in the queue
	all := packedEvents(b, 12000, 0)
	for _, e := range all {
		open.wfCaps[e.event.WorkflowID] = 1 << 20
		none.wfCaps[e.event.WorkflowID] = 1 << 20
	}
	round := func(b *testing.B, o nodeOpts, events []hashedEvent, nodes int) []ocrtypes.AttributedObservation {
		aos := make([]ocrtypes.AttributedObservation, 0, n)
		for node := range n {
			var seen []hashedEvent
			if node < nodes {
				seen = events
			}
			aos = append(aos, ao(b, node, nodeObs(b, o, seen...)))
		}
		return aos
	}

	run := func(b *testing.B, c *consensusQueuePlugin, seed memKV, aos []ocrtypes.AttributedObservation, metric func(memKV) float64, unit string) {
		var last memKV
		b.ResetTimer()
		for range b.N {
			b.StopTimer()
			kv := maps.Clone(seed)
			b.StartTimer()
			require.NoError(b, stateTransition(b, c, kv, aos...))
			last = kv
		}
		b.StopTimer()
		b.ReportMetric(metric(last), unit)
	}

	b.Run("observed set overflows", func(b *testing.B) {
		c := newSTPlugin(b)
		c.limits.MaxKeyValueModifiedKeysPlusValuesBytes = 100 << 10 // room for about 300 of the 1000 observed events
		aos := round(b, open, all[:1000], 2)                        // f+1 votes: observed, not queued
		run(b, c, memKV{}, aos, func(kv memKV) float64 { return float64(len(loadState(b, kv).observedEvents)) }, "kept/op")
	})

	b.Run("full queue receives 3000 more", func(b *testing.B) {
		c := newSTPlugin(b)
		seed := memKV{}
		for i := range 9 { // fill the queue to just under its value limit
			require.NoError(b, stateTransition(b, c, seed, round(b, none, all[i*1000:(i+1)*1000], n)...))
		}
		queued := func(kv memKV) float64 { return float64(len(loadState(b, kv).queuedEvents)) }
		b.ReportMetric(queued(seed), "queuedBefore")
		run(b, c, seed, round(b, none, all[9000:12000], n), queued, "queuedAfter/op")
	})
}

// A peer can pad a blob with empty events, to make finding each of its real events by scanning from the start cost a
// lot: validating 1000 events behind 2M empty ones took over half a second on 16 cores (seconds of CPU). Such a blob
// has more events than any observation can use, so it is rejected, cheaply.
func TestValidateObservation_RejectsBlobPaddedWithEmptyEvents(t *testing.T) {
	t.Parallel()

	blobs := &memBlobs{}
	c := testPlugin(t, sliceQueue{}, blobs, ocr3_1types.MaxMaxObservationBytes)

	padding := 2 * maxEventsPerBlob
	var payload []byte
	for range padding {
		payload = protowire.AppendTag(payload, coordinatedEventsFieldNumber, protowire.BytesType)
		payload = protowire.AppendBytes(payload, nil)
	}
	obs := &pb.ConsensusQueueObservation{BlobHandles: [][]byte{handleBytes(1)}, Expiration: ts(100), WorkflowCapacities: map[string]uint32{testWF: 1}}
	for i := range 10 {
		ce, qe := blobEvent(t, testWF, fmt.Sprintf("evt-%02d", i))
		qe.Deadline = ts(int64(2000 + i))
		qe.Blobs = []*pb.CoordinatedEventBlobHandle{{Blob: 0, Element: uint32(padding + i)}}
		payload = append(payload, marshalEvents(t, ce)...)
		obs.Events = append(obs.Events, qe)
	}
	blobs.store(handleBytes(1), payload)
	raw, err := proto.Marshal(obs)
	require.NoError(t, err)

	err = c.ValidateObservation(t.Context(), 1, ocrtypes.AttributedQuery{}, ocrtypes.AttributedObservation{Observation: raw}, memKV{}, blobs)
	require.ErrorContains(t, err, "more than")
}

// A node that votes for the same event round after round, with a new blob each time, must not make the stored state
// grow. Each round adds one blob handle for it, until the event is pruned or promoted, which it need not be.
func TestStateTransition_RepeatedVotesDoNotGrowState(t *testing.T) {
	t.Parallel()

	c := newSTPlugin(t) // N=4, F=1: two votes are tracked, and it takes three to promote
	kv := memKV{}
	e := mkEvents(t, 1, 2000)
	vote := func(node int, round int) ocrtypes.AttributedObservation {
		obs := nodeObs(t, stdNode, e...)
		obs.BlobHandles[0] = handleBytes(byte(10 + round)) //nolint:gosec // G115: a small test value, and a new blob each round
		return ao(t, node, obs)
	}
	none := func(node int) ocrtypes.AttributedObservation { return ao(t, node, nodeObs(t, stdNode)) }

	require.NoError(t, stateTransition(t, c, kv, vote(0, 0), vote(1, 0), none(2)))

	var sizeAfter10 int
	for round := 1; round <= 100; round++ {
		require.NoError(t, stateTransition(t, c, kv, vote(0, round), none(1), none(2)))
		if round == 10 {
			sizeAfter10 = kvTotal(kv)
		}
	}

	st := loadState(t, kv)
	require.Contains(t, st.observedEvents, string(e[0].hash), "still waiting for a third vote")
	assert.Equal(t, []byte{0, 1}, st.observedEvents[string(e[0].hash)].NodeIds)
	assert.Len(t, st.observedBlobs[string(e[0].hash)], 2, "one blob per voter")
	assert.Equal(t, blobRef{string(handleBytes(110)), 0}, st.observedBlobs[string(e[0].hash)][0], "node 0's latest")
	assert.Equal(t, sizeAfter10, kvTotal(kv), "90 more rounds of votes did not add a byte")

	// and when a third node votes, the event is promoted with just those two nodes' blobs and the new one
	require.NoError(t, stateTransition(t, c, kv, vote(2, 101), none(1), none(0)))
	cons, handles, err := readConsensusState(kv)
	require.NoError(t, err)
	require.Len(t, cons, 1)
	assert.Len(t, cons[0].Blobs, 3, "a blob for each of the three voters")
	assert.Len(t, handles, 3)
}

// gatedBlobs holds back fetches of some blobs until their gate is closed (or the fetch is cancelled).
type gatedBlobs struct {
	memBlobs
	gates map[string]chan struct{} // marshaled handle -> gate
}

func (g *gatedBlobs) FetchBlob(ctx context.Context, h ocr3_1types.BlobHandle) ([]byte, error) {
	if gate, ok := g.gates[marshalHandle(h)]; ok {
		select {
		case <-gate:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return g.memBlobs.FetchBlob(ctx, h)
}

// Any one blob of an event will do, so they are tried together: a blob that is slow or unavailable cannot hold up the
// event, however it sorts, and the attempts that are not needed must not be cancelled, as the fetch of a blob is shared
// by every event that uses it.
func TestCommitted_FetchesAnEventsBlobsConcurrently(t *testing.T) {
	t.Parallel()

	// two blobs with the same two events. X is the first of an event's blobs, and is held back.
	ceA, qeA := blobEvent(t, testWF, "a")
	ceB, qeB := blobEvent(t, testWF, "b")
	payload := marshalEvents(t, ceA, ceB)
	x, y := handleBytes(1), handleBytes(2)
	// A has both blobs, X first. B has only X.
	qeA.Blobs = []*pb.CoordinatedEventBlobHandle{{Blob: 0, Element: 0}, {Blob: 1, Element: 0}}
	qeB.Blobs = []*pb.CoordinatedEventBlobHandle{{Blob: 0, Element: 1}}

	gate := make(chan struct{})
	var once sync.Once
	release := func() { once.Do(func() { close(gate) }) }
	t.Cleanup(release)
	blobs := &gatedBlobs{}
	blobs.store(x, payload)
	blobs.store(y, payload)
	var hx ocr3_1types.BlobHandle
	require.NoError(t, hx.UnmarshalBinary(x))
	blobs.gates = map[string]chan struct{}{marshalHandle(hx): gate}

	events, err := proto.Marshal(&pb.ConsensusQueueEventList{Events: []*pb.ConsensusQueueEvent{qeA, qeB}})
	require.NoError(t, err)
	handles, err := proto.Marshal(&pb.ConsensusQueueBlobHandleList{BlobHandles: [][]byte{x, y}})
	require.NoError(t, err)
	kv := memKV{consensusEventsKey: events, blobHandlesKey: handles}

	coord := &recCoordinator{}
	c := newCommitPlugin(t, coord, &recQueue{}, blobs)
	committed := make(chan error, 1)
	go func() { committed <- c.Committed(t.Context(), 1, kv) }()

	// A is started through Y, while X is still held back, and so is not waited for
	require.Eventually(t, func() bool { return len(coord.startedEvents()) == 1 }, 5*time.Second, 5*time.Millisecond,
		"an event is not held up by its other, unavailable blob")
	assert.Equal(t, "a", coord.startedEvents()[0].Event.Event.ID)
	select {
	case err := <-committed:
		t.Fatalf("B can only be read from X, so committing waits for it, but it returned %v (started %d)", err, len(coord.startedEvents()))
	default:
	}

	// X was the loser for A, and must not have been cancelled for it: B needs it
	release()
	select {
	case err := <-committed:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Committed did not finish")
	}
	started := coord.startedEvents()
	require.Len(t, started, 2)
	ids := []string{started[0].Event.Event.ID, started[1].Event.Event.ID}
	assert.ElementsMatch(t, []string{"a", "b"}, ids, "B was started from X after it was released")
}

// removalQueue is a queue that serves events, and records those removed from it.
type removalQueue struct {
	sliceQueue
	mu      sync.Mutex
	removed []string
}

func (q *removalQueue) Remove(_ context.Context, executionIDs iter.Seq[string]) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.removed = slices.AppendSeq(q.removed, executionIDs)
}

func (q *removalQueue) removedIDs() []string {
	q.mu.Lock()
	defer q.mu.Unlock()
	return slices.Clone(q.removed)
}

// drainSwitch is a KV drain mode that can be flipped while the plugin runs.
type drainSwitch struct{ on atomic.Bool }

func (d *drainSwitch) check(context.Context, time.Time) (bool, error) { return d.on.Load(), nil }

func TestKVDrainActive(t *testing.T) {
	t.Parallel()
	active := func(settingsJSON string, at time.Time) (bool, error) {
		var getter settings.Getter
		if settingsJSON != "" {
			var err error
			getter, err = settings.NewJSONGetter([]byte(settingsJSON))
			require.NoError(t, err)
		}
		return KVDrainActive(limits.Factory{Logger: logger.Test(t), Settings: getter})(t.Context(), at)
	}
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	window := fmt.Sprintf(`{"global":{"FeatureConsensusQueueKVDrainActivePeriod":"[%d,%d]"}}`, start.Unix(), start.Add(time.Hour).Unix())

	on, err := active("", start)
	require.NoError(t, err)
	assert.False(t, on, "off by default")

	for at, want := range map[time.Time]bool{
		start.Add(-time.Second): false, start: true, start.Add(time.Hour): true, start.Add(time.Hour + time.Second): false,
	} {
		on, err = active(window, at)
		require.NoError(t, err)
		assert.Equal(t, want, on, at.String())
	}
}

func TestObservation_KVDrain(t *testing.T) {
	t.Parallel()

	events := mkEvents(t, 3, time.Now().Add(time.Hour).Unix())
	q := &removalQueue{sliceQueue: sliceQueue{events: events}}
	blobs := &memBlobs{}
	c := testPlugin(t, q.sliceQueue, blobs, 1<<19)
	c.triggerQueue = q
	drain := &drainSwitch{}
	c.kvDrain = drain.check
	kv := memKV{}

	drain.on.Store(true)
	obs, err := observe(t, c, kv, blobs)
	require.NoError(t, err)
	assert.Empty(t, obs.Events, "no new events are reported")
	assert.Empty(t, obs.BlobHandles, "so none are broadcast")
	assert.True(t, obs.KvDrainActive, "and the mode is reported, for consensus")
	assert.NotNil(t, obs.Expiration, "though what consensus needs is still reported")
	assert.Len(t, q.removedIDs(), 3, "the local queue is emptied")
	assert.Empty(t, blobs.payloads, "nothing was broadcast")

	// a flip takes effect on the next round
	drain.on.Store(false)
	obs, err = observe(t, c, kv, blobs)
	require.NoError(t, err)
	assert.Len(t, obs.Events, 3, "events are observed again")
	assert.False(t, obs.KvDrainActive)
}

func TestStateTransition_KVDrain(t *testing.T) {
	t.Parallel()

	c := newSTPlugin(t)
	kv := memKV{}
	events := mkEvents(t, 4, 2000)
	all := func(o nodeOpts, ev ...hashedEvent) []ocrtypes.AttributedObservation {
		return []ocrtypes.AttributedObservation{ao(t, 0, nodeObs(t, o, ev...)), ao(t, 1, nodeObs(t, o, ev...)), ao(t, 2, nodeObs(t, o, ev...))}
	}

	// no capacity: all four wait in the queue
	noCapacity := nodeOpts{capacity: 0, wfCaps: map[string]uint32{testWF: 0}, expiration: 100}
	require.NoError(t, stateTransition(t, c, kv, all(noCapacity, events...)...))
	require.Len(t, loadState(t, kv).queuedEvents, 4)

	// drain: still no capacity, and a faulty node still reports new events, but nothing is added and all queued are taken
	late := mkEvents(t, 2, 3000)
	for i := range late {
		_, late[i].event = newCoordinatedEvent(t, testWF, fmt.Sprintf("late-%d", i))
		late[i].event.Deadline = ts(3000 + int64(i))
	}
	draining := noCapacity
	draining.drain = true
	require.NoError(t, stateTransition(t, c, kv,
		ao(t, 0, nodeObs(t, draining)), ao(t, 1, nodeObs(t, draining)), ao(t, 2, nodeObs(t, noCapacity, late...))))
	state := loadState(t, kv)
	assert.Empty(t, state.queuedEvents, "every queued event was taken")
	assert.Empty(t, state.observedEvents, "and nothing new was added")
	assert.Len(t, consensusIDs(t, kv), 4, "to execute")

	// the next round without the flag adds events again
	require.NoError(t, stateTransition(t, c, kv, all(noCapacity, late...)...))
	assert.Len(t, loadState(t, kv).queuedEvents, 2)
}

// KV drain mode takes the votes of F+1 nodes: at least one honest. N=4, F=1.
func TestStateTransition_KVDrainNeedsFPlusOneVotes(t *testing.T) {
	t.Parallel()

	for votes, wantDrain := range map[int]bool{0: false, 1: false, 2: true, 3: true} {
		t.Run(fmt.Sprintf("%d of 3 vote", votes), func(t *testing.T) {
			t.Parallel()
			c := newSTPlugin(t)
			kv := memKV{}
			noCapacity := nodeOpts{capacity: 0, wfCaps: map[string]uint32{testWF: 0}, expiration: 100}
			events := mkEvents(t, 2, 2000)
			require.NoError(t, stateTransition(t, c, kv,
				ao(t, 0, nodeObs(t, noCapacity, events...)), ao(t, 1, nodeObs(t, noCapacity, events...)), ao(t, 2, nodeObs(t, noCapacity, events...))))
			require.Len(t, loadState(t, kv).queuedEvents, 2)

			aos := make([]ocrtypes.AttributedObservation, 0, 3)
			for i := range 3 {
				o := noCapacity
				o.drain = i < votes
				aos = append(aos, ao(t, i, nodeObs(t, o)))
			}
			require.NoError(t, stateTransition(t, c, kv, aos...))
			if wantDrain {
				assert.Len(t, consensusIDs(t, kv), 2)
			} else {
				assert.Empty(t, consensusIDs(t, kv), "no capacity, so nothing is taken")
			}
		})
	}
}

func TestStateTransition_KVDrainDropsExpired(t *testing.T) {
	t.Parallel()

	c := newSTPlugin(t)
	kv := memKV{}
	events := mkEvents(t, 3, 2000) // deadlines 2000, 2001, 2002
	noCapacity := nodeOpts{capacity: 0, wfCaps: map[string]uint32{testWF: 0}, expiration: 100}
	require.NoError(t, stateTransition(t, c, kv,
		ao(t, 0, nodeObs(t, noCapacity, events...)), ao(t, 1, nodeObs(t, noCapacity, events...)), ao(t, 2, nodeObs(t, noCapacity, events...))))

	expired := nodeOpts{capacity: 0, wfCaps: map[string]uint32{testWF: 0}, expiration: 2001, drain: true} // only the first has expired
	require.NoError(t, stateTransition(t, c, kv,
		ao(t, 0, nodeObs(t, expired)), ao(t, 1, nodeObs(t, expired)), ao(t, 2, nodeObs(t, expired))))
	assert.Equal(t, []string{execIDOf(t, "evt-01"), execIDOf(t, "evt-02")}, consensusIDs(t, kv), "the expired event is dropped, not taken")
	assert.Empty(t, loadState(t, kv).queuedEvents)
}

func TestValidateObservation_KVDrain(t *testing.T) {
	t.Parallel()

	c := newSTPlugin(t)
	validate := func(obs *pb.ConsensusQueueObservation) error {
		return c.ValidateObservation(t.Context(), 1, ocrtypes.AttributedQuery{}, ao(t, 1, obs), memKV{}, &memBlobs{})
	}
	drain := nodeOpts{capacity: 3, wfCaps: map[string]uint32{}, expiration: 100, drain: true}

	empty := nodeObs(t, drain)
	empty.BlobHandles = nil // as the helper always references a shared blob
	require.NoError(t, validate(empty), "no events and no blob handles")

	withEvents := nodeObs(t, drain, mkEvents(t, 1, 2000)...)
	require.ErrorContains(t, validate(withEvents), "KV drain mode reports no new events")

	withHandles := nodeObs(t, drain)
	withHandles.Events = nil
	withHandles.BlobHandles = [][]byte{handleBytes(1)}
	assert.ErrorContains(t, validate(withHandles), "KV drain mode reports no new events")
}
