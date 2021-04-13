package orm

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/dialects"

	"gorm.io/gorm"

	"github.com/multiformats/go-multiaddr"

	ocr "github.com/smartcontractkit/libocr/offchainreporting"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	ethCore "github.com/ethereum/go-ethereum/core"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

// this permission grants read / write accccess to file owners only
const readWritePerms = os.FileMode(0600)

var (
	ErrUnset   = errors.New("env var unset")
	ErrInvalid = errors.New("env var invalid")

	configFileNotFoundError = reflect.TypeOf(viper.ConfigFileNotFoundError{})

	// keyed by ChainID
	ChainSpecificDefaults map[int64]ChainSpecificDefaultSet
	// If the chain is unknown, fallback to the general defaults
	GeneralDefaults ChainSpecificDefaultSet
)

type (
	// Config holds parameters used by the application which can be overridden by
	// setting environment variables.
	//
	// If you add an entry here which does not contain sensitive information, you
	// should also update presenters.ConfigWhitelist and cmd_test.TestClient_RunNodeShowsEnv.
	Config struct {
		viper           *viper.Viper
		SecretGenerator SecretGenerator
		runtimeStore    *ORM
		Dialect         dialects.DialectName
		AdvisoryLockID  int64
	}

	// ChainSpecificDefaultSet us a list of defaults specific to a particular chain ID
	ChainSpecificDefaultSet struct {
		EthGasBumpThreshold              uint64
		EthGasBumpWei                    *big.Int
		EthGasPriceDefault               *big.Int
		EthMaxGasPriceWei                *big.Int
		EthFinalityDepth                 uint
		EthHeadTrackerHistoryDepth       uint
		EthBalanceMonitorBlockDelay      uint16
		EthTxResendAfterThreshold        time.Duration
		GasUpdaterBlockDelay             uint16
		GasUpdaterBlockHistorySize       uint16
		HeadTimeBudget                   time.Duration
		MinIncomingConfirmations         uint32
		MinRequiredOutgoingConfirmations uint64
	}
)

func init() {
	ChainSpecificDefaults = make(map[int64]ChainSpecificDefaultSet)

	mainnet := ChainSpecificDefaultSet{
		EthGasBumpThreshold:              3,
		EthGasBumpWei:                    big.NewInt(5000000000),    // 5 Gwei
		EthGasPriceDefault:               big.NewInt(20000000000),   // 20 Gwei
		EthMaxGasPriceWei:                big.NewInt(1500000000000), // 1.5 Twei
		EthFinalityDepth:                 50,
		EthHeadTrackerHistoryDepth:       100,
		EthBalanceMonitorBlockDelay:      1,
		EthTxResendAfterThreshold:        30 * time.Second,
		GasUpdaterBlockDelay:             1,
		GasUpdaterBlockHistorySize:       24,
		HeadTimeBudget:                   13 * time.Second,
		MinIncomingConfirmations:         3,
		MinRequiredOutgoingConfirmations: 12,
	}

	// NOTE: There are probably other variables we can tweak for Kovan and other
	// test chains, but requires more in-depth research on their consensus
	// mechanisms. For now, mainnet defaults ought to be safe
	kovan := mainnet
	kovan.HeadTimeBudget = 4 * time.Second

	// BSC uses Clique consensus with ~3s block times
	// Clique offers finality within (N/2)+1 blocks where N is number of signers
	// There are 21 BSC validators so theoretically finality should occur after 21/2+1 = 11 blocks
	bscMainnet := ChainSpecificDefaultSet{
		EthGasBumpThreshold:              12,                       // mainnet * 4 (3s vs 13s block time)
		EthGasBumpWei:                    big.NewInt(5000000000),   // 5 Gwei
		EthGasPriceDefault:               big.NewInt(5000000000),   // 5 Gwei
		EthMaxGasPriceWei:                big.NewInt(500000000000), // 500 Gwei
		EthFinalityDepth:                 50,                       // Keeping this > 11 because it's not expensive and gives us a safety margin
		EthHeadTrackerHistoryDepth:       100,
		EthBalanceMonitorBlockDelay:      2,
		EthTxResendAfterThreshold:        15 * time.Second,
		GasUpdaterBlockDelay:             2,
		GasUpdaterBlockHistorySize:       24,
		HeadTimeBudget:                   3 * time.Second,
		MinIncomingConfirmations:         3,
		MinRequiredOutgoingConfirmations: 12,
	}

	hecoMainnet := bscMainnet

	// Matic has a 1s block time and looser finality guarantees than Ethereum.
	polygonMatic := ChainSpecificDefaultSet{
		EthGasBumpThreshold:              39,                       // mainnet * 13
		EthGasBumpWei:                    big.NewInt(5000000000),   // 5 Gwei
		EthGasPriceDefault:               big.NewInt(1000000000),   // 1 Gwei
		EthMaxGasPriceWei:                big.NewInt(500000000000), // 500 Gwei
		EthFinalityDepth:                 200,                      // A sprint is 64 blocks long and doesn't guarantee finality. To be safe, we take three sprints (192 blocks) plus a safety margin
		EthHeadTrackerHistoryDepth:       250,                      // EthFinalityDepth + safety margin
		EthBalanceMonitorBlockDelay:      13,                       // equivalent of 1 eth block seems reasonable
		EthTxResendAfterThreshold:        5 * time.Minute,          // 5 minutes is roughly 300 blocks on Matic. Since re-orgs occur often and can be deep, we want to avoid overloading the node with a ton of re-sent unconfirmed transactions.
		GasUpdaterBlockDelay:             32,                       // Delay needs to be large on matic since re-orgs are so frequent at the top level
		GasUpdaterBlockHistorySize:       128,
		HeadTimeBudget:                   1 * time.Second,
		MinIncomingConfirmations:         39, // mainnet * 13 (1s vs 13s block time)
		MinRequiredOutgoingConfirmations: 39, // mainnet * 13
	}

	GeneralDefaults = mainnet
	ChainSpecificDefaults[1] = mainnet
	ChainSpecificDefaults[42] = kovan
	ChainSpecificDefaults[56] = bscMainnet
	ChainSpecificDefaults[128] = hecoMainnet
	ChainSpecificDefaults[80001] = polygonMatic
}

func chainSpecificConfig(c Config) ChainSpecificDefaultSet {
	chainID := c.ChainID().Int64()
	if cset, exists := ChainSpecificDefaults[chainID]; exists {
		return cset
	}
	return GeneralDefaults
}

