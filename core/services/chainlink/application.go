package chainlink

import (
	"context"
	stderr "errors"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/services/cron"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/core/services/gasupdater"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/periodicbackup"
	"github.com/smartcontractkit/chainlink/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	"gorm.io/gorm"

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
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/models/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"
)

//go:generate mockery --name ExternalInitiatorManager --output ../../internal/mocks/ --case=underscore
type (
	// ExternalInitiatorManager manages HTTP requests to remote external initiators
	ExternalInitiatorManager interface {
		Notify(models.JobSpec, *strpkg.Store) error
		DeleteJob(db *gorm.DB, jobID models.JobID) error
	}
)

//go:generate mockery --name Application --output ../../internal/mocks/ --case=underscore

// Application implements the common functions used in the core node.
type Application interface {
	Start() error
	Stop() error
	GetLogger() *logger.Logger
	GetHealthChecker() health.Checker
	GetStore() *strpkg.Store
	GetOCRKeyStore() *offchainreporting.KeyStore // TODO: this should be replaced with a generic GetKeystore()
	GetStatsPusher() synchronization.StatsPusher
	GetHeadBroadcaster() httypes.HeadBroadcasterRegistry
	WakeSessionReaper()
	AddServiceAgreement(*models.ServiceAgreement) error
	NewBox() packr.Box

	// V1 Jobs (JSON specified)
	services.RunManager // For managing job runs.
	AddJob(job models.JobSpec) error
	ArchiveJob(models.JobID) error
	GetExternalInitiatorManager() ExternalInitiatorManager

	// V2 Jobs (TOML specified)
	GetJobORM() job.ORM
	AddJobV2(ctx context.Context, job job.Job, name null.String) (int32, error)
	DeleteJobV2(ctx context.Context, jobID int32) error
	RunWebhookJobV2(ctx context.Context, jobUUID models.JobID, pipelineInput interface{}, meta pipeline.JSONSerializable) (int64, error)
	// Testing only
	RunJobV2(ctx context.Context, jobID int32, meta map[string]interface{}) (int64, error)
	SetServiceLogger(ctx context.Context, service string, level zapcore.Level) error
}

