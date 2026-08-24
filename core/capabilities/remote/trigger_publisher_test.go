package remote_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/synctest"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const capID = "cap_id@1"

func TestTriggerPublisher_Register(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	capabilityDONID, workflowDONID := uint32(1), uint32(2)

	underlyingTriggerCap, publisher, _, peers := newServices(t, capabilityDONID, workflowDONID, 1, time.Second)

	// invalid sender case - node 0 is not a member of the workflow DON, registration shouldn't happen
	regEvent := newRegisterTriggerMessage(t, workflowDONID, peers[0])
	publisher.Receive(ctx, regEvent)
	require.Empty(t, underlyingTriggerCap.registrationsCh)

	// valid registration
	regEvent = newRegisterTriggerMessage(t, workflowDONID, peers[1])
	publisher.Receive(ctx, regEvent)
	require.Eventually(t, func() bool { return len(underlyingTriggerCap.registrationsCh) > 0 }, 2*time.Second, 10*time.Millisecond)
	forwarded := <-underlyingTriggerCap.registrationsCh
	require.Equal(t, workflowID1, forwarded.Metadata.WorkflowID)
}

func TestTriggerPublisher_ReceiveTriggerEvents_NoBatching(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	capabilityDONID, workflowDONID := uint32(1), uint32(2)

	underlyingTriggerCap, publisher, dispatcher, peers := newServices(t, capabilityDONID, workflowDONID, 1, time.Second)
	regEvent := newRegisterTriggerMessage(t, workflowDONID, peers[1])
	publisher.Receive(ctx, regEvent)
	require.Eventually(t, func() bool { return len(underlyingTriggerCap.registrationsCh) > 0 }, 2*time.Second, 10*time.Millisecond)

	// send a trigger event and expect that it gets delivered right away
	awaitOutgoingMessageCh := make(chan struct{})
	dispatcher.On("Send", peers[1], mock.Anything).Run(func(args mock.Arguments) {
		awaitOutgoingMessageCh <- struct{}{}
	}).Return(nil)
	underlyingTriggerCap.eventCh <- commoncap.TriggerResponse{}
	<-awaitOutgoingMessageCh
}

func TestTriggerPublisher_ReceiveTriggerEvents_BatchingEnabled(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		const batchPeriod = 50 * time.Millisecond
		ctx := t.Context()
		capabilityDONID, workflowDONID := uint32(1), uint32(2)

		underlyingTriggerCap, publisher, dispatcher, peers := newServices(t, capabilityDONID, workflowDONID, 2, batchPeriod)
		regEvent := newRegisterTriggerMessage(t, workflowDONID, peers[1])
		publisher.Receive(ctx, regEvent)
		require.Eventually(t, func() bool { return len(underlyingTriggerCap.registrationsCh) > 0 }, 2*time.Second, 10*time.Millisecond)

		awaitOutgoingMessageCh := make(chan struct{}, 1)

		// send two trigger events and expect them to be delivered in a batch
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
		time.Sleep(batchPeriod)
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
		time.Sleep(batchPeriod)
		<-awaitOutgoingMessageCh
	})
}

func TestTriggerPublisher_ReceiveTriggerEventAcks(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	capabilityDONID, workflowDONID := uint32(1), uint32(2)
	underlyingTriggerCap, publisher, _, peers := newServices(t, capabilityDONID, workflowDONID, 2, time.Second)
	eventID := "123"
	triggerID := "trigA"
	regEvent := newAckEventMessage(t, eventID, triggerID, workflowDONID, peers[1])
	publisher.Receive(ctx, regEvent)

	require.Eventually(t, func() bool { return underlyingTriggerCap.ackCount.Load() == 1 }, 2*time.Second, 10*time.Millisecond)
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
		t.Parallel()

		dispatcher := mocks.NewDispatcher(t)
		publisher := remote.NewTriggerPublisher(capInfo.ID, "method", dispatcher, lggr)
		config := &commoncap.RemoteTriggerConfig{}
		err := publisher.SetConfig(config, nil, capDonInfo, workflowDONs)
		require.Error(t, err)
		require.Contains(t, err.Error(), "underlying trigger capability cannot be nil")
	})

	t.Run("handles nil config", func(t *testing.T) {
		t.Parallel()

		dispatcher := mocks.NewDispatcher(t)
		publisher := remote.NewTriggerPublisher(capInfo.ID, "method", dispatcher, lggr)
		// Set config as nil - should use defaults
		err := publisher.SetConfig(nil, underlying, capDonInfo, workflowDONs)
		require.NoError(t, err)

		// Verify config works
		ctx := t.Context()
		require.NoError(t, publisher.Start(ctx))
		require.NoError(t, publisher.Close())
	})

	t.Run("handles nil workflowDONs", func(t *testing.T) {
		t.Parallel()

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
		ctx := t.Context()
		require.NoError(t, publisher.Start(ctx))
		require.NoError(t, publisher.Close())
	})

	t.Run("updates existing config", func(t *testing.T) {
		t.Parallel()

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
		ctx := t.Context()
		require.NoError(t, publisher.Start(ctx))
		require.NoError(t, publisher.Close())
	})
}

