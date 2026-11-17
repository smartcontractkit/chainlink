package v2_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"
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
	}
	for _, mode := range []lifecycleEngineMode{lifecycleLegacy, lifecycleCoordinated} {
		for _, s := range scenarios {
			t.Run(string(mode)+"/"+s.name, func(t *testing.T) {
				t.Parallel()
				s.run(t, mode, newDropMeter(t))
			})
		}
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

// TestEngine_AdmissionDenialParity asserts that the legacy and the coordinated engine handle an
// event the admission check denies the same way: it is ACKed, never executes, emits no execution
// events and is dropped by its denial reason. A shard denial also counts on its shard-denial
// counter. On the legacy path the engine handles the denial, on the coordinated path the trigger
// coordinator does.
func TestEngine_AdmissionDenialParity(t *testing.T) {
	t.Parallel()

	const (
		shardNotOwnerTotal = "platform_engine_shard_execution_denied_not_owner"
		shardOrchErrTotal  = "platform_engine_shard_execution_denied_orchestrator_error"
		// A denial is handled as soon as the event is read, so a failure shows within seconds
		// instead of waiting out the test deadline.
		denialWait = 5 * time.Second
	)

	tests := []struct {
		name           string
		admissionErr   error
		wantDropReason string
		wantNotOwner   int64
		wantOrchErr    int64
	}{
		{
			name:           "shard denial not owner",
			admissionErr:   v2.ErrShardDeniedNotOwner,
			wantDropReason: monitoring.TriggerDropReasonShardDeniedNotOwner,
			wantNotOwner:   1,
		},
		{
			name:           "shard denial orchestrator error",
			admissionErr:   v2.ErrShardDeniedOrchestrator,
			wantDropReason: monitoring.TriggerDropReasonShardDeniedOrchestrator,
			wantOrchErr:    1,
		},
		{
			name:           "other denial",
			admissionErr:   fmt.Errorf("not owned by this node: %w", triggers.ErrAdmissionDenied),
			wantDropReason: monitoring.TriggerDropReasonAdmissionDenied,
		},
	}
	for _, mode := range []lifecycleEngineMode{lifecycleLegacy, lifecycleCoordinated} {
		t.Run(string(mode), func(t *testing.T) {
			t.Parallel()
			for _, tc := range tests {
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()
					meter := newDropMeter(t)
					h := newEngineLifecycleHarness(t, mode, withMeter(meter.meter),
						withAdmission(func(context.Context, triggers.CoordinatedEvent) error { return tc.admissionErr }))
					h.start()

					h.deliver(denialWait)

					require.EventuallyWithT(t, func(c *assert.CollectT) {
						assert.Equal(c, map[string]int64{tc.wantDropReason: 1}, meter.drops(t), "drops by reason")
						assert.Equal(c, int32(1), h.ackCalls.Load(), "a denied event is ACKed")
					}, denialWait, 5*time.Millisecond)
					// Both paths count the shard denial before the drop, so it is final by now.
					assert.Equal(t, tc.wantNotOwner, meter.total(t, shardNotOwnerTotal), shardNotOwnerTotal)
					assert.Equal(t, tc.wantOrchErr, meter.total(t, shardOrchErrTotal), shardOrchErrTotal)
					assert.Zero(t, h.moduleExecutions.Load(), "a denied event never reaches the module")
					assert.Empty(t, h.rec.Messages(t, "beholder_entity", v2ExecutionStartedEntity), "execution started events")
					assert.Empty(t, h.rec.Messages(t, "beholder_entity", v2ExecutionFinishedEntity), "execution finished events")
				})
			}
		})
	}
}

