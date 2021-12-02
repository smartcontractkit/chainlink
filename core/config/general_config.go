package config

import (
	"fmt"
	"log"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sync"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name GeneralConfig --output ./mocks/ --case=underscore

// this permission grants read / write access to file owners only
const readWritePerms = os.FileMode(0600)

var (
	ErrUnset        = errors.New("env var unset")
	ErrInvalid      = errors.New("env var invalid")
	DefaultLogLevel = LogLevel{zapcore.InfoLevel}

	configFileNotFoundError = reflect.TypeOf(viper.ConfigFileNotFoundError{})
)

type GeneralOnlyConfig interface {
	AdminCredentialsFile() string
	AllowOrigins() string
	AuthenticatedRateLimit() int64
	AuthenticatedRateLimitPeriod() models.Duration
	AutoPprofEnabled() bool
	AutoPprofProfileRoot() string
	AutoPprofPollInterval() models.Duration
	AutoPprofGatherDuration() models.Duration
	AutoPprofGatherTraceDuration() models.Duration
	AutoPprofMaxProfileSize() utils.FileSize
	AutoPprofCPUProfileRate() int
	AutoPprofMemProfileRate() int
	AutoPprofBlockProfileRate() int
	AutoPprofMutexProfileFraction() int
	AutoPprofMemThreshold() utils.FileSize
	AutoPprofGoroutineThreshold() int
	BlockBackfillDepth() uint64
	BlockBackfillSkip() bool
	BridgeResponseURL() *url.URL
	CertFile() string
	ClientNodeURL() string
	DatabaseBackupDir() string
	DatabaseBackupFrequency() time.Duration
	DatabaseBackupMode() DatabaseBackupMode
	DatabaseBackupURL() *url.URL
	DatabaseListenerMaxReconnectDuration() time.Duration
	DatabaseListenerMinReconnectInterval() time.Duration
	DatabaseLockingMode() string
	DatabaseURL() url.URL
	DefaultChainID() *big.Int
	DefaultHTTPAllowUnrestrictedNetworkAccess() bool
	DefaultHTTPLimit() int64
	DefaultHTTPTimeout() models.Duration
	DefaultMaxHTTPAttempts() uint
	Dev() bool
	EVMDisabled() bool
	EthereumDisabled() bool
	EthereumHTTPURL() *url.URL
	EthereumSecondaryURLs() []url.URL
	EthereumURL() string
	ExplorerAccessKey() string
	ExplorerSecret() string
	ExplorerURL() *url.URL
	FMDefaultTransactionQueueDepth() uint32
	FMSimulateTransactions() bool
	FeatureExternalInitiators() bool
	FeatureOffchainReporting() bool
	FeatureOffchainReporting2() bool
	FeatureUICSAKeys() bool
	FeatureUIFeedsManager() bool
	GetAdvisoryLockIDConfiguredOrDefault() int64
	GetDatabaseDialectConfiguredOrDefault() dialects.DialectName
	GlobalLockRetryInterval() models.Duration
	HTTPServerWriteTimeout() time.Duration
	InsecureFastScrypt() bool
	InsecureSkipVerify() bool
	JSONConsole() bool
	JobPipelineMaxRunDuration() time.Duration
	JobPipelineReaperInterval() time.Duration
	JobPipelineReaperThreshold() time.Duration
	JobPipelineResultWriteQueueDepth() uint64
	KeeperDefaultTransactionQueueDepth() uint32
	KeeperGasPriceBufferPercent() uint32
	KeeperGasTipCapBufferPercent() uint32
	KeeperMaximumGracePeriod() int64
	KeeperRegistryCheckGasOverhead() uint64
	KeeperRegistryPerformGasOverhead() uint64
	KeeperRegistrySyncInterval() time.Duration
	KeeperRegistrySyncUpkeepQueueSize() uint32
	KeyFile() string
	LeaseLockRefreshInterval() time.Duration
	LeaseLockDuration() time.Duration
	LogLevel() zapcore.Level
	DefaultLogLevel() zapcore.Level
	LogSQLMigrations() bool
	LogSQL() bool
	LogToDisk() bool
	LogUnixTimestamps() bool
	MigrateDatabase() bool
	ORMMaxIdleConns() int
	ORMMaxOpenConns() int
	Port() uint16
	RPID() string
	RPOrigin() string
	ReaperExpiration() models.Duration
	ReplayFromBlock() int64
	RootDir() string
	SecureCookies() bool
	SessionOptions() sessions.Options
	SessionSecret() ([]byte, error)
	SessionTimeout() models.Duration
	SetDialect(dialects.DialectName)
	SetLogLevel(lvl zapcore.Level) error
	SetLogSQL(logSQL bool)
	StatsPusherLogging() bool
	TLSCertPath() string
	TLSDir() string
	TLSHost() string
	TLSKeyPath() string
	TLSPort() uint16
	TLSRedirect() bool
	TelemetryIngressLogging() bool
	TelemetryIngressServerPubKey() string
	TelemetryIngressURL() *url.URL
	TriggerFallbackDBPollInterval() time.Duration
	UnAuthenticatedRateLimit() int64
	UnAuthenticatedRateLimitPeriod() models.Duration
	UseLegacyEthEnvVars() bool
	Validate() error

	OCRConfig
	P2PNetworking
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
	GlobalBlockHistoryEstimatorTransactionPercentile() (uint16, bool)
	GlobalEthTxReaperInterval() (time.Duration, bool)
	GlobalEthTxReaperThreshold() (time.Duration, bool)
	GlobalEthTxResendAfterThreshold() (time.Duration, bool)
	GlobalEvmDefaultBatchSize() (uint32, bool)
	GlobalEvmEIP1559DynamicFees() (bool, bool)
	GlobalEvmFinalityDepth() (uint32, bool)
	GlobalEvmGasBumpPercent() (uint16, bool)
	GlobalEvmGasBumpThreshold() (uint64, bool)
	GlobalEvmGasBumpTxDepth() (uint16, bool)
	GlobalEvmGasBumpWei() (*big.Int, bool)
	GlobalEvmGasLimitDefault() (uint64, bool)
	GlobalEvmGasLimitMultiplier() (float32, bool)
	GlobalEvmGasLimitTransfer() (uint64, bool)
	GlobalEvmGasPriceDefault() (*big.Int, bool)
	GlobalEvmGasTipCapDefault() (*big.Int, bool)
	GlobalEvmGasTipCapMinimum() (*big.Int, bool)
	GlobalEvmHeadTrackerHistoryDepth() (uint32, bool)
	GlobalEvmHeadTrackerMaxBufferSize() (uint32, bool)
	GlobalEvmHeadTrackerSamplingInterval() (time.Duration, bool)
	GlobalEvmLogBackfillBatchSize() (uint32, bool)
	GlobalEvmMaxGasPriceWei() (*big.Int, bool)
	GlobalEvmMaxInFlightTransactions() (uint32, bool)
	GlobalEvmMaxQueuedTransactions() (uint64, bool)
	GlobalEvmMinGasPriceWei() (*big.Int, bool)
	GlobalEvmNonceAutoSync() (bool, bool)
	GlobalEvmRPCDefaultBatchSize() (uint32, bool)
	GlobalFlagsContractAddress() (string, bool)
	GlobalGasEstimatorMode() (string, bool)
	GlobalChainType() (string, bool)
	GlobalLinkContractAddress() (string, bool)
	GlobalMinIncomingConfirmations() (uint32, bool)
	GlobalMinRequiredOutgoingConfirmations() (uint64, bool)
	GlobalMinimumContractPayment() (*assets.Link, bool)
	GlobalOCRContractConfirmations() (uint16, bool)
	GlobalOCRContractTransmitterTransmitTimeout() (time.Duration, bool)
	GlobalOCRDatabaseTimeout() (time.Duration, bool)
	GlobalOCRObservationGracePeriod() (time.Duration, bool)
	GlobalOCR2ContractConfirmations() (uint16, bool)
}

type GeneralConfig interface {
	GeneralOnlyConfig
	GlobalConfig
}

// generalConfig holds parameters used by the application which can be overridden by
// setting environment variables.
//
// If you add an entry here which does not contain sensitive information, you
// should also update presenters.ConfigWhitelist and cmd_test.TestClient_RunNodeShowsEnv.
type generalConfig struct {
	viper            *viper.Viper
	secretGenerator  SecretGenerator
	randomP2PPort    uint16
	randomP2PPortMtx *sync.RWMutex
	dialect          dialects.DialectName
	advisoryLockID   int64
	logLevel         zapcore.Level
	defaultLogLevel  zapcore.Level
	logSQL           bool
	logMutex         sync.RWMutex
}

// NewGeneralConfig returns the config with the environment variables set to their
// respective fields, or their defaults if environment variables are not set.
func NewGeneralConfig() GeneralConfig {
	v := viper.New()
	c := newGeneralConfigWithViper(v)
	c.secretGenerator = FilePersistedSecretGenerator{}
	c.dialect = dialects.Postgres
	return c
}

func newGeneralConfigWithViper(v *viper.Viper) *generalConfig {
	schemaT := reflect.TypeOf(ConfigSchema{})
	for index := 0; index < schemaT.NumField(); index++ {
		item := schemaT.FieldByIndex([]int{index})
		name := item.Tag.Get("env")
		def, exists := item.Tag.Lookup("default")
		if exists {
			v.SetDefault(name, def)
		}
		_ = v.BindEnv(name, name)
	}

	config := &generalConfig{
		viper:            v,
		randomP2PPortMtx: new(sync.RWMutex),
		defaultLogLevel:  DefaultLogLevel.Level,
	}

	if err := utils.EnsureDirAndMaxPerms(config.RootDir(), os.FileMode(0700)); err != nil {
		logger.Fatalf(`Error creating root directory "%s": %+v`, config.RootDir(), err)
	}

	v.SetConfigName("chainlink")
	v.AddConfigPath(config.RootDir())
	err := v.ReadInConfig()
	if err != nil && reflect.TypeOf(err) != configFileNotFoundError {
		logger.Warnf("Unable to load config file: %v\n", err)
	}

	if v.IsSet(EnvVarName("LogLevel")) {
		str := v.GetString(EnvVarName("LogLevel"))
		ll, err := ParseLogLevel(str)
		if err != nil {
			logger.Errorf("error parsing log level: %s, falling back to %s", str, DefaultLogLevel.Level)
		} else {
			config.defaultLogLevel = ll.(LogLevel).Level
		}
	}
	config.logLevel = config.defaultLogLevel
	config.logSQL = viper.GetBool(EnvVarName("LogSQL"))
	config.logMutex = sync.RWMutex{}

	return config
}

// Validate performs basic sanity checks on config and returns error if any
// misconfiguration would be fatal to the application
func (c *generalConfig) Validate() error {
	if c.P2PAnnouncePort() != 0 && c.P2PAnnounceIP() == nil {
		return errors.Errorf("P2P_ANNOUNCE_PORT was given as %v but P2P_ANNOUNCE_IP was unset. You must also set P2P_ANNOUNCE_IP if P2P_ANNOUNCE_PORT is set", c.P2PAnnouncePort())
	}

	if _, exists := os.LookupEnv("MINIMUM_CONTRACT_PAYMENT"); exists {
		return errors.Errorf("MINIMUM_CONTRACT_PAYMENT is deprecated, use MINIMUM_CONTRACT_PAYMENT_LINK_JUELS instead.")
	}

	if _, err := c.OCRKeyBundleID(); errors.Cause(err) == ErrInvalid {
		return err
	}
	if _, err := c.OCRTransmitterAddress(); errors.Cause(err) == ErrInvalid {
		return err
	}
	if peers, err := c.P2PBootstrapPeers(); err == nil {
		for i := range peers {
			if _, err := multiaddr.NewMultiaddr(peers[i]); err != nil {
				return errors.Errorf("p2p bootstrap peer %d is invalid: err %v", i, err)
			}
		}
	}
	if me := c.OCRMonitoringEndpoint(); me != "" {
		if _, err := url.Parse(me); err != nil {
			return errors.Wrapf(err, "invalid monitoring url: %s", me)
		}
	}
	if ct, set := c.GlobalChainType(); set && !chains.ChainType(ct).IsValid() {
		return errors.Errorf("CHAIN_TYPE is invalid: %s", ct)
	}

	if !c.UseLegacyEthEnvVars() {
		if c.EthereumURL() != "" {
			logger.Warn("ETH_URL has no effect when USE_LEGACY_ETH_ENV_VARS=false")
		}
		if c.EthereumHTTPURL() != nil {
			logger.Warn("ETH_HTTP_URL has no effect when USE_LEGACY_ETH_ENV_VARS=false")
		}
		if len(c.EthereumSecondaryURLs()) > 0 {
			logger.Warn("ETH_SECONDARY_URL/ETH_SECONDARY_URLS have no effect when USE_LEGACY_ETH_ENV_VARS=false")
		}
	}
	// Warn on legacy OCR env vars
	if c.OCRDHTLookupInterval() != 0 {
		logger.Warn("OCR_DHT_LOOKUP_INTERVAL is deprecated, use P2P_DHT_LOOKUP_INTERVAL instead")
	}
	if c.OCRBootstrapCheckInterval() != 0 {
		logger.Warn("OCR_BOOTSTRAP_CHECK_INTERVAL is deprecated, use P2P_BOOTSTRAP_CHECK_INTERVAL instead")
	}
	if c.OCRIncomingMessageBufferSize() != 0 {
		logger.Warn("OCR_INCOMING_MESSAGE_BUFFER_SIZE is deprecated, use P2P_INCOMING_MESSAGE_BUFFER_SIZE instead")
	}
	if c.OCROutgoingMessageBufferSize() != 0 {
		logger.Warn("OCR_OUTGOING_MESSAGE_BUFFER_SIZE is deprecated, use P2P_OUTGOING_MESSAGE_BUFFER_SIZE instead")
	}
	if c.OCRNewStreamTimeout() != 0 {
		logger.Warn("OCR_NEW_STREAM_TIMEOUT is deprecated, use P2P_NEW_STREAM_TIMEOUT instead")
	}

	switch c.DatabaseLockingMode() {
	case "dual", "lease", "advisorylock", "none":
	default:
		return errors.Errorf("unrecognised value for DATABASE_LOCKING_MODE: %s (valid options are 'dual', 'lease', 'advisorylock' or 'none')", c.DatabaseLockingMode())
	}

	if c.LeaseLockRefreshInterval() > c.LeaseLockDuration()/2 {
		return errors.Errorf("LEASE_LOCK_REFRESH_INTERVAL must be less than or equal to half of LEASE_LOCK_DURATION (got LEASE_LOCK_REFRESH_INTERVAL=%s, LEASE_LOCK_DURATION=%s)", c.LeaseLockRefreshInterval().String(), c.LeaseLockDuration().String())
	}

	return nil
}

func (c *generalConfig) SetDialect(d dialects.DialectName) {
	c.dialect = d
}

func (c *generalConfig) GetAdvisoryLockIDConfiguredOrDefault() int64 {
	return c.advisoryLockID
}

func (c *generalConfig) GetDatabaseDialectConfiguredOrDefault() dialects.DialectName {
	return c.dialect
}

// AllowOrigins returns the CORS hosts used by the frontend.
func (c *generalConfig) AllowOrigins() string {
	return c.viper.GetString(EnvVarName("AllowOrigins"))
}

// AdminCredentialsFile points to text file containing admin credentials for logging in
func (c *generalConfig) AdminCredentialsFile() string {
	fieldName := "AdminCredentialsFile"
	file := c.viper.GetString(EnvVarName(fieldName))
	defaultValue, _ := defaultValue(fieldName)
	if file == defaultValue {
		return filepath.Join(c.RootDir(), "apicredentials")
	}
	return file
}

// AuthenticatedRateLimit defines the threshold to which requests authenticated requests get limited
func (c *generalConfig) AuthenticatedRateLimit() int64 {
	return c.viper.GetInt64(EnvVarName("AuthenticatedRateLimit"))
}

// AuthenticatedRateLimitPeriod defines the period to which authenticated requests get limited
func (c *generalConfig) AuthenticatedRateLimitPeriod() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("AuthenticatedRateLimitPeriod", ParseDuration).(time.Duration))
}

