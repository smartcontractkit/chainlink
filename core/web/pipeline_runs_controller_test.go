package web_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/pelletier/go-toml"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	"github.com/smartcontractkit/chainlink/core/web"
)

func TestPipelineRunsController_CreateWithBody_HappyPath(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	app.Config.Set("MAX_HTTP_ATTEMPTS", "1")
	app.Config.Set("DEFAULT_HTTP_TIMEOUT", "2s")
	app.Config.Set("TRIGGER_FALLBACK_DB_POLL_INTERVAL", "10ms")
	require.NoError(t, app.Start())

	// Setup the bridge
	{
		mockServer, cleanup := cltest.NewHTTPMockServerWithRequest(t, 200, `{}`, func(r *http.Request) {
			defer r.Body.Close()
			bs, err := ioutil.ReadAll(r.Body)
			require.NoError(t, err)
			require.Equal(t, `{"result":"12345"}`, string(bs))
		})
		defer cleanup()

		_, bridge := cltest.NewBridgeType(t, "my_bridge", mockServer.URL)
		require.NoError(t, app.Store.DB.Create(bridge).Error)
	}

	// Add the job
	var uuid uuid.UUID
	{
		tree, err := toml.LoadFile("../testdata/tomlspecs/webhook-job-spec-with-body.toml")
		require.NoError(t, err)
		webhookJobSpecFromFile, err := webhook.ValidatedWebhookSpec(tree.String(), app.GetExternalInitiatorManager())
		require.NoError(t, err)

		_, err = app.AddJobV2(context.Background(), webhookJobSpecFromFile, null.String{})
		require.NoError(t, err)

		uuid = webhookJobSpecFromFile.ExternalJobID
	}

	// Give the job.Spawner ample time to discover the job and start its service
	// (because Postgres events don't seem to work here)
	time.Sleep(3 * time.Second)

	// Make the request
	{
		client := app.NewHTTPClient()
		body := strings.NewReader(`{"data":{"result":"123.45"}}`)
		response, cleanup := client.Post("/v2/jobs/"+uuid.String()+"/runs", body)
		defer cleanup()
		cltest.AssertServerResponse(t, response, http.StatusOK)

		var parsedResponse pipeline.Run
		err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &parsedResponse)
		bs, _ := json.MarshalIndent(parsedResponse, "", "    ")
		fmt.Println(string(bs))
		assert.NoError(t, err)
		assert.NotNil(t, parsedResponse.ID)
		assert.NotNil(t, parsedResponse.CreatedAt)
		assert.NotNil(t, parsedResponse.FinishedAt)
		require.Len(t, parsedResponse.PipelineTaskRuns, 3)
	}
}

func TestPipelineRunsController_CreateNoBody_HappyPath(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	app.Config.Set("MAX_HTTP_ATTEMPTS", "1")
	app.Config.Set("DEFAULT_HTTP_TIMEOUT", "2s")
	app.Config.Set("TRIGGER_FALLBACK_DB_POLL_INTERVAL", "10ms")
	require.NoError(t, app.Start())

	// Setup the bridges
	{
		mockServer, cleanup := cltest.NewHTTPMockServer(t, 200, "POST", `{"data":{"result":"123.45"}}`)
		defer cleanup()

		_, bridge := cltest.NewBridgeType(t, "fetch_bridge", mockServer.URL)
		require.NoError(t, app.Store.DB.Create(bridge).Error)

		mockServer, cleanup = cltest.NewHTTPMockServerWithRequest(t, 200, `{}`, func(r *http.Request) {
			defer r.Body.Close()
			bs, err := ioutil.ReadAll(r.Body)
			require.NoError(t, err)
			require.Equal(t, `{"result":"12345"}`, string(bs))
		})
		defer cleanup()

		_, bridge = cltest.NewBridgeType(t, "submit_bridge", mockServer.URL)
		require.NoError(t, app.Store.DB.Create(bridge).Error)
	}

	// Add the job
	var uuid uuid.UUID
	{
		tree, err := toml.LoadFile("../testdata/tomlspecs/webhook-job-spec-no-body.toml")
		require.NoError(t, err)
		webhookJobSpecFromFile, err := webhook.ValidatedWebhookSpec(tree.String(), app.GetExternalInitiatorManager())
		require.NoError(t, err)

		_, err = app.AddJobV2(context.Background(), webhookJobSpecFromFile, null.String{})
		require.NoError(t, err)

		uuid = webhookJobSpecFromFile.ExternalJobID
	}

	// Give the job.Spawner ample time to discover the job and start its service
	// (because Postgres events don't seem to work here)
	time.Sleep(3 * time.Second)

	// Make the request
	{
		client := app.NewHTTPClient()
		response, cleanup := client.Post("/v2/jobs/"+uuid.String()+"/runs", nil)
		defer cleanup()
		cltest.AssertServerResponse(t, response, http.StatusOK)

		var parsedResponse pipeline.Run
		err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &parsedResponse)
		bs, _ := json.MarshalIndent(parsedResponse, "", "    ")
		fmt.Println(string(bs))
		assert.NoError(t, err)
		assert.NotNil(t, parsedResponse.ID)
		assert.NotNil(t, parsedResponse.CreatedAt)
		assert.NotNil(t, parsedResponse.FinishedAt)
		require.Len(t, parsedResponse.PipelineTaskRuns, 4)
	}
}

