package chainlink

import (
	"bytes"
	"context"
	stderrors "errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/grafana/pyroscope-go"
	"github.com/jonboulle/clockwork"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/credentials"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/billing"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	nodeauthjwt "github.com/smartcontractkit/chainlink-common/pkg/nodeauth/jwt"
	commonsrv "github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/otelhealth"
	"github.com/smartcontractkit/chainlink-common/pkg/services/promhealth"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/storage"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/dontime"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	evmutils "github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/build"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/compute"
	gatewayconnector "github.com/smartcontractkit/chainlink/v2/core/capabilities/gateway_connector"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockheaderfeeder"
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
	externalp2p "github.com/smartcontractkit/chainlink/v2/core/services/p2p"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	p2pwrapper "github.com/smartcontractkit/chainlink/v2/core/services/p2p/wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/services/periodicbackup"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	registrysyncerV1 "github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
	registrysyncerV2 "github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf"
	"github.com/smartcontractkit/chainlink/v2/core/services/webhook"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	artifactsV1 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts"
	artifactsV2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	workflowstore "github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	syncerV1 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer"
	syncerV2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	wftypes "github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
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
	RunJobV2(ctx context.Context, jobID int32, meta map[string]interface{}) (int64, error)

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

	started     bool
	startStopMu sync.Mutex
}

type ApplicationOpts struct {
	// CREOpts is the options for the CRE services
	CREOpts

	Config                   GeneralConfig
	Logger                   logger.Logger
	Registerer               prometheus.Registerer
	DS                       sqlutil.DataSource
	KeyStore                 keystore.Master
	AuditLogger              audit.AuditLogger
	CloseLogger              func() error
	ExternalInitiatorManager webhook.ExternalInitiatorManager
	Version                  string
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
	LimitsFactory            limits.Factory
}

