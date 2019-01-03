package services

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gobuffalo/packr"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"go.uber.org/multierr"
)

// Application implements the common functions used in the core node.
type Application interface {
	Start() error
	Stop() error
	OnConnect(func())
	GetStore() *store.Store
	WakeSessionReaper()
	WakeBulkRunDeleter()
	AddJob(job models.JobSpec) error
	AddAdapter(bt *models.BridgeType) error
	RemoveAdapter(bt *models.BridgeType) error
	NewBox() packr.Box
}

// ChainlinkApplication contains fields for the JobSubscriber, Scheduler,
// and Store. The JobSubscriber and Scheduler are also available
// in the services package, but the Store has its own package.
type ChainlinkApplication struct {
	Exiter                                            func(int)
	HeadTracker                                       *HeadTracker
	JobRunner                                         JobRunner
	JobSubscriber                                     JobSubscriber
	Scheduler                                         *Scheduler
	Store                                             *store.Store
	SessionReaper                                     SleeperTask
	BulkRunDeleter                                    SleeperTask
	pendingConnectionResumer                          *pendingConnectionResumer
	bridgeTypeMutex                                   sync.Mutex
	jobSubscriberID, txManagerID, connectionResumerID string
}

// NewApplication initializes a new store if one is not already
// present at the configured root directory (default: ~/.chainlink),
// the logger at the same directory and returns the Application to
// be used by the node.
func NewApplication(config store.Config) Application {
	store := store.NewStore(config)
	ht := NewHeadTracker(store)
	return &ChainlinkApplication{
		HeadTracker:              ht,
		JobSubscriber:            NewJobSubscriber(store),
		JobRunner:                NewJobRunner(store),
		Scheduler:                NewScheduler(store),
		Store:                    store,
		SessionReaper:            NewStoreReaper(store),
		BulkRunDeleter:           NewBulkRunDeleter(store),
		Exiter:                   os.Exit,
		pendingConnectionResumer: newPendingConnectionResumer(store),
	}
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

	app.txManagerID = app.HeadTracker.Attach(app.Store.TxManager)
	app.jobSubscriberID = app.HeadTracker.Attach(app.JobSubscriber)
	app.connectionResumerID = app.HeadTracker.Attach(app.pendingConnectionResumer)

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
		app.BulkRunDeleter.Start(),
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
	merr = multierr.Append(merr, app.BulkRunDeleter.Stop())
	app.HeadTracker.Detach(app.jobSubscriberID)
	app.HeadTracker.Detach(app.txManagerID)
	app.HeadTracker.Detach(app.connectionResumerID)
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

// WakeBulkRunDeleter wakes up the bulk task runner to process tasks.
func (app *ChainlinkApplication) WakeBulkRunDeleter() {
	app.BulkRunDeleter.WakeUp()
}

// AddJob adds a job to the store and the scheduler. If there was
// an error from adding the job to the store, the job will not be
// added to the scheduler.
func (app *ChainlinkApplication) AddJob(job models.JobSpec) error {
	err := app.Store.SaveJob(&job)
	if err != nil {
		return err
	}

	app.Scheduler.AddJob(job)
	return app.JobSubscriber.AddJob(job, nil) // nil for latest
}

// AddAdapter adds an adapter to the store. If another
// adapter with the same name already exists the adapter
// will not be added.
func (app *ChainlinkApplication) AddAdapter(bt *models.BridgeType) error {
	store := app.GetStore()

	bt.IncomingToken = utils.NewBytes32ID()
	bt.OutgoingToken = utils.NewBytes32ID()

	app.bridgeTypeMutex.Lock()
	defer app.bridgeTypeMutex.Unlock()

	if err := ValidateAdapter(bt, store); err != nil {
		return models.NewValidationError(err.Error())
	}

	if err := store.SaveBridgeType(bt); err != nil {
		return models.NewDatabaseAccessError(err.Error())
	}

	return nil
}

// RemoveAdapter removes an adapter from the store.
func (app *ChainlinkApplication) RemoveAdapter(bt *models.BridgeType) error {
	store := app.GetStore()

	app.bridgeTypeMutex.Lock()
	defer app.bridgeTypeMutex.Unlock()

	if err := store.DeleteStruct(bt); err != nil {
		return models.NewDatabaseAccessError(err.Error())
	}

	return nil
}

// NewBox returns the packr.Box instance that holds the static assets to
// be delivered by the router.
func (app *ChainlinkApplication) NewBox() packr.Box {
	return packr.NewBox("../gui/dist")
}

// OnConnect invokes the passed callback when connected to the block chain.
func (app *ChainlinkApplication) OnConnect(callback func()) {
	app.HeadTracker.Attach(newHeadTrackableCallback(callback))
}

type headTrackableCallback struct {
	onConnect func()
}

func newHeadTrackableCallback(callback func()) store.HeadTrackable {
	return &headTrackableCallback{onConnect: callback}
}

func (c *headTrackableCallback) Connect(*models.IndexableBlockNumber) error {
	c.onConnect()
	return nil
}

func (c *headTrackableCallback) Disconnect()                   {}
func (c *headTrackableCallback) OnNewHead(*models.BlockHeader) {}

type pendingConnectionResumer struct {
	store   *store.Store
	resumer func(*models.JobRun, *store.Store) (*models.JobRun, error)
}

func newPendingConnectionResumer(store *store.Store) *pendingConnectionResumer {
	return &pendingConnectionResumer{store: store, resumer: ResumeConnectingTask}
}

func (p *pendingConnectionResumer) Connect(head *models.IndexableBlockNumber) error {
	pendingRuns, err := p.store.JobRunsWithStatus(models.RunStatusPendingConnection)
	if err != nil {
		return multierr.Append(errors.New("error resuming pending connections"), err)
	}

	var merr error
	for _, jr := range pendingRuns {
		_, err := p.resumer(&jr, p.store)
		if err != nil {
			merr = multierr.Append(merr, err)
		}
	}
	return merr
}

func (p *pendingConnectionResumer) Disconnect()                   {}
func (p *pendingConnectionResumer) OnNewHead(*models.BlockHeader) {}
