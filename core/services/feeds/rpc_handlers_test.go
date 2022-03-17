package feeds_test

import (
	"context"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/services/feeds/mocks"
	pb "github.com/smartcontractkit/chainlink/core/services/feeds/proto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TestRPCHandlers struct {
	*feeds.RPCHandlers

	svc            *mocks.Service
	feedsManagerID int64
}

func setupTestHandlers(t *testing.T) *TestRPCHandlers {
	var (
		svc            = &mocks.Service{}
		feedsManagerID = int64(1)
	)

	t.Cleanup(func() {
		mock.AssertExpectationsForObjects(t,
			svc,
		)
	})

	return &TestRPCHandlers{
		RPCHandlers:    feeds.NewRPCHandlers(svc, feedsManagerID),
		svc:            svc,
		feedsManagerID: feedsManagerID,
	}
}

func Test_RPCHandlers_ProposeJob(t *testing.T) {
	var (
		ctx     = context.Background()
		jobID   = uuid.NewV4()
		spec    = TestSpec
		version = int64(1)
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
