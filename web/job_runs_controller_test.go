package web_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/web"
	"github.com/stretchr/testify/assert"
)

type JobRunsJSON struct {
	Runs []JobRun `json:"runs"`
}

type JobRun struct {
	ID string `json:"id"`
}

func BenchmarkJobRunsController_Index(b *testing.B) {
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()
	run1, _, _ := setupJobRunsControllerIndex(b, app)
	client := app.NewHTTPClient()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		resp, cleanup := client.Get("/v2/runs?jobSpecId=" + run1.JobID)
		defer cleanup()
		assert.Equal(b, 200, resp.StatusCode, "Response should be successful")
	}
}

func TestJobRunsController_Index(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	run1, run2, run3 := setupJobRunsControllerIndex(t, app)

	resp, cleanup := client.Get("/v2/runs?size=x&jobSpecId=" + run1.JobID)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 422)

	resp, cleanup = client.Get("/v2/runs?size=1&jobSpecId=" + run1.JobID)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	var links jsonapi.Links
	var runs []models.JobRun

	err := web.ParsePaginatedResponse(cltest.ParseResponseBody(resp), &runs, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, runs, 1)
	assert.Equal(t, run1.ID, runs[0].ID, "expected runs order by createdAt ascending")

	resp, cleanup = client.Get(links["next"].Href)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	var nextPageLinks jsonapi.Links
	var nextPageRuns = []models.JobRun{}

	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(resp), &nextPageRuns, &nextPageLinks)
	assert.NoError(t, err)
	assert.Empty(t, nextPageLinks["next"])
	assert.NotEmpty(t, nextPageLinks["prev"])

	assert.Len(t, nextPageRuns, 1)
	assert.Equal(t, run2.ID, nextPageRuns[0].ID, "expected runs order by createdAt ascending")

	resp, cleanup = client.Get("/v2/runs?sort=-createdAt")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	var allJobRunLinks jsonapi.Links
	var allJobRuns []models.JobRun

	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(resp), &allJobRuns, &allJobRunLinks)
	assert.NoError(t, err)
	assert.Empty(t, allJobRunLinks["next"].Href)
	assert.Empty(t, allJobRunLinks["prev"].Href)

	assert.Len(t, allJobRuns, 3)
	assert.Equal(t, run3.ID, allJobRuns[0].ID, "expected runs ordered by created at descending")
	assert.Equal(t, run2.ID, allJobRuns[1].ID, "expected runs ordered by created at descending")
	assert.Equal(t, run1.ID, allJobRuns[2].ID, "expected runs ordered by created at descending")
}

func setupJobRunsControllerIndex(t assert.TestingT, app *cltest.TestApplication) (*models.JobRun, *models.JobRun, *models.JobRun) {
	j1, initr := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.SaveJob(&j1))

	j2, initr := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.SaveJob(&j2))

	run1 := j1.NewRun(initr)
	run1.ID = "runA"
	assert.Nil(t, app.Store.Save(&run1))

	run2 := j1.NewRun(initr)
	run2.ID = "runB"
	run2.CreatedAt = run1.CreatedAt.Add(time.Second)
	assert.Nil(t, app.Store.Save(&run2))

	run3 := j2.NewRun(initr)
	run3.ID = "runC"
	run3.CreatedAt = run1.CreatedAt.Add(2 * time.Second)
	assert.Nil(t, app.Store.Save(&run3))

	return &run1, &run2, &run3
}

func TestJobRunsController_Create_Success(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()

	j, _ := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.SaveJob(&j))

	jr := cltest.CreateJobRunViaWeb(t, app, j, `{"value":"100"}`)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)
	val, err := jr.Result.Value()
	assert.NoError(t, err)
	assert.Equal(t, "100", val)
}

func TestJobRunsController_Create_EmptyBody(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()

	j, _ := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.SaveJob(&j))

	jr := cltest.CreateJobRunViaWeb(t, app, j)
	cltest.WaitForJobRunToComplete(t, app.Store, jr)
}

func TestJobRunsController_Create_InvalidBody(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	j, _ := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.SaveJob(&j))

	resp, cleanup := client.Post("/v2/specs/"+j.ID+"/runs", bytes.NewBufferString(`{`))
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 500)
}

