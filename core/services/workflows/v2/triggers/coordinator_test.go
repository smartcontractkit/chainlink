package triggers

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
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
	// execute, if set, runs inside ExecuteTrigger on the reader goroutine.
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
	capReg := regmocks.NewCapabilitiesRegistry(t)
	engines := newFakeEngineRegistry()
	wfLimits := &fakeWorkflowLimits{}
	clock := clockwork.NewFakeClock()

	wid, err := types.WorkflowIDFromHex(validWorkflowID)
	require.NoError(t, err)
	engine := &fakeEngine{coordinated: true}
	engines.set(wid, engine)

	c := NewCoordinator(newRegisterDeps(t, capReg, limits.NewGateLimiter(true)), engines, wfLimits, clock)
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
	f.c.mu.Lock()
	defer f.c.mu.Unlock()
	_, ok := f.c.workflows[validWorkflowID]
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

func TestCoordinator_RegisterTriggers(t *testing.T) {
	t.Parallel()

	t.Run("registers triggers and delivers events to the engine", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		var gotReq capabilities.TriggerRegistrationRequest
		eventCh := make(chan capabilities.TriggerResponse)
		trigger := capmocks.NewTriggerCapability(t)
		trigger.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).
			Run(func(_ context.Context, req capabilities.TriggerRegistrationRequest) { gotReq = req }).
			Return((<-chan capabilities.TriggerResponse)(eventCh), nil).Once()
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Maybe()
		f.capReg.EXPECT().GetTrigger(mock.Anything, testTriggerCapID).Return(trigger, nil).Once()

		triggerIDs, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)
		assert.Equal(t, []string{testTriggerCapID}, triggerIDs)
		assert.Equal(t, 1, f.limits.inUse())
		assert.True(t, f.registered())

		assert.Equal(t, RegistrationID(validWorkflowID, 0), gotReq.TriggerID)
		assert.Equal(t, validWorkflowID, gotReq.Metadata.WorkflowID)
		assert.Equal(t, testOwner, gotReq.Metadata.WorkflowOwner)
		assert.Equal(t, testDonID, gotReq.Metadata.WorkflowDonID)
		assert.Equal(t, uint32(pinnedWorkflowDonConfigVersion), gotReq.Metadata.WorkflowDonConfigVersion)

		eventCh <- triggerEvent("evt-1")
		require.Eventually(t, func() bool { return len(f.engine.executed()) == 1 }, 5*time.Second, 10*time.Millisecond)

		got := f.engine.executed()[0]
		assert.Equal(t, validWorkflowID, got.WorkflowID)
		assert.Equal(t, testTriggerCapID, got.TriggerCapID)
		assert.Equal(t, 0, got.TriggerIndex)
		assert.Equal(t, "evt-1", got.Event.Event.ID)
		assert.Equal(t, f.clock.Now(), got.ObservedAt)

		f.engine.mu.Lock()
		assert.Equal(t, testTenant, f.engine.tenants[0], "execution ctx must carry the workflow's CRE")
		f.engine.mu.Unlock()
	})

	t.Run("invalid workflow ID", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		sub := newTestSubscriber()
		sub.tenant.Workflow = "not-hex"

		_, err := f.c.RegisterTriggers(t.Context(), sub, testParams(t))
		require.ErrorContains(t, err, "invalid workflow id")
		assert.Equal(t, 0, sub.calls)
		assert.Equal(t, 0, f.limits.inUse())
	})

	t.Run("workflow limit reached maps to the scoped sentinel and never subscribes", func(t *testing.T) {
		t.Parallel()
		for _, tc := range []struct {
			scope settings.Scope
			want  error
		}{
			{settings.ScopeOwner, types.ErrPerOwnerWorkflowCountLimitReached},
			{settings.ScopeGlobal, types.ErrGlobalWorkflowCountLimitReached},
		} {
			f := newCoordinatorFixture(t)
			f.limits.useErr = limits.ErrorResourceLimited[int]{Scope: tc.scope, Limit: 1, Used: 1, Amount: 1}
			sub := newTestSubscriber()

			_, err := f.c.RegisterTriggers(t.Context(), sub, testParams(t))
			require.ErrorIs(t, err, tc.want)
			assert.Equal(t, 0, sub.calls)
			assert.False(t, f.registered())
		}
	})

	t.Run("the limit is acquired under the workflow's tenant", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		f.limits.useErr = errors.New("boom")

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.Error(t, err)
		f.limits.mu.Lock()
		defer f.limits.mu.Unlock()
		require.Len(t, f.limits.tenants, 1)
		assert.Equal(t, testTenant, f.limits.tenants[0])
	})

	t.Run("subscribe failure frees the limit", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		sub := newTestSubscriber()
		sub.err = errors.New("wasm boom")

		_, err := f.c.RegisterTriggers(t.Context(), sub, testParams(t))
		require.ErrorContains(t, err, "failed to subscribe to triggers")
		assert.Equal(t, 0, f.limits.inUse())
		assert.False(t, f.registered())
	})

	t.Run("registration failure frees the limit and leaves nothing registered", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		f.capReg.EXPECT().GetTrigger(mock.Anything, testTriggerCapID).Return(nil, errors.New("not found")).Once()

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.ErrorContains(t, err, "trigger capability not found")
		assert.Equal(t, 0, f.limits.inUse())
		assert.False(t, f.registered())
	})
}

