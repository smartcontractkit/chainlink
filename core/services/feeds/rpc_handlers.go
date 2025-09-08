package feeds

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	pb "github.com/smartcontractkit/chainlink-protos/orchestrator/feedsmanager"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	MaxJobRunsLimit     = 100
	DefaultJobRunsLimit = 50
)

// RPCHandlers define handlers for RPC method calls from the Feeds Manager
type RPCHandlers struct {
	svc            Service
	feedsManagerID int64
	lggr           logger.Logger
}

func NewRPCHandlers(svc Service, feedsManagerID int64, lggr logger.Logger) *RPCHandlers {
	return &RPCHandlers{
		svc:            svc,
		feedsManagerID: feedsManagerID,
		lggr:           lggr.Named("RPCHandlers"),
	}
}

// ProposeJob creates a new job proposal record for the feeds manager
func (h *RPCHandlers) ProposeJob(ctx context.Context, req *pb.ProposeJobRequest) (*pb.ProposeJobResponse, error) {
	remoteUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	_, err = h.svc.ProposeJob(ctx, &ProposeJobArgs{
		Spec:           req.GetSpec(),
		FeedsManagerID: h.feedsManagerID,
		RemoteUUID:     remoteUUID,
		Version:        int32(req.GetVersion()),
		Multiaddrs:     req.GetMultiaddrs(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.ProposeJobResponse{}, nil
}

// DeleteJob deletes a job proposal record.
func (h *RPCHandlers) DeleteJob(ctx context.Context, req *pb.DeleteJobRequest) (*pb.DeleteJobResponse, error) {
	remoteUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	_, err = h.svc.DeleteJob(ctx, &DeleteJobArgs{
		FeedsManagerID: h.feedsManagerID,
		RemoteUUID:     remoteUUID,
	})
	if err != nil {
		return nil, err
	}

	return &pb.DeleteJobResponse{}, nil
}

// RevokeJob revokes a pending job proposal record.
func (h *RPCHandlers) RevokeJob(ctx context.Context, req *pb.RevokeJobRequest) (*pb.RevokeJobResponse, error) {
	remoteUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	_, err = h.svc.RevokeJob(ctx, &RevokeJobArgs{
		FeedsManagerID: h.feedsManagerID,
		RemoteUUID:     remoteUUID,
	})
	if err != nil {
		return nil, err
	}

	return &pb.RevokeJobResponse{}, nil
}

// GetJobRuns fetches job run history for the specified job proposal
func (h *RPCHandlers) GetJobRuns(ctx context.Context, req *pb.GetJobRunsRequest) (*pb.GetJobRunsResponse, error) {
	remoteUUID, err := uuid.Parse(req.Id)
	limit := req.Limit
	if err != nil {
		return nil, fmt.Errorf("unable to parse request id (%s): %w", req.Id, err)
	}

	if limit == 0 || limit > MaxJobRunsLimit {
		h.lggr.Warnw("Invalid limit provided, using default",
			"requestedLimit", limit,
			"defaultLimit", DefaultJobRunsLimit,
		)
		limit = DefaultJobRunsLimit
	}

	summaries, err := h.svc.GetJobRuns(ctx, &GetJobRunsArgs{
		FeedsManagerID: h.feedsManagerID,
		RemoteUUID:     remoteUUID,
		Limit:          limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get job runs: %w", err)
	}

	return &pb.GetJobRunsResponse{
		Runs: summaries,
	}, nil
}
