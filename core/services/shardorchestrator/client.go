package shardorchestrator

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator/pb"
)

// Client wraps gRPC client for communicating with shard 0's orchestrator service
type Client struct {
	conn   *grpc.ClientConn
	client pb.ShardOrchestratorServiceClient
	logger logger.Logger
}

// NewClient creates a new gRPC client to communicate with the shard orchestrator on shard 0
func NewClient(ctx context.Context, address string, lggr logger.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create shard orchestrator client for %s: %w", address, err)
	}

	return &Client{
		conn:   conn,
		client: pb.NewShardOrchestratorServiceClient(conn),
		logger: logger.Named(lggr, "ShardOrchestratorClient"),
	}, nil
}

// GetWorkflowShardMapping queries shard 0 for workflow-to-shard mappings
func (c *Client) GetWorkflowShardMapping(ctx context.Context, workflowIDs []string) (*pb.GetWorkflowShardMappingResponse, error) {
	c.logger.Debugw("Calling GetWorkflowShardMapping", "workflowCount", len(workflowIDs))

	req := &pb.GetWorkflowShardMappingRequest{
		WorkflowIds: workflowIDs,
	}

	resp, err := c.client.GetWorkflowShardMapping(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("gRPC GetWorkflowShardMapping failed: %w", err)
	}

	c.logger.Debugw("GetWorkflowShardMapping response received",
		"mappingCount", len(resp.Mappings),
		"version", resp.MappingVersion)

	return resp, nil
}

// ReportWorkflowTriggerRegistration reports workflow trigger registration to shard 0
func (c *Client) ReportWorkflowTriggerRegistration(ctx context.Context, req *pb.ReportWorkflowTriggerRegistrationRequest) (*pb.ReportWorkflowTriggerRegistrationResponse, error) {
	c.logger.Debugw("Calling ReportWorkflowTriggerRegistration",
		"shardID", req.SourceShardId,
		"workflowCount", len(req.RegisteredWorkflows))

	resp, err := c.client.ReportWorkflowTriggerRegistration(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("gRPC ReportWorkflowTriggerRegistration failed: %w", err)
	}

	c.logger.Debugw("ReportWorkflowTriggerRegistration response received",
		"success", resp.Success)

	return resp, nil
}

// Close closes the gRPC connection
func (c *Client) Close() error {
	c.logger.Info("Closing ShardOrchestrator gRPC client")
	return c.conn.Close()
}
