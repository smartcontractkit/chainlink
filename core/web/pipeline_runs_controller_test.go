package web_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pelletier/go-toml"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/webhook"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestPipelineRunsController_CreateWithBody_HappyPath(t *testing.T) {
	t.Parallel()

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.JobPipeline.HTTPRequest.DefaultTimeout = models.MustNewDuration(2 * time.Second)
		c.Database.Listener.FallbackPollInterval = models.MustNewDuration(10 * time.Millisecond)
	})

	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Setup the bridge
	mockServer := cltest.NewHTTPMockServerWithRequest(t, 200, `{}`, func(r *http.Request) {
		defer r.Body.Close()
		bs, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, `{"result":"12345"}`, string(bs))
	})

	_, bridge := cltest.MustCreateBridge(t, app.GetSqlxDB(), cltest.BridgeOpts{URL: mockServer.URL}, app.GetConfig())

	// Add the job
	var uuid uuid.UUID
	{
		tomlStr := fmt.Sprintf(testspecs.WebhookSpecWithBody, bridge.Name.String())
		jb, err := webhook.ValidatedWebhookSpec(tomlStr, app.GetExternalInitiatorManager())
		require.NoError(t, err)

		err = app.AddJobV2(testutils.Context(t), &jb)
		require.NoError(t, err)

		uuid = jb.ExternalJobID
	}

	// Give the job.Spawner ample time to discover the job and start its service
	// (because Postgres events don't seem to work here)
	time.Sleep(3 * time.Second)

	// Make the request
	{
		client := app.NewHTTPClient(cltest.APIEmailAdmin)
		body := strings.NewReader(`{"data":{"result":"123.45"}}`)
		response, cleanup := client.Post("/v2/jobs/"+uuid.String()+"/runs", body)
		defer cleanup()
		cltest.AssertServerResponse(t, response, http.StatusOK)

		var parsedResponse presenters.PipelineRunResource
		err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &parsedResponse)
		assert.NoError(t, err)
		assert.NotNil(t, parsedResponse.ID)
		assert.NotNil(t, parsedResponse.CreatedAt)
		assert.NotNil(t, parsedResponse.FinishedAt)
		require.Len(t, parsedResponse.TaskRuns, 3)
	}
}

func TestPipelineRunsController_CreateNoBody_HappyPath(t *testing.T) {
	t.Parallel()

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.JobPipeline.HTTPRequest.DefaultTimeout = models.MustNewDuration(2 * time.Second)
		c.Database.Listener.FallbackPollInterval = models.MustNewDuration(10 * time.Millisecond)
	})

	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Setup the bridges
	mockServer := cltest.NewHTTPMockServer(t, 200, "POST", `{"data":{"result":"123.45"}}`)

	_, bridge := cltest.MustCreateBridge(t, app.GetSqlxDB(), cltest.BridgeOpts{URL: mockServer.URL}, app.GetConfig())

	mockServer = cltest.NewHTTPMockServerWithRequest(t, 200, `{}`, func(r *http.Request) {
		defer r.Body.Close()
		bs, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, `{"result":"12345"}`, string(bs))
	})

	_, submitBridge := cltest.MustCreateBridge(t, app.GetSqlxDB(), cltest.BridgeOpts{URL: mockServer.URL}, app.GetConfig())

	// Add the job
	var uuid uuid.UUID
	{
		tomlStr := fmt.Sprintf(testspecs.WebhookSpecNoBody, bridge.Name.String(), submitBridge.Name.String())
		jb, err := webhook.ValidatedWebhookSpec(tomlStr, app.GetExternalInitiatorManager())
		require.NoError(t, err)

		err = app.AddJobV2(testutils.Context(t), &jb)
		require.NoError(t, err)

		uuid = jb.ExternalJobID
	}

	// Give the job.Spawner ample time to discover the job and start its service
	// (because Postgres events don't seem to work here)
	time.Sleep(3 * time.Second)

	// Make the request (authorized as user)
	{
		client := app.NewHTTPClient(cltest.APIEmailAdmin)
		response, cleanup := client.Post("/v2/jobs/"+uuid.String()+"/runs", nil)
		defer cleanup()
		cltest.AssertServerResponse(t, response, http.StatusOK)

		var parsedResponse presenters.PipelineRunResource
		err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &parsedResponse)
		bs, _ := json.MarshalIndent(parsedResponse, "", "    ")
		t.Log(string(bs))
		assert.NoError(t, err)
		assert.NotNil(t, parsedResponse.ID)
		assert.NotNil(t, parsedResponse.CreatedAt)
		assert.NotNil(t, parsedResponse.FinishedAt)
		require.Len(t, parsedResponse.TaskRuns, 4)
	}
}

