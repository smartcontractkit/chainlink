package remote_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const capID = "cap_id@1"

func TestTriggerPublisher_Register(t *testing.T) {
	ctx := testutils.Context(t)
	capabilityDONID, workflowDONID := uint32(1), uint32(2)

	underlyingTriggerCap, publisher, _, peers := newServices(t, capabilityDONID, workflowDONID, 1)

	// invalid sender case - node 0 is not a member of the workflow DON, registration shoudn't happen
	regEvent := newRegisterTriggerMessage(t, workflowDONID, peers[0])
	publisher.Receive(ctx, regEvent)
	require.Empty(t, underlyingTriggerCap.registrationsCh)

	// valid registration
	regEvent = newRegisterTriggerMessage(t, workflowDONID, peers[1])
	publisher.Receive(ctx, regEvent)
	require.NotEmpty(t, underlyingTriggerCap.registrationsCh)
	forwarded := <-underlyingTriggerCap.registrationsCh
	require.Equal(t, workflowID1, forwarded.Metadata.WorkflowID)

	require.NoError(t, publisher.Close())
}

func TestTriggerPublisher_ReceiveTriggerEvents_NoBatching(t *testing.T) {
	ctx := testutils.Context(t)
	capabilityDONID, workflowDONID := uint32(1), uint32(2)

	underlyingTriggerCap, publisher, dispatcher, peers := newServices(t, capabilityDONID, workflowDONID, 1)
	regEvent := newRegisterTriggerMessage(t, workflowDONID, peers[1])
	publisher.Receive(ctx, regEvent)
	require.NotEmpty(t, underlyingTriggerCap.registrationsCh)

	// send a trigger event and expect that it gets delivered right away
	awaitOutgoingMessageCh := make(chan struct{})
	dispatcher.On("Send", peers[1], mock.Anything).Run(func(args mock.Arguments) {
		awaitOutgoingMessageCh <- struct{}{}
	}).Return(nil)
	underlyingTriggerCap.eventCh <- commoncap.TriggerResponse{}
	<-awaitOutgoingMessageCh

	require.NoError(t, publisher.Close())
}

func TestTriggerPublisher_ReceiveTriggerEvents_BatchingEnabled(t *testing.T) {
	ctx := testutils.Context(t)
	capabilityDONID, workflowDONID := uint32(1), uint32(2)

	underlyingTriggerCap, publisher, dispatcher, peers := newServices(t, capabilityDONID, workflowDONID, 2)
	regEvent := newRegisterTriggerMessage(t, workflowDONID, peers[1])
	publisher.Receive(ctx, regEvent)
	require.NotEmpty(t, underlyingTriggerCap.registrationsCh)

	// send two trigger events and expect them to be delivered in a batch
	awaitOutgoingMessageCh := make(chan struct{})
	dispatcher.On("Send", peers[1], mock.Anything).Run(func(args mock.Arguments) {
		msg := args.Get(1).(*remotetypes.MessageBody)
		require.Equal(t, capID, msg.CapabilityId)
		require.Equal(t, remotetypes.MethodTriggerEvent, msg.Method)
		require.NotEmpty(t, msg.Payload)
		metadata := msg.Metadata.(*remotetypes.MessageBody_TriggerEventMetadata)
		require.Len(t, metadata.TriggerEventMetadata.WorkflowIds, 2)
		awaitOutgoingMessageCh <- struct{}{}
	}).Return(nil).Once()
	underlyingTriggerCap.eventCh <- commoncap.TriggerResponse{}
	underlyingTriggerCap.eventCh <- commoncap.TriggerResponse{}
	<-awaitOutgoingMessageCh

	// if there are fewer pending event than the batch size,
	// the events should still be sent after the batch collection period
	dispatcher.On("Send", peers[1], mock.Anything).Run(func(args mock.Arguments) {
		msg := args.Get(1).(*remotetypes.MessageBody)
		metadata := msg.Metadata.(*remotetypes.MessageBody_TriggerEventMetadata)
		require.Len(t, metadata.TriggerEventMetadata.WorkflowIds, 1)
		awaitOutgoingMessageCh <- struct{}{}
	}).Return(nil).Once()
	underlyingTriggerCap.eventCh <- commoncap.TriggerResponse{}
	<-awaitOutgoingMessageCh

	require.NoError(t, publisher.Close())
}

