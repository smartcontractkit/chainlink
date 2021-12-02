package configtest

import (
	"fmt"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	null "gopkg.in/guregu/null.v4"
)

const (
	// RootDir the root directory for test
	RootDir       = "/tmp/chainlink_test"
	DefaultPeerID = "12D3KooWPjceQrSwdWXPyLLeABRXmuqt69Rg3sBYbU1Nft9HyQ6X"
)

var _ config.GeneralConfig = &TestGeneralConfig{}

type GeneralConfigOverrides struct {
	AdminCredentialsFile                      null.String
	AdvisoryLockID                            null.Int
	AllowOrigins                              null.String
	BlockBackfillDepth                        null.Int
	BlockBackfillSkip                         null.Bool
	ClientNodeURL                             null.String
	DatabaseURL                               null.String
	DefaultChainID                            *big.Int
	DefaultHTTPAllowUnrestrictedNetworkAccess null.Bool
	DefaultHTTPTimeout                        *time.Duration
	DefaultMaxHTTPAttempts                    null.Int
	Dev                                       null.Bool
	Dialect                                   dialects.DialectName
	EVMDisabled                               null.Bool
	EthereumDisabled                          null.Bool
	FeatureExternalInitiators                 null.Bool
	GlobalBalanceMonitorEnabled               null.Bool
	GlobalChainType                           null.String
	GlobalEthTxReaperThreshold                *time.Duration
	GlobalEthTxResendAfterThreshold           *time.Duration
	GlobalEvmEIP1559DynamicFees               null.Bool
	GlobalEvmFinalityDepth                    null.Int
	GlobalEvmGasBumpPercent                   null.Int
	GlobalEvmGasBumpTxDepth                   null.Int
	GlobalEvmGasBumpWei                       *big.Int
	GlobalEvmGasLimitDefault                  null.Int
	GlobalEvmGasLimitMultiplier               null.Float
	GlobalEvmGasPriceDefault                  *big.Int
	GlobalEvmGasTipCapDefault                 *big.Int
	GlobalEvmGasTipCapMinimum                 *big.Int
	GlobalEvmHeadTrackerHistoryDepth          null.Int
	GlobalEvmHeadTrackerMaxBufferSize         null.Int
	GlobalEvmHeadTrackerSamplingInterval      *time.Duration
	GlobalEvmLogBackfillBatchSize             null.Int
	GlobalEvmMaxGasPriceWei                   *big.Int
	GlobalEvmMinGasPriceWei                   *big.Int
	GlobalEvmNonceAutoSync                    null.Bool
	GlobalEvmRPCDefaultBatchSize              null.Int
	GlobalFlagsContractAddress                null.String
	GlobalGasEstimatorMode                    null.String
	GlobalMinIncomingConfirmations            null.Int
	GlobalMinRequiredOutgoingConfirmations    null.Int
	GlobalMinimumContractPayment              *assets.Link
	GlobalOCRObservationGracePeriod           time.Duration
	KeeperMaximumGracePeriod                  null.Int
	KeeperRegistrySyncInterval                *time.Duration
	KeeperRegistrySyncUpkeepQueueSize         null.Int
	LogLevel                                  *config.LogLevel
	DefaultLogLevel                           *config.LogLevel
	LogSQL                                    null.Bool
	LogToDisk                                 null.Bool
	SecretGenerator                           config.SecretGenerator
	TriggerFallbackDBPollInterval             *time.Duration
	KeySpecific                               map[string]types.ChainCfg
	FeatureOffchainReporting                  null.Bool
	FeatureOffchainReporting2                 null.Bool

	// OCR v2
	OCR2DatabaseTimeout *time.Duration

	// OCR v1
	OCRKeyBundleID            null.String
	OCRDatabaseTimeout        *time.Duration
	OCRObservationGracePeriod *time.Duration
	OCRObservationTimeout     *time.Duration
	OCRTransmitterAddress     *ethkey.EIP55Address

	// P2P v1 and V2
	P2PPeerID          p2pkey.PeerID
	P2PPeerIDError     error
	P2PNetworkingStack ocrnetworking.NetworkingStack

	// P2P v1
	P2PBootstrapCheckInterval *time.Duration
	P2PBootstrapPeers         []string
	P2PListenPort             null.Int

	// P2PV2
	P2PV2ListenAddresses   []string
	P2PV2AnnounceAddresses []string
	P2PV2Bootstrappers     []ocrcommontypes.BootstrapperLocator
	P2PV2DeltaDial         *time.Duration
	P2PV2DeltaReconcile    *time.Duration
}

