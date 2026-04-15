package chainlink

import (
	"bytes"
	"context"
	stderrors "errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	otelpyroscope "github.com/grafana/otel-profiling-go"
	"github.com/grafana/pyroscope-go"
	"github.com/jonboulle/clockwork"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	nodeauthjwt "github.com/smartcontractkit/chainlink-common/pkg/nodeauth/jwt"
	commonsrv "github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/otelhealth"
	"github.com/smartcontractkit/chainlink-common/pkg/services/promhealth"
	commoncresettings "github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/dontime"
	"github.com/smartcontractkit/chainlink-data-streams/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/mercury"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	evmutils "github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/build"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockheaderfeeder"
	"github.com/smartcontractkit/chainlink/v2/core/services/ccv/ccvcommitteeverifier"
	"github.com/smartcontractkit/chainlink/v2/core/services/ccv/ccvexecutor"
	"github.com/smartcontractkit/chainlink/v2/core/services/cre"
	"github.com/smartcontractkit/chainlink/v2/core/services/cresettings"
	"github.com/smartcontractkit/chainlink/v2/core/services/cron"
	"github.com/smartcontractkit/chainlink/v2/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
	"github.com/smartcontractkit/chainlink/v2/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/headreporter"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keeper"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/retirement"
	"github.com/smartcontractkit/chainlink/v2/core/services/nodestatusreporter/bridgestatus"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/periodicbackup"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/services/ring"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf"
	"github.com/smartcontractkit/chainlink/v2/core/services/webhook"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	workflowstore "github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/sessions/ldapauth"
	"github.com/smartcontractkit/chainlink/v2/core/sessions/localauth"
	"github.com/smartcontractkit/chainlink/v2/core/sessions/oidcauth"
	"github.com/smartcontractkit/chainlink/v2/core/static"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

// Application implements the common functions used in the core node.
type Application interface {
	Start(ctx context.Context) error
	Stop() error
	GetLogger() logger.SugaredLogger
	GetAuditLogger() audit.AuditLogger
	GetHealthChecker() services.Checker
	GetDB() sqlutil.DataSource
	GetConfig() GeneralConfig
	SetLogLevel(lvl zapcore.Level) error
	GetKeyStore() keystore.Master
	WakeSessionReaper()
	GetWebAuthnConfiguration() sessions.WebAuthnConfiguration

	GetCapabilitiesRegistry() *capabilities.Registry

	GetExternalInitiatorManager() webhook.ExternalInitiatorManager
	GetRelayers() RelayerChainInteroperators
	GetLoopRegistry() *plugins.LoopRegistry
	GetLoopRegistrarConfig() plugins.RegistrarConfig

	// V2 Jobs (TOML specified)
	JobSpawner() job.Spawner
	JobORM() job.ORM
	PipelineORM() pipeline.ORM
	BridgeORM() bridges.ORM
	BasicAdminUsersORM() sessions.BasicAdminUsersORM
	AuthenticationProvider() sessions.AuthenticationProvider
	TxmStorageService() txmgr.EvmTxStore
	AddJobV2(ctx context.Context, job *job.Job) error
	DeleteJob(ctx context.Context, jobID int32) error
	RunWebhookJobV2(ctx context.Context, jobUUID uuid.UUID, requestBody string, meta jsonserializable.JSONSerializable) (int64, error)
	ResumeJobV2(ctx context.Context, taskID uuid.UUID, result pipeline.Result) error
	// Testing only
	RunJobV2(ctx context.Context, jobID int32, meta map[string]any) (int64, error)

	// Feeds
	GetFeedsService() feeds.Service

	// ReplayFromBlock replays logs from on or after the given block number. If forceBroadcast (evm only)
	// is set to true, consumers will reprocess data even if it has already been processed.
	ReplayFromBlock(ctx context.Context, chainFamily string, chainID string, number uint64, forceBroadcast bool) error

	// ID is unique to this particular application instance
	ID() uuid.UUID

	SecretGenerator() SecretGenerator

	// FindLCA - finds last common ancestor for LogPoller's chain available in the database and RPC chain
	FindLCA(ctx context.Context, chainID *big.Int) (*logpoller.Block, error)
	// DeleteLogPollerDataAfter - delete LogPoller state starting from the specified block
	DeleteLogPollerDataAfter(ctx context.Context, chainID *big.Int, start int64) error
}

// ChainlinkApplication contains fields for the JobSubscriber, Scheduler,
// and Store. The JobSubscriber and Scheduler are also available
// in the services package, but the Store has its own package.
type ChainlinkApplication struct {
	relayers                 *CoreRelayerChainInteroperators
	jobORM                   job.ORM
	jobSpawner               job.Spawner
	pipelineORM              pipeline.ORM
	pipelineRunner           pipeline.Runner
	bridgeORM                bridges.ORM
	localAdminUsersORM       sessions.BasicAdminUsersORM
	authenticationProvider   sessions.AuthenticationProvider // Note: this will be OIDC instance
	txmStorageService        txmgr.EvmTxStore
	FeedsService             feeds.Service
	webhookJobRunner         webhook.JobRunner
	Config                   GeneralConfig
	KeyStore                 keystore.Master
	ExternalInitiatorManager webhook.ExternalInitiatorManager
	SessionReaper            *utils.SleeperTask
	shutdownOnce             sync.Once
	srvcs                    []services.ServiceCtx
	HealthChecker            services.Checker
	logger                   logger.SugaredLogger
	AuditLogger              audit.AuditLogger
	closeLogger              func() error
	ds                       sqlutil.DataSource
	secretGenerator          SecretGenerator
	profiler                 *pyroscope.Profiler
	loopRegistry             *plugins.LoopRegistry
	loopRegistrarConfig      plugins.RegistrarConfig
	capabilitiesRegistry     *capabilities.Registry
	shardOrchestratorClient  shardorchestrator.ClientInterface

	started     bool
	startStopMu sync.Mutex
}

type ApplicationOpts struct {
	// CREOpts is the options for the CRE services
	cre.Opts

	Config                   GeneralConfig
	Logger                   logger.Logger
	Registerer               prometheus.Registerer
	DS                       sqlutil.DataSource
	KeyStore                 keystore.Master
	AuditLogger              audit.AuditLogger
	CloseLogger              func() error
	ExternalInitiatorManager webhook.ExternalInitiatorManager
	Version                  string
	VersionTag               string
	DockerTag                string
	RestrictedHTTPClient     *http.Client
	UnrestrictedHTTPClient   *http.Client
	SecretGenerator          SecretGenerator
	GRPCOpts                 loop.GRPCOpts
	MercuryPool              wsrpc.Pool
	RetirementReportCache    retirement.RetirementReportCache
	LLOTransmissionReaper    services.ServiceCtx
	NewOracleFactoryFn       standardcapabilities.NewOracleFactoryFn
	EVMFactoryConfigFn       func(*EVMFactoryConfig)
	DonTimeStore             *dontime.Store
}

// NewApplication initializes a new store if one is not already
// present at the configured root directory (default: ~/.chainlink),
// the logger at the same directory and returns the Application to
// be used by the node.
// TODO: Inject more dependencies here to save booting up useless stuff in tests
func NewApplication(ctx context.Context, opts ApplicationOpts) (Application, error) {
	var srvcs []services.ServiceCtx

	heartbeat := NewHeartbeat(NewHeartbeatConfig(opts))
	nodePlatformBuildInfo := NewNodePlatformBuildInfoService(NewNodePlatformBuildInfoConfig(opts))
	srvcs = append(srvcs, &heartbeat, &nodePlatformBuildInfo)

	auditLogger := opts.AuditLogger
	cfg := opts.Config
	externalInitiatorManager := opts.ExternalInitiatorManager
	globalLogger := logger.Sugared(opts.Logger)
	keyStore := opts.KeyStore
	restrictedHTTPClient := opts.RestrictedHTTPClient
	unrestrictedHTTPClient := opts.UnrestrictedHTTPClient

	mailMon := mailbox.NewMonitor(cfg.AppID().String(), globalLogger.Named("Mailbox"))

	if opts.CapabilitiesRegistry == nil {
		// for tests only, in prod Registry should always be set at this point
		opts.CapabilitiesRegistry = capabilities.NewRegistry(globalLogger)
	}

	if opts.DonTimeStore == nil {
		opts.DonTimeStore = dontime.NewStore(dontime.DefaultRequestTimeout)
	}

	var ringStoreForShard0 *ring.Store
	var shardOrchestratorClient shardorchestrator.ClientInterface
	if cfg.Sharding().ShardingEnabled() {
		shardIdx := cfg.Sharding().ShardIndex()
		if shardIdx == 0 {
			ringStoreForShard0 = ring.NewStore()
			server := shardorchestrator.NewServer(ringStoreForShard0, globalLogger)
			shardOrchestratorClient = shardorchestrator.NewLocalClient(server, globalLogger)
			globalLogger.Infow("ShardOrchestrator in-process client created", "shardID", shardIdx)
		} else {
			shardOrchestratorAddr := cfg.Sharding().ShardOrchestratorAddress()
			if shardOrchestratorAddr == nil {
				return nil, fmt.Errorf("shard %d requires ShardOrchestratorAddress when sharding is enabled", shardIdx)
			}
			client, err := shardorchestrator.NewClient(shardOrchestratorAddr.String(), globalLogger.Named("ShardOrchestratorClient"))
			if err != nil {
				return nil, fmt.Errorf("failed to create ShardOrchestrator gRPC client: %w", err)
			}
			shardOrchestratorClient = client
			globalLogger.Infow("ShardOrchestrator gRPC client created", "shardID", shardIdx, "serverAddress", shardOrchestratorAddr.String())
		}
	} else {
		globalLogger.Debug("Sharding not enabled, running without shard orchestrator client")
	}

	creSettingsTOML, err := toml.Marshal(commoncresettings.Default)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal cre settings TOML: %w", err)
	}
	globalLogger.Debugf("# CRESettings defaults: \n%s", creSettingsTOML)
	atomicSettings := loop.NewAtomicSettings(commoncresettings.DefaultGetter)
	limitsFactory := limits.Factory{
		Meter:    beholder.GetMeter(),
		Logger:   globalLogger.Named("Limits"),
		Settings: atomicSettings,
	}

	csaKeystore := &keystore.CSASigner{CSA: keyStore.CSA()}
	beholderAuthHeaders, csaPubKeyHex, err := keystore.BuildBeholderAuth(ctx, keyStore.CSA())
	if err != nil {
		return nil, fmt.Errorf("failed to build Beholder auth: %w", err)
	}
	loopRegistry := plugins.NewLoopRegistry(globalLogger, cfg.AppID().String(), cfg.Feature().LogPoller(),
		cfg.Database(), cfg.Mercury(), cfg.Pyroscope(), cfg.AutoPprof(), cfg.Tracing(), cfg.Telemetry(),
		beholderAuthHeaders, csaPubKeyHex, cfg.LOOPP())

	relayerFactory := RelayerFactory{
		Logger:                opts.Logger,
		Registerer:            opts.Registerer,
		LoopRegistry:          loopRegistry,
		GRPCOpts:              opts.GRPCOpts,
		MercuryPool:           opts.MercuryPool,
		CapabilitiesRegistry:  opts.CapabilitiesRegistry,
		HTTPClient:            opts.UnrestrictedHTTPClient,
		RetirementReportCache: opts.RetirementReportCache,
	}

	evmFactoryCfg := EVMFactoryConfig{
		ChainOpts: legacyevm.ChainOpts{
			ChainConfigs:   cfg.EVMConfigs(),
			DatabaseConfig: cfg.Database(),
			ListenerConfig: cfg.Database().Listener(),
			FeatureConfig:  cfg.Feature(),
			MailMon:        mailMon,
			DS:             opts.DS,
		},
		EthKeystore:   keyStore.Eth(),
		CSAKeystore:   csaKeystore,
		MercuryConfig: cfg.Mercury(),
	}

	if opts.EVMFactoryConfigFn != nil {
		opts.EVMFactoryConfigFn(&evmFactoryCfg)
	}

	// evm always enabled for backward compatibility
	// TODO BCF-2510 this needs to change in order to clear the path for EVM extraction
	initOps := []CoreRelayerChainInitFunc{InitDummy(relayerFactory), InitEVM(relayerFactory, evmFactoryCfg)}

	if cfg.CosmosEnabled() {
		initOps = append(initOps, InitCosmos(relayerFactory, keyStore.Cosmos(), keyStore.CSA(), cfg.CosmosConfigs()))
	}
	if cfg.SolanaEnabled() {
		solanaCfg := SolanaFactoryConfig{
			TOMLConfigs: cfg.SolanaConfigs(),
			DS:          opts.DS,
		}
		initOps = append(initOps, InitSolana(relayerFactory, keyStore.Solana(), keyStore.CSA(), solanaCfg))
	}
	if cfg.StarkNetEnabled() {
		initOps = append(initOps, InitStarknet(relayerFactory, keyStore.StarkNet(), keyStore.CSA(), cfg.StarknetConfigs()))
	}
	if cfg.AptosEnabled() {
		initOps = append(initOps, InitAptos(relayerFactory, keyStore.Aptos(), keyStore.CSA(), cfg.AptosConfigs()))
	}
	if cfg.TronEnabled() {
		initOps = append(initOps, InitTron(relayerFactory, keyStore.Tron(), keyStore.CSA(), cfg.TronConfigs()))
	}
	if cfg.TONEnabled() {
		initOps = append(initOps, InitTON(relayerFactory, keyStore.TON(), keyStore.CSA(), cfg.TONConfigs()))
	}
	if cfg.SuiEnabled() {
		initOps = append(initOps, InitSui(relayerFactory, keyStore.Sui(), keyStore.CSA(), cfg.SuiConfigs()))
	}

	relayChainInterops, err := NewCoreRelayerChainInteroperators(initOps...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize relayer chain interoperators: %w", err)
	}

	var peerWrapper *ocrcommon.SingletonPeerWrapper
	if !cfg.OCR().Enabled() && !cfg.OCR2().Enabled() {
		globalLogger.Debug("P2P stack not needed")
	} else {
		if !cfg.P2P().Enabled() {
			return nil, errors.New("P2P stack required for OCR or OCR2")
		}
		if err2 := ocrcommon.ValidatePeerWrapperConfig(cfg.P2P()); err != nil {
			return nil, fmt.Errorf("invalid P2P config: %w", err2)
		}
		peerWrapper = ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), opts.DS, globalLogger)
		srvcs = append(srvcs, peerWrapper)
	}

	workflowKey, err := keystore.GetDefault(ctx, keyStore.Workflow())
	if err != nil {
		return nil, fmt.Errorf("failed to get default workflow key: %w", err)
	}

	csaSigner, csaPubKey, err := keystore.BuildNodeAuth(ctx, csaKeystore)
	if err != nil {
		return nil, fmt.Errorf("failed to build node auth: %w", err)
	}
	jwtGenerator := nodeauthjwt.NewNodeJWTGenerator(csaSigner, csaPubKey)

	creServices, err := cre.NewServices(
		globalLogger,
		opts.DS,
		keyStore,
		cfg,
		relayChainInterops,
		peerWrapper,
		cre.Opts{
			CapabilitiesRegistry:    opts.CapabilitiesRegistry,
			CapabilitiesDispatcher:  opts.CapabilitiesDispatcher,
			CapabilitiesPeerWrapper: opts.CapabilitiesPeerWrapper,
			FetcherFunc:             opts.FetcherFunc,
			FetcherFactoryFn:        opts.FetcherFactoryFn,
			BillingClient:           opts.BillingClient,
			LinkingClient:           opts.LinkingClient,
			StorageClient:           opts.StorageClient,
			DonTimeStore:            opts.DonTimeStore,
			LimitsFactory:           limitsFactory,
			UseLocalTimeProvider:    opts.UseLocalTimeProvider,
			WorkflowKey:             workflowKey,
			JWTGenerator:            jwtGenerator,
			ShardOrchestratorClient: shardOrchestratorClient,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initilize CRE: %w", err)
	}
	srvcs = append(srvcs, creServices)

	// If the audit logger is enabled
	if auditLogger.Ready() == nil {
		srvcs = append(srvcs, auditLogger)
	}

	var profiler *pyroscope.Profiler
	if cfg.Pyroscope().ServerAddress() != "" {
		globalLogger.Debug("Pyroscope (automatic pprof profiling) is enabled")
		var err error
		profiler, err = logger.StartPyroscope(cfg.Pyroscope(), cfg.AutoPprof())
		if err != nil {
			return nil, errors.Wrap(err, "starting pyroscope (automatic pprof profiling) failed")
		}

		if cfg.Pyroscope().LinkTracesToProfiles() && cfg.Tracing().Enabled() {
			// Enable span profiling - link OTel traces to Pyroscope profiles
			// This wraps the global tracer provider to record span IDs alongside profile samples
			otel.SetTracerProvider(
				otelpyroscope.NewTracerProvider(
					otel.GetTracerProvider(),
					otelpyroscope.WithAppName("chainlink-node"),
					otelpyroscope.WithPyroscopeURL(cfg.Pyroscope().ServerAddress()),
				),
			)
		}
	} else {
		globalLogger.Debug("Pyroscope (automatic pprof profiling) is disabled")
	}

	ap := cfg.AutoPprof()
	if ap.Enabled() {
		globalLogger.Info("Nurse service (automatic pprof profiling) is enabled")
		srvcs = append(srvcs, services.NewNurse(ap, globalLogger))
	} else {
		globalLogger.Info("Nurse service (automatic pprof profiling) is disabled")
	}

	telemetryManager := telemetry.NewManager(cfg.TelemetryIngress(), csaKeystore, globalLogger)
	srvcs = append(srvcs, telemetryManager)

	backupCfg := cfg.Database().Backup()
	if backupCfg.Mode() != config.DatabaseBackupModeNone && backupCfg.Frequency() > 0 {
		globalLogger.Infow("DatabaseBackup: periodic database backups are enabled", "frequency", backupCfg.Frequency())

		databaseBackup, err := periodicbackup.NewDatabaseBackup(cfg.Database().URL(), cfg.RootDir(), backupCfg, globalLogger)
		if err != nil {
			return nil, errors.Wrap(err, "NewApplication: failed to initialize database backup")
		}
		srvcs = append(srvcs, databaseBackup)
	} else {
		globalLogger.Info("DatabaseBackup: periodic database backups are disabled. To enable automatic backups, set Database.Backup.Mode=lite or Database.Backup.Mode=full")
	}

	// pool must be started before all relayers and stopped after them
	if opts.MercuryPool != nil {
		srvcs = append(srvcs, opts.MercuryPool)
	}
	if opts.RetirementReportCache != nil {
		srvcs = append(srvcs, opts.RetirementReportCache)
	}
	if opts.LLOTransmissionReaper != nil {
		srvcs = append(srvcs, opts.LLOTransmissionReaper)
	}

	// EVM chains are used all over the place. This will need to change for fully EVM extraction
	// TODO: BCF-2510, BCF-2511

	legacyEVMChains := relayChainInterops.LegacyEVMChains()
	if legacyEVMChains == nil {
		return nil, errors.New("no evm chains found")
	}

	srvcs = append(srvcs, mailMon)
	srvcs = append(srvcs, relayChainInterops.Services()...)

	// Initialize Local Users ORM and Authentication Provider specified in config
	// BasicAdminUsersORM is initialized and required regardless of separate Authentication Provider
	localAdminUsersORM := localauth.NewORM(opts.DS, cfg.WebServer().SessionTimeout().Duration(), globalLogger, auditLogger)

	// Initialize Sessions ORM based on environment configured authenticator
	// localDB auth, LDAP auth, or OIDC auth
	authMethod := cfg.WebServer().AuthenticationMethod()
	var authenticationProvider sessions.AuthenticationProvider
	var sessionReaper *utils.SleeperTask

	switch sessions.AuthenticationProviderName(authMethod) {
	case sessions.LDAPAuth:
		var err error
		authenticationProvider, err = ldapauth.NewLDAPAuthenticator(
			opts.DS, cfg.WebServer().LDAP(), cfg.Insecure().DevWebServer(), globalLogger, auditLogger,
		)
		if err != nil {
			return nil, errors.Wrap(err, "NewApplication: failed to initialize LDAP Authentication module")
		}
		syncer := ldapauth.NewLDAPServerStateSyncer(opts.DS, cfg.WebServer().LDAP(), globalLogger)
		srvcs = append(srvcs, syncer)
		sessionReaper = utils.NewSleeperTaskCtx(syncer)
	case sessions.OIDCAuth:
		var err error
		authenticationProvider, err = oidcauth.NewOIDCAuthenticator(
			opts.DS, cfg.WebServer().OIDC(), globalLogger, auditLogger,
		)
		if err != nil {
			return nil, errors.Wrap(err, "NewApplication: failed to initialize OIDC Authentication module")
		}
		sessionReaper = oidcauth.NewSessionReaper(opts.DS, cfg.WebServer(), globalLogger)
	case sessions.LocalAuth:
		authenticationProvider = localauth.NewORM(opts.DS, cfg.WebServer().SessionTimeout().Duration(), globalLogger, auditLogger)
		sessionReaper = localauth.NewSessionReaper(opts.DS, cfg.WebServer(), globalLogger)
	default:
		return nil, errors.Errorf("NewApplication: Unexpected 'AuthenticationMethod': %s supported values: %s, %s", authMethod, sessions.LocalAuth, sessions.LDAPAuth)
	}

	var (
		pipelineORM    = pipeline.NewORM(opts.DS, globalLogger, cfg.JobPipeline().MaxSuccessfulRuns())
		bridgeORM      = bridges.NewORM(opts.DS)
		mercuryORM     = mercury.NewORM(opts.DS)
		pipelineRunner = pipeline.NewRunner(pipelineORM, bridgeORM, cfg.JobPipeline(), cfg.WebServer(), legacyEVMChains, keyStore.Eth(), keyStore.VRF(), globalLogger, restrictedHTTPClient, unrestrictedHTTPClient)
		jobORM         = job.NewORM(opts.DS, pipelineORM, bridgeORM, keyStore, globalLogger)
		txmORM         = txmgr.NewTxStore(opts.DS, globalLogger)
		streamRegistry = streams.NewRegistry(globalLogger, pipelineRunner)
		workflowORM    = workflowstore.NewInMemoryStore(globalLogger, clockwork.NewRealClock())
	)
	srvcs = append(srvcs, workflowORM)

	promReporter := headreporter.NewLegacyEVMPrometheusReporter(opts.DS, legacyEVMChains)
	evmChainIDs := make([]*big.Int, len(cfg.EVMConfigs()))
	for i, chain := range cfg.EVMConfigs() {
		evmChainIDs[i] = chain.ChainID.ToInt()
	}

	legacyEVMTelemReporter := headreporter.NewLegacyEVMTelemetryReporter(telemetryManager, globalLogger, evmChainIDs...)
	loopTelemReporter := headreporter.NewTelemetryReporter(telemetryManager, globalLogger, relayChainInterops.GetIDToRelayerMap())
	headReporter := headreporter.NewHeadReporterService(opts.DS, globalLogger, promReporter, legacyEVMTelemReporter, loopTelemReporter)
	srvcs = append(srvcs, headReporter)
	for _, chain := range legacyEVMChains.Slice() {
		legacyChain, ok := chain.(legacyevm.Chain)
		if !ok {
			continue
		}
		legacyChain.HeadBroadcaster().Subscribe(headReporter)
		legacyChain.TxManager().RegisterResumeCallback(pipelineRunner.ResumeRun)
	}

	srvcs = append(srvcs, pipelineORM)

	loopRegistrarConfig := plugins.NewRegistrarConfig(opts.GRPCOpts, loopRegistry.Register, loopRegistry.Unregister)

	var (
		delegates = map[job.Type]job.Delegate{
			job.DirectRequest: directrequest.NewDelegate(
				globalLogger,
				pipelineRunner,
				pipelineORM,
				legacyEVMChains,
				mailMon),
			job.Keeper: keeper.NewDelegate(
				cfg,
				opts.DS,
				jobORM,
				pipelineRunner,
				globalLogger,
				legacyEVMChains,
				mailMon),
			job.VRF: vrf.NewDelegate(
				opts.DS,
				keyStore,
				pipelineRunner,
				pipelineORM,
				legacyEVMChains,
				globalLogger,
				mailMon),
			job.Webhook: webhook.NewDelegate(
				pipelineRunner,
				externalInitiatorManager,
				globalLogger),
			job.Cron: cron.NewDelegate(
				pipelineRunner,
				globalLogger),
			job.BlockhashStore: blockhashstore.NewDelegate(
				cfg,
				globalLogger,
				legacyEVMChains,
				keyStore.Eth()),
			job.BlockHeaderFeeder: blockheaderfeeder.NewDelegate(
				cfg,
				globalLogger,
				legacyEVMChains,
				keyStore.Eth()),
			job.Gateway: gateway.NewDelegate(
				legacyEVMChains,
				keyStore.Eth(),
				opts.DS,
				opts.CapabilitiesRegistry,
				creServices.WorkflowRegistrySyncer,
				globalLogger,
				limitsFactory,
			),
			job.Stream: streams.NewDelegate(
				globalLogger,
				streamRegistry,
				pipelineRunner,
				cfg.JobPipeline(),
			),
			job.CCVCommitteeVerifier: ccvcommitteeverifier.NewDelegate(
				globalLogger,
				opts.DS,
				cfg.CCV(),
				keyStore.OCR2(),
				relayChainInterops.LegacyEVMChains().Slice(),
			),
			job.CCVExecutor: ccvexecutor.NewDelegate(
				globalLogger,
				cfg.CCV(),
				keyStore.Eth(),
				relayChainInterops.LegacyEVMChains().Slice(),
			),
		}
		webhookJobRunner = delegates[job.Webhook].(*webhook.Delegate).WebhookJobRunner()
	)

	delegates[job.Workflow] = workflows.NewDelegate(
		globalLogger,
		opts.CapabilitiesRegistry,
		opts.DonTimeStore,
		workflowORM,
		creServices.WorkflowRateLimiter,
		creServices.WorkflowLimits,
		workflows.WithBillingClient(creServices.BillingClient),
		workflows.WithWorkflowRegistry(cfg.Capabilities().WorkflowRegistry().Address(), cfg.Capabilities().WorkflowRegistry().ChainID()),
	)

	// Flux monitor requires ethereum just to boot, silence errors with a null delegate
	if !cfg.EVMConfigs().RPCEnabled() {
		delegates[job.FluxMonitor] = &job.NullDelegate{Type: job.FluxMonitor}
	} else {
		delegates[job.FluxMonitor] = fluxmonitorv2.NewDelegate(
			cfg,
			keyStore.Eth(),
			jobORM,
			pipelineORM,
			pipelineRunner,
			opts.DS,
			legacyEVMChains,
			globalLogger,
		)
	}

	delegates[job.CRESettings] = cresettings.NewDelegate(globalLogger, atomicSettings)

	// If peer wrapper is initialized, Oracle Factory dependency will be available to standard capabilities
	stdcapDelegate := standardcapabilities.NewDelegate(
		globalLogger,
		opts.DS, jobORM,
		opts.CapabilitiesRegistry,
		loopRegistrarConfig,
		telemetryManager,
		pipelineRunner,
		relayChainInterops,
		creServices.GatewayConnectorWrapper,
		keyStore,
		creServices.GetPeerID,
		peerWrapper,
		opts.NewOracleFactoryFn,
		opts.FetcherFactoryFn,
		creServices.OrgResolver,
		atomicSettings,
		creServices.OCRConfigService,
		cfg.Capabilities().Local(),
	)
	delegates[job.StandardCapabilities] = stdcapDelegate
	if creServices.SetDelegatesDeps != nil {
		depSvc, depErr := creServices.SetDelegatesDeps(stdcapDelegate)
		if depErr != nil {
			return nil, fmt.Errorf("failed to set CRE delegates dependencies: %w", depErr)
		}
		if depSvc != nil {
			srvcs = append(srvcs, depSvc)
		}
	}

	if cfg.OCR().Enabled() {
		delegates[job.OffchainReporting] = ocr.NewDelegate(
			opts.DS,
			jobORM,
			keyStore.Eth(),
			keyStore.OCR(),
			pipelineRunner,
			peerWrapper,
			telemetryManager,
			legacyEVMChains,
			globalLogger,
			cfg,
			mailMon,
		)
	} else {
		globalLogger.Debug("Off-chain reporting disabled")
	}

	if cfg.OCR2().Enabled() {
		globalLogger.Debug("Off-chain reporting v2 enabled")

		ocr2DelegateConfig := ocr2.NewDelegateConfig(cfg.OCR2(), cfg.Mercury(), cfg.Threshold(), cfg.Insecure(), cfg.JobPipeline(), loopRegistrarConfig, cfg.Sharding(), ringStoreForShard0)

		ocr2Delegate := ocr2.NewDelegate(
			ocr2.DelegateOpts{
				Ds:                             opts.DS,
				JobORM:                         jobORM,
				BridgeORM:                      bridgeORM,
				MercuryORM:                     mercuryORM,
				PipelineRunner:                 pipelineRunner,
				StreamRegistry:                 streamRegistry,
				PeerWrapper:                    peerWrapper,
				MonitoringEndpointGen:          telemetryManager,
				LegacyChains:                   legacyEVMChains,
				Lggr:                           globalLogger,
				Ks:                             keyStore.OCR2(),
				EthKs:                          keyStore.Eth(),
				WorkflowKs:                     keyStore.Workflow(),
				DKGRecipientKs:                 keyStore.DKGRecipient(),
				Relayers:                       relayChainInterops,
				MailMon:                        mailMon,
				CapabilitiesRegistry:           opts.CapabilitiesRegistry,
				DonTimeStore:                   opts.DonTimeStore,
				RetirementReportCache:          opts.RetirementReportCache,
				GatewayConnectorServiceWrapper: creServices.GatewayConnectorWrapper,
				WorkflowRegistrySyncer:         creServices.WorkflowRegistrySyncer,
				OrgResolver:                    creServices.OrgResolver,
				LimitsFactory:                  limitsFactory,
				OCRConfigService:               creServices.OCRConfigService,
			},
			ocr2DelegateConfig,
		)
		if ocr2Delegate == nil {
			return nil, errors.New("ocr2.NewDelegate() returned nil")
		}
		delegates[job.OffchainReporting2] = ocr2Delegate
		delegates[job.Bootstrap] = ocrbootstrap.NewDelegateBootstrap(
			opts.DS,
			jobORM,
			peerWrapper,
			globalLogger,
			cfg.OCR2(),
			cfg.Insecure(),
			relayChainInterops,
		)
		delegates[job.CCIP] = ccip.NewDelegate(
			globalLogger,
			loopRegistrarConfig,
			pipelineRunner,
			relayChainInterops,
			opts.KeyStore,
			opts.DS,
			peerWrapper,
			telemetryManager,
			cfg.Capabilities(),
			cfg.EVMConfigs(),
		)
	} else {
		globalLogger.Debug("Off-chain reporting v2 disabled")
	}

	bridgeStatusReporter := bridgestatus.NewBridgeStatusReporter(
		cfg.BridgeStatusReporter(),
		bridgeORM,
		jobORM,
		unrestrictedHTTPClient,
		beholder.GetEmitter(),
		globalLogger,
	)
	srvcs = append(srvcs, bridgeStatusReporter)

	healthCfg := commonsrv.HealthCheckerConfig{Ver: static.Version, Sha: static.Sha}
	healthCfg = promhealth.ConfigureHooks(healthCfg)
	healthCfg, err = otelhealth.ConfigureHooks(healthCfg, beholder.GetMeter())
	if err != nil {
		return nil, fmt.Errorf("failed to configure health checker otel hooks: %w", err)
	}
	healthChecker := healthCfg.New()

	var lbs []utils.DependentAwaiter
	for _, c := range legacyEVMChains.Slice() {
		legacyChain, ok := c.(legacyevm.Chain)
		if !ok {
			continue
		}
		lbs = append(lbs, legacyChain.LogBroadcaster())
	}
	jobSpawner := job.NewSpawner(jobORM, cfg.Database(), healthChecker, delegates, globalLogger, lbs)
	srvcs = append(srvcs, jobSpawner, pipelineRunner)

	var feedsService feeds.Service
	if cfg.Feature().FeedsManager() {
		feedsORM := feeds.NewORM(opts.DS, globalLogger)
		feedsService = feeds.NewService(
			feedsORM,
			jobORM,
			opts.DS,
			jobSpawner,
			keyStore,
			cfg,
			cfg.JobDistributor(),
			cfg.Feature(),
			cfg.Insecure(),
			cfg.JobPipeline(),
			cfg.OCR(),
			cfg.OCR2(),
			legacyEVMChains,
			globalLogger,
			opts.VersionTag,
			loopRegistrarConfig,
		)
	} else {
		feedsService = &feeds.NullService{}
	}

	for _, s := range srvcs {
		if s == nil {
			panic("service unexpectedly nil")
		}
		if err := healthChecker.Register(s); err != nil {
			return nil, fmt.Errorf("failed to register health check for service %T: %w", s, err)
		}
	}

	return &ChainlinkApplication{
		relayers:                 relayChainInterops,
		jobORM:                   jobORM,
		jobSpawner:               jobSpawner,
		pipelineRunner:           pipelineRunner,
		pipelineORM:              pipelineORM,
		bridgeORM:                bridgeORM,
		localAdminUsersORM:       localAdminUsersORM,
		authenticationProvider:   authenticationProvider,
		txmStorageService:        txmORM,
		FeedsService:             feedsService,
		Config:                   cfg,
		webhookJobRunner:         webhookJobRunner,
		KeyStore:                 keyStore,
		SessionReaper:            sessionReaper,
		ExternalInitiatorManager: externalInitiatorManager,
		HealthChecker:            healthChecker,
		logger:                   globalLogger,
		AuditLogger:              auditLogger,
		closeLogger:              opts.CloseLogger,
		secretGenerator:          opts.SecretGenerator,
		profiler:                 profiler,
		loopRegistry:             loopRegistry,
		loopRegistrarConfig:      loopRegistrarConfig,
		capabilitiesRegistry:     opts.CapabilitiesRegistry,
		ds:                       opts.DS,
		shardOrchestratorClient:  shardOrchestratorClient,

		srvcs: srvcs,
	}, nil
}

