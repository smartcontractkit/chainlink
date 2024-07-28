package web_test

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/pelletier/go-toml"
	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils/tomlutils"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestJobsController_Create_ValidationFailure_OffchainReportingSpec(t *testing.T) {
	var (
		contractAddress = cltest.NewEIP55Address()
	)

	var peerID ragep2ptypes.PeerID
	require.NoError(t, peerID.UnmarshalText([]byte(configtest.DefaultPeerID)))
	randomBytes := testutils.Random32Byte()

	var tt = []struct {
		name        string
		pid         p2pkey.PeerID
		kb          string
		taExists    bool
		expectedErr error
	}{
		{
			name:        "invalid keybundle",
			pid:         p2pkey.PeerID(peerID),
			kb:          hex.EncodeToString(randomBytes[:]),
			taExists:    true,
			expectedErr: job.ErrNoSuchKeyBundle,
		},
		{
			name:        "invalid transmitter address",
			pid:         p2pkey.PeerID(peerID),
			kb:          cltest.DefaultOCRKeyBundleID,
			taExists:    false,
			expectedErr: job.ErrNoSuchTransmitterKey,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testutils.Context(t)
			ta, client := setupJobsControllerTests(t)

			var address types.EIP55Address
			if tc.taExists {
				key, _ := cltest.MustInsertRandomKey(t, ta.KeyStore.Eth())
				address = key.EIP55Address
			} else {
				address = cltest.NewEIP55Address()
			}

			require.NoError(t, ta.KeyStore.OCR().Add(ctx, cltest.DefaultOCRKey))

			sp := cltest.MinimalOCRNonBootstrapSpec(contractAddress, address, tc.pid, tc.kb)
			body, _ := json.Marshal(web.CreateJobRequest{
				TOML: sp,
			})
			resp, cleanup := client.Post("/v2/jobs", bytes.NewReader(body))
			t.Cleanup(cleanup)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			b, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(b), tc.expectedErr.Error())
		})
	}
}

func TestJobController_Create_DirectRequest_Fast(t *testing.T) {
	ctx := testutils.Context(t)
	app, client := setupJobsControllerTests(t)
	require.NoError(t, app.KeyStore.OCR().Add(ctx, cltest.DefaultOCRKey))

	n := 10

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			body, err := json.Marshal(web.CreateJobRequest{
				TOML: fmt.Sprintf(testspecs.DirectRequestSpecNoExternalJobID, i),
			})
			require.NoError(t, err)

			t.Logf("POSTing %d", i)
			r, cleanup := client.Post("/v2/jobs", bytes.NewReader(body))
			defer cleanup()
			require.Equal(t, http.StatusOK, r.StatusCode)
		}(i)
	}
	wg.Wait()
	cltest.AssertCount(t, app.GetDB(), "direct_request_specs", int64(n))
}

func mustInt32FromString(t *testing.T, s string) int32 {
	i, err := strconv.ParseInt(s, 10, 32)
	require.NoError(t, err)
	return int32(i)
}