func TestJobRunsController_Create_WithoutWebInitiator(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	j := cltest.NewJob()
	assert.Nil(t, app.Store.SaveJob(&j))

	resp, cleanup := client.Post("/v2/specs/"+j.ID+"/runs", bytes.NewBuffer([]byte{}))
	defer cleanup()
	assert.Equal(t, 403, resp.StatusCode, "Response should be forbidden")
}

func TestJobRunsController_Create_NotFound(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	resp, cleanup := client.Post("/v2/specs/garbageID/runs", bytes.NewBuffer([]byte{}))
	defer cleanup()
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
}

func TestJobRunsController_Update_Success(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()

	bt := cltest.NewBridgeType()
	assert.Nil(t, app.Store.Save(&bt))
	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.Nil(t, app.Store.Save(&j))
	jr := cltest.MarkJobRunPendingBridge(j.NewRun(initr), 0)
	assert.Nil(t, app.Store.Save(&jr))

	body := fmt.Sprintf(`{"id":"%v","data":{"value": "100"}}`, jr.ID)
	headers := map[string]string{"Authorization": "Bearer " + bt.IncomingToken}
	url := app.Config.ClientNodeURL + "/v2/runs/" + jr.ID
	resp, cleanup := cltest.UnauthenticatedPatch(url, bytes.NewBufferString(body), headers)
	defer cleanup()

	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	jrID := cltest.ParseCommonJSON(resp.Body).ID
	assert.Equal(t, jr.ID, jrID)

	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)
	val, err := jr.Result.Value()
	assert.NoError(t, err)
	assert.Equal(t, "100", val)
}

func TestJobRunsController_Update_WrongAccessToken(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	bt := cltest.NewBridgeType()
	assert.Nil(t, app.Store.Save(&bt))
	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.Nil(t, app.Store.Save(&j))
	jr := cltest.MarkJobRunPendingBridge(j.NewRun(initr), 0)
	assert.Nil(t, app.Store.Save(&jr))

	body := fmt.Sprintf(`{"id":"%v","data":{"value": "100"}}`, jr.ID)
	headers := map[string]string{"Authorization": "Bearer " + "wrongaccesstoken"}
	resp, cleanup := client.Patch("/v2/runs/"+jr.ID, bytes.NewBufferString(body), headers)
	defer cleanup()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Response should be unauthorized")
	assert.Nil(t, app.Store.One("ID", jr.ID, &jr))
	assert.Equal(t, models.RunStatusPendingBridge, jr.Status)
}

func TestJobRunsController_Update_NotPending(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	bt := cltest.NewBridgeType()
	assert.Nil(t, app.Store.Save(&bt))
	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.Nil(t, app.Store.Save(&j))
	jr := j.NewRun(initr)
	assert.Nil(t, app.Store.Save(&jr))

	body := fmt.Sprintf(`{"id":"%v","data":{"value": "100"}}`, jr.ID)
	headers := map[string]string{"Authorization": "Bearer " + bt.IncomingToken}
	resp, cleanup := client.Patch("/v2/runs/"+jr.ID, bytes.NewBufferString(body), headers)
	defer cleanup()
	assert.Equal(t, 405, resp.StatusCode, "Response should be unsuccessful")
}

func TestJobRunsController_Update_WithError(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	bt := cltest.NewBridgeType()
	assert.Nil(t, app.Store.Save(&bt))
	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.Nil(t, app.Store.Save(&j))
	jr := cltest.MarkJobRunPendingBridge(j.NewRun(initr), 0)
	assert.Nil(t, app.Store.Save(&jr))

	body := fmt.Sprintf(`{"id":"%v","error":"stack overflow","data":{"value": "0"}}`, jr.ID)
	headers := map[string]string{"Authorization": "Bearer " + bt.IncomingToken}
	resp, cleanup := client.Patch("/v2/runs/"+jr.ID, bytes.NewBufferString(body), headers)
	defer cleanup()
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	jrID := cltest.ParseCommonJSON(resp.Body).ID
	assert.Equal(t, jr.ID, jrID)

	jr = cltest.WaitForJobRunStatus(t, app.Store, jr, models.RunStatusErrored)
	val, err := jr.Result.Value()
	assert.NoError(t, err)
	assert.Equal(t, "0", val)
}

