package executable_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const (
	stepReferenceID1     = "step1"
	workflowID1          = "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	workflowID2          = "25c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1"
	workflowExecutionID1 = "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed"
	workflowExecutionID2 = "85ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeee"
	workflowOwnerID      = "0xAA"
)

func Test_Client_DonTopologies(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/CAPPL-322")

	ctx := testutils.Context(t)

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_OneAtATime,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		if assert.NoError(t, responseError) {
			mp, err := response.Value.Unwrap()
			if assert.NoError(t, err) {
				assert.Equal(t, "aValue1", mp.(map[string]any)["response"].(string))
			}
		}
	}

	capability := &TestCapability{}

	responseTimeOut := 10 * time.Minute

	var methods []func(caller commoncap.ExecutableCapability)

	methods = append(methods, func(caller commoncap.ExecutableCapability) {
		executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
		if assert.NoError(t, err) {
			executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
		}
	})

	for _, method := range methods {
		testClient(t, 1, responseTimeOut, 1, 0,
			capability, method)

		testClient(t, 10, responseTimeOut, 1, 0,
			capability, method)

		testClient(t, 1, responseTimeOut, 10, 3,
			capability, method)

		testClient(t, 10, responseTimeOut, 10, 3,
			capability, method)

		testClient(t, 10, responseTimeOut, 10, 9,
			capability, method)
	}
}

func Test_Client_TransmissionSchedules(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-104")
	ctx := testutils.Context(t)

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		if assert.NoError(t, responseError) {
			mp, err := response.Value.Unwrap()
			if assert.NoError(t, err) {
				assert.Equal(t, "aValue1", mp.(map[string]any)["response"].(string))
			}
		}
	}

	capability := &TestCapability{}

	responseTimeOut := 10 * time.Minute

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_OneAtATime,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	testClient(t, 1, responseTimeOut, 1, 0,
		capability, func(caller commoncap.ExecutableCapability) {
			executeInputs, err2 := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			if assert.NoError(t, err2) {
				executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
			}
		})
	testClient(t, 10, responseTimeOut, 10, 3,
		capability, func(caller commoncap.ExecutableCapability) {
			executeInputs, err2 := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			if assert.NoError(t, err2) {
				executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
			}
		})

	transmissionSchedule, err = values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	testClient(t, 1, responseTimeOut, 1, 0,
		capability, func(caller commoncap.ExecutableCapability) {
			executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			if assert.NoError(t, err) {
				executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
			}
		})
	testClient(t, 10, responseTimeOut, 10, 3,
		capability, func(caller commoncap.ExecutableCapability) {
			executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			if assert.NoError(t, err) {
				executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
			}
		})
}

func Test_Client_TimesOutIfInsufficientCapabilityPeerResponses(t *testing.T) {
	ctx := testutils.Context(t)

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		assert.ErrorIs(t, responseError, executable.ErrRequestExpired)
	}

	capability := &TestCapability{}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	// number of capability peers is less than F + 1

	testClient(t, 10, 1*time.Second, 10, 11,
		capability,
		func(caller commoncap.ExecutableCapability) {
			executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			if assert.NoError(t, err) {
				executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
			}
		})
}

func Test_Client_ContextCanceledBeforeQuorumReached(t *testing.T) {
	ctx, cancel := context.WithCancel(testutils.Context(t))

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		assert.ErrorIs(t, responseError, executable.ErrContextDoneBeforeResponseQuorum)
	}

	capability := &TestCapability{}
	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "20s",
	})
	require.NoError(t, err)

	cancel()
	testClient(t, 2, 20*time.Second, 2, 2,
		capability,
		func(caller commoncap.ExecutableCapability) {
			executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			if assert.NoError(t, err) {
				executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
			}
		})
}

