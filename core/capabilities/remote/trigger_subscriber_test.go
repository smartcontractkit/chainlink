package remote_test

import (
	"testing"

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
	peerID1       = "12D3KooWF3dVeJ6YoT5HFnYhmwQWWMoEwVFzJQ5kKCMX3ZityxMC"
	peerID2       = "12D3KooWQsmok6aD8PZqt3RnJhQRrNzKHLficq7zYFRp7kZ1hHP8"
	workflowID1   = "workflowID1"
	triggerEvent1 = "triggerEvent1"
	triggerEvent2 = "triggerEvent2"
)

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
		ID:      "capability-don",
		Members: []p2ptypes.PeerID{p1},
		F:       0,
	}
	workflowDonInfo := commoncap.DON{
		ID:      "workflow-don",
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
	config := remotetypes.RemoteTriggerConfig{
		RegistrationRefreshMs:   100,
		RegistrationExpiryMs:    100,
		MinResponsesToAggregate: 1,
		MessageExpiryMs:         100_000,
	}
	subscriber := remote.NewTriggerSubscriber(config, capInfo, capDonInfo, workflowDonInfo, dispatcher, nil, lggr)
	require.NoError(t, subscriber.Start(ctx))

	req := commoncap.CapabilityRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}
	triggerEventCallbackCh, err := subscriber.RegisterTrigger(ctx, req)
	require.NoError(t, err)
	<-awaitRegistrationMessageCh

	// receive trigger event
	triggerEventValue, err := values.Wrap(triggerEvent1)
	require.NoError(t, err)
	capResponse := commoncap.CapabilityResponse{
		Value: triggerEventValue,
		Err:   nil,
	}
	marshaled, err := pb.MarshalCapabilityResponse(capResponse)
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
	subscriber.Receive(triggerEvent)
	response := <-triggerEventCallbackCh
	require.Equal(t, response.Value, triggerEventValue)

	require.NoError(t, subscriber.UnregisterTrigger(ctx, req))
	require.NoError(t, subscriber.UnregisterTrigger(ctx, req))
	require.NoError(t, subscriber.Close())
}
