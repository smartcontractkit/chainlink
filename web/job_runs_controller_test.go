package web_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/store/models"
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
	app := cltest.NewApplication()
	server := app.NewServer()
	defer app.Stop()

	j := cltest.NewJobWithSchedule("9 9 9 9 6")
	assert.Nil(t, app.Store.Save(&j))
	jr := j.NewRun()
	assert.Nil(t, app.Store.Save(&jr))

	resp, err := cltest.BasicAuthGet(server.URL + "/jobs/" + j.ID + "/runs")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	var respJSON JobRunsJSON
	json.Unmarshal(b, &respJSON)
	assert.Equal(t, 1, len(respJSON.Runs), "expected no runs to be created")
	assert.Equal(t, jr.ID, respJSON.Runs[0].ID, "expected the run IDs to match")
}

func TestJobRunsCreate(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)

	app := cltest.NewApplication()
	server := app.NewServer()
	defer app.Stop()

	j := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.SaveJob(j))

	url := server.URL + "/jobs/" + j.ID + "/runs"
	resp, err := cltest.BasicAuthPost(url, "application/json", bytes.NewBuffer([]byte{}))
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	respJSON := cltest.JobJSONFromResponse(resp.Body)

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		app.Store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))

	jr := models.JobRun{}
	assert.Nil(t, app.Store.One("ID", respJSON.ID, &jr))
	assert.Equal(t, jr.ID, respJSON.ID)
}

func TestJobRunsCreateWithoutWebInitiator(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplication()
	server := app.NewServer()
	defer app.Stop()

	j := cltest.NewJobWithSchedule("* * * * *")
	assert.Nil(t, app.Store.SaveJob(j))

	url := server.URL + "/jobs/" + j.ID + "/runs"
	resp, err := cltest.BasicAuthPost(url, "application/json", bytes.NewBuffer([]byte{}))
	assert.Nil(t, err)
	assert.Equal(t, 403, resp.StatusCode, "Response should be forbidden")
}

func TestJobRunsCreateNotFound(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplication()
	server := app.NewServer()
	defer app.Stop()

	url := server.URL + "/jobs/garbageID/runs"
	resp, err := cltest.BasicAuthPost(url, "application/json", bytes.NewBuffer([]byte{}))
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
}
