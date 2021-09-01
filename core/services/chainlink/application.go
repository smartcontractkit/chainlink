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
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	loggerPkg "github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/cron"
	"github.com/smartcontractkit/chainlink/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/core/services/health"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/periodicbackup"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/core/services/versioning"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name Application --output ../../internal/mocks/ --case=underscore

// Application implements the common functions used in the core node.
type Application interface {
	Start() error
	Stop() error
	GetLogger() *loggerPkg.Logger
	GetHealthChecker() health.Checker
	GetStore() *strpkg.Store
	GetConfig() config.GeneralConfig
	GetKeyStore() keystore.Master
	GetEventBroadcaster() postgres.EventBroadcaster
	WakeSessionReaper()
	NewBox() packr.Box

	GetExternalInitiatorManager() webhook.ExternalInitiatorManager
	GetChainSet() evm.ChainSet

	// V2 Jobs (TOML specified)
	JobSpawner() job.Spawner
	JobORM() job.ORM
	EVMORM() evmtypes.ORM
	PipelineORM() pipeline.ORM
	AddJobV2(ctx context.Context, job job.Job, name null.String) (job.Job, error)
	DeleteJob(ctx context.Context, jobID int32) error
	RunWebhookJobV2(ctx context.Context, jobUUID uuid.UUID, requestBody string, meta pipeline.JSONSerializable) (int64, error)
	ResumeJobV2(ctx context.Context, taskID uuid.UUID, result interface{}) error
	// Testing only
	RunJobV2(ctx context.Context, jobID int32, meta map[string]interface{}) (int64, error)
	SetServiceLogger(ctx context.Context, service string, level zapcore.Level) error

	// Feeds
	GetFeedsService() feeds.Service

	// ReplayFromBlock of blocks
	ReplayFromBlock(chainID *big.Int, number uint64) error
}

// ChainlinkApplication contains fields for the JobSubscriber, Scheduler,
// and Store. The JobSubscriber and Scheduler are also available
// in the services package, but the Store has its own package.
type ChainlinkApplication struct {
	Exiter                   func(int)
	ChainSet                 evm.ChainSet
	EventBroadcaster         postgres.EventBroadcaster
	jobORM                   job.ORM
	jobSpawner               job.Spawner
	pipelineORM              pipeline.ORM
	pipelineRunner           pipeline.Runner
	FeedsService             feeds.Service
	webhookJobRunner         webhook.JobRunner
	evmORM                   evmtypes.ORM
	Store                    *strpkg.Store
	Config                   config.GeneralConfig
	KeyStore                 keystore.Master
	ExternalInitiatorManager webhook.ExternalInitiatorManager
	SessionReaper            utils.SleeperTask
	shutdownOnce             sync.Once
	shutdownSignal           gracefulpanic.Signal
	explorerClient           synchronization.ExplorerClient
	subservices              []service.Service
	HealthChecker            health.Checker
	logger                   *loggerPkg.Logger

	started     bool
	startStopMu sync.Mutex
}

type ApplicationOpts struct {
	Config                   config.GeneralConfig
	AdvisoryLocker           postgres.AdvisoryLocker
	EventBroadcaster         postgres.EventBroadcaster
	ShutdownSignal           gracefulpanic.Signal
	Store                    *strpkg.Store
	GormDB                   *gorm.DB
	KeyStore                 keystore.Master
	ChainSet                 evm.ChainSet
	Logger                   *loggerPkg.Logger
	ExternalInitiatorManager webhook.ExternalInitiatorManager
}