func TestJobController_Create_HappyPath(t *testing.T) {
	ctx := testutils.Context(t)
	app, client := setupJobsControllerTests(t)
	b1, b2 := setupBridges(t, app.GetDB())
	require.NoError(t, app.KeyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	var pks []vrfkey.KeyV2
	var k []p2pkey.KeyV2
	{
		var err error
		pks, err = app.KeyStore.VRF().GetAll()
		require.NoError(t, err)
		require.Len(t, pks, 1)
		k, err = app.KeyStore.P2P().GetAll()
		require.NoError(t, err)
		require.Len(t, k, 1)
	}

	jorm := app.JobORM()
	var tt = []struct {
		name         string
		tomlTemplate func(nameAndExternalJobID string) string
		assertion    func(t *testing.T, nameAndExternalJobID string, r *http.Response)
	}{
		{
			name: "offchain reporting",
			tomlTemplate: func(nameAndExternalJobID string) string {
				return testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
					TransmitterAddress: app.Keys[0].Address.Hex(),
					DS1BridgeName:      b1,
					DS2BridgeName:      b2,
					Name:               nameAndExternalJobID,
				}).Toml()
			},
			assertion: func(t *testing.T, nameAndExternalJobID string, r *http.Response) {
				require.Equal(t, http.StatusOK, r.StatusCode)

				resource := presenters.JobResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, r), &resource)
				assert.NoError(t, err)

				jb, err := jorm.FindJob(testutils.Context(t), mustInt32FromString(t, resource.ID))
				require.NoError(t, err)
				require.NotNil(t, resource.OffChainReportingSpec)

				assert.Equal(t, nameAndExternalJobID, jb.Name.ValueOrZero())
				assert.Equal(t, jb.OCROracleSpec.P2PV2Bootstrappers, resource.OffChainReportingSpec.P2PV2Bootstrappers)
				assert.Equal(t, jb.OCROracleSpec.IsBootstrapPeer, resource.OffChainReportingSpec.IsBootstrapPeer)
				assert.Equal(t, jb.OCROracleSpec.EncryptedOCRKeyBundleID, resource.OffChainReportingSpec.EncryptedOCRKeyBundleID)
				assert.Equal(t, jb.OCROracleSpec.TransmitterAddress, resource.OffChainReportingSpec.TransmitterAddress)
				assert.Equal(t, jb.OCROracleSpec.ObservationTimeout, resource.OffChainReportingSpec.ObservationTimeout)
				assert.Equal(t, jb.OCROracleSpec.BlockchainTimeout, resource.OffChainReportingSpec.BlockchainTimeout)
				assert.Equal(t, jb.OCROracleSpec.ContractConfigTrackerSubscribeInterval, resource.OffChainReportingSpec.ContractConfigTrackerSubscribeInterval)
				assert.Equal(t, jb.OCROracleSpec.ContractConfigTrackerSubscribeInterval, resource.OffChainReportingSpec.ContractConfigTrackerSubscribeInterval)
				assert.Equal(t, jb.OCROracleSpec.ContractConfigConfirmations, resource.OffChainReportingSpec.ContractConfigConfirmations)
				assert.NotNil(t, resource.PipelineSpec.DotDAGSource)
				// Sanity check to make sure it inserted correctly
				require.Equal(t, types.EIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C"), jb.OCROracleSpec.ContractAddress)
			},
		},
		{
			name: "keeper",
			tomlTemplate: func(nameAndExternalJobID string) string {
				return fmt.Sprintf(`
                                  type                        = "keeper"
                                  schemaVersion               = 1
                                  name                        = "%s"
                                  contractAddress             = "0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba"
                                  fromAddress                 = "0xa8037A20989AFcBC51798de9762b351D63ff462e"
                                  evmChainID                  = 0
                                  minIncomingConfigurations   = 1
                                  externalJobID               = "%s"
                             `, nameAndExternalJobID, nameAndExternalJobID)
			},
			assertion: func(t *testing.T, nameAndExternalJobID string, r *http.Response) {
				require.Equal(t, http.StatusInternalServerError, r.StatusCode)

				errs := cltest.ParseJSONAPIErrors(t, r.Body)
				require.NotNil(t, errs)
				require.Len(t, errs.Errors, 1)
				// services failed to start
				require.Contains(t, errs.Errors[0].Detail, "no contract code at given address")
				// but the job should still exist
				ctx := testutils.Context(t)
				jb, err := jorm.FindJobByExternalJobID(ctx, uuid.MustParse(nameAndExternalJobID))
				require.NoError(t, err)
				require.NotNil(t, jb.KeeperSpec)

				require.Equal(t, types.EIP55Address("0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba"), jb.KeeperSpec.ContractAddress)
				require.Equal(t, types.EIP55Address("0xa8037A20989AFcBC51798de9762b351D63ff462e"), jb.KeeperSpec.FromAddress)
				assert.Equal(t, nameAndExternalJobID, jb.Name.ValueOrZero())

				// Sanity check to make sure it inserted correctly
				require.Equal(t, types.EIP55Address("0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba"), jb.KeeperSpec.ContractAddress)
				require.Equal(t, types.EIP55Address("0xa8037A20989AFcBC51798de9762b351D63ff462e"), jb.KeeperSpec.FromAddress)
			},
		},
		{
			name: "cron",
			tomlTemplate: func(nameAndExternalJobID string) string {
				return fmt.Sprintf(testspecs.CronSpecTemplate, nameAndExternalJobID)
			},
			assertion: func(t *testing.T, nameAndExternalJobID string, r *http.Response) {
				require.Equal(t, http.StatusOK, r.StatusCode)
				resource := presenters.JobResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, r), &resource)
				assert.NoError(t, err)

				jb, err := jorm.FindJob(testutils.Context(t), mustInt32FromString(t, resource.ID))
				require.NoError(t, err)
				require.NotNil(t, jb.CronSpec)

				assert.NotNil(t, resource.PipelineSpec.DotDAGSource)
				require.Equal(t, "CRON_TZ=UTC * 0 0 1 1 *", jb.CronSpec.CronSchedule)
			},
		},
		{
			name: "cron-dot-separator",
			tomlTemplate: func(nameAndExternalJobID string) string {
				return fmt.Sprintf(testspecs.CronSpecDotSepTemplate, nameAndExternalJobID)
			},
			assertion: func(t *testing.T, nameAndExternalJobID string, r *http.Response) {
				require.Equal(t, http.StatusOK, r.StatusCode)
				resource := presenters.JobResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, r), &resource)
				assert.NoError(t, err)

				jb, err := jorm.FindJob(testutils.Context(t), mustInt32FromString(t, resource.ID))
				require.NoError(t, err)
				require.NotNil(t, jb.CronSpec)

				assert.NotNil(t, resource.PipelineSpec.DotDAGSource)
				require.Equal(t, "CRON_TZ=UTC * 0 0 1 1 *", jb.CronSpec.CronSchedule)
			},
		},
		{
			name: "cron-evm-chain-id",
			tomlTemplate: func(nameAndExternalJobID string) string {
				return fmt.Sprintf(testspecs.CronSpecEVMChainIDTemplate, nameAndExternalJobID)
			},
			assertion: func(t *testing.T, nameAndExternalJobID string, r *http.Response) {
				require.Equal(t, http.StatusOK, r.StatusCode)
				resource := presenters.JobResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, r), &resource)
				assert.NoError(t, err)

				jb, err := jorm.FindJob(testutils.Context(t), mustInt32FromString(t, resource.ID))
				require.NoError(t, err)
				require.NotNil(t, jb.CronSpec)

				assert.NotNil(t, resource.PipelineSpec.DotDAGSource)
				require.Equal(t, ubig.NewI(42), jb.CronSpec.EVMChainID)
			},
		},
		{
			name: "directrequest",
			tomlTemplate: func(nameAndExternalJobID string) string {
				return testspecs.GetDirectRequestSpecWithUUID(uuid.MustParse(nameAndExternalJobID))
			},
			assertion: func(t *testing.T, nameAndExternalJobID string, r *http.Response) {
				require.Equal(t, http.StatusOK, r.StatusCode)
				resource := presenters.JobResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, r), &resource)
				assert.NoError(t, err)

				jb, err := jorm.FindJob(testutils.Context(t), mustInt32FromString(t, resource.ID))
				require.NoError(t, err)
				require.NotNil(t, jb.DirectRequestSpec)

				assert.Equal(t, nameAndExternalJobID, jb.Name.ValueOrZero())
				assert.NotNil(t, resource.PipelineSpec.DotDAGSource)
				// Sanity check to make sure it inserted correctly
				require.Equal(t, types.EIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C"), jb.DirectRequestSpec.ContractAddress)
				require.Equal(t, jb.ExternalJobID.String(), nameAndExternalJobID)
			},
		},
		{
			name: "directrequest-with-requesters-and-min-contract-payment",
			tomlTemplate: func(nameAndExternalJobID string) string {
				return fmt.Sprintf(testspecs.DirectRequestSpecWithRequestersAndMinContractPaymentTemplate, nameAndExternalJobID, nameAndExternalJobID)
			},
			assertion: func(t *testing.T, nameAndExternalJobID string, r *http.Response) {
				require.Equal(t, http.StatusOK, r.StatusCode)
				resource := presenters.JobResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, r), &resource)
				assert.NoError(t, err)

				jb, err := jorm.FindJob(testutils.Context(t), mustInt32FromString(t, resource.ID))
				require.NoError(t, err)
				require.NotNil(t, jb.DirectRequestSpec)

				assert.Equal(t, nameAndExternalJobID, jb.Name.ValueOrZero())
				assert.NotNil(t, resource.PipelineSpec.DotDAGSource)
				assert.NotNil(t, resource.DirectRequestSpec.Requesters)
				assert.Equal(t, "1000000000000000000000", resource.DirectRequestSpec.MinContractPayment.String())
				// Check requesters got saved properly
				require.EqualValues(t, []common.Address{common.HexToAddress("0xAaAA1F8ee20f5565510b84f9353F1E333e753B7a"), common.HexToAddress("0xBbBb70f0E81c6F3430dfDc9fa02fB22bDD818c4E")}, jb.DirectRequestSpec.Requesters)
				require.Equal(t, "1000000000000000000000", jb.DirectRequestSpec.MinContractPayment.String())
				require.Equal(t, jb.ExternalJobID.String(), nameAndExternalJobID)
			},
		},
		{
			name: "fluxmonitor",
			tomlTemplate: func(nameAndExternalJobID string) string {
				return fmt.Sprintf(testspecs.FluxMonitorSpecTemplate, nameAndExternalJobID, nameAndExternalJobID)
			},
			assertion: func(t *testing.T, nameAndExternalJobID string, r *http.Response) {
				require.Equal(t, http.StatusInternalServerError, r.StatusCode)

				errs := cltest.ParseJSONAPIErrors(t, r.Body)
				require.NotNil(t, errs)
				require.Len(t, errs.Errors, 1)
				// services failed to start
				require.Contains(t, errs.Errors[0].Detail, "no contract code at given address")
				// but the job should still exist
				ctx := testutils.Context(t)
				jb, err := jorm.FindJobByExternalJobID(ctx, uuid.MustParse(nameAndExternalJobID))
				require.NoError(t, err)
				require.NotNil(t, jb.FluxMonitorSpec)

				assert.Equal(t, nameAndExternalJobID, jb.Name.ValueOrZero())
				assert.NotNil(t, jb.PipelineSpec.DotDagSource)
				assert.Equal(t, types.EIP55Address("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"), jb.FluxMonitorSpec.ContractAddress)
				assert.Equal(t, time.Second, jb.FluxMonitorSpec.IdleTimerPeriod)
				assert.Equal(t, false, jb.FluxMonitorSpec.IdleTimerDisabled)
				assert.Equal(t, tomlutils.Float32(0.5), jb.FluxMonitorSpec.Threshold)
				assert.Equal(t, tomlutils.Float32(0), jb.FluxMonitorSpec.AbsoluteThreshold)
			},
		},
		{
			name: "vrf",
			tomlTemplate: func(_ string) string {
				return testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{PublicKey: pks[0].PublicKey.String()}).Toml()
			},
			assertion: func(t *testing.T, nameAndExternalJobID string, r *http.Response) {
				require.Equal(t, http.StatusOK, r.StatusCode)
				resp := cltest.ParseResponseBody(t, r)
				resource := presenters.JobResource{}
				err := web.ParseJSONAPIResponse(resp, &resource)
				require.NoError(t, err)

				jb, err := jorm.FindJob(testutils.Context(t), mustInt32FromString(t, resource.ID))
				require.NoError(t, err)
				require.NotNil(t, jb.VRFSpec)

				assert.NotNil(t, resource.PipelineSpec.DotDAGSource)
				assert.Equal(t, uint32(6), resource.VRFSpec.MinIncomingConfirmations)
				assert.Equal(t, jb.VRFSpec.MinIncomingConfirmations, resource.VRFSpec.MinIncomingConfirmations)
				assert.Equal(t, "0xABA5eDc1a551E55b1A570c0e1f1055e5BE11eca7", resource.VRFSpec.CoordinatorAddress.Hex())
				assert.Equal(t, jb.VRFSpec.CoordinatorAddress.Hex(), resource.VRFSpec.CoordinatorAddress.Hex())
			},
		},
		{
			name: "stream",
			tomlTemplate: func(_ string) string {
				return testspecs.GenerateStreamSpec(testspecs.StreamSpecParams{Name: "ETH/USD", StreamID: 32}).Toml()
			},
			assertion: func(t *testing.T, nameAndExternalJobID string, r *http.Response) {
				require.Equal(t, http.StatusOK, r.StatusCode)
				resp := cltest.ParseResponseBody(t, r)
				resource := presenters.JobResource{}
				err := web.ParseJSONAPIResponse(resp, &resource)
				require.NoError(t, err)

				jb, err := jorm.FindJob(testutils.Context(t), mustInt32FromString(t, resource.ID))
				require.NoError(t, err)
				require.NotNil(t, jb.PipelineSpec)

				assert.NotNil(t, resource.PipelineSpec.DotDAGSource)
				assert.Equal(t, jb.Name.ValueOrZero(), resource.Name)
				assert.Equal(t, jb.StreamID, resource.StreamID)
			},
		},
		{
			name: "workflow",
			tomlTemplate: func(_ string) string {
				workflow := `
triggers:
  - id: "mercury-trigger@1.0.0"
    config:
      feedIds:
        - "0x1111111111111111111100000000000000000000000000000000000000000000"
        - "0x2222222222222222222200000000000000000000000000000000000000000000"
        - "0x3333333333333333333300000000000000000000000000000000000000000000"

consensus:
  - id: "offchain_reporting@2.0.0"
    ref: "evm_median"
    inputs:
      observations:
        - "$(trigger.outputs)"
    config:
      aggregation_method: "data_feeds_2_0"
      aggregation_config:
        "0x1111111111111111111100000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: 3600
        "0x2222222222222222222200000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: 3600
        "0x3333333333333333333300000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: 3600
      encoder: "EVM"
      encoder_config:
        abi: "mercury_reports bytes[]"

targets:
  - id: "write_polygon-testnet-mumbai@3.0.0"
    inputs:
      report: "$(evm_median.outputs.report)"
    config:
      address: "0x3F3554832c636721F1fD1822Ccca0354576741Ef"
      params: ["$(report)"]
      abi: "receive(report bytes)"
  - id: "write_ethereum-testnet-sepolia@4.0.0"
    inputs:
      report: "$(evm_median.outputs.report)"
    config:
      address: "0x54e220867af6683aE6DcBF535B4f952cB5116510"
      params: ["$(report)"]
      abi: "receive(report bytes)"
`
				return testspecs.GenerateWorkflowJobSpec(t, workflow).Toml()
			},
			assertion: func(t *testing.T, nameAndExternalJobID string, r *http.Response) {
				require.Equal(t, http.StatusOK, r.StatusCode)
				resp := cltest.ParseResponseBody(t, r)
				resource := presenters.JobResource{}
				err := web.ParseJSONAPIResponse(resp, &resource)
				require.NoError(t, err, "failed to parse response body: %s", resp)

				jb, err := jorm.FindJob(testutils.Context(t), mustInt32FromString(t, resource.ID))
				require.NoError(t, err)
				require.NotNil(t, jb.WorkflowSpec)

				assert.Equal(t, jb.WorkflowSpec.Workflow, resource.WorkflowSpec.Workflow)
				assert.Equal(t, jb.WorkflowSpec.WorkflowID, resource.WorkflowSpec.WorkflowID)
				assert.Equal(t, jb.WorkflowSpec.WorkflowOwner, resource.WorkflowSpec.WorkflowOwner)
				assert.Equal(t, jb.WorkflowSpec.WorkflowName, resource.WorkflowSpec.WorkflowName)
			},
		},
	}
	for _, tc := range tt {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			nameAndExternalJobID := uuid.New().String()
			toml := c.tomlTemplate(nameAndExternalJobID)
			t.Log("Job toml:", toml)
			body, err := json.Marshal(web.CreateJobRequest{
				TOML: toml,
			})
			require.NoError(t, err)
			response, cleanup := client.Post("/v2/jobs", bytes.NewReader(body))
			defer cleanup()
			c.assertion(t, nameAndExternalJobID, response)
		})
	}
}

