package resolver

import "github.com/smartcontractkit/chainlink/core/config"

type ConfigResolver struct {
	cfg config.GeneralConfig
}

func NewConfig(cfg config.GeneralConfig) *ConfigResolver {
	return &ConfigResolver{cfg: cfg}
}

func (r *ConfigResolver) AllowOrigins() string {
	return r.cfg.AllowOrigins()
}

func (r *ConfigResolver) BlockBackfillDepth() int32 {
	return int32(r.cfg.BlockBackfillDepth())
}

func (r *ConfigResolver) BridgeResponseURL() string {
	return r.cfg.BridgeResponseURL().String()
}

func (r *ConfigResolver) ClientNodeURL() string {
	return r.cfg.ClientNodeURL()
}

func (r *ConfigResolver) DatabaseBackupFrequency() string {
	return r.cfg.DatabaseBackupFrequency().String()
}

func (r *ConfigResolver) DatabaseBackupMode() string {
	return string(r.cfg.DatabaseBackupMode())
}

func (r *ConfigResolver) DatabaseMaximumTxDuration() string {
	return r.cfg.DatabaseMaximumTxDuration().String()
}

func (r *ConfigResolver) DatabaseTimeout() string {
	return r.cfg.DatabaseTimeout().String()
}

func (r *ConfigResolver) DatabaseLockingMode() string {
	return r.cfg.DatabaseLockingMode()
}

func (r *ConfigResolver) DefaultChainID() string {
	return r.cfg.DefaultChainID().String()
}

func (r *ConfigResolver) DefaultHTTPLimit() int32 {
	return int32(r.cfg.DefaultHTTPLimit())
}

func (r *ConfigResolver) DefaultHTTPTimeout() string {
	return r.cfg.DefaultHTTPTimeout().String()
}

func (r *ConfigResolver) Dev() bool {
	return r.cfg.Dev()
}

func (r *ConfigResolver) EthereumDisabled() bool {
	return r.cfg.EthereumDisabled()
}

func (r *ConfigResolver) EthereumHTTPURL() string {
	if r.cfg.EthereumHTTPURL() != nil {
		url := r.cfg.EthereumHTTPURL().String()

		return url
	}

	return ""
}

func (r *ConfigResolver) EthereumSecondaryURLs() []string {
	var urls []string

	for _, url := range r.cfg.EthereumSecondaryURLs() {
		urls = append(urls, url.String())
	}

	return urls
}

func (r *ConfigResolver) EthereumURL() string {
	return r.cfg.EthereumURL()
}

func (r *ConfigResolver) ExplorerURL() string {
	if r.cfg.ExplorerURL() != nil {
		url := r.cfg.ExplorerURL().String()

		return url
	}

	return ""
}

func (r *ConfigResolver) FMDefaultTransactionQueueDepth() int32 {
	return int32(r.cfg.FMDefaultTransactionQueueDepth())
}

func (r *ConfigResolver) FeatureExternalInitiators() bool {
	return r.cfg.FeatureExternalInitiators()
}

func (r *ConfigResolver) FeatureOffchainReporting() bool {
	return r.cfg.FeatureOffchainReporting()
}

func (r *ConfigResolver) InsecureFastScrypt() bool {
	return r.cfg.InsecureFastScrypt()
}

func (r *ConfigResolver) JSONConsole() bool {
	return r.cfg.JSONConsole()
}

func (r *ConfigResolver) JobPipelineReaperInterval() string {
	return r.cfg.JobPipelineReaperInterval().String()
}

func (r *ConfigResolver) JobPipelineReaperThreshold() string {
	return r.cfg.JobPipelineReaperThreshold().String()
}

func (r *ConfigResolver) KeeperDefaultTransactionQueueDepth() int32 {
	return int32(r.cfg.KeeperDefaultTransactionQueueDepth())
}

func (r *ConfigResolver) KeeperGasPriceBufferPercent() int32 {
	return int32(r.cfg.KeeperGasPriceBufferPercent())
}

func (r *ConfigResolver) KeeperGasTipCapBufferPercent() int32 {
	return int32(r.cfg.KeeperGasTipCapBufferPercent())
}

func (r *ConfigResolver) LogLevel() LogLevel {
	return ToLogLevel(r.cfg.LogLevel().String())
}

func (r *ConfigResolver) LogSQLMigrations() bool {
	return r.cfg.LogSQLMigrations()
}

func (r *ConfigResolver) LogSQL() bool {
	return r.cfg.LogSQL()
}

func (r *ConfigResolver) LogToDisk() bool {
	return r.cfg.LogToDisk()
}

