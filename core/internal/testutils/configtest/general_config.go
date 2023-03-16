package configtest

import (
	"fmt"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const (
	// RootDir the root directory for test
	RootDir       = "/tmp/chainlink_test"
	DefaultPeerID = "12D3KooWPjceQrSwdWXPyLLeABRXmuqt69Rg3sBYbU1Nft9HyQ6X"
)

var _ config.GeneralConfig = &TestGeneralConfig{}

// Deprecated: https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
type GeneralConfigOverrides struct {
	AdvisoryLockCheckInterval                       *time.Duration
	AdvisoryLockID                                  null.Int
	AllowOrigins                                    null.String
	BlockBackfillDepth                              null.Int
	BlockBackfillSkip                               null.Bool
	DatabaseURL                                     null.String
	DatabaseLockingMode                             null.String
	DatabaseDefaultLockTimeout                      *time.Duration
	DefaultChainID                                  *big.Int
	BridgeCacheTTL                                  *time.Duration
	DefaultHTTPTimeout                              *time.Duration
	HTTPServerWriteTimeout                          *time.Duration
	Dev                                             null.Bool
	ShutdownGracePeriod                             *time.Duration
	Dialect                                         dialects.DialectName
	EthereumURL                                     null.String
	GlobalBalanceMonitorEnabled                     null.Bool
	GlobalBlockEmissionIdleWarningThreshold         *time.Duration
	GlobalBlockHistoryEstimatorCheckInclusionBlocks null.Int
	GlobalChainType                                 null.String
	GlobalEthTxReaperThreshold                      *time.Duration
	GlobalEthTxResendAfterThreshold                 *time.Duration
	GlobalEvmEIP1559DynamicFees                     null.Bool
	GlobalEvmFinalityDepth                          null.Int
	GlobalEvmGasBumpPercent                         null.Int
	GlobalEvmGasBumpTxDepth                         null.Int
	GlobalEvmGasBumpWei                             *assets.Wei
	GlobalEvmGasFeeCapDefault                       *assets.Wei
	GlobalEvmGasLimitDefault                        null.Int
	GlobalEvmGasLimitMax                            null.Int
	GlobalEvmGasLimitMultiplier                     null.Float
	GlobalEvmGasPriceDefault                        *assets.Wei
	GlobalEvmGasTipCapDefault                       *assets.Wei
	GlobalEvmGasTipCapMinimum                       *assets.Wei
	GlobalEvmGasLimitOCRJobType                     null.Int
	GlobalEvmGasLimitDRJobType                      null.Int
	GlobalEvmGasLimitVRFJobType                     null.Int
	GlobalEvmGasLimitFMJobType                      null.Int
	GlobalEvmGasLimitKeeperJobType                  null.Int
	GlobalEvmHeadTrackerHistoryDepth                null.Int
	GlobalEvmHeadTrackerMaxBufferSize               null.Int
	GlobalEvmHeadTrackerSamplingInterval            *time.Duration
	GlobalEvmLogBackfillBatchSize                   null.Int
	GlobalEvmLogPollInterval                        *time.Duration
	GlobalEvmMaxGasPriceWei                         *assets.Wei
	GlobalEvmMinGasPriceWei                         *assets.Wei
	GlobalEvmNonceAutoSync                          null.Bool
	GlobalEvmRPCDefaultBatchSize                    null.Int
	GlobalEvmUseForwarders                          null.Bool
	GlobalFlagsContractAddress                      null.String
	GlobalGasEstimatorMode                          null.String
	GlobalMinIncomingConfirmations                  null.Int
	GlobalMinimumContractPayment                    *assets.Link
	GlobalOCRObservationGracePeriod                 time.Duration
	GlobalOCR2AutomationGasLimit                    null.Int
	KeeperRegistryMaxPerformDataSize                null.Int
	KeeperMaximumGracePeriod                        null.Int
	KeeperRegistrySyncInterval                      *time.Duration
	KeeperRegistrySyncUpkeepQueueSize               null.Int
	KeeperTurnLookBack                              null.Int
	LeaseLockDuration                               *time.Duration
	LeaseLockRefreshInterval                        *time.Duration
	LogFileDir                                      null.String
	LogLevel                                        *zapcore.Level
	DefaultLogLevel                                 *zapcore.Level
	LogSQL                                          null.Bool
	LogFileMaxSize                                  null.String
	LogFileMaxAge                                   null.Int
	LogFileMaxBackups                               null.Int
	TriggerFallbackDBPollInterval                   *time.Duration
	KeySpecific                                     map[string]types.ChainCfg
	LinkContractAddress                             null.String
	OperatorFactoryAddress                          null.String
	NodeNoNewHeadsThreshold                         *time.Duration
	JobPipelineReaperInterval                       *time.Duration

	// Feature Flags
	FeatureExternalInitiators null.Bool
	FeatureFeedsManager       null.Bool
	FeatureOffchainReporting  null.Bool
	FeatureOffchainReporting2 null.Bool
	FeatureLogPoller          null.Bool
	EVMEnabled                null.Bool
	EVMRPCEnabled             null.Bool
	P2PEnabled                null.Bool
	SolanaEnabled             null.Bool
	StarkNetEnabled           null.Bool

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
// SetTriggerFallbackDBPollInterval sets test override value for TriggerFallbackDBPollInterval
func (o *GeneralConfigOverrides) SetTriggerFallbackDBPollInterval(d time.Duration) {
	o.TriggerFallbackDBPollInterval = &d
}

// SetOCRBootstrapCheckInterval sets test override value for P2PBootstrapCheckInterval
func (o *GeneralConfigOverrides) SetOCRBootstrapCheckInterval(d time.Duration) {
	o.P2PBootstrapCheckInterval = &d
}

// SetOCRObservationTimeout sets test override value for OCRObservationTimeout
func (o *GeneralConfigOverrides) SetOCRObservationTimeout(d time.Duration) {
	o.OCRObservationTimeout = &d
}

// SetDefaultHTTPTimeout sets test override value for DefaultHTTPTimeout
func (o *GeneralConfigOverrides) SetDefaultHTTPTimeout(d time.Duration) {
	o.DefaultHTTPTimeout = &d
}

// SetBridgeCacheTTL sets test override value for BridgeCacheTTL
func (o *GeneralConfigOverrides) SetBridgeCacheTTL(d time.Duration) {
	o.BridgeCacheTTL = &d
}

// SetP2PV2DeltaDial sets test override value for P2PV2DeltaDial
func (o *GeneralConfigOverrides) SetP2PV2DeltaDial(d time.Duration) {
	o.P2PV2DeltaDial = &d
}

// SetP2PV2DeltaReconcile sets test override value for P2PV2DeltaReconcile
func (o *GeneralConfigOverrides) SetP2PV2DeltaReconcile(d time.Duration) {
	o.P2PV2DeltaReconcile = &d
}

// SetDefaultDatabaseLockTimeout sets test override value for DefaultDatabaseLockTimeout
func (o *GeneralConfigOverrides) SetDefaultDatabaseLockTimeout(d time.Duration) {
	o.DatabaseDefaultLockTimeout = &d
}

// TestGeneralConfig defaults to whatever config.NewGeneralConfig()
// gives but allows overriding certain methods
// Deprecated: https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
type TestGeneralConfig struct {
	config.GeneralConfig
	t         testing.TB
	rootdir   string
	Overrides GeneralConfigOverrides
}

// Deprecated: see v2.NewTestGeneralConfig
// NewTestGeneralConfig returns a legacy *TestGeneralConfig. Use v2.NewTestGeneralConfig instead.
// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func NewTestGeneralConfig(t *testing.T) *TestGeneralConfig {
	return NewTestGeneralConfigWithOverrides(t, GeneralConfigOverrides{})
}

// Deprecated: see v2.NewGeneralConfig
// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func NewTestGeneralConfigWithOverrides(t testing.TB, overrides GeneralConfigOverrides) *TestGeneralConfig {
	cfg := config.NewGeneralConfig(logger.TestLogger(t))
	return &TestGeneralConfig{
		cfg,
		t,
		genRootDir(t),
		overrides,
	}
}

