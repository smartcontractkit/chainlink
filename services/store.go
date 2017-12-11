package services

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/models"
)

type Store struct {
	*models.ORM
	Scheduler *Scheduler
	Config    Config
	KeyStore  *KeyStore
	sigs      chan os.Signal
}

func NewStore(config Config) *Store {
	orm := models.NewORM(config.RootDir)
	return &Store{
		ORM:       orm,
		Scheduler: NewScheduler(orm),
		Config:    config,
		KeyStore:  NewKeyStore(config.KeysDir()),
	}
}

func (self *Store) Start() error {
	self.sigs = make(chan os.Signal, 1)
	signal.Notify(self.sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-self.sigs
		self.Close()
		os.Exit(0)
	}()
	return self.Scheduler.Start()
}

func (self *Store) Close() {
	logger.Info("Gracefully exiting...")
	self.Scheduler.Stop()
	self.ORM.Close()
}

func (self *Store) AddJob(job models.Job) error {
	err := job.Validate()
	if err != nil {
		return err
	}

	err = self.Save(&job)
	if err != nil {
		return err
	}

	self.Scheduler.AddJob(job)
	return nil
}

func (self *Store) JobRunsFor(job models.Job) ([]models.JobRun, error) {
	var runs []models.JobRun
	err := self.Where("JobID", job.ID, &runs)
	return runs, err
}
