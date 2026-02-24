package v2

import (
	"context"
	"encoding/hex"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	pb "github.com/smartcontractkit/chainlink-protos/workflows/go/sources"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Test constants for workflow metadata
const (
	grpcTestOwnerHex  = "0102030405060708091011121314151617181920"
	grpcTestBinaryURL = "https://example.com/binary.wasm"
	grpcTestConfigURL = "https://example.com/config.json"
)

// grpcTestBinaryContent and grpcTestConfigContent are mock content used for canonical workflowID calculation
var (
	grpcTestBinaryContent = []byte("mock-wasm-binary-content")
	grpcTestConfigContent = []byte("{}")
)

// mockGRPCClient is a mock implementation of grpcClient for testing.
type mockGRPCClient struct {
	// allWorkflows contains all workflows to be returned (used for stateless pagination)
	allWorkflows []*pb.WorkflowMetadata
	// err is the error to return (if set, takes precedence)
	err error
	// errSequence allows returning different errors on successive calls (for retry testing)
	errSequence []error
	// callCount tracks how many times ListWorkflowMetadata was called
	callCount atomic.Int32
	// closed tracks if Close was called
	closed bool
	// closeErr is the error to return from Close
	closeErr error
}

func (m *mockGRPCClient) ListWorkflowMetadata(_ context.Context, _ []string, offset, limit int64) ([]*pb.WorkflowMetadata, bool, error) {
	callNum := int(m.callCount.Add(1)) - 1 // 0-indexed call number

	// Check if there's a specific error for this call number
	if callNum < len(m.errSequence) && m.errSequence[callNum] != nil {
		return nil, false, m.errSequence[callNum]
	}

	// Check for general error
	if m.err != nil {
		return nil, false, m.err
	}

	// Stateless pagination based on offset/limit
	start := int(offset)
	if start >= len(m.allWorkflows) {
		return []*pb.WorkflowMetadata{}, false, nil
	}

	end := start + int(limit)
	if end > len(m.allWorkflows) {
		end = len(m.allWorkflows)
	}

	hasMore := end < len(m.allWorkflows)
	return m.allWorkflows[start:end], hasMore, nil
}

func (m *mockGRPCClient) Close() error {
	m.closed = true
	return m.closeErr
}

// CallCount returns the number of times ListWorkflowMetadata was called
func (m *mockGRPCClient) CallCount() int {
	return int(m.callCount.Load())
}

// createTestProtoWorkflow creates a test protobuf WorkflowMetadata for testing.
// It uses the canonical workflow ID calculation to ensure test data is realistic.
func createTestProtoWorkflow(name string, family string) *pb.WorkflowMetadata {
	owner, err := hex.DecodeString(grpcTestOwnerHex)
	if err != nil {
		panic("failed to decode owner hex: " + err.Error())
	}

	// Use canonical workflow ID calculation
	workflowID, err := workflows.GenerateWorkflowID(owner, name, grpcTestBinaryContent, grpcTestConfigContent, "")
	if err != nil {
		panic("failed to generate workflow ID: " + err.Error())
	}

	return &pb.WorkflowMetadata{
		WorkflowId:   workflowID[:],
		Owner:        grpcTestOwnerHex, // Proto uses hex string, not bytes
		CreatedAt:    1234567890,
		Status:       pb.WorkflowStatus_WORKFLOW_STATUS_ACTIVE,
		WorkflowName: name,
		BinaryUrl:    grpcTestBinaryURL,
		ConfigUrl:    grpcTestConfigURL,
		Tag:          "v1.0.0",
		Attributes:   []byte("{}"),
		DonFamily:    family,
	}
}

func TestGRPCWorkflowSource_NewGRPCWorkflowSource_EmptyURL(t *testing.T) {
	lggr := logger.TestLogger(t)

	_, err := NewGRPCWorkflowSource(lggr, GRPCWorkflowSourceConfig{
		Name: "test-source",
		URL:  "",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "GRPC URL is required")
}

func TestGRPCWorkflowSource_NewGRPCWorkflowSource_EmptyName(t *testing.T) {
	lggr := logger.TestLogger(t)

	_, err := NewGRPCWorkflowSource(lggr, GRPCWorkflowSourceConfig{
		Name: "",
		URL:  "localhost:50051",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "source name is required")
}

func TestGRPCWorkflowSourceWithClient_EmptyName(t *testing.T) {
	lggr := logger.TestLogger(t)

	_, err := NewGRPCWorkflowSourceWithClient(lggr, &mockGRPCClient{}, GRPCWorkflowSourceConfig{
		Name: "",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "source name is required")
}

func TestGRPCWorkflowSource_ListWorkflowMetadata_Success(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	mockClient := &mockGRPCClient{
		allWorkflows: []*pb.WorkflowMetadata{
			createTestProtoWorkflow("workflow-1", "family-a"),
			createTestProtoWorkflow("workflow-2", "family-a"),
		},
	}

	source, err := NewGRPCWorkflowSourceWithClient(lggr, mockClient, GRPCWorkflowSourceConfig{
		Name: "test-source",
	})
	require.NoError(t, err)

	don := capabilities.DON{
		ID:       1,
		Families: []string{"family-a"},
	}

	wfs, head, err := source.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Len(t, wfs, 2)
	require.NotNil(t, head)
	assert.NotEmpty(t, head.Height)
	assert.Equal(t, []byte("grpc-source"), head.Hash)
	assert.Equal(t, "workflow-1", wfs[0].WorkflowName)
	assert.Equal(t, "workflow-2", wfs[1].WorkflowName)
	assert.Equal(t, 1, mockClient.CallCount())
}

func TestGRPCWorkflowSource_ListWorkflowMetadata_Pagination(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	// Configure mock with all workflows - pagination is handled stateless via offset/limit
	mockClient := &mockGRPCClient{
		allWorkflows: []*pb.WorkflowMetadata{
			createTestProtoWorkflow("workflow-1", "family-a"),
			createTestProtoWorkflow("workflow-2", "family-a"),
			createTestProtoWorkflow("workflow-3", "family-a"),
		},
	}

	source, err := NewGRPCWorkflowSourceWithClient(lggr, mockClient, GRPCWorkflowSourceConfig{
		Name:     "test-source",
		PageSize: 2, // Small page size to test pagination
	})
	require.NoError(t, err)

	don := capabilities.DON{
		ID:       1,
		Families: []string{"family-a"},
	}

	wfs, head, err := source.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Len(t, wfs, 3) // 2 from first page + 1 from second page
	require.NotNil(t, head)
	assert.NotEmpty(t, head.Height)
	assert.Equal(t, 2, mockClient.CallCount()) // Two pages fetched
}

func TestGRPCWorkflowSource_ListWorkflowMetadata_InvalidWorkflow(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	// Create a workflow with invalid ID (not 32 bytes)
	invalidWorkflow := &pb.WorkflowMetadata{
		WorkflowId:   []byte{1, 2, 3}, // Invalid: only 3 bytes
		WorkflowName: "invalid-workflow",
	}

	mockClient := &mockGRPCClient{
		allWorkflows: []*pb.WorkflowMetadata{
			createTestProtoWorkflow("valid-workflow", "family-a"),
			invalidWorkflow,
		},
	}

	source, err := NewGRPCWorkflowSourceWithClient(lggr, mockClient, GRPCWorkflowSourceConfig{
		Name: "test-source",
	})
	require.NoError(t, err)

	don := capabilities.DON{
		ID:       1,
		Families: []string{"family-a"},
	}

	wfs, head, err := source.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Len(t, wfs, 1) // Only valid workflow is returned
	assert.Equal(t, "valid-workflow", wfs[0].WorkflowName)
	require.NotNil(t, head)
	assert.NotEmpty(t, head.Height)
}

func TestGRPCWorkflowSource_Retry_Unavailable(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	// Use errSequence to return errors on first two calls, then succeed
	mockClient := &mockGRPCClient{
		allWorkflows: []*pb.WorkflowMetadata{
			createTestProtoWorkflow("workflow-1", "family-a"),
		},
		errSequence: []error{
			status.Error(codes.Unavailable, "server unavailable"),
			status.Error(codes.Unavailable, "server unavailable"),
			nil, // Third call succeeds
		},
	}

	source, err := NewGRPCWorkflowSourceWithClient(lggr, mockClient, GRPCWorkflowSourceConfig{
		Name:           "test-source",
		MaxRetries:     2,
		RetryBaseDelay: 1 * time.Millisecond, // Fast retries for testing
		RetryMaxDelay:  10 * time.Millisecond,
	})
	require.NoError(t, err)

	don := capabilities.DON{
		ID:       1,
		Families: []string{"family-a"},
	}

	wfs, head, err := source.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Len(t, wfs, 1)
	require.NotNil(t, head)
	assert.NotEmpty(t, head.Height)
	assert.Equal(t, 3, mockClient.CallCount()) // 2 failures + 1 success
}

func TestGRPCWorkflowSource_Retry_ResourceExhausted(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	mockClient := &mockGRPCClient{
		allWorkflows: []*pb.WorkflowMetadata{
			createTestProtoWorkflow("workflow-1", "family-a"),
		},
		errSequence: []error{
			status.Error(codes.ResourceExhausted, "rate limited"),
			nil, // Second call succeeds
		},
	}

	source, err := NewGRPCWorkflowSourceWithClient(lggr, mockClient, GRPCWorkflowSourceConfig{
		Name:           "test-source",
		MaxRetries:     2,
		RetryBaseDelay: 1 * time.Millisecond,
	})
	require.NoError(t, err)

	don := capabilities.DON{
		ID:       1,
		Families: []string{"family-a"},
	}

	wfs, _, err := source.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Len(t, wfs, 1)
	assert.Equal(t, 2, mockClient.CallCount()) // 1 failure + 1 success
}

func TestGRPCWorkflowSource_Retry_MaxExceeded(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	// Always return unavailable error
	mockClient := &mockGRPCClient{
		err: status.Error(codes.Unavailable, "server unavailable"),
	}

	source, err := NewGRPCWorkflowSourceWithClient(lggr, mockClient, GRPCWorkflowSourceConfig{
		Name:           "test-source",
		MaxRetries:     2,
		RetryBaseDelay: 1 * time.Millisecond,
	})
	require.NoError(t, err)

	don := capabilities.DON{
		ID:       1,
		Families: []string{"family-a"},
	}

	_, _, err = source.ListWorkflowMetadata(ctx, don)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "max retries")
	assert.Equal(t, 3, mockClient.CallCount()) // 1 initial + 2 retries
}

func TestGRPCWorkflowSource_Retry_NonRetryable(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	mockClient := &mockGRPCClient{
		err: status.Error(codes.InvalidArgument, "bad request"),
	}

	source, err := NewGRPCWorkflowSourceWithClient(lggr, mockClient, GRPCWorkflowSourceConfig{
		Name:           "test-source",
		MaxRetries:     2,
		RetryBaseDelay: 1 * time.Millisecond,
	})
	require.NoError(t, err)

	don := capabilities.DON{
		ID:       1,
		Families: []string{"family-a"},
	}

	_, _, err = source.ListWorkflowMetadata(ctx, don)
	require.Error(t, err)
	assert.Equal(t, 1, mockClient.CallCount()) // No retries for non-retryable errors
}

func TestGRPCWorkflowSource_Backoff_Jitter(t *testing.T) {
	lggr := logger.TestLogger(t)

	source, err := NewGRPCWorkflowSourceWithClient(lggr, &mockGRPCClient{}, GRPCWorkflowSourceConfig{
		Name:           "test-source",
		RetryBaseDelay: 100 * time.Millisecond,
		RetryMaxDelay:  2 * time.Second,
	})
	require.NoError(t, err)

	// Test backoff calculation
	backoff1 := source.calculateBackoff(1)
	backoff2 := source.calculateBackoff(2)
	backoff3 := source.calculateBackoff(3)

	// Backoff should increase exponentially (with jitter)
	// Attempt 1: baseDelay * 2^0 * jitter = 100ms * 1 * [0.5, 1.5] = [50ms, 150ms]
	assert.GreaterOrEqual(t, backoff1, 50*time.Millisecond)
	assert.LessOrEqual(t, backoff1, 150*time.Millisecond)

	// Attempt 2: baseDelay * 2^1 * jitter = 100ms * 2 * [0.5, 1.5] = [100ms, 300ms]
	assert.GreaterOrEqual(t, backoff2, 100*time.Millisecond)
	assert.LessOrEqual(t, backoff2, 300*time.Millisecond)

	// Attempt 3: baseDelay * 2^2 * jitter = 100ms * 4 * [0.5, 1.5] = [200ms, 600ms]
	assert.GreaterOrEqual(t, backoff3, 200*time.Millisecond)
	assert.LessOrEqual(t, backoff3, 600*time.Millisecond)
}

func TestGRPCWorkflowSource_ContextCancellation(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx, cancel := context.WithCancel(t.Context())

	// Always return unavailable to trigger retries
	mockClient := &mockGRPCClient{
		err: status.Error(codes.Unavailable, "server unavailable"),
	}

	source, err := NewGRPCWorkflowSourceWithClient(lggr, mockClient, GRPCWorkflowSourceConfig{
		Name:           "test-source",
		MaxRetries:     5, // High retry count
		RetryBaseDelay: 100 * time.Millisecond,
	})
	require.NoError(t, err)

	don := capabilities.DON{
		ID:       1,
		Families: []string{"family-a"},
	}

	// Cancel context immediately after first call
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	_, _, err = source.ListWorkflowMetadata(ctx, don)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestGRPCWorkflowSource_ConfigDefaults(t *testing.T) {
	lggr := logger.TestLogger(t)

	// Name is required, but other config options have defaults
	source, err := NewGRPCWorkflowSourceWithClient(lggr, &mockGRPCClient{}, GRPCWorkflowSourceConfig{
		Name: "test-source",
	})
	require.NoError(t, err)

	// Verify defaults are applied for non-required fields
	assert.Equal(t, defaultPageSize, source.pageSize)
	assert.Equal(t, defaultMaxRetries, source.maxRetries)
	assert.Equal(t, defaultRetryBaseDelay, source.retryBaseDelay)
	assert.Equal(t, defaultRetryMaxDelay, source.retryMaxDelay)
	assert.Equal(t, "test-source", source.name) // Name from config
}

func TestGRPCWorkflowSource_Ready(t *testing.T) {
	lggr := logger.TestLogger(t)

	source, err := NewGRPCWorkflowSourceWithClient(lggr, &mockGRPCClient{}, GRPCWorkflowSourceConfig{
		Name: "test-source",
	})
	require.NoError(t, err)

	// Initially ready
	assert.NoError(t, source.Ready())

	// After close, not ready
	err = source.Close()
	require.NoError(t, err)
	assert.Error(t, source.Ready())
}

func TestGRPCWorkflowSource_ListWorkflowMetadata_NotReady(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	source, err := NewGRPCWorkflowSourceWithClient(lggr, &mockGRPCClient{}, GRPCWorkflowSourceConfig{
		Name: "test-source",
	})
	require.NoError(t, err)

	// Close the source to make it not ready
	err = source.Close()
	require.NoError(t, err)

	don := capabilities.DON{
		ID:       1,
		Families: []string{"family-a"},
	}

	_, _, err = source.ListWorkflowMetadata(ctx, don)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not ready")
}

func TestGRPCWorkflowSource_Close(t *testing.T) {
	lggr := logger.TestLogger(t)

	mockClient := &mockGRPCClient{}

	source, err := NewGRPCWorkflowSourceWithClient(lggr, mockClient, GRPCWorkflowSourceConfig{
		Name: "test-source",
	})
	require.NoError(t, err)

	// Initially ready
	assert.NoError(t, source.Ready())
	assert.False(t, mockClient.closed)

	// Close
	err = source.Close()
	require.NoError(t, err)

	// Now not ready and client is closed
	require.Error(t, source.Ready())
	assert.True(t, mockClient.closed)
}

func TestGRPCWorkflowSource_Name(t *testing.T) {
	lggr := logger.TestLogger(t)

	source, err := NewGRPCWorkflowSourceWithClient(lggr, &mockGRPCClient{}, GRPCWorkflowSourceConfig{
		Name: "my-custom-source",
	})
	require.NoError(t, err)

	assert.Equal(t, "my-custom-source", source.Name())
}

func TestGRPCWorkflowSource_Name_Required(t *testing.T) {
	lggr := logger.TestLogger(t)

	// Empty name should return an error
	_, err := NewGRPCWorkflowSourceWithClient(lggr, &mockGRPCClient{}, GRPCWorkflowSourceConfig{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "source name is required")
}

func TestGRPCWorkflowSource_syntheticHead(t *testing.T) {
	lggr := logger.TestLogger(t)

	source, err := NewGRPCWorkflowSourceWithClient(lggr, &mockGRPCClient{}, GRPCWorkflowSourceConfig{
		Name: "test-source",
	})
	require.NoError(t, err)

	head := source.syntheticHead()
	require.NotNil(t, head)
	// Should return synthetic head with current timestamp
	assert.NotEmpty(t, head.Height)
	assert.Equal(t, []byte("grpc-source"), head.Hash)
	assert.Positive(t, head.Timestamp)
}
