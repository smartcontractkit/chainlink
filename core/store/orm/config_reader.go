package orm

import (
	"math/big"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"go.uber.org/zap"
)

// ConfigReader represents just the read side of the config
type ConfigReader interface {
	AllowOrigins() string
	BridgeResponseURL() *url.URL
	ChainID() uint64
	ClientNodeURL() string
	DatabaseTimeout() time.Duration
	DatabaseURL() string
	DefaultHTTPLimit() int64
	Dev() bool
	MaximumServiceDuration() time.Duration
	MinimumServiceDuration() time.Duration
	EthGasBumpThreshold() uint64
	EthGasBumpWei() *big.Int
	EthGasPriceDefault() *big.Int
	SetEthGasPriceDefault(value *big.Int) error
	EthereumURL() string
	JSONConsole() bool
	LinkContractAddress() string
	ExplorerURL() *url.URL
	ExplorerAccessKey() string
	ExplorerSecret() string
	OracleContractAddress() *common.Address
	LogLevel() LogLevel
	LogToDisk() bool
	LogSQLStatements() bool
	MinIncomingConfirmations() uint32
	MinOutgoingConfirmations() uint64
	MinimumContractPayment() *assets.Link
	MinimumRequestExpiration() uint64
	Port() uint16
	ReaperExpiration() time.Duration
	RootDir() string
	SecureCookies() bool
	SessionTimeout() time.Duration
	TLSCertPath() string
	TLSHost() string
	TLSKeyPath() string
	TLSPort() uint16
	TLSRedirect() bool
	TxAttemptLimit() uint16
	KeysDir() string
	tlsDir() string
	KeyFile() string
	CertFile() string
	CreateProductionLogger() *zap.Logger
	SessionSecret() ([]byte, error)
	SessionOptions() sessions.Options
}
