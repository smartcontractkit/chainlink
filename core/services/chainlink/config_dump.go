package chainlink

import (
	"encoding/csv"
	"net"
	"net/url"
	"strings"

	"github.com/pelletier/go-toml/v2"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/config/envvar"
	"github.com/smartcontractkit/chainlink/core/config/parse"
	tcfg "github.com/smartcontractkit/chainlink/core/config/toml"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func (app *ChainlinkApplication) ConfigDump() (string, error) {
	var c Config

	if err := c.loadChainsAndNodes(app.Chains); err != nil {
		return "", err
	}

	c.loadLegacyEVMEnv()

	c.loadLegacyCoreEnv()

	b, err := toml.Marshal(c)
	if err != nil {
		return "", err
	}
	//TODO what about secrets?

	return string(b), nil
}

//TODO doc
func (c *Config) loadChainsAndNodes(chains Chains) error {
	//TODO copy chains and nodes from database
	//TODO must iterate over pure ORM methods for all chains and nodes
	//TODO even disabled chains?
	c.EVM = nil
	c.Solana = nil
	c.Terra = nil
	return nil
}

//TODO doc
func (c *Config) loadLegacyEVMEnv() {
	//TODO read legacy ETH/EVM global overrides, clobbering persisted values already set
	//TODO ignore EVM_NODES? since DB already updated at this point?
}

// loadLegacyCoreEnv loads CoreConfig values from legacy environment variables.
//TODO test
func (c *Config) loadLegacyCoreEnv() {
	c.Dev = envvar.NewBool("Dev").ParsePtr()
	c.ExplorerURL = envURL("ExplorerURL")
	c.InsecureFastScrypt = envvar.NewBool("InsecureFastScrypt").ParsePtr()
	c.ReaperExpiration = envDuration("ReaperExpiration")
	c.RootDir = envvar.RootDir.ParsePtr()
	c.ShutdownGracePeriod = envDuration("ShutdownGracePeriod")

	c.Database = &tcfg.DatabaseConfig{
		ListenerMaxReconnectDuration:  envDuration("DatabaseListenerMaxReconnectDuration"),
		ListenerMinReconnectInterval:  envDuration("DatabaseListenerMinReconnectInterval"),
		Migrate:                       envvar.NewBool("MigrateDatabase").ParsePtr(),
		ORMMaxIdleConns:               envvar.NewInt64("ORMMaxIdleConns").ParsePtr(),
		ORMMaxOpenConns:               envvar.NewInt64("ORMMaxOpenConns").ParsePtr(),
		TriggerFallbackDBPollInterval: envDuration("TriggerFallbackDBPollInterval"),
		AdvisoryLockCheckInterval:     envDuration("AdvisoryLockCheckInterval"),
		AdvisoryLockID:                envvar.AdvisoryLockID.ParsePtr(),
		LockingMode:                   envvar.NewString("DatabaseLockingMode").ParsePtr(),
		LeaseLockDuration:             envDuration("LeaseLockDuration"),
		LeaseLockRefreshInterval:      envDuration("LeaseLockRefreshInterval"),
		BackupDir:                     envvar.NewString("DatabaseBackupDir").ParsePtr(),
		BackupFrequency:               envDuration("DatabaseBackupFrequency"),
		BackupMode:                    config.DatabaseBackupModeEnvVar.ParsePtr(),
		BackupOnVersionUpgrade:        envvar.NewBool("DatabaseBackupOnVersionUpgrade").ParsePtr(),
		BackupURL:                     envURL("DatabaseBackupDir"),
	}

	c.TelemetryIngress = &tcfg.TelemetryIngressConfig{
		UniConn:      envvar.NewBool("TelemetryIngressUniConn").ParsePtr(),
		Logging:      envvar.NewBool("TelemetryIngressLogging").ParsePtr(),
		ServerPubKey: envvar.NewString("TelemetryIngressServerPubKey").ParsePtr(),
		URL:          envURL("TelemetryIngressURL"),
		BufferSize:   envvar.NewUint16("TelemetryIngressBufferSize").ParsePtr(),
		MaxBatchSize: envvar.NewUint16("TelemetryIngressMaxBatchSize").ParsePtr(),
		SendInterval: envDuration("TelemetryIngressSendInterval"),
		SendTimeout:  envDuration("TelemetryIngressSendTimeout"),
		UseBatchSend: envvar.NewBool("TelemetryIngressUseBatchSend").ParsePtr(),
	}

	c.Log = &tcfg.LogConfig{
		JSONConsole:    envvar.JSONConsole.ParsePtr(),
		FileDir:        envvar.NewString("LogFileDir").ParsePtr(),
		Level:          envvar.LogLevel.ParsePtr(),
		SQL:            envvar.NewBool("LogSQL").ParsePtr(),
		FileMaxSize:    envvar.LogFileMaxSize.ParsePtr(),
		FileMaxAgeDays: envvar.LogFileMaxAge.ParsePtr(),
		FileMaxBackups: envvar.LogFileMaxBackups.ParsePtr(),
		UnixTS:         envvar.LogUnixTS.ParsePtr(),
	}

	c.WebServer = &tcfg.WebServerConfig{
		AllowOrigins:                   envvar.NewString("AllowOrigins").ParsePtr(),
		AuthenticatedRateLimit:         envvar.NewInt64("AuthenticatedRateLimit").ParsePtr(),
		AuthenticatedRateLimitPeriod:   envDuration("AuthenticatedRateLimitPeriod"),
		BridgeResponseURL:              envURL("BridgeResponseURL"),
		HTTPWriteTimeout:               envDuration("HTTPServerWriteTimeout"),
		Port:                           envvar.NewUint16("Port").ParsePtr(),
		SecureCookies:                  envvar.NewBool("SecureCookies").ParsePtr(),
		SessionTimeout:                 envDuration("SessionTimeout"),
		UnAuthenticatedRateLimit:       envvar.NewInt64("UnAuthenticatedRateLimit").ParsePtr(),
		UnAuthenticatedRateLimitPeriod: envDuration("UnAuthenticatedRateLimitPeriod"),
		RPID:                           envvar.NewString("RPID").ParsePtr(),
		RPOrigin:                       envvar.NewString("RPOrigin").ParsePtr(),
		TLSCertPath:                    envvar.NewString("TLSCertPath").ParsePtr(),
		TLSHost:                        envvar.NewString("TLSHost").ParsePtr(),
		TLSKeyPath:                     envvar.NewString("TLSKeyPath").ParsePtr(),
		TLSPort:                        envvar.NewUint16("TLSPort").ParsePtr(),
		TLSRedirect:                    envvar.NewBool("TLSRedirect").ParsePtr(),
	}

	c.FeatureFeedsManager = envvar.NewBool("FeatureFeedsManager").ParsePtr()
	c.FeatureUICSAKeys = envvar.NewBool("FeatureUICSAKeys").ParsePtr()

	c.JobPipeline = &tcfg.JobPipelineConfig{
		DefaultHTTPLimit:          envvar.NewInt64("DefaultHTTPLimit").ParsePtr(),
		DefaultHTTPTimeout:        envDuration("DefaultHTTPTimeout"),
		FeatureExternalInitiators: envvar.NewBool("FeatureExternalInitiators").ParsePtr(),
		MaxRunDuration:            envDuration("JobPipelineMaxRunDuration"),
		ReaperInterval:            envDuration("JobPipelineReaperInterval"),
		ReaperThreshold:           envDuration("JobPipelineReaperThreshold"),
		ResultWriteQueueDepth:     envvar.NewUint32("JobPipelineResultWriteQueueDepth").ParsePtr(),
	}

	c.FMDefaultTransactionQueueDepth = envvar.NewUint32("FMDefaultTransactionQueueDepth").ParsePtr()
	c.FMSimulateTransactions = envvar.NewBool("FMSimulateTransactions").ParsePtr()

	c.FeatureOffchainReporting2 = envvar.NewBool("FeatureOffchainReporting2").ParsePtr()
	c.OCR2 = &tcfg.OCR2Config{
		ContractConfirmations:              envvar.NewUint32("OCR2ContractConfirmations").ParsePtr(),
		BlockchainTimeout:                  envDuration("OCR2BlockchainTimeout"),
		ContractPollInterval:               envDuration("OCR2ContractPollInterval"),
		ContractSubscribeInterval:          envDuration("OCR2ContractSubscribeInterval"),
		ContractTransmitterTransmitTimeout: envDuration("OCR2ContractTransmitterTransmitTimeout"),
		DatabaseTimeout:                    envDuration("OCR2DatabaseTimeout"),
		KeyBundleID:                        envvar.New("OCR2KeyBundleID", models.Sha256HashFromHex).ParsePtr(),
		MonitoringEndpoint:                 envvar.NewString("OCR2MonitoringEndpoint").ParsePtr(),
	}

	c.FeatureOffchainReporting = envvar.NewBool("FeatureOffchainReporting").ParsePtr()
	c.OCR = &tcfg.OCRConfig{
		ObservationTimeout:           envDuration("OCRObservationTimeout"),
		BlockchainTimeout:            envDuration("OCRBlockchainTimeout"),
		ContractPollInterval:         envDuration("OCRContractPollInterval"),
		ContractSubscribeInterval:    envDuration("OCRContractSubscribeInterval"),
		DefaultTransactionQueueDepth: envvar.NewUint32("OCRDefaultTransactionQueueDepth").ParsePtr(),
		KeyBundleID:                  envvar.New("OCRKeyBundleID", models.Sha256HashFromHex).ParsePtr(),
		MonitoringEndpoint:           envvar.NewString("OCRMonitoringEndpoint").ParsePtr(),
		SimulateTransactions:         envvar.NewBool("OCRSimulateTransactions").ParsePtr(),
		TraceLogging:                 envvar.NewBool("OCRTraceLogging").ParsePtr(),
		TransmitterAddress:           envvar.New("OCRTransmitterAddress", ethkey.NewEIP55Address).ParsePtr(),
		OutgoingMessageBufferSize:    envvar.NewInt64("OCROutgoingMessageBufferSize").ParsePtr(),
		IncomingMessageBufferSize:    envvar.NewInt64("OCRIncomingMessageBufferSize").ParsePtr(),
		DHTLookupInterval:            envvar.NewInt64("OCRDHTLookupInterval").ParsePtr(),
		BootstrapCheckInterval:       envDuration("OCRBootstrapCheckInterval"),
		NewStreamTimeout:             envDuration("OCRNewStreamTimeout"),
	}

	c.P2P = &tcfg.P2PConfig{
		NetworkingStack: envvar.New("P2PNetworkingStack", func(s string) (ns ocrnetworking.NetworkingStack, err error) {
			err = ns.UnmarshalText([]byte(s))
			return
		}).ParsePtr(),
		IncomingMessageBufferSize:        envvar.NewInt64("P2PIncomingMessageBufferSize").ParsePtr(),
		OutgoingMessageBufferSize:        envvar.NewInt64("P2POutgoingMessageBufferSize").ParsePtr(),
		AnnounceIP:                       envIP("P2PAnnounceIP"),
		AnnouncePort:                     envvar.NewUint16("P2PAnnouncePort").ParsePtr(),
		BootstrapCheckInterval:           envDuration("P2PBootstrapCheckInterval"),
		BootstrapPeers:                   envStringSlice("P2PBootstrapPeers"),
		DHTAnnouncementCounterUserPrefix: envvar.NewUint32("P2PDHTAnnouncementCounterUserPrefix").ParsePtr(),
		DHTLookupInterval:                envvar.NewInt64("P2PDHTLookupInterval").ParsePtr(),
		ListenIP:                         envIP("P2PListenIP"),
		ListenPort:                       envvar.NewUint16("P2PListenPort").ParsePtr(),
		NewStreamTimeout:                 envDuration("P2PNewStreamTimeout"),
		PeerID:                           envvar.New("P2PPeerID", p2pkey.MakePeerID).ParsePtr(),
		PeerstoreWriteInterval:           envDuration("P2PPeerstoreWriteInterval"),
		V2AnnounceAddresses:              envStringSlice("P2PV2AnnounceAddresses"),
		V2Bootstrappers:                  envStringSlice("P2PV2Bootstrappers"),
		V2DeltaDial:                      envDuration("P2PV2DeltaDial"),
		V2DeltaReconcile:                 envDuration("P2PV2DeltaReconcile"),
		V2ListenAddresses:                envStringSlice("P2PV2ListenAddresses"),
	}

	c.Keeper = &tcfg.KeeperConfig{
		CheckUpkeepGasPriceFeatureEnabled: envvar.NewBool("KeeperCheckUpkeepGasPriceFeatureEnabled").ParsePtr(),
		DefaultTransactionQueueDepth:      envvar.NewUint32("KeeperDefaultTransactionQueueDepth").ParsePtr(),
		GasPriceBufferPercent:             envvar.NewUint32("KeeperGasPriceBufferPercent").ParsePtr(),
		GasTipCapBufferPercent:            envvar.NewUint32("KeeperGasTipCapBufferPercent").ParsePtr(),
		BaseFeeBufferPercent:              envvar.NewUint32("KeeperBaseFeeBufferPercent").ParsePtr(),
		MaximumGracePeriod:                envvar.NewInt64("KeeperMaximumGracePeriod").ParsePtr(),
		RegistryCheckGasOverhead:          envBig("KeeperRegistryCheckGasOverhead"),
		RegistryPerformGasOverhead:        envBig("KeeperRegistryPerformGasOverhead"),
		RegistrySyncInterval:              envDuration("KeeperRegistrySyncInterval"),
		RegistrySyncUpkeepQueueSize:       envvar.KeeperRegistrySyncUpkeepQueueSize.ParsePtr(),
		TurnLookBack:                      envvar.NewInt64("KeeperTurnLookBack").ParsePtr(),
		TurnFlagEnabled:                   envvar.NewBool("KeeperTurnFlagEnabled").ParsePtr(),
	}

	c.AutoPprof = &tcfg.AutoPprofConfig{
		Enabled:              envvar.NewBool("AutoPprofEnabled").ParsePtr(),
		ProfileRoot:          envvar.NewString("AutoPprofProfileRoot").ParsePtr(),
		PollInterval:         envDuration("AutoPprofPollInterval"),
		GatherDuration:       envDuration("AutoPprofGatherDuration"),
		GatherTraceDuration:  envDuration("AutoPprofGatherTraceDuration"),
		MaxProfileSize:       envvar.New("AutoPprofMaxProfileSize", parse.FileSize).ParsePtr(),
		CPUProfileRate:       envvar.NewInt64("AutoPprofCPUProfileRate").ParsePtr(),
		MemProfileRate:       envvar.NewInt64("AutoPprofMemProfileRate").ParsePtr(),
		BlockProfileRate:     envvar.NewInt64("AutoPprofBlockProfileRate").ParsePtr(),
		MutexProfileFraction: envvar.NewInt64("AutoPprofMutexProfileFraction").ParsePtr(),
		MemThreshold:         envvar.New("AutoPprofMemThreshold", parse.FileSize).ParsePtr(),
		GoroutineThreshold:   envvar.NewInt64("AutoPprofGoroutineThreshold").ParsePtr(),
	}
}

func envDuration(s string) *models.Duration {
	if p := envvar.NewDuration(s).ParsePtr(); p != nil {
		d := *p
		if d >= 0 {
			return models.MustNewDuration(d)
		}
	}
	return nil
}

func envURL(s string) *tcfg.URL {
	if p := envvar.New(s, url.Parse).ParsePtr(); p != nil {
		return (*tcfg.URL)(*p)
	}
	return nil
}

func envIP(s string) *net.IP {
	return envvar.New(s, func(s string) (net.IP, error) {
		return net.ParseIP(s), nil
	}).ParsePtr()
}

func envStringSlice(s string) *[]string {
	return envvar.New(s, func(s string) ([]string, error) {
		// matching viper stringSlice logic
		t := strings.TrimSuffix(strings.TrimPrefix(s, "["), "]")
		return csv.NewReader(strings.NewReader(t)).Read()
	}).ParsePtr()
}

func envBig(s string) *utils.Big {
	return envvar.New(s, func(s string) (b utils.Big, err error) {
		err = b.UnmarshalText([]byte(s))
		return
	}).ParsePtr()
}
