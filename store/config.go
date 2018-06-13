package store

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"path"
	"reflect"

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
	LogLevel               LogLevel        `env:"LOG_LEVEL" envDefault:"info"`
	RootDir                string          `env:"ROOT" envDefault:"~/.chainlink"`
	Port                   string          `env:"CHAINLINK_PORT" envDefault:"6688"`
	GuiPort                string          `env:"GUI_PORT" envDefault:"6689"`
	BasicAuthUsername      string          `env:"USERNAME" envDefault:"chainlink"`
	BasicAuthPassword      string          `env:"PASSWORD" envDefault:"twochains"`
	EthereumURL            string          `env:"ETH_URL" envDefault:"ws://localhost:8546"`
	ChainID                uint64          `env:"ETH_CHAIN_ID" envDefault:"0"`
	ClientNodeURL          string          `env:"CLIENT_NODE_URL" envDefault:"http://localhost:6688"`
	TxMinConfirmations     uint64          `env:"TX_MIN_CONFIRMATIONS" envDefault:"12"`
	TaskMinConfirmations   uint64          `env:"TASK_MIN_CONFIRMATIONS" envDefault:"0"`
	EthGasBumpThreshold    uint64          `env:"ETH_GAS_BUMP_THRESHOLD" envDefault:"12"`
	EthGasBumpWei          big.Int         `env:"ETH_GAS_BUMP_WEI" envDefault:"5000000000"`
	EthGasPriceDefault     big.Int         `env:"ETH_GAS_PRICE_DEFAULT" envDefault:"20000000000"`
	LinkContractAddress    string          `env:"LINK_CONTRACT_ADDRESS" envDefault:"0x514910771AF9Ca656af840dff83E8264EcF986CA"`
	MinimumContractPayment big.Int         `env:"MINIMUM_CONTRACT_PAYMENT" envDefault:"1000000000000000000"`
	OracleContractAddress  *common.Address `env:"ORACLE_CONTRACT_ADDRESS"`
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

func (c Config) String() string {
	fmtConfig := "LOG_LEVEL: %v\n" +
		"ROOT: %s\n" +
		"CHAINLINK_PORT: %s\n" +
		"GUI_PORT: %s\n" +
		"USERNAME: %s\n" +
		"ETH_URL: %s\n" +
		"ETH_CHAIN_ID: %d\n" +
		"CLIENT_NODE_URL: %s\n" +
		"TX_MIN_CONFIRMATIONS: %d\n" +
		"TASK_MIN_CONFIRMATIONS: %d\n" +
		"ETH_GAS_BUMP_THRESHOLD: %d\n" +
		"ETH_GAS_BUMP_WEI: %s\n" +
		"ETH_GAS_PRICE_DEFAULT: %s\n" +
		"LINK_CONTRACT_ADDRESS: %s\n" +
		"MINIMUM_CONTRACT_PAYMENT: %s\n" +
		"ORACLE_CONTRACT_ADDRESS: %s\n"

	oracleContractAddress := ""
	if c.OracleContractAddress != nil {
		oracleContractAddress = c.OracleContractAddress.String()
	}

	return fmt.Sprintf(
		fmtConfig,
		c.LogLevel,
		c.RootDir,
		c.Port,
		c.GuiPort,
		c.BasicAuthUsername,
		c.EthereumURL,
		c.ChainID,
		c.ClientNodeURL,
		c.TxMinConfirmations,
		c.TaskMinConfirmations,
		c.EthGasBumpThreshold,
		c.EthGasBumpWei.String(),
		c.EthGasPriceDefault.String(),
		c.LinkContractAddress,
		c.MinimumContractPayment.String(),
		oracleContractAddress,
	)
}

func parseEnv(cfg interface{}) error {
	return env.ParseWithFuncs(cfg, env.CustomParsers{
		reflect.TypeOf(&common.Address{}): addressParser,
		reflect.TypeOf(big.Int{}):         bigIntParser,
		reflect.TypeOf(LogLevel{}):        levelParser,
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
