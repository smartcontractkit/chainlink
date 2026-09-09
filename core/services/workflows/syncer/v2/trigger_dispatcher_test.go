package v2

import (
	"context"
	"errors"
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

// fakePutEngine is a services.Service that also implements the dispatcher's
// narrow eventSink interface, so a RegisterTriggers reader can resolve and
// deliver to it via the EngineRegistry, exactly as it would a real engine.
type fakePutEngine struct {
	fakeService
	putCh chan enginev2.RoutedTriggerEvent
}

func newFakePutEngine() *fakePutEngine {
	return &fakePutEngine{putCh: make(chan enginev2.RoutedTriggerEvent, 4)}
}

func (f *fakePutEngine) Put(_ context.Context, event enginev2.RoutedTriggerEvent) error {
	f.putCh <- event
	return nil
}

var _ eventSink = (*fakePutEngine)(nil)

func newTestDispatcher(t *testing.T, capReg *regmocks.CapabilitiesRegistry, registry *EngineRegistry) TriggerDispatcher {
	t.Helper()
	dispatcherMetrics, err := monitoring.InitMonitoringResources()
	require.NoError(t, err)
	metricsLabeler := monitoring.NewWorkflowsMetricLabeler(commonmetrics.NewLabeler(), dispatcherMetrics)
	d := NewTriggerDispatcher(logger.TestLogger(t), capReg, registry, limits.NewTimeLimiter(5*time.Second), metricsLabeler, clockwork.NewFakeClock())
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

// Test_RegisterTriggers_Success covers AC 7's "registration": each
// subscription is registered with the capability registry, the returned
// trigger capability IDs are handed back in subscription order, and a reader
// is started per subscription (verified in a separate test below).
func Test_RegisterTriggers_Success(t *testing.T) {
	t.Parallel()

	capReg := regmocks.NewCapabilitiesRegistry(t)
	registry := NewEngineRegistry()
	d := newTestDispatcher(t, capReg, registry)

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

// Test_RegisterTriggers_RollbackOnFailure covers the rollback path: if any
// subscription fails to register (as opposed to failing the earlier
// "does this trigger exist" validation, which fails before any registration
// is attempted), every subscription that DID successfully register is
// unregistered again, and no handles are retained for the workflow (a
// subsequent Ack for a "successful" registration must fail to find it).
func Test_RegisterTriggers_RollbackOnFailure(t *testing.T) {
	t.Parallel()

	capReg := regmocks.NewCapabilitiesRegistry(t)
	registry := NewEngineRegistry()
	d := newTestDispatcher(t, capReg, registry)

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
	err = d.Ack(t.Context(), "id_0", regID, "event-1")
	require.Error(t, err)
}

// Test_RegisterTriggers_ReaderDeliversToEngineViaRegistry covers "reading"
// and "engine lookup via registry": the dispatcher never holds an engine
// reference — its reader goroutine resolves the target engine from the
// EngineRegistry at delivery time, by workflow ID, for every event.
func Test_RegisterTriggers_ReaderDeliversToEngineViaRegistry(t *testing.T) {
	t.Parallel()

	capReg := regmocks.NewCapabilitiesRegistry(t)
	registry := NewEngineRegistry()
	d := newTestDispatcher(t, capReg, registry)

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

// Test_RegisterTriggers_ReaderExitsWhenEngineGone covers the "never holds an
// engine reference" half of AC 1: if the engine isn't in the registry (e.g.
// it was popped before this event arrived), the reader must not panic or
// deliver anywhere — it silently drops the event and returns.
func Test_RegisterTriggers_ReaderExitsWhenEngineGone(t *testing.T) {
	t.Parallel()

	capReg := regmocks.NewCapabilitiesRegistry(t)
	registry := NewEngineRegistry()
	d := newTestDispatcher(t, capReg, registry)

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

// Test_Ack_DelegatesToTriggerCapability covers "ACK delegation": Ack resolves
// the registration's handle and calls AckEvent on the underlying trigger
// capability, the point at which the event is confirmed handled.
func Test_Ack_DelegatesToTriggerCapability(t *testing.T) {
	t.Parallel()

	capReg := regmocks.NewCapabilitiesRegistry(t)
	registry := NewEngineRegistry()
	d := newTestDispatcher(t, capReg, registry)

	trigger := capmocks.NewTriggerCapability(t)
	capReg.EXPECT().GetTrigger(mock.Anything, "id_0").Return(trigger, nil).Once()
	trigger.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).Return(make(chan capabilities.TriggerResponse), nil).Once()

	wid := testWorkflowID(5)
	cre := contexts.CRE{Owner: "owner-a", Workflow: wid.Hex()}
	_, err := d.RegisterTriggers(t.Context(), cre, RegistrationParams{}, []*sdkpb.TriggerSubscription{testSub("id_0")})
	require.NoError(t, err)

	regID := enginev2.TriggerRegistrationID(wid.Hex(), 0)
	trigger.EXPECT().AckEvent(mock.Anything, regID, "event-1", "Trigger").Return(nil).Once()
	require.NoError(t, d.Ack(t.Context(), "id_0", regID, "event-1"))
}

// Test_Ack_UnknownRegistration_ReturnsError covers Ack's failure path: an
// eventID/registrationID the dispatcher never registered (or already
// released) resolves to nothing, and Ack must return an error rather than
// panicking or silently succeeding.
func Test_Ack_UnknownRegistration_ReturnsError(t *testing.T) {
	t.Parallel()

	capReg := regmocks.NewCapabilitiesRegistry(t)
	registry := NewEngineRegistry()
	d := newTestDispatcher(t, capReg, registry)

	err := d.Ack(t.Context(), "id_0", "trigger_reg_does_not_exist_0", "event-1")
	require.Error(t, err)
}

// Test_UnregisterTriggers_RetainsHandlesForInFlightAck is the AC 7
// requirement most specific to CRE-6176's design: unregistering stops
// ingress immediately (UnregisterTrigger is called on the capability) but
// the handle is RETAINED, so an execution that was already in flight when
// unregistration started can still successfully Ack. Only ReleaseHandles
// (called by the syncer once the engine has fully drained and closed) makes
// that Ack impossible.
func Test_UnregisterTriggers_RetainsHandlesForInFlightAck(t *testing.T) {
	t.Parallel()

	capReg := regmocks.NewCapabilitiesRegistry(t)
	registry := NewEngineRegistry()
	d := newTestDispatcher(t, capReg, registry)

	trigger := capmocks.NewTriggerCapability(t)
	capReg.EXPECT().GetTrigger(mock.Anything, "id_0").Return(trigger, nil).Once()
	trigger.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).Return(make(chan capabilities.TriggerResponse), nil).Once()

	wid := testWorkflowID(6)
	cre := contexts.CRE{Owner: "owner-a", Workflow: wid.Hex()}
	_, err := d.RegisterTriggers(t.Context(), cre, RegistrationParams{}, []*sdkpb.TriggerSubscription{testSub("id_0")})
	require.NoError(t, err)

	// Unregistering stops ingress with the capability...
	trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()
	require.NoError(t, d.UnregisterTriggers(wid.Hex()))

	// ...but the handle is retained: an execution still in flight (started
	// before unregistration) can still Ack successfully.
	regID := enginev2.TriggerRegistrationID(wid.Hex(), 0)
	trigger.EXPECT().AckEvent(mock.Anything, regID, "in-flight-event", "Trigger").Return(nil).Once()
	require.NoError(t, d.Ack(t.Context(), "id_0", regID, "in-flight-event"))

	// Once the syncer has drained and closed the engine, it releases the
	// handles. From this point on, Ack for the same registration must fail —
	// there is nothing left to resolve it to.
	require.NoError(t, d.ReleaseHandles(wid.Hex()))
	err = d.Ack(t.Context(), "id_0", regID, "too-late-event")
	require.Error(t, err)
}

// Test_ReleaseHandles_IsIdempotent covers "unregister/cleanup": ReleaseHandles
// may be called more than once (e.g. a retried cleanup) without error.
func Test_ReleaseHandles_IsIdempotent(t *testing.T) {
	t.Parallel()

	capReg := regmocks.NewCapabilitiesRegistry(t)
	registry := NewEngineRegistry()
	d := newTestDispatcher(t, capReg, registry)

	wid := testWorkflowID(7)
	require.NoError(t, d.ReleaseHandles(wid.Hex()))
	require.NoError(t, d.ReleaseHandles(wid.Hex()))
}

// Test_UnregisterTriggers_UnknownWorkflow_ReturnsError covers unregistering a
// workflow the dispatcher never registered (e.g. a duplicate/retried cleanup
// event racing with an already-completed one).
func Test_UnregisterTriggers_UnknownWorkflow_ReturnsError(t *testing.T) {
	t.Parallel()

	capReg := regmocks.NewCapabilitiesRegistry(t)
	registry := NewEngineRegistry()
	d := newTestDispatcher(t, capReg, registry)

	wid := testWorkflowID(8)
	require.Error(t, d.UnregisterTriggers(wid.Hex()))
}
