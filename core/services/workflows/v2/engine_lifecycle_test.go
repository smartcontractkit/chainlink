package v2_test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonmetrics "github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	metmocks "github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
)

// lifecycleEngineMode is the kind of engine a lifecycle test runs against.
type lifecycleEngineMode string

const (
	// lifecycleLegacy is the engine that registers its own triggers and runs its own queue.
	lifecycleLegacy lifecycleEngineMode = "legacy"
	// lifecycleCoordinated is the execution-only engine driven by a real trigger coordinator.
	lifecycleCoordinated lifecycleEngineMode = "coordinated"
)

// lifecycleTriggerCapID is the trigger every lifecycle workflow subscribes to.
const lifecycleTriggerCapID = "id_0"

const testTimeout = 10 * time.Second

// lifecycleEvent is something the engine did that a lifecycle test orders its assertions on.
type lifecycleEvent string

const (
	// eventExecutionCanceled is recorded when an execution stops because its context was canceled.
	eventExecutionCanceled lifecycleEvent = "execution-canceled"

	// eventExecutionFinished is recorded by the OnExecutionFinished hook, which the engine runs
	// after the module returns and before the execution call returns.
	eventExecutionFinished lifecycleEvent = "execution-finished"

	// eventModuleClosed is recorded when the engine closes its module.
	eventModuleClosed lifecycleEvent = "module-closed"

	// eventCloseReturned is recorded when the engine's Close returns.
	eventCloseReturned lifecycleEvent = "close-returned"
)

// TestEngine_Lifecycle runs the same Close and Drain scenarios against both engines. Each
// scenario asserts the behavior every engine is expected to have, so the two must not diverge.
//
// The scenarios hold executions in the module on channels and order their assertions on a log of
// engine events.
func TestEngine_Lifecycle(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name string
		run  func(*testing.T, lifecycleEngineMode)
	}{
		{"Close_Idle", testLifecycleCloseIdle},
		{"Close_CancelsInFlightExecution", testLifecycleCloseCancelsInFlightExecution},
		{"Close_WaitsForExecutionBeforeClosingModule", testLifecycleCloseWaitsForExecutionBeforeClosingModule},
		{"ExecuteTrigger_AfterClose", testLifecycleExecuteTriggerAfterClose},
		{"ExecuteTrigger_DuringClose", testLifecycleExecuteTriggerDuringClose},
		{"ExecuteTrigger_ConcurrentWithClose", testLifecycleExecuteTriggerConcurrentWithClose},
		{"ExecuteTrigger_WhileDraining", testLifecycleExecuteTriggerWhileDraining},
		{"Drain_KeepsInFlightExecutionRunning", testLifecycleDrainKeepsInFlightExecutionRunning},
		{"Drain_ThenClose_CancelsInFlightExecution", testLifecycleDrainThenCloseCancelsInFlightExecution},
	}
	for _, mode := range []lifecycleEngineMode{lifecycleLegacy, lifecycleCoordinated} {
		t.Run(string(mode), func(t *testing.T) {
			t.Parallel()
			for _, s := range scenarios {
				t.Run(s.name, func(t *testing.T) {
					t.Parallel()
					s.run(t, mode)
				})
			}
		})
	}
}

// testLifecycleCloseIdle closes an engine with nothing in flight. It asserts that Close succeeds,
// closes the module once, and that a second Close reports the engine is already stopped without
// closing the module again.
func testLifecycleCloseIdle(t *testing.T, mode lifecycleEngineMode) {
	h := newEngineLifecycleHarness(t, mode)
	h.start()

	require.NoError(t, h.close())
	require.Equal(t, int32(1), h.moduleCloses.Load(), "module Close calls after the first engine Close")

	require.ErrorIs(t, h.close(), services.ErrAlreadyStopped)
	require.Equal(t, int32(1), h.moduleCloses.Load(), "module Close calls after the second engine Close")
}

// testLifecycleCloseCancelsInFlightExecution closes an engine while an execution is running in
// the module. The execution returns only when its context is canceled, so it asserts that Close
// cancels it, rather than leaving it to run on or waiting for it to finish by itself.
func testLifecycleCloseCancelsInFlightExecution(t *testing.T, mode lifecycleEngineMode) {
	h := newEngineLifecycleHarness(t, mode)
	h.start()
	h.startExecution()

	require.NoError(t, h.close())

	require.Contains(t, h.log.events(), eventExecutionCanceled)
}