// FIXME: This is a hack, the proper fix is here: https://app.clubhouse.io/chainlinklabs/story/15103/use-in-memory-event-broadcaster-instead-of-postgres-event-broadcaster-in-transactional-tests-so-it-actually-works
func (o *GeneralConfigOverrides) SetTriggerFallbackDBPollInterval(d time.Duration) {
	o.TriggerFallbackDBPollInterval = &d
}
func (o *GeneralConfigOverrides) SetOCRBootstrapCheckInterval(d time.Duration) {
	o.P2PBootstrapCheckInterval = &d
}
func (o *GeneralConfigOverrides) SetOCRObservationTimeout(d time.Duration) {
	o.OCRObservationTimeout = &d
}
func (o *GeneralConfigOverrides) SetDefaultHTTPTimeout(d time.Duration) {
	o.DefaultHTTPTimeout = &d
}
func (o *GeneralConfigOverrides) SetP2PV2DeltaDial(d time.Duration) {
	o.P2PV2DeltaDial = &d
}
func (o *GeneralConfigOverrides) SetP2PV2DeltaReconcile(d time.Duration) {
	o.P2PV2DeltaReconcile = &d
}

// TestGeneralConfig defaults to whatever config.NewGeneralConfig()
// gives but allows overriding certain methods
type TestGeneralConfig struct {
	config.GeneralConfig
	t         testing.TB
	rootdir   string
	Overrides GeneralConfigOverrides
}

func NewTestGeneralConfig(t *testing.T) *TestGeneralConfig {
	return NewTestGeneralConfigWithOverrides(t, GeneralConfigOverrides{})
}

func NewTestGeneralConfigWithOverrides(t testing.TB, overrides GeneralConfigOverrides) *TestGeneralConfig {
	cfg := config.NewGeneralConfig()
	return &TestGeneralConfig{
		cfg,
		t,
		genRootDir(t),
		overrides,
	}
}

func genRootDir(t testing.TB) string {
	name := fmt.Sprintf("%d-%d", time.Now().UnixNano(), 0)
	dir := filepath.Join(RootDir, name)
	if err := utils.EnsureDirAndMaxPerms(dir, os.FileMode(0700)); err != nil {
		t.Fatalf(`Error creating root directory "%s": %+v`, dir, err)
	}
	return dir
}

func (c *TestGeneralConfig) GetAdvisoryLockIDConfiguredOrDefault() int64 {
	if c.Overrides.AdvisoryLockID.Valid {
		return c.Overrides.AdvisoryLockID.Int64
	}
	return c.GeneralConfig.GetAdvisoryLockIDConfiguredOrDefault()
}

func (c *TestGeneralConfig) BridgeResponseURL() *url.URL {
	uri, err := url.Parse("http://localhost:6688")
	require.NoError(c.t, err)
	return uri
}

func (c *TestGeneralConfig) DefaultChainID() *big.Int {
	if c.Overrides.DefaultChainID != nil {
		return c.Overrides.DefaultChainID
	}
	return big.NewInt(eth.NullClientChainID)
}

func (c *TestGeneralConfig) Dev() bool {
	if c.Overrides.Dev.Valid {
		return c.Overrides.Dev.Bool
	}
	return true
}

func (c *TestGeneralConfig) MigrateDatabase() bool {
	return false
}

func (c *TestGeneralConfig) RootDir() string {
	return c.rootdir
}

func (c *TestGeneralConfig) SessionTimeout() models.Duration {
	return models.MustMakeDuration(2 * time.Minute)
}

func (c *TestGeneralConfig) InsecureFastScrypt() bool {
	return true
}

func (c *TestGeneralConfig) GlobalLockRetryInterval() models.Duration {
	return models.MustMakeDuration(10 * time.Millisecond)
}

func (c *TestGeneralConfig) ORMMaxIdleConns() int {
	return 5
}