// ChainlinkApplication contains fields for the JobSubscriber, Scheduler,
// and Store. The JobSubscriber and Scheduler are also available
// in the services package, but the Store has its own package.
type ChainlinkApplication struct {
	Exiter          func(int)
	HeadTracker     *services.HeadTracker
	HeadBroadcaster *headtracker.HeadBroadcaster
	BPTXM           *bulletprooftxmanager.BulletproofTxManager
	StatsPusher     synchronization.StatsPusher
	services.RunManager
	RunQueue         services.RunQueue
	JobSubscriber    services.JobSubscriber
	LogBroadcaster   log.Broadcaster
	EventBroadcaster postgres.EventBroadcaster
	JobORM           job.ORM
	jobSpawner       job.Spawner
	pipelineRunner   pipeline.Runner
	FluxMonitor      fluxmonitor.Service
	webhookJobRunner webhook.JobRunner
	Scheduler        *services.Scheduler
	Store            *strpkg.Store
	// TODO:
	// moved OCR keystore from store to application in order to resolve:
	// https://app.clubhouse.io/chainlinklabs/story/10097/remove-ocr-as-dependency-of-store-package

	// waiting on this before combining and moving other keystores
	// https://github.com/smartcontractkit/chainlink/pull/4447

	// finally, keystore unification will be completed by:
	// https://app.clubhouse.io/chainlinklabs/story/7735/combine-keystores
	OCRKeyStore              *offchainreporting.KeyStore
	ExternalInitiatorManager ExternalInitiatorManager
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
func NewApplication(config *orm.Config, ethClient eth.Client, advisoryLocker postgres.AdvisoryLocker, keyStoreGenerator strpkg.KeyStoreGenerator, externalInitiatorManager ExternalInitiatorManager, onConnectCallbacks ...func(Application)) (Application, error) {
	var subservices []service.Service

	shutdownSignal := gracefulpanic.NewSignal()
	store, err := strpkg.NewStore(config, ethClient, advisoryLocker, shutdownSignal, keyStoreGenerator)
	if err != nil {
		return nil, err
	}

	setupConfig(config, store)

	headBroadcaster := headtracker.NewHeadBroadcaster()

	healthChecker := health.NewChecker()

	explorerClient := synchronization.ExplorerClient(&synchronization.NoopExplorerClient{})
	statsPusher := synchronization.StatsPusher(&synchronization.NoopStatsPusher{})
	monitoringEndpoint := ocrtypes.MonitoringEndpoint(&telemetry.NoopAgent{})

	if config.ExplorerURL() != nil {
		explorerClient = synchronization.NewExplorerClient(config.ExplorerURL(), config.ExplorerAccessKey(), config.ExplorerSecret(), config.StatsPusherLogging())
		statsPusher = synchronization.NewStatsPusher(store.DB, explorerClient)
		monitoringEndpoint = telemetry.NewAgent(explorerClient)
	}
	subservices = append(subservices, explorerClient, statsPusher)

	if store.Config.GasUpdaterEnabled() {
		logger.Debugw("GasUpdater: dynamic gas updates are enabled", "ethGasPriceDefault", store.Config.EthGasPriceDefault())
		gasUpdater := gasupdater.NewGasUpdater(store.EthClient, store.Config)
		headBroadcaster.Subscribe(gasUpdater)
		subservices = append(subservices, gasUpdater)
	} else {
		logger.Debugw("GasUpdater: dynamic gas updating is disabled", "ethGasPriceDefault", store.Config.EthGasPriceDefault())
	}

	if store.Config.DatabaseBackupMode() != orm.DatabaseBackupModeNone && store.Config.DatabaseBackupFrequency() > 0 {
		logger.Infow("DatabaseBackup: periodic database backups are enabled", "frequency", store.Config.DatabaseBackupFrequency())

		databaseBackup := periodicbackup.NewDatabaseBackup(store.Config, logger.Default)
		subservices = append(subservices, databaseBackup)
	} else {
		logger.Info("DatabaseBackup: periodic database backups are disabled")
	}

	// Init service loggers
	globalLogger := config.CreateProductionLogger()
	globalLogger.SetDB(store.DB)
	serviceLogLevels, err := globalLogger.GetServiceLogLevels()
	if err != nil {
		logger.Fatalf("error getting log levels: %v", err)
	}
	headTrackerLogger, err := globalLogger.InitServiceLevelLogger(logger.HeadTracker, serviceLogLevels[logger.HeadTracker])
	if err != nil {
		logger.Fatal("error starting logger for head tracker", err)
	}

	headTracker := services.NewHeadTracker(headTrackerLogger, store, headBroadcaster)

	var runExecutor services.RunExecutor
	var runQueue services.RunQueue
	var runManager services.RunManager
	var jobSubscriber services.JobSubscriber
	if config.EnableLegacyJobPipeline() {
		runExecutor = services.NewRunExecutor(store, statsPusher)
		runQueue = services.NewRunQueue(runExecutor)
		runManager = services.NewRunManager(runQueue, config, store.ORM, statsPusher, store.Clock)
		jobSubscriber = services.NewJobSubscriber(store, runManager)
	} else {
		runExecutor = &services.NullRunExecutor{}
		runQueue = &services.NullRunQueue{}
		runManager = &services.NullRunManager{}
		jobSubscriber = &services.NullJobSubscriber{}
	}

	// Highest seen head height is used as part of the start of LogBroadcaster backfill range
	highestSeenHead, err := headTracker.HighestSeenHeadFromDB()
	if err != nil {
		return nil, err
	}

	logBroadcaster := log.NewBroadcaster(log.NewORM(store.DB), ethClient, store.Config, highestSeenHead)
	eventBroadcaster := postgres.NewEventBroadcaster(config.DatabaseURL(), config.DatabaseListenerMinReconnectInterval(), config.DatabaseListenerMaxReconnectDuration())
	fluxMonitor := fluxmonitor.New(store, runManager, logBroadcaster)

	bptxm := bulletprooftxmanager.NewBulletproofTxManager(store.DB, ethClient, store.Config, store.KeyStore, advisoryLocker, eventBroadcaster)

	promReporter := services.NewPromReporter(store.MustSQLDB())

	subservices = append(subservices,
		logBroadcaster,
		eventBroadcaster,
		fluxMonitor,
		jobSubscriber,
	)

	var balanceMonitor services.BalanceMonitor
	if config.BalanceMonitorEnabled() {
		balanceMonitor = services.NewBalanceMonitor(store)
	} else {
		balanceMonitor = &services.NullBalanceMonitor{}
	}

	subservices = append(subservices, balanceMonitor)

	subservices = append(subservices, promReporter)

	var (
		pipelineORM    = pipeline.NewORM(store.ORM.DB, store.Config, eventBroadcaster)
		pipelineRunner = pipeline.NewRunner(pipelineORM, store.Config)
		jobORM         = job.NewORM(store.ORM.DB, store.Config, pipelineORM, eventBroadcaster, advisoryLocker)
	)

	var (
		delegates = map[job.Type]job.Delegate{
			job.DirectRequest: directrequest.NewDelegate(
				logBroadcaster,
				pipelineRunner,
				pipelineORM,
				ethClient,
				store.DB,
				config,
			),
			job.Keeper: keeper.NewDelegate(store.DB, jobORM, pipelineRunner, store.EthClient, headBroadcaster, logBroadcaster, config),
			job.VRF:    vrf.NewDelegate(vrf.NewORM(store.DB), pipelineRunner, pipelineORM),
		}
	)

	if config.Dev() || config.FeatureFluxMonitorV2() {
		delegates[job.FluxMonitor] = fluxmonitorv2.NewDelegate(
			store,
			jobORM,
			pipelineORM,
			pipelineRunner,
			store.DB,
			ethClient,
			logBroadcaster,
			fluxmonitorv2.Config{
				DefaultHTTPTimeout:       store.Config.DefaultHTTPTimeout().Duration(),
				FlagsContractAddress:     store.Config.FlagsContractAddress(),
				MinContractPayment:       store.Config.MinimumContractPayment(),
				EthGasLimit:              store.Config.EthGasLimitDefault(),
				EthMaxQueuedTransactions: store.Config.EthMaxQueuedTransactions(),
			},
		)
	}

	scryptParams := utils.GetScryptParams(config)
	ocrKeyStore := offchainreporting.NewKeyStore(store.DB, scryptParams)

	if (config.Dev() && config.P2PListenPort() > 0) || config.FeatureOffchainReporting() {
		logger.Debug("Off-chain reporting enabled")
		concretePW := offchainreporting.NewSingletonPeerWrapper(ocrKeyStore, config, store.DB)
		subservices = append(subservices, concretePW)
		delegates[job.OffchainReporting] = offchainreporting.NewDelegate(
			store.DB,
			jobORM,
			config,
			ocrKeyStore,
			pipelineRunner,
			ethClient,
			logBroadcaster,
			concretePW,
			monitoringEndpoint,
		)
	} else {
		logger.Debug("Off-chain reporting disabled")
	}

	var webhookJobRunner webhook.JobRunner
	if config.Dev() || config.FeatureWebhookV2() {
		delegate := webhook.NewDelegate(pipelineRunner)
		delegates[job.Webhook] = delegate
		webhookJobRunner = delegate.WebhookJobRunner()
	}

	if config.Dev() || config.FeatureCronV2() {
		delegates[job.Cron] = cron.NewDelegate(pipelineRunner)
	}

	jobSpawner := job.NewSpawner(jobORM, store.Config, delegates)
	subservices = append(subservices, jobSpawner, pipelineRunner, bptxm, headBroadcaster)

	app := &ChainlinkApplication{
		HeadBroadcaster:          headBroadcaster,
		BPTXM:                    bptxm,
		JobSubscriber:            jobSubscriber,
		LogBroadcaster:           logBroadcaster,
		EventBroadcaster:         eventBroadcaster,
		JobORM:                   jobORM,
		jobSpawner:               jobSpawner,
		pipelineRunner:           pipelineRunner,
		FluxMonitor:              fluxMonitor,
		StatsPusher:              statsPusher,
		RunManager:               runManager,
		RunQueue:                 runQueue,
		webhookJobRunner:         webhookJobRunner,
		Scheduler:                services.NewScheduler(store, runManager),
		Store:                    store,
		OCRKeyStore:              ocrKeyStore,
		SessionReaper:            services.NewStoreReaper(store),
		Exiter:                   os.Exit,
		ExternalInitiatorManager: externalInitiatorManager,
		shutdownSignal:           shutdownSignal,
		balanceMonitor:           balanceMonitor,
		explorerClient:           explorerClient,
		HealthChecker:            healthChecker,
		logger:                   globalLogger,
		// NOTE: Can keep things clean by putting more things in subservices
		// instead of manually start/closing
		subservices: subservices,
	}

	headBroadcaster.Subscribe(logBroadcaster)
	headBroadcaster.Subscribe(bptxm)
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
	app.HeadTracker = headTracker

	// Log Broadcaster waits for other services' registrations
	// until app.LogBroadcaster.DependentReady() call (see below)
	logBroadcaster.AddDependents(1)

	for _, service := range app.subservices {
		if err = app.HealthChecker.Register(reflect.TypeOf(service).String(), service); err != nil {
			return nil, err
		}
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

func setupConfig(config *orm.Config, store *strpkg.Store) {
	config.SetRuntimeStore(store.ORM)

	if !config.P2PPeerIDIsSet() {
		var keys []p2pkey.EncryptedP2PKey
		err := store.DB.Order("created_at asc, id asc").Find(&keys).Error
		if err != nil {
			logger.Warnw("Failed to load keys", "err", err)
		} else {
			if len(keys) > 0 {
				peerID := keys[0].PeerID
				logger.Debugw("P2P_PEER_ID was not set, using the first available key", "peerID", peerID.String())
				config.Set("P2P_PEER_ID", peerID)
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
	if err := app.Store.EthClient.Dial(context.TODO()); err != nil {
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

		logger.Info("Exited all services")

		app.started = false
	})
	return merr
}

// GetStore returns the pointer to the store for the ChainlinkApplication.
func (app *ChainlinkApplication) GetStore() *strpkg.Store {
	return app.Store
}

func (app *ChainlinkApplication) GetOCRKeyStore() *offchainreporting.KeyStore {
	return app.OCRKeyStore
}

func (app *ChainlinkApplication) GetLogger() *logger.Logger {
	return app.logger
}

func (app *ChainlinkApplication) GetHealthChecker() health.Checker {
	return app.HealthChecker
}

func (app *ChainlinkApplication) GetJobORM() job.ORM {
	return app.JobORM
}

func (app *ChainlinkApplication) GetExternalInitiatorManager() ExternalInitiatorManager {
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

func (app *ChainlinkApplication) AddJobV2(ctx context.Context, job job.Job, name null.String) (int32, error) {
	return app.jobSpawner.CreateJob(ctx, job, name)
}

func (app *ChainlinkApplication) DeleteJobV2(ctx context.Context, jobID int32) error {
	return app.jobSpawner.DeleteJob(ctx, jobID)
}

func (app *ChainlinkApplication) RunWebhookJobV2(ctx context.Context, jobUUID models.JobID, pipelineInput interface{}, meta pipeline.JSONSerializable) (int64, error) {
	return app.webhookJobRunner.RunJob(ctx, jobUUID, pipelineInput, meta)
}

// Only used for testing, not supported by the UI.
func (app *ChainlinkApplication) RunJobV2(
	ctx context.Context,
	jobID int32,
	meta map[string]interface{},
) (int64, error) {
	if !app.Store.Config.Dev() {
		return 0, errors.New("manual job runs only supported in dev mode - export CHAINLINK_DEV=true to use.")
	}
	jb, err := app.JobORM.FindJob(jobID)
	if err != nil {
		return 0, errors.Wrapf(err, "job ID %v", jobID)
	}
	var runID int64

	// Keeper jobs are special in that they do not have a task graph.
	if jb.Type == job.Keeper {
		t := time.Now()
		runID, err = app.pipelineRunner.InsertFinishedRun(app.Store.DB.WithContext(ctx), pipeline.Run{
			PipelineSpecID: jb.PipelineSpecID,
			Errors:         pipeline.RunErrors{null.String{}},
			Outputs:        pipeline.JSONSerializable{Val: "queued eth transaction"},
			CreatedAt:      t,
			FinishedAt:     &t,
		}, nil, false)
	} else {
		meta := pipeline.JSONSerializable{
			Val:  meta,
			Null: false,
		}
		runID, _, err = app.pipelineRunner.ExecuteAndInsertFinishedRun(ctx, *jb.PipelineSpec, nil, meta, *logger.Default, false)
	}
	return runID, err
}

// ArchiveJob silences the job from the system, preventing future job runs.
// It is idempotent and can be run as many times as you like.
func (app *ChainlinkApplication) ArchiveJob(ID models.JobID) error {
	err := app.JobSubscriber.RemoveJob(ID)
	if err != nil {
		logger.Warnw("Error removing job from JobSubscriber", "error", err)
	}
	app.FluxMonitor.RemoveJob(ID)

	if err = app.ExternalInitiatorManager.DeleteJob(app.Store.DB, ID); err != nil {
		err = errors.Wrapf(err, "failed to delete job with id %s from external initiator", ID)
	}
	return multierr.Combine(err, app.Store.ArchiveJob(ID))
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

// NewBox returns the packr.Box instance that holds the static assets to
// be delivered by the router.
func (app *ChainlinkApplication) NewBox() packr.Box {
	return packr.NewBox("../../../operator_ui/dist")
}
