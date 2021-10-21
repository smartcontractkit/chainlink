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
		jobID = uuid.NewV4()
		spec  = TestSpec
	)
	h := setupTestHandlers(t)

	h.svc.
		On("ProposeJob", &feeds.JobProposal{
			Spec:           spec,
			FeedsManagerID: h.feedsManagerID,
			RemoteUUID:     jobID,
		}).
		Return(int64(1), nil)

	_, err := h.ProposeJob(context.Background(), &pb.ProposeJobRequest{
		Id:   jobID.String(),
		Spec: spec,
	})
	require.NoError(t, err)
}
