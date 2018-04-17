package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
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
	defer cleanup()

	j := setupJobRunsControllerIndex(t, app)
	resp := cltest.BasicAuthGet(app.Server.URL + "/v2/specs/" + j.ID + "/runs")

	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	var respJSON JobRunsJSON
	assert.Nil(t, json.Unmarshal(cltest.ParseResponseBody(resp), &respJSON))

	jr, err := app.Store.JobRunsFor(j.ID)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(respJSON.Runs), "expected no runs to be created")
	assert.Equal(t, jr[0].ID, respJSON.Runs[0].ID, "expected runs ordered by created at(descending)")
	assert.Equal(t, jr[1].ID, respJSON.Runs[1].ID, "expected runs ordered by created at(descending)")
}

func setupJobRunsControllerIndex(t assert.TestingT, app *cltest.TestApplication) *models.JobSpec {
	j, initr := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.SaveJob(&j))
	jr1 := j.NewRun(initr)
	jr1.ID = "2"
	assert.Nil(t, app.Store.Save(&jr1))
	jr2 := j.NewRun(initr)
	jr2.ID = "1"
	jr2.CreatedAt = jr1.CreatedAt.Add(time.Second)
	assert.Nil(t, app.Store.Save(&jr2))

	return &j
}

func TestJobRunsController_Create_Success(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j, _ := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.SaveJob(&j))

	jr := cltest.CreateJobRunViaWeb(t, app, j, `{"result":"100"}`)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)
	val, err := jr.Result.Result()
	assert.Nil(t, err)
	assert.Equal(t, "100", val)
}

func TestJobRunsController_Create_EmptyBody(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j, _ := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.SaveJob(&j))

	jr := cltest.CreateJobRunViaWeb(t, app, j)
	cltest.WaitForJobRunToComplete(t, app.Store, jr)
}

func TestJobRunsController_Create_InvalidBody(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j, _ := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.SaveJob(&j))

	url := app.Server.URL + "/v2/specs/" + j.ID + "/runs"
	resp := cltest.BasicAuthPost(url, "application/json", bytes.NewBufferString(`{`))
	defer resp.Body.Close()
	cltest.CheckStatusCode(t, resp, 500)
}

func TestJobRunsController_Create_WithoutWebInitiator(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
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
	defer cleanup()

	url := app.Server.URL + "/v2/specs/garbageID/runs"
	resp := cltest.BasicAuthPost(url, "application/json", bytes.NewBuffer([]byte{}))
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
}

func TestJobRunsController_Update_Success(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	bt := cltest.NewBridgeType()
	assert.Nil(t, app.Store.Save(&bt))
	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{cltest.NewTask(bt.Name)}
	assert.Nil(t, app.Store.Save(&j))
	jr := cltest.MarkJobRunPendingBridge(j.NewRun(initr), 0)
	assert.Nil(t, app.Store.Save(&jr))

	url := app.Server.URL + "/v2/runs/" + jr.ID
	body := fmt.Sprintf(`{"id":"%v","data":{"result": "100"}}`, jr.ID)
	resp := cltest.BasicAuthPatch(url, "application/json", bytes.NewBufferString(body))
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	jrID := cltest.ParseCommonJSON(resp.Body).ID
	assert.Equal(t, jr.ID, jrID)

	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)
	val, err := jr.Result.Result()
	assert.Nil(t, err)
	assert.Equal(t, "100", val)
}

func TestJobRunsController_Update_NotPending(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	bt := cltest.NewBridgeType()
	assert.Nil(t, app.Store.Save(&bt))
	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{cltest.NewTask(bt.Name)}
	assert.Nil(t, app.Store.Save(&j))
	jr := j.NewRun(initr)
	assert.Nil(t, app.Store.Save(&jr))

	url := app.Server.URL + "/v2/runs/" + jr.ID
	body := fmt.Sprintf(`{"id":"%v","data":{"result": "100"}}`, jr.ID)
	resp := cltest.BasicAuthPatch(url, "application/json", bytes.NewBufferString(body))
	assert.Equal(t, 405, resp.StatusCode, "Response should be unsuccessful")
}

func TestJobRunsController_Update_WithError(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	bt := cltest.NewBridgeType()
	assert.Nil(t, app.Store.Save(&bt))
	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{cltest.NewTask(bt.Name)}
	assert.Nil(t, app.Store.Save(&j))
	jr := cltest.MarkJobRunPendingBridge(j.NewRun(initr), 0)
	assert.Nil(t, app.Store.Save(&jr))

	url := app.Server.URL + "/v2/runs/" + jr.ID
	body := fmt.Sprintf(`{"id":"%v","error":"stack overflow","data":{"result": "0"}}`, jr.ID)
	resp := cltest.BasicAuthPatch(url, "application/json", bytes.NewBufferString(body))
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	jrID := cltest.ParseCommonJSON(resp.Body).ID
	assert.Equal(t, jr.ID, jrID)

	jr = cltest.WaitForJobRunStatus(t, app.Store, jr, models.RunStatusErrored)
	val, err := jr.Result.Result()
	assert.Nil(t, err)
	assert.Equal(t, "0", val)
}

func TestJobRunsController_Update_BadInput(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	bt := cltest.NewBridgeType()
	assert.Nil(t, app.Store.Save(&bt))
	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{cltest.NewTask(bt.Name)}
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
	defer cleanup()

	bt := cltest.NewBridgeType()
	assert.Nil(t, app.Store.Save(&bt))
	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{cltest.NewTask(bt.Name)}
	assert.Nil(t, app.Store.Save(&j))
	jr := cltest.MarkJobRunPendingBridge(j.NewRun(initr), 0)
	assert.Nil(t, app.Store.Save(&jr))

	url := app.Server.URL + "/v2/runs/" + jr.ID + "1"
	body := fmt.Sprintf(`{"id":"%v","data":{"result": "100"}}`, jr.ID)
	resp := cltest.BasicAuthPatch(url, "application/json", bytes.NewBufferString(body))
	assert.Equal(t, 404, resp.StatusCode, "Response should be successful")
	assert.Nil(t, app.Store.One("ID", jr.ID, &jr))
	assert.Equal(t, models.RunStatusPendingBridge, jr.Status)
}