func newServices(t *testing.T, capabilityDONID, workflowDONID, maxBatchSize uint32, batchCollectionPeriod time.Duration) (*testTrigger, remotetypes.ReceiverService, *mocks.Dispatcher, []p2ptypes.PeerID) {
	t.Helper()
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
	allowRegistrationChecks(dispatcher)
	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: 1,
		MessageExpiry:           100 * time.Second,
		MaxBatchSize:            maxBatchSize,
		BatchCollectionPeriod:   batchCollectionPeriod,
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
	servicetest.Run(t, publisher)
	return underlying, publisher, dispatcher, peers
}

func allowRegistrationChecks(dispatcher *mocks.Dispatcher) {
	dispatcher.On("Send", mock.Anything, mock.MatchedBy(func(m *remotetypes.MessageBody) bool {
		return m.Method == remotetypes.MethodTriggerRegistrationCheck
	})).Return(nil).Maybe()
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

func newAckEventMessage(t *testing.T, eventID, triggerID string, callerDonID uint32, sender p2ptypes.PeerID) *remotetypes.MessageBody {
	return &remotetypes.MessageBody{
		Sender:      sender[:],
		Method:      remotetypes.MethodTriggerEventAck,
		CallerDonId: callerDonID,
		Metadata: &remotetypes.MessageBody_TriggerEventMetadata{
			TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
				TriggerEventId: eventID,
				TriggerIds:     []string{triggerID},
			},
		},
	}
}

type testTrigger struct {
	info            commoncap.CapabilityInfo
	registrationsCh chan commoncap.TriggerRegistrationRequest
	eventCh         chan commoncap.TriggerResponse
	eventAckd       bool
	ackCount        atomic.Int32
	ackBlock        chan struct{}
}

func (tr *testTrigger) Info(_ context.Context) (commoncap.CapabilityInfo, error) {
	return tr.info, nil
}

// waitForRegistration blocks until async RegisterTrigger delivers to the underlying mock.
func waitForRegistration(t *testing.T, ch <-chan commoncap.TriggerRegistrationRequest) commoncap.TriggerRegistrationRequest {
	t.Helper()
	require.Eventually(t, func() bool { return len(ch) > 0 }, 2*time.Second, 10*time.Millisecond)
	return <-ch
}

// waitForUnregister blocks until UnregisterTrigger runs on the underlying mock.
func waitForUnregister(t *testing.T, ch <-chan string) string {
	t.Helper()
	require.Eventually(t, func() bool { return len(ch) > 0 }, 2*time.Second, 10*time.Millisecond)
	return <-ch
}

func (tr *testTrigger) RegisterTrigger(_ context.Context, request commoncap.TriggerRegistrationRequest) (<-chan commoncap.TriggerResponse, error) {
	tr.registrationsCh <- request
	return tr.eventCh, nil
}

func (tr *testTrigger) UnregisterTrigger(_ context.Context, request commoncap.TriggerRegistrationRequest) error {
	return nil
}

func (tr *testTrigger) AckEvent(_ context.Context, triggerID, eventID, method string) error {
	tr.eventAckd = true
	tr.ackCount.Add(1)
	if tr.ackBlock != nil {
		<-tr.ackBlock
	}
	return nil
}

