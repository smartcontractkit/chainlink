package remote_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/aggregation"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/e2etesting"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// End-to-end tests for remote trigger registration status propagation between
// Workflow DON (TriggerSubscriber) and Capability DON (TriggerPublisher).
//
// Why this file exists:
// Unit tests in trigger_subscriber_test.go and trigger_publisher_test.go exercise
// each side in isolation with direct message injection. They do not prove that a
// real RegisterTrigger call on the WF side blocks until CAP peers return
// MethodTriggerRegistrationStatus, or that quorum aggregation matches across both
// DONs. These tests wire multiple subscriber and publisher nodes through an
// async message broker (e2etesting.TestAsyncMessageBroker) so don2don send/receive,
// registration quorum, and status round-trip run together without mocks.
//
// What this file should test:
//   - Success: WF RegisterTrigger returns after CAP quorum registers the trigger.
//   - User error: underlying.RegisterTrigger returns OriginUser → error propagates
//     to all WF subscribers; registration is terminal on the CAP side.
//   - System error: non-user failure → serialized system error returned to WF.
//   - Status timeout: when status messages are dropped, WF gets
//     ErrUnableToDetermineRegistrationStatus and still returns the event channel
//     (legacy optimistic path for mixed-version rollout).
//
// Intentionally out of scope here (covered elsewhere or deferred):
//   - Trigger event forwarding after registration
//   - UnregisterTrigger status flow
//   - Production-scale topology (10 nodes, F=4); this suite uses 4 CAP / 4 WF peers, F=1

func Test_RemoteTriggerRegistration_Success(t *testing.T) {
	ctx := t.Context()

	underlying, err := newUnderlyingTriggers(t, 4, nil)
	require.NoError(t, err)

	subscribers, _, _ := setupRegistrationE2E(t, underlying, 4, 1, limits.NewTimeLimiter(10*time.Second))

	req := commoncap.TriggerRegistrationRequest{
		TriggerID: "e2e-trigger",
		Metadata:  commoncap.RequestMetadata{WorkflowID: workflowID1},
	}
	for i, reg := range registerOnAllSubscribers(ctx, subscribers, req) {
		require.NoError(t, reg.err, "subscriber %d", i)
		require.NotNil(t, reg.respCh)
	}

	waitForUnderlyingRegistration(t, underlying)
}

func Test_RemoteTriggerRegistration_UserError(t *testing.T) {
	ctx := t.Context()

	capErr := caperrors.NewPublicUserError(errors.New("some error"), caperrors.Unknown)
	underlying, err := newUnderlyingTriggers(t, 4, capErr)
	require.NoError(t, err)

	subscribers, _, _ := setupRegistrationE2E(t, underlying, 4, 1, limits.NewTimeLimiter(10*time.Second))

	req := commoncap.TriggerRegistrationRequest{
		TriggerID: "e2e-trigger",
		Metadata:  commoncap.RequestMetadata{WorkflowID: workflowID1},
	}
	for i, reg := range registerOnAllSubscribers(ctx, subscribers, req) {
		require.Equal(t, capErr, reg.err, "subscriber %d", i)
	}
}

func Test_RemoteTriggerRegistration_SystemError(t *testing.T) {
	ctx := t.Context()

	regErr := errors.New("some error")
	underlying, err := newUnderlyingTriggers(t, 4, regErr)
	require.NoError(t, err)

	subscribers, _, _ := setupRegistrationE2E(t, underlying, 4, 1, limits.NewTimeLimiter(10*time.Second))

	req := commoncap.TriggerRegistrationRequest{
		TriggerID: "e2e-trigger",
		Metadata:  commoncap.RequestMetadata{WorkflowID: workflowID1},
	}
	expected := caperrors.NewPublicSystemError(regErr, caperrors.Unknown)
	for i, reg := range registerOnAllSubscribers(ctx, subscribers, req) {
		require.Equal(t, expected, reg.err, "subscriber %d", i)
	}
}

func Test_RemoteTriggerRegistration_StatusTimeout(t *testing.T) {
	ctx := t.Context()

	underlying, err := newUnderlyingTriggers(t, 4, nil)
	require.NoError(t, err)

	subscribers, _, broker := setupRegistrationE2E(t, underlying, 4, 1, limits.NewTimeLimiter(200*time.Millisecond))

	broker.SetMessageFilter(func(msg *remotetypes.MessageBody) bool {
		return msg.Method != remotetypes.MethodTriggerRegistrationStatus
	})

	req := commoncap.TriggerRegistrationRequest{
		TriggerID: "e2e-trigger",
		Metadata:  commoncap.RequestMetadata{WorkflowID: workflowID1},
	}
	for i, reg := range registerOnAllSubscribers(ctx, subscribers, req) {
		require.ErrorIs(t, reg.err, commoncap.ErrUnableToDetermineRegistrationStatus, "subscriber %d", i)
		require.NotNil(t, reg.respCh, "legacy optimistic path should still return event channel")
	}
}

type registerTriggerResult struct {
	respCh <-chan commoncap.TriggerResponse
	err    error
}

