package orm

import (
	"encoding"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/url"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gorilla/securecookie"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Configger is the interface which represents all values that can be enquired from a config store
type Configger interface {
	AllowOrigins() string
	BridgeResponseURL() *url.URL
	CertFile() string
	ChainID() uint64
	ClientNodeURL() string
	DatabaseTimeout() time.Duration
	DatabaseURL() string
	DefaultHTTPLimit() int64
	Dev() bool
	EthGasBumpThreshold() uint64
	EthGasBumpWei() *big.Int
	EthGasPriceDefault() *big.Int
	EthereumURL() *url.URL
	ExplorerAccessKey() string
	ExplorerSecret() string
	ExplorerURL() *url.URL
	JSONConsole() bool
	KeyFile() string
	KeysDir() string
	LinkContractAddress() string
	LogLevel() LogLevel
	LogSQLStatements() bool
	LogToDisk() bool
	MaximumServiceDuration() time.Duration
	MinIncomingConfirmations() uint32
	MinOutgoingConfirmations() uint64
	MinimumContractPayment() *assets.Link
	MinimumRequestExpiration() uint64
	MinimumServiceDuration() time.Duration
	OracleContractAddress() *common.Address
	Port() uint16
	ReaperExpiration() time.Duration
	RootDir() string
	SecureCookies() bool
	SessionOptions() sessions.Options
	SessionSecret() ([]byte, error)
	SessionTimeout() time.Duration
	TLSCertPath() string
	TLSHost() string
	TLSKeyPath() string
	TLSPort() uint16
	TxAttemptLimit() uint16
}

// ConfigStore simply represents a store of configuration key value pairs, it
// carries around the "state" of a configuration but is not meant to be used
// directly
type ConfigStore interface {
	Get(name string, value encoding.TextUnmarshaler) error

	SetMarshaler(name string, value encoding.TextMarshaler) error
	SetString(name, value string) error
	SetStringer(name string, value fmt.Stringer) error
}

// Config represents the implementation of Configger
type Config struct {
	store ConfigStore
}

// NewConfig returns a Config instance
func NewConfig(store ConfigStore) *Config {
	return &Config{store: store}
}

// BridgeResponseURL represents the URL for bridges to send a response to.
func (c Config) BridgeResponseURL() *url.URL {
	var value urlUnmarshaler
	c.store.Get("BridgeResponseURL", &value)
	return (*url.URL)(&value)
}

// ChainID represents the chain ID to use for transactions.
func (c Config) ChainID() uint64 {
	var value uint64Unmarshaler
	c.store.Get("ChainID", &value)
	return uint64(value)
}

// DefaultHTTPLimit defines the limit for HTTP requests.
func (c Config) DefaultHTTPLimit() int64 {
	var value uint64Unmarshaler
	c.store.Get("DefaultHTTPLimit", &value)
	return int64(value)
}

// MaximumServiceDuration is the maximum time that a service agreement can run
// from after the time it is created. Default 1 year = 365 * 24h = 8760h
func (c Config) MaximumServiceDuration() time.Duration {
	var value durationUnmarshaler
	c.store.Get("MaximumServiceDuration", &value)
	return time.Duration(value)
}

// MinimumServiceDuration is the shortest duration from now that a service is
// allowed to run.
func (c Config) MinimumServiceDuration() time.Duration {
	var value durationUnmarshaler
	c.store.Get("MinimumServiceDuration", &value)
	return time.Duration(value)
}

// EthGasBumpThreshold represents the maximum amount a transaction's ETH amount
// should be increased in order to facilitate a transaction.
func (c Config) EthGasBumpThreshold() uint64 {
	var value uint64Unmarshaler
	c.store.Get("EthGasBumpThreshold", &value)
	return uint64(value)
}

// EthGasBumpWei represents the intervals in which ETH should be increased when
// doing gas bumping.
func (c Config) EthGasBumpWei() *big.Int {
	var value models.Big
	c.store.Get("EthGasBumpWei", &value)
	return value.ToInt()
}

// EthGasPriceDefault represents the default gas price for transactions.
func (c Config) EthGasPriceDefault() *big.Int {
	var value models.Big
	c.store.Get("EthGasPriceDefault", &value)
	return value.ToInt()
}

// EthereumURL represents the URL of the Ethereum node to connect Chainlink to.
func (c Config) EthereumURL() *url.URL {
	var value urlUnmarshaler
	c.store.Get("EthereumURL", &value)
	return (*url.URL)(&value)
}

// LinkContractAddress represents the address
func (c Config) LinkContractAddress() string {
	var value stringUnmarshaler
	c.store.Get("LinkContractAddress", &value)
	return string(value)
}

// ExplorerURL returns the websocket URL for this node to push stats to, or nil.
func (c Config) ExplorerURL() *url.URL {
	var value urlUnmarshaler
	c.store.Get("ExplorerURL", &value)
	return (*url.URL)(&value)

	//rval := c.getWithFallback("ExplorerURL", parseURL)
	//switch t := rval.(type) {
	//case nil:
	//return nil
	//case *url.URL:
	//return t
	//default:
	//logger.Panicf("invariant: ExplorerURL returned as type %T", rval)
	//return nil
	//}
}

// ExplorerAccessKey returns the access key for authenticating with explorer
func (c Config) ExplorerAccessKey() string {
	var value stringUnmarshaler
	c.store.Get("ExplorerAccessKey", &value)
	return string(value)
}

// ExplorerSecret returns the secret for authenticating with explorer
func (c Config) ExplorerSecret() string {
	var value stringUnmarshaler
	c.store.Get("ExplorerSecret", &value)
	return string(value)
}

// OracleContractAddress represents the deployed Oracle contract's address.
func (c Config) OracleContractAddress() *common.Address {
	//if c.viper.GetString("OracleContractAddress")) == "" {
	//return nil
	//}
	//return c.getWithFallback("OracleContractAddress", parseAddress).(*common.Address)

	var value addressUnmarshaler
	c.store.Get("OracleContractAddress", &value)
	return (*common.Address)(&value)
}

// MinIncomingConfirmations represents the minimum number of block
// confirmations that need to be recorded since a job run started before a task
// can proceed.
func (c Config) MinIncomingConfirmations() uint32 {
	var value uint64Unmarshaler
	c.store.Get("MinIncomingConfirmations", &value)
	return uint32(value)
}

// MinOutgoingConfirmations represents the minimum number of block
// confirmations that need to be recorded on an outgoing transaction before a
// task is completed.
func (c Config) MinOutgoingConfirmations() uint64 {
	var value uint64Unmarshaler
	c.store.Get("MinOutgoingConfirmations", &value)
	return uint64(value)
}

// MinimumContractPayment represents the minimum amount of ETH that must be
// supplied for a contract to be considered.
func (c Config) MinimumContractPayment() *assets.Link {
	var value models.Big
	c.store.Get("MinimumContractPayment", &value)
	return (*assets.Link)(value.ToInt())
}

// MinimumRequestExpiration is the minimum allowed request expiration for a Service Agreement.
func (c Config) MinimumRequestExpiration() uint64 {
	var value uint64Unmarshaler
	c.store.Get("MinimumRequestExpiration", &value)
	return uint64(value)
}

// ReaperExpiration represents
func (c Config) ReaperExpiration() time.Duration {
	var value durationUnmarshaler
	c.store.Get("ReaperExpiration", &value)
	return time.Duration(value)
}

// SecureCookies allows toggling of the secure cookies HTTP flag
func (c Config) SecureCookies() bool {
	var value boolUnmarshaler
	c.store.Get("SecureCookies", &value)
	return bool(value)
}

// SessionTimeout is the maximum duration that a user session can persist without any activity.
func (c Config) SessionTimeout() time.Duration {
	var value durationUnmarshaler
	c.store.Get("SessionTimeout", &value)
	return time.Duration(value)
}

// TLSCertPath represents the file system location of the TLS certificate
// Chainlink should use for HTTPS.
func (c Config) TLSCertPath() string {
	var value stringUnmarshaler
	c.store.Get("TLSCertPath", &value)
	return string(value)
}

// TLSHost represents the hostname to use for TLS clients. This should match
// the TLS certificate.
func (c Config) TLSHost() string {
	var value stringUnmarshaler
	c.store.Get("TLSHost", &value)
	return string(value)
}

// TLSKeyPath represents the file system location of the TLS key Chainlink
// should use for HTTPS.
func (c Config) TLSKeyPath() string {
	var value stringUnmarshaler
	c.store.Get("TLSKeyPath", &value)
	return string(value)
}

// TxAttemptLimit represents the maximum number of transaction attempts that
// the TxManager should allow to for a transaction
func (c Config) TxAttemptLimit() uint16 {
	var value uint16Unmarshaler
	c.store.Get("TxAttemptLimit", &value)
	return uint16(value)
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

// SessionSecret returns a sequence of bytes to be used as a private key for
// session signing or encryption.
func (c Config) SessionSecret() ([]byte, error) {
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

// SessionOptions returns the sesssions.Options struct used to configure
// the session store.
func (c Config) SessionOptions() sessions.Options {
	return sessions.Options{
		Secure:   c.SecureCookies(),
		HttpOnly: true,
		MaxAge:   86400 * 30,
	}
}

// AllowOrigins returns the CORS hosts used by the frontend.
func (c Config) AllowOrigins() string {
	var value stringUnmarshaler
	c.store.Get("AllowOrigins", &value)
	return string(value)
}

// ClientNodeURL is the URL of the Ethereum node this Chainlink node should connect to.
func (c Config) ClientNodeURL() string {
	var value stringUnmarshaler
	c.store.Get("ClientNodeURL", &value)
	return string(value)
}

// DatabaseTimeout represents how long to tolerate non response from the DB.
func (c Config) DatabaseTimeout() time.Duration {
	var value durationUnmarshaler
	c.store.Get("DatabaseTimeout", &value)
	return time.Duration(value)
}

// DatabaseURL configures the URL for chainlink to connect to. This must be
// a properly formatted URL, with a valid scheme (postgres://, file://), or
// an empty string, so the application defaults to .chainlink/db.sqlite.
func (c Config) DatabaseURL() string {
	var value stringUnmarshaler
	c.store.Get("DatabaseURL", &value)
	url := string(value)
	if url == "" {
		return filepath.ToSlash(filepath.Join(c.RootDir(), "db.sqlite3"))
	}
	return url
}

// Dev configures "development" mode for chainlink.
func (c Config) Dev() bool {
	var value boolUnmarshaler
	c.store.Get("Dev", &value)
	return bool(value)
}

// JSONConsole enables the JSON console.
func (c Config) JSONConsole() bool {
	var value boolUnmarshaler
	c.store.Get("JSONConsole", &value)
	return bool(value)
}

// KeysDir returns the path of the keys directory (used for keystore files).
func (c Config) KeysDir() string {
	return filepath.Join(c.RootDir(), "tempkeys")
}

// LogLevel represents the maximum level of log messages to output.
func (c Config) LogLevel() LogLevel {
	var value logLevelUnmarshaler
	c.store.Get("LogLevel", &value)
	return LogLevel(value)
}

// LogSQLStatements tells chainlink to log all SQL statements made using the default logger
func (c Config) LogSQLStatements() bool {
	var value boolUnmarshaler
	c.store.Get("LogSQLStatements", &value)
	return bool(value)
}

// LogToDisk configures disk preservation of logs.
func (c Config) LogToDisk() bool {
	var value boolUnmarshaler
	c.store.Get("LogToDisk", &value)
	return bool(value)
}

// Port represents the port Chainlink should listen on for client requests.
func (c Config) Port() uint16 {
	var value uint16Unmarshaler
	c.store.Get("Port", &value)
	return uint16(value)
}

// RootDir represents the location on the file system where Chainlink should
// keep its files.
func (c Config) RootDir() string {
	var value stringUnmarshaler
	c.store.Get("RootDir", &value)
	return string(value)
}

// TLSPort represents the port Chainlink should listen on for encrypted client requests.
func (c Config) TLSPort() uint16 {
	var value uint16Unmarshaler
	c.store.Get("TLSPort", &value)
	return uint16(value)
}