func TestTriggerPublisher_AckEvent_DoesNotBlockReceive(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	capabilityDONID, workflowDONID := uint32(1), uint32(2)
	underlying, publisher, _, peers := newServices(t, capabilityDONID, workflowDONID, 1, time.Second)
	underlying.ackBlock = make(chan struct{})

	eventID := "evt-block-test"
	triggerID := "trigA"
	publisher.Receive(ctx, newAckEventMessage(t, eventID, triggerID, workflowDONID, peers[1]))

	regMsg := newRegisterTriggerMessage(t, workflowDONID, peers[1])
	regDone := make(chan struct{})
	go func() {
		publisher.Receive(ctx, regMsg)
		close(regDone)
	}()

	select {
	case <-regDone:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("RegisterTrigger Receive blocked while AckEvent in flight")
	}

	close(underlying.ackBlock)
	require.Eventually(t, func() bool { return underlying.ackCount.Load() == 1 }, 2*time.Second, 10*time.Millisecond)
}

func TestTriggerPublisher_MultipleTriggersSameWorkflow(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
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
	allowRegistrationChecks(dispatcher)
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
	reg1 := waitForRegistration(t, underlying.registrationsCh)
	require.Equal(t, "trigger1", reg1.TriggerID)
	require.Equal(t, workflowID1, reg1.Metadata.WorkflowID)

	// Register trigger2 for the same workflow
	regEvent2 := newRegisterTriggerMessageWithTriggerID(t, workflowDONID, peers[1], "trigger2")
	publisher.Receive(ctx, regEvent2)
	reg2 := waitForRegistration(t, underlying.registrationsCh)
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

func TestTriggerPublisher_ExplicitUnregister(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
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

	dispatcher := mocks.NewDispatcher(t)
	allowRegistrationChecks(dispatcher)

	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: 1,
		MessageExpiry:           100 * time.Second,
		MaxBatchSize:            1,
		BatchCollectionPeriod:   time.Second,
	}

	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}

	underlying := newMultiTrigger(capInfo)

	publisher := remote.NewTriggerPublisher(capInfo.ID, "", dispatcher, lggr)
	require.NoError(t, publisher.SetConfig(config, underlying, capDonInfo, workflowDONs))
	servicetest.Run(t, publisher)

	// Register trigger
	regEvent := newRegisterTriggerMessageWithTriggerID(t, workflowDONID, peers[1], "triggerA")
	publisher.Receive(ctx, regEvent)

	waitForRegistration(t, underlying.registrationsCh)

	// Send unregister
	unregMsg := &remotetypes.MessageBody{
		Sender:      peers[1][:],
		Method:      remotetypes.MethodUnregisterTrigger,
		CallerDonId: workflowDONID,
		Metadata: &remotetypes.MessageBody_TriggerEventMetadata{
			TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
				WorkflowIds: []string{workflowID1},
				TriggerIds:  []string{"triggerA"},
			},
		},
	}
	publisher.Receive(ctx, unregMsg)
	require.Equal(t, "triggerA", waitForUnregister(t, underlying.unregisterCalled))
}

func TestTriggerPublisher_SendsRegistrationChecks(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
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

	underlying := newMultiTrigger(capInfo)
	dispatcher := mocks.NewDispatcher(t)

	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: 1,
		MessageExpiry:           100 * time.Second,
		MaxBatchSize:            1,
		BatchCollectionPeriod:   time.Second,
	}

	publisher := remote.NewTriggerPublisher(capInfo.ID, "", dispatcher, lggr)
	require.NoError(t, publisher.SetConfig(config, underlying, capDonInfo, workflowDONs))

	checkReceived := make(chan *remotetypes.MessageBody, 10)
	dispatcher.On("Send", mock.Anything, mock.MatchedBy(func(m *remotetypes.MessageBody) bool {
		return m.Method == remotetypes.MethodTriggerRegistrationCheck
	})).Run(func(args mock.Arguments) {
		msg := args.Get(1).(*remotetypes.MessageBody)
		checkReceived <- msg
	}).Return(nil).Maybe()

	// Start before Receive so initMetrics() has run (RegisterTrigger success path records metrics).
	servicetest.Run(t, publisher)

	regEvent := newRegisterTriggerMessageWithTriggerID(t, workflowDONID, peers[1], "triggerA")
	publisher.Receive(ctx, regEvent)
	waitForRegistration(t, underlying.registrationsCh)

	select {
	case msg := <-checkReceived:
		meta := msg.GetTriggerEventMetadata()
		require.NotNil(t, meta)
		require.Equal(t, []string{workflowID1}, meta.WorkflowIds)
		require.Equal(t, []string{"triggerA"}, meta.TriggerIds)
		require.Equal(t, capabilityDONID, msg.CapabilityDonId)
		require.Equal(t, workflowDONID, msg.CallerDonId)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for registration check message")
	}
}

