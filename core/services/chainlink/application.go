package chainlink

import (
	"bytes"
	"context"
	stderr "errors"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/services/cron"
	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/periodicbackup"
	"github.com/smartcontractkit/chainlink/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/services/webhook"

	"github.com/gobuffalo/packr"
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitor"
	"github.com/smartcontractkit/chainlink/core/services/health"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"
)

//go:generate mockery --name Application --output ../../internal/mocks/ --case=underscore

// Application implements the common functions used in the core node.
type Application interface {
	Start() error
	Stop() error
	GetLogger() *logger.Logger
	GetHealthChecker() health.Checker
	GetStore() *strpkg.Store
	GetEthClient() eth.Client
	GetConfig() *config.Config
	GetKeyStore() *keystore.Master
	GetStatsPusher() synchronization.StatsPusher
	GetHeadBroadcaster() httypes.HeadBroadcasterRegistry
	WakeSessionReaper()
	AddServiceAgreement(*models.ServiceAgreement) error
	NewBox() packr.Box

	// V1 Jobs (JSON specified)
	services.RunManager // For managing job runs.
	AddJob(job models.JobSpec) error
	ArchiveJob(models.JobID) error
	GetExternalInitiatorManager() webhook.ExternalInitiatorManager

	// V2 Jobs (TOML specified)
	JobSpawner() job.Spawner
	JobORM() job.ORM
	PipelineORM() pipeline.ORM
	AddJobV2(ctx context.Context, job job.Job, name null.String) (job.Job, error)
	DeleteJobV2(ctx context.Context, jobID int32) error
	RunWebhookJobV2(ctx context.Context, jobUUID uuid.UUID, requestBody string, meta pipeline.JSONSerializable) (int64, error)
	ResumeJobV2(ctx context.Context, run *pipeline.Run) (bool, error)
	// Testing only
	RunJobV2(ctx context.Context, jobID int32, meta map[string]interface{}) (int64, error)
	SetServiceLogger(ctx context.Context, service string, level zapcore.Level) error

	// Feeds
	GetFeedsService() feeds.Service

	// ReplayFromBlock of blocks
	ReplayFromBlock(number uint64) error
}

// ChainlinkApplication contains fields for the JobSubscriber, Scheduler,
// and Store. The JobSubscriber and Scheduler are also available
// in the services package, but the Store has its own package.
type ChainlinkApplication struct {
	Exiter          func(int)
	HeadTracker     httypes.Tracker
	HeadBroadcaster httypes.HeadBroadcaster
	TxManager       bulletprooftxmanager.TxManager
	StatsPusher     synchronization.StatsPusher
	services.RunManager
	RunQueue                 services.RunQueue
	JobSubscriber            services.JobSubscriber
	LogBroadcaster           log.Broadcaster
	EventBroadcaster         postgres.EventBroadcaster
	jobORM                   job.ORM
	jobSpawner               job.Spawner
	pipelineORM              pipeline.ORM
	pipelineRunner           pipeline.Runner
	FluxMonitor              fluxmonitor.Service
	FeedsService             feeds.Service
	webhookJobRunner         webhook.JobRunner
	Scheduler                *services.Scheduler
	ethClient                eth.Client
	Store                    *strpkg.Store
	Config                   *config.Config
	KeyStore                 *keystore.Master
	ExternalInitiatorManager webhook.ExternalInitiatorManager
	SessionReaper            utils.SleeperTask
	shutdownOnce             sync.Once
	shutdownSignal           gracefulpanic.Signal
	balanceMonitor           services.BalanceMonitor
	explorerClient           synchronization.ExplorerClient
	subservices              []service.Service
	HealthChecker            health.Checker
	logger                   *logger.Logger

	started     bool
	startStopMu sync.Mutex
}

