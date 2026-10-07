package v2

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/workflowkey"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capreg "github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonmetrics "github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/dontime"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/confidentialrelay"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	metmocks "github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	workflowstore "github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
)

// engineMode is the kind of engine a lifecycle test runs against.
type engineMode string

const (
	// modeLegacy is the engine that registers its own triggers and owns the workflow-limit slot.
	modeLegacy engineMode = "legacy"
	// modeCoordinated is the execution-only engine whose triggers and workflow-limit slot are
	// owned by a real trigger coordinator.
	modeCoordinated engineMode = "coordinated"
)

// Test_workflowLifecycle drives the real event handler, engine registry and workflow limiter
// through each workflow lifecycle path, once with the legacy engine and once with the
// coordinated engine behind a real trigger coordinator. Only the capability registry, trigger,
// WASM module and billing client are mocked. Each subtest reads top to bottom: send lifecycle
// events, then check the three outcomes that matter: the workflow-limit slots in use, whether
// the workflow's engine is registered, and how many times its engine was closed.
func Test_workflowLifecycle(t *testing.T) {
	t.Parallel()

	for _, mode := range []engineMode{modeLegacy, modeCoordinated} {
		t.Run(string(mode), func(t *testing.T) {
			t.Parallel()

			run := func(name string, test func(*testing.T, engineMode)) {
				t.Run(name, func(t *testing.T) { test(t, mode) })
			}
			run("register holds a slot and rejects a second workflow", testRegisterHoldsSlotAndRejectsSecondWorkflow)
			run("pause frees the slot and redelivery is a no-op", testPauseFreesSlotAndRedeliveryIsNoOp)
			run("activate after pause holds one slot", testActivateAfterPauseHoldsOneSlot)
			run("delete frees the slot and admits the next workflow", testDeleteFreesSlotAndAdmitsNextWorkflow)
			run("delete with an execution in flight is deferred", testDeleteWithExecutionInFlightIsDeferred)
			run("pause with an execution in flight is deferred", testPauseWithExecutionInFlightIsDeferred)
			run("registered event with paused status releases the slot", testPausedRegisteredEventReleasesSlot)
			run("replacing a draining engine holds one slot", testReplacingDrainingEngineHoldsOneSlot)
			run("failed registration leaks nothing", testFailedRegistrationLeaksNothing)
			run("activate while an execution is in flight replaces the engine", testActivateWhileExecutionInFlightReplacesEngine)
			run("delete then register at the limit admits the next workflow", testDeleteThenRegisterAtLimitAdmitsNextWorkflow)
		})
	}
}

// testRegisterHoldsSlotAndRejectsSecondWorkflow registers a workflow on a node that allows one
// workflow. It asserts the workflow holds the only slot and has an engine, and that a second
// workflow is refused with ErrGlobalWorkflowCountLimitReached, leaving the slot count
// unchanged, no engine registered, and its engine closed once.
func testRegisterHoldsSlotAndRejectsSecondWorkflow(t *testing.T, mode engineMode) {
	t.Parallel()
	h := newLifecycleHarness(t, mode, 1)
	wfA, wfB := h.workflow("wf-a"), h.workflow("wf-b")

	require.NoError(t, h.registered(wfA))
	h.requireSlotsInUse(1)
	h.requireEngineRegistered(wfA)

	// The node is at its workflow limit, so the second workflow is refused.
	err := h.registered(wfB)
	require.ErrorIs(t, err, types.ErrGlobalWorkflowCountLimitReached)
	h.requireSlotsInUse(1)
	h.requireNoEngine(wfB)
	h.requireEngineCloses(wfB, 1)
}

// testPauseFreesSlotAndRedeliveryIsNoOp pauses a running workflow. It asserts the slot is
// freed, the engine is removed and closed once, and that delivering the same pause event again
// changes nothing: no second close and no double free.
func testPauseFreesSlotAndRedeliveryIsNoOp(t *testing.T, mode engineMode) {
	t.Parallel()
	h := newLifecycleHarness(t, mode, 1)
	wf := h.workflow("wf-a")

	require.NoError(t, h.registered(wf))
	h.requireSlotsInUse(1)

	require.NoError(t, h.paused(wf))
	h.requireSlotsInUse(0)
	h.requireNoEngine(wf)
	h.requireEngineCloses(wf, 1)

	// Redelivery of pause finds no engine
	require.NoError(t, h.paused(wf))
	h.requireSlotsInUse(0)
	h.requireEngineCloses(wf, 1)
}

