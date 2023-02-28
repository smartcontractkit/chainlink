package chainlink

import (
	_ "embed"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/libocr/commontypes"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"

	simplelogger "github.com/smartcontractkit/chainlink-relay/pkg/logger"

	evmcfg "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/core/chains/solana"
	"github.com/smartcontractkit/chainlink/core/chains/starknet"
	coreconfig "github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/config/parse"
	v2 "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/logger/audit"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type ConfigV2 interface {
	// ConfigTOML returns both the user provided and effective configuration as TOML.
	ConfigTOML() (user, effective string)
}

// generalConfig is a wrapper to adapt Config to the config.GeneralConfig interface.
type generalConfig struct {
	lggr simplelogger.Logger

	inputTOML     string // user input, normalized via de/re-serialization
	effectiveTOML string // with default values included
	secretsTOML   string // with env overdies includes, redacted

	c       *Config // all fields non-nil (unless the legacy method signature return a pointer)
	secrets *Secrets

	logLevelDefault zapcore.Level

	appIDOnce sync.Once

	randomP2PPort     uint16
	randomP2PPortOnce sync.Once

	logMu sync.RWMutex // for the mutable fields Log.Level & Log.SQL

	passwordMu sync.RWMutex // passwords are set after initialization
}

// GeneralConfigOpts holds configuration options for creating a coreconfig.GeneralConfig via New().
//
// See ParseTOML to initilialize Config and Secrets from TOML.
type GeneralConfigOpts struct {
	Config
	Secrets

	// OverrideFn is a *test-only* hook to override effective values.
	OverrideFn func(*Config, *Secrets)

	SkipEnv bool
}

// ParseTOML sets Config and Secrets from the given TOML strings.
func (o *GeneralConfigOpts) ParseTOML(config, secrets string) (err error) {
	return multierr.Combine(o.ParseConfig(config), o.ParseSecrets(secrets))
}

// ParseConfig sets Config from the given TOML string, overriding any existing duplicate Config fields.
func (o *GeneralConfigOpts) ParseConfig(config string) error {
	var c Config
	if err2 := v2.DecodeTOML(strings.NewReader(config), &c); err2 != nil {
		return fmt.Errorf("failed to decode config TOML: %w", err2)
	}
	o.Config.SetFrom(&c)
	return nil
}

// ParseSecrets sets Secrets from the given TOML string.
func (o *GeneralConfigOpts) ParseSecrets(secrets string) (err error) {
	if err2 := v2.DecodeTOML(strings.NewReader(secrets), &o.Secrets); err2 != nil {
		return fmt.Errorf("failed to decode secrets TOML: %w", err2)
	}
	return nil
}

// New returns a coreconfig.GeneralConfig for the given options.
func (o GeneralConfigOpts) New(lggr logger.Logger) (coreconfig.GeneralConfig, error) {
	cfg, err := o.init()
	if err != nil {
		return nil, err
	}
	cfg.lggr = lggr
	return cfg, nil
}

// NewAndLogger returns a coreconfig.GeneralConfig for the given options, and a logger.Logger (with close func).
func (o GeneralConfigOpts) NewAndLogger() (coreconfig.GeneralConfig, logger.Logger, func() error, error) {
	cfg, err := o.init()
	if err != nil {
		return nil, nil, nil, err
	}

	// placeholder so we can call config methods to bootstrap the real logger
	cfg.lggr, err = simplelogger.New()
	if err != nil {
		return nil, nil, nil, err
	}
	lggrCfg := logger.Config{
		LogLevel:       cfg.LogLevel(),
		Dir:            cfg.LogFileDir(),
		JsonConsole:    cfg.JSONConsole(),
		UnixTS:         cfg.LogUnixTimestamps(),
		FileMaxSizeMB:  int(cfg.LogFileMaxSize() / utils.MB),
		FileMaxAgeDays: int(cfg.LogFileMaxAge()),
		FileMaxBackups: int(cfg.LogFileMaxBackups()),
	}
	lggr, closeLggr := lggrCfg.New()

	cfg.lggr = lggr
	return cfg, lggr, closeLggr, nil
}

// new returns a new generalConfig, but with a nil lggr.
func (o *GeneralConfigOpts) init() (*generalConfig, error) {
	input, err := o.Config.TOMLString()
	if err != nil {
		return nil, err
	}

	o.Config.setDefaults()
	if !o.SkipEnv {
		o.Config.DevMode = v2.EnvDev.IsTrue()

		err = o.Secrets.setEnv()
		if err != nil {
			return nil, err
		}
	}

	if fn := o.OverrideFn; fn != nil {
		fn(&o.Config, &o.Secrets)
	}

	effective, err := o.Config.TOMLString()
	if err != nil {
		return nil, err
	}

	secrets, err := o.Secrets.TOMLString()
	if err != nil {
		return nil, err
	}

	cfg := &generalConfig{
		inputTOML: input, effectiveTOML: effective, secretsTOML: secrets,
		c: &o.Config, secrets: &o.Secrets,
	}
	if lvl := o.Config.Log.Level; lvl != nil {
		cfg.logLevelDefault = zapcore.Level(*lvl)
	}

	if err2 := utils.EnsureDirAndMaxPerms(cfg.RootDir(), os.FileMode(0700)); err2 != nil {
		return nil, fmt.Errorf(`failed to create root directory %q: %w`, cfg.RootDir(), err2)
	}

	return cfg, nil
}

func (g *generalConfig) EVMConfigs() evmcfg.EVMConfigs {
	return g.c.EVM
}

func (g *generalConfig) SolanaConfigs() solana.SolanaConfigs {
	return g.c.Solana
}

func (g *generalConfig) StarknetConfigs() starknet.StarknetConfigs {
	return g.c.Starknet
}

func (g *generalConfig) Validate() error {
	_, err := utils.MultiErrorList(multierr.Combine(
		validateEnv(),
		g.c.Validate(),
		g.secrets.Validate()))
	return err
}

//go:embed cfgtest/dump/empty-strings.env
var emptyStringsEnv string

var legacyEnvToV2 = map[string]string{
	"CHAINLINK_DEV": "CL_DEV",

	"DATABASE_URL":                            "CL_DATABASE_URL",
	"DATABASE_BACKUP_URL":                     "CL_DATABASE_BACKUP_URL",
	"SKIP_DATABASE_PASSWORD_COMPLEXITY_CHECK": "CL_DATABASE_ALLOW_SIMPLE_PASSWORDS",

	"EXPLORER_ACCESS_KEY": "CL_EXPLORER_ACCESS_KEY",
	"EXPLORER_SECRET":     "CL_EXPLORER_SECRET",

	"PYROSCOPE_AUTH_TOKEN": "CL_PYROSCOPE_AUTH_TOKEN",

	"LOG_COLOR":          "CL_LOG_COLOR",
	"LOG_SQL_MIGRATIONS": "CL_LOG_SQL_MIGRATIONS",
}

// validateEnv returns an error if any legacy environment variables are set, unless a v2 equivalent exists with the same value.
func validateEnv() (err error) {
	defer func() {
		if err != nil {
			_, err = utils.MultiErrorList(err)
			err = fmt.Errorf("invalid environment: %w", err)
		}
	}()
	for _, kv := range strings.Split(emptyStringsEnv, "\n") {
		if strings.TrimSpace(kv) == "" {
			continue
		}
		i := strings.Index(kv, "=")
		if i == -1 {
			return errors.Errorf("malformed .env file line: %s", kv)
		}
		k := kv[:i]
		if k == "LOG_LEVEL" {
			continue // exceptional case of permitting legacy env w/o equivalent v2
		}
		v, ok := os.LookupEnv(k)
		if ok {
			if k2, ok2 := legacyEnvToV2[k]; ok2 {
				if v2 := os.Getenv(k2); v != v2 {
					err = multierr.Append(err, fmt.Errorf("environment variables %s and %s must be equal, or %s must not be set", k, k2, k2))
				}
			} else {
				err = multierr.Append(err, fmt.Errorf("environment variable %s must not be set: %v", k, v2.ErrUnsupported))
			}
		}
	}
	return
}

func (g *generalConfig) LogConfiguration(log coreconfig.LogFn) {
	log("Secrets:\n", g.secretsTOML)
	log("Input Configuration:\n", g.inputTOML)
	log("Effective Configuration, with defaults applied:\n", g.effectiveTOML)
}

// ConfigTOML implements chainlink.ConfigV2
func (g *generalConfig) ConfigTOML() (user, effective string) {
	return g.inputTOML, g.effectiveTOML
}

func (g *generalConfig) Dev() bool {
	return g.c.DevMode
}

func (g *generalConfig) FeatureExternalInitiators() bool {
	return *g.c.JobPipeline.ExternalInitiatorsEnabled
}

func (g *generalConfig) FeatureFeedsManager() bool {
	return *g.c.Feature.FeedsManager
}

func (g *generalConfig) FeatureOffchainReporting() bool {
	return *g.c.OCR.Enabled
}

func (g *generalConfig) FeatureOffchainReporting2() bool {
	return *g.c.OCR2.Enabled
}

func (g *generalConfig) FeatureLogPoller() bool {
	return *g.c.Feature.LogPoller
}

func (g *generalConfig) FeatureUICSAKeys() bool {
	return *g.c.Feature.UICSAKeys
}

func (g *generalConfig) AutoPprofEnabled() bool {
	return *g.c.AutoPprof.Enabled
}

func (g *generalConfig) EVMEnabled() bool {
	for _, c := range g.c.EVM {
		if c.IsEnabled() {
			return true
		}
	}
	return false
}

func (g *generalConfig) EVMRPCEnabled() bool {
	for _, c := range g.c.EVM {
		if c.IsEnabled() {
			if len(c.Nodes) > 0 {
				return true
			}
		}
	}
	return false
}

func (g *generalConfig) DefaultChainID() *big.Int {
	for _, c := range g.c.EVM {
		if c.IsEnabled() {
			return (*big.Int)(c.ChainID)
		}
	}
	return nil
}

func (g *generalConfig) EthereumHTTPURL() *url.URL {
	for _, c := range g.c.EVM {
		if c.IsEnabled() {
			for _, n := range c.Nodes {
				if n.SendOnly == nil || !*n.SendOnly {
					return (*url.URL)(n.HTTPURL)
				}
			}
		}
	}
	return nil

}
func (g *generalConfig) EthereumSecondaryURLs() (us []url.URL) {
	for _, c := range g.c.EVM {
		if c.IsEnabled() {
			for _, n := range c.Nodes {
				if n.HTTPURL != nil {
					us = append(us, (url.URL)(*n.HTTPURL))
				}
			}
		}
	}
	return nil

}
func (g *generalConfig) EthereumURL() string {
	for _, c := range g.c.EVM {
		if c.IsEnabled() {
			for _, n := range c.Nodes {
				if n.SendOnly == nil || !*n.SendOnly {
					if n.WSURL != nil {
						return n.WSURL.String()
					}
				}
			}
		}
	}
	return ""
}

func (g *generalConfig) P2PEnabled() bool {
	p := g.c.P2P
	return *p.V1.Enabled || *p.V2.Enabled
}

func (g *generalConfig) SolanaEnabled() bool {
	for _, c := range g.c.Solana {
		if c.IsEnabled() {
			return true
		}
	}
	return false
}

func (g *generalConfig) StarkNetEnabled() bool {
	for _, c := range g.c.Starknet {
		if c.IsEnabled() {
			return true
		}
	}
	return false
}

func (g *generalConfig) AllowOrigins() string {
	return *g.c.WebServer.AllowOrigins
}

func (g *generalConfig) AuditLoggerEnabled() bool {
	return *g.c.AuditLogger.Enabled
}

func (g *generalConfig) AuditLoggerForwardToUrl() (models.URL, error) {
	return *g.c.AuditLogger.ForwardToUrl, nil
}

func (g *generalConfig) AuditLoggerHeaders() (audit.ServiceHeaders, error) {
	return *g.c.AuditLogger.Headers, nil
}

func (g *generalConfig) AuditLoggerEnvironment() string {
	if g.Dev() {
		return "develop"
	}
	return "production"
}

func (g *generalConfig) AuditLoggerJsonWrapperKey() string {
	return *g.c.AuditLogger.JsonWrapperKey
}

func (g *generalConfig) AuthenticatedRateLimit() int64 {
	return *g.c.WebServer.RateLimit.Authenticated
}

func (g *generalConfig) AuthenticatedRateLimitPeriod() models.Duration {
	return *g.c.WebServer.RateLimit.AuthenticatedPeriod
}

func (g *generalConfig) AutoPprofBlockProfileRate() int {
	return int(*g.c.AutoPprof.BlockProfileRate)
}

func (g *generalConfig) AutoPprofCPUProfileRate() int {
	return int(*g.c.AutoPprof.CPUProfileRate)
}

func (g *generalConfig) AutoPprofGatherDuration() models.Duration {
	return models.MustMakeDuration(g.c.AutoPprof.GatherDuration.Duration())
}

func (g *generalConfig) AutoPprofGatherTraceDuration() models.Duration {
	return models.MustMakeDuration(g.c.AutoPprof.GatherTraceDuration.Duration())
}

func (g *generalConfig) AutoPprofGoroutineThreshold() int {
	return int(*g.c.AutoPprof.GoroutineThreshold)
}

func (g *generalConfig) AutoPprofMaxProfileSize() utils.FileSize {
	return *g.c.AutoPprof.MaxProfileSize
}

func (g *generalConfig) AutoPprofMemProfileRate() int {
	return int(*g.c.AutoPprof.MemProfileRate)
}

func (g *generalConfig) AutoPprofMemThreshold() utils.FileSize {
	return *g.c.AutoPprof.MemThreshold
}

func (g *generalConfig) AutoPprofMutexProfileFraction() int {
	return int(*g.c.AutoPprof.MutexProfileFraction)
}

func (g *generalConfig) AutoPprofPollInterval() models.Duration {
	return *g.c.AutoPprof.PollInterval
}

func (g *generalConfig) AutoPprofProfileRoot() string {
	s := *g.c.AutoPprof.ProfileRoot
	if s == "" {
		s = filepath.Join(g.RootDir(), "pprof")
	}
	return s
}

func (g *generalConfig) BlockBackfillDepth() uint64 { panic(v2.ErrUnsupported) }

func (g *generalConfig) BlockBackfillSkip() bool { panic(v2.ErrUnsupported) }

func (g *generalConfig) BridgeResponseURL() *url.URL {
	if g.c.WebServer.BridgeResponseURL.IsZero() {
		return nil
	}
	return g.c.WebServer.BridgeResponseURL.URL()
}

func (g *generalConfig) BridgeCacheTTL() time.Duration {
	return g.c.WebServer.BridgeCacheTTL.Duration()
}

func (g *generalConfig) CertFile() string {
	s := *g.c.WebServer.TLS.CertPath
	if s == "" {
		s = filepath.Join(g.TLSDir(), "server.crt")
	}
	return s
}

func (g *generalConfig) DatabaseBackupDir() string {
	return *g.c.Database.Backup.Dir
}

func (g *generalConfig) DatabaseBackupFrequency() time.Duration {
	return g.c.Database.Backup.Frequency.Duration()
}

func (g *generalConfig) DatabaseBackupMode() coreconfig.DatabaseBackupMode {
	return *g.c.Database.Backup.Mode
}

func (g *generalConfig) DatabaseBackupOnVersionUpgrade() bool {
	return *g.c.Database.Backup.OnVersionUpgrade
}

func (g *generalConfig) DatabaseListenerMaxReconnectDuration() time.Duration {
	return g.c.Database.Listener.MaxReconnectDuration.Duration()
}

func (g *generalConfig) DatabaseListenerMinReconnectInterval() time.Duration {
	return g.c.Database.Listener.MinReconnectInterval.Duration()
}

func (g *generalConfig) MigrateDatabase() bool {
	return *g.c.Database.MigrateOnStartup
}

func (g *generalConfig) ORMMaxIdleConns() int {
	return int(*g.c.Database.MaxIdleConns)
}

func (g *generalConfig) ORMMaxOpenConns() int {
	return int(*g.c.Database.MaxOpenConns)
}

func (g *generalConfig) DatabaseDefaultLockTimeout() time.Duration {
	return g.c.Database.DefaultLockTimeout.Duration()
}

func (g *generalConfig) DatabaseDefaultQueryTimeout() time.Duration {
	return g.c.Database.DefaultQueryTimeout.Duration()
}

func (g *generalConfig) DatabaseDefaultIdleInTxSessionTimeout() time.Duration {
	return g.c.Database.DefaultIdleInTxSessionTimeout.Duration()
}

func (g *generalConfig) DefaultHTTPLimit() int64 {
	return int64(*g.c.JobPipeline.HTTPRequest.MaxSize)
}

func (g *generalConfig) DefaultHTTPTimeout() models.Duration {
	return *g.c.JobPipeline.HTTPRequest.DefaultTimeout
}

func (g *generalConfig) ShutdownGracePeriod() time.Duration {
	return g.c.ShutdownGracePeriod.Duration()
}

func (g *generalConfig) ExplorerURL() *url.URL {
	u := (*url.URL)(g.c.ExplorerURL)
	if *u == zeroURL {
		u = nil
	}
	return u
}

func (g *generalConfig) FMDefaultTransactionQueueDepth() uint32 {
	return *g.c.FluxMonitor.DefaultTransactionQueueDepth
}

func (g *generalConfig) FMSimulateTransactions() bool {
	return *g.c.FluxMonitor.SimulateTransactions
}

func (g *generalConfig) GetDatabaseDialectConfiguredOrDefault() dialects.DialectName {
	return g.c.Database.Dialect
}

func (g *generalConfig) HTTPServerWriteTimeout() time.Duration {
	return g.c.WebServer.HTTPWriteTimeout.Duration()
}

func (g *generalConfig) InsecureFastScrypt() bool {
	return *g.c.InsecureFastScrypt
}

func (g *generalConfig) JSONConsole() bool {
	return *g.c.Log.JSONConsole
}

func (g *generalConfig) JobPipelineMaxRunDuration() time.Duration {
	return g.c.JobPipeline.MaxRunDuration.Duration()
}

func (g *generalConfig) JobPipelineMaxSuccessfulRuns() uint64 {
	return *g.c.JobPipeline.MaxSuccessfulRuns
}

func (g *generalConfig) JobPipelineReaperInterval() time.Duration {
	return g.c.JobPipeline.ReaperInterval.Duration()
}

func (g *generalConfig) JobPipelineReaperThreshold() time.Duration {
	return g.c.JobPipeline.ReaperThreshold.Duration()
}

func (g *generalConfig) JobPipelineResultWriteQueueDepth() uint64 {
	return uint64(*g.c.JobPipeline.ResultWriteQueueDepth)
}

func (g *generalConfig) KeeperDefaultTransactionQueueDepth() uint32 {
	return *g.c.Keeper.DefaultTransactionQueueDepth
}

func (g *generalConfig) KeeperGasPriceBufferPercent() uint16 {
	return *g.c.Keeper.GasPriceBufferPercent
}

func (g *generalConfig) KeeperGasTipCapBufferPercent() uint16 {
	return *g.c.Keeper.GasTipCapBufferPercent
}

func (g *generalConfig) KeeperBaseFeeBufferPercent() uint16 {
	return *g.c.Keeper.BaseFeeBufferPercent
}

func (g *generalConfig) KeeperMaximumGracePeriod() int64 {
	return *g.c.Keeper.MaxGracePeriod
}

func (g *generalConfig) KeeperRegistryCheckGasOverhead() uint32 {
	return *g.c.Keeper.Registry.CheckGasOverhead
}

func (g *generalConfig) KeeperRegistryPerformGasOverhead() uint32 {
	return *g.c.Keeper.Registry.PerformGasOverhead
}

func (g *generalConfig) KeeperRegistryMaxPerformDataSize() uint32 {
	return *g.c.Keeper.Registry.MaxPerformDataSize
}

func (g *generalConfig) KeeperRegistrySyncInterval() time.Duration {
	return g.c.Keeper.Registry.SyncInterval.Duration()
}

func (g *generalConfig) KeeperRegistrySyncUpkeepQueueSize() uint32 {
	return *g.c.Keeper.Registry.SyncUpkeepQueueSize
}

func (g *generalConfig) KeeperTurnLookBack() int64 {
	return *g.c.Keeper.TurnLookBack
}

func (g *generalConfig) KeyFile() string {
	if g.TLSKeyPath() == "" {
		return filepath.Join(g.TLSDir(), "server.key")
	}
	return g.TLSKeyPath()
}

func (g *generalConfig) DatabaseLockingMode() string { return g.c.Database.LockingMode() }

func (g *generalConfig) LeaseLockDuration() time.Duration {
	return g.c.Database.Lock.LeaseDuration.Duration()
}

func (g *generalConfig) LeaseLockRefreshInterval() time.Duration {
	return g.c.Database.Lock.LeaseRefreshInterval.Duration()
}

func (g *generalConfig) LogFileDir() string {
	s := *g.c.Log.File.Dir
	if s == "" {
		s = g.RootDir()
	}
	return s
}

func (g *generalConfig) LogFileMaxSize() utils.FileSize {
	return *g.c.Log.File.MaxSize
}

func (g *generalConfig) LogFileMaxAge() int64 {
	return *g.c.Log.File.MaxAgeDays
}

func (g *generalConfig) LogFileMaxBackups() int64 {
	return *g.c.Log.File.MaxBackups
}

func (g *generalConfig) LogUnixTimestamps() bool {
	return *g.c.Log.UnixTS
}

func (g *generalConfig) OCRBlockchainTimeout() time.Duration {
	return g.c.OCR.BlockchainTimeout.Duration()
}

func (g *generalConfig) OCRContractPollInterval() time.Duration {
	return g.c.OCR.ContractPollInterval.Duration()
}

func (g *generalConfig) OCRContractSubscribeInterval() time.Duration {
	return g.c.OCR.ContractSubscribeInterval.Duration()
}

func (g *generalConfig) OCRKeyBundleID() (string, error) {
	b := g.c.OCR.KeyBundleID
	if *b == zeroSha256Hash {
		return "", nil
	}
	return b.String(), nil
}

func (g *generalConfig) OCRObservationTimeout() time.Duration {
	return g.c.OCR.ObservationTimeout.Duration()
}

func (g *generalConfig) OCRSimulateTransactions() bool {
	return *g.c.OCR.SimulateTransactions
}

func (g *generalConfig) OCRTransmitterAddress() (ethkey.EIP55Address, error) {
	a := *g.c.OCR.TransmitterAddress
	if a.IsZero() {
		return a, errors.Wrap(coreconfig.ErrEnvUnset, "OCRTransmitterAddress is not set")
	}
	return a, nil
}

func (g *generalConfig) OCRTraceLogging() bool {
	return *g.c.P2P.TraceLogging
}

func (g *generalConfig) OCRDefaultTransactionQueueDepth() uint32 {
	return *g.c.OCR.DefaultTransactionQueueDepth
}

func (g *generalConfig) OCR2ContractConfirmations() uint16 {
	return uint16(*g.c.OCR2.ContractConfirmations)
}

func (g *generalConfig) OCR2ContractTransmitterTransmitTimeout() time.Duration {
	return g.c.OCR2.ContractTransmitterTransmitTimeout.Duration()
}

func (g *generalConfig) OCR2BlockchainTimeout() time.Duration {
	return g.c.OCR2.BlockchainTimeout.Duration()
}

func (g *generalConfig) OCR2DatabaseTimeout() time.Duration {
	return g.c.OCR2.DatabaseTimeout.Duration()
}

func (g *generalConfig) OCR2ContractPollInterval() time.Duration {
	return g.c.OCR2.ContractPollInterval.Duration()
}

func (g *generalConfig) OCR2ContractSubscribeInterval() time.Duration {
	return g.c.OCR2.ContractSubscribeInterval.Duration()
}

func (g *generalConfig) OCR2KeyBundleID() (string, error) {
	b := g.c.OCR2.KeyBundleID
	if *b == zeroSha256Hash {
		return "", nil
	}
	return b.String(), nil
}

func (g *generalConfig) OCR2TraceLogging() bool {
	return *g.c.P2P.TraceLogging
}

func (g *generalConfig) P2PNetworkingStack() (n ocrnetworking.NetworkingStack) {
	return g.c.P2P.NetworkStack()
}

func (g *generalConfig) P2PNetworkingStackRaw() string {
	return g.c.P2P.NetworkStack().String()
}

func (g *generalConfig) P2PPeerID() p2pkey.PeerID {
	return *g.c.P2P.PeerID
}

func (g *generalConfig) P2PPeerIDRaw() string {
	return g.c.P2P.PeerID.String()
}

func (g *generalConfig) P2PIncomingMessageBufferSize() int {
	return int(*g.c.P2P.IncomingMessageBufferSize)
}

func (g *generalConfig) P2POutgoingMessageBufferSize() int {
	return int(*g.c.P2P.OutgoingMessageBufferSize)
}

func (g *generalConfig) P2PAnnounceIP() net.IP {
	return *g.c.P2P.V1.AnnounceIP
}

func (g *generalConfig) P2PAnnouncePort() uint16 {
	return *g.c.P2P.V1.AnnouncePort
}

func (g *generalConfig) P2PBootstrapPeers() ([]string, error) {
	p := *g.c.P2P.V1.DefaultBootstrapPeers
	if p == nil {
		p = []string{}
	}
	return p, nil
}

func (g *generalConfig) P2PDHTAnnouncementCounterUserPrefix() uint32 {
	return *g.c.P2P.V1.DHTAnnouncementCounterUserPrefix
}

func (g *generalConfig) P2PListenIP() net.IP {
	return *g.c.P2P.V1.ListenIP
}

func (g *generalConfig) P2PListenPort() uint16 {
	v1 := g.c.P2P.V1
	p := *v1.ListenPort
	if p == 0 && *v1.Enabled {
		g.randomP2PPortOnce.Do(func() {
			addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
			if err != nil {
				panic(fmt.Errorf("unexpected ResolveTCPAddr error generating random P2PListenPort: %w", err))
			}
			l, err := net.ListenTCP("tcp", addr)
			if err != nil {
				panic(fmt.Errorf("unexpected ListenTCP error generating random P2PListenPort: %w", err))
			}
			defer l.Close()
			g.randomP2PPort = uint16(l.Addr().(*net.TCPAddr).Port)
			g.lggr.Warnw(fmt.Sprintf("P2PListenPort was not set, listening on random port %d. A new random port will be generated on every boot, for stability it is recommended to set P2PListenPort to a fixed value in your environment", g.randomP2PPort), "p2pPort", g.randomP2PPort)
		})
		return g.randomP2PPort
	}
	return p
}

func (g *generalConfig) P2PListenPortRaw() string {
	p := *g.c.P2P.V1.ListenPort
	if p == 0 {
		return ""
	}
	return strconv.Itoa(int(p))
}

func (g *generalConfig) P2PNewStreamTimeout() time.Duration {
	return g.c.P2P.V1.NewStreamTimeout.Duration()
}

func (g *generalConfig) P2PBootstrapCheckInterval() time.Duration {
	return g.c.P2P.V1.BootstrapCheckInterval.Duration()
}

func (g *generalConfig) P2PDHTLookupInterval() int {
	return int(*g.c.P2P.V1.DHTLookupInterval)
}

func (g *generalConfig) P2PPeerstoreWriteInterval() time.Duration {
	return g.c.P2P.V1.PeerstoreWriteInterval.Duration()
}

func (g *generalConfig) P2PV2AnnounceAddresses() []string {
	if v := g.c.P2P.V2.AnnounceAddresses; v != nil {
		return *v
	}
	return nil
}

func (g *generalConfig) P2PV2Bootstrappers() (locators []commontypes.BootstrapperLocator) {
	if v := g.c.P2P.V2.DefaultBootstrappers; v != nil {
		return *v
	}
	return nil
}

func (g *generalConfig) P2PV2BootstrappersRaw() (s []string) {
	if v := g.c.P2P.V2.DefaultBootstrappers; v != nil {
		for _, b := range *v {
			t, err := b.MarshalText()
			if err != nil {
				// log panic matches old behavior - only called for UI presentation
				panic(fmt.Sprintf("Failed to marshal bootstrapper: %v", err))
			}
			s = append(s, string(t))
		}
	}
	return
}

func (g *generalConfig) P2PV2DeltaDial() models.Duration {
	if v := g.c.P2P.V2.DeltaDial; v != nil {
		return *v
	}
	return models.Duration{}
}

func (g *generalConfig) P2PV2DeltaReconcile() models.Duration {
	if v := g.c.P2P.V2.DeltaReconcile; v != nil {
		return *v

	}
	return models.Duration{}
}

func (g *generalConfig) P2PV2ListenAddresses() []string {
	if v := g.c.P2P.V2.ListenAddresses; v != nil {
		return *v
	}
	return nil
}

func (g *generalConfig) PyroscopeServerAddress() string {
	return *g.c.Pyroscope.ServerAddress
}

func (g *generalConfig) PyroscopeEnvironment() string {
	return *g.c.Pyroscope.Environment
}
func (g *generalConfig) Port() uint16 {
	return *g.c.WebServer.HTTPPort
}

func (g *generalConfig) RPID() string {
	return *g.c.WebServer.MFA.RPID
}

func (g *generalConfig) RPOrigin() string {
	return *g.c.WebServer.MFA.RPOrigin
}

func (g *generalConfig) ReaperExpiration() models.Duration {
	return *g.c.WebServer.SessionReaperExpiration
}

func (g *generalConfig) RootDir() string {
	d := *g.c.RootDir
	h, err := parse.HomeDir(d)
	if err != nil {
		g.lggr.Error("Failed to expand RootDir. You may need to set an explicit path", "err", err)
		return d
	}
	return h
}

func (g *generalConfig) SecureCookies() bool {
	return *g.c.WebServer.SecureCookies
}

func (g *generalConfig) SessionOptions() sessions.Options {
	return sessions.Options{
		Secure:   g.SecureCookies(),
		HttpOnly: true,
		MaxAge:   86400 * 30,
		SameSite: http.SameSiteStrictMode,
	}
}

func (g *generalConfig) SessionTimeout() models.Duration {
	return models.MustMakeDuration(g.c.WebServer.SessionTimeout.Duration())
}

func (g *generalConfig) SentryDSN() string {
	return *g.c.Sentry.DSN
}

func (g *generalConfig) SentryDebug() bool {
	return *g.c.Sentry.Debug
}

func (g *generalConfig) SentryEnvironment() string {
	return *g.c.Sentry.Environment
}

func (g *generalConfig) SentryRelease() string {
	return *g.c.Sentry.Release
}

func (g *generalConfig) TLSCertPath() string {
	return *g.c.WebServer.TLS.CertPath
}

func (g *generalConfig) TLSDir() string {
	return filepath.Join(g.RootDir(), "tls")
}

func (g *generalConfig) TLSHost() string {
	return *g.c.WebServer.TLS.Host
}

func (g *generalConfig) TLSKeyPath() string {
	return *g.c.WebServer.TLS.KeyPath
}

func (g *generalConfig) TLSPort() uint16 {
	return *g.c.WebServer.TLS.HTTPSPort
}

func (g *generalConfig) TLSRedirect() bool {
	return *g.c.WebServer.TLS.ForceRedirect
}

func (g *generalConfig) TelemetryIngressLogging() bool {
	return *g.c.TelemetryIngress.Logging
}

func (g *generalConfig) TelemetryIngressUniConn() bool {
	return *g.c.TelemetryIngress.UniConn
}

func (g *generalConfig) TelemetryIngressServerPubKey() string {
	return *g.c.TelemetryIngress.ServerPubKey
}

func (g *generalConfig) TelemetryIngressURL() *url.URL {
	if g.c.TelemetryIngress.URL.IsZero() {
		return nil
	}
	return g.c.TelemetryIngress.URL.URL()
}

func (g *generalConfig) TelemetryIngressBufferSize() uint {
	return uint(*g.c.TelemetryIngress.BufferSize)
}

func (g *generalConfig) TelemetryIngressMaxBatchSize() uint {
	return uint(*g.c.TelemetryIngress.MaxBatchSize)
}

func (g *generalConfig) TelemetryIngressSendInterval() time.Duration {
	return g.c.TelemetryIngress.SendInterval.Duration()
}

func (g *generalConfig) TelemetryIngressSendTimeout() time.Duration {
	return g.c.TelemetryIngress.SendTimeout.Duration()
}

func (g *generalConfig) TelemetryIngressUseBatchSend() bool {
	return *g.c.TelemetryIngress.UseBatchSend
}

func (g *generalConfig) TriggerFallbackDBPollInterval() time.Duration {
	return g.c.Database.Listener.FallbackPollInterval.Duration()
}

func (g *generalConfig) UnAuthenticatedRateLimit() int64 {
	return *g.c.WebServer.RateLimit.Unauthenticated
}

func (g *generalConfig) UnAuthenticatedRateLimitPeriod() models.Duration {
	return *g.c.WebServer.RateLimit.UnauthenticatedPeriod
}

var (
	zeroURL        = url.URL{}
	zeroSha256Hash = models.Sha256Hash{}
)
