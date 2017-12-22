package services

import (
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/store/models"
)

type Application struct {
	Scheduler *Scheduler
	Store     *store.Store
}

func NewApplication(config store.Config) *Application {
	store := store.NewStore(config)
	logger.SetLoggerDir(config.RootDir)
	return &Application{
		Scheduler: NewScheduler(store),
		Store:     store,
	}
}

func (self *Application) Start() error {
	self.Store.Start()
	return self.Scheduler.Start()
}

func (self *Application) Stop() error {
	defer logger.Sync()
	logger.Info("Gracefully exiting...")
	self.Scheduler.Stop()
	return self.Store.Close()
}

func (self *Application) AddJob(job models.Job) error {
	err := self.Store.SaveJob(job)
	if err != nil {
		return err
	}

	self.Scheduler.AddJob(job)
	return nil
}