// testActivateAfterPauseHoldsOneSlot pauses a workflow and activates it again. It asserts the
// pause frees the slot and the activation takes exactly one slot back with a new registered
// engine, so the old engine was closed only once.
func testActivateAfterPauseHoldsOneSlot(t *testing.T, mode engineMode) {
	t.Parallel()
	h := newLifecycleHarness(t, mode, 1)
	wf := h.workflow("wf-a")

	require.NoError(t, h.registered(wf))
	require.NoError(t, h.paused(wf))
	h.requireSlotsInUse(0)

	require.NoError(t, h.activated(wf))
	h.requireSlotsInUse(1)
	h.requireEngineRegistered(wf)
	h.requireEngineCloses(wf, 1)
}

// testDeleteFreesSlotAndAdmitsNextWorkflow deletes a running workflow on a node that allows
// one workflow. It asserts the slot is freed and the engine removed and closed once, that
// redelivering the delete is a no-op, and that a second workflow can then take the freed slot.
func testDeleteFreesSlotAndAdmitsNextWorkflow(t *testing.T, mode engineMode) {
	t.Parallel()
	h := newLifecycleHarness(t, mode, 1)
	wfA, wfB := h.workflow("wf-a"), h.workflow("wf-b")

	require.NoError(t, h.registered(wfA))
	require.NoError(t, h.deleted(wfA))
	h.requireSlotsInUse(0)
	h.requireNoEngine(wfA)
	h.requireEngineCloses(wfA, 1)

	// Redelivery of delete wfA finds no engine
	require.NoError(t, h.deleted(wfA))
	h.requireSlotsInUse(0)
	h.requireEngineCloses(wfA, 1)

	// The freed slot is available to another workflow.
	require.NoError(t, h.registered(wfB))
	h.requireSlotsInUse(1)
	h.requireEngineRegistered(wfB)
}

// testDeleteWithExecutionInFlightIsDeferred deletes a workflow while one of its executions is
// running. It asserts the delete returns ErrDrainInProgress and leaves the slot held and the
// engine registered and open. Once the execution finishes, it asserts retrying the delete
// frees the slot and removes and closes the engine.
func testDeleteWithExecutionInFlightIsDeferred(t *testing.T, mode engineMode) {
	t.Parallel()
	h := newLifecycleHarness(t, mode, 1)
	wf := h.workflow("wf-a")

	require.NoError(t, h.registered(wf))
	h.startExecution(wf)

	// The delete drains the engine but cannot close it while the execution runs.
	err := h.deleted(wf)
	require.ErrorIs(t, err, ErrDrainInProgress)
	h.requireSlotsInUse(1)
	h.requireEngineRegistered(wf)
	h.requireEngineCloses(wf, 0)

	h.finishExecution(wf)

	// The retry finds no active executions and completes the delete.
	require.NoError(t, h.deleted(wf))
	h.requireSlotsInUse(0)
	h.requireNoEngine(wf)
	h.requireEngineCloses(wf, 1)
}

// testPauseWithExecutionInFlightIsDeferred pauses a workflow while one of its executions is
// running. It asserts the pause returns ErrDrainInProgress and leaves the slot held and the
// engine registered and open. Once the execution finishes, it asserts retrying the pause frees
// the slot and removes and closes the engine.
func testPauseWithExecutionInFlightIsDeferred(t *testing.T, mode engineMode) {
	t.Parallel()
	h := newLifecycleHarness(t, mode, 1)
	wf := h.workflow("wf-a")

	require.NoError(t, h.registered(wf))
	h.startExecution(wf)

	// The pause drains the engine but cannot close it while the execution runs.
	err := h.paused(wf)
	require.ErrorIs(t, err, ErrDrainInProgress)
	h.requireSlotsInUse(1)
	h.requireEngineRegistered(wf)
	h.requireEngineCloses(wf, 0)

	h.finishExecution(wf)

	// The retry finds no active executions and completes the pause.
	require.NoError(t, h.paused(wf))
	h.requireSlotsInUse(0)
	h.requireNoEngine(wf)
	h.requireEngineCloses(wf, 1)
}

