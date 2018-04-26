package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/web"
	"github.com/stretchr/testify/assert"
)

func BenchmarkJobSpecsController_Index(b *testing.B) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	setupJobSpecsControllerIndex(app)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		resp := cltest.BasicAuthGet(app.Server.URL + "/v2/specs")
		assert.Equal(b, 200, resp.StatusCode, "Response should be successful")
	}
}

func TestJobSpecsController_Index(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j1, err := setupJobSpecsControllerIndex(app)
	assert.NoError(t, err)

	resp := cltest.BasicAuthGet(app.Server.URL + "/v2/specs?size=x")
	cltest.AssertServerResponse(t, resp, 422)

	resp = cltest.BasicAuthGet(app.Server.URL + "/v2/specs?size=1")
	cltest.AssertServerResponse(t, resp, 200)

	var links jsonapi.Links
	jobs := []models.JobSpec{}

	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(resp), &jobs, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, jobs, 1)
	assert.Equal(t, j1.Initiators[0].Schedule, jobs[0].Initiators[0].Schedule, "should have the same schedule")

	resp = cltest.BasicAuthGet(app.Server.URL + links["next"].Href)
	cltest.AssertServerResponse(t, resp, 200)

	jobs = []models.JobSpec{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(resp), &jobs, &links)
	assert.NoError(t, err)
	assert.Empty(t, links["next"])
	assert.NotEmpty(t, links["prev"])

	assert.Len(t, jobs, 1)
	assert.Equal(t, models.InitiatorWeb, jobs[0].Initiators[0].Type, "should have the same type")
	assert.NotEqual(t, true, jobs[0].Initiators[0].Ran, "should ignore fields for other initiators")
}

func setupJobSpecsControllerIndex(app *cltest.TestApplication) (*models.JobSpec, error) {
	j1, _ := cltest.NewJobWithSchedule("9 9 9 9 6")
	j1.CreatedAt = models.Time{Time: time.Now().AddDate(0, 0, -1)}
	err := app.Store.SaveJob(&j1)
	if err != nil {
		return nil, err
	}
	j2, _ := cltest.NewJobWithWebInitiator()
	j2.Initiators[0].Ran = true
	err = app.Store.SaveJob(&j2)
	return &j1, err
}

func TestJobSpecsController_Create(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/hello_world_job.json")

	adapter1, _ := adapters.For(j.Tasks[0], app.Store)
	httpGet := cltest.UnwrapAdapter(adapter1).(*adapters.HTTPGet)
	assert.Equal(t, httpGet.URL.String(), "https://bitstamp.net/api/ticker/")

	adapter2, _ := adapters.For(j.Tasks[1], app.Store)
	jsonParse := cltest.UnwrapAdapter(adapter2).(*adapters.JSONParse)
	assert.Equal(t, jsonParse.Path, []string{"last"})

	adapter4, _ := adapters.For(j.Tasks[3], app.Store)
	signTx := cltest.UnwrapAdapter(adapter4).(*adapters.EthTx)
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
	httpGet := cltest.UnwrapAdapter(adapter1).(*adapters.HTTPGet)
	assert.Equal(t, httpGet.URL.String(), "https://bitstamp.net/api/ticker/")

	adapter2, _ := adapters.For(j.Tasks[1], app.Store)
	jsonParse := cltest.UnwrapAdapter(adapter2).(*adapters.JSONParse)
	assert.Equal(t, jsonParse.Path, []string{"last"})

	assert.Equal(t, "ethbytes32", j.Tasks[2].Type)

	adapter4, _ := adapters.For(j.Tasks[3], app.Store)
	signTx := cltest.UnwrapAdapter(adapter4).(*adapters.EthTx)
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

func TestJobSpecsController_Create_Initiator_Only(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/initiator_only_job.json")
	resp := cltest.BasicAuthPost(
		app.Server.URL+"/v2/specs",
		"application/json",
		bytes.NewBuffer(jsonStr),
	)

	assert.Equal(t, 400, resp.StatusCode, "Response should be caller error")

	expected := `{"errors":["job validation: Must have at least one Initiator and one Task"]}`
	assert.Equal(t, expected, string(cltest.ParseResponseBody(resp)))
}

func TestJobSpecsController_Create_Task_Only(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/task_only_job.json")
	resp := cltest.BasicAuthPost(
		app.Server.URL+"/v2/specs",
		"application/json",
		bytes.NewBuffer(jsonStr),
	)

	assert.Equal(t, 400, resp.StatusCode, "Response should be caller error")

	expected := `{"errors":["job validation: Must have at least one Initiator and one Task"]}`
	assert.Equal(t, expected, string(cltest.ParseResponseBody(resp)))
}

func BenchmarkJobSpecsController_Show(b *testing.B) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	j := setupJobSpecsControllerShow(b, app)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		resp := cltest.BasicAuthGet(app.Server.URL + "/v2/specs/" + j.ID)
		assert.Equal(b, 200, resp.StatusCode, "Response should be successful")
	}
}

func TestJobSpecsController_Show(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j := setupJobSpecsControllerShow(t, app)

	jr, err := app.Store.JobRunsFor(j.ID)
	assert.Nil(t, err)

	resp := cltest.BasicAuthGet(app.Server.URL + "/v2/specs/" + j.ID)
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")

	var respJob presenters.JobSpec
	json.Unmarshal(cltest.ParseResponseBody(resp), &respJob)
	assert.Equal(t, respJob.Initiators[0].Schedule, j.Initiators[0].Schedule, "should have the same schedule")
	assert.Equal(t, respJob.Runs[0].ID, jr[0].ID, "should have job runs ordered by created at(descending)")
	assert.Equal(t, respJob.Runs[1].ID, jr[1].ID, "should have job runs ordered by created at(descending)")
}

func setupJobSpecsControllerShow(t assert.TestingT, app *cltest.TestApplication) *models.JobSpec {
	j, initr := cltest.NewJobWithSchedule("9 9 9 9 6")
	app.Store.SaveJob(&j)

	jr1 := j.NewRun(initr)
	jr1.ID = "2"
	assert.Nil(t, app.Store.Save(&jr1))
	jr2 := j.NewRun(initr)
	jr2.ID = "1"
	jr2.CreatedAt = jr1.CreatedAt.Add(time.Second)
	assert.Nil(t, app.Store.Save(&jr2))

	return &j
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