// NewApplication initializes a new store if one is not already
// present at the configured root directory (default: ~/.chainlink),
// the logger at the same directory and returns the Application to
// be used by the node.
func NewApplication(cfg *config.Config, ethClient eth.Client, advisoryLocker postgres.AdvisoryLocker, onConnectCallbacks ...func(Application)) (Application, error) {
	var subservices []service.Service

	shutdownSignal := gracefulpanic.NewSignal()
	store, err := strpkg.NewStore(cfg, ethClient, advisoryLocker, shutdownSignal)
	if err != nil {
		return nil, err
	}
	gormTxm := postgres.NewGormTransactionManager(store.DB)

	setupConfig(cfg, store.DB)

	healthChecker := health.NewChecker()

	scryptParams := utils.GetScryptParams(cfg)
	keyStore := keystore.New(store.DB, scryptParams)

	explorerClient := synchronization.ExplorerClient(&synchronization.NoopExplorerClient{})
	statsPusher := synchronization.StatsPusher(&synchronization.NoopStatsPusher{})
	monitoringEndpoint := ocrtypes.MonitoringEndpoint(&telemetry.NoopAgent{})

	if cfg.ExplorerURL() != nil {
		explorerClient = synchronization.NewExplorerClient(cfg.ExplorerURL(), cfg.ExplorerAccessKey(), cfg.ExplorerSecret(), cfg.StatsPusherLogging())
		statsPusher = synchronization.NewStatsPusher(store.DB, explorerClient)
		monitoringEndpoint = telemetry.NewAgent(explorerClient)
	}
	subservices = append(subservices, explorerClient, statsPusher)

	if cfg.DatabaseBackupMode() != config.DatabaseBackupModeNone && cfg.DatabaseBackupFrequency() > 0 {
		logger.Infow("DatabaseBackup: periodic database backups are enabled", "frequency", cfg.DatabaseBackupFrequency())

		databaseBackup := periodicbackup.NewDatabaseBackup(cfg, logger.Default)
		subservices = append(subservices, databaseBackup)
	} else {
		logger.Info("DatabaseBackup: periodic database backups are disabled. To enable automatic backups, set DATABASE_BACKUP_MODE=lite or DATABASE_BACKUP_MODE=full")
	}

	// Init service loggers
	globalLogger := cfg.CreateProductionLogger()
	globalLogger.SetDB(store.DB)
	serviceLogLevels, err := globalLogger.GetServiceLogLevels()
	if err != nil {
		logger.Fatalf("error getting log levels: %v", err)
	}
	headTrackerLogger, err := globalLogger.InitServiceLevelLogger(logger.HeadTracker, serviceLogLevels[logger.HeadTracker])
	if err != nil {
		logger.Fatal("error starting logger for head tracker", err)
	}

	var headBroadcaster httypes.HeadBroadcaster
	var headTracker httypes.Tracker
	if cfg.EthereumDisabled() {
		headBroadcaster = &headtracker.NullBroadcaster{}
		headTracker = &headtracker.NullTracker{}
	} else {
		headBroadcaster = headtracker.NewHeadBroadcaster()
		orm := headtracker.NewORM(store.DB)
		headTracker = headtracker.NewHeadTracker(headTrackerLogger, ethClient, cfg, orm, headBroadcaster)
	}

	var runExecutor services.RunExecutor
	var runQueue services.RunQueue
	var runManager services.RunManager
	var jobSubscriber services.JobSubscriber
	if cfg.EnableLegacyJobPipeline() {
		runExecutor = services.NewRunExecutor(store, ethClient, keyStore, statsPusher)
		runQueue = services.NewRunQueue(runExecutor)
		runManager = services.NewRunManager(runQueue, cfg, store.ORM, statsPusher, store.Clock)
		jobSubscriber = services.NewJobSubscriber(store, runManager, ethClient)
	} else {
		runExecutor = &services.NullRunExecutor{}
		runQueue = &services.NullRunQueue{}
		runManager = &services.NullRunManager{}
		jobSubscriber = &services.NullJobSubscriber{}
	}

	eventBroadcaster := postgres.NewEventBroadcaster(cfg.DatabaseURL(), cfg.DatabaseListenerMinReconnectInterval(), cfg.DatabaseListenerMaxReconnectDuration())
	subservices = append(subservices, eventBroadcaster)

	var txManager bulletprooftxmanager.TxManager
	var logBroadcaster log.Broadcaster
	if cfg.EthereumDisabled() {
		txManager = &bulletprooftxmanager.NullTxManager{ErrMsg: "TxManager is not running because Ethereum is disabled"}
		logBroadcaster = &log.NullBroadcaster{ErrMsg: "LogBroadcaster is not running because Ethereum is disabled"}
	} else {
		// Highest seen head height is used as part of the start of LogBroadcaster backfill range
		highestSeenHead, err2 := headTracker.HighestSeenHeadFromDB()
		if err2 != nil {
			return nil, err2
		}

		logBroadcaster = log.NewBroadcaster(log.NewORM(store.DB), ethClient, cfg, highestSeenHead)
		txManager = bulletprooftxmanager.NewBulletproofTxManager(store.DB, ethClient, cfg, keyStore.Eth(), advisoryLocker, eventBroadcaster)
		subservices = append(subservices, logBroadcaster, txManager)
	}

	fluxMonitor := fluxmonitor.New(store, keyStore.Eth(), runManager, logBroadcaster, ethClient)

	subservices = append(subservices,
		fluxMonitor,
		jobSubscriber,
	)

	var balanceMonitor services.BalanceMonitor
	if cfg.BalanceMonitorEnabled() {
		balanceMonitor = services.NewBalanceMonitor(store.DB, ethClient, keyStore.Eth())
	} else {
		balanceMonitor = &services.NullBalanceMonitor{}
	}
	subservices = append(subservices, balanceMonitor)

	promReporter := services.NewPromReporter(store.MustSQLDB())
	subservices = append(subservices, promReporter)

	var (
		pipelineORM    = pipeline.NewORM(store.DB)
		pipelineRunner = pipeline.NewRunner(pipelineORM, cfg, ethClient, keyStore.Eth(), keyStore.VRF(), txManager)
		jobORM         = job.NewORM(store.ORM.DB, cfg, pipelineORM, eventBroadcaster, advisoryLocker)
	)

	var (
		delegates = map[job.Type]job.Delegate{
			job.DirectRequest: directrequest.NewDelegate(
				logBroadcaster,
				pipelineRunner,
				pipelineORM,
				ethClient,
				store.DB,
				cfg,
			),
			job.Keeper: keeper.NewDelegate(store.DB, txManager, jobORM, pipelineRunner, ethClient, headBroadcaster, logBroadcaster, cfg),
			job.VRF: vrf.NewDelegate(
				store.DB,
				txManager,
				keyStore,
				pipelineRunner,
				pipelineORM,
				logBroadcaster,
				headBroadcaster,
				ethClient,
				cfg),
		}
	)

	// Flux monitor requires ethereum just to boot, silence errors with a null delegate
	if cfg.EthereumDisabled() {
		delegates[job.FluxMonitor] = &job.NullDelegate{Type: job.FluxMonitor}
	} else if cfg.Dev() || cfg.FeatureFluxMonitorV2() {
		delegates[job.FluxMonitor] = fluxmonitorv2.NewDelegate(
			txManager,
			keyStore.Eth(),
			jobORM,
			pipelineORM,
			pipelineRunner,
			store.DB,
			ethClient,
			logBroadcaster,
			fluxmonitorv2.Config{
				DefaultHTTPTimeout:             cfg.DefaultHTTPTimeout().Duration(),
				FlagsContractAddress:           cfg.FlagsContractAddress(),
				MinContractPayment:             cfg.MinimumContractPayment(),
				EthGasLimit:                    cfg.EthGasLimitDefault(),
				EthMaxQueuedTransactions:       cfg.EthMaxQueuedTransactions(),
				FMDefaultTransactionQueueDepth: cfg.FMDefaultTransactionQueueDepth(),
			},
		)
	}

	if (cfg.Dev() && cfg.P2PListenPort() > 0) || cfg.FeatureOffchainReporting() {
		logger.Debug("Off-chain reporting enabled")
		concretePW := offchainreporting.NewSingletonPeerWrapper(keyStore.OCR(), cfg, store.DB)
		subservices = append(subservices, concretePW)
		delegates[job.OffchainReporting] = offchainreporting.NewDelegate(
			store.DB,
			txManager,
			jobORM,
			cfg,
			keyStore.OCR(),
			pipelineRunner,
			ethClient,
			logBroadcaster,
			concretePW,
			monitoringEndpoint,
			cfg.Chain(),
			headBroadcaster,
		)
	} else {
		logger.Debug("Off-chain reporting disabled")
	}

	externalInitiatorManager := webhook.NewExternalInitiatorManager(store.DB, utils.UnrestrictedClient)

	var webhookJobRunner webhook.JobRunner
	if cfg.Dev() || cfg.FeatureWebhookV2() {
		delegate := webhook.NewDelegate(pipelineRunner, externalInitiatorManager)
		delegates[job.Webhook] = delegate
		webhookJobRunner = delegate.WebhookJobRunner()
	}

	if cfg.Dev() || cfg.FeatureCronV2() {
		delegates[job.Cron] = cron.NewDelegate(pipelineRunner)
	}

	jobSpawner := job.NewSpawner(jobORM, cfg, delegates, gormTxm)
	subservices = append(subservices, jobSpawner, pipelineRunner, headBroadcaster)

	feedsORM := feeds.NewORM(store.DB)
	feedsService := feeds.NewService(feedsORM, gormTxm, jobSpawner, keyStore.CSA(), keyStore.Eth(), cfg)

	app := &ChainlinkApplication{
		ethClient:                ethClient,
		HeadBroadcaster:          headBroadcaster,
		TxManager:                txManager,
		JobSubscriber:            jobSubscriber,
		LogBroadcaster:           logBroadcaster,
		EventBroadcaster:         eventBroadcaster,
		jobORM:                   jobORM,
		jobSpawner:               jobSpawner,
		pipelineRunner:           pipelineRunner,
		pipelineORM:              pipelineORM,
		FluxMonitor:              fluxMonitor,
		FeedsService:             feedsService,
		StatsPusher:              statsPusher,
		Config:                   cfg,
		RunManager:               runManager,
		RunQueue:                 runQueue,
		webhookJobRunner:         webhookJobRunner,
		Scheduler:                services.NewScheduler(store, runManager),
		Store:                    store,
		KeyStore:                 keyStore,
		SessionReaper:            services.NewSessionReaper(store.DB, cfg),
		Exiter:                   os.Exit,
		ExternalInitiatorManager: externalInitiatorManager,
		shutdownSignal:           shutdownSignal,
		balanceMonitor:           balanceMonitor,
		explorerClient:           explorerClient,
		HealthChecker:            healthChecker,
		HeadTracker:              headTracker,
		logger:                   globalLogger,
		// NOTE: Can keep things clean by putting more things in subservices
		// instead of manually start/closing
		subservices: subservices,
	}

	headBroadcaster.Subscribe(logBroadcaster)
	headBroadcaster.Subscribe(txManager)
	headBroadcaster.Subscribe(promReporter)
	headBroadcaster.Subscribe(jobSubscriber)
	headBroadcaster.Subscribe(balanceMonitor)

	headBroadcaster.Subscribe(&httypes.HeadTrackableCallback{OnConnect: func() error {
		return runManager.ResumeAllPendingConnection()
	}})

	for _, onConnectCallback := range onConnectCallbacks {
		headBroadcaster.Subscribe(&httypes.HeadTrackableCallback{OnConnect: func() error {
			onConnectCallback(app)
			return nil
		}})
	}

	// Log Broadcaster waits for other services' registrations
	// until app.LogBroadcaster.DependentReady() call (see below)
	logBroadcaster.AddDependents(1)

	for _, service := range app.subservices {
		if err = app.HealthChecker.Register(reflect.TypeOf(service).String(), service); err != nil {
			return nil, err
		}
	}

	if err = app.HealthChecker.Register(reflect.TypeOf(headTracker).String(), headTracker); err != nil {
		return nil, err
	}

	return app, nil
}

