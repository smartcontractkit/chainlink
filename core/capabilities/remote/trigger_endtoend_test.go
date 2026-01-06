package remote_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/aggregation"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/e2etesting"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"

	"github.com/smartcontractkit/cre-sdk-go/internal_testing/capabilities/basictrigger"
)

func Test_RemoteTriggerCapability_EventForwarding(t *testing.T) {
	ctx := t.Context()

	var err error
	numCapabilityPeers := 10
	var underlyingTriggers []commoncap.TriggerCapability
	broadcaster := NewTriggerBroadcaster(t, 100, nil)
	underlyingTriggers, err = broadcaster.CreateTriggerCapabilities(numCapabilityPeers)
	require.NoError(t, err)

	f := uint8(4)
	cfg := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Second,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: uint32(f + 1),
		MessageExpiry:           100 * time.Second,
		MaxBatchSize:            1, // no batching
		BatchCollectionPeriod:   time.Second,
	}

	clientTriggers, _, _ := setupRemoteTriggers(t, cfg, underlyingTriggers, 10, f, numCapabilityPeers, f)

	req := commoncap.TriggerRegistrationRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}
	result := registerOnAllTriggers(ctx, clientTriggers, req)
	triggerRegistrations := result.GetAll()
	for i, reg := range triggerRegistrations {
		require.NoError(t, reg.err, "expected nil error for trigger registration at index %d", i)
	}

	payload := sendTriggerEvent(t, "Hello, ", broadcaster)

	waitForEventFromAllTriggers(t, triggerRegistrations, payload)
}

func Test_RemoteTriggerCapability_EventForwarding_WithTriggerID(t *testing.T) {
	ctx := t.Context()

	var err error
	numCapabilityPeers := 10
	var underlyingTriggers []commoncap.TriggerCapability
	broadcaster := NewTriggerBroadcaster(t, 100, nil)
	underlyingTriggers, err = broadcaster.CreateTriggerCapabilities(numCapabilityPeers)
	require.NoError(t, err)

	f := uint8(4)
	cfg := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Second,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: uint32(f + 1),
		MessageExpiry:           100 * time.Second,
		MaxBatchSize:            1, // no batching
		BatchCollectionPeriod:   time.Second,
	}

	clientTriggers, _, _ := setupRemoteTriggers(t, cfg, underlyingTriggers, 10, f, numCapabilityPeers, f)

	req := commoncap.TriggerRegistrationRequest{
		TriggerID: "custom-trigger-id-123",
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}
	result := registerOnAllTriggers(ctx, clientTriggers, req)
	triggerRegistrations := result.GetAll()
	for i, reg := range triggerRegistrations {
		require.NoError(t, reg.err, "expected nil error for trigger registration at index %d", i)
	}

	payload := sendTriggerEvent(t, "Hello, ", broadcaster)

	waitForEventFromAllTriggers(t, triggerRegistrations, payload)
}

func Test_RemoteTriggerCapability_WaitForTriggerRegistration(t *testing.T) {
	ctx := t.Context()

	var err error
	numCapabilityPeers := 10
	var underlyingTriggers []commoncap.TriggerCapability
	broadcaster := NewTriggerBroadcaster(t, 100, nil)
	underlyingTriggers, err = broadcaster.CreateTriggerCapabilities(numCapabilityPeers)
	require.NoError(t, err)

	f := uint8(4)
	cfg := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Second,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: uint32(f + 1),
		MessageExpiry:           100 * time.Second,
		MaxBatchSize:            1, // no batching
		BatchCollectionPeriod:   time.Second,
	}

	clientTriggers, _, _ := setupRemoteTriggers(t, cfg, underlyingTriggers, 10, f, numCapabilityPeers, f)

	req := commoncap.TriggerRegistrationRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}
	result := registerOnAllTriggers(ctx, clientTriggers, req)
	triggerRegistrations := result.GetAll()
	for i, reg := range triggerRegistrations {
		require.NoError(t, reg.err, "expected nil error for trigger registration at index %d", i)
	}

	payload := sendTriggerEvent(t, "Hello, ", broadcaster)

	waitForEventFromAllTriggers(t, triggerRegistrations, payload)
}