func (c *generalConfig) AutoPprofEnabled() bool {
	return c.viper.GetBool(EnvVarName("AutoPprofEnabled"))
}

func (c *generalConfig) AutoPprofProfileRoot() string {
	root := c.viper.GetString(EnvVarName("AutoPprofProfileRoot"))
	if root == "" {
		return c.RootDir()
	}
	return root
}

func (c *generalConfig) AutoPprofPollInterval() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("AutoPprofPollInterval", ParseDuration).(time.Duration))
}

func (c *generalConfig) AutoPprofGatherDuration() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("AutoPprofGatherDuration", ParseDuration).(time.Duration))
}

func (c *generalConfig) AutoPprofGatherTraceDuration() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("AutoPprofGatherTraceDuration", ParseDuration).(time.Duration))
}

func (c *generalConfig) AutoPprofMaxProfileSize() utils.FileSize {
	return c.getWithFallback("AutoPprofMaxProfileSize", ParseFileSize).(utils.FileSize)
}

func (c *generalConfig) AutoPprofCPUProfileRate() int {
	return c.viper.GetInt(EnvVarName("AutoPprofCPUProfileRate"))
}

func (c *generalConfig) AutoPprofMemProfileRate() int {
	return c.viper.GetInt(EnvVarName("AutoPprofMemProfileRate"))
}