// testLifecycleCloseWaitsForExecutionBeforeClosingModule closes an engine while an execution is
// running in the module. It asserts that the execution finishes completely, finalizers included,
// before the module is closed, and that Close returns only after both.
func testLifecycleCloseWaitsForExecutionBeforeClosingModule(t *testing.T, mode lifecycleEngineMode) {
	h := newEngineLifecycleHarness(t, mode)
	h.start()
	h.startExecution()

	require.NoError(t, h.close())

	requireEventOrder(t, h.log.events(), eventExecutionFinished, eventModuleClosed, eventCloseReturned)
}

// testLifecycleExecuteTriggerAfterClose delivers a trigger event directly to a closed engine. It
// asserts that the event is rejected without ever reaching the module. The gate is open, so an
// engine that wrongly runs the event finishes it and is counted, rather than holding the test.
func testLifecycleExecuteTriggerAfterClose(t *testing.T, mode lifecycleEngineMode) {
	h := newEngineLifecycleHarness(t, mode)
	h.start()
	require.NoError(t, h.close())
	h.gate.open()

	err := h.engine.ExecuteTrigger(h.tenantCtx(), h.event("after-close"))

	require.ErrorIs(t, err, v2.ErrEngineClosed, "ExecuteTrigger on a closed engine")
	require.Equal(t, int32(0), h.moduleExecutions.Load(), "executions that reached the module")
}

// testLifecycleExecuteTriggerDuringClose delivers a trigger event directly to an engine whose Close
// is pending. An execution held in the module is kept there after Close cancels it, so Close, which
// waits for its in-flight executions, cannot return until the test lets the execution go. It asserts
// that the event is rejected without reaching the module.
func testLifecycleExecuteTriggerDuringClose(t *testing.T, mode lifecycleEngineMode) {
	h := newEngineLifecycleHarness(t, mode)
	h.start()
	h.gate.holdCanceled()
	h.startExecution()

	closed := make(chan error, 1)
	go func() {
		err := h.engine.Close()
		h.log.record(eventCloseReturned)
		closed <- err
	}()
	awaitLifecycle(t, h.gate.canceled, "Close to cancel the in-flight execution")
	require.Empty(t, closed, "Close returned while an execution was still in flight")

	err := h.engine.ExecuteTrigger(h.tenantCtx(), h.event("during-close"))

	require.ErrorIs(t, err, v2.ErrEngineClosed, "ExecuteTrigger while the engine is closing")
	require.Equal(t, int32(1), h.moduleExecutions.Load(), "executions that reached the module")

	h.gate.finishCanceled()
	require.NoError(t, awaitLifecycle(t, closed, "engine Close to return"))
}

// testLifecycleExecuteTriggerConcurrentWithClose delivers trigger events from many goroutines
// while the engine is closed. Every call races Close, so each may run to completion or be rejected.
// It asserts that Close returns, that a rejected call reports the engine is closed, and that the
// module is closed once. It is meant to be run with the race detector.
func testLifecycleExecuteTriggerConcurrentWithClose(t *testing.T, mode lifecycleEngineMode) {
	const callers = 16

	h := newEngineLifecycleHarness(t, mode)
	h.start()
	h.gate.open()

	errs := make([]error, callers)
	var wg sync.WaitGroup
	for i := range callers {
		wg.Go(func() {
			errs[i] = h.engine.ExecuteTrigger(h.tenantCtx(), h.event(fmt.Sprintf("racing-%d", i)))
		})
	}
	closeErr := h.close()
	wg.Wait()

	require.NoError(t, closeErr)
	for i, err := range errs {
		if err != nil {
			require.ErrorIs(t, err, v2.ErrEngineClosed, "ExecuteTrigger call %d", i)
		}
	}
	require.Equal(t, int32(1), h.moduleCloses.Load(), "module Close calls")
}

