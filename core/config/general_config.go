package config

import (
	"fmt"
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
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/config/envvar"
	"github.com/smartcontractkit/chainlink/core/config/parse"
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
	ErrUnset   = errors.New("env var unset")
	ErrInvalid = errors.New("env var invalid")

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

	AutoPprofEnabled() bool
	EVMEnabled() bool
	EVMRPCEnabled() bool
	KeeperCheckUpkeepGasPriceFeatureEnabled() bool
	P2PEnabled() bool
	SolanaEnabled() bool
	TerraEnabled() bool
}

type GeneralOnlyConfig interface {
	Validate() error
	SetLogLevel(lvl zapcore.Level) error
	SetLogSQL(logSQL bool)

	FeatureFlags

	AdminCredentialsFile() string
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
	CertFile() string
	ClientNodeURL() string
	DatabaseBackupDir() string
	DatabaseBackupFrequency() time.Duration
	DatabaseBackupMode() DatabaseBackupMode
	DatabaseBackupOnVersionUpgrade() bool
	DatabaseBackupURL() *url.URL
	DatabaseListenerMaxReconnectDuration() time.Duration
	DatabaseListenerMinReconnectInterval() time.Duration
	DatabaseLockingMode() string
	DatabaseURL() url.URL
	DefaultChainID() *big.Int
	DefaultHTTPAllowUnrestrictedNetworkAccess() bool
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
	GetAdvisoryLockIDConfiguredOrDefault() int64
	GetDatabaseDialectConfiguredOrDefault() dialects.DialectName
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
	LeaseLockDuration() time.Duration
	LeaseLockRefreshInterval() time.Duration
	LogFileDir() string
	LogLevel() zapcore.Level
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
	RootDir() string
	SecureCookies() bool
	SessionOptions() sessions.Options
	SessionSecret() ([]byte, error)
	SessionTimeout() models.Duration
	TLSCertPath() string
	TLSDir() string
	TLSHost() string
	TLSKeyPath() string
	TLSPort() uint16
	TLSRedirect() bool
	TelemetryIngressLogging() bool
	TelemetryIngressServerPubKey() string
	TelemetryIngressURL() *url.URL
	TelemetryIngressBufferSize() uint
	TelemetryIngressMaxBatchSize() uint
	TelemetryIngressSendInterval() time.Duration
	TelemetryIngressUseBatchSend() bool
	TriggerFallbackDBPollInterval() time.Duration
	UnAuthenticatedRateLimit() int64
	UnAuthenticatedRateLimitPeriod() models.Duration
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
	GlobalBlockHistoryEstimatorTransactionPercentile() (uint16, bool)
	GlobalChainType() (string, bool)
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
	GlobalEvmGasFeeCapDefault() (*big.Int, bool)
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
	GlobalLinkContractAddress() (string, bool)
	GlobalMinIncomingConfirmations() (uint32, bool)
	GlobalMinRequiredOutgoingConfirmations() (uint64, bool)
	GlobalMinimumContractPayment() (*assets.Link, bool)

	OCR1Config
	OCR2Config
	P2PNetworking
	P2PV1Networking
	P2PV2Networking
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
	lggr             logger.Logger
	viper            *viper.Viper
	secretGenerator  SecretGenerator
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
}

