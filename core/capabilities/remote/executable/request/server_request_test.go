package request_test

import (
	"context"
	"crypto/rand"
	"errors"
	"testing"
	"time"

	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func Test_ServerRequest_MessageValidation(t *testing.T) {
	lggr := logger.TestLogger(t)
	capability := TestCapability{}
	capabilityPeerID := NewP2PPeerID(t)

	numWorkflowPeers := 2
	workflowPeers := make([]p2ptypes.PeerID, numWorkflowPeers)
	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeers[i] = NewP2PPeerID(t)
	}

	callingDon := commoncap.DON{
		Members: workflowPeers,
		ID:      1,
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
		req, err := request.NewServerRequest(capability, types.MethodExecute, "capabilityID", 2,
			capabilityPeerID, callingDon, "requestMessageID", dispatcher, 10*time.Minute, lggr)
		require.NoError(t, err)

		err = sendValidRequest(req, workflowPeers, capabilityPeerID, rawRequest)
		require.NoError(t, err)
		err = sendValidRequest(req, workflowPeers, capabilityPeerID, rawRequest)
		assert.Error(t, err)
	})

	t.Run("Send message with non calling don peer", func(t *testing.T) {
		req, err := request.NewServerRequest(capability, types.MethodExecute, "capabilityID", 2,
			capabilityPeerID, callingDon, "requestMessageID", dispatcher, 10*time.Minute, lggr)
		require.NoError(t, err)

		err = sendValidRequest(req, workflowPeers, capabilityPeerID, rawRequest)
		require.NoError(t, err)

		nonDonPeer := NewP2PPeerID(t)
		err = req.OnMessage(context.Background(), &types.MessageBody{
			Version:         0,
			Sender:          nonDonPeer[:],
			Receiver:        capabilityPeerID[:],
			MessageId:       []byte("workflowID" + "workflowExecutionID"),
			CapabilityId:    "capabilityID",
			CapabilityDonId: 2,
			CallerDonId:     1,
			Method:          types.MethodExecute,
			Payload:         rawRequest,
		})

		assert.Error(t, err)
	})

	t.Run("Send message invalid payload", func(t *testing.T) {
		req, err := request.NewServerRequest(capability, types.MethodExecute, "capabilityID", 2,
			capabilityPeerID, callingDon, "requestMessageID", dispatcher, 10*time.Minute, lggr)
		require.NoError(t, err)

		err = sendValidRequest(req, workflowPeers, capabilityPeerID, rawRequest)
		require.NoError(t, err)

		err = req.OnMessage(context.Background(), &types.MessageBody{
			Version:         0,
			Sender:          workflowPeers[1][:],
			Receiver:        capabilityPeerID[:],
			MessageId:       []byte("workflowID" + "workflowExecutionID"),
			CapabilityId:    "capabilityID",
			CapabilityDonId: 2,
			CallerDonId:     1,
			Method:          types.MethodExecute,
			Payload:         append(rawRequest, []byte("asdf")...),
		})
		require.NoError(t, err)
		assert.Len(t, dispatcher.msgs, 2)
		assert.Equal(t, types.Error_INTERNAL_ERROR, dispatcher.msgs[0].Error)
		assert.Equal(t, types.Error_INTERNAL_ERROR, dispatcher.msgs[1].Error)
	})

	t.Run("Send second valid request when capability errors", func(t *testing.T) {
		dispatcher := &testDispatcher{}
		req, err := request.NewServerRequest(TestErrorCapability{err: errors.New("an error")}, types.MethodExecute, "capabilityID", 2,
			capabilityPeerID, callingDon, "requestMessageID", dispatcher, 10*time.Minute, lggr)
		require.NoError(t, err)

		err = sendValidRequest(req, workflowPeers, capabilityPeerID, rawRequest)
		require.NoError(t, err)

		err = req.OnMessage(context.Background(), &types.MessageBody{
			Version:         0,
			Sender:          workflowPeers[1][:],
			Receiver:        capabilityPeerID[:],
			MessageId:       []byte("workflowID" + "workflowExecutionID"),
			CapabilityId:    "capabilityID",
			CapabilityDonId: 2,
			CallerDonId:     1,
			Method:          types.MethodExecute,
			Payload:         rawRequest,
		})
		require.NoError(t, err)
		assert.Len(t, dispatcher.msgs, 2)
		assert.Equal(t, types.Error_INTERNAL_ERROR, dispatcher.msgs[0].Error)
		assert.Equal(t, "failed to execute capability", dispatcher.msgs[0].ErrorMsg)
		assert.Equal(t, types.Error_INTERNAL_ERROR, dispatcher.msgs[1].Error)
		assert.Equal(t, "failed to execute capability", dispatcher.msgs[1].ErrorMsg)
	})

	t.Run("Reportable errors are returned to the caller", func(t *testing.T) {
		dispatcher := &testDispatcher{}
		req, err := request.NewServerRequest(TestErrorCapability{err: commoncap.NewRemoteReportableError(errors.New("error details"))}, types.MethodExecute, "capabilityID", 2,
			capabilityPeerID, callingDon, "requestMessageID", dispatcher, 10*time.Minute, lggr)
		require.NoError(t, err)

		err = sendValidRequest(req, workflowPeers, capabilityPeerID, rawRequest)
		require.NoError(t, err)

		err = req.OnMessage(context.Background(), &types.MessageBody{
			Version:         0,
			Sender:          workflowPeers[1][:],
			Receiver:        capabilityPeerID[:],
			MessageId:       []byte("workflowID" + "workflowExecutionID"),
			CapabilityId:    "capabilityID",
			CapabilityDonId: 2,
			CallerDonId:     1,
			Method:          types.MethodExecute,
			Payload:         rawRequest,
		})
		require.NoError(t, err)
		assert.Len(t, dispatcher.msgs, 2)
		assert.Equal(t, types.Error_INTERNAL_ERROR, dispatcher.msgs[0].Error)
		assert.Equal(t, "failed to execute capability: error details", dispatcher.msgs[0].ErrorMsg)
		assert.Equal(t, types.Error_INTERNAL_ERROR, dispatcher.msgs[1].Error)
		assert.Equal(t, "failed to execute capability: error details", dispatcher.msgs[1].ErrorMsg)
	})

	t.Run("Execute capability", func(t *testing.T) {
		dispatcher := &testDispatcher{}
		req, err := request.NewServerRequest(capability, types.MethodExecute, "capabilityID", 2,
			capabilityPeerID, callingDon, "requestMessageID", dispatcher, 10*time.Minute, lggr)
		require.NoError(t, err)

		err = sendValidRequest(req, workflowPeers, capabilityPeerID, rawRequest)
		require.NoError(t, err)

		err = req.OnMessage(context.Background(), &types.MessageBody{
			Version:         0,
			Sender:          workflowPeers[1][:],
			Receiver:        capabilityPeerID[:],
			MessageId:       []byte("workflowID" + "workflowExecutionID"),
			CapabilityId:    "capabilityID",
			CapabilityDonId: 2,
			CallerDonId:     1,
			Method:          types.MethodExecute,
			Payload:         rawRequest,
		})
		require.NoError(t, err)
		assert.Len(t, dispatcher.msgs, 2)
		assert.Equal(t, types.Error_OK, dispatcher.msgs[0].Error)
		assert.Equal(t, types.Error_OK, dispatcher.msgs[1].Error)
	})
}

