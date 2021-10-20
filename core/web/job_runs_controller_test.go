package web_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/static"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkJobRunsController_Index(b *testing.B) {
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(b)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(b,
		ethClient,
	)
	defer cleanup()
	app.Start()
	run1, _, _ := setupJobRunsControllerIndex(b, app)
	client := app.NewHTTPClient()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		resp, cleanup := client.Get("/v2/runs?jobSpecId=" + run1.JobSpecID.String())
		defer cleanup()
		assert.Equal(b, http.StatusOK, resp.StatusCode, "Response should be successful")
	}
}

func TestJobRunsController_Index(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	app.Start()
	client := app.NewHTTPClient()

	runA, runB, runC := setupJobRunsControllerIndex(t, app)

	resp, cleanup := client.Get("/v2/runs?size=x&jobSpecId=" + runA.JobSpecID.String())
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusUnprocessableEntity)

	resp, cleanup = client.Get("/v2/runs?size=1&jobSpecId=" + runA.JobSpecID.String())
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	var links jsonapi.Links
	var meta jsonapi.Meta
	var runs []models.JobRun

	err := web.ParsePaginatedResponseWithMeta(cltest.ParseResponseBody(t, resp), &runs, &links, &meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	require.Len(t, runs, 1)
	assert.Equal(t, runA.ID, runs[0].ID, "expected runs order by createdAt ascending")
	assert.EqualValues(t, meta["errored"], 1, "expect there to be 1 errored run")
	assert.EqualValues(t, meta["completed"], 1, "expect there to be 1 completed run")
	assert.EqualValues(t, meta["count"], 2, "expect there to be 2 runs in total")

	resp, cleanup = client.Get(links["next"].Href)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	var nextPageLinks jsonapi.Links
	var nextPageRuns = []models.JobRun{}

	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &nextPageRuns, &nextPageLinks)
	assert.NoError(t, err)
	assert.Empty(t, nextPageLinks["next"])
	assert.NotEmpty(t, nextPageLinks["prev"])

	assert.Len(t, nextPageRuns, 1)
	assert.Equal(t, runB.ID, nextPageRuns[0].ID, "expected runs order by createdAt ascending")

	resp, cleanup = client.Get("/v2/runs?sort=-createdAt")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	var allJobRunLinks jsonapi.Links
	var allJobRuns []models.JobRun

	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &allJobRuns, &allJobRunLinks)
	assert.NoError(t, err)
	assert.Empty(t, allJobRunLinks["next"].Href)
	assert.Empty(t, allJobRunLinks["prev"].Href)

	assert.Len(t, allJobRuns, 3)
	assert.Equal(t, runC.ID, allJobRuns[0].ID, "expected runs ordered by created at descending")
	assert.Equal(t, runB.ID, allJobRuns[1].ID, "expected runs ordered by created at descending")
	assert.Equal(t, runA.ID, allJobRuns[2].ID, "expected runs ordered by created at descending")
}

func setupJobRunsControllerIndex(t assert.TestingT, app *cltest.TestApplication) (*models.JobRun, *models.JobRun, *models.JobRun) {
	j1 := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.CreateJob(&j1))

	j2 := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.CreateJob(&j2))

	now := time.Now()

	runA := cltest.NewJobRun(j1)
	runA.ID = uuid.NewV4()
	runA.CreatedAt = now.Add(-2 * time.Second)
	runA.Status = models.RunStatusErrored
	assert.Nil(t, app.Store.CreateJobRun(&runA))

	runB := cltest.NewJobRun(j1)
	runB.ID = uuid.NewV4()
	runB.CreatedAt = now.Add(-time.Second)
	runB.Status = models.RunStatusCompleted
	assert.Nil(t, app.Store.CreateJobRun(&runB))

	runC := cltest.NewJobRun(j2)
	runC.ID = uuid.NewV4()
	runC.CreatedAt = now
	runC.Status = models.RunStatusCompleted
	assert.Nil(t, app.Store.CreateJobRun(&runC))

	return &runA, &runB, &runC
}

func TestJobRunsController_Create_Success(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	app.Start()

	j := cltest.NewJobWithWebInitiator()
	assert.NoError(t, app.Store.CreateJob(&j))

	jr := cltest.CreateJobRunViaWeb(t, app, j, `{"result":"100"}`)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)
	value := cltest.MustResultString(t, jr.Result)
	assert.Equal(t, "100", value)
}

