package chainlink

import (
	"context"
	stderr "errors"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitor"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"

	"github.com/gobuffalo/packr"
	"go.uber.org/multierr"
)

// headTrackableCallback is a simple wrapper around an On Connect callback
type headTrackableCallback struct {
	onConnect func()
}

func (c *headTrackableCallback) Connect(*models.Head) error {
	c.onConnect()
	return nil
}

func (c *headTrackableCallback) Disconnect()                   {}
func (c *headTrackableCallback) OnNewLongestChain(models.Head) {}

//go:generate mockery --name Application --output ../internal/mocks/ --case=underscore

// Application implements the common functions used in the core node.
type Application interface {
	Start() error
	Stop() error
	GetStore() *strpkg.Store
	GetStatsPusher() synchronization.StatsPusher
	WakeSessionReaper()
	AddJob(job models.JobSpec) error
	ArchiveJob(*models.ID) error
	AddServiceAgreement(*models.ServiceAgreement) error
	NewBox() packr.Box
	services.RunManager
}

// ChainlinkApplication contains fields for the JobSubscriber, Scheduler,
// and Store. The JobSubscriber and Scheduler are also available
// in the services package, but the Store has its own package.
type ChainlinkApplication struct {
	Exiter      func(int)
	HeadTracker *services.HeadTracker
	StatsPusher synchronization.StatsPusher
	services.RunManager
	RunQueue                 services.RunQueue
	JobSubscriber            services.JobSubscriber
	GasUpdater               services.GasUpdater
	EthBroadcaster           bulletprooftxmanager.EthBroadcaster
	LogBroadcaster           eth.LogBroadcaster
	FluxMonitor              fluxmonitor.Service
	Scheduler                *services.Scheduler
	Store                    *strpkg.Store
	SessionReaper            services.SleeperTask
	pendingConnectionResumer *pendingConnectionResumer
	shutdownOnce             sync.Once
	shutdownSignal           gracefulpanic.Signal
}

// NewApplication initializes a new store if one is not already
// present at the configured root directory (default: ~/.chainlink),
// the logger at the same directory and returns the Application to
// be used by the node.
func NewApplication(config *orm.Config, onConnectCallbacks ...func(Application)) Application {
	shutdownSignal := gracefulpanic.NewSignal()
	store := strpkg.NewStore(config, shutdownSignal)
	config.SetRuntimeStore(store.ORM)

	statsPusher := synchronization.NewStatsPusher(
		store.ORM, config.ExplorerURL(), config.ExplorerAccessKey(), config.ExplorerSecret(),
	)
	runExecutor := services.NewRunExecutor(store, statsPusher)
	runQueue := services.NewRunQueue(runExecutor)
	runManager := services.NewRunManager(runQueue, config, store.ORM, statsPusher, store.TxManager, store.Clock)
	jobSubscriber := services.NewJobSubscriber(store, runManager)
	gasUpdater := services.NewGasUpdater(store)
	logBroadcaster := eth.NewLogBroadcaster(store.TxManager, store.ORM, store.Config.BlockBackfillDepth())
	fluxMonitor := fluxmonitor.New(store, runManager, logBroadcaster)
	ethBroadcaster := bulletprooftxmanager.NewEthBroadcaster(store, config)
	ethConfirmer := bulletprooftxmanager.NewEthConfirmer(store, config)
	balanceMonitor := services.NewBalanceMonitor(store)

	store.NotifyNewEthTx = ethBroadcaster

	pendingConnectionResumer := newPendingConnectionResumer(runManager)

	app := &ChainlinkApplication{
		JobSubscriber:            jobSubscriber,
		GasUpdater:               gasUpdater,
		EthBroadcaster:           ethBroadcaster,
		LogBroadcaster:           logBroadcaster,
		FluxMonitor:              fluxMonitor,
		StatsPusher:              statsPusher,
		RunManager:               runManager,
		RunQueue:                 runQueue,
		Scheduler:                services.NewScheduler(store, runManager),
		Store:                    store,
		SessionReaper:            services.NewStoreReaper(store),
		Exiter:                   os.Exit,
		pendingConnectionResumer: pendingConnectionResumer,
		shutdownSignal:           shutdownSignal,
	}

	headTrackables := []strpkg.HeadTrackable{gasUpdater}

	if store.Config.EnableBulletproofTxManager() {
		headTrackables = append(headTrackables, ethConfirmer)
	} else {
		headTrackables = append(headTrackables, store.TxManager)
	}

	headTrackables = append(
		headTrackables,
		jobSubscriber,
		pendingConnectionResumer,
		balanceMonitor,
	)

	for _, onConnectCallback := range onConnectCallbacks {
		headTrackable := &headTrackableCallback{func() {
			onConnectCallback(app)
		}}
		headTrackables = append(headTrackables, headTrackable)
	}
	app.HeadTracker = services.NewHeadTracker(store, headTrackables)

	return app
}

