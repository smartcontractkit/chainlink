package store

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"path"
	"reflect"

	"github.com/gin-gonic/gin"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/smartcontractkit/env"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	LogLevel            LogLevel `env:"LOG_LEVEL" envDefault:"info"`
	RootDir             string   `env:"ROOT" envDefault:"~/.chainlink"`
	BasicAuthUsername   string   `env:"USERNAME" envDefault:"chainlink"`
	BasicAuthPassword   string   `env:"PASSWORD" envDefault:"twochains"`
	EthereumURL         string   `env:"ETH_URL" envDefault:"ws://localhost:8545"`
	ChainID             uint64   `env:"ETH_CHAIN_ID" envDefault:"0"`
	PollingSchedule     string   `env:"POLLING_SCHEDULE" envDefault:"*/15 * * * * *"`
	ClientNodeURL       string   `env:"CLIENT_NODE_URL" envDefault:"http://localhost:8080"`
	EthMinConfirmations uint64   `env:"ETH_MIN_CONFIRMATIONS" envDefault:"12"`
	EthGasBumpThreshold uint64   `env:"ETH_GAS_BUMP_THRESHOLD" envDefault:"12"`
	EthGasBumpWei       big.Int  `env:"ETH_GAS_BUMP_WEI" envDefault:"5000000000"`
	EthGasPriceDefault  big.Int  `env:"ETH_GAS_PRICE_DEFAULT" envDefault:"20000000000"`
}

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

type LogLevel struct {
	zapcore.Level
}

func (ll LogLevel) ForGin() string {
	switch {
	case ll.Level < zapcore.InfoLevel:
		return gin.DebugMode
	default:
		return gin.ReleaseMode
	}
}
