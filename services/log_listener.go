package services

import (
	"fmt"

	"github.com/asdine/storm/q"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

type LogListener struct {
	Store *store.Store
	logs  chan store.EventLog
}

func (self *LogListener) Start() error {
	jobs, err := self.Store.Jobs()
	if err != nil {
		return err
	}

	self.logs = make(chan store.EventLog)
	go self.listenToLogs()
	for _, j := range jobs {
		self.AddJob(j)
	}
	return nil
}

func (self *LogListener) Stop() error {
	if self.logs != nil {
		close(self.logs) // this will crash if called without start
	}
	return nil
}

func (self *LogListener) AddJob(job models.Job) error {
	for _, initr := range job.InitiatorsFor("ethLog") {
		address := initr.Address.String()
		if err := self.Store.Eth.Subscribe(self.logs, address); err != nil {
			return err
		}
	}
	return nil
}

func (self *LogListener) listenToLogs() {
	for l := range self.logs {
		for _, initr := range self.initrsWithLogAndAddress(l.Address) {
			if job, err := self.Store.FindJob(initr.JobID); err != nil {
				msg := fmt.Sprintf("Initiating job from log: %v", err)
				logger.Errorw(msg, "job", initr.JobID, "initiator", initr.ID)
			} else {
				StartJob(job.NewRun(), self.Store)
			}
		}
	}
}

func (self *LogListener) initrsWithLogAndAddress(address common.Address) []models.Initiator {
	initrs := []models.Initiator{}
	query := self.Store.Select(q.And(
		q.Eq("Address", address),
		q.Re("Type", "ethLog"),
	))
	if err := query.Find(&initrs); err != nil {
		msg := fmt.Sprintf("Initiating job from log: %v", err)
		logger.Errorw(msg, "address", address.String())
	}
	return initrs
}
