package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/utils"
	"go.uber.org/multierr"
)

// Descriptive indices of a RunLog's Topic array
const (
	EventTopicSignature = iota
	EventTopicRequestID
	EventTopicJobID
)

// RunLogTopic is the signature for the Request(uint256,bytes32,string) event
// which Chainlink RunLog initiators watch for.
// See https://github.com/smartcontractkit/chainlink/blob/master/solidity/contracts/Oracle.sol
var RunLogTopic = common.HexToHash("0x06f4bf36b4e011a5c499cef1113c2d166800ce4013f6c2509cab1a0e92b83fb2")

// Listens to event logs being pushed from the Ethereum Node specific to a job.
type JobSubscription struct {
	Job           models.Job
	unsubscribers []Unsubscriber
}

// Constructor of JobSubscription that to starts listening to and keeps track of
// event logs corresponding to a job.
func StartJobSubscription(job models.Job, head *models.IndexableBlockNumber, store *store.Store) (JobSubscription, error) {
	var merr error
	var initSubs []Unsubscriber
	for _, initr := range job.InitiatorsFor(models.InitiatorEthLog) {
		sub, err := StartEthLogSubscription(initr, job, head, store)
		merr = multierr.Append(merr, err)
		if err == nil {
			initSubs = append(initSubs, sub)
		}
	}

	for _, initr := range job.InitiatorsFor(models.InitiatorRunLog) {
		sub, err := StartRunLogSubscription(initr, job, head, store)
		merr = multierr.Append(merr, err)
		if err == nil {
			initSubs = append(initSubs, sub)
		}
	}

	if len(initSubs) == 0 {
		return JobSubscription{}, multierr.Append(merr, errors.New("Job must have a valid log initiator"))
	}

	js := JobSubscription{Job: job, unsubscribers: initSubs}
	return js, merr
}

// Stops the subscription and cleans up associated resources.
func (js JobSubscription) Unsubscribe() {
	for _, sub := range js.unsubscribers {
		sub.Unsubscribe()
	}
}

// Interface for all subscriptions made specific to a subscription.
type Unsubscriber interface {
	Unsubscribe()
}

// Encapsulates all functionality needed to wrap an ethereum subscription
// for use with a Chainlink Initiator. Initiator specific functionality is delegated
// to the ReceiveLog callback using a strategy pattern.
type RPCLogSubscription struct {
	Job             models.Job
	Initiator       models.Initiator
	ReceiveLog      func(RPCLogEvent)
	store           *store.Store
	logs            chan types.Log
	errors          chan error
	ethSubscription models.EthSubscription
}

// Create a new RPCLogSubscription that feeds received logs to the callback func parameter.
func NewRPCLogSubscription(
	initr models.Initiator,
	job models.Job,
	head *models.IndexableBlockNumber,
	store *store.Store,
	callback func(RPCLogEvent),
) (RPCLogSubscription, error) {
	sub := RPCLogSubscription{Job: job, Initiator: initr, store: store, ReceiveLog: callback}
	sub.errors = make(chan error)
	sub.logs = make(chan types.Log)

	logListening(initr, head)
	fq := utils.ToFilterQueryFor(head.ToInt(), []common.Address{initr.Address})
	rpc, err := store.TxManager.SubscribeToLogs(sub.logs, fq)
	if err != nil {
		return sub, err
	}
	sub.ethSubscription = rpc
	go sub.listenToSubscriptionErrors()
	go sub.listenToLogs()
	return sub, nil
}

// Close channels and clean up resources.
func (sub RPCLogSubscription) Unsubscribe() {
	if sub.ethSubscription != nil && sub.ethSubscription.Err() != nil {
		sub.ethSubscription.Unsubscribe()
	}
	close(sub.logs)
	close(sub.errors)
}

func (sub RPCLogSubscription) listenToSubscriptionErrors() {
	for err := range sub.errors {
		logger.Errorw(fmt.Sprintf("Error in log subscription for job %v", sub.Job.ID), "err", err, "initr", sub.Initiator)
	}
}

func (sub RPCLogSubscription) listenToLogs() {
	for el := range sub.logs {
		sub.ReceiveLog(RPCLogEvent{
			Job:       sub.Job,
			Initiator: sub.Initiator,
			Log:       el,
			store:     sub.store,
		})
	}
}

// Starts an RPCLogSubscription tailored for use with RunLogs.
func StartRunLogSubscription(initr models.Initiator, job models.Job, head *models.IndexableBlockNumber, store *store.Store) (Unsubscriber, error) {
	return NewRPCLogSubscription(initr, job, head, store, ReceiveRunLog)
}