// testPausedRegisteredEventReleasesSlot sends a registered event that carries the paused
// status for a running workflow. It asserts the engine is stopped, the slot is freed, and the
// engine is removed and closed once, as it would be for a pause event.
func testPausedRegisteredEventReleasesSlot(t *testing.T, mode engineMode) {
	t.Parallel()
	h := newLifecycleHarness(t, mode, 1)
	wf := h.workflow("wf-a")

	require.NoError(t, h.registered(wf))
	h.requireSlotsInUse(1)

	// A registered event that carries the paused status stops the running engine.
	require.NoError(t, h.registeredAsPaused(wf))
	h.requireSlotsInUse(0)
	h.requireNoEngine(wf)
	h.requireEngineCloses(wf, 1)
}

// testReplacingDrainingEngineHoldsOneSlot pauses a workflow while an execution is running, so
// the pause is deferred and leaves the engine marked as draining. The workflow is then activated
// again before the pause is retried, on a node that allows one workflow. It asserts the handler
// replaces the draining engine instead of treating it as healthy: the old engine is closed and
// its slot freed before the new engine takes one, so a single slot suffices and exactly one
// engine is registered afterwards.
func testReplacingDrainingEngineHoldsOneSlot(t *testing.T, mode engineMode) {
	t.Parallel()
	h := newLifecycleHarness(t, mode, 1)
	wf := h.workflow("wf-a")

	require.NoError(t, h.registered(wf))
	h.startExecution(wf)

	// The pause is deferred by the running execution. The engine is now marked as draining and
	// stays registered and open.
	require.ErrorIs(t, h.paused(wf), ErrDrainInProgress)
	h.requireEngineRegistered(wf)
	h.requireEngineCloses(wf, 0)

	// Finish execution to unblock the remainder of the test.
	// Finishing the execution does not undo the drain: the engine stays marked as draining
	// until it is closed.
	h.finishExecution(wf)

	// Activating a workflow whose engine is draining makes the handler replace that engine: it
	// closes the old one, which frees its slot, then creates a new one that takes a slot.
	require.NoError(t, h.activated(wf))
	h.requireSlotsInUse(1)
	h.requireEngineRegistered(wf)
	h.requireEngineCloses(wf, 1)
}

// testFailedRegistrationLeaksNothing registers a workflow whose trigger registration fails. It
// asserts the registration returns an error, no slot stays held, no engine is registered, and
// the engine that was created is closed once.
func testFailedRegistrationLeaksNothing(t *testing.T, mode engineMode) {
	t.Parallel()
	h := newLifecycleHarness(t, mode, 1)
	wf := h.workflow("wf-a")
	h.failTriggerRegistrations(errors.New("trigger registration failed"))

	require.Error(t, h.registered(wf))
	h.requireSlotsInUse(0)
	h.requireNoEngine(wf)
	h.requireEngineCloses(wf, 1)
}

// testActivateWhileExecutionInFlightReplacesEngine pauses a workflow while an execution is
// running, so the pause is deferred and leaves the engine draining, then activates the workflow
// again without waiting for the execution to finish. It asserts the draining engine is replaced
// and the slot handed over, leaving one engine registered and one slot held with the old engine
// closed once, and that the running execution was canceled rather than waited for.
func testActivateWhileExecutionInFlightReplacesEngine(t *testing.T, mode engineMode) {
	t.Parallel()
	h := newLifecycleHarness(t, mode, 1)
	wf := h.workflow("wf-a")

	require.NoError(t, h.registered(wf))
	h.startExecution(wf)
	require.ErrorIs(t, h.paused(wf), ErrDrainInProgress)

	// The activation arrives while the execution is still running. Replacing the draining
	// engine closes it, which cancels the execution instead of waiting for it to finish.
	require.NoError(t, h.activated(wf))
	h.requireSlotsInUse(1)
	h.requireEngineRegistered(wf)
	h.requireEngineCloses(wf, 1)
	h.requireExecutionCanceled(wf)
}

