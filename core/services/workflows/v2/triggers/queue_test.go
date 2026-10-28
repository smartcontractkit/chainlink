package triggers

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

const queueFull = monitoring.CentralTriggerQueueDropReasonQueueFull

// fakeLimit is a BoundLimiter[int] whose limit and read error can change between calls.
type fakeLimit struct {
	mu    sync.Mutex
	limit int
	err   error
}

func (l *fakeLimit) Limit(context.Context) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.limit, l.err
}

func (l *fakeLimit) Check(context.Context, int) error { return nil }
func (l *fakeLimit) Close() error                     { return nil }

func (l *fakeLimit) set(limit int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.limit, l.err = limit, err
}

func newTestCentralQueue(t *testing.T, limit *fakeLimit) *centralQueue {
	t.Helper()
	qm, _ := testQueueMetrics(t)
	return NewQueue(limit, qm, logger.Test(t))
}

// baseLabel is the only label the queue's metrics may carry, apart from drop_reason on the dropped counter.
const baseLabel = "test_label"

// testQueueMetrics returns a labeler with one base label, backed by a ManualReader so emitted values can be
// read back.
func testQueueMetrics(t *testing.T) (*monitoring.CentralTriggerQueueMetricLabeler, *sdkmetric.ManualReader) {
	t.Helper()
	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { _ = mp.Shutdown(context.Background()) })
	qm, err := monitoring.NewCentralTriggerQueueMetrics(mp.Meter("test"))
	require.NoError(t, err)
	return monitoring.NewCentralTriggerQueueMetricLabeler(metrics.NewLabeler().With(baseLabel, "base"), qm), reader
}

// centralEvent returns an event for the given registration. Deadline and SequenceNumber are set so a test
// can tell if the queue changes them.
func centralEvent(workflowID string, triggerIndex int, id string) CoordinatedEvent {
	observed := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	return CoordinatedEvent{
		WorkflowID:     workflowID,
		TriggerCapID:   fmt.Sprintf("trigger-%d@1.0.0", triggerIndex),
		TriggerIndex:   triggerIndex,
		ObservedAt:     observed,
		Deadline:       observed.Add(time.Minute),
		SequenceNumber: 42,
		Event:          capabilities.TriggerResponse{Event: capabilities.TriggerEvent{ID: id}},
	}
}

func putAll(t *testing.T, q *centralQueue, events ...CoordinatedEvent) {
	t.Helper()
	for _, e := range events {
		require.NoError(t, q.Put(t.Context(), e))
	}
}

func drainIDs(t *testing.T, q *centralQueue) []string {
	t.Helper()
	events, err := q.Drain(t.Context())
	require.NoError(t, err)
	ids := make([]string, 0, len(events))
	for _, e := range events {
		ids = append(ids, e.Event.Event.ID)
	}
	return ids
}

// TestQueue_PutDrain asserts events come back unchanged and in insertion order. Drain empties the queue and
// later puts do not change a slice it already returned.
func TestQueue_PutDrain(t *testing.T) {
	t.Parallel()
	q := newTestCentralQueue(t, &fakeLimit{limit: 10})

	want := []CoordinatedEvent{
		centralEvent("wf-a", 0, "evt-1"),
		centralEvent("wf-b", 0, "evt-2"),
		centralEvent("wf-a", 1, "evt-3"),
	}
	putAll(t, q, want...)

	got, err := q.Drain(t.Context())
	require.NoError(t, err)
	assert.Equal(t, want, got)

	putAll(t, q, centralEvent("wf-c", 0, "evt-4"))
	assert.Equal(t, want, got, "a later put must not write into a drained slice")

	assert.Equal(t, []string{"evt-4"}, drainIDs(t, q))
	assert.Empty(t, drainIDs(t, q))
}