func (app *ChainlinkApplication) SetLogLevel(lvl zapcore.Level) error {
	if err := app.Config.SetLogLevel(lvl); err != nil {
		return err
	}
	app.logger.SetLogLevel(lvl)
	return nil
}

// Start all necessary services. If successful, nil will be returned.
// Start sequence is aborted if the context gets cancelled.
func (app *ChainlinkApplication) Start(ctx context.Context) error {
	app.startStopMu.Lock()
	defer app.startStopMu.Unlock()
	if app.started {
		panic("application is already started")
	}

	var span trace.Span
	ctx, span = otel.Tracer("").Start(ctx, "Start", trace.WithAttributes(
		attribute.String("app-id", app.ID().String()),
		attribute.String("version", static.Version),
		attribute.String("commit", static.Sha),
	))
	defer span.End()

	if app.FeedsService != nil {
		if err := app.FeedsService.Start(ctx); err != nil {
			app.logger.Errorf("[Feeds Service] Failed to start %v", err)
			app.FeedsService = &feeds.NullService{} // so we don't try to Close() later
		}
	}

	var ms services.MultiStart
	for _, service := range app.srvcs {
		if ctx.Err() != nil {
			err := errors.Wrap(ctx.Err(), "aborting start")
			return stderrors.Join(err, ms.Close())
		}

		app.logger.Infow("Starting service...", "name", service.Name())

		if err := ms.Start(ctx, service); err != nil {
			return err
		}
	}

	// Start HealthChecker last, so that the other services had the chance to
	// start enough to immediately pass the readiness check.
	if err := app.HealthChecker.Start(); err != nil {
		return err
	}

	app.started = true

	return nil
}

