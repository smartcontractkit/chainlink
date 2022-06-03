package chainlink

import (
	"context"
	"encoding/csv"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmcfg "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	legacy "github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/config/envvar"
	"github.com/smartcontractkit/chainlink/core/config/parse"
	config "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func (app *ChainlinkApplication) ConfigDump(ctx context.Context) (string, error) {
	var c Config

	if err := c.loadChainsAndNodes(ctx, app.Chains); err != nil {
		return "", err
	}

	c.loadLegacyEVMEnv()

	c.loadLegacyCoreEnv()

	b, err := toml.Marshal(c)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// loadChainsAndNodes initializes chains & nodes from configurations persisted in the database.
//TODO doc
func (c *Config) loadChainsAndNodes(ctx context.Context, chains Chains) error {
	{
		dbChains, _, err := chains.EVM.Index(0, -1)
		if err != nil {
			return err
		}
		for _, dbChain := range dbChains {
			dbNodes, _, err := chains.EVM.GetNodesForChain(ctx, dbChain.ID, 0, -1)
			if err != nil {
				return errors.Wrapf(err, "failed to get nodes for chain %v", dbChain.ID)
			}
			var evmChain EVMConfig
			if err := evmChain.setFromDB(dbChain, dbNodes); err != nil {
				return errors.Wrapf(err, "failed to convert db config for chain %v", dbChain.ID)
			}
			if *evmChain.Enabled {
				// no need to persist if enabled
				evmChain.Enabled = nil
			}
			c.EVM = append(c.EVM, evmChain)
		}
	}

	{
		dbChains, _, err := chains.Solana.Index(0, -1)
		if err != nil {
			return err
		}
		for _, dbChain := range dbChains {
			dbNodes, _, err := chains.Solana.GetNodesForChain(ctx, dbChain.ID, 0, -1)
			if err != nil {
				return errors.Wrapf(err, "failed to get nodes for chain %s", dbChain.ID)
			}
			var solChain SolanaConfig
			if err := solChain.setFromDB(dbChain, dbNodes); err != nil {
				return errors.Wrapf(err, "failed to convert db config for chain %s", dbChain.ID)
			}
			if *solChain.Enabled {
				// no need to persist if enabled
				solChain.Enabled = nil
			}
			c.Solana = append(c.Solana, solChain)
		}
	}

	{
		dbChains, _, err := chains.Terra.Index(0, -1)
		if err != nil {
			return err
		}
		for _, dbChain := range dbChains {
			dbNodes, _, err := chains.Terra.GetNodesForChain(ctx, dbChain.ID, 0, -1)
			if err != nil {
				return errors.Wrapf(err, "failed to get nodes for chain %s", dbChain.ID)
			}
			var terChain TerraConfig
			if err := terChain.setFromDB(dbChain, dbNodes); err != nil {
				return errors.Wrapf(err, "failed to convert db config for chain %s", dbChain.ID)
			}
			if *terChain.Enabled {
				// no need to persist if enabled
				terChain.Enabled = nil
			}
			c.Terra = append(c.Terra, terChain)
		}
	}

	return nil
}

// loadLegacyEVMEnv reads legacy ETH/EVM global overrides from the environment and updates all EVM chains.
//TODO test
func (c *Config) loadLegacyEVMEnv() {
	if e := envvar.NewBool("BalanceMonitorEnabled").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].BalanceMonitorEnabled = e
		}
	}
	if e := envvar.NewUint32("BlockBackfillDepth").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].BlockBackfillDepth = e
		}
	}
	if e := envvar.NewBool("BlockBackfillSkip").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].BlockBackfillSkip = e
		}
	}
	if e := envvar.NewDuration("BlockEmissionIdleWarningThreshold").ParsePtr(); e != nil {
		d := models.MustNewDuration(*e)
		for i := range c.EVM {
			c.EVM[i].BlockEmissionIdleWarningThreshold = d
		}
	}
	if e := envvar.NewDuration("EthTxReaperInterval").ParsePtr(); e != nil {
		d := models.MustNewDuration(*e)
		for i := range c.EVM {
			c.EVM[i].TxReaperInterval = d
		}
	}
	if e := envvar.NewDuration("EthTxReaperThreshold").ParsePtr(); e != nil {
		d := models.MustNewDuration(*e)
		for i := range c.EVM {
			c.EVM[i].TxReaperThreshold = d
		}
	}
	if e := envvar.NewDuration("EthTxResendAfterThreshold").ParsePtr(); e != nil {
		d := models.MustNewDuration(*e)
		for i := range c.EVM {
			c.EVM[i].TxResendAfterThreshold = d
		}
	}
	if e := envvar.NewUint32("EvmFinalityDepth").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].FinalityDepth = e
		}
	}
	if e := envvar.NewUint32("EvmHeadTrackerHistoryDepth").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].HeadTrackerHistoryDepth = e
		}
	}
	if e := envvar.NewUint32("EvmHeadTrackerMaxBufferSize").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].HeadTrackerMaxBufferSize = e
		}
	}
	if e := envvar.NewDuration("EvmHeadTrackerSamplingInterval").ParsePtr(); e != nil {
		d := models.MustNewDuration(*e)
		for i := range c.EVM {
			c.EVM[i].HeadTrackerSamplingInterval = d
		}
	}
	if e := envvar.NewUint32("EvmLogBackfillBatchSize").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].LogBackfillBatchSize = e
		}
	}
	if e := envvar.NewDuration("EvmLogPollInterval").ParsePtr(); e != nil {
		d := models.MustNewDuration(*e)
		for i := range c.EVM {
			c.EVM[i].LogPollInterval = d
		}
	}
	if e := envvar.NewUint32("EvmRPCDefaultBatchSize").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].RPCDefaultBatchSize = e
		}
	}
	if e := envvar.New("LinkContractAddress", ethkey.NewEIP55Address).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].LinkContractAddress = e
		}
	}
	if e := envvar.New("OperatorFactoryAddress", ethkey.NewEIP55Address).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].OperatorFactoryAddress = e
		}
	}
	if e := envvar.NewUint32("MinIncomingConfirmations").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].MinIncomingConfirmations = e
		}
	}
	if e := envvar.New("MinimumContractPayment", func(s string) (l assets.Link, err error) {
		err = l.UnmarshalText([]byte(s))
		return
	}).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].MinimumContractPayment = e
		}
	}
	if e := envvar.NewDuration("NodeNoNewHeadsThreshold").ParsePtr(); e != nil {
		d := models.MustNewDuration(*e)
		for i := range c.EVM {
			c.EVM[i].NodeNoNewHeadsThreshold = d
		}
	}
	if e := envvar.NewUint32("NodePollFailureThreshold").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].NodePollFailureThreshold = e
		}
	}
	if e := envvar.NewDuration("NodePollInterval").ParsePtr(); e != nil {
		d := models.MustNewDuration(*e)
		for i := range c.EVM {
			c.EVM[i].NodePollInterval = d
		}
	}
	if e := envvar.NewBool("EvmEIP1559DynamicFees").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].EIP1559DynamicFees = e
		}
	}
	if e := envvar.NewUint16("EvmGasBumpPercent").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasBumpPercent = e
		}
	}
	if e := envvar.New("EvmGasBumpThreshold", parse.BigInt).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasBumpThreshold = utils.NewBig(*e)
		}
	}
	if e := envvar.New("EvmGasBumpWei", parse.BigInt).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasBumpWei = utils.NewBig(*e)
		}
	}
	if e := envvar.New("EvmGasFeeCapDefault", parse.BigInt).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasFeeCapDefault = utils.NewBig(*e)
		}
	}
	if e := envvar.New("EvmGasLimitDefault", parse.BigInt).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasLimitDefault = utils.NewBig(*e)
		}
	}
	if e := envvar.New("EvmGasLimitMultiplier", decimal.NewFromString).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasLimitMultiplier = e
		}
	}
	if e := envvar.New("EvmGasLimitTransfer", parse.BigInt).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasLimitTransfer = utils.NewBig(*e)
		}
	}
	if e := envvar.New("EvmGasPriceDefault", parse.BigInt).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasPriceDefault = utils.NewBig(*e)
		}
	}
	if e := envvar.New("EvmGasTipCapDefault", parse.BigInt).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasTipCapDefault = utils.NewBig(*e)
		}
	}
	if e := envvar.New("EvmGasTipCapMinimum", parse.BigInt).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasTipCapMinimum = utils.NewBig(*e)
		}
	}
	if e := envvar.New("EvmMaxGasPriceWei", parse.BigInt).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].MaxGasPriceWei = utils.NewBig(*e)
		}
	}
	if e := envvar.New("EvmMinGasPriceWei", parse.BigInt).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].MinGasPriceWei = utils.NewBig(*e)
		}
	}
	if e := envvar.NewString("GasEstimatorMode").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimatorMode = e
		}
	} else if e, ok := os.LookupEnv("GAS_UPDATER_ENABLED"); ok {
		v := "FixedPrice"
		if b, err := strconv.ParseBool(e); err != nil && b {
			v = "BlockHistory"
		}
		for i := range c.EVM {
			c.EVM[i].GasEstimatorMode = &v
		}
	}
	if e := envvar.NewUint32("BlockHistoryEstimatorBatchSize").ParsePtr(); e != nil {
		for i := range c.EVM {
			if c.EVM[i].BlockHistoryEstimator == nil {
				c.EVM[i].BlockHistoryEstimator = &evmcfg.BlockHistoryEstimator{}
			}
			c.EVM[i].BlockHistoryEstimator.BatchSize = e
		}
	} else if s, ok := os.LookupEnv("GAS_UPDATER_BATCH_SIZE"); ok {
		l, err := parse.Uint32(s)
		if err == nil {
			for i := range c.EVM {
				if c.EVM[i].BlockHistoryEstimator == nil {
					c.EVM[i].BlockHistoryEstimator = &evmcfg.BlockHistoryEstimator{}
				}
				c.EVM[i].BlockHistoryEstimator.BatchSize = &l
			}
		}
	}
	if e := envvar.NewUint16("BlockHistoryEstimatorBlockDelay").ParsePtr(); e != nil {
		for i := range c.EVM {
			if c.EVM[i].BlockHistoryEstimator == nil {
				c.EVM[i].BlockHistoryEstimator = &evmcfg.BlockHistoryEstimator{}
			}
			c.EVM[i].BlockHistoryEstimator.BlockDelay = e
		}
	} else if s, ok := os.LookupEnv("GAS_UPDATER_BLOCK_DELAY"); ok {
		l, err := parse.Uint16(s)
		if err == nil {
			for i := range c.EVM {
				if c.EVM[i].BlockHistoryEstimator == nil {
					c.EVM[i].BlockHistoryEstimator = &evmcfg.BlockHistoryEstimator{}
				}
				c.EVM[i].BlockHistoryEstimator.BlockDelay = &l
			}
		}
	}
	if e := envvar.NewUint16("BlockHistoryEstimatorBlockHistorySize").ParsePtr(); e != nil {
		for i := range c.EVM {
			if c.EVM[i].BlockHistoryEstimator == nil {
				c.EVM[i].BlockHistoryEstimator = &evmcfg.BlockHistoryEstimator{}
			}
			c.EVM[i].BlockHistoryEstimator.BlockHistorySize = e
		}
	} else if s, ok := os.LookupEnv("GAS_UPDATER_BLOCK_HISTORY_SIZE"); ok {
		l, err := parse.Uint16(s)
		if err == nil {
			for i := range c.EVM {
				if c.EVM[i].BlockHistoryEstimator == nil {
					c.EVM[i].BlockHistoryEstimator = &evmcfg.BlockHistoryEstimator{}
				}
				c.EVM[i].BlockHistoryEstimator.BlockHistorySize = &l
			}
		}
	}
	if e := envvar.NewUint16("BlockHistoryEstimatorEIP1559FeeCapBufferBlocks").ParsePtr(); e != nil {
		for i := range c.EVM {
			if c.EVM[i].BlockHistoryEstimator == nil {
				c.EVM[i].BlockHistoryEstimator = &evmcfg.BlockHistoryEstimator{}
			}
			c.EVM[i].BlockHistoryEstimator.EIP1559FeeCapBufferBlocks = e
		}
	}
	if e := envvar.NewUint16("BlockHistoryEstimatorTransactionPercentile").ParsePtr(); e != nil {
		for i := range c.EVM {
			if c.EVM[i].BlockHistoryEstimator == nil {
				c.EVM[i].BlockHistoryEstimator = &evmcfg.BlockHistoryEstimator{}
			}
			c.EVM[i].BlockHistoryEstimator.TransactionPercentile = e
		}
	} else if s, ok := os.LookupEnv("GAS_UPDATER_TRANSACTION_PERCENTILE"); ok {
		l, err := parse.Uint16(s)
		if err == nil {
			for i := range c.EVM {
				if c.EVM[i].BlockHistoryEstimator == nil {
					c.EVM[i].BlockHistoryEstimator = &evmcfg.BlockHistoryEstimator{}
				}
				c.EVM[i].BlockHistoryEstimator.TransactionPercentile = &l
			}
		}
	}
	if e := envvar.NewUint16("EvmGasBumpTxDepth").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasBumpTxDepth = e
		}
	}
	if e := envvar.NewUint32("EvmMaxInFlightTransactions").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].MaxInFlightTransactions = e
		}
	}
	if e := envvar.NewUint32("EvmMaxQueuedTransactions").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].MaxInFlightTransactions = e
		}
	}
	if e := envvar.NewBool("EvmNonceAutoSync").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].NonceAutoSync = e
		}
	}
	if e := envvar.NewBool("EvmUseForwarders").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].UseForwarders = e
		}
	}
}

