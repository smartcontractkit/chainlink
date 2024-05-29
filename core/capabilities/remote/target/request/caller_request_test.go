package request_test

import (
	"context"
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

func Test_CallerRequest_MessageValidation(t *testing.T) {
	lggr := logger.TestLogger(t)

	numCapabilityPeers := 2
	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityPeers[i] = NewP2PPeerID(t)
	}

	capDonInfo := commoncap.DON{
		ID:      "capability-don",
		Members: capabilityPeers,
		F:       1,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Remote Target",
		Version:        "0.0.1",
		DON:            &capDonInfo,
	}

	numWorkflowPeers := 2
	workflowPeers := make([]p2ptypes.PeerID, numWorkflowPeers)
	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeers[i] = NewP2PPeerID(t)
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      "workflow-don",
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
	capabilityRequest := commoncap.CapabilityRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID:          "workflowID",
			WorkflowExecutionID: "workflowExecutionID",
		},
		Inputs: executeInputs,
		Config: transmissionSchedule,
	}

	capabilityResponse := commoncap.CapabilityResponse{
		Value: values.NewString("response1"),
		Err:   nil,
	}

	rawResponse, err := pb.MarshalCapabilityResponse(capabilityResponse)
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

		messageID, err := target.GetMessageIDForRequest(capabilityRequest)
		require.NoError(t, err)

		dispatcher := &callerRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
		request, err := request.NewCallerRequest(ctx, lggr, capabilityRequest, messageID, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute)
		require.NoError(t, err)

		capabilityResponse2 := commoncap.CapabilityResponse{
			Value: values.NewString("response2"),
			Err:   nil,
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

		err = request.AddResponse(capabilityPeers[0], msg)
		require.NoError(t, err)
		err = request.AddResponse(capabilityPeers[1], msg2)
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

		messageID, err := target.GetMessageIDForRequest(capabilityRequest)
		require.NoError(t, err)

		dispatcher := &callerRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
		request, err := request.NewCallerRequest(ctx, lggr, capabilityRequest, messageID, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute)
		require.NoError(t, err)

		err = request.AddResponse(capabilityPeers[0], msg)
		require.NoError(t, err)
		err = request.AddResponse(NewP2PPeerID(t), msg)
		require.NotNil(t, err)
	})

	t.Run("Send second message from same peer as first message", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		messageID, err := target.GetMessageIDForRequest(capabilityRequest)
		require.NoError(t, err)

		dispatcher := &callerRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
		request, err := request.NewCallerRequest(ctx, lggr, capabilityRequest, messageID, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute)
		require.NoError(t, err)

		err = request.AddResponse(capabilityPeers[0], msg)
		require.NoError(t, err)
		err = request.AddResponse(capabilityPeers[0], msg)
		require.NotNil(t, err)
	})

	t.Run("Send second message with same error as first", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		messageID, err := target.GetMessageIDForRequest(capabilityRequest)
		require.NoError(t, err)

		dispatcher := &callerRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
		request, err := request.NewCallerRequest(ctx, lggr, capabilityRequest, messageID, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute)
		require.NoError(t, err)

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

		err = request.AddResponse(capabilityPeers[0], msgWithError)
		require.NoError(t, err)
		err = request.AddResponse(capabilityPeers[1], msgWithError)
		require.NoError(t, err)

		response := <-request.ResponseChan()

		assert.Equal(t, "an error", response.Err.Error())
	})

	t.Run("Send second valid message", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		messageID, err := target.GetMessageIDForRequest(capabilityRequest)
		require.NoError(t, err)

		dispatcher := &callerRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
		request, err := request.NewCallerRequest(ctx, lggr, capabilityRequest, messageID, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute)
		require.NoError(t, err)

		<-dispatcher.msgs
		<-dispatcher.msgs
		assert.Equal(t, 0, len(dispatcher.msgs))

		err = request.AddResponse(capabilityPeers[0], msg)
		require.NoError(t, err)
		err = request.AddResponse(capabilityPeers[1], msg)
		require.NoError(t, err)

		response := <-request.ResponseChan()

		assert.Equal(t, response.Value, values.NewString("response1"))
	})
}

type callerRequestTestDispatcher struct {
	msgs chan *types.MessageBody
}

func (t *callerRequestTestDispatcher) SetReceiver(capabilityId string, donId string, receiver types.Receiver) error {
	return nil
}

func (t *callerRequestTestDispatcher) RemoveReceiver(capabilityId string, donId string) {}

func (t *callerRequestTestDispatcher) Send(peerID p2ptypes.PeerID, msgBody *types.MessageBody) error {
	t.msgs <- msgBody
	return nil
}