func TestPipelineRunsController_Index_GlobalHappyPath(t *testing.T) {
	client, jobID, runIDs := setupPipelineRunsControllerTests(t)

	response, cleanup := client.Get("/v2/pipeline/runs")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	var parsedResponse []presenters.PipelineRunResource
	responseBytes := cltest.ParseResponseBody(t, response)
	assert.Contains(t, string(responseBytes), `"outputs":["3"],"errors":[null],"allErrors":["uh oh"],"fatalErrors":[null],"inputs":{"answer":"3","ds1":"{\"USD\": 1}","ds1_multiply":"3","ds1_parse":1,"ds2":"{\"USD\": 1}","ds2_multiply":"3","ds2_parse":1,"ds3":{},"jobRun":{"meta":null}`)

	err := web.ParseJSONAPIResponse(responseBytes, &parsedResponse)
	assert.NoError(t, err)

	require.Len(t, parsedResponse, 2)
	// Job Run ID is returned in descending order by run ID (most recent run first)
	assert.Equal(t, parsedResponse[1].ID, strconv.Itoa(int(runIDs[0])))
	assert.NotNil(t, parsedResponse[1].CreatedAt)
	assert.NotNil(t, parsedResponse[1].FinishedAt)
	assert.Equal(t, jobID, parsedResponse[1].PipelineSpec.JobID)
	require.Len(t, parsedResponse[1].TaskRuns, 8)
}

func TestPipelineRunsController_Index_HappyPath(t *testing.T) {
	client, jobID, runIDs := setupPipelineRunsControllerTests(t)

	response, cleanup := client.Get("/v2/jobs/" + fmt.Sprintf("%v", jobID) + "/runs")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	var parsedResponse []presenters.PipelineRunResource
	responseBytes := cltest.ParseResponseBody(t, response)
	assert.Contains(t, string(responseBytes), `"outputs":["3"],"errors":[null],"allErrors":["uh oh"],"fatalErrors":[null],"inputs":{"answer":"3","ds1":"{\"USD\": 1}","ds1_multiply":"3","ds1_parse":1,"ds2":"{\"USD\": 1}","ds2_multiply":"3","ds2_parse":1,"ds3":{},"jobRun":{"meta":null}`)

	err := web.ParseJSONAPIResponse(responseBytes, &parsedResponse)
	assert.NoError(t, err)

	require.Len(t, parsedResponse, 2)
	// Job Run ID is returned in descending order by run ID (most recent run first)
	assert.Equal(t, parsedResponse[1].ID, strconv.Itoa(int(runIDs[0])))
	assert.NotNil(t, parsedResponse[1].CreatedAt)
	assert.NotNil(t, parsedResponse[1].FinishedAt)
	assert.Equal(t, jobID, parsedResponse[1].PipelineSpec.JobID)
	require.Len(t, parsedResponse[1].TaskRuns, 8)
}

func TestPipelineRunsController_Index_Pagination(t *testing.T) {
	client, jobID, runIDs := setupPipelineRunsControllerTests(t)

	response, cleanup := client.Get("/v2/jobs/" + fmt.Sprintf("%v", jobID) + "/runs?page=1&size=1")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	var parsedResponse []presenters.PipelineRunResource
	responseBytes := cltest.ParseResponseBody(t, response)
	assert.Contains(t, string(responseBytes), `"outputs":["3"],"errors":[null],"allErrors":["uh oh"],"fatalErrors":[null],"inputs":{"answer":"3","ds1":"{\"USD\": 1}","ds1_multiply":"3","ds1_parse":1,"ds2":"{\"USD\": 1}","ds2_multiply":"3","ds2_parse":1,"ds3":{},"jobRun":{"meta":null}`)
	assert.Contains(t, string(responseBytes), `"meta":{"count":2}`)

	err := web.ParseJSONAPIResponse(responseBytes, &parsedResponse)
	assert.NoError(t, err)

	require.Len(t, parsedResponse, 1)
	assert.Equal(t, parsedResponse[0].ID, strconv.Itoa(int(runIDs[1])))
	assert.NotNil(t, parsedResponse[0].CreatedAt)
	assert.NotNil(t, parsedResponse[0].FinishedAt)
	require.Len(t, parsedResponse[0].TaskRuns, 8)
}