// loadLegacyCoreEnv loads Core values from legacy environment variables.
//TODO test
func (c *Config) loadLegacyCoreEnv() {
	c.Dev = envvar.NewBool("Dev").ParsePtr()
	c.ExplorerURL = envURL("ExplorerURL")
	c.InsecureFastScrypt = envvar.NewBool("InsecureFastScrypt").ParsePtr()
	c.ReaperExpiration = envDuration("ReaperExpiration")
	c.Root = envvar.RootDir.ParsePtr()
	c.ShutdownGracePeriod = envDuration("ShutdownGracePeriod")

	c.Database = &config.Database{
		ListenerMaxReconnectDuration:  envDuration("DatabaseListenerMaxReconnectDuration"),
		ListenerMinReconnectInterval:  envDuration("DatabaseListenerMinReconnectInterval"),
		Migrate:                       envvar.NewBool("MigrateDatabase").ParsePtr(),
		ORMMaxIdleConns:               envvar.NewInt64("ORMMaxIdleConns").ParsePtr(),
		ORMMaxOpenConns:               envvar.NewInt64("ORMMaxOpenConns").ParsePtr(),
		TriggerFallbackDBPollInterval: envDuration("TriggerFallbackDBPollInterval"),
		Lock: &config.DatabaseLock{
			Mode:                  envvar.NewString("DatabaseLockingMode").ParsePtr(),
			AdvisoryCheckInterval: envDuration("AdvisoryLockCheckInterval"),
			AdvisoryID:            envvar.AdvisoryLockID.ParsePtr(),
			LeaseDuration:         envDuration("LeaseLockDuration"),
			LeaseRefreshInterval:  envDuration("LeaseLockRefreshInterval"),
		},
		Backup: &config.DatabaseBackup{
			Dir:              envvar.NewString("DatabaseBackupDir").ParsePtr(),
			Frequency:        envDuration("DatabaseBackupFrequency"),
			Mode:             legacy.DatabaseBackupModeEnvVar.ParsePtr(),
			OnVersionUpgrade: envvar.NewBool("DatabaseBackupOnVersionUpgrade").ParsePtr(),
			URL:              envURL("DatabaseBackupDir"),
		},
	}

	c.TelemetryIngress = &config.TelemetryIngress{
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

	c.Log = &config.Log{
		JSONConsole:    envvar.JSONConsole.ParsePtr(),
		FileDir:        envvar.NewString("LogFileDir").ParsePtr(),
		Level:          envvar.LogLevel.ParsePtr(),
		SQL:            envvar.NewBool("LogSQL").ParsePtr(),
		FileMaxSize:    envvar.LogFileMaxSize.ParsePtr(),
		FileMaxAgeDays: envvar.LogFileMaxAge.ParsePtr(),
		FileMaxBackups: envvar.LogFileMaxBackups.ParsePtr(),
		UnixTS:         envvar.LogUnixTS.ParsePtr(),
	}

	c.WebServer = &config.WebServer{
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
		MFA: &config.WebServerMFA{
			RPID:     envvar.NewString("RPID").ParsePtr(),
			RPOrigin: envvar.NewString("RPOrigin").ParsePtr(),
		},
		TLS: &config.WebServerTLS{
			CertPath: envvar.NewString("TLSCertPath").ParsePtr(),
			Host:     envvar.NewString("TLSHost").ParsePtr(),
			KeyPath:  envvar.NewString("TLSKeyPath").ParsePtr(),
			Port:     envvar.NewUint16("TLSPort").ParsePtr(),
			Redirect: envvar.NewBool("TLSRedirect").ParsePtr(),
		},
	}

	c.FeatureFeedsManager = envvar.NewBool("FeatureFeedsManager").ParsePtr()
	c.FeatureUICSAKeys = envvar.NewBool("FeatureUICSAKeys").ParsePtr()

	c.JobPipeline = &config.JobPipeline{
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
	c.OCR2 = &config.OCR2{
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
	c.OCR = &config.OCR{
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
	}

	c.P2P = &config.P2P{
		IncomingMessageBufferSize: first(envvar.NewInt64("OCRIncomingMessageBufferSize"), envvar.NewInt64("P2PIncomingMessageBufferSize")),
		OutgoingMessageBufferSize: first(envvar.NewInt64("OCROutgoingMessageBufferSize"), envvar.NewInt64("P2POutgoingMessageBufferSize")),
	}
	if p := envvar.New("P2PNetworkingStack", func(s string) (ns ocrnetworking.NetworkingStack, err error) {
		err = ns.UnmarshalText([]byte(s))
		return
	}).ParsePtr(); p != nil {
		ns := *p
		var v1, v2, v1v2 = ocrnetworking.NetworkingStackV1, ocrnetworking.NetworkingStackV2, ocrnetworking.NetworkingStackV1V2
		if ns == v1 || ns == v1v2 {
			c.P2P.V1 = &config.P2PV1{
				AnnounceIP:                       envIP("P2PAnnounceIP"),
				AnnouncePort:                     envvar.NewUint16("P2PAnnouncePort").ParsePtr(),
				BootstrapCheckInterval:           envDuration("OCRBootstrapCheckInterval", "P2PBootstrapCheckInterval"),
				BootstrapPeers:                   envStringSlice("P2PBootstrapPeers"),
				DHTAnnouncementCounterUserPrefix: envvar.NewUint32("P2PDHTAnnouncementCounterUserPrefix").ParsePtr(),
				DHTLookupInterval:                first(envvar.NewInt64("OCRDHTLookupInterval"), envvar.NewInt64("P2PDHTLookupInterval")),
				ListenIP:                         envIP("P2PListenIP"),
				ListenPort:                       envvar.NewUint16("P2PListenPort").ParsePtr(),
				NewStreamTimeout:                 envDuration("OCRNewStreamTimeout", "P2PNewStreamTimeout"),
				PeerID:                           envvar.New("P2PPeerID", p2pkey.MakePeerID).ParsePtr(),
				PeerstoreWriteInterval:           envDuration("P2PPeerstoreWriteInterval"),
			}
		}
		if ns == v2 || ns == v1v2 {
			c.P2P.V2 = &config.P2PV2{
				AnnounceAddresses: envStringSlice("P2PV2AnnounceAddresses"),
				Bootstrappers:     envStringSlice("P2PV2Bootstrappers"),
				DeltaDial:         envDuration("P2PV2DeltaDial"),
				DeltaReconcile:    envDuration("P2PV2DeltaReconcile"),
				ListenAddresses:   envStringSlice("P2PV2ListenAddresses"),
			}
		}
	}

	c.Keeper = &config.Keeper{
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

	c.AutoPprof = &config.AutoPprof{
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

func first[T any](es ...*envvar.EnvVar[T]) *T {
	for _, e := range es {
		if p := e.ParsePtr(); p != nil {
			return p
		}
	}
	return nil
}

func envDuration(ns ...string) *models.Duration {
	for _, n := range ns {
		if p := envvar.NewDuration(n).ParsePtr(); p != nil {
			d := *p
			if d >= 0 {
				return models.MustNewDuration(d)
			}
		}
	}
	return nil
}

func envURL(s string) *models.URL {
	if p := envvar.New(s, models.ParseURL).ParsePtr(); p != nil {
		return *p
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
