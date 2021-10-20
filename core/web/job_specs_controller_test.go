package web_test

import (
	"bytes"
	"encoding/json"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkJobSpecsController_Index(b *testing.B) {
	app, cleanup := cltest.NewApplication(b)
	defer cleanup()
	client := app.NewHTTPClient()
	setupJobSpecsControllerIndex(app)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		resp, cleanup := client.Get("/v2/specs")
		defer cleanup()
		assert.Equal(b, http.StatusOK, resp.StatusCode, "Response should be successful")
	}
}

func TestJobSpecsController_Index_noSort(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	j1, err := setupJobSpecsControllerIndex(app)
	assert.NoError(t, err)

	resp, cleanup := client.Get("/v2/specs?size=x")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusUnprocessableEntity)

	resp, cleanup = client.Get("/v2/specs?size=1")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)
	body := cltest.ParseResponseBody(t, resp)

	metaCount, err := cltest.ParseJSONAPIResponseMetaCount(body)
	require.NoError(t, err)
	require.Equal(t, 2, metaCount)

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
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	jobs = []models.JobSpec{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &jobs, &links)
	assert.NoError(t, err)
	assert.Empty(t, links["next"])
	assert.NotEmpty(t, links["prev"])

	require.Len(t, jobs, 1)
	assert.Equal(t, models.InitiatorWeb, jobs[0].Initiators[0].Type, "should have the same type")
	assert.NotEqual(t, true, jobs[0].Initiators[0].Ran, "should ignore fields for other initiators")
}

func TestJobSpecsController_Index_sortCreatedAt(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	j2 := cltest.NewJobWithWebInitiator()
	j2.CreatedAt = time.Now().AddDate(0, 0, 1)
	require.NoError(t, app.Store.CreateJob(&j2))

	j3 := cltest.NewJobWithWebInitiator()
	j3.CreatedAt = time.Now().AddDate(0, 0, 2)
	require.NoError(t, app.Store.CreateJob(&j3))

	j1 := cltest.NewJobWithWebInitiator() // deliberately out of order
	j1.CreatedAt = time.Now().AddDate(0, 0, -1)
	require.NoError(t, app.Store.CreateJob(&j1))

	jobs := []models.JobSpec{j1, j2, j3}

	resp, cleanup := client.Get("/v2/specs?sort=createdAt&size=2")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)
	body := cltest.ParseResponseBody(t, resp)

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
	cltest.AssertServerResponse(t, resp, http.StatusOK)
	body = cltest.ParseResponseBody(t, resp)

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
	j1 := cltest.NewJobWithSchedule("CRON_TZ=UTC 9 9 9 9 6")
	j1.CreatedAt = time.Now().AddDate(0, 0, -1)
	err := app.Store.CreateJob(&j1)
	if err != nil {
		return nil, err
	}
	j2 := cltest.NewJobWithWebInitiator()
	j2.Initiators[0].Ran = true
	err = app.Store.CreateJob(&j2)
	return &j1, err
}

func TestJobSpecsController_Create_HappyPath(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(cltest.MustReadFile(t, "../testdata/jsonspecs/hello_world_job.json")))
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	// Check Response
	var j models.JobSpec
	err := cltest.ParseJSONAPIResponse(t, resp, &j)
	require.NoError(t, err)

	adapter1, _ := adapters.For(j.Tasks[0], app.Store.Config, app.Store.ORM, nil)
	httpGet := adapter1.BaseAdapter.(*adapters.HTTPGet)
	assert.Equal(t, httpGet.GetURL(), "https://bitstamp.net/api/ticker/")

	adapter2, _ := adapters.For(j.Tasks[1], app.Store.Config, app.Store.ORM, nil)
	jsonParse := adapter2.BaseAdapter.(*adapters.JSONParse)
	assert.Equal(t, []string(jsonParse.Path), []string{"last"})

	adapter4, _ := adapters.For(j.Tasks[3], app.Store.Config, app.Store.ORM, nil)
	signTx := adapter4.BaseAdapter.(*adapters.EthTx)
	assert.Equal(t, "0x356a04bCe728ba4c62A30294A55E6A8600a320B3", signTx.ToAddress.String())
	assert.Equal(t, "0x609ff1bd", signTx.FunctionSelector.String())

	initr := j.Initiators[0]
	assert.Equal(t, models.InitiatorWeb, initr.Type)
	assert.NotEqual(t, models.AnyTime{}, j.CreatedAt)

	// Check ORM
	orm := app.GetStore().ORM
	j, err = orm.FindJobSpec(j.ID)
	require.NoError(t, err)
	require.Len(t, j.Initiators, 1)
	assert.Equal(t, models.InitiatorWeb, j.Initiators[0].Type)

	adapter1, _ = adapters.For(j.Tasks[0], app.Store.Config, app.Store.ORM, nil)
	httpGet = adapter1.BaseAdapter.(*adapters.HTTPGet)
	assert.Equal(t, httpGet.GetURL(), "https://bitstamp.net/api/ticker/")
}