// testDeleteThenRegisterAtLimitAdmitsNextWorkflow deletes a workflow on a node that allows one
// workflow and registers another immediately afterwards, as a reconcile batch with both events
// would. It asserts the second registration is admitted without waiting for the freed slot to
// become available, and ends with only the second workflow holding the slot.
func testDeleteThenRegisterAtLimitAdmitsNextWorkflow(t *testing.T, mode engineMode) {
	t.Parallel()
	h := newLifecycleHarness(t, mode, 1)
	wfA, wfB := h.workflow("wf-a"), h.workflow("wf-b")

	require.NoError(t, h.registered(wfA))
	require.NoError(t, h.deleted(wfA))
	require.NoError(t, h.registered(wfB))

	h.requireSlotsInUse(1)
	h.requireNoEngine(wfA)
	h.requireEngineRegistered(wfB)
}

// lifecycleTriggerID is the one trigger every lifecycle workflow subscribes to.
const lifecycleTriggerID = "id_0"

// lifecycleWorkflow identifies one workflow a lifecycle test drives.
type lifecycleWorkflow struct {
	id    types.WorkflowID
	owner []byte
	name  string
}

// registered returns a registered event for w carrying the given status.
func (w lifecycleWorkflow) registered(status uint8) WorkflowRegisteredEvent {
	return WorkflowRegisteredEvent{
		Status:        status,
		WorkflowID:    w.id,
		WorkflowOwner: w.owner,
		WorkflowName:  w.name,
		CreatedAt:     1,
	}
}

// lifecycleHarness is an event handler wired to a real engine registry and workflow
// limiter. The engine factory builds a real engine of the harness's mode, whose executions
// block until the test lets them finish. In coordinated mode a real trigger coordinator owns
// trigger registration, event delivery and the workflow-limit slot.
type lifecycleHarness struct {
	t           *testing.T
	mode        engineMode
	coordinator triggers.Coordinator // nil in legacy mode
	handler     *eventHandler
	registry    *EngineRegistry
	limits      *countingLimiter

	capRegistry  *regmocks.CapabilitiesRegistry
	limiters     *v2.EngineLimiters
	featureFlags *v2.EngineFeatureFlags

	mu         sync.Mutex
	engines    map[types.WorkflowID][]*trackedEngine
	eventChs   map[string]chan capabilities.TriggerResponse // workflow ID -> latest trigger event channel
	gates      map[string]*executionGate                    // workflow ID -> execution gate
	triggerErr error                                        // returned from RegisterTrigger when set
	nextEvent  int
}

// newLifecycleHarness returns a harness for the given engine mode whose node allows globalLimit
// workflows.
func newLifecycleHarness(t *testing.T, mode engineMode, globalLimit int) *lifecycleHarness {
	t.Helper()
	lggr := logger.TestLogger(t)
	lf := limits.Factory{Logger: lggr}

	workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{Global: int32(globalLimit), PerOwner: 200}, lf) //nolint:gosec // G115: small test value
	require.NoError(t, err)
	limiters, err := v2.NewLimiters(lf, nil)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, limiters.Close()) })
	featureFlags, err := v2.NewFeatureFlags(lf, nil)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, featureFlags.Close()) })

	h := &lifecycleHarness{
		t:            t,
		mode:         mode,
		registry:     NewEngineRegistry(),
		limits:       &countingLimiter{ResourceLimiter: workflowLimits},
		capRegistry:  regmocks.NewCapabilitiesRegistry(t),
		limiters:     limiters,
		featureFlags: featureFlags,
		engines:      make(map[types.WorkflowID][]*trackedEngine),
		eventChs:     make(map[string]chan capabilities.TriggerResponse),
		gates:        make(map[string]*executionGate),
	}
	h.expectCapabilities()
	// Cleanups run last-in first-out, so this closes any engine the test left
	// running before the limiters those engines use are closed.
	t.Cleanup(h.closeEngines)

	handlerOpts := []func(*eventHandler){WithEngineFactoryFn(h.engineFactory)}
	if mode == modeCoordinated {
		h.coordinator = h.newCoordinator()
		handlerOpts = append(handlerOpts, WithTriggerCoordinator(h.coordinator))
	}

	rl, err := ratelimiter.NewRateLimiter(rlConfig)
	require.NoError(t, err)
	h.handler, err = NewEventHandler(
		lggr,
		workflowstore.NewInMemoryStore(lggr, clockwork.NewFakeClock()),
		nil,
		true,
		h.capRegistry,
		&confidentialrelay.ExecutionHandlers{},
		h.registry,
		custmsg.NewLabeler(),
		limiters,
		nil,
		rl,
		h.limits,
		&lifecycleStore{byID: make(map[string]*job.WorkflowSpec)},
		workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
		&testDonNotifier{},
		handlerOpts...,
	)
	require.NoError(t, err)
	return h
}

