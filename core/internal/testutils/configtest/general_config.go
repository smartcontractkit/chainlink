package configtest

import (
	"fmt"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v4"
)

const (
	// RootDir the root directory for test
	RootDir = "/tmp/chainlink_test"
	// DefaultPeerID is the peer ID of the default p2p key
	DefaultPeerID              = "12D3KooWPjceQrSwdWXPyLLeABRXmuqt69Rg3sBYbU1Nft9HyQ6X"
	HeadSamplingIntervalInTest = 0 * time.Millisecond
)

var _ config.GeneralConfig = &TestGeneralConfig{}

type GeneralConfigOverrides struct {
	AdminCredentialsFile                      null.String
	AdvisoryLockID                            null.Int
	AllowOrigins                              null.String
	BlockBackfillDepth                        null.Int
	BlockBackfillSkip                         null.Bool
	ClientNodeURL                             null.String
	DatabaseTimeout                           *time.Duration
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
	GlobalEthTxReaperThreshold                *time.Duration
	GlobalEthTxResendAfterThreshold           *time.Duration
	GlobalEvmFinalityDepth                    null.Int
	GlobalEvmGasBumpPercent                   null.Int
	GlobalEvmGasBumpTxDepth                   null.Int
	GlobalEvmGasBumpWei                       *big.Int
	GlobalEvmGasLimitDefault                  null.Int
	GlobalEvmGasLimitMultiplier               null.Float
	GlobalEvmGasPriceDefault                  *big.Int
	GlobalEvmHeadTrackerHistoryDepth          null.Int
	GlobalEvmHeadTrackerMaxBufferSize         null.Int
	GlobalEvmHeadTrackerSamplingInterval      *time.Duration
	GlobalEvmLogBackfillBatchSize             null.Int
	GlobalEvmMaxGasPriceWei                   *big.Int
	GlobalEvmNonceAutoSync                    null.Bool
	GlobalEvmRPCDefaultBatchSize              null.Int
	GlobalFlagsContractAddress                null.String
	GlobalGasEstimatorMode                    null.String
	GlobalMinIncomingConfirmations            null.Int
	GlobalMinRequiredOutgoingConfirmations    null.Int
	GlobalMinimumContractPayment              *assets.Link
	KeeperMaximumGracePeriod                  null.Int
	KeeperMinimumRequiredConfirmations        null.Int
	KeeperRegistrySyncInterval                *time.Duration
	LogLevel                                  *config.LogLevel
	LogSQLStatements                          null.Bool
	LogToDisk                                 null.Bool
	OCRBootstrapCheckInterval                 *time.Duration
	OCRKeyBundleID                            null.String
	OCRObservationGracePeriod                 *time.Duration
	OCRObservationTimeout                     *time.Duration
	OCRTransmitterAddress                     *ethkey.EIP55Address
	P2PBootstrapPeers                         []string
	P2PListenPort                             null.Int
	P2PPeerID                                 *p2pkey.PeerID
	P2PPeerIDError                            error
	SecretGenerator                           config.SecretGenerator
	TriggerFallbackDBPollInterval             *time.Duration
}

