package ring

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
)

type mockArbiterScaler struct {
	called  bool
	nShards uint32
	err     error
}

func (m *mockArbiterScaler) Status(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ringpb.ReplicaStatus, error) {
	return &ringpb.ReplicaStatus{}, nil
}

func (m *mockArbiterScaler) ConsensusWantShards(ctx context.Context, req *ringpb.ConsensusWantShardsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.called = true
	m.nShards = req.NShards
	if m.err != nil {
		return nil, m.err
	}
	return &emptypb.Empty{}, nil
}

func TestTransmitter_NewTransmitter(t *testing.T) {
	lggr := logger.Test(t)
	store := NewStore()
	tx := NewTransmitter(lggr, store, nil, "test-account")
	require.NotNil(t, tx)
}

func TestTransmitter_FromAccount(t *testing.T) {
	lggr := logger.Test(t)
	store := NewStore()
	tx := NewTransmitter(lggr, store, nil, "my-account")

	account, err := tx.FromAccount(t.Context())
	require.NoError(t, err)
	require.Equal(t, types.Account("my-account"), account)
}

func TestTransmitter_Transmit(t *testing.T) {
	lggr := logger.Test(t)
	store := NewStore()
	mock := &mockArbiterScaler{}
	tx := NewTransmitter(lggr, store, mock, "test-account")

	outcome := &ringpb.Outcome{
		State: &ringpb.RoutingState{
			Id:    1,
			State: &ringpb.RoutingState_RoutableShards{RoutableShards: 3},
		},
		Routes: map[string]*ringpb.WorkflowRoute{
			"wf-1": {Shard: 0},
			"wf-2": {Shard: 1},
		},
	}
	outcomeBytes, err := proto.Marshal(outcome)
	require.NoError(t, err)

	report := ocr3types.ReportWithInfo[[]byte]{Report: outcomeBytes}
	err = tx.Transmit(t.Context(), types.ConfigDigest{}, 0, report, nil)
	require.NoError(t, err)

	// Verify arbiter was notified
	require.True(t, mock.called)
	require.Equal(t, uint32(3), mock.nShards)

	// Verify store was updated
	require.Equal(t, uint32(3), store.GetRoutingState().GetRoutableShards())
	routes := store.GetAllRoutingState()
	require.Equal(t, uint32(0), routes["wf-1"])
	require.Equal(t, uint32(1), routes["wf-2"])
}

func TestTransmitter_Transmit_NilArbiter(t *testing.T) {
	lggr := logger.Test(t)
	store := NewStore()
	tx := NewTransmitter(lggr, store, nil, "test-account")

	outcome := &ringpb.Outcome{
		State: &ringpb.RoutingState{
			Id:    1,
			State: &ringpb.RoutingState_RoutableShards{RoutableShards: 2},
		},
		Routes: map[string]*ringpb.WorkflowRoute{"wf-1": {Shard: 0}},
	}
	outcomeBytes, _ := proto.Marshal(outcome)

	err := tx.Transmit(t.Context(), types.ConfigDigest{}, 0, ocr3types.ReportWithInfo[[]byte]{Report: outcomeBytes}, nil)
	require.NoError(t, err)
}

func TestTransmitter_Transmit_TransitionState(t *testing.T) {
	lggr := logger.Test(t)
	store := NewStore()
	mock := &mockArbiterScaler{}
	tx := NewTransmitter(lggr, store, mock, "test-account")

	outcome := &ringpb.Outcome{
		State: &ringpb.RoutingState{
			Id: 1,
			State: &ringpb.RoutingState_Transition{
				Transition: &ringpb.Transition{WantShards: 5},
			},
		},
	}
	outcomeBytes, _ := proto.Marshal(outcome)

	err := tx.Transmit(t.Context(), types.ConfigDigest{}, 0, ocr3types.ReportWithInfo[[]byte]{Report: outcomeBytes}, nil)
	require.NoError(t, err)
	require.Equal(t, uint32(5), mock.nShards)
}