func TestTriggerPublisher_SetConfig_Basic(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	capInfo := commoncap.CapabilityInfo{
		ID:             capID,
		CapabilityType: commoncap.CapabilityTypeTrigger,
		Description:    "Remote Trigger",
	}
	peers := make([]p2ptypes.PeerID, 2)
	require.NoError(t, peers[0].UnmarshalText([]byte(peerID1)))
	require.NoError(t, peers[1].UnmarshalText([]byte(peerID2)))
	capDonInfo := commoncap.DON{
		ID:      1,
		Members: []p2ptypes.PeerID{peers[0]},
		F:       0,
	}
	workflowDonInfo := commoncap.DON{
		ID:      2,
		Members: []p2ptypes.PeerID{peers[1]},
		F:       0,
	}
	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}
	underlying := &testTrigger{
		info:            capInfo,
		registrationsCh: make(chan commoncap.TriggerRegistrationRequest, 2),
		eventCh:         make(chan commoncap.TriggerResponse, 2),
	}

	t.Run("returns error when underlying trigger capability is nil", func(t *testing.T) {
		dispatcher := mocks.NewDispatcher(t)
		publisher := remote.NewTriggerPublisher(capInfo.ID, "method", dispatcher, lggr)
		config := &commoncap.RemoteTriggerConfig{}
		err := publisher.SetConfig(config, nil, capDonInfo, workflowDONs)
		require.Error(t, err)
		require.Contains(t, err.Error(), "underlying trigger capability cannot be nil")
	})

	t.Run("handles nil config", func(t *testing.T) {
		dispatcher := mocks.NewDispatcher(t)
		publisher := remote.NewTriggerPublisher(capInfo.ID, "method", dispatcher, lggr)
		// Set config as nil - should use defaults
		err := publisher.SetConfig(nil, underlying, capDonInfo, workflowDONs)
		require.NoError(t, err)

		// Verify config works
		ctx := testutils.Context(t)
		require.NoError(t, publisher.Start(ctx))
		require.NoError(t, publisher.Close())
	})

	t.Run("handles nil workflowDONs", func(t *testing.T) {
		dispatcher := mocks.NewDispatcher(t)
		publisher := remote.NewTriggerPublisher(capInfo.ID, "method", dispatcher, lggr)
		config := &commoncap.RemoteTriggerConfig{
			RegistrationRefresh:     100 * time.Millisecond,
			RegistrationExpiry:      100 * time.Second,
			MinResponsesToAggregate: 1,
			MessageExpiry:           100 * time.Second,
		}
		// Set workflowDONs as nil - should create empty map
		err := publisher.SetConfig(config, underlying, capDonInfo, nil)
		require.NoError(t, err)

		// Verify config works
		ctx := testutils.Context(t)
		require.NoError(t, publisher.Start(ctx))
		require.NoError(t, publisher.Close())
	})

	t.Run("updates existing config", func(t *testing.T) {
		dispatcher := mocks.NewDispatcher(t)
		publisher := remote.NewTriggerPublisher(capInfo.ID, "method", dispatcher, lggr)
		// Set initial config
		initialConfig := &commoncap.RemoteTriggerConfig{
			RegistrationRefresh:     100 * time.Millisecond,
			RegistrationExpiry:      100 * time.Second,
			MinResponsesToAggregate: 1,
			MessageExpiry:           100 * time.Second,
			MaxBatchSize:            1,
			BatchCollectionPeriod:   100 * time.Millisecond,
		}
		err := publisher.SetConfig(initialConfig, underlying, capDonInfo, workflowDONs)
		require.NoError(t, err)

		// Update with new config
		newConfig := &commoncap.RemoteTriggerConfig{
			RegistrationRefresh:     500 * time.Millisecond,
			RegistrationExpiry:      500 * time.Second,
			MinResponsesToAggregate: 3,
			MessageExpiry:           500 * time.Second,
			MaxBatchSize:            5,
			BatchCollectionPeriod:   500 * time.Millisecond,
		}
		err = publisher.SetConfig(newConfig, underlying, capDonInfo, workflowDONs)
		require.NoError(t, err)

		// Verify updated config works
		ctx := testutils.Context(t)
		require.NoError(t, publisher.Start(ctx))
		require.NoError(t, publisher.Close())
	})
}