// FIXME: This is a hack, the proper fix is here: https://app.clubhouse.io/chainlinklabs/story/15103/use-in-memory-event-broadcaster-instead-of-postgres-event-broadcaster-in-transactional-tests-so-it-actually-works
func (o *GeneralConfigOverrides) SetTriggerFallbackDBPollInterval(d time.Duration) {
	o.TriggerFallbackDBPollInterval = &d
}
func (o *GeneralConfigOverrides) SetOCRBootstrapCheckInterval(d time.Duration) {
	o.OCRBootstrapCheckInterval = &d
}
func (o *GeneralConfigOverrides) SetOCRObservationGracePeriod(d time.Duration) {
	o.OCRObservationGracePeriod = &d
}
func (o *GeneralConfigOverrides) SetOCRObservationTimeout(d time.Duration) {
	o.OCRObservationTimeout = &d
}
func (o *GeneralConfigOverrides) SetDefaultHTTPTimeout(d time.Duration) {
	o.DefaultHTTPTimeout = &d
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

func (c *TestGeneralConfig) P2PListenPort() uint16 {
	if c.Overrides.P2PListenPort.Valid {
		return uint16(c.Overrides.P2PListenPort.Int64)
	}
	return 12345
}

func (c *TestGeneralConfig) P2PPeerID() (p2pkey.PeerID, error) {
	if c.Overrides.P2PPeerIDError != nil {
		return "", c.Overrides.P2PPeerIDError
	}
	if c.Overrides.P2PPeerID != nil {
		return *c.Overrides.P2PPeerID, nil
	}
	defaultP2PPeerID, err := p2ppeer.Decode(DefaultPeerID)
	require.NoError(c.t, err)
	return p2pkey.PeerID(defaultP2PPeerID), nil
}

func (c *TestGeneralConfig) DatabaseTimeout() models.Duration {
	if c.Overrides.DatabaseTimeout != nil {
		return models.MustMakeDuration(*c.Overrides.DatabaseTimeout)
	}
	return models.MustMakeDuration(5 * time.Second)
}

func (c *TestGeneralConfig) GlobalLockRetryInterval() models.Duration {
	return models.MustMakeDuration(10 * time.Millisecond)
}

func (c *TestGeneralConfig) ORMMaxIdleConns() int {
	return 5
}

func (c *TestGeneralConfig) ORMMaxOpenConns() int {
	return 5
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

func (c *TestGeneralConfig) TriggerFallbackDBPollInterval() time.Duration {
	if c.Overrides.TriggerFallbackDBPollInterval != nil {
		return *c.Overrides.TriggerFallbackDBPollInterval
	}
	return c.GeneralConfig.TriggerFallbackDBPollInterval()
}

func (c *TestGeneralConfig) OCRBootstrapCheckInterval() time.Duration {
	if c.Overrides.OCRBootstrapCheckInterval != nil {
		return *c.Overrides.OCRBootstrapCheckInterval
	}
	return c.GeneralConfig.OCRBootstrapCheckInterval()
}

func (c *TestGeneralConfig) OCRObservationGracePeriod() time.Duration {
	if c.Overrides.OCRObservationGracePeriod != nil {
		return *c.Overrides.OCRObservationGracePeriod
	}
	return c.GeneralConfig.OCRObservationGracePeriod()
}

func (c *TestGeneralConfig) OCRObservationTimeout() time.Duration {
	if c.Overrides.OCRObservationTimeout != nil {
		return *c.Overrides.OCRObservationTimeout
	}
	return c.GeneralConfig.OCRObservationTimeout()
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

func (c *TestGeneralConfig) P2PBootstrapPeers() ([]string, error) {
	if c.Overrides.P2PBootstrapPeers != nil {
		return c.Overrides.P2PBootstrapPeers, nil
	}
	return c.GeneralConfig.P2PBootstrapPeers()
}

func (c *TestGeneralConfig) OCRKeyBundleID() (string, error) {
	if c.Overrides.OCRKeyBundleID.Valid {
		return c.Overrides.OCRKeyBundleID.String, nil
	}
	return c.GeneralConfig.OCRKeyBundleID()
}

func (c *TestGeneralConfig) OCRTransmitterAddress() (ethkey.EIP55Address, error) {
	if c.Overrides.OCRTransmitterAddress != nil {
		return *c.Overrides.OCRTransmitterAddress, nil
	}
	return c.GeneralConfig.OCRTransmitterAddress()
}

// CreateProductionLogger returns a custom logger for the config's root
// directory and LogLevel, with pretty printing for stdout. If LOG_TO_DISK is
// false, the logger will only log to stdout.
func (c *TestGeneralConfig) CreateProductionLogger() *logger.Logger {
	return logger.CreateProductionLogger(c.RootDir(), c.JSONConsole(), c.LogLevel().Level, c.LogToDisk())
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

func (c *TestGeneralConfig) BlockBackfillDepth() uint64 {
	if c.Overrides.BlockBackfillDepth.Valid {
		return uint64(c.Overrides.BlockBackfillDepth.Int64)
	}
	return c.GeneralConfig.BlockBackfillDepth()
}

func (c *TestGeneralConfig) KeeperMinimumRequiredConfirmations() uint64 {
	if c.Overrides.KeeperMinimumRequiredConfirmations.Valid {
		return uint64(c.Overrides.KeeperMinimumRequiredConfirmations.Int64)
	}
	return c.GeneralConfig.KeeperMinimumRequiredConfirmations()
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

func (c *TestGeneralConfig) LogLevel() config.LogLevel {
	if c.Overrides.LogLevel != nil {
		return *c.Overrides.LogLevel
	}
	return c.GeneralConfig.LogLevel()
}

func (c *TestGeneralConfig) LogSQLStatements() bool {
	if c.Overrides.LogSQLStatements.Valid {
		return c.Overrides.LogSQLStatements.Bool
	}
	return c.GeneralConfig.LogSQLStatements()
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

func (c *TestGeneralConfig) SetDialect(d dialects.DialectName) {
	c.Overrides.Dialect = d
}
