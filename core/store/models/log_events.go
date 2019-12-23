package models

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"chainlink/core/assets"
	"chainlink/core/eth"
	"chainlink/core/logger"
	"chainlink/core/services/vrf"
	"chainlink/core/utils"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// Descriptive indices of a RunLog's Topic array
const (
	RequestLogTopicSignature = iota
	RequestLogTopicJobID
	RequestLogTopicRequester
	RequestLogTopicPayment
)

// Descriptive indices of the PrepaidAggregator's NewRound Topic Array:
// event NewRound(uint256 indexed roundId, address indexed startedBy);
const (
	NewRoundTopicSignature = iota
	NewRoundTopicRoundID
)

const (
	evmWordSize      = common.HashLength
	requesterSize    = evmWordSize
	idSize           = evmWordSize
	paymentSize      = evmWordSize
	versionSize      = evmWordSize
	callbackAddrSize = evmWordSize
	callbackFuncSize = evmWordSize
	expirationSize   = evmWordSize
	dataLocationSize = evmWordSize
	dataLengthSize   = evmWordSize
)

var (
	// RunLogTopic0original was the original topic to filter for Oracle.sol RunRequest events.
	RunLogTopic0original = utils.MustHash("RunRequest(bytes32,address,uint256,uint256,uint256,bytes)")
	// RunLogTopic20190123withFullfillmentParams was the new RunRequest filter topic as of 2019-01-23,
	// when callback address, callback function, and expiration were added to the data payload.
	RunLogTopic20190123withFullfillmentParams = utils.MustHash("RunRequest(bytes32,address,uint256,uint256,uint256,address,bytes4,uint256,bytes)")
	// RunLogTopic20190207withoutIndexes was the new RunRequest filter topic as of 2019-01-28,
	// after renaming Solidity variables, moving data version, and removing the cast of requestId to uint256
	RunLogTopic20190207withoutIndexes = utils.MustHash("OracleRequest(bytes32,address,bytes32,uint256,address,bytes4,uint256,uint256,bytes)")
	// ServiceAgreementExecutionLogTopic is the signature for the
	// Coordinator.RunRequest(...) events which Chainlink nodes watch for. See
	// https://chainlink/blob/master/evm/contracts/Coordinator.sol#RunRequest
	ServiceAgreementExecutionLogTopic = utils.MustHash("ServiceAgreementExecution(bytes32,address,uint256,uint256,uint256,bytes)")
	// RandomnessRequestLogTopic is the signature for the event log
	// VRFCoordinator.RandomnessRequest.
	RandomnessRequestLogTopic = vrf.RandomnessRequestLogTopic()
	// RandomnessRequestFullfillSelector is the signature for the method randomnessRequestFul
	RandomnessRequestFullfillSelector = vrf.FulfillSelector
	// OracleFullfillmentFunctionID0original is the original function selector for fulfilling Ethereum requests.
	OracleFullfillmentFunctionID0original = utils.MustHash("fulfillData(uint256,bytes32)").Hex()[:10]
	// OracleFulfillmentFunctionID20190123withFulfillmentParams is the function selector for fulfilling Ethereum requests,
	// as updated on 2019-01-23, accepting all fulfillment callback parameters.
	OracleFulfillmentFunctionID20190123withFulfillmentParams = utils.MustHash("fulfillData(uint256,uint256,address,bytes4,uint256,bytes32)").Hex()[:10]
	// OracleFulfillmentFunctionID20190128withoutCast is the function selector for fulfilling Ethereum requests,
	// as updated on 2019-01-28, removing the cast to uint256 for the requestId.
	OracleFulfillmentFunctionID20190128withoutCast = utils.MustHash("fulfillOracleRequest(bytes32,uint256,address,bytes4,uint256,bytes32)").Hex()[:10]
	// AggregatorNewRoundLogTopic20191220 is the NewRound filter topic for
	// the PrepaidAggregator as of Dec. 20th 2019. Eagerly fails if not found.
	AggregatorNewRoundLogTopic20191220 = eth.MustGetV5ContractEventID("PrepaidAggregator", "NewRound")
)