// NewConfig returns the config with the environment variables set to their
// respective fields, or their defaults if environment variables are not set.
func NewConfig() *Config {
	v := viper.New()
	return newConfigWithViper(v)
}

func newConfigWithViper(v *viper.Viper) *Config {
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

	config := &Config{
		viper:           v,
		SecretGenerator: filePersistedSecretGenerator{},
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

	return config
}

// Validate performs basic sanity checks on config and returns error if any
// misconfiguration would be fatal to the application
func (c *Config) Validate() error {
	ethGasBumpPercent := c.EthGasBumpPercent()
	if uint64(ethGasBumpPercent) < ethCore.DefaultTxPoolConfig.PriceBump {
		return errors.Errorf(
			"ETH_GAS_BUMP_PERCENT of %v may not be less than Geth's default of %v",
			c.EthGasBumpPercent(),
			ethCore.DefaultTxPoolConfig.PriceBump,
		)
	}

	if c.EthHeadTrackerHistoryDepth() < c.EthFinalityDepth() {
		return errors.New("ETH_HEAD_TRACKER_HISTORY_DEPTH must be equal to or greater than ETH_FINALITY_DEPTH")
	}

	if c.P2PAnnouncePort() != 0 && c.P2PAnnounceIP() == nil {
		return errors.Errorf("P2P_ANNOUNCE_PORT was given as %v but P2P_ANNOUNCE_IP was unset. You must also set P2P_ANNOUNCE_IP if P2P_ANNOUNCE_PORT is set", c.P2PAnnouncePort())
	}

	if c.FeatureOffchainReporting() && c.P2PListenPort() == 0 {
		return errors.New("P2P_LISTEN_PORT must be set to a non-zero value if FEATURE_OFFCHAIN_REPORTING is enabled")
	}

	var override time.Duration
	lc := ocrtypes.LocalConfig{
		BlockchainTimeout:                      c.OCRBlockchainTimeout(override),
		ContractConfigConfirmations:            c.OCRContractConfirmations(0),
		ContractConfigTrackerPollInterval:      c.OCRContractPollInterval(override),
		ContractConfigTrackerSubscribeInterval: c.OCRContractSubscribeInterval(override),
		ContractTransmitterTransmitTimeout:     c.OCRContractTransmitterTransmitTimeout(),
		DatabaseTimeout:                        c.OCRDatabaseTimeout(),
		DataSourceTimeout:                      c.OCRObservationTimeout(override),
		DataSourceGracePeriod:                  c.OCRObservationGracePeriod(),
	}
	if err := ocr.SanityCheckLocalConfig(lc); err != nil {
		return err
	}
	if _, err := c.P2PPeerID(nil); errors.Cause(err) == ErrInvalid {
		return err
	}
	if _, err := c.OCRKeyBundleID(nil); errors.Cause(err) == ErrInvalid {
		return err
	}
	if _, err := c.OCRTransmitterAddress(nil); errors.Cause(err) == ErrInvalid {
		return err
	}
	if peers, err := c.P2PBootstrapPeers(nil); err == nil {
		for i := range peers {
			if _, err := multiaddr.NewMultiaddr(peers[i]); err != nil {
				return errors.Errorf("p2p bootstrap peer %d is invalid: err %v", i, err)
			}
		}
	}
	if me := c.OCRMonitoringEndpoint(""); me != "" {
		if _, err := url.Parse(me); err != nil {
			return errors.Wrapf(err, "invalid monitoring url: %s", me)
		}
	}
	return nil
}

// SetRuntimeStore tells the configuration system to use a store for retrieving
// configuration variables that can be configured at runtime.
func (c *Config) SetRuntimeStore(orm *ORM) {
	c.runtimeStore = orm
}

// Set a specific configuration variable
func (c Config) Set(name string, value interface{}) {
	schemaT := reflect.TypeOf(ConfigSchema{})
	for index := 0; index < schemaT.NumField(); index++ {
		item := schemaT.FieldByIndex([]int{index})
		envName := item.Tag.Get("env")
		if envName == name {
			c.viper.Set(name, value)
			return
		}
	}
	logger.Panicf("No configuration parameter for %s", name)
}

const defaultPostgresAdvisoryLockID int64 = 1027321974924625846

func (c Config) GetAdvisoryLockIDConfiguredOrDefault() int64 {
	if c.AdvisoryLockID == 0 {
		return defaultPostgresAdvisoryLockID
	}
	return c.AdvisoryLockID
}

func (c Config) GetDatabaseDialectConfiguredOrDefault() dialects.DialectName {
	if c.Dialect == "" {
		return dialects.Postgres
	}
	return c.Dialect
}

// AllowOrigins returns the CORS hosts used by the frontend.
func (c Config) AllowOrigins() string {
	return c.viper.GetString(EnvVarName("AllowOrigins"))
}

// AdminCredentialsFile points to text file containing admnn credentials for logging in
func (c Config) AdminCredentialsFile() string {
	fieldName := "AdminCredentialsFile"
	file := c.viper.GetString(EnvVarName(fieldName))
	defaultValue, _ := defaultValue(fieldName)
	if file == defaultValue {
		return filepath.Join(c.RootDir(), "apicredentials")
	}
	return file
}

// AuthenticatedRateLimit defines the threshold to which requests authenticated requests get limited
func (c Config) AuthenticatedRateLimit() int64 {
	return c.viper.GetInt64(EnvVarName("AuthenticatedRateLimit"))
}

// AuthenticatedRateLimitPeriod defines the period to which authenticated requests get limited
func (c Config) AuthenticatedRateLimitPeriod() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("AuthenticatedRateLimitPeriod", parseDuration).(time.Duration))
}

// BalanceMonitorEnabled enables the balance monitor
func (c Config) BalanceMonitorEnabled() bool {
	return c.viper.GetBool(EnvVarName("BalanceMonitorEnabled"))
}

// BlockBackfillDepth specifies the number of blocks before the current HEAD that the
// log broadcaster will try to re-consume logs from
func (c Config) BlockBackfillDepth() uint64 {
	return c.getWithFallback("BlockBackfillDepth", parseUint64).(uint64)
}

// BridgeResponseURL represents the URL for bridges to send a response to.
func (c Config) BridgeResponseURL() *url.URL {
	return c.getWithFallback("BridgeResponseURL", parseURL).(*url.URL)
}

// ChainID represents the chain ID to use for transactions.
func (c Config) ChainID() *big.Int {
	return c.getWithFallback("ChainID", parseBigInt).(*big.Int)
}