// Deprecated: https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
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

func (c *TestGeneralConfig) BridgeCacheTTL() time.Duration {
	if c.Overrides.BridgeCacheTTL != nil {
		return *c.Overrides.BridgeCacheTTL
	}
	return c.GeneralConfig.BridgeCacheTTL()
}

func (c *TestGeneralConfig) DefaultChainID() *big.Int {
	if c.Overrides.DefaultChainID != nil {
		return c.Overrides.DefaultChainID
	}
	return big.NewInt(evmclient.NullClientChainID)
}

func (c *TestGeneralConfig) Dev() bool {
	if c.Overrides.Dev.Valid {
		return c.Overrides.Dev.Bool
	}
	return true
}

// ShutdownGracePeriod returns shutdown grace period duration.
func (c *TestGeneralConfig) ShutdownGracePeriod() time.Duration {
	if c.Overrides.ShutdownGracePeriod != nil {
		return *c.Overrides.ShutdownGracePeriod
	}
	return c.GeneralConfig.ShutdownGracePeriod()
}

func (c *TestGeneralConfig) MigrateDatabase() bool {
	return false
}

func (c *TestGeneralConfig) RootDir() string {
	return c.rootdir
}

// SetRootDir Added in order to not get a different dir on certain tests that validate this value
func (c *TestGeneralConfig) SetRootDir(dir string) {
	c.rootdir = dir
}