// newCoordinator starts a real trigger coordinator over the harness's engine registry,
// capability registry and workflow limiter.
func (h *lifecycleHarness) newCoordinator() triggers.Coordinator {
	h.t.Helper()
	em, err := monitoring.InitMonitoringResources()
	require.NoError(h.t, err)
	coordinator := triggers.NewCoordinator(
		triggers.RegisterDeps{
			CapRegistry:  h.capRegistry,
			RegTimeout:   limits.NewTimeLimiter(5 * time.Second),
			ChainAllowed: limits.NewGateLimiter(true),
			Logger:       commonlogger.Test(h.t),
			Metrics:      monitoring.NewWorkflowsMetricLabeler(commonmetrics.NewLabeler(), em),
		},
		NewTriggerEngineRegistry(h.registry),
		h.limits,
		clockwork.NewFakeClock(),
	)
	servicetest.Run(h.t, coordinator)
	return coordinator
}

// workflow returns a workflow whose ID matches the artifacts the stub store
// serves (binary "binary", config "config"), so it passes validation. The owner
// is a full 20-byte address, which the engine requires.
func (h *lifecycleHarness) workflow(name string) lifecycleWorkflow {
	h.t.Helper()
	owner := bytes.Repeat([]byte{0xaa}, 20)
	id, err := pkgworkflows.GenerateWorkflowID(owner, name, []byte("binary"), []byte("config"), "")
	require.NoError(h.t, err)
	return lifecycleWorkflow{id: types.WorkflowID(id), owner: owner, name: name}
}

// Lifecycle events, as the syncer delivers them to the handler.

// registered delivers a registered event with the active status for w, which creates and
// starts its engine, and returns the handler's error.
func (h *lifecycleHarness) registered(w lifecycleWorkflow) error {
	return h.handler.workflowRegisteredEvent(h.t.Context(), w.registered(WorkflowStatusActive))
}

// registeredAsPaused delivers a registered event with the paused status for w, which stops
// its engine if one is running, and returns the handler's error.
func (h *lifecycleHarness) registeredAsPaused(w lifecycleWorkflow) error {
	return h.handler.workflowRegisteredEvent(h.t.Context(), w.registered(WorkflowStatusPaused))
}

// activated delivers an activated event for w, which starts a new engine for a paused
// workflow, and returns the handler's error.
func (h *lifecycleHarness) activated(w lifecycleWorkflow) error {
	return h.handler.workflowActivatedEvent(h.t.Context(), WorkflowActivatedEvent(w.registered(WorkflowStatusActive)))
}

// paused delivers a paused event for w, which drains and closes its engine, and returns the
// handler's error. It is ErrDrainInProgress while executions are still running.
func (h *lifecycleHarness) paused(w lifecycleWorkflow) error {
	return h.handler.workflowPausedEvent(h.t.Context(), WorkflowPausedEvent{WorkflowID: w.id})
}

// deleted delivers a deleted event for w, which drains and closes its engine and removes its
// spec, and returns the handler's error. It is ErrDrainInProgress while executions are still
// running.
func (h *lifecycleHarness) deleted(w lifecycleWorkflow) error {
	return h.handler.workflowDeletedEvent(h.t.Context(), WorkflowDeletedEvent{WorkflowID: w.id}, hex.EncodeToString(w.owner))
}