func TestJobSpecsController_Create_CustomName(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	fixtureBytes := cltest.MustReadFile(t, "../testdata/jsonspecs/hello_world_job.json")
	jsr := cltest.JSONFromBytes(t, fixtureBytes)
	jsr, err := jsr.MultiAdd(map[string]interface{}{"name": "CustomJobName"})
	require.NoError(t, err)
	requestBody, err := json.Marshal(jsr)
	require.NoError(t, err)

	t.Run("it creates the job spec with the specified custom name", func(t *testing.T) {
		resp, cleanup := client.Post("/v2/specs", bytes.NewReader(requestBody))
		defer cleanup()
		cltest.AssertServerResponse(t, resp, http.StatusOK)

		var j models.JobSpec
		err = cltest.ParseJSONAPIResponse(t, resp, &j)
		require.NoError(t, err)

		orm := app.GetStore().ORM
		j, err = orm.FindJobSpec(j.ID)
		require.NoError(t, err)
		assert.Equal(t, j.Name, "CustomJobName")
	})
}

func TestJobSpecsController_CreateExternalInitiator_Success(t *testing.T) {
	t.Parallel()

	var eiReceived webhook.JobSpecNotice
	eiMockServer, assertCalled := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", "",
		func(header http.Header, body string) {
			err := json.Unmarshal([]byte(body), &eiReceived)
			require.NoError(t, err)
		},
	)
	defer assertCalled()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient, cltest.UseRealExternalInitiatorManager)
	defer cleanup()
	app.Start()

	url := cltest.WebURL(t, eiMockServer.URL)
	eir := models.ExternalInitiatorRequest{
		Name: "someCoin",
		URL:  &url,
	}
	eia := auth.NewToken()
	ei, err := models.NewExternalInitiator(eia, &eir)
	require.NoError(t, err)
	err = app.GetStore().CreateExternalInitiator(ei)
	require.NoError(t, err)

	jobSpec := cltest.FixtureCreateJobViaWeb(t, app, "./../testdata/jsonspecs/external_initiator_job.json")
	expected := webhook.JobSpecNotice{
		JobID:  jobSpec.ID,
		Type:   models.InitiatorExternal,
		Params: cltest.JSONFromString(t, `{"foo":"bar"}`),
	}
	assert.Equal(t, expected, eiReceived)

	jobRun := cltest.CreateJobRunViaExternalInitiator(t, app, jobSpec, *eia, "")
	cltest.WaitForJobRunToComplete(t, app.Store, jobRun)
}

