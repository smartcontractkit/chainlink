package store

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"path"
	"reflect"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	homedir "github.com/mitchellh/go-homedir"
	"go.uber.org/zap/zapcore"
)

// Config holds parameters used by the application which can be overridden
// by setting environment variables.
type Config struct {
	LogLevel            LogLevel `env:"LOG_LEVEL" envDefault:"info"`
	RootDir             string   `env:"ROOT" envDefault:"~/.chainlink"`
	Port                string   `env:"PORT" envDefault:"6688"`
	BasicAuthUsername   string   `env:"USERNAME" envDefault:"chainlink"`
	BasicAuthPassword   string   `env:"PASSWORD" envDefault:"twochains"`
	EthereumURL         string   `env:"ETH_URL" envDefault:"ws://localhost:8546"`
	ChainID             uint64   `env:"ETH_CHAIN_ID" envDefault:"0"`
	ClientNodeURL       string   `env:"CLIENT_NODE_URL" envDefault:"http://localhost:6688"`
	EthMinConfirmations uint64   `env:"ETH_MIN_CONFIRMATIONS" envDefault:"12"`
	EthGasBumpThreshold uint64   `env:"ETH_GAS_BUMP_THRESHOLD" envDefault:"12"`
	EthGasBumpWei       big.Int  `env:"ETH_GAS_BUMP_WEI" envDefault:"5000000000"`
	EthGasPriceDefault  big.Int  `env:"ETH_GAS_PRICE_DEFAULT" envDefault:"20000000000"`
	LinkContractAddress string   `env:"LINK_CONTRACT_ADDRESS" envDefault:"0x514910771AF9Ca656af840dff83E8264EcF986CA"`
}

// NewConfig returns the config with the environment variables set to their
// respective fields, or defaults if not present.
func NewConfig() Config {
	config := Config{}
	if err := parseEnv(&config); err != nil {
		log.Fatal(err)
	}
	dir, err := homedir.Expand(config.RootDir)
	if err != nil {
		log.Fatal(err)
	}
	if err = os.MkdirAll(dir, os.FileMode(0700)); err != nil {
		log.Fatal(err)
	}
	config.RootDir = dir
	return config
}

// KeysDir returns the path of the keys directory (used for keystore files).
func (c Config) KeysDir() string {
	return path.Join(c.RootDir, "keys")
}

func parseEnv(cfg interface{}) error {
	return env.ParseWithFuncs(cfg, env.CustomParsers{
		reflect.TypeOf(big.Int{}):  bigIntParser,
		reflect.TypeOf(LogLevel{}): levelParser,
	})
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
