package store

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/url"
	"os"
	"path"
	"reflect"
	"strconv"
	"time"

	"github.com/caarlos0/env"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/mitchellh/go-homedir"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds parameters used by the application which can be overridden
// by setting environment variables.
type Config struct {
	AllowOrigins      string        `env:"ALLOW_ORIGINS" envDefault:"http://localhost:3000,http://localhost:6688"`
	BridgeResponseURL models.WebURL `env:"BRIDGE_RESPONSE_URL" envDefault:""`
	ChainID           uint64        `env:"ETH_CHAIN_ID" envDefault:"0"`
	ClientNodeURL     string        `env:"CLIENT_NODE_URL" envDefault:"http://localhost:6688"`
	DatabaseTimeout   Duration      `env:"DATABASE_TIMEOUT" envDefault:"500ms"`
	Dev               bool          `env:"CHAINLINK_DEV" envDefault:"false"`
	// How long from now that a service agreement is allowed to run. Default 1 year = 365 * 24h = 8760h
	MaximumServiceDuration Duration `env:"MAXIMUM_SERVICE_DURATION" envDefault:"8760h"`
	// Shortest duration from now that a service is allowed to run.
	MinimumServiceDuration   Duration        `env:"MINIMUM_SERVICE_DURATION" envDefault:"0s"`
	EthGasBumpThreshold      uint64          `env:"ETH_GAS_BUMP_THRESHOLD" envDefault:"12"`
	EthGasBumpWei            big.Int         `env:"ETH_GAS_BUMP_WEI" envDefault:"5000000000"`
	EthGasPriceDefault       big.Int         `env:"ETH_GAS_PRICE_DEFAULT" envDefault:"20000000000"`
	EthereumURL              string          `env:"ETH_URL" envDefault:"ws://localhost:8546"`
	JSONStdout               bool            `env:"JSON_STDOUT" envDefault:"false"`
	LinkContractAddress      string          `env:"LINK_CONTRACT_ADDRESS" envDefault:"0x514910771AF9Ca656af840dff83E8264EcF986CA"`
	LogLevel                 LogLevel        `env:"LOG_LEVEL" envDefault:"info"`
	MinIncomingConfirmations uint64          `env:"MIN_INCOMING_CONFIRMATIONS" envDefault:"0"`
	MinOutgoingConfirmations uint64          `env:"MIN_OUTGOING_CONFIRMATIONS" envDefault:"12"`
	MinimumContractPayment   assets.Link     `env:"MINIMUM_CONTRACT_PAYMENT" envDefault:"1000000000000000000"`
	MinimumRequestExpiration uint64          `env:"MINIMUM_REQUEST_EXPIRATION" envDefault:"300"`
	OracleContractAddress    *common.Address `env:"ORACLE_CONTRACT_ADDRESS"`
	Port                     uint16          `env:"CHAINLINK_PORT" envDefault:"6688"`
	ReaperExpiration         Duration        `env:"REAPER_EXPIRATION" envDefault:"240h"`
	RootDir                  string          `env:"ROOT" envDefault:"~/.chainlink"`
	SessionTimeout           Duration        `env:"SESSION_TIMEOUT" envDefault:"15m"`
	TLSCertPath              string          `env:"TLS_CERT_PATH" envDefault:""`
	TLSHost                  string          `env:"CHAINLINK_TLS_HOST" envDefault:""`
	TLSKeyPath               string          `env:"TLS_KEY_PATH" envDefault:""`
	TLSPort                  uint16          `env:"CHAINLINK_TLS_PORT" envDefault:"6689"`
	SecretGenerator          SecretGenerator
}

// NewConfig returns the config with the environment variables set to their
// respective fields, or defaults if not present.
func NewConfig() Config {
	config := Config{}
	if err := parseEnv(&config); err != nil {
		log.Fatal(fmt.Errorf("error parsing environment: %+v", err))
	}
	dir, err := homedir.Expand(config.RootDir)
	if err != nil {
		log.Fatal(fmt.Errorf("error expanding $HOME: %+v", err))
	}
	if err = os.MkdirAll(dir, os.FileMode(0700)); err != nil {
		log.Fatal(fmt.Errorf("error creating %s: %+v", dir, err))
	}
	config.RootDir = dir
	config.SecretGenerator = filePersistedSecretGenerator{}
	return config
}