// ClientNodeURL is the URL of the Ethereum node this Chainlink node should connect to.
func (c Config) ClientNodeURL() string {
	return c.viper.GetString(EnvVarName("ClientNodeURL"))
}

func (c Config) DatabaseListenerMinReconnectInterval() time.Duration {
	return c.getWithFallback("DatabaseListenerMinReconnectInterval", parseDuration).(time.Duration)
}

func (c Config) DatabaseListenerMaxReconnectDuration() time.Duration {
	return c.getWithFallback("DatabaseListenerMaxReconnectDuration", parseDuration).(time.Duration)
}

func (c Config) DatabaseMaximumTxDuration() time.Duration {
	return c.getWithFallback("DatabaseMaximumTxDuration", parseDuration).(time.Duration)
}

// DatabaseBackupMode sets the database backup mode
func (c Config) DatabaseBackupMode() DatabaseBackupMode {
	return c.getWithFallback("DatabaseBackupMode", parseDatabaseBackupMode).(DatabaseBackupMode)
}

// DatabaseBackupFrequency turns on the periodic database backup if set to a positive value
// DatabaseBackupMode must be then set to a value other than "none"
func (c Config) DatabaseBackupFrequency() time.Duration {
	return c.getWithFallback("DatabaseBackupFrequency", parseDuration).(time.Duration)
}

// DatabaseBackupURL configures the URL for the database to backup, if it's to be different from the main on
func (c Config) DatabaseBackupURL() *url.URL {
	s := c.viper.GetString(EnvVarName("DatabaseBackupURL"))
	if s == "" {
		return nil
	}
	uri, err := url.Parse(s)
	if err != nil {
		logger.Error("invalid database backup url %s", s)
		return nil
	}
	return uri
}

// DatabaseTimeout represents how long to tolerate non response from the DB.
func (c Config) DatabaseTimeout() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("DatabaseTimeout", parseDuration).(time.Duration))
}

// GlobalLockRetryInterval represents how long to wait before trying again to get the global advisory lock.
func (c Config) GlobalLockRetryInterval() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("GlobalLockRetryInterval", parseDuration).(time.Duration))
}

// DatabaseURL configures the URL for chainlink to connect to. This must be
// a properly formatted URL, with a valid scheme (postgres://)
func (c Config) DatabaseURL() url.URL {
	s := c.viper.GetString(EnvVarName("DatabaseURL"))
	uri, err := url.Parse(s)
	if err != nil {
		logger.Error("invalid database url %s", s)
		return url.URL{}
	}
	if uri.String() == "" {
		return *uri
	}
	static.SetConsumerName(uri, "Default")
	return *uri
}

// MigrateDatabase determines whether the database will be automatically
// migrated on application startup if set to true
func (c Config) MigrateDatabase() bool {
	return c.viper.GetBool(EnvVarName("MigrateDatabase"))
}

// DefaultMaxHTTPAttempts defines the limit for HTTP requests.
func (c Config) DefaultMaxHTTPAttempts() uint {
	return uint(c.getWithFallback("DefaultMaxHTTPAttempts", parseUint64).(uint64))
}

// DefaultHTTPLimit defines the size limit for HTTP requests and responses
func (c Config) DefaultHTTPLimit() int64 {
	return c.viper.GetInt64(EnvVarName("DefaultHTTPLimit"))
}

// DefaultHTTPTimeout defines the default timeout for http requests
func (c Config) DefaultHTTPTimeout() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("DefaultHTTPTimeout", parseDuration).(time.Duration))
}

// DefaultHTTPAllowUnrestrictedNetworkAccess controls whether http requests are unrestricted by default
// It is recommended that this be left disabled
func (c Config) DefaultHTTPAllowUnrestrictedNetworkAccess() bool {
	return c.viper.GetBool(EnvVarName("DefaultHTTPAllowUnrestrictedNetworkAccess"))
}

// Dev configures "development" mode for chainlink.
func (c Config) Dev() bool {
	return c.viper.GetBool(EnvVarName("Dev"))
}

// EnableExperimentalAdapters enables support for experimental adapters
func (c Config) EnableExperimentalAdapters() bool {
	return c.viper.GetBool(EnvVarName("EnableExperimentalAdapters"))
}

// FeatureExternalInitiators enables the External Initiator feature.
func (c Config) FeatureExternalInitiators() bool {
	return c.viper.GetBool(EnvVarName("FeatureExternalInitiators"))
}

// FeatureFluxMonitor enables the Flux Monitor feature.
func (c Config) FeatureFluxMonitor() bool {
	return c.viper.GetBool(EnvVarName("FeatureFluxMonitor"))
}

// FeatureFluxMonitorV2 enables the Flux Monitor v2 feature.
func (c Config) FeatureFluxMonitorV2() bool {
	return c.getWithFallback("FeatureFluxMonitorV2", parseBool).(bool)
}

// FeatureOffchainReporting enables the Flux Monitor feature.
func (c Config) FeatureOffchainReporting() bool {
	return c.viper.GetBool(EnvVarName("FeatureOffchainReporting"))
}

// MaximumServiceDuration is the maximum time that a service agreement can run
// from after the time it is created. Default 1 year = 365 * 24h = 8760h
func (c Config) MaximumServiceDuration() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("MaximumServiceDuration", parseDuration).(time.Duration))
}

// MinimumServiceDuration is the shortest duration from now that a service is
// allowed to run.
func (c Config) MinimumServiceDuration() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("MinimumServiceDuration", parseDuration).(time.Duration))
}

// EthBalanceMonitorBlockDelay is the number of blocks that the balance monitor
// trails behind head. This is required e.g. for Infura because they will often
// announce a new head, then route a request to a different node which does not
// have this head yet.
func (c Config) EthBalanceMonitorBlockDelay() uint16 {
	if c.viper.IsSet(EnvVarName("EthBalanceMonitorBlockDelay")) {
		return uint16(c.viper.GetUint32(EnvVarName("EthBalanceMonitorBlockDelay")))
	}
	return chainSpecificConfig(c).EthBalanceMonitorBlockDelay
}

// EthReceiptFetchBatchSize controls the number of receipts fetched in each
// request in the EthConfirmer
func (c Config) EthReceiptFetchBatchSize() uint32 {
	return c.viper.GetUint32(EnvVarName("EthReceiptFetchBatchSize"))
}

// EthGasBumpThreshold is the number of blocks to wait before bumping gas again on unconfirmed transactions
// Set to 0 to disable gas bumping
func (c Config) EthGasBumpThreshold() uint64 {
	if c.viper.IsSet(EnvVarName("EthGasBumpThreshold")) {
		return c.viper.GetUint64(EnvVarName("EthGasBumpThreshold"))
	}
	return chainSpecificConfig(c).EthGasBumpThreshold
}