// testLifecycleExecuteTriggerWhileDraining delivers a trigger event directly to a draining engine.
// It asserts that the event is rejected without ever reaching the module. The gate is open, so an
// engine that wrongly runs the event finishes it and is counted, rather than holding the test.
func testLifecycleExecuteTriggerWhileDraining(t *testing.T, mode lifecycleEngineMode) {
	h := newEngineLifecycleHarness(t, mode)
	h.start()
	require.True(t, h.engine.Drain(), "first Drain")
	h.gate.open()

	err := h.engine.ExecuteTrigger(h.tenantCtx(), h.event("while-draining"))

	require.ErrorIs(t, err, v2.ErrEngineDraining, "ExecuteTrigger on a draining engine")
	require.Equal(t, int32(0), h.moduleExecutions.Load(), "executions that reached the module")
}

// testLifecycleDrainKeepsInFlightExecutionRunning drains an engine while an execution is running.
// It asserts that Drain reports the transition once and that the execution, when it resumes, still
// has a live context and finishes normally.
func testLifecycleDrainKeepsInFlightExecutionRunning(t *testing.T, mode lifecycleEngineMode) {
	h := newEngineLifecycleHarness(t, mode)
	h.start()
	h.startExecution()

	require.True(t, h.engine.Drain(), "first Drain")
	require.False(t, h.engine.Drain(), "second Drain")
	_, draining := h.engine.DrainStartedAt()
	require.True(t, draining, "DrainStartedAt reports a drain")
	require.Equal(t, int32(1), h.engine.ActiveExecutions(), "active executions while draining")

	h.gate.open()

	require.Equal(t, "completed", h.waitExecutionFinished())
	require.NoError(t, h.gate.contextErrAtRelease(), "the execution's context when it resumed after Drain")
	require.NotContains(t, h.log.events(), eventExecutionCanceled)
}

// testLifecycleDrainThenCloseCancelsInFlightExecution drains an engine and then closes it while an
// execution is still running. It asserts that Close cancels the execution, and that the module is
// closed only after the execution has finished.
func testLifecycleDrainThenCloseCancelsInFlightExecution(t *testing.T, mode lifecycleEngineMode) {
	h := newEngineLifecycleHarness(t, mode)
	h.start()
	h.startExecution()
	require.True(t, h.engine.Drain())

	require.NoError(t, h.close())

	requireEventOrder(t, h.log.events(),
		eventExecutionCanceled, eventExecutionFinished, eventModuleClosed, eventCloseReturned)
}

// engineLifecycleHarness is a real engine of the harness's mode over a mock module whose
// executions are held on a gate until the test releases them or their context is canceled. In
// coordinated mode a real trigger coordinator registers the engine's triggers, delivers its events
// and acknowledges them; in legacy mode the engine does all of that itself.
type engineLifecycleHarness struct {
	t      *testing.T
	mode   lifecycleEngineMode
	cfg    *v2.EngineConfig
	engine v2.WorkflowEngine

	coordinator triggers.Coordinator // nil in legacy mode
	engines     *lifecycleEngineRegistry
	eventCh     chan capabilities.TriggerResponse
	gate        *lifecycleGate
	log         *lifecycleLog

	initDoneCh   chan error
	subscribedCh chan []string
	finishedCh   chan string // status of each finished execution

	moduleCloses     atomic.Int32
	moduleExecutions atomic.Int32 // executions that reached the module, subscribe calls excluded
	nextEvent        atomic.Int32
}