// startExecution delivers a trigger event to w and waits until its execution is
// running. The execution stays in flight until finishExecution.
func (h *lifecycleHarness) startExecution(w lifecycleWorkflow) {
	h.t.Helper()
	workflowID := w.id.Hex()
	h.mu.Lock()
	ch := h.eventChs[workflowID]
	h.nextEvent++
	eventID := fmt.Sprintf("event-%d", h.nextEvent)
	h.mu.Unlock()
	require.NotNil(h.t, ch, "workflow %s has no registered trigger", w.name)

	event := capabilities.TriggerResponse{Event: capabilities.TriggerEvent{TriggerType: lifecycleTriggerID, ID: eventID}}
	select {
	case ch <- event:
	case <-time.After(tests.WaitTimeout(h.t)):
		require.FailNow(h.t, "trigger event was not accepted")
	}
	select {
	case <-h.gate(workflowID).started:
	case <-time.After(tests.WaitTimeout(h.t)):
		require.FailNow(h.t, "execution did not start")
	}
}

// finishExecution lets w's in-flight executions return and waits until its
// engine reports none active.
func (h *lifecycleHarness) finishExecution(w lifecycleWorkflow) {
	h.t.Helper()
	h.gate(w.id.Hex()).open()
	engine := h.latestEngine(w)
	require.NotNil(h.t, engine, "workflow %s has no engine", w.name)
	require.Eventually(h.t, func() bool { return engine.ActiveExecutions() == 0 },
		tests.WaitTimeout(h.t), 10*time.Millisecond, "executions did not finish")
}

// failTriggerRegistrations makes every subsequent trigger registration fail.
func (h *lifecycleHarness) failTriggerRegistrations(err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.triggerErr = err
}

// Outcome checks.

// requireSlotsInUse fails the test unless exactly want workflow-limit slots are in use,
// summed over all workflows.
func (h *lifecycleHarness) requireSlotsInUse(want int) {
	h.t.Helper()
	if h.mode == modeCoordinated {
		// The coordinator frees a slot on its own goroutine once the workflow's readers drain.
		require.Eventually(h.t, func() bool { return h.limits.inUse() == want },
			tests.WaitTimeout(h.t), 10*time.Millisecond, "workflow limit slots in use, want %d", want)
		return
	}
	require.Equal(h.t, want, h.limits.inUse(), "workflow limit slots in use")
}

// requireEngineRegistered fails the test unless w has an engine in the engine registry.
func (h *lifecycleHarness) requireEngineRegistered(w lifecycleWorkflow) {
	h.t.Helper()
	_, ok := h.registry.Get(w.id)
	require.True(h.t, ok, "workflow %s should have an engine in the registry", w.name)
}

// requireNoEngine fails the test if w has an engine in the engine registry.
func (h *lifecycleHarness) requireNoEngine(w lifecycleWorkflow) {
	h.t.Helper()
	_, ok := h.registry.Get(w.id)
	require.False(h.t, ok, "workflow %s should have no engine in the registry", w.name)
}

// requireExecutionCanceled fails the test unless w's in-flight execution ended because its
// context was canceled, as happens when the engine running it is closed.
func (h *lifecycleHarness) requireExecutionCanceled(w lifecycleWorkflow) {
	h.t.Helper()
	require.True(h.t, h.gate(w.id.Hex()).canceled.Load(), "workflow %s's execution should have been canceled", w.name)
}

// requireEngineCloses checks the Close calls summed over every engine created for w.
func (h *lifecycleHarness) requireEngineCloses(w lifecycleWorkflow, want int32) {
	h.t.Helper()
	h.mu.Lock()
	defer h.mu.Unlock()
	var got int32
	for _, e := range h.engines[w.id] {
		got += e.closes.Load()
	}
	require.Equal(h.t, want, got, "engine Close calls for workflow %s", w.name)
}

// Mocks and engine construction.

// expectCapabilities mocks the capabilities registry and the one trigger every
// workflow subscribes to. RegisterTrigger hands back a fresh event channel per
// registration, kept so tests can deliver events to the workflow.
func (h *lifecycleHarness) expectCapabilities() {
	node, err := (&capreg.TestRegistryMetadata{}).LocalNode(h.t.Context())
	require.NoError(h.t, err)
	h.capRegistry.EXPECT().LocalNode(mock.Anything).Return(node, nil).Maybe()

	trigger := capmocks.NewTriggerCapability(h.t)
	h.capRegistry.EXPECT().GetTrigger(mock.Anything, lifecycleTriggerID).Return(trigger, nil).Maybe()
	trigger.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
			h.mu.Lock()
			defer h.mu.Unlock()
			if h.triggerErr != nil {
				return nil, h.triggerErr
			}
			ch := make(chan capabilities.TriggerResponse)
			h.eventChs[req.Metadata.WorkflowID] = ch
			return ch, nil
		}).Maybe()
	trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Maybe()
	trigger.EXPECT().AckEvent(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
}

