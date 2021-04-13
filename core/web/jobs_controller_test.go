package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/pelletier/go-toml"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestJobsController_Create_ValidationFailure_OffchainReportingSpec(t *testing.T) {
	var (
		contractAddress = cltest.NewEIP55Address()
	)

	var tt = []struct {
		name        string
		pid         models.PeerID
		kb          models.Sha256Hash
		taExists    bool
		expectedErr error
	}{
		{
			name:        "invalid keybundle",
			pid:         models.PeerID(cltest.DefaultP2PPeerID),
			kb:          models.Sha256Hash(cltest.Random32Byte()),
			taExists:    true,
			expectedErr: job.ErrNoSuchKeyBundle,
		},
		{
			name:        "invalid peerID",
			pid:         models.PeerID(cltest.NonExistentP2PPeerID),
			kb:          cltest.DefaultOCRKeyBundleIDSha256,
			taExists:    true,
			expectedErr: job.ErrNoSuchPeerID,
		},
		{
			name:        "invalid transmitter address",
			pid:         models.PeerID(cltest.DefaultP2PPeerID),
			kb:          cltest.DefaultOCRKeyBundleIDSha256,
			taExists:    false,
			expectedErr: job.ErrNoSuchTransmitterAddress,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ta, client, cleanup := setupJobsControllerTests(t)
			defer cleanup()

			var address models.EIP55Address
			if tc.taExists {
				key := cltest.MustInsertRandomKey(t, ta.Store.DB)
				address = key.Address
			} else {
				address = cltest.NewEIP55Address()
			}

			sp := cltest.MinimalOCRNonBootstrapSpec(contractAddress, address, tc.pid, tc.kb)
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

	toml := string(cltest.MustReadFile(t, "testdata/oracle-spec.toml"))
	toml = strings.Replace(toml, "0xF67D0290337bca0847005C7ffD1BC75BA9AAE6e4", app.Key.Address.Hex(), 1)
	body, _ := json.Marshal(models.CreateJobSpecRequest{
		TOML: toml,
	})
	response, cleanup := client.Post("/v2/jobs", bytes.NewReader(body))
	defer cleanup()
	require.Equal(t, http.StatusOK, response.StatusCode)

	jb := job.Job{}
	require.NoError(t, app.Store.DB.Preload("OffchainreportingOracleSpec").First(&jb).Error)

	resource := presenters.JobResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	assert.NoError(t, err)

	assert.Equal(t, "web oracle spec", jb.Name.ValueOrZero())
	assert.Equal(t, jb.OffchainreportingOracleSpec.P2PPeerID, resource.OffChainReportingSpec.P2PPeerID)
	assert.Equal(t, jb.OffchainreportingOracleSpec.P2PBootstrapPeers, resource.OffChainReportingSpec.P2PBootstrapPeers)
	assert.Equal(t, jb.OffchainreportingOracleSpec.IsBootstrapPeer, resource.OffChainReportingSpec.IsBootstrapPeer)
	assert.Equal(t, jb.OffchainreportingOracleSpec.EncryptedOCRKeyBundleID, resource.OffChainReportingSpec.EncryptedOCRKeyBundleID)
	assert.Equal(t, jb.OffchainreportingOracleSpec.TransmitterAddress, resource.OffChainReportingSpec.TransmitterAddress)
	assert.Equal(t, jb.OffchainreportingOracleSpec.ObservationTimeout, resource.OffChainReportingSpec.ObservationTimeout)
	assert.Equal(t, jb.OffchainreportingOracleSpec.BlockchainTimeout, resource.OffChainReportingSpec.BlockchainTimeout)
	assert.Equal(t, jb.OffchainreportingOracleSpec.ContractConfigTrackerSubscribeInterval, resource.OffChainReportingSpec.ContractConfigTrackerSubscribeInterval)
	assert.Equal(t, jb.OffchainreportingOracleSpec.ContractConfigTrackerSubscribeInterval, resource.OffChainReportingSpec.ContractConfigTrackerSubscribeInterval)
	assert.Equal(t, jb.OffchainreportingOracleSpec.ContractConfigConfirmations, resource.OffChainReportingSpec.ContractConfigConfirmations)
	assert.NotNil(t, resource.PipelineSpec.DotDAGSource)

	// Sanity check to make sure it inserted correctly
	require.Equal(t, models.EIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C"), jb.OffchainreportingOracleSpec.ContractAddress)
}

func TestJobsController_Create_HappyPath_KeeperSpec(t *testing.T) {
	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	tomlBytes := cltest.MustReadFile(t, "testdata/keeper-spec.toml")
	body, _ := json.Marshal(models.CreateJobSpecRequest{
		TOML: string(tomlBytes),
	})
	response, cleanup := client.Post("/v2/jobs", bytes.NewReader(body))
	defer cleanup()
	require.Equal(t, http.StatusOK, response.StatusCode)

	jb := job.Job{}
	require.NoError(t, app.Store.DB.Preload("KeeperSpec").First(&jb).Error)

	resource := presenters.JobResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	assert.NoError(t, err)

	require.Equal(t, resource.KeeperSpec.ContractAddress, jb.KeeperSpec.ContractAddress)
	require.Equal(t, resource.KeeperSpec.FromAddress, jb.KeeperSpec.FromAddress)
	assert.Equal(t, "example keeper spec", jb.Name.ValueOrZero())

	// Sanity check to make sure it inserted correctly
	require.Equal(t, models.EIP55Address("0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba"), jb.KeeperSpec.ContractAddress)
	require.Equal(t, models.EIP55Address("0xa8037A20989AFcBC51798de9762b351D63ff462e"), jb.KeeperSpec.FromAddress)
}

func TestJobsController_Create_HappyPath_DirectRequestSpec(t *testing.T) {
	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()
	require.NoError(t, app.Start())
	gethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(cltest.EmptyMockSubscription(), nil)

	client := app.NewHTTPClient()

	tomlBytes := cltest.MustReadFile(t, "testdata/direct-request-spec.toml")
	body, _ := json.Marshal(models.CreateJobSpecRequest{
		TOML: string(tomlBytes),
	})
	response, cleanup := client.Post("/v2/jobs", bytes.NewReader(body))
	defer cleanup()
	require.Equal(t, http.StatusOK, response.StatusCode)

	jb := job.Job{}
	require.NoError(t, app.Store.DB.Preload("DirectRequestSpec").First(&jb).Error)

	resource := presenters.JobResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	assert.NoError(t, err)

	assert.Equal(t, "example eth request event spec", jb.Name.ValueOrZero())
	assert.NotNil(t, resource.PipelineSpec.DotDAGSource)

	// Sanity check to make sure it inserted correctly
	require.Equal(t, models.EIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C"), jb.DirectRequestSpec.ContractAddress)

	require.NotZero(t, jb.DirectRequestSpec.OnChainJobSpecID[:])
}

func TestJobsController_Create_HappyPath_FluxMonitorSpec(t *testing.T) {
	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()
	require.NoError(t, app.Start())
	gethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(cltest.EmptyMockSubscription(), nil)

	client := app.NewHTTPClient()

	tomlBytes := cltest.MustReadFile(t, "testdata/flux-monitor-spec.toml")
	body, _ := json.Marshal(models.CreateJobSpecRequest{
		TOML: string(tomlBytes),
	})

	response, cleanup := client.Post("/v2/jobs", bytes.NewReader(body))
	defer cleanup()
	require.Equal(t, http.StatusOK, response.StatusCode)

	jb := job.Job{}
	require.NoError(t, app.Store.DB.Preload("FluxMonitorSpec").First(&jb).Error)

	resource := presenters.JobResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	assert.NoError(t, err)
	t.Log()

	assert.Equal(t, "example flux monitor spec", jb.Name.ValueOrZero())
	assert.NotNil(t, resource.PipelineSpec.DotDAGSource)
	assert.Equal(t, models.EIP55Address("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"), jb.FluxMonitorSpec.ContractAddress)
	assert.Equal(t, time.Second, jb.FluxMonitorSpec.IdleTimerPeriod)
	assert.Equal(t, false, jb.FluxMonitorSpec.IdleTimerDisabled)
	assert.Equal(t, int32(2), jb.FluxMonitorSpec.Precision)
	assert.Equal(t, float32(0.5), jb.FluxMonitorSpec.Threshold)
	assert.Equal(t, float32(0), jb.FluxMonitorSpec.AbsoluteThreshold)
}
func TestJobsController_Index_HappyPath(t *testing.T) {
	client, cleanup, ocrJobSpecFromFile, _, ereJobSpecFromFile, _ := setupJobSpecsControllerTestsWithJobs(t)
	defer cleanup()

	response, cleanup := client.Get("/v2/jobs")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	resources := []presenters.JobResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	assert.NoError(t, err)

	require.Len(t, resources, 2)

	runOCRJobSpecAssertions(t, ocrJobSpecFromFile, resources[0])
	runDirectRequestJobSpecAssertions(t, ereJobSpecFromFile, resources[1])
}

func TestJobsController_Show_HappyPath(t *testing.T) {
	client, cleanup, ocrJobSpecFromFile, jobID, ereJobSpecFromFile, jobID2 := setupJobSpecsControllerTestsWithJobs(t)
	defer cleanup()

	response, cleanup := client.Get("/v2/jobs/" + fmt.Sprintf("%v", jobID))
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	ocrJob := presenters.JobResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &ocrJob)
	assert.NoError(t, err)

	runOCRJobSpecAssertions(t, ocrJobSpecFromFile, ocrJob)

	response, cleanup = client.Get("/v2/jobs/" + fmt.Sprintf("%v", jobID2))
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	ereJob := presenters.JobResource{}
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &ereJob)
	assert.NoError(t, err)

	runDirectRequestJobSpecAssertions(t, ereJobSpecFromFile, ereJob)
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

func runOCRJobSpecAssertions(t *testing.T, ocrJobSpecFromFileDB job.Job, ocrJobSpecFromServer presenters.JobResource) {
	ocrJobSpecFromFile := ocrJobSpecFromFileDB.OffchainreportingOracleSpec
	assert.Equal(t, ocrJobSpecFromFile.ContractAddress, ocrJobSpecFromServer.OffChainReportingSpec.ContractAddress)
	assert.Equal(t, ocrJobSpecFromFile.P2PPeerID, ocrJobSpecFromServer.OffChainReportingSpec.P2PPeerID)
	assert.Equal(t, ocrJobSpecFromFile.P2PBootstrapPeers, ocrJobSpecFromServer.OffChainReportingSpec.P2PBootstrapPeers)
	assert.Equal(t, ocrJobSpecFromFile.IsBootstrapPeer, ocrJobSpecFromServer.OffChainReportingSpec.IsBootstrapPeer)
	assert.Equal(t, ocrJobSpecFromFile.EncryptedOCRKeyBundleID, ocrJobSpecFromServer.OffChainReportingSpec.EncryptedOCRKeyBundleID)
	assert.Equal(t, ocrJobSpecFromFile.TransmitterAddress, ocrJobSpecFromServer.OffChainReportingSpec.TransmitterAddress)
	assert.Equal(t, ocrJobSpecFromFile.ObservationTimeout, ocrJobSpecFromServer.OffChainReportingSpec.ObservationTimeout)
	assert.Equal(t, ocrJobSpecFromFile.BlockchainTimeout, ocrJobSpecFromServer.OffChainReportingSpec.BlockchainTimeout)
	assert.Equal(t, ocrJobSpecFromFile.ContractConfigTrackerSubscribeInterval, ocrJobSpecFromServer.OffChainReportingSpec.ContractConfigTrackerSubscribeInterval)
	assert.Equal(t, ocrJobSpecFromFile.ContractConfigTrackerSubscribeInterval, ocrJobSpecFromServer.OffChainReportingSpec.ContractConfigTrackerSubscribeInterval)
	assert.Equal(t, ocrJobSpecFromFile.ContractConfigConfirmations, ocrJobSpecFromServer.OffChainReportingSpec.ContractConfigConfirmations)
	assert.Equal(t, ocrJobSpecFromFileDB.Pipeline.DOTSource, ocrJobSpecFromServer.PipelineSpec.DotDAGSource)

	// Check that create and update dates are non empty values.
	// Empty date value is "0001-01-01 00:00:00 +0000 UTC" so we are checking for the
	// millenia and century characters to be present
	assert.Contains(t, ocrJobSpecFromServer.OffChainReportingSpec.CreatedAt.String(), "20")
	assert.Contains(t, ocrJobSpecFromServer.OffChainReportingSpec.UpdatedAt.String(), "20")
}

func runDirectRequestJobSpecAssertions(t *testing.T, ereJobSpecFromFile job.Job, ereJobSpecFromServer presenters.JobResource) {
	assert.Equal(t, ereJobSpecFromFile.DirectRequestSpec.ContractAddress, ereJobSpecFromServer.DirectRequestSpec.ContractAddress)
	assert.Equal(t, ereJobSpecFromFile.Pipeline.DOTSource, ereJobSpecFromServer.PipelineSpec.DotDAGSource)
	// Check that create and update dates are non empty values.
	// Empty date value is "0001-01-01 00:00:00 +0000 UTC" so we are checking for the
	// millenia and century characters to be present
	assert.Contains(t, ereJobSpecFromServer.DirectRequestSpec.CreatedAt.String(), "20")
	assert.Contains(t, ereJobSpecFromServer.DirectRequestSpec.UpdatedAt.String(), "20")
}

func setupJobsControllerTests(t *testing.T) (*cltest.TestApplication, cltest.HTTPClientCleaner, func()) {
	t.Parallel()
	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	require.NoError(t, app.Start())

	_, bridge := cltest.NewBridgeType(t, "voter_turnout", "http://blah.com")
	require.NoError(t, app.Store.DB.Create(bridge).Error)
	_, bridge2 := cltest.NewBridgeType(t, "election_winner", "http://blah.com")
	require.NoError(t, app.Store.DB.Create(bridge2).Error)
	client := app.NewHTTPClient()
	return app, client, cleanup
}

func setupJobSpecsControllerTestsWithJobs(t *testing.T) (cltest.HTTPClientCleaner, func(), job.Job, int32, job.Job, int32) {
	t.Parallel()
	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	require.NoError(t, app.Start())

	_, bridge := cltest.NewBridgeType(t, "voter_turnout", "http://blah.com")
	require.NoError(t, app.Store.DB.Create(bridge).Error)
	_, bridge2 := cltest.NewBridgeType(t, "election_winner", "http://blah.com")
	require.NoError(t, app.Store.DB.Create(bridge2).Error)

	client := app.NewHTTPClient()

	var ocrJobSpecFromFileDB job.Job
	tree, err := toml.LoadFile("testdata/oracle-spec.toml")
	require.NoError(t, err)
	err = tree.Unmarshal(&ocrJobSpecFromFileDB)
	require.NoError(t, err)
	var ocrSpec job.OffchainReportingOracleSpec
	err = tree.Unmarshal(&ocrSpec)
	require.NoError(t, err)
	ocrJobSpecFromFileDB.OffchainreportingOracleSpec = &ocrSpec
	ocrJobSpecFromFileDB.OffchainreportingOracleSpec.TransmitterAddress = &app.Key.Address
	jobID, _ := app.AddJobV2(context.Background(), ocrJobSpecFromFileDB, null.String{})

	ereJobSpecFromFileDB, err := directrequest.ValidatedDirectRequestSpec(string(cltest.MustReadFile(t, "testdata/direct-request-spec.toml")))
	require.NoError(t, err)
	jobID2, _ := app.AddJobV2(context.Background(), ereJobSpecFromFileDB, null.String{})

	return client, cleanup, ocrJobSpecFromFileDB, jobID, ereJobSpecFromFileDB, jobID2
}
