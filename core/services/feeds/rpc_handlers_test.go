package feeds_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	pb "github.com/smartcontractkit/chainlink-protos/orchestrator/feedsmanager"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
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
	)

	return &TestRPCHandlers{
		RPCHandlers:    feeds.NewRPCHandlers(svc, feedsManagerID),
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
