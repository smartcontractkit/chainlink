package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/freeport"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestPipelineRunsController_CreateWebhookJobRejected(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)
	// Bridge names are arbitrary; webhook creation is rejected before bridge validation.
	tomlStr := testspecs.GetWebhookSpecNoBody(uuid.New(), uuid.NewString(), uuid.NewString())
	body, err := json.Marshal(web.CreateJobRequest{TOML: tomlStr})
	require.NoError(t, err)
	response, cleanup := client.Post("/v2/jobs", bytes.NewReader(body))
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusUnprocessableEntity)
	require.Contains(t, string(cltest.ParseResponseBody(t, response)), "job type webhook has been removed")
}

func TestPipelineRunsController_RunExistingWebhookJobRejected(t *testing.T) {
	t.Parallel()

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(big.NewInt(0), nil)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.JobPipeline.HTTPRequest.DefaultTimeout = commonconfig.MustNewDuration(2 * time.Second)
	})

	app := cltest.NewApplicationWithConfig(t, cfg, ethClient)
	require.NoError(t, app.Start(testutils.Context(t)))

	jobUUID := uuid.New()
	cltest.MustInsertWebhookSpec(t, app.GetDB(), jobUUID)

	client := app.NewHTTPClient(nil)
	body := strings.NewReader(`{"data":{"result":"123.45"}}`)
	response, cleanup := client.Post("/v2/jobs/"+jobUUID.String()+"/runs", body)
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusUnprocessableEntity)
	require.Contains(t, string(cltest.ParseResponseBody(t, response)), "webhook")
}

func TestPipelineRunsController_Index_GlobalHappyPath(t *testing.T) {
	client, jobID, runIDs := setupPipelineRunsControllerTests(t)

	url := url.URL{Path: "/v2/pipeline/runs"}
	query := url.Query()
	query.Set("evmChainID", cltest.FixtureChainID.String())
	url.RawQuery = query.Encode()

	response, cleanup := client.Get(url.String())
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

	response, cleanup := client.Get("/v2/jobs/" + strconv.Itoa(int(jobID)) + "/runs")
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

	response, cleanup := client.Get("/v2/jobs/" + strconv.Itoa(int(jobID)) + "/runs?page=1&size=1")
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

	response, cleanup := client.Get("/v2/jobs/" + strconv.Itoa(int(jobID)) + "/runs/" + strconv.FormatInt(runIDs[0], 10))
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
	client := app.NewHTTPClient(nil)

	response, cleanup := client.Get("/v2/jobs/1/runs/invalid-run-ID")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusUnprocessableEntity)
}

func setupPipelineRunsControllerTests(t *testing.T) (cltest.HTTPClientCleaner, int32, []int64) {
	t.Parallel()
	ctx := testutils.Context(t)
	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	ethClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(nil, nil)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.OCR.Enabled = ptr(true)
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", freeport.GetOne(t))}
		c.P2P.PeerID = &cltest.DefaultP2PPeerID
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	app := cltest.NewApplicationWithConfigAndKey(t, cfg, ethClient, cltest.DefaultP2PKey)
	require.NoError(t, app.Start(ctx))
	require.NoError(t, app.KeyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	client := app.NewHTTPClient(nil)

	key, _ := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())

	nameAndExternalJobID := uuid.New()
	sp := fmt.Sprintf(`
	type               = "offchainreporting"
	schemaVersion      = 1
	externalJobID       = "%s"
	name               = "%s"
	contractAddress    = "%s"
	evmChainID		   = "%s"
	p2pv2Bootstrappers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
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
	`, nameAndExternalJobID, nameAndExternalJobID, testutils.NewAddress().Hex(), cltest.FixtureChainID.String(), cltest.DefaultOCRKeyBundleID, key.Address.Hex())
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
