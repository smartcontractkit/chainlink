package services

import (
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"go.uber.org/multierr"
)

// Application implements the common functions used in the core node.
type Application interface {
	Start() error
	Stop() error
	GetStore() *store.Store
}

// ChainlinkApplication contains fields for the EthereumListener, Scheduler,
// and Store. The EthereumListener and Scheduler are also available
// in the services package, but the Store has its own package.
type ChainlinkApplication struct {
	HeadTracker  *HeadTracker
	EthereumListener *EthereumListener
	Scheduler    *Scheduler
	Store        *store.Store
}

// NewApplication initializes a new store if one is not already
// present at the configured root directory (default: ~/.chainlink),
// the logger at the same directory and returns the Application to
// be used by the node.
func NewApplication(config store.Config) Application {
	store := store.NewStore(config)
	logger.Reconfigure(config.RootDir, config.LogLevel.Level)
	ht := NewHeadTracker(store)
	return &ChainlinkApplication{
		HeadTracker:  ht,
		EthereumListener: &EthereumListener{Store: store, HeadTracker: ht},
		Scheduler:    NewScheduler(store),
		Store:        store,
	}
}

// Start runs the Store, EthereumListener, and Scheduler. If successful,
// nil will be returned.
func (app *ChainlinkApplication) Start() error {
	app.Store.Start()
	return multierr.Combine(
		app.HeadTracker.Start(),
		app.EthereumListener.Start(),
		app.Scheduler.Start())
}

// Stop allows the application to exit by halting schedules, closing
// logs, and closing the DB connection.
func (app *ChainlinkApplication) Stop() error {
	defer logger.Sync()
	logger.Info("Gracefully exiting...")
	app.Scheduler.Stop()
	app.EthereumListener.Stop()
	app.HeadTracker.Stop()
	return app.Store.Close()
}

// GetStore returns the pointer to the store for the ChainlinkApplication.
func (app *ChainlinkApplication) GetStore() *store.Store {
	return app.Store
}

// AddJob adds a job to the store and the scheduler. If there was
// an error from adding the job to the store, the job will not be
// added to the scheduler.
func (app *ChainlinkApplication) AddJob(job models.Job) error {
	err := app.Store.SaveJob(&job)
	if err != nil {
		return err
	}

	app.Scheduler.AddJob(job)
	return app.EthereumListener.AddJob(job)
}
