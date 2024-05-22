package remote_test

import (
	"context"
	"crypto/rand"
	"fmt"
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

func Test_TargetReceiverConsensusWithMultipleCallers(t *testing.T) {

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse) {
		responseValue, err := response.Value.Unwrap()
		require.NoError(t, err)
		assert.Equal(t, "aValue1", responseValue.(string))
	}

	// Test scenarios where the number of submissions is greater than or equal to F + 1
	testRemoteTargetConsensus(t, 1, 0, 10*time.Minute, responseTest)
	testRemoteTargetConsensus(t, 4, 3, 10*time.Minute, responseTest)
	testRemoteTargetConsensus(t, 10, 3, 10*time.Minute, responseTest)

	/*
		errResponseTest := func(t *testing.T, response commoncap.CapabilityResponse) {
			_, err := response.Value.Unwrap()
			assert.NotNil(t, err)
			//require.NoError(t, err)
			//assert.Equal(t, "aValue1", responseValue.(string))
		}

		// Test scenario where number of submissions is less than F + 1
		// TODO implement the timeout handling and cleanup logic of the execute requests cache
		testRemoteTargetConsensus(t, 4, 6, 1*time.Second, errResponseTest)

	*/
}

func testRemoteTargetConsensus(t *testing.T, numWorkflowPeers int, workflowDonF uint8,
	consensusTimeout time.Duration, responseTest func(t *testing.T, response commoncap.CapabilityResponse)) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Remote Target",
		Version:        "0.0.1",
	}
	capabilityPeerID := p2ptypes.PeerID{}
	require.NoError(t, capabilityPeerID.UnmarshalText([]byte(newPeerID())))

	capDonInfo := commoncap.DON{
		ID:      "capability-don",
		Members: []p2ptypes.PeerID{capabilityPeerID},
		F:       0,
	}

	// Define the number of workflow peers

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

	dispatcher := newTestTargetReceiverDispatcher(capabilityPeerID)

	workflowDONs := map[string]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}
	underlying := &testTargetReceiver{}

	receiver := remote.NewRemoteTargetReceiver(ctx, lggr, underlying, capInfo, &capDonInfo, workflowDONs, dispatcher, consensusTimeout)
	dispatcher.RegisterReceiver(receiver)

	callers := make([]commoncap.TargetCapability, numWorkflowPeers)
	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeerDispatcher := dispatcher.GetDispatcherForCaller(workflowPeers[i])
		caller, err := remote.NewRemoteTargetCaller(lggr, capInfo, capDonInfo, workflowDonInfo, workflowPeerDispatcher)
		require.NoError(t, err)
		dispatcher.RegisterCaller(workflowPeers[i], caller)
		callers[i] = caller
	}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "100ms",
	})
	require.NoError(t, err)

	executeInputs, err := values.NewMap(
		map[string]any{
			"executeValue1": "aValue1",
		},
	)

	wg := &sync.WaitGroup{}
	wg.Add(len(callers))

	// Fire off all the requests
	for _, caller := range callers {
		go func(caller commoncap.TargetCapability) {
			responseCh, err := caller.Execute(ctx,
				commoncap.CapabilityRequest{
					Metadata: commoncap.RequestMetadata{},
					Config:   transmissionSchedule,
					Inputs:   executeInputs,
				})

			require.NoError(t, err)

			response := <-responseCh
			responseTest(t, response)
			wg.Done()
		}(caller)
	}

	wg.Wait()
}

// Confirm that the target receiver return a response only when sufficient requests have been received

// Also confirm that any request received after the first response is replied to

// Check request times out if insufficient requests are received in a timely manner

// Check request errors as expected and all error responses are received

//  Check that requests from an incorrect don are ignored?

// Check that multiple requests from the same sender are ignored

type testTargetReceiverDispatcher struct {
	abstractDispatcher
	receiver       remotetypes.Receiver
	callers        map[p2ptypes.PeerID]remotetypes.Receiver
	receiverPeerID p2ptypes.PeerID
}

