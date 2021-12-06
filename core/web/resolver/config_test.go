package resolver

import (
	"testing"

	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
)

func TestResolver_Config(t *testing.T) {
	t.Parallel()

	query := `
		query GetConfiguration {
			config {
				items {
					key
					config
				}
			}
		}`

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "config"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				// Using the default config value for now just to validate that it works
				// Mocking this would require complying to the whole interface
				// Which means mocking each method here, which I'm not sure we would like to do
				cfg := configtest.NewTestGeneralConfigWithOverrides(t, configtest.GeneralConfigOverrides{
					AdminCredentialsFile: null.StringFrom("test"),
					AdvisoryLockID:       null.IntFrom(1),
					AllowOrigins:         null.StringFrom("test"),
					BlockBackfillDepth:   null.IntFrom(1),
					BlockBackfillSkip:    null.BoolFrom(false),
					ClientNodeURL:        null.StringFrom("test"),
					DatabaseURL:          null.StringFrom("test"),
					DefaultChainID:       nil,
					DefaultHTTPAllowUnrestrictedNetworkAccess: null.BoolFrom(true),
					DefaultHTTPTimeout:                        nil,
					DefaultMaxHTTPAttempts:                    null.IntFrom(1),
					Dev:                                       null.BoolFrom(true),
					Dialect:                                   "",
					EVMDisabled:                               null.BoolFrom(true),
					EthereumDisabled:                          null.BoolFrom(true),
					EthereumURL:                               null.StringFrom(""),
					FeatureExternalInitiators:                 null.BoolFrom(true),
					GlobalBalanceMonitorEnabled:               null.BoolFrom(true),
					GlobalChainType:                           null.StringFrom(""),
					GlobalEthTxReaperThreshold:                nil,
					GlobalEthTxResendAfterThreshold:           nil,
					GlobalEvmEIP1559DynamicFees:               null.BoolFrom(true),
					GlobalEvmFinalityDepth:                    null.IntFrom(1),
					GlobalEvmGasBumpPercent:                   null.IntFrom(1),
					GlobalEvmGasBumpTxDepth:                   null.IntFrom(1),
					GlobalEvmGasBumpWei:                       nil,
					GlobalEvmGasLimitDefault:                  null.IntFrom(1),
					GlobalEvmGasLimitMultiplier:               null.FloatFrom(1),
					GlobalEvmGasPriceDefault:                  nil,
					GlobalEvmGasTipCapDefault:                 nil,
					GlobalEvmGasTipCapMinimum:                 nil,
					GlobalEvmHeadTrackerHistoryDepth:          null.IntFrom(1),
					GlobalEvmHeadTrackerMaxBufferSize:         null.IntFrom(1),
					GlobalEvmHeadTrackerSamplingInterval:      nil,
					GlobalEvmLogBackfillBatchSize:             null.IntFrom(1),
					GlobalEvmMaxGasPriceWei:                   nil,
					GlobalEvmMinGasPriceWei:                   nil,
					GlobalEvmNonceAutoSync:                    null.BoolFrom(false),
					GlobalEvmRPCDefaultBatchSize:              null.IntFrom(1),
					GlobalFlagsContractAddress:                null.StringFrom("test"),
					GlobalGasEstimatorMode:                    null.StringFrom("test"),
					GlobalMinIncomingConfirmations:            null.IntFrom(1),
					GlobalMinRequiredOutgoingConfirmations:    null.IntFrom(1),
					GlobalMinimumContractPayment:              nil,
					KeeperMaximumGracePeriod:                  null.IntFrom(1),
					KeeperRegistrySyncInterval:                nil,
					KeeperRegistrySyncUpkeepQueueSize:         null.IntFrom(1),
					LogLevel:                                  nil,
					DefaultLogLevel:                           nil,
					LogSQL:                                    null.BoolFrom(true),
					LogToDisk:                                 null.BoolFrom(true),
					OCRKeyBundleID:                            null.StringFrom("test"),
					OCRObservationTimeout:                     nil,
					OCRTransmitterAddress:                     nil,
					P2PBootstrapPeers:                         nil,
					P2PListenPort:                             null.IntFrom(1),
					P2PPeerID:                                 "",
					P2PPeerIDError:                            nil,
					SecretGenerator:                           nil,
					TriggerFallbackDBPollInterval:             nil,
				})
				cfg.SetRootDir("/tmp/chainlink_test/gql-test")

				f.App.On("GetConfig").Return(cfg)
			},
			query: query,
			result: `{
				"config": {
					"items": [{
						"config": {
							"value": "test"
						},
						"key": "ALLOW_ORIGINS"
					}, {
						"config": {
							"value": "1"
						},
						"key": "BLOCK_BACKFILL_DEPTH"
					}, {
						"config": {
							"value": "0"
						},
						"key": "BLOCK_HISTORY_ESTIMATOR_BLOCK_DELAY"
					}, {
						"config": {
							"value": "0"
						},
						"key": "BLOCK_HISTORY_ESTIMATOR_BLOCK_HISTORY_SIZE"
					}, {
						"config": {
							"value": "0"
						},
						"key": "BLOCK_HISTORY_ESTIMATOR_TRANSACTION_PERCENTILE"
					}, {
						"config": {
							"value": "http://localhost:6688"
						},
						"key": "BRIDGE_RESPONSE_URL"
					}, {
						"config": {
							"value": ""
						},
						"key": "CHAIN_TYPE"
					}, {
						"config": {
							"value": "test"
						},
						"key": "CLIENT_NODE_URL"
					}, {
						"config": {
							"value": "1h0m0s"
						},
						"key": "DATABASE_BACKUP_FREQUENCY"
					}, {
						"config": {
							"value": "none"
						},
						"key": "DATABASE_BACKUP_MODE"
					}, {
						"config": {
							"value": "none"
						},
						"key": "DATABASE_LOCKING_MODE"
					}, {
						"config": {
							"value": "0"
						},
						"key": "ETH_CHAIN_ID"
					}, {
						"config": {
							"value": "32768"
						},
						"key": "DEFAULT_HTTP_LIMIT"
					}, {
						"config": {
							"value": "15s"
						},
						"key": "DEFAULT_HTTP_TIMEOUT"
					}, {
						"config": {
							"value": "true"
						},
						"key": "CHAINLINK_DEV"
					}, {
						"config": {
							"value": "true"
						},
						"key": "ETH_DISABLED"
					}, {
						"config": {
							"value": ""
						},
						"key": "ETH_HTTP_URL"
					}, {
						"config": {
							"value": "[]"
						},
						"key": "ETH_SECONDARY_URLS"
					}, {
						"config": {
							"value": ""
						},
						"key": "ETH_URL"
					}, {
						"config": {
							"value": ""
						},
						"key": "EXPLORER_URL"
					}, {
						"config": {
							"value": "1"
						},
						"key": "FM_DEFAULT_TRANSACTION_QUEUE_DEPTH"
					}, {
						"config": {
							"value": "true"
						},
						"key": "FEATURE_EXTERNAL_INITIATORS"
					}, {
						"config": {
							"value": "false"
						},
						"key": "FEATURE_OFFCHAIN_REPORTING"
					}, {
						"config": {
							"value": ""
						},
						"key": "GAS_ESTIMATOR_MODE"
					}, {
						"config": {
							"value": "true"
						},
						"key": "INSECURE_FAST_SCRYPT"
					}, {
						"config": {
							"value": "false"
						},
						"key": "JSON_CONSOLE"
					}, {
						"config": {
							"value": "1h0m0s"
						},
						"key": "JOB_PIPELINE_REAPER_INTERVAL"
					}, {
						"config": {
							"value": "24h0m0s"
						},
						"key": "JOB_PIPELINE_REAPER_THRESHOLD"
					}, {
						"config": {
							"value": "1"
						},
						"key": "KEEPER_DEFAULT_TRANSACTION_QUEUE_DEPTH"
					}, {
						"config": {
							"value": "20"
						},
						"key": "KEEPER_GAS_PRICE_BUFFER_PERCENT"
					}, {
						"config": {
							"value": "20"
						},
						"key": "KEEPER_GAS_TIP_CAP_BUFFER_PERCENT"
					}, {
						"config": {
							"value": "0"
						},
						"key": "KEEPER_MAXIMUM_GRACE_PERIOD"
					}, {
						"config": {
							"value": "0"
						},
						"key": "KEEPER_REGISTRY_CHECK_GAS_OVERHEAD"
					}, {
						"config": {
							"value": "0"
						},
						"key": "KEEPER_REGISTRY_PERFORM_GAS_OVERHEAD"
					}, {
						"config": {
							"value": "0"
						},
						"key": "KEEPER_REGISTRY_SYNC_UPKEEP_QUEUE_SIZE"
					}, {
						"config": {
							"value": "30s"
						},
						"key": "LEASE_LOCK_DURATION"
					}, {
						"config": {
							"value": "1s"
						},
						"key": "LEASE_LOCK_REFRESH_INTERVAL"
					}, {
						"config": {
							"value": ""
						},
						"key": "LINK_CONTRACT_ADDRESS"
					}, {
						"config": {
							"value": ""
						},
						"key": "FLAGS_CONTRACT_ADDRESS"
					}, {
						"config": {
							"value": "debug"
						},
						"key": "LOG_LEVEL"
					}, {
						"config": {
							"value": "false"
						},
						"key": "LOG_SQL_MIGRATIONS"
					}, {
						"config": {
							"value": "true"
						},
						"key": "LOG_SQL"
					}, {
						"config": {
							"value": "true"
						},
						"key": "LOG_TO_DISK"
					}, {
						"config": {
							"value": "30s"
						},
						"key": "TRIGGER_FALLBACK_DB_POLL_INTERVAL"
					}, {
						"config": {
							"value": "1"
						},
						"key": "OCR_DEFAULT_TRANSACTION_QUEUE_DEPTH"
					}, {
						"config": {
							"value": "false"
						},
						"key": "OCR_TRACE_LOGGING"
					}, {
						"config": {
							"value": "V1"
						},
						"key": "P2P_NETWORKING_STACK"
					}, {
						"config": {
							"value": ""
						},
						"key": "P2P_PEER_ID"
					}, {
						"config": {
							"value": "10"
						},
						"key": "P2P_INCOMING_MESSAGE_BUFFER_SIZE"
					}, {
						"config": {
							"value": "10"
						},
						"key": "P2P_OUTGOING_MESSAGE_BUFFER_SIZE"
					}, {
						"config": {
							"value": "[]"
						},
						"key": "P2P_BOOTSTRAP_PEERS"
					}, {
						"config": {
							"value": "0.0.0.0"
						},
						"key": "P2P_LISTEN_IP"
					}, {
						"config": {
							"value": ""
						},
						"key": "P2P_LISTEN_PORT"
					}, {
						"config": {
							"value": "10s"
						},
						"key": "P2P_NEW_STREAM_TIMEOUT"
					}, {
						"config": {
							"value": "10"
						},
						"key": "P2P_DHT_LOOKUP_INTERVAL"
					}, {
						"config": {
							"value": "20s"
						},
						"key": "P2P_BOOTSTRAP_CHECK_INTERVAL"
					}, {
						"config": {
							"value": "[]"
						},
						"key": "P2PV2_ANNOUNCE_ADDRESSES"
					}, {
						"config": {
							"value": "[]"
						},
						"key": "P2PV2_BOOTSTRAPPERS"
					}, {
						"config": {
							"value": "15s"
						},
						"key": "P2PV2_DELTA_DIAL"
					}, {
						"config": {
							"value": "1m0s"
						},
						"key": "P2PV2_DELTA_RECONCILE"
					}, {
						"config": {
							"value": "[]"
						},
						"key": "P2PV2_LISTEN_ADDRESSES"
					}, {
						"config": {
							"value": "6688"
						},
						"key": "CHAINLINK_PORT"
					}, {
						"config": {
							"value": "240h0m0s"
						},
						"key": "REAPER_EXPIRATION"
					}, {
						"config": {
							"value": "-1"
						},
						"key": "REPLAY_FROM_BLOCK"
					}, {
						"config": {
							"value": "/tmp/chainlink_test/gql-test"
						},
						"key": "ROOT"
					}, {
						"config": {
							"value": "true"
						},
						"key": "SECURE_COOKIES"
					}, {
						"config": {
							"value": "2m0s"
						},
						"key": "SESSION_TIMEOUT"
					}, {
						"config": {
							"value": "false"
						},
						"key": "TELEMETRY_INGRESS_LOGGING"
					}, {
						"config": {
							"value": ""
						},
						"key": "TELEMETRY_INGRESS_SERVER_PUB_KEY"
					}, {
						"config": {
							"value": ""
						},
						"key": "TELEMETRY_INGRESS_URL"
					}, {
						"config": {
							"value": ""
						},
						"key": "CHAINLINK_TLS_HOST"
					}, {
						"config": {
							"value": "6689"
						},
						"key": "CHAINLINK_TLS_PORT"
					}, {
						"config": {
							"value": "false"
						},
						"key": "CHAINLINK_TLS_REDIRECT"
					}]
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}
