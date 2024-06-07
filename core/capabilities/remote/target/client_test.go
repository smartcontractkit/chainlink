package target_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/target"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func Test_Client_DonTopologies(t *testing.T) {
	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_OneAtATime,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	responseTest := func(t *testing.T, responseCh <-chan commoncap.CapabilityResponse, responseError error) {
		require.NoError(t, responseError)
		response := <-responseCh
		responseValue, err := response.Value.Unwrap()
		require.NoError(t, err)
		assert.Equal(t, "aValue1", responseValue.(string))
	}

	capability := &TestCapability{}

	responseTimeOut := 10 * time.Minute

	testClient(ctx, t, 1, responseTimeOut, 1, 0,
		capability, transmissionSchedule, responseTest)

	testClient(ctx, t, 10, responseTimeOut, 1, 0,
		capability, transmissionSchedule, responseTest)

	testClient(ctx, t, 1, responseTimeOut, 10, 3,
		capability, transmissionSchedule, responseTest)

	testClient(ctx, t, 10, responseTimeOut, 10, 3,
		capability, transmissionSchedule, responseTest)

	testClient(ctx, t, 10, responseTimeOut, 10, 9,
		capability, transmissionSchedule, responseTest)
}

func Test_Client_TransmissionSchedules(t *testing.T) {
	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	responseTest := func(t *testing.T, responseCh <-chan commoncap.CapabilityResponse, responseError error) {
		require.NoError(t, responseError)
		response := <-responseCh
		responseValue, err := response.Value.Unwrap()
		require.NoError(t, err)
		assert.Equal(t, "aValue1", responseValue.(string))
	}

	capability := &TestCapability{}

	responseTimeOut := 10 * time.Minute

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_OneAtATime,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	testClient(ctx, t, 1, responseTimeOut, 1, 0,
		capability, transmissionSchedule, responseTest)
	testClient(ctx, t, 10, responseTimeOut, 10, 3,
		capability, transmissionSchedule, responseTest)

	transmissionSchedule, err = values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	testClient(ctx, t, 1, responseTimeOut, 1, 0,
		capability, transmissionSchedule, responseTest)
	testClient(ctx, t, 10, responseTimeOut, 10, 3,
		capability, transmissionSchedule, responseTest)
}

func Test_Client_TimesOutIfInsufficientCapabilityPeerResponses(t *testing.T) {
	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	responseTest := func(t *testing.T, responseCh <-chan commoncap.CapabilityResponse, responseError error) {
		require.NoError(t, responseError)
		response := <-responseCh
		assert.NotNil(t, response.Err)
	}

	capability := &TestCapability{}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	// number of capability peers is less than F + 1

	testClient(ctx, t, 10, 1*time.Second, 10, 11,
		capability, transmissionSchedule, responseTest)
}

func testClient(ctx context.Context, t *testing.T, numWorkflowPeers int, workflowNodeResponseTimeout time.Duration,
	numCapabilityPeers int, capabilityDonF uint8, underlying commoncap.TargetCapability, transmissionSchedule *values.Map,
	responseTest func(t *testing.T, responseCh <-chan commoncap.CapabilityResponse, responseError error)) {
	lggr := logger.TestLogger(t)

	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityPeers[i] = NewP2PPeerID(t)
	}

	capDonInfo := commoncap.DON{
		ID:      "capability-don",
		Members: capabilityPeers,
		F:       capabilityDonF,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1.0.0",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Remote Target",
		DON:            &capDonInfo,
	}

	workflowPeers := make([]p2ptypes.PeerID, numWorkflowPeers)
	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeers[i] = NewP2PPeerID(t)
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      "workflow-don",
	}

	broker := newTestMessageBroker()

	receivers := make([]remotetypes.Receiver, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityDispatcher := broker.NewDispatcherForNode(capabilityPeers[i])
		receiver := newTestServer(capabilityPeers[i], capabilityDispatcher, workflowDonInfo, underlying)
		broker.RegisterReceiverNode(capabilityPeers[i], receiver)
		receivers[i] = receiver
	}

	callers := make([]commoncap.TargetCapability, numWorkflowPeers)
	srvcs := make([]services.Service, numWorkflowPeers)
	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeerDispatcher := broker.NewDispatcherForNode(workflowPeers[i])
		caller := target.NewClient(capInfo, workflowDonInfo, workflowPeerDispatcher, workflowNodeResponseTimeout, lggr)
		require.NoError(t, caller.Start(ctx))
		broker.RegisterReceiverNode(workflowPeers[i], caller)
		callers[i] = caller
		srvcs[i] = caller
	}

	executeInputs, err := values.NewMap(
		map[string]any{
			"executeValue1": "aValue1",
		},
	)

	require.NoError(t, err)

	wg := &sync.WaitGroup{}
	wg.Add(len(callers))

	// Fire off all the requests
	for _, caller := range callers {
		go func(caller commoncap.TargetCapability) {
			responseCh, err := caller.Execute(ctx,
				commoncap.CapabilityRequest{
					Metadata: commoncap.RequestMetadata{
						WorkflowID:          "workflowID",
						WorkflowExecutionID: "workflowExecutionID",
					},
					Config: transmissionSchedule,
					Inputs: executeInputs,
				})

			responseTest(t, responseCh, err)
			wg.Done()
		}(caller)
	}

	wg.Wait()
	for i := 0; i < numWorkflowPeers; i++ {
		require.NoError(t, srvcs[i].Close())
	}
}

