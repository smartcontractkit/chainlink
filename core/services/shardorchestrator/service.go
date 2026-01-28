package shardorchestrator

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
)

// Server implements the gRPC ShardOrchestratorService
// This runs on shard zero and serves requests from other shards
type Server struct {
	ringpb.UnimplementedShardOrchestratorServiceServer
	store  *Store
	logger logger.Logger
}

func NewServer(store *Store, lggr logger.Logger) *Server {
	return &Server{
		store:  store,
		logger: logger.Named(lggr, "ShardOrchestratorServer"),
	}
}

// RegisterWithGRPCServer registers this service with a gRPC server
func (s *Server) RegisterWithGRPCServer(grpcServer *grpc.Server) {
	ringpb.RegisterShardOrchestratorServiceServer(grpcServer, s)
	s.logger.Info("Registered ShardOrchestrator gRPC service")
}

// GetWorkflowShardMapping handles batch requests for workflow-to-shard mappings
// This is called by other shards to determine where to route workflow executions
func (s *Server) GetWorkflowShardMapping(ctx context.Context, req *ringpb.GetWorkflowShardMappingRequest) (*ringpb.GetWorkflowShardMappingResponse, error) {
	s.logger.Debugw("GetWorkflowShardMapping called", "workflowCount", len(req.WorkflowIds))

	if len(req.WorkflowIds) == 0 {
		return nil, errors.New("workflow_ids is required and must not be empty")
	}

	// Retrieve batch from store
	mappings, version, err := s.store.GetWorkflowMappingsBatch(ctx, req.WorkflowIds)
	if err != nil {
		s.logger.Errorw("Failed to get workflow mappings", "error", err)
		return nil, fmt.Errorf("failed to get workflow mappings: %w", err)
	}

	// Build simple mappings map (workflow_id -> shard_id)
	simpleMappings := make(map[string]uint32, len(mappings))
	// Build detailed mapping states
	mappingStates := make(map[string]*ringpb.WorkflowMappingState, len(mappings))

	for workflowID, mapping := range mappings {
		// Simple mapping: just the current shard
		simpleMappings[workflowID] = mapping.NewShardID

		// Detailed state: includes transition information
		mappingStates[workflowID] = &ringpb.WorkflowMappingState{
			OldShardId:   mapping.OldShardID,
			NewShardId:   mapping.NewShardID,
			InTransition: mapping.TransitionState.InTransition(),
		}
	}

	return &ringpb.GetWorkflowShardMappingResponse{
		Mappings:       simpleMappings,
		MappingStates:  mappingStates,
		MappingVersion: version,
	}, nil
}

// ReportWorkflowTriggerRegistration handles shard registration reports
// Shards call this to inform shard zero about which workflows they have loaded
func (s *Server) ReportWorkflowTriggerRegistration(ctx context.Context, req *ringpb.ReportWorkflowTriggerRegistrationRequest) (*ringpb.ReportWorkflowTriggerRegistrationResponse, error) {
	s.logger.Debugw("ReportWorkflowTriggerRegistration called",
		"shardID", req.SourceShardId,
		"workflowCount", len(req.RegisteredWorkflows),
		"totalActive", req.TotalActiveWorkflows,
	)

	// Extract workflow IDs from the map
	workflowIDs := make([]string, 0, len(req.RegisteredWorkflows))
	for workflowID := range req.RegisteredWorkflows {
		workflowIDs = append(workflowIDs, workflowID)
	}

	err := s.store.ReportShardRegistration(ctx, req.SourceShardId, workflowIDs)
	if err != nil {
		s.logger.Errorw("Failed to update shard registrations",
			"shardID", req.SourceShardId,
			"error", err,
		)
		return &ringpb.ReportWorkflowTriggerRegistrationResponse{
			Success: false,
		}, nil
	}

	s.logger.Infow("Successfully registered workflows",
		"shardID", req.SourceShardId,
		"workflowCount", len(workflowIDs),
	)

	return &ringpb.ReportWorkflowTriggerRegistrationResponse{
		Success: true,
	}, nil
}