func TestJobSpecsController_Create_CaseInsensitiveTypes(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	require.NoError(t, app.Start())

	j := cltest.FixtureCreateJobViaWeb(t, app, "../testdata/jsonspecs/caseinsensitive_hello_world_job.json")

	adapter1, _ := adapters.For(j.Tasks[0], app.Store.Config, app.Store.ORM, nil)
	httpGet := adapter1.BaseAdapter.(*adapters.HTTPGet)
	assert.Equal(t, httpGet.GetURL(), "https://bitstamp.net/api/ticker/")

	adapter2, _ := adapters.For(j.Tasks[1], app.Store.Config, app.Store.ORM, nil)
	jsonParse := adapter2.BaseAdapter.(*adapters.JSONParse)
	assert.Equal(t, []string(jsonParse.Path), []string{"last"})

	assert.Equal(t, "ethbytes32", j.Tasks[2].Type.String())

	adapter4, _ := adapters.For(j.Tasks[3], app.Store.Config, app.Store.ORM, nil)
	signTx := adapter4.BaseAdapter.(*adapters.EthTx)
	assert.Equal(t, "0x356a04bCe728ba4c62A30294A55E6A8600a320B3", signTx.ToAddress.String())
	assert.Equal(t, "0x609ff1bd", signTx.FunctionSelector.String())

	assert.Equal(t, models.InitiatorWeb, j.Initiators[0].Type)
	assert.Equal(t, models.InitiatorRunAt, j.Initiators[1].Type)
}

func TestJobSpecsController_Create_NonExistentTaskJob(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	jsonStr := cltest.MustReadFile(t, "../testdata/jsonspecs/nonexistent_task_job.json")
	resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(jsonStr))
	defer cleanup()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Response should be caller error")

	expected := `{"errors":[{"detail":"idonotexist is not a supported adapter type"}]}`
	body := string(cltest.ParseResponseBody(t, resp))
	assert.Equal(t, expected, strings.TrimSpace(body))
}

func TestJobSpecsController_Create_FluxMonitor_disabled(t *testing.T) {
	config := cltest.NewTestConfig(t)
	config.Set("CHAINLINK_DEV", "FALSE")
	config.Set("FEATURE_FLUX_MONITOR", "FALSE")
	config.Set("GAS_ESTIMATOR_MODE", "FixedPrice")
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()

	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	jsonStr := cltest.MustReadFile(t, "../testdata/jsonspecs/flux_monitor_job.json")
	resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(jsonStr))
	defer cleanup()

	require.Equal(t, http.StatusText(http.StatusNotImplemented), http.StatusText(resp.StatusCode))
	expected := `{"errors":[{"detail":"The Flux Monitor feature is disabled by configuration"}]}`
	body := string(cltest.ParseResponseBody(t, resp))
	assert.Equal(t, expected, strings.TrimSpace(body))
}

func TestJobSpecsController_Create_FluxMonitor_enabled(t *testing.T) {
	config := cltest.NewTestConfig(t)
	config.Set("CHAINLINK_DEV", "FALSE")
	config.Set("FEATURE_FLUX_MONITOR", "TRUE")
	config.Set("GAS_ESTIMATOR_MODE", "FixedPrice")
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()

	getOraclesResult, err := cltest.GenericEncode([]string{"address[]"}, []common.Address{})
	require.NoError(t, err)
	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "getOracles").
		Return(getOraclesResult, nil)
	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "latestRoundData").
		Return(nil, errors.New("first round"))
	result := cltest.MakeRoundStateReturnData(2, true, 10000, 7, 0, 1000, 100, 1)
	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "oracleRoundState").
		Return(result, nil).Maybe()
	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "minSubmissionValue").
		Return(cltest.MustGenericEncode([]string{"uint256"}, big.NewInt(0)), nil)
	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "maxSubmissionValue").
		Return(cltest.MustGenericEncode([]string{"uint256"}, big.NewInt(10000000)), nil)

	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	jsonStr := cltest.MustReadFile(t, "../testdata/jsonspecs/flux_monitor_job.json")
	resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(jsonStr))
	defer cleanup()

	cltest.AssertServerResponse(t, resp, http.StatusOK)
}

