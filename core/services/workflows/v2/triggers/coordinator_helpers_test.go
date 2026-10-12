package triggers

import (
	"context"
	"sync"
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
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
