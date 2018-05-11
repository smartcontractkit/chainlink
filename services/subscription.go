package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
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
	EventTopicInternalID
	EventTopicJobID
	EventTopicAmount
)

// RunLogTopic is the signature for the RunRequest(...) event
// which Chainlink RunLog initiators watch for.
// See https://github.com/smartcontractkit/chainlink/blob/master/solidity/contracts/Oracle.sol
// If updating this, be sure to update the truffle suite's "expected event signature" test.
var RunLogTopic = common.HexToHash("0x3fab86a1207bdcfe3976d0d9df25f263d45ae8d381a60960559771a2b223974d")

// JobSubscription listens to event logs being pushed from the Ethereum Node to a job.
type JobSubscription struct {
	Job           models.JobSpec
	unsubscribers []Unsubscriber
}

// StartJobSubscription is the constructor of JobSubscription that to starts
// listening to and keeps track of event logs corresponding to a job.
func StartJobSubscription(job models.JobSpec, head *models.IndexableBlockNumber, store *store.Store) (JobSubscription, error) {
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

// Unsubscribe stops the subscription and cleans up associated resources.
func (js JobSubscription) Unsubscribe() {
	for _, sub := range js.unsubscribers {
		sub.Unsubscribe()
	}
}

// Unsubscriber is the interface for all subscriptions, allowing one to unsubscribe.
type Unsubscriber interface {
	Unsubscribe()
}

// RPCLogSubscription encapsulates all functionality needed to wrap an ethereum subscription
// for use with a Chainlink Initiator. Initiator specific functionality is delegated
// to the ReceiveLog callback using a strategy pattern.
type RPCLogSubscription struct {
	Job             models.JobSpec
	Initiator       models.Initiator
	ReceiveLog      func(RPCLogEvent)
	store           *store.Store
	logs            chan types.Log
	errors          chan error
	ethSubscription models.EthSubscription
}

// RPCLogSubscriber represents the function for processing a received RPCLog and the filter.
type RPCLogSubscriber struct {
	Filter   ethereum.FilterQuery
	Callback func(RPCLogEvent)
}

// NewRPCLogSubscription returns a new RPCLogSubscriber with initialized filter.
func NewRPCLogSubscriber(
	initr models.Initiator,
	head *models.IndexableBlockNumber,
	topics [][]common.Hash,
	callback func(RPCLogEvent)) *RPCLogSubscriber {

	listenFromNumber := head.NextInt()
	q := utils.ToFilterQueryFor(listenFromNumber, []common.Address{initr.Address})

	return &RPCLogSubscriber{
		Filter:   q,
		Callback: callback,
	}
}

// NewRPCLogSubscription creates a new RPCLogSubscription that feeds received
// logs to the callback func parameter.
func NewRPCLogSubscription(
	initr models.Initiator,
	job models.JobSpec,
	store *store.Store,
	subscriber *RPCLogSubscriber,
) (RPCLogSubscription, error) {
	if !initr.IsLogInitiated() {
		return RPCLogSubscription{}, errors.New("Can only create an RPC log subscription for log initiators")
	}

	sub := RPCLogSubscription{Job: job, Initiator: initr, store: store, ReceiveLog: subscriber.Callback}
	sub.errors = make(chan error)
	sub.logs = make(chan types.Log)

	loggerLogListening(initr, subscriber.Filter.FromBlock)
	es, err := store.TxManager.SubscribeToLogs(sub.logs, subscriber.Filter)
	if err != nil {
		return sub, err
	}

	sub.ethSubscription = es
	go sub.listenToSubscriptionErrors()
	go sub.listenToLogs(subscriber.Filter)
	return sub, nil
}

// Unsubscribe closes channels and clean up resources.
func (sub RPCLogSubscription) Unsubscribe() {
	if sub.ethSubscription != nil {
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

func (sub RPCLogSubscription) listenToLogs(q ethereum.FilterQuery) {
	backfilledSet := sub.backfillLogs(q)
	for el := range sub.logs {
		if _, present := backfilledSet[el.BlockHash.String()]; !present {
			sub.dispatchLog(el)
		}
	}
}

func (sub RPCLogSubscription) backfillLogs(q ethereum.FilterQuery) map[string]bool {
	backfilledSet := map[string]bool{}
	if q.FromBlock.Cmp(big.NewInt(0)) <= 0 {
		return backfilledSet
	}

	logs, err := sub.store.TxManager.GetLogs(q)
	if err != nil {
		logger.Errorw("Unable to backfill logs", "err", err)
		return backfilledSet
	}

	for _, log := range logs {
		sub.dispatchLog(log)
		backfilledSet[log.BlockHash.String()] = true
	}
	return backfilledSet
}

func (sub RPCLogSubscription) dispatchLog(log types.Log) {
	sub.ReceiveLog(RPCLogEvent{
		Job:       sub.Job,
		Initiator: sub.Initiator,
		Log:       log,
		store:     sub.store,
	})
}

// StartRunLogSubscription starts an RPCLogSubscription tailored for use with RunLogs.
func StartRunLogSubscription(initr models.Initiator, job models.JobSpec, head *models.IndexableBlockNumber, store *store.Store) (Unsubscriber, error) {
	hashJobID := common.BytesToHash([]byte(job.ID))
	topics := [][]common.Hash{{RunLogTopic}, {}, {hashJobID}}
	subscriber := NewRPCLogSubscriber(initr, head, topics, receiveRunLog)
	return NewRPCLogSubscription(initr, job, store, subscriber)
}

// StartEthLogSubscription starts an RPCLogSubscription tailored for use with EthLogs.
func StartEthLogSubscription(initr models.Initiator, job models.JobSpec, head *models.IndexableBlockNumber, store *store.Store) (Unsubscriber, error) {
	subscriber := NewRPCLogSubscriber(initr, head, nil, receiveEthLog)
	return NewRPCLogSubscription(initr, job, store, subscriber)
}

func loggerLogListening(initr models.Initiator, blockNumber *big.Int) {
	msg := fmt.Sprintf(
		"Listening for %v from block %v for address %v for job %v",
		initr.Type,
		presenters.FriendlyBigInt(blockNumber),
		presenters.LogListeningAddress(initr.Address),
		initr.JobID)
	logger.Infow(msg)
}

// Parse the log and run the job specific to this initiator log event.
func receiveRunLog(le RPCLogEvent) {
	if !le.ValidateRunLog() {
		return
	}

	le.ToDebug()
	data, err := le.RunLogJSON()
	if err != nil {
		logger.Errorw(err.Error(), le.ForLogger()...)
		return
	}

	runJob(le, data, le.Initiator)
}

// Parse the log and run the job specific to this initiator log event.
func receiveEthLog(le RPCLogEvent) {
	le.ToDebug()
	data, err := le.EthLogJSON()
	if err != nil {
		logger.Errorw(err.Error(), le.ForLogger()...)
		return
	}

	runJob(le, data, le.Initiator)
}

func runJob(le RPCLogEvent, data models.JSON, initr models.Initiator) {
	payment, err := le.ContractPayment()
	if err != nil {
		logger.Errorw(err.Error(), le.ForLogger()...)
		return
	}
	input := models.RunResult{
		Data:   data,
		Amount: payment,
	}
	if _, err := BeginRunAtBlock(le.Job, initr, input, le.store, le.ToIndexableBlockNumber()); err != nil {
		logger.Errorw(err.Error(), le.ForLogger()...)
	}
}

// RPCLogEvent encapsulates all information as a result of a received log from an
// RPCLogSubscription.
type RPCLogEvent struct {
	Log       types.Log
	Job       models.JobSpec
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

// ToDebug prints this event via logger.Debug.
func (le RPCLogEvent) ToDebug() {
	friendlyAddress := presenters.LogListeningAddress(le.Initiator.Address)
	msg := fmt.Sprintf("Received log from block #%v for address %v for job %v", le.Log.BlockNumber, friendlyAddress, le.Job.ID)
	logger.Debugw(msg, le.ForLogger()...)
}

func (le RPCLogEvent) ToIndexableBlockNumber() *models.IndexableBlockNumber {
	num := new(big.Int)
	num.SetUint64(le.Log.BlockNumber)
	return models.NewIndexableBlockNumber(num, le.Log.BlockHash)
}

// ValidateRunLog returns whether or not the contained log is a RunLog,
// a specific Chainlink event trigger from smart contracts.
func (le RPCLogEvent) ValidateRunLog() bool {
	el := le.Log
	if !isRunLog(el) {
		logger.Errorw("Skipping; Unable to retrieve runlog parameters from log", le.ForLogger()...)
		return false
	}

	jid, err := jobIDFromLog(el)
	if err != nil {
		logger.Errorw("Failed to retrieve Job ID from log", le.ForLogger("err", err.Error())...)
		return false
	} else if jid != le.Job.ID {
		logger.Errorw(fmt.Sprintf("Run Log didn't have matching job ID: %v != %v", jid, le.Job.ID), le.ForLogger()...)
		return false
	}
	return true
}

// fulfillDataFunctionID is the signature for the fulfillData(uint256,bytes32) function located in Oracle.sol
// fulfillDataFunctionID is calculated in the following way: bytes4(keccak256("fulfillData(uint256,bytes32)"))
// See https://github.com/smartcontractkit/chainlink/blob/master/solidity/contracts/Oracle.sol
var fulfillDataFunctionID = "76005c26"

// RunLogJSON extracts data from the log's topics and data specific to the format defined
// by RunLogs.
func (le RPCLogEvent) RunLogJSON() (models.JSON, error) {
	el := le.Log
	js, err := decodeABIToJSON(el.Data)
	if err != nil {
		return js, err
	}

	fullfillmentJSON, err := fulfillmentToJSON(le)
	if err != nil {
		return js, err
	}
	return js.Merge(fullfillmentJSON)
}

func fulfillmentToJSON(le RPCLogEvent) (models.JSON, error) {
	el := le.Log
	var js models.JSON
	js, err := js.Add("address", el.Address.String())
	if err != nil {
		return js, err
	}

	js, err = js.Add("dataPrefix", el.Topics[EventTopicInternalID].String())
	if err != nil {
		return js, err
	}

	return js.Add("functionSelector", fulfillDataFunctionID)
}

// EthLogJSON reformats the log as JSON.
func (le RPCLogEvent) EthLogJSON() (models.JSON, error) {
	el := le.Log
	var out models.JSON
	b, err := json.Marshal(el)
	if err != nil {
		return out, err
	}
	return out, json.Unmarshal(b, &out)
}

// ContractPayment returns the amount attached to a contract to pay the Oracle upon fulfillment.
func (le RPCLogEvent) ContractPayment() (*big.Int, error) {
	if !isRunLog(le.Log) {
		return nil, nil
	}
	encodedAmount := le.Log.Topics[EventTopicAmount].Hex()
	payment, ok := new(big.Int).SetString(encodedAmount, 0)
	if !ok {
		return payment, fmt.Errorf("unable to decoded amount from RunLog: %s", encodedAmount)
	}
	return payment, nil
}

func decodeABIToJSON(data hexutil.Bytes) (models.JSON, error) {
	versionSize := 32
	varLocationSize := 32
	varLengthSize := 32
	prefix := versionSize + varLocationSize + varLengthSize
	hex := []byte(string([]byte(data)[prefix:]))
	return models.ParseCBOR(hex)
}

func isRunLog(log types.Log) bool {
	return len(log.Topics) == 4 && log.Topics[0] == RunLogTopic
}

func jobIDFromLog(log types.Log) (string, error) {
	return utils.HexToString(log.Topics[EventTopicJobID].Hex())
}
