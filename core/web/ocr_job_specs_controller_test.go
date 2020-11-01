package web_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOCRJobSpecsController_Create_ValidationFailure(t *testing.T) {
	_, client, cleanup := setupOCRJobSpecsControllerTests(t)
	defer cleanup()

	fixtureBytes := cltest.MustReadFile(t, "testdata/oracle-spec-invalid-key.toml")

	resp, cleanup := client.Post("/v2/ocr/specs", bytes.NewReader(fixtureBytes))
	defer cleanup()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	b, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "{\"errors\":[{\"detail\":\"unrecognised key: isBootstrapNode\"}]}", string(b))
}

func TestOCRJobSpecsController_Create_HappyPath(t *testing.T) {
	app, client, cleanup := setupOCRJobSpecsControllerTests(t)
	defer cleanup()

	fixtureBytes := cltest.MustReadFile(t, "testdata/oracle-spec.toml")

	resp, cleanup := client.Post("/v2/ocr/specs", bytes.NewReader(fixtureBytes))
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	job := models.JobSpecV2{}
	require.NoError(t, app.Store.DB.Preload("OffchainreportingOracleSpec").First(&job).Error)

	b, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("{\"jobID\":%v}", job.ID), string(b))

	// Sanity check to make sure it inserted correctly
	require.Equal(t, models.EIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C"), job.OffchainreportingOracleSpec.ContractAddress)
}

func TestOCRJobSpecsController_Index_HappyPath(t *testing.T) {
	app, client, cleanup := setupOCRJobSpecsControllerTests(t)
	defer cleanup()

	var ocrJobSpec offchainreporting.OracleSpec
	toml.DecodeFile("testdata/oracle-spec.toml", &ocrJobSpec)
	app.AddJobV2(context.Background(), ocrJobSpec)

	response, cleanup := client.Get("/v2/ocr/specs")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	ocrJobSpecs := []models.JobSpecV2{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &ocrJobSpecs)
	assert.NoError(t, err)

	require.Len(t, ocrJobSpecs, 1)

	assert.Equal(t, ocrJobSpec.ContractAddress, ocrJobSpecs[0].OffchainreportingOracleSpec.ContractAddress)
	assert.Equal(t, ocrJobSpec.P2PPeerID, ocrJobSpecs[0].OffchainreportingOracleSpec.P2PPeerID)
	assert.Equal(t, ocrJobSpec.P2PBootstrapPeers, ocrJobSpecs[0].OffchainreportingOracleSpec.P2PBootstrapPeers)
	assert.Equal(t, ocrJobSpec.IsBootstrapPeer, ocrJobSpecs[0].OffchainreportingOracleSpec.IsBootstrapPeer)
	assert.Equal(t, ocrJobSpec.EncryptedOCRKeyBundleID, ocrJobSpecs[0].OffchainreportingOracleSpec.EncryptedOCRKeyBundleID)
	assert.Equal(t, ocrJobSpec.MonitoringEndpoint, ocrJobSpecs[0].OffchainreportingOracleSpec.MonitoringEndpoint)
	assert.Equal(t, ocrJobSpec.TransmitterAddress, ocrJobSpecs[0].OffchainreportingOracleSpec.TransmitterAddress)
	assert.Equal(t, ocrJobSpec.ObservationTimeout, ocrJobSpecs[0].OffchainreportingOracleSpec.ObservationTimeout)
	assert.Equal(t, ocrJobSpec.BlockchainTimeout, ocrJobSpecs[0].OffchainreportingOracleSpec.BlockchainTimeout)
	assert.Equal(t, ocrJobSpec.ContractConfigTrackerSubscribeInterval, ocrJobSpecs[0].OffchainreportingOracleSpec.ContractConfigTrackerSubscribeInterval)
	assert.Equal(t, ocrJobSpec.ContractConfigTrackerSubscribeInterval, ocrJobSpecs[0].OffchainreportingOracleSpec.ContractConfigTrackerSubscribeInterval)
	assert.Equal(t, ocrJobSpec.ContractConfigConfirmations, ocrJobSpecs[0].OffchainreportingOracleSpec.ContractConfigConfirmations)

	// Check that create and update dates are non empty values.
	// Empty date value is "0001-01-01 00:00:00 +0000 UTC" so we are checking for the
	// millenia and century characters to be present
	assert.Contains(t, ocrJobSpecs[0].OffchainreportingOracleSpec.CreatedAt.String(), "20")
	assert.Contains(t, ocrJobSpecs[0].OffchainreportingOracleSpec.UpdatedAt.String(), "20")
}

func setupOCRJobSpecsControllerTests(t *testing.T) (*cltest.TestApplication, cltest.HTTPClientCleaner, func()) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	return app, client, cleanup
}
