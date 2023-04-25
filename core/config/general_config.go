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

	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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
	CosmosEnabled() bool
	SolanaEnabled() bool
	StarkNetEnabled() bool
}

type LogfFn func(string, ...any)

type BasicConfig interface {
	Validate() error
	LogConfiguration(log LogfFn)
	SetLogLevel(lvl zapcore.Level) error
	SetLogSQL(logSQL bool)
	SetPasswords(keystore, vrf *string)

	FeatureFlags
	audit.Config

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
	EthereumSecondaryURLs() []url.URL
	EthereumURL() string
	ExplorerAccessKey() string
	ExplorerSecret() string
	ExplorerURL() *url.URL
	FMDefaultTransactionQueueDepth() uint32
	FMSimulateTransactions() bool
	GetDatabaseDialectConfiguredOrDefault() dialects.DialectName
	HTTPServerWriteTimeout() time.Duration
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

	// Insecure config
	DevWebServer() bool
	InsecureFastScrypt() bool
	OCRDevelopmentMode() bool
	DisableRateLimiting() bool
	InfiniteDepthQueries() bool

	OCR1Config
	OCR2Config

	P2PNetworking
	P2PV1Networking
	P2PV2Networking
}

type GeneralConfig interface {
	BasicConfig
	ValidateDB() error
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