// newEngineLifecycleHarness builds the engine for mode without starting it. Whatever the test
// leaves running is released and closed at cleanup.
func newEngineLifecycleHarness(t *testing.T, mode lifecycleEngineMode) *engineLifecycleHarness {
	t.Helper()

	log := &lifecycleLog{}
	h := &engineLifecycleHarness{
		t:            t,
		mode:         mode,
		engines:      &lifecycleEngineRegistry{},
		eventCh:      make(chan capabilities.TriggerResponse),
		gate:         newLifecycleGate(log),
		log:          log,
		initDoneCh:   make(chan error, 1),
		subscribedCh: make(chan []string, 1),
		finishedCh:   make(chan string, 16),
	}

	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil).Maybe()
	trigger := capmocks.NewTriggerCapability(t)
	capreg.EXPECT().GetTrigger(matches.AnyContext, lifecycleTriggerCapID).Return(trigger, nil).Maybe()
	trigger.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(h.eventCh, nil).Maybe()
	trigger.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Maybe()
	trigger.EXPECT().AckEvent(matches.AnyContext, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	cfg := defaultTestConfig(t, nil)
	cfg.CapRegistry = capreg
	cfg.Module = h.newModule()
	cfg.BillingClient = newLifecycleBillingClient(t)
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized:          func(err error) { nonBlockingSend(h.initDoneCh, err) },
		OnSubscribedToTriggers: func(ids []string) { nonBlockingSend(h.subscribedCh, ids) },
		OnExecutionFinished: func(_, status string) {
			h.log.record(eventExecutionFinished)
			nonBlockingSend(h.finishedCh, status)
		},
	}
	h.cfg = cfg

	construct := v2.NewEngine
	if mode == lifecycleCoordinated {
		h.coordinator = h.newCoordinator(capreg)
		cfg.TriggerAcknowledger = h.coordinator
		construct = v2.NewCoordinatedEngine
	}
	engine, err := construct(cfg)
	require.NoError(t, err)
	h.engine = engine
	// Set before the engine starts, so the coordinator's readers never see it unset.
	h.engines.engine = engine

	// Open the gate first so a failed test's held execution does not hold up Close, then close
	// the engine before the mocks and limiters it uses go away.
	t.Cleanup(func() {
		h.gate.open()
		_ = engine.Close()
	})
	return h
}

// newModule returns the module mock. It answers the subscribe call with one trigger and holds
// every other execution on the gate, counting executions and logging when the module is closed.
func (h *engineLifecycleHarness) newModule() *modulemocks.ModuleV2 {
	module := modulemocks.NewModuleV2(h.t)
	module.EXPECT().Start().Maybe()
	module.EXPECT().Close().Run(func() {
		h.log.record(eventModuleClosed)
		h.moduleCloses.Add(1)
	}).Maybe()
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, req *sdkpb.ExecuteRequest, _ host.ExecutionHelper) (*sdkpb.ExecutionResult, error) {
			if _, ok := req.Request.(*sdkpb.ExecuteRequest_Subscribe); ok {
				return newTriggerSubs(1), nil
			}
			h.moduleExecutions.Add(1)
			if err := h.gate.block(ctx); err != nil {
				return nil, err
			}
			return &sdkpb.ExecutionResult{Result: &sdkpb.ExecutionResult_Value{}}, nil
		}).Maybe()
	return module
}

// newCoordinator starts a real trigger coordinator whose engine registry holds only the engine
// under test.
func (h *engineLifecycleHarness) newCoordinator(capreg *regmocks.CapabilitiesRegistry) triggers.Coordinator {
	h.t.Helper()
	lggr := logger.TestLogger(h.t)
	workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: 10, PerOwner: 10},
		limits.Factory{Logger: lggr})
	require.NoError(h.t, err)
	h.t.Cleanup(func() { require.NoError(h.t, workflowLimits.Close()) })

	em, err := monitoring.InitMonitoringResources()
	require.NoError(h.t, err)
	coordinator := triggers.NewCoordinator(
		triggers.RegisterDeps{
			CapRegistry:  capreg,
			RegTimeout:   limits.NewTimeLimiter(5 * time.Second),
			ChainAllowed: limits.NewGateLimiter(true),
			Logger:       commonlogger.Test(h.t),
			Metrics:      monitoring.NewWorkflowsMetricLabeler(commonmetrics.NewLabeler(), em),
		},
		h.engines,
		workflowLimits,
		clockwork.NewFakeClock(),
	)
	servicetest.Run(h.t, coordinator)
	return coordinator
}