func TestJobsController_Create_WebhookSpec(t *testing.T) {
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	_, fetchBridge := cltest.MustCreateBridge(t, app.GetDB(), cltest.BridgeOpts{})
	_, submitBridge := cltest.MustCreateBridge(t, app.GetDB(), cltest.BridgeOpts{})

	client := app.NewHTTPClient(nil)

	tomlStr := testspecs.GetWebhookSpecNoBody(uuid.New(), fetchBridge.Name.String(), submitBridge.Name.String())
	body, _ := json.Marshal(web.CreateJobRequest{
		TOML: tomlStr,
	})
	response, cleanup := client.Post("/v2/jobs", bytes.NewReader(body))
	defer cleanup()
	require.Equal(t, http.StatusOK, response.StatusCode)
	resource := presenters.JobResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	require.NoError(t, err)
	assert.NotNil(t, resource.PipelineSpec.DotDAGSource)

	jorm := app.JobORM()
	_, err = jorm.FindJob(testutils.Context(t), mustInt32FromString(t, resource.ID))
	require.NoError(t, err)
}

//go:embed webhook-spec-template.yml
var webhookSpecTemplate string

func TestJobsController_FailToCreate_EmptyJsonAttribute(t *testing.T) {
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)

	nameAndExternalJobID := uuid.New()
	spec := fmt.Sprintf(webhookSpecTemplate, nameAndExternalJobID, nameAndExternalJobID)
	body, err := json.Marshal(web.CreateJobRequest{
		TOML: spec,
	})
	require.NoError(t, err)
	response, cleanup := client.Post("/v2/jobs", bytes.NewReader(body))
	defer cleanup()

	b, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Contains(t, string(b), "syntax is not supported. Please use \\\"{}\\\" instead")
}

