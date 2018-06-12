package services

import (
	"fmt"
	"math/big"
	"strings"
	"sync"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/tidwall/gjson"
)

// OracleFulfillmentFunctionID is the function id of the oracle fulfillment
// method used by EthTx: bytes4(keccak256("fulfillData(uint256,bytes32)"))
// Kept in sync with solidity/contracts/Oracle.sol
const OracleFulfillmentFunctionID = "0x76005c26"

// Descriptive indices of a SpecAndRun's Topic array
const (
	SpecAndRunTopicSignature = iota
	SpecAndRunTopicInternalID
	SpecAndRunTopicAmount
)

// SpecAndRunTopic is the signature of the Oracle.sol's SpecAndRun event used for filtering.
// See https://github.com/smartcontractkit/chainlink/blob/master/solidity/contracts/Oracle.sol
var SpecAndRunTopic = common.HexToHash("0x40a86f3bd301164dcd67d63d081ecb2db540ac73bafb27eea27d65b3a2694f39")

// SpecAndRunSubscriber listens to push notifications from the ethereum node
// for JobSpec and Run requests originating on chain.
type SpecAndRunSubscriber struct {
	store         *store.Store
	subscription  *ManagedSubscription
	mutex         sync.Mutex
	listenAddress *common.Address
}

// NewSpecAndRunSubscriber creates a new instance of SpecAndRunSubscriber.
func NewSpecAndRunSubscriber(store *store.Store, listenAddress *common.Address) *SpecAndRunSubscriber {
	return &SpecAndRunSubscriber{store: store, listenAddress: listenAddress}
}

// Connect creates a subscription to SpecAndRunTopic to listen for new jobs.
func (scl *SpecAndRunSubscriber) Connect(head *models.IndexableBlockNumber) error {
	scl.mutex.Lock()
	defer scl.mutex.Unlock()

	topics := [][]common.Hash{{SpecAndRunTopic}}
	filter := ethereum.FilterQuery{
		FromBlock: head.NextInt(),
		Topics:    topics,
	}
	if scl.listenAddress != nil {
		log.Info("Only accepting SpecAndRun requests from contract at: %s", scl.listenAddress.String())
		filter.Addresses = []common.Address{*scl.listenAddress}
	}
	sub, err := NewManagedSubscription(scl.store, filter, scl.dispatchLog)
	scl.subscription = sub
	return err
}

// OnNewHead is a noop. Resuming of jobs based on new head activity aka confirmations,
// is handled by JobSubscriber, regardless of Initiator.
func (scl *SpecAndRunSubscriber) OnNewHead(*models.BlockHeader) {}

// Disconnect unsubscribes the SpecAndRunTopic subscription, stopping new jobs
// from on-chain.
func (scl *SpecAndRunSubscriber) Disconnect() {
	scl.mutex.Lock()
	defer scl.mutex.Unlock()

	if scl.subscription != nil {
		scl.subscription.Unsubscribe()
		scl.subscription = nil
	}
}

func (scl *SpecAndRunSubscriber) dispatchLog(log types.Log) {
	store := scl.store
	le, err := NewSpecAndRunLogEvent(log)
	if err != nil {
		logger.Errorw("SpecAndRunInitiator: Unable to retrieve spec and run parameters from log", le.forLogger("err", err)...)
		return
	}

	err = ValidateJob(le.Job, store)
	if err != nil {
		logger.Errorw("SpecAndRunInitiator: Invalid job from spec and run log", "err", err)
		return
	}

	err = store.SaveJob(&le.Job)
	if err != nil {
		logger.Errorw("SpecAndRunInitiator: Unable to save job", "err", err)
		return
	}

	go le.StartJob(store)
}

// SpecAndRunLogEvent encapsulates all information coming from the on-chain
// job and run request.
type SpecAndRunLogEvent struct {
	Log    types.Log
	Job    models.JobSpec
	Params models.JSON
}

// NewSpecAndRunLogEvent creates a new SpecAndRunLogEvent from a log,
// or returns an error if not possible.
func NewSpecAndRunLogEvent(log types.Log) (SpecAndRunLogEvent, error) {
	js, err := decodeABIToJSON(log.Data)
	if err != nil {
		return SpecAndRunLogEvent{}, err
	}

	job := models.NewJob()
	job.Initiators = []models.Initiator{{JobID: job.ID, Type: models.InitiatorSpecAndRun}}
	jsTasks := js.Get("tasks").Array()
	tasks := make([]models.TaskSpec, len(jsTasks))
	for i, tt := range jsTasks {
		tasks[i] = models.TaskSpec{
			Type:   tt.String(),
			Params: defaultParamsFor(tt.String(), log),
		}
	}
	job.Tasks = tasks

	return SpecAndRunLogEvent{
		Log:    log,
		Job:    job,
		Params: models.JSON{js.Get("params")},
	}, nil
}

func defaultParamsFor(task string, log types.Log) models.JSON {
	switch strings.ToLower(task) {
	case "ethtx":
		js := fmt.Sprintf(
			`{"address":"%s", "functionSelector":"%s", "dataPrefix":"%s"}`,
			log.Address.String(),
			OracleFulfillmentFunctionID,
			log.Topics[SpecAndRunTopicInternalID].String(),
		)
		return models.JSON{gjson.Parse(js)}
	default:
		return models.JSON{}
	}
}

// ToIndexableBlockNumber returns the IndexableBlockNumber associated with this
// log event.
func (le SpecAndRunLogEvent) ToIndexableBlockNumber() *models.IndexableBlockNumber {
	num := new(big.Int)
	num.SetUint64(le.Log.BlockNumber)
	return models.NewIndexableBlockNumber(num, le.Log.BlockHash)
}

// StartJob runs the job associated with the SpecAndRunLogEvent in job runner.
func (le SpecAndRunLogEvent) StartJob(store *store.Store) {
	jr, err := BuildRun(le.Job, le.Job.Initiators[0], store)
	if err != nil {
		logger.Error("SpecAndRunInitiator: unable to start job", err.Error())
	}

	rr := models.RunResult{Data: le.Params}
	if _, err := ExecuteRunAtBlock(jr, store, rr, le.ToIndexableBlockNumber()); err != nil {
		logger.Error("SpecAndRunInitiator: unable to start job", err.Error())
	}
}

// ForLogger formats the SpecAndRunLogEvent for easy formatting in logs.
func (le SpecAndRunLogEvent) forLogger(kvs ...interface{}) []interface{} {
	output := []interface{}{
		"job", le.Job,
		"log", le.Log,
		"params", le.Params,
	}

	return append(kvs, output...)
}