type logRequestParser interface {
	parseJSON(eth.Log) (JSON, error)
	parseRequestID(eth.Log) string
}

// topicFactoryMap maps the log topic to a factory method that returns an
// implementation of the interface LogRequest. The concrete implementations
// are polymorphic and can have difference behaviors for methods like JSON().
var topicFactoryMap = map[common.Hash]logRequestParser{
	ServiceAgreementExecutionLogTopic:         parseRunLog0original{},
	RunLogTopic0original:                      parseRunLog0original{},
	RunLogTopic20190123withFullfillmentParams: parseRunLog20190123withFulfillmentParams{},
	RunLogTopic20190207withoutIndexes:         parseRunLog20190207withoutIndexes{},
	RandomnessRequestLogTopic:                 parseRandomnessRequest{},
}

// TopicFiltersForRunLog generates the two variations of RunLog IDs that could
// possibly be entered on a RunLog or a ServiceAgreementExecutionLog. There is the ID,
// hex encoded and the ID zero padded.
func TopicFiltersForRunLog(logTopics []common.Hash, jobID *ID) [][]common.Hash {
	return [][]common.Hash{logTopics, {IDToTopic(jobID), IDToHexTopic(jobID)}}
}

// FilterQueryFactory returns the ethereum FilterQuery for this initiator.
func FilterQueryFactory(i Initiator, from *big.Int) (ethereum.FilterQuery, error) {
	q := ethereum.FilterQuery{
		FromBlock: from,
		Addresses: utils.WithoutZeroAddresses([]common.Address{i.Address}),
	}

	switch i.Type {
	case InitiatorEthLog:
		if from == nil {
			q.FromBlock = i.InitiatorParams.FromBlock.ToInt()
		} else if i.InitiatorParams.FromBlock != nil { // Know from != nil
			q.FromBlock = utils.MaxBigs(from, i.InitiatorParams.FromBlock.ToInt())
		}
		q.ToBlock = i.InitiatorParams.ToBlock.ToInt()

		if q.FromBlock != nil && q.ToBlock != nil && q.FromBlock.Cmp(q.ToBlock) >= 0 {
			return q, fmt.Errorf("cannot generate a FilterQuery with fromBlock >= toBlock")
		}

		q.Topics = make([][]common.Hash, len(i.Topics))
		copy(q.Topics, i.Topics) // Simply coercing i.Topics to the underlying type confuses reflect.DeepEqual

		return q, nil

	case InitiatorRunLog:
		topics := []common.Hash{RunLogTopic20190207withoutIndexes, RunLogTopic20190123withFullfillmentParams, RunLogTopic0original}
		q.Topics = TopicFiltersForRunLog(topics, i.JobSpecID)
		return q, nil

	case InitiatorServiceAgreementExecutionLog:
		topics := []common.Hash{ServiceAgreementExecutionLogTopic}
		q.Topics = TopicFiltersForRunLog(topics, i.JobSpecID)
		return q, nil
	case InitiatorFluxMonitor:
		q.Topics = [][]common.Hash{{AggregatorNewRoundLogTopic20191220}}
		return q, nil
	case InitiatorRandomnessLog:
		topics := []common.Hash{RandomnessRequestLogTopic}
		q.Topics = TopicFiltersForRunLog(topics, i.JobSpecID)
		return q, nil // XXX:  DRY this up
	default:
		return ethereum.FilterQuery{}, fmt.Errorf("Cannot generate a FilterQuery for initiator of type %T", i)
	}
}

