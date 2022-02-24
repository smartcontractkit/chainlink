package presenters_test

import (
	"fmt"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/lib/pq"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func TestJob(t *testing.T) {
	// Used in multiple tests
	timestamp := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	contractAddress, err := ethkey.NewEIP55Address("0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba")
	require.NoError(t, err)
	cronSchedule := "0 0 0 1 1 *"
	evmChainID := utils.NewBigI(42)
	fromAddress, err := ethkey.NewEIP55Address("0xa8037A20989AFcBC51798de9762b351D63ff462e")
	require.NoError(t, err)

	// Used in OCR tests
	var ocrKeyBundleID = "f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5"
	ocrKeyID := models.MustSha256HashFromHex(ocrKeyBundleID)
	transmitterAddress, err := ethkey.NewEIP55Address("0x27548a32b9aD5D64c5945EaE9Da5337bc3169D15")
	require.NoError(t, err)

	// Used in blockhashstore test
	v1CoordAddress, err := ethkey.NewEIP55Address("0x16988483b46e695f6c8D58e6e1461DC703e008e1")
	require.NoError(t, err)

	v2CoordAddress, err := ethkey.NewEIP55Address("0x2C409DD6D4eBDdA190B5174Cc19616DD13884262")
	require.NoError(t, err)

	testCases := []struct {
		name string
		job  job.Job
		want string
	}{
		{
			name: "direct request spec",
			job: job.Job{
				ID: 1,
				DirectRequestSpec: &job.DirectRequestSpec{
					ContractAddress: contractAddress,
					CreatedAt:       timestamp,
					UpdatedAt:       timestamp,
					EVMChainID:      evmChainID,
				},
				ExternalJobID: uuid.FromStringOrNil("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"),
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: `ds1 [type=http method=GET url="https://pricesource1.com"`,
				},
				Type:            job.Type("directrequest"),
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
						"keeperSpec": null,
                        "cronSpec": null,
                        "vrfSpec": null,
						"webhookSpec": null,
						"blockhashStoreSpec": null,
						"bootstrapSpec": null,
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
				ExternalJobID: uuid.FromStringOrNil("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"),
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: `ds1 [type=http method=GET url="https://pricesource1.com"`,
				},
				Type:            job.Type("fluxmonitor"),
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
						"offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
						"directRequestSpec": null,
						"keeperSpec": null,
                        "cronSpec": null,
                        "vrfSpec": null,
						"webhookSpec": null,
						"blockhashStoreSpec": null,
						"bootstrapSpec": null,
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
					P2PBootstrapPeers:                      pq.StringArray{"/dns4/chain.link/tcp/1234/p2p/xxx"},
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
				ExternalJobID: uuid.FromStringOrNil("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"),
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: `ds1 [type=http method=GET url="https://pricesource1.com"`,
				},
				Type:            job.Type("offchainreporting"),
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
							"p2pBootstrapPeers": ["/dns4/chain.link/tcp/1234/p2p/xxx"],
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
						"directRequestSpec": null,
						"keeperSpec": null,
                        "cronSpec": null,
                        "vrfSpec": null,
						"webhookSpec": null,
						"blockhashStoreSpec": null,
						"bootstrapSpec": null,
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
				ExternalJobID: uuid.FromStringOrNil("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"),
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: "",
				},
				Type:            job.Type("keeper"),
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
						"directRequestSpec": null,
						"cronSpec": null,
						"webhookSpec": null,
						"offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
                        "cronSpec": null,
                        "vrfSpec": null,
						"blockhashStoreSpec": null,
						"bootstrapSpec": null,
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
				ExternalJobID: uuid.FromStringOrNil("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"),
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: "",
				},
				Type:            job.Type("cron"),
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
                        "directRequestSpec": null,
                        "keeperSpec": null,
                        "offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
						"vrfSpec": null,
                        "webhookSpec": null,
						"blockhashStoreSpec": null,
						"bootstrapSpec": null,
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
				ExternalJobID: uuid.FromStringOrNil("0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"),
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: "",
				},
				Type:            job.Type("webhook"),
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
						"fluxMonitorSpec": null,
						"directRequestSpec": null,
						"keeperSpec": null,
						"cronSpec": null,
						"offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
                        "vrfSpec": null,
						"blockhashStoreSpec": null,
						"bootstrapSpec": null,
						"errors": []
					}
				}
			}`,
		},
		{
			name: "blockhash store spec",
			job: job.Job{
				ID: 1,
				BlockhashStoreSpec: &job.BlockhashStoreSpec{
					ID:                    1,
					CoordinatorV1Address:  &v1CoordAddress,
					CoordinatorV2Address:  &v2CoordAddress,
					WaitBlocks:            123,
					LookbackBlocks:        223,
					BlockhashStoreAddress: contractAddress,
					PollPeriod:            25 * time.Second,
					RunTimeout:            10 * time.Second,
					EVMChainID:            utils.NewBigI(4),
					FromAddress:           &fromAddress,
				},
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: "",
				},
				ExternalJobID: uuid.FromStringOrNil("0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"),
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
						"cronSpec": null,
						"offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
						"keeperSpec": null,
						"vrfSpec": null,
						"webhookSpec": null,
						"blockhashStoreSpec": {
							"coordinatorV1Address": "0x16988483b46e695f6c8D58e6e1461DC703e008e1",
							"coordinatorV2Address": "0x2C409DD6D4eBDdA190B5174Cc19616DD13884262",
							"waitBlocks": 123,
							"lookbackBlocks": 223,
							"blockhashStoreAddress": "0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba",
							"pollPeriod": 25000000000,
							"runTimeout": 10000000000,
							"evmChainID": "4",
							"fromAddress": "0xa8037A20989AFcBC51798de9762b351D63ff462e",
							"createdAt": "0001-01-01T00:00:00Z",
							"updatedAt": "0001-01-01T00:00:00Z"
						},
						"bootstrapSpec": null,
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
				ExternalJobID: uuid.FromStringOrNil("0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"),
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
						"cronSpec": null,
						"offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
						"keeperSpec": null,
						"vrfSpec": null,
						"webhookSpec": null,
						"blockhashStoreSpec": null,
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
				ExternalJobID: uuid.FromStringOrNil("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"),
				PipelineSpec: &pipeline.Spec{
					ID:           1,
					DotDagSource: "",
				},
				Type:            job.Type("keeper"),
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
						"directRequestSpec": null,
						"cronSpec": null,
						"webhookSpec": null,
						"offChainReportingOracleSpec": null,
						"offChainReporting2OracleSpec": null,
						"vrfSpec": null,
						"blockhashStoreSpec": null,
						"bootstrapSpec": null,
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
