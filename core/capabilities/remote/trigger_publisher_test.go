package remote_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
	underlyingTriggerCap.eventCh <- commoncap.TriggerResponse{}
	awaitOutgoingMessageCh := make(chan struct{})
	dispatcher.On("Send", peers[1], mock.Anything).Run(func(args mock.Arguments) {
		awaitOutgoingMessageCh <- struct{}{}
	}).Return(nil)
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
	underlyingTriggerCap.eventCh <- commoncap.TriggerResponse{}
	underlyingTriggerCap.eventCh <- commoncap.TriggerResponse{}
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
	<-awaitOutgoingMessageCh

	// if there are fewer pending event than the batch size,
	// the events should still be sent after the batch collection period
	underlyingTriggerCap.eventCh <- commoncap.TriggerResponse{}
	dispatcher.On("Send", peers[1], mock.Anything).Run(func(args mock.Arguments) {
		msg := args.Get(1).(*remotetypes.MessageBody)
		metadata := msg.Metadata.(*remotetypes.MessageBody_TriggerEventMetadata)
		require.Len(t, metadata.TriggerEventMetadata.WorkflowIds, 1)
		awaitOutgoingMessageCh <- struct{}{}
	}).Return(nil).Once()
	<-awaitOutgoingMessageCh

	require.NoError(t, publisher.Close())
}

func newServices(t *testing.T, capabilityDONID uint32, workflowDONID uint32, maxBatchSize uint32) (*testTrigger, remotetypes.ReceiverService, *mocks.Dispatcher, []p2ptypes.PeerID) {
	lggr := logger.TestLogger(t)
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
	publisher := remote.NewTriggerPublisher(config, underlying, capInfo, capDonInfo, workflowDONs, dispatcher, lggr)
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
