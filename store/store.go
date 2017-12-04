package store

import (
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/smartcontractkit/chainlink-go/scheduler"
)

type Store struct {
	models.ORM
	Scheduler *scheduler.Scheduler
}

func New() Store {
	orm := models.InitORM("production")
	return Store{
		ORM:       orm,
		Scheduler: scheduler.New(orm),
	}
}

func (self Store) Start() error {
	return self.Scheduler.Start()
}

func (self Store) Close() {
	self.ORM.Close()
	self.Scheduler.Stop()
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