// NewApplication initializes a new store if one is not already
// present at the configured root directory (default: ~/.chainlink),
// the logger at the same directory and returns the Application to
// be used by the node.
// TODO: Inject more dependencies here to save booting up useless stuff in tests
func NewApplication(ctx context.Context, opts ApplicationOpts) (Application, error) {
	var srvcs []services.ServiceCtx

	heartbeat := NewHeartbeat(NewHeartbeatConfig(opts))
	srvcs = append(srvcs, &heartbeat)

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

	csaKeystore := &keystore.CSASigner{CSA: keyStore.CSA()}
	beholderAuthHeaders, csaPubKeyHex, err := keystore.BuildBeholderAuth(ctx, keyStore.CSA())
	if err != nil {
		return nil, fmt.Errorf("failed to build Beholder auth: %w", err)
	}
	loopRegistry := plugins.NewLoopRegistry(globalLogger, cfg.AppID().String(), cfg.Feature().LogPoller(), cfg.Database(), cfg.Mercury(), cfg.Tracing(), cfg.Telemetry(), beholderAuthHeaders, csaPubKeyHex)

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

	relayChainInterops, err := NewCoreRelayerChainInteroperators(initOps...)
	if err != nil {
		return nil, err
	}

	var billingClient metering.BillingClient

	signer, CSAPubKey, err := keystore.BuildNodeAuth(ctx, csaKeystore)
	if err != nil {
		return nil, fmt.Errorf("failed to build node auth: %w", err)
	}
	jwtGenerator := nodeauthjwt.NewNodeJWTGenerator(signer, CSAPubKey)

	if opts.BillingClient != nil {
		billingClient = opts.BillingClient
	} else if cfg.Billing().URL() != "" {
		workflowOpts := []billing.WorkflowClientOpt{}

		if opts.Config.Billing().TLSEnabled() {
			workflowOpts = append(workflowOpts, billing.WithWorkflowTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
		}

		workflowOpts = append(workflowOpts, billing.WithJWTGenerator(jwtGenerator))

		billingClient, err = billing.NewWorkflowClient(globalLogger, opts.Config.Billing().URL(), workflowOpts...)
		if err != nil {
			globalLogger.Infof("NewApplication: failed to create billing client; %s", err)
		}
	}

	var storageClient storage.WorkflowClient
	if opts.StorageClient != nil {
		storageClient = opts.StorageClient
	} else if cfg.Capabilities().WorkflowRegistry().WorkflowStorage().URL() != "" {
		workflowOpts := []storage.WorkflowClientOpt{}

		if opts.Config.Capabilities().WorkflowRegistry().WorkflowStorage().TLSEnabled() {
			workflowOpts = append(workflowOpts, storage.WithWorkflowTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
		}

		workflowOpts = append(workflowOpts, storage.WithJWTGenerator(jwtGenerator))

		storageClient, err = storage.NewWorkflowClient(globalLogger, opts.Config.Capabilities().WorkflowRegistry().WorkflowStorage().URL(), workflowOpts...)
		if err != nil {
			globalLogger.Infof("NewApplication: failed to create storage client; %s", err)
		}
	}

	creServices, err := newCREServices(ctx, globalLogger, opts.DS, keyStore, cfg.Capabilities(), cfg.Workflows(), relayChainInterops, opts.CREOpts, billingClient, storageClient, opts.DonTimeStore, opts.LimitsFactory)
	if err != nil {
		return nil, fmt.Errorf("failed to initilize CRE: %w", err)
	}
	srvcs = append(srvcs, creServices.srvs...)
	// LOOPs can be created as options, in the  case of LOOP relayers, or
	// as OCR2 job implementations, in the case of Median today.
	// We will have a non-nil registry here in LOOP relayers are being used, otherwise
	// we need to initialize in case we serve OCR2 LOOPs
	if loopRegistry == nil {
		loopRegistry = plugins.NewLoopRegistry(globalLogger, opts.Config.AppID().String(), opts.Config.Feature().LogPoller(), opts.Config.Database(), opts.Config.Mercury(), opts.Config.Tracing(), opts.Config.Telemetry(), beholderAuthHeaders, csaPubKeyHex)
	}

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
				creServices.workflowRegistrySyncer,
				globalLogger),
			job.Stream: streams.NewDelegate(
				globalLogger,
				streamRegistry,
				pipelineRunner,
				cfg.JobPipeline(),
			),
		}
		webhookJobRunner = delegates[job.Webhook].(*webhook.Delegate).WebhookJobRunner()
	)

	delegates[job.Workflow] = workflows.NewDelegate(
		globalLogger,
		opts.CapabilitiesRegistry,
		opts.DonTimeStore,
		workflowORM,
		creServices.workflowRateLimiter,
		creServices.workflowLimits,
		workflows.WithBillingClient(billingClient),
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

	var peerWrapper *ocrcommon.SingletonPeerWrapper
	if !cfg.OCR().Enabled() && !cfg.OCR2().Enabled() {
		globalLogger.Debug("P2P stack not needed")
	} else if cfg.P2P().Enabled() {
		if err := ocrcommon.ValidatePeerWrapperConfig(cfg.P2P()); err != nil {
			return nil, err
		}
		peerWrapper = ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), opts.DS, globalLogger)
		srvcs = append(srvcs, peerWrapper)
	} else {
		return nil, errors.New("P2P stack required for OCR or OCR2")
	}

	// If peer wrapper is initialized, Oracle Factory dependency will be available to standard capabilities
	delegates[job.StandardCapabilities] = standardcapabilities.NewDelegate(
		globalLogger,
		opts.DS, jobORM,
		opts.CapabilitiesRegistry,
		loopRegistrarConfig,
		telemetryManager,
		pipelineRunner,
		relayChainInterops,
		creServices.gatewayConnectorWrapper,
		keyStore,
		creServices.externalPeerWrapper,
		peerWrapper,
		opts.NewOracleFactoryFn,
		opts.FetcherFactoryFn,
	)

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

		ocr2DelegateConfig := ocr2.NewDelegateConfig(cfg.OCR2(), cfg.Mercury(), cfg.Threshold(), cfg.Insecure(), cfg.JobPipeline(), loopRegistrarConfig)

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
				Relayers:                       relayChainInterops,
				MailMon:                        mailMon,
				CapabilitiesRegistry:           opts.CapabilitiesRegistry,
				DonTimeStore:                   opts.DonTimeStore,
				RetirementReportCache:          opts.RetirementReportCache,
				GatewayConnectorServiceWrapper: creServices.gatewayConnectorWrapper,
				WorkflowRegistrySyncer:         creServices.workflowRegistrySyncer,
			},
			ocr2DelegateConfig,
		)
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

	// We start the log poller after the job spawner
	// so jobs have a chance to apply their initial log filters.
	if cfg.Feature().LogPoller() {
		for _, c := range legacyEVMChains.Slice() {
			legacyChain, ok := c.(legacyevm.Chain)
			if !ok {
				continue
			}
			srvcs = append(srvcs, legacyChain.LogPoller())
		}
	}

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
			opts.Version,
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
			return nil, err
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

		ds: opts.DS,

		// NOTE: Can keep things clean by putting more things in srvcs instead of manually start/closing
		srvcs: srvcs,
	}, nil
}