func TestJobRunsController_Create_Wrong_ExternalInitiator(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	app.Start()

	eir_url := cltest.WebURL(t, "http://localhost:8888")

	eir := &models.ExternalInitiatorRequest{
		Name: "bitcoin",
		URL:  &eir_url,
	}
	eia := auth.NewToken()
	ei, err := models.NewExternalInitiator(eia, eir)
	require.NoError(t, err)
	assert.NoError(t, app.Store.CreateExternalInitiator(ei))

	j := cltest.NewJobWithExternalInitiator(ei)
	assert.NoError(t, app.Store.CreateJob(&j))

	wrongEIR := &models.ExternalInitiatorRequest{
		Name: "someCoin",
		URL:  &eir_url,
	}
	wrongEIA := auth.NewToken()
	wrongEI, err := models.NewExternalInitiator(wrongEIA, wrongEIR)
	require.NoError(t, err)
	assert.NoError(t, app.Store.CreateExternalInitiator(wrongEI))

	// Set up AUTH
	headers := make(map[string]string)
	headers[static.ExternalInitiatorAccessKeyHeader] = wrongEIA.AccessKey
	headers[static.ExternalInitiatorSecretHeader] = wrongEIA.Secret

	url := app.Config.ClientNodeURL() + "/v2/specs/" + j.ID.String() + "/runs"
	bodyBuf := bytes.NewBufferString(`{"result":"100"}`)
	resp, cleanup := cltest.UnauthenticatedPost(t, url, bodyBuf, headers)
	defer cleanup()
	assert.Equal(t, http.StatusForbidden, resp.StatusCode, "Response should be forbidden")
}

func TestJobRunsController_Create_ExternalInitiator_Success(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	app.Start()

	url := cltest.WebURL(t, "http://localhost:8888")
	eia := auth.NewToken()
	eir := &models.ExternalInitiatorRequest{
		Name: "bitcoin",
		URL:  &url,
	}
	ei, err := models.NewExternalInitiator(eia, eir)
	require.NoError(t, err)
	assert.NoError(t, app.Store.CreateExternalInitiator(ei))

	j := cltest.NewJobWithExternalInitiator(ei)
	assert.NoError(t, app.Store.CreateJob(&j))

	jr := cltest.CreateJobRunViaExternalInitiator(
		t, app,
		j, *eia, `{"result":"100"}`,
	)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)
	value := cltest.MustResultString(t, jr.Result)
	assert.Equal(t, "100", value)
}

func TestJobRunsController_Create_Archived(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	app.Start()

	j := cltest.NewJobWithWebInitiator()
	require.NoError(t, app.Store.CreateJob(&j))
	require.NoError(t, app.Store.ArchiveJob(j.ID))

	client := app.NewHTTPClient()
	resp, cleanup := client.Post("/v2/specs/"+j.ID.String()+"/runs", bytes.NewBufferString(`{"result":"100"}`))
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusGone)
}

func TestJobRunsController_Create_EmptyBody(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	app.Start()

	j := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.CreateJob(&j))

	jr := cltest.CreateJobRunViaWeb(t, app, j)
	cltest.WaitForJobRunToComplete(t, app.Store, jr)
}

func TestJobRunsController_Create_InvalidBody(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	app.Start()
	client := app.NewHTTPClient()

	j := cltest.NewJobWithWebInitiator()
	assert.Nil(t, app.Store.CreateJob(&j))

	resp, cleanup := client.Post("/v2/specs/"+j.ID.String()+"/runs", bytes.NewBufferString(`{`))
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusInternalServerError)
}

func TestJobRunsController_Create_WithoutWebInitiator(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	app.Start()
	client := app.NewHTTPClient()

	j := cltest.NewJob()
	assert.Nil(t, app.Store.CreateJob(&j))

	resp, cleanup := client.Post("/v2/specs/"+j.ID.String()+"/runs", bytes.NewBuffer([]byte{}))
	defer cleanup()
	assert.Equal(t, http.StatusForbidden, resp.StatusCode, "Response should be forbidden")
}

func TestJobRunsController_Create_NotFound(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	app.Start()
	client := app.NewHTTPClient()

	resp, cleanup := client.Post("/v2/specs/4C95A8FA-EEAC-4BD5-97D9-27806D200D3C/runs", bytes.NewBuffer([]byte{}))
	defer cleanup()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Response should be not found")
}

