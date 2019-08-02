package web_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkJobRunsController_Index(b *testing.B) {
	app, cleanup := cltest.NewApplication(b)
	defer cleanup()
	app.Start()
	run1, _, _ := setupJobRunsControllerIndex(b, app)
	client := app.NewHTTPClient()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		resp, cleanup := client.Get("/v2/runs?jobSpecId=" + run1.JobSpecID)
		defer cleanup()
		assert.Equal(b, 200, resp.StatusCode, "Response should be successful")
	}
}

func TestJobRunsController_Index(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	runA, runB, runC := setupJobRunsControllerIndex(t, app)

	resp, cleanup := client.Get("/v2/runs?size=x&jobSpecId=" + runA.JobSpecID)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 422)

	resp, cleanup = client.Get("/v2/runs?size=1&jobSpecId=" + runA.JobSpecID)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	var links jsonapi.Links
	var runs []models.JobRun

	err := web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &runs, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, runs, 1)
	assert.Equal(t, runA.ID, runs[0].ID, "expected runs order by createdAt ascending")

	resp, cleanup = client.Get(links["next"].Href)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	var nextPageLinks jsonapi.Links
	var nextPageRuns = []models.JobRun{}

	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &nextPageRuns, &nextPageLinks)
	assert.NoError(t, err)
	assert.Empty(t, nextPageLinks["next"])
	assert.NotEmpty(t, nextPageLinks["prev"])

	assert.Len(t, nextPageRuns, 1)
	assert.Equal(t, runB.ID, nextPageRuns[0].ID, "expected runs order by createdAt ascending")

	resp, cleanup = client.Get("/v2/runs?sort=-createdAt")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	var allJobRunLinks jsonapi.Links
	var allJobRuns []models.JobRun

	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &allJobRuns, &allJobRunLinks)
	assert.NoError(t, err)
	assert.Empty(t, allJobRunLinks["next"].Href)
	assert.Empty(t, allJobRunLinks["prev"].Href)

	assert.Len(t, allJobRuns, 3)
	assert.Equal(t, runC.ID, allJobRuns[0].ID, "expected runs ordered by created at descending")
	assert.Equal(t, runB.ID, allJobRuns[1].ID, "expected runs ordered by created at descending")
	assert.Equal(t, runA.ID, allJobRuns[2].ID, "expected runs ordered by created at descending")
}

func setupJobRunsControllerIndex(t assert.TestingT, app *cltest.TestApplication) (*models.JobRun, *models.JobRun, *models.JobRun) {
	j1 := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.CreateJob(&j1))

	j2 := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.CreateJob(&j2))

	now := time.Now()

	runA := j1.NewRun(j1.Initiators[0])
	runA.ID = "runA"
	runA.CreatedAt = now.Add(-2 * time.Second)
	assert.Nil(t, app.Store.CreateJobRun(&runA))

	runB := j1.NewRun(j1.Initiators[0])
	runB.ID = "runB"
	runB.CreatedAt = now.Add(-time.Second)
	assert.Nil(t, app.Store.CreateJobRun(&runB))

	runC := j2.NewRun(j2.Initiators[0])
	runC.ID = "runC"
	runC.CreatedAt = now
	assert.Nil(t, app.Store.CreateJobRun(&runC))

	return &runA, &runB, &runC
}

func TestJobRunsController_Create_Success(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	app.Start()
	defer cleanup()

	j := cltest.NewJobWithWebInitiator()
	assert.NoError(t, app.Store.CreateJob(&j))

	jr := cltest.CreateJobRunViaWeb(t, app, j, `{"result":"100"}`)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)
	val, err := jr.Result.ResultString()
	assert.NoError(t, err)
	assert.Equal(t, "100", val)
}

func TestJobRunsController_Create_Archived(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	app.Start()
	defer cleanup()

	j := cltest.NewJobWithWebInitiator()
	require.NoError(t, app.Store.CreateJob(&j))
	require.NoError(t, app.Store.ArchiveJob(j.ID))

	client := app.NewHTTPClient()
	resp, cleanup := client.Post("/v2/specs/"+j.ID+"/runs", bytes.NewBufferString(`{"result":"100"}`))
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 404)
}