func TestJobsController_Index_HappyPath(t *testing.T) {
	_, client, ocrJobSpecFromFile, _, ereJobSpecFromFile, _ := setupJobSpecsControllerTestsWithJobs(t)

	url := url.URL{Path: "/v2/jobs"}
	query := url.Query()
	query.Set("evmChainID", cltest.FixtureChainID.String())
	url.RawQuery = query.Encode()

	response, cleanup := client.Get(url.String())
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	var resources []presenters.JobResource
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	assert.NoError(t, err)

	require.Len(t, resources, 2)

	runDirectRequestJobSpecAssertions(t, ereJobSpecFromFile, resources[0])
	runOCRJobSpecAssertions(t, ocrJobSpecFromFile, resources[1])
}

func TestJobsController_Show_HappyPath(t *testing.T) {
	_, client, ocrJobSpecFromFile, jobID, ereJobSpecFromFile, jobID2 := setupJobSpecsControllerTestsWithJobs(t)

	response, cleanup := client.Get("/v2/jobs/" + fmt.Sprintf("%v", jobID))
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	ocrJob := presenters.JobResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &ocrJob)
	assert.NoError(t, err)

	runOCRJobSpecAssertions(t, ocrJobSpecFromFile, ocrJob)

	response, cleanup = client.Get("/v2/jobs/" + ocrJobSpecFromFile.ExternalJobID.String())
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	ocrJob = presenters.JobResource{}
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &ocrJob)
	assert.NoError(t, err)

	runOCRJobSpecAssertions(t, ocrJobSpecFromFile, ocrJob)

	response, cleanup = client.Get("/v2/jobs/" + fmt.Sprintf("%v", jobID2))
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	ereJob := presenters.JobResource{}
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &ereJob)
	assert.NoError(t, err)

	runDirectRequestJobSpecAssertions(t, ereJobSpecFromFile, ereJob)

	response, cleanup = client.Get("/v2/jobs/" + ereJobSpecFromFile.ExternalJobID.String())
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	ereJob = presenters.JobResource{}
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &ereJob)
	assert.NoError(t, err)

	runDirectRequestJobSpecAssertions(t, ereJobSpecFromFile, ereJob)
}