// Simple client that only responds once it has received a message from each workflow peer
type clientTestServer struct {
	peerID             p2ptypes.PeerID
	dispatcher         remotetypes.Dispatcher
	workflowDonInfo    commoncap.DON
	messageIDToSenders map[string]map[p2ptypes.PeerID]bool

	targetCapability commoncap.TargetCapability

	mux sync.Mutex
}

func newTestServer(peerID p2ptypes.PeerID, dispatcher remotetypes.Dispatcher, workflowDonInfo commoncap.DON,
	targetCapability commoncap.TargetCapability) *clientTestServer {
	return &clientTestServer{
		dispatcher:         dispatcher,
		workflowDonInfo:    workflowDonInfo,
		peerID:             peerID,
		messageIDToSenders: make(map[string]map[p2ptypes.PeerID]bool),
		targetCapability:   targetCapability,
	}
}

func (t *clientTestServer) Receive(msg *remotetypes.MessageBody) {
	t.mux.Lock()
	defer t.mux.Unlock()

	sender := toPeerID(msg.Sender)
	messageID := target.GetMessageID(msg)

	if t.messageIDToSenders[messageID] == nil {
		t.messageIDToSenders[messageID] = make(map[p2ptypes.PeerID]bool)
	}

	sendersOfMessageID := t.messageIDToSenders[messageID]
	if sendersOfMessageID[sender] {
		panic("received duplicate message")
	}

	sendersOfMessageID[sender] = true

	if len(t.messageIDToSenders[messageID]) == len(t.workflowDonInfo.Members) {
		capabilityRequest, err := pb.UnmarshalCapabilityRequest(msg.Payload)
		if err != nil {
			panic(err)
		}

		respCh, responseErr := t.targetCapability.Execute(context.Background(), capabilityRequest)
		resp := <-respCh

		for receiver := range t.messageIDToSenders[messageID] {
			var responseMsg = &remotetypes.MessageBody{
				CapabilityId:    "cap_id@1.0.0",
				CapabilityDonId: "capability-don",
				CallerDonId:     t.workflowDonInfo.ID,
				Method:          remotetypes.MethodExecute,
				MessageId:       []byte(messageID),
				Sender:          t.peerID[:],
				Receiver:        receiver[:],
			}

			if responseErr != nil {
				responseMsg.Error = remotetypes.Error_INTERNAL_ERROR
			} else {
				payload, marshalErr := pb.MarshalCapabilityResponse(resp)
				if marshalErr != nil {
					panic(marshalErr)
				}
				responseMsg.Payload = payload
			}

			err = t.dispatcher.Send(receiver, responseMsg)
			if err != nil {
				panic(err)
			}
		}
	}
}

type TestDispatcher struct {
	sentMessagesCh chan *remotetypes.MessageBody
	receiver       remotetypes.Receiver
}

func NewTestDispatcher() *TestDispatcher {
	return &TestDispatcher{
		sentMessagesCh: make(chan *remotetypes.MessageBody, 1),
	}
}

func (t *TestDispatcher) SendToReceiver(msgBody *remotetypes.MessageBody) {
	t.receiver.Receive(msgBody)
}

func (t *TestDispatcher) SetReceiver(capabilityId string, donId string, receiver remotetypes.Receiver) error {
	t.receiver = receiver
	return nil
}

func (t *TestDispatcher) RemoveReceiver(capabilityId string, donId string) {}

func (t *TestDispatcher) Send(peerID p2ptypes.PeerID, msgBody *remotetypes.MessageBody) error {
	t.sentMessagesCh <- msgBody
	return nil
}
