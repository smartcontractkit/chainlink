package services

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gobuffalo/packr"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"go.uber.org/multierr"
)

//go:generate mockery -name Application -output ../internal/mocks/ -case=underscore

// Application implements the common functions used in the core node.
type Application interface {
	Start() error
	Stop() error
	GetStore() *store.Store
	WakeSessionReaper()
	AddJob(job models.JobSpec) error
	ArchiveJob(*models.ID) error
	AddServiceAgreement(*models.ServiceAgreement) error
	NewBox() packr.Box
	RunManager
}

// ChainlinkApplication contains fields for the JobSubscriber, Scheduler,
// and Store. The JobSubscriber and Scheduler are also available
// in the services package, but the Store has its own package.
type ChainlinkApplication struct {
	Exiter      func(int)
	HeadTracker *HeadTracker
	RunManager
	RunQueue                 RunQueue
	JobSubscriber            JobSubscriber
	Scheduler                *Scheduler
	Store                    *store.Store
	SessionReaper            SleeperTask
	pendingConnectionResumer *pendingConnectionResumer
	shutdownOnce             sync.Once
}

// NewApplication initializes a new store if one is not already
// present at the configured root directory (default: ~/.chainlink),
// the logger at the same directory and returns the Application to
// be used by the node.
func NewApplication(config *orm.Config, onConnectCallbacks ...func(Application)) Application {
	store := store.NewStore(config)
	config.SetRuntimeStore(store.ORM)

	runExecutor := NewRunExecutor(store)
	runQueue := NewRunQueue(runExecutor)
	runManager := NewRunManager(runQueue, config, store.ORM, store.TxManager, store.Clock)
	jobSubscriber := NewJobSubscriber(store, runManager)

	pendingConnectionResumer := newPendingConnectionResumer(runManager)

	app := &ChainlinkApplication{
		JobSubscriber:            jobSubscriber,
		RunManager:               runManager,
		RunQueue:                 runQueue,
		Scheduler:                NewScheduler(store, runManager),
		Store:                    store,
		SessionReaper:            NewStoreReaper(store),
		Exiter:                   os.Exit,
		pendingConnectionResumer: pendingConnectionResumer,
	}

	headTrackables := []strpkg.HeadTrackable{
		store.TxManager,
		jobSubscriber,
		pendingConnectionResumer,
	}
	for _, onConnectCallback := range onConnectCallbacks {
		headTrackable := &headTrackableCallback{func() {
			onConnectCallback(app)
		}}
		headTrackables = append(headTrackables, headTrackable)
	}
	app.HeadTracker = NewHeadTracker(store, headTrackables)

	return app
}

// Start runs the JobSubscriber and Scheduler. If successful,
// nil will be returned.
// Also listens for interrupt signals from the operating system so
// that the application can be properly closed before the application
// exits.
func (app *ChainlinkApplication) Start() error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		app.Stop()
		app.Exiter(0)
	}()

	return multierr.Combine(
		app.Store.Start(),
		app.RunQueue.Start(),
		app.RunManager.ResumeAllInProgress(),

		// HeadTracker deliberately started after
		// RunQueue#resumeRunsSinceLastShutdown since it Connects JobSubscriber
		// which leads to writes of JobRuns RunStatus to the db.
		// https://www.pivotaltracker.com/story/show/162230780
		app.HeadTracker.Start(),

		app.Scheduler.Start(),
		app.SessionReaper.Start(),
	)
}

// Stop allows the application to exit by halting schedules, closing
// logs, and closing the DB connection.
func (app *ChainlinkApplication) Stop() error {
	var merr error
	app.shutdownOnce.Do(func() {
		defer logger.Sync()
		logger.Info("Gracefully exiting...")

		app.Scheduler.Stop()
		merr = multierr.Append(merr, app.HeadTracker.Stop())
		app.RunQueue.Stop()
		merr = multierr.Append(merr, app.SessionReaper.Stop())
		merr = multierr.Append(merr, app.Store.Close())
	})
	return merr
}

// GetStore returns the pointer to the store for the ChainlinkApplication.
func (app *ChainlinkApplication) GetStore() *store.Store {
	return app.Store
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
	return app.JobSubscriber.AddJob(job, nil) // nil for latest
}

// ArchiveJob silences the job from the system, preventing future job runs.
func (app *ChainlinkApplication) ArchiveJob(ID *models.ID) error {
	_ = app.JobSubscriber.RemoveJob(ID)
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
	return app.JobSubscriber.AddJob(sa.JobSpec, nil) // nil for latest
}

// NewBox returns the packr.Box instance that holds the static assets to
// be delivered by the router.
func (app *ChainlinkApplication) NewBox() packr.Box {
	return packr.NewBox("../../operator_ui/dist")
}

type pendingConnectionResumer struct {
	runManager RunManager
}

func newPendingConnectionResumer(runManager RunManager) *pendingConnectionResumer {
	return &pendingConnectionResumer{runManager: runManager}
}

func (p *pendingConnectionResumer) Connect(head *models.Head) error {
	return p.runManager.ResumeAllConnecting()
}

func (p *pendingConnectionResumer) Disconnect()            {}
func (p *pendingConnectionResumer) OnNewHead(*models.Head) {}
