package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/pelletier/go-toml"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestJobsController_Create_ValidationFailure(t *testing.T) {
	var (
		contractAddress    = cltest.NewEIP55Address()
		monitoringEndpoint = "chain.link:101"
	)

	var tt = []struct {
		name        string
		pid         models.PeerID
		kb          models.Sha256Hash
		ta          models.EIP55Address
		expectedErr error
	}{
		{
			name:        "invalid keybundle",
			pid:         models.PeerID(cltest.DefaultP2PPeerID),
			kb:          models.Sha256Hash(cltest.Random32Byte()),
			ta:          cltest.DefaultKeyAddressEIP55,
			expectedErr: job.ErrNoSuchKeyBundle,
		},
		{
			name:        "invalid peerID",
			pid:         models.PeerID(cltest.NonExistentP2PPeerID),
			kb:          cltest.DefaultOCRKeyBundleIDSha256,
			ta:          cltest.DefaultKeyAddressEIP55,
			expectedErr: job.ErrNoSuchPeerID,
		},
		{
			name:        "invalid transmitter address",
			pid:         models.PeerID(cltest.DefaultP2PPeerID),
			kb:          cltest.DefaultOCRKeyBundleIDSha256,
			ta:          cltest.NewEIP55Address(),
			expectedErr: job.ErrNoSuchTransmitterAddress,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, client, cleanup := setupJobsControllerTests(t)
			defer cleanup()
			sp := cltest.MinimalOCRNonBootstrapSpec(contractAddress, tc.ta, tc.pid, monitoringEndpoint, tc.kb)
			body, _ := json.Marshal(models.CreateJobSpecRequest{
				TOML: sp,
			})
			resp, cleanup := client.Post("/v2/jobs", bytes.NewReader(body))
			defer cleanup()
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			b, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(b), tc.expectedErr.Error())
		})
	}
}

func TestJobsController_Create_HappyPath_OffchainReportingSpec(t *testing.T) {
	app, client, cleanup := setupJobsControllerTests(t)
	defer cleanup()

	body, _ := json.Marshal(models.CreateJobSpecRequest{
		TOML: string(cltest.MustReadFile(t, "testdata/oracle-spec.toml")),
	})
	response, cleanup := client.Post("/v2/jobs", bytes.NewReader(body))
	defer cleanup()
	require.Equal(t, http.StatusOK, response.StatusCode)

	job := models.JobSpecV2{}
	require.NoError(t, app.Store.DB.Preload("OffchainreportingOracleSpec").First(&job).Error)

	ocrJobSpec := models.JobSpecV2{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &ocrJobSpec)
	assert.NoError(t, err)

	assert.Equal(t, "web oracle spec", job.Name.ValueOrZero())
	assert.Equal(t, job.OffchainreportingOracleSpec.P2PPeerID, ocrJobSpec.OffchainreportingOracleSpec.P2PPeerID)
	assert.Equal(t, job.OffchainreportingOracleSpec.P2PBootstrapPeers, ocrJobSpec.OffchainreportingOracleSpec.P2PBootstrapPeers)
	assert.Equal(t, job.OffchainreportingOracleSpec.IsBootstrapPeer, ocrJobSpec.OffchainreportingOracleSpec.IsBootstrapPeer)
	assert.Equal(t, job.OffchainreportingOracleSpec.EncryptedOCRKeyBundleID, ocrJobSpec.OffchainreportingOracleSpec.EncryptedOCRKeyBundleID)
	assert.Equal(t, job.OffchainreportingOracleSpec.MonitoringEndpoint, ocrJobSpec.OffchainreportingOracleSpec.MonitoringEndpoint)
	assert.Equal(t, job.OffchainreportingOracleSpec.TransmitterAddress, ocrJobSpec.OffchainreportingOracleSpec.TransmitterAddress)
	assert.Equal(t, job.OffchainreportingOracleSpec.ObservationTimeout, ocrJobSpec.OffchainreportingOracleSpec.ObservationTimeout)
	assert.Equal(t, job.OffchainreportingOracleSpec.BlockchainTimeout, ocrJobSpec.OffchainreportingOracleSpec.BlockchainTimeout)
	assert.Equal(t, job.OffchainreportingOracleSpec.ContractConfigTrackerSubscribeInterval, ocrJobSpec.OffchainreportingOracleSpec.ContractConfigTrackerSubscribeInterval)
	assert.Equal(t, job.OffchainreportingOracleSpec.ContractConfigTrackerSubscribeInterval, ocrJobSpec.OffchainreportingOracleSpec.ContractConfigTrackerSubscribeInterval)
	assert.Equal(t, job.OffchainreportingOracleSpec.ContractConfigConfirmations, ocrJobSpec.OffchainreportingOracleSpec.ContractConfigConfirmations)
	assert.NotNil(t, ocrJobSpec.PipelineSpec.DotDagSource)

	// Sanity check to make sure it inserted correctly
	require.Equal(t, models.EIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C"), job.OffchainreportingOracleSpec.ContractAddress)
}

func TestJobsController_Create_HappyPath_EthRequestEventSpec(t *testing.T) {
	app, client, cleanup := setupJobsControllerTests(t)
	defer cleanup()

	body, _ := json.Marshal(models.CreateJobSpecRequest{
		TOML: string(cltest.MustReadFile(t, "testdata/eth-request-event-spec.toml")),
	})
	response, cleanup := client.Post("/v2/jobs", bytes.NewReader(body))
	defer cleanup()
	require.Equal(t, http.StatusOK, response.StatusCode)

	job := models.JobSpecV2{}
	require.NoError(t, app.Store.DB.Preload("EthRequestEventSpec").First(&job).Error)

	jobSpec := models.JobSpecV2{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &jobSpec)
	assert.NoError(t, err)

	assert.Equal(t, "example eth request event spec", job.Name.ValueOrZero())
	assert.NotNil(t, jobSpec.PipelineSpec.DotDagSource)

	// Sanity check to make sure it inserted correctly
	require.Equal(t, models.EIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C"), job.EthRequestEventSpec.ContractAddress)
}

func TestJobsController_Index_HappyPath(t *testing.T) {
	client, cleanup, ocrJobSpecFromFile, _, ereJobSpecFromFile, _ := setupJobSpecsControllerTestsWithJobs(t)
	defer cleanup()

	response, cleanup := client.Get("/v2/jobs")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	jobSpecs := []models.JobSpecV2{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &jobSpecs)
	assert.NoError(t, err)

	require.Len(t, jobSpecs, 2)

	runOCRJobSpecAssertions(t, ocrJobSpecFromFile, jobSpecs[0])
	runEthRequestEventJobSpecAssertions(t, ereJobSpecFromFile, jobSpecs[1])
}

func TestJobsController_Show_HappyPath(t *testing.T) {
	client, cleanup, ocrJobSpecFromFile, jobID, ereJobSpecFromFile, jobID2 := setupJobSpecsControllerTestsWithJobs(t)
	defer cleanup()

	response, cleanup := client.Get("/v2/jobs/" + fmt.Sprintf("%v", jobID))
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	ocrJobSpec := models.JobSpecV2{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &ocrJobSpec)
	assert.NoError(t, err)

	runOCRJobSpecAssertions(t, ocrJobSpecFromFile, ocrJobSpec)

	response, cleanup = client.Get("/v2/jobs/" + fmt.Sprintf("%v", jobID2))
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	ereJobSpec := models.JobSpecV2{}
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &ereJobSpec)
	assert.NoError(t, err)

	runEthRequestEventJobSpecAssertions(t, ereJobSpecFromFile, ereJobSpec)
}

func TestJobsController_Show_InvalidID(t *testing.T) {
	client, cleanup, _, _, _, _ := setupJobSpecsControllerTestsWithJobs(t)
	defer cleanup()

	response, cleanup := client.Get("/v2/jobs/uuidLikeString")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusUnprocessableEntity)
}