func TestCoordinator_Deliver(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		mutate func(f *coordinatorFixture)
	}{
		{"drops events once the engine is gone from the registry", func(f *coordinatorFixture) {
			f.engines.mu.Lock()
			delete(f.engines.engines, f.wid)
			f.engines.mu.Unlock()
		}},
		{"drops events for an engine that is not coordinated", func(f *coordinatorFixture) {
			f.engine.coordinated = false
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := newCoordinatorFixture(t)
			trigger, eventCh := f.expectTrigger(t)
			trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Maybe()

			_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
			require.NoError(t, err)
			tc.mutate(f)

			// The channel is unbuffered: the second send only completes once the
			// reader has finished handling the first.
			eventCh <- triggerEvent("evt-1")
			eventCh <- triggerEvent("evt-2")
			assert.Empty(t, f.engine.executed())
		})
	}
}

func TestCoordinator_Ack(t *testing.T) {
	t.Parallel()

	t.Run("routes to the registration's handle", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		trigger, _ := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Maybe()
		regID := RegistrationID(validWorkflowID, 0)
		trigger.EXPECT().AckEvent(mock.Anything, regID, "evt-1", testMethod).Return(nil).Once()

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)
		require.NoError(t, f.c.Ack(t.Context(), testTriggerCapID, regID, "evt-1"))
	})

	t.Run("unknown workflow", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		err := f.c.Ack(t.Context(), testTriggerCapID, RegistrationID(validWorkflowID, 0), "evt-1")
		require.ErrorContains(t, err, "failed to find trigger")
	})

	t.Run("malformed registration ID", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		err := f.c.Ack(t.Context(), testTriggerCapID, "not-a-registration", "evt-1")
		require.ErrorContains(t, err, "invalid trigger registration ID")
	})
}