// Starts an RPCLogSubscription tailored for use with EthLogs.
func StartEthLogSubscription(initr models.Initiator, job models.Job, head *models.IndexableBlockNumber, store *store.Store) (Unsubscriber, error) {
	return NewRPCLogSubscription(initr, job, head, store, ReceiveEthLog)
}

func logListening(initr models.Initiator, number *models.IndexableBlockNumber) {
	msg := fmt.Sprintf(
		"Listening for %v from address %v from %v for job %v",
		initr.Type,
		presenters.LogListeningAddress(initr.Address),
		number.FriendlyString(),
		initr.JobID)
	logger.Infow(msg)
}

// Parse the log and run the job specific to this initiator log event.
func ReceiveRunLog(le RPCLogEvent) {
	if !le.ValidateRunLog() {
		return
	}

	friendlyAddress := presenters.LogListeningAddress(le.Initiator.Address)
	msg := fmt.Sprintf("Received log for address %v for job %v", friendlyAddress, le.Job.ID)
	logger.Infow(msg, le.ForLogger()...)

	data, err := le.RunLogJSON()
	if err != nil {
		logger.Errorw(err.Error(), le.ForLogger()...)
		return
	}

	runJob(le, data)
}

// Parse the log and run the job specific to this initiator log event.
func ReceiveEthLog(le RPCLogEvent) {
	friendlyAddress := presenters.LogListeningAddress(le.Initiator.Address)
	msg := fmt.Sprintf("Received log for address %v for job %v", friendlyAddress, le.Job.ID)
	logger.Infow(msg, le.ForLogger()...)

	data, err := le.EthLogJSON()
	if err != nil {
		logger.Errorw(err.Error(), le.ForLogger()...)
		return
	}

	runJob(le, data)
}

func runJob(le RPCLogEvent, data models.JSON) {
	input := models.RunResult{Data: data}
	if _, err := BeginRun(le.Job, le.store, input); err != nil {
		logger.Errorw(err.Error(), le.ForLogger()...)
	}
}

// Encapsulates all information as a result of a received log from an
// RPCLogSubscription.
type RPCLogEvent struct {
	Log       types.Log
	Job       models.Job
	Initiator models.Initiator
	store     *store.Store
}

// ForLogger formats the RPCLogEvent for easy common formatting in logs (trace statements, not ethereum events).
func (le RPCLogEvent) ForLogger(kvs ...interface{}) []interface{} {
	output := []interface{}{
		"job", le.Job,
		"log", le.Log,
		"initiator", le.Initiator,
	}

	return append(kvs, output...)
}

// Return whether or not the contained log is a RunLog, a specific Chainlink event trigger
// from smart contracts.
func (le RPCLogEvent) ValidateRunLog() bool {
	el := le.Log
	if !isRunLog(el) {
		logger.Debugw("Skipping; Unable to retrieve runlog parameters from log", le.ForLogger()...)
		return false
	}

	jid, err := jobIDFromLog(el)
	if err != nil {
		logger.Warnw("Failed to retrieve Job ID from log", le.ForLogger("err", err.Error())...)
		return false
	} else if jid != le.Job.ID {
		logger.Warnw(fmt.Sprintf("Run Log didn't have matching job ID: %v != %v", jid, le.Job.ID), le.ForLogger()...)
		return false
	}
	return true
}

// Extract data from the log's topics and data specific to the format defined
// by RunLogs.
func (le RPCLogEvent) RunLogJSON() (models.JSON, error) {
	el := le.Log
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

// Reformat the log as JSON.
func (le RPCLogEvent) EthLogJSON() (models.JSON, error) {
	el := le.Log
	var out models.JSON
	b, err := json.Marshal(el)
	if err != nil {
		return out, err
	}
	return out, json.Unmarshal(b, &out)
}

func decodeABIToJSON(data hexutil.Bytes) (models.JSON, error) {
	varLocationSize := 32
	varLengthSize := 32
	hex := []byte(string([]byte(data)[varLocationSize+varLengthSize:]))
	return models.ParseJSON(bytes.TrimRight(hex, "\x00"))
}

func isRunLog(log types.Log) bool {
	return len(log.Topics) == 3 && log.Topics[0] == RunLogTopic
}

func jobIDFromLog(log types.Log) (string, error) {
	return utils.HexToString(log.Topics[EventTopicJobID].Hex())
}