func Test_RemoteTriggerCapability_WaitForTriggerRegistrationCapabilityError(t *testing.T) {
	ctx := t.Context()

	capError := caperrors.NewPublicUserError(errors.New("some error"), caperrors.Unknown)

	var err error
	numCapabilityPeers := 10
	var underlyingTriggers []commoncap.TriggerCapability
	broadcaster := NewTriggerBroadcaster(t, 100, capError)
	underlyingTriggers, err = broadcaster.CreateTriggerCapabilities(numCapabilityPeers)
	require.NoError(t, err)

	f := uint8(4)
	cfg := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Second,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: uint32(f + 1),
		MessageExpiry:           100 * time.Second,
		MaxBatchSize:            1, // no batching
		BatchCollectionPeriod:   time.Second,
	}

	clientTriggers, _, _ := setupRemoteTriggers(t, cfg, underlyingTriggers, 10, f, numCapabilityPeers, f)

	req := commoncap.TriggerRegistrationRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}
	result := registerOnAllTriggers(ctx, clientTriggers, req)
	triggerRegistrations := result.GetAll()
	for i, reg := range triggerRegistrations {
		require.Equal(t, capError, reg.err, "expected error for trigger registration at index %d", i)
	}
}

func Test_RemoteTriggerCapability_WaitForTriggerRegistrationError(t *testing.T) {
	ctx := t.Context()

	regError := errors.New("some error")

	var err error
	numCapabilityPeers := 10
	var underlyingTriggers []commoncap.TriggerCapability
	broadcaster := NewTriggerBroadcaster(t, 100, regError)
	underlyingTriggers, err = broadcaster.CreateTriggerCapabilities(numCapabilityPeers)
	require.NoError(t, err)

	f := uint8(4)
	cfg := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Second,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: uint32(f + 1),
		MessageExpiry:           100 * time.Second,
		MaxBatchSize:            1, // no batching
		BatchCollectionPeriod:   time.Second,
	}

	clientTriggers, _, _ := setupRemoteTriggers(t, cfg, underlyingTriggers, 10, f, numCapabilityPeers, f)

	req := commoncap.TriggerRegistrationRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}

	expectedError := caperrors.NewPublicSystemError(regError, caperrors.Unknown)
	result := registerOnAllTriggers(ctx, clientTriggers, req)
	triggerRegistrations := result.GetAll()
	for i, reg := range triggerRegistrations {
		require.Equal(t, expectedError, reg.err, "expected error for trigger registration at index %d", i)
	}
}

func Test_RemoteTriggerCapability_WaitForTriggerRegistrationTimeout(t *testing.T) {
	ctx := t.Context()

	var err error
	numCapabilityPeers := 10
	var underlyingTriggers []commoncap.TriggerCapability
	broadcaster := NewTriggerBroadcaster(t, 100, nil)
	underlyingTriggers, err = broadcaster.CreateTriggerCapabilities(numCapabilityPeers)
	require.NoError(t, err)

	f := uint8(4)
	cfg := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Second,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: uint32(f + 1),
		MessageExpiry:           100 * time.Second,
		MaxBatchSize:            1, // no batching
		BatchCollectionPeriod:   time.Second,
	}

	subscribers, _, broker := setupRemoteTriggers(t, cfg, underlyingTriggers, 10, f, numCapabilityPeers, f)

	// Block registration responses to simulate timeout when running against legacy don
	broker.SetMessageFilter(func(msg *types.MessageBody) bool {
		return msg.Method != types.TriggerRegistrationStatus
	})

	req := commoncap.TriggerRegistrationRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}

	result := registerOnAllTriggers(ctx, subscribers, req)
	triggerRegistrations := result.GetAll()
	for i, reg := range triggerRegistrations {
		require.ErrorIs(t, commoncap.ErrUnableToDetermineRegistrationStatus, reg.err, "expected error for trigger registration at index %d", i)
	}
}