func TestJobsController_Show_InvalidID(t *testing.T) {
	_, client, _, _, _, _ := setupJobSpecsControllerTestsWithJobs(t)

	response, cleanup := client.Get("/v2/jobs/uuidLikeString")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusUnprocessableEntity)
}

func TestJobsController_Show_NonExistentID(t *testing.T) {
	_, client, _, _, _, _ := setupJobSpecsControllerTestsWithJobs(t)

	response, cleanup := client.Get("/v2/jobs/999999999")
	t.Cleanup(cleanup)

	cltest.AssertServerResponse(t, response, http.StatusNotFound)
}

func TestJobsController_Update_HappyPath(t *testing.T) {
	ctx := testutils.Context(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.OCR.Enabled = ptr(true)
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", freeport.GetOne(t))}
		c.P2P.PeerID = &cltest.DefaultP2PPeerID
	})
	app := cltest.NewApplicationWithConfigAndKey(t, cfg, cltest.DefaultP2PKey)

	require.NoError(t, app.KeyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	require.NoError(t, app.Start(ctx))

	_, bridge := cltest.MustCreateBridge(t, app.GetDB(), cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, app.GetDB(), cltest.BridgeOpts{})

	client := app.NewHTTPClient(nil)

	var jb job.Job
	ocrspec := testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
		DS1BridgeName: bridge.Name.String(),
		DS2BridgeName: bridge2.Name.String(),
		Name:          "old OCR job",
	})
	err := toml.Unmarshal([]byte(ocrspec.Toml()), &jb)
	require.NoError(t, err)

	// BCF-2095
	// disable fkey checks until the end of the test transaction
	require.NoError(t, utils.JustError(
		app.GetDB().ExecContext(ctx, `SET CONSTRAINTS job_spec_errors_v2_job_id_fkey DEFERRED`)))

	var ocrSpec job.OCROracleSpec
	err = toml.Unmarshal([]byte(ocrspec.Toml()), &ocrSpec)
	require.NoError(t, err)
	jb.OCROracleSpec = &ocrSpec
	jb.OCROracleSpec.TransmitterAddress = &app.Keys[0].EIP55Address
	err = app.AddJobV2(ctx, &jb)
	require.NoError(t, err)
	dbJb, err := app.JobORM().FindJob(ctx, jb.ID)
	require.NoError(t, err)
	require.Equal(t, dbJb.Name.String, ocrspec.Name)

	// test Calling update on the job id with changed values should succeed.
	updatedSpec := testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
		DS1BridgeName:      bridge2.Name.String(),
		DS2BridgeName:      bridge.Name.String(),
		Name:               "updated OCR job",
		TransmitterAddress: app.Keys[0].Address.Hex(),
	})
	require.NoError(t, err)
	body, _ := json.Marshal(web.UpdateJobRequest{
		TOML: updatedSpec.Toml(),
	})
	response, cleanup := client.Put("/v2/jobs/"+fmt.Sprintf("%v", jb.ID), bytes.NewReader(body))
	t.Cleanup(cleanup)

	dbJb, err = app.JobORM().FindJob(ctx, jb.ID)
	require.NoError(t, err)
	require.Equal(t, dbJb.Name.String, updatedSpec.Name)

	cltest.AssertServerResponse(t, response, http.StatusOK)
}