// TestQueue_Capacity asserts how the limit, a changed limit, and a failed limit read decide when the queue
// is full.
func TestQueue_Capacity(t *testing.T) {
	t.Parallel()
	defaultLimit := cresettings.Default.CentralTriggerEventQueueLimit.DefaultValue

	t.Run("a put at the limit is rejected and not stored", func(t *testing.T) {
		t.Parallel()
		q := newTestCentralQueue(t, &fakeLimit{limit: 2})
		putAll(t, q, centralEvent("wf-a", 0, "evt-1"), centralEvent("wf-a", 0, "evt-2"))

		require.ErrorIs(t, q.Put(t.Context(), centralEvent("wf-a", 0, "evt-3")), ErrQueueFull)
		assert.Equal(t, []string{"evt-1", "evt-2"}, drainIDs(t, q))
	})

	t.Run("lowering the limit keeps queued events and rejects new ones", func(t *testing.T) {
		t.Parallel()
		limit := &fakeLimit{limit: 3}
		q := newTestCentralQueue(t, limit)
		putAll(t, q, centralEvent("wf-a", 0, "evt-1"), centralEvent("wf-a", 0, "evt-2"), centralEvent("wf-a", 0, "evt-3"))

		limit.set(1, nil)
		require.ErrorIs(t, q.Put(t.Context(), centralEvent("wf-a", 0, "evt-4")), ErrQueueFull)
		assert.Equal(t, []string{"evt-1", "evt-2", "evt-3"}, drainIDs(t, q))
	})

	t.Run("raising the limit accepts puts again", func(t *testing.T) {
		t.Parallel()
		limit := &fakeLimit{limit: 1}
		q := newTestCentralQueue(t, limit)
		putAll(t, q, centralEvent("wf-a", 0, "evt-1"))
		require.ErrorIs(t, q.Put(t.Context(), centralEvent("wf-a", 0, "evt-2")), ErrQueueFull)

		limit.set(2, nil)
		putAll(t, q, centralEvent("wf-a", 0, "evt-3"))
		assert.Equal(t, []string{"evt-1", "evt-3"}, drainIDs(t, q))
	})

	t.Run("a failed limit read uses the default", func(t *testing.T) {
		t.Parallel()
		// The value returned with the error is ignored.
		q := newTestCentralQueue(t, &fakeLimit{limit: 1, err: assert.AnError})
		putAll(t, q, centralEvent("wf-a", 0, "evt-1"), centralEvent("wf-a", 0, "evt-2"))

		stats, err := q.Stats(t.Context())
		require.NoError(t, err)
		assert.EqualValues(t, defaultLimit, stats.Capacity)
		assert.Equal(t, uint64(2), stats.Len)
	})

	t.Run("a failed limit read with more queued than the default reports full", func(t *testing.T) {
		t.Parallel()
		limit := &fakeLimit{limit: defaultLimit + 5}
		q := newTestCentralQueue(t, limit)
		for i := range defaultLimit + 5 {
			putAll(t, q, centralEvent("wf-a", 0, fmt.Sprintf("evt-%d", i)))
		}

		limit.set(0, assert.AnError)
		require.ErrorIs(t, q.Put(t.Context(), centralEvent("wf-a", 0, "one-more")), ErrQueueFull)

		stats, err := q.Stats(t.Context())
		require.NoError(t, err)
		assert.EqualValues(t, defaultLimit+5, stats.Len)
		assert.EqualValues(t, defaultLimit+5, stats.Capacity, "capacity never reads below the current length")
	})
}

// TestQueue_Stats asserts depth by workflow and trigger follows puts and drains. Drop counts accumulate per
// reason across drains and are not split by workflow, so they stay bounded. The returned maps are copies.
func TestQueue_Stats(t *testing.T) {
	t.Parallel()
	limit := &fakeLimit{limit: 3}
	q := newTestCentralQueue(t, limit)

	putAll(t, q, centralEvent("wf-a", 0, "evt-1"), centralEvent("wf-a", 1, "evt-2"), centralEvent("wf-b", 0, "evt-3"))
	require.ErrorIs(t, q.Put(t.Context(), centralEvent("wf-b", 0, "evt-4")), ErrQueueFull)
	require.ErrorIs(t, q.Put(t.Context(), centralEvent("wf-a", 1, "evt-5")), ErrQueueFull)

	stats, err := q.Stats(t.Context())
	require.NoError(t, err)
	assert.Equal(t, uint64(3), stats.Len)
	assert.Equal(t, uint64(3), stats.Capacity)
	assert.Equal(t, map[string]uint64{"wf-a": 2, "wf-b": 1}, stats.DepthByWorkflow)
	assert.Equal(t, map[string]uint64{"trigger-0@1.0.0": 2, "trigger-1@1.0.0": 1}, stats.DepthByTriggerType)
	assert.Equal(t, map[string]uint64{queueFull: 2}, stats.DroppedByReason)

	// Changing the returned maps must not change the queue.
	want, err := q.Stats(t.Context())
	require.NoError(t, err)
	stats.DepthByWorkflow["wf-a"] = 99
	stats.DepthByTriggerType["trigger-0@1.0.0"] = 99
	stats.DroppedByReason[queueFull] = 99
	stats.DroppedByReason["other"] = 1
	got, err := q.Stats(t.Context())
	require.NoError(t, err)
	assert.Equal(t, want, got)

	drainIDs(t, q)
	limit.set(1, nil)
	putAll(t, q, centralEvent("wf-c", 0, "evt-6"))
	require.ErrorIs(t, q.Put(t.Context(), centralEvent("wf-b", 0, "evt-7")), ErrQueueFull)

	stats, err = q.Stats(t.Context())
	require.NoError(t, err)
	assert.Equal(t, uint64(1), stats.Len)
	assert.Equal(t, map[string]uint64{"wf-c": 1}, stats.DepthByWorkflow)
	assert.Equal(t, map[string]uint64{"trigger-0@1.0.0": 1}, stats.DepthByTriggerType)
	assert.Equal(t, map[string]uint64{queueFull: 3}, stats.DroppedByReason)
}