type serverRequest interface {
	OnMessage(ctx context.Context, msg *types.MessageBody) error
}

func sendValidRequest(request serverRequest, workflowPeers []p2ptypes.PeerID, capabilityPeerID p2ptypes.PeerID,
	rawRequest []byte) error {
	return request.OnMessage(context.Background(), &types.MessageBody{
		Version:         0,
		Sender:          workflowPeers[0][:],
		Receiver:        capabilityPeerID[:],
		MessageId:       []byte("workflowID" + "workflowExecutionID"),
		CapabilityId:    "capabilityID",
		CapabilityDonId: 2,
		CallerDonId:     1,
		Method:          types.MethodExecute,
		Payload:         rawRequest,
	})
}

type testDispatcher struct {
	msgs []*types.MessageBody
}

func (t *testDispatcher) Name() string {
	return "testDispatcher"
}

func (t *testDispatcher) Start(ctx context.Context) error {
	return nil
}

func (t *testDispatcher) Close() error {
	return nil
}

func (t *testDispatcher) Ready() error {
	return nil
}

func (t *testDispatcher) HealthReport() map[string]error {
	return nil
}

func (t *testDispatcher) SetReceiver(capabilityID string, donID uint32, receiver types.Receiver) error {
	return nil
}

func (t *testDispatcher) RemoveReceiver(capabilityID string, donID uint32) {}

func (t *testDispatcher) Send(peerID p2ptypes.PeerID, msgBody *types.MessageBody) error {
	t.msgs = append(t.msgs, msgBody)
	return nil
}

type abstractTestCapability struct {
}

func (t abstractTestCapability) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	return commoncap.CapabilityInfo{}, nil
}

func (t abstractTestCapability) RegisterToWorkflow(ctx context.Context, request commoncap.RegisterToWorkflowRequest) error {
	return nil
}

func (t abstractTestCapability) UnregisterFromWorkflow(ctx context.Context, request commoncap.UnregisterFromWorkflowRequest) error {
	return nil
}

type TestCapability struct {
	abstractTestCapability
}

func (t TestCapability) Execute(ctx context.Context, request commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	value := request.Inputs.Underlying["executeValue1"]

	response, err := values.NewMap(map[string]any{"response": value})
	if err != nil {
		return commoncap.CapabilityResponse{}, err
	}

	return commoncap.CapabilityResponse{
		Value: response,
	}, nil
}

type TestErrorCapability struct {
	abstractTestCapability
	err error
}

func (t TestErrorCapability) Execute(ctx context.Context, request commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	return commoncap.CapabilityResponse{}, t.err
}

func (t TestErrorCapability) RegisterToWorkflow(ctx context.Context, request commoncap.RegisterToWorkflowRequest) error {
	return t.err
}

func (t TestErrorCapability) UnregisterFromWorkflow(ctx context.Context, request commoncap.UnregisterFromWorkflowRequest) error {
	return t.err
}

func NewP2PPeerID(t *testing.T) p2ptypes.PeerID {
	id := p2ptypes.PeerID{}
	require.NoError(t, id.UnmarshalText([]byte(NewPeerID())))
	return id
}

func NewPeerID() string {
	var privKey [32]byte
	_, err := rand.Read(privKey[:])
	if err != nil {
		panic(err)
	}

	peerID := append(libp2pMagic(), privKey[:]...)

	return base58.Encode(peerID)
}

func libp2pMagic() []byte {
	return []byte{0x00, 0x24, 0x08, 0x01, 0x12, 0x20}
}
