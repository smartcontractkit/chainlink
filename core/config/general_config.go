package config

import (
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"

	"github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/config/envvar"
	"github.com/smartcontractkit/chainlink/core/config/parse"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/logger/audit"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --quiet --name GeneralConfig --output ./mocks/ --case=underscore

// nolint
var (
	ErrEnvUnset   = errors.New("env var unset")
	ErrEnvInvalid = errors.New("env var invalid")

	configFileNotFoundError = reflect.TypeOf(viper.ConfigFileNotFoundError{})
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

// generalConfig holds parameters used by the application which can be overridden by
// setting environment variables.
//
// If you add an entry here which does not contain sensitive information, you
// should also update presenters.ConfigWhitelist and cmd_test.TestClient_RunNodeShowsEnv.
type generalConfig struct {
	lggr             logger.Logger
	viper            *viper.Viper
	randomP2PPort    uint16
	randomP2PPortMtx sync.RWMutex
	dialect          dialects.DialectName
	advisoryLockID   int64
	logLevel         zapcore.Level
	defaultLogLevel  zapcore.Level
	logSQL           bool
	logMutex         sync.RWMutex
	genAppID         sync.Once
	appID            uuid.UUID

	passwordKeystore, passwordVRF string
	passwordMu                    sync.RWMutex // passwords are set after initialization
}

// NewGeneralConfig returns the config with the environment variables set to their
// respective fields, or their defaults if environment variables are not set.
// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func NewGeneralConfig(lggr logger.Logger) GeneralConfig {
	v := viper.New()
	c := newGeneralConfigWithViper(v, lggr.Named("GeneralConfig"))
	c.dialect = dialects.Postgres
	return c
}

func newGeneralConfigWithViper(v *viper.Viper, lggr logger.Logger) (config *generalConfig) {
	schemaT := reflect.TypeOf(envvar.ConfigSchema{})
	for index := 0; index < schemaT.NumField(); index++ {
		item := schemaT.FieldByIndex([]int{index})
		name := item.Tag.Get("env")
		def, exists := item.Tag.Lookup("default")
		if exists {
			v.SetDefault(name, def)
		}
		_ = v.BindEnv(name, name)
	}

	config = &generalConfig{
		lggr:  lggr,
		viper: v,
	}

	if err := utils.EnsureDirAndMaxPerms(config.RootDir(), os.FileMode(0700)); err != nil {
		lggr.Fatalf(`Error creating root directory "%s": %+v`, config.RootDir(), err)
	}

	v.SetConfigName("chainlink")
	v.AddConfigPath(config.RootDir())
	err := v.ReadInConfig()
	if err != nil && reflect.TypeOf(err) != configFileNotFoundError {
		lggr.Warnf("Unable to load config file: %v\n", err)
	}

	ll, invalid := envvar.LogLevel.Parse()
	if invalid != "" {
		lggr.Error(invalid)
	}
	config.defaultLogLevel = ll

	config.logLevel = config.defaultLogLevel
	config.logSQL = viper.GetBool(envvar.Name("LogSQL"))

	return
}

func (c *generalConfig) LogConfiguration(log LogFn) {
	log("Environment variables\n", NewConfigPrinter(c))
}

// Validate performs basic sanity checks on config and returns error if any
// misconfiguration would be fatal to the application
func (c *generalConfig) Validate() error {
	if c.P2PAnnouncePort() != 0 && c.P2PAnnounceIP() == nil {
		return errors.Errorf("P2P_ANNOUNCE_PORT was given as %v but P2P_ANNOUNCE_IP was unset. You must also set P2P_ANNOUNCE_IP if P2P_ANNOUNCE_PORT is set", c.P2PAnnouncePort())
	}

	if _, exists := os.LookupEnv("MINIMUM_CONTRACT_PAYMENT"); exists {
		return errors.Errorf("MINIMUM_CONTRACT_PAYMENT is deprecated, use MINIMUM_CONTRACT_PAYMENT_LINK_JUELS instead")
	}

	if _, exists := os.LookupEnv("ETH_DISABLED"); exists {
		c.lggr.Error(`ETH_DISABLED is deprecated.

This will become a fatal error in a future release. Please switch to using one of the two options below instead:

- EVM_ENABLED=false - set this if you wish to completely disable all EVM chains and jobs and prevent them from ever loading (this is probably the one you want).
- EVM_RPC_ENABLED=false - set this if you wish to load all EVM chains and jobs, but prevent any RPC calls to the eth node (the old behaviour).
`)
	}
	if _, exists := os.LookupEnv("EVM_DISABLED"); exists {
		c.lggr.Error(`EVM_DISABLED is deprecated and superceded by EVM_ENABLED.

This will become a fatal error in a future release. Please use the following instead to disable EVM chains:

EVM_ENABLED=false
`)
	}

	if _, err := c.OCRKeyBundleID(); errors.Is(errors.Cause(err), ErrEnvInvalid) {
		return err
	}
	if _, err := c.OCRTransmitterAddress(); errors.Is(errors.Cause(err), ErrEnvInvalid) {
		return err
	}
	if peers, err := c.P2PBootstrapPeers(); err == nil {
		for i := range peers {
			if _, err := multiaddr.NewMultiaddr(peers[i]); err != nil {
				return errors.Errorf("p2p bootstrap peer %d is invalid: err %v", i, err)
			}
		}
	}
	if ct, set := c.GlobalChainType(); set && !ChainType(ct).IsValid() {
		return errors.Errorf("CHAIN_TYPE is invalid: %s", ct)
	}

	if c.EthereumURL() == "" {
		if c.EthereumHTTPURL() != nil {
			c.lggr.Warn("ETH_HTTP_URL has no effect when ETH_URL is not set")
		}
		if len(c.EthereumSecondaryURLs()) > 0 {
			c.lggr.Warn("ETH_SECONDARY_URL/ETH_SECONDARY_URLS have no effect when ETH_URL is not set")
		}
	} else if c.EthereumNodes() != "" {
		return errors.Errorf("It is not permitted to set both ETH_URL (got %s) and EVM_NODES (got %s). Please set either one or the other", c.EthereumURL(), c.EthereumNodes())
	}
	// Warn on legacy OCR env vars
	if c.ocrDHTLookupInterval() != 0 {
		c.lggr.Error("OCR_DHT_LOOKUP_INTERVAL is deprecated, use P2P_DHT_LOOKUP_INTERVAL instead")
	}
	if c.ocrBootstrapCheckInterval() != 0 {
		c.lggr.Error("OCR_BOOTSTRAP_CHECK_INTERVAL is deprecated, use P2P_BOOTSTRAP_CHECK_INTERVAL instead")
	}
	if c.ocrIncomingMessageBufferSize() != 0 {
		c.lggr.Error("OCR_INCOMING_MESSAGE_BUFFER_SIZE is deprecated, use P2P_INCOMING_MESSAGE_BUFFER_SIZE instead")
	}
	if c.ocrOutgoingMessageBufferSize() != 0 {
		c.lggr.Error("OCR_OUTGOING_MESSAGE_BUFFER_SIZE is deprecated, use P2P_OUTGOING_MESSAGE_BUFFER_SIZE instead")
	}
	if c.ocrNewStreamTimeout() != 0 {
		c.lggr.Error("OCR_NEW_STREAM_TIMEOUT is deprecated, use P2P_NEW_STREAM_TIMEOUT instead")
	}

	switch c.DatabaseLockingMode() {
	case "dual", "lease", "advisorylock", "none":
	default:
		return errors.Errorf("unrecognised value for DATABASE_LOCKING_MODE: %s (valid options are 'dual', 'lease', 'advisorylock' or 'none')", c.DatabaseLockingMode())
	}

	if c.LeaseLockRefreshInterval() > c.LeaseLockDuration()/2 {
		return errors.Errorf("LEASE_LOCK_REFRESH_INTERVAL must be less than or equal to half of LEASE_LOCK_DURATION (got LEASE_LOCK_REFRESH_INTERVAL=%s, LEASE_LOCK_DURATION=%s)", c.LeaseLockRefreshInterval().String(), c.LeaseLockDuration().String())
	}

	if c.viper.GetString(envvar.Name("LogFileDir")) != "" && c.LogFileMaxSize() <= 0 {
		c.lggr.Warn("LOG_FILE_DIR is ignored and has no effect when LOG_FILE_MAX_SIZE is not set to a value greater than zero")
	}

	{
		str := os.Getenv("SKIP_DATABASE_PASSWORD_COMPLEXITY_CHECK")
		var skipDatabasePasswordComplexityCheck bool
		if str != "" {
			var err error
			skipDatabasePasswordComplexityCheck, err = strconv.ParseBool(str)
			if err != nil {
				return errors.Errorf("SKIP_DATABASE_PASSWORD_COMPLEXITY_CHECK has invalid value for bool: %s", str)
			}
		}
		if !(c.Dev() || skipDatabasePasswordComplexityCheck) {
			if err := ValidateDBURL(c.DatabaseURL()); err != nil {
				// TODO: Make this a hard error in some future version of Chainlink > 1.4.x
				c.lggr.Errorf("DEPRECATION WARNING: Database has missing or insufficiently complex password: %s.\nDatabase should be secured by a password matching the following complexity requirements:\n%s\nThis error will PREVENT BOOT in a future version of Chainlink. To bypass this check at your own risk, you may set SKIP_DATABASE_PASSWORD_COMPLEXITY_CHECK=true\n\n", err, utils.PasswordComplexityRequirements)
			}
		}
	}

	if str := c.viper.GetString("MIN_OUTGOING_CONFIRMATIONS"); str != "" {
		c.lggr.Errorf("MIN_OUTGOING_CONFIRMATIONS has been removed and no longer has any effect. ETH_FINALITY_DEPTH is now used as the default for ethtx confirmations instead. You may override this on a per-task basis by setting `minConfirmations` e.g. `foo [type=ethtx minConfirmations=%s ...]`", str)
	}

	return nil
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

func (c *generalConfig) GetAdvisoryLockIDConfiguredOrDefault() int64 {
	return c.advisoryLockID
}

func (c *generalConfig) GetDatabaseDialectConfiguredOrDefault() dialects.DialectName {
	return c.dialect
}

// AllowOrigins returns the CORS hosts used by the frontend.
func (c *generalConfig) AllowOrigins() string {
	return c.viper.GetString(envvar.Name("AllowOrigins"))
}

func (c *generalConfig) AppID() uuid.UUID {
	c.genAppID.Do(func() {
		c.appID = uuid.NewV4()
	})
	return c.appID
}

func (c *generalConfig) AuditLoggerEnabled() bool {
	return c.viper.GetBool(envvar.Name("AuditLoggerEnabled"))
}

func (c *generalConfig) AuditLoggerForwardToUrl() (models.URL, error) {
	url, err := models.ParseURL(c.viper.GetString(envvar.Name("AuditLoggerForwardToUrl")))
	if err != nil {
		return models.URL{}, err
	}
	return *url, nil
}

func (c *generalConfig) AuditLoggerEnvironment() string {
	if c.Dev() {
		return "develop"
	}
	return "production"
}

func (c *generalConfig) AuditLoggerJsonWrapperKey() string {
	return c.viper.GetString(envvar.Name("AuditLoggerJsonWrapperKey"))
}

func (c *generalConfig) AuditLoggerHeaders() (audit.ServiceHeaders, error) {
	sh, invalid := audit.AuditLoggerHeaders.Parse()
	if invalid != "" {
		return nil, errors.New(invalid)
	}
	return sh, nil
}

// AuthenticatedRateLimit defines the threshold to which authenticated requests
// get limited. More than this many requests per AuthenticatedRateLimitPeriod will be rejected.
func (c *generalConfig) AuthenticatedRateLimit() int64 {
	return c.viper.GetInt64(envvar.Name("AuthenticatedRateLimit"))
}

// AuthenticatedRateLimitPeriod defines the period to which authenticated requests get limited
func (c *generalConfig) AuthenticatedRateLimitPeriod() models.Duration {
	return models.MustMakeDuration(getEnvWithFallback(c, envvar.AuthenticatedRateLimitPeriod))
}

func (c *generalConfig) AutoPprofEnabled() bool {
	return c.viper.GetBool(envvar.Name("AutoPprofEnabled"))
}

func (c *generalConfig) AutoPprofProfileRoot() string {
	root := c.viper.GetString(envvar.Name("AutoPprofProfileRoot"))
	if root == "" {
		return filepath.Join(c.RootDir(), "pprof")
	}
	return root
}

func (c *generalConfig) AutoPprofPollInterval() models.Duration {
	return models.MustMakeDuration(getEnvWithFallback(c, envvar.AutoPprofPollInterval))
}

func (c *generalConfig) AutoPprofGatherDuration() models.Duration {
	return models.MustMakeDuration(getEnvWithFallback(c, envvar.AutoPprofGatherDuration))
}

func (c *generalConfig) AutoPprofGatherTraceDuration() models.Duration {
	return models.MustMakeDuration(getEnvWithFallback(c, envvar.AutoPprofGatherTraceDuration))
}

func (c *generalConfig) AutoPprofMaxProfileSize() utils.FileSize {
	return getEnvWithFallback(c, envvar.New("AutoPprofMaxProfileSize", parse.FileSize))
}

func (c *generalConfig) AutoPprofCPUProfileRate() int {
	return c.viper.GetInt(envvar.Name("AutoPprofCPUProfileRate"))
}

func (c *generalConfig) AutoPprofMemProfileRate() int {
	return c.viper.GetInt(envvar.Name("AutoPprofMemProfileRate"))
}

func (c *generalConfig) AutoPprofBlockProfileRate() int {
	return c.viper.GetInt(envvar.Name("AutoPprofBlockProfileRate"))
}

func (c *generalConfig) AutoPprofMutexProfileFraction() int {
	return c.viper.GetInt(envvar.Name("AutoPprofMutexProfileFraction"))
}

func (c *generalConfig) AutoPprofMemThreshold() utils.FileSize {
	return getEnvWithFallback(c, envvar.New("AutoPprofMemThreshold", parse.FileSize))
}

func (c *generalConfig) AutoPprofGoroutineThreshold() int {
	return c.viper.GetInt(envvar.Name("AutoPprofGoroutineThreshold"))
}

// PyroscopeAuthToken specifies the Auth Token used to send profiling info to Pyroscope
func (c *generalConfig) PyroscopeAuthToken() string {
	return c.viper.GetString(envvar.Name("PyroscopeAuthToken"))
}

// PyroscopeServerAddress specifies the Server Address where the Pyroscope instance lives
func (c *generalConfig) PyroscopeServerAddress() string {
	return c.viper.GetString(envvar.Name("PyroscopeServerAddress"))
}

// PyroscopeEnvironment specifies the Environment where the Pyroscope logs will be categorized
func (c *generalConfig) PyroscopeEnvironment() string {
	return c.viper.GetString(envvar.Name("PyroscopeEnvironment"))
}

// BlockBackfillDepth specifies the number of blocks before the current HEAD that the
// log broadcaster will try to re-consume logs from
func (c *generalConfig) BlockBackfillDepth() uint64 {
	return getEnvWithFallback(c, envvar.BlockBackfillDepth)
}

// BlockBackfillSkip enables skipping of very long log backfills
func (c *generalConfig) BlockBackfillSkip() bool {
	return getEnvWithFallback(c, envvar.NewBool("BlockBackfillSkip"))
}

// BridgeResponseURL represents the URL for bridges to send a response to.
func (c *generalConfig) BridgeResponseURL() *url.URL {
	return getEnvWithFallback(c, envvar.New("BridgeResponseURL", url.Parse))
}

// BridgeCacheTTL represents the max acceptable duration for a cached bridge value to be used in case of intermittent failure.
func (c *generalConfig) BridgeCacheTTL() time.Duration {
	return getEnvWithFallback(c, envvar.NewDuration("BridgeCacheTTL"))
}

// FeatureUICSAKeys enables the CSA Keys UI Feature.
func (c *generalConfig) FeatureUICSAKeys() bool {
	return getEnvWithFallback(c, envvar.NewBool("FeatureUICSAKeys"))
}

func (c *generalConfig) DatabaseListenerMinReconnectInterval() time.Duration {
	return getEnvWithFallback(c, envvar.NewDuration("DatabaseListenerMinReconnectInterval"))
}

func (c *generalConfig) DatabaseListenerMaxReconnectDuration() time.Duration {
	return getEnvWithFallback(c, envvar.NewDuration("DatabaseListenerMaxReconnectDuration"))
}

var DatabaseBackupModeEnvVar = envvar.New("DatabaseBackupMode", parseDatabaseBackupMode)

// DatabaseBackupMode sets the database backup mode
func (c *generalConfig) DatabaseBackupMode() DatabaseBackupMode {
	return getEnvWithFallback(c, DatabaseBackupModeEnvVar)
}

// DatabaseBackupFrequency turns on the periodic database backup if set to a positive value
// DatabaseBackupMode must be then set to a value other than "none"
func (c *generalConfig) DatabaseBackupFrequency() time.Duration {
	return getEnvWithFallback(c, envvar.NewDuration("DatabaseBackupFrequency"))
}

// DatabaseBackupURL configures the URL for the database to backup, if it's to be different from the main on
func (c *generalConfig) DatabaseBackupURL() *url.URL {
	s := c.viper.GetString(envvar.Name("DatabaseBackupURL"))
	if s == "" {
		return nil
	}
	uri, err := url.Parse(s)
	if err != nil {
		c.lggr.Errorf("Invalid database backup url %s", s)
		return nil
	}
	return uri
}

// DatabaseBackupOnVersionUpgrade controls whether an automatic backup will be
// taken before migrations are run, if the node version has been bumped
func (c *generalConfig) DatabaseBackupOnVersionUpgrade() bool {
	return getEnvWithFallback(c, envvar.NewBool("DatabaseBackupOnVersionUpgrade"))
}

// DatabaseBackupDir configures the directory for saving the backup file, if it's to be different from default one located in the RootDir
func (c *generalConfig) DatabaseBackupDir() string {
	return c.viper.GetString(envvar.Name("DatabaseBackupDir"))
}

func (c *generalConfig) DatabaseDefaultIdleInTxSessionTimeout() time.Duration {
	return pg.DefaultIdleInTxSessionTimeout
}

func (c *generalConfig) DatabaseDefaultLockTimeout() time.Duration {
	return pg.DefaultLockTimeout
}

func (c *generalConfig) DatabaseDefaultQueryTimeout() time.Duration {
	return pg.DefaultQueryTimeout
}

// DatabaseURL configures the URL for chainlink to connect to. This must be
// a properly formatted URL, with a valid scheme (postgres://)
func (c *generalConfig) DatabaseURL() url.URL {
	s := c.viper.GetString(envvar.Name("DatabaseURL"))
	uri, err := url.Parse(s)
	if err != nil {
		c.lggr.Error("invalid database url %s", s)
		return url.URL{}
	}
	if uri.String() == "" {
		return *uri
	}
	static.SetConsumerName(uri, "Default", nil)
	return *uri
}

// MigrateDatabase determines whether the database will be automatically
// migrated on application startup if set to true
func (c *generalConfig) MigrateDatabase() bool {
	return c.viper.GetBool(envvar.Name("MigrateDatabase"))
}

// DefaultHTTPLimit defines the size limit for HTTP requests and responses
func (c *generalConfig) DefaultHTTPLimit() int64 {
	return c.viper.GetInt64(envvar.Name("DefaultHTTPLimit"))
}

// DefaultHTTPTimeout defines the default timeout for http requests
func (c *generalConfig) DefaultHTTPTimeout() models.Duration {
	return models.MustMakeDuration(getEnvWithFallback(c, envvar.NewDuration("DefaultHTTPTimeout")))
}

// Dev configures "development" mode for chainlink.
func (c *generalConfig) Dev() bool {
	return c.viper.GetBool(envvar.Name("Dev"))
}

// ShutdownGracePeriod is the maximum duration of graceful application shutdown.
// If exceeded, it will try closing DB lock and connection and exit immediately to avoid SIGKILL.
func (c *generalConfig) ShutdownGracePeriod() time.Duration {
	return getEnvWithFallback(c, envvar.NewDuration("ShutdownGracePeriod"))
}

// FeatureExternalInitiators enables the External Initiator feature.
func (c *generalConfig) FeatureExternalInitiators() bool {
	return c.viper.GetBool(envvar.Name("FeatureExternalInitiators"))
}

// FeatureFeedsManager enables the feeds manager
func (c *generalConfig) FeatureFeedsManager() bool {
	return c.viper.GetBool(envvar.Name("FeatureFeedsManager"))
}

func (c *generalConfig) FeatureLogPoller() bool {
	return c.viper.GetBool(envvar.Name("FeatureLogPoller"))
}

// FeatureOffchainReporting enables the OCR job type.
func (c *generalConfig) FeatureOffchainReporting() bool {
	return getEnvWithFallback(c, envvar.NewBool("FeatureOffchainReporting"))
}

// FeatureOffchainReporting2 enables the OCR2 job type.
func (c *generalConfig) FeatureOffchainReporting2() bool {
	return getEnvWithFallback(c, envvar.NewBool("FeatureOffchainReporting2"))
}

// FMDefaultTransactionQueueDepth controls the queue size for DropOldestStrategy in Flux Monitor
// Set to 0 to use SendEvery strategy instead
func (c *generalConfig) FMDefaultTransactionQueueDepth() uint32 {
	return c.viper.GetUint32(envvar.Name("FMDefaultTransactionQueueDepth"))
}

// FMSimulateTransactions enables using eth_call transaction simulation before
// sending when set to true
func (c *generalConfig) FMSimulateTransactions() bool {
	return c.viper.GetBool(envvar.Name("FMSimulateTransactions"))
}

// EthereumURL represents the URL of the Ethereum node to connect Chainlink to.
func (c *generalConfig) EthereumURL() string {
	return c.viper.GetString(envvar.Name("EthereumURL"))
}

// EthereumHTTPURL is an optional but recommended url that points to the HTTP port of the primary node
func (c *generalConfig) EthereumHTTPURL() (uri *url.URL) {
	urlStr := c.viper.GetString(envvar.Name("EthereumHTTPURL"))
	if urlStr == "" {
		return nil
	}
	var err error
	uri, err = url.Parse(urlStr)
	if err != nil || !(uri.Scheme == "http" || uri.Scheme == "https") {
		c.lggr.Fatalf("Invalid Ethereum HTTP URL: %s, got error: %s", urlStr, err)
	}
	return
}

// EthereumNodes is a hack to allow node operators to give a JSON string that
// sets up multiple nodes
func (c *generalConfig) EthereumNodes() string {
	return c.viper.GetString(envvar.Name("EthereumNodes"))
}

// EthereumSecondaryURLs is an optional backup RPC URL
// Must be http(s) format
// If specified, transactions will also be broadcast to this ethereum node
func (c *generalConfig) EthereumSecondaryURLs() []url.URL {
	oldConfig := c.viper.GetString(envvar.Name("EthereumSecondaryURL"))
	newConfig := c.viper.GetString(envvar.Name("EthereumSecondaryURLs"))

	config := ""
	if newConfig != "" {
		config = newConfig
	} else if oldConfig != "" {
		config = oldConfig
	}

	urlStrings := regexp.MustCompile(`\s*[;,]\s*`).Split(config, -1)
	var urls []url.URL
	for _, urlString := range urlStrings {
		if urlString == "" {
			continue
		}
		url, err := url.Parse(urlString)
		if err != nil {
			c.lggr.Fatalf("Invalid Secondary Ethereum URL: %s, got error: %v", urlString, err)
		}
		urls = append(urls, *url)
	}

	return urls
}

// EVMRPCEnabled if false prevents any calls to any EVM-based chain RPC node
func (c *generalConfig) EVMRPCEnabled() bool {
	if ethDisabled, exists := os.LookupEnv("ETH_DISABLED"); exists {
		res, err := parse.Bool(ethDisabled)
		if err == nil {
			return !res.(bool)
		}
		c.lggr.Warnw("Failed to parse value for ETH_DISABLED", "err", err)
	}
	rpcEnabled := c.viper.GetBool(envvar.Name("EVMRPCEnabled"))
	return rpcEnabled
}

// EVMEnabled allows EVM chains to be used
func (c *generalConfig) EVMEnabled() bool {
	if evmDisabled, exists := os.LookupEnv("EVM_DISABLED"); exists {
		res, err := parse.Bool(evmDisabled)
		if err == nil {
			return res.(bool)
		}
		c.lggr.Warnw("Failed to parse value for EVM_DISABLED", "err", err)
	}
	return c.viper.GetBool(envvar.Name("EVMEnabled"))
}

// SolanaEnabled allows Solana to be used
func (c *generalConfig) SolanaEnabled() bool {
	return c.viper.GetBool(envvar.Name("SolanaEnabled"))
}

// StarkNetEnabled allows StarkNet to be used
func (c *generalConfig) StarkNetEnabled() bool {
	return c.viper.GetBool(envvar.Name("StarknetEnabled"))
}

// P2PEnabled controls whether Chainlink will run as a P2P peer for OCR protocol
func (c *generalConfig) P2PEnabled() bool {
	// We need p2p networking if either ocr1 or ocr2 is enabled
	return c.P2PListenPort() > 0 || c.FeatureOffchainReporting() || c.FeatureOffchainReporting2()
}

// InsecureFastScrypt causes all key stores to encrypt using "fast" scrypt params instead
// This is insecure and only useful for local testing. DO NOT SET THIS IN PRODUCTION
func (c *generalConfig) InsecureFastScrypt() bool {
	return c.viper.GetBool(envvar.Name("InsecureFastScrypt"))
}

func (c *generalConfig) TriggerFallbackDBPollInterval() time.Duration {
	return getEnvWithFallback(c, envvar.NewDuration("TriggerFallbackDBPollInterval"))
}

// JobPipelineMaxRunDuration is the maximum time that a job run may take
func (c *generalConfig) JobPipelineMaxRunDuration() time.Duration {
	return getEnvWithFallback(c, envvar.JobPipelineMaxRunDuration)
}

func (c *generalConfig) JobPipelineResultWriteQueueDepth() uint64 {
	return getEnvWithFallback(c, envvar.JobPipelineResultWriteQueueDepth)
}

func (c *generalConfig) JobPipelineMaxSuccessfulRuns() uint64 {
	return getEnvWithFallback(c, envvar.JobPipelineMaxSuccessfulRuns)
}

func (c *generalConfig) JobPipelineReaperInterval() time.Duration {
	return getEnvWithFallback(c, envvar.JobPipelineReaperInterval)
}

func (c *generalConfig) JobPipelineReaperThreshold() time.Duration {
	return getEnvWithFallback(c, envvar.JobPipelineReaperThreshold)
}

// KeeperRegistryCheckGasOverhead is the amount of extra gas to provide checkUpkeep() calls
// to account for the gas consumed by the keeper registry
func (c *generalConfig) KeeperRegistryCheckGasOverhead() uint32 {
	return getEnvWithFallback(c, envvar.KeeperRegistryCheckGasOverhead)
}

// KeeperRegistryPerformGasOverhead is the amount of extra gas to provide performUpkeep() calls
// to account for the gas consumed by the keeper registry
func (c *generalConfig) KeeperRegistryPerformGasOverhead() uint32 {
	return getEnvWithFallback(c, envvar.KeeperRegistryPerformGasOverhead)
}

// KeeperRegistryMaxPerformDataSize is the max perform data size we allow in our pipeline for an
// upkeep to be performed with
func (c *generalConfig) KeeperRegistryMaxPerformDataSize() uint32 {
	return getEnvWithFallback(c, envvar.KeeperRegistryMaxPerformDataSize)
}

// KeeperDefaultTransactionQueueDepth controls the queue size for DropOldestStrategy in Keeper
// Set to 0 to use SendEvery strategy instead
func (c *generalConfig) KeeperDefaultTransactionQueueDepth() uint32 {
	return c.viper.GetUint32(envvar.Name("KeeperDefaultTransactionQueueDepth"))
}

// KeeperGasPriceBufferPercent adds the specified percentage to the gas price
// used for checking whether to perform an upkeep. Only applies in legacy mode.
func (c *generalConfig) KeeperGasPriceBufferPercent() uint16 {
	return uint16(c.viper.GetUint32(envvar.Name("KeeperGasPriceBufferPercent")))
}

// KeeperGasTipCapBufferPercent adds the specified percentage to the gas price
// used for checking whether to perform an upkeep. Only applies in EIP-1559 mode.
func (c *generalConfig) KeeperGasTipCapBufferPercent() uint16 {
	return uint16(c.viper.GetUint32(envvar.Name("KeeperGasTipCapBufferPercent")))
}

// KeeperBaseFeeBufferPercent adds the specified percentage to the base fee
// used for checking whether to perform an upkeep. Only applies in EIP-1559 mode.
func (c *generalConfig) KeeperBaseFeeBufferPercent() uint16 {
	return uint16(c.viper.GetUint32(envvar.Name("KeeperBaseFeeBufferPercent")))
}

// KeeperRegistrySyncInterval is the interval in which the RegistrySynchronizer performs a full
// sync of the keeper registry contract it is tracking *after* the most recent update triggered
// by an on-chain log.
func (c *generalConfig) KeeperRegistrySyncInterval() time.Duration {
	return getEnvWithFallback(c, envvar.KeeperRegistrySyncInterval)
}

// KeeperMaximumGracePeriod is the maximum number of blocks that a keeper will wait after performing
// an upkeep before it resumes checking that upkeep
func (c *generalConfig) KeeperMaximumGracePeriod() int64 {
	return c.viper.GetInt64(envvar.Name("KeeperMaximumGracePeriod"))
}

// KeeperRegistrySyncUpkeepQueueSize represents the maximum number of upkeeps that can be synced in parallel
func (c *generalConfig) KeeperRegistrySyncUpkeepQueueSize() uint32 {
	return getEnvWithFallback(c, envvar.KeeperRegistrySyncUpkeepQueueSize)
}

// KeeperTurnLookBack represents the number of blocks in the past to loo back when getting block for turn
func (c *generalConfig) KeeperTurnLookBack() int64 {
	return c.viper.GetInt64(envvar.Name("KeeperTurnLookBack"))
}

// JSONConsole when set to true causes logging to be made in JSON format
// If set to false, logs in console format
func (c *generalConfig) JSONConsole() bool {
	return getEnvWithFallback(c, envvar.JSONConsole)
}

// ExplorerURL returns the websocket URL for this node to push stats to, or nil.
func (c *generalConfig) ExplorerURL() *url.URL {
	return getEnvWithFallback(c, envvar.New("ExplorerURL", url.Parse))
}

// ExplorerAccessKey returns the access key for authenticating with explorer
func (c *generalConfig) ExplorerAccessKey() string {
	return c.viper.GetString(envvar.Name("ExplorerAccessKey"))
}

// ExplorerSecret returns the secret for authenticating with explorer
func (c *generalConfig) ExplorerSecret() string {
	return c.viper.GetString(envvar.Name("ExplorerSecret"))
}

// SolanaNodes is a hack to allow node operators to give a JSON string that
// sets up multiple nodes
func (c *generalConfig) SolanaNodes() string {
	return c.viper.GetString(envvar.Name("SolanaNodes"))
}

// StarkNetNodes is a hack to allow node operators to give a JSON string that
// sets up multiple nodes
func (c *generalConfig) StarkNetNodes() string {
	return c.viper.GetString(envvar.Name("StarknetNodes"))
}

// TelemetryIngressURL returns the WSRPC URL for this node to push telemetry to, or nil.
func (c *generalConfig) TelemetryIngressURL() *url.URL {
	return getEnvWithFallback(c, envvar.New("TelemetryIngressURL", url.Parse))
}

// TelemetryIngressServerPubKey returns the public key to authenticate the telemetry ingress server
func (c *generalConfig) TelemetryIngressServerPubKey() string {
	return c.viper.GetString(envvar.Name("TelemetryIngressServerPubKey"))
}

// TelemetryIngressBufferSize is the number of telemetry messages to buffer before dropping new ones
func (c *generalConfig) TelemetryIngressBufferSize() uint {
	return c.viper.GetUint(envvar.Name("TelemetryIngressBufferSize"))
}

// TelemetryIngressMaxBatchSize is the maximum number of messages to batch into one telemetry request
func (c *generalConfig) TelemetryIngressMaxBatchSize() uint {
	return c.viper.GetUint(envvar.Name("TelemetryIngressMaxBatchSize"))
}

// TelemetryIngressSendInterval is the cadence on which batched telemetry is sent to the ingress server
func (c *generalConfig) TelemetryIngressSendInterval() time.Duration {
	return c.getDuration("TelemetryIngressSendInterval")
}

// TelemetryIngressSendTimeout is the max duration to wait for the request to complete when sending batch telemetry
func (c *generalConfig) TelemetryIngressSendTimeout() time.Duration {
	return c.getDuration("TelemetryIngressSendTimeout")
}

// TelemetryIngressUseBatchSend toggles sending telemetry using the batch client to the ingress server
func (c *generalConfig) TelemetryIngressUseBatchSend() bool {
	return c.viper.GetBool(envvar.Name("TelemetryIngressUseBatchSend"))
}

// TelemetryIngressLogging toggles very verbose logging of raw telemetry messages for the TelemetryIngressClient
func (c *generalConfig) TelemetryIngressLogging() bool {
	return getEnvWithFallback(c, envvar.NewBool("TelemetryIngressLogging"))
}

// TelemetryIngressUniconn toggles which ws connection style is used.
func (c *generalConfig) TelemetryIngressUniConn() bool {
	return c.getWithFallback("TelemetryIngressUniConn", parse.Bool).(bool)
}

func (c *generalConfig) ORMMaxOpenConns() int {
	return int(getEnvWithFallback(c, envvar.NewUint16("ORMMaxOpenConns")))
}

func (c *generalConfig) ORMMaxIdleConns() int {
	return int(getEnvWithFallback(c, envvar.NewUint16("ORMMaxIdleConns")))
}

// LogLevel represents the maximum level of log messages to output.
func (c *generalConfig) LogLevel() zapcore.Level {
	c.logMutex.RLock()
	defer c.logMutex.RUnlock()
	return c.logLevel
}

// DefaultLogLevel returns default log level.
func (c *generalConfig) DefaultLogLevel() zapcore.Level {
	return c.defaultLogLevel
}

// SetLogLevel saves a runtime value for the default logger level
func (c *generalConfig) SetLogLevel(lvl zapcore.Level) error {
	c.logMutex.Lock()
	defer c.logMutex.Unlock()
	c.logLevel = lvl
	return nil
}

// LogFileMaxSize configures disk preservation of logs max size (in megabytes) before file rotation.
func (c *generalConfig) LogFileMaxSize() utils.FileSize {
	return getEnvWithFallback(c, envvar.LogFileMaxSize)
}

// LogFileMaxAge configures disk preservation of logs max age (in days) before file rotation.
func (c *generalConfig) LogFileMaxAge() int64 {
	return getEnvWithFallback(c, envvar.LogFileMaxAge)
}

// LogFileMaxBackups configures disk preservation of the max amount of old log files to retain.
// If this is set to 0, the node will retain all old log files instead.
func (c *generalConfig) LogFileMaxBackups() int64 {
	return getEnvWithFallback(c, envvar.LogFileMaxBackups)
}

// LogSQL tells chainlink to log all SQL statements made using the default logger
func (c *generalConfig) LogSQL() bool {
	c.logMutex.RLock()
	defer c.logMutex.RUnlock()
	return c.logSQL
}

// SetLogSQL saves a runtime value for enabling/disabling logging all SQL statements on the default logger
func (c *generalConfig) SetLogSQL(logSQL bool) {
	c.logMutex.Lock()
	defer c.logMutex.Unlock()
	c.logSQL = logSQL
}

// LogUnixTimestamps if set to true will log with timestamp in unix format, otherwise uses ISO8601
func (c *generalConfig) LogUnixTimestamps() bool {
	return getEnvWithFallback(c, envvar.LogUnixTS)
}

func (c *generalConfig) MercuryCredentials(url string) (username, password string, err error) {
	return "", "", errors.New("legacy config does not support Mercury credentials; use V2 TOML config to enable this feature")
}

// Port represents the port Chainlink should listen on for client requests.
func (c *generalConfig) Port() uint16 {
	return getEnvWithFallback(c, envvar.NewUint16("Port"))
}

// DefaultChainID represents the chain ID which jobs will use if one is not explicitly specified
func (c *generalConfig) DefaultChainID() *big.Int {
	str := c.viper.GetString(envvar.Name("DefaultChainID"))
	if str != "" {
		v, err := parse.BigInt(str)
		if err != nil {
			c.lggr.Errorw(
				"Ignoring invalid value provided for ETH_CHAIN_ID",
				"value", str,
				"error", err)
			return nil
		}
		return v

	}
	return nil
}

// HTTPServerWriteTimeout controls how long chainlink's API server may hold a
// socket open for writing a response to an HTTP request. This sometimes needs
// to be increased for pprof.
func (c *generalConfig) HTTPServerWriteTimeout() time.Duration {
	return getEnvWithFallback(c, envvar.HTTPServerWriteTimeout)
}

// ReaperExpiration represents how long a session is held in the DB before being cleared
func (c *generalConfig) ReaperExpiration() models.Duration {
	return models.MustMakeDuration(getEnvWithFallback(c, envvar.NewDuration("ReaperExpiration")))
}

// RootDir represents the location on the file system where Chainlink should
// keep its files.
func (c *generalConfig) RootDir() string {
	return getEnvWithFallback(c, envvar.RootDir)
}

// RPID Fetches the RPID used for WebAuthn sessions. The RPID value should be the FQDN (localhost)
func (c *generalConfig) RPID() string {
	return c.viper.GetString(envvar.Name("RPID"))
}

// RPOrigin Fetches the RPOrigin used to configure WebAuthn sessions. The RPOrigin value should be
// the origin URL where WebAuthn requests initiate (http://localhost:6688/)
func (c *generalConfig) RPOrigin() string {
	return c.viper.GetString(envvar.Name("RPOrigin"))
}

// SecureCookies allows toggling of the secure cookies HTTP flag
func (c *generalConfig) SecureCookies() bool {
	return c.viper.GetBool(envvar.Name("SecureCookies"))
}

// SessionTimeout is the maximum duration that a user session can persist without any activity.
func (c *generalConfig) SessionTimeout() models.Duration {
	return models.MustMakeDuration(getEnvWithFallback(c, envvar.NewDuration("SessionTimeout")))
}

func (c *generalConfig) SentryDSN() string {
	return os.Getenv("SENTRY_DSN")
}

func (c *generalConfig) SentryDebug() bool {
	return os.Getenv("SENTRY_DEBUG") == "true"
}

func (c *generalConfig) SentryEnvironment() string {
	return os.Getenv("SENTRY_ENVIRONMENT")
}

func (c *generalConfig) SentryRelease() string {
	return os.Getenv("SENTRY_RELEASE")
}

// TLSCertPath represents the file system location of the TLS certificate
// Chainlink should use for HTTPS.
func (c *generalConfig) TLSCertPath() string {
	return c.viper.GetString(envvar.Name("TLSCertPath"))
}

// TLSHost represents the hostname to use for TLS clients. This should match
// the TLS certificate.
func (c *generalConfig) TLSHost() string {
	return c.viper.GetString(envvar.Name("TLSHost"))
}

// TLSKeyPath represents the file system location of the TLS key Chainlink
// should use for HTTPS.
func (c *generalConfig) TLSKeyPath() string {
	return c.viper.GetString(envvar.Name("TLSKeyPath"))
}

// TLSPort represents the port Chainlink should listen on for encrypted client requests.
func (c *generalConfig) TLSPort() uint16 {
	return getEnvWithFallback(c, envvar.NewUint16("TLSPort"))
}

// TLSRedirect forces TLS redirect for unencrypted connections
func (c *generalConfig) TLSRedirect() bool {
	return c.viper.GetBool(envvar.Name("TLSRedirect"))
}

// UnAuthenticatedRateLimit defines the threshold to which requests unauthenticated requests get limited
func (c *generalConfig) UnAuthenticatedRateLimit() int64 {
	return c.viper.GetInt64(envvar.Name("UnAuthenticatedRateLimit"))
}

// UnAuthenticatedRateLimitPeriod defines the period to which unauthenticated requests get limited
func (c *generalConfig) UnAuthenticatedRateLimitPeriod() models.Duration {
	return models.MustMakeDuration(getEnvWithFallback(c, envvar.NewDuration("UnAuthenticatedRateLimitPeriod")))
}

func (c *generalConfig) TLSDir() string {
	return filepath.Join(c.RootDir(), "tls")
}

// KeyFile returns the path where the server key is kept
func (c *generalConfig) KeyFile() string {
	if c.TLSKeyPath() == "" {
		return filepath.Join(c.TLSDir(), "server.key")
	}
	return c.TLSKeyPath()
}

// CertFile returns the path where the server certificate is kept
func (c *generalConfig) CertFile() string {
	if c.TLSCertPath() == "" {
		return filepath.Join(c.TLSDir(), "server.crt")
	}
	return c.TLSCertPath()
}

// SessionOptions returns the sessions.Options struct used to configure
// the session store.
func (c *generalConfig) SessionOptions() sessions.Options {
	return sessions.Options{
		Secure:   c.SecureCookies(),
		HttpOnly: true,
		MaxAge:   86400 * 30,
		SameSite: http.SameSiteStrictMode,
	}
}

// Deprecated - prefer getEnvWithFallback with an EnvVar
func (c *generalConfig) getWithFallback(name string, parser func(string) (interface{}, error)) interface{} {
	return getEnvWithFallback(c, envvar.New(name, parser))
}

func getEnvWithFallback[T any](c *generalConfig, e *envvar.EnvVar[T]) T {
	v, invalid, err := e.ParseFrom(c.viper.GetString)
	if err != nil {
		c.lggr.Panic(err)
	}
	if invalid != "" {
		c.lggr.Error(invalid)
	}
	return v
}

type DatabaseBackupMode string

var (
	DatabaseBackupModeNone DatabaseBackupMode = "none"
	DatabaseBackupModeLite DatabaseBackupMode = "lite"
	DatabaseBackupModeFull DatabaseBackupMode = "full"
)

func parseDatabaseBackupMode(s string) (DatabaseBackupMode, error) {
	switch DatabaseBackupMode(s) {
	case DatabaseBackupModeNone, DatabaseBackupModeLite, DatabaseBackupModeFull:
		return DatabaseBackupMode(s), nil
	default:
		return "", fmt.Errorf("unable to parse %v into DatabaseBackupMode. Must be one of values: \"%s\", \"%s\", \"%s\"", s, DatabaseBackupModeNone, DatabaseBackupModeLite, DatabaseBackupModeFull)
	}
}

func lookupEnv[T any](c *generalConfig, k string, parse func(string) (T, error)) (t T, ok bool) {
	s, ok := os.LookupEnv(k)
	if !ok {
		return
	}
	val, err := parse(s)
	if err == nil {
		return val, true
	}
	c.lggr.Errorw(fmt.Sprintf("Invalid value provided for %s, falling back to default.", s),
		"value", s, "key", k, "error", err)
	return
}

// EVM methods

func (c *generalConfig) GlobalBalanceMonitorEnabled() (bool, bool) {
	return lookupEnv(c, envvar.Name("BalanceMonitorEnabled"), strconv.ParseBool)
}
func (c *generalConfig) GlobalBlockEmissionIdleWarningThreshold() (time.Duration, bool) {
	return lookupEnv(c, envvar.Name("BlockEmissionIdleWarningThreshold"), time.ParseDuration)
}
func (c *generalConfig) GlobalBlockHistoryEstimatorBatchSize() (uint32, bool) {
	return lookupEnv(c, envvar.Name("BlockHistoryEstimatorBatchSize"), parse.Uint32)
}
func (c *generalConfig) GlobalBlockHistoryEstimatorBlockDelay() (uint16, bool) {
	return lookupEnv(c, envvar.Name("BlockHistoryEstimatorBlockDelay"), parse.Uint16)
}
func (c *generalConfig) GlobalBlockHistoryEstimatorBlockHistorySize() (uint16, bool) {
	return lookupEnv(c, envvar.Name("BlockHistoryEstimatorBlockHistorySize"), parse.Uint16)
}
func (c *generalConfig) GlobalBlockHistoryEstimatorTransactionPercentile() (uint16, bool) {
	return lookupEnv(c, envvar.Name("BlockHistoryEstimatorTransactionPercentile"), parse.Uint16)
}
func (c *generalConfig) GlobalEthTxReaperInterval() (time.Duration, bool) {
	return lookupEnv(c, envvar.Name("EthTxReaperInterval"), time.ParseDuration)
}
func (c *generalConfig) GlobalEthTxReaperThreshold() (time.Duration, bool) {
	return lookupEnv(c, envvar.Name("EthTxReaperThreshold"), time.ParseDuration)
}
func (c *generalConfig) GlobalEthTxResendAfterThreshold() (time.Duration, bool) {
	return lookupEnv(c, envvar.Name("EthTxResendAfterThreshold"), time.ParseDuration)
}
func (c *generalConfig) GlobalEvmFinalityDepth() (uint32, bool) {
	return lookupEnv(c, envvar.Name("EvmFinalityDepth"), parse.Uint32)
}
func (c *generalConfig) GlobalEvmGasBumpPercent() (uint16, bool) {
	return lookupEnv(c, envvar.Name("EvmGasBumpPercent"), parse.Uint16)
}
func (c *generalConfig) GlobalEvmGasBumpThreshold() (uint64, bool) {
	return lookupEnv(c, envvar.Name("EvmGasBumpThreshold"), parse.Uint64)
}
func (c *generalConfig) GlobalEvmGasBumpTxDepth() (uint16, bool) {
	return lookupEnv(c, envvar.Name("EvmGasBumpTxDepth"), parse.Uint16)
}
func (c *generalConfig) GlobalEvmGasBumpWei() (*assets.Wei, bool) {
	return lookupEnv(c, envvar.Name("EvmGasBumpWei"), parse.Wei)
}
func (c *generalConfig) GlobalEvmGasFeeCapDefault() (*assets.Wei, bool) {
	return lookupEnv(c, envvar.Name("EvmGasFeeCapDefault"), parse.Wei)
}
func (c *generalConfig) GlobalBlockHistoryEstimatorEIP1559FeeCapBufferBlocks() (uint16, bool) {
	return lookupEnv(c, envvar.Name("BlockHistoryEstimatorEIP1559FeeCapBufferBlocks"), parse.Uint16)
}
func (c *generalConfig) GlobalBlockHistoryEstimatorCheckInclusionBlocks() (uint16, bool) {
	return lookupEnv(c, envvar.Name("BlockHistoryEstimatorCheckInclusionBlocks"), parse.Uint16)
}
func (c *generalConfig) GlobalBlockHistoryEstimatorCheckInclusionPercentile() (uint16, bool) {
	return lookupEnv(c, envvar.Name("BlockHistoryEstimatorCheckInclusionPercentile"), parse.Uint16)
}
func (c *generalConfig) GlobalEvmGasLimitDefault() (uint32, bool) {
	return lookupEnv(c, envvar.Name("EvmGasLimitDefault"), parse.Uint32)
}
func (c *generalConfig) GlobalEvmGasLimitMax() (uint32, bool) {
	return lookupEnv(c, envvar.Name("EvmGasLimitMax"), parse.Uint32)
}
func (c *generalConfig) GlobalEvmGasLimitMultiplier() (float32, bool) {
	return lookupEnv(c, envvar.Name("EvmGasLimitMultiplier"), parse.F32)
}
func (c *generalConfig) GlobalEvmGasLimitTransfer() (uint32, bool) {
	return lookupEnv(c, envvar.Name("EvmGasLimitTransfer"), parse.Uint32)
}
func (c *generalConfig) GlobalEvmGasPriceDefault() (*assets.Wei, bool) {
	return lookupEnv(c, envvar.Name("EvmGasPriceDefault"), parse.Wei)
}
func (c *generalConfig) GlobalEvmGasLimitOCRJobType() (uint32, bool) {
	return lookupEnv(c, envvar.Name("EvmGasLimitOCRJobType"), parse.Uint32)
}
func (c *generalConfig) GlobalEvmGasLimitDRJobType() (uint32, bool) {
	return lookupEnv(c, envvar.Name("EvmGasLimitDRJobType"), parse.Uint32)
}
func (c *generalConfig) GlobalEvmGasLimitVRFJobType() (uint32, bool) {
	return lookupEnv(c, envvar.Name("EvmGasLimitVRFJobType"), parse.Uint32)
}
func (c *generalConfig) GlobalEvmGasLimitFMJobType() (uint32, bool) {
	return lookupEnv(c, envvar.Name("EvmGasLimitFMJobType"), parse.Uint32)
}
func (c *generalConfig) GlobalEvmGasLimitKeeperJobType() (uint32, bool) {
	return lookupEnv(c, envvar.Name("EvmGasLimitKeeperJobType"), parse.Uint32)
}
func (c *generalConfig) GlobalEvmHeadTrackerHistoryDepth() (uint32, bool) {
	return lookupEnv(c, envvar.Name("EvmHeadTrackerHistoryDepth"), parse.Uint32)
}
func (c *generalConfig) GlobalEvmHeadTrackerMaxBufferSize() (uint32, bool) {
	return lookupEnv(c, envvar.Name("EvmHeadTrackerMaxBufferSize"), parse.Uint32)
}
func (c *generalConfig) GlobalEvmHeadTrackerSamplingInterval() (time.Duration, bool) {
	return lookupEnv(c, envvar.Name("EvmHeadTrackerSamplingInterval"), time.ParseDuration)
}
func (c *generalConfig) GlobalEvmLogBackfillBatchSize() (uint32, bool) {
	return lookupEnv(c, envvar.Name("EvmLogBackfillBatchSize"), parse.Uint32)
}
func (c *generalConfig) GlobalEvmLogPollInterval() (time.Duration, bool) {
	return lookupEnv(c, envvar.Name("EvmLogPollInterval"), time.ParseDuration)
}
func (c *generalConfig) GlobalEvmLogKeepBlocksDepth() (uint32, bool) {
	return lookupEnv(c, envvar.Name("EvmLogKeepBlocksDepth"), parse.Uint32)
}
func (c *generalConfig) GlobalEvmMaxGasPriceWei() (*assets.Wei, bool) {
	return lookupEnv(c, envvar.Name("EvmMaxGasPriceWei"), parse.Wei)
}
func (c *generalConfig) GlobalEvmMaxInFlightTransactions() (uint32, bool) {
	return lookupEnv(c, envvar.Name("EvmMaxInFlightTransactions"), parse.Uint32)
}
func (c *generalConfig) GlobalEvmMaxQueuedTransactions() (uint64, bool) {
	return lookupEnv(c, envvar.Name("EvmMaxQueuedTransactions"), parse.Uint64)
}
func (c *generalConfig) GlobalEvmMinGasPriceWei() (*assets.Wei, bool) {
	return lookupEnv(c, envvar.Name("EvmMinGasPriceWei"), parse.Wei)
}
func (c *generalConfig) GlobalEvmNonceAutoSync() (bool, bool) {
	return lookupEnv(c, envvar.Name("EvmNonceAutoSync"), strconv.ParseBool)
}
func (c *generalConfig) GlobalEvmUseForwarders() (bool, bool) {
	return lookupEnv(c, envvar.Name("EvmUseForwarders"), strconv.ParseBool)
}
func (c *generalConfig) GlobalEvmRPCDefaultBatchSize() (uint32, bool) {
	return lookupEnv(c, envvar.Name("EvmRPCDefaultBatchSize"), parse.Uint32)
}
func (c *generalConfig) GlobalFlagsContractAddress() (string, bool) {
	return lookupEnv(c, envvar.Name("FlagsContractAddress"), parse.String)
}
func (c *generalConfig) GlobalGasEstimatorMode() (string, bool) {
	return lookupEnv(c, envvar.Name("GasEstimatorMode"), parse.String)
}

// GlobalChainType overrides all chains and forces them to act as a particular
// chain type. List of chain types is given in `chaintype.go`.
func (c *generalConfig) GlobalChainType() (string, bool) {
	return lookupEnv(c, envvar.Name("ChainType"), parse.String)
}
func (c *generalConfig) GlobalLinkContractAddress() (string, bool) {
	return lookupEnv(c, envvar.Name("LinkContractAddress"), parse.String)
}
func (c *generalConfig) GlobalOperatorFactoryAddress() (string, bool) {
	return lookupEnv(c, envvar.Name("OperatorFactoryAddress"), parse.String)
}
func (c *generalConfig) GlobalMinIncomingConfirmations() (uint32, bool) {
	return lookupEnv(c, envvar.Name("MinIncomingConfirmations"), parse.Uint32)
}
func (c *generalConfig) GlobalMinimumContractPayment() (*assets.Link, bool) {
	return lookupEnv(c, envvar.Name("MinimumContractPayment"), parse.Link)
}
func (c *generalConfig) GlobalEvmEIP1559DynamicFees() (bool, bool) {
	return lookupEnv(c, envvar.Name("EvmEIP1559DynamicFees"), strconv.ParseBool)
}
func (c *generalConfig) GlobalEvmGasTipCapDefault() (*assets.Wei, bool) {
	return lookupEnv(c, envvar.Name("EvmGasTipCapDefault"), parse.Wei)
}
func (c *generalConfig) GlobalEvmGasTipCapMinimum() (*assets.Wei, bool) {
	return lookupEnv(c, envvar.Name("EvmGasTipCapMinimum"), parse.Wei)
}

func (c *generalConfig) GlobalNodeNoNewHeadsThreshold() (time.Duration, bool) {
	return lookupEnv(c, envvar.Name("NodeNoNewHeadsThreshold"), time.ParseDuration)
}

func (c *generalConfig) GlobalNodePollFailureThreshold() (uint32, bool) {
	return lookupEnv(c, envvar.Name("NodePollFailureThreshold"), parse.Uint32)
}

func (c *generalConfig) GlobalNodePollInterval() (time.Duration, bool) {
	return lookupEnv(c, envvar.Name("NodePollInterval"), time.ParseDuration)
}

func (c *generalConfig) GlobalNodeSelectionMode() (string, bool) {
	return lookupEnv(c, envvar.Name("NodeSelectionMode"), parse.String)
}

func (c *generalConfig) GlobalNodeSyncThreshold() (uint32, bool) {
	return lookupEnv(c, envvar.Name("NodeSyncThreshold"), parse.Uint32)
}

func (c *generalConfig) GlobalOCR2AutomationGasLimit() (uint32, bool) {
	return lookupEnv(c, envvar.Name("OCR2AutomationGasLimit"), parse.Uint32)
}

// DatabaseLockingMode can be one of 'dual', 'advisorylock', 'lease' or 'none'
// It controls which mode to use to enforce that only one Chainlink application can use the database
func (c *generalConfig) DatabaseLockingMode() string {
	return getEnvWithFallback(c, envvar.NewString("DatabaseLockingMode"))
}

// LeaseLockRefreshInterval controls how often the node should attempt to
// refresh the lease lock
func (c *generalConfig) LeaseLockRefreshInterval() time.Duration {
	return c.getDuration("LeaseLockRefreshInterval")
}

// LeaseLockDuration controls when the lock is set to expire on each refresh
// (this many seconds from now in the future)
func (c *generalConfig) LeaseLockDuration() time.Duration {
	return c.getDuration("LeaseLockDuration")
}

// AdvisoryLockID is the application advisory lock ID. Should match all other
// chainlink applications that might access this database
func (c *generalConfig) AdvisoryLockID() int64 {
	return getEnvWithFallback(c, envvar.AdvisoryLockID)
}

// AdvisoryLockCheckInterval controls how often Chainlink will check to make
// sure it still holds the advisory lock. If it no longer holds it, it will try
// to re-acquire it and if that fails the application will exit
func (c *generalConfig) AdvisoryLockCheckInterval() time.Duration {
	return c.getDuration("AdvisoryLockCheckInterval")
}

// LogFileDir if set will override RootDir as the output path for log files
func (c *generalConfig) LogFileDir() string {
	s := c.viper.GetString(envvar.Name("LogFileDir"))
	if s == "" {
		return c.RootDir()
	}
	return s
}

func (c *generalConfig) SetPasswords(keystore, vrf *string) {
	c.passwordMu.Lock()
	defer c.passwordMu.Unlock()
	if keystore != nil {
		c.passwordKeystore = *keystore
	}
	if vrf != nil {
		c.passwordVRF = *vrf
	}
}

func (c *generalConfig) KeystorePassword() string {
	c.passwordMu.RLock()
	defer c.passwordMu.RUnlock()
	return c.passwordKeystore
}

func (c *generalConfig) VRFPassword() string {
	c.passwordMu.RLock()
	defer c.passwordMu.RUnlock()
	return c.passwordVRF
}