// EthGasBumpTxDepth is the number of transactions to gas bump starting from oldest.
// Set to 0 for no limit (i.e. bump all)
func (c Config) EthGasBumpTxDepth() uint16 {
	return c.getWithFallback("EthGasBumpTxDepth", parseUint16).(uint16)
}

// EthGasBumpPercent is the minimum percentage by which gas is bumped on each transaction attempt
// Change with care since values below geth's default will fail with "underpriced replacement transaction"
func (c Config) EthGasBumpPercent() uint16 {
	return c.getWithFallback("EthGasBumpPercent", parseUint16).(uint16)
}

// EthGasBumpWei is the minimum fixed amount of wei by which gas is bumped on each transaction attempt
func (c Config) EthGasBumpWei() *big.Int {
	str := c.viper.GetString(EnvVarName("EthGasBumpWei"))
	if str != "" {
		n, err := parseBigInt(str)
		if err != nil {
			logger.Errorw(
				"Invalid value provided for EthGasBumpWei, falling back to default.",
				"value", str,
				"error", err)
		} else {
			return n.(*big.Int)
		}
	}
	return chainSpecificConfig(c).EthGasBumpWei
}

// EthMaxGasPriceWei is the maximum amount in Wei that a transaction will be
// bumped to before abandoning it and marking it as errored.
func (c Config) EthMaxGasPriceWei() *big.Int {

	str := c.viper.GetString(EnvVarName("EthMaxGasPriceWei"))
	if str != "" {
		n, err := parseBigInt(str)
		if err != nil {
			logger.Errorw(
				"Invalid value provided for EthMaxGasPriceWei, falling back to default.",
				"value", str,
				"error", err)
		} else {
			return n.(*big.Int)
		}
	}
	return chainSpecificConfig(c).EthMaxGasPriceWei
}

// EthMaxUnconfirmedTransactions is the maximum number of unconfirmed
// transactions per key that are allowed to be in flight before jobs will start
// failing and rejecting send of any further transactions.
// 0 value disables
func (c Config) EthMaxUnconfirmedTransactions() uint64 {
	return c.getWithFallback("EthMaxUnconfirmedTransactions", parseUint64).(uint64)
}

// EthGasLimitDefault sets the default gas limit for outgoing transactions.
func (c Config) EthGasLimitDefault() uint64 {
	return c.getWithFallback("EthGasLimitDefault", parseUint64).(uint64)
}

// EthGasPriceDefault is the starting gas price for every transaction
func (c Config) EthGasPriceDefault() *big.Int {
	if c.runtimeStore != nil {
		var value big.Int
		if err := c.runtimeStore.GetConfigValue("EthGasPriceDefault", &value); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warnw("Error while trying to fetch EthGasPriceDefault.", "error", err)
		} else if err == nil {
			return &value
		}
	}
	str := c.viper.GetString(EnvVarName("EthGasPriceDefault"))
	if str != "" {
		n, err := parseBigInt(str)
		if err != nil {
			logger.Errorw(
				"Invalid value provided for EthGasPriceDefault, falling back to default.",
				"value", str,
				"error", err)
		} else {
			return n.(*big.Int)
		}
	}
	return chainSpecificConfig(c).EthGasPriceDefault
}

// SetEthGasPriceDefault saves a runtime value for the default gas price for transactions
func (c Config) SetEthGasPriceDefault(value *big.Int) error {
	if c.runtimeStore == nil {
		return errors.New("No runtime store installed")
	}
	return c.runtimeStore.SetConfigValue("EthGasPriceDefault", value)
}

// EthFinalityDepth is the number of blocks after which an ethereum transaction is considered "final"
// BlocksConsideredFinal determines how deeply we look back to ensure that transactions are confirmed onto the longest chain
// There is not a large performance penalty to setting this relatively high (on the order of hundreds)
// It is practically limited by the number of heads we store in the database and should be less than this with a comfortable margin.
// If a transaction is mined in a block more than this many blocks ago, and is reorged out, we will NOT retransmit this transaction and undefined behaviour can occur including gaps in the nonce sequence that require manual intervention to fix.
// Therefore this number represents a number of blocks we consider large enough that no re-org this deep will ever feasibly happen.
func (c Config) EthFinalityDepth() uint {
	if c.viper.IsSet(EnvVarName("EthFinalityDepth")) {
		return uint(c.viper.GetUint64(EnvVarName("EthFinalityDepth")))
	}
	return chainSpecificConfig(c).EthFinalityDepth
}

// EthHeadTrackerHistoryDepth tracks the top N block numbers to keep in the `heads` database table.
// Note that this can easily result in MORE than N records since in the case of re-orgs we keep multiple heads for a particular block height.
// This number should be at least as large as `EthFinalityDepth`.
// There may be a small performance penalty to setting this to something very large (10,000+)
func (c Config) EthHeadTrackerHistoryDepth() uint {
	if c.viper.IsSet(EnvVarName("EthHeadTrackerHistoryDepth")) {
		return uint(c.viper.GetUint64(EnvVarName("EthHeadTrackerHistoryDepth")))
	}
	return chainSpecificConfig(c).EthHeadTrackerHistoryDepth
}

// EthHeadTrackerMaxBufferSize is the maximum number of heads that may be
// buffered in front of the head tracker before older heads start to be
// dropped. You may think of it as something like the maximum permittable "lag"
// for the head tracker before we start dropping heads to keep up.
func (c Config) EthHeadTrackerMaxBufferSize() uint {
	return uint(c.getWithFallback("EthHeadTrackerMaxBufferSize", parseUint64).(uint64))
}

// EthTxResendAfterThreshold controls how long the ethResender will wait before
// re-sending the latest eth_tx_attempt. This is designed a as a fallback to
// protect against the eth nodes dropping txes (it has been anecdotally
// observed to happen), networking issues or txes being ejected from the
// mempool.
// See eth_resender.go for more details
func (c Config) EthTxResendAfterThreshold() time.Duration {
	str := c.viper.GetString(EnvVarName("EthTxResendAfterThreshold"))
	if str != "" {
		n, err := parseDuration(str)
		if err != nil {
			logger.Errorw(
				"Invalid value provided for EthTxResendAfterThreshold, falling back to default.",
				"value", str,
				"error", err)
		} else {
			return n.(time.Duration)
		}
	}
	return chainSpecificConfig(c).EthTxResendAfterThreshold
}