// TestEngine_AdmissionCutShortParity asserts that the legacy and the coordinated engine handle an
// admission check that is still running when ingress stops the same way, e.g. a shard ownership
// lookup against an unresponsive orchestrator. Stopping ingress (closing the legacy engine,
// unregistering the coordinated one) must cancel the check and return promptly. The check failed
// only because it was cut short, so the event is not ACKed, not counted as a denial, and dropped
// as unregistered. It never takes an execution slot.
func TestEngine_AdmissionCutShortParity(t *testing.T) {
	t.Parallel()

	const (
		shardOrchErrTotal = "platform_engine_shard_execution_denied_orchestrator_error"
		// Stopping ingress cancels the check at once, so a failure shows within seconds.
		cutShortWait = 5 * time.Second
	)

	for _, mode := range []lifecycleEngineMode{lifecycleLegacy, lifecycleCoordinated} {
		t.Run(string(mode), func(t *testing.T) {
			t.Parallel()
			meter := newDropMeter(t)
			admitting := make(chan struct{})
			checkErr := make(chan error, 1)
			h := newEngineLifecycleHarness(t, mode, withMeter(meter.meter),
				withAdmission(func(ctx context.Context, _ triggers.CoordinatedEvent) error {
					close(admitting)
					<-ctx.Done()
					checkErr <- ctx.Err()
					// A resolver reports the cut-short lookup as a failure, which would read as a denial.
					return v2.ErrShardDeniedOrchestrator
				}))
			h.start()

			h.deliver(cutShortWait)
			select {
			case <-admitting:
			case <-time.After(cutShortWait):
				t.Fatal("admission check never started")
			}

			stopped := make(chan error, 1)
			go func() {
				if mode == lifecycleCoordinated {
					stopped <- h.coordinator.UnregisterTriggers(h.tenantCtx(), h.cfg.WorkflowID)
					return
				}
				stopped <- h.close()
			}()
			select {
			case err := <-stopped:
				require.NoError(t, err)
			case <-time.After(cutShortWait):
				t.Fatal("stopping ingress did not return while an admission check was running")
			}
			select {
			case err := <-checkErr:
				require.ErrorIs(t, err, context.Canceled)
			case <-time.After(cutShortWait):
				t.Fatal("admission check was not cancelled")
			}

			require.EventuallyWithT(t, func(c *assert.CollectT) {
				assert.Equal(c, map[string]int64{monitoring.TriggerDropReasonUnregistered: 1}, meter.drops(t), "drops by reason")
			}, cutShortWait, 5*time.Millisecond)
			assert.Zero(t, h.ackCalls.Load(), "a cut-short check decided nothing, so the event is not ACKed")
			assert.Zero(t, meter.total(t, shardOrchErrTotal), shardOrchErrTotal)
			assert.Zero(t, h.slots.waits.Load(), "an admission check never waits for an execution slot")
			assert.Zero(t, h.moduleExecutions.Load(), "the event never reaches the module")
		})
	}
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

// parityLog is a log line a delivery log test expects. A nil value in fields only checks that the
// key is present.
type parityLog struct {
	level  zapcore.Level
	msg    string
	fields map[string]any
}

// TestEngine_DeliveryLogParity asserts that the legacy and the coordinated engine log each
// admission and delivery outcome with the same message, level and fields, so log queries work for
// both. Every line in both is asserted against both engines. legacy and coordinated hold the lines
// only one of them logs, where the two deliberately differ. No other line is logged at Error.
func TestEngine_DeliveryLogParity(t *testing.T) {
	t.Parallel()

	const eventID = "evt-1"
	ids := map[string]any{"triggerID": lifecycleTriggerCapID, "eventID": eventID}
	idsAndErr := map[string]any{"triggerID": lifecycleTriggerCapID, "eventID": eventID, "err": nil}
	enqueued := parityLog{zapcore.DebugLevel, "Enqueued trigger event", ids}
	popped := parityLog{zapcore.DebugLevel, "Popped a trigger event from the queue", map[string]any{"eventID": eventID, "eventAgeMs": nil}}
	// Legacy logs every admission error again at Error under this message.
	putFailed := parityLog{zapcore.ErrorLevel, "Failed to put routed trigger event", idsAndErr}
	denied := parityLog{zapcore.InfoLevel, "Trigger event denied admission, dropping", idsAndErr}

	tests := []struct {
		name      string
		admission func(context.Context, triggers.CoordinatedEvent) error
		ackErr    error
		// cutShort blocks the admission check until ingress stops.
		cutShort bool
		// drain drains the engine before the event is delivered.
		drain bool
		// duplicate delivers the event a second time once its first execution is running.
		duplicate   bool
		both        []parityLog
		legacy      []parityLog
		coordinated []parityLog
	}{
		{
			name: "allowed",
			both: []parityLog{enqueued, popped},
		},
		{
			name:      "cached",
			admission: func(context.Context, triggers.CoordinatedEvent) error { return v2.ErrAdmissionCache },
			both:      []parityLog{{zapcore.InfoLevel, "Trigger event cached for failover, not enqueuing", ids}},
			legacy:    []parityLog{putFailed},
		},
		{
			name:        "denied not owner",
			admission:   func(context.Context, triggers.CoordinatedEvent) error { return v2.ErrShardDeniedNotOwner },
			legacy:      []parityLog{putFailed},
			coordinated: []parityLog{denied},
		},
		{
			name:        "denied orchestrator",
			admission:   func(context.Context, triggers.CoordinatedEvent) error { return v2.ErrShardDeniedOrchestrator },
			legacy:      []parityLog{putFailed},
			coordinated: []parityLog{denied},
		},
		{
			name:      "ACK fails after denial",
			admission: func(context.Context, triggers.CoordinatedEvent) error { return v2.ErrShardDeniedNotOwner },
			ackErr:    errors.New("ack boom"),
			legacy: []parityLog{
				putFailed,
				{zapcore.ErrorLevel, "failed to ACK trigger after admission denial", map[string]any{"eventID": eventID, "err": nil}},
			},
			coordinated: []parityLog{
				denied,
				{zapcore.ErrorLevel, "Failed to ACK trigger event after admission denial", idsAndErr},
			},
		},
		{
			name:     "ctx cancelled",
			cutShort: true,
			both:     []parityLog{{zapcore.InfoLevel, "Trigger admission cut short, dropping event without ACK", idsAndErr}},
			legacy:   []parityLog{putFailed},
		},
		{
			name:        "unknown error",
			admission:   func(context.Context, triggers.CoordinatedEvent) error { return errors.New("boom") },
			legacy:      []parityLog{putFailed},
			coordinated: []parityLog{{zapcore.ErrorLevel, "Trigger admission check failed, dropping event", idsAndErr}},
		},
		{
			name:  "draining",
			drain: true,
			legacy: []parityLog{
				{zapcore.InfoLevel, "Engine is draining, dropping trigger event before enqueue", ids},
				{zapcore.InfoLevel, "Dropping trigger event: engine draining", ids},
			},
			// The coordinator queues the event, and the engine rejects it when it executes.
			coordinated: []parityLog{enqueued, popped, {zapcore.InfoLevel, "Engine is draining, dropping trigger event", ids}},
		},
		{
			name:      "duplicate execution",
			duplicate: true,
			both:      []parityLog{{zapcore.InfoLevel, "Skipping trigger event execution", idsAndErr}},
		},
	}
	for _, mode := range []lifecycleEngineMode{lifecycleLegacy, lifecycleCoordinated} {
		t.Run(string(mode), func(t *testing.T) {
			t.Parallel()
			for _, tc := range tests {
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()
					lggr, logs := commonlogger.TestObserved(t, zapcore.DebugLevel)
					admitting := make(chan struct{})
					admission := tc.admission
					if tc.cutShort {
						admission = func(ctx context.Context, _ triggers.CoordinatedEvent) error {
							close(admitting)
							<-ctx.Done()
							return v2.ErrShardDeniedOrchestrator
						}
					}
					opts := []harnessOption{withLogger(lggr), withAckErr(tc.ackErr)}
					if admission != nil {
						opts = append(opts, withAdmission(admission))
					}
					h := newEngineLifecycleHarness(t, mode, opts...)
					h.start()
					if tc.drain {
						require.True(t, h.engine.Drain())
					}

					h.deliverEvent(eventID, testutils.WaitTimeout(t))
					switch {
					case tc.cutShort:
						awaitLifecycle(t, admitting, "admission check to start")
						if mode == lifecycleCoordinated {
							require.NoError(t, h.coordinator.UnregisterTriggers(h.tenantCtx(), h.cfg.WorkflowID))
						} else {
							require.NoError(t, h.close())
						}
					case tc.duplicate:
						awaitLifecycle(t, h.blocker.started, "first execution to start")
						h.deliverEvent(eventID, testutils.WaitTimeout(t))
					}

					want := append(append([]parityLog{}, tc.both...), tc.legacy...)
					if mode == lifecycleCoordinated {
						want = append(append([]parityLog{}, tc.both...), tc.coordinated...)
					}
					requireLogs(t, logs, want)
				})
			}
		})
	}
}