func TestJobSpecsController_Create_FluxMonitor_Bridge(t *testing.T) {
	config := cltest.NewTestConfig(t)
	config.Set("CHAINLINK_DEV", "FALSE")
	config.Set("FEATURE_FLUX_MONITOR", "TRUE")
	config.Set("GAS_ESTIMATOR_MODE", "FixedPrice")
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithConfig(t, config,
		ethClient,
	)
	defer cleanup()

	getOraclesResult, err := cltest.GenericEncode([]string{"address[]"}, []common.Address{})
	require.NoError(t, err)
	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "getOracles").
		Return(getOraclesResult, nil)
	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "latestRoundData").
		Return(nil, errors.New("first round"))
	result := cltest.MakeRoundStateReturnData(2, true, 10000, 7, 0, 1000, 100, 1)
	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "oracleRoundState").
		Return(result, nil).Maybe()
	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "minSubmissionValue").
		Return(cltest.MustGenericEncode([]string{"uint256"}, big.NewInt(0)), nil)
	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "maxSubmissionValue").
		Return(cltest.MustGenericEncode([]string{"uint256"}, big.NewInt(10000000)), nil)

	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	bridge := &models.BridgeType{
		Name: models.MustNewTaskType("testbridge"),
		URL:  cltest.WebURL(t, "https://testing.com/bridges"),
	}
	require.NoError(t, app.Store.CreateBridgeType(bridge))

	jsonStr := cltest.MustReadFile(t, "../testdata/jsonspecs/flux_monitor_bridge_job_noop.json")
	resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(jsonStr))
	defer cleanup()

	require.Equal(t, http.StatusText(http.StatusOK), http.StatusText(resp.StatusCode))
}

func TestJobSpecsController_Create_FluxMonitor_NoBridgeError(t *testing.T) {
	config := cltest.NewTestConfig(t)
	config.Set("CHAINLINK_DEV", "FALSE")
	config.Set("FEATURE_FLUX_MONITOR", "TRUE")
	config.Set("GAS_ESTIMATOR_MODE", "FixedPrice")
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithConfig(t, config,
		ethClient,
	)
	defer cleanup()

	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	jsonStr := cltest.MustReadFile(t, "../testdata/jsonspecs/flux_monitor_bridge_job_noop.json")
	resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(jsonStr))
	defer cleanup()

	require.Equal(t, http.StatusText(http.StatusBadRequest), http.StatusText(resp.StatusCode))
}

func TestJobSpecsController_Create_InvalidJob(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	jsonStr := cltest.MustReadFile(t, "../testdata/jsonspecs/run_at_wo_time_job.json")
	resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(jsonStr))
	defer cleanup()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Response should be caller error")

	expected := `{"errors":[{"detail":"RunAt must have a time"}]}`
	body := string(cltest.ParseResponseBody(t, resp))
	assert.Equal(t, expected, strings.TrimSpace(body))
}

func TestJobSpecsController_Create_InvalidCron(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	jsonStr := cltest.MustReadFile(t, "../testdata/jsonspecs/invalid_cron.json")
	resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(jsonStr))
	defer cleanup()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Response should be caller error")

	expected := `{"errors":[{"detail":"Cron: failed to parse int from !: strconv.Atoi: parsing \"!\": invalid syntax"}]}`
	body := string(cltest.ParseResponseBody(t, resp))
	assert.Equal(t, expected, strings.TrimSpace(body))
}

func TestJobSpecsController_Create_Initiator_Only(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	jsonStr := cltest.MustReadFile(t, "../testdata/jsonspecs/initiator_only_job.json")
	resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(jsonStr))
	defer cleanup()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Response should be caller error")

	expected := `{"errors":[{"detail":"Must have at least one Initiator and one Task"}]}`
	body := string(cltest.ParseResponseBody(t, resp))
	assert.Equal(t, expected, strings.TrimSpace(body))
}

func TestJobSpecsController_Create_Task_Only(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	jsonStr := cltest.MustReadFile(t, "../testdata/jsonspecs/task_only_job.json")
	resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(jsonStr))
	defer cleanup()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Response should be caller error")

	expected := `{"errors":[{"detail":"Must have at least one Initiator and one Task"}]}`
	body := string(cltest.ParseResponseBody(t, resp))
	assert.Equal(t, expected, strings.TrimSpace(body))
}

