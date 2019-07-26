package services

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/gobuffalo/packr"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"go.uber.org/multierr"
)

// Application implements the common functions used in the core node.
type Application interface {
	Start() error
	Stop() error
	GetStore() *store.Store
	WakeSessionReaper()
	AddJob(job models.JobSpec) error
	ArchiveJob(ID string) error
	AddServiceAgreement(*models.ServiceAgreement) error
	NewBox() packr.Box
}

// ChainlinkApplication contains fields for the JobSubscriber, Scheduler,
// and Store. The JobSubscriber and Scheduler are also available
// in the services package, but the Store has its own package.
type ChainlinkApplication struct {
	Exiter                   func(int)
	HeadTracker              *HeadTracker
	JobRunner                JobRunner
	JobSubscriber            JobSubscriber
	Scheduler                *Scheduler
	Store                    *store.Store
	SessionReaper            SleeperTask
	pendingConnectionResumer *pendingConnectionResumer
}

// NewApplication initializes a new store if one is not already
// present at the configured root directory (default: ~/.chainlink),
// the logger at the same directory and returns the Application to
// be used by the node.
func NewApplication(config orm.Depot, onConnectCallbacks ...func(Application)) Application {
	store := store.NewStore(config)

	jobSubscriber := NewJobSubscriber(store)
	pendingConnectionResumer := newPendingConnectionResumer(store)

	app := &ChainlinkApplication{
		JobSubscriber:            jobSubscriber,
		JobRunner:                NewJobRunner(store),
		Scheduler:                NewScheduler(store),
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

		// Deliberately started immediately after Store, to start the RunChannel consumer
		app.JobRunner.Start(),
		app.JobRunner.resumeRunsSinceLastShutdown(), // Started before any other service writes RunStatus to db.

		// HeadTracker deliberately started after JobRunner#resumeRunsSinceLastShutdown
		// since it Connects JobSubscriber which leads to writes of JobRuns RunStatus to the db.
		// https://www.pivotaltracker.com/story/show/162230780
		app.HeadTracker.Start(),
		app.Scheduler.Start(),
		app.SessionReaper.Start(),
	)
}

// Stop allows the application to exit by halting schedules, closing
// logs, and closing the DB connection.
func (app *ChainlinkApplication) Stop() error {
	defer logger.Sync()
	logger.Info("Gracefully exiting...")

	var merr error
	app.Scheduler.Stop()
	merr = multierr.Append(merr, app.HeadTracker.Stop())
	app.JobRunner.Stop()
	merr = multierr.Append(merr, app.SessionReaper.Stop())
	return multierr.Append(merr, app.Store.Close())
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
func (app *ChainlinkApplication) ArchiveJob(ID string) error {
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
	store   *store.Store
	resumer func(*models.JobRun, *store.Store) error
}

func newPendingConnectionResumer(store *store.Store) *pendingConnectionResumer {
	return &pendingConnectionResumer{store: store, resumer: ResumeConnectingTask}
}

func (p *pendingConnectionResumer) Connect(head *models.Head) error {
	var merr error
	err := p.store.UnscopedJobRunsWithStatus(func(run *models.JobRun) {
		err := p.resumer(run, p.store.Unscoped())
		if err != nil {
			merr = multierr.Append(merr, err)
		}
	}, models.RunStatusPendingConnection)

	if err != nil {
		return multierr.Append(errors.New("error resuming pending connections"), err)
	}

	return merr
}

func (p *pendingConnectionResumer) Disconnect()            {}
func (p *pendingConnectionResumer) OnNewHead(*models.Head) {}
