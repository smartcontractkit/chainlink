package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/web"
	"github.com/stretchr/testify/assert"
)

func TestIndexJobs(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j1 := cltest.NewJobWithSchedule("9 9 9 9 6")
	j1.CreatedAt = models.Time{time.Now().AddDate(0, 0, -1)}
	app.Store.SaveJob(j1)
	j2 := cltest.NewJobWithWebInitiator()
	app.Store.Save(j2)

	resp := cltest.BasicAuthGet(app.Server.URL + "/v2/jobs")
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")

	var jobs []models.Job
	json.Unmarshal(cltest.ParseResponseBody(resp), &jobs)
	assert.Equal(t, j1.Initiators[0].Schedule, jobs[0].Initiators[0].Schedule, "should have the same schedule")
	assert.Equal(t, "web", jobs[1].Initiators[0].Type, "should have the same type")
}

func TestCreateJobs(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/hello_world_job.json")
	resp := cltest.BasicAuthPost(app.Server.URL+"/v2/jobs", "application/json", bytes.NewBuffer(jsonStr))
	defer resp.Body.Close()
	respJSON := cltest.JobJSONFromResponse(resp.Body)
	assert.Equal(t, 200, resp.StatusCode, "Response should be success")

	j, _ := app.Store.FindJob(respJSON.ID)
	assert.Equal(t, j.ID, respJSON.ID, "Wrong job returned")

	adapter1, _ := adapters.For(j.Tasks[0])
	httpGet := adapter1.(*adapters.HttpGet)
	assert.Equal(t, httpGet.Endpoint, "https://bitstamp.net/api/ticker/")

	adapter2, _ := adapters.For(j.Tasks[1])
	jsonParse := adapter2.(*adapters.JsonParse)
	assert.Equal(t, jsonParse.Path, []string{"last"})

	adapter4, _ := adapters.For(j.Tasks[3])
	signTx := adapter4.(*adapters.EthTx)
	assert.Equal(t, signTx.Address, "0x356a04bce728ba4c62a30294a55e6a8600a320b3")
	assert.Equal(t, signTx.FunctionID, "12345679")

	var initr models.Initiator
	app.Store.One("JobID", j.ID, &initr)
	assert.Equal(t, "web", initr.Type)
}

func TestCreateInvalidJobs(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/invalid_job.json")
	resp := cltest.BasicAuthPost(
		app.Server.URL+"/v2/jobs",
		"application/json",
		bytes.NewBuffer(jsonStr),
	)

	assert.Equal(t, 500, resp.StatusCode, "Response should be internal error")

	expected := `{"errors":["IdoNotExist is not a supported adapter type"]}`
	assert.Equal(t, expected, string(cltest.ParseResponseBody(resp)))
}

func TestCreateInvalidCron(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/invalid_cron.json")
	resp := cltest.BasicAuthPost(
		app.Server.URL+"/v2/jobs",
		"application/json",
		bytes.NewBuffer(jsonStr),
	)

	assert.Equal(t, 500, resp.StatusCode, "Response should be internal error")

	expected := `{"errors":["Cron: Failed to parse int from !: strconv.Atoi: parsing \"!\": invalid syntax"]}`
	assert.Equal(t, expected, string(cltest.ParseResponseBody(resp)))
}

func TestShowJobs(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j := cltest.NewJobWithSchedule("9 9 9 9 6")
	app.Store.SaveJob(j)
	jr := j.NewRun()
	app.Store.Save(jr)

	resp := cltest.BasicAuthGet(app.Server.URL + "/v2/jobs/" + j.ID)
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")

	var respJob web.JobPresenter
	json.Unmarshal(cltest.ParseResponseBody(resp), &respJob)
	assert.Equal(t, respJob.Initiators[0].Schedule, j.Initiators[0].Schedule, "should have the same schedule")
	assert.Equal(t, respJob.Runs[0].ID, jr.ID, "should have the job runs")
}

func TestShowNotFoundJobs(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp := cltest.BasicAuthGet(app.Server.URL + "/v2/jobs/" + "garbage")
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
}

func TestShowJobUnauthenticated(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp, err := http.Get(app.Server.URL + "/v2/jobs/" + "garbage")
	assert.Nil(t, err)
	assert.Equal(t, 401, resp.StatusCode, "Response should be forbidden")
}
