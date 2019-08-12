package orm

import (
	"log"
	"math/big"
	"net/url"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/store/assets"
)

// ConfigSchema records the schema of configuration at the type level
type ConfigSchema struct {
	AllowOrigins             string         `env:"ALLOW_ORIGINS" default:"http://localhost:3000,http://localhost:6688"`
	BridgeResponseURL        url.URL        `env:"BRIDGE_RESPONSE_URL"`
	ChainID                  uint64         `env:"ETH_CHAIN_ID" default:"0"`
	ClientNodeURL            string         `env:"CLIENT_NODE_URL" default:"http://localhost:6688"`
	DatabaseTimeout          time.Duration  `env:"DATABASE_TIMEOUT" default:"500ms"`
	DatabaseURL              string         `env:"DATABASE_URL"`
	DefaultHTTPLimit         int64          `env:"DEFAULT_HTTP_LIMIT" default:"32768"`
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
	MinIncomingConfirmations uint32         `env:"MIN_INCOMING_CONFIRMATIONS" default:"3"`
	MinOutgoingConfirmations uint64         `env:"MIN_OUTGOING_CONFIRMATIONS" default:"12"`
	MinimumContractPayment   assets.Link    `env:"MINIMUM_CONTRACT_PAYMENT" default:"1000000000000000000"`
	MinimumRequestExpiration uint64         `env:"MINIMUM_REQUEST_EXPIRATION" default:"300" `
	OracleContractAddress    common.Address `env:"ORACLE_CONTRACT_ADDRESS"`
	Port                     uint16         `env:"CHAINLINK_PORT" default:"6688"`
	ReaperExpiration         time.Duration  `env:"REAPER_EXPIRATION" default:"240h"`
	RootDir                  string         `env:"ROOT" default:"~/.chainlink"`
	SecureCookies            bool           `env:"SECURE_COOKIES" default:"true"`
	SessionTimeout           time.Duration  `env:"SESSION_TIMEOUT" default:"15m"`
	TLSCertPath              string         `env:"TLS_CERT_PATH" `
	TLSHost                  string         `env:"CHAINLINK_TLS_HOST" `
	TLSKeyPath               string         `env:"TLS_KEY_PATH" `
	TLSPort                  uint16         `env:"CHAINLINK_TLS_PORT" default:"6689"`
	TxAttemptLimit           uint16         `env:"CHAINLINK_TX_ATTEMPT_LIMIT" default:"10"`
}

// EnvVarName gets the environment variable name for a config schema field
func EnvVarName(field string) string {
	schemaT := reflect.TypeOf(ConfigSchema{})
	item, ok := schemaT.FieldByName(field)
	if !ok {
		log.Panicf("Invariant violated, no field of name %s found on ConfigSchema", field)
	}
	return item.Tag.Get("env")
}

func defaultValue(name string) (string, bool) {
	schemaT := reflect.TypeOf(ConfigSchema{})
	if item, ok := schemaT.FieldByName(name); ok {
		return item.Tag.Lookup("default")
	}
	log.Panicf("Invariant violated, no field of name %s found for defaultValue", name)
	return "", false
}

func zeroValue(name string) interface{} {
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
