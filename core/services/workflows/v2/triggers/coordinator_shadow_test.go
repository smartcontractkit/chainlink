package triggers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

// shadowStart is where the shadow tests' fake clock starts. It is a whole second so window edges are exact.
var shadowStart = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// newShadowMode returns the shadow mode limiter for the real setting. settingsJSON uses the cresettings
// getter format. An empty string leaves the setting at its default.
func newShadowMode(t *testing.T, settingsJSON string) limits.RangeLimiter[commonconfig.Timestamp] {
	t.Helper()
	if settingsJSON == "" {
		return newShadowModeFrom(t, nil)
	}
	return newShadowModeFrom(t, newJSONGetter(t, settingsJSON))
}

// newShadowModeFrom returns the shadow mode limiter for the real setting, read from getter. A nil getter
// leaves the setting at its default.
func newShadowModeFrom(t *testing.T, getter settings.Getter) limits.RangeLimiter[commonconfig.Timestamp] {
	t.Helper()
	l, err := limits.MakeRangeLimiter[commonconfig.Timestamp](limits.Factory{Logger: logger.Test(t), Settings: getter},
		cresettings.Default.PerWorkflow.FeatureShadowTriggerModeActivePeriod)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, l.Close()) })
	return l
}

func newJSONGetter(t *testing.T, settingsJSON string) settings.Getter {
	t.Helper()
	getter, err := settings.NewJSONGetter([]byte(settingsJSON))
	require.NoError(t, err)
	return getter
}

// swapGetter is a settings.Getter whose settings can be replaced while the node runs.
type swapGetter struct {
	current atomic.Pointer[settings.Getter]
}

func (g *swapGetter) set(t *testing.T, settingsJSON string) {
	t.Helper()
	getter := newJSONGetter(t, settingsJSON)
	g.current.Store(&getter)
}

func (g *swapGetter) GetScoped(ctx context.Context, scope settings.Scope, key string) (string, error) {
	return (*g.current.Load()).GetScoped(ctx, scope, key)
}

func window(from, to time.Time) string { return fmt.Sprintf("[%d,%d]", from.Unix(), to.Unix()) }

// globalWindow returns settings that set the shadow window for every workflow.
func globalWindow(from, to time.Time) string {
	return fmt.Sprintf(`{"global":{"PerWorkflow":{"FeatureShadowTriggerModeActivePeriod":%q}}}`, window(from, to))
}

// orgWindow returns settings that set the shadow window for one org only.
func orgWindow(org string, from, to time.Time) string {
	return fmt.Sprintf(`{"org":{%q:{"PerWorkflow":{"FeatureShadowTriggerModeActivePeriod":%q}}}}`, org, window(from, to))
}

// sendEvent sends the event and waits for it to execute. Shadow routing happens before execution, so the
// central queue has seen the event once this returns.
func sendEvent(t *testing.T, f *coordinatorFixture, eventCh chan<- capabilities.TriggerResponse, id string) {
	t.Helper()
	executed := len(f.engine.executed())
	eventCh <- triggerEvent(id)
	require.Eventually(t, func() bool { return len(f.engine.executed()) == executed+1 }, 5*time.Second, 5*time.Millisecond)
}