func (c *TestGeneralConfig) ORMMaxOpenConns() int {
	// HACK: txdb does not appear to use connection pooling properly, so that
	// if this value is not large enough instead of waiting for a connection the
	// database call will fail with "conn busy" or some other cryptic error
	return 20
}

func (c *TestGeneralConfig) LogSQLMigrations() bool {
	return false
}

func (c *TestGeneralConfig) EthereumDisabled() bool {
	if c.Overrides.EthereumDisabled.Valid {
		return c.Overrides.EthereumDisabled.Bool
	}
	return c.GeneralConfig.EthereumDisabled()
}

func (c *TestGeneralConfig) SessionSecret() ([]byte, error) {
	if c.Overrides.SecretGenerator != nil {
		return c.Overrides.SecretGenerator.Generate(c.RootDir())
	}
	return c.GeneralConfig.SessionSecret()
}

func (c *TestGeneralConfig) GetDatabaseDialectConfiguredOrDefault() dialects.DialectName {
	if c.Overrides.Dialect != "" {
		return c.Overrides.Dialect
	}
	// Always return txdb for tests, if you want a non-transactional database
	// you must set an override explicitly
	return "txdb"
}

func (c *TestGeneralConfig) ClientNodeURL() string {
	if c.Overrides.ClientNodeURL.Valid {
		return c.Overrides.ClientNodeURL.String
	}
	return c.GeneralConfig.ClientNodeURL()
}

func (c *TestGeneralConfig) DatabaseURL() url.URL {
	if c.Overrides.DatabaseURL.Valid {
		uri, err := url.Parse(c.Overrides.DatabaseURL.String)
		require.NoError(c.t, err)
		return *uri
	}
	return c.GeneralConfig.DatabaseURL()
}

func (c *TestGeneralConfig) FeatureExternalInitiators() bool {
	if c.Overrides.FeatureExternalInitiators.Valid {
		return c.Overrides.FeatureExternalInitiators.Bool
	}
	return c.GeneralConfig.FeatureExternalInitiators()
}

func (c *TestGeneralConfig) FeatureOffchainReporting() bool {
	if c.Overrides.FeatureOffchainReporting.Valid {
		return c.Overrides.FeatureOffchainReporting.Bool
	}
	return c.GeneralConfig.FeatureOffchainReporting()
}

func (c *TestGeneralConfig) FeatureOffchainReporting2() bool {
	if c.Overrides.FeatureOffchainReporting2.Valid {
		return c.Overrides.FeatureOffchainReporting2.Bool
	}
	return c.GeneralConfig.FeatureOffchainReporting2()
}

func (c *TestGeneralConfig) TriggerFallbackDBPollInterval() time.Duration {
	if c.Overrides.TriggerFallbackDBPollInterval != nil {
		return *c.Overrides.TriggerFallbackDBPollInterval
	}
	return c.GeneralConfig.TriggerFallbackDBPollInterval()
}

func (c *TestGeneralConfig) LogToDisk() bool {
	if c.Overrides.LogToDisk.Valid {
		return c.Overrides.LogToDisk.Bool
	}
	return c.GeneralConfig.LogToDisk()
}

func (c *TestGeneralConfig) DefaultMaxHTTPAttempts() uint {
	if c.Overrides.DefaultMaxHTTPAttempts.Valid {
		return uint(c.Overrides.DefaultMaxHTTPAttempts.Int64)
	}
	return c.GeneralConfig.DefaultMaxHTTPAttempts()
}

func (c *TestGeneralConfig) AdminCredentialsFile() string {
	if c.Overrides.AdminCredentialsFile.Valid {
		return c.Overrides.AdminCredentialsFile.String
	}
	return c.GeneralConfig.AdminCredentialsFile()
}

func (c *TestGeneralConfig) DefaultHTTPAllowUnrestrictedNetworkAccess() bool {
	if c.Overrides.DefaultHTTPAllowUnrestrictedNetworkAccess.Valid {
		return c.Overrides.DefaultHTTPAllowUnrestrictedNetworkAccess.Bool
	}
	return c.GeneralConfig.DefaultHTTPAllowUnrestrictedNetworkAccess()
}