// NewGeneralConfig returns the config with the environment variables set to their
// respective fields, or their defaults if environment variables are not set.
func NewGeneralConfig(lggr logger.Logger) GeneralConfig {
	v := viper.New()
	c := newGeneralConfigWithViper(v, lggr.Named("GeneralConfig"))
	c.secretGenerator = FilePersistedSecretGenerator{}
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

	ll, invalid := envvar.LogLevel.ParseLogLevel()
	if invalid != "" {
		lggr.Error(invalid)
	}
	config.defaultLogLevel = ll

	config.logLevel = config.defaultLogLevel
	config.logSQL = viper.GetBool(envvar.Name("LogSQL"))

	return
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

	if _, exists := os.LookupEnv("ETH_DISABLED"); exists {
		c.lggr.Warn(`DEPRECATION WARNING: ETH_DISABLED has been deprecated.

This warning will become a fatal error in a future release. Please switch to using one of the two options below instead:

- EVM_ENABLED=false - set this if you wish to completely disable all EVM chains and jobs and prevent them from ever loading (this is probably the one you want).
- EVM_RPC_ENABLED=false - set this if you wish to load all EVM chains and jobs, but prevent any RPC calls to the eth node (the old behaviour).
`)
	}
	if _, exists := os.LookupEnv("EVM_DISABLED"); exists {
		c.lggr.Warn(`DEPRECATION WARNING: EVM_DISABLED has been deprecated and superceded by EVM_ENABLED.

This warning will become a fatal error in a future release. Please use the following instead to disable EVM chains:

EVM_ENABLED=false
`)
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

	if c.EthereumURL() == "" {
		if c.EthereumHTTPURL() != nil {
			c.lggr.Warn("ETH_HTTP_URL has no effect when ETH_URL is not set")
		}
		if len(c.EthereumSecondaryURLs()) > 0 {
			c.lggr.Warn("ETH_SECONDARY_URL/ETH_SECONDARY_URLS have no effect when ETH_URL is not set")
		}
	}
	// Warn on legacy OCR env vars
	if c.OCRDHTLookupInterval() != 0 {
		c.lggr.Warn("OCR_DHT_LOOKUP_INTERVAL is deprecated, use P2P_DHT_LOOKUP_INTERVAL instead")
	}
	if c.OCRBootstrapCheckInterval() != 0 {
		c.lggr.Warn("OCR_BOOTSTRAP_CHECK_INTERVAL is deprecated, use P2P_BOOTSTRAP_CHECK_INTERVAL instead")
	}
	if c.OCRIncomingMessageBufferSize() != 0 {
		c.lggr.Warn("OCR_INCOMING_MESSAGE_BUFFER_SIZE is deprecated, use P2P_INCOMING_MESSAGE_BUFFER_SIZE instead")
	}
	if c.OCROutgoingMessageBufferSize() != 0 {
		c.lggr.Warn("OCR_OUTGOING_MESSAGE_BUFFER_SIZE is deprecated, use P2P_OUTGOING_MESSAGE_BUFFER_SIZE instead")
	}
	if c.OCRNewStreamTimeout() != 0 {
		c.lggr.Warn("OCR_NEW_STREAM_TIMEOUT is deprecated, use P2P_NEW_STREAM_TIMEOUT instead")
	}

	switch c.DatabaseLockingMode() {
	case "dual", "lease", "advisorylock", "none":
	default:
		return errors.Errorf("unrecognised value for DATABASE_LOCKING_MODE: %s (valid options are 'dual', 'lease', 'advisorylock' or 'none')", c.DatabaseLockingMode())
	}

	if c.LeaseLockRefreshInterval() > c.LeaseLockDuration()/2 {
		return errors.Errorf("LEASE_LOCK_REFRESH_INTERVAL must be less than or equal to half of LEASE_LOCK_DURATION (got LEASE_LOCK_REFRESH_INTERVAL=%s, LEASE_LOCK_DURATION=%s)", c.LeaseLockRefreshInterval().String(), c.LeaseLockDuration().String())
	}

	if c.viper.GetString(envvar.Name("LogFileDir")) != "" && !c.LogToDisk() {
		c.lggr.Warn("LOG_FILE_DIR is ignored and has no effect when LOG_TO_DISK is not enabled")
	}

	return nil
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

// AdminCredentialsFile points to text file containing admin credentials for logging in
func (c *generalConfig) AdminCredentialsFile() string {
	fieldName := "AdminCredentialsFile"
	file := c.viper.GetString(envvar.Name(fieldName))
	defaultValue, _ := envvar.DefaultValue(fieldName)
	if file == defaultValue {
		return filepath.Join(c.RootDir(), "apicredentials")
	}
	return file
}

// AuthenticatedRateLimit defines the threshold to which authenticated requests
// get limited. More than this many requests per AuthenticatedRateLimitPeriod will be rejected.
func (c *generalConfig) AuthenticatedRateLimit() int64 {
	return c.viper.GetInt64(envvar.Name("AuthenticatedRateLimit"))
}

// AuthenticatedRateLimitPeriod defines the period to which authenticated requests get limited
func (c *generalConfig) AuthenticatedRateLimitPeriod() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("AuthenticatedRateLimitPeriod", parse.Duration).(time.Duration))
}

