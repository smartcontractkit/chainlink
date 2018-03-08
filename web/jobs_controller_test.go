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
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/stretchr/testify/assert"
)

func TestJobSpecsController_Index(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j1 := cltest.NewJobWithSchedule("9 9 9 9 6")
	j1.CreatedAt = models.Time{time.Now().AddDate(0, 0, -1)}
	app.Store.SaveJob(&j1)
	j2 := cltest.NewJobWithWebInitiator()
	j2.Initiators[0].Ran = true
	app.Store.SaveJob(&j2)

	resp := cltest.BasicAuthGet(app.Server.URL + "/v2/specs")
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")

	var jobs []models.JobSpec
	json.Unmarshal(cltest.ParseResponseBody(resp), &jobs)
	assert.Equal(t, j1.Initiators[0].Schedule, jobs[0].Initiators[0].Schedule, "should have the same schedule")
	assert.Equal(t, models.InitiatorWeb, jobs[1].Initiators[0].Type, "should have the same type")
	assert.NotEqual(t, true, jobs[1].Initiators[0].Ran, "should ignore fields for other initiators")
}

func TestJobSpecsController_Create(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/hello_world_job.json")

	adapter1, _ := adapters.For(j.Tasks[0], app.Store)
	httpGet := adapter1.(*adapters.HTTPGet)
	assert.Equal(t, httpGet.URL.String(), "https://bitstamp.net/api/ticker/")

	adapter2, _ := adapters.For(j.Tasks[1], app.Store)
	jsonParse := adapter2.(*adapters.JSONParse)
	assert.Equal(t, jsonParse.Path, []string{"last"})

	adapter4, _ := adapters.For(j.Tasks[3], app.Store)
	signTx := adapter4.(*adapters.EthTx)
	assert.Equal(t, "0x356a04bCe728ba4c62A30294A55E6A8600a320B3", signTx.Address.String())
	assert.Equal(t, "0x609ff1bd", signTx.FunctionSelector.String())

	var initr models.Initiator
	app.Store.One("JobID", j.ID, &initr)
	assert.Equal(t, models.InitiatorWeb, initr.Type)
}

func TestJobSpecsController_Create_CaseInsensitiveTypes(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/caseinsensitive_hello_world_job.json")

	adapter1, _ := adapters.For(j.Tasks[0], app.Store)
	httpGet := adapter1.(*adapters.HTTPGet)
	assert.Equal(t, httpGet.URL.String(), "https://bitstamp.net/api/ticker/")

	adapter2, _ := adapters.For(j.Tasks[1], app.Store)
	jsonParse := adapter2.(*adapters.JSONParse)
	assert.Equal(t, jsonParse.Path, []string{"last"})

	assert.Equal(t, "ethbytes32", j.Tasks[2].Type)

	adapter4, _ := adapters.For(j.Tasks[3], app.Store)
	signTx := adapter4.(*adapters.EthTx)
	assert.Equal(t, "0x356a04bCe728ba4c62A30294A55E6A8600a320B3", signTx.Address.String())
	assert.Equal(t, "0x609ff1bd", signTx.FunctionSelector.String())

	assert.Equal(t, models.InitiatorWeb, j.Initiators[0].Type)
	assert.Equal(t, models.InitiatorRunAt, j.Initiators[1].Type)
}

func TestJobSpecsController_Create_NonExistentTaskJob(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/nonexistent_task_job.json")
	resp := cltest.BasicAuthPost(
		app.Server.URL+"/v2/specs",
		"application/json",
		bytes.NewBuffer(jsonStr),
	)

	assert.Equal(t, 400, resp.StatusCode, "Response should be caller error")

	expected := `{"errors":["job validation: task validation: idonotexist is not a supported adapter type"]}`
	assert.Equal(t, expected, string(cltest.ParseResponseBody(resp)))
}

func TestJobSpecsController_Create_InvalidJob(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/run_at_wo_time_job.json")
	resp := cltest.BasicAuthPost(
		app.Server.URL+"/v2/specs",
		"application/json",
		bytes.NewBuffer(jsonStr),
	)

	assert.Equal(t, 400, resp.StatusCode, "Response should be caller error")

	expected := `{"errors":["job validation: initiator validation: runat must have a time"]}`
	assert.Equal(t, expected, string(cltest.ParseResponseBody(resp)))
}

func TestJobSpecsController_Create_InvalidCron(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/invalid_cron.json")
	resp := cltest.BasicAuthPost(
		app.Server.URL+"/v2/specs",
		"application/json",
		bytes.NewBuffer(jsonStr),
	)

	assert.Equal(t, 400, resp.StatusCode, "Response should be caller error")

	expected := `{"errors":["Cron: Failed to parse int from !: strconv.Atoi: parsing \"!\": invalid syntax"]}`
	assert.Equal(t, expected, string(cltest.ParseResponseBody(resp)))
}

func TestJobSpecsController_Show(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j := cltest.NewJobWithSchedule("9 9 9 9 6")
	app.Store.SaveJob(&j)

	jr1 := j.NewRun()
	jr1.ID = "2"
	assert.Nil(t, app.Store.Save(&jr1))
	jr2 := j.NewRun()
	jr2.ID = "1"
	jr2.CreatedAt = jr1.CreatedAt.Add(time.Second)
	assert.Nil(t, app.Store.Save(&jr2))

	resp := cltest.BasicAuthGet(app.Server.URL + "/v2/specs/" + j.ID)
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")

	var respJob presenters.JobSpec
	json.Unmarshal(cltest.ParseResponseBody(resp), &respJob)
	assert.Equal(t, respJob.Initiators[0].Schedule, j.Initiators[0].Schedule, "should have the same schedule")
	assert.Equal(t, respJob.Runs[0].ID, jr2.ID, "should have job runs ordered by created at(descending)")
	assert.Equal(t, respJob.Runs[1].ID, jr1.ID, "should have job runs ordered by created at(descending)")
}

func TestJobSpecsController_Show_NotFound(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp := cltest.BasicAuthGet(app.Server.URL + "/v2/specs/" + "garbage")
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
}

func TestJobSpecsController_Show_Unauthenticated(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp, err := http.Get(app.Server.URL + "/v2/specs/" + "garbage")
	assert.Nil(t, err)
	assert.Equal(t, 401, resp.StatusCode, "Response should be forbidden")
}