// EthLogBackfillBatchSize sets the batch size for calling FilterLogs when we backfill missing logs
func (c Config) EthLogBackfillBatchSize() uint32 {
	return c.getWithFallback("EthLogBackfillBatchSize", parseUint32).(uint32)
}

// EthereumURL represents the URL of the Ethereum node to connect Chainlink to.
func (c Config) EthereumURL() string {
	return c.viper.GetString(EnvVarName("EthereumURL"))
}

// EthereumSecondaryURLs is an optional backup RPC URL
// Must be http(s) format
// If specified, transactions will also be broadcast to this ethereum node
func (c Config) EthereumSecondaryURLs() []url.URL {
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

// EthereumDisabled shows whether Ethereum interactions are supported.
func (c Config) EthereumDisabled() bool {
	return c.viper.GetBool(EnvVarName("EthereumDisabled"))
}

// FlagsContractAddress represents the Flags contract address
func (c Config) FlagsContractAddress() string {
	return c.viper.GetString(EnvVarName("FlagsContractAddress"))
}

// GasUpdaterBlockDelay is the number of blocks that the gas updater trails behind head.
// E.g. if this is set to 3, and we receive block 10, gas updater will
// fetch block 7.
// CAUTION: You might be tempted to set this to 0 to use the latest possible
// block, but it is possible to receive a head BEFORE that block is actually
// available from the connected node via RPC. In this case you will get false
// "zero" blocks that are missing transactions.
func (c Config) GasUpdaterBlockDelay() uint16 {
	if c.viper.IsSet(EnvVarName("GasUpdaterBlockDelay")) {
		return uint16(c.viper.GetUint32(EnvVarName("GasUpdaterBlockDelay")))
	}
	return chainSpecificConfig(c).GasUpdaterBlockDelay
}

// GasUpdaterBlockHistorySize is the number of past blocks to keep in memory to
// use as a basis for calculating a percentile gas price
func (c Config) GasUpdaterBlockHistorySize() uint16 {
	if c.viper.IsSet(EnvVarName("GasUpdaterBlockHistorySize")) {
		return uint16(c.viper.GetUint32(EnvVarName("GasUpdaterBlockHistorySize")))
	}
	return chainSpecificConfig(c).GasUpdaterBlockHistorySize
}

// GasUpdaterTransactionPercentile is the percentile gas price to choose. E.g.
// if the past transaction history contains four transactions with gas prices:
// [100, 200, 300, 400], picking 25 for this number will give a value of 200
func (c Config) GasUpdaterTransactionPercentile() uint16 {
	return c.getWithFallback("GasUpdaterTransactionPercentile", parseUint16).(uint16)
}

// GasUpdaterEnabled turns on the automatic gas updater if set to true
// It is disabled by default
func (c Config) GasUpdaterEnabled() bool {
	return c.viper.GetBool(EnvVarName("GasUpdaterEnabled"))
}

// InsecureFastScrypt causes all key stores to encrypt using "fast" scrypt params instead
// This is insecure and only useful for local testing. DO NOT SET THIS IN PRODUCTION
func (c Config) InsecureFastScrypt() bool {
	return c.viper.GetBool(EnvVarName("InsecureFastScrypt"))
}

func (c Config) TriggerFallbackDBPollInterval() time.Duration {
	return c.getWithFallback("TriggerFallbackDBPollInterval", parseDuration).(time.Duration)
}

// JobPipelineMaxRunDuration is the maximum time that a job run may take
func (c Config) JobPipelineMaxRunDuration() time.Duration {
	return c.getWithFallback("JobPipelineMaxRunDuration", parseDuration).(time.Duration)
}

func (c Config) JobPipelineResultWriteQueueDepth() uint64 {
	return c.getWithFallback("JobPipelineResultWriteQueueDepth", parseUint64).(uint64)
}

// JobPipelineParallelism controls how many workers the pipeline.Runner
// uses in parallel (how many pipeline runs may simultaneously be executing)
func (c Config) JobPipelineParallelism() uint8 {
	return c.getWithFallback("JobPipelineParallelism", parseUint8).(uint8)
}

func (c Config) JobPipelineReaperInterval() time.Duration {
	return c.getWithFallback("JobPipelineReaperInterval", parseDuration).(time.Duration)
}

func (c Config) JobPipelineReaperThreshold() time.Duration {
	return c.getWithFallback("JobPipelineReaperThreshold", parseDuration).(time.Duration)
}

func (c Config) KeeperRegistrySyncInterval() time.Duration {
	return c.getWithFallback("KeeperRegistrySyncInterval", parseDuration).(time.Duration)
}

func (c Config) KeeperMinimumRequiredConfirmations() uint64 {
	return c.viper.GetUint64(EnvVarName("KeeperMinimumRequiredConfirmations"))
}

func (c Config) KeeperMaximumGracePeriod() int64 {
	return c.viper.GetInt64(EnvVarName("KeeperMaximumGracePeriod"))
}

// JSONConsole enables the JSON console.
func (c Config) JSONConsole() bool {
	return c.viper.GetBool(EnvVarName("JSONConsole"))
}

// LinkContractAddress represents the address
func (c Config) LinkContractAddress() string {
	return c.viper.GetString(EnvVarName("LinkContractAddress"))
}

// ExplorerURL returns the websocket URL for this node to push stats to, or nil.
func (c Config) ExplorerURL() *url.URL {
	rval := c.getWithFallback("ExplorerURL", parseURL)
	switch t := rval.(type) {
	case nil:
		return nil
	case *url.URL:
		return t
	default:
		logger.Panicf("invariant: ExplorerURL returned as type %T", rval)
		return nil
	}
}

// ExplorerAccessKey returns the access key for authenticating with explorer
func (c Config) ExplorerAccessKey() string {
	return c.viper.GetString(EnvVarName("ExplorerAccessKey"))
}

// ExplorerSecret returns the secret for authenticating with explorer
func (c Config) ExplorerSecret() string {
	return c.viper.GetString(EnvVarName("ExplorerSecret"))
}

// FIXME: Add comments to all of these
func (c Config) OCRBootstrapCheckInterval() time.Duration {
	return c.getWithFallback("OCRBootstrapCheckInterval", parseDuration).(time.Duration)
}

func (c Config) OCRContractTransmitterTransmitTimeout() time.Duration {
	return c.getWithFallback("OCRContractTransmitterTransmitTimeout", parseDuration).(time.Duration)
}

func (c Config) getDurationWithOverride(override time.Duration, field string) time.Duration {
	if override != time.Duration(0) {
		return override
	}
	return c.getWithFallback(field, parseDuration).(time.Duration)
}

func (c Config) OCRObservationTimeout(override time.Duration) time.Duration {
	return c.getDurationWithOverride(override, "OCRObservationTimeout")
}

func (c Config) OCRObservationGracePeriod() time.Duration {
	return c.getWithFallback("OCRObservationGracePeriod", parseDuration).(time.Duration)
}

func (c Config) OCRBlockchainTimeout(override time.Duration) time.Duration {
	return c.getDurationWithOverride(override, "OCRBlockchainTimeout")
}

func (c Config) OCRContractSubscribeInterval(override time.Duration) time.Duration {
	return c.getDurationWithOverride(override, "OCRContractSubscribeInterval")
}

func (c Config) OCRContractPollInterval(override time.Duration) time.Duration {
	return c.getDurationWithOverride(override, "OCRContractPollInterval")
}

func (c Config) OCRContractConfirmations(override uint16) uint16 {
	if override != uint16(0) {
		return override
	}
	return c.getWithFallback("OCRContractConfirmations", parseUint16).(uint16)
}

func (c Config) OCRDatabaseTimeout() time.Duration {
	return c.getWithFallback("OCRDatabaseTimeout", parseDuration).(time.Duration)
}

func (c Config) OCRDHTLookupInterval() int {
	return int(c.getWithFallback("OCRDHTLookupInterval", parseUint16).(uint16))
}

func (c Config) OCRIncomingMessageBufferSize() int {
	return int(c.getWithFallback("OCRIncomingMessageBufferSize", parseUint16).(uint16))
}

func (c Config) OCRNewStreamTimeout() time.Duration {
	return c.getWithFallback("OCRNewStreamTimeout", parseDuration).(time.Duration)
}

func (c Config) OCROutgoingMessageBufferSize() int {
	return int(c.getWithFallback("OCRIncomingMessageBufferSize", parseUint16).(uint16))
}

// OCRTraceLogging determines whether OCR logs at TRACE level are enabled. The
// option to turn them off is given because they can be very verbose
func (c Config) OCRTraceLogging() bool {
	return c.viper.GetBool(EnvVarName("OCRTraceLogging"))
}

func (c Config) OCRMonitoringEndpoint(override string) string {
	if override != "" {
		return override
	}
	return c.viper.GetString(EnvVarName("OCRMonitoringEndpoint"))
}

func (c Config) OCRTransmitterAddress(override *models.EIP55Address) (models.EIP55Address, error) {
	if override != nil {
		return *override, nil
	}
	taStr := c.viper.GetString(EnvVarName("OCRTransmitterAddress"))
	if taStr != "" {
		ta, err := models.NewEIP55Address(taStr)
		if err != nil {
			return "", errors.Wrapf(ErrInvalid, "OCR_TRANSMITTER_ADDRESS is invalid EIP55 %v", err)
		}
		return ta, nil
	}
	return "", errors.Wrap(ErrUnset, "OCR_TRANSMITTER_ADDRESS")
}

func (c Config) OCRKeyBundleID(override *models.Sha256Hash) (models.Sha256Hash, error) {
	if override != nil {
		return *override, nil
	}
	kbStr := c.viper.GetString(EnvVarName("OCRKeyBundleID"))
	if kbStr != "" {
		kb, err := models.Sha256HashFromHex(kbStr)
		if err != nil {
			return models.Sha256Hash{}, errors.Wrapf(ErrInvalid, "OCR_KEY_BUNDLE_ID is an invalid sha256 hash hex string %v", err)
		}
		return kb, nil
	}
	return models.Sha256Hash{}, errors.Wrap(ErrUnset, "OCR_KEY_BUNDLE_ID")
}

func (c Config) ORMMaxOpenConns() int {
	return int(c.getWithFallback("ORMMaxOpenConns", parseUint16).(uint16))
}

func (c Config) ORMMaxIdleConns() int {
	return int(c.getWithFallback("ORMMaxIdleConns", parseUint16).(uint16))
}

// OperatorContractAddress represents the address where the Operator.sol
// contract is deployed, this is used for filtering RunLog requests
func (c Config) OperatorContractAddress() common.Address {
	if c.viper.GetString(EnvVarName("OperatorContractAddress")) == "" {
		return common.Address{}
	}
	address, ok := c.getWithFallback("OperatorContractAddress", parseAddress).(*common.Address)
	if !ok {
		return common.Address{}
	}
	return *address
}

// LogLevel represents the maximum level of log messages to output.
func (c Config) LogLevel() LogLevel {
	if c.runtimeStore != nil {
		var value LogLevel
		if err := c.runtimeStore.GetConfigValue("LogLevel", &value); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warnw("Error while trying to fetch LogLevel.", "error", err)
		} else if err == nil {
			return value
		}
	}
	return c.getWithFallback("LogLevel", parseLogLevel).(LogLevel)
}