func (c *generalConfig) AutoPprofBlockProfileRate() int {
	return c.viper.GetInt(EnvVarName("AutoPprofBlockProfileRate"))
}

func (c *generalConfig) AutoPprofMutexProfileFraction() int {
	return c.viper.GetInt(EnvVarName("AutoPprofMutexProfileFraction"))
}

func (c *generalConfig) AutoPprofMemThreshold() utils.FileSize {
	return c.getWithFallback("AutoPprofMemThreshold", ParseFileSize).(utils.FileSize)
}

func (c *generalConfig) AutoPprofGoroutineThreshold() int {
	return c.viper.GetInt(EnvVarName("AutoPprofGoroutineThreshold"))
}

// BlockBackfillDepth specifies the number of blocks before the current HEAD that the
// log broadcaster will try to re-consume logs from
func (c *generalConfig) BlockBackfillDepth() uint64 {
	return c.getWithFallback("BlockBackfillDepth", ParseUint64).(uint64)
}

// BlockBackfillSkip enables skipping of very long log backfills
func (c *generalConfig) BlockBackfillSkip() bool {
	return c.getWithFallback("BlockBackfillSkip", ParseBool).(bool)
}

// BridgeResponseURL represents the URL for bridges to send a response to.
func (c *generalConfig) BridgeResponseURL() *url.URL {
	return c.getWithFallback("BridgeResponseURL", ParseURL).(*url.URL)
}

// ClientNodeURL is the URL of the Ethereum node this Chainlink node should connect to.
func (c *generalConfig) ClientNodeURL() string {
	return c.viper.GetString(EnvVarName("ClientNodeURL"))
}

// FeatureUICSAKeys enables the CSA Keys UI Feature.
func (c *generalConfig) FeatureUICSAKeys() bool {
	return c.getWithFallback("FeatureUICSAKeys", ParseBool).(bool)
}

func (c *generalConfig) FeatureUIFeedsManager() bool {
	return c.getWithFallback("FeatureUIFeedsManager", ParseBool).(bool)
}

func (c *generalConfig) DatabaseListenerMinReconnectInterval() time.Duration {
	return c.getWithFallback("DatabaseListenerMinReconnectInterval", ParseDuration).(time.Duration)
}