func testClient(t *testing.T, numWorkflowPeers int, workflowNodeResponseTimeout time.Duration,
	numCapabilityPeers int, capabilityDonF uint8, underlying commoncap.ExecutableCapability,
	method func(caller commoncap.ExecutableCapability)) {
	lggr := logger.TestLogger(t)

	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityPeers[i] = NewP2PPeerID(t)
	}

	capDonInfo := commoncap.DON{
		ID:      1,
		Members: capabilityPeers,
		F:       capabilityDonF,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1.0.0",
		CapabilityType: commoncap.CapabilityTypeTrigger,
		Description:    "Remote Executable Capability",
		DON:            &capDonInfo,
	}

	workflowPeers := make([]p2ptypes.PeerID, numWorkflowPeers)
	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeers[i] = NewP2PPeerID(t)
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      2,
	}

	broker := newTestAsyncMessageBroker(t, 100)

	receivers := make([]remotetypes.Receiver, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityDispatcher := broker.NewDispatcherForNode(capabilityPeers[i])
		receiver := newTestServer(capabilityPeers[i], capabilityDispatcher, workflowDonInfo, underlying)
		broker.RegisterReceiverNode(capabilityPeers[i], receiver)
		receivers[i] = receiver
	}

	callers := make([]commoncap.ExecutableCapability, numWorkflowPeers)

	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeerDispatcher := broker.NewDispatcherForNode(workflowPeers[i])
		caller := executable.NewClient(capInfo, workflowDonInfo, workflowPeerDispatcher, workflowNodeResponseTimeout, lggr)
		servicetest.Run(t, caller)
		broker.RegisterReceiverNode(workflowPeers[i], caller)
		callers[i] = caller
	}

	servicetest.Run(t, broker)

	wg := &sync.WaitGroup{}
	wg.Add(len(callers))

	// Fire off all the requests
	for _, caller := range callers {
		go func(caller commoncap.ExecutableCapability) {
			defer wg.Done()
			method(caller)
		}(caller)
	}

	wg.Wait()
}

func executeMethod(ctx context.Context, caller commoncap.ExecutableCapability, transmissionSchedule *values.Map,
	executeInputs *values.Map, responseTest func(t *testing.T, responseCh commoncap.CapabilityResponse, responseError error), t *testing.T) {
	responseCh, err := caller.Execute(ctx,
		commoncap.CapabilityRequest{
			Metadata: commoncap.RequestMetadata{
				WorkflowID:          workflowID1,
				WorkflowExecutionID: workflowExecutionID1,
				WorkflowOwner:       workflowOwnerID,
			},
			Config: transmissionSchedule,
			Inputs: executeInputs,
		})

	responseTest(t, responseCh, err)
}

// Simple client that only responds once it has received a message from each workflow peer
type clientTestServer struct {
	peerID             p2ptypes.PeerID
	dispatcher         remotetypes.Dispatcher
	workflowDonInfo    commoncap.DON
	messageIDToSenders map[string]map[p2ptypes.PeerID]bool

	executableCapability commoncap.ExecutableCapability

	mux sync.Mutex
}

func newTestServer(peerID p2ptypes.PeerID, dispatcher remotetypes.Dispatcher, workflowDonInfo commoncap.DON,
	executableCapability commoncap.ExecutableCapability) *clientTestServer {
	return &clientTestServer{
		dispatcher:           dispatcher,
		workflowDonInfo:      workflowDonInfo,
		peerID:               peerID,
		messageIDToSenders:   make(map[string]map[p2ptypes.PeerID]bool),
		executableCapability: executableCapability,
	}
}

func (t *clientTestServer) Receive(_ context.Context, msg *remotetypes.MessageBody) {
	t.mux.Lock()
	defer t.mux.Unlock()

	sender := toPeerID(msg.Sender)
	messageID, err := executable.GetMessageID(msg)
	if err != nil {
		panic(err)
	}

	if t.messageIDToSenders[messageID] == nil {
		t.messageIDToSenders[messageID] = make(map[p2ptypes.PeerID]bool)
	}

	sendersOfMessageID := t.messageIDToSenders[messageID]
	if sendersOfMessageID[sender] {
		panic("received duplicate message")
	}

	sendersOfMessageID[sender] = true

	if len(t.messageIDToSenders[messageID]) == len(t.workflowDonInfo.Members) {
		switch msg.Method {
		case remotetypes.MethodExecute:
			capabilityRequest, err := pb.UnmarshalCapabilityRequest(msg.Payload)
			if err != nil {
				panic(err)
			}
			resp, responseErr := t.executableCapability.Execute(context.Background(), capabilityRequest)
			payload, marshalErr := pb.MarshalCapabilityResponse(resp)
			t.sendResponse(messageID, responseErr, payload, marshalErr)
		default:
			panic("unknown method")
		}
	}
}

func (t *clientTestServer) sendResponse(messageID string, responseErr error,
	payload []byte, marshalErr error) {
	for receiver := range t.messageIDToSenders[messageID] {
		var responseMsg = &remotetypes.MessageBody{
			CapabilityId:    "cap_id@1.0.0",
			CapabilityDonId: 1,
			CallerDonId:     t.workflowDonInfo.ID,
			Method:          remotetypes.MethodExecute,
			MessageId:       []byte(messageID),
			Sender:          t.peerID[:],
			Receiver:        receiver[:],
		}

		if responseErr != nil {
			responseMsg.Error = remotetypes.Error_INTERNAL_ERROR
		} else {
			if marshalErr != nil {
				panic(marshalErr)
			}
			responseMsg.Payload = payload
		}

		err := t.dispatcher.Send(receiver, responseMsg)
		if err != nil {
			panic(err)
		}
	}
}