func TestCoordinator_UnregisterTriggers(t *testing.T) {
	t.Parallel()

	t.Run("unknown workflow", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		require.ErrorIs(t, f.c.UnregisterTriggers(validWorkflowID), ErrWorkflowNotCoordinated)
	})

	t.Run("releases handles once idle", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		trigger, _ := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.MatchedBy(func(req capabilities.TriggerRegistrationRequest) bool {
			return req.TriggerID == RegistrationID(validWorkflowID, 0) && req.Method == testMethod
		})).Return(nil).Once()

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)
		require.NoError(t, f.c.UnregisterTriggers(validWorkflowID))

		require.Eventually(t, func() bool { return !f.registered() && f.limits.inUse() == 0 }, 5*time.Second, 10*time.Millisecond)
		assert.Equal(t, 1, f.limits.freeCount())
		f.limits.mu.Lock()
		assert.Equal(t, testTenant, f.limits.tenants[len(f.limits.tenants)-1], "the limit is freed under the workflow's tenant")
		f.limits.mu.Unlock()
	})

	t.Run("keeps handles while an execution is in flight so it can still ACK", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		trigger, eventCh := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()
		regID := RegistrationID(validWorkflowID, 0)
		trigger.EXPECT().AckEvent(mock.Anything, regID, "evt-1", testMethod).Return(nil).Once()
		started, release := f.blockingExecution()

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)
		eventCh <- triggerEvent("evt-1")
		<-started

		require.NoError(t, f.c.UnregisterTriggers(validWorkflowID))

		// Ingress is stopped but the in-flight execution still resolves its handle.
		require.NoError(t, f.c.Ack(t.Context(), testTriggerCapID, regID, "evt-1"))
		assert.True(t, f.registered())
		assert.Equal(t, 1, f.limits.inUse())

		close(release)
		require.Eventually(t, func() bool { return !f.registered() && f.limits.inUse() == 0 }, 5*time.Second, 10*time.Millisecond)
	})

	t.Run("repeated calls unregister with the capability once", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		trigger, _ := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)
		require.NoError(t, f.c.UnregisterTriggers(validWorkflowID))
		// The second call may land before or after the release; both are fine.
		err = f.c.UnregisterTriggers(validWorkflowID)
		if err != nil {
			require.ErrorIs(t, err, ErrWorkflowNotCoordinated)
		}
		require.Eventually(t, func() bool { return f.limits.inUse() == 0 }, 5*time.Second, 10*time.Millisecond)
		assert.Equal(t, 1, f.limits.freeCount())
	})

	t.Run("a failed unregister is returned and retried on the next call", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		trigger, eventCh := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(errors.New("capability down")).Once()
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()
		// Keep a reader busy so the retry happens before the release.
		started, release := f.blockingExecution()

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)
		eventCh <- triggerEvent("evt-1")
		<-started

		require.ErrorContains(t, f.c.UnregisterTriggers(validWorkflowID), "failed to unregister 1 of 1 triggers")
		require.NoError(t, f.c.UnregisterTriggers(validWorkflowID))

		close(release)
		require.Eventually(t, func() bool { return f.limits.inUse() == 0 }, 5*time.Second, 10*time.Millisecond)
		assert.Equal(t, 1, f.limits.freeCount(), "only one release waiter may free the slot")
	})

	t.Run("releases after the drain timeout even if an execution hangs", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		trigger, eventCh := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()
		started, release := f.blockingExecution()
		defer close(release)

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)
		eventCh <- triggerEvent("evt-1")
		<-started

		require.NoError(t, f.c.UnregisterTriggers(validWorkflowID))
		require.NoError(t, f.clock.BlockUntilContext(t.Context(), 1))
		assert.True(t, f.registered())

		f.clock.Advance(defaultDrainTimeout)
		require.Eventually(t, func() bool { return !f.registered() && f.limits.inUse() == 0 }, 5*time.Second, 10*time.Millisecond)
	})
}

func TestCoordinator_ReRegisterWhileDraining(t *testing.T) {
	t.Parallel()
	f := newCoordinatorFixture(t)

	oldTrigger, oldCh := f.expectTrigger(t)
	oldTrigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()
	started, release := f.blockingExecution()

	_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
	require.NoError(t, err)
	oldCh <- triggerEvent("evt-1")
	<-started
	require.NoError(t, f.c.UnregisterTriggers(validWorkflowID))

	f.c.mu.Lock()
	oldWT := f.c.workflows[validWorkflowID]
	f.c.mu.Unlock()

	// The replacement engine registers while the old registration is still draining.
	f.engine.execute = nil
	newTrigger, _ := f.expectTrigger(t)
	newTrigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Maybe()
	_, err = f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
	require.NoError(t, err)

	f.c.mu.Lock()
	newWT := f.c.workflows[validWorkflowID]
	f.c.mu.Unlock()
	require.NotSame(t, oldWT, newWT, "registering again must replace, not reuse, the old state")
	assert.Equal(t, 2, f.limits.inUse(), "the old slot is held until its drain completes")

	close(release)
	oldWT.readers.Wait() // the old registration's own readers have now exited

	// The old release waiter, woken by the same Wait(), may still run and
	// (wrongly, if buggy) drop the new registration. Poll instead of a single
	// sleep-then-check: this fails the instant the bug appears rather than
	// only if it happens to land inside a guessed window.
	require.Never(t, func() bool {
		f.c.mu.Lock()
		defer f.c.mu.Unlock()
		return f.c.workflows[validWorkflowID] != newWT
	}, 100*time.Millisecond, 2*time.Millisecond, "the old waiter must not drop the new registration")
	require.Eventually(t, func() bool { return f.limits.inUse() == 1 }, time.Second, time.Millisecond,
		"the old waiter must still free its own slot")

	regID := RegistrationID(validWorkflowID, 0)
	newTrigger.EXPECT().AckEvent(mock.Anything, regID, "evt-2", testMethod).Return(nil).Once()
	require.NoError(t, f.c.Ack(t.Context(), testTriggerCapID, regID, "evt-2"))
}