// Test that the underlying triggers unregister as expected when the trigger subscriber stops send registration requests.
func Test_RemoteTriggerCapability_TriggerUnregisters(t *testing.T) {
	ctx := t.Context()

	var err error
	numCapabilityPeers := 10
	var underlyingTriggers []commoncap.TriggerCapability
	broadcaster := NewTriggerBroadcaster(t, 100, nil)
	underlyingTriggers, err = broadcaster.CreateTriggerCapabilities(numCapabilityPeers)
	require.NoError(t, err)

	f := uint8(4)
	cfg := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      100 * time.Millisecond,
		MinResponsesToAggregate: uint32(f + 1),
		MessageExpiry:           100 * time.Millisecond,
		MaxBatchSize:            1, // no batching
		BatchCollectionPeriod:   time.Second,
	}

	subscribers, _, broker := setupRemoteTriggers(t, cfg, underlyingTriggers, 10, f, numCapabilityPeers, f)

	req := commoncap.TriggerRegistrationRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}

	result := registerOnAllTriggers(ctx, subscribers, req)
	triggerRegistrations := result.GetAll()
	for i, reg := range triggerRegistrations {
		require.NoError(t, reg.err, "expected error for trigger registration at index %d", i)
	}

	waitForTriggerRegistration(t, underlyingTriggers)

	// Now block all trigger registration requests to simulate stopping the subscribers
	broker.SetMessageFilter(func(msg *types.MessageBody) bool {
		return msg.Method != types.MethodRegisterTrigger
	})

	// Verify that all underlying triggers unregister
	waitForTriggersToUnregister(t, underlyingTriggers)
}

// Test that the underlying triggers unregister as expected when the trigger subscriber stops send registration requests.
func Test_RemoteTriggerCapability_TriggerUnregisters_WithTriggerID(t *testing.T) {
	ctx := t.Context()

	var err error
	numCapabilityPeers := 10
	var underlyingTriggers []commoncap.TriggerCapability
	broadcaster := NewTriggerBroadcaster(t, 100, nil)
	underlyingTriggers, err = broadcaster.CreateTriggerCapabilities(numCapabilityPeers)
	require.NoError(t, err)

	f := uint8(4)
	cfg := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      100 * time.Millisecond,
		MinResponsesToAggregate: uint32(f + 1),
		MessageExpiry:           100 * time.Millisecond,
		MaxBatchSize:            1, // no batching
		BatchCollectionPeriod:   time.Second,
	}

	subscribers, _, broker := setupRemoteTriggers(t, cfg, underlyingTriggers, 10, f, numCapabilityPeers, f)

	req := commoncap.TriggerRegistrationRequest{
		TriggerID: "custom-trigger-id-123",
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}

	result := registerOnAllTriggers(ctx, subscribers, req)
	triggerRegistrations := result.GetAll()
	for i, reg := range triggerRegistrations {
		require.NoError(t, reg.err, "expected error for trigger registration at index %d", i)
	}

	waitForTriggerRegistration(t, underlyingTriggers)

	// Now block all trigger registration requests to simulate stopping the subscribers
	broker.SetMessageFilter(func(msg *types.MessageBody) bool {
		return msg.Method != types.MethodRegisterTrigger
	})

	// Verify that all underlying triggers unregister
	waitForTriggersToUnregister(t, underlyingTriggers)
}

func waitForTriggersToUnregister(t *testing.T, underlyingTriggers []commoncap.TriggerCapability) {
	for i := 0; i < len(underlyingTriggers); i++ {
		underlyingTrigger := underlyingTriggers[i].(*broadcasterTestTrigger)
		require.Eventually(t, func() bool {
			return !underlyingTrigger.IsTriggerRegistered()
		}, 10*time.Second, 100*time.Millisecond, "underlying trigger did not have UnregisterTrigger called")
	}
}

func waitForEventFromAllTriggers(t *testing.T, triggerRegistrations []registerTriggerResponse, payload *anypb.Any) {
	for i := 0; i < len(triggerRegistrations); i++ {
		triggerRegistration := triggerRegistrations[i]
		require.Eventually(t, func() bool {
			select {
			case response := <-triggerRegistration.respCh:
				return proto.Equal(response.Event.Payload, payload)
			default:
				return false
			}
		}, 10*time.Second, 100*time.Millisecond, "did not receive expected trigger response on client trigger %d", i)
	}
}

