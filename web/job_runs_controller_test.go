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

func TestJobRunsController_Index(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j := cltest.NewJob()
	assert.Nil(t, app.Store.SaveJob(&j))
	jr1 := j.NewRun()
	jr1.ID = "2"
	assert.Nil(t, app.Store.Save(&jr1))
	jr2 := j.NewRun()
	jr2.ID = "1"
	jr2.CreatedAt = jr1.CreatedAt.Add(time.Second)
	assert.Nil(t, app.Store.Save(&jr2))

	resp := cltest.BasicAuthGet(app.Server.URL + "/v2/jobs/" + j.ID + "/runs")
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	var respJSON JobRunsJSON
	assert.Nil(t, json.Unmarshal(cltest.ParseResponseBody(resp), &respJSON))

	assert.Equal(t, 2, len(respJSON.Runs), "expected no runs to be created")
	assert.Equal(t, jr2.ID, respJSON.Runs[0].ID, "expected runs ordered by created at(descending)")
	assert.Equal(t, jr1.ID, respJSON.Runs[1].ID, "expected runs ordered by created at(descending)")
}

func TestJobRunsController_Create(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.SaveJob(&j))

	jr := cltest.CreateJobRunViaWeb(t, app, j)
	cltest.WaitForJobRunToComplete(t, app, jr)
}

func TestJobRunsController_Create_WithoutWebInitiator(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j := cltest.NewJob()
	assert.Nil(t, app.Store.SaveJob(&j))

	url := app.Server.URL + "/v2/jobs/" + j.ID + "/runs"
	resp := cltest.BasicAuthPost(url, "application/json", bytes.NewBuffer([]byte{}))
	assert.Equal(t, 403, resp.StatusCode, "Response should be forbidden")
}

func TestJobRunsController_Create_NotFound(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	url := app.Server.URL + "/v2/jobs/garbageID/runs"
	resp := cltest.BasicAuthPost(url, "application/json", bytes.NewBuffer([]byte{}))
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
}

func TestJobRunsController_Update(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	bt := models.BridgeType{
		Name: "slowcomputation",
		URL:  cltest.WebURL("http://localhost:12345"),
	}
	assert.Nil(t, app.Store.Save(&bt))
	j := cltest.NewJob()
	j.Tasks = []models.Task{{
		Type:   bt.Name,
		Params: cltest.JSONFromString(`{"type":"%v"}`, bt.Name),
	}}
	assert.Nil(t, app.Store.Save(&j))
	jr := j.NewRun()
	jr.Status = models.StatusPending
	jr.TaskRuns[0].Status = models.StatusPending
	jr.TaskRuns[0].Result.Pending = true
	assert.Nil(t, app.Store.Save(&jr))

	url := app.Server.URL + "/v2/runs/" + jr.ID
	body := fmt.Sprintf(`{"id":"%v","data":{"value": "100"}}`, jr.ID)
	resp := cltest.BasicAuthPatch(url, "application/json", bytes.NewBufferString(body))
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	jrID := cltest.ParseCommonJSON(resp.Body).ID
	assert.Equal(t, jr.ID, jrID)

	jr = cltest.WaitForJobRunToComplete(t, app, jr)
	val, err := jr.Result.Value()
	assert.Nil(t, err)
	assert.Equal(t, "100", val)
}
