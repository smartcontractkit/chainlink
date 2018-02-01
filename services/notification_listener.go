package services

import (
	"fmt"

	"github.com/asdine/storm/q"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// NotificationListener contains fields for the pointer of the store and
// a channel to the EthNotification (as the field 'logs').
type NotificationListener struct {
	Store *store.Store
	logs  chan store.EthNotification
}

// Start obtains the jobs from the store and begins execution
// of the jobs' given runs.
func (nl *NotificationListener) Start() error {
	jobs, err := nl.Store.Jobs()
	if err != nil {
		return err
	}

	nl.logs = make(chan store.EthNotification)
	go nl.listenToLogs()
	for _, j := range jobs {
		nl.AddJob(j)
	}
	return nil
}

// Stop gracefully closes its access to the store's EthNotifications.
func (nl *NotificationListener) Stop() error {
	if nl.logs != nil {
		close(nl.logs)
	}
	return nil
}

// AddJob looks for "ethlog" Initiators for a given job and watches
// the Ethereum blockchain for the addresses in the job.
func (nl *NotificationListener) AddJob(job *models.Job) error {
	for _, initr := range job.InitiatorsFor(models.InitiatorEthLog) {
		address := initr.Address.String()
		if err := nl.Store.TxManager.Subscribe(nl.logs, address); err != nil {
			return err
		}
	}
	return nil
}

func (nl *NotificationListener) listenToLogs() {
	for l := range nl.logs {
		el, err := l.UnmarshalLog()
		if err != nil {
			logger.Errorw("Unable to unmarshal log", "log", l)
			continue
		}

		for _, initr := range nl.initrsWithLogAndAddress(el.Address) {
			if job, err := nl.Store.FindJob(initr.JobID); err != nil {
				msg := fmt.Sprintf("Initiating job from log: %v", err)
				logger.Errorw(msg, "job", initr.JobID, "initiator", initr.ID)
			} else {
				BeginRun(job, nl.Store)
			}
		}
	}
}

func (nl *NotificationListener) initrsWithLogAndAddress(address common.Address) []models.Initiator {
	initrs := []models.Initiator{}
	query := nl.Store.Select(q.And(
		q.Eq("Address", address),
		q.Re("Type", models.InitiatorEthLog),
	))
	if err := query.Find(&initrs); err != nil {
		msg := fmt.Sprintf("Initiating job from log: %v", err)
		logger.Errorw(msg, "address", address.String())
	}
	return initrs
}
