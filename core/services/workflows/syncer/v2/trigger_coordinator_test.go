package v2

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	commonmetrics "github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	enginev2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

// fakePutEngine is a services.Service that also implements the coordinator's
// narrow eventSink and activeExecutionsReporter interfaces, so a
// RegisterTriggers reader can resolve and deliver to it via the
// EngineRegistry, and UnregisterTriggers can wait on it, exactly as it would
// a real engine.
type fakePutEngine struct {
	fakeService
	putCh  chan enginev2.RoutedTriggerEvent
	active atomic.Int32
}

func newFakePutEngine() *fakePutEngine {
	return &fakePutEngine{putCh: make(chan enginev2.RoutedTriggerEvent, 4)}
}

func (f *fakePutEngine) Put(_ context.Context, event enginev2.RoutedTriggerEvent) error {
	f.putCh <- event
	return nil
}

func (f *fakePutEngine) ActiveExecutions() int32 { return f.active.Load() }

var (
	_ eventSink                = (*fakePutEngine)(nil)
	_ activeExecutionsReporter = (*fakePutEngine)(nil)
)

func newTestCoordinator(t *testing.T, capReg *regmocks.CapabilitiesRegistry, registry *EngineRegistry) TriggerCoordinator {
	t.Helper()
	coordinatorMetrics, err := monitoring.InitMonitoringResources()
	require.NoError(t, err)
	metricsLabeler := monitoring.NewWorkflowsMetricLabeler(commonmetrics.NewLabeler(), coordinatorMetrics)
	d := NewTriggerCoordinator(logger.TestLogger(t), capReg, registry, limits.NewTimeLimiter(5*time.Second), metricsLabeler, clockwork.NewFakeClock())
	require.NoError(t, d.Start(t.Context()))
	t.Cleanup(func() { _ = d.Close() })
	return d
}

func testWorkflowID(b byte) types.WorkflowID {
	return types.WorkflowID([32]byte{b})
}

func testSub(id string) *sdkpb.TriggerSubscription {
	return &sdkpb.TriggerSubscription{Id: id, Method: "Trigger"}
}

// Each subscription is registered with the capability registry, the returned
// trigger capability IDs are handed back in subscription order, and a reader
// is started per subscription.
func Test_RegisterTriggers_Success(t *testing.T) {
	t.Parallel()

	capReg := regmocks.NewCapabilitiesRegistry(t)
	registry := NewEngineRegistry()
	d := newTestCoordinator(t, capReg, registry)

	trigger0, trigger1 := capmocks.NewTriggerCapability(t), capmocks.NewTriggerCapability(t)
	capReg.EXPECT().GetTrigger(mock.Anything, "id_0").Return(trigger0, nil).Once()
	capReg.EXPECT().GetTrigger(mock.Anything, "id_1").Return(trigger1, nil).Once()
	trigger0.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).Return(make(chan capabilities.TriggerResponse), nil).Once()
	trigger1.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).Return(make(chan capabilities.TriggerResponse), nil).Once()

	wid := testWorkflowID(1)
	cre := contexts.CRE{Owner: "owner-a", Workflow: wid.Hex()}
	triggerIDs, err := d.RegisterTriggers(t.Context(), cre, RegistrationParams{WorkflowName: "wf"}, []*sdkpb.TriggerSubscription{testSub("id_0"), testSub("id_1")})
	require.NoError(t, err)
	require.Equal(t, []string{"id_0", "id_1"}, triggerIDs)

	// Registered triggers stop delivering once unregistered; UnregisterTrigger
	// must be called for both on cleanup regardless of order.
	trigger0.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()
	trigger1.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()
	require.NoError(t, d.UnregisterTriggers(wid.Hex()))
}

