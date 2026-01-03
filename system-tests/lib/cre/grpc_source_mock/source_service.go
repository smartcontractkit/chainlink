package grpc_source_mock

import (
	"context"
	"crypto/sha256"
	"log/slog"
	"os"
	"strconv"
	"time"

	sourcesv1 "github.com/smartcontractkit/chainlink-protos/workflows/go/sources"
)

var sourceLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug,
})).With("logger", "grpc_source_mock.SourceService")

// SourceService implements the WorkflowMetadataSourceService gRPC service
type SourceService struct {
	sourcesv1.UnimplementedWorkflowMetadataSourceServiceServer
	store *WorkflowStore
}

// NewSourceService creates a new SourceService
func NewSourceService(store *WorkflowStore) *SourceService {
	return &SourceService{
		store: store,
	}
}

// ListWorkflowMetadata returns all workflow metadata for the given DON
func (s *SourceService) ListWorkflowMetadata(ctx context.Context, req *sourcesv1.ListWorkflowMetadataRequest) (*sourcesv1.ListWorkflowMetadataResponse, error) {
	sourceLogger.Debug("ListWorkflowMetadata called",
		"donFamilies", req.GetDonFamilies(),
		"start", req.GetStart(),
		"limit", req.GetLimit(),
	)

	// Get all workflows matching the filter
	workflows := s.store.List(req.GetDonFamilies())

	sourceLogger.Debug("ListWorkflowMetadata results",
		"donFamiliesFilter", req.GetDonFamilies(),
		"workflowCount", len(workflows),
	)

	// Apply pagination
	start := req.GetStart()
	limit := req.GetLimit()
	if limit == 0 {
		limit = 1000 // default limit
	}

	// Calculate pagination bounds
	totalCount := int64(len(workflows))
	if start >= totalCount {
		// No results for this page
		return &sourcesv1.ListWorkflowMetadataResponse{
			Workflows: []*sourcesv1.WorkflowMetadata{},
			Head:      s.createHead(),
			HasMore:   false,
		}, nil
	}

	end := min(start+limit, totalCount)

	// Convert to proto messages
	protoWorkflows := make([]*sourcesv1.WorkflowMetadata, 0, end-start)
	for i := start; i < end; i++ {
		wf := workflows[i]
		protoWorkflows = append(protoWorkflows, &sourcesv1.WorkflowMetadata{
			WorkflowId:   wf.Registration.WorkflowID[:],
			Owner:        wf.Registration.Owner,
			CreatedAt:    uint64(wf.CreatedAt), // Convert millisecond timestamp to uint64
			Status:       uint32(wf.Status),
			WorkflowName: wf.Registration.WorkflowName,
			BinaryUrl:    wf.Registration.BinaryURL,
			ConfigUrl:    wf.Registration.ConfigURL,
			Tag:          wf.Registration.Tag,
			Attributes:   wf.Registration.Attributes,
			DonFamily:    wf.Registration.DonFamily,
		})
	}

	return &sourcesv1.ListWorkflowMetadataResponse{
		Workflows: protoWorkflows,
		Head:      s.createHead(),
		HasMore:   end < totalCount,
	}, nil
}

// createHead creates a synthetic head for the response
func (s *SourceService) createHead() *sourcesv1.Head {
	now := time.Now()
	height := strconv.FormatInt(now.UnixNano(), 10)
	hash := sha256.Sum256([]byte(height))

	return &sourcesv1.Head{
		Height:    height,
		Hash:      hash[:],
		Timestamp: uint64(now.Unix()),
	}
}
