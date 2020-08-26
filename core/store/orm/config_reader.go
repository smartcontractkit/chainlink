package orm

import (
	"math/big"
	"net/url"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/contrib/sessions"
	"go.uber.org/zap"
)

// ConfigReader represents just the read side of the config
type ConfigReader interface {
	AllowOrigins() string
	BlockBackfillDepth() uint64
	BridgeResponseURL() *url.URL
	ChainID() *big.Int
	ClientNodeURL() string
	DatabaseTimeout() models.Duration
	DatabaseURL() string
	DefaultMaxHTTPAttempts() uint
	DefaultHTTPLimit() int64
	DefaultHTTPTimeout() models.Duration
	Dev() bool
	FeatureExternalInitiators() bool
	FeatureFluxMonitor() bool
	MaximumServiceDuration() models.Duration
	MinimumServiceDuration() models.Duration
	EnableExperimentalAdapters() bool
	EnableBulletproofTxManager() bool
	EthBalanceMonitorBlockDelay() uint16
	EthGasBumpPercent() uint16
	EthGasBumpThreshold() uint64
	EthGasBumpWei() *big.Int
	EthGasLimitDefault() uint64
	EthGasPriceDefault() *big.Int
	EthMaxGasPriceWei() *big.Int
	EthFinalityDepth() uint
	EthHeadTrackerHistoryDepth() uint
	SetEthGasPriceDefault(value *big.Int) error
	EthereumURL() string
	GasUpdaterBlockDelay() uint16
	GasUpdaterBlockHistorySize() uint16
	GasUpdaterTransactionPercentile() uint16
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
	MinRequiredOutgoingConfirmations() uint64
	MinimumContractPayment() *assets.Link
	MinimumRequestExpiration() uint64
	MigrateDatabase() bool
	Port() uint16
	ReaperExpiration() models.Duration
	RootDir() string
	SecureCookies() bool
	SessionTimeout() models.Duration
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