func TestJobRunsController_Create_InvalidID(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	app.Start()
	client := app.NewHTTPClient()

	resp, cleanup := client.Post("/v2/specs/garbageID/runs", bytes.NewBuffer([]byte{}))
	defer cleanup()
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode, "Response should be unprocessable entity")
}

func TestJobRunsController_Update_Success(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	app.Start()

	tests := []struct {
		name     string
		archived bool
	}{
		{"normal_job", false},
		{"archived_job", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bta, bt := cltest.NewBridgeType(t, test.name)
			require.NoError(t, app.Store.CreateBridgeType(bt))
			j := cltest.NewJobWithWebInitiator()
			j.Tasks = []models.TaskSpec{{Type: bt.Name}}
			require.NoError(t, app.Store.CreateJob(&j))
			jr := cltest.NewJobRunPendingBridge(j)
			require.NoError(t, app.Store.CreateJobRun(&jr))

			if test.archived {
				require.NoError(t, app.Store.ArchiveJob(j.ID))
			}

			// resume run
			body := fmt.Sprintf(`{"id":"%v","data":{"result": "100"}}`, jr.ID.String())
			headers := map[string]string{"Authorization": "Bearer " + bta.IncomingToken}
			url := app.Config.ClientNodeURL() + "/v2/runs/" + jr.ID.String()
			resp, cleanup := cltest.UnauthenticatedPatch(t, url, bytes.NewBufferString(body), headers)
			defer cleanup()

			require.Equal(t, http.StatusOK, resp.StatusCode, "Response should be successful")
			var respJobRun presenters.JobRun
			assert.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &respJobRun))
			require.Equal(t, jr.ID, respJobRun.ID)

			jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)
			value := cltest.MustResultString(t, jr.Result)
			assert.Equal(t, "100", value)
		})
	}
}

func TestJobRunsController_Update_WrongAccessToken(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	app.Start()
	client := app.NewHTTPClient()

	_, bt := cltest.NewBridgeType(t)
	assert.Nil(t, app.Store.CreateBridgeType(bt))
	j := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.Nil(t, app.Store.CreateJob(&j))
	jr := cltest.NewJobRunPendingBridge(j)
	assert.Nil(t, app.Store.CreateJobRun(&jr))

	body := fmt.Sprintf(`{"id":"%v","data":{"result": "100"}}`, jr.ID.String())
	headers := map[string]string{"Authorization": "Bearer " + "wrongaccesstoken"}
	resp, cleanup := client.Patch("/v2/runs/"+jr.ID.String(), bytes.NewBufferString(body), headers)
	defer cleanup()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Response should be unauthorized")
	jr, err := app.Store.FindJobRun(jr.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.RunStatusPendingBridge, jr.GetStatus())
}

func TestJobRunsController_Update_NotPending(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	app.Start()
	client := app.NewHTTPClient()

	bta, bt := cltest.NewBridgeType(t)
	assert.Nil(t, app.Store.CreateBridgeType(bt))
	j := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.Nil(t, app.Store.CreateJob(&j))
	jr := cltest.NewJobRun(j)
	assert.Nil(t, app.Store.CreateJobRun(&jr))

	body := fmt.Sprintf(`{"id":"%v","data":{"result": "100"}}`, jr.ID.String())
	headers := map[string]string{"Authorization": "Bearer " + bta.IncomingToken}
	resp, cleanup := client.Patch("/v2/runs/"+jr.ID.String(), bytes.NewBufferString(body), headers)
	defer cleanup()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "Response should be unsuccessful")
}

func TestJobRunsController_Update_WithError(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	app.Start()
	client := app.NewHTTPClient()

	bta, bt := cltest.NewBridgeType(t)
	assert.Nil(t, app.Store.CreateBridgeType(bt))
	j := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.Nil(t, app.Store.CreateJob(&j))
	jr := cltest.NewJobRunPendingBridge(j)
	assert.Nil(t, app.Store.CreateJobRun(&jr))

	body := fmt.Sprintf(`{"id":"%v","error":"stack overflow","data":{"result": "0"}}`, jr.ID.String())
	headers := map[string]string{"Authorization": "Bearer " + bta.IncomingToken}
	resp, cleanup := client.Patch("/v2/runs/"+jr.ID.String(), bytes.NewBufferString(body), headers)
	defer cleanup()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response should be successful")
	var respJobRun presenters.JobRun
	assert.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &respJobRun))
	assert.Equal(t, jr.ID, respJobRun.ID)

	jr = cltest.WaitForJobRunStatus(t, app.Store, jr, models.RunStatusErrored)
	value := cltest.MustResultString(t, jr.Result)
	assert.Equal(t, "0", value)
}

