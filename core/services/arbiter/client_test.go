package arbiter

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"

	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// mockArbiterScalerServer implements ringpb.ArbiterScalerServer for testing.
type mockArbiterScalerServer struct {
	ringpb.UnimplementedArbiterScalerServer
	statusResp      *ringpb.ReplicaStatus
	statusErr       error
	consensusErr    error
	consensusCalled bool
	lastNShards     uint32
}

func (m *mockArbiterScalerServer) Status(ctx context.Context, _ *emptypb.Empty) (*ringpb.ReplicaStatus, error) {
	if m.statusErr != nil {
		return nil, m.statusErr
	}
	return m.statusResp, nil
}

func (m *mockArbiterScalerServer) ConsensusWantShards(ctx context.Context, req *ringpb.ConsensusWantShardsRequest) (*emptypb.Empty, error) {
	m.consensusCalled = true
	m.lastNShards = req.GetNShards()
	if m.consensusErr != nil {
		return nil, m.consensusErr
	}
	return &emptypb.Empty{}, nil
}

func TestRingArbiterClient_Status(t *testing.T) {
	lggr := logger.TestLogger(t)

	t.Run("returns status from server", func(t *testing.T) {
		mockServer := &mockArbiterScalerServer{
			statusResp: &ringpb.ReplicaStatus{
				WantShards: 5,
				Status: map[uint32]*ringpb.ShardStatus{
					0: {IsHealthy: true},
					1: {IsHealthy: true},
					2: {IsHealthy: false},
				},
			},
		}

		client := NewRingArbiterClient(mockServer, lggr)
		resp, err := client.Status(context.Background(), &emptypb.Empty{})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, uint32(5), resp.WantShards)
		assert.Len(t, resp.Status, 3)
		assert.True(t, resp.Status[0].IsHealthy)
		assert.True(t, resp.Status[1].IsHealthy)
		assert.False(t, resp.Status[2].IsHealthy)
	})

	t.Run("returns error from server", func(t *testing.T) {
		expectedErr := errors.New("server error")
		mockServer := &mockArbiterScalerServer{
			statusErr: expectedErr,
		}

		client := NewRingArbiterClient(mockServer, lggr)
		resp, err := client.Status(context.Background(), &emptypb.Empty{})

		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, expectedErr, err)
	})
}

func TestRingArbiterClient_ConsensusWantShards(t *testing.T) {
	lggr := logger.TestLogger(t)

	t.Run("calls server with correct request", func(t *testing.T) {
		mockServer := &mockArbiterScalerServer{}

		client := NewRingArbiterClient(mockServer, lggr)
		req := &ringpb.ConsensusWantShardsRequest{NShards: 10}
		resp, err := client.ConsensusWantShards(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.True(t, mockServer.consensusCalled)
		assert.Equal(t, uint32(10), mockServer.lastNShards)
	})

	t.Run("returns error from server", func(t *testing.T) {
		expectedErr := errors.New("consensus error")
		mockServer := &mockArbiterScalerServer{
			consensusErr: expectedErr,
		}

		client := NewRingArbiterClient(mockServer, lggr)
		req := &ringpb.ConsensusWantShardsRequest{NShards: 5}
		resp, err := client.ConsensusWantShards(context.Background(), req)

		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, expectedErr, err)
	})
}