func TestJobSpecsController_Create_EthDisabled(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	config.Set("ETH_DISABLED", true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config, ethClient)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())
	db := app.Store.DB

	client := app.NewHTTPClient()

	t.Run("VRF", func(t *testing.T) {
		jsonStr := cltest.MustReadFile(t, "./../testdata/jsonspecs/randomness_job.json")
		resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(jsonStr))
		t.Cleanup(cleanup)

		assert.Equal(t, 200, resp.StatusCode)
		cltest.AssertCount(t, db, models.JobSpec{}, 1)
	})

	t.Run("runlog", func(t *testing.T) {
		jsonStr := cltest.MustReadFile(t, "../testdata/jsonspecs/runlog_noop_job.json")
		resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(jsonStr))
		t.Cleanup(cleanup)

		assert.Equal(t, 200, resp.StatusCode)
		cltest.AssertCount(t, db, models.JobSpec{}, 2)
	})

	t.Run("ethlog", func(t *testing.T) {
		jsonStr := cltest.MustReadFile(t, "../testdata/jsonspecs/runlog_noop_job.json")
		resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(jsonStr))
		t.Cleanup(cleanup)

		assert.Equal(t, 200, resp.StatusCode)
		cltest.AssertCount(t, db, models.JobSpec{}, 3)
	})
}

func BenchmarkJobSpecsController_Show(b *testing.B) {
	app, cleanup := cltest.NewApplication(b)
	defer cleanup()
	require.NoError(b, app.Start())

	client := app.NewHTTPClient()
	j := setupJobSpecsControllerShow(b, app)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		resp, _ := client.Get("/v2/specs/" + j.ID.String())
		assert.Equal(b, http.StatusOK, resp.StatusCode, "Response should be successful")
	}
}

func TestJobSpecsController_Show(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	j := setupJobSpecsControllerShow(t, app)

	resp, cleanup := client.Get("/v2/specs/" + j.ID.String())
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	var respJob presenters.JobSpec
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &respJob))
	require.Len(t, j.Initiators, 1)
	require.Len(t, respJob.Initiators, 1)
	require.Len(t, respJob.Errors, 1)
	assert.Equal(t, j.Initiators[0].Schedule, respJob.Initiators[0].Schedule, "should have the same schedule")
}

func TestJobSpecsController_Show_FluxMonitorJob(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	j := cltest.NewJobWithFluxMonitorInitiator()
	app.Store.CreateJob(&j)

	resp, cleanup := client.Get("/v2/specs/" + j.ID.String())
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	var respJob presenters.JobSpec
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &respJob))
	require.Equal(t, len(respJob.Initiators), len(j.Initiators))
	require.Equal(t, respJob.Initiators[0].Address, j.Initiators[0].Address)
	require.Equal(t, respJob.Initiators[0].RequestData, j.Initiators[0].RequestData)
	require.Equal(t, respJob.Initiators[0].Feeds, j.Initiators[0].Feeds)
	require.Equal(t, respJob.Initiators[0].Threshold, j.Initiators[0].Threshold)
	require.Equal(t, respJob.Initiators[0].AbsoluteThreshold, j.Initiators[0].AbsoluteThreshold)
	require.Equal(t, respJob.Initiators[0].IdleTimer, j.Initiators[0].IdleTimer)
	require.Equal(t, respJob.Initiators[0].PollTimer, j.Initiators[0].PollTimer)
	require.Equal(t, respJob.Initiators[0].Precision, j.Initiators[0].Precision)
}

func TestJobSpecsController_Show_MultipleTasks(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	// Create a task with multiple jobs
	j := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{
		{Type: models.MustNewTaskType("Task1")},
		{Type: models.MustNewTaskType("Task2")},
		{Type: models.MustNewTaskType("Task3")},
		{Type: models.MustNewTaskType("Task4")},
	}
	assert.NoError(t, app.Store.CreateJob(&j))

	resp, cleanup := client.Get("/v2/specs/" + j.ID.String())
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	var respJob presenters.JobSpec
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &respJob))
	assert.Equal(t, string(respJob.Tasks[0].Type), "task1")
	assert.Equal(t, string(respJob.Tasks[1].Type), "task2")
	assert.Equal(t, string(respJob.Tasks[2].Type), "task3")
	assert.Equal(t, string(respJob.Tasks[3].Type), "task4")
}