func (app *ChainlinkApplication) StopIfStarted() error {
	app.startStopMu.Lock()
	defer app.startStopMu.Unlock()
	if app.started {
		return app.stop()
	}
	return nil
}

func (app *ChainlinkApplication) GetLoopRegistry() *plugins.LoopRegistry {
	return app.loopRegistry
}

func (app *ChainlinkApplication) GetLoopRegistrarConfig() plugins.RegistrarConfig {
	return app.loopRegistrarConfig
}

// Stop allows the application to exit by halting schedules, closing
// logs, and closing the DB connection.
func (app *ChainlinkApplication) Stop() error {
	app.startStopMu.Lock()
	defer app.startStopMu.Unlock()
	return app.stop()
}

func (app *ChainlinkApplication) stop() (err error) {
	if !app.started {
		panic("application is already stopped")
	}
	app.shutdownOnce.Do(func() {
		shutdownStart := time.Now()
		defer func() {
			if app.closeLogger == nil {
				return
			}
			if lerr := app.closeLogger(); lerr != nil {
				err = stderrors.Join(err, lerr)
			}
		}()
		app.logger.Info("Gracefully exiting...")

		// Stop services in the reverse order from which they were started
		for i := len(app.srvcs) - 1; i >= 0; i-- {
			service := app.srvcs[i]
			app.logger.Debugw("Closing service...", "name", service.Name())
			err = stderrors.Join(err, service.Close())
		}

		app.logger.Debug("Stopping SessionReaper...")
		err = stderrors.Join(err, app.SessionReaper.Stop())
		app.logger.Debug("Closing HealthChecker...")
		err = stderrors.Join(err, app.HealthChecker.Close())
		if app.FeedsService != nil {
			app.logger.Debug("Closing Feeds Service...")
			err = stderrors.Join(err, app.FeedsService.Close())
		}

		if app.profiler != nil {
			err = stderrors.Join(err, app.profiler.Stop())
		}

		app.logger.Debugf("Closed application in %v", time.Since(shutdownStart))

		app.started = false
	})
	return err
}