func TestPipelineRunsController_Index_HappyPath(t *testing.T) {
	client, jobID, runIDs, cleanup := setupPipelineRunsControllerTests(t)
	defer cleanup()

	response, cleanup := client.Get("/v2/jobs/" + fmt.Sprintf("%v", jobID) + "/runs")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	var parsedResponse []pipeline.Run
	responseBytes := cltest.ParseResponseBody(t, response)
	assert.Contains(t, string(responseBytes), `"errors":[null],"inputs":{"answer":"3","ds":"{\"USD\": 1}","ds_multiply":"3","ds_parse":1,"jobRun":{"meta":null}},"outputs":["3"]`)

	err := web.ParseJSONAPIResponse(responseBytes, &parsedResponse)
	assert.NoError(t, err)

	require.Len(t, parsedResponse, 2)
	assert.Equal(t, parsedResponse[1].ID, runIDs[0])
	assert.NotNil(t, parsedResponse[1].CreatedAt)
	assert.NotNil(t, parsedResponse[1].FinishedAt)
	// Successful pipeline runs does not save task runs.
	require.Len(t, parsedResponse[1].PipelineTaskRuns, 0)
}

func TestPipelineRunsController_Index_Pagination(t *testing.T) {
	client, jobID, runIDs, cleanup := setupPipelineRunsControllerTests(t)
	defer cleanup()

	response, cleanup := client.Get("/v2/jobs/" + fmt.Sprintf("%v", jobID) + "/runs?page=1&size=1")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	var parsedResponse []pipeline.Run
	responseBytes := cltest.ParseResponseBody(t, response)
	assert.Contains(t, string(responseBytes), `"errors":[null],"inputs":{"answer":"3","ds":"{\"USD\": 1}","ds_multiply":"3","ds_parse":1,"jobRun":{"meta":null}},"outputs":["3"]`)
	assert.Contains(t, string(responseBytes), `"meta":{"count":2}`)

	err := web.ParseJSONAPIResponse(responseBytes, &parsedResponse)
	assert.NoError(t, err)

	require.Len(t, parsedResponse, 1)
	assert.Equal(t, parsedResponse[0].ID, runIDs[1])
	assert.NotNil(t, parsedResponse[0].CreatedAt)
	assert.NotNil(t, parsedResponse[0].FinishedAt)
	require.Len(t, parsedResponse[0].PipelineTaskRuns, 0)
}

func TestPipelineRunsController_Show_HappyPath(t *testing.T) {
	client, jobID, runIDs, cleanup := setupPipelineRunsControllerTests(t)
	defer cleanup()

	response, cleanup := client.Get("/v2/jobs/" + fmt.Sprintf("%v", jobID) + "/runs/" + fmt.Sprintf("%v", runIDs[0]))
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	var parsedResponse pipeline.Run
	responseBytes := cltest.ParseResponseBody(t, response)
	assert.Contains(t, string(responseBytes), `"errors":[null],"inputs":{"answer":"3","ds":"{\"USD\": 1}","ds_multiply":"3","ds_parse":1,"jobRun":{"meta":null}},"outputs":["3"]`)

	err := web.ParseJSONAPIResponse(responseBytes, &parsedResponse)
	assert.NoError(t, err)

	assert.Equal(t, parsedResponse.ID, runIDs[0])
	assert.NotNil(t, parsedResponse.CreatedAt)
	assert.NotNil(t, parsedResponse.FinishedAt)
	require.Len(t, parsedResponse.PipelineTaskRuns, 0)
}

func TestPipelineRunsController_ShowRun_InvalidID(t *testing.T) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	response, cleanup := client.Get("/v2/jobs/1/runs/invalid-run-ID")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusUnprocessableEntity)
}

func setupPipelineRunsControllerTests(t *testing.T) (cltest.HTTPClientCleaner, int32, []int64, func()) {
	t.Parallel()
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()
	mockHTTP, cleanupHTTP := cltest.NewHTTPMockServer(t, http.StatusOK, "GET", `{"USD": 1}`)

	key := cltest.MustInsertRandomKey(t, app.Store.DB)

	sp := fmt.Sprintf(`
	type               = "offchainreporting"
	schemaVersion      = 1
	externalJobID       = "0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"
	contractAddress    = "%s"
	p2pPeerID          = "%s"
	p2pBootstrapPeers  = [
		"/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju",
	]
	keyBundleID        = "%s"
	transmitterAddress = "%s"
	observationSource = """
		// data source 1
		ds          [type=http method=GET url="%s"];
		ds_parse    [type=jsonparse path="USD"];
		ds_multiply [type=multiply times=3];

		ds -> ds_parse -> ds_multiply -> answer;

		answer [type=median index=0];
	"""
	`, cltest.NewAddress().Hex(), cltest.DefaultP2PPeerID, cltest.DefaultOCRKeyBundleID, key.Address.Hex(), mockHTTP.URL)
	var ocrJobSpec job.Job
	err := toml.Unmarshal([]byte(sp), &ocrJobSpec)
	require.NoError(t, err)
	var os job.OffchainReportingOracleSpec
	err = toml.Unmarshal([]byte(sp), &os)
	require.NoError(t, err)
	ocrJobSpec.OffchainreportingOracleSpec = &os

	err = app.GetKeyStore().OCR().Unlock(cltest.Password)
	require.NoError(t, err)

	jobID, err := app.AddJobV2(context.Background(), ocrJobSpec, null.String{})
	require.NoError(t, err)

	firstRunID, err := app.RunJobV2(context.Background(), jobID, nil)
	require.NoError(t, err)
	secondRunID, err := app.RunJobV2(context.Background(), jobID, nil)
	require.NoError(t, err)

	return client, jobID, []int64{firstRunID, secondRunID}, func() {
		cleanup()
		cleanupHTTP()
	}
}