func TestTriggerPublisher_RegistrationChecksChunkByMaxBatchSize(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Parallel()
	ctx := t.Context()
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
	capDonInfo := commoncap.DON{ID: capabilityDONID, Members: []p2ptypes.PeerID{peers[0]}, F: 0}
	workflowDonInfo := commoncap.DON{ID: workflowDONID, Members: []p2ptypes.PeerID{peers[1]}, F: 0}
	workflowDONs := map[uint32]commoncap.DON{workflowDonInfo.ID: workflowDonInfo}

	const maxBatchSize uint32 = 100
	const nRegs = 250

	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: 1,
		MessageExpiry:           100 * time.Second,
		MaxBatchSize:            maxBatchSize,
		BatchCollectionPeriod:   time.Second,
	}

	underlying := newMultiTrigger(capInfo)
	dispatcher := mocks.NewDispatcher(t)

	var mu sync.Mutex
	var chunkLens []int
	dispatcher.On("Send", peers[1], mock.MatchedBy(func(m *remotetypes.MessageBody) bool {
		return m.Method == remotetypes.MethodTriggerRegistrationCheck
	})).Run(func(args mock.Arguments) {
		msg := args.Get(1).(*remotetypes.MessageBody)
		meta := msg.GetTriggerEventMetadata()
		require.NotNil(t, meta)
		mu.Lock()
		chunkLens = append(chunkLens, len(meta.WorkflowIds))
		mu.Unlock()
	}).Return(nil).Maybe()

	publisher := remote.NewTriggerPublisher(capInfo.ID, "", dispatcher, lggr)
	require.NoError(t, publisher.SetConfig(config, underlying, capDonInfo, workflowDONs))
	servicetest.Run(t, publisher)

	for i := range nRegs {
		publisher.Receive(ctx, newRegisterTriggerMessageWithTriggerID(t, workflowDONID, peers[1], fmt.Sprintf("trigger_%d", i)))
		waitForRegistration(t, underlying.registrationsCh)
	}

	// 250 registrations at MaxBatchSize=100 → chunk lengths 100, 100, 50 per tick per peer.
	require.Eventually(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		var has100, has50 bool
		var n100 int
		for _, n := range chunkLens {
			if n == 100 {
				has100 = true
				n100++
			}
			if n == 50 {
				has50 = true
			}
		}
		return has100 && has50 && n100 >= 2 && len(chunkLens) >= 3
	}, 3*time.Second, 20*time.Millisecond)
}

