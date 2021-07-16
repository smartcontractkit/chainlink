package web_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/feeds"
	pbMocks "github.com/smartcontractkit/chainlink/core/services/feeds/mocks"
	pb "github.com/smartcontractkit/chainlink/core/services/feeds/proto"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_JobProposalsController_Reject(t *testing.T) {
	t.Parallel()

	var (
		jp1 = feeds.JobProposal{
			ID:             1,
			RemoteUUID:     uuid.NewV4(),
			Spec:           "some spec",
			Status:         feeds.JobProposalStatusPending,
			JobID:          uuid.NullUUID{},
			FeedsManagerID: 10,
		}
		expected = jp1
	)
	expected.Status = feeds.JobProposalStatusRejected

	testCases := []struct {
		name           string
		before         func(t *testing.T, app *cltest.TestApplication, id *string, rpcClient *pbMocks.FeedsManagerClient)
		want           *feeds.JobProposal
		wantStatusCode int
	}{
		{
			name: "success",
			before: func(t *testing.T, app *cltest.TestApplication, id *string, rpcClient *pbMocks.FeedsManagerClient) {
				fsvc := app.GetFeedsService()

				jp1ID, err := fsvc.CreateJobProposal(&jp1)
				require.NoError(t, err)

				*id = strconv.Itoa(int(jp1ID))

				rpcClient.On("RejectedJob", mock.MatchedBy(func(c context.Context) bool { return true }), &pb.RejectedJobRequest{
					Uuid: jp1.RemoteUUID.String(),
				}).Return(&pb.RejectedJobResponse{}, nil)
			},
			wantStatusCode: http.StatusOK,
			want:           &expected,
		},
		{
			name: "invalid id",
			before: func(t *testing.T, app *cltest.TestApplication, id *string, rpcClient *pbMocks.FeedsManagerClient) {
				*id = "notanint"
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "not found",
			before: func(t *testing.T, app *cltest.TestApplication, id *string, rpcClient *pbMocks.FeedsManagerClient) {
				*id = "999999999"
			},
			wantStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app, client := setupJobProposalsTest(t)
			rpcClient := &pbMocks.FeedsManagerClient{}
			app.FeedsService.Unsafe_SetFMSClient(rpcClient)

			// Defer the FK requirement of a feeds manager.
			require.NoError(t, app.Store.DB.Exec(
				`SET CONSTRAINTS fk_feeds_manager DEFERRED`,
			).Error)

			var id string
			if tc.before != nil {
				tc.before(t, app, &id, rpcClient)
			}

			resp, cleanup := client.Post(fmt.Sprintf("/v2/job_proposals/%s/reject", id), bytes.NewReader([]byte{}))
			t.Cleanup(cleanup)
			require.Equal(t, tc.wantStatusCode, resp.StatusCode)

			if tc.want != nil {
				resource := presenters.JobProposalResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource)
				require.NoError(t, err)

				assert.Equal(t, id, resource.ID)
				assert.Equal(t, tc.want.Spec, resource.Spec)
				assert.Equal(t, tc.want.Status, resource.Status)
			}
		})
	}
}

func setupJobProposalsTest(t *testing.T) (*cltest.TestApplication, cltest.HTTPClientCleaner) {
	app, cleanup := cltest.NewApplication(t)
	t.Cleanup(cleanup)
	app.Start()

	client := app.NewHTTPClient()

	return app, client
}
