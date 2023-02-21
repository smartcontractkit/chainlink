package envvar

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestConfigSchema(t *testing.T) {
	items := map[string]string{
		"AdvisoryLockCheckInterval":                      "ADVISORY_LOCK_CHECK_INTERVAL",
		"AdvisoryLockID":                                 "ADVISORY_LOCK_ID",
		"AuditLoggerEnabled":                             "AUDIT_LOGGER_ENABLED",
		"AuditLoggerForwardToUrl":                        "AUDIT_LOGGER_FORWARD_TO_URL",
		"AuditLoggerHeaders":                             "AUDIT_LOGGER_HEADERS",
		"AuditLoggerJsonWrapperKey":                      "AUDIT_LOGGER_JSON_WRAPPER_KEY",
		"AllowOrigins":                                   "ALLOW_ORIGINS",
		"AuthenticatedRateLimit":                         "AUTHENTICATED_RATE_LIMIT",
		"AuthenticatedRateLimitPeriod":                   "AUTHENTICATED_RATE_LIMIT_PERIOD",
		"AutoPprofBlockProfileRate":                      "AUTO_PPROF_BLOCK_PROFILE_RATE",
		"AutoPprofCPUProfileRate":                        "AUTO_PPROF_CPU_PROFILE_RATE",
		"AutoPprofEnabled":                               "AUTO_PPROF_ENABLED",
		"AutoPprofGatherDuration":                        "AUTO_PPROF_GATHER_DURATION",
		"AutoPprofGatherTraceDuration":                   "AUTO_PPROF_GATHER_TRACE_DURATION",
		"AutoPprofGoroutineThreshold":                    "AUTO_PPROF_GOROUTINE_THRESHOLD",
		"AutoPprofMaxProfileSize":                        "AUTO_PPROF_MAX_PROFILE_SIZE",
		"AutoPprofMemProfileRate":                        "AUTO_PPROF_MEM_PROFILE_RATE",
		"AutoPprofMemThreshold":                          "AUTO_PPROF_MEM_THRESHOLD",
		"AutoPprofMutexProfileFraction":                  "AUTO_PPROF_MUTEX_PROFILE_FRACTION",
		"AutoPprofPollInterval":                          "AUTO_PPROF_POLL_INTERVAL",
		"AutoPprofProfileRoot":                           "AUTO_PPROF_PROFILE_ROOT",
		"BalanceMonitorEnabled":                          "BALANCE_MONITOR_ENABLED",
		"BlockBackfillDepth":                             "BLOCK_BACKFILL_DEPTH",
		"BlockBackfillSkip":                              "BLOCK_BACKFILL_SKIP",
		"BlockEmissionIdleWarningThreshold":              "BLOCK_EMISSION_IDLE_WARNING_THRESHOLD",
		"BlockHistoryEstimatorBatchSize":                 "BLOCK_HISTORY_ESTIMATOR_BATCH_SIZE",
		"BlockHistoryEstimatorBlockDelay":                "BLOCK_HISTORY_ESTIMATOR_BLOCK_DELAY",
		"BlockHistoryEstimatorBlockHistorySize":          "BLOCK_HISTORY_ESTIMATOR_BLOCK_HISTORY_SIZE",
		"BlockHistoryEstimatorCheckInclusionBlocks":      "BLOCK_HISTORY_ESTIMATOR_CHECK_INCLUSION_BLOCKS",
		"BlockHistoryEstimatorCheckInclusionPercentile":  "BLOCK_HISTORY_ESTIMATOR_CHECK_INCLUSION_PERCENTILE",
		"BlockHistoryEstimatorEIP1559FeeCapBufferBlocks": "BLOCK_HISTORY_ESTIMATOR_EIP1559_FEE_CAP_BUFFER_BLOCKS",
		"BlockHistoryEstimatorTransactionPercentile":     "BLOCK_HISTORY_ESTIMATOR_TRANSACTION_PERCENTILE",
		"BridgeResponseURL":                              "BRIDGE_RESPONSE_URL",
		"BridgeCacheTTL":                                 "BRIDGE_CACHE_TTL",
		"ChainType":                                      "CHAIN_TYPE",
		"DatabaseBackupDir":                              "DATABASE_BACKUP_DIR",
		"DatabaseBackupFrequency":                        "DATABASE_BACKUP_FREQUENCY",
		"DatabaseBackupMode":                             "DATABASE_BACKUP_MODE",
		"DatabaseBackupOnVersionUpgrade":                 "DATABASE_BACKUP_ON_VERSION_UPGRADE",
		"DatabaseBackupURL":                              "DATABASE_BACKUP_URL",
		"DatabaseListenerMaxReconnectDuration":           "DATABASE_LISTENER_MAX_RECONNECT_DURATION",
		"DatabaseListenerMinReconnectInterval":           "DATABASE_LISTENER_MIN_RECONNECT_INTERVAL",
		"DatabaseLockingMode":                            "DATABASE_LOCKING_MODE",
		"DatabaseURL":                                    "DATABASE_URL",
		"DefaultChainID":                                 "ETH_CHAIN_ID",
		"DefaultHTTPLimit":                               "DEFAULT_HTTP_LIMIT",
		"DefaultHTTPTimeout":                             "DEFAULT_HTTP_TIMEOUT",
		"Dev":                                            "CHAINLINK_DEV",
		"EVMEnabled":                                     "EVM_ENABLED",
		"EVMRPCEnabled":                                  "EVM_RPC_ENABLED",
		"EthTxReaperInterval":                            "ETH_TX_REAPER_INTERVAL",
		"EthTxReaperThreshold":                           "ETH_TX_REAPER_THRESHOLD",
		"EthTxResendAfterThreshold":                      "ETH_TX_RESEND_AFTER_THRESHOLD",
		"EthereumHTTPURL":                                "ETH_HTTP_URL",
		"EthereumSecondaryURL":                           "ETH_SECONDARY_URL",
		"EthereumSecondaryURLs":                          "ETH_SECONDARY_URLS",
		"EthereumURL":                                    "ETH_URL",
		"EthereumNodes":                                  "EVM_NODES",
		"EvmBalanceMonitorBlockDelay":                    "ETH_BALANCE_MONITOR_BLOCK_DELAY",
		"EvmEIP1559DynamicFees":                          "EVM_EIP1559_DYNAMIC_FEES",
		"EvmFinalityDepth":                               "ETH_FINALITY_DEPTH",
		"EvmGasBumpPercent":                              "ETH_GAS_BUMP_PERCENT",
		"EvmGasBumpThreshold":                            "ETH_GAS_BUMP_THRESHOLD",
		"EvmGasBumpTxDepth":                              "ETH_GAS_BUMP_TX_DEPTH",
		"EvmGasBumpWei":                                  "ETH_GAS_BUMP_WEI",
		"EvmGasFeeCapDefault":                            "EVM_GAS_FEE_CAP_DEFAULT",
		"EvmGasLimitDefault":                             "ETH_GAS_LIMIT_DEFAULT",
		"EvmGasLimitMax":                                 "ETH_GAS_LIMIT_MAX",
		"EvmGasLimitMultiplier":                          "ETH_GAS_LIMIT_MULTIPLIER",
		"EvmGasLimitTransfer":                            "ETH_GAS_LIMIT_TRANSFER",
		"EvmGasLimitOCRJobType":                          "ETH_GAS_LIMIT_OCR_JOB_TYPE",
		"EvmGasLimitDRJobType":                           "ETH_GAS_LIMIT_DR_JOB_TYPE",
		"EvmGasLimitVRFJobType":                          "ETH_GAS_LIMIT_VRF_JOB_TYPE",
		"EvmGasLimitFMJobType":                           "ETH_GAS_LIMIT_FM_JOB_TYPE",
		"EvmGasLimitKeeperJobType":                       "ETH_GAS_LIMIT_KEEPER_JOB_TYPE",
		"EvmGasPriceDefault":                             "ETH_GAS_PRICE_DEFAULT",
		"EvmGasTipCapDefault":                            "EVM_GAS_TIP_CAP_DEFAULT",
		"EvmGasTipCapMinimum":                            "EVM_GAS_TIP_CAP_MINIMUM",
		"EvmHeadTrackerHistoryDepth":                     "ETH_HEAD_TRACKER_HISTORY_DEPTH",
		"EvmHeadTrackerMaxBufferSize":                    "ETH_HEAD_TRACKER_MAX_BUFFER_SIZE",
		"EvmHeadTrackerSamplingInterval":                 "ETH_HEAD_TRACKER_SAMPLING_INTERVAL",
		"EvmLogBackfillBatchSize":                        "ETH_LOG_BACKFILL_BATCH_SIZE",
		"EvmLogPollInterval":                             "ETH_LOG_POLL_INTERVAL",
		"EvmLogKeepBlocksDepth":                          "ETH_LOG_KEEP_BLOCKS_DEPTH",
		"EvmMaxGasPriceWei":                              "ETH_MAX_GAS_PRICE_WEI",
		"EvmMaxInFlightTransactions":                     "ETH_MAX_IN_FLIGHT_TRANSACTIONS",
		"EvmMaxQueuedTransactions":                       "ETH_MAX_QUEUED_TRANSACTIONS",
		"EvmMinGasPriceWei":                              "ETH_MIN_GAS_PRICE_WEI",
		"EvmNonceAutoSync":                               "ETH_NONCE_AUTO_SYNC",
		"EvmUseForwarders":                               "ETH_USE_FORWARDERS",
		"EvmRPCDefaultBatchSize":                         "ETH_RPC_DEFAULT_BATCH_SIZE",
		"ExplorerAccessKey":                              "EXPLORER_ACCESS_KEY",
		"ExplorerSecret":                                 "EXPLORER_SECRET",
		"ExplorerURL":                                    "EXPLORER_URL",
		"FMDefaultTransactionQueueDepth":                 "FM_DEFAULT_TRANSACTION_QUEUE_DEPTH",
		"FMSimulateTransactions":                         "FM_SIMULATE_TRANSACTIONS",
		"FeatureExternalInitiators":                      "FEATURE_EXTERNAL_INITIATORS",
		"FeatureFeedsManager":                            "FEATURE_FEEDS_MANAGER",
		"FeatureLogPoller":                               "FEATURE_LOG_POLLER",
		"FeatureOffchainReporting":                       "FEATURE_OFFCHAIN_REPORTING",
		"FeatureOffchainReporting2":                      "FEATURE_OFFCHAIN_REPORTING2",
		"FeatureUICSAKeys":                               "FEATURE_UI_CSA_KEYS",
		"FlagsContractAddress":                           "FLAGS_CONTRACT_ADDRESS",
		"GasEstimatorMode":                               "GAS_ESTIMATOR_MODE",
		"GasUpdaterBatchSize":                            "GAS_UPDATER_BATCH_SIZE",
		"GasUpdaterBlockDelay":                           "GAS_UPDATER_BLOCK_DELAY",
		"GasUpdaterBlockHistorySize":                     "GAS_UPDATER_BLOCK_HISTORY_SIZE",
		"GasUpdaterTransactionPercentile":                "GAS_UPDATER_TRANSACTION_PERCENTILE",
		"HTTPServerWriteTimeout":                         "HTTP_SERVER_WRITE_TIMEOUT",
		"InsecureFastScrypt":                             "INSECURE_FAST_SCRYPT",
		"JSONConsole":                                    "JSON_CONSOLE",
		"JobPipelineMaxRunDuration":                      "JOB_PIPELINE_MAX_RUN_DURATION",
		"JobPipelineMaxSuccessfulRuns":                   "JOB_PIPELINE_MAX_SUCCESSFUL_RUNS",
		"JobPipelineReaperInterval":                      "JOB_PIPELINE_REAPER_INTERVAL",
		"JobPipelineReaperThreshold":                     "JOB_PIPELINE_REAPER_THRESHOLD",
		"JobPipelineResultWriteQueueDepth":               "JOB_PIPELINE_RESULT_WRITE_QUEUE_DEPTH",
		"KeeperDefaultTransactionQueueDepth":             "KEEPER_DEFAULT_TRANSACTION_QUEUE_DEPTH",
		"KeeperGasPriceBufferPercent":                    "KEEPER_GAS_PRICE_BUFFER_PERCENT",
		"KeeperGasTipCapBufferPercent":                   "KEEPER_GAS_TIP_CAP_BUFFER_PERCENT",
		"KeeperBaseFeeBufferPercent":                     "KEEPER_BASE_FEE_BUFFER_PERCENT",
		"KeeperMaximumGracePeriod":                       "KEEPER_MAXIMUM_GRACE_PERIOD",
		"KeeperRegistryCheckGasOverhead":                 "KEEPER_REGISTRY_CHECK_GAS_OVERHEAD",
		"KeeperRegistryPerformGasOverhead":               "KEEPER_REGISTRY_PERFORM_GAS_OVERHEAD",
		"KeeperRegistryMaxPerformDataSize":               "KEEPER_REGISTRY_MAX_PERFORM_DATA_SIZE",
		"KeeperRegistrySyncInterval":                     "KEEPER_REGISTRY_SYNC_INTERVAL",
		"KeeperRegistrySyncUpkeepQueueSize":              "KEEPER_REGISTRY_SYNC_UPKEEP_QUEUE_SIZE",
		"KeeperTurnLookBack":                             "KEEPER_TURN_LOOK_BACK",
		"LeaseLockDuration":                              "LEASE_LOCK_DURATION",
		"LeaseLockRefreshInterval":                       "LEASE_LOCK_REFRESH_INTERVAL",
		"LinkContractAddress":                            "LINK_CONTRACT_ADDRESS",
		"OCR2AutomationGasLimit":                         "OCR2_AUTOMATION_GAS_LIMIT",
		"OperatorFactoryAddress":                         "OPERATOR_FACTORY_ADDRESS",
		"LogFileDir":                                     "LOG_FILE_DIR",
		"LogLevel":                                       "LOG_LEVEL",
		"LogSQL":                                         "LOG_SQL",
		"LogFileMaxSize":                                 "LOG_FILE_MAX_SIZE",
		"LogFileMaxAge":                                  "LOG_FILE_MAX_AGE",
		"LogFileMaxBackups":                              "LOG_FILE_MAX_BACKUPS",
		"LogUnixTS":                                      "LOG_UNIX_TS",
		"MaximumServiceDuration":                         "MAXIMUM_SERVICE_DURATION",
		"MigrateDatabase":                                "MIGRATE_DATABASE",
		"MinIncomingConfirmations":                       "MIN_INCOMING_CONFIRMATIONS",
		"MinimumContractPayment":                         "MINIMUM_CONTRACT_PAYMENT_LINK_JUELS",
		"MinimumServiceDuration":                         "MINIMUM_SERVICE_DURATION",
		"NodeNoNewHeadsThreshold":                        "NODE_NO_NEW_HEADS_THRESHOLD",
		"NodePollFailureThreshold":                       "NODE_POLL_FAILURE_THRESHOLD",
		"NodePollInterval":                               "NODE_POLL_INTERVAL",
		"NodeSelectionMode":                              "NODE_SELECTION_MODE",
		"NodeSyncThreshold":                              "NODE_SYNC_THRESHOLD",
		"ORMMaxIdleConns":                                "ORM_MAX_IDLE_CONNS",
		"ORMMaxOpenConns":                                "ORM_MAX_OPEN_CONNS",
		"OptimismGasFees":                                "OPTIMISM_GAS_FEES",
		"Port":                                           "CHAINLINK_PORT",
		"RPCEnabled":                                     "RPC_ENABLED",
		"RPID":                                           "MFA_RPID",
		"RPOrigin":                                       "MFA_RPORIGIN",
		"ReaperExpiration":                               "REAPER_EXPIRATION",
		"RootDir":                                        "ROOT",
		"SecureCookies":                                  "SECURE_COOKIES",
		"SessionTimeout":                                 "SESSION_TIMEOUT",
		"ShutdownGracePeriod":                            "SHUTDOWN_GRACE_PERIOD",
		"SolanaEnabled":                                  "SOLANA_ENABLED",
		"SolanaNodes":                                    "SOLANA_NODES",
		"StarknetEnabled":                                "STARKNET_ENABLED",
		"StarknetNodes":                                  "STARKNET_NODES",
		"TLSCertPath":                                    "TLS_CERT_PATH",
		"TLSHost":                                        "CHAINLINK_TLS_HOST",
		"TLSKeyPath":                                     "TLS_KEY_PATH",
		"TLSPort":                                        "CHAINLINK_TLS_PORT",
		"TLSRedirect":                                    "CHAINLINK_TLS_REDIRECT",
		"TelemetryIngressBufferSize":                     "TELEMETRY_INGRESS_BUFFER_SIZE",
		"TelemetryIngressLogging":                        "TELEMETRY_INGRESS_LOGGING",
		"TelemetryIngressUniConn":                        "TELEMETRY_INGRESS_UNICONN",
		"TelemetryIngressMaxBatchSize":                   "TELEMETRY_INGRESS_MAX_BATCH_SIZE",
		"TelemetryIngressSendInterval":                   "TELEMETRY_INGRESS_SEND_INTERVAL",
		"TelemetryIngressSendTimeout":                    "TELEMETRY_INGRESS_SEND_TIMEOUT",
		"TelemetryIngressServerPubKey":                   "TELEMETRY_INGRESS_SERVER_PUB_KEY",
		"TelemetryIngressURL":                            "TELEMETRY_INGRESS_URL",
		"TelemetryIngressUseBatchSend":                   "TELEMETRY_INGRESS_USE_BATCH_SEND",
		"TriggerFallbackDBPollInterval":                  "TRIGGER_FALLBACK_DB_POLL_INTERVAL",
		"UnAuthenticatedRateLimit":                       "UNAUTHENTICATED_RATE_LIMIT",
		"UnAuthenticatedRateLimitPeriod":                 "UNAUTHENTICATED_RATE_LIMIT_PERIOD",

		// OCR v2
		"OCR2BlockchainTimeout":                  "OCR2_BLOCKCHAIN_TIMEOUT",
		"OCR2ContractPollInterval":               "OCR2_CONTRACT_POLL_INTERVAL",
		"OCR2ContractSubscribeInterval":          "OCR2_CONTRACT_SUBSCRIBE_INTERVAL",
		"OCR2ContractTransmitterTransmitTimeout": "OCR2_CONTRACT_TRANSMITTER_TRANSMIT_TIMEOUT",
		"OCR2DatabaseTimeout":                    "OCR2_DATABASE_TIMEOUT",
		"OCR2ContractConfirmations":              "OCR2_CONTRACT_CONFIRMATIONS",
		"OCR2KeyBundleID":                        "OCR2_KEY_BUNDLE_ID",
		"OCR2TraceLogging":                       "OCR2_TRACE_LOGGING",

		// OCR v1
		"OCRBlockchainTimeout":                  "OCR_BLOCKCHAIN_TIMEOUT",
		"OCRContractPollInterval":               "OCR_CONTRACT_POLL_INTERVAL",
		"OCRContractSubscribeInterval":          "OCR_CONTRACT_SUBSCRIBE_INTERVAL",
		"OCRContractTransmitterTransmitTimeout": "OCR_CONTRACT_TRANSMITTER_TRANSMIT_TIMEOUT",
		"OCRDatabaseTimeout":                    "OCR_DATABASE_TIMEOUT",
		"OCRContractConfirmations":              "OCR_CONTRACT_CONFIRMATIONS",
		"OCRKeyBundleID":                        "OCR_KEY_BUNDLE_ID",
		"OCRDefaultTransactionQueueDepth":       "OCR_DEFAULT_TRANSACTION_QUEUE_DEPTH",
		"OCRTraceLogging":                       "OCR_TRACE_LOGGING",
		"OCRObservationGracePeriod":             "OCR_OBSERVATION_GRACE_PERIOD",
		"OCRObservationTimeout":                 "OCR_OBSERVATION_TIMEOUT",
		"OCRTransmitterAddress":                 "OCR_TRANSMITTER_ADDRESS",
		"OCRSimulateTransactions":               "OCR_SIMULATE_TRANSACTIONS",

		// P2P v1 and v2 networking
		"P2PNetworkingStack":           "P2P_NETWORKING_STACK",
		"P2PPeerID":                    "P2P_PEER_ID",
		"P2PIncomingMessageBufferSize": "P2P_INCOMING_MESSAGE_BUFFER_SIZE",
		"P2POutgoingMessageBufferSize": "P2P_OUTGOING_MESSAGE_BUFFER_SIZE",

		// P2P v1 networking
		"P2PBootstrapCheckInterval":           "P2P_BOOTSTRAP_CHECK_INTERVAL",
		"P2PDHTLookupInterval":                "P2P_DHT_LOOKUP_INTERVAL",
		"P2PNewStreamTimeout":                 "P2P_NEW_STREAM_TIMEOUT",
		"P2PAnnounceIP":                       "P2P_ANNOUNCE_IP",
		"P2PAnnouncePort":                     "P2P_ANNOUNCE_PORT",
		"P2PBootstrapPeers":                   "P2P_BOOTSTRAP_PEERS",
		"P2PDHTAnnouncementCounterUserPrefix": "P2P_DHT_ANNOUNCEMENT_COUNTER_USER_PREFIX",
		"P2PListenIP":                         "P2P_LISTEN_IP",
		"P2PListenPort":                       "P2P_LISTEN_PORT",
		"P2PPeerstoreWriteInterval":           "P2P_PEERSTORE_WRITE_INTERVAL",

		// P2P v2 networking
		"P2PV2AccountAddresses":  "P2PV2_ANNOUNCE_ADDRESSES",
		"P2PV2AnnounceAddresses": "P2PV2_ANNOUNCE_ADDRESSES",
		"P2PV2Bootstrappers":     "P2PV2_BOOTSTRAPPERS",
		"P2PV2DeltaDial":         "P2PV2_DELTA_DIAL",
		"P2PV2DeltaReconcile":    "P2PV2_DELTA_RECONCILE",
		"P2PV2ListenAddresses":   "P2PV2_LISTEN_ADDRESSES",

		// Pyroscope profiling
		"PyroscopeAuthToken":     "PYROSCOPE_AUTH_TOKEN",
		"PyroscopeServerAddress": "PYROSCOPE_SERVER_ADDRESS",
		"PyroscopeEnvironment":   "PYROSCOPE_ENVIRONMENT",

		// P2P deprecated
		"OCRNewStreamTimeout":          "OCR_NEW_STREAM_TIMEOUT",
		"OCRBootstrapCheckInterval":    "OCR_BOOTSTRAP_CHECK_INTERVAL",
		"OCRDHTLookupInterval":         "OCR_DHT_LOOKUP_INTERVAL",
		"OCRIncomingMessageBufferSize": "OCR_INCOMING_MESSAGE_BUFFER_SIZE",
		"OCROutgoingMessageBufferSize": "OCR_OUTGOING_MESSAGE_BUFFER_SIZE",
	}

	schemaT := reflect.TypeOf(ConfigSchema{})
	for i := 0; i < schemaT.NumField(); i++ {
		field := schemaT.FieldByIndex([]int{i})
		item, found := items[field.Name]

		//
		// ╭──╮   ╭────────────────────────────────────╮
		// │  │   │ It looks like you're trying to add │
		// @  @  ╭│ a new configuration variable!      │
		// ││ ││ ││                                    │
		// ││ ││ ╯╰────────────────────────────────────╯
		// │╰─╯│
		// ╰───╯

		msg := utils.BoxOutput(`New configuration variable detected: '%s'. Please take the following steps:
0. Are you SURE that this config variable ought to apply to ALL chains? If this
   needs to be a chain-specific variable, you should follow the pattern of making
   it a Global and using the chain-specific overridable version in the chain
   scoped config instead (for more info see notion doc here: https://www.notion.so/chainlink/Config-in-Chainlink-4f36fb8f28f241f7b87cd1857df54db7).
1. Make sure that the method has a comment explaining in detail what the new config var does
2. Update the Changelog
3. Update the ConfigPrinter found in core/config/presenters.go if you
   think this variable needs to be shown in the UI
4. Make a PR into the documentation page if this is not a "nodoc" env var
   this (found at https://github.com/smartcontractkit/documentation/blob/main/docs/Node%%20Operators/configuration-variables.md).
   You may be able to add to an existing PR for the next version.
   Don't forget to update TOC.
5. Add your new config variable to this test`, field.Name)
		assert.True(t, found, msg)
		env := field.Tag.Get("env")
		assert.Equal(t, item, env)
	}
}