func (c *TestGeneralConfig) DefaultHTTPTimeout() models.Duration {
	if c.Overrides.DefaultHTTPTimeout != nil {
		return models.MustMakeDuration(*c.Overrides.DefaultHTTPTimeout)
	}
	return c.GeneralConfig.DefaultHTTPTimeout()
}

func (c *TestGeneralConfig) KeeperRegistrySyncInterval() time.Duration {
	if c.Overrides.KeeperRegistrySyncInterval != nil {
		return *c.Overrides.KeeperRegistrySyncInterval
	}
	return c.GeneralConfig.KeeperRegistrySyncInterval()
}

func (c *TestGeneralConfig) KeeperRegistrySyncUpkeepQueueSize() uint32 {
	if c.Overrides.KeeperRegistrySyncUpkeepQueueSize.Valid {
		return uint32(c.Overrides.KeeperRegistrySyncUpkeepQueueSize.Int64)
	}
	return c.GeneralConfig.KeeperRegistrySyncUpkeepQueueSize()
}

func (c *TestGeneralConfig) BlockBackfillDepth() uint64 {
	if c.Overrides.BlockBackfillDepth.Valid {
		return uint64(c.Overrides.BlockBackfillDepth.Int64)
	}
	return c.GeneralConfig.BlockBackfillDepth()
}

func (c *TestGeneralConfig) KeeperMaximumGracePeriod() int64 {
	if c.Overrides.KeeperMaximumGracePeriod.Valid {
		return c.Overrides.KeeperMaximumGracePeriod.Int64
	}
	return c.GeneralConfig.KeeperMaximumGracePeriod()
}

func (c *TestGeneralConfig) BlockBackfillSkip() bool {
	if c.Overrides.BlockBackfillSkip.Valid {
		return c.Overrides.BlockBackfillSkip.Bool
	}
	return c.GeneralConfig.BlockBackfillSkip()
}

func (c *TestGeneralConfig) AllowOrigins() string {
	if c.Overrides.AllowOrigins.Valid {
		return c.Overrides.AllowOrigins.String
	}
	return c.GeneralConfig.AllowOrigins()
}

func (c *TestGeneralConfig) LogLevel() zapcore.Level {
	if c.Overrides.LogLevel != nil {
		return c.Overrides.LogLevel.Level
	}
	return c.GeneralConfig.LogLevel()
}

func (c *TestGeneralConfig) DefaultLogLevel() zapcore.Level {
	if c.Overrides.DefaultLogLevel != nil {
		return c.Overrides.DefaultLogLevel.Level
	}
	return c.GeneralConfig.DefaultLogLevel()
}

func (c *TestGeneralConfig) LogSQL() bool {
	if c.Overrides.LogSQL.Valid {
		return c.Overrides.LogSQL.Bool
	}
	return c.GeneralConfig.LogSQL()
}

func (c *TestGeneralConfig) EVMDisabled() bool {
	if c.Overrides.EVMDisabled.Valid {
		return c.Overrides.EVMDisabled.Bool
	}
	return c.GeneralConfig.EVMDisabled()
}

func (c *TestGeneralConfig) GlobalGasEstimatorMode() (string, bool) {
	if c.Overrides.GlobalGasEstimatorMode.Valid {
		return c.Overrides.GlobalGasEstimatorMode.String, true
	}
	return c.GeneralConfig.GlobalGasEstimatorMode()
}

func (c *TestGeneralConfig) GlobalChainType() (string, bool) {
	if c.Overrides.GlobalChainType.Valid {
		return c.Overrides.GlobalChainType.String, true
	}
	return c.GeneralConfig.GlobalChainType()
}

func (c *TestGeneralConfig) GlobalEvmNonceAutoSync() (bool, bool) {
	if c.Overrides.GlobalEvmNonceAutoSync.Valid {
		return c.Overrides.GlobalEvmNonceAutoSync.Bool, true
	}
	return c.GeneralConfig.GlobalEvmNonceAutoSync()
}
func (c *TestGeneralConfig) GlobalBalanceMonitorEnabled() (bool, bool) {
	if c.Overrides.GlobalBalanceMonitorEnabled.Valid {
		return c.Overrides.GlobalBalanceMonitorEnabled.Bool, true
	}
	return c.GeneralConfig.GlobalBalanceMonitorEnabled()
}