func TestJobsController_Update_NonExistentID(t *testing.T) {
	ctx := testutils.Context(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.OCR.Enabled = ptr(true)
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", freeport.GetOne(t))}
		c.P2P.PeerID = &cltest.DefaultP2PPeerID
	})
	app := cltest.NewApplicationWithConfigAndKey(t, cfg, cltest.DefaultP2PKey)

	require.NoError(t, app.KeyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	require.NoError(t, app.Start(ctx))

	_, bridge := cltest.MustCreateBridge(t, app.GetDB(), cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, app.GetDB(), cltest.BridgeOpts{})

	client := app.NewHTTPClient(nil)

	var jb job.Job
	ocrspec := testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
		DS1BridgeName: bridge.Name.String(),
		DS2BridgeName: bridge2.Name.String(),
		Name:          "old OCR job",
	})
	err := toml.Unmarshal([]byte(ocrspec.Toml()), &jb)
	require.NoError(t, err)
	var ocrSpec job.OCROracleSpec
	err = toml.Unmarshal([]byte(ocrspec.Toml()), &ocrSpec)
	require.NoError(t, err)
	jb.OCROracleSpec = &ocrSpec
	jb.OCROracleSpec.TransmitterAddress = &app.Keys[0].EIP55Address
	err = app.AddJobV2(ctx, &jb)
	require.NoError(t, err)

	// test Calling update on the job id with changed values should succeed.
	updatedSpec := testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
		DS1BridgeName:      bridge2.Name.String(),
		DS2BridgeName:      bridge.Name.String(),
		Name:               "updated OCR job",
		TransmitterAddress: app.Keys[0].EIP55Address.String(),
	})
	require.NoError(t, err)
	body, _ := json.Marshal(web.UpdateJobRequest{
		TOML: updatedSpec.Toml(),
	})
	response, cleanup := client.Put("/v2/jobs/99999", bytes.NewReader(body))
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusNotFound)
}

