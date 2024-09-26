package request_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/target"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/target/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const (
	workflowID1          = "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	workflowExecutionID1 = "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed"
)

func Test_ClientRequest_MessageValidation(t *testing.T) {
	lggr := logger.TestLogger(t)

	numCapabilityPeers := 2
	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityPeers[i] = NewP2PPeerID(t)
	}

	capDonInfo := commoncap.DON{
		ID:      1,
		Members: capabilityPeers,
		F:       1,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1.0.0",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Remote Target",
		DON:            &capDonInfo,
	}

	numWorkflowPeers := 2
	workflowPeers := make([]p2ptypes.PeerID, numWorkflowPeers)
	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeers[i] = NewP2PPeerID(t)
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      2,
	}

	executeInputs, err := values.NewMap(
		map[string]any{
			"executeValue1": "aValue1",
		},
	)
	require.NoError(t, err)

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_OneAtATime,
		"deltaStage": "1000ms",
	})
	require.NoError(t, err)

	capabilityRequest := commoncap.CapabilityRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID:          workflowID1,
			WorkflowExecutionID: workflowExecutionID1,
		},
		Inputs: executeInputs,
		Config: transmissionSchedule,
	}

	m, err := values.NewMap(map[string]any{"response": "response1"})
	require.NoError(t, err)
	capabilityResponse := commoncap.CapabilityResponse{
		Value: m,
	}

	rawResponse, err := pb.MarshalCapabilityResponse(capabilityResponse)
	require.NoError(t, err)

	messageID, err := target.GetMessageIDForRequest(capabilityRequest)
	require.NoError(t, err)

	msg := &types.MessageBody{
		CapabilityId:    capInfo.ID,
		CapabilityDonId: capDonInfo.ID,
		CallerDonId:     workflowDonInfo.ID,
		Method:          types.MethodExecute,
		Payload:         rawResponse,
		MessageId:       []byte("messageID"),
	}

	t.Run("Send second message with different response", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
		request, err := request.NewClientRequest(ctx, lggr, capabilityRequest, messageID, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute)
		defer request.Cancel(errors.New("test end"))

		require.NoError(t, err)

		nm, err := values.NewMap(map[string]any{"response": "response2"})
		require.NoError(t, err)
		capabilityResponse2 := commoncap.CapabilityResponse{
			Value: nm,
		}

		rawResponse2, err := pb.MarshalCapabilityResponse(capabilityResponse2)
		require.NoError(t, err)
		msg2 := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse2,
			MessageId:       []byte("messageID"),
		}

		msg.Sender = capabilityPeers[0][:]
		err = request.OnMessage(ctx, msg)
		require.NoError(t, err)

		msg2.Sender = capabilityPeers[1][:]
		err = request.OnMessage(ctx, msg2)
		require.NoError(t, err)

		select {
		case <-request.ResponseChan():
			t.Fatal("expected no response")
		default:
		}
	})

	t.Run("Send second message from non calling Don peer", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
		request, err := request.NewClientRequest(ctx, lggr, capabilityRequest, messageID, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute)
		require.NoError(t, err)
		defer request.Cancel(errors.New("test end"))

		msg.Sender = capabilityPeers[0][:]
		err = request.OnMessage(ctx, msg)
		require.NoError(t, err)

		nonDonPeer := NewP2PPeerID(t)
		msg.Sender = nonDonPeer[:]
		err = request.OnMessage(ctx, msg)
		require.NotNil(t, err)

		select {
		case <-request.ResponseChan():
			t.Fatal("expected no response")
		default:
		}
	})

	t.Run("Send second message from same peer as first message", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
		request, err := request.NewClientRequest(ctx, lggr, capabilityRequest, messageID, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute)
		require.NoError(t, err)
		defer request.Cancel(errors.New("test end"))

		msg.Sender = capabilityPeers[0][:]
		err = request.OnMessage(ctx, msg)
		require.NoError(t, err)
		err = request.OnMessage(ctx, msg)
		require.NotNil(t, err)

		select {
		case <-request.ResponseChan():
			t.Fatal("expected no response")
		default:
		}
	})

	t.Run("Send second message with same error as first", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
		request, err := request.NewClientRequest(ctx, lggr, capabilityRequest, messageID, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute)
		require.NoError(t, err)
		defer request.Cancel(errors.New("test end"))

		<-dispatcher.msgs
		<-dispatcher.msgs
		assert.Equal(t, 0, len(dispatcher.msgs))

		msgWithError := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
			Error:           types.Error_INTERNAL_ERROR,
			ErrorMsg:        "an error",
		}

		msgWithError.Sender = capabilityPeers[0][:]
		err = request.OnMessage(ctx, msgWithError)
		require.NoError(t, err)

		msgWithError.Sender = capabilityPeers[1][:]
		err = request.OnMessage(ctx, msgWithError)
		require.NoError(t, err)

		response := <-request.ResponseChan()

		assert.Equal(t, "an error", response.Err.Error())
	})

	t.Run("Send second message with different error to first", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
		request, err := request.NewClientRequest(ctx, lggr, capabilityRequest, messageID, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute)
		require.NoError(t, err)
		defer request.Cancel(errors.New("test end"))

		<-dispatcher.msgs
		<-dispatcher.msgs
		assert.Equal(t, 0, len(dispatcher.msgs))

		msgWithError := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
			Error:           types.Error_INTERNAL_ERROR,
			ErrorMsg:        "an error",
			Sender:          capabilityPeers[0][:],
		}

		msgWithError2 := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
			Error:           types.Error_INTERNAL_ERROR,
			ErrorMsg:        "an error2",
			Sender:          capabilityPeers[1][:],
		}

		err = request.OnMessage(ctx, msgWithError)
		require.NoError(t, err)
		err = request.OnMessage(ctx, msgWithError2)
		require.NoError(t, err)

		select {
		case <-request.ResponseChan():
			t.Fatal("expected no response")
		default:
		}
	})

	t.Run("Send second valid message", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
		request, err := request.NewClientRequest(ctx, lggr, capabilityRequest, messageID, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute)
		require.NoError(t, err)
		defer request.Cancel(errors.New("test end"))

		<-dispatcher.msgs
		<-dispatcher.msgs
		assert.Equal(t, 0, len(dispatcher.msgs))

		msg.Sender = capabilityPeers[0][:]
		err = request.OnMessage(ctx, msg)
		require.NoError(t, err)

		msg.Sender = capabilityPeers[1][:]
		err = request.OnMessage(ctx, msg)
		require.NoError(t, err)

		response := <-request.ResponseChan()
		resp := response.Value.Underlying["response"]

		assert.Equal(t, resp, values.NewString("response1"))
	})
}

type clientRequestTestDispatcher struct {
	msgs chan *types.MessageBody
}

func (t *clientRequestTestDispatcher) Name() string {
	return "clientRequestTestDispatcher"
}

func (t *clientRequestTestDispatcher) Start(ctx context.Context) error {
	return nil
}

func (t *clientRequestTestDispatcher) Close() error {
	return nil
}

func (t *clientRequestTestDispatcher) Ready() error {
	return nil
}

func (t *clientRequestTestDispatcher) HealthReport() map[string]error {
	return nil
}

func (t *clientRequestTestDispatcher) SetReceiver(capabilityId string, donId uint32, receiver types.Receiver) error {
	return nil
}

func (t *clientRequestTestDispatcher) RemoveReceiver(capabilityId string, donId uint32) {}

func (t *clientRequestTestDispatcher) Send(peerID p2ptypes.PeerID, msgBody *types.MessageBody) error {
	t.msgs <- msgBody
	return nil
}
