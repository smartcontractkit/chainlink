package remote_test

import (
	"context"
	"crypto/rand"
	"sync"
	"testing"
	"time"

	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func Test_TargetRemoteTarget(t *testing.T) {

	responseTest := func(t *testing.T, responseCh <-chan commoncap.CapabilityResponse, responseError error) {
		require.NoError(t, responseError)
		response := <-responseCh
		responseValue, err := response.Value.Unwrap()
		require.NoError(t, err)
		assert.Equal(t, "aValue1", responseValue.(string))
	}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "100ms",
	})
	require.NoError(t, err)

	// Test scenarios where the number of submissions is greater than or equal to F + 1
	testRemoteTarget(t, 1, 0, 10*time.Minute, 1, 0, 10*time.Minute, transmissionSchedule, responseTest)
	testRemoteTarget(t, 4, 3, 10*time.Minute, 4, 3, 10*time.Minute, transmissionSchedule, responseTest)
	testRemoteTarget(t, 10, 3, 10*time.Minute, 10, 3, 10*time.Minute, transmissionSchedule, responseTest)

	transmissionSchedule, err = values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_OneAtATime,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	testRemoteTarget(t, 1, 0, 10*time.Minute, 1, 0, 10*time.Minute, transmissionSchedule, responseTest)
	testRemoteTarget(t, 10, 3, 10*time.Minute, 10, 3, 10*time.Minute, transmissionSchedule, responseTest)

	//here - below tests plus additional tests for the remoteTargetCapability test

	// test capability don F handling

	/*
		here - these errors tests failing still? why?

		errResponseTest := func(t *testing.T, responseCh <-chan commoncap.CapabilityResponse, responseError error) {
			require.NoError(t, responseError)
			response := <-responseCh
			assert.NotNil(t, response.Err)
		}

		// Test scenario where number of submissions is less than F + 1

		// How to make these tests less time dependent? risk of being flaky
		testRemoteTargetConsensus(t, 4, 6, 5*time.Second, 1, 0, 1*time.Second, errResponseTest)
		testRemoteTargetConsensus(t, 10, 10, 5*time.Second, 1, 0, 1*time.Second, errResponseTest)
	*/
	//tyring to modify tests to test the caller F number handling?

	//also having issues with error test cases - since the client F handling?

	//then got threading to do

	// Context cancellation test - use an underlying capability that blocks until the context is cancelled

	// Check request errors as expected and all error responses are received

	//  Check that requests from an incorrect don are ignored?

	// Check that multiple requests from the same sender are ignored

	// Test with different transmission schedules ?

}

func testRemoteTarget(t *testing.T, numWorkflowPeers int, workflowDonF uint8, workflowNodeTimeout time.Duration,
	numCapabilityPeers int, capabilityDonF uint8, capabilityNodeResponseTimeout time.Duration, transmissionSchedule *values.Map,
	responseTest func(t *testing.T, responseCh <-chan commoncap.CapabilityResponse, responseError error)) {
	lggr := logger.TestLogger(t)
	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityPeerID := p2ptypes.PeerID{}
		require.NoError(t, capabilityPeerID.UnmarshalText([]byte(newPeerID())))
		capabilityPeers[i] = capabilityPeerID
	}

	capabilityPeerID := p2ptypes.PeerID{}
	require.NoError(t, capabilityPeerID.UnmarshalText([]byte(newPeerID())))

	capDonInfo := commoncap.DON{
		ID:      "capability-don",
		Members: capabilityPeers,
		F:       capabilityDonF,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Remote Target",
		Version:        "0.0.1",
		DON:            &capDonInfo,
	}

	workflowPeers := make([]p2ptypes.PeerID, numWorkflowPeers)
	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeerID := p2ptypes.PeerID{}
		require.NoError(t, workflowPeerID.UnmarshalText([]byte(newPeerID())))
		workflowPeers[i] = workflowPeerID
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      "workflow-don",
		F:       workflowDonF,
	}

	broker := newTestMessageBroker()

	workflowDONs := map[string]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}
	underlying := &testCapability{}

	receivers := make([]remotetypes.Receiver, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityPeer := capabilityPeers[i]
		capabilityDispatcher := broker.NewDispatcherForNode(capabilityPeer)
		receiver := remote.NewRemoteTargetReceiver(ctx, lggr, capabilityPeer, underlying, capInfo, capDonInfo, workflowDONs, capabilityDispatcher,
			capabilityNodeResponseTimeout)
		broker.RegisterReceiverNode(capabilityPeer, receiver)
		receivers[i] = receiver
	}

	callers := make([]commoncap.TargetCapability, numWorkflowPeers)
	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeerDispatcher := broker.NewDispatcherForNode(workflowPeers[i])
		caller := remote.NewRemoteTargetCaller(ctx, lggr, capInfo, workflowDonInfo, workflowPeerDispatcher, workflowNodeTimeout)
		broker.RegisterReceiverNode(workflowPeers[i], caller)
		callers[i] = caller
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
}