func (r *ConfigResolver) OCRBootstrapCheckInterval() string {
	return r.cfg.OCRBootstrapCheckInterval().String()
}

func (r *ConfigResolver) OCRContractTransmitterTransmitTimeout() string {
	return r.cfg.OCRContractTransmitterTransmitTimeout().String()
}

func (r *ConfigResolver) OCRDHTLookupInterval() int32 {
	return int32(r.cfg.OCRDHTLookupInterval())
}

func (r *ConfigResolver) OCRDatabaseTimeout() string {
	return r.cfg.OCRDatabaseTimeout().String()
}

func (r *ConfigResolver) OCRDefaultTransactionQueueDepth() int32 {
	return int32(r.cfg.OCRDefaultTransactionQueueDepth())
}

func (r *ConfigResolver) OCRIncomingMessageBufferSize() int32 {
	return int32(r.cfg.OCRIncomingMessageBufferSize())
}

func (r *ConfigResolver) OCRNewStreamTimeout() string {
	return r.cfg.OCRNewStreamTimeout().String()
}

func (r *ConfigResolver) OCROutgoingMessageBufferSize() int32 {
	return int32(r.cfg.OCROutgoingMessageBufferSize())
}

func (r *ConfigResolver) OCRTraceLogging() bool {
	return r.cfg.OCRTraceLogging()
}

func (r *ConfigResolver) P2PBootstrapPeers() []string {
	peers, err := r.cfg.P2PBootstrapPeers()
	if err != nil {
		return []string{}
	}

	return peers
}

func (r *ConfigResolver) P2PListenIP() string {
	return r.cfg.P2PListenIP().String()
}

func (r *ConfigResolver) P2PListenPort() string {
	return r.cfg.P2PListenPortRaw()
}

func (r *ConfigResolver) P2PNetworkingStack() string {
	return r.cfg.P2PNetworkingStackRaw()
}

func (r *ConfigResolver) P2PPeerID() string {
	return r.cfg.P2PPeerIDRaw()
}

func (r *ConfigResolver) P2PV2AnnounceAddresses() []string {
	return r.cfg.P2PV2AnnounceAddressesRaw()
}

func (r *ConfigResolver) P2PV2Bootstrappers() []string {
	return r.cfg.P2PV2BootstrappersRaw()
}

func (r *ConfigResolver) P2PV2DeltaDial() string {
	return r.cfg.P2PV2DeltaDial().String()
}

func (r *ConfigResolver) P2PV2DeltaReconcile() string {
	return r.cfg.P2PV2DeltaReconcile().String()
}

func (r *ConfigResolver) P2PV2ListenAddresses() []string {
	return r.cfg.P2PV2ListenAddresses()
}

func (r *ConfigResolver) Port() int32 {
	return int32(r.cfg.Port())
}

func (r *ConfigResolver) ReaperExpiration() string {
	return r.cfg.ReaperExpiration().String()
}

func (r *ConfigResolver) ReplayFromBlock() int32 {
	return int32(r.cfg.ReplayFromBlock())
}

func (r *ConfigResolver) RootDir() string {
	return r.cfg.RootDir()
}

func (r *ConfigResolver) SecureCookies() bool {
	return r.cfg.SecureCookies()
}

func (r *ConfigResolver) SessionTimeout() string {
	return r.cfg.SessionTimeout().String()
}

func (r *ConfigResolver) TLSHost() string {
	return r.cfg.TLSHost()
}

func (r *ConfigResolver) TLSPort() int32 {
	return int32(r.cfg.TLSPort())
}

func (r *ConfigResolver) TLSRedirect() bool {
	return r.cfg.TLSRedirect()
}

func (r *ConfigResolver) TelemetryIngressLogging() bool {
	return r.cfg.TelemetryIngressLogging()
}

func (r *ConfigResolver) TelemetryIngressServerPubKey() string {
	return r.cfg.TelemetryIngressServerPubKey()
}

func (r *ConfigResolver) TelemetryIngressURL() string {
	if r.cfg.TelemetryIngressURL() != nil {
		url := r.cfg.TelemetryIngressURL().String()

		return url
	}

	return ""
}

func (r *ConfigResolver) TriggerFallbackDBPollInterval() string {
	return r.cfg.TriggerFallbackDBPollInterval().String()
}

type ConfigPayloadResolver struct {
	cfg config.GeneralConfig
}

func NewConfigPayload(cfg config.GeneralConfig) *ConfigPayloadResolver {
	return &ConfigPayloadResolver{cfg: cfg}
}

func (r *ConfigPayloadResolver) ToConfig() (*ConfigResolver, bool) {
	return NewConfig(r.cfg), true
}