// SetServiceLogger sets the Logger for a given service and stores the setting in the db
func (app *ChainlinkApplication) SetServiceLogger(ctx context.Context, serviceName string, level zapcore.Level) error {
	newL, err := app.logger.InitServiceLevelLogger(serviceName, level.String())
	if err != nil {
		return err
	}

	// TODO: Implement other service loggers
	switch serviceName {
	case logger.HeadTracker:
		app.HeadTracker.SetLogger(newL)
	case logger.FluxMonitor:
		app.FluxMonitor.SetLogger(newL)
	default:
		return fmt.Errorf("no service found with name: %s", serviceName)
	}

	return app.logger.Orm.SetServiceLogLevel(ctx, serviceName, level)
}

func setupConfig(cfg *config.Config, db *gorm.DB) {
	orm := config.NewORM(db)
	cfg.SetRuntimeStore(orm)

	if !cfg.P2PPeerIDIsSet() {
		var keys []p2pkey.EncryptedP2PKey
		err := db.Order("created_at asc, id asc").Find(&keys).Error
		if err != nil {
			logger.Warnw("Failed to load keys", "err", err)
		} else {
			if len(keys) > 0 {
				peerID := keys[0].PeerID
				logger.Debugw("P2P_PEER_ID was not set, using the first available key", "peerID", peerID.String())
				cfg.Set("P2P_PEER_ID", peerID)
				if len(keys) > 1 {
					logger.Warnf("Found more than one P2P key in the database, but no P2P_PEER_ID was specified. Defaulting to first key: %s. Please consider setting P2P_PEER_ID explicitly.", peerID.String())
				}
			}
		}
	}
}