func runOCRJobSpecAssertions(t *testing.T, ocrJobSpecFromFileDB job.Job, ocrJobSpecFromServer presenters.JobResource) {
	ocrJobSpecFromFile := ocrJobSpecFromFileDB.OCROracleSpec
	assert.Equal(t, ocrJobSpecFromFile.ContractAddress, ocrJobSpecFromServer.OffChainReportingSpec.ContractAddress)
	assert.Equal(t, ocrJobSpecFromFile.P2PV2Bootstrappers, ocrJobSpecFromServer.OffChainReportingSpec.P2PV2Bootstrappers)
	assert.Equal(t, ocrJobSpecFromFile.IsBootstrapPeer, ocrJobSpecFromServer.OffChainReportingSpec.IsBootstrapPeer)
	assert.Equal(t, ocrJobSpecFromFile.EncryptedOCRKeyBundleID, ocrJobSpecFromServer.OffChainReportingSpec.EncryptedOCRKeyBundleID)
	assert.Equal(t, ocrJobSpecFromFile.TransmitterAddress, ocrJobSpecFromServer.OffChainReportingSpec.TransmitterAddress)
	assert.Equal(t, ocrJobSpecFromFile.ObservationTimeout, ocrJobSpecFromServer.OffChainReportingSpec.ObservationTimeout)
	assert.Equal(t, ocrJobSpecFromFile.BlockchainTimeout, ocrJobSpecFromServer.OffChainReportingSpec.BlockchainTimeout)
	assert.Equal(t, ocrJobSpecFromFile.ContractConfigTrackerSubscribeInterval, ocrJobSpecFromServer.OffChainReportingSpec.ContractConfigTrackerSubscribeInterval)
	assert.Equal(t, ocrJobSpecFromFile.ContractConfigTrackerSubscribeInterval, ocrJobSpecFromServer.OffChainReportingSpec.ContractConfigTrackerSubscribeInterval)
	assert.Equal(t, ocrJobSpecFromFile.ContractConfigConfirmations, ocrJobSpecFromServer.OffChainReportingSpec.ContractConfigConfirmations)
	assert.Equal(t, ocrJobSpecFromFileDB.Pipeline.Source, ocrJobSpecFromServer.PipelineSpec.DotDAGSource)

	// Check that create and update dates are non empty values.
	// Empty date value is "0001-01-01 00:00:00 +0000 UTC" so we are checking for the
	// millennia and century characters to be present
	assert.Contains(t, ocrJobSpecFromServer.OffChainReportingSpec.CreatedAt.String(), "20")
	assert.Contains(t, ocrJobSpecFromServer.OffChainReportingSpec.UpdatedAt.String(), "20")
}

