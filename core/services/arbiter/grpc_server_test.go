package arbiter

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// mockShardConfigReaderForGRPC is a mock implementation of ShardConfigReader for testing.
type mockShardConfigReaderForGRPC struct {
	services.StateMachine
	shardCount uint64
	err        error
}

func (m *mockShardConfigReaderForGRPC) Start(ctx context.Context) error {
	return nil
}

func (m *mockShardConfigReaderForGRPC) Close() error {
	return nil
}

func (m *mockShardConfigReaderForGRPC) Ready() error {
	return nil
}

func (m *mockShardConfigReaderForGRPC) HealthReport() map[string]error {
	return nil
}

func (m *mockShardConfigReaderForGRPC) Name() string {
	return "mockShardConfigReaderForGRPC"
}

func (m *mockShardConfigReaderForGRPC) GetDesiredShardCount(ctx context.Context) (uint64, error) {
	return m.shardCount, m.err
}

func TestGRPCServer_GetDesiredReplicas_Success(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockShardConfigReaderForGRPC{shardCount: 5}
	state := NewState()

	server := NewGRPCServer(mockReader, state, lggr)

	req := &ringpb.ShardStatusRequest{
		Status: map[uint32]*ringpb.ShardStatus{
			0: {IsHealthy: true},
			1: {IsHealthy: true},
			2: {IsHealthy: false},
		},
	}

	resp, err := server.GetDesiredReplicas(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, uint32(5), resp.WantShards)
}

func TestGRPCServer_GetDesiredReplicas_EmptyRequest(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockShardConfigReaderForGRPC{shardCount: 3}
	state := NewState()

	server := NewGRPCServer(mockReader, state, lggr)

	req := &ringpb.ShardStatusRequest{
		Status: map[uint32]*ringpb.ShardStatus{},
	}

	resp, err := server.GetDesiredReplicas(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, uint32(3), resp.WantShards)
}

func TestGRPCServer_GetDesiredReplicas_ShardConfigError(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockShardConfigReaderForGRPC{
		err: errors.New("contract read failed"),
	}
	state := NewState()

	server := NewGRPCServer(mockReader, state, lggr)

	req := &ringpb.ShardStatusRequest{
		Status: map[uint32]*ringpb.ShardStatus{
			0: {IsHealthy: true},
		},
	}

	resp, err := server.GetDesiredReplicas(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestGRPCServer_GetDesiredReplicas_ZeroShards(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockShardConfigReaderForGRPC{shardCount: 0}
	state := NewState()

	server := NewGRPCServer(mockReader, state, lggr)

	req := &ringpb.ShardStatusRequest{
		Status: map[uint32]*ringpb.ShardStatus{},
	}

	resp, err := server.GetDesiredReplicas(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, uint32(0), resp.WantShards)
}

func TestGRPCServer_GetDesiredReplicas_LargeShardCount(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockShardConfigReaderForGRPC{shardCount: 100}
	state := NewState()

	server := NewGRPCServer(mockReader, state, lggr)

	// Simulate many healthy shards
	status := make(map[uint32]*ringpb.ShardStatus)
	for i := uint32(0); i < 100; i++ {
		status[i] = &ringpb.ShardStatus{IsHealthy: true}
	}

	req := &ringpb.ShardStatusRequest{
		Status: status,
	}

	resp, err := server.GetDesiredReplicas(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, uint32(100), resp.WantShards)
}
