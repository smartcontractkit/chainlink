package store

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
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds parameters used by the application which can be overridden by
// setting environment variables.
//
// If you add an entry here which does not contain sensitive information, you
// should also update presenters.ConfigWhitelist and cmd_test.TestClient_RunNodeShowsEnv.
type Config struct {
	viper           *viper.Viper
	SecretGenerator SecretGenerator
}

// ConfigSchema records the schema of configuration at the type level
type ConfigSchema struct {
	AllowOrigins             string         `env:"ALLOW_ORIGINS" default:"http://localhost:3000,http://localhost:6688"`
	BridgeResponseURL        url.URL        `env:"BRIDGE_RESPONSE_URL"`
	ChainID                  uint64         `env:"ETH_CHAIN_ID" default:"0"`
	ClientNodeURL            string         `env:"CLIENT_NODE_URL" default:"http://localhost:6688"`
	DatabaseTimeout          time.Duration  `env:"DATABASE_TIMEOUT" default:"500ms"`
	DatabaseURL              string         `env:"DATABASE_URL"`
	Dev                      bool           `env:"CHAINLINK_DEV" default:"false"`
	MaximumServiceDuration   time.Duration  `env:"MAXIMUM_SERVICE_DURATION" default:"8760h" `
	MinimumServiceDuration   time.Duration  `env:"MINIMUM_SERVICE_DURATION" default:"0s" `
	EthGasBumpThreshold      uint64         `env:"ETH_GAS_BUMP_THRESHOLD" default:"12" `
	EthGasBumpWei            big.Int        `env:"ETH_GAS_BUMP_WEI" default:"5000000000"`
	EthGasPriceDefault       big.Int        `env:"ETH_GAS_PRICE_DEFAULT" default:"20000000000"`
	EthereumURL              string         `env:"ETH_URL" default:"ws://localhost:8546"`
	JSONConsole              bool           `env:"JSON_CONSOLE" default:"false"`
	LinkContractAddress      string         `env:"LINK_CONTRACT_ADDRESS" default:"0x514910771AF9Ca656af840dff83E8264EcF986CA"`
	ExplorerURL              *url.URL       `env:"EXPLORER_URL"`
	ExplorerAccessKey        string         `env:"EXPLORER_ACCESS_KEY"`
	ExplorerSecret           string         `env:"EXPLORER_SECRET"`
	LogLevel                 LogLevel       `env:"LOG_LEVEL" default:"info"`
	LogToDisk                bool           `env:"LOG_TO_DISK" default:"true"`
	LogSQLStatements         bool           `env:"LOG_SQL" default:"false"`
	MinIncomingConfirmations uint64         `env:"MIN_INCOMING_CONFIRMATIONS" default:"3"`
	MinOutgoingConfirmations uint64         `env:"MIN_OUTGOING_CONFIRMATIONS" default:"12"`
	MinimumContractPayment   assets.Link    `env:"MINIMUM_CONTRACT_PAYMENT" default:"1000000000000000000"`
	MinimumRequestExpiration uint64         `env:"MINIMUM_REQUEST_EXPIRATION" default:"300" `
	OracleContractAddress    common.Address `env:"ORACLE_CONTRACT_ADDRESS"`
	Port                     uint16         `env:"CHAINLINK_PORT" default:"6688"`
	ReaperExpiration         time.Duration  `env:"REAPER_EXPIRATION" default:"240h"`
	RootDir                  string         `env:"ROOT" default:"~/.chainlink"`
	SessionTimeout           time.Duration  `env:"SESSION_TIMEOUT" default:"15m"`
	TLSCertPath              string         `env:"TLS_CERT_PATH" `
	TLSHost                  string         `env:"CHAINLINK_TLS_HOST" `
	TLSKeyPath               string         `env:"TLS_KEY_PATH" `
	TLSPort                  uint16         `env:"CHAINLINK_TLS_PORT" default:"6689"`
}

var configFileNotFoundError = reflect.TypeOf(viper.ConfigFileNotFoundError{})

// NewConfig returns the config with the environment variables set to their
// respective fields, or their defaults if environment variables are not set.
func NewConfig() Config {
	v := viper.New()
	return newConfigWithViper(v)
}

