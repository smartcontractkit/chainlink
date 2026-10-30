package v2_test

import (
	"context"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
)

const (
	droppedTotalMetric  = "platform_engine_trigger_event_dropped_total"
	enqueuedTotalMetric = "platform_engine_trigger_event_enqueued_total"
)

// TestEngine_DropReasonParity asserts that the legacy and the coordinated engine drop a trigger
// event for the same reason. On the legacy path the engine admits events, on the coordinated
// path the trigger coordinator does. Each scenario sends events through the real delivery path
// and checks platform_engine_trigger_event_dropped_total by drop reason. The engine and the
// coordinator record on one test meter, so the test does not depend on which of them counts
// the drop.
func TestEngine_DropReasonParity(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name string
		run  func(t *testing.T, mode lifecycleEngineMode, meter *dropMeter)
	}{
		{"QueueTimeoutUnreadable", testDropQueueTimeoutUnreadable},
		{"Expired", testDropExpired},
		{"GlobalPoolSaturated_QueueFull", testDropGlobalPoolSaturated},
		{"ShardDenied", testDropShardDenied},
	}
	for _, mode := range []lifecycleEngineMode{lifecycleLegacy, lifecycleCoordinated} {
		t.Run(string(mode), func(t *testing.T) {
			t.Parallel()
			for _, s := range scenarios {
				t.Run(s.name, func(t *testing.T) {
					t.Parallel()
					s.run(t, mode, newDropMeter(t))
				})
			}
		})
	}
}

// testDropQueueTimeoutUnreadable: when the queue timeout has no usable value, the event's
// deadline cannot be stamped, so the event is dropped before it is queued.
func testDropQueueTimeoutUnreadable(t *testing.T, mode lifecycleEngineMode, meter *dropMeter) {
	limiters := newSharedLimiters(t, `{}`)
	limiters.TriggerEventQueueTimeout = &failAfterNBoundLimiter[time.Duration]{
		ok:      time.Minute,
		failErr: limits.ErrMissingTenant{Scope: settings.ScopeWorkflow},
	}
	h := newEngineLifecycleHarness(t, mode, withLimiters(limiters), withMeter(meter.meter))
	h.start()

	h.deliver(testutils.WaitTimeout(t))

	meter.requireDrops(t, map[string]int64{monitoring.TriggerDropReasonQueueAgeLimitReadFailed: 1})
	assert.Zero(t, h.moduleExecutions.Load(), "a dropped event never reaches the module")
}

// testDropExpired: an event that waits in the queue past its deadline is dropped when it is
// taken off the queue, and never executes.
func testDropExpired(t *testing.T, mode lifecycleEngineMode, meter *dropMeter) {
	clock := clockwork.NewFakeClock()
	limiters := newSharedLimiters(t, `{"global":{"WorkflowExecutionConcurrencyLimit":"1"}}`)
	h := newEngineLifecycleHarness(t, mode, withLimiters(limiters), withMeter(meter.meter), withClock(clock))
	h.start()

	// The first event holds the only slot. The second leaves the queue, passing its deadline
	// check, and waits for the slot. The third stays in the queue.
	h.startExecution()
	h.deliver(testutils.WaitTimeout(t))
	require.Eventually(t, func() bool { return h.slots.waits.Load() == 2 }, testutils.WaitTimeout(t), 5*time.Millisecond,
		"the second event waits for the slot")
	h.deliver(testutils.WaitTimeout(t))
	require.Eventually(t, func() bool { return meter.total(t, enqueuedTotalMetric) == 3 }, testutils.WaitTimeout(t), 5*time.Millisecond,
		"the third event is queued")

	clock.Advance(cresettings.Default.PerWorkflow.TriggerEventQueueTimeout.DefaultValue + time.Second)
	h.blocker.release()

	meter.requireDrops(t, map[string]int64{monitoring.TriggerDropReasonExpired: 1})
	require.Eventually(t, func() bool { return h.moduleExecutions.Load() == 2 }, testutils.WaitTimeout(t), 5*time.Millisecond,
		"the events that passed their deadline check execute")
}