func (c *generalConfig) AutoPprofEnabled() bool {
	return c.viper.GetBool(envvar.Name("AutoPprofEnabled"))
}

func (c *generalConfig) AutoPprofProfileRoot() string {
	root := c.viper.GetString(envvar.Name("AutoPprofProfileRoot"))
	if root == "" {
		return c.RootDir()
	}
	return root
}

func (c *generalConfig) AutoPprofPollInterval() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("AutoPprofPollInterval", parse.Duration).(time.Duration))
}

func (c *generalConfig) AutoPprofGatherDuration() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("AutoPprofGatherDuration", parse.Duration).(time.Duration))
}

func (c *generalConfig) AutoPprofGatherTraceDuration() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("AutoPprofGatherTraceDuration", parse.Duration).(time.Duration))
}

func (c *generalConfig) AutoPprofMaxProfileSize() utils.FileSize {
	return c.getWithFallback("AutoPprofMaxProfileSize", parse.FileSize).(utils.FileSize)
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
	return c.getWithFallback("AutoPprofMemThreshold", parse.FileSize).(utils.FileSize)
}

func (c *generalConfig) AutoPprofGoroutineThreshold() int {
	return c.viper.GetInt(envvar.Name("AutoPprofGoroutineThreshold"))
}

// BlockBackfillDepth specifies the number of blocks before the current HEAD that the
// log broadcaster will try to re-consume logs from
func (c *generalConfig) BlockBackfillDepth() uint64 {
	return c.getWithFallback("BlockBackfillDepth", parse.Uint64).(uint64)
}

// BlockBackfillSkip enables skipping of very long log backfills
func (c *generalConfig) BlockBackfillSkip() bool {
	return c.getWithFallback("BlockBackfillSkip", parse.Bool).(bool)
}

// BridgeResponseURL represents the URL for bridges to send a response to.
func (c *generalConfig) BridgeResponseURL() *url.URL {
	return c.getWithFallback("BridgeResponseURL", parse.URL).(*url.URL)
}

// ClientNodeURL is the URL of the Ethereum node this Chainlink node should connect to.
func (c *generalConfig) ClientNodeURL() string {
	return c.viper.GetString(envvar.Name("ClientNodeURL"))
}

// FeatureUICSAKeys enables the CSA Keys UI Feature.
func (c *generalConfig) FeatureUICSAKeys() bool {
	return c.getWithFallback("FeatureUICSAKeys", parse.Bool).(bool)
}

func (c *generalConfig) DatabaseListenerMinReconnectInterval() time.Duration {
	return c.getWithFallback("DatabaseListenerMinReconnectInterval", parse.Duration).(time.Duration)
}

func (c *generalConfig) DatabaseListenerMaxReconnectDuration() time.Duration {
	return c.getWithFallback("DatabaseListenerMaxReconnectDuration", parse.Duration).(time.Duration)
}

// DatabaseBackupMode sets the database backup mode
func (c *generalConfig) DatabaseBackupMode() DatabaseBackupMode {
	return c.getWithFallback("DatabaseBackupMode", parseDatabaseBackupMode).(DatabaseBackupMode)
}

// DatabaseBackupFrequency turns on the periodic database backup if set to a positive value
// DatabaseBackupMode must be then set to a value other than "none"
func (c *generalConfig) DatabaseBackupFrequency() time.Duration {
	return c.getWithFallback("DatabaseBackupFrequency", parse.Duration).(time.Duration)
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
	return c.getWithFallback("DatabaseBackupOnVersionUpgrade", parse.Bool).(bool)
}

// DatabaseBackupDir configures the directory for saving the backup file, if it's to be different from default one located in the RootDir
func (c *generalConfig) DatabaseBackupDir() string {
	return c.viper.GetString(envvar.Name("DatabaseBackupDir"))
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
	return models.MustMakeDuration(c.getWithFallback("DefaultHTTPTimeout", parse.Duration).(time.Duration))
}

// DefaultHTTPAllowUnrestrictedNetworkAccess controls whether http requests are unrestricted by default
// It is recommended that this be left disabled
func (c *generalConfig) DefaultHTTPAllowUnrestrictedNetworkAccess() bool {
	return c.viper.GetBool(envvar.Name("DefaultHTTPAllowUnrestrictedNetworkAccess"))
}

