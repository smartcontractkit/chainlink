package presenters_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/lib/pq"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	evmassets "github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	clnull "github.com/smartcontractkit/chainlink/v2/core/null"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestJob(t *testing.T) {
	// Used in multiple tests
	timestamp := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	contractAddress, err := types.NewEIP55Address("0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba")
	require.NoError(t, err)
	cronSchedule := "0 0 0 1 1 *"
	evmChainID := big.NewI(42)
	fromAddress, err := types.NewEIP55Address("0xa8037A20989AFcBC51798de9762b351D63ff462e")
	require.NoError(t, err)

	// Used in OCR tests
	var ocrKeyBundleID = "f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5"
	ocrKeyID := models.MustSha256HashFromHex(ocrKeyBundleID)
	transmitterAddress, err := types.NewEIP55Address("0x27548a32b9aD5D64c5945EaE9Da5337bc3169D15")
	require.NoError(t, err)

	// Used in blockhashstore test
	v1CoordAddress, err := types.NewEIP55Address("0x16988483b46e695f6c8D58e6e1461DC703e008e1")
	require.NoError(t, err)

	v2CoordAddress, err := types.NewEIP55Address("0x2C409DD6D4eBDdA190B5174Cc19616DD13884262")
	require.NoError(t, err)

	v2PlusCoordAddress, err := types.NewEIP55Address("0x92B5e28Ac583812874e4271380c7d070C5FB6E6b")
	require.NoError(t, err)

	// Used in blockheaderfeeder test
	batchBHSAddress, err := types.NewEIP55Address("0xF6bB415b033D19EFf24A872a4785c6e1C4426103")
	require.NoError(t, err)

	trustedBlockhashStoreAddress, err := types.NewEIP55Address("0x0ad9FE7a58216242a8475ca92F222b0640E26B63")
	require.NoError(t, err)
	trustedBlockhashStoreBatchSize := int32(20)

	var specGasLimit uint32 = 1000
	vrfPubKey, _ := secp256k1.NewPublicKeyFromHex("0xede539e216e3a50e69d1c68aa9cc472085876c4002f6e1e6afee0ea63b50a78b00")

	testCases := []struct {
		name string
		job  job.Job
		want string
	}{
		{
			name: "direct request spec",
			job: job.Job{
				ID:                1,
				GasLimit:          clnull.Uint32From(specGasLimit),
				ForwardingAllowed: false,
				DirectRequestSpec: &job.DirectRequestSpec{
					ContractAddress: contractAddress,
					CreatedAt:       timestamp,
					UpdatedAt:       timestamp,
					EVMChainID:      evmChainID,
				},
				ExternalJobID: uuid.MustParse("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"),
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: `ds1 [type=http method=GET url="https://pricesource1.com"`,
				},
				Type:            job.DirectRequest,
				SchemaVersion:   1,
				Name:            null.StringFrom("test"),
				MaxTaskDuration: models.Interval(1 * time.Minute),
			},
			want: fmt.Sprintf(`
			{
				"data":{
					"type":"jobs",
					"id":"1",
					"attributes":{
						"name": "test",
						"schemaVersion": 1,
						"type": "directrequest",
						"maxTaskDuration": "1m0s",
					    "externalJobID":"0eec7e1d-d0d2-476c-a1a8-72dfb6633f46",
						"pipelineSpec": {
							"id": 1,
							"dotDagSource": "ds1 [type=http method=GET url=\"https://pricesource1.com\"",
							"jobID": 0
						},
						"directRequestSpec": {
							"contractAddress": "%s",
							"minIncomingConfirmations": null,
							"minContractPaymentLinkJuels": null,
							"requesters": null,
							"initiator": "runlog",
							"createdAt":"2000-01-01T00:00:00Z",
							"updatedAt":"2000-01-01T00:00:00Z",
							"evmChainID": "42"
						},
						"offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
						"fluxMonitorSpec": null,
						"gasLimit": 1000,
						"forwardingAllowed": false,
						"keeperSpec": null,
                        "cronSpec": null,
                        "vrfSpec": null,
						"webhookSpec": null,
						"workflowSpec": null,
						"blockhashStoreSpec": null,
						"blockHeaderFeederSpec": null,
						"bootstrapSpec": null,
						"gatewaySpec": null,
						"errors": []
					}
				}
			}`, contractAddress),
		},
		{
			name: "fluxmonitor spec",
			job: job.Job{
				ID: 1,
				FluxMonitorSpec: &job.FluxMonitorSpec{
					ContractAddress:   contractAddress,
					Threshold:         0.5,
					IdleTimerPeriod:   1 * time.Minute,
					IdleTimerDisabled: false,
					PollTimerPeriod:   1 * time.Second,
					PollTimerDisabled: false,
					MinPayment:        assets.NewLinkFromJuels(1),
					CreatedAt:         timestamp,
					UpdatedAt:         timestamp,
					EVMChainID:        evmChainID,
				},
				ExternalJobID: uuid.MustParse("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"),
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: `ds1 [type=http method=GET url="https://pricesource1.com"`,
				},
				Type:            job.FluxMonitor,
				SchemaVersion:   1,
				Name:            null.StringFrom("test"),
				MaxTaskDuration: models.Interval(1 * time.Minute),
			},
			want: fmt.Sprintf(`
			{
				"data":{
					"type":"jobs",
					"id":"1",
					"attributes":{
						"name": "test",
						"schemaVersion": 1,
						"type": "fluxmonitor",
						"maxTaskDuration": "1m0s",
					    "externalJobID":"0eec7e1d-d0d2-476c-a1a8-72dfb6633f46",
						"pipelineSpec": {
							"id": 1,
							"dotDagSource": "ds1 [type=http method=GET url=\"https://pricesource1.com\"",
							"jobID": 0
						},
						"fluxMonitorSpec": {
							"contractAddress": "%s",
							"threshold": 0.5,
							"absoluteThreshold": 0,
							"idleTimerPeriod": "1m0s",
							"idleTimerDisabled": false,
							"pollTimerPeriod": "1s",
							"pollTimerDisabled": false,
              				"drumbeatEnabled": false,
              				"drumbeatRandomDelay": null,
              				"drumbeatSchedule": null,
							"minPayment": "1",
							"createdAt":"2000-01-01T00:00:00Z",
							"updatedAt":"2000-01-01T00:00:00Z",
							"evmChainID": "42"
						},
						"gasLimit": null,
						"forwardingAllowed": false,
						"offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
						"directRequestSpec": null,
						"keeperSpec": null,
                        "cronSpec": null,
                        "vrfSpec": null,
						"webhookSpec": null,
						"workflowSpec": null,
						"blockhashStoreSpec": null,
						"blockHeaderFeederSpec": null,
						"bootstrapSpec": null,
						"gatewaySpec": null,
						"errors": []
					}
				}
			}`, contractAddress),
		},
		{
			name: "ocr spec",
			job: job.Job{
				ID: 1,
				OCROracleSpec: &job.OCROracleSpec{
					ContractAddress:                        contractAddress,
					P2PV2Bootstrappers:                     pq.StringArray{"xxx:5001"},
					IsBootstrapPeer:                        true,
					EncryptedOCRKeyBundleID:                &ocrKeyID,
					TransmitterAddress:                     &transmitterAddress,
					ObservationTimeout:                     models.Interval(1 * time.Minute),
					BlockchainTimeout:                      models.Interval(1 * time.Minute),
					ContractConfigTrackerSubscribeInterval: models.Interval(1 * time.Minute),
					ContractConfigTrackerPollInterval:      models.Interval(1 * time.Minute),
					ContractConfigConfirmations:            1,
					CreatedAt:                              timestamp,
					UpdatedAt:                              timestamp,
					EVMChainID:                             evmChainID,
					DatabaseTimeout:                        models.NewInterval(2 * time.Second),
					ObservationGracePeriod:                 models.NewInterval(3 * time.Second),
					ContractTransmitterTransmitTimeout:     models.NewInterval(444 * time.Millisecond),
				},
				ExternalJobID: uuid.MustParse("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"),
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: `ds1 [type=http method=GET url="https://pricesource1.com"`,
				},
				Type:              job.OffchainReporting,
				SchemaVersion:     1,
				Name:              null.StringFrom("test"),
				GasLimit:          clnull.Uint32From(123),
				ForwardingAllowed: true,
				MaxTaskDuration:   models.Interval(1 * time.Minute),
			},
			want: fmt.Sprintf(`
			{
				"data":{
					"type":"jobs",
					"id":"1",
					"attributes":{
						"name": "test",
						"schemaVersion": 1,
						"type": "offchainreporting",
						"maxTaskDuration": "1m0s",
					  "externalJobID":"0eec7e1d-d0d2-476c-a1a8-72dfb6633f46",
						"pipelineSpec": {
							"id": 1,
							"dotDagSource": "ds1 [type=http method=GET url=\"https://pricesource1.com\"",
							"jobID": 0
						},
						"offChainReportingOracleSpec": {
							"contractAddress": "%s",
							"p2pv2Bootstrappers": ["xxx:5001"],
							"isBootstrapPeer": true,
							"keyBundleID": "%s",
							"transmitterAddress": "%s",
							"observationTimeout": "1m0s",
							"blockchainTimeout": "1m0s",
							"contractConfigTrackerSubscribeInterval": "1m0s",
							"contractConfigTrackerPollInterval": "1m0s",
							"contractConfigConfirmations": 1,
							"createdAt":"2000-01-01T00:00:00Z",
							"updatedAt":"2000-01-01T00:00:00Z",
							"evmChainID": "42",
							"databaseTimeout": "2s",
							"observationGracePeriod": "3s",
							"contractTransmitterTransmitTimeout": "444ms"
						},
						"offChainReporting2OracleSpec": null,
						"fluxMonitorSpec": null,
						"gasLimit": 123,
						"forwardingAllowed": true,
						"directRequestSpec": null,
						"keeperSpec": null,
                        "cronSpec": null,
                        "vrfSpec": null,
						"webhookSpec": null,
						"workflowSpec": null,
						"blockhashStoreSpec": null,
						"blockHeaderFeederSpec": null,
						"bootstrapSpec": null,
						"gatewaySpec": null,
						"errors": []
					}
				}
			}`, contractAddress, ocrKeyBundleID, transmitterAddress),
		},
		{
			name: "keeper spec",
			job: job.Job{
				ID: 1,
				KeeperSpec: &job.KeeperSpec{
					ContractAddress: contractAddress,
					FromAddress:     fromAddress,
					CreatedAt:       timestamp,
					UpdatedAt:       timestamp,
					EVMChainID:      evmChainID,
				},
				ExternalJobID: uuid.MustParse("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"),
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: "",
				},
				Type:            job.Keeper,
				SchemaVersion:   1,
				Name:            null.StringFrom("test"),
				MaxTaskDuration: models.Interval(1 * time.Minute),
			},
			want: fmt.Sprintf(`
			{
				"data":{
					"type":"jobs",
					"id":"1",
					"attributes":{
						"name": "test",
						"schemaVersion": 1,
						"type": "keeper",
						"maxTaskDuration": "1m0s",
					    "externalJobID":"0eec7e1d-d0d2-476c-a1a8-72dfb6633f46",
						"pipelineSpec": {
							"id": 1,
							"dotDagSource": "",
							"jobID": 0
						},
						"keeperSpec": {
							"contractAddress": "%s",
							"fromAddress": "%s",
							"createdAt":"2000-01-01T00:00:00Z",
							"updatedAt":"2000-01-01T00:00:00Z",
							"evmChainID": "42"
						},
						"fluxMonitorSpec": null,
						"gasLimit": null,
						"forwardingAllowed": false,
						"directRequestSpec": null,
						"cronSpec": null,
						"webhookSpec": null,
						"workflowSpec": null,
						"offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
                        "cronSpec": null,
                        "vrfSpec": null,
						"blockhashStoreSpec": null,
						"blockHeaderFeederSpec": null,
						"bootstrapSpec": null,
						"gatewaySpec": null,
						"errors": []
					}
				}
			}`, contractAddress, fromAddress),
		},
		{
			name: "cron spec",
			job: job.Job{
				ID: 1,
				CronSpec: &job.CronSpec{
					CronSchedule: cronSchedule,
					CreatedAt:    timestamp,
					UpdatedAt:    timestamp,
				},
				ExternalJobID: uuid.MustParse("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"),
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: "",
				},
				Type:            job.Cron,
				SchemaVersion:   1,
				Name:            null.StringFrom("test"),
				MaxTaskDuration: models.Interval(1 * time.Minute),
			},
			want: fmt.Sprintf(`
            {
                "data":{
                    "type":"jobs",
                    "id":"1",
                    "attributes":{
                        "name": "test",
                        "schemaVersion": 1,
                        "type": "cron",
                        "maxTaskDuration": "1m0s",
					    "externalJobID":"0eec7e1d-d0d2-476c-a1a8-72dfb6633f46",
                        "pipelineSpec": {
                            "id": 1,
                            "dotDagSource": "",
														"jobID": 0
                        },
                        "cronSpec": {
                            "schedule": "%s",
                            "createdAt":"2000-01-01T00:00:00Z",
                            "updatedAt":"2000-01-01T00:00:00Z"
                        },
                        "fluxMonitorSpec": null,
						"gasLimit": null,
						"forwardingAllowed": false,
                        "directRequestSpec": null,
                        "keeperSpec": null,
                        "offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
						"vrfSpec": null,
                        "webhookSpec": null,
						"workflowSpec": null,
						"blockhashStoreSpec": null,
						"blockHeaderFeederSpec": null,
						"bootstrapSpec": null,
						"gatewaySpec": null,
                        "errors": []
                    }
                }
            }`, cronSchedule),
		},
		{
			name: "webhook spec",
			job: job.Job{
				ID: 1,
				WebhookSpec: &job.WebhookSpec{
					CreatedAt: timestamp,
					UpdatedAt: timestamp,
				},
				ExternalJobID: uuid.MustParse("0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"),
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: "",
				},
				Type:            job.Webhook,
				SchemaVersion:   1,
				Name:            null.StringFrom("test"),
				MaxTaskDuration: models.Interval(1 * time.Minute),
			},
			want: `
			{
				"data":{
					"type":"jobs",
					"id":"1",
					"attributes":{
						"name": "test",
						"schemaVersion": 1,
						"type": "webhook",
						"maxTaskDuration": "1m0s",
					    "externalJobID":"0eec7e1d-d0d2-476c-a1a8-72dfb6633f46",
						"pipelineSpec": {
							"id": 1,
							"dotDagSource": "",
							"jobID": 0
						},
						"webhookSpec": {
							"createdAt":"2000-01-01T00:00:00Z",
							"updatedAt":"2000-01-01T00:00:00Z"
						},
						"workflowSpec": null,
						"fluxMonitorSpec": null,
						"gasLimit": null,
						"forwardingAllowed": false,
						"directRequestSpec": null,
						"keeperSpec": null,
						"cronSpec": null,
						"offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
                        "vrfSpec": null,
						"blockhashStoreSpec": null,
						"blockHeaderFeederSpec": null,
						"bootstrapSpec": null,
						"gatewaySpec": null,
						"errors": []
					}
				}
			}`,
		},
		{
			name: "vrf job spec",
			job: job.Job{
				ID:            1,
				Name:          null.StringFrom("vrf_test"),
				Type:          job.VRF,
				SchemaVersion: 1,
				ExternalJobID: uuid.MustParse("0eec7e1d-d0d2-476c-a1a8-72dfb6633f47"),
				VRFSpec: &job.VRFSpec{
					BatchCoordinatorAddress:       &contractAddress,
					BatchFulfillmentEnabled:       true,
					CustomRevertsPipelineEnabled:  true,
					MinIncomingConfirmations:      1,
					CoordinatorAddress:            contractAddress,
					CreatedAt:                     timestamp,
					UpdatedAt:                     timestamp,
					EVMChainID:                    evmChainID,
					FromAddresses:                 []types.EIP55Address{fromAddress},
					PublicKey:                     vrfPubKey,
					RequestedConfsDelay:           10,
					ChunkSize:                     25,
					BatchFulfillmentGasMultiplier: 1,
					GasLanePrice:                  evmassets.GWei(200),
					VRFOwnerAddress:               nil,
				},
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: "",
				},
			},
			want: fmt.Sprintf(`
			{
				"data": {
					"type": "jobs",
					"id": "1",
					"attributes": {
						"name": "vrf_test",
						"type": "vrf",
						"schemaVersion": 1,
						"maxTaskDuration": "0s",
						"externalJobID": "0eec7e1d-d0d2-476c-a1a8-72dfb6633f47",
						"directRequestSpec": null,
						"fluxMonitorSpec": null,
						"gasLimit": null,
						"forwardingAllowed": false,
						"cronSpec": null,
						"offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
						"keeperSpec": null,
						"vrfSpec": {
							"batchCoordinatorAddress": "%s",
							"batchFulfillmentEnabled": true,
							"customRevertsPipelineEnabled":  true,
							"confirmations":      1,
							"coordinatorAddress":            "%s",
							"createdAt":                     "2000-01-01T00:00:00Z",
							"updatedAt":                     "2000-01-01T00:00:00Z",
							"evmChainID":                    "42",
							"fromAddresses":                 ["%s"],
							"pollPeriod":                    "0s",
							"publicKey":                     "%s",
							"requestedConfsDelay":           10,
							"requestTimeout":                "0s",
							"chunkSize":                     25,
							"batchFulfillmentGasMultiplier": 1,
							"backoffInitialDelay":           "0s",
							"backoffMaxDelay":               "0s",
							"gasLanePrice":                  "200 gwei"
						},
						"webhookSpec": null,
						"workflowSpec": null,
						"blockhashStoreSpec": null,
						"blockHeaderFeederSpec": null,
						"bootstrapSpec": null,
						"pipelineSpec": {
							"id": 1,
							"jobID": 0,
							"dotDagSource": ""
						},
						"gatewaySpec": null,
						"errors": []
					}
				}
			}`, contractAddress, contractAddress, fromAddress, vrfPubKey.String()),
		},
		{
			name: "blockhash store spec",
			job: job.Job{
				ID: 1,
				BlockhashStoreSpec: &job.BlockhashStoreSpec{
					ID:                             1,
					CoordinatorV1Address:           &v1CoordAddress,
					CoordinatorV2Address:           &v2CoordAddress,
					CoordinatorV2PlusAddress:       &v2PlusCoordAddress,
					WaitBlocks:                     123,
					LookbackBlocks:                 223,
					HeartbeatPeriod:                375 * time.Second,
					BlockhashStoreAddress:          contractAddress,
					PollPeriod:                     25 * time.Second,
					RunTimeout:                     10 * time.Second,
					EVMChainID:                     big.NewI(4),
					FromAddresses:                  []types.EIP55Address{fromAddress},
					TrustedBlockhashStoreAddress:   &trustedBlockhashStoreAddress,
					TrustedBlockhashStoreBatchSize: trustedBlockhashStoreBatchSize,
				},
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: "",
				},
				ExternalJobID: uuid.MustParse("0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"),
				Type:          job.BlockhashStore,
				SchemaVersion: 1,
				Name:          null.StringFrom("test"),
			},
			want: `
			{
				"data": {
					"type": "jobs",
					"id": "1",
					"attributes": {
						"name": "test",
						"type": "blockhashstore",
						"schemaVersion": 1,
						"maxTaskDuration": "0s",
						"externalJobID": "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46",
						"directRequestSpec": null,
						"fluxMonitorSpec": null,
						"gasLimit": null,
						"forwardingAllowed": false,
						"cronSpec": null,
						"offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
						"keeperSpec": null,
						"vrfSpec": null,
						"webhookSpec": null,
						"workflowSpec": null,
						"blockhashStoreSpec": {
							"coordinatorV1Address": "0x16988483b46e695f6c8D58e6e1461DC703e008e1",
							"coordinatorV2Address": "0x2C409DD6D4eBDdA190B5174Cc19616DD13884262",
							"coordinatorV2PlusAddress": "0x92B5e28Ac583812874e4271380c7d070C5FB6E6b",
							"waitBlocks": 123,
							"lookbackBlocks": 223,
							"heartbeatPeriod": 375000000000,
							"blockhashStoreAddress": "0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba",
							"trustedBlockhashStoreAddress": "0x0ad9FE7a58216242a8475ca92F222b0640E26B63",
							"trustedBlockhashStoreBatchSize": 20,
							"pollPeriod": 25000000000,
							"runTimeout": 10000000000,
							"evmChainID": "4",
							"fromAddresses": ["0xa8037A20989AFcBC51798de9762b351D63ff462e"],
							"createdAt": "0001-01-01T00:00:00Z",
							"updatedAt": "0001-01-01T00:00:00Z"
						},
						"blockHeaderFeederSpec": null,
						"bootstrapSpec": null,
						"pipelineSpec": {
							"id": 1,
							"jobID": 0,
							"dotDagSource": ""
						},
						"gatewaySpec": null,
						"errors": []
					}
				}
			}`,
		},
		{
			name: "block header feeder spec",
			job: job.Job{
				ID: 1,
				BlockHeaderFeederSpec: &job.BlockHeaderFeederSpec{
					ID:                         1,
					CoordinatorV1Address:       &v1CoordAddress,
					CoordinatorV2Address:       &v2CoordAddress,
					CoordinatorV2PlusAddress:   &v2PlusCoordAddress,
					WaitBlocks:                 123,
					LookbackBlocks:             223,
					BlockhashStoreAddress:      contractAddress,
					BatchBlockhashStoreAddress: batchBHSAddress,
					PollPeriod:                 25 * time.Second,
					RunTimeout:                 10 * time.Second,
					EVMChainID:                 big.NewI(4),
					FromAddresses:              []types.EIP55Address{fromAddress},
					GetBlockhashesBatchSize:    5,
					StoreBlockhashesBatchSize:  10,
				},
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: "",
				},
				ExternalJobID: uuid.MustParse("0eec7e1d-d0d2-476c-a1a8-72dfb6633f47"),
				Type:          job.BlockHeaderFeeder,
				SchemaVersion: 1,
				Name:          null.StringFrom("blockheaderfeeder"),
			},
			want: `
			{
				"data": {
					"type": "jobs",
					"id": "1",
					"attributes": {
						"name": "blockheaderfeeder",
						"type": "blockheaderfeeder",
						"schemaVersion": 1,
						"maxTaskDuration": "0s",
						"externalJobID": "0eec7e1d-d0d2-476c-a1a8-72dfb6633f47",
						"directRequestSpec": null,
						"fluxMonitorSpec": null,
						"gasLimit": null,
						"forwardingAllowed": false,
						"cronSpec": null,
						"offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
						"keeperSpec": null,
						"vrfSpec": null,
						"webhookSpec": null,
						"workflowSpec": null,
						"blockhashStoreSpec": null,
						"blockHeaderFeederSpec": {
							"coordinatorV1Address": "0x16988483b46e695f6c8D58e6e1461DC703e008e1",
							"coordinatorV2Address": "0x2C409DD6D4eBDdA190B5174Cc19616DD13884262",
							"coordinatorV2PlusAddress": "0x92B5e28Ac583812874e4271380c7d070C5FB6E6b",
							"waitBlocks": 123,
							"lookbackBlocks": 223,
							"blockhashStoreAddress": "0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba",
							"batchBlockhashStoreAddress": "0xF6bB415b033D19EFf24A872a4785c6e1C4426103",
							"pollPeriod": 25000000000,
							"runTimeout": 10000000000,
							"evmChainID": "4",
							"fromAddresses": ["0xa8037A20989AFcBC51798de9762b351D63ff462e"],
							"getBlockhashesBatchSize": 5,
							"storeBlockhashesBatchSize": 10,
							"createdAt": "0001-01-01T00:00:00Z",
							"updatedAt": "0001-01-01T00:00:00Z"
						},
						"bootstrapSpec": null,
						"pipelineSpec": {
							"id": 1,
							"jobID": 0,
							"dotDagSource": ""
						},
						"gatewaySpec": null,
						"errors": []
					}
				}
			}`,
		},
		{
			name: "bootstrap spec",
			job: job.Job{
				ID: 1,
				BootstrapSpec: &job.BootstrapSpec{
					ID:          1,
					ContractID:  "0x16988483b46e695f6c8D58e6e1461DC703e008e1",
					Relay:       "evm",
					RelayConfig: map[string]interface{}{"chainID": 1337},
				},
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: "",
				},
				ExternalJobID: uuid.MustParse("0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"),
				Type:          job.Bootstrap,
				SchemaVersion: 1,
				Name:          null.StringFrom("test"),
			},
			want: `
			{
				"data": {
					"type": "jobs",
					"id": "1",
					"attributes": {
						"name": "test",
						"type": "bootstrap",
						"schemaVersion": 1,
						"maxTaskDuration": "0s",
						"externalJobID": "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46",
						"directRequestSpec": null,
						"fluxMonitorSpec": null,
						"gasLimit": null,
						"forwardingAllowed": false,
						"cronSpec": null,
						"offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
						"keeperSpec": null,
						"vrfSpec": null,
						"webhookSpec": null,
						"workflowSpec": null,
						"blockhashStoreSpec": null,
						"blockHeaderFeederSpec": null,
						"bootstrapSpec": {
							"blockchainTimeout":"0s", 
							"contractConfigConfirmations":0, 
							"contractConfigTrackerPollInterval":"0s", 
							"contractConfigTrackerSubscribeInterval":"0s", 
							"contractID":"0x16988483b46e695f6c8D58e6e1461DC703e008e1", 
							"createdAt":"0001-01-01T00:00:00Z", 
							"relay":"evm", 
							"relayConfig":{"chainID":1337}, 
							"updatedAt":"0001-01-01T00:00:00Z"
						},
						"pipelineSpec": {
							"id": 1,
							"jobID": 0,
							"dotDagSource": ""
						},
						"gatewaySpec": null,
						"errors": []
					}
				}
			}`,
		},
		{
			name: "gateway spec",
			job: job.Job{
				ID: 1,
				GatewaySpec: &job.GatewaySpec{
					ID: 3,
					GatewayConfig: map[string]interface{}{
						"NodeServerConfig": map[string]interface{}{},
					},
				},
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: "",
				},
				ExternalJobID: uuid.MustParse("0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"),
				Type:          job.Gateway,
				SchemaVersion: 1,
				Name:          null.StringFrom("gateway test"),
			},
			want: `
			{
				"data": {
					"type": "jobs",
					"id": "1",
					"attributes": {
						"name": "gateway test",
						"type": "gateway",
						"schemaVersion": 1,
						"maxTaskDuration": "0s",
						"externalJobID": "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46",
						"directRequestSpec": null,
						"fluxMonitorSpec": null,
						"gasLimit": null,
						"forwardingAllowed": false,
						"cronSpec": null,
						"offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
						"keeperSpec": null,
						"vrfSpec": null,
						"webhookSpec": null,
						"workflowSpec": null,
						"blockhashStoreSpec": null,
						"blockHeaderFeederSpec": null,
						"bootstrapSpec": null,
						"gatewaySpec": {
							"gatewayConfig": {
								"NodeServerConfig": {
								}
							},
							"createdAt":"0001-01-01T00:00:00Z",
							"updatedAt":"0001-01-01T00:00:00Z"
						},
						"pipelineSpec": {
							"id": 1,
							"jobID": 0,
							"dotDagSource": ""
						},
						"errors": []
					}
				}
			}`,
		},
		{
			name: "workflow spec",
			job: job.Job{
				ID: 1,
				WorkflowSpec: &job.WorkflowSpec{
					ID:            3,
					WorkflowID:    "<test-workflow-id>",
					Workflow:      `<test-workflow-spec>`,
					WorkflowOwner: "<test-workflow-owner>",
				},
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: "",
				},
				ExternalJobID: uuid.MustParse("0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"),
				Type:          job.Workflow,
				SchemaVersion: 1,
				Name:          null.StringFrom("workflow test"),
			},
			want: `
			{
				"data": {
					"type": "jobs",
					"id": "1",
					"attributes": {
						"name": "workflow test",
						"type": "workflow",
						"schemaVersion": 1,
						"maxTaskDuration": "0s",
						"externalJobID": "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46",
						"directRequestSpec": null,
						"fluxMonitorSpec": null,
						"gasLimit": null,
						"forwardingAllowed": false,
						"cronSpec": null,
						"offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
						"keeperSpec": null,
						"vrfSpec": null,
						"webhookSpec": null,
						"workflowSpec": {
							"workflow": "<test-workflow-spec>",
							"workflowId": "<test-workflow-id>",
							"workflowOwner": "<test-workflow-owner>",
							"createdAt":"0001-01-01T00:00:00Z",
							"updatedAt":"0001-01-01T00:00:00Z"
						},
						"blockhashStoreSpec": null,
						"blockHeaderFeederSpec": null,
						"bootstrapSpec": null,
						"gatewaySpec": null,
						"pipelineSpec": {
							"id": 1,
							"jobID": 0,
							"dotDagSource": ""
						},
						"errors": []
					}
				}
			}`,
		},
		{
			name: "with errors",
			job: job.Job{
				ID: 1,
				KeeperSpec: &job.KeeperSpec{
					ContractAddress: contractAddress,
					FromAddress:     fromAddress,
					CreatedAt:       timestamp,
					UpdatedAt:       timestamp,
					EVMChainID:      evmChainID,
				},
				ExternalJobID: uuid.MustParse("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"),
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: "",
				},
				Type:            job.Keeper,
				SchemaVersion:   1,
				Name:            null.StringFrom("test"),
				MaxTaskDuration: models.Interval(1 * time.Minute),
				JobSpecErrors: []job.SpecError{
					{
						ID:          200,
						JobID:       1,
						Description: "some error",
						Occurrences: 1,
						CreatedAt:   timestamp,
						UpdatedAt:   timestamp,
					},
				},
			},
			want: fmt.Sprintf(`
			{
				"data":{
					"type":"jobs",
					"id":"1",
					"attributes":{
						"name": "test",
						"schemaVersion": 1,
						"type": "keeper",
						"maxTaskDuration": "1m0s",
					    "externalJobID":"0eec7e1d-d0d2-476c-a1a8-72dfb6633f46",
						"pipelineSpec": {
							"id": 1,
							"dotDagSource": "",
							"jobID": 0
						},
						"keeperSpec": {
							"contractAddress": "%s",
							"fromAddress": "%s",
							"createdAt":"2000-01-01T00:00:00Z",
							"updatedAt":"2000-01-01T00:00:00Z",
							"evmChainID": "42"
						},
						"fluxMonitorSpec": null,
						"gasLimit": null,
						"forwardingAllowed": false,
						"directRequestSpec": null,
						"cronSpec": null,
						"webhookSpec": null,
						"workflowSpec": null,
						"offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
						"vrfSpec": null,
						"blockhashStoreSpec": null,
						"blockHeaderFeederSpec": null,
						"bootstrapSpec": null,
						"gatewaySpec": null,
						"errors": [{
							"id": 200,
							"description": "some error",
							"occurrences": 1,
							"createdAt":"2000-01-01T00:00:00Z",
							"updatedAt":"2000-01-01T00:00:00Z"
						}]
					}
				}
			}`, contractAddress, fromAddress),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := presenters.NewJobResource(tc.job)
			b, err := jsonapi.Marshal(r)
			require.NoError(t, err)

			assert.JSONEq(t, tc.want, string(b))
		})
	}
}
