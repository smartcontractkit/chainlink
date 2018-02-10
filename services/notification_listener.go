package services

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/asdine/storm/q"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

const (
	EventTopicSignature = iota
	EventTopicNonce
	EventTopicJobID
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
		nl.AddJob(&j)
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

// AddJob looks for "chainlinklog" and "ethlog" Initiators for a given job
// and watches the Ethereum blockchain for the addresses in the job.
func (nl *NotificationListener) AddJob(job *models.Job) error {
	for _, initr := range job.InitiatorsFor(models.InitiatorEthLog, models.InitiatorChainlinkLog) {
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
			job, err := nl.Store.FindJob(initr.JobID)
			if err != nil {
				msg := fmt.Sprintf("Initiating job from log: %v", err)
				logger.Errorw(msg, "job", initr.JobID, "initiator", initr.ID)
				continue
			}

			input, err := FormatLogOutput(initr, el)
			if err != nil {
				logger.Errorw(err.Error(), "job", initr.JobID, "initiator", initr.ID)
				continue
			}

			BeginRun(job, nl.Store, input)
		}
	}
}

// FormatLogOutput uses the Initiator to decide how to format the EventLog
// as an Output object.
func FormatLogOutput(initr models.Initiator, el store.EventLog) (models.JSON, error) {
	if initr.Type == models.InitiatorEthLog {
		return convertEventLogToOutput(el)
	} else if initr.Type == models.InitiatorChainlinkLog {
		out, err := parseEventLogJSON(el)
		return out, err
	}
	return models.JSON{}, fmt.Errorf("no supported initiator type was found")
}

func convertEventLogToOutput(el store.EventLog) (models.JSON, error) {
	var out models.JSON
	b, err := json.Marshal(el)
	if err != nil {
		return out, err
	}
	return out, json.Unmarshal(b, &out)
}

func parseEventLogJSON(el store.EventLog) (models.JSON, error) {
	js, err := decodeABIToJSON(el.Data)
	if err != nil {
		return js, err
	}

	js, err = js.Add("address", el.Address.String())
	if err != nil {
		return js, err
	}

	js, err = js.Add("dataPrefix", el.Topics[EventTopicNonce].String())
	if err != nil {
		return js, err
	}

	return js.Add("functionId", "76005c26")
}

func (nl *NotificationListener) initrsWithLogAndAddress(address common.Address) []models.Initiator {
	initrs := []models.Initiator{}
	query := nl.Store.Select(q.Or(
		q.And(q.Eq("Address", address), q.Re("Type", models.InitiatorChainlinkLog)),
		q.And(q.Eq("Address", address), q.Re("Type", models.InitiatorEthLog)),
	))
	if err := query.Find(&initrs); err != nil {
		msg := fmt.Sprintf("Initiating job from log: %v", err)
		logger.Errorw(msg, "address", address.String())
	}
	return initrs
}

func decodeABIToJSON(data hexutil.Bytes) (models.JSON, error) {
	varLocationSize := 32
	varLengthSize := 32
	var js models.JSON
	hex := []byte(string([]byte(data)[varLocationSize+varLengthSize:]))
	return js, json.Unmarshal(bytes.TrimRight(hex, "\x00"), &js)
}