// LogRequest is the interface to allow polymorphic functionality of different
// types of LogEvents.
// i.e. EthLogEvent, RunLogEvent, ServiceAgreementLogEvent, OracleLogEvent
type LogRequest interface {
	GetLog() eth.Log
	GetJobSpecID() *ID
	GetInitiator() Initiator

	Validate() bool
	JSON() (JSON, error)
	ToDebug()
	ForLogger(kvs ...interface{}) []interface{}
	ValidateRequester() error
	BlockNumber() *big.Int
	RunRequest() (RunRequest, error)
}

// InitiatorLogEvent encapsulates all information as a result of a received log from an
// InitiatorSubscription.
type InitiatorLogEvent struct {
	Log       eth.Log
	Initiator Initiator
}

var _ LogRequest = InitiatorLogEvent{}

// LogRequest is a factory method that coerces this log event to the correct
// type based on Initiator.Type, exposed by the LogRequest interface.
func (le InitiatorLogEvent) LogRequest() LogRequest {
	switch le.Initiator.Type {
	case InitiatorEthLog:
		return EthLogEvent{InitiatorLogEvent: le}
	case InitiatorServiceAgreementExecutionLog:
		fallthrough
	case InitiatorRunLog:
		return RunLogEvent{le}
	case InitiatorRandomnessLog:
		return RandomnessLogEvent{le}
	}
	logger.Warnw("LogRequest: Unable to discern initiator type for log request", le.ForLogger()...)
	return EthLogEvent{InitiatorLogEvent: le}
}

// GetLog returns the log.
func (le InitiatorLogEvent) GetLog() eth.Log {
	return le.Log
}

// GetJobSpecID returns the associated JobSpecID
func (le InitiatorLogEvent) GetJobSpecID() *ID {
	return le.Initiator.JobSpecID
}

// GetInitiator returns the initiator.
func (le InitiatorLogEvent) GetInitiator() Initiator {
	return le.Initiator
}

// ForLogger formats the InitiatorSubscriptionLogEvent for easy common
// formatting in logs (trace statements, not ethereum events).
func (le InitiatorLogEvent) ForLogger(kvs ...interface{}) []interface{} {
	output := []interface{}{
		"job", le.Initiator.JobSpecID.String(),
		"log", le.Log.BlockNumber,
		"initiator", le.Initiator,
	}
	for index, topic := range le.Log.Topics {
		output = append(output, fmt.Sprintf("topic%d", index), topic.Hex())
	}

	return append(kvs, output...)
}

// ToDebug prints this event via logger.Debug.
func (le InitiatorLogEvent) ToDebug() {
	friendlyAddress := utils.LogListeningAddress(le.Initiator.Address)
	logger.Debugw(
		fmt.Sprintf("Received log from block #%v for address %v", le.Log.BlockNumber, friendlyAddress),
		le.ForLogger()...,
	)
}

// BlockNumber returns the block number for the given InitiatorSubscriptionLogEvent.
func (le InitiatorLogEvent) BlockNumber() *big.Int {
	num := new(big.Int)
	num.SetUint64(le.Log.BlockNumber)
	return num
}