func newServices(t *testing.T, capabilityDONID uint32, workflowDONID uint32, maxBatchSize uint32) (*testTrigger, remotetypes.ReceiverService, *mocks.Dispatcher, []p2ptypes.PeerID) {
	lggr := logger.Test(t)
	ctx := testutils.Context(t)
	capInfo := commoncap.CapabilityInfo{
		ID:             capID,
		CapabilityType: commoncap.CapabilityTypeTrigger,
		Description:    "Remote Trigger",
	}
	peers := make([]p2ptypes.PeerID, 2)
	require.NoError(t, peers[0].UnmarshalText([]byte(peerID1)))
	require.NoError(t, peers[1].UnmarshalText([]byte(peerID2)))
	capDonInfo := commoncap.DON{
		ID:      capabilityDONID,
		Members: []p2ptypes.PeerID{peers[0]}, // peer 0 is in the capability DON
		F:       0,
	}
	workflowDonInfo := commoncap.DON{
		ID:      workflowDONID,
		Members: []p2ptypes.PeerID{peers[1]}, // peer 1 is in the workflow DON
		F:       0,
	}

	dispatcher := mocks.NewDispatcher(t)
	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: 1,
		MessageExpiry:           100 * time.Second,
		MaxBatchSize:            maxBatchSize,
		BatchCollectionPeriod:   time.Second,
	}
	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}
	underlying := &testTrigger{
		info:            capInfo,
		registrationsCh: make(chan commoncap.TriggerRegistrationRequest, 2),
		eventCh:         make(chan commoncap.TriggerResponse, 2),
	}
	publisher := remote.NewTriggerPublisher(capInfo.ID, "", dispatcher, lggr)
	require.NoError(t, publisher.SetConfig(config, underlying, capDonInfo, workflowDONs))
	require.NoError(t, publisher.Start(ctx))
	return underlying, publisher, dispatcher, peers
}

func newRegisterTriggerMessage(t *testing.T, callerDonID uint32, sender p2ptypes.PeerID) *remotetypes.MessageBody {
	// trigger registration event
	triggerRequest := commoncap.TriggerRegistrationRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}
	marshaled, err := pb.MarshalTriggerRegistrationRequest(triggerRequest)
	require.NoError(t, err)
	return &remotetypes.MessageBody{
		Sender:      sender[:],
		Method:      remotetypes.MethodRegisterTrigger,
		CallerDonId: callerDonID,
		Payload:     marshaled,
	}
}

type testTrigger struct {
	info            commoncap.CapabilityInfo
	registrationsCh chan commoncap.TriggerRegistrationRequest
	eventCh         chan commoncap.TriggerResponse
}

func (tr *testTrigger) Info(_ context.Context) (commoncap.CapabilityInfo, error) {
	return tr.info, nil
}

func (tr *testTrigger) RegisterTrigger(_ context.Context, request commoncap.TriggerRegistrationRequest) (<-chan commoncap.TriggerResponse, error) {
	tr.registrationsCh <- request
	return tr.eventCh, nil
}

func (tr *testTrigger) UnregisterTrigger(_ context.Context, request commoncap.TriggerRegistrationRequest) error {
	return nil
}