func (c *generalConfig) DatabaseListenerMaxReconnectDuration() time.Duration {
	return c.getWithFallback("DatabaseListenerMaxReconnectDuration", ParseDuration).(time.Duration)
}

// DatabaseBackupMode sets the database backup mode
func (c *generalConfig) DatabaseBackupMode() DatabaseBackupMode {
	return c.getWithFallback("DatabaseBackupMode", parseDatabaseBackupMode).(DatabaseBackupMode)
}

// DatabaseBackupFrequency turns on the periodic database backup if set to a positive value
// DatabaseBackupMode must be then set to a value other than "none"
func (c *generalConfig) DatabaseBackupFrequency() time.Duration {
	return c.getWithFallback("DatabaseBackupFrequency", ParseDuration).(time.Duration)
}

// DatabaseBackupURL configures the URL for the database to backup, if it's to be different from the main on
func (c *generalConfig) DatabaseBackupURL() *url.URL {
	s := c.viper.GetString(EnvVarName("DatabaseBackupURL"))
	if s == "" {
		return nil
	}
	uri, err := url.Parse(s)
	if err != nil {
		logger.Errorf("invalid database backup url %s", s)
		return nil
	}
	return uri
}

// DatabaseBackupDir configures the directory for saving the backup file, if it's to be different from default one located in the RootDir
func (c *generalConfig) DatabaseBackupDir() string {
	return c.viper.GetString(EnvVarName("DatabaseBackupDir"))
}

// GlobalLockRetryInterval represents how long to wait before trying again to get the global advisory lock.
func (c *generalConfig) GlobalLockRetryInterval() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("GlobalLockRetryInterval", ParseDuration).(time.Duration))
}

// DatabaseURL configures the URL for chainlink to connect to. This must be
// a properly formatted URL, with a valid scheme (postgres://)
func (c *generalConfig) DatabaseURL() url.URL {
	s := c.viper.GetString(EnvVarName("DatabaseURL"))
	uri, err := url.Parse(s)
	if err != nil {
		logger.Error("invalid database url %s", s)
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
	return c.viper.GetBool(EnvVarName("MigrateDatabase"))
}

// DefaultMaxHTTPAttempts defines the limit for HTTP requests.
func (c *generalConfig) DefaultMaxHTTPAttempts() uint {
	return uint(c.getWithFallback("DefaultMaxHTTPAttempts", ParseUint64).(uint64))
}

// DefaultHTTPLimit defines the size limit for HTTP requests and responses
func (c *generalConfig) DefaultHTTPLimit() int64 {
	return c.viper.GetInt64(EnvVarName("DefaultHTTPLimit"))
}

// DefaultHTTPTimeout defines the default timeout for http requests
func (c *generalConfig) DefaultHTTPTimeout() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("DefaultHTTPTimeout", ParseDuration).(time.Duration))
}

// DefaultHTTPAllowUnrestrictedNetworkAccess controls whether http requests are unrestricted by default
// It is recommended that this be left disabled
func (c *generalConfig) DefaultHTTPAllowUnrestrictedNetworkAccess() bool {
	return c.viper.GetBool(EnvVarName("DefaultHTTPAllowUnrestrictedNetworkAccess"))
}

// Dev configures "development" mode for chainlink.
func (c *generalConfig) Dev() bool {
	return c.viper.GetBool(EnvVarName("Dev"))
}

// FeatureExternalInitiators enables the External Initiator feature.
func (c *generalConfig) FeatureExternalInitiators() bool {
	return c.viper.GetBool(EnvVarName("FeatureExternalInitiators"))
}

// FeatureOffchainReporting enables the OCR job type.
func (c *generalConfig) FeatureOffchainReporting() bool {
	return c.getWithFallback("FeatureOffchainReporting", ParseBool).(bool)
}

// FeatureOffchainReporting2 enables the OCR2 job type.
func (c *generalConfig) FeatureOffchainReporting2() bool {
	return c.getWithFallback("FeatureOffchainReporting2", ParseBool).(bool)
}

// FMDefaultTransactionQueueDepth controls the queue size for DropOldestStrategy in Flux Monitor
// Set to 0 to use SendEvery strategy instead
func (c *generalConfig) FMDefaultTransactionQueueDepth() uint32 {
	return c.viper.GetUint32(EnvVarName("FMDefaultTransactionQueueDepth"))
}

// FMSimulateTransactions enables using eth_call transaction simulation before
// sending when set to true
func (c *generalConfig) FMSimulateTransactions() bool {
	return c.viper.GetBool(EnvVarName("FMSimulateTransactions"))
}

// EthereumURL represents the URL of the Ethereum node to connect Chainlink to.
func (c *generalConfig) EthereumURL() string {
	return c.viper.GetString(EnvVarName("EthereumURL"))
}

// EthereumHTTPURL is an optional but recommended url that points to the HTTP port of the primary node
func (c *generalConfig) EthereumHTTPURL() (uri *url.URL) {
	urlStr := c.viper.GetString(EnvVarName("EthereumHTTPURL"))
	if urlStr == "" {
		return nil
	}
	var err error
	uri, err = url.Parse(urlStr)
	if err != nil || !(uri.Scheme == "http" || uri.Scheme == "https") {
		logger.Fatalf("Invalid Ethereum HTTP URL: %s, got error: %s", urlStr, err)
	}
	return
}

// EthereumSecondaryURLs is an optional backup RPC URL
// Must be http(s) format
// If specified, transactions will also be broadcast to this ethereum node
func (c *generalConfig) EthereumSecondaryURLs() []url.URL {
	oldConfig := c.viper.GetString(EnvVarName("EthereumSecondaryURL"))
	newConfig := c.viper.GetString(EnvVarName("EthereumSecondaryURLs"))

	config := ""
	if newConfig != "" {
		config = newConfig
	} else if oldConfig != "" {
		config = oldConfig
	}

	urlStrings := regexp.MustCompile(`\s*[;,]\s*`).Split(config, -1)
	urls := []url.URL{}
	for _, urlString := range urlStrings {
		if urlString == "" {
			continue
		}
		url, err := url.Parse(urlString)
		if err != nil {
			logger.Fatalf("Invalid Secondary Ethereum URL: %s, got error: %v", urlString, err)
		}
		urls = append(urls, *url)
	}

	return urls
}

// EthereumDisabled will substitute null Eth clients if set
func (c *generalConfig) EthereumDisabled() bool {
	return c.viper.GetBool(EnvVarName("EthereumDisabled"))
}

// EVMDisabled prevents any evm_chains from being loaded at all if set
func (c *generalConfig) EVMDisabled() bool {
	return c.viper.GetBool(EnvVarName("EVMDisabled"))
}