// RunRequest returns a run request instance with the transaction hash,
// present on all log initiated runs.
func (le InitiatorLogEvent) RunRequest() (RunRequest, error) {
	txHash := common.BytesToHash(le.Log.TxHash.Bytes())
	blockHash := common.BytesToHash(le.Log.BlockHash.Bytes())
	return RunRequest{
		BlockHash: &blockHash,
		TxHash:    &txHash,
	}, nil
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
func (le InitiatorLogEvent) JSON() (JSON, error) {
	el := le.Log
	var out JSON
	b, err := json.Marshal(el)
	if err != nil {
		return out, err
	}
	return out, json.Unmarshal(b, &out)
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
	jobSpecID := le.Initiator.JobSpecID
	topic := le.Log.Topics[RequestLogTopicJobID]

	if IDToTopic(jobSpecID) != topic && IDToHexTopic(jobSpecID) != topic {
		logger.Errorw("Run Log didn't have matching job ID", le.ForLogger("id", jobSpecID.String())...)
		return false
	}
	return true
}

// ContractPayment returns the amount attached to a contract to pay the Oracle upon fulfillment.
func contractPayment(log eth.Log) (*assets.Link, error) {
	version, err := log.GetTopic(0)
	if err != nil {
		return nil, fmt.Errorf("Missing RunLogEvent Topic#0: %v", err)
	}

	var encodedAmount common.Hash
	if oldRequestVersion(version) {
		encodedAmount = log.Topics[RequestLogTopicPayment]
	} else {
		paymentStart := requesterSize + idSize
		encodedAmount = common.BytesToHash(log.Data[paymentStart : paymentStart+paymentSize])
	}

	payment, ok := new(assets.Link).SetString(encodedAmount.Hex(), 0)
	if !ok {
		return payment, fmt.Errorf("unable to decoded amount from RunLog: %s", encodedAmount.Hex())
	}
	return payment, nil
}

func oldRequestVersion(version common.Hash) bool {
	return version == RunLogTopic0original ||
		version == RunLogTopic20190123withFullfillmentParams ||
		version == ServiceAgreementExecutionLogTopic
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
	version, err := le.Log.GetTopic(0)
	if err != nil {
		return common.Address{}
	}

	if oldRequestVersion(version) {
		return common.BytesToAddress(le.Log.Topics[RequestLogTopicRequester].Bytes())
	}
	return common.BytesToAddress(le.Log.Data[:requesterSize])
}

// RunRequest returns an RunRequest instance with all parameters
// from a run log topic, like RequestID.
func (le RunLogEvent) RunRequest() (RunRequest, error) {
	parser, err := parserFromLog(le.Log)
	if err != nil {
		return RunRequest{}, err
	}

	payment, err := contractPayment(le.Log)
	if err != nil {
		return RunRequest{}, err
	}

	txHash := common.BytesToHash(le.Log.TxHash.Bytes())
	blockHash := common.BytesToHash(le.Log.BlockHash.Bytes())
	str := parser.parseRequestID(le.Log)
	requester := le.Requester()
	return RunRequest{
		RequestID: &str,
		TxHash:    &txHash,
		BlockHash: &blockHash,
		Requester: &requester,
		Payment:   payment,
	}, nil
}

// JSON decodes the RunLogEvent's data converts it to a JSON object.
func (le RunLogEvent) JSON() (JSON, error) {
	return ParseRunLog(le.Log)
}

func parserFromLog(log eth.Log) (logRequestParser, error) {
	topic, err := log.GetTopic(0)
	if err != nil {
		return nil, errors.Wrap(err, "log#GetTopic(0)")
	}
	parser, ok := topicFactoryMap[topic]
	if !ok {
		return nil, fmt.Errorf("No parser for the RunLogEvent topic %s", topic.String())
	}
	return parser, nil
}

// ParseRunLog decodes the CBOR in the ABI of the log event.
func ParseRunLog(log eth.Log) (JSON, error) {
	parser, err := parserFromLog(log)
	if err != nil {
		return JSON{}, err
	}
	return parser.parseJSON(log)
}

func mustAdd(object JSON, key string, value string) JSON {
	rv, err := object.Add(key, value)
	if err != nil {
		panic(errors.Wrapf(err, "could not add %s under key %s in object %#+v",
			value, key, object))
	}
	return rv
}

// parseRunLog0original parses the original OracleRequest log format.
// It responds with only the request ID and data.
type parseRunLog0original struct{}

func (p parseRunLog0original) parseJSON(log eth.Log) (JSON, error) {
	data := log.Data
	start := idSize + versionSize + dataLocationSize + dataLengthSize

	js, err := ParseCBOR(data[start:])
	if err != nil {
		return js, err
	}
	js = mustAdd(js, "address", log.Address.String())
	js = mustAdd(js, "dataPrefix", bytesToHex(data[:idSize]))
	return mustAdd(js, "functionSelector", OracleFullfillmentFunctionID0original), nil
}

func (parseRunLog0original) parseRequestID(log eth.Log) string {
	return common.BytesToHash(log.Data[:idSize]).Hex()
}

// parseRunLog20190123withFulfillmentParams parses the OracleRequest log format
// which includes the callback, the payment amount, and expiration time
// The fulfillment also includes the callback, payment amount, and expiration,
// in addition to the request ID and data.
type parseRunLog20190123withFulfillmentParams struct{}

func (parseRunLog20190123withFulfillmentParams) parseJSON(log eth.Log) (JSON, error) {
	data := log.Data
	cborStart := idSize + versionSize + callbackAddrSize + callbackFuncSize + expirationSize + dataLocationSize + dataLengthSize

	js, err := ParseCBOR(data[cborStart:])
	if err != nil {
		return js, err
	}
	js = mustAdd(js, "address", log.Address.String())
	callbackAndExpStart := idSize + versionSize
	callbackAndExpEnd := callbackAndExpStart + callbackAddrSize + callbackFuncSize + expirationSize
	dataPrefix := bytesToHex(append(append(data[:idSize],
		log.Topics[RequestLogTopicPayment].Bytes()...),
		data[callbackAndExpStart:callbackAndExpEnd]...))
	js = mustAdd(js, "dataPrefix", dataPrefix)
	return mustAdd(js, "functionSelector",
		OracleFulfillmentFunctionID20190123withFulfillmentParams), nil
}

func (parseRunLog20190123withFulfillmentParams) parseRequestID(log eth.Log) string {
	return common.BytesToHash(log.Data[:idSize]).Hex()
}

// parseRunLog20190207withoutIndexes parses the OracleRequest log format after
// the sender and payment amount indexes were removed.
// Additionally, the version field for the data payload was moved next to the
// data that it corresponds to. The fulfillment is made up of the request ID,
// payment amount, callback, expiration, and data.
type parseRunLog20190207withoutIndexes struct{}

func (parseRunLog20190207withoutIndexes) parseJSON(log eth.Log) (JSON, error) {
	data := log.Data
	idStart := requesterSize
	expirationEnd := idStart + idSize + paymentSize + callbackAddrSize + callbackFuncSize + expirationSize
	cborStart := expirationEnd + versionSize + dataLocationSize + dataLengthSize
	js, err := ParseCBOR(data[cborStart:])
	if err != nil {
		return js, fmt.Errorf("Error parsing CBOR: %v", err)
	}
	js = mustAdd(js, "address", log.Address.String())
	js = mustAdd(js, "dataPrefix", bytesToHex(data[idStart:expirationEnd]))
	return mustAdd(js, "functionSelector", OracleFulfillmentFunctionID20190128withoutCast), nil
}

func (parseRunLog20190207withoutIndexes) parseRequestID(log eth.Log) string {
	start := requesterSize
	return common.BytesToHash(log.Data[start : start+idSize]).Hex()
}

// ParseNewRoundLog pulls the round from the aggregator log event.
func ParseNewRoundLog(log eth.Log) (*big.Int, error) {
	encodedRound := log.Topics[NewRoundTopicRoundID]
	round, ok := new(big.Int).SetString(encodedRound.Hex(), 0)
	if !ok {
		return nil, fmt.Errorf("unable to parse new round log from %s", encodedRound.Hex())
	}
	return round, nil
}

func bytesToHex(data []byte) string {
	return utils.AddHexPrefix(hex.EncodeToString(data))
}

// IDToTopic encodes the bytes representation of the ID padded to fit into a
// bytes32
func IDToTopic(id *ID) common.Hash {
	return common.BytesToHash(common.RightPadBytes(id.Bytes(), utils.EVMWordByteLen))
}

// IDToHexTopic encodes the string representation of the ID
func IDToHexTopic(id *ID) common.Hash {
	return common.BytesToHash([]byte(id.String()))
}