func TestTriggerPublisher_UnregisterValidatesSenderMembership(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	lggr := logger.Test(t)

	capabilityDONID, workflowDONID := uint32(1), uint32(2)

	capInfo := commoncap.CapabilityInfo{
		ID:             capID,
		CapabilityType: commoncap.CapabilityTypeTrigger,
		Description:    "Remote Trigger",
	}

	peers := make([]p2ptypes.PeerID, 3)
	require.NoError(t, peers[0].UnmarshalText([]byte(peerID1)))
	require.NoError(t, peers[1].UnmarshalText([]byte(peerID2)))
	// peers[2] is a random peer not in any DON
	peers[2] = p2ptypes.PeerID{0xff}

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

	underlying := newMultiTrigger(capInfo)
	dispatcher := mocks.NewDispatcher(t)
	allowRegistrationChecks(dispatcher)

	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: 1,
		MessageExpiry:           100 * time.Second,
		MaxBatchSize:            1,
		BatchCollectionPeriod:   time.Second,
	}

	publisher := remote.NewTriggerPublisher(capInfo.ID, "", dispatcher, lggr)
	require.NoError(t, publisher.SetConfig(config, underlying, capDonInfo, workflowDONs))
	servicetest.Run(t, publisher)

	// Register a trigger from the valid workflow DON member
	regEvent := newRegisterTriggerMessageWithTriggerID(t, workflowDONID, peers[1], "triggerA")
	publisher.Receive(ctx, regEvent)
	waitForRegistration(t, underlying.registrationsCh)

	// Send unregister from a peer NOT in the workflow DON — should be ignored
	unregMsg := &remotetypes.MessageBody{
		Sender:      peers[2][:],
		Method:      remotetypes.MethodUnregisterTrigger,
		CallerDonId: workflowDONID,
		Metadata: &remotetypes.MessageBody_TriggerEventMetadata{
			TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
				WorkflowIds: []string{workflowID1},
				TriggerIds:  []string{"triggerA"},
			},
		},
	}
	publisher.Receive(ctx, unregMsg)

	// UnregisterTrigger should NOT have been called on the underlying
	select {
	case trigID := <-underlying.unregisterCalled:
		t.Fatalf("expected no unregister, but got unregister for %s", trigID)
	default:
		// expected: no unregister
	}

	// Now send from the valid member — should succeed
	unregMsg.Sender = peers[1][:]
	publisher.Receive(ctx, unregMsg)
	require.Equal(t, "triggerA", waitForUnregister(t, underlying.unregisterCalled))
}

func TestTriggerPublisher_UnregisterRequiresQuorum(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	lggr := logger.Test(t)

	capabilityDONID, workflowDONID := uint32(1), uint32(2)

	capInfo := commoncap.CapabilityInfo{
		ID:             capID,
		CapabilityType: commoncap.CapabilityTypeTrigger,
		Description:    "Remote Trigger",
	}

	peers := make([]p2ptypes.PeerID, 4)
	require.NoError(t, peers[0].UnmarshalText([]byte(peerID1)))
	require.NoError(t, peers[1].UnmarshalText([]byte(peerID2)))
	require.NoError(t, peers[2].UnmarshalText([]byte(peerID3)))
	require.NoError(t, peers[3].UnmarshalText([]byte("12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw")))

	capDonInfo := commoncap.DON{
		ID:      capabilityDONID,
		Members: []p2ptypes.PeerID{peers[0]},
		F:       0,
	}
	workflowDonInfo := commoncap.DON{
		ID:      workflowDONID,
		Members: []p2ptypes.PeerID{peers[1], peers[2], peers[3]},
		F:       1,
	}
	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}

	dispatcher := mocks.NewDispatcher(t)
	allowRegistrationChecks(dispatcher)

	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: 1,
		MessageExpiry:           100 * time.Second,
		MaxBatchSize:            1,
		BatchCollectionPeriod:   time.Second,
	}

	underlying := newMultiTrigger(capInfo)
	publisher := remote.NewTriggerPublisher(capInfo.ID, "", dispatcher, lggr)
	require.NoError(t, publisher.SetConfig(config, underlying, capDonInfo, workflowDONs))
	servicetest.Run(t, publisher)

	for _, p := range []p2ptypes.PeerID{peers[1], peers[2], peers[3]} {
		publisher.Receive(ctx, newRegisterTriggerMessageWithTriggerID(t, workflowDONID, p, "triggerA"))
	}
	require.Equal(t, "triggerA", waitForRegistration(t, underlying.registrationsCh).TriggerID)

	unregMeta := &remotetypes.MessageBody_TriggerEventMetadata{
		TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
			WorkflowIds: []string{workflowID1},
			TriggerIds:  []string{"triggerA"},
		},
	}
	recvUnreg := func(sender p2ptypes.PeerID) {
		publisher.Receive(ctx, &remotetypes.MessageBody{
			Sender:      sender[:],
			Method:      remotetypes.MethodUnregisterTrigger,
			CallerDonId: workflowDONID,
			Metadata:    unregMeta,
		})
	}

	recvUnreg(peers[1])
	recvUnreg(peers[2])
	select {
	case id := <-underlying.unregisterCalled:
		t.Fatalf("unregister after 2 of 3 nodes should not run; got %q", id)
	default:
	}

	recvUnreg(peers[3])
	require.Equal(t, "triggerA", waitForUnregister(t, underlying.unregisterCalled))
}

