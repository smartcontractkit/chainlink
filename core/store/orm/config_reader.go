package orm

import (
	"math/big"
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/contrib/sessions"
)

// ConfigReader represents just the read side of the config
type ConfigReader interface {
	AllowOrigins() string
	BlockBackfillDepth() uint64
	BridgeResponseURL() *url.URL
	CertFile() string
	ChainID() *big.Int
	ClientNodeURL() string
	CreateProductionLogger() *logger.Logger
	DatabaseMaximumTxDuration() time.Duration
	DatabaseTimeout() models.Duration
	DatabaseURL() url.URL
	DefaultHTTPAllowUnrestrictedNetworkAccess() bool
	DefaultHTTPLimit() int64
	DefaultHTTPTimeout() models.Duration
	DefaultMaxHTTPAttempts() uint
	Dev() bool
	EnableExperimentalAdapters() bool
	EthBalanceMonitorBlockDelay() uint16
	EthFinalityDepth() uint
	EthGasBumpPercent() uint16
	EthGasBumpThreshold() uint64
	EthGasBumpTxDepth() uint16
	EthGasBumpWei() *big.Int
	EthGasLimitDefault() uint64
	EthGasLimitMultiplier() float32
	EthGasPriceDefault() *big.Int
	EthHeadTrackerHistoryDepth() uint
	EthHeadTrackerMaxBufferSize() uint
	EthLogBackfillBatchSize() uint32
	EthMaxGasPriceWei() *big.Int
	EthNonceAutoSync() bool
	EthRPCDefaultBatchSize() uint32
	EthTxReaperInterval() time.Duration
	EthTxReaperThreshold() time.Duration
	EthTxResendAfterThreshold() time.Duration
	EthereumSecondaryURLs() []url.URL
	EthereumURL() string
	ExplorerAccessKey() string
	ExplorerSecret() string
	ExplorerURL() *url.URL
	FeatureExternalInitiators() bool
	FeatureFluxMonitor() bool
	FeatureOffchainReporting() bool
	BlockHistoryEstimatorBlockDelay() uint16
	BlockHistoryEstimatorBlockHistorySize() uint16
	BlockHistoryEstimatorTransactionPercentile() uint16
	InsecureSkipVerify() bool
	JSONConsole() bool
	KeeperDefaultTransactionQueueDepth() uint32
	KeeperRegistryCheckGasOverhead() uint64
	KeeperRegistryPerformGasOverhead() uint64
	KeeperRegistrySyncInterval() time.Duration
	KeeperMinimumRequiredConfirmations() uint64
	KeeperMaximumGracePeriod() int64
	KeyFile() string
	LinkContractAddress() string
	LogLevel() config.LogLevel
	LogSQLStatements() bool
	LogToDisk() bool
	MaximumServiceDuration() models.Duration
	MigrateDatabase() bool
	MinIncomingConfirmations() uint32
	MinRequiredOutgoingConfirmations() uint64
	MinimumContractPayment() *assets.Link
	MinimumRequestExpiration() uint64
	MinimumServiceDuration() models.Duration
	OCRTraceLogging() bool
	OperatorContractAddress() common.Address
	Port() uint16
	ReaperExpiration() models.Duration
	RootDir() string
	SecureCookies() bool
	SessionOptions() sessions.Options
	SessionSecret() ([]byte, error)
	SessionTimeout() models.Duration
	SetEthGasPriceDefault(value *big.Int) error
	TLSCertPath() string
	TLSDir() string
	TLSHost() string
	TLSKeyPath() string
	TLSPort() uint16
	TLSRedirect() bool
	TriggerFallbackDBPollInterval() time.Duration
}