func TestJobRunsController_Update_BadInput(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	app.Start()
	client := app.NewHTTPClient()

	_, bt := cltest.NewBridgeType(t)
	assert.Nil(t, app.Store.CreateBridgeType(bt))
	j := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.Nil(t, app.Store.CreateJob(&j))
	jr := cltest.NewJobRunPendingBridge(j)
	assert.Nil(t, app.Store.CreateJobRun(&jr))

	body := fmt.Sprint(`{`, jr.ID.String())
	resp, cleanup := client.Patch("/v2/runs/"+jr.ID.String(), bytes.NewBufferString(body))
	defer cleanup()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Response should be successful")
	jr, err := app.Store.FindJobRun(jr.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.RunStatusPendingBridge, jr.GetStatus())
}

func TestJobRunsController_Update_NotFound(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	app.Start()
	client := app.NewHTTPClient()

	_, bt := cltest.NewBridgeType(t)
	assert.Nil(t, app.Store.CreateBridgeType(bt))
	j := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.Nil(t, app.Store.CreateJob(&j))
	jr := cltest.NewJobRunPendingBridge(j)
	assert.Nil(t, app.Store.CreateJobRun(&jr))

	body := fmt.Sprintf(`{"id":"%v","data":{"result": "100"}}`, jr.ID.String())
	resp, cleanup := client.Patch("/v2/runs/4C95A8FA-EEAC-4BD5-97D9-27806D200D3C", bytes.NewBufferString(body))
	defer cleanup()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Response should be not found")
	jr, err := app.Store.FindJobRun(jr.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.RunStatusPendingBridge, jr.GetStatus())
}

func TestJobRunsController_Show_Found(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	app.Start()
	client := app.NewHTTPClient()

	j := cltest.NewJobWithSchedule("CRON_TZ=UTC 9 9 9 9 6")
	app.Store.CreateJob(&j)

	jr := cltest.NewJobRun(j)
	assert.NoError(t, app.Store.CreateJobRun(&jr))

	resp, cleanup := client.Get("/v2/runs/" + jr.ID.String())
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode, "Response should be successful")

	var respJobRun presenters.JobRun
	assert.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &respJobRun))
	assert.Equal(t, jr.Initiator.Schedule, respJobRun.Initiator.Schedule, "should have the same schedule")
	assert.Equal(t, jr.ID, respJobRun.ID, "should have job run id")
}

func TestJobRunsController_Show_NotFound(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	app.Start()
	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v2/runs/4C95A8FA-EEAC-4BD5-97D9-27806D200D3C")
	defer cleanup()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Response should be not found")
}

func TestJobRunsController_Show_InvalidID(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	app.Start()
	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v2/runs/garbage")
	defer cleanup()
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode, "Response should be unprocessable entity")
}

func TestJobRunsController_Show_Unauthenticated(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	app.Start()

	resp, err := http.Get(app.Server.URL + "/v2/runs/notauthorized")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Response should be forbidden")
}

func TestJobRunsController_Cancel(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	app.Start()

	client := app.NewHTTPClient()

	t.Run("invalid run id", func(t *testing.T) {
		response, cleanup := client.Put("/v2/runs/xxx/cancellation", nil)
		defer cleanup()
		cltest.AssertServerResponse(t, response, http.StatusUnprocessableEntity)
	})

	t.Run("missing run", func(t *testing.T) {
		resp, cleanup := client.Put("/v2/runs/29023583-0D39-4844-9696-451102590936/cancellation", nil)
		defer cleanup()
		cltest.AssertServerResponse(t, resp, http.StatusNotFound)
	})

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, app.Store.CreateJob(&job))
	run := cltest.NewJobRun(job)
	require.NoError(t, app.Store.CreateJobRun(&run))

	t.Run("valid run", func(t *testing.T) {
		resp, cleanup := client.Put(fmt.Sprintf("/v2/runs/%s/cancellation", run.ID), nil)
		defer cleanup()
		cltest.AssertServerResponse(t, resp, http.StatusOK)

		r, err := app.Store.FindJobRun(run.ID)
		assert.NoError(t, err)
		assert.Equal(t, models.RunStatusCancelled, r.GetStatus())
	})
}