// Start all necessary services. If successful, nil will be returned.  Also
// listens for interrupt signals from the operating system so that the
// application can be properly closed before the application exits.
func (app *ChainlinkApplication) Start() error {
	app.startStopMu.Lock()
	defer app.startStopMu.Unlock()
	if app.started {
		panic("application is already started")
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-sigs:
		case <-app.shutdownSignal.Wait():
		}
		logger.ErrorIf(app.Stop())
		app.Exiter(0)
	}()

	// EthClient must be dialed first because it is required in subtasks
	if err := app.ethClient.Dial(context.Background()); err != nil {
		return err
	}

	if err := app.Store.Start(); err != nil {
		return err
	}

	if err := app.RunQueue.Start(); err != nil {
		return err
	}

	if err := app.RunManager.ResumeAllInProgress(); err != nil {
		return err
	}

	if err := app.FeedsService.Start(); err != nil {
		logger.Infof("[Feeds Service] %v", err)
	}

	for _, subservice := range app.subservices {
		logger.Debugw("Starting service...", "serviceType", reflect.TypeOf(subservice))
		if err := subservice.Start(); err != nil {
			return err
		}
	}

	// Log Broadcaster fully starts after all initial Register calls are done from other starting services
	// to make sure the initial backfill covers those subscribers.
	app.LogBroadcaster.DependentReady()

	// HeadTracker deliberately started afterwards since several tasks are
	// registered as callbacks and it's sensible to have started them before
	// calling the first OnNewHead
	// For example:
	// RunManager.ResumeAllInProgress since it Connects JobSubscriber
	// which leads to writes of JobRuns RunStatus to the db.
	// https://www.pivotaltracker.com/story/show/162230780
	if err := app.HeadTracker.Start(); err != nil {
		return err
	}

	if err := app.Scheduler.Start(); err != nil {
		return err
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

// Stop allows the application to exit by halting schedules, closing
// logs, and closing the DB connection.
func (app *ChainlinkApplication) Stop() error {
	app.startStopMu.Lock()
	defer app.startStopMu.Unlock()
	return app.stop()
}

func (app *ChainlinkApplication) stop() error {
	if !app.started {
		panic("application is already stopped")
	}
	var merr error
	app.shutdownOnce.Do(func() {
		defer func() {
			if err := logger.Sync(); err != nil {
				if stderr.Unwrap(err).Error() != os.ErrInvalid.Error() &&
					stderr.Unwrap(err).Error() != "inappropriate ioctl for device" &&
					stderr.Unwrap(err).Error() != "bad file descriptor" {
					merr = multierr.Append(merr, err)
				}
			}
		}()
		logger.Info("Gracefully exiting...")

		// Stop services in the reverse order from which they were started

		logger.Debug("Stopping Scheduler...")
		merr = multierr.Append(merr, app.Scheduler.Stop())

		logger.Debug("Stopping HeadTracker...")
		merr = multierr.Append(merr, app.HeadTracker.Stop())

		for i := len(app.subservices) - 1; i >= 0; i-- {
			service := app.subservices[i]
			logger.Debugw("Closing service...", "serviceType", reflect.TypeOf(service))
			merr = multierr.Append(merr, service.Close())
		}

		logger.Debug("Closing RunQueue...")
		merr = multierr.Append(merr, app.RunQueue.Close())
		logger.Debug("Stopping SessionReaper...")
		merr = multierr.Append(merr, app.SessionReaper.Stop())
		logger.Debug("Closing Store...")
		merr = multierr.Append(merr, app.Store.Close())
		logger.Debug("Closing HealthChecker...")
		merr = multierr.Append(merr, app.HealthChecker.Close())
		logger.Debug("Closing Feeds Service...")
		merr = multierr.Append(merr, app.FeedsService.Close())

		logger.Info("Exited all services")

		app.started = false
	})
	return merr
}

// GetStore returns the pointer to the store for the ChainlinkApplication.
func (app *ChainlinkApplication) GetStore() *strpkg.Store {
	return app.Store
}

func (app *ChainlinkApplication) GetEthClient() eth.Client {
	return app.ethClient
}

func (app *ChainlinkApplication) GetConfig() *config.Config {
	return app.Config
}

func (app *ChainlinkApplication) GetKeyStore() *keystore.Master {
	return app.KeyStore
}

func (app *ChainlinkApplication) GetLogger() *logger.Logger {
	return app.logger
}

func (app *ChainlinkApplication) GetHealthChecker() health.Checker {
	return app.HealthChecker
}

func (app *ChainlinkApplication) JobSpawner() job.Spawner {
	return app.jobSpawner
}

func (app *ChainlinkApplication) JobORM() job.ORM {
	return app.jobORM
}

func (app *ChainlinkApplication) PipelineORM() pipeline.ORM {
	return app.pipelineORM
}

func (app *ChainlinkApplication) GetExternalInitiatorManager() webhook.ExternalInitiatorManager {
	return app.ExternalInitiatorManager
}

func (app *ChainlinkApplication) GetHeadBroadcaster() httypes.HeadBroadcasterRegistry {
	return app.HeadBroadcaster
}

func (app *ChainlinkApplication) GetStatsPusher() synchronization.StatsPusher {
	return app.StatsPusher
}

// WakeSessionReaper wakes up the reaper to do its reaping.
func (app *ChainlinkApplication) WakeSessionReaper() {
	app.SessionReaper.WakeUp()
}

// AddJob adds a job to the store and the scheduler. If there was
// an error from adding the job to the store, the job will not be
// added to the scheduler.
func (app *ChainlinkApplication) AddJob(job models.JobSpec) error {
	err := app.Store.CreateJob(&job)
	if err != nil {
		return err
	}

	app.Scheduler.AddJob(job)
	logger.ErrorIf(app.FluxMonitor.AddJob(job))
	logger.ErrorIf(app.JobSubscriber.AddJob(job, nil))
	return nil
}

func (app *ChainlinkApplication) AddJobV2(ctx context.Context, j job.Job, name null.String) (job.Job, error) {
	return app.jobSpawner.CreateJob(ctx, j, name)
}

func (app *ChainlinkApplication) DeleteJobV2(ctx context.Context, jobID int32) error {
	return app.jobSpawner.DeleteJob(ctx, jobID)
}

func (app *ChainlinkApplication) RunWebhookJobV2(ctx context.Context, jobUUID uuid.UUID, requestBody string, meta pipeline.JSONSerializable) (int64, error) {
	return app.webhookJobRunner.RunJob(ctx, jobUUID, requestBody, meta)
}

// Only used for local testing, not supported by the UI.
func (app *ChainlinkApplication) RunJobV2(
	ctx context.Context,
	jobID int32,
	meta map[string]interface{},
) (int64, error) {
	if !app.Store.Config.Dev() {
		return 0, errors.New("manual job runs only supported in dev mode - export CHAINLINK_DEV=true to use")
	}
	jb, err := app.jobORM.FindJob(ctx, jobID)
	if err != nil {
		return 0, errors.Wrapf(err, "job ID %v", jobID)
	}
	var runID int64

	// Some jobs are special in that they do not have a task graph.
	isBootstrap := jb.Type == job.OffchainReporting && jb.OffchainreportingOracleSpec != nil && jb.OffchainreportingOracleSpec.IsBootstrapPeer
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
					utils.NewHash().Bytes(),                  // sender
					utils.NewHash().Bytes(),                  // fee
					utils.NewHash().Bytes()},                 // requestID
					[]byte{}),
				Topics:      []common.Hash{{}, jb.ExternalIDEncodeBytesToTopic()}, // jobID BYTES
				TxHash:      utils.NewHash(),
				BlockNumber: 10,
				BlockHash:   utils.NewHash(),
			}
			vars = map[string]interface{}{
				"jobSpec": map[string]interface{}{
					"databaseID":    jb.ID,
					"externalJobID": jb.ExternalJobID,
					"name":          jb.Name.ValueOrZero(),
					"publicKey":     jb.VRFSpec.PublicKey[:],
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
		runID, _, err = app.pipelineRunner.ExecuteAndInsertFinishedRun(ctx, *jb.PipelineSpec, pipeline.NewVarsFrom(vars), *logger.Default, saveTasks)
	} else {
		// This is a weird situation, even if a job doesn't have a pipeline it needs a pipeline_spec_id in order to insert the run
		// TODO: Once all jobs have a pipeline this can be removed
		// See: https://app.clubhouse.io/chainlinklabs/story/6065/hook-keeper-up-to-use-tasks-in-the-pipeline
		runID, err = app.pipelineRunner.TestInsertFinishedRun(app.Store.DB.WithContext(ctx), jb.ID, jb.Name.String, jb.Type.String(), jb.PipelineSpecID)
	}
	return runID, err
}

