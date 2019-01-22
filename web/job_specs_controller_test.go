package web_test

import (
	"bytes"
	"fmt"
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
	"github.com/stretchr/testify/require"
)

func BenchmarkJobSpecsController_Index(b *testing.B) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()
	setupJobSpecsControllerIndex(app)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		resp, cleanup := client.Get("/v2/specs")
		defer cleanup()
		assert.Equal(b, 200, resp.StatusCode, "Response should be successful")
	}
}

func TestJobSpecsController_Index_noSort(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	j1, err := setupJobSpecsControllerIndex(app)
	assert.NoError(t, err)

	resp, cleanup := client.Get("/v2/specs?size=x")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 422)

	resp, cleanup = client.Get("/v2/specs?size=1")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)
	body := cltest.ParseResponseBody(resp)

	metaCount, err := cltest.ParseJSONAPIResponseMetaCount(body)
	assert.NoError(t, err)
	assert.Equal(t, 2, metaCount)

	var links jsonapi.Links
	jobs := []models.JobSpec{}
	err = web.ParsePaginatedResponse(body, &jobs, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, jobs, 1)
	assert.Equal(t, j1.ID, jobs[0].ID)

	resp, cleanup = client.Get(links["next"].Href)
	defer cleanup()
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

func TestJobSpecsController_Index_sortCreatedAt(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	jobs := createJobs(app, 3)

	resp, cleanup := client.Get("/v2/specs?sort=createdAt&size=2")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)
	body := cltest.ParseResponseBody(resp)

	metaCount, err := cltest.ParseJSONAPIResponseMetaCount(body)
	assert.NoError(t, err)
	assert.Equal(t, 3, metaCount)

	var links jsonapi.Links
	ascJobs := []models.JobSpec{}
	err = web.ParsePaginatedResponse(body, &ascJobs, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, ascJobs, 2)
	assert.Equal(t, jobs[0].ID, ascJobs[0].ID)
	assert.Equal(t, jobs[1].ID, ascJobs[1].ID)

	resp, cleanup = client.Get("/v2/specs?sort=-createdAt&size=2")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)
	body = cltest.ParseResponseBody(resp)

	metaCount, err = cltest.ParseJSONAPIResponseMetaCount(body)
	assert.NoError(t, err)
	assert.Equal(t, 3, metaCount)

	descJobs := []models.JobSpec{}
	err = web.ParsePaginatedResponse(body, &descJobs, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, descJobs, 2)
	assert.Equal(t, jobs[2].ID, descJobs[0].ID)
	assert.Equal(t, jobs[1].ID, descJobs[1].ID)
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

func createJobs(app *cltest.TestApplication, n int) (jobs []*models.JobSpec) {
	for i := 0; i < n; i++ {
		j, _ := cltest.NewJobWithWebInitiator()
		err := app.Store.SaveJob(&j)
		if err != nil {
			panic(fmt.Sprintf("Could not save job: %v", err))
		}
		jobs = append(jobs, &j)
	}

	return jobs
}

func TestJobSpecsController_Create_HappyPath(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(cltest.LoadJSON("../internal/fixtures/web/hello_world_job.json")))
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	// Check Response
	var j models.JobSpec
	err := cltest.ParseJSONAPIResponse(resp, &j)
	require.NoError(t, err)

	adapter1, _ := adapters.For(j.Tasks[0], app.Store)
	httpGet := adapter1.BaseAdapter.(*adapters.HTTPGet)
	assert.Equal(t, httpGet.GetURL(), "https://bitstamp.net/api/ticker/")

	adapter2, _ := adapters.For(j.Tasks[1], app.Store)
	jsonParse := adapter2.BaseAdapter.(*adapters.JSONParse)
	assert.Equal(t, []string(jsonParse.Path), []string{"last"})

	adapter4, _ := adapters.For(j.Tasks[3], app.Store)
	signTx := adapter4.BaseAdapter.(*adapters.EthTx)
	assert.Equal(t, "0x356a04bCe728ba4c62A30294A55E6A8600a320B3", signTx.Address.String())
	assert.Equal(t, "0x609ff1bd", signTx.FunctionSelector.String())

	initr := j.Initiators[0]
	assert.Equal(t, models.InitiatorWeb, initr.Type)
	assert.NotEqual(t, models.Time{}, j.CreatedAt)

	// Check ORM
	orm := app.GetStore().ORM
	j, err = orm.FindJob(j.ID)
	assert.NoError(t, err)
	assert.Len(t, j.Initiators, 1)
	assert.Equal(t, models.InitiatorWeb, j.Initiators[0].Type)

	adapter1, _ = adapters.For(j.Tasks[0], app.Store)
	httpGet = adapter1.BaseAdapter.(*adapters.HTTPGet)
	assert.Equal(t, httpGet.GetURL(), "https://bitstamp.net/api/ticker/")
}