func registerOnAllSubscribers(ctx context.Context, subscribers []triggerSubscriber, req commoncap.TriggerRegistrationRequest) []registerTriggerResult {
	results := make([]registerTriggerResult, len(subscribers))
	var wg sync.WaitGroup
	for i, sub := range subscribers {
		wg.Add(1)
		go func(idx int, trig triggerSubscriber) {
			defer wg.Done()
			respCh, err := trig.RegisterTrigger(ctx, req)
			results[idx] = registerTriggerResult{respCh: respCh, err: err}
		}(i, sub)
	}
	wg.Wait()
	return results
}

type triggerSubscriber interface {
	commoncap.TriggerCapability
	Start(ctx context.Context) error
	Close() error
}

type underlyingTestTrigger struct {
	t          *testing.T
	err        error
	ch         chan commoncap.TriggerResponse
	registered atomic.Bool
}

func (t *underlyingTestTrigger) Info(context.Context) (commoncap.CapabilityInfo, error) {
	return commoncap.CapabilityInfo{ID: "underlying-trigger"}, nil
}

func (t *underlyingTestTrigger) RegisterTrigger(_ context.Context, req commoncap.TriggerRegistrationRequest) (<-chan commoncap.TriggerResponse, error) {
	require.NotNil(t.t, req)
	t.registered.Store(true)
	return t.ch, t.err
}

func (t *underlyingTestTrigger) UnregisterTrigger(_ context.Context, req commoncap.TriggerRegistrationRequest) error {
	require.NotNil(t.t, req)
	t.registered.Store(false)
	return nil
}

func (t *underlyingTestTrigger) AckEvent(context.Context, string, string, string) error {
	return nil
}

func (t *underlyingTestTrigger) isRegistered() bool {
	return t.registered.Load()
}

func newUnderlyingTriggers(t *testing.T, n int, regErr error) ([]commoncap.TriggerCapability, error) {
	triggers := make([]commoncap.TriggerCapability, n)
	for i := range n {
		triggers[i] = &underlyingTestTrigger{
			t:   t,
			err: regErr,
			ch:  make(chan commoncap.TriggerResponse, 4),
		}
	}
	return triggers, nil
}

func waitForUnderlyingRegistration(t *testing.T, underlying []commoncap.TriggerCapability) {
	t.Helper()
	for i, trig := range underlying {
		ut := trig.(*underlyingTestTrigger)
		require.Eventually(t, ut.isRegistered, 5*time.Second, 50*time.Millisecond,
			"underlying trigger %d was not registered", i)
	}
}

func setupRegistrationE2E(
	t *testing.T,
	underlying []commoncap.TriggerCapability,
	numWorkflowPeers int,
	donF uint8,
	statusTimeout limits.TimeLimiter,
) ([]triggerSubscriber, []interface{ Close() error }, *e2etesting.TestAsyncMessageBroker) {
	t.Helper()
	lggr := logger.Test(t)

	numCapabilityPeers := len(underlying)
	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := range numCapabilityPeers {
		capabilityPeers[i] = e2etesting.NewP2PPeerID(t)
	}

	capDonInfo := commoncap.DON{
		ID:      2,
		Members: capabilityPeers,
		F:       donF,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1.0.0",
		CapabilityType: commoncap.CapabilityTypeTrigger,
		Description:    "Remote Trigger",
		DON:            &capDonInfo,
	}

	workflowPeers := make([]p2ptypes.PeerID, numWorkflowPeers)
	for i := range numWorkflowPeers {
		workflowPeers[i] = e2etesting.NewP2PPeerID(t)
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      1,
		F:       donF,
	}

	cfg := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Second,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: uint32(donF + 1),
		MessageExpiry:           100 * time.Second,
		MaxBatchSize:            1,
		BatchCollectionPeriod:   time.Second,
	}

	broker := e2etesting.NewTestAsyncMessageBroker(t, 1000)
	workflowDONs := map[uint32]commoncap.DON{workflowDonInfo.ID: workflowDonInfo}

	publishers := make([]interface{ Close() error }, numCapabilityPeers)
	for i := range numCapabilityPeers {
		capabilityPeer := capabilityPeers[i]
		capabilityDispatcher := broker.NewDispatcherForNode(capabilityPeer)
		publisher := remote.NewTriggerPublisher(capInfo.ID, "", capabilityDispatcher, lggr)
		require.NoError(t, publisher.SetConfig(cfg, underlying[i], capDonInfo, workflowDONs))
		servicetest.Run(t, publisher)
		broker.RegisterReceiverNode(capabilityPeer, publisher)
		publishers[i] = publisher
	}

	subscribers := make([]triggerSubscriber, numWorkflowPeers)
	for i := range numWorkflowPeers {
		workflowDispatcher := broker.NewDispatcherForNode(workflowPeers[i])
		subscriber := remote.NewTriggerSubscriber(capInfo.ID, "", workflowDispatcher, lggr, statusTimeout)
		agg := aggregation.NewDefaultModeAggregator(cfg.MinResponsesToAggregate)
		require.NoError(t, subscriber.SetConfig(cfg, capInfo, workflowDonInfo.ID, capDonInfo, agg))
		servicetest.Run(t, subscriber)
		broker.RegisterReceiverNode(workflowPeers[i], subscriber)
		subscribers[i] = subscriber
	}

	servicetest.Run(t, broker)

	return subscribers, publishers, broker
}
