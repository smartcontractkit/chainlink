package feeds_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	pb "github.com/smartcontractkit/chainlink-protos/orchestrator/feedsmanager"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds/mocks"
)

type TestRPCHandlers struct {
	*feeds.RPCHandlers

	svc            *mocks.Service
	feedsManagerID int64
}

func setupTestHandlers(t *testing.T) *TestRPCHandlers {
	var (
		svc            = mocks.NewService(t)
		feedsManagerID = int64(1)
		lggr           = logger.TestLogger(t)
	)

	return &TestRPCHandlers{
		RPCHandlers:    feeds.NewRPCHandlers(svc, feedsManagerID, lggr),
		svc:            svc,
		feedsManagerID: feedsManagerID,
	}
}

func Test_RPCHandlers_ProposeJob(t *testing.T) {
	var (
		ctx                  = testutils.Context(t)
		jobID                = uuid.New()
		nameAndExternalJobID = uuid.New()
		spec                 = fmt.Sprintf(FluxMonitorTestSpecTemplate, nameAndExternalJobID, nameAndExternalJobID)
		version              = int64(1)
	)
	h := setupTestHandlers(t)

	h.svc.
		On("ProposeJob", ctx, &feeds.ProposeJobArgs{
			FeedsManagerID: h.feedsManagerID,
			RemoteUUID:     jobID,
			Spec:           spec,
			Version:        int32(version),
		}).
		Return(int64(1), nil)

	_, err := h.ProposeJob(ctx, &pb.ProposeJobRequest{
		Id:      jobID.String(),
		Spec:    spec,
		Version: version,
	})
	require.NoError(t, err)
}

func Test_RPCHandlers_DeleteJob(t *testing.T) {
	var (
		ctx   = testutils.Context(t)
		jobID = uuid.New()
	)
	h := setupTestHandlers(t)

	h.svc.
		On("DeleteJob", ctx, &feeds.DeleteJobArgs{
			FeedsManagerID: h.feedsManagerID,
			RemoteUUID:     jobID,
		}).
		Return(int64(1), nil)

	_, err := h.DeleteJob(ctx, &pb.DeleteJobRequest{
		Id: jobID.String(),
	})
	require.NoError(t, err)
}

func Test_RPCHandlers_RevokeJob(t *testing.T) {
	var (
		ctx   = testutils.Context(t)
		jobID = uuid.New()
	)
	h := setupTestHandlers(t)

	h.svc.
		On("RevokeJob", ctx, &feeds.RevokeJobArgs{
			FeedsManagerID: h.feedsManagerID,
			RemoteUUID:     jobID,
		}).
		Return(int64(1), nil)

	_, err := h.RevokeJob(ctx, &pb.RevokeJobRequest{
		Id: jobID.String(),
	})
	require.NoError(t, err)
}

func Test_RPCHandlers_GetJobRuns(t *testing.T) {
	var (
		ctx   = testutils.Context(t)
		jobID = uuid.New()
		limit = uint32(10)
	)
	h := setupTestHandlers(t)

	mockSummaries := []*pb.JobRunSummary{
		{
			RunId:       1,
			State:       "completed",
			AllErrors:   []string{},
			FatalErrors: []string{},
		},
		{
			RunId:       2,
			State:       "errored",
			AllErrors:   []string{"error1"},
			FatalErrors: []string{"fatal1"},
		},
	}

	h.svc.EXPECT().
		GetJobRuns(ctx, &feeds.GetJobRunsArgs{
			FeedsManagerID: h.feedsManagerID,
			RemoteUUID:     jobID,
			Limit:          limit,
		}).
		Return(mockSummaries, nil)

	response, err := h.GetJobRuns(ctx, &pb.GetJobRunsRequest{
		Id:    jobID.String(),
		Limit: limit,
	})

	require.NoError(t, err)
	require.NotNil(t, response)
	require.Len(t, response.Runs, 2)

	run1 := response.Runs[0]
	require.Equal(t, int64(1), run1.RunId)
	require.Equal(t, "completed", run1.State)
	require.Empty(t, run1.AllErrors)
	require.Empty(t, run1.FatalErrors)

	run2 := response.Runs[1]
	require.Equal(t, int64(2), run2.RunId)
	require.Equal(t, "errored", run2.State)
	require.Len(t, run2.AllErrors, 1)
	require.Equal(t, "error1", run2.AllErrors[0])
	require.Len(t, run2.FatalErrors, 1)
	require.Equal(t, "fatal1", run2.FatalErrors[0])
}

func Test_RPCHandlers_GetJobRuns_InvalidUUID(t *testing.T) {
	ctx := testutils.Context(t)
	h := setupTestHandlers(t)

	_, err := h.GetJobRuns(ctx, &pb.GetJobRunsRequest{
		Id:    "invalid-uuid",
		Limit: 10,
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid UUID")
}

func Test_RPCHandlers_GetJobRuns_DefaultLimit(t *testing.T) {
	var (
		ctx   = testutils.Context(t)
		jobID = uuid.New()
	)
	h := setupTestHandlers(t)

	h.svc.EXPECT().GetJobRuns(ctx, &feeds.GetJobRunsArgs{
		FeedsManagerID: h.feedsManagerID,
		RemoteUUID:     jobID,
		Limit:          50,
	}).
		Return([]*pb.JobRunSummary{}, nil)

	_, err := h.GetJobRuns(ctx, &pb.GetJobRunsRequest{
		Id:    jobID.String(),
		Limit: 0,
	})

	require.NoError(t, err)
}

func Test_RPCHandlers_GetJobRuns_ExceedsMaxLimit(t *testing.T) {
	var (
		ctx   = testutils.Context(t)
		jobID = uuid.New()
	)
	h := setupTestHandlers(t)

	h.svc.EXPECT().GetJobRuns(ctx, &feeds.GetJobRunsArgs{
		FeedsManagerID: h.feedsManagerID,
		RemoteUUID:     jobID,
		Limit:          50,
	}).
		Return([]*pb.JobRunSummary{}, nil)

	_, err := h.GetJobRuns(ctx, &pb.GetJobRunsRequest{
		Id:    jobID.String(),
		Limit: 150,
	})

	require.NoError(t, err)
}
