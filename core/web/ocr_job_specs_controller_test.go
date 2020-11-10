package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOCRJobSpecsController_Create_ValidationFailure(t *testing.T) {
	_, client, cleanup := setupOCRJobSpecsControllerTests(t)
	defer cleanup()

	body, _ := json.Marshal(models.CreateOCRJobSpecRequest{
		TOML: string(cltest.MustReadFile(t, "testdata/oracle-spec-invalid-key.toml")),
	})
	resp, cleanup := client.Post("/v2/ocr/specs", bytes.NewReader(body))
	defer cleanup()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	b, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "{\"errors\":[{\"detail\":\"unrecognised key: isBootstrapNode\"}]}", string(b))
}

func TestOCRJobSpecsController_Create_HappyPath(t *testing.T) {
	app, client, cleanup := setupOCRJobSpecsControllerTests(t)
	defer cleanup()

	body, _ := json.Marshal(models.CreateOCRJobSpecRequest{
		TOML: string(cltest.MustReadFile(t, "testdata/oracle-spec.toml")),
	})
	response, cleanup := client.Post("/v2/ocr/specs", bytes.NewReader(body))
	defer cleanup()
	require.Equal(t, http.StatusOK, response.StatusCode)

	job := models.JobSpecV2{}
	require.NoError(t, app.Store.DB.Preload("OffchainreportingOracleSpec").First(&job).Error)

	ocrJobSpec := models.JobSpecV2{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &ocrJobSpec)
	assert.NoError(t, err)

	assert.Equal(t, ocrJobSpec.OffchainreportingOracleSpec.ContractAddress, ocrJobSpec.OffchainreportingOracleSpec.ContractAddress)
	assert.Equal(t, ocrJobSpec.OffchainreportingOracleSpec.P2PPeerID, ocrJobSpec.OffchainreportingOracleSpec.P2PPeerID)
	assert.Equal(t, ocrJobSpec.OffchainreportingOracleSpec.P2PBootstrapPeers, ocrJobSpec.OffchainreportingOracleSpec.P2PBootstrapPeers)
	assert.Equal(t, ocrJobSpec.OffchainreportingOracleSpec.IsBootstrapPeer, ocrJobSpec.OffchainreportingOracleSpec.IsBootstrapPeer)
	assert.Equal(t, ocrJobSpec.OffchainreportingOracleSpec.EncryptedOCRKeyBundleID, ocrJobSpec.OffchainreportingOracleSpec.EncryptedOCRKeyBundleID)
	assert.Equal(t, ocrJobSpec.OffchainreportingOracleSpec.MonitoringEndpoint, ocrJobSpec.OffchainreportingOracleSpec.MonitoringEndpoint)
	assert.Equal(t, ocrJobSpec.OffchainreportingOracleSpec.TransmitterAddress, ocrJobSpec.OffchainreportingOracleSpec.TransmitterAddress)
	assert.Equal(t, ocrJobSpec.OffchainreportingOracleSpec.ObservationTimeout, ocrJobSpec.OffchainreportingOracleSpec.ObservationTimeout)
	assert.Equal(t, ocrJobSpec.OffchainreportingOracleSpec.BlockchainTimeout, ocrJobSpec.OffchainreportingOracleSpec.BlockchainTimeout)
	assert.Equal(t, ocrJobSpec.OffchainreportingOracleSpec.ContractConfigTrackerSubscribeInterval, ocrJobSpec.OffchainreportingOracleSpec.ContractConfigTrackerSubscribeInterval)
	assert.Equal(t, ocrJobSpec.OffchainreportingOracleSpec.ContractConfigTrackerSubscribeInterval, ocrJobSpec.OffchainreportingOracleSpec.ContractConfigTrackerSubscribeInterval)
	assert.Equal(t, ocrJobSpec.OffchainreportingOracleSpec.ContractConfigConfirmations, ocrJobSpec.OffchainreportingOracleSpec.ContractConfigConfirmations)

	// Sanity check to make sure it inserted correctly
	require.Equal(t, models.EIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C"), job.OffchainreportingOracleSpec.ContractAddress)
}

func TestOCRJobSpecsController_Index_HappyPath(t *testing.T) {
	client, cleanup, ocrJobSpecFromFile, _ := setupOCRJobSpecsWControllerTestsWithJob(t)
	defer cleanup()

	response, cleanup := client.Get("/v2/ocr/specs")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	ocrJobSpecs := []models.JobSpecV2{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &ocrJobSpecs)
	assert.NoError(t, err)

	require.Len(t, ocrJobSpecs, 1)

	runOCRJobSpecAssertions(t, ocrJobSpecFromFile, ocrJobSpecs[0])
}

func TestOCRJobSpecsController_Show_HappyPath(t *testing.T) {
	client, cleanup, ocrJobSpecFromFile, jobID := setupOCRJobSpecsWControllerTestsWithJob(t)
	defer cleanup()

	response, cleanup := client.Get("/v2/ocr/specs/" + fmt.Sprintf("%v", jobID))
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	ocrJobSpec := models.JobSpecV2{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &ocrJobSpec)
	assert.NoError(t, err)

	runOCRJobSpecAssertions(t, ocrJobSpecFromFile, ocrJobSpec)
}

func TestOCRJobSpecsController_Show_InvalidID(t *testing.T) {
	client, cleanup, _, _ := setupOCRJobSpecsWControllerTestsWithJob(t)
	defer cleanup()

	response, cleanup := client.Get("/v2/ocr/specs/uuidLikeString")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusUnprocessableEntity)
}