type testMessageBroker struct {
	receivers map[p2ptypes.PeerID]remotetypes.Receiver
}

func newTestMessageBroker() *testMessageBroker {
	return &testMessageBroker{
		receivers: make(map[p2ptypes.PeerID]remotetypes.Receiver),
	}
}

func (r *testMessageBroker) NewDispatcherForNode(nodePeerID p2ptypes.PeerID) remotetypes.Dispatcher {
	return &nodeDispatcher{
		callerPeerID: nodePeerID,
		broker:       r,
	}
}

func (r *testMessageBroker) RegisterReceiverNode(nodePeerID p2ptypes.PeerID, node remotetypes.Receiver) {
	if _, ok := r.receivers[nodePeerID]; ok {
		panic("node already registered")
	}

	r.receivers[nodePeerID] = node
}

func (r *testMessageBroker) Send(msg *remotetypes.MessageBody) {
	receiverId := toPeerID(msg.Receiver)

	if receiver, ok := r.receivers[receiverId]; ok {
		receiver.Receive(msg)
	} else {
		panic("receiver not found for peer id")
	}

}

func toPeerID(id []byte) p2ptypes.PeerID {
	return [32]byte(id)
}

type nodeDispatcher struct {
	callerPeerID p2ptypes.PeerID
	broker       *testMessageBroker
}

func (t *nodeDispatcher) Send(peerID p2ptypes.PeerID, msgBody *remotetypes.MessageBody) error {
	msgBody.Version = 1
	msgBody.Sender = t.callerPeerID[:]
	msgBody.Receiver = peerID[:]
	msgBody.Timestamp = time.Now().UnixMilli()
	t.broker.Send(msgBody)
	return nil
}

func (t *nodeDispatcher) SetReceiver(capabilityId string, donId string, receiver remotetypes.Receiver) error {
	return nil
}
func (t *nodeDispatcher) RemoveReceiver(capabilityId string, donId string) {}

type testCapability struct {
}

func (t testCapability) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	return commoncap.CapabilityInfo{}, nil
}

func (t testCapability) RegisterToWorkflow(ctx context.Context, request commoncap.RegisterToWorkflowRequest) error {
	return nil
}

func (t testCapability) UnregisterFromWorkflow(ctx context.Context, request commoncap.UnregisterFromWorkflowRequest) error {
	return nil
}

func (t testCapability) Execute(ctx context.Context, request commoncap.CapabilityRequest) (<-chan commoncap.CapabilityResponse, error) {
	ch := make(chan commoncap.CapabilityResponse, 1)

	value := request.Inputs.Underlying["executeValue1"]

	ch <- commoncap.CapabilityResponse{
		Value: value,
	}

	return ch, nil
}

func newPeerID() string {
	var privKey [32]byte
	_, err := rand.Read(privKey[:])
	if err != nil {
		panic(err)
	}

	peerID := append(libp2pMagic(), privKey[:]...)

	return base58.Encode(peerID[:])
}

func libp2pMagic() []byte {
	return []byte{0x00, 0x24, 0x08, 0x01, 0x12, 0x20}
}