func newConfigWithViper(v *viper.Viper) Config {
	schemaT := reflect.TypeOf(ConfigSchema{})
	for index := 0; index < schemaT.NumField(); index++ {
		item := schemaT.FieldByIndex([]int{index})
		name := item.Tag.Get("env")
		v.SetDefault(name, item.Tag.Get("default"))
		v.BindEnv(name, name)
	}

	config := Config{
		viper:           v,
		SecretGenerator: filePersistedSecretGenerator{},
	}

	if err := os.MkdirAll(config.RootDir(), os.FileMode(0700)); err != nil {
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

// AllowOrigins returns the CORS hosts used by the frontend.
func (c Config) AllowOrigins() string {
	return c.viper.GetString(c.envVarName("AllowOrigins"))
}

// BridgeResponseURL represents the URL for bridges to send a response to.
func (c Config) BridgeResponseURL() *url.URL {
	return c.getWithFallback("BridgeResponseURL", parseURL).(*url.URL)
}

// ChainID represents the chain ID to use for transactions.
func (c Config) ChainID() uint64 {
	return uint64(c.viper.GetInt64(c.envVarName("ChainID")))
}

// ClientNodeURL is the URL of the Ethereum node this Chainlink node should connect to.
func (c Config) ClientNodeURL() string {
	return c.viper.GetString(c.envVarName("ClientNodeURL"))
}

// DatabaseTimeout represents how long to tolerate non response from the DB.
func (c Config) DatabaseTimeout() time.Duration {
	return c.viper.GetDuration(c.envVarName("DatabaseTimeout"))
}

// DatabaseURL configures the URL for chainlink to connect to. This must be
// a properly formatted URL, with a valid scheme (postgres://, file://), or
// an empty string, so the application defaults to .chainlink/db.sqlite.
func (c Config) DatabaseURL() string {
	return c.viper.GetString(c.envVarName("DatabaseURL"))
}

// Dev configures "development" mode for chainlink.
func (c Config) Dev() bool {
	return c.viper.GetBool(c.envVarName("Dev"))
}

// MaximumServiceDuration is the maximum time that a service agreement can run
// from after the time it is created. Default 1 year = 365 * 24h = 8760h
func (c Config) MaximumServiceDuration() time.Duration {
	return c.viper.GetDuration(c.envVarName("MaximumServiceDuration"))
}

// MinimumServiceDuration is the shortest duration from now that a service is
// allowed to run.
func (c Config) MinimumServiceDuration() time.Duration {
	return c.viper.GetDuration(c.envVarName("MinimumServiceDuration"))
}

// EthGasBumpThreshold represents the maximum amount a transaction's ETH amount
// should be increased in order to facilitate a transaction.
func (c Config) EthGasBumpThreshold() uint64 {
	return uint64(c.viper.GetInt64(c.envVarName("EthGasBumpThreshold")))
}

// EthGasBumpWei represents the intervals in which ETH should be increased when
// doing gas bumping.
func (c Config) EthGasBumpWei() *big.Int {
	return c.getWithFallback("EthGasBumpWei", parseBigInt).(*big.Int)
}

// EthGasPriceDefault represents the default gas price for transactions.
func (c Config) EthGasPriceDefault() *big.Int {
	return c.getWithFallback("EthGasPriceDefault", parseBigInt).(*big.Int)
}

// EthereumURL represents the URL of the Ethereum node to connect Chainlink to.
func (c Config) EthereumURL() string {
	return c.viper.GetString(c.envVarName("EthereumURL"))
}

// JSONConsole enables the JSON console.
func (c Config) JSONConsole() bool {
	return c.viper.GetBool(c.envVarName("JSONConsole"))
}

// LinkContractAddress represents the address
func (c Config) LinkContractAddress() string {
	return c.viper.GetString(c.envVarName("LinkContractAddress"))
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
	return c.viper.GetString(c.envVarName("ExplorerAccessKey"))
}

// ExplorerSecret returns the secret for authenticating with explorer
func (c Config) ExplorerSecret() string {
	return c.viper.GetString(c.envVarName("ExplorerSecret"))
}

// OracleContractAddress represents the deployed Oracle contract's address.
func (c Config) OracleContractAddress() *common.Address {
	if c.viper.GetString(c.envVarName("OracleContractAddress")) == "" {
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
	return c.viper.GetBool(c.envVarName("LogToDisk"))
}

// LogSQLStatements tells chainlink to log all SQL statements made using the default logger
func (c Config) LogSQLStatements() bool {
	return c.viper.GetBool(c.envVarName("LogSQLStatements"))
}

// MinIncomingConfirmations represents the minimum number of block
// confirmations that need to be recorded since a job run started before a task
// can proceed.
func (c Config) MinIncomingConfirmations() uint64 {
	return uint64(c.viper.GetInt64(c.envVarName("MinIncomingConfirmations")))
}

// MinOutgoingConfirmations represents the minimum number of block
// confirmations that need to be recorded on an outgoing transaction before a
// task is completed.
func (c Config) MinOutgoingConfirmations() uint64 {
	return uint64(c.viper.GetInt64(c.envVarName("MinOutgoingConfirmations")))
}

// MinimumContractPayment represents the minimum amount of ETH that must be
// supplied for a contract to be considered.
func (c Config) MinimumContractPayment() *assets.Link {
	return c.getWithFallback("MinimumContractPayment", parseLink).(*assets.Link)
}

// MinimumRequestExpiration is the minimum allowed request expiration for a Service Agreement.
func (c Config) MinimumRequestExpiration() uint64 {
	return uint64(c.viper.GetInt64(c.envVarName("MinimumRequestExpiration")))
}

// Port represents the port Chainlink should listen on for client requests.
func (c Config) Port() uint16 {
	return c.getWithFallback("Port", parsePort).(uint16)
}

// ReaperExpiration represents
func (c Config) ReaperExpiration() time.Duration {
	return c.viper.GetDuration(c.envVarName("ReaperExpiration"))
}

// RootDir represents the location on the file system where Chainlink should
// keep its files.
func (c Config) RootDir() string {
	return c.getWithFallback("RootDir", parseHomeDir).(string)
}

// SessionTimeout is the maximum duration that a user session can persist without any activity.
func (c Config) SessionTimeout() time.Duration {
	return c.viper.GetDuration(c.envVarName("SessionTimeout"))
}

// TLSCertPath represents the file system location of the TLS certificate
// Chainlink should use for HTTPS.
func (c Config) TLSCertPath() string {
	return c.viper.GetString(c.envVarName("TLSCertPath"))
}

// TLSHost represents the hostname to use for TLS clients. This should match
// the TLS certificate.
func (c Config) TLSHost() string {
	return c.viper.GetString(c.envVarName("TLSHost"))
}

// TLSKeyPath represents the file system location of the TLS key Chainlink
// should use for HTTPS.
func (c Config) TLSKeyPath() string {
	return c.viper.GetString(c.envVarName("TLSKeyPath"))
}

// TLSPort represents the port Chainlink should listen on for encrypted client requests.
func (c Config) TLSPort() uint16 {
	return c.getWithFallback("TLSPort", parsePort).(uint16)
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
		Secure:   !c.Dev(),
		HttpOnly: true,
		MaxAge:   86400 * 30,
	}
}

func (c Config) envVarName(field string) string {
	schemaT := reflect.TypeOf(ConfigSchema{})
	item, ok := schemaT.FieldByName(field)
	if !ok {
		log.Panicf("Invariant violated, no field of name %s found on ConfigSchema", field)
	}
	return item.Tag.Get("env")
}

func (c Config) defaultValue(name string) (string, bool) {
	schemaT := reflect.TypeOf(ConfigSchema{})
	if item, ok := schemaT.FieldByName(name); ok {
		return item.Tag.Lookup("default")
	}
	log.Panicf("Invariant violated, no field of name %s found for defaultValue", name)
	return "", false
}

func (c Config) zeroValue(name string) interface{} {
	schemaT := reflect.TypeOf(ConfigSchema{})
	if item, ok := schemaT.FieldByName(name); ok {
		if item.Type.Kind() == reflect.Ptr {
			return nil
		}
		return reflect.New(item.Type).Interface()
	}
	log.Panicf("Invariant violated, no field of name %s found for zeroValue", name)
	return nil
}

func (c Config) getWithFallback(name string, parser func(string) (interface{}, error)) interface{} {
	str := c.viper.GetString(c.envVarName(name))
	defaultValue, hasDefault := c.defaultValue(name)
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
		return c.zeroValue(name)
	}

	v, err := parser(defaultValue)
	if err != nil {
		log.Fatalf(fmt.Sprintf(`Invalid default for %s: "%s"`, name, defaultValue))
	}
	return v
}

// NormalizedDatabaseURL returns the DatabaseURL with the empty default
// coerced to a sqlite3 URL.
func (c Config) NormalizedDatabaseURL() string {
	if c.DatabaseURL() == "" {
		return filepath.ToSlash(filepath.Join(c.RootDir(), "db.sqlite3"))
	}
	return c.DatabaseURL()
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
	return key, ioutil.WriteFile(sessionPath, []byte(str), 0644)
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
	return nil, fmt.Errorf("Unable to parse '%s' into EIP55-compliant address", str)
}

func parseLink(str string) (interface{}, error) {
	i, ok := new(assets.Link).SetString(str, 10)
	if !ok {
		return i, fmt.Errorf("Unable to parse '%v' into *assets.Link(base 10)", str)
	}
	return i, nil
}

func parseLogLevel(str string) (interface{}, error) {
	var lvl LogLevel
	err := lvl.Set(str)
	return lvl, err
}

func parsePort(str string) (interface{}, error) {
	d, err := strconv.ParseUint(str, 10, 16)
	return uint16(d), err
}

func parseURL(s string) (interface{}, error) {
	return url.Parse(s)
}

func parseBigInt(str string) (interface{}, error) {
	i, ok := new(big.Int).SetString(str, 10)
	if !ok {
		return i, fmt.Errorf("Unable to parse %v into *big.Int(base 10)", str)
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
