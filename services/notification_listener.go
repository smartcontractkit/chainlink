package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/asdine/storm/q"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"go.uber.org/multierr"
)

const (
	EventTopicSignature = iota
	EventTopicRequestID
	EventTopicJobID
)

// NotificationListener contains fields for the pointer of the store and
// a channel to the EthNotification (as the field 'logs').
type NotificationListener struct {
	Store             *store.Store
	subscriptions     []*rpc.ClientSubscription
	logNotifications  chan []types.Log
	headNotifications chan types.Header
	errors            chan error
	mutex             sync.Mutex
}

// Start obtains the jobs from the store and begins execution
// of the jobs' given runs.
func (nl *NotificationListener) Start() error {
	jobs, err := nl.Store.Jobs()
	if err != nil {
		return err
	}

	nl.errors = make(chan error)
	nl.logNotifications = make(chan []types.Log)
	nl.headNotifications = make(chan types.Header)
	var merr error
	for _, j := range jobs {
		merr = multierr.Append(merr, nl.AddJob(&j))
	}
	if merr != nil {
		return merr
	}

	go nl.listenToSubscriptionErrors()
	go nl.listenToLogs()
	go nl.listenToNewHeads()
	return nil
}

// Stop gracefully closes its access to the store's EthNotifications.
func (nl *NotificationListener) Stop() error {
	if nl.logNotifications != nil {
		nl.unsubscribe()
		close(nl.errors)
		close(nl.logNotifications)
	}
	return nil
}

// AddJob looks for "runlog" and "ethlog" Initiators for a given job
// and watches the Ethereum blockchain for the addresses in the job.
func (nl *NotificationListener) AddJob(job *models.Job) error {
	var addresses []common.Address
	for _, initr := range job.InitiatorsFor(models.InitiatorEthLog, models.InitiatorRunLog) {
		logger.Debugw(fmt.Sprintf("Listening for logs from address %v", initr.Address.String()))
		addresses = append(addresses, initr.Address)
	}

	if len(addresses) == 0 {
		return nil
	}

	sub, err := nl.Store.TxManager.SubscribeToLogs(nl.logNotifications, addresses)
	if err != nil {
		return err
	}
	nl.addSubscription(sub)
	return nil
}

func (nl *NotificationListener) listenToNewHeads() error {
	sub, err := nl.Store.TxManager.SubscribeToNewHeads(nl.headNotifications)
	if err != nil {
		return err
	}
	nl.addSubscription(sub)
	return nil
}

func (nl *NotificationListener) addSubscription(sub *rpc.ClientSubscription) {
	nl.mutex.Lock()
	defer nl.mutex.Unlock()
	nl.subscriptions = append(nl.subscriptions, sub)
	go func() {
		nl.errors <- (<-sub.Err())
	}()
}

func (nl *NotificationListener) unsubscribe() {
	nl.mutex.Lock()
	defer nl.mutex.Unlock()
	for _, sub := range nl.subscriptions {
		if sub.Err() != nil {
			sub.Unsubscribe()
		}
	}
}

func (nl *NotificationListener) listenToSubscriptionErrors() {
	for err := range nl.errors {
		logger.Errorw("Error in log subscription", "err", err)
	}
}

func (nl *NotificationListener) listenToLogs() {
	for logs := range nl.logNotifications {
		for _, el := range logs {
			if err := nl.receiveLog(el); err != nil {
				logger.Errorw(err.Error())
			}
		}
	}
}

func (nl *NotificationListener) receiveLog(el types.Log) error {
	var merr error
	msg := fmt.Sprintf("Received log from %v", el.Address.String())
	logger.Debugw(msg, "log", el)
	for _, initr := range nl.initrsWithLogAndAddress(el.Address) {
		job, err := nl.Store.FindJob(initr.JobID)
		if err != nil {
			msg := fmt.Sprintf("Error initiating job from log: %v", err)
			logger.Errorw(msg, "job", initr.JobID, "initiator", initr.ID)
			merr = multierr.Append(merr, err)
			continue
		}

		input, err := FormatLogJSON(initr, el)
		if err != nil {
			logger.Errorw(err.Error(), "job", initr.JobID, "initiator", initr.ID)
			merr = multierr.Append(merr, err)
			continue
		}

		_, err = BeginRun(job, nl.Store, input)
		merr = multierr.Append(merr, err)
	}
	return merr
}

// FormatLogJSON uses the Initiator to decide how to format the EventLog
// as a JSON object.
func FormatLogJSON(initr models.Initiator, el types.Log) (models.JSON, error) {
	if initr.Type == models.InitiatorEthLog {
		return ethLogJSON(el)
	} else if initr.Type == models.InitiatorRunLog {
		out, err := runLogJSON(el)
		return out, err
	}
	return models.JSON{}, fmt.Errorf("no supported initiator type was found")
}

// make our own types.Log for better serialization
func ethLogJSON(el types.Log) (models.JSON, error) {
	var out models.JSON
	b, err := json.Marshal(el)
	if err != nil {
		return out, err
	}
	return out, json.Unmarshal(b, &out)
}

func runLogJSON(el types.Log) (models.JSON, error) {
	js, err := decodeABIToJSON(el.Data)
	if err != nil {
		return js, err
	}

	js, err = js.Add("address", el.Address.String())
	if err != nil {
		return js, err
	}

	js, err = js.Add("dataPrefix", el.Topics[EventTopicRequestID].String())
	if err != nil {
		return js, err
	}

	return js.Add("functionSelector", "76005c26")
}

func (nl *NotificationListener) initrsWithLogAndAddress(address common.Address) []models.Initiator {
	initrs := []models.Initiator{}
	query := nl.Store.Select(q.Or(
		q.And(q.Eq("Address", address), q.Re("Type", models.InitiatorRunLog)),
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