// TestQueue_ConcurrentPutDrain asserts that with many producers, one consumer and a concurrent Stats reader no
// event is lost or duplicated, and each producer's events keep their order. Run with -race.
func TestQueue_ConcurrentPutDrain(t *testing.T) {
	t.Parallel()
	const (
		producers   = 8
		perProducer = 500
	)
	tests := []struct {
		name  string
		limit int
		// drainDuringPuts runs the consumer alongside the producers. Otherwise it starts once they finish,
		// so drops are certain.
		drainDuringPuts bool
	}{
		{name: "consumer drains during puts", limit: 100, drainDuringPuts: true},
		{name: "consumer starts after puts", limit: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := newTestCentralQueue(t, &fakeLimit{limit: tt.limit})

			var accepted, dropped [producers]int
			var wg sync.WaitGroup
			for p := range producers {
				wg.Go(func() {
					for seq := range perProducer {
						e := centralEvent(fmt.Sprintf("wf-%d", p), 0, "evt")
						e.SequenceNumber = uint64(seq)
						err := q.Put(t.Context(), e)
						if err == nil {
							accepted[p]++
							continue
						}
						assert.ErrorIs(t, err, ErrQueueFull)
						dropped[p]++
					}
				})
			}

			done := make(chan struct{})
			go func() {
				wg.Wait()
				close(done)
			}()

			// Stats runs alongside, and walks every map, so -race catches a Stats that returns internal maps.
			var readers sync.WaitGroup
			readers.Go(func() {
				for {
					select {
					case <-done:
						return
					default:
					}
					stats, err := q.Stats(t.Context())
					if err != nil {
						t.Errorf("Stats() error = %v", err)
						return
					}
					walkStats(stats)
				}
			})

			if !tt.drainDuringPuts {
				<-done
			}
			var drained []CoordinatedEvent
			for finished := false; !finished; {
				select {
				case <-done:
					finished = true
				default:
				}
				events, err := q.Drain(t.Context())
				require.NoError(t, err)
				drained = append(drained, events...)
			}
			readers.Wait()

			var totalAccepted, totalDropped int
			for p := range producers {
				assert.Equal(t, perProducer, accepted[p]+dropped[p], "producer %d", p)
				totalAccepted += accepted[p]
				totalDropped += dropped[p]
			}
			assert.Len(t, drained, totalAccepted)
			if !tt.drainDuringPuts {
				assert.Positive(t, totalDropped, "the drop path never ran")
			}

			last := map[string]uint64{}
			for _, e := range drained {
				prev, seen := last[e.WorkflowID]
				if seen {
					assert.Greater(t, e.SequenceNumber, prev, "order for %s", e.WorkflowID)
				}
				last[e.WorkflowID] = e.SequenceNumber
			}

			stats, err := q.Stats(t.Context())
			require.NoError(t, err)
			assert.Equal(t, uint64(totalDropped), stats.DroppedByReason[queueFull])
		})
	}
}

// walkStats reads every value in the maps of stats and returns their sum.
func walkStats(stats Stats) uint64 {
	var sum uint64
	for _, n := range stats.DepthByWorkflow {
		sum += n
	}
	for _, n := range stats.DepthByTriggerType {
		sum += n
	}
	for _, n := range stats.DroppedByReason {
		sum += n
	}
	return sum
}