// registerWorkflow registers sub with the fixture's coordinator and returns its trigger channel.
func registerWorkflow(t *testing.T, f *coordinatorFixture, sub *fakeSubscriber) chan capabilities.TriggerResponse {
	t.Helper()
	trigger, eventCh := f.expectTrigger(t)
	trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Maybe()
	trigger.EXPECT().AckEvent(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	_, err := f.c.RegisterTriggers(t.Context(), sub, testParams(t))
	require.NoError(t, err)
	return eventCh
}

// TestCoordinator_ShadowDisabledByDefault asserts that with the default window no event reaches the central
// queue and every event still executes.
func TestCoordinator_ShadowDisabledByDefault(t *testing.T) {
	t.Parallel()
	q := newTestCentralQueue(t, &fakeLimit{limit: 10})
	f := newCoordinatorFixtureAt(t, clockwork.NewFakeClockAt(shadowStart), newTestAdmission(),
		WithShadowMode(q, newShadowMode(t, "")))
	eventCh := registerWorkflow(t, f, newTestSubscriber())

	eventCh <- triggerEvent("evt-1")
	eventCh <- triggerEvent("evt-2")
	require.Eventually(t, func() bool { return len(f.engine.executed()) == 2 }, 5*time.Second, 5*time.Millisecond)

	assert.Empty(t, drainIDs(t, q))
}

// TestCoordinator_ShadowWindow asserts only events observed inside the window reach the central queue.
// The window has one-second resolution, so an event in the second before the start is out and one in the
// second of the end is in. Moving the clock is enough, so nothing is cached. Every event executes, and
// shadowed events already carry their deadline.
func TestCoordinator_ShadowWindow(t *testing.T) {
	t.Parallel()
	q := newTestCentralQueue(t, &fakeLimit{limit: 10})
	shadow := newShadowMode(t, globalWindow(shadowStart.Add(time.Minute), shadowStart.Add(3*time.Minute)))
	f := newCoordinatorFixtureAt(t, clockwork.NewFakeClockAt(shadowStart), newTestAdmission(), WithShadowMode(q, shadow))
	eventCh := registerWorkflow(t, f, newTestSubscriber())

	// Each event is handled before the clock moves, so its observed time is exact.
	send := func(id string, advance time.Duration) {
		t.Helper()
		f.clock.Advance(advance)
		sendEvent(t, f, eventCh, id)
	}
	send("before", 0)
	send("just-before-start", time.Minute-time.Millisecond)
	send("at-start", time.Millisecond)
	send("inside", time.Minute)
	send("at-end", time.Minute)
	send("just-past-end", 999*time.Millisecond)
	send("after", time.Second)

	shadowed, err := q.Drain(t.Context())
	require.NoError(t, err)
	assert.Equal(t, []string{"at-start", "inside", "at-end", "just-past-end"}, eventIDs(shadowed))
	for _, e := range shadowed {
		assert.Equal(t, e.ObservedAt.Add(testQueueTimeout), e.Deadline, "deadline of %s", e.Event.Event.ID)
	}
	assert.Equal(t, []string{"before", "just-before-start", "at-start", "inside", "at-end", "just-past-end", "after"},
		f.engine.executedIDs())
}

// TestCoordinator_ShadowSettingChange asserts a change to the shadow setting applies to the next event
// without a restart. The clock never moves, so only the setting decides.
func TestCoordinator_ShadowSettingChange(t *testing.T) {
	t.Parallel()
	q := newTestCentralQueue(t, &fakeLimit{limit: 10})
	covering := globalWindow(shadowStart.Add(-time.Hour), shadowStart.Add(time.Hour))
	future := globalWindow(shadowStart.Add(time.Hour), shadowStart.Add(2*time.Hour))
	getter := &swapGetter{}
	getter.set(t, covering)
	f := newCoordinatorFixtureAt(t, clockwork.NewFakeClockAt(shadowStart), newTestAdmission(),
		WithShadowMode(q, newShadowModeFrom(t, getter)))
	eventCh := registerWorkflow(t, f, newTestSubscriber())

	sendEvent(t, f, eventCh, "on-1")
	getter.set(t, future)
	sendEvent(t, f, eventCh, "off")
	getter.set(t, covering)
	sendEvent(t, f, eventCh, "on-2")

	assert.Equal(t, []string{"on-1", "on-2"}, drainIDs(t, q))
	assert.Equal(t, []string{"on-1", "off", "on-2"}, f.engine.executedIDs())
}

// TestCoordinator_ShadowScopedByTenant asserts a window set for one org shadows that org's workflows only.
// The shadowed counter is recorded per workflow.
func TestCoordinator_ShadowScopedByTenant(t *testing.T) {
	t.Parallel()
	q := newTestCentralQueue(t, &fakeLimit{limit: 10})
	shadow := newShadowMode(t, orgWindow(testOrg, shadowStart.Add(-time.Hour), shadowStart.Add(time.Hour)))
	wfMetrics, reader := newTestMeterMetrics(t)
	f := newCoordinatorFixtureWithMetrics(t, clockwork.NewFakeClockAt(shadowStart), newTestAdmission(), wfMetrics,
		WithShadowMode(q, shadow))

	otherWorkflowID := strings.Repeat("cd", 32)
	otherWID, err := types.WorkflowIDFromHex(otherWorkflowID)
	require.NoError(t, err)
	otherEngine := &fakeEngine{coordinated: true}
	f.engines.set(otherWID, otherEngine)
	otherSub := newTestSubscriber()
	otherSub.tenant = contexts.CRE{Org: "org-2", Owner: testOwner, Workflow: otherWorkflowID}

	eventCh := registerWorkflow(t, f, newTestSubscriber())
	otherEventCh := registerWorkflow(t, f, otherSub)
	eventCh <- triggerEvent("org-1-evt")
	otherEventCh <- triggerEvent("org-2-evt")
	require.Eventually(t, func() bool {
		return len(f.engine.executed()) == 1 && len(otherEngine.executed()) == 1
	}, 5*time.Second, 5*time.Millisecond)

	assert.Equal(t, []string{"org-1-evt"}, drainIDs(t, q))
	assert.Equal(t, int64(1), counterTotal(t, reader, "platform_engine_trigger_event_shadowed_total",
		map[string]string{platform.KeyWorkflowID: validWorkflowID}))
	assert.Equal(t, int64(0), counterTotal(t, reader, "platform_engine_trigger_event_shadowed_total",
		map[string]string{platform.KeyWorkflowID: otherWorkflowID}))
}

// orderLog records puts to the central and workflow queues in the order they happen.
type orderLog struct {
	mu      sync.Mutex
	entries []string
}

func (l *orderLog) add(entry string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, entry)
}

