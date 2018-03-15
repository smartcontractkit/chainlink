package services

import (
	"os"
	"os/signal"
	"syscall"

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
	HeadTracker      *HeadTracker
	EthereumListener *EthereumListener
	Scheduler        *Scheduler
	Store            *store.Store
	Exiter           func(int)
	attachmentID     string
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
		HeadTracker:      ht,
		EthereumListener: &EthereumListener{Store: store},
		Scheduler:        NewScheduler(store),
		Store:            store,
		Exiter:           os.Exit,
	}
}

// Start runs the EthereumListener and Scheduler. If successful,
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
		app.Exiter(1)
	}()

	app.attachmentID = app.HeadTracker.Attach(app.EthereumListener)
	return multierr.Combine(app.HeadTracker.Start(), app.Scheduler.Start())
}

// Stop allows the application to exit by halting schedules, closing
// logs, and closing the DB connection.
func (app *ChainlinkApplication) Stop() error {
	defer logger.Sync()
	logger.Info("Gracefully exiting...")
	app.Scheduler.Stop()
	app.HeadTracker.Stop()
	app.HeadTracker.Detach(app.attachmentID)
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
	return app.EthereumListener.AddJob(job, app.HeadTracker.LastRecord())
}
