package services

import (
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/smartcontractkit/chainlink-go/store"
)

type Application struct {
	Scheduler *Scheduler
	Store     *store.Store
}

func NewApplication(config store.Config) *Application {
	store := store.NewStore(config)
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
	logger.Info("Gracefully exiting...")
	self.Scheduler.Stop()
	return self.Store.Close()
}

func (self *Application) AddJob(job models.Job) error {
	err := self.Store.Save(&job)
	if err != nil {
		return err
	}

	self.Scheduler.AddJob(job)
	return nil
}