func TestJobRunsController_Update_WithMergeError(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()

	bt := cltest.NewBridgeType()
	assert.Nil(t, app.Store.Save(&bt))
	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.Nil(t, app.Store.Save(&j))
	jr := cltest.MarkJobRunPendingBridge(j.NewRun(initr), 0)
	jr.Overrides = jr.Overrides.WithError(errors.New("Already errored")) // easy way to force Merge error
	assert.Nil(t, app.Store.Save(&jr))

	body := fmt.Sprintf(`{"id":"%v","data":{"value": "100"}}`, jr.ID)
	headers := map[string]string{"Authorization": "Bearer " + bt.IncomingToken}
	url := app.Config.ClientNodeURL + "/v2/runs/" + jr.ID
	resp, cleanup := cltest.UnauthenticatedPatch(url, bytes.NewBufferString(body), headers)
	defer cleanup()

	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	jrID := cltest.ParseCommonJSON(resp.Body).ID
	assert.Equal(t, jr.ID, jrID)

	jr = cltest.WaitForJobRunStatus(t, app.Store, jr, models.RunStatusErrored)
	assert.Contains(t, jr.Result.ErrorMessage.String, "Cannot merge")
}

func TestJobRunsController_Update_BadInput(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	bt := cltest.NewBridgeType()
	assert.Nil(t, app.Store.Save(&bt))
	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.Nil(t, app.Store.Save(&j))
	jr := cltest.MarkJobRunPendingBridge(j.NewRun(initr), 0)
	assert.Nil(t, app.Store.Save(&jr))

	body := fmt.Sprint(`{`, jr.ID)
	resp, cleanup := client.Patch("/v2/runs/"+jr.ID, bytes.NewBufferString(body))
	defer cleanup()
	assert.Equal(t, 500, resp.StatusCode, "Response should be successful")
	assert.Nil(t, app.Store.One("ID", jr.ID, &jr))
	assert.Equal(t, models.RunStatusPendingBridge, jr.Status)
}

func TestJobRunsController_Update_NotFound(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	bt := cltest.NewBridgeType()
	assert.Nil(t, app.Store.Save(&bt))
	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.Nil(t, app.Store.Save(&j))
	jr := cltest.MarkJobRunPendingBridge(j.NewRun(initr), 0)
	assert.Nil(t, app.Store.Save(&jr))

	body := fmt.Sprintf(`{"id":"%v","data":{"value": "100"}}`, jr.ID)
	resp, cleanup := client.Patch("/v2/runs/"+jr.ID+"1", bytes.NewBufferString(body))
	defer cleanup()
	assert.Equal(t, 404, resp.StatusCode, "Response should be successful")
	assert.Nil(t, app.Store.One("ID", jr.ID, &jr))
	assert.Equal(t, models.RunStatusPendingBridge, jr.Status)
}

func TestJobRunsController_Show_Found(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	j, initr := cltest.NewJobWithSchedule("9 9 9 9 6")
	app.Store.SaveJob(&j)

	jr := j.NewRun(initr)
	jr.ID = "jobrun1"
	assert.Nil(t, app.Store.Save(&jr))

	resp, cleanup := client.Get("/v2/runs/" + jr.ID)
	defer cleanup()
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")

	var respJobRun presenters.JobRun
	assert.NoError(t, cltest.ParseJSONAPIResponse(resp, &respJobRun))
	assert.Equal(t, jr.Initiator.Schedule, respJobRun.Initiator.Schedule, "should have the same schedule")
	assert.Equal(t, jr.ID, respJobRun.ID, "should have job run id")
}

func TestJobRunsController_Show_NotFound(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()
	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v2/runs/garbage")
	defer cleanup()
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
}

func TestJobRunsController_Show_Unauthenticated(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()

	resp, err := http.Get(app.Server.URL + "/v2/runs/notauthorized")
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode, "Response should be forbidden")
}
