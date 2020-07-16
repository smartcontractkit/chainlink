package orm

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// this permission grants read / write accccess to file owners only
const readWritePerms = os.FileMode(0600)

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
		v.BindEnv(name, name)
	}

	config := &Config{
		viper:           v,
		SecretGenerator: filePersistedSecretGenerator{},
	}

	if err := utils.EnsureDirAndPerms(config.RootDir(), os.FileMode(0700)); err != nil {
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
		logger.Warnf(
			"ETH_GAS_BUMP_PERCENT of %v is less than Geth's default of %v, transactions may fail with underpriced replacement errors",
			c.EthGasBumpPercent(),
			ethCore.DefaultTxPoolConfig.PriceBump,
		)
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
	} else {
		return c.AdvisoryLockID
	}
}

func (c Config) GetDatabaseDialectConfiguredOrDefault() DialectName {
	if c.Dialect == "" {
		return DialectPostgres
	} else {
		return c.Dialect
	}
}

// AllowOrigins returns the CORS hosts used by the frontend.
func (c Config) AllowOrigins() string {
	return c.viper.GetString(EnvVarName("AllowOrigins"))
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

// DefaultHTTPLimit defines the limit for HTTP requests.
func (c Config) DefaultHTTPLimit() int64 {
	return c.viper.GetInt64(EnvVarName("DefaultHTTPLimit"))
}

// DefaultHTTPTimeout defines the default timeout for http requests
func (c Config) DefaultHTTPTimeout() models.Duration {
	return c.getDuration("DefaultHTTPTimeout")
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

// MaxRPCCallsPerSecond returns the rate at which RPC calls can be fired
func (c Config) MaxRPCCallsPerSecond() uint64 {
	return c.viper.GetUint64(EnvVarName("MaxRPCCallsPerSecond"))
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

// EthGasBumpThreshold is the number of blocks to wait for confirmations before bumping gas again
func (c Config) EthGasBumpThreshold() uint64 {
	return c.viper.GetUint64(EnvVarName("EthGasBumpThreshold"))
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

// EthGasLimitDefault  sets the default gas limit for outgoing transactions.
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

// EthereumURL represents the URL of the Ethereum node to connect Chainlink to.
func (c Config) EthereumURL() string {
	return c.viper.GetString(EnvVarName("EthereumURL"))
}

// EthereumDisabled shows whether Ethereum interactions are supported.
func (c Config) EthereumDisabled() bool {
	return c.viper.GetBool(EnvVarName("EthereumDisabled"))
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

// OracleContractAddress represents the deployed Oracle contract's address.
func (c Config) OracleContractAddress() *common.Address {
	if c.viper.GetString(EnvVarName("OracleContractAddress")) == "" {
		return nil
	}
	return c.getWithFallback("OracleContractAddress", parseAddress).(*common.Address)
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

// MinOutgoingConfirmations represents the minimum number of block
// confirmations that need to be recorded on an outgoing transaction before a
// task is completed.
func (c Config) MinOutgoingConfirmations() uint64 {
	return c.viper.GetUint64(EnvVarName("MinOutgoingConfirmations"))
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
func (c Config) CreateProductionLogger() *zap.Logger {
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
	err := utils.WriteFileWithPerms(sessionPath, []byte(str), readWritePerms)
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

func parseURL(s string) (interface{}, error) {
	return url.Parse(s)
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
