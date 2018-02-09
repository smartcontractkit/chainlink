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
	EventTopicNonce
	EventTopicJobID
)

// NotificationListener contains fields for the pointer of the store and
// a channel to the EthNotification (as the field 'logs').
type NotificationListener struct {
	Store        *store.Store
	logs         chan []types.Log
	subscription *rpc.ClientSubscription
}

// Start obtains the jobs from the store and begins execution
// of the jobs' given runs.
func (nl *NotificationListener) Start() error {
	jobs, err := nl.Store.Jobs()
	if err != nil {
		return err
	}

	nl.logs = make(chan []types.Log)
	go nl.listenToLogs()
	err = nil
	for _, j := range jobs {
		err = multierr.Append(err, nl.AddJob(&j))
	}
	return err
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
	var addresses []common.Address
	for _, initr := range job.InitiatorsFor(models.InitiatorEthLog, models.InitiatorChainlinkLog) {
		logger.Debugw(fmt.Sprintf("Listening for logs from address %v", initr.Address.String()))
		addresses = append(addresses, initr.Address)
	}

	sub, err := nl.Store.TxManager.Subscribe(nl.logs, addresses)
	if err != nil {
		return err
	}
	nl.subscription = sub
	go func() {
		select {
		case err := <-nl.subscription.Err():
			logger.Panic(err)
		}
	}()
	return nil
}

func (nl *NotificationListener) listenToLogs() {
	for {
		select {
		case l := <-nl.logs:
			fmt.Println("***", l)
			el := l[0]
			msg := fmt.Sprintf("Received log from %v", el.Address.String())
			logger.Debugw(msg, "log", el)
			for _, initr := range nl.initrsWithLogAndAddress(el.Address) {
				job, err := nl.Store.FindJob(initr.JobID)
				if err != nil {
					msg := fmt.Sprintf("Error initiating job from log: %v", err)
					logger.Errorw(msg, "job", initr.JobID, "initiator", initr.ID)
					continue
				}

				input, err := FormatLogOutput(initr, el)
				if err != nil {
					logger.Errorw(err.Error(), "job", initr.JobID, "initiator", initr.ID)
					continue
				}

				fmt.Println("**** blow up on begin run?", input)
				BeginRun(job, nl.Store, input)
			}
			//case err := <-nl.subscription.Err():
			//logger.Panic(err)
		}
	}
}

// FormatLogOutput uses the Initiator to decide how to format the EventLog
// as an Output object.
func FormatLogOutput(initr models.Initiator, el types.Log) (models.JSON, error) {
	if initr.Type == models.InitiatorEthLog {
		return convertEventLogToOutput(el)
	} else if initr.Type == models.InitiatorChainlinkLog {
		out, err := parseEventLogJSON(el)
		return out, err
	}
	return models.JSON{}, fmt.Errorf("no supported initiator type was found")
}

// make our own types.Log
func convertEventLogToOutput(el types.Log) (models.JSON, error) {
	var out models.JSON
	b, err := json.Marshal(el)
	if err != nil {
		return out, err
	}
	var middle map[string]interface{}
	err = json.Unmarshal(b, &middle)
	if err != nil {
		return out, err
	}

	delete(middle, "removed")
	b, err = json.Marshal(middle)
	if err != nil {
		return out, err
	}
	return out, json.Unmarshal(b, &out)
}

func parseEventLogJSON(el types.Log) (models.JSON, error) {
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

// https://medium.com/justforfunc/two-ways-of-merging-n-channels-in-go-43c0b57cd1de
func merge(cs ...<-chan error) <-chan error {
	out := make(chan error)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan error) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