func (app *ChainlinkApplication) GetConfig() GeneralConfig {
	return app.Config
}

func (app *ChainlinkApplication) GetKeyStore() keystore.Master {
	return app.KeyStore
}

func (app *ChainlinkApplication) GetLogger() logger.SugaredLogger {
	return app.logger
}

func (app *ChainlinkApplication) GetAuditLogger() audit.AuditLogger {
	return app.AuditLogger
}

func (app *ChainlinkApplication) GetHealthChecker() services.Checker {
	return app.HealthChecker
}

func (app *ChainlinkApplication) JobSpawner() job.Spawner {
	return app.jobSpawner
}

func (app *ChainlinkApplication) JobORM() job.ORM {
	return app.jobORM
}

func (app *ChainlinkApplication) BridgeORM() bridges.ORM {
	return app.bridgeORM
}

func (app *ChainlinkApplication) BasicAdminUsersORM() sessions.BasicAdminUsersORM {
	return app.localAdminUsersORM
}

func (app *ChainlinkApplication) AuthenticationProvider() sessions.AuthenticationProvider {
	return app.authenticationProvider
}

func (app *ChainlinkApplication) PipelineORM() pipeline.ORM {
	return app.pipelineORM
}

func (app *ChainlinkApplication) TxmStorageService() txmgr.EvmTxStore {
	return app.txmStorageService
}

