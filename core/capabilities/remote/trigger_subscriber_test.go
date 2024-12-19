package remote_test

import (
	"strconv"
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
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const (
	peerID1     = "12D3KooWF3dVeJ6YoT5HFnYhmwQWWMoEwVFzJQ5kKCMX3ZityxMC"
	peerID2     = "12D3KooWQsmok6aD8PZqt3RnJhQRrNzKHLficq7zYFRp7kZ1hHP8"
	workflowID1 = "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	workflowID2 = "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce2"
)

var (
	triggerEvent1 = map[string]any{"event": "triggerEvent1"}
)

func TestTriggerSubscriber_SlowConsumer(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1",
		CapabilityType: commoncap.CapabilityTypeTrigger,
		Description:    "Remote Trigger",
	}
	p1 := p2ptypes.PeerID{}
	require.NoError(t, p1.UnmarshalText([]byte(peerID1)))
	capDonInfo := commoncap.DON{
		ID:      1,
		Members: []p2ptypes.PeerID{p1},
		F:       0,
	}
	workflowDonInfo := commoncap.DON{
		ID:      2,
		Members: []p2ptypes.PeerID{p1},
		F:       0,
	}
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
	subscriber := remote.NewTriggerSubscriber(config, capInfo, capDonInfo, workflowDonInfo, dispatcher, nil, lggr)
	require.NoError(t, subscriber.Start(ctx))

	req := commoncap.TriggerRegistrationRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}
	triggerEventCallbackCh, err := subscriber.RegisterTrigger(ctx, req)
	require.NoError(t, err)
	<-awaitRegistrationMessageCh

	req2 := commoncap.TriggerRegistrationRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID2,
		},
	}
	triggerEventCallbackCh2, err := subscriber.RegisterTrigger(ctx, req2)
	require.NoError(t, err)
	<-awaitRegistrationMessageCh

	// receive trigger event
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

	// Simulate a slow consumer by not pulling events from triggerEventCallbackCh after a few events
	// triggerEventCallbackCh2 should still receive events, note the number of messages

	for i := 0; i < remote.DefaultSendChannelBufferSize*2; i++ {
		triggerEvent := &remotetypes.MessageBody{

			Sender: p1[:],
			Method: remotetypes.MethodTriggerEvent,
			Metadata: &remotetypes.MessageBody_TriggerEventMetadata{
				TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
					TriggerEventId: strconv.Itoa(i),
					WorkflowIds:    []string{workflowID1, workflowID2},
				},
			},
			Payload: marshaled,
		}

		subscriber.Receive(ctx, triggerEvent)
		if i < 5 {
			select {
			case <-triggerEventCallbackCh:
			case <-triggerEventCallbackCh2:
			}
		} else {
			<-triggerEventCallbackCh2
		}
	}

	// Confirm the slow consumer channel was closed
	for {
		_, ok := <-triggerEventCallbackCh
		if !ok {
			break
		}
	}

	require.NoError(t, subscriber.UnregisterTrigger(ctx, req))
	require.NoError(t, subscriber.UnregisterTrigger(ctx, req2))
	require.NoError(t, subscriber.Close())
}

func TestTriggerSubscriber_RegisterAndReceive(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1",
		CapabilityType: commoncap.CapabilityTypeTrigger,
		Description:    "Remote Trigger",
	}
	p1 := p2ptypes.PeerID{}
	require.NoError(t, p1.UnmarshalText([]byte(peerID1)))
	p2 := p2ptypes.PeerID{}
	require.NoError(t, p2.UnmarshalText([]byte(peerID2)))
	capDonInfo := commoncap.DON{
		ID:      1,
		Members: []p2ptypes.PeerID{p1},
		F:       0,
	}
	workflowDonInfo := commoncap.DON{
		ID:      2,
		Members: []p2ptypes.PeerID{p2},
		F:       0,
	}
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
	subscriber := remote.NewTriggerSubscriber(config, capInfo, capDonInfo, workflowDonInfo, dispatcher, nil, lggr)
	require.NoError(t, subscriber.Start(ctx))

	req := commoncap.TriggerRegistrationRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}
	triggerEventCallbackCh, err := subscriber.RegisterTrigger(ctx, req)
	require.NoError(t, err)
	<-awaitRegistrationMessageCh

	// receive trigger event
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
	triggerEvent := &remotetypes.MessageBody{
		Sender: p1[:],
		Method: remotetypes.MethodTriggerEvent,
		Metadata: &remotetypes.MessageBody_TriggerEventMetadata{
			TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
				WorkflowIds: []string{workflowID1},
			},
		},
		Payload: marshaled,
	}
	subscriber.Receive(ctx, triggerEvent)
	response := <-triggerEventCallbackCh
	require.Equal(t, response.Event.Outputs, triggerEventValue)

	require.NoError(t, subscriber.UnregisterTrigger(ctx, req))
	require.NoError(t, subscriber.UnregisterTrigger(ctx, req))
	require.NoError(t, subscriber.Close())
}