// creKeystore is the minimal interface needed from keystore for CRE
type creKeystore interface {
	Eth() keystore.Eth
	P2P() keystore.P2P
	Workflow() keystore.Workflow
}

// CREOpts are the options for the CRE services that are exposed by the application
type CREOpts struct {
	CapabilitiesRegistry    *capabilities.Registry
	CapabilitiesDispatcher  remotetypes.Dispatcher
	CapabilitiesPeerWrapper p2ptypes.PeerWrapper

	FetcherFunc      wftypes.FetcherFunc
	FetcherFactoryFn compute.FetcherFactory

	BillingClient metering.BillingClient

	StorageClient storage.WorkflowClient

	UseLocalTimeProvider bool // Set this to true if the DON Time Plugin is not running
}

// creServiceConfig contains the configuration required to create the CRE services
type creServiceConfig struct {
	CREOpts

	capabilityCfg        config.Capabilities
	workflowsCfg         config.Workflows
	keystore             creKeystore
	logger               logger.Logger
	relayerChainInterops *CoreRelayerChainInteroperators
	DS                   sqlutil.DataSource
}

type CREServices struct {
	// workflowRateLimiter is the rate limiter for workflows
	// it is exposed because there are contingent services in the application
	workflowRateLimiter limits.RateLimiter

	// workflowLimits is the syncer limiter for workflows
	// it will specify the amount of global and per owner workflows that can be registered
	workflowLimits limits.ResourceLimiter[int]

	// gatewayConnectorWrapper is the wrapper for the gateway connector
	// it is exposed because there are contingent services in the application
	gatewayConnectorWrapper *gatewayconnector.ServiceWrapper

	// externalPeerWrapper is the wrapper for external peering
	// it is exposed because there are contingent services in the application
	externalPeerWrapper p2ptypes.PeerWrapper

	// srvs are all the services that are created, including those that are explicitly exposed
	srvs []services.ServiceCtx

	workflowRegistrySyncer syncerV2.WorkflowRegistrySyncer
}