func sendTriggerEvent(t *testing.T, payloadString string, broadcaster *triggerBroadcaster) *anypb.Any {
	trigger := &basictrigger.Outputs{CoolOutput: payloadString}
	payload, err := anypb.New(trigger)
	require.NoError(t, err)

	broadcaster.SendTriggerResponse(commoncap.TriggerResponse{
		Event: commoncap.TriggerEvent{
			TriggerType: "",
			ID:          "",
			Outputs:     nil,
			Payload:     payload,
			OCREvent:    nil,
		},
		Err: nil,
	})
	return payload
}

func waitForTriggerRegistration(t *testing.T, underlyingTriggers []commoncap.TriggerCapability) {
	for i := 0; i < len(underlyingTriggers); i++ {
		underlyingTrigger := underlyingTriggers[i].(*broadcasterTestTrigger)
		require.Eventually(t, func() bool {
			return underlyingTrigger.IsTriggerRegistered()
		}, 10*time.Second, 100*time.Millisecond, "underlying trigger did not have RegisterTrigger called")
	}
}

func registerOnAllTriggers(ctx context.Context, clientTriggers []triggerSubscriber, req commoncap.TriggerRegistrationRequest) *ThreadSafeSlice[registerTriggerResponse] {
	result := &ThreadSafeSlice[registerTriggerResponse]{}
	var wg sync.WaitGroup
	for _, trigger := range clientTriggers {
		wg.Add(1)
		go func(trig commoncap.TriggerCapability) {
			defer wg.Done()
			respCh, err := trigger.RegisterTrigger(ctx, req)
			result.Append(registerTriggerResponse{
				respCh: respCh,
				err:    err,
			})
		}(trigger)
	}
	wg.Wait()
	return result
}

type ThreadSafeSlice[T any] struct {
	mu    sync.Mutex
	slice []T
}

func (s *ThreadSafeSlice[T]) Append(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.slice = append(s.slice, item)
}

func (s *ThreadSafeSlice[T]) GetAll() []T {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]T, len(s.slice))
	copy(result, s.slice)
	return result
}

type triggerBroadcaster struct {
	t               *testing.T
	err             error
	triggers        []*broadcasterTestTrigger
	eventBufferSize int
}

func NewTriggerBroadcaster(t *testing.T, eventBufferSize int, err error) *triggerBroadcaster {
	return &triggerBroadcaster{
		t:               t,
		err:             err,
		eventBufferSize: eventBufferSize,
	}
}

func (t *triggerBroadcaster) SetUnderlyingTriggerError(err error) {
	for _, trigger := range t.triggers {
		trigger.SetError(err)
	}
}

func (t *triggerBroadcaster) CreateTriggerCapabilities(num int) ([]commoncap.TriggerCapability, error) {
	var capabilities []commoncap.TriggerCapability
	for i := 0; i < num; i++ {
		capability, err := t.newTriggerCapability()
		if err != nil {
			return nil, err
		}
		capabilities = append(capabilities, capability)
	}
	return capabilities, nil
}

func (t *triggerBroadcaster) newTriggerCapability() (commoncap.TriggerCapability, error) {
	trigger := &broadcasterTestTrigger{
		t:   t.t,
		err: t.err,
		ch:  make(chan commoncap.TriggerResponse, t.eventBufferSize),
	}
	t.triggers = append(t.triggers, trigger)

	return trigger, nil
}

func (t *triggerBroadcaster) SendTriggerResponse(resp commoncap.TriggerResponse) {
	for _, trigger := range t.triggers {
		trigger.ch <- resp
	}
}

type broadcasterTestTrigger struct {
	mu sync.Mutex

	t   *testing.T
	err error
	ch  chan commoncap.TriggerResponse

	registered atomic.Bool
}

func (t *broadcasterTestTrigger) SetError(err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.err = err
}

func (t *broadcasterTestTrigger) GetError() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.err
}

func (t *broadcasterTestTrigger) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	return commoncap.CapabilityInfo{
		ID: "broadcast-trigger-capability",
	}, nil
}

type registerTriggerResponse struct {
	respCh <-chan commoncap.TriggerResponse
	err    error
}

func (t *broadcasterTestTrigger) RegisterTrigger(ctx context.Context, req commoncap.TriggerRegistrationRequest) (<-chan commoncap.TriggerResponse, error) {
	require.NotNil(t.t, req)
	t.registered.Store(true)
	return t.ch, t.GetError()
}

