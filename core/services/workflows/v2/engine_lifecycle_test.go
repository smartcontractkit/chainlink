package v2_test

import (
	"context"
	"errors"
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
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	metmocks "github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
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

// inFlightShutdown is how an InFlightTriggers scenario shuts the engine down while trigger events
// are still arriving.
type inFlightShutdown int

const (
	// shutdownClose closes the engine directly.
	shutdownClose inFlightShutdown = iota + 1

	// shutdownDrainThenClose drains the engine and closes it straight away.
	shutdownDrainThenClose

	// shutdownDrainWaitThenClose drains the engine, waits for its active executions to finish, then
	// closes it.
	shutdownDrainWaitThenClose
)

// TestEngine_Lifecycle runs the same Close, Drain and replacement scenarios against both engines. Each
// scenario asserts the behavior every engine is expected to have, so the two must not diverge.
//
// The scenarios hold executions in the module on channels, or let them run to completion, and order
// their assertions on a log of engine events. The concurrent scenarios are meant to be run with the
// race detector.
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
		{"ExecuteTrigger_ConcurrentWithClose", testLifecycleExecuteTriggerConcurrentWithClose},
		{"InFlightTriggers_Close", testLifecycleInFlightTriggersClose},
		{"InFlightTriggers_DrainThenClose", testLifecycleInFlightTriggersDrainThenClose},
		{"InFlightTriggers_DrainWaitThenClose", testLifecycleInFlightTriggersDrainWaitThenClose},
		{"ExecuteTrigger_ConcurrentWithDrainThenClose", testLifecycleExecuteTriggerConcurrentWithDrainThenClose},
		{"ExecuteTrigger_WhileDraining", testLifecycleExecuteTriggerWhileDraining},
		{"Drain_KeepsInFlightExecutionRunning", testLifecycleDrainKeepsInFlightExecutionRunning},
		{"Drain_ThenClose_CancelsInFlightExecution", testLifecycleDrainThenCloseCancelsInFlightExecution},
		{"Replace_NextEventRunsOnNewEngine", testLifecycleReplaceNextEventRunsOnNewEngine},
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

// testLifecycleExecuteTriggerAfterClose calls ExecuteTrigger directly on a closed engine. It asserts
// that the event is rejected without ever reaching the module.
func testLifecycleExecuteTriggerAfterClose(t *testing.T, mode lifecycleEngineMode) {
	h := newEngineLifecycleHarness(t, mode)
	h.start()
	require.NoError(t, h.close())

	err := h.engine.ExecuteTrigger(h.tenantCtx(), h.event("after-close"))

	require.ErrorIs(t, err, v2.ErrEngineClosed, "ExecuteTrigger on a closed engine")
	require.Equal(t, int32(0), h.moduleExecutions.Load(), "executions that reached the module")
}

// testLifecycleExecuteTriggerConcurrentWithClose calls ExecuteTrigger directly from many test
// goroutines while an engine is closed. Every call races Close, so each may run to completion or be
// rejected. It asserts that Close returns, that a rejected call reports the engine is closed, and
// that the module is closed once.
func testLifecycleExecuteTriggerConcurrentWithClose(t *testing.T, mode lifecycleEngineMode) {
	const callers = 16

	h := newEngineLifecycleHarness(t, mode)
	h.start()
	h.blocker.release()

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

// testLifecycleInFlightTriggersClose closes an engine while trigger events are still arriving
// through its real delivery path, so Close races the executions the engine starts for itself rather
// than calls from test goroutines.
func testLifecycleInFlightTriggersClose(t *testing.T, mode lifecycleEngineMode) {
	testLifecycleInFlightTriggersShutdown(t, mode, shutdownClose)
}

// testLifecycleInFlightTriggersDrainThenClose drains an engine and closes it straight away while
// trigger events are still arriving through its real delivery path.
func testLifecycleInFlightTriggersDrainThenClose(t *testing.T, mode lifecycleEngineMode) {
	testLifecycleInFlightTriggersShutdown(t, mode, shutdownDrainThenClose)
}

// testLifecycleInFlightTriggersDrainWaitThenClose drains an engine, waits for its active executions
// to finish, then closes it, while trigger events are still arriving through its real delivery path.
func testLifecycleInFlightTriggersDrainWaitThenClose(t *testing.T, mode lifecycleEngineMode) {
	testLifecycleInFlightTriggersShutdown(t, mode, shutdownDrainWaitThenClose)
}

// testLifecycleInFlightTriggersShutdown shuts an engine down as shutdown describes while a feeder
// keeps sending trigger events through its real delivery path. Executions run to completion when
// they reach the module, keeping the delivery path busy. It asserts that Drain, when called, reports
// the transition, that Close returns, that the module is closed once and only after every execution
// that reached it has finished, and that no execution is left active.
func testLifecycleInFlightTriggersShutdown(t *testing.T, mode lifecycleEngineMode, shutdown inFlightShutdown) {
	// executions that must reach the module before shutdown, so delivery is under way when it starts
	const warmup = 4

	h := newEngineLifecycleHarness(t, mode)
	h.start()
	h.blocker.release()

	// The feeder also stops when the test ends, so a test that fails before closing stop doesn't
	// leave it running.
	stop := make(chan struct{})
	fed := h.feedEvents(stop)
	for range warmup {
		awaitLifecycle(t, h.blocker.started, "fed execution to start")
	}

	if shutdown != shutdownClose {
		require.True(t, h.engine.Drain(), "Drain")
	}
	if shutdown == shutdownDrainWaitThenClose {
		require.Eventually(t, func() bool { return h.engine.ActiveExecutions() == 0 },
			testutils.WaitTimeout(t), 10*time.Millisecond, "active executions to finish after Drain")
	}
	require.NoError(t, h.close())
	close(stop)
	awaitLifecycle(t, fed, "feeder to exit")

	require.Equal(t, int32(1), h.moduleCloses.Load(), "module Close calls")
	require.Equal(t, int32(0), h.engine.ActiveExecutions(), "active executions after Close")
	events := h.log.events()
	requireAllBefore(t, events, eventExecutionFinished, eventModuleClosed)
	requireEventOrder(t, events, eventModuleClosed, eventCloseReturned)
}

// testLifecycleExecuteTriggerConcurrentWithDrainThenClose runs the standard shutdown, Drain then
// Close, while many direct ExecuteTrigger calls from test goroutines are in flight or arriving. The
// executions are held in the module, so Close has to cancel them. A first wave is held in the module
// before shutdown begins. A second wave races the shutdown, so each of its calls may be admitted and
// canceled, or rejected as draining or closed. It asserts that Close cancels the held executions,
// that no call fails any way but a rejection (a canceled execution reports through its status, not
// its error), that every execution that reached the module finishes before the module is closed,
// and that Close returns.
func testLifecycleExecuteTriggerConcurrentWithDrainThenClose(t *testing.T, mode lifecycleEngineMode) {
	const (
		inflight = 50
		racing   = inflight * 2
	)

	h := newEngineLifecycleHarness(t, mode, withSignalBuffer(inflight+racing))
	h.start()

	inflightErrs := make([]error, inflight)
	var inflightWG sync.WaitGroup
	for i := range inflight {
		inflightWG.Go(func() {
			inflightErrs[i] = h.engine.ExecuteTrigger(h.tenantCtx(), h.event(fmt.Sprintf("inflight-%d", i)))
		})
	}
	for range inflight {
		awaitLifecycle(h.t, h.blocker.started, "in-flight execution to start")
	}

	// Release the racing callers and the shutdown together.
	begin := make(chan struct{})
	racingErrs := make([]error, racing)
	var racingWG sync.WaitGroup
	for i := range racing {
		racingWG.Go(func() {
			<-begin
			racingErrs[i] = h.engine.ExecuteTrigger(h.tenantCtx(), h.event(fmt.Sprintf("racing-%d", i)))
		})
	}
	var drained bool
	var closeErr error
	var shutdownWG sync.WaitGroup
	shutdownWG.Go(func() {
		<-begin
		drained = h.engine.Drain()
		closeErr = h.close()
	})
	close(begin)

	shutdownWG.Wait()
	racingWG.Wait()
	inflightWG.Wait()

	require.True(t, drained, "Drain")
	require.NoError(t, closeErr)
	require.Equal(t, int32(1), h.moduleCloses.Load(), "module Close calls")

	for i, err := range inflightErrs {
		require.NoError(t, err, "in-flight ExecuteTrigger call %d", i)
	}
	for i, err := range racingErrs {
		if err != nil && !errors.Is(err, v2.ErrEngineDraining) && !errors.Is(err, v2.ErrEngineClosed) {
			require.NoError(t, err, "racing ExecuteTrigger call %d must succeed or be rejected", i)
		}
	}

	executions := h.moduleExecutions.Load()
	require.GreaterOrEqual(t, executions, int32(inflight), "executions that reached the module")
	require.LessOrEqual(t, executions, int32(inflight+racing), "executions that reached the module")

	// The blocker is never released, so every execution that reached the module was held until
	// Close canceled it, which the engine reports as an errored status.
	statuses := h.finishedStatuses()
	require.Len(t, statuses, int(executions), "executions that finished")
	for i, status := range statuses {
		require.Equal(t, store.StatusErrored, status, "status of finished execution %d", i)
	}
	require.Equal(t, int32(0), h.engine.ActiveExecutions(), "active executions after Close")

	events := h.log.events()
	require.Contains(t, events, eventExecutionCanceled)
	requireAllBefore(t, events, eventExecutionFinished, eventModuleClosed)
	requireEventOrder(t, events, eventModuleClosed, eventCloseReturned)
}

// testLifecycleExecuteTriggerWhileDraining calls ExecuteTrigger directly on a draining engine. It
// asserts that the event is rejected without ever reaching the module.
func testLifecycleExecuteTriggerWhileDraining(t *testing.T, mode lifecycleEngineMode) {
	h := newEngineLifecycleHarness(t, mode)
	h.start()
	require.True(t, h.engine.Drain(), "first Drain")

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

	h.blocker.release()

	require.Equal(t, "completed", h.waitExecutionFinished())
	require.NoError(t, h.blocker.contextErrAtRelease(), "the execution's context when it resumed after Drain")
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

// testLifecycleReplaceNextEventRunsOnNewEngine closes an engine and replaces it with a new one for
// the same workflow. It asserts that the next trigger event runs on the new engine. In coordinated
// mode the workflow's trigger registration stays live across the swap, so the coordinator must
// route the event to the new engine; in legacy mode the new engine registers its own triggers.
func testLifecycleReplaceNextEventRunsOnNewEngine(t *testing.T, mode lifecycleEngineMode) {
	h := newEngineLifecycleHarness(t, mode)
	h.start()
	h.blocker.release()

	h.deliver(admissionTimeout)
	require.Eventually(t, func() bool { return h.moduleExecutions.Load() == 1 }, admissionTimeout, 5*time.Millisecond,
		"the first event to run on the original engine")

	h.replaceEngine()

	h.deliver(admissionTimeout)
	require.Eventually(t, func() bool { return h.moduleExecutions.Load() == 2 }, admissionTimeout, 5*time.Millisecond,
		"the next event to run on the replacement engine")
}

// engineLifecycleHarness is a real engine of the harness's mode over a mock module whose executions
// are held by an executionBlocker until the test releases them or their context is canceled. In
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
	blocker     *executionBlocker
	log         *lifecycleLog

	initDoneCh   chan error
	subscribedCh chan []string
	finishedCh   chan string // status of each finished execution

	moduleCloses     atomic.Int32
	moduleExecutions atomic.Int32 // executions that reached the module, subscribe calls excluded
	moduleInFlight   atomic.Int32 // executions currently held in the module
	moduleMaxFlight  atomic.Int32 // the most executions ever held in the module at once

	// slots counts requests for an execution-concurrency slot, made by whichever component
	// admits events: the legacy engine or the coordinator.
	slots *slotCounter
	// sharedFlight, if set, also counts this harness's executions, so scenarios with several
	// harnesses can see how many run at once across all of them.
	sharedFlight *flightTracker
	nextEvent    atomic.Int32

	finishedMu  sync.Mutex
	finishedLog []string // status of every OnExecutionFinished call, which finishedCh drops when full
}

// defaultSignalBuffer is the capacity of the channels that carry finished and blocker signals.
// Sends never block, so a scenario that expects more signals than this must raise it.
const defaultSignalBuffer = 16

// harnessOptions holds the settings harnessOption values change.
type harnessOptions struct {
	signalBuffer  int
	eventBuffer   int
	settings      func(*cresettings.Workflows)
	limiters      *v2.EngineLimiters
	workflowID    string
	workflowOwner string
	orgID         string
	billing       *metmocks.BillingClient
	sharedFlight  *flightTracker
}

// harnessOption configures newEngineLifecycleHarness.
type harnessOption func(*harnessOptions)

// withSignalBuffer sets the capacity of the finished and blocker signal channels.
func withSignalBuffer(n int) harnessOption {
	return func(o *harnessOptions) { o.signalBuffer = n }
}

// withEventBuffer sets the capacity of the trigger capability's event channel. It is unbuffered
// by default, so a delivery only completes once the node's reader takes the event.
func withEventBuffer(n int) harnessOption {
	return func(o *harnessOptions) { o.eventBuffer = n }
}

// withSettings overrides the CRE settings the engine's limiters are built from.
func withSettings(fn func(*cresettings.Workflows)) harnessOption {
	return func(o *harnessOptions) { o.settings = fn }
}

// withLimiters makes the engine use limiters, so several harnesses can share them as the
// engines on a node do. It takes precedence over withSettings.
func withLimiters(l *v2.EngineLimiters) harnessOption {
	return func(o *harnessOptions) { o.limiters = l }
}

// withWorkflow sets the workflow ID and owner, so several harnesses can run distinct workflows.
func withWorkflow(id, owner string) harnessOption {
	return func(o *harnessOptions) { o.workflowID, o.workflowOwner = id, owner }
}

// withOrg makes the engine resolve its workflow owner to orgID.
func withOrg(orgID string) harnessOption {
	return func(o *harnessOptions) { o.orgID = orgID }
}

// withBillingClient replaces the default billing client, which accepts every call.
func withBillingClient(c *metmocks.BillingClient) harnessOption {
	return func(o *harnessOptions) { o.billing = c }
}

// withSharedFlight also counts the harness's executions in f.
func withSharedFlight(f *flightTracker) harnessOption {
	return func(o *harnessOptions) { o.sharedFlight = f }
}

// newEngineLifecycleHarness builds the engine for mode without starting it. Whatever the test
// leaves running is canceled and closed at cleanup.
func newEngineLifecycleHarness(t *testing.T, mode lifecycleEngineMode, opts ...harnessOption) *engineLifecycleHarness {
	t.Helper()

	o := harnessOptions{signalBuffer: defaultSignalBuffer}
	for _, opt := range opts {
		opt(&o)
	}

	log := &lifecycleLog{}
	h := &engineLifecycleHarness{
		t:            t,
		mode:         mode,
		engines:      &lifecycleEngineRegistry{},
		eventCh:      make(chan capabilities.TriggerResponse, o.eventBuffer),
		blocker:      newExecutionBlocker(t, log, o.signalBuffer),
		log:          log,
		initDoneCh:   make(chan error, 1),
		subscribedCh: make(chan []string, 1),
		finishedCh:   make(chan string, o.signalBuffer),
	}

	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil).Maybe()
	trigger := capmocks.NewTriggerCapability(t)
	capreg.EXPECT().GetTrigger(matches.AnyContext, lifecycleTriggerCapID).Return(trigger, nil).Maybe()
	trigger.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(h.eventCh, nil).Maybe()
	trigger.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Maybe()
	trigger.EXPECT().AckEvent(matches.AnyContext, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	cfg := defaultTestConfig(t, o.settings)
	if o.limiters != nil {
		// A copy, so wrapping its ExecutionConcurrency below does not change the shared struct.
		// The copy still holds the shared limiter instances.
		shared := *o.limiters
		cfg.LocalLimiters = &shared
	}
	if o.workflowID != "" {
		cfg.WorkflowID, cfg.WorkflowOwner = o.workflowID, o.workflowOwner
	}
	if o.orgID != "" {
		cfg.OrgResolver = &mockOrgResolver{orgID: o.orgID}
	}
	h.sharedFlight = o.sharedFlight
	h.slots = &slotCounter{ResourcePoolLimiter: cfg.LocalLimiters.ExecutionConcurrency}
	cfg.LocalLimiters.ExecutionConcurrency = h.slots
	cfg.CapRegistry = capreg
	cfg.Module = h.newModule()
	cfg.BillingClient = newLifecycleBillingClient(t)
	if o.billing != nil {
		cfg.BillingClient = o.billing
	}
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized:          func(err error) { nonBlockingSend(h.initDoneCh, err) },
		OnSubscribedToTriggers: func(ids []string) { nonBlockingSend(h.subscribedCh, ids) },
		OnExecutionFinished: func(_, status string) {
			h.log.record(eventExecutionFinished)
			h.recordFinished(status)
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
	h.engines.set(engine)

	// Close the engine before the mocks and limiters it uses go away. The test's context is
	// already canceled by now, which returns any execution still held by the blocker. An engine
	// that deadlocks here hangs until the suite timeout, which dumps every goroutine's stack.
	t.Cleanup(func() {
		_ = engine.Close()
	})
	return h
}

// newModule returns the module mock. It answers the subscribe call with one trigger and holds
// every other execution on the blocker, counting executions and logging when the module is closed.
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
			inFlight := h.moduleInFlight.Add(1)
			defer h.moduleInFlight.Add(-1)
			for maxFlight := h.moduleMaxFlight.Load(); inFlight > maxFlight && !h.moduleMaxFlight.CompareAndSwap(maxFlight, inFlight); {
				maxFlight = h.moduleMaxFlight.Load()
			}
			if h.sharedFlight != nil {
				h.sharedFlight.enter()
				defer h.sharedFlight.exit()
			}
			if err := h.blocker.block(ctx); err != nil {
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
		// The engine's own limiters, as in production, so the limits a scenario sets apply to both.
		h.cfg.LocalLimiters.Admission(),
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

// startExecution sends one trigger event through the engine's real delivery path and waits until
// its execution is running in the module. It stays there until the blocker is released or its
// context is canceled.
func (h *engineLifecycleHarness) startExecution() {
	h.t.Helper()
	h.deliver(testutils.WaitTimeout(h.t))
	awaitLifecycle(h.t, h.blocker.started, "execution to start")
}

// deliver hands the next trigger event to the node as the trigger capability does, failing the
// test if the node's reader does not take it within timeout.
func (h *engineLifecycleHarness) deliver(timeout time.Duration) {
	h.t.Helper()
	h.deliverEvent(fmt.Sprintf("event-%d", h.nextEvent.Add(1)), timeout)
}

// deliverEvent is deliver with a chosen event ID, for scenarios that send the same event twice.
func (h *engineLifecycleHarness) deliverEvent(id string, timeout time.Duration) {
	h.t.Helper()
	select {
	case h.eventCh <- h.event(id).Event:
	case <-time.After(timeout):
		require.FailNowf(h.t, "trigger event was not accepted",
			"%s was not read within %s: the node stopped draining the trigger's event channel", id, timeout)
	}
}

// feedEvents sends trigger events through the engine's real delivery path until stop is closed or
// the test ends. Unlike startExecution, it neither waits for an execution nor stops on its own, so
// delivery can overlap shutdown. The returned channel is closed once the feeder has exited.
func (h *engineLifecycleHarness) feedEvents(stop <-chan struct{}) <-chan struct{} {
	exited := make(chan struct{})
	go func() {
		defer close(exited)
		for {
			event := h.event(fmt.Sprintf("fed-%d", h.nextEvent.Add(1))).Event
			select {
			case h.eventCh <- event:
			case <-stop:
				return
			case <-h.t.Context().Done():
				return
			}
		}
	}()
	return exited
}

// close closes the engine, logging when Close returns, and returns its error.
func (h *engineLifecycleHarness) close() error {
	err := h.engine.Close()
	h.log.record(eventCloseReturned)
	return err
}

// replaceEngine closes the engine under test and starts a new one of the same mode for the same
// workflow. In coordinated mode the new engine only replaces the old one in the registry, leaving
// the coordinator registration untouched, as when a workflow's engine changes while its
// registration stays live. In legacy mode the new engine registers its own triggers.
func (h *engineLifecycleHarness) replaceEngine() {
	h.t.Helper()
	require.NoError(h.t, h.close())

	cfg := *h.cfg
	cfg.Module = h.newModule()
	construct := v2.NewCoordinatedEngine
	if h.mode == lifecycleLegacy {
		// NewEngine set the old engine as its own acknowledger; the new one must acknowledge through itself.
		cfg.TriggerAcknowledger = nil
		construct = v2.NewEngine
	}
	engine, err := construct(&cfg)
	require.NoError(h.t, err)
	h.t.Cleanup(func() { _ = engine.Close() })
	h.engine = engine
	h.engines.set(engine)

	require.NoError(h.t, engine.Start(h.t.Context()))
	require.NoError(h.t, awaitLifecycle(h.t, h.initDoneCh, "replacement engine to initialize"))
	if h.mode == lifecycleLegacy {
		require.Equal(h.t, []string{lifecycleTriggerCapID}, awaitLifecycle(h.t, h.subscribedCh, "replacement engine to register triggers"))
	}
}

// recordFinished is the OnExecutionFinished hook's record of an execution's status.
func (h *engineLifecycleHarness) recordFinished(status string) {
	h.finishedMu.Lock()
	defer h.finishedMu.Unlock()
	h.finishedLog = append(h.finishedLog, status)
}

// finishedStatuses returns the status of every execution that has finished, in finishing order.
func (h *engineLifecycleHarness) finishedStatuses() []string {
	h.finishedMu.Lock()
	defer h.finishedMu.Unlock()
	return append([]string(nil), h.finishedLog...)
}

// waitExecutionFinished returns the status of the next execution to finish.
func (h *engineLifecycleHarness) waitExecutionFinished() string {
	h.t.Helper()
	return awaitLifecycle(h.t, h.finishedCh, "execution to finish")
}

// event returns a trigger event for the workflow's one trigger. Callers pass it to ExecuteTrigger
// directly, or send its Event through the real delivery path.
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

// awaitLifecycle returns the next value from ch, or the zero value once ch is closed, failing the
// test if neither happens within testutils.WaitTimeout.
func awaitLifecycle[T any](t *testing.T, ch <-chan T, what string) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(testutils.WaitTimeout(t)):
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

// requireAllBefore fails the test unless every occurrence of before in got precedes the first
// occurrence of after, and after appears.
func requireAllBefore(t *testing.T, got []lifecycleEvent, before, after lifecycleEvent) {
	t.Helper()
	afterAt := -1
	for i, e := range got {
		if e == after {
			afterAt = i
			break
		}
	}
	require.NotEqual(t, -1, afterAt, "events %v do not contain %v", got, after)
	for i, e := range got {
		if e == before {
			require.Less(t, i, afterAt, "%v at %d follows %v at %d in %v", before, i, after, afterAt, got)
		}
	}
}

// lifecycleLog records the events a scenario orders its assertions on.
type lifecycleLog struct {
	mu     sync.Mutex
	logged []lifecycleEvent
}

// record appends e to the log.
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

// lifecycleEngineRegistry is a triggers.EngineRegistry holding the one engine under test. The
// engine can be swapped while the coordinator reads it.
type lifecycleEngineRegistry struct {
	mu     sync.Mutex
	engine v2.WorkflowEngine
}

// Get returns the engine under test if wid is its workflow.
func (r *lifecycleEngineRegistry) Get(wid types.WorkflowID) (triggers.RegisteredEngine, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.engine == nil || wid.Hex() != r.engine.Tenant().Workflow {
		return nil, false
	}
	return r.engine, true
}

// set makes e the engine under test.
func (r *lifecycleEngineRegistry) set(e v2.WorkflowEngine) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.engine = e
}

// executionBlocker holds executions in the module until the test releases them or their context
// is canceled. An execution nothing releases or cancels fails the test after backstop.
type executionBlocker struct {
	t        *testing.T
	backstop time.Duration
	started  chan struct{}
	canceled chan struct{} // signaled when a held execution observes its context's cancellation
	released chan struct{}
	step     chan struct{} // releases one held execution per receive
	once     sync.Once
	log      *lifecycleLog

	mu           sync.Mutex
	errAtRelease error // the context's error when a released execution resumed
}

// newExecutionBlocker returns a blocker whose started and canceled signals hold up to signalBuffer
// unread values. Its backstop is half the test's wait budget, so a stuck execution fails with its
// own message before the suite timeout.
func newExecutionBlocker(t *testing.T, log *lifecycleLog, signalBuffer int) *executionBlocker {
	return &executionBlocker{
		t:        t,
		backstop: testutils.WaitTimeout(t) / 2,
		started:  make(chan struct{}, signalBuffer),
		canceled: make(chan struct{}, signalBuffer),
		released: make(chan struct{}),
		step:     make(chan struct{}),
		log:      log,
	}
}

// block signals that an execution started and holds it. It returns nil once the blocker is
// released, and the context's error, after logging the cancellation, if the context is canceled
// first. It fails the test if neither happens within the backstop.
func (b *executionBlocker) block(ctx context.Context) error {
	// signal that a blocking execution has started
	nonBlockingSend(b.started, struct{}{})

	backstop := time.NewTimer(b.backstop)
	defer backstop.Stop()

	select {
	case <-b.step:
		return nil
	case <-b.released:
		b.mu.Lock()
		b.errAtRelease = ctx.Err()
		b.mu.Unlock()
		return nil
	case <-ctx.Done():
		b.log.record(eventExecutionCanceled)
		nonBlockingSend(b.canceled, struct{}{})
		return ctx.Err()
	case <-backstop.C:
		b.t.Errorf("execution was held for %s without release or cancellation", b.backstop)
		return errors.New("execution blocker backstop elapsed")
	case <-b.t.Context().Done():
		return b.t.Context().Err()
	}
}

// releaseOne lets exactly one held execution return, waiting until one takes the release.
func (b *executionBlocker) releaseOne() {
	b.t.Helper()
	select {
	case b.step <- struct{}{}:
	case <-time.After(b.backstop):
		require.FailNow(b.t, "no held execution to release")
	}
}

// release lets every held execution return.
func (b *executionBlocker) release() {
	b.once.Do(func() { close(b.released) })
}

// contextErrAtRelease returns the context's error as seen by an execution that resumed when the
// blocker was released.
func (b *executionBlocker) contextErrAtRelease() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.errAtRelease
}

// slotCounter counts Wait calls on the execution-concurrency limiter, so a scenario can tell when
// an event is waiting for a slot.
type slotCounter struct {
	limits.ResourcePoolLimiter[int]
	waits    atomic.Int32 // Wait calls
	acquired atomic.Int32 // slots granted
	released atomic.Int32 // slots given back
}

func (s *slotCounter) Wait(ctx context.Context, amount int) (func(), error) {
	s.waits.Add(1)
	free, err := s.ResourcePoolLimiter.Wait(ctx, amount)
	if err != nil {
		return free, err
	}
	s.acquired.Add(1)
	return func() {
		free()
		s.released.Add(1)
	}, nil
}

// flightTracker counts executions running at once across several harnesses, and the most seen.
type flightTracker struct {
	inFlight atomic.Int32
	max      atomic.Int32
}

func (f *flightTracker) enter() {
	n := f.inFlight.Add(1)
	for m := f.max.Load(); n > m && !f.max.CompareAndSwap(m, n); {
		m = f.max.Load()
	}
}

func (f *flightTracker) exit() { f.inFlight.Add(-1) }

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
	return newTestBillingClient(t, true)
}

// newTestBillingClient is newLifecycleBillingClient whose credit reservations succeed only if
// reserveOK. A failed reservation makes the engine abort the execution before the module runs.
func newTestBillingClient(t *testing.T, reserveOK bool) *metmocks.BillingClient {
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
		Return(&billing.ReserveCreditsResponse{Success: reserveOK, Credits: "10000"}, nil).Maybe()
	client.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).Return(&emptypb.Empty{}, nil).Maybe()
	return client
}