// InsecureFastScrypt causes all key stores to encrypt using "fast" scrypt params instead
// This is insecure and only useful for local testing. DO NOT SET THIS IN PRODUCTION
func (c *generalConfig) InsecureFastScrypt() bool {
	return c.viper.GetBool(EnvVarName("InsecureFastScrypt"))
}

// InsecureSkipVerify disables SSL certificate verification when connection to
// a chainlink client using the remote client, i.e. when executing most remote
// commands in the CLI.
//
// This is mostly useful for people who want to use TLS on localhost.
func (c *generalConfig) InsecureSkipVerify() bool {
	return c.viper.GetBool(EnvVarName("InsecureSkipVerify"))
}

func (c *generalConfig) TriggerFallbackDBPollInterval() time.Duration {
	return c.getWithFallback("TriggerFallbackDBPollInterval", ParseDuration).(time.Duration)
}

// JobPipelineMaxRunDuration is the maximum time that a job run may take
func (c *generalConfig) JobPipelineMaxRunDuration() time.Duration {
	return c.getWithFallback("JobPipelineMaxRunDuration", ParseDuration).(time.Duration)
}

func (c *generalConfig) JobPipelineResultWriteQueueDepth() uint64 {
	return c.getWithFallback("JobPipelineResultWriteQueueDepth", ParseUint64).(uint64)
}

func (c *generalConfig) JobPipelineReaperInterval() time.Duration {
	return c.getWithFallback("JobPipelineReaperInterval", ParseDuration).(time.Duration)
}

func (c *generalConfig) JobPipelineReaperThreshold() time.Duration {
	return c.getWithFallback("JobPipelineReaperThreshold", ParseDuration).(time.Duration)
}

// KeeperRegistryCheckGasOverhead is the amount of extra gas to provide checkUpkeep() calls
// to account for the gas consumed by the keeper registry
func (c *generalConfig) KeeperRegistryCheckGasOverhead() uint64 {
	return c.getWithFallback("KeeperRegistryCheckGasOverhead", ParseUint64).(uint64)
}

// KeeperRegistryPerformGasOverhead is the amount of extra gas to provide performUpkeep() calls
// to account for the gas consumed by the keeper registry
func (c *generalConfig) KeeperRegistryPerformGasOverhead() uint64 {
	return c.getWithFallback("KeeperRegistryPerformGasOverhead", ParseUint64).(uint64)
}

// KeeperDefaultTransactionQueueDepth controls the queue size for DropOldestStrategy in Keeper
// Set to 0 to use SendEvery strategy instead
func (c *generalConfig) KeeperDefaultTransactionQueueDepth() uint32 {
	return c.viper.GetUint32(EnvVarName("KeeperDefaultTransactionQueueDepth"))
}

// KeeperGasPriceBufferPercent adds the specified percentage to the gas price
// used for checking whether to perform an upkeep. Only applies in legacy mode.
func (c *generalConfig) KeeperGasPriceBufferPercent() uint32 {
	return c.viper.GetUint32(EnvVarName("KeeperGasPriceBufferPercent"))
}

// KeeperGasTipCapBufferPercent adds the specified percentage to the gas price
// used for checking whether to perform an upkeep. Only applies in EIP-1559 mode.
func (c *generalConfig) KeeperGasTipCapBufferPercent() uint32 {
	return c.viper.GetUint32(EnvVarName("KeeperGasTipCapBufferPercent"))
}

// KeeperRegistrySyncInterval is the interval in which the RegistrySynchronizer performs a full
// sync of the keeper registry contract it is tracking
func (c *generalConfig) KeeperRegistrySyncInterval() time.Duration {
	return c.getWithFallback("KeeperRegistrySyncInterval", ParseDuration).(time.Duration)
}

// KeeperMaximumGracePeriod is the maximum number of blocks that a keeper will wait after performing
// an upkeep before it resumes checking that upkeep
func (c *generalConfig) KeeperMaximumGracePeriod() int64 {
	return c.viper.GetInt64(EnvVarName("KeeperMaximumGracePeriod"))
}

// KeeperRegistrySyncUpkeepQueueSize represents the maximum number of upkeeps that can be synced in parallel
func (c *generalConfig) KeeperRegistrySyncUpkeepQueueSize() uint32 {
	return c.getWithFallback("KeeperRegistrySyncUpkeepQueueSize", ParseUint32).(uint32)
}

// JSONConsole when set to true causes logging to be made in JSON format
// If set to false, logs in console format
func (c *generalConfig) JSONConsole() bool {
	return c.viper.GetBool(EnvVarName("JSONConsole"))
}

// ExplorerURL returns the websocket URL for this node to push stats to, or nil.
func (c *generalConfig) ExplorerURL() *url.URL {
	rval := c.getWithFallback("ExplorerURL", ParseURL)
	switch t := rval.(type) {
	case nil:
		return nil
	case *url.URL:
		return t
	default:
		panic(fmt.Sprintf("invariant: ExplorerURL returned as type %T", rval))
	}
}

// ExplorerAccessKey returns the access key for authenticating with explorer
func (c *generalConfig) ExplorerAccessKey() string {
	return c.viper.GetString(EnvVarName("ExplorerAccessKey"))
}

// ExplorerSecret returns the secret for authenticating with explorer
func (c *generalConfig) ExplorerSecret() string {
	return c.viper.GetString(EnvVarName("ExplorerSecret"))
}

// TelemetryIngressURL returns the WSRPC URL for this node to push telemetry to, or nil.
func (c *generalConfig) TelemetryIngressURL() *url.URL {
	rval := c.getWithFallback("TelemetryIngressURL", ParseURL)
	switch t := rval.(type) {
	case nil:
		return nil
	case *url.URL:
		return t
	default:
		panic(fmt.Sprintf("invariant: TelemetryIngressURL returned as type %T", rval))
	}
}

// TelemetryIngressServerPubKey returns the public key to authenticate the telemetry ingress server
func (c *generalConfig) TelemetryIngressServerPubKey() string {
	return c.viper.GetString(EnvVarName("TelemetryIngressServerPubKey"))
}

// TelemetryIngressLogging toggles very verbose logging of raw telemetry messages for the TelemetryIngressClient
func (c *generalConfig) TelemetryIngressLogging() bool {
	return c.getWithFallback("TelemetryIngressLogging", ParseBool).(bool)
}

func (c *generalConfig) ORMMaxOpenConns() int {
	return int(c.getWithFallback("ORMMaxOpenConns", ParseUint16).(uint16))
}