// start starts the engine and waits until its triggers are registered, so a trigger event
// delivered next reaches it.
func (h *engineLifecycleHarness) start() {
	h.t.Helper()
	require.NoError(h.t, h.engine.Start(h.t.Context()))
	require.NoError(h.t, awaitLifecycle(h.t, h.initDoneCh, "engine to initialize"))

	if h.mode == lifecycleCoordinated {
		ids, err := h.coordinator.RegisterTriggers(h.t.Context(), h.engine, triggers.RegistrationParams{
			WorkflowOwner:                 h.cfg.WorkflowOwner,
			WorkflowName:                  h.cfg.WorkflowName,
			WorkflowTag:                   h.cfg.WorkflowTag,
			WorkflowRegistryChainSelector: h.cfg.WorkflowRegistryChainSelector,
			WorkflowRegistryAddress:       h.cfg.WorkflowRegistryAddress,
		})
		require.NoError(h.t, err)
		require.Equal(h.t, []string{lifecycleTriggerCapID}, ids)
		return
	}
	require.Equal(h.t, []string{lifecycleTriggerCapID}, awaitLifecycle(h.t, h.subscribedCh, "engine to register triggers"))
}

// startExecution delivers a trigger event through the engine's real delivery path and waits until
// its execution is running in the module. It stays there until the gate opens or its context is
// canceled.
func (h *engineLifecycleHarness) startExecution() {
	h.t.Helper()
	event := h.event(fmt.Sprintf("event-%d", h.nextEvent.Add(1))).Event
	select {
	case h.eventCh <- event:
	case <-time.After(testTimeout):
		require.FailNow(h.t, "trigger event was not accepted")
	}
	awaitLifecycle(h.t, h.gate.started, "execution to start")
}

// close closes the engine, logging when Close returns, and returns its error. Close runs on its
// own goroutine so an engine that never returns fails the test instead of hanging it.
func (h *engineLifecycleHarness) close() error {
	h.t.Helper()
	done := make(chan error, 1)
	go func() {
		err := h.engine.Close()
		h.log.record(eventCloseReturned)
		done <- err
	}()
	return awaitLifecycle(h.t, done, "engine Close to return")
}

// waitExecutionFinished returns the status of the next execution to finish.
func (h *engineLifecycleHarness) waitExecutionFinished() string {
	h.t.Helper()
	return awaitLifecycle(h.t, h.finishedCh, "execution to finish")
}

// event returns a trigger event for the workflow's one trigger, for delivery straight to the engine.
func (h *engineLifecycleHarness) event(id string) triggers.CoordinatedEvent {
	return triggers.CoordinatedEvent{
		WorkflowID:   h.cfg.WorkflowID,
		TriggerCapID: lifecycleTriggerCapID,
		TriggerIndex: 0,
		ObservedAt:   time.Now(),
		Event: capabilities.TriggerResponse{
			Event: capabilities.TriggerEvent{TriggerType: "basic-trigger@1.0.0", ID: id},
		},
	}
}

// tenantCtx returns a context carrying the workflow's tenant, as the engine's callers give it.
func (h *engineLifecycleHarness) tenantCtx() context.Context {
	return contexts.WithCRE(h.t.Context(), h.engine.Tenant())
}

// awaitLifecycle returns the next value from ch, failing the test if none arrives within the
// backstop.
func awaitLifecycle[T any](t *testing.T, ch <-chan T, what string) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(testTimeout):
		require.FailNow(t, "timed out waiting for "+what)
		panic("unreachable")
	}
}

// requireEventOrder fails the test unless every event in want appears in got, in that order.
// Other events may appear in between.
func requireEventOrder(t *testing.T, got []lifecycleEvent, want ...lifecycleEvent) {
	t.Helper()
	next := 0
	for _, e := range got {
		if next < len(want) && e == want[next] {
			next++
		}
	}
	require.Equal(t, len(want), next, "events %v do not contain %v in order", got, want)
}

// lifecycleLog records the events a scenario orders its assertions on.
type lifecycleLog struct {
	mu     sync.Mutex
	logged []lifecycleEvent
}

func (l *lifecycleLog) record(e lifecycleEvent) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logged = append(l.logged, e)
}

// events returns the events recorded so far, in order.
func (l *lifecycleLog) events() []lifecycleEvent {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]lifecycleEvent(nil), l.logged...)
}

// lifecycleEngineRegistry is a triggers.EngineRegistry holding the one engine under test.
type lifecycleEngineRegistry struct {
	engine v2.WorkflowEngine
}