func (c *TestGeneralConfig) GlobalEvmGasLimitDefault() (uint64, bool) {
	if c.Overrides.GlobalEvmGasLimitDefault.Valid {
		return uint64(c.Overrides.GlobalEvmGasLimitDefault.Int64), true
	}
	return c.GeneralConfig.GlobalEvmGasLimitDefault()
}

func (c *TestGeneralConfig) GlobalEvmGasLimitMultiplier() (float32, bool) {
	if c.Overrides.GlobalEvmGasLimitMultiplier.Valid {
		return float32(c.Overrides.GlobalEvmGasLimitMultiplier.Float64), true
	}
	return c.GeneralConfig.GlobalEvmGasLimitMultiplier()
}

func (c *TestGeneralConfig) GlobalEvmGasBumpWei() (*big.Int, bool) {
	if c.Overrides.GlobalEvmGasBumpWei != nil {
		return c.Overrides.GlobalEvmGasBumpWei, true
	}
	return c.GeneralConfig.GlobalEvmGasBumpWei()
}

func (c *TestGeneralConfig) GlobalEvmGasBumpPercent() (uint16, bool) {
	if c.Overrides.GlobalEvmGasBumpPercent.Valid {
		return uint16(c.Overrides.GlobalEvmGasBumpPercent.Int64), true
	}
	return c.GeneralConfig.GlobalEvmGasBumpPercent()
}

func (c *TestGeneralConfig) GlobalEvmGasPriceDefault() (*big.Int, bool) {
	if c.Overrides.GlobalEvmGasPriceDefault != nil {
		return c.Overrides.GlobalEvmGasPriceDefault, true
	}
	return c.GeneralConfig.GlobalEvmGasPriceDefault()
}

func (c *TestGeneralConfig) GlobalEvmRPCDefaultBatchSize() (uint32, bool) {
	if c.Overrides.GlobalEvmRPCDefaultBatchSize.Valid {
		return uint32(c.Overrides.GlobalEvmRPCDefaultBatchSize.Int64), true
	}
	return c.GeneralConfig.GlobalEvmRPCDefaultBatchSize()
}

func (c *TestGeneralConfig) GlobalEvmFinalityDepth() (uint32, bool) {
	if c.Overrides.GlobalEvmFinalityDepth.Valid {
		return uint32(c.Overrides.GlobalEvmFinalityDepth.Int64), true
	}
	return c.GeneralConfig.GlobalEvmFinalityDepth()
}

func (c *TestGeneralConfig) GlobalEvmLogBackfillBatchSize() (uint32, bool) {
	if c.Overrides.GlobalEvmLogBackfillBatchSize.Valid {
		return uint32(c.Overrides.GlobalEvmLogBackfillBatchSize.Int64), true
	}
	return c.GeneralConfig.GlobalEvmLogBackfillBatchSize()
}

func (c *TestGeneralConfig) GlobalEvmMaxGasPriceWei() (*big.Int, bool) {
	if c.Overrides.GlobalEvmMaxGasPriceWei != nil {
		return c.Overrides.GlobalEvmMaxGasPriceWei, true
	}
	return c.GeneralConfig.GlobalEvmMaxGasPriceWei()
}

func (c *TestGeneralConfig) GlobalEvmMinGasPriceWei() (*big.Int, bool) {
	if c.Overrides.GlobalEvmMinGasPriceWei != nil {
		return c.Overrides.GlobalEvmMinGasPriceWei, true
	}
	return c.GeneralConfig.GlobalEvmMinGasPriceWei()
}

func (c *TestGeneralConfig) GlobalEvmGasBumpTxDepth() (uint16, bool) {
	if c.Overrides.GlobalEvmGasBumpTxDepth.Valid {
		return uint16(c.Overrides.GlobalEvmGasBumpTxDepth.Int64), true
	}
	return c.GeneralConfig.GlobalEvmGasBumpTxDepth()
}

func (c *TestGeneralConfig) GlobalEthTxResendAfterThreshold() (time.Duration, bool) {
	if c.Overrides.GlobalEthTxResendAfterThreshold != nil {
		return *c.Overrides.GlobalEthTxResendAfterThreshold, true
	}
	return c.GeneralConfig.GlobalEthTxResendAfterThreshold()
}