func (app *ChainlinkApplication) GetExternalInitiatorManager() webhook.ExternalInitiatorManager {
	return app.ExternalInitiatorManager
}

func (app *ChainlinkApplication) GetCapabilitiesRegistry() *capabilities.Registry {
	return app.capabilitiesRegistry
}

func (app *ChainlinkApplication) SecretGenerator() SecretGenerator {
	return app.secretGenerator
}

// WakeSessionReaper wakes up the reaper to do its reaping.
func (app *ChainlinkApplication) WakeSessionReaper() {
	app.SessionReaper.WakeUp()
}

func (app *ChainlinkApplication) AddJobV2(ctx context.Context, j *job.Job) error {
	return app.jobSpawner.CreateJob(ctx, nil, j)
}

func (app *ChainlinkApplication) DeleteJob(ctx context.Context, jobID int32) error {
	// Do not allow the job to be deleted if it is managed by the Feeds Manager
	isManaged, err := app.FeedsService.IsJobManaged(ctx, int64(jobID))
	if err != nil {
		return err
	}

	if isManaged {
		return errors.New("job must be deleted in the feeds manager")
	}

	return app.jobSpawner.DeleteJob(ctx, nil, jobID)
}

func (app *ChainlinkApplication) RunWebhookJobV2(ctx context.Context, jobUUID uuid.UUID, requestBody string, meta jsonserializable.JSONSerializable) (int64, error) {
	return app.webhookJobRunner.RunJob(ctx, jobUUID, requestBody, meta)
}

