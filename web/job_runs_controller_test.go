package web_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
	j := setupJobRunsControllerIndex(b, app)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		resp := cltest.BasicAuthGet(app.Server.URL + "/v2/specs/" + j.ID + "/runs")
		assert.Equal(b, 200, resp.StatusCode, "Response should be successful")
	}
}

func TestJobRunsController_Index(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()

	j := setupJobRunsControllerIndex(t, app)
	jr, err := app.Store.JobRunsFor(j.ID)
	assert.NoError(t, err)

	resp := cltest.BasicAuthGet(app.Server.URL + "/v2/specs/" + j.ID + "/runs?size=x")
	cltest.AssertServerResponse(t, resp, 422)

	resp = cltest.BasicAuthGet(app.Server.URL + "/v2/specs/" + j.ID + "/runs?size=1")
	cltest.AssertServerResponse(t, resp, 200)

	var links jsonapi.Links
	var runs []models.JobRun

	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(resp), &runs, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, runs, 1)
	assert.Equal(t, jr[1].ID, runs[0].ID, "expected runs ordered by created at(descending)")

	resp = cltest.BasicAuthGet(app.Server.URL + links["next"].Href)
	cltest.AssertServerResponse(t, resp, 200)

	runs = []models.JobRun{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(resp), &runs, &links)
	assert.NoError(t, err)
	assert.Empty(t, links["next"])
	assert.NotEmpty(t, links["prev"])

	assert.Len(t, runs, 1)
	assert.Equal(t, jr[0].ID, runs[0].ID, "expected runs ordered by created at(descending)")
}

func setupJobRunsControllerIndex(t assert.TestingT, app *cltest.TestApplication) *models.JobSpec {
	j, initr := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.SaveJob(&j))
	jr1 := j.NewRun(initr)
	jr1.ID = "runB"
	assert.Nil(t, app.Store.Save(&jr1))
	jr2 := j.NewRun(initr)
	jr2.ID = "runA"
	jr2.CreatedAt = jr1.CreatedAt.Add(time.Second)
	assert.Nil(t, app.Store.Save(&jr2))

	j2, initr := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.SaveJob(&j2))
	jr3 := j2.NewRun(initr)
	jr3.ID = "runC"
	assert.Nil(t, app.Store.Save(&jr3))

	return &j
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

	j, _ := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.SaveJob(&j))

	url := app.Server.URL + "/v2/specs/" + j.ID + "/runs"
	resp := cltest.BasicAuthPost(url, "application/json", bytes.NewBufferString(`{`))
	defer resp.Body.Close()
	cltest.AssertServerResponse(t, resp, 500)
}

func TestJobRunsController_Create_WithoutWebInitiator(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()

	j := cltest.NewJob()
	assert.Nil(t, app.Store.SaveJob(&j))

	url := app.Server.URL + "/v2/specs/" + j.ID + "/runs"
	resp := cltest.BasicAuthPost(url, "application/json", bytes.NewBuffer([]byte{}))
	assert.Equal(t, 403, resp.StatusCode, "Response should be forbidden")
}

func TestJobRunsController_Create_NotFound(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()

	url := app.Server.URL + "/v2/specs/garbageID/runs"
	resp := cltest.BasicAuthPost(url, "application/json", bytes.NewBuffer([]byte{}))
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
	j.Tasks = []models.TaskSpec{{Type: models.NewTaskType(bt.Name)}}
	assert.Nil(t, app.Store.Save(&j))
	jr := cltest.MarkJobRunPendingBridge(j.NewRun(initr), 0)
	assert.Nil(t, app.Store.Save(&jr))

	url := app.Server.URL + "/v2/runs/" + jr.ID
	body := fmt.Sprintf(`{"id":"%v","data":{"value": "100"}}`, jr.ID)
	resp := cltest.BasicAuthPatch(url, "application/json", bytes.NewBufferString(body))
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	jrID := cltest.ParseCommonJSON(resp.Body).ID
	assert.Equal(t, jr.ID, jrID)

	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)
	val, err := jr.Result.Value()
	assert.NoError(t, err)
	assert.Equal(t, "100", val)
}