func (c *TestGeneralConfig) GlobalMinIncomingConfirmations() (uint32, bool) {
	if c.Overrides.GlobalMinIncomingConfirmations.Valid {
		return uint32(c.Overrides.GlobalMinIncomingConfirmations.Int64), true
	}
	return c.GeneralConfig.GlobalMinIncomingConfirmations()
}

func (c *TestGeneralConfig) GlobalMinimumContractPayment() (*assets.Link, bool) {
	if c.Overrides.GlobalMinimumContractPayment != nil {
		return c.Overrides.GlobalMinimumContractPayment, true
	}
	return c.GeneralConfig.GlobalMinimumContractPayment()
}

func (c *TestGeneralConfig) GlobalFlagsContractAddress() (string, bool) {
	if c.Overrides.GlobalFlagsContractAddress.Valid {
		return c.Overrides.GlobalFlagsContractAddress.String, true
	}
	return c.GeneralConfig.GlobalFlagsContractAddress()
}

func (c *TestGeneralConfig) GlobalMinRequiredOutgoingConfirmations() (uint64, bool) {
	if c.Overrides.GlobalMinRequiredOutgoingConfirmations.Valid {
		return uint64(c.Overrides.GlobalMinRequiredOutgoingConfirmations.Int64), true
	}
	return c.GeneralConfig.GlobalMinRequiredOutgoingConfirmations()
}

func (c *TestGeneralConfig) GlobalEvmHeadTrackerMaxBufferSize() (uint32, bool) {
	if c.Overrides.GlobalEvmHeadTrackerMaxBufferSize.Valid {
		return uint32(c.Overrides.GlobalEvmHeadTrackerMaxBufferSize.Int64), true
	}
	return c.GeneralConfig.GlobalEvmHeadTrackerMaxBufferSize()
}

func (c *TestGeneralConfig) GlobalEvmHeadTrackerHistoryDepth() (uint32, bool) {
	if c.Overrides.GlobalEvmHeadTrackerHistoryDepth.Valid {
		return uint32(c.Overrides.GlobalEvmHeadTrackerHistoryDepth.Int64), true
	}
	return c.GeneralConfig.GlobalEvmHeadTrackerHistoryDepth()
}

func (c *TestGeneralConfig) GlobalEvmHeadTrackerSamplingInterval() (time.Duration, bool) {
	if c.Overrides.GlobalEvmHeadTrackerSamplingInterval != nil {
		return *c.Overrides.GlobalEvmHeadTrackerSamplingInterval, true
	}
	return c.GeneralConfig.GlobalEvmHeadTrackerSamplingInterval()
}

func (c *TestGeneralConfig) GlobalEthTxReaperThreshold() (time.Duration, bool) {
	if c.Overrides.GlobalEthTxReaperThreshold != nil {
		return *c.Overrides.GlobalEthTxReaperThreshold, true
	}
	return c.GeneralConfig.GlobalEthTxReaperThreshold()
}

func (c *TestGeneralConfig) GlobalEvmEIP1559DynamicFees() (bool, bool) {
	if c.Overrides.GlobalEvmEIP1559DynamicFees.Valid {
		return c.Overrides.GlobalEvmEIP1559DynamicFees.Bool, true
	}
	return c.GeneralConfig.GlobalEvmEIP1559DynamicFees()
}

func (c *TestGeneralConfig) GlobalEvmGasTipCapDefault() (*big.Int, bool) {
	if c.Overrides.GlobalEvmGasTipCapDefault != nil {
		return c.Overrides.GlobalEvmGasTipCapDefault, true
	}
	return c.GeneralConfig.GlobalEvmGasTipCapDefault()
}

func (c *TestGeneralConfig) GlobalEvmGasTipCapMinimum() (*big.Int, bool) {
	if c.Overrides.GlobalEvmGasTipCapMinimum != nil {
		return c.Overrides.GlobalEvmGasTipCapMinimum, true
	}
	return c.GeneralConfig.GlobalEvmGasTipCapMinimum()
}

func (c *TestGeneralConfig) SetDialect(d dialects.DialectName) {
	c.Overrides.Dialect = d
}

// There is no need for any database application locking in tests
func (c *TestGeneralConfig) DatabaseLockingMode() string {
	return "none"
}
