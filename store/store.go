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