func newCREServices(
	ctx context.Context,
	globalLogger logger.Logger,
	ds sqlutil.DataSource,
	keyStore creKeystore,
	capCfg config.Capabilities,
	wCfg config.Workflows,
	relayerChainInterops *CoreRelayerChainInteroperators,
	opts CREOpts,
	billingClient metering.BillingClient,
	storageClient storage.WorkflowClient,
	dontimeStore *dontime.Store,
	lf limits.Factory,
) (*CREServices, error) {
	var srvcs []services.ServiceCtx
	engineLimiters, err := v2.NewLimiters(lf, nil)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate engine limiters: %w", err)
	}
	srvcs = append(srvcs, closerService{name: "WorkflowEngineLimiters", Closer: engineLimiters})

	workflowRateLimiter, err := ratelimiter.NewRateLimiter(ratelimiter.Config{
		GlobalRPS:      capCfg.RateLimit().GlobalRPS(),
		GlobalBurst:    capCfg.RateLimit().GlobalBurst(),
		PerSenderRPS:   capCfg.RateLimit().PerSenderRPS(),
		PerSenderBurst: capCfg.RateLimit().PerSenderBurst(),
	}, lf)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate workflow rate limiter: %w", err)
	}
	srvcs = append(srvcs, closerService{name: "WorkflowRateLimiter", Closer: workflowRateLimiter})

	if len(wCfg.Limits().PerOwnerOverrides()) > 0 {
		globalLogger.Debugw("loaded per owner overrides", "overrides", wCfg.Limits().PerOwnerOverrides())
	}

	workflowLimits, err := syncerlimiter.NewWorkflowLimits(globalLogger, syncerlimiter.Config{
		Global:            wCfg.Limits().Global(),
		PerOwner:          wCfg.Limits().PerOwner(),
		PerOwnerOverrides: wCfg.Limits().PerOwnerOverrides(),
	}, lf)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate workflow syncer limiter: %w", err)
	}
	srvcs = append(srvcs, closerService{name: "WorkflowExecutionLimiter", Closer: workflowLimits})

	var gatewayConnectorWrapper *gatewayconnector.ServiceWrapper
	if capCfg.GatewayConnector().DonID() != "" {
		globalLogger.Debugw("Creating GatewayConnector wrapper", "donID", capCfg.GatewayConnector().DonID())
		chainID, ok := new(big.Int).SetString(capCfg.GatewayConnector().ChainIDForNodeKey(), 0)
		if !ok {
			return nil, fmt.Errorf("failed to parse gateway connector chain ID as integer: %s", capCfg.GatewayConnector().ChainIDForNodeKey())
		}
		gatewayConnectorWrapper = gatewayconnector.NewGatewayConnectorServiceWrapper(
			capCfg.GatewayConnector(),
			keys.NewStore(keystore.NewEthSigner(keyStore.Eth(), chainID)),
			clockwork.NewRealClock(),
			globalLogger)
		srvcs = append(srvcs, gatewayConnectorWrapper)
	}

	var externalPeerWrapper p2ptypes.PeerWrapper
	var workflowRegistrySyncer syncerV2.WorkflowRegistrySyncer
	if capCfg.Peering().Enabled() {
		var dispatcher remotetypes.Dispatcher
		if opts.CapabilitiesDispatcher == nil {
			externalPeerWrapper = p2pwrapper.NewExternalPeerWrapper(keyStore.P2P(), capCfg.Peering(), ds, globalLogger)
			signer := externalp2p.NewSigner(keyStore.P2P(), capCfg.Peering().PeerID())
			remoteDispatcher, err := remote.NewDispatcher(capCfg.Dispatcher(), externalPeerWrapper, signer, opts.CapabilitiesRegistry, globalLogger)
			if err != nil {
				return nil, fmt.Errorf("could not create dispatcher: %w", err)
			}
			dispatcher = remoteDispatcher
		} else {
			dispatcher = opts.CapabilitiesDispatcher
			externalPeerWrapper = opts.CapabilitiesPeerWrapper
		}

		srvcs = append(srvcs, externalPeerWrapper, dispatcher)

		if capCfg.ExternalRegistry().Address() != "" {
			rid := capCfg.ExternalRegistry().RelayID()
			registryAddress := capCfg.ExternalRegistry().Address()
			relayer, err := relayerChainInterops.Get(rid)
			if err != nil {
				return nil, fmt.Errorf("could not fetch relayer %s configured for capabilities registry: %w", rid, err)
			}

			externalRegistryVersion, err := semver.NewVersion(capCfg.ExternalRegistry().ContractVersion())
			if err != nil {
				return nil, err
			}

			workflowDonNotifier := capabilities.NewDonNotifier()
			wfLauncher := capabilities.NewLauncher(
				globalLogger,
				externalPeerWrapper,
				dispatcher,
				opts.CapabilitiesRegistry,
				workflowDonNotifier,
			)

			switch externalRegistryVersion.Major() {
			case 1:
				registrySyncer, err := registrysyncerV1.New(
					globalLogger,
					func() (p2ptypes.PeerID, error) {
						p := externalPeerWrapper.GetPeer()
						if p == nil {
							return p2ptypes.PeerID{}, errors.New("could not get peer")
						}

						return p.ID(), nil
					},
					relayer,
					registryAddress,
					registrysyncerV1.NewORM(ds, globalLogger),
				)
				if err != nil {
					return nil, fmt.Errorf("could not configure syncer: %w", err)
				}

				registrySyncer.AddListener(wfLauncher)
				srvcs = append(srvcs, wfLauncher, registrySyncer)
			case 2:
				registrySyncer, err := registrysyncerV2.New(
					globalLogger,
					func() (p2ptypes.PeerID, error) {
						p := externalPeerWrapper.GetPeer()
						if p == nil {
							return p2ptypes.PeerID{}, errors.New("could not get peer")
						}

						return p.ID(), nil
					},
					relayer,
					registryAddress,
					registrysyncerV1.NewORM(ds, globalLogger),
				)
				if err != nil {
					return nil, fmt.Errorf("could not configure syncer: %w", err)
				}

				registrySyncer.AddListener(wfLauncher)
				srvcs = append(srvcs, wfLauncher, registrySyncer)
			default:
				return nil, fmt.Errorf("could not configure capability registry syncer with version: %d", externalRegistryVersion.Major())
			}

			if capCfg.WorkflowRegistry().Address() != "" {
				globalLogger.Debugw("Creating WorkflowRegistrySyncer")
				lggr := globalLogger.Named("WorkflowRegistrySyncer")

				key, err := keystore.GetDefault(ctx, keyStore.Workflow())
				if err != nil {
					return nil, fmt.Errorf("failed to get all workflow keys: %w", err)
				}

				wfRegRid := capCfg.WorkflowRegistry().RelayID()
				wfRegRelayer, err := relayerChainInterops.Get(wfRegRid)
				if err != nil {
					return nil, fmt.Errorf("could not fetch relayer %s configured for workflow registry: %w", rid, err)
				}

				crFactory := func(ctx context.Context, bytes []byte) (commontypes.ContractReader, error) {
					return wfRegRelayer.NewContractReader(ctx, bytes)
				}

				wrVersion, vErr := semver.NewVersion(capCfg.WorkflowRegistry().ContractVersion())
				if vErr != nil {
					return nil, vErr
				}

				switch wrVersion.Major() {
				case 1:
					var fetcherFunc wftypes.FetcherFunc
					if opts.FetcherFunc == nil {
						if gatewayConnectorWrapper == nil {
							return nil, errors.New("unable to create workflow registry syncer without gateway connector")
						}
						fetcher := syncerV1.NewFetcherService(lggr, gatewayConnectorWrapper)
						fetcherFunc = fetcher.Fetch
						srvcs = append(srvcs, fetcher)
					} else {
						fetcherFunc = opts.FetcherFunc
					}

					artifactsStore := artifactsV1.NewStore(lggr, artifactsV1.NewWorkflowRegistryDS(ds, globalLogger),
						fetcherFunc,
						clockwork.NewRealClock(), key, custmsg.NewLabeler(), artifactsV1.WithMaxArtifactSize(
							artifactsV1.ArtifactConfig{
								MaxBinarySize:  uint64(capCfg.WorkflowRegistry().MaxBinarySize()),
								MaxSecretsSize: uint64(capCfg.WorkflowRegistry().MaxEncryptedSecretsSize()),
								MaxConfigSize:  uint64(capCfg.WorkflowRegistry().MaxConfigSize()),
							},
						))

					engineRegistry := syncerV1.NewEngineRegistry()

					eventHandler, err := syncerV1.NewEventHandler(
						lggr,
						workflowstore.NewInMemoryStore(lggr, clockwork.NewRealClock()),
						opts.CapabilitiesRegistry,
						dontimeStore,
						opts.UseLocalTimeProvider,
						engineRegistry,
						custmsg.NewLabeler(),
						engineLimiters,
						workflowRateLimiter,
						workflowLimits,
						artifactsStore,
						key,
						syncerV1.WithBillingClient(billingClient),
						syncerV1.WithWorkflowRegistry(capCfg.WorkflowRegistry().Address(), capCfg.WorkflowRegistry().ChainID()),
					)
					if err != nil {
						return nil, fmt.Errorf("unable to create workflow registry event handler: %w", err)
					}

					wfSyncer, err := syncerV1.NewWorkflowRegistry(
						lggr,
						crFactory,
						capCfg.WorkflowRegistry().Address(),
						syncerV1.Config{
							QueryCount:   100,
							SyncStrategy: syncerV1.SyncStrategy(capCfg.WorkflowRegistry().SyncStrategy()),
						},
						eventHandler,
						workflowDonNotifier,
						engineRegistry,
					)
					if err != nil {
						return nil, fmt.Errorf("unable to create workflow registry syncer: %w", err)
					}

					srvcs = append(srvcs, wfSyncer)
					globalLogger.Debugw("Created WorkflowRegistrySyncer V1")

				case 2:
					var fetcherFunc wftypes.FetcherFunc
					var retrieverFunc wftypes.LocationRetrieverFunc
					if opts.FetcherFunc == nil {
						if gatewayConnectorWrapper == nil {
							return nil, errors.New("unable to create workflow registry syncer without gateway connector")
						}
						if storageClient == nil {
							return nil, errors.New("unable to create workflow registry syncer without storage client")
						}
						fetcher := syncerV2.NewFetcherService(lggr, gatewayConnectorWrapper, storageClient)
						fetcherFunc = fetcher.Fetch
						retrieverFunc = fetcher.RetrieveURL
						srvcs = append(srvcs, fetcher)
					} else {
						fetcherFunc = opts.FetcherFunc
						retrieverFunc = nil
					}

					artifactsStore, err := artifactsV2.NewStore(lggr, artifactsV2.NewWorkflowRegistryDS(ds, globalLogger),
						fetcherFunc,
						retrieverFunc,
						clockwork.NewRealClock(), key, custmsg.NewLabeler(), artifactsV2.WithMaxArtifactSize(
							artifactsV2.ArtifactConfig{
								MaxBinarySize:  uint64(capCfg.WorkflowRegistry().MaxBinarySize()),
								MaxSecretsSize: uint64(capCfg.WorkflowRegistry().MaxEncryptedSecretsSize()),
								MaxConfigSize:  uint64(capCfg.WorkflowRegistry().MaxConfigSize()),
							},
						),
						artifactsV2.WithConfig(artifactsV2.StoreConfig{
							ArtifactStorageHost: capCfg.WorkflowRegistry().WorkflowStorage().ArtifactStorageHost(),
						}))
					if err != nil {
						return nil, fmt.Errorf("unable to create artifact store: %w", err)
					}

					engineRegistry := syncerV2.NewEngineRegistry()

					eventHandler, err := syncerV2.NewEventHandler(
						lggr,
						workflowstore.NewInMemoryStore(lggr, clockwork.NewRealClock()),
						dontimeStore,
						opts.UseLocalTimeProvider,
						opts.CapabilitiesRegistry,
						engineRegistry,
						custmsg.NewLabeler(),
						engineLimiters,
						workflowRateLimiter,
						workflowLimits,
						artifactsStore,
						key,
						syncerV2.WithBillingClient(billingClient),
						syncerV2.WithWorkflowRegistry(capCfg.WorkflowRegistry().Address(), capCfg.WorkflowRegistry().ChainID()),
					)
					if err != nil {
						return nil, fmt.Errorf("unable to create workflow registry event handler: %w", err)
					}

					wfSyncer, err := syncerV2.NewWorkflowRegistry(
						lggr,
						crFactory,
						capCfg.WorkflowRegistry().Address(),
						syncerV2.Config{
							QueryCount:   100,
							SyncStrategy: syncerV2.SyncStrategy(capCfg.WorkflowRegistry().SyncStrategy()),
						},
						eventHandler,
						workflowDonNotifier,
						engineRegistry,
					)
					if err != nil {
						return nil, fmt.Errorf("unable to create workflow registry syncer: %w", err)
					}

					srvcs = append(srvcs, wfSyncer)
					globalLogger.Debugw("Created WorkflowRegistrySyncer V2")

				default:
					return nil, fmt.Errorf("unsupported WorkflowRegistry contract version %s", wrVersion)
				}
			}
		}
	} else {
		globalLogger.Debug("External registry not configured, skipping registry syncer and starting with an empty registry")
		opts.CapabilitiesRegistry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
	}
	return &CREServices{
		workflowRateLimiter:     workflowRateLimiter,
		workflowLimits:          workflowLimits,
		gatewayConnectorWrapper: gatewayConnectorWrapper,
		externalPeerWrapper:     externalPeerWrapper,
		srvs:                    srvcs,
		workflowRegistrySyncer:  workflowRegistrySyncer,
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
	meta map[string]interface{},
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
		var vars map[string]interface{}
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
					evmutils.NewHash().Bytes()},              // requestID
					[]byte{}),
				Topics:      []common.Hash{{}, jb.ExternalIDEncodeBytesToTopic()}, // jobID BYTES
				TxHash:      evmutils.NewHash(),
				BlockNumber: 10,
				BlockHash:   evmutils.NewHash(),
			}
			vars = map[string]interface{}{
				"jobSpec": map[string]interface{}{
					"databaseID":    jb.ID,
					"externalJobID": jb.ExternalJobID,
					"name":          jb.Name.ValueOrZero(),
					"publicKey":     jb.VRFSpec.PublicKey[:],
					"evmChainID":    jb.VRFSpec.EVMChainID.String(),
				},
				"jobRun": map[string]interface{}{
					"meta":           meta,
					"logBlockHash":   testLog.BlockHash[:],
					"logBlockNumber": testLog.BlockNumber,
					"logTxHash":      testLog.TxHash,
					"logTopics":      testLog.Topics,
					"logData":        testLog.Data,
				},
			}
		} else {
			vars = map[string]interface{}{
				"jobRun": map[string]interface{}{
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

var _ services.ServiceCtx = closerService{}

// closerService extends an io.Closer to implement [services.ServiceCtx]
type closerService struct {
	name string
	io.Closer
}

func (c closerService) Start(ctx context.Context) error { return nil }

func (c closerService) Ready() error { return nil }

func (c closerService) HealthReport() map[string]error { return map[string]error{c.Name(): nil} }

func (c closerService) Name() string { return c.name }
