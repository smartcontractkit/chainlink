package request_test

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const (
	workflowID1          = "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	workflowExecutionID1 = "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed"
	stepRef1             = "stepRef1"
)

func Test_ClientRequest_MessageValidation(t *testing.T) {
	lggr := logger.TestLogger(t)

	numWorkflowPeers := 2
	workflowPeers := make([]p2ptypes.PeerID, numWorkflowPeers)
	for i := range numWorkflowPeers {
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
			ReferenceID:         stepRef1,
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

	t.Run("Send second message with different response", func(t *testing.T) {
		ctx := t.Context()
		capabilityPeers, capDonInfo, capInfo := capabilityDon(t, 2, 1)

		dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
		req, err := request.NewClientExecuteRequest(ctx, lggr, capabilityRequest, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute)
		defer req.Cancel(errors.New("test end"))

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

		msg := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
		}
		msg.Sender = capabilityPeers[0][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		msg2.Sender = capabilityPeers[1][:]
		err = req.OnMessage(ctx, msg2)
		require.NoError(t, err)

		select {
		case <-req.ResponseChan():
			t.Fatal("expected no response")
		default:
		}
	})

	t.Run("Send second message from non calling Don peer", func(t *testing.T) {
		ctx := t.Context()
		capabilityPeers, capDonInfo, capInfo := capabilityDon(t, 2, 1)

		dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
		req, err := request.NewClientExecuteRequest(ctx, lggr, capabilityRequest, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		msg := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
		}
		msg.Sender = capabilityPeers[0][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		nonDonPeer := NewP2PPeerID(t)
		msg.Sender = nonDonPeer[:]
		err = req.OnMessage(ctx, msg)
		require.Error(t, err)

		select {
		case <-req.ResponseChan():
			t.Fatal("expected no response")
		default:
		}
	})

	t.Run("Send second message from same peer as first message", func(t *testing.T) {
		ctx := t.Context()
		capabilityPeers, capDonInfo, capInfo := capabilityDon(t, 2, 1)

		dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
		req, err := request.NewClientExecuteRequest(ctx, lggr, capabilityRequest, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		msg := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
		}
		msg.Sender = capabilityPeers[0][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)
		err = req.OnMessage(ctx, msg)
		require.Error(t, err)

		select {
		case <-req.ResponseChan():
			t.Fatal("expected no response")
		default:
		}
	})

	t.Run("Send second message with same error as first", func(t *testing.T) {
		ctx := t.Context()
		capabilityPeers, capDonInfo, capInfo := capabilityDon(t, 4, 1)

		dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
		req, err := request.NewClientExecuteRequest(ctx, lggr, capabilityRequest, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		<-dispatcher.msgs
		<-dispatcher.msgs
		assert.Empty(t, dispatcher.msgs)

		msgWithError := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
			Error:           types.Error_INTERNAL_ERROR,
			ErrorMsg:        assert.AnError.Error(),
		}

		msgWithError.Sender = capabilityPeers[0][:]
		err = req.OnMessage(ctx, msgWithError)
		require.NoError(t, err)

		msgWithError.Sender = capabilityPeers[1][:]
		err = req.OnMessage(ctx, msgWithError)
		require.NoError(t, err)

		response := <-req.ResponseChan()

		assert.Equal(t, fmt.Sprintf("%s : %s", types.Error_INTERNAL_ERROR, assert.AnError.Error()), response.Err.Error())
	})

	t.Run("Send three messages with different errors", func(t *testing.T) {
		ctx := t.Context()
		capabilityPeers, capDonInfo, capInfo := capabilityDon(t, 4, 1)

		dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
		req, err := request.NewClientExecuteRequest(ctx, lggr, capabilityRequest, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		<-dispatcher.msgs
		<-dispatcher.msgs
		assert.Empty(t, dispatcher.msgs)

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

		msgWithError3 := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
			Error:           types.Error_INTERNAL_ERROR,
			ErrorMsg:        "an error3",
			Sender:          capabilityPeers[2][:],
		}

		err = req.OnMessage(ctx, msgWithError)
		require.NoError(t, err)
		err = req.OnMessage(ctx, msgWithError2)
		require.NoError(t, err)
		err = req.OnMessage(ctx, msgWithError3)
		require.NoError(t, err)

		response := <-req.ResponseChan()
		assert.Equal(t, "received 3 errors, last error INTERNAL_ERROR : an error3", response.Err.Error())
	})

	t.Run("Execute Request", func(t *testing.T) {
		ctx := t.Context()
		capabilityPeers, capDonInfo, capInfo := capabilityDon(t, 4, 1)

		dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
		req, err := request.NewClientExecuteRequest(ctx, lggr, capabilityRequest, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		<-dispatcher.msgs
		<-dispatcher.msgs
		assert.Empty(t, dispatcher.msgs)

		msg := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
		}
		msg.Sender = capabilityPeers[0][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		msg.Sender = capabilityPeers[1][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		response := <-req.ResponseChan()
		capResponse, err := pb.UnmarshalCapabilityResponse(response.Result)
		require.NoError(t, err)

		resp := capResponse.Value.Underlying["response"]

		assert.Equal(t, resp, values.NewString("response1"))
	})

	t.Run("Executes full schedule", func(t *testing.T) {
		beholderTester := tests.Beholder(t)
		lggr, obs := logger.TestLoggerObserved(t, zapcore.DebugLevel)

		capPeers, capDonInfo, capInfo := capabilityDon(t, 3, 1)

		ctx := t.Context()
		ctxWithCancel, cancelFn := context.WithCancel(t.Context())

		// cancel the context immediately so we can verify
		// that the schedule is still executed entirely.
		cancelFn()

		// Buffered channel so the goroutines block
		// when executing the schedule
		dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody)}
		req, err := request.NewClientExecuteRequest(
			ctxWithCancel,
			lggr,
			capabilityRequest,
			capInfo,
			workflowDonInfo,
			dispatcher,
			10*time.Minute,
		)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		// Despite the context being cancelled,
		// we still send the full schedule.
		<-dispatcher.msgs
		<-dispatcher.msgs
		<-dispatcher.msgs
		assert.Empty(t, dispatcher.msgs)

		msg := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
		}
		msg.Sender = capPeers[0][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		msg.Sender = capPeers[1][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		response := <-req.ResponseChan()
		capResponse, err := pb.UnmarshalCapabilityResponse(response.Result)
		require.NoError(t, err)

		resp := capResponse.Value.Underlying["response"]

		assert.Equal(t, resp, values.NewString("response1"))

		logs := obs.FilterMessage("sending request to peers").All()
		assert.Len(t, logs, 1)

		log := logs[0]
		for _, k := range log.Context {
			if k.Key == "originalTimeout" {
				assert.Equal(t, int64(0), k.Integer)
			}

			if k.Key == "effectiveTimeout" {
				assert.Greater(t, k.Integer, int64(10*time.Second))
			}
		}

		// Verify the TransmissionsScheduledEvent data
		assert.Equal(t, 1, beholderTester.Len(t, "beholder_entity", fmt.Sprintf("%v.%v", request.TransmissionEventProtoPkg, request.TransmissionEventEntity)))

		// Get the messages for the transmission event
		messages := beholderTester.Messages(t, "beholder_entity", fmt.Sprintf("%v.%v", request.TransmissionEventProtoPkg, request.TransmissionEventEntity))
		assert.Len(t, messages, 1)

		// Unmarshal the message to verify its contents
		var event events.TransmissionsScheduledEvent
		err = proto.Unmarshal(messages[0].Body, &event)
		require.NoError(t, err)

		// Verify the event fields
		assert.Equal(t, transmission.Schedule_OneAtATime, event.ScheduleType)
		assert.Equal(t, workflowExecutionID1, event.WorkflowExecutionID)
		assert.Equal(t, "cap_id@1.0.0", event.CapabilityID)
		assert.Equal(t, stepRef1, event.StepRef)
		assert.Equal(t, fmt.Sprintf("Execute:%v:%v", workflowExecutionID1, stepRef1), event.TransmissionID)
		assert.NotEmpty(t, event.Timestamp)

		// Verify the peer delays
		assert.Len(t, event.PeerTransmissionDelays, 3)

		// Convert map to slice of delays and sort them
		var delays []int64
		for _, delay := range event.PeerTransmissionDelays {
			delays = append(delays, delay)
		}
		sort.Slice(delays, func(i, j int) bool {
			return delays[i] < delays[j]
		})

		// Verify delays are sorted and increment by 1000ms
		for i := 1; i < len(delays); i++ {
			assert.Equal(t, delays[i-1]+1000, delays[i], "delays should increment by 1000ms")
		}

		// Verify each peer ID exists in capability peers
		for peerID := range event.PeerTransmissionDelays {
			found := false
			for _, peer := range capPeers {
				if peer.String() == peerID {
					found = true
					break
				}
			}
			assert.True(t, found, "peer ID %s not found in capability peers", peerID)
		}
	})

	t.Run("Uses passed in time out if larger than schedule", func(t *testing.T) {
		lggr, obs := logger.TestLoggerObserved(t, zapcore.DebugLevel)

		capPeers, capDonInfo, capInfo := capabilityDon(t, 3, 1)

		ctx := t.Context()
		ctx, cancelFn := context.WithTimeout(ctx, 15*time.Second)
		defer cancelFn()

		// Buffered channel so the goroutines block
		// when executing the schedule
		dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody)}
		req, err := request.NewClientExecuteRequest(
			ctx,
			lggr,
			capabilityRequest,
			capInfo,
			workflowDonInfo,
			dispatcher,
			10*time.Minute,
		)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		// Despite the context being cancelled,
		// we still send the full schedule.
		<-dispatcher.msgs
		<-dispatcher.msgs
		<-dispatcher.msgs
		assert.Empty(t, dispatcher.msgs)

		msg := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
		}
		msg.Sender = capPeers[0][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		msg.Sender = capPeers[1][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		response := <-req.ResponseChan()
		capResponse, err := pb.UnmarshalCapabilityResponse(response.Result)
		require.NoError(t, err)

		resp := capResponse.Value.Underlying["response"]

		assert.Equal(t, resp, values.NewString("response1"))

		logs := obs.FilterMessage("sending request to peers").All()
		assert.Len(t, logs, 1)

		log := logs[0]
		for _, k := range log.Context {
			if k.Key == "effectiveTimeout" {
				// Greater than what it would otherwise be
				// i.e. 2 *deltaStage + margin = 12s
				assert.Greater(t, k.Integer, int64(12*time.Second))
			}
		}
	})

	// tests that (once added) metering data in the capability responses
	// will not cause the identical response calculation to break;
	// also locks in no validation of SpendUnit/SpendValue at that layer.
	t.Run("with metering metadata", func(t *testing.T) {
		capabilityPeers, capDonInfo, capInfo := capabilityDon(t, 4, 1)

		capabilityResponseWithMetering1 := commoncap.CapabilityResponse{
			Value: m,
			Metadata: commoncap.ResponseMetadata{
				Metering: []commoncap.MeteringNodeDetail{
					{SpendUnit: "testunit_a", SpendValue: "15"},
				},
			},
		}

		capabilityResponseWithMetering2 := commoncap.CapabilityResponse{
			Value: m,
			Metadata: commoncap.ResponseMetadata{
				Metering: []commoncap.MeteringNodeDetail{
					{SpendUnit: "testunit_b", SpendValue: "17"},
				},
			},
		}

		payload1, err2 := pb.MarshalCapabilityResponse(capabilityResponseWithMetering1)
		require.NoError(t, err2)

		payload2, err2 := pb.MarshalCapabilityResponse(capabilityResponseWithMetering2)
		require.NoError(t, err2)

		msg := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         rawResponse,
			MessageId:       []byte("messageID"),
		}
		msg.Payload = payload1

		ctx := t.Context()

		dispatcher := &clientRequestTestDispatcher{msgs: make(chan *types.MessageBody, 100)}
		req, err := request.NewClientExecuteRequest(ctx, lggr, capabilityRequest, capInfo,
			workflowDonInfo, dispatcher, 10*time.Minute)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		msg.Sender = capabilityPeers[0][:]
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		msg.Sender = capabilityPeers[1][:]
		msg.Payload = payload2
		err = req.OnMessage(ctx, msg)
		require.NoError(t, err)

		response := <-req.ResponseChan()
		capResponse, err := pb.UnmarshalCapabilityResponse(response.Result)
		require.NoError(t, err)

		resp := capResponse.Value.Underlying["response"]
		assert.Equal(t, resp, values.NewString("response1"))

		assert.Len(t, capResponse.Metadata.Metering, 2)
		spendUnit := capResponse.Metadata.Metering[0].SpendUnit
		spendValue := capResponse.Metadata.Metering[0].SpendValue
		p2pID := capResponse.Metadata.Metering[0].Peer2PeerID

		assert.Equal(t, "testunit_a", spendUnit)
		assert.Equal(t, "15", spendValue)
		assert.Equal(t, capabilityPeers[0].String(), p2pID)

		spendUnit = capResponse.Metadata.Metering[1].SpendUnit
		spendValue = capResponse.Metadata.Metering[1].SpendValue
		p2pID = capResponse.Metadata.Metering[1].Peer2PeerID

		assert.Equal(t, "testunit_b", spendUnit)
		assert.Equal(t, "17", spendValue)
		assert.Equal(t, capabilityPeers[1].String(), p2pID)
	})
}

func capabilityDon(t *testing.T, numCapabilityPeers int, f uint8) ([]p2ptypes.PeerID, commoncap.DON, commoncap.CapabilityInfo) {
	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := range numCapabilityPeers {
		capabilityPeers[i] = NewP2PPeerID(t)
	}

	capDonInfo := commoncap.DON{
		ID:      1,
		Members: capabilityPeers,
		F:       f,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1.0.0",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Remote Target",
		DON:            &capDonInfo,
	}
	return capabilityPeers, capDonInfo, capInfo
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

func (t *clientRequestTestDispatcher) SetReceiver(capabilityID string, donID uint32, receiver types.Receiver) error {
	return nil
}

func (t *clientRequestTestDispatcher) RemoveReceiver(capabilityID string, donID uint32) {}

func (t *clientRequestTestDispatcher) Send(peerID p2ptypes.PeerID, msgBody *types.MessageBody) error {
	t.msgs <- msgBody
	return nil
}