func TestTriggerPublisher_MultipleTriggersSameWorkflow(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.Test(t)
	capabilityDONID, workflowDONID := uint32(1), uint32(2)

	capInfo := commoncap.CapabilityInfo{
		ID:             capID,
		CapabilityType: commoncap.CapabilityTypeTrigger,
		Description:    "Remote Trigger",
	}
	peers := make([]p2ptypes.PeerID, 2)
	require.NoError(t, peers[0].UnmarshalText([]byte(peerID1)))
	require.NoError(t, peers[1].UnmarshalText([]byte(peerID2)))
	capDonInfo := commoncap.DON{
		ID:      capabilityDONID,
		Members: []p2ptypes.PeerID{peers[0]},
		F:       0,
	}
	workflowDonInfo := commoncap.DON{
		ID:      workflowDONID,
		Members: []p2ptypes.PeerID{peers[1]},
		F:       0,
	}
	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}

	// Create a trigger that tracks registrations by triggerID
	underlying := newMultiTrigger(capInfo)

	dispatcher := mocks.NewDispatcher(t)
	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: 1,
		MessageExpiry:           100 * time.Second,
		MaxBatchSize:            1, // no batching
		BatchCollectionPeriod:   time.Second,
	}

	publisher := remote.NewTriggerPublisher(capInfo.ID, "", dispatcher, lggr)
	require.NoError(t, publisher.SetConfig(config, underlying, capDonInfo, workflowDONs))
	require.NoError(t, publisher.Start(ctx))

	// Register trigger1
	regEvent1 := newRegisterTriggerMessageWithTriggerID(t, workflowDONID, peers[1], "trigger1")
	publisher.Receive(ctx, regEvent1)
	reg1 := <-underlying.registrationsCh
	require.Equal(t, "trigger1", reg1.TriggerID)
	require.Equal(t, workflowID1, reg1.Metadata.WorkflowID)

	// Register trigger2 for the same workflow
	regEvent2 := newRegisterTriggerMessageWithTriggerID(t, workflowDONID, peers[1], "trigger2")
	publisher.Receive(ctx, regEvent2)
	reg2 := <-underlying.registrationsCh
	require.Equal(t, "trigger2", reg2.TriggerID)
	require.Equal(t, workflowID1, reg2.Metadata.WorkflowID) // same workflowID

	trigger1EventReceived := make(chan struct{})
	trigger2EventReceived := make(chan struct{})

	dispatcher.On("Send", peers[1], mock.Anything).Run(func(args mock.Arguments) {
		msg := args.Get(1).(*remotetypes.MessageBody)
		require.Equal(t, capID, msg.CapabilityId)
		require.Equal(t, remotetypes.MethodTriggerEvent, msg.Method)
		metadata := msg.Metadata.(*remotetypes.MessageBody_TriggerEventMetadata)
		require.Len(t, metadata.TriggerEventMetadata.WorkflowIds, 1)
		require.Len(t, metadata.TriggerEventMetadata.TriggerIds, 1)
		triggerID := metadata.TriggerEventMetadata.TriggerIds[0]
		eventID := metadata.TriggerEventMetadata.TriggerEventId
		if triggerID == "trigger1" && eventID == "event1" {
			close(trigger1EventReceived)
		} else if triggerID == "trigger2" && eventID == "event2" {
			close(trigger2EventReceived)
		}
	}).Return(nil)

	// Send both events and expect them to be delivered separately
	underlying.SendEvent("trigger1", commoncap.TriggerResponse{
		Event: commoncap.TriggerEvent{ID: "event1"},
	})
	underlying.SendEvent("trigger2", commoncap.TriggerResponse{
		Event: commoncap.TriggerEvent{ID: "event2"},
	})

	<-trigger1EventReceived
	<-trigger2EventReceived

	require.NoError(t, publisher.Close())
}

func newRegisterTriggerMessageWithTriggerID(t *testing.T, callerDonID uint32, sender p2ptypes.PeerID, triggerID string) *remotetypes.MessageBody {
	triggerRequest := commoncap.TriggerRegistrationRequest{
		TriggerID: triggerID,
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}
	marshaled, err := pb.MarshalTriggerRegistrationRequest(triggerRequest)
	require.NoError(t, err)
	return &remotetypes.MessageBody{
		Sender:      sender[:],
		Method:      remotetypes.MethodRegisterTrigger,
		CallerDonId: callerDonID,
		Payload:     marshaled,
	}
}

// multiTrigger is a test trigger that supports multiple trigger registrations
// and can send events to specific triggers by triggerID
type multiTrigger struct {
	info            commoncap.CapabilityInfo
	registrationsCh chan commoncap.TriggerRegistrationRequest
	eventChans      map[string]chan commoncap.TriggerResponse
	mu              sync.Mutex
}

func newMultiTrigger(info commoncap.CapabilityInfo) *multiTrigger {
	return &multiTrigger{
		info:            info,
		registrationsCh: make(chan commoncap.TriggerRegistrationRequest, 10),
		eventChans:      make(map[string]chan commoncap.TriggerResponse),
	}
}

func (tr *multiTrigger) Info(_ context.Context) (commoncap.CapabilityInfo, error) {
	return tr.info, nil
}

func (tr *multiTrigger) RegisterTrigger(_ context.Context, request commoncap.TriggerRegistrationRequest) (<-chan commoncap.TriggerResponse, error) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	ch := make(chan commoncap.TriggerResponse, 10)
	tr.eventChans[request.TriggerID] = ch
	tr.registrationsCh <- request
	return ch, nil
}

func (tr *multiTrigger) UnregisterTrigger(_ context.Context, request commoncap.TriggerRegistrationRequest) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	if ch, ok := tr.eventChans[request.TriggerID]; ok {
		close(ch)
		delete(tr.eventChans, request.TriggerID)
	}
	return nil
}

func (tr *multiTrigger) SendEvent(triggerID string, event commoncap.TriggerResponse) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	if ch, ok := tr.eventChans[triggerID]; ok {
		ch <- event
	}
}