func TestOCRJobSpecsController_Show_NonExistentID(t *testing.T) {
	client, cleanup, _, _ := setupOCRJobSpecsWControllerTestsWithJob(t)
	defer cleanup()

	response, cleanup := client.Get("/v2/ocr/specs/999999999")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusNotFound)
}

func TestOCRJobSpecsController_Run_HappyPath(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	var ocrJobSpecFromFile offchainreporting.OracleSpec
	toml.DecodeFile("testdata/oracle-spec.toml", &ocrJobSpecFromFile)
	jobID, _ := app.AddJobV2(context.Background(), ocrJobSpecFromFile)

	response, cleanup := client.Post("/v2/ocr/specs/"+fmt.Sprintf("%v", jobID)+"/runs", nil)
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	parsedResponse := models.OCRJobRun{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &parsedResponse)
	assert.NoError(t, err)
	assert.NotNil(t, parsedResponse.ID)
}

func TestOCRJobSpecsController_Runs_HappyPath(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()
	mockHTTP, cleanupHTTP := cltest.NewHTTPMockServer(t, http.StatusOK, "GET", `{"USD": 1}`)
	defer cleanupHTTP()
	httpURL := mockHTTP.URL

	var ocrJobSpec offchainreporting.OracleSpec
	toml.Decode(fmt.Sprintf(`
	type               = "offchainreporting"
	schemaVersion      = 1
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
	`, cltest.NewAddress().Hex(), cltest.DefaultP2PPeerID, cltest.DefaultOCRKeyBundleID, cltest.DefaultKey, httpURL), &ocrJobSpec)
	jobID, err := app.AddJobV2(context.Background(), ocrJobSpec)
	require.NoError(t, err)
	runID, err := app.RunJobV2(context.Background(), jobID, nil)
	require.NoError(t, err)
	err = app.AwaitRun(context.Background(), runID)
	require.NoError(t, err)

	response, cleanup := client.Get("/v2/ocr/specs/" + fmt.Sprintf("%v", jobID) + "/runs")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	var parsedResponse []pipeline.Run
	responseBytes := cltest.ParseResponseBody(t, response)
	assert.Contains(t, string(responseBytes), `"meta":null,"errors":[null],"outputs":["3"]`)

	err = web.ParseJSONAPIResponse(responseBytes, &parsedResponse)
	assert.NoError(t, err)

	assert.Equal(t, parsedResponse[0].ID, runID)
	assert.NotNil(t, parsedResponse[0].CreatedAt)
	assert.NotNil(t, parsedResponse[0].FinishedAt)
	require.Len(t, parsedResponse[0].PipelineTaskRuns, 4)
}

