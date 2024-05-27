package target_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/target"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func Test_RemoteTargetCapability_InsufficientWorkflowCallers(t *testing.T) {
	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	numCapabilityPeers := 4

	callers := testRemoteTargetReceiver(t, ctx, 10, 10, numCapabilityPeers, 3, 100*time.Millisecond)

	for _, caller := range callers {
		caller.Execute(context.Background(),
			commoncap.CapabilityRequest{
				Metadata: commoncap.RequestMetadata{
					WorkflowID:          "workflowID",
					WorkflowExecutionID: "workflowExecutionID",
				},
			})
	}

	for _, caller := range callers {
		for i := 0; i < numCapabilityPeers; i++ {
			msg := <-caller.receivedMessages
			assert.Equal(t, remotetypes.Error_TIMEOUT, msg.Error)
		}
	}
}

// Test stuck capability times out

func Test_RemoteTargetCapability_IgnoresRequestFromIncorrectPeer(t *testing.T) {
	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	numCapabilityPeers := 4

	callers := testRemoteTargetReceiver(t, ctx, 10, 9, numCapabilityPeers, 3, 100*time.Millisecond)

	for _, caller := range callers {
		caller.Execute(context.Background(),
			commoncap.CapabilityRequest{
				Metadata: commoncap.RequestMetadata{
					WorkflowID:          "workflowID",
					WorkflowExecutionID: "workflowExecutionID",
				},
			})
	}

	for _, caller := range callers {
		for i := 0; i < numCapabilityPeers; i++ {
			msg := <-caller.receivedMessages
			assert.Equal(t, remotetypes.Error_TIMEOUT, msg.Error)
		}
	}
}

func testRemoteTargetReceiver(t *testing.T, ctx context.Context, numWorkflowPeers int, workflowDonF uint8,
	numCapabilityPeers int, capabilityDonF uint8, capabilityNodeResponseTimeout time.Duration) []*receiverTestCaller {

	lggr := logger.TestLogger(t)

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

	capabilityNodes := make([]remotetypes.Receiver, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityPeer := capabilityPeers[i]
		capabilityDispatcher := broker.NewDispatcherForNode(capabilityPeer)
		capabilityNode := target.NewRemoteTargetReceiver(ctx, lggr, capabilityPeer, underlying, capInfo, capDonInfo, workflowDONs, capabilityDispatcher,
			capabilityNodeResponseTimeout)
		broker.RegisterReceiverNode(capabilityPeer, capabilityNode)
		capabilityNodes[i] = capabilityNode
	}

	workflowNodes := make([]*receiverTestCaller, numWorkflowPeers)
	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeerDispatcher := broker.NewDispatcherForNode(workflowPeers[i])
		workflowNode := newReceiverTestCaller(workflowPeers[i], capDonInfo, workflowPeerDispatcher)
		broker.RegisterReceiverNode(workflowPeers[i], workflowNode)
		workflowNodes[i] = workflowNode
	}

	return workflowNodes
}

type receiverTestCaller struct {
	peerID            p2ptypes.PeerID
	dispatcher        remotetypes.Dispatcher
	capabilityDonInfo commoncap.DON
	receivedMessages  chan *remotetypes.MessageBody
	callerDonID       string
}

func (r *receiverTestCaller) Receive(msg *remotetypes.MessageBody) {
	r.receivedMessages <- msg
}

func newReceiverTestCaller(peerID p2ptypes.PeerID, capabilityDonInfo commoncap.DON,
	dispatcher remotetypes.Dispatcher) *receiverTestCaller {
	return &receiverTestCaller{peerID: peerID, dispatcher: dispatcher, capabilityDonInfo: capabilityDonInfo,
		receivedMessages: make(chan *remotetypes.MessageBody, 100), callerDonID: "workflow-don"}
}

func (r *receiverTestCaller) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	panic("not implemented")
}

func (r *receiverTestCaller) RegisterToWorkflow(ctx context.Context, request commoncap.RegisterToWorkflowRequest) error {
	panic("not implemented")
}

func (r *receiverTestCaller) UnregisterFromWorkflow(ctx context.Context, request commoncap.UnregisterFromWorkflowRequest) error {
	panic("not implemented")
}

func (r *receiverTestCaller) Execute(ctx context.Context, req commoncap.CapabilityRequest) (<-chan commoncap.CapabilityResponse, error) {

	rawRequest, err := pb.MarshalCapabilityRequest(req)
	if err != nil {
		return nil, err
	}

	messageID, err := target.GetRequestID(req)
	if err != nil {
		return nil, err
	}

	for _, node := range r.capabilityDonInfo.Members {
		message := &remotetypes.MessageBody{
			CapabilityId:    "capability-id",
			CapabilityDonId: "capability-don",
			CallerDonId:     "workflow-don",
			Method:          remotetypes.MethodExecute,
			Payload:         rawRequest,
			MessageId:       []byte(messageID),
			Sender:          r.peerID[:],
			Receiver:        node[:],
		}

		if err = r.dispatcher.Send(node, message); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
