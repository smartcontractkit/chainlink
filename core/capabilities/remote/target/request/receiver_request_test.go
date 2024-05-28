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
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func Test_ReceiverRequest_MessageValidation(t *testing.T) {
	lggr := logger.TestLogger(t)
	capability := target.testCapability{}
	capabilityPeerID := target.newP2PPeerID(t)

	numWorkflowPeers := 2
	workflowPeers := make([]p2ptypes.PeerID, numWorkflowPeers)
	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeers[i] = target.newP2PPeerID(t)
	}

	callingDon := commoncap.DON{
		Members: workflowPeers,
		ID:      "workflow-don",
		F:       1,
	}

	dispatcher := &testDispatcher{}

	executeInputs, err := values.NewMap(
		map[string]any{
			"executeValue1": "aValue1",
		},
	)
	require.NoError(t, err)

	capabilityRequest := commoncap.CapabilityRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID:          "workflowID",
			WorkflowExecutionID: "workflowExecutionID",
		},
		Inputs: executeInputs,
	}

	rawRequest, err := pb.MarshalCapabilityRequest(capabilityRequest)
	require.NoError(t, err)

	t.Run("Send duplicate message", func(t *testing.T) {
		request := request.NewReceiverRequest(lggr, capability, "capabilityID", "capabilityDonID",
			capabilityPeerID, callingDon, "requestMessageID", dispatcher, 10*time.Minute)

		err := sendValidRequest(request, workflowPeers, capabilityPeerID, rawRequest)
		require.NoError(t, err)
		err = sendValidRequest(request, workflowPeers, capabilityPeerID, rawRequest)
		assert.NotNil(t, err)
	})

	t.Run("Send message with non calling don peer", func(t *testing.T) {
		request := request.NewReceiverRequest(lggr, capability, "capabilityID", "capabilityDonID",
			capabilityPeerID, callingDon, "requestMessageID", dispatcher, 10*time.Minute)

		err := sendValidRequest(request, workflowPeers, capabilityPeerID, rawRequest)
		require.NoError(t, err)

		nonDonPeer := target.newP2PPeerID(t)
		err = request.Receive(context.Background(), &types.MessageBody{
			Version:         0,
			Sender:          nonDonPeer[:],
			Receiver:        capabilityPeerID[:],
			MessageId:       []byte("workflowID" + "workflowExecutionID"),
			CapabilityId:    "capabilityID",
			CapabilityDonId: "capabilityDonID",
			CallerDonId:     "workflow-don",
			Method:          types.MethodExecute,
			Payload:         rawRequest,
		})

		assert.NotNil(t, err)
	})

	t.Run("Send message invalid payload", func(t *testing.T) {
		request := request.NewReceiverRequest(lggr, capability, "capabilityID", "capabilityDonID",
			capabilityPeerID, callingDon, "requestMessageID", dispatcher, 10*time.Minute)

		err := sendValidRequest(request, workflowPeers, capabilityPeerID, rawRequest)
		require.NoError(t, err)

		err = request.Receive(context.Background(), &types.MessageBody{
			Version:         0,
			Sender:          workflowPeers[1][:],
			Receiver:        capabilityPeerID[:],
			MessageId:       []byte("workflowID" + "workflowExecutionID"),
			CapabilityId:    "capabilityID",
			CapabilityDonId: "capabilityDonID",
			CallerDonId:     "workflow-don",
			Method:          types.MethodExecute,
			Payload:         append(rawRequest, []byte("asdf")...),
		})
		assert.NoError(t, err)
		assert.Equal(t, 2, len(dispatcher.msgs))
		assert.Equal(t, dispatcher.msgs[0].Error, types.Error_INTERNAL_ERROR)
		assert.Equal(t, dispatcher.msgs[1].Error, types.Error_INTERNAL_ERROR)

	})

	t.Run("Send second valid request when capability errors", func(t *testing.T) {

		dispatcher := &testDispatcher{}
		request := request.NewReceiverRequest(lggr, target.testErrorCapability{}, "capabilityID", "capabilityDonID",
			capabilityPeerID, callingDon, "requestMessageID", dispatcher, 10*time.Minute)

		err := sendValidRequest(request, workflowPeers, capabilityPeerID, rawRequest)
		require.NoError(t, err)

		err = request.Receive(context.Background(), &types.MessageBody{
			Version:         0,
			Sender:          workflowPeers[1][:],
			Receiver:        capabilityPeerID[:],
			MessageId:       []byte("workflowID" + "workflowExecutionID"),
			CapabilityId:    "capabilityID",
			CapabilityDonId: "capabilityDonID",
			CallerDonId:     "workflow-don",
			Method:          types.MethodExecute,
			Payload:         rawRequest,
		})
		assert.NoError(t, err)
		assert.Equal(t, 2, len(dispatcher.msgs))
		assert.Equal(t, dispatcher.msgs[0].Error, types.Error_INTERNAL_ERROR)
		assert.Equal(t, dispatcher.msgs[1].Error, types.Error_INTERNAL_ERROR)

	})

	t.Run("Send second valid request", func(t *testing.T) {
		dispatcher := &testDispatcher{}
		request := request.NewReceiverRequest(lggr, capability, "capabilityID", "capabilityDonID",
			capabilityPeerID, callingDon, "requestMessageID", dispatcher, 10*time.Minute)

		err := sendValidRequest(request, workflowPeers, capabilityPeerID, rawRequest)
		require.NoError(t, err)

		err = request.Receive(context.Background(), &types.MessageBody{
			Version:         0,
			Sender:          workflowPeers[1][:],
			Receiver:        capabilityPeerID[:],
			MessageId:       []byte("workflowID" + "workflowExecutionID"),
			CapabilityId:    "capabilityID",
			CapabilityDonId: "capabilityDonID",
			CallerDonId:     "workflow-don",
			Method:          types.MethodExecute,
			Payload:         rawRequest,
		})
		assert.NoError(t, err)
		assert.Equal(t, 2, len(dispatcher.msgs))
		assert.Equal(t, dispatcher.msgs[0].Error, types.Error_OK)
		assert.Equal(t, dispatcher.msgs[1].Error, types.Error_OK)
	})
}

type receiverRequest interface {
	Receive(ctx context.Context, msg *types.MessageBody) error
}

func sendValidRequest(request receiverRequest, workflowPeers []p2ptypes.PeerID, capabilityPeerID p2ptypes.PeerID,
	rawRequest []byte) error {
	return request.Receive(context.Background(), &types.MessageBody{
		Version:         0,
		Sender:          workflowPeers[0][:],
		Receiver:        capabilityPeerID[:],
		MessageId:       []byte("workflowID" + "workflowExecutionID"),
		CapabilityId:    "capabilityID",
		CapabilityDonId: "capabilityDonID",
		CallerDonId:     "workflow-don",
		Method:          types.MethodExecute,
		Payload:         rawRequest,
	})

}

type testDispatcher struct {
	msgs []*types.MessageBody
}

func (t *testDispatcher) SetReceiver(capabilityId string, donId string, receiver types.Receiver) error {
	return nil
}

func (t *testDispatcher) RemoveReceiver(capabilityId string, donId string) {}

func (t *testDispatcher) Send(peerID p2ptypes.PeerID, msgBody *types.MessageBody) error {
	t.msgs = append(t.msgs, msgBody)
	return nil
}
