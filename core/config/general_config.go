package config

import (
	"math/big"
	"net/url"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger/audit"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// nolint
var (
	ErrEnvUnset = errors.New("env var unset")
)

// FeatureFlags contains bools that toggle various features or chains
// TODO: document the new ones
type FeatureFlags interface {
	FeatureExternalInitiators() bool
	FeatureFeedsManager() bool
	FeatureOffchainReporting() bool
	FeatureOffchainReporting2() bool
	FeatureUICSAKeys() bool
	FeatureLogPoller() bool

	AutoPprofEnabled() bool
	EVMEnabled() bool
	EVMRPCEnabled() bool
	P2PEnabled() bool
	SolanaEnabled() bool
	StarkNetEnabled() bool
}

type LogFn func(...any)

type BasicConfig interface {
	Validate() error
	LogConfiguration(log LogFn)
	SetLogLevel(lvl zapcore.Level) error
	SetLogSQL(logSQL bool)
	SetPasswords(keystore, vrf *string)

	FeatureFlags
	audit.Config

	AdvisoryLockCheckInterval() time.Duration
	AdvisoryLockID() int64
	AllowOrigins() string
	AppID() uuid.UUID
	AuthenticatedRateLimit() int64
	AuthenticatedRateLimitPeriod() models.Duration
	AutoPprofBlockProfileRate() int
	AutoPprofCPUProfileRate() int
	AutoPprofGatherDuration() models.Duration
	AutoPprofGatherTraceDuration() models.Duration
	AutoPprofGoroutineThreshold() int
	AutoPprofMaxProfileSize() utils.FileSize
	AutoPprofMemProfileRate() int
	AutoPprofMemThreshold() utils.FileSize
	AutoPprofMutexProfileFraction() int
	AutoPprofPollInterval() models.Duration
	AutoPprofProfileRoot() string
	BlockBackfillDepth() uint64
	BlockBackfillSkip() bool
	BridgeResponseURL() *url.URL
	BridgeCacheTTL() time.Duration
	CertFile() string
	DatabaseBackupDir() string
	DatabaseBackupFrequency() time.Duration
	DatabaseBackupMode() DatabaseBackupMode
	DatabaseBackupOnVersionUpgrade() bool
	DatabaseBackupURL() *url.URL
	DatabaseDefaultIdleInTxSessionTimeout() time.Duration
	DatabaseDefaultLockTimeout() time.Duration
	DatabaseDefaultQueryTimeout() time.Duration
	DatabaseListenerMaxReconnectDuration() time.Duration
	DatabaseListenerMinReconnectInterval() time.Duration
	DatabaseLockingMode() string
	DatabaseURL() url.URL
	DefaultChainID() *big.Int
	DefaultHTTPLimit() int64
	DefaultHTTPTimeout() models.Duration
	DefaultLogLevel() zapcore.Level
	Dev() bool
	ShutdownGracePeriod() time.Duration
	EthereumHTTPURL() *url.URL
	EthereumNodes() string
	EthereumSecondaryURLs() []url.URL
	EthereumURL() string
	ExplorerAccessKey() string
	ExplorerSecret() string
	ExplorerURL() *url.URL
	FMDefaultTransactionQueueDepth() uint32
	FMSimulateTransactions() bool
	GetAdvisoryLockIDConfiguredOrDefault() int64
	GetDatabaseDialectConfiguredOrDefault() dialects.DialectName
	HTTPServerWriteTimeout() time.Duration
	InsecureFastScrypt() bool
	JSONConsole() bool
	JobPipelineMaxRunDuration() time.Duration
	JobPipelineMaxSuccessfulRuns() uint64
	JobPipelineReaperInterval() time.Duration
	JobPipelineReaperThreshold() time.Duration
	JobPipelineResultWriteQueueDepth() uint64
	KeeperDefaultTransactionQueueDepth() uint32
	KeeperGasPriceBufferPercent() uint16
	KeeperGasTipCapBufferPercent() uint16
	KeeperBaseFeeBufferPercent() uint16
	KeeperMaximumGracePeriod() int64
	KeeperRegistryCheckGasOverhead() uint32
	KeeperRegistryPerformGasOverhead() uint32
	KeeperRegistryMaxPerformDataSize() uint32
	KeeperRegistrySyncInterval() time.Duration
	KeeperRegistrySyncUpkeepQueueSize() uint32
	KeeperTurnLookBack() int64
	KeyFile() string
	KeystorePassword() string
	LeaseLockDuration() time.Duration
	LeaseLockRefreshInterval() time.Duration
	LogFileDir() string
	LogLevel() zapcore.Level
	LogSQL() bool
	LogFileMaxSize() utils.FileSize
	LogFileMaxAge() int64
	LogFileMaxBackups() int64
	LogUnixTimestamps() bool
	MercuryCredentials(url string) (username, password string, err error)
	MigrateDatabase() bool
	ORMMaxIdleConns() int
	ORMMaxOpenConns() int
	Port() uint16
	PrometheusAuthToken() string
	PyroscopeAuthToken() string
	PyroscopeServerAddress() string
	PyroscopeEnvironment() string
	RPID() string
	RPOrigin() string
	ReaperExpiration() models.Duration
	RootDir() string
	SecureCookies() bool
	SentryDSN() string
	SentryDebug() bool
	SentryEnvironment() string
	SentryRelease() string
	SessionOptions() sessions.Options
	SessionTimeout() models.Duration
	SolanaNodes() string
	StarkNetNodes() string
	TLSCertPath() string
	TLSDir() string
	TLSHost() string
	TLSKeyPath() string
	TLSPort() uint16
	TLSRedirect() bool
	TelemetryIngressLogging() bool
	TelemetryIngressUniConn() bool
	TelemetryIngressServerPubKey() string
	TelemetryIngressURL() *url.URL
	TelemetryIngressBufferSize() uint
	TelemetryIngressMaxBatchSize() uint
	TelemetryIngressSendInterval() time.Duration
	TelemetryIngressSendTimeout() time.Duration
	TelemetryIngressUseBatchSend() bool
	TriggerFallbackDBPollInterval() time.Duration
	UnAuthenticatedRateLimit() int64
	UnAuthenticatedRateLimitPeriod() models.Duration
	VRFPassword() string

	OCR1Config
	OCR2Config

	P2PNetworking
	P2PV1Networking
	P2PV2Networking
}

// GlobalConfig holds global ENV overrides for EVM chains
// If set the global ENV will override everything
// The second bool indicates if it is set or not
type GlobalConfig interface {
	GlobalBalanceMonitorEnabled() (bool, bool)
	GlobalBlockEmissionIdleWarningThreshold() (time.Duration, bool)
	GlobalBlockHistoryEstimatorBatchSize() (uint32, bool)
	GlobalBlockHistoryEstimatorBlockDelay() (uint16, bool)
	GlobalBlockHistoryEstimatorBlockHistorySize() (uint16, bool)
	GlobalBlockHistoryEstimatorEIP1559FeeCapBufferBlocks() (uint16, bool)
	GlobalBlockHistoryEstimatorCheckInclusionBlocks() (uint16, bool)
	GlobalBlockHistoryEstimatorCheckInclusionPercentile() (uint16, bool)
	GlobalBlockHistoryEstimatorTransactionPercentile() (uint16, bool)
	GlobalChainType() (string, bool)
	GlobalEthTxReaperInterval() (time.Duration, bool)
	GlobalEthTxReaperThreshold() (time.Duration, bool)
	GlobalEthTxResendAfterThreshold() (time.Duration, bool)
	GlobalEvmEIP1559DynamicFees() (bool, bool)
	GlobalEvmFinalityDepth() (uint32, bool)
	GlobalEvmGasBumpPercent() (uint16, bool)
	GlobalEvmGasBumpThreshold() (uint64, bool)
	GlobalEvmGasBumpTxDepth() (uint16, bool)
	GlobalEvmGasBumpWei() (*assets.Wei, bool)
	GlobalEvmGasFeeCapDefault() (*assets.Wei, bool)
	GlobalEvmGasLimitDefault() (uint32, bool)
	GlobalEvmGasLimitMax() (uint32, bool)
	GlobalEvmGasLimitMultiplier() (float32, bool)
	GlobalEvmGasLimitTransfer() (uint32, bool)
	GlobalEvmGasLimitOCRJobType() (uint32, bool)
	GlobalEvmGasLimitDRJobType() (uint32, bool)
	GlobalEvmGasLimitVRFJobType() (uint32, bool)
	GlobalEvmGasLimitFMJobType() (uint32, bool)
	GlobalEvmGasLimitKeeperJobType() (uint32, bool)
	GlobalEvmGasPriceDefault() (*assets.Wei, bool)
	GlobalEvmGasTipCapDefault() (*assets.Wei, bool)
	GlobalEvmGasTipCapMinimum() (*assets.Wei, bool)
	GlobalEvmHeadTrackerHistoryDepth() (uint32, bool)
	GlobalEvmHeadTrackerMaxBufferSize() (uint32, bool)
	GlobalEvmHeadTrackerSamplingInterval() (time.Duration, bool)
	GlobalEvmLogBackfillBatchSize() (uint32, bool)
	GlobalEvmLogPollInterval() (time.Duration, bool)
	GlobalEvmLogKeepBlocksDepth() (uint32, bool)
	GlobalEvmMaxGasPriceWei() (*assets.Wei, bool)
	GlobalEvmMaxInFlightTransactions() (uint32, bool)
	GlobalEvmMaxQueuedTransactions() (uint64, bool)
	GlobalEvmMinGasPriceWei() (*assets.Wei, bool)
	GlobalEvmNonceAutoSync() (bool, bool)
	GlobalEvmUseForwarders() (bool, bool)
	GlobalEvmRPCDefaultBatchSize() (uint32, bool)
	GlobalFlagsContractAddress() (string, bool)
	GlobalGasEstimatorMode() (string, bool)
	GlobalLinkContractAddress() (string, bool)
	GlobalOCRContractConfirmations() (uint16, bool)
	GlobalOCRContractTransmitterTransmitTimeout() (time.Duration, bool)
	GlobalOCRDatabaseTimeout() (time.Duration, bool)
	GlobalOCRObservationGracePeriod() (time.Duration, bool)
	GlobalOCR2AutomationGasLimit() (uint32, bool)
	GlobalOperatorFactoryAddress() (string, bool)
	GlobalMinIncomingConfirmations() (uint32, bool)
	GlobalMinimumContractPayment() (*assets.Link, bool)
	GlobalNodeNoNewHeadsThreshold() (time.Duration, bool)
	GlobalNodePollFailureThreshold() (uint32, bool)
	GlobalNodePollInterval() (time.Duration, bool)
	GlobalNodeSelectionMode() (string, bool)
	GlobalNodeSyncThreshold() (uint32, bool)
}

type GeneralConfig interface {
	BasicConfig
	GlobalConfig
}

func ValidateDBURL(dbURI url.URL) error {
	if strings.Contains(dbURI.Redacted(), "_test") {
		return nil
	}

	// url params take priority if present, multiple params are ignored by postgres (it picks the first)
	q := dbURI.Query()
	// careful, this is a raw database password
	pw := q.Get("password")
	if pw == "" {
		// fallback to user info
		userInfo := dbURI.User
		if userInfo == nil {
			return errors.Errorf("DB URL must be authenticated; plaintext URLs are not allowed")
		}
		var pwSet bool
		pw, pwSet = userInfo.Password()
		if !pwSet {
			return errors.Errorf("DB URL must be authenticated; password is required")
		}
	}

	return utils.VerifyPasswordComplexity(pw)
}

type DatabaseBackupMode string

var (
	DatabaseBackupModeNone DatabaseBackupMode = "none"
	DatabaseBackupModeLite DatabaseBackupMode = "lite"
	DatabaseBackupModeFull DatabaseBackupMode = "full"
)
