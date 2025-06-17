package remote_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	remoteMocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	peerID1     = "12D3KooWF3dVeJ6YoT5HFnYhmwQWWMoEwVFzJQ5kKCMX3ZityxMC"
	peerID2     = "12D3KooWQsmok6aD8PZqt3RnJhQRrNzKHLficq7zYFRp7kZ1hHP8"
	workflowID1 = "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
)

var (
	triggerEvent1 = map[string]any{"event": "triggerEvent1"}
)

func TestTriggerSubscriber_RegisterAndReceive(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	capInfo, capDon, workflowDon := buildTwoTestDONs(t, 1, 1)
	dispatcher := remoteMocks.NewDispatcher(t)
	awaitRegistrationMessageCh := make(chan struct{})
	dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		select {
		case awaitRegistrationMessageCh <- struct{}{}:
		default:
		}
	})

	// register trigger
	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: 1,
		MessageExpiry:           100 * time.Second,
	}
	subscriber := remote.NewTriggerSubscriber(config, capInfo, capDon, workflowDon, dispatcher, nil, lggr)
	require.NoError(t, subscriber.Start(t.Context()))

	req := commoncap.TriggerRegistrationRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}
	triggerEventCallbackCh, err := subscriber.RegisterTrigger(t.Context(), req)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, subscriber.UnregisterTrigger(t.Context(), req))
		// calling UnregisterTrigger repeatedly is safe
		require.NoError(t, subscriber.UnregisterTrigger(t.Context(), req))
		require.NoError(t, subscriber.Close())
	})
	<-awaitRegistrationMessageCh

	// receive trigger event
	triggerEventValue, err := values.NewMap(triggerEvent1)
	require.NoError(t, err)
	triggerEvent := buildTriggerEvent(t, capDon.Members[0][:])
	subscriber.Receive(t.Context(), triggerEvent)
	response := <-triggerEventCallbackCh
	require.Equal(t, response.Event.Outputs, triggerEventValue)
}

func TestTriggerSubscriber_CorrectEventExpiryCheck(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	capInfo, capDon, workflowDon := buildTwoTestDONs(t, 3, 1)
	awaitRegistrationMessageCh := make(chan struct{})
	dispatcher := remoteMocks.NewDispatcher(t)
	dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		select {
		case awaitRegistrationMessageCh <- struct{}{}:
		default:
		}
	})

	// register trigger
	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      10 * time.Second,
		MinResponsesToAggregate: 2,
		MessageExpiry:           10 * time.Second,
	}
	subscriber := remote.NewTriggerSubscriber(config, capInfo, capDon, workflowDon, dispatcher, nil, lggr)

	require.NoError(t, subscriber.Start(t.Context()))
	regReq := commoncap.TriggerRegistrationRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}
	triggerEventCallbackCh, err := subscriber.RegisterTrigger(t.Context(), regReq)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, subscriber.UnregisterTrigger(t.Context(), regReq))
		require.NoError(t, subscriber.Close())
	})
	<-awaitRegistrationMessageCh

	// receive trigger events:
	// cleanup loop happens every 10 seconds, at 0:00, 0:10, 0:20, etc.
	// send the event from the first node around 0:02 (this is a bad node
	// that sends it too early)
	triggerEvent := buildTriggerEvent(t, capDon.Members[0][:])
	time.Sleep(2 * time.Second)
	subscriber.Receive(t.Context(), triggerEvent)

	// send events from nodes 2 & 3 (the good ones) around 0:15 so that
	// the diff between 0:02 and 0:15 exceeds the expiry threshold but
	// we don't hit the cleanup loop yet
	time.Sleep(13 * time.Second)
	triggerEvent.Sender = capDon.Members[1][:]
	subscriber.Receive(t.Context(), triggerEvent)
	// the aggregation shouldn't happen after events 1 and 2 as they
	// were received too far apart in time
	require.Empty(t, triggerEventCallbackCh)
	triggerEvent.Sender = capDon.Members[2][:]
	subscriber.Receive(t.Context(), triggerEvent)

	// event should be processed
	response := <-triggerEventCallbackCh
	triggerEventValue, err := values.NewMap(triggerEvent1)
	require.NoError(t, err)
	require.Equal(t, response.Event.Outputs, triggerEventValue)
}

func buildTwoTestDONs(t *testing.T, capDonSize int, workflowDonSize int) (commoncap.CapabilityInfo, commoncap.DON, commoncap.DON) {
	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1",
		CapabilityType: commoncap.CapabilityTypeTrigger,
		Description:    "Remote Trigger",
	}

	capDon := commoncap.DON{
		ID:      1,
		Members: []p2ptypes.PeerID{},
		F:       0,
	}
	for range capDonSize {
		pid := utils.MustNewPeerID()
		peer := p2ptypes.PeerID{}
		require.NoError(t, peer.UnmarshalText([]byte(pid)))
		capDon.Members = append(capDon.Members, peer)
	}

	workflowDon := commoncap.DON{
		ID:      2,
		Members: []p2ptypes.PeerID{},
		F:       0,
	}
	for range workflowDonSize {
		pid := utils.MustNewPeerID()
		peer := p2ptypes.PeerID{}
		require.NoError(t, peer.UnmarshalText([]byte(pid)))
		workflowDon.Members = append(workflowDon.Members, peer)
	}
	return capInfo, capDon, workflowDon
}

func buildTriggerEvent(t *testing.T, sender []byte) *remotetypes.MessageBody {
	triggerEventValue, err := values.NewMap(triggerEvent1)
	require.NoError(t, err)
	capResponse := commoncap.TriggerResponse{
		Event: commoncap.TriggerEvent{
			Outputs: triggerEventValue,
		},
		Err: nil,
	}
	marshaled, err := pb.MarshalTriggerResponse(capResponse)
	require.NoError(t, err)

	return &remotetypes.MessageBody{
		Sender: sender,
		Method: remotetypes.MethodTriggerEvent,
		Metadata: &remotetypes.MessageBody_TriggerEventMetadata{
			TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
				WorkflowIds: []string{workflowID1},
			},
		},
		Payload: marshaled,
	}
}