func runDirectRequestJobSpecAssertions(t *testing.T, ereJobSpecFromFile job.Job, ereJobSpecFromServer presenters.JobResource) {
	assert.Equal(t, ereJobSpecFromFile.DirectRequestSpec.ContractAddress, ereJobSpecFromServer.DirectRequestSpec.ContractAddress)
	assert.Equal(t, ereJobSpecFromFile.Pipeline.Source, ereJobSpecFromServer.PipelineSpec.DotDAGSource)
	// Check that create and update dates are non empty values.
	// Empty date value is "0001-01-01 00:00:00 +0000 UTC" so we are checking for the
	// millennia and century characters to be present
	assert.Contains(t, ereJobSpecFromServer.DirectRequestSpec.CreatedAt.String(), "20")
	assert.Contains(t, ereJobSpecFromServer.DirectRequestSpec.UpdatedAt.String(), "20")
}

func setupBridges(t *testing.T, ds sqlutil.DataSource) (b1, b2 string) {
	_, bridge := cltest.MustCreateBridge(t, ds, cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, ds, cltest.BridgeOpts{})
	return bridge.Name.String(), bridge2.Name.String()
}

func setupJobsControllerTests(t *testing.T) (ta *cltest.TestApplication, cc cltest.HTTPClientCleaner) {
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.OCR.Enabled = ptr(true)
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", freeport.GetOne(t))}
		c.P2P.PeerID = &cltest.DefaultP2PPeerID
	})
	ec := setupEthClientForControllerTests(t)
	app := cltest.NewApplicationWithConfigAndKey(t, cfg, cltest.DefaultP2PKey, ec)
	ctx := testutils.Context(t)
	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	vrfKeyStore := app.GetKeyStore().VRF()
	_, err := vrfKeyStore.Create(ctx)
	require.NoError(t, err)
	return app, client
}

func setupEthClientForControllerTests(t *testing.T) *evmclimocks.Client {
	ec := cltest.NewEthMocksWithStartupAssertions(t)
	ec.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil).Maybe()
	ec.On("LatestBlockHeight", mock.Anything).Return(big.NewInt(100), nil).Maybe()
	ec.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Once().Return(big.NewInt(0), nil).Maybe()
	return ec
}

func setupJobSpecsControllerTestsWithJobs(t *testing.T) (*cltest.TestApplication, cltest.HTTPClientCleaner, job.Job, int32, job.Job, int32) {
	ctx := testutils.Context(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.OCR.Enabled = ptr(true)
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", freeport.GetOne(t))}
		c.P2P.PeerID = &cltest.DefaultP2PPeerID
	})
	app := cltest.NewApplicationWithConfigAndKey(t, cfg, cltest.DefaultP2PKey)

	require.NoError(t, app.KeyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	require.NoError(t, app.Start(ctx))

	_, bridge := cltest.MustCreateBridge(t, app.GetDB(), cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, app.GetDB(), cltest.BridgeOpts{})

	client := app.NewHTTPClient(nil)

	var jb job.Job
	ocrspec := testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{DS1BridgeName: bridge.Name.String(), DS2BridgeName: bridge2.Name.String(), EVMChainID: testutils.FixtureChainID.String()})
	err := toml.Unmarshal([]byte(ocrspec.Toml()), &jb)
	require.NoError(t, err)
	var ocrSpec job.OCROracleSpec
	err = toml.Unmarshal([]byte(ocrspec.Toml()), &ocrSpec)
	require.NoError(t, err)
	jb.OCROracleSpec = &ocrSpec
	jb.OCROracleSpec.TransmitterAddress = &app.Keys[0].EIP55Address
	err = app.AddJobV2(ctx, &jb)
	require.NoError(t, err)

	drSpec := fmt.Sprintf(`
		type                = "directrequest"
		schemaVersion       = 1
		evmChainID          = "0"
		name                = "example eth request event spec"
		contractAddress     = "0x613a38AC1659769640aaE063C651F48E0250454C"
		externalJobID       = "%s"
		observationSource   = """
		    ds1          [type=http method=GET url="http://example.com" allowunrestrictednetworkaccess="true"];
		    ds1_merge    [type=merge left="{}"]
		    ds1_parse    [type=jsonparse path="USD"];
		    ds1_multiply [type=multiply times=100];
		    ds1 -> ds1_parse -> ds1_multiply;
		"""
		`, uuid.New())

	erejb, err := directrequest.ValidatedDirectRequestSpec(drSpec)
	require.NoError(t, err)
	err = app.AddJobV2(ctx, &erejb)
	require.NoError(t, err)

	return app, client, jb, jb.ID, erejb, erejb.ID
}