func TestJobSpecsController_Create_CaseInsensitiveTypes(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/caseinsensitive_hello_world_job.json")

	adapter1, _ := adapters.For(j.Tasks[0], app.Store)
	httpGet := adapter1.BaseAdapter.(*adapters.HTTPGet)
	assert.Equal(t, httpGet.GetURL(), "https://bitstamp.net/api/ticker/")

	adapter2, _ := adapters.For(j.Tasks[1], app.Store)
	jsonParse := adapter2.BaseAdapter.(*adapters.JSONParse)
	assert.Equal(t, []string(jsonParse.Path), []string{"last"})

	assert.Equal(t, "ethbytes32", j.Tasks[2].Type.String())

	adapter4, _ := adapters.For(j.Tasks[3], app.Store)
	signTx := adapter4.BaseAdapter.(*adapters.EthTx)
	assert.Equal(t, "0x356a04bCe728ba4c62A30294A55E6A8600a320B3", signTx.Address.String())
	assert.Equal(t, "0x609ff1bd", signTx.FunctionSelector.String())

	assert.Equal(t, models.InitiatorWeb, j.Initiators[0].Type)
	assert.Equal(t, models.InitiatorRunAt, j.Initiators[1].Type)
}

func TestJobSpecsController_Create_NonExistentTaskJob(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/nonexistent_task_job.json")
	resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(jsonStr))
	defer cleanup()

	assert.Equal(t, 400, resp.StatusCode, "Response should be caller error")

	expected := `{"errors":[{"detail":"idonotexist is not a supported adapter type"}]}`
	assert.Equal(t, expected, string(cltest.ParseResponseBody(resp)))
}

func TestJobSpecsController_Create_InvalidJob(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/run_at_wo_time_job.json")
	resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(jsonStr))
	defer cleanup()

	assert.Equal(t, 400, resp.StatusCode, "Response should be caller error")

	expected := `{"errors":[{"detail":"RunAt must have a time"}]}`
	assert.Equal(t, expected, string(cltest.ParseResponseBody(resp)))
}

func TestJobSpecsController_Create_InvalidCron(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/invalid_cron.json")
	resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(jsonStr))
	defer cleanup()

	assert.Equal(t, 400, resp.StatusCode, "Response should be caller error")

	expected := `{"errors":[{"detail":"Cron: Failed to parse int from !: strconv.Atoi: parsing \"!\": invalid syntax"}]}`
	assert.Equal(t, expected, string(cltest.ParseResponseBody(resp)))
}

func TestJobSpecsController_Create_Initiator_Only(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/initiator_only_job.json")
	resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(jsonStr))
	defer cleanup()

	assert.Equal(t, 400, resp.StatusCode, "Response should be caller error")

	expected := `{"errors":[{"detail":"Must have at least one Initiator and one Task"}]}`
	assert.Equal(t, expected, string(cltest.ParseResponseBody(resp)))
}

func TestJobSpecsController_Create_Task_Only(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/task_only_job.json")
	resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(jsonStr))
	defer cleanup()

	assert.Equal(t, 400, resp.StatusCode, "Response should be caller error")

	expected := `{"errors":[{"detail":"Must have at least one Initiator and one Task"}]}`
	assert.Equal(t, expected, string(cltest.ParseResponseBody(resp)))
}

func BenchmarkJobSpecsController_Show(b *testing.B) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()
	j := setupJobSpecsControllerShow(b, app)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		resp, _ := client.Get("/v2/specs/" + j.ID)
		assert.Equal(b, 200, resp.StatusCode, "Response should be successful")
	}
}

func TestJobSpecsController_Show(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	j := setupJobSpecsControllerShow(t, app)

	jr, err := app.Store.JobRunsFor(j.ID)
	assert.NoError(t, err)

	resp, cleanup := client.Get("/v2/specs/" + j.ID)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	var respJob presenters.JobSpec
	assert.NoError(t, cltest.ParseJSONAPIResponse(resp, &respJob))
	assert.Equal(t, j.Initiators[0].Schedule, respJob.Initiators[0].Schedule, "should have the same schedule")
	assert.Equal(t, jr[0].ID, respJob.Runs[0].ID, "should have job runs ordered by created at(descending)")
	assert.Equal(t, jr[1].ID, respJob.Runs[1].ID, "should have job runs ordered by created at(descending)")
}

func setupJobSpecsControllerShow(t assert.TestingT, app *cltest.TestApplication) *models.JobSpec {
	j, _ := cltest.NewJobWithSchedule("9 9 9 9 6")
	app.Store.SaveJob(&j)
	initr := j.Initiators[0]

	jr1 := j.NewRun(initr)
	jr1.ID = "2"
	assert.Nil(t, app.Store.SaveJobRun(&jr1))
	jr2 := j.NewRun(initr)
	jr2.ID = "1"
	jr2.CreatedAt = jr1.CreatedAt.Add(time.Second)
	assert.Nil(t, app.Store.SaveJobRun(&jr2))

	return &j
}

func TestJobSpecsController_Show_NotFound(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v2/specs/" + "garbage")
	defer cleanup()
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
}

func TestJobSpecsController_Show_Unauthenticated(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp, err := http.Get(app.Server.URL + "/v2/specs/" + "garbage")
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode, "Response should be forbidden")
}