func TestTransmitter_Transmit_InvalidReport(t *testing.T) {
	lggr := logger.Test(t)
	store := NewStore()
	tx := NewTransmitter(lggr, store, nil, "test-account")

	// Send invalid protobuf data
	report := ocr3types.ReportWithInfo[[]byte]{Report: []byte("invalid protobuf")}
	err := tx.Transmit(t.Context(), types.ConfigDigest{}, 0, report, nil)
	require.Error(t, err)
}

func TestTransmitter_Transmit_ArbiterError(t *testing.T) {
	lggr := logger.Test(t)
	store := NewStore()
	mock := &mockArbiterScaler{err: context.DeadlineExceeded}
	tx := NewTransmitter(lggr, store, mock, "test-account")

	outcome := &ringpb.Outcome{
		State: &ringpb.RoutingState{
			Id:    1,
			State: &ringpb.RoutingState_RoutableShards{RoutableShards: 3},
		},
	}
	outcomeBytes, _ := proto.Marshal(outcome)

	err := tx.Transmit(t.Context(), types.ConfigDigest{}, 0, ocr3types.ReportWithInfo[[]byte]{Report: outcomeBytes}, nil)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestTransmitter_Transmit_StaleEntriesPruned(t *testing.T) {
	lggr := logger.Test(t)
	store := NewStore()
	mock := &mockArbiterScaler{}
	tx := NewTransmitter(lggr, store, mock, "test-account")

	outcome1 := &ringpb.Outcome{
		State: &ringpb.RoutingState{
			Id:    1,
			State: &ringpb.RoutingState_RoutableShards{RoutableShards: 3},
		},
		Routes: map[string]*ringpb.WorkflowRoute{
			"wf-1": {Shard: 0},
			"wf-2": {Shard: 1},
			"wf-3": {Shard: 2},
		},
	}
	outcomeBytes, err := proto.Marshal(outcome1)
	require.NoError(t, err)

	err = tx.Transmit(t.Context(), types.ConfigDigest{}, 0, ocr3types.ReportWithInfo[[]byte]{Report: outcomeBytes}, nil)
	require.NoError(t, err)
	require.Len(t, store.GetAllRoutingState(), 3)

	outcome2 := &ringpb.Outcome{
		State: &ringpb.RoutingState{
			Id:    2,
			State: &ringpb.RoutingState_RoutableShards{RoutableShards: 3},
		},
		Routes: map[string]*ringpb.WorkflowRoute{
			"wf-1": {Shard: 0},
		},
	}
	outcomeBytes, err = proto.Marshal(outcome2)
	require.NoError(t, err)

	err = tx.Transmit(t.Context(), types.ConfigDigest{}, 0, ocr3types.ReportWithInfo[[]byte]{Report: outcomeBytes}, nil)
	require.NoError(t, err)

	routes := store.GetAllRoutingState()
	require.Len(t, routes, 1)
	require.Equal(t, uint32(0), routes["wf-1"])
	require.NotContains(t, routes, "wf-2")
	require.NotContains(t, routes, "wf-3")
}

func TestTransmitter_Transmit_NilState(t *testing.T) {
	lggr := logger.Test(t)
	store := NewStore()
	tx := NewTransmitter(lggr, store, nil, "test-account")

	outcome := &ringpb.Outcome{
		State:  nil,
		Routes: map[string]*ringpb.WorkflowRoute{"wf-1": {Shard: 0}},
	}
	outcomeBytes, _ := proto.Marshal(outcome)

	err := tx.Transmit(t.Context(), types.ConfigDigest{}, 0, ocr3types.ReportWithInfo[[]byte]{Report: outcomeBytes}, nil)
	require.NoError(t, err)

	// Routes should still be applied
	routes := store.GetAllRoutingState()
	require.Equal(t, uint32(0), routes["wf-1"])
}