// SetLogLevel saves a runtime value for the default logger level
func (c Config) SetLogLevel(ctx context.Context, value string) error {
	if c.runtimeStore == nil {
		return errors.New("No runtime store installed")
	}
	var ll LogLevel
	err := ll.Set(value)
	if err != nil {
		return err
	}
	return c.runtimeStore.SetConfigStrValue(ctx, "LogLevel", ll.String())
}

// LogToDisk configures disk preservation of logs.
func (c Config) LogToDisk() bool {
	return c.viper.GetBool(EnvVarName("LogToDisk"))
}

// LogSQLStatements tells chainlink to log all SQL statements made using the default logger
func (c Config) LogSQLStatements() bool {
	if c.runtimeStore != nil {
		logSqlStatements, err := c.runtimeStore.GetConfigBoolValue("LogSQLStatements")
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warnw("Error while trying to fetch LogSQLStatements.", "error", err)
		} else if err == nil {
			return *logSqlStatements
		}
	}
	return c.viper.GetBool(EnvVarName("LogSQLStatements"))
}

// SetLogSQLStatements saves a runtime value for enabling/disabling logging all SQL statements on the default logger
func (c Config) SetLogSQLStatements(ctx context.Context, sqlEnabled bool) error {
	if c.runtimeStore == nil {
		return errors.New("No runtime store installed")
	}

	return c.runtimeStore.SetConfigStrValue(ctx, "LogSQLStatements", strconv.FormatBool(sqlEnabled))
}

// LogSQLMigrations tells chainlink to log all SQL migrations made using the default logger
func (c Config) LogSQLMigrations() bool {
	return c.viper.GetBool(EnvVarName("LogSQLMigrations"))
}

