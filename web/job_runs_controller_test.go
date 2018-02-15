package web_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
)

type JobRunsJSON struct {
	Runs []JobRun `json:"runs"`
}

type JobRun struct {
	ID string `json:"id"`
}

func TestJobRunsIndex(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j := cltest.NewJobWithSchedule("9 9 9 9 6")
	assert.Nil(t, app.Store.SaveJob(&j))
	jr := j.NewRun()
	assert.Nil(t, app.Store.Save(&jr))

	resp := cltest.BasicAuthGet(app.Server.URL + "/v2/jobs/" + j.ID + "/runs")
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")

	var respJSON JobRunsJSON
	json.Unmarshal(cltest.ParseResponseBody(resp), &respJSON)
	assert.Equal(t, 1, len(respJSON.Runs), "expected no runs to be created")
	assert.Equal(t, jr.ID, respJSON.Runs[0].ID, "expected the run IDs to match")
}

func TestJobRunsCreateSuccessfully(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.SaveJob(&j))

	jr := cltest.CreateJobRunViaWeb(t, app, j)
	cltest.WaitForJobRunToComplete(t, app, jr)
}

func TestJobRunsCreateWithoutWebInitiator(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j := cltest.NewJob()
	assert.Nil(t, app.Store.SaveJob(&j))

	url := app.Server.URL + "/v2/jobs/" + j.ID + "/runs"
	resp := cltest.BasicAuthPost(url, "application/json", bytes.NewBuffer([]byte{}))
	assert.Equal(t, 403, resp.StatusCode, "Response should be forbidden")
}

func TestJobRunsCreateNotFound(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	url := app.Server.URL + "/v2/jobs/garbageID/runs"
	resp := cltest.BasicAuthPost(url, "application/json", bytes.NewBuffer([]byte{}))
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
}