// NewApplication initializes a new store if one is not already
// present at the configured root directory (default: ~/.chainlink),
// the logger at the same directory and returns the Application to
// be used by the node.
// TODO: Inject more dependencies here to save booting up useless stuff in tests
func NewApplication(opts ApplicationOpts) (Application, error) {
	var subservices []service.Service
	db := opts.GormDB
	gormTxm := postgres.NewGormTransactionManager(db)
	cfg := opts.Config
	advisoryLocker := opts.AdvisoryLocker
	store := opts.Store
	shutdownSignal := opts.ShutdownSignal
	keyStore := opts.KeyStore
	chainSet := opts.ChainSet
	globalLogger := opts.Logger
	eventBroadcaster := opts.EventBroadcaster
	externalInitiatorManager := opts.ExternalInitiatorManager

	healthChecker := health.NewChecker()

	telemetryIngressClient := synchronization.TelemetryIngressClient(&synchronization.NoopTelemetryIngressClient{})
	explorerClient := synchronization.ExplorerClient(&synchronization.NoopExplorerClient{})
	monitoringEndpointGen := telemetry.MonitoringEndpointGenerator(&telemetry.NoopAgent{})

	if cfg.ExplorerURL() != nil {
		explorerClient = synchronization.NewExplorerClient(cfg.ExplorerURL(), cfg.ExplorerAccessKey(), cfg.ExplorerSecret(), cfg.StatsPusherLogging())
		monitoringEndpointGen = telemetry.NewExplorerAgent(explorerClient)
	}

	// Use Explorer over TelemetryIngress if both URLs are set
	if cfg.ExplorerURL() == nil && cfg.TelemetryIngressURL() != nil {
		telemetryIngressClient = synchronization.NewTelemetryIngressClient(cfg.TelemetryIngressURL(), cfg.TelemetryIngressServerPubKey(), keyStore.CSA(), cfg.TelemetryIngressLogging())
		monitoringEndpointGen = telemetry.NewIngressAgentWrapper(telemetryIngressClient)
	}
	subservices = append(subservices, explorerClient, telemetryIngressClient)

	if cfg.DatabaseBackupMode() != config.DatabaseBackupModeNone && cfg.DatabaseBackupFrequency() > 0 {
		logger.Infow("DatabaseBackup: periodic database backups are enabled", "frequency", cfg.DatabaseBackupFrequency())

		databaseBackup := periodicbackup.NewDatabaseBackup(cfg, globalLogger)
		subservices = append(subservices, databaseBackup)
	} else {
		logger.Info("DatabaseBackup: periodic database backups are disabled. To enable automatic backups, set DATABASE_BACKUP_MODE=lite or DATABASE_BACKUP_MODE=full")
	}

	subservices = append(subservices, eventBroadcaster, chainSet)
	promReporter := services.NewPromReporter(postgres.MustSQLDB(db))
	subservices = append(subservices, promReporter)

	var (
		pipelineORM    = pipeline.NewORM(db)
		pipelineRunner = pipeline.NewRunner(pipelineORM, cfg, chainSet, keyStore.Eth(), keyStore.VRF())
		jobORM         = job.NewORM(db, chainSet, pipelineORM, eventBroadcaster, advisoryLocker, keyStore)
	)

	for _, chain := range chainSet.Chains() {
		chain.HeadBroadcaster().Subscribe(promReporter)
		chain.TxManager().RegisterResumeCallback(pipelineRunner.ResumeRun)
	}

	var (
		delegates = map[job.Type]job.Delegate{
			job.DirectRequest: directrequest.NewDelegate(
				globalLogger,
				pipelineRunner,
				pipelineORM,
				db,
				chainSet),
			job.Keeper: keeper.NewDelegate(
				db,
				jobORM,
				pipelineRunner,
				globalLogger,
				cfg,
				chainSet),
			job.VRF: vrf.NewDelegate(
				db,
				keyStore,
				pipelineRunner,
				pipelineORM,
				chainSet),
		}
	)

	// Flux monitor requires ethereum just to boot, silence errors with a null delegate
	if cfg.EthereumDisabled() {
		delegates[job.FluxMonitor] = &job.NullDelegate{Type: job.FluxMonitor}
	} else if cfg.Dev() || cfg.FeatureFluxMonitorV2() {
		delegates[job.FluxMonitor] = fluxmonitorv2.NewDelegate(
			keyStore.Eth(),
			jobORM,
			pipelineORM,
			pipelineRunner,
			db,
			chainSet,
		)
	}

	if (cfg.Dev() && cfg.P2PListenPort() > 0) || cfg.FeatureOffchainReporting() {
		logger.Debug("Off-chain reporting enabled")
		concretePW := offchainreporting.NewSingletonPeerWrapper(keyStore, cfg, db)
		subservices = append(subservices, concretePW)
		delegates[job.OffchainReporting] = offchainreporting.NewDelegate(
			db,
			jobORM,
			keyStore.OCR(),
			pipelineRunner,
			concretePW,
			monitoringEndpointGen,
			chainSet,
		)
	} else {
		logger.Debug("Off-chain reporting disabled")
	}

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
	subservices = append(subservices, jobSpawner, pipelineRunner)

	feedsORM := feeds.NewORM(db)
	verORM := versioning.NewORM(postgres.WrapDbWithSqlx(
		postgres.MustSQLDB(db)),
	)

	// TODO: Make feeds manager compatible with multiple chains
	// See: https://app.clubhouse.io/chainlinklabs/story/14615/add-ability-to-set-chain-id-in-all-pipeline-tasks-that-interact-with-evm
	var feedsService feeds.Service
	chain, err := chainSet.Default()
	if err != nil {
		logger.Warnw("Unable to load feeds service; no default chain available", "err", err)
	} else {
		feedsService = feeds.NewService(feedsORM, verORM, gormTxm, jobSpawner, keyStore.CSA(), keyStore.Eth(), chain.Config(), chainSet)
	}

	app := &ChainlinkApplication{
		ChainSet:                 chainSet,
		Store:                    store,
		EventBroadcaster:         eventBroadcaster,
		jobORM:                   jobORM,
		jobSpawner:               jobSpawner,
		pipelineRunner:           pipelineRunner,
		pipelineORM:              pipelineORM,
		FeedsService:             feedsService,
		Config:                   cfg,
		webhookJobRunner:         webhookJobRunner,
		KeyStore:                 keyStore,
		SessionReaper:            services.NewSessionReaper(db, cfg),
		Exiter:                   os.Exit,
		ExternalInitiatorManager: externalInitiatorManager,
		shutdownSignal:           shutdownSignal,
		explorerClient:           explorerClient,
		HealthChecker:            healthChecker,
		logger:                   globalLogger,
		// NOTE: Can keep things clean by putting more things in subservices
		// instead of manually start/closing
		subservices: subservices,
	}

	for _, service := range app.subservices {
		if err := app.HealthChecker.Register(reflect.TypeOf(service).String(), service); err != nil {
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
	// FIXME: How to reconcile with multichain? This is broken and will currently overwrite the old label/value that sets evmChainID logger label
	// See: https://app.clubhouse.io/chainlinklabs/story/15452/initservicelevellogger-appears-to-discard-label-values-set-by-previous-with
	switch serviceName {
	case loggerPkg.HeadTracker:
		for _, c := range app.ChainSet.Chains() {
			c.HeadTracker().SetLogger(newL)
		}
	case loggerPkg.FluxMonitor:
		// TODO: Set FMv2?
	case loggerPkg.Keeper:
	default:
		return fmt.Errorf("no service found with name: %s", serviceName)
	}

	return app.logger.Orm.SetServiceLogLevel(ctx, serviceName, level)
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
		app.logger.ErrorIf(app.Stop())
		app.Exiter(0)
	}()

	if err := app.Store.Start(); err != nil {
		return err
	}

	if app.FeedsService != nil {
		if err := app.FeedsService.Start(); err != nil {
			app.logger.Infof("[Feeds Service] %v", err)
		}
	}

	for _, subservice := range app.subservices {
		app.logger.Debugw("Starting service...", "serviceType", reflect.TypeOf(subservice))
		if err := subservice.Start(); err != nil {
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
		done := make(chan error)
		go func() {
			var merr error
			defer func() {
				if lerr := app.logger.Sync(); lerr != nil {
					if stderr.Unwrap(lerr).Error() != os.ErrInvalid.Error() &&
						stderr.Unwrap(lerr).Error() != "inappropriate ioctl for device" &&
						stderr.Unwrap(lerr).Error() != "bad file descriptor" {
						merr = multierr.Append(merr, lerr)
					}
				}
			}()
			app.logger.Info("Gracefully exiting...")

			// Stop services in the reverse order from which they were started
			for i := len(app.subservices) - 1; i >= 0; i-- {
				service := app.subservices[i]
				app.logger.Debugw("Closing service...", "serviceType", reflect.TypeOf(service))
				merr = multierr.Append(merr, service.Close())
			}

			app.logger.Debug("Stopping SessionReaper...")
			merr = multierr.Append(merr, app.SessionReaper.Stop())
			app.logger.Debug("Closing Store...")
			merr = multierr.Append(merr, app.Store.Close())
			app.logger.Debug("Closing HealthChecker...")
			merr = multierr.Append(merr, app.HealthChecker.Close())
			if app.FeedsService != nil {
				app.logger.Debug("Closing Feeds Service...")
				merr = multierr.Append(merr, app.FeedsService.Close())
			}

			app.logger.Info("Exited all services")

			app.started = false
			done <- err
		}()
		select {
		case merr := <-done:
			err = merr
		case <-time.After(15 * time.Second):
			err = errors.New("application timed out shutting down")
		}
	})
	return err
}

func (app *ChainlinkApplication) GetConfig() config.GeneralConfig {
	return app.Config
}

func (app *ChainlinkApplication) GetKeyStore() keystore.Master {
	return app.KeyStore
}

func (app *ChainlinkApplication) GetLogger() *loggerPkg.Logger {
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

func (app *ChainlinkApplication) EVMORM() evmtypes.ORM {
	return app.ChainSet.ORM()
}

func (app *ChainlinkApplication) PipelineORM() pipeline.ORM {
	return app.pipelineORM
}

func (app *ChainlinkApplication) GetExternalInitiatorManager() webhook.ExternalInitiatorManager {
	return app.ExternalInitiatorManager
}

// WakeSessionReaper wakes up the reaper to do its reaping.
func (app *ChainlinkApplication) WakeSessionReaper() {
	app.SessionReaper.WakeUp()
}

func (app *ChainlinkApplication) AddJobV2(ctx context.Context, j job.Job, name null.String) (job.Job, error) {
	return app.jobSpawner.CreateJob(ctx, j, name)
}

func (app *ChainlinkApplication) DeleteJob(ctx context.Context, jobID int32) error {
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
		runID, _, err = app.pipelineRunner.ExecuteAndInsertFinishedRun(ctx, *jb.PipelineSpec, pipeline.NewVarsFrom(vars), *app.logger, saveTasks)
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
	taskID uuid.UUID,
	result interface{},
) error {
	return app.pipelineRunner.ResumeRun(taskID, result)
}

func (app *ChainlinkApplication) GetFeedsService() feeds.Service {
	return app.FeedsService
}

// NewBox returns the packr.Box instance that holds the static assets to
// be delivered by the router.
func (app *ChainlinkApplication) NewBox() packr.Box {
	return packr.NewBox("../../../operator_ui/dist")
}

func (app *ChainlinkApplication) ReplayFromBlock(chainID *big.Int, number uint64) error {
	chain, err := app.ChainSet.Get(chainID)
	if err != nil {
		return err
	}
	chain.LogBroadcaster().ReplayFromBlock(int64(number))
	return nil
}

func (app *ChainlinkApplication) GetChainSet() evm.ChainSet {
	return app.ChainSet
}

func (app *ChainlinkApplication) GetStore() *strpkg.Store {
	return app.Store
}

func (app *ChainlinkApplication) GetEventBroadcaster() postgres.EventBroadcaster {
	return app.EventBroadcaster
}

func (app *ChainlinkApplication) GetDB() *gorm.DB {
	return app.Store.DB
}
