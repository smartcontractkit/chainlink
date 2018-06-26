package store

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"path"
	"reflect"
	"time"

	"github.com/caarlos0/env"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/smartcontractkit/chainlink/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds parameters used by the application which can be overridden
// by setting environment variables.
type Config struct {
	LogLevel                 LogLevel        `env:"LOG_LEVEL" envDefault:"info"`
	RootDir                  string          `env:"ROOT" envDefault:"~/.chainlink"`
	Port                     string          `env:"CHAINLINK_PORT" envDefault:"6688"`
	GuiPort                  string          `env:"GUI_PORT" envDefault:"6689"`
	BasicAuthUsername        string          `env:"USERNAME" envDefault:"chainlink"`
	BasicAuthPassword        string          `env:"PASSWORD" envDefault:"twochains"`
	EthereumURL              string          `env:"ETH_URL" envDefault:"ws://localhost:8546"`
	ChainID                  uint64          `env:"ETH_CHAIN_ID" envDefault:"0"`
	ClientNodeURL            string          `env:"CLIENT_NODE_URL" envDefault:"http://localhost:6688"`
	MinIncomingConfirmations uint64          `env:"MIN_INCOMING_CONFIRMATIONS" envDefault:"0"`
	MinOutgoingConfirmations uint64          `env:"MIN_OUTGOING_CONFIRMATIONS" envDefault:"12"`
	EthGasBumpThreshold      uint64          `env:"ETH_GAS_BUMP_THRESHOLD" envDefault:"12"`
	EthGasBumpWei            big.Int         `env:"ETH_GAS_BUMP_WEI" envDefault:"5000000000"`
	EthGasPriceDefault       big.Int         `env:"ETH_GAS_PRICE_DEFAULT" envDefault:"20000000000"`
	LinkContractAddress      string          `env:"LINK_CONTRACT_ADDRESS" envDefault:"0x514910771AF9Ca656af840dff83E8264EcF986CA"`
	MinimumContractPayment   big.Int         `env:"MINIMUM_CONTRACT_PAYMENT" envDefault:"1000000000000000000"`
	OracleContractAddress    *common.Address `env:"ORACLE_CONTRACT_ADDRESS"`
	DatabasePollInterval     Duration        `env:"DATABASE_POLL_INTERVAL" envDefault:"500ms"`
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
	return config
}

// KeysDir returns the path of the keys directory (used for keystore files).
func (c Config) KeysDir() string {
	return path.Join(c.RootDir, "keys")
}

// CreateProductionLogger returns a custom logger for the config's root directory
// and LogLevel, with pretty printing for stdout.
func (c Config) CreateProductionLogger() *zap.Logger {
	return logger.CreateProductionLogger(c.RootDir, c.LogLevel.Level)
}

func parseEnv(cfg interface{}) error {
	return env.ParseWithFuncs(cfg, env.CustomParsers{
		reflect.TypeOf(&common.Address{}): addressParser,
		reflect.TypeOf(big.Int{}):         bigIntParser,
		reflect.TypeOf(LogLevel{}):        levelParser,
		reflect.TypeOf(Duration{}):        durationParser,
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

func levelParser(str string) (interface{}, error) {
	var lvl LogLevel
	err := lvl.Set(str)
	return lvl, err
}

func durationParser(str string) (interface{}, error) {
	d, err := time.ParseDuration(str)
	return Duration{Duration: d}, err
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
