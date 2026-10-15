package triggers

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

const (
	testTriggerCapID = "cron-trigger@1.0.0"
	testMethod       = "Trigger"
	testOwner        = "abcdef" // already normalized: contexts.WithCRE strips a 0x prefix
	testOrg          = "org-1"
	testDonID        = uint32(7)
)

var testTenant = contexts.CRE{Org: testOrg, Owner: testOwner, Workflow: validWorkflowID}

type fakeSubscriber struct {
	tenant contexts.CRE
	subs   []*sdkpb.TriggerSubscription
	err    error
	calls  int
}

func (s *fakeSubscriber) Subscribe(context.Context) ([]*sdkpb.TriggerSubscription, error) {
	s.calls++
	return s.subs, s.err
}

func (s *fakeSubscriber) Tenant() contexts.CRE { return s.tenant }

type fakeEngine struct {
	coordinated bool
	// execute, if set, runs inside ExecuteTrigger on the execution's goroutine.
	execute func(ctx context.Context, event CoordinatedEvent)

	mu      sync.Mutex
	events  []CoordinatedEvent
	tenants []contexts.CRE
}

func (e *fakeEngine) ExecuteTrigger(ctx context.Context, event CoordinatedEvent) error {
	e.mu.Lock()
	e.events = append(e.events, event)
	e.tenants = append(e.tenants, contexts.CREValue(ctx))
	e.mu.Unlock()
	if e.execute != nil {
		e.execute(ctx, event)
	}
	return nil
}

func (e *fakeEngine) IsCoordinated() bool { return e.coordinated }

// executedIDs returns the event IDs of every ExecuteTrigger call, in call order.
func (e *fakeEngine) executedIDs() []string {
	executed := e.executed()
	ids := make([]string, 0, len(executed))
	for _, ev := range executed {
		ids = append(ids, ev.Event.Event.ID)
	}
	return ids
}

func (e *fakeEngine) executed() []CoordinatedEvent {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]CoordinatedEvent(nil), e.events...)
}

type fakeEngineRegistry struct {
	mu      sync.Mutex
	engines map[types.WorkflowID]RegisteredEngine
}

func newFakeEngineRegistry() *fakeEngineRegistry {
	return &fakeEngineRegistry{engines: make(map[types.WorkflowID]RegisteredEngine)}
}

func (r *fakeEngineRegistry) Get(wid types.WorkflowID) (RegisteredEngine, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.engines[wid]
	return e, ok
}

func (r *fakeEngineRegistry) set(wid types.WorkflowID, e RegisteredEngine) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.engines[wid] = e
}

// fakeWorkflowLimits counts slots, recording the tenant each call was scoped to.
type fakeWorkflowLimits struct {
	mu      sync.Mutex
	used    int
	frees   int
	useErr  error
	tenants []contexts.CRE
}

func (l *fakeWorkflowLimits) Close() error                           { return nil }
func (l *fakeWorkflowLimits) Limit(context.Context) (int, error)     { return 100, nil }
func (l *fakeWorkflowLimits) Available(context.Context) (int, error) { return 100, nil }
func (l *fakeWorkflowLimits) Use(ctx context.Context, amount int) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.tenants = append(l.tenants, contexts.CREValue(ctx))
	if l.useErr != nil {
		return l.useErr
	}
	l.used += amount
	return nil
}
func (l *fakeWorkflowLimits) Free(ctx context.Context, amount int) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.tenants = append(l.tenants, contexts.CREValue(ctx))
	l.used -= amount
	l.frees++
	return nil
}
func (l *fakeWorkflowLimits) inUse() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.used
}
func (l *fakeWorkflowLimits) freeCount() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.frees
}

type coordinatorFixture struct {
	c       *coordinator
	capReg  *regmocks.CapabilitiesRegistry
	engines *fakeEngineRegistry
	limits  *fakeWorkflowLimits
	clock   *clockwork.FakeClock
	wid     types.WorkflowID
	engine  *fakeEngine
}

