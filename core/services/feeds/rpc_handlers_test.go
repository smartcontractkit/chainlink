package feeds_test

import (
	"context"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/services/feeds"
	pb "github.com/smartcontractkit/chainlink/core/services/feeds/proto"
	"github.com/stretchr/testify/require"
)

func Test_RPCHandlers_ProposeJob(t *testing.T) {
	svc := setupTestService(t)

	var (
		jobID          = uuid.NewV4()
		spec           = `some spec`
		feedsManagerID = int64(1)
	)
	h := feeds.NewRPCHandlers(svc, feedsManagerID)

	svc.orm.
		On("CreateJobProposal", context.Background(), &feeds.JobProposal{
			Spec:           spec,
			Status:         feeds.JobProposalStatusPending,
			FeedsManagerID: feedsManagerID,
			RemoteUUID:     jobID,
		}).
		Return(int64(1), nil)

	_, err := h.ProposeJob(context.Background(), &pb.ProposeJobRequest{
		Id:   jobID.String(),
		Spec: spec,
	})
	require.NoError(t, err)
}
