package resolver

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
)

var (
	//go:embed testdata/config-empty-effective.toml
	configEmptyEffective string
	//go:embed testdata/config-full.toml
	configFull string
	//go:embed testdata/config-multi-chain.toml
	configMulti string
	//go:embed testdata/config-multi-chain-effective.toml
	configMultiEffective string
)

func TestResolver_ConfigV2(t *testing.T) {
	t.Parallel()

	query := `
		query FetchConfigV2 {
			configv2 {
				user
				effective
			}
	  	}`

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "configv2"),
		{
			name:          "empty",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				var opts chainlink.GeneralConfigOpts
				require.NoError(f.t, opts.ParseTOML("", ""))
				cfg, err := opts.New(logger.TestLogger(f.t))
				require.NoError(t, err)
				f.App.On("GetConfig").Return(cfg)
			},
			query:  query,
			result: fmt.Sprintf(`{"configv2":{"user":"","effective":%s}}`, mustJSONMarshal(t, configEmptyEffective)),
		},
		{
			name:          "full",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				var opts chainlink.GeneralConfigOpts
				require.NoError(f.t, opts.ParseTOML(configFull, ""))
				cfg, err := opts.New(logger.TestLogger(f.t))
				require.NoError(t, err)
				f.App.On("GetConfig").Return(cfg)
			},
			query:  query,
			result: fmt.Sprintf(`{"configv2":{"user":%s,"effective":%s}}`, mustJSONMarshal(t, configFull), mustJSONMarshal(t, configFull)),
		},
		{
			name:          "partial",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				var opts chainlink.GeneralConfigOpts
				require.NoError(f.t, opts.ParseTOML(configMulti, ""))
				cfg, err := opts.New(logger.TestLogger(f.t))
				require.NoError(t, err)
				f.App.On("GetConfig").Return(cfg)
			},
			query:  query,
			result: fmt.Sprintf(`{"configv2":{"user":%s,"effective":%s}}`, mustJSONMarshal(t, configMulti), mustJSONMarshal(t, configMultiEffective)),
		},
	}

	RunGQLTests(t, testCases)
}

