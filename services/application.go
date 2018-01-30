package services

import (
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Application implements the common functions used in the core node.
type Application interface {
	Start() error
	Stop() error
	GetStore() *store.Store
}

// ChainlinkApplication contains fields for the LogListener, Scheduler,
// and Store. The LogListener and Scheduler are also available
// in the services package, but the Store has its own package.
type ChainlinkApplication struct {
	LogListener *LogListener
	Scheduler   *Scheduler
	Store       *store.Store
}

// NewApplication initializes a new store if one is not already
// present at the configured root directory (default: ~/.chainlink),
// the logger at the same directory and returns the Application to
// be used by the node.
func NewApplication(config store.Config) Application {
	store := store.NewStore(config)
	logger.SetLoggerDir(config.RootDir)
	return &ChainlinkApplication{
		LogListener: &LogListener{Store: store},
		Scheduler:   NewScheduler(store),
		Store:       store,
	}
}

// Start runs the Store, LogListener, and Scheduler. If successful,
// nil will be returned.
func (app *ChainlinkApplication) Start() error {
	app.Store.Start()
	app.LogListener.Start()
	return app.Scheduler.Start()
}

// Stop allows the application to exit by halting schedules, closing
// logs, and closing the DB connection.
func (app *ChainlinkApplication) Stop() error {
	defer logger.Sync()
	logger.Info("Gracefully exiting...")
	app.Scheduler.Stop()
	app.LogListener.Stop()
	return app.Store.Close()
}

// GetStore returns the pointer to the store for the ChainLinkApplication.
func (app *ChainlinkApplication) GetStore() *store.Store {
	return app.Store
}

// AddJob adds a job to the store and the scheduler. If there was
// an error from adding the job to the store, the job will not be
// added to the scheduler.
func (app *ChainlinkApplication) AddJob(job *models.Job) error {
	err := app.Store.SaveJob(job)
	if err != nil {
		return err
	}

	app.Scheduler.AddJob(job)
	return nil
}