// If any subscription fails to register, every subscription that DID successfully
// register is unregistered again, and no handles are retained for the workflow
// (a subsequent Ack for a "successful" registration must fail to find it).
func Test_RegisterTriggers_RollbackOnFailure(t *testing.T) {
	t.Parallel()

	capReg := regmocks.NewCapabilitiesRegistry(t)
	registry := NewEngineRegistry()
	d := newTestCoordinator(t, capReg, registry)

	trigger0, trigger1 := capmocks.NewTriggerCapability(t), capmocks.NewTriggerCapability(t)
	capReg.EXPECT().GetTrigger(mock.Anything, "id_0").Return(trigger0, nil).Once()
	capReg.EXPECT().GetTrigger(mock.Anything, "id_1").Return(trigger1, nil).Once()
	trigger0.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).Return(make(chan capabilities.TriggerResponse), nil).Once()
	trigger1.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).Return(nil, errors.New("registration failed")).Once()
	trigger0.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()

	wid := testWorkflowID(2)
	cre := contexts.CRE{Owner: "owner-a", Workflow: wid.Hex()}
	triggerIDs, err := d.RegisterTriggers(t.Context(), cre, RegistrationParams{}, []*sdkpb.TriggerSubscription{testSub("id_0"), testSub("id_1")})
	require.Error(t, err)
	require.Nil(t, triggerIDs)

	// No handle was retained for id_0 either: Ack for it must fail.
	regID := enginev2.TriggerRegistrationID(wid.Hex(), 0)
	err = d.Ack(t.Context(), wid.Hex(), "id_0", regID, "event-1")
	require.Error(t, err)
}

// Test_RegisterTriggers_ReaderDeliversToEngineViaRegistry covers "reading"
// and "engine lookup via registry": the coordinator never holds an engine
// reference — its reader goroutine resolves the target engine from the
// EngineRegistry at delivery time, by workflow ID, for every event.
func Test_RegisterTriggers_ReaderDeliversToEngineViaRegistry(t *testing.T) {
	t.Parallel()

	capReg := regmocks.NewCapabilitiesRegistry(t)
	registry := NewEngineRegistry()
	d := newTestCoordinator(t, capReg, registry)

	trigger := capmocks.NewTriggerCapability(t)
	capReg.EXPECT().GetTrigger(mock.Anything, "id_0").Return(trigger, nil).Once()
	eventCh := make(chan capabilities.TriggerResponse, 1)
	trigger.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).Return(eventCh, nil).Once()

	wid := testWorkflowID(3)
	engine := newFakePutEngine()
	require.NoError(t, registry.Add(wid, "test-source", engine))

	cre := contexts.CRE{Owner: "owner-a", Workflow: wid.Hex()}
	_, err := d.RegisterTriggers(t.Context(), cre, RegistrationParams{}, []*sdkpb.TriggerSubscription{testSub("id_0")})
	require.NoError(t, err)

	eventCh <- capabilities.TriggerResponse{Event: capabilities.TriggerEvent{TriggerType: "basic-trigger@1.0.0", ID: "event-1"}}

	select {
	case routed := <-engine.putCh:
		require.Equal(t, wid.String(), routed.WorkflowID)
		require.Equal(t, "id_0", routed.TriggerCapID)
		require.Equal(t, "event-1", routed.Event.Event.ID)
	case <-time.After(2 * time.Second):
		t.Fatal("event was not delivered to the engine resolved from the registry")
	}
}

// If the engine isn't in the registry (e.g. it was popped before this event arrived),
// the reader must not panic or deliver anywhere — it drops the event and returns.
func Test_RegisterTriggers_ReaderExitsWhenEngineGone(t *testing.T) {
	t.Parallel()

	capReg := regmocks.NewCapabilitiesRegistry(t)
	registry := NewEngineRegistry()
	d := newTestCoordinator(t, capReg, registry)

	trigger := capmocks.NewTriggerCapability(t)
	capReg.EXPECT().GetTrigger(mock.Anything, "id_0").Return(trigger, nil).Once()
	eventCh := make(chan capabilities.TriggerResponse, 1)
	trigger.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).Return(eventCh, nil).Once()

	wid := testWorkflowID(4)
	// Deliberately never added to the registry.
	cre := contexts.CRE{Owner: "owner-a", Workflow: wid.Hex()}
	_, err := d.RegisterTriggers(t.Context(), cre, RegistrationParams{}, []*sdkpb.TriggerSubscription{testSub("id_0")})
	require.NoError(t, err)

	// Must not panic or hang; give the reader goroutine a moment to observe
	// the event and exit.
	eventCh <- capabilities.TriggerResponse{Event: capabilities.TriggerEvent{TriggerType: "basic-trigger@1.0.0", ID: "event-1"}}
	time.Sleep(100 * time.Millisecond)
}

