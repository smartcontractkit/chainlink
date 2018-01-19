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

func (ll *LogListener) Start() error {
	jobs, err := ll.Store.Jobs()
	if err != nil {
		return err
	}

	ll.logs = make(chan store.EventLog)
	go ll.listenToLogs()
	for _, j := range jobs {
		ll.AddJob(j)
	}
	return nil
}

func (ll *LogListener) Stop() error {
	if ll.logs != nil {
		close(ll.logs)
	}
	return nil
}

func (ll *LogListener) AddJob(job *models.Job) error {
	for _, initr := range job.InitiatorsFor(models.InitiatorEthLog) {
		address := initr.Address.String()
		if err := ll.Store.TxManager.Subscribe(ll.logs, address); err != nil {
			return err
		}
	}
	return nil
}

func (ll *LogListener) listenToLogs() {
	for l := range ll.logs {
		for _, initr := range ll.initrsWithLogAndAddress(l.Address) {
			if job, err := ll.Store.FindJob(initr.JobID); err != nil {
				msg := fmt.Sprintf("Initiating job from log: %v", err)
				logger.Errorw(msg, "job", initr.JobID, "initiator", initr.ID)
			} else {
				BeginRun(job, ll.Store)
			}
		}
	}
}

func (ll *LogListener) initrsWithLogAndAddress(address common.Address) []models.Initiator {
	initrs := []models.Initiator{}
	query := ll.Store.Select(q.And(
		q.Eq("Address", address),
		q.Re("Type", models.InitiatorEthLog),
	))
	if err := query.Find(&initrs); err != nil {
		msg := fmt.Sprintf("Initiating job from log: %v", err)
		logger.Errorw(msg, "address", address.String())
	}
	return initrs
}