// newCoordinatorFixture starts a coordinator with one coordinated engine
// already in the registry for testTenant's workflow.
func newCoordinatorFixture(t *testing.T) *coordinatorFixture {
	t.Helper()
	return newCoordinatorFixtureWithAdmission(t, newTestAdmission())
}

// newCoordinatorFixtureWithAdmission is newCoordinatorFixture with the given admission limiters.
func newCoordinatorFixtureWithAdmission(t *testing.T, admission AdmissionLimiters) *coordinatorFixture {
	t.Helper()
	capReg := regmocks.NewCapabilitiesRegistry(t)
	engines := newFakeEngineRegistry()
	wfLimits := &fakeWorkflowLimits{}
	clock := clockwork.NewFakeClock()

	wid, err := types.WorkflowIDFromHex(validWorkflowID)
	require.NoError(t, err)
	engine := &fakeEngine{coordinated: true}
	engines.set(wid, engine)

	c := NewCoordinator(newRegisterDeps(t, capReg, limits.NewGateLimiter(true)), admission, engines, wfLimits, clock)
	servicetest.Run(t, c)

	return &coordinatorFixture{
		c: c.(*coordinator), capReg: capReg, engines: engines, limits: wfLimits,
		clock: clock, wid: wid, engine: engine,
	}
}

// expectTrigger wires a trigger capability that hands back a fresh event channel on registration.
func (f *coordinatorFixture) expectTrigger(t *testing.T) (*capmocks.TriggerCapability, chan capabilities.TriggerResponse) {
	t.Helper()
	trigger := capmocks.NewTriggerCapability(t)
	eventCh := make(chan capabilities.TriggerResponse)
	trigger.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).
		Return((<-chan capabilities.TriggerResponse)(eventCh), nil).Once()
	f.capReg.EXPECT().GetTrigger(mock.Anything, testTriggerCapID).Return(trigger, nil).Once()
	return trigger, eventCh
}

// newTestAdmission returns unscoped admission limiters with the production defaults.
func newTestAdmission() AdmissionLimiters {
	return AdmissionLimiters{
		NewTriggerEventQueue: newTestQueue(50),
		ExecutionConcurrency: limits.GlobalResourcePoolLimiter(5),
	}
}

// newTestQueue returns an AdmissionLimiters.NewTriggerEventQueue that builds an
// unscoped queue of the given capacity for each registration. It goes through
// the factory, as production does, because the static limits.NewQueueLimiter
// panics on Close, and each registration closes its queue.
func newTestQueue(capacity int) func() (limits.QueueLimiter[CoordinatedEvent], error) {
	return func() (limits.QueueLimiter[CoordinatedEvent], error) {
		return limits.MakeQueueLimiter[CoordinatedEvent](limits.Factory{}, settings.Int(capacity))
	}
}

// closeTrackingQueue records whether its registration closed it.
type closeTrackingQueue struct {
	limits.QueueLimiter[CoordinatedEvent]
	closed atomic.Bool
}

func (q *closeTrackingQueue) Close() error {
	q.closed.Store(true)
	return q.QueueLimiter.Close()
}

// queueTracker hands each registration a fresh closeTrackingQueue and keeps
// them in creation order.
type queueTracker struct {
	mu     sync.Mutex
	queues []*closeTrackingQueue
}

func (q *queueTracker) newQueue() (limits.QueueLimiter[CoordinatedEvent], error) {
	inner, err := newTestQueue(50)()
	if err != nil {
		return nil, err
	}
	queue := &closeTrackingQueue{QueueLimiter: inner}
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queues = append(q.queues, queue)
	return queue, nil
}

// get returns the queue of the i-th registration, failing the test if there is none.
func (q *queueTracker) get(t *testing.T, i int) *closeTrackingQueue {
	t.Helper()
	q.mu.Lock()
	defer q.mu.Unlock()
	require.Greater(t, len(q.queues), i, "queues created")
	return q.queues[i]
}

