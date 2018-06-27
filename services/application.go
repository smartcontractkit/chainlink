package services

import (
	"os"
	"os/signal"
	"sync"
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

// ChainlinkApplication contains fields for the JobSubscriber, Scheduler,
// and Store. The JobSubscriber and Scheduler are also available
// in the services package, but the Store has its own package.
type ChainlinkApplication struct {
	HeadTracker          *HeadTracker
	JobSubscriber        JobSubscriber
	RunManager           *RunManager
	Scheduler            *Scheduler
	Store                *store.Store
	Exiter               func(int)
	jobSubscriberID      string
	specAndRunSubscriber *SpecAndRunSubscriber
	specSubscriberID     string
	bridgeTypeMutex      sync.Mutex
}

// NewApplication initializes a new store if one is not already
// present at the configured root directory (default: ~/.chainlink),
// the logger at the same directory and returns the Application to
// be used by the node.
func NewApplication(config store.Config) Application {
	store := store.NewStore(config)
	ht := NewHeadTracker(store)
	return &ChainlinkApplication{
		HeadTracker:          ht,
		JobSubscriber:        NewJobSubscriber(store),
		RunManager:           NewRunManager(store),
		Scheduler:            NewScheduler(store),
		Store:                store,
		Exiter:               os.Exit,
		specAndRunSubscriber: NewSpecAndRunSubscriber(store, config.OracleContractAddress),
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
	app.specSubscriberID = app.HeadTracker.Attach(app.specAndRunSubscriber)

	return multierr.Combine(app.Store.Start(),
		app.HeadTracker.Start(),
		app.Scheduler.Start(),
		app.RunManager.Start())
}

// Stop allows the application to exit by halting schedules, closing
// logs, and closing the DB connection.
func (app *ChainlinkApplication) Stop() error {
	defer logger.Sync()
	defer app.RunManager.Stop()
	logger.Info("Gracefully exiting...")
	app.JobSubscriber.Stop()
	app.Scheduler.Stop()
	app.HeadTracker.Stop()
	app.HeadTracker.Detach(app.jobSubscriberID)
	app.HeadTracker.Detach(app.specSubscriberID)
	return app.Store.Stop()
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
	return app.JobSubscriber.AddJob(job, app.HeadTracker.LastRecord())
}

// AddAdapter adds an adapter to the store. If another
// adapter with the same name already exists the adapter
// will not be added.
func (app *ChainlinkApplication) AddAdapter(bt *models.BridgeType) error {
	store := app.GetStore()

	app.bridgeTypeMutex.Lock()
	defer app.bridgeTypeMutex.Unlock()

	if err := ValidateAdapter(bt, store); err != nil {
		return NewValidationError(err.Error())
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

// RunManager safely handles coordinating job runs.
type RunManager struct {
	store   *store.Store
	Workers map[string]chan store.RunRequest
}

// NewRunManager initializes a RunManager.
func NewRunManager(str *store.Store) *RunManager {
	return &RunManager{
		store:   str,
		Workers: make(map[string]chan store.RunRequest),
	}
}

func (rm *RunManager) Start() error {
	go rm.executeRunQueue()
	if err := rm.ResumeSleepingRuns(); err != nil {
		return err
	}

	return nil
}

func (rm *RunManager) ResumeSleepingRuns() error {
	pendingRuns, err := rm.store.JobRunsWithStatus(models.RunStatusPendingSleep)
	if err != nil {
		return err
	}
	for _, run := range pendingRuns {
		rm.store.RunQueue <- store.RunRequest{
			Input: models.RunResult{JobRunID: run.ID},
		}
	}
	return nil
}

func (rm *RunManager) executeRunQueue() {
	for rr := range rm.store.RunQueue {
		go func(rr store.RunRequest) {
			jr, err := rm.store.FindJobRun(rr.Input.JobRunID)
			if err != nil {
				logger.Warn("Application Run Channel Executor: error finding run", "ID", rr.Input.JobRunID)
				return
			}

			rm.WorkerChannelFor(jr) <- rr
		}(rr)
	}
	logger.Debug("Run Manager Closed")
}

// WorkerChannelFor accepts a JobRun and returns a worker channel dedicated
// to that JobRun. The channel accepts new block heights for triggering runs,
// and ensures that the block height confirmations are run syncronously.
func (rm *RunManager) WorkerChannelFor(jr models.JobRun) chan store.RunRequest {
	workerChannel, present := rm.Workers[jr.ID]
	if !present {
		workerChannel = make(chan store.RunRequest, 1000)
		rm.Workers[jr.ID] = workerChannel

		go func() {
			for rr := range workerChannel {
				logger.Debug("Woke up", jr.ID, "worker to process", rr.BlockNumber.ToInt())
				jr, err := ExecuteRunAtBlock(jr, rm.store, rr.Input, rr.BlockNumber)
				if err != nil {
					logger.Warnw("Application Run Channel Executor: error executing run", jr.ForLogger("error", err)...)
				}

				if jr.Status.Completed() {
					delete(rm.Workers, jr.ID)
					close(workerChannel)
				}
			}
			logger.Debug("Stopped worker for", jr.ID)
		}()
	}
	return workerChannel
}

// Stop closes all workers that have been started to process Job Runs on new
// heads and waits for them to finish.
func (rm *RunManager) Stop() {
	for _, workerChannel := range rm.Workers {
		close(workerChannel)
	}
}