func TestPipelineRunsController_Show_HappyPath(t *testing.T) {
	client, jobID, runIDs := setupPipelineRunsControllerTests(t)

	response, cleanup := client.Get("/v2/jobs/" + fmt.Sprintf("%v", jobID) + "/runs/" + fmt.Sprintf("%v", runIDs[0]))
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	var parsedResponse presenters.PipelineRunResource
	responseBytes := cltest.ParseResponseBody(t, response)
	assert.Contains(t, string(responseBytes), `"outputs":["3"],"errors":[null],"allErrors":["uh oh"],"fatalErrors":[null],"inputs":{"answer":"3","ds1":"{\"USD\": 1}","ds1_multiply":"3","ds1_parse":1,"ds2":"{\"USD\": 1}","ds2_multiply":"3","ds2_parse":1,"ds3":{},"jobRun":{"meta":null}`)
	err := web.ParseJSONAPIResponse(responseBytes, &parsedResponse)
	require.NoError(t, err)

	assert.Equal(t, parsedResponse.ID, strconv.Itoa(int(runIDs[0])))
	assert.NotNil(t, parsedResponse.CreatedAt)
	assert.NotNil(t, parsedResponse.FinishedAt)
	require.Len(t, parsedResponse.TaskRuns, 8)
}

func TestPipelineRunsController_ShowRun_InvalidID(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	response, cleanup := client.Get("/v2/jobs/1/runs/invalid-run-ID")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusUnprocessableEntity)
}

func setupPipelineRunsControllerTests(t *testing.T) (cltest.HTTPClientCleaner, int32, []int64) {
	t.Parallel()
	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.OCR.Enabled = ptr(true)
		c.P2P.V1.Enabled = ptr(true)
		c.P2P.PeerID = &cltest.DefaultP2PPeerID
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfigAndKey(t, cfg, ethClient, cltest.DefaultP2PKey)
	require.NoError(t, app.Start(testutils.Context(t)))
	require.NoError(t, app.KeyStore.OCR().Add(cltest.DefaultOCRKey))
	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	key, _ := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())

	sp := fmt.Sprintf(`
	type               = "offchainreporting"
	schemaVersion      = 1
	externalJobID       = "0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"
	contractAddress    = "%s"
	p2pBootstrapPeers  = [
		"/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju",
	]
	p2pv2Bootstrappers = []
	keyBundleID        = "%s"
	transmitterAddress = "%s"
	observationSource = """
		// data source 1
		ds1          [type=memo value=<"{\\"USD\\": 1}">];
		ds1_parse    [type=jsonparse path="USD"];
		ds1_multiply [type=multiply times=3];

		ds2          [type=memo value=<"{\\"USD\\": 1}">];
		ds2_parse    [type=jsonparse path="USD"];
		ds2_multiply [type=multiply times=3];

		ds3          [type=fail msg="uh oh"];

		ds1 -> ds1_parse -> ds1_multiply -> answer;
		ds2 -> ds2_parse -> ds2_multiply -> answer;
		ds3 -> answer;

		answer [type=median index=0];
	"""
	`, testutils.NewAddress().Hex(), cltest.DefaultOCRKeyBundleID, key.Address.Hex())
	var jb job.Job
	err := toml.Unmarshal([]byte(sp), &jb)
	require.NoError(t, err)
	var os job.OCROracleSpec
	err = toml.Unmarshal([]byte(sp), &os)
	require.NoError(t, err)
	jb.OCROracleSpec = &os

	err = app.AddJobV2(testutils.Context(t), &jb)
	require.NoError(t, err)

	firstRunID, err := app.RunJobV2(testutils.Context(t), jb.ID, nil)
	require.NoError(t, err)
	secondRunID, err := app.RunJobV2(testutils.Context(t), jb.ID, nil)
	require.NoError(t, err)

	return client, jb.ID, []int64{firstRunID, secondRunID}
}
