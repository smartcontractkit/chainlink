package chainlink

import (
	"fmt"
	"net"
	"net/url"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/pelletier/go-toml/v2"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/libocr/commontypes"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"

	coreconfig "github.com/smartcontractkit/chainlink/core/config"
	v2 "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// generalConfig is a wrapper to adapt Config to the config.GeneralConfig interface.
type generalConfig struct {
	input     string // user input, normalized via de/re-serialization
	effective string // with default values included
	c         *Config

	// state
	appID     uuid.UUID
	appIDOnce sync.Once

	logLevelDefault zapcore.Level
	logLevel        zapcore.Level
	logSQL          bool
	logMu           sync.RWMutex
}

func NewGeneralConfig(tomlString string, lggr logger.Logger) (coreconfig.GeneralConfig, error) {
	lggr = lggr.Named("Config")
	var c Config
	err := toml.Unmarshal([]byte(tomlString), &c)
	if err != nil {
		return nil, err
	}
	input, err := c.TOMLString()
	if err != nil {
		return nil, err
	}

	c.SetDefaults()

	effective, err := c.TOMLString()
	if err != nil {
		return nil, err
	}
	return &generalConfig{c: &c, input: input, effective: effective}, nil
}

func (g *generalConfig) Validate() error {
	return g.c.Validate()
}

func (g *generalConfig) LogConfiguration(log func(...any)) {
	log("Input Configuration:\n", g.input)
	log("Effective Configuration, with defaults applied:\n", g.effective)
}

func (g *generalConfig) Dev() bool {
	return v2.CLDev
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
	return *g.c.Feature.UICSA
}

func (g *generalConfig) AutoPprofEnabled() bool {
	return *g.c.AutoPprof.Enabled
}

func (g *generalConfig) EVMEnabled() bool {
	for _, c := range g.c.EVM {
		if e := c.Enabled; e != nil && *e == true {
			return true
		}
	}
	return false
}

func (g *generalConfig) KeeperCheckUpkeepGasPriceFeatureEnabled() bool {
	return *g.c.Keeper.UpkeepCheckGasPriceEnabled
}

func (g *generalConfig) P2PEnabled() bool {
	p := g.c.P2P
	return p.V1 != nil || p.V2 != nil //TODO or Disabled off switch?
}

func (g *generalConfig) SolanaEnabled() bool {
	for _, c := range g.c.Solana {
		if e := c.Enabled; e != nil && *e == true {
			return true
		}
	}
	return false
}

func (g *generalConfig) TerraEnabled() bool {
	for _, c := range g.c.Terra {
		if e := c.Enabled; e != nil && *e == true {
			return true
		}
	}
	return false
}

func (g *generalConfig) StarkNetEnabled() bool {
	//TODO implement me
	panic("implement me")
}

func (g *generalConfig) AllowOrigins() string {
	return *g.c.WebServer.AllowOrigins
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
	return *g.c.AutoPprof.ProfileRoot
}

func (g *generalConfig) BlockBackfillDepth() uint64 {
	//TODO implement me
	panic("implement me")
}

func (g *generalConfig) BlockBackfillSkip() bool {
	//TODO implement me
	panic("implement me")
}

func (g *generalConfig) BridgeResponseURL() *url.URL {
	return (*url.URL)(g.c.WebServer.BridgeResponseURL)
}

func (g *generalConfig) CertFile() string {
	return *g.c.WebServer.TLS.CertPath
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

func (g *generalConfig) DatabaseBackupURL() *url.URL {
	return (*url.URL)(g.c.Database.Backup.URL)
}

func (g *generalConfig) DatabaseListenerMaxReconnectDuration() time.Duration {
	return g.c.Database.Listener.MaxReconnectDuration.Duration()
}

func (g *generalConfig) DatabaseListenerMinReconnectInterval() time.Duration {
	return g.c.Database.Listener.MinReconnectInterval.Duration()
}

func (g *generalConfig) DefaultHTTPLimit() int64 {
	return int64(*g.c.JobPipeline.HTTPRequestMaxSize)
}

func (g *generalConfig) DefaultHTTPTimeout() models.Duration {
	return *g.c.JobPipeline.DefaultHTTPRequestTimeout
}

func (g *generalConfig) ShutdownGracePeriod() time.Duration {
	return g.c.ShutdownGracePeriod.Duration()
}

func (g *generalConfig) ExplorerURL() *url.URL {
	return (*url.URL)(g.c.ExplorerURL)
}

func (g *generalConfig) FMDefaultTransactionQueueDepth() uint32 {
	return *g.c.FluxMonitor.DefaultTransactionQueueDepth
}

func (g *generalConfig) FMSimulateTransactions() bool {
	return *g.c.FluxMonitor.SimulateTransactions
}

func (g *generalConfig) GetDatabaseDialectConfiguredOrDefault() dialects.DialectName {
	//TODO implement me
	panic("implement me")
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

func (g *generalConfig) KeeperGasPriceBufferPercent() uint32 {
	return *g.c.Keeper.GasPriceBufferPercent
}

func (g *generalConfig) KeeperGasTipCapBufferPercent() uint32 {
	return *g.c.Keeper.GasTipCapBufferPercent
}

func (g *generalConfig) KeeperBaseFeeBufferPercent() uint32 {
	return *g.c.Keeper.BaseFeeBufferPercent
}

func (g *generalConfig) KeeperMaximumGracePeriod() int64 {
	return *g.c.Keeper.MaximumGracePeriod
}

func (g *generalConfig) KeeperRegistryCheckGasOverhead() uint64 {
	return g.c.Keeper.RegistryCheckGasOverhead.ToInt().Uint64()
}

func (g *generalConfig) KeeperRegistryPerformGasOverhead() uint64 {
	return g.c.Keeper.RegistryPerformGasOverhead.ToInt().Uint64()
}

func (g *generalConfig) KeeperRegistrySyncInterval() time.Duration {
	return g.c.Keeper.RegistrySyncInterval.Duration()
}

func (g *generalConfig) KeeperRegistrySyncUpkeepQueueSize() uint32 {
	return *g.c.Keeper.RegistrySyncUpkeepQueueSize
}

func (g *generalConfig) KeeperTurnLookBack() int64 {
	return *g.c.Keeper.TurnLookBack
}

func (g *generalConfig) KeeperTurnFlagEnabled() bool {
	return *g.c.Keeper.TurnFlagEnabled
}

func (g *generalConfig) KeyFile() string {
	if g.TLSKeyPath() == "" {
		return filepath.Join(g.TLSDir(), "server.key")
	}
	return g.TLSKeyPath()
}

func (g *generalConfig) LeaseLockDuration() time.Duration {
	return g.c.Database.Lock.LeaseDuration.Duration()
}

func (g *generalConfig) LeaseLockRefreshInterval() time.Duration {
	return g.c.Database.Lock.LeaseRefreshInterval.Duration()
}

func (g *generalConfig) LogFileDir() string {
	return *g.c.Log.FileDir
}

func (g *generalConfig) LogFileMaxSize() utils.FileSize {
	return *g.c.Log.FileMaxSize
}

func (g *generalConfig) LogFileMaxAge() int64 {
	return *g.c.Log.FileMaxAgeDays
}

func (g *generalConfig) LogFileMaxBackups() int64 {
	return *g.c.Log.FileMaxBackups
}

func (g *generalConfig) LogUnixTimestamps() bool {
	return *g.c.Log.UnixTS
}

func (g *generalConfig) MigrateDatabase() bool {
	return *g.c.Database.MigrateOnStartup
}

func (g *generalConfig) ORMMaxIdleConns() int {
	return int(*g.c.Database.ORMMaxIdleConns)
}

func (g *generalConfig) ORMMaxOpenConns() int {
	return int(*g.c.Database.ORMMaxOpenConns)
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
	return *g.c.RootDir
}

func (g *generalConfig) SecureCookies() bool {
	return *g.c.WebServer.SecureCookies
}

func (g *generalConfig) SessionOptions() sessions.Options {
	return sessions.Options{
		Secure:   g.SecureCookies(),
		HttpOnly: true,
		MaxAge:   86400 * 30,
	}
}

func (g *generalConfig) SessionTimeout() models.Duration {
	return models.MustMakeDuration(g.c.WebServer.SessionTimeout.Duration())
}

func (g *generalConfig) TLSCertPath() string {
	return *g.c.WebServer.TLS.CertPath
}

func (g *generalConfig) TLSDir() string {
	return filepath.Join(*g.c.RootDir, "tls")
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
	return (*url.URL)(g.c.TelemetryIngress.URL)
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
	return g.c.OCR.KeyBundleID.String(), nil
}

func (g *generalConfig) OCRObservationTimeout() time.Duration {
	return g.c.OCR.ObservationTimeout.Duration()
}

func (g *generalConfig) OCRSimulateTransactions() bool {
	return *g.c.OCR.SimulateTransactions
}

func (g *generalConfig) OCRTransmitterAddress() (ethkey.EIP55Address, error) {
	return *g.c.OCR.TransmitterAddress, nil
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
	return g.c.OCR2.KeyBundleID.String(), nil
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
	return *g.c.P2P.V1.PeerID
}

func (g *generalConfig) P2PPeerIDRaw() string {
	return g.c.P2P.V1.PeerID.String()
}

func (g *generalConfig) P2PIncomingMessageBufferSize() int {
	return int(*g.c.P2P.IncomingMessageBufferSize)
}

func (g *generalConfig) P2POutgoingMessageBufferSize() int {
	return int(*g.c.P2P.OutgoingMessageBufferSize)
}

func (g *generalConfig) OCRNewStreamTimeout() time.Duration {
	return g.c.P2P.V1.NewStreamTimeout.Duration()
}

func (g *generalConfig) OCRBootstrapCheckInterval() time.Duration {
	return g.c.P2P.V1.BootstrapCheckInterval.Duration()
}

func (g *generalConfig) OCRDHTLookupInterval() int {
	return int(*g.c.P2P.V1.DHTLookupInterval)
}

func (g *generalConfig) OCRIncomingMessageBufferSize() int {
	return int(*g.c.P2P.IncomingMessageBufferSize)
}

func (g *generalConfig) OCROutgoingMessageBufferSize() int {
	return int(*g.c.P2P.OutgoingMessageBufferSize)
}

func (g *generalConfig) P2PAnnounceIP() net.IP {
	return *g.c.P2P.V1.AnnounceIP
}

func (g *generalConfig) P2PAnnouncePort() uint16 {
	return *g.c.P2P.V1.AnnouncePort
}

func (g *generalConfig) P2PBootstrapPeers() ([]string, error) {
	return *g.c.P2P.V1.DefaultBootstrapPeers, nil
}

func (g *generalConfig) P2PDHTAnnouncementCounterUserPrefix() uint32 {
	return *g.c.P2P.V1.DHTAnnouncementCounterUserPrefix
}

func (g *generalConfig) P2PListenIP() net.IP {
	return *g.c.P2P.V1.ListenIP
}

func (g *generalConfig) P2PListenPort() uint16 {
	return *g.c.P2P.V1.ListenPort
}

func (g *generalConfig) P2PListenPortRaw() string {
	return strconv.Itoa(int(*g.c.P2P.V1.ListenPort))
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
	if p := g.c.P2P; p != nil {
		if v2 := p.V2; v2 != nil {
			if v := v2.AnnounceAddresses; v != nil {
				return *v
			}
		}
	}
	return nil
}

func (g *generalConfig) P2PV2Bootstrappers() (locators []commontypes.BootstrapperLocator) {
	if p := g.c.P2P; p != nil {
		if v2 := p.V2; v2 != nil {
			if v := v2.DefaultBootstrappers; v != nil {
				return *v
			}
		}
	}
	return nil
}

func (g *generalConfig) P2PV2BootstrappersRaw() (s []string) {
	if p := g.c.P2P; p != nil {
		if v2 := p.V2; v2 != nil {
			if v := v2.DefaultBootstrappers; v != nil {
				for _, b := range *v {
					t, err := b.MarshalText()
					if err != nil {
						// log panic matches old behavior - only called for UI presentation
						panic(fmt.Sprintf("Failed to marshal bootstrapper: %v", err))
					}
					s = append(s, string(t))
				}
			}
		}
	}
	return
}

func (g *generalConfig) P2PV2DeltaDial() models.Duration {
	if p := g.c.P2P; p != nil {
		if v2 := p.V2; v2 != nil {
			if v := v2.DeltaDial; v != nil {
				return *v
			}
		}
	}
	return models.Duration{}
}

func (g *generalConfig) P2PV2DeltaReconcile() models.Duration {
	if p := g.c.P2P; p != nil {
		if v2 := p.V2; v2 != nil {
			if v := v2.DeltaReconcile; v != nil {
				return *v
			}
		}
	}
	return models.Duration{}
}

func (g *generalConfig) P2PV2ListenAddresses() []string {
	if p := g.c.P2P; p != nil {
		if v2 := p.V2; v2 != nil {
			if v := v2.ListenAddresses; v != nil {
				return *v
			}
		}
	}
	return nil
}