// Only used for local testing, not supported by the UI.
func (app *ChainlinkApplication) RunJobV2(
	ctx context.Context,
	jobID int32,
	meta map[string]any,
) (int64, error) {
	if build.IsProd() {
		return 0, errors.New("manual job runs not supported on secure builds")
	}
	jb, err := app.jobORM.FindJob(ctx, jobID)
	if err != nil {
		return 0, errors.Wrapf(err, "job ID %v", jobID)
	}
	var runID int64

	// Some jobs are special in that they do not have a task graph.
	isBootstrap := jb.Type == job.OffchainReporting && jb.OCROracleSpec != nil && jb.OCROracleSpec.IsBootstrapPeer
	if jb.Type.RequiresPipelineSpec() || !isBootstrap {
		var vars map[string]any
		var saveTasks bool
		if jb.Type == job.VRF {
			saveTasks = true
			// Create a dummy log to trigger a run
			testLog := types.Log{
				Data: bytes.Join([][]byte{
					jb.VRFSpec.PublicKey.MustHash().Bytes(),  // key hash
					common.BigToHash(big.NewInt(42)).Bytes(), // seed
					evmutils.NewHash().Bytes(),               // sender
					evmutils.NewHash().Bytes(),               // fee
					evmutils.NewHash().Bytes(),
				}, // requestID
					[]byte{}),
				Topics:      []common.Hash{{}, jb.ExternalIDEncodeBytesToTopic()}, // jobID BYTES
				TxHash:      evmutils.NewHash(),
				BlockNumber: 10,
				BlockHash:   evmutils.NewHash(),
			}
			vars = map[string]any{
				"jobSpec": map[string]any{
					"databaseID":    jb.ID,
					"externalJobID": jb.ExternalJobID,
					"name":          jb.Name.ValueOrZero(),
					"publicKey":     jb.VRFSpec.PublicKey[:],
					"evmChainID":    jb.VRFSpec.EVMChainID.String(),
				},
				"jobRun": map[string]any{
					"meta":           meta,
					"logBlockHash":   testLog.BlockHash[:],
					"logBlockNumber": testLog.BlockNumber,
					"logTxHash":      testLog.TxHash,
					"logTopics":      testLog.Topics,
					"logData":        testLog.Data,
				},
			}
		} else {
			vars = map[string]any{
				"jobRun": map[string]any{
					"meta": meta,
				},
			}
		}
		runID, _, err = app.pipelineRunner.ExecuteAndInsertFinishedRun(ctx, *jb.PipelineSpec, pipeline.NewVarsFrom(vars), saveTasks)
	}
	return runID, err
}