func (l *orderLog) get() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.entries...)
}

// recordingQueue logs each Put to the central queue before passing it on.
type recordingQueue struct {
	Queue
	log *orderLog
}

func (q *recordingQueue) Put(ctx context.Context, e CoordinatedEvent) error {
	q.log.add("central " + e.Event.Event.ID)
	return q.Queue.Put(ctx, e)
}

// recordingWorkflowQueue logs each Put to a workflow's own queue before passing it on.
type recordingWorkflowQueue struct {
	limits.QueueLimiter[CoordinatedEvent]
	log *orderLog
}

func (q *recordingWorkflowQueue) Put(ctx context.Context, e CoordinatedEvent) error {
	q.log.add("workflow " + e.Event.Event.ID)
	return q.QueueLimiter.Put(ctx, e)
}

// failingTimeout is a TriggerEventQueueTimeout limiter that never resolves a value.
type failingTimeout struct{}

func (failingTimeout) Limit(context.Context) (time.Duration, error) {
	return 0, limits.ErrMissingTenant{Scope: settings.ScopeWorkflow}
}
func (failingTimeout) Check(context.Context, time.Duration) error { return nil }
func (failingTimeout) Close() error                               { return nil }

// TestCoordinator_ShadowPreAdmission asserts the central put happens before the workflow queue put, so the
// central queue sees events the engine path later drops. An event without a deadline reaches neither queue.
func TestCoordinator_ShadowPreAdmission(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		queueCapacity int
		setup         func(a *AdmissionLimiters)
		wantLog       []string
		wantCentral   []string
		wantExecuted  []string
	}{
		{
			name:          "accepted",
			queueCapacity: 10,
			wantLog:       []string{"central evt-1", "workflow evt-1"},
			wantCentral:   []string{"evt-1"},
			wantExecuted:  []string{"evt-1"},
		},
		{
			name:          "workflow queue full",
			queueCapacity: 0,
			wantLog:       []string{"central evt-1", "workflow evt-1"},
			wantCentral:   []string{"evt-1"},
		},
		{
			name:          "timeout read fails",
			queueCapacity: 10,
			setup: func(a *AdmissionLimiters) {
				a.TriggerEventQueueTimeout = failingTimeout{}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			log := &orderLog{}
			central := newTestCentralQueue(t, &fakeLimit{limit: 10})
			admission := newTestAdmission()
			admission.NewTriggerEventQueue = func() (limits.QueueLimiter[CoordinatedEvent], error) {
				inner, err := newTestQueue(tt.queueCapacity)()
				return &recordingWorkflowQueue{QueueLimiter: inner, log: log}, err
			}
			if tt.setup != nil {
				tt.setup(&admission)
			}
			shadow := newShadowMode(t, globalWindow(shadowStart.Add(-time.Hour), shadowStart.Add(time.Hour)))
			f := newCoordinatorFixtureAt(t, clockwork.NewFakeClockAt(shadowStart), admission,
				WithShadowMode(&recordingQueue{Queue: central, log: log}, shadow))
			engine := f.engine
			eventCh := registerWorkflow(t, f, newTestSubscriber())

			eventCh <- triggerEvent("evt-1")
			// The reader handles one event at a time, so once it takes the probe evt-1 has been handled.
			eventCh <- capabilities.TriggerResponse{Err: errors.New("probe")}
			require.Eventually(t, func() bool { return len(engine.executed()) == len(tt.wantExecuted) }, 5*time.Second, 5*time.Millisecond)

			assert.Equal(t, tt.wantLog, log.get())
			assert.Equal(t, tt.wantCentral, drainIDs(t, central))
			assert.ElementsMatch(t, tt.wantExecuted, engine.executedIDs())
		})
	}
}