// TestQueue_Metrics asserts the queue's counters and gauges have the expected values and match Stats. Drain
// records depth as well as Put. Every data point carries only the base label, plus drop_reason on the dropped
// counter.
func TestQueue_Metrics(t *testing.T) {
	t.Parallel()
	qm, reader := testQueueMetrics(t)
	q := NewQueue(&fakeLimit{limit: 2}, qm, logger.Test(t))

	putAll(t, q, centralEvent("wf-a", 0, "evt-1"), centralEvent("wf-b", 0, "evt-2"))
	require.ErrorIs(t, q.Put(t.Context(), centralEvent("wf-a", 0, "evt-3")), ErrQueueFull)
	drainIDs(t, q)
	assert.Equal(t, int64(0), int64Gauge(t, reader, "platform_engine_central_trigger_queue_depth"), "after drain")
	putAll(t, q, centralEvent("wf-a", 0, "evt-4"))

	stats, err := q.Stats(t.Context())
	require.NoError(t, err)

	assert.Equal(t, int64(3), int64Sum(t, reader, "platform_engine_central_trigger_queue_puts_total"))
	dropped := int64Sum(t, reader, "platform_engine_central_trigger_queue_dropped_events_total")
	assert.Equal(t, int64(1), dropped)
	assert.EqualValues(t, stats.DroppedByReason[queueFull], dropped)
	assert.Equal(t, int64(2), int64Sum(t, reader, "platform_engine_central_trigger_queue_drained_events_total"))
	depth := int64Gauge(t, reader, "platform_engine_central_trigger_queue_depth")
	assert.Equal(t, int64(1), depth)
	assert.EqualValues(t, stats.Len, depth)
	capacity := int64Gauge(t, reader, "platform_engine_central_trigger_queue_capacity")
	assert.Equal(t, int64(2), capacity)
	assert.EqualValues(t, stats.Capacity, capacity)

	droppedMetric := findMetric(t, reader, "platform_engine_central_trigger_queue_dropped_events_total")
	require.NotNil(t, droppedMetric)
	droppedSum, ok := droppedMetric.Data.(metricdata.Sum[int64])
	require.True(t, ok)
	require.Len(t, droppedSum.DataPoints, 1)
	reason, ok := droppedSum.DataPoints[0].Attributes.Value(attribute.Key(platform.KeyTriggerDropReason))
	require.True(t, ok)
	assert.Equal(t, queueFull, reason.AsString())

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	require.NotEmpty(t, rm.ScopeMetrics)
	for _, m := range rm.ScopeMetrics[0].Metrics {
		for _, attrs := range dataPointAttributes(t, m) {
			for _, kv := range attrs.ToSlice() {
				allowed := kv.Key == baseLabel || kv.Key == attribute.Key(platform.KeyTriggerDropReason)
				assert.True(t, allowed, "%s carries label %s", m.Name, kv.Key)
			}
			assert.True(t, attrs.HasValue(baseLabel), "%s is missing the base label", m.Name)
		}
	}
}

// int64Sum returns the total of the named Int64 counter across its data points.
func int64Sum(t *testing.T, reader *sdkmetric.ManualReader, name string) int64 {
	t.Helper()
	m := findMetric(t, reader, name)
	require.NotNil(t, m, name)
	sum, ok := m.Data.(metricdata.Sum[int64])
	require.True(t, ok, "%s is %T", name, m.Data)
	var total int64
	for _, dp := range sum.DataPoints {
		total += dp.Value
	}
	return total
}

// int64Gauge returns the value of the named Int64 gauge, which must have one data point.
func int64Gauge(t *testing.T, reader *sdkmetric.ManualReader, name string) int64 {
	t.Helper()
	m := findMetric(t, reader, name)
	require.NotNil(t, m, name)
	gauge, ok := m.Data.(metricdata.Gauge[int64])
	require.True(t, ok, "%s is %T", name, m.Data)
	require.Len(t, gauge.DataPoints, 1, name)
	return gauge.DataPoints[0].Value
}

func dataPointAttributes(t *testing.T, m metricdata.Metrics) []attribute.Set {
	t.Helper()
	var sets []attribute.Set
	switch data := m.Data.(type) {
	case metricdata.Sum[int64]:
		for _, dp := range data.DataPoints {
			sets = append(sets, dp.Attributes)
		}
	case metricdata.Gauge[int64]:
		for _, dp := range data.DataPoints {
			sets = append(sets, dp.Attributes)
		}
	default:
		t.Fatalf("unexpected data type %T for %s", m.Data, m.Name)
	}
	return sets
}