func (app *ChainlinkApplication) ResumeJobV2(
	ctx context.Context,
	run *pipeline.Run,
) (bool, error) {
	return app.pipelineRunner.Run(ctx, run, *logger.Default, false)
}

// ArchiveJob silences the job from the system, preventing future job runs.
// It is idempotent and can be run as many times as you like.
func (app *ChainlinkApplication) ArchiveJob(ID models.JobID) error {
	err := app.JobSubscriber.RemoveJob(ID)
	if err != nil {
		logger.Warnw("Error removing job from JobSubscriber", "error", err)
	}
	app.FluxMonitor.RemoveJob(ID)

	if err = app.ExternalInitiatorManager.DeleteJob(ID); err != nil {
		logger.Warnf("failed to delete job with id %s from external initiator: %v", ID, err)
	}
	return app.Store.ArchiveJob(ID)
}

// AddServiceAgreement adds a Service Agreement which includes a job that needs
// to be scheduled.
func (app *ChainlinkApplication) AddServiceAgreement(sa *models.ServiceAgreement) error {
	err := app.Store.CreateServiceAgreement(sa)
	if err != nil {
		return err
	}

	app.Scheduler.AddJob(sa.JobSpec)

	// XXX: Add mechanism to asynchronously communicate when a job spec has
	// an ethereum interaction error.
	// https://www.pivotaltracker.com/story/show/170349568
	logger.ErrorIf(app.FluxMonitor.AddJob(sa.JobSpec))
	logger.ErrorIf(app.JobSubscriber.AddJob(sa.JobSpec, nil))
	return nil
}

func (app *ChainlinkApplication) GetFeedsService() feeds.Service {
	return app.FeedsService
}

// NewBox returns the packr.Box instance that holds the static assets to
// be delivered by the router.
func (app *ChainlinkApplication) NewBox() packr.Box {
	return packr.NewBox("../../../operator_ui/dist")
}

func (app *ChainlinkApplication) ReplayFromBlock(number uint64) error {

	app.LogBroadcaster.ReplayFromBlock(int64(number))

	return nil
}
