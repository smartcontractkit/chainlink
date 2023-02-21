package chainlink

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"

	soldb "github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	stkdb "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/core/chains/starknet"
	stktyp "github.com/smartcontractkit/chainlink/core/chains/starknet/types"
	"github.com/smartcontractkit/chainlink/core/logger/audit"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmcfg "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	evmtyp "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/chains/solana"
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
	var chains dbData
	if app != nil {
		var err error
		if app.Chains.EVM != nil {
			chains.EVMChains, _, err = app.Chains.EVM.Index(0, -1)
			if err != nil {
				return "", err
			}
			chains.EVMNodes = make(map[string][]evmtyp.Node)
			for _, dbChain := range chains.EVMChains {
				chains.EVMNodes[dbChain.ID.String()], _, err = app.Chains.EVM.GetNodesForChain(ctx, dbChain.ID, 0, -1)
				if err != nil {
					return "", errors.Wrapf(err, "failed to get nodes for evm chain %v", dbChain.ID)
				}
			}
		}

		if app.Chains.Solana != nil {
			chains.SolanaChains, _, err = app.Chains.Solana.Index(0, -1)
			if err != nil {
				return "", err
			}

			chains.SolanaNodes = make(map[string][]soldb.Node)
			for _, dbChain := range chains.SolanaChains {
				chains.SolanaNodes[dbChain.ID], _, err = app.Chains.Solana.GetNodesForChain(ctx, dbChain.ID, 0, -1)
				if err != nil {
					return "", errors.Wrapf(err, "failed to get nodes for solana chain %s", dbChain.ID)
				}
			}
		}
	}
	return configDump(chains)
}

type dbData struct {
	EVMChains []evmtyp.DBChain
	EVMNodes  map[string][]evmtyp.Node

	SolanaChains []solana.DBChain
	SolanaNodes  map[string][]soldb.Node

	StarknetChains []stktyp.DBChain
	StarknetNodes  map[string][]stkdb.Node
}

func configDump(data dbData) (string, error) {
	var c Config

	if err := c.loadChainsAndNodes(data); err != nil {
		return "", err
	}

	c.loadLegacyEVMEnv()

	if err := c.loadLegacyCoreEnv(); err != nil {
		return "", err
	}

	return c.TOMLString()
}

// loadChainsAndNodes initializes chains & nodes from configurations persisted in the database.
func (c *Config) loadChainsAndNodes(dbData dbData) error {
	for _, dbChain := range dbData.EVMChains {
		var evmChain evmcfg.EVMConfig
		if err := evmChain.SetFromDB(dbChain, dbData.EVMNodes[dbChain.ID.String()]); err != nil {
			return errors.Wrapf(err, "failed to convert db config for evm chain %v", dbChain.ID)
		}
		if *evmChain.Enabled {
			// no need to persist if enabled
			evmChain.Enabled = nil
		}
		c.EVM = append(c.EVM, &evmChain)
	}

	for _, dbChain := range dbData.SolanaChains {
		var solChain solana.SolanaConfig
		if err := solChain.SetFromDB(dbChain, dbData.SolanaNodes[dbChain.ID]); err != nil {
			return errors.Wrapf(err, "failed to convert db config for solana chain %s", dbChain.ID)
		}
		if *solChain.Enabled {
			// no need to persist if enabled
			solChain.Enabled = nil
		}
		c.Solana = append(c.Solana, &solChain)
	}

	for _, dbChain := range dbData.StarknetChains {
		var starkChain starknet.StarknetConfig
		if err := starkChain.SetFromDB(dbChain, dbData.StarknetNodes[dbChain.ID]); err != nil {
			return errors.Wrapf(err, "failed to convert db config for starknet chain %s", dbChain.ID)
		}
		if *starkChain.Enabled {
			// no need to persist if enabled
			starkChain.Enabled = nil
		}
		c.Starknet = append(c.Starknet, &starkChain)
	}

	return nil
}