func (app *ChainlinkApplication) ResumeJobV2(
	ctx context.Context,
	taskID uuid.UUID,
	result pipeline.Result,
) error {
	return app.pipelineRunner.ResumeRun(ctx, taskID, result.Value, result.Error)
}

func (app *ChainlinkApplication) GetFeedsService() feeds.Service {
	return app.FeedsService
}

// ReplayFromBlock implements the Application interface.
func (app *ChainlinkApplication) ReplayFromBlock(ctx context.Context, chainFamily string, chainID string, number uint64, forceBroadcast bool) error {
	if chainFamily == relay.NetworkEVM {
		// TODO: Implement EVM Replay on Relayer instead of using LegacyChains - BCFR-1160
		chain, err := app.GetRelayers().LegacyEVMChains().Get(chainID)
		if err != nil {
			return err
		}
		//nolint:gosec // this won't overflow
		fromBlock := int64(number)

		if legacyChain, ok := chain.(legacyevm.Chain); ok {
			legacyChain.LogBroadcaster().ReplayFromBlock(fromBlock, forceBroadcast)
			if app.Config.Feature().LogPoller() {
				legacyChain.LogPoller().ReplayAsync(fromBlock)
			}
			return nil
		}
		// else LOOPP mode, so fall back to default
	}
	relayer, err := app.GetRelayers().Get(commontypes.RelayID{
		Network: chainFamily,
		ChainID: chainID,
	})
	if err != nil {
		return err
	}
	return relayer.Replay(ctx, strconv.FormatUint(number, 10), map[string]any{})
}