func TestTriggerPublisher_UnregisterInvalidMetadata(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	_, publisher, _, peers := newServices(t, 1, 2, 1, time.Second)

	cases := []struct {
		name string
		meta *remotetypes.MessageBody_TriggerEventMetadata
	}{
		{
			name: "nil metadata",
			meta: nil,
		},
		{
			name: "empty workflow IDs",
			meta: &remotetypes.MessageBody_TriggerEventMetadata{
				TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
					WorkflowIds: []string{},
					TriggerIds:  []string{"triggerA"},
				},
			},
		},
		{
			name: "empty trigger IDs",
			meta: &remotetypes.MessageBody_TriggerEventMetadata{
				TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
					WorkflowIds: []string{workflowID1},
					TriggerIds:  []string{},
				},
			},
		},
		{
			name: "multiple workflow IDs",
			meta: &remotetypes.MessageBody_TriggerEventMetadata{
				TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
					WorkflowIds: []string{workflowID1, workflowID1},
					TriggerIds:  []string{"triggerA"},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			msg := &remotetypes.MessageBody{
				Sender:      peers[1][:],
				Method:      remotetypes.MethodUnregisterTrigger,
				CallerDonId: 2,
			}
			if tc.meta != nil {
				msg.Metadata = tc.meta
			}
			// Should not panic
			publisher.Receive(ctx, msg)
		})
	}
}

func TestTriggerPublisher_AckCacheCleanup(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Parallel()
	ctx := t.Context()
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
		Members: []p2ptypes.PeerID{peers[0], peers[1]},
		F:       0,
	}
	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}

	underlying := newMultiTrigger(capInfo)
	dispatcher := mocks.NewDispatcher(t)

	triggerEventSent := make(chan p2ptypes.PeerID, 10)
	dispatcher.On("Send", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		msg := args.Get(1).(*remotetypes.MessageBody)
		if msg.Method == remotetypes.MethodTriggerEvent {
			triggerEventSent <- args.Get(0).(p2ptypes.PeerID)
		}
	}).Return(nil).Maybe()

	// Very short MessageExpiry so ACK entries expire quickly
	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     200 * time.Millisecond,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: 1,
		MessageExpiry:           200 * time.Millisecond,
		MaxBatchSize:            1,
		BatchCollectionPeriod:   time.Second,
	}

	publisher := remote.NewTriggerPublisher(capInfo.ID, "", dispatcher, lggr)
	require.NoError(t, publisher.SetConfig(config, underlying, capDonInfo, workflowDONs))

	// Start before Receive so initMetrics() has run on registration and ACK paths.
	servicetest.Run(t, publisher)

	regEvent := newRegisterTriggerMessageWithTriggerID(t, workflowDONID, peers[0], "triggerA")
	publisher.Receive(ctx, regEvent)
	waitForRegistration(t, underlying.registrationsCh)

	ackMsg := newAckEventMessage(t, "event1", "triggerA", workflowDONID, peers[0])
	publisher.Receive(ctx, ackMsg)

	// cacheCleanupLoop ticks on MessageExpiry (200ms) and removes expired ack cache
	// entries so a later trigger event is not suppressed.

	// Wait long enough for the ack cache entry to expire and be cleaned up
	time.Sleep(500 * time.Millisecond)

	// Send a new trigger event for the same event ID and verify it gets sent to
	// peers[0] again (not suppressed by the old ACK), proving the cleanup worked.
	underlying.SendEvent("triggerA", commoncap.TriggerResponse{
		Event: commoncap.TriggerEvent{ID: "event1"},
	})

	sentTo := make(map[p2ptypes.PeerID]bool)
	for range 2 {
		select {
		case peer := <-triggerEventSent:
			sentTo[peer] = true
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for trigger event sends")
		}
	}
	require.True(t, sentTo[peers[0]], "event should be re-sent to peers[0] after ack cache cleanup")
	require.True(t, sentTo[peers[1]], "event should be sent to peers[1]")
}

