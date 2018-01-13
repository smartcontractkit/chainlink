package web_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

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
	app.Store.Save(&j1)
	j2 := cltest.NewJobWithWebInitiator()
	app.Store.Save(&j2)

	resp, err := cltest.BasicAuthGet(app.Server.URL + "/v2/jobs")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var jobs []models.Job
	json.Unmarshal(b, &jobs)
	assert.Equal(t, jobs[0].Initiators[0].Schedule, j1.Initiators[0].Schedule, "should have the same schedule")
	assert.Equal(t, jobs[1].Initiators[0].Type, "web", "should have the same type")
}

func TestCreateJobs(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/hello_world_job.json")
	resp, _ := cltest.BasicAuthPost(app.Server.URL+"/v2/jobs", "application/json", bytes.NewBuffer(jsonStr))
	defer resp.Body.Close()
	respJSON := cltest.JobJSONFromResponse(resp.Body)
	assert.Equal(t, 200, resp.StatusCode, "Response should be success")

	var j models.Job
	app.Store.One("ID", respJSON.ID, &j)
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
	resp, err := cltest.BasicAuthPost(
		app.Server.URL+"/v2/jobs",
		"application/json",
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 500, resp.StatusCode, "Response should be internal error")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, `{"errors":["IdoNotExist is not a supported adapter type"]}`, string(body), "Response should return JSON")
}

func TestCreateInvalidCron(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/invalid_cron.json")
	resp, err := cltest.BasicAuthPost(
		app.Server.URL+"/v2/jobs",
		"application/json",
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 500, resp.StatusCode, "Response should be internal error")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, `{"errors":["Cron: Failed to parse int from !: strconv.Atoi: parsing \"!\": invalid syntax"]}`, string(body), "Response should return JSON")
}

func TestShowJobs(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j := cltest.NewJobWithSchedule("9 9 9 9 6")
	app.Store.Save(&j)
	jr := j.NewRun()
	app.Store.Save(&jr)

	resp, err := cltest.BasicAuthGet(app.Server.URL + "/v2/jobs/" + j.ID)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var respJob web.JobPresenter
	json.Unmarshal(b, &respJob)
	assert.Equal(t, respJob.Initiators[0].Schedule, j.Initiators[0].Schedule, "should have the same schedule")
	assert.Equal(t, respJob.Runs[0].ID, jr.ID, "should have the job runs")
}

func TestShowNotFoundJobs(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp, err := cltest.BasicAuthGet(app.Server.URL + "/v2/jobs/" + "garbage")
	assert.Nil(t, err)
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