// requireLogs waits until every line in want is logged, then asserts that each was logged once,
// with its level and fields, and that no other line was logged at Error. Every outcome logs within
// moments of delivery, so a missing line fails within seconds instead of waiting out the test
// deadline.
func requireLogs(t *testing.T, logs *observer.ObservedLogs, want []parityLog) {
	t.Helper()
	const logWait = 5 * time.Second
	require.Eventually(t, func() bool {
		for _, w := range want {
			if logs.FilterMessage(w.msg).Len() == 0 {
				return false
			}
		}
		return true
	}, logWait, 5*time.Millisecond, "every wanted line is logged")

	var wantErrors []string
	for _, w := range want {
		entries := logs.FilterMessage(w.msg).All()
		require.Len(t, entries, 1, "%q is logged once", w.msg)
		assert.Equal(t, w.level, entries[0].Level, "level of %q", w.msg)
		got := entries[0].ContextMap()
		for k, v := range w.fields {
			if assert.Contains(t, got, k, "field %q of %q", k, w.msg) && v != nil {
				assert.Equal(t, v, got[k], "field %q of %q", k, w.msg)
			}
		}
		if w.level == zapcore.ErrorLevel {
			wantErrors = append(wantErrors, w.msg)
		}
	}
	errorEntries := logs.FilterLevelExact(zapcore.ErrorLevel).All()
	gotErrors := make([]string, 0, len(errorEntries))
	for _, entry := range errorEntries {
		gotErrors = append(gotErrors, entry.Message)
	}
	assert.ElementsMatch(t, wantErrors, gotErrors, "Error lines")
}