// testDropGlobalPoolSaturated: with the global concurrency pool full and the queue full, the
// reader keeps reading the trigger's event channel, and every extra event is dropped as
// enqueue_queue_full. A reader that stopped would fail h.deliver.
func testDropGlobalPoolSaturated(t *testing.T, mode lifecycleEngineMode, meter *dropMeter) {
	const overflow = 3
	limiters := newSharedLimiters(t,
		`{"global":{"WorkflowExecutionConcurrencyLimit":"1","PerWorkflow":{"TriggerEventQueueLimit":"1"}}}`)
	h := newEngineLifecycleHarness(t, mode, withLimiters(limiters), withMeter(meter.meter))
	h.start()

	// The first event holds the only global slot. The second leaves the queue and waits for
	// the slot. The third fills the queue.
	h.startExecution()
	h.deliver(testutils.WaitTimeout(t))
	require.Eventually(t, func() bool { return h.slots.waits.Load() == 2 }, testutils.WaitTimeout(t), 5*time.Millisecond,
		"the second event waits for the slot")
	h.deliver(testutils.WaitTimeout(t))

	for range overflow {
		h.deliver(testutils.WaitTimeout(t))
	}
	meter.requireDrops(t, map[string]int64{monitoring.TriggerDropReasonEnqueueQueueFull: overflow})

	if h.coordinator != nil {
		got, err := h.coordinator.AvailableCapacity(t.Context())
		require.NoError(t, err)
		assert.Zero(t, got.Global, "the capacity report shows the global pool full")
	}

	h.blocker.release()
	require.Eventually(t, func() bool { return h.moduleExecutions.Load() == 3 }, testutils.WaitTimeout(t), 5*time.Millisecond,
		"the events that were not dropped execute")
	meter.requireDrops(t, map[string]int64{monitoring.TriggerDropReasonEnqueueQueueFull: overflow})
}

// testDropShardDenied: an event the shard ownership check denies is ACKed, dropped as
// shard_denied_not_owner and counted on the shard-denial counter.
func testDropShardDenied(t *testing.T, mode lifecycleEngineMode, meter *dropMeter) {
	h := newEngineLifecycleHarness(t, mode, withMeter(meter.meter),
		withAdmission(func(context.Context, triggers.CoordinatedEvent) error { return v2.ErrShardDeniedNotOwner }))
	h.start()

	h.deliver(testutils.WaitTimeout(t))

	meter.requireDrops(t, map[string]int64{monitoring.TriggerDropReasonShardDeniedNotOwner: 1})
	assert.Equal(t, int64(1), meter.total(t, "platform_engine_shard_execution_denied_not_owner"))
	assert.Equal(t, int32(1), h.ackCalls.Load(), "a denied event is ACKed")
	assert.Zero(t, h.moduleExecutions.Load(), "a denied event never reaches the module")
}

// dropMeter is a test meter that the engine and the coordinator both record on.
type dropMeter struct {
	meter  metric.Meter
	reader *sdkmetric.ManualReader
}

func newDropMeter(t *testing.T) *dropMeter {
	meter, reader := newTestMeter(t)
	return &dropMeter{meter: meter, reader: reader}
}

// total sums every data point of the named counter.
func (m *dropMeter) total(t *testing.T, name string) int64 {
	var total int64
	for _, dp := range m.dataPoints(t, name) {
		total += dp.Value
	}
	return total
}

// drops returns platform_engine_trigger_event_dropped_total by drop reason.
func (m *dropMeter) drops(t *testing.T) map[string]int64 {
	got := map[string]int64{}
	for _, dp := range m.dataPoints(t, droppedTotalMetric) {
		reason, _ := dp.Attributes.Value(attribute.Key("dropReason"))
		got[reason.AsString()] += dp.Value
	}
	return got
}

// requireDrops waits until the drop counts by reason are exactly want.
func (m *dropMeter) requireDrops(t *testing.T, want map[string]int64) {
	t.Helper()
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		assert.Equal(c, want, m.drops(t))
	}, testutils.WaitTimeout(t), 5*time.Millisecond, "drops by reason")
}

func (m *dropMeter) dataPoints(t *testing.T, name string) []metricdata.DataPoint[int64] {
	var rm metricdata.ResourceMetrics
	require.NoError(t, m.reader.Collect(context.Background(), &rm))
	var points []metricdata.DataPoint[int64]
	for _, sm := range rm.ScopeMetrics {
		for _, md := range sm.Metrics {
			if md.Name != name {
				continue
			}
			data, ok := md.Data.(metricdata.Sum[int64])
			require.True(t, ok, "expected Sum[int64] for %s, got %T", name, md.Data)
			points = append(points, data.DataPoints...)
		}
	}
	return points
}