// Start all necessary services. If successful, nil will be returned.  Also
// listens for interrupt signals from the operating system so that the
// application can be properly closed before the application exits.
func (app *ChainlinkApplication) Start() error {
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

	err := app.Store.EthClient.Dial(context.TODO())
	if err != nil {
		return err
	}

	// XXX: Change to exit on first encountered error.
	return multierr.Combine(
		err,
		app.Store.Start(),
		app.StatsPusher.Start(),
		app.RunQueue.Start(),
		app.RunManager.ResumeAllInProgress(),
		app.LogBroadcaster.Start(),
		app.FluxMonitor.Start(),
		app.EthBroadcaster.Start(),

		// HeadTracker deliberately started after
		// RunManager.ResumeAllInProgress since it Connects JobSubscriber
		// which leads to writes of JobRuns RunStatus to the db.
		// https://www.pivotaltracker.com/story/show/162230780
		app.HeadTracker.Start(),

		app.Scheduler.Start(),
	)
}

// Stop allows the application to exit by halting schedules, closing
// logs, and closing the DB connection.
func (app *ChainlinkApplication) Stop() error {
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

		app.Scheduler.Stop()
		merr = multierr.Append(merr, app.HeadTracker.Stop())
		merr = multierr.Append(merr, app.JobSubscriber.Stop())
		app.FluxMonitor.Stop()
		merr = multierr.Append(merr, app.EthBroadcaster.Stop())
		app.RunQueue.Stop()
		merr = multierr.Append(merr, app.StatsPusher.Close())
		merr = multierr.Append(merr, app.SessionReaper.Stop())
		merr = multierr.Append(merr, app.Store.Close())
	})
	return merr
}

// GetStore returns the pointer to the store for the ChainlinkApplication.
func (app *ChainlinkApplication) GetStore() *strpkg.Store {
	return app.Store
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

// ArchiveJob silences the job from the system, preventing future job runs.
func (app *ChainlinkApplication) ArchiveJob(ID *models.ID) error {
	_ = app.JobSubscriber.RemoveJob(ID)
	app.FluxMonitor.RemoveJob(ID)
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

// NewBox returns the packr.Box instance that holds the static assets to
// be delivered by the router.
func (app *ChainlinkApplication) NewBox() packr.Box {
	return packr.NewBox("../../../operator_ui/dist")
}

type pendingConnectionResumer struct {
	runManager services.RunManager
}

func newPendingConnectionResumer(runManager services.RunManager) *pendingConnectionResumer {
	return &pendingConnectionResumer{runManager: runManager}
}

func (p *pendingConnectionResumer) Connect(head *models.Head) error {
	return p.runManager.ResumeAllPendingConnection()
}

func (p *pendingConnectionResumer) Disconnect()                   {}
func (p *pendingConnectionResumer) OnNewLongestChain(models.Head) {}