func (t *broadcasterTestTrigger) IsTriggerRegistered() bool {
	return t.registered.Load()
}

func (t *broadcasterTestTrigger) UnregisterTrigger(ctx context.Context, req commoncap.TriggerRegistrationRequest) error {
	require.NotNil(t.t, req)
	t.registered.Store(false)
	return nil
}

func (t *broadcasterTestTrigger) AckEvent(ctx context.Context, triggerID string, eventID string, method string) error {
	return nil
}

type triggerSubscriber interface {
	commoncap.TriggerCapability
	Start(ctx context.Context) error
	Close() error
}

type triggerPublisher interface {
	Close() error
}

func setupRemoteTriggers(t *testing.T, cfg *commoncap.RemoteTriggerConfig, underlying []commoncap.TriggerCapability, numWorkflowPeers int, workflowDonF uint8,
	numCapabilityPeers int, capabilityDonF uint8) ([]triggerSubscriber, []triggerPublisher, *e2etesting.TestAsyncMessageBroker) {
	lggr := logger.Test(t)

	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := range numCapabilityPeers {
		capabilityPeerID := p2ptypes.PeerID{}
		require.NoError(t, capabilityPeerID.UnmarshalText([]byte(e2etesting.NewPeerID())))
		capabilityPeers[i] = capabilityPeerID
	}

	capabilityPeerID := p2ptypes.PeerID{}
	require.NoError(t, capabilityPeerID.UnmarshalText([]byte(e2etesting.NewPeerID())))

	capDonInfo := commoncap.DON{
		ID:      2,
		Members: capabilityPeers,
		F:       capabilityDonF,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1.0.0",
		CapabilityType: commoncap.CapabilityTypeTrigger,
		Description:    "Remote Target",
		DON:            &capDonInfo,
	}

	workflowPeers := make([]p2ptypes.PeerID, numWorkflowPeers)
	for i := range numWorkflowPeers {
		workflowPeerID := p2ptypes.PeerID{}
		require.NoError(t, workflowPeerID.UnmarshalText([]byte(e2etesting.NewPeerID())))
		workflowPeers[i] = workflowPeerID
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      1,
		F:       workflowDonF,
	}

	broker := e2etesting.NewTestAsyncMessageBroker(t, 1000)

	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}

	capabilityNodeTriggerPublishers := make([]triggerPublisher, numCapabilityPeers)

	if len(underlying) != numCapabilityPeers {
		t.Fatalf("expected %d underlying capability triggers, got %d", numCapabilityPeers, len(underlying))
	}

	for i := range numCapabilityPeers {
		capabilityPeer := capabilityPeers[i]
		capabilityDispatcher := broker.NewDispatcherForNode(capabilityPeer)
		capabilityNode := remote.NewTriggerPublisher(capInfo.ID, "", capabilityDispatcher, lggr)
		require.NoError(t, capabilityNode.SetConfig(cfg, underlying[i], capDonInfo, workflowDONs))
		//	if runPublishers {
		servicetest.Run(t, capabilityNode)
		//	}
		broker.RegisterReceiverNode(capabilityPeer, capabilityNode)
		capabilityNodeTriggerPublishers[i] = capabilityNode
	}

	workflowNodeTriggerSubscribers := make([]triggerSubscriber, numWorkflowPeers)
	for i := range numWorkflowPeers {
		workflowPeerDispatcher := broker.NewDispatcherForNode(workflowPeers[i])
		workflowNode := remote.NewTriggerSubscriber(capInfo.ID, "", workflowPeerDispatcher, lggr,
			limits.NewTimeLimiter(60*time.Second))
		agg := aggregation.NewDefaultModeAggregator(cfg.MinResponsesToAggregate)
		err := workflowNode.SetConfig(cfg, capInfo, workflowDonInfo.ID, capDonInfo, agg)
		require.NoError(t, err)
		servicetest.Run(t, workflowNode)
		broker.RegisterReceiverNode(workflowPeers[i], workflowNode)
		workflowNodeTriggerSubscribers[i] = workflowNode
	}

	servicetest.Run(t, broker)

	return workflowNodeTriggerSubscribers, capabilityNodeTriggerPublishers, broker
}