func TestJobRunsController_Create_EmptyBody(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	app.Start()
	defer cleanup()

	j := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.CreateJob(&j))

	jr := cltest.CreateJobRunViaWeb(t, app, j)
	cltest.WaitForJobRunToComplete(t, app.Store, jr)
}

func TestJobRunsController_Create_InvalidBody(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	j := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.CreateJob(&j))

	resp, cleanup := client.Post("/v2/specs/"+j.ID+"/runs", bytes.NewBufferString(`{`))
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 500)
}

func TestJobRunsController_Create_WithoutWebInitiator(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	j := cltest.NewJob()
	assert.Nil(t, app.Store.CreateJob(&j))

	resp, cleanup := client.Post("/v2/specs/"+j.ID+"/runs", bytes.NewBuffer([]byte{}))
	defer cleanup()
	assert.Equal(t, 403, resp.StatusCode, "Response should be forbidden")
}

func TestJobRunsController_Create_NotFound(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	resp, cleanup := client.Post("/v2/specs/garbageID/runs", bytes.NewBuffer([]byte{}))
	defer cleanup()
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
}

func TestJobRunsController_Update_Success(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	app.Start()
	defer cleanup()

	tests := []struct {
		name     string
		archived bool
	}{
		{"normal_job", false},
		{"archived_job", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bta, bt := cltest.NewBridgeType(t, test.name)
			require.NoError(t, app.Store.CreateBridgeType(bt))
			j := cltest.NewJobWithWebInitiator()
			j.Tasks = []models.TaskSpec{{Type: bt.Name}}
			require.NoError(t, app.Store.CreateJob(&j))
			jr := cltest.MarkJobRunPendingBridge(j.NewRun(j.Initiators[0]), 0)
			require.NoError(t, app.Store.CreateJobRun(&jr))

			if test.archived {
				require.NoError(t, app.Store.ArchiveJob(j.ID))
			}

			// resume run
			body := fmt.Sprintf(`{"id":"%v","data":{"result": "100"}}`, jr.ID)
			headers := map[string]string{"Authorization": "Bearer " + bta.IncomingToken}
			url := app.Configger.ClientNodeURL() + "/v2/runs/" + jr.ID
			resp, cleanup := cltest.UnauthenticatedPatch(t, url, bytes.NewBufferString(body), headers)
			defer cleanup()

			require.Equal(t, 200, resp.StatusCode, "Response should be successful")
			var respJobRun presenters.JobRun
			assert.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &respJobRun))
			require.Equal(t, jr.ID, respJobRun.ID)

			jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)
			val, err := jr.Result.ResultString()
			assert.NoError(t, err)
			assert.Equal(t, "100", val)
		})
	}
}

func TestJobRunsController_Update_WrongAccessToken(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	_, bt := cltest.NewBridgeType(t)
	assert.Nil(t, app.Store.CreateBridgeType(bt))
	j := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.Nil(t, app.Store.CreateJob(&j))
	jr := cltest.MarkJobRunPendingBridge(j.NewRun(j.Initiators[0]), 0)
	assert.Nil(t, app.Store.CreateJobRun(&jr))

	body := fmt.Sprintf(`{"id":"%v","data":{"result": "100"}}`, jr.ID)
	headers := map[string]string{"Authorization": "Bearer " + "wrongaccesstoken"}
	resp, cleanup := client.Patch("/v2/runs/"+jr.ID, bytes.NewBufferString(body), headers)
	defer cleanup()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Response should be unauthorized")
	jr, err := app.Store.FindJobRun(jr.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.RunStatusPendingBridge, jr.Status)
}

func TestJobRunsController_Update_NotPending(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	bta, bt := cltest.NewBridgeType(t)
	assert.Nil(t, app.Store.CreateBridgeType(bt))
	j := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.Nil(t, app.Store.CreateJob(&j))
	jr := j.NewRun(j.Initiators[0])
	assert.Nil(t, app.Store.CreateJobRun(&jr))

	body := fmt.Sprintf(`{"id":"%v","data":{"result": "100"}}`, jr.ID)
	headers := map[string]string{"Authorization": "Bearer " + bta.IncomingToken}
	resp, cleanup := client.Patch("/v2/runs/"+jr.ID, bytes.NewBufferString(body), headers)
	defer cleanup()
	assert.Equal(t, 405, resp.StatusCode, "Response should be unsuccessful")
}