// gate returns the execution gate for the workflow with the given hex ID, creating it on
// first use. Every gate is opened at test cleanup so no execution is left blocked.
func (h *lifecycleHarness) gate(workflowID string) *executionGate {
	h.mu.Lock()
	defer h.mu.Unlock()
	g, ok := h.gates[workflowID]
	if !ok {
		g = newExecutionGate()
		h.gates[workflowID] = g
		h.t.Cleanup(g.open)
	}
	return g
}

// latestEngine returns the most recent engine the factory created for w, or nil if none.
func (h *lifecycleHarness) latestEngine(w lifecycleWorkflow) *trackedEngine {
	h.mu.Lock()
	defer h.mu.Unlock()
	engines := h.engines[w.id]
	if len(engines) == 0 {
		return nil
	}
	return engines[len(engines)-1]
}

// closeEngines closes every engine the harness created. Closing an engine that
// is already closed is harmless.
func (h *lifecycleHarness) closeEngines() {
	h.mu.Lock()
	var all []*trackedEngine
	for _, engines := range h.engines {
		all = append(all, engines...)
	}
	h.mu.Unlock()
	for _, e := range all {
		_ = e.Close()
	}
}

// engineFactory builds the real engine for a workflow in the harness's mode: the legacy engine,
// or the coordinated engine acknowledging through the harness's coordinator. Its module mock
// returns one trigger subscription and holds every execution on the workflow's gate.
func (h *lifecycleHarness) engineFactory(
	_ context.Context, wfid, owner string, name types.WorkflowName, tag string, config, _ []byte, _ string, _ []byte,
	initDone chan<- error,
) (v2.WorkflowEngine, error) {
	t := h.t
	gate := h.gate(wfid)

	module := modulemocks.NewModuleV2(t)
	module.EXPECT().Start().Maybe()
	module.EXPECT().Close().Maybe()
	module.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, req *sdk.ExecuteRequest, _ host.ExecutionHelper) (*sdk.ExecutionResult, error) {
			if _, ok := req.Request.(*sdk.ExecuteRequest_Subscribe); ok {
				return &sdk.ExecutionResult{
					Result: &sdk.ExecutionResult_TriggerSubscriptions{
						TriggerSubscriptions: &sdk.TriggerSubscriptionRequest{
							Subscriptions: []*sdk.TriggerSubscription{{Id: lifecycleTriggerID, Method: "method"}},
						},
					},
				}, nil
			}
			gate.block(ctx)
			return nil, nil
		}).Maybe()

	donSubscriber := capmocks.NewDonSubscriber(t)
	donSubscriber.EXPECT().Subscribe(mock.Anything).Return(make(<-chan capabilities.DON), func() {}, nil).Maybe()

	lggr := logger.TestLogger(t)
	cfg := &v2.EngineConfig{
		Lggr:                          lggr,
		Module:                        module,
		CapRegistry:                   h.capRegistry,
		DonTimeStore:                  dontime.NewStore(dontime.DefaultRequestTimeout),
		UseLocalTimeProvider:          true,
		DonSubscriber:                 donSubscriber,
		ExecutionsStore:               workflowstore.NewInMemoryStore(lggr, clockwork.NewRealClock()),
		WorkflowID:                    wfid,
		WorkflowOwner:                 owner,
		WorkflowName:                  name,
		WorkflowTag:                   tag,
		WorkflowConfig:                config,
		WorkflowEncryptionKey:         workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
		LocalLimits:                   v2.EngineLimits{},
		LocalLimiters:                 h.limiters,
		FeatureFlags:                  h.featureFlags,
		GlobalWorkflowLimit:           h.limits,
		BeholderEmitter:               custmsg.NewLabeler(),
		BillingClient:                 newLifecycleBillingClient(t),
		WorkflowRegistryAddress:       "0x123",
		WorkflowRegistryChainSelector: "11155111",
		Hooks: v2.LifecycleHooks{
			OnInitialized: func(err error) {
				select {
				case initDone <- err:
				default:
				}
			},
		},
	}

	construct := v2.NewEngine
	if h.mode == modeCoordinated {
		// The coordinated engine registers no triggers and holds no limit slot of its own, so
		// it acknowledges events through the coordinator that owns them.
		cfg.TriggerAcknowledger = h.coordinator
		construct = v2.NewCoordinatedEngine
	}
	inner, err := construct(cfg)
	if err != nil {
		return nil, err
	}

	engine := &trackedEngine{WorkflowEngine: inner}
	wid, err := types.WorkflowIDFromHex(wfid)
	if err != nil {
		return nil, err
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.engines[wid] = append(h.engines[wid], engine)
	return engine, nil
}

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