// Dev configures "development" mode for chainlink.
func (c *generalConfig) Dev() bool {
	return c.viper.GetBool(envvar.Name("Dev"))
}

// ShutdownGracePeriod is the maximum duration of graceful application shutdown.
// If exceeded, it will try closing DB lock and connection and exit immediately to avoid SIGKILL.
func (c *generalConfig) ShutdownGracePeriod() time.Duration {
	return c.getWithFallback("ShutdownGracePeriod", parse.Duration).(time.Duration)
}

// FeatureExternalInitiators enables the External Initiator feature.
func (c *generalConfig) FeatureExternalInitiators() bool {
	return c.viper.GetBool(envvar.Name("FeatureExternalInitiators"))
}

// FeatureFeedsManager enables the feeds manager
func (c *generalConfig) FeatureFeedsManager() bool {
	return c.viper.GetBool(envvar.Name("FeatureFeedsManager"))
}

// FeatureOffchainReporting enables the OCR job type.
func (c *generalConfig) FeatureOffchainReporting() bool {
	return c.getWithFallback("FeatureOffchainReporting", parse.Bool).(bool)
}

// FeatureOffchainReporting2 enables the OCR2 job type.
func (c *generalConfig) FeatureOffchainReporting2() bool {
	return c.getWithFallback("FeatureOffchainReporting2", parse.Bool).(bool)
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

// TerraEnabled allows Terra to be used
func (c *generalConfig) TerraEnabled() bool {
	return c.viper.GetBool(envvar.Name("TerraEnabled"))
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

// InsecureSkipVerify disables SSL certificate verification when connection to
// a chainlink client using the remote client, i.e. when executing most remote
// commands in the CLI.
//
// This is mostly useful for people who want to use TLS on localhost.
func (c *generalConfig) InsecureSkipVerify() bool {
	return c.viper.GetBool(envvar.Name("InsecureSkipVerify"))
}

func (c *generalConfig) TriggerFallbackDBPollInterval() time.Duration {
	return c.getWithFallback("TriggerFallbackDBPollInterval", parse.Duration).(time.Duration)
}

// JobPipelineMaxRunDuration is the maximum time that a job run may take
func (c *generalConfig) JobPipelineMaxRunDuration() time.Duration {
	return c.getWithFallback("JobPipelineMaxRunDuration", parse.Duration).(time.Duration)
}

func (c *generalConfig) JobPipelineResultWriteQueueDepth() uint64 {
	return c.getWithFallback("JobPipelineResultWriteQueueDepth", parse.Uint64).(uint64)
}

func (c *generalConfig) JobPipelineReaperInterval() time.Duration {
	return c.getWithFallback("JobPipelineReaperInterval", parse.Duration).(time.Duration)
}

func (c *generalConfig) JobPipelineReaperThreshold() time.Duration {
	return c.getWithFallback("JobPipelineReaperThreshold", parse.Duration).(time.Duration)
}

// KeeperRegistryCheckGasOverhead is the amount of extra gas to provide checkUpkeep() calls
// to account for the gas consumed by the keeper registry
func (c *generalConfig) KeeperRegistryCheckGasOverhead() uint64 {
	return c.getWithFallback("KeeperRegistryCheckGasOverhead", parse.Uint64).(uint64)
}

// KeeperRegistryPerformGasOverhead is the amount of extra gas to provide performUpkeep() calls
// to account for the gas consumed by the keeper registry
func (c *generalConfig) KeeperRegistryPerformGasOverhead() uint64 {
	return c.getWithFallback("KeeperRegistryPerformGasOverhead", parse.Uint64).(uint64)
}

// KeeperDefaultTransactionQueueDepth controls the queue size for DropOldestStrategy in Keeper
// Set to 0 to use SendEvery strategy instead
func (c *generalConfig) KeeperDefaultTransactionQueueDepth() uint32 {
	return c.viper.GetUint32(envvar.Name("KeeperDefaultTransactionQueueDepth"))
}

// KeeperGasPriceBufferPercent adds the specified percentage to the gas price
// used for checking whether to perform an upkeep. Only applies in legacy mode.
func (c *generalConfig) KeeperGasPriceBufferPercent() uint32 {
	return c.viper.GetUint32(envvar.Name("KeeperGasPriceBufferPercent"))
}

// KeeperGasTipCapBufferPercent adds the specified percentage to the gas price
// used for checking whether to perform an upkeep. Only applies in EIP-1559 mode.
func (c *generalConfig) KeeperGasTipCapBufferPercent() uint32 {
	return c.viper.GetUint32(envvar.Name("KeeperGasTipCapBufferPercent"))
}

// KeeperRegistrySyncInterval is the interval in which the RegistrySynchronizer performs a full
// sync of the keeper registry contract it is tracking
func (c *generalConfig) KeeperRegistrySyncInterval() time.Duration {
	return c.getWithFallback("KeeperRegistrySyncInterval", parse.Duration).(time.Duration)
}

// KeeperMaximumGracePeriod is the maximum number of blocks that a keeper will wait after performing
// an upkeep before it resumes checking that upkeep
func (c *generalConfig) KeeperMaximumGracePeriod() int64 {
	return c.viper.GetInt64(envvar.Name("KeeperMaximumGracePeriod"))
}

// KeeperRegistrySyncUpkeepQueueSize represents the maximum number of upkeeps that can be synced in parallel
func (c *generalConfig) KeeperRegistrySyncUpkeepQueueSize() uint32 {
	return c.getWithFallback("KeeperRegistrySyncUpkeepQueueSize", parse.Uint32).(uint32)
}

// KeeperCheckUpkeepGasPriceFeatureEnabled enables keepers to include a gas price when running checkUpkeep
func (c *generalConfig) KeeperCheckUpkeepGasPriceFeatureEnabled() bool {
	return c.getWithFallback("KeeperCheckUpkeepGasPriceFeatureEnabled", parse.Bool).(bool)
}

// JSONConsole when set to true causes logging to be made in JSON format
// If set to false, logs in console format
func (c *generalConfig) JSONConsole() bool {
	return c.getEnvWithFallback(envvar.JSONConsole).(bool)
}

// ExplorerURL returns the websocket URL for this node to push stats to, or nil.
func (c *generalConfig) ExplorerURL() *url.URL {
	rval := c.getWithFallback("ExplorerURL", parse.URL)
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
	return c.viper.GetString(envvar.Name("ExplorerAccessKey"))
}

// ExplorerSecret returns the secret for authenticating with explorer
func (c *generalConfig) ExplorerSecret() string {
	return c.viper.GetString(envvar.Name("ExplorerSecret"))
}

// TelemetryIngressURL returns the WSRPC URL for this node to push telemetry to, or nil.
func (c *generalConfig) TelemetryIngressURL() *url.URL {
	rval := c.getWithFallback("TelemetryIngressURL", parse.URL)
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

// TelemetryIngressUseBatchSend toggles sending telemetry using the batch client to the ingress server
func (c *generalConfig) TelemetryIngressUseBatchSend() bool {
	return c.viper.GetBool(envvar.Name("TelemetryIngressUseBatchSend"))
}

// TelemetryIngressLogging toggles very verbose logging of raw telemetry messages for the TelemetryIngressClient
func (c *generalConfig) TelemetryIngressLogging() bool {
	return c.getWithFallback("TelemetryIngressLogging", parse.Bool).(bool)
}

func (c *generalConfig) ORMMaxOpenConns() int {
	return int(c.getWithFallback("ORMMaxOpenConns", parse.Uint16).(uint16))
}

func (c *generalConfig) ORMMaxIdleConns() int {
	return int(c.getWithFallback("ORMMaxIdleConns", parse.Uint16).(uint16))
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
	return c.getEnvWithFallback(envvar.LogToDisk).(bool)
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
	return c.getEnvWithFallback(envvar.LogUnixTS).(bool)
}

// Port represents the port Chainlink should listen on for client requests.
func (c *generalConfig) Port() uint16 {
	return c.getWithFallback("Port", parse.Uint16).(uint16)
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
		return v.(*big.Int)

	}
	return nil
}

// HTTPServerWriteTimeout controls how long chainlink's API server may hold a
// socket open for writing a response to an HTTP request. This sometimes needs
// to be increased for pprof.
func (c *generalConfig) HTTPServerWriteTimeout() time.Duration {
	return c.getWithFallback("HTTPServerWriteTimeout", parse.Duration).(time.Duration)
}

// ReaperExpiration represents how long a session is held in the DB before being cleared
func (c *generalConfig) ReaperExpiration() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("ReaperExpiration", parse.Duration).(time.Duration))
}