// Ack resolves the registration's handle and calls AckEvent on the underlying
// trigger capability, the point at which the event is confirmed handled.
func Test_Ack_DelegatesToTriggerCapability(t *testing.T) {
	t.Parallel()

	capReg := regmocks.NewCapabilitiesRegistry(t)
	registry := NewEngineRegistry()
	d := newTestCoordinator(t, capReg, registry)

	trigger := capmocks.NewTriggerCapability(t)
	capReg.EXPECT().GetTrigger(mock.Anything, "id_0").Return(trigger, nil).Once()
	trigger.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).Return(make(chan capabilities.TriggerResponse), nil).Once()

	wid := testWorkflowID(5)
	cre := contexts.CRE{Owner: "owner-a", Workflow: wid.Hex()}
	_, err := d.RegisterTriggers(t.Context(), cre, RegistrationParams{}, []*sdkpb.TriggerSubscription{testSub("id_0")})
	require.NoError(t, err)

	regID := enginev2.TriggerRegistrationID(wid.Hex(), 0)
	trigger.EXPECT().AckEvent(mock.Anything, regID, "event-1", "Trigger").Return(nil).Once()
	require.NoError(t, d.Ack(t.Context(), wid.Hex(), "id_0", regID, "event-1"))
}

// An eventID/registrationID the coordinator never registered (or already
// released) resolves to nothing, and Ack must return an error rather than
// panicking or silently succeeding.
func Test_Ack_UnknownRegistration_ReturnsError(t *testing.T) {
	t.Parallel()

	capReg := regmocks.NewCapabilitiesRegistry(t)
	registry := NewEngineRegistry()
	d := newTestCoordinator(t, capReg, registry)

	wid := testWorkflowID(9)
	err := d.Ack(t.Context(), wid.Hex(), "id_0", "trigger_reg_does_not_exist_0", "event-1")
	require.Error(t, err)
}

// Unregistering stops ingress immediately (UnregisterTrigger is called on the capability) but
// the handle is RETAINED for as long as the engine reports active executions, so an execution
// that was already in flight when unregistration started can still successfully Ack. Once the
// engine reports it has drained, UnregisterTriggers releases the handle in the background and
// that Ack becomes impossible.
func Test_UnregisterTriggers_RetainsHandlesUntilEngineDrains(t *testing.T) {
	t.Parallel()

	capReg := regmocks.NewCapabilitiesRegistry(t)
	registry := NewEngineRegistry()
	d := newTestCoordinator(t, capReg, registry)

	trigger := capmocks.NewTriggerCapability(t)
	capReg.EXPECT().GetTrigger(mock.Anything, "id_0").Return(trigger, nil).Once()
	trigger.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).Return(make(chan capabilities.TriggerResponse), nil).Once()

	wid := testWorkflowID(6)
	engine := newFakePutEngine()
	engine.active.Store(1) // an execution is in flight
	require.NoError(t, registry.Add(wid, "test-source", engine))

	cre := contexts.CRE{Owner: "owner-a", Workflow: wid.Hex()}
	_, err := d.RegisterTriggers(t.Context(), cre, RegistrationParams{}, []*sdkpb.TriggerSubscription{testSub("id_0")})
	require.NoError(t, err)

	// Unregistering stops ingress with the capability...
	trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()
	require.NoError(t, d.UnregisterTriggers(wid.Hex()))

	// ...but the handle is retained while the engine still reports an active
	// execution: it can still Ack successfully.
	regID := enginev2.TriggerRegistrationID(wid.Hex(), 0)
	trigger.EXPECT().AckEvent(mock.Anything, regID, "in-flight-event", "Trigger").Return(nil).Once()
	require.NoError(t, d.Ack(t.Context(), wid.Hex(), "id_0", regID, "in-flight-event"))

	// Once the engine reports it has drained, the coordinator releases the
	// handle in the background. From this point on, Ack for the same
	// registration must eventually fail — there is nothing left to resolve
	// it to. The poll loop's own Ack calls may land before release finishes,
	// so this expectation is optional (.Maybe()) and only exercised on those.
	trigger.EXPECT().AckEvent(mock.Anything, regID, "too-late-event", "Trigger").Return(nil).Maybe()
	engine.active.Store(0)
	require.Eventually(t, func() bool {
		return d.Ack(t.Context(), wid.Hex(), "id_0", regID, "too-late-event") != nil
	}, 3*time.Second, 20*time.Millisecond, "handle was never released after the engine drained")
}

// Unregistering a workflow the coordinator never registered (e.g. a duplicate/retried
// cleanup event racing with an already-completed one).
func Test_UnregisterTriggers_UnknownWorkflow_ReturnsError(t *testing.T) {
	t.Parallel()

	capReg := regmocks.NewCapabilitiesRegistry(t)
	registry := NewEngineRegistry()
	d := newTestCoordinator(t, capReg, registry)

	wid := testWorkflowID(8)
	require.Error(t, d.UnregisterTriggers(wid.Hex()))
}