func TestOCRJobSpecsController_ShowRun_HappyPath(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()
	mockHTTP, cleanupHTTP := cltest.NewHTTPMockServer(t, http.StatusOK, "GET", `{"USD": 1}`)
	defer cleanupHTTP()
	httpURL := mockHTTP.URL

	var ocrJobSpec offchainreporting.OracleSpec
	toml.Decode(fmt.Sprintf(`
	type               = "offchainreporting"
	schemaVersion      = 1
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
	`, cltest.NewAddress().Hex(), cltest.DefaultP2PPeerID, cltest.DefaultOCRKeyBundleID, cltest.DefaultKey, httpURL), &ocrJobSpec)
	jobID, err := app.AddJobV2(context.Background(), ocrJobSpec)
	require.NoError(t, err)
	runID, err := app.RunJobV2(context.Background(), jobID, nil)
	require.NoError(t, err)
	err = app.AwaitRun(context.Background(), runID)
	require.NoError(t, err)

	response, cleanup := client.Get("/v2/ocr/specs/" + fmt.Sprintf("%v", jobID) + "/runs/" + fmt.Sprintf("%v", runID))
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	var parsedResponse pipeline.Run
	responseBytes := cltest.ParseResponseBody(t, response)
	assert.Contains(t, string(responseBytes), `"meta":null,"errors":[null],"outputs":["3"]`)

	err = web.ParseJSONAPIResponse(responseBytes, &parsedResponse)
	assert.NoError(t, err)

	assert.Equal(t, parsedResponse.ID, runID)
	assert.NotNil(t, parsedResponse.CreatedAt)
	assert.NotNil(t, parsedResponse.FinishedAt)
	require.Len(t, parsedResponse.PipelineTaskRuns, 4)
}

func runOCRJobSpecAssertions(t *testing.T, ocrJobSpecFromFile offchainreporting.OracleSpec, ocrJobSpecFromServer models.JobSpecV2) {
	assert.Equal(t, ocrJobSpecFromFile.ContractAddress, ocrJobSpecFromServer.OffchainreportingOracleSpec.ContractAddress)
	assert.Equal(t, ocrJobSpecFromFile.P2PPeerID, ocrJobSpecFromServer.OffchainreportingOracleSpec.P2PPeerID)
	assert.Equal(t, ocrJobSpecFromFile.P2PBootstrapPeers, ocrJobSpecFromServer.OffchainreportingOracleSpec.P2PBootstrapPeers)
	assert.Equal(t, ocrJobSpecFromFile.IsBootstrapPeer, ocrJobSpecFromServer.OffchainreportingOracleSpec.IsBootstrapPeer)
	assert.Equal(t, ocrJobSpecFromFile.EncryptedOCRKeyBundleID, ocrJobSpecFromServer.OffchainreportingOracleSpec.EncryptedOCRKeyBundleID)
	assert.Equal(t, ocrJobSpecFromFile.MonitoringEndpoint, ocrJobSpecFromServer.OffchainreportingOracleSpec.MonitoringEndpoint)
	assert.Equal(t, ocrJobSpecFromFile.TransmitterAddress, ocrJobSpecFromServer.OffchainreportingOracleSpec.TransmitterAddress)
	assert.Equal(t, ocrJobSpecFromFile.ObservationTimeout, ocrJobSpecFromServer.OffchainreportingOracleSpec.ObservationTimeout)
	assert.Equal(t, ocrJobSpecFromFile.BlockchainTimeout, ocrJobSpecFromServer.OffchainreportingOracleSpec.BlockchainTimeout)
	assert.Equal(t, ocrJobSpecFromFile.ContractConfigTrackerSubscribeInterval, ocrJobSpecFromServer.OffchainreportingOracleSpec.ContractConfigTrackerSubscribeInterval)
	assert.Equal(t, ocrJobSpecFromFile.ContractConfigTrackerSubscribeInterval, ocrJobSpecFromServer.OffchainreportingOracleSpec.ContractConfigTrackerSubscribeInterval)
	assert.Equal(t, ocrJobSpecFromFile.ContractConfigConfirmations, ocrJobSpecFromServer.OffchainreportingOracleSpec.ContractConfigConfirmations)

	// Check that create and update dates are non empty values.
	// Empty date value is "0001-01-01 00:00:00 +0000 UTC" so we are checking for the
	// millenia and century characters to be present
	assert.Contains(t, ocrJobSpecFromServer.OffchainreportingOracleSpec.CreatedAt.String(), "20")
	assert.Contains(t, ocrJobSpecFromServer.OffchainreportingOracleSpec.UpdatedAt.String(), "20")
}

func setupOCRJobSpecsControllerTests(t *testing.T) (*cltest.TestApplication, cltest.HTTPClientCleaner, func()) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	return app, client, cleanup
}

func setupOCRJobSpecsWControllerTestsWithJob(t *testing.T) (cltest.HTTPClientCleaner, func(), offchainreporting.OracleSpec, int32) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	var ocrJobSpecFromFile offchainreporting.OracleSpec
	toml.DecodeFile("testdata/oracle-spec.toml", &ocrJobSpecFromFile)
	jobID, _ := app.AddJobV2(context.Background(), ocrJobSpecFromFile)
	return client, cleanup, ocrJobSpecFromFile, jobID
}