func TestTriggerPublisher_SecondDeliveryAfterFullAck_ReachesAllPeers(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
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
		Members: []p2ptypes.PeerID{peers[0], peers[1]},
		F:       0,
	}
	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}

	underlying := newMultiTrigger(capInfo)
	dispatcher := mocks.NewDispatcher(t)
	allowRegistrationChecks(dispatcher)

	var triggerEventMu sync.Mutex
	triggerEventSendCount := 0
	dispatcher.On("Send", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		msg := args.Get(1).(*remotetypes.MessageBody)
		if msg.Method == remotetypes.MethodTriggerEvent {
			triggerEventMu.Lock()
			triggerEventSendCount++
			triggerEventMu.Unlock()
		}
	}).Return(nil).Maybe()

	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: 1,
		MessageExpiry:           100 * time.Second,
		MaxBatchSize:            1,
		BatchCollectionPeriod:   time.Second,
	}

	publisher := remote.NewTriggerPublisher(capInfo.ID, "", dispatcher, lggr)
	require.NoError(t, publisher.SetConfig(config, underlying, capDonInfo, workflowDONs))
	servicetest.Run(t, publisher)

	regEvent := newRegisterTriggerMessageWithTriggerID(t, workflowDONID, peers[0], "triggerA")
	publisher.Receive(ctx, regEvent)
	waitForRegistration(t, underlying.registrationsCh)

	eventID := "shared-event-second-delivery-regression"
	underlying.SendEvent("triggerA", commoncap.TriggerResponse{
		Event: commoncap.TriggerEvent{ID: eventID},
	})

	require.Eventually(t, func() bool {
		triggerEventMu.Lock()
		defer triggerEventMu.Unlock()
		return triggerEventSendCount >= 2
	}, 2*time.Second, 10*time.Millisecond, "first delivery should emit MethodTriggerEvent once per workflow peer")

	triggerEventMu.Lock()
	firstRoundTotal := triggerEventSendCount
	triggerEventMu.Unlock()
	require.Equal(t, 2, firstRoundTotal, "workflow DON has two members")

	publisher.Receive(ctx, newAckEventMessage(t, eventID, "triggerA", workflowDONID, peers[0]))
	publisher.Receive(ctx, newAckEventMessage(t, eventID, "triggerA", workflowDONID, peers[1]))

	underlying.SendEvent("triggerA", commoncap.TriggerResponse{
		Event: commoncap.TriggerEvent{ID: eventID},
	})

	require.Eventually(t, func() bool {
		triggerEventMu.Lock()
		defer triggerEventMu.Unlock()
		return triggerEventSendCount >= firstRoundTotal+2
	}, 2*time.Second, 10*time.Millisecond,
		"second delivery with same triggerEventID must still fan out to all workflow peers (regression: ack-cache peer skip used to send zero)")

	triggerEventMu.Lock()
	finalTotal := triggerEventSendCount
	triggerEventMu.Unlock()
	require.Equal(t, 4, finalTotal, "two delivery rounds × two workflow peers")
}

