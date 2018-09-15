package services

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	GetStore() *store.Store
}

// ChainlinkApplication contains fields for the JobSubscriber, Scheduler,
// and Store. The JobSubscriber and Scheduler are also available
// in the services package, but the Store has its own package.
type ChainlinkApplication struct {
	Exiter          func(int)
	HeadTracker     *HeadTracker
	JobMetrics      JobMetrics
	JobRunner       JobRunner
	JobSubscriber   JobSubscriber
	Scheduler       *Scheduler
	Store           *store.Store
	Reaper          Reaper
	bridgeTypeMutex sync.Mutex
	jobSubscriberID string
}

// NewApplication initializes a new store if one is not already
// present at the configured root directory (default: ~/.chainlink),
// the logger at the same directory and returns the Application to
// be used by the node.
func NewApplication(config store.Config) Application {
	store := store.NewStore(config)
	ht := NewHeadTracker(store)
	jm := NewJobMetrics(store)
	return &ChainlinkApplication{
		HeadTracker:   ht,
		JobMetrics:    jm,
		JobSubscriber: NewJobSubscriber(store),
		JobRunner:     NewJobRunner(store, jm),
		Scheduler:     NewScheduler(store),
		Store:         store,
		Reaper:        NewStoreReaper(store),
		Exiter:        os.Exit,
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

	app.jobSubscriberID = app.HeadTracker.Attach(app.JobSubscriber)

	return multierr.Combine(
		app.Store.Start(),
		app.HeadTracker.Start(),
		app.Scheduler.Start(),
		app.JobMetrics.Start(),
		app.JobRunner.Start(),
		app.Reaper.Start(),
	)
}

// Stop allows the application to exit by halting schedules, closing
// logs, and closing the DB connection.
func (app *ChainlinkApplication) Stop() error {
	defer logger.Sync()
	logger.Info("Gracefully exiting...")
	app.Scheduler.Stop()
	app.HeadTracker.Stop()
	app.JobRunner.Stop()
	app.Reaper.Stop()
	app.HeadTracker.Detach(app.jobSubscriberID)
	return app.Store.Close()
}

// GetStore returns the pointer to the store for the ChainlinkApplication.
func (app *ChainlinkApplication) GetStore() *store.Store {
	return app.Store
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
	app.JobMetrics.Add(job)
	return app.JobSubscriber.AddJob(job, app.HeadTracker.LastRecord())
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

	if err := store.Save(bt); err != nil {
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