func (c *TestGeneralConfig) InsecureFastScrypt() bool {
	return true
}

func (c *TestGeneralConfig) SessionTimeout() models.Duration {
	return models.MustMakeDuration(2 * time.Minute)
}

func (c *TestGeneralConfig) ORMMaxIdleConns() int {
	return 20
}

func (c *TestGeneralConfig) ORMMaxOpenConns() int {
	// Set this to a reasonable number to enable test parallelisation (it requires one conn per db in tests)
	return 20
}

// EVMRPCEnabled overrides
func (c *TestGeneralConfig) EVMRPCEnabled() bool {
	if c.Overrides.EVMRPCEnabled.Valid {
		return c.Overrides.EVMRPCEnabled.Bool
	}
	return c.GeneralConfig.EVMRPCEnabled()
}

// SolanaEnabled allows Solana to be used
func (c *TestGeneralConfig) SolanaEnabled() bool {
	if c.Overrides.SolanaEnabled.Valid {
		return c.Overrides.SolanaEnabled.Bool
	}
	return c.GeneralConfig.SolanaEnabled()
}

// StarkNetEnabled allows StarkNet to be used
func (c *TestGeneralConfig) StarkNetEnabled() bool {
	if c.Overrides.StarkNetEnabled.Valid {
		return c.Overrides.StarkNetEnabled.Bool
	}
	return c.GeneralConfig.StarkNetEnabled()
}

func (c *TestGeneralConfig) EthereumURL() string {
	if c.Overrides.EthereumURL.Valid {
		return c.Overrides.EthereumURL.String
	}
	return c.GeneralConfig.EthereumURL()
}

func (c *TestGeneralConfig) GetDatabaseDialectConfiguredOrDefault() dialects.DialectName {
	if c.Overrides.Dialect != "" {
		return c.Overrides.Dialect
	}
	// Always return txdb for tests, if you want a non-transactional database
	// you must set an override explicitly
	return dialects.TransactionWrappedPostgres
}

func (c *TestGeneralConfig) DatabaseURL() url.URL {
	if c.Overrides.DatabaseURL.Valid {
		uri, err := url.Parse(c.Overrides.DatabaseURL.String)
		require.NoError(c.t, err)
		return *uri
	}
	return c.GeneralConfig.DatabaseURL()
}

// DatabaseLockingMode returns either overridden DatabaseLockingMode value or "none"
func (c *TestGeneralConfig) DatabaseLockingMode() string {
	if c.Overrides.DatabaseLockingMode.Valid {
		return c.Overrides.DatabaseLockingMode.String
	}
	// tests do not need DB locks, except for LockedDB tests
	return "none"
}

func (c *TestGeneralConfig) DatabaseDefaultLockTimeout() time.Duration {

	if c.Overrides.DatabaseDefaultLockTimeout != nil {
		return *c.Overrides.DatabaseDefaultLockTimeout

	}
	return c.GeneralConfig.DatabaseDefaultLockTimeout()
}

func (c *TestGeneralConfig) FeatureExternalInitiators() bool {
	if c.Overrides.FeatureExternalInitiators.Valid {
		return c.Overrides.FeatureExternalInitiators.Bool
	}
	return c.GeneralConfig.FeatureExternalInitiators()
}