// loadLegacyEVMEnv reads legacy ETH/EVM global overrides from the environment and updates all EVM chains.
func (c *Config) loadLegacyEVMEnv() {
	if e := envvar.NewBool("BalanceMonitorEnabled").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].BalanceMonitor.Enabled = e
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
	if e := envvar.NewString("ChainType").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].ChainType = e
		}
	}
	if e := envvar.NewDuration("EthTxReaperInterval").ParsePtr(); e != nil {
		d := models.MustNewDuration(*e)
		for i := range c.EVM {
			c.EVM[i].Transactions.ReaperInterval = d
		}
	}
	if e := envvar.NewDuration("EthTxReaperThreshold").ParsePtr(); e != nil {
		d := models.MustNewDuration(*e)
		for i := range c.EVM {
			c.EVM[i].Transactions.ReaperThreshold = d
		}
	}
	if e := envvar.NewDuration("EthTxResendAfterThreshold").ParsePtr(); e != nil {
		d := models.MustNewDuration(*e)
		for i := range c.EVM {
			c.EVM[i].Transactions.ResendAfterThreshold = d
		}
	}
	if e := envvar.NewUint32("EvmFinalityDepth").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].FinalityDepth = e
		}
	}
	if e := envvar.NewUint32("EvmHeadTrackerHistoryDepth").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].HeadTracker.HistoryDepth = e
		}
	}
	if e := envvar.NewUint32("EvmHeadTrackerMaxBufferSize").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].HeadTracker.MaxBufferSize = e
		}
	}
	if e := envvar.NewDuration("EvmHeadTrackerSamplingInterval").ParsePtr(); e != nil {
		d := models.MustNewDuration(*e)
		for i := range c.EVM {
			c.EVM[i].HeadTracker.SamplingInterval = d
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
	if e := envvar.NewUint32("EvmLogKeepBlocksDepth").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].LogKeepBlocksDepth = e
		}
	}
	if e := envvar.NewUint32("EvmRPCDefaultBatchSize").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].RPCDefaultBatchSize = e
		}
	}
	if e := envvar.New("FlagsContractAddress", ethkey.NewEIP55Address).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].FlagsContractAddress = e
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
			c.EVM[i].MinContractPayment = e
		}
	}
	if e := envvar.NewDuration("NodeNoNewHeadsThreshold").ParsePtr(); e != nil {
		d := models.MustNewDuration(*e)
		for i := range c.EVM {
			c.EVM[i].NoNewHeadsThreshold = d
		}
	} else if e := envvar.NewDuration("BlockEmissionIdleWarningThreshold").ParsePtr(); e != nil {
		d := models.MustNewDuration(*e)
		for i := range c.EVM {
			c.EVM[i].NoNewHeadsThreshold = d
		}
	}
	if e := envvar.NewUint32("NodePollFailureThreshold").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].NodePool.PollFailureThreshold = e
		}
	}
	if e := envvar.NewDuration("NodePollInterval").ParsePtr(); e != nil {
		d := models.MustNewDuration(*e)
		for i := range c.EVM {
			c.EVM[i].NodePool.PollInterval = d
		}
	}
	if e := envvar.NewString("NodeSelectionMode").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].NodePool.SelectionMode = e
		}
	}
	if e := envvar.NewUint32("NodeSyncThreshold").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].NodePool.SyncThreshold = e
		}
	}
	if e := envvar.NewBool("EvmEIP1559DynamicFees").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.EIP1559DynamicFees = e
		}
	}
	if e := envvar.NewUint16("EvmGasBumpPercent").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.BumpPercent = e
		}
	}
	if e := envvar.NewUint32("EvmGasBumpThreshold").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.BumpThreshold = e
		}
	}
	if e := envvar.New("EvmGasBumpWei", parse.BigInt).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.BumpMin = assets.NewWei(*e)
		}
	}
	if e := envvar.New("EvmGasFeeCapDefault", parse.BigInt).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.FeeCapDefault = assets.NewWei(*e)
		}
	}
	if e := envvar.NewUint32("EvmGasLimitDefault").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.LimitDefault = e
		}
	}
	if e := envvar.NewUint32("EvmGasLimitMax").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.LimitMax = e
		}
	}
	if e := envvar.NewUint32("EvmGasLimitOCRJobType").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.LimitJobType.OCR = e
		}
	}
	if e := envvar.NewUint32("EvmGasLimitDRJobType").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.LimitJobType.DR = e
		}
	}
	if e := envvar.NewUint32("EvmGasLimitVRFJobType").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.LimitJobType.VRF = e
		}
	}
	if e := envvar.NewUint32("EvmGasLimitFMJobType").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.LimitJobType.FM = e
		}
	}
	if e := envvar.NewUint32("EvmGasLimitKeeperJobType").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.LimitJobType.Keeper = e
		}
	}
	if e := envvar.New("EvmGasLimitMultiplier", decimal.NewFromString).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.LimitMultiplier = e
		}
	}
	if e := envvar.NewUint32("EvmGasLimitTransfer").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.LimitTransfer = e
		}
	}
	if e := envvar.New("EvmGasPriceDefault", parse.BigInt).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.PriceDefault = assets.NewWei(*e)
		}
	}
	if e := envvar.New("EvmGasTipCapDefault", parse.BigInt).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.TipCapDefault = assets.NewWei(*e)
		}
	}
	if e := envvar.New("EvmGasTipCapMinimum", parse.BigInt).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.TipCapMin = assets.NewWei(*e)
		}
	}
	if e := envvar.New("EvmMaxGasPriceWei", parse.BigInt).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.PriceMax = assets.NewWei(*e)
		}
	}
	if e := envvar.New("EvmMinGasPriceWei", parse.BigInt).ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.PriceMin = assets.NewWei(*e)
		}
	}
	if e := envvar.NewString("GasEstimatorMode").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.Mode = e
		}
	} else if e, ok := os.LookupEnv("GAS_UPDATER_ENABLED"); ok && e != "" {
		v := "FixedPrice"
		if b, err := strconv.ParseBool(e); err != nil && b {
			v = "BlockHistory"
		}
		for i := range c.EVM {
			c.EVM[i].GasEstimator.Mode = &v
		}
	}
	if e := envvar.NewUint16("EvmGasBumpTxDepth").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.BumpTxDepth = e
		}
	}
	if e := envvar.NewUint32("BlockHistoryEstimatorBatchSize").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.BlockHistory.BatchSize = e
		}
	} else if s, ok := os.LookupEnv("GAS_UPDATER_BATCH_SIZE"); ok {
		l, err := parse.Uint32(s)
		if err == nil {
			for i := range c.EVM {
				c.EVM[i].GasEstimator.BlockHistory.BatchSize = &l
			}
		}
	}
	if e := envvar.NewUint16("BlockHistoryEstimatorBlockDelay").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].RPCBlockQueryDelay = e
		}
	} else if s, ok := os.LookupEnv("GAS_UPDATER_BLOCK_DELAY"); ok {
		l, err := parse.Uint16(s)
		if err == nil {
			for i := range c.EVM {
				c.EVM[i].RPCBlockQueryDelay = &l
			}
		}
	}
	if e := envvar.NewUint16("BlockHistoryEstimatorBlockHistorySize").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.BlockHistory.BlockHistorySize = e
		}
	} else if s, ok := os.LookupEnv("GAS_UPDATER_BLOCK_HISTORY_SIZE"); ok {
		l, err := parse.Uint16(s)
		if err == nil {
			for i := range c.EVM {
				c.EVM[i].GasEstimator.BlockHistory.BlockHistorySize = &l
			}
		}
	}
	if e := envvar.NewUint16("BlockHistoryEstimatorCheckInclusionBlocks").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.BlockHistory.CheckInclusionBlocks = e
		}
	}
	if e := envvar.NewUint16("BlockHistoryEstimatorCheckInclusionPercentile").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.BlockHistory.CheckInclusionPercentile = e
		}
	}
	if e := envvar.NewUint16("BlockHistoryEstimatorEIP1559FeeCapBufferBlocks").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.BlockHistory.EIP1559FeeCapBufferBlocks = e
		}
	}
	if e := envvar.NewUint16("BlockHistoryEstimatorTransactionPercentile").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].GasEstimator.BlockHistory.TransactionPercentile = e
		}
	} else if s, ok := os.LookupEnv("GAS_UPDATER_TRANSACTION_PERCENTILE"); ok {
		l, err := parse.Uint16(s)
		if err == nil {
			for i := range c.EVM {
				c.EVM[i].GasEstimator.BlockHistory.TransactionPercentile = &l
			}
		}
	}
	if e := envvar.NewUint32("EvmMaxInFlightTransactions").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].Transactions.MaxInFlight = e
		}
	}
	if e := envvar.NewUint32("EvmMaxQueuedTransactions").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].Transactions.MaxQueued = e
		}
	}
	if e := envvar.NewBool("EvmNonceAutoSync").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].NonceAutoSync = e
		}
	}
	if e := envvar.NewBool("EvmUseForwarders").ParsePtr(); e != nil {
		for i := range c.EVM {
			c.EVM[i].Transactions.ForwardersEnabled = e
		}
	}
}

