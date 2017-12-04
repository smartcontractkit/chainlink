package store

import (
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/smartcontractkit/chainlink-go/services"
)

type Store struct {
	models.ORM
	Scheduler *services.Scheduler
}

func New() Store {
	orm := models.InitORM("production")
	return Store{
		ORM:       orm,
		Scheduler: services.NewScheduler(orm),
	}
}

func (self Store) Start() error {
	return self.Scheduler.Start()
}

func (self Store) Close() {
	self.Scheduler.Stop()
	self.ORM.Close()
}

func (self Store) AddJob(job models.Job) error {
	err := self.Save(&job)
	if err != nil {
		return err
	}

	self.Scheduler.AddJob(job)
	return nil
}

func (self Store) JobRunsFor(job models.Job) ([]models.JobRun, error) {
	var runs []models.JobRun
	err := self.Where("JobID", job.ID, &runs)
	return runs, err
}