// MinIncomingConfirmations represents the minimum number of block
// confirmations that need to be recorded since a job run started before a task
// can proceed.
func (c Config) MinIncomingConfirmations() uint32 {
	if c.viper.IsSet(EnvVarName("MinIncomingConfirmations")) {
		return c.viper.GetUint32(EnvVarName("MinIncomingConfirmations"))
	}
	return chainSpecificConfig(c).MinIncomingConfirmations
}

// MinRequiredOutgoingConfirmations represents the default minimum number of block
// confirmations that need to be recorded on an outgoing ethtx task before the run can move onto the next task.
// This can be overridden on a per-task basis by setting the `MinRequiredOutgoingConfirmations` parameter.
func (c Config) MinRequiredOutgoingConfirmations() uint64 {
	if c.viper.IsSet(EnvVarName("MinRequiredOutgoingConfirmations")) {
		return c.viper.GetUint64(EnvVarName("MinRequiredOutgoingConfirmations"))
	}
	return chainSpecificConfig(c).MinRequiredOutgoingConfirmations
}

// MinimumContractPayment represents the minimum amount of LINK that must be
// supplied for a contract to be considered.
func (c Config) MinimumContractPayment() *assets.Link {
	return c.getWithFallback("MinimumContractPayment", parseLink).(*assets.Link)
}

// MinimumRequestExpiration is the minimum allowed request expiration for a Service Agreement.
func (c Config) MinimumRequestExpiration() uint64 {
	return c.getWithFallback("MinimumRequestExpiration", parseUint64).(uint64)
}

// P2PListenIP is the ip that libp2p willl bind to and listen on
func (c Config) P2PListenIP() net.IP {
	return c.getWithFallback("P2PListenIP", parseIP).(net.IP)
}

// P2PListenPort is the port that libp2p willl bind to and listen on
func (c Config) P2PListenPort() uint16 {
	return uint16(c.viper.GetUint32(EnvVarName("P2PListenPort")))
}

// P2PAnnounceIP is an optional override. If specified it will force the p2p
// layer to announce this IP as the externally reachable one to the DHT
// If this is set, P2PAnnouncePort MUST also be set.
func (c Config) P2PAnnounceIP() net.IP {
	str := c.viper.GetString(EnvVarName("P2PAnnounceIP"))
	return net.ParseIP(str)
}

// P2PAnnouncePort is an optional override. If specified it will force the p2p
// layer to announce this port as the externally reachable one to the DHT.
// If this is set, P2PAnnounceIP MUST also be set.
func (c Config) P2PAnnouncePort() uint16 {
	return uint16(c.viper.GetUint32(EnvVarName("P2PAnnouncePort")))
}

// P2PDHTAnnouncementCounterUserPrefix can be used to restore the node's
// ability to announce its IP/port on the P2P network after a database
// rollback. Make sure to only increase this value, and *never* decrease it.
// Don't use this variable unless you really know what you're doing, since you
// could semi-permanently exclude your node from the P2P network by
// misconfiguring it.
func (c Config) P2PDHTAnnouncementCounterUserPrefix() uint32 {
	return c.viper.GetUint32(EnvVarName("P2PDHTAnnouncementCounterUserPrefix"))
}

func (c Config) P2PPeerstoreWriteInterval() time.Duration {
	return c.getWithFallback("P2PPeerstoreWriteInterval", parseDuration).(time.Duration)
}

func (c Config) P2PPeerID(override *models.PeerID) (models.PeerID, error) {
	if override != nil {
		return *override, nil
	}
	pidStr := c.viper.GetString(EnvVarName("P2PPeerID"))
	if pidStr != "" {
		var pid models.PeerID
		err := pid.UnmarshalText([]byte(pidStr))
		if err != nil {
			return "", errors.Wrapf(ErrInvalid, "P2P_PEER_ID is invalid %v", err)
		}
		return pid, nil
	}
	return "", errors.Wrap(ErrUnset, "P2P_PEER_ID")
}

func (c Config) P2PPeerIDIsSet() bool {
	return c.viper.GetString(EnvVarName("P2PPeerID")) != ""
}

func (c Config) P2PBootstrapPeers(override []string) ([]string, error) {
	if override != nil {
		return override, nil
	}
	if c.viper.IsSet(EnvVarName("P2PBootstrapPeers")) {
		bps := c.viper.GetStringSlice(EnvVarName("P2PBootstrapPeers"))
		if bps != nil {
			return bps, nil
		}
		return nil, errors.Wrap(ErrUnset, "P2P_BOOTSTRAP_PEERS")
	}
	return []string{}, nil
}

// Port represents the port Chainlink should listen on for client requests.
func (c Config) Port() uint16 {
	return c.getWithFallback("Port", parseUint16).(uint16)
}

func (c Config) HTTPServerWriteTimeout() time.Duration {
	return c.getWithFallback("HTTPServerWriteTimeout", parseDuration).(time.Duration)
}

// ReaperExpiration represents
func (c Config) ReaperExpiration() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("ReaperExpiration", parseDuration).(time.Duration))
}

func (c Config) ReplayFromBlock() int64 {
	return c.viper.GetInt64(EnvVarName("ReplayFromBlock"))
}

// RootDir represents the location on the file system where Chainlink should
// keep its files.
func (c Config) RootDir() string {
	return c.getWithFallback("RootDir", parseHomeDir).(string)
}

// SecureCookies allows toggling of the secure cookies HTTP flag
func (c Config) SecureCookies() bool {
	return c.viper.GetBool(EnvVarName("SecureCookies"))
}

// SessionTimeout is the maximum duration that a user session can persist without any activity.
func (c Config) SessionTimeout() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("SessionTimeout", parseDuration).(time.Duration))
}

// StatsPusherLogging toggles very verbose logging of raw messages for the StatsPusher (also telemetry)
func (c Config) StatsPusherLogging() bool {
	return c.getWithFallback("StatsPusherLogging", parseBool).(bool)
}

// TLSCertPath represents the file system location of the TLS certificate
// Chainlink should use for HTTPS.
func (c Config) TLSCertPath() string {
	return c.viper.GetString(EnvVarName("TLSCertPath"))
}

// TLSHost represents the hostname to use for TLS clients. This should match
// the TLS certificate.
func (c Config) TLSHost() string {
	return c.viper.GetString(EnvVarName("TLSHost"))
}

// TLSKeyPath represents the file system location of the TLS key Chainlink
// should use for HTTPS.
func (c Config) TLSKeyPath() string {
	return c.viper.GetString(EnvVarName("TLSKeyPath"))
}