func TestJobsController_Show_NonExistentID(t *testing.T) {
	client, cleanup, _, _, _, _ := setupJobSpecsControllerTestsWithJobs(t)
	defer cleanup()

	response, cleanup := client.Get("/v2/jobs/999999999")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusNotFound)
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
	assert.Equal(t, ocrJobSpecFromFile.Pipeline.DOTSource, ocrJobSpecFromServer.PipelineSpec.DotDagSource)

	// Check that create and update dates are non empty values.
	// Empty date value is "0001-01-01 00:00:00 +0000 UTC" so we are checking for the
	// millenia and century characters to be present
	assert.Contains(t, ocrJobSpecFromServer.OffchainreportingOracleSpec.CreatedAt.String(), "20")
	assert.Contains(t, ocrJobSpecFromServer.OffchainreportingOracleSpec.UpdatedAt.String(), "20")
}

func runEthRequestEventJobSpecAssertions(t *testing.T, ereJobSpecFromFile services.EthRequestEventSpec, ereJobSpecFromServer models.JobSpecV2) {
	assert.Equal(t, ereJobSpecFromFile.ContractAddress, ereJobSpecFromServer.EthRequestEventSpec.ContractAddress)
	assert.Equal(t, ereJobSpecFromFile.Pipeline.DOTSource, ereJobSpecFromServer.PipelineSpec.DotDagSource)
	// Check that create and update dates are non empty values.
	// Empty date value is "0001-01-01 00:00:00 +0000 UTC" so we are checking for the
	// millenia and century characters to be present
	assert.Contains(t, ereJobSpecFromServer.EthRequestEventSpec.CreatedAt.String(), "20")
	assert.Contains(t, ereJobSpecFromServer.EthRequestEventSpec.UpdatedAt.String(), "20")
}

func setupJobsControllerTests(t *testing.T) (*cltest.TestApplication, cltest.HTTPClientCleaner, func()) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	return app, client, cleanup
}

func setupJobSpecsControllerTestsWithJobs(t *testing.T) (cltest.HTTPClientCleaner, func(), offchainreporting.OracleSpec, int32, services.EthRequestEventSpec, int32) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	var ocrJobSpecFromFile offchainreporting.OracleSpec
	tree, err := toml.LoadFile("testdata/oracle-spec.toml")
	require.NoError(t, err)
	err = tree.Unmarshal(&ocrJobSpecFromFile)
	require.NoError(t, err)
	jobID, _ := app.AddJobV2(context.Background(), ocrJobSpecFromFile, null.String{})

	var ereJobSpecFromFile services.EthRequestEventSpec
	tree, err = toml.LoadFile("testdata/eth-request-event-spec.toml")
	require.NoError(t, err)
	err = tree.Unmarshal(&ereJobSpecFromFile)
	require.NoError(t, err)
	jobID2, _ := app.AddJobV2(context.Background(), ereJobSpecFromFile, null.String{})

	return client, cleanup, ocrJobSpecFromFile, jobID, ereJobSpecFromFile, jobID2
}