// failingQueue is a central queue whose Put always fails with an error other than ErrQueueFull.
type failingQueue struct {
	puts atomic.Int32
}

func (q *failingQueue) Put(context.Context, CoordinatedEvent) error {
	q.puts.Add(1)
	return errors.New("boom")
}

func (q *failingQueue) Drain(context.Context) ([]CoordinatedEvent, error) { return nil, nil }
func (q *failingQueue) Stats(context.Context) (Stats, error)              { return Stats{}, nil }

// TestCoordinator_ShadowBestEffort asserts a failed central put, for any reason, never stops execution and
// is not counted by the engine's dropped events counter. A full queue counts the drop itself.
func TestCoordinator_ShadowBestEffort(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		// setup returns the central queue and a check to run once both events have executed.
		setup func(t *testing.T) (Queue, func(t *testing.T))
	}{
		{
			name: "queue full",
			setup: func(t *testing.T) (Queue, func(t *testing.T)) {
				qm, queueReader := testQueueMetrics(t)
				return NewQueue(&fakeLimit{limit: 0}, qm, logger.Test(t)), func(t *testing.T) {
					assert.Equal(t, int64(2), int64Sum(t, queueReader, "platform_engine_central_trigger_queue_dropped_events_total"))
				}
			},
		},
		{
			name: "put fails",
			setup: func(*testing.T) (Queue, func(t *testing.T)) {
				q := &failingQueue{}
				return q, func(t *testing.T) {
					assert.Equal(t, int32(2), q.puts.Load(), "every event must be put to the central queue")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			wfMetrics, reader := newTestMeterMetrics(t)
			central, check := tt.setup(t)
			shadow := newShadowMode(t, globalWindow(shadowStart.Add(-time.Hour), shadowStart.Add(time.Hour)))
			f := newCoordinatorFixtureWithMetrics(t, clockwork.NewFakeClockAt(shadowStart), newTestAdmission(), wfMetrics,
				WithShadowMode(central, shadow))
			eventCh := registerWorkflow(t, f, newTestSubscriber())

			eventCh <- triggerEvent("evt-1")
			eventCh <- triggerEvent("evt-2")
			require.Eventually(t, func() bool { return len(f.engine.executed()) == 2 }, 5*time.Second, 5*time.Millisecond)

			assert.Equal(t, int64(2), counterTotal(t, reader, "platform_engine_trigger_event_enqueued_total", nil),
				"the reader must see the coordinator's metrics")
			assert.Equal(t, int64(0), counterTotal(t, reader, "platform_engine_trigger_event_dropped_total", nil))
			assert.Equal(t, int64(0), counterTotal(t, reader, "platform_engine_trigger_event_shadowed_total", nil),
				"a failed central put is not counted as shadowed")
			check(t)
		})
	}
}