func (app *ChainlinkApplication) GetRelayers() RelayerChainInteroperators {
	return app.relayers
}

func (app *ChainlinkApplication) GetDB() sqlutil.DataSource {
	return app.ds
}

// Returns the configuration to use for creating and authenticating
// new WebAuthn credentials
func (app *ChainlinkApplication) GetWebAuthnConfiguration() sessions.WebAuthnConfiguration {
	rpid := app.Config.WebServer().MFA().RPID()
	rporigin := app.Config.WebServer().MFA().RPOrigin()
	if rpid == "" {
		app.GetLogger().Errorf("RPID is not set, WebAuthn will likely not work as intended")
	}

	if rporigin == "" {
		app.GetLogger().Errorf("RPOrigin is not set, WebAuthn will likely not work as intended")
	}

	return sessions.WebAuthnConfiguration{
		RPID:     rpid,
		RPOrigin: rporigin,
	}
}

func (app *ChainlinkApplication) ID() uuid.UUID {
	return app.Config.AppID()
}

var ErrUnsupportedInLOOPPMode = fmt.Errorf("legacy command not available in LOOP Plugin mode: %w", stderrors.ErrUnsupported)

// FindLCA - finds last common ancestor
func (app *ChainlinkApplication) FindLCA(ctx context.Context, chainID *big.Int) (*logpoller.Block, error) {
	chain, err := app.GetRelayers().LegacyEVMChains().Get(chainID.String())
	if err != nil {
		return nil, err
	}
	if !app.Config.Feature().LogPoller() {
		return nil, errors.New("FindLCA is only available if LogPoller is enabled")
	}
	legacyChain, ok := chain.(legacyevm.Chain)
	if !ok {
		return nil, ErrUnsupportedInLOOPPMode
	}

	lca, err := legacyChain.LogPoller().FindLCA(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find lca: %w", err)
	}

	return lca, nil
}

// DeleteLogPollerDataAfter - delete LogPoller state starting from the specified block
func (app *ChainlinkApplication) DeleteLogPollerDataAfter(ctx context.Context, chainID *big.Int, start int64) error {
	chain, err := app.GetRelayers().LegacyEVMChains().Get(chainID.String())
	if err != nil {
		return err
	}
	if !app.Config.Feature().LogPoller() {
		return errors.New("DeleteLogPollerDataAfter is only available if LogPoller is enabled")
	}
	legacyChain, ok := chain.(legacyevm.Chain)
	if !ok {
		return ErrUnsupportedInLOOPPMode
	}

	err = legacyChain.LogPoller().DeleteLogsAndBlocksAfter(ctx, start)
	if err != nil {
		return fmt.Errorf("failed to recover LogPoller: %w", err)
	}

	return nil
}
