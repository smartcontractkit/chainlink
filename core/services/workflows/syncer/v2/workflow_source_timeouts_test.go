package v2

import (
	"context"
	"net"
	"strconv"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	jwtmocks "github.com/smartcontractkit/chainlink-common/pkg/nodeauth/jwt/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/storage"
	storagepb "github.com/smartcontractkit/chainlink-protos/storage-service/go"
	sourcepb "github.com/smartcontractkit/chainlink-protos/workflows/go/sources"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	wftypes "github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

//nolint:paralleltest // Subtests exercise the same clients and server sequentially.
func TestWorkflowSourceTimeouts_GRPC(t *testing.T) {
	t.Parallel()
	var stalled atomic.Bool
	entered := make(chan context.Context, 1)
	server := grpc.NewServer(grpc.UnaryInterceptor(func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, _ grpc.UnaryHandler) (any, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		assert.Equal(t, []string{"Bearer test-token"}, md.Get("authorization"))
		if stalled.Load() {
			entered <- ctx
			<-ctx.Done()
			return nil, status.FromContextError(ctx.Err()).Err()
		}
		switch req.(type) {
		case *sourcepb.ListWorkflowMetadataRequest:
			return &sourcepb.ListWorkflowMetadataResponse{Workflows: []*sourcepb.WorkflowMetadata{createTestProtoWorkflow("workflow-1", "family-a")}}, nil
		default:
			return &storagepb.DownloadArtifactResponse{Url: "https://example.com/signed"}, nil
		}
	}))
	sourcepb.RegisterWorkflowMetadataSourceServiceServer(server, &sourcepb.UnimplementedWorkflowMetadataSourceServiceServer{})
	storagepb.RegisterNodeServiceServer(server, &storagepb.UnimplementedNodeServiceServer{})
	listener, err := net.Listen("tcp", "0.0.0.0:0")
	require.NoError(t, err)
	go func() { assert.NoError(t, server.Serve(listener)) }()
	t.Cleanup(server.Stop)
	address := net.JoinHostPort("127.0.0.1", strconv.Itoa(listener.Addr().(*net.TCPAddr).Port))

	jwt := jwtmocks.NewJWTGenerator(t)
	jwt.EXPECT().CreateJWTForRequest(mock.Anything).Return("test-token", nil)
	cfg := grpcTestSourceConfig("test-source")
	cfg.URL, cfg.JWTGenerator = address, jwt
	source, err := NewGRPCWorkflowSource(logger.TestLogger(t), cfg)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, source.Close()) })
	client, err := storage.NewWorkflowClient(logger.TestLogger(t), address, storage.WithJWTGenerator(jwt))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	fetcher := NewFetcherService(logger.TestLogger(t), nil, client)

	for name, call := range map[string]func(context.Context) error{
		"listing": func(ctx context.Context) error {
			views, _, err := source.ListWorkflowMetadata(ctx, capabilities.DON{})
			if err == nil {
				assert.Len(t, views, 1)
			}
			return err
		},
		"artifact URL": func(ctx context.Context) error {
			url, err := fetcher.RetrieveURL(ctx, &storagepb.DownloadArtifactRequest{Id: "workflow-id", Type: storagepb.ArtifactType_ARTIFACT_TYPE_CONFIG})
			if err == nil {
				assert.Equal(t, "https://example.com/signed", url)
			}
			return err
		},
	} {
		t.Run(name, func(t *testing.T) {
			for _, callerDeadline := range []bool{false, true} {
				stalled.Store(true)
				ctx, cancel := context.WithCancel(t.Context())
				t.Cleanup(cancel)
				budget := 20 * time.Second
				if callerDeadline {
					budget = time.Second
					ctx, cancel = context.WithTimeout(ctx, budget)
					t.Cleanup(cancel)
				}
				result := make(chan error, 1)
				go func() { result <- call(ctx) }()
				select {
				case serverCtx := <-entered:
					deadline, ok := serverCtx.Deadline()
					require.True(t, ok, "deadline must reach the server through JWT injection")
					require.InDelta(t, budget.Seconds(), time.Until(deadline).Seconds(), 0.5)
				case <-time.After(5 * time.Second):
					t.Fatal("RPC did not reach server")
				}
				wantCode := codes.DeadlineExceeded
				if !callerDeadline {
					cancel()
					wantCode = codes.Canceled
				}
				require.Equal(t, wantCode, status.Code(<-result))
				stalled.Store(false)
				require.NoError(t, call(t.Context()), "same client must recover on the next operation")
			}
		})
	}
}

func TestWorkflowSourceTimeouts_Reconciliation(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		wf := createTestProtoWorkflow("workflow-1", "family-a")
		client := &mockGRPCClient{allWorkflows: []*sourcepb.WorkflowMetadata{wf}}
		client.listFn = func(ctx context.Context, _ []string, offset, _ int64) ([]*sourcepb.WorkflowMetadata, bool, error) {
			if offset == 0 {
				return []*sourcepb.WorkflowMetadata{wf}, true, nil
			}
			<-ctx.Done()
			return nil, false, ctx.Err()
		}
		source, err := NewGRPCWorkflowSourceWithClient(logger.TestLogger(t), client, grpcTestSourceConfig("test-source"))
		require.NoError(t, err)
		liveID := wftypes.WorkflowID{2}
		orphanID := wftypes.WorkflowID{3}
		er := NewEngineRegistry()
		require.NoError(t, er.Add(liveID, source.SourceIdentifier(), &mockService{}))
		h := &orphanSweepFakeHandler{specs: []*job.WorkflowSpec{{WorkflowID: orphanID.Hex()}}}
		wr, tick := newSweepLoopRegistry(t, h, er, []WorkflowMetadataSource{
			source,
			&fakeMetadataSource{name: "healthy", metadata: []WorkflowMetadataView{{
				WorkflowID: wftypes.WorkflowID{4}, Status: WorkflowStatusActive, Source: "healthy",
			}}},
		})
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		done := make(chan struct{})
		go func() {
			defer close(done)
			wr.syncUsingReconciliationStrategy(ctx)
		}()
		tick <- time.Now()
		time.Sleep(20 * time.Second)
		synctest.Wait()
		require.Empty(t, h.OrphanDeletes(), "timeout must preserve engines and persisted specs")
		handled := h.Handled()
		require.Len(t, handled, 1, "later sources must still reconcile")
		require.Equal(t, WorkflowActivated, handled[0].Name)

		client.listFn = nil
		tick <- time.Now()
		synctest.Wait()
		require.ElementsMatch(t, []string{liveID.Hex(), orphanID.Hex()}, h.OrphanDeletes(), "complete listing permits deletions again")
		cancel()
		<-done
	})
}
