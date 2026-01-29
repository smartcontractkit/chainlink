package shardorchestrator

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
)

const bufSize = 1024 * 1024

// mockShardOrchestratorServer implements the gRPC server for testing
type mockShardOrchestratorServer struct {
	ringpb.UnimplementedShardOrchestratorServiceServer
	mappings           map[string]uint32
	registrationCalled bool
}

func (m *mockShardOrchestratorServer) GetWorkflowShardMapping(ctx context.Context, req *ringpb.GetWorkflowShardMappingRequest) (*ringpb.GetWorkflowShardMappingResponse, error) {
	mappings := make(map[string]uint32)
	mappingStates := make(map[string]*ringpb.WorkflowMappingState)

	for _, wfID := range req.WorkflowIds {
		if shardID, ok := m.mappings[wfID]; ok {
			mappings[wfID] = shardID
			mappingStates[wfID] = &ringpb.WorkflowMappingState{
				OldShardId:   0,
				NewShardId:   shardID,
				InTransition: false,
			}
		}
	}

	return &ringpb.GetWorkflowShardMappingResponse{
		Mappings:       mappings,
		MappingStates:  mappingStates,
		MappingVersion: 1,
	}, nil
}

func (m *mockShardOrchestratorServer) ReportWorkflowTriggerRegistration(ctx context.Context, req *ringpb.ReportWorkflowTriggerRegistrationRequest) (*ringpb.ReportWorkflowTriggerRegistrationResponse, error) {
	m.registrationCalled = true
	return &ringpb.ReportWorkflowTriggerRegistrationResponse{
		Success: true,
	}, nil
}

// setupTestServer creates a test gRPC server using bufconn
func setupTestServer(t *testing.T, mock *mockShardOrchestratorServer) (*grpc.Server, *bufconn.Listener) {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	ringpb.RegisterShardOrchestratorServiceServer(s, mock)

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Logf("Server exited with error: %v", err)
		}
	}()

	return s, lis
}

// createTestClient creates a client connected to the test server
func createTestClient(t *testing.T, lis *bufconn.Listener) *Client {
	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	lggr := logger.Test(t)
	return &Client{
		conn:   conn,
		client: ringpb.NewShardOrchestratorServiceClient(conn),
		logger: logger.Named(lggr, "TestClient"),
	}
}

func TestClient_GetWorkflowShardMapping(t *testing.T) {
	ctx := context.Background()

	mock := &mockShardOrchestratorServer{
		mappings: map[string]uint32{
			"workflow-1": 0,
			"workflow-2": 1,
			"workflow-3": 2,
		},
	}

	grpcServer, lis := setupTestServer(t, mock)
	defer grpcServer.Stop()

	client := createTestClient(t, lis)
	defer client.Close()

	t.Run("successful mapping query", func(t *testing.T) {
		workflowIDs := []string{"workflow-1", "workflow-2", "workflow-3"}
		resp, err := client.GetWorkflowShardMapping(ctx, workflowIDs)
		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Len(t, resp.Mappings, 3)
		assert.Equal(t, uint32(0), resp.Mappings["workflow-1"])
		assert.Equal(t, uint32(1), resp.Mappings["workflow-2"])
		assert.Equal(t, uint32(2), resp.Mappings["workflow-3"])

		assert.Len(t, resp.MappingStates, 3)
		assert.Equal(t, uint64(1), resp.MappingVersion)
	})

	t.Run("partial workflow query", func(t *testing.T) {
		workflowIDs := []string{"workflow-1", "workflow-unknown"}
		resp, err := client.GetWorkflowShardMapping(ctx, workflowIDs)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Should only return mappings for known workflows
		assert.Len(t, resp.Mappings, 1)
		assert.Equal(t, uint32(0), resp.Mappings["workflow-1"])
		_, exists := resp.Mappings["workflow-unknown"]
		assert.False(t, exists)
	})

	t.Run("empty workflow list", func(t *testing.T) {
		resp, err := client.GetWorkflowShardMapping(ctx, []string{})
		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Empty(t, resp.Mappings)
	})
}

func TestClient_ReportWorkflowTriggerRegistration(t *testing.T) {
	ctx := context.Background()

	mock := &mockShardOrchestratorServer{
		mappings: map[string]uint32{},
	}

	grpcServer, lis := setupTestServer(t, mock)
	defer grpcServer.Stop()

	client := createTestClient(t, lis)
	defer client.Close()

	t.Run("successful registration report", func(t *testing.T) {
		req := &ringpb.ReportWorkflowTriggerRegistrationRequest{
			SourceShardId: 1,
			RegisteredWorkflows: map[string]uint32{
				"workflow-1": 1,
				"workflow-2": 1,
			},
			TotalActiveWorkflows: 2,
		}

		resp, err := client.ReportWorkflowTriggerRegistration(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.True(t, resp.Success)
		assert.True(t, mock.registrationCalled)
	})
}

func TestClient_Close(t *testing.T) {
	mock := &mockShardOrchestratorServer{
		mappings: map[string]uint32{},
	}

	grpcServer, lis := setupTestServer(t, mock)
	defer grpcServer.Stop()

	client := createTestClient(t, lis)

	err := client.Close()
	require.NoError(t, err)

	// Verify connection is closed by attempting to use it
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err = client.GetWorkflowShardMapping(ctx, []string{"test"})
	assert.Error(t, err, "should fail after client is closed")
}

func TestNewClient(t *testing.T) {
	ctx := context.Background()
	lggr := logger.Test(t)

	t.Run("creates client successfully", func(t *testing.T) {
		// Note: This creates a client but doesn't connect immediately with grpc.NewClient
		client, err := NewClient(ctx, "localhost:50051", lggr)
		require.NoError(t, err)
		require.NotNil(t, client)
		defer client.Close()

		assert.NotNil(t, client.conn)
		assert.NotNil(t, client.client)
		assert.NotNil(t, client.logger)
	})
}