// TLSPort represents the port Chainlink should listen on for encrypted client requests.
func (c Config) TLSPort() uint16 {
	return c.getWithFallback("TLSPort", parseUint16).(uint16)
}

// TLSRedirect forces TLS redirect for unencrypted connections
func (c Config) TLSRedirect() bool {
	return c.viper.GetBool(EnvVarName("TLSRedirect"))
}

// UnAuthenticatedRateLimit defines the threshold to which requests unauthenticated requests get limited
func (c Config) UnAuthenticatedRateLimit() int64 {
	return c.viper.GetInt64(EnvVarName("UnAuthenticatedRateLimit"))
}

// UnAuthenticatedRateLimitPeriod defines the period to which unauthenticated requests get limited
func (c Config) UnAuthenticatedRateLimitPeriod() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("UnAuthenticatedRateLimitPeriod", parseDuration).(time.Duration))
}

// KeysDir returns the path of the keys directory (used for keystore files).
func (c Config) KeysDir() string {
	return filepath.Join(c.RootDir(), "tempkeys")
}

func (c Config) tlsDir() string {
	return filepath.Join(c.RootDir(), "tls")
}

// KeyFile returns the path where the server key is kept
func (c Config) KeyFile() string {
	if c.TLSKeyPath() == "" {
		return filepath.Join(c.tlsDir(), "server.key")
	}
	return c.TLSKeyPath()
}

// CertFile returns the path where the server certificate is kept
func (c Config) CertFile() string {
	if c.TLSCertPath() == "" {
		return filepath.Join(c.tlsDir(), "server.crt")
	}
	return c.TLSCertPath()
}

// HeadTimeBudget returns the time allowed for context timeout in head tracker
func (c Config) HeadTimeBudget() time.Duration {
	str := c.viper.GetString(EnvVarName("HeadTimeBudget"))
	if str != "" {
		n, err := parseDuration(str)
		if err != nil {
			logger.Errorw(
				"Invalid value provided for HeadTimeBudget, falling back to default.",
				"value", str,
				"error", err)
		} else {
			return n.(time.Duration)
		}
	}
	return chainSpecificConfig(c).HeadTimeBudget
}

// CreateProductionLogger returns a custom logger for the config's root
// directory and LogLevel, with pretty printing for stdout. If LOG_TO_DISK is
// false, the logger will only log to stdout.
func (c Config) CreateProductionLogger() *logger.Logger {
	return logger.CreateProductionLogger(c.RootDir(), c.JSONConsole(), c.LogLevel().Level, c.LogToDisk())
}

// SessionSecret returns a sequence of bytes to be used as a private key for
// session signing or encryption.
func (c Config) SessionSecret() ([]byte, error) {
	return c.SecretGenerator.Generate(c)
}

// SessionOptions returns the sesssions.Options struct used to configure
// the session store.
func (c Config) SessionOptions() sessions.Options {
	return sessions.Options{
		Secure:   c.SecureCookies(),
		HttpOnly: true,
		MaxAge:   86400 * 30,
	}
}

func (c Config) getWithFallback(name string, parser func(string) (interface{}, error)) interface{} {
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
		log.Fatalf(fmt.Sprintf(`Invalid default for %s: "%s"`, name, defaultValue))
	}
	return v
}

// SecretGenerator is the interface for objects that generate a secret
// used to sign or encrypt.
type SecretGenerator interface {
	Generate(Config) ([]byte, error)
}

type filePersistedSecretGenerator struct{}

func (f filePersistedSecretGenerator) Generate(c Config) ([]byte, error) {
	sessionPath := filepath.Join(c.RootDir(), "secret")
	if utils.FileExists(sessionPath) {
		data, err := ioutil.ReadFile(sessionPath)
		if err != nil {
			return data, err
		}
		return base64.StdEncoding.DecodeString(string(data))
	}
	key := securecookie.GenerateRandomKey(32)
	str := base64.StdEncoding.EncodeToString(key)
	err := utils.WriteFileWithMaxPerms(sessionPath, []byte(str), readWritePerms)
	return key, err
}

func parseAddress(str string) (interface{}, error) {
	if str == "" {
		return nil, nil
	} else if common.IsHexAddress(str) {
		val := common.HexToAddress(str)
		return &val, nil
	} else if i, ok := new(big.Int).SetString(str, 10); ok {
		val := common.BigToAddress(i)
		return &val, nil
	}
	return nil, fmt.Errorf("unable to parse '%s' into EIP55-compliant address", str)
}

func parseLink(str string) (interface{}, error) {
	i, ok := new(assets.Link).SetString(str, 10)
	if !ok {
		return i, fmt.Errorf("unable to parse '%v' into *assets.Link(base 10)", str)
	}
	return i, nil
}

func parseLogLevel(str string) (interface{}, error) {
	var lvl LogLevel
	err := lvl.Set(str)
	return lvl, err
}

func parseUint8(s string) (interface{}, error) {
	v, err := strconv.ParseUint(s, 10, 8)
	return uint8(v), err
}

func parseUint16(s string) (interface{}, error) {
	v, err := strconv.ParseUint(s, 10, 16)
	return uint16(v), err
}

func parseUint32(s string) (interface{}, error) {
	v, err := strconv.ParseUint(s, 10, 32)
	return uint32(v), err
}

func parseUint64(s string) (interface{}, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	return v, err
}

func parseURL(s string) (interface{}, error) {
	return url.Parse(s)
}

func parseIP(s string) (interface{}, error) {
	return net.ParseIP(s), nil
}

func parseDuration(s string) (interface{}, error) {
	return time.ParseDuration(s)
}

func parseBool(s string) (interface{}, error) {
	return strconv.ParseBool(s)
}

func parseBigInt(str string) (interface{}, error) {
	i, ok := new(big.Int).SetString(str, 10)
	if !ok {
		return i, fmt.Errorf("unable to parse %v into *big.Int(base 10)", str)
	}
	return i, nil
}

func parseHomeDir(str string) (interface{}, error) {
	exp, err := homedir.Expand(str)
	if err != nil {
		return nil, err
	}
	return filepath.ToSlash(exp), nil
}

// LogLevel determines the verbosity of the events to be logged.
type LogLevel struct {
	zapcore.Level
}

// ForGin keeps Gin's mode at the appropriate level with the LogLevel.
func (ll LogLevel) ForGin() string {
	switch {
	case ll.Level < zapcore.InfoLevel:
		return gin.DebugMode
	default:
		return gin.ReleaseMode
	}
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