func TestJobRunsController_Update_WithError(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	bta, bt := cltest.NewBridgeType(t)
	assert.Nil(t, app.Store.CreateBridgeType(bt))
	j := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.Nil(t, app.Store.CreateJob(&j))
	jr := cltest.MarkJobRunPendingBridge(j.NewRun(j.Initiators[0]), 0)
	assert.Nil(t, app.Store.CreateJobRun(&jr))

	body := fmt.Sprintf(`{"id":"%v","error":"stack overflow","data":{"result": "0"}}`, jr.ID)
	headers := map[string]string{"Authorization": "Bearer " + bta.IncomingToken}
	resp, cleanup := client.Patch("/v2/runs/"+jr.ID, bytes.NewBufferString(body), headers)
	defer cleanup()
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	var respJobRun presenters.JobRun
	assert.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &respJobRun))
	assert.Equal(t, jr.ID, respJobRun.ID)

	jr = cltest.WaitForJobRunStatus(t, app.Store, jr, models.RunStatusErrored)
	val, err := jr.Result.ResultString()
	assert.NoError(t, err)
	assert.Equal(t, "0", val)
}

func TestJobRunsController_Update_BadInput(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	_, bt := cltest.NewBridgeType(t)
	assert.Nil(t, app.Store.CreateBridgeType(bt))
	j := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.Nil(t, app.Store.CreateJob(&j))
	jr := cltest.MarkJobRunPendingBridge(j.NewRun(j.Initiators[0]), 0)
	assert.Nil(t, app.Store.CreateJobRun(&jr))

	body := fmt.Sprint(`{`, jr.ID)
	resp, cleanup := client.Patch("/v2/runs/"+jr.ID, bytes.NewBufferString(body))
	defer cleanup()
	assert.Equal(t, 500, resp.StatusCode, "Response should be successful")
	jr, err := app.Store.FindJobRun(jr.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.RunStatusPendingBridge, jr.Status)
}

func TestJobRunsController_Update_NotFound(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	_, bt := cltest.NewBridgeType(t)
	assert.Nil(t, app.Store.CreateBridgeType(bt))
	j := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.Nil(t, app.Store.CreateJob(&j))
	jr := cltest.MarkJobRunPendingBridge(j.NewRun(j.Initiators[0]), 0)
	assert.Nil(t, app.Store.CreateJobRun(&jr))

	body := fmt.Sprintf(`{"id":"%v","data":{"result": "100"}}`, jr.ID)
	resp, cleanup := client.Patch("/v2/runs/"+jr.ID+"1", bytes.NewBufferString(body))
	defer cleanup()
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
	jr, err := app.Store.FindJobRun(jr.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.RunStatusPendingBridge, jr.Status)
}

func TestJobRunsController_Show_Found(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	j := cltest.NewJobWithSchedule("9 9 9 9 6")
	app.Store.CreateJob(&j)

	jr := j.NewRun(j.Initiators[0])
	assert.NoError(t, app.Store.CreateJobRun(&jr))

	resp, cleanup := client.Get("/v2/runs/" + jr.ID)
	defer cleanup()
	require.Equal(t, 200, resp.StatusCode, "Response should be successful")

	var respJobRun presenters.JobRun
	assert.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &respJobRun))
	assert.Equal(t, jr.Initiator.Schedule, respJobRun.Initiator.Schedule, "should have the same schedule")
	assert.Equal(t, jr.ID, respJobRun.ID, "should have job run id")
}

func TestJobRunsController_Show_NotFound(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v2/runs/garbage")
	defer cleanup()
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
}

func TestJobRunsController_Show_Unauthenticated(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	app.Start()
	defer cleanup()

	resp, err := http.Get(app.Server.URL + "/v2/runs/notauthorized")
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode, "Response should be forbidden")
}