func (c *generalConfig) ORMMaxIdleConns() int {
	return int(c.getWithFallback("ORMMaxIdleConns", ParseUint16).(uint16))
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

// LogToDisk configures disk preservation of logs.
func (c *generalConfig) LogToDisk() bool {
	return c.viper.GetBool(EnvVarName("LogToDisk"))
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

// LogSQLMigrations tells chainlink to log all SQL migrations made using the default logger
func (c *generalConfig) LogSQLMigrations() bool {
	return c.viper.GetBool(EnvVarName("LogSQLMigrations"))
}

// LogUnixTimestamps if set to true will log with timestamp in unix format, otherwise uses ISO8601
func (c *generalConfig) LogUnixTimestamps() bool {
	return c.viper.GetBool(EnvVarName("LogUnixTS"))
}

// Port represents the port Chainlink should listen on for client requests.
func (c *generalConfig) Port() uint16 {
	return c.getWithFallback("Port", ParseUint16).(uint16)
}

// DefaultChainID represents the chain ID which jobs will use if one is not explicitly specified
func (c *generalConfig) DefaultChainID() *big.Int {
	str := c.viper.GetString(EnvVarName("DefaultChainID"))
	if str != "" {
		v, err := ParseBigInt(str)
		if err != nil {
			logger.Errorw(
				"Ignoring invalid value provided for ETH_CHAIN_ID",
				"value", str,
				"error", err)
			return nil
		}
		return v.(*big.Int)

	}
	return nil
}

func (c *generalConfig) HTTPServerWriteTimeout() time.Duration {
	return c.getWithFallback("HTTPServerWriteTimeout", ParseDuration).(time.Duration)
}

// ReaperExpiration represents
func (c *generalConfig) ReaperExpiration() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("ReaperExpiration", ParseDuration).(time.Duration))
}

func (c *generalConfig) ReplayFromBlock() int64 {
	return c.viper.GetInt64(EnvVarName("ReplayFromBlock"))
}

// RootDir represents the location on the file system where Chainlink should
// keep its files.
func (c *generalConfig) RootDir() string {
	return c.getWithFallback("RootDir", ParseHomeDir).(string)
}

// RPID Fetches the RPID used for WebAuthn sessions. The RPID value should be the FQDN (localhost)
func (c *generalConfig) RPID() string {
	return c.viper.GetString(EnvVarName("RPID"))
}

// RPOrigin Fetches the RPOrigin used to configure WebAuthn sessions. The RPOrigin valiue should be
// the origin URL where WebAuthn requests initiate (http://localhost:6688/)
func (c *generalConfig) RPOrigin() string {
	return c.viper.GetString(EnvVarName("RPOrigin"))
}

// SecureCookies allows toggling of the secure cookies HTTP flag
func (c *generalConfig) SecureCookies() bool {
	return c.viper.GetBool(EnvVarName("SecureCookies"))
}

// SessionTimeout is the maximum duration that a user session can persist without any activity.
func (c *generalConfig) SessionTimeout() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("SessionTimeout", ParseDuration).(time.Duration))
}

// StatsPusherLogging toggles very verbose logging of raw messages for the StatsPusher (also telemetry)
func (c *generalConfig) StatsPusherLogging() bool {
	return c.getWithFallback("StatsPusherLogging", ParseBool).(bool)
}

// TLSCertPath represents the file system location of the TLS certificate
// Chainlink should use for HTTPS.
func (c *generalConfig) TLSCertPath() string {
	return c.viper.GetString(EnvVarName("TLSCertPath"))
}

// TLSHost represents the hostname to use for TLS clients. This should match
// the TLS certificate.
func (c *generalConfig) TLSHost() string {
	return c.viper.GetString(EnvVarName("TLSHost"))
}

// TLSKeyPath represents the file system location of the TLS key Chainlink
// should use for HTTPS.
func (c *generalConfig) TLSKeyPath() string {
	return c.viper.GetString(EnvVarName("TLSKeyPath"))
}

// TLSPort represents the port Chainlink should listen on for encrypted client requests.
func (c *generalConfig) TLSPort() uint16 {
	return c.getWithFallback("TLSPort", ParseUint16).(uint16)
}

// TLSRedirect forces TLS redirect for unencrypted connections
func (c *generalConfig) TLSRedirect() bool {
	return c.viper.GetBool(EnvVarName("TLSRedirect"))
}

// UnAuthenticatedRateLimit defines the threshold to which requests unauthenticated requests get limited
func (c *generalConfig) UnAuthenticatedRateLimit() int64 {
	return c.viper.GetInt64(EnvVarName("UnAuthenticatedRateLimit"))
}

// UnAuthenticatedRateLimitPeriod defines the period to which unauthenticated requests get limited
func (c *generalConfig) UnAuthenticatedRateLimitPeriod() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("UnAuthenticatedRateLimitPeriod", ParseDuration).(time.Duration))
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

// SessionSecret returns a sequence of bytes to be used as a private key for
// session signing or encryption.
func (c *generalConfig) SessionSecret() ([]byte, error) {
	return c.secretGenerator.Generate(c.RootDir())
}

// SessionOptions returns the sessions.Options struct used to configure
// the session store.
func (c *generalConfig) SessionOptions() sessions.Options {
	return sessions.Options{
		Secure:   c.SecureCookies(),
		HttpOnly: true,
		MaxAge:   86400 * 30,
	}
}

func (c *generalConfig) getWithFallback(name string, parser func(string) (interface{}, error)) interface{} {
	str := c.viper.GetString(EnvVarName(name))
	defaultValue, hasDefault := defaultValue(name)
	if str != "" {
		v, err := parser(str)
		if err == nil {
			return v
		}
		logger.Errorw(
			fmt.Sprintf("Invalid value provided for %s, falling back to default.", name),
			"value", str,
			"default", defaultValue,
			"error", err)
	}

	if !hasDefault {
		return zeroValue(name)
	}

	v, err := parser(defaultValue)
	if err != nil {
		log.Fatalf(`Invalid default for %s: "%s" (%s)`, name, defaultValue, err)
	}
	return v
}

// LogLevel determines the verbosity of the events to be logged.
type LogLevel struct {
	zapcore.Level
}

type DatabaseBackupMode string

var (
	DatabaseBackupModeNone DatabaseBackupMode = "none"
	DatabaseBackupModeLite DatabaseBackupMode = "lite"
	DatabaseBackupModeFull DatabaseBackupMode = "full"
)

func parseDatabaseBackupMode(s string) (interface{}, error) {
	switch DatabaseBackupMode(s) {
	case DatabaseBackupModeNone, DatabaseBackupModeLite, DatabaseBackupModeFull:
		return DatabaseBackupMode(s), nil
	default:
		return "", fmt.Errorf("unable to parse %v into DatabaseBackupMode. Must be one of values: \"%s\", \"%s\", \"%s\"", s, DatabaseBackupModeNone, DatabaseBackupModeLite, DatabaseBackupModeFull)
	}
}