func TestTriggerPublisher_RegisterTrigger_FailureShortCircuit(t *testing.T) {
	t.Parallel()

	t.Run("user error suppresses retries", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		lggr := logger.Test(t)
		capabilityDONID, workflowDONID := uint32(1), uint32(2)

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

		capInfo := commoncap.CapabilityInfo{
			ID:             capID,
			CapabilityType: commoncap.CapabilityTypeTrigger,
		}
		userErr := caperrors.NewPublicUserError(errors.New("bad workflow config"), caperrors.InvalidArgument)
		underlying := &errTrigger{info: capInfo, err: userErr}

		dispatcher := mocks.NewDispatcher(t)
		config := &commoncap.RemoteTriggerConfig{
			RegistrationRefresh:     100 * time.Millisecond,
			RegistrationExpiry:      100 * time.Second,
			MinResponsesToAggregate: 1,
			MessageExpiry:           100 * time.Second,
			MaxBatchSize:            1,
			BatchCollectionPeriod:   time.Second,
		}
		publisher := remote.NewTriggerPublisher(capInfo.ID, "", dispatcher, lggr)
		require.NoError(t, publisher.SetConfig(config, underlying, capDonInfo, workflowDONs))
		servicetest.Run(t, publisher)

		regMsg := newRegisterTriggerMessage(t, workflowDONID, peers[1])
		publisher.Receive(ctx, regMsg)
		require.Eventually(t, func() bool { return underlying.callCount.Load() == 1 }, 2*time.Second, 10*time.Millisecond)

		publisher.Receive(ctx, regMsg)
		publisher.Receive(ctx, regMsg)
		require.Eventually(t, func() bool { return underlying.callCount.Load() == 1 }, 2*time.Second, 10*time.Millisecond)
	})

	t.Run("system error allows retries", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		lggr := logger.Test(t)
		capabilityDONID, workflowDONID := uint32(1), uint32(2)

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

		capInfo := commoncap.CapabilityInfo{
			ID:             capID,
			CapabilityType: commoncap.CapabilityTypeTrigger,
		}
		underlying := &errTrigger{info: capInfo, err: errors.New("transient system failure")}

		dispatcher := mocks.NewDispatcher(t)
		config := &commoncap.RemoteTriggerConfig{
			RegistrationRefresh:     100 * time.Millisecond,
			RegistrationExpiry:      100 * time.Second,
			MinResponsesToAggregate: 1,
			MessageExpiry:           100 * time.Second,
			MaxBatchSize:            1,
			BatchCollectionPeriod:   time.Second,
		}
		publisher := remote.NewTriggerPublisher(capInfo.ID, "", dispatcher, lggr)
		require.NoError(t, publisher.SetConfig(config, underlying, capDonInfo, workflowDONs))
		servicetest.Run(t, publisher)

		regMsg := newRegisterTriggerMessage(t, workflowDONID, peers[1])
		publisher.Receive(ctx, regMsg)
		require.Eventually(t, func() bool { return underlying.callCount.Load() == 1 }, 2*time.Second, 10*time.Millisecond)

		require.Eventually(t, func() bool {
			if underlying.callCount.Load() >= 2 {
				return true
			}
			publisher.Receive(ctx, regMsg)
			return underlying.callCount.Load() >= 2
		}, 2*time.Second, 10*time.Millisecond)

		require.Eventually(t, func() bool {
			if underlying.callCount.Load() >= 3 {
				return true
			}
			publisher.Receive(ctx, regMsg)
			return underlying.callCount.Load() >= 3
		}, 2*time.Second, 10*time.Millisecond)
	})
}

// errTrigger is a TriggerCapability that always returns an error from RegisterTrigger.
type errTrigger struct {
	info      commoncap.CapabilityInfo
	err       error
	callCount atomic.Int32
}

func (tr *errTrigger) Info(_ context.Context) (commoncap.CapabilityInfo, error) {
	return tr.info, nil
}

func (tr *errTrigger) RegisterTrigger(_ context.Context, _ commoncap.TriggerRegistrationRequest) (<-chan commoncap.TriggerResponse, error) {
	tr.callCount.Add(1)
	return nil, tr.err
}

func (tr *errTrigger) UnregisterTrigger(_ context.Context, _ commoncap.TriggerRegistrationRequest) error {
	return nil
}

func (tr *errTrigger) AckEvent(_ context.Context, _, _, _ string) error {
	return nil
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
	info             commoncap.CapabilityInfo
	registrationsCh  chan commoncap.TriggerRegistrationRequest
	eventChans       map[string]chan commoncap.TriggerResponse
	unregisterCalled chan string
	mu               sync.Mutex
}

func newMultiTrigger(info commoncap.CapabilityInfo) *multiTrigger {
	return &multiTrigger{
		info:             info,
		registrationsCh:  make(chan commoncap.TriggerRegistrationRequest, 128),
		eventChans:       make(map[string]chan commoncap.TriggerResponse),
		unregisterCalled: make(chan string, 1),
	}
}

func (tr *multiTrigger) Info(_ context.Context) (commoncap.CapabilityInfo, error) {
	return tr.info, nil
}

func (tr *multiTrigger) AckEvent(_ context.Context, triggerID, eventID, method string) error {
	return nil
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
	tr.unregisterCalled <- request.TriggerID
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
