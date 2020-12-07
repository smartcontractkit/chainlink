package orm

import (
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

	"github.com/multiformats/go-multiaddr"

	"github.com/libp2p/go-libp2p-core/peer"

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
)

// Config holds parameters used by the application which can be overridden by
// setting environment variables.
//
// If you add an entry here which does not contain sensitive information, you
// should also update presenters.ConfigWhitelist and cmd_test.TestClient_RunNodeShowsEnv.
type Config struct {
	viper           *viper.Viper
	SecretGenerator SecretGenerator
	runtimeStore    *ORM
	Dialect         DialectName
	AdvisoryLockID  int64
}

var configFileNotFoundError = reflect.TypeOf(viper.ConfigFileNotFoundError{})

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
		v.SetDefault(name, item.Tag.Get("default"))
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
	if c.EthGasBumpWei().Cmp(big.NewInt(5000000000)) < 0 {
		return errors.Errorf("ETH_GAS_BUMP_WEI of %s Wei may not be less than the minimum allowed value of 5 GWei", c.EthGasBumpWei().String())
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
	}
	if err := ocr.SanityCheckLocalConfig(lc); err != nil {
		return err
	}
	if _, err := c.P2PPeerID(""); errors.Cause(err) == ErrInvalid {
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

func (c Config) GetDatabaseDialectConfiguredOrDefault() DialectName {
	if c.Dialect == "" {
		return DialectPostgres
	}
	return c.Dialect
}

// AllowOrigins returns the CORS hosts used by the frontend.
func (c Config) AllowOrigins() string {
	return c.viper.GetString(EnvVarName("AllowOrigins"))
}

// BalanceMonitorEnabled enables the balance monitor
func (c Config) BalanceMonitorEnabled() bool {
	return c.viper.GetBool(EnvVarName("BalanceMonitorEnabled"))
}

// BlockBackfillDepth specifies the number of blocks before the current HEAD that the
// log broadcaster will try to re-consume logs from
func (c Config) BlockBackfillDepth() uint64 {
	return c.viper.GetUint64(EnvVarName("BlockBackfillDepth"))
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

func (c Config) getDuration(s string) models.Duration {
	rv, err := models.MakeDuration(c.viper.GetDuration(EnvVarName(s)))
	if err != nil {
		panic(errors.Wrapf(err, "bad duration for config value %s: %s", s, rv))
	}
	return rv
}

func (c Config) DatabaseListenerMinReconnectInterval() time.Duration {
	return c.viper.GetDuration(EnvVarName("DatabaseListenerMinReconnectInterval"))
}

func (c Config) DatabaseListenerMaxReconnectDuration() time.Duration {
	return c.viper.GetDuration(EnvVarName("DatabaseListenerMaxReconnectDuration"))
}

func (c Config) DatabaseMaximumTxDuration() time.Duration {
	return c.viper.GetDuration(EnvVarName("DatabaseMaximumTxDuration"))
}

// DatabaseTimeout represents how long to tolerate non response from the DB.
func (c Config) DatabaseTimeout() models.Duration {
	return c.getDuration("DatabaseTimeout")
}

// DatabaseURL configures the URL for chainlink to connect to. This must be
// a properly formatted URL, with a valid scheme (postgres://)
func (c Config) DatabaseURL() string {
	return c.viper.GetString(EnvVarName("DatabaseURL"))
}

// MigrateDatabase determines whether the database will be automatically
// migrated on application startup if set to true
func (c Config) MigrateDatabase() bool {
	return c.viper.GetBool(EnvVarName("MigrateDatabase"))
}

// DefaultMaxHTTPAttempts defines the limit for HTTP requests.
func (c Config) DefaultMaxHTTPAttempts() uint {
	return c.viper.GetUint(EnvVarName("DefaultMaxHTTPAttempts"))
}

// DefaultHTTPLimit defines the size limit for HTTP requests and responses
func (c Config) DefaultHTTPLimit() int64 {
	return c.viper.GetInt64(EnvVarName("DefaultHTTPLimit"))
}

// DefaultHTTPTimeout defines the default timeout for http requests
func (c Config) DefaultHTTPTimeout() models.Duration {
	return c.getDuration("DefaultHTTPTimeout")
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

// FeatureOffchainReporting enables the Flux Monitor feature.
func (c Config) FeatureOffchainReporting() bool {
	return c.viper.GetBool(EnvVarName("FeatureOffchainReporting"))
}

// MaximumServiceDuration is the maximum time that a service agreement can run
// from after the time it is created. Default 1 year = 365 * 24h = 8760h
func (c Config) MaximumServiceDuration() models.Duration {
	return c.getDuration("MaximumServiceDuration")
}

// MinimumServiceDuration is the shortest duration from now that a service is
// allowed to run.
func (c Config) MinimumServiceDuration() models.Duration {
	return c.getDuration("MinimumServiceDuration")
}

// EthBalanceMonitorBlockDelay is the number of blocks that the balance monitor
// trails behind head. This is required e.g. for Infura because they will often
// announce a new head, then route a request to a different node which does not
// have this head yet.
func (c Config) EthBalanceMonitorBlockDelay() uint16 {
	return c.getWithFallback("EthBalanceMonitorBlockDelay", parseUint16).(uint16)
}

// EthGasBumpThreshold is the number of blocks to wait for confirmations before bumping gas again
func (c Config) EthGasBumpThreshold() uint64 {
	return c.viper.GetUint64(EnvVarName("EthGasBumpThreshold"))
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
	return c.getWithFallback("EthGasBumpWei", parseBigInt).(*big.Int)
}

// EthMaxGasPriceWei is the maximum amount in Wei that a transaction will be
// bumped to before abandoning it and marking it as errored.
func (c Config) EthMaxGasPriceWei() *big.Int {
	return c.getWithFallback("EthMaxGasPriceWei", parseBigInt).(*big.Int)
}

// EthGasLimitDefault sets the default gas limit for outgoing transactions.
func (c Config) EthGasLimitDefault() uint64 {
	return c.viper.GetUint64(EnvVarName("EthGasLimitDefault"))
}

// EthGasPriceDefault is the starting gas price for every transaction
func (c Config) EthGasPriceDefault() *big.Int {
	if c.runtimeStore != nil {
		var value big.Int
		if err := c.runtimeStore.GetConfigValue("EthGasPriceDefault", &value); err != nil && errors.Cause(err) != ErrorNotFound {
			logger.Warnw("Error while trying to fetch EthGasPriceDefault.", "error", err)
		} else if err == nil {
			return &value
		}
	}
	return c.getWithFallback("EthGasPriceDefault", parseBigInt).(*big.Int)
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
	return c.viper.GetUint(EnvVarName("EthFinalityDepth"))
}

// EthHeadTrackerHistoryDepth is the number of heads to keep in the `heads` database table.
// This number should be at least as large as `EthFinalityDepth`.
// There may be a small performance penalty to setting this to something very large (10,000+)
func (c Config) EthHeadTrackerHistoryDepth() uint {
	return c.viper.GetUint(EnvVarName("EthHeadTrackerHistoryDepth"))
}

// EthHeadTrackerMaxBufferSize is the maximum number of heads that may be
// buffered in front of the head tracker before older heads start to be
// dropped. You may think of it as something like the maximum permittable "lag"
// for the head tracker before we start dropping heads to keep up.
func (c Config) EthHeadTrackerMaxBufferSize() uint {
	return c.viper.GetUint(EnvVarName("EthHeadTrackerMaxBufferSize"))
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
	return c.getWithFallback("GasUpdaterBlockDelay", parseUint16).(uint16)
}

// GasUpdaterBlockHistorySize is the number of past blocks to keep in memory to
// use as a basis for calculating a percentile gas price
func (c Config) GasUpdaterBlockHistorySize() uint16 {
	return c.getWithFallback("GasUpdaterBlockHistorySize", parseUint16).(uint16)
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
	return c.viper.GetDuration(EnvVarName("TriggerFallbackDBPollInterval"))
}

func (c Config) JobPipelineMaxTaskDuration() time.Duration {
	return c.viper.GetDuration(EnvVarName("JobPipelineMaxTaskDuration"))
}

// JobPipelineParallelism controls how many workers the pipeline.Runner
// uses in parallel
func (c Config) JobPipelineParallelism() uint8 {
	return c.getWithFallback("JobPipelineParallelism", parseUint8).(uint8)
}

func (c Config) JobPipelineReaperInterval() time.Duration {
	return c.viper.GetDuration(EnvVarName("JobPipelineReaperInterval"))
}

func (c Config) JobPipelineReaperThreshold() time.Duration {
	return c.viper.GetDuration(EnvVarName("JobPipelineReaperThreshold"))
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
	return c.viper.GetDuration(EnvVarName("OCRBootstrapCheckInterval"))
}

func (c Config) OCRContractTransmitterTransmitTimeout() time.Duration {
	return c.viper.GetDuration(EnvVarName("OCRContractTransmitterTransmitTimeout"))
}

func (c Config) getDurationWithOverride(override time.Duration, field string) time.Duration {
	if override != time.Duration(0) {
		return override
	}
	return c.viper.GetDuration(EnvVarName(field))
}

func (c Config) OCRObservationTimeout(override time.Duration) time.Duration {
	return c.getDurationWithOverride(override, "OCRObservationTimeout")
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
	return c.viper.GetDuration(EnvVarName("OCRDatabaseTimeout"))
}

func (c Config) OCRDHTLookupInterval() int {
	return c.viper.GetInt(EnvVarName("OCRDHTLookupInterval"))
}

func (c Config) OCRIncomingMessageBufferSize() int {
	return c.viper.GetInt(EnvVarName("OCRIncomingMessageBufferSize"))
}

func (c Config) OCRNewStreamTimeout() time.Duration {
	return c.viper.GetDuration(EnvVarName("OCRNewStreamTimeout"))
}

func (c Config) OCROutgoingMessageBufferSize() int {
	return c.viper.GetInt(EnvVarName("OCROutgoingMessageBufferSize"))
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
	return c.getWithFallback("LogLevel", parseLogLevel).(LogLevel)
}

// LogToDisk configures disk preservation of logs.
func (c Config) LogToDisk() bool {
	return c.viper.GetBool(EnvVarName("LogToDisk"))
}

// LogSQLStatements tells chainlink to log all SQL statements made using the default logger
func (c Config) LogSQLStatements() bool {
	return c.viper.GetBool(EnvVarName("LogSQLStatements"))
}

// LogSQLMigrations tells chainlink to log all SQL migrations made using the default logger
func (c Config) LogSQLMigrations() bool {
	return c.viper.GetBool(EnvVarName("LogSQLMigrations"))
}

// MinIncomingConfirmations represents the minimum number of block
// confirmations that need to be recorded since a job run started before a task
// can proceed.
func (c Config) MinIncomingConfirmations() uint32 {
	return c.viper.GetUint32(EnvVarName("MinIncomingConfirmations"))
}

// MinRequiredOutgoingConfirmations represents the default minimum number of block
// confirmations that need to be recorded on an outgoing ethtx task before the run can move onto the next task.
// This can be overridden on a per-task basis by setting the `MinRequiredOutgoingConfirmations` parameter.
func (c Config) MinRequiredOutgoingConfirmations() uint64 {
	return c.viper.GetUint64(EnvVarName("MinRequiredOutgoingConfirmations"))
}

// MinimumContractPayment represents the minimum amount of LINK that must be
// supplied for a contract to be considered.
func (c Config) MinimumContractPayment() *assets.Link {
	return c.getWithFallback("MinimumContractPayment", parseLink).(*assets.Link)
}

// MinimumRequestExpiration is the minimum allowed request expiration for a Service Agreement.
func (c Config) MinimumRequestExpiration() uint64 {
	return c.viper.GetUint64(EnvVarName("MinimumRequestExpiration"))
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

func (c Config) P2PPeerstoreWriteInterval() time.Duration {
	return c.viper.GetDuration(EnvVarName("P2PPeerstoreWriteInterval"))
}

func (c Config) P2PPeerID(override models.PeerID) (models.PeerID, error) {
	if override != "" {
		return override, nil
	}
	pidStr := c.viper.GetString(EnvVarName("P2PPeerID"))
	if pidStr != "" {
		pid, err := peer.Decode(pidStr)
		if err != nil {
			return "", errors.Wrapf(ErrInvalid, "P2P_PEER_ID is invalid %v", err)
		}
		return models.PeerID(pid), nil
	}
	return "", errors.Wrap(ErrUnset, "P2P_PEER_ID")
}

func (c Config) P2PBootstrapPeers(override []string) ([]string, error) {
	if override != nil {
		return override, nil
	}
	bps := c.viper.GetStringSlice(EnvVarName("P2PBootstrapPeers"))
	if bps != nil {
		return bps, nil
	}
	return nil, errors.Wrap(ErrUnset, "P2P_BOOTSTRAP_PEERS")
}

// Port represents the port Chainlink should listen on for client requests.
func (c Config) Port() uint16 {
	return c.getWithFallback("Port", parseUint16).(uint16)
}

// ReaperExpiration represents
func (c Config) ReaperExpiration() models.Duration {
	return c.getDuration("ReaperExpiration")
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
	return c.getDuration("SessionTimeout")
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

// TxAttemptLimit is the maximum number of transaction attempts (gas bumps)
// that will occur before giving a transaction up as errored
// NOTE: That initial transactions are retried forever until they succeed
func (c Config) TxAttemptLimit() uint16 {
	return c.getWithFallback("TxAttemptLimit", parseUint16).(uint16)
}

// TLSRedirect forces TLS redirect for unencrypted connections
func (c Config) TLSRedirect() bool {
	return c.viper.GetBool(EnvVarName("TLSRedirect"))
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

// CreateProductionLogger returns a custom logger for the config's root
// directory and LogLevel, with pretty printing for stdout. If LOG_TO_DISK is
// false, the logger will only log to stdout.
func (c Config) CreateProductionLogger() *logger.Logger {
	return logger.CreateProductionLogger(
		c.RootDir(), c.JSONConsole(), c.LogLevel().Level, c.LogToDisk())
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

func parseUint16(str string) (interface{}, error) {
	d, err := strconv.ParseUint(str, 10, 16)
	return uint16(d), err
}

func parseUint8(str string) (interface{}, error) {
	d, err := strconv.ParseUint(str, 10, 8)
	return uint8(d), err
}

func parseURL(s string) (interface{}, error) {
	return url.Parse(s)
}

func parseIP(s string) (interface{}, error) {
	return net.ParseIP(s), nil
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