func TestJobRunsController_Update_NotPending(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()

	bt := cltest.NewBridgeType()
	assert.Nil(t, app.Store.Save(&bt))
	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: models.NewTaskType(bt.Name)}}
	assert.Nil(t, app.Store.Save(&j))
	jr := j.NewRun(initr)
	assert.Nil(t, app.Store.Save(&jr))

	url := app.Server.URL + "/v2/runs/" + jr.ID
	body := fmt.Sprintf(`{"id":"%v","data":{"value": "100"}}`, jr.ID)
	resp := cltest.BasicAuthPatch(url, "application/json", bytes.NewBufferString(body))
	assert.Equal(t, 405, resp.StatusCode, "Response should be unsuccessful")
}

func TestJobRunsController_Update_WithError(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()

	bt := cltest.NewBridgeType()
	assert.Nil(t, app.Store.Save(&bt))
	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: models.NewTaskType(bt.Name)}}
	assert.Nil(t, app.Store.Save(&j))
	jr := cltest.MarkJobRunPendingBridge(j.NewRun(initr), 0)
	assert.Nil(t, app.Store.Save(&jr))

	url := app.Server.URL + "/v2/runs/" + jr.ID
	body := fmt.Sprintf(`{"id":"%v","error":"stack overflow","data":{"value": "0"}}`, jr.ID)
	resp := cltest.BasicAuthPatch(url, "application/json", bytes.NewBufferString(body))
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	jrID := cltest.ParseCommonJSON(resp.Body).ID
	assert.Equal(t, jr.ID, jrID)

	jr = cltest.WaitForJobRunStatus(t, app.Store, jr, models.RunStatusErrored)
	val, err := jr.Result.Value()
	assert.NoError(t, err)
	assert.Equal(t, "0", val)
}

func TestJobRunsController_Update_BadInput(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()

	bt := cltest.NewBridgeType()
	assert.Nil(t, app.Store.Save(&bt))
	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: models.NewTaskType(bt.Name)}}
	assert.Nil(t, app.Store.Save(&j))
	jr := cltest.MarkJobRunPendingBridge(j.NewRun(initr), 0)
	assert.Nil(t, app.Store.Save(&jr))

	url := app.Server.URL + "/v2/runs/" + jr.ID
	body := fmt.Sprint(`{`, jr.ID)
	resp := cltest.BasicAuthPatch(url, "application/json", bytes.NewBufferString(body))
	assert.Equal(t, 500, resp.StatusCode, "Response should be successful")
	assert.Nil(t, app.Store.One("ID", jr.ID, &jr))
	assert.Equal(t, models.RunStatusPendingBridge, jr.Status)
}

func TestJobRunsController_Update_NotFound(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()

	bt := cltest.NewBridgeType()
	assert.Nil(t, app.Store.Save(&bt))
	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{cltest.NewTask(bt.Name)}
	assert.Nil(t, app.Store.Save(&j))
	jr := cltest.MarkJobRunPendingBridge(j.NewRun(initr), 0)
	assert.Nil(t, app.Store.Save(&jr))

	url := app.Server.URL + "/v2/runs/" + jr.ID + "1"
	body := fmt.Sprintf(`{"id":"%v","data":{"value": "100"}}`, jr.ID)
	resp := cltest.BasicAuthPatch(url, "application/json", bytes.NewBufferString(body))
	assert.Equal(t, 404, resp.StatusCode, "Response should be successful")
	assert.Nil(t, app.Store.One("ID", jr.ID, &jr))
	assert.Equal(t, models.RunStatusPendingBridge, jr.Status)
}

func TestJobRunsController_Show_Found(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()

	j, initr := cltest.NewJobWithSchedule("9 9 9 9 6")
	app.Store.SaveJob(&j)

	jr := j.NewRun(initr)
	jr.ID = "jobrun1"
	assert.Nil(t, app.Store.Save(&jr))

	resp := cltest.BasicAuthGet(app.Server.URL + "/v2/runs/" + jr.ID)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")

	var respJobRun presenters.JobRun
	b, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NoError(t, web.ParseJSONAPIResponse(b, &respJobRun))
	assert.Equal(t, jr.Initiator.Schedule, respJobRun.Initiator.Schedule, "should have the same schedule")
	assert.Equal(t, jr.ID, respJobRun.ID, "should have job run id")
}

func TestJobRunsController_Show_NotFound(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	app.Start()
	defer cleanup()

	resp := cltest.BasicAuthGet(app.Server.URL + "/v2/runs/garbage")
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