func lookupEnv(k string, parse func(string) (interface{}, error)) (interface{}, bool) {
	s, ok := os.LookupEnv(k)
	if ok {
		val, err := parse(s)
		if err != nil {
			logger.Errorw(
				fmt.Sprintf("Invalid value provided for %s, falling back to default.", s),
				"value", s,
				"key", k,
				"error", err)
			return nil, false
		}
		return val, true
	}
	return nil, false
}

// EVM methods

func (*generalConfig) GlobalBalanceMonitorEnabled() (bool, bool) {
	val, ok := lookupEnv(EnvVarName("BalanceMonitorEnabled"), ParseBool)
	if val == nil {
		return false, false
	}
	return val.(bool), ok
}
func (*generalConfig) GlobalBlockEmissionIdleWarningThreshold() (time.Duration, bool) {
	val, ok := lookupEnv(EnvVarName("BlockEmissionIdleWarningThreshold"), ParseDuration)
	if val == nil {
		return 0, false
	}
	return val.(time.Duration), ok
}
func (*generalConfig) GlobalBlockHistoryEstimatorBatchSize() (uint32, bool) {
	val, ok := lookupEnv(EnvVarName("BlockHistoryEstimatorBatchSize"), ParseUint32)
	if val == nil {
		return 0, false
	}
	return val.(uint32), ok
}
func (*generalConfig) GlobalBlockHistoryEstimatorBlockDelay() (uint16, bool) {
	val, ok := lookupEnv(EnvVarName("BlockHistoryEstimatorBlockDelay"), ParseUint16)
	if val == nil {
		return 0, false
	}
	return val.(uint16), ok
}
func (*generalConfig) GlobalBlockHistoryEstimatorBlockHistorySize() (uint16, bool) {
	val, ok := lookupEnv(EnvVarName("BlockHistoryEstimatorBlockHistorySize"), ParseUint16)
	if val == nil {
		return 0, false
	}
	return val.(uint16), ok
}
func (*generalConfig) GlobalBlockHistoryEstimatorTransactionPercentile() (uint16, bool) {
	val, ok := lookupEnv(EnvVarName("BlockHistoryEstimatorTransactionPercentile"), ParseUint16)
	if val == nil {
		return 0, false
	}
	return val.(uint16), ok
}
func (*generalConfig) GlobalEthTxReaperInterval() (time.Duration, bool) {
	val, ok := lookupEnv(EnvVarName("EthTxReaperInterval"), ParseDuration)
	if val == nil {
		return 0, false
	}
	return val.(time.Duration), ok
}
func (*generalConfig) GlobalEthTxReaperThreshold() (time.Duration, bool) {
	val, ok := lookupEnv(EnvVarName("EthTxReaperThreshold"), ParseDuration)
	if val == nil {
		return 0, false
	}
	return val.(time.Duration), ok
}
func (*generalConfig) GlobalEthTxResendAfterThreshold() (time.Duration, bool) {
	val, ok := lookupEnv(EnvVarName("EthTxResendAfterThreshold"), ParseDuration)
	if val == nil {
		return 0, false
	}
	return val.(time.Duration), ok
}
func (*generalConfig) GlobalEvmDefaultBatchSize() (uint32, bool) {
	val, ok := lookupEnv(EnvVarName("EvmDefaultBatchSize"), ParseUint32)
	if val == nil {
		return 0, false
	}
	return val.(uint32), ok
}
func (*generalConfig) GlobalEvmFinalityDepth() (uint32, bool) {
	val, ok := lookupEnv(EnvVarName("EvmFinalityDepth"), ParseUint32)
	if val == nil {
		return 0, false
	}
	return val.(uint32), ok
}
func (*generalConfig) GlobalEvmGasBumpPercent() (uint16, bool) {
	val, ok := lookupEnv(EnvVarName("EvmGasBumpPercent"), ParseUint16)
	if val == nil {
		return 0, false
	}
	return val.(uint16), ok
}
func (*generalConfig) GlobalEvmGasBumpThreshold() (uint64, bool) {
	val, ok := lookupEnv(EnvVarName("EvmGasBumpThreshold"), ParseUint64)
	if val == nil {
		return 0, false
	}
	return val.(uint64), ok
}
func (*generalConfig) GlobalEvmGasBumpTxDepth() (uint16, bool) {
	val, ok := lookupEnv(EnvVarName("EvmGasBumpTxDepth"), ParseUint16)
	if val == nil {
		return 0, false
	}
	return val.(uint16), ok
}
func (*generalConfig) GlobalEvmGasBumpWei() (*big.Int, bool) {
	val, ok := lookupEnv(EnvVarName("EvmGasBumpWei"), ParseBigInt)
	if val == nil {
		return nil, false
	}
	return val.(*big.Int), ok
}
func (*generalConfig) GlobalEvmGasLimitDefault() (uint64, bool) {
	val, ok := lookupEnv(EnvVarName("EvmGasLimitDefault"), ParseUint64)
	if val == nil {
		return 0, false
	}
	return val.(uint64), ok
}
func (*generalConfig) GlobalEvmGasLimitMultiplier() (float32, bool) {
	val, ok := lookupEnv(EnvVarName("EvmGasLimitMultiplier"), ParseF32)
	if val == nil {
		return 0, false
	}
	return val.(float32), ok
}
func (*generalConfig) GlobalEvmGasLimitTransfer() (uint64, bool) {
	val, ok := lookupEnv(EnvVarName("EvmGasLimitTransfer"), ParseUint64)
	if val == nil {
		return 0, false
	}
	return val.(uint64), ok
}
func (*generalConfig) GlobalEvmGasPriceDefault() (*big.Int, bool) {
	val, ok := lookupEnv(EnvVarName("EvmGasPriceDefault"), ParseBigInt)
	if val == nil {
		return nil, false
	}
	return val.(*big.Int), ok
}
func (*generalConfig) GlobalEvmHeadTrackerHistoryDepth() (uint32, bool) {
	val, ok := lookupEnv(EnvVarName("EvmHeadTrackerHistoryDepth"), ParseUint32)
	if val == nil {
		return 0, false
	}
	return val.(uint32), ok
}
func (*generalConfig) GlobalEvmHeadTrackerMaxBufferSize() (uint32, bool) {
	val, ok := lookupEnv(EnvVarName("EvmHeadTrackerMaxBufferSize"), ParseUint32)
	if val == nil {
		return 0, false
	}
	return val.(uint32), ok
}
func (*generalConfig) GlobalEvmHeadTrackerSamplingInterval() (time.Duration, bool) {
	val, ok := lookupEnv(EnvVarName("EvmHeadTrackerSamplingInterval"), ParseDuration)
	if val == nil {
		return 0, false
	}
	return val.(time.Duration), ok
}
func (*generalConfig) GlobalEvmLogBackfillBatchSize() (uint32, bool) {
	val, ok := lookupEnv(EnvVarName("EvmLogBackfillBatchSize"), ParseUint32)
	if val == nil {
		return 0, false
	}
	return val.(uint32), ok
}
func (*generalConfig) GlobalEvmMaxGasPriceWei() (*big.Int, bool) {
	val, ok := lookupEnv(EnvVarName("EvmMaxGasPriceWei"), ParseBigInt)
	if val == nil {
		return nil, false
	}
	return val.(*big.Int), ok
}
func (*generalConfig) GlobalEvmMaxInFlightTransactions() (uint32, bool) {
	val, ok := lookupEnv(EnvVarName("EvmMaxInFlightTransactions"), ParseUint32)
	if val == nil {
		return 0, false
	}
	return val.(uint32), ok
}
func (*generalConfig) GlobalEvmMaxQueuedTransactions() (uint64, bool) {
	val, ok := lookupEnv(EnvVarName("EvmMaxQueuedTransactions"), ParseUint64)
	if val == nil {
		return 0, false
	}
	return val.(uint64), ok
}
func (*generalConfig) GlobalEvmMinGasPriceWei() (*big.Int, bool) {
	val, ok := lookupEnv(EnvVarName("EvmMinGasPriceWei"), ParseBigInt)
	if val == nil {
		return nil, false
	}
	return val.(*big.Int), ok
}
func (*generalConfig) GlobalEvmNonceAutoSync() (bool, bool) {
	val, ok := lookupEnv(EnvVarName("EvmNonceAutoSync"), ParseBool)
	if val == nil {
		return false, false
	}
	return val.(bool), ok
}
func (*generalConfig) GlobalEvmRPCDefaultBatchSize() (uint32, bool) {
	val, ok := lookupEnv(EnvVarName("EvmRPCDefaultBatchSize"), ParseUint32)
	if val == nil {
		return 0, false
	}
	return val.(uint32), ok
}
func (*generalConfig) GlobalFlagsContractAddress() (string, bool) {
	val, ok := lookupEnv(EnvVarName("FlagsContractAddress"), ParseString)
	if val == nil {
		return "", false
	}
	return val.(string), ok
}
func (*generalConfig) GlobalGasEstimatorMode() (string, bool) {
	val, ok := lookupEnv(EnvVarName("GasEstimatorMode"), ParseString)
	if val == nil {
		return "", false
	}
	return val.(string), ok
}
func (*generalConfig) GlobalChainType() (string, bool) {
	val, ok := lookupEnv(EnvVarName("ChainType"), ParseString)
	if val == nil {
		return "", false
	}
	return val.(string), ok
}
func (*generalConfig) GlobalLinkContractAddress() (string, bool) {
	val, ok := lookupEnv(EnvVarName("LinkContractAddress"), ParseString)
	if val == nil {
		return "", false
	}
	return val.(string), ok
}
func (*generalConfig) GlobalMinIncomingConfirmations() (uint32, bool) {
	val, ok := lookupEnv(EnvVarName("MinIncomingConfirmations"), ParseUint32)
	if val == nil {
		return 0, false
	}
	return val.(uint32), ok
}
func (*generalConfig) GlobalMinRequiredOutgoingConfirmations() (uint64, bool) {
	val, ok := lookupEnv(EnvVarName("MinRequiredOutgoingConfirmations"), ParseUint64)
	if val == nil {
		return 0, false
	}
	return val.(uint64), ok
}
func (*generalConfig) GlobalMinimumContractPayment() (*assets.Link, bool) {
	val, ok := lookupEnv(EnvVarName("MinimumContractPayment"), ParseLink)
	if val == nil {
		return nil, false
	}
	return val.(*assets.Link), ok
}
func (*generalConfig) GlobalOCRContractTransmitterTransmitTimeout() (time.Duration, bool) {
	val, ok := lookupEnv(EnvVarName("OCRContractTransmitterTransmitTimeout"), ParseDuration)
	if val == nil {
		return 0, false
	}
	return val.(time.Duration), ok
}
func (*generalConfig) GlobalOCRDatabaseTimeout() (time.Duration, bool) {
	val, ok := lookupEnv(EnvVarName("OCRDatabaseTimeout"), ParseDuration)
	if val == nil {
		return 0, false
	}
	return val.(time.Duration), ok
}
func (*generalConfig) GlobalOCRObservationGracePeriod() (time.Duration, bool) {
	val, ok := lookupEnv(EnvVarName("OCRObservationGracePeriod"), ParseDuration)
	if val == nil {
		return 0, false
	}
	return val.(time.Duration), ok
}
func (*generalConfig) GlobalOCRContractConfirmations() (uint16, bool) {
	val, ok := lookupEnv(EnvVarName("OCRContractConfirmations"), ParseUint16)
	if val == nil {
		return 0, false
	}
	return val.(uint16), ok
}
func (*generalConfig) GlobalOCR2ContractConfirmations() (uint16, bool) {
	val, ok := lookupEnv(EnvVarName("OCR2ContractConfirmations"), ParseUint16)
	if val == nil {
		return 0, false
	}
	return val.(uint16), ok
}
func (*generalConfig) GlobalEvmEIP1559DynamicFees() (bool, bool) {
	val, ok := lookupEnv(EnvVarName("EvmEIP1559DynamicFees"), ParseBool)
	if val == nil {
		return false, false
	}
	return val.(bool), ok
}
func (*generalConfig) GlobalEvmGasTipCapDefault() (*big.Int, bool) {
	val, ok := lookupEnv(EnvVarName("EvmGasTipCapDefault"), ParseBigInt)
	if val == nil {
		return nil, false
	}
	return val.(*big.Int), ok
}
func (*generalConfig) GlobalEvmGasTipCapMinimum() (*big.Int, bool) {
	val, ok := lookupEnv(EnvVarName("EvmGasTipCapMinimum"), ParseBigInt)
	if val == nil {
		return nil, false
	}
	return val.(*big.Int), ok
}

// UseLegacyEthEnvVars will upsert a new chain using the DefaultChainID and
// upsert nodes corresponding to the given ETH_URL and ETH_SECONDARY_URLS
func (c *generalConfig) UseLegacyEthEnvVars() bool {
	return c.viper.GetBool(EnvVarName("UseLegacyEthEnvVars"))
}

// DatabaseLockingMode can be one of 'dual', 'advisorylock', 'lease' or 'none'
// It controls which mode to use to enforce that only one Chainlink application can use the database
func (c *generalConfig) DatabaseLockingMode() string {
	return c.getWithFallback("DatabaseLockingMode", ParseString).(string)
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
