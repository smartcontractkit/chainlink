package sharding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	dispatchermocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
)

func makePeerID(b byte) ragetypes.PeerID {
	var id ragetypes.PeerID
	id[0] = b
	return id
}

func TestExecutionStatusUpdateSender_SendsToAllMembers(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	dispatcher := dispatchermocks.NewDispatcher(t)
	dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil)

	secondary := commoncap.DON{
		ID:      2,
		F:       1,
		Members: []ragetypes.PeerID{makePeerID(1), makePeerID(2), makePeerID(3)},
	}

	sender := NewExecutionStatusUpdateSender(dispatcher, 1, secondary, logger.Test(t))
	err := sender.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { _ = sender.Close() })

	msg := &ringpb.ExecutionStatusUpdate{
		WorkflowId:     "wf-1",
		TriggerEventId: "evt-1",
		TriggerIndex:   0,
		Status:         ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS,
		PrimaryDonId:   1,
	}

	sender.Send(ctx, msg)

	dispatcher.AssertNumberOfCalls(t, "Send", len(secondary.Members))
	for _, peer := range secondary.Members {
		dispatcher.AssertCalled(t, "Send", peer, mock.Anything)
	}
}

func TestExecutionStatusUpdateReceiver_QuorumAggregation(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	primary := commoncap.DON{
		ID:      1,
		F:       1,
		Members: []ragetypes.PeerID{makePeerID(1), makePeerID(2), makePeerID(3)},
	}

	received := make(chan *ringpb.ExecutionStatusUpdate, 1)
	handler := func(msg *ringpb.ExecutionStatusUpdate) {
		received <- msg
	}

	rcvr := NewExecutionStatusUpdateReceiver(primary, handler, logger.Test(t))
	err := rcvr.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { _ = rcvr.Close() })

	msg := &ringpb.ExecutionStatusUpdate{
		WorkflowId:     "wf-1",
		TriggerEventId: "evt-1",
		TriggerIndex:   0,
		Status:         ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS,
		PrimaryDonId:   1,
	}
	payload, err := proto.Marshal(msg)
	require.NoError(t, err)

	body := &remotetypes.MessageBody{
		Method:  remotetypes.MethodExecutionStatusUpdate,
		Payload: payload,
		Sender:  primary.Members[0][:],
	}

	rcvr.Receive(ctx, body)
	select {
	case <-received:
		t.Fatal("handler should not be called before quorum")
	default:
	}

	body.Sender = primary.Members[1][:]
	rcvr.Receive(ctx, body)

	select {
	case got := <-received:
		assert.Equal(t, "wf-1", got.WorkflowId)
		assert.Equal(t, ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS, got.Status)
	case <-ctx.Done():
		t.Fatal("timeout waiting for handler after quorum")
	}
}

func TestExecutionStatusUpdateReceiver_IgnoresUnknownPeers(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	primary := commoncap.DON{
		ID:      1,
		F:       0,
		Members: []ragetypes.PeerID{makePeerID(1)},
	}

	handlerCalled := false
	handler := func(msg *ringpb.ExecutionStatusUpdate) {
		handlerCalled = true
	}

	rcvr := NewExecutionStatusUpdateReceiver(primary, handler, logger.Test(t))
	err := rcvr.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { _ = rcvr.Close() })

	msg := &ringpb.ExecutionStatusUpdate{
		WorkflowId:     "wf-1",
		TriggerEventId: "evt-1",
		Status:         ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS,
	}
	payload, err := proto.Marshal(msg)
	require.NoError(t, err)

	unknownPeer := makePeerID(99)
	body := &remotetypes.MessageBody{
		Method:  remotetypes.MethodExecutionStatusUpdate,
		Payload: payload,
		Sender:  unknownPeer[:],
	}

	rcvr.Receive(ctx, body)
	assert.False(t, handlerCalled, "handler should not be called for unknown peer")
}

func TestShardHeartbeatReceiver_UpdatesLastSeen(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	primary := commoncap.DON{
		ID:      1,
		F:       0,
		Members: []ragetypes.PeerID{makePeerID(1)},
	}

	handler := func(msg *ringpb.ShardHeartbeat) {}

	rcvr := NewShardHeartbeatReceiver(primary, handler, logger.Test(t))
	err := rcvr.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { _ = rcvr.Close() })

	assert.Equal(t, int64(0), rcvr.LastSeen())

	hb := &ringpb.ShardHeartbeat{
		PrimaryDonId: 1,
		Timestamp:    12345,
	}
	payload, err := proto.Marshal(hb)
	require.NoError(t, err)

	body := &remotetypes.MessageBody{
		Method:  remotetypes.MethodShardHeartbeat,
		Payload: payload,
		Sender:  primary.Members[0][:],
	}

	rcvr.Receive(ctx, body)
	assert.NotEqual(t, int64(0), rcvr.LastSeen())
}
