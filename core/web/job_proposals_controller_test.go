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

func Test_JobProposalsController_Index(t *testing.T) {
	t.Parallel()

	var (
		spec = string(cltest.MustReadFile(t, "../testdata/tomlspecs/flux-monitor-spec.toml"))
		jp1  = feeds.JobProposal{
			ID:             1,
			RemoteUUID:     uuid.NewV4(),
			Spec:           spec,
			Status:         feeds.JobProposalStatusPending,
			ExternalJobID:  uuid.NullUUID{},
			FeedsManagerID: 10,
		}
	)

	testCases := []struct {
		name           string
		before         func(t *testing.T, app *cltest.TestApplication, id *string)
		want           *feeds.JobProposal
		wantStatusCode int
	}{
		{
			name:           "success",
			wantStatusCode: http.StatusOK,
			want:           &jp1,
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

			// Create the job proposal
			fsvc := app.GetFeedsService()
			id, err := fsvc.CreateJobProposal(&jp1)
			require.NoError(t, err)

			resp, cleanup := client.Get("/v2/job_proposals")
			t.Cleanup(cleanup)
			require.Equal(t, tc.wantStatusCode, resp.StatusCode)

			if tc.want != nil {
				resources := []presenters.JobProposalResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resources)
				require.NoError(t, err)
				require.Len(t, resources, 1)

				actual := resources[0]
				assert.Equal(t, strconv.FormatInt(id, 10), actual.ID)
				assert.Equal(t, tc.want.Spec, actual.Spec)
				assert.Equal(t, tc.want.Status, actual.Status)
			}
		})
	}
}

func Test_JobProposalsController_Show(t *testing.T) {
	t.Parallel()

	var (
		spec = string(cltest.MustReadFile(t, "../testdata/tomlspecs/flux-monitor-spec.toml"))
		jp   = feeds.JobProposal{
			ID:             1,
			RemoteUUID:     uuid.NewV4(),
			Spec:           spec,
			Status:         feeds.JobProposalStatusPending,
			ExternalJobID:  uuid.NullUUID{},
			FeedsManagerID: 10,
		}
	)

	testCases := []struct {
		name           string
		before         func(t *testing.T, app *cltest.TestApplication, id *string)
		want           *feeds.JobProposal
		wantStatusCode int
	}{
		{
			name:           "success",
			wantStatusCode: http.StatusOK,
			want:           &jp,
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

			// Create the job proposal
			fsvc := app.GetFeedsService()
			id, err := fsvc.CreateJobProposal(&jp)
			require.NoError(t, err)

			resp, cleanup := client.Get(fmt.Sprintf("/v2/job_proposals/%d", id))
			t.Cleanup(cleanup)
			require.Equal(t, tc.wantStatusCode, resp.StatusCode)

			if tc.want != nil {
				actual := presenters.JobProposalResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &actual)
				require.NoError(t, err)

				assert.Equal(t, strconv.FormatInt(id, 10), actual.ID)
				assert.Equal(t, tc.want.Spec, actual.Spec)
				assert.Equal(t, tc.want.Status, actual.Status)
			}
		})
	}
}

func Test_JobProposalsController_Approve(t *testing.T) {
	t.Parallel()

	var (
		spec = string(cltest.MustReadFile(t, "../testdata/tomlspecs/flux-monitor-spec.toml"))
		jp1  = feeds.JobProposal{
			ID:             1,
			RemoteUUID:     uuid.NewV4(),
			Spec:           spec,
			Status:         feeds.JobProposalStatusPending,
			ExternalJobID:  uuid.NullUUID{},
			FeedsManagerID: 10,
		}
		expected = jp1
	)
	expected.Status = feeds.JobProposalStatusApproved

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

				rpcClient.On("ApprovedJob", mock.MatchedBy(func(c context.Context) bool { return true }), &pb.ApprovedJobRequest{
					Uuid: jp1.RemoteUUID.String(),
				}).Return(&pb.ApprovedJobResponse{}, nil)
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

			resp, cleanup := client.Post(fmt.Sprintf("/v2/job_proposals/%s/approve", id), bytes.NewReader([]byte{}))
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

func Test_JobProposalsController_Reject(t *testing.T) {
	t.Parallel()

	var (
		jp1 = feeds.JobProposal{
			ID:             1,
			RemoteUUID:     uuid.NewV4(),
			Spec:           "some spec",
			Status:         feeds.JobProposalStatusPending,
			ExternalJobID:  uuid.NullUUID{},
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

func Test_JobProposalsController_UpdateSpec(t *testing.T) {
	t.Parallel()

	var (
		jp1 = feeds.JobProposal{
			ID:             1,
			RemoteUUID:     uuid.NewV4(),
			Spec:           "some spec",
			Status:         feeds.JobProposalStatusPending,
			ExternalJobID:  uuid.NullUUID{},
			FeedsManagerID: 10,
		}
		reqBody  = `{"spec": "updated spec"}`
		expected = jp1
	)
	expected.Spec = "updated spec"

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

			// Defer the FK requirement of a feeds manager.
			require.NoError(t, app.Store.DB.Exec(
				`SET CONSTRAINTS fk_feeds_manager DEFERRED`,
			).Error)

			var id string
			if tc.before != nil {
				tc.before(t, app, &id, rpcClient)
			}

			resp, cleanup := client.Patch(
				fmt.Sprintf("/v2/job_proposals/%s/spec", id),
				bytes.NewReader([]byte(reqBody)),
			)
			t.Cleanup(cleanup)
			require.Equal(t, tc.wantStatusCode, resp.StatusCode)

			if tc.want != nil {
				resource := presenters.JobProposalResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource)
				require.NoError(t, err)

				assert.Equal(t, id, resource.ID)
				assert.Equal(t, tc.want.Spec, resource.Spec)
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