func (c *TestGeneralConfig) FeatureFeedsManager() bool {
	if c.Overrides.FeatureFeedsManager.Valid {
		return c.Overrides.FeatureFeedsManager.Bool
	}
	return c.GeneralConfig.FeatureFeedsManager()
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

func (c *TestGeneralConfig) FeatureLogPoller() bool {
	if c.Overrides.FeatureLogPoller.Valid {
		return c.Overrides.FeatureLogPoller.Bool
	}
	return c.GeneralConfig.FeatureLogPoller()
}

// TriggerFallbackDBPollInterval returns the test configured value for TriggerFallbackDBPollInterval
func (c *TestGeneralConfig) TriggerFallbackDBPollInterval() time.Duration {
	if c.Overrides.TriggerFallbackDBPollInterval != nil {
		return *c.Overrides.TriggerFallbackDBPollInterval
	}
	return c.GeneralConfig.TriggerFallbackDBPollInterval()
}

// LogFileMaxSize allows to override the log file's max size before file rotation.
func (c *TestGeneralConfig) LogFileMaxSize() utils.FileSize {
	if c.Overrides.LogFileMaxSize.Valid {
		var val utils.FileSize

		err := val.UnmarshalText([]byte(c.Overrides.LogFileMaxSize.String))
		require.NoError(c.t, err)

		return val
	}
	return c.GeneralConfig.LogFileMaxSize()
}

// LogFileMaxAge allows to override the log file's max age before file rotation.
func (c *TestGeneralConfig) LogFileMaxAge() int64 {
	if c.Overrides.LogFileMaxAge.Valid {
		return c.Overrides.LogFileMaxAge.Int64
	}
	return int64(c.GeneralConfig.LogFileMaxAge())
}

// LogFileMaxBackups allows to override the max amount of old log files to retain.
func (c *TestGeneralConfig) LogFileMaxBackups() int64 {
	if c.Overrides.LogFileMaxBackups.Valid {
		return c.Overrides.LogFileMaxBackups.Int64
	}
	return int64(c.GeneralConfig.LogFileMaxBackups())
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

func (c *TestGeneralConfig) KeeperRegistryMaxPerformDataSize() uint32 {
	if c.Overrides.KeeperRegistryMaxPerformDataSize.Valid {
		return uint32(c.Overrides.KeeperRegistryMaxPerformDataSize.Int64)
	}
	return c.GeneralConfig.KeeperRegistryMaxPerformDataSize()
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

func (c *TestGeneralConfig) KeeperTurnLookBack() int64 {
	if c.Overrides.KeeperTurnLookBack.Valid {
		return c.Overrides.KeeperTurnLookBack.Int64
	}
	return c.GeneralConfig.KeeperTurnLookBack()
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
		return *c.Overrides.LogLevel
	}
	return c.GeneralConfig.LogLevel()
}

func (c *TestGeneralConfig) DefaultLogLevel() zapcore.Level {
	if c.Overrides.DefaultLogLevel != nil {
		return *c.Overrides.DefaultLogLevel
	}
	return c.GeneralConfig.DefaultLogLevel()
}

func (c *TestGeneralConfig) LogSQL() bool {
	if c.Overrides.LogSQL.Valid {
		return c.Overrides.LogSQL.Bool
	}
	return c.GeneralConfig.LogSQL()
}

// EVMEnabled overrides
func (c *TestGeneralConfig) EVMEnabled() bool {
	if c.Overrides.EVMEnabled.Valid {
		return c.Overrides.EVMEnabled.Bool
	}
	return c.GeneralConfig.EVMEnabled()
}

// P2PEnabled overrides
func (c *TestGeneralConfig) P2PEnabled() bool {
	if c.Overrides.P2PEnabled.Valid {
		return c.Overrides.P2PEnabled.Bool
	}
	return c.GeneralConfig.P2PEnabled()
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

// GlobalEvmGasFeeCapDefault is the override for EvmGasFeeCapDefault
func (c *TestGeneralConfig) GlobalEvmGasFeeCapDefault() (*assets.Wei, bool) {
	if c.Overrides.GlobalEvmGasFeeCapDefault != nil {
		return c.Overrides.GlobalEvmGasFeeCapDefault, true
	}
	return c.GeneralConfig.GlobalEvmGasFeeCapDefault()
}

func (c *TestGeneralConfig) GlobalEvmGasLimitDefault() (uint32, bool) {
	if c.Overrides.GlobalEvmGasLimitDefault.Valid {
		return uint32(c.Overrides.GlobalEvmGasLimitDefault.Int64), true
	}
	return c.GeneralConfig.GlobalEvmGasLimitDefault()
}

func (c *TestGeneralConfig) GlobalEvmGasLimitMax() (uint32, bool) {
	if c.Overrides.GlobalEvmGasLimitMax.Valid {
		return uint32(c.Overrides.GlobalEvmGasLimitMax.Int64), true
	}
	return c.GeneralConfig.GlobalEvmGasLimitMax()
}

func (c *TestGeneralConfig) GlobalEvmGasLimitOCRJobType() (uint32, bool) {
	if c.Overrides.GlobalEvmGasLimitOCRJobType.Valid {
		return uint32(c.Overrides.GlobalEvmGasLimitOCRJobType.Int64), true
	}
	return c.GeneralConfig.GlobalEvmGasLimitOCRJobType()
}

func (c *TestGeneralConfig) GlobalEvmGasLimitDRJobType() (uint32, bool) {
	if c.Overrides.GlobalEvmGasLimitDRJobType.Valid {
		return uint32(c.Overrides.GlobalEvmGasLimitDRJobType.Int64), true
	}
	return c.GeneralConfig.GlobalEvmGasLimitDRJobType()
}

func (c *TestGeneralConfig) GlobalEvmGasLimitVRFJobType() (uint32, bool) {
	if c.Overrides.GlobalEvmGasLimitVRFJobType.Valid {
		return uint32(c.Overrides.GlobalEvmGasLimitVRFJobType.Int64), true
	}
	return c.GeneralConfig.GlobalEvmGasLimitVRFJobType()
}

func (c *TestGeneralConfig) GlobalEvmGasLimitFMJobType() (uint32, bool) {
	if c.Overrides.GlobalEvmGasLimitFMJobType.Valid {
		return uint32(c.Overrides.GlobalEvmGasLimitFMJobType.Int64), true
	}
	return c.GeneralConfig.GlobalEvmGasLimitFMJobType()
}

func (c *TestGeneralConfig) GlobalEvmGasLimitKeeperJobType() (uint32, bool) {
	if c.Overrides.GlobalEvmGasLimitKeeperJobType.Valid {
		return uint32(c.Overrides.GlobalEvmGasLimitKeeperJobType.Int64), true
	}
	return c.GeneralConfig.GlobalEvmGasLimitKeeperJobType()
}

func (c *TestGeneralConfig) GlobalEvmGasLimitMultiplier() (float32, bool) {
	if c.Overrides.GlobalEvmGasLimitMultiplier.Valid {
		return float32(c.Overrides.GlobalEvmGasLimitMultiplier.Float64), true
	}
	return c.GeneralConfig.GlobalEvmGasLimitMultiplier()
}

func (c *TestGeneralConfig) GlobalEvmGasBumpWei() (*assets.Wei, bool) {
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

func (c *TestGeneralConfig) GlobalEvmGasPriceDefault() (*assets.Wei, bool) {
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

func (c *TestGeneralConfig) GlobalEvmLogPollInterval() (time.Duration, bool) {
	if c.Overrides.GlobalEvmLogPollInterval != nil {
		return *c.Overrides.GlobalEvmLogPollInterval, true
	}
	return c.GeneralConfig.GlobalEvmLogPollInterval()
}

func (c *TestGeneralConfig) GlobalEvmMaxGasPriceWei() (*assets.Wei, bool) {
	if c.Overrides.GlobalEvmMaxGasPriceWei != nil {
		return c.Overrides.GlobalEvmMaxGasPriceWei, true
	}
	return c.GeneralConfig.GlobalEvmMaxGasPriceWei()
}

func (c *TestGeneralConfig) GlobalEvmMinGasPriceWei() (*assets.Wei, bool) {
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

func (c *TestGeneralConfig) GlobalEvmGasTipCapDefault() (*assets.Wei, bool) {
	if c.Overrides.GlobalEvmGasTipCapDefault != nil {
		return c.Overrides.GlobalEvmGasTipCapDefault, true
	}
	return c.GeneralConfig.GlobalEvmGasTipCapDefault()
}

func (c *TestGeneralConfig) GlobalEvmGasTipCapMinimum() (*assets.Wei, bool) {
	if c.Overrides.GlobalEvmGasTipCapMinimum != nil {
		return c.Overrides.GlobalEvmGasTipCapMinimum, true
	}
	return c.GeneralConfig.GlobalEvmGasTipCapMinimum()
}

func (c *TestGeneralConfig) LeaseLockRefreshInterval() time.Duration {
	if c.Overrides.LeaseLockRefreshInterval != nil {
		return *c.Overrides.LeaseLockRefreshInterval
	}
	return c.GeneralConfig.LeaseLockRefreshInterval()
}

func (c *TestGeneralConfig) LeaseLockDuration() time.Duration {
	if c.Overrides.LeaseLockDuration != nil {
		return *c.Overrides.LeaseLockDuration
	}
	return c.GeneralConfig.LeaseLockDuration()
}

func (c *TestGeneralConfig) AdvisoryLockCheckInterval() time.Duration {
	if c.Overrides.AdvisoryLockCheckInterval != nil {
		return *c.Overrides.AdvisoryLockCheckInterval
	}
	return c.GeneralConfig.AdvisoryLockCheckInterval()
}

func (c *TestGeneralConfig) GlobalBlockEmissionIdleWarningThreshold() (time.Duration, bool) {
	if c.Overrides.GlobalBlockEmissionIdleWarningThreshold != nil {
		return *c.Overrides.GlobalBlockEmissionIdleWarningThreshold, true
	}
	return c.GeneralConfig.GlobalBlockEmissionIdleWarningThreshold()
}

func (c *TestGeneralConfig) GlobalBlockHistoryEstimatorCheckInclusionBlocks() (uint16, bool) {
	if c.Overrides.GlobalBlockHistoryEstimatorCheckInclusionBlocks.Valid {
		return uint16(c.Overrides.GlobalBlockHistoryEstimatorCheckInclusionBlocks.Int64), true
	}
	return c.GeneralConfig.GlobalBlockHistoryEstimatorCheckInclusionBlocks()
}

func (c *TestGeneralConfig) LogFileDir() string {
	if c.Overrides.LogFileDir.Valid {
		return c.Overrides.LogFileDir.String
	}
	return c.RootDir()
}

// GlobalLinkContractAddress allows to override the LINK contract address
func (c *TestGeneralConfig) GlobalLinkContractAddress() (string, bool) {
	if c.Overrides.LinkContractAddress.Valid {
		return c.Overrides.LinkContractAddress.String, true
	}
	return c.GeneralConfig.GlobalLinkContractAddress()
}

func (c *TestGeneralConfig) GlobalOCR2AutomationGasLimit() (uint32, bool) {
	if c.Overrides.GlobalOCR2AutomationGasLimit.Valid {
		return uint32(c.Overrides.GlobalOCR2AutomationGasLimit.Int64), true
	}
	return c.GeneralConfig.GlobalOCR2AutomationGasLimit()
}

// GlobalOperatorFactoryAddress allows to override the OperatorFactory contract address
func (c *TestGeneralConfig) GlobalOperatorFactoryAddress() (string, bool) {
	if c.Overrides.OperatorFactoryAddress.Valid {
		return c.Overrides.OperatorFactoryAddress.String, true
	}
	return c.GeneralConfig.GlobalOperatorFactoryAddress()
}

// GlobalNodeNoNewHeadsThreshold overrides NodeNoNewHeadsThreshold for all chains
func (c *TestGeneralConfig) GlobalNodeNoNewHeadsThreshold() (time.Duration, bool) {
	if c.Overrides.NodeNoNewHeadsThreshold != nil {
		return *c.Overrides.NodeNoNewHeadsThreshold, true
	}
	return c.GeneralConfig.GlobalNodeNoNewHeadsThreshold()
}

func (c *TestGeneralConfig) JobPipelineReaperInterval() time.Duration {
	if c.Overrides.JobPipelineReaperInterval != nil {
		return *c.Overrides.JobPipelineReaperInterval
	}
	return c.GeneralConfig.JobPipelineReaperInterval()
}

func (c *TestGeneralConfig) GlobalEvmUseForwarders() (bool, bool) {
	if c.Overrides.GlobalEvmUseForwarders.Valid {
		return c.Overrides.GlobalEvmUseForwarders.Bool, true
	}
	return c.GeneralConfig.GlobalEvmUseForwarders()
}