// RootDir represents the location on the file system where Chainlink should
// keep its files.
func (c *generalConfig) RootDir() string {
	return c.getEnvWithFallback(envvar.RootDir).(string)
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
	return models.MustMakeDuration(c.getWithFallback("SessionTimeout", parse.Duration).(time.Duration))
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
	return c.getWithFallback("TLSPort", parse.Uint16).(uint16)
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
	return models.MustMakeDuration(c.getWithFallback("UnAuthenticatedRateLimitPeriod", parse.Duration).(time.Duration))
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

// Deprecated - prefer getEnvWithFallback with an EnvVar
func (c *generalConfig) getWithFallback(name string, parser func(string) (interface{}, error)) interface{} {
	return c.getEnvWithFallback(envvar.New(name, parser))
}

func (c *generalConfig) getEnvWithFallback(e *envvar.EnvVar) interface{} {
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

func parseDatabaseBackupMode(s string) (interface{}, error) {
	switch DatabaseBackupMode(s) {
	case DatabaseBackupModeNone, DatabaseBackupModeLite, DatabaseBackupModeFull:
		return DatabaseBackupMode(s), nil
	default:
		return "", fmt.Errorf("unable to parse %v into DatabaseBackupMode. Must be one of values: \"%s\", \"%s\", \"%s\"", s, DatabaseBackupModeNone, DatabaseBackupModeLite, DatabaseBackupModeFull)
	}
}

func (c *generalConfig) lookupEnv(k string, parse func(string) (interface{}, error)) (interface{}, bool) {
	s, ok := os.LookupEnv(k)
	if !ok {
		return nil, false
	}
	val, err := parse(s)
	if err == nil {
		return val, true
	}
	c.lggr.Errorw(fmt.Sprintf("Invalid value provided for %s, falling back to default.", s),
		"value", s, "key", k, "error", err)
	return nil, false
}

// EVM methods

func (c *generalConfig) GlobalBalanceMonitorEnabled() (bool, bool) {
	val, ok := c.lookupEnv(envvar.Name("BalanceMonitorEnabled"), parse.Bool)
	if val == nil {
		return false, false
	}
	return val.(bool), ok
}
func (c *generalConfig) GlobalBlockEmissionIdleWarningThreshold() (time.Duration, bool) {
	val, ok := c.lookupEnv(envvar.Name("BlockEmissionIdleWarningThreshold"), parse.Duration)
	if val == nil {
		return 0, false
	}
	return val.(time.Duration), ok
}
func (c *generalConfig) GlobalBlockHistoryEstimatorBatchSize() (uint32, bool) {
	val, ok := c.lookupEnv(envvar.Name("BlockHistoryEstimatorBatchSize"), parse.Uint32)
	if val == nil {
		return 0, false
	}
	return val.(uint32), ok
}
func (c *generalConfig) GlobalBlockHistoryEstimatorBlockDelay() (uint16, bool) {
	val, ok := c.lookupEnv(envvar.Name("BlockHistoryEstimatorBlockDelay"), parse.Uint16)
	if val == nil {
		return 0, false
	}
	return val.(uint16), ok
}
func (c *generalConfig) GlobalBlockHistoryEstimatorBlockHistorySize() (uint16, bool) {
	val, ok := c.lookupEnv(envvar.Name("BlockHistoryEstimatorBlockHistorySize"), parse.Uint16)
	if val == nil {
		return 0, false
	}
	return val.(uint16), ok
}
func (c *generalConfig) GlobalBlockHistoryEstimatorTransactionPercentile() (uint16, bool) {
	val, ok := c.lookupEnv(envvar.Name("BlockHistoryEstimatorTransactionPercentile"), parse.Uint16)
	if val == nil {
		return 0, false
	}
	return val.(uint16), ok
}
func (c *generalConfig) GlobalEthTxReaperInterval() (time.Duration, bool) {
	val, ok := c.lookupEnv(envvar.Name("EthTxReaperInterval"), parse.Duration)
	if val == nil {
		return 0, false
	}
	return val.(time.Duration), ok
}
func (c *generalConfig) GlobalEthTxReaperThreshold() (time.Duration, bool) {
	val, ok := c.lookupEnv(envvar.Name("EthTxReaperThreshold"), parse.Duration)
	if val == nil {
		return 0, false
	}
	return val.(time.Duration), ok
}
func (c *generalConfig) GlobalEthTxResendAfterThreshold() (time.Duration, bool) {
	val, ok := c.lookupEnv(envvar.Name("EthTxResendAfterThreshold"), parse.Duration)
	if val == nil {
		return 0, false
	}
	return val.(time.Duration), ok
}
func (c *generalConfig) GlobalEvmDefaultBatchSize() (uint32, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmDefaultBatchSize"), parse.Uint32)
	if val == nil {
		return 0, false
	}
	return val.(uint32), ok
}
func (c *generalConfig) GlobalEvmFinalityDepth() (uint32, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmFinalityDepth"), parse.Uint32)
	if val == nil {
		return 0, false
	}
	return val.(uint32), ok
}
func (c *generalConfig) GlobalEvmGasBumpPercent() (uint16, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmGasBumpPercent"), parse.Uint16)
	if val == nil {
		return 0, false
	}
	return val.(uint16), ok
}
func (c *generalConfig) GlobalEvmGasBumpThreshold() (uint64, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmGasBumpThreshold"), parse.Uint64)
	if val == nil {
		return 0, false
	}
	return val.(uint64), ok
}
func (c *generalConfig) GlobalEvmGasBumpTxDepth() (uint16, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmGasBumpTxDepth"), parse.Uint16)
	if val == nil {
		return 0, false
	}
	return val.(uint16), ok
}
func (c *generalConfig) GlobalEvmGasBumpWei() (*big.Int, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmGasBumpWei"), parse.BigInt)
	if val == nil {
		return nil, false
	}
	return val.(*big.Int), ok
}
func (c *generalConfig) GlobalEvmGasFeeCapDefault() (*big.Int, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmGasFeeCapDefault"), parse.BigInt)
	if val == nil {
		return nil, false
	}
	return val.(*big.Int), ok
}
func (c *generalConfig) GlobalBlockHistoryEstimatorEIP1559FeeCapBufferBlocks() (uint16, bool) {
	val, ok := c.lookupEnv(envvar.Name("BlockHistoryEstimatorEIP1559FeeCapBufferBlocks"), parse.Uint16)
	if val == nil {
		return 0, false
	}
	return val.(uint16), ok
}
func (c *generalConfig) GlobalEvmGasLimitDefault() (uint64, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmGasLimitDefault"), parse.Uint64)
	if val == nil {
		return 0, false
	}
	return val.(uint64), ok
}
func (c *generalConfig) GlobalEvmGasLimitMultiplier() (float32, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmGasLimitMultiplier"), parse.F32)
	if val == nil {
		return 0, false
	}
	return val.(float32), ok
}
func (c *generalConfig) GlobalEvmGasLimitTransfer() (uint64, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmGasLimitTransfer"), parse.Uint64)
	if val == nil {
		return 0, false
	}
	return val.(uint64), ok
}
func (c *generalConfig) GlobalEvmGasPriceDefault() (*big.Int, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmGasPriceDefault"), parse.BigInt)
	if val == nil {
		return nil, false
	}
	return val.(*big.Int), ok
}
func (c *generalConfig) GlobalEvmHeadTrackerHistoryDepth() (uint32, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmHeadTrackerHistoryDepth"), parse.Uint32)
	if val == nil {
		return 0, false
	}
	return val.(uint32), ok
}
func (c *generalConfig) GlobalEvmHeadTrackerMaxBufferSize() (uint32, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmHeadTrackerMaxBufferSize"), parse.Uint32)
	if val == nil {
		return 0, false
	}
	return val.(uint32), ok
}
func (c *generalConfig) GlobalEvmHeadTrackerSamplingInterval() (time.Duration, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmHeadTrackerSamplingInterval"), parse.Duration)
	if val == nil {
		return 0, false
	}
	return val.(time.Duration), ok
}
func (c *generalConfig) GlobalEvmLogBackfillBatchSize() (uint32, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmLogBackfillBatchSize"), parse.Uint32)
	if val == nil {
		return 0, false
	}
	return val.(uint32), ok
}
func (c *generalConfig) GlobalEvmMaxGasPriceWei() (*big.Int, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmMaxGasPriceWei"), parse.BigInt)
	if val == nil {
		return nil, false
	}
	return val.(*big.Int), ok
}
func (c *generalConfig) GlobalEvmMaxInFlightTransactions() (uint32, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmMaxInFlightTransactions"), parse.Uint32)
	if val == nil {
		return 0, false
	}
	return val.(uint32), ok
}
func (c *generalConfig) GlobalEvmMaxQueuedTransactions() (uint64, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmMaxQueuedTransactions"), parse.Uint64)
	if val == nil {
		return 0, false
	}
	return val.(uint64), ok
}
func (c *generalConfig) GlobalEvmMinGasPriceWei() (*big.Int, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmMinGasPriceWei"), parse.BigInt)
	if val == nil {
		return nil, false
	}
	return val.(*big.Int), ok
}
func (c *generalConfig) GlobalEvmNonceAutoSync() (bool, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmNonceAutoSync"), parse.Bool)
	if val == nil {
		return false, false
	}
	return val.(bool), ok
}
func (c *generalConfig) GlobalEvmRPCDefaultBatchSize() (uint32, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmRPCDefaultBatchSize"), parse.Uint32)
	if val == nil {
		return 0, false
	}
	return val.(uint32), ok
}
func (c *generalConfig) GlobalFlagsContractAddress() (string, bool) {
	val, ok := c.lookupEnv(envvar.Name("FlagsContractAddress"), parse.String)
	if val == nil {
		return "", false
	}
	return val.(string), ok
}
func (c *generalConfig) GlobalGasEstimatorMode() (string, bool) {
	val, ok := c.lookupEnv(envvar.Name("GasEstimatorMode"), parse.String)
	if val == nil {
		return "", false
	}
	return val.(string), ok
}

// GlobalChainType overrides all chains and forces them to act as a particular
// chain type. List of chain types is given in `chaintype.go`.
func (c *generalConfig) GlobalChainType() (string, bool) {
	val, ok := c.lookupEnv(envvar.Name("ChainType"), parse.String)
	if val == nil {
		return "", false
	}
	return val.(string), ok
}
func (c *generalConfig) GlobalLinkContractAddress() (string, bool) {
	val, ok := c.lookupEnv(envvar.Name("LinkContractAddress"), parse.String)
	if val == nil {
		return "", false
	}
	return val.(string), ok
}
func (c *generalConfig) GlobalMinIncomingConfirmations() (uint32, bool) {
	val, ok := c.lookupEnv(envvar.Name("MinIncomingConfirmations"), parse.Uint32)
	if val == nil {
		return 0, false
	}
	return val.(uint32), ok
}
func (c *generalConfig) GlobalMinRequiredOutgoingConfirmations() (uint64, bool) {
	val, ok := c.lookupEnv(envvar.Name("MinRequiredOutgoingConfirmations"), parse.Uint64)
	if val == nil {
		return 0, false
	}
	return val.(uint64), ok
}
func (c *generalConfig) GlobalMinimumContractPayment() (*assets.Link, bool) {
	val, ok := c.lookupEnv(envvar.Name("MinimumContractPayment"), parse.Link)
	if val == nil {
		return nil, false
	}
	return val.(*assets.Link), ok
}
func (c *generalConfig) GlobalEvmEIP1559DynamicFees() (bool, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmEIP1559DynamicFees"), parse.Bool)
	if val == nil {
		return false, false
	}
	return val.(bool), ok
}
func (c *generalConfig) GlobalEvmGasTipCapDefault() (*big.Int, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmGasTipCapDefault"), parse.BigInt)
	if val == nil {
		return nil, false
	}
	return val.(*big.Int), ok
}
func (c *generalConfig) GlobalEvmGasTipCapMinimum() (*big.Int, bool) {
	val, ok := c.lookupEnv(envvar.Name("EvmGasTipCapMinimum"), parse.BigInt)
	if val == nil {
		return nil, false
	}
	return val.(*big.Int), ok
}

// DatabaseLockingMode can be one of 'dual', 'advisorylock', 'lease' or 'none'
// It controls which mode to use to enforce that only one Chainlink application can use the database
func (c *generalConfig) DatabaseLockingMode() string {
	return c.getWithFallback("DatabaseLockingMode", parse.String).(string)
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
	return c.getWithFallback("AdvisoryLockID", parse.Int64).(int64)
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
