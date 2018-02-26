package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/utils"
)

const (
	EventTopicSignature = iota
	EventTopicRequestID
	EventTopicJobID
)

// RunLogTopic is the signature for the Request(uint256,bytes32,string) event
// which Chainlink RunLog initiators watch for.
// See https://github.com/smartcontractkit/chainlink/blob/master/solidity/contracts/Oracle.sol
var RunLogTopic = common.HexToHash("0x06f4bf36b4e011a5c499cef1113c2d166800ce4013f6c2509cab1a0e92b83fb2")

// Listens to event logs being pushed from the Ethereum Node specific to this job.
type Subscription struct {
	Job              models.Job
	store            *store.Store
	logNotifications chan types.Log
	errors           chan error
	rpcSubscription  *rpc.ClientSubscription
}

// Constructor of Subscription that to starts listening to and keeps track of
// event logs corresponding to a job.
func StartSubscription(job models.Job, store *store.Store) (Subscription, error) {
	var addresses []common.Address
	for _, initr := range job.InitiatorsFor(models.InitiatorEthLog, models.InitiatorRunLog) {
		msg := fmt.Sprintf("Listening for logs from address %v", presenters.LogListeningAddress(initr.Address))
		logger.Debugw(msg)
		addresses = append(addresses, initr.Address)
	}

	if len(addresses) == 0 {
		return Subscription{}, errors.New("Job must have a log initiator")
	}

	sub := Subscription{Job: job, store: store}
	sub.errors = make(chan error)
	sub.logNotifications = make(chan types.Log)

	fq := utils.ToFilterQueryFor(store.HeadTracker.Get().ToInt(), addresses)
	rpc, err := store.TxManager.SubscribeToLogs(sub.logNotifications, fq)
	if err != nil {
		return sub, err
	}
	sub.rpcSubscription = rpc
	go sub.listenToSubscriptionErrors()
	go sub.listenToLogs()
	return sub, nil
}

// Close channels and clean up resources.
func (sub Subscription) Unsubscribe() {
	if sub.rpcSubscription != nil && sub.rpcSubscription.Err() != nil {
		sub.rpcSubscription.Unsubscribe()
	}
	close(sub.logNotifications)
	close(sub.errors)
}

func (sub Subscription) listenToSubscriptionErrors() {
	for err := range sub.errors {
		logger.Errorw("Error in log subscription", "err", err)
	}
}

func (sub Subscription) listenToLogs() {
	for el := range sub.logNotifications {
		sub.receiveLog(el)
	}
}

func (sub Subscription) receiveLog(el types.Log) {
	msg := fmt.Sprintf("Received log from %v for job %v", el.Address.String(), sub.Job.ID)
	logger.Debugw(msg, "log", el, "job", sub.Job)

	for _, initr := range sub.Job.Initiators {
		if !sub.validateLog(initr.Type, el) {
			continue
		}
		data, err := FormatLogJSON(initr, el)
		if err != nil {
			logger.Errorw(err.Error(), "job", initr.JobID, "initiator", initr.ID)
			continue
		}

		input := models.RunResult{Data: data}
		if _, err = BeginRun(sub.Job, sub.store, input); err != nil {
			logger.Errorw(err.Error(), "job", initr.JobID, "initiator", initr.ID)
		}
	}
}

func (sub Subscription) validateLog(logType string, el types.Log) bool {
	switch logType {
	case models.InitiatorRunLog:
		return isRunLog(el) && sub.validateRunLog(el)
	}
	return true
}

func (sub Subscription) validateRunLog(el types.Log) bool {
	jid, err := jobIDFromLog(el)
	if err != nil {
		logger.Warnw("Failed to retrieve Job ID from log", "err", err.Error(), "jobID", sub.Job.ID)
		return false
	} else if jid != sub.Job.ID {
		logger.Warnw(fmt.Sprintf("Run Log didn't have matching job ID: %v != %v", jid, sub.Job.ID))
		return false
	}
	return true
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

func decodeABIToJSON(data hexutil.Bytes) (models.JSON, error) {
	varLocationSize := 32
	varLengthSize := 32
	var js models.JSON
	hex := []byte(string([]byte(data)[varLocationSize+varLengthSize:]))
	return js, json.Unmarshal(bytes.TrimRight(hex, "\x00"), &js)
}

func isRunLog(log types.Log) bool {
	return len(log.Topics) == 3 && log.Topics[0] == RunLogTopic
}

func jobIDFromLog(log types.Log) (string, error) {
	return utils.HexToString(log.Topics[EventTopicJobID].Hex())
}
