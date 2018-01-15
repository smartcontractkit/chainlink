package services

import (
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

type Application struct {
	LogListener *LogListener
	Scheduler   *Scheduler
	Store       *store.Store
}

func NewApplication(config store.Config) *Application {
	store := store.NewStore(config)
	logger.SetLoggerDir(config.RootDir)
	return &Application{
		LogListener: &LogListener{Store: store},
		Scheduler:   NewScheduler(store),
		Store:       store,
	}
}

func (app *Application) Start() error {
	app.Store.Start()
	app.LogListener.Start()
	return app.Scheduler.Start()
}

func (app *Application) Stop() error {
	defer logger.Sync()
	logger.Info("Gracefully exiting...")
	app.Scheduler.Stop()
	app.LogListener.Stop()
	return app.Store.Close()
}

func (app *Application) AddJob(job models.Job) error {
	err := app.Store.SaveJob(job)
	if err != nil {
		return err
	}

	app.Scheduler.AddJob(job)
	return nil
}