func setupJobSpecsControllerShow(t assert.TestingT, app *cltest.TestApplication) *models.JobSpec {
	j := cltest.NewJobWithSchedule("CRON_TZ=UTC 9 9 9 9 6")
	app.Store.CreateJob(&j)

	app.Store.UpsertErrorFor(j.ID, "job spec error description")

	jr1 := cltest.NewJobRun(j)
	assert.Nil(t, app.Store.CreateJobRun(&jr1))
	jr2 := cltest.NewJobRun(j)
	jr2.CreatedAt = jr1.CreatedAt.Add(time.Second)
	assert.Nil(t, app.Store.CreateJobRun(&jr2))

	return &j
}

func TestJobSpecsController_Show_NotFound(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v2/specs/190AE4CE-40B6-4D60-A3DA-061C5ACD32D0")
	defer cleanup()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Response should be not found")
}

func TestJobSpecsController_Show_InvalidUuid(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v2/specs/garbage")
	defer cleanup()
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode, "Response should be unprocessable entity")
}

func TestJobSpecsController_Show_Unauthenticated(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	resp, err := http.Get(app.Server.URL + "/v2/specs/garbage")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Response should be forbidden")
}

func TestJobSpecsController_Destroy(t *testing.T) {
	t.Parallel()
	ethClient, s, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(s, nil)

	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	job := cltest.NewJobWithLogInitiator()
	require.NoError(t, app.Store.CreateJob(&job))

	resp, cleanup := client.Delete("/v2/specs/" + job.ID.String())
	defer cleanup()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Error(t, utils.JustError(app.Store.FindJobSpec(job.ID)))
	assert.Equal(t, 0, len(app.ChainlinkApplication.JobSubscriber.Jobs()))
}

func TestJobSpecsController_DestroyAdd(t *testing.T) {
	ethClient, s, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(s, nil)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	job := cltest.NewJobWithLogInitiator()
	job.Name = "testjob"
	require.NoError(t, app.Store.CreateJob(&job))

	resp, cleanup := client.Delete("/v2/specs/" + job.ID.String())
	defer cleanup()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Error(t, utils.JustError(app.Store.FindJobSpec(job.ID)))
	assert.Equal(t, 0, len(app.ChainlinkApplication.JobSubscriber.Jobs()))

	job = cltest.NewJobWithLogInitiator()
	job.Name = "testjob"
	require.NoError(t, app.Store.CreateJob(&job))

	// Can delete this new job
	resp, cleanup = client.Delete("/v2/specs/" + job.ID.String())
	defer cleanup()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Error(t, utils.JustError(app.Store.FindJobSpec(job.ID)))
	assert.Equal(t, 0, len(app.ChainlinkApplication.JobSubscriber.Jobs()))
}

func TestJobSpecsController_Destroy_MultipleJobs(t *testing.T) {
	t.Parallel()
	ethClient, s, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(s, nil)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	job1 := cltest.NewJobWithLogInitiator()
	job2 := cltest.NewJobWithLogInitiator()
	require.NoError(t, app.Store.CreateJob(&job1))
	require.NoError(t, app.Store.CreateJob(&job2))

	resp, cleanup := client.Delete("/v2/specs/" + job1.ID.String())
	defer cleanup()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Error(t, utils.JustError(app.Store.FindJobSpec(job1.ID)))
	assert.Equal(t, 0, len(app.ChainlinkApplication.JobSubscriber.Jobs()))

	resp, cleanup = client.Delete("/v2/specs/" + job2.ID.String())
	defer cleanup()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Error(t, utils.JustError(app.Store.FindJobSpec(job2.ID)))
	assert.Equal(t, 0, len(app.ChainlinkApplication.JobSubscriber.Jobs()))
}