func (r *lifecycleEngineRegistry) Get(wid types.WorkflowID) (triggers.RegisteredEngine, bool) {
	if r.engine == nil || wid.Hex() != r.engine.Tenant().Workflow {
		return nil, false
	}
	return r.engine, true
}

// lifecycleGate holds executions in the module until it is opened or their context is canceled.
// It has no timeout: an execution nothing releases or cancels stays held until the test ends.
type lifecycleGate struct {
	started  chan struct{}
	canceled chan struct{} // signaled when a held execution observes its context's cancellation
	release  chan struct{}
	finish   chan struct{} // holds a canceled execution, once holdCanceled is set, until it is closed
	once     sync.Once
	finOnce  sync.Once
	hold     atomic.Bool
	log      *lifecycleLog

	mu           sync.Mutex
	errAtRelease error // the context's error when a released execution resumed
}

func newLifecycleGate(log *lifecycleLog) *lifecycleGate {
	return &lifecycleGate{
		started:  make(chan struct{}, 16),
		canceled: make(chan struct{}, 16),
		release:  make(chan struct{}),
		finish:   make(chan struct{}),
		log:      log,
	}
}

// holdCanceled makes the next execution that observes its context's cancellation stay held, after
// signaling canceled, until finishCanceled or open is called. It lets a test keep an engine's Close
// pending, since Close waits for its in-flight executions. Only one execution is held, so an
// execution that should have been rejected fails the test instead of hanging it.
func (g *lifecycleGate) holdCanceled() { g.hold.Store(true) }

// finishCanceled lets executions held by holdCanceled return.
func (g *lifecycleGate) finishCanceled() { g.finOnce.Do(func() { close(g.finish) }) }

// block signals that an execution started and holds it. It returns nil once the gate opens, and
// the context's error, after logging the cancellation, if the context is canceled first.
func (g *lifecycleGate) block(ctx context.Context) error {
	nonBlockingSend(g.started, struct{}{})

	select {
	case <-g.release:
		g.mu.Lock()
		g.errAtRelease = ctx.Err()
		g.mu.Unlock()
		return nil
	case <-ctx.Done():
		g.log.record(eventExecutionCanceled)
		nonBlockingSend(g.canceled, struct{}{})
		if g.hold.CompareAndSwap(true, false) {
			<-g.finish
		}
		return ctx.Err()
	}
}

// open releases every held execution, including ones held after cancellation.
func (g *lifecycleGate) open() {
	g.once.Do(func() { close(g.release) })
	g.finishCanceled()
}

// contextErrAtRelease returns the context's error as seen by an execution that resumed when the
// gate opened.
func (g *lifecycleGate) contextErrAtRelease() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.errAtRelease
}

// nonBlockingSend sends v on ch unless it is full, so a hook never stalls the engine.
func nonBlockingSend[T any](ch chan<- T, v T) {
	select {
	case ch <- v:
	default:
	}
}

// newLifecycleBillingClient returns a billing client that accepts, but does not require, the calls
// an execution makes, since not every scenario runs one.
func newLifecycleBillingClient(t *testing.T) *metmocks.BillingClient {
	t.Helper()
	client := metmocks.NewBillingClient(t)
	client.EXPECT().GetWorkflowExecutionRates(mock.Anything, mock.Anything).
		Return(&billing.GetWorkflowExecutionRatesResponse{
			RateCards: []*billing.RateCard{
				{
					ResourceType:    billing.ResourceType_RESOURCE_TYPE_COMPUTE,
					MeasurementUnit: billing.MeasurementUnit_MEASUREMENT_UNIT_MILLISECONDS,
					UnitsPerCredit:  "0.0001",
				},
				{
					ResourceType:    billing.ResourceType_RESOURCE_TYPE_NETWORK,
					MeasurementUnit: billing.MeasurementUnit_MEASUREMENT_UNIT_COST,
					UnitsPerCredit:  "0.01",
				},
			},
		}, nil).Maybe()
	client.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
		Return(&billing.ReserveCreditsResponse{Success: true, Credits: "10000"}, nil).Maybe()
	client.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).Return(&emptypb.Empty{}, nil).Maybe()
	return client
}
