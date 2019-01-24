package services

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/utils"
)

// LogRequest is the interface to allow polymorphic functionality of different
// types of LogEvents.
// i.e. EthLogEvent, RunLogEvent, ServiceAgreementLogEvent, OracleLogEvent
type LogRequest interface {
	GetLog() models.Log
	GetJobSpec() models.JobSpec
	GetInitiator() models.Initiator

	Validate() bool
	JSON() (models.JSON, error)
	ToDebug()
	ForLogger(kvs ...interface{}) []interface{}
	ContractPayment() (*assets.Link, error)
	ValidateRequester() error
	ToIndexableBlockNumber() *models.IndexableBlockNumber
}

// InitiatorLogEvent encapsulates all information as a result of a received log from an
// InitiatorSubscription.
type InitiatorLogEvent struct {
	Log       models.Log
	JobSpec   models.JobSpec
	Initiator models.Initiator
}

// LogRequest is a factory method that coerces this log event to the correct
// type based on Initiator.Type, exposed by the LogRequest interface.
func (le InitiatorLogEvent) LogRequest() LogRequest {
	switch le.Initiator.Type {
	case models.InitiatorServiceAgreementExecutionLog:
		fallthrough
	case models.InitiatorRunLog:
		return RunLogEvent{InitiatorLogEvent: le}
	case models.InitiatorEthLog:
		return EthLogEvent{InitiatorLogEvent: le}
	default:
		logger.Warnw("LogRequest: Unable to discern initiator type for log request", le.ForLogger()...)
		return EthLogEvent{InitiatorLogEvent: le}
	}
}

// GetLog returns the log.
func (le InitiatorLogEvent) GetLog() models.Log {
	return le.Log
}

// GetJobSpec returns the associated JobSpec
func (le InitiatorLogEvent) GetJobSpec() models.JobSpec {
	return le.JobSpec
}

// GetInitiator returns the initiator.
func (le InitiatorLogEvent) GetInitiator() models.Initiator {
	return le.Initiator
}

// ForLogger formats the InitiatorSubscriptionLogEvent for easy common formatting in logs (trace statements, not ethereum events).
func (le InitiatorLogEvent) ForLogger(kvs ...interface{}) []interface{} {
	output := []interface{}{
		"job", le.JobSpec.ID,
		"log", le.Log.BlockNumber,
		"initiator", le.Initiator,
	}

	return append(kvs, output...)
}

// ToDebug prints this event via logger.Debug.
func (le InitiatorLogEvent) ToDebug() {
	friendlyAddress := presenters.LogListeningAddress(le.Initiator.Address)
	msg := fmt.Sprintf("Received log from block #%v for address %v for job %v", le.Log.BlockNumber, friendlyAddress, le.JobSpec.ID)
	logger.Debugw(msg, le.ForLogger()...)
}

// ToIndexableBlockNumber returns an IndexableBlockNumber for the given InitiatorSubscriptionLogEvent Block
func (le InitiatorLogEvent) ToIndexableBlockNumber() *models.IndexableBlockNumber {
	num := new(big.Int)
	num.SetUint64(le.Log.BlockNumber)
	return models.NewIndexableBlockNumber(num, le.Log.BlockHash)
}

// Validate returns true, no validation on this log event type.
func (le InitiatorLogEvent) Validate() bool {
	return true
}

// ValidateRequester returns true since all requests are valid for base
// initiator log events.
func (le InitiatorLogEvent) ValidateRequester() error {
	return nil
}

// JSON returns the eth log as JSON.
func (le InitiatorLogEvent) JSON() (models.JSON, error) {
	el := le.Log
	var out models.JSON
	b, err := json.Marshal(el)
	if err != nil {
		return out, err
	}
	return out, json.Unmarshal(b, &out)
}

// ContractPayment returns the amount attached to a contract to pay the Oracle upon fulfillment.
func (le InitiatorLogEvent) ContractPayment() (*assets.Link, error) {
	return nil, nil
}

// EthLogEvent provides functionality specific to a log event emitted
// for an eth log initiator.
type EthLogEvent struct {
	InitiatorLogEvent
}

// RunLogEvent provides functionality specific to a log event emitted
// for a run log initiator.
type RunLogEvent struct {
	InitiatorLogEvent
}

// Validate returns whether or not the contained log has a properly encoded
// job id.
func (le RunLogEvent) Validate() bool {
	el := le.Log
	jid := jobIDFromHexEncodedTopic(el)
	if jid != le.JobSpec.ID && jobIDFromImproperEncodedTopic(el) != le.JobSpec.ID {
		logger.Errorw(fmt.Sprintf("Run Log didn't have matching job ID: %v != %v", jid, le.JobSpec.ID), le.ForLogger()...)
		return false
	}

	return true
}

// ContractPayment returns the amount attached to a contract to pay the Oracle upon fulfillment.
func (le RunLogEvent) ContractPayment() (*assets.Link, error) {
	encodedAmount := le.Log.Topics[models.RequestLogTopicAmount].Hex()
	payment, ok := new(assets.Link).SetString(encodedAmount, 0)
	if !ok {
		return payment, fmt.Errorf("unable to decoded amount from RunLog: %s", encodedAmount)
	}
	return payment, nil
}

// ValidateRequester returns true if the requester matches the one associated
// with the initiator.
func (le RunLogEvent) ValidateRequester() error {
	if len(le.Initiator.Requesters) == 0 {
		return nil
	}
	for _, r := range le.Initiator.Requesters {
		if le.Requester() == r {
			return nil
		}
	}
	return fmt.Errorf("Run Log didn't have have a valid requester: %v", le.Requester().Hex())
}

// Requester pulls the requesting address out of the LogEvent's topics.
func (le RunLogEvent) Requester() common.Address {
	b := le.Log.Topics[models.RequestLogTopicRequester].Bytes()
	return common.BytesToAddress(b)
}

// JSON decodes the CBOR in the ABI of the log event.
func (le RunLogEvent) JSON() (models.JSON, error) {
	el := le.Log
	js, err := decodeABIToJSON(el.Data)
	if err != nil {
		return js, err
	}

	fullfillmentJSON, err := fulfillmentToJSON(el)
	if err != nil {
		return js, err
	}
	return js.Merge(fullfillmentJSON)
}

func fulfillmentToJSON(el models.Log) (models.JSON, error) {
	var js models.JSON
	js, err := js.Add("address", el.Address.String())
	if err != nil {
		return js, err
	}

	js, err = js.Add("dataPrefix", encodeRequestID(el.Data))
	if err != nil {
		return js, err
	}

	return js.Add("functionSelector", models.OracleFulfillmentFunctionID)
}

func encodeRequestID(data []byte) string {
	return utils.AddHexPrefix(hex.EncodeToString(data[:common.HashLength]))
}

func decodeABIToJSON(data []byte) (models.JSON, error) {
	idSize := common.HashLength
	versionSize := common.HashLength
	varLocationSize := common.HashLength
	varLengthSize := common.HashLength
	start := idSize + versionSize + varLocationSize + varLengthSize
	return models.ParseCBOR(data[start:])
}

func jobIDFromHexEncodedTopic(log models.Log) string {
	return string(log.Topics[models.RequestLogTopicJobID].Bytes())
}

func jobIDFromImproperEncodedTopic(log models.Log) string {
	return log.Topics[models.RequestLogTopicJobID].String()[2:34]
}
