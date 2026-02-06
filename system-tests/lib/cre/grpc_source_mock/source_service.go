package grpcsourcemock

import (
	"context"
	"encoding/hex"
	"log/slog"
	"os"

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
			HasMore:   false,
		}, nil
	}

	end := min(start+limit, totalCount)

	// Convert to proto messages
	protoWorkflows := make([]*sourcesv1.WorkflowMetadata, 0, end-start)
	for i := start; i < end; i++ {
		wf := workflows[i]
		var createdAt uint64
		if wf.CreatedAt >= 0 {
			createdAt = uint64(wf.CreatedAt) // #nosec G115 -- CreatedAt is always positive timestamp
		}
		protoWorkflows = append(protoWorkflows, &sourcesv1.WorkflowMetadata{
			WorkflowId:   wf.Registration.WorkflowID[:],
			Owner:        hex.EncodeToString(wf.Registration.Owner),
			CreatedAt:    createdAt,
			Status:       workflowStatusToProto(wf.Status),
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
		HasMore:   end < totalCount,
	}, nil
}

// workflowStatusToProto converts store WorkflowStatus to proto WorkflowStatus.
// Store uses: Active=0, Paused=1
// Proto uses: UNSPECIFIED=0, ACTIVE=1, PAUSED=2
func workflowStatusToProto(status WorkflowStatus) sourcesv1.WorkflowStatus {
	switch status {
	case WorkflowStatusActive:
		return sourcesv1.WorkflowStatus_WORKFLOW_STATUS_ACTIVE
	case WorkflowStatusPaused:
		return sourcesv1.WorkflowStatus_WORKFLOW_STATUS_PAUSED
	default:
		return sourcesv1.WorkflowStatus_WORKFLOW_STATUS_UNSPECIFIED
	}
}
