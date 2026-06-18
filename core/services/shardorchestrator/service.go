package shardorchestrator

import (
	"context"
	"errors"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/services/ring"
)

// Server implements the gRPC ShardOrchestratorService
// This runs on shard zero and serves requests from other shards
type Server struct {
	ringpb.UnimplementedShardOrchestratorServiceServer
	ringStore *ring.Store
	logger    logger.Logger
}

func NewServer(ringStore *ring.Store, lggr logger.Logger) *Server {
	return &Server{
		ringStore: ringStore,
		logger:    logger.Named(lggr, "ShardOrchestratorServer"),
	}
}

// RegisterWithGRPCServer registers this service with a gRPC server
func (s *Server) RegisterWithGRPCServer(grpcServer *grpc.Server) {
	ringpb.RegisterShardOrchestratorServiceServer(grpcServer, s)
	s.logger.Info("Registered ShardOrchestrator gRPC service")
}

// GetWorkflowShardMapping handles batch requests for workflow-to-shard mappings
// This is called by other shards to determine where to route workflow executions
func (s *Server) GetWorkflowShardMapping(_ context.Context, req *ringpb.GetWorkflowShardMappingRequest) (*ringpb.GetWorkflowShardMappingResponse, error) {
	s.logger.Debugw("GetWorkflowShardMapping called", "workflowCount", len(req.WorkflowIds))

	if len(req.WorkflowIds) == 0 {
		return nil, errors.New("workflow_ids is required and must not be empty")
	}

	mappings, version := s.ringStore.GetWorkflowMappingsBatch(req.WorkflowIds)

	var missing []string
	for _, wfID := range req.WorkflowIds {
		if _, exists := mappings[wfID]; !exists {
			missing = append(missing, wfID)
		}
	}

	if len(missing) > 0 {
		dropped := s.ringStore.SubmitWorkflowsForAllocation(missing)
		s.logger.Debugw("Submitted missing workflows for allocation", "count", len(missing))
		if dropped > 0 {
			s.logger.Warnw("Allocation request channel full, workflows dropped after retries", "dropped", dropped)
		}
	}

	simpleMappings := make(map[string]uint32, len(mappings))
	mappingStates := make(map[string]*ringpb.WorkflowMappingState, len(mappings))

	for workflowID, meta := range mappings {
		simpleMappings[workflowID] = meta.NewShardID
		mappingStates[workflowID] = &ringpb.WorkflowMappingState{
			OldShardId:   meta.OldShardID,
			NewShardId:   meta.NewShardID,
			InTransition: meta.InTransition,
		}
	}

	rs := s.ringStore.GetRoutingState()
	resp := &ringpb.GetWorkflowShardMappingResponse{
		Mappings:       simpleMappings,
		MappingStates:  mappingStates,
		MappingVersion: version,
		RoutingSteady:  ring.IsInSteadyState(rs),
	}
	if rs != nil {
		resp.RoutingStateId = rs.Id
	}
	return resp, nil
}

// ReportWorkflowTriggerRegistration handles shard registration reports
// Shards call this to inform shard zero about which workflows they have loaded
func (s *Server) ReportWorkflowTriggerRegistration(_ context.Context, req *ringpb.ReportWorkflowTriggerRegistrationRequest) (*ringpb.ReportWorkflowTriggerRegistrationResponse, error) {
	s.logger.Debugw("ReportWorkflowTriggerRegistration called",
		"shardID", req.SourceShardId,
		"workflowCount", len(req.RegisteredWorkflows),
		"totalActive", req.TotalActiveWorkflows,
	)

	workflowIDs := make([]string, 0, len(req.RegisteredWorkflows))
	for workflowID := range req.RegisteredWorkflows {
		workflowIDs = append(workflowIDs, workflowID)
	}

	s.ringStore.RegisterWorkflowsFromShard(req.SourceShardId, workflowIDs)

	s.logger.Infow("Successfully registered workflows",
		"shardID", req.SourceShardId,
		"workflowCount", len(workflowIDs),
	)

	return &ringpb.ReportWorkflowTriggerRegistrationResponse{
		Success: true,
	}, nil
}