// KeysDir returns the path of the keys directory (used for keystore files).
func (c Config) KeysDir() string {
	return path.Join(c.RootDir, "keys")
}

func (c Config) tlsDir() string {
	return path.Join(c.RootDir, "tls")
}

// KeyFile returns the path where the server key is kept
func (c Config) KeyFile() string {
	if c.TLSKeyPath == "" {
		return path.Join(c.tlsDir(), "server.key")
	}
	return c.TLSKeyPath
}

// CertFile returns the path where the server certificate is kept
func (c Config) CertFile() string {
	if c.TLSCertPath == "" {
		return path.Join(c.tlsDir(), "server.crt")
	}
	return c.TLSCertPath
}

// CreateProductionLogger returns a custom logger for the config's root directory
// and LogLevel, with pretty printing for stdout.
func (c Config) CreateProductionLogger() *zap.Logger {
	return logger.CreateProductionLogger(c.RootDir, c.JSONStdout, c.LogLevel.Level)
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
		Secure:   c.Dev == false,
		HttpOnly: true,
		MaxAge:   86400 * 30,
	}
}

// SecretGenerator is the interface for objects that generate a secret
// used to sign or encrypt.
type SecretGenerator interface {
	Generate(Config) ([]byte, error)
}

type filePersistedSecretGenerator struct{}

func (f filePersistedSecretGenerator) Generate(c Config) ([]byte, error) {
	sessionPath := path.Join(c.RootDir, "secret")
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

func parseEnv(cfg interface{}) error {
	return env.ParseWithFuncs(cfg, env.CustomParsers{
		reflect.TypeOf(&common.Address{}): addressParser,
		reflect.TypeOf(big.Int{}):         bigIntParser,
		reflect.TypeOf(assets.Link{}):     linkParser,
		reflect.TypeOf(LogLevel{}):        levelParser,
		reflect.TypeOf(Duration{}):        durationParser,
		reflect.TypeOf(models.WebURL{}):   urlParser,
		reflect.TypeOf(uint16(0)):         portParser,
	})
}

func addressParser(str string) (interface{}, error) {
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

func bigIntParser(str string) (interface{}, error) {
	i, ok := new(big.Int).SetString(str, 10)
	if !ok {
		return i, fmt.Errorf("Unable to parse %v into *big.Int(base 10)", str)
	}
	return *i, nil
}

func linkParser(str string) (interface{}, error) {
	i, ok := new(assets.Link).SetString(str, 10)
	if !ok {
		return i, fmt.Errorf("Unable to parse %v into *assets.Link(base 10)", str)
	}
	return *i, nil
}

func levelParser(str string) (interface{}, error) {
	var lvl LogLevel
	err := lvl.Set(str)
	return lvl, err
}

func durationParser(str string) (interface{}, error) {
	d, err := time.ParseDuration(str)
	return Duration{Duration: d}, err
}

func portParser(str string) (interface{}, error) {
	d, err := strconv.ParseUint(str, 10, 16)
	return uint16(d), err
}

func urlParser(s string) (interface{}, error) {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return nil, err
	}
	return models.WebURL(*u), nil
}

// LogLevel determines the verbosity of the events to be logged.
type LogLevel struct {
	zapcore.Level
}

// Duration returns a time duration with the supported
// units of "ns", "us", "ms", "s", "m", "h".
type Duration struct {
	time.Duration
}

// MarshalText returns the byte slice of the formatted duration e.g. "500ms"
func (d Duration) MarshalText() ([]byte, error) {
	b := []byte(d.Duration.String())
	return b, nil
}

// UnmarshalText parses the time.Duration and assigns it
func (d *Duration) UnmarshalText(text []byte) error {
	td, err := time.ParseDuration((string)(text))
	d.Duration = td
	return err
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