func mustJSONMarshal(t *testing.T, s string) string {
	b, err := json.Marshal(s)
	require.NoError(t, err)
	return string(b)
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func TestResolver_Config(t *testing.T) {
	t.Parallel()

	query := `
		query GetConfiguration {
			config {
				items {
					key
					value
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
				logLevel := zapcore.ErrorLevel
				cfg := configtest.NewTestGeneralConfigWithOverrides(t, configtest.GeneralConfigOverrides{
					AdvisoryLockID:                       null.IntFrom(1),
					AllowOrigins:                         null.StringFrom("test"),
					BlockBackfillDepth:                   null.IntFrom(1),
					BlockBackfillSkip:                    null.BoolFrom(false),
					DatabaseURL:                          null.StringFrom("test"),
					DefaultChainID:                       nil,
					DefaultHTTPTimeout:                   nil,
					Dev:                                  null.BoolFrom(true),
					ShutdownGracePeriod:                  nil,
					Dialect:                              "",
					EVMEnabled:                           null.BoolFrom(false),
					EVMRPCEnabled:                        null.BoolFrom(false),
					EthereumURL:                          null.StringFrom(""),
					FeatureExternalInitiators:            null.BoolFrom(true),
					GlobalBalanceMonitorEnabled:          null.BoolFrom(true),
					GlobalChainType:                      null.StringFrom(""),
					GlobalEthTxReaperThreshold:           nil,
					GlobalEthTxResendAfterThreshold:      nil,
					GlobalEvmEIP1559DynamicFees:          null.BoolFrom(true),
					GlobalEvmFinalityDepth:               null.IntFrom(1),
					GlobalEvmGasBumpPercent:              null.IntFrom(1),
					GlobalEvmGasBumpTxDepth:              null.IntFrom(1),
					GlobalEvmGasBumpWei:                  nil,
					GlobalEvmGasLimitDefault:             null.IntFrom(1),
					GlobalEvmGasLimitMax:                 null.IntFrom(10),
					GlobalEvmGasLimitMultiplier:          null.FloatFrom(1),
					GlobalEvmGasPriceDefault:             nil,
					GlobalEvmGasTipCapDefault:            nil,
					GlobalEvmGasTipCapMinimum:            nil,
					GlobalEvmGasLimitOCRJobType:          null.IntFrom(10),
					GlobalEvmGasLimitDRJobType:           null.IntFrom(11),
					GlobalEvmGasLimitVRFJobType:          null.IntFrom(12),
					GlobalEvmGasLimitFMJobType:           null.IntFrom(13),
					GlobalEvmGasLimitKeeperJobType:       null.IntFrom(14),
					GlobalEvmHeadTrackerHistoryDepth:     null.IntFrom(1),
					GlobalEvmHeadTrackerMaxBufferSize:    null.IntFrom(1),
					GlobalEvmHeadTrackerSamplingInterval: nil,
					GlobalEvmLogBackfillBatchSize:        null.IntFrom(1),
					GlobalEvmMaxGasPriceWei:              nil,
					GlobalEvmMinGasPriceWei:              nil,
					GlobalEvmNonceAutoSync:               null.BoolFrom(false),
					GlobalEvmRPCDefaultBatchSize:         null.IntFrom(1),
					GlobalFlagsContractAddress:           null.StringFrom("test"),
					GlobalGasEstimatorMode:               null.StringFrom("test"),
					GlobalMinIncomingConfirmations:       null.IntFrom(1),
					GlobalMinimumContractPayment:         nil,
					KeeperMaximumGracePeriod:             null.IntFrom(1),
					KeeperRegistrySyncInterval:           nil,
					KeeperRegistrySyncUpkeepQueueSize:    null.IntFrom(1),
					KeeperTurnLookBack:                   null.IntFrom(0),
					LogLevel:                             &logLevel,
					DefaultLogLevel:                      nil,
					LogFileDir:                           null.StringFrom("foo"),
					LogSQL:                               null.BoolFrom(true),
					LogFileMaxSize:                       null.StringFrom("100mb"),
					LogFileMaxAge:                        null.IntFrom(3),
					LogFileMaxBackups:                    null.IntFrom(12),
					OCRKeyBundleID:                       null.StringFrom("test"),
					OCRObservationTimeout:                nil,
					OCRTransmitterAddress:                nil,
					P2PBootstrapPeers:                    nil,
					P2PListenPort:                        null.IntFrom(1),
					P2PPeerID:                            "",
					TriggerFallbackDBPollInterval:        nil,
				})
				cfg.SetRootDir("/tmp/chainlink_test/gql-test")

				f.App.On("GetConfig").Return(cfg)
			},
			query: query,
			result: `
{
  "config": {
    "items": [
      {
        "key": "ADVISORY_LOCK_CHECK_INTERVAL",
        "value": "1s"
      },
      {
        "key": "ADVISORY_LOCK_ID",
        "value": "1027321974924625846"
      },
      {
        "key": "ALLOW_ORIGINS",
        "value": "test"
      },
      {
        "key": "BLOCK_BACKFILL_DEPTH",
        "value": "1"
      },
      {
        "key": "BLOCK_HISTORY_ESTIMATOR_BLOCK_DELAY",
        "value": "0"
      },
      {
        "key": "BLOCK_HISTORY_ESTIMATOR_BLOCK_HISTORY_SIZE",
        "value": "0"
      },
      {
        "key": "BLOCK_HISTORY_ESTIMATOR_TRANSACTION_PERCENTILE",
        "value": "0"
      },
      {
        "key": "BRIDGE_RESPONSE_URL",
        "value": "http://localhost:6688"
      },
      {
        "key": "CHAIN_TYPE",
        "value": ""
      },
      {
        "key": "DATABASE_BACKUP_FREQUENCY",
        "value": "1h0m0s"
      },
      {
        "key": "DATABASE_BACKUP_MODE",
        "value": "none"
      },
      {
        "key": "DATABASE_BACKUP_ON_VERSION_UPGRADE",
        "value": "true"
      },
      {
        "key": "DATABASE_LOCKING_MODE",
        "value": "none"
      },
      {
        "key": "ETH_CHAIN_ID",
        "value": "0"
      },
      {
        "key": "DEFAULT_HTTP_LIMIT",
        "value": "32768"
      },
      {
        "key": "DEFAULT_HTTP_TIMEOUT",
        "value": "15s"
      },
      {
        "key": "CHAINLINK_DEV",
        "value": "true"
      },
	  {
		"key":"SHUTDOWN_GRACE_PERIOD",
		"value":"5s"
	  },
      {
        "key": "EVM_RPC_ENABLED",
        "value": "false"
      },
      {
        "key": "ETH_HTTP_URL",
        "value": ""
      },
      {
        "key": "ETH_SECONDARY_URLS",
        "value": "[]"
      },
      {
        "key": "ETH_URL",
        "value": ""
      },
      {
        "key": "EXPLORER_URL",
        "value": ""
      },
      {
        "key": "FM_DEFAULT_TRANSACTION_QUEUE_DEPTH",
        "value": "1"
      },
      {
        "key": "FEATURE_EXTERNAL_INITIATORS",
        "value": "true"
      },
      {
        "key": "FEATURE_OFFCHAIN_REPORTING",
        "value": "false"
      },
      {
        "key": "GAS_ESTIMATOR_MODE",
        "value": ""
      },
      {
        "key": "INSECURE_FAST_SCRYPT",
        "value": "true"
      },
      {
        "key": "JSON_CONSOLE",
        "value": "false"
      },
      {
        "key": "JOB_PIPELINE_REAPER_INTERVAL",
        "value": "1h0m0s"
      },
      {
        "key": "JOB_PIPELINE_REAPER_THRESHOLD",
        "value": "24h0m0s"
      },
      {
        "key": "KEEPER_DEFAULT_TRANSACTION_QUEUE_DEPTH",
        "value": "1"
      },
      {
        "key": "KEEPER_GAS_PRICE_BUFFER_PERCENT",
        "value": "20"
      },
      {
        "key": "KEEPER_GAS_TIP_CAP_BUFFER_PERCENT",
        "value": "20"
      },
      {
        "key": "KEEPER_BASE_FEE_BUFFER_PERCENT",
        "value": "20"
      },
      {
        "key": "KEEPER_MAXIMUM_GRACE_PERIOD",
        "value": "1"
      },
      {
        "key": "KEEPER_REGISTRY_CHECK_GAS_OVERHEAD",
        "value": "200000"
      },
      {
        "key": "KEEPER_REGISTRY_PERFORM_GAS_OVERHEAD",
        "value": "300000"
      },
      {
        "key": "KEEPER_REGISTRY_MAX_PERFORM_DATA_SIZE",
        "value": "5000"
      },
	  {
		"key":"KEEPER_REGISTRY_SYNC_INTERVAL",
		"value":"30m0s"
	  },
      {
        "key": "KEEPER_REGISTRY_SYNC_UPKEEP_QUEUE_SIZE",
        "value": "1"
      },
      {
        "key": "KEEPER_TURN_LOOK_BACK",
        "value": "0"
      },
      {
        "key": "LEASE_LOCK_DURATION",
        "value": "10s"
      },
      {
        "key": "LEASE_LOCK_REFRESH_INTERVAL",
        "value": "1s"
      },
      {
        "key": "FLAGS_CONTRACT_ADDRESS",
        "value": ""
      },
      {
        "key": "LINK_CONTRACT_ADDRESS",
        "value": ""
      },
      {
        "key": "LOG_FILE_DIR",
        "value": "foo"
      },
      {
        "key": "LOG_LEVEL",
        "value": "error"
      },
      {
        "key": "LOG_SQL",
        "value": "true"
      },
      {
        "key": "LOG_FILE_MAX_SIZE",
        "value": "100.00mb"
      },
      {
        "key": "LOG_FILE_MAX_AGE",
        "value": "3"
      },
      {
        "key": "LOG_FILE_MAX_BACKUPS",
        "value": "12"
      },
      {
        "key": "TRIGGER_FALLBACK_DB_POLL_INTERVAL",
        "value": "30s"
      },
      {
        "key": "AUDIT_LOGGER_ENABLED",
        "value": "false"
      },
      {
        "key": "AUDIT_LOGGER_FORWARD_TO_URL",
        "value": ""
      },
      {
        "key": "AUDIT_LOGGER_JSON_WRAPPER_KEY",
        "value": ""
      },
      {
        "key": "AUDIT_LOGGER_HEADERS",
        "value": ""
      },
      {
        "key": "OCR_DEFAULT_TRANSACTION_QUEUE_DEPTH",
        "value": "1"
      },
      {
        "key": "OCR_TRACE_LOGGING",
        "value": "false"
      },
      {
        "key": "P2P_NETWORKING_STACK",
        "value": "V1"
      },
      {
        "key": "P2P_PEER_ID",
        "value": ""
      },
      {
        "key": "P2P_INCOMING_MESSAGE_BUFFER_SIZE",
        "value": "10"
      },
      {
        "key": "P2P_OUTGOING_MESSAGE_BUFFER_SIZE",
        "value": "10"
      },
      {
        "key": "P2P_BOOTSTRAP_PEERS",
        "value": "[]"
      },
      {
        "key": "P2P_LISTEN_IP",
        "value": "0.0.0.0"
      },
      {
        "key": "P2P_LISTEN_PORT",
        "value": ""
      },
      {
        "key": "P2P_NEW_STREAM_TIMEOUT",
        "value": "10s"
      },
      {
        "key": "P2P_DHT_LOOKUP_INTERVAL",
        "value": "10"
      },
      {
        "key": "P2P_BOOTSTRAP_CHECK_INTERVAL",
        "value": "20s"
      },
      {
        "key": "P2PV2_ANNOUNCE_ADDRESSES",
        "value": "[]"
      },
      {
        "key": "P2PV2_BOOTSTRAPPERS",
        "value": "[]"
      },
      {
        "key": "P2PV2_DELTA_DIAL",
        "value": "15s"
      },
      {
        "key": "P2PV2_DELTA_RECONCILE",
        "value": "1m0s"
      },
      {
        "key": "P2PV2_LISTEN_ADDRESSES",
        "value": "[]"
      },
      {
        "key": "CHAINLINK_PORT",
        "value": "6688"
      },
      {
        "key": "REAPER_EXPIRATION",
        "value": "240h0m0s"
      },
      {
        "key": "ROOT",
        "value": "/tmp/chainlink_test/gql-test"
      },
      {
        "key": "SECURE_COOKIES",
        "value": "true"
      },
      {
        "key": "SESSION_TIMEOUT",
        "value": "2m0s"
      },
      {
        "key": "TELEMETRY_INGRESS_LOGGING",
        "value": "false"
      },
      {
        "key": "TELEMETRY_INGRESS_SERVER_PUB_KEY",
        "value": ""
      },
      {
        "key": "TELEMETRY_INGRESS_URL",
        "value": ""
      },
      {
        "key": "CHAINLINK_TLS_HOST",
        "value": ""
      },
      {
        "key": "CHAINLINK_TLS_PORT",
        "value": "6689"
      },
      {
        "key": "CHAINLINK_TLS_REDIRECT",
        "value": "false"
      }
    ]
  }
}
`,
		},
	}

	RunGQLTests(t, testCases)
}