// lifecycleExecutionTimeout bounds how long a blocked execution waits, standing in for the real
// per-execution timeout so a test whose engine never cancels the execution fails instead of hanging.
const lifecycleExecutionTimeout = 5 * time.Second

// executionGate holds a workflow's executions in flight until opened, their context is canceled,
// or lifecycleExecutionTimeout passes.
type executionGate struct {
	started  chan struct{}
	release  chan struct{}
	once     sync.Once
	canceled atomic.Bool // an execution ended because its context was canceled
}

func newExecutionGate() *executionGate {
	return &executionGate{started: make(chan struct{}, 16), release: make(chan struct{})}
}

// block signals that an execution started and holds it until the gate opens, ctx is canceled
// (as it is when the engine running the execution is closed), or the timeout passes.
func (g *executionGate) block(ctx context.Context) {
	select {
	case g.started <- struct{}{}:
	default:
	}

	timeout := time.NewTimer(lifecycleExecutionTimeout)
	defer timeout.Stop()
	select {
	case <-g.release:
	case <-ctx.Done():
		g.canceled.Store(true)
	case <-timeout.C:
	}
}

func (g *executionGate) open() { g.once.Do(func() { close(g.release) }) }

// trackedEngine counts Close calls on the engine it wraps.
type trackedEngine struct {
	v2.WorkflowEngine
	closes atomic.Int32
}

func (e *trackedEngine) Close() error {
	e.closes.Add(1)
	return e.WorkflowEngine.Close()
}

// countingLimiter tracks the slots in use on top of a real workflow limiter.
type countingLimiter struct {
	limits.ResourceLimiter[int]
	mu   sync.Mutex
	used int
}

func (l *countingLimiter) Use(ctx context.Context, amount int) error {
	if err := l.ResourceLimiter.Use(ctx, amount); err != nil {
		return err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.used += amount
	return nil
}

func (l *countingLimiter) Free(ctx context.Context, amount int) error {
	if err := l.ResourceLimiter.Free(ctx, amount); err != nil {
		return err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.used -= amount
	return nil
}

func (l *countingLimiter) inUse() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.used
}

// lifecycleStore is an in-memory artifacts store keyed by workflow ID, so more
// than one workflow can live in it.
type lifecycleStore struct {
	stubWorkflowArtifactsStore
	mu   sync.Mutex
	byID map[string]*job.WorkflowSpec
}

func (s *lifecycleStore) GetWorkflowSpec(_ context.Context, id string) (*job.WorkflowSpec, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	spec, ok := s.byID[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	cp := *spec
	return &cp, nil
}

func (s *lifecycleStore) UpsertWorkflowSpec(_ context.Context, spec *job.WorkflowSpec) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := *spec
	s.byID[spec.WorkflowID] = &cp
	return 1, nil
}

func (s *lifecycleStore) DeleteWorkflowArtifacts(_ context.Context, id string) (*job.WorkflowSpec, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	spec, ok := s.byID[id]
	if !ok {
		return nil, nil
	}
	delete(s.byID, id)
	return spec, nil
}

func (s *lifecycleStore) PauseWorkflowArtifacts(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if spec, ok := s.byID[id]; ok {
		spec.Status = job.WorkflowSpecStatusPaused
		spec.Workflow = ""
		spec.Config = ""
	}
	return nil
}