// newTrackedAdmission returns test admission limiters whose queues record Close.
func newTrackedAdmission() (AdmissionLimiters, *queueTracker) {
	queues := &queueTracker{}
	admission := newTestAdmission()
	admission.NewTriggerEventQueue = queues.newQueue
	return admission, queues
}

func newTestSubscriber() *fakeSubscriber {
	return &fakeSubscriber{
		tenant: testTenant,
		subs:   []*sdkpb.TriggerSubscription{{Id: testTriggerCapID, Method: testMethod}},
	}
}

func testParams(t *testing.T) RegistrationParams {
	t.Helper()
	name, err := types.NewWorkflowName("my-workflow")
	require.NoError(t, err)
	return RegistrationParams{WorkflowOwner: testOwner, WorkflowName: name, WorkflowDonID: testDonID}
}

func (f *coordinatorFixture) registered() bool {
	_, ok := f.c.workflows.get(validWorkflowID)
	return ok
}

func triggerEvent(id string) capabilities.TriggerResponse {
	return capabilities.TriggerResponse{Event: capabilities.TriggerEvent{TriggerType: testTriggerCapID, ID: id}}
}

// blockingExecution makes the fixture engine park inside ExecuteTrigger until
// release is closed, signalling started once it is in flight.
func (f *coordinatorFixture) blockingExecution() (started chan struct{}, release chan struct{}) {
	started, release = make(chan struct{}, 1), make(chan struct{})
	f.engine.execute = func(context.Context, CoordinatedEvent) {
		started <- struct{}{}
		<-release
	}
	return started, release
}

// holdExecutions makes the fixture engine hold every execution until release is
// called. Unlike blockingExecution it never blocks on signalling a start, so it
// holds any number of executions. Release also runs at cleanup.
func (f *coordinatorFixture) holdExecutions(t *testing.T) (release func()) {
	ch := make(chan struct{})
	var once sync.Once
	release = func() { once.Do(func() { close(ch) }) }
	t.Cleanup(release)
	f.engine.execute = func(context.Context, CoordinatedEvent) { <-ch }
	return release
}

// countingPool counts Wait calls on a ResourcePoolLimiter, so a test can tell
// when the admitter has taken an event off the queue and waits for a slot.
type countingPool struct {
	limits.ResourcePoolLimiter[int]
	waits atomic.Int32
}

func (p *countingPool) Wait(ctx context.Context, amount int) (func(), error) {
	p.waits.Add(1)
	return p.ResourcePoolLimiter.Wait(ctx, amount)
}

// useTestMeter points the global beholder meter at a manual reader, so workflow
// metrics created after this call can be read back. It swaps a global, so a test
// that uses it must not run in parallel.
func useTestMeter(t *testing.T) *sdkmetric.ManualReader {
	t.Helper()
	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { _ = mp.Shutdown(context.Background()) })

	prev := beholder.GetClient()
	t.Cleanup(func() { beholder.SetClient(prev) })
	client := beholder.NoopClientConfig{Lggr: logger.Test(t)}.New()
	client.Meter = mp.Meter("beholder")
	client.MeterProvider = mp
	beholder.SetClient(client)
	return reader
}

// counterTotal sums the named counter over the data points that carry every
// attribute in attrs.
func counterTotal(t *testing.T, reader *sdkmetric.ManualReader, name string, attrs map[string]string) int64 {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	var total int64
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name != name {
				continue
			}
			sum, ok := m.Data.(metricdata.Sum[int64])
			require.True(t, ok, "expected Sum[int64] for %s, got %T", name, m.Data)
			for _, dp := range sum.DataPoints {
				matches := true
				for k, want := range attrs {
					if got, ok := dp.Attributes.Value(attribute.Key(k)); !ok || got.AsString() != want {
						matches = false
					}
				}
				if matches {
					total += dp.Value
				}
			}
		}
	}
	return total
}