// loadLegacyCoreEnv loads Core values from legacy environment variables.
func (c *Config) loadLegacyCoreEnv() error {
	c.ExplorerURL = envURL("ExplorerURL")
	c.InsecureFastScrypt = envvar.NewBool("InsecureFastScrypt").ParsePtr()
	c.RootDir = envvar.RootDir.ParsePtr()
	c.ShutdownGracePeriod = envDuration("ShutdownGracePeriod")

	c.Feature = config.Feature{
		FeedsManager: envvar.NewBool("FeatureFeedsManager").ParsePtr(),
		LogPoller:    envvar.NewBool("FeatureLogPoller").ParsePtr(),
		UICSAKeys:    envvar.NewBool("FeatureUICSAKeys").ParsePtr(),
	}
	c.AuditLogger = audit.AuditLoggerConfig{
		Enabled:        envvar.NewBool("AuditLoggerEnabled").ParsePtr(),
		ForwardToUrl:   envURL("AuditLoggerForwardToUrl"),
		JsonWrapperKey: envvar.NewString("AuditLoggerJsonWrapperKey").ParsePtr(),
		Headers:        (*[]audit.ServiceHeader)(audit.AuditLoggerHeaders.ParsePtr()),
	}

	var lockEnabled *bool
	if mode := envvar.NewString("DatabaseLockingMode").ParsePtr(); mode == nil { // dual default
		lockEnabled = nil // lease default
	} else {
		switch *mode {
		case "advisorylock":
			return fmt.Errorf("%w: '%s' mode: must use 'lease' or 'none'", config.ErrUnsupported, *mode)
		case "none":
			lockEnabled = ptr(false)
		case "lease", "dual":
			lockEnabled = nil // lease default
		default:
			return fmt.Errorf("%w: unrecognized mode '%s': must use one of 'lease', 'dual', or 'none'", config.ErrUnsupported, *mode)
		}
	}
	c.Database = config.Database{
		DefaultIdleInTxSessionTimeout: mustParseDuration(os.Getenv("DATABASE_DEFAULT_IDLE_IN_TX_SESSION_TIMEOUT")),
		DefaultLockTimeout:            mustParseDuration(os.Getenv("DATABASE_DEFAULT_LOCK_TIMEOUT")),
		DefaultQueryTimeout:           mustParseDuration(os.Getenv("DATABASE_DEFAULT_QUERY_TIMEOUT")),
		LogQueries:                    envvar.LogSQL.ParsePtr(),
		MigrateOnStartup:              envvar.NewBool("MigrateDatabase").ParsePtr(),
		MaxIdleConns:                  envvar.NewInt64("ORMMaxIdleConns").ParsePtr(),
		MaxOpenConns:                  envvar.NewInt64("ORMMaxOpenConns").ParsePtr(),
		Listener: config.DatabaseListener{
			MaxReconnectDuration: envDuration("DatabaseListenerMaxReconnectDuration"),
			MinReconnectInterval: envDuration("DatabaseListenerMinReconnectInterval"),
			FallbackPollInterval: envDuration("TriggerFallbackDBPollInterval"),
		},
		Lock: config.DatabaseLock{
			Enabled:              lockEnabled,
			LeaseDuration:        envDuration("LeaseLockDuration"),
			LeaseRefreshInterval: envDuration("LeaseLockRefreshInterval"),
		},
		Backup: config.DatabaseBackup{
			Dir:              envvar.NewString("DatabaseBackupDir").ParsePtr(),
			Frequency:        envDuration("DatabaseBackupFrequency"),
			Mode:             legacy.DatabaseBackupModeEnvVar.ParsePtr(),
			OnVersionUpgrade: envvar.NewBool("DatabaseBackupOnVersionUpgrade").ParsePtr(),
		},
	}

	c.TelemetryIngress = config.TelemetryIngress{
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

	c.Log = config.Log{
		Level:       (*config.LogLevel)(envvar.LogLevel.ParsePtr()),
		JSONConsole: envvar.JSONConsole.ParsePtr(),
		UnixTS:      envvar.LogUnixTS.ParsePtr(),
		File: config.LogFile{
			Dir:        envvar.NewString("LogFileDir").ParsePtr(),
			MaxSize:    envvar.LogFileMaxSize.ParsePtr(),
			MaxAgeDays: envvar.LogFileMaxAge.ParsePtr(),
			MaxBackups: envvar.LogFileMaxBackups.ParsePtr(),
		},
	}

	c.WebServer = config.WebServer{
		AllowOrigins:            envvar.NewString("AllowOrigins").ParsePtr(),
		BridgeResponseURL:       envURL("BridgeResponseURL"),
		BridgeCacheTTL:          envDuration("BridgeCacheTTL"),
		HTTPWriteTimeout:        envDuration("HTTPServerWriteTimeout"),
		HTTPPort:                envvar.NewUint16("Port").ParsePtr(),
		SecureCookies:           envvar.NewBool("SecureCookies").ParsePtr(),
		SessionTimeout:          envDuration("SessionTimeout"),
		SessionReaperExpiration: envDuration("ReaperExpiration"),
		MFA: config.WebServerMFA{
			RPID:     envvar.NewString("RPID").ParsePtr(),
			RPOrigin: envvar.NewString("RPOrigin").ParsePtr(),
		},
		RateLimit: config.WebServerRateLimit{
			Authenticated:         envvar.NewInt64("AuthenticatedRateLimit").ParsePtr(),
			AuthenticatedPeriod:   envDuration("AuthenticatedRateLimitPeriod"),
			Unauthenticated:       envvar.NewInt64("UnAuthenticatedRateLimit").ParsePtr(),
			UnauthenticatedPeriod: envDuration("UnAuthenticatedRateLimitPeriod"),
		},
		TLS: config.WebServerTLS{
			CertPath:      envvar.NewString("TLSCertPath").ParsePtr(),
			Host:          envvar.NewString("TLSHost").ParsePtr(),
			KeyPath:       envvar.NewString("TLSKeyPath").ParsePtr(),
			HTTPSPort:     envvar.NewUint16("TLSPort").ParsePtr(),
			ForceRedirect: envvar.NewBool("TLSRedirect").ParsePtr(),
		},
	}

	c.JobPipeline = config.JobPipeline{
		ExternalInitiatorsEnabled: envvar.NewBool("FeatureExternalInitiators").ParsePtr(),
		MaxRunDuration:            envDuration("JobPipelineMaxRunDuration"),
		MaxSuccessfulRuns:         envvar.NewUint64("JobPipelineMaxSuccessfulRuns").ParsePtr(),
		ReaperInterval:            envDuration("JobPipelineReaperInterval"),
		ReaperThreshold:           envDuration("JobPipelineReaperThreshold"),
		ResultWriteQueueDepth:     envvar.NewUint32("JobPipelineResultWriteQueueDepth").ParsePtr(),
		HTTPRequest: config.JobPipelineHTTPRequest{
			DefaultTimeout: envDuration("DefaultHTTPTimeout"),
		},
	}
	if p := envvar.NewInt64("DefaultHTTPLimit").ParsePtr(); p != nil {
		b := utils.FileSize(*p)
		c.JobPipeline.HTTPRequest.MaxSize = &b
	}

	c.FluxMonitor = config.FluxMonitor{
		DefaultTransactionQueueDepth: envvar.NewUint32("FMDefaultTransactionQueueDepth").ParsePtr(),
		SimulateTransactions:         envvar.NewBool("FMSimulateTransactions").ParsePtr(),
	}

	c.OCR2 = config.OCR2{
		Enabled:                            envvar.NewBool("FeatureOffchainReporting2").ParsePtr(),
		ContractConfirmations:              envvar.NewUint32("OCR2ContractConfirmations").ParsePtr(),
		BlockchainTimeout:                  envDuration("OCR2BlockchainTimeout"),
		ContractPollInterval:               envDuration("OCR2ContractPollInterval"),
		ContractSubscribeInterval:          envDuration("OCR2ContractSubscribeInterval"),
		ContractTransmitterTransmitTimeout: envDuration("OCR2ContractTransmitterTransmitTimeout"),
		DatabaseTimeout:                    envDuration("OCR2DatabaseTimeout"),
		KeyBundleID:                        envvar.New("OCR2KeyBundleID", models.Sha256HashFromHex).ParsePtr(),
	}

	c.OCR = config.OCR{
		Enabled:                      envvar.NewBool("FeatureOffchainReporting").ParsePtr(),
		ObservationTimeout:           envDuration("OCRObservationTimeout"),
		BlockchainTimeout:            envDuration("OCRBlockchainTimeout"),
		ContractPollInterval:         envDuration("OCRContractPollInterval"),
		ContractSubscribeInterval:    envDuration("OCRContractSubscribeInterval"),
		DefaultTransactionQueueDepth: envvar.NewUint32("OCRDefaultTransactionQueueDepth").ParsePtr(),
		KeyBundleID:                  envvar.New("OCRKeyBundleID", models.Sha256HashFromHex).ParsePtr(),
		SimulateTransactions:         envvar.NewBool("OCRSimulateTransactions").ParsePtr(),
		TransmitterAddress:           envvar.New("OCRTransmitterAddress", ethkey.NewEIP55Address).ParsePtr(),
	}

	c.P2P = config.P2P{
		IncomingMessageBufferSize: first(envvar.NewInt64("OCRIncomingMessageBufferSize"), envvar.NewInt64("P2PIncomingMessageBufferSize")),
		OutgoingMessageBufferSize: first(envvar.NewInt64("OCROutgoingMessageBufferSize"), envvar.NewInt64("P2POutgoingMessageBufferSize")),
		PeerID:                    envvar.New("P2PPeerID", p2pkey.MakePeerID).ParsePtr(),
		TraceLogging:              envvar.NewBool("OCRTraceLogging").ParsePtr(),
	}
	p := envvar.New("P2PNetworkingStack", func(s string) (ns ocrnetworking.NetworkingStack, err error) {
		err = ns.UnmarshalText([]byte(s))
		return
	}).ParsePtr()
	var v1, v2, v1v2 = ocrnetworking.NetworkingStackV1, ocrnetworking.NetworkingStackV2, ocrnetworking.NetworkingStackV1V2
	if p == nil {
		p = &v1 // legacy default
	}
	ns := *p
	c.P2P.V1 = config.P2PV1{
		AnnounceIP:                       envIP("P2PAnnounceIP"),
		AnnouncePort:                     envvar.NewUint16("P2PAnnouncePort").ParsePtr(),
		BootstrapCheckInterval:           envDuration("OCRBootstrapCheckInterval", "P2PBootstrapCheckInterval"),
		DefaultBootstrapPeers:            envStringSlice("P2PBootstrapPeers"),
		DHTAnnouncementCounterUserPrefix: envvar.NewUint32("P2PDHTAnnouncementCounterUserPrefix").ParsePtr(),
		DHTLookupInterval:                first(envvar.NewInt64("OCRDHTLookupInterval"), envvar.NewInt64("P2PDHTLookupInterval")),
		ListenIP:                         envIP("P2PListenIP"),
		ListenPort:                       envvar.NewUint16("P2PListenPort").ParsePtr(),
		NewStreamTimeout:                 envDuration("OCRNewStreamTimeout", "P2PNewStreamTimeout"),
		PeerstoreWriteInterval:           envDuration("P2PPeerstoreWriteInterval"),
	}
	if ns == v2 {
		c.P2P.V1.Enabled = ptr(false)
	}

	c.P2P.V2 = config.P2PV2{
		AnnounceAddresses: envStringSlice("P2PV2AnnounceAddresses"),
		DefaultBootstrappers: envSlice("P2PV2Bootstrappers", func(v *ocrcommontypes.BootstrapperLocator, b []byte) error {
			return v.UnmarshalText(b)
		}),
		DeltaDial:       envDuration("P2PV2DeltaDial"),
		DeltaReconcile:  envDuration("P2PV2DeltaReconcile"),
		ListenAddresses: envStringSlice("P2PV2ListenAddresses"),
	}
	if (ns == v2 || ns == v1v2) && c.P2P.V2 != (config.P2PV2{}) {
		c.P2P.V2.Enabled = ptr(true)
	}

	c.Keeper = config.Keeper{
		DefaultTransactionQueueDepth: envvar.NewUint32("KeeperDefaultTransactionQueueDepth").ParsePtr(),
		GasPriceBufferPercent:        envvar.NewUint16("KeeperGasPriceBufferPercent").ParsePtr(),
		GasTipCapBufferPercent:       envvar.NewUint16("KeeperGasTipCapBufferPercent").ParsePtr(),
		BaseFeeBufferPercent:         envvar.NewUint16("KeeperBaseFeeBufferPercent").ParsePtr(),
		MaxGracePeriod:               envvar.NewInt64("KeeperMaximumGracePeriod").ParsePtr(),
		TurnLookBack:                 envvar.NewInt64("KeeperTurnLookBack").ParsePtr(),
		Registry: config.KeeperRegistry{
			CheckGasOverhead:    envvar.NewUint32("KeeperRegistryCheckGasOverhead").ParsePtr(),
			PerformGasOverhead:  envvar.NewUint32("KeeperRegistryPerformGasOverhead").ParsePtr(),
			MaxPerformDataSize:  envvar.NewUint32("KeeperRegistryMaxPerformDataSize").ParsePtr(),
			SyncInterval:        envDuration("KeeperRegistrySyncInterval"),
			SyncUpkeepQueueSize: envvar.KeeperRegistrySyncUpkeepQueueSize.ParsePtr(),
		},
	}

	c.AutoPprof = config.AutoPprof{
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

	c.Pyroscope = config.Pyroscope{
		ServerAddress: envvar.NewString("PyroscopeServerAddress").ParsePtr(),
		Environment:   envvar.NewString("PyroscopeEnvironment").ParsePtr(),
	}

	if dsn := os.Getenv("SENTRY_DSN"); dsn != "" {
		c.Sentry.DSN = &dsn
	}
	if debug := os.Getenv("SENTRY_DEBUG") == "true"; debug {
		c.Sentry.Debug = &debug
	}
	if env := os.Getenv("SENTRY_ENVIRONMENT"); env != "" {
		c.Sentry.Environment = &env
	}
	if rel := os.Getenv("SENTRY_RELEASE"); rel != "" {
		c.Sentry.Release = &rel
	}
	return nil
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
	ip := envvar.New(s, func(s string) (net.IP, error) {
		return net.ParseIP(s), nil
	}).ParsePtr()
	if ip != nil {
		// ensure non-zero
		for _, b := range *ip {
			if b != 0 {
				return ip
			}
		}
	}
	return nil
}

func envStringSlice(s string) *[]string {
	return envvar.New(s, func(s string) ([]string, error) {
		// matching viper cast.ToStringSliceE logic
		return strings.Fields(s), nil
	}).ParsePtr()
}

func envSlice[T any](s string, parse func(*T, []byte) error) *[]T {
	return envvar.New(s, func(v string) ([]T, error) {
		// matching viper cast.ToStringSliceE logic
		ss := strings.Fields(v)
		var ts []T
		for _, s := range ss {
			var t T
			err := parse(&t, []byte(s))
			if err != nil {
				return nil, err
			}
			ts = append(ts, t)
		}
		return ts, nil
	}).ParsePtr()
}

func mustParseDuration(s string) *models.Duration {
	if s == "" {
		return nil
	}
	d, err := models.ParseDuration(s)
	if err != nil {
		panic(fmt.Sprintf("invalid duration %s: %v", s, err))
	}
	return &d
}

func ptr[T any](t T) *T { return &t }