func newTestTargetReceiverDispatcher(receiverPeerID p2ptypes.PeerID) *testTargetReceiverDispatcher {
	return &testTargetReceiverDispatcher{
		receiverPeerID: receiverPeerID,
		callers:        make(map[p2ptypes.PeerID]remotetypes.Receiver),
	}
}

func (r *testTargetReceiverDispatcher) RegisterReceiver(receiver remotetypes.Receiver) {
	if r.receiver != nil {
		panic("receiver already registered")
	}

	r.receiver = receiver
}

func (r *testTargetReceiverDispatcher) GetDispatcherForCaller(callerPeerID p2ptypes.PeerID) remotetypes.Dispatcher {
	dispatcher := &callerDispatcher{
		callerPeerID: callerPeerID,
		broker:       r,
	}
	return dispatcher
}

func (r *testTargetReceiverDispatcher) RegisterCaller(callerPeerID p2ptypes.PeerID, caller remotetypes.Receiver) {
	if _, ok := r.callers[callerPeerID]; ok {
		panic("caller already registered")
	}

	r.callers[callerPeerID] = caller
}

func (r *testTargetReceiverDispatcher) SendToReceiver(peerID p2ptypes.PeerID, msg *remotetypes.MessageBody) {
	if peerID != r.receiverPeerID {
		panic("receiver peer id mismatch")
	}

	msg.Receiver = r.receiverPeerID[:]

	r.receiver.Receive(msg)
}

func (r *testTargetReceiverDispatcher) Send(callerPeerID p2ptypes.PeerID, msgBody *remotetypes.MessageBody) error {

	msgBody.Version = 1
	msgBody.Sender = r.receiverPeerID[:]
	msgBody.Receiver = callerPeerID[:]
	msgBody.Timestamp = time.Now().UnixMilli()

	if caller, ok := r.callers[callerPeerID]; ok {
		caller.Receive(msgBody)
	} else {
		return fmt.Errorf("caller not found for caller peer id %s", callerPeerID.String())
	}

	return nil
}

type callerDispatcher struct {
	abstractDispatcher
	callerPeerID p2ptypes.PeerID
	broker       *testTargetReceiverDispatcher
}

func (t *callerDispatcher) Send(peerID p2ptypes.PeerID, msgBody *remotetypes.MessageBody) error {
	msgBody.Version = 1
	msgBody.Sender = t.callerPeerID[:]
	msgBody.Timestamp = time.Now().UnixMilli()
	t.broker.SendToReceiver(peerID, msgBody)
	return nil
}

type abstractDispatcher struct {
}

func (t *abstractDispatcher) SetReceiver(capabilityId string, donId string, receiver remotetypes.Receiver) error {
	return nil
}
func (t *abstractDispatcher) RemoveReceiver(capabilityId string, donId string) {}

type testTargetReceiver struct {
}

func (t testTargetReceiver) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	return commoncap.CapabilityInfo{}, nil
}

func (t testTargetReceiver) RegisterToWorkflow(ctx context.Context, request commoncap.RegisterToWorkflowRequest) error {
	return nil
}

func (t testTargetReceiver) UnregisterFromWorkflow(ctx context.Context, request commoncap.UnregisterFromWorkflowRequest) error {
	return nil
}

func (t testTargetReceiver) Execute(ctx context.Context, request commoncap.CapabilityRequest) (<-chan commoncap.CapabilityResponse, error) {
	ch := make(chan commoncap.CapabilityResponse, 1)

	value := request.Inputs.Underlying["executeValue1"]

	ch <- commoncap.CapabilityResponse{
		Value: value,
	}

	return ch, nil
}

func libp2pMagic() []byte {
	return []byte{0x00, 0x24, 0x08, 0x01, 0x12, 0x20}
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
