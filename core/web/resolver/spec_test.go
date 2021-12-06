package resolver

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Specs are only embedded on the job and are not fetchable by it's own id, so
// we test the spec resolvers by fetching a job by id.

func TestResolver_CronSpec(t *testing.T) {
	var (
		id = int32(1)
	)

	testCases := []GQLTestCase{
		{
			name:          "cron spec success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobTx", id).Return(job.Job{
					Type: job.Cron,
					CronSpec: &job.CronSpec{
						CronSchedule: "CRON_TZ=UTC 0 0 1 1 *",
						CreatedAt:    f.Timestamp(),
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on CronSpec {
									schedule
									createdAt
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "CronSpec",
							"schedule": "CRON_TZ=UTC 0 0 1 1 *",
							"createdAt": "2021-01-01T00:00:00Z"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_DirectRequestSpec(t *testing.T) {
	var (
		id               = int32(1)
		requesterAddress = common.HexToAddress("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	)
	contractAddress, err := ethkey.NewEIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C")
	require.NoError(t, err)

	testCases := []GQLTestCase{
		{
			name:          "direct request spec success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobTx", id).Return(job.Job{
					Type: job.DirectRequest,
					DirectRequestSpec: &job.DirectRequestSpec{
						ContractAddress:             contractAddress,
						CreatedAt:                   f.Timestamp(),
						EVMChainID:                  utils.NewBigI(42),
						MinIncomingConfirmations:    clnull.NewUint32(1, true),
						MinIncomingConfirmationsEnv: true,
						MinContractPayment:          assets.NewLinkFromJuels(1000),
						Requesters:                  models.AddressCollection{requesterAddress},
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on DirectRequestSpec {
									contractAddress
									createdAt
									evmChainID
									minIncomingConfirmations
									minIncomingConfirmationsEnv
									minContractPayment
									requesters
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "DirectRequestSpec",
							"contractAddress": "0x613a38AC1659769640aaE063C651F48E0250454C",
							"createdAt": "2021-01-01T00:00:00Z",
							"evmChainID": "42",
							"minIncomingConfirmations": 1,
							"minIncomingConfirmationsEnv": true,
							"minContractPayment": "1000",
							"requesters": ["0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"]
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_FluxMonitorSpec(t *testing.T) {
	var (
		id = int32(1)
	)
	contractAddress, err := ethkey.NewEIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C")
	require.NoError(t, err)

	testCases := []GQLTestCase{
		{
			name:          "flux monitor spec with standard timers",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobTx", id).Return(job.Job{
					Type: job.FluxMonitor,
					FluxMonitorSpec: &job.FluxMonitorSpec{
						ContractAddress:   contractAddress,
						CreatedAt:         f.Timestamp(),
						EVMChainID:        utils.NewBigI(42),
						DrumbeatEnabled:   false,
						IdleTimerDisabled: false,
						IdleTimerPeriod:   time.Duration(1 * time.Hour),
						MinPayment:        assets.NewLinkFromJuels(1000),
						PollTimerDisabled: false,
						PollTimerPeriod:   time.Duration(1 * time.Minute),
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on FluxMonitorSpec {
									absoluteThreshold
									contractAddress
									createdAt
									drumbeatEnabled
									drumbeatRandomDelay
									drumbeatSchedule
									evmChainID
									idleTimerDisabled
									idleTimerPeriod
									minPayment
									pollTimerDisabled
									pollTimerPeriod
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "FluxMonitorSpec",
							"absoluteThreshold": 0,
							"contractAddress": "0x613a38AC1659769640aaE063C651F48E0250454C",
							"createdAt": "2021-01-01T00:00:00Z",
							"drumbeatEnabled": false,
							"drumbeatRandomDelay": null,
							"drumbeatSchedule": null,
							"evmChainID": "42",
							"idleTimerDisabled": false,
							"idleTimerPeriod": "1h0m0s",
							"minPayment": "1000",
							"pollTimerDisabled": false,
							"pollTimerPeriod": "1m0s"
						}
					}
				}
			`,
		},
		{
			name:          "flux monitor spec with drumbeat",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobTx", id).Return(job.Job{
					Type: job.FluxMonitor,
					FluxMonitorSpec: &job.FluxMonitorSpec{
						ContractAddress:     contractAddress,
						CreatedAt:           f.Timestamp(),
						EVMChainID:          utils.NewBigI(42),
						DrumbeatEnabled:     true,
						DrumbeatRandomDelay: time.Duration(1 * time.Second),
						DrumbeatSchedule:    "CRON_TZ=UTC 0 0 1 1 *",
						IdleTimerDisabled:   true,
						IdleTimerPeriod:     time.Duration(1 * time.Hour),
						MinPayment:          assets.NewLinkFromJuels(1000),
						PollTimerDisabled:   true,
						PollTimerPeriod:     time.Duration(1 * time.Minute),
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on FluxMonitorSpec {
									absoluteThreshold
									contractAddress
									createdAt
									drumbeatEnabled
									drumbeatRandomDelay
									drumbeatSchedule
									evmChainID
									idleTimerDisabled
									idleTimerPeriod
									minPayment
									pollTimerDisabled
									pollTimerPeriod
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "FluxMonitorSpec",
							"absoluteThreshold": 0,
							"contractAddress": "0x613a38AC1659769640aaE063C651F48E0250454C",
							"createdAt": "2021-01-01T00:00:00Z",
							"drumbeatEnabled": true,
							"drumbeatRandomDelay": "1s",
							"drumbeatSchedule": "CRON_TZ=UTC 0 0 1 1 *",
							"evmChainID": "42",
							"idleTimerDisabled": true,
							"idleTimerPeriod": "1h0m0s",
							"minPayment": "1000",
							"pollTimerDisabled": true,
							"pollTimerPeriod": "1m0s"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_KeeperSpec(t *testing.T) {
	var (
		id          = int32(1)
		fromAddress = common.HexToAddress("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	)
	contractAddress, err := ethkey.NewEIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C")
	require.NoError(t, err)

	testCases := []GQLTestCase{
		{
			name:          "keeper spec",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobTx", id).Return(job.Job{
					Type: job.Keeper,
					KeeperSpec: &job.KeeperSpec{
						ContractAddress: contractAddress,
						CreatedAt:       f.Timestamp(),
						EVMChainID:      utils.NewBigI(42),
						FromAddress:     ethkey.EIP55AddressFromAddress(fromAddress),
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on KeeperSpec {
									contractAddress
									createdAt
									evmChainID
									fromAddress
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "KeeperSpec",
							"contractAddress": "0x613a38AC1659769640aaE063C651F48E0250454C",
							"createdAt": "2021-01-01T00:00:00Z",
							"evmChainID": "42",
							"fromAddress": "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_OCRSpec(t *testing.T) {
	var (
		id = int32(1)
	)
	contractAddress, err := ethkey.NewEIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C")
	require.NoError(t, err)

	transmitterAddress, err := ethkey.NewEIP55Address("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	require.NoError(t, err)

	keyBundleID := models.MustSha256HashFromHex("f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5")

	p2pPeerID, err := p2pkey.MakePeerID("12D3KooWL3XJ9EMCyZvmmGXL2LMiVBtrVa2BuESsJiXkSj7333Jw")
	require.NoError(t, err)

	testCases := []GQLTestCase{
		{
			name:          "OCR spec",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobTx", id).Return(job.Job{
					Type: job.OffchainReporting,
					OffchainreportingOracleSpec: &job.OffchainReportingOracleSpec{
						BlockchainTimeout:                         models.Interval(1 * time.Minute),
						BlockchainTimeoutEnv:                      false,
						ContractAddress:                           contractAddress,
						ContractConfigConfirmations:               1,
						ContractConfigConfirmationsEnv:            true,
						ContractConfigTrackerPollInterval:         models.Interval(1 * time.Minute),
						ContractConfigTrackerPollIntervalEnv:      false,
						ContractConfigTrackerSubscribeInterval:    models.Interval(2 * time.Minute),
						ContractConfigTrackerSubscribeIntervalEnv: true,
						DatabaseTimeout:                           models.Interval(3 * time.Second),
						DatabaseTimeoutEnv:                        true,
						ObservationGracePeriod:                    models.Interval(4 * time.Second),
						ObservationGracePeriodEnv:                 true,
						ContractTransmitterTransmitTimeout:        models.Interval(555 * time.Millisecond),
						ContractTransmitterTransmitTimeoutEnv:     true,
						CreatedAt:                                 f.Timestamp(),
						EVMChainID:                                utils.NewBigI(42),
						IsBootstrapPeer:                           false,
						EncryptedOCRKeyBundleID:                   &keyBundleID,
						ObservationTimeout:                        models.Interval(2 * time.Minute),
						ObservationTimeoutEnv:                     false,
						P2PPeerID:                                 p2pPeerID,
						P2PPeerIDEnv:                              true,
						P2PBootstrapPeers:                         pq.StringArray{"/dns4/test.com/tcp/2001/p2pkey"},
						TransmitterAddress:                        &transmitterAddress,
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on OCRSpec {
									blockchainTimeout
									blockchainTimeoutEnv
									contractAddress
									contractConfigConfirmations
									contractConfigConfirmationsEnv
									contractConfigTrackerPollInterval
									contractConfigTrackerPollIntervalEnv
									contractConfigTrackerSubscribeInterval
									contractConfigTrackerSubscribeIntervalEnv
									databaseTimeout
									databaseTimeoutEnv
									observationGracePeriod
									observationGracePeriodEnv
									contractTransmitterTransmitTimeout
									contractTransmitterTransmitTimeoutEnv
									createdAt
									evmChainID
									isBootstrapPeer
									keyBundleID
									observationTimeout
									observationTimeoutEnv
									p2pPeerID
									p2pPeerIDEnv
									p2pBootstrapPeers
									transmitterAddress
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "OCRSpec",
							"blockchainTimeout": "1m0s",
							"blockchainTimeoutEnv": false,
							"contractAddress": "0x613a38AC1659769640aaE063C651F48E0250454C",
							"contractConfigConfirmations": 1,
							"contractConfigConfirmationsEnv": true,
							"contractConfigTrackerPollInterval": "1m0s",
							"contractConfigTrackerPollIntervalEnv": false,
							"contractConfigTrackerSubscribeInterval": "2m0s",
							"contractConfigTrackerSubscribeIntervalEnv": true,
							"databaseTimeout": "3s",
							"databaseTimeoutEnv": true,
							"observationGracePeriod": "4s",
							"observationGracePeriodEnv": true,
							"contractTransmitterTransmitTimeout": "555ms",
							"contractTransmitterTransmitTimeoutEnv": true,
							"createdAt": "2021-01-01T00:00:00Z",
							"evmChainID": "42",
							"isBootstrapPeer": false,
							"keyBundleID": "f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5",
							"observationTimeout": "2m0s",
							"observationTimeoutEnv": false,
							"p2pPeerID": "p2p_12D3KooWL3XJ9EMCyZvmmGXL2LMiVBtrVa2BuESsJiXkSj7333Jw",
							"p2pPeerIDEnv": true,
							"p2pBootstrapPeers": ["/dns4/test.com/tcp/2001/p2pkey"],
							"transmitterAddress": "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_OCR2Spec(t *testing.T) {
	var (
		id = int32(1)
	)
	contractAddress, err := ethkey.NewEIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C")
	require.NoError(t, err)

	transmitterAddress, err := ethkey.NewEIP55Address("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	require.NoError(t, err)

	keyBundleID := models.MustSha256HashFromHex("f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5")

	p2pPeerID, err := p2pkey.MakePeerID("12D3KooWL3XJ9EMCyZvmmGXL2LMiVBtrVa2BuESsJiXkSj7333Jw")
	require.NoError(t, err)

	testCases := []GQLTestCase{
		{
			name:          "OCR 2 spec",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobTx", id).Return(job.Job{
					Type: job.OffchainReporting2,
					Offchainreporting2OracleSpec: &job.OffchainReporting2OracleSpec{
						BlockchainTimeout:                      models.Interval(1 * time.Minute),
						ContractAddress:                        contractAddress,
						ContractConfigConfirmations:            1,
						ContractConfigTrackerPollInterval:      models.Interval(1 * time.Minute),
						ContractConfigTrackerSubscribeInterval: models.Interval(1 * time.Minute),
						CreatedAt:                              f.Timestamp(),
						EVMChainID:                             utils.NewBigI(42),
						IsBootstrapPeer:                        false,
						JuelsPerFeeCoinPipeline:                "100000000",
						EncryptedOCRKeyBundleID:                null.StringFrom(keyBundleID.String()),
						MonitoringEndpoint:                     null.StringFrom("https://monitor.endpoint"),
						P2PPeerID:                              &p2pPeerID,
						P2PBootstrapPeers:                      pq.StringArray{"/dns4/test.com/tcp/2001/p2pkey"},
						TransmitterAddress:                     &transmitterAddress,
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on OCR2Spec {
									blockchainTimeout
									contractAddress
									contractConfigConfirmations
									contractConfigTrackerPollInterval
									contractConfigTrackerSubscribeInterval
									createdAt
									evmChainID
									isBootstrapPeer
									juelsPerFeeCoinSource
									keyBundleID
									monitoringEndpoint
									p2pPeerID
									p2pBootstrapPeers
									transmitterAddress
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "OCR2Spec",
							"blockchainTimeout": "1m0s",
							"contractAddress": "0x613a38AC1659769640aaE063C651F48E0250454C",
							"contractConfigConfirmations": 1,
							"contractConfigTrackerPollInterval": "1m0s",
							"contractConfigTrackerSubscribeInterval": "1m0s",
							"createdAt": "2021-01-01T00:00:00Z",
							"evmChainID": "42",
							"isBootstrapPeer": false,
							"juelsPerFeeCoinSource": "100000000",
							"keyBundleID": "f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5",
							"monitoringEndpoint": "https://monitor.endpoint",
							"p2pPeerID": "p2p_12D3KooWL3XJ9EMCyZvmmGXL2LMiVBtrVa2BuESsJiXkSj7333Jw",
							"p2pBootstrapPeers": ["/dns4/test.com/tcp/2001/p2pkey"],
							"transmitterAddress": "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_VRFSpec(t *testing.T) {
	var (
		id = int32(1)
	)
	coordinatorAddress, err := ethkey.NewEIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C")
	require.NoError(t, err)

	fromAddress, err := ethkey.NewEIP55Address("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	require.NoError(t, err)

	pubKey, err := secp256k1.NewPublicKeyFromHex("0x9dc09a0f898f3b5e8047204e7ce7e44b587920932f08431e29c9bf6923b8450a01")
	require.NoError(t, err)

	testCases := []GQLTestCase{
		{
			name:          "vrf spec",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobTx", id).Return(job.Job{
					Type: job.VRF,
					VRFSpec: &job.VRFSpec{
						MinIncomingConfirmations: 1,
						CoordinatorAddress:       coordinatorAddress,
						CreatedAt:                f.Timestamp(),
						EVMChainID:               utils.NewBigI(42),
						FromAddress:              &fromAddress,
						PollPeriod:               1 * time.Minute,
						PublicKey:                pubKey,
						RequestedConfsDelay:      10,
						RequestTimeout:           24 * time.Hour,
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on VRFSpec {
									coordinatorAddress
									createdAt
									evmChainID
									fromAddress
									minIncomingConfirmations
									pollPeriod
									publicKey
									requestedConfsDelay
									requestTimeout
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "VRFSpec",
							"coordinatorAddress": "0x613a38AC1659769640aaE063C651F48E0250454C",
							"createdAt": "2021-01-01T00:00:00Z",
							"evmChainID": "42",
							"fromAddress": "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42",
							"minIncomingConfirmations": 1,
							"pollPeriod": "1m0s",
							"publicKey": "0x9dc09a0f898f3b5e8047204e7ce7e44b587920932f08431e29c9bf6923b8450a01",
							"requestedConfsDelay": 10,
							"requestTimeout": "24h0m0s"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_WebhookSpec(t *testing.T) {
	var (
		id = int32(1)
	)

	testCases := []GQLTestCase{
		{
			name:          "webhook spec",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobTx", id).Return(job.Job{
					Type: job.Webhook,
					WebhookSpec: &job.WebhookSpec{
						CreatedAt: f.Timestamp(),
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on WebhookSpec {
									createdAt
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "WebhookSpec",
							"createdAt": "2021-01-01T00:00:00Z"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}
