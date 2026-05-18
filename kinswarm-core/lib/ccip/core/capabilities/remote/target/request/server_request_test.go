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
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/target/request"
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
		req := request.NewServerRequest(capability, "capabilityID", 2,
			capabilityPeerID, callingDon, "requestMessageID", dispatcher, 10*time.Minute, lggr)

		err := sendValidRequest(req, workflowPeers, capabilityPeerID, rawRequest)
		require.NoError(t, err)
		err = sendValidRequest(req, workflowPeers, capabilityPeerID, rawRequest)
		assert.NotNil(t, err)
	})

	t.Run("Send message with non calling don peer", func(t *testing.T) {
		req := request.NewServerRequest(capability, "capabilityID", 2,
			capabilityPeerID, callingDon, "requestMessageID", dispatcher, 10*time.Minute, lggr)

		err := sendValidRequest(req, workflowPeers, capabilityPeerID, rawRequest)
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

		assert.NotNil(t, err)
	})

	t.Run("Send message invalid payload", func(t *testing.T) {
		req := request.NewServerRequest(capability, "capabilityID", 2,
			capabilityPeerID, callingDon, "requestMessageID", dispatcher, 10*time.Minute, lggr)

		err := sendValidRequest(req, workflowPeers, capabilityPeerID, rawRequest)
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
		assert.NoError(t, err)
		assert.Equal(t, 2, len(dispatcher.msgs))
		assert.Equal(t, dispatcher.msgs[0].Error, types.Error_INTERNAL_ERROR)
		assert.Equal(t, dispatcher.msgs[1].Error, types.Error_INTERNAL_ERROR)
	})

	t.Run("Send second valid request when capability errors", func(t *testing.T) {
		dispatcher := &testDispatcher{}
		req := request.NewServerRequest(TestErrorCapability{}, "capabilityID", 2,
			capabilityPeerID, callingDon, "requestMessageID", dispatcher, 10*time.Minute, lggr)

		err := sendValidRequest(req, workflowPeers, capabilityPeerID, rawRequest)
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
		assert.NoError(t, err)
		assert.Equal(t, 2, len(dispatcher.msgs))
		assert.Equal(t, dispatcher.msgs[0].Error, types.Error_INTERNAL_ERROR)
		assert.Equal(t, dispatcher.msgs[0].ErrorMsg, "failed to execute capability: an error")
		assert.Equal(t, dispatcher.msgs[1].Error, types.Error_INTERNAL_ERROR)
		assert.Equal(t, dispatcher.msgs[1].ErrorMsg, "failed to execute capability: an error")
	})

	t.Run("Send second valid request", func(t *testing.T) {
		dispatcher := &testDispatcher{}
		request := request.NewServerRequest(capability, "capabilityID", 2,
			capabilityPeerID, callingDon, "requestMessageID", dispatcher, 10*time.Minute, lggr)

		err := sendValidRequest(request, workflowPeers, capabilityPeerID, rawRequest)
		require.NoError(t, err)

		err = request.OnMessage(context.Background(), &types.MessageBody{
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
		assert.NoError(t, err)
		assert.Equal(t, 2, len(dispatcher.msgs))
		assert.Equal(t, dispatcher.msgs[0].Error, types.Error_OK)
		assert.Equal(t, dispatcher.msgs[1].Error, types.Error_OK)
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

func (t *testDispatcher) SetReceiver(capabilityId string, donId uint32, receiver types.Receiver) error {
	return nil
}

func (t *testDispatcher) RemoveReceiver(capabilityId string, donId uint32) {}

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
}

func (t TestErrorCapability) Execute(ctx context.Context, request commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	return commoncap.CapabilityResponse{}, errors.New("an error")
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

	return base58.Encode(peerID[:])
}

func libp2pMagic() []byte {
	return []byte{0x00, 0x24, 0x08, 0x01, 0x12, 0x20}
}
